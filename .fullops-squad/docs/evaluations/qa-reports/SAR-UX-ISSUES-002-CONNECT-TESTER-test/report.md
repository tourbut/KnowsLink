---
title: 연결 안내·재개·보낸 요청 코드 끝 4자 독립 QA
status: review
updated: 2026-10-20
owner: tester
tasks: [SAR-UX-ISSUES-002-CONNECT-TESTER]
summary: "재확인 abc2191(제품 d6369a4) 재요청 3경로 동일·D-CT-01·low PASS, 새 결함 없음. 이전 8d8bad4(제품 37febf9) 격리 fixture 독립 QA: U8 선택지·재개 행·중복 생성·5분 재확인, U7 연결 화면, O-NM-03 PASS. critical·high 없음, medium 1건(D-CT-01 승인 뒤 미저장 시도의 연결 화면 안내가 카드와 모순), low 관찰 4건"
---

# SAR-UX-ISSUES-002-CONNECT-TESTER 독립 QA

## 재확인 (attempt 2): DEV `abc219120509027b143b0e4ad121511a5a18395b`(제품 `d6369a4`)

**결론: 재요청 3경로 동일화, D-CT-01, low 재확인 모두 PASS. 새 결함 없음, critical·high·medium 0건.** 아래 ‘이전 결과’는 8d8bad4 기준이며 이 절이 D-CT-01·O-CT-01·02·04를 대체한다. 대상 SHA는 dev 워크트리 HEAD가 `abc2191`임을 읽기 전용으로 확인했다. 서버가 이 SHA인지는 coor 지시와 manifest를 따랐고 독립 확인하지 않았다. fixture는 `http://127.0.0.1:18932`의 초기 상태(계정 a~e)에서 시작했다. 제품 코드는 고치지 않았다. 테스트 레벨 lite.

| 영역 | 실행 | 결과 |
|---|---|---|
| HTTP 읽기 [checks.py](checks.py) → [http-result.json](http-result.json) | 95 | 95 PASS |
| 재요청 동일화 [rerequest.py](rerequest.py) → [rerequest-result.json](rerequest-result.json) | 36 | 36 PASS |
| mine 대기 상한 [cap.py](cap.py) → [cap-result.json](cap-result.json) | 4 | 4 PASS |
| 5분 재확인 [stale.py](stale.py) → [stale-result.json](stale-result.json) | 1 | PASS (422, agent 4개 그대로, O-CT-04 문구 5항 true) |
| 생성 [mutate.py](mutate.py) → [mutate-result.json](mutate-result.json) | 8 | 8 PASS |
| Chrome 155.0.8059.80 layout / keys / after | 76 / 33 / 8 | 모두 PASS |
| Edge 155.0.4283.62 layout / keys / after | 76 / 33 / 8 | 모두 PASS (keys는 재실행, 아래) |

실행 순서(비가역 경로 때문에 고정): checks.py → rerequest.py → cap.py → stale.py → mutate.py → browser.mjs after. 종료코드는 파이프 없이 확인했다(checks 0, rerequest 1→아래 설명, cap 0, stale 0, mutate 0, 브라우저 0). 브라우저는 `--headless=new --disable-extensions`, 비모바일 viewport 1440·390·320, 실제 키 입력은 CDP `Input.dispatchKeyEvent`다. 도구 변경: `checks.py`의 O-NM-03 기대값은 이전 fixture의 임의 코드 접미(`XXK2`·`Z-3N`)와 행 수 5를 하드코딩해 새 fixture 인스턴스에서 거짓 FAIL이 났다. 접미를 b·a의 친구 코드에서 읽고 행 수를 ‘원본 7 이상’으로 바꿨으며 판정 기준은 같다. D-CT-01·O-CT-02 검사는 낮추지 않고 더 엄격하게 했다(아래).

### 재요청 3경로 동일 (F1·V-CN-01·O-CT-01) — PASS

