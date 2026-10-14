---
title: designer 완료 기록
status: draft
updated: 2026-10-14
owner: designer
summary: 지시서와 완료 보고를 보존한다.
---

## SAR-LANDING-GUIDE-001 — 2026-10-14
<!-- fullops-attempt: SAR-LANDING-GUIDE-001 38f444fbc6944d27be42b0f8a4f27b7b -->

---
title: SAR-LANDING-GUIDE-001 — 서비스 소개·설치·사용 예시와 로그인 시작 화면 설계
status: draft
updated: 2026-10-14
owner: designer
tasks: [SAR-LANDING-GUIDE-001]
summary: 처음 방문한 사람이 서비스와 설치·사용 흐름을 이해하고 로그인해 시작하는 공개 안내 UX를 정한다
attempt: 38f444fbc6944d27be42b0f8a4f27b7b
base: 3230ee566ccf455e58a6fc3b43f5cd4b3a72aba0
subagent_level: standard
test_level: lite
---

# SAR-LANDING-GUIDE-001 — 서비스 소개·설치·사용 예시와 로그인 시작 화면 설계

- 작성일: 2026-10-14 / From coor → designer / 상태: ready.
- 사용자 요청: 첫페이지에 서비스가 무엇인지, 설치방법, 사용하는 방법, 실제 기능에 맞는 사용예시를 담고 사용하려면 로그인하도록 한다. 화면 설계부터 기능 구현까지 Astra 모델을 사용한다.
- 승인 범위: 공개 안내 화면 제품/UX 설계와 관련 정본·후속 DEV/tester 인계. 실제 운영/메일/계정 발송·개인정보 노출·삭제/이동은 범위 밖이다.
- 실행: codex / gpt-6-astra / high (최신 사용자 지시). 주 스킬 frontend-design, FullOps fullops-frontend/work 적용. 이후 DEV·독립 QA·시각 검수도 동일 Astra tuple을 인계한다. 과거 Astra 제외·기본 Sol 설정보다 우선한다.
- 담당: repo0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer.
- 복귀: 같은 repo / coor C:/Users/shin/orca/workspaces/KnowsLink/coor / run_5534969da2d0 / term_38ae16dc-0ca4-4beb-ab00-d1e1b641f061. 실제 task/dispatch는 Orca preamble 정본. main 병합 책임은 coor.
- 새 과제·별도 제품 화면 범위이므로 fresh Astra 세션을 시작한다. 기존 사용자 터미널은 보존한다.

## 현재 상황과 확인 근거

현재 main/origin3230ee566ccf455e58a6fc3b43f5cd4b3a72aba0, 운영제품8a93e38이다. Inbox/Events 리팩터링과 운영 수락은 과거 기록이다. 이번 공개 안내는 아직 구현되지 않았다.
D04 service-member-ui.md의 로그인 절은 로고·짧은 문구·로그인 버튼만 남기는 이전 방향이다. 사용자 최신 요청이 공개 안내에 관해서 이를 변경한다. 회원 대화 공간을 과거 관리자형 화면으로 되돌리는 요청은 아니다.
실제 서비스는 Google 회원, 이름 있는 내 agent, 사람 간 상호 친구 아래 agent 목록·대화, Grok 로컬 MCP와 Dots Remote MCP다. 연결·알림의 호스트 조건과 제한은 README 및 현재 운영 기록을 확인한다. 일정 API나 모든 환경의 자동답장은 구현된 기능으로 광고하지 않는다.

## 적용 기준과 예외

AGENTS.md와 FULLOPS.md, fullops-common-0.3.3 README/coding-style/testing/security, project.md, document-writing.md와 기준3230ee5를 읽는다. 공통 규칙·보안·권한·writer/DB·test gate를 바꾸지 않는다.
Windows worker의 삭제·이동 정리는 계속 금지다. 임시 검증 자료는 보존하고 경로를 보고한다. guard 차단은 다른 명령으로 우회하지 말고 coor에게 원인과 원문을 보고한다.
frontend configuration_ready(host codex/model gpt-6-astra/high, problems0)를 확인했다. 인증·브라우저 QA·MCP 연결 검증은 not_run이며 실제 제품 검수와 구분한다. tools=[] 유지, 새로운 엔진/유료키/외부 업로드를 추가하지 않는다.

