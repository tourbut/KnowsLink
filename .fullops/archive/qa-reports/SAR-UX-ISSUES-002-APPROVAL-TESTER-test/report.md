---
title: 승인 요약·승인 화면 접근성 독립 QA
status: review
updated: 2026-10-20
owner: tester
tasks: [SAR-UX-ISSUES-002-APPROVAL-TESTER]
summary: "9d13a05(제품 f43d05c) 격리 fixture 독립 QA: U5 승인 요약·승인·거절·중복 클릭, U7 승인 화면, O-NM-04, O-NM-08 모두 PASS. critical·high·medium 없음, low 관찰 3건"
---

# SAR-UX-ISSUES-002-APPROVAL-TESTER 독립 QA

## 대상과 결론

대상은 fullops/dev `9d13a05ba93c7c085034213590ce661305606010`(제품 `f43d05c`, fixture 수정 `ea94128`)다. 기준 main `0c0fc96`(지시서 base `278001a`). fixture `http://127.0.0.1:18931`(manifest `approval-ux002-dev/fixture-manifest.json`)이다. 서버가 이 SHA인지는 coor 지시와 manifest를 따랐고 tester가 독립 확인하지는 않았다. 제품 코드는 고치지 않았다. 테스트 레벨 lite. 등록 `make lint`·`make test`·integration은 DEV 같은 SHA 결과(원격 Linux EXIT=0)를 재사용했다(Windows에 make 없음).

**결론: 모든 지시 check PASS. critical·high·medium 결함 없음. low 관찰 3건(아래 O-AP-01~03).** 스크린리더·실제 브라우저 200%·보조 브라우저 전 조합은 not_run이다. 전체 접근성 완료로 보고하지 않는다.

| 영역 | 실행 | 결과 |
|---|---|---|
| HTTP 읽기 전용 [checks.py](checks.py) → [http-result.json](http-result.json) | 105 | 105 PASS |
| 상태 변경 HTTP [mutate.py](mutate.py) 1·2 → [1](mutate-result-1.json), [2](mutate-result-2.json) | 27 + 8 | 35 PASS |
| Chrome 155.0.8059.80 layout [json](browser-chrome-layout.json) / keys [json](browser-chrome-keys.json) / flow [json](browser-chrome-flow.json) / after [json](browser-chrome-after.json) / onm04 / p409 / contrast | 60 / 22 / 5 / 17 / 8 / 4 / 8 | 모두 PASS |
| Edge 155.0.4283.62 layout / keys / after / onm04 | 60 / 22 / 17 / 8 | 모두 PASS |

실행 명령: `python checks.py`, `python mutate.py 1`, `node browser.mjs 9333 chrome flow`, `python mutate.py 2`, `node browser.mjs <port> <chrome|edge> <layout|keys|after|onm04|p409|contrast>`(종료코드 0, 파이프 없이). 브라우저는 `--headless=new --disable-extensions`, 비모바일 viewport 1440×900·390×844·320×640, 실제 키 입력은 `Input.dispatchKeyEvent`다. 시나리오는 [SAR-UX-ISSUES-002-APPROVAL-TESTER.md](../../scenarios/SAR-UX-ISSUES-002-APPROVAL-TESTER.md)다.

## 확인한 내용 (지시서 항목별)

**승인 목록 (U5).** C의 빈 목록은 설계 문구 그대로이고 행이 없다. A 목록은 7행이고 행마다 요청한 에이전트·받은 내 에이전트·요청 유형·절대 기한(KST)·상태·‘요청 검토’가 이 순서로 있다. 요청자는 허용된 이름이 없어 ‘확인할 수 없는 에이전트 / agent_a’, 받은 쪽은 이름·종류·구분값이다. 원요청 ID는 목록에 없고(상세의 기술 상세로 이동) gate ID는 href에만 있다. 상태는 대기 4·만료 1·승인 1·확인할 수 없음 1이고 대화 화면의 ‘승인 4건’은 대기만 센다.