A가 `blocked`(D가 차단), `receive_limit`(E가 10건 수신), `unknown_cancelled`(없는 코드) 코드를 다시 보냈다. 시작 상태는 A의 보낸 행 7개(대기 4, 종료 3: 세 코드)다.

- 재요청 직후(P1): 세 경로 POST 응답이 모두 `303 /home/friends?n=requested`로 같다. A의 보낸 행은 7개로 늘지 않고 세 코드 각각 정확히 1행이 ‘요청 대기 / 입력한 코드 …끝4자 / 요청 기한 … KST / 요청 취소’로 바뀌었다. 정규화(끝 4자·시각만 치환)한 행 문장이 세 경로에서 같고, 취소 버튼 접근성 이름 형식(‘입력한 코드 …xxxx 친구 요청 취소, 기한 …’)과 기한 형식·24시간 길이도 같다.
- 행 순서: 세 코드의 위치(3·4·5)는 재요청 전후와 취소·재요청 뒤에도 같다.
- 대기 중 반복(P2): 같은 코드를 또 보내도 세 경로 응답이 같고 A 화면이 바이트 단위로 같다(무동작).
- 취소(P3): 세 행이 모두 ‘친구 요청 기록 / — 취소 / 입력한 코드 …끝4자 / 다시 연결하려면 새 친구 요청과 수락이 필요합니다.’로 같고 행 수는 7, 숨은 generation 값도 세 경로 모두 2다. 취소 응답도 `303 /home/friends`로 같다.
- 취소 뒤 재요청(P4): 다시 7행, P1과 같은 문장·순서다.
- **D·E 비추론**: D·E와 대조군 B·C의 `/home/friends` main HTML이 시작·P1·P2·P3·P4 모든 단계 뒤에 바이트 단위로 같았다. D는 계속 ‘— 차단 … 차단 해제’ 하나(결정 폼 `unblock`뿐, 회원 A 요청 카드 없음), E는 수신 대기 10건과 기존 ‘— 만료’ 기록 그대로다.
- 한계: 이 비교는 A·D·E의 서버 렌더 HTML 기준이다. 응답 시간·타이밍 채널은 보지 않았다.

첫 실행에서 `rerequest.py` 마지막 검사 1건(`d_still_shows_block_record_no_request_card`)이 FAIL이었다. 원인은 내 기대값 결함(‘수락’ 부분 문자열이 D 화면의 ‘새 친구 요청과 수락이 필요합니다’ 안내와 겹치고 `합성 회원 A`를 잘못 가정)이다. 제품 동작이 아니고 같은 단계의 바이트 동등 검사는 PASS였다. 기대를 결정 폼 값(`unblock`만)으로 고치고 그 검사만 재실행해 PASS다([첫 실행](rerequest-result-first-run-expectation-error.json)). 재요청은 비가역이라 전체를 다시 돌리지 않았다. 그래서 `rerequest-result.json`은 35개는 첫 실행, 1개는 재실행 값이다.

### mine 대기 상한이 보류 행을 센다 — PASS(관측)

[cap.py](cap.py): 재요청 뒤 A의 대기 보낸 행은 7개(그중 `blocked`·`receive_limit` 2개는 `blocked-pending` 보류 행)다. 없는 합성 코드 4개(`qa-cap-fake-0001~0004`)를 보내니 대기가 8·9·10까지 늘고 4번째는 조용히 버려졌다(POST 응답은 4번 모두 같음). 보류 행을 세지 않았다면 4번째도 나타났을 것이므로(7 대신 5가 센 값) 보류 행이 같은 방식으로 센다. 소스도 같다: `human.go:185` `open` 함수가 `blocked-pending`(기한 전)을 대기와 같이 세고 `mine >= 10`에서 `return`한다. 추가한 3행은 취소했다(A의 기록에 종료 행 3개 `…0001~0003`이 남음).

### D-CT-01 — PASS (이전 FAIL 재판정)

