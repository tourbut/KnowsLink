---
title: SAR-SESSION-PEER-002-DEV — Claude Code 세션 메시지 후속 8건 DEV 검증
status: draft
updated: 2026-10-17
owner: dev
tasks: [SAR-SESSION-PEER-002-DEV]
summary: "승인 창 2000자 한도·대화상자 중 해제·비가시 문자·재연결 대기 수·친구/agent 이름·relay 통신 불가·세션 수신 확인을 자동 검사와 격리 relay의 실제 Claude Code 두 세션으로 확인한 기록"
---

# SAR-SESSION-PEER-002-DEV — DEV 검증 기록

DEV 자체 검증이다. 독립 리뷰·tester 재확인·운영 수락이 아니다. 실제 회원·운영 relay·실제 상대에게 보내지 않았다.

## 환경

| 항목 | 값 |
|---|---|
| 호스트 | Windows 11, Claude Code 2.1.302, Node 22.22.2(KnowsLinkDevTools), Docker Desktop |
| relay | 작업 트리 소스의 Compose 프로젝트 `knowslink-sp002-dev`, `http://127.0.0.1:18142`, `KNOWSLINK_SYNTHETIC_SIGNUP=1`, 메일 sink(loopback). 6·8 관측 전 같은 프로젝트의 relay를 변경 소스로 다시 빌드했다 |
| 회원·agent | 이메일 합성 회원 `합성 소유자 A`/`합성 소유자 B`, 친구 수락. agent A `agent_e754f9668c8e7692ff9e7b`(이름 `Claude 세션 A`), B `agent_bdaaef675789d1982f90e0`(`Claude 세션 B`) |
| 설치 | `package_plugin.py --host claude --install-dir … --connect` 두 개, RELAY_URL만 격리 relay. 변경한 `hooks/register.ts`·`dist/plugin.js`를 복사해 사용 |
| 세션 | Orca 터미널 두 개(`KL2-A`, `KL2-B`)에서 `claude --plugin-dir <설치>/knowslink`. 6·8 관측 전 두 세션을 다시 시작했다 |

## 항목별 결과

| # | 항목 | 변경 | 자동 검사 | 실제 세션 관측 |
|---|---|---|---|---|
| 1 | O-RT-01 승인 창 잘림 | 질문 전체가 2000자를 넘으면 대화상자를 열지 않고 deny하며 나눠 보낼 길이를 안내한다. 기존 4096자 한도는 이 한도에 포함된다 | 질문이 정확히 2000자인 본문은 대화상자로 가고 끝까지 포함, 1자 초과·4096자는 deny(`register.test.ts`) | 한글 1171자 본문(약 20줄)이 끝 `끝표시KTAIL」 을 그대로 보낼까요?`까지 보였다: [dialog-korean-1171.txt](dialog-korean-1171.txt). 줄 수가 아니라 글자 수가 한도임을 확인했다. 1880자 본문은 대화상자 없이 `승인 창은 2000자까지만 보여 … 약 1879자 이하로 나눠` deny: [send-fit-and-over-deny.txt](send-fit-and-over-deny.txt) |
| 2 | 리뷰 #7 대화상자 중 해제 | 대화상자 전 `link`를 잡아 두고 응답 뒤 같은 연결이 아니면 deny | 대화상자 안에서 `off` 뒤 `보내기` → deny, 발송 0 | 대화상자를 연 채 CLI watcher `--take`(05:56:30Z bound) → 15초 뒤 `보내기` → `승인하는 동안 이 세션의 연결이 해제돼 발송을 막았습니다`, B inbox 증가 없음: [release-during-dialog-before.txt](release-during-dialog-before.txt), [release-during-dialog-after.txt](release-during-dialog-after.txt) |
| 3 | 리뷰 #8 비가시 문자 | `\p{Cf}`·`\p{Default_Ignorable_Code_Point}`(변이 선택자 FE00–FE0F·E0100–E01EF, ZWSP, bidi 제어 포함)를 `⟨U+XXXX⟩`로 바꾸고 개수를 질문에 적는다. peer·reply_to·이름에도 적용 | FE0F·200D·202E·E0101 4개 표시 | `보이는글자⟨U+FE0F⟩숨김⟨U+200B⟩끝⟨U+202E⟩역순`, `보이지 않는 문자 3개`: [dialog-invisible.txt](dialog-invisible.txt). 취소 |
| 4 | 리뷰 #9 최대 길이 표시 | 1번 한도 | 1번 | 질문 정확히 2000자(영문 본문 1879자)가 끝 `END-ASCII-TAIL」 을 그대로 보낼까요?`까지 보였다: [dialog-ascii-2000.txt](dialog-ascii-2000.txt). `보내기` → `01a1486c-8b5f-71ea-bb02-7e74c4bf443d` queued, B 상태줄 `대기 1` |
| 5 | O-02 재연결 대기 수 | watcher `knowslink_session_bound`에 inbox 미확인 수 `pending`, 모드가 연결 문구·상태줄에 반영 | bound `pending:2` → `대기 2건`·상태줄 `대기 2`, watcher CLI bound 줄의 `pending` | B `off` → `on` 직후 `대기 1건`·상태줄 `대기 1`(inbox `.json` 1개): [reconnect-pending.txt](reconnect-pending.txt) |
| 6 | O-03 친구·agent 이름 | relay `GET /v1/contacts/{agent}`(사용자 승인). watcher toast, `knowslink_text_receive`의 `sender`, 대화상자(`plugin.js contact`)에 `친구 <이름> · agent <이름> (<ID>)`. 조회 실패 시 ID만 | Go: 친구 전 403·친구 뒤 4필드만·자기 agent 403. 모드: 대화상자·toast 이름. node: `sender` | 대화상자 `→ 친구 합성 소유자 B · agent Claude 세션 B (agent_bdaa…)`: [dialog-names.txt](dialog-names.txt). B toast `발신 친구 합성 소유자 A · agent Claude 세션 A (agent_e754…)`: [toast-names.png](toast-names.png). receive `sender:{friend:"합성 소유자 A",name:"Claude 세션 A",client:"node-local"}`: [receive-names-seen.txt](receive-names-seen.txt) |
| 7 | O-05 relay 통신 불가 | fetch에 HTTP 응답이 없으면 `relay unreachable`(시간 초과도 같은 분류). MCP 실패 결과는 unreachable·5xx일 때 '키·친구 문제가 아님' 안내. watcher `knowslink_relay` ok/unreachable/error → 상태줄 `통신 불가 · 재시도 중`·`수신 오류 · 키·친구 확인`, toast, status | 500 → 안내 문구, 닫힌 포트 → `relay unreachable`, 모드 상태줄·toast·status | relay 컨테이너 정지(05:58:06Z) 뒤 35초 안에 두 세션 상태줄 `통신 불가 · 재시도 중`. receive 결과 `error:"relay unreachable"`, `next:"… not a key or friend problem …"`, status `relay 통신 불가(자동 재시도 중)`: [relay-outage.txt](relay-outage.txt). 05:59:34Z 재시작 → 06:00:51Z 상태줄 `연결` 복귀 |
| 8 | O-06 세션 수신 확인 | relay `POST /v1/text/seen`(수신자 전용, delivered·현재 세대, 첫 시각 유지)과 receipt `session_received_at`(사용자 승인). `knowslink_text_receive`가 표시 뒤 보고하고 `sessionReceipt: recorded|not_recorded`를 돌려준다 | Go: delivered 전 403, 발신자 403, 수신자 200, 재호출 시 시각 불변, 발신 receipt에 별도 필드. node: seen 1회·`recorded` | B receive 두 건 모두 `recorded`. A receipt `01a1487c-…`: `transport:"delivered"`, `completion:"received"`, `session_received_at:"2026-10-17T06:12:57.079152Z"`: [receipt-session-received.txt](receipt-session-received.txt) |

