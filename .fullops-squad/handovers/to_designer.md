---
title: SAR-SERVICE-UI-001-DESIGN — 회원 웹 상용 디자인과 개발 인계
status: draft
updated: 2026-10-11
owner: designer
tasks: [SAR-SERVICE-UI-001-DESIGN]
summary: 회원 웹 상용 디자인과 개발 인계
attempt: d24a0a05b3984fddbd62f78f24bedd07
base: 7a8c2b9194edb77b52db41a14b2325af743f0d6e
subagent_level: standard
test_level: lite
---

# SAR-SERVICE-UI-001-DESIGN — 회원 웹 상용 디자인과 개발 인계

- 작성일: 2026-10-11
- From / To: coor / designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 담당 repo id: 0b08ec4c-9e3a-4613-8197-5a835e545335
- 담당 워크트리: C:/Users/shin/orca/workspaces/KnowsLink/designer; 브랜치: fullops/designer
- 병합 책임자: coor; 기본 브랜치: main
- 복귀: C:/Users/shin/orca/workspaces/KnowsLink/coor, term_e61d3e14-29e9-4954-943a-4a75707c82de. 실제 run/task/dispatch는 Orca preamble을 사용한다.

## 적용 기준과 예외
fullops-common-0.3.3: rules/common/README.md와 coding-style.md, testing.md, security.md를 읽는다. project.md와 FULLOPS.md, contexts/designer.md, docs/agents/document-writing.md를 읽는다. 기준 SHA는 front matter의 base다. Ponytail full 적용. 실제 사용 중인 기존 에이전트 세션은 종료하지 않는다.
추가 과금 금지. 로컬 Docker 금지. 검증은 기존 Linux 서버에서 고유 후보 체크아웃으로 실행한다. 기존 DB/키/agent/pair/Google/Tunnel 자격/공유 서비스를 보존한다. 파일·컨테이너·임시 증거 삭제/이동/정리 금지. F:/KnowsLink-Incident-20261011-110130 보존. 비밀값은 로그/Git/채팅에 쓰지 않는다.
웹은 Orca 내장 브라우저만 사용한다. Computer Use는 Grok Bot 설정 제어에만 사용한다. Grok UI를 열어 수신을 유발하거나 수신 결과를 읽지 않는다. 제품 배포·Bot 설정 변경은 후속 검수 이후 coor/ops 담당이다.

## 완료 기준과 검증
현재 범위 완료 후 docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md에 기록하고, packet-outcomes의 모든 필수 path/category와 partial 불확실성을 처리한다. work.py finish로 전문 보존·빈 인박스를 확인한다. 커밋 후 fullops lint를 base 기준으로 실행하고 SHA/ERROR/WARNING/불가 사유를 기록한다. worker_done은 actual preamble로 보낸다. 구현과 제품 최종 수락을 구분한다.
멈출 조건: 파괴 작업, 범위 밖 파일 변경, 설명되지 않는 필수 검증 실패, 제품 규칙 변경만 ask/escalation한다. 비파괴 작업은 중간 승인을 기다리지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-SERVICE-UI-001-DESIGN-packet.json` / SHA `7a8c2b9194edb77b52db41a14b2325af743f0d6e` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 30, 32 · inferred · 필수 · {'relevant': 0.37, 'evidence': 0.56, 'contradicts': 0.29, 'injection': 0.1, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/agents/document-writing.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.66, 'evidence': 0.8, 'contradicts': 0.1, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 31, 31 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 74 · inferred · 필수 · {'relevant': 0.71, 'evidence': 0.86, 'contradicts': 0.31, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.38, 'evidence': 0.59, 'contradicts': 0.2, 'injection': 0.06, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 17 · inferred · 필수 · {'relevant': 0.66, 'evidence': 0.72, 'contradicts': 0.14, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.67, 'evidence': 0.8, 'contradicts': 0.1, 'injection': 0.08, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.6, 'evidence': 0.76, 'contradicts': 0.15, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `internal/relay/member.go` (impact_check) · 줄 전체/미확인 · unknown · 필수 · {'relevant': 0.47, 'evidence': 0.49, 'contradicts': 0.27, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `internal/relay/member_agents.go` (impact_check) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.23, 'evidence': 0.44, 'contradicts': 0.39, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `.fullops-squad/PLANS.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 52 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 119 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 141 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 81 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 18, 78 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 35, 111 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md` (document_read) · 줄 21, 92 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 6, 14, 31, 41 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/evaluations/qa-reports/COOR/dev-fix-2-final/coordinator-handover-supplement.md` (document_read) · 줄 20, 33 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-REVIEW-review/report.md` (document_read) · 줄 21 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 16, 17, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW-review/report.md` (document_read) · 줄 16, 101 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-REVIEW-review/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-2-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 18, 23, 89 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-RECORD-REVIEW-review/report.md` (document_read) · 줄 33 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-SERVICE-001-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.14.md` (document_read) · 줄 26, 27 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PREP-002.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md` (document_read) · 줄 63 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX.md` (document_read) · 줄 76 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 35, 80 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2.md` (document_read) · 줄 14, 86 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-3.md` (document_read) · 줄 15, 78 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX.md` (document_read) · 줄 14, 69 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (document_read) · 줄 6, 10, 26, 28, 32, 48, 49, 66 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPEN-PREP.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV.md` (document_read) · 줄 39 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_read) · 줄 7, 305, 307, 312, 336 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 7, 89 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 6, 12, 14, 20, 23, 24, 27 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 7, 40, 56, 66 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 7, 25, 83, 154 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 7, 15, 23, 63, 76, 125, 128, 152 · inferred
- `.fullops-squad/handovers/_TEMPLATE.md` (document_read) · 줄 53 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_designer.md` (document_read) · 줄 9, 16, 20, 79, 87, 114 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_dev.md` (document_read) · 줄 16, 33, 35, 37, 51, 93, 125, 181 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_ops.md` (document_read) · 줄 90, 93, 97, 101, 116, 130, 144, 147, 151, 155, 171 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_tester.md` (document_read) · 줄 16, 33, 35, 37, 49 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 35, 74, 102, 105, 106, 142, 248, 249, 306, 386 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 37, 176, 206, 245, 246, 304, 357, 364 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_tester.md` (document_read) · 줄 100, 175, 176, 228 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_designer.md` (document_read) · 줄 29, 46, 69, 76, 133, 182, 183, 184 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 24, 32, 64, 108, 113, 148, 149, 229, 230 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_tester.md` (document_read) · 줄 48, 70, 75, 139, 180, 181, 182 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_designer.md` (document_read) · 줄 257, 258, 268 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_ops.md` (document_read) · 줄 204, 205, 216 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_tester.md` (document_read) · 줄 254, 361, 362, 363, 375 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-DEV-REVIEW.md` (document_read) · 줄 45 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-INTEGRATION-REVIEW.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/lint/README.md` (document_read) · 줄 31, 81, 83, 94, 104 · inferred
- `.fullops-squad/review/_REPORT.md` (document_read) · 줄 20 · inferred
- `adapters/.grok-plugin/plugin.json` (impact_check) · 줄 전체/미확인 · inferred
- `adapters/README.md` (document_read) · 줄 전체/미확인 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 전체/미확인 · inferred
- `adapters/src/trial-cli.ts` (impact_check) · 줄 전체/미확인 · inferred
- `cmd/relay/main.go` (impact_check) · 줄 전체/미확인 · inferred
- `internal/config/config.go` (impact_check) · 줄 전체/미확인 · inferred
- `scripts/check_public_ingress.py` (impact_check) · 줄 전체/미확인 · inferred
- `scripts/shell_guard.py` (impact_check) · 줄 전체/미확인 · inferred
미확인 9건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고
완료 시 변경 이유/실제 SHA/검증 명령과 종료코드/미검증/후속/산출물·로그를 전문으로 기록한다.

