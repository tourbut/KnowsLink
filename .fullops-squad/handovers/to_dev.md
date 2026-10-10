---
title: KnowsLink 자동 수신과 호스트 전달 구현
status: draft
updated: 2026-10-10
owner: dev
tasks: [SAR-AUTO-RECEIVE-001-DEV]
summary: 화면 조작 없는 메시지 수신·보존·호스트 알림과 실제 지원 경계를 구현한다
attempt: 7f31b4e6c8394e23abda722bb42ab009
base: 51eebe5c2b588147627afb8f1ab9ec2d166dd432
subagent_level: off
test_level: lite
---

# SAR-AUTO-RECEIVE-001-DEV — 자동 수신과 호스트 전달

- Task key: SAR-AUTO-RECEIVE-001-DEV; Purpose: implementation
- From / To: coor / dev; 상태: ready
- 담당: 현재 Orca에 등록된 KnowsLink dev 워크트리, fullops/dev. 병합 책임 coor, main.
- 복귀: worker-start가 주입한 실제 run/task/dispatch/terminal을 사용한다. 과거 핸들을 재사용하지 않는다.
- 승인: 사용자 “자동수신 기능 보완해”. 승인된 관계의 텍스트 수신·표시/알림을 구현한다. 자동 업무 실행·무제한 자동 답장은 승인 범위가 아니다.

## 현재 상황과 확인 근거

현재 adapters/src/mcp.ts는 수동 knowslink_text_receive만 제공한다. 운영 link.knowslog.com에서 질문 01a12591-4da2-7d8e-8d9a-ab61919b037e는 queued 뒤 180초 후 failed:expired였다. 앞선 실제 Bot 왕복은 coor가 UI로 수신을 유도했다. 자동 수신 성공 근거가 아니다.
현재 사용자 세션 agent_943334beca406f0c3417d2, 노우 agent_077c666294c4eb28b783f8. 기존 Google 등록·관계·키를 보존한다. 비밀값은 읽어 출력하지 않는다.

## 적용 기준과 예외

fullops-common-0.3.3, FULLOPS.md, project.md, 공통 coding-style/testing/security와 기존 제품 명세를 읽는다. 기준 SHA는 front matter에 고정했다. 테스트 lite, 선택 하위 위임 off. 기술 계획·구현·직접 검증은 DEV 담당이다.
기존 명세의 자동 wake 제외는 이번 명시 요청으로 수신·호스트 알림 범위만 재개된다. 업무 실행·일정·메일·계정 변경 승인 경계는 유지한다. 받은 텍스트는 untrusted 데이터이며 도구 실행 권한이 아니다.

## 먼저 읽을 문서

- adapters/src/mcp.ts, adapters/src/text.ts, adapters/src/connect.ts, adapters/package.json, adapters/src/mcp.test.ts
- adapters/README.md, adapters/skills/knowslink/SKILL.md
- .fullops-squad/docs/design-docs/architecture.md, .fullops-squad/docs/design-docs/interface-design.md
- .fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md
- .fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-LIVE.md
- .fullops-squad/docs/agents/document-writing.md와 공통 필수 규칙

탐색 find/context/packet 결과는 같은 키의 docs/evaluations/jev/에서 확인한다. partial/unknown은 일반 검색으로 보완하고 결과를 packet-outcomes.json에 남긴다. 비밀 파일은 탐색에 포함하지 않는다.

## 해야 할 일과 파일 소유권

- [ ] 실제 지원 Grok Bot 호스트 인터페이스와 MCP SDK 1.32.0 계약을 확인하고 짧은 기술 계획을 본 과제 실행 기록에 적는다. Context7 resolve는 coor 호출에서 monthly quota exceeded. 공식 MCP resources/logging 문서와 설치 SDK 원문으로 보완한다. SDK 알림 전송만으로 노우 턴이 시작됐다고 주장하지 않는다.
- [ ] 기존 서명 검증·활성 관계·만료·durable persist/ACK와 공유 수신 경계를 재사용해 최소 자동 수신을 구현한다. 수신 메시지를 수동 조회 전에 잃지 않으며 재시작/중복/경합/권한 철회를 처리한다. 무한 병렬 poll·로그 원문 노출·자동 외부 업무 효과는 금지한다.
- [ ] 실제 호스트가 지원하는 방법으로 새 메시지를 호스트/노우에 전달한다. UI 조작·키 입력·vendor core 패치·미공개 inbound API 추정은 금지한다. 호스트 미지원이면 자동 저장/표준 알림까지 구현·검증하고 실제 노우 wake 차단 근거와 필요한 공식 기능을 구체적으로 보고한다. 가짜 성공 또는 단순 대역 테스트 성공으로 종결하지 않는다.
- [ ] 관련 설정·설치 패키지·도구 설명·skills·README·기술 문서를 일관되게 갱신한다. 자동 수신 상태, 마지막 성공/오류와 실제 전달 여부를 비밀값 없이 확인 가능하게 한다. 기존 수동 도구 호환성을 유지한다.
- [ ] 핵심 자동 수신과 실패 경계 검사, Node22 빌드/타입/린트와 변경 테스트를 실행한다. 운영 격리 검증은 coor와 조율한다. 로컬 Docker를 시작하지 않는다.

