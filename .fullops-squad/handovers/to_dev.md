---
title: SAR-SERVICE-UI-001-DEV — 로그인·내 에이전트·연결 관리 웹의 상용 디자인 구현
status: draft
updated: 2026-10-11
owner: dev
tasks: [SAR-SERVICE-UI-001-DEV]
summary: 로그인·내 에이전트·연결 관리 웹의 상용 디자인 구현
attempt: 69804b27a2274875a19179a73e833ae9
base: b6a4f4da545cfdd3f57862f8724f3eacbe00c8e2
subagent_level: standard
test_level: lite
---

# SAR-SERVICE-UI-001-DEV — 회원 웹 상용 UI 구현

- 상태: ready; Test level: lite; Subagent level: standard
- 담당 repo: 0b08ec4c-9e3a-4613-8197-5a835e545335; 워크트리 C:/Users/shin/orca/workspaces/KnowsLink/dev; fullops/dev
- 복귀: run_7176a2ad6ef1, coor term_e61d3e14-29e9-4954-943a-4a75707c82de. task/dispatch/worker_done은 실제 preamble.
## 적용 기준과 예외
fullops-common-0.3.3: rules/common/README.md와 coding-style.md/testing.md/security.md, FULLOPS.md/project.md/contexts/dev.md/docs/agents/document-writing.md를 읽는다. Ponytail full. 새 유료 서비스/키 결제/로컬 Docker 금지. 기존 DB/키/agent/pair/Google/Tunnel 자격/공유 서비스와 F: 사고 증거 보존. 삭제·이동·컨테이너 정리 금지. 비밀값 로그/Git/채팅 금지. 웹은 Orca 내장 브라우저만, Computer Use는 Grok Bot 설정만이며 수신 확인 위해 Grok UI를 열지 않는다. 아직 webhook 후보는 main 미수락·미배포다. 운영 변경은 OPS/coor 후속.
## 완료 기록
본인 phase/contexts/QA/packet-outcomes를 기록하고 work.py finish로 전문 보존·인박스 비움·커밋 후 fullops lint의 base/head/ERROR/WARNING/실행불가 및 실제 명령 종료코드를 보고한다. 정상 작업은 중간 확인 없이 끝까지 진행한다. 파괴/범위 밖/설명 안 되는 필수 실패만 ask/escalation한다. PLANS/board는 coor 소유다.
<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-SERVICE-UI-001-DEV-packet.json` / SHA `b6a4f4da545cfdd3f57862f8724f3eacbe00c8e2` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.25, 'evidence': 0.57, 'contradicts': 0.27, 'injection': 0.07, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/agents/document-writing.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.32, 'evidence': 0.48, 'contradicts': 0.14, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read, document_update) · 줄 7, 213, 213 · inferred · 필수 · {'relevant': 0.9, 'evidence': 0.9, 'contradicts': 0.27, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 6, 14 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.78, 'evidence': 0.89, 'contradicts': 0.24, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.31, 'evidence': 0.57, 'contradicts': 0.17, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.78, 'evidence': 0.86, 'contradicts': 0.1, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.73, 'evidence': 0.86, 'contradicts': 0.11, 'injection': 0.09, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.7, 'evidence': 0.85, 'contradicts': 0.14, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `internal/relay/member_receipt.go` (impact_check) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.37, 'evidence': 0.48, 'contradicts': 0.23, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `.fullops-squad/PLANS.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 55 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 30, 32, 70, 72 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read, document_update) · 줄 119 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/database-design.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read, document_update) · 줄 141 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read, document_update) · 줄 7 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 18, 42 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 35 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md` (document_read) · 줄 21 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 6, 14, 31, 41 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read, document_update) · 줄 145, 145, 145 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/COOR/dev-fix-2-final/coordinator-handover-supplement.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 16, 17 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW-review/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-REVIEW-review/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-SERVICE-001-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6 · inferred · {'relevant': 0.62, 'evidence': 0.71, 'contradicts': 0.28, 'injection': 0.08, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/evaluations/scenarios/README.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-3.md` (document_read) · 줄 15 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (document_read) · 줄 6, 10, 26, 28, 32, 48, 49, 66 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPEN-PREP.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SETUP-001.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/generated/db-schema.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_read, document_update) · 줄 7, 305, 307, 312, 336 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read, document_update) · 줄 7, 89 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 6, 12, 14, 20, 23, 24, 27 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 7, 40, 56, 66 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 7, 25, 83, 154 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 7, 15, 23, 63, 76, 125, 128, 152 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/protocol.md` (document_read) · 줄 209 · inferred
- `.fullops-squad/handovers/_TEMPLATE.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_designer.md` (document_read) · 줄 9, 16, 20, 79, 87, 114 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_dev.md` (document_read) · 줄 16, 33, 35, 37, 51, 93, 125, 181 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_ops.md` (document_read) · 줄 90, 93, 97, 101, 116, 130, 144, 147, 151, 155, 171 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_tester.md` (document_read) · 줄 16, 33, 35, 37, 49 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 35, 102, 105, 106, 142, 248, 249, 261, 306, 386 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 37, 176, 245, 246, 304, 357 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_tester.md` (document_read) · 줄 100, 175, 176, 228 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_designer.md` (document_read) · 줄 29, 69, 76, 133, 182, 183, 184 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 24, 108, 113, 148, 149, 229, 230 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_tester.md` (document_read) · 줄 48, 70, 75, 139, 180, 181, 182 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_designer.md` (document_read) · 줄 257, 258, 268, 304, 305, 308, 312, 320, 337, 343, 365, 377, 391, 392, 393, 400, 433 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_ops.md` (document_read) · 줄 204, 205, 216 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_tester.md` (document_read) · 줄 254, 361, 362, 363, 375 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-DEV-REVIEW.md` (document_read) · 줄 45 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-INTEGRATION-REVIEW.md` (document_read) · 줄 42 · inferred
- `Makefile` (impact_check) · 줄 전체/미확인 · inferred
- `adapters/README.md` (document_read) · 줄 전체/미확인 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 16 · inferred
- `adapters/src/login.test.ts` (impact_check) · 줄 88, 95 · inferred
- `adapters/src/mcp.ts` (impact_check) · 줄 78 · inferred
- `cmd/relay/main.go` (impact_check) · 줄 전체/미확인 · inferred
- `internal/config/config.go` (impact_check) · 줄 전체/미확인 · inferred
- `internal/relay/connections_test.go` (impact_check) · 줄 335, 338, 407 · inferred
- `internal/relay/device.go` (impact_check) · 줄 55, 111, 117 · inferred
- `internal/relay/http.go` (direct_edit, impact_check) · 줄 전체/미확인 · unknown · {'relevant': 0.34, 'evidence': 0.44, 'contradicts': 0.18, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `internal/relay/member.go` (direct_edit, impact_check) · 줄 전체/미확인 · unknown · {'relevant': 0.62, 'evidence': 0.68, 'contradicts': 0.18, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `internal/relay/member_agents.go` (direct_edit, impact_check) · 줄 36, 40, 43, 69, 94, 107, 115 · inferred · {'relevant': 0.59, 'evidence': 0.72, 'contradicts': 0.15, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `internal/relay/policy_test.go` (impact_check) · 줄 22, 31, 59, 65, 68, 68 · inferred
- `scripts/install_bot_mcp.sh` (impact_check) · 줄 83 · inferred
미확인 9건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고
완료 시 실제 SHA, 검증 결과와 미검증, 파일과 증거, 후속을 전문으로 기록한다.

## 목표와 범위
main 수락 디자인 service-member-ui.md(R-01–11, A-01–10, I-01–08)에 맞게 Google 로그인·이메일 로그인(설정 시)·회원 홈·연결 관리·기기 승인·receipt/gate/오류의 상용 UI를 구현한다. 기술 계획과 구조는 DEV가 결정한다. 기존 SSR·공통 CSS/기본 HTML부터 재사용하고 새 프레임워크/라이브러리/외부 폰트/이미지/JS를 추가하지 않는다.
## 해야 할 일과 파일 소유권
- [ ] frontend-design 실제 스킬 C:/Users/shin/.agents/skills/frontend-design/SKILL.md 및 제품 디자인 정본을 읽고 영향 template·caller·test·fixture를 확인해 짧은 기술 계획을 같은 과제에 기록한다.
- [ ] 로그인 서비스 소개·Google 주 동작, 에이전트 우선 목록/접힘 설정·위험 구역, 공통 타입/색/반응형, 한국어 줄바꿈·접근성을 구현한다. 장식용 통계·알약 버튼·무작위 그라디언트·같은 카드의 반복은 피한다. 온라인/최근 접속 같은 없는 사실을 만들지 않는다.
- [ ] form action/method/hidden 값·CSRF·CSP·owner 분리·5분 재인증·키/관계/한도 규칙을 보존한다. 디자인의 R-10/11은 저비용 기존 데이터로 가능할 때 포함한다. 보안상 필요한 지문/경고를 접어서 놓치지 않게 한다.
- [ ] 실제 고정 후보를 Linux 고유 체크아웃으로 product lint/test, 핵심 폼/owner/CSP 회귀와 1440/390/320 overflow/keyboard 확인을 수행한다. Go 문자열 CSS는 FullOps DESIGN 검사 미지원이므로 실제 렌더 검증을 후속 시각 담당이 할 수 있게 합성 fixture와 비밀 없는 검수 환경/방법을 제공한다.
소유권: internal/relay/member.go, 필요한 member_agents.go/http.go의 UI 및 관련 tests/fixtures, 본인 phase/contexts/packet outcomes, 필요한 기술 원천(D10/D09). 기존 adapters webhook 제품 코드는 고정 QA 후보이므로 수정하지 않는다. designer 디자인 정본과 인박스 수정 금지.
## 검증 담당과 후속
DEV 자체검증 후 고정 SHA에서 tester 독립 기능 QA·코드 리뷰, designer 직접 1440/390 시각 검수, 수락 뒤 main/OPS 배포다. 자체 캡처는 민감정보 없는 합성 데이터만. 전체 실계정 OAuth 왕복·제품 최종 수락을 자체 검사로 대신하지 않는다.
## 갱신할 산출물
D10/D09 영향 원천만 필요한 범위 갱신. D04는 read-only 정본. UI 토큰 정본은 service-member-ui.md 8절이며 실제 공유 CSS 위치와 검사/예외를 기록한다.
