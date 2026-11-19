---
title: Remote OAuth 승인 뒤 브라우저 callback 복귀 수정
status: draft
updated: 2026-10-11
owner: dev
tasks: [SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV]
summary: Chromium form-action이 consent POST의 303 callback을 막는 원인을 실제 브라우저로 재현하고 동일 출처 반환 화면으로 최소 수정한다.
---

# SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV

## 적용 기준

기준 SHA는 `ae5d9f0db878216ee011610430615067488fbe4a`다. attempt는 `0a90d721ecf04034a0bdda333f24c7a9`다. 적용 규약은 `fullops-common-0.3.3`, project.md, delegation.md, document-writing.md, Ponytail full이다. 테스트 레벨은 lite다. 하위 위임 레벨은 standard이지만 하위 에이전트를 배정하지 않았다. 원인과 수정이 한 handler에 있어 직접 처리했다.

## 기술 계획

1. 실제 Chromium으로 합성 consent 승인 POST부터 외부 callback 도착까지 재현한다.
2. 승인·거절 결과를 form 제출의 redirect가 아닌 새 동일 출처 문서의 이동으로 바꾼다.
3. 만료·사용된 consent는 JSON 대신 비밀 없는 HTML 안내를 표시한다.
4. 공통 CSP·CSRF·Google 재확인·TTL·client/redirect/resource/issuer binding은 바꾸지 않는다.

## 원인

consent 화면의 공통 CSP는 `form-action 'self' https://accounts.google.com`이다. 승인 POST는 등록된 호스트 callback으로 303 redirect를 반환했다. Chromium은 form 제출에 이어지는 redirect에도 form-action을 적용한다. 따라서 callback 이동이 차단된다. [MDN form-action](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/form-action)의 Warning은 "Firefox 57 doesn't block the redirects whereas Chrome 63 does"라고 명시한다.

서버는 승인 transaction에서 install과 1분 code를 만들고 attempt를 삭제한 뒤 commit한다. 브라우저가 callback으로 가지 못하면 code가 1분 뒤 만료된다. 사용자가 다시 시도하면 삭제된 attempt 때문에 `authorization_failed`를 본다. 운영 관찰(install 1, token 0, code 0, pending attempt 0)과 일치한다. 운영 DB·실제 계정은 이 과제에서 조회·변경하지 않았다.

Chromium 콘솔의 실제 위반 문구는 다음과 같다. 원래 form 주소를 보고하며 callback query는 포함하지 않는다.

`Sending form data to 'https://example.com/remote/consent/<id>' violates the following Content Security Policy directive: "form-action 'self' https://accounts.google.com". The request has been blocked.`

## 수정

- `internal/relay/remote_auth.go`: 승인·거절 POST 결과를 HTTP 200 반환 화면으로 응답한다. 화면은 `<meta http-equiv="refresh">`와 같은 주소의 링크로 등록 callback에 이동한다. 이 이동은 form 제출이 아니므로 form-action 대상이 아니다. callback 주소는 기존처럼 attempt의 정확한 등록 redirect에 `state`·`iss`·`code` 또는 `error`만 추가한다. 응답은 기존 middleware의 `Cache-Control: no-store`, `Referrer-Policy: no-referrer`, `frame-ancestors 'none'`을 유지한다. html/template이 주소를 escape한다.
- 같은 handler의 실패 응답을 기존 `refusal` 화면으로 바꿨다. 상태 코드는 기존 `statusFor` 값(만료·사용 401, CSRF 403)을 유지한다. 안내는 "연결을 시작한 앱에서 KnowsLink 연결을 다시 시작하세요"다. callback·code·state를 표시하지 않는다.
- 전역 form-action 확장, callback origin을 CSP에 결합, SameSite·CSRF 완화, 로그인 생략, TTL 연장, attempt 재사용은 하지 않았다. 공통 security middleware는 바꾸지 않았다.

UI: 반환 화면과 실패 안내는 기존 `memberStyle`과 `refusal` 화면을 재사용한다. 새 CSS·토큰은 없다. 프로젝트의 DESIGN lint는 Go 문자열 CSS를 검사하지 않으므로 해당 없음이다. 반환 화면은 자동 이동 전 잠깐 표시된다. 이 과제에서 별도 캡처 검수는 하지 않았고 tester·designer 후속에서 확인한다.

## 검증

Linux 서버의 기존 격리 경로를 사용했다. 고유 `/tmp/knowslink-oauth001-ctx1e648875` clone, 고정 digest Postgres17 임시 container(loopback 임시 port·256MiB·1CPU)를 사용하고 회수했다. 실행 전후 기존 container 이름·ID 목록이 같았다(`cmp` 종료코드 0). local Docker·운영 DB·Tunnel·실제 Google·Dots 계정은 사용하지 않았다.

