---
title: SAR-GOOGLE-CONNECT-002-DEV — 추가 과금 없이 기존 Tunnel·도메인·플러그인·Google 인증 연결을 최소 보완한다
status: draft
updated: 2026-10-10
owner: dev
tasks: [SAR-GOOGLE-CONNECT-002-DEV]
summary: 추가 과금 없이 기존 Tunnel·도메인·플러그인·Google 인증 연결을 최소 보완한다
attempt: 185beb2fe9c8422d922d88bfe45093ae
base: 5af28c9fcfc289a845f9730aed34d32308e8b5a6
subagent_level: off
test_level: lite
---

# SAR-GOOGLE-CONNECT-002-DEV — 비용 없는 Tunnel·Google 연결 최소 보완

- Purpose: implementation. From coor / To dev. Test level lite; Subagent level off.
- 기준: 5af28c9. 워크트리 C:/Users/shin/orca/workspaces/KnowsLink/dev, 브랜치 fullops/dev.
- 복귀: run_86e0e674b5a0 / term_e61d3e14-29e9-4954-943a-4a75707c82de. 실제 task/dispatch는 새 preamble 정본.
- 사용자 확정: 기존 Linux 서버를 Cloudflare Tunnel로 인터넷에 노출하고 자기 도메인에 연결한다. 봇에 KnowsLink 플러그인을 설치한 뒤 Google 로그인으로 등록하고 두 독립 agent가 대화한다. 최소 구현. 공개 고정 IP 구매가 아니라 Tunnel HTTPS 도메인으로 목적 충족 여부를 기술적으로 확인한다.
- 비용 금지: 유료 요금제, 무료 초과 자동 과금 동의, 새 유료 서비스·구독 모두 제외. Cloudflare Access 활성화/결제 방식은 사용자가 거부했다. 이를 필수 선행으로 두지 않는다.
- 기존 구현: DEV152217f, QA44e9e457은 main/origin에 통합됐다. 운영 서버 relay04a65b2 배포·backup·health·localverify 통과. 기존 Tunnel fallback Access required:true, root Access 앱이 모든 Google/home 경로를 막는다. 현재 Cloudflare dashboard는 개요에 앱2개가 보이지만 앱 관리에서 활성요금제 요구. 원래 owner-only 앱·trial 앱의 제거/범위변경과 Tunnel ingress 조정 등 비용 없는 최소 대안을 실제 근거로 제시한다. coor가 로그인된 Orca 브라우저와 운영 서버 적용을 맡는다.
- Google client/secret/callback은 운영 .env에 이미 설정. 기존 callback https://link.knowslog.com/auth/google/callback 재사용 우선. 기존 회원·DB·키·OAuth 자격·공유 서비스 보존.
- 노우 Bot: 사용자 PC Grok Bot 앱에 설치된 MCP는 새 checkout /workspace/KnowsLink-it-20261010 04a65b2, public-node, RELAY_URL https://link.knowslog.com, folder /workspace/.knowslink-connect/nou/google-it-20261010. 실제 도구9개(connect/status포함) 준비됐고 로그인 아직 시작 전. Bot 조작은 coor의 computer-use만, 웹은 Orca브라우저만. Google 로그인은 사용자가 직접 한다.

## 적용 기준과 읽을 원천

fullops-common-0.3.3의 rules/common README·coding-style·testing·security, project.md, docs/agents/document-writing.md, FULLOPS.md, contexts/dev.md. 기준 SHA의 같은 규약 사용. 기존 미해결 critical/high 차단 유지. ponytail full. 예외 없음.
필수 원천: adapters/README.md, adapters/src/login.ts, adapters/src/mcp.ts, deploy/knowslink/tunnel/public-ingress.yml, deploy/knowslink/beta.sh, deploy/knowslink/access_apply.py, .fullops-squad/docs/operations/ops-guide.md, .fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-001-TESTER-review/report.md.

## 해야 할 일과 소유권

- [ ] 기존 구현을 재사용해 가장 작은 기술 계획을 같은 exec-plan에 적고 필요한 코드·설정·플러그인 설치/안내·테스트만 보완한다. 비용 없는 Tunnel은 자체 Google 세션/서명·회원 소유권을 유지한다. 관리자/owner/test API는 인터넷에서 origin 인증만 믿고 모두 노출하지 말고 최소 deny 또는 기존 보호를 보존하는 적용안을 제공한다.
- [ ] Google 연결 도구가 설치 직후 도메인 설정으로 시작 가능한지 검토한다. 설치패키지/manifest/skills의 오래된 held 안내가 실제 연결을 막으면 필요한 최소 범위 수정. 기존 trial 호환·기본 비명시 발송 금지는 유지.
- [ ] 앞 리뷰의 medium 익명 Device 요청 cap2000/24h 보존과 low 타 기존회원 callback 분기 미검사 중 현재 공개수락에 필요한 최소 보완을 판단하고 구현한다. 별도 큰 리팩터링 금지. 지문 비교 경고·사용자 명시적 동의 유지.
- [ ] 기술적으로 불필요한 기능은 추가하지 않는다. 새 auth framework/remote OAuth 서버/유료 보안계층/자동wake/dots는 범위 밖이다.
- [ ] DEV 소유 제품/배포자료/기술문서 D03 D10 D12 및 필요 D11 D13 갱신. coor는 운영실제 적용과 PLANS 최신 사용자기록만 소유. 다른 인박스 수정 금지.