`GET /home/connections/{device_lost}`(b): 제목 ‘에이전트 연결’, 본문 ‘확인 필요 · 승인 뒤 클라이언트 저장을 확인하지 못함. 시작한 클라이언트의 저장 결과를 먼저 확인하세요. 쓰지 않는 에이전트라면 직접 철회할 수 있습니다. 자동 철회는 없습니다.’와 ‘연결·키 관리’·‘내 에이전트로 돌아가기’ 링크다. 카드와 같은 문구이고 ‘만료되었습니다’·‘새 키가 연결되지 않았습니다’·‘새로 시작’이 없으며 단계 `aria-current`도 없고 새 시작 버튼도 없다. 기존 FAIL 검사(`device_lost_connection_page_agrees_with_card_unknown_state`)의 기대는 그대로 두고 단정 부재 검사 3건을 더했다. 대조로 `expired`·`cancelled` 연결 화면은 기존 ‘만료/취소되었습니다. 새 키가 연결되지 않았습니다’를 유지한다(승인 전 종료라 맞다).

### low 회귀 — PASS

- V-CN-02: 없는 연결 ID와 남의 연결 ID 모두 403 ‘연결을 조회할 수 없습니다’ + ‘내 에이전트로 돌아가기’(`/home/agents`), 내부 정보 없음.
- O-CT-02: `/home/legacy` 목차는 `<nav class="toc" aria-labelledby="toc-h"><span id="toc-h" class="toc-title">이 화면에서 찾기</span>…`로 보이는 글자가 있고 `aria-current`는 없다.
- O-CT-04: 5분 지난 세션의 생성 POST는 422, 문구 ‘아래 ‘로그인 다시 확인’으로 신원을 확인한 뒤 이 작업을 시작한 화면에서 다시 시도하세요.’, ‘로그인 다시 확인’ 버튼과 ‘내 에이전트로 돌아가기’ 링크가 있고 ‘홈에서’는 없다. agent는 4개 그대로.
- F3(키 문구): a·b의 /home/agents 어디에도 ‘키 키’ 중복이 없다. ‘기존 키 등록됨 · 새 키 등록 전’은 키가 있는 agent의 새 시도에만 나오고 fixture에 그런 agent가 없어 문구는 소스(`connect_guide.go:56`)와 DEV 통합 검사(`connect_guide_integration_test.go`)로만 확인했다(not_run on fixture).
- F2(요약 없는 gate 버튼 없음): 소스 `http.go:578` `view["Active"] = active && summary != nil`. fixture가 `/gate`를 막아 화면은 못 열었고 DEV도 도달 경로가 없어 새 검사가 없다고 보고했다(소스 근거만).

### 이전 PASS 핵심 회귀 — PASS

5선택지·문구·보조 문구·선택 전 미리 선택 없음·위조 query 무시, 다섯 선택 안내(4번 `/knowslink-session on`·‘앱에서 확인 필요’, 5번 고정 문장·버튼 없음), OpenClaw·Meta Muse·온라인 배지 없음, 재개 행 전부(waiting/prepared/approved/키 등록/만료/취소/승인 뒤 미저장/remote 3종), 요약 줄(‘보낸 친구 요청 N건’은 이제 화면의 대기 행 수에 맞춰 검사), O-NM-03(있는/없는 코드 같은 모양·대소문자 유지·형식 다름·받는 쪽 비노출·전체 코드 비노출), 연결 화면 7개 단계·복귀·접근성, 중복 생성(같은 attempt 동시 2회 → 1개, 재전송 같은 위치, 새 attempt 한도 거절, 기존 키 카드 불변, 다른 회원에 비노출), 5분 재확인. 1440·390·320 가로 넘침 0·잘림 0·44px 미만 0(Chrome·Edge 각 76 PASS), 키보드 Tab·Enter·화살표 PASS, 요청은 모두 같은 origin GET뿐.

