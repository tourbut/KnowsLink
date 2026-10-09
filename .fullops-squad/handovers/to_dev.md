---
title: SAR-GOOGLE-CONNECT-001-DEV — Google 로그인만으로 자기 클라이언트 agent 연결을 완료하고 Cloudflare 로그인과 token 파일 수동 전달을 제거한다
status: draft
updated: 2026-10-09
owner: dev
tasks: [SAR-GOOGLE-CONNECT-001-DEV]
summary: Google 로그인만으로 자기 클라이언트 agent 연결을 완료하고 Cloudflare 로그인과 token 파일 수동 전달을 제거한다
attempt: d5a1fab9baf24084b01f5c4e66f359fa
base: f9f7675be6c9502bd6ab2f810bd42ff62a26a174
test_level: lite
---

# SAR-GOOGLE-CONNECT-001-DEV — 양쪽 Google 로그인으로 agent 연결 완료

- From / To: coor / dev. 상태: ready.
- Task key: SAR-GOOGLE-CONNECT-001-DEV; Purpose: implementation.
- 기준 SHA: c3f918dfb89c6c670075cddf7a1cc351d1397bf2; 테스트 레벨: lite.
- 담당: repo 0b08ec4c-9e3a-4613-8197-5a835e545335, C:/Users/shin/orca/workspaces/KnowsLink/dev, fullops/dev.
- 복귀: C:/Users/shin/orca/workspaces/KnowsLink/coor, term_e61d3e14-29e9-4954-943a-4a75707c82de, Run run_86e0e674b5a0. task/dispatch는 실제 preamble을 따른다.

## Windows 재개 지시 — 2026-10-09

사용자 “대기 작업 진행해”로 이 과제의 중지를 해제했다. 제품 목표·확정 규칙·범위·완료 조건은 아래 원문을 유지한다. 기술 계획·구현·검증은 DEV가 같은 과제에서 끝낸다.

- Purpose: implementation; Test level: lite; Subagent level: off. 선택형 재위임은 하지 않는다.
- 현재 기준 ref: `f9f7675be6c9502bd6ab2f810bd42ff62a26a174`. 아래 과거 c3f918 기준과 Linux 복귀 주소는 역사 기록이며 이번 실행에는 사용하지 않는다. 최종 lint --from은 현재 기준 ref를 사용한다.
- 원천·정본은 자기 체크아웃의 공통 규칙·project.md·아래 제품 규칙이다. API·라이브러리 문서가 필요하면 Context7으로 확인한다. 패킷의 오래된 producer와 partial은 근거를 확인하고 직접 탐색으로 보완한다.
- Windows의 make 부재는 확인됐다. 필수 product-lint/product-test를 생략하거나 완화하지 않는다. 담당자가 기존 검사에 필요한 로컬 실행 환경을 복구하고 명령 자체의 종료코드를 보존한다. 설치·실행이 불가능하면 구현과 실행 가능한 검사를 끝낸 뒤 정확한 장애·남은 검사·재개 조건을 보고한다.
- 비공개 `.env`는 main 파일의 심볼릭 링크다. 값을 출력하거나 Git에 넣지 않는다. 운영 서버·Google/Cloudflare 계정·DB 변경은 coor/ops에 인계한다. 운영 자격을 구현 선행조건으로 만들지 않는다.
- 부모는 완료 후 고정 SHA 독립 리뷰·필요한 tester QA·운영 인계를 처리한다. 종료 전에 packet outcomes·work.py finish·커밋·등록 lint/test·authentic worker_done을 완료한다.

## 목표·확정 규칙·범위

