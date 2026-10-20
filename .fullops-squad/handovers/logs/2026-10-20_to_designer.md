---
title: designer 완료 기록
status: draft
updated: 2026-10-20
owner: designer
summary: 지시서와 완료 보고를 보존한다.
---

## SAR-UX-ISSUES-002-APPROVAL-VISUAL — 2026-10-20
<!-- fullops-attempt: SAR-UX-ISSUES-002-APPROVAL-VISUAL 32e83db28e9b491c9737488a241f4305 -->

---
title: SAR-UX-ISSUES-002-APPROVAL-VISUAL — 승인 요약 후보 9d13a05의 직접 시각 검수
status: draft
updated: 2026-10-20
owner: designer
tasks: [SAR-UX-ISSUES-002-APPROVAL-VISUAL]
summary: 승인 요약 후보 9d13a05의 직접 시각 검수
attempt: 32e83db28e9b491c9737488a241f4305
base: 278001a45880f57134db4d2f86856497933b46bd
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-APPROVAL-VISUAL — 승인 요약 후보 9d13a05의 직접 시각 검수

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: fixture 허용 GET 관찰과 VISUAL 보고·필요한 제품 정본 보완, designer 브랜치 커밋. 제품 코드 수정은 하지 않는다.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

3단위(U5 승인 요약+U7) DEV 후보 `9d13a05ba93c7c085034213590ce661305606010`(fullops/dev, 제품 `f43d05c`, fixture 수정 `ea94128`). 원격 Linux make lint/test/integration EXIT=0, FullOps lint ERROR 0. 고정 SHA 독립 코드 리뷰가 병행 중이다(사용자 D003 '이번만 리뷰').

- DEV 보고(읽기 전용): `C:/Users/shin/orca/workspaces/KnowsLink/dev/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-DEV-test/report.md`(검증 필드↔문구 대응표, 렌더 `render/`, fixture 절, not_run).
- DEV 완료 기록: `C:/Users/shin/orca/workspaces/KnowsLink/dev/.fullops-squad/handovers/logs/`의 SAR-UX-ISSUES-002-APPROVAL-DEV 절.
- 제품 소스는 dev 워크트리에서 읽기 전용으로 확인한다.

## fixture

origin `http://127.0.0.1:18931`(coor SSH 포워딩, 2026-10-20 기동, fixture 수명 12시간). 서버 clone `/tmp/knowslink-approval002-9d13a05ba93c-a877fcc1`, container `kl-approval002-fixture-9d13a05ba93c-5cd13abf`, 고정 SHA 9d13a05. manifest(origin·accounts·gates·expires_at): `C:/Users/shin/AppData/Local/KnowsLinkDevTools/approval-ux002-dev/fixture-manifest.json`.

- 계정: 합성 A(승인자), B(요청자), C(gate 없음). `POST /fixture/as/{a|b|c}`(동일 출처만).
- gate 7개: `pending_query_hint`(악성 hint), `pending_commit_long`(긴 자유문·링크), `pending_decide_1`, `pending_decide_2`(승인·거절·중복 클릭용), `expired`, `processed`, `original_missing`. A와 B는 legacy pair라 요청자는 '확인할 수 없는 에이전트 / 구분값'으로 보인다(이름 있는 요청자는 DEV 렌더 fixture `gate-active`).
- C→A 친구 요청 1건 대기(O-NM-08 확인용).
- 허용 GET: `/`, `/home`, `/home/legacy`, `/home/friends`, `/home/account`, `/home/gates/*`. 허용 POST: `/home/gates/*`(승인·거절)와 `/fixture/as/*`. 그 밖은 403.
- tester와 designer가 함께 쓴다. 승인·거절 POST는 `pending_decide_1`·`pending_decide_2`에만, tester가 마지막에 실행한다.

## 적용 기준과 예외

fullops-common-0.3.3, `.fullops-squad/project.md`, 제품 정본 `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` U5·U7, D04. 기준 main 278001a. 등록 make lint/test는 DEV 같은 SHA 결과를 재사용한다(Windows make 없음). frontend.json codex host는 Codex 제외 지시로 적용하지 않는다.

## 먼저 읽을 문서

FULLOPS.md, common 세 규칙, project.md, 역할 contexts, SAR-UX-ISSUES-002.md U5·U7, D04, 위 DEV 보고.

## 해야 할 일과 파일 소유권

