---
title: 연결 안내·재개·보낸 요청 코드 끝 4자 독립 QA
status: review
updated: 2026-10-20
owner: tester
tasks: [SAR-UX-ISSUES-002-CONNECT-TESTER]
summary: "8d8bad4(제품 37febf9) 격리 fixture 독립 QA: U8 선택지·재개 행·중복 생성·5분 재확인, U7 연결 화면, O-NM-03 PASS. critical·high 없음, medium 1건(D-CT-01 승인 뒤 미저장 시도의 연결 화면 안내가 카드와 모순), low 관찰 4건"
---

# SAR-UX-ISSUES-002-CONNECT-TESTER 독립 QA

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
