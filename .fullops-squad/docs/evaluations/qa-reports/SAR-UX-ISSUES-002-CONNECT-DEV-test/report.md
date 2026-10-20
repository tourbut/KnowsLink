---
title: SAR-UX-ISSUES-002-CONNECT-DEV DEV 보고
status: review
updated: 2026-10-20
owner: dev
tasks: [SAR-UX-ISSUES-002-CONNECT-DEV]
summary: "회원 연결 경로 선택·재개(U8), 나머지 U7, O-NM-03, low 3건과 2차 수정(재요청 행 하나 재사용·저장 미확인 연결 화면·low 6건)의 구현·검사·렌더·fixture 기동 근거와 미실행 항목"
---

# SAR-UX-ISSUES-002-CONNECT-DEV DEV 보고

DEV 검사 기록이다. 독립 리뷰·tester QA·designer 시각 수락이 아니다. 실제 호스트 설치·지원 완료 표시는 범위 밖이며 다섯 호스트를 지원 완료로 표시하지 않았다.

## 기술 계획 (실행함)

1. 경로 선택: `/home/agents?use=&start=` GET 폼. `use`는 web·grok·dots·claude·other, `start`는 first·continue만 받는다. 표시 전용이며 저장하지 않는다. 선택지 순서·문구는 `connect_guide.go:useChoices` 한 곳에 둔다.
2. 재개 상태: `agentResume`가 agent마다 저장된 상태만으로 단계와 다음 행동 하나를 계산한다. 근거는 최신 연결 시도(waiting·prepared·approved, 만료·취소 tombstone 24h), 활성 키, remote install의 미교환 code(1분)·token, device 승인 뒤 미저장이다. 불확실하면 ‘확인 필요’다. 온라인·앱 설치·도구 검색은 관측하지 않는다.
3. 세션 연결 관측: relay는 `/knowslink-session on`을 관측하지 않는다(SSE 구독은 세션 화면 성공의 근거가 아님, 설계 U8). 따라서 Claude Code 안내는 ‘앱에서 확인 필요’와 `/knowslink-session status`로 남긴다.
4. O-NM-03: 상태 JSON `Friend.CodeTail`(승인된 필드 추가). 요청 시점 trim 입력의 끝 4자, 4자 미만·`[A-Za-z0-9_-]` 밖 문자는 `?`, 이전 기록은 빈 값. 있는/없는/차단/자기 코드 모두 같은 함수로 저장한다. 요청자의 보낸 대기·거절·취소·만료 행에만 보인다. 전체 코드는 저장·URL·로그에 넣지 않는다.
5. 중복 클릭: 생성 폼의 기존 `Attempt` 값을 기존 `Idempotency` 맵에 `#agent-create/owner/attempt → agent`로 기록한다. 새 필드·migration은 없다. agent 삭제 sweep이 해당 항목을 지운다.

## 변경 요약

- U8: 내 에이전트에 ‘어디에서 사용할까요?’(다섯 선택지, 보조 문구 포함)·‘처음 연결/기존 연결 계속’·고른 곳의 준비물·권한·다음 행동, ‘내 에이전트와 연결 상태’(agent별 단계·다음 행동·보낸 친구 요청 대기 수·요청 상태 조회 안내). 5번은 설치 명령·버튼·배지·도구 이름 없이 웹 사용으로 안내한다. OpenClaw·Meta Muse 이름과 ‘온라인’ 배지는 없다.
- U7: 연결 화면 consumed는 ‘키 등록됨’으로 바꾸고 단계 `aria-current=step`은 서버 상태의 한 단계만 둔다. 연결 화면·취소는 `/home/agents#agent-{id}`, confirm은 해당 연결 화면, 생성은 새 agent 위치로 복귀한다. 회원 연결 동작 거부는 내 에이전트(legacy 관계는 기존 연결 목록)로 복귀한다. 호출자 return URL은 받지 않는다. remote 앱 복귀 화면은 완료가 아님과 만료 시 앱에서 다시 시작을 안내한다. legacy 앵커 목록은 이미 ‘이 화면에서 찾기’였고 그대로 두었다.
- low 3건: 익명 조회 nav의 ‘내 에이전트’ 링크 제거, `gate_summary.go` comma-ok, 승인 목록 ‘요청 검토’ 접근성 이름에 받은 내 에이전트 추가.
- 문서: README·adapters 연결 절, D05·D10·D11·D13. D02/D04·D12는 바꾸지 않았다(D12는 새 운영 절차 없음).