**브라우저 재실행 1건**: Edge keys 첫 실행에서 `conn-waiting`·`conn-prepared`의 ‘복귀 링크가 자기 에이전트 카드에 도착’ 검사가 FAIL(스크롤 위치를 앵커 스크롤 전에 측정, 첫 attempt와 같은 시점 경쟁, [첫 실행 로그](edge-keys-attempt2-first-run-return-link-scroll-timing.log)). 검사를 약화하지 않고 재실행해 33/33 PASS. 같은 시점 경쟁이 두 attempt 연속 나므로 도구 안정화(스크롤 완료 대기)를 후속으로 권한다. 제품 결함 증거는 없다.

### 관찰과 남은 한계

- **O-CT-05 (low, 정보) 종료 행 3개 위치와 ID 순서.** 세 코드의 위치(3·4·5)는 경로 사이에 달랐지만(예: `receive_limit` 3, `blocked` 4, `unknown_cancelled` 5) 재요청·취소·재요청 전후 모두 같았다. 순서가 코드의 존재·차단 여부를 따르는지는 이 fixture(경로당 1행)로는 판단할 수 없다. 정렬 기준이 서버 ID라면 무작위라 문제없다. 이번에는 확인하지 못했다(소스 정렬 키 미읽음).
- 기존 데이터(이전 차단 경로의 쌍 행+대리 행 2행, 이전 코드로 되돌리기 시 보류 행 노출)와 DEV가 적은 남은 한계 4건은 fixture에 없어 not_run이다.
- not_run: 스크린리더 음성, 실제 브라우저 200%, 대비 측정, remote 동의·앱 복귀 화면(401), 승인 목록·익명 429 화면, 서비스 재시작 뒤 복구, 실제 호스트 설치·왕복, Windows tester 워크트리에서 Go 테스트 재실행(DEV 원격 Linux EXIT=0 재사용).

### 상태 변경과 정리 (attempt 2)

fixture에 만든 것: A의 보낸 요청 7행이 모두 ‘대기’(세 코드는 재요청으로 새 24시간 기한, generation 3), 종료 행 3개(`qa-cap-fake-0001~0003` 취소), A의 agent 1개(`Connect QA 신규 에이전트`, ID는 [mutate-result.json](mutate-result.json), 5/5). 비가역이다. D·E·B·C 화면은 바뀌지 않았다. 재요청 경로(`rerequest_codes`)는 이 fixture에서 다시 쓸 수 없다. 브라우저 프로필은 session 임시 디렉터리에 두었고 headless 프로세스는 종료했다.

---

## 이전 결과 (8d8bad4, 제품 37febf9)

## 대상과 결론

대상은 fullops/dev `8d8bad4cd75a37546ba0f803e3c0c8290ddb6a87`(제품 `37febf9`)이다. dev 워크트리 HEAD가 이 SHA임을 읽기 전용으로 확인했다. 기준 main `4d15c77`(지시서 base `d3e676a`). fixture `http://127.0.0.1:18932`(manifest `connect-ux002-dev/fixture-manifest.json`)이다. 서버가 이 SHA인지는 coor 지시와 manifest를 따랐고 tester가 독립 확인하지는 않았다. 제품 코드는 고치지 않았다. 테스트 레벨은 lite다.

**결론: critical·high 결함 없음. medium 결함 1건(D-CT-01). low 관찰 4건(O-CT-01~04).** 스크린리더·실제 브라우저 200%·remote 앱 복귀/동의 화면은 not_run이다. 전체 접근성 완료나 호스트 지원 완료로 보고하지 않는다.

| 영역 | 실행 | 결과 |
|---|---|---|
| HTTP 읽기 전용 [checks.py](checks.py) → [http-result.json](http-result.json) | 89 | 88 PASS, **1 FAIL (D-CT-01)** |
| 5분 재확인 [stale.py](stale.py) → [stale-result.json](stale-result.json) | 1 | PASS (422, agent 4개 그대로) |
| 생성 [mutate.py](mutate.py) → [mutate-result.json](mutate-result.json) | 8 | 8 PASS |
| Chrome 155.0.8059.80 layout [json](browser-chrome-layout.json) / keys [json](browser-chrome-keys.json) / after [json](browser-chrome-after.json) | 76 / 33 / 8 | 모두 PASS |
| Edge 155.0.4283.62 layout [json](browser-edge-layout.json) / keys [json](browser-edge-keys.json) / after [json](browser-edge-after.json) | 76 / 33 / 8 | 모두 PASS |

