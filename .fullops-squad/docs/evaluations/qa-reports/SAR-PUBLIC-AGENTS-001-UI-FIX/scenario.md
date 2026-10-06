---
title: SAR-PUBLIC-AGENTS-001-UI-FIX — 직접 시각 검수 시나리오
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-AGENTS-001-UI-FIX]
summary: F-UI-01–04와 POLICY 안내의 실제 직접 검수 재현 조건을 기록한다
---

# SAR-PUBLIC-AGENTS-001-UI-FIX — 직접 시각 검수 시나리오

이 문서는 실행한 좁은 수동 시각 검수의 재현 조건이다. Jev 자동 조작 테스트를 실행했다고 주장하지 않는다. 지시서가 직접 화면 검수를 지정했다. API·race 독립 QA는 TESTER 담당이다. 모든 단계의 제품 fixed SHA는 `458798c2ee15c179edacfd6f94ebb9896d26f411`이다.

1. 고정 detached clone에 새 Postgres·SMTP sink·loopback relay를 준비한다. synthetic signup/allowlist·운영 환경을 사용하지 않는다.
2. 실제 이메일 코드 경로로 합성 회원 A/B를 각각 새 비영속 context에서 로그인한다. 각 context는 별도 cookie를 쓴다.
3. A의 첫 agent를 실제 Node prepare·지문 대조·브라우저 승인·complete로 연결한다. desktop1280/mobile390에서 prepared·home 지문과 select를 확인한다.
4. waiting을 실제 취소한다. 별도 waiting의 Exp를 relay 정지 중 past로 바꾼다. cancelled/expired/consumed에 취소 버튼이 없는지 확인한다.
5. 자기 form에 unsupported option과 없는 agent를 각각 넣는다. 거부 뒤 홈 링크를 실제 누른다. Verified만 6m past로 바꾸고 재확인 버튼으로 실제 이메일 확인을 완료한다.
6. A의 두 agent와 B의 agent로 pending 같은/반대 방향 반복, 받은/보낸 결정 위치, active 반복 notice를 확인한다.
7. 실제 수신 거절·송신 철회·수신 철회와 controlled Exp 만료를 각각 확인한다. 종료 카드의 새 초대 버튼 뒤 새 pending과 새 수락 대기를 확인한다.
8. 같은 owner 두 agent도 첫 초대는 pending이고 실제 수락 뒤만 active인지 확인한다.
9. own relay 정지 중 live agent 키 기록을 20개로 만든다. 새 기록은 revoked public-key-only fixture이며 25h 이상이다. 포화 안내와 실제 connect 거부의 새 agent/홈 버튼을 확인한다.
10. 실제 선택 키를 철회한다. 키 기록은 20개로 남고 새 공간을 보장하지 않는지 확인한다. 오류의 새 agent 버튼으로 별도 ID를 만든다. 실제 정상 연결과 B에 대한 별도 새 pending/새 수락을 확인한다.
11. 실제 활성 5개 생성 경계와 owner 기록 10개 fixture 경계를 각각 관측한다. 최소 24h·실제 정리 뒤 수동 재시도와 기존 관리 복귀 안내를 확인한다.
12. 포화 중 실제 agent 철회를 수행한다. pre-revoke 경고와 retained 철회 상태·복구 불가/영구 백업 삭제 아님 안내를 확인한다. 기록 slot이 바로 해제되지 않는지 확인한다.
13. own relay 정지 중 revoked Changed만 25h past로 바꾼다. 제품의 실제 sweep 뒤 옛 ID/관련 pair가 목록에서 사라지고 수동 새 생성이 가능한지 확인한다.
14. 세션 없는 새 context의 거부에서 Tab으로 로그인 링크에 focus한다. 링크로 실제 시작 화면에 복귀한다.
15. 인증 입력을 raster capture 전에 mask한다. 공개 fingerprint는 남긴다. PNG를 직접 열어 판정하고 manifest에 실제 viewport·fixture·fixedSHA를 연결한다.
16. 자기 browser/page/process/Compose/private files만 회수한다. 기존 서비스 ID·원본 UI/QA/OPS의 byte 동일을 확인한다.

각 단계의 결과와 PNG는 [직접 보고](../../../design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md)와 [manifest](evidence.json)에 있다. 실제 24h 경과·실메일·실플랫폼·운영 공개는 실행하지 않았다.
