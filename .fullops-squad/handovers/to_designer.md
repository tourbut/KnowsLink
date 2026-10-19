---
title: SAR-UX-ISSUES-002-NAMES-VISUAL — DEV 수정 5fa3333의 새 대화·내 에이전트·친구 요청·이름 오류 화면을 재검수한다
status: draft
updated: 2026-10-19
owner: designer
tasks: [SAR-UX-ISSUES-002-NAMES-VISUAL]
summary: DEV 수정 5fa3333의 새 대화·내 에이전트·친구 요청·이름 오류 화면을 재검수한다
attempt: 084aab73136444028cb931e5c8747c3c
base: b923825a5b4ed05675d4c3a169f65f63882bed22
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-NAMES-VISUAL — DEV 수정 5fa3333의 새 대화·내 에이전트·친구 요청·이름 오류 화면을 재검수한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-19
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: fixture 허용 GET과 금지 문자 이름 POST(저장 안 됨)로 관찰, VISUAL 보고 갱신과 designer 브랜치 커밋. 제품 코드 수정, 운영·실계정·외부 발송, 삭제·이동, force push는 하지 않는다.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

같은 과제의 재확인 attempt다. 이전 결과는 아래 이전 완료 기록에 있다.

- DEV 수정: fullops/dev `5fa3333b5948fd36792662b202915bc6243cac14`(제품 `97a8899`, b2e683a 대비 internal 8파일 +65/-17). 원격 Linux make lint/test/integration EXIT=0, FullOps lint ERROR 0. 소스·DEV 보고(수정 절)는 `C:/Users/shin/orca/workspaces/KnowsLink/dev`에서 읽기 전용.
- 수정 4건: ① 새 대화 발신·수신 radio(`이름 / 종류 / 구분값`) ② 내 에이전트 상세 '요청 상태 조회'·'연결·키 관리' 링크 줄 분리 ③ 이름 금지 문자 전용 문구와 금지 문자 제거 입력 복구(저장 안 함, 내 이름·에이전트 이름·에이전트 추가·device·remote 동의 5경로) ④ 받은 요청 거절 앞 안내(aria-describedby).
- designer 판정·결정: main의 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-VISUAL/report.md`, 정본 `SAR-UX-ISSUES-002.md` U2·U4.
- coor 확인: 수정분 need=skip(8파일 82줄), 병합 책임자가 diff를 읽었다. 권한 경계 변경 없음.

## fixture

origin `http://127.0.0.1:18930`(coor SSH 포워딩, 2026-10-19 재기동, 24시간). 서버 clone `/tmp/knowslink-names002-5fa3333b5948-cd4acd86`, container `kl-names002-fixture-5fa3333b5948-385e64d3`, 고정 SHA 5fa3333. manifest `C:/Users/shin/AppData/Local/KnowsLinkDevTools/names-ux002-dev/fixture-manifest.json`. 합성 A/B/C·허용 GET/POST·계정 전환은 이전 attempt와 같다(DEV 보고 'tester/designer용 격리 fixture' 절). tester와 designer가 같은 fixture를 함께 쓴다.


## 적용 기준과 예외

이전 attempt와 같다(fullops-common-0.3.3, project.md, SAR-UX-ISSUES-002.md U2·U4·U7, D04). 등록 make lint/test는 DEV 같은 SHA 결과를 재사용한다.

## 먼저 읽을 문서

이전 완료 기록, DEV 보고 수정 절, designer 판정 보고, SAR-UX-ISSUES-002.md U2 결정 문단·U4 표.

## 해야 할 일과 파일 소유권

- [ ] 새 대화 발신·수신 radio 화면(3폭), 내 에이전트 상세 링크 분리, 받은 요청 거절 앞 안내, 이름 금지 문자 오류 화면 1장을 직접 판정(PASS/FAIL/not_run, 근거 캡처). 이전 조건부 수락의 조건 충족 여부를 결론으로 적는다.
- [ ] 캡처는 tester 재확인 캡처(`C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/`) 또는 fixture 허용 GET으로 직접 만든다. 오류 화면은 금지 문자가 든 이름 POST(저장되지 않음)만 허용한다. 친구 결정·코드 교체·정상 이름 저장 POST는 하지 않는다.
- [ ] VISUAL 보고에 재검수 절 추가 → 완료 보고 → work.py finish → 커밋 → worker_done(결론과 남은 DEV 수정 여부).

소유 파일: `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-VISUAL/` 하위, own route/inbox/log.

## 완료 기준과 검증

재검수 결론(수락/조건부/거부)과 근거 캡처, 남은 수정이 있으면 DEV 기준. 스크린리더·실제 200%는 not_run 유지 가능.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 해당 없음(제품 코드 수정 없음).
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 테마 전환 없음.

## 갱신할 산출물

없음(정본 결정 변경이 필요하면 SAR-UX-ISSUES-002.md만).

탐색 패킷의 task/role/attempt·기준 SHA와 direct_edit/impact_check/document_read/document_update 경로·구간·근거·partial/unknown을 연결한다. 명시적 없음과 미확인을 구분한다. 필수·충돌·주의 후보와 budget 잔여 확인 목록을 보존한다.

## 기대 산출물

기존 보고서의 재확인 절(또는 하위 보고), 결과 JSON·캡처, 완료 보고.

## 제약·협업·후속

제품 코드·다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다. 운영·실계정·외부 발송·삭제/이동·force push는 하지 않는다.

## 이전 완료 기록

- .fullops-squad/handovers/logs/2026-10-19_to_designer.md

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