실행 명령(종료코드는 파이프 없이 확인): `python -I -X utf8 checks.py`(마지막 실행 종료코드 1, D-CT-01 때문), `python -I -X utf8 stale.py`(0), `python -I -X utf8 mutate.py`(0), `CAP=1 node browser.mjs 9333 chrome layout`(0), `node browser.mjs <9333|9334> <chrome|edge> keys`(0), `node browser.mjs <port> <label> after`(0). 브라우저는 `--headless=new --disable-extensions`, 비모바일 viewport 1440×900·390×844·320×640(`innerWidth` 측정), 실제 키 입력은 CDP `Input.dispatchKeyEvent`다. 시나리오는 [SAR-UX-ISSUES-002-CONNECT-TESTER.md](../../scenarios/SAR-UX-ISSUES-002-CONNECT-TESTER.md)다.

## 첫 실행 실패 보존 (검사·환경 결함, 제품 결함 아님)

| 로그 | 원인 | 처리 |
|---|---|---|
| [http-result-first-run-expectation-errors.json](http-result-first-run-expectation-errors.json) | 내 기대값 결함 3건: ‘연결 취소’ 부분 문자열이 제목 ‘연결 취소되었습니다’와 겹침, 취소 버튼 이름 정규식이 ‘입력한 코드’ 접두를 빠뜨림, legacy ‘이 화면에서 찾기’는 `aria-label`이라 본문 텍스트에 없음. 4번째 FAIL(D-CT-01)은 실제 결함이라 유지 | 기대값만 고침. D-CT-01은 그대로 FAIL |
| [chrome-keys-first-run-429-on-click-navigation.log](chrome-keys-first-run-429-on-click-navigation.log) | 클릭 이동이 `go()`의 429 재시도를 거치지 않아 회원당 분당 40회 한도 화면을 측정 | `enterNavRL`(429면 기다린 뒤 reload)로 고침 |
| [edge-keys-first-run-return-link-scroll-timing.log](edge-keys-first-run-return-link-scroll-timing.log) | waiting 복귀 뒤 카드 top 1825(앵커 스크롤 전 측정). 재실행은 PASS라 시점 경쟁으로 본다 | 검사는 약화하지 않고 재실행 |
| [chrome-after-first-run-expectation-error-create-returns-to-agents-card.log](chrome-after-first-run-expectation-error-create-returns-to-agents-card.log), [-unused-wording.log](chrome-after-first-run-expectation-error-unused-wording.log)(edge 동일) | 생성 복귀 위치와 신규 카드 문구를 내가 잘못 가정(생성은 `/home/agents#agent-{id}`로 복귀, 카드는 ‘미사용 · 등록된 키 없음’) | 설계 문장(생성은 새 agent 위치로 복귀)에 맞게 기대값 수정 |

## 결함

### D-CT-01 (medium) 승인 뒤 클라이언트 저장이 확인되지 않은 연결 화면이 ‘새 키가 연결되지 않았습니다 · 새로 시작하세요’라고 단정한다

