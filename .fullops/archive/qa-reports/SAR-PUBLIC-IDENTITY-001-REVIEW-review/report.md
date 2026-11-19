---
title: SAR-PUBLIC-IDENTITY-001 독립 보안 리뷰
status: draft
updated: 2026-10-05
owner: dev
tasks: [SAR-PUBLIC-IDENTITY-001-REVIEW]
summary: 이메일 신원·세션 후보 59b66ad의 독립 리뷰 결과와 rate 포화 medium 발견을 기록한다
---

# SAR-PUBLIC-IDENTITY-001-REVIEW 리뷰

- 검토자 / CLI / 모델: Orca dispatch `ctx_5d0cf2d4375e`(task `task_6cfcf7f4d2f6`, Run `run_8ca8bc058ab7`) / Claude Code / `claude-opus-5-5`. 리뷰 세션 ID는 `0bf1508f-7a43-4fc1-ba68-5f396b0546c4`이다.
- base SHA / head SHA / merge-base: `94533b207b456c0560800fe30a7c90b2b5887c6e` / `59b66ada8b36802484cc6d7e22523257b50572cc` / `94533b207b456c0560800fe30a7c90b2b5887c6e`.
- 준비 커밋: `746ecd9`(preview·rules·result·report 템플릿).
- OCR 버전 / 적용 규칙: open-code-review v1.12.11 (a758d9c) / `review/rule.json`(rule_sha256 `ab2116fb…2aa0`). 공통 규칙 `fullops-common-0.3.2`(`rules/common/README.md`와 coding-style·testing·security). 프로젝트 정본 `project.md`.
- 요구사항·완료 기준 원천: D02 `docs/planning/product-specs/SAR-PUBLIC-SERVICE.md`(PS-01–04, PS-11 신원 부분, 운영 기본값 표), UX `docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md`(UX-01–03), DEV 실행 기록 `docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV.md`, 기록 체크아웃의 `handovers/to_dev.md`.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 31 / 30 / 1 / 31 / 0.
- lint(`lint.json`) ERROR / WARNING / 실행 불가: 0 / 5 / 0. 종료코드 0.

## 독립성 근거

- 구현자: Claude provider 세션 `e0666abc-1cf2-49ac-a288-45a8043404bc`. coor가 session JSONL로 확인한 값이다.
- 이 검토자: Claude Code 세션 `0bf1508f-7a43-4fc1-ba68-5f396b0546c4`. `CLAUDE_CODE_SESSION_ID`와 `~/.claude/projects/-home-shin-orca-workspaces-KnowsLink-fullops-coor/` 최신 JSONL 이름으로 확인했다. 구현자 세션과 다르다.
- snapshot `/tmp/knowslink-public-identity-review-59b66ad`: HEAD `59b66ad…`, detached, `git status --porcelain` 빈 값이다. 리뷰 시작 때와 검사 실행 뒤에 같은 상태를 확인했다. snapshot에서 build·test·install을 실행하지 않았다.
- 재현 검사는 `git archive 59b66ad`로 만든 scratch 사본에서 실행했다. 결과는 이 기록 체크아웃의 리뷰 디렉터리에만 썼다.

## 검토 범위

result.json의 31개 `(path, status)`를 모두 `reviewed`로 기록했다. OCR 제외 파일 `.env.example`(unsupported_ext)도 직접 읽었다.

- 코드(Go 규칙 그룹): `identity.go`, `member.go`, `mail.go`, `http.go`, `store.go`, `cmd/relay/main.go`와 세 테스트 파일. diff와 함께 `principal`, `transaction`, `sweep`, `current`, `cleanup.go`, `ownerPage`를 직접 읽었다.
- 설정: `compose.yaml`, `.env.example`, `Makefile`, `scripts/mail_sink.py`. 변경되지 않은 `deploy/knowslink/beta.sh`의 `.env` 생성과 seed 경로를 영향 범위로 확인했다.
- 문서·JSON: D03–D10 설계 문서, user-guide, project.md, PLANS·board·contexts, DEV 실행 기록, handover logs, Jev route JSON. 문서의 경로·상태코드·cookie·한도·미구현 범위 주장을 코드와 대조했다.

확인한 경계와 결론:

| 항목 | 관측 | 결론 |
|---|---|---|
| 코드 발급·검증 | `crypto/rand` 6자리. 저장값 `SHA256(pending:code)`. `hmac.Equal` 비교. 5회 오답·성공·10분 만료 시 삭제. 실패 횟수는 operation이 nil 오류로 반환해 커밋된다 | 적합 |
| Challenge 식별 | pending 32B token cookie(`__Host-`, Strict). 같은 이메일 새 요청이 이전 Challenge 삭제 | 적합. F2 잔여 위험 |
| SMTP 실패 | lock 밖 발송. 실패 시 Challenge 삭제, budget 유지, 503, 고정 로그 문구 | 적합 |
| 추측·타이밍·열거 | 가입 여부와 무관한 동일 발송 흐름. capacity는 코드 확인 뒤에만 노출 | 적합 |
| 신원 연속성·동시 첫 가입 | `issuer|email` 매핑. 전역 row lock 직렬화. 통합 검사 PASS | 적합 |
| rate·재시작 | Rates를 Postgres JSON에 저장해 재시작 후 유지. `(t−window,t]` | F1 |
| 세션 fixation·CSRF | 확인 성공 시 제시 세션 삭제 후 새 token. `__Host-`·Secure·HttpOnly·Strict. `http.CrossOriginProtection`(Go 1.27.1) 전체 mux. gate는 세션 token HMAC | 적합 |
| 현재·전체 logout | 현재 token 삭제. 전체는 `Verified` 5분 이내만 회원의 모든 세션 삭제. 정리 budget 사용 | 적합 |
| cross-owner·gate | `g.Owner != owner` 검사 유지. 회원 cookie는 `/v1`·`/owner`에서 401 | 적합. F3 |
| 공개 synthetic 우회 | `/v1/owners`는 `SyntheticSignup`일 때만. 회원 owner의 빈 Credential은 `principal`의 빈 token 거부와 해시 비교로 도달 불가 | 적합. F4 운영 위험 |
| 민감정보 | 코드·token 원문 미저장. rate key의 이메일은 SHA256. 로그는 고정 문구·boolean만 | 적합 |
| SMTP 헤더 주입 | 수신자는 `NormalizeEmail`(quoted·이름 거부)과 `net/smtp` Rcpt 검증. 제목 B-encoding | 적합 |
| frozen 업무 회귀 | gate 결정 로직은 ownerAuth 분리 외 동일. 기존 PostgresSafety·Gate·Trial·TS synthetic/seed/trial-check PASS | 적합 |

## 발견 사항

### F1 medium — 단일 source가 전체 신규 rate를 포화시킨다 (미해결)

- 위치: `internal/relay/identity.go:87-112`(`hit`, `anonymousRate`, `memberRate`).
- 원인: `hit`은 bucket을 순서대로 기록한다. 공유 `http:new`(200/60s)를 per-IP(30/60s)·per-member(40/60s)보다 먼저 기록한다. IP 한도로 이미 거부되는 요청도 전체 bucket을 소비한다.
- 영향: 한 익명 source가 60초에 200회(약 3.4 req/s)를 보내면 모든 source의 `/auth/start`·`/auth/verify`가 429가 된다. `memberRate`도 `http:new`를 공유하므로 모든 회원의 `/home`·`/auth/reauth`도 429가 된다. 공격을 계속하면 무기한 유지된다. 로그아웃은 정리 budget이라 영향이 없다.
- 재현(scratch 사본, 리뷰어 작성 검사, PASS = 재현 확인):

```go
func TestReviewSingleSourceSaturatesGlobal(t *testing.T) {
	st := newState()
	now := time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 200; i++ {
		st.hit(now, anonymousRate("198.51.100.1")...)
	}
	if ok, _ := st.hit(now, anonymousRate("203.0.113.9")...); ok {
		t.Fatal("other source allowed")
	}
	if ok, _ := st.hit(now, memberRate("mem_x")...); ok {
		t.Fatal("member home allowed")
	}
}
```

- 기준 관계: D02 HTTP rate 행은 "실패·거부·재시도도 집계한다"고 정한다. principal 단계에서 거부된 요청을 전체 bucket에도 넣어야 하는지는 명시하지 않는다. 현재 구현은 per-IP 한도를 공유 서비스 보호 수단으로 무력화한다.
- 권고: principal bucket을 먼저 평가하고, 그 단계에서 거부된 요청은 전체 bucket에 기록하지 않는다. 또는 designer가 전체 bucket 집계 범위를 확정한다. OPS의 edge rate limit(DEV OPS 인계 6)은 완화책이며 대체 수정은 아니다.

### F2 low — 재발송으로 피해자 코드 무효화와 누적 추측 (잔여 위험)

- 위치: `internal/relay/identity.go:162-178`.
- 타인이 피해자 이메일로 코드를 요청하면 피해자의 진행 중 코드가 무효가 된다. 이메일 5/h 발송 budget도 소모되어 피해자 로그인을 방해할 수 있다.
- 이메일별 누적 오답 상한이 없다. 5회 오답 × 5발송/h로 하루 최대 600회 추측이 가능하다. 하루 탈취 확률은 약 6×10⁻⁴, 1년 지속 시 약 20%다.
- 구현은 D02 기본값과 일치한다. 코드 결함이 아니라 제품 기본값의 잔여 위험이다. designer의 이메일별 일 누적 실패 상한 검토를 권고한다.