사용자 원문: “그냥 서로 양쪽에서 메일로 로그인하면 되게해”, 이어 “클라우드 플레어는 뭐고 그냥 구글 로그인하면 되게해”, “fullops에서 테스트 레벨은lit로 변경해”.
Google 로그인만으로 자기 클라이언트의 agent 연결을 끝내는 가장 작은 완결 흐름을 구현한다. 사용자가 token 파일·개인키·credential을 옮기거나 연결 token을 복사하지 않는다. 서비스의 사용자 경로에서 Cloudflare 로그인은 요구하지 않는다. 기존 Google issuer/sub 회원·세션과 Node/Command MCP를 재사용한다. 양쪽은 같은 Google 회원 아래 별도 agent·키·credential을 사용한다. 개인키는 각 클라이언트에서 생성·보관한다. Google 로그인 자체가 권한 없는 타인 agent 연결을 허용하지 않는다. 기존 명시적 동의·관계 수락은 유지하되 사용자에게 식별자/비밀 문자열 복사 대신 가능한 기존 UI 동작을 제공한다. 임의 자동 pairing 등 제품 규칙을 새로 정하지 않는다.
기술 파악·짧은 계획·구현·lite 검사·필요한 기술/사용 문서 갱신은 같은 DEV 과제에서 수행한다. 새 제품 규칙/범위 판단만 coor에 ask한다. 기술 방식 선택은 DEV가 맡는다.

## 현재 상황과 근거

Google 실제 회원 로그인과 자기 홈은 운영 a9a0db49745542727611ee863e0beeb03c53288e에서 동작한다. 공개 /v1/connect/{info,prepare,complete}, /v1/text/{send,pull,persist,ack}, /v1/keys/{agent}/{kid}, /v1/receipts/{id}만 relay 인증까지 도달한다. /home·기존 root는 아직 Cloudflare owner-only 보호다. 실제 Node grant info 200과 음성 경계·공유 회귀는 COOR/open-readiness/nou-public-access.json에 있다.
기존 노우는 /workspace/KnowsLink와 Node22 설치·빌드를 끝냈다. /workspace/.knowslink-connect/nou/는 비어 있는 0700 폴더다. 새 일반 회원 agent agent_af951c12bcd21106357539는 미연결이다. 이전 10분 grant는 자동 만료되므로 수동 token 전달을 재개하지 않는다. 실제 Bot 컴퓨터는 이 워크트리와 다르다. 로컬 성공을 실제 계정 성공으로 보고하지 않는다.

## 적용 기준과 먼저 읽을 문서

공통 기준 fullops-common-0.3.3: .fullops-squad/rules/common/README.md, coding-style.md, testing.md, security.md. 프로젝트 정본 .fullops-squad/project.md. 모두 위 기준 SHA다. test_level=lite이며 보안·데이터 손실 방지·필수 product-lint/product-test·독립 리뷰는 유지한다. 이전 광범위 QA는 재실행하지 않는다.
먼저 .fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md, .fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md, README.md, adapters/README.md, internal/relay/google.go, internal/relay/member_agents.go, internal/relay/connections.go, adapters/src/connect.ts를 존재 확인 후 필요한 범위만 읽는다. 파일명 추측이 틀리면 rg로 확인한다. find/context/packet의 추가 필수·충돌·주의 후보도 확인한다. 과거 token 수동 전달 안내는 사용자 최신 지시로 대체되며 보안 원칙은 유지한다.

## 해야 할 일과 소유권

- [ ] 기존 인증·연결·MCP 호출자/설정/fixture를 확인하고 가장 작은 기술 계획을 실행 기록에 적는다.
- [ ] 사용자 Google 로그인으로 클라이언트 연결·키 확인/동의·자격 저장까지 이어지고 token 파일 복사가 없는 경로를 구현한다. replay·다른 회원·다른 클라이언트에 대한 연결 탈취를 막는다.
- [ ] 기존 Google 가입/재로그인, 수동 connect와 기본 held/시험 모드의 필요한 호환을 유지한다. 필요한 Node/MCP 설정 및 문서를 같은 과제에서 완결한다.
- [ ] 사용자 경로의 Cloudflare 로그인 제거를 위한 실제 필요한 URL·ingress/Access 경계·적용/복귀·검증 방법을 배포 인계에 적는다. 제품 deploy 코드 수정은 담당 범위이며 운영 Cloudflare/DB 직접 변경은 coor/ops가 수행한다.
- [ ] 핵심 성공·실패/만료·replay·타 회원/키 분리의 좁은 검사와 필수 등록 lint/test를 통과하고 패킷 outcomes·완료 보고·archive·커밋 후 authentic worker_done을 보낸다.