- 근거: U8 표 ‘승인 후 응답 유실’(‘결과가 불확실하면 확인 필요. 맹목적 새 등록/기존 키 덮어쓰기 금지’), ‘local approved … 승인만으로 연결 완료를 표시하지 않음’.
- 재현: `POST /fixture/as/b` → `GET /home/agents`의 ‘합성 저장 확인 도우미’ 카드는 “연결 상태: 확인 필요 · 승인 뒤 클라이언트 저장을 확인하지 못함 / 시작한 클라이언트의 저장 결과를 먼저 확인하세요”다. 같은 agent의 연결 화면 `GET /home/connections/YkFECkExYLDoQ43m7qGtUgsj_A5Ke3YmfEQ6HjKYGJI`(manifest `connections.device_lost`)는 “만료되었습니다. 새 키가 연결되지 않았습니다. 이전 연결 수단은 다시 쓸 수 없습니다. 내 에이전트의 연결·키 관리에서 새로 시작하세요.”다.
- 기대: 두 화면이 같은 ‘확인 필요’(클라이언트 저장 결과를 먼저 확인)를 안내한다. 서버는 클라이언트가 키를 저장했는지 알 수 없다.
- 실제: 연결 화면은 만료·취소만 구분해 승인 여부를 쓰지 않고 키가 없다고 단정하며 새 시작을 권한다. 클라이언트가 이미 저장했다면 사용자는 새 키를 만들 수 있다(키 한도 3개·기존 키 불변은 DEV 검사가 보호).
- 원인(소스 읽기): `internal/relay/member_pages.go:68`의 템플릿이 `expired`/`cancelled` 두 상태만 보고 ‘승인 뒤 만료’를 구분하지 않는다. 카드는 `connect_guide.go:67`이 구분한다.
- SHA: dev `8d8bad4`(제품 `37febf9`). 증거: [http-result.json](http-result.json) `device_lost_connection_page_agrees_with_card_unknown_state`, 화면 `chrome-conn-device-lost-{1440,390,320}.png`.
- 담당: dev(제품 코드). tester는 고치지 않았다.

## 관찰 (결함 아님, 판단 대상)

- **O-CT-01 (low, 제품 판단) 차단된 요청의 기한 뒤 재요청이 행 수 차이를 남긴다.** DEV 질문과 같다. fixture에는 차단 경로가 없어 재현하지 못했고 소스로만 확인했다. `internal/relay/human.go:171-187`: 코드가 없으면 `undisclosed:`+해시 자리표시 행을 재사용(같은 행이 다시 대기)한다. 실제 회원이 차단했으면 만료 뒤에도 옛 행이 `blocked-pending`으로 남고(`sweepFriends`는 `pending`만 만료) 재요청은 새 자리표시 행을 만든다. 요청자가 기한(24h) 뒤 같은 코드를 다시 보낼 때 행 1개(없는 코드)와 2개(차단한 회원) 차이가 생길 수 있다. 정보량은 ‘이 코드의 회원이 나를 차단했다’ 한 비트다. coor/designer가 수용·수정 여부를 정한다.
- **O-CT-02 (low) legacy ‘이 화면에서 찾기’는 보이는 글자가 없다.** `/home/legacy`의 목차는 `<nav class="toc" aria-label="이 화면에서 찾기">`이고 `aria-current`는 없다(거짓 현재 표시 없음, PASS). 보이는 제목이 없어 눈으로 보는 사용자에게는 구분 문구가 나타나지 않는다. 설계는 ‘구분한다’이며 DEV 변경 없음(기존 상태).
- **O-CT-03 (low) 에이전트 카드 제목은 3줄 뒤 ‘…’로 줄여 보인다.** 40자 긴 이름(‘긴한글이름’×8)이 390px에서 잘려 보인다(`.name-preview` 기존 U2 미리보기 패턴, 전체 이름은 ‘이름·연결 관리’ 안에 있다). 가로 넘침은 0이다. 이번 단위 변경이 아니다.
- **O-CT-04 (low) 재확인 거부 안내가 ‘홈에서 확인’이라고 쓴다.** 5분이 지난 세션의 생성 POST는 422 ‘이 작업은 5분 안의 로그인 신원 재확인이 필요합니다. 홈에서 로그인 신원을 다시 확인하세요.’와 ‘로그인 다시 확인’·‘내 에이전트로 돌아가기’를 보인다. 동작은 맞다(agent 0 증가). 문구가 ‘홈’이라 내 에이전트에서 시작한 사용자에게는 한 단계 멀다.

## 확인한 내용 (지시서 항목별)

