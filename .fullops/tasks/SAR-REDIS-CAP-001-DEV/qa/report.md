---
title: SAR-REDIS-CAP-001-DEV QA 보고
status: draft
updated: 2026-11-28
owner: dev
summary: "Redis 메시지 사본의 경로별 최악 계산(현재·2026-12-28 뒤), 바이트 상한 32 MiB·owner 6 MiB와 근거, verify.py redis 관측, 테스트 결과"
---

# SAR-REDIS-CAP-001-DEV QA 보고

- base: dee77d7df450cc3251dd8bc4177afb087c67cc56
- 근거: 주인 승인 순서 ②(쿠어 msg_2e4feb01e049). 이전 어림은 [SAR-MSG-TTL-1H-001-DEV R3](../../SAR-MSG-TTL-1H-001-DEV/qa/report.md), 운영 값은 [SAR-MCP-DEPLOY-001-OPS](../../SAR-MCP-DEPLOY-001-OPS/qa/report.md)(`maxmemory` 128 MiB, `noeviction`, 사용 1.5 MiB)다.
- 테스트 레벨: lite. 운영 배포·Redis 설정 변경은 하지 않았다.

## 1. 저장 구조와 입력값(코드)

- Redis에는 `ephemeral{Messages, Idempotency, Inboxes}` JSON 한 개만 있다. commit은 새 키에 SET한 뒤 이전 키를 DEL하므로 commit 중 사본이 2개다(`ephemeral.go`). 그래서 Redis 정점 ≈ 2 × 사본 + 기본 1.5 MiB다.
- 속도: 새 작업 member당 40/분, 전체 200/분 → 시간당 member 2,400건, 전체 12,000건. 기록 수 < 20,000. 미배달(queued·leased) 전체 < 100, owner별 < 20. 원격 설치 ≤ 500, 설치당 Inbox ≤ 32.
- 크기: text 본문 ≤ 4,096 B, `/v1/send` raw envelope ≤ 32 KiB. 저장 JSON은 Go `encoding/json`이 `<`·`>`·`&`를 `<`로 쓴다. text envelope는 이미 이 형태라 ≤ 약 25 KiB다. `/v1/send` envelope는 JCS 정규화가 `<`를 그대로 두므로 Redis JSON에서 최대 6배까지 커진다(아래 표는 raw 기준이고 괄호가 escape 최악이다).
- 기록 수명: 모든 기록(메타·멱등·웹 대화 사본)은 수락 1시간 뒤 sweep. 원문(envelope·`Inbox` 사본·원격 Inbox)은 전달 기한까지다. 지금(D076)은 자기 agent 사이 1시간, 그 밖 text 180초·envelope 300초. 2026-12-28 뒤는 모두 1시간이다.

## 2. 최악 계산 — 상한 없을 때(base)

경로별 사본 1개 기준. "전체"는 시간당 12,000건, "member 1명"은 2,400건이다.

| 경로 | 기록당 | 지금 수명 | 지금 최악(전체) | 12-28 뒤 수명 | 12-28 뒤 최악(전체) | member 1명(둘 다) |
|---|---|---|---|---|---|---|
| 메타·멱등 | ≤ 1.5 KiB | 1시간 | 18 MiB | 1시간 | 18 MiB | 3.5 MiB |
| 웹 대화 사본(text·이름) | ≤ 4 KiB + 이름(escape 24 KiB) | 1시간 | 47 MiB(281 MiB) | 1시간 | 47 MiB(281 MiB) | 9.4 MiB(56 MiB) |
| 원문 미배달·미ACK(queued·leased, envelope+Inbox) | ≤ 2 × 32 KiB | 전달 기한 | 큐 100건 = 6.3 MiB | 같음 | 6.3 MiB | 20건 = 1.3 MiB |
| 원문 ACK 뒤 업무 envelope(결과 전, envelope+Inbox) | ≤ 2 × 32 KiB(escape 384 KiB) | 자기 1시간 / 그 밖 300초 | 자기 사이로 몰면 750 MiB, 그 밖만이면 1,000건 63 MiB | 1시간(모든 pair) | 750 MiB(4.4 GiB) | 150 MiB |
| 원격 Inbox 사본(받고 `knowslink_ack` 안 함) | ≤ 25 KiB | 자기 1시간 / 그 밖 180초 | 자기 293 MiB, 그 밖만 600건 15 MiB | 1시간 | 293 MiB | 59 MiB |

