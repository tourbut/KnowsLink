---
title: SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG 실행 기록
status: draft
updated: 2026-10-06
owner: dev
tasks: [SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG]
summary: TestTrialHTTP invalid_lease 간헐 실패의 원인과 verify-mvp 격리 수정 및 검증을 기록한다
---

# SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG — TestTrialHTTP invalid_lease 원인과 수정

## 기준

- 브랜치 `fullops/dev`. 시작 HEAD `ff1e6704ec78ff1e44294bb194fb08fb6cc1b9a9`. RATE-FIX 후보 `f364d48417b58c69969a4765b88324724eeb5c78`은 시작 HEAD의 조상이다. lint 기준 ref `9c915dc71e2a872243ffec294126d4668b4d32a4`.
- Task `task_8703a6250fa8`, Dispatch `ctx_5f78e35e77c2`, coordinator `term_6895aaf1-7b43-4fe0-a416-76f1255a5946`, Run `run_8ca8bc058ab7`.
- 적용 규칙: `fullops-common-0.3.3`, FULLOPS.md, project.md, coding-style/testing/security. Ponytail full. 예외 없음.
- 원래 실패: [RATE-FIX 실행 기록](SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX.md)의 `make verify-mvp` 1회차 `test_messages_integration_test.go:54: POST /v1/test/persist got 409 invalid_lease`. 원본 기록은 수정하지 않았다.

## 기술 계획

1. 실패 경로 추적: `TestTrialHTTP` → `/v1/test/persist` → `operation("test-persist")` → `State.leased`. `invalid_lease` 조건과 lease를 지우는 호출자(`sweep`, `ack`)를 찾는다.
2. 같은 DB의 다른 writer를 찾는다. `scripts/verify_mvp.py`와 `compose.yaml`, `cmd/relay/main.go`의 `Cleanup`을 확인한다.
3. 결정적 재현 검사와 통계 재현(relay 실행/정지 비교)으로 원인을 확정한다.
4. 제품 규칙을 바꾸지 않고 검사 격리를 고친다.

## 원인

- `verify_mvp.py`는 `compose up --wait relay` 뒤 같은 project의 Postgres로 Go integration 검사를 실행했다. 이때 Compose relay가 계속 실행 중이었다.
- 그 relay는 `KNOWSLINK_TEST_AGENTS=trial_codex,trial_grok`으로 시작한다. `Service.Cleanup`은 1초마다 빈 transaction을 실행한다. `transaction`은 `state.TestAgents = s.TestAgents`를 설정하고 `sweep`한 뒤 상태 행을 저장한다.
- `TestTrialHTTP`는 자기 Service의 allowlist를 `agent_a,agent_b`로 둔다. relay의 sweep에서 이 메시지는 `current`가 false다(`relay.test.message`의 양쪽이 allowlist에 없음). 그래서 sweep이 `failed:revoked`로 바꾸고 `Envelope`·`LeaseToken`을 지운다.
- sweep tick이 검사 send~claim 구간(약 0.1초)에 들어가면 실패한다. 위치에 따라 증상이 다르다. pull 전이면 envelope nil panic이다. pull 뒤면 persist/ack의 `409 invalid_lease`다. ack 뒤면 claim의 `403 sender_not_allowed`다.
- 분류: 제품 결함이 아니다. allowlist 밖 시험 메시지를 회수하는 것은 기존 제품 규칙이다(`current`). 운영은 relay 하나가 같은 설정으로 실행된다. 원인은 검사 격리 결함(같은 DB의 설정이 다른 두 writer)과 1초 tick 시간 조건이다. 테스트 코드의 계약은 맞다.
- 다른 Go integration 검사는 `relay.test.message`를 쓰지 않아 이 sweep의 영향을 받지 않았다. 그래도 같은 DB의 외부 writer를 없애므로 모든 integration 검사의 격리가 같이 확보된다.

## 재현 (수정 전 verify 순서, 격리 Compose project `knowslink-diag-*`)