## 검사

| 검사 | SHA | 결과 | 근거 |
|---|---|---|---|
| 원격 `make lint` | 37febf9 | EXIT=0 | [lint](37febf953bcb1d1c5ce3d7e58b4302058ab88ef7-lint.log) |
| 원격 `make test` | 37febf9 | EXIT=0 | [test](37febf953bcb1d1c5ce3d7e58b4302058ab88ef7-test.log) |
| 원격 integration 전체(`-race`, `-skip TestRemoteConsentBrowserReturn`) | 37febf9 | EXIT=0 | [integration](37febf953bcb1d1c5ce3d7e58b4302058ab88ef7-integration.log) |
| fixture 컴파일·실제 기동·probe | 37febf9 | EXIT=0 | fixture 절 |
| 리뷰 수정 뒤 integration 전체 | 4d0fe06 | EXIT=1 | [로그](4d0fe061f7bba4ae68f16b87eebc774f5acd77c4-integration.log). 생성 재전송 키의 NUL 접두사를 Postgres JSONB가 거부해 `POST /home/agents`가 503, 이어 2건 연쇄 실패. `37febf9`에서 `#` 접두사로 고쳤다 |
| 첫 실행 `TestConnectionResumeAndCodeTail` | bec551f | EXIT=1 ([로그](bec551f557723f489638621a72770fcdb5a34e7d-integration-TestConnectionResumeAndCodeTail.log)) | 429가 공용 admission에서 먼저 나 기대한 복귀 링크가 없었다(테스트 기대 문제). 403 2건은 첫 실패로 agent ID가 비어 생긴 연쇄 |
| 수정 뒤 `TestConnectionResumeAndCodeTail` | 137a7f5 | EXIT=0 | 하위 6건 PASS ([로그](137a7f58293c953a177bff52abee2e6a29bf5dc3-integration-TestConnectionResumeAndCodeTail.log)) |

U8-A 자동 검사 `TestConnectionResumeAndCodeTail`(Postgres): 선택지 순서·보류 호스트 미표시·5번 설치 경로 없음·읽기/새로고침으로 agent 0; 429 뒤 agent 0과 복구 뒤 같은 attempt로 정확히 1개, 중복 제출은 같은 위치; local waiting→prepared→approved 구분과 재개 링크, 응답 유실 뒤 같은 token 재완료 401·키 1개 유지, 취소·만료(rotate 포함) 뒤 기존 키 불변, 만료된 시도와 미사용 구분, 서비스 재시작 뒤 같은 상태; device 승인 뒤 미저장 ‘확인 필요’, remote 거절은 agent 0·재사용 불가, 승인 뒤 code 대기·유실(확인 필요)·교환 구분; O-NM-03 있는/없는/차단/자기 코드 행의 동일 모양, 형식 다름, 이전 요청, 상대 코드 교체 뒤 불변, 상대 화면 비노출, 전체 코드 비노출, 취소 행 유지; 친구·기존 키 보존. 단위 `TestCodeTailKeepsOnlyPlainTails`, `TestConnectGuidePagesRender`.

## 렌더

`TestConnectGuidePagesRender`가 `KNOWSLINK_UI_FIXTURE_DIR`로 만든 합성 HTML을 headless Chrome(비모바일 viewport) 1440·390·320에서 [render/probe.mjs](render/probe.mjs)로 렌더했다. 대상 9개(선택 없음·다섯 선택지 안내·보낸 요청 끝 4자·키 등록됨 연결 화면·remote 앱 복귀), 27장 모두 페이지 가로 넘침 0([probe.json](render/probe.json)). 320px에서 종료 행의 ‘입력한 코드 …abcd’가 줄 중간에서 끊겨 `nowrap`으로 묶었다. DEV 확인이며 designer 시각 수락이 아니다.

## fixture

`TestConnectBrowserFixture`(태그 `integration,connectfixture`, [runner](remote.py) `fixture`)를 원격 격리 clone `/tmp/knowslink-connect002-37febf953bcb-d5033ad1`와 자기 Postgres 컨테이너 `kl-connect002-fixture-37febf953bcb-cd981d46`에서 실제 기동했다. `fixture.log`에 `connect fixture ready; synthetic accounts only`, 프로세스 생존, `fixture.json` 권한 600. loopback probe: `GET /` 200, `POST /fixture/as/a` 204, 그 세션으로 `GET /home/agents`·`/home/agents?use=other&start=first`·`/home/friends` 200, 범위 밖 `GET /v1/pull` 403. 근거: [start](37febf953bcb1d1c5ce3d7e58b4302058ab88ef7-fixture-start.log), [probe](37febf953bcb1d1c5ce3d7e58b4302058ab88ef7-fixture-probe.log), [compile](37febf953bcb1d1c5ce3d7e58b4302058ab88ef7-fixture-compile.log).