- 한 기록은 위 경로 중 하나의 원문만 갖는다(text는 ACK에서 원문을 지우고, 업무 envelope는 대화 사본이 없다). 그래서 합은 기록당 최대 약 1.5 + 64 = 65.5 KiB × 12,000 ≈ **768 MiB**(지금·12-28 뒤 같음)다. commit 중 2배라 Redis 정점은 약 1.5 GiB다. 128 MiB를 약 12배 넘는다.
- 지금과 12-28 뒤의 차이는 "누가 채울 수 있나"다. 지금은 750 MiB·293 MiB 경로를 쓰려면 같은 owner의 자기 agent 사이여야 한다. 12-28 뒤에는 친구 사이 메시지도 1시간 머문다. 수신 client가 ACK하지 않으면 발신자가 아닌 쪽 때문에도 쌓인다.
- 참고: 위 값 근처에서는 relay(256 MiB)의 JSON 직렬화가 먼저 메모리를 다 쓸 수 있다. 측정하지 않았다.

## 3. 바이트 상한 — 값과 근거

| 상수(`internal/relay/capacity.go`) | 값 | 근거 |
|---|---|---|
| `messageBytesLimit` | 32 MiB(사본 1개) | commit 중 2개 → 64 MiB = `maxmemory` 128 MiB의 50%. 기본 1.5 MiB를 더해도 51%다 |
| `ownerMessageBytesLimit` | 6 MiB(보내는 owner 1명) | 32/6 → owner 5명이 가득 채워도 6번째가 2 MiB를 쓴다. `ownerQueueLimit`(20/100, 최소 5 owner)와 같은 비율 |
| `messageOverhead` | 1,536 B | 고정 메타·멱등 항목·Inbox 키. 가장 긴 키·ID로 잰 실제 최대 약 1.4 KiB. 20,000건 × 1.5 KiB = 30 MiB라 기존 건수 상한이 바이트 상한보다 먼저 걸린다 |

- 셈(`retainedBytes`): 기록마다 `messageBytes` = overhead + envelope + `Inbox` 사본 + 웹 대화 사본(본문·이름 4개)의 JSON 크기 상한(`storedLen`). 아직 persist하지 않은 envelope는 이후 생길 `Inbox`·원격 Inbox 사본 1부를 미리 센다. 원격 설치 Inbox 사본은 원 메시지의 발신 owner에 센다. 발신 agent가 없어진 기록은 전체에만 센다.
- 판정(`messageCapacity`): 새 기록 크기를 더해 owner 몫 > 6 MiB면 그 owner만, 전체 > 32 MiB면 모두 `capacity`(429, 기존 응답)다. text(`ingestPublicText`)·웹 owner text(`sendOwnerText`)·`/v1/send`(`ingest`) 세 입구가 같은 함수를 쓴다. 수신·ACK·정리는 막지 않으므로 상한에 걸려도 빠져나갈 수 있다. 오류와 로그에 본문·발신자를 남기지 않는다(기존 `capacity` 경로 그대로).
- 상한 뒤 최악: 사본 ≤ 32 MiB → Redis 정점 ≤ 64 + 1.5 = **65.5 MiB(51%)**, commit 사이 ≤ 33.5 MiB(26%). 전달 기한과 무관하므로 **2026-12-28 뒤에도 같다**. 한 owner가 남에게 주는 영향은 6 MiB까지다.
- 정상 사용 여유: 본문 약 200 B text는 기록당 약 1.9 KiB다. 전체 시간당 12,000건 = 22 MiB, member 시간당 2,400건 = 4.5 MiB로 둘 다 상한 안이다. 4 KiB 본문(기록당 약 5.9 KiB)을 분당 40건 계속 보내는 member는 약 25분 뒤 owner 몫에 걸리고, 1시간 sweep이 오래된 기록을 지우는 만큼 다시 보낸다.
- 비용: 새 기록마다 보관 중 바이트를 한 번 훑는다(최대 32 MiB). commit마다 이미 사본 전체를 직렬화하므로 그보다 작다.
- 운영 Redis 설정 추천: 바꾸지 않는다(현행 128 MiB·noeviction 유지). 상한이 정점을 51%로 묶는다.