제품 코드 cmd/, internal/, adapters/, db/, scripts/, deploy/ 및 관련 README·기술 문서가 DEV 소유다. 필요 없는 schema/dependency 추가는 피한다. 운영 state의 비밀·실제 Google OAuth/Cloudflare 계정·Tunnel·DB는 변경하지 않는다. 새 유료 서비스·외부 메시지 발송·강제 Git 명령은 승인 범위 밖이다.

## 완료 기준·검증·후속

사용자가 자기 클라이언트에서 연결을 시작하고 자기 브라우저의 Google 로그인/필요한 동의를 마치면 해당 클라이언트에만 자격이 자동 저장된다. token 파일·키·credential의 사용자 수동 전달은 없다. 서로 다른 클라이언트가 같은 Google 회원에 독립 agent로 연결되며 다른 회원/클라이언트가 연결 결과를 가로채지 못한다. 기존 로그인/동작 회귀를 좁게 검사한다. 외부 Bot의 앱/MCP 권한 제약은 실제 근거로 구분한다.
DEV는 커밋한 현재 SHA에서 python3 <설치 플러그인>/scripts/lint.py --repo <자기 체크아웃> --from c3f918dfb89c6c670075cddf7a1cc351d1397bf2를 통과한다. 등록 product-lint/product-test와 각 명령 종료코드를 보존한다. 전체 재검증·장시간 외부 QA는 추가하지 않는다. 후속 담당 coor는 별도 세션 고정 SHA 인증 delta 리뷰·필요한 좁은 QA/UI·main 통합/push·운영 적용 후 실제 사용자 경로 1개를 확인한다. 제품 최종 운영 왕복은 별도 실제 근거로 판정한다.
UI는 기존 Go template/CSS를 재사용한다. 새 디자인 시스템·프레임워크·테마를 도입하지 않는다. 설정 내부 용어·보안 token·Cloudflare를 사용자 연결 화면에 노출하지 않는다. 디자인 lint는 기존 프로젝트 적용 범위만 검사하고 변경 화면의 성공/오류/모바일 접근성을 짧게 확인한다.

## 갱신할 산출물과 기대 결과

