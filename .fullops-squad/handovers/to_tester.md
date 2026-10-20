---
title: SAR-UX-ISSUES-002-LANG-TESTER — 언어·설정(7704e6d)을 격리 fixture에서 독립 QA한다
status: draft
updated: 2026-10-20
owner: tester
tasks: [SAR-UX-ISSUES-002-LANG-TESTER]
summary: 언어·설정(7704e6d)을 격리 fixture에서 독립 QA한다
attempt: 4b1aa18f07504d9fb2e6e091242ea605
base: 08f26982f8ccceb4921b5d82462f57411fb2a112
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-LANG-TESTER — 언어·설정(7704e6d)을 격리 fixture에서 독립 QA한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → tester
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 격리 fixture 사용 검증, 시나리오·검사·보고·캡처 작성과 tester 브랜치 커밋. 제품 코드 수정, 새 의존성은 하지 않는다.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/tester / fullops/tester
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

- [ ] 시나리오 `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-LANG-TESTER.md`(covers·독립 checks·origin·허용 동작·합성 en/ko/none·관측 수준).
- [ ] 설정 화면: 헤더 '설정', 순서(내 이름·내 친구 코드·언어·로그인과 보안), '한국어 / English'·'언어 저장 / Save language', 저장 성공 뒤 같은 위치 선택 언어, 거부(잘못된 값) 시 이전 저장값과 선택 유지.
- [ ] 선호 우선순위: 저장값 → 비로그인 명시 선호(`/lang/{ko|en}` 쿠키) → 한국어, 계정 교대(같은 브라우저) 시 이전 계정 선호 미복사, 첫 로그인 시 비로그인 선호는 화면만 적용·저장은 명시 저장, 손상 locale fallback, 쿠키 차단 시 현재 응답만 적용·안내, 외부 back 거부.
- [ ] 지정 화면 두 언어: 공개·로그인(`/auth/verify`)·이름·대화·친구·내 에이전트·설정·legacy·조회·승인·local 연결·remote 동의·복귀를 en/ko로 열어 HTML lang 일치, 영어판에 한글 잔존 없음(원문·ID·이름·명령 제외), 접근성 이름 번역, KST 기한 문자열 동일, 원문/ID 불변.
- [ ] 언어 전환이 친구 요청·발송·승인을 실행하거나 입력을 잃지 않음(전후 상태 수 비교).
- [ ] 키보드·focus·44px, 1440·390·320 두 언어 넘침 0(긴 영문). 스크린리더·실제 200%는 가능하면 시도, 못 하면 not_run.
- [ ] 보고서 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-TESTER-test/report.md`, designer용 두 언어 3폭 캡처, 완료 보고 → work.py finish → 커밋 → worker_done.

소유 파일: 시나리오, `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-TESTER-test/` 하위, own route/inbox/log.

## 완료 기준과 검증

각 check PASS/FAIL/not_run, 결함의 재현·기대·실제·SHA·원본 로그. critical/high는 완료 보고에 명확히 표시. 검사를 약화하지 않는다.

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
