---
title: SAR-CODEX-RECV-001-DEV — Codex 세션 수신 기술 계획
status: draft
updated: 2026-10-22
owner: dev
tasks: [SAR-CODEX-RECV-001-DEV]
summary: Codex 세션 수신 경로로 MCP 세션 결속과 소유자 elicitation 승인을 고른 근거와 한계
---

# SAR-CODEX-RECV-001-DEV — Codex 세션 수신 기술 계획

사용자 결정 D027에 따라 C-01 Codex 세션 수신을 구현한다. 제품 기준은 [SAR-SESSION-PEER-001](../../design-docs/SAR-SESSION-PEER-001.md)의 SP-01–12·V02다. 기준 커밋은 main `9c397d98`이다(지시서 커밋 `c78f0da7`). 실측 결과는 [DEV QA 보고](../../evaluations/qa-reports/SAR-CODEX-RECV-001-DEV-test/report.md)에 있다.

## 확인한 근거

- 설치 Codex CLI는 `codex-cli 0.163.0`(Windows 11, ChatGPT 로그인)이다. Context7은 월간 한도 초과로 조회하지 못했다. 대신 설치본의 `codex --help`·`codex queue --help`·`codex features list`·`codex app-server generate-json-schema`를 읽었다.
- `codex features list`: `tool_call_mcp_elicitation stable true`. App Server 프로토콜에 서버 요청 `mcpServer/elicitation/request`가 있다(`McpServerElicitationRequestParams`, 응답 `accept|decline|cancel`).
- App Server 클라이언트 요청에 `turn/start`·`turn/steer`·`thread/inject_items`가 있다. `thread/inject_items`는 "model-visible history"에 Responses API item을 넣는다.
- `codex queue --thread <id> --message <text>`는 기존 세션에 사용자 메시지를 넣는다.

## 경로 선택

| 후보 | 판정 | 이유 |
|---|---|---|
| `codex queue` | 사용 안 함 | 사용자(소유자) 입력 경로다. 친구 본문을 넣으면 SP-07 위반이다. 고정 doorbell만 넣어도 자동 turn을 시작해 SP-08의 자동 재촉에 가깝다. shared daemon에 붙은 세션만 대상이며 독립 TUI 적용은 확인하지 못했다. |
| App Server `turn/start`·`turn/steer`·`thread/inject_items` | 사용 안 함 | 모두 모델이 보는 대화 기록에 직접 넣는다. 친구 본문이 소유자 지시나 사실 기록으로 들어간다(SP-07). 현재 turn ID 경합도 있다. App Server(remote) 세션에만 적용된다. |
| MCP 세션 결속 + 다음 조회 | **선택** | Codex는 세션마다 MCP 서버 프로세스를 따로 띄운다. 그 프로세스를 세션으로 보고 기존 `session-lock`·`AutoReceiver`를 재사용한다. 본문은 `knowslink_text_receive` 도구 결과(`untrusted:true`)로만 모델에 들어간다. 독립 TUI와 App Server 세션 모두 같은 경로다. |

## 설계

- `scripts/package_plugin.py --host codex --connect`는 `.mcp.json` env에 `KNOWSLINK_SESSION_HOST=codex`를 더한다. 이 모드의 MCP 서버는 시작 때 자동 수신하지 않는다.
- 새 도구 `knowslink_session {action: on|off|status, take?}`(Codex 모드에서만 등록). `on`은 소유자 elicitation을 통과해야 `claimSession`으로 agent 폴더를 잡고 `AutoReceiver`를 시작한다(SP-01, SP-02). 다른 살아 있는 세션이 잡고 있으면 묻지 않고 `busy`다. `take:true`는 전환 사실을 대화상자에 적는다. 10초마다 lock을 확인하고 잃으면 수신을 멈춘다(SP-09). `off`와 MCP 연결 종료는 수신을 멈추고 자기 lock을 지운다.
- 구분값은 MCP 프로세스마다 새로 만든 `codex-<hex6>`다. 비밀이 아니며 새 세션에 이어지지 않는다.
- Codex 모드의 `knowslink_text_send/receive/receipt`는 결속된 세션에서만 동작한다. 그 밖은 `{state:"denied"}`다.
- 발송 승인(SP-05, SP-06): 모델의 `confirmed:true`는 승인이 아니다. MCP 서버가 직접 만든 대화상자에 친구 이름·agent 이름·ID, 새 메시지/답장(reply_to), 본문 전체 글자 수, 비가시 문자 `⟨U+XXXX⟩` 표시와 본문 전체를 보인다. 대화상자가 열린 동안 결속이 풀리면 거부한다.
- 확인 코드: 실측에서 Codex의 YOLO(전체 자동) 모드는 MCP elicitation을 소유자에게 보이지 않고 자동 수락했다. 그래서 단순 accept는 승인이 아니다. 대화상자마다 새 4자리 코드를 보이고 소유자가 그 코드를 입력해야 수락한다. 잘못된 코드·입력 없는 accept·거절·취소·시간 초과(5분)·elicitation 미지원 클라이언트는 모두 거부다.
- 수신 알림: 기존 `notifications/message`를 그대로 보낸다. Codex 0.163.0 TUI에서 도착 표시는 관측하지 못했다. 그래서 본문 확인 수준은 '다음 조회'다.

## 신뢰 경계(SP-07)

친구 본문은 `knowslink_text_receive`의 도구 결과에만 들어간다. 결과는 `untrusted:true`와 발신 친구·agent 이름을 붙인다. 사용자 입력(`codex queue`), 대화 기록 주입(`thread/inject_items`), 현재 turn steer, 시스템·개발자 지시에는 넣지 않는다. 승인 대화상자는 MCP 서버가 만든 문장이며 본문은 「」 안의 데이터로만 보인다. 받은 본문이 '보내라'·'승인하라'를 담아도 발송에는 소유자의 코드 입력이 필요하다. Codex 자체의 도구 실행 승인(`Allow the knowslink MCP server to run tool …`)은 별개 대화상자이며 KnowsLink 승인으로 세지 않는다.

## 한계

- 하위 에이전트(SP-04): Codex의 하위 에이전트가 부모 MCP 프로세스를 공유하는지 검증하지 않았다. 공유하더라도 발송은 소유자 코드 입력이 필요하다. 수신 조회는 결속된 프로세스에서 가능할 수 있다. '미검증'으로 둔다.
- Windows에서 TUI `/quit`은 lock 파일을 남겼다. 다음 `on`이 죽은 pid의 lock을 회수한다(Claude Code와 같은 한계).
- 대화상자 표시 한도: 1171자급 본문은 확인하지 않았다. 30자 내외 본문만 TUI에서 전체 표시를 확인했다.
- 플러그인 설치(`codex plugin add`) 경로 대신 같은 MCP 설정을 `-c mcp_servers.knowslink.*` 실행 인자로 주어 시험했다. 사용자 `~/.codex` 설정은 바꾸지 않았다.
