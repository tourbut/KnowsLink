---
title: SAR-HUMAN-PEER-001-DEV — 사람 친구와 이름 있는 에이전트의 대화 UI 전체 구현 및 안전한 데이터 이행
status: draft
updated: 2026-10-12
owner: dev
tasks: [SAR-HUMAN-PEER-001-DEV]
summary: 사람 친구와 이름 있는 에이전트의 대화 UI 전체 구현 및 안전한 데이터 이행
attempt: 091fb01ddedc4dcda2db743061393dfb
base: f20b33a11addeb3debcb23e2e835be31c7dd1454
subagent_level: standard
test_level: lite
---

# SAR-HUMAN-PEER-001-DEV — 사람 친구·이름·대화의 전체 구현과 안전 이행

- 상태: waiting-design-integration. 현재 설계의 독립 문서 리뷰·main 통합 뒤 coor가 준비 SHA와 실행 tuple을 고정하고 ready로 바꾼다.
- From / To: designer / dev. Purpose: implementation. Task key: SAR-HUMAN-PEER-001-DEV.
- 담당: C:/Users/shin/orca/workspaces/KnowsLink/dev, fullops/dev. 제품 코드·기술 문서·관련 테스트·자기 증거 소유.
- 복귀: repo 0b08ec4c-9e3a-4613-8197-5a835e545335, coor 워크트리, term_47da6145-4e1e-4cc1-93ea-20bc0dd590b9, Run run_5534969da2d0. 실제 실행 task/dispatch는 새 preamble을 사용한다.
- 새 구현은 이 인박스 원문과 준비 커밋을 동일하게 받은 뒤 시작한다. 이 문서의 base는 기획 기준이며 dispatch 직전 coor가 실제 고정 기준과 패킷을 갱신한다.

## 현재 상황

기존 회원 UI는 agentID·키·agent pair 관리다. public_text.go는 signed text이며 웹 작성은 없다. identity.go는 issuer/sub 신원을 사용하고 Member/Agent에 표시 이름이 없다. 사용자는 사람 친구 아래 agent가 자동 연결되고 웹에서 대화하는 구조를 요청했다. 간격/CSS만 고치는 과제가 아니다.

## 해야 할 일과 파일 소유권

- [ ] 기존 코드·인가 호출자·테스트·설치 흐름을 확인하고 짧은 기술 계획을 자기 phase에 쓴다. API/저장/웹 발송의 실제 주체·legacy 권한/세대·원문 보존/rollback 경계를 결정한다.
- [ ] HP-01–07: 회원/agent 이름, Google 신원 연속성, 친구 코드·요청·수락·거절·취소·해제·차단, 친구 agent 목록·새 agent 노출·자기 agent 연결을 구현한다.
- [ ] HP-08–12: 친구 아래 agent 두 개의 대화, 웹 owner 작성과 agent 작성의 구분, 발신/수신 선택, 실제 전달/수신·새로고침·상태/승인 진입·멱등·4096 bytes/180s TTL·최대24h 본문 보존을 구현한다. 로컬 개인키를 브라우저/서버로 옮기거나 agent 서명을 위조하지 않는다.
- [ ] HP-13–15: 해제/차단/철회가 enqueue·lease·ACK·claim·gate consume·result·SSE·remote·legacy 모두에 적용되도록 공유 인가를 완결한다. 과거 대화는 참여자 자기 복사본만 잔여24h 읽기 전용으로 허용한다.
- [ ] HP 이행 표의 이전 JSON/agent/key/pair/pending/receipt/gate/합성 owner를 격리 DB로 검증한다. cross-owner 자동 친구 확대를 금지하고 기존 정확한 pair 범위를 보존한다. legacy 신규 pair 경로는 새 권한을 만들지 못하게 한다.
- [ ] 현행 D04 전면 UI를 로그인·이름·대화·친구·agent·계정·연결 승인·gate·receipt/오류에 적용한다. 기존 direct link/권한·명령 사용자를 끊지 않도록 읽기/복귀 경로를 제공한다.
- [ ] 핵심 자동 검사·필수 제품 lint/test와 필요한 회귀를 완료한다. 구현에 맞게 기술 문서/README/adapter 사용자 안내를 갱신한다.

제품 코드 전체(cmd/internal/adapters/db/scripts 및 필요한 기존 설정)와 기술 정본은 dev 소유다. API/schema/구조/파일 세부 배정은 DEV가 정한다. coor PLANS/board와 designer HP/D04 의미를 임의 수정하지 않는다. 규모는 여러 relay/adapter/test·기술 문서에 걸친 기능 변경으로 예상한다. 기능 완결/권한 이행을 나누지 않은 이유와 SIZE 경고를 기록한다. 새 프레임워크/의존성/유료 도구는 요구하지 않는다.

