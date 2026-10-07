---
title: 일반 사용자 정리 동작과 snapshot 일관성 수정 기록
status: draft
updated: 2026-10-07
owner: dev
tasks: [SAR-PUBLIC-MESSAGES-001-DEV-FIX-3]
summary: Sol 관측(같은 commit 시각 snapshot 역전·갱신 전 포화 거부·agent ACK의 owner 정리 차단)의 원인·최소 수정·회귀·검증을 기록한다
---

# SAR-PUBLIC-MESSAGES-001-DEV-FIX-3 — 일반 사용자 정리 동작과 snapshot 일관성 수정 기록

## 기준과 기술 계획

- 기준 ref: fixed `758e9f638501f11000ed17c558a9eb54b4350eb3`. 제품 코드는 지시서 HEAD `73deea9`와 같다(758 이후 변경은 `.fullops-squad` 운영 기록뿐이다).
- 읽은 문서: FULLOPS·project·fullops-common-0.3.3(README·coding-style·testing·security)·contexts/dev, SAR-PUBLIC-SERVICE PS-11과 운영 기본값 표(신규16·정리4·rate), DEV-FIX-2 실행 기록, Sol 진행분 `provenance.json`·`sol-narrow-tests.log`·`sol-review-probe_test.go.txt`(읽기만 했다. 실행·수정하지 않았다).
- Sol 관측의 PASS는 관측 성립이다. 제품 적합성 PASS가 아니다. 원 증거와 334 원본 hash는 바꾸지 않았다.
- 계획: 공통 입장 경계(`boundedHTTP`)와 분류(`cleanupOwner`·`cleanupTarget`·`requestBuckets`)·snapshot 저장(`remember`)만 고친다. 제품 수치(신규16·정리4·rate·10s·30s·8/32KiB)·wire·응답 코드·CSRF·현재 키·관계 세대·crash 회수는 바꾸지 않는다. 무제한 fallback·대기열을 만들지 않는다. 제품 판단 변경이 없어 designer 질문은 하지 않았다.
- 라이브러리 근거: 새 의존성은 없다. 표준 `sync/atomic.Pointer`·채널과 기존 pgx/v5·sqlc v1.30.0(`make generate`)을 재사용했다. Context7 조회 대상이 아니다.

## 원인과 수정

### 1. 같은 DB commit 시각에서 늦은 이전 snapshot이 최신을 덮음

- 원인: `remember`가 commit 시각(`clock_timestamp`, µs)만 비교하고 같으면 교체했다. 연속 commit은 같은 시각을 가질 수 있다.
- 수정: 순서 키를 (commit 시각, row epoch)로 바꿨다. epoch는 `SaveRelay`가 commit마다 1 올린다. 시각을 앞에 두어 낮은 epoch로 복원된 DB도 다음 commit에서 snapshot이 갱신된다. 같은 (시각, epoch)는 같은 commit이라 교체하지 않는다.

### 2. 다른 instance의 새 자격·재시작 nil snapshot에서 신규16 포화 중 유효 자기 정리가 429

- 원인: snapshot이 자격이나 대상을 모르면 신규 채널로 보냈다. 신규가 가득하면 즉시 429였다. 갱신은 다음 로컬 commit(1s sweep)뿐이었다.
- 수정: 신규 채널이 가득 찬 때에만 증명 못한 정리 자격에 대해 커밋 상태를 한 번 다시 읽고 다시 판정한다(`refresh`, sqlc `ReadRelay`, lock 없는 SELECT). 읽기는 프로세스당 동시 1개다. 요청 도착 뒤 시작한 읽기를 공유한다. 읽기는 슬롯·DB lock·rate 기록을 쓰지 않는다. 대기자는 자기 연결만 쥐고 10s 기한을 따른다. 위조·타 owner 정리는 다시 읽어도 신규로 남아 429다.
- 중간 실패(보존): 첫 수정 `955cbae`·`9620723`은 신규 여유와 무관하게 다시 읽었다. `9620723`의 verify-mvp에서 `TestValidCredentialCleanupFlood` A가 `attack admission: new 2 clean 1`로 실패했다(`mvp-9620723-flake.log`). row lock에 막힌 transaction들이 pool 연결을 잡아 재조회가 pool을 기다렸다. 그 사이 신규 입장이 늦어졌다. 실제 거부가 생기는 신규 포화 때에만 다시 읽도록 `5d1924c`에서 좁혔다. 여유가 있으면 신규로 대기하고 transaction이 공유 budget을 정확히 판정하며 그 commit이 snapshot을 갱신한다.

