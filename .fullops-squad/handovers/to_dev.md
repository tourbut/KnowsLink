---
title: SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV — Dots OAuth 승인 후 안전한 callback 복귀와 실패 경계 최소 수정
status: draft
updated: 2026-10-11
owner: dev
tasks: [SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV]
summary: Dots OAuth 승인 후 안전한 callback 복귀와 실패 경계 최소 수정
attempt: 0a90d721ecf04034a0bdda333f24c7a9
base: ae5d9f0db878216ee011610430615067488fbe4a
subagent_level: standard
test_level: lite
---

# SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV — Dots OAuth 승인 후 안전한 callback 복귀

- Purpose: implementation. ready. Test lite / Subagent standard.
- 역할 dev, C:/Users/shin/orca/workspaces/KnowsLink/dev, fullops/dev.
- coor 복귀 term_e61d3e14-29e9-4954-943a-4a75707c82de / run_dc9b4e6dd18a. 실제 task/dispatch는 preamble.
- 기준 ae5d9f0db878216ee011610430615067488fbe4a. fresh Claude Opus5.5 high(Jev).
- 사용자가 기존 Dots Google 인증 연결과 오류 해결을 요청했다. 기존 Linux 서버의 승인된 격리 테스트를 사용할 수 있다. 배포는 OPS 후속이다.

## 관찰과 목표

사용자는 실제 Dots custom MCP KnowsLink를 만들고 Google 신원 확인·설치 승인을 진행했다. 최종 remote consent에서 authorization_failed를 관찰했다. 승인 직후 서버 읽기 전용 집계: install 1, token 0, subscription 0, pending attempt 0. install 생성 시각은 2026-10-11T06:22:26Z, 06:24 조회에는 code 0이다. 기존 코드의 code 만료1분/attempt10분과 대조한다. 미완료 설치를 삭제하거나 DB 원문·키·token을 출력하지 않는다. 실제 app ID는 asdk_app_6acb0615e7f88191b07413e8595308ed, URL https://link.knowslog.com/mcp.

읽기 전용 별도 조사: consent GET의 공통 CSP form-action self+Google와 consent POST의 등록 callback303이 Chromium에서 충돌할 가능성이 높다. 승인 시 install+code 생성 후 attempt삭제가 먼저 커밋되어 callback이 차단되면 새로고침에 authorization_failed가 될 수 있다. 실제 브라우저 console capture가 꺼져 있어 확정하지 않았다. 단순 Go303검사는 이 경계를 검증하지 않는다. MDN 공식 form-action 문서(https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy/form-action)는 POST 이후 redirect 차단의 브라우저 차이를 명시한다. DEV가 기술 원인을 재현하고 최소 안전 수정방법을 정한다.

기존 제품 요구 내 목표: 명시적 승인 이후 등록된 호스트로 복귀하고 PKCE token 교환이 완료돼야 한다. 기존 Google 신원 재확인·client/resource/redirect/issuer binding·CSRF·TTL·키/관계 분리와 권한 검사를 유지한다. 전역 form-action wildcard, 임의 callback origin 문자열 결합, SameSite/CSRF 완화, 로그인 생략을 하지 않는다. 만료/이미 사용된 요청은 비밀 없는 사용자 안내와 새 연결 시작 방법을 제공한다. 시간 연장이나 기존 요청 재사용으로 보안을 약화하지 않는다.

## 소유권과 완료 기준

- internal/relay의 원인·관련 테스트, 필요한 D05/D10/D12 및 자신의 exec plan/QA/contexts/dev. 범위를 변경할 때 이유를 기록한다. 설치·callback 경계만 최소 수정한다. 공통 security middleware를 건드리면 영향을 받는 Google/remote/SSE 회귀 범위를 검증한다.
- 현재 요구·호출자·테스트를 먼저 확인하고 짧은 기술 계획을 같은 과제에 기록한다.
- Chromium 실제 브라우저로 합성 consent 승인 POST에서 외부 callback 도착까지 재현/수정 후 확인한다. synthetic provider/host와 격리 DB 사용. 실제 사용자 Google·Dots계정·기존 운영DB/키·메시지 본문은 사용하지 않는다. 로그에 callbackcode/credential을 남기지 않는다.
- 만료/재사용 실패 경계, 허용 callback만 복귀, 기존 Google와 OAuth token 교환 관련 회귀를 확인한다. 필수 registered lint/test를 exact final SHA에서 통과시키고 결과·HEAD·attempt·종료코드를 남긴다. local Docker는 금지이며 기존 Linux서버의 승인된 격리 테스트 경로를 재사용한다. 새 유료 서비스·새 dependency 없이 구현한다.
- fullops-common-0.3.3과 project.md, rules/delegation.md, document-writing.md 적용. Context7의 quota 실패 공식 대체근거는 기존 REMOTE/INSTALL계획을 재사용하되 새로운 API가 필요하면 규약대로 조회한다.
- 패킷의 partial/unknown/비선택 path/category를 완료 전에 outcomes로 처리한다. work.py finish 뒤 명확한 최종40SHA worker_done.
- 배포/다른서비스/복구/운영데이터 삭제·이동/force reset/clean/cleanup은 금지한다. 사용자 detachedcleanup 승인은 부모가 이미 완료했고 이 과제의 다른 파일 삭제 권한이 아니다. scratch·backup은 보존한다.

