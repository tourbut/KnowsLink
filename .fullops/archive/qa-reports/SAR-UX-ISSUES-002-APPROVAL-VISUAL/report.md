---
title: 승인 요약 3단위 후보의 직접 시각 판정
status: review
updated: 2026-10-20
owner: designer
tasks: [SAR-UX-ISSUES-002-APPROVAL-VISUAL]
summary: 후보 9d13a05의 승인 목록·상세·익명 429·친구 거절 안내를 1440·390·320 캡처로 판정한다. 결론은 수락이며 DEV 수정 필요 항목은 없다
---

# SAR-UX-ISSUES-002-APPROVAL-VISUAL 직접 시각 판정

**판정은 수락이다.** 승인 목록(대기 있음·빈 목록), 상세 6상태, 익명 429 조회, 친구 거절 안내 배치가 U5·U7 해당 범위와 D04를 충족한다. critical·high·medium은 없다. DEV 수정 필요 항목은 없다. 비차단 관찰 low 2건과 범위 밖 정보 1건을 남긴다.

대상은 고정 후보 `9d13a05ba93c7c085034213590ce661305606010`(fullops/dev, 제품 `f43d05c`, fixture 수정 `ea94128`)이다. 기준 main은 `278001a`, 지시서 커밋은 `0c0fc96`이다. 제품 템플릿은 dev 워크트리에서 `git show 9d13a05:<경로>`로만 읽었다. 제품 코드는 고치지 않았다.

## 증거와 방법

- fixture: origin `http://127.0.0.1:18931`(coor 기동, 컨테이너 `kl-approval002-fixture-9d13a05ba93c-5cd13abf`). gate ID는 manifest `C:/Users/shin/AppData/Local/KnowsLinkDevTools/approval-ux002-dev/fixture-manifest.json`에서 읽었다.
- 요청: 허용 GET(`/home`, `/home/legacy`, `/home/friends`, `/home/gates/*`)과 세션용 `POST /fixture/as/{a|c}`만 보냈다. 승인·거절 POST는 보내지 않았다.
- 도구: [capture.mjs](capture.mjs). headless Chrome(비모바일 viewport, 배율 1)으로 1440·390·320 전체 높이 캡처와 측정값을 만든다. 결과는 [captures/capture.json](captures/capture.json)과 `captures/*.png` 39장이다(전체 33장, 승인 목록 크롭 6장). 실행 결과는 `NO_OVERFLOW 33`, 종료코드 0이다.
- 측정값: 페이지 가로 넘침, 44px 미만 동작, main 안의 가운데·오른쪽 정렬, 그라디언트·그림자·4px 초과 모서리, 사용 색, 제목 순서, details 열림, `aria-current`, alert/status, 링크·버튼 이름.
- 이름 있는 요청자와 익명 429는 DEV 렌더 `C:/Users/shin/orca/workspaces/KnowsLink/dev/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-DEV-test/render/`(아래 `R/`)를 직접 열었다. 익명 429 조회 경로는 fixture 허용 GET이 아니다.
- 운영 기록: 첫 실행은 세션 쿠키(`__Host-`, Secure)를 headless Chrome이 http loopback에서 저장하지 않아 비로그인 GET이 반복됐다. 이 때문에 fixture의 요청 제한(429)이 약 1분 걸렸다. 쿠키를 요청 헤더로 넣고 탐색 간격을 2.5초로 늘려 다시 실행했다. 이 실행의 결과만 판정에 사용한다. 같은 시간대 tester 요청이 429를 받았다면 이 원인이다.

## 화면별 판정

D04 기준은 여섯 토큰과 danger, 시스템 한글 글꼴, 왼쪽 정렬, 44px 동작, 4px 모서리다. 카드 그림자·그라디언트·알약 버튼을 쓰지 않는다. 측정한 모든 화면에서 사용 색은 토큰 값(`#182B31`, `#FFFFFF`, `#005A62`, `#E8F3F3`, `#52646B`, `#A01616`)뿐이다. 그라디언트·그림자·4px 초과 모서리는 0건이다. 33개 캡처 모두 가로 넘침이 0이다. main의 가운데 정렬은 버튼 안의 라벨뿐이고 본문은 왼쪽 정렬이다.

