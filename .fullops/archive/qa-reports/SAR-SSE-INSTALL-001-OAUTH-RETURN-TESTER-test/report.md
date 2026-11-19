---
title: SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER 독립 Chromium QA
status: draft
updated: 2026-10-11
owner: tester
tasks: [SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER]
summary: DEV 최종 SHA eac603b의 OAuth 승인 callback 복귀를 실제 Chromium·격리 Postgres·합성 Google/호스트로 독립 검증한 결과와 증거
---

# SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER 독립 Chromium QA

## 대상과 조건

- 대상 제품 SHA: `eac603b6a41d2ce002e698600ac66b8fc67ad762` (DEV 최종). 비교 기준: `c399a5e52c714d4ee61a7cbca03d32ffbfed9bcf`. 기준은 대상의 조상이다.
- 제품 변경은 4개 파일이다: `remote_auth.go`(+6/−2), `remote_browser_integration_test.go`(+113), `remote_integration_test.go`(+17/−2), `testdata/consent_browser.mjs`(+80). `http.go`·`capacity.go`·`google.go` 등 공통 경로는 바뀌지 않았다.
- 리뷰 세션: `d7e2c059-f0af-44ff-82ce-6e0f2f6ea834`. 구현자 세션: `e20cc179-625e-43fa-84fd-b431d432b9c6`. 두 세션은 다르다.
- 읽기 전용 detached snapshot: `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-9f195c88a523468eb4bd4fe91452e8ae` (HEAD = 대상 SHA, `git status` 변경 0). snapshot은 수정하지 않았다.
- 실행 환경: 기존 Linux 서버의 고유 경로 `/tmp/knowslink-oauth001-tester-d7e2c059`, 고유 Postgres 17 컨테이너 `knowslink-tester-oauth001-d7e2c059`(loopback 임시 port·256MiB·1CPU), Google Chrome 148.0.7778.96 headless, Node 22.22.2. 운영 DB·Tunnel·실제 Google/Dots 계정은 사용하지 않았다.
- 테스트 레벨 lite. 범위는 변경한 핵심 동작과 필요한 실패 경계다.
- 삭제·이동·cleanup·trap은 하지 않았다. 임시 컨테이너와 clone은 보존했다.

## 결과 요약

| # | 검사 | 결과 |
|---|---|---|
| 1 | DEV 테스트를 같은 SHA에서 독립 재실행(`TestRemoteConsentBrowserReturn`, `TestRemoteOAuthHTTP`, `TestGoogleHTTP`, `TestGoogleDeviceHTTP`, `TestTextEventsHTTP`) | PASS, 종료코드 0 |
| 2 | 기준 handler(`c399a5e`의 `remote_auth.go`) RED 재현 | FAIL. Chromium이 `form-action 'self' https://accounts.google.com` 위반을 보고하고 callback에 도착하지 않았다 |
| 3 | 수정 SHA에서 같은 시나리오 GREEN | PASS. 승인 POST 200 → 같은 출처 반환 화면 → `https://callback.test/oauth` 자동 도착 → PKCE token 교환 200 |
| 4 | 허용 callback·state·iss·code, 거절, 만료/사용된 consent, CSRF, 로그인 없음, Google 계정(owner) binding, agent 한도 | PASS(아래 표). 발견 F1(medium), F2(low) |
| 5 | token 교환 경계(wrong verifier/redirect_uri/resource/client/scope, 짧은·없는 verifier, 만료 code, code 재사용) | 모두 400 `invalid_grant`. 정상 교환 200 |
| 6 | authorize 경계(미등록 redirect, 다른 resource, plain PKCE, challenge 없음, 모르는 client, 다른 scope, `response_type=token`) | 모두 400 `invalid_request`, attempt 미생성 |
| 7 | DEV exact SHA 등록 `product-lint`/`product-test` | PASS. 이 SHA에서 `lint.py`를 다시 실행했다(ERROR 0, WARNING 2, 실행 불가 0) |
| 8 | 반환 화면·오류 화면 desktop(1280×800)/mobile(390×844) 캡처와 링크 실제 클릭 | 캡처 19장. 링크 클릭으로 callback 도착 확인. 시각 수락은 designer 담당 |

## 독립 RED/GREEN 재현

1. RED: 격리 복사본 `/tmp/knowslink-oauth001-tester-d7e2c059/repro`에서 `internal/relay/remote_auth.go`만 `git show c399a5e:internal/relay/remote_auth.go`로 덮었다(sha256 `850987e6…05bec`). 같은 테스트로 실행했다.
   - 로그 [red-run.log](red-run.log): `Sending form data to 'https://example.com/remote/consent/<id>' violates the following Content Security Policy directive: "form-action 'self' https://accounts.google.com". The request has been blocked.` 승인 뒤 callback 도착 0건이다.
   - 같은 실행에서 consent 재방문은 401이었다. 승인 POST가 서버에서 처리되어 attempt가 소멸했고, 브라우저는 callback으로 가지 못했다. 이는 운영 관찰(install 1, token 0, 1분 code 만료, 재시도 시 `authorization_failed`)과 일치한다. 최초 실패 원인은 만료가 아니라 callback 차단이다.
