---
title: "SAR-UX-ISSUES-002-NAMES-DEV — tester D-NM-02(medium): 새 대화 수신 그룹 legend의 긴 친구 이름 가로 넘침을 고친다"
status: draft
updated: 2026-10-19
owner: dev
tasks: [SAR-UX-ISSUES-002-NAMES-DEV]
summary: "tester D-NM-02(medium): 새 대화 수신 그룹 legend의 긴 친구 이름 가로 넘침을 고친다"
attempt: fdbc7013b4884f0f85468e7dedd4b93d
base: 8d0db4bf59e905c9f398123c4e2e5bc13ed4c3fa
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-NAMES-DEV — tester D-NM-02(medium): 새 대화 수신 그룹 legend의 긴 친구 이름 가로 넘침을 고친다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-19
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: legend 줄바꿈 수정과 필요한 회귀 검사·DEV 보고 갱신. 그 밖의 변경은 coor에 ask.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

같은 과제의 세 번째 attempt다. 후보 5fa3333(제품 97a8899)은 designer 재검수 수락, tester 재확인에서 수정 4건 PASS. 새 결함 D-NM-02(medium): 40자 이름 친구와 새 대화를 열면(B→A `/home/new?peer=...`) 수신 그룹 `<legend>`(`{{.Peer.Name}}의 수신 에이전트`)가 줄바꿈되지 않아 390px에서 282px, 320px에서 352px 가로 넘침. 근거: `C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/report.md`(D-NM-02), 캡처 `C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/chrome-fix-new-peer-b-390.png`. 공용 CSS는 `internal/relay/member_pages.go`의 `legend{font-weight:600;padding:0;margin:16px 0 4px}`.

## 적용 기준과 예외

이전 attempt와 같다. D04 토큰 유지, 새 JS·토큰 없음.

## 먼저 읽을 문서

tester 보고 D-NM-02 절, `internal/relay/member_pages.go` 공용 CSS, `internal/relay/human_pages.go` human-new 템플릿.

## 해야 할 일과 파일 소유권

- [ ] legend가 좁은 폭에서 줄바꿈되게 공용 CSS를 고친다(조회·새 대화 등 모든 radio fieldset legend에 적용).
- [ ] 3폭(1440·390·320) 렌더로 40자 한글·이모지 친구 이름의 넘침 0 확인. 기존 렌더 probe 재사용.
- [ ] Linux make lint/test(가능하면 기존 remote helper), FullOps lint.py --from 8d0db4b, DEV 보고에 짧은 절 추가, packet-outcomes, finish, 커밋, worker_done.

소유 파일: `internal/relay/member_pages.go`(CSS), DEV 보고, own 기록. 예상 1–5줄.

## 완료 기준과 검증

D-NM-02 재현 경로에서 넘침 0, 기존 테스트 통과, lint ERROR 0. 이후 tester가 해당 경로만 재확인한다.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 공용 legend 규칙에 줄바꿈 속성만 추가.
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 테마 전환 없음.

## 갱신할 산출물

없음(필요 시 D11 짧은 갱신).

탐색 패킷의 task/role/attempt·기준 SHA와 direct_edit/impact_check/document_read/document_update 경로·구간·근거·partial/unknown을 연결한다. 명시적 없음과 미확인을 구분한다. 필수·충돌·주의 후보와 budget 잔여 확인 목록을 보존한다.

## 기대 산출물

CSS 수정, 렌더 확인 결과, DEV 보고 절, 완료 보고.

## 제약·협업·후속

미추적 `videos/service-intro/`·`qa-reports/SAR-MOTION-INTRO-001-DEV-test/` 하위 미추적 폴더는 건드리지 않는다. 다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다.

## 이전 완료 기록

- .fullops-squad/handovers/logs/2026-10-19_to_dev.md

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
