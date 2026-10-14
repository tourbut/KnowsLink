---
title: SAR-REFACTOR-002 어댑터 감사와 Inbox 오류 복구
status: draft
updated: 2026-10-14
owner: dev
tasks: [SAR-REFACTOR-002-DEV]
summary: 읽음 처리 뒤 정리 오류가 본문을 잃게 하는 결함을 재현하고 읽기 전 검사로 수정했다.
---

# 어댑터 감사와 수정 결과

기준 SHA는 `8f9370bb62c9133c73e10ea2307b4cd696e10f7e`다.
Purpose는 implementation이며 테스트 레벨은 lite다. 이 하위 작업은 재위임하지 않았다.
동일 dev 워크트리에서 adapters와 두 패키징 스크립트만 감사했다.
Git·공유 기술 문서·전체 gate는 부모가 통합한다. 이 보고서는 하위 작업의 미커밋 변경 근거다.

## 감사 범위와 부하

KnowsLink 로컬 어댑터는 승인된 text를 검증하고 private Inbox에 저장한 뒤 relay persist/ACK를 보낸다.
MCP와 watcher는 같은 연결 폴더를 사용할 수 있다. 자동 pull과 수동 receive가 함께 일어나는 정상 부하를 전제로 했다.
자동 수신의 기본 주기는 10초이며 읽음 marker는 24시간 남는다.

직접 확인한 주요 경로:

- `adapters/src/core.ts`: HTTP 응답 상한·timeout·redirect·서명 입력과 synthetic 경로.
- `adapters/src/text.ts`: owner/agent 검증·TTL·저장 후 ACK·SSE parser·watch 및 CLI 수신.
- `adapters/src/inbox.ts`: store·take·중복 marker·wake claim·poll/SSE 직렬화·backoff.
- `adapters/src/mcp.ts`: startup·로그인·자동 receiver·metadata notice·상태와 수동 receive 호출자.
- `adapters/src/private-files.ts`, `connect.ts`, `login.ts`, `grok-wake.ts`: private 파일/ACL·새 연결·명시 동의·wake 결과 경계.
- `adapters/src/inbox.test.ts`, `events.test.ts`, `owner-text.test.ts`, `trial-boundaries.test.ts`: 관련 회귀와 정리 동작의 존재.
- `adapters/package.json`, `tsconfig.json`, host manifest와 MCP 설정: 고정 도구·held 기본값·설치별 경로.
- `scripts/package_plugin.py`: allowlist·standalone 묶음·private sibling 거부·cloud app 참조.
- `scripts/install_bot_mcp.sh`: Node checksum·staging·기존 경로 검사·설치 순서.

현재 `handovers/to_dev.md`, `FULLOPS.md`, common 0.3.3 세 규칙, delegation, project, dev context,
`docs/system-design.md`, `adapters/README.md`, 문서 작성 규칙을 확인했다.
Ponytail full·audit와 codebase-design 기준으로 새 계층이나 의존성을 추가하지 않았다.
부모의 packet·공유 문서와 Go 변경은 수정하지 않았다.

## 재현된 결함과 최소 수정

### 읽음 처리 뒤 housekeeping 오류가 미표시 본문을 잃게 함

기존 `Inbox.take()`는 `.json`을 `.read`로 claim하고 내용을 읽은 뒤 `.read`를 비웠다.
그 다음 `prune()`가 디렉터리를 순회했다.
이때 stat에 EACCES 또는 I/O 오류가 나면 take는 본문을 반환하지 않고 실패했다.
이미 본문을 비웠으므로 다음 receive도 해당 메시지를 반환하지 못했다.
또한 다른 receiver가 readdir와 stat 사이에서 tmp를 rename하거나 marker를 정리하면 ENOENT가 발생할 수 있었다.

합성 메시지 1개와 fresh tmp marker 1개로 재현했다.
실제 filesystem의 Inbox를 사용하고 marker stat에만 EACCES를 주입했다.
수정 전 take는 EACCES로 실패하고 pending이 `[]`가 되어 unread 보존 assertion이 실패했다.
반환되지 않은 본문은 원본 `.json`에서도 사라졌다.
프로덕션 파일·실계정·외부 발송은 사용하지 않았다.

수정은 다음 두 가지다.

1. unread 후보가 있으면 claim 전에 prune을 완료한다. prune 실패 시 unread 파일을 보존하고 오류를 전달한다.
2. prune의 stat ENOENT만 무시한다. 다른 오류는 그대로 전달한다.

새 테스트는 EACCES 후 원본 파일·본문 유지, ENOENT 후 같은 메시지 수신, 빈 Inbox, 읽은 ID의 중복 저장 거부를 검사한다.
`package.json`의 기존 test 체인에 이 회귀를 추가했다.
take 호출자는 MCP 수동 receive와 text CLI다. 공개 메서드·반환형·호출 순서는 바꾸지 않았다.
store-before-persist/ACK·서명·TTL·wake 상태·재시도 정책·MCP stdout 계약은 수정하지 않았다.

## 검사와 원본 종료코드

실행 위치는 dev 워크트리 루트다. Node는 다음 고정 실행 파일을 사용했다.