2. GREEN: 대상 SHA 클론(`clone`)에서 같은 시나리오는 CSP 위반 0건이며 callback에 도착한다. 로그 [qa-run-final.log](qa-run-final.log), DEV 테스트 재실행 [dev-tests-rerun.log](dev-tests-rerun.log).

## 수정 후 시나리오 결과 (대상 SHA, Chromium)

원 code·token은 기록하지 않았다. 표의 "code 있음"은 존재 여부만 뜻한다.

| 시나리오 | 뷰포트 | 관찰 | 판정 |
|---|---|---|---|
| approve-auto (승인 → 자동 이동) | 1280×800 | 승인 POST 200, CSP 위반 0, 자동으로 `callback.test/oauth` 도착, query 키 `code`·`iss`·`state`, state=`qa-state`, iss=`https://example.com`, code 있음, install 정확히 1건 증가 | PASS |
| approve-auto 이후 같은 consent 재방문 | 1280×800 | GET 401, 재시작 안내 있음, `callback.test`·`code=` 미노출 | PASS |
| approve 링크 클릭(자동 이동을 204로 보류) | desktop, mobile | 반환 화면 h1 "앱으로 돌아가는 중", `meta refresh 0;url=https://callback.test/oauth?…`, 링크 rect desktop [260,207,229,24]·mobile [20,232,229,24], 가로 스크롤 없음, 실제 마우스/터치 클릭으로 callback 도착, PKCE 교환 200 | PASS |
| deny (취소) | 1280×800 | callback query 키 `error`·`iss`·`state`, error=`access_denied`, code 없음, install 0건 증가, 재방문 401 안내 | PASS |
| 만료된 attempt(DB의 attempt.Exp를 과거로) | 1280×800 | POST 401, 재시작 안내 화면, callback 도착 0, install 0 | PASS |
| Google 계정 binding 불일치(attempt.Owner를 다른 값으로) | 390×844 | POST 401, 재시작 안내 화면, callback 0, install 0 | PASS |
| 로그인 없음(쿠키 삭제) | 1280×800 | POST 401, 안내 화면, callback 0, install 0 | PASS |
| CSRF 변조 후 정상 재시도 | 390×844 | 첫 POST 403 안내 화면, callback 0. attempt는 유지되어 정상 재시도 POST 200 → callback 1회, install 1건 | PASS |
| agent 한도 초과(owner에 agent 10건) | 1280×800 | POST 409 안내 화면, callback 0, install 0. 안내 문구가 원인과 다르다(F1) | PASS(동작), 발견 F1 |

모든 consent POST 응답(승인 200과 거절 401/403)은 `Cache-Control: no-store`, `Referrer-Policy: no-referrer`, `X-Content-Type-Options: nosniff`, 변경 없는 CSP(`default-src 'none'; style-src 'unsafe-inline'; form-action 'self' https://accounts.google.com; frame-ancestors 'none'; base-uri 'none'`), `text/html`이다. 테스트가 문자열 전체를 비교한다.

## 보안·범위 확인

- 전역 wildcard CSP, callback origin을 CSP에 결합, SameSite·CSRF 완화, 로그인 생략, TTL 연장, attempt 재사용은 없다. diff가 `remote_auth.go`의 consent handler에서 끝난다.
- 반환 화면의 callback 주소는 attempt에 저장된 등록 redirect에 `state`·`iss`·`code`/`error`만 붙인다. html/template이 `<meta content>`와 `<a href>`를 escape한다. 등록 redirect는 `https`만 허용된다(`remoteURL`).
- 새 노출: 승인 반환 화면의 HTML에 code가 있다. 기존 Location header와 같은 값이다. no-store·no-referrer가 유지된다. 자동 이동 흐름에서 반환 화면은 탐색 기록 항목으로 남지 않았다(기록: about:blank, consent, Google callback, consent, callback). 수동 보류 시나리오는 이 기록 검사에 쓰지 않았다.
- 공통 middleware 변경이 없다. 영향 확인으로 `TestGoogleHTTP`·`TestGoogleDeviceHTTP`·`TestTextEventsHTTP`를 같은 SHA에서 재실행했고 모두 PASS다.
- 로그·증거에는 URL query, code, token, 쿠키 값을 기록하지 않았다. consent id는 합성 DB의 소멸한 attempt이며 증거 로그에서도 `<id>`로 지웠다.

## 등록 lint/test 근거

