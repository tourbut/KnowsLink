---
title: 요청 상태 조회 복구의 독립 QA
status: review
updated: 2026-10-17
owner: tester
tasks: [SAR-UX-ISSUES-002-RECEIPT-TESTER]
summary: U6 조회 복구와 계정 격리는 통과하고 닫힌 선택기의 구분값 잘림 한 건을 medium으로 보고한다
---

# SAR-UX-ISSUES-002-RECEIPT-TESTER 독립 QA

U6의 입력 보존·안전 복귀·계정 격리와 U7의 키보드·focus·320/390/1440 흐름은 통과했다. critical·high 결함은 없다. medium 한 건이 남는다. 긴 이름의 에이전트는 닫힌 선택기에서 이름·종류·구분값 중 구분값이 잘린다(U2). 이 결함이 닫히기 전에는 U2 이름 표시 수락을 PASS로 보고하지 않는다. 제품 코드는 고치지 않았다.

## 입력과 실행 경계

- 고정 SHA는 `9ae29ec4ddfcfacb59ed48871a14bb558342c26e`다. dev 워크트리 HEAD가 같다. 제품 source `1e49033`과 이 SHA 사이의 변경 21개는 모두 `.fullops-squad/` 아래다. 제품 파일은 같다.
- origin은 `http://127.0.0.1:18928`이다. 합성 A/B/C와 manifest의 container `kl-receipt002-fixture-9ae29ec4ddfc-6620c395`만 썼다. 실제 이메일·계정·친구·메시지·운영 서비스는 건드리지 않았다.
- 허용한 요청은 GET `/home/receipts`, `/`와 fixture 전용 POST(`as/a|b`, `limit`, `recover`, `outage`, `expire`, `revoke`)다. 브라우저 요청은 fixture origin의 GET뿐이다.
- DEV의 `make lint/test`는 이 SHA에서 exit 0이다. DEV의 실제 Linux/Postgres 검사는 source `1e49033`에서 나왔다. 두 증거를 DEV 보고에서 가져왔고 이번 독립 실행으로 바꾸지 않는다.
- `fixture/revoke`는 되돌릴 수 없다. 모든 다른 검사 뒤에 한 번 실행했다. 이후 상태는 별도 스크립트로 검사했다.

## 실행기와 사유

- HTTP: `checks.py`(16개). 쿠키를 직접 다루는 Python 클라이언트다. 변조 쿠키와 계정 교대는 브라우저로 만들 수 없는 값이다.
- 브라우저: `browser.mjs`(51개). 설치된 Chrome 155와 Edge 155를 CDP로 실행했다. 실제 키 이벤트(`Input.dispatchKeyEvent`)로 Tab·Shift+Tab·Enter·Space·방향키를 입력했다.
- 기본 Orca 브라우저는 탐색·snapshot·cookie 확인에만 썼다. 키보드 검증에는 쓰지 못했다. Orca 탭은 `document.hasFocus()`가 false이고 Tab 입력이 focus를 옮기지 않았다. 페이지 CSP가 `default-src 'none'`이라 페이지 안 fetch로 fixture POST도 할 수 없었다. 이전 과제에서 승인된 Chrome CDP 대안을 같은 사유로 재사용했다.
- Jev 모델 시나리오(`jev_test_web.py`)는 실행하지 않았다. 모델 키가 필요하고 위 입력 검사가 결정론적이다.

## 통과한 수락 조건

