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

## 완료 보고

실제 final SHA·변경 목적·최소 검증 결과·실제 Google 로그인 미검증 여부·정확한 운영 설정/redirect·남은 차단을 전문 보고하고 work.py finish 뒤 직접 worker_done을 보낸다.
