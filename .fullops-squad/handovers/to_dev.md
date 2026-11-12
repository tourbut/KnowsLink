---
title: SAR-MSG-EPHEMERAL-001 — 독립 리뷰 medium 1(개인정보 문구의 승인 요청 metadata)과 low 1(interface-design 장애 범위 문장)을 사실대로 고친다
status: draft
updated: 2026-11-12
owner: dev
tasks: [SAR-MSG-EPHEMERAL-001]
summary: 독립 리뷰 medium 1(개인정보 문구의 승인 요청 metadata)과 low 1(interface-design 장애 범위 문장)을 사실대로 고친다
attempt: 5a794a6294d44300ab18e3cb686347c6
base: 9a5a07bd2239d10ff794149329270c411c4bdf23
subagent_level: standard
test_level: lite
---

# SAR-MSG-EPHEMERAL-001 — 독립 리뷰 medium 1(개인정보 문구의 승인 요청 metadata)과 low 1(interface-design 장애 범위 문장)을 사실대로 고친다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-12
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위: 아래 문구·문서·시험 기대값 수정과 `fullops/dev` 커밋.
- 추가 승인이 필요한 행위: 운영 접속·배포, main 병합·push, 승인 요청(Gate)을 Redis로 옮기는 코드 변경(사용자가 D054에서 고르지 않음).
- 상위 의도: 없음
- 담당 repo id / 워크트리 / 브랜치: `0b08ec4c-9e3a-4613-8197-5a835e545335` / `C:/Users/shin/orca/workspaces/KnowsLink/dev` / `fullops/dev`
- 병합 책임자 / 기본 브랜치: coor / `main`
- 복귀 repo id / 워크트리 / 터미널 / run·task·dispatch id: `0b08ec4c-9e3a-4613-8197-5a835e545335` / `C:/Users/shin/orca/workspaces/KnowsLink/coor` / Orca preamble의 값

## 현재 상황과 확인 근거

- 첫 시도는 `9a5a07bd`(코드 `6fe42d32`)로 끝났다. 이전 지시서와 완료 보고는 `.fullops-squad/handovers/logs/2026-11-12_to_dev.md`에 있다.
- 독립 리뷰 `SAR-MSG-EPHEMERAL-001-REVIEW`는 조건부 수락했다(critical/high 0). 보고서는 `.fullops-squad/docs/evaluations/qa-reports/SAR-MSG-EPHEMERAL-001-REVIEW-review/report.md`, 지적은 같은 폴더 `result.json`의 findings다.
- 이번에 고칠 지적:
  - medium 1(`internal/relay/policy/privacy.ko.html`·`privacy.en.html` 29줄): 개인정보처리방침이 '승인 요청'을 Redis에만 둔다고 적는다. 실제로는 본문 없는 Gate metadata(ID·Owner·Parent·Digest·From·To·Policy·State·Exp)가 Postgres `relay_state`에 요청 사본이 사라질 때까지(최대 약 1시간) 남는다. `pg_dump` 백업에는 최대 14일 남는다. Dots 구독의 Seen/Pending/EventID 메시지 ID도 Postgres에 있다.
  - low 1(`.fullops-squad/docs/design-docs/interface-design.md` 16줄, dev 보고서의 같은 문장): 'Redis에 닿지 못하면 모든 relay transaction이 unavailable'은 행에 Ephemeral 포인터가 있을 때만 맞다. 메시지 기록이 없거나 쓰지 않는 transaction은 Redis 장애 중에도 commit된다.
- 사용자 결정 D054(FullOps 결정 기록): 코드는 그대로 두고 문구를 사실대로 고친다.

## 적용 기준과 예외

- 공통 기준: `.fullops-squad/rules/common/README.md`와 연결된 규칙, `AGENTS.md`, `.fullops-squad/project.md`, 기준 `9a5a07bd2239d10ff794149329270c411c4bdf23`.
- 산출물 ID `D01~D13`은 `.fullops-squad/docs/deliverables/README.md`를 따른다.
- 예외 없음. 보관 동작은 바꾸지 않는다.

## 먼저 읽을 문서

- `.fullops-squad/docs/evaluations/qa-reports/SAR-MSG-EPHEMERAL-001-REVIEW-review/report.md`, `result.json`
- `internal/relay/policy/privacy.ko.html`, `internal/relay/policy/privacy.en.html`, `internal/relay/policy_pages_test.go`
- `.fullops-squad/docs/design-docs/interface-design.md`, `.fullops-squad/docs/design-docs/database-design.md`, `.fullops-squad/docs/evaluations/qa-reports/SAR-MSG-EPHEMERAL-001-test/report.md`
- 같은 내용을 담은 다른 안내(README, `docs/architecture.md`, `docs/system-design.md`, 회원 화면 i18n)는 검색해 같은 오류가 있으면 함께 고친다.

## 해야 할 일과 파일 소유권

파일 소유권: 위 문서·정책 HTML·i18n·시험 파일, `qa-reports/SAR-MSG-EPHEMERAL-001-test/report.md`(기존 보고서에 절을 덧붙임).

