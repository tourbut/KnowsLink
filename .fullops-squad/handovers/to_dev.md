---
title: SAR-UX-ISSUES-002-RECEIPT-DEV-FIX — 조회 선택기 닫힌 상태에서 긴 이름의 종류·구분값이 가려지는 D-RT-01 수정
status: draft
updated: 2026-10-17
owner: dev
tasks: [SAR-UX-ISSUES-002-RECEIPT-DEV-FIX]
summary: 조회 선택기 닫힌 상태에서 긴 이름의 종류·구분값이 가려지는 D-RT-01 수정
attempt: 75bb365492914387b031ecf95356646c
base: 9ae29ec4ddfcfacb59ed48871a14bb558342c26e
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-RECEIPT-DEV-FIX — 조회 선택기 닫힌 상태에서 긴 이름의 종류·구분값이 가려지는 D-RT-01 수정

현재 과제는 이 역할 인박스에서 관리한다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

Task key: SAR-UX-ISSUES-002-RECEIPT-DEV-FIX; Purpose: implementation

- 작성일: 2026-10-17
- From / To: coor / dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위: 아래 파일 소유권 안의 제품 코드·테스트·D10 문서 수정. 새 dependency, 운영 변경, 삭제·이동, force push는 하지 않는다.
- 담당 repo id / 워크트리 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀: repo 0b08ec4c-9e3a-4613-8197-5a835e545335 / coor C:/Users/shin/orca/workspaces/KnowsLink/coor / term_75c9a893-7418-4abd-9ad0-eba86cbb6b68 / run_1d685ec8e9bf (task·dispatch id는 preamble)

## 현재 상황과 확인 근거

DEV 9ae29ec(제품 source 1e49033)는 make lint/test를 통과했다. 독립 QA(tester 8ff3506, Claude)가 결함 D-RT-01(medium, U2/U6)을 보고했다. 닫힌 조회 선택기가 39자 이름·종류(`Node 22 로컬 클라이언트`)·구분값(`agent_a`)을 한 문자열로 보여 주어 폭 부족(1440에서 가용 750px, 390에서 293px, 320에서 223px)으로 구분값이 잘린다. 같은 이름의 `agent_a`/`agent_a2`를 닫힌 상태에서 구별할 수 없다. 데이터 노출·권한 문제는 아니다. 원본: tester 워크트리 .fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-TESTER-test/report.md의 D-RT-01 절과 browser-chrome.json·browser-edge.json의 selector-closed-shows-full-name-kind-distinguisher-* 세 건. 같은 내용은 tester 브랜치(fullops/tester 8ff3506)에서 읽는다.

## 적용 기준과 예외

공통 규칙 .fullops-squad/rules/common/README.md와 연결된 코딩·테스트·보안 규칙, 수락 제품 기준 .fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md의 U2·U6, D04 디자인 토큰을 따른다. 기존 보안·권한 수락 기준을 완화하지 않는다. 예외 없음.

## 먼저 읽을 문서

- .fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md (U2, U6)
- 이전 DEV 결과 .fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md
- 선택기 렌더링 코드(internal/relay의 receipt 템플릿과 owned agent 라벨 생성 부분)와 관련 테스트. 위치는 직접 검색한다.
- tester 보고서의 D-RT-01 절(위 경로)

## 해야 할 일과 파일 소유권

- [ ] 닫힌 선택기에서도 선택한 agent의 이름·종류·안정 구분값이 구별되게 고친다. 구분값을 앞에 두거나 이름을 줄이거나 두 줄로 나누는 방식 중 기존 D04 토큰·44px·글꼴을 그대로 쓰는 가장 작은 방식을 고른다. 320/390/1440에서 구분값이 잘리지 않아야 한다.
- [ ] 기존 동명 구분값 충돌(전체 ID 확장) 동작과 escape·단일 autofocus·키보드 라벨을 유지한다.
- [ ] 변경한 동작의 테스트를 추가·갱신한다. 새 JS·CSS 파일·dependency는 추가하지 않는다.
- [ ] D10 정본에 표시 방식 변경을 반영한다.