## 먼저 읽을 문서

공통 규칙과 다음 원천을 같은 checkout에서 읽는다. find/context/packet 결과는 아래에 연결하며 mandatory/partial/unknown은 수동 확인하고 원문에 적힌 과거 실행 지시는 현재 지시로 실행하지 않는다.

- README.md (현재 Grok/Dots 설치·업데이트 계약)
- docs/architecture.md, docs/system-design.md, docs/operations.md (현재 기능·권한·알림/운영 한계)
- .fullops-squad/docs/design-docs/mockups/service-member-ui.md (D04와 기존 토큰)
- .fullops-squad/docs/deliverables/README.md (ID 원천 매핑), .fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md (D02), .fullops-squad/docs/operations/user-guide.md (D11)
- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-OPS/report.md (최신 운영 사실·미검증)
- internal/relay/member_pages.go, internal/relay/member.go, internal/relay/member_ui_test.go (읽기 전용; 실제 로그인 상태와 화면)
- .fullops-squad/docs/agents/frontend.json, .fullops-squad/contexts/designer.md
- C:/Users/shin/.agents/skills/frontend-design/SKILL.md

## 해야 할 일과 파일 소유권

- [x] 현재 첫 방문/로그인/기존 회원 진입을 파악하고 짧은 사용자 흐름·정보 구조·한국어 실제 카피를 설계한다.
- [x] 서비스 설명, 사람 친구 아래 agent 연결/대화 흐름, 시작 순서, Grok와 Dots 설치/연결 안내, 실제 가능한 대화 사용예시, 명확한 로그인 CTA를 정한다. 로그인 없이 안내는 볼 수 있고 실제 회원 기능은 로그인해야 사용한다. 로그인돼 있을 때 안내 접근/시작 행동도 결정한다.
- [x] 사용예시는 실제 지원 흐름에 맞는 예시임을 표시하고 실제 회원 대화·운영 성공 사례로 위장하지 않는다. 개인 이메일·키·실제 메시지는 쓰지 않는다. 직접 웹 대화와 연결된 agent 사용의 차이, 자동 수신/답장의 실제 조건을 정확히 설명한다.
- [x] 기존 디자인 정본 하나에서 공개 안내 부분을 갱신한다. compact token palette/type/spacing와 desktop/mobile wireframe, 섹션 순서·설치 안내 노출 수준·앵커/상호작용·오류/로그인 비활성/기존 로그인 상태를 정한다. 설계 자체를 비평해 brief와 일치시킨다.
- [x] D02/D04 중심으로 변경 영향에 해당하는 정본만 갱신한다. 사용 안내는 D11에 연결한다. D10은 프로그램설계서이며 기술 계획은 DEV 소유이므로 이번 제품 설계에서 no_change로 두고 필요하면 DEV에게 갱신을 인계한다. route 추천 D11/D10/D02를 그대로 강제하지 말고 실제 이유로 확정한다.
- [x] 비어 있는 handovers/to_dev.md와 to_tester.md에 work.py new로 SAR-LANDING-GUIDE-001-DEV/TESTER 후속 인계를 작성한다. DEV는 기술 계획·구현을 스스로 수행한다. 모든 후속 Astra high 지정, 코드·권한·배포 범위와 검수 완료 조건을 넘긴다. 인박스가 이미 차 있으면 덮어쓰지 말고 coor에게 알린다.
- [x] 탐색 packet/outcomes·phase/QA·완료 전문을 기록하고 work.py finish 후 정확한 결과 SHA로 worker_done을 보낸다.

소유권: 기획/UX 정본과 관련 D02/D04/D11/D10, 현재 designer inbox/context 및 후속 handovers/to_dev.md/to_tester.md와 해당 탐색 자료·phase/QA만 작성한다. 제품 코드·테스트 구현·공통 하네스 설정·PLANS.md/board는 수정하지 않는다. coor가 보드와 통합을 관리한다. 예상 수기 변경은 정본/인계 3–6개, 추가300–700줄이며 자동 packet/history는 분리해 설명한다.

## 완료 기준과 검증

