---
title: SAR-SERVICE-UI-001-OPS — 수락한 회원 웹·플러그인의 상태 보존 운영 배포
status: draft
updated: 2026-10-11
owner: ops
tasks: [SAR-SERVICE-UI-001-OPS]
summary: 수락한 회원 웹·플러그인의 상태 보존 운영 배포
attempt: f34e8cdec6d3408eacc74ae2b5170e08
base: 95df83c22c1559453a7c73504aa248088a8e4ed0
subagent_level: standard
test_level: lite
---

# SAR-SERVICE-UI-001-OPS — 수락한 회원 웹·플러그인 보존 배포

- 상태: ready. 수락 main 고정 SHA 95df83c22c1559453a7c73504aa248088a8e4ed0. 이 SHA만 상태 보존 배포한다.
- 역할: ops, C:/Users/shin/orca/workspaces/KnowsLink/ops, fullops/ops. 복귀 run_7176a2ad6ef1, coor term_e61d3e14-29e9-4954-943a-4a75707c82de. 실제 task/dispatch는 preamble.
- 공통 기준: fullops-common-0.3.3와 세 규칙, FULLOPS.md/project.md/document-writing. Ponytail full; lite/subagent standard. 등록 Claude Sonnet5.5 high fresh. Codex Trust 임의 승인 금지.
- 사용자 승인: 기존 .env.server의 Linux 서버에서 직접 시험·진행, link.knowslog.com 상용 UI 개선·기존 플러그인 구현의 운영 반영. 추가 유료서비스·로컬 Docker·파일/컨테이너 삭제·이동·광범위 정리 금지. F: 사고 증거 보존.

## 보존과 배포

운영 위치 /home/shin/deploy/knowslink, 상태 /home/shin/deploy/knowslink-state. 기존 .env.server 자격은 메모리에서 읽고 출력하지 않는다. 기존 coor_remote.py SSH helper와 OPS-LIVE/D12 runbook을 직접 읽어 재사용한다. 이전 운영 제품 eac603b6a41d2ce002e698600ac66b8fc67ad762를 복귀 기준으로 확인한다. 새 기준은 실제 작업 직전 수집한다. 이전 수치에 맞추려고 DB를 바꾸지 않는다.

배포 전후 relay DB·member/owner/agent/key/pair·기존 메시지/구독·Google/Tunnel/상태 자격과 config 해시·볼륨·공유 서비스 Orca200/s8200/mcp401/myportfolio·기존 포트·Tunnel ingress를 비교한다. 민감 데이터 원문·SQL dump는 보존/출력하지 않고 집계/비밀값을 제외한 digest만 기록한다. 새 integration 시험에서 생긴 변화는 원인과 구분한다. 기존 state·DB·keys·cert·Google/Tunnel credential·shared resources를 재생성/철회/삭제하지 않는다.

수락 SHA에 실제 Linux 배포 checkout이 clean detached인지 확인한다. 필요 artifact와 binary/container만 기존 runbook대로 갱신한다. 외부 의존성/유료 서비스를 추가하지 않는다. DB migration 변경0이면 기존 schema·state 유지하고 초기화하지 않는다. 실패하면 확인된 기존 배포로 복귀하고 사유와 증거를 보존한다. 새 설치 plugin zip이 webhook wake 코드와 정확한 배포 SHA에서 생성됐는지 package verify/hash로 확인한다. Bot runtime/env는 coor 후속이며 건드리지 않는다.

## 공개 검증과 완료

공개 / 및 /home·member 경로의 상태/CSP/로그인 경로, readiness·public deny·기존 회귀를 기존 도구로 확인한다. 새 Google 실제계정 로그인은 coor Orca browser 후속이며 자동 PASS로 기록하지 않는다. 웹 조작 필요 시 Orca 내장 브라우저만. Grok Bot UI 접근 금지. secret input/token이 든 페이지 캡처 금지.