### F3 low — 회원 gate 경로에 principal rate가 없다 (다음 단계)

- 위치: `internal/relay/member.go:64-65`.
- `/home/gates/{id}` GET·POST는 신규 40·정리 20 budget을 세지 않는다. DEV 실행 기록이 미구현으로 명시했다.
- 현재 회원 owner는 agent를 연결할 수 없어 gate가 생기지 않는다. SAR-PUBLIC-AGENTS-001에서 approve는 신규, deny는 정리 budget으로 적용해야 한다.

### F4 low — 합성 가입 예시값과 beta seed 회귀 (운영 설정)

- 위치: `.env.example:8`.
- 예시 env가 `KNOWSLINK_SYNTHETIC_SIGNUP=1`을 둔다. 공개 후보 `.env`를 예시에서 복사하면 `/v1/owners`가 신원 확인 없이 owner credential을 발급한다.
- 변경되지 않은 `deploy/knowslink/beta.sh`가 만드는 `.env`에는 이 값이 없다. 새 코드 배포 뒤 beta `seed`의 `/v1/owners`는 403이 된다. DEV OPS 인계 4에 기록되어 있다.
- 권고: OPS 공개 검증에 `/v1/owners` 403 확인을 추가한다.

critical/high 발견 사항은 없다.

## 검증 및 남은 제약

| 검사 | 대상 SHA | 위치 | 종료코드 | 결과 |
|---|---|---|---|---|
| `go test -race -count=1 ./internal/relay/` | `59b66ad` + 리뷰 repro 1개 | scratch 사본 | 0 | 단위 검사와 F1 재현 PASS |
| `npm ci --prefix adapters` | `59b66ad` | scratch 사본 | 0 | — |
| `make verify-mvp` | `59b66ad` + 리뷰 repro 1개 | scratch 사본, 격리 Compose `knowslink-mvp-e442f32134` | 0 | `TestEmailIdentity` 5개 하위 검사, PostgresSafety·Gate·Trial, TS synthetic/seed/trial-check PASS. `down --volumes`로 정리했다 |
| `lint.py --repo fullops-dev --from 94533b2` | `59b66ad`(DEV checkout clean) | `/home/shin/orca/workspaces/KnowsLink/fullops-dev` | 0 | ERROR 0, WARNING 5(SIZE-001), 실행 불가 0 |
| `review.py check --task-key SAR-PUBLIC-IDENTITY-001-DEV` | `94533b2..59b66ad` | 기록 체크아웃 | 0 | 기록 검사 통과: reviewed=31, skipped=0, total=31, lint WARNING=5. jev_find recall 0.714 / precision 0.417(수락 판단에 쓰지 않음) |

lint WARNING 5건은 모두 SIZE-001이다. `PLANS.md` 684줄, `http.go` 544줄, `store.go` 448줄, `member.go` 311줄, `identity_integration_test.go` 303줄이다. 동작 영향은 없다. `http.go`와 `store.go`의 분할은 후속 정리 후보다.

실행하지 않은 검증:

- 실제 운영 SMTP 발송과 일반 사용자 이메일의 로그인·로그아웃·재로그인: 미실행. 운영 SMTP 설정과 사용자 이메일이 없다. sink·fixture 결과를 실제 PASS로 쓰지 않는다.
- `make lint`·`make test` 전체와 `make verify`·`make verify-runtime`: 리뷰어가 별도로 실행하지 않았다. lint 게이트의 product-lint(`make lint`)는 DEV checkout에서 PASS했다. DEV 보고의 실행 결과는 구현자 참고 근거다.
- designer 직접 UX 검수와 TESTER 독립 QA: 별도 과제다.
- 회원 100명 규모 처리량·지연: 미측정(DEV 기록과 동일).
- Cloudflare edge rate limit·Tunnel의 `CF-Connecting-IP` 신뢰 경계: OPS 범위다.

skipped 파일은 없다.

## 검토 결론

기술 리뷰 기준으로 조건부 수락 가능하다. 미해결 critical/high는 없다.

- F1 medium은 공개 수락 전에 수정하거나, designer 판단과 edge 완화로 위험 수용을 기록해야 한다.
- F2–F4 low는 제품 잔여 위험·다음 단계 범위·운영 설정 항목이다.
- 실제 일반 이메일 사람 확인은 미실행이며 PS-13 수락 근거가 아니다.

이 결과는 AI 검토자의 판단이다. OCR의 자동 판정이 아니다. `review.py check` 통과는 기록 일치 검사이며 의미 수락이나 테스트 성공을 자동 보증하지 않는다.
