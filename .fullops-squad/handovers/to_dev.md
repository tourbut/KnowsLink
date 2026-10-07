---
title: SAR-GOOGLE-LOGIN-001-DEV — Google 첫 로그인 자동 가입을 기존 회원·세션·agent 기능에 연결한다
status: draft
updated: 2026-10-07
owner: dev
tasks: [SAR-GOOGLE-LOGIN-001-DEV]
summary: Google 첫 로그인 자동 가입을 기존 회원·세션·agent 기능에 연결한다
attempt: 58fe2507fee041d797df00bafce55d3d
base: abf2de1e3f0501e56b41f9d909514f591b656d4c
---

# SAR-GOOGLE-LOGIN-001-DEV — Google 로그인으로 자동 회원가입과 기존 기능을 연결한다

- 작성일: 2026-10-07. From coor / To dev. 상태 ready.
- repo 818c78e5-d51c-4ff4-aa88-70e9ee185fbb; /home/shin/orca/workspaces/KnowsLink/fullops-dev; fullops/dev.
- coor 복귀: /home/shin/orca/workspaces/KnowsLink/fullops-coor; term_a8a1fa04-50ab-448d-94e7-11e8ee3c77f1; Run run_8ca8bc058ab7. Task/Dispatch는 실제 주입 preamble을 따른다.
- 승인: 구현·관련 좁은 검사·기술 정본 갱신·일반 역할 push. coor가 main 통합과 기존 서버 배포를 진행한다. 무료만 사용하며 결제/구독/새 유료서비스·사용자 자료 삭제·공유 Tunnel 변경은 금지한다.

## 현재 상황과 확인 근거

사용자가 Google 로그인 하나로 먼저 출시하기로 확정했고 기능 구현에 집중해 빠르게 진행하라고 요청했다. SMTP는 없으므로 이메일 코드가 있어도 실제 가입이 안 된다. Google 첫 로그인에서 일반 회원과 owner를 자동 생성하고 재로그인 시 같은 신원으로 기존 agent/관계/메시지를 사용하게 한다. Google 계정의 비밀번호를 이 서버에 받지 않는다. 기존 제품 코드 d089가 보호 서버에 배포됐고 main abf2de1에 검수 종료 기록까지 통합됐다. 과거 QA/UI 미완료와 UTF8 medium을 이번 과제에서 재실행/확대하지 않는다.

## 적용 기준과 예외

fullops-common-0.3.3 README/coding-style/testing/security, project.md, 문서 작성 규칙과 기준 abf2de1e3f0501e56b41f9d909514f591b656d4c. 사용자 최신 결정이 이메일 코드 전용 가입과 전체 QA 대기보다 우선한다. 넓은 인증 프레임워크·다중 IdP·메일서버를 만들지 않는다. 알려진 권한 노출·계정 탈취·데이터 손실 문제는 생략하지 않는다. 기술 계획과 구현/검사는 같은 DEV 과제다.

## 먼저 읽을 문서

- .fullops-squad/FULLOPS.md
- .fullops-squad/rules/common/README.md
- .fullops-squad/rules/common/coding-style.md
- .fullops-squad/rules/common/testing.md
- .fullops-squad/rules/common/security.md
- .fullops-squad/project.md
- .fullops-squad/contexts/dev.md
- .fullops-squad/docs/agents/document-writing.md
- .fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md (PS01/02/04와 이번 사용자 Google 선택을 구분)
- internal/relay/member.go
- internal/relay/store.go
- cmd/relay/main.go
- compose.yaml
- .fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md

공식 근거: https://developers.google.com/identity/openid-connect/openid-connect 를 coor가 2026-10-07 직접 확인했다. 프로젝트/OAuth client·redirect URI·ID token 검증·state·issuer/sub 신원 연결이 필요하다. SDK를 새로 쓸 때 Context7이 있으면 사용하고 없으면 공식 문서를 확인한다. Google 콘솔 실제 등록과 자격은 coor 담당이며 현재 로그인 화면에서 사용자 로그인 대기다. 비밀값을 chat/Git/argv/log에 쓰지 않는다.

## 해야 할 일과 파일 소유권

