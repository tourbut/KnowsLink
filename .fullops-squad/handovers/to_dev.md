---
title: SAR-SERVICE-UI-001-FIX-DEV — D04 글자 링크 44px 터치영역 V-01 최소 수정
status: draft
updated: 2026-10-11
owner: dev
tasks: [SAR-SERVICE-UI-001-FIX-DEV]
summary: D04 글자 링크 44px 터치영역 V-01 최소 수정
attempt: 63b11e3217af4852a1abe19766e7f2bd
base: c2fd5f4c37f17c3833ea2e5d9134feeb01561257
subagent_level: standard
test_level: lite
---
# SAR-SERVICE-UI-001-FIX-DEV — D04 글자 링크 44px 터치영역 V-01 최소 수정

- From / To: coor / dev. 상태 ready. Purpose: implementation. Test level lite, Subagent level standard.
- 워크트리 C:/Users/shin/orca/workspaces/KnowsLink/dev, branch fullops/dev. 병합 책임 coor/main.
- 복귀 run run_7176a2ad6ef1, terminal term_e61d3e14-29e9-4954-943a-4a75707c82de. task/dispatch는 실제 worker-start 영수증을 사용한다.
- 사용자 승인: 상용 UI 기능 보완·무과금·기존 Linux 서버 테스트·상태 보존. 새 유료 서비스/로컬 Docker/삭제/force/전역설정 변경 금지. hook Trust를 대신 승인하지 않는다.

## 현재 상황과 기준
고정 후보 c7828b227b5715f7bdf02bff46bf9cf102a758d9는 독립 기능 QA와 시각 PASS를 받았다. designer 고정 결과1188f569111154c2e92103acbb4ee832ed6642dd의 V-01 medium은 글자 링크11종의 높이21px가 D04 9절44px 기준에 미달함을 확인했다. coor는 이 결함을 배포 전에 최소 수정하기로 했다. 원래 실패/한계/기존 리뷰 기록은 보존한다. 기존 low와 새 기능은 이 과제 범위 밖이다.

## 적용 기준과 예외
fullops-common-0.3.3 및 project.md, D04 service-member-ui.md. 기준 c2fd5f4c37f17c3833ea2e5d9134feeb01561257와 준비 커밋에 통합된 시각 결과를 읽는다. 원 제품 인증·폼·CSRF·owner·recent5분·CSP 계약을 유지한다. 라이브러리/API 변경 없음, 기존 CSS/Go template 사용. Codex Trust 미승인으로 등록 Claude 후보를 사용한다.

## 먼저 읽을 문서
- .fullops-squad/FULLOPS.md / project.md / contexts/dev.md
- .fullops-squad/rules/common/README.md 및 coding-style.md / testing.md / security.md
- .fullops-squad/docs/agents/document-writing.md
- .fullops-squad/docs/design-docs/mockups/service-member-ui.md
- .fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md (특히 V-01)
- .fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-review/report.md
- .fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md (Linux gate/fixture 재사용)
- internal/relay/member_pages.go 및 관련 UI template/fixture/test 호출자

## 해야 할 일과 파일 소유권
- [ ] 관련 모든 링크 호출자/공통 CSS를 확인하고 짧은 기술 계획을 phase에 기록한다.
- [ ] V-01 글자 링크의 실제 터치 영역이44px 이상이 되도록 공통 스타일에서 최소 수정한다. 작은 화면 줄바꿈·기존 버튼/inline코드/skip-link에 부작용이 없게 한다.
- [ ] 관련 테스트·fixtures와 기술 기록을 필요한 만큼만 갱신한다. 새 라이브러리/JS/색상/추가 기능은 넣지 않는다. 예상 제품 변경은 CSS와 필요시 회귀 증거의 수십 줄이다.
- [ ] 고정 최종 SHA에서 등록 lint/test와 관련 실패경계/반응형 측정 결과를 보존하고 완료 아카이브·packet outcomes·커밋/push·worker_done까지 처리한다.

## 완료 기준과 검증
DEV: 해당 글자 링크11종을 포함한 a[href] 측정을320/390/1440에서 수행해 높이44px 미만0·overflow0 확인. 기존 버튼/summary/키보드 동작·폼 golden 계약을 확인한다. 이전 테스트의 버튼만 측정한 누락을 재사용하지 않는다. 기존 Linux 보존 검증 도구를 사용하고 원 명령 종료코드를 남긴다. 기존 서버자료/서비스/컨테이너/임시증거도 삭제하지 않는다. make가 없는Windows에서 unavailable을PASS로 바꾸지 않는다.
독립 tester: 새 최종 SHA의 작은 diff 고정 snapshot 리뷰와 관련 link/overflow 회귀. designer: V-01 해당 링크만 직접 재측정한다. 이전 전체 검증은 제품 바이트 동일 부분의 SHA·근거를 연결하고 중복하지 않는다. 실제Google·운영렌더는 OPS후 coor 검증이다.
공통 memberStyle/token 정본을 재사용한다. 별도 디자인 lint/다크테마는 project.md 미구성으로 해당없음이며 실제 DOM 측정으로 보완한다. D04를 완화하거나 V-01을 단순 문서상 resolved 처리하지 않는다.

