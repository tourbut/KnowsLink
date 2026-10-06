---
title: designer 컨텍스트
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-SETUP-001, SAR-PREP-002, SAR-MVP-PUBLIC-POLICY-001, SAR-MVP-001-UI, SAR-PUBLIC-SERVICE-001, SAR-PUBLIC-AGENTS-001-UI, SAR-PUBLIC-AGENTS-001-POLICY]
summary: 제품 판단과 일반 서비스 및 직접 UI 검수의 수락과 보류 경계를 보존한다
---

# designer 컨텍스트

결정·교훈을 항목당 3줄 이내로 기록한다.

- 2026-10-03: SAR-SETUP-001은 초기 구성·lint 과제다. 첫 기능 MVP의 human-gate 포함 결정은 유지하며 이번 범위로 당기지 않는다.
  원천의 확정 기술 선택과 project.md의 이전 미정 상태를 구분한다. 기술 세부 결정은 dev에게 맡긴다.
  상세 근거: [SAR-SETUP-001 실행 기록](../docs/exec-plans/phases/SAR-SETUP-001.md).

- 2026-10-03: SAR-PREP-002는 최신 7bc9ea1의 개발 준비다. D01·새 MVP D02·기능 백로그와 queued 인계를 작성하고 초기 setup·원천·제품 코드는 보존한다.
  첫 기능은 로컬 등록·수락·안전 전달·human-gate다. 무정책 일정 공개는 deny이며 실데이터 silent 성공·가격/quota·운영 공개는 담당과 재개 조건을 가진 보류다.
  상세 결정·검증·인계: [SAR-PREP-002 실행 기록](../docs/exec-plans/phases/SAR-PREP-002.md).

- 2026-10-03: SAR-MVP-PUBLIC-POLICY-001은 link.knowslog.com의 누구나 가입 가능한 인증 합성 파일럿 기준이다. 가입 초대 제한을 추가하지 않으며 pairing 수락은 유지한다.
  DEC-03은 수용량/rate/size/concurrency 권장안 단계다. 신규 수락과 안전 정리 budget을 분리하고 DEV/OPS 근거 및 사용자 수치 승인 전 공개 held를 유지한다.
  결정·검증·재개 조건: [공개 정책 기록](../docs/exec-plans/phases/SAR-MVP-PUBLIC-POLICY-001.md).

- 2026-10-03: SAR-MVP-001-UI의 기존 PNG 7개와 보완 PNG 1개를 직접 확인했다. 합성 V-01–04 시각 판정은 PASS다.
  `invalid_auth`는 인증 실패 차단만 증명한다. `deliver:human` 인증 경계 high는 DEV 수정·독립 재검증 전 유지한다.
  상세 관찰·재개·인계: [UI 실행 기록](../docs/exec-plans/phases/SAR-MVP-001-UI.md).

- 2026-10-05: SAR-PUBLIC-SERVICE-001은 일반 서비스 기반 수락 뒤 동일 이메일의 Grok Bot “노우”↔OpenAI dot “다닷” 운영 시험이다. agentID·키·credential과 관계 수락은 각각 유지한다.
  새 운영 기본값·복구 목표는 이번 위임 결정이다. 과거 DEC-03 미승인과 실일정·상품 held는 소급하지 않는다. 인증/API/DB/MCP의 기술 선택은 DEV/OPS 책임이다.
  상세 판단·OPS 근거·기능 인계·검증: [일반 서비스 실행 기록](../docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md).

- 2026-10-06: SAR-PUBLIC-AGENTS-001-UI의 d1eef9b 연결·키·관계를 격리 fixture와 실제 브라우저에서 직접 확인했다. 모바일 지문 overflow와 오류 뒤 복귀 부재로 UX 시각 수락은 FAIL/보류다.
  제품 코드·제품 규칙과 사용자 자격·운영 자원은 보존했다. Orca blank PNG는 실패 근거이며 실제 Playwright PNG와 구분한다.
  상세 관측·마스킹·미검증·DEV 후속: [직접 검수 보고](../docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md).

- 2026-10-06: SAR-PUBLIC-AGENTS-001-POLICY는 pending/active 반복의 새 초대를 금지하고 거절·만료·양측 철회 뒤 수동 재초대를 허용한다. 새 세대와 수신 owner의 새 수락 전 메시지 거부를 유지한다.
  차단·쿨다운·새 수치는 추가하지 않았다. 수신측 반복 노출 위험과 원본 UI FAIL/보류·실제 운영 후속은 유지한다.
  제품 답·DEV/QA 관찰 조건·검증: [POLICY 실행 기록](../docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md).

- 2026-10-06: 같은 POLICY 추가 지시에서 키 기록 포화는 새 agent 연결·새 관계 수락으로 처리한다. owner 기록 포화는 필요한 철회·최소 24h 보존·실제 정리 뒤 수동 재시도를 안내한다.
  기술 보호값은 DEV 소유이고 상품 quota가 아니다. 기존 활성 자격은 임의 철회하지 않으며 철회 목록의 사라짐은 권한 복구나 영구 백업 삭제가 아니다.
  추가 제품 답과 독립 관찰 조건: [POLICY 실행 기록](../docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md#기록-보호의-추가-제품-답과-관찰-조건).
