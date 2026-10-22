---
title: Silent Agent Relay — 결정 로그
status: draft
updated: 2026-10-22
owner: designer
tasks: [SAR-HUMAN-PEER-001, SAR-COMPETE-001]
summary: service-design be7c00e 반입 원문이다. 2026-10-19 메신저 재정의와 2026-10-22 포지셔닝 잠금을 포함한다
---

# Silent Agent Relay — 결정 로그

## 2026-10-02

### Grill points 1–7 — locked

1. Agent signup + pubkey; owner revoke/rotate  
2. Invite A→B; B-human must accept (auto only for already-paired reinvite) — *2026-10-19 superseded → 사람 친구 요청·상호 수락 (legacy agent pair만 유지)*  
3. Short TTL queue only; expired → `failed: expired`; ~~no long-term chat history~~ — *2026-10-19 amended → 웹 대화 사본 24h, 그 이상 아카이브 없음*  
4. Deny-by-default silent ops; human required for payment, delete, grants, outbound sends, commitments  
5. Core = HTTPS relay; plugins pull-default, webhook if available — *2026-10-19: 서명된 도착 metadata webhook 허용*  
6. Idempotency: same key+payload → cache hit; conflict → require new key  
7. Contacts = agents; A picks target agent (humans are owners only) — *2026-10-19 superseded → 친구 단위 = 사람(회원), 대화 단위 = agent 두 개*

### Business model — locked direction

- Tailscale-like freemium: personal free + friend/pair slot limits; Personal Pro unlocks more slots.
- Frame as “friends up to N”. Exact **N** and Pro price: **TBD**.
- Message-volume / Slack-seat / marketplace framing omitted for now.
- Slot-unit refinement deferred; **do not invent N/price**.

### Protocol / product drops

- **`deliver: both` dropped.** Only `agent` | `human` (default `agent` at app layer; **wire must set explicitly**).
- Human-side NL is a **render** of the structured envelope (`render.hint`), not the wire format.
- Invent thin `relay.v1` (ASMTP-shaped mailbox + A2A Task–inspired); not a full ASMTP/A2A clone.

### A2A gap review — locked (2026-10-03)

- A2A gap review is complete against the public **A2A v0.3.0 concepts**. This is not a wire-compatibility claim; the delta from any later/latest revision was **not rechecked**.
- Verdict: `relay.v1` is **Task-inspired, not an A2A clone**. **No A2A wire is imported.**
- Do **not** add `taskId`, `contextId`, `parts`, `artifacts`, or A2A state enums to `relay.v1`.
- Do not treat `role:user`, an AgentCard, or push as owner approval. Push is not a reason to enable the MVP webhook; the MVP webhook remains OFF. *(2026-10-19 amended: 서명된 도착 metadata webhook만 허용. A2A push를 근거로 연 것이 아니며, push·wake ≠ owner approval 원칙은 유지.)*
- Optional later documentation only (not written now): intent capability notes, an out-of-envelope observation view, and per-intent result schema examples.

### Scope

- **Hospital / closed-net out of scope.** Personal public-net only. No self-host hospital assumptions.

### Build shape

- **Relay + plugins/adapters.** Do **not** modify vendor agent cores.
- Do not claim those products expose official cross-agent inbound APIs.

### Non-goals recorded

- Latent KV / token-id LLM handoff: research-only, not MVP.
- Auth tokens are not message body.

### relay.v1 schema freeze — locked

Source: ChatGPT grill 「relay.v1 스키마 비평」 + 서비스 디렉터 정리. Docs: `protocol.md`.

- Envelope: `v,id,from,to,intent,body,deliver,exp,sig,idempotency_key` required; optional `priority,evidence,reply_to,trace,render,ext`.
- Reject: `kind`, `body.op`/`args`, request-envelope `status`, canonical `body.text`, `deliver:both`, unknown fields outside `ext`.
- Relay owns transport state outside envelope (`queued→leased→delivered` / expired / max_attempts). Agent completion via new signed `relay.result` (`reply_to`).
- Escalate = **new** signed `deliver:human` message (no mutate of original).
- Signing: RFC8785 JCS + Ed25519; domain-separated prefix `SILENT-AGENT-RELAY\\0relay.v1\\0Ed25519\\0` + kid + JCS(E).
- Idempotency scope `(from, idempotency_key)` + semantic digest (excludes `id,exp,sig,trace,render`).

### Operational constants — locked

