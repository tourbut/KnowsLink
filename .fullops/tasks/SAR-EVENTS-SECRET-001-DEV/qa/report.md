---
id: qa-SAR-EVENTS-SECRET-001-DEV
title: SAR-EVENTS-SECRET-001-DEV QA 보고
status: approved
updated: 2026-11-28
owner: dev
summary: Events secret 봉인·Beta dual 광고·poll의 표적 검증 결과
---

# SAR-EVENTS-SECRET-001-DEV QA 보고

테스트 레벨 lite. 변경 동작의 표적 단위 테스트와 격리 Postgres·Redis integration(`TestRemote*`)으로 확인했다. 운영 배포·운영 env 변경은 하지 않았다.

## 결정(추천안으로 진행, 되돌릴 수 있음)

| 결정 | 선택 | 근거 | 되돌리는 법 |
|---|---|---|---|
| 키 이름·형식 | `KNOWSLINK_SEAL_KEY`, base64 32바이트, AES-256-GCM, AAD=구독 ID | 이전 `KNOWSLINK_CALENDAR_KEY`와 같은 운영 형식. 다음 과제 A가 agent ID AAD로 재사용 | env 이름만 `cmd/relay/main.go`에서 바꾼다 |
| 키가 없을 때 | relay는 기동하고 `events/subscribe`만 `events_webhook_unavailable`로 거절. Beta `events/list`는 `["poll"]`만 광고 | 평문 저장 금지. poll은 secret이 없어 계속 쓸 수 있다. 개발·테스트 환경이 키 없이 뜬다 | — |
| 잘못된 키 | 기동 실패 | 조용히 기능이 꺼지는 오설정 방지 | — |
| 평문 이행 | 읽기는 `whsec_` 평문 fallback, 성공한 subscribe/refresh는 항상 봉인 값으로 다시 쓴다 | DB migration 불필요. TTL 최대 1시간이라 갱신되지 않는 평문도 1시간 안에 정리된다 | fallback 한 줄(`remoteStoredSecret`) |
| 열리지 않는 봉인 값(키 분실·교체) | 배달 때 구독을 410처럼 삭제. 다음 refresh가 새로 검증·봉인 | 서명할 수 없는 구독을 남기지 않는다 | — |
| dual | 지금 형식(OpenAI·Dots)은 응답 바이트까지 그대로. Beta 동작은 요청 `_meta` clientCapabilities에 확장을 선언한 2026-07-28 요청에만 적용 | 두 문서가 `webhook-id`(eventId와 같음 vs 별개)와 검증 실패 코드(-32015 vs -32027)에서 충돌 | `remoteDeclaresEvents` |
| `events/poll` | 구현(상태 없음, cursor = 대기 집합 digest) | Beta 스펙의 가장 싼 배달 | — |
| `events/stream` | 하지 않음 | 구독자마다 긴 SSE 요청·heartbeat를 붙잡아 relay 동시 요청 예산을 쓴다. 실사용자 Dots는 webhook만 쓴다 | 별도 과제 |
| 알 수 없는 event 이름 | 기존과 같은 -32602 `invalid_schema` | Beta의 NotFound(-32023) 세분화는 이번 범위 밖 | — |

## 완료 기준 대응

| 기준 | 테스트 |
|---|---|
| 암호화 왕복 | `TestSealerRoundTripAndRefusals`: 왕복, nonce 재사용 없음, 평문 미포함, 다른 AAD·다른 키·변조·키 없음·평문 입력은 열리지 않음, 잘못된 키 값은 오류에 되풀이되지 않음 |
| 키가 없을 때 | `TestSealerRoundTripAndRefusals`(빈 키 = nil, 오류 없음), `TestRemoteSubscribeRefusedWithoutSealKey`(트랜잭션 전 거절, 응답에 secret 없음), `TestRemoteEventsDualShapeAndMetadataOnlyBody`(Beta 목록 `["poll"]`) |
| 평문 이행 | `TestRemoteSecretMigrationReadsPlainTextUntilResealed`(fallback 읽기·비교, 키 분실 시 봉인 값 거부) + integration `remoteLegacySecretReseal`(평문 구독 저장 → 실제 `/mcp` OpenAI 형식 refresh → 저장 값이 봉인되고 같은 secret으로 열림, 응답에 secret 없음 → Beta url-only unsubscribe로 삭제) |
| dual 광고 | integration `TestRemoteOAuthHTTP`: `server/discover`에 `"events":{}`와 `"extensions":{"io.modelcontextprotocol/events":{}}` 둘 다. `remoteEraChecks`: 선언한 Beta 요청의 `events/list`는 `["webhook","poll"]`, 선언하지 않은 요청은 `["webhook"]` |
| poll 결과 | `TestRemotePollIsStatelessAndMetadataOnly`: null·생략 cursor는 0건+cursor, 도착 뒤 1건, 같은 집합 재poll 0건, 잘못된 입력 거절. integration: modern `events/poll` 결과 필드·`Mcp-Name` 불일치 400 -32020 |
| Dots 회귀 | 단위: 미선언 목록 `["webhook"]`, `webhook-id`=`eventId`. integration: 기존 OpenAI 형식 subscribe의 127.0.0.1 callback은 여전히 -32015, 같은 형식 refresh 성공. 기존 `TestRemote*` 전부 통과 |
| 본문·발신자 없음 | `TestRemoteEventsDualShapeAndMetadataOnlyBody`(webhook 본문 키가 정확히 `cursor,data,eventId,name,timestamp`, `data`=`{"arrival":true}`), poll 결과에 본문·발신자·메시지 ID 없음 |

## 실행 기록

```
go build ./... && go vet ./internal/relay/                       # 통과
go vet -tags=integration ./internal/relay/ ./cmd/...             # 통과
go test -count=1 ./internal/... ./cmd/...                        # ok (relay 3.0s)
# 격리 컨테이너 postgres:17-alpine·redis:8.10.2-alpine(고유 이름, 끝난 뒤 삭제), go run ./cmd/migrate 뒤
TEST_SYNTHETIC_DATABASE=1 TEST_DATABASE_URL=<격리> TEST_REDIS_URL=<격리> \
  go test -count=1 -tags=integration -run 'TestPostgresSafety|TestRemote' ./internal/relay/   # ok 5.3s
```

같은 DB에 두 번째로 돌리면 `TestPostgresSafety`가 실패한다. 이 테스트는 빈 DB를 전제한다(첫 실행은 통과, 변경과 무관). 컨테이너는 실행 뒤 삭제했다.

## 하지 않은 것

- 실제 Dots 계정의 구독 갱신·배달(운영 배포 뒤 tester/ops).
- `events/stream`, `notifications/events/list_changed`, Beta `deliveryStatus`, `maxAgeMs` 적용(받기만 하고 무시).
- 운영 `.env` 키 생성·주입과 이행 확인(ops, D12 절차).
