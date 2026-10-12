---
title: SAR-HUMAN-PEER-001-VISUAL — 고정 구현 후보의 사람 친구와 대화 UI 직접 시각 수락
status: draft
updated: 2026-10-12
owner: designer
tasks: [SAR-HUMAN-PEER-001-VISUAL]
summary: 고정 구현 후보의 사람 친구와 대화 UI 직접 시각 수락
attempt: a772658c0880455a933958b3cc02f211
base: f849215bdefd3984dbf385c9f0973a74525be545
subagent_level: standard
test_level: lite
---

# SAR-HUMAN-PEER-001-VISUAL — 사람 친구·대화 UI 직접 시각 검수

- 상태: waiting-dev-sha. coor가 최종 DEV 40자리 SHA·캡처 manifest·환경을 고정한 뒤 ready로 전환한다. 현재 base는 설계 통합 기준이며 구현 SHA가 아니다.
- From / To: coor / designer. Purpose: review. 실행 codex / gpt-6.1-sol / high, 주 스킬 frontend-design. 이전 설계 작성 세션과 다른 fresh 세션을 사용한다.
- 담당: C:/Users/shin/orca/workspaces/KnowsLink/designer, fullops/designer. D04 시각 판정·자기 phase/보고서/증거/context만 소유. 제품 코드·coor PLANS/board·tester QA 변경 금지.
- 복귀: repo0b08ec4c-9e3a-4613-8197-5a835e545335, Run run_5534969da2d0, coor term_47da6145-4e1e-4cc1-93ea-20bc0dd590b9. 실제 task/dispatch는 새 preamble.

## 범위와 수락 기준

현행 D02 SAR-HUMAN-PEER HP-01–16과 D04 service-member-ui.md의 현행 SAR-HUMAN-PEER 절을 기준으로 실제 새 구현을 판정한다. 과거 polish PASS나 정적 fixture를 실제 제품 수락으로 재사용하지 않는다.

- [ ] 고정 후보·캡처 manifest·실제 loopback UI의 동일성, 명세 토큰/정보 구조를 확인한다.
- [ ] 로그인/이름·홈·친구 요청/목록·내 agent·대화/작성·해제 읽기전용·gate/오류를 390/1440 실제 캡처로 직접 본다. 320폭·긴/동명 이름·200%·포커스·44px 근거를 확인한다. 정적 렌더와 실제 auth/API/DB 캡처를 구분한다.
- [ ] 세 메뉴, 사람 아래 agent, 고정 발신/수신과 owner/agent 작성 구분이 바로 이해되는지 확인한다. 기술 ID/키/명령/긴 안내가 다시 주 화면을 차지하지 않아야 한다. 장식 카드·통계·그라디언트·알약 버튼·가짜 온라인/대화는 제외한다. 관련 보안·보존·실패 안내는 유지한다.
- [ ] qa-reports/SAR-HUMAN-PEER-001-VISUAL/report.md에 관찰 이미지 경로·범위·PASS/FAIL/not_run·심각도·수정 요청·남은 제약을 기록한다. 발견 결함을 직접 코드 수정하지 않고 coor를 통해 DEV에 전달한다.

## 적용 기준·실행 경계

fullops-common-0.3.3 README/coding-style/testing/security, project.md, FULLOPS.md, orca-agents.md, docs/agents/document-writing.md, fullops-frontend를 읽는다. 선택 test lite/subagent standard 유지. frontend.py check의 configuration_ready와 실제 인증/브라우저 검증을 구분한다. 도구·호스트 자동 전환/설치 없음. 다크 테마는 범위 밖이다.

기존 Go template/CSS/memberStyle과 6색 정본을 적용한다. 기술 의미/API/인가는 DEV와 tester/별도 코드 reviewer의 책임이다. 직접 시각 판정만 맡고 기능 QA/키보드/스크린 리더를 이미지로 PASS하지 않는다. 필요한 읽기/스크롤/캡처만 실제 합성 loopback에서 수행한다. 상태를 바꾸는 POST는 tester 소유 fixture에서 명시 협의 후만 실행한다. 사용자 실계정/운영/실메일/외부 발송 없음. 삭제·이동 정리는 금지하고 임시 자료를 보존한다.

먼저 읽을 문서: .fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md, docs/design-docs/mockups/service-member-ui.md, docs/exec-plans/phases/SAR-HUMAN-PEER-001.md, DEV 완료 보고/캡처 manifest(고정 후 연결), docs/evaluations/scenarios/README.md, project.md와 공통 규칙. 경로는 .fullops-squad/ 기준이다. 부분 packet/unknown은 일반 검색으로 보완하고 항목별 outcomes를 기록한다.

## 완료와 복귀

고정 SHA의 실제 화면 판정과 D04에 결과를 연결한다. 시각 조건이 충족되지 않으면 FAIL과 재현을 정직하게 보고한다. 미실행은 사유/영향을 기록한다. 자기 문서만 커밋하고 최종 HEAD의 FullOps lint와 문서 검사 결과를 남긴다. 제품 코드/설정 변경 없는 경우 제품 검사 재사용 또는 실행 불가 이유를 명시하며 독립 제품 QA로 확대하지 않는다. packet-outcomes·완료 전문·work.py finish로 인박스를 아카이빙한 뒤 실제 preamble worker_done을 한 번 보낸다.

## 완료 보고

미완료. 최종 후보와 실제 시각 근거가 아직 고정되지 않았다.
