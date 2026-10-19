---
title: SAR-UX-ISSUES-002-NAMES-DEV — designer 조건부 수락의 수정 4건(새 대화 radio·링크 분리·이름 거부 문구와 입력 복구·거절 안내)을 반영한다
status: draft
updated: 2026-10-19
owner: dev
tasks: [SAR-UX-ISSUES-002-NAMES-DEV]
summary: designer 조건부 수락의 수정 4건(새 대화 radio·링크 분리·이름 거부 문구와 입력 복구·거절 안내)을 반영한다
attempt: a82026636fb14599945245451d90f702
base: 8d0db4bf59e905c9f398123c4e2e5bc13ed4c3fa
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-NAMES-DEV — designer 조건부 수락의 수정 4건(새 대화 radio·링크 분리·이름 거부 문구와 입력 복구·거절 안내)을 반영한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-19
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 아래 수정 4건과 관련 테스트·기술 문서. 새 스키마·API·권한 확대·실계정·운영 배포·삭제/이동·force push는 coor에 ask.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

같은 과제의 수정 attempt다. 이전 후보 b2e683a(제품 516bf62)는 tester QA 통과(low D-NM-01), 고정 SHA 리뷰 수락(critical/high/medium 0, low D-NM-01), designer 조건부 수락을 받았다. designer 결정은 main b923825에 통합됐고 fullops/dev에 merge(a6a3ebd)했다.

- designer 판정·결정: `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-VISUAL/report.md`, 갱신된 정본 `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` U2(2026-10-19 NAMES-VISUAL 결정 문단)·U4 표(받은 대기·거절 행).
- tester: `C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/report.md`(D-NM-01, O-NM-02). 리뷰: `C:/Users/shin/orca/workspaces/KnowsLink/coor/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-CODE-REVIEW-review/report.md`.

## 적용 기준과 예외

이전 attempt와 같다(fullops-common-0.3.3, project.md, Windows make 없음 → 격리 Linux clone에서 make lint/test/integration). 기준 8d0db4b로 lint해 단위 전체를 덮는다.

## 먼저 읽을 문서

위 designer 판정 보고, SAR-UX-ISSUES-002.md U2·U4, 이전 완료 기록(아래), 수정 대상 소스.

## 해야 할 일과 파일 소유권

- [ ] ① (medium) 새 대화 발신·수신 select를 조회와 같은 radio 라벨(`이름 / 종류 / 구분값`, 줄바꿈)로 교체. 390·320에서 동명 구별 가능해야 한다.
- [ ] ② (low, V-NM-01) 내 에이전트 상세의 '요청 상태 조회'·'연결·키 관리' 링크가 붙어 보이는 문제 분리.
- [ ] ③ (low, D-NM-01) 금지 문자 거부 시 문구 '이름에 줄바꿈·제어 문자·방향 제어 문자는 쓸 수 없습니다. 보이는 글자로 1–40자를 입력하세요.', 길이·빈 이름은 '이름은 1–40자로 입력하세요' 유지. 입력란에는 금지 문자를 제거한 나머지를 남기고 자동 저장하지 않는다. 에이전트 이름·에이전트 추가·내 이름에 같은 규칙.
- [ ] ④ (low) 받은 요청 거절 앞 안내 한 줄: '거절하면 요청한 사람의 기록에 거절로 표시됩니다. 알리지 않으려면 응답하지 않아도 기한이 지나면 요청이 끝납니다.'
- [ ] 핵심 실패 경계 테스트(새 대화 radio 선택·권한 경계 유지, 금지 문자 제거 복구·저장 안 함). 문구 복제 테스트는 쓰지 않는다.
- [ ] names fixture가 새 화면(새 대화)과 거부 흐름을 다시 검증할 수 있는지 확인·보완(tester·designer 재검수용).
- [ ] 영향받는 기술 문서 갱신, 원격 Linux make lint/test/integration, FullOps lint.py --from 8d0db4b, packet-outcomes, 완료 보고, work.py finish, 커밋, worker_done.

소유 파일: 이전 attempt와 같은 범위(`internal/relay/` 회원·친구 화면/handler/테스트·fixture, 기술 문서, own 기록). 예상 100–250줄.

## 완료 기준과 검증

4건 반영, 기존 테스트·권한 경계 회귀 없음, Linux make lint/test/integration EXIT=0, FullOps lint ERROR 0. 이후 coor가 수정분(b2e683a..새 SHA) 리뷰·tester 해당 항목 재확인·designer 해당 화면 재검수를 조정한다.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04, SAR-UX-ISSUES-002.md.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 조회 선택기의 radio fieldset 마크업·스타일을 재사용한다. 새 JS·CSS 토큰 없음.
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 테마 전환 없음. 새 대화 3폭 렌더 확인.

## 갱신할 산출물

D10, D11, D13 중 실제 변경된 것.

탐색 패킷의 task/role/attempt·기준 SHA와 direct_edit/impact_check/document_read/document_update 경로·구간·근거·partial/unknown을 연결한다. 명시적 없음과 미확인을 구분한다. 필수·충돌·주의 후보와 budget 잔여 확인 목록을 보존한다.

## 기대 산출물

코드·테스트, 기술 문서, DEV 보고 `SAR-UX-ISSUES-002-NAMES-DEV-test/report.md`의 수정 절(또는 fix 하위 보고), 완료 보고.

## 제약·협업·후속

하위 작업은 [위임 기준](../rules/delegation.md)을 따른다. 미추적 `videos/service-intro/`·`qa-reports/SAR-MOTION-INTRO-001-DEV-test/` 하위 미추적 폴더는 사용자 판단 대기 자료이므로 수정·삭제·커밋하지 않는다. O-NM-03·04는 후속 단위, 05·06·07은 현행 유지(designer 결정)이므로 이번에 고치지 않는다. 다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다.

## 이전 완료 기록

- .fullops-squad/handovers/logs/2026-10-19_to_dev.md

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
