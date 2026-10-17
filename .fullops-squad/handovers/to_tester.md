---
title: SAR-UX-ISSUES-002-RECEIPT-RETEST — D-RT-01 수정 후보 71dcac7의 닫힌 조회 선택기 재검증
status: draft
updated: 2026-10-17
owner: tester
tasks: [SAR-UX-ISSUES-002-RECEIPT-RETEST]
summary: D-RT-01 수정 후보 71dcac7의 닫힌 조회 선택기 재검증
attempt: 165b25c204074e81a13d6cdf0725c130
base: 7f341230a26abc4afa97850687e624f2a583715c
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-RECEIPT-RETEST — D-RT-01 수정 후보 71dcac7의 닫힌 조회 선택기 재검증

현재 과제는 이 역할 인박스에서 관리한다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

Task key: SAR-UX-ISSUES-002-RECEIPT-RETEST; Purpose: implementation

- 작성일: 2026-10-17
- From / To: coor / tester
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위: 기존 합성 fixture 사용 검증, 보고서·캡처 작성. 제품 코드 수정, 운영 서비스·실계정·외부 발송, 삭제·이동, force push는 하지 않는다.
- 담당 repo id / 워크트리 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/tester / fullops/tester
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀: repo 0b08ec4c-9e3a-4613-8197-5a835e545335 / coor C:/Users/shin/orca/workspaces/KnowsLink/coor / term_75c9a893-7418-4abd-9ad0-eba86cbb6b68 / run_1d685ec8e9bf (task·dispatch id는 preamble)

## 현재 상황과 확인 근거

이전 QA(SAR-UX-ISSUES-002-RECEIPT-TESTER, tester 8ff3506)가 D-RT-01(medium)을 보고했고 DEV가 71dcac7(제품 source da53e96)에서 고쳤다. 선택기 옵션이 구분값을 앞에 둔다. 독립 코드 리뷰는 수락(critical/high 0, low 3)이다. low: (1) 구분값을 앞에 두어 U2의 이름 선두 순서와 다르고 긴 이름은 좁은 폭에서 종류가 잘릴 수 있다. (2) 익명 한도 429가 '조회할 내 에이전트가 없습니다'로 잘못 안내한다. 리뷰 보고서: 이 브랜치의 .fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-CODE-REVIEW-review/report.md. 이전 QA 보고서와 스크립트는 같은 브랜치의 qa-reports/SAR-UX-ISSUES-002-RECEIPT-TESTER-test/에 있다.

## 적용 기준과 예외

.fullops-squad/rules/common/README.md와 연결된 규칙, 수락 제품 기준 .fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md의 U2·U6, D04 토큰. 예외 없음. 기존 안전 경계를 완화하지 않는다.

## 먼저 읽을 문서

- 이전 QA 보고서 qa-reports/SAR-UX-ISSUES-002-RECEIPT-TESTER-test/report.md 와 browser-chrome.json·browser-edge.json, 재사용할 검사 스크립트
- 코드 리뷰 보고서(위 경로)
- DEV 수정 보고서 C:/Users/shin/orca/workspaces/KnowsLink/dev/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-FIX-test/report.md
- DEV fixture 인계 절: dev의 qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md

## 해야 할 일과 파일 소유권

- [ ] fixture를 확인하고 이전 QA의 selector-closed-shows-full-name-kind-distinguisher 세 건(320/390/1440)을 같은 방식으로 다시 실행한다. 기존 검사를 약화하지 않는다. 검사가 새 표시 순서(구분값 먼저)를 기대하지 않아 실패하면 기대를 바꾸지 말고 관찰 결과와 판단 근거를 보고한다.
- [ ] 구분값이 닫힌 상태에서 보이는지, agent_a와 agent_a2가 구별되는지, 종류가 잘리는 정도를 폭별로 기록한다.
- [ ] U6 핵심 회귀(입력 보존, 없는/타인 동일 문구, 계정 격리, 429)를 짧게 다시 확인한다. 변경 없는 전체 검사를 복제하지 않는다.
- [ ] 리뷰 low 2(익명 한도 429 문구)를 fixture로 재현 가능하면 관찰만 기록한다.
- [ ] 캡처는 designer 시각 판정에 필요한 닫힌 선택기 상태만 만든다(320/390/1440). 이미지 생성이 필요하면 grok을 쓴다.