- DEV 게이트: `dev-eac603b6-final-gate.json`의 head는 `eac603b6a41d2ce002e698600ac66b8fc67ad762`, `product-lint`·`product-test` 종료코드 0, ERROR 0/WARNING 3/실행 불가 0이다. 읽고 대조했다.
- 새 실행: 기준 `c399a5e`에서 snapshot 체크아웃으로 `lint.py`를 실행했다. 등록 명령 `make lint`/`make test`는 tester 전용 shim(`qa/make.cmd`, `qa/tester_make.py`)이 snapshot HEAD와 같은 SHA의 깨끗한 Linux 클론(`/tmp/knowslink-oauth001-tester-d7e2c059/lint`)에서 실행했다. shim은 HEAD 불일치·dirty를 거부한다. 결과: product-lint passed, product-test passed, ERROR 0, WARNING 2(SIZE-001 `remote_auth.go` 505줄, SLOP-004 `consent_browser.mjs:13`), 실행 불가 0. 결과 파일은 리뷰 폴더의 `lint.json`이다.
- SIZE-001은 기존 OAuth 경계 파일의 +2줄이다. SLOP-004는 드라이버가 결과 JSON을 stdout에 쓰는 정상 출력이다. 둘 다 수락을 막지 않는다.

## 발견 사항

| ID | 심각도 | 내용 | 상태 |
|---|---|---|---|
| F1 | medium | `remoteConsent`의 모든 실패가 같은 안내("연결 요청이 만료되었거나 이미 처리되었습니다 … 연결을 시작한 앱에서 다시 시작하세요")를 보여 준다. 409 agent 한도 초과는 앱에서 다시 시작해도 해결되지 않고 `/home`에서 agent를 정리해야 한다. 403 CSRF 불일치는 attempt가 살아 있는데 "다시 사용할 수 없다"고 안내한다. 401 로그인 만료·계정 불일치도 같은 문구다. 안전(비밀 비노출, callback 미도달)은 유지된다. 한도 초과는 이전 실패로 남은 install이 agent 한도를 쓰는 운영 상황에서 실제로 보일 수 있다. | 미해결. 후속 DEV 권고: `err` 종류별 문구(capacity는 "agent 한도"와 `/home` 안내) |
| F2 | low | 반환 화면의 링크는 높이 24px 한 줄이다. 별도 카드 스타일 없이 기본 본문으로 표시된다. 자동 이동 전 화면이 매우 짧다. 접근성·시각 판단은 designer 몫이다. | 미해결. designer 확인 |
| F3 | info | 기존 medium(SSE 라우트 보안 헤더 누락)·`verify.py` baseline 결함은 이번 변경과 무관하며 새 결함으로 세지 않았다. 이번 실행에서 다루지 않았다. | 기존 항목 유지 |

미해결 critical/high: 0건.

## 미실행·한계

- 실제 Dots·ChatGPT 호스트의 callback 처리와 token 교환, 실제 Google 계정·운영 DB·운영 도메인은 실행하지 않았다. 성공으로 보고하지 않는다. 배포 뒤 사용자 계정 확인은 부모·OPS 후속이다.
- Chromium(Linux Chrome 148) 한 엔진만 확인했다. Firefox·Safari·모바일 앱 WebView는 확인하지 않았다. DEV 문서의 MDN 설명대로 Firefox는 redirect를 막지 않는다.
- 반환 화면 시각 품질·접근성의 직접 판정은 designer 몫이다. 캡처만 넘긴다.
- 대기 시간 기반 만료(실시간 10분/1분)는 기다리지 않고 DB 상태의 만료 시각을 조작해 검증했다. 서버의 시계 소스 변경은 검증하지 않았다.

## 증거와 재현

- 증거 폴더: 이 디렉터리. `qa/oauth_return_driver.mjs`(CDP 드라이버), `qa/zz_tester_oauth_return_integration_test.go`(Go 통합 테스트, 대상 SHA 클론 옆에 복사해 실행. 제품 SHA에는 넣지 않았다), `qa/tester_remote.py`·`qa/tester_make.py`·`qa/make.cmd`(전용 SSH helper와 shim. 비밀값은 메모리에서만 읽고 출력을 지운다), `qa/outcomes.py`(패킷 처리 기록 생성기).
- 캡처: `screenshots/`(19장). `approve-link-{desktop,mobile}-return.png`가 반환 화면, `*-refusal.png`·`*-revisit.png`가 안내 화면, `*-consent.png`가 동의 화면이다. 캡처에는 code·token이 보이지 않는다.
- 로그: `dev-tests-rerun.log`, `red-run.log`, `qa-run.log`(초기 실행), `qa-run-full.log`(capacity 추가), `qa-run-final.log`(최종, History 기록 포함), `migration.log`, `lint-install.log`, `services-*.txt`.
- 서비스 보존: 실행 전 컨테이너 8개. 실행 뒤 기존 8개의 이름·ID는 `cmp` 동일이며 tester 전용 PG 1개만 늘었다(`services-compare.txt`). 이 PG는 의도적으로 실행 중이다.
- 재현: 서버에서 clone을 `eac603b`에 detached checkout하고 `zz_tester_oauth_return_integration_test.go`를 `internal/relay`에 복사한다. `TEST_DATABASE_URL`(격리 PG)·`TEST_SYNTHETIC_DATABASE=1`·`TEST_BROWSER=/usr/bin/google-chrome`·`QA_DRIVER`·`QA_OUT`를 설정하고 `go test -race -tags 'integration browser' ./internal/relay -run TestQAOAuthReturn -count=1 -v`를 실행한다. RED는 `remote_auth.go`를 `c399a5e` 버전으로 덮은 복사본에서 실행한다.
