---
title: SAR-UX-ISSUES-002-NAMES-DEV — 에이전트 이름 중심 표시·친구 코드 공유·사람 친구 요청 상태 안내(U2/U4+U7)를 구현한다
status: draft
updated: 2026-10-18
owner: dev
tasks: [SAR-UX-ISSUES-002-NAMES-DEV]
summary: 에이전트 이름 중심 표시·친구 코드 공유·사람 친구 요청 상태 안내(U2/U4+U7)를 구현한다
attempt: 9a5aabbe90544e148909a807a7a14d18
base: 8d0db4bf59e905c9f398123c4e2e5bc13ed4c3fa
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-NAMES-DEV — 에이전트 이름 중심 표시·친구 코드 공유·사람 친구 요청 상태 안내(U2/U4+U7)를 구현한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-18
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 아래 U2/U4+해당 U7 구현·검사·기술 문서 갱신은 승인됨. 새 DB 스키마·새 공개 API·권한 확대·실계정/실메일/상대 발송·운영 배포·삭제/이동·force push는 coor에 ask.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

사용자 2026-10-18 지시 "잔여작업 모두 수행해". 수락된 제품 정본 `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md`의 순차 인계 2단위다. 1단위(U6 조회 복구)는 DEV 71dcac7·tester e1a92b3로 main 78bd0b3에 통합됐다. 이후 세션 메시지·운영 배포 1dd6800·소개 모션·CI 로고가 main에 들어왔다. 기준은 main 8d0db4b다.

1단위 후속으로 남은 것 중 이 단위와 겹치는 항목을 함께 처리한다: 조회 선택기의 종류 표시가 390·320px에서 잘림(tester e1a92b3 보고), 선택기 순서가 U2 이름 선두와 다름(사용자는 예외 수락; U2 이름 체계 통일 시 정합화), 익명 한도 429 문구 오안내 `internal/relay/member_receipt.go:27-31`(리뷰 low), HTTP 구 순서 단언 1건 갱신 필요.

## 적용 기준과 예외

