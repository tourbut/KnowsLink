---
title: SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL — tester의 같은 후보 캡처로 OAuth 반환·오류 화면 직접 시각 수락
status: draft
updated: 2026-10-11
owner: designer
tasks: [SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL]
summary: tester의 같은 후보 캡처로 OAuth 반환·오류 화면 직접 시각 수락
attempt: b61277b0559d44d992788a7d19d490cd
base: eac603b6a41d2ce002e698600ac66b8fc67ad762
subagent_level: standard
test_level: lite
---
# SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL — OAuth 반환·오류 화면 직접 시각 수락

- Purpose visual acceptance. 상태: tester 합성 캡처와 독립 QA 전달 완료, ready.
- designer C:/Users/shin/orca/workspaces/KnowsLink/designer, fullops/designer. fresh Claude Sonnet5.5 high(운영 대체). test lite/subagent standard. Jev Codex high는 실제 hook 신뢰 검토창에서 readiness 실패하여 작업 미착수였다. 기존 검증된 Claude hook 아래 같은 역할·범위로 재배정한다. Codex trust 값을 수정하지 않는다.
- 고정 제품 SHA eac603b6a41d2ce002e698600ac66b8fc67ad762. 실제 구현자 e20cc179-625e-43fa-84fd-b431d432b9c6. 소스·제품 기준을 바꾸지 않는다.
- 복귀 coor term_e61d3e14-29e9-4954-943a-4a75707c82de / run_dc9b4e6dd18a. task/dispatch는 preamble을 따른다.

## 요구와 검수 범위

OAuth consent POST303이 Chromium CSP에 막히는 오류를 DEV가 반환 화면으로 고쳤다. 기능 QA/fixed-SHA code review는 tester 담당이다. designer는 같은 통합 후보의 tester 캡처를 직접 보고 반환·재시작 안내가 사용 가능한지 판정한다. 새로운 기획·스타일·화면을 만들지 않는다.

기존 memberStyle/refusal 스타일과 한국어 안내를 유지한다. 별도 마케팅 화면·장식·새 CSS·테마·애니메이션을 추가하지 않는다. 반환 링크가 보이고 클릭 대상임을 알 수 있는지, desktop/mobile에 글자 넘침·가림이 없는지, expired/used 화면에 재시작 경로가 명확한지, 화면에 민감한 callback/code/state가 노출되지 않는지 확인한다. 자동 복귀·토큰 교환 성공은 tester 근거를 연결한다.

## 먼저 읽을 문서·적용 기준

지정 캡처·SHA256는 `docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-COOR/visual-inputs.json`에 있다. 반환 desktop/mobile, expired desktop, CSRF/owner-binding mobile의 5장을 직접 연다. 원 근거는 `SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-test/report.md`와 같은 폴더 screenshots(총19장)다. tester 기록 SHA는 96c7b8236f8335b4e0c7344e09c00e3cb36c0e8b다. F1 medium은 오류 원인별 문구가 같다는 비차단 잔여로 남긴다. F2 low의 링크24px은 실제 캡처와 클릭 성공을 함께 보고 시각 수락 가능 여부를 판단한다.

제품 gate와 tester 최종 기록 gate는 같은 COOR 폴더의 JSON에 보존됐다. `tester-final-product-shim.py`는 tester가 실제 실행한 최종 shim 보존본이다. 자신의 clean HEAD와 고정 eac603b6의 제품 바이트 전체(.fullops-squad 제외)가 같은지 검사한 뒤, 기존 격리 Linux의 같은 eac603b6 clean clone에서 등록 make target만 실행한다. Windows make 부재 시 자기 task-local make.cmd에서 이 보존본을 호출해 사용할 수 있다. 제품 바이트와 lint 정책이 다르면 재사용하지 않는다. 원문을 고치거나 exact Git HEAD 실행이라고 오인하지 않는다. 검사를 실행하면 byte-identical product 검증임을 결과에 적는다.

FULLOPS, project, orca-agents, rules/common README+세 규칙, delegation, document-writing, interface-design, DEV OAuth 복귀 실행 기록을 읽는다. fullops-common-0.3.3, test lite를 적용한다. 제품 정본은 고정 SHA의 remote_auth.go와 기존 memberStyle/refusal이며 캡처 경로·hash는 tester 완료 후 연결한다. 기존 시각 PASS를 새 화면의 수락으로 해석하지 않는다.

## 소유권과 완료 기준

직접 이미지 도구로 desktop/mobile 지정 캡처를 확인하고 본인의 exec plan·QA 보고·contexts/designer만 작성한다. 제품 source·다른 역할 증거·운영/실계정/Google·Grok 조작·삭제/이동/cleanup은 금지한다. 부족한 캡처는 coor에 구체적으로 요청한다. 전체 QA와 이미 완료된 자동 검사를 다시 실행하지 않는다. 영상은 필요 없다.

