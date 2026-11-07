---
title: "SAR-DB-WAIT-001 — relay DB 대기를 타임아웃, 연결 수, 저장 형식 변경으로 줄인다"
status: draft
updated: 2026-11-07
owner: dev
tasks: [SAR-DB-WAIT-001]
summary: "relay DB 대기를 타임아웃, 연결 수, 저장 형식 변경으로 줄인다"
attempt: b60b839a146841d5b700e8f5bd8be577
base: ebba73c53c284a742b96ef1ee097b0628c8d9fe7
subagent_level: standard
test_level: lite
---

# SAR-DB-WAIT-001 — relay DB 대기를 타임아웃, 연결 수, 저장 형식 변경으로 줄인다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-07
- From / To: coor / dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard. 하위 에이전트 effort는 사용자 지시로 명시하며 기본 high, 기계 작업 medium, max 금지다.
- 승인된 범위 / 추가 승인이 필요한 행위: 이 지시서의 개선 1·2·3과 `fullops/dev` 커밋은 승인됐다. main 병합·push, 운영 배포, 운영 서버 접속은 승인 범위 밖이다.
- 상위 의도: 없음(단일 역할 과제)
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / preamble 값을 사용한다

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

사용자가 2026-11-07 coor의 DB 대기 개선안에 "권장대로 진행해"라고 답했다. 권장안은 "개선 1·2를 먼저 넣고, 개선 3은 측정을 보강한 뒤 진행"이다. 개선 4·5·6은 이번 범위가 아니다.

coor가 ebba73c5에서 직접 확인한 내용:

- 모든 쓰기는 `internal/relay/store.go:147` `transaction`에서 전역 행을 `FOR UPDATE`로 잠그고 행 전체를 다시 쓴다. 요청 하나가 입장(`capacity.go:229`), 처리, 해제(`capacity.go:272`)에서 3번 잠근다. 1초 정리(`cleanup.go:19`)도 프로세스마다 1번 잠근다.
- `lock_timeout`, `idle_in_transaction_session_timeout`, TCP keepalive 설정이 레포에 없다. Rust 평가 측정 때 강제 종료된 클라이언트의 세션이 3분 46초 넘게 남았고, 세션 하나가 행 잠금을 쥐어 새 측정을 막았다. 운영에서 재현한 적은 없다.
- `cmd/relay/main.go:38` `pgxpool.New`는 기본 MaxConns를 쓴다. 12 CPU 측정 호스트에서 12였다. `/healthz`(`pool.Ping`)와 `refresh`도 같은 풀을 쓴다. LISTEN 연결(`text_events.go:77`)은 풀 밖의 별도 연결이지만 풀 설정을 복사한다.
- 측정 근거(로컬 Docker, relay 0.5 CPU/256 MB, Postgres 1 CPU/512 MB, WSL2 디스크):
  - `D:/workspace/KnowsLink/.devtools/rust-eval/results/go.md`: 단계 분해. 메시지 1,000건에서 46.6 ms 중 `save` 24.9 ms, 읽기 8.1 ms다.
  - `D:/workspace/KnowsLink/.devtools/rust-eval/go/`: 실제 `transaction()`을 재는 기존 하니스(`rusteval_test.go.src`, `overlay.json`, `stage_ab.sh`). 픽스처는 `../fixtures/`에 있다.
  - `D:/workspace/KnowsLink/.devtools/db-wait/results/storage.jsonl`, 코드는 `../bench/main.go`: 읽기+저장+commit p50이다. 현재 `jsonb`는 10.1/37.2/330.9 ms, `json`+lz4는 4.6/9.9/61.9 ms다(100/1,000/10,000건). 단일 클라이언트 1회 실행이고 실제 `transaction()` 경로가 아니다.
  - 연결 수 12→4는 처리량 +45%, p99 약 절반이었으나 제품 경로에서 확인하지 못했다. `results/verify/`를 참고한다.

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md`와 연결된 코딩·테스트·보안 규칙이다. 프로젝트 정본은 `.fullops-squad/project.md`, 기준 커밋은 ebba73c53c284a742b96ef1ee097b0628c8d9fe7이다. 공유 시스템 정본은 `docs/system-design.md`의 `DB 동시성·저장·보관`, `Schema·writer 보호와 운영 격리` 절이다.
`D01~D13`의 정의는 `.fullops-squad/docs/deliverables/README.md`에 있다. 예외 없음.

## 먼저 읽을 문서

코드: `internal/relay/store.go`, `capacity.go`, `cleanup.go`, `cleanup_admission.go`, `text_events.go`, `cmd/relay/main.go`, `internal/config/config.go`, `db/queries/relay.sql`, `db/migrations/00001_relay.sql`, `sqlc.yaml`, `scripts/schema.py`, `scripts/verify_runtime.py`, `deploy/knowslink/beta.sh`, `deploy/knowslink/compose.ops.yaml`.
문서: `docs/system-design.md`, `docs/architecture.md`, `.fullops-squad/docs/design-docs/database-design.md`, `data-model.md`, `tech-stack.md`, `.fullops-squad/docs/generated/db-schema.md`, `.fullops-squad/project.md`, 공통 규칙 4개.
필요 시 확인(제외 추천): `db/migrations/00002_human_writer.sql`. 다만 가드를 옮기는 개선 3에서는 반드시 읽는다.
지시 전제와 충돌 — 먼저 확인: 위 문서에 "JSONB 한 행"이라는 설명이 있다. 개선 3 구현 때 현재 상태에 맞게 고칠 대상이고, 지시를 뒤집지 않는다.
결과 파일: `docs/evaluations/jev/SAR-DB-WAIT-001-context.json`, `-find.json`, `-documents-find.json`, `-packet.json`.

## 해야 할 일과 파일 소유권

파일 소유권은 `cmd/`, `internal/`, `db/`, `scripts/`, `sqlc.yaml`, `compose.yaml`, `deploy/knowslink/`, 그리고 `docs/`와 `.fullops-squad/docs/`의 해당 줄이다. 먼저 이 과제 기록에 짧은 기술 계획을 적는다.

기술 계획(dev, 2026-11-07):

1. 설정은 새 파일 `internal/relay/db_limits.go`의 `PoolConfig` 하나에 둔다. `cmd/relay`가 이 함수로 풀을 만들고, LISTEN 연결은 기존 `ConnConfig.Copy()`로 같은 값을 받는다. 서버 쪽 세션 값은 RuntimeParams로, relay 쪽 keepalive는 `net.Dialer.KeepAliveConfig`로 넣는다.
2. 동시 트랜잭션 제한은 `transaction()` 입구의 채널 슬롯 3개다. 대기는 ctx와 상한 시간을 따르고, 초과하면 `unavailable`이다.
3. 한 요청의 최악 시간은 입장 대기 + 처리 대기 + 해제 2초다. 이 합이 WriteTimeout 10초 안에 들도록 값을 정하고, 시험으로 고정한다.
4. 측정은 기준 커밋 바이너리(12 연결), 개선 1·2 바이너리(4 연결), 개선 3 후보 바이너리를 같은 Postgres 컨테이너의 별도 DB에서 번갈아 돌린다. 하니스는 `.devtools/db-wait/`에 두고 실제 `transaction()`만 잰다.
5. 개선 3은 작업 트리에서 후보를 만들어 측정하고, 판정을 통과할 때만 커밋한다.

- [ ] 개선 1 멈춤 차단: relay의 모든 DB 연결에 `lock_timeout`, `idle_in_transaction_session_timeout`, TCP keepalive를 설정한다. 값은 측정으로 정한다. `idle_in_transaction_session_timeout`은 상한 20,000건에서 lock 안 계산 시간보다 커야 한다. `lock_timeout`은 HTTP `WriteTimeout` 10초 안에서 실패를 돌려줄 수 있어야 한다. 시간 초과는 기존 `unavailable` 응답으로 돌아가야 한다.
- [ ] 개선 2 연결 수: 풀 상한을 4로 내린다. 프로세스 안에서 동시 트랜잭션 수를 제한해 연결 하나를 비워 둔다. 대기는 ctx를 따르고 상한 시간이 있어야 한다. 트랜잭션이 모두 잠금 대기 중이어도 `/healthz`는 응답해야 한다. 실제 `transaction()` 경로에서 12와 4를 비교해 효과를 기록한다. 효과가 없거나 나빠지면 이유와 숫자를 적고 되돌린다.
- [ ] 측정 보강: 기존 하니스로 실제 `transaction()` 경로의 현재 `jsonb`와 `json`+lz4를 비교한다. 메시지 100/1,000/10,000건, 1·4·16 클라이언트, 운영 제한 조건이다. 같은 셀을 번갈아 3회 이상 돌려 호스트 흔들림을 적는다.
- [ ] 개선 3 판정: 1,000건 단일 클라이언트 total p50이 현재의 50% 이하이고 16 클라이언트 처리량이 낮아지지 않으면 진행한다. 조건을 못 채우면 3을 구현하지 않고 숫자를 보고한다. 이 경우에도 개선 1·2는 완료한다.
- [ ] 개선 3 구현(판정 통과 시): 새 goose migration으로 `data`를 `json` + lz4로 바꾸고 기존 행을 보존한다. HumanModel 가드는 숫자 컬럼으로 옮겨 `data`를 파싱하지 않게 한다. 구버전 writer가 가드를 우회해 쓰지 못하도록 컬럼 이름도 바꿔 SQL 오류로 실패하게 한다. Down은 지금처럼 실패 고정으로 둔다. `db/queries/relay.sql`, 생성된 `internal/database`, `text_events.go:167`의 직접 SELECT, 통합 시험의 raw SQL과 temp table, `scripts/verify_runtime.py`, `deploy/knowslink/beta.sh`를 함께 맞춘다. `make generate`와 `make schema` 결과를 커밋한다.
- [ ] 문서: `docs/system-design.md`와 산출물의 해당 줄을 고친다. 개선 3을 구현했으면 배포 때 필요한 백업과 되돌리기 제약을 운영 문서에 적는다.

예상 규모: 개선 1·2는 코드 약 100줄과 시험이다. 개선 3은 migration, 쿼리, 시험 수정이다. SIZE-002가 나오면 실제 차이와 분할하지 않은 이유를 보고한다.

## 완료 기준과 검증

모두 채울 때까지 중간 확인 없이 진행한다. 담당은 전부 DEV이고, 대상은 완료 SHA다.

- 잠금을 쥔 세션이 멈춘 상황에서 새 요청이 설정한 상한 안에 `unavailable`로 끝난다. 유휴 트랜잭션 세션은 종료되고 이후 쓰기가 다시 성공한다. Postgres 통합 시험으로 증명한다.
- 잠금 대기 중에도 `/healthz`가 2초 안에 응답한다. 시험으로 증명한다.
- 정상 부하에서 새 타임아웃으로 인한 오류가 0이다. 1,000건 16 클라이언트 30초로 확인한다.
- 개선 3을 구현한 경우 다음을 증명한다.
  - 기존 `jsonb` 행이 migration 뒤 같은 State로 읽힌다.
  - `pg_column_compression`이 lz4다.
  - HumanModel을 낮추는 쓰기와 구버전 쿼리가 거부된다.
  - envelope digest, Canonical 시험이 통과한다.
  - `beta.sh`의 백업·복원 확인 경로가 새 스키마에서 동작한다.
- 등록 검사: Windows에는 make가 없으므로 완료 SHA에서 Linux 컨테이너로 `make lint`, `make test`, `make build`, Postgres race integration을 돌려 종료코드를 남긴다. 이미지는 `kl-refactor003-linux`를 재사용할 수 있다. 파이프로 종료코드를 가리지 않는다. Docker VM 시계 보정으로 `abnormal_clock` 시험 2건이 간헐 실패한 이력이 있다. 같은 실패가 나오면 재실행 결과를 함께 적는다.
- FullOps lint: `python3 <스킬 기준 ../../scripts/lint.py> --repo <dev 체크아웃> --from ebba73c53c284a742b96ef1ee097b0628c8d9fe7` ERROR 0. make 실행 불가는 사유와 Linux 결과를 적는다.

증명 방법: `docs/evaluations/qa-reports/SAR-DB-WAIT-001-DEV-test/report.md`에 위 결과와 측정표(현재, 개선 후, 셀별 p50/p95/tx/s, 표본 수, 호스트 흔들림)를 적는다. 측정 원본은 `D:/workspace/KnowsLink/.devtools/db-wait/`에 둔다.
독립 리뷰와 main 병합은 coor가 완료 보고 뒤 진행한다. 제품 최종 수락과 운영 효과 확인은 운영 배포 뒤의 별도 판단이다.

## 갱신할 산출물

route 추천은 D10, D07, D08이다. 개선 3을 구현하면 D08은 `make schema`로 다시 만든다. 운영 절차를 바꾸면 D12의 해당 줄을 더한다. 본문 검토 뒤 해당 없는 ID는 이유를 보고한다.

## 기대 산출물

코드와 시험, migration(개선 3 구현 시), `docs/evaluations/qa-reports/SAR-DB-WAIT-001-DEV-test/report.md`, `docs/evaluations/jev/SAR-DB-WAIT-001-packet-outcomes.json`, 갱신한 문서다. 작업 도구와 로그는 레포 밖 임시 경로가 아니라 gitignore된 `D:/workspace/KnowsLink/.devtools/db-wait/`에 둔다.

## 제약·협업·후속

- 범위 밖이며 하지 않는다: 입장·해제 분리(개선 4), 잠금 없는 읽기(개선 5), 정규화(개선 6), 전역 `synchronous_commit=off`, relay CPU 제한 변경, Rust 전환, 운영 배포, 운영 서버 접속.
- 사용자가 개선 4·5용 결정을 확정했다. 이 과제에서 구현하지 않으며 PLANS.md에 기록돼 있다.
- 다른 프로젝트의 컨테이너는 건드리지 않는다. `kl-` 접두어가 아닌 컨테이너를 멈추지 않는다.
- shell_guard가 `docker rm`을 막는다. 새 컨테이너는 `--rm`으로 띄우고, 우회하지 않는다.
- 비밀번호, 키를 로그와 문서에 쓰지 않는다. 시험 계정 이메일을 레포 문서에 쓰지 않는다.
- 멈추고 물을 조건: 워크트리 밖 변경, 기존 데이터 삭제가 필요한 migration, force-push, 파일 소유권 밖 수정, 제품 규칙 변경, 설명되지 않는 검증 실패.
- 라이브러리 설정은 프로젝트 버전(pgx v5.10.0, Postgres 17)을 기준으로 Context7에서 확인하고 근거를 보고서에 적는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-DB-WAIT-001-packet.json` / SHA `7cbfa6221f0f24c6a13f3758a4fdab29ea22e749` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `docs/system-design.md` (document_read) · 줄 75, 76, 76, 79, 99, 168, 245, 266, 267, 269 · inferred · 필수
- `.fullops-squad/docs/design-docs/database-design.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md` (document_read) · 줄 52 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX.md` (document_read) · 줄 29 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-003.md` (document_read) · 줄 30 · inferred
- `.fullops-squad/docs/generated/db-schema.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 37, 37, 62, 62, 80, 84 · inferred
- `cmd/relay/main.go` (direct_edit) · 줄 20, 29, 103 · inferred
- `cmd/relay/main_test.go` (impact_check) · 줄 24 · inferred
- `db/queries/relay.sql` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/capacity.go` (direct_edit) · 줄 71, 77, 107, 297, 302 · inferred
- `internal/relay/public_messages_integration_test.go` (impact_check) · 줄 418 · inferred
- `internal/relay/store.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/text_events.go` (direct_edit) · 줄 전체/미확인 · inferred
나머지 참조 후보 0건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 9건: route: task/role/SHA/attempt/instruction differs or absent; regenerate/bind route; find: task/role/SHA/attempt/instruction differs or absent; regenerate/bind route; documents-find: task/role/SHA/attempt/instruction differs or absent; regenerate/bind route; context: task/role/SHA/attempt/instruction differs or absent; regenerate/bind route; .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection …. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-test/report.md`(역할별로 나눈 키는 `<키>-DEV` 그대로) 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. worker_done 본문 첫 줄은 `[완료] <과제 키> | 브랜치 <이름> | SHA <브랜치 HEAD>` 한 줄로 쓰고 SHA는 하나만 적는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
