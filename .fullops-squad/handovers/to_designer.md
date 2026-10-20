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

- [ ] '어디에서 사용할까요?' 5선택지와 각 선택 뒤 안내, 처음/기존 연결, 상태별 재개 행, 보낸 요청 '입력한 코드 …abcd', remote 앱 복귀·실패 안내, legacy '이 화면에서 찾기'를 1440·390·320으로 직접 판정(PASS/FAIL/not_run, 근거 캡처). DEV 렌더 `render/`도 참고한다.
- [ ] 정본 문구·순서, 대상 외 호스트 미표시, '온라인' 배지 없음, D04 토큰·44px·왼쪽 정렬·피할 스타일 준수.
- [ ] 제품 판단: 차단된 요청이 만료된 뒤 같은 코드 재요청은 새 행, 없는 코드는 기존 행 재사용 — 이 차이를 허용할지, 같은 모양으로 맞출 기준(예: 둘 다 새 행, 둘 다 재사용)을 정한다. 결정을 SAR-UX-ISSUES-002.md U4에 반영.
- [ ] 캡처는 fixture 허용 GET으로 직접 만든다. 상태를 바꾸는 POST는 하지 않는다. tester 캡처(`C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-TESTER-test/`)는 생기면 참고한다.
- [ ] 판정 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-VISUAL/report.md` → 완료 보고 → work.py finish → 커밋 → worker_done(결론과 DEV 수정 필요 목록).

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

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