본인 OPS-LIVE에 baseline/after/명령·종료코드·SHA·artifact hash·불변식·복귀 방법·미검증을 보존한다. 필요한 D12/D13 상태만 실제 결과로 갱신한다. phase/contexts/packet outcomes, 대화 미참조 인계 점검과 lint 근거를 기록한다. work.py finish로 전문 아카이브·인박스 비움·커밋 후 [완료] SAR-SERVICE-UI-001-OPS | 최종 SHA <실제 결과40자리>로 보고한다. 결과 SHA·운영 SHA를 구분한다. main 결과 통합은 coor 담당이다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-SERVICE-UI-001-OPS-packet.json` / SHA `95df83c22c1559453a7c73504aa248088a8e4ed0` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 73 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-TESTER-review/report.md` (document_read) · 줄 2, 6, 10, 42, 46 · inferred · 필수 · {'relevant': 0.8, 'evidence': 0.81, 'contradicts': 0.26, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL/report.md` (document_read) · 줄 6, 10, 16, 31, 41, 85 · inferred · 필수 · {'relevant': 0.69, 'evidence': 0.78, 'contradicts': 0.14, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/operations/ops-guide.md` (document_read, document_update) · 줄 7, 18, 20, 32, 37, 42, 44, 72, 83, 85, 85, 194, 194, 194, 194, 204, 204, 208, 224, 251, 289, 297, 303, 303, 303, 305, 307, 312, 336, 336, 344, 353, 356, 368, 374 · inferred · 필수 · {'relevant': 0.5, 'evidence': 0.78, 'contradicts': 0.77, 'injection': 0.14, 'decision': 'conflict', 'reason': 'contradicts task'}
- `.fullops-squad/docs/operations/transition.md` (document_read, document_update) · 줄 7, 15, 15, 21, 30, 49, 49, 81, 85, 89, 89, 97, 103, 106 · inferred · 필수 · {'relevant': 0.36, 'evidence': 0.65, 'contradicts': 0.69, 'injection': 0.07, 'decision': 'conflict', 'reason': 'contradicts task'}
- `.fullops-squad/handovers/to_ops.md` (document_read) · 줄 2, 6, 14, 16, 27, 31, 33, 33, 33, 42, 42, 42, 43, 44, 47, 47 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 15, 40, 50, 72 · inferred · 필수 · {'relevant': 0.67, 'evidence': 0.84, 'contradicts': 0.24, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.dockerignore` (impact_check) · 줄 1 · inferred
- `.fullops-squad/PLANS.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 29, 30 · inferred
- `.fullops-squad/contexts/coor.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 30, 32, 70, 72, 74, 76, 78, 80 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 28, 32, 38 · inferred
- `.fullops-squad/contexts/ops.md` (document_read) · 줄 63, 69, 71, 75, 77, 77 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 6, 15, 28, 59, 69, 71, 71, 74 · inferred
- `.fullops-squad/docs/agents/document-writing.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 16, 20, 25, 81, 119, 160, 179, 187, 202 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 21 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 25, 25, 44, 54, 55, 141, 175, 185, 188, 230, 236, 241, 249, 255 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 7, 97 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 25, 29, 35, 81, 97, 118 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md` (document_read) · 줄 21 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 6, 14, 31, 41, 63 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 7, 314, 406, 408 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 7, 19, 24, 88, 99, 116, 119, 168, 177, 191, 193, 203, 221, 234, 234 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 16, 32, 43, 49, 53, 62, 68, 85, 87, 91 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/COOR/dev-fix-2-final/coordinator-handover-supplement.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.2.0-review/report.md` (document_read) · 줄 32, 36, 54 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.3.0-FINAL-review/report.md` (document_read) · 줄 50 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-LIVE/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-REVIEW-review/report.md` (document_read) · 줄 20, 35, 46 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-TESTER.md` (document_read) · 줄 24, 28, 36 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-FINAL-review/report.md` (document_read) · 줄 12, 27, 27, 36, 49, 85, 86, 95 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-N1-review/report.md` (document_read) · 줄 12, 73, 73, 73, 98 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-review/report.md` (document_read) · 줄 12, 32, 58, 74, 75, 86 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-TESTER-PUBLIC.md` (document_read) · 줄 20, 21, 60, 61, 75, 76, 77 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-TESTER.md` (document_read) · 줄 30, 50, 51, 62, 65, 71, 72, 73, 74, 84, 91 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-002-TESTER.md` (document_read) · 줄 103 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-DEPLOY-001-OPS-FINAL-review/report.md` (document_read) · 줄 29, 34, 61 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-001-TESTER-test/report.md` (document_read) · 줄 16, 25, 28, 29, 30, 31, 32, 37, 79, 83 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-DEV-review/report.md` (document_read) · 줄 19, 24, 28, 40 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-DEV/report.md` (document_read) · 줄 18, 21, 22, 24, 32 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-LIVE.md` (document_read) · 줄 12, 14, 14, 16, 26 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-LOGIN-001-REVIEW-review/report.md` (document_read) · 줄 26, 36 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GROK-WEBHOOK-001-TESTER-RESULT-REVIEW-review/report.md` (document_read) · 줄 14, 28 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GROK-WEBHOOK-001-TESTER-REVIEW-review/report.md` (document_read) · 줄 16, 20, 20, 74, 91 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GROK-WEBHOOK-001-TESTER-test/report.md` (document_read) · 줄 13 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FINAL-review/report.md` (document_read) · 줄 22, 117, 151 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FIX-review/report.md` (document_read) · 줄 143, 187 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FINAL.md` (document_read) · 줄 86 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX.md` (document_read) · 줄 27 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md` (document_read) · 줄 25, 52 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER.md` (document_read) · 줄 60, 62, 90, 104, 109, 115, 118, 130, 135 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-FIX-review/report.md` (document_read) · 줄 38, 48 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-TESTER.md` (document_read) · 줄 92, 98, 124, 131 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-review/report.md` (document_read) · 줄 49 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-DEV-TESTER.md` (document_read) · 줄 40, 46, 103 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-DEV-review/report.md` (document_read) · 줄 25, 33, 60 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-TESTER.md` (document_read) · 줄 17, 18, 51, 59, 69, 76, 76, 110, 117, 124, 126, 127, 127 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-review/report.md` (document_read) · 줄 14, 23, 35, 35, 58, 65 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TESTER.md` (document_read) · 줄 45, 71, 98 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER.md` (document_read) · 줄 27, 35, 57, 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-review/report.md` (document_read) · 줄 37 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-FINAL-RECORDS-review/report.md` (document_read) · 줄 40, 42 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-FIX-TESTER.md` (document_read) · 줄 37 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 16, 17, 21, 53, 76, 86, 99, 100 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-TESTER-test/probe-failures.md` (document_read) · 줄 25 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-TESTER.md` (document_read) · 줄 116, 139 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/scenario.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-FIX-TESTER.md` (document_read) · 줄 38, 45 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW-review/report.md` (document_read) · 줄 16, 20, 38, 78, 89, 90, 93, 131, 148, 164 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-REVIEW-review/report.md` (document_read) · 줄 16, 24, 31 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TESTER.md` (document_read) · 줄 15, 32, 39 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW-review/report.md` (document_read) · 줄 21 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-UI.md` (document_read) · 줄 30 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER.md` (document_read) · 줄 33 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 18, 20, 22, 29, 95, 114 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER.md` (document_read) · 줄 23, 32 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX-2.md` (document_read) · 줄 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 28, 39 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-SERVICE-001-REVIEW-review/report.md` (document_read) · 줄 6, 12, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-TESTER-RESULT-REVIEW-review/report.md` (document_read) · 줄 6, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL-RESULT-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-RESULT-REVIEW-review/report.md` (document_read) · 줄 6, 20, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-review/report.md` (document_read) · 줄 2, 6, 10, 15, 19, 45, 58, 67, 68, 70 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-test/report.md` (document_read) · 줄 2, 6, 10, 25, 26, 34, 62 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-2-review/report.md` (document_read) · 줄 6, 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-review/report.md` (document_read) · 줄 6, 10, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md` (document_read) · 줄 6, 10, 28, 34, 66, 72, 107, 109, 121, 179, 181 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-DEV-099-review/report.md` (document_read) · 줄 37, 43 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-INTEGRATION-FINAL-review/report.md` (document_read) · 줄 23, 31, 39 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-INTEGRATION-SOURCE-review/report.md` (document_read) · 줄 7, 13, 23, 32, 34, 45, 47, 63, 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-TESTER.md` (document_read) · 줄 24, 31, 31, 34, 46, 51, 79, 100, 101, 107, 109 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS-REVIEW-review/report.md` (document_read) · 줄 7, 16, 16, 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-review/report.md` (document_read) · 줄 19, 26, 32, 32, 45, 54, 55 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-test/report.md` (document_read) · 줄 38, 39, 40, 48, 51, 52, 53, 54, 55, 56, 80, 95, 95, 97 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL-REVIEW-review/report.md` (document_read) · 줄 7, 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL/report.md` (document_read) · 줄 49, 83, 91 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OPS-REVIEW-review/report.md` (document_read) · 줄 20, 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-REVIEW-review/report.md` (document_read) · 줄 19, 40 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-TESTER-test/report.md` (document_read) · 줄 12, 16, 31 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-VISUAL/report.md` (document_read) · 줄 82, 84 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-WINDOWS-SHELL-GUARD-001-FIX-REVIEW-review/report.md` (document_read) · 줄 18, 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-WINDOWS-SHELL-GUARD-001-REVIEW-review/report.md` (document_read) · 줄 19, 37, 49 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-WINDOWS-SHELL-GUARD-001-TESTER-test/report.md` (document_read) · 줄 58, 67 · inferred
- `.fullops-squad/docs/evaluations/scenarios/README.md` (document_read) · 줄 17, 68 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-BETA-001-TESTER-PUBLIC.md` (document_read) · 줄 24, 24, 24, 37 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-BETA-001-TESTER.md` (document_read) · 줄 13, 24, 24, 24, 29, 40, 50, 52 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-BETA-002-TESTER.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER-FINAL.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER-FIX.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER.md` (document_read) · 줄 13 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-AGENTS-001-TESTER.md` (document_read) · 줄 29 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-IDENTITY-001-TESTER.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-SETUP-001-TESTER.md` (document_read) · 줄 19, 26, 29, 38, 40 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER.md` (document_read) · 줄 26, 30, 31, 32, 33, 34, 35 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.10.md` (document_read) · 줄 14, 17 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.11.md` (document_read) · 줄 16, 18 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.12.md` (document_read) · 줄 14, 22, 32, 36, 38, 42, 42, 48 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.13.md` (document_read) · 줄 26, 32 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.14.md` (document_read) · 줄 18, 25, 39, 44, 49, 51, 61 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-099.md` (document_read) · 줄 15, 41, 46 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.2.0.md` (document_read) · 줄 18, 32, 57, 57 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.3.0.md` (document_read) · 줄 34, 40, 42, 46, 59, 60 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-AUTO-RECEIVE-001-DEV.md` (document_read) · 줄 14, 44, 98, 116 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 14, 31, 42, 46, 50, 84, 84, 97, 113, 113, 115, 115, 115 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-DEPLOY-001-OPS.md` (document_read) · 줄 22, 36, 44 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-001-DEV.md` (document_read) · 줄 22, 32, 36 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-002-DEV.md` (document_read) · 줄 16, 21, 25, 33, 35 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md` (document_read) · 줄 25, 26, 27, 27, 37, 45 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-DEV.md` (document_read) · 줄 16, 24, 33, 48, 67, 76 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-TESTER.md` (document_read) · 줄 21, 33 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-001-DEV.md` (document_read) · 줄 59, 79, 80, 124, 153, 158, 159, 170 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-001-UI.md` (document_read) · 줄 100 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV-FIX.md` (document_read) · 줄 38, 48, 49, 51, 57, 81 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md` (document_read) · 줄 50, 71, 71, 72, 88, 102 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-DEV-TESTER.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md` (document_read) · 줄 27, 34, 47, 55, 57, 58, 58, 70, 87, 88, 98, 99, 111 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW.md` (document_read) · 줄 14, 19, 53, 89 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md` (document_read) · 줄 14, 31, 35, 35, 35, 45, 46, 46, 80, 101, 133, 191, 212, 223, 223, 223 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT.md` (document_read) · 줄 20, 32, 44, 64 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md` (document_read) · 줄 14, 36, 38, 50, 70, 88, 110 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-TRIAL-CLEANUP.md` (document_read) · 줄 16, 50, 50, 50, 66 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PREP-002.md` (document_read) · 줄 40 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 14, 16, 80 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX.md` (document_read) · 줄 18, 91, 109 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG.md` (document_read) · 줄 22, 41, 43 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW.md` (document_read) · 줄 29, 36 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-TESTER.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW.md` (document_read) · 줄 28 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-UI.md` (document_read) · 줄 22, 30, 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2.md` (document_read) · 줄 14, 93 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-3.md` (document_read) · 줄 15, 57, 88 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX.md` (document_read) · 줄 14, 76 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 38, 50 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (document_read) · 줄 6, 10, 26, 28, 32, 40, 48, 49, 66, 66, 74, 74, 80 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPEN-PREP.md` (document_read) · 줄 2, 6, 10, 14, 18, 21, 21, 22, 27 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 2, 6, 10, 30, 30, 30, 128, 169, 170, 171 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 6, 10, 15, 69, 92, 101, 105, 108 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-DEV.md` (document_read) · 줄 6, 10, 15, 31, 41, 44, 47, 56 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-TESTER.md` (document_read) · 줄 2, 6, 10, 12 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-VISUAL.md` (document_read) · 줄 2, 6, 10, 22, 33 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-TESTER.md` (document_read) · 줄 2, 6, 10, 20, 26, 30, 32, 39 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-VISUAL.md` (document_read) · 줄 2, 6, 10, 25 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SETUP-001-TESTER.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-DEV.md` (document_read) · 줄 25, 28, 55, 61, 62, 63, 68, 86 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-INSTALL.md` (document_read) · 줄 16, 28, 28, 31, 31, 31, 32, 33, 38, 39, 40, 49, 50, 51, 52, 53, 54, 62, 63, 64, 75, 77, 80, 82, 97 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV.md` (document_read) · 줄 27, 43, 68, 69 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS.md` (document_read) · 줄 25, 29, 29, 34, 34 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER.md` (document_read) · 줄 33 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OPS.md` (document_read) · 줄 23, 29, 32, 33, 33, 38, 38 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-REMOTE.md` (document_read) · 줄 23, 32, 36, 49, 60, 66, 72 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-WINDOWS-SHELL-GUARD-001-DEV.md` (document_read) · 줄 66, 88 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 7, 27, 125 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 6, 12, 14, 18, 20, 23, 24, 27 · inferred
- `.fullops-squad/docs/planning/SAR-PREP-002-request.md` (document_read) · 줄 14, 16 · inferred
- `.fullops-squad/docs/planning/SAR-SETUP-001-request.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 7, 40, 56, 66 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 7, 25, 35, 67, 83, 109, 154 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 7, 15, 23, 63, 76, 125, 128, 152 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 36 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/README.md` (document_read) · 줄 30, 38 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md` (document_read) · 줄 30, 32, 38, 118, 121, 129, 129, 136 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/product.md` (document_read) · 줄 113 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/protocol.md` (document_read) · 줄 364, 370, 384 · inferred
- `.fullops-squad/handovers/SAR-MVP-001-REVIEW.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/handovers/_TEMPLATE.md` (document_read) · 줄 18, 80 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_designer.md` (document_read) · 줄 28, 104, 110 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_dev.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/handovers/logs/2026-10-03_to_ops.md` (document_read) · 줄 52, 68, 84, 94, 115, 148, 189, 213, 215, 219, 257 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_tester.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/handovers/logs/2026-10-04_to_dev.md` (document_read) · 줄 27, 33, 37, 90, 125, 127, 149, 161, 196, 197, 244, 255, 257, 259, 275, 285, 290, 294 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_ops.md` (document_read) · 줄 24, 39, 49, 72, 100, 161, 197, 229, 236, 237, 253, 253, 253, 254, 269, 270, 283, 296, 296, 296, 299 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_tester.md` (document_read) · 줄 24, 39, 40, 53, 75, 125, 136, 167, 207, 222 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_designer.md` (document_read) · 줄 9, 16, 20, 79, 85, 87, 93, 114, 126, 131, 156 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_dev.md` (document_read) · 줄 16, 33, 35, 37, 43, 51, 93, 93, 121, 125, 134, 140, 154, 181 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_ops.md` (document_read) · 줄 22, 39, 45, 64, 80, 80, 87, 87, 87, 90, 93, 97, 101, 105, 116, 130, 142, 142, 142, 144, 147, 151, 155, 159, 171 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_tester.md` (document_read) · 줄 16, 33, 35, 37, 49, 76, 85, 86 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 32, 35, 102, 105, 106, 142, 235, 248, 249, 257, 276, 284, 306, 386 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 28, 37, 69, 73, 77, 108, 125, 173, 176, 203, 205, 232, 245, 246, 254, 280, 280, 282, 283, 304, 312, 331, 336, 357, 380, 382, 384 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_tester.md` (document_read) · 줄 36, 97, 100, 142, 162, 175, 176, 184, 207, 228, 240, 252, 254 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_designer.md` (document_read) · 줄 29, 29, 58, 69, 76, 133, 182, 183, 184, 220, 222, 226, 228 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 24, 37, 88, 92, 108, 111, 113, 148, 149, 229, 230, 259, 261, 265, 267, 301, 301, 309 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_tester.md` (document_read) · 줄 48, 54, 60, 70, 75, 139, 180, 181, 182, 267 · inferred
- `.fullops-squad/handovers/logs/2026-10-10_to_dev.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/handovers/logs/2026-10-10_to_tester.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/handovers/logs/2026-10-11_to_designer.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/handovers/logs/2026-10-11_to_dev.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/handovers/logs/2026-10-11_to_ops.md` (document_read) · 줄 31, 43, 43, 45, 53, 54, 54, 61, 65, 65, 68, 77, 80, 83, 85, 85, 91, 93, 99, 104, 112, 114, 124, 136, 138, 138, 147, 147, 204, 205, 216, 235, 249, 254, 256, 256, 262 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_tester.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/handovers/logs/SAR-SETUP-001-DEV-REVIEW.md` (document_read) · 줄 29, 37, 45 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-INTEGRATION-REVIEW.md` (document_read) · 줄 7, 12, 13, 20, 22, 22, 26, 28, 32, 36, 40, 42 · inferred
- `.fullops-squad/lint/README.md` (document_read) · 줄 29, 48, 49, 57, 61, 62, 85 · inferred
- `.fullops-squad/orca-agents.md` (document_read) · 줄 26, 27, 28, 29, 32, 56, 67, 128, 147 · inferred
- `.fullops-squad/rules/delegation.md` (document_read) · 줄 39 · inferred
- `.gitignore` (impact_check) · 줄 전체/미확인 · inferred
- `AGENTS.md` (document_read) · 줄 전체/미확인 · inferred
- `CLAUDE.md` (document_read) · 줄 전체/미확인 · inferred
- `Dockerfile` (impact_check) · 줄 12 · inferred
- `GEMINI.md` (document_read) · 줄 전체/미확인 · inferred
- `Makefile` (impact_check) · 줄 2, 4, 17, 36, 46 · inferred
- `adapters/.claude-plugin/plugin.json` (impact_check) · 줄 9 · inferred
- `adapters/.cursor-plugin/plugin.json` (impact_check) · 줄 전체/미확인 · inferred
- `adapters/README.md` (document_read) · 줄 29, 34, 36, 36, 41, 42, 48, 48, 63, 74, 80, 81, 96, 107, 169, 177, 183, 183, 187, 222, 244, 263, 271, 334, 355, 357 · inferred
- `adapters/manifest.json` (impact_check) · 줄 24 · inferred
- `adapters/package.json` (impact_check) · 줄 11, 11, 14, 15, 17 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 3, 12, 14, 14, 16, 20, 26, 36, 40, 40, 44, 44, 46 · inferred
- `adapters/src/grok-wake.test.ts` (impact_check) · 줄 25, 79, 95 · inferred
- `adapters/src/inbox.ts` (impact_check) · 줄 1, 170, 171, 172, 246 · inferred
- `adapters/src/index.ts` (impact_check) · 줄 6, 25 · inferred
- `adapters/src/login.test.ts` (impact_check) · 줄 25, 38, 41, 43, 43, 119 · inferred
- `adapters/src/mcp.ts` (impact_check) · 줄 1, 22, 27, 64, 78, 89, 130, 151, 156, 313 · inferred
- `adapters/src/public-check.ts` (impact_check) · 줄 1, 15, 18, 23, 23, 29, 47, 59, 68, 84, 91, 93, 99 · inferred
- `adapters/src/trial-check.ts` (impact_check) · 줄 12, 113 · inferred
- `adapters/src/trial-cli.ts` (impact_check) · 줄 3, 33 · inferred
- `adapters/src/trial-setup.ts` (impact_check) · 줄 19, 50, 78, 117, 124 · inferred
- `cmd/migrate/main.go` (impact_check) · 줄 2, 2, 18, 19, 25 · inferred
- `cmd/migrate/main_test.go` (impact_check) · 줄 2, 2, 14, 17, 21 · inferred
- `cmd/relay/main.go` (impact_check) · 줄 1, 2, 2, 19, 22, 28 · inferred
- `cmd/relay/main_test.go` (impact_check) · 줄 2, 2 · inferred
- `deploy/knowslink/access_trial_plan.py` (impact_check) · 줄 26, 49, 91 · inferred
- `deploy/knowslink/beta.sh` (direct_edit, impact_check) · 줄 전체/미확인 · unknown · {'relevant': 0.51, 'evidence': 0.58, 'contradicts': 0.68, 'injection': 0.06, 'decision': 'conflict', 'reason': 'contradicts task'}
- `deploy/knowslink/compose.ops.yaml` (impact_check) · 줄 32 · inferred
- `deploy/knowslink/test_public_config.sh` (impact_check) · 줄 2, 18, 19, 25, 32, 35, 38 · inferred
- `deploy/knowslink/verify.py` (direct_edit, impact_check) · 줄 15, 19, 26, 35, 41, 49, 66 · inferred · {'relevant': 0.76, 'evidence': 0.79, 'contradicts': 0.32, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `internal/config/config.go` (impact_check) · 줄 2 · inferred
- `internal/database/db.go` (impact_check) · 줄 5 · inferred
- `internal/database/models.go` (impact_check) · 줄 5 · inferred
- `internal/database/relay.sql.go` (impact_check) · 줄 6 · inferred
- `internal/relay/admission_integration_test.go` (impact_check) · 줄 4, 121, 155 · inferred
- `internal/relay/api_rate.go` (impact_check) · 줄 2 · inferred
- `internal/relay/capacity.go` (impact_check) · 줄 2, 249 · inferred
- `internal/relay/cleanup.go` (impact_check) · 줄 2 · inferred
- `internal/relay/cleanup_flood_integration_test.go` (impact_check) · 줄 4, 139, 146 · inferred
- `internal/relay/cleanup_unit_integration_test.go` (impact_check) · 줄 4 · inferred
- `internal/relay/connections_test.go` (impact_check) · 줄 2, 29, 32, 303, 304, 333 · inferred
- `internal/relay/device.go` (impact_check) · 줄 1, 2, 21, 26, 116, 117 · inferred
- `internal/relay/google_test.go` (impact_check) · 줄 2 · inferred
- `internal/relay/mail.go` (impact_check) · 줄 2 · inferred
- `internal/relay/member_agents.go` (impact_check) · 줄 2, 61, 62, 304, 305 · inferred
- `internal/relay/member_pages.go` (impact_check) · 줄 2, 75, 80, 83, 84, 85, 95, 96, 102, 103, 124, 125, 133, 134, 136, 139, 142, 153, 181 · inferred
- `internal/relay/member_receipt.go` (impact_check) · 줄 2 · inferred
- `internal/relay/member_ui_test.go` (impact_check) · 줄 2 · inferred
- `internal/relay/policy_test.go` (impact_check) · 줄 2 · inferred
- `internal/relay/protocol.go` (impact_check) · 줄 2, 264, 266, 289, 290 · inferred
- `internal/relay/protocol_test.go` (impact_check) · 줄 2, 48, 51, 61 · inferred
- `internal/relay/public_messages_integration_test.go` (impact_check) · 줄 4, 43, 54, 85, 88, 108, 121, 162, 175, 194 · inferred
- `internal/relay/public_text.go` (impact_check) · 줄 2 · inferred
- `internal/relay/registry.go` (impact_check) · 줄 2 · inferred
- `internal/relay/test_messages.go` (impact_check) · 줄 2 · inferred
- `internal/relay/test_messages_integration_test.go` (impact_check) · 줄 4 · inferred
- `internal/relay/testdata/consent_browser.mjs` (impact_check) · 줄 5, 33 · inferred
- `internal/relay/text_events.go` (impact_check) · 줄 2, 156 · inferred
- `internal/relay/text_events_integration_test.go` (impact_check) · 줄 4 · inferred
- `internal/relay/text_events_test.go` (impact_check) · 줄 2 · inferred
- `scripts/check_public_ingress.py` (impact_check) · 줄 5, 7, 8, 15, 20 · inferred
- `scripts/install_bot_mcp.sh` (impact_check) · 줄 25, 35, 60, 64, 77 · inferred
- `scripts/mail_sink.py` (impact_check) · 줄 38, 41, 53 · inferred
- `scripts/package_plugin.py` (direct_edit, impact_check) · 줄 25, 32, 99, 130, 164 · inferred · {'relevant': 0.32, 'evidence': 0.36, 'contradicts': 0.33, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `scripts/run_trial.py` (impact_check) · 줄 9, 33 · inferred
- `scripts/shell_guard.py` (impact_check) · 줄 13, 193, 213, 222, 234 · inferred
- `scripts/test_shell_guard.py` (impact_check) · 줄 2, 83, 120, 121, 145, 150, 153, 154, 155, 167, 179 · inferred
- `scripts/verify_grok_plugin.py` (impact_check) · 줄 13, 15, 18, 24, 29, 30, 30, 37, 38, 42 · inferred
- `scripts/verify_setup.py` (impact_check) · 줄 13, 14, 23, 24, 27, 36, 36, 45, 48, 53 · inferred
- `sqlc.yaml` (impact_check) · 줄 9 · inferred
미확인 12건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 `fullops/ops` / 운영 제품 SHA `95df83c22c1559453a7c73504aa248088a8e4ed0`(서버 배포 체크아웃 `git rev-parse HEAD`로 확인, detached, clean) / 결과 SHA는 이 기록을 커밋한 최종 HEAD이며 worker_done 첫 줄에 실제 40자로 적는다. 준비 main `0045089f2a2c3efe3e210eedc0c5cd4868487e42`는 제품 동일(`.fullops-squad` 밖 diff 0). 결과 SHA와 운영 SHA는 다르다.

변경 이유: 수락한 회원 웹 UI(V-01 수정 포함)와 Grok webhook wake 플러그인 코드를 기존 Linux 서버(`/home/shin/deploy/knowslink`)에 상태 보존으로 반영했다. 이전 운영 제품은 `eac603b6a41d2ce002e698600ac66b8fc67ad762`다.

수행: 읽기 전용 baseline → 로컬 bundle(`HEAD ^eac603b`, SHA256 `132262f16674b23cde05dc7a19504c320ce7753065274b2edbdab4c36aae2b59`) SFTP 전달·해시 일치·`git bundle verify`·fetch·조상 관계 → detached checkout `95df83c` → `deploy-history.log` 한 줄 → `beta.sh up`(relay만 재빌드·재생성) → `verify.py local`. 종료코드는 모두 0이다. migration·compose·deploy·Dockerfile·go.mod diff 0건, goose 버전 1 유지.

금지 준수: `prepare`·`deploy`·`backup`·`restore-verify`·`seed`·`owner-login`, raw `pg_dump`, Tunnel create/render/expose, 파일·컨테이너 삭제·정리를 실행하지 않았다. 새 DB backup은 없다(복귀는 코드 전용: `eac603b` checkout + `beta.sh up`, DB 미덮어쓰기). 비밀값·원문 state·로그 원문은 출력하지 않았다. 사용자의 Grok/Dots UI는 조작하지 않았다. Bot runtime·watcher·MCP는 건드리지 않았다.

보존(전→후): `.env`·`tunnel.uuid`·`access.aud`·Tunnel credential·`tunnel/config.yml`(`f5abe860…8512`)·compose 파일 SHA256 동일. relay_state 키별 개수·digest(Agents 28, Members 2, Messages 30, Owners 24, Pairs 13, Identities 2, Idempotency 30, Remote 6, Sessions 1)와 콘텐츠 md5 동일. `epoch`·`Rates`만 relay 동작으로 변함. Postgres 컨테이너(8일 유지)·볼륨·cloudflared·host PID 506937·myportfolio·orca 200/s8 200/mcp 401·다른 역할 컨테이너 변화 없음. 이전 LIVE 수치(Agents 26 등)와 달랐지만 이 과제 시작 전의 변화이며 맞추려고 DB를 바꾸지 않았다.

Plugin package: `git archive 95df83c` 사본(Node 22.22.2)에서 `npm ci`·build·`package_plugin.py --host all --verify` 종료 0. Grok zip SHA256 `762411df0e49dd0568c38b118d67c947910b39f56df5208afd897b5dc434cf94`, Claude zip `5388516edcfb992f777ef13f1c748df6b52a703fdc7200e391676da3e3ba8b9f`, Codex zip `d79abee4dc1639e8e83a12b10063a2fb63a04ba3f6492b3db62a53dab82c3410`, Desktop mcpb `bfd13777aa21dadabbf46a7f0439a6d8a9a04c7a90b0bff2401cb7aabb5aacb2`. Grok zip `dist/plugin.js`에서 wake 심볼을 확인했다. 서버 `knowslink-state/serviceui-20261011/pkg/build/`에 보존했다.

공개 smoke(무인증): `/` 200(새 디자인 토큰·건너뛰기 링크 확인, CSP·no-store·no-referrer), `/home` 303, `/healthz`·`/owner`·`/v1/test/pull` 404, OAuth discovery 200, `POST /mcp` 401, 합성 consent 401 HTML, SSE 401. 서버 `verify.py local` 통과, relay healthy·재시작 0·error/panic 0건. `verify.py regression`은 기존 baseline 결함(`projects`)으로 종료 1(codes·PID·containers는 동일).

검증하지 못한 것: 새 Google 실계정 로그인·`/home` 회원 목록·모바일 렌더·키보드 포커스·새 플러그인 실제 설치와 webhook 왕복은 coor 후속(Orca 내장 브라우저)이며 PASS로 기록하지 않았다. 운영 SHA `95df83c22c1559453a7c73504aa248088a8e4ed0`, URL `https://link.knowslog.com/`. 서버 `make test`·전체 QA는 반복하지 않고 고정 SHA의 tester·visual 증거를 재사용한다. 새 relay의 상태 JSON 필드 쓰기는 런타임에서 확인하지 않았다.

미검토·미해결: critical/high 없음. DEP-001: 의존성 변경 없음(이 과제는 코드·lock 무변경; `npm ci`는 lock 그대로). 하위 에이전트는 쓰지 않았다.

산출물·로그: [운영 LIVE 기록](../docs/evaluations/qa-reports/SAR-SERVICE-UI-001-OPS-LIVE.md), [실행 기록](../docs/exec-plans/phases/SAR-SERVICE-UI-001-OPS.md), D12·D13 갱신, `contexts/ops.md` 교훈, `docs/evaluations/jev/SAR-SERVICE-UI-001-OPS-packet-outcomes.json`.

### 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과 | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | 이 지시서, phases/SAR-SERVICE-UI-001-OPS.md | 확인 | 없음 |
| 구조와 구현/미완료 상태 | LIVE 기록 적용 요약 | 확인 | 코드 변경 없음, 운영 반영 완료 |
| 실행·검증 방법과 증거 | LIVE 기록 명령·종료코드·보존·smoke 절 | 확인 | baseline/after 원문은 서버 세션에서 redacted 요약만 남김 |
| 운영·복구 | LIVE 기록 복귀 방법 | 확인 | 새 backup 없음(지시서 금지). DB 복원 필요 시 coor 결정 |
| 다음 작업·담당·재개 조건 | phases 남은 일 | 확인 | coor 실계정 확인, plugin zip 전달 |
| 로컬 링크·절 접근/지원 한계 | 상대 링크 | 확인 | 서버 산출물은 이 PC 밖 |
| 정리 후 정본 접근 | main 통합 후 | 미확인 | coor 통합 후 |

### Lint·검증 근거

(lint 실행 뒤 아래에 기록)

## 지시 전제와 충돌 — 먼저 확인
Jev context는 D12/D13/이전 OPS 로그를 conflict로 표시했다. 이 문서들을 먼저 읽는다. D12의 1–10장은 역사적 조사·계획이고 SSE 운영 후보 절의 미적용 문장은 상단 적용 결과가 이전 사실로 명시했다. D13/최근 OPS 로그의 실제 운영 제품은 eac603b다. 신규 UI는 독립 검수 후 수락 SHA를 추가할 때까지 미배포다. 이전 held/fail·미검증을 새 PASS로 바꾸지 않는다. 배포 직전 실제 서버가 정본과 다르면 coor에게 보고하고 무작정 checkout/up하지 않는다. 기존 승인된 공개 Google/member 서비스와 Tunnel allowlist를 유지하는 제품 갱신이며 Access/새인증 정책 변경은 없다.

## 수락 고정 SHA와 검수 — 실행 정본
main/origin/main 95df83c22c1559453a7c73504aa248088a8e4ed0에 UI 원후보 c7828b2·V-01 수정제품609fcfd/최종후보295417a·독립 tester c1ab9976·직접 시각결과fe47f5c7가 통합됐고 fetch/모든 조상 관계 확인했다. 새 product C/H0, V-01 resolved. 고정 product review: SAR-SERVICE-UI-001-FIX-TESTER-review/report.md 및 root FIX-TESTER-RESULT-REVIEW/report.md(실제 구현 FIX-DEV task key로 check2/9통과); 시각: FIX-VISUAL/report.md 및 FIX-VISUAL-RESULT-REVIEW/report.md(check11/26통과). 기존 전체 UI 독립QA SAR-SERVICE-UI-001-TESTER-test/report.md와 원시각보고도 연결한다. 실제 Google·운영은 아직 미검증이다.
배포 대상은 위95df83c만이며 후속 coor 운영 준비커밋은 제품동일 metadata다. 원격 source SHA가 다르면 실제 제품 동일성·clean·승인 main조상을 확인하고 보고한다. Node22 npm/build/package verify는 기존 절차대로. 원래 Bot 설치의 실제 webhook 왕복은 LIVE/evidence.json에 있고 Bot watcher1 PID88404/private webhook 설정을 coor가 유지한다. OPS가 Bot/MCP를 재설치하거나 추가 watcher를 띄우지 않는다.
운영 전 fresh baseline·전후 보존 비교·실제 명령 원 종료코드를 기록한다. CSS만 바뀐 검수재실행/DB통합은 이전SHA결과와 구분한다. 실제 로그인·회원 목록·모바일 렌더·키보드 포커스 확인은 배포후 coor에게 정확한 운영SHA와 URL을 돌려준다. 기존 Google/Tunnel 정책을 바꾸지 않는다.

## beta.sh의 구체 충돌 — 먼저 확인
coor가 현재 beta.sh를 직접 읽었다. `prepare`는 state와새비밀 생성, `deploy`는 backup을 호출하며 `backup`은 raw pg_dump와 docker run --rm·실패시rm을 포함한다. `restore-verify`는 docker rm/trap을 포함한다. 이 과제에서는 이 entrypoint들을 실행하지 않는다. raw SQL dump·원문 state 수집·파일/컨테이너 정리 삭제를 요청하지 않았다. 기존 read-only baseline/집계·digest와 migration diff0 확인 후 이전 OPS의 bundle전달→검증된detached checkout→beta.sh up(relay만 갱신) 절차를 사용한다. up에 따른 승인된 relay 서비스 갱신 외 기존/임시 컨테이너나 파일을 삭제·정리하지 않는다. state/DB 복원·seed·owner-login/자격출력·Tunnel create/render/expose도 금지다. 정상 배포의 코드복귀는 이전eac603b로만하며 DB를 덮지 않는다. 실제기준 불일치·migration변경·설명되지않는검증실패는 coor에 ask한다.
