---
title: SAR-UX-ISSUES-002-LANG-DEV — 리뷰·designer·tester 수정 묶음(admission 언어·README 언어 표시·저장값 회원 링크·영어 문구)을 반영한다
status: draft
updated: 2026-10-20
owner: dev
tasks: [SAR-UX-ISSUES-002-LANG-DEV]
summary: 리뷰·designer·tester 수정 묶음(admission 언어·README 언어 표시·저장값 회원 링크·영어 문구)을 반영한다
attempt: d151cd6ef79c4ccb8bec68adb6ff72b6
base: 131544873193795e77a2886dcdd81e8a2658c0f6
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-LANG-DEV — 리뷰·designer·tester 수정 묶음(admission 언어·README 언어 표시·저장값 회원 링크·영어 문구)을 반영한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 아래 수정 묶음과 관련 테스트·기술 문서. 새 경로·새 의존성·권한 확대·DB migration·실계정·운영 배포·삭제/이동·force push는 coor에 ask.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

같은 과제의 두 번째 attempt다. 후보 7704e6d(제품 2ff9e44)는 고정 SHA 리뷰(critical/high 0, medium F1·F2, low F3), designer 조건부 수락(V-LANG-01 medium, V-LANG-02·03 low; admission 429/503은 브라우저 쿠키 언어를 따르도록 U3 결정), tester(전 check PASS, medium D-LG-01, low O-LG-01~04)를 받았다. designer 결정은 main 3f81a7d에 통합됐다. fullops/dev에 main을 merge해 정본을 받는다.

