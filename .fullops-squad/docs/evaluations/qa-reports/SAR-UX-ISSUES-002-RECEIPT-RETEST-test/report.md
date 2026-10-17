---
title: 조회 선택기 구분값 선두 표시 재검증
status: review
updated: 2026-10-17
owner: tester
tasks: [SAR-UX-ISSUES-002-RECEIPT-RETEST]
summary: D-RT-01의 구분값 가림은 해소됐다. 변경 없는 선택기 검사 세 건은 그대로 FAIL이고 종류 잘림은 남는다
---

# SAR-UX-ISSUES-002-RECEIPT-RETEST 재검증

고정 소스는 `71dcac7822b84acc2ae458a533dda96af1e94fe4`(제품 source `da53e96`)다. 제품 코드는 고치지 않았다. critical·high·새 결함은 없다.

## 요약

- **이전 선택기 검사 세 건(320/390/1440)은 기대를 바꾸지 않고 다시 실행했고 모두 FAIL이다.** Chrome 155와 Edge 155가 같다. 검사는 "열린 이름+종류+구분값 전체 문자열이 닫힌 선택기 폭에 들어가는가"를 본다. 전체 폭 902px는 가용 750/293/223px를 넘는다. 새 표시 순서를 기대하지 않았으나 이 검사의 판정 기준은 순서와 무관해서 실패 이유가 같다.
- **D-RT-01의 핵심 영향은 해소됐다.** 구분값이 모든 폭·두 브라우저에서 닫힌 상태에 보이고 `agent_a`와 `agent_a2`는 첫 토큰이 다르다. 동명 두 agent를 닫힌 선택기에서 구별할 수 있다.
- **남는 차이.** 옵션 순서가 `구분값 / 이름 / 종류`라 U2의 "이름 먼저"와 다르다. 종류는 1440에서 앞 3~4자(`Node`)만 보이고 390·320에서는 전혀 보이지 않는다. 이름도 390에서 13/39자, 320에서 8~9/39자만 보인다. 이 순서 예외와 종류 잘림을 수락할지는 designer가 판정한다.
- U6 회귀(입력 보존·없는/타인 동일 문구·계정 격리·429·503·세션·변조 쿠키·철회)는 통과했다. 변경 전과 다른 결과는 위 순서 단언 한 건뿐이다.

## 선택기 세 건 — 변경 없는 검사

| 검사 | 폭 | 가용 폭 | 전체 폭 | Chrome | Edge |
|---|---|---|---|---|---|
| selector-closed-shows-full-name-kind-distinguisher-1440 | 1440 | 750px | 902px | FAIL | FAIL |
| …-390 | 390 | 293px | 902px | FAIL | FAIL |
| …-320 | 320 | 223px | 902px | FAIL | FAIL |

`browser.mjs`는 이전 QA 파일을 그대로 복사했다(`browser-chrome.json`·`browser-edge.json`). 각각 51건이고 Chrome 3건, Edge 4건이 실패다. Edge의 네 번째는 이전과 같은 `requests-same-origin-get-only`다. 실패한 요청은 모두 Edge 내부 `origin: null` GET(`/icons/sprite.svg`, `/build/content.css`)이다. fixture origin의 요청은 전부 GET이다.

## 폭별 관찰 (`selector-observe.mjs`, 관찰 전용)

선택한 옵션 글자가 닫힌 선택기에서 몇 자 보이는지 같은 글꼴로 측정했다. 이름은 합성 39자 한글, 종류는 `Node 22 로컬 클라이언트`(16자)다. 측정값은 canvas 추정이며 실제 렌더 캡처와 대조했다(1440의 종류는 추정 3자, 캡처에서 `Node` 4자).

| 폭 | agent | 구분값 | 이름 | 종류 | 보이는 글자 |
|---|---|---|---|---|---|
| 1440 | agent_a | 보임 | 39/39 | 3/16 | `agent_a / 긴한글…긴한글 / Nod` |
| 1440 | agent_a2 | 보임 | 39/39 | 2/16 | `agent_a2 / 긴한글…긴한글 / No` |
| 390 | agent_a | 보임 | 13/39 | 0/16 | `agent_a / 긴한글×4+긴` |
| 390 | agent_a2 | 보임 | 13/39 | 0/16 | `agent_a2 / 긴한글×4+긴` |
| 320 | agent_a | 보임 | 9/39 | 0/16 | `agent_a / 긴한글긴한글긴한글` |
| 320 | agent_a2 | 보임 | 8/39 | 0/16 | `agent_a2 / 긴한글긴한글긴한` |

Edge는 같다. 390의 agent_a2만 가용 폭이 308px라 이름 14자가 보인다(Chrome 293px). 옵션 HTML은 `<option value="agent_a" selected>agent_a / 긴한글… / Node 22 로컬 클라이언트</option>`다.

## U6 회귀 (짧게)