| 항목 | 내용 |
|---|---|
| 계정 | 합성 a·b·c(`connect-*@example.test`), `POST /fixture/as/{a,b,c}`로 로그인(동일 출처·POST만) |
| a | 연결 시도 waiting·prepared·approved agent 각 1, 키 등록 agent 1(긴 영문 이름), 남은 1칸으로 생성·중복 클릭 시험. 보낸 요청: b의 실제 코드·없는 코드·형식 다름(`ab`)·이전 요청·취소 1, c에게서 받은 요청 1 |
| b | 만료·취소·device 승인 뒤 미저장·remote 앱 복귀 대기·remote 교환 유실·remote 교환 완료 agent. a의 받은 요청 1 |
| c | agent 0개 |
| 허용 | GET `/`·`/home`·`/home/agents`(query 포함)·`/home/friends`·`/home/account`·`/home/legacy`·`/home/new`·`/home/receipts`·`/home/connections/*`·`/connect/*`; POST `/home/agents`·`/home/connect`·`/home/confirm`·`/home/cancel`·`/home/friends/request`·`/home/friends/decision`. 그 밖은 403 |
| 수명 | 12시간(열린 연결 시도·remote code도 같은 수명으로 늘림). manifest에 origin·agent/연결 ID·계정·만료 시각 |

키·연결 변경(생성·연결)은 로그인 뒤 5분 재확인 창 안에서만 된다. 시험 전에 `/fixture/as/{name}`로 다시 로그인한다. 이전 SHA `4d0fe06` fixture 프로세스도 남아 있다(coor가 새 SHA로 다시 기동한다).

## 하위 리뷰

읽기 전용 하위 리뷰 1건(Purpose: review, 고정 `bec551f`)을 배정했다. coor의 독립 고정 SHA 리뷰를 대신하지 않는다. 결과와 처리(`d4a1026`):

| 등급 | 지적 | 처리 |
|---|---|---|
| high | `connections_integration_test.go:317`의 거부 복귀 기대가 옛 ‘홈으로 돌아가기’ | 새 복귀 문구로 고침 |
| high | 받은 사람이 대기 요청을 차단하면 요청자 화면에 ‘해제’로 보여 코드 존재·행동이 드러남(기존 동작, O-NM-03 동일 모양 위반) | 차단된 대기 요청을 `blocked-pending`으로 저장하고 요청자에게 기한까지 대기·뒤 만료로 보인다. 요청자 취소는 허용, 같은 코드 반복은 대기 중 반복처럼 무동작. 받은 사람은 ‘차단’, 해제 뒤에는 ‘해제’로 보이며 수락할 수 없다. 통합 검사에 같은 모양·행 수 단언 추가 |
| medium | 키가 있는 agent의 추가·교체 시도 중 ‘키 등록 전’ 표시 | ‘기존 키 등록됨 · 새 키 …’로 표시 |
| medium | 생성 재전송 키와 메시지 멱등 키의 이론적 충돌 | `#`로 시작하는 `agentCreateKey`로 분리. agent ID는 `#`로 시작할 수 없어 메시지 키(`from/key`)와 겹치지 않는다. 처음 고친 NUL 접두사는 Postgres JSONB가 거부해 생성이 503이 됐고(`4d0fe06` integration 실패), `37febf9`에서 바꿨다 |
| low | 매 sweep의 Idempotency 전체 순회 | 삭제된 agent가 있을 때만 순회 |
| low | 철회된 agent로 재전송 복귀, 검사 공백 | 철회된 agent는 재전송 대상에서 제외. ‘온라인’ 단언은 배지 마크업 기준으로 유지(본문 설명 문장에 단어가 있음). Inbox 보존은 이번 변경이 메시지·Inbox 경로를 바꾸지 않아 기존 통합 검사(`TestHumanPostgres` 등) 통과로 대신한다 |

남은 한계: 차단된 요청이 기한을 넘긴 뒤 같은 코드로 다시 보내면 새 행이 생기고, 없는 코드는 같은 행이 다시 대기가 된다. 행 수 차이가 남으며 coor/designer 판단 대상이다.