- [ ] 승인 목록(대기 있음/빈 목록), 상세(정상·악성 hint·긴 자유문과 링크·만료·처리됨·원문 부재), 익명 429 조회, 친구 요청 거절 안내 배치를 1440·390·320으로 직접 판정(PASS/FAIL/not_run, 근거 캡처). 이름 있는 요청자 화면은 DEV 렌더 `render/` 캡처를 쓴다.
- [ ] 읽기 순서·고정 문구·'요청자가 적은 내용' 분리·기술 상세 접기가 설계 정본과 맞는지, D04 토큰·44px·왼쪽 정렬·피할 스타일 준수.
- [ ] 캡처는 fixture 허용 GET으로 직접 만든다(headless 가능). 승인·거절 POST는 하지 않는다. tester 캡처(`C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-TESTER-test/`)가 생기면 참고할 수 있다.
- [ ] 판정 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-VISUAL/report.md`, 제품 결정이 필요하면 SAR-UX-ISSUES-002.md 반영 → 완료 보고 → work.py finish → 커밋 → worker_done(결론과 DEV 수정 필요 목록).

소유 파일: `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-VISUAL/` 하위, 필요 시 SAR-UX-ISSUES-002.md, own route/inbox/log.

## 완료 기준과 검증

시각 판정 결론(수락/조건부/거부)과 근거 캡처, DEV 수정 필요 항목. 스크린리더·실제 200%는 not_run 가능.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 해당 없음(제품 코드 수정 없음).
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 테마 전환 없음.

## 갱신할 산출물

없음(정본 결정 변경 시 SAR-UX-ISSUES-002.md만).

## 기대 산출물

보고서·결과·캡처, 완료 보고.

## 제약·협업·후속

제품 코드·다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다. 운영·실계정·실메일·외부 발송·삭제/이동·force push는 하지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-UX-ISSUES-002-APPROVAL-VISUAL-packet.json` / SHA `252531daceb674d9dc02f991cf88f5cc33f294c8` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_update, document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read, document_update) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 2, 6, 6, 6, 14, 14, 14, 22, 31, 31, 32, 32, 40, 40, 40, 40, 43, 43, 47, 51, 57, 57, 58, 58, 58, 60, 60, 60, 75 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/archive/designer-2026-10-14-landing.md` (document_read) · 줄 6, 72, 80, 84, 88, 92, 94, 96 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 6, 27, 29, 35, 37, 39, 41, 41, 43, 45 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 58, 60 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 6, 38, 46 · inferred
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` (document_read) · 줄 86, 86, 111 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 14, 20 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 7, 32, 48, 50, 52, 56 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-DEV-test/report.md` (document_read) · 줄 33, 48 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-TESTER-review/report.md` (document_read) · 줄 65 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-VISUAL-DOCS-review/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-VISUAL/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL-DOC-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL-FINAL-DOC-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/final-report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-FINAL-RECORDS-REVIEW-review/report.md` (document_read) · 줄 41 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-TESTER-test/css-followup/report.md` (document_read) · 줄 57 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/css-followup/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/report.md` (document_read) · 줄 6, 85 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-TESTER-review/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL-RESULT-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL/report.md` (document_read) · 줄 6, 10, 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-LIVE/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-2-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-review/report.md` (document_read) · 줄 6, 10, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md` (document_read) · 줄 6, 10, 179, 181 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL/report.md` (document_read) · 줄 6, 10, 77 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-VISUAL/report.md` (document_read) · 줄 6, 10, 72 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN/report.md` (document_read) · 줄 6, 10, 14, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 14, 61, 62, 63 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-DEV-test/report.md` (document_read) · 줄 2, 6, 10, 18, 94, 96, 96 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-VISUAL/report.md` (document_read) · 줄 6, 6, 10, 10, 18, 89, 89 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 15, 58, 59, 60 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-FIX-test/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-RETEST-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-TESTER-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-NAMES-TESTER.md` (document_read) · 줄 2, 6, 10, 12, 46 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-FRONTEND-POLISH-001-DEV.md` (document_read) · 줄 31 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-FRONTEND-POLISH-001-VISUAL.md` (document_read) · 줄 6, 10, 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-OPS.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-VISUAL.md` (document_read) · 줄 2, 6, 10, 24, 38 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001-VISUAL.md` (document_read) · 줄 6, 10, 16, 26, 28, 32, 47 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 104 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-DEV.md` (document_read) · 줄 15 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-VISUAL.md` (document_read) · 줄 2, 6, 10, 22 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-OPS.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-VISUAL.md` (document_read) · 줄 2, 6, 10, 25 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL.md` (document_read) · 줄 2, 6, 10, 21 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OPS.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-DESIGN.md` (document_read) · 줄 6, 10, 14, 14, 38, 42 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-NAMES-DEV.md` (document_read) · 줄 6, 7, 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-NAMES-TESTER.md` (document_read) · 줄 6, 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-RECEIPT-DEV.md` (document_read) · 줄 6, 20 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 7, 13, 15 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 7, 35, 45 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_dev.md` (document_read) · 줄 102, 112, 128, 129, 136, 137, 162, 183 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_ops.md` (document_read) · 줄 166, 170 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_tester.md` (document_read) · 줄 216, 221, 227 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_designer.md` (document_read) · 줄 136, 137, 138, 139, 144, 196, 197, 200, 204, 212, 249, 262, 283, 284, 285, 286, 287, 288, 294, 295, 297, 298, 299, 300, 301, 302, 304, 306, 311, 313, 315, 341 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_dev.md` (document_read) · 줄 211, 212, 249, 250, 251, 252, 258, 271, 288, 290, 467, 468, 471, 473, 479 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_dev.md` (document_read) · 줄 47, 136, 137, 140, 144, 152, 156, 185, 199 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_tester.md` (document_read) · 줄 203, 349 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_designer.md` (document_read) · 줄 69, 90 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_ops.md` (document_read) · 줄 123, 207, 237, 248, 257, 258, 259, 299, 300, 311, 312, 314, 317, 318, 319, 320, 321, 353, 356, 384, 396 · inferred
- `.fullops-squad/handovers/logs/2026-10-18_to_dev.md` (document_read) · 줄 158, 195, 198, 199 · inferred
- `.fullops-squad/handovers/logs/2026-10-19_to_designer.md` (document_read) · 줄 9, 9, 10, 10, 13, 13, 17, 17, 25, 25, 42, 43, 48, 52, 61, 61, 64, 64, 83, 92, 92, 94, 114, 115, 116, 117, 118, 119, 123, 124, 127, 128, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 145, 147, 148, 152, 154, 156, 158, 159, 203, 206, 206, 209, 209, 210, 210, 213, 213, 217, 217, 225, 225, 233, 244, 244, 254, 258, 263, 264, 266, 266, 281, 300, 300, 302, 320, 321, 322, 323, 324, 325, 328, 329, 331, 332, 334, 335, 336, 337, 338, 339, 340, 341, 342, 343, 343, 344, 345, 346, 347, 348, 350, 352, 353, 356, 358, 360, 362, 363, 392, 395, 395 · inferred
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 34, 34, 40, 42, 47, 61, 61, 69, 91, 91, 93, 115, 116, 117, 118, 119, 120, 120, 121, 122, 123, 124, 125, 126, 128, 129, 130, 131 · inferred
- `.fullops-squad/handovers/to_tester.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 31, 31, 32, 32, 40, 40, 40, 40, 43, 43, 47, 51, 55, 55, 59, 62, 62, 64, 64, 92, 92, 94, 196, 198, 199, 200, 201, 202, 203, 203, 204, 205, 206, 207, 208, 227, 264, 265, 266, 267 · inferred
미확인 21건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

