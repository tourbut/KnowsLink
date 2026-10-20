---
title: SAR-UX-ISSUES-002-LANG-TESTER — DEV 수정 b11d5b4(제품 665895a)의 admission 언어·링크 대체·문구를 재확인한다
status: draft
updated: 2026-10-20
owner: tester
tasks: [SAR-UX-ISSUES-002-LANG-TESTER]
summary: DEV 수정 b11d5b4(제품 665895a)의 admission 언어·링크 대체·문구를 재확인한다
attempt: a59427218e9045169be0b418572ba44f
base: 08f26982f8ccceb4921b5d82462f57411fb2a112
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-LANG-TESTER — DEV 수정 b11d5b4(제품 665895a)의 admission 언어·링크 대체·문구를 재확인한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → tester
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: fixture 재확인 검사·보고 갱신과 tester 브랜치 커밋. 제품 코드 수정, 운영·실계정·외부 발송, 삭제·이동, force push, 새 의존성은 하지 않는다. workspace 밖에 파일을 만들지 않는다.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/tester / fullops/tester
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

같은 과제의 재확인 attempt다. 직전 결과 cca284e(전 check PASS, medium D-LG-01, low O-LG-01~04). DEV 수정 `b11d5b49bddb94f97461d7f33be543f9d003fbc5`(fullops/dev, 제품 `665895a`): admission 429/503이 브라우저 쿠키 언어·Content-Language를 따름(D-LG-01·리뷰 F1), 한국어 README 링크 `hreflang=ko`와 '(in Korean)'(F2), 저장값 회원의 공개·remote 동의 화면 언어 링크를 '설정에서 변경'으로 대체(V-LANG-01·F3), 영어 'your' 문구(V-LANG-02), 설정 '미저장' 표시(V-LANG-03), details 장식 문자 접근성 이름 제외(O-LG-03). Linux make lint/test·전체 integration EXIT=0. coor가 수정분 diff를 확인했다(need skip). designer 재검수는 필요 없다(designer 판정). DEV 보고 수정 절: `C:/Users/shin/orca/workspaces/KnowsLink/dev/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-DEV-test/report.md`.

## fixture

origin `http://127.0.0.1:18933`(coor SSH 포워딩, b11d5b4로 2026-10-20 재기동, 초기 상태). 서버 clone `/tmp/knowslink-lang002-b11d5b49bddb-db564d08`, container `kl-lang002-fixture-b11d5b49bddb-66cf1e77`. manifest `C:/Users/shin/AppData/Local/KnowsLinkDevTools/lang-ux002-dev/fixture-manifest.json`. 계정 `en`·`ko`·`none`, 허용 경로는 이전 attempt와 같다.

## 적용 기준과 예외

이전 attempt와 같다. 제품 정본 SAR-UX-ISSUES-002.md U3(designer 갱신: admission은 쿠키 언어).

## 먼저 읽을 문서

DEV 보고 수정 절, 이전 완료 기록, 기존 tester 보고 D-LG-01·O-LG 절, designer VISUAL 보고 V-LANG-01~03.

## 해야 할 일과 파일 소유권

- [ ] D-LG-01 재판정: 회원 요청 한도 429(및 재현 가능하면 503) 화면이 쿠키 언어(en/ko)를 따르고 `Content-Language`·html lang 일치. 기존 FAIL 검사의 기대를 낮추지 않는다.
- [ ] 저장값 회원(`en`·`ko`)의 공개 첫페이지·remote 동의 화면 언어 표시가 '설정에서 변경' 링크로 대체되고 동작하는지, 저장값 없는 `none`과 비로그인은 기존 `/lang/` 링크가 동작하는지.
- [ ] 영어 화면의 한국어 문서 링크에 `hreflang=ko`와 언어 안내.
- [ ] 설정 '미저장' 표시(`none`), 영어 문구 표본(designer 목록 중 일부), details 요약 접근성 이름에 장식 문자 없음(가능하면 AX 측정).
- [ ] 이전 PASS 핵심(설정 순서·우선순위·계정 교대·두 언어 lang·넘침 0) 빠른 회귀.
- [ ] 보고 재확인 절 → 완료 보고 → work.py finish → 커밋 → worker_done.

소유 파일: `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-TESTER-test/` 하위, 시나리오, own route/inbox/log.

## 완료 기준과 검증

각 항목 PASS/FAIL/not_run과 근거, 새 결함의 재현·기대·실제·SHA. 검사를 약화하지 않는다.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 해당 없음(제품 코드 수정 없음).
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 테마 전환 없음.

## 갱신할 산출물

없음.

탐색 패킷의 task/role/attempt·기준 SHA와 direct_edit/impact_check/document_read/document_update 경로·구간·근거·partial/unknown을 연결한다. 명시적 없음과 미확인을 구분한다. 필수·충돌·주의 후보와 budget 잔여 확인 목록을 보존한다.

## 기대 산출물

보고 재확인 절, 결과·캡처, 완료 보고.

## 제약·협업·후속

제품 코드·다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다.

## 이전 완료 기록

- .fullops-squad/handovers/logs/2026-10-20_to_tester.md

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
