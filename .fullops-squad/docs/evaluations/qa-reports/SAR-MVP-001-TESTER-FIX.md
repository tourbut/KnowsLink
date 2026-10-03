---
title: SAR-MVP-001-TESTER-FIX — 수정 후보의 독립 QA 보고서
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-MVP-001-TESTER-FIX]
summary: 제품 4262d02의 human 전달 차단 통과와 legacy claim high 실패를 기록한다
---

# SAR-MVP-001-TESTER-FIX — 수정 후보의 독립 QA 보고서

## 판정

제품 후보는 `4262d02fdd7b0b57a804d9e550597852950ffeae`다. 검사 시작 트리의 제품 경로는 이 커밋 및 `d7e2149`와 같다.
H가 아닌 `deliver:human`의 신규 send와, 저장된 human 행의 pull·persist·ACK·claim 차단은 이번 HTTP probe에서 통과했다.
같은 실행에서 owner gate, agent delivery, approval/result, current-auth, 철회, TTL, lease, CSRF 회귀도 통과했다.
기존 claimed 부모의 authorize와 `relay.result` 수락은 실패다. 이 high는 제품 수락을 차단한다.
신규 transport 통과는 이 high의 해소가 아니다. 수정 뒤 새 제품 후보의 재검증이 필요하다.
이 판정은 병합 승인이나 공개 배포 완료가 아니다. 제품 코드는 수정하지 않았다.
시나리오: [SAR-MVP-001-TESTER-FIX.md](../scenarios/SAR-MVP-001-TESTER-FIX.md). 로그: [SAR-MVP-001-TESTER-FIX-test/](SAR-MVP-001-TESTER-FIX-test/).

## 환경

실행 날짜는 2026-10-03이다. 브랜치는 `fullops/tester`다. 워크트리는 `/home/shin/orca/workspaces/KnowsLink/fullops-tester`다.
Go 1.27.1, Node v22.22.2, Python 3.12.3, GNU Make 4.3, Docker 29.4.3, Compose v5.1.3이다.
적용 규칙은 `fullops-common-0.3.2`다. 프로젝트 정본은 [project.md](../../../project.md)다. 프로토콜 기준은 frozen C1–C5다.
명령은 기존 `run.py`가 셸 없이 실행하고 `[exit N]`을 남긴다.
합성 credential, CSRF, claim token, DB 비밀번호는 로그에 남기지 않았다. Tunnel profile은 켜지 않았다. Compose project는 종료 시 그 project만 `down --volumes` 했다.
버전 종료코드는 [versions.log](SAR-MVP-001-TESTER-FIX-test/versions.log)에 있다.

## 이번 실행 결과

판정 실행은 probe.log의 세 번째 `isolated.py`다. 마지막 줄은 `PROBE fail 1 pass 17`이다. 그 실행의 `[exit 1]`이 제품 판정이다.
첫 실행의 `[exit 1]`은 tester SQL이 `relay_state` 없이 `data`를 읽어 중단된 것이다. 제품 결과로 세지 않는다. 기록은 [probe-results-abort.json](SAR-MVP-001-TESTER-FIX-test/probe-results-abort.json)이다.
두 번째 실행은 동작 확인이다. `UPDATE 1` 태그가 shape 비교에 섞여 새 claim 항목이 실패로 보였다. 같은 줄의 새 claim 상태는 403이다.
세 번째 실행이 [probe-results.json](SAR-MVP-001-TESTER-FIX-test/probe-results.json)이다.

