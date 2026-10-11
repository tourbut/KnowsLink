---
title: SAR-SERVICE-UI-001-FIX-VISUAL — V-01 글자 링크44px 수정 직접 재측정·시각 수락
status: draft
updated: 2026-10-11
owner: designer
tasks: [SAR-SERVICE-UI-001-FIX-VISUAL]
summary: V-01 글자 링크44px 수정 직접 재측정·시각 수락
attempt: 6e41d168df4f466e8401cf47727973b0
base: 2e8b8b228568baa2d849831a0859063de55aa473206009cd87129ab2faf2c5001ba058b1d
subagent_level: standard
test_level: lite
---
# SAR-SERVICE-UI-001-FIX-VISUAL — V-01 글자 링크44px 직접 재측정·제품 시각 수락

- 상태 queued: DEV 최종 고정 SHA 및 실제 구현자 세션 확인 후 coor가 기입해 dispatch한다. 지금 임의 후보를 수락하지 않는다.
- Purpose: product visual acceptance. Test level lite; Subagent level standard. 등록 Claude Sonnet5.5 high fresh. Codex Trust 미승인으로 대체하며 임의 신뢰 승인 금지.
- 담당 C:/Users/shin/orca/workspaces/KnowsLink/designer, fullops/designer. 복귀 run_7176a2ad6ef1/coor term_e61d3e14-29e9-4954-943a-4a75707c82de, task/dispatch는 실제 preamble.
- 공통 기준 fullops-common-0.3.3/project.md/FULLOPS.md/document-writing, 기존 D04 요구 유지. 기준 a072630 및 후속 준비 커밋. Ponytail full.

## 먼저 읽을 문서
- .fullops-squad/FULLOPS.md / project.md / contexts/designer.md
- .fullops-squad/rules/common/README.md 및 coding-style.md/testing.md/security.md
- .fullops-squad/docs/agents/document-writing.md
- .fullops-squad/docs/design-docs/mockups/service-member-ui.md
- .fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md
- .fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-review/report.md
- .fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-test/report.md
- .fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-DEV.md (완료 후 추가)

## 현재 상황과 범위
기존 후보 c7828b2는 독립 QA/코드리뷰와 시각PASS를 받았으나 글자 링크11종 높이21px가 D04 9절44px 요구 미달이다. 기존 결과1188f569의 V-01을 실제 수정한 새 고정 후보만 재검수한다. 기존 low11/F1-F3과 미검증·과거 실패는 보존한다. 동일 제품 바이트 부분은 이전 SHA/근거 연결로 재사용하며 전체 QA/캡처를 복제하지 않는다.

## 승인·제약
실제 Linux .env.server 서버 사용 가능, 로컬 Docker·추가 유료서비스·삭제/이동/광범위정리·DB/state/key/agent/pair/credential 변경 금지. helper 자격은 메모리 처리하고 비밀값 출력 금지. 웹은 Orca 내장 브라우저만. Grok Bot 수신 UI 접근 금지. 합성 fixtures만 새 고유 실행 폴더에 생성하고 리뷰 snapshot은 read-only로 유지한다. 기존 사용자 세션과 증거를 보존한다.

## 소유권과 완료 기준
- [ ] frontend-design 스킬과 기존 D04·시각 보고 V-01을 읽고 수정된 글자 링크11종만 Orca 내장 브라우저에서320/390/1440 직접 측정한다. 링크44px이상, 줄바꿈·overflow·인접 클릭영역·기존 버튼/skip-link에 시각 부작용 없음 확인. 필요한 대표 PNG만 저장·manifest hash, base64 출력 금지.
- [ ] 제품 시각 판정과 직접 픽셀/DOM측정/DEV보고를 구분한다. V-01 resolved는 실제 근거 있을 때만 후속 기록으로 연결하고 기존 발견을 삭제하지 않는다. 실제 Google/운영/스크린리더/200%는 하지 않은 경우 미검증 유지.
- [ ] 제품 코드 수정 금지. 이 QA/시각보고·필요시 D04 후속 수락 기록·phase/contexts만 소유. 별도 디자인lint/다크테마 미구성으로 해당없음이며 실제측정으로 보완한다.

