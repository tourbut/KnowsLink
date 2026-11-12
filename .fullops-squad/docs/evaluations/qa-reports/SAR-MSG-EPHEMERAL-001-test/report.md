---
title: SAR-MSG-EPHEMERAL-001 메시지 휘발 저장 검증
status: draft
updated: 2026-11-12
owner: dev
tasks: [SAR-MSG-EPHEMERAL-001]
summary: 메시지·대화 사본·Dots Inbox를 Redis에만 두는 구현의 검증 결과와 남은 일
---

# SAR-MSG-EPHEMERAL-001 메시지 휘발 저장 검증

## 결론

- 메시지 기록(본문·envelope·Inbox 사본·웹 대화 사본·receipt·lease/claim·멱등)과 Dots 설치 Inbox를 Postgres에서 빼고 디스크 저장 없는 Redis에만 둔다. 최종 코드 SHA는 `6fe42d32bc96`(브랜치 `fullops/dev`)이다.
- 등록 검사와 실제 Postgres·Redis integration이 모두 종료코드 0이다. DATA RACE 0이다.
- 운영 적용·Redis 운영 기동·기존 메시지 삭제는 하지 않았다. 절차 초안은 [D12](../../../operations/ops-guide.md#sar-msg-ephemeral-001-메시지-휘발-저장-이행-운영-미적용)에 있다.

## 구조 요약

- Postgres 행에는 Redis 키 이름(`State.Ephemeral`, random token)만 남는다. Redis 값은 `{Messages, Idempotency, Inboxes}` JSON 한 개이고 TTL은 1시간이다.
- 기존 행 lock 안에서 Redis를 읽고, 바뀐 경우에만 새 키를 쓴 뒤 행을 commit하고 이전 키를 지운다. commit 전 실패는 새 키를 지운다. 그래서 여러 relay process의 lease·claim·ACK·멱등 직렬화는 바뀌지 않는다.
- Redis 오류·미설정은 `unavailable`이고 Postgres 대체 저장은 없다. `/healthz`는 Postgres와 Redis가 모두 응답해야 200이다.
- 모든 메시지 기록은 수락 1시간 뒤 sweep이 지운다(기존 24시간). text 전달 기한 180초와 `/v1/send` 300초, 20,000건·대기 100건 상한은 그대로다.
- 요청 사본이 Redis에서 사라진 gate는 sweep이 지운다. LISTEN/NOTIFY는 빈 payload라 유지했다. lock 없이 상태를 읽는 두 경로(cleanup snapshot, SSE 도착 신호)도 Redis 사본을 붙인다.
- migration 4가 기존 행의 `Messages`·`Idempotency`·설치 `Inbox`를 지운다. Down은 실패 고정이다.

## 등록 검사

실행 위치는 로컬 Linux 컨테이너(`kl-msgeph-linux`, Go 1.27.1), Postgres 17 격리 DB, `redis:8.10.2-alpine`(`--save ""`·`--appendonly no`·`noeviction`)이다. 대상은 `6fe42d32bc96`이다. 원시 로그는 `D:/workspace/KnowsLink/.devtools/msg-ephemeral/gate/6fe42d32bc96…/`에 있다.

| 명령 | 종료코드 |
|---|---:|
| `make install` | 0 |
| `make lint` | 0 |
| `make test` | 0 |
| `make build` | 0 |
| `python3 scripts/package_plugin.py --verify` | 0 |
| `go run ./cmd/migrate` | 0 |
| `go test -race -count=1 -tags=integration ./internal/... ./cmd/...` | 0 (DATA RACE 0) |
| `make verify-runtime` (실제 Compose, 대상 `6fe42d32`) | 0 |
| FullOps `lint.py --from 17acfa5e` | ERROR 0, WARNING 9, Windows `make` 실행 불가 2 |

`make verify-runtime` 로그는 `D:/workspace/KnowsLink/.devtools/msg-ephemeral/runtime/6fe42d32…/verify-runtime.log`다. 같은 과제의 이전 게이트(`ea72bbab`, `e077ed8a`)는 실패했다. 실패 원인과 수정은 아래 "발견하고 고친 것"에 있다.

## 완료 기준별 근거

| 기준 | 근거 |
|---|---|
| Postgres에 본문·envelope·사본이 없음 | `TestMessagesLiveOnlyInRedis`: envelope·Inbox·웹 사본·Dots Inbox를 쓴 뒤 public 전체 테이블을 텍스트로 훑어 본문과 메시지 ID가 없음을 확인한다. `TestHumanPostgres`: 실제 웹 송신·agent 수신·서명 text 뒤 같은 검사를 한다. `TestStateEncodingSplitsMessages`: Postgres 직렬화에 `Messages`·`Idempotency`·`Inbox`가 없다. |
| 1시간·180초 뒤 삭제 | 같은 시험이 Redis 키 TTL이 59분 초과 1시간 이하임과 수락 1시간 된 기록이 다음 sweep에서 사라짐을 확인한다. 180초 경계는 기존 `TestPublicTextTTLAndLeaseFailures`가 그대로 통과한다. |
| Redis 재시작·장애 | 시험: 키 소실 뒤 빈 저장소와 새 기록 정상, 도달 불가 Redis·미설정 Redis는 `unavailable`이고 epoch 불변. `verify-runtime`: 실제 `redis` 컨테이너 stop 중 `/healthz` 503, start 뒤 200, 재시작 전 넣은 probe 키가 재시작 뒤 0, `REDIS_URL` 누락·도달 불가 시 relay 기동 거부. |
| 철회·친구 해제·키 철회 뒤 미전달 | `TestHumanPostgres`의 친구 해제 뒤 persist 409, `TestPublicTextRoundtripAndCurrentAuth`의 키·agent·회원·pair 변경 뒤 `current=false`가 Redis 경유 상태에서 통과한다. 검사 코드(`current`)는 바꾸지 않았다. |
| 여러 process 원자성 | `TestHTTPConcurrencyAcrossInstances`, `TestSharedExecutionClaimHTTP`, `TestSharedMessageCapacityHTTP`가 같은 Redis를 쓰는 두 Service로 통과한다. |
| migration 4 | `TestStateFormatMigration`: legacy 행(메시지 1,001건·멱등·Dots Inbox)에 00003 뒤 00004를 적용하면 세 값만 빠지고 나머지 JSON이 같다. `human_model` 1, lz4, guard, Down 실패를 확인한다. |
| 페이지 문구 | `TestPolicyPages`가 개인정보·약관 한국어·영어에서 "영구 저장하지 않습니다/never stored permanently", "최대 1시간", 재시작, Redis·180초·백업 문장을 확인하고 24시간 문장이 없음을 확인한다. 회원 대화·랜딩 시험도 새 문장을 확인한다. |

## 문구가 사실인지

- 로그: relay의 `log` 호출은 기동 설정, 정리 실패, agent 용량 경고, 메일 발송 실패뿐이다. 본문·envelope을 쓰지 않는다. 오류 응답은 고정 코드다.
- 새 백업: `pg_dump`는 Postgres만 담고 Redis는 백업하지 않는다. `beta.sh backup`이 `message_state_rows`를 출력한다. 이 조회는 이전 `data`(jsonb)와 새 `state`(json) 열 모두에서 맞게 동작함을 격리 Postgres로 확인했다(메시지 있음 1, 없음 0).
- 이전 운영 백업: 메시지를 담고 있다. `prune-backups.sh`(D021)는 수정 시각 14일이 지난 `*.dump`를 매일 03:17 KST에 지운다. 따라서 이행 전 dump는 만든 날부터 최대 약 15일 남는다. 단 가장 최근 dump 1개는 지우지 않는다. 이행 뒤 새 dump를 만들지 않으면 이행 전 dump가 계속 남는다. 홈서버의 `c6dfc39a-pre-vultr-20261112T124516Z.dump`도 같다. 삭제는 하지 않았다.
- 개인정보처리방침의 "백업에는 메시지를 넣지 않습니다"는 이행 뒤 사실이다. 이행 전 백업이 메시지를 담을 수 있다는 문장을 같은 단락에 두었다.
- 서명 agent 업무 envelope(`/v1/send`)는 기존대로 최대 300초다. 개인정보처리방침은 이를 따로 적었다. 회원 화면·약관의 180초는 text 대화 기준이다.

## 발견하고 고친 것

- `ea72bbab`: compose Prettier 형식 오류, 시험 2건(현재 Redis 키를 다시 읽지 않음, Postgres 직렬화로 envelope을 복원). 시험과 형식만 고쳤다.
- `e077ed8a`: SSE 도착 신호가 lock 없이 Postgres 행만 읽어 Redis 메시지를 보지 못했다(실제 결함, `TestTextEventsHTTP`). `readEvents`가 Redis 사본을 붙이고, 이미 교체된 키면 행을 다시 읽는다.
- 자체 점검으로 commit 전 실패 때 새 Redis 키가 1시간 고아로 남는 문제와, 요청 취소 때 이전 키 삭제가 실패하는 문제를 고쳤다.

## 스스로 정한 선택과 되돌리는 방법

- Redis 키 하나에 전체 메시지 기록을 둔다(Lua·WATCH 없음). 행 lock을 그대로 쓰므로 보장이 같다. 되돌리기: `ephemeral.go`와 `transaction`의 Redis 호출을 키 단위 저장으로 바꾼다.
- Redis 장애 때 메시지 외 작업(로그인·친구)도 `unavailable`이다. 모든 transaction이 sweep에서 메시지를 보기 때문이다. `/healthz`도 503이다. 되돌리기: `main.go` health 함수에서 Redis ping을 뺀다(작업 실패는 그대로).
- Redis 비밀번호 없이 internal network만 쓴다. `database` 망에는 postgres·migrate·relay·redis만 있다. 되돌리기: `--requirepass`와 `REDIS_URL` 자격을 추가한다.
- `maxmemory 128mb`, 운영 `mem_limit 192m`·`cpus 0.25`. commit 중 사본 2개와 relay 256 MB 한도를 근거로 정했다. 실제 크기는 운영 수락 때 측정한다.
- 정책 문서의 시행일은 바꾸지 않았다. 개인정보처리방침 변경 고지 여부는 제품·법률 판단이다.

## 의존성 (DEP-001)

- `go.mod`·`go.sum`: `github.com/redis/go-redis/v9 v9.23.0` 추가와 간접 의존성(`cespare/xxhash/v2`, `dgryski/go-rendezvous`, `go.uber.org/atomic`), `golang.org/x/sys` 0.48.0. 필요성과 표준 라이브러리 대안 부재는 [D03 기술 스택](../../../design-docs/tech-stack.md)에 적었다. 런타임 이미지는 `redis:8.10.2-alpine`(digest 고정, 2026-10-21 게시)다.

## lint WARNING 설명

- SIZE-001 `store.go` 569줄, `ops-guide.md`, `PLANS.md`(main 병합분): 기존 초과 파일이다. 이번 추가는 Redis 연결 몇 줄이고 새 책임은 `ephemeral.go`로 분리했다.
- SIZE-002 추가 909줄: 지시서 예상 600~1,200줄 안이다. 시험·문서가 절반 이상이다.
- SEC-001 `ephemeral_test.go:32`: 변수 이름 `secret`의 오탐이다. 값은 합성 대화 문장이다.
- SLOP-004 세 스크립트의 `print`: 검사 결과 출력이다.

## Docker 자원

`kl-msgeph-` 컨테이너·볼륨(`gomod`, `gocache`)·네트워크·이미지(`kl-msgeph-linux`), 내려받은 `redis:8.10.2-alpine`, `verify-runtime`이 만든 `knowslink-check-*` 이미지를 모두 지웠다. 남긴 자원은 없다. 시험을 위해 꺼져 있던 Docker Desktop을 켰고 켜진 상태로 두었다.

## 남은 일

- ops: 사용자 승인 뒤 Vultr 운영 이행(Redis 기동, migration 3·4, 기존 메시지 삭제, 이행 직후 새 백업).
- tester: 운영 이행 전 독립 QA(웹 대화·Dots 수신·재시작 체감).
- 사용자: 이전 백업 보존 정책 판단, 개인정보처리방침 개정 시행일·고지 여부.