- 규칙: fullops-common-0.3.3(`.fullops-squad/rules/common/README.md`와 연결 세 규칙), `.fullops-squad/project.md`, 문서 규칙 `.fullops-squad/docs/agents/document-writing.md`. 기준 커밋 8d0db4b.
- 제품 정본: SAR-UX-ISSUES-002.md의 U2·U4·U7, 사람 친구 D02 `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` HP-01–16, D04 `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`.
- 예외: `docs/agents/frontend.json`의 host codex/gpt-6.1-sol은 사용자 2026-10-17 Codex 제외 지시로 적용하지 않는다. 이 과제는 claude/claude-opus-5-5/high(user)다. 주 스킬 frontend-design과 D04 정본은 유지한다.
- Windows에는 make가 없다. 1단위와 같이 격리 원격 Linux clone에서 등록 `make lint`/`make test`를 실행하고 종료코드를 남긴다(방법은 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md` 참조).

## 먼저 읽을 문서

아래 탐색 패킷 절이 정본이다. 우선 순서: SAR-UX-ISSUES-002.md(U2·U4·U7·순차 인계 표 2행), D04, SAR-HUMAN-PEER.md, `internal/relay/human_pages.go`, `human_http.go`, `human.go`, `member_agents.go`, `member_pages.go`, `member_receipt.go`, `http.go`(gateView 이름 표시 영향만), 관련 `*_test.go`, `.fullops-squad/contexts/dev.md`.

지시 전제와 충돌 — 먼저 확인: `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md`(이슈 원문의 개별 에이전트 초대 요구와 다름; SAR-UX-ISSUES-002.md 적용 범위 절이 HP-01–16 우선을 확정), `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-RECEIPT-DEV.md`(1단위 범위 기록; 이번 단위 범위는 이 지시서가 정본). 다르게 읽히면 착수 전 ask.

Jev 분류 결과: `.fullops-squad/docs/evaluations/jev/SAR-UX-ISSUES-002-NAMES-DEV-context.json`. D04는 크기 때문에 분류되지 않았으므로 원문을 직접 읽는다.

## 해야 할 일과 파일 소유권

- [ ] 짧은 기술 계획(이 인박스 완료 보고에 기록)
- [ ] U2: 카드·선택·승인·조회·연결·legacy 화면에서 이름 먼저 + 클라이언트 + 비밀 아닌 안정 구분값. 구분값 충돌 시 확장. 권한상 이름 불가면 '확인할 수 없는 에이전트'+허용 식별자. 상대 비공개 프로필 추가 조회 금지. 사용자 지정 handle 미도입.
- [ ] U2: 내 계정·친구 추가 화면에 '내 친구 코드' 읽기 전용 입력(보이는 라벨, 키보드 전체 선택·복사). 복사 JS·CSP 완화 없음. 코드 앞뒤 공백만 제거, 빈 입력 로컬 형식 오류.
- [ ] U4: 친구 추가 화면 구성(내 이름·상대 코드·연결 범위·명시 동의·'친구 요청 보내기'), 상태 표(요청 제출 동일 문구·받은/보낸 대기·친구·종료·차단·에이전트 0개)대로 표시와 행동. 존재·차단 비노출 유지. 에이전트 0개라도 요청 가능, 빈 발신 선택기 금지.
- [ ] U7(해당 화면): 고유 title·h1·main·skip, 정확한 내부 복귀, status/alert, 반복 버튼 접근성 이름에 대상·구분값.
- [ ] 1단위 잔여: 조회 선택기 종류 잘림·순서 정합, 429 문구(member_receipt.go:27-31), HTTP 구 순서 단언 갱신.
- [ ] U2-A·U4-A의 핵심 실패 경계 자동 검사(동명·40자 한글/이모지·제어문자 거절·이름 변경 후 ID/관계 불변·세 생성 경로·코드 교체/차단/세대 경합·GET만으로 관계 생성 없음·agent 0개 요청). 문구 문자열만 복제하는 테스트는 쓰지 않는다.
- [ ] 영향받는 D05/D10/D11·D13과 필요한 사용 안내 갱신(route 추천 D10·D13·D11).
- [ ] 안정 후보에서 원격 Linux `make lint`·`make test` 한 번, FullOps `lint.py --from 8d0db4b` 실행·결과 보존.
- [ ] tester/designer용 격리 fixture(합성 A/B/C·동명·legacy·agent 0개, loopback origin/포트, 허용 GET/명시 POST, 수명)를 준비해 완료 보고에 연결.
- [ ] packet-outcomes 기록 → 완료 보고 → work.py finish → 커밋 → worker_done.

소유 파일: `internal/relay/` 회원·사람 친구 화면/handler/템플릿과 해당 테스트, 기술 문서(D05/D10/D11/D13), own route/find/context/packet/outcomes·QA dev 보고(`docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-DEV-test/`). 예상 규모 500–900줄, 10–15파일. SIZE 경고가 나면 실제 차이와 분할하지 않은 이유를 보고한다.

## 완료 기준과 검증

- U2-A·U4-A가 자동 검사 또는 fixture 실행으로 확인되고 실패 경계가 기록됐다.
- 기존 사람 친구 권한·서명·관계 세대·존재 비노출·계정 격리가 회귀 없이 유지된다(기존 테스트 통과).
- 원격 Linux `make lint`/`make test` exit 0, FullOps lint ERROR 0. 원 exit·WARNING·unavailable 사유를 보존한다.
- worker 완료는 제품 최종 수락이 아니다. 후속 독립 tester QA·designer 직접 시각·고정 SHA 리뷰·main 통합은 coor 책임이다.
- 멈추고 묻는 경우: 새 스키마/API/권한, 설명되지 않는 검증 실패, 제품 규칙 변경·범위 확대.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04 `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`, SAR-UX-ISSUES-002.md '화면 배치와 자기 비평'.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 기존 Go html/template·`humanAgentView`·native form/details. 새 테마·폰트·이미지·프레임워크·JS 없음.
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps `lint.py`(cwd 레포 루트), DESIGN 경고는 수정하거나 사유 기록.
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 테마 전환 없음. 1440·390·320px에서 이름·구분값·친구 코드 입력·상태 안내의 줄바꿈·넘침을 DEV가 가능한 범위로 확인하고 미실행은 not_run으로 남긴다.

## 갱신할 산출물

D10, D13, D11(route 추천). 실제 변경 원천에 따라 D05를 더한다. D02/D04 제품 정본은 수정하지 않는다.

탐색 패킷의 task/role/attempt·기준 SHA와 direct_edit/impact_check/document_read/document_update 경로·구간·근거·partial/unknown을 연결한다. 명시적 없음과 미확인을 구분한다. 필수·충돌·주의 후보와 budget 잔여 확인 목록을 보존한다.

## 기대 산출물

제품 코드·테스트, 기술 문서, `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-DEV-test/report.md`(검사 로그·fixture·not_run), packet-outcomes, 완료 보고.

## 제약·협업·후속

하위 작업은 [위임 기준](../rules/delegation.md)을 따른다(standard, 최대 2). 배정 여부와 이유를 보고한다.

다른 역할 인박스·PLANS·board·GitHub 이슈 상태는 수정하지 않는다. 승인 요약(U5)·연결 온보딩(U8)·언어(U3)는 다음 단위이므로 섞지 않는다. 피할 스타일: 상태 카드·통계 배지·그라디언트·알약 버튼·새 색 토큰·가운데 정렬 본문. D04의 canvas/ink/muted/accent/selected/edge·danger, 44px 제어, 4px 모서리를 유지한다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-UX-ISSUES-002-NAMES-DEV-packet.json` / SHA `8d0db4bf59e905c9f398123c4e2e5bc13ed4c3fa` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_read) · 줄 6, 10, 14, 26 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 29, 36, 38, 42, 44, 46, 46, 62, 62, 74, 75, 87, 87 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 55 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 39, 41 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 58 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 6, 38 · inferred
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` (document_read) · 줄 86, 111 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 14 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read, document_update) · 줄 7, 32, 38, 212, 212, 212, 295 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-DEV-test/runtime/report.md` (document_read) · 줄 257, 265 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-FINAL-MAIN-REVIEW-review/report.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-review/report.md` (document_read) · 줄 36 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-test/report.md` (document_read) · 줄 80 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN/report.md` (document_read) · 줄 6, 10, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 15, 24, 24, 58, 59, 60 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-FIX-test/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-RETEST-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-TESTER-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-DEV.md` (document_read) · 줄 33 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001.md` (document_read) · 줄 67 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001-VISUAL.md` (document_read) · 줄 26 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md` (document_read) · 줄 46, 48, 49 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 24, 33, 55 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-DESIGN.md` (document_read) · 줄 6, 10, 14, 38, 42 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-RECEIPT-DEV.md` (document_read) · 줄 6, 20 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read, document_update) · 줄 7, 35 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/protocol.md` (document_read) · 줄 209 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 261 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_designer.md` (document_read) · 줄 173 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_dev.md` (document_read) · 줄 136, 137, 140, 144, 152, 156, 185, 199 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_designer.md` (document_read) · 줄 69, 90 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_ops.md` (document_read) · 줄 123, 207, 317, 318, 319, 320, 321, 396 · inferred
- `.fullops-squad/handovers/logs/2026-10-18_to_dev.md` (document_read) · 줄 198, 199 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 16 · inferred
- `adapters/src/login.test.ts` (impact_check) · 줄 88, 95 · inferred
- `adapters/src/mcp.ts` (impact_check) · 줄 83 · inferred
- `docs/system-design.md` (document_read) · 줄 70, 70, 98 · inferred
- `internal/relay/connections_test.go` (impact_check) · 줄 335, 338, 415 · inferred
- `internal/relay/device.go` (impact_check) · 줄 55, 115, 121 · inferred
- `internal/relay/human.go` (direct_edit) · 줄 26, 39, 295 · inferred
- `internal/relay/human_http.go` (direct_edit) · 줄 21, 47, 237, 248 · inferred
- `internal/relay/human_pages.go` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/member_agents.go` (direct_edit) · 줄 61, 65, 68, 94, 119, 132, 140 · inferred
- `internal/relay/member_pages.go` (direct_edit) · 줄 31 · inferred
- `internal/relay/member_receipt.go` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/member_ui_test.go` (impact_check) · 줄 84 · inferred
- `internal/relay/policy_test.go` (impact_check) · 줄 22, 31, 59, 65, 68, 68 · inferred
- `internal/relay/receipt_context.go` (impact_check) · 줄 95, 100 · inferred
- `internal/relay/receipt_context_test.go` (impact_check) · 줄 82 · inferred
- `scripts/install_bot_mcp.sh` (impact_check) · 줄 83 · inferred
- `videos/knowslink-motion/index.html` (impact_check) · 줄 171, 181 · inferred
미확인 23건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
