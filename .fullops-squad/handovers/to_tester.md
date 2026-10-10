---
title: SAR-AUTO-RECEIVE-001-TESTER — 자동수신 고정 후보의 독립 QA와 snapshot 코드 리뷰
status: draft
updated: 2026-10-10
owner: tester
tasks: [SAR-AUTO-RECEIVE-001-TESTER]
summary: 자동수신 고정 후보의 독립 QA와 snapshot 코드 리뷰
attempt: 1c5542ce04914b9d93cfad64624462a7
base: 51eebe5c2b588147627afb8f1ab9ec2d166dd432
subagent_level: off
test_level: lite
---

# SAR-AUTO-RECEIVE-001-TESTER — 자동 수신 독립 검증

- Task key: SAR-AUTO-RECEIVE-001-TESTER; Purpose: review and independent QA
- 상태: 대기. DEV 완료 고정 SHA와 실제 별도 reviewer session/snapshot이 준비되기 전 착수하지 않는다.
- From / To: coor / tester. 테스트 lite, 선택 하위 위임 off.
- 소유권: 본 인박스, QA/시나리오/리뷰 결과. 제품 코드 수정 금지. coor PLANS/board 수정 금지.
- 승인: 사용자 자동수신 보완 요청. 추가 비용·로컬 Docker·그록봇 UI 조작·vendor core 수정·자동 업무 실행 금지.
- 복귀: 실제 worker-start preamble의 현재 run/task/dispatch를 그대로 사용한다.

## 적용 기준과 먼저 읽을 문서

fullops-common-0.3.3과 FULLOPS.md, project.md, rules/common/README.md 및 coding-style/testing/security, docs/agents/document-writing.md를 읽는다. fullops-review/fullops-test 및 open-code-review-delegate를 적용한다. 구현자와 다른 실제 세션/clean detached snapshot을 읽기 전용으로 검토한다. 결과는 snapshot 밖 기록 워크트리에 쓴다.
제품 정본: docs/deliverables/README.md에서 시작해 현재 자동수신 요구·기술 구조·인터페이스·운영/복구·다음 일로 이어지는 경로를 확인한다. adapters/README.md, adapters/skills/knowslink/SKILL.md, DEV 실행 기록 SAR-AUTO-RECEIVE-001-DEV.md와 같은 키의 로그/packet-outcomes를 확인한다. 수동 receive를 UI로 유도한 기존 왕복은 자동수신 증거가 아니다.

## 해야 할 일과 완료 기준

- [ ] 고정 SHA diff와 설치된 SDK 계약을 검토한다. 실제 호스트 자동 전달 지원 근거를 확인한다. MCP 알림 전송을 노우 턴 시작으로 표시하면 결함이다.
- [ ] 수동 receive 호출 없이 자동 수신/서명·관계 검증/지속 보존/ACK/호스트 알림을 관찰한다. 이후 수동 조회에서 메시지 원문이 유실되지 않는지 확인한다.
- [ ] 중복·재시작·경합·권한 철회·만료·네트워크 실패/백오프 핵심 경계를 짧게 검증한다. 원문과 비밀값은 로그에 남기지 않는다. 실제 운영 증거는 coor와 조율하며 새 회원/키/관계를 임의 변경하지 않는다.
- [ ] DEV의 동일 SHA 필수 lint/빌드/검사 증거를 확인한다. 영향 없는 전체 회귀는 반복하지 않는다. Node22 명시 경로를 사용한다.
- [ ] OCR prepare/check와 독립 리뷰 result/report/lint를 기록한다. 기존 동일 HEAD lint 재사용은 출처와 hash를 명시한다. 미해결 critical/high 또는 필수 실패는 수락 차단이다.
- [ ] QA 결과와 리뷰의 파일별 reviewed/skipped 사유, 한계와 재현 절차, 실제 host 미검증 범위를 남긴다. 자동 저장, 호스트 알림, 노우 확인/답장을 구분한다.

## 기대 산출물과 후속

D10 QA 보고서, 본 과제 실행 기록, 독립 snapshot 리뷰. D12는 실제 운영 검증 영향이 없으면 변경하지 않고 이유를 적는다. 구현자 self-review로 대체하지 않는다. 고정 SHA/구현자 세션/reviewer 세션/snapshot 경로는 coor가 후보 완료 후 이 인박스에 기록한다.
완료 전문을 채우고 work.py finish로 로그/빈 인박스를 확인한 뒤 실제 worker_done을 전송한다. 코드 결함은 수정하지 말고 재현 근거로 보고한다. 최종 운영 설치·실제 Bot 수신 수락은 coor가 이어서 진행한다.

## 완료 보고

고정 SHA / 독립 세션과 snapshot / QA 결과 / 리뷰 결론 / lint 출처 / 미검증·후속.