**‘어디에서 사용할까요?’ (U8).** 선택 전에는 안내·미리 선택이 없다. 다섯 선택지의 순서·문구·보조 문구가 설계표와 같다(web·grok·dots·claude·other, 값 포함). 선택 뒤 안내가 준비물·권한·다음 행동을 보이고 README 앵커는 dev 워크트리 README의 실제 제목과 맞다(`Grok 신규 설치`, `Dots 신규 연결`, adapters `Claude Code 작업 세션 연결 — SAR-SESSION-PEER-001-DEV`). 4번은 `/knowslink-session on`과 ‘앱에서 확인 필요’(`status`, `knowslink_text_receive`), ‘기존 연결 계속’은 새 에이전트·새 폴더를 만들지 않는다고 안내한다. 5번은 고정 문장과 웹 링크만 있고 버튼·코드·명령·배지·호스트 이름이 없다. 모든 변형 화면에 OpenClaw·Meta Muse 없음, ‘온라인’은 ‘확인하지 않습니다’ 부정문에만 나타나고 배지·단정 없음. 위조 query(`use=openclaw<b>&start=zzz`, 중복 `use`)는 무시되고 반영되지 않는다. 선택·새로고침·뒤로가기·키보드 제출로 agent 수가 변하지 않는다.

**상태별 재개 행.** a: waiting(‘키 등록 전 · 내 컴퓨터의 준비 대기’, 연결 수단 미재표시), prepared(‘지문 확인 필요’), approved(‘브라우저 동의 완료 · 클라이언트 저장 확인 필요’, ‘승인만으로는 연결되지 않습니다’), 키 등록됨(‘앱에서 연결 상태를 확인하세요’, 기존 키 덮어쓰지 말 것). b: 만료·취소(‘이전 시도 … · 키 없음’, 재사용 불가), 승인 뒤 미저장(‘확인 필요’, 자동 철회 없음), remote 앱 복귀 대기(‘앱으로 이동한 것만으로 연결되지 않습니다’), 교환 유실(‘확인 필요 · 앱의 토큰 교환을 확인하지 못함’), 교환 완료(‘키 등록됨 · 원격 앱의 토큰 교환 확인’). 완료·온라인 단정 없음. 카드에 명령·키 비밀값 없음. 요약은 ‘보낸 친구 요청 4건’을 정확히 센다. 다만 승인 뒤 미저장은 연결 화면과 모순이다(D-CT-01).

**연결 화면과 복귀 (U7).** 연결 화면 7개 중 waiting/prepared/approved/키 등록됨은 서버 상태의 한 단계만 `aria-current=step`(1·2·3·3), 만료·취소·승인 뒤 미저장은 단계가 없다. 키 등록됨은 ‘키 등록됨’ 라벨이고 ‘연결 취소’가 없다. 다른 회원의 연결 화면은 열리지 않고 익명은 로그인으로 이동한다. 네 개의 복귀 링크는 Enter로 각자의 `/home/agents#agent-{id}` 카드에 도착하고 스크롤된다. title·h1(1개)·main·skip·`lang=ko`, 서로 다른 title, 현재 메뉴 표시(`aria-current=page`는 내 에이전트·친구만)가 PASS다. Tab 순회에 함정이 없고 모든 정지점에 이름·보이는 focus·44px이다(선택기 화살표·묶음 순환·Enter 제출 포함, 제출 뒤 focus는 `use-guide`).

**새로고침·뒤로가기·중복 클릭.** `mutate.py`로 같은 attempt를 동시에 2회 제출해 둘 다 303(같은 위치), agent는 정확히 1개 증가(4→5). 나중의 재전송도 같은 위치로 가고 새 agent가 없다. 새 attempt의 여섯째 생성은 거부되고 agent는 5개로 남는다. 기존 키 카드(`71513175`)는 그대로, 새 카드는 ‘미사용 · 등록된 키 없음’이며 다른 회원에게 보이지 않는다. 브라우저 back/forward/reload로 카드 수가 변하지 않는다. 5분 재확인 창 밖(310초 경과 세션)의 생성은 422로 거부되고 agent 수가 같다.