- 브랜치·SHA: `fullops/designer`. 고정 SHA는 이 보고를 포함한 완료 커밋이며 worker_done에 적는다. 대상 후보는 `9d13a05`다.
- 결론: **수락**. 승인 목록(대기·빈 목록), 상세 6상태(정상·악성 hint·긴 자유문과 링크·만료·처리됨·원문 부재), 이름 있는 요청자(DEV 렌더), 익명 429(O-NM-04), 친구 거절 안내(O-NM-08)가 1440·390·320에서 PASS다. critical·high·medium은 없다. DEV 수정 필요 항목은 없다.
- 관찰: V-AP-01(low) 확인 불가 목록 행의 서버 정책 유형 표시는 현행 수락으로 결정했다. V-AP-02(low) 익명 429의 위치 nav ‘내 에이전트’ 링크는 후속 U8·U7 다듬기 후보다. V-AP-03(정보) 헤더 ‘설정’ 이름 변경은 범위 밖이다.
- 변경 이유: 지시서의 직접 시각 판정. 제품 코드·SAR-UX-ISSUES-002.md는 바꾸지 않았다(새 제품 규칙 없음).
- 검증한 것: fixture `127.0.0.1:18931`에 허용 GET과 `POST /fixture/as/{a|c}`만 보냈다. headless Chrome 캡처 33장과 목록 크롭 6장, 측정 JSON, `NO_OVERFLOW 33` 종료코드 0. 토큰 색만 사용, 그라디언트·그림자·4px 초과 모서리 0, 44px, 왼쪽 정렬, 대비를 확인했다.
- 못 한 것(not_run): 실제 키보드, 실제 200%, 스크린리더, 보조 브라우저, 승인·거절 POST 결과 화면(tester 몫).
- 주의: 첫 캡처 실행이 세션 쿠키 미저장으로 비로그인 GET을 반복해 fixture 요청 제한(429)이 약 1분 걸렸다(2026-10-20 08:54–08:57 KST). 재실행은 헤더 쿠키와 2.5초 간격으로 정상이다.
- 산출물: [보고](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-VISUAL/report.md), [캡처 스크립트](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-VISUAL/capture.mjs), `captures/`.
- 남은 일: coor의 tester QA·고정 SHA 리뷰 취합과 통합 판단.