| 검사 | 대상 | 결과 |
|---|---|---|
| `TestRemoteConsentBrowserReturn` 재현 | `4960186`의 test와 기준 `remote_auth.go`를 합친 로컬 전용 commit `c9ecb77` | exit1. Chrome 148이 위 CSP 위반을 보고했다. callback 미도착 |
| `TestRemoteConsentBrowserReturn` 수정 | `4960186a934a1da6fbf04b0c51391447f1bfed79` | PASS. 합성 Google 재확인→승인→`https://callback.test/oauth` 도착, CSP 위반 0, state·iss·code 확인, PKCE token 교환 200, 사용한 consent 재방문 401 HTML 안내와 callback 미노출 |
| 회귀 | 같은 SHA, `go test -race -tags 'integration browser'` | `TestRemoteOAuthHTTP`, `TestGoogleHTTP`, `TestGoogleDeviceHTTP`, `TestTextEventsHTTP` PASS(종료코드 0) |
| Windows unit | `go test ./internal/relay`, `go vet` (integration tag 포함) | 종료코드 0 |

브라우저 검사는 `internal/relay/remote_browser_integration_test.go`에 있다. `go test -tags "integration browser"`와 `TEST_BROWSER`(Chromium 경로)로 실행한다. build tag로 분리하여 skip 표식을 쓰지 않는다. `TEST_BROWSER`가 없으면 실패한다. 드라이버 `internal/relay/testdata/consent_browser.mjs`는 Node 22 내장 WebSocket으로 CDP를 사용한다. 새 의존성은 없다. 하나의 TLS 테스트 서버가 Host로 relay·합성 Google·등록 호스트를 나눈다. Chromium `--host-resolver-rules`가 세 이름을 이 서버로 보낸다. 드라이버는 URL query를 출력하지 않는다. 기존 `TestRemoteOAuthHTTP`는 303 대신 반환 화면의 링크에서 callback을 읽도록 바꿨다.

### 실패와 수정

1. Linux Chrome이 SSH 세션에서 첫 network navigation 때 OS keyring을 기다려 응답하지 않았다. `--password-store=basic`을 드라이버에 추가했다. 제품 결함이 아니다.
2. Chromium이 callback host의 `/favicon.ico`도 요청하여 기록한 query를 덮었다. test가 `/oauth` 경로만 기록하도록 고쳤다.
3. 원인 확인을 위해 DB 없는 scratch mock(동일 CSP·흐름)도 Windows Chrome과 Linux Chrome에서 실행했다. 303은 차단되고 반환 화면은 callback에 도착했다. mock은 레포에 넣지 않았다.

실행 출력은 [browser-runs.log](../../evaluations/dev-checks/SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV/browser-runs.log)에 있다. 이전 SHA `9e9fccf`에서도 같은 결과였다. 첫 등록 gate(`98c53e1`)는 product-lint·product-test 통과, ERROR ANTI-004(`t.Skip`)였다. 검사를 완화하지 않고 test를 build tag 파일로 옮겼다. WARNING SIZE-001(remote_auth.go 505줄, +2줄)은 기존 OAuth 경계 파일의 최소 수정이라 분할하지 않았다. PLANS.md 크기 경고는 coor 문서다. SLOP-004는 드라이버의 결과 JSON 출력 채널이며 디버그 잔재가 아니다.

등록 lint/test의 최종 SHA 결과는 완료 보고와 Git 디렉터리의 `fullops-gate`가 정본이다.

## 한계와 후속

- 실제 Dots·ChatGPT 호스트의 callback 처리와 token 교환은 이 과제에서 실행하지 않았다. 배포 뒤 사용자 계정으로 확인해야 한다.
- 승인 뒤 반환 화면은 code를 HTML에 포함한다. 기존 Location header와 같은 값이며 no-store·no-referrer로 보호한다. 자동 이동이 안 되면 사용자가 1분 안에 링크를 눌러야 한다.
- 이미 생성된 운영 install 1건(token 0)은 정리하지 않았다. 30일 install 수명과 기존 sweep 규칙을 따른다. 이 install의 새 agent는 회원 agent 한도에 포함될 수 있다. 운영 데이터 판단은 OPS·부모 후속이다.
- 독립 고정 SHA 리뷰, tester QA, main 통합, OPS 배포, 실제 계정 연결은 부모가 조율한다.
