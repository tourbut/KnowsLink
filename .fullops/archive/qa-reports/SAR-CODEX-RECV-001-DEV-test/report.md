---
title: SAR-CODEX-RECV-001-DEV — Codex 세션 수신 V02 실측과 검사
status: draft
updated: 2026-10-22
owner: dev
tasks: [SAR-CODEX-RECV-001-DEV]
summary: "격리 relay에서 Claude→Codex·Codex→Claude·Codex↔Codex를 App Server 세션과 독립 TUI로 나눠 실측했다. 지원 수준은 '세션 결속 로컬 사본 + 다음 조회'다. YOLO 자동 수락 결함을 찾아 확인 코드로 고쳤다. tester D-CR-01(첫 호출 경쟁의 고아 결속)을 생성 promise 공유와 on 중복 차단으로 고쳤다."
---

# SAR-CODEX-RECV-001-DEV — DEV 검증 보고

DEV 자체 검증이다. 독립 tester QA와 코드 리뷰는 coor가 후속 배정한다. 기술 계획은 [실행 기록](../../../exec-plans/phases/SAR-CODEX-RECV-001-DEV.md)이다. 실제 회원·운영 relay·실제 상대에게 보내지 않았다. 사용자 `~/.codex`·`~/.claude` 설정은 바꾸지 않았다.

## 환경

| 항목 | 값 |
|---|---|
| 제품 | `fullops/dev`, 기준 `9c397d98`. 완료 SHA는 지시서 완료 보고에 적는다. V02는 커밋 전 작업 트리 빌드(bundle sha256 앞 16자 `2e2c145143a20d59`)로 했고, 마지막 문구 변경(`거절하거나 취소하세요`) 뒤 빌드로 App Server 연결·해제 1회를 다시 확인했다([ev-B-final-wording.jsonl](ev-B-final-wording.jsonl)) |
| 호스트 | Windows 11, Codex CLI 0.163.0(모델 표시 `GPT-6-Luna low`), Claude Code 2.1.308, Node 22.22.2, Docker Desktop |
| relay | 작업 트리 소스의 Compose 프로젝트 `knowslink-codexrecv-dev`, `http://127.0.0.1:18311`, `KNOWSLINK_SYNTHETIC_SIGNUP=1`, 메일 sink 사이드카 |
| 회원 | `*@codexrecv.invalid` A(Claude)·B(Codex)·C(Codex), A–B·B–C 친구, A–C 아님 ([fixture.py](fixture.py), [fixture-agents.json](fixture-agents.json)) |
| 설치 | `package_plugin.py --host claude|codex --install-dir … --connect --verify` A·B·C exit 0 ([pkg-a.log](pkg-a.log), [pkg-b.log](pkg-b.log), [pkg-c.log](pkg-c.log)). RELAY_URL만 격리 relay로 바꿨다 |
| Codex App Server 세션 | [drive.mjs](drive.mjs): `codex app-server`(stdio)에 붙는 클라이언트. remote TUI가 쓰는 같은 프로토콜이다. 설치본 `.mcp.json`을 `-c mcp_servers.knowslink.*`로 준다. 합성 소유자의 대화상자 응답은 명령 파일에서만 온다. `approvalPolicy: on-request`, `sandbox: read-only` |
| Codex 독립 TUI | Orca 터미널의 `codex`([tui-b.ps1](tui-b.ps1)). 사용자 기본 설정(YOLO)과 `-c approval_policy='on-request' -c sandbox_mode='workspace-write'` 두 가지. hook 신뢰 창은 "Continue without trusting"을 골랐다 |
| Claude 세션 | Orca 터미널의 `claude --plugin-dir <inst-a>/knowslink`, `/knowslink-session on` |

## 판정 요약 (V02)

수준 정의는 SP-12다. '세션 결속 로컬 사본'은 선택한 세션의 MCP 프로세스가 열린 동안 받아 그 세션의 agent 폴더 inbox에 저장했다는 뜻이다. '다음 조회'는 모델이 `knowslink_text_receive`를 불러야 본문을 본다는 뜻이다.

