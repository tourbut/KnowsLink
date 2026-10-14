---
title: 공개 안내·로그인 진입 구현과 검증
status: draft
updated: 2026-10-14
owner: dev
tasks: [SAR-LANDING-GUIDE-001-DEV]
summary: 공개 안내와 설정별 인증 복구를 구현하고 실제 제품 게이트 및 격리 브라우저를 인계한다
---

# 공개 안내 구현

## 기준과 기술 계획

실제 task는 task_f0d345aa84d2, dispatch는 ctx_65212ab451f7이다. 착수 HEAD는 b0f6596b6b7bfd4ea4e3b8cfed6d8b1e20a9b93e다. 비교 기준은 eae9a8bffae706410f684e4c837e396e7b0bcf26이다. 설계 독립 리뷰와 main/origin 통합 완료를 받은 실제 dispatch로 queued에서 running으로 전환했다.

공식 jev_packet.py로 startup packet을 갱신하고 check(required=True)를 통과했다. 원 지시 fingerprint는 e1b0d0c5e21d50b0e9751d263794111d503a818ee796e28807b691d6d35479e6, startup packet fingerprint는 497a85ed4167c3fd06ef93ec7c311293b9ad254c2287fcc1f038315c6cbfe612다. 과거 producer의 HEAD·지시 불일치와 partial/unknown을 보존했다. 원문 검토와 좁은 검색으로 보완한다.

1. startView를 Service 메서드로 바꿔 모든 start 오류 호출자가 실제 Google/Mail 활성 상태를 공유한다. startPage는 유효 세션 여부만 view에 전달한다.
2. 공개 전용 landingTemplates를 기존 memberPages에 연결한다. 기존 memberStyle의 팔레트·타입·네이티브 폼/details를 재사용한다. 공개 hero·h2·최대 폭·구획 간격만 D04 예외를 적용한다.
3. 회원 내 계정에 서비스 안내 링크를 추가한다. 인증 POST·원 redirect·Google 재확인·Remote 동의·회원 이름 설정은 유지한다.
4. 상태별 handler/template와 실제 DB 경계를 확인한다. 고정 SHA의 make lint/test/build, 격리 DB·서버와 브라우저 증거를 제공한다.

변경 범위는 member.go·google.go의 view 호출자, member_pages.go·landing_pages.go·human_pages.go, 관련 테스트와 fixture, D10/D11/D04 구현 연결 및 DEV 기록이다. 새 schema·프레임워크·패키지·클라이언트 프로토콜은 없다. 단일 공유 view와 template 변경이라 선택 하위 위임은 사용하지 않았다.

## 디자인·사실 대조

D04 공개 첫페이지 절의 카피·다섯 구획·순서·접기·예시·로그인 상태가 정본이다. README의 Grok 신규/업데이트, Dots 신규/업데이트 제목과 installer의 Command/Arguments·신규 폴더 조건을 대조했다. 수신 process·Events·화면 표시·발송 승인·180초·24h는 docs/architecture.md, docs/system-design.md, docs/operations.md와 대조했다. 과거 운영 성공·실패·미검증을 새 페이지의 성공 보장으로 바꾸지 않는다.

실행 선택은 codex/gpt-6-astra/high와 frontend-design이다. frontend.py check는 configuration_ready다. 기존 frontend.json의 Sol 기본값은 변경하지 않았다. 실제 인증·브라우저 검증과 구성 확인은 별개다. 새 도구·외부 리소스·테마는 추가하지 않는다. 별도 디자인 lint와 다중테마는 없으며 Go 문자열 CSS는 일반 DESIGN 규칙의 검사 대상이 아니다.

## 구현과 검사 중 발견

기존 email 오류의 startView는 EmailLogin=true를 고정했다. Service 설정을 공유하도록 바꿔 Google-only/none 오류가 비활성 이메일 폼을 노출하지 않게 한다. DB·세션·권한 로직은 유지한다.

기존 member form golden·Google·human 단위 회귀는 초기 구현에서 통과했다. 새 landing 단위 검사는 nil DB로 전체 Service.Handler를 호출해 전역 rate wrapper에서 panic했다. 이는 테스트 준비 결함이다. startPage 직접 호출과 실제 DB integration의 cache/CSP 검사를 분리한다. 첫 실패와 후속 결과는 QA에 보존한다.

비파괴 테스트 편집의 inline Python 셸 명령이 자동 승인 검토에서 `blocked by policy`로 거부됐다. 구체 이유는 없었다. 우회하지 않고 coordinator ask로 apply_patch 수정 허가를 요청했다.

coor 답변 msg_d75b88015752는 원 차단 call_a0b84ac8954c41208d0df30778fa2355와 두 테스트를 직접 확인했다. 셸 재실행 없이 명시적 apply_patch 수정을 허용했다. 익명 GET과 유효세션 GET 모두 실제 DB integration에 cache/CSP assertion을 유지했다. 중복 폼과 활성 여부 assertion도 유지했다. 수정 후 gofmt와 landing/member/Google/human 관련 단위 회귀는 종료0이다.

## 후속

DEV 자체 검사 완료 뒤 독립 tester QA, designer 직접 시각 판정, 별도 세션 고정 SHA 리뷰가 필요하다. 운영 배포·실계정 Google/메일·실상대 메시지·설치 실행은 범위 밖이다. PLANS/board와 기본 브랜치 통합은 coor 소유다. 임시 검증 자료는 보존한다.