`C:/Users/shin/AppData/Local/KnowsLinkDevTools/node-v22.22.2-win-x64/node.exe`

기본 PATH의 Node는 v24.13.0이므로 검사에 사용하지 않았다.
아래 `node22`는 위 실행 파일을 나타내는 보고서 표기다.
셸은 각 명령의 `$LASTEXITCODE`를 그대로 반환했다. 로그 축약 파이프는 사용하지 않았다.

| 검사 | 명령 | 원본 exit | 결과 |
|---|---|---:|---|
| 수정 전 컴파일 | `node22 adapters/node_modules/typescript/bin/tsc --project adapters/tsconfig.json` | 0 | 새 회귀가 정상 컴파일됨 |
| 수정 전 회귀 | `node22 adapters/dist/inbox-files.test.js` | 1 | `actual: []`, expected는 unread `.json` 1개 |
| 수정 후 컴파일 | 위 tsc 명령 | 0 | 변경 코드와 회귀 컴파일 성공 |
| 수정 후 회귀 | `node22 adapters/dist/inbox-files.test.js` | 0 | housekeeping failure preserves text; missing marker permits receive; duplicate stays read |
| owner 수신 계약 | `node22 adapters/dist/owner-text.test.js` | 0 | authorship·대상·schema 거부·copy-before-ACK 통과 |
| 변경 파일 lint | `node22 adapters/node_modules/eslint/bin/eslint.js adapters/src/inbox.ts adapters/src/inbox-files.test.ts --config adapters/eslint.config.mjs` | 0 | 경고·오류 출력 없음 |
| 변경 파일 format | `node22 adapters/node_modules/prettier/bin/prettier.cjs --check adapters/src/inbox.ts adapters/src/inbox-files.test.ts adapters/package.json` | 0 | All matched files use Prettier code style |
| patch 공백 | `git diff --check` | 0 | 오류 없음 |

오류 주입에는 Node 기본 `mock.method`와 `syncBuiltinESMExports`를 사용했다.
Context7 resolve 요청은 `Monthly quota exceeded`로 실패하여 query를 실행할 library ID를 받지 못했다.
공식 Node v22.22.2의 [test 문서](https://github.com/nodejs/node/blob/v22.22.2/doc/api/test.md)와
[module 문서](https://github.com/nodejs/node/blob/v22.22.2/doc/api/module.md#modulesyncbuiltinesmexports)를 확인했다.
테스트는 stat 원본을 mock 전에 보관하고 finally에서 mock과 ESM 바인딩을 복원한다.

성능 수치는 수집하지 않았다. 이 수정은 오류 시 본문 보존을 고치며 순회 횟수나 자료구조를 최적화하지 않는다.
처리량 향상이나 큰 Inbox의 성능 보장을 주장하지 않는다. 부모의 Go/DB 성능 근거와 구분한다.

## 보존 자료와 실행 제약

- 수정 전 합성 fixture: `C:/Users/shin/AppData/Local/Temp/kl-inbox-files-pNIZsc`.
- 수정 후 합성 fixture: `C:/Users/shin/AppData/Local/Temp/kl-inbox-files-3U4oIj`.
- 컴파일 출력: `adapters/dist/`. 이 하위 작업은 파일이나 디렉터리를 정리하지 않았다.
- 테스트의 모든 marker는 fresh다. prune의 오래된 파일 제거 분기는 실행하지 않았다.
- 기존 통합 테스트와 패키지 verify는 임시 폴더 제거 코드를 포함한다. 이 하위 작업에서는 실행하지 않았다.
- read-only rg 검색의 삭제 API 문자열과 이를 포함한 보고 명령이 shell_guard에 거부됐다.
  차단된 삭제 실행은 없었다. 검색 표현을 우회 실행하지 않고 소스 직접 읽기를 계속했다.
  짧은 별도 escalation은 부모에게 정상 전달했다.

## 변경 파일과 부모 인계

- `adapters/src/inbox.ts`
- `adapters/src/inbox-files.test.ts` (신규)
- `adapters/package.json`
- 이 보고서

공유 문서 수정 제안: `docs/system-design.md`와 D05의 로컬 Inbox 절에
"읽음 claim 전에 housekeeping을 완료하고 다른 receiver가 없앤 marker의 ENOENT는 건너뛴다"를 추가한다.
저장/ACK·TTL·wake 설명과 제품 규칙은 그대로 유지한다. adapter README의 설치 계약은 변경하지 않아도 된다.

부모 후속: 전체 lint/test/build·FullOps gate, Git 커밋, packet-outcomes와 공유 기술 문서, 독립 고정 SHA 리뷰와 QA.
독립 검토·실계정·실제 호스트 UI·운영 배포·장시간/다중 프로세스 스트레스·패키지 설치 실행은 이 보고서의 통과 범위가 아니다.
프로세스가 읽음 claim 뒤 종료되는 기존 crash window와 payload read/write 자체의 I/O 실패는 이번 수정 범위가 아니다.
패키징·설정·설치 스크립트에서는 수정이 필요한 결함을 별도로 재현하지 않아 변경하지 않았다.