## 2차 수정 (리뷰 F1–F4, designer V-CN-01·02, tester D-CT-01·O-CT-01·02·04)

두 번째 attempt다. 기준은 `be5c0f6`(main `a7f6b10` 통합)이다. 제품 커밋은 `d61cc6f`(수정 묶음), `7e4ca39`(렌더 변형), `d6369a4`(하위 리뷰 후속)다. 이 절이 위 ‘남은 한계’ 문단을 대체한다.

| ID | 등급 | 처리 | 자동 검사 |
|---|---|---|---|
| F1 = V-CN-01 = O-CT-01 | medium | `requestFriend`가 코드마다 요청자 행 하나를 재사용한다. 끝난 요청을 같은 코드로 다시 보내면 그 행이 세대 +1, 새 24h 기한, 같은 ID로 대기가 된다. 상대가 차단했거나 상대 수신 대기가 10건 이상이면 같은 행을 `blocked-pending`으로 되살려 상대에게 보이지 않게 한다. 대리 행(`undisclosed:`)은 없는 코드·자기 코드·요청자 본인이 차단한 상대에만 쓴다. 발신 상한 10건과 pending 전체 상한은 기한 전 `blocked-pending`도 센다. 수신 한도 보류 행은 `Friend.Held`에 받은 사람이 직전에 보던 종료 상태(처음이면 `-`)를 남겨, 요청자가 취소·만료해도 받은 사람 화면이 바뀌지 않는다 | `TestFriendRerequestReusesOneRow`(있는 코드 만료·취소, 없는 코드, 자기 코드, 차단 뒤 기한, 차단 뒤 취소, 차단 해제 뒤, 수신 한도 8경우: 요청자 행 수·상태·라벨 모양·같은 ID·세대 +1·새 기한이 모두 같음, 차단한 사람·받은 사람의 렌더 화면이 재요청·취소 전후 바이트 단위로 같음, 대기 중 반복 무동작), `TestFriendBlockedRequestCountsTowardSentLimit`, `TestFriendLegacyPlaceholderAndStaleUnblock`. 이전 코드로 되돌리면 이 검사들이 실패함을 확인했다 |
| 기존 저장 데이터 | — | 지우거나 합치지 않는다(삭제는 승인 범위 밖이고 요청자는 이미 그 행을 봤다). 쌍 행이 없고 이전 대리 행만 끝난 채 있으면 재요청 때 그 행을 쌍 키로 옮겨 ID·세대를 유지한다(이전 수신 한도 경로). 대리 행이 대기 중이면 재요청은 무동작이다. 이전 차단 경로가 만든 쌍 행+대리 행 2행은 그대로 둔다. 상태 JSON 추가 필드는 `Held`(생략 가능) 하나이며 migration은 없다. 되돌리면 이전 코드는 `Held`를 몰라 받은 사람에게 보류 행을 종료 기록으로 보인다 | 같은 검사 |
| D-CT-01 | medium | 승인 뒤 만료된 device 연결(에이전트 활성 키 0)의 `/home/connections/{id}`가 카드와 같은 `deviceUnsure` 문구(‘확인 필요 · 승인 뒤 클라이언트 저장을 확인하지 못함’, 저장 결과 먼저 확인, 자동 철회 없음)와 연결·키 관리 링크를 보인다. 단계 목록·‘새로 시작하세요’는 없다. 승인 전에 만료된 device 연결은 에이전트가 없어 기존 만료 문구다 | 통합 `device_and_remote_uncertain_results_stay_check_needed`, 단위 `TestConnectGuidePagesRender`, fixture `device_lost` |
| V-CN-02 | low | 없는/남의 연결 ID 거부 화면 복귀를 ‘내 에이전트로 돌아가기’(`/home/agents`)로 바꿨다 | 같은 통합 하위 검사 |
| F2 | low | gate 요약이 nil이면 `Active`를 false로 두어 승인·거절 폼을 그리지 않는다 | 기존 gate 통합 검사. Parse가 모양을 보장해 현재 도달 경로가 없어 새 검사는 없다 |
| F3 | low | 키가 있는 agent의 새 시도 상태를 ‘기존 키 등록됨 · 새 키 등록 전 · {단계}’로 바꿔 ‘키 키’ 중복을 없앴다 | 문구 변경, 새 검사 없음 |
| O-CT-02 | low | legacy 목차에 보이는 제목 ‘이 화면에서 찾기’(`aria-labelledby`)를 넣었다. `aria-current`는 없다 | `TestConnectGuidePagesRender` |
| O-CT-04 | low | 5분 재확인 거부 문구를 ‘아래 ‘로그인 다시 확인’으로 신원을 확인한 뒤 이 작업을 시작한 화면에서 다시 시도하세요.’로 바꿨다. 거부 화면에 재확인 버튼과 ‘내 에이전트로 돌아가기’가 함께 있다. 재확인 완료 뒤 이동 위치(`/home`)는 바꾸지 않았다 | 문구 변경 |
| F4 | low | 이번 원시 로그는 LF로 정규화해 복사했고 manifest는 커밋 blob과 같은 바이트로 계산했다 | `raw-log-sha256.txt`의 2차 항목 |
| O-CT-03 | low | 지시서대로 바꾸지 않았다 | — |

