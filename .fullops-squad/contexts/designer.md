---
title: designer 컨텍스트
status: draft
updated: 2026-10-14
owner: designer
tasks: [SAR-SETUP-001, SAR-PREP-002, SAR-MVP-PUBLIC-POLICY-001, SAR-MVP-001-UI, SAR-PUBLIC-SERVICE-001, SAR-PUBLIC-AGENTS-001-UI, SAR-PUBLIC-AGENTS-001-POLICY, SAR-PUBLIC-AGENTS-001-UI-FIX, SAR-PUBLIC-MESSAGES-001-UI, SAR-PUBLIC-MESSAGES-001-UI-FIX, SAR-PUBLIC-MESSAGES-001-UI-FIX-2, SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL, SAR-SERVICE-UI-001-DESIGN, SAR-SERVICE-UI-001-VISUAL, SAR-SERVICE-UI-001-FIX-VISUAL, SAR-HUMAN-PEER-001, SAR-HUMAN-PEER-001-VISUAL, SAR-LANDING-GUIDE-001]
summary: 제품 판단과 직접 UI 검수의 수락·중단·보류 경계를 보존한다
---

# designer 컨텍스트

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
