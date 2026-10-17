---
title: SAR-SESSION-PEER-001-DEV — Claude↔Claude 세션 왕복 DEV 검증
status: draft
updated: 2026-10-17
owner: dev
tasks: [SAR-SESSION-PEER-001-DEV, SAR-SESSION-PEER-001-DEV-FIX]
summary: "격리 relay와 합성 회원 두 명, Claude Code 세션 두 개로 V01 왕복·toast·발송 승인·세션 결속을 관측한 DEV 자동 검증 기록"
---

# SAR-SESSION-PEER-001-DEV — DEV 검증 기록

DEV 자체 검증이다. 독립 QA(V01–07)·리뷰·운영 수락이 아니다. 실제 회원·운영 relay·실제 상대에게 보내지 않았다.

## 환경

| 항목 | 값 |
|---|---|
| 호스트 | Windows 11, Claude Code 2.1.302, Node 22.22.2(KnowsLinkDevTools), Docker Desktop 29.2.1 |
| relay | 작업 트리 소스의 Compose 프로젝트 `knowslink-sp001-dev`, `http://127.0.0.1:18141`, `KNOWSLINK_SYNTHETIC_SIGNUP=1`, 메일 sink 사이드카(loopback) |
| 회원 | 이메일 코드 로그인 합성 회원 A/B(`*@sessionpeer.invalid`), 각자 `node-local` agent 1개를 `connect.js prepare/complete`로 연결, 친구 코드 요청·수락 |
| agent | A `agent_bf57b32c65d5da4ee7cf4f`, B `agent_5977cc13254c57956a9cd4` ([fixture-agents.json](fixture-agents.json)) |
| 설치 | `package_plugin.py --host claude --install-dir … --connect` 두 개, RELAY_URL만 격리 relay로 수정 |
| 세션 | Orca 터미널 두 개에서 `claude --plugin-dir <설치>/knowslink`, 빈 작업 폴더, 세션 구분값 A `71a93be4`, B `87399d7b` |

## 관측 기록 (V01)

| 시각(UTC) | 세션 | 관측 |
|---|---|---|
| 03:18:10 | B | `/knowslink-session on` → `연결됨: 세션 87399d7b · agent_5977…`, lock `{session:87399d7b,pid:51544}` |
| 03:19:xx | B | 100초 Bash 작업 시작(작업 중 상태) |
| 03:19:55 | A | `knowslink_text_send` 승인 대화상자 → `보내기` → `{"id":"01a147df-8d01-78bf-b1a0-9859afdfa99d","transport":"queued","exp":"03:22:55Z"}` |
| 03:19 작업 중 | B | 상태줄 `KnowsLink 87399d7b 대기 1`, inbox `01a147df….json` 생성, Bash 작업은 1분 48초에 `long task done`으로 정상 완료 |
| 03:22 | B | 사용자 요청으로 `knowslink_text_receive` → 같은 ID `01a147df…`, `untrusted:true`, `authorship:agent`. 본문 지시 없음, 답장 0 |
| 03:23 | B | `reply_to=01a147df…` 답장 승인 → relay `422 expired`(원 TTL 경과). 재시도 없이 보고 |
| 03:25:00 | B→A | 새 메시지 승인 → `01a147e4-335b-7c28-9245-c40d751436e7` queued. A(대기 상태) 상태줄 `대기 1`은 승인 7초 뒤 관측 |
| 03:25 | A | `knowslink_text_receive` → `01a147e4…` untrusted 표시, 답장 0 |
| 03:26 | A | `/knowslink-session off` 뒤 `knowslink_text_receive` 호출 → deny `이 세션은 연결되지 않았습니다` |
| 03:27–03:28 | B→A | 세 번째·네 번째 메시지 승인 → A toast 캡처: [toast-a-terminal.png](toast-a-terminal.png) (`KnowsLink 외부 메시지 도착 · 발신 agent agent_5977… · 대기 2건 · 확인: knowslink_text_receive`, 상태줄 `대기 2`) |

종료 시 두 세션 모두 `off` 후 `/exit`. 최종 inbox: A에 미확인 2건(`01a147e6…`, `01a147e7…`)과 읽음 표시 1건, B에 읽음 표시 1건. lock 파일은 죽은 pid로 남았다(Windows 종료 특성).

## 자동 검사

| 검사 | 결과 |
|---|---|
| `claude plugin validate adapters` | 통과. hooks: session.start, command.run{knowslink-session}, tool.call(.catch 있음) |
| `claude plugin test adapters` ([mod-test.log](mod-test.log)) | `hooks/register.test.ts` 3 pass: 미연결 send/receive deny·도구 미도달, status/off 문구, Bash 통과. `src/*.test.ts` 9개는 node 검사 파일이라 kit 로드 실패로 표시되며 판정에서 제외 |
| 패키지 추출 폴더 `claude plugin test` | 3 pass 0 fail |
| `npm run check --prefix adapters` ([adapters-check-win.log](adapters-check-win.log)) | exit 0 (prettier·eslint·tsc) |
| `npm test --prefix adapters` ([adapters-test-win.log](adapters-test-win.log)) | exit 0, PASS 10줄. `session-lock.test.js` 포함, `mcp.test.js` textBoundaries(V06 4096/4097) 포함 |
| watcher CLI smoke ([watcher-smoke.log](watcher-smoke.log)) | 로그 3줄: `smokeA1` bound → 옛 watcher `knowslink_session_replaced` → `smokeA3 --take` bound. 두 번째 watcher의 busy 출력은 이 로그에 없다. busy exit 2는 FIX 재검증의 [fix/watch-busy.out](fix/watch-busy.out)에 원 출력·종료코드로 남겼다 |
| `package_plugin.py --host claude --verify` | 추출 standalone MCP 검사 exit 0, 준비/private 경계 PASS, ZIP에 hooks 3파일 포함 |