관찰별 PASS/FAIL/미확인과 실제 대상 SHA·이미지 hash·한계를 기록한다. 이미지가 새 코드에서 실제 생성됐는지 tester manifest와 source를 대조한다. 문서 결과의 등록 lint를 실행하고 Windows make 부재는 원문을 보존해 사유를 적는다. 제품 바이트가 같으므로 원 exact-SHA product gate 재사용을 명시하며 통과를 꾸미지 않는다. packet의 모든 비선택 path/category를 outcomes로 처리하고 work.py finish로 원 지시서·보고를 보존한다. 최종40자리 기록 SHA 하나를 worker_done에 보고한다. main 통합·배포는 부모가 맡는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL-packet.json` / SHA `ab924bab38bcbcbbcea5dd43ece4c996c4f48a70` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-test/report.md` (document_read) · 줄 2, 2, 2, 6, 6, 6, 7, 10, 10, 10, 72 · inferred · 필수
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 2, 2, 6, 6, 6, 6, 7, 13, 13, 13, 13, 22, 28, 28, 28, 32 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 6 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `internal/relay/member.go` (impact_check) · 줄 전체/미확인 · unknown · 필수
- `internal/relay/remote_auth.go` (impact_check) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/PLANS.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 70 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 6, 6, 28, 28, 28, 28, 29, 30 · inferred
- `.fullops-squad/contexts/ops.md` (document_read) · 줄 6, 34, 67 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 6, 6, 6, 66, 66, 66 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 7, 14, 20, 23, 25, 155, 163, 179 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read) · 줄 7, 14, 38 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 7, 14, 19 · inferred
- `.fullops-squad/docs/design-docs/database-design.md` (document_read) · 줄 7, 14, 39 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 7, 7, 14, 23, 28, 190 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 87 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 82 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 7, 7, 7, 14, 24, 55, 188, 199 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 7, 14, 16, 19, 21, 21 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-001-TESTER-test/report.md` (document_read) · 줄 43, 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-LOGIN-001-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-TESTER.md` (document_read) · 줄 2, 6, 10, 23, 27, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 130 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-review/report.md` (document_read) · 줄 2, 6, 10, 17, 51 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 93 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-COOR/visual-inputs.json` (direct_edit) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-review/report.md` (document_read) · 줄 2, 2, 2, 6, 6, 6, 7, 10, 10, 10, 15, 15, 15, 20, 20, 20, 41, 41, 41, 43, 43, 43 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OPS-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-REVIEW-review/report.md` (document_read) · 줄 2, 6, 10, 15, 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-TESTER-test/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-VISUAL/report.md` (document_read) · 줄 6, 6, 10, 10, 72, 72 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER.md` (document_read) · 줄 2, 2, 2, 2, 6, 6, 6, 10, 10, 10, 12, 12, 12, 20, 20, 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.2.0.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 99 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-001-DEV.md` (document_read) · 줄 24, 34 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-002-DEV.md` (document_read) · 줄 17 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md` (document_read) · 줄 27, 33, 35, 39 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md` (document_read) · 줄 18, 61 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md` (document_read) · 줄 2, 6, 10, 17, 55, 56, 57, 59, 60, 61 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md` (document_read) · 줄 68, 72, 93 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 58 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (document_read) · 줄 36, 56 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 24, 53, 204, 221 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-DEV.md` (document_read) · 줄 6, 10, 38, 48, 55, 76, 82, 82, 86, 92, 97 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-INSTALL.md` (document_read) · 줄 2, 6, 10, 29, 55 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV.md` (document_read) · 줄 2, 6, 6, 6, 10, 10, 10, 60, 60, 60, 60 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER.md` (document_read) · 줄 2, 2, 2, 6, 6, 6, 7, 10, 10, 10, 28, 28, 28 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OPS.md` (document_read) · 줄 2, 6, 10, 21, 22, 22, 24, 37 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-REMOTE.md` (document_read) · 줄 2, 2, 6, 10, 16, 22, 27, 44, 48, 57, 72, 76 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_read, document_update) · 줄 7, 7, 7, 14, 16, 20, 22, 24, 40, 42, 44, 44, 44, 47, 51, 280, 314, 340 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read, document_update) · 줄 7, 13, 15, 15, 17, 19, 21, 23, 28, 51, 93, 103 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 7, 13, 25 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 136 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_ops.md` (document_read) · 줄 191 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_dev.md` (document_read) · 줄 9, 12, 16, 20, 33, 39, 71, 104, 126 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_ops.md` (document_read) · 줄 9, 12, 16, 20, 29, 34, 38, 41, 52, 54 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_tester.md` (document_read) · 줄 9, 12, 16, 20, 29, 34, 43, 51, 57 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_dev.md` (document_read) · 줄 60, 89 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 76 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 210 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_designer.md` (document_read) · 줄 113, 157 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 115, 179, 210, 307 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_tester.md` (document_read) · 줄 116, 158 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_designer.md` (document_read) · 줄 9, 9, 10, 10, 13, 13, 13, 17, 17, 18, 25, 25, 25, 34, 40, 49, 49, 69, 69, 77, 77, 81, 81, 143, 143 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_ops.md` (document_read) · 줄 9, 10, 14, 18, 26, 35, 54, 61, 80, 83, 85, 87, 91, 93 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_tester.md` (document_read) · 줄 9, 10, 13, 17, 25, 36, 36, 46, 48, 71, 83, 97, 103, 175, 175, 175, 176, 176, 176, 179, 179, 179, 179, 183, 183, 183, 184, 191, 191, 191, 191, 201, 201, 201, 207, 228, 228, 228, 230, 230, 230, 231, 231, 231, 286, 287, 313, 314, 315, 316, 316, 342, 366, 367, 368, 369, 434, 434, 434, 440, 440, 440, 445, 445, 445 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-DEV-REVIEW.md` (document_read) · 줄 45 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-INTEGRATION-REVIEW.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/handovers/to_ops.md` (document_read) · 줄 2, 2, 2, 2, 6, 6, 6, 7, 13, 13, 13, 13, 21, 27, 27, 39, 47 · inferred
- `adapters/README.md` (document_read) · 줄 6, 14, 24, 50 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 14 · inferred
- `scripts/package_plugin.py` (impact_check) · 줄 146 · inferred
미확인 11건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

최종 기록과 남은 제한을 전문으로 보존한다.