route 추천 D10·D11·D12를 실제 영향에 맞춰 갱신한다. D03·D05·D06·D09·D13은 영향 확인 후 필요한 것만 갱신하고 불필요하면 이유를 적는다. 실행 기록은 .fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-001-DEV.md다. 기존 실패/보류·원본 증거를 보존한다.
완료 시 현재 인박스에 변경 이유·기준/최종 SHA·검사 결과/한계·외부 설정 인계·후속을 쓰고 work.py finish로 전문을 보존한다. 수정 규모/의존성 경고가 있으면 최소 변경 근거를 남긴다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad/docs/evaluations/jev/SAR-GOOGLE-CONNECT-001-DEV-packet.json` / SHA `c3f918dfb89c6c670075cddf7a1cc351d1397bf2` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md` (document_read) · 줄 2, 6, 7, 10, 14, 16, 18, 27, 27, 35, 39 · inferred · 필수
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read, document_update) · 줄 25, 27, 36, 37, 38, 40, 41, 42, 48, 54, 57, 67, 71, 73, 74, 76, 80, 82, 92, 93, 94, 110, 112, 128, 144, 158, 160, 160 · inferred · 필수
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 2, 2, 6, 6, 7, 7, 13, 13, 13, 13, 16, 16, 24, 24, 29, 29, 30, 35, 40, 41, 45, 49, 55, 55, 68, 69, 69 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 57 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 22, 32, 50, 55, 57, 67, 69, 69, 70, 70 · inferred
- `.fullops-squad/contexts/ops.md` (document_read) · 줄 65 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 7, 21 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read, document_update) · 줄 27, 44, 50, 51, 98, 102, 128, 138, 140, 147, 148, 149, 154, 157 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read, document_update) · 줄 20, 22, 27, 40, 41, 43, 49, 50, 60, 73, 76 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read, document_update) · 줄 25, 35, 40, 46, 50, 51, 55, 56, 58, 59, 64 · inferred
- `.fullops-squad/docs/design-docs/database-design.md` (document_read, document_update) · 줄 38, 44, 49 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read, document_update) · 줄 9, 17, 28, 30, 32, 33, 34, 35, 36, 41, 48, 49, 50, 51, 52, 53, 54, 56, 72, 73, 77, 92, 104, 109, 133, 139, 152, 158, 159, 163, 164, 165, 166, 167, 168, 172, 172, 173, 180, 181, 192, 193, 201, 202, 203, 204, 208, 213 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 32, 41, 59, 60, 61, 62, 63, 64, 65, 66, 69, 73 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md` (document_read) · 줄 24, 41, 42, 46, 48, 51, 52, 53, 68, 70 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 45, 46 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md` (document_read) · 줄 36, 45, 48 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 7, 22, 23, 24, 25, 27, 39, 41, 43, 45, 47 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read, document_update) · 줄 7, 23, 33, 35, 70, 71, 92, 98, 99, 100, 104, 111, 113, 115, 118, 119, 128, 128, 128, 129, 131, 132, 146, 148, 166, 168, 170, 172, 173, 173, 174, 176, 176 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/COOR/dev-fix-2-final/coordinator-handover-supplement.md` (document_read) · 줄 26 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.2.0-review/report.md` (document_read) · 줄 22, 54, 58 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-FINAL-review/report.md` (document_read) · 줄 40 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-review/report.md` (document_read) · 줄 60 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-TESTER.md` (document_read) · 줄 54, 55 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-002-TESTER.md` (document_read) · 줄 17, 36, 79 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-DEPLOY-001-OPS-FINAL-review/report.md` (document_read) · 줄 31 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-LOGIN-001-REVIEW-review/report.md` (document_read) · 줄 2, 6, 7, 10, 16, 18, 20, 26, 28, 32, 36 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FINAL-review/report.md` (document_read) · 줄 71, 77, 83, 104, 106, 135 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FIX-review/report.md` (document_read) · 줄 78, 83, 92, 106, 117, 118, 124, 126, 139 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FINAL.md` (document_read) · 줄 7, 15, 17, 34, 39, 45, 49, 50, 51 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX.md` (document_read) · 줄 16, 46, 48, 49, 56, 90, 95 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md` (document_read) · 줄 34, 41 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TESTER.md` (document_read) · 줄 63 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-review/report.md` (document_read) · 줄 28, 30 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-FIX-TESTER.md` (document_read) · 줄 65, 68, 74, 80 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 7, 35, 36, 37, 38, 39, 40, 41, 48, 50, 51, 52, 54, 77 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-TESTER-test/probe-failures.md` (document_read) · 줄 25, 43 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-TESTER.md` (document_read) · 줄 2, 10, 60, 62, 66, 84, 89, 90, 103, 110, 111, 112, 116, 117 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/scenario.md` (document_read) · 줄 22, 23, 25 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TESTER.md` (document_read) · 줄 47 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-UI.md` (document_read) · 줄 59, 60, 62, 92 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER.md` (document_read) · 줄 48, 56, 57 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 39, 64, 72, 82, 84 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX-2.md` (document_read) · 줄 28, 37 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-SERVICE-001-REVIEW-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-BETA-001-TESTER.md` (document_read) · 줄 25 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER-FINAL.md` (document_read) · 줄 7, 25, 27, 29, 38 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER-FIX.md` (document_read) · 줄 38, 39, 44, 45 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER.md` (document_read) · 줄 30, 37 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-AGENTS-001-FIX-TESTER.md` (document_read) · 줄 37 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-AGENTS-001-TESTER.md` (document_read) · 줄 29, 39 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-IDENTITY-001-TESTER.md` (document_read) · 줄 30, 31 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER.md` (document_read) · 줄 15 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.2.0.md` (document_read) · 줄 45 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 44 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-DEPLOY-001-OPS.md` (document_read) · 줄 20, 46 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-001-DEV.md` (document_read) · 줄 20, 54, 110, 114, 116, 132, 143, 148, 149, 150, 154, 166 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md` (document_read) · 줄 29 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md` (document_read) · 줄 29 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md` (document_read) · 줄 21, 117 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md` (document_read) · 줄 24, 32, 47, 70, 90 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-PUBLIC-POLICY-001.md` (document_read) · 줄 22, 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PREP-002.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md` (document_read) · 줄 7, 24, 26, 32, 34, 36, 38, 82, 84, 92, 93 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX.md` (document_read) · 줄 28, 32, 35, 38, 45, 47, 65, 90 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 2, 7, 16, 18, 24, 27, 29, 30, 34, 50 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` (document_read) · 줄 24, 49, 51, 53, 69, 75, 78, 80, 82 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2.md` (document_read) · 줄 24, 26, 46, 47, 62, 66 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-3.md` (document_read) · 줄 7, 33, 35, 37, 43, 65, 70, 72 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX.md` (document_read) · 줄 32, 59 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (document_read) · 줄 36, 46, 48 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPEN-PREP.md` (document_read) · 줄 53 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 18, 35, 48, 51, 54, 88, 103, 113, 220, 221 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SETUP-001.md` (document_read) · 줄 15, 69 · inferred
- `.fullops-squad/docs/generated/db-schema.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_read, document_update) · 줄 79, 119, 120, 184, 217, 221, 261, 288, 302, 304, 306, 308 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read, document_update) · 줄 74, 76, 78, 80 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read, document_update) · 줄 85, 91 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 25, 35, 45 · inferred
- `.fullops-squad/docs/planning/SAR-PREP-002-request.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/planning/SAR-SETUP-001-request.md` (document_read) · 줄 12, 14 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read, document_update) · 줄 20, 24, 38 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read, document_update) · 줄 17, 31, 35, 40, 41, 52, 56, 107, 123, 129, 145 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read, document_update) · 줄 20, 21, 30, 31, 32, 33, 34, 35, 36 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/README.md` (document_read) · 줄 8, 19, 37 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md` (document_read) · 줄 13, 24, 42, 43, 83, 172 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/product.md` (document_read) · 줄 7, 56, 94, 95, 121 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/protocol.md` (document_read) · 줄 9, 13, 14, 30, 45, 72, 74, 140, 149, 162, 198, 209, 248, 261 · inferred
- `.fullops-squad/handovers/SAR-MVP-001-REVIEW.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_designer.md` (document_read) · 줄 34, 145, 146, 147, 148, 149, 150, 151, 152, 221, 222, 223 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_ops.md` (document_read) · 줄 76, 90 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_dev.md` (document_read) · 줄 232, 233, 234, 236, 241, 251, 255 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_ops.md` (document_read) · 줄 169, 240 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_designer.md` (document_read) · 줄 49, 73, 122, 128 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_dev.md` (document_read) · 줄 28, 57, 60, 66, 82, 86, 97 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_ops.md` (document_read) · 줄 23, 45, 113, 115, 138 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_tester.md` (document_read) · 줄 55, 56, 63 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 12, 17, 20, 64, 109, 185, 186, 195, 197, 201, 205, 261, 417 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 153, 158, 161, 194, 207, 259, 278, 310, 329, 362, 378 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_tester.md` (document_read) · 줄 77, 82, 85, 118, 130, 189, 232 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_designer.md` (document_read) · 줄 194 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 47, 57, 69, 70, 73, 73, 73, 77, 78, 78, 83, 83, 92, 92, 108, 115, 120, 121, 121, 122, 129, 131, 139, 139, 139, 144, 287, 297, 297, 299, 301, 305, 307, 309, 313, 314, 315, 316, 319 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_tester.md` (document_read) · 줄 32, 282, 290, 291 · inferred
- `.fullops-squad/orca-agents.md` (document_read) · 줄 37 · inferred
- `adapters/README.md` (document_read) · 줄 143, 165, 167, 179, 188, 189, 190, 191, 196, 204 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 14 · inferred
- `adapters/src/connect.ts` (direct_edit) · 줄 전체/미확인 · unknown
- `adapters/src/mcp.ts` (impact_check) · 줄 149, 212 · inferred
- `adapters/src/trial-check.ts` (impact_check) · 줄 22, 24, 34, 37, 65, 74, 86, 89 · inferred
- `adapters/src/trial-setup.ts` (impact_check) · 줄 1, 22, 35, 44, 50, 53, 61, 63, 68, 72, 73, 89, 99, 103 · inferred
- `cmd/relay/main.go` (impact_check) · 줄 56, 84 · inferred
- `deploy/knowslink/access_trial_plan.py` (impact_check) · 줄 69 · inferred
- `internal/relay/admission_integration_test.go` (impact_check) · 줄 39, 141, 148 · inferred
- `internal/relay/capacity.go` (impact_check) · 줄 28, 53, 120 · inferred
- `internal/relay/cleanup_flood_integration_test.go` (impact_check) · 줄 26, 36, 88, 89, 157, 159, 199 · inferred
- `internal/relay/cleanup_unit_integration_test.go` (impact_check) · 줄 48, 52, 53, 63, 64, 82, 93, 108, 120, 138, 144, 146 · inferred
- `internal/relay/connections.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/connections_test.go` (impact_check) · 줄 62, 64, 89, 93, 101, 105, 145, 156, 235, 238, 269, 278, 291, 295, 335, 335, 338, 338, 342, 360, 381, 392, 395, 397, 402, 404, 405, 407, 407 · inferred
- `internal/relay/google.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/google_test.go` (impact_check) · 줄 1 · inferred
- `internal/relay/member_agents.go` (direct_edit) · 줄 36, 40, 66, 91, 104, 112 · inferred
- `internal/relay/member_receipt.go` (impact_check) · 줄 14, 29, 30, 31, 45 · inferred
- `internal/relay/policy_test.go` (impact_check) · 줄 22, 27, 31, 41, 51, 54, 59, 65, 68, 68, 234, 253, 256 · inferred
- `internal/relay/protocol.go` (impact_check) · 줄 133, 164, 180, 289, 290 · inferred
- `internal/relay/protocol_test.go` (impact_check) · 줄 20, 54 · inferred
- `internal/relay/public_messages_integration_test.go` (impact_check) · 줄 57, 65, 69, 73, 79, 83, 84, 91, 109, 111, 123, 124, 152, 153, 167, 169, 170, 171, 178, 179, 381 · inferred
- `internal/relay/public_text.go` (impact_check) · 줄 70, 72, 87, 160, 176, 190 · inferred
- `internal/relay/test_messages.go` (impact_check) · 줄 18, 20, 21, 22, 24, 32, 36, 110, 114 · inferred
- `internal/relay/testdata/legacy_claims.json` (impact_check) · 줄 230, 252 · inferred
- `scripts/run_trial.py` (impact_check) · 줄 24 · inferred
미확인 50건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

