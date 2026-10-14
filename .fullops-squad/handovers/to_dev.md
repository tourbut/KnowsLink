---
title: SAR-LANDING-GUIDE-001-DEV — 독립 QA의 기존 대화 CTA 대비와 공개 예시 하단 간격을 최소 CSS 변경으로 수정하고 실제 영향 검사를 완료한다
status: draft
updated: 2026-10-14
owner: dev
tasks: [SAR-LANDING-GUIDE-001-DEV]
summary: 독립 QA의 기존 대화 CTA 대비와 공개 예시 하단 간격을 최소 CSS 변경으로 수정하고 실제 영향 검사를 완료한다
attempt: 9a48a4180041493b93ddce9b5d1bf6f2
base: fe2387a7676f7997aea76a2606f4cc2dbcd4b56b
subagent_level: standard
test_level: lite
---

# SAR-LANDING-GUIDE-001-DEV — 독립 QA의 두 CSS 결함 수정

- 작성일: 2026-10-14 / From coor → dev / 상태: ready.
- 실행: codex/gpt-6-astra/high, 사용자 지정. test_level=lite/subagent_level=standard.
- 기준: fe2387a7676f7997aea76a2606f4cc2dbcd4b56b. 최초 attempt/원 실패/4795 fixture·code review·QA는 보존하고 이번 reopen attempt에 연결한다.
- 담당: repo0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev.
- 복귀: coor / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_38ae16dc-0ca4-4beb-ab00-d1e1b641f061 / run_5534969da2d0 / 실제 새 task·dispatch.
- 병합 책임: coor / main·origin/main. 운영·실계정·실메일·실상대 송신·삭제/이동 정리·새 dependency는 범위 밖이다.

## 적용 기준과 근거

AGENTS/FULLOPS·common0.3.3·project.md·문서 규칙을 유지한다. fullops-frontend 주 스킬 frontend-design, D04 공개 첫페이지와 기존 여섯 토큰·44px·대비 기준을 지킨다. 제품 규칙/문구/인증 handler를 다시 설계하지 않는다. 최소 CSS로 기존 대비와 예시 간격을 고친다. 기술 원인·구현 방법은 DEV가 결정한다.

tester escalation msg_18106237290b: 기존 .conversation-list .primary의 실제 색24,43,49/0,90,98·16px 대비1.846967로4.5 실패다. 공개 hero primary는7.958060 PASS이며 기존 공용 CSS 우선순위 결함이다. 증거는 C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-TESTER-test/contrast-results.json 및 member-home-1440.png다.
designer V-LG-01 low: 예시 블록과 그 다음 '대화에서 알아둘 점' 제목 사이 간격0px, 1440/390/320에서 관측했다. 제목/본문은 읽히지만 독립 구획을 구분하는 기존 D04 간격 토큰 적용이 필요하다. 원 보고는 C:/Users/shin/orca/workspaces/KnowsLink/designer/.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/report.md의 구체적 발견 절이다. 다른 역할 파일은 읽기만 한다.

## 먼저 읽을 문서

- internal/relay/member_pages.go, internal/relay/landing_pages.go.
- internal/relay/member_ui_test.go, internal/relay/human_pages_test.go, internal/relay/landing_test.go.
- .fullops-squad/docs/design-docs/mockups/service-member-ui.md.
- .fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/report.md.
- .fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CODE-REVIEW-review/report.md.
- .fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001-DEV.md.
- .fullops-squad/docs/agents/frontend.json.
- 위 tester/designer 원 증거 경로.

## 해야 할 일과 완료 기준

### 지시 전제와 충돌 — 먼저 확인

context는 최초 DEV QA를 conflict로 표시했다. 그 공개 CTA 대비 PASS는 이번에 발견한 회원 conversation-list CTA의 대비를 검사한 것이 아니다. 옛 make/공개 검사 PASS와 원 실패는 유지하고 새 회원 medium·예시 low의 범위를 추가한다. 옛 native keyboard not_run은 독립 tester의 실제 Chrome CDP 결과와 구분하고, native browser200%는 아직 별도 후속이다. 과거 Sol 기본보다 최신 사용자 Astra high가 우선한다. 큰 D04·partial/unknown/잔여 후보는 관련 원문을 직접 확인하며 제외한 필수 문서는 없다.