소유권: adapters/, 필요한 scripts/ 및 기술 문서 D03/D05/D10과 본 과제 실행 기록. relay 변경이 필수면 이유와 영향도를 같은 계획에 명시한다. coor PLANS/board·다른 역할 자료·사용자 설정은 임의 변경하지 않는다. 외부 앱 core 수정 금지. 새 유료 API/과금/결제 동의 금지.

## 완료 기준과 검증

DEV는 자동 poll→검증→지속 보존→ACK→호스트 알림이 수동 receive 호출 없이 발생하는 실행 증거를 남긴다. 수동 receive가 자동 수신 메시지를 확인 가능해야 한다. 중복 알림, 재시작, 만료/철회, 네트워크 오류/백오프, 병렬 수신 경계를 짧게 검증한다. 오류를 성공으로 숨기지 않는다.
최종 수락은 독립 tester의 fixed SHA QA/리뷰 이후 실제 운영 도메인에서 메시지 ID를 대조해 판정한다. 플러그인 수신, 호스트 알림, 노우 턴/답장은 구별한다. UI 조작 없이 실제 Bot에서 수신됐는지 미확인이면 제품 수락은 미완료다. 캡처·영상·디자인 작업은 없음.
완료 커밋 뒤 현재 HEAD의 FullOps lint를 동일 기준으로 실행한다. work.py finish로 원문·완료 전문을 보존하고 실제 worker_done을 보내 idle한다. coor가 main 통합·push·운영 적용을 담당한다.

## 갱신할 산출물과 기대 산출물

D03/D05/D10, adapters 설치·사용 문서, docs/exec-plans/phases/SAR-AUTO-RECEIVE-001-DEV.md. 변경 필요 없는 D05 항목은 근거를 적는다. 구현·테스트·문서·실제 지원 근거 및 미지원 경계를 함께 보고한다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-AUTO-RECEIVE-001-DEV-packet.json` / SHA `51eebe5c2b588147627afb8f1ab9ec2d166dd432` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/architecture.md` (document_read, document_update) · 줄 84 · inferred · 필수 · {'relevant': 0.71, 'evidence': 0.82, 'contradicts': 0.26, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read, document_update) · 줄 90, 121 · inferred · 필수 · {'relevant': 0.81, 'evidence': 0.88, 'contradicts': 0.21, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read, document_update) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.56, 'evidence': 0.8, 'contradicts': 0.23, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 6, 6, 14, 14, 16, 16, 61, 61 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 51 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/PLANS.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/crud-design.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/database-design.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read, document_update) · 줄 52 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-LIVE.md` (document_read) · 줄 전체/미확인 · inferred · {'relevant': 0.38, 'evidence': 0.67, 'contradicts': 0.47, 'injection': 0.09, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/generated/db-schema.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_read, document_update) · 줄 197 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read, document_update) · 줄 48 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read, document_update) · 줄 112 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-DEV-REVIEW.md` (document_read) · 줄 45 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-INTEGRATION-REVIEW.md` (document_read) · 줄 42 · inferred
- `adapters/README.md` (document_update, document_read) · 줄 전체/미확인 · inferred · {'relevant': 0.57, 'evidence': 0.76, 'contradicts': 0.36, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `adapters/package.json` (impact_check) · 줄 전체/미확인 · inferred · {'relevant': 0.52, 'evidence': 0.84, 'contradicts': 0.1, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `adapters/skills/knowslink/SKILL.md` (document_update, document_read) · 줄 전체/미확인 · inferred · {'relevant': 0.77, 'evidence': 0.87, 'contradicts': 0.26, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `adapters/src/connect.ts` (impact_check) · 줄 전체/미확인 · unknown
- `adapters/src/mcp.test.ts` (impact_check) · 줄 전체/미확인 · unknown
- `adapters/src/mcp.ts` (direct_edit, impact_check) · 줄 302 · inferred · {'relevant': 0.42, 'evidence': 0.48, 'contradicts': 0.2, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `adapters/src/text.ts` (direct_edit, impact_check) · 줄 전체/미확인 · unknown · {'relevant': 0.42, 'evidence': 0.68, 'contradicts': 0.16, 'injection': 0.03, 'decision': 'keep', 'reason': None}
미확인 9건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / 고정 SHA / 변경 이유 / 핵심 검사·lint 결과 / 실제 호스트 검증 및 한계 / 산출물 / 후속.