## 완료 기준과 검증

현행 D02의 '대화의 연속성과 과거 열람 경계'를 적용한다. 같은 대화에서 세 번 이상 왕복하고 앞선 메시지 TTL 경과 뒤 새 발송도 성공해야 한다. 각 발송의 TTL·24h 보존은 독립적이다. 과거 열람은 유효 회원 세션의 참여 기록만 허용하며 철회 키나 옛 gate의 권한을 복구하지 않는다. 친구 코드 교체와 응답·한도 안내의 존재 추측 방지도 검증한다.

HP-01–16을 구현/검사에 매핑한다. 핵심 happy path는 이름 저장→두 회원 친구 수락→새 agent 보임→양방향 text→해제 후 신규 전달 차단이다. 필요한 음성 경계는 미수락/C 제삼자/발신 위조/CSRF/차단·옛 pair/철회/재수락 옛 세대/중복·충돌/4096bytes·180s·24h/이전 데이터 보존이다. signed transport와 업무 gate의 기존 필수 검사를 유지한다. 모든 커밋 뒤 FullOps lint.py --from <실제 기준 SHA>가 product-lint·product-test를 실행해 같은 HEAD에서 통과해야 한다. 명령 자신의 종료코드를 기록한다.

D04의 두 자기 비평을 구현에 적용하고 실제 동작 캡처를 만든다. 정적 template fixture PASS와 실제 auth/API/DB PASS를 분리한다. DEV 완료는 전체 제품 수락이 아니다. 고정 SHA·검사/캡처·환경·미실행을 tester와 designer에게 인계한다.

## 갱신할 산출물

D03 architecture/tech-stack, D05 interface, D06 data-model, D07 database-design, D08 실제 schema/생성 여부, D09 CRUD, D10 module-design을 실제 구현에 맞춰 갱신한다. README·project.md·adapters README·D11 기술 절차를 갱신한다. D12/D13은 운영 실행이 필요한 내용과 담당/재개 조건을 OPS로 인계한다. D07은 DB 설계이며 사용자 안내는 D11이다. 새 테이블이 없으면 D08 변경 불필요 근거를 남긴다.

## 적용 기준과 예외

fullops-common-0.3.3의 README·coding-style·testing·security, FULLOPS.md, orca-agents.md, project.md, docs/agents/document-writing.md를 직접 읽는다. 기준 ref는 f20b33a11addeb3debcb23e2e835be31c7dd1454다. 테스트 lite·하위 위임 standard를 유지한다. 제품 변경 근거는 2026-10-12 사람 peer·이름·대화 및 frontend-design 전면 개편 사용자 요청이다. 보안·필수 검사·독립 리뷰 예외는 없다.

## 먼저 읽을 문서

필수: docs/planning/product-specs/SAR-HUMAN-PEER.md, docs/design-docs/mockups/service-member-ui.md의 현행 절, docs/planning/product-specs/SAR-PUBLIC-SERVICE.md의 현행 안내, docs/exec-plans/phases/SAR-HUMAN-PEER-001.md, docs/operations/user-guide.md의 새 흐름, project.md와 공통 규칙. 경로는 .fullops-squad/ 기준이다.

코드 시작점: internal/relay/member_agents.go, member_pages.go, identity.go, public_text.go, test_messages.go, store.go, member.go. 필요에 따라 현재 인가의 모든 호출자, adapter local/remote·SSE·gate 경로를 일반 검색으로 보완한다. test_messages.go는 시험 경로다. public_text.go가 일반 signed text 경로다. 존재하지 않는 member_messages.go를 전제로 삼지 않는다. 라이브러리 API/설정 확인은 Context7 resolve-library-id 다음 query-docs를 사용한다. 문서·패킷 partial/unknown을 부재로 해석하지 않는다.

## 지시 전제와 충돌 — 먼저 확인

과거 원천의 Contacts=agents·agent별 첫 수락·composer 금지는 현행 HP/D04가 대체한다. 업무 relay.v1·서명·현재 인가·인간 게이트는 유지한다. 새 대화 text의 24h 보존만 기존 완료 뒤 원문 삭제를 바꾼다. 과거 메시지·업무 M/H/R을 소급 보존하지 않는다. 새 사람 친구 권한이 기존 pair를 자동 확대하지 않도록 HP 이행 표를 적용한다. 과거 리뷰/QA PASS는 새 기능의 증거가 아니다.

## 협업과 실행 제약

삭제·이동 정리는 금지한다. 임시 검증 자료는 보존하고 경로를 보고한다. 실제 이메일/Google 계정·외부 상대 발송·운영 배포·DNS/신뢰 설정·비밀값 출력은 이 지시의 실행 범위가 아니다. 실계정 확인이 필요한 항목은 로컬 합성 검증과 분리해 인계한다. 이미 위임된 제품 세부 기준은 반복 승인받지 않는다. 기술 판단은 DEV가 처리하며 설명되지 않은 제품 충돌만 coor에게 ask한다.