재현 스크립트는 scratch의 `repro.py`다. 고유 project를 만들고 relay 실행 상태와 정지 상태에서 같은 검사를 반복한 뒤 `down --volumes`로 회수한다. 기존 `knowslink-*` 컨테이너와 Tunnel은 건드리지 않았다.

| 조건 | 명령 | 결과 |
|---|---|---|
| relay 실행 | `go test -tags=integration -race -count=1 -v -run 'TestTrialHTTP$'` 40회 | 4/40 실패: `:55 POST /v1/test/ack got 409 invalid_lease` 1, `:56 POST /v1/test/claim got 403` 1, pull nil panic 2 |
| relay 정지 | 같은 명령 40회 | 0/40 실패 |
| relay 실행 | `-count=200 -run 'TestTrialHTTP$\|TestTrialForeignAllowlistRevokesLease'` 1회 | exit 1. 원래와 같은 `:54: POST /v1/test/persist got 409 {"error":"invalid_lease"}` 포함 FAIL 3줄, panic으로 중단 |
| relay 정지 | 같은 명령 1회 | exit 0. 두 검사 각 200회 PASS |

결정적 검사 `TestTrialForeignAllowlistRevokesLease`는 relay 정지 상태에서도 원인을 재현한다. pull 직후 verify-mvp relay와 같은 allowlist의 두 번째 Service가 `Cleanup`과 같은 빈 transaction을 한 번 실행한다. 그 뒤 persist가 `409 invalid_lease`다.

## 수정

| 파일 | 변경 |
|---|---|
| `scripts/verify_mvp.py` | Go integration 검사 전에 `compose stop relay`, 뒤에 `compose up --wait relay`. TS 검사는 다시 시작된 relay로 그대로 실행한다 |
| `internal/relay/test_messages_integration_test.go` | `TestTrialForeignAllowlistRevokesLease`: 원인 메커니즘과 allowlist 밖 lease 회수 규칙을 고정한다 |
| `README.md`, D10 `module-design.md` | verify-mvp가 Go 검사 동안 relay를 멈추는 이유 |

제품 코드(`internal/relay/*.go` 비시험 파일)는 바꾸지 않았다. `TestTrialHTTP`도 바꾸지 않았다. UI 영향 없음.

## 검증 (작업 트리 = 최종 코드 커밋 내용)

명령은 레포 루트에서 실행했다. 종료코드는 명령 자신의 값이다. 로그는 레포 밖 scratch에 두었다.

| 명령 | 종료코드 | 결과 |
|---|---|---|
| `make lint` | 0 | PASS (gofmt·vet·adapters·config) |
| `go vet -tags=integration ./internal/relay/` | 0 | 새 integration 검사 컴파일·vet |
| `make test` | 0 | PASS |
| `make verify-mvp` (수정 후) | 0 | `stop relay` exit 0 → Go integration 전체 PASS(`TestTrialHTTP`, `TestTrialForeignAllowlistRevokesLease`, `TestEmailIdentity`, `TestRatePrincipalIsolation`, `TestRateStateBoundedUnderRotation`, PostgresSafety·Gate·Trial 포함) → `up --wait relay` exit 0 → TS synthetic/seed/trial-check exit 0 → `down --volumes` exit 0. 8개 command 모두 exit 0 |
| `git diff --check` | 0 | — |

`make build`·`make verify`·`make verify-runtime`·`make generate`·`make schema`·`make verify-grok-plugin`은 제품 코드·SQL·adapter·런타임 설정 변경이 없어 실행하지 않았다. `make verify-mvp`의 build 단계는 실행되었다.

## 한계와 후속

- 통계 재현은 확률 근거다. 결정적 근거는 `TestTrialForeignAllowlistRevokesLease`와 코드 경로다.
- 합성 fixture 검증이다. 실제 운영 relay·노우↔다닷 시험·배포·Cloudflare 쓰기는 수행하지 않았다.
- 독립 delta 리뷰와 좁은 QA는 coor가 배정한다. 범위: `verify_mvp.py` 순서, 새 integration 검사, README·D10 문구.
