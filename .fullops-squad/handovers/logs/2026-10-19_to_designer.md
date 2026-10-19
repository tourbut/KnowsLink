---
title: designer 완료 기록
status: draft
updated: 2026-10-19
owner: designer
summary: 지시서와 완료 보고를 보존한다.
---

## SAR-UX-ISSUES-002-NAMES-VISUAL — 2026-10-19
<!-- fullops-attempt: SAR-UX-ISSUES-002-NAMES-VISUAL daf215708ff94a038a0a73a6da39763b -->

---
title: SAR-UX-ISSUES-002-NAMES-VISUAL — 2단위 후보 b2e683a의 직접 시각 검수와 거절 표시·select 잘림·D-NM-01 제품 판단
status: draft
updated: 2026-10-19
owner: designer
tasks: [SAR-UX-ISSUES-002-NAMES-VISUAL]
summary: 2단위 후보 b2e683a의 직접 시각 검수와 거절 표시·select 잘림·D-NM-01 제품 판단
attempt: daf215708ff94a038a0a73a6da39763b
base: 7a66f71341f00f9044f3d1e7da981ef4b57cfca5
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-NAMES-VISUAL — 2단위 후보 b2e683a의 직접 시각 검수와 거절 표시·select 잘림·D-NM-01 제품 판단

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-19
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: tester 캡처·보고와 fixture 화면 관찰로 직접 시각 판정, 제품 판단 기록, D04·D11 제품 정본의 필요한 보완. 제품 코드 수정, 운영·실계정·외부 발송, 삭제·이동, force push는 하지 않는다.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

2단위(U2/U4+U7) DEV 후보 `b2e683ae677c78cd7bdc47c75820bcb191b94838`(제품 516bf62, fullops/dev)이 tester 독립 QA를 통과했다(tester 48ff7b1, critical/high 0, low D-NM-01). 고정 SHA 코드 리뷰는 병행 중이다.

- tester 보고·캡처(읽기 전용): `C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/report.md`와 같은 폴더의 `chrome-*.png`(1440·390·320), `browser-chrome-layout.json`. 관찰 O-NM-01~07은 보고서 '관찰' 표.
- DEV 보고(읽기 전용): `C:/Users/shin/orca/workspaces/KnowsLink/dev/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-DEV-test/report.md`('남은 위험과 판단 필요'), DEV 렌더 캡처 `render/*.png`.
- fixture `http://127.0.0.1:18930`은 tester 비가역 검사 뒤 상태다. 추가 관찰은 허용 GET만 사용한다(목록은 tester 지시서 로그). 상태를 바꾸는 POST는 하지 않는다.

## 적용 기준과 예외

fullops-common-0.3.3, 문서 규칙 `.fullops-squad/docs/agents/document-writing.md`, 제품 정본 `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` U2·U4·U7, D04 `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`, D02 `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md`. 기준 main 7a66f71. frontend.json의 codex host는 사용자 Codex 제외 지시로 적용하지 않는다.

## 먼저 읽을 문서

SAR-UX-ISSUES-002.md(U2·U4·U7·'화면 배치와 자기 비평'), D04, SAR-HUMAN-PEER.md(HP-04 상태 구분), 위 tester·DEV 보고, `.fullops-squad/contexts/designer.md`.

## 해야 할 일과 파일 소유권

- [x] 직접 시각 판정: 친구 목록·빈 상태(C agent 0개)·친구 추가·내 계정 친구 코드·내 에이전트·조회 radio 선택기·legacy·연결 화면의 이름/종류/구분값 표시를 3폭 캡처로 판정(PASS/FAIL/not_run과 근거 캡처 경로). D04 토큰·44px·왼쪽 정렬·피할 스타일 준수.
- [x] 제품 판단 1(O-NM-01): 요청자에게 '거절'을 보일지. HP-04 상태 구분과 U4 존재 비노출 사이 결정. 바꾸면 기대 문구·상태를 DEV가 구현할 수 있게 확정한다.
- [x] 제품 판단 2(O-NM-02): 새 대화 발신/수신 select에서 40자 이름일 때 종류·구분값이 닫힌 상태에 안 보임. 허용 또는 수정 기준(예: 조회처럼 radio, 또는 표시 순서·생략 규칙) 결정.
- [x] 제품 판단 3(D-NM-01 low): 제어·방향 문자 거부 시 사유 설명·입력 되돌리기 기준 문구 확정.
- [x] O-NM-03~07 중 제품 결정이 필요한 것은 결정하고, 아니면 '현행 유지' 또는 '후속 단위'로 분류한다.
- [x] 판정·결정 기록: `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-VISUAL/report.md`. 결정이 제품 정본을 바꾸면 SAR-UX-ISSUES-002.md 또는 D04에 반영(front matter 갱신).
- [x] 완료 보고 → work.py finish → 커밋 → worker_done. body에 시각 판정 결론과 DEV 수정 필요 항목 목록을 적는다.