- [ ] selector·호출되는 화면·기존 hover/focus/44px를 확인하고 짧은 원인/영향을 기록한다.
- [ ] 두 CSS 문제만 최소 변경으로 고친다. 실제 계산 스타일·4.5 이상 대비와 D04 간격을 확인한다. implementation 문자열만 따라 쓰는 무의미한 새 테스트는 추가하지 않는다.
- [ ] 기존 관련 UI/폼 검사와 이번 코드 SHA의 등록 make lint/test를 실행하고 원 종료코드를 보존한다. 소스 입력이 바뀌었으므로 옛 gate PASS로 대체하지 않는다. DB/auth/adapter 입력 불변은 직접 대조하고 같은 검사를 무관하게 반복하지 않는다.
- [ ] 수정된 실제 fixture/캡처·색상/간격 측정과 exact 제품SHA를 tester/designer에 넘길 수 있게 준비한다. 기존4795 fixture는 두 QA가 아직 쓰므로 coor 확인 없이 종료/교체하지 않는다. 고유 포트·DB·프로필 또는 QA 완료 뒤 기존 실행을 stop 후 새 소스로 재기동하는 무삭제 경로를 정한다. 기존 실패/evidence는 보존한다.
- [ ] own 수정 QA/phase/context/outcomes를 기록하고 work.py finish·최종 lint·고정 SHA·실제 worker_done으로 회신한다.

수정 소유권은 두 CSS 원천과 꼭 필요한 기존 UI 검사·own QA/phase/context/inbox/outcomes다. 다른 역할 보고/PLANS/board/운영은 수정하지 않는다. 새 fixture helper 변경이 필요하면 테스트 경계 안에서 최소로 하고 이유를 적는다. 독립 tester는 대비/공개 영향·native200%, designer는 예시 경계/대표 화면을 다시 판정한다. 코드 고정 리뷰와 main 통합은 coor가 한다.

현재 설치 cache는 C:/Users/shin/.codex/plugins/cache/fullops-squad/fullops-squad/1.4.3이다. 필요한 해당 SKILL/help를 읽고 호환 도구를 사용하지만 레포하네스1.4.2·선택 수준·설정은 바꾸지 않는다. 이전 외부 gate 자원 AppData/Local/KnowsLinkDevTools/landing-guide001-dev는 원본 참고로만 읽고 새 attempt 자료는 별도 landing-guide001-css-dev에 보존한다. known_hosts/RejectPolicy·비밀 비출력·무삭제 계약을 유지한다.

기존 user-owned DEV terminal은 완료·유휴로 확인했으며 종료/프롬프트 주입하지 않는다. 문맥이 큰 이전 세션을 재개하지 않고 fresh Astra high로 같은 과제의 좁은 후속을 수행한다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-LANDING-GUIDE-001-DEV-packet.json` / SHA `c6154d3557d1704ea0b50b7db950d2f11054fa71` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read, document_update) · 줄 전체/미확인 · unknown · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CODE-REVIEW-review/report.md` (document_read) · 줄 6, 6, 10, 10 · inferred · 필수
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 27, 27, 28, 28, 35, 35, 36, 36, 37, 37 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/archive/designer-2026-10-14-landing.md` (document_read) · 줄 6, 6, 98, 98, 100, 100 · inferred
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 6, 26, 26, 28, 28 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 52, 52 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 7, 14, 14 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 7, 7, 14, 14, 20, 20, 268 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-DEV-test/runtime/report.md` (document_read) · 줄 257, 265 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CODE-REVIEW-review/author-final-handoff.md` (document_read) · 줄 1, 1, 11, 11 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6, 6, 10, 10, 56, 56 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/report.md` (document_read) · 줄 6, 6, 12, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-test/report.md` (document_read) · 줄 6, 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-review/report.md` (document_read) · 줄 36 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001-DEV.md` (document_read) · 줄 6, 6, 43, 43 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md` (document_read) · 줄 6, 6, 10, 10, 50, 50, 51, 51, 52, 52, 53, 53, 57, 57, 59, 59 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 24, 33 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 7, 7, 13, 13, 15, 15 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 7, 7, 13, 13 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_designer.md` (document_read) · 줄 9, 9, 10, 10, 13, 13, 17, 17, 25, 25, 67, 67, 101, 101, 106, 106, 188, 188, 189, 189, 190, 190 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_dev.md` (document_read) · 줄 357, 357, 358, 358, 361, 361, 365, 365, 373, 373, 384, 384, 386, 386, 393, 393, 394, 394, 429, 429, 448, 448, 469, 469, 470, 470, 476, 476, 514, 514 · inferred
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 33, 33, 36, 36, 51, 51, 64, 64, 67, 67, 91, 91, 92, 92, 107, 107, 108, 108 · inferred
- `.fullops-squad/handovers/to_ops.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 33, 33, 34, 34 · inferred
- `.fullops-squad/handovers/to_tester.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 25, 25, 27, 27, 34, 34, 35, 35, 85, 85, 92, 92, 95, 95, 141, 141, 142, 142, 218, 218, 219, 219 · inferred
- `docs/system-design.md` (document_read) · 줄 70 · inferred
- `internal/relay/human_pages.go` (impact_check) · 줄 14, 21, 22, 25, 29, 40, 42 · inferred
- `internal/relay/landing_pages.go` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/relay/landing_test.go` (impact_check) · 줄 39 · inferred
- `internal/relay/member_pages.go` (direct_edit) · 줄 31 · inferred
미확인 18건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

미착수. 수정 전 실패와 수정 후 실제 측정·새 gate·미실행·fixture 유지 책임을 구분한다.
