---
title: designer 컨텍스트
status: draft
updated: 2026-10-14
owner: designer
tasks: [SAR-SETUP-001, SAR-PREP-002, SAR-MVP-PUBLIC-POLICY-001, SAR-MVP-001-UI, SAR-PUBLIC-SERVICE-001, SAR-PUBLIC-AGENTS-001-UI, SAR-PUBLIC-AGENTS-001-POLICY, SAR-PUBLIC-AGENTS-001-UI-FIX, SAR-PUBLIC-MESSAGES-001-UI, SAR-PUBLIC-MESSAGES-001-UI-FIX, SAR-PUBLIC-MESSAGES-001-UI-FIX-2, SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL, SAR-SERVICE-UI-001-DESIGN, SAR-SERVICE-UI-001-VISUAL, SAR-SERVICE-UI-001-FIX-VISUAL, SAR-HUMAN-PEER-001, SAR-HUMAN-PEER-001-VISUAL, SAR-LANDING-GUIDE-001]
summary: 제품 판단과 직접 UI 검수의 수락·중단·보류 경계를 보존한다
---

# designer 컨텍스트

결정·교훈을 항목당 3줄 이내로 기록한다.

- 2026-10-12: SAR-HUMAN-PEER-001은 사람 친구의 상호 수락 아래 현재·향후 agent와 연속 대화, 최대 24h 원문 보존을 정했다. legacy pair는 정확한 기존 범위만 보존하며 자동 친구 확대를 금지한다.
  D04를 대화·친구·내 에이전트 중심으로 전면 개편했다. 과거 UI PASS는 새 구조 수락이 아니며 구현·독립 QA·실제 시각 검수는 후속이다.
  제품 결정·이행·DEV/QA 인계와 검증 한계: [실행 기록](../../docs/exec-plans/phases/SAR-HUMAN-PEER-001.md).

- 2026-10-03: SAR-SETUP-001은 초기 구성·lint 과제다. 첫 기능 MVP의 human-gate 포함 결정은 유지하며 이번 범위로 당기지 않는다.
  원천의 확정 기술 선택과 project.md의 이전 미정 상태를 구분한다. 기술 세부 결정은 dev에게 맡긴다.
  상세 근거: [SAR-SETUP-001 실행 기록](../../docs/exec-plans/phases/SAR-SETUP-001.md).

- 2026-10-03: SAR-PREP-002는 최신 7bc9ea1의 개발 준비다. D01·새 MVP D02·기능 백로그와 queued 인계를 작성하고 초기 setup·원천·제품 코드는 보존한다.
  첫 기능은 로컬 등록·수락·안전 전달·human-gate다. 무정책 일정 공개는 deny이며 실데이터 silent 성공·가격/quota·운영 공개는 담당과 재개 조건을 가진 보류다.
  상세 결정·검증·인계: [SAR-PREP-002 실행 기록](../../docs/exec-plans/phases/SAR-PREP-002.md).

- 2026-10-03: SAR-MVP-PUBLIC-POLICY-001은 link.knowslog.com의 누구나 가입 가능한 인증 합성 파일럿 기준이다. 가입 초대 제한을 추가하지 않으며 pairing 수락은 유지한다.
  DEC-03은 수용량/rate/size/concurrency 권장안 단계다. 신규 수락과 안전 정리 budget을 분리하고 DEV/OPS 근거 및 사용자 수치 승인 전 공개 held를 유지한다.
  결정·검증·재개 조건: [공개 정책 기록](../../docs/exec-plans/phases/SAR-MVP-PUBLIC-POLICY-001.md).

- 2026-10-03: SAR-MVP-001-UI의 기존 PNG 7개와 보완 PNG 1개를 직접 확인했다. 합성 V-01–04 시각 판정은 PASS다.
  `invalid_auth`는 인증 실패 차단만 증명한다. `deliver:human` 인증 경계 high는 DEV 수정·독립 재검증 전 유지한다.
  상세 관찰·재개·인계: [UI 실행 기록](../../docs/exec-plans/phases/SAR-MVP-001-UI.md).

- 2026-10-05: SAR-PUBLIC-SERVICE-001은 일반 서비스 기반 수락 뒤 동일 이메일의 Grok Bot “노우”↔OpenAI dot “다닷” 운영 시험이다. agentID·키·credential과 관계 수락은 각각 유지한다.
  새 운영 기본값·복구 목표는 이번 위임 결정이다. 과거 DEC-03 미승인과 실일정·상품 held는 소급하지 않는다. 인증/API/DB/MCP의 기술 선택은 DEV/OPS 책임이다.
  상세 판단·OPS 근거·기능 인계·검증: [일반 서비스 실행 기록](../../docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md).

