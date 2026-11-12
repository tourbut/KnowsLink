---
title: "SAR-MSG-EPHEMERAL-001 — 메시지·대화 사본·Dots Inbox를 Postgres 대신 휘발성 Redis에만 TTL로 두고, 영구 저장하지 않음을 명시한다"
status: draft
updated: 2026-11-12
owner: dev
tasks: [SAR-MSG-EPHEMERAL-001]
summary: "메시지·대화 사본·Dots Inbox를 Postgres 대신 휘발성 Redis에만 TTL로 두고, 영구 저장하지 않음을 명시한다"
attempt: 233358aeaabb4ee9999c668776723109
base: 17acfa5e0b9800522854163b9d5df67b9870a159
subagent_level: standard
test_level: lite
---

# SAR-MSG-EPHEMERAL-001 — 메시지·대화 사본·Dots Inbox를 Postgres 대신 휘발성 Redis에만 TTL로 두고, 영구 저장하지 않음을 명시한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-12
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard (fullops.json 선택값)
- 승인된 범위: 아래 `해야 할 일`의 코드·migration·문서·문구 변경, 로컬 Docker(Redis·Postgres) 시험, `fullops/dev` 브랜치 커밋.
- 추가 승인이 필요한 행위: 운영 서버 접속·배포, main 병합·push, 기존 운영 백업 삭제, 아래 결정과 다른 보관 기한·범위.
- 상위 의도: 없음(단일 과제)
- 담당 repo id / 워크트리 / 브랜치: `0b08ec4c-9e3a-4613-8197-5a835e545335` / `C:/Users/shin/orca/workspaces/KnowsLink/dev` / `fullops/dev`
- 병합 책임자 / 기본 브랜치: coor / `main`
- 복귀 repo id / 워크트리 / 터미널 / run·task·dispatch id: `0b08ec4c-9e3a-4613-8197-5a835e545335` / `C:/Users/shin/orca/workspaces/KnowsLink/coor` / Orca preamble의 값 / preamble의 값

## 현재 상황과 확인 근거

- 사용자 요청: "메시지를 postgre에 저장하지말고 redis에만 임시저장시키는 방식으로하자 db에 영구저장은 필요없을꺼같고 저장도 안한다고 명시하자".
- 사용자 결정(FullOps 결정 기록 `D049`~`D052`이며 산출물 ID `D01~D13`과 다르다):
  - D049: 웹 대화 사본은 최대 1시간(현재 24시간), 전달 기한은 최대 180초를 유지한다.
  - D050: Dots 원격 Inbox도 Redis 임시 보관으로 옮긴다. 서명키·구독 같은 연결 정보는 Postgres에 남긴다.
  - D051: Redis는 휘발성이다(디스크 저장 없음). 재시작 때 미전달 메시지와 사본 유실을 받아들인다. 문구는 "영구 저장하지 않음"이다.
  - D052: 운영 전환 때 Postgres에 남은 기존 메시지·사본·Inbox는 삭제한다.
- 현재 저장 위치(기준 `17acfa5e`):
  - `internal/relay/store.go:54` `Message`(Envelope, Inbox, Conversation 사본, lease·claim)와 `store.go:90-91` `State.Messages`·`State.Idempotency`가 `relay_state` 한 행에 있다.
  - `internal/relay/remote_auth.go:40-45` `remoteInstall.Inbox`가 Dots 설치별 Inbox다.
  - 상한과 보관: `internal/relay/capacity.go:84-100`(20,000건·대기 100건), sweep은 `store.go:285-330`.
  - 도착 신호: `internal/relay/text_events.go`의 Postgres LISTEN/NOTIFY.
  - 문구: `internal/relay/policy_pages.go`(개인정보·약관), `i18n_en_human.go:19`, `i18n_en_member.go:243`(24시간 보관 안내), `README.md` "현재 기능과 확인 범위", `docs/architecture.md` "저장과 보관 기한".
  - 백업: `deploy/knowslink/beta.sh`의 `pg_dump`와 `prune-backups.sh`. 지금 백업에는 메시지가 들어 있다.