| ID | 결과 | 관찰 |
|---|---|---|
| FIX-direct-human-query | 통과 | `schedule.query` `deliver:human` 403 `sender_not_allowed`. 행 없음. pull 본문 `null`. claim 403. receipt 403. gate 0 |
| FIX-direct-human-commit | 통과 | `schedule.commit` `deliver:human` 403 `sender_not_allowed`. 행 없음 |
| FIX-result-human-stays-schema | 통과 | `relay.result` `deliver:human` 422 `invalid_schema` |
| FIX-stored-human-transport | 통과 | 저장 `Deliver=human`은 pull되지 않음. 위조 lease의 persist·ACK 409 `invalid_lease`. delivered claim 403. Claimed false. gate 0 |
| FIX-legacy-new-claim-denied | 통과 | 봉투 `deliver=human`이고 `Deliver` 필드가 없는 claimed 행의 새 claim 403 |
| FIX-legacy-claim-authorize-result | 실패 | authorize 200. `relay.result` 200. completion `denied`. gate 0 |
| REG-current-auth-agent-credential | 통과 | agent credential의 owner revoke 401 |
| REG-pair | 통과 | invite 200. accept `active` |
| REG-agent-delivery | 통과 | lease 간격 30.0초. attempts 1. ACK-before-persist 409. persist 200. ACK 200. transport `delivered`. claim 200. 재claim 409 |
| REG-agent-restart | 통과 | 재시작 후 health 200. 재claim 409. 새 `deliver:agent` 임대 |
| REG-ttl | 통과 | 301초 422 `ttl_too_long`. 같은 key의 유효 TTL 200 |
| REG-lease-max-attempts | 통과 | transport `failed:max_attempts` |
| REG-revoke | 통과 | revoke 후 token 없음. 교체 key도 기존 대기를 임대하지 않음. 새 전송은 임대 |
| REG-unpair-generation | 통과 | unpair 중 ACK 409. 재수락 세대 2. 이전 replay 403. 새 전송 200 |
| REG-clock-unknown-auth | 통과 | 시계 +1시간 503 `abnormal_clock`. 복구 후 health 200. 없는 credential 401 |
| REG-owner-gate-open | 통과 | H는 claim과 함께 200이고 저장 `Deliver`는 `human`. pull은 H를 임대하지 않음. claim 없는 H는 403. ACK 전 transport는 `leased` |
| REG-gate-csrf | 통과 | typed body와 정책 표시. hint는 화면에 없음. GET은 pending. agent 401. Basic 200. 잘못된 CSRF 403. approve 303. 재결정 409. consume의 executable false, disclosure false |
| REG-approval-result | 통과 | commit stub 비실행. `done` 403. 잘못된 방향 403. optional result 422. denied 200. 두 번째 result 403 `sender_not_allowed`. completion `denied`. 결과 ACK 200. 다음 pull에 token 없음 |

`make verify-mvp` 종료코드는 0이다. 고유 project `knowslink-mvp-3c62d6cca4`에서 migration, Postgres integration, TypeScript stub, owner UI seed가 실행됐다.
`TestPostgresSafety/agent_cannot_process_human_delivery`는 통과했다. gate CSRF, unpair·만료·clock, approval binding, 동시 accept, gate 실패 상태, frozen parse도 같은 종료코드 0에 포함된다.
Go integration은 Compose Postgres에 연결된다. `TEST_DATABASE_URL`이 없으면 테스트가 실패한다. TypeScript stub은 `http://127.0.0.1:54587`에 대해 종료코드 0이다.
`make test`는 integration tag를 포함하지 않는다. `internal/relay`는 1.060초에 실행됐다.

## 명령과 종료코드

