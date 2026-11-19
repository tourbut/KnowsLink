---
title: SAR-PUBLIC-IDENTITY-001-RATE-REVIEW 독립 delta 리뷰
status: draft
updated: 2026-10-06
owner: ops
tasks: [SAR-PUBLIC-IDENTITY-001-RATE-REVIEW]
summary: 고정 후보 9c915dc의 RATE-FIX delta와 원본 TESTER·UI 기록의 독립 리뷰 결과를 기록한다
---

# SAR-PUBLIC-IDENTITY-001-RATE-REVIEW 리뷰

- 검토자 / CLI / 모델: OPS Orca dispatch `ctx_dda6a6213308`(task `task_85f7c3846cae`, Run `run_8ca8bc058ab7`) / Claude Code / `claude-opus-5-5` high. 리뷰 세션 ID는 `886fc5cf-5392-4129-a2e6-f8bfdbf36950`이다.
- base SHA / head SHA / merge-base: `59b66ada8b36802484cc6d7e22523257b50572cc` / `9c915dc71e2a872243ffec294126d4668b4d32a4` / `59b66ada8b36802484cc6d7e22523257b50572cc`.
- 준비 커밋: `62abd1c`(지시서). prepare는 이 리뷰 세션이 실행했다.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11 (a758d9c) / `review/rule.json`(rule_sha256 `ab2116fb…2aa0`). 공통 규칙은 후보의 `fullops-common-0.3.3`(`rules/common/README.md`와 coding-style·testing·security)이다. 프로젝트 정본은 `project.md`다. 기준 완화와 예외는 없다.
- 요구사항·완료 기준 원천: 기록 checkout의 `handovers/to_ops.md`. D02 `docs/planning/product-specs/SAR-PUBLIC-SERVICE.md`(PS-01–04, PS-11 운영 기본값). 원본 리뷰 `25b110f`의 F1–F4. RATE-FIX 지시서 전문(`handovers/logs/2026-10-05_to_dev.md`)과 실행 기록.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 117 / 64 / 53 / 64 / 53.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: 0 / 5 / 0. lint.py 종료코드 0.
- SIZE-002: 지시서의 예상 변경 규모 / 실제 추가 줄 수 / 차이·분할하지 않은 이유: 아래 `lint와 규모` 절.
- DEP-001: 파일별 의존성 변경 여부 / 완료 보고의 필요성·표준 라이브러리 대안: 발생하지 않았다. `go.mod`·`go.sum`·`adapters/package*.json` 변경은 0줄이다.
- UI 디자인: 아래 `UI 영향` 절. 제품 UI 변경이 없어 디자인 lint·테마 전환은 해당 없음이다.

## 독립성 근거

- 원본 구현자: Claude 세션 `e0666abc-1cf2-49ac-a288-45a8043404bc`(59b66ad). 원본 리뷰 report와 PLANS의 DEV 완료 기록에서 확인했다.
- RATE-FIX 구현자: Claude 세션 `0695c42b-0a85-4cc7-90dc-f227379db7b2`(term `term_b516f045…`, 커밋 `689ba3f`·`8c7ea94`). `~/.claude/projects/-home-shin-orca-workspaces-KnowsLink-fullops-dev/`에서 `task_8c3fd6fcc56b`와 두 커밋 명령을 가진 JSONL로 확인했다. result.json의 `implementer_session`은 delta 구현자인 이 값이다.
- 원본 리뷰어: `0bf1508f-7a43-4fc1-ba68-5f396b0546c4`.
- 이 검토자: `886fc5cf-5392-4129-a2e6-f8bfdbf36950`. scratchpad 경로와 `~/.claude/projects/-home-shin-orca-workspaces-KnowsLink-fullops-ops/`의 JSONL 이름으로 확인했다. 위 세 세션과 모두 다르다.
- snapshot `/tmp/knowslink-identity-rate-review-9c915dc`: HEAD `9c915dc…`, detached, `git status --porcelain` 빈 값이다. 시작과 check 실행 때 같은 상태를 확인했다. snapshot에서 build·test·install을 실행하지 않았다.
- lint·test는 scratchpad의 별도 clone(`checkout --detach 9c915dc`)에서 실행했다. 실행 전후 HEAD 같음과 clean을 확인했다. 결과는 이 기록 checkout의 리뷰 디렉터리에만 썼다.

