---
title: SAR-UX-ISSUES-002-CONNECT-TESTER — DEV 수정 abc2191(제품 d6369a4)의 재요청 동일화·D-CT-01·low를 재확인한다
status: draft
updated: 2026-10-20
owner: tester
tasks: [SAR-UX-ISSUES-002-CONNECT-TESTER]
summary: DEV 수정 abc2191(제품 d6369a4)의 재요청 동일화·D-CT-01·low를 재확인한다
attempt: b496587ea088470abe5ac3b1cf97864d
base: d3e676a0a062cd2e77511e2689436e2cfe5613b6
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-CONNECT-TESTER — DEV 수정 abc2191(제품 d6369a4)의 재요청 동일화·D-CT-01·low를 재확인한다

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

같은 과제의 재확인 attempt다. 직전 결과 822415d(medium D-CT-01, low O-CT-01~04). DEV 수정 `abc219120509027b143b0e4ad121511a5a18395b`(fullops/dev, 제품 `d6369a4`): 재요청 행 동일화(F1·V-CN-01·O-CT-01), D-CT-01 저장 미확인 안내, V-CN-02 복귀, F2 요약 없는 gate 버튼 제거, F3 문구, O-CT-02 legacy 목차 보이는 제목, O-CT-04 재확인 안내. Linux make lint/test/integration EXIT=0. 수정분 독립 delta 리뷰가 병행 중이다. DEV 보고 수정 절: `C:/Users/shin/orca/workspaces/KnowsLink/dev/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-DEV-test/report.md`('fixture 추가 경로 (2차)'·'남은 한계 (2차)').

## fixture

origin `http://127.0.0.1:18932`(coor SSH 포워딩, abc2191로 2026-10-20 재기동, 수명 12시간, 초기 상태). 서버 clone `/tmp/knowslink-connect002-abc219120509-0b21a980`, container `kl-connect002-fixture-abc219120509-354beae3`. manifest(origin·accounts·agents·connections·rerequest_codes·expires_at): `C:/Users/shin/AppData/Local/KnowsLinkDevTools/connect-ux002-dev/fixture-manifest.json`. 계정 a·b·c와 새 d·e(`POST /fixture/as/{a..e}`). `rerequest_codes`: `blocked`(D가 A의 대기 요청을 차단, A 화면 만료), `receive_limit`(E가 대기 요청 10건 수신, A의 이전 요청 만료), `unknown_cancelled`(없는 코드 취소 행, 비교 기준). 쿠키는 `__Host-` Secure라 http cookie jar가 버릴 수 있다(DEV 보고 참고).

## 적용 기준과 예외

이전 attempt와 같다. 제품 정본 SAR-UX-ISSUES-002.md U4 재요청 결정·U8 상태 표.

## 먼저 읽을 문서

DEV 보고 수정 절, 이전 완료 기록, U4 재요청 결정, 기존 tester 보고 D-CT-01·O-CT 절.

## 해야 할 일과 파일 소유권

- [ ] 재요청 3경로(blocked·receive_limit·unknown_cancelled): A로 같은 코드를 다시 보낸 뒤 A 화면의 행 수·순서·상태·끝 4자·만료 형식이 세 경로에서 같고, D·E 화면이 변하지 않음(차단·수신 한도 비추론). 반복 요청·취소 뒤 재요청도 확인.
- [ ] mine 대기 상한이 보류 행을 같은 방식으로 센다는 것을 관측 가능한 범위로 확인(불가하면 소스·DEV 테스트 근거로 기록).
- [ ] D-CT-01 재판정: 승인 뒤 저장 미확인 시도의 연결 화면이 카드의 '확인 필요'와 일치하고 '만료·새로 시작' 단정이 없음. 기존 FAIL 검사의 기대를 낮추지 않는다.
- [ ] low 회귀: 없는 연결 ID 거부 복귀 '내 에이전트', legacy 목차 보이는 '이 화면에서 찾기'(aria-current 없음), 재확인 안내, 키 문구, 요약 없는 gate 버튼 없음(소스·테스트 근거 가능).
- [ ] 이전 PASS 핵심(5선택지·재개 행·중복 생성 방지·O-NM-03·3폭 넘침 0)을 빠르게 회귀.
- [ ] 보고에 재확인 절 추가 → 완료 보고 → work.py finish → 커밋 → worker_done.

소유 파일: `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-TESTER-test/` 하위, 시나리오, own route/inbox/log.

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