| 명령 | 종료코드 | 로그 |
|---|---|---|
| 제품 diff `4262d02`와 HEAD, `d7e2149`와 `4262d02` | 0 | [versions.log](SAR-MVP-001-TESTER-FIX-test/versions.log). 제품 경로 출력은 비어 있다 |
| `a6a10c7`과 `4262d02`의 설정 경로 diff | 0 | 같은 로그. Compose, Makefile, adapter, cmd, db는 같다 |
| PNG blob 8개 | 0 | [png-identity.log](SAR-MVP-001-TESTER-FIX-test/png-identity.log) |
| `make lint` | 0 | [lint.log](SAR-MVP-001-TESTER-FIX-test/lint.log) |
| FullOps `lint.py --from 4262d02` | 0 | [fullops-lint.json](SAR-MVP-001-TESTER-FIX-test/fullops-lint.json). HEAD `bcc8d34009f00e749975546f66e162b72e8b2ba2`, ERROR 0, WARNING 0, 실행 불가 0. 같은 기준의 `7c3e5ecd4f12f081c3b093bd116353d05737b5bc` 재실행도 종료코드 0, ERROR 0, WARNING 0, 실행 불가 0 |
| `make test` | 0 | [unit.log](SAR-MVP-001-TESTER-FIX-test/unit.log) |
| `make verify-mvp` | 0 | [verify-mvp.log](SAR-MVP-001-TESTER-FIX-test/verify-mvp.log). 마지막 판정 줄은 PASS |
| `isolated.py` 첫 실행 | 1 | [probe.log](SAR-MVP-001-TESTER-FIX-test/probe.log). tester SQL 중단. 제품 판정 아님 |
| `isolated.py` 두 번째 실행 | 1 | 같은 로그. shape 태그 혼선. 새 claim 상태는 403 |
| `isolated.py` 세 번째 실행 | 1 | 같은 로그의 마지막 `[exit 1]`. `down --volumes` 종료코드 0. 제품 실패 1건 |

`make verify`와 `make verify-runtime`은 실행하지 않았다. 이번 실제 DB 명령은 `make verify-mvp`와 고유 Compose의 HTTP probe다.

## 결함

미해결 high가 하나 있다. 제품 수락은 차단된다.

reviewer 증거는 dev 체크아웃의 `SAR-MVP-001-REVIEW-FIX-review/legacy-claim.log`와 `legacy-claim-probe.txt`다. 이 파일은 읽기만 했다.
`a6a10c7`에서 test-only overlay가 `deliver:human`을 pull·persist·ACK·claim해 상태를 저장했다. 그 로그의 종료코드 0은 재현 테스트의 종료다.
`4262d02`에서 그 상태를 읽으면 새 claim은 거부된다. authorize와 `relay.result`는 수락된다. completion은 `denied`다. owner gate는 0개다.

이번 probe는 같은 서버 형태를 실제 relay HTTP와 Postgres에서 확인했다.
정상 `deliver:agent`를 claim한 뒤, SQL로 `Deliver` 필드를 지우고 봉투 `deliver`를 `human`으로 바꿨다. 이 형태는 옛 직렬화에 `Deliver` 필드가 없을 때의 적재 형태다.
새 claim은 403이다. authorize는 200이다. `relay.result`는 200이다. completion은 `denied`다. gate는 0이다.
이 probe의 claim token은 형태를 바꾸기 전에 `4262d02`가 발급했다. `a6a10c7` 발급 token의 증거는 reviewer 로그다.
두 증거는 같은 서버 동작을 가리킨다. delivered이고 claimed인 human 부모가 authorize와 result를 통과한다.

C1은 agent credential의 `deliver:human` 처리와 owner 승인 대체를 금지한다. 신규 전송 차단은 그 문장의 새 행에 해당한다. 이미 발급된 claim의 authorize와 result 수락은 남아 있다.
`parentRouting`과 authorize는 부모의 deliver 경로를 거부하지 않는다. 제품 코드는 수정하지 않았다.
수정 후보가 나오면 그 SHA에서 이 항목을 다시 검증한다. 이번 통과 항목을 그 후보의 통과로 옮기지 않는다.

## 재사용

`a6a10c7`에서 `4262d02`까지의 제품 차이는 `internal/relay/http.go`, `store.go`, `integration_test.go`, `protocol_test.go`다.
`http.go`의 변경은 claim 조건 한 줄이다. 승인 화면 템플릿 줄은 변경 diff에 없다.
설정 파일과 adapter, cmd, db의 diff 종료코드는 0이다.
아래 PNG blob은 원래 커밋과 HEAD가 같다. 이번 실행은 새 캡처를 만들지 않았다. 시각 판정을 새로 통과로 기록하지 않는다.