## 먼저 읽을 원천

FULLOPS.md, rules/common README와3규칙, project.md, orca-agents.md, rules/delegation.md, docs/agents/document-writing.md. internal/relay/remote_auth.go, google.go, http.go, capacity.go, remote_integration_test.go, remote_test.go, google_test.go. docs/exec-plans/phases/SAR-SSE-INSTALL-001-REMOTE.md와DEV.md, evaluations/qa-reports/SAR-SSE-INSTALL-001-OPS-LIVE.md, operations/ops-guide.md. 기존 원검수 medium SSE헤더와 verify.py baseline 실패를 이번 신규결함과 혼동하지 않는다.

## 후속

작성자와 다른 fixedSHA 독립 리뷰·tester QA 뒤 main 통합/push·OPS 최소배포·실제 사용자 Google/Dots설치는 부모가 조율한다. DEV 완료가 전체 실계정 연결 수락을 대신하지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV-packet.json` / SHA `ae5d9f0db878216ee011610430615067488fbe4a` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OPS-LIVE.md` (document_read) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-REMOTE.md` (document_read) · 줄 2, 2, 6, 10, 16, 22, 22, 27, 30, 31, 32, 44, 45, 48, 51, 57, 58, 62, 68, 72, 72, 72, 76 · inferred · 필수 · {'relevant': 0.42, 'evidence': 0.75, 'contradicts': 0.47, 'injection': 0.06, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 2, 2, 6, 6, 6, 7, 7, 7, 14, 14, 14, 14, 20, 24, 28, 32, 34, 35, 42, 54, 54 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 6 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 70, 76 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 28, 29 · inferred
- `.fullops-squad/contexts/ops.md` (document_read) · 줄 6, 34, 67 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 7, 14, 20, 23, 25, 155, 163, 179 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read) · 줄 7, 14, 38 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 7, 14, 19, 21 · inferred
- `.fullops-squad/docs/design-docs/database-design.md` (document_read) · 줄 7, 14, 39, 41 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 14, 23, 28, 190 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 87 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 82 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read, document_update) · 줄 7, 14, 24, 55, 57, 188, 189, 191, 197, 199 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 7, 14, 16, 16, 19, 21, 21 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-001-TESTER-test/report.md` (document_read) · 줄 43, 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-DEV-review/report.md` (document_read) · 줄 37, 40 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-DEV/report.md` (document_read) · 줄 34 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-LOGIN-001-REVIEW-review/report.md` (document_read) · 줄 20, 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-TESTER.md` (document_read) · 줄 2, 6, 10, 23, 27, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 130 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-review/report.md` (document_read) · 줄 2, 6, 10, 17, 51 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 93 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OPS-REVIEW-review/report.md` (document_read) · 줄 6, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-REVIEW-review/report.md` (document_read) · 줄 2, 6, 10, 15, 18, 46, 46 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-TESTER-test/report.md` (document_read) · 줄 2, 6, 10, 27, 27 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-VISUAL/report.md` (document_read) · 줄 6, 10, 72, 84, 84 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.2.0.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 99 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-001-DEV.md` (document_read) · 줄 24, 34 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-002-DEV.md` (document_read) · 줄 17, 18 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md` (document_read) · 줄 18, 27, 33, 33, 35, 39, 39 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md` (document_read) · 줄 18, 61 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md` (document_read) · 줄 2, 6, 10, 17, 55, 56, 57, 59, 60, 61 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md` (document_read) · 줄 68, 72, 93 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 58 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (document_read) · 줄 36, 56, 56 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 24, 53, 204, 221 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-DEV.md` (document_read) · 줄 6, 7, 10, 16, 34, 36, 38, 48, 55, 76, 80, 82, 82, 82, 82, 86, 86, 92, 97 · inferred · {'relevant': 0.37, 'evidence': 0.62, 'contradicts': 0.71, 'injection': 0.08, 'decision': 'conflict', 'reason': 'contradicts task'}
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-INSTALL.md` (document_read) · 줄 2, 6, 10, 29, 30, 30, 33, 54, 55, 55, 57, 97, 102 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OPS.md` (document_read) · 줄 2, 6, 10, 21, 22, 24, 37, 37 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_read, document_update) · 줄 7, 14, 16, 20, 22, 22, 24, 28, 40, 42, 42, 45, 49, 51, 51, 278, 312, 338, 340, 346, 356, 368 · inferred · {'relevant': 0.52, 'evidence': 0.58, 'contradicts': 0.68, 'injection': 0.17, 'decision': 'conflict', 'reason': 'contradicts task'}
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 7, 13, 15, 15, 17, 17, 19, 21, 23, 23, 28, 30, 51, 93, 95, 103 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 7, 13, 25, 27 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 50 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 136, 138 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 35 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/README.md` (document_read) · 줄 44 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md` (document_read) · 줄 134 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/product.md` (document_read) · 줄 122 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_ops.md` (document_read) · 줄 191 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_dev.md` (document_read) · 줄 9, 12, 16, 20, 33, 39, 71, 104, 126 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_ops.md` (document_read) · 줄 9, 12, 16, 20, 29, 34, 38, 41, 52, 54 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_tester.md` (document_read) · 줄 9, 12, 16, 20, 29, 34, 43, 51, 57 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_dev.md` (document_read) · 줄 60, 89 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 76 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 210 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_designer.md` (document_read) · 줄 113, 157 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 115, 179, 210, 299, 307, 307, 309, 316 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_tester.md` (document_read) · 줄 116, 158 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_designer.md` (document_read) · 줄 9, 10, 13, 13, 17, 18, 25, 25, 34, 34, 40, 45, 47, 49, 67, 69, 77, 81, 143, 155, 155 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_dev.md` (document_read) · 줄 9, 10, 13, 17, 25, 36, 71, 103, 105, 110, 114, 116, 118, 120, 120, 288, 336, 363, 364, 365 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_ops.md` (document_read) · 줄 9, 10, 14, 18, 26, 35, 54, 54, 61, 80, 83, 85, 85, 87, 91, 93, 93 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_tester.md` (document_read) · 줄 9, 10, 13, 17, 25, 36, 36, 36, 46, 48, 55, 55, 65, 71, 83, 83, 97, 101, 103 · inferred
- `.fullops-squad/rules/delegation.md` (document_read) · 줄 전체/미확인 · inferred
- `adapters/README.md` (document_read) · 줄 6, 14, 24, 24, 50, 50, 52, 56, 56 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 14, 14 · inferred
- `adapters/src/login.test.ts` (impact_check) · 줄 전체/미확인 · inferred
- `cmd/relay/main.go` (impact_check) · 줄 전체/미확인 · inferred
- `deploy/knowslink/tunnel/public-ingress.yml` (impact_check) · 줄 5 · inferred
- `internal/config/config.go` (impact_check) · 줄 전체/미확인 · inferred
- `internal/relay/admission_integration_test.go` (impact_check) · 줄 전체/미확인 · inferred
- `internal/relay/capacity.go` (impact_check) · 줄 전체/미확인 · inferred · {'relevant': 0.16, 'evidence': 0.32, 'contradicts': 0.18, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `internal/relay/connections_test.go` (impact_check) · 줄 전체/미확인 · inferred
- `internal/relay/device.go` (impact_check) · 줄 전체/미확인 · inferred · {'relevant': 0.08, 'evidence': 0.14, 'contradicts': 0.28, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `internal/relay/google.go` (impact_check) · 줄 전체/미확인 · unknown
- `internal/relay/google_test.go` (impact_check) · 줄 165 · inferred · {'relevant': 0.27, 'evidence': 0.42, 'contradicts': 0.18, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `internal/relay/http.go` (direct_edit, impact_check) · 줄 전체/미확인 · unknown · {'relevant': 0.18, 'evidence': 0.24, 'contradicts': 0.19, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `internal/relay/remote_auth.go` (direct_edit, impact_check) · 줄 전체/미확인 · unknown · {'relevant': 0.44, 'evidence': 0.5, 'contradicts': 0.17, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `internal/relay/remote_integration_test.go` (direct_edit, impact_check) · 줄 전체/미확인 · unknown
- `internal/relay/remote_test.go` (impact_check) · 줄 전체/미확인 · unknown
- `scripts/check_public_ingress.py` (impact_check) · 줄 9 · inferred
- `scripts/package_plugin.py` (impact_check) · 줄 144, 146 · inferred
미확인 15건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

작업 완료 뒤 전문을 작성하고 work.py finish로 보존한다.

## 지시 전제와 충돌 — 먼저 확인

DEV 기록은 실제 Dots 미검증, ops-guide는 과거 미적용/held와 최신 운영 적용을 함께 보존한다. 현재 수락된 제품은121cde9이고 실제 운영 적용은OPS-LIVE가 정본이다. 이번 관찰은 이후 실계정 연결의 실패이며 예전 unit/HTTP성공이나 배포 성공이 브라우저 callback 성공을 보장하지 않는다. 원문 전체를 읽고 현재 설치 결과와 과거 미검증을 분리한다. shared HTTP header 변경 시 기존Google/SSE경계를 다시 확인한다.
