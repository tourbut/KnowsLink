---
title: SAR-SESSION-PEER-001-DEV — Claude↔Claude 세션 왕복 DEV 검증
status: draft
updated: 2026-10-17
owner: dev
tasks: [SAR-SESSION-PEER-001-DEV]
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
| watcher CLI smoke ([watcher-smoke.log](watcher-smoke.log)) | `--session smokeA1` bound → 두 번째 `smokeA2` busy exit 2 → `smokeA3 --take` bound, 옛 watcher `knowslink_session_replaced` |
| `package_plugin.py --host claude --verify` | 추출 standalone MCP 검사 exit 0, 준비/private 경계 PASS, ZIP에 hooks 3파일 포함 |

등록 gate(`make lint`·`make test`)의 실행 결과는 최종 커밋 SHA에서 외부 로그로 보존하고 worker_done에 적는다. Windows에는 make가 없어 Linux 격리 clone(remote helper)으로 실행한다.

## 판정과 한계

- 지원 수준: '기존 세션 수신'(도착 toast·상태줄·세션 결속 로컬 사본, 작업 중·대기 양쪽). 본문 모델 표시는 '다음 조회'. Codex·Hermes 미실행.
- 모드 소스 hot reload가 열린 세션의 watcher를 끝내 두 번 재 `on`이 필요했다. 설치 플러그인 폴더가 고정된 제품 사용에서는 발생하지 않는다.
- 첫 send에서 모델이 잘못된 인자명을 써 빈 peer 대화상자가 한 번 떴다. 인자 누락 deny를 추가한 뒤 같은 상황은 대화상자 없이 거부됐다(03:23 관측).
- toast 글자는 Orca 터미널 screen read에 나타나지 않아 창 캡처로 확인했다. 캡처는 터미널 영역만 잘랐다.
- V03(동명·자식 세션·제삼자)·V04(철회·재수락)·V07(정책 비활성 등)은 이번 DEV 검증 범위 밖이며 tester가 독립 수행한다.