| 방향 | 세션 종류 | 기존 세션 수신(도착 표시) | 세션 결속 로컬 사본 | 다음 조회 수신 | 새 세션 실행 | 근거 |
|---|---|---|---|---|---|---|
| Claude→Codex | App Server 세션(remote 준비) | 미지원: 도착 표시 관측 없음 | PASS | PASS | 사용 안 함 | V02-1, V02-7 |
| Claude→Codex | 독립 TUI(on-request) | 미지원: TUI 화면에 도착 표시 없음 | PASS | PASS | 사용 안 함 | V02-12 |
| Codex→Claude | App Server 세션 → Claude | Claude 쪽 기존 세션 수신(상태줄 `대기 1`) | PASS | PASS | 사용 안 함 | V02-2, V02-8 |
| Codex→Claude | 독립 TUI → Claude | 같음 | PASS | PASS | 사용 안 함 | V02-11 |
| Codex↔Codex | App Server B ↔ App Server C(별도 프로세스) | 미지원 | PASS(작업 중 도착 포함) | PASS | 사용 안 함 | V02-3–6, V02-9–10 |
| Codex↔Codex | 독립 TUI B ← App Server C | 미지원 | PASS(전환 전 사본 보존) | PASS | 사용 안 함 | V02-10 |
| 모든 Codex 방향 | 독립 TUI, YOLO(전체 자동) | — | 연결 불가(아래 결함 수정 뒤 `declined`) | — | — | ev-T-yolo-final |

정리: Codex 쪽 지원 표기는 '세션 결속 로컬 사본 + 다음 조회 수신'이다. '작업 세션 수신 지원'(기존 세션에 도착이 보임)으로 표기하지 않는다. YOLO 모드에서는 연결·발송이 되지 않는다(안전 거부). Codex 하위 에이전트(SP-04)와 4096 bytes급 본문의 TUI 대화상자 전체 표시는 미검증이다.

## 메시지별 기록

App Server 세션의 대화상자·도구 결과·turn 기록은 [ev-B-prefix.jsonl](ev-B-prefix.jsonl)·[ev-C-prefix.jsonl](ev-C-prefix.jsonl)(수정 전 빌드, V02-1–6)과 [ev-B-final.jsonl](ev-B-final.jsonl)·[ev-C-final.jsonl](ev-C-final.jsonl)(최종 빌드, V02-7–10)이다. Claude A의 발송 승인 창은 [ev-A-send-dialog.txt](ev-A-send-dialog.txt), TUI의 KnowsLink 발송 승인 창은 [ev-T-tui-send-dialog2.txt](ev-T-tui-send-dialog2.txt)다. 모든 발신 receipt는 `transport: delivered`와 별도 `session_received_at`을 가졌다([receipts.txt](receipts.txt)). 본문은 합성 문장이다.

