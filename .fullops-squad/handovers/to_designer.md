---
title: SAR-SSE-INSTALL-001-VISUAL — 새 Google OAuth 동의 화면의 독립 desktop/mobile 직접 시각 검수
status: draft
updated: 2026-10-11
owner: designer
tasks: [SAR-SSE-INSTALL-001-VISUAL]
summary: 새 Google OAuth 동의 화면의 독립 desktop/mobile 직접 시각 검수
attempt: 10e9aad1f3f24addad543a18cc892f06
base: febd1282df2e84a29fa0a53834569e7b9cc028e2
subagent_level: standard
test_level: lite
---

# SAR-SSE-INSTALL-001-VISUAL — 새 Google OAuth 동의 화면의 독립 desktop/mobile 직접 시각 검수

- From / To: coor / designer. 상태 ready. Test lite / subagent standard.
- 담당 repo 0b08ec4c-9e3a-4613-8197-5a835e545335, C:/Users/shin/orca/workspaces/KnowsLink/designer, fullops/designer.
- 복귀 coor C:/Users/shin/orca/workspaces/KnowsLink/coor / term_e61d3e14-29e9-4954-943a-4a75707c82de / run_dc9b4e6dd18a. 실제 task/dispatch는 preamble.
- 고정 후보 febd1282df2e84a29fa0a53834569e7b9cc028e2. 부모 구현자01a12678-994d-7680-a15c-83577943b52d와 다른 시각 검수 세션. 동작 QA/코드 리뷰는 tester가 별도로 진행 중이다.

## 현재 상황·기준

사용자 요청 SSE/API 수신과 Google 인증 호스트 설치 후보다. 새 remote OAuth 동의 화면은 기존 memberStyle을 재사용한다. DEV는 mobile synthetic 렌더와 긴 URI overflow를 확인했지만 최종 desktop preview는 background helper 자동승인 차단으로 미실행이다. 실제 Google 로그인/외부 Dots 설치는 합성 UI 검수와 구분한다.

fullops-common-0.3.3 및 연결 coding/testing/security, project.md, document-writing.md, rules/delegation.md, orca-agents.md의 designer 직접 시각 검수 책임을 따른다. 예외·보호완화 없음. 기존 디자인 기준은 D04와 memberStyle이며 새 디자인/기획 승인은 범위 밖이다.

## 먼저 읽을 문서

.fullops-squad/FULLOPS.md, project.md, rules/common/README.md 및 세 연결 규칙, docs/agents/document-writing.md, docs/exec-plans/phases/SAR-SSE-INSTALL-001-DEV.md의 UI 검수 절, docs/exec-plans/phases/SAR-SSE-INSTALL-001-REMOTE.md, docs/design-docs/mockups/SAR-MVP-001-UI.md, internal/relay/remote_auth.go의 remotePage 및 memberStyle 정의만 읽는다. 연결 packet partial/unknown은 실제 UI 영향 범위로 점검하고 역사적 unrelated 자료는 no_change 이유로 구분한다.

## 해야 할 일·소유권·완료 기준

- [ ] 고정 후보 UI source와 실제 렌더 대상이 동일한지 확인한다. 기존 DEV synthetic preview 자료가 있으면 재사용하되 새로운 시각 판정임을 구분한다.
- [ ] desktop1280×800 및 mobile390×844에서 최초 Google 확인과 설치 승인/취소 화면을 직접 렌더해 가독성·버튼구분·긴 callback URI 줄바꿈·가로 overflow·기존 memberStyle 재사용을 판정한다. 실제 로그인·승인 제출 없이 합성 fixture를 사용한다. 의미 있는 화면별 캡처와 viewport/source SHA를 보존한다.
- [ ] 웹은 Orca embedded browser만 사용한다. OS computer-use는 Grok Bot만 허용되어 웹 검수에 사용하지 않는다. 안전한 file preview 또는 foreground localhost synthetic 렌더는 가능하다. background helper 자동승인 거부를 동일 방법으로 재시도하거나 보안 설정을 완화하지 않는다. 필요한 임시 화면은 검사 후 정리한다.
- [ ] 실패/미실행/시각 finding을 구체적 근거·영향으로 보고한다. 실제 Google 로그인/동의·Dots 실행은 미실행이며 PASS로 바꾸지 않는다. 안전한 렌더가 불가하면 구체적 오류·미판정 항목·후속 담당 coor/사용자·재개조건을 남긴다.

소유권은 qa-reports/SAR-SSE-INSTALL-001-VISUAL/의 검수 보고서·캡처·필요 synthetic QA fixture와 시나리오뿐이다. 제품 Go/adapter/기획/테마/호스트 전역설정/운영 DB·키·관계·서비스는 수정하지 않는다. 실제 계정 접속·송신·유료 API·local Docker 금지. 인증/동작 보안 분석은 tester 담당이며 제품 결함은 보고만 한다.

## 산출물·완료 보고

D04 관련 직접 시각 검증 근거를 QA 보고서로 연결하며 기존 mockup 원천은 바꾸지 않는다. 그림만으로 판정할 수 있는 항목만 캡처하고 동작QA를 복제하지 않는다. packet path/category별 outcomes·partial uncertainty 근거를 기록한다. 실제 source SHA·viewport·직접 관찰·미실행·findings·증거 경로·lint/문서검사 범위를 남기고 work.py finish로 정본과 전문을 보존한다. 실제 preamble worker_done으로 복귀한다.

수동 탐색 보완: jev_find code가80파일 지도에서 unclear를 반환했다. 실제 rg는 internal/relay/remote_auth.go:312의 remotePage와 remoteConsent를 확인했다. 후보의 직접 시각 원천은 remotePage 및 memberStyle이며 없다고 단정하지 않는다. packet의 seed/read는 검사 대상일 뿐 제품 수정 권한이 아니다.

memberStyle 정본은 rg로 확인한 internal/relay/member.go다. 최초 context의 member_auth.go는 존재하지 않는 후보로 거부되어 실제 member.go로 보완했다. 후보부재를 미해결 제품 결함으로 해석하지 않는다.

