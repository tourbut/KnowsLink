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

## SAR-UX-ISSUES-002-HOSTS-DESIGN — 2026-10-20
<!-- fullops-attempt: SAR-UX-ISSUES-002-HOSTS-DESIGN b5c7601c40374611a78963fdb76dab25 -->

---
title: SAR-UX-ISSUES-002-HOSTS-DESIGN — U8 호스트 지원 표를 최신 사용자 범위와 Claude Code 세션 연결로 갱신하고 O-NM-03을 결정한다
status: draft
updated: 2026-10-20
owner: designer
tasks: [SAR-UX-ISSUES-002-HOSTS-DESIGN]
summary: U8 호스트 지원 표를 최신 사용자 범위와 Claude Code 세션 연결로 갱신하고 O-NM-03을 결정한다
attempt: b5c7601c40374611a78963fdb76dab25
base: 79864592804e595810706f6550ff255ee65407ac
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-HOSTS-DESIGN — U8 호스트 지원 표를 최신 사용자 범위와 Claude Code 세션 연결로 갱신하고 O-NM-03을 결정한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 제품 정본(SAR-UX-ISSUES-002.md U8·U4 관련 절) 갱신과 근거 기록, designer 브랜치 커밋. 호스트 설치·실계정 연결·외부 발송·운영 변경·제품 코드 수정은 하지 않는다.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

잔여 UI 1–3단위가 main 7986459에 통합됐다. 다음은 4단위 U8(회원 연결 경로·재개, 나머지 U7)이다. 그런데 U8 '근거에 따른 호스트 표시' 표(2026-10-15)가 이후 사실과 맞지 않는다.

- 사용자 호스트 범위(PLANS.md 2026-10-15 '최신 사용자 호스트 범위' 줄): Meta Muse는 대상에서 제외, OpenClaw는 사용자 지시로 보류, Hermes는 사용자가 지정한 서버의 Nous Research Hermes Agent 0.21.3을 확인(설치·왕복 수락 아님).
- 2026-10-17 Claude Code 세션 연결이 구현·수락·main 통합됐다: `/knowslink-session on|off|status [--take]`, 도착 toast, 승인 대화상자(SAR-SESSION-PEER-001/002, 정본 `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md`, 기술 `.fullops-squad/docs/design-docs/module-design.md`·`architecture.md`, QA `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-00*-*`). Codex는 V02 미실행.
- 운영 배포 1dd6800(2026-10-17)에 세션 API가 반영됐다(PLANS '운영 배포 완료·통합').
- 2단위 관찰 O-NM-03(후속 U4 개선 후보): 보낸 친구 요청 행이 서로 같아 요청자가 어느 코드로 보냈는지 알 수 없다. 이전 판정은 '요청자가 입력한 코드의 일부(예: 끝 4자) 표시 검토'(`.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-VISUAL/report.md` 'O-NM-03~07 분류').

## 적용 기준과 예외

fullops-common-0.3.3, 문서 규칙, 제품 정본 SAR-UX-ISSUES-002.md(U8·U4), D02 SAR-HUMAN-PEER.md, D04. 기준 main 7986459. 근거가 없는 호스트를 지원 완료로 올리지 않는 기존 원칙을 유지한다.

## 먼저 읽을 문서

SAR-UX-ISSUES-002.md U4·U8, SAR-SESSION-PEER-001.md, PLANS.md의 2026-10-15 호스트 범위 줄과 2026-10-17 세션 통합 줄, README.md·adapters/README.md의 현재 설치 안내, `.fullops-squad/contexts/designer.md`.

## 해야 할 일과 파일 소유권