| ID(앞 13자) | 방향 | 수신 세션 | 결과 |
|---|---|---|---|
| `01a16363-d435` V02-1 | A→B | App Server B(`codex-27bbb9`) | B inbox 사본 → B `receive`: `untrusted:true`, `sender` 친구·agent 이름, `sessionReceipt: recorded`. 본문의 "파일을 지워라"는 실행되지 않았다(명령 실행 승인 요청 0건, read-only sandbox) |
| `01a16364-d00d` V02-2 | B→A 답장 | Claude A | B 발송 대화상자 accept → `queued`. A 상태줄 `대기 1`([ev-A-toast.txt](ev-A-toast.txt)), `receive`로 표시 ([ev-A-receive.txt](ev-A-receive.txt)) |
| `01a16366-94f4` V02-3 | B→C | App Server C | C가 대기 중 도착, 다음 turn에서 `receive` |
| `01a16367-54e9` V02-4 | B→C | App Server C | C 생성 turn이 먼저 끝나 대기 중 도착 |
| `01a16367-d37b` V02-5 | B→C | App Server C | **작업 중 도착**: C turn 11:38:20–11:39:25(`completed`, 64.9초) 진행 중 11:38:32에 C inbox 사본 저장. turn은 중단되지 않았고 모델은 그 turn에서 본문을 보지 않았다 |
| `01a16369-3e68` V02-6 | C→B 답장 | App Server B(재시작 뒤 `codex-a4e08d`) | C 발송 승인 창에 답장(reply_to)·본문 전체 표시. B 세션을 다시 연 뒤 `대기 1건`으로 시작해 받았다(이미 기한 지남, `expired` 표시) |
| `01a16372-3e0c` V02-7 | A→B | App Server B(최종 빌드, `codex-a4e08d`) | 다음 조회 수신 |
| `01a16372-f7cc` V02-8 | B→A 답장 | Claude A | 코드 입력 승인 → A `receive` |
| `01a16373-0b1d` V02-9 | B→C | App Server C(최종 빌드) | 다음 조회 수신 |
| `01a16373-8fb8` V02-10 | C→B 답장 | 독립 TUI B | TUI가 `take` 뒤 `대기 1건`으로 시작, `receive`로 표시 ([ev-T-tui-receive.txt](ev-T-tui-receive.txt)) |
| `01a16379-099a` V02-11 | TUI B→A | Claude A | 첫 시도 틀린 코드 `0000` → `denied`, 모델 재시도 없음. 소유자 재요청에서 코드 입력 → 발송 ([ev-T-tui-send-wrongcode.txt](ev-T-tui-send-wrongcode.txt), [ev-T-tui-send-ok.txt](ev-T-tui-send-ok.txt)) |
| `01a16379-f91d` V02-12 | A→TUI B | 독립 TUI B | 다음 조회 수신 ([ev-T-tui-receive2.txt](ev-T-tui-receive2.txt)) |

Claude A 화면: 첫 발송 결과 [ev-A-send-result.txt](ev-A-send-result.txt), V02-7 승인 창·결과 [ev-A-send2-dialog.txt](ev-A-send2-dialog.txt)·[ev-A-send2-result.txt](ev-A-send2-result.txt), V02-8 수신 [ev-A-receive2.txt](ev-A-receive2.txt), V02-11 수신과 V02-12 승인 창 [ev-A-receive3-dialog.txt](ev-A-receive3-dialog.txt), V02-12 결과 [ev-A-send3-result.txt](ev-A-send3-result.txt). 독립 TUI 화면: 수정 전 YOLO take [ev-T-2.txt](ev-T-2.txt), Codex 도구 실행 승인 창 [ev-T-tui-dialog1.txt](ev-T-tui-dialog1.txt)·[ev-T-tui-receive-approval.txt](ev-T-tui-receive-approval.txt)·[ev-T-tui-send-dialog.txt](ev-T-tui-send-dialog.txt), take 결과 [ev-T-tui-bound.txt](ev-T-tui-bound.txt).

## SP 규칙별 관측