## 검토 범위

result.json의 117개 `(path, status)`를 모두 기록했다. reviewed 64, skipped 53이다. 커버리지는 100%다.

- 제품 코드: `internal/relay/identity.go`의 `hit`과 세 rate 함수. `take`, `within`, `sweepIdentity`의 1시간 prune, `member.go`의 다섯 호출부(117·172·223·272·305행)를 직접 읽었다. 세 테스트 변경을 검토했다.
- 설정·스크립트: `.env.example`(OCR user_exclude, 직접 읽음), `scripts/check_compose.py`, `scripts/verify_mvp.py`, `README.md`.
- 기술 문서: D03 architecture·D05 interface·D10 module-design의 rate·합성 가입 서술을 실제 코드와 대조했다. contexts·PLANS·board의 상태 서술에서 main 수락 주장을 확인했다.
- 원본 TESTER(`9e2654d`)·UI(`cf0ab09`) 기록: 보고서·실행 기록·시나리오·증거 요약을 대조했다. 아래 `원본 TESTER·UI 기록 검토` 절.
- FullOps 0.9.14 운영 변경: FULLOPS.md·project.md·rules/common·lint/README·lint.json·_TEMPLATE·_REPORT·board·실행 기록·세 lint 증거 JSON.
- Jev route/find/context JSON 13개: JSON 유효성과 비밀값 부재만 확인했다. 원본을 보존했다. 수락 판단에 쓰지 않는다.

skipped 53개는 모두 원본 UI 증거 폴더의 이미지·원시 텍스트·로그다(rule.json exclude). 확인 방법은 다음과 같다.

- PNG 18개: 시그니처와 크기(390·1280 폭)를 전수 확인했다. 보고서가 참조한 PNG의 누락은 없다. `validation.json`의 `json_and_png_format: PASS`와 대조했다. `pw-04-resend-limited-mobile.png`와 `pw-11-expired-member-gate-desktop.png`는 직접 열어 reviewed로 기록했다.
- 원시 텍스트·jsonl·로그: 전 파일에서 이메일 주소와 6자리 코드 패턴이 없음을 검사했다. `visual.exit`=0, `visual.log`의 완료 문구, `cleanup.json`을 보고서와 대조했다.
- 생략 영향: 원본 59b66ad의 개별 화면 재판정을 하지 않았다. 이 리뷰는 기록의 무결성과 범위를 검토하며 새 시각 검수가 아니다.

## RATE-FIX 코드 검토

### `hit`의 새 동작

| 경우 | 기록 | 반환 |
|---|---|---|
| 모든 bucket 허용 | 모든 bucket에 `now` 추가 | 허용 |
| principal(첫 bucket)이 가득 | principal만 `limit+1` 보관으로 기록 | 거부, `full[1]+window` |
| principal 허용, 공유가 가득, principal에 창 안 기록 있음 | principal에 추가, 공유는 `limit+1` 보관 | 거부, `full[1]+window` |
| principal 허용, 공유가 가득, principal 기록 없음 | 기록 없음 | 거부, 공유 `kept[len-limit]+window` |

확인한 성질:

- 한 principal이 공유 bucket에 넣는 기록은 rolling 창마다 자기 한도(30·40·20) 이하다. principal bucket이 한도에 이르면 첫 단계에서 거부되고 공유 bucket을 건드리지 않는다.
- 새 rate key는 허용 요청으로만 생긴다. 원래 shared-first 구현과 같은 상한이다. sweep은 1시간 밖 key를 지운다.
- 기록하지 않는 거부의 재시도 시각은 공유 창 안 기록이 한도 아래로 내려가는 시각이다. 기록하는 거부는 기존 의미(다음 요청도 집계)를 유지한다.
- slice 별칭: `kept[i]`는 새로 만든 slice이므로 `st.Rates` 원본과 공유하지 않는다. `for j := range i`는 Go 1.22 이상 문법이며 go.mod는 1.27.1이다.
- `take`(발송 한도)는 all-or-nothing이라 순서 변경의 영향이 없다. `TestRollingWindowBoundary`·`TestSendLimits` PASS로 확인했다.