- 리뷰: `C:/Users/shin/orca/workspaces/KnowsLink/coor/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-CODE-REVIEW-review/report.md`(F1–F3).
- designer: main의 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-VISUAL/report.md`(V-LANG-01~03, V-LANG-02 대상 문구 목록), 정본 SAR-UX-ISSUES-002.md U3 갱신.
- tester: `C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-TESTER-test/report.md`(D-LG-01, O-LG-01~04).

## 적용 기준과 예외

이전 attempt와 같다. 원격 Linux make lint/test/integration과 fixture 실제 기동 확인을 유지한다. lint는 깨끗한 임시 detached worktree에서 실행하고 force 없이 정리한다. lint 기준 1315448.

## 먼저 읽을 문서

위 세 보고의 해당 절, SAR-UX-ISSUES-002.md U3 갱신 문단, admission(요청 한도·일시 장애) 화면 코드, i18n 관련 파일, 설정 화면 템플릿, 이전 완료 기록.

## 해야 할 일과 파일 소유권

- [ ] (medium, 리뷰 F1 = tester D-LG-01) 회원 요청 한도 429·일시 장애 503 등 admission 화면이 designer 결정대로 브라우저 쿠키 언어를 따르고 `Content-Language`와 html lang이 일치하게 한다(현재 항상 한국어).
- [ ] (medium, 리뷰 F2) 영어 화면에서 한국어 README 등 한국어 외부 문서로 가는 링크에 언어 표시(U3 '외부 한국어 문서 링크는 언어를 표시한다').
- [ ] (medium, designer V-LANG-01 = 리뷰 F3) 저장값이 있는 회원에게 공개·remote 동의 화면의 '한국어 / English' 링크가 쿠키만 바꾸고 동작하지 않는 문제: designer 결정대로 '설정에서 변경' 링크로 대체.
- [ ] (low, V-LANG-02) 영어 1인칭 문구 약 25개(designer 보고 목록)와 설정 도움말 수정.
- [ ] (low, V-LANG-03) 저장값 없는 회원 설정 화면에 '미저장' 상태 표시(정본 문구).
- [ ] (low, O-LG-03) `details` 요약 장식 문자(`::before`)가 접근성 이름에 들어가지 않게(공용 스타일, 두 언어).
- [ ] 바뀐 경계의 자동 검사(429/503 두 언어·Content-Language, 저장값 회원의 링크 대체, README 링크 언어 표시). 문구 복제 테스트는 쓰지 않는다.
- [ ] 원격 Linux make lint/test/integration(최종 제품 SHA), fixture 실제 기동, FullOps lint.py --from 131544873193795e77a2886dcdd81e8a2658c0f6, DEV 보고 수정 절, packet-outcomes, finish, 커밋, worker_done.

O-LG-01(로그인 전 `/auth/verify`는 쿠키 언어), O-LG-02(저장 버튼 단일 언어), O-LG-04(legacy 기술 용어 원문)는 현행 유지한다.

소유 파일: 이전 attempt와 같은 범위. 예상 100–300줄.

## 완료 기준과 검증

medium 3건이 자동 검사로 확인되고 low 3건이 반영됐다. 기존 회귀 없음, Linux make lint/test/integration EXIT=0(최종 제품 SHA), lint ERROR 0, fixture 실제 기동. 이후 coor 수정분 확인·tester 재확인(designer 재검수 불필요).

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04, SAR-UX-ISSUES-002.md U3.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 기존 템플릿·i18n 구조. 새 토큰·JS 없음.
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(임시 worktree).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 테마 전환 없음. 바뀐 화면 두 언어 3폭 넘침 0.

## 갱신할 산출물

D10, D13 중 실제 변경된 것.

탐색 패킷의 task/role/attempt·기준 SHA와 direct_edit/impact_check/document_read/document_update 경로·구간·근거·partial/unknown을 연결한다. 명시적 없음과 미확인을 구분한다. 필수·충돌·주의 후보와 budget 잔여 확인 목록을 보존한다.

## 기대 산출물

코드·테스트, 기술 문서, DEV 보고 수정 절, 완료 보고.

## 제약·협업·후속

미추적 `videos/service-intro/`·`qa-reports/SAR-MOTION-INTRO-001-DEV-test/` 하위 미추적 폴더는 건드리지 않는다. workspace 밖에 파일을 만들지 않는다(private 도구 디렉터리 제외). 다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다.

## 이전 완료 기록

- .fullops-squad/handovers/logs/2026-10-20_to_dev.md

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-UX-ISSUES-002-LANG-DEV-packet.json` / SHA `af16ac52044b3958d3a89afe6900a0bc0051f15a` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 30, 37, 41, 50, 91 · inferred · 필수
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_read) · 줄 5, 6, 6, 10, 14, 32, 40, 64, 76, 153, 161, 161, 176, 178, 184, 193, 263, 263, 270, 279 · inferred · 필수
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 2, 2, 6, 6, 7, 7, 7, 7, 14, 14, 14, 14, 29, 29, 29, 29, 31, 31, 32, 32, 32, 33, 33, 33, 41, 41, 45, 45, 45, 46, 47, 47, 48, 48, 49, 51, 60, 60, 64 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 19, 20, 32, 34, 56 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 3, 17 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 3 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 26 · inferred · 필수
- `.fullops-squad/contexts/archive/designer-2026-10-14-landing.md` (document_read) · 줄 2, 5, 10, 70, 99 · inferred
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 52 · inferred
- `.fullops-squad/contexts/archive/tester-2026-10-14-refactor002.md` (document_read) · 줄 2, 5, 10, 17, 61, 78 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 2, 5, 6, 10, 17, 19, 32, 36, 39, 41, 45 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 58, 60, 62, 64 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 2, 5, 6, 10, 12, 18, 38, 46, 50 · inferred
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` (document_read) · 줄 5, 86, 111 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read, document_update) · 줄 115 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read, document_update) · 줄 7, 7, 14, 20, 30, 32, 169, 200, 229 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read, document_update) · 줄 6 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 5, 18, 83 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 5, 18 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 5, 34, 103 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md` (document_read) · 줄 5, 20, 22 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_read, document_update) · 줄 7, 7, 30, 32, 53, 55, 61, 61, 62, 62, 62, 63, 63, 70, 72, 72, 76, 117, 117, 174, 238, 276, 286, 296, 333 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read, document_update) · 줄 55 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/COOR/dev-fix-2-final/coordinator-handover-supplement.md` (document_read) · 줄 41 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.2.0-review/report.md` (document_read) · 줄 16, 44 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.3.0-FINAL-review/report.md` (document_read) · 줄 26, 28, 42 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.4.0-review/report.md` (document_read) · 줄 16, 32 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.4.2-review/report.md` (document_read) · 줄 24, 30, 40, 48, 52 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.4.4-FINAL-REVIEW-review/report.md` (document_read) · 줄 27, 51, 55 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-COPY-REVIEW-review/report.md` (document_read) · 줄 18, 34 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-REVIEW-review/report.md` (document_read) · 줄 5, 12, 16, 31, 36, 37, 38, 58, 62, 63, 64 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-TESTER.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-N1-review/report.md` (document_read) · 줄 73 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-review/report.md` (document_read) · 줄 17 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-TESTER-PUBLIC.md` (document_read) · 줄 5, 32 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-TESTER.md` (document_read) · 줄 5, 34 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-002-TESTER.md` (document_read) · 줄 5, 26 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-HOOK-QUOTE-001-REVIEW-review/report.md` (document_read) · 줄 5, 12, 14, 14, 75, 96 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-DEPLOY-001-OPS-FINAL-review/report.md` (document_read) · 줄 14, 36 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-QA-DOCS-review/report.md` (document_read) · 줄 7, 26 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-TESTER-review/report.md` (document_read) · 줄 5, 12, 13, 15, 18, 21, 51, 65, 67, 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-TESTER-test/report.md` (document_read) · 줄 5, 16, 24, 26, 27, 30 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-VISUAL-DOCS-review/report.md` (document_read) · 줄 25, 25 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-VISUAL/report.md` (document_read) · 줄 5, 60 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-001-TESTER-test/report.md` (document_read) · 줄 5, 12, 17, 31, 37, 41, 58, 72, 81 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-DEV-review/report.md` (document_read) · 줄 5, 12, 14, 14, 63 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-LIVE.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-LOGIN-001-REVIEW-review/report.md` (document_read) · 줄 32 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GROK-WEBHOOK-001-TESTER-RESULT-REVIEW-review/report.md` (document_read) · 줄 14, 24, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GROK-WEBHOOK-001-TESTER-REVIEW-review/report.md` (document_read) · 줄 5, 12, 15, 25, 60, 66, 67, 68, 69, 70, 82, 86, 87 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GROK-WEBHOOK-001-TESTER-test/report.md` (document_read) · 줄 5, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-DESIGN-review/report.md` (document_read) · 줄 16, 26, 28, 28, 29 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-DEV-test/runtime/report.md` (document_read) · 줄 92, 92, 257, 265 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-DEV-test/runtime/v01-v02-report.md` (document_read) · 줄 20, 20, 38 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-LIVE/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-TESTER-DOC-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-TESTER-review/report.md` (document_read) · 줄 5, 14, 20, 62 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-TESTER-test/report.md` (document_read) · 줄 5, 16, 16, 33, 35, 39 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-UI-DELTA-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/final-report.md` (document_read) · 줄 5, 18, 26, 40, 61, 76 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/report.md` (document_read) · 줄 5, 18, 20, 92, 98 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CODE-REVIEW-review/author-final-handoff.md` (document_read) · 줄 23, 23, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CODE-REVIEW-review/report.md` (document_read) · 줄 12, 12, 24, 54, 57, 58, 58, 59, 62 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CSS-REVIEW-review/report.md` (document_read) · 줄 12, 12, 31, 47, 73, 73 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DESIGN-REVIEW-review/report.md` (document_read) · 줄 24, 40, 40, 42, 54, 57, 59, 66 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/css-resume/report.md` (document_read) · 줄 33, 35, 45, 49, 49, 59 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/report.md` (document_read) · 줄 18, 61 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-FINAL-RECORDS-REVIEW-review/report.md` (document_read) · 줄 23, 23, 27, 47, 48, 49, 53, 53, 59, 59, 65, 65, 71, 77, 80 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-TESTER-test/css-followup/report.md` (document_read) · 줄 5, 12, 23, 29, 57, 57, 77 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-TESTER-test/report.md` (document_read) · 줄 5, 12, 30, 34, 45, 47, 57, 94, 96 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-TESTER-test/web-20261014-232426-595664/report.md` (document_read) · 줄 5, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/css-followup/dev-final-handoff.md` (document_read) · 줄 38, 38 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/css-followup/report.md` (document_read) · 줄 5, 14, 40, 44, 50 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/report.md` (document_read) · 줄 5, 12, 19, 35, 45, 73, 79, 79, 81, 92, 92 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-test/report.md` (document_read) · 줄 5, 20, 22, 23, 24, 27, 35, 37, 61 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FINAL-review/report.md` (document_read) · 줄 41, 129, 137, 143 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FIX-review/report.md` (document_read) · 줄 43, 48, 148, 157, 199 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FINAL.md` (document_read) · 줄 5, 25, 26, 76 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX.md` (document_read) · 줄 5, 24, 34, 75 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md` (document_read) · 줄 5, 16, 23, 82, 108 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-UI.md` (document_read) · 줄 5, 27 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER.md` (document_read) · 줄 5, 27, 29, 108 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-FIX-review/report.md` (document_read) · 줄 14, 23, 27, 33, 34, 35, 37, 48, 50 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-TESTER.md` (document_read) · 줄 5, 28, 102, 104 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-review/report.md` (document_read) · 줄 14, 34, 44, 45, 46, 47, 48, 51 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-DEV-TESTER.md` (document_read) · 줄 5, 23, 25, 46 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-DEV-review/report.md` (document_read) · 줄 17, 23, 39 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-TESTER.md` (document_read) · 줄 5, 34, 59 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-review/report.md` (document_read) · 줄 23, 35, 36, 41, 42, 52, 58 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TESTER.md` (document_read) · 줄 5, 15, 29, 106 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER.md` (document_read) · 줄 5, 15, 25, 33, 34 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-review/report.md` (document_read) · 줄 22, 42, 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-review/report.md` (document_read) · 줄 15 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PREP-002-review/report.md` (document_read) · 줄 12, 16, 24, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-FINAL-RECORDS-review/report.md` (document_read) · 줄 14, 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-FIX-TESTER.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 16, 24, 66, 100 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-TESTER-test/probe-failures.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-TESTER.md` (document_read) · 줄 5, 17, 128, 133 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/scenario.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-FIX-TESTER.md` (document_read) · 줄 5, 29 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-QA-RECORD-REVIEW-review/report.md` (document_read) · 줄 19, 29, 37 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW-review/report.md` (document_read) · 줄 15, 37, 40, 90, 102, 112, 137 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-REVIEW-review/report.md` (document_read) · 줄 15, 119, 129 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TESTER.md` (document_read) · 줄 5, 19 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW-review/report.md` (document_read) · 줄 27, 45, 66 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-UI.md` (document_read) · 줄 5, 45 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER.md` (document_read) · 줄 5, 23, 53, 56, 59, 64 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 17, 23, 107, 117 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER-test/probe-failures.md` (document_read) · 줄 5, 26 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER.md` (document_read) · 줄 5, 27, 69 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX-2.md` (document_read) · 줄 5, 62, 63, 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-RECORD-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-SERVICE-001-REVIEW-review/report.md` (document_read) · 줄 12, 20, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-DEV-test/report.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-FINAL-review/report.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-TESTER-test/report.md` (document_read) · 줄 5, 18, 34, 79 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-CLEANUP/report.md` (document_read) · 줄 16, 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-CODE-REVIEW-review/report.md` (document_read) · 줄 24, 72, 79 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/adapter-report.md` (document_read) · 줄 36 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/audit.md` (document_read) · 줄 20, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/remote-report.md` (document_read) · 줄 17 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-FINAL-MAIN-REVIEW-review/report.md` (document_read) · 줄 24, 62, 70, 72, 78 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-FINAL-MAIN-REVIEW-review/verified-qa-report.md` (document_read) · 줄 5, 19, 98 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-INTEGRATION/report.md` (document_read) · 줄 34 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-OPS-REVIEW-review/report.md` (document_read) · 줄 19, 51, 56, 59 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-OPS/report.md` (document_read) · 줄 86 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-QA-REVIEW-review/report.md` (document_read) · 줄 14, 16, 18, 52 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-TESTER-test/report.md` (document_read) · 줄 5, 19, 98 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-RELEASE-001-OPS/report.md` (document_read) · 줄 98, 103 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REPO-DOCS-001-DEV/report.md` (document_read) · 줄 2, 10, 15, 34, 58, 59, 67, 78 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REPO-DOCS-001-REVIEW-review/report.md` (document_read) · 줄 2, 18, 24, 34 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-DESIGN-REVIEW-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-TESTER-RESULT-REVIEW-review/report.md` (document_read) · 줄 7, 10, 10, 12, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-TESTER-review/report.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL-RESULT-REVIEW-review/report.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL/report.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-RESULT-REVIEW-review/report.md` (document_read) · 줄 14, 18, 22, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-review/report.md` (document_read) · 줄 5, 12, 14, 14, 20, 26, 36, 47, 55, 56, 57, 62, 67, 70, 76 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-test/report.md` (document_read) · 줄 5, 17, 24, 45, 81, 84, 91, 93, 99 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-2-review/report.md` (document_read) · 줄 10, 14, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-review/report.md` (document_read) · 줄 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md` (document_read) · 줄 5, 14, 110, 115, 117, 119, 148, 174 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-CODE-REVIEW-2-review/report.md` (document_read) · 줄 14, 25, 48, 50, 62, 63, 72, 73 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-CODE-REVIEW-3-review/report.md` (document_read) · 줄 16, 30, 30, 50, 59, 65, 68, 69, 75 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-CODE-REVIEW-review/report.md` (document_read) · 줄 14, 25, 38, 56, 65, 66, 68 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-DESIGN/report.md` (document_read) · 줄 5, 16, 41, 43, 57 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-DEV-test/report.md` (document_read) · 줄 77, 85, 100, 100 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-RETEST-test/report.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-TESTER-test/report.md` (document_read) · 줄 5, 69, 101, 104 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-002-DEV-test/report.md` (document_read) · 줄 12, 56 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-002-TESTER-test/report.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-INTEGRATION-SOURCE-review/report.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-TESTER.md` (document_read) · 줄 5, 38, 62, 81, 92 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS-REVIEW-review/report.md` (document_read) · 줄 14, 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-review/report.md` (document_read) · 줄 5, 12, 14, 33, 41, 43 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-test/report.md` (document_read) · 줄 5, 18, 33, 37, 71, 79, 88, 96 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL-REVIEW-review/report.md` (document_read) · 줄 14, 14, 16, 18, 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL/report.md` (document_read) · 줄 5, 16, 23, 27, 49, 57, 58, 76, 81, 83, 92 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OPS-REVIEW-review/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-REVIEW-review/report.md` (document_read) · 줄 5, 20, 22, 40, 46, 48, 56, 60 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-TESTER-test/report.md` (document_read) · 줄 5, 12, 27 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-VISUAL/report.md` (document_read) · 줄 5, 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SYSTEM-DIAGRAM-001-DEV/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SYSTEM-DIAGRAM-001-REVIEW-review/report.md` (document_read) · 줄 22, 26, 36 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 14, 15, 21, 21, 21, 42, 53, 53, 59, 61, 63, 63, 65 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-DEV-test/report.md` (document_read) · 줄 2, 6, 10, 14, 14, 55, 73, 75, 92, 92, 116, 116 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-TESTER-test/report.md` (document_read) · 줄 5, 6, 10, 25, 54, 56, 56, 70 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-VISUAL/report.md` (document_read) · 줄 5, 6, 10, 22, 23, 71, 73, 77 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 6, 10, 15, 16, 17, 22, 22, 22, 28, 29, 42, 67, 71, 90, 90, 98, 102, 102, 103 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-DEV-test/report.md` (document_read) · 줄 2, 6, 10, 12, 12, 27, 45, 75, 77, 77, 109, 115, 136, 160, 163 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-FIX-REVIEW-review/report.md` (document_read) · 줄 2, 9, 14, 15, 15, 15, 21, 21, 81, 82, 88, 93, 98, 98 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-TESTER-test/report.md` (document_read) · 줄 5, 6, 10, 68, 92, 124, 136 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-VISUAL/report.md` (document_read) · 줄 5, 6, 10, 12, 21, 22, 23, 24, 68, 68, 75, 75 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6, 10, 23, 39, 39, 51, 66 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN/report.md` (document_read) · 줄 5, 6, 10, 14, 42, 42, 58 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-HOSTS-DESIGN/report.md` (document_read) · 줄 5, 6, 10, 14, 23 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 2, 9, 9, 14, 15, 15, 21, 21, 34, 36, 38, 38, 38, 40, 46, 59, 60, 63, 64, 64, 64, 65, 77, 78, 78, 84, 85, 88, 88, 89, 90, 90 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-DEV-test/report.md` (document_read) · 줄 2, 2, 6, 6, 10, 10, 12, 12, 16, 16, 16, 19, 19, 28, 44, 48, 52, 60, 60, 66, 66 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-VISUAL/report.md` (document_read) · 줄 5, 6, 6, 10, 10, 12, 23, 23, 38, 39, 41, 45, 47, 49, 56, 56, 56, 57, 64, 64, 65, 65, 66, 68, 78, 106, 107, 107, 107 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 13, 14, 14, 20, 20, 36, 54, 55, 61, 62, 63 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-DEV-test/report.md` (document_read) · 줄 2, 6, 10, 12, 12, 18, 69, 69, 71, 85, 90, 94, 96, 96, 142, 142, 147 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/report.md` (document_read) · 줄 5, 6, 10, 33, 59, 86, 117, 150, 232, 244 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-VISUAL/report.md` (document_read) · 줄 5, 6, 10, 16, 18, 18, 22, 89, 95, 96, 97, 103, 103, 109, 110, 111, 112 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 14, 15, 15, 18, 20, 20, 28, 39, 50, 52, 52, 58, 59, 60, 62, 62, 68, 68 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-FIX-test/report.md` (document_read) · 줄 6, 18, 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md` (document_read) · 줄 6, 22, 62, 62 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-RETEST-test/report.md` (document_read) · 줄 5, 6, 10, 55, 63, 64, 66, 68 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-TESTER-test/report.md` (document_read) · 줄 5, 6, 10, 87, 91 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-WINDOWS-SHELL-GUARD-001-FIX-REVIEW-review/report.md` (document_read) · 줄 18, 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-WINDOWS-SHELL-GUARD-001-REVIEW-review/report.md` (document_read) · 줄 5, 12, 47, 59 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-WINDOWS-SHELL-GUARD-001-TESTER-test/report.md` (document_read) · 줄 5, 14, 25, 32, 35, 67 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-BETA-001-TESTER-PUBLIC.md` (document_read) · 줄 5, 12 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-BETA-001-TESTER.md` (document_read) · 줄 5, 12 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-BETA-002-TESTER.md` (document_read) · 줄 5, 12 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-HUMAN-PEER-001-TESTER.md` (document_read) · 줄 5, 12, 82 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER-FINAL.md` (document_read) · 줄 5, 13 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER-FIX.md` (document_read) · 줄 5, 13, 53, 53 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER.md` (document_read) · 줄 5, 13, 41 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-AGENTS-001-FIX-TESTER.md` (document_read) · 줄 5, 14 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-AGENTS-001-TESTER.md` (document_read) · 줄 5, 13 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-IDENTITY-001-FIX-TESTER.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-IDENTITY-001-TESTER.md` (document_read) · 줄 5, 35 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-MESSAGES-001-TESTER.md` (document_read) · 줄 5, 13, 51 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-SETUP-001-TESTER.md` (document_read) · 줄 5, 20, 44, 51 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER.md` (document_read) · 줄 5, 44 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-APPROVAL-TESTER.md` (document_read) · 줄 2, 5, 6, 10, 12, 19, 42 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-CONNECT-TESTER.md` (document_read) · 줄 2, 5, 6, 10, 12, 28, 43 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-NAMES-TESTER.md` (document_read) · 줄 2, 5, 6, 10, 12, 42, 46 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.10.md` (document_read) · 줄 31, 31, 38, 58, 61, 64 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.12.md` (document_read) · 줄 24, 24, 46 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.13.md` (document_read) · 줄 16, 16 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.14.md` (document_read) · 줄 25, 25, 33, 53, 56, 59 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-099.md` (document_read) · 줄 25, 35, 37, 48, 49, 50, 59, 65 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.2.0.md` (document_read) · 줄 57 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.3.0.md` (document_read) · 줄 7, 26 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.0.md` (document_read) · 줄 28 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.2.md` (document_read) · 줄 14, 16, 39, 47, 49 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.4.md` (document_read) · 줄 30, 44, 46, 46 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-AUTO-RECEIVE-001-DEV.md` (document_read) · 줄 114 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 97 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-CODEX-HOOK-QUOTE-001-TESTER.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-DEPLOY-001-OPS.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-FRONTEND-POLISH-001-DEV.md` (document_read) · 줄 46, 80, 80 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-FRONTEND-POLISH-001-VISUAL.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-002-DEV.md` (document_read) · 줄 17, 21 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md` (document_read) · 줄 37 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-DEV.md` (document_read) · 줄 34, 77 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-TESTER.md` (document_read) · 줄 5, 35, 36 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-DEV.md` (document_read) · 줄 37, 37, 43, 43, 56, 82, 82 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-VISUAL.md` (document_read) · 줄 5, 14, 28 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001.md` (document_read) · 줄 5, 16, 22, 26, 65, 78, 79, 89, 95, 95, 104, 108, 110, 118, 120 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001-DEV.md` (document_read) · 줄 20, 20, 63, 63 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001-TESTER.md` (document_read) · 줄 5, 14, 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001-VISUAL.md` (document_read) · 줄 5, 16, 18, 28, 40, 42, 51 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md` (document_read) · 줄 2, 5, 14, 38, 59, 59, 65, 67 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-001-DEV.md` (document_read) · 줄 85, 89, 89, 135 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-001-UI.md` (document_read) · 줄 5, 14, 17, 17, 35 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV-FIX.md` (document_read) · 줄 33, 44, 57, 132 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md` (document_read) · 줄 59, 150 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-DEV-TESTER.md` (document_read) · 줄 5, 15 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md` (document_read) · 줄 43 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT.md` (document_read) · 줄 70 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-PUBLIC-POLICY-001.md` (document_read) · 줄 5, 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PREP-002.md` (document_read) · 줄 5, 28, 32, 44, 44, 70, 72 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md` (document_read) · 줄 16, 38, 63, 67, 92, 97 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX.md` (document_read) · 줄 90 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 35, 60, 73 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` (document_read) · 줄 5, 16, 61, 97, 104, 110 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 5, 32 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX.md` (document_read) · 줄 61, 69, 83, 136, 137 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG.md` (document_read) · 줄 54, 76 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-FIX-TESTER.md` (document_read) · 줄 5, 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-TESTER.md` (document_read) · 줄 5, 40 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW.md` (document_read) · 줄 35 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-UI.md` (document_read) · 줄 5, 16, 16, 18, 56 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2.md` (document_read) · 줄 25, 27, 38, 58, 93 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-3.md` (document_read) · 줄 15, 17, 31, 37, 88 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX.md` (document_read) · 줄 7, 15, 76 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-UI-FIX-2.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 5, 17, 40, 56, 62 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-UI.md` (document_read) · 줄 5, 38, 46 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (document_read) · 줄 5, 14, 36, 49, 49, 50, 50, 62, 80 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 14, 44, 55, 107, 111, 113, 116, 220 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-001-DEV.md` (document_read) · 줄 64 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-002-DEV.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REFACTOR-002-TESTER.md` (document_read) · 줄 5, 12 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-REPO-DOCS-001-DEV.md` (document_read) · 줄 7, 33 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md` (document_read) · 줄 5, 17, 45 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 24, 33, 39, 53, 87, 101, 101, 105 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-DEV.md` (document_read) · 줄 54 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-TESTER.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-VISUAL.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-TESTER.md` (document_read) · 줄 5, 16, 24, 38, 39, 40 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-VISUAL.md` (document_read) · 줄 5, 14, 37 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SESSION-PEER-001-DESIGN.md` (document_read) · 줄 5, 12 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SESSION-PEER-001-DEV.md` (document_read) · 줄 25 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SETUP-001-TESTER.md` (document_read) · 줄 5, 17, 36 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SETUP-001.md` (document_read) · 줄 5, 17, 25, 36, 37, 39, 62, 67 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-DEV.md` (document_read) · 줄 16, 25, 28 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-INSTALL.md` (document_read) · 줄 16, 102 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV.md` (document_read) · 줄 39, 39, 69 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER.md` (document_read) · 줄 5, 21, 33, 34 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL.md` (document_read) · 줄 5, 35 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OPS.md` (document_read) · 줄 21 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-REMOTE.md` (document_read) · 줄 18, 66 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-DESIGN.md` (document_read) · 줄 5, 6, 10, 14, 16, 38, 42 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-NAMES-DEV.md` (document_read) · 줄 6, 7, 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-NAMES-TESTER.md` (document_read) · 줄 5, 6, 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-RECEIPT-DEV.md` (document_read) · 줄 6, 20 · inferred
- `.fullops-squad/docs/generated/db-schema.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_read, document_update) · 줄 7, 7, 13, 13, 15, 15, 15, 15, 17, 19, 19, 19, 21, 23, 25, 25, 25, 27, 29, 29, 29, 87, 94, 129 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read, document_update) · 줄 7, 7, 24, 35, 45, 53, 57, 67, 82, 87, 93, 129, 137, 235 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 5, 14, 23, 28, 28, 29, 49, 49, 50, 51, 51, 52, 53, 54, 74, 87, 91 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 6, 63 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 6, 30, 34, 34, 38, 82, 125 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 6, 79 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 6, 29, 131 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 6, 70, 94, 94, 95, 96, 97 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md` (document_read) · 줄 5, 137 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/product.md` (document_read) · 줄 5 · inferred
- `.fullops-squad/handovers/SAR-MVP-001-REVIEW.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/handovers/_TEMPLATE.md` (document_read) · 줄 26 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_designer.md` (document_read) · 줄 2, 5, 16, 24, 27, 44, 51, 74, 76, 92, 99, 101, 121, 123, 138, 146, 154, 163, 177, 187, 197, 213, 226, 235, 243, 257, 266, 268, 268, 274 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_ops.md` (document_read) · 줄 37, 70, 156, 189, 205, 255 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_dev.md` (document_read) · 줄 33, 37, 49, 67, 100, 115, 126, 129, 153, 163, 176, 179, 207, 277 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_ops.md` (document_read) · 줄 29, 40, 76, 99, 124, 182, 198, 212, 215, 233, 275 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_tester.md` (document_read) · 줄 2, 5, 15, 29, 31, 42, 65, 79, 81, 115, 122, 160, 166, 200, 207, 208 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_designer.md` (document_read) · 줄 2, 5, 53, 57, 79, 89, 101, 109, 131, 156 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_dev.md` (document_read) · 줄 28, 32, 43, 45, 51, 88, 125, 142, 147, 154, 181, 196, 197, 202, 207 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_ops.md` (document_read) · 줄 105, 136 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_tester.md` (document_read) · 줄 2, 5, 15, 28, 32, 43, 43, 49, 79, 86, 88 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 2, 5, 15, 34, 37, 39, 50, 86, 93, 94, 100, 104, 104, 156, 166, 205, 215, 223, 231, 234, 244, 283, 284, 285, 301, 306, 324, 373, 381, 386, 407, 441 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 33, 71, 74, 106, 110, 133, 139, 145, 175, 178, 180, 203, 206, 207, 212, 241, 266, 283, 304, 357, 377 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_tester.md` (document_read) · 줄 2, 5, 15, 22, 23, 34, 38, 80, 99, 102, 104, 115, 142, 150, 158, 161, 171, 223, 228, 246, 252 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_designer.md` (document_read) · 줄 2, 5, 16, 28, 29, 34, 59, 79, 85, 88, 215, 218, 219, 275, 275 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 24, 24, 37, 67, 96, 101, 122, 125, 139, 152, 158, 255, 258, 305, 317 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_tester.md` (document_read) · 줄 2, 5, 16, 32, 36, 45, 61, 78, 84, 87, 212, 213, 257, 287, 290, 293, 298 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_dev.md` (document_read) · 줄 39, 53, 73, 81, 141, 218, 221, 270, 276, 281, 288, 291, 304, 306, 310, 311, 313, 327, 331 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_ops.md` (document_read) · 줄 44, 108, 112, 117, 120, 215, 267 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_tester.md` (document_read) · 줄 2, 5, 16, 30, 36, 47, 56, 73, 85, 90, 101, 123, 130, 131, 234, 320, 347, 399 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_designer.md` (document_read) · 줄 2, 5, 16, 27, 29, 31, 39, 43, 51, 54, 57, 70, 114, 118, 126, 157, 169, 179, 190, 191, 192, 203, 214, 217, 232, 245, 255, 269, 275, 276, 333, 337, 337 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_dev.md` (document_read) · 줄 34, 48, 56, 62, 74, 112, 126, 133, 147, 173, 182, 185, 320, 321, 326, 355, 375, 379, 395, 400, 407, 425, 425, 443, 443, 456, 460, 461, 462, 490, 508, 512, 516 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_dev.md` (document_read) · 줄 38, 46, 47, 58, 58, 69, 88, 95, 96, 134, 136, 137, 140, 144, 152, 156, 185, 191, 199, 205 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_tester.md` (document_read) · 줄 2, 5, 16, 44, 103, 109, 110, 111, 113, 292, 300, 304, 315, 321, 349, 349, 369, 384 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_designer.md` (document_read) · 줄 2, 5, 16, 67, 69, 69, 73, 90, 97, 100 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_ops.md` (document_read) · 줄 30, 30, 55, 74, 84, 123, 125, 130, 139, 159, 163, 197, 202, 204, 207, 317, 318, 319, 320, 321, 396, 435, 437, 438, 498 · inferred
- `.fullops-squad/handovers/logs/2026-10-18_to_dev.md` (document_read) · 줄 48, 70, 97, 115, 119, 121, 198, 199, 202, 242, 271, 273, 278, 297 · inferred
- `.fullops-squad/handovers/logs/2026-10-19_to_designer.md` (document_read) · 줄 2, 5, 9, 10, 13, 16, 17, 25, 30, 33, 34, 40, 42, 42, 43, 44, 48, 52, 52, 52, 61, 64, 92, 94, 98, 102, 103, 105, 137, 138, 139, 140, 141, 142, 143, 158, 159, 199, 203, 204, 206, 207, 209, 210, 213, 216, 217, 225, 230, 233, 234, 244, 244, 249, 254, 258, 258, 263, 263, 266, 281, 300, 302, 305, 309, 310, 312, 340, 341, 342, 343, 344, 345, 346, 347, 348, 362, 363, 389, 392, 393, 394, 395, 395, 396 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-DEV-REVIEW.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-INTEGRATION-REVIEW.md` (document_read) · 줄 17, 24 · inferred
- `.fullops-squad/handovers/to_tester.md` (document_read) · 줄 2, 2, 5, 6, 6, 14, 14, 19, 22, 23, 31, 31, 31, 31, 45, 49, 53, 53, 59, 59, 61, 61, 89, 89, 91, 94, 98, 99, 101, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 218, 234, 235, 236, 273, 274, 275, 276, 312 · inferred
- `.fullops-squad/lint/README.md` (document_read) · 줄 86 · inferred
- `.fullops-squad/orca-agents.md` (document_read) · 줄 21, 21, 26, 29, 31, 41, 44, 50, 59, 60, 68, 69, 70, 77, 78, 79, 82, 83, 85, 86, 111, 113, 121 · inferred
- `README.md` (document_read) · 줄 20, 21, 26, 167, 195 · inferred
- `adapters/README.md` (document_read) · 줄 16, 77, 126, 271 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 12, 14 · inferred
- `docs/README.md` (document_read) · 줄 12, 19, 20, 21, 24 · inferred
- `docs/development.md` (document_read) · 줄 4 · inferred
- `docs/operations.md` (document_read) · 줄 82 · inferred
- `docs/system-design.md` (document_read) · 줄 14, 70 · inferred
- `internal/relay/admission_integration_test.go` (impact_check) · 줄 3, 19, 74, 100, 146, 156 · inferred
- `internal/relay/api_rate.go` (impact_check) · 줄 1 · inferred
- `internal/relay/capacity.go` (impact_check) · 줄 1, 13, 137, 156, 160, 167, 222 · inferred
- `internal/relay/cleanup_flood_integration_test.go` (impact_check) · 줄 3, 19, 49, 68, 135, 178 · inferred
- `internal/relay/cleanup_unit_integration_test.go` (impact_check) · 줄 3, 19, 93, 150 · inferred
- `internal/relay/human_http.go` (impact_check) · 줄 97, 449, 522 · inferred
- `internal/relay/human_pages.go` (impact_check) · 줄 10, 15, 20, 21, 22, 22, 23, 27, 27, 33, 33, 38, 42, 43, 44, 55, 55, 57, 57 · inferred
- `internal/relay/i18n.go` (direct_edit) · 줄 22, 26, 40, 50, 58, 84, 120, 153, 166, 174, 191, 201, 250, 269 · inferred
- `internal/relay/i18n_en_human.go` (impact_check) · 줄 48, 50, 92, 135, 137, 143, 145, 152, 248, 314 · inferred
- `internal/relay/i18n_test.go` (impact_check) · 줄 48, 137, 266, 272, 282, 285, 285, 327, 332, 333, 360, 368 · inferred
- `internal/relay/landing_pages.go` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/member_pages.go` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/public_messages_integration_test.go` (impact_check) · 줄 272, 275, 290, 429, 445, 451, 474 · inferred
- `internal/relay/text_events.go` (impact_check) · 줄 36, 82, 109, 133 · inferred
- `internal/relay/text_events_integration_test.go` (impact_check) · 줄 51 · inferred
- `scripts/package_plugin.py` (impact_check) · 줄 15, 16, 155 · inferred
- `videos/knowslink-motion/assets/fonts/Geist-OFL.txt` (impact_check) · 줄 43 · inferred
- `videos/knowslink-motion/assets/fonts/Pretendard-OFL.txt` (impact_check) · 줄 45 · inferred
- `videos/knowslink-motion/index.html` (impact_check) · 줄 171, 181 · inferred
- `videos/knowslink-motion/render.mjs` (impact_check) · 줄 11 · inferred
미확인 26건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
