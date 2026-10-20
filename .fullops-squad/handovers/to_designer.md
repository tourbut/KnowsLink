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

- [ ] 설정 화면(한국어·영어), 공개 첫페이지 두 언어와 '한국어 / English' 링크 위치, 주요 회원 화면(친구·내 에이전트·대화·승인·조회·연결)의 영어판을 1440·390·320으로 직접 판정(PASS/FAIL/not_run, 근거 캡처). DEV 렌더 `render/`도 참고한다.
- [ ] 정본 U3: 설정 순서·문구, 언어 선택 형태(국기·드롭다운 JS 없음), 공개 페이지 정보 구조·설치 사실 불변, 영어판이 같은 의미·기능이며 새 약속을 하지 않는지, 영어 문구 품질(표본: 버튼·상태·오류·빈 상태·도움말), 긴 영문 줄바꿈, D04 토큰·44px·왼쪽 정렬.
- [ ] 제품 판단: DEV 알려진 한계 중 admission 429/503 페이지가 회원 저장값이 아니라 브라우저 쿠키 언어를 따르는 것을 허용할지. 영어 용어집 수정이 필요하면 목록으로 남긴다(DEV 수정 항목).
- [ ] 캡처는 fixture 허용 GET으로 직접 만든다(언어 저장 POST는 하지 않는다; `en`·`ko`·`none` 계정과 `/lang/{ko|en}` GET으로 언어를 바꾼다). tester 캡처(`C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-TESTER-test/`)는 생기면 참고한다.
- [ ] 판정 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-VISUAL/report.md` → 완료 보고 → work.py finish → 커밋 → worker_done(결론과 DEV 수정 필요 목록).

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

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