- [ ] 기존 가입·로그인·세션·재확인·owner 생성과 모든 호출자를 좁게 읽고 같은 과제에 짧은 기술 계획을 남긴다.
- [ ] Google 로그인 버튼과 안전한 로그인 흐름을 구현한다. 새 회원 생성·같은 Google 신원 재로그인·기존 owner 기능을 연결한다. 이메일이 같다는 이유만으로 이메일 코드 회원이나 다른 발급자의 회원을 자동 병합하지 않는다. 인증된 Google issuer/sub를 안정적인 신원 키로 사용한다. 권한 작업의 최근 인증도 Google로 동작해 SMTP가 다시 필요해지지 않게 한다.
- [ ] 기존 cookie/CSRF/기간/철회·로그아웃·회원/agent 한도와 권한을 재사용한다. 설정 미완료/Google 취소·실패는 안전하게 안내하고 회원 생성/로그인 성공으로 표시하지 않는다. 세션 고정·state/replay·token 서명/issuer/audience/expiry/검증 이메일 조건을 놓치지 않는다.
- [ ] 환경 설정과 compose 전달·README 실제 가입 사용법을 연결한다. coor에게 정확한 redirect URI·Google 콘솔 값·필요 env 이름을 가능한 한 초기 status로 알린다. 비밀 없이 로컬/운영 경로를 구분한다.
- [ ] 영향받은 작은 기술 정본과 실행 기록만 갱신하고 최소 runnable 로그인/거부 회귀·빌드·lint를 확인한다. work.py finish 전문 archive/빈 inbox·최종 SHA·역할 일반 push와 authentic worker_done을 완료한다.

소유권: DEV 제품 경로 internal/cmd/compose/README/.env.example/Go module 및 영향받은 docs/design-docs 기술 정본. 기획 정본이나 OPS 운영 문서의 의미 변경은 직접 하지 말고 coor에 짧은 변경 필요만 전달한다. 기존 UI template/style 재사용; 새로운 대시보드/테마/컴포넌트 라이브러리 없음. 예상 코드/회귀 3–8파일, 제품 +200–350줄 이내를 목표로 하되 실제 안전한 구현에 필요한 규모는 근거를 남긴다.

## 완료 기준과 검증

Google-only 설정에서 메일러 없이 가입·로그인·로그아웃·권한 재확인이 가능하며 같은 Google 신원은 중복 회원/owner를 만들지 않는다. 실패/위조/타인 신원은 회원 권한을 만들지 않는다. Google ID token 확인을 stdlib로 재발명하지 말고 기존/검증된 작은 라이브러리를 먼저 검토한다. 실제 Google 자격이 없으면 코드와 재현 가능한 좁은 검사·빌드를 마치고 외부 실로그인은 미검증으로 보고한다. 자격 대기를 이유로 코드를 중단하지 않는다.

DEV가 변경 인증의 정상/거부/재로그인/최근인증 연결만 검증한다. 전체 verify-mvp·전체 UI 캡처·부하/장시간/외부메일 검증을 추가하지 않는다. 이미 통과한 검사를 반복하지 않는다. 기존 required lint/test 명령이 호출되면 결과를 재사용하고 필요한 1회만 수행한다. 후속 검토는 인증 delta의 고정 SHA에 한정하고 일반 서비스 전체 재검수는 배포 선행조건이 아니다. 실제 Google 브라우저 로그인 한 경로는 coor와 사용자 계정으로 확인한다. snapshot/UI 대량 캡처 불필요. 새 critical/high가 발견되면 같은 DEV 과제에서 수정한다.

## 갱신할 산출물

route 추천 D10/D12/D13 중 DEV는 D10 영향 절만 갱신한다. D12/D13 설정/배포 변경은 정확한 env/redirect·코드 SHA 인계로 coor가 반영한다. 실제 기술 영향에 따라 D03/D05/D06 중 필요한 절만 갱신하며 전체 13종 재작성은 금지한다. 원 이메일 기반 기획의 변경 필요는 coor에 알린다.

## 기대 산출물·제약

