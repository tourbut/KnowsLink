---
title: SAR-SCHEDULE-001-DEV — Google Calendar 빈 시간 조회 연동(SAR-SCHEDULE-001 SC-01~08)을 구현한다
status: draft
updated: 2026-10-20
owner: dev
tasks: [SAR-SCHEDULE-001-DEV]
summary: Google Calendar 빈 시간 조회 연동(SAR-SCHEDULE-001 SC-01~08)을 구현한다
attempt: 34301b9d99c149c7acc9383d75a7155e
base: 826ea80be09c09d67fda75f39fd6a7e7fd32f494
subagent_level: standard
test_level: lite
---

# SAR-SCHEDULE-001-DEV — Google Calendar 빈 시간 조회 연동(SAR-SCHEDULE-001 SC-01~08)을 구현한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: SAR-SCHEDULE-001 확정 규칙의 기술 계획·구현·테스트·기술 문서. 캘린더 연결 상태와 암호화 토큰 보관을 위한 상태 저장 확장, 회원 화면의 '캘린더 연결/해제' 경로, Google OAuth 증분 동의·freebusy 호출 코드, 합성 제공자 fixture는 승인됨. 실제 Google 계정 연결·실데이터 조회·OAuth 클라이언트 설정 변경(Cloud Console)·운영 배포·DB migration·새 외부 의존성·권한 확대는 coor에 ask(D011 열림).
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

사용자 결정 D008·D009·D010(2026-10-20)으로 실제 일정 API 연동을 구현한다. 확정 정본은 `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md`(SD-01~06 모두 A, 확정 수치, 확정 제품 규칙, 완료 조건 SC-01~SC-09, 범위 밖, DEV 인계 기준). 결정 근거와 공식 문서 출처 E1~E10: `.fullops-squad/docs/evaluations/qa-reports/SAR-SCHEDULE-001-DESIGN/report.md`.

요약: Google Calendar만, KnowsLink 서버가 refresh token을 암호화 보관하고 서버에서 빈 시간(freebusy)만 조회, 결과는 바쁨 구간을 요청 단위(15·30·60분, 기본 30)로 반올림해 공개, 최대 14일, 친구당 24시간 5회, 매 요청 사람 승인 = 그 요청 한 번의 조회와 한 번의 공개, 내 계정 '캘린더 연결' 별도 증분 동의·같은 화면 해제. `schedule.commit`은 계속 stub. 운영은 2026-10-20부터 제품 7b4f9bc.

D011(열림): Google Cloud Console OAuth 앱 게시 상태와 calendar.freebusy 민감 범위 심사 여부는 사용자가 나중에 확인한다. 이 과제는 운영 연결 없이 합성 제공자와 격리 환경으로 완료한다.

## 적용 기준과 예외

- 규칙: fullops-common-0.3.3(특히 security: 토큰 암호화·비밀 비출력·로그 최소화), `.fullops-squad/project.md`, 문서 규칙. 기준 main(이 지시서 커밋).
- 외부 API(Google OAuth 증분 동의, Calendar freebusy, 토큰 갱신·철회)는 Context7 `resolve-library-id` → `query-docs`로 확인하고 실패하면 공식 문서로 확인한다. 근거(라이브러리 ID·버전·URL)를 DEV 보고에 남긴다. 새 외부 Go 의존성은 쓰지 않는 쪽을 우선하고 필요하면 ask.
- 암호화 키·OAuth 클라이언트 비밀은 환경 변수로만 받고 저장소·로그·화면에 넣지 않는다. 기존 Google 로그인 설정 방식과 맞춘다.
- Windows make 없음: 이전 단위와 같이 격리 원격 Linux clone에서 make lint/test/integration, fixture 실제 기동. lint.py는 깨끗한 임시 detached worktree에서 실행하고 force 없이 정리.

## 먼저 읽을 문서

