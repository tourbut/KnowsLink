---
title: SAR-AUTO-RECEIVE-001-TESTER 독립 QA
status: draft
updated: 2026-10-10
owner: tester
tasks: [SAR-AUTO-RECEIVE-001-TESTER]
summary: 자동수신 고정 SHA의 lite QA. 자동 저장과 호스트 알림과 노우 턴을 구분한다.
---

# SAR-AUTO-RECEIVE-001-TESTER 독립 QA

테스트 레벨은 lite다. 대상 SHA는 `af7627d8cca54df55e856f48d226f78dbf2b9d0f`이고 기준은 `dda6130d7023903175bbfa4036a7e052011e2c83`이다. 제품 코드는 바꾸지 않았다. 산출물 인덱스의 D10은 `docs/design-docs/module-design.md`이며 후보 안에 이미 자동수신 모듈이 있다. 이 파일은 그 설계서를 대체하지 않는 QA 기록이다. D12 `ops-guide.md`는 실제 운영을 적용하지 않아 수정하지 않았다.

## 구분

자동 저장은 서명과 관계 검증 뒤 로컬 inbox에 쓰고 ACK하기 전에 원문을 남기는 동작이다. 호스트 알림은 MCP `notifications/message`로 `event`, `id`, `from`, `pending`, `next`만 보내는 동작이다. 노우 확인과 답장은 호스트가 턴을 시작해 원문을 보고 사용자가 승인한 뒤에만 가능하다. `hostNotice: sent_unverified`와 wake의 `accepted_unverified`는 앞의 두 단계가 아니다.

설치본 `@modelcontextprotocol/sdk` 1.32.0은 `McpServer`에 logging capability가 있고 session level이 없으면 `sendLoggingMessage`를 알림으로 보낸다. 그 알림을 모델 턴으로 시작한다는 계약은 없다. Context7 조회는 월간 한도로 실패해 설치본 소스를 읽었다.

## 실행

| 구분 | 명령 | 환경 | 종료 | 메모 |
|---|---|---|---|---|
| 직접 | `npm run build`, `node dist/mcp.test.js` | Windows Node 22.22.2, detached af7627d | 0 | synthetic boundary PASS |
| 직접 | `python scripts/package_plugin.py --verify` | 같은 worktree | 0 | zip `666704d3d8da0dd20360ebf63d20581216819ede1b2dc3aeb203f90c0837aa7f` |
| 직접 | `node dist/inbox.test.js` | 같은 worktree | 0 | 53.05초. 자동 저장, 중복, 재시작, 만료, 철회, backoff, wake |
| 직접 | `npm ci`, `lint.py --from dda6130…` | Linux Node v22.22.2, 고유 `/tmp` clone | 0 | ERROR 0, WARNING 7. lint.json sha256 `0648aca688a400ae1f5c9c1fa7485adcde77ede21410edb09db2dca1101f1581` |
| 관찰 | `coor-af7627d-linux-package.log` | Linux, 새 npm ci 188패키지 | 0 | zip `f7d09684870514b0cfd368458aeaa99d968493abb92c52850ba6d46edb0c4eb6`, 추출 MCP 검사 exit 0 |

Windows worktree의 `node_modules`는 같은 lockfile blob `e69dc9c0d2f42a584717527a2e090aeaee789052`의 junction이다. Linux 패키지 로그는 그 junction 없이 `npm ci`를 했다. zip 해시가 다른 것은 아카이브 줄바꿈 차이고, 두 검사 모두 exit 0이다.

`inbox.test.js` 통과 문장은 수동 receive 없이 poll, 검증, ACK 전 로컬 복사, 알림, 중복, 재시작, read marker, 만료, 철회, backoff, opt-out, CLI watcher, 단일 in-flight pull, loopback wake의 거절·수락·중복·재시작·off·401·reset·invalid, 원문과 token 비노출을 포함한다. SEC-001 72행과 243행 인근은 합성 fixture다. 제품 비밀값이 아니다.

## 재현되지 않은 실패

coor msg_9ed204c7307a는 같은 SHA에서 Node 22 `npm run build` exit 0 뒤에 `package_plugin.py --verify`와 `adapters/dist/mcp.test.js`가 exit 1이라고 했다. 기본 PATH 탓으로 단정하지 말 것을 요청했다. 이 세션의 clean detached worktree와 명시 Node 22.22.2에서는 두 명령이 exit 0이다. Linux 패키지 로그도 추출 검사 exit 0이다. 테스트는 실패 시 assertion 본문을 출력하지 않고 exit 1만 남기므로, coor 쪽 예외 종류는 그 로그 없이는 특정하지 못한다. 이 SHA의 소스 회귀로 확정하지 않는다. 설치와 통합은 이 QA가 하지 않는다.

의존성 없이 돌린 Linux lint는 `prettier`와 `tsc`가 없어 ERROR 2였다. `npm ci` 뒤 같은 SHA와 같은 `--from`으로 다시 실행한 결과가 위의 통과 lint.json이다.

## 하지 않은 검증

실서버 배포, 실계정과 관계 변경, 그록봇 UI, 실제 gateway token wake, 추가 과금은 하지 않았다. 공개 relay preflight JSON은 수동 receive 없이 `lastSuccessAt`가 갱신되고 pending 0이다. 호스트 알림이나 노우 수신의 증거가 아니다. 실제 Bot 수신 수락은 coor 후속이다.

리뷰 발견은 `SAR-AUTO-RECEIVE-001-REVIEW-review/report.md`에 있다. medium 1건과 low 2건은 미해결이고 critical과 high는 없다.