판단(8): 기존 receipt의 `completion:"received"`는 relay ACK와 같은 시점이라 '세션 수신 확인' 근거가 아니었다. 서버 계약 추가가 필요해 coordinator에 물었고 사용자 결정(2026-10-17, 6·8 모두 승인, DB 스키마 변경 없음)으로 구현했다. 값은 기존 relay 상태 JSON의 Message 필드(`SessionSeen`, omitzero)에 저장한다. 사람의 읽음은 뜻하지 않는다.

## 자동 검사

| 검사 | 결과 |
|---|---|
| `claude plugin test adapters` ([mod-test.log](mod-test.log), 원 출력 [mod-test.full.log](mod-test.full.log)) | `hooks/register.test.ts` 11 pass. `src/*.test.ts`는 node 검사 파일이라 kit 로드 실패로 표시되며 기존과 같이 판정 제외 |
| `npm run check`(adapters) ([adapters-check-win.log](adapters-check-win.log)) | exit 0 |
| `npm test`(adapters) ([adapters-test-win.log](adapters-test-win.log)) | exit 0, PASS 10줄 |
| relay 통합 검사 `go test -tags=integration -count=1 ./internal/relay` ([relay-integration.log](relay-integration.log)) | `ok … 104.721s`. `scripts/verify_mvp.py`와 같은 임시 Compose(Postgres)에서 실행, 끝에 그 임시 프로젝트를 내렸다. Windows라 `-race`·synthetic.js 단계는 빼고 실행 |
| 등록 gate product-lint/product-test | 최종 SHA에서 Linux 격리 clone(remote helper)으로 실행하고 결과를 worker_done에 적는다 |

## 한계와 남은 일

- 2000자는 Claude Code 2.1.302의 관측값이다. 호스트가 바뀌어 한도가 줄면 다시 관측해 `ASK_MAX`를 낮춘다. 대화상자의 줄 수 한도는 120·175열 터미널에서 관측되지 않았다.
- 이름 조회는 대화상자마다 `plugin.js contact`를 한 번 실행한다(최대 5초). 실패하면 ID만 보인다. 이름은 친구가 정한 글이라 잘라 보이고 비가시 문자를 표시한다.
- 운영 relay는 아직 새 엔드포인트가 없다. 배포 전까지 이름은 ID로, `sessionReceipt`는 `not_recorded`로 보인다. 운영 배포는 이번 범위가 아니다.
- toast 캡처는 Orca 창 전체를 찍은 원본 `toast-names-1.png`–`toast-names-6.png`에서 toast 영역만 잘라 [toast-names.png](toast-names.png)로 남겼다. 원본 6개는 저장소 규칙(worker 삭제·이동 금지)에 따라 지우지 않고 커밋하지 않은 채 이 폴더에 남겼다.
- 첫 시도에서 모델 안전장치가 반복 본문의 그대로 옮기기를 멈췄다. 대신 모델이 직접 쓴 한글 수필(1171자)과 파일의 자연 영문(1879자)으로 관측했다. 한글 최대(4096바이트, 약 1365자)는 직접 관측하지 못했고, 1171자가 약 20줄로 끝까지 보인 결과와 글자 수 한도 관측에 맡긴다.
- 독립 리뷰·tester 재확인·main 통합은 coor가 조정한다.