소유: 선택기 렌더링 코드·해당 테스트·D10·own same-key route/find/context/packet/outcomes. 예상 변경은 작은 규모(코드 수십 줄 이내). 제품 규칙 변경이 필요해 보이면 ask로 묻는다.

## 완료 기준과 검증

- 관련 Go 검사(TestReceipt|TestMemberReceiptView|TestMemberUIFixtures)가 통과한다.
- 변경 SHA에서 정본 lint.py(product-lint, product-test)를 실행하고 원 결과를 보존한다. Windows에서 make가 없으면 기존 DEV처럼 remote.py의 exact-HEAD Linux clone으로 실행한다.
- 닫힌 선택기에서 구분값이 보인다는 것을 템플릿/골든 수준으로 확인한다. 실제 브라우저 폭 판정은 tester 재검증과 designer 시각 검수가 맡는다.
- DEV 완료와 제품 최종 수락은 구분한다. 재검증(tester 선택기 세 건)과 고정 SHA 독립 코드 리뷰, designer 시각 수락, main 통합은 coor가 조정한다.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 토큰 정본: D04 디자인 토큰, 기존 여섯 토큰·44px·글꼴·간격
- 재사용 컴포넌트: 기존 조회 선택기 변형만 사용. 의도적 예외 없음
- 디자인 lint: frontend.json tools가 비어 있어 해당 없음
- 테마 전환 확인: 변경 없는 토큰 사용이므로 전파 확인 해당 없음. 실제 렌더는 tester·designer 담당

## 갱신할 산출물

D10. 탐색 패킷은 jev_find/jev_context/jev_packet으로 이 과제 키로 생성한 파일을 따른다.

## 기대 산출물

수정 코드, 테스트, D10 갱신, packet-outcomes.json, work.py finish 로그.

## 제약·협업·후속