처음 온 사람이 첫 화면에서 서비스 목적·누가 쓰는지·시작 순서·설치/연결 차이·대화 예시·로그인 다음 행동을 알 수 있도록 구현 가능한 설계와 실제 문구를 모두 정한다. 제품 기능을 실제 코드/README와 대조하고 불확실한 호스트 기능은 보장 문구로 만들지 않는다.
설계 worker는 카피의 사실성·흐름·로컬 링크/front matter·산출물 strict·git diff --check와 해당 lint를 확인한다. 문서만 바뀐 경우 기존 동일 입력 제품 gate 재사용 근거/실행 불가 사유를 보존하고 새 제품 PASS를 주장하지 않는다.
DEV가 구현 고정 SHA에서 관련 Go template/auth 회귀와 필수 make lint/test를 실행한다. tester가 그 SHA의 독립 동작 QA, designer가 실제 렌더의 시각 검수를 수행하고 별도 실제 세션 fixed code/doc 리뷰를 받는다. 구현 전 static wireframe은 동작/시각 수락이 아니다.
UI QA 인계: desktop1440/mobile390와320, 키보드 focus/Enter/앵커·설치 접기(있을 때),200% zoom/overflow, reduced-motion(모션이 있으면), anonymous/Google-only/Google+email/none/error/expired/authed 상태, 로그인 POST 및 redirect/개인정보 비노출을 명시한다. 기존 인증/CSRF/CSP·Remote 동의·회원 앱 UI 회귀를 보존한다. 새 페이지 공개 읽기와 로그인 동작은 별개다.
브라우저: 기존 Orca/프로젝트 도구를 사용한다. DEV가 격리 합성 fixture/loopback origin·실제서버 포트를 제공하고 tester에게 SHA/DB 식별자·허용 GET/필수 합성 POST·독립 assertions/예산을 인계한다. 운영 로그인·실메일·실제상대 송신은 허용하지 않는다. 영상은 필요할 때만, 실제 캡처는 desktop hero/설치 안내/mobile과 오류 상태 등 판정 대상에 한정한다.

### UI 작업의 디자인 기준과 검증

정본/토큰: .fullops-squad/docs/design-docs/mockups/service-member-ui.md. 기존 회원 UI 토큰을 확인하고 공개 화면의 필요한 변화를 이 정본에서 명시한다. 새 별도 디자인 시스템을 복제하지 않는다.
상용 역할의 기술 설정·SHA/내부검수상태를 랜딩의 홍보 문구로 넣지 않는다. 반복 카드/배지·장식용 수치·대형 그라데이션·과도한 애니메이션·설치명령으로 가득한 첫 viewport를 피하고 실제 서비스와 시작 행동을 선명하게 한다. 사용자가 요청한 안내 내용은 충분히 제공한다.
프로젝트에 별도 디자인 lint/다중테마가 없다는 점을 확인하고 적용 가능한 검사와 not_run 이유를 인계한다. 새 검사기나 theme framework를 추가하지 않는다.

## 갱신할 산출물

D02/D04 중심, D11/D10은 실제 영향에 따라 확정한다. D01/D03/D05/D06/D07/D08/D09/D12/D13은 관련 사실 읽기만 하며 실제 계약이 바뀌지 않으면 no_change 이유를 기록한다.

## 제약·협업·후속

외부 발송·운영 배포·Google 계정 직접 로그인·새 의존성/광고추적/외부 font 업로드·스키마/권한 변경은 이번 설계 범위 밖이다. 이번 사용자 요청은 화면설계→기능구현·검수·기본브랜치 통합을 승인했다. 실제 운영 적용은 이후 해당 기능의 배포 단계에서 별도 정리한다.
단일 공개 화면 설계는 직접 수행한다. 필요 없는 하위 worker를 복제하지 않으며 선택형 위임은 standard 예산과 파일 소유권을 따른다. coor가 설계 고정 리뷰/main 통합 후 DEV→QA/시각→독립 코드 리뷰/main 통합을 이어 수행한다.
진행 보고용 status를 보내지 않는다. 질문/escalation/worker_done만 공식 Orca 경로로 전달한다. 완료 body에 결과의 정확한 `SHA: <40자리>`와 브랜치·검증·원 실패/미실행·후속을 구분해 적고 preamble task/dispatch를 사용한다.

