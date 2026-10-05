---
title: SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX 실행 기록
status: draft
updated: 2026-10-05
owner: dev
tasks: [SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX]
summary: 독립 리뷰 F1 rate 격리와 F4 예시 설정 수정 및 F2 판단과 검증 및 인계를 기록한다
---

# SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX — 로그인 rate 격리·공개 기본값 수정 실행 기록

## 기준

- 준비 SHA `b643e73fb2d740b8845db53c2f2f73f894d87b24`(브랜치 `fullops/dev`). 원본 후보 `59b66ada8b36802484cc6d7e22523257b50572cc`. 독립 리뷰 `25b110fb694d9ccdce6a3b445d6936159b7437d5`. lint 기준 `59b66ada8b36802484cc6d7e22523257b50572cc`.
- 수정 코드 SHA `689ba3f090200c144648cd41a40585902fe2ab8c`. 이 기록·인박스 보고는 뒤의 문서 커밋이다. 문서 커밋은 제품 코드를 바꾸지 않는다.
- Task `task_8c3fd6fcc56b`, Dispatch `ctx_6407d8fcefa7`, coordinator `term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9`, Run `run_8ca8bc058ab7`.
- 적용 규칙: `fullops-common-0.3.2`, FULLOPS.md, project.md, orca-agents.md, document-writing.md, coding-style/testing/security. Ponytail full. 예외 없음.
- 제품 기준: D02 [일반 이메일 서비스](../../planning/product-specs/SAR-PUBLIC-SERVICE.md) PS-01–04·PS-11 운영 기본값, frozen C1–C5. 제품 수치는 바꾸지 않았다.
- 원천: [독립 리뷰](../../evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-REVIEW-review/report.md) F1–F4, [원본 실행 기록](SAR-PUBLIC-IDENTITY-001-DEV.md).

## F1 medium — principal rate 격리 (수정)

### 원인

`hit`은 bucket을 순서대로 기록하고 첫 거부에서 멈춘다. `anonymousRate`·`memberRate`는 공유 `http:new`(200/60s)를 principal bucket보다 앞에 두었다. `cleanupRate`도 공유 `cleanup`(100/60s)을 회원 정리 bucket(20/60s)보다 앞에 두었다. 그래서 자기 principal 한도로 이미 거부되는 요청도 공유 budget을 소비했다. 한 익명 IP가 60초에 200회를 보내면 모든 source의 `/auth/start`·`/auth/verify`와 모든 회원의 `/home`·`/auth/reauth`가 429가 되었다. 정리 경로도 같은 구조였다. 한 회원이 `/auth/logout-all`(5분 재인증 전 403 거부도 정리 budget을 센다)을 100회 보내면 다른 회원의 로그아웃이 429가 되었다.

### 수정

세 함수의 bucket 순서를 principal 먼저로 바꿨다(`internal/relay/identity.go`). `hit`의 동작은 그대로다.

- principal 한도 안의 요청: principal과 공유 bucket을 모두 센다. 공유 bucket이 차 있으면 거부되고 두 bucket 모두에 집계된다. D02의 "실패·거부·재시도도 집계"를 유지한다.
- principal 한도로 거부된 요청: principal bucket에만 집계된다. 공유 budget은 쓰지 않는다.
- 결과: 한 principal은 rolling 60s에 공유 budget에 최대 자기 한도(익명 30·회원 40·정리 20)만 기여한다. 전체 200·정리 100은 독립 principal 7개(익명)·5개(정리) 이상의 합으로만 포화한다.

수치 30/40/20/200/100, `(t−window,t]`, limit+1 보관, 재시도 시각, `State.Rates`의 Postgres 영속은 바꾸지 않았다. 발송 한도 `take`는 all-or-nothing이라 같은 문제가 없다(기존 `TestRollingWindowBoundary`의 partial spend 검사).

### 남은 위험

- 서로 다른 IP 7개 이상을 가진 공격자는 여전히 전체 신규 200을 포화시킬 수 있다. 이것은 D02 전체 한도의 의도된 동작이다. IPv6는 주소를 쉽게 바꿀 수 있다. edge rate limit·bot 보호는 OPS 범위다(원본 기록 OPS 인계 6).
- 전체 신규 포화 중에는 회원 `/home`도 429다. D02가 익명과 회원의 신규 작업에 같은 전체 200을 정했다. 로그아웃은 별도 정리 budget이라 영향이 없다.

## F4 low — 예시 설정의 합성 가입 기본값 (수정)

