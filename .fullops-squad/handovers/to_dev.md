---
title: SAR-FRONTEND-POLISH-001-DEV — 기존 회원 UI 검토와 결함 수정으로 완성도 개선
status: draft
updated: 2026-10-12
owner: dev
tasks: [SAR-FRONTEND-POLISH-001-DEV]
summary: 기존 회원 UI 검토와 결함 수정으로 완성도 개선
attempt: e5b3085323c84ed69a5d155b9a006ed6
base: 5c6485b4c7c1c573a14bf3e683701b2399284666
subagent_level: standard
test_level: lite
---

# SAR-FRONTEND-POLISH-001-DEV — 기존 회원 UI 검토와 결함 수정으로 완성도 개선

- From / To: coor / dev. 상태: ready. Purpose: implementation.
- 테스트 lite / 하위 에이전트 standard. 코드 소유: internal/relay/ 회원 UI와 관련 회귀 검사. 기존 파일 삭제·이동 정리 금지.
- 작업 경로: C:/Users/shin/orca/workspaces/KnowsLink/dev, fullops/dev.
- 복귀: run_2d134e8802eb, coor C:/Users/shin/orca/workspaces/KnowsLink/coor. task/dispatch는 Orca preamble 정본.
- 구현 tuple: claude / claude-opus-5-5 / medium. fresh: 새 독립 과제이며 직전 완료 세션을 재사용하지 않는다.

## 목표·범위·규칙

사용자 요청: “프론트앤드 기능을 업데이트 해줬어 지금 프론트앤드 검토해서 완성도를 높여봐”.
기존 로그인·회원 홈·연결 관리·승인 관련 화면을 검토하고 실제 결함을 최소 수정한다. 반응형, 키보드/포커스, 접근성, 오류·빈 상태, 동작 안정성을 확인한다.
기존 브랜드·제품 문구의 사실·보안·권한·관계·rate·키 한도를 유지한다. 새 기능·프레임워크·테마·유료 도구·외부 폰트는 범위 밖이다.
제품 방향·공유 규칙 변경만 coor에게 질문한다. 기술 분석·계획·구현·직접 검증은 같은 과제에서 끝낸다.
운영 서비스의 실제 계정/키/관계/메시지 변경·로그아웃·배포는 금지한다. 로컬 합성 fixture와 정확한 localhost origin으로 검증한다. 비밀값을 출력하거나 캡처하지 않는다.
기존 hook 수정 후보와 FullOps 업데이트 수락 보류는 별개다. 관련 코드를 수정하거나 통과로 재해석하지 않는다.

## 적용 기준·먼저 읽을 문서

fullops-common-0.3.3: .fullops-squad/rules/common/README.md 및 coding-style.md, testing.md, security.md.
.fullops-squad/FULLOPS.md, project.md, contexts/dev.md, docs/agents/frontend.json.
디자인 정본: .fullops-squad/docs/design-docs/mockups/service-member-ui.md. 기존 Go template CSS·토큰·공용 패턴을 재사용한다.
fullops-frontend와 fullops-work, fullops-test를 사용한다. 주 디자인 스킬은 impeccable 하나다. 기존 화면 polish/audit로 적용한다.
라이브러리·API·CLI 사용 근거가 필요하면 Context7 resolve-library-id → query-docs를 사용한다. 미지원이면 공식 문서를 사용하고 출처를 기록한다.

## 해야 할 일·완료 기준

- [ ] 대상 UI 구현·호출자·기존 검사·최근 변경과 정본을 확인하고 짧은 기술 계획 및 결함 목록을 기록한다.
- [ ] frontend.py --host claude-code --model claude-opus-5-5 --effort medium --check로 준비를 확인한다. 실제 인증/브라우저 성공과 구분한다.
- [ ] impeccable context를 한 번 실행한다. 실패 시 스킬의 fallback을 따른다. UI edit 직전에 craft-floor를 읽는다.
- [ ] desktop 1440×900, mobile 390×844에서 핵심 화면·키보드·포커스·가로 overflow·오류/빈 상태를 합성 fixture로 확인한다. 실제 결함을 한 번에 수정한다.
- [ ] 변경 핵심 동작과 필요한 실패 경계에 짧은 검사. 프로젝트 필수 product-lint/product-test를 유지한다. non-trivial 로직은 작은 회귀 검사로 보호한다.
- [ ] 변경을 커밋하고 지시서 base에서 현재 HEAD의 FullOps lint ERROR0 및 등록 product-lint/product-test exit0 증거를 남긴다. 파이프로 종료코드를 가리지 않는다.
- [ ] 실행 보고/비평/캡처(합성 자료만)/명령과 종료코드를 qa-reports/SAR-FRONTEND-POLISH-001-DEV-test/에 기록한다. 테마 전환·별도 디자인 lint가 없으면 해당 없음의 이유를 적는다.
- [ ] packet-outcomes를 완성하고 work.py finish로 인박스와 완료 보고를 보존한다. 최종 커밋 SHA·변경·검증·미검사·독립 QA/시각 검수 필요를 worker_done으로 회신한다.

