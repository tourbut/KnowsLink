---
title: SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER — OAuth 승인 callback 복귀 Chromium 시나리오
status: draft
updated: 2026-10-11
owner: tester
tasks: [SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER]
summary: 합성 Google·등록 호스트·격리 Postgres에서 실제 Chromium으로 승인 POST부터 callback 도착과 token 교환까지 확인하는 시나리오와 통과 조건
---

# SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER — Chromium 시나리오

대상은 DEV 최종 SHA `eac603b6a41d2ce002e698600ac66b8fc67ad762`다. Jev 웹 조작 도구 대신 Chrome DevTools Protocol 드라이버를 쓴다. 이유: 합성 호스트 해석(`--host-resolver-rules`), callback 이동 보류(Fetch 가로채기), 응답 헤더·CSP 위반 수집이 필요하다. 판정은 Go 통합 테스트의 단정(assert)이다. 결과는 [QA 보고서](../qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-test/report.md)에 있다.

## 공통 조건

1. 실행 위치는 레포 밖 임시 clone이다. 대상 SHA에 detached checkout하고 `git status --porcelain`이 비어 있어야 한다.
2. 격리 Postgres 17만 쓴다(`TEST_SYNTHETIC_DATABASE=1`). 운영 DB·실제 Google·Dots 계정은 쓰지 않는다.
3. 하나의 TLS 서버가 Host로 relay(`example.com`)·합성 Google(`accounts.google.com`)·등록 호스트(`callback.test`)를 나눈다.
4. 출력에 URL query·code·token·쿠키 값을 쓰지 않는다. code는 존재 여부만 기록한다.
5. 소스는 `qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-test/qa/`에 있다. 테스트 파일은 제품 SHA에 넣지 않고 clone 옆에 복사해 실행한다.

## 시나리오와 통과 조건

| 이름 | 입력 | 통과 조건 |
|---|---|---|
| approve-auto-desktop | 1280×800, 승인, 자동 이동 | CSP 위반 0. callback 1회 도착. state·iss 일치. code 있음. install +1. 응답 헤더(no-store·no-referrer·nosniff·CSP 전체) 동일. 같은 consent 재방문 401과 재시작 안내 |
| token 경계 | 위 code로 교환 | wrong verifier·redirect_uri·resource·client·scope·짧은/없는 verifier는 400. 정상 200. 재사용 400 |
| expired-code | code 만료 시각을 과거로 | 교환 400 |
| approve-link-desktop / -mobile | 자동 이동을 204로 보류 | 반환 화면 meta refresh 0초와 링크 표시, 가로 스크롤 없음, 실제 클릭·터치로 callback 도착, 교환 200 |
| deny-desktop | 취소 | callback에 error=access_denied·state·iss, code 없음, install 0 |
| expired-attempt-desktop | attempt 만료 | POST 401 안내 화면, 재시작 문구, callback 0, install 0 |
| owner-binding-mobile | attempt owner를 다른 값으로 | POST 401 안내 화면, callback 0, install 0 |
| no-login-desktop | 쿠키 삭제 | POST 401 안내 화면, callback 0, install 0 |
| csrf-tamper-mobile | csrf 값 변조 후 정상 재시도 | 첫 POST 403 안내, callback 0, 재시도 200으로 callback 1회, install +1 |
| capacity-desktop | owner agent 10건 | POST 409 HTML 안내, callback 0, install 0 (문구 원인 불일치는 발견 F1) |
| authorize 경계(HTTP) | 미등록 redirect·다른 resource·plain PKCE·challenge 없음·모르는 client·다른 scope·`response_type=token` | 모두 400, attempt 미생성. 유효 요청은 303 |

## RED 기준

기준 `c399a5e`의 `remote_auth.go`로 덮은 복사본에서 approve-auto-desktop은 `form-action` CSP 위반으로 실패해야 하고 callback에 도착하지 않아야 한다. 이 조건이 유지되어야 회귀 시나리오로 유효하다.

## 담당 구분

- 이 시나리오는 tester의 독립 자동 검사다. 반환·안내 화면의 시각 수락은 designer 담당이며 캡처를 넘긴다.
- 실제 Dots·ChatGPT 호스트와 실계정은 범위 밖이다. 배포 뒤 부모·OPS가 확인한다.