| Constant | Value |
|----------|-------|
| `MAX_TTL` | 300s (typical request 60–120s) |
| Idempotency / receipt retention | 24h |
| Pull lease | 30s; `lease_until = min(exp, lease_granted_at + 30s)` |
| Delivery attempts | 3 (increment on successful lease grant only; 3rd lease usable in its window) |
| Intent registry | central, ADD-only revisions; name remap forbidden; breaking → new intent/version |
| Key revoke metadata retention | ≥ 24h (idempotency retention); keep revoke history |

### Round-2 ops / intent-seed freeze — locked

Source: ChatGPT Extra High grill round-2. Docs: `protocol.md`, `architecture.md`. **BM N/price unchanged (TBD).**

1. **No parent-link top-level field.** Escalation correlation = `reply_to` + typed body digest only.
2. **Intent registry seed:** `schedule.query`, `schedule.commit` (stub, not executable), `relay.result`, `relay.approval.request` (`deliver=human`, `reply_to` required, body `{reason, request_digest}`, does not execute tools). `schedule.query` still needs silent disclosure policy (TBD).
3. **Intent names immutable;** registry ADD-only; receivers implement subset, never redefine locally.
4. **`delivered` = after durable inbox persist + ACK** (not on lease grant). Sequence: lease → durable inbox → ACK → transport delivered → process.
5. **Lease / attempts:** `lease_until = min(exp, lease_granted_at + 30s)`; attempts++ only on successful lease grant; 3rd lease usable for its window.
6. **Ingest order:** structure/sig/auth/routing → semantic digest → atomic idempotency → exp/MAX_TTL/id collision → queue+receipt same txn. Cache hit after auth returns **receipt** (not full envelope): `id, from/to, intent, digest, exp, accepted_at, transport state`. No re-enqueue / no TTL extend on replay.
7. **Escalation H rules:** `intent=relay.approval.request`, `deliver=human`, `reply_to=M.id`, `from=to=M.to`, `body.request_digest=semantic_digest(M)`, `H.exp≤M.exp`; no recursive escalate from result/approval; B→B only as narrow exception tied to verified parent A↔B pairing + current B-owner accept (not pairing restore on revoke).
8. **Approval record outside envelope;** agent sig / `render.hint` ≠ owner approval proof; **amend out of MVP**.
9. **Relay transport failures are NOT forged as B-signed `relay.result`.**
10. **Reject:** NL wire, `deliver:both`, inventing rate-limit / extra envelope-size numbers, BM slot-unit / N / Pro price changes.

### Security MUST C1–C5 — locked

Source: ChatGPT Extra High security grill (`relay-v1-security-grill.reply.md`). Docs: `protocol.md`, `architecture.md`, `product.md`, `mvp-checklist.md`. **BM N/price/slot-unit untouched (TBD).** Freeze 유지하되 C1–C5 잠금 전에는 실데이터 silent 경로를 열지 않음.

**C1 Identity + current auth — Adopt**
- Key lookup = registered `(from, kid)` only; no URL/path kid; no kid rebind.
- Register/rotate = owner auth + PoP. Agent cred ≠ owner accept/approve/revoke.
- Re-check current key/pair/owner at enqueue · lease/ACK · approval consume · tool exec · result disclosure.
- Authorization epoch CAS; unpair re-accept = new pair generation; unknown auth state → block.
- Policy enforced at adapter/tool gateway (not model prompt).

**C2 Approval binding — Adopt**
- Human-gate UI from verified M typed body + policy, **not** `render.hint`.
- Approval record outside envelope; no GET-link approve; same digest ≠ reuse across `M.id`.
- Self B→B only with authoritative parent + actual B receipt; suppress duplicate pending gates per parent.
- Approve ≠ enable `schedule.commit` stub.

**C3 Atomic ingest + shared exec claim — Adopt**
- Locked order: `structure/sig/auth/routing → digest → atomic idempotency → exp/TTL/id → queue+receipt`.
- Fix: discard `product.md` legacy `validate(exp)→idempotency→auth`.
- Receipt-only replay; rollback idempotency reservation on TTL/id failure; no pre-auth cache probe.
- Lease token bound to recipient/msg/generation; durable persist → ACK success → shared exec claim → process.
- No multi-adapter without shared coordination.

**C4 Disclosure + result allowlist — Adopt**
- `schedule.query` default private; no policy = deny (scope numbers TBD).
- `relay.result`: `R.reply_to==M.id`, `R.from==M.to`, `R.to==M.from` + per-intent output allowlist.
- Data minimization on inbox/logs/traces/model context.