## 현재 상황과 확인 근거
사용자는 link.knowslog.com에서 Google 로그인한 이메일에 등록된 에이전트를 보는 화면을 상용 서비스 수준으로 개선하라고 요청했다. Go SSR internal/relay/member.go의 inline template+memberStyle이며 별도 프레임워크는 없다. 기존 기능은 실제 운영 중이다.

## 해야 할 일과 파일 소유권
- [ ] frontend-design 스킬을 읽고 현재 공개 로그인 및 실제 로그인 회원 화면을 Orca 내장 브라우저로 확인한다. 계정/키/개인정보는 캡처와 로그에서 최소화·가린다.
- [ ] 회원 로그인→내 에이전트→연결 관리 흐름의 제품 방향·정보 구조·한국어 문구·색상/타입/여백/반응형/빈 상태·오류 상태·접근성 수락 기준을 docs/design-docs/mockups/service-member-ui.md에 작성한다. D04 원천을 확인해 갱신한다.
- [ ] 구현 방법·파일별 기술 계획은 DEV가 결정하도록 제품 요구를 인계한다. dev 인박스는 webhook 구현이 사용하므로 덮어쓰지 말고 이 디자인 문서의 개발 인계 절에 범위·완료 기준을 쓴다.
파일 소유권: 지정 디자인 문서, D04 원천, 본인 phase/contexts/로그/packet outcomes. 제품 코드 수정 금지. PLANS/board는 coor 소유다.

## UI 작업의 디자인 기준과 검증
Go SSR 재사용. 추가 CSS/JS 라이브러리·외부 폰트·유료 에셋 없이 가능하도록 기준을 정한다. 무작위 그라디언트, 모든 정보의 동일한 카드 배치, 지나친 알약 버튼, 장식용 통계, 긴 한국어의 무조건 anywhere 줄바꿈은 피한다. 실제 상태 없는 연결됨/보안검증 표시를 만들지 않는다.
기존 CSRF·권한/owner 분리·재인증·파괴 동작 확인과 CSP를 보존한다. 최근 로그인한 개인 계정과 노우 계정은 서로 다른 owner다. 다른 owner 에이전트 노출 금지.
캡처는 현재 로그인/회원 화면 기준 1440px/390px의 정보 위계·overflow·기능 발견성을 판정하기 위한 범위로 한정한다. 구현 후 직접 시각 검수는 후속 지정 fixed SHA에서 수행한다. 이번 결과는 디자인 인계 문서이며 최종 UI 수락이 아니다.

## 갱신할 산출물
D04. 문서 디자인 lint와 front matter 확인. 제품 코드 변경이 없어 make lint/test 반복은 해당 없음으로 명시한다.
