---
title: SAR-SERVICE-UI-001-VISUAL — 회원 웹 고정 후보 직접 시각 수락 검수
status: draft
updated: 2026-10-11
owner: designer
tasks: [SAR-SERVICE-UI-001-VISUAL]
summary: 회원 웹 고정 후보 직접 시각 수락 검수
attempt: f1d96e3cc69e458f9ccd7521778f59f2
base: 8539ef5b6d3e7ed99259762c62d164b8ba40fab0
subagent_level: standard
test_level: lite
---

# SAR-SERVICE-UI-001-VISUAL — 회원 웹 직접 시각 수락 검수

- 상태: queued. DEV 완료 고정 SHA와 실제 구현자 세션을 coor가 착수 전 추가한다. 지금 dispatch하지 않는다.
- 역할/소유권: designer, C:/Users/shin/orca/workspaces/KnowsLink/designer, fullops/designer. 본인 시각 QA·PNG·phase·contexts·packet outcomes와 인박스/로그만. 제품 코드·DEV·tester 공간 수정 금지. 판정 뒤 디자인 정본의 실제 수락 상태는 갱신 가능하다.
- 복귀: run_7176a2ad6ef1, coor term_e61d3e14-29e9-4954-943a-4a75707c82de. 실제 task/dispatch는 preamble.
- 기준: fullops-common-0.3.3과 세 규칙, project.md, FULLOPS.md, 문서 작성 규칙. frontend-design 실제 경로 C:/Users/shin/.agents/skills/frontend-design/SKILL.md. Ponytail full, test lite/subagent standard. Codex hook 사용자 Trust 미승인으로 등록 Claude Sonnet5.5 high 새 세션 예외. Trust 자동 승인 금지.
- 제약: 추가 유료 서비스·로컬 Docker·삭제·이동·운영 변경 금지. DB/자격/agent/관계/공유 서비스/F: 증거 보존. 비밀값 출력 금지. 웹은 Orca 내장 브라우저만. Grok Bot 접근 금지. 실제 세션 로그아웃/agent 생성·폐기 금지.

## 제품 시각 수락

디자인 정본 docs/design-docs/mockups/service-member-ui.md의 IA·토큰·필수 R-01–09·수락 조건과 고정 후보를 직접 대조한다. DEV synthetic fixture와 Orca 내장 브라우저로 실제 login/home(4 agent, 빈 상태 가능)/연결·device·receipt/오류의1440·390 PNG를 보존하고320 overflow·키보드focus·44px 조작·contrast·한국어줄바꿈을 확인한다. 스크린샷 CLI JSON의 base64를 로그에 출력하지 말고 파일로 저장한다. 값은 합성 fixture만 사용하고 live 개인 정보는 캡처하지 않는다. 디자인 lint/테마 전환은 기존 지원 범위로 검증하거나 해당 없음 사유를 남긴다.

직접 픽셀 관찰·DOM 측정·worker 보고를 구분한다. 실제 Google 인증 왕복 성공으로 확대하지 않는다. 좁은 디자인 수정 권고는 담당 DEV에 coor 경유 전달하고 직접 제품 코드를 바꾸지 않는다. 미해결 critical/high는 수락 차단이며 낮은 등급도 남긴다. 기능 QA/기술 리뷰는 별도 tester 담당이다.

산출물 인덱스부터 대화 미참조 인계 점검, screenshot manifest·phase·contexts·packet outcomes·미검증을 기록한다. work.py finish로 전문 아카이브·인박스 비움 후 커밋한다. 완료 보고 첫 부분은 [완료] SAR-SERVICE-UI-001-VISUAL | 최종 SHA <실제 결과40자리>로 적는다. main/배포는 coor 후속이다.

