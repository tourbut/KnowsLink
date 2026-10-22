---
title: designer 컨텍스트
status: draft
updated: 2026-10-22
owner: designer
tasks: [SAR-SETUP-001, SAR-PREP-002, SAR-MVP-PUBLIC-POLICY-001, SAR-MVP-001-UI, SAR-PUBLIC-SERVICE-001, SAR-PUBLIC-AGENTS-001-UI, SAR-PUBLIC-AGENTS-001-POLICY, SAR-PUBLIC-AGENTS-001-UI-FIX, SAR-PUBLIC-MESSAGES-001-UI, SAR-PUBLIC-MESSAGES-001-UI-FIX, SAR-PUBLIC-MESSAGES-001-UI-FIX-2, SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL, SAR-SERVICE-UI-001-DESIGN, SAR-SERVICE-UI-001-VISUAL, SAR-SERVICE-UI-001-FIX-VISUAL, SAR-HUMAN-PEER-001, SAR-HUMAN-PEER-001-VISUAL, SAR-LANDING-GUIDE-001, SAR-UX-ISSUES-002-DESIGN, SAR-SESSION-PEER-001-DESIGN, SAR-COMPETE-001]
summary: 제품 판단과 직접 UI 검수의 수락·중단·보류 경계를 보존한다
---

# designer 컨텍스트

- 2026-10-15 SAR-SESSION-PEER-001-DESIGN: 명시 선택한 작업 세션의 수신을 다음 조회·새 세션 실행과 구별한다. HP-07/10/15·legacy·자원 한도를 유지하며 기본 자동 답장은 0이다.
  Claude Channels·Codex remote/App Server/queue는 검증 후보이며 Hermes logging callback은 로그 처리다. 실제 연결·왕복·운영은 미검증이다.
  [제품 규칙과 DEV 인계](../docs/design-docs/SAR-SESSION-PEER-001.md), [검증 보고](../docs/evaluations/qa-reports/SAR-SESSION-PEER-001-DESIGN/report.md). 조회 복구 수락 뒤 세션 통신을 우선하고 나머지 UI 이슈는 대기한다.
  2026-10-17 Claude로 재개해 완성했다. 세션 선택기는 U2 이름 선두를 기본으로 하고 좁은 닫힌 선택기만 U6 구분값 선두 예외를 허용한다.

전문이10KB를 넘어 원문을 [2026-10-14 보존본](archive/designer-2026-10-14-landing.md)에 복사했다. 원 파일을 삭제하거나 이동하지 않았다. 아래는 다음 과제에 필요한 현재 원칙이다. 과거 FAIL·한정 PASS·미검증의 전문과 링크는 보존본을 따른다.

- 제품 기획은 designer, 기술 계획·구현·API·버그 해결은 DEV 책임이다. 실제 화면 검수는 고정 SHA와 합성/실계정 범위를 구분한다.
  작업 완료·문서 검사·동작 QA·직접 시각 수락·운영 수락은 서로 다르다. 과거 PASS는 새 화면의 수락 근거가 아니다.
  공통 기준: [FULLOPS](../FULLOPS.md), [프로젝트](../project.md).

- 사람 이름 아래 agent가 연결되는 구조와 대화/친구/내 에이전트 주 이동을 유지한다. 친구 수락은 현재·향후 활성 agent 범위 동의이며 매 발송 승인을 대신하지 않는다.
  legacy pair의 정확한 범위·철회/차단·키·신원·24h 웹 사본과180초 전달 기한을 보존한다. 이름 변경은 ID나 권한 변경이 아니다.
  정본: [D02](../docs/planning/product-specs/SAR-HUMAN-PEER.md), [D04](../docs/design-docs/mockups/service-member-ui.md).

- SAR-HUMAN-PEER-001-VISUAL 원254ff6a의 FAIL은 보존한다. 후보12dc0c7의 제공된 V01–V03 화면만 resolved/시각 PASS다.
  실제 keyboard/200%·미제공 상태는 not_run이었다. 새 공개 안내 QA의 통과로 간주하지 않는다.
  [원 보고](../docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/report.md), [수정 후보](../docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/final-report.md).