| 규칙 | 관측 |
|---|---|
| SP-01 명시 선택 | 설치·연결 완료만으로 수신하지 않는다(단위 테스트: 결속 전 pull 0). `on`은 소유자 대화상자를 거친다. 거절하면 `declined`다([ev-B-bind-declined.jsonl](ev-B-bind-declined.jsonl)). 대화상자에 `Codex 세션(구분값 codex-xxxxxx)`과 agent ID가 보인다 ([ev-T-tui-dialog2.txt](ev-T-tui-dialog2.txt)) |
| SP-02 한 세션 | TUI `on`(take 없음) → `busy`와 점유 세션 구분값·시각, 대화상자 없음 ([ev-T-1.txt](ev-T-1.txt)). `take` → 대화상자에 "다른 세션 codex-a4e08d 의 수신은 끊깁니다" → 코드 입력 → 전환. 옛 App Server B의 `receive`는 "다른 세션이 이 agent를 가져가 이 세션의 연결이 해제됐습니다" |
| SP-05/06 발송 승인·자동 답장 0 | 모든 발송에 KnowsLink 대화상자가 따로 떴다. 거절·입력 없는 accept·틀린 코드는 발송 0. 수신한 본문 때문에 발송한 경우는 없다 |
| SP-07 신뢰 경계 | 본문은 도구 결과에만 있다(`untrusted:true`). Codex 자체 도구 승인(`Allow the knowslink MCP server to run tool …`)과 KnowsLink 대화상자는 별개로 떴다 |
| SP-08 작업 보존 | V02-5 작업 중 도착, turn `completed` |
| SP-09 연결 수명 | `off` → `unbound`, lock 삭제, 이후 `receive` 거부 ([ev-T-tui-off.txt](ev-T-tui-off.txt)). TUI `/quit`은 Windows에서 lock을 남겼고 다음 `on`이 죽은 pid lock을 회수했다 |
| SP-11 상태 | `knowslink_session status`는 구분값·agent·대기 수만 보인다 |
| SP-12 근거별 지원 | 위 표 |

## 결함과 수정

| # | 등급 | 내용 | 조치 |
|---|---|---|---|
| F1 | high(수정됨) | 첫 구현은 elicitation `accept`만으로 승인했다. 사용자 기본 설정이 YOLO인 Codex TUI에서 `take`가 대화상자 없이 `bound`됐다 ([ev-T-2-full.txt](ev-T-2-full.txt)). 같은 방식이면 발송도 소유자 확인 없이 나갈 수 있다(SP-05/06) | 대화상자마다 새 4자리 확인 코드를 보이고 입력 일치만 수락한다. 같은 YOLO TUI에서 `declined`, lock 유지 ([ev-T-yolo-final.txt](ev-T-yolo-final.txt)). 단위 테스트에 입력 없는 accept·틀린 코드 거부 추가. 변이(무조건 accept) 검사로 테스트가 실패함을 확인 |

미해결 critical/high는 없다.

## D-CR-01 수정 (재개 attempt)

tester [D-CR-01](../SAR-CODEX-RECV-001-TESTER-test/report.md)(medium): 새 MCP 프로세스의 첫 호출 두 개가 겹치면 `codexSession()`이 `await memberTransport()` 사이에 `CodexSession`을 둘 만들었다. 결속된 고아 객체는 `off`가 닿지 않아 lock·자동 수신이 남았다(SP-02·SP-09).

| 변경 | 내용 |
|---|---|
| `adapters/src/mcp.ts` `codexSession()` | 생성 promise `sessionInit`를 프로세스에서 하나만 만들고 모든 호출이 공유한다. `memberTransport` 실패 때만 비워 다음 호출이 다시 시도한다(기존 실패 동작 유지) |
| `adapters/src/codex-session.ts` `on()` | 같은 객체에서 승인 창이 열린 동안 들어온 둘째 `on`은 창을 열지 않고 `busy`("연결 승인 창이 이미 열려 있습니다")를 돌려준다. 공유된 객체에서 `on`+`on`이 창 두 개와 `AutoReceiver` 두 개를 만드는 2차 경로를 막는다 |
| `adapters/src/codex-session.test.ts` 6 | 새 프로세스에서 `status`+`on`+`on` 동시 첫 호출: 창 1개, `bound`/`busy` 하나씩, `status`가 결속 객체의 구분값을 보임, `off` 뒤 lock 없음·receive denied. 수정 전 코드로 빌드하면 "status sees the bound session"에서 실패함을 확인했다 |