| 화면 | 근거 캡처 | 판정 | 근거 |
|---|---|---|---|
| 승인 목록, 대기 있음(A, 7행) | `captures/legacy-a-{1440,390,320}-gates.png` | PASS | 행마다 요청한 에이전트·받은 내 에이전트·요청 유형·기한·상태·‘요청 검토’가 있다. 원요청 ID는 목록 HTML에 0회다. 링크 접근성 이름은 ‘요청 검토: 유형 / 요청 에이전트 / 기한’이다. 39자 이름은 320에서 줄바꿈된다. |
| 승인 목록, 빈 목록(C) | `captures/legacy-c-empty-{1440,390,320}-gates.png` | PASS | 문구가 U5 빈 목록 문구와 같다. |
| 대화 화면 승인 수 | `captures/home-a-*.png` | PASS | ‘승인 4건’은 대기 4건(query_hint·commit_long·decide_1·decide_2)과 같다. 만료·처리됨·확인 불가를 세지 않는다. |
| 상세 정상(decide_1) | `captures/gate-pending_decide_1-*.png` | PASS | 읽기 순서가 위치 → h1 → 누가 누구에게 → 요청한 내용 → 승인하면 허용되는 범위 → 기한·상태 → 행동 → 기술 상세 → 복귀다. 고정 설명과 버튼 ‘이 요청 승인’·‘거절’이 U5와 같다. |
| 상세 악성 hint(query_hint) | `captures/gate-pending_query_hint-*.png` | PASS | hint는 ‘참고 문구 — 승인 근거 아님’ 접힘 안에 escape된 글자로만 있다. 허용 범위 절에 없다. 링크가 생기지 않는다. |
| 상세 긴 자유문·링크(commit_long) | `captures/gate-pending_commit_long-*.png` | PASS | ‘요청자가 적은 내용’ h3와 설명이 요약 필드와 분리된다. `<a href=…>`는 글자로 보이고 링크가 아니다. 긴 문자열이 320에서 줄바꿈되고 넘침이 없다. |
| 상세 만료 | `captures/gate-expired-*.png` | PASS | 상태 ‘— 만료’, 결정 버튼 없음, status ‘지금은 승인하거나 거절할 수 없습니다…’, 복귀 링크. |
| 상세 처리됨 | `captures/gate-processed-*.png` | PASS | 상태 ‘✓ 승인’, 결정 버튼 없음, 같은 status 안내. 중복 행동을 유도하지 않는다. |
| 상세 원문 부재 | `captures/gate-original_missing-*.png` | PASS | alert ‘승인할 내용을 확인할 수 없습니다…’, 기한·상태, 복귀만 있다. 요약·기술 상세·버튼이 없다. |
| 상세 이름 있는 요청자 | `R/gate-active-{1440,390,320}.png` | PASS | ‘합성 친구 도우미 / node-local / 5e0c7a12’처럼 이름·종류·구분값이 보인다. 순서와 문구는 fixture 화면과 같다. |
| 익명 429 조회(O-NM-04) | `R/receipt-anonymous-limited-{1440,390,320}.png` | PASS(관찰 V-AP-02) | 선택기·조회 폼·하단 링크가 없다. 재시도 시각과 ‘로그인 화면으로 이동’(`/`)만 있다. alert에 autofocus focus 링이 보인다. |
| 친구 거절 안내(O-NM-08) | `captures/friends-a-{1440,390,320}.png` | PASS | 안내가 ‘수락’ 아래 16px, ‘거절’ 위 4px에 있어 ‘거절’과 묶여 읽힌다. `aria-describedby`가 유지된다. |

## U5·U7 기준 대조