| PNG | 원래 커밋 | blob |
|---|---|---|
| [v01-pending.png](SAR-MVP-001-TESTER-test/ui/v01-pending.png) | `c59537b` | `adfaec781bed4cd1c39a17eaf650d270de3678e1` |
| [v02-approved.png](SAR-MVP-001-TESTER-test/ui/v02-approved.png) | `c59537b` | `b514d2c43c4a7233498421531e120cd8228155de` |
| [v02-denied.png](SAR-MVP-001-TESTER-test/ui/v02-denied.png) | `c59537b` | `e939ca5f9c53474e15d8851278d5b07b53aedbda` |
| [v03-expired.png](SAR-MVP-001-TESTER-test/ui/v03-expired.png) | `c59537b` | `15c84cccda23b2033226d9a0a347652a717940b9` |
| [v03-revoked.png](SAR-MVP-001-TESTER-test/ui/v03-revoked.png) | `c59537b` | `9f5b1549d77373d65692c358a04ea0f997ea419e` |
| [v04-unavailable.png](SAR-MVP-001-TESTER-test/ui/v04-unavailable.png) | `c59537b` | `fdb291ed60627a3766d7a74f4cc7e9d9a4eda220` |
| [v04-unauthorized.png](SAR-MVP-001-TESTER-test/ui/v04-unauthorized.png) | `c59537b` | `87538f2208f0b2db69044fe45183abf95adc854e` |
| [v01-pending-full.png](../../design-docs/mockups/SAR-MVP-001-UI-v01-pending-full.png) | `e238777` | `97efa9185256e419ea1e1e3e9ee950d919914085` |

원래 시각 조건은 [SAR-MVP-001-TESTER.md](SAR-MVP-001-TESTER.md)와 [SAR-MVP-001-UI.md](SAR-MVP-001-UI.md)에 있다. designer의 V-01–04 판정은 `e238777`의 기록이다.
`protocol.go`는 제품 diff에 없다. QA-03의 HTTP 변형 7종은 원래 probe `c59537b`의 조건으로 남긴다. 이번 `make verify-mvp`의 frozen parse 통과는 그 재사용과 별도인 이번 실행이다.
QA-11의 24시간 원문 정리와 WAL 항목은 이번 SHA에서 다시 실행하지 않았다. `store.go`가 바뀌었으므로 원래 통과를 이번 통과로 옮기지 않는다. 원래 결과는 `c59537b`에 남는다.
QA-10의 실연결과 webhook HTTP는 이번 probe에서 반복하지 않았다. adapter 트리는 `a6a10c7`과 같다. 미설정 adapter의 원래 관찰은 `c59537b`에 남는다.
전체 QA-01–11 probe와 화면 캡처는 반복하지 않았다.

## held

원래 held를 유지한다. 이번 통과로 바꾸지 않는다.

- QA-02-free-n. Free N, 가격, slot 단위는 미확정이다.
- QA-06-epoch-cas. 트랜잭션 밖의 고의 stale epoch 503은 이번에도 만들지 않았다.
- QA-10-real-adapter와 QA-10-a2a. 벤더 연결과 A2A 현행 개정 검토는 이번 범위가 아니다.
- QA-11-wal. 행 삭제를 실행하지 않았고, 실행해도 WAL 또는 backup 삭제가 아니다.
- DEC-02. 실데이터 positive silent done은 held다.
- DEC-03. 공개 한도, 실제 신원 증명, singleton 처리량은 held다.
- 다른 플랫폼과 `make verify-runtime`은 실행하지 않았다.

## 후속

dev가 legacy claim의 authorize와 result 수락을 수정한다. 수정 제품 SHA에서 tester가 그 항목과 영향 회귀를 다시 실행한다.
coor는 이 high가 열린 동안 제품 수락과 병합을 차단한다. 원래 `a6a10c7` QA의 held와 실패 기록은 유지한다.
독립 QA 완료는 제품 최종 수락이 아니다.