### 하위 리뷰 (2차)

읽기 전용 하위 리뷰 1건(Purpose: review, 고정 `d61cc6f`). coor의 독립 고정 SHA 리뷰를 대신하지 않는다. critical·high 0. 처리(`d6369a4`):

| 등급 | 지적 | 처리 |
|---|---|---|
| medium | 재요청이 공유 행의 세대를 올려 차단한 사람 화면의 차단 해제 폼 숨은 값이 바뀌고, 열어 둔 화면의 해제가 실패한다 | 차단 행은 세대 0으로 보이고 차단 해제는 본인 차단 플래그로만 판단한다(세대 대조 없음) |
| medium | 수신 한도 보류 행을 요청자가 취소하면 받은 사람에게 ‘취소’ 기록으로 드러난다. 보류 재사용 때 받은 사람의 이전 거절 기록이 사라진다 | `Friend.Held`로 받은 사람 화면을 이전 그대로 유지 |
| low | 이전 수신 한도 대리 행만 있는 데이터에서 재요청 시 새 쌍 행이 생겨 2행 | 끝난 대리 행을 쌍 키로 옮김 |
| low | 수신 한도로 보류된 T→R 행을 R이 T의 코드로 요청하면 R의 새 행 세대가 1이 아님 | 남은 한계(드문 경로, 받은 사람 R이 보류를 알 수 있는 정보는 세대 숫자 하나) |
| low | 되돌리기 시 보류 행 노출 | 남은 한계로 기록(D11 data-model) |
| low | 차단 해제 뒤 차단한 사람 화면이 재요청 여부에 따라 다름 | 차단에서 나온 `blocked-pending`은 다시 ‘해제’로 보인다(1차 동작 복원). 보류만 숨긴다 |

### 렌더 (2차)

바뀐 화면 3개(device 저장 미확인 연결 화면, legacy 목차, 연결 거부)를 1440·390·320에서 [render-fix/probe.mjs 결과](render-fix/probe.json)로 렌더했다. 9장 모두 페이지 가로 넘침 0. 저장 미확인 화면(390)과 legacy 목차(1440)를 직접 열어 문구·순서를 확인했다. DEV 확인이며 designer 시각 수락이 아니다. 렌더 도구는 기존 [render/probe.mjs](render/probe.mjs)이며 Chrome 임시 프로필을 OS 임시 폴더에 만든다.

### fixture 추가 경로 (2차)

합성 계정 d·e를 추가했다(`POST /fixture/as/{d,e}`). manifest의 `rerequest_codes`에 코드가 있다.

| 키 | 상태 | tester 재현 |
|---|---|---|
| `blocked` | A가 25시간 전에 D에게 보낸 요청을 D가 대기 중 차단, A 화면은 만료 | A로 이 코드를 다시 보내면 A의 행은 하나이고 대기로 바뀐다. D 화면은 그대로 ‘차단’ |
| `receive_limit` | A가 25시간 전에 E에게 보낸 요청이 만료, E는 합성 발신자 10명의 대기 요청을 받음 | A로 다시 보내면 A의 행은 하나이고 대기. E 화면은 A의 만료 기록 그대로 |
| `unknown_cancelled` | 없는 코드 `fixture-cancel-ca11`의 취소 행 | 비교 기준: 다시 보내면 같은 행이 대기 |

### 검사 (2차)

원격 Linux 격리 clone(`/tmp/knowslink-connect002-<sha12>-*`)과 task 소유 Postgres 컨테이너에서 [remote.py](remote.py)로 실행했다. 원시 로그는 LF로 정규화해 복사했고 `raw-log-sha256.txt` 끝 14줄이 그 hash다(F4).