**C5 Boundaries — Adopt**
- Strict JSON; same parse object for verify+exec.
- **MVP OFF:** evidence auto-fetch/preview; optional webhook. *(2026-10-19 amended: webhook은 서명된 도착 metadata 알림만 허용; evidence OFF 유지)*
- Resource/rate limits TBD (never unlimited). Pending invites do not consume active slots.

**Defer:** evidence remote fetch; ~~optional webhook~~ (2026-10-19 도착 알림 한정 허용); `schedule.commit` executable; amend.  
**Reject:** wire/schema bypasses; auth-meaning bypasses; ops-invariant breaks; “HTTPS ref is safe”; invent Free N / Pro price / rate / extra size numbers; BM slot-unit change.

### Tech stack / deploy — locked

Source: user 성철 신 (2026-10-02). Docs: `architecture.md`, `product.md`, `mvp-checklist.md`, `README.md`.

- **MVP hosting:** home **mini-server** (self-host process). **No paid cloud app server** for MVP.
- **Public HTTPS:** **Cloudflare Tunnel** as **edge ingress only** (TLS/public entry → local relay). Public host is ~~`relay.knowslog.com`~~ **`link.knowslog.com`** (2026-10-19 superseded) on `knowslog.com`; Tunnel ingress routes it to the home mini-server. Tunnel is not the application runtime.
- **Reject as MVP hosting assumption:** Cloudflare Workers / Durable Objects (and similar serverless relay cores). May revisit later; not the current plan.
- **Database:** **Postgres** — locked (keys/pairs/intent registry/queue/lease/receipts/approval records; atomic txn requirements). **Not SQLite.**
- **Relay core:** **Go** — the HTTPS relay process/runtime.
- **Adapters/plugins:** **TypeScript** — including MCP pull clients.
- **First MVP adapters (locked, priority order):** **Grok Bot → Claude Code → Codex → Dots** (TypeScript, pull-default).
- These are adapter targets, not claims of official inbound APIs. No integration details are implied here.
- Unchanged: HTTPS relay core; personal public-net; no hospital/closed-net; no vendor core mods; pull-default adapters; webhook/evidence MVP OFF *(webhook: 2026-10-19 amended)*; `relay.v1` freeze; BM Free N / Pro price TBD (do not invent).

### Go MVP first implementation scope — locked

- 첫 구현은 ingest/queue-only가 아니라 **human-gate까지의 end-to-end 범위**다.
- 포함: validate/auth/ingest, short TTL queue, lease/ACK/shared exec claim, receipts, `relay.approval.request`, 그리고 envelope 밖의 approve/deny 경로.
- human-gate는 frozen `relay.v1` 및 C1–C5 규칙(검증된 typed body 기반 UI, 봉투 밖 approval record, current-auth 재검사)을 따른다.

### Mini-server configuration — deferred / Postgres tooling — locked

- 미니서버의 port·path·process manager·Cloudflare Tunnel 설정은 이 SoT에서 정하지 않는다. 향후 구현 레포의 env/config로 둔다.
- **Postgres tooling (locked):** Go driver는 `jackc/pgx/v5` + `pgxpool`을 사용한다.
- **Migrations:** `pressly/goose/v3`, SQL-only; 별도 `cmd/migrate`에서 실행하며 API 기동 시 auto-up하지 않는다.
- **Query codegen:** `sqlc` with `sql_package: pgx/v5`.
- **MVP에서 제외:** Atlas auto-diff/apply, ORM AutoMigrate, Go-code migrations as primary path.
### MVP deploy packaging — locked

- **Docker Compose**로 MVP를 패키징한다.
- Compose 서비스는 **`postgres`**, one-shot **`migrate`** (`cmd/migrate`), **`relay`** (Go), **`cloudflared`** (Tunnel ingress only)로 구성한다.
- Postgres는 퍼블릭에 노출하지 않는다. 포트·이미지/다이제스트·환경변수 등 구현 세부는 구현 레포에서 정하며, 이 SoT에서 발명하지 않는다.
- `migrate`는 `pressly/goose` SQL migration을 한 번 실행하고 종료하며, relay API 기동 시 auto-up하지 않는다.
- Go/Postgres 경계는 앞서 잠근 `pgx/v5` + `pgxpool` 및 `sqlc` (`sql_package: pgx/v5`)를 따른다.

### Human-gate UI — locked