선택형 하위 작업은 rules/delegation.md를 읽고 Purpose: research|review|implementation, 테스트 lite와 파일 소유권을 선언한다. 최대 두 개이며 부모가 회수한다. 독립 리뷰를 자기 하위 자기 리뷰로 대체하지 않는다.

Windows make 부재 참고: coordinator msg_907ca28147c9가 C:/Users/shin/AppData/Local/KnowsLinkDevTools/frontpolish001-coor/remote.py 및 make.cmd의 기존 Linux 검증 경로를 안내했다. ROOT가 coor이므로 그대로 실행하지 않는다. 자기 대상 SHA/ROOT·known hosts·비밀값 미출력과 실제 동기화 범위를 확인한 후 기술 담당자가 재사용 여부를 결정한다. 과거 helper 존재는 새 후보의 검사 통과가 아니다.

## UI 선택과 검증 계약

제품 디자인 정본은 docs/design-docs/mockups/service-member-ui.md의 현행 SAR-HUMAN-PEER-001 절이다. 사용자 선택 frontend-design, codex / gpt-6.1-sol / high를 유지한다. 다른 호스트는 coor가 실제 지원 tuple을 명시하기 전 자동 전환하지 않는다. 새 도구 설치 없음. 기존 Go templates/CSS·Orca 브라우저를 재사용한다. memberStyle은 기존 공용 스타일 시작점이며 현행 D04 토큰으로 바꾼다. DESIGN 검사가 Go 문자열 CSS를 판정하지 못하는 한계를 기록한다. 다크 테마는 범위 밖이다.

장식 카드 묶음·긴 기술 설명·ID 위주 홈·상시 보안 배너·그라디언트·알약 버튼·가짜 대화/온라인을 만들지 않는다. 보안/보존/실패 안내는 실제 동작 옆에 유지한다. 320/390/1440폭, 44px 동작, 포커스, 키보드, 이름 중복/긴 이름, 한국어 줄바꿈을 확인한다.

웹 QA는 자기 loopback origin·고유 포트/fixture/DB·정확한 허용 GET/명시 POST·독립 assertions를 기록한다. 운영 origin/외부 발송은 금지한다. 권장 예산은 핵심 동작 시나리오 120 steps·20분이며 예산 중단은 not_run이다. 필수 보안 검사는 예산으로 생략하지 않는다. E2E/Orca 실행기 버전, 요청/시나리오 지문과 실제 완료 범위를 보고한다. 탐색·shadow·CUAWright는 OFF다. 캡처는 시각 항목만 만든다.

## 완료와 복귀