소유 파일: 위 VISUAL 보고, SAR-UX-ISSUES-002.md·D04·D02의 필요한 절, own route/inbox/log.

## 완료 기준과 검증

시각 판정 결론(수락/조건부 수락/거부)과 근거 캡처, 판단 1–3의 결정과 DEV 구현 기준이 기록됐다. 미관측(스크린리더·실제 200%·보조 브라우저)은 not_run으로 둔다. 문서 lint ERROR 0.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 해당 없음(판정·제품 기준만).
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 테마 전환 없음.

## 갱신할 산출물

D04, D11(route 추천). 실제 결정이 바꾸는 정본만 갱신한다.

## 기대 산출물

VISUAL 보고서, 필요한 제품 정본 보완, 완료 보고.

## 제약·협업·후속

기술 구현 방법은 정하지 않는다(DEV 담당). 다른 역할 인박스·PLANS·board·GitHub·제품 코드는 수정하지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-UX-ISSUES-002-NAMES-VISUAL-packet.json` / SHA `7a66f71341f00f9044f3d1e7da981ef4b57cfca5` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_update, document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.75, 'evidence': 0.9, 'contradicts': 0.3, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read, document_update) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 2, 2, 6, 6, 6, 7, 14, 14, 14, 14, 31, 31, 32, 32, 37, 41, 50, 50, 50, 53, 53, 72 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/archive/designer-2026-10-14-landing.md` (document_read) · 줄 6, 72, 80, 84, 88, 92, 94, 96 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 6, 27, 29, 35, 37, 39, 41, 41, 43, 45 · inferred · {'relevant': 0.51, 'evidence': 0.62, 'contradicts': 0.13, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 58 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 6, 38 · inferred
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` (document_read) · 줄 86, 86, 111 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 14 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 44 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 7, 32, 38, 200 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-DEV-test/report.md` (document_read) · 줄 33, 42, 48 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-TESTER-review/report.md` (document_read) · 줄 65 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-TESTER-test/report.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-VISUAL-DOCS-review/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-VISUAL/report.md` (document_read) · 줄 6, 10, 33 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL-DOC-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL-FINAL-DOC-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/final-report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/report.md` (document_read) · 줄 6, 37 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-FINAL-RECORDS-REVIEW-review/report.md` (document_read) · 줄 41 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-TESTER-test/css-followup/report.md` (document_read) · 줄 57 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-TESTER-test/report.md` (document_read) · 줄 61 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/css-followup/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/report.md` (document_read) · 줄 6, 85 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 47 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-TESTER-review/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL-RESULT-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL/report.md` (document_read) · 줄 6, 10, 16, 32, 52 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-LIVE/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-test/report.md` (document_read) · 줄 29, 67, 85 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-2-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-review/report.md` (document_read) · 줄 6, 10, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md` (document_read) · 줄 6, 10, 64, 100, 159, 179, 181 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL/report.md` (document_read) · 줄 6, 10, 77 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-VISUAL/report.md` (document_read) · 줄 6, 10, 72 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN/report.md` (document_read) · 줄 6, 10, 14, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 15, 58, 59, 60 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-FIX-test/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-RETEST-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-TESTER-test/report.md` (document_read) · 줄 6, 10, 71 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-FRONTEND-POLISH-001-DEV.md` (document_read) · 줄 31, 38 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-FRONTEND-POLISH-001-VISUAL.md` (document_read) · 줄 6, 10, 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-OPS.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-VISUAL.md` (document_read) · 줄 2, 6, 10, 24, 38 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001-VISUAL.md` (document_read) · 줄 6, 10, 16, 26, 28, 32, 47 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md` (document_read) · 줄 78 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 44, 57, 68, 82, 104 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-DEV.md` (document_read) · 줄 15, 43 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-VISUAL.md` (document_read) · 줄 2, 6, 10, 22 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-OPS.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-VISUAL.md` (document_read) · 줄 2, 6, 10, 25 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL.md` (document_read) · 줄 2, 6, 10, 21 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OPS.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-DESIGN.md` (document_read) · 줄 6, 10, 14, 14, 38, 42 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-RECEIPT-DEV.md` (document_read) · 줄 6, 20 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read, document_update) · 줄 7, 35 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 전체/미확인 · inferred · {'relevant': 0.67, 'evidence': 0.88, 'contradicts': 0.26, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 278 · inferred
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
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 29, 36, 38, 42, 44, 46, 46, 62, 62, 74, 87, 87, 98, 98, 100, 122, 123, 124, 125, 126, 127, 128, 131, 134, 135 · inferred
- `.fullops-squad/handovers/to_tester.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 29, 29, 31, 31, 45, 49, 51, 55, 55, 62, 62, 93, 93, 95, 197, 199, 200, 201, 202, 203, 204, 205, 260, 261 · inferred
- `cmd/relay/main.go` (impact_check) · 줄 86 · inferred
- `internal/relay/capacity.go` (impact_check) · 줄 215, 276 · inferred
- `internal/relay/cleanup.go` (impact_check) · 줄 14 · inferred
- `internal/relay/connections_test.go` (impact_check) · 줄 356 · inferred
- `internal/relay/human_pages.go` (impact_check) · 줄 29, 33, 38 · inferred
- `internal/relay/human_pages_test.go` (impact_check) · 줄 120, 168, 181 · inferred
- `internal/relay/landing_fixture_test.go` (impact_check) · 줄 99 · inferred
- `internal/relay/landing_pages.go` (impact_check) · 줄 7, 38, 48 · inferred
- `internal/relay/member_pages.go` (impact_check) · 줄 18, 21, 77, 101, 102, 123, 125, 127 · inferred
- `internal/relay/member_ui_test.go` (impact_check) · 줄 115, 145, 147, 200 · inferred
- `internal/relay/public_messages_integration_test.go` (impact_check) · 줄 439, 454 · inferred
- `internal/relay/text_events_integration_test.go` (impact_check) · 줄 69 · inferred
- `videos/knowslink-motion/render.mjs` (impact_check) · 줄 97 · inferred
미확인 21건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

- 브랜치 / SHA: fullops/designer / 판정·정본 커밋 `8f9beb60cfd6522af1e0d125ca046969f393db62`(기준 main 7a66f71, 지시서 00b6836). 대상 후보 `b2e683ae677c78cd7bdc47c75820bcb191b94838`.
- 결론: **조건부 수락.** 친구·친구 추가·내 계정 코드·조회 radio·legacy·친구 상세·0개 상태의 이름/종류/구분값 표시는 3폭 시각 PASS다. 새 대화 select는 390·320에서 동명 발신·수신을 구별하지 못해 FAIL(medium)이다. 내 에이전트 상세의 두 링크가 붙어 보이는 V-NM-01(low)을 새로 찾았다. critical·high 없음.
- 결정: (1) 요청자에게 ‘거절’ 표시를 유지하고 받은 요청 거절 앞에 안내 한 줄을 둔다. (2) 새 대화 발신·수신 선택기를 조회와 같은 radio 라벨로 바꾼다. (3) 금지 문자 거부 문구 ‘이름에 줄바꿈·제어 문자·방향 제어 문자는 쓸 수 없습니다. 보이는 글자로 1–40자를 입력하세요.’와 금지 문자 제거 후 입력 복구. O-NM-03·04는 후속 단위, O-NM-05·06·07은 현행 유지.
- DEV 수정 필요: ① 새 대화 radio(medium) ② 두 링크 분리(low) ③ 이름 거부 문구·입력 복구(low) ④ 거절 안내 한 줄(low). 수정 후보에서 designer가 해당 화면만 재검수한다.
- 변경 이유 / 정본: SAR-UX-ISSUES-002.md U2·U4에 결정 반영. D02·D04·D11은 바뀌는 기준이 없어 수정하지 않았다.
- 검증한 것: tester Chrome 캡처 직접 열람(1440·390·320), 고정 SHA 템플릿·CSS 읽기, FullOps lint ERROR 0(WARNING 1은 기존 PLANS.md 크기, product-lint는 Windows에 make 없음으로 실행 불가, product-test는 문서 전용이라 생략).
- 못 한 것(not_run): 승인·조회 성공·충돌 확장·내 이름 화면 캡처 없음, 스크린리더·실제 200%·Firefox/Safari·실기, Edge 캡처 직접 열람, fixture 재관찰.
- 산출물: [시각 판정 보고](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-VISUAL/report.md), [정본](../docs/design-docs/SAR-UX-ISSUES-002.md).
- 남은 일: coor가 DEV 수정 4건을 배정하고 수정 후보의 designer 재검수·코드 리뷰를 진행한다. O-NM-04는 1단위 잔여 후속 과제로 PLANS에 둔다.