**O-NM-03.** a의 보낸 행 5개: ‘…XXK2’(실제 코드, 대소문자 유지)·‘…zz99’(없는 코드)는 같은 모양, ‘형식 다름’(`ab`, 원문 미표시), ‘이전 요청’, 취소 행 ‘…ca11’(종료 뒤에도 유지, ‘새 친구 요청과 수락이 필요합니다’). 취소 버튼 접근성 이름은 ‘입력한 코드 …XXK2 친구 요청 취소, 기한 …’처럼 끝 4자를 포함하고 서로 다르다. 받는 쪽(b)·a가 받은 c의 요청 카드에는 ‘입력한 코드’나 끝 4자가 없다. b가 가진 ‘…XXK2’는 b 자신의 코드뿐이다. 요청자 화면과 URL·링크·hidden 값에 전체 코드가 없다. c의 보낸 행은 ‘…Z-3N’이고 a의 전체 코드는 없다. 있는/없는 코드 구분 검증은 이 fixture의 대기 행 2개 비교까지이며 차단·상대 수신 한도·자기 코드 행은 not_run(fixture에 경로 없음, DEV 통합 검사가 근거)이다.

**레이아웃.** 1440·390·320에서 17개 화면(a 10, b 5, c 2, 연결 화면 7개 포함) 가로 넘침 0, 잘림 0, 44px 미만 동작 0(Chrome·Edge 각 76 PASS). 320px에서 ‘입력한 코드 …ca11’은 한 줄로 유지된다. 캡처: `chrome-*.png` 57장(designer용 3폭, 파일명 `chrome-<화면>-<폭>.png`).

**low 3건(소스·Go 테스트 근거, fixture에서 화면 열 수 없음).** 익명 429: `member_pages.go:98`의 `{{if not .Anonymous}}`가 ‘내 에이전트’ 링크를 감싼다. 승인 목록: `member_pages.go:137`의 `aria-label="요청 검토: 종류 / 요청한 에이전트 → 받은 내 에이전트 이름 / 기한"`, 검증 `gate_summary_integration_test.go:77`. `gate_summary.go:47`부터 `v, isMap :=`·`minutes, isNumber :=`로 comma-ok. DEV 원격 Linux make test/integration EXIT=0(같은 제품 SHA)을 재사용했고 이 Windows tester 워크트리에서 Go 테스트를 다시 돌리지 않았다.

## not_run

- 스크린리더 음성, 실제 브라우저 200% 확대(viewport/CSS 확대는 대체가 아니다), 대비 측정(이번에 새로 측정하지 않음, DEV 렌더 probe만).
- remote 동의·앱 복귀 실패 화면(`/connect/*`는 유효한 remote 거래 없이 401): 소스 템플릿 `remote_auth.go:314`(‘앱으로 돌아가는 중’, ‘이 화면은 연결 완료를 뜻하지 않습니다’, 링크와 1분 만료 안내)와 DEV 렌더로만 확인.
- `/home/cancel`·`/home/confirm`·`/home/friends/request`·`decision` 실행(공유 fixture 상태를 바꾸므로 하지 않음). 취소 뒤 복귀·승인 뒤 복귀는 DEV 통합 검사가 근거이다.
- 차단·상대 수신 한도·자기 코드 행의 동일 모양, 서비스 재시작 뒤 복구, 실제 호스트 설치·왕복(범위 밖).
- 승인 목록·익명 429 화면(fixture가 `/gate`·`/home/receipts` 비허용).

## 상태 변경과 정리

fixture에 만든 것은 a의 agent 1개(`Connect QA 신규 에이전트`, ID는 [mutate-result.json](mutate-result.json))다. 남은 1칸을 쓴 비가역 변경이며 designer가 보는 a의 `/home/agents`에 카드 1개(5/5)가 늘었다. 거절된 POST(재확인 창, 여섯째)는 아무것도 만들지 않았다. 브라우저 프로필은 session 임시 디렉터리에 두었다.