| 조건 | 근거 |
|---|---|
| 처음·성공·형식 오류·없는/타인/비참여 종료에서 같은 폼과 내 agent 선택 유지 | HTTP·Chrome. 라벨은 ‘내 에이전트’·‘요청 ID 또는 답장 ID’·‘상태 확인’이다. 오류 뒤 입력이 남고 같은 화면에서 수정·재조회했다. |
| 없는 ID·타인 간 ID·비참여 내 agent의 구분 불가 | 세 경우 모두 403, 같은 문구와 같은 마크업(ID만 다름). agent_b/agent_c·타인 ID를 본문에 싣지 않는다. |
| 소유하지 않은 agent의 URL | `agent=agent_b|agent_c`는 선택·ID를 반사하지 않는다. |
| 안전 길이·escape | `"><script>`는 escape되어 422로 남는다. 129 바이트(한글 43자)·NUL·개행·RLO·isolate는 ID만 비우고 선택은 유지한다. |
| 성공 화면 | 전달·처리·관련 답장·기한을 분리한다. 폼은 GET이고 POST 폼이 없다. 반복·새로고침한 결과가 같다. 읽음·실행 완료는 아니라는 문구가 있다. ‘새 key로 송신’은 없다. |
| 새로고침·뒤로·복귀 | Chrome에서 새로고침·뒤로·‘조회 입력으로 돌아가기’(입력으로 focus 이동)·기본 진입 복원을 확인했다. |
| 429 | 입력 유지, 재시도 시각 표시, 자동 refresh·script 없음. 한도 중에는 존재 여부 문구가 달라지지 않는다. 복구 뒤 정상 조회. 390/320에서도 확인했다. |
| 503 | 이전 입력과 agent를 표시하지 않는다. `no-store`다. 쿠키와 URL을 유지한다. 복구 뒤 복원한다. |
| 세션 만료·재로그인 | 만료 뒤 A의 내용이 없다. 같은 회원 A로 재로그인하면 복원한다. 다른 회원 B는 복원하지 못한다. |
| 계정 B·익명 | B 선택기에는 B의 agent만 있다. B가 A의 agent를 조회하면 거부된다. 익명이 A의 쿠키만 보내도 내용이 없다. |
| 변조된 복구 쿠키 | Member를 B로 바꾼 쿠키, 타 agent·미소유 agent, 500자 ID, HTML·bidi ID, 만료 시각, 추가 필드(Email·Result), 깨진 base64·거대 값을 넣었다. 계정과 agent를 넘지 못했다. |
| 쿠키 | `__Host-`, Secure, HttpOnly, SameSite=Strict, Path=/, Max-Age≤4200. 필드는 Member·Agent·ID·Exp뿐이다. 세션 token·결과·이름은 없다. |
| 승인되지 않은 return | `return`·`next`·`redirect`·`url` 등에 외부 URL·`javascript:`를 넣어도 Location·반사가 없다. 모든 링크가 내부다. |
| 철회 | 철회된 agent는 disabled·‘철회’ 표시와 오류 안내로 남고 agent_a2로 바뀌지 않는다. 키보드 Enter는 agent 선택 오류로 막히고 ArrowUp으로 직접 고른 뒤 같은 화면에서 재조회했다. B에게 비노출이다. |
| 키보드·focus | 첫 화면 autofocus는 선택기 한 곳이다. Tab 순서는 ID→상태 확인→링크다. 모든 focus가 3px 윤곽이고 44px 이상이다. skip link가 보이고 main으로 이동한다. 오류 때 focus는 해당 필드에 있고 `aria-invalid`·`role=alert`가 있다. 결과는 `role=status`다. details를 Space·Enter로 열고 닫는다. |
| 접근성 트리 | banner·main·navigation(‘현재 위치’)·h1·combobox·textbox·button 이름이 맞다. ‘요청 상태 조회’만 `aria-current=page`다. |
| 320/390/1440·모의 200% | 첫·성공(details 열림)·거부·긴 ID(100자) 네 상태에서 가로 넘침 0, 44px 미만 동작 0, 잘림 0. 모의 200%는 720px·DPR 2다. |
| 대비 | 본문·보조·링크·버튼 ≥4.5:1, 입력 테두리·오류 선 ≥3:1(계산값은 `browser-chrome.json`). |
| 두 번째 브라우저 | Edge 155에서 같은 51개를 실행했다. Chrome과 같은 Chromium 계열이다. |

## 결함

**D-RT-01 (medium, U2/U6): 닫힌 선택기가 긴 이름의 구분값을 가린다.**
- 재현: 합성 A로 `/home/receipts?agent=agent_a&id=0199a3f2-4c10-7a11-8b22-334455667788`을 연다. 선택기는 39자 이름과 `Node 22 로컬 클라이언트`, `agent_a`를 하나의 문자열로 보여 준다.
- 기대: 선택한 agent의 이름·종류·안정 구분값이 닫힌 상태에서도 구별된다. U2는 동명 agent를 구분값으로 구별하게 한다.
- 실제: 텍스트 폭 902px에 가용 폭이 1440에서 750px, 390에서 293px, 320에서 223px이다. 구분값 `agent_a`와 `agent_a2`는 잘려 보이지 않는다. 같은 이름의 두 agent를 닫힌 상태에서 구별할 수 없다. 이름은 최대 40자라 일반 입력으로도 생긴다.
- 영향: 열린 목록의 표시는 확인하지 못했다. 결과 화면의 기술 상세에는 agent ID가 있어 잘못 고른 선택은 발견할 수 있다. 데이터 노출·권한 문제는 아니다.
- 증거: `browser-chrome.json`·`browser-edge.json`의 `selector-closed-shows-full-name-kind-distinguisher-*` 세 건 실패, 캡처 `chrome-success-1440.png`·`chrome-refusal-320.png`. 고정 소스는 위 SHA다.
- 담당: DEV. 구분값을 앞에 두거나 이름을 줄이는 방식은 designer와 DEV가 정한다.