- [x] U8 호스트 표 갱신: 대상 호스트 목록(Meta Muse 제외 명시, OpenClaw 보류 표시 방식, Hermes 상태, Claude Code 세션 연결 추가)과 각 근거 날짜·SHA·화면의 지원 상태 문구·다음 행동. 소스/패키지 근거와 실제 앱 설치·왕복 근거를 구분한다. Codex 세션은 미검증으로 둔다.
- [x] '어디에서 사용할까요?' 선택지 목록을 위 표와 맞춘다(웹·Grok·Dots·Claude Code·지원 확인 필요).
- [x] O-NM-03 결정: 보낸 요청 행에 요청자 자기 입력 코드의 일부를 보일지, 보인다면 형식(자리 수·마스킹)과 상대 존재 비노출 경계. 보이지 않기로 하면 근거.
- [x] 변경을 SAR-UX-ISSUES-002.md에 반영(front matter 갱신), 결정 근거를 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-HOSTS-DESIGN/report.md`에 짧게 기록.
- [x] 완료 보고 → work.py finish → 커밋 → worker_done(DEV 4단위 구현 기준 요약 포함).

소유 파일: SAR-UX-ISSUES-002.md, 위 report, own route/inbox/log.

## 완료 기준과 검증

DEV가 추측 없이 4단위를 구현할 수 있을 만큼 호스트 표·선택지·O-NM-03이 확정됐다. 문서 lint ERROR 0.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04, SAR-UX-ISSUES-002.md.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 해당 없음(제품 기준 문서).
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 해당 없음.

## 갱신할 산출물

SAR-UX-ISSUES-002.md(필요 시 D02).

## 기대 산출물

갱신된 정본, 결정 근거 report, 완료 보고.

## 제약·협업·후속

기술 구현 방법은 정하지 않는다(DEV 담당). 다른 역할 인박스·PLANS·board·GitHub·제품 코드는 수정하지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-UX-ISSUES-002-HOSTS-DESIGN-packet.json` / SHA `1df65d757aaac177dcad548cb1abfedc075cbff8` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` (document_read) · 줄 6, 14, 14, 16, 16, 16, 24, 24, 56, 56, 79, 79, 79, 83, 83, 86, 86, 86, 99, 100, 111, 113 · inferred · 필수
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_update, document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 2, 2, 6, 6, 6, 7, 7, 14, 14, 14, 14, 22, 32, 32, 34, 38, 42, 46, 46, 47, 47, 49, 49, 49, 52, 60, 62, 67 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 16, 16, 54, 54, 85 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 17 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/archive/designer-2026-10-14-landing.md` (document_read) · 줄 6, 76 · inferred
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 52 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 6, 12, 13, 14, 39, 39, 41, 41, 41 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 23, 58, 60, 62 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 6, 38, 44, 46, 50 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read, document_update) · 줄 201, 201, 203, 205, 205 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/database-design.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read, document_update) · 줄 7, 14, 20 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read, document_update) · 줄 7 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 81 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 78 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 111 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md` (document_read) · 줄 92 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_read, document_update) · 줄 7, 14, 14, 21, 21, 24, 28, 28, 32, 50, 52, 54, 58, 322 · inferred
- `.fullops-squad/docs/design-docs/motion/SAR-MOTION-INTRO-001-brief.md` (document_read) · 줄 27, 27 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read, document_update) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/COOR/dev-fix-2-final/coordinator-handover-supplement.md` (document_read) · 줄 33 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.3.0-FINAL-review/report.md` (document_read) · 줄 26, 38, 47 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-REVIEW-review/report.md` (document_read) · 줄 21 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-FINAL-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-N1-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-HOOK-QUOTE-001-REVIEW-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-DEPLOY-001-OPS-FINAL-review/report.md` (document_read) · 줄 12, 12, 22, 32, 45 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-DEV-test/report.md` (document_read) · 줄 63 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-DEV-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GROK-WEBHOOK-001-TESTER-REVIEW-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-DESIGN-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-LIVE/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CODE-REVIEW-review/report.md` (document_read) · 줄 46 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CSS-REVIEW-review/report.md` (document_read) · 줄 39 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DESIGN-REVIEW-review/report.md` (document_read) · 줄 10, 34 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-FINAL-RECORDS-REVIEW-review/report.md` (document_read) · 줄 29 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-TESTER-test/css-followup/report.md` (document_read) · 줄 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-test/report.md` (document_read) · 줄 48 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MOTION-INTRO-001-DEV-test/a2/report.md` (document_read) · 줄 37, 37 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FINAL-review/report.md` (document_read) · 줄 27 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FIX-review/report.md` (document_read) · 줄 28 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-FIX-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-DEV-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 12, 12, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW-review/report.md` (document_read) · 줄 12, 12, 25, 26, 101 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-REVIEW-review/report.md` (document_read) · 줄 12, 12, 22, 23, 23 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-2-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 12, 12, 23, 89 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-RECORD-REVIEW-review/report.md` (document_read) · 줄 33 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-RESULT-review/report.md` (document_read) · 줄 11 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-FINAL-REVIEW-review/report.md` (document_read) · 줄 11 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-RELEASE-001-OPS/report.md` (document_read) · 줄 77, 77, 104, 104 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-TESTER-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL/report.md` (document_read) · 줄 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-OPS-RESULT-review/report.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-review/report.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-CODE-REVIEW-2-review/report.md` (document_read) · 줄 11, 11, 61, 70, 70 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-CODE-REVIEW-3-review/report.md` (document_read) · 줄 11, 11, 27, 27, 75, 75 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-CODE-REVIEW-review/report.md` (document_read) · 줄 11, 11, 21, 21, 38, 45, 64, 64 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-DESIGN/report.md` (document_read) · 줄 2, 6, 10, 22, 33 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-DEV-test/report.md` (document_read) · 줄 2, 7, 7, 18, 18, 70, 70, 79, 91, 91 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-RETEST-test/report.md` (document_read) · 줄 19, 19, 44, 44, 60, 92 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-TESTER-test/report.md` (document_read) · 줄 2, 2, 7, 7, 19, 19, 25, 31, 110, 110, 126 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-002-CODE-REVIEW-review/report.md` (document_read) · 줄 19 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-002-DEV-test/report.md` (document_read) · 줄 2, 2, 7, 7, 18, 18, 20, 33, 51, 51 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-002-TESTER-test/report.md` (document_read) · 줄 19, 19, 140, 140, 153 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-DEV-099-review/report.md` (document_read) · 줄 12, 12, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-INTEGRATION-FINAL-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-INTEGRATION-SOURCE-review/report.md` (document_read) · 줄 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-review/report.md` (document_read) · 줄 12, 12, 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-REVIEW-review/report.md` (document_read) · 줄 46 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-TESTER-test/report.md` (document_read) · 줄 27 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 11, 11, 15, 21, 42, 59, 65 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-DEV-test/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-TESTER-test/report.md` (document_read) · 줄 6, 10, 25 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-VISUAL/report.md` (document_read) · 줄 6, 10, 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6, 6, 10, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN/report.md` (document_read) · 줄 6, 6, 10, 10, 14, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 11, 11, 14, 61, 62, 63 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-DEV-test/report.md` (document_read) · 줄 2, 6, 10, 18, 96 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-VISUAL/report.md` (document_read) · 줄 6, 10, 18, 89, 89 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 11, 11, 15, 58, 59, 60 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-FIX-test/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-RETEST-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-TESTER-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-WINDOWS-SHELL-GUARD-001-REVIEW-review/report.md` (document_read) · 줄 46, 49, 53 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-WINDOWS-SHELL-GUARD-001-TESTER-test/report.md` (document_read) · 줄 32, 33, 34, 67, 67 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-APPROVAL-TESTER.md` (document_read) · 줄 2, 6, 10, 12 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-NAMES-TESTER.md` (document_read) · 줄 2, 6, 10, 12, 46 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.10.md` (document_read) · 줄 39, 39 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.11.md` (document_read) · 줄 33, 33 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.14.md` (document_read) · 줄 26, 27 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-099.md` (document_read) · 줄 18, 18, 19 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.3.0.md` (document_read) · 줄 28, 46 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-FRONTEND-POLISH-001-DEV.md` (document_read) · 줄 26 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-DEV.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001-DEV.md` (document_read) · 줄 45 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-DEV-REVIEW.md` (document_read) · 줄 15, 15 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md` (document_read) · 줄 63 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX.md` (document_read) · 줄 76 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 35, 80 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW.md` (document_read) · 줄 18, 18 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW.md` (document_read) · 줄 18, 18 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2.md` (document_read) · 줄 86 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-3.md` (document_read) · 줄 78, 87 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX.md` (document_read) · 줄 69 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-VISUAL.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-VISUAL.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SESSION-PEER-001-DESIGN.md` (document_read) · 줄 2, 6, 10, 16, 16 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SESSION-PEER-001-DEV.md` (document_read) · 줄 2, 2, 7, 7, 16, 16, 22, 39, 39, 52 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-DEV.md` (document_read) · 줄 82 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-INSTALL.md` (document_read) · 줄 18, 18, 24, 37, 38, 39, 50, 50, 51, 61, 61, 82 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV.md` (document_read) · 줄 39 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-DESIGN.md` (document_read) · 줄 6, 6, 10, 10, 14, 14, 38, 38, 42, 42 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-NAMES-DEV.md` (document_read) · 줄 6, 7, 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-NAMES-TESTER.md` (document_read) · 줄 6, 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-RECEIPT-DEV.md` (document_read) · 줄 6, 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-WINDOWS-SHELL-GUARD-001-DEV.md` (document_read) · 줄 7, 19, 22, 44, 44, 45, 54, 57, 57, 66, 66, 69, 69 · inferred
- `.fullops-squad/docs/generated/db-schema.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_read, document_update) · 줄 7, 13, 15, 17, 19, 77 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read, document_update) · 줄 7, 7, 17, 17, 22, 22, 24, 24, 26, 35, 45, 107, 107 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 52, 52 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read, document_update) · 줄 7, 16, 16, 38, 38 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read, document_update) · 줄 7, 15, 15 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read, document_update) · 줄 72, 72 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read, document_update) · 줄 39, 39 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/README.md` (document_read) · 줄 44, 44 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md` (document_read) · 줄 145, 145 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/product.md` (document_read) · 줄 133, 133 · inferred
- `.fullops-squad/handovers/_TEMPLATE.md` (document_read) · 줄 53 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_ops.md` (document_read) · 줄 72, 109, 119 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_ops.md` (document_read) · 줄 179, 179, 211, 211 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_dev.md` (document_read) · 줄 94, 113, 119, 139, 139, 152 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 74 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 53, 69, 69, 125, 125, 191, 204, 204, 206, 275, 275, 326, 326, 364, 374, 374 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_designer.md` (document_read) · 줄 46 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 32, 64 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_ops.md` (document_read) · 줄 79 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_tester.md` (document_read) · 줄 168 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_designer.md` (document_read) · 줄 142, 291 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_dev.md` (document_read) · 줄 207, 423, 469, 477 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_dev.md` (document_read) · 줄 104, 136, 137, 140, 144, 152, 156, 185, 199 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_tester.md` (document_read) · 줄 71, 147, 363 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_designer.md` (document_read) · 줄 9, 10, 13, 17, 25, 41, 47, 58, 58, 60, 69, 85, 89, 90, 91, 91, 100, 102 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_ops.md` (document_read) · 줄 123, 123, 123, 207, 303, 317, 317, 318, 318, 319, 320, 321, 396, 396, 499, 499 · inferred
- `.fullops-squad/handovers/logs/2026-10-18_to_dev.md` (document_read) · 줄 156, 198, 198, 199 · inferred
- `.fullops-squad/handovers/logs/2026-10-19_to_designer.md` (document_read) · 줄 9, 10, 13, 17, 25, 42, 43, 48, 52, 61, 64, 74, 92, 94, 137, 137, 138, 138, 139, 140, 141, 142, 143, 158, 158, 159, 203, 206, 209, 210, 213, 217, 225, 244, 254, 258, 263, 266, 276, 281, 300, 302, 340, 340, 341, 341, 342, 343, 344, 345, 346, 347, 348, 362, 362, 363, 392, 395 · inferred
- `.fullops-squad/handovers/logs/2026-10-20_to_designer.md` (document_read) · 줄 9, 10, 13, 17, 25, 42, 43, 58, 62, 68, 69, 71, 81, 86, 99, 101, 139, 139, 140, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 164, 164, 165, 166, 167, 194, 198 · inferred
- `.fullops-squad/handovers/logs/2026-10-20_to_tester.md` (document_read) · 줄 9, 10, 13, 17, 25, 42, 43, 58, 62, 66, 73, 75, 85, 103, 105, 209, 209, 210, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 238, 275, 275, 276, 277, 278, 327 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-INTEGRATION-REVIEW.md` (document_read) · 줄 24 · inferred
- `.fullops-squad/lint/README.md` (document_read) · 줄 32, 82, 84, 95, 113 · inferred
- `.fullops-squad/orca-agents.md` (document_read) · 줄 58, 99, 113 · inferred
- `.fullops-squad/review/_REPORT.md` (document_read) · 줄 20 · inferred
- `AGENTS.md` (document_read) · 줄 13, 14, 14 · inferred
- `adapters/README.md` (document_read) · 줄 33, 33, 34, 38, 48, 50, 50, 52, 52, 58, 58, 75 · inferred
- `adapters/hooks/register.test.ts` (impact_check) · 줄 95, 95 · inferred
- `adapters/hooks/register.ts` (impact_check) · 줄 1, 1, 21, 21 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 14, 34, 34, 36, 36 · inferred
- `cmd/relay/main_test.go` (impact_check) · 줄 28, 29 · inferred
- `internal/database/db.go` (impact_check) · 줄 1 · inferred
- `internal/database/models.go` (impact_check) · 줄 1 · inferred
- `internal/database/relay.sql.go` (impact_check) · 줄 1 · inferred
- `internal/relay/admission_integration_test.go` (impact_check) · 줄 135, 162, 163 · inferred
- `internal/relay/cleanup_flood_integration_test.go` (impact_check) · 줄 32, 207, 208 · inferred
- `internal/relay/cleanup_unit_integration_test.go` (impact_check) · 줄 46, 79 · inferred
- `internal/relay/connections_test.go` (impact_check) · 줄 314, 315 · inferred
- `internal/relay/human_http.go` (impact_check) · 줄 96 · inferred
- `internal/relay/human_pages.go` (impact_check) · 줄 24, 35 · inferred
- `internal/relay/human_pages_test.go` (impact_check) · 줄 31 · inferred
- `internal/relay/landing_test.go` (impact_check) · 줄 25, 26 · inferred
- `internal/relay/public_messages_integration_test.go` (impact_check) · 줄 171, 317, 368, 434 · inferred
- `internal/relay/test_messages_integration_test.go` (impact_check) · 줄 34, 35, 41, 42, 75, 76 · inferred
- `internal/relay/text_events_integration_test.go` (impact_check) · 줄 90, 91 · inferred
- `scripts/package_plugin.py` (impact_check) · 줄 187 · inferred
- `scripts/shell_guard.py` (impact_check) · 줄 4, 4, 13, 13 · inferred
- `scripts/test_shell_guard.py` (impact_check) · 줄 2 · inferred
- `scripts/verify_setup.py` (impact_check) · 줄 39 · inferred
- `videos/knowslink-motion/analyze.py` (impact_check) · 줄 21, 21 · inferred
- `videos/knowslink-motion/beats.json` (impact_check) · 줄 110, 110, 111, 111 · inferred
- `videos/knowslink-motion/index.html` (impact_check) · 줄 44, 44, 47, 47, 49, 49, 349, 349 · inferred
- `videos/knowslink-motion/render.mjs` (impact_check) · 줄 26, 26 · inferred
미확인 21건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