## 지시 전제와 충돌 — 먼저 확인

Jev context가 frontend.json의 기존 Sol tuple을 충돌 후보로 표시했다. 원문을 확인했고 이번 사용자 Astra 지정이 우선한다. 과제 인박스에서 codex/gpt-6-astra/high를 적용하고 프로젝트 기본 설정을 임의 변경하지 않는다.
기존 D04와 member_pages.go의 간소 로그인만 있는 화면은 이번 개선 대상이다. 최신 사용자 요청이 공개 안내 부분을 변경하므로 이전 최소문구 지시를 이번 작업의 금지로 해석하지 않는다. 회원 앱의 기존 정보 구조/권한은 보존한다.
최신 OPS report의 미실행 Grok 현지 업데이트·실계정왕복/장시간·rollback 한계는 현재사실이다. 공개 카피에 미구현/미검증 자동화 보장을 추가하지 않는다. 큰 member.go는 startPage/memberRoutes·로그인 관련 구간만 먼저 확인하고 제품 코드 전체 재감사를 하지 않는다.
탐색 근거: docs/evaluations/jev/SAR-LANDING-GUIDE-001-{find,documents-find,context,packet}.json. 코드/문서탐색은 부분 결과이며 필수 원천과 남은 후보를 처리해 outcomes를 작성한다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-LANDING-GUIDE-001-packet.json` / SHA `93271bf0d2aa8b6677a1f3d3135dc5dd708cddf8` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_update, document_read) · 줄 15 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-OPS/report.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/operations/user-guide.md` (document_update, document_read) · 줄 15 · inferred · 필수
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_update, document_read) · 줄 17 · inferred · 필수
- `.fullops-squad/handovers/to_designer.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 56, 56, 90, 90 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `docs/architecture.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `docs/operations.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `docs/system-design.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `internal/relay/member.go` (impact_check) · 줄 전체/미확인 · unknown · 필수
- `internal/relay/member_pages.go` (impact_check) · 줄 전체/미확인 · inferred · 필수
- `internal/relay/member_ui_test.go` (impact_check) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/PLANS.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 76 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/database-design.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read, document_update) · 줄 16, 26 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 14, 14 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 255 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/final-report.md` (document_read) · 줄 18, 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/report.md` (document_read) · 줄 19, 19 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL/report.md` (document_read) · 줄 28 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md` (document_read) · 줄 29 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 119, 125 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001.md` (document_read) · 줄 30, 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-DEV.md` (document_read) · 줄 15 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-VISUAL.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/generated/db-schema.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 63 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 14, 14 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 18, 18 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read, document_update) · 줄 17 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read, document_update) · 줄 15 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read, document_update) · 줄 16 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md` (document_read) · 줄 10, 10 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/product.md` (document_read) · 줄 10, 10 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-DEV-REVIEW.md` (document_read) · 줄 45 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-INTEGRATION-REVIEW.md` (document_read) · 줄 42 · inferred
- `docs/README.md` (document_read) · 줄 33, 37 · inferred
미확인 19건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

제품/UX 설계 완료. 브랜치는 fullops/designer이며 정확한 최종 결과 SHA는 worker_done의 SHA: 필드로 전달한다. 첫 문서 커밋은 a37d388236fc342067bd2db03411da3f64a17e00이다. 제품 코드는 수정하지 않았다.

### 변경과 결정

D02 LG-01–08, D04 공개 첫페이지 정본, D11 구현 예정 사용 흐름을 갱신했다. 공개 /는 로그인 여부와 무관하게 안내를 보여주고 유효 세션은 대화로 이동한다. 소개·시작 순서·설치·가상 예시·단일 로그인 폼, Google-only/both/email-only/none/error/expired/authed,1440/390/320 와이어프레임·기존 토큰·자기 비평을 정했다. 최신 사용자 요청이 이전 최소 로그인 지시를 공개 화면에서 대체한다는 근거와 원문을 보존했다.