- `.env.example`의 `KNOWSLINK_SYNTHETIC_SIGNUP`을 빈 값으로 바꿨다. 예시를 복사한 `.env`는 `/v1/owners`를 403으로 닫는다.
- `scripts/verify_mvp.py`는 격리 Compose 환경에 `KNOWSLINK_SYNTHETIC_SIGNUP=1`을 명시해 기존 synthetic/seed/trial-check 검사를 유지한다. 셸 값이 env-file 값보다 우선한다.
- `scripts/check_compose.py`(`make lint`의 lint-config)는 셸 값을 제거한 뒤 relay의 값이 빈 값인지 검사한다.
- README 독립 QA 실행 명령에 `KNOWSLINK_SYNTHETIC_SIGNUP=1`을 추가했다. 값이 없으면 seed가 403으로 실패한다는 안내를 넣었다.
- `deploy/knowslink/beta.sh`는 OPS 소유라 수정하지 않았다. 아래 OPS 인계를 참고한다.

## F2 low — 재발송 무효화·누적 추측 (변경 없음)

- 재발송 무효화: D02 이메일 확인 행이 "새 요청이 이전 확인 수단을 무효화한다"를 정한다. 제품 규칙이므로 유지했다. 타인의 재요청이 피해자 코드를 무효화하고 이메일 5/h 발송 budget을 쓰는 위험은 남는다.
- 누적 추측: 코드당 오답 5회 × 이메일 5발송/h로 하루 최대 600회 추측이 가능하다. 하루 성공 확률은 약 6×10⁻⁴다. 1년 지속 공격에서 약 20%다. 공격 중에는 피해자 메일함에 코드 메일이 반복 도착한다.
- 기술 대안 검토: 코드 길이는 D02가 고정하지 않는다. 하지만 "6자리"는 회원 화면 문구·입력 pattern, D05 메일 계약, user-guide에 있다. 원본 59b66ad의 UI 직접 검수·QA가 병행 중이다. 길이를 바꾸면 그 화면 증거와 문서를 다시 검수해야 한다. 이번에는 바꾸지 않았다. 다음 후보에서 8자리 숫자 코드로 올리면 같은 수치에서 하루 성공 확률이 약 6×10⁻⁶이 된다.
- 이메일별 일 누적 실패 상한은 새 제품 수치다. DEV가 정하지 않는다. designer 검토 대상으로 남긴다.

## F3 low — 회원 gate 경로 rate (다음 단계)

변경하지 않았다. SAR-PUBLIC-AGENTS-001에서 approve는 신규 budget, deny는 정리 budget으로 적용한다.

## 변경 파일

| 파일 | 변경 |
|---|---|
| `internal/relay/identity.go` | `anonymousRate`·`memberRate`·`cleanupRate`의 bucket 순서와 `hit` 주석 |
| `internal/relay/identity_test.go` | `TestRatePrincipalIsolation` |
| `internal/relay/identity_integration_test.go` | `TestEmailIdentity/refused_principal_does_not_spend_shared_rate` |
| `.env.example`, `scripts/verify_mvp.py`, `scripts/check_compose.py` | 합성 가입 기본 닫힘, 격리 검사 opt-in, 기본값 검사 |
| `README.md`, D03 architecture, D05 interface, D10 module-design, `contexts/dev.md` | rate 순서·기본 닫힘·QA opt-in·검사 목록 |

## 검증 (HEAD `689ba3f`)

명령은 레포 루트에서 실행했다. 종료코드는 명령 자신의 값이다. 로그는 레포 밖 scratch에 두었다.

| 명령 | 대상 | 종료코드 | 결과 |
|---|---|---|---|
| `go test -count=1 -run TestRatePrincipalIsolation ./internal/relay/` | 수정 전 코드 + 새 검사 | 1 | RED: `identity_test.go:309: one refused principal blocked others` |
| `make verify-mvp` | 수정 전 코드 + 새 검사 | 2 | RED: `refused_principal_does_not_spend_shared_rate`가 다른 source `/auth/start`에서 `got 429 want 303`. 다른 하위 검사와 기존 검사는 PASS |
| `check_compose.py` | 수정 전 `.env.example` + 새 검사 | 1 | RED: `AssertionError: Synthetic signup must be closed by default` |
| `go test -race -count=1 ./internal/relay/` | 수정 후 | 0 | PASS |
| `python3 scripts/check_compose.py` | 수정 후 | 0 | PASS, `synthetic signup closed` |
| `make lint` | `689ba3f` 작업 트리 | 0 | PASS |
| `make test` | 같음 | 0 | PASS |
| `make build` | 같음 | 0 | PASS |
| `make verify` | 같음 | 0 | PASS |
| `make verify-runtime` | 같음 | 0 | PASS |
| `make verify-mvp` 1회차 | 같음 | 2 | `TestEmailIdentity` 6개·`TestRatePrincipalIsolation` PASS. 기존 `TestTrialHTTP`가 `test_messages_integration_test.go:54: POST /v1/test/persist got 409 invalid_lease` 1회 실패 |
| `make verify-mvp` 2회차 | 같음 | 0 | 전체 PASS: `TestEmailIdentity` 6개, `TestRatePrincipalIsolation`, PostgresSafety·Gate·Trial(`TestTrialHTTP` 포함), TS synthetic/seed/trial-check. `down --volumes` 정리 |
| `git diff --cached --check` | 커밋 전 | 0 | — |

