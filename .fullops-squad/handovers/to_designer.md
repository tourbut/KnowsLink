---
title: SAR-HUMAN-PEER-001 — 사람 간 peer 아래 이름 있는 에이전트와 대화를 제공하도록 제품 규칙과 UI/UX를 전면 개편하고 DEV/QA 인계를 작성한다
status: draft
updated: 2026-10-12
owner: designer
tasks: [SAR-HUMAN-PEER-001]
summary: 사람 간 peer 아래 이름 있는 에이전트와 대화를 제공하도록 제품 규칙과 UI/UX를 전면 개편하고 DEV/QA 인계를 작성한다
attempt: ffaa6967c2b04a6daae94dfa27aa467c
base: f20b33a11addeb3debcb23e2e835be31c7dd1454
subagent_level: standard
test_level: lite
---

# SAR-HUMAN-PEER-001 — 사람 중심 연결 모델과 UI/UX 전면 개편

- 작성일: 2026-10-12
- From / To: coor / designer
- 상태: ready
- 테스트 레벨: lite. 하위 에이전트 레벨: standard.
- 담당: C:/Users/shin/orca/workspaces/KnowsLink/designer, fullops/designer.
- 병합: coor가 main 검토·통합·origin push를 맡는다.
- 복귀: repo 0b08ec4c-9e3a-4613-8197-5a835e545335, coor 워크트리, term_47da6145-4e1e-4cc1-93ea-20bc0dd590b9, Run run_5534969da2d0. task/dispatch는 실제 주입 preamble을 따른다.

## 사용자 원문과 목표

사용자: “UI가 뭐 바뀐게 없는거 같은데 좀 더 심플하게 하고 불필요한 문구는 안보이게하고 에이전트의 명칭도 정하게하고 실제 내 메일로 나의 에이전트들을 등록하고 내가 친구로 등록한 사람의 에이전트 목록들이 보여서 서로 대화르를 주고받는 형식으로 하자 별도로 에이전트끼리 관계를 등록한다기보다 사람과 사람이 서로 peer가 되면 그 하위 에이전트들은 연결이 되서 대화를 주고 받을수 있는 형식으로”
후속 사용자: “UI/UX도 한번 다시한번 디자인스킬써서 전면개편해보고”

앞선 polish는 간격·접근성 수정이었다. 이번 요청은 사람 간 peer 관계로 제품 모델을 바꾸고 간결한 UI/UX를 전면 재설계하는 것이다. 문서 작성으로 끝내지 않고 후속 DEV 구현·독립 QA·실제 시각 검수까지 이어간다. 현재 표시 중인 정적 fixture를 실제 동작으로 오인시키지 않는다.

## 적용 기준과 예외

기준 SHA f20b33a11addeb3debcb23e2e835be31c7dd1454. fullops-common-0.3.3의 README와 coding-style/testing/security, FULLOPS.md, orca-agents.md, project.md를 읽는다. 문서는 docs/agents/document-writing.md를 적용한다. 이 사용자 요청은 기존 에이전트 단위 관계 제품 규칙을 바꾸는 새 근거다. 과거 기준을 무조건 유지하거나 모순된 상태로 덧붙이지 말고 변경/유지/이행 범위를 명확히 한다. 기술 계획·구현 파일 배정의 세부 방식은 DEV 책임이다.

## 먼저 읽을 문서

필수: .fullops-squad/project.md; .fullops-squad/docs/design-docs/mockups/service-member-ui.md; .fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md; .fullops-squad/docs/planning/sources/silent-agent-relay/product.md; .fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md; .fullops-squad/docs/operations/user-guide.md.
현 상태 확인: internal/relay/member_agents.go, internal/relay/member_pages.go, internal/relay/test_messages.go, internal/relay/identity.go. 필요에 따라 일반 검색으로 보완한다. find의 낮은 점수·partial을 부재로 취급하지 않는다.

## 지시 전제와 충돌 — 먼저 확인

기존 product.md는 Contacts=agents 및 장기 채팅 이력 없음이 잠겨 있다. coor가 원문을 확인했다. 사용자의 2026-10-12 새 요청이 사람 간 peer와 대화 UI 변경의 명시적 근거다. 이 과제에서 designer가 어떤 기존 잠금을 교체하고 어떤 전달/보안/보존 규칙을 유지할지 정한다. 기술 구현으로 이 잠금을 몰래 우회하지 않는다. test_messages.go는 현재 메시지 동작의 실제 파일이다. 잘못 지정했던 member_messages.go는 존재하지 않아 제외하고 실패 기록을 보존한다.