- 선행: SAR-DB-WAIT-001(migration 3, `state` json+lz4)은 main에 있고 운영 미적용이다. 이번 migration은 그 뒤 번호다. 메시지를 빼면 상태 행이 작아져 DB 잠금 시간도 줄어든다.

## 적용 기준과 예외

- 공통 기준: `.fullops-squad/rules/common/README.md`와 연결된 코딩·테스트·보안 규칙, `AGENTS.md`(작업 파일 위치, Docker 정리), `.fullops-squad/project.md`, 기준 `17acfa5e0b9800522854163b9d5df67b9870a159`.
- 산출물 ID `D01~D13`은 `.fullops-squad/docs/deliverables/README.md`, lint 코드는 `.fullops-squad/lint/README.md`를 따른다.
- 보안·권한 기준은 낮추지 않는다. 전달 때 현재 키·owner·권한·generation 검사, 철회·친구 해제·차단 뒤 전달 차단, 서명·멱등 의미는 그대로다.
- 예외: 새 Go 의존성(Redis client) 추가를 허용한다. 표준 라이브러리 대안이 없는 이유와 버전 고정을 DEP-001 항목으로 적는다. 라이브러리 API는 Context7로 확인한다.

## 먼저 읽을 문서

- `internal/relay/store.go`, `internal/relay/capacity.go`, `internal/relay/human_messages.go`, `internal/relay/public_text.go`, `internal/relay/http.go`(메시지 경로), `internal/relay/remote_auth.go`, `internal/relay/remote_mcp.go`, `internal/relay/remote_events.go`, `internal/relay/text_events.go`, `internal/relay/cleanup.go`
- `internal/relay/policy_pages.go`, `internal/relay/i18n_en_human.go`, `internal/relay/i18n_en_member.go`
- `db/migrations/00003_relay_state_json.sql`, `db/queries/relay.sql`, `cmd/relay/main.go`, `compose.yaml`, `deploy/knowslink/compose.ops.yaml`, `deploy/knowslink/beta.sh`
- `docs/architecture.md`, `docs/system-design.md`, `.fullops-squad/docs/design-docs/database-design.md`
- 탐색 패킷이 있으면 `docs/evaluations/jev/SAR-MSG-EPHEMERAL-001-packet.json`을 함께 본다.

## 해야 할 일과 파일 소유권

파일 소유권: 이 레포의 `cmd/`, `internal/`, `db/`, `adapters/`(필요할 때만), `compose.yaml`, `deploy/knowslink/`, `scripts/`, `docs/`, `README.md`, `.fullops-squad/docs/`. `.fullops-squad/handovers/`는 완료 보고만 쓴다.

