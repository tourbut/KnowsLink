---
title: SAR-LANDING-GUIDE-001-VISUAL — 고정 구현 후보의 공개 첫페이지를 실제 브라우저에서 시각 검수하고 D04 수락 여부를 기록한다
status: draft
updated: 2026-10-14
owner: designer
tasks: [SAR-LANDING-GUIDE-001-VISUAL]
summary: 고정 구현 후보의 공개 첫페이지를 실제 브라우저에서 시각 검수하고 D04 수락 여부를 기록한다
attempt: 93adb3b5b2124aed99b9c02f45e39281
base: fe2387a7676f7997aea76a2606f4cc2dbcd4b56b
subagent_level: standard
test_level: lite
---

# SAR-LANDING-GUIDE-001-VISUAL — 고정 구현 후보의 공개 첫페이지 직접 시각 검수

- 작성일: 2026-10-14 / From coor → designer / 상태: ready.
- 선행 충족: DEV worker_done msg_0a44968602b5의 최종 fe2387a7676f7997aea76a2606f4cc2dbcd4b56b를 coor 기록 브랜치에만 통합했다. runtime/fixture exact4795d619513fde4fbcd4a55b94ec91b056ebbaf0와 입력 diff0이다. Google/both/email/none는 http://127.0.0.1:18781/ ~18784/. DEV-test/report.md의 실제 capture manifest·원실패·fixture 허용 경로를 읽는다. 별도 isolated profile/page에서 관측하고 기존 exactSHA 이미지는 직접 열어 판정할 수 있다. 현재 새로운 dispatch를 기다린다.
- 실행: codex/gpt-6-astra/high, 사용자 최신 지정. test_level=lite/subagent_level=standard.
- 담당: repo0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer.
- 복귀: coor / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_38ae16dc-0ca4-4beb-ab00-d1e1b641f061 / run_5534969da2d0. 실제 새 task/dispatch preamble을 사용한다.
- 병합 책임: coor / main·origin/main. 제품 코드 변경·배포·실계정·실메일·외부 발송은 범위 밖이다.

## 적용 기준과 예외

AGENTS.md, FULLOPS.md, fullops-common-0.3.3 공통 코딩·테스트·보안, project.md, document-writing.md, fullops-frontend/delivery.md를 적용한다. frontend.json 기존 주 스킬 frontend-design과 D04 공개 첫페이지 정본을 사용한다. 고정 실행 tuple은 사용자 지정이 우선하며 새 도구 설치·유료 키·이미지 에셋은 추가하지 않는다. 과거 최소 로그인 절은 새 공개 첫페이지 정본에서 대체한 범위만 달라진다. 실제 QA와 미실행을 구분하고 기존 실패를 보존한다.

## 먼저 읽을 문서

- .fullops-squad/docs/design-docs/mockups/service-member-ui.md — 공개 첫페이지 정본.
- .fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md — LG-01~08.
- .fullops-squad/docs/agents/frontend.json.
- .fullops-squad/docs/operations/user-guide.md.
- .fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md.
- README.md — 실제 설치·수신 조건.
- internal/relay/landing_pages.go, internal/relay/member_pages.go — 고정 후보만 읽는다.
- .fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/report.md — 실제 DEV evidence와 한계.
- .fullops-squad/handovers/to_tester.md — 검사 담당 분리.

## 해야 할 일과 소유권

### 지시 전제와 충돌 — 먼저 확인

context는 landing_pages.go에 구체 문장 없이 contradicts task를 표시했다. 이 구현은 확정 D04에 대한 검수 대상이고 coor가 시각 PASS를 가정하지 않는다. 원문·실제 화면을 비교하고 제품 기준 차이는 증거로 보고한다. 큰 D04의 관련 절과 partial/unknown/잔여 후보는 직접 읽는다. frontend의 과거 Sol 기본값보다 최신 사용자 Astra high가 우선하며 기준을 바꾸지 않는다. 제외한 필수 문서는 없다.

- [ ] 실제 고정 후보의 loopback 페이지와 SHA를 대조하고 desktop1440/mobile390/320을 직접 캡처·열어 확인한다.
- [ ] hero·관계 설명·시작 순서·Grok/Dots 펼침·가상 예시·로그인/오류/none을 D04와 비교한다. 기존 회원 CSS에 영향을 주지 않았는지 대표 화면을 확인한다.
- [ ] 텍스트 위계·여백·한글 줄바꿈·긴 명령 overflow·44px 대상·focus/대비를 판정한다. 키보드/200% 실제 조작은 tester 결과와 연결하며 사진만으로 실행 PASS를 주장하지 않는다.
- [ ] 결함은 재현·viewport·증거·심각도와 수정 담당 DEV를 기록한다. 제품 코드를 수정하지 않는다.
- [ ] own phase/QA/manifest/context/packet outcomes를 작성하고 work.py finish 및 실제 worker_done으로 고정 기록 SHA를 회신한다.

파일 소유권은 own 인박스/context, docs/exec-plans/phases/SAR-LANDING-GUIDE-001-VISUAL.md 및 docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/와 own 탐색 기록이다. D04 변경·다른 역할 결과·제품 코드는 읽기 전용이다.

## 실행·완료 기준