### 3. 같은 회원의 agent ACK가 회원 자기 key-revoke를 로컬 입장에서 429로 막음

- 원인 A(공정성 단위): agent·세션·owner 자격이 모두 owner 하나의 정리 슬롯을 공유했다. agent의 ACK 반복이 그 owner의 슬롯을 계속 잡아 owner가 그 agent를 철회하지 못했다(Sol 18/18).
- 원인 B(분류·handler 불일치): `cleanupTarget`은 persist 전 lease의 ACK도 정리로 분류했다. handler는 persist 전 ACK를 `409 invalid_lease`로 거절한다(Sol 결정적 관측).
- 수정 A: 공정성 단위를 owner 자신의 제어(owner ID)와 그 owner의 agent ACK(`<owner>/agents`, `agentsUnit`) 둘로 나눴다. 로컬(`cleaning`)과 공유(`admission.Owner`) 모두 단위당 동시 1개다. 한 회원은 정리 4개 중 최대 2개만 쓴다. 다른 owner를 혼자 막지 못하는 FIX-2 기준을 유지한다. 같은 owner의 여러 agent는 ACK 슬롯 하나를 공유한다. 반복 ACK를 보내는 agent는 owner가 철회할 수 있다.
- 수정 B: ACK 정리 분류에 `m.Persisted`를 더해 handler와 맞췄다. persist 전 ACK는 신규 입장·신규 rate다.

### 정리 caller와 snapshot caller 대조

- `remember` caller는 `transaction` commit과 `refresh` 둘이다. `cleanupOwner` caller는 `boundedHTTP`와 시험이다. `requestBuckets`·`enterHTTP`·`cleanupTarget`의 제품 caller는 `boundedHTTP` 한 곳이다. 모든 HTTP는 `boundedHTTP`를 지난다(FIX-2 대조 유지).
- 정리 handler별 권한·lease·상태를 분류와 다시 대조했다. ACK(`/v1/ack`·`/v1/test/ack`·`/v1/text/ack`)는 같은 `operateAs` persist→ack 경로이며 persist 요구가 유일한 차이였다. 키/agent/owner 철회·unpair·cancel·invite/gate deny·logout의 분류는 handler와 같아 바꾸지 않았다.
- `cmd/relay`의 1s sweep(`Cleanup`)·crash 30s 회수·종료 실패 `orphans` 회수는 바꾸지 않았다. 재조회 상태에는 sweep을 적용하지 않는다. 판정은 세션 만료·lease 기한을 요청 시각으로 다시 본다.

## 자동 검증 증거

증거: [QA 증거](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-DEV-FIX-3/). `.exit`는 명령 반환값이다. 통합 시험은 자기 격리 Postgres 컨테이너 `kl-devfix3-pg`(같은 migration)와 verify-mvp의 고유 Compose project에서만 실행했다. 공유 서버 `knowslink-*`·Tunnel은 건드리지 않았다.

| 검사 | 대상 | 결과 |
|---|---|---|
| 새 통합 회귀 2개 | 758 제품 코드 | exit1 = RED. stale/restart 모두 `own revoke before the sweep: 429 capacity`, owner 제어 `cleanup slots 1 want 2` (red-758.log) |
| 같은 commit 시각 순서(758 API 임시 시험, 커밋 안 함) | 758 | exit1 = RED `a late older snapshot replaced the newer one` (red-758-equal-time.log, 원문 _test.go.txt) |
| 새 통합 회귀 2개 | `955cbae` | exit0, owner 철회 거부 0/18 (green-new-955cbae.log) |
| 변형: ACK persist 검사 제거 | `955cbae` 변형 | exit1 unit 행렬 (mutant-persisted.log) |
| 변형: epoch 동률 비교 제거 | `955cbae` 변형 | exit1 `older commit replaced newer snapshot` (mutant-epoch.log) |
| 변형: 로컬 agentsUnit 제거 | `955cbae` 변형 | exit1 `cleanup slots 1 want 2` (mutant-unit-local.log) |
| 변형: 공유 agentsUnit 제거 | `955cbae` 변형 | exit1 `after the lock: 429` (mutant-unit-shared.log) |
| 전체 통합 `-race` | `955cbae` | exit0, `--- PASS` 47, FAIL/SKIP/DATA RACE 0 (integration-955cbae.log) |
| make lint / test / verify-mvp | `955cbae` | exit0 / exit0 / exit0 (lint/test/mvp-955cbae) |
| make verify-mvp | `9620723` | exit2 중간 실패, 위 2의 원인 (mvp-9620723-flake) |
| 관련 통합 3회 반복 | `5d1924c` 작업본 | exit0 ×3 |
| make lint | 최종 코드 `5d1924c` | exit0 (lint-final) |
| make test | 최종 코드 `5d1924c` | exit0, Go race·adapter 시험 (test-final) |
| make verify-mvp | 최종 코드 `5d1924c` | exit0. Go 통합 `--- PASS` 47·FAIL/SKIP/DATA RACE 0, 실제 Node/MCP·Go owner UI. flood B 4종 유효 정리 거부 0/18(최대 2.6s), agent ACK flood 중 owner 철회 거부 0/18. 자기 Compose 자원 회수 (mvp-final) |
| FullOps lint.py --from 758, deliverables strict, diff check | 기록 마지막 HEAD | 완료 보고에 기록한다. 기록 안에 자기 SHA를 순환 기록하지 않는다 |

