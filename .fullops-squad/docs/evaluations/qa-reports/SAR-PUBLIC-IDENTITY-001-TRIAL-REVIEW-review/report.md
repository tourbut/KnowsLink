# SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW 리뷰

- 검토자 / CLI / 모델: OPS, Claude Code, `claude-opus-5-5` high. 세션 `84d2e5f7-312b-4757-b578-df436e08f493`. Task `task_5c53d2c3739c`, Dispatch `ctx_ade6aee10fe1`.
- 구현자 세션: DEV-TRIAL-DIAG `fdc4e778-7118-413f-a46f-a9714ff64b1f`(DEV 프로젝트 JSONL에서 `task_8703a6250fa8`와 `2111ff4` 커밋 메시지로 확인). 검토자와 다르다.
- base SHA / head SHA / merge-base: `9c915dc71e2a872243ffec294126d4668b4d32a4` / `eb2e34b93fe8d20fa1cd9166f73ff68d14bf17de` / `9c915dc71e2a872243ffec294126d4668b4d32a4`.
- snapshot: `/tmp/knowslink-identity-trial-review-eb2e34b`, detached HEAD `eb2e34b`, clean, read-only. 설치·빌드·테스트·수정은 하지 않았다.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11, `review/rule.json`(rule_sha256 `ab2116fb…`). 공통 기준 `fullops-common-0.3.3`과 coding-style/testing/security, `project.md`. 예외 없음.
- 요구사항·완료 기준 원천: `handovers/to_ops.md`(준비 커밋 `f45b3e9`), DEV-TRIAL-DIAG 실행 기록, RATE-REVIEW report.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 23 / 23 / 0 / 15 / 8.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: 0 / 1 / 0. WARNING은 `PLANS.md` SIZE-001(764줄, 기존 길이 경고)이다. 이번 리뷰 범위의 수정 대상이 아니다.
- SIZE-002: 해당 없음. lint 대상 추가 180줄 중 코드·README는 33줄이고 나머지는 리뷰·Jev·계획 기록이다.
- DEP-001: 해당 없음. `go.mod`·`package.json`·lockfile 변경이 없다.
- UI 디자인: 해당 없음. UI 템플릿·상태 문구·공용 컴포넌트 변경이 없다. 디자인 lint와 공용 테마 전환은 미구성이다. 원본 직접 시각 검수 UI `cf0ab09`는 fixed59 결과로만 연결한다. 새 시각 검수로 표시하지 않는다.

## 검토 범위

- 코드: `scripts/verify_mvp.py`, `internal/relay/test_messages_integration_test.go`. 관련 경로 `store.go`의 `transaction`·`sweep`·`current`·`leased`, `cleanup.go`, `integration_test.go`의 `setup`을 snapshot에서 읽었다.
- 기술 문서: `README.md`, D10 `module-design.md`, `contexts/dev.md`, DEV-TRIAL-DIAG 실행 기록.
- 결과 기록: RATE-REVIEW 리뷰 디렉터리(77dd464)와 실행 기록, 아카이브 로그 2개, PLANS·board 갱신. 원본 범위는 재리뷰하지 않고 기록 정합성만 확인했다.
- skipped 8: Jev 라우팅·탐색 생성물 5개, RATE-REVIEW의 preview/rules/jev-find-score 생성물 3개. 제품 동작과 수락 판단에 영향이 없다. 존재·형식·참조를 확인했다.
- 변경 Markdown 9개의 상대 링크를 snapshot에서 검사했다.

## 발견 사항

| ID | 심각도 | 파일·줄 | 내용 | 상태 |
|---|---|---|---|---|
| L2 | low | `handovers/logs/2026-10-06_to_dev.md:72` | 완료 보고의 `../docs/exec-plans/...` 링크가 logs/ 아카이브 뒤 깨진다. 올바른 경로는 `../../docs/...`다. L1과 같은 유형이다 | 미해결. 원본 소유자(coor/dev) 또는 finish 처리에서 정정 |

critical/high/medium 없음.

### 코드 검토 근거