- [ ] 기술 계획을 이 인박스 끝에 짧게 적는다(Redis 키 구조, TTL, 원자 연산 방식, LISTEN/NOTIFY 유지 여부).
- [ ] 메시지 단위 기록(본문·envelope·Inbox 사본·대화 사본·receipt·lease/claim·멱등 기록)과 Dots 원격 Inbox를 Redis로 옮긴다. Postgres에는 어떤 형태로도 메시지 본문·envelope·사본을 쓰지 않는다.
- [ ] TTL: 전달 기한 180초 유지, 웹 대화 사본 최대 1시간, Redis 안 모든 메시지 기록의 상한 1시간. 기존 20,000건·대기 100건 상한은 유지하거나 같은 의미로 옮긴다.
- [ ] Redis는 휘발성으로 둔다: RDB·AOF 끔, 외부 포트 미게시, 내부 네트워크만. 메모리 상한을 두고 가득 차면 조용히 지우지 않고 기존 `unavailable`로 거부한다(`noeviction` 등).
- [ ] Redis 장애·재시작: 송수신은 기존 `unavailable`로 끝나고 Postgres로 대신 저장하지 않는다. 복구 뒤 새 송신은 정상 동작한다. `/healthz`의 처리는 dev가 정하고 문서에 적는다.
- [ ] 여러 relay 프로세스가 같은 Redis를 쓸 때 lease·claim·ACK·멱등이 지금의 보장과 같게 원자적으로 동작한다.
- [ ] 새 goose migration으로 `relay_state`에서 메시지·멱등·대화 사본·Dots Inbox를 지운다(D052). Down은 기존처럼 실패 고정이다. `make generate`·`make schema`를 갱신한다.
- [ ] `compose.yaml`과 `deploy/knowslink/compose.ops.yaml`에 Redis를 추가하고(운영 적용은 하지 않음), `scripts/verify_runtime.py`·`beta.sh` 확인 경로를 맞춘다. 운영 자원 한도(relay 0.5 CPU·256 MB, Postgres 1 CPU·512 MB)를 고려해 Redis 메모리 한도를 정하고 근거를 적는다.
- [ ] 문구를 바꾼다(한국어·영어): 개인정보·약관 페이지, 회원·대화 화면 안내, README, `docs/architecture.md`, `docs/system-design.md`, `docs/operations.md`. 뜻은 다음과 같다. "메시지는 서버에 영구 저장하지 않는다. 전달 대기는 최대 180초, 웹 대화 사본은 최대 1시간 서버 메모리(Redis)에만 두고 자동 삭제한다. 서버 재시작 때 그 전에 사라질 수 있다. 수신한 클라이언트가 내 기기에 보관한 사본은 그 기기에 남는다." 실제 구현과 다른 문장은 쓰지 않는다.
- [ ] 문구가 사실인지 확인한다: 로그·오류 메시지에 본문이 남지 않는지, 새 백업에 메시지가 없는지 확인한다. 이전 운영 백업에는 메시지가 남아 있으므로 그 보존 기간을 `prune-backups.sh` 기준으로 보고서에 적는다(삭제는 하지 않음).
- [ ] 영향을 받는 산출물 원천 문서를 갱신한다(아래 `갱신할 산출물`).

예상 규모: Go 600~1,200줄(시험 포함), 주요 파일은 `store.go`, `capacity.go`, `human_messages.go`, `public_text.go`, `remote_*.go`, 새 Redis 저장 파일, migration 1개. SIZE-002가 나면 이유를 적는다.

## 완료 기준과 검증

- Postgres에 메시지 본문·envelope·사본이 남지 않는다: 송수신·웹 대화·Dots 수신 뒤 `relay_state.state`와 모든 테이블에 해당 값이 없음을 시험으로 보인다.
- 대화 사본은 1시간 뒤, 전달 대기는 180초 뒤 사라진다(시계 주입 또는 Redis TTL 확인).
- Redis 재시작 뒤 이전 메시지는 없고, 재시작 중 송신은 `unavailable`, 복구 뒤 송신은 성공한다.
- 철회·친구 해제·키 철회 뒤 Redis에 남은 메시지가 전달되지 않는다.
- migration이 기존 행의 메시지·멱등·사본·Dots Inbox를 지우고, 회원·친구·agent·키·권한·Dots 연결 정보와 `human_model` guard는 그대로다.
- 개인정보 페이지와 안내 문구에 위 뜻의 한국어·영어 문장이 있다(페이지 시험).
- 검증 명령(Linux 컨테이너, 최종 코드 SHA): `make lint`, `make test`, `make build`, `go test -race -count=1 -tags=integration ./internal/... ./cmd/...`. 실제 Redis와 Postgres 컨테이너를 쓴다.
- FullOps lint: 체크아웃 루트에서 `python3 <플러그인>/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/dev --from 17acfa5e0b9800522854163b9d5df67b9870a159`, ERROR 0. Windows의 make 실행 불가는 Linux 결과로 대체한다.
- 보고서: `.fullops-squad/docs/evaluations/qa-reports/SAR-MSG-EPHEMERAL-001-test/report.md`(front matter 필수). 원시 로그는 `D:/workspace/KnowsLink/.devtools/msg-ephemeral/`.
- verifier: lite라 실행할 때만 지정한다. 독립 QA는 운영 이행 전 tester 후속으로 남긴다. 이 과제 완료는 제품 최종 수락이 아니다.
- Docker: 시험 자원은 `kl-msgeph-` 접두어로 만들고 완료 보고 전에 지운다(`AGENTS.md` 작업 파일 위치). 남기면 이름과 이유를 적는다. 시험 중 메모리 부족 알림이 오면 재시작하지 말고 진행 범위를 적어 `ask`로 묻는다.
- 멈추고 묻는 경우: 위 결정과 다른 보관 기한·범위가 필요할 때, 보안 보장을 유지할 수 없을 때, 설명되지 않는 검증 실패.