진행은 자기 인박스 한 곳에 체크한다. 상세 로그와 packet-outcomes에 모든 비선택 path/category의 completed/no_change와 구체적 이유를 쓴다. partial 잔여를 수동 검토한다. 결과·검증·미검증·후속을 완료 보고에 전문으로 기록하고 커밋·work.py finish를 수행한다. task/dispatch/terminal은 실제 주입 preamble을 따른다. worker_done은 3문장과 최종 SHA 하나로 정확히 한 번 보낸다. coor가 고정 후보 리뷰·main 통합·origin push를 조정한다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-HUMAN-PEER-001-DEV-packet.json` / SHA `5d5d26cd02e98a6a0978ca681a7565a6b6e5a29d` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read, document_update) · 줄 7, 7, 13, 13, 15, 15, 381, 381 · inferred · 필수 · {'relevant': 0.91, 'evidence': 0.92, 'contradicts': 0.27, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001.md` (document_read) · 줄 2, 2, 6, 6, 10, 10, 30, 30, 65, 65, 67, 71, 71 · inferred · 필수 · {'relevant': 0.9, 'evidence': 0.9, 'contradicts': 0.32, 'injection': 0.07, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/operations/user-guide.md` (document_read, document_update) · 줄 7, 7, 13, 13, 15, 15 · inferred · 필수 · {'relevant': 0.39, 'evidence': 0.39, 'contradicts': 0.42, 'injection': 0.15, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 7, 7, 15, 15 · inferred · 필수 · {'relevant': 0.93, 'evidence': 0.94, 'contradicts': 0.2, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 7, 7, 15, 15 · inferred · 필수 · {'relevant': 0.24, 'evidence': 0.38, 'contradicts': 0.38, 'injection': 0.07, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 17, 17, 56, 56, 74, 74 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.51, 'evidence': 0.82, 'contradicts': 0.27, 'injection': 0.06, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `internal/relay/test_messages.go` (impact_check) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 55 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 6, 14, 14, 16, 16 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read, document_update) · 줄 전체/미확인 · inferred · {'relevant': 0.81, 'evidence': 0.87, 'contradicts': 0.29, 'injection': 0.06, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read, document_update) · 줄 전체/미확인 · inferred · {'relevant': 0.75, 'evidence': 0.88, 'contradicts': 0.28, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/design-docs/data-model.md` (document_read, document_update) · 줄 전체/미확인 · inferred · {'relevant': 0.77, 'evidence': 0.91, 'contradicts': 0.27, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/design-docs/database-design.md` (document_read, document_update) · 줄 전체/미확인 · inferred · {'relevant': 0.64, 'evidence': 0.82, 'contradicts': 0.29, 'injection': 0.07, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read, document_update) · 줄 전체/미확인 · inferred · {'relevant': 0.76, 'evidence': 0.89, 'contradicts': 0.42, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read, document_update) · 줄 7, 7, 16, 16 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 6, 6, 14, 14 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read, document_update) · 줄 145, 145, 145, 228 · inferred · {'relevant': 0.48, 'evidence': 0.81, 'contradicts': 0.25, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PREP-002-review/report.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-SERVICE-001-REVIEW-review/report.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-DESIGN-REVIEW-review/report.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-review/report.md` (document_read) · 줄 36 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-test/report.md` (document_read) · 줄 80 · inferred
- `.fullops-squad/docs/evaluations/scenarios/README.md` (document_read) · 줄 전체/미확인 · inferred · {'relevant': 0.42, 'evidence': 0.73, 'contradicts': 0.15, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 24, 33, 55 · inferred
- `.fullops-squad/docs/generated/db-schema.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_read, document_update) · 줄 전체/미확인 · inferred · {'relevant': 0.14, 'evidence': 0.29, 'contradicts': 0.4, 'injection': 0.09, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 6, 6, 12, 12, 14, 14 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 7, 7, 18, 18 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 7, 7, 17, 17 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 7, 7, 16, 16 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md` (document_read) · 줄 6, 6, 10, 10 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/product.md` (document_read) · 줄 6, 6, 10, 10 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/protocol.md` (document_read) · 줄 209 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 261 · inferred
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 50, 50, 53, 53, 65, 65, 74, 74 · inferred
- `.fullops-squad/handovers/to_tester.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 17, 17, 24, 24, 32, 32, 54, 54, 72, 72 · inferred
- `adapters/README.md` (document_read) · 줄 전체/미확인 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 16 · inferred
- `adapters/src/login.test.ts` (impact_check) · 줄 88, 95 · inferred
- `adapters/src/mcp.ts` (impact_check) · 줄 78 · inferred
- `cmd/relay/main.go` (impact_check) · 줄 전체/미확인 · inferred
- `db/queries/relay.sql` (impact_check) · 줄 전체/미확인 · inferred
- `internal/config/config.go` (impact_check) · 줄 전체/미확인 · inferred
- `internal/relay/connections_test.go` (impact_check) · 줄 335, 338, 409 · inferred
- `internal/relay/device.go` (impact_check) · 줄 55, 111, 117 · inferred
- `internal/relay/identity.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/member.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/member_agents.go` (direct_edit, impact_check) · 줄 61, 65, 68, 94, 119, 132, 140 · inferred
- `internal/relay/member_pages.go` (direct_edit, impact_check) · 줄 62 · inferred
- `internal/relay/member_ui_test.go` (impact_check) · 줄 84 · inferred
- `internal/relay/policy_test.go` (impact_check) · 줄 22, 31, 59, 65, 68, 68 · inferred
- `internal/relay/public_text.go` (direct_edit, impact_check) · 줄 30 · inferred
- `internal/relay/store.go` (direct_edit) · 줄 전체/미확인 · unknown
- `scripts/install_bot_mcp.sh` (impact_check) · 줄 83 · inferred
- 작업자 의견(원본 추천 유지): `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` change — designer가 확정한 HP 제품 정본이다. DEV는 구현 추적만 하고 제품 의미 변경은 coor에 전달한다.
- 작업자 의견(원본 추천 유지): `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` change — 현행 D04를 구현한다. 과거 비현행 절이나 과거 QA를 새 수락으로 사용하지 않는다.
- 작업자 의견(원본 추천 유지): `.fullops-squad/docs/operations/ops-guide.md` change — 새 권한/보존/rollback 기술 영향만 DEV가 인계한다. 운영 실행은 OPS 후속이다.
- 작업자 의견(원본 추천 유지): `.fullops-squad/docs/operations/transition.md` change — 실제 운영 이행은 이번 DEV/QA 실행 범위가 아니다. 고정 수락 후보 뒤 OPS가 수행한다.
미확인 17건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

미완료. 고정 SHA·검증·미실행·후속을 실제 완료 때 작성한다.
