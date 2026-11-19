---
title: SAR-PUBLIC-IDENTITY-001-DEV 실행 기록
status: draft
updated: 2026-10-05
owner: dev
tasks: [SAR-PUBLIC-IDENTITY-001-DEV]
summary: 일반 이메일 코드 신원과 세션 및 자기 owner 홈의 기술 선택과 구현 및 검증과 인계를 기록한다
---

# SAR-PUBLIC-IDENTITY-001-DEV — 일반 이메일 신원과 세션 및 자기 owner 화면 실행 기록

## 기준

- 준비 SHA `d7745d584a4cfe43c8a11bb90b037f94b88bdadb`(브랜치 `fullops/dev`). lint 기준 `94533b207b456c0560800fe30a7c90b2b5887c6e`. 제품 고정 `1233e4c3167f722d51f99cb2ef495691734be714`(독립 `SAR-PUBLIC-SERVICE-001-REVIEW` 수락).
- 구현 코드 SHA `a446d89ff288c4243ad6d7f8780a778517154584`. 이 기록·인박스 보고는 뒤의 문서 커밋이다. 문서 커밋은 제품 코드를 바꾸지 않는다.
- Task `task_568f0a5f227c`, Dispatch `ctx_cae2f16a8f1d`, coordinator `term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9`, Run `run_8ca8bc058ab7`.
- 적용 규칙: `fullops-common-0.3.2`, FULLOPS.md, project.md, orca-agents.md, document-writing.md, coding-style/testing/security. Ponytail full. 예외 없음.
- 제품 기준: D02 [일반 이메일 서비스](../../planning/product-specs/SAR-PUBLIC-SERVICE.md) PS-01–04·신원 PS-11, [UX-01–03](../../design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md), frozen [SAR-MVP](../../planning/product-specs/SAR-MVP.md) C1–C5.

## 기술 선택과 근거

### 신원 제공자

relay가 6자리 이메일 코드를 직접 발급·검증한다. 발송은 표준 SMTP submission이다. Go 표준 라이브러리만 사용한다. 새 의존성은 없다.

| 선택지 | 판단 | 근거 |
|---|---|---|
| A. Cloudflare Access One-time PIN + JWT | 채택하지 않음 | OTP는 Cloudflare가 집행한다. PIN 유효 10분은 같지만 오답 5회·재발송 60s·이메일/IP/전체 시간당 한도를 relay가 집행하거나 관측할 수 없다. relay 로그아웃 뒤 Access cookie로 재확인 없이 재진입한다. root를 모든 OTP 이메일에 열어야 한다. 좌석 한도·요금 미확인(OPS `0313deae`). |
| B. relay 이메일 코드 + SMTP | 채택 | D02 기본값 전체를 relay가 집행하고 자동 검사할 수 있다. 사용자에게 Cloudflare 계정을 요구하지 않는다. 제공자 독립이다. |
| C. Cloudflare 계정 IdP | 채택하지 않음 | 일반 사용자에게 Cloudflare 계정을 요구한다. |

운영 발송 후보는 Cloudflare Email Service의 SMTP다. 공식 근거(2026-10-05 조회):

