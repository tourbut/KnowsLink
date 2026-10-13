---
title: Sol 리팩토링 독립 경계 QA
status: review
updated: 2026-10-13
owner: tester
tasks: [SAR-REFACTOR-001-TESTER]
summary: 고정 cd4dae8의 원본 게이트를 대조하고 Linux에서 연결 응답·사람 선택·text·MCP 경계를 독립 검증했다.
---

# SAR-REFACTOR-001-TESTER 독립 QA

판정은 이번 기능 보존 리팩토링 범위에서 PASS다. 새 제품 결함과 미해결 critical/high 발견은 0건이다. 제품 코드는 변경하지 않았다. 전체 운영 수락과 coordinator의 독립 코드 리뷰는 별도다.

## 대상과 독립성

- 기준 diff: `0eb7735f0963272b62266992791f9f09152f9266` → `cd4dae810ae262759bcff29895165f1a64f8c51a`.
- 준비 HEAD: `eea94316418cfd8ccc501e6ca66ed82a2abc3565`. 제품 경로는 후보와 동일하다. DEV 최종 문서 `e65146b`도 제품이 동일하다.
- tester 실제 세션: `01a13534-4add-7a22-bb3f-5e40c8addb09`. Codex `gpt-6.1-sol` / high. Task `task_bf0eee88d603`, Dispatch `ctx_6e75a377f476`.
- 원 구현 부모 `01a13509-e1a3-7b62-bf9d-5a48405d706a`, Go `01a1350c-1383-7323-9d75-748304b33792`, TS `01a1350c-7f6b-7f03-822b-e3a3ec516f30`, 완료 재개 `01a1351f-c007-7f22-ac7c-1837dca51f0c`와 다른 세션이다.
- 적용 기준: `fullops-common-0.3.3`, AGENTS.md, FULLOPS.md, project.md, 공통 세 규칙, D02 SAR-HUMAN-PEER, D10 module-design, DEV 실행 기록과 QA 원본. Test level lite. 좁은 QA는 직접 수행했다.

## 원본 전체 게이트 재사용