### 리뷰어 추가 검사

git archive 9c915dc 사본에 리뷰어 전용 `TestReviewHitInvariants`를 넣었다. 후보에는 포함되지 않는다. 200 seed × 3000 요청에서 IP 12개·회원 8개의 익명·회원·정리 요청을 무작위로 섞고 시각을 0–400ms씩 진행했다. 다음 불변식을 매 요청 확인했다.

- 모든 bucket 길이 ≤ `limit+1`.
- principal 한도로 거부된 요청은 공유 bucket 길이를 바꾸지 않는다.
- 창 안 기록이 없던 principal의 거부는 새 key를 만들지 않는다.
- 허용 요청 수는 모든 bucket의 rolling 창에서 한도 이하다.
- `ok`와 `retry.IsZero()`가 일치하고 거부의 retry는 `now`보다 뒤다.

`go test -race -count=1 -run 'TestReviewHitInvariants|TestRate' ./internal/relay/` 종료코드는 0이다.

### 남은 위험(결함 아님)

- 서로 다른 IP 7개 이상(정리는 회원 5명 이상)의 합은 전체 200·정리 100을 포화시킨다. D02 전체 한도의 의도된 동작이다.
- 창 안 기록이 있는 기존 principal은 공유 포화 중에도 자기 한도까지 공유 bucket을 갱신해 포화를 연장할 수 있다. 연장량은 principal 수 × 자기 한도로 제한된다. DEV 실행 기록의 서술과 일치한다. IPv6 회전과 함께 edge rate limit·bot 보호(OPS)로 완화한다.

## 원본 TESTER·UI 기록 검토

| 항목 | 확인 | 결론 |
|---|---|---|
| TESTER 대상 | `59b66ad`, QA checkout `/tmp/knowslink-public-identity-qa-59b66ad`, 준비 `746ecd9` | 원본 후보 결과다. 9c915dc의 QA로 쓰지 않는다 |
| TESTER 집계 | `result.json` 159건: ok 157, skipped 2(`p06-public-boundary`, `p07-human-email-and-designer-ui`), 실패 0 | 보고서와 일치 |
| TESTER F1 | `f1-result.json`: A 30×422, 170×429, `http:new` 200, B 첫 요청 429, `/home` 429, logout 303 | 원본 리뷰 F1과 독립 관측이 일치. 9c915dc에서는 이 관측이 바뀌는 것이 정상이다 |
| TESTER 보존 | to_tester 전문이 logs에 보존됐다. P06·P07은 미체크와 BLOCKED를 유지한다. 증거에 이메일 원문·DB URL이 없다 | 실패·미실행 보존 적합 |
| UI 대상 | `59b66ad`, UI checkout `/tmp/knowslink-public-identity-ui-59b66ad`, 준비 `1207bdf` | 원본 후보 결과다 |
| UI 판정 | 로컬 fixture UX-01–03 직접 시각 PASS, F-UI-01 low, Orca blank 장애와 승인된 별도 Chromium 보완 구분 | 근거·범위·의존성 서술이 증거와 일치 |
| UI 보존 | 실제 이메일·운영 SMTP·공개·P06·P07·노우↔다닷 미실행을 PASS로 쓰지 않는다. `cleanup.json`은 자기 자원 정리를 기록한다 | 미실행 보존 적합 |

## UI 영향

- 59b66ad→9c915dc의 Go 변경은 `identity.go`의 rate 코드와 테스트뿐이다. `member.go`·`http.go`의 템플릿·문구·상태코드 매핑은 바뀌지 않았다. 바뀐 것은 429의 발생 조건이다.
- 원본 UI PASS는 59b66ad의 직접 시각 결과다. 이 리뷰는 그 기록을 9c915dc의 새 직접 시각 검수로 표시하지 않는다. 화면 변경이 없으므로 UI 재검수 필요성은 낮다. 429 안내 화면(`pw-04`)의 문구와 경로는 같다.
- 테마·디자인 lint: project.md가 밝힌 대로 서비스의 디자인 lint와 공용 테마 전환은 미구성이다. DESIGN 기본 규칙은 Go 문자열 안의 CSS를 검사하지 않는다. 이 delta에서 DESIGN 경고는 0건이다. 테마 전환 검증은 해당 없음이다.
- README QA 실행 영향: 9c915dc에서 seed를 실행할 때 셸에 `KNOWSLINK_SYNTHETIC_SIGNUP=1`이 필요하다. 원본 UI는 합성 가입을 쓰지 않았으므로 영향이 없다.