**상세 (U5).** 읽기 순서는 누가 누구에게 → 요청한 내용 → 승인하면 허용되는 범위 → 기한·상태 → 승인/거절 → 기술 상세 → 복귀다. query는 확인할 기간·시간 단위가, commit은 일정·요청에 적힌 시간대가 서명된 원문 JSON과 KST 환산까지 일치한다. 고정 설명 ‘승인은 이 요청의 인간 게이트만 통과시킵니다…’가 범위 절에 있고, 정책 ID·원요청 ID·내부 종류는 기술 상세에만 있다. 악성 hint(`<script>`·URL·‘모든 권한 허용’)는 ‘참고 문구 — 승인 근거 아님’ 접힘 영역 1곳에 escape되어 있고 범위·상태·버튼에 오르지 않는다. commit 자유문(링크 포함 HTML)은 ‘요청자가 적은 내용’에만 escape되어 있고 링크·태그가 만들어지지 않으며 범위 절에 없다. 페이지에 `<script>`가 없다.

**닫힌 상태.** 만료·처리됨은 요약을 읽기 전용으로 보이고 form·버튼·CSRF가 없고 status ‘지금은 승인하거나 거절할 수 없습니다…’가 있다. 원문 부재는 요약·기술 상세 없이 alert ‘승인할 내용을 확인할 수 없습니다…’, 상태 ‘확인할 수 없음’, 복귀 링크만 있다.

**권한·재검사.** 다른 회원(B·C)의 GET은 403이고 이름·ID·본문이 새지 않는다. 존재하지 않는 gate도 같은 응답이다. 익명은 303이다. 다른 gate의 token(approve·deny 둘 다), 틀린·빈·잘린·없는 CSRF, 다른 세션 token, 알 수 없는 decision 값, 쿼리 문자열은 모두 상태를 바꾸지 못한다(403 또는 4xx, 두 gate 대기 유지).

**승인·거절·중복·뒤로가기 (마지막에 실행).**
- `pending_decide_1` A 승인: 303, 상세 ‘승인’·버튼 없음, 목록 행 ‘✓ 승인’, 대기 수 4→3. 같은 세션의 중복 승인·거절은 둘 다 409와 ‘이번 제출은 적용되지 않았습니다…’, 승인 목록 복귀(title·h1·alert 있음), 첫 결정 유지, 새 CSRF token 없음. B에게는 결정 뒤에도 403이다.
- `pending_decide_2`: Chrome에서 실제 키(`거절`에 focus, Enter) 1회 POST(`/home/gates/<id>/deny`) → 상세 ‘거절’·컨트롤 없음. 히스토리 뒤로가기는 `no-store`라 낡은 폼이 아니라 서버의 최신 닫힌 상태를 보인다. 목록 행 ‘거절’, 대기 수 2(HTTP 확인). 다른 세션의 stage-1 token 재제출은 approve·deny 모두 403과 같은 안전 문구이며 거절 상태가 유지된다.
- 요청 두 건 외의 gate(`pending_query_hint`·`pending_commit_long`)는 여전히 대기다.

**U7.** 승인 상세·403·원문 부재·만료·처리됨 모두 title `요청 승인 · KnowsLink`, h1 1, main 1, skip link(첫 Tab 정지점, 보이는 focus, Enter로 `#main`), 위치 nav(`aria-current=page`는 ‘요청 승인’ 하나), 하단 ‘승인 목록으로 돌아가기’다. 두 복귀 링크 모두 키보드 Enter로 `/home/legacy#gates`에 가고 `:target`이 승인 절이다. Tab 순서: skip → KnowsLink → 내 계정 → 승인 목록 → 참고 문구 → 이 요청 승인 → 거절 → 기술 상세 → 복귀(한 바퀴, 키보드 함정 없음). 모든 정지점에 접근성 이름, 보이는 focus, 높이 44px 이상. Space·Enter로 details 열고 닫아도 focus가 summary에 남는다. Shift+Tab이 거꾸로 이동한다. 접근성 트리에서 main 1·h1 ‘요청 승인’·nav ‘현재 위치’·heading 순서가 맞다. 결과 안내는 status, 오류는 alert다(원문 부재·403·409). 1440·390·320 모두 가로 넘침·잘림·말줄임·44px 미만 0이다(14개 화면 변형 × 3폭 × 2 브라우저). 640px viewport(실제 200% 등가 reflow)도 넘침 0이다. 40자 이름·긴 자유문·긴 JSON은 줄바꿈된다. 본문 대비 최소 6.19:1이다(고체 배경 텍스트 기준, 8개 화면).