- [SMTP submission changelog 2026-06-08](https://developers.cloudflare.com/changelog/post/2026-06-08-smtp-submission/): `smtp.mx.cloudflare.net:465` implicit TLS, AUTH PLAIN/LOGIN, username `api_token`, password는 **Email Sending: Edit** 권한 API token. beta.
- [Email Service 개요](https://developers.cloudflare.com/email-service/)·[Pricing](https://developers.cloudflare.com/email-service/platform/pricing/): Email Sending은 Workers Paid plan에서만 임의 수신자에게 보낸다. 월 3,000건 포함, 이후 1,000건당 $0.35. 계정의 verified destination address로 보내는 메일은 모든 plan에서 무료다.
- [One-time PIN](https://developers.cloudflare.com/cloudflare-one/integrations/identity-providers/one-time-pin/): Access PIN은 최초 요청 10분 뒤 만료. 선택지 A 비교에만 사용했다.

따라서 일반 공개 발송은 Workers Paid의 신규 비용이 필요할 수 있다. 사용자 본인 이메일 하나의 최종 시험은 verified destination으로 무료 경로가 가능하다. 이 판단과 계정 설정은 OPS·coor·사용자 확인 대상이다. 구현은 SMTP 제공자를 바꿔도 같은 설정으로 동작한다.

### 구조·데이터·세션

- 저장: 기존 `relay_state` singleton JSON에 Member·Identity·Session·Challenge·Rate map을 추가했다. migration은 없다. 같은 global lock이 동시 첫 가입을 직렬화한다.
- 신원 매핑: `knowslink-email-otp|<정규화 이메일>` → 회원 ID `mem_…`. 정규화는 trim과 전체 소문자다. `+tag`·점은 별도 신원이다. 다른 issuer의 같은 이메일은 병합하지 않는다.
- owner 발급: 코드 확인 성공 transaction에서만 회원과 owner를 만든다. 회원 owner는 bearer credential이 없다.
- 세션: random 32 bytes token, DB에는 SHA256. `__Host-kl_session` `Secure; HttpOnly; SameSite=Strict; Path=/`. 절대 12h, 무활동 60분. 확인 성공 시 제시된 이전 세션을 교체한다(fixation 방지).
- CSRF: Go 1.25+ 표준 `http.CrossOriginProtection`을 전체 mux에 적용했다. `Sec-Fetch-Site`/`Origin`으로 cross-origin 브라우저 쓰기를 403으로 거부한다. Node adapter 같은 비브라우저 호출은 두 header가 없어 영향이 없다(verify-mvp 확인). 회원 gate 결정은 기존 HMAC CSRF를 세션 token에 결속한다.
- 신뢰할 인증 근거: 이메일 코드 확인 결과와 서버 저장 세션만 신뢰한다. Access JWT·이메일 header·xAI/OpenAI 계정 이메일 표시는 신원 근거로 읽지 않는다.
- 합성 가입 우회 차단: `/v1/owners`는 `KNOWSLINK_SYNTHETIC_SIGNUP=1`일 때만 동작한다. 기본은 403이다. 회원 세션 cookie는 `/v1` owner·agent credential이 아니다.
- client IP: 기본은 socket 주소다. `KNOWSLINK_CLIENT_IP_HEADER=CF-Connecting-IP`는 relay가 Tunnel로만 닿을 때 OPS가 설정한다. 다른 값은 기동 실패다.

### 운영 기본값의 적용

| D02 기본값 | 구현 |
|---|---|
| 코드 최대 10분·1회·오답 5회 무효·새 요청이 이전 무효화 | `codeTTL`, 성공·5회 오답 시 삭제, 같은 이메일의 기존 Challenge 삭제 |
| 재발송 60s, 1h 이메일 5·IP 20·전체 100, 실패 포함 | `take`(all-or-nothing). 발송 실패도 budget 유지. 같은 안전 안내 |
| 세션 절대 12h·무활동 60분, 전체 로그아웃 5분 재인증 | `session()`, `reauthWindow` |
| 회원 100명 | 확인 성공 시 신규 회원만 거부. 기존 회원 로그인 유지 |
| HTTP 익명 IP 30/60s·인증 principal 40/60s·전체 신규 200/60s·정리 20/100 | `hit`(거부 포함 집계). 시작·확인은 익명, 홈·재확인은 principal, 로그아웃은 정리 budget |
| rolling `(t−window,t]` | `within`과 경계 단위 검사 |

`hit`은 첫 거부 bucket에서 멈추고 bucket마다 최대 limit+1개만 보관한다. flood로 상태 JSON이 무한히 커지지 않는다. 안내 재시도 시각에 요청하면 허용된다(단위 검사).

제공자의 더 강한 제한: Cloudflare Email Service의 발송 limit은 [limits](https://developers.cloudflare.com/email-service/platform/limits/)를 따른다. 이번에 수치를 조회하지 않았다. OPS가 계정 설정 때 확인한다.

### 이번 범위에서 미구현인 항목

- PS-11의 동시 처리 상한(HTTP 16·정리 4·claim 4), agent·관계·queue 한도: 다음 기능(AGENTS/MESSAGES)의 범위다.
- 회원 홈의 gate 화면 요청은 principal rate를 세지 않는다. 기존 gate CSRF·소유 검사는 유지한다.
- 계정 비활성화·30일 정리·agent 연결·키·pair: 화면에 "상태: 준비 중"만 표시한다. 성공 버튼은 없다.

### 다음 agent·원격 MCP/OAuth 신원 경계

- 회원 ID가 사용자별 권한 주체다. 다음 SAR-PUBLIC-AGENTS-001은 세션과 최근 재인증으로 1회·10분 연결 승인을 발급한다. 승인은 회원 owner·대상 agent·클라이언트에 결속한다. 노우와 다닷은 같은 회원 owner 아래 별도 agent·키·credential이다.
- [OpenAI plugin authentication](https://developers.openai.com/plugins/build/auth)은 사용자별 원격 MCP에 OAuth 2.1·PKCE S256·resource metadata·issuer/audience/scope 검증을 요구하고 기존 IdP를 권장한다(coor 전달, 2026-10-05). 후보는 (1) relay가 이메일 코드 로그인을 인증 단계로 쓰는 OAuth authorization server, (2) Access Managed OAuth, (3) 외부 IdP다. 어느 경우든 OAuth subject는 회원 ID에 연결하고 이메일을 token claim의 권한 근거로 쓰지 않는다. 이번에 OAuth를 구현하지 않았다. 선택 전 대상 client의 실제 지원을 확인한다.
- [Grok Bot computer and apps](https://docs.x.ai/grok-bot/computer-and-apps)는 Grok의 실제 computer/browser/terminal 경로다(coor 전달). 노우 연결은 이 경로의 브라우저에서 회원이 연결 승인을 확인하는 흐름을 다음 기능에서 설계한다.

## 구현

| 파일 | 변경 |
|---|---|
| `internal/relay/identity.go` | 정규화·마스킹·`take`/`hit`·코드 발급/검증·회원/owner 바인딩·세션·정리 |
| `internal/relay/member.go` | `/`, `/auth/start`, `/auth/verify`, `/auth/reauth`, `/auth/logout`, `/auth/logout-all`, `/home`, `/home/gates/{id}` |
| `internal/relay/mail.go` | `net/smtp` submission. smtps implicit TLS, smtp STARTTLS. TLS 없는 비loopback 서버 거부 |
| `internal/relay/http.go` | CrossOriginProtection, 합성 가입 설정, Basic/세션 공용 gate 화면 |
| `internal/relay/store.go` | State map·Service 설정 추가, sweep에 신원 정리 |
| `cmd/relay/main.go`, `compose.yaml`, `.env.example` | 설정 4개 전달. 값은 출력하지 않음 |
| `scripts/mail_sink.py`, `Makefile` | 로컬 QA SMTP sink(0600 저장), py_compile 대상 |

## 검증 (HEAD `a446d89`)

명령은 레포 루트에서 실행했다. 종료코드는 명령 자신의 값이다. 로그는 레포 밖 scratch에 두었다.

| 명령 | 결과 |
|---|---|
| `git diff --cached --check` | exit 0 |
| `make lint` | exit 0 |
| `make test` | exit 0 (Go unit/race 신원 단위 검사 포함, adapter test) |
| `make build` | exit 0 |
| `make verify` | exit 0 |
| `make verify-runtime` | exit 0 |
| `make verify-mvp` | exit 0. `TestEmailIdentity` 5개 하위 검사, 기존 PostgresSafety·Gate·Trial·TS synthetic/seed/trial-check 모두 PASS |
| 실제 프로세스 smoke | host `build/relay` + 임시 Postgres 17 + `mail_sink.py` + curl: `/` 200, `/v1/owners` 403, start 303, 60s 내 재요청 429, verify 화면 200(마스킹), 실제 SMTP로 받은 코드로 verify 303, home 200(회원 ID·빈 agent), logout 303, 이전 cookie home 303 expired, cross-site 403. relay·sink 로그에 이메일·코드 0건. 자원은 정리했다 |

`make generate`·`make schema`는 SQL·migration 변경이 없어 실행하지 않았다. `make verify-grok-plugin`은 adapter 변경이 없어 실행하지 않았다.

### 검사가 증명하는 것

- PS-01: 확인 전 owner/회원 미생성, 위조(오답)·만료·재사용·5회 오답 거부, 형식 오류 422.
- PS-02: 재로그인·대소문자·새 브라우저·재시작 후 같은 회원. 동시 첫 가입 두 Challenge → 회원 1개. 다른 issuer 미병합.
- PS-03: 로그아웃 후 이전 cookie 재사용 차단, 전체 로그아웃의 5분 재인증, 무활동·절대 만료, agent credential 유지.
- PS-04: 다른 회원 gate 조회 403·결정 409, 회원 cookie로 `/v1`·`/owner` 401, 합성 가입 403. 회원 gate 결정과 기존 consume 효과 false.
- PS-11 신원 부분: 60s/5/20/100, 익명 30/60s 거부 집계, NAT(다른 source 독립), client IP header 설정 경계, 회원 100 수용량.
- UX-01–03 자동 HTML 근거: 관리자 아이디·SSH 불요 문구, 발송과 확인 구분 문구, 마스킹·회원 ID·준비 중 상태·로그아웃과 agent 철회 구분 문구. 직접 시각 검수는 designer 몫이다.

## 미실행·blocked

- 실제 일반 이메일 확인·로그인·로그아웃·재로그인: **미실행**. 운영 SMTP 설정과 사용자 이메일 입력이 없다. fixture·sink 결과를 실제 PASS로 쓰지 않는다.
- 운영 배포·Cloudflare 설정: 이번 과제 범위 밖이다.
- 독립 코드 리뷰·TESTER 독립 QA·designer 직접 UX 검수: coor 후속.
- 회원 100명 규모 처리량·지연: 미측정. 세션 무활동 갱신은 회원 요청마다 singleton 행을 다시 쓴다.

## OPS 인계 (비밀값 없음)

1. 발송 계정: Cloudflare Email Service Email Sending 활성화와 발신 도메인(`knowslog.com` 후보) SPF/DKIM/DMARC. 임의 수신자 발송은 Workers Paid 필요. 본인 이메일 하나 시험은 verified destination 경로 확인.
2. API token: **Email Sending: Edit** 권한만. 상태 `.env`에 `KNOWSLINK_SMTP_URL=smtps://api_token:<token>@smtp.mx.cloudflare.net:465`, `KNOWSLINK_MAIL_FROM=<발신 주소>`. 값은 Git·로그·보고에 쓰지 않는다.
3. `KNOWSLINK_CLIENT_IP_HEADER=CF-Connecting-IP`는 relay가 Tunnel로만 닿을 때 설정한다. 현재 relay는 `127.0.0.1:8080`에도 게시되므로 호스트 로컬 사용자는 header를 위조할 수 있다. NAT 영향과 함께 검증한다.
4. 합성 가입: 공개 후보 `.env`에는 `KNOWSLINK_SYNTHETIC_SIGNUP`을 두지 않는다. 현재 owner-only 베타의 `beta.sh seed`는 `/v1/owners`를 쓰므로 새 코드 배포 뒤 베타 seed를 유지하려면 그 베타 `.env`에만 `=1`을 둔다. 공개 경로를 연 뒤에는 두지 않는다.
5. 경로 보호: 회원 경로 `/`, `/auth/*`, `/home*`는 relay 인증이 담당한다. `/owner*`와 `/v1/*` owner·agent 경로의 기존 Access 보호를 단순 해제하지 않는다. `deploy/knowslink/verify.py`의 공개 302 기대는 일반 이메일 positive 근거가 아니다. 경로 분리 후 갱신이 필요하다.
6. 가장자리 rate limit·bot 보호는 relay 한도와 별도로 OPS가 결정한다.

## QA·UI 인계

- TESTER: 고정 후보에서 `make verify-mvp`와 README의 sink 절차로 독립 QA. 다른 회원 음성 QA는 독립 fixture 사용 가능. 실제 이메일 확인과 분리해 기록한다.
- designer: README "일반 이메일 로그인 확인" 절차로 UX-01(시작), UX-02(확인 대기·오답·만료·제한·발송 실패), UX-03(홈·로그아웃·전체 로그아웃·재확인)을 직접 확인한다. 코드·이메일은 캡처 전에 가린다. 발송 실패 화면은 SMTP 미설정 relay로 재현한다.
