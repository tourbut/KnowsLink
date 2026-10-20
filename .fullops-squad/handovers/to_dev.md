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

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