- [ ] 개인정보처리방침 ko·en: 메시지 본문·envelope·사본은 영구 저장하지 않고 Redis에만 둔다는 문장은 유지한다. 본문 없는 승인 요청·전달 metadata(ID·내용 해시·대상·상태)는 Postgres에 최대 약 1시간, 백업에 최대 14일 남는다는 문장을 따로 적는다. 실제 코드 동작과 숫자를 다시 확인해 적는다.
- [ ] `policy_pages_test.go` 기대값을 맞추고, 같은 오류가 있는 다른 안내 문장을 함께 고친다.
- [ ] interface-design과 dev 보고서의 Redis 장애 범위 문장을 '메시지 기록이 있거나 새로 쓰는 transaction'으로 고친다.
- [ ] 기존 보고서에 수정 절(무엇을, 어떤 근거로, 검증 결과)을 덧붙인다.

예상 규모: 문서·HTML 20~60줄, 시험 10줄 안팎.

## 완료 기준과 검증

- 개인정보처리방침 ko·en이 Gate·전달 metadata의 Postgres·백업 보관을 사실대로 적고, 메시지 본문에 대한 '영구 저장하지 않음'은 그대로 남는다.
- `go test ./internal/relay/ -run 'Policy|Landing|Human|Member'`가 통과한다. 변경이 시험과 문서뿐이면 Linux 전체 gate는 다시 돌리지 않아도 된다. 다른 Go 코드를 바꾸면 Linux `make lint`·`make test`를 최종 SHA에서 실행한다.
- FullOps lint: `python3 <플러그인>/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/dev --from 9a5a07bd2239d10ff794149329270c411c4bdf23`, ERROR 0.
- verifier: lite라 실행할 때만 지정한다. 이 과제 완료는 제품 최종 수락이 아니다.
- Docker를 쓰면 `kl-msgeph-` 접두어로 만들고 완료 보고 전에 지운다.

### UI 작업의 디자인 기준과 검증

해당 없음. 정책 페이지의 문장만 바꾸고 화면 구조는 그대로다.

## 갱신할 산출물

- D07(저장 구조 문장 확인), D06(interface-design 장애 범위 문장), D10(운영 안내 문구 확인). 바꾸지 않은 ID는 "확인만"으로 적는다.

## 기대 산출물

- 수정한 정책 HTML ko·en, 시험, 문서, 보고서 추가 절.

## 제약·협업·후속

- 범위 밖: 리뷰 medium 2(Redis 한도 도달 시 축소 불가)는 coor가 PLANS.md 후속으로 등재한다. 이번 과제에서 고치지 않는다.
- 범위 밖: 운영 배포, main 병합·push, Gate 저장 위치 변경.

## 이전 완료 기록

- .fullops-squad/handovers/logs/2026-11-12_to_dev.md

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad/docs/evaluations/jev/SAR-MSG-EPHEMERAL-001-packet.json` / SHA `b81218841a6b14bdc5ef8e2c0521667978cb561c` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MSG-EPHEMERAL-001-REVIEW-review/report.md` (document_read) · 줄 39, 45, 52, 82 · inferred · 필수
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 7, 14, 34, 34, 35, 48, 57, 76 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/ops.md` (document_read) · 줄 81, 87, 89 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 69 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read, document_update) · 줄 61 · inferred
- `.fullops-squad/docs/design-docs/database-design.md` (document_read, document_update) · 줄 18, 63, 81, 98 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_update, document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-OPS.md` (document_read) · 줄 25, 33, 36 · inferred
- `.fullops-squad/docs/generated/db-schema.md` (document_read) · 줄 22, 28, 47, 56, 57, 58, 59, 60, 65, 66, 74, 83, 89, 90 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 전체/미확인 · inferred
- `db/migrations/00003_relay_state_json.sql` (impact_check) · 줄 4, 5, 6, 7, 8, 13, 14, 22 · inferred
- `docs/README.md` (document_read) · 줄 35, 35 · inferred
- `docs/architecture.md` (document_read) · 줄 96, 96, 120 · inferred
- `docs/system-design.md` (document_read) · 줄 262, 262, 266, 339, 339 · inferred
- `internal/database/db.go` (impact_check) · 줄 14, 15, 16, 17 · inferred
- `internal/relay/ephemeral_integration_test.go` (impact_check) · 줄 67, 90, 98, 104, 128 · inferred
- `internal/relay/policy/privacy.en.html` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/policy/privacy.ko.html` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/policy_pages_test.go` (direct_edit) · 줄 14, 31, 80, 113 · inferred
- `internal/relay/state_migration_integration_test.go` (impact_check) · 줄 98, 110, 124, 145, 146, 184, 187, 200, 211 · inferred
나머지 참조 후보 29건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 5건: route: task/role/SHA/attempt/instruction differs or absent; regenerate/bind route; .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 스스로 정한 선택과 되돌리는 방법 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/SAR-MSG-EPHEMERAL-001-test/report.md` 링크로 대신한다. worker_done 본문 첫 줄은 `[완료] SAR-MSG-EPHEMERAL-001 | 브랜치 fullops/dev | SHA <브랜치 HEAD>` 한 줄로 쓰고 SHA는 하나만 적는다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