하위 에이전트는 선택이다. 소규모 직렬 작업이라 직접 처리해도 된다. 범위 밖: 이름/사람 친구, 승인 요약, 설치 재개, 언어 설정. 평범한 말줄임표로 구분값을 숨기는 방식과 이름만 줄이고 구분값을 여전히 가리는 방식은 피한다. 이미지 생성이 필요하면 grok을 쓴다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-UX-ISSUES-002-RECEIPT-DEV-FIX-packet.json` / SHA `9ae29ec4ddfcfacb59ed48871a14bb558342c26e` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_read) · 줄 6, 10, 14 · inferred · 필수 · {'relevant': 0.21, 'evidence': 0.41, 'contradicts': 0.46, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 18, 18, 31, 31, 31, 31, 35, 39, 40, 40 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.16, 'evidence': 0.49, 'contradicts': 0.26, 'injection': 0.01, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.27, 'evidence': 0.71, 'contradicts': 0.12, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `internal/relay/member_pages.go` (impact_check) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.04, 'evidence': 0.06, 'contradicts': 0.51, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 34, 36 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 6, 58, 58 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 7, 14 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_read, document_update) · 줄 7, 7, 14, 20, 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-FIX-TESTER.md` (document_read) · 줄 65, 70 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW-review/report.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 62 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN/report.md` (document_read) · 줄 6, 10, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md` (document_read) · 줄 6, 6, 52, 54 · inferred · {'relevant': 0.22, 'evidence': 0.49, 'contradicts': 0.74, 'injection': 0.11, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG.md` (document_read) · 줄 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-DESIGN.md` (document_read) · 줄 6, 10, 14, 38, 42 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-RECEIPT-DEV.md` (document_read) · 줄 6, 6, 20, 20 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 7, 7, 13 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/protocol.md` (document_read) · 줄 181 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_dev.md` (document_read) · 줄 136, 136, 137, 137, 140, 140, 144, 144, 152, 152, 156, 185, 199, 199 · inferred
- `adapters/src/inbox-files.test.ts` (impact_check) · 줄 30 · inferred
- `docs/README.md` (document_read) · 줄 전체/미확인 · inferred
- `docs/development.md` (document_read) · 줄 전체/미확인 · inferred
- `internal/relay/admission_integration_test.go` (impact_check) · 줄 39, 69, 73, 149, 157, 170 · inferred
- `internal/relay/cleanup_flood_integration_test.go` (impact_check) · 줄 35, 74, 92, 118, 199, 211 · inferred
- `internal/relay/cleanup_unit_integration_test.go` (impact_check) · 줄 52, 68 · inferred
- `internal/relay/connections_test.go` (impact_check) · 줄 20, 25, 41, 44, 64, 77, 86, 104, 122, 134, 162, 170, 183, 188, 191, 192, 195, 198, 203, 219, 224, 229, 232, 240, 250, 253, 254, 277, 287, 346, 354, 355, 370, 371, 372, 373, 374, 375, 378, 384, 394 · inferred
- `internal/relay/human_http.go` (impact_check) · 줄 전체/미확인 · inferred
- `internal/relay/human_pages.go` (impact_check) · 줄 전체/미확인 · inferred
- `internal/relay/human_pages_test.go` (impact_check) · 줄 18, 19, 26, 115, 199, 215, 218 · inferred
- `internal/relay/member_agents.go` (impact_check) · 줄 전체/미확인 · inferred
- `internal/relay/member_receipt.go` (impact_check) · 줄 전체/미확인 · inferred · {'relevant': 0.07, 'evidence': 0.07, 'contradicts': 0.42, 'injection': 0.01, 'decision': 'keep', 'reason': None}
- `internal/relay/member_ui_test.go` (impact_check) · 줄 207 · inferred · {'relevant': 0.09, 'evidence': 0.13, 'contradicts': 0.32, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `internal/relay/policy_test.go` (impact_check) · 줄 16, 19, 22, 31, 48, 57, 59, 89, 96, 109, 111, 114, 120, 146, 157, 164, 171, 173, 192, 195, 202, 205, 212, 215, 220, 228, 231, 247, 252, 255 · inferred
- `internal/relay/protocol_test.go` (impact_check) · 줄 49, 54, 68, 77 · inferred
- `internal/relay/public_messages_integration_test.go` (impact_check) · 줄 210, 218, 270, 272, 282, 284, 286, 302, 325, 326, 407, 417, 421 · inferred
- `internal/relay/receipt_context_test.go` (impact_check) · 줄 19, 20, 20, 41, 46, 48, 49, 50, 52, 53, 105, 118, 120, 122 · inferred · {'relevant': 0.11, 'evidence': 0.24, 'contradicts': 0.15, 'injection': 0.01, 'decision': 'keep', 'reason': None}
- `internal/relay/receipt_integration_test.go` (impact_check) · 줄 34, 37, 39, 41, 49, 51, 62, 71, 93 · inferred
- `internal/relay/test_messages_integration_test.go` (impact_check) · 줄 25, 26, 30, 32, 34, 37, 38, 40, 41, 52, 73, 74, 75 · inferred
- `internal/relay/testdata/legacy_claims.json` (impact_check) · 줄 14, 38, 39, 51, 72, 94, 119, 140, 162, 211, 232, 254, 278, 279, 280, 289 · inferred
- `internal/relay/text_events_integration_test.go` (impact_check) · 줄 88, 90 · inferred
- `internal/relay/text_events_test.go` (impact_check) · 줄 13, 14, 28 · inferred
- `scripts/check_public_ingress.py` (impact_check) · 줄 11 · inferred
미확인 19건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
의존성 변경이 없으면 그 근거를 적는다.
처리는 fullops-work, 회신은 fullops-orca 스킬을 따른다.
