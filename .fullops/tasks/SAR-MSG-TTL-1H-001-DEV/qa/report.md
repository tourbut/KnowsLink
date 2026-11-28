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