`TestTrialHTTP` 실패는 변경하지 않은 `/v1/test` lease 경로다. 같은 세션의 수정 전 RED 실행과 2회차 실행에서는 PASS했다. 이번 수정과 무관한 기존 간헐 실패로 판단한다. 원인은 조사하지 않았다. 후속 확인 대상으로 남긴다.

`make generate`·`make schema`는 SQL·migration 변경이 없어 실행하지 않았다. `make verify-grok-plugin`은 adapter 변경이 없어 실행하지 않았다.

### 검사가 증명하는 것

- 단위: 한 IP·회원이 1000회 보내도 자기 30/40/20번째까지 허용하고 다음을 거부한다. 다른 IP·회원의 신규·정리는 허용된다. 독립 principal의 합이 정확히 200에서 허용, 201번째에서 거부된다. 정리는 100에서 같다. 60초 뒤 모두 회복한다.
- 통합(실제 Postgres HTTP): 한 IP의 verify 250회와 한 회원의 home 250회·logout-all 150회 뒤에도 자기 요청은 429다. 같은 DB의 새 Service(재시작)에서도 429가 유지된다. 다른 IP의 가입 시작 303, 다른 회원의 home 200·로그아웃 303이 허용된다.
- 재시도 시각·limit+1 보관·all-or-nothing 발송은 기존 `TestRollingWindowBoundary`·`TestSendLimits`가 계속 PASS한다.

## 미실행

- 실제 운영 SMTP 발송과 일반 사용자 이메일의 로그인·로그아웃·재로그인: **미실행**. 운영 SMTP 설정과 사용자 이메일이 없다.
- 원본 전체 QA·UI 직접 검수: 59b66ad에서 병행 중이다. 이 결과로 대체하지 않았다.
- 서버 배포·Cloudflare 쓰기·`beta.sh` 변경: 범위 밖이다.

## QA·UI 인계

- 화면 HTML·문구·상태코드 매핑은 바꾸지 않았다. UX-01–03 화면 증거는 그대로 유효하다.
- 바뀐 동작은 429 발생 조건뿐이다. 한 source·회원의 자기 한도 초과 뒤에도 다른 source·회원은 429가 되지 않는다. 59b66ad에서 "다른 source도 429"를 기대한 QA 시나리오가 있으면 이 후보에서는 실패가 정상이다.
- 수정 후보에서 README 독립 QA 실행(seed)을 할 때는 `KNOWSLINK_SYNTHETIC_SIGNUP=1`을 셸에 준다. 59b66ad의 `.env.example`은 기본 `1`이라 이 차이가 있다. 값 없이 실행하면 seed는 403으로 실패한다.
- 후속 delta 리뷰·QA 범위: `identity.go` bucket 순서, 새 두 검사, `.env.example`·두 검사 스크립트, README 실행 명령.

## OPS 인계 (비밀값 없음)

1. `beta.sh`가 만드는 `.env`에는 `KNOWSLINK_SYNTHETIC_SIGNUP`이 없다. 새 코드 배포 뒤 owner-only 베타의 `beta.sh seed`(`/v1/owners`)는 403이 된다. 베타 seed를 유지하려면 그 베타 `.env`에만 `KNOWSLINK_SYNTHETIC_SIGNUP=1`을 둔다. 공개 경로를 연 뒤에는 두지 않는다.
2. 공개 검증에 `/v1/owners` POST가 403인지 확인을 추가한다.
3. `.env.example`을 복사해 만든 `.env`는 이제 합성 가입이 닫힌다. 예시 값으로 seed를 실행하는 운영 절차가 있으면 위 1을 적용한다.
4. relay 한도는 단일 source 고갈을 막는다. 독립 IP 다수의 합 포화와 IPv6 주소 회전은 edge rate limit·bot 보호로 완화한다.