coor 제공 고유 fixture origin에서 GET·앵커·native details 펼침과 viewport만 허용한다. 합성 오류는 제공 fixture의 테스트 경로만 사용한다. 실제 Google/운영/메일/설치 실행·계정 생성·발송은 금지한다. Orca CLI와 frontend 기존 도구를 사용하고 페이지별 browserPageId로 tester와 충돌을 피한다. 최대40개 의도적 조작, 20분/시나리오로 제한한다. capture 기존 evidence는 exactSHA·viewport·상태 확인 후 재사용하고 실제 이미지를 열어 판정한다. 빠진 상태만 새 캡처한다. 영상은 필요 없다.

보고서에 고정 제품 SHA, author와 다른 실제 session, 실행기·origin·page·viewport·zoom·파일 경로, 수락/FAIL/not_run, D04 대비 차이와 미해결 C/H를 기록한다. static wireframe·브라우저 설치·CLI 성공으로 실제 시각 PASS를 대신하지 않는다. 디자인 미해결 critical/high는 수락을 차단한다. 독립 tester 기능 QA·코드 리뷰·운영 배포 성공을 이번 시각 검수로 주장하지 않는다.

Windows 삭제·이동 정리는 실행하지 않는다. 검증 자료를 보존하고 경로·fixture 유지 담당을 남긴다. PLANS/board와 main 통합은 coor 담당이다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-LANDING-GUIDE-001-VISUAL-packet.json` / SHA `1312658557c77e21de852edee6a860ea8cc69988` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/report.md` (document_read) · 줄 6, 6, 12, 12 · inferred · 필수
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 7, 7, 13, 13 · inferred · 필수
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 2, 6, 6, 6, 14, 14, 14, 33, 33, 36, 36, 51, 51, 51 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `internal/relay/landing_fixture_test.go` (impact_check) · 줄 전체/미확인 · inferred · 필수
- `internal/relay/landing_pages.go` (impact_check) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/archive/designer-2026-10-14-landing.md` (document_read) · 줄 6, 6, 6, 72, 80, 84, 88, 92, 94, 96, 98, 98, 100, 100 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 6, 6, 22, 24, 26, 26, 28, 28 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 52, 52 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 7, 14, 14 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 7, 7, 14, 14, 20, 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-DEV-test/report.md` (document_read) · 줄 33, 48 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-TESTER-review/report.md` (document_read) · 줄 65 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-VISUAL-DOCS-review/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-VISUAL/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL-DOC-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL-FINAL-DOC-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/final-report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6, 6, 10, 10, 56, 56 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-test/report.md` (document_read) · 줄 6, 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-TESTER-review/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL-RESULT-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL/report.md` (document_read) · 줄 6, 10, 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-LIVE/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-2-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL-RESULT-REVIEW-review/report.md` (document_read) · 줄 6, 10, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md` (document_read) · 줄 6, 10, 179, 181 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL-REVIEW-review/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL/report.md` (document_read) · 줄 6, 10, 77 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-VISUAL/report.md` (document_read) · 줄 6, 10, 72 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-FRONTEND-POLISH-001-DEV.md` (document_read) · 줄 31 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-FRONTEND-POLISH-001-VISUAL.md` (document_read) · 줄 6, 10, 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-OPS.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-VISUAL.md` (document_read) · 줄 2, 6, 10, 24, 38 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001-DEV.md` (document_read) · 줄 6, 6, 43, 43 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md` (document_read) · 줄 6, 6, 10, 10, 50, 50, 51, 51, 52, 52, 53, 53, 57, 57, 59, 59 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 104 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-DEV.md` (document_read) · 줄 15 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-FIX-VISUAL.md` (document_read) · 줄 2, 6, 10, 22 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-OPS.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-VISUAL.md` (document_read) · 줄 2, 6, 10, 25 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL.md` (document_read) · 줄 2, 6, 10, 21 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OPS.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 7, 7, 13, 13, 15, 15 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_dev.md` (document_read) · 줄 102, 112, 128, 129, 136, 137, 162, 183 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_ops.md` (document_read) · 줄 166, 170 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_tester.md` (document_read) · 줄 216, 221, 227 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_designer.md` (document_read) · 줄 9, 9, 10, 10, 13, 13, 17, 17, 25, 25, 67, 67, 101, 101, 106, 106, 136, 137, 138, 139, 144, 188, 188, 189, 189, 190, 190 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_dev.md` (document_read) · 줄 211, 212, 249, 250, 251, 252, 258, 271, 288, 290, 357, 357, 358, 358, 361, 361, 365, 365, 373, 373, 384, 384, 386, 386, 393, 393, 394, 394, 429, 429, 448, 448, 467, 468, 469, 469, 470, 470, 471, 473, 476, 476, 479, 514, 514 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_tester.md` (document_read) · 줄 182 · inferred
- `.fullops-squad/handovers/to_tester.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 25, 25, 27, 27, 34, 34, 35, 35, 85, 85, 92, 92, 95, 95, 141, 141, 142, 142, 194, 218, 218, 219, 219 · inferred
- `internal/relay/landing_test.go` (impact_check) · 줄 39 · inferred
미확인 18건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

미착수. DEV 고정 최종 SHA·fixture와 실제 dispatch를 기다린다.
