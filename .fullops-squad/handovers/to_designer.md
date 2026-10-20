---
title: SAR-SCHEDULE-001-DESIGN — 실제 일정 API 연동의 제품 범위 결정 선택지와 규칙 초안을 정리한다
status: draft
updated: 2026-10-20
owner: designer
tasks: [SAR-SCHEDULE-001-DESIGN]
summary: 실제 일정 API 연동의 제품 범위 결정 선택지와 규칙 초안을 정리한다
attempt: eaa1057e93ce4166a4cff3bd9b02d93a
base: 760dc1f4bad73df3edb03cd02f0e5cb7075d0341
subagent_level: standard
test_level: lite
---

# SAR-SCHEDULE-001-DESIGN — 실제 일정 API 연동의 제품 범위 결정 선택지와 규칙 초안을 정리한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 실제 일정 API 연동의 제품 기획(목표·사용자 흐름·범위·동의·실패 규칙·완료 조건)과 사용자 결정이 필요한 선택지 정리, 제품 정본 문서 작성, designer 브랜치 커밋. 외부 캘린더 계정 연결·실데이터 조회·API 키 발급·제품 코드 수정·운영 변경은 하지 않는다.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

사용자 결정 D008(2026-10-20)로 '실제 일정 API 연동'을 진행한다. 현재 제품은 에이전트끼리 `schedule.query`·`schedule.commit` 요청을 서명·전달하고 사람 승인 게이트를 거치지만, 승인 뒤에도 실제 캘린더 조회·일정 생성은 실행하지 않는다(정책상 disclosure·executable false). 보드의 '실제 어댑터·일정 조회' 단계가 이 후속을 남겨 두었다. 제품 범위(어떤 캘린더, 읽기만인지 쓰기까지인지, 동의·공개 범위)는 정해지지 않았다.

- 기존 정본: `.fullops-squad/docs/planning/product-specs/SAR-MVP.md`(첫 기능 범위·실제 calendar 제외 근거), `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md`, `.fullops-squad/docs/design-docs/architecture.md`(gate approve 뒤 executable/disclosure false), `interface-design.md`(gate 화면·요약), `SAR-UX-ISSUES-002.md` U5(승인 요약·'일정 실행을 허용하지 않습니다' 고정 문구).
- 로그인은 Google이다(운영 1dd6800). Google 연결 범위 기록: `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md`, `SAR-GOOGLE-CONNECT-002-DEV.md`.

## 적용 기준과 예외

fullops-common-0.3.3, 문서 규칙 `.fullops-squad/docs/agents/document-writing.md`, D01/D02 제품 정본. 기준 main(이 지시서 커밋). 외부 API의 실제 기능·권한 범위는 Context7 또는 공식 문서로 확인하고 출처를 남긴다. 확인하지 못한 API 동작을 확정 사실로 쓰지 않는다.

## 먼저 읽을 문서

위 정본들, `.fullops-squad/contexts/designer.md`, PLANS.md의 '실제 어댑터·일정 조회' 관련 절.

## 해야 할 일과 파일 소유권

- [ ] 사용자 결정이 필요한 선택지를 근거와 함께 정리: 대상 캘린더(예: Google Calendar만/그 밖), 기능 범위(빈 시간 조회만 / 일정 생성까지), 공개 범위(바쁨·한가함만 / 제목 등), 계정 연결·동의·철회 방식, 승인 게이트와의 관계(승인 = 그 요청 한 번의 실행 허용인지).
- [ ] 각 선택지의 위험·비용·운영 영향(외부 OAuth scope 심사 필요 여부, 토큰 보관, 실패·재시도, 감사 기록)을 짧게 비교하고 designer 추천안을 적는다.
- [ ] 추천안 기준의 제품 규칙 초안: 사용자 흐름, 화면 문구 원칙, 실패·만료·철회, 완료 조건(DEV/QA가 쓸 수락 기준), 범위 밖.
- [ ] 정본 `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md`(D02 성격) 작성, 결정 근거 `.fullops-squad/docs/evaluations/qa-reports/SAR-SCHEDULE-001-DESIGN/report.md`.
- [ ] 사용자 결정이 필요한 항목은 완료 보고에 질문 목록(선택지·추천)으로 남긴다. 결정 전까지 DEV 인계는 하지 않는다.
- [ ] 완료 보고 → work.py finish → 커밋 → worker_done.

소유 파일: 위 두 문서, 필요 시 D01/D02 인덱스 연결, own route/inbox/log.

## 완료 기준과 검증

coor가 사용자에게 바로 물을 수 있는 결정 목록(각 2–4 선택지·추천·영향)과 추천안 기준 제품 규칙 초안이 있다. 외부 API 근거 출처가 있다. 문서 lint ERROR 0.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04(화면 원칙만 참조).
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 해당 없음(기획).
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 해당 없음.

## 갱신할 산출물

D01, D02(새 제품 정본 연결).

## 기대 산출물