## 결과와 후속
route D ID 없음. packet partial/unknown은 실제 필요한 영향 후보를 수동확인하여 outcomes에 적는다. 하위 위임은 작은 검수로 불필요하다. registered lint를 고정 결과에서 실행하고 unavailable은 이유를 남기며PASS하지 않는다. 완료 보고 전문·work.py finish·아카이브/빈인박스·커밋/push·worker_done까지 진행한다. 첫 완료문자열 [완료] SAR-SERVICE-UI-001-FIX-VISUAL | 최종 SHA <실제40자리>, 후보 SHA도 명시한다. main 통합/배포는 coor.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-SERVICE-UI-001-FIX-VISUAL-packet.json` / SHA `295417acfdcd415ad0c768d164bd71dc07214ea9` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 7, 7, 406, 406 · inferred · 필수 · {'relevant': 0.27, 'evidence': 0.39, 'contradicts': 0.26, 'injection': 0.06, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-review/report.md` (document_read) · 줄 2, 6, 10, 15, 36, 54, 58, 67, 68 · inferred · 필수 · {'relevant': 0.5, 'evidence': 0.74, 'contradicts': 0.5, 'injection': 0.06, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md` (document_read) · 줄 6, 6, 10, 10, 34 · inferred · 필수 · {'relevant': 0.67, 'evidence': 0.77, 'contradicts': 0.61, 'injection': 0.08, 'decision': 'conflict', 'reason': 'contradicts task'}
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-DEV.md` (document_read) · 줄 6, 10, 15, 15, 23, 31, 41, 47 · inferred · 필수 · {'relevant': 0.74, 'evidence': 0.84, 'contradicts': 0.35, 'injection': 0.07, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 6, 6, 13, 13, 25, 25, 26, 27, 28, 42, 42 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.5, 'evidence': 0.78, 'contradicts': 0.19, 'injection': 0.06, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 6, 30, 32, 68, 70, 72, 74, 76, 76 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 28, 38 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 6, 71 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 119 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 141 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 7 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 35 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md` (document_read) · 줄 21 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 6, 14, 31, 41 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 7, 116, 221, 227, 228, 232, 234, 234 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/COOR/dev-fix-2-final/coordinator-handover-supplement.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 16, 17 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW-review/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-REVIEW-review/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-SERVICE-001-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-RESULT-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-test/report.md` (document_read) · 줄 2, 6, 10, 25, 26 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-2-review/report.md` (document_read) · 줄 6, 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-review/report.md` (document_read) · 줄 6, 6, 10, 10, 12, 12 · inferred
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
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 6, 10, 24, 33, 67, 92, 104, 105, 108 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-TESTER.md` (document_read) · 줄 2, 6, 10, 20, 26, 32 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-VISUAL.md` (document_read) · 줄 2, 2, 6, 6, 10, 10, 21, 25, 25 · inferred
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
- `.fullops-squad/handovers/logs/2026-10-05_to_designer.md` (document_read) · 줄 9, 16, 20, 79, 87, 114 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_dev.md` (document_read) · 줄 16, 33, 35, 37, 51, 93, 125, 181 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_ops.md` (document_read) · 줄 90, 93, 97, 101, 116, 130, 144, 147, 151, 155, 171 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_tester.md` (document_read) · 줄 16, 33, 35, 37, 49 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 35, 102, 105, 106, 142, 248, 249, 306, 386 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 37, 176, 245, 246, 304, 357 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_tester.md` (document_read) · 줄 100, 175, 176, 228 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_designer.md` (document_read) · 줄 29, 69, 76, 133, 182, 183, 184 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 24, 108, 113, 148, 149, 229, 230 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_tester.md` (document_read) · 줄 48, 70, 75, 139, 180, 181, 182 · inferred
- `.fullops-squad/handovers/logs/2026-10-11_to_ops.md` (document_read) · 줄 124, 155, 186, 191, 204, 205, 210, 216 · inferred
- `.fullops-squad/handovers/to_ops.md` (document_read) · 줄 2, 6, 14, 33, 38, 64, 122, 123, 131, 132, 136, 192, 193, 194, 195, 196, 212 · inferred
- `.fullops-squad/handovers/to_tester.md` (document_read) · 줄 2, 6, 13, 25, 25, 26, 27, 28, 42, 47, 50, 51, 51, 69, 118, 119, 120, 121, 122, 122, 123, 123, 129, 169, 170, 171, 172, 173, 174, 175, 175, 187 · inferred
- `internal/relay/member_pages.go` (direct_edit, impact_check) · 줄 60 · inferred · {'relevant': 0.39, 'evidence': 0.61, 'contradicts': 0.37, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `internal/relay/member_ui_test.go` (direct_edit, impact_check) · 줄 21, 96, 120, 161 · inferred · {'relevant': 0.12, 'evidence': 0.17, 'contradicts': 0.31, 'injection': 0.06, 'decision': 'keep', 'reason': None}
- `internal/relay/remote_auth.go` (impact_check) · 줄 전체/미확인 · unknown · {'relevant': 0.12, 'evidence': 0.13, 'contradicts': 0.5, 'injection': 0.05, 'decision': 'keep', 'reason': None}
미확인 12건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고
착수 후 체크리스트를 갱신하고 완료 기준 전체를 채운 뒤 결과·한계·재개 조건을 쓴다.

## 지시 전제와 충돌 — 먼저 확인
Jev context는 member_pages.go/member_ui_test.go/remote_auth.go에 일반 conflict 수치를 반환했으며 구체적 reason은 없다. 현재 파일은 수정 전 후보의21px 링크와 버튼만44px 측정한 이전 테스트를 포함한다. 이는 V-01 발견 근거이며 미달을 수락 기준 완화로 해결하지 않는다. remote_auth.go는 새 UI 마크업 범위 밖이나 공통 스타일에 영향받는 원격 동의 글자 링크도44px 측정 대상이다. 새 최종 후보 diff와 D04 9절·V-01 원문을 먼저 대조한다. 과거 PASS는 V-01 잔여를 포함한 시각 판정이며 수정 후 실제 재측정이 필요하다. 실제 제품 규칙 충돌은 coor에 ask하고 임의 정본 변경하지 않는다.

## 최종 고정 후보 — 착수 정본
Review SHA: 295417acfdcd415ad0c768d164bd71dc07214ea9. Implementer session: 98481ab6-9f92-4bdd-ba02-c9881c19312c (ctx_4567096a3f21 실제 provider transcript 파일명 확인). Reviewer session은 본인의 실제 새 세션 ID를 사용한다. DEV 결과1e73308a2b3ebb5aeb14eb484651318deae4117c 및 제품609fcfd9b5049f093cec4162a3a83e4801d2e2df 포함. 후보295417a와 제품609fcfd 사이 제품 바이트 동일, 나머지는 coor 검수 준비·기록이다. review.py from은2e8b8b228568baa2d849831a0859063de55aa473, to는295417ac... 고정한다. 합성 fixtures는 이 후보에서 새 고유 실행 공간에 생성한다.
DEV Windows lint의 DOC-003 1건은 coor 최초 중단 리뷰 report의 front matter 누락이었다. 4a65bac에서 원본 템플릿을 보존하고 중단·superseded 메타를 붙여 수정했다. 후보295에 포함되어 있으므로 최신 lint에서 오류0 확인한다. 과거 DEV lint 실패를 PASS로 바꾸지 않는다. Linux make gate는609fcfd의 실제 종료0이며 제품 동일성을 대조해 재사용한다. 새 DBintegration 미재실행은 CSS만 변경한 영향 경계로 원래 근거를 유지한다.
Codex hook Trust 미승인으로 등록 Claude Sonnet5.5를 사용한다. tester medium/시각 high는 기존 독립 QA와 직접 시각 책임을 유지한다. queued 조건은 이 절의 고정 후보·실제 구현자 확인으로 충족됐다. 다른 보류는 main 통합 후 coor가 처리한다.