## 산출물과 후속
route 갱신할 D ID 없음. 구현/검증·phase·contexts/dev.md 및 해당 UI 기술 원천만 필요한 범위로 갱신한다. 정본 보고서의 과거 발견을 삭제하지 말고 후속 수정 근거를 연결한다. packet partial/unknown은 실제 영향 후보를 확인해 outcomes에 적는다. 하위 위임은 작은 직렬 CSS 수정이므로 불필요하다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-SERVICE-UI-001-FIX-DEV-packet.json` / SHA `5febcb355a175e89e6bc5ac87bfff42b6907e7db` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 7, 406 · inferred · 필수 · {'relevant': 0.21, 'evidence': 0.33, 'contradicts': 0.27, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-review/report.md` (document_read) · 줄 2, 6, 10, 15, 36, 54, 58, 67, 68 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md` (document_read) · 줄 6, 10, 34 · inferred · 필수 · {'relevant': 0.89, 'evidence': 0.91, 'contradicts': 0.14, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 6, 10, 24, 33, 67, 92, 105 · inferred · 필수 · {'relevant': 0.21, 'evidence': 0.59, 'contradicts': 0.43, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 6, 13, 31, 32, 33, 51 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.44, 'evidence': 0.79, 'contradicts': 0.35, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 30, 32, 70, 72, 74, 76 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 28 · inferred · {'relevant': 0.24, 'evidence': 0.51, 'contradicts': 0.18, 'injection': 0.07, 'decision': 'keep', 'reason': None}
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
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-2-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-review/report.md` (document_read) · 줄 1 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-3.md` (document_read) · 줄 15 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (document_read) · 줄 6, 10, 26, 28, 32, 48, 49, 66 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPEN-PREP.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-TESTER.md` (document_read) · 줄 2, 6, 10, 20, 26, 32 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-VISUAL.md` (document_read) · 줄 2, 6, 10, 21, 25 · inferred
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
- `.fullops-squad/handovers/logs/2026-10-11_to_ops.md` (document_read) · 줄 204, 205, 216 · inferred
- `.fullops-squad/handovers/to_ops.md` (document_read) · 줄 2, 6, 14, 33, 38, 64, 122, 123, 192, 193, 194, 195, 196, 212 · inferred
- `internal/relay/member_pages.go` (direct_edit, impact_check) · 줄 60 · inferred · {'relevant': 0.68, 'evidence': 0.87, 'contradicts': 0.2, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `internal/relay/member_ui_test.go` (direct_edit, impact_check) · 줄 21, 96, 120, 161 · inferred · {'relevant': 0.07, 'evidence': 0.11, 'contradicts': 0.27, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `internal/relay/remote_auth.go` (impact_check) · 줄 전체/미확인 · unknown · {'relevant': 0.02, 'evidence': 0.03, 'contradicts': 0.5, 'injection': 0.03, 'decision': 'keep', 'reason': None}
미확인 11건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고
[완료] SAR-SERVICE-UI-001-FIX-DEV | 최종 SHA 08e97e2b26ba730a86eeadc6f2b9fe52da50efd5 (문서 기준. 이후 archive 커밋이 최종 HEAD이며 worker_done의 결과 SHA를 따른다)

- 제품 변경 SHA `609fcfd9b5049f093cec4162a3a83e4801d2e2df`: `internal/relay/member_pages.go`의 공통 `a` 규칙에 `display:inline-flex;align-items:center;min-height:44px`를 추가하고 중복 `.toc a` 규칙을 지웠다. `member_ui_test.go`는 fixture 29개 렌더마다 이 규칙을 확인한다. 마크업·form·CSRF·owner·recent 5분·CSP는 바꾸지 않았다.
- 링크 호출자: 모든 회원 템플릿, gate, OAuth 동의·반환 화면이 `memberStyle` 하나를 쓴다. 호출자별 수정은 없다.
- 검증(실제 종료코드): Windows `go test ./internal/relay -run TestMemberUIFixtures -count=1` 0, `gofmt`·`git diff --check` 0. Orca 내장 브라우저에서 fixture 29개 x 1440/390/320, details 펼침 상태로 `a[href]` 198건과 버튼·입력·select·summary를 측정했다: 44px 미만 0건, overflow 요소 0건, 종료코드 0. 기존 Linux 서버 고유 clone `/tmp/knowslink-uifix-609fcfd`(대상 SHA detach)에서 `make install`·`make lint`·`make test` 모두 0, clone clean. Windows에는 make가 없어 Linux 결과를 대체로 쓰지 않았다.
- 기록: 실행 기록 `docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-DEV.md`, 측정 `docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-DEV/`(measure_links.py, measurements.json), packet outcomes `docs/evaluations/jev/SAR-SERVICE-UI-001-FIX-DEV-packet-outcomes.json`, `contexts/dev.md`. 원 V-01 보고와 이전 DEV 기록은 보존하고 후속 연결만 추가했다.
- 미검증: Postgres integration·`make verify-mvp`는 이번 SHA에서 재실행하지 않았다(CSS 한 규칙, Go 로직 변경 없음). tester의 새 SHA 작은 diff 리뷰·link/overflow 회귀, designer의 V-01 링크 직접 재측정, 390px "취소" 버튼 아래 간격의 실제 모습, 실제 Google·운영 렌더(OPS 뒤 coor), main 통합·배포는 후속이다. 기존 low V-02–V-12는 범위 밖이다. 서버의 기존 자료·컨테이너·임시 증거는 삭제하지 않았다. Codex hook Trust는 승인하지 않았다.
