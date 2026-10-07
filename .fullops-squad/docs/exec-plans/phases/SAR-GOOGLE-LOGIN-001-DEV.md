---
title: Google 로그인 구현 실행 기록
status: draft
updated: 2026-10-07
owner: dev
tasks: [SAR-GOOGLE-LOGIN-001-DEV]
summary: Google 신원과 기존 회원·세션 연결의 좁은 구현 및 검증·운영 설정을 인계한다
---

# SAR-GOOGLE-LOGIN-001-DEV 실행 기록

## 범위와 기술 판단

기준 `abf2de1e3f0501e56b41f9d909514f591b656d4c`, 준비 packet SHA `3b183a5304734d4c50f60d3e26da23e057fd9971`을 사용했다. preamble의 실제 task/dispatch로 회신한다. Google issuer/sub 신원·자동 가입·기존 owner·세션·Google 재확인만 구현한다. 기존 미커밋 지시서와 Jev 자료는 보존한다.

Google ID token 서명 검증을 재발명하지 않으려고 `github.com/coreos/go-oidc/v3` v3.17.0과 `golang.org/x/oauth2` v0.36.0을 추가했다. 하위 의존은 `go-jose/v4` v4.1.3이다. 기존 라이브러리에는 OIDC verifier가 없다. Context7 resolve는 monthly quota 오류였다. 대신 [Google OIDC](https://developers.google.com/identity/openid-connect/openid-connect), [go-oidc API](https://pkg.go.dev/github.com/coreos/go-oidc/v3/oidc), [oauth2 API](https://pkg.go.dev/golang.org/x/oauth2)와 설치한 고정 버전 소스를 확인했다.

기존 JSON state와 세션 발급을 공유한다. 10분 시도·nonce·PKCE·state/cookie·single-use 소비로 code 흐름을 보호한다. 최근 ID token을 다시 확인하는 Google 신원 확인이며 Google 비밀번호 재입력을 보장하지 않는다. 기존 Google 회원의 재확인은 원 세션·회원에 묶인다. 동일 이메일의 이메일 회원은 병합하지 않는다. Strict 세션 cookie를 유지하려고 callback 완료 화면의 홈 링크를 사용한다. 새로운 UI 스타일·DB migration·IdP 추상화·메일 서버는 만들지 않는다.

제품 변경은 새 인증 파일 하나와 기존 세션/UI/설정 연결이다. 안전한 거부·state 소비·실제 서명 검사·Postgres HTTP 회귀 때문에 권장 200–350줄 목표를 넘는 검증 코드가 필요했다. UI 문구 변경 때문에 기존 이메일 integration 기대 문구 1개를 함께 수정했다. 무관한 UTF8·전체 QA·외부 메일 검사는 하지 않았다.

## 검증

- `go test ./internal/relay ./cmd/relay`: exit 0.
- `go test ./internal/relay -run '^TestGoogle' -count=1`: exit 0. RSA/JWK 서명 검증, 위조/거부, 재로그인·owner·세션·state·설정 경계를 확인했다.
- 격리 `postgres:17` 컨테이너 `knowslink-google-dev-b6d97`, loopback 임시 포트, 빈 DB에서 `go run ./cmd/migrate`: exit 0. 합성 자격만 사용했다.
- `TEST_SYNTHETIC_DATABASE=1 TEST_DATABASE_URL=<isolated synthetic DB> go test -race -tags integration ./internal/relay -run '^(TestGoogleHTTP|TestEmailIdentity)$' -count=1`: exit 0. 실제 Postgres와 signed local OAuth provider를 사용했다. Google SMTP 없는 가입·agent 생성·재확인·잘못된 신원·동일 회원 재로그인·로그아웃·state replay·취소·cross-site POST 거부와 기존 이메일 인증을 확인했다. 이 검사는 실제 Google 콘솔 실로그인을 증명하지 않는다.
- `go build ./cmd/...`, `git diff --check`, D10 `deliverables.py --strict`: exit 0.
- required FullOps product-lint/product-test는 완료 아카이브 커밋 뒤 한 번 실행한다. 최종 SHA의 결과 JSON은 체크아웃 Git 디렉터리 `fullops-gate/google-login-lint.json`에 보존하며 worker_done에 종료코드·ERROR·WARNING을 보고한다. 이 문서 작성 시에는 아직 실행하지 않았다.

## 운영 인계와 남은 확인

Web application OAuth client와 `openid email` scope를 사용한다. Client Secret이 필요하다. JS authorized origins는 사용하지 않는다. 운영 Authorized redirect URI와 `KNOWSLINK_GOOGLE_REDIRECT_URL`은 `https://link.knowslog.com/auth/google/callback`이다. 로컬 HTTPS 예시는 `https://localhost:8443/auth/google/callback`이며 별도 콘솔 등록·TLS proxy가 필요하다.

필수 env는 `KNOWSLINK_GOOGLE_CLIENT_ID`, `KNOWSLINK_GOOGLE_CLIENT_SECRET`, `KNOWSLINK_GOOGLE_REDIRECT_URL`이다. secret은 ignored `.env`에만 둔다. SMTP 없이 동작한다. OAuth testing이면 사용자 계정을 Test users에 추가한다. 서버의 Google discovery/token/JWK HTTPS egress가 필요하다. 현재 edge Access 보호와 Tunnel은 수정하지 않았다. callback의 edge 접근과 실제 Google 브라우저 1경로는 coor가 확인한다.

coor는 고정 SHA의 인증 delta 독립 리뷰·필요한 로그인 직접 시각 확인·main 통합/기존 서버 배포를 맡는다. D02 이메일 전용 설명과 D12/D13 환경 안내의 갱신 필요를 coor에 전달했다. DEV는 D10 영향 절과 README/env만 갱신했다. 과거 전체 QA/UTF8 판정은 유지한다.

coor는 2026-10-07 status로 실제 OAuth Web client 생성과 운영 callback 등록·비공개 env 저장 완료를 알렸다. 실제 자격은 DEV에 전달하지 않았고 DEV 검사는 합성 provider를 유지했다. 로컬 Google callback은 실제 client에 아직 등록하지 않았다.
