---
title: SAR-UX-ISSUES-002-CONNECT-DEV — 리뷰·designer·tester 수정 묶음(재요청 행 동일화·저장 미확인 안내 등)을 반영한다
status: draft
updated: 2026-10-20
owner: dev
tasks: [SAR-UX-ISSUES-002-CONNECT-DEV]
summary: 리뷰·designer·tester 수정 묶음(재요청 행 동일화·저장 미확인 안내 등)을 반영한다
attempt: 79003121506e4d4f9d1c667f0fc595b1
base: 022c748e6fd153dea024785ebb045d28fee5a948
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-CONNECT-DEV — 리뷰·designer·tester 수정 묶음(재요청 행 동일화·저장 미확인 안내 등)을 반영한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 아래 수정 묶음과 관련 테스트·fixture·기술 문서. 상태 JSON의 기존 필드 사용 방식 변경은 기존 데이터 호환을 지키는 범위에서 승인. DB migration·새 공개 API·권한 확대·실계정·운영 배포·삭제/이동·force push는 coor에 ask.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

같은 과제의 두 번째 attempt다. 후보 8d8bad4(제품 37febf9)는 고정 SHA 리뷰(critical/high 0, medium F1, low F2–F4), designer 조건부 수락(V-CN-01 medium, V-CN-02 low), tester(critical/high 0, medium D-CT-01, low O-CT-01~04)를 받았다. designer 결정은 main a7f6b10에 통합됐다(SAR-UX-ISSUES-002.md U4 '재요청' 결정). fullops/dev에 main을 merge해 정본을 받는다.

- 리뷰: `C:/Users/shin/orca/workspaces/KnowsLink/coor/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-CODE-REVIEW-review/report.md`(F1–F4).
- designer: main의 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-VISUAL/report.md`.
- tester: `C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-TESTER-test/report.md`(D-CT-01, O-CT-01~04).

## 적용 기준과 예외

이전 attempt와 같다. 원격 Linux make lint/test/integration과 fixture 실제 기동 확인을 유지한다. lint 기준 022c748.

## 먼저 읽을 문서

위 세 보고의 해당 절, SAR-UX-ISSUES-002.md U4 재요청 결정·U8 상태 표, `internal/relay/human.go`(requestFriend·sweepFriends), `internal/relay/member_pages.go`(연결 화면 68행 부근·legacy toc), 연결 거부 복귀 handler, gate 상세 템플릿, 이전 완료 기록.

## 해야 할 일과 파일 소유권

- [ ] (medium, F1=V-CN-01=O-CT-01) 재요청 행: designer 결정대로 기존 요청자 행을 재사용해 있는/없는/차단/상대 수신 한도/자기 코드에서 행 수·순서·상태·끝 4자·만료·세대가 같게 한다. 차단 기한 뒤·차단 취소 뒤·수신 한도 10건 경로의 대리 행 2행 제거. mine 대기 상한이 blocked-pending도 같은 방식으로 센다. 기존 저장 데이터(이미 2행인 경우) 처리 방식을 정하고 기록한다.
- [ ] (medium, D-CT-01) 승인 뒤 클라이언트 저장이 확인되지 않은 시도의 연결 화면이 '만료 · 새 키가 연결되지 않았습니다 · 새로 시작하세요'로 단정하지 않게 한다. 같은 agent 카드의 '확인 필요'와 U8 상태 표('승인 후 응답 유실': 결과 대조 먼저, 맹목적 새 등록 금지)에 맞춘다.
- [ ] (low, V-CN-02) `/home/connections/{없는 ID}` 거부 뒤 복귀를 '내 에이전트'로.
- [ ] (low, F2) 요약을 만들 수 없는 gate 상세에 승인 버튼이 남는 문제(3단위 '확인할 수 없습니다' 규칙과 일치).
- [ ] (low, F3) 새 키 상태 문구의 '키' 중복.
- [ ] (low, O-CT-02) legacy 목차에 보이는 '이 화면에서 찾기' 제목(aria-current 없음 유지).
- [ ] (low, O-CT-04) 5분 재확인 거부 안내가 시작 위치(내 에이전트)에서 바로 재확인할 수 있게 문구·링크 정리.
- [ ] F1·D-CT-01 실패 경계 자동 검사(행 수·순서·상태 동일성, 차단 비추론, 저장 미확인 안내). fixture에 차단·수신 한도 경로가 없으면 추가해 tester가 재현할 수 있게 한다(실제 기동 확인).
- [ ] 원격 Linux make lint/test/integration, FullOps lint.py --from 022c748e6fd153dea024785ebb045d28fee5a948, DEV 보고 수정 절, packet-outcomes, finish, 커밋, worker_done.

소유 파일: `internal/relay/` 친구·연결·gate 화면/handler/테스트·fixture, 기술 문서, DEV 보고, own 기록. O-CT-03(카드 제목 3줄 생략)은 기존 패턴이라 이번에 바꾸지 않는다.

## 완료 기준과 검증

medium 2건이 자동 검사와 fixture로 확인되고 low 5건이 반영됐다. 기존 회귀 없음, Linux make lint/test/integration EXIT=0, lint ERROR 0, fixture 실제 기동. 이후 coor가 수정분 diff 확인·tester 재확인을 조정한다.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04, SAR-UX-ISSUES-002.md.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 기존 템플릿·문구 패턴. 새 토큰·JS 없음.
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 테마 전환 없음. 바뀐 화면 3폭 렌더 넘침 0.

## 갱신할 산출물

D10, D11, D13 중 실제 변경된 것.

탐색 패킷의 task/role/attempt·기준 SHA와 direct_edit/impact_check/document_read/document_update 경로·구간·근거·partial/unknown을 연결한다. 명시적 없음과 미확인을 구분한다. 필수·충돌·주의 후보와 budget 잔여 확인 목록을 보존한다.

## 기대 산출물

코드·테스트·fixture, 기술 문서, DEV 보고 수정 절, 완료 보고.

## 제약·협업·후속

미추적 `videos/service-intro/`·`qa-reports/SAR-MOTION-INTRO-001-DEV-test/` 하위 미추적 폴더는 건드리지 않는다. workspace 밖에 파일을 만들지 않는다. 다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다.

## 이전 완료 기록

- .fullops-squad/handovers/logs/2026-10-20_to_dev.md

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