- 원인 일치: Compose relay는 `KNOWSLINK_TEST_AGENTS=trial_codex,trial_grok`으로 `Cleanup`을 1초마다 실행한다. `transaction`은 자기 allowlist로 `sweep`한 뒤 상태 행 전체를 저장한다. `TestTrialHTTP`의 `agent_a→agent_b` 시험 메시지는 그 allowlist에서 `current=false`다. 그래서 `failed:revoked`가 되고 `leased`가 `invalid_lease`를 반환한다. 원래 `:54 persist 409 invalid_lease` 실패와 맞다.
- 격리 수정: Go integration 전 `compose stop relay`, 뒤 `compose up --wait relay`. relay의 DB writer는 Cleanup뿐이다. Postgres는 계속 실행되어 Go 검사 연결이 유지된다. Go 검사 실패 시에도 `finally`의 `down --volumes`가 회수한다. 제품 비시험 코드 변경이 없다.
- 새 검사: `TestTrialForeignAllowlistRevokesLease`는 pull 직후 다른 allowlist Service로 빈 transaction을 한 번 실행하고 persist `409 invalid_lease`를 확인한다. 시간 조건이 없어 결정적이다. 같은 allowlist면 sweep이 회수하지 않아 persist가 200이 되어 실패한다(코드 추론, 음성 대조는 실행하지 않음). `setup`이 상태 행을 초기화하므로 다른 검사와 독립이다.
- 한계(발견 아님): 이 검사는 메커니즘을 고정한다. `verify_mvp.py`에서 stop을 지워도 이 검사 자체는 실패하지 않는다. 격리 회귀는 README/D10 설명과 verify-mvp의 간헐 실패로만 드러난다. 현재 규모에서 수락 가능하다.

### 기록 정합성

- DEV 실행 기록의 원인·재현·수정·검증과 코드가 일치한다. 원래 실패를 보존하고 PASS로 재작성하지 않았다. `2111ff4`와 `00a1384`의 제품 내용 동일을 `git diff 2111ff4 00a1384 -- ':!.fullops-squad'` 빈 결과로 확인할 수 있다(완료 커밋은 문서만 바꿈).
- RATE-REVIEW 실행 기록·report·result·lint·test-run은 서로 일치한다. PLANS 수신 기록(117개, critical/high 0, F1/F4 해소)과 일치한다.
- PLANS·board의 보류 조건과 실제 이메일의 운영 후속 경계가 지시서와 일치한다.

## 검증 및 남은 제약

scratch clone(`eb2e34b` detached, 실행 전후 clean)에서 실행했다. 종료코드는 명령 자신의 값이다.

| 명령 | 종료코드 | 결과 |
|---|---|---|
| `npm ci --prefix adapters` | 0 | 의존성 복원 |
| `lint.py --from 9c915dc --out lint.json` | 0 | product-lint exit 0, product-test(`make test`, kind: test) exit 0, ERROR 0, WARNING 1, 실행 불가 0 |
| `make verify-mvp` | 0 | project `knowslink-mvp-a3723f3371`. stop relay 0 → Go integration 20개 PASS(`TestTrialHTTP`, `TestTrialForeignAllowlistRevokesLease` 포함) → up --wait relay 0 → TS 3개 0 → down --volumes 0 |

- coordinator의 `SAR-PUBLIC-IDENTITY-001-final/candidate-lint.json`(같은 HEAD, ERROR 0, WARNING 1)과 독립 lint 결과가 같다.
- 실행하지 않음: `make build`·`make verify`·`make verify-runtime`(제품 코드·SQL·adapter·런타임 설정 변경 없음), `TestTrialHTTP` 반복 통계 재현(DEV 40/40·count=200 기록 재사용), 같은 allowlist 음성 대조.
- 실제 이메일·운영 공개·배포·Cloudflare 쓰기·노우↔다닷 시험은 수행하지 않았다. 전체 일반 서비스 수락의 운영 후속이며 로컬 코드 수락의 선행 조건이 아니다.
- 좁은 QA는 별도 tester 결과와 결합한다. 남은 low: F2·F3·F-UI-01(원본), L1(RATE 리뷰), L2(신규).

## 검토 결론

고정 후보 `eb2e34b`의 `9c915dc` 이후 delta는 수락 가능하다. 미해결 critical/high가 없다. 격리 수정은 원인과 맞고 제품 비시험 코드를 바꾸지 않는다. 새 결정적 검사가 원인 메커니즘을 고정한다. 독립 lint·test·verify-mvp가 통과했다. 원본 TESTER `9e2654d`/UI `cf0ab09`는 fixed59 결과로, RATE 리뷰 `77dd464`는 `9c915dc` 결과로만 재사용한다. AI 검토 판단이며 OCR의 자동 판정이 아니다. check 통과는 기록 검사이며 테스트 성공을 자동 보증하지 않는다.