## 해야 할 일과 파일 소유권

- [x] 사용자 흐름과 현재 소유/연결/메시지/이메일 인증 구조를 확인한다. 제품 정본의 충돌을 판정한다.
- [x] 이메일 소유 확인 아래 내 에이전트 등록·이름 지정/변경, 사람 간 친구 요청·상호 수락·거절·취소·해제·차단, 친구 에이전트 목록 노출과 대화 흐름을 확정한다. 노출 범위·새로 추가한 에이전트·철회·친구 해제 후 권한·과거 대화 접근·기존 관계 이행을 명시한다. 단순 친구 검색만으로 타인의 정보/에이전트를 공개하지 않는다.
- [x] frontend-design 스킬로 UI/UX를 전면 재설계한다. 대화·친구·내 에이전트의 정보 구조를 사용자 요구에 맞게 정하고, 간결한 텍스트와 이름 중심 목록 및 대화를 첫 경험에서 볼 수 있게 한다. 내·친구 에이전트 선택과 메시지 전달/수신·승인 상태를 이해할 수 있게 설계한다. 위험 안내와 실패 복구 문구는 필요한 곳에 유지한다.
- [x] D04 정본을 실질적으로 개편한다. 토큰·글꼴·레이아웃·desktop/mobile 와이어프레임·상태·주요 동작을 두 번의 디자인 계획/비평으로 기록한다. 기존 장식적 카드 묶음·긴 기술 설명·식별자 위주 첫 화면·상시 보안 설명·그라디언트 장식을 피한다. 취향만 바꾸는 CSS 수정으로 끝내지 않는다.
- [x] 현재 빈 DEV/tester 인박스에 후속 구현/독립 QA 지시서를 쓴다. fullops-work find/context/packet과 정확한 역할별 키 SAR-HUMAN-PEER-001-DEV / SAR-HUMAN-PEER-001-TESTER를 사용한다. DEV는 짧은 기술 계획→전체 구현→필수 회귀/권한 경계→기술 문서 갱신을 같은 과제에서 끝낸다. tester는 고정 DEV SHA를 기다린다.
- [x] D01/D02/D04와 관련 사용자 안내/실행 기록을 현재 제품 규칙에 맞게 갱신한다. D03/D05/D08/API/DB 이행 문서의 기술 갱신은 DEV에 인계한다. 과거 검증 기록을 새 구조 수락 증거로 쓰지 않는다.

소유권: .fullops-squad/docs/planning/, docs/design-docs/mockups/, docs/operations/user-guide.md의 제품 흐름 부분, docs/exec-plans/phases/SAR-HUMAN-PEER-001.md, handovers/to_dev.md, handovers/to_tester.md, contexts/designer.md 및 자기 탐색/결과. 제품 코드를 구현하지 않는다. coor 소유 PLANS/board는 coordinator에게 상태를 보고해 갱신시킨다.

## UI 실행 선택과 검증

주 스킬 frontend-design: C:/Users/shin/.agents/skills/frontend-design/SKILL.md. 호스트 codex, 모델 gpt-6.1-sol, effort high. 기존 D04 정본을 사용자 전면개편 근거로 변경한다. 선택형 도구 설치 없음. 기존 Go template/CSS와 Orca 내장 브라우저를 재사용한다. 전역 frontend.json enabled=false는 과거 기본값이며 이번 인박스에서 활성 디자인 과제를 지정한다. 설치/인증/브라우저 성공을 구분한다.

현재 단계는 제품/디자인 설계 검증이다. 문서·링크·실제 화면과 충돌 확인, 토큰과 와이어프레임의 사용자 요구 대조가 완료 조건이다. DEV 완료와 제품 수락은 별도다. 후속 tester가 핵심 실제 회원 두 명의 격리 fixture/API/DB 흐름과 권한 실패 경계, 320/390/1440폭, 키보드 조작을 확인하고 designer가 실제 구현 캡처로 전면 개편 결과를 직접 판정한다. 핵심은 이름 등록→친구 요청/수락→친구 에이전트 보임→대화 송수신→해제 시 신규 전달 차단이다. 실제 이메일/Google 계정이나 외부 사람에게 발송·운영 배포는 별도 실행 단계이며 임의 실행하지 않는다. 사용자의 실제 이메일은 저장소/로그에서 추측하거나 출력하지 않는다. 기존 인증을 실제 계정에 연결할 동작을 구현하고 실계정 수동 확인이 필요한 부분을 인계한다.