- 읽기 순서·고정 문구: 위 표와 같이 설계 순서와 일치한다. 범위 절에 ‘친구 동의나 대화 메시지는 이 요청의 승인 증명이 아닙니다.’가 더 있다. 이 문장은 승인 범위를 넓히지 않으므로 수락한다.
- 기술 상세 접기: 원요청 ID·승인 요청 ID·두 agent ID·intent·정책·검증된 원문 JSON이 닫힌 ‘기술 상세’ 안에 있다. 요약 값(기간·단위)이 원문 JSON과 같다.
- 위치(U7): 상세 title은 ‘요청 승인 · KnowsLink’이고 skip·header·main이 있다. 현재 위치 nav ‘승인 목록 / 요청 승인’의 마지막만 `aria-current=page`다. 주 메뉴를 현재 위치로 표시하지 않는다. 복귀는 `/home/legacy#gates`다.
- 행동 구분: ‘이 요청 승인’(accent 채움)과 ‘거절’(danger 테두리)은 글자로도 구분된다. 390·320에서 승인 버튼만 한 줄 전체 폭이다. D04의 기존 primary 버튼 형태이며 친구 ‘수락’과 같아 수락한다.
- 대비: muted `#52646B` 6.19:1, accent `#005A62` 7.96:1, danger `#A01616` 8.01:1(흰 바탕)로 4.5:1 이상이다. focus는 `:focus-visible` 3px accent 외곽선이다.
- 44px: 모든 링크·버튼·summary가 44px 이상이다. 친구 추가의 동의 체크박스는 24×44px이며 라벨이 같이 눌린다. 이번 변경 밖의 기존 화면이다.

## 관찰과 제품 결정

### V-AP-01 (low, 비차단) — 목록의 확인 불가 행 유형

원문 부재 gate의 목록 행은 유형 ‘일정 확정’과 상태 ‘! 확인할 수 없음’을 보인다. 상세는 유형 없이 ‘승인할 내용을 확인할 수 없습니다’만 보인다. 결정 기록 gate의 유형은 gate 생성 때 서버가 정한 정책 라벨에서 온다(`gate_summary.go:gateRows`). 요청자 자유문이 아니다. 상태가 명시되고 결정 버튼이 없어 권한이 넓어지지 않는다.

결정: 현행을 수락한다. DEV 수정은 요구하지 않는다. 대기 gate가 검증에 실패하면 ‘확인할 수 없는 요청’으로 보이는 현행 규칙도 유지한다.

### V-AP-02 (low, 비차단) — 익명 429의 현재 위치 링크

익명 429 조회 화면의 현재 위치 nav에 ‘내 에이전트’ 링크가 남아 있다. 익명 사용자에게 회원 화면 진입을 보인다. O-NM-04 결정(선택기·빈 href·복원 문구 제거)은 충족했다. 후속 U8·나머지 U7 단위에서 익명 화면의 위치 표시를 함께 정리할 후보로 남긴다.

### V-AP-03 (정보, 범위 밖)

헤더는 아직 ‘내 계정’이다. U7의 ‘설정’ 이름 변경은 이번 3단위 범위 밖이다(DEV 보고 not_run과 같다).

## 미실행(not_run)

- 실제 키보드(Tab/Shift+Tab/Enter/Space), 브라우저 실제 200% 확대, 스크린리더 낭독, 보조 브라우저: designer는 실행하지 않았다. tester 증거가 생기면 별도로 연결한다.
- 승인·거절·중복 클릭 결과 화면: POST 금지에 따라 보지 않았다. tester가 `pending_decide_1·2`로 마지막에 실행한다.
- tester 캡처: 판정 시점에 tester 보고 폴더에는 `client.py`만 있어 참고하지 않았다.

## 남은 일

coor가 tester 독립 QA·고정 SHA 리뷰와 함께 통합을 판단한다. V-AP-02는 후속 단위의 다듬기 후보다.