## 발견 사항

critical/high 발견 사항은 없다.

| ID | 심각도 | 위치 | 내용 | 상태 |
|---|---|---|---|---|
| F1(원본) | medium | `internal/relay/identity.go:85-130` | 단일 source의 공유 budget 고갈 | **해소 확인**. 단위·통합·리뷰어 불변식 검사 PASS |
| F4(원본) | low | `.env.example:8` | 예시 env의 합성 가입 기본 열림 | **해소 확인**. `check_compose.py` 단언 PASS. beta.sh seed·공개 403 확인은 OPS 인계 |
| F2(원본) | low | `internal/relay/identity.go` | 재발송 무효화·이메일별 누적 오답 상한 부재 | 미해결 유지. DEV 판단(D02 규칙, 6자리 문구 의존, 8자리 대안) 기록 적합. 일 누적 상한은 designer |
| F3(원본) | low | `internal/relay/member.go:64-65` | 회원 gate 경로 rate 미적용 | 미해결 유지. AGENTS 단계 |
| F-UI-01(원본) | low | 회원 gate 만료 표시 | UTC 기본 문자열과 `+0000` 중복 | 미해결 유지. `pw-11` 직접 확인. DEV 후속 |
| L1(신규) | low | `handovers/logs/2026-10-05_to_dev.md:83,148,194`, `2026-10-05_to_tester.md:81` | 인박스에서 쓴 `../docs/...` 상대 링크가 logs 아카이브 뒤 깨진다. 올바른 경로는 `../../docs/...`다 | 미해결. 원본 기록 소유자 또는 `work.py finish` 처리에서 정정한다. 이 리뷰는 원본을 수정하지 않았다 |

L1 재현: 변경 Markdown의 상대 링크를 파일 위치 기준으로 해석해 존재를 검사했다. 위 4개만 실패했다. 같은 로그의 RATE-FIX 추가 보고 절은 이미 `../../docs/`를 쓴다.

## 검증 및 남은 제약

| 검사 | 대상 SHA | 위치 | 종료코드 | 결과 |
|---|---|---|---|---|
| `review.py prepare --from 59b66ad --to 9c915dc` | 59b66ad..9c915dc | 기록 checkout | 0 | preview 117(대상 64, 제외 53) |
| `npm ci --prefix adapters` | 9c915dc | scratch clone | 0 | — |
| `lint.py --repo <clone> --from 59b66ad --out lint.json` | 9c915dc(clean) | scratch clone | 0 | product-lint(`make lint`) passed exit 0. ERROR 0, WARNING 5, 실행 불가 0 |
| `make test`(product-test) | 9c915dc(clean) | scratch clone | 0 | Go race 4 패키지 ok, adapter 테스트 PASS |
| `make verify-mvp` | 9c915dc(clean) | scratch clone, Compose `knowslink-mvp-481713b41f` | 0 | `TestEmailIdentity` 6개(새 `refused_principal_does_not_spend_shared_rate` 포함), `TestRatePrincipalIsolation`, `TestRateStateBoundedUnderRotation`, PostgresSafety·Gate·`TestTrialHTTP`, TS synthetic/seed/trial-check PASS. `down --volumes` 0 |
| 리뷰어 불변식 검사 | 9c915dc + 리뷰어 파일 1개 | git archive 사본 | 0 | 위 절 |
| `review.py check --task-key SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX` | 59b66ad..9c915dc | 기록 checkout | 0 | 기록 검사 통과: reviewed=64, skipped=53, total=117, lint WARNING=5. snapshot head·detached·clean 검증 포함. jev_find recall 0.667 / precision 0.4(`jev-find-score.json`, 수락 판단에 쓰지 않음) |

명령 종료코드는 파이프 없이 명령 자신의 값이다. `make test`·`verify-mvp`의 HEAD·시각·통과 목록은 `test-run.json`에 보존했다.

