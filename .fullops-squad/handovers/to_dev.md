---
title: SAR-UX-ISSUES-002-APPROVAL-DEV — 격리 승인 fixture가 실제 기동 중 POST /v1/claim 409 capacity로 실패하는 문제를 고친다
status: draft
updated: 2026-10-19
owner: dev
tasks: [SAR-UX-ISSUES-002-APPROVAL-DEV]
summary: 격리 승인 fixture가 실제 기동 중 POST /v1/claim 409 capacity로 실패하는 문제를 고친다
attempt: 21e8bb63fedb4af787034cabc2d95b54
base: 688305054dc4f3570f0c9f0d1e7917c35a19db3a
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-APPROVAL-DEV — 격리 승인 fixture가 실제 기동 중 POST /v1/claim 409 capacity로 실패하는 문제를 고친다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-19
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 승인 fixture 테스트 코드 수정과 실제 기동 확인(원격 격리 clone·자기 Postgres 컨테이너, 기존 remote.py fixture). 제품 동작·정책·한도 값 변경이 필요하면 coor에 ask.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

같은 과제의 두 번째 attempt다. 후보 ba047b7(제품 f43d05c)에서 coor가 `remote.py fixture`로 기동했더니 fixture가 바로 종료됐다. 원격 `/tmp/knowslink-approval002-ba047b74581f-195bf370/fixture.log`:

```
=== RUN   TestApprovalBrowserFixture
    approval_browser_fixture_test.go:59: POST /v1/claim got 409 {"error":"capacity"}
         want 200
--- FAIL: TestApprovalBrowserFixture (0.49s)
```

59행은 gate spec 루프의 `claim := f.deliver(parent)`다. 이전 attempt는 컴파일만 확인했다(not_run). 제품 코드는 고정 SHA 리뷰 전이다.

## 적용 기준과 예외

이전 attempt와 같다. 원인이 fixture 구성(예: 수신 한도·claim 순서)이면 테스트 코드만 고친다. 제품 한도나 동작을 바꿔야 하면 멈추고 ask한다.

## 먼저 읽을 문서

`internal/relay/approval_browser_fixture_test.go`, claim/capacity 경로(`internal/relay/capacity.go`, `connections.go`, `device.go` 등 실제 409 capacity 발생원), 이전 완료 기록.

## 해야 할 일과 파일 소유권

- [ ] 409 capacity 원인 분석과 fixture 수정.
- [ ] `remote.py fixture`로 실제 기동해 manifest(`fixture.json`)가 쓰이고 loopback origin이 HTTP 200을 주는지 확인(원격 fixture.log 근거). 확인 뒤 프로세스는 그대로 두어도 된다(coor가 새 SHA로 다시 기동한다).
- [ ] 원격 Linux make lint/test/integration(제품 변경이 없으면 make test·fixture compile만이라도 같은 SHA로), FullOps lint.py --from 688305054dc4, DEV 보고에 절 추가, packet-outcomes, finish, 커밋, worker_done.

소유 파일: fixture 테스트 파일, DEV 보고, own 기록.

## 완료 기준과 검증

fixture가 실제로 기동해 manifest와 HTTP 200이 확인됐다. 기존 검사 통과, lint ERROR 0.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: 해당 없음(테스트 fixture 수정).
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 해당 없음.
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 해당 없음.

## 갱신할 산출물

없음.

탐색 패킷의 task/role/attempt·기준 SHA와 direct_edit/impact_check/document_read/document_update 경로·구간·근거·partial/unknown을 연결한다. 명시적 없음과 미확인을 구분한다. 필수·충돌·주의 후보와 budget 잔여 확인 목록을 보존한다.

## 기대 산출물

fixture 수정, 기동 근거, DEV 보고 절, 완료 보고.

## 제약·협업·후속

미추적 `videos/service-intro/`·`qa-reports/SAR-MOTION-INTRO-001-DEV-test/` 하위 미추적 폴더는 건드리지 않는다. 다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다.

## 이전 완료 기록

- .fullops-squad/handovers/logs/2026-10-19_to_dev.md

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
