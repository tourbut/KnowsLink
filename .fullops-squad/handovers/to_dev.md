---
title: "SAR-PUBLISH-001-P03 — 개인정보처리방침 백업 문장을 최대 14일 보관으로 수정(D021, P-03)"
status: draft
updated: 2026-10-22
owner: dev
tasks: [SAR-PUBLISH-001-P03]
summary: "개인정보처리방침 백업 문장을 최대 14일 보관으로 수정(D021, P-03)"
attempt: 392326d106924fd9a19217c8782f818d
base: d5695ff1b7c6c738f527a05978e69248d1404e16
subagent_level: standard
test_level: lite
---

# SAR-PUBLISH-001-P03 — 개인정보처리방침 백업 문장을 최대 14일 보관으로 수정(D021, P-03)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-22
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 사용자 결정 D021(2026-10-22). 개인정보처리방침 KO/EN의 백업 문장과 정본 SAR-PUBLISH-001.md의 해당 문장·P-03 행을 바꾸고 테스트를 맞춘다. / 추가 승인 필요: 다른 문안 변경, 운영 배포, force-push
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

운영은 main 7c6ef373(SAR-PUBLISH-001 배포 완료)다. ops 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-OPS-test/report.md` P-03 절: 정책 문안은 백업을 정기 교체·삭제한다고 하지만 실제로는 자동 정리가 없고 `knowslink-state/backups/`에 백업 20개(14일 초과 6개)가 쌓여 있다. 2026-10-22 사용자 결정 D021: 14일 초과 백업을 자동 삭제하고 기존 초과분도 지우며 문안을 '최대 14일 보관'으로 바꾼다.

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md`와 연결 규칙(testing.md `담당과 반복 범위`)이다. 기준 커밋 main `d5695ff1`. 정책 정본·완료 조건 PC-01–06·RC-01–05: `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`. 화면 정본 D04: `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`. 결정 D017·D018은 `decision.py list`. 산출물 ID D01~D13: `docs/deliverables/README.md`. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`(PC·RC 표와 KO/EN 문안)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md`(격리 fixture 절·화면 근거·검증하지 못한 것)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-OPS-test/report.md`(P-03·P-04 절)
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`(정책 페이지·첫페이지 하단·설정 절)

## 해야 할 일과 파일 소유권

- [ ] `internal/relay/policy/privacy.ko.html` 33행: '장애 복구용 백업에는 정보가 더 오래 남을 수 있습니다. 백업은 만든 날부터 14일이 지나면 삭제합니다.'
- [ ] `internal/relay/policy/privacy.en.html` 33행: 'Backups for recovering from failures can keep information longer. We delete backups 14 days after they are made.'
- [ ] `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` 111·199행의 `{백업 보존 기간}`·`{backup retention period}`를 14일 문장으로 바꾸고, 386행 P-03 행을 '확정: 14일(D021, ops 자동 정리 SAR-PUBLISH-001-RETENTION)'로 고친다.
- [ ] 정책 문안 일치 테스트가 있으면 함께 고친다. 시행일(P-01)은 바꾸지 않는다.

소유 파일: 위 두 HTML, 정본 문서, 관련 테스트, DEV 보고.

## 완료 기준과 검증

- 원격 Linux `make lint`·`make test` EXIT=0(dev `remote.py` 방식). `/privacy` KO/EN에 새 문장이 보이고 영어판에 한국어가 없다.
- 기록 커밋 뒤 `lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/dev --from d5695ff1` ERROR 0(경로는 ops 지시서와 같다). gate가 미추적 파일로 막히면 깨끗한 임시 detached worktree에서 돌리고 결과 경로를 보고한다. `work.py finish`, `worker_done`.

### UI 작업의 디자인 기준과 검증

문장 교체만 한다. 화면 구조·스타일은 바꾸지 않는다.

## 갱신할 산출물

D02(SAR-PUBLISH-001 P-03).

## 기대 산출물

HTML 두 개, 정본, 테스트, `qa-reports/SAR-PUBLISH-001-P03-test/report.md`.

## 제약·협업·후속

ops 백업 정리(SAR-PUBLISH-001-RETENTION)와 병렬. 운영 배포는 coor가 둘 다 끝난 뒤 배정한다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-PUBLISH-001-P03-packet.json` / SHA `d5695ff1b7c6c738f527a05978e69248d1404e16` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-OPS-test/report.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.4.md` (document_read) · 줄 44, 68 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-DEV.md` (document_read) · 줄 70 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-DEV.md` (document_read) · 줄 72 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001.md` (document_read) · 줄 114 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-001-DEV.md` (document_read) · 줄 101 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 77 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` (document_read) · 줄 67 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 184 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SESSION-PEER-001-DESIGN.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 133 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_update) · 줄 전체/미확인 · inferred
나머지 참조 후보 4건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-<역할>-test/report.md` 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