- HTTP `checks.py` 16건(이전 파일 그대로): 15 PASS, 1 FAIL. 실패는 `first_view_form_and_selector`의 `opt.endswith("agent_a")`다. 이 단언은 구분값이 뒤에 오는 이전 순서를 기대한다. 기대를 바꾸지 않았다. 새 순서의 옵션(`agent_a / 이름 / Node 22 …`)은 구분자가 2개이고 `Node 22`를 포함하며 첫 토큰이 구분값이다. 이 단언 외의 같은 검사(옵션 목록 `["", agent_a, agent_a2]`, agent_b/agent_c 비노출, 폼 GET, autofocus 1곳 등)는 통과했다.
- 입력 보존·없는/타인/비참여 동일 문구·계정 B·익명 격리·변조 쿠키·429(입력 유지, 존재 oracle 없음)·503·세션 만료 후 재로그인: PASS.
- 브라우저(Chrome·Edge `browser.mjs`): 선택기 외 키보드·focus·a11y·320/390/1440 레이아웃·대비·429 레이아웃·details 검사가 모두 통과했다.
- 철회(맨 마지막): `/fixture/revoke` 뒤 `revoke_check.py` 3건과 `revoke-browser.mjs` 4건 PASS. 철회된 agent는 `disabled`·`철회`로 남고 agent_a2로 바뀌지 않으며 B에게 비노출이다. 옵션 문자열이 바뀌어도 정규식 `… / 철회</option>`이 그대로 맞는다.

## 리뷰 low 2 — 익명 한도 429 문구 (관찰만)

fixture에서 세션이 없거나 무효 세션 쿠키로 `/home/receipts`를 31번 보내면 30번 303(로그인 이동) 뒤 31번째가 429다(`anon429-result.json`, `anon429-invalid-cookie-result.json`). 두 경우 모두 응답에 "조회할 내 에이전트가 없습니다"는 없다. 실제 문구는 "요청이 많습니다. … 이후 다시 시도하세요. 저장된 입력은 로그인 확인이 가능해지면 복원됩니다."이고 선택기는 빈 옵션 하나, 링크는 `/home/agents`("내 에이전트로 돌아가기")와 `href=""`인 "로그인 확인 후 다시 조회"다. 리뷰가 설명한 문구는 이 경로에서 재현하지 못했다. 다른 분기(코드 `member_receipt.go:27-31`)가 있는지는 소스를 읽지 않아 판단하지 않는다. 관찰: 익명에게 "저장된 입력 복원" 문구와 빈 `href` 링크가 보인다. 데이터 노출은 없다. 결함으로 올리지 않고 designer·DEV 판단 자료로 남긴다.

## 판정

| 항목 | 결과 |
|---|---|
| 선택기 세 건(변경 없는 검사) | FAIL 3/3 (Chrome·Edge) |
| 구분값이 닫힌 상태에서 보임, agent_a/agent_a2 구별 | PASS (관찰, 모든 폭·두 브라우저) |
| 종류 표시 | 1440 일부, 390·320 없음 — designer 판정 필요 |
| U2 "이름 먼저" | 순서 예외 — designer 판정 필요 |
| U6 핵심 회귀 | PASS (순서 단언 1건만 새 순서로 인해 FAIL, 기대 유지) |
| D-RT-01(medium) | 구분값 가림은 해소. 완전 수락은 designer 시각 판정 후 |

D-RT-01을 PASS로 닫지 않는다. 변경 없는 검사가 실패이고 U2의 순서·종류 표시 조건이 designer 결정에 걸려 있다. 새 DEV 결함은 없다.

## 실행하지 않은 항목 (not_run)

- 열린 선택 목록의 표시(네이티브 popup, headless 캡처 불가). 열린 목록에서 이름·종류가 전부 보이는지는 미확인이다.
- 보존 종료·전달 기한 만료 receipt, 명시적 로그아웃 쿠키 삭제(fixture가 허용하지 않음), 브라우저 UI의 실제 200% 확대(모의 720px/DPR 2만), 스크린리더, Firefox·Safari, 실제 이메일·Google 로그인.
- DEV `make lint/test`는 재실행하지 않았다. 기존 증거(da53e96 Linux EXIT=0, 71dcac7 직접 실행 아님)를 쓴다.
- Jev 모델 시나리오는 실행하지 않았다. 입력 검사가 결정론적이다.

## 실행 경계와 상태

- origin `http://127.0.0.1:18928`, fixture는 같은 SHA다. 합성 A/B만 썼다. 실제 이메일·계정·친구·메시지·운영 서비스는 건드리지 않았다. 허용한 GET과 fixture 전용 POST(`as/a|b`, `limit`, `recover`, `outage`, `expire`, `revoke`)만 보냈다.
- **fixture는 이제 revoke 상태다.** 이 fixture로 철회 전 흐름을 다시 재현할 수 없다.
- 브라우저는 새 임시 프로필의 Chrome 155·Edge 155 headless CDP이고 종료했다.
- 익명 429 관찰로 fixture 서버 쪽 익명 한도를 두 번 채웠다. 약 1분 뒤 해제됐다.

## 증거 파일

`browser-chrome.json`·`browser-edge.json`·`chrome.log`·`edge.log`, `http-result.json`(`checks.py`), `selector-observe-chrome.json`·`selector-observe-edge.json`(`selector-observe.mjs`), `anon429*.json`(`anon429.py`), `revoke-result.json`·`revoke-browser-result.json`·`revoke-browser.log`, 캡처 `*-selector-closed-agent_a|agent_a2-{1440,390,320}.png`(닫힌 선택기 확대 2배)와 `*-viewport-agent_a-{1440,390,320}.png`, 이전 파일과 같은 성공·거부 캡처. 명령은 `python3 -I -X utf8 checks.py`, `node browser.mjs <port> <chrome|edge>`, `node selector-observe.mjs <port> <label>`, `python3 -I -X utf8 anon429.py`, 철회 뒤 `python3 -I -X utf8 revoke_check.py`와 `node revoke-browser.mjs <port>`다. 종료코드는 `checks.py` 1, Chrome 1, Edge 1, 나머지 0이다.