## 4. 관측

- `deploy/knowslink/verify.py redis`(새, 읽기 전용): `docker exec knowslink-redis-1 redis-cli INFO memory`에서 `used_memory`·`maxmemory`·`maxmemory_policy`만 읽어 비율을 출력한다. `verify.py local`도 끝에 실행한다(`beta.sh deploy` 경로).
- 경고 기준: ≥ 40%면 stderr `WARNING`(통과). commit 사이 정상 최대 26%를 넘었다는 뜻이다: `retire` 실패로 남은 이전 사본이나 예산이 가득 찬 순간이다. ≥ 75%이거나 `maxmemory` 0, 정책이 `noeviction`이 아니면 실패(종료 1).
- 문서: D12 맨 앞 절(관측·경고·대응), D12 파일 표의 `verify.py` 행, D07 SAR-MSG-EPHEMERAL-001 절의 "용량" 항목.

## 5. 테스트 발췌

```
go vet ./... && go test -count=1 ./...
ok  	github.com/tourbut/KnowsLink/cmd/migrate
ok  	github.com/tourbut/KnowsLink/cmd/relay
ok  	github.com/tourbut/KnowsLink/internal/config
ok  	github.com/tourbut/KnowsLink/internal/relay	2.990s
go vet -tags=integration ./internal/relay/ ./cmd/...   # 통과(통합 테스트 컴파일)

go test -run 'TestOwnerByteLimit|TestMessageBytesBoundStoredJSON|TestOwnerQueueLimit|TestSharedCapacitiesAndCleanupClassification|TestPublicText|TestDeliveryCeiling' -v ./internal/relay/
--- PASS: TestDeliveryCeilingBySelfAgents
--- PASS: TestPublicTextWire
--- PASS: TestPublicTextRoundtripAndCurrentAuth
--- PASS: TestPublicTextTTLAndLeaseFailures
--- PASS: TestOwnerQueueLimit
--- PASS: TestOwnerByteLimit
--- PASS: TestMessageBytesBoundStoredJSON
--- PASS: TestSharedCapacitiesAndCleanupClassification
```

표적 단언:
- `TestOwnerByteLimit`: owner a의 보관 바이트를 6 MiB − 16 KiB까지 채우면 a의 4,096 B(`<` escape) text는 `capacity`. owner b는 같은 크기를 보낸다. a의 기록 하나를 지우면 a가 다시 보낸다. 발신 agent가 없는 기록으로 전체 32 MiB를 넘기면 c와 무소유 발신 모두 거부된다.
- `TestMessageBytesBoundStoredJSON`: escape 최악 본문(`<`×4096, U+2028×1365, `"`×4096, 줄바꿈×4095)과 최장 메타·U+2028 이름 40자 기록에서, 실제 `json.Marshal(ephemeralPart)` 크기 ≤ `retainedBytes` ≤ 2 × 실제. queued(예약 포함)와 persist 뒤 값이 같다.
- `TestSharedCapacitiesAndCleanupClassification`: 빈 기록 20,000건 건수 상한이 바이트 상한보다 먼저 걸린다(경계 19,999건 통과).
- `verify.py` 자체 확인(가짜 `INFO memory`): 1.5 MiB → `(1%)` 통과, 57.2 MiB → `(45%)` `WARNING` 통과, 104.9 MiB → `verify failed: Redis used_memory 82% >= 75% of maxmemory`.

## 6. 남은 것

- 운영 미적용. 배포 뒤 `verify.py redis` 값 확인은 ops 몫이다.
- 통합 테스트(실 Redis·Postgres)는 lite라 돌리지 않았다. 바이트 상한은 메모리 상태 판정이라 단위 테스트로 덮었다.
- `retire` 실패로 남는 이전 사본은 상한 밖이다(TTL 1시간이 경계). 관측의 40% 경고로 잡는다.