[gate-audit.json](gate-audit.json)은 DEV의 원본 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/refactor001/cd4dae810ae262759bcff29895165f1a64f8c51a-<target>.log`와 Git의 [DEV gate](../SAR-REFACTOR-001-DEV-test/gate/)를 대조한다. install/lint/test/build의 내용은 LF 정규화 뒤 동일하다. 원본과 사본의 SHA256을 각각 남겼다. Windows checkout의 개행 변환 때문에 바이트 동일성을 주장하지 않는다.

| 검사 | 원본 결과 | 재사용 근거 |
|---|---|---|
| `make install` | EXIT=0 | 고정 후보의 lockfile 설치 |
| `make lint` | EXIT=0 | format/vet/modules, adapter format/ESLint/types, 설정·공개 ingress 검사 |
| `make test` | EXIT=0 | Go race, TS 등록 suite, shell guard 4 tests. Go 결과 일부는 정상 캐시다. |
| `make build` | EXIT=0 | Go relay/migrate, TypeScript, plugin bundle |
| FullOps 1.4.1 등록 게이트 | product-lint/product-test exit 0 | JSON head가 정확한 후보다. ERROR 0, WARNING 3, unavailable 0이다. |

DEV FullOps 기준은 `bc0fd031efd252d9755681d8b429eb31c045b7e4`다. tester 기준 `0eb7735f0963272b62266992791f9f09152f9266`와 제품 경로가 같음을 확인했다. 경고는 PLANS/core/human_test 크기 3건이다. 원래 사라진 1.4.0 경로의 실패는 DEV 기록에 그대로 둔다. lint.log의 public config 거절 문구는 예상 부정 검사이며 전체 EXIT=0과 함께 판정했다. 재사용 결과를 tester의 새 전체 게이트 실행으로 표시하지 않는다.

## 이번 독립 실행

실행 위치는 기존 고유 clone `/tmp/knowslink-refactor001-cd4dae810ae2-45997998/repo`다. 검증 전후 HEAD는 정확한 후보이며 `git status --porcelain`은 비어 있다. Node `v22.22.2`, npm `10.9.7`, Go `go1.27.1 linux/amd64`, CGO_ENABLED=`1`을 확인했다. SSH helper의 known_hosts와 RejectPolicy, 비밀값 가림을 유지했다. 접속값이나 인증 URL은 기록하지 않았다.

전체 명령과 실제 종료코드는 [commands.json](commands.json)에 있다. 모든 실행은 `set -e`와 exact HEAD·clean 검사를 먼저 수행했다. SSH channel의 종료코드를 직접 기록했다. 파이프로 명령 종료코드를 가리지 않았다.

| 검사 | 결과·증거 |
|---|---|
| 지원 환경·exact SHA·clean 검사 | 0 · [environment-verified.log](environment-verified.log) |
| `go test -race ./internal/relay -count=1 -v -run '^(TestConnectionResponseKeepsRetryAndErrorPriority\|TestDeviceFailureBoundaries\|TestHumanAgentListsKeepOwnerAndFriendBoundaries\|TestHumanFriendTransitionsAndDiscovery\|TestHumanGateRevocationAndReaccept\|TestFrozenParsingAndSigning\|TestDeviceConsentTemplate\|TestMemberUIFixtures)$'` | 0 · 실제 최상위 8 tests 통과 · [go-boundaries.log](go-boundaries.log) |
| `npm run build --prefix adapters` | 0 · 현재 소스의 TS·MCP bundle을 다시 만들어 focused 검사의 stale artifact 가능성을 없앴다 · [ts-build.log](ts-build.log) |
| `node adapters/dist/mcp.test.js && node adapters/dist/owner-text.test.js` | 0 · [ts-boundaries.log](ts-boundaries.log) |
| 독립 `node --input-type=module` text matrix | 0 · 15개 경계 · [text-matrix.log](text-matrix.log) |
| 실행 뒤 exact SHA·clean 확인 | 0 · [post-clean.log](post-clean.log) |

첫 환경 프로브는 JS 확인식의 셸 따옴표 오류로 EXIT=1이었다. [environment.log](environment.log)와 원 명령을 보존했다. 이는 제품 실패가 아니다. 확인식을 고친 뒤 환경과 나머지 검사가 통과했다. 증거 대조의 첫 바이트 동일성 assert도 checkout 개행 차이로 실패했으며 LF 정규화 대조로 구분했다.

## 보존한 동작과 경계

| 대상 | 독립 확인 |
|---|---|
| `connectionResponse` | deviceAPI/connectAPI 두 호출자를 읽었다. time.Time retry가 transaction 오류보다 먼저다. 명령 error map이 transaction 오류보다 먼저다. 일반 transaction 오류와 성공 응답도 유지한다. 429 JSON과 Retry-After의 반올림이 일치한다. |
| `humanAgents` | owner 필터·철회 제외·ID 순 정렬을 유지한다. 키 없는 agent는 목록에서 연결 필요로 보인다. friend/new의 현재 친구·block·상대 회원 자격 검사는 helper 호출 전에 유지된다. 차단/해제·키/회원 철회 gate 회귀도 통과했다. |
| `validText` | 일반·trial 송신과 두 수신 경로가 공유 함수를 호출한다. 4096 bytes ASCII/Korean/emoji는 허용한다. 초과·빈 값·ASCII/Unicode 공백·고립/삽입 surrogate는 거부한다. 정상 surrogate pair는 허용한다. 앞뒤 공백도 원래 byte 한도에 포함한다. 명시 confirm=false는 relay 요청 없이 거부한다. |
| MCP 포장 | SDK stdio 실제 bundle의 도구 9개 discovery, text content JSON의 state/transport/actualConnection/webhook/evidenceFetch, isError와 기본 held를 검사했다. invalid config와 redirect를 거부하고 held에서 network/key read를 하지 않는다. |
| 서명·Inbox·ACK | frozen signing 회귀가 통과했다. owner wire는 agent 서명을 만들지 않는다. recipient/schema/surrogate 실패는 persist/ACK를 하지 않는다. 정상 owner text는 local store 이후 persist/ACK를 한다. 공유 helper 밖의 서명·저장 순서는 diff가 없다. |
| API·CSRF·TTL·Events·UI | 변경 제품 파일 11개의 diff를 확인했다. handler 외부 API·CSRF·TTL·MCP Events 계약, template/CSS/문구, schema, 의존성·배포 설정의 변경은 없다. 기존 exact SHA 전체 테스트를 재사용하고 frozen/template 회귀를 새로 실행했다. |

## 패킷 보완과 한계

packet 183개 경로와 189개 path/category 쌍의 정본 존재·검색 근거를 [packet-audit.json](packet-audit.json)에 남겼다. 큰 문서·코드의 unknown과 documents-find 잔여 후보는 일반 `rg`로 보완했다. 추가 보완 180개 경로는 모두 존재한다. 과거 기록은 시점과 관련 검색 구간을 확인했으며 전문 재독·현재 SHA의 새 실행을 주장하지 않는다. packet direct_edit는 DEV 변경 seed이므로 tester는 제품을 편집하지 않았다. 모든 쌍의 처리 결과는 [packet outcomes](../../jev/SAR-REFACTOR-001-TESTER-packet-outcomes.json)에 있다.

Windows Go race는 cgo 부재 때문에 실행하지 않았다. DEV Node24/npm11 결과는 보조 근거이며 최종 지원환경 증거로 쓰지 않았다. Docker/실제 SQL migration·운영 배포·실계정 메시지·Google/외부 Bot 왕복·부하·백업 복원·장시간 wake는 not_run이다. 관련 계약이 바뀌지 않은 이번 lite 리팩토링 범위를 넘는다. UI 변화가 없어 새 직접 시각 판정은 해당 없음이다. 전체 제품 HP-16 수락과 과거 held 항목을 이번 QA로 해제하지 않는다.

새 코드·테스트 파일은 만들지 않았다. 실행용 코드는 stdin으로만 전달했고 재현 명령을 JSON에 보존했다. 따라서 이번 문서·로그 변경에 새 제품 등록 게이트는 적용하지 않는다. 지시대로 exact 후보의 기존 전체 게이트를 재사용했다. `git diff --check`와 완료 패킷 검사를 수행한다. 최종 기록 커밋은 worker_done의 SHA로 식별한다. bundle·clone·기존 증거와 이번 로그를 모두 보존한다. 후속 fixed-SHA 리뷰·main 통합은 coordinator가 담당한다.
