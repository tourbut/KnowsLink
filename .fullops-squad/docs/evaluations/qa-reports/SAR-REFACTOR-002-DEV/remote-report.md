---
title: Remote OAuth·MCP·Events·SSE 감사와 최소 수정
status: draft
updated: 2026-10-14
owner: dev
tasks: [SAR-REFACTOR-002-DEV]
summary: callback crash 뒤 새 도착 wake를 복구하고 같은 agent의 반복 pending 검색을 줄였다.
---

# Remote 범위 완료 보고

## 범위와 기준

부모 과제 기준 SHA는 `8f9370bb62c9133c73e10ea2307b4cd696e10f7e`다. 같은 dev 워크트리의 미커밋 변경으로 검사했다. 부모와 다른 worker의 파일은 수정하지 않았다.

- 직접 확인: `remote_auth.go`, `remote_mcp.go`, `remote_events.go`, `text_events.go`와 관련 단위 검사. 공유 `store.go`의 transaction·current·sweep·principal 및 text fixture를 필요한 범위에서 확인했다.
- 근거: 현재 `handovers/to_dev.md`, `FULLOPS.md`, `project.md`, `rules/common` 세 규칙과 README, `rules/delegation.md`, `docs/system-design.md`, `docs/architecture.md`, 개발·운영 안내의 관련 절, D05·D10 Remote/SSE 절.
- 기준: fullops-common-0.3.3, ponytail full, codebase-design, ponytail-audit. 감사 뒤 근거 있는 구현만 적용했다. 새 공개 interface나 저장 형식은 만들지 않았다.
- 예상 부하: 설치당 최대 4개 구독, 총 128개 구독, 24시간 보존 receipt가 쌓인 상태. 벤치는 단일 설치의 1/4개 구독과 100/10000개 기록을 비교한다. 운영 부하를 실측한 값은 아니다.
- 새 라이브러리 API 확인은 필요하지 않았다. 기존 Go 제어 흐름과 호출 계약만 변경했다.

## 감사 결과와 기술 계획

### 1. 중간 심각도: 마지막 callback lease 만료 뒤 새 도착 wake 지연

원인은 `claimRemote`의 attempts 소진 분기다. 기존 Pending을 Seen으로 옮긴 뒤 다른 현재 pending이 있어도 구독 만료만 다음 시각으로 남겼다. 새 도착의 NOTIFY를 활성 lease 중 이미 처리한 경우 후속 wake가 사라진다.

재현은 동일 event를 t+0/16/32초에 claim하고 sender가 매번 finish 전에 중단되는 입력이다. t+40초에 두 번째 메시지가 도착한다. 이때 claim은 t+47초 lease 만료를 기다린다. t+47초에는 첫 event의 3회 한도만 소진해야 한다. 기존 코드는 새 도착 재확인 대신 구독 만료 t+3600초를 예약했다. 추가 NOTIFY나 재시작이 없으면 새 메시지의 180초 기한을 넘길 수 있다. 본문 수동 receive는 가능하지만 도착 신호가 지연된다.

최소 수정은 attempts 소진 뒤 현재 pending과 Seen이 다를 때 `Next=now`를 예약하는 것이다. 다음 loop가 현재 권한·기한을 다시 검사하고 별도 event를 시작한다. 동일 event의 최대 3회와 활성 lease를 유지한다. 변경은 `remote_events.go`의 4줄이다.

`TestRemoteExhaustedLeaseWakesNewArrival`는 변경 전 실제 실패했다. 수정 후 즉시 wake, 새 event ID, 새 event의 첫 시도, 활성 lease 보호, 성공 뒤 idle 대기를 확인한다.

### 2. 성능: 같은 agent의 구독마다 전체 message map 재검색

`claimRemote`는 같은 설치의 구독 4개에도 `eventPending`을 각각 호출했다. 해당 함수는 전달 완료된 기록을 포함한 모든 message를 순회한다. 외부 요청으로 wake가 발생할 때 이 작업은 singleton transaction 안에서 반복된다.

한 claim 호출의 agent별 pending 결과만 재사용했다. 함수가 반환되면 map은 버린다. 구독별 설치·기한·lease 검사는 유지한다. `eventPending`, SSE, 실제 전송 직전 권한 재검사는 변경하지 않았다. 메시지나 권한을 요청 사이에 캐시하지 않는다.

`TestRemoteClaimsRecheckCurrentAuthority`는 빈 queue 검사 후 새 도착을 인식하는지 검사한다. 첫 구독 성공 후 관계를 철회하면 두 번째 구독이 이전 권한을 재사용하지 않는지도 검사한다.

### 보존한 경계

OAuth의 등록 redirect·PKCE·resource·scope·code 일회성·refresh family 회전/철회·설치 키 검사는 그대로다. MCP의 서명 재검증과 원자적 persist/Inbox/transport ACK도 그대로다. 설치 Inbox의 32개 상한·명시 ACK·응답 유실 복구·현재 권한·TTL sweep은 변경하지 않았다.

callback의 HTTPS443·전체 DNS 응답 검사·IP pinning·redirect 거부·challenge·HMAC·secret 회전은 그대로다. 410/413 제거와 429/서버 실패 backoff, 최대 3회, DB lease도 그대로다. SSE의 LISTEN 선행·빈 payload·권한 확인·stream 한도·heartbeat는 수정하지 않았다.

