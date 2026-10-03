# Silent Agent Relay — 결정 로그

## 2026-10-02

### Grill points 1–7 — locked

1. Agent signup + pubkey; owner revoke/rotate  
2. Invite A→B; B-human must accept (auto only for already-paired reinvite)  
3. Short TTL queue only; expired → `failed: expired`; no long-term chat history  
4. Deny-by-default silent ops; human required for payment, delete, grants, outbound sends, commitments  
5. Core = HTTPS relay; plugins pull-default, webhook if available  
6. Idempotency: same key+payload → cache hit; conflict → require new key  
7. Contacts = agents; A picks target agent (humans are owners only)

### Business model — locked direction

- Tailscale-like freemium: personal free + friend/pair slot limits; Personal Pro unlocks more slots.
- Frame as “friends up to N”. Exact **N** and Pro price: **TBD**.
- Message-volume / Slack-seat / marketplace framing omitted for now.
- Slot-unit refinement deferred; **do not invent N/price**.

### Protocol / product drops

- **`deliver: both` dropped.** Only `agent` | `human` (default `agent` at app layer; **wire must set explicitly**).
- Human-side NL is a **render** of the structured envelope (`render.hint`), not the wire format.
- Invent thin `relay.v1` (ASMTP-shaped mailbox + A2A Task–inspired); not a full ASMTP/A2A clone.

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
- **MVP OFF:** evidence auto-fetch/preview; optional webhook.
- Resource/rate limits TBD (never unlimited). Pending invites do not consume active slots.

**Defer:** evidence remote fetch; optional webhook; `schedule.commit` executable; amend.  
**Reject:** wire/schema bypasses; auth-meaning bypasses; ops-invariant breaks; “HTTPS ref is safe”; invent Free N / Pro price / rate / extra size numbers; BM slot-unit change.

### Tech stack / deploy — locked

Source: user 성철 신 (2026-10-02). Docs: `architecture.md`, `product.md`, `mvp-checklist.md`, `README.md`.

- **MVP hosting:** home **mini-server** (self-host process). **No paid cloud app server** for MVP.
- **Public HTTPS:** **Cloudflare Tunnel** as **edge ingress only** (TLS/public entry → local relay). Public host is **`relay.knowslog.com`** on `knowslog.com`; Tunnel ingress routes it to the home mini-server. Tunnel is not the application runtime.
- **Reject as MVP hosting assumption:** Cloudflare Workers / Durable Objects (and similar serverless relay cores). May revisit later; not the current plan.
- **Database:** **Postgres** — locked (keys/pairs/intent registry/queue/lease/receipts/approval records; atomic txn requirements). **Not SQLite.**
- **Relay core:** **Go** — the HTTPS relay process/runtime.
- **Adapters/plugins:** **TypeScript** — including MCP pull clients.
- **First MVP adapters (locked, priority order):** **Grok Bot → Claude Code → Codex → Dots** (TypeScript, pull-default).
- These are adapter targets, not claims of official inbound APIs. No integration details are implied here.
- Unchanged: HTTPS relay core; personal public-net; no hospital/closed-net; no vendor core mods; pull-default adapters; webhook/evidence MVP OFF; `relay.v1` freeze; BM Free N / Pro price TBD (do not invent).

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
- KR: Knows 가족에서 승인된 에이전트끼리만 짧게 잇는 연결(채팅 아님; human-gated handoff).
- EN: Selective, human-gated handoff between paired agents in the Knows family—not a chat app.
- 내부 코드/폴더명 `silent-agent-relay`와 기술 호스트 `relay.knowslog.com`은 유지한다.
