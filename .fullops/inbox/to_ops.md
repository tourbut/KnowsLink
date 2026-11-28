---
id: handover-ops
title: "SAR-MCP-DEPLOY-001-OPS — d3670f2f를 운영에 배포하고 verify, Redis maxmemory, policyEffective, 기존 연결을 확인한다"
status: approved
updated: 2026-11-28
owner: ops
summary: "d3670f2f를 운영에 배포하고 verify, Redis maxmemory, policyEffective, 기존 연결을 확인한다"
key: SAR-MCP-DEPLOY-001-OPS
role: ops
goal: "d3670f2f를 운영에 배포하고 verify, Redis maxmemory, policyEffective, 기존 연결을 확인한다"
base: d3670f2fa3c85a8d320542055188296948aa783b
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-MCP-DEPLOY-001-OPS — d3670f2f를 운영에 배포하고 verify, Redis maxmemory, policyEffective, 기존 연결을 확인한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-28
- From / To: coor(lead) / ops
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: 사용자가 이 터미널 질문 카드로 승인했다("포함해 push + 배포", 2026-11-28). 순서는 다음과 같다:
  1. 백업
  2. `d3670f2f` 운영 배포
  3. verify
  4. 실패하면 직전 배포 `8272fb91`로 되돌린다
  ops 브랜치 push(`fullops done`)도 승인됐다.
- 추가 승인이 필요한 행위: DB 마이그레이션, 데이터 정리, 부하 테스트, `.env`·Tunnel·DNS·방화벽 변경, Redis 설정 변경. Redis maxmemory가 부족해 보여도 바꾸지 않고 보고한다. 이런 일이 생기면 멈추고 `ask`로 묻는다.
- 상위 의도: 없음

## 현재 상황과 확인 근거

- 운영 서버는 Vultr `knowslink-relay-01`이고, 지금 `8272fb91`이 돈다(SAR-ORCA-HOST-001-OPS).
- 배포 대상은 `d3670f2f`(origin/main)이다. 들어간 과제는 다음과 같다:
  - SAR-MCP-UPGRADE-001: 원격 /mcp 2026-07-28 dual-era. 리뷰 accept
  - SAR-MSG-TTL-1H-001-DEV: 같은 owner 자기 agent 사이 메시지만 보관 1시간, owner별 미배달 상한 20, 개인정보 방침 ko·en 개정과 2026-12-28 고지. 리뷰 accept
  - SAR-MCP-SPEC-CHECK-001-DEV: 필수 누락 3건 수정. modern 헤더 필수, clientCapabilities 필수, DCR application_type 허용
- lead가 확인한 `8272fb91..d3670f2f`의 변경: `db/`·`deploy/`·`compose.yaml`·`Dockerfile`·`tunnel/` 변경 없음. 배포 전에 직접 다시 확인한다.
- 배포 순서: relay를 먼저 배포하고 adapter는 나중에 한다(README). 이 PC의 Orca 설치본 갱신은 범위 밖이다.

## 적용 기준과 예외

- D12 운영 절차와 공통 보안 규칙을 따른다. SSH 스크립트는 파일로 올려 실행한다. 서버 명령은 Bash 도구로 실행한다.

## 먼저 읽을 문서

- `.fullops/guide/roles/ops.md`
- `.fullops/tasks/SAR-ORCA-HOST-001-OPS/qa/report.md`와 `deploy-8272fb91.sh`. 같은 절차를 대상 SHA만 바꿔 쓴다.
- `.fullops/tasks/SAR-MSG-TTL-1H-001-DEV/qa/report.md`의 용량 어림(R3)

## 해야 할 일과 파일 소유권

- [ ] 배포 전 `beta.sh backup`과 rollback 태그를 만들고, 해시를 기록한다.
- [ ] `d3670f2f`을 배포한다. `verify.py local`·`public` 통과와 relay error 0줄을 확인한다.
- [ ] 운영 Redis의 `maxmemory`, `maxmemory-policy`, 현재 `used_memory`를 읽기 전용으로 기록한다. R3의 이론 최악값과 견준다.
- [ ] `policyEffective`를 확인한다. 개인정보 방침 시행일 상수는 2026-11-28이다. 배포일이 다르면 멈추고 `ask`로 묻는다.
- [ ] 공개 응답에서 `/.well-known/oauth-authorization-server`와 `/mcp` legacy initialize가 정상인지 확인한다. 예전 방식 요청이 2025-06-18이면 2025-06-18로 답해야 한다.
- [ ] 운영 읽기 전용으로 쿠어·Grok Bot agent와 grant 3개가 그대로인지 확인한다.
- [ ] D12에 이번 배포 절을 짧게 더한다.

## 완료 기준과 검증

- `tasks/SAR-MCP-DEPLOY-001-OPS/qa/report.md`에 다음을 적는다:
  - 배포 SHA
  - 백업 해시
  - verify 발췌
  - Redis 값
  - policyEffective
  - /mcp 확인 발췌
  - agent·grant 상태
- 실패하면 `8272fb91`로 되돌리고 원인을 보고한다.

### UI 작업의 디자인 기준과 검증

해당 없음.

## 갱신할 산출물

- D12 운영자설명서

## 기대 산출물

- qa 보고, D12

## 제약·협업·후속

- 범위 밖:
  - 코드 수정
  - 이 PC Orca 설치본 갱신
  - 쿠어의 Grok decision 재시험. 배포 뒤 lead가 쿠어에게 요청한다.

## 완료 보고

`fullops done`을 실행하고 그 출력 한 줄을 worker_done 본문으로 쓴다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