이 범위에서 추가 확정 결함은 찾지 못했다. 이것은 전체 OAuth/MCP 표준 적합성이나 운영 보안 인증을 뜻하지 않는다.

## 실행과 종료코드

환경은 Windows amd64, Go 1.27.1, AMD Ryzen 5 5600이다. cwd는 dev 워크트리 루트다. 아래 로그는 삭제하지 않고 보존한다. PowerShell은 Go 출력을 `Tee-Object`에 기록한 뒤 `$LASTEXITCODE`를 즉시 저장하고 같은 값으로 종료했다.

| 검사 | 명령 | 원본 종료코드 | 근거 |
|---|---|---:|---|
| 기존 focused 회귀 | `go test ./internal/relay -run 'TestRemote\|TestTextEvent' -count=1` | 0 | 수정 전 tool 실행 결과, 0.658s |
| 결함 RED | `go test ./internal/relay -run '^TestRemoteExhaustedLeaseWakesNewArrival$' -count=1` | 1 | [red.txt](remote-evidence/red.txt), 예상 Next 불일치 |
| 수정 후 focused 회귀 | `go test ./internal/relay -run 'TestRemote\|TestTextEvent' -count=1` | 0 | [green.txt](remote-evidence/green.txt), 0.394s |
| 전후 벤치 | `go test ./internal/relay -run '^$' -bench '^BenchmarkRemoteClaimSubscriptions$' -benchmem -benchtime=200ms -count=3` | 각각 0 | [before.txt](remote-evidence/before.txt), [after.txt](remote-evidence/after.txt) |
| 포맷 | `gofmt -w internal/relay/remote_events.go internal/relay/remote_events_test.go` | 0 | 실행 후 제품 diff 확인 |
| 공백 | `git diff --check -- internal/relay/remote_events.go internal/relay/remote_events_test.go` | 0 | 부모 통합 검사와 별도 focused 검사 |

표의 정규식 `\|`는 Markdown 표 구분자를 escape한 표기다. 실제 명령은 `TestRemote|TestTextEvent`다.

벤치는 같은 코드의 fixture·컴파일러·데이터 수·반복 수를 사용했다. before에는 신규 벤치만 추가하고 제품 수정은 적용하지 않았다. 부모의 별도 sweep 변경은 이 벤치가 호출하지 않는다.

| 보존 기록 / 같은 설치 구독 | 이전 ns/op 중앙값 | 이후 ns/op 중앙값 | 이전/이후 B/op | 이전/이후 allocs/op |
|---|---:|---:|---:|---:|
| 100 / 1 | 1143 | 1438 | 0 / 0 | 0 / 0 |
| 100 / 4 | 4558 | 1754 | 64 / 64 | 1 / 1 |
| 10000 / 1 | 148226 | 156203 | 0 / 0 | 0 / 0 |
| 10000 / 4 | 590438 | 162572 | 64 / 64 | 1 / 1 |

4개 구독에서 100개 기록은 약 61.5%, 10000개 기록은 약 72.5% 감소했다. 1개 구독에서는 검색 재사용의 이득이 없으며 관측치는 각각 약 25.8%, 5.4% 증가했다. 짧은 비격리 실행의 변동과 map 조회 비용을 분리 측정하지 않았다. 모든 규모에서 개선됐다고 주장하지 않는다. 4개 구독의 반복 검색 제거가 목표다.

이는 DB·TLS·네트워크를 제외한 claim CPU 벤치다. 합성 delivered 기록은 map 순회 비용만 대표한다. 유효 queued 메시지의 current 검사, 여러 설치, 전역 JSONB 직렬화·행 lock 처리량은 이 측정에 포함하지 않았다.

## 변경 파일과 부모 인계

- 제품: `internal/relay/remote_events.go` — 10줄 추가, 1줄 삭제.
- 검사: `internal/relay/remote_events_test.go` — 회귀 2개와 벤치.
- 보고: 이 파일과 `remote-evidence/{red,green,before,after}.txt`.
- D10과 `docs/system-design.md`에 callback 마지막 lease 소진 뒤 새 pending 재예약, claim 호출 내 pending 재사용을 반영하도록 부모에게 제안한다. HTTP·저장 schema·제품 규칙은 그대로여서 D02/D04/D05/D07 계약 변경은 없다.

재위임, Git commit/branch 변경, 공유 인박스·packet·문서 수정, 삭제·이동 정리를 실행하지 않았다. 별도 의존성 추가나 외부 계정 발송·배포는 없다.

하위 과제의 범위는 focused 단위 회귀와 측정이다. 전체 lint/test/build·race, Linux 실제 DB integration, 고정 SHA 독립 리뷰, 커밋, FullOps gate, `work.py finish`는 부모의 통합 책임이다. 해당 검사를 이번 하위 실행에서 통과한 것으로 표시하지 않는다. UI 변경이 없어 직접 시각 검수는 이 diff에 적용하지 않았다. live OAuth, 공개 callback DNS/TLS, 실제 sender process crash와 다중 process DB lease는 실행하지 않았다. 상태 메서드 회귀가 그 환경 검증을 대체하지 않는다.