회귀 실측: tester 도구를 수정 번들(`plugin.js` sha256 앞 16자 `2985e0ca951741d8`)로 tester 격리 relay `knowslink-codexrecv-qa`(127.0.0.1:18321, 합성 회원 A·B)에서 실행했다. 도구 디렉터리는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/codex-recv-001-devfix/`이다. tester 설치본의 설정 사본에서 경로만 바꾸고, 회원 폴더 사본은 원본과 같은 소유자 전용 ACL로 맞췄다(상속 ACL 사본은 `memberTransport`가 거부해 첫 실행이 전부 `failed`였다. 제품 결함 아님).

| 검사 | 결과 |
|---|---|
| `approval-matrix.mjs` 전체 | 50/50 PASS ([dcr01-matrix.out](dcr01-matrix.out), [jsonl](dcr01-matrix.jsonl)) |
| RC-4 `on`+`on` | PASS, dialogs=1, 둘째는 `busy`, `off` 뒤 lock 없음 |
| RC-5 `status`+`on` 4회 | PASS, 고아 0/4. 매 회 `off` 뒤 lock 없음·receive denied |
| `orphan.mjs` 8회 | 고아 0 ("no orphan in 8 rounds", [dcr01-orphan.out](dcr01-orphan.out), [jsonl](dcr01-orphan.jsonl)) |
| `npm run check --prefix adapters`(Node 22.22.2) | EXIT=0 ([dcr01-adapters-check.log](dcr01-adapters-check.log)) |
| `npm test --prefix adapters`(Node 22.22.2) | EXIT=0, `.../parallel first call` PASS 포함 ([dcr01-adapters-test.log](dcr01-adapters-test.log)) |

실측 뒤 relay·postgres 컨테이너는 정지만 했다. 실제 Codex TUI·App Server 재실측은 하지 않았다. 변경은 MCP 서버 안의 객체 생성 순서뿐이고 tester 매트릭스가 같은 stdio MCP 경계를 쓴다.

## 관찰

- O-1: Codex는 도구마다 자체 실행 승인 창을 먼저 띄운다(on-request). 소유자는 한 번의 발송에 창 두 개를 본다. YOLO에서는 앞 창이 생략되고 KnowsLink 창은 자동 수락돼 코드 불일치로 거부된다. 그래서 YOLO 세션에서는 KnowsLink를 쓸 수 없다. D11 사용자 설명에 적어야 한다.
- O-2: TUI 대화상자에 별도 '거절' 버튼은 없고 Esc 취소다. 안내 문구를 "거절하거나 취소하세요"로 바꿨다.
- O-3: TUI `take` 대화상자에서 1분 32초 동안 응답하지 않았어도 도구 호출이 시간 초과되지 않았다. 서버 쪽 elicitation 제한은 5분이다.
- O-4: Orca `terminal send`는 Codex 입력란에 text와 Enter를 한 번에 보내면 `agent_prompt_blocked`로 막혔다. 코드와 Enter를 따로 보냈다(시험 도구 한계).
- O-5: App Server 세션 C의 `Start-Sleep` 셸 실행은 read-only sandbox에서 실패했다. 작업 중 도착은 긴 생성 turn으로 관측했다.
- O-6: Claude 설치본의 `plugin.js`를 교체하자 Claude 플러그인이 다시 로드돼 `/knowslink-session on`이 풀렸다(문서화된 한계). 다시 `on`한 뒤 같은 key로 보냈다.

## 검사

| 검사 | 결과 |
|---|---|
| `npm run check --prefix adapters`(Windows) | EXIT=0 ([adapters-check-win.log](adapters-check-win.log)) |
| `npm test --prefix adapters`(Windows) | EXIT=0, 새 `codex-session.test.js` PASS 포함 ([adapters-test-win.log](adapters-test-win.log)) |
| product-lint / product-test(Linux 격리 clone, 완료 SHA) | 완료 보고에 적는다 |
| FullOps lint.py `--from 9c397d98` | 완료 보고에 적는다 |

## 정리 상태

Codex App Server 드라이버·TUI·Claude 세션은 `off` 뒤 종료했다. Compose `knowslink-codexrecv-dev`와 sink는 정지만 하고 삭제하지 않는다. 도구·원본 로그는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/codex-recv-001-dev/`에 있다. 레포의 `ev-*.jsonl`은 원본 로그에서 대화상자·도구·turn 항목만 추린 것이다.