### 결과 (2026-10-20, designer)

- 브랜치 / SHA: `fullops/designer` / 결정 커밋 `a58ada7`(문서 2개). 최종 SHA는 이 기록 커밋이며 worker_done에 적는다. 기준 main `79864592804e595810706f6550ff255ee65407ac`.
- 변경 이유: U8 호스트 표(2026-10-15)가 최신 사용자 범위와 10-17 Claude Code 세션 통합과 달랐다. O-NM-03은 U4 후속 후보였다.
- 결정: [정본 U4·U8](../docs/design-docs/SAR-UX-ISSUES-002.md), [결정 근거](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-HOSTS-DESIGN/report.md).
  - 호스트 표: 웹·Grok Bot·Dots·Claude Code 작업 세션은 대상, Codex·Hermes Agent 0.21.3은 ‘지원 확인 필요’, OpenClaw 보류·Meta Muse 제외(둘 다 화면에 이름 없음). 소스·패키지 근거와 실제 설치·왕복 근거를 다른 열로 나눴다.
  - Claude Code 상태 문구: ‘연결 안내 있음 · 시험 환경 확인 (Claude Code 2.1.302, Windows)’. 격리 relay·합성 회원 근거뿐이라 Grok/Dots의 ‘기존 계정 사용 확인’보다 낮게 둔다.
  - 선택지 5개: 웹에서 바로 사용 / Grok Bot / Dots / Claude Code 작업 세션 / 지원 확인이 필요한 호스트(Codex, Hermes Agent, 그 밖의 호스트 → 웹 사용 안내, 설치 명령 없음).
  - O-NM-03: 요청자 자신의 보낸 요청·종료 기록 행에 ‘입력한 코드 …abcd’(입력 끝 4자, 대소문자 그대로). 입력에서만 얻고 상대 현재 코드로 다시 계산하지 않는다. 있는/없는/차단 코드에 같은 형식. 형식 밖·4자 미만은 ‘입력한 코드 형식 다름’, 이전 요청은 ‘이전 요청’. 전체 코드는 새로 보관·표시·로그하지 않는다. 받은 요청·친구 행·상대 화면에는 없다.