구현·좁은 회귀·README/환경 안내·docs/exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md. 회원 생성/agent API의 기존 권한 계약 유지. 비공개 Google client 자격과 실제 메일주소는 로그에 넣지 않는다. 원 실패/보존 자료를 변경하지 않는다. Google 등록의 확정 경로·설정은 coor로 조율하고 구현은 계속 진행한다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad/docs/evaluations/jev/SAR-GOOGLE-LOGIN-001-DEV-packet.json` / SHA `3b183a5304734d4c50f60d3e26da23e057fd9971` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 62 · inferred · 필수 · {'relevant': 0.26, 'evidence': 0.65, 'contradicts': 0.21, 'injection': 0.06, 'decision': 'keep', 'reason': None}
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 29 · inferred · 필수 · {'relevant': 0.26, 'evidence': 0.49, 'contradicts': 0.52, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/agents/document-writing.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.14, 'evidence': 0.38, 'contradicts': 0.13, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (document_read) · 줄 63 · inferred · 필수 · {'relevant': 0.08, 'evidence': 0.12, 'contradicts': 0.64, 'injection': 0.05, 'decision': 'keep', 'reason': None}
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.23, 'evidence': 0.32, 'contradicts': 0.79, 'injection': 0.07, 'decision': 'keep', 'reason': None}
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 2, 6, 6, 7, 12, 12, 12, 17, 21, 21, 37, 40, 44, 49, 50, 51, 58, 60, 68, 68, 68, 72 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 15 · inferred · 필수 · {'relevant': 0.47, 'evidence': 0.88, 'contradicts': 0.21, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.11, 'evidence': 0.35, 'contradicts': 0.18, 'injection': 0.06, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.24, 'evidence': 0.62, 'contradicts': 0.1, 'injection': 0.02, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.76, 'evidence': 0.83, 'contradicts': 0.1, 'injection': 0.07, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수 · {'relevant': 0.36, 'evidence': 0.7, 'contradicts': 0.12, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `.fullops-squad/contexts/coor.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/contexts/ops.md` (document_read) · 줄 63 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 15, 59 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 93 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 25, 29, 97, 118 · inferred
- `.fullops-squad/docs/design-docs/module-design.md` (document_read, document_update) · 줄 82, 160 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 44, 53, 59, 76 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-FINAL-review/report.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-N1-review/report.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-review/report.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-TESTER.md` (document_read) · 줄 30, 71, 84, 91 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-002-TESTER.md` (document_read) · 줄 17, 36 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-DEPLOY-001-OPS-FINAL-review/report.md` (document_read) · 줄 29, 34, 61 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FINAL-review/report.md` (document_read) · 줄 22, 151 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FINAL.md` (document_read) · 줄 86 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX.md` (document_read) · 줄 27 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md` (document_read) · 줄 25 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-FIX-TESTER.md` (document_read) · 줄 130 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-TESTER.md` (document_read) · 줄 124 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-review/report.md` (document_read) · 줄 58 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER.md` (document_read) · 줄 27, 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-review/report.md` (document_read) · 줄 37 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-FINAL-RECORDS-review/report.md` (document_read) · 줄 40, 42 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-FIX-TESTER.md` (document_read) · 줄 37 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md` (document_read) · 줄 53, 99, 100 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-FIX-TESTER.md` (document_read) · 줄 38 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW-review/report.md` (document_read) · 줄 38, 78, 131, 148, 164 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-REVIEW-review/report.md` (document_read) · 줄 31 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TESTER.md` (document_read) · 줄 39 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-REVIEW-review/report.md` (document_read) · 줄 29, 114 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER.md` (document_read) · 줄 32 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX-2.md` (document_read) · 줄 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 28, 39 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-DEV-099-review/report.md` (document_read) · 줄 43 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-INTEGRATION-FINAL-review/report.md` (document_read) · 줄 23, 31, 39 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-INTEGRATION-SOURCE-review/report.md` (document_read) · 줄 7, 13, 23, 32, 34, 45, 47, 63, 71 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-TESTER.md` (document_read) · 줄 24, 31, 79, 100, 101 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-BETA-001-TESTER.md` (document_read) · 줄 13, 52 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER-FINAL.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER-FIX.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER.md` (document_read) · 줄 13 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-SETUP-001-TESTER.md` (document_read) · 줄 19, 29 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.10.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.12.md` (document_read) · 줄 14, 22, 32, 36, 38, 42, 48 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.13.md` (document_read) · 줄 26, 32 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.14.md` (document_read) · 줄 25, 39, 44, 49, 51, 61 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-099.md` (document_read) · 줄 15, 41, 46 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-BETA-001-OPS.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-DEPLOY-001-OPS.md` (document_read) · 줄 22, 36, 44 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-001-DEV.md` (document_read) · 줄 153, 158, 159, 170 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md` (document_read) · 줄 47, 70 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT.md` (document_read) · 줄 64 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md` (document_read) · 줄 14, 110 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-TRIAL-CLEANUP.md` (document_read) · 줄 66 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PREP-002.md` (document_read) · 줄 40 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX.md` (document_read) · 줄 91, 109 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG.md` (document_read) · 줄 22, 41, 43 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW.md` (document_read) · 줄 29, 36 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW.md` (document_read) · 줄 28 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-UI.md` (document_read) · 줄 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-2.md` (document_read) · 줄 93 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX-3.md` (document_read) · 줄 88 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV-FIX.md` (document_read) · 줄 76 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-UI-FIX.md` (document_read) · 줄 38, 50 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (document_read) · 줄 40, 80 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPEN-PREP.md` (document_read) · 줄 2, 10, 14, 18, 27 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SETUP-001-TESTER.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/docs/operations/ops-guide.md` (document_read) · 줄 17, 47, 49, 300 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 72 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/planning/SAR-PREP-002-request.md` (document_read) · 줄 14, 16 · inferred
- `.fullops-squad/docs/planning/SAR-SETUP-001-request.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/handovers/SAR-MVP-001-REVIEW.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/handovers/_TEMPLATE.md` (document_read) · 줄 18, 80 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_designer.md` (document_read) · 줄 28, 104, 110 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_ops.md` (document_read) · 줄 68, 84, 94, 115, 148, 213, 219 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_dev.md` (document_read) · 줄 27, 90, 149, 196, 197, 259, 275, 294 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_ops.md` (document_read) · 줄 24, 72, 161, 197, 229, 237, 254, 269, 270, 299 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_tester.md` (document_read) · 줄 24, 75, 125, 167, 207 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_designer.md` (document_read) · 줄 93, 156 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_dev.md` (document_read) · 줄 43, 93, 154 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_ops.md` (document_read) · 줄 22, 39, 64, 105, 159 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_tester.md` (document_read) · 줄 85 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 32, 235, 257, 276, 284 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 28, 73, 77, 108, 173, 203, 232, 254, 280, 283, 336, 384 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_tester.md` (document_read) · 줄 36, 97, 142, 162, 184, 207, 240, 252, 254 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_designer.md` (document_read) · 줄 29, 226, 228 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 37 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-DEV-REVIEW.md` (document_read) · 줄 37 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-INTEGRATION-REVIEW.md` (document_read) · 줄 7, 12, 13, 22, 22, 26, 28, 32, 36, 40 · inferred
- `.fullops-squad/lint/README.md` (document_read) · 줄 48, 49, 61, 62, 85 · inferred
- `.fullops-squad/orca-agents.md` (document_read) · 줄 18, 19, 20, 21, 24, 29, 40, 99 · inferred
- `Makefile` (impact_check) · 줄 17, 33, 43 · inferred
- `README.md` (document_read) · 줄 전체/미확인 · unknown
- `adapters/package.json` (impact_check) · 줄 14, 15, 17 · inferred
- `adapters/src/index.ts` (impact_check) · 줄 6, 25 · inferred
- `adapters/src/public-check.ts` (impact_check) · 줄 91, 93, 99 · inferred
- `adapters/src/trial-check.ts` (impact_check) · 줄 12, 113 · inferred
- `adapters/src/trial-cli.ts` (impact_check) · 줄 3, 33 · inferred
- `adapters/src/trial-setup.ts` (impact_check) · 줄 78, 124 · inferred
- `cmd/migrate/main.go` (impact_check) · 줄 2, 18, 19, 25, 26 · inferred
- `cmd/migrate/main_test.go` (impact_check) · 줄 2, 14, 21 · inferred
- `cmd/relay/main.go` (direct_edit, impact_check) · 줄 19, 28, 95 · inferred · {'relevant': 0.18, 'evidence': 0.36, 'contradicts': 0.39, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `cmd/relay/main_test.go` (impact_check) · 줄 2, 27 · inferred
- `compose.yaml` (direct_edit, impact_check) · 줄 전체/미확인 · unknown
- `deploy/knowslink/access_trial_plan.py` (impact_check) · 줄 49, 91 · inferred
- `deploy/knowslink/compose.ops.yaml` (impact_check) · 줄 32 · inferred
- `deploy/knowslink/verify.py` (impact_check) · 줄 15, 16, 28, 29, 32, 54, 56, 59 · inferred
- `go.mod` (direct_edit) · 줄 전체/미확인 · inferred
- `internal/config/config.go` (direct_edit) · 줄 13, 25 · inferred
- `internal/relay/identity.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/member.go` (direct_edit, impact_check) · 줄 전체/미확인 · unknown
- `internal/relay/store.go` (impact_check) · 줄 전체/미확인 · unknown · {'relevant': 0.07, 'evidence': 0.1, 'contradicts': 0.19, 'injection': 0.03, 'decision': 'keep', 'reason': None}
- `scripts/install_bot_mcp.sh` (impact_check) · 줄 25, 35, 60, 64 · inferred
- `scripts/mail_sink.py` (impact_check) · 줄 38, 53 · inferred
- `scripts/package_plugin.py` (impact_check) · 줄 22, 29, 57, 64 · inferred
- `scripts/run_trial.py` (impact_check) · 줄 9, 33 · inferred
- `scripts/verify_grok_plugin.py` (impact_check) · 줄 13, 25, 38, 43 · inferred
- `scripts/verify_setup.py` (impact_check) · 줄 13, 14, 23, 24, 27, 36, 45, 48, 53 · inferred
미확인 49건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified

### 패킷 해석과 최신 사용자 결정
- context의 옛 PS01/02·UX 이메일 전용과 dev 과거 상태는 현재 사용자 Google 선택보다 이전이다. 충돌 원문을 보존하되 최신 사용자 확정으로 Google 로그인 전환을 구현한다. 추가 제품 승인 대기는 없다.
- find가 확인한 internal/relay/identity.go·internal/config/config.go를 member.go와 먼저 읽는다. 필수 문서와 이 핵심 경로의 초기 입력은 20개 이하로 유지한다.
- producer partial/unknown은 전체 부재가 아니다. 검색 영향 목록은 범위 확대 지시가 아니다. 관련성을 좁게 판단해 unchanged/no_change의 실제 근거를 outcomes에 적고 무관한 과거 산출물·QA를 다시 실행/작성하지 않는다. 공통 store 호출자 중 인증 변화 영향만 확인한다.

<!-- fullops-packet:end -->

## 완료 보고

실제 final SHA·변경 목적·최소 검증 결과·실제 Google 로그인 미검증 여부·정확한 운영 설정/redirect·남은 차단을 전문 보고하고 work.py finish 뒤 직접 worker_done을 보낸다.

## DEV 기술 계획

1. `coreos/go-oidc/v3` v3.17.0과 `oauth2` v0.36.0으로 Google code 교환과 서명·issuer·audience·expiry를 검증한다. Context7은 quota 오류이므로 Google 및 라이브러리 공식 문서와 해당 버전 소스를 확인했다.
2. 기존 JSON state에 10분짜리 Google 시도만 추가한다. state 해시·nonce·PKCE·현재 세션 해시를 저장한다. callback은 브라우저의 Secure/HttpOnly/Lax 시도 cookie와 state를 대조하고 1회 소비한다. 기존 세션 cookie는 Strict를 유지한다.
3. issuer/sub 기반 회원·owner 생성과 기존 세션 발급을 공유한다. 재확인은 현재 Google 회원과 동일한 신원만 허용한다. 최근 발급된 nonce token과 계정 선택·동의로 신원을 다시 확인한다. Google 비밀번호 재입력은 요구하거나 보장하지 않는다.
4. 기존 템플릿·CSRF·rate·한도를 유지한다. Google-only 환경에서 이메일 폼을 숨기고 Google 회원의 재확인을 Google로 보낸다. callback 완료 화면의 홈 링크는 Strict 세션이 cross-site redirect chain에서 누락되는 문제를 피한다.
5. 인증 정상·거부·재로그인·최근인증과 설정 검사를 추가한다. required lint/test 1회와 Go 빌드 및 좁은 Postgres HTTP 확인 후 D10·README·env를 인계한다.
