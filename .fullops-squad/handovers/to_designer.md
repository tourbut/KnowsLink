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

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