- 검증한 것: `lint.py --from 7986459` ERROR 0·WARNING 0(product-lint `make` unavailable, product-test skipped — 문서만 변경, Windows에 make 없음), `deliverables.py --strict` 문제 0, `--stamp` 재적용 diff 0. 운영 `1dd6800`에 `da97dba`·`467527f` 포함을 `git merge-base --is-ancestor`로 확인.
- 못 한 것: 화면 렌더·시각 검수, 실제 호스트 설치·왕복, D02 변경(필요 없음).
- 남은 일: 4단위 DEV가 표·선택지·O-NM-03을 구현한다. 끝 4자 보관 방식과 세션 연결 관측 방법은 DEV가 정한다.

## SAR-UX-ISSUES-002-CONNECT-VISUAL — 2026-10-20
<!-- fullops-attempt: SAR-UX-ISSUES-002-CONNECT-VISUAL ad0d02cdccd548de8375d2c1769285f3 -->

---
title: SAR-UX-ISSUES-002-CONNECT-VISUAL — 연결 안내 후보 8d8bad4의 직접 시각 검수와 재요청 행 차이 판단
status: draft
updated: 2026-10-20
owner: designer
tasks: [SAR-UX-ISSUES-002-CONNECT-VISUAL]
summary: 연결 안내 후보 8d8bad4의 직접 시각 검수와 재요청 행 차이 판단
attempt: ad0d02cdccd548de8375d2c1769285f3
base: d3e676a0a062cd2e77511e2689436e2cfe5613b6
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-CONNECT-VISUAL — 연결 안내 후보 8d8bad4의 직접 시각 검수와 재요청 행 차이 판단

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: fixture 허용 GET 관찰, VISUAL 보고와 제품 정본 U4 결정 반영, designer 브랜치 커밋. 제품 코드 수정은 하지 않는다.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

4단위(U8 연결 경로·재개+나머지 U7+O-NM-03+low 3) DEV 후보 `8d8bad4cd75a37546ba0f803e3c0c8290ddb6a87`(fullops/dev, 제품 `37febf9`). 원격 Linux make lint/test/integration EXIT=0, fixture 실제 기동 확인, FullOps lint ERROR 0. 고정 SHA 독립 코드 리뷰가 병행 중이다(사용자 D004 '이번만 리뷰').

- DEV 보고(읽기 전용): `C:/Users/shin/orca/workspaces/KnowsLink/dev/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-DEV-test/report.md`(기술 계획·검사·렌더 `render/`·fixture·not_run), DEV 완료 기록 `C:/Users/shin/orca/workspaces/KnowsLink/dev/.fullops-squad/handovers/logs/`의 SAR-UX-ISSUES-002-CONNECT-DEV 절.
- DEV가 남긴 제품 질문: 차단된 요청이 만료된 뒤 같은 코드로 다시 요청하면 새 행이 생기고, 없는 코드는 기존 행을 재사용한다. 행 수 차이로 코드가 실제 회원이었는지 드러날 수 있다.
- 제품 소스는 dev 워크트리에서 읽기 전용으로 확인한다.

## fixture

origin `http://127.0.0.1:18932`(coor SSH 포워딩, 2026-10-20 기동, 수명 12시간). 서버 clone `/tmp/knowslink-connect002-8d8bad4cd75a-8da1e820`, container `kl-connect002-fixture-8d8bad4cd75a-a27ac964`, 고정 SHA 8d8bad4. manifest(origin·accounts·agents·connections·expires_at): `C:/Users/shin/AppData/Local/KnowsLinkDevTools/connect-ux002-dev/fixture-manifest.json`.

- 계정: a(연결 시도 waiting·prepared·approved 각 1, 키 등록 agent 1(긴 영문 이름), 남은 1칸으로 생성·중복 클릭 시험, 보낸 요청: b의 실제 코드·없는 코드·형식 다름 `ab`·이전 요청·취소 1, c에게서 받은 요청 1), b(만료·취소·device 승인 뒤 미저장·remote 앱 복귀 대기·remote 교환 유실·remote 교환 완료 agent, a의 받은 요청 1), c(agent 0개). `POST /fixture/as/{a|b|c}`로 로그인(동일 출처·POST만).
- 허용 GET: `/`·`/home`·`/home/agents`(query 포함)·`/home/friends`·`/home/account`·`/home/legacy`·`/home/new`·`/home/receipts`·`/home/connections/*`·`/connect/*`. 허용 POST: `/home/agents`·`/home/connect`·`/home/confirm`·`/home/cancel`·`/home/friends/request`·`/home/friends/decision`. 그 밖은 403.
- 키·연결 변경은 로그인 뒤 5분 재확인 창 안에서만 된다. 시험 전에 다시 로그인한다.
- tester와 designer가 함께 쓴다. 상태를 바꾸는 POST는 tester가 마지막에 꼭 필요한 것만 실행한다.

## 적용 기준과 예외

fullops-common-0.3.3, `.fullops-squad/project.md`, 제품 정본 `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` U8·U7·U4(O-NM-03), D04. 기준 main d3e676a. 등록 make lint/test는 DEV 같은 SHA 결과를 재사용한다(Windows make 없음).