## 관찰

- 기본 폼 GET도 회원당 40회 한도에 들어간다(41번째가 429). 한도 안에서 폼은 계속 보인다. DEV 보고의 ‘40회 rate’와 맞고 UX 영향은 작다.
- 철회 뒤 복구 쿠키는 다음 진입에서도 같은 철회 안내를 반복한다. 새 조회로 덮어쓸 때까지 유지하며 다른 agent로는 바뀌지 않는다.

## 보존한 첫 실패 (제품 결함이 아님)

- 첫 HTTP 실행 여러 번은 내 검사의 결함으로 실패했다. A의 조회 수가 40회 한도를 넘었고 미소유 agent 쿠키·id_html·성공 화면의 정당한 agent_b 표시에 기대값이 틀렸다. `new()`가 조회 전에 recover하도록 하고 기대값을 제품 계약에 맞게 고쳤다. 제품 응답은 바꾸지 않았다.
- 첫 브라우저 실행은 select 필수 검증으로 막혔다. 선택기를 ArrowDown으로 고른 뒤 입력하도록 바꿨다. 두 번째 실행은 이전 실행의 복구 쿠키가 남아 첫 화면 검사가 어긋났다. 시작 때 쿠키를 지운다.
- Edge의 `requests-same-origin-get-only`는 실패로 보존한다. fixture origin 요청은 모두 GET이다. 나머지 74건은 Edge 내부의 `origin: null` 자원(`/icons/sprite.svg`, `/build/content.css`)이다. 이후 `revoke`가 되돌릴 수 없어 같은 상태로 재실행하지 않았다.
- 최종 HTTP 실행의 `revoked_agent_is_not_silently_replaced`는 내 기대가 틀려 실패했다(철회 agent는 마지막에 disabled로 표시됨). 결과는 `http-result.json`이다. 상세한 철회 검사는 `revoke_check.py`와 `revoke-browser.mjs`로 옮겼다. 그 뒤의 전체 재실행은 철회로 앞선 검사가 깨지는 것이 당연해 `http-result-invalid-rerun-after-revoke.json`에 무효 기록으로만 둔다.
- `revoke-browser.mjs`의 첫 실행에서 ArrowDown이 마지막 disabled 옵션에서 움직이지 않았다. 내 기대가 틀렸고 ArrowUp으로 바꿨다.

## 실행하지 않은 항목(not_run)

- ‘전달 기한 만료’와 보존 종료의 실제 receipt: fixture에 해당 ID가 없다. DEV의 Linux/Postgres 검사가 근거이며 이번 독립 확인은 아니다.
- 명시적 로그아웃의 쿠키 삭제: fixture가 `/auth/logout` POST를 허용하지 않는다.
- 브라우저 UI의 실제 200% 확대. 모의 720px/DPR 2만 실행했다.
- 스크린리더 음성, Firefox·Safari, 실제 이메일·Google 로그인과 복귀 redirect.
- 열린 선택 목록의 표시.
- 전체 접근성 완료를 주장하지 않는다.

## 증거 파일

`checks.py`·`client.py`·`http-result.json`, `browser.mjs`·`browser-chrome.json`·`browser-edge.json`, `revoke_check.py`·`revoke-result.json`, `revoke-browser.mjs`·`revoke-browser-result.json`, 캡처 `chrome-*.png`·`edge-*.png`(성공 1440/390/320, 거부 390/320). 캡처는 designer 시각 판정용으로 coor에 전달한다. 명령은 `python3 -I -X utf8 checks.py`, `node browser.mjs <port> <chrome|edge>`, `python3 -I -X utf8 revoke_check.py`, `node revoke-browser.mjs <port>`다. 종료코드는 HTTP 1(마지막 revoke 기대 결함), Chrome 1(선택기 3건), Edge 1(선택기 3건·내부 요청 1건), revoke 두 스크립트 0이다. 이 fixture는 revoke 때문에 다시 시작하지 않으면 전체 흐름을 재현할 수 없다.

## 다음 담당

DEV가 D-RT-01을 고친다. 고친 뒤 선택기 세 건과 필요한 경우에만 새 fixture로 재검증한다. designer 시각 판정과 고정 SHA 코드 리뷰, main 통합은 coor가 조정한다.