- 2026-10-06: SAR-PUBLIC-AGENTS-001-UI의 d1eef9b 연결·키·관계를 격리 fixture와 실제 브라우저에서 직접 확인했다. 모바일 지문 overflow와 오류 뒤 복귀 부재로 UX 시각 수락은 FAIL/보류다.
  제품 코드·제품 규칙과 사용자 자격·운영 자원은 보존했다. Orca blank PNG는 실패 근거이며 실제 Playwright PNG와 구분한다.
  상세 관측·마스킹·미검증·DEV 후속: [직접 검수 보고](../../docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md).

- 2026-10-06: SAR-PUBLIC-AGENTS-001-POLICY는 pending/active 반복의 새 초대를 금지하고 거절·만료·양측 철회 뒤 수동 재초대를 허용한다. 새 세대와 수신 owner의 새 수락 전 메시지 거부를 유지한다.
  차단·쿨다운·새 수치는 추가하지 않았다. 수신측 반복 노출 위험과 원본 UI FAIL/보류·실제 운영 후속은 유지한다.
  제품 답·DEV/QA 관찰 조건·검증: [POLICY 실행 기록](../../docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md).

- 2026-10-06: 같은 POLICY 추가 지시에서 키 기록 포화는 새 agent 연결·새 관계 수락으로 처리한다. owner 기록 포화는 필요한 철회·최소 24h 보존·실제 정리 뒤 수동 재시도를 안내한다.
  기술 보호값은 DEV 소유이고 상품 quota가 아니다. 기존 활성 자격은 임의 철회하지 않으며 철회 목록의 사라짐은 권한 복구나 영구 백업 삭제가 아니다.
  추가 제품 답과 독립 관찰 조건: [POLICY 실행 기록](../../docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md#기록-보호의-추가-제품-답과-관찰-조건).

- 2026-10-06: SAR-PUBLIC-AGENTS-001-UI-FIX는 fixed 458798c의 F-UI-01–04와 POLICY 안내를 actual desktop1280/mobile390에서 직접 확인해 좁은 UI PASS로 기록했다.
  원본 d1/d165 FAIL·QA/OPS 실행 65개 파일은 byte 동일이다. 제어 25h fixture와 실제 sweep을 실운영 24h로 확대하지 않는다.
  직접 근거·한계·인계: [UI-FIX 보고](../../docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md).

- 2026-10-06: SAR-PUBLIC-MESSAGES-001-UI는 fixed09c의 실제 일반회원 receipt·gate를 자기SQL/SMTP·Chrome에서 직접 검수했다. UX06 로컬PASS, UX07 Deny 뒤405로 FAIL/보류다.
  F-UI-MSG-01 medium은 거절 저장 뒤 결과·홈 복귀 부재다. 제품 코드·규칙은 보존하며 원본OPS high/medium과 fixture 초기실패도 유지한다.
  직접 증거49PNG·manifest·회수·DEV후속: [MESSAGES 보고](../../docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md).

- 2026-10-06: SAR-PUBLIC-MESSAGES-001-UI-FIX는 fixed dfc70ca의 정상 Deny 클릭/Enter desktop1280/mobile390와 owner desktop을 직접 확인해 좁은 UI PASS로 기록했다.
  신규 rate 포화 뒤 denied 저장/결과GET429는 즉시결과PASS에서 제외한다. 원09c FAIL·viewport980·OPS H1/M1·216개원천을 보존했다.
  직접42PNG·실패/회수·재사용 경계·인계: [UI-FIX 보고](../../docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md).

- 2026-10-06: UI-FIX archive 뒤 coor msg_23b8de8f85ac가 dfc의 OPS H2 high와 유효정리18/18의429를 인계했다. 정상Deny 시각PASS는 전체제품수락이 아니다.
  같은DEV후속수정·새fixed독립수락 전 main을 차단한다. 원dfc캡처/완료전문·빈인박스를 보존하고 다음후보는UI의존성/영향만 확인한다.
  원메시지·담당·재개조건: [UI-FIX phase](../../docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-UI-FIX.md).

- 2026-10-07: UI-FIX-2는 실제 별도 세션/managed d089 snapshot과 자기 빌드/migration 뒤 사용자 추가 검증 중단으로 마감했다. helper 초기화 exit1이며 새 브라우저/PNG/UI PASS는 없다.
  원09c FAIL·원dfc 좁은 PASS/GET429 제한·pending/운영/벤더 보류를 유지한다. 제품은 읽기 전용이며 자기 fixture만 회수했다.
  현재 결과·재개 담당: [UI-FIX-2 보고](../../docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX-2.md).

- 2026-10-11: SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL은 eac603b6의 tester 합성 캡처 5장을 직접 열어 반환·재시작 안내 화면을 좁은 시각 PASS로 기록했다.
  F2 링크24px는 실제 클릭 성공과 24px 최소 기준으로 수락 가능하다. F1 medium(원인별 문구 없음)은 비차단 잔여로 유지하고 CSRF·계정 binding 캡처가 같은 파일임을 기록했다.
  증거·한계·인계: [시각 수락 보고](../../docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL/report.md).

- 2026-10-11: SAR-SERVICE-UI-001-DESIGN은 로그인·내 에이전트·연결 관리의 상용 UX·비주얼 명세와 DEV 인계(R-01–R-11)를 작성했다. 홈 문서는 6394/7986px로 길고 한국어 anywhere 줄바꿈 결함이 있다.
  제품 코드·서버 규칙·CSP·재인증은 보존했고 JS·이미지 없는 CSS 기준만 정했다. 구현 후 직접 시각 검수와 UI 수락은 후속이다.
  명세·관찰·인계: [회원 웹 명세](../../docs/design-docs/mockups/service-member-ui.md), [실행 기록](../../docs/exec-plans/phases/SAR-SERVICE-UI-001-DESIGN.md).

- 2026-10-11: SAR-SERVICE-UI-001-VISUAL은 고정 c7828b2의 합성 화면 29개를 Orca 내장 브라우저 1440·390·320에서 직접 검수해 시각 PASS로 기록했다. critical·high 0건이다.
  D04 9절 44px 기준을 글자 링크 21px가 채우지 못한 medium 1건(V-01)과 low 11건이 남는다. 실제 Google 왕복·실계정·200%·스크린 리더는 미검증이다.
  도구 교훈: `goto`가 뷰포트를 초기화하고 `full-screenshot`이 에뮬레이션에서 타일을 반복하므로 스크롤 조각으로 찍는다. 증거·한계: [시각 수락 보고](../../docs/evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md).

- 2026-10-11: SAR-SERVICE-UI-001-FIX-VISUAL은 고정 295417a의 글자 링크 11종을 Orca 내장 브라우저 320·390·1440에서 직접 재측정했다. 모두 44px 이상이며 V-01은 시각 기준 resolved다.
  원 보고 V-01–V-12는 보존했고 low N-01만 추가했다. 키보드 포커스 재입력·실제 Google·스크린 리더·200%는 미검증이다.
  증거·한계: [재측정 보고](../../docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL/report.md).

- 2026-10-12: SAR-FRONTEND-POLISH-001-VISUAL은 후보 9597bad의 제품 동일성과 manifest를 확인하고 after PNG 26장을 직접 검수했다. 제공된 변경 화면은 시각 PASS다.
  V-02·V-03·V-08·V-12와 mobile 생성 버튼 폭·명령 줄바꿈은 비차단 잔여다. 포커스·320px는 DEV 근거 재사용이며 운영·실계정·200%·스크린 리더는 미검증이다.
  근거·후속: [직접 시각 검수 보고](../../docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-VISUAL/report.md).

- 2026-10-12: SAR-HUMAN-PEER-001-VISUAL은 최종 후보 254ff6a의 memberStyle/대화 template 동일성을 확인한 뒤 원래 4eaefa6의 실제 합성 PNG 21장을 모두 직접 검수했다. 시각 수락은 FAIL이다.
  V-01 작성란 기술 설명·반복 대상 문장, V-02 긴 목록 이름/desktop 작성 위치는 medium 수정 대상이다. V-03 반복 제목은 low다. 실제 surface·200%·focus의 완료 범위와 not_run을 이미지 판정과 구분한다.
  근거·수정 수락 조건·후속: [사람 친구·대화 시각 보고](../../docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/report.md). 제품 코드는 수정하지 않았다.

- 2026-10-12: 같은 VISUAL 새 attempt161c5363은 후보12dc0c7의 원본14PNG와 별도 실제15상태를 직접 확인했다. V01–V03은 제공 범위 resolved/시각 PASS이며 원254ff6a FAIL은 보존한다.
  실 keyboard/200%·미제공 상태는 not_run이다. 별도page의 요청/실제 CSS 치수를 대조하고 보정했으며 viewport를 실확대로 주장하지 않는다. Orca 세션 토큰 도구로그 노출은 값 없이 coor에 인계했다.
  새 근거·한계·후속: [수정 후보 재검수 보고](../../docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/final-report.md). 제품·다른역할·PLANS/board는 수정하지 않았다.

- 2026-10-14: SAR-LANDING-GUIDE-001은 최신 요청에 따라 공개 첫페이지를 소개·시작 순서·Grok/Dots 연결·가상 예시·로그인으로 확장한다. 유효 세션도 안내를 읽고 대화로 이동하며 회원 앱 구조와 권한은 유지한다.
  D02/D04/D11과 Astra high DEV/tester 인계를 작성했다. 원 partial/unknown·최소 로그인 지시·운영 미검증을 보존했으며 새 제품/브라우저/시각 수락은 아니다.
  결정·검증·후속: [실행 기록](../../docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md).