SAR-SCHEDULE-001.md 전체, 결정 근거 report, `.fullops-squad/docs/design-docs/architecture.md`(gate·policy·disclosure), `interface-design.md`·`data-model.md`, U5 승인 요약 코드(`internal/relay/gate_summary.go`, `http.go` gateView/gate 처리), Google 로그인 코드(Google OAuth 처리 파일), 회원 설정 화면(`member_pages.go`), `schedule.query` 정책·스키마, `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md`, `.fullops-squad/contexts/dev.md`.

## 해야 할 일과 파일 소유권

- [ ] 짧은 기술 계획(완료 보고): 토큰 암호화·키 관리, 연결 상태 저장, 증분 동의 흐름, freebusy 호출·반올림, 승인 1회 실행 보장(멱등), 친구당 횟수 집계, 실패 분류(denied/failed), 철회 처리, 합성 제공자 구조.
- [ ] SC-01~SC-08을 구현한다(SC-09는 수락 절차).
- [ ] 각 SC의 실패 경계 자동 검사(합성 제공자): 승인 1회=조회 1회(중복·재접근·재시작), 공개 필드 최소화(제목·참석자·위치·설명·캘린더 ID·이메일이 결과·로그·화면에 없음), 범위 초과 거부, 철회·친구 해제·차단 직후 무조회, 실패 구분, 미연결 회원 승인 버튼 없음, commit stub 유지.
- [ ] tester/designer용 격리 fixture(합성 제공자, 연결/미연결 회원, 친구 쌍, 승인 대기·만료·처리됨, 바쁨 데이터) 준비·실제 기동 확인.
- [ ] 영향받는 기술 문서(D03·D05·D09·D10·D11·D13 중 실제 변경) 갱신, 운영 연결에 필요한 환경 변수·OAuth scope·D011 의존을 운영 문서 초안으로 남긴다(실행은 하지 않음).
- [ ] 원격 Linux make lint/test/integration, FullOps lint.py --from (지시서 base), DEV 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-SCHEDULE-001-DEV-test/report.md`, packet-outcomes, finish, 커밋, worker_done.

소유 파일: `internal/relay/` 일정·캘린더 연결·gate 실행 경로와 테스트·fixture, 회원 설정의 캘린더 연결 절, 기술 문서, own 기록·DEV 보고. 예상 규모가 크다. SIZE 경고는 분할하지 않은 이유와 함께 보고한다.

## 완료 기준과 검증

SC-01~SC-08이 자동 검사와 실제 기동 fixture로 확인되고, 기존 gate·친구·승인 요약 회귀 없음. 원격 Linux make lint/test/integration EXIT=0, FullOps lint ERROR 0. worker 완료는 제품 수락이 아니다(리뷰·tester·designer·실계정 1경로 확인(D011 뒤)·통합은 coor).

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04, SAR-SCHEDULE-001.md '화면 문구 원칙'.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 기존 설정 화면 절·승인 상세 구조, 두 언어 i18n 구조(새 문구는 영어판 포함). 새 테마·JS 없음.
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(임시 worktree).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 테마 전환 없음. 바뀐 화면 두 언어 3폭 렌더 넘침 0.

## 갱신할 산출물

route 추천과 실제 변경 기준으로 D03·D05·D09·D10·D11·D13 중 해당 원천.

탐색 패킷의 task/role/attempt·기준 SHA와 direct_edit/impact_check/document_read/document_update 경로·구간·근거·partial/unknown을 연결한다. 명시적 없음과 미확인을 구분한다. 필수·충돌·주의 후보와 budget 잔여 확인 목록을 보존한다.

## 기대 산출물

제품 코드·테스트·fixture, 기술·운영 문서 초안, DEV 보고(기술 계획·외부 API 근거·검사·렌더·fixture·not_run), packet-outcomes, 완료 보고.

## 제약·협업·후속

하위 작업은 [위임 기준](../rules/delegation.md)을 따른다(standard, 최대 2). 미추적 `videos/service-intro/`·`qa-reports/SAR-MOTION-INTRO-001-DEV-test/` 하위 미추적 폴더는 건드리지 않는다. workspace 밖에 파일을 만들지 않는다(private 도구 디렉터리 제외). 다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다. 피할 스타일: 상태 카드·통계 배지·캘린더 그리드 위젯·JS 드롭다운·새 색 토큰.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
