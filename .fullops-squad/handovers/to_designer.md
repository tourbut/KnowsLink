---
title: SAR-SCHEDULE-001-DESIGN — 사용자 결정 D009·D010을 SAR-SCHEDULE-001에 기록하고 DEV 인계 기준을 확정한다
status: draft
updated: 2026-10-20
owner: designer
tasks: [SAR-SCHEDULE-001-DESIGN]
summary: 사용자 결정 D009·D010을 SAR-SCHEDULE-001에 기록하고 DEV 인계 기준을 확정한다
attempt: 223741f4de4b485fa8cad8a5906c5352
base: 96bf670a7bd38ce5e5856b68d4cafec7a56b35cb
subagent_level: standard
test_level: lite
---

# SAR-SCHEDULE-001-DESIGN — 사용자 결정 D009·D010을 SAR-SCHEDULE-001에 기록하고 DEV 인계 기준을 확정한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: SAR-SCHEDULE-001.md에 사용자 결정 기록과 DEV 인계 기준 확정, designer 브랜치 커밋. 제품 코드·외부 계정·운영 변경은 하지 않는다.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

같은 과제의 두 번째 attempt다. 초안 80b0960이 main 96bf670에 통합됐다. coor가 2026-10-20 사용자에게 물은 결과:

- D009: SD-01~SD-06 추천안(A) 전체 채택 — Google Calendar만, KnowsLink 서버가 토큰 암호화 보관·조회, 빈 시간 조회만(`schedule.commit`은 stub 유지), 바쁨 구간만 요청 단위 반올림 공개, 매 요청 사람 승인 = 그 요청 한 번의 조회와 한 번의 공개, 내 계정 '캘린더 연결' 별도 증분 동의·같은 화면 해제.
- D010: 수치 제안대로 승인 — 요청 시각부터 최대 14일, 15·30·60분(기본 30분), 친구당 24시간 5회, 승인 기한은 기존 gate 기한.
- D011(열림): Google Cloud Console OAuth 앱 게시 상태와 calendar.freebusy 민감 범위 심사 여부는 사용자가 나중에 확인한다. 개발·격리 QA는 진행하고 운영 연결 전에 다시 묻는다.

## 적용 기준과 예외

이전 attempt와 같다.

## 먼저 읽을 문서

SAR-SCHEDULE-001.md, 이전 완료 기록.

## 해야 할 일과 파일 소유권

- [ ] SAR-SCHEDULE-001.md에 결정 기록 절(D009·D010 원문 요지, 날짜, 결정자 user)을 넣고 SD-01~06 각 항목에 확정 표시, 수치 표를 확정값으로, front matter status를 확정 상태로 갱신한다.
- [ ] '추천안 기준 제품 규칙 초안'을 확정 규칙으로 정리하고 DEV 인계 기준(완료 조건 SC-01~09, 범위 밖, 운영 연결 전 D011 확인 필요)을 명시한다.
- [ ] D011이 열려 있음과 운영 영향(Testing이면 7일 만료, In production이면 심사 가능)을 '운영·심사 영향' 절에 남긴다.
- [ ] 결정 근거 report에 짧은 절 추가 → 완료 보고 → work.py finish → 커밋 → worker_done.

소유 파일: SAR-SCHEDULE-001.md, 결정 근거 report, own route/inbox/log.

## 완료 기준과 검증

DEV가 추측 없이 구현할 수 있게 결정과 수락 기준이 확정됐다. 문서 lint ERROR 0.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: 해당 없음(기획 기록).
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 해당 없음.
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 해당 없음.

## 갱신할 산출물

D02(SAR-SCHEDULE-001).

## 기대 산출물

확정된 SAR-SCHEDULE-001.md, report 절, 완료 보고.

## 제약·협업·후속

다른 역할 인박스·PLANS·board·GitHub·제품 코드는 수정하지 않는다. workspace 밖에 파일을 만들지 않는다.

## 이전 완료 기록