### lint와 규모

- product-test 연결: lint.py는 merge-base(59b66ad)의 lint 설정을 쓴다. 그 설정에는 `product-test`가 없다. 후보의 0.9.14 설정이 추가했고 LINT-001 WARNING이 이를 안내한다. 그래서 `make test`를 같은 clean HEAD에서 따로 실행했고 종료코드 0을 확인했다.
- SIZE-001 3건: `PLANS.md` 744줄(기존 경고), `identity_integration_test.go` 333줄, `identity_test.go` 332줄. 테스트 파일은 상한 300을 조금 넘는다. 동작 영향은 없다. 분할은 후속 정리 후보다.
- SIZE-002: 검사 대상 추가 744줄 > 400. RATE-FIX 지시서는 0.9.14 이전 양식이라 예상 규모 항목이 없다. 이 범위는 한 과제가 아니라 RATE-FIX·TESTER·UI·0.9.14 운영 업데이트·PLANS 기록의 통합 delta다. 제품 코드는 `identity.go` +27/−13, 테스트 +101, 스크립트 +4/−2, README·.env.example +5/−4다. 나머지는 기록 문서다. 분할하지 않은 것은 coor가 원본59 이후 통합 후보 하나를 고정했기 때문이며 수용 가능하다. 추가 수정 요청은 없다.

### 기존 증거 재사용

- 원본 리뷰 `25b110f`: 59b66ad 범위 판정과 F1–F4를 그대로 연결했다. 이 delta는 59b66ad 이후만 검토했다.
- DEV의 RED/GREEN 기록(689ba3f·8c7ea94): 구현자 근거로만 참고했다. 수락 근거는 이 리뷰의 9c915dc 재실행이다.
- TESTER·UI: 59b66ad 결과로만 연결했다. 9c915dc의 QA·시각 검수로 표시하지 않는다.

### invalid_lease

- 이 리뷰의 `make verify-mvp` 1회에서 `TestTrialHTTP`는 PASS했다. `invalid_lease`는 관측되지 않았다.
- 1회 PASS는 간헐 실패의 해소 근거가 아니다. 원인 진단은 별도 DEV(SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG)가 진행 중이다. 이 리뷰는 그 실패의 해소나 main 수락을 선언하지 않는다.

### 실행하지 않은 검증

- 실제 운영 SMTP 발송과 일반 사용자 이메일의 로그인·로그아웃: 미실행. 운영 SMTP와 사용자 이메일이 없다.
- 9c915dc의 좁은 독립 QA와 UI 직접 시각 재검수: 미실행. 별도 과제다. invalid_lease 진단 뒤 새 후보가 생기면 그 최신 SHA로 수행한다.
- `make build`·`make verify`·`make verify-runtime`: 리뷰어는 실행하지 않았다. DEV 실행 기록의 8c7ea94 결과(각 0)는 구현자 참고 근거다.
- 공개 경계 QA-P06, Cloudflare edge rate limit, 배포, Cloudflare 쓰기, 노우↔다닷 시험: 범위 밖이다. Workers Free 제한을 유지한다.

## 검토 결론

delta 기술 리뷰 기준으로 조건부 수락 가능하다. 미해결 critical/high는 없다.

- 원본 F1 medium과 F4 low는 9c915dc에서 해소됐다. 독립 재실행과 리뷰어 불변식 검사로 확인했다.
- 원본 TESTER·UI 기록은 근거·범위·실패/미실행 보존이 적합하다. 둘 다 59b66ad 결과이며 9c915dc의 새 QA·시각 검수가 아니다.
- 남은 low는 F2·F3·F-UI-01(원본)과 L1(신규, 아카이브 링크)이다.
- main 수락 전 조건: invalid_lease 진단 결과, 그 뒤 최신 SHA의 delta 리뷰와 좁은 QA, 실제 이메일 조건. 이 리뷰는 이 조건을 충족했다고 보고하지 않는다.

이 결과는 AI 검토자의 판단이다. OCR의 자동 판정이 아니다. `review.py check` 통과는 기록 일치 검사이며 의미 수락이나 테스트 성공을 자동 보증하지 않는다.
