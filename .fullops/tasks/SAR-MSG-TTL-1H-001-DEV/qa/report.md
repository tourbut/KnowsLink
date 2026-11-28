---
title: SAR-MSG-TTL-1H-001-DEV QA 보고
status: draft
updated: 2026-11-28
owner: dev
summary: "메시지 배달 기한 1시간 전환의 180초 구분, 정책 문구 비교, 용량 어림, MCP low 3건, 테스트 결과"
---

# SAR-MSG-TTL-1H-001-DEV QA 보고

- base: 112e8810ac7a3c025b6a7815e2fd850871b4d889
- 근거: 주인 결정 2026-11-28(쿠어 msg_01c39e2e0d7e). 모든 KnowsLink 메시지의 relay 보관·배달 기한을 1시간으로 한다.
- 테스트 레벨: lite

## 1. 180초(및 300초) 구분 표

| 파일:줄(base 기준) | 뜻 | 바꿈 |
|---|---|---|
| internal/relay/public_text.go:17 `publicTextTTL = 180s` | text·웹 owner text·원격 MCP send 배달 상한 | 예 → `messageRetention`(1h) |
| internal/relay/store.go:549 `300 * time.Second` | `/v1/send` 업무 envelope(relay.request·relay.result 등) 배달 상한 | 예 → `envelopeTTL(intent)` = 1h |
| 같은 줄, `relay.approval.request` | 사람 승인 gate 대기(Gate.Exp = 승인 요청 exp) | 아니오. `approvalTTL = 300s` 유지 |
| internal/relay/ephemeral.go:15 `messageRetention = time.Hour` | 기록 보관·Redis 키 TTL·sweep | 값 유지, 배달 상한의 단일 상수로 지정 |
| internal/relay/remote_mcp.go:223 `now.Add(publicTextTTL)` | 원격 MCP send exp | 상수 따라 1h(파일은 remote_tools.go로 이동) |
| internal/relay/human_messages.go:62 | 웹 발송 exp | 상수 따라 1h |
| adapters/src/text.ts:130, 165 `170000` | 로컬 adapter 발송 기한(MCP send, orca send) | 예 → `SEND_TTL_MS = 3590000`(10초 여유 유지) |
| adapters/src/mcp.ts:418 "TTL 180s" | 로컬 MCP 도구 설명 | 예 → "TTL up to 1 hour" |
| adapters/src/inbox.ts:17 주석 | read marker 근거(24h 보관) | 주석만 1h로. 값 24h 유지 |
| adapters/src/synthetic.ts:86, test-transport.ts:56 `180000` | 합성·test-remote 발송이 고른 TTL(상한 아래) | 아니오. 상한 안이라 유지 |
| human_pages.go:12, 62 / i18n_en_human.go:19, 229 | 보존 안내·composer 안내 | 예 |
| member_pages.go:95, 124 / i18n_en_member.go:114, 158 | receipt TTL 상세·연결 확인 안내 | 예 |
| landing_pages.go:50 / i18n_en_member.go:233 | 랜딩 전달 기한 | 예 |
| policy/privacy.ko.html:29, privacy.en.html:29 | 개인정보 보관 기한 | 예(아래 비교) |
| remote_mcp.go:159 receive 설명 "kept at most 1 hour" | 원격 Inbox 보관 | 문구 유지. 이제 실제와 맞음 |
| remote_era.go:69 `remoteCacheTTL = 300000` | MCP discover/list 캐시 | 아니오(다른 뜻) |
| inbox.ts:263 `300000`, mcp.ts:74 `timeout: 300000` | 수신 backoff·MCP 호출 timeout | 아니오(다른 뜻) |
| videos/**, `180°` 셔터 | 영상 | 아니오(무관) |
| user-guide.md:233 "gate는 생성 후 180초에 만료" | 사람 승인 gate | 아니오 |
| 테스트 픽스처 `now.Add(180*time.Second)` 다수 | 상한 안의 발송 예시 | 아니오(여전히 유효) |

문서(README, adapters/README, SKILL.md, docs/*, D03, D05, D07, D09, user-guide, ops-guide)의 180초·300초 배달 기한 문구는 1시간으로 고쳤다. `.fullops/docs/planning/**`(과거 기획 원천)·평가 시나리오·archive는 이력이라 고치지 않았다.

## 2. 정책 문구 이전·새 비교

| | 이전 | 새 |
|---|---|---|
| privacy.ko | 영구 저장하지 않습니다. 서버 메모리(Redis)에만 두고 자동 삭제합니다. 전달 대기는 최대 180초(에이전트 업무 요청은 최대 300초), 그 밖의 기록은 서비스가 받은 때부터 최대 1시간입니다. 서버가 재시작하면 그 전에 사라질 수 있습니다. | 영구 저장하지 않습니다. 서버 메모리(Redis)에만 두고 자동 삭제합니다. 전달 대기와 그 밖의 기록 모두 서비스가 받은 때부터 최대 1시간입니다. 서버가 재시작하면 그 전에 사라질 수 있습니다. |
| privacy.en | Never stored permanently. Kept only in server memory (Redis) and deleted automatically: pending delivery lasts at most 180 seconds (agent task requests at most 300 seconds), and other records at most 1 hour after the Service receives them. A server restart can remove them sooner. | Never stored permanently. Kept only in server memory (Redis) and deleted automatically: pending delivery and all other records last at most 1 hour after the Service receives them. A server restart can remove them sooner. |
| 웹 보존 안내(ko) | 전달 대기는 최대 180초, 웹 대화 사본은 수락 시각부터 최대 1시간 서버 메모리(Redis)에만 두고 자동 삭제합니다. | 전달 대기와 웹 대화 사본은 수락 시각부터 최대 1시간 서버 메모리(Redis)에만 두고 자동 삭제합니다. |
| 랜딩·composer | 전달 기한은 최대 180초입니다. | 전달 기한은 최대 1시간입니다. |

terms(ko·en)는 이미 "최대 1시간 안에 자동 삭제"라 바꾸지 않았다. 보관 위치(Redis)·삭제 방식·암호화는 그대로다. receipt·로그에 본문을 넣는 경로는 추가하지 않았다.

검토 필요: 개인정보 방침 59행은 "변경 시 시행 7일 전 알림, 불리한 변경 30일 전"이다. 원문 보관 시간이 길어지는 변경이라 `policyEffective`(2026-10-22) 갱신과 고지 기간은 주인·ops가 정한다. 이 과제에서는 날짜를 바꾸지 않았다.

## 3. 용량 어림

입력값(코드): text ≤4096 B, raw envelope ≤32768 B(`rawEnvelopeLimit`), HTTP 새 작업 member당 40/분·전체 200/분(시간당 최대 12000건), 기록 수 `len(Messages) < 20000`, 미배달 큐(queued·leased) `< 100`, 원격 Inbox 설치당 ≤32, 기록 보관 1h(이전과 같음).

- 기록·웹 대화 사본: 이미 1시간 보관이었다. 이 과제로 바뀌지 않는다. 최대 12000건 × (사본 ≤4 KiB + 메타 약 1 KiB) ≈ 60 MiB.
- 새로 늘어나는 것은 원문 payload(Envelope·Inbox 사본)의 수명이다. 180초 → 최대 1시간(×20).
  - 미배달 큐: ≤100 × 32 KiB ≈ 3.2 MiB. 상한은 같고 머무는 시간만 길다.
  - persist 후 ACK 안 한 원문: 큐 상한 밖이다. 이론 최악 12000 × (Envelope+Inbox 2부) — 보통 text ≈4.5 KiB면 약 105 MiB, escape로 32 KiB 꽉 차면 약 750 MiB. 이전(180초)은 같은 조건에서 약 600건 ≈ 5–37 MiB였다.
  - 원격 Inbox: 설치당 ≤32 × 32 KiB = 1 MiB.
- 운영 비교: 저장소 compose는 Redis `--maxmemory 128mb`, `noeviction`, ops override `mem_limit: 192m`, relay `mem_limit: 256m`. 주석대로 commit 중 2부를 담으므로 실효 약 64 MiB다. 운영 서버(Vultr knowslink-relay-01)의 실제 maxmemory·사용량은 **미확인**이다(접속하지 않음).
- 판단: 정상 사용(ACK까지 바로 진행)에서는 영향이 작다. 다만 이론 최악(다수 member가 보내고 수신자가 persist 후 ACK 안 함)에서는 128 MiB를 넘을 수 있고, noeviction이라 쓰기가 `unavailable`이 된다. 미배달 큐 100건 전역 상한도 이제 최대 1시간 점유될 수 있어, 오프라인 수신자에게 보낸 100건이 1시간 동안 새 발송을 `capacity`로 막을 수 있다(이전 180초). 상한·정책 변경이라 이 과제에서는 고치지 않았다. 후속 후보: 발신자별 큐 상한, persist-미ACK 원문 상한, 운영 Redis 사용량 관측.

## 4. 리뷰 low 3건(SAR-MCP-UPGRADE-001)

- modern `events/subscribe`·`events/unsubscribe` 성공 결과에도 `resultType`·`_meta.serverInfo`를 붙인다(`remoteReply` 헬퍼, remote_era.go). D05 서술과 일치.
- `Mcp-Name`은 `params.name`을 가진 메서드(`tools/call`, `prompts/get`, `events/subscribe`, `events/unsubscribe`)에서만 비교한다. `resources/read`는 404 `-32601`.
- remote_mcp.go 325줄 → 141줄. 도구 정의·호출은 remote_tools.go(188줄)로 옮겼다(코드 이동만).

## 5. 테스트 발췌

```
go test -count=1 ./internal/...
ok  	github.com/tourbut/KnowsLink/internal/relay	2.945s

TEST_SYNTHETIC_DATABASE=1 TEST_DATABASE_URL=<격리 postgres> TEST_REDIS_URL=<격리 redis> \
  go test -count=1 -tags=integration -run 'TestPostgresSafety|TestRemote' ./internal/relay/
ok  	github.com/tourbut/KnowsLink/internal/relay	4.802s

adapters: npm run build && node dist/mcp.test.js && node dist/orca.test.js && npm run lint
PASS: MCP initialize/discovery, ... no secrets in tool output
PASS: orca schema, ... one watch per folder
eslint: 오류 없음
```

표적 단언:
- `TestPublicTextTTLAndLeaseFailures`: exp 180s·1h 수락, 0·1h+1s 거부(`ttl_too_long`).
- `test_messages_test.go`: `/v1/send` 1h+1s 거부, `envelopeTTL` = 1h, 승인 요청 300s 유지.
- `integration_test.go`: `/v1/send` exp 1h+1m → 422.
- `adapters/src/mcp.test.ts`: `SEND_TTL_MS == 3590000`, 실제 발송 exp가 3580–3590초 안.
- 문구: landing_test, human_pages_test, policy_pages_test를 새 문구로.
- low: remote_era_test `resources/read` + 다른 `Mcp-Name` → 헤더 검사 통과. remote_integration_test modern `events/unsubscribe` → `resultType`·serverInfo, `resources/read` → 404 -32601.

integration 컨테이너(postgres·redis)는 실행 뒤 지웠다. 배포 뒤 Grok 실제 확인은 범위 밖(쿠어).

## 재작업 1차 (D076, 2026-11-28)

1차 `0410f759`의 "모든 메시지 1시간"을 주인 결정 D076대로 좁혔다. 위 1–5절은 1차 기록이고, 아래가 현재 규칙이다.

### R1. 규칙과 코드 위치

| 경로 | 자기 agent 사이(같은 owner, self 권한) | 그 밖(친구·legacy·다른 owner) | 1시간 초과 |
|---|---|---|---|
| text `/v1/text/send`(`ingestPublicText`) | 최대 1시간 | 180초. 180초–1시간 요청은 180초로 잘라 받음(clamp) | `ttl_too_long` |
| 업무 envelope `/v1/send`(`ingest`) | 최대 1시간 | 300초(base와 같음), 넘으면 `ttl_too_long` | `ttl_too_long` |
| `relay.approval.request` | 300초 | 300초 | `ttl_too_long` |
| 원격 MCP send·웹 owner 발송 | exp = now+1시간 | exp = now+180초 | — |

- 판정: `State.selfAgents(from, to)`(grants.go). 두 agent의 owner가 같고 `authority`가 `self:` scope일 때만 참이다. owner grant(`grantLive`·`createGrant`)도 이 함수를 쓰도록 바꿨다. 보내는 쪽 입력은 보지 않는다.
- 상한: `deliveryTTL(from, to, peer)`(public_text.go). 자기 agent면 `messageRetention`(1h), 아니면 `peerTextTTL`(180s)·`peerEnvelopeTTL`(300s). 1차의 `publicTextTTL`·`approvalTTL`은 지웠다.
- clamp: 서명된 envelope의 exp는 그대로 두고 receipt의 `Exp`만 줄인다. sweep·lease·답장 부모 기한은 receipt 값을 쓴다.
- adapter: `SEND_TTL_MS = 3590000` 유지(1시간 10초 전). 받는 쪽이 친구면 relay가 180초로 자른다. 이전 클라이언트의 170초는 그대로 받는다.
- 날짜 분기 코드는 없다. 2026-12-28 전환은 별도 과제다.

### R2. 정책 문구 이전(1차)·새

| | 1차(병합 안 함) | 새 |
|---|---|---|
| privacy.ko 29행 | 전달 대기와 그 밖의 기록 모두 서비스가 받은 때부터 최대 1시간입니다. | 전달 대기는 내 에이전트끼리 최대 1시간, 그 밖은 최대 180초(에이전트 업무 요청은 최대 300초)입니다. 그 밖의 기록은 서비스가 받은 때부터 최대 1시간입니다. 서버가 재시작하면 그 전에 사라질 수 있습니다. 2026-12-28부터 모든 메시지의 전달 대기가 최대 1시간으로 바뀝니다. |
| privacy.en 29행 | pending delivery and all other records last at most 1 hour … | pending delivery lasts at most 1 hour between your own agents and otherwise at most 180 seconds (agent task requests at most 300 seconds), and other records at most 1 hour … From 2026-12-28, pending delivery for every message will last at most 1 hour. |
| 웹 보존 안내 | 전달 대기와 웹 대화 사본은 … 최대 1시간 | 전달 대기는 내 에이전트끼리 최대 1시간, 그 밖은 최대 180초입니다. 웹 대화 사본은 … 최대 1시간 |
| 랜딩 | 전달 기한은 최대 1시간입니다. | 전달 기한은 최대 180초, 내 에이전트끼리는 최대 1시간입니다. 2026-12-28부터 모든 메시지가 최대 1시간입니다. |
| composer·receipt 상세·연결 확인 | 최대 1시간 | 180초, 내 에이전트끼리 1시간 |

base(2026-10-22) 문구와 비교하면 privacy는 "내 에이전트끼리 최대 1시간"과 변경 예정 고지 한 문장만 늘었다.

`policyEffective`: 2026-10-22 → **2026-11-28**. 근거: 바뀌는 조항은 자기 agent 사이 1시간뿐이다. 주인이 D076에서 즉시 시행으로 정했다. 친구·다른 owner 사이(불리할 수 있는 변경)는 바꾸지 않고 2026-12-28 예정으로만 고지했다(30일 전). 배포일이 다르면 ops가 이 상수만 배포일로 바꾼다(되돌리기: `policy_pages.go`의 값).

도구 설명: 원격 `knowslink_send`에 "Delivery waits at most 1 hour between the owner's own agents, else 180 seconds."를 넣었다. 로컬 `knowslink_text_send`는 "TTL up to 1 hour between your own agents, else the relay cuts it to 180s". README·user-guide·SKILL.md·adapters/README·docs·D03·D05(interface/crud/database/architecture)·ops-guide도 같은 규칙으로 고쳤다. README·user-guide에는 2026-12-28 예정 한 줄을 넣었다.

### R3. 용량 어림 다시(검토 필요)

- 1시간 원문 수명은 자기 agent 사이 메시지만 받는다. 친구·다른 owner 사이 원문은 base와 같은 180초(envelope 300초)다. 그 몫은 base 어림(약 600건 ≈ 5–37 MiB)으로 돌아간다.
- 자기 agent 사이 몫의 이론 최악은 1차와 같다. 한 member 40건/분 × 60 = 2400건/시간이고, 여러 member가 동시에 하면 전체 200건/분 상한(12000건/시간)까지 간다. 보통 text ≈4.5 KiB면 2400건 ≈ 21 MiB(Envelope+Inbox 2부), 12000건 ≈ 105 MiB다. 32 KiB를 채우면 각각 약 150 MiB·750 MiB다.
- 위험이 줄어든 점: 이 몫을 쓰려면 보내는 쪽과 받는 쪽을 같은 사람이 가져야 한다. 남의 agent를 오프라인으로 둔다고 남의 원문이 1시간 쌓이지는 않는다. orca 고정 양식은 grant 30건/시간 상한도 받는다.
- 남은 위험: 미배달 큐 전역 100건 상한은 그대로다. 한 member가 꺼 둔 자기 agent에 100건을 보내면 1시간 동안 모두의 새 발송이 `capacity`가 될 수 있다(친구 사이만이면 180초). Redis 128 MiB·noeviction(저장소 compose 기준)도 같다. 운영 Redis maxmemory·사용량은 **미확인**이다.
- 큐·Redis 상한 정책은 바꾸지 않았다. 후속 후보: 발신 member별 미배달 큐 상한, persist 후 미ACK 원문 상한, 운영 Redis 사용량 관측. 2026-12-28 전 모든 메시지 1시간 전환 과제에서 함께 보는 것을 권한다.

### R4. 테스트

```
go vet ./internal/relay/ && go test -count=1 ./internal/...
ok  	github.com/tourbut/KnowsLink/internal/config
ok  	github.com/tourbut/KnowsLink/internal/relay

adapters: npm run build && node dist/mcp.test.js && node dist/orca.test.js && npm run lint
PASS: MCP initialize/discovery, held without network/key reads, invalid config, redirect rejection; no secrets in tool output
PASS: orca schema, … one watch per folder
eslint: 오류 없음
```

표적 단언:
- `TestDeliveryCeilingBySelfAgents`(grants_test.go, 새): grok→orca(자기) 1h 요청 → receipt 1h. grok→peer(친구) 1h 요청 → receipt 180s(clamp). 친구에 170s(이전 클라이언트) → 170s 그대로. 둘 다 1h+1s → `ttl_too_long`. envelope 상한: 자기 1h, 친구 300s, 자기 사이 승인 요청 300s. `SelfConnect=false`면 자기 판정이 거짓이고 180s.
- `TestPublicTextTTLAndLeaseFailures`: legacy pair는 0 거부, 180s·1h 수락이되 receipt 180s, 1h+1s `ttl_too_long`, 180s 시점 `failed:expired`.
- `test_messages_test.go`: 다른 owner `/v1/send` 301s 거부(300s 유지). `integration_test.go`의 301s → 422 단언은 base로 되돌렸다.
- `adapters/src/mcp.test.ts`: `SEND_TTL_MS == 3590000`, 실제 exp 3580–3590초(1차 그대로).
- 문구: landing_test, policy_pages_test(ko·en 새 문구와 2026-12-28 고지), human_pages_test(composer "최대 1시간" 포함).

못 한 것: integration 태그 테스트(postgres·redis)는 이번에 다시 돌리지 않았다. 바뀐 integration 단언은 base로 돌린 301s 한 줄과 상수 이름 한 줄뿐이다. 배포 뒤 Grok 확인은 범위 밖(쿠어).
