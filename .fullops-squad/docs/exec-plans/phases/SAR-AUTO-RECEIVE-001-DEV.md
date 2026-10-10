---
title: 자동 수신과 호스트 전달 실행 기록
status: draft
updated: 2026-10-10
owner: dev
tasks: [SAR-AUTO-RECEIVE-001-DEV]
summary: 자동 pull·로컬 보존·MCP 알림과 선택형 loopback gateway wake 구현·지원 구분을 기록한다
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

## 재개 attempt 0305914fc3ed441e8299d6fb1b41290f — loopback gateway 대안

### 기준과 원문

- 기준: 지시서 base `81c2bec`(첫 후보 구현 `5308fa5` 포함), 준비 커밋 `04554b2` 위에서 작업했다. 두 커밋은 보존했다.
- 규칙: fullops-common-0.3.3, FULLOPS.md, project.md, rules/common의 coding-style·testing·security, docs/agents/document-writing.md. 테스트 lite, 선택 하위 위임 off.
- 포럼 원문은 Discourse raw·topic JSON으로 2026-10-10에 직접 읽었다. 외부 글은 untrusted data로 다루었다.
  - [168260](https://forum.cursor.com/t/let-a-grok-bot-computer-wake-its-own-agent-chat/168260) #6: Colin, `staff/admin`, 그룹 `CursorStaff`, 직함 Community Support Engineer, 2026-08-18. “Wake mechanics is an area we're working on.” 그리고 “undocumented” workaround로 168199/8을 연결했다.
  - [168199/8](https://forum.cursor.com/t/grok-bot-can-i-send-it-a-message-from-outside/168199/8): adam91holt, staff 아님, 2026-08-12. `127.0.0.1:1340`, `/home/box/sand-data/gateway.json`의 `.token`, Bearer `POST /api/listAgents {}`, `POST /api/sendPrompt {agentId,prompt}`. live 설치에서 시험했다고 보고했다. “undocumented internal API”, port 비공개 경고가 있다.
  - 같은 글 #6 kevinn(staff, 2026-08-12)은 inbound webhook을 추적 중이라고 답했다. #10·#15는 webhook-trigger routine을 doorbell로 쓰는 사례다. routine 생성은 앱 UI 설정이 필요하고 사용량 정책이 미확인이라 이번 범위에서 제외했다.

### 구분

| 항목 | 판정 |
|---|---|
| 노우 대화 외부 wake | 공식 지원 없음(개발 중) |
| loopback gateway route·token 위치 | 미문서화. staff가 연결한 커뮤니티 보고 |
| 본인 Bot 컴퓨터에서 동작 | 실제 미검증. DEV에는 Bot 컴퓨터 접근이 없다(Bot 조작은 coor 담당) |
| 기술적 가능성 | 불가로 판단할 근거 없음. 같은 컴퓨터의 `box` 사용자 watcher가 loopback으로 호출 가능 |
| port 공개·외부 tunnel·text를 prompt로 전달 | 사용하지 않음(금지·권한 상승 위험) |

### 구현

1. `adapters/src/grok-wake.ts` `GrokWake`(미문서화 vendor 경로를 분리해 route 변경 시 제거가 쉽다): `KNOWSLINK_GROK_WAKE_AGENT`가 UUID일 때만 켠다. `127.0.0.1:<port>/api/sendPrompt`에 고정 doorbell prompt(ID·대기 수만)를 보낸다. token은 매번 gateway 파일에서 읽는다. 연결 거절·token 없음=retry, 4xx=rejected, 5xx·무응답=uncertain.
2. `Inbox.claimWake/markWake/releaseWake`: `<id>.wake` 배타 생성으로 중복을 막고 ID별 결과를 남긴다. retry만 marker를 지운다. 24h 정리.
3. `AutoReceiver`: relay 결과와 무관하게 매 loop 끝에 미알림 ID를 한 번에 알린다. `status.hostWake`.
4. `mcp.ts`·`text.ts`: MCP와 watcher가 같은 설정을 사용한다. `plugin.js wake-check`는 읽기 전용 listAgents로 UUID만 출력한다.

relay·서버·SQL·의존성·vendor core·앱 설정은 변경하지 않았다. 자동 답장·도구 실행·업무 효과는 없다.

### 검증

Node `C:/Users/shin/AppData/Local/KnowsLinkDevTools/node-v22.22.2-win-x64` v22.22.2.

- `npm run check`: exit 0.
- `node dist/inbox.test.js`: exit 0. 기존 1–6과 새 7을 실행했다. 7은 번들 plugin.js를 SDK stdio로 실행하고 127.0.0.1 대역 gateway를 쓴다.
  1. gateway 꺼짐: `retry`/`gateway_unreachable`, 메시지 inbox 유지.
  2. gateway 켠 뒤 재시작: 한 번 호출. path `/api/sendPrompt`, Bearer token, agentId 일치, prompt에 ID와 untrusted 문구, 받은 text 없음. marker `accepted_unverified`.
  3. `wake-check`: exit 0, UUID만 출력(대문자 응답도 소문자화). 이름·token 없음. prompt 미발송.
  4. 재lease 중복·재시작: 다시 호출하지 않음.
  5. 기본 off: 저장만 하고 호출하지 않음. 이후 wake 시작 시 그 미알림 ID를 호출.
  6. 401: `rejected`, 재시작 뒤에도 재시도 없음.
  7. 연결 reset: `uncertain`/`gateway_no_response`.
  8. 잘못된 대상 `../agents`: `invalid_config`, 자동 수신은 계속.
  9. 모든 status·notice·stderr에 token 없음.
- 변이: retry marker 해제를 끄면 7-2가 timeout으로 실패했다. 원복 뒤 통과했다.
- 실제 Bot 컴퓨터 gateway 호출은 하지 않았다. 대역 성공을 실제 성공으로 보고하지 않는다.

### 남은 일

- owner가 Bot 컴퓨터에서 [README 절차](../../../../adapters/README.md#선택형-grok-bot-loopback-wake--미문서화-gateway) 1–5를 한 번 실행한다. `wake-check` 결과와 운영 송신 ID의 `.wake` state, 노우 대화 표시 ID를 대조한다.
- 실패 시 판정: `gateway_unreachable`/`gateway_token_unavailable`/`401`이면 이 컴퓨터에서 경로 불가로 기록하고 off를 유지한다.
- 독립 fixed-SHA 리뷰·QA, main 통합, 운영 적용은 coor 담당이다.