| 검사 | SHA | 결과 | 근거 |
|---|---|---|---|
| 원격 `make lint` | d61cc6f · 7e4ca39 · d6369a4 | EXIT=0 | [d6369a4 lint](d6369a43c5c41b7bf97823187a25e12f492cde60-lint.log) |
| 원격 `make test`(`-race` 포함) | d61cc6f · 7e4ca39 · d6369a4 | EXIT=0 | [d6369a4 test](d6369a43c5c41b7bf97823187a25e12f492cde60-test.log) |
| 원격 integration 전체(`-race`, `-skip TestRemoteConsentBrowserReturn`) | d61cc6f · d6369a4 | EXIT=0 (`ok … 147.174s`) | [d6369a4 integration](d6369a43c5c41b7bf97823187a25e12f492cde60-integration.log), [d61cc6f](d61cc6f449cb47def3a65fa5e60e07bf49b5a148-integration.log) |
| fixture 컴파일 | d61cc6f · d6369a4 | EXIT=0 | [compile](d6369a43c5c41b7bf97823187a25e12f492cde60-fixture-compile.log) |
| fixture 실제 기동·probe | d6369a4 | EXIT=0. manifest 권한 600, `GET /` 200, `POST /fixture/as/a` 204, 회원 화면 3개 200, 범위 밖 `/v1/pull` 403 | [start](d6369a43c5c41b7bf97823187a25e12f492cde60-fixture-start.log), [probe](d6369a43c5c41b7bf97823187a25e12f492cde60-fixture-probe.log) |
| fixture 재요청 왕복 | d6369a4 | EXIT=0. `blocked`·`receive_limit`·`unknown_cancelled` 모두 A의 라벨 언급 1→2(행 1개가 대기로 바뀌어 취소 버튼 이름에도 나옴), D·E 화면 불변 | [rerequest](d6369a43c5c41b7bf97823187a25e12f492cde60-fixture-rerequest.log), 도구 [rerequest_probe.py](rerequest_probe.py). 첫 실행은 `__Host-` Secure 쿠키를 cookie jar가 http에서 버려 비로그인 요청만 나갔고(상태 변화 없음) 429로 끝났다. 쿠키를 직접 실어 70초 뒤 재실행했다 |
| 로컬 단위 `go test ./internal/relay`, `go vet`(태그 없음·integration·connectfixture) | d6369a4 | PASS | 로컬 Windows go1.27.1 |
| FullOps lint.py `--from 022c748` | 최종 HEAD | 아래 ‘FullOps lint’ | coor 지시(D)대로 깨끗한 detached 임시 worktree에서 실행 |

fixture 상태: 위 재요청 왕복이 A의 세 행을 이미 대기로 바꿨다. tester 재현에는 fixture 재기동이 필요하다. `7e4ca39` fixture 프로세스(같은 원격 호스트, [start](7e4ca39ba7ad16a98b9b2fc485d41206b9dc4b71-fixture-start.log))와 이전 attempt 프로세스도 남아 있다. 정리·재기동은 coor가 정한다. 포워딩 18932는 건드리지 않았다.

### 남은 한계 (2차)

- 이전 차단 경로가 이미 만든 쌍 행+대리 행 2행 데이터는 합치지 않았다.
- 수신 한도로 보류된 T→R 행을 R이 T의 코드로 요청하면 R의 새 행 세대가 1이 아니다(리뷰 low 4).
- 되돌리기 시 이전 코드가 보류 행을 받은 사람에게 보인다.
- 수신 한도 보류 재사용 때, 받은 사람이 예전에 보낸 쪽이었던 행은 받은 사람 화면에서 받은 요청 기록 모양으로 바뀔 수 있다(드문 경로).
- F2는 현재 도달 경로가 없어 자동 검사를 더하지 않았다.
- 렌더 도구는 OS 임시 폴더에 Chrome 프로필을 만들고 지우지 않는다(기존 도구 동작).

## 실행하지 못한 것 (not_run)

- 실제 키보드·브라우저 200% 확대·스크린리더·보조 브라우저(U7-A), fixture 브라우저 왕복(tester 범위).
- 실제 호스트(Grok Bot·Dots·Claude Code·Codex·Hermes) 설치·인증·왕복. 범위 밖이다.
- `TestRemoteConsentBrowserReturn`(실제 Chromium 필요, 이전 단위와 같이 skip).
- 공용 admission 429 화면의 복귀 링크는 기존 `/home`이다. 이번 범위에서 바꾸지 않았다.