## 후속·산출물

갱신할 D01–D13: 없음. 실제 정본에 영향이 있으면 이유와 해당 문서만 갱신한다.
DEV 완료 뒤 coor가 tester 독립 동작 QA, designer 직접 시각 검수, 고정 SHA 독립 리뷰를 배정한다. 미해결 critical/high와 필수 실패는 병합 차단이다.
수정 필요 없음도 실제 화면/동작 근거가 있어야 한다. 캡처·비평은 최초와 수정 확인 최대 2회로 제한한다. 불필요한 추가 polish를 반복하지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-FRONTEND-POLISH-001-DEV-packet.json` / SHA `dcfc9ad66dc2ce2c2f0d7144a7f4ae9128432781` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 213, 213 · inferred · 필수 · {'relevant': 0.9, 'evidence': 0.92, 'contradicts': 0.19, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 47, 47 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.73, 'evidence': 0.88, 'contradicts': 0.27, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 55 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 145, 145, 145, 227, 228, 232, 234 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-DESIGN-REVIEW-review/report.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-TESTER-review/report.md` (document_read) · 줄 32 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL/report.md` (document_read) · 줄 31, 85 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-review/report.md` (document_read) · 줄 36, 54 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-test/report.md` (document_read) · 줄 26, 80 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md` (document_read) · 줄 34 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 24, 33, 55, 67, 92 · inferred · {'relevant': 0.55, 'evidence': 0.71, 'contradicts': 0.61, 'injection': 0.05, 'decision': 'conflict', 'reason': 'contradicts task'}
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-DEV.md` (document_read) · 줄 23, 41 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-VISUAL.md` (document_read) · 줄 19 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-TESTER.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-VISUAL.md` (document_read) · 줄 21 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/protocol.md` (document_read) · 줄 209 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 261 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 16 · inferred
- `adapters/src/login.test.ts` (impact_check) · 줄 88, 95 · inferred
- `adapters/src/mcp.ts` (impact_check) · 줄 78 · inferred
- `internal/relay/connections_test.go` (impact_check) · 줄 335, 338, 409 · inferred
- `internal/relay/device.go` (impact_check) · 줄 55, 111, 117 · inferred
- `internal/relay/member_agents.go` (direct_edit, impact_check) · 줄 61, 65, 68, 94, 119, 132, 140 · inferred · {'relevant': 0.25, 'evidence': 0.34, 'contradicts': 0.2, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `internal/relay/member_pages.go` (direct_edit, impact_check) · 줄 60 · inferred · {'relevant': 0.63, 'evidence': 0.68, 'contradicts': 0.17, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `internal/relay/member_ui_test.go` (direct_edit, impact_check) · 줄 21, 96, 120, 161 · inferred · {'relevant': 0.61, 'evidence': 0.7, 'contradicts': 0.2, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `internal/relay/policy_test.go` (impact_check) · 줄 22, 31, 59, 65, 68, 68 · inferred
- `scripts/install_bot_mcp.sh` (impact_check) · 줄 83 · inferred
미확인 11건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

작업 완료 때 전문을 기록한다.

## 지시 전제와 충돌 — 먼저 확인

- .fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md: 과거 구현 기록의 frontend-design 선택과 옛 검사 명령은 당시 근거다. 이번 작업은 impeccable polish를 사용한다. 제품 규칙 변경이 아니며 기존 결과를 이번 SHA의 검증으로 주장하지 않는다. 과거 컨테이너·증거의 삭제 금지는 유지한다.