- MVP human-gate UI는 Go relay core의 **`net/http` + `html/template`**로 구현한다.
- 별도 TypeScript frontend는 MVP human-gate에 사용하지 않는다.
- relay가 최소 approve/deny UI를 직접 서빙한다.
- UI는 C2에 따라 검증된 서명 메시지 **M의 typed body + 적용 정책**으로 구성한다. `render.hint`는 권위 있는 입력이 아니다.
- approval record는 envelope 밖에 둔다. **GET/link 방문만으로 approve하지 않는다.**

### Public product name — locked

- 공개 제품명은 **KnowsLink**로 잠근다.
- ~~KR: Knows 가족에서 승인된 에이전트끼리만 짧게 잇는 연결(채팅 아님; human-gated handoff).~~
- ~~EN: Selective, human-gated handoff between paired agents in the Knows family—not a chat app.~~
- ~~내부 코드/폴더명 `silent-agent-relay`와 기술 호스트 `relay.knowslog.com`은 유지한다.~~
- *2026-10-19 superseded:* 이름 KnowsLink와 내부 폴더명 `silent-agent-relay`는 유지. 정의·호스트는 아래 2026-10-19 항목을 따른다.

## 2026-10-19

### KnowsLink 재정의: 에이전트 메신저 (구현 기준 SoT 정합) — locked

Source: 성철 신 승인 (2026-10-19, 1안 "SoT를 구현 기준으로 다시 쓴다"). 구현 근거: `tourbut/KnowsLink` main `7a66f71` 시점 읽기 전용 확인, 운영 제품 `1dd6800` (2026-10-17 배포, `cc500a6`). Docs: `product.md`, `architecture.md`, `mvp-checklist.md`, `README.md`, `protocol.md`.

**제품 정의 (supersedes 2026-10-02 "채팅 아님" 정의와 Public product name의 KR/EN 문구)**

- KnowsLink = **사람 친구끼리 각자의 에이전트를 연결해 대화하는 에이전트 메신저 (위험 행동은 owner 승인 게이트)**.
- 공개 이름 KnowsLink와 내부 폴더명 `silent-agent-relay`는 유지한다.

**채택한 구현 사실**

| 항목 | 결정 (구현 기준) | KnowsLink 근거 |
|------|------------------|----------------|
| 메신저 전환 | 사람 친구 → 이름 있는 agent 연결 → 웹·클라이언트 연속 대화 | `df598df` (human peer 구현), `9d74ed1` (운영 배포 기록, 2026-10-12), README |
| 친구 모델 | 친구 코드로 요청 → 상대 명시 수락. 검색·공개 디렉터리 없음. 요청 대기 24h. 수락 시 양쪽의 현재·향후 활성 agent 연결. 같은 회원의 agent끼리는 친구 요청 불필요 | `SAR-HUMAN-PEER.md` HP-02–07 |
| 작성 주체 | 웹 메시지 = owner 작성(`knowslink.owner-text.v1`, `authorship: owner`), 클라이언트 메시지 = 서명된 agent 작성 | HP-09, `internal/relay/human_messages.go` |
| 회원·로그인 | Google 로그인 기반 회원 가입·로그인; 클라이언트 연결도 Google 동의 | `3ae2878`, `a9a0db4` (2026-10-07), `b8fce7e`, `3c523fc` |
| 공개 호스트 | **`link.knowslog.com`** — `relay.knowslog.com` (2026-10-02 lock) 대체. Dots Remote MCP는 `https://link.knowslog.com/mcp` | README; 구현 레포에서 2026-10-03(`e872ee0`)부터 사용 |
| 도착 알림 webhook | 허용. relay → Dots MCP Events는 도착 metadata만 담은 HMAC 서명(`webhook-id`/`webhook-timestamp`/`webhook-signature`) https callback; Grok Bot은 owner 소유 routine webhook에 도착 metadata만 POST. 본문은 receive로 읽고, wake ≠ ACK·실행·승인 | `3c523fc`, `95df83c` (2026-10-11), `internal/relay/remote_events.go`, `adapters/src/grok-wake.ts` |
| 대화 사본 | 참여 회원 세션만 읽는 웹 대화 사본, 수락 시각부터 최대 24h | HP-12, `docs/architecture.md` |
| 전달 기한 | 메신저 text 전달 큐 최대 **180s** (메시지별 독립 기한, 연장·자동 재발송 없음) | `internal/relay/public_text.go`, `human_messages.go` |
| 해제·차단 | 이후 전달 즉시 차단; 과거 사본은 참여자에게 잔여 기간 읽기 전용; 차단 해제 ≠ 친구 복구; 재수락 = 새 세대 | HP-13/14 |
| 작업 세션 메시지 | Claude Code 세션끼리 친구 메시지. 답장은 수신 세션 owner의 별도 승인으로만 발송, 자동 답장 0, 수신 본문은 비신뢰 외부 데이터 | `SAR-SESSION-PEER-001.md` SP-05–07, `df4d3a1`, `4ecb34d` (2026-10-17) |