### UI 작업의 디자인 기준과 검증

해당 없음. 화면 구조는 바꾸지 않고 기존 페이지의 안내 문장만 바꾼다. 문장은 기존 페이지 형식을 따른다.

## 갱신할 산출물

- 갱신: D10, D07, D06 (route 추천). D07은 저장 구조, D06은 인터페이스·오류 의미, D10은 운영·보관 안내다.
- 확인 후 필요하면 갱신: D03, D05, D08, D09, D11, D12, D13. D08(테이블 정의)은 migration이 바뀌므로 `make schema` 결과를 반영한다.

## 기대 산출물

- 코드·migration·시험, 갱신한 문서와 문구, 보고서 `qa-reports/SAR-MSG-EPHEMERAL-001-test/report.md`, 운영 이행 절차 초안(D12, 운영 미적용 표시).

## 제약·협업·후속

- 범위 밖: 운영 배포·Redis 운영 기동, main 병합·push, 기존 백업 삭제, DB 대기 개선 4~6, Rust.
- 하위 작업은 `.fullops-squad/rules/delegation.md`의 standard 범위 안에서 정하고, 파일 소유권이 겹치지 않게 나눈다.
- 후속(남은 일에 적을 것): ops의 운영 이행(Redis 기동, migration, 기존 메시지 삭제)은 사용자 승인 뒤, tester 독립 QA, 이전 백업 보존 정책 판단.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad/docs/evaluations/jev/SAR-MSG-EPHEMERAL-001-packet.json` / SHA `17acfa5e0b9800522854163b9d5df67b9870a159` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 38, 43, 71, 81, 85 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 69 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read, document_update) · 줄 54 · inferred
- `.fullops-squad/docs/design-docs/database-design.md` (document_read, document_update) · 줄 16, 17, 47, 47, 65, 82 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/generated/db-schema.md` (document_read) · 줄 22, 28, 47, 56, 57, 58, 58, 59, 59, 60, 65, 66, 74, 74, 75 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 전체/미확인 · inferred
- `compose.yaml` (direct_edit) · 줄 전체/미확인 · unknown
- `deploy/knowslink/compose.ops.yaml` (direct_edit) · 줄 전체/미확인 · inferred
- `docs/development.md` (document_read) · 줄 9 · inferred
- `docs/system-design.md` (document_read) · 줄 72, 72, 74, 74, 76, 80, 99, 118, 125, 168, 215, 223, 245, 264, 264, 267, 281, 315, 317 · inferred
- `internal/relay/capacity.go` (direct_edit) · 줄 71, 77, 107, 297, 302 · inferred
- `internal/relay/human_messages.go` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/policy_pages.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/public_text.go` (direct_edit) · 줄 31 · inferred
- `internal/relay/remote_auth.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/remote_events.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/remote_mcp.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/store.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/text_events.go` (direct_edit) · 줄 전체/미확인 · inferred
나머지 참조 후보 26건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 10건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection; compose.yaml: missing, oversized, sensitive or invalid; manual inspection …. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 스스로 정한 선택과 되돌리는 방법 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/SAR-MSG-EPHEMERAL-001-test/report.md` 링크로 대신한다. worker_done 본문 첫 줄은 `[완료] SAR-MSG-EPHEMERAL-001 | 브랜치 fullops/dev | SHA <브랜치 HEAD>` 한 줄로 쓰고 SHA는 하나만 적는다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