- .fullops-squad/handovers/logs/2026-10-20_to_designer.md

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-SCHEDULE-001-DESIGN-packet.json` / SHA `ceb42257523b8908c30b9e8ed61e9809190ea54a` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_update, document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 2, 6, 6, 7, 7, 14, 14, 14, 22, 31, 41, 45, 45, 50, 60, 65, 69 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 85 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 17 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/archive/designer-2026-10-14-landing.md` (document_read) · 줄 6, 76 · inferred
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 52 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 12, 14, 39, 41 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/deliverables/README.md` (document_read) · 줄 36, 36 · inferred
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` (document_read) · 줄 6, 16, 79 · inferred
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_read) · 줄 6, 14, 64 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 7 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 81 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 78 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 111 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md` (document_read) · 줄 92 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 340 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/COOR/dev-fix-2-final/coordinator-handover-supplement.md` (document_read) · 줄 33 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-REVIEW-review/report.md` (document_read) · 줄 21 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-DEV-test/report.md` (document_read) · 줄 63 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-DESIGN-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CODE-REVIEW-review/report.md` (document_read) · 줄 46 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CSS-REVIEW-review/report.md` (document_read) · 줄 39 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DESIGN-REVIEW-review/report.md` (document_read) · 줄 10, 34 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-FINAL-RECORDS-REVIEW-review/report.md` (document_read) · 줄 29 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-TESTER-test/css-followup/report.md` (document_read) · 줄 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-test/report.md` (document_read) · 줄 48 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW-review/report.md` (document_read) · 줄 101 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-2-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 23, 89 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-RECORD-REVIEW-review/report.md` (document_read) · 줄 33 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-RESULT-review/report.md` (document_read) · 줄 11 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-FINAL-REVIEW-review/report.md` (document_read) · 줄 11 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SCHEDULE-001-DESIGN/report.md` (document_read) · 줄 6, 6, 10, 10, 14, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-review/report.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-DESIGN/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-002-CODE-REVIEW-review/report.md` (document_read) · 줄 19 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-HOSTS-DESIGN/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.14.md` (document_read) · 줄 26, 27 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-FRONTEND-POLISH-001-DEV.md` (document_read) · 줄 26 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-DEV.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001-DEV.md` (document_read) · 줄 45 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md` (document_read) · 줄 63 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX.md` (document_read) · 줄 76 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 35, 80 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2.md` (document_read) · 줄 86 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-3.md` (document_read) · 줄 78 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX.md` (document_read) · 줄 69 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SESSION-PEER-001-DESIGN.md` (document_read) · 줄 2, 6, 10, 16 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV.md` (document_read) · 줄 39 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-DESIGN.md` (document_read) · 줄 6, 10, 38, 42 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 7 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 7, 7, 16, 16 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 7 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 7, 7, 91, 91 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/handovers/_TEMPLATE.md` (document_read) · 줄 53 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 74 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 206, 364 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_designer.md` (document_read) · 줄 46 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 32, 64 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_ops.md` (document_read) · 줄 79 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_tester.md` (document_read) · 줄 168 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_designer.md` (document_read) · 줄 142, 291 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_dev.md` (document_read) · 줄 207, 423, 469, 477 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_dev.md` (document_read) · 줄 104 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_tester.md` (document_read) · 줄 71, 147, 363 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_designer.md` (document_read) · 줄 9, 10, 13, 17, 25, 89, 100 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_ops.md` (document_read) · 줄 303, 317, 318, 396 · inferred
- `.fullops-squad/handovers/logs/2026-10-18_to_dev.md` (document_read) · 줄 156, 198 · inferred
- `.fullops-squad/handovers/logs/2026-10-19_to_designer.md` (document_read) · 줄 74, 137, 138, 158, 276, 340, 341, 362 · inferred
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 58 · inferred
- `.fullops-squad/handovers/to_ops.md` (document_read) · 줄 203, 227, 228, 229, 320 · inferred
- `.fullops-squad/lint/README.md` (document_read) · 줄 32, 82, 84, 95, 113 · inferred
- `.fullops-squad/review/_REPORT.md` (document_read) · 줄 20 · inferred
미확인 24건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