## 지시 전제와 충돌 — 먼저 확인
Jev context는 service-member-ui.md를 conflict로 표시했으나 구체적 이유는 contradicts task만 반환했다. coor가 정본12절을 대조했다. R-01–09는 필수, R-10–11은 권장, R-12 다크모드·실접속·클라이언트복사는 후속으로 이번 범위 밖이다. 해당 없음인 후속 기능을 요구하지 않는다. 이 원천을 먼저 읽고 실제 불일치가 있으면 coor 경유 확인한다. 12.4의 실제 흐름 검수는 합성 fixture에서 수행하며 실계정 상태 변경은 하지 않는다. Google 인증 실제 왕복 미검증은 솔직히 남긴다. 이 해석은 기존 제품 범위를 줄이거나 넓히지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-SERVICE-UI-001-VISUAL-packet.json` / SHA `457b980413baf41d8c5bc661f906bcc2d71bffca` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.52, 'evidence': 0.77, 'contradicts': 0.15, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/agents/document-writing.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.24, 'evidence': 0.55, 'contradicts': 0.14, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 7, 213, 213 · inferred · 필수 · {'relevant': 0.89, 'evidence': 0.86, 'contradicts': 0.66, 'injection': 0.06, 'decision': 'conflict', 'reason': 'contradicts task'}
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 6, 10, 24, 33, 55, 67, 92, 104, 105 · inferred · 필수 · {'relevant': 0.86, 'evidence': 0.91, 'contradicts': 0.36, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 28, 28, 37, 37 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.73, 'evidence': 0.81, 'contradicts': 0.3, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.14, 'evidence': 0.32, 'contradicts': 0.18, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.46, 'evidence': 0.53, 'contradicts': 0.14, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.43, 'evidence': 0.56, 'contradicts': 0.13, 'injection': 0.1, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.55, 'evidence': 0.68, 'contradicts': 0.16, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `.fullops-squad/PLANS.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 55 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 6, 30, 32, 68, 70, 72 · inferred · {'relevant': 0.22, 'evidence': 0.34, 'contradicts': 0.35, 'injection': 0.12, 'decision': 'keep', 'reason': None}
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 28 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 119 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 141 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 7 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 18, 42 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 35 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md` (document_read) · 줄 21 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 6, 14, 31, 41 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 7, 116, 145, 145, 145, 221, 227, 228, 232, 234, 234 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/COOR/dev-fix-2-final/coordinator-handover-supplement.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-001-TESTER-test/report.md` (document_read) · 줄 60 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 16, 17 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW-review/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-REVIEW-review/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-SERVICE-001-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL/report.md` (document_read) · 줄 6, 10, 77 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-VISUAL/report.md` (document_read) · 줄 6, 10, 72 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-3.md` (document_read) · 줄 15 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (document_read) · 줄 6, 10, 26, 28, 32, 48, 49, 66 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPEN-PREP.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL.md` (document_read) · 줄 2, 6, 10, 21 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OPS.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_read) · 줄 7, 305, 307, 312, 336 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 7, 89 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 7 · inferred
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
- `.fullops-squad/handovers/logs/2026-10-11_to_designer.md` (document_read) · 줄 9, 10, 13, 17, 25, 49, 69, 77, 81, 143, 165, 166, 169, 173, 180, 210, 246, 257, 258, 268, 302, 304, 305, 308, 312, 320, 337, 343, 365, 377, 391, 392, 393, 400, 433 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_ops.md` (document_read) · 줄 124, 155, 186, 191, 204, 205, 210, 216 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-DEV-REVIEW.md` (document_read) · 줄 45 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-INTEGRATION-REVIEW.md` (document_read) · 줄 42 · inferred
- `adapters/README.md` (document_read) · 줄 전체/미확인 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 16 · inferred
- `adapters/src/login.test.ts` (impact_check) · 줄 88, 95 · inferred
- `adapters/src/mcp.ts` (impact_check) · 줄 78 · inferred
- `internal/relay/connections_test.go` (impact_check) · 줄 335, 338, 409 · inferred
- `internal/relay/device.go` (direct_edit) · 줄 11, 15 · inferred
- `internal/relay/http.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/member.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/member_agents.go` (direct_edit) · 줄 61, 65, 68, 94, 119, 132, 140 · inferred
- `internal/relay/member_pages.go` (direct_edit, impact_check) · 줄 60 · inferred
- `internal/relay/member_receipt.go` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/member_ui_test.go` (direct_edit, impact_check) · 줄 21, 96, 120, 161 · inferred
- `internal/relay/policy_test.go` (impact_check) · 줄 22, 31, 59, 65, 68, 68 · inferred
- `scripts/install_bot_mcp.sh` (impact_check) · 줄 83 · inferred
미확인 11건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고
작업 완료 뒤 실제 결과 SHA·후보 SHA·검증 근거·미검증·후속을 이 절에 전문으로 기록한다. 지금은 DEV 고정 SHA를 기다린다.

## 고정 후보와 구현자 — 착수 정본
고정 후보 c7828b227b5715f7bdf02bff46bf9cf102a758d9. 실제 구현자 Claude 세션5953ccd3-7f4f-4c5e-b63d-c02d36ee35ac, dispatch ctx_d554b501df04. DEV clean HEAD 일치 확인. 구현자 phase docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md의 시각 검수 환경 절을 먼저 읽는다. KNOWSLINK_UI_FIXTURE_DIR=<새 빈 폴더> go test ./internal/relay -run TestMemberUIFixtures -count=1로 29개 합성 HTML을 만든다. snapshot 자체에는 생성물을 쓰지 말고 고유 실행 clone 또는 별도 기록 역할의 동일 제품 바이트에서 생성한다. 기준 main8539ef5와 제품 diff를 확인한다. DEV는 agent 서버 안내문·짧은 날짜·200% 확대·실스크린리더 미적용/미검증을 기록했다. 수락 기준과 대조하고 차이의 영향과 후속을 남긴다. form golden의 추가 헤더 logout 1건은 명세4.2의 의도적 추가이며 실제 서버 보안 검사를 보존한다.