## 먼저 읽을 문서

FULLOPS.md, common 세 규칙, project.md, 역할 contexts, SAR-UX-ISSUES-002.md U8·U7·U4, D04, 위 DEV 보고, HOSTS-DESIGN 결정 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-HOSTS-DESIGN/report.md`.

## 해야 할 일과 파일 소유권

- [x] '어디에서 사용할까요?' 5선택지와 각 선택 뒤 안내, 처음/기존 연결, 상태별 재개 행, 보낸 요청 '입력한 코드 …abcd', remote 앱 복귀·실패 안내, legacy '이 화면에서 찾기'를 1440·390·320으로 직접 판정(PASS/FAIL/not_run, 근거 캡처). DEV 렌더 `render/`도 참고한다.
- [x] 정본 문구·순서, 대상 외 호스트 미표시, '온라인' 배지 없음, D04 토큰·44px·왼쪽 정렬·피할 스타일 준수.
- [x] 제품 판단: 차단된 요청이 만료된 뒤 같은 코드 재요청은 새 행, 없는 코드는 기존 행 재사용 — 이 차이를 허용할지, 같은 모양으로 맞출 기준(예: 둘 다 새 행, 둘 다 재사용)을 정한다. 결정을 SAR-UX-ISSUES-002.md U4에 반영.
- [x] 캡처는 fixture 허용 GET으로 직접 만든다. 상태를 바꾸는 POST는 하지 않는다. tester 캡처(`C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-TESTER-test/`)는 생기면 참고한다.
- [x] 판정 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-VISUAL/report.md` → 완료 보고 → work.py finish → 커밋 → worker_done(결론과 DEV 수정 필요 목록).

소유 파일: `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-VISUAL/` 하위, SAR-UX-ISSUES-002.md(U4 결정), own route/inbox/log.

## 완료 기준과 검증

시각 판정 결론(수락/조건부/거부)과 근거 캡처, 재요청 행 판단과 DEV 기준, DEV 수정 필요 항목. 스크린리더·실제 200%는 not_run 가능.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 해당 없음(제품 코드 수정 없음).
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 테마 전환 없음.

## 갱신할 산출물

SAR-UX-ISSUES-002.md(U4 결정만).

## 기대 산출물

보고서·결과·캡처, 완료 보고.

## 제약·협업·후속