**O-NM-08 PASS.** 받은 친구 요청에서 수락 → 거절 안내 → 거절 순서이고 안내는 Tab 정지점이 아니다. 거절 버튼의 `aria-describedby`가 안내 문구를 가리키고 접근성 설명이 같다. 안내 위 간격 16px이다. 수락·거절 이름에 상대 이름과 구분값이 있다.

**O-NM-04 PASS(템플릿 기준).** fixture가 `/home/receipts`를 403으로 막아 실제 익명 429 응답은 열 수 없다. 고정 SHA 소스를 `git archive`로 scratch에 풀어 `KNOWSLINK_UI_FIXTURE_DIR`로 렌더한 [render/receipt-anonymous-limited.html](render/receipt-anonymous-limited.html)을 브라우저로 열어 확인했다. 선택기·폼·fieldset·legend 0, 재시도 시각 있음, 실제 로그인 링크 `/|로그인 화면으로 이동`, 빈 `href` 0, ‘저장된 입력… 복원’ 문구 없음(소스상 로그인 회원에게만). 로그인 회원 429 렌더 [render/receipt-limited.html](render/receipt-limited.html)은 `#receipt-id`로 이동하는 링크이며 빈 href가 없다. 3폭 넘침 0.

## 관찰 (결함 아님, 모두 low)

| ID | 내용 | 재현 | 비고 |
|---|---|---|---|
| O-AP-01 | 목록 ‘요청 검토’ 링크 7개의 접근성 이름이 3종뿐이다. 이름이 ‘요청 검토: 유형 / 요청 에이전트 / 기한’이라 요청자·유형·기한이 같으면 같은 이름이고 **받은 내 에이전트는 이름에 없다**. | `/home/legacy` A, 링크 이름 비교(`browser-*-keys.json` `listReviewNames`) | 실제로는 기한이 초 단위라 구분되는 경우가 많다. 내 에이전트가 여럿인 회원이 같은 유형 요청을 받으면 구분이 약하다. U7 ‘반복 버튼 이름에 대상과 구분값’ 관점의 보강 후보 |
| O-AP-02 | 확인할 수 없음 상태 행(`original_missing`)이 요청 유형 ‘일정 확정’을 그대로 보인다. 상세는 요약을 숨긴다. | 목록 마지막 행 | 소스(`gate_summary.go gateRows`)상 대기가 아닌 행의 유형은 저장된 정책 문자열에서 온다. 권한 확장은 아니다. 일관성 차원의 확인 후보 |
| O-AP-03 | 익명 429 화면 위치 nav에 ‘내 에이전트’(`/home/agents`)·‘요청 상태 조회’ 링크가 남는다. 익명에게는 로그인으로 가는 링크다. | `render/receipt-anonymous-limited.html` | 설계의 ‘실제 로그인 링크’는 충족한다. 로그인 전 화면에서 회원 영역 링크가 보이는 점은 designer 판단 |

시각 판정(거절 버튼이 폭 좁게 보임 등)은 tester 범위 밖이며 designer 몫이다.

## 실패·보정 기록 (원 결과 보존)