새·변경 검사:
- `TestCleanupProvenWhileSnapshotStaleAndNewFull`(통합, stale·restart): 다른 Service가 만든 owner·agent에 대해 신규16 포화 중 자기 키 철회가 sweep 없이 200이다. 같은 자격의 타 owner 철회와 위조 token은 429 capacity다.
- `TestOwnerControlNotHeldByOwnAgentAcks`(통합): row lock 동안 persist된 lease의 agent ACK와 owner의 자기 key-revoke가 각각 정리 슬롯 1개를 잡는다. 같은 단위의 두 번째 ACK·두 번째 owner 제어는 429 capacity다. lock 해제 뒤 둘 다 200이다. lock 없이 48 연결 ACK flood 중 owner 철회 18회가 모두 200이고 공유 기록이 남지 않는다.
- `TestRememberKeepsNewestCommit`(unit): 이전 시각·같은 시각 낮은/같은 epoch는 교체하지 않는다. 같은 시각 높은 epoch와 낮은 epoch의 더 늦은 시각(복원)은 교체한다.
- `TestCleanupAdmissionNeedsVerifiedOwnRecord`(unit): persist된 ACK는 정리, persist 전 ACK는 신규다. 세션 제어는 owner 단위, agent ACK는 `<owner>/agents` 단위다.
- `TestCleanupSnapshotAcrossProcessesAndRestart`: 설명만 바꿨다. `cleanupOwner` 단독은 여전히 commit 전 정리 슬롯을 주지 않는다.

## 경고 처리와 산출물

- 예상 규모: 제품 Go 약 +90/−30줄, 회귀 시험 약 +190줄. 실제: 제품 +95/−34(생성 sqlc 17·SQL 3 포함), 시험 +188/−13이다.
- SIZE-001: `store.go` 524(이전 522)는 기존 대상이며 증가는 필드 2개다. 새 회귀는 `cleanup_unit_integration_test.go`(153줄)로 분리해 `cleanup_flood_integration_test.go`를 221줄로 유지했다. DEP 0: 새 의존성 없음. SLOP·DESIGN 해당 없음: template·CSS 변경 없음. 최종 FullOps lint의 ERROR·WARNING은 완료 보고에 원문대로 남긴다.
- 기술 산출물: D03 architecture, D05 interface-design, D06 data-model, D07 database-design(`ReadRelay` 읽기 질의, table·migration 무변경), D09 crud-design, D10 module-design을 갱신했다. D03 tech-stack은 새 기술이 없어 변경 없음이다. D08 테이블 정의는 변경 없음이다.
- D12/D13은 OPS 소유라 수정하지 않았다. 이행 조건은 아래 미검증·인계에 연결한다. 운영 수락을 만들지 않는다.

## 미검증·인계

- 미검증: 운영 Tunnel/edge·실부하·상태 크기·CPU(L-1, 재조회는 신규 포화 중 전체 JSON 단일 읽기를 더한다), 다수 실제 계정 공모, 운영 배포·실메일·외부 계정·노우↔다닷·실24h. `make verify-grok-plugin`은 adapter 변경이 없어 실행하지 않았다.
- 보존: Sol 진행분·플랫폼 차단 기록·334 원본 hash, OPS H-2 원 high·QA09c 실패, Claude quota 부분 결과, L-1·L-A·L-B. 이 DEV 결과는 독립 QA나 제품 수락이 아니다.
- 후속(coor 배정): 새 fixed SHA의 별도 정적 리뷰, 사용자 승인 Sol 독립 QA, designer UI 영향 확인(화면·문구 변경 없음, 429 안내 동일) 뒤 main 판정. 미해결 critical/high는 계속 main을 차단한다. D12/D13 이행 조건은 OPS 공개 준비에서 반영한다.
