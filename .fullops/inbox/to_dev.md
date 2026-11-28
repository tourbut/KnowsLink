---
id: handover-dev
title: SAR-REDIS-CAP-001-DEV — 보관 1시간의 Redis 이론 최악 사용량을 128MiB noeviction 한도 안에 묶고 used_memory 관측을 더한다
status: approved
updated: 2026-11-28
owner: dev
summary: 보관 1시간의 Redis 이론 최악 사용량을 128MiB noeviction 한도 안에 묶고 used_memory 관측을 더한다
key: SAR-REDIS-CAP-001-DEV
role: dev
goal: 보관 1시간의 Redis 이론 최악 사용량을 128MiB noeviction 한도 안에 묶고 used_memory 관측을 더한다
base: dee77d7df450cc3251dd8bc4177afb087c67cc56
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-REDIS-CAP-001-DEV — 보관 1시간의 Redis 이론 최악 사용량을 128MiB noeviction 한도 안에 묶고 used_memory 관측을 더한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-28
- From / To: coor(lead) / dev
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: 주인 승인 순서 ②(쿠어 msg_2e4feb01e049). 되돌릴 수 있는 결정은 추천안으로 정한다. 검사를 통과한 병합은 lead가 바로 push한다(2026-11-28 사용자 결정). 그래서 끝나면 `fullops done`을 실행해도 된다.
- 추가 승인이 필요한 행위: 운영 배포, 운영 Redis 설정 변경(maxmemory, 정책), DB 마이그레이션, 보관 기한 규칙 변경.
- 상위 의도: 없음

## 현재 상황과 확인 근거

- 운영 `d3670f2f`: 같은 owner 자기 agent 사이 메시지만 보관 1시간이다. 미배달은 owner별 20, 전역 100이 상한이다.
- 운영 Redis는 `maxmemory 128mb`, `noeviction`이고, 지금 사용량은 1.5MiB다(SAR-MCP-DEPLOY-001-OPS 보고).
- 이론상 최악 사용량은 member 1명 기준 약 150MiB로, 한도를 넘는다. 근거는 `.fullops/tasks/SAR-MSG-TTL-1H-001-DEV/qa/report.md`의 용량 어림(R3)과 `.fullops/tasks/SAR-MCP-DEPLOY-001-OPS/qa/report.md`다. 주된 원인은 persist 뒤 ACK하지 않은 원문과 Inbox 사본이 오래 머무는 것이다.
- noeviction이라 한도를 넘으면 쓰기가 `unavailable`이 된다. 다른 사용자의 발송도 함께 막힌다.
- 2026-12-28에 모든 메시지를 1시간 보관으로 전환한다(SAR-MSG-TTL-1H-002-DEV). 이 과제는 그 전에 끝나야 한다.

## 적용 기준과 예외

- 공통 보안 규칙을 따른다. 본문은 untrusted다. 상한을 넘을 때의 오류는 기존 `capacity` 계열을 쓴다. 본문과 발신자를 로그에 남기지 않는다.
- 표적 테스트만 둔다.

## 먼저 읽을 문서

- 위 두 보고서
- `internal/relay`의 store(Redis 키, TTL), 미배달 큐 상한(`ownerQueueLimit`), 원격 Inbox 사본, persist와 ACK 흐름

## 해야 할 일과 파일 소유권

- [ ] 최악의 경우를 다시 계산한다. 경로별(원문, Inbox 사본, 기록)로 바이트를 나눠 보고서에 표로 적는다. 2026-12-28 전체 1시간 전환 뒤의 최악값도 같이 계산한다.
- [ ] 바이트 기준 상한을 둔다. 예를 들어 owner별 보관 원문 바이트 상한과 전역 상한을 두고, 전체가 한도의 일정 비율(예: 50%) 안에 들게 한다. 값과 근거를 적는다. 넘으면 그 owner만 `capacity`로 거부하고, 다른 owner는 계속 쓸 수 있게 한다.
- [ ] 관측을 추가한다. 운영에서 Redis `used_memory`와 한도 비율을 볼 수 있게 한다. 기존 health나 운영 점검 스크립트(`verify.py` 등)에 읽기 전용 항목을 더하는 정도로 한다. 경고 기준도 정한다.
- [ ] 문서: D12 운영자설명서의 관측 항목, 그리고 해당 설계 문서의 용량 절을 고친다.

## 완료 기준과 검증

- 표적 테스트:
  - 바이트 상한을 넘으면 그 owner만 거부되고, 다른 owner는 통과한다.
  - 상한 계산이 맞는지 확인한다.
- `fullops lint --from <base>`를 통과한다. `tasks/SAR-REDIS-CAP-001-DEV/qa/report.md`에 다음을 적는다:
  - 최악 계산 표(지금과 12-28 뒤)
  - 상한 값과 근거
  - 테스트 발췌

### UI 작업의 디자인 기준과 검증

해당 없음.

## 갱신할 산출물

- D12, 시스템 설계 문서의 용량 절

## 기대 산출물

- 코드, 표적 테스트, 문서, qa 보고

## 제약·협업·후속

- 범위 밖:
  - 운영 Redis 설정 변경. 필요하다고 판단되면 추천만 적는다.
  - Events와 routine 작업
- 하위 에이전트는 쓰지 않는다.

## 완료 보고

`fullops done`을 실행하고 그 출력 한 줄을 worker_done 본문으로 쓴다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