- Chrome·Edge `keys` 첫 실행 각 1건 FAIL: 기대 Tab 순서에서 헤더 링크(‘KnowsLink’·‘내 계정’)를 빠뜨린 **내 기대값 오류**였다. 제품 결함 아님. 기대값을 고쳐 다시 실행해 PASS. 첫 결과 [chrome](browser-chrome-keys-first-run-expectation-error.json)·[edge](browser-edge-keys-first-run-expectation-error.json) 보존.
- `mutate.py 2` 첫 실행 2건 FAIL: stage-1 token을 브라우저 거절 세션과 다른 세션으로 재제출해 409를 기대했다. 서버는 403(다른 세션 CSRF)과 같은 안전 문구로 거절했다. **내 기대 오류**. 같은 세션 409는 `mutate.py 1`(decide_1 승인 뒤 중복 승인·거절)에서 확인했고 stage 2는 403 기대로 고쳤다. 첫 결과 [보존](mutate-result-2-first-run-cross-session-expectation.json).
- `onm04` 첫 실행 1건 FAIL: 검사 스크립트의 정규식 이스케이프 오류. 수정 뒤 PASS. [보존](browser-chrome-onm04-first-run-regex-escape-error.json).
- `flow` 실행은 모든 검사를 저장한 뒤 마지막 관찰 `fetch()`가 페이지 CSP(`default-src 'none'`) 때문에 막혀 종료했다(제품 동작, 검사 5건은 PASS). 대기 수는 HTTP로 읽었다. 스크립트에서 그 fetch를 뺐다. flow는 비가역이라 다시 실행하지 않았다.
- 이 흐름 뒤 읽기 전용 `chrome layout`·`chrome keys`를 다시 돌려 `browser-chrome-layout.json`·`browser-chrome-keys.json`이 **결정 뒤 상태**의 결과다(PASS 60·22). 그래서 `chrome-detail-decide-1-*.png`는 승인된 상태, `detail-decide-2`는 캡처하지 않았다. Edge layout·keys와 `http-result.json`은 결정 전 상태다. 대기 상태의 상세 화면은 `detail-query-hint*`·`detail-commit-long*`가 같은 UI다.

## designer용 캡처

Chrome 1440·390·320, 파일명 `chrome-<화면>-<폭>.png`(이 폴더). 화면: `list-c-empty`, `list-a`, `home-a`, `detail-query-hint`(+`-open`), `detail-commit-long`(+`-open`), `detail-decide-1`(승인됨), `detail-expired`, `detail-processed`, `detail-original-missing`, `detail-foreign-b`, `friends-a`(O-NM-08), `gate-409-refused`(중복 클릭 409 본문), `onm04-anonymous-429`, 결정 뒤 `after-list-a`·`after-detail-decide-1-approved`·`after-detail-decide-2-denied`·`after-home-a`, `flow-after-deny-1440`. 캡처 hash는 각 `browser-chrome-*.json`의 `captures`에 있다.

**fixture 상태 변경 알림(designer):** 이 QA로 `pending_decide_1`은 **승인됨**, `pending_decide_2`는 **거절됨**으로 바뀌었고 되돌릴 수 없다. 대기 상태로 남은 것은 `pending_query_hint`·`pending_commit_long`이며 GET만 하면 상태가 바뀌지 않는다. 친구 요청 C→A도 그대로다.

## not_run

- 스크린리더 낭독, 실제 브라우저 200%(Chrome 설정 Page zoom) — 640px viewport reflow로 등가 확인만 했다. Firefox·Safari·모바일 실기.
- 이름 있는 요청자의 gate 화면, 이름 변경 뒤 같은 원요청 식별, 바뀐 digest·관계 철회 화면: fixture가 만들지 않거나 POST가 막혀 있어 DEV integration(`TestApprovalSummaryBoundaries`) 근거에 의존한다.
- 실제 익명 429 응답(fixture가 `/home/receipts`를 막음): 템플릿 렌더로 대체했다.
- `make lint`·`make test`·integration: DEV 같은 SHA(원격 Linux EXIT=0)를 재사용했고 다시 실행하지 않았다.
- U7 일반의 ‘내 계정’→‘설정’ 이름 변경(지시서 범위 밖): 헤더에 아직 ‘내 계정’이 보인다.

## 파일

`checks.py`(HTTP), `mutate.py`(상태 변경 1·2, `stash.json`은 token을 가린 기록), `browser.mjs`(CDP 전 단계), `client.py`, `render/`(O-NM-04·409 HTML), 결과 JSON·캡처. `browser_head.mjs`는 `browser.mjs` 조립 과정의 잔여 파일로 내용을 한 줄 주석으로 바꿨다(셸 가드가 삭제를 막아 남김). 서버는 건드리지 않았고 요청은 fixture origin뿐이다.