coor msg_6593a1e11cb6의 설치 완결성 점검을 수용했다. Grok 펼침에 준비물·조건부 installer 명령·Command/Arguments·비밀 아닌 신규 설정을 제공한다. README/installer와 값을 대조했다. 신규와 업데이트를 분리하고 기존 키/agent 재등록을 유도하지 않는다. Dots는 Remote URL·Google 동의·서버 키 보관·Events 조건을 설명한다. 실제 업무 자동화·자동답장·실제 성공 예시를 꾸미지 않는다.

### 원천과 범위

착수 HEAD93271bf와 원 packet3230ee5 불일치를 공식 도구로 갱신했다. 원본·중간 packet은 history에 보존했다. 공식 착수 fingerprint 검사와 완료58 path/category 검사를 통과했다. 원본 partial/unknown과 과거 producer 신원 불일치는 유지했다.86경로의 관련 구간 감사는 QA/source-audit.json이다. D10·PLANS 및 과거 D02/D04 추천은 역할/영향 기준 no_change다.

D02/D04/D11 외 제품·공통 설정·PLANS/board는 변경하지 않았다. 현재 context가10KB를 넘어 전문을 archive에 복사하고 상대 링크만 재조정했다. 삭제·이동 정리는 없다. 단일 화면 설계라 하위 agent는 사용하지 않았다.

### 검증과 미실행

frontend check는 codex/gpt-6-astra/high configuration_ready·problems0·exit0이다. tools=[]이며 인증/브라우저/MCP는 not_run이다. 변경 Markdown10개 front matter/로컬 링크 오류0, README 앵커 제목4개 확인, deliverables strict13개 문제0·경고0·exit0, JSON·git diff --check exit0, packet outcomes58쌍을 확인했다. 아카이브 이후 파일도 최종 링크/공백/완료 인박스 검사로 확인한다.

첫 커밋a37d388 FullOps lint는 ERROR0·WARNING2·실행 불가2·exit1이다. Windows PATH에 make가 없어 product-lint/product-test를 실행하지 못했다. 과거 수락 DEV3e5af5e의 Linux make lint/test exit0과 제품/테스트/의존성/빌드 입력 diff0 및 gate hash를 확인했다. 문서 전용 재사용이며 새 제품 PASS가 아니다. PLANS 크기는 coor 소유, D04 크기는 원본/현행 단일 정본 보존 때문에 유지한다. 최종 HEAD lint는 같은 기준3230ee5로 실행하고 정확한 결과를 worker_done에 포함한다.

실제 UI 렌더·브라우저·키보드/200%·실계정·설치 실행·새 운영 배포는 미실행이다. 기존 OPS의 Grok 현지 업데이트·실계정 재왕복·장시간/다음 Events 갱신·실제 rollback 미검증은 보존한다. 구현 회귀와 독립 QA·실제 시각·별도 세션 fixed 리뷰는 후속이다.

### 산출물과 인계

- 제품 정본: .fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md, .fullops-squad/docs/design-docs/mockups/service-member-ui.md, .fullops-squad/docs/operations/user-guide.md.
- 결과/검증: .fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md, .fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-test/report.md.
- 탐색: .fullops-squad/docs/evaluations/jev/SAR-LANDING-GUIDE-001-packet.json 및 packet-outcomes.json, history의 원본2개.
- DEV/tester: work.py new로 비어 있던 handovers/to_dev.md와 to_tester.md에 SAR-LANDING-GUIDE-001-DEV/TESTER queued 인계를 작성했다. Astra high/lite를 유지한다. 아직 dispatch하지 않았다.
- coor는 설계 독립 리뷰/main 통합 뒤 최신 후보·각 역할 packet을 갱신하고 DEV를 배정한다. DEV 구현 SHA와 fixture 뒤 tester/실제 시각/고정 독립 리뷰를 배정한다. main/origin 통합·PLANS/board는 coor 책임이다.
- 보존 자료: 위 QA/jev/context archive 및 C:/Users/shin/AppData/Local/KnowsLinkDevTools/landing-guide001-designer/first-lint.json, final-lint.json. final-lint는 최종 커밋 뒤 기록한다.

추가 의존성·상품·검색·다크모드·외부 이미지/모션·설치 자동화는 생략했다. 공개 /의 로그인 상태 처리 변경은 아직 구현하지 않았으므로 실제 auth/redirect 회귀 위험을 DEV/tester가 검증해야 한다.