등록 gate(`make lint`·`make test`)의 실행 결과는 최종 커밋 SHA에서 외부 로그로 보존하고 worker_done에 적는다. Windows에는 make가 없어 Linux 격리 clone(remote helper)으로 실행한다.

## 판정과 한계

- 지원 수준: '기존 세션 수신'(도착 toast·상태줄·세션 결속 로컬 사본, 작업 중·대기 양쪽). 본문 모델 표시는 '다음 조회'. Codex·Hermes 미실행.
- 모드 소스 hot reload가 열린 세션의 watcher를 끝내 두 번 재 `on`이 필요했다. 설치 플러그인 폴더가 고정된 제품 사용에서는 발생하지 않는다.
- 첫 send에서 모델이 잘못된 인자명을 써 빈 peer 대화상자가 한 번 떴다. 인자 누락 deny를 추가한 뒤 같은 상황은 대화상자 없이 거부됐다(03:23 관측).
- toast 글자는 Orca 터미널 screen read에 나타나지 않아 창 캡처로 확인했다. 캡처는 터미널 영역만 잘랐다.
- V03(동명·자식 세션·제삼자)·V04(철회·재수락)·V07(정책 비활성 등)은 이번 DEV 검증 범위 밖이며 tester가 독립 수행한다.

## FIX 재검증 — SAR-SESSION-PEER-001-DEV-FIX

독립 코드 리뷰(high 1, medium 2, low 3) 반영 후 DEV 자체 확인이다. 같은 격리 relay(`knowslink-sp001-dev`, 확인 뒤 다시 정지)와 합성 agent A 설치만 사용했다. 실제 회원·운영 relay에는 보내지 않았다.

| 항목 | 변경 | 근거 |
|---|---|---|
| 발송 승인 본문(high) | 대화상자가 본문 전체와 `본문 전체 N자`를 보인다. 4096자를 넘으면 묻지 않고 deny한다. | `register.test.ts` 201자 이상 본문 전체·4097자 deny. Claude Code 2.1.302 세션 A에서 255자 본문이 끝(`끝표시TAIL250`)까지 대화상자에 보임: [fix/dialog-full-body.txt](fix/dialog-full-body.txt). `취소` 선택 → deny, 모델 재시도 없음 |
| replaced 즉시 해제(medium) | `knowslink_session_replaced`에서 `link === current`이면 연결 해제·상태줄 초기화·자식 정지 | `register.test.ts` replaced 뒤 send/receive deny·status 미연결. 세션 A(`06199ecc`)에서 CLI `--take`(04:10:08 bound) 뒤 04:10:19 상태줄 해제와 toast(debug 로그 [fix/debug-ui-lines.txt](fix/debug-ui-lines.txt)), 이후 `knowslink_text_receive` deny와 status `연결되어 있지 않습니다`: [fix/replaced-deny.txt](fix/replaced-deny.txt), [fix/watch-take-bg.out](fix/watch-take-bg.out) |
| V05 자동 검사(medium) | 연결 상태 fixture(설정 파일·세션 ID·watcher 이벤트·소유자 대화상자를 test hook이 대신함) | `취소`·대화상자 실패(dismiss/`-p`)·`Other` 자유 입력·인자 누락이 MCP send에 도달하지 않음. 명령·파일 경로·도구 승인 요청을 담은 외부 본문 수신 뒤 Bash 호출 0·send 0·대화상자 0·추가 spawn 0, toast에 본문 없음 |
| reply_to 표시(low) | `새 메시지` 또는 `답장(reply_to <ID>)` | 같은 테스트 |
| smoke 증거(low) | 위 표 서술을 로그 범위로 고침, busy 출력 추가 | [fix/watch-busy.out](fix/watch-busy.out): `knowslink_session_busy` 한 줄, `exit 2` |
| 주석(low) | `package_plugin.py` 주석 `/knowslink-session on` | diff |

자동화하지 않은 부분: 모델이 외부 본문의 지시(명령·경로·승인 요청)를 따르지 않는지는 hook 테스트로 검사할 수 없다. 모델 판단이기 때문이다. 근거는 기존 V01 관측(본문 지시 없음·답장 0)과 `untrusted:true` 결과·SKILL 지시, 그리고 모드 경계다. 모델이 지시를 따르려 해도 send는 소유자 대화상자를 지나야 하고, 그 대화상자는 이제 본문 전체를 보인다. 악성 본문 수동 관측은 tester V05 범위로 남긴다.

관측된 이상: 첫 takeover 시도(세션 `c0d543f9`)에서 foreground `timeout 25`가 takeover watcher를 멈춘 시각(04:08:21 전후)에 그 Claude 세션도 종료됐다. 원인은 확인하지 못했다. watcher를 background로 띄운 두 번째 시도(`06199ecc`)에서는 replaced 해제 뒤에도 세션이 살아 있었고 `/exit`으로 정상 종료했다.

FIX 자동 검사: `claude plugin validate adapters` 통과, `claude plugin test adapters`의 `hooks/register.test.ts` 7 pass([fix/mod-test.log](fix/mod-test.log), 원 출력 [fix/mod-test.full.log](fix/mod-test.full.log)의 `src/*.test.ts` 로드 실패는 기존과 같이 판정 제외), `npm run check`·`npm test`(adapters) exit 0([fix/adapters-check-win.log](fix/adapters-check-win.log), [fix/adapters-test-win.log](fix/adapters-test-win.log)). 등록 gate는 최종 SHA에서 Linux 격리 clone으로 실행하고 결과를 worker_done에 적는다.