- 2026-10-14 SAR-LANDING-GUIDE-001은 최신 요청에 따라 공개 첫페이지를 소개·순서·Grok/Dots 설치·가상 예시·로그인으로 확장한다. 유효 세션도 안내를 읽고 대화로 이동한다.
  D02/D04/D11과 Astra high DEV/tester 인계를 작성했다. Grok 접기에 준비물·installer 조건·Command/Arguments·신규 설정을 제공하며 업데이트는 분리한다. 새 제품/시각 수락은 아니다.
  결정·검증·후속: [실행 기록](../docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md). 원 partial/unknown·최소 로그인 지시·운영 미검증을 보존한다.

- 2026-10-14 SAR-LANDING-GUIDE-001-VISUAL: 후보fe2387a/fixture4795의 공개 화면은 직접 시각 PASS이며 예시 뒤 간격0px V-LG-01 low가 남는다.
  회원 대표 화면의 기존 CTA 대비1.846967:1은 V-LG-02 medium FAIL이다. designer native 입력 미실행과 tester의 별도 증거, native200% 대기를 구분한다. 전체 수락은 아니다.
  [직접 시각 보고](../docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/report.md). 제품/D04를 바꾸지 않았고 공유 fixture는 coor/tester가 유지한다.

- 2026-10-15 SAR-UX-ISSUES-002-DESIGN: #2–#8을 사람 친구 규칙으로 대조했다. handle/초대 링크 대신 친구 코드, 조회 복구·검증된 승인 요약·언어 저장·설치 재개 기준을 확정했다.
  기존 D04 공개 규칙·제품 코드·타 역할 인박스는 보존했다. Grok/Dots 기존 계정 관측을 신규 설치 완료로 확대하지 않는다. Meta Muse/Hermes/OpenClaw는 지원 확인 필요다.
  [설계·순차 DEV/QA 인계](../docs/design-docs/SAR-UX-ISSUES-002.md), [검증 기록](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN/report.md). 구현·새 UI 수락·운영은 후속이다.

- 2026-10-15 SAR-LANDING-GUIDE-001-VISUAL CSS 후속: source f2b8cb6/final fc4ad58의 PNG7 직접 검수로 V-LG-01/02는 resolved / 시각 PASS다. 원 실패·수치·PNG는 보존한다.
  DEV38 checks·24px·7.958:1과 Linux gate 종료0은 동일성 확인 후 재사용이다. 새 native keyboard/200%·전체 제품 수락은 아니다.
  [후속 보고](../docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/css-followup/report.md). 코드 리뷰·tester 취합·main 통합은 coor 책임이다.

- 2026-10-22 SAR-PUBLISH-001-VISUAL: 정책 페이지·첫페이지·설정 KO/EN 3폭을 CDP 고정 폭 Chrome으로 직접 검수해 PASS(medium 1·low 3). 전역 `a{inline-flex;min-height:44px}`가 문장 속 링크 줄 간격을 벌리는 것이 새 인라인 링크에서 드러났다.
  정책 하단 '홈'은 서비스 약관 뒤 승인, D04 배치도에는 누락(D-01). [보고](../docs/evaluations/qa-reports/SAR-PUBLISH-001-VISUAL-designer-test/report.md)

- 2026-10-22 SAR-COMPETE-001: 이슈 #9 Reason Network는 팀 코딩 세션 버스, KnowsLink는 사람 1:1 수락·발송마다 승인·서버 기본 거부·짧은 보존의 개인 메신저로 구분했다. 기능 수 경쟁 대신 신뢰 경계를 차별점으로 둔다.
  P0 루프 방지 불변식·SoT 포지셔닝, P1 Codex 수신·자동 설정·서명≠암호화 고지, 첨부·그룹 P3, P2P·팀 초대·승인 생략·실시간 활동 공유는 non-goal. Q1–Q8 사용자 결정 전 제안이다.
  [경쟁 비교 D02](../docs/planning/product-specs/SAR-COMPETE-001.md). 경쟁사 수치는 이슈 본문만 쓰고 이슈의 미확인 3건은 미확인으로 유지한다.