## 산출물과 완료

D01, D02, D04, D07 사용자 안내의 영향 확인/갱신. D03/D05/D08/D09/D10 영향과 담당을 확정한다. 사용자 요청을 불필요하게 쪼개 매번 승인을 요구하지 않는다. 최소한의 완결 기능 단위로 인계한다. 예상 규모는 여러 제품 문서와 역할 지시서이며 새 프레임워크/유료 도구 도입은 요구하지 않는다.

상세 기록: docs/exec-plans/phases/SAR-HUMAN-PEER-001.md. 패킷 모든 필수 path/category outcomes를 쓴다. 적용 규칙·확정/미확정·검증/미검증·기존 데이터 권한 이행·후속 역할을 완료 보고에 남긴다. 커밋한 뒤 work.py finish로 원본/보고를 보존하고 실제 worker_done에 최종 SHA 하나를 명시한다. 완료 뒤 독립 문서 검토와 main 통합 후 DEV를 착수시킨다.

## 제약

Windows 파일 삭제·이동 정리 금지, 임시 검증 자료 보존. 비밀값·실제 이메일·토큰을 출력하거나 커밋하지 않는다. 기존 권한·서명·CSRF·owner gate·rate limit을 임의 완화하지 않는다. 기존 관계 데이터는 증거 없이 사람 간 전체 권한으로 자동 확장하지 않는다. 기존 문서와 사용자 방향 충돌은 사용자 새 요청이 근거이며 합리적 제품 세부 결정은 designer가 진행한다. 불가피한 제품 모호성만 coor에 ask한다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-HUMAN-PEER-001-packet.json` / SHA `f20b33a11addeb3debcb23e2e835be31c7dd1454` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/operations/user-guide.md` (document_read, document_update) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/planning/sources/silent-agent-relay/product.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 50, 50, 53, 53, 65, 65 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `internal/relay/identity.go` (impact_check) · 줄 전체/미확인 · unknown · 필수
- `internal/relay/member_agents.go` (impact_check) · 줄 전체/미확인 · inferred · 필수
- `internal/relay/member_pages.go` (impact_check) · 줄 전체/미확인 · inferred · 필수
- `internal/relay/test_messages.go` (impact_check) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 72 · inferred
- `.fullops-squad/docs/agents/frontend.json` (direct_edit) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read, document_update) · 줄 119 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/database-design.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read, document_update) · 줄 141 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read, document_update) · 줄 22 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 35 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md` (document_read) · 줄 21 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 14, 31, 41 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read, document_update) · 줄 223 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL/report.md` (document_read) · 줄 28 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md` (document_read) · 줄 29 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (document_read) · 줄 26 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 15 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-VISUAL.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/generated/db-schema.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read, document_update) · 줄 40, 56 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read, document_update) · 줄 25 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_update) · 줄 전체/미확인 · inferred
미확인 15건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

제품 기획·UI/UX 설계·DEV/QA 인계 범위를 완료했다. 실제 기능·운영·새 구조의 제품 수락은 후속이다. 재개한 초안의 소유권은 coor msg_59d0bbe24839로 확인했고 기존 내용과 과거 실패/수락 기록을 보존했다.

확정 결정: 사람 친구를 공유 코드로 요청하고 상호 수락한다. 현재·향후 활성 agent를 이름으로 선택해 연속 대화한다. owner 작성과 agent 작성을 구분한다. 각 메시지의 전달 TTL은 최대 180s이고 새 text 원문 보존은 수락 후 최대 24h다. 친구 해제/차단·agent 철회가 모든 새 전달/실행/승인/result/legacy 경로보다 우선한다. 현재 회원의 본인 참여 과거 사본만 잔여 보존 기간 읽기 전용으로 허용한다. legacy pair는 정확한 기존 두 agent 범위만 유지하며 전체 사람 권한으로 자동 승격하지 않는다. 기술 구조·API·DB·웹 작성 인증·이행/복귀 방법은 DEV 책임이다.

디자인: frontend-design으로 기존 D04를 전면 개편했다. 대화 / 친구 / 내 에이전트 구조, 사람 아래 agent 이름 목록, 두 차례 계획·자기 비평, 6색 토큰·글꼴·간격·desktop/mobile 와이어프레임·상태·오류 복구·접근성과 수락 담당을 기록했다. 기술 ID·키·명령·계정은 필요한 상세로 옮겼다. 장식 카드·그라디언트·상시 보안 설명과 가짜 대화는 제외했다. 실제 보안/보존/실패 안내는 해당 행동 옆에 유지했다.

산출물: D01·D02·D04·D11 사용자 흐름 및 기획/원천 우선순위 안내를 갱신했다. 지시서의 D07 사용자 안내 표기는 인덱스와 달라 D07=DB 설계, D11=사용자설명서로 바로잡았다. D03/D05/D06/D07/D08/D09/D10 및 README/adapter 기술 갱신은 DEV다. D12/D13의 이행/복구/운영 실행은 DEV 기술 인계 뒤 OPS다. PLANS/board는 coor가 갱신한다.

검증: deliverables strict exit 0(13종, 문제/경고 0), 수정 뒤 git diff --check exit 0, 변경 Markdown 로컬 링크 169개/누락 0, 명세 토큰 대비 7개 조합 PASS다. 과거 product/decisions 원본문과 D04 이전 절 보존을 assert했다. frontend.py check는 configuration_ready exit 0이며 로그인·MCP·새 제품 브라우저 QA 성공을 뜻하지 않는다. designer packet completion 계약 exit 0, 비선택 path/category 52개와 unknown 15건·remaining 11개의 수동 검토를 outcomes에 기록했다. 최종 HEAD의 FullOps lint는 별도 검증 파일과 worker_done에서 고정하며 제품 테스트는 새 동작이 없어 해당 없음이다.

FullOps lint 실제 결과: 커밋 d6a4ef55df44에서 exit 1, ERROR 0·WARNING 3·unavailable 2다. 등록 make lint/make test는 Windows make 부재로 실행 불가이며 PASS가 아니다. 제품 코드·설정 변경이 없는 문서 과제로 제품 검사는 해당 없음이며 후속 DEV의 구현 고정 SHA에서 필수 실행한다. SIZE 경고는 inherited PLANS와 현행/과거 원문을 보존한 D04, 완결 제품 문서 추가 규모다. 분할하지 않은 이유는 phase에 기록했다. 원 결과는 C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-designer/pre-finish-lint.json, 최종 HEAD 재검사는 같은 폴더 final-lint.json에 남긴다.

후속 인계: to_dev.md는 SAR-HUMAN-PEER-001-DEV waiting-design-integration, to_tester.md는 SAR-HUMAN-PEER-001-TESTER waiting-dev-sha다. 각각 code/documents find·context·route bind·packet을 작성했다. packet 생성은 exit 0(DEV 후보65, tester203, 각각 partial/unknown17)이다. 준비 packet check는 상대 인박스 원천 hash 변경으로 exit 1이며 PASS가 아니다. coor가 문서 독립 리뷰·main 통합 뒤 DEV 실제 준비 HEAD의 패킷을 갱신하고, tester는 DEV 고정 구현 SHA에서 갱신해야 한다. tester의 검색 추천이 제품 코드 수정 권한을 주지 않는다.

미검증·보류: 실제 새 기능·키보드/200%/스크린 리더·실 Google/노우/다닷 새 흐름·운영 DB 이행/복원·배포는 실행하지 않았다. 기존 제품의 과거 PASS를 새 수락 근거로 쓰지 않는다. DEV 핵심/필수 검사, tester 독립 A/B/C와 legacy fixture QA, 같은 SHA의 designer 실제 구현 시각 검수 및 별도 세션 리뷰가 남았다. 실계정·외부 발송·운영 실행은 별도 단계다. 새 다크 모드·공개 검색·첨부·장기 보관·읽음·agent별 공유·상품 기능은 생략했다.

정본·상세 결과: [HP 요구사항](../docs/planning/product-specs/SAR-HUMAN-PEER.md), [현행 D04](../docs/design-docs/mockups/service-member-ui.md), [실행 기록](../docs/exec-plans/phases/SAR-HUMAN-PEER-001.md), [designer outcomes](../docs/evaluations/jev/SAR-HUMAN-PEER-001-packet-outcomes.json). 제품/임시 자료 삭제·이동은 실행하지 않았으며 모든 생성 탐색 JSON과 과거 캡처는 그대로 보존한다. 기획 본문은 5d5d26c부터 커밋했고 결과·finish 아카이브를 포함한 최종 SHA 하나를 실제 worker_done에 보고한다.