작업 후 전문을 작성한다.

## 전달 스냅샷

기준 HEAD와 worker HEAD는 c3f918dfb89c6c670075cddf7a1cc351d1397bf2다. coordinator의 현재 미커밋 지시서·route/find/context/packet을 동일 경로로 worker에 복사해 전달한다. 원천 정본은 동일 기준의 커밋 파일이다. 실제 전달 시 인박스 SHA-256 동일성을 확인한다. 작업 중 coor 문서가 바뀌어도 worker는 전달받은 자기 인박스와 preamble을 정본으로 사용한다.

## 지시 전제와 충돌 — 먼저 확인

- deploy/knowslink/access_apply.py는 이전 owner-only beta 생성/검증 도구다. 사용자 최신 지시는 일반 사용자에게 Cloudflare 로그인을 요구하지 않는 것이다. 기존 관리자 보호/복구와 새 Google 사용자 경계를 구분하고 기술 인계에서 해결한다. 이전 도구를 무조건 재실행하지 않는다.
- SAR-GOOGLE-LOGIN-001-DEV.md는 Google 회원 로그인만 구현한 이전 과제다. 지금은 그 인증을 재사용하여 클라이언트 연결을 단순화하는 확정 후속이다. 이전 실제 Google 가입·원 검증은 보존하고 token 수동 전달을 새 UX의 전제로 삼지 않는다.
- 위 두 충돌은 최신 사용자 지시와 과거 범위의 차이다. 새 제품 승인을 기다릴 사유가 아니며 기술 구현은 DEV가 같은 과제에서 결정한다.