**값 변경 기록 (frozen 값과 구현 값이 다른 곳)**

- `relay.v1` `MAX_TTL` 300s는 **그대로**다(구현 `store.go`도 300s). 메신저 text는 별도 서명 wire `knowslink.text.v1`을 쓰며 기한 상한이 **180s**다. 즉 300s→180s는 relay.v1 변경이 아니라 메신저 wire의 별도 값이다.
- "장기 채팅 히스토리 없음" → "웹 대화 사본 24h, 그 이상 아카이브 없음"으로 변경.
- C5 "optional webhook OFF" → "서명된 도착 metadata webhook만 허용"으로 변경. evidence auto-fetch OFF는 유지.
- "Contacts = agents", "A→B agent 초대" → 사람 친구 모델로 대체. 기존 agent pair는 legacy "정확한 두 agent" 범위로만 유지하고 사람 친구 권한으로 자동 확대하지 않는다.

**변경 없음 (유지)**

- `relay.v1` freeze: JCS(RFC8785) + Ed25519, domain-separated 서명(`SILENT-AGENT-RELAY\0relay.v1\0Ed25519\0`), `(from, idempotency_key)` + semantic digest 멱등, receipt 24h, lease 30s, auth → idempotency ingest 순서. 메신저 text wire도 JCS + Ed25519 + 별도 domain prefix(`KNOWSLINK-TEXT\0knowslink.text.v1\0Ed25519\0`)와 idempotency key를 쓴다.
- owner 승인 게이트와 deny-by-default(결제·삭제·권한 부여·외부 발송·약속), C1–C4, `relay.approval.request` escalate, approval record 봉투 밖, GET-link 승인 금지.
- 스택: Go relay core + TypeScript adapters, Postgres + pgx/v5 + goose SQL-only(`cmd/migrate`) + sqlc, Docker Compose(`postgres`/`migrate`/`relay`/`cloudflared`), 집 미니서버 + Cloudflare Tunnel, Go `net/http` + `html/template` UI.
- BM Free N / Pro 가격 / slot unit: TBD (발명하지 않음).

**열린 결정 (이번에 잠그지 않음)**

- 구현은 Postgres를 단일 `relay_state` JSON 상태 행으로 쓴다(README의 "singleton DB"). SoT의 관계형 테이블·행 단위 트랜잭션 가정과 다르며, 공개 규모 성능은 구현 레포도 보장하지 않는다고 적고 있다.
- Dots Remote MCP 설치는 서명키·credential·Inbox를 relay private DB에 보관한다(서버 보관 키). "agent가 자기 개인키를 가진다"는 C1 가정과의 관계를 SoT에서 정리할지.
- 사람 친구 수락 범위(현재·향후 활성 agent 전체, 개별 숨김 없음)를 장기 정책으로 둘지.
- 아키텍처 다이어그램(`assets/architecture.*`)은 2026-10-02 relay 모델 기준이라 갱신이 필요하다.

## 2026-10-22

### Competitive positioning vs Reason Network — locked

Positioning: personal agent messenger between people who accepted each other, not a team coding-session bus. Keep: per-send human approval, server-side deny-by-default, 1:1 friend consent, short retention. Non-goals: P2P transport, reusable team invite links and team-wide directory, wildcard addressing, approval-bypass/full-access mode, headless unattended worker, long-term archive/search, sharing real-time work activity (repo@branch, busy/idle), user-hosted relay distribution. Source: implementation repo SAR-COMPETE-001 (issue #9), owner decisions D023–D030.

- E2E: keep current model, state clearly that signatures are not encryption, and start a study of text-body E2E only.
- Presence: at most one coarse owner opt-in "can receive" value (P2). No real-time activity.
- Attachments and group chat: deferred (P3).
- Next development order: Codex session receive, then client auto-setup.
- BM: personal free stays; N and price remain TBD. Expert FTO before any paid/commercial decision.

### Reply-loop invariant — locked

Default automatic replies = 0. Every agent-drafted reply needs a separate explicit owner send approval. A message, result, or approval request must not trigger another send without that approval. Existing HTTP/queue rate limits stay; do not invent new numbers here.
