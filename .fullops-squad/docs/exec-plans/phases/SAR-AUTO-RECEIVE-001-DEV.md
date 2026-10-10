---
title: 자동 수신과 호스트 전달 실행 기록
status: draft
updated: 2026-10-10
owner: dev
tasks: [SAR-AUTO-RECEIVE-001-DEV]
summary: 자동 pull·로컬 보존·ACK·MCP 알림 구현과 Grok Bot 노우 wake 미지원 근거를 기록한다
---

# SAR-AUTO-RECEIVE-001-DEV — 자동 수신과 호스트 전달

## 기준

- 기준 main `dda6130`(handover front matter base `51eebe5`)을 포함한 `fullops/dev`에서 작업했다.
- 규칙: fullops-common-0.3.3, FULLOPS.md, project.md, rules/common의 coding-style·testing·security. 테스트 레벨 lite, 선택 하위 위임 off.
- SDK: 설치된 `@modelcontextprotocol/sdk` 1.32.0 원문(`dist/esm/server/index.js`의 `sendLoggingMessage`, `isMessageIgnored`). Context7은 coor 호출에서 월 한도 초과였다. 공식 문서로 보완했다.

## 실제 호스트 지원 확인

| 질문 | 근거 | 결론 |
|---|---|---|
| MCP 알림이 모델 턴을 시작하는가 | [MCP resources 명세](https://modelcontextprotocol.io/specification/2025-06-18/server/resources): 사용 방식은 host가 결정한다 | 보장 없음 |
| Grok Bot이 외부 이벤트로 대화를 깨우는가 | [Skills and routines](https://docs.x.ai/grok-bot/skills-routines-and-automations): routine은 schedule(최소 5분) 또는 Cursor 계정 통합(Slack·GitHub) event | MCP 알림·webhook·외부 inbound API는 없음 |
| 알림 설정 | [Settings and notifications](https://docs.x.ai/grok-bot/settings-and-notifications): Bot이 끝나거나 입력이 필요할 때 사용자에게 OS 알림 | inbound trigger 아님 |
| 공식 로드맵 | [포럼 요청](https://forum.cursor.com/t/let-a-grok-bot-computer-wake-its-own-agent-chat/168260): 팀 답변 “Wake mechanics is an area we're working on” | 미지원. 미공개 gateway 우회는 사용 금지 |
| Command MCP process 수명 | 공식 문서에 없음 | 미확정. 상시 watcher를 별도로 제공 |

5분 routine은 TTL 180s 안의 관련 답장을 보장하지 못한다. 사용량도 소비한다. 그래서 해결로 인정하지 않고 만들지 않았다.

## 기술 계획과 구현

1. `text.ts` `receive(store)`: 서명·만료·수신자 검증 뒤 persist/ACK 전에 로컬 저장을 호출한다. 저장 실패는 ACK를 막는다.
2. `inbox.ts` `Inbox`: `<연결 폴더>/inbox/`에 private 저장(fsync·rename), ID 중복 제거, rename claim 기반 1회 조회, 원문 없는 읽음 표시 24h.
3. `inbox.ts` `AutoReceiver`: process당 직렬 pull, idle 10s, 오류 backoff 20s–300s와 `retry_at`, drain, 상태 기록.
4. `mcp.ts`: logging capability, connect 뒤·로그인 완료 뒤 loop 시작, `notifications/message` metadata 알림, `knowslink_status.autoReceive`, 수동 receive의 inbox 우선, `KNOWSLINK_AUTO_RECEIVE=off`.
5. `plugin.js watch <폴더>`와 `text.js watch|receive`: 배포 bundle 하나로 상시 수신, 같은 inbox 사용.

relay·SQL·의존성·vendor core는 변경하지 않았다. 자동 답장·도구 실행·업무 효과는 없다.

## 검증

Node `C:/Users/shin/AppData/Local/KnowsLinkDevTools/node-v22.22.2-win-x64/node.exe` v22.22.2, npm 10.9.7을 사용했다.

- `npm run check`(prettier·eslint·tsc): 통과.
- `npm test`: 기존 4개와 새 `inbox.test.js` 모두 PASS, exit 0.
- `inbox.test.js`는 번들 `dist/plugin.js`를 SDK `StdioClientTransport`로 실제 실행한다. 로컬 relay 대역으로 다음을 확인했다.
  1. 수동 receive 호출 없이 pull→검증→inbox 기록→persist→ACK→알림이 일어난다. ACK 시점에 inbox 파일이 이미 있다. 알림에는 text가 없다.
  2. 재lease한 중복은 ACK만 하고 다시 알리지 않는다. 재시작 뒤 미확인 메시지를 수동 receive가 반환한다. 읽은 ID는 다시 표시하지 않는다.
  3. 만료 envelope와 철회된 발신 key(403)는 inbox와 ACK에 들어가지 않는다.
  4. relay 500은 `backoff`와 `lastError`로 보고된다. 성공으로 숨기지 않는다.
  5. `KNOWSLINK_AUTO_RECEIVE=off`는 pull하지 않는다. 수동 receive는 동작한다.
  6. `plugin.js watch`가 같은 inbox에 저장하고 metadata만 출력한다.
  7. 전체 시험에서 relay pull 동시 실행은 최대 1이다.
- 변이 확인: store를 ACK 뒤로 옮기면 `inbox.test.js`가 실패한다. 원복 뒤 다시 통과했다.
- 로컬 Docker는 시작하지 않았다.

## 한계와 후속

- 플러그인 자동 수신과 호스트 알림 전송은 로컬 실제 process에서 확인했다. Grok Bot 화면 표시·노우 턴 시작·자동 답장은 미지원이며 확인하지 않았다.
- 실제 Bot의 자동 수신은 이 후보를 Bot 컴퓨터에 설치한 뒤 운영 도메인 메시지 ID로 대조해야 한다. coor가 UI 조작 없이 조율한다. 미확인이면 제품 수락은 미완료다.
- 서로 다른 process(MCP 여러 개, watcher)는 각자 10s pull을 실행한다. 중복 표시는 ID로 막는다. 필요하면 폴더 lock으로 단일 poller를 강제한다.
- 노우 wake에 필요한 공식 기능: MCP 알림으로 대화 턴을 시작하는 host 계약, 또는 인증된 외부 trigger(webhook routine 등)의 공식 API.
