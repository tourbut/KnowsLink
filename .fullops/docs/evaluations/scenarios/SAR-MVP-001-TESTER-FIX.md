---
title: SAR-MVP-001-TESTER-FIX — 수정 후보의 독립 QA 시나리오
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-MVP-001-TESTER-FIX]
summary: "고정 제품 4262d02의 human 전달 차단과 관련 회귀 절차, 재사용 조건, held 경계를 정의한다"
---

# SAR-MVP-001-TESTER-FIX — 수정 후보의 독립 QA 시나리오

제품 대상은 `4262d02fdd7b0b57a804d9e550597852950ffeae`다. 이 트리의 제품 경로는 그 커밋과 같아야 한다.
실행 위치는 tester 워크트리 `/home/shin/orca/workspaces/KnowsLink/fullops-tester`다.
명령은 기존 `run.py`로 실행하고 종료코드를 그대로 남긴다.
합성 credential, claim token, CSRF 값은 로그에 남기지 않는다. Tunnel과 운영 DB는 사용하지 않는다.
이 시나리오는 원래 [SAR-MVP-001-TESTER 시나리오](SAR-MVP-001-TESTER.md)를 바꾸지 않는다.

## 절차

1. `git diff --exit-code 4262d02 HEAD --` 제품 경로가 비어 있는지 확인한다. `d7e2149`와 `4262d02`의 제품 경로도 같아야 한다.
2. `a6a10c7` 대비 제품 차이는 `internal/relay/http.go`, `store.go`, `integration_test.go`, `protocol_test.go`만인지 확인한다. Compose, Makefile, adapter, cmd, db는 같아야 한다.
3. 기존 PNG 8개의 Git blob이 원래 커밋과 같은지 확인한다. 화면 템플릿 줄이 바뀌지 않았을 때만 그 PNG를 원래 SHA와 조건으로 재사용한다. 새 캡처는 만들지 않는다.
4. `make lint`와 `make test`를 실행한다.
5. `make verify-mvp`로 고유 Compose project의 실제 Postgres integration과 relay HTTP adapter를 실행한다.
6. `isolated.py`가 다른 고유 project를 띄우고 `probe.mjs`를 실행한다. 종료 시 그 project만 `down --volumes` 한다.
7. 환경 오류는 제품 결함이나 통과로 바꾸지 않는다. 제품 결함은 재현 상태와 종료코드로 남기고 제품 코드는 고치지 않는다.

## 이번 실행의 판정

| ID | 통과 조건 | 실행 |
|---|---|---|
| FIX-direct-human-query | H가 아닌 `schedule.query`의 `deliver:human`은 403 `sender_not_allowed`다. 행이 없다. pull 본문은 `null`이다. claim과 receipt는 403이다. gate는 0이다. | probe |
| FIX-direct-human-commit | H가 아닌 `schedule.commit`의 `deliver:human`도 같은 403이며 행이 없다. | probe |
| FIX-result-human-stays-schema | `relay.result`의 `deliver:human`은 기존 422 `invalid_schema`다. 새 수락이 아니다. | probe |
| FIX-stored-human-transport | 저장된 `Deliver=human`은 pull되지 않는다. 위조 lease의 persist와 ACK는 409이다. delivered 상태의 claim은 403이다. claim과 gate는 생기지 않는다. | probe |
| FIX-legacy-new-claim-denied | 봉투 `deliver=human`이고 `Deliver` 필드가 없는 기존 claimed 행의 새 claim은 403이다. | probe |
| FIX-legacy-claim-authorize-result | 그 행의 기존 ClaimToken으로 authorize와 `relay.result`가 거부된다. 통과는 403이다. 수락은 결함이다. | probe |
| REG-agent-delivery | `deliver:agent`의 lease 30초, ACK 전 persist 실패, persist, ACK, transport `delivered`, claim 1회와 재claim 거부가 유지된다. | probe와 `make verify-mvp` |
| REG-agent-restart | relay 재시작 뒤 기존 claim은 거부되고 새 `deliver:agent`는 임대된다. | probe |
| REG-ttl | 301초는 422 `ttl_too_long`이고 같은 key의 유효 TTL은 200이다. | probe |
| REG-lease-max-attempts | 만료 lease 3회 뒤 transport는 `failed:max_attempts`다. | probe |
| REG-revoke | key revoke 확정 후 기존 대기는 임대되지 않는다. 교체 key도 그 대기를 복구하지 않는다. 새 전송은 임대된다. | probe |
| REG-unpair-generation | unpair 중 ACK는 409이다. 재수락은 새 세대다. 이전 replay는 403이고 새 전송은 200이다. | probe |
| REG-clock-unknown-auth | 시계 +1시간은 503 `abnormal_clock`이다. 복구 후 health는 200이다. 없는 credential은 401이다. agent credential의 owner revoke는 401이다. | probe |
| REG-owner-gate-open | H는 claim이 있을 때 200이고 저장 `Deliver`는 `human`이다. agent pull은 H를 임대하지 않는다. claim 없는 H는 403이다. ACK 전 transport는 `delivered`가 아니다. | probe |
| REG-gate-csrf | 검증 body와 정책이 보인다. hint는 판단 근거가 아니다. GET은 pending을 유지한다. agent는 401, Basic은 200, 잘못된 CSRF는 403, approve는 303, 재결정은 409이다. consume의 executable과 disclosure는 false다. | probe |
| REG-approval-result | commit stub은 비실행이다. `done`은 403이다. 잘못된 방향은 403, optional result는 422, denied 1회는 200, 두 번째 result는 403 `sender_not_allowed`다. completion은 `denied`다. 결과 pull은 추가 명령을 남기지 않는다. | probe |

## 재사용과 held

변화 없는 동작과 PNG는 위 동일성 확인 뒤에만 원래 SHA와 조건으로 연결한다. 이 실행의 통과로 바꾸지 않는다.
원래 held는 DEC-02, DEC-03, Free N, 실adapter, A2A 현행 검토, WAL 또는 backup 삭제, 고의 stale epoch다.
원래 tester 보고서의 시각 held는 designer 검수 기록과 별도다. 이번 실행은 시각을 다시 판정하지 않는다.
고의 stale epoch는 이번에도 외부에서 강제하지 않는다. 그 held를 유지한다.