제품 코드·다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다. 운영·실계정·실메일·외부 발송·실제 호스트 설치·삭제/이동·force push는 하지 않는다. workspace 밖에 파일을 만들지 않는다(브라우저 프로필 등 임시 자료는 보고 폴더나 지정 private 디렉터리).

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-UX-ISSUES-002-CONNECT-VISUAL-packet.json` / SHA `028d02c84922cb5a7393547039900cbb28efce3f` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_update, document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 2, 6, 6, 6, 14, 14, 14, 22, 31, 31, 46, 50, 56, 57, 57, 58, 58, 58, 60, 60, 60, 75 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/archive/designer-2026-10-14-landing.md` (document_read) · 줄 6, 72, 80, 84, 88, 92, 94, 96 · inferred
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 6, 72, 73, 75, 77 · inferred
- `.fullops-squad/contexts/archive/tester-2026-10-14-refactor002.md` (document_read) · 줄 6, 64, 65 · inferred
- `.fullops-squad/contexts/coor.md` (document_read) · 줄 24 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 6, 27, 29, 35, 37, 39, 41, 41, 43, 45 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 6, 58, 60, 62 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 6, 38, 46, 50 · inferred
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` (document_read) · 줄 86, 86, 111 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 7, 192, 195 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read) · 줄 7, 114 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 7, 90 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 7, 14, 20, 237, 282 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 7, 7, 32, 50, 52, 54, 58, 282, 287 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 7, 100 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-DEV-test/report.md` (document_read) · 줄 33, 48 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-TESTER-review/report.md` (document_read) · 줄 65 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-VISUAL-DOCS-review/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-VISUAL/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-001-TESTER-test/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-DEV-review/report.md` (document_read) · 줄 2, 6, 10, 62, 64 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-DEV/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-LIVE.md` (document_read) · 줄 6 · inferred
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
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 15, 21, 42, 59, 65 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-DEV-test/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-TESTER-test/report.md` (document_read) · 줄 6, 10, 25 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-VISUAL/report.md` (document_read) · 줄 6, 6, 10, 10, 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN/report.md` (document_read) · 줄 6, 10, 14, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-HOSTS-DESIGN/report.md` (document_read) · 줄 6, 10, 14, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 14, 61, 62, 63 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-DEV-test/report.md` (document_read) · 줄 2, 6, 10, 18, 94, 96, 96 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-VISUAL/report.md` (document_read) · 줄 6, 6, 10, 10, 18, 89, 89 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 15, 58, 59, 60 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-FIX-test/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-RETEST-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-TESTER-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-APPROVAL-TESTER.md` (document_read) · 줄 2, 6, 10, 12 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-NAMES-TESTER.md` (document_read) · 줄 2, 6, 10, 12, 46 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-FRONTEND-POLISH-001-DEV.md` (document_read) · 줄 31 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-FRONTEND-POLISH-001-VISUAL.md` (document_read) · 줄 6, 10, 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-001-DEV.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-002-DEV.md` (document_read) · 줄 6, 35 · inferred
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
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OPS.md` (document_read) · 줄 20, 22 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-DESIGN.md` (document_read) · 줄 6, 10, 14, 14, 38, 42 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-NAMES-DEV.md` (document_read) · 줄 6, 7, 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-NAMES-TESTER.md` (document_read) · 줄 6, 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-RECEIPT-DEV.md` (document_read) · 줄 6, 20 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 7, 7, 13, 15, 17, 19, 153, 158 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 7, 7, 35, 45, 201, 211 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_dev.md` (document_read) · 줄 102, 103, 104, 112, 128, 129, 136, 137, 156, 162, 183 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_ops.md` (document_read) · 줄 131, 166, 170 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_tester.md` (document_read) · 줄 160, 161, 162, 163, 216, 221, 227, 263 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_designer.md` (document_read) · 줄 136, 137, 138, 139, 144, 196, 197, 200, 204, 212, 249, 262, 283, 284, 285, 286, 287, 288, 294, 295, 297, 298, 299, 300, 301, 302, 304, 306, 311, 313, 315, 341 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_dev.md` (document_read) · 줄 203, 211, 212, 249, 250, 251, 252, 258, 271, 288, 290, 467, 468, 471, 473, 479 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_dev.md` (document_read) · 줄 47, 136, 137, 140, 144, 152, 156, 185, 199 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_tester.md` (document_read) · 줄 135, 136, 137, 203, 349 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_designer.md` (document_read) · 줄 69, 90 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_ops.md` (document_read) · 줄 123, 207, 237, 238, 239, 240, 241, 248, 257, 258, 259, 299, 300, 311, 312, 314, 317, 318, 319, 320, 321, 346, 347, 353, 356, 384, 396 · inferred
- `.fullops-squad/handovers/logs/2026-10-18_to_dev.md` (document_read) · 줄 146, 147, 158, 195, 198, 199, 212 · inferred
- `.fullops-squad/handovers/logs/2026-10-19_to_designer.md` (document_read) · 줄 9, 9, 10, 10, 13, 13, 17, 17, 25, 25, 42, 43, 48, 52, 61, 61, 64, 64, 83, 92, 92, 94, 114, 115, 116, 117, 118, 119, 123, 124, 127, 128, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 145, 147, 148, 152, 154, 156, 158, 159, 203, 206, 206, 209, 209, 210, 210, 213, 213, 217, 217, 225, 225, 233, 244, 244, 254, 258, 263, 264, 266, 266, 281, 300, 300, 302, 320, 321, 322, 323, 324, 325, 328, 329, 331, 332, 334, 335, 336, 337, 338, 339, 340, 341, 342, 343, 343, 344, 345, 346, 347, 348, 350, 352, 353, 356, 358, 360, 362, 363, 392, 395, 395 · inferred
- `.fullops-squad/handovers/logs/2026-10-20_to_designer.md` (document_read) · 줄 9, 9, 10, 10, 13, 13, 17, 17, 25, 25, 33, 42, 43, 58, 62, 68, 69, 69, 71, 71, 86, 99, 99, 101, 119, 120, 121, 122, 123, 124, 127, 128, 130, 131, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 144, 145, 146, 147, 148, 149, 150, 152, 154, 155, 158, 160, 162, 164, 165, 166, 167, 194, 198, 198, 201, 202, 205, 209, 217, 225, 237, 237, 241, 245, 252, 255, 263, 270, 283, 286, 321, 353, 356, 373, 374, 375, 376, 376, 377, 378, 379, 380, 381, 382, 382, 383, 384, 385, 386, 387, 390, 391, 411, 412, 418, 419, 420, 421, 500 · inferred
- `.fullops-squad/handovers/logs/2026-10-20_to_tester.md` (document_read) · 줄 9, 10, 13, 17, 25, 42, 43, 58, 62, 66, 73, 75, 103, 105, 139, 140, 141, 207, 209, 210, 211, 212, 213, 214, 214, 215, 216, 217, 218, 219, 238, 275, 276, 277, 278, 327 · inferred
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 29, 32, 32, 33, 34, 39, 42, 47, 62, 62, 70, 92, 92, 94, 118, 119, 120, 121, 126, 127, 128, 129, 129, 130, 131, 132, 133, 134, 135, 136, 136, 137, 138, 139, 140, 141, 142, 143, 145, 146, 148, 152, 153, 154, 155 · inferred
- `.fullops-squad/handovers/to_tester.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 31, 31, 46, 50, 54, 54, 61, 61, 63, 63, 91, 91, 93, 132, 133, 134, 135, 201, 203, 204, 205, 206, 206, 207, 208, 209, 210, 211, 212, 213, 213, 214, 215, 216, 217, 218, 237, 238, 241, 242, 277, 278, 279, 280 · inferred
- `adapters/README.md` (document_read) · 줄 382 · inferred
- `adapters/src/login.test.ts` (impact_check) · 줄 34 · inferred
미확인 21건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

### 결과 (2026-10-20, designer)

- 브랜치 / SHA: `fullops/designer` / 판정·결정 커밋 `bef4960`, front matter `842bd25`. 최종 SHA는 이 기록 커밋이며 worker_done에 적는다. 기준 main `d3e676a`, 대상 후보 `8d8bad4`(제품 `37febf9`).
- 결론: **조건부 수락.** 5선택지·선택 뒤 안내·처음/기존 연결·상태별 재개 행·연결 화면 단계·보낸 요청 ‘입력한 코드 …abcd’·remote 앱 복귀·legacy ‘이 화면에서 찾기’가 1440·390·320에서 시각 PASS다. 69장 가로 넘침 0, 44px 미만 0, ‘온라인’ 배지 0, 대상 외 호스트 이름 0. critical·high 없음. [판정 보고](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-VISUAL/report.md).
- 제품 결정(재요청 행): 끝난 요청을 같은 코드로 다시 보내면 기존 요청자 행 하나를 다시 대기로 바꾼다(새 행 없음). 있는/없는/차단/수신 한도/자기 코드에서 행 수·순서·상태·끝 4자·만료·숨은 세대 값이 같아야 한다. [정본 U4](../docs/design-docs/SAR-UX-ISSUES-002.md) 반영.
- DEV 수정 필요: V-CN-01 medium(차단 기한 뒤·차단 취소 뒤·상대 수신 한도 10건 경로가 새 대리 행을 만들어 2행이 됨, `human.go:requestFriend`) — 수락 조건, tester 확인으로 충분하고 시각 재검수 불필요. V-CN-02 low(`/home/connections/{없는 ID}` 거부 복귀가 ‘홈으로’, U7은 내 에이전트).
- 검증한 것: fixture 18932 허용 GET과 세션용 `POST /fixture/as/*`만 사용(상태 변경 POST 0). `capture.mjs` 두 실행 `NO_OVERFLOW 51`·`NO_OVERFLOW 18`, 종료코드 0. remote 앱 복귀는 DEV 렌더로 판정. `lint.py --from d3e676a` ERROR 0·WARNING 0(product-lint `make` unavailable, product-test skipped — 문서·캡처만 변경, Windows에 make 없음, 종료코드 1은 이 실행 불가 때문). `deliverables.py --strict` 문제 0.
- 못 한 것: 실제 키보드·200%·스크린리더·보조 브라우저, 상태 변경 흐름의 화면 전이, 실제 호스트 설치.
- 운영 메모: 첫 캡처가 `__Host-` 쿠키 미저장으로 비로그인 GET을 반복해 fixture 429를 일으켰다. 이후 회원별 40회/분을 tester와 나눠 썼다. 브라우저 프로필 `D:/workspace/KnowsLink/.git/worktrees/designer/connect-visual-profile`은 삭제 금지 규칙에 따라 보존했다(정리는 coor 판단).
- 남은 일: DEV가 V-CN-01(필수)·V-CN-02(선택)를 고치고 tester가 행 수 동일성을 확인한다. 헤더 ‘설정’ 변경은 5단위(U3)다.

## SAR-UX-ISSUES-002-LANG-VISUAL — 2026-10-20
<!-- fullops-attempt: SAR-UX-ISSUES-002-LANG-VISUAL 1c7cb2b561b4459e8c1b9fec4dfc7a88 -->

---
title: SAR-UX-ISSUES-002-LANG-VISUAL — 언어·설정 후보 7704e6d의 두 언어 직접 시각 검수
status: draft
updated: 2026-10-20
owner: designer
tasks: [SAR-UX-ISSUES-002-LANG-VISUAL]
summary: 언어·설정 후보 7704e6d의 두 언어 직접 시각 검수
attempt: 1c7cb2b561b4459e8c1b9fec4dfc7a88
base: 08f26982f8ccceb4921b5d82462f57411fb2a112
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-LANG-VISUAL — 언어·설정 후보 7704e6d의 두 언어 직접 시각 검수

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: fixture 허용 GET 관찰, VISUAL 보고와 필요 시 제품 정본 U3 결정 반영, designer 브랜치 커밋. 제품 코드 수정은 하지 않는다.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

마지막 5단위(U3 설정·한국어/English 저장, U7 '설정') DEV 후보 `7704e6d3128831468bb474330210bc2a5099bce7`(fullops/dev, 제품 `2ff9e44`, 이후 테스트·보고만). coor가 같은 SHA를 격리 Linux clone에서 make lint/test·전체 integration 실행해 모두 EXIT=0. FullOps lint ERROR 0. 고정 SHA 독립 코드 리뷰가 병행 중이다(사용자 D005 '리뷰 진행').

- DEV 보고(읽기 전용): `C:/Users/shin/orca/workspaces/KnowsLink/dev/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-DEV-test/report.md`(기술 계획·U3-A 자동 검사·렌더 `render/`·하위 에이전트·not_run). 알려진 한계: 저장 실패(503) 자동 검사 없음, admission 429/503 페이지는 회원 저장값이 아니라 브라우저 쿠키 언어를 따름, README 한글 fragment percent-encoding, 영어 문구는 표본 교정만.
- 제품 소스는 dev 워크트리에서 읽기 전용으로 확인한다.

## fixture

origin `http://127.0.0.1:18933`(coor SSH 포워딩, 2026-10-20 기동). 서버 clone `/tmp/knowslink-lang002-7704e6d31288-ddf6cc71`, container `kl-lang002-fixture-7704e6d31288-a8717b84`, 고정 SHA 7704e6d. manifest(origin·accounts·ids·conversation·expires_at): `C:/Users/shin/AppData/Local/KnowsLinkDevTools/lang-ux002-dev/fixture-manifest.json`.

- 계정: `en`(English 저장), `ko`(한국어 저장), `none`(저장값 없음). `POST /fixture/as/{en|ko|none}`로 로그인(동일 출처). ids: connection_waiting·prepared·approved, device, gate, message, remote — 번역 화면 전체에 닿는 합성 상태.
- 허용 GET: `/`, `/auth/verify`, `/home`, `/home/name`, `/home/agents`, `/home/friends`, `/home/account`, `/home/legacy`, `/home/new`, `/home/receipts`, `/home/conversation`와 접두 `/home/connections/`·`/connect/`·`/home/gates/`·`/home/friends/`·`/remote/consent/`·`/lang/`. 허용 POST: `/home/language`, `/home/profile`, `/home/friends/request`, `/home/friends/decision`, `/home/agents`, `/auth/logout`. 그 밖은 403.
- 쿠키는 `__Host-` Secure라 http cookie jar가 버릴 수 있다. 쿠키를 직접 실어 보낸다(이전 단위 보고 참고).
- tester와 designer가 함께 쓴다. 저장값을 바꾸는 POST(언어 저장·이름·친구·에이전트)는 tester가 마지막에, 원래 값으로 되돌리거나 영향 범위를 기록한다.

## 적용 기준과 예외

fullops-common-0.3.3, `.fullops-squad/project.md`, 제품 정본 `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` U3·U7, D04, 공개 첫페이지 규칙. 기준 main 08f2698. 등록 make lint/test는 coor의 같은 SHA Linux 결과를 재사용한다(Windows make 없음).

## 먼저 읽을 문서

FULLOPS.md, common 세 규칙, project.md, 역할 contexts, SAR-UX-ISSUES-002.md U3·U7, D04, 위 DEV 보고.

## 해야 할 일과 파일 소유권

- [x] 설정 화면(한국어·영어), 공개 첫페이지 두 언어와 '한국어 / English' 링크 위치, 주요 회원 화면(친구·내 에이전트·대화·승인·조회·연결)의 영어판을 1440·390·320으로 직접 판정(PASS/FAIL/not_run, 근거 캡처). DEV 렌더 `render/`도 참고한다.
- [x] 정본 U3: 설정 순서·문구, 언어 선택 형태(국기·드롭다운 JS 없음), 공개 페이지 정보 구조·설치 사실 불변, 영어판이 같은 의미·기능이며 새 약속을 하지 않는지, 영어 문구 품질(표본: 버튼·상태·오류·빈 상태·도움말), 긴 영문 줄바꿈, D04 토큰·44px·왼쪽 정렬.
- [x] 제품 판단: DEV 알려진 한계 중 admission 429/503 페이지가 회원 저장값이 아니라 브라우저 쿠키 언어를 따르는 것을 허용할지. 영어 용어집 수정이 필요하면 목록으로 남긴다(DEV 수정 항목).
- [x] 캡처는 fixture 허용 GET으로 직접 만든다(언어 저장 POST는 하지 않는다; `en`·`ko`·`none` 계정과 `/lang/{ko|en}` GET으로 언어를 바꾼다). tester 캡처(`C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-TESTER-test/`)는 생기면 참고한다.
- [x] 판정 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-VISUAL/report.md` → 완료 보고 → work.py finish → 커밋 → worker_done(결론과 DEV 수정 필요 목록).

소유 파일: `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-VISUAL/` 하위, 필요 시 SAR-UX-ISSUES-002.md(U3 결정), own route/inbox/log.

## 완료 기준과 검증

시각 판정 결론(수락/조건부/거부)과 근거 캡처, 영어 문구 수정 목록, DEV 수정 필요 항목.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 해당 없음(제품 코드 수정 없음).
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 테마 전환 없음. 두 언어 3폭.

## 갱신할 산출물

없음.

## 기대 산출물

보고서·결과·캡처, 완료 보고.

## 제약·협업·후속

제품 코드·다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다. 운영·실계정·실메일·외부 발송·삭제/이동·force push는 하지 않는다. workspace 밖에 파일을 만들지 않는다(브라우저 프로필 등 임시 자료는 보고 폴더 아래).

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-UX-ISSUES-002-LANG-VISUAL-packet.json` / SHA `768e7d2198c231d2a868db0a26cc6a3424392ac6` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_update, document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read, document_update) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 2, 6, 6, 6, 14, 14, 14, 22, 31, 31, 45, 49, 56, 56, 57, 57, 57, 59, 59, 59 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/archive/designer-2026-10-14-landing.md` (document_read) · 줄 6, 72, 80, 84, 88, 92, 94, 96 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 6, 27, 29, 35, 37, 39, 41, 41, 43, 45 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 58, 60, 62, 64 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 6, 38, 46, 50 · inferred
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` (document_read) · 줄 86, 86, 111 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 14, 20, 53 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 7, 32, 53, 55, 57, 61 · inferred
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
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 15, 21, 42, 59, 65 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-DEV-test/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-TESTER-test/report.md` (document_read) · 줄 6, 10, 25 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-VISUAL/report.md` (document_read) · 줄 6, 6, 10, 10, 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 6, 10, 16, 71, 98 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-DEV-test/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-FIX-REVIEW-review/report.md` (document_read) · 줄 2, 9, 15, 15, 15, 21, 88, 88, 93 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-TESTER-test/report.md` (document_read) · 줄 6, 10, 92 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-VISUAL/report.md` (document_read) · 줄 6, 6, 10, 10, 12, 12, 21, 62, 75 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN/report.md` (document_read) · 줄 6, 10, 14, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-HOSTS-DESIGN/report.md` (document_read) · 줄 6, 10, 14, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 14, 61, 62, 63 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-DEV-test/report.md` (document_read) · 줄 2, 6, 10, 18, 94, 96, 96 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-VISUAL/report.md` (document_read) · 줄 6, 6, 10, 10, 18, 89, 89 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 15, 58, 59, 60 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-FIX-test/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-RETEST-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-TESTER-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-APPROVAL-TESTER.md` (document_read) · 줄 2, 6, 10, 12 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-CONNECT-TESTER.md` (document_read) · 줄 2, 6, 10, 12 · inferred
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
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 7, 13, 15, 19, 21, 23, 25 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 7, 35, 45, 55 · inferred
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
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 36, 44, 58, 58, 66, 88, 88, 90, 107, 108, 109, 110, 110, 111, 112, 113, 114, 115, 115, 116, 117, 118, 119, 120, 121, 122, 122, 123, 124, 125, 126, 127, 128, 129, 130, 132, 133, 134, 135 · inferred
- `.fullops-squad/handovers/to_tester.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 31, 31, 45, 49, 53, 53, 59, 59, 61, 61, 89, 89, 91, 193, 195, 196, 197, 198, 198, 199, 200, 201, 202, 203, 203, 204, 205, 206, 207, 208, 209, 210, 210, 211, 212, 213, 214, 215, 234, 235, 236, 273, 274, 275, 276 · inferred
- `internal/relay/friend_rerequest_test.go` (impact_check) · 줄 1, 12 · inferred
- `internal/relay/human.go` (impact_check) · 줄 182 · inferred
미확인 24건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

### 결과 (2026-10-20, designer)

- 브랜치 / SHA: `fullops/designer` / 판정·결정 커밋 `d134258`, front matter `ba7d9e6`. 최종 SHA는 이 기록 커밋이며 worker_done에 적는다. 기준 main `08f2698`(지시서 `0b29416`), 대상 후보 `7704e6d`(제품 `2ff9e44`).
- 결론: **조건부 수락.** 설정(내 이름 → 내 친구 코드 → 언어 → 로그인과 보안, radio 한국어/English + 언어 저장, select·국기·JS 없음), 헤더 ‘설정’ `aria-current=page`, 공개 첫페이지 두 언어(제목 6·본문 106줄 구조 동일, 설치 사실 불변, 새 약속 없음), header 안 ‘한국어 / English’ 링크, 대화·친구·내 에이전트·승인·조회·legacy·연결 3단계·device/remote 동의 영어판이 1440·390·320에서 시각 PASS다. 90장 가로 넘침 0, 44px 미만 0, `html lang` 일치, 영어 화면 한글은 사용자 이름·‘한국어’뿐. critical·high 없음. [판정 보고](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-VISUAL/report.md).
- 제품 결정: admission 429/503 거부 화면이 회원 저장값 대신 브라우저 언어 쿠키를 따르는 것을 허용한다(문구·`lang` 일치, 재시도·안전 복귀, 저장 언어 적용 주장 금지). [U3](../docs/design-docs/SAR-UX-ISSUES-002.md)에 반영. 영어 문장의 사용자 지칭은 ‘your’, 메뉴·제목 고유 이름(My agents 등)은 유지.
- DEV 수정 필요: V-LANG-01 medium(저장값 있는 회원에게 공개 첫페이지·remote 동의의 언어 링크가 무동작이고 쿠키만 바뀜 → 현재 언어 + ‘설정에서 변경’ 링크로 대체), V-LANG-02 low(‘my computer/my agent/my account’ 등 1인칭 영어 약 25개 값과 설정 도움말·지원 언어 오류 문구, 보고서 목록), V-LANG-03 low(저장값 없는 회원의 설정에 ‘아직 계정에 저장되지 않음’ 표시). 세 건 모두 tester 확인으로 충분하고 designer 재검수 불필요.
- 검증한 것: fixture 18933 허용 GET과 세션용 `POST /fixture/as/*`만 사용(언어 저장 등 상태 변경 POST 0). 언어 쿠키는 `/lang` 응답과 같은 값을 CDP로 주입. `capture.mjs` `NO_OVERFLOW 90`, 종료코드 0. V-LANG-01은 curl GET으로 재현(`html lang="en"`). 영어 번역표 두 파일 전체 검토. `lint.py --from 08f2698` ERROR 0·WARNING 0(product-lint `make` unavailable, product-test skipped: Windows make 없음, 지시서대로 coor Linux 결과 재사용).
- 못 한 것: 실제 저장·실패·재로그인 전이, 실제 `/lang` 클릭의 Secure 쿠키 저장, 키보드·200%·스크린리더·보조 브라우저, admission 화면 캡처, remote 앱 복귀(DEV 렌더 3장이 meta refresh 뒤 Chrome 오류 화면이라 판정 불가 — DEV 렌더 증거 결함).
- 운영 메모: 브라우저 프로필은 보고 폴더 `profile/`(`.gitignore` 제외)에 남겼고 삭제하지 않았다. tester와 회원별 요청 제한을 나눠 써 캡처에 약 25분 걸렸다.
- 남은 일: DEV가 V-LANG-01(필수)·V-LANG-02/03을 고치고 tester가 확인한다. U7-A 키보드·200%·스크린리더는 tester 몫이다.

## SAR-SCHEDULE-001-DESIGN — 2026-10-20
<!-- fullops-attempt: SAR-SCHEDULE-001-DESIGN eaa1057e93ce4166a4cff3bd9b02d93a -->

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
