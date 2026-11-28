---
id: handover-ops
title: SAR-EVENTS-DEPLOY-001-OPS — KNOWSLINK_SEAL_KEY를 운영 .env에 넣고 65a752fe를 배포해 Redis 상한·Events secret 봉인을 확인한다
status: approved
updated: 2026-11-28
owner: ops
summary: KNOWSLINK_SEAL_KEY를 운영 .env에 넣고 65a752fe를 배포해 Redis 상한·Events secret 봉인을 확인한다
key: SAR-EVENTS-DEPLOY-001-OPS
role: ops
goal: KNOWSLINK_SEAL_KEY를 운영 .env에 넣고 65a752fe를 배포해 Redis 상한·Events secret 봉인을 확인한다
base: 65a752fe2e88e0da1cbf1a844cb0000c51595295
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-EVENTS-DEPLOY-001-OPS — KNOWSLINK_SEAL_KEY를 운영 .env에 넣고 65a752fe를 배포해 Redis 상한·Events secret 봉인을 확인한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-28
- From / To: coor(lead) / ops
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: 사용자가 이 터미널 질문 카드로 승인했다("키 넣고 배포", 2026-11-28). 주인도 쿠어 msg_45d31927c16c로 같은 승인을 줬다. 순서는 다음과 같다:
  1. 백업과 rollback 태그
  2. 서버에서 `KNOWSLINK_SEAL_KEY` 생성과 운영 `.env` 주입
  3. `65a752fe` 배포
  4. verify(Redis 점검 포함)
  5. 기동 로그 확인, 1시간 뒤 평문 0건 확인
  실패하면 직전 배포 `d3670f2f`로 되돌린다. 이번 승인은 `.env`에 이 키 한 줄을 더하는 것까지 포함한다. ops 브랜치 push(`fullops done`)도 승인됐다.
- 추가 승인이 필요한 행위: 그 밖의 `.env`·Tunnel·DNS·방화벽 변경, Redis 설정 변경, DB 마이그레이션, 데이터 정리. 필요하면 멈추고 `ask`로 묻는다.
- 상위 의도: 없음

## 현재 상황과 확인 근거

- 운영 서버는 Vultr `knowslink-relay-01`이고 `d3670f2f`가 돌고 있다(SAR-MCP-DEPLOY-001-OPS).
- 배포 대상 `65a752fe`(origin/main)에 들어간 과제:
  - SAR-REDIS-CAP-001-DEV: 메시지 사본 바이트 상한(전체 32 MiB, owner 6 MiB), `verify.py redis` 점검(40% 경고, 75% 실패)
  - SAR-EVENTS-SECRET-001-DEV: Events `whsec` secret을 `KNOWSLINK_SEAL_KEY`(AES-256-GCM)로 봉인, Beta 광고와 `events/poll`. 리뷰 accept
- lead가 `d3670f2f..65a752fe`에서 확인한 배포 관련 변경:
  - `compose.yaml`: relay env에 `KNOWSLINK_SEAL_KEY: ${KNOWSLINK_SEAL_KEY:-}` 한 줄
  - `deploy/knowslink/verify.py`: redis 점검
  - `db/`·`Dockerfile`·`tunnel/`: 변경 없음
  배포 전에 직접 다시 확인한다.
- 키가 없으면 Dots의 다음 구독 갱신이 거절된다. 그래서 키를 먼저 넣는다.
- 기존 평문 secret은 읽을 때는 fallback으로 받고, 다음 갱신 때 다시 봉인된다. 1시간 뒤 평문 0건 확인 절차는 D12에 있다(SAR-EVENTS-SECRET-001-DEV가 적음).

## 적용 기준과 예외

- 키 값은 서버에서 `openssl rand -base64 32`로 만들어 `.env`에 바로 쓴다. 다음에는 키 값을 남기지 않는다: argv, 터미널 출력, 로그, 보고서, 메일, 이 PC.
  - 예: 서버 스크립트 안에서 변수로 받아 파일에 append한다.
  - `.env` 권한(0600)을 유지한다.
- 키 백업: 서버 밖 백업 위치와 방법은 D12에 '위치와 절차'만 적는다. 값은 적지 않는다.
  - 기존 `.env` 백업 관례가 있으면 그대로 따른다.
  - 없으면 추천안을 보고에 적는다. 새 외부 저장소를 만들지 않는다.
- SSH 스크립트는 파일로 올려 실행한다. 서버 명령은 Bash 도구로 실행한다.

## 먼저 읽을 문서

- `.fullops/guide/roles/ops.md`
- `.fullops/tasks/SAR-MCP-DEPLOY-001-OPS/qa/report.md`와 배포 스크립트. 같은 절차를 대상 SHA만 바꿔 쓴다.
- D12의 SAR-EVENTS-SECRET-001-DEV 절(새 env 키, 이행 확인), SAR-REDIS-CAP-001-DEV 절(verify redis)

## 해야 할 일과 파일 소유권

- [ ] 배포 전 `beta.sh backup`, rollback 태그, `.env` 백업(서버 안, 0600)을 하고 해시를 기록한다.
- [ ] 서버에서 키를 만들고 `.env`에 주입한다. 키가 이미 있으면 덮어쓰지 않고 멈춘다.
- [ ] `65a752fe`를 배포한다. `verify.py local`·`public`(redis 포함) 통과와 relay error 0줄을 확인한다.
- [ ] 기동 로그에서 `events webhook configured=true`(또는 문서에 적힌 같은 뜻의 줄)를 확인한다.
- [ ] 1시간 뒤 평문 0건을 확인한다(D12 절차, 읽기 전용). 1시간을 기다리기 어려우면 그 시점에 다시 확인하는 절차를 남기고, 보고에 '미확인'으로 적는다.
- [ ] 운영 읽기 전용으로 쿠어·Grok Bot agent와 grant 3개가 그대로인지 확인한다.
- [ ] D12에 이번 배포 절과 키 백업 위치·절차를 짧게 더한다.

## 완료 기준과 검증

- `tasks/SAR-EVENTS-DEPLOY-001-OPS/qa/report.md`에 다음을 적는다(키 값은 넣지 않는다):
  - 배포 SHA
  - 백업 해시
  - 키 주입 여부(길이만)
  - verify 발췌
  - Redis 값
  - 기동 로그 줄
  - 평문 0건 확인 결과
  - agent·grant 상태
- 실패하면 `d3670f2f`로 되돌린다. 주입한 키는 남겨도 무해하다. 남겼는지 보고에 적는다.

### UI 작업의 디자인 기준과 검증

해당 없음.

## 갱신할 산출물

- D12

## 기대 산출물

- qa 보고, D12

## 제약·협업·후속

- 범위 밖: 코드 수정, 실제 Dots 갱신과 배달 확인(배포 뒤 주인 단계), Orca 설치본 갱신

## 완료 보고

`fullops done`을 실행하고, 그 출력 한 줄을 worker_done 본문으로 쓴다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