## 검증·완료 기준

- 로컬 Docker 금지. .env.server의 비밀 SSH로 기존 운영 서버의 고유 임시 checkout/격리DB/포트에서 등록 make lint/test와 영향 있는 좁은 경계를 검증한다. 운영 실제 서비스를 worker가 변경하지 않는다. 기존 테스트증거는 코드/조건 동일할 때만 재사용한다.
- Linux 서버/Windows 클라이언트와 Bot설치에 필요한 작은 검사. command 실제 exit code 보존. 개인키/token/Google/SSH/CF 자격 출력·argv·Git 금지.
- 사용자 관측 완료: 플러그인 설치→도메인으로 connect→사용자 Google 로그인/동의→독립agent/키 로컬저장→관계수락→양쪽 실제대화. DEV는 준비와 합성검증을 완료하며 실제로그인 성공으로 표기하지 않는다. coor가 독립 fixed-SHA 리뷰·좁은QA 후 운영적용·실계정 수락한다.
- UI 재설계 없음. 기존 UI/컴포넌트 유지. 새 시각 변화가 필요하면 영향과 직접검수 항목 명시.
- 탐색 packet을 확인하고 빠진 후보·partial/unknown을 수동 점검. packet-outcomes 모든 비선택쌍 기록. work.py finish/archive 후 commit 및 lint.py 최종HEAD gate. ERROR0, 실행불가0; 경고근거 명시.
- 완료는 실제 새 task/dispatch의 worker_done로 정규 한줄 '[완료] SAR-GOOGLE-CONNECT-002-DEV | 브랜치 fullops/dev | SHA <최종결과SHA> | ...' 전달. body의 SHA 라벨은 결과 하나만. 구현자 세션ID·lint·운영 적용/복구 인계 포함.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-GOOGLE-CONNECT-002-DEV-packet.json` / SHA `5af28c9fcfc289a845f9730aed34d32308e8b5a6` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/operations/ops-guide.md` (document_read, document_update) · 줄 7, 7, 14, 36, 38, 46, 47, 49, 51, 71, 85, 99, 101, 122, 133, 167, 175, 197, 205, 213, 235, 259, 267, 302, 304, 306, 306, 308, 310, 310, 310, 312, 312, 316, 320, 322 · inferred · 필수
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 2, 2, 6, 6, 7, 7, 14, 14, 14, 14, 19, 19, 21, 21, 22, 23, 28, 28, 32, 33, 41, 42, 45, 45 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 45, 51, 52 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 6, 7, 69, 69, 70, 70, 72, 72, 72, 73, 73 · inferred
- `.fullops-squad/contexts/ops.md` (document_read) · 줄 29 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read, document_update) · 줄 7, 7, 84, 161, 161, 161, 163 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read) · 줄 7, 7, 82, 82, 82, 84 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 7, 7, 66, 66, 66, 68 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 7, 90, 121, 172, 216, 216, 216, 224 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 24 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 106 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md` (document_read) · 줄 86 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read, document_update) · 줄 7, 7, 28, 52, 72, 166, 168, 170, 172, 173, 174, 176, 176, 178, 178, 178 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/COOR/dev-fix-2-final/coordinator-handover-supplement.md` (document_read) · 줄 33 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.2.0-review/report.md` (document_read) · 줄 54, 58 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-review/report.md` (document_read) · 줄 32, 43, 53 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-TESTER.md` (document_read) · 줄 96 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-002-TESTER.md` (document_read) · 줄 17, 36 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-DEPLOY-001-OPS-FINAL-review/report.md` (document_read) · 줄 32, 33, 35, 47, 56, 63 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-001-TESTER-test/report.md` (document_read) · 줄 2, 2, 6, 6, 7, 10, 10, 17, 33, 54, 71, 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-LOGIN-001-REVIEW-review/report.md` (document_read) · 줄 2, 6, 7, 10, 16, 18, 20, 26, 28, 32, 36 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FINAL.md` (document_read) · 줄 28 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX.md` (document_read) · 줄 28 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md` (document_read) · 줄 26 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-TESTER.md` (document_read) · 줄 47 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-QA-RECORD-REVIEW-review/report.md` (document_read) · 줄 25 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER.md` (document_read) · 줄 63 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 107 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER.md` (document_read) · 줄 36 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX-2.md` (document_read) · 줄 55 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-TESTER.md` (document_read) · 줄 51, 76, 109 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER-FINAL.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-SETUP-001-TESTER.md` (document_read) · 줄 18, 50 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.2.0.md` (document_read) · 줄 45 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 20, 47, 48, 75, 109, 110 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-DEPLOY-001-OPS.md` (document_read) · 줄 23, 34, 44 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-001-DEV.md` (document_read) · 줄 2, 2, 6, 6, 7, 10, 10, 17, 22, 24, 28, 34, 36, 40 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md` (document_read) · 줄 2, 6, 7, 10, 14, 16, 18, 27, 35, 39 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-001-DEV.md` (document_read) · 줄 23, 43 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW.md` (document_read) · 줄 14, 34 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md` (document_read) · 줄 16, 25, 69, 195 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md` (document_read) · 줄 56, 68, 78 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-TRIAL-CLEANUP.md` (document_read) · 줄 14, 20, 25, 33, 44, 45 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-PUBLIC-POLICY-001.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` (document_read) · 줄 63 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2.md` (document_read) · 줄 91 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-3.md` (document_read) · 줄 86 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX.md` (document_read) · 줄 74 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPEN-PREP.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 28, 148, 202 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 7, 7, 21, 23, 24, 31, 35, 46, 48, 74, 76, 78, 78, 80, 82, 82, 82, 84, 86 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 7, 7, 100, 100, 105, 107 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 51, 73, 75 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 32, 58 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 158, 160, 160 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 46, 62, 84 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/README.md` (document_read) · 줄 40, 50 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md` (document_read) · 줄 129, 146, 154 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/product.md` (document_read) · 줄 119, 125 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_designer.md` (document_read) · 줄 204 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_ops.md` (document_read) · 줄 12, 20, 24, 30, 72, 78, 94, 121, 125, 165, 191, 199, 207, 213, 215 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_dev.md` (document_read) · 줄 236, 251, 257 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_ops.md` (document_read) · 줄 236, 237 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_dev.md` (document_read) · 줄 147 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_ops.md` (document_read) · 줄 45, 72, 87, 113 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 43, 124, 148, 244, 318 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 114, 184, 241, 317, 335, 383 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_tester.md` (document_read) · 줄 42, 108, 171, 240 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 37, 66, 69, 70, 73, 73, 77, 78, 83, 83, 88, 92, 108, 115, 120, 121, 122, 129, 131, 139, 139, 144, 287, 297, 299, 305, 307, 309, 313, 314, 315, 316, 319 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_tester.md` (document_read) · 줄 297 · inferred
- `.fullops-squad/handovers/logs/2026-10-10_to_dev.md` (document_read) · 줄 9, 9, 10, 10, 13, 13, 13, 17, 17, 18, 24, 24, 24, 27, 27, 40, 46, 51, 57, 62, 63, 67, 67, 71, 77, 77, 83, 83, 85, 120, 257, 257, 257, 262, 262, 263, 267, 267, 275, 276, 276 · inferred
- `.fullops-squad/handovers/logs/2026-10-10_to_tester.md` (document_read) · 줄 9, 9, 10, 10, 13, 13, 13, 17, 17, 18, 25, 25, 28, 28, 35, 37, 37, 37, 41, 41, 45, 45, 57, 57, 60, 61, 61, 63, 63 · inferred
- `.fullops-squad/orca-agents.md` (document_read) · 줄 47 · inferred
- `adapters/README.md` (document_update) · 줄 전체/미확인 · inferred
- `adapters/src/login.test.ts` (impact_check) · 줄 1, 34 · inferred
- `adapters/src/mcp.ts` (direct_edit) · 줄 302 · inferred
- `cmd/relay/main.go` (impact_check) · 줄 56, 66, 84 · inferred
- `deploy/knowslink/access_trial_plan.py` (impact_check) · 줄 1 · inferred
- `deploy/knowslink/beta.sh` (direct_edit) · 줄 전체/미확인 · unknown
- `deploy/knowslink/tunnel/public-ingress.yml` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/device.go` (impact_check) · 줄 1, 114, 117 · inferred
- `internal/relay/google_test.go` (impact_check) · 줄 1 · inferred
- `scripts/check_public_ingress.py` (impact_check) · 줄 1, 19 · inferred
미확인 5건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

작성 후 finish한다.