소유: 재검증 시나리오·보고서·캡처, own same-key route/find/context/packet/outcomes, context/inbox/log.

## 완료 기준과 검증

선택기 세 건의 PASS/FAIL, U6 회귀 결과, 미검사 목록, 신규 결함의 재현·기대·실제·고정 소스를 보고한다. 새 결함은 DEV 담당으로 coor에 보낸다. 고정 소스 SHA는 71dcac7이고 fixture는 같은 SHA다. DEV 전체 gate는 재실행하지 않는다.

### UI 작업의 디자인 기준과 검증

디자인 정본은 D04. 이 작업은 검증이며 UI 변경을 하지 않는다. 실제 렌더 판정은 기존 Chrome·Edge CDP 방법을 재사용하고 designer 시각 수락은 별도다.

## fixture

origin http://127.0.0.1:18928 (SSH 포워딩, 2026-10-17 시작 24시간). 서버 clone /tmp/knowslink-receipt002-71dcac7822b8-b2c15735, 합성 DB container kl-receipt002-fixture-71dcac7822b8-cb5a6311. 합성 POST·A/B/C·ID 표는 DEV fixture 인계 절과 같다. 이 fixture는 새로 시작해 revoke 상태가 아니다. revoke는 맨 마지막에만 실행한다.

## 갱신할 산출물

없음(검증 보고서와 시나리오). 탐색 패킷은 이 과제 키 파일을 따른다.

## 기대 산출물

재검증 보고서, 검사 결과 JSON, 캡처, packet-outcomes.json, work.py finish 로그.

## 제약·협업·후속

제품 코드를 수정하지 않는다. 범위 밖: 이름/사람 친구, 승인 요약, 설치 재개, 언어 설정, 변경 없는 첫페이지 검사. 보존 종료 receipt, 명시 로그아웃, 실제 UI 200%, 스크린리더는 이전 미실행 항목이며 이번에도 fixture로 가능한 것만 시도하고 못 한 것은 미검사로 남긴다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-UX-ISSUES-002-RECEIPT-RETEST-packet.json` / SHA `2c68b4cf381da0eaa55f899e4c505ecbdaf06a85` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_read) · 줄 6, 10, 14 · inferred · 필수 · {'relevant': 0.15, 'evidence': 0.22, 'contradicts': 0.36, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-TESTER-test/report.md` (document_read) · 줄 6, 6, 10, 10 · inferred · 필수 · {'relevant': 0.82, 'evidence': 0.94, 'contradicts': 0.38, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_tester.md` (document_read) · 줄 2, 2, 2, 6, 6, 6, 14, 14, 14, 18, 18, 18, 31, 31, 35, 39, 39, 41, 41, 42, 42 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.13, 'evidence': 0.35, 'contradicts': 0.12, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.16, 'evidence': 0.38, 'contradicts': 0.1, 'injection': 0.35, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.67, 'evidence': 0.83, 'contradicts': 0.12, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 34, 36 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 6, 6, 38, 38 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-TESTER-test/report.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-SERVICE-001-REVIEW-review/report.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN/report.md` (document_read) · 줄 6, 10, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 2, 9, 9, 15, 15, 18, 58, 59, 60, 60 · inferred · {'relevant': 0.89, 'evidence': 0.94, 'contradicts': 0.79, 'injection': 0.03, 'decision': 'conflict', 'reason': 'contradicts task'}
- `.fullops-squad/docs/exec-plans/phases/SAR-REPO-DOCS-001-DEV.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-TESTER.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-DESIGN.md` (document_read) · 줄 6, 10, 14, 38, 42 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/protocol.md` (document_read) · 줄 181 · inferred
- `.fullops-squad/handovers/_TEMPLATE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_tester.md` (document_read) · 줄 9, 9, 10, 10, 13, 13, 17, 17, 25, 25, 40, 40, 44, 63, 63, 72, 72 · inferred
- `.fullops-squad/review/_REPORT.md` (document_read) · 줄 전체/미확인 · inferred
미확인 19건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
처리는 fullops-work, 회신은 fullops-orca 스킬을 따른다.