SAR-SCHEDULE-001.md, 결정 근거 report, 완료 보고(결정 질문 목록 포함).

## 제약·협업·후속

기술 구현 방법은 정하지 않는다(DEV 담당). 다른 역할 인박스·PLANS·board·GitHub·제품 코드는 수정하지 않는다. workspace 밖에 파일을 만들지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-SCHEDULE-001-DESIGN-packet.json` / SHA `4321ce8f8e5d8f11cc6ec419ac9f4783b517e1ba` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/architecture.md` (document_read, document_update) · 줄 47 · inferred · 필수
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_update, document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 47, 47, 61, 70 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 85 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 17 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/archive/designer-2026-10-14-landing.md` (document_read) · 줄 6, 76 · inferred
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 52 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 12, 14, 39, 41 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/deliverables/README.md` (document_read) · 줄 34 · inferred
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` (document_read) · 줄 6, 16, 79 · inferred
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_read) · 줄 6, 14, 64 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read, document_update) · 줄 40 · inferred
- `.fullops-squad/docs/design-docs/database-design.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read, document_update) · 줄 93 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read, document_update) · 줄 7, 31 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 81 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 78 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 35, 111 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md` (document_read) · 줄 21, 92 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_read, document_update) · 줄 131, 340 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read, document_update) · 줄 18, 31 · inferred
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
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-UI.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW-review/report.md` (document_read) · 줄 101 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-2-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 23, 89 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-RECORD-REVIEW-review/report.md` (document_read) · 줄 33 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-RESULT-review/report.md` (document_read) · 줄 11 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-FINAL-REVIEW-review/report.md` (document_read) · 줄 11 · inferred
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
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-PUBLIC-POLICY-001.md` (document_read) · 줄 24 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PREP-002.md` (document_read) · 줄 21 · inferred
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
- `.fullops-squad/docs/operations/user-guide.md` (document_read, document_update) · 줄 7 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 39, 75 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read, document_update) · 줄 7, 40, 56 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read, document_update) · 줄 7 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read, document_update) · 줄 25 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_update) · 줄 전체/미확인 · inferred
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
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

### 결과 — 2026-10-20

- 브랜치 / 문서 SHA: `fullops/designer` / `780bbbd`(정본·보고), 최종 SHA는 finish 커밋이며 worker_done에 적는다.
- 변경: [SAR-SCHEDULE-001](../docs/planning/product-specs/SAR-SCHEDULE-001.md) D02 초안(SD-01–06 결정표, 수치 제안, 추천안 기준 흐름·문구·실패/철회·기록·SC-01–09·범위 밖). [결정 근거](../docs/evaluations/qa-reports/SAR-SCHEDULE-001-DESIGN/report.md)(외부 API E1–E10 출처, 미확인 6건, 선택지 비교). D01·D02 인덱스(business-plan, deliverables README, SAR-MVP DEC-02 문단)에 연결 한 줄씩 추가.
- 외부 근거: Context7 monthly quota 초과. 하위 research 에이전트 1개가 Google·Microsoft·Apple 공식 문서를 WebFetch로 확인했다. `calendar.freebusy`의 sensitive 여부와 현재 OAuth 앱 게시 상태는 미확인이며 Cloud Console 확인이 필요하다.
- 검증: FullOps lint `--from d0fe090` head 780bbbd ERROR 0 / WARNING 0 / 실행 불가 1. product-lint는 Windows에 make가 없어 unavailable, product-test skipped. 제품 코드 변경 없음(`.fullops-squad` 밖 diff 0). 새 동작 테스트 없음(문서 과제).
- 사용자 결정 질문(첫 선택지가 추천):
  1. SD-01 대상 캘린더: A Google만 / B +Outlook / C +Outlook+iCloud.
  2. SD-02 조회 실행 위치: A KnowsLink 서버가 토큰 보관·조회 / B 받는 사람 agent의 자체 도구.
  3. SD-03 기능: A 빈 시간 조회만 / B +승인 뒤 내 캘린더에 생성 / C +상대 초대 발송.
  4. SD-04 공개: A 바쁨 구간 단위 반올림 / B 가능한 시간 후보 3개 / C 바쁨+제목.
  5. SD-05 승인 관계: A 매 요청 승인=그 요청 1회 조회·1회 공개 / B 친구별 사전 허용 시 자동 / C 친구 수락만으로 자동(비추천).
  6. SD-06 연결: A 내 계정 '캘린더 연결' 별도 증분 동의·같은 화면 해제 / B 로그인 때 함께 / C 다른 Google 계정 허용.
  7. 수치 제안 승인: 기간 최대 14일, 단위 15/30/60분(기본 30), 친구당 24시간 5회.
- 남은 일: 사용자 결정 기록 뒤 designer가 정본 확정·DEV 인계. coor/ops가 Cloud Console에서 scope 분류·앱 게시 상태 확인. 결정 전 DEV 인계·실조회 없음.
