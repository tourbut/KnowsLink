---
title: SAR-PUBLISH-001-P03-DEPLOY — 개인정보처리방침 14일 백업 문장(main 4fdd149b)을 운영에 배포
status: draft
updated: 2026-10-22
owner: ops
tasks: [SAR-PUBLISH-001-P03-DEPLOY]
summary: 개인정보처리방침 14일 백업 문장(main 4fdd149b)을 운영에 배포
attempt: 84c1758b68e5417f84ae015cf802f036
base: 4fdd149b27d3ff251a73128120dce9f8bb136620
subagent_level: standard
test_level: lite
---

# SAR-PUBLISH-001-P03-DEPLOY — 개인정보처리방침 14일 백업 문장(main 4fdd149b)을 운영에 배포

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-22
- From / To: coor → ops
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 사용자 결정 D020·D021. 운영 relay를 main `4fdd149b` 제품 코드로 전환한다. 배포 전 과제 고유 backup과 rollback 태그(7c6ef373)를 만든다. env·ingress는 바꾸지 않는다. / 추가 승인 필요: env·ingress·DB 변경, Cloud Console, 실계정 사용, 백업 삭제, force-push
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/ops / fullops/ops
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

운영은 main 7c6ef373이다. 사용자 결정 D021에 따라 백업 14일 자동 정리가 설치됐고(SAR-PUBLISH-001-RETENTION), main 4fdd149b에 개인정보처리방침 KO/EN 백업 문장을 '만든 날부터 14일이 지나면 삭제'로 바꾼 변경(SAR-PUBLISH-001-P03, 정책 HTML 두 줄)이 있다. 배포 절차는 SAR-PUBLISH-001-OPS와 같다(`.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-OPS-test/report.md`).

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md`와 연결 규칙(testing.md `담당과 반복 범위`)이다. 기준 커밋 main `4fdd149b`. 정책 정본·완료 조건 PC-01–06·RC-01–05: `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`. 화면 정본 D04: `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`. 결정 D017·D018은 `decision.py list`. 산출물 ID D01~D13: `docs/deliverables/README.md`. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`(PC·RC 표와 KO/EN 문안)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md`(격리 fixture 절·화면 근거·검증하지 못한 것)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-OPS-test/report.md`(배포 절차 근거)
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`(정책 페이지·첫페이지 하단·설정 절)

## 해야 할 일과 파일 소유권

- [x] 배포 전 backup·rollback 태그, healthz 기록.
- [x] main `4fdd149b`로 relay를 전환하고 healthy·restart 0을 확인한다.
- [x] live `https://link.knowslog.com/privacy`(KO)와 영어판(`/lang/en?back=/privacy`)에 14일 문장이 보이는지, `/terms`·`/` 200인지 확인한다.
- [x] D12 `.fullops-squad/docs/operations/ops-guide.md` 현행 운영 SHA를 갱신한다.

소유 파일: D12, `qa-reports/SAR-PUBLISH-001-P03-DEPLOY-test/`, 이 인박스.

## 완료 기준과 검증

- 운영이 `4fdd149b` 제품 코드이고 healthy, live 정책에 14일 문장. 실패하면 7c6ef373으로 되돌리고 보고한다.
- 보고 `qa-reports/SAR-PUBLISH-001-P03-DEPLOY-test/report.md`(비밀 값 없음). 기록 커밋 뒤 `python3 C:/Users/shin/AppData/Roaming/orca/claude-profiles/6dd9013c-003c-4a1b-8dc5-f4ba45568e5e/home/plugins/cache/fullops-squad/fullops-squad/1.8.0/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/ops --from 4fdd149b` ERROR 0. `work.py finish`, `worker_done`.

### UI 작업의 디자인 기준과 검증

해당 없음(배포).

## 갱신할 산출물

D12 ops-guide.

## 기대 산출물

운영 반영, D12, 보고서.

## 제약·협업·후속

이전 RETENTION 과제의 기준 SHA·파일 소유권·멈출 조건은 이 지시서 기준으로 다시 정한다. Cloud Console은 coor가 사용자와 한다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-PUBLISH-001-P03-DEPLOY-packet.json` / SHA `4fdd149b27d3ff251a73128120dce9f8bb136620` / partial=True
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
- `.fullops-squad/docs/operations/ops-guide.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_read, document_update) · 줄 133 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
나머지 참조 후보 5건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-<역할>-test/report.md` 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

- 브랜치 / SHA: `fullops/ops` / 완료 SHA `3b97f80e30f70d3ced80fd2f3ee0761beb60a32d`(이 보고 기록 커밋은 그 뒤). 운영 제품은 main `4fdd149b27d3ff251a73128120dce9f8bb136620`이다.
- 변경 이유: D020·D021. 개인정보처리방침 KO/EN의 백업 문장을 ‘만든 날부터 14일이 지나면 삭제’로 바꾼 정책 HTML 두 줄을 운영 relay에 반영했다. env·ingress·DB는 바꾸지 않았다.
- 지시와 달라진 판단: 없다. 첫 smoke에서 `/lang/en` 리다이렉트 기대값을 `/privacy`로 잘못 둬 1건이 실패로 찍혔다(실제 `/privacy?lang=en`, 제품 결함 아님). 기대값을 고쳐 전체를 다시 실행해 7/7을 얻었다.
- 검토 필요: 새 백업 `7c6ef373-pre-p03deploy-20261021T235716Z.dump`는 2026-11-05 이후 14일 정리 대상이 된다. 더 둘 복구점은 따로 복사한다.
- 검증한 것: bundle hash·조상·infra diff 0, rollback 태그 `knowslink-rollback:7c6ef373-p03-20261022`, 새 0600 백업·TOC 25줄·격리 복원 digest 일치, relay 중단 약17초(23:57:12Z~23:57:29Z)·healthy·restart 0·migrate 0, `verify.py local` 통과, live `/privacy`(한국어)·영어판 14일 문장 있음·옛 문장 없음, `/terms`·`/` 200, 이웃 404, 상태 digest 보존(Rates만 volatile). 못 한 것: 브라우저 렌더링, 실계정 로그인·발송, 실제 rollback. 의존성 변경 없음(제품 코드·go.mod 변경 0).
- FullOps 검사: HEAD `3b97f80e`, 기준 `4fdd149b`, ERROR 0·WARNING 2(`PLANS.md`·`ops-guide.md` SIZE-001, 기존 초과 문서)·실행 불가 1(`make lint`, Windows `make` 부재, `make test`는 skipped). 제품 코드·테스트는 바꾸지 않았다.
- 산출물·로그: [보고서](../docs/evaluations/qa-reports/SAR-PUBLISH-001-P03-DEPLOY-test/report.md), D12 `ops-guide.md`(P03 운영 적용 절). 서버: `knowslink-state/p03deploy-20261022/`, `backups/7c6ef373-pre-p03deploy-20261021T235716Z.dump`.
- 남은 일: coor가 D13 `transition.md` '현재 운영 상태'의 운영 SHA(`7c6ef373`)를 `4fdd149b`로 고친다(ops 소유 밖이라 고치지 않았다). 2026-10-23 03:17 KST 이후 `backup-prune.log`에 첫 자동 실행 줄이 생겼는지 확인한다. Cloud Console은 coor와 사용자가 진행한다.
