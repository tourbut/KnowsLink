# KnowsLink (Silent Agent Relay) — MVP 체크리스트

개인 퍼블릭넷 최소 구현. 병원·폐쇄망·벤더 코어 수정·공식 cross-agent API 가정 없음.

## Must (MVP)

- [x] **메신저 재정의 (2026-10-19)** — 사람 친구끼리 각자의 에이전트를 연결해 대화하는 에이전트 메신저 (위험 행동은 owner 승인 게이트). 구현 `tourbut/KnowsLink` 운영 `1dd6800`
- [x] **Google 로그인 회원** — 가입·로그인·클라이언트 연결 동의 (구현됨)
- [x] **사람 친구** — 친구 코드 요청 → 상대 수락; 요청 24h; 해제·차단 → 신규 전달 차단, 과거 사본 읽기 전용 (구현됨, HP-02–14)
- [x] **이름 있는 내 에이전트 + 대화** — agent 두 개 단위 대화; 웹 = owner 작성, 클라이언트 = 서명된 agent; text ≤ 4096 bytes; 전달 ≤ 180s; 웹 사본 ≤ 24h (구현됨)
- [x] **발송마다 사용자 승인 · 자동 답장 0** — 수신·ACK·친구 수락 ≠ 답장 승인 (구현됨)
- [x] **Claude Code 작업 세션 친구 메시지** — owner 승인 답장만, 수신 본문 비신뢰 (구현됨, SAR-SESSION-PEER-001/002; 독립 코드 리뷰 미실시 기록 있음)
- [ ] **Agent register + pubkey (C1)** — 가입, owner 바인딩, revoke/rotate + **PoP**; key lookup `(from,kid)` only; agent cred ≠ owner approve
- [ ] **Invite accept flow (C1/C5)** — *(2026-10-19: legacy agent pair; 신규는 아래 사람 친구)* A→B 초대, B-human 수락; pending ≠ active slots; unpair 재수락 = 새 세대; auto = current active pair only
- [ ] **Short queue + exp** — TTL 큐 (`relay.v1` `MAX_TTL` 300s; 메신저 text ≤ 180s), 만료 시 `failed: expired`, ~~장기 히스토리 없음~~ 웹 대화 사본 24h 초과 보관 없음
- [ ] **Deny-by-default silent rules** — payment / delete / grants / outbound sends / commitments → 인간 게이트
- [ ] **HTTPS relay core** — structure/sig/**auth**/routing → digest → idempotency → exp/TTL/id → queue+receipt; human gate; lease 30s / 3 attempts
- [ ] **Go first implementation scope** — 위 relay core를 ingest/queue-only로 제한하지 않고 human-gate까지 구현: `relay.approval.request` + envelope 밖 approve/deny 경로 (frozen protocol/C1–C5 준수)
- [ ] **ACK + shared exec claim (C3)** — lease (token=recipient/msg/gen) → durable persist → ACK 성공 → **exec claim** → process; no uncoordinated multi-adapter
- [ ] **Receipt store** — 24h retention; cache hit returns receipt metadata (`id, from/to, intent, digest, exp, accepted_at, transport state`); no re-enqueue / no TTL extend on replay
- [ ] **First adapter targets (priority order)** — **Grok Bot → Claude Code → Codex → Dots**; TypeScript, pull-default, no vendor core modifications or assumed official inbound APIs; one pull plugin stub for MVP; ~~MVP webhook OFF~~ 도착 metadata webhook만 허용 (2026-10-19)
- [ ] **Idempotency** — `(from, key)` + semantic digest; conflict → 새 key; retention 24h
- [ ] ~~**Contacts = agents**~~ — *(2026-10-19 superseded)* 친구 = 사람, 대화 = agent 두 개; 친구/자기 연결 범위 안에서만 릴레이
- [ ] **Paths:** silent `deliver:agent` → `relay.result` done / escalate = `relay.approval.request` (`deliver:human`)
- [ ] **Seed intents** — `schedule.query` (**default private**, policy TBD), `schedule.commit` (stub), `relay.result` (endpoint binding+allowlist), `relay.approval.request`
- [ ] **Approval binding (C2)** — UI from verified M typed body (not hint); record outside envelope; no GET-link approve; approve ≠ enable commit stub; amend out of MVP
- [ ] **MVP human-gate UI** — Go relay serves a minimal approve/deny UI via `net/http` + `html/template`; no separate TypeScript frontend; verified M typed body + policy only (C2)
- [ ] **Current-auth re-check (C1)** — enqueue/lease/ACK/approval/exec/result; authorization epoch CAS
- [ ] **Result allowlist (C4)** — `R.reply_to==M.id`, `R.from==M.to`, `R.to==M.from` + output allowlist; log minimization
- [ ] **Boundaries (C5)** — strict JSON same parse object; evidence auto-fetch OFF; ~~webhook OFF~~ 서명된 도착 webhook만; resource limits TBD (not unlimited)
- [ ] **No forged B-signed `relay.result`** on transport failure
- [x] **`relay.v1` schema freeze** — see `protocol.md` (2026-10-02)
- [x] **Round-2 ops/seed freeze** — see `protocol.md` / `decisions.md` (2026-10-02)
- [x] **Security MUST C1–C5 freeze** — see `protocol.md` / `decisions.md` (2026-10-02)
- [x] **MVP deploy stack** — Docker Compose on the home mini-server: `postgres`, one-shot `migrate` (`cmd/migrate`), Go `relay`, and `cloudflared` (Tunnel ingress only); Postgres is not publicly exposed. Cloudflare Tunnel serves **`link.knowslog.com`** (`knowslog.com`; 2026-10-19 replaces `relay.knowslog.com`); **Go relay core + TypeScript adapters/plugins**; **not** Workers/DO as MVP host (2026-10-02). Ports, image/digest, and env details remain deferred to the implementation repo.

## Open / TBD

- [ ] Free pair slot count **N** (비즈니스) — 미잠금 (**발명 금지**)
- [ ] Personal Pro 가격 — 미잠금 (**발명 금지**)
- [ ] BM slot-unit refinement — deferred (변경하지 않음)
- [ ] 추가 플러그인(어댑터) 구현 매트릭스 — 첫 대상은 **Grok Bot, Claude Code, Codex, Dots**; MVP는 stub 1개; 공유 claim 없이 다중 활성 금지
- [ ] `schedule.query` silent **disclosure policy** 세부 (반환 필드·window·granularity·누적 한도) — default = deny; 수치 TBD
- [ ] Resource / rate / size / concurrency 상한 수치 — TBD (무제한 배포 금지)
- [x] **Postgres tooling** — `jackc/pgx/v5` + `pgxpool`; `pressly/goose/v3` SQL-only migrations via separate `cmd/migrate` (API auto-up 금지); `sqlc` with `sql_package: pgx/v5`. Atlas auto-diff/apply, ORM AutoMigrate, Go-code migrations는 MVP primary path에서 제외 (2026-10-02)
- [ ] `schedule.commit` 실행 가능화 (effect / auth / commit idempotency) — stub; owner approve만으로 대체 불가
- [ ] Evidence remote fetch sandbox — Deferred (MVP OFF)
- [ ] ~~Optional webhook hardening — Deferred (MVP OFF)~~ → 2026-10-19: 도착 알림 webhook 구현됨 (HMAC 서명, https callback). 장시간·만료 갱신 안정성은 구현 레포 기준 미검증
- [ ] Postgres 단일 `relay_state` JSON 상태 행 vs SoT 관계형 가정 — 열린 결정
- [ ] Dots Remote MCP의 서버 보관 agent 키와 C1 관계 정리 — 열린 결정
- [ ] 아키텍처 다이어그램 메신저 반영 갱신
- [ ] Codex 세션 수신 지원 판정 — SP-12 수준별 근거 (2026-10-22, 다음 개발 1순위)
- [ ] text 본문 종단 간 암호화 조사 — 서버 게이트·웹 대화와의 양립 조건
- [ ] 거친 '수신 가능' 상태의 opt-in 공개 범위 — 실시간 활동 공유 금지

## Explicit non-goals (MVP)

- [x] ~~Hospital / closed-net~~ — out of scope
- [x] ~~Vendor agent core patches~~ — plugins only
- [x] ~~`deliver: both`~~ — dropped
- [x] ~~Latent KV / token-id LLM handoff~~ — research only
- [x] ~~Auth tokens in message body~~ — forbidden
- [x] ~~Message-volume / Slack-seat / marketplace framing~~ — BM에서 일단 생략
- [x] ~~parent-link top-level field~~ — use `reply_to` + body digest
- [x] ~~amend~~ — out of MVP
- [x] ~~Invented rate limits / BM slot-unit / Free N / Pro price numbers~~ — rejected
- [x] ~~MVP evidence auto-fetch~~ — OFF (C5); optional webhook은 2026-10-19 도착 알림 한정 허용
- [x] ~~장기(24h 초과) 대화 아카이브~~ — 비목표 (2026-10-19)
- [x] ~~Pre-auth idempotency probe / full-envelope cache hit / lease=delivered~~ — rejected (C3)
- [x] ~~GET-link approve / hint-as-approval / approve⇒commit stub~~ — rejected (C2)
- [x] ~~P2P transport / user-hosted relay~~ — central relay enforces gate (2026-10-22)
- [x] ~~Team invite links / team-wide directory / wildcard address~~ — 1:1 friend consent (2026-10-22)
- [x] ~~Approval-bypass full-access mode / headless unattended worker~~ — per-send approval (2026-10-22)
- [x] ~~Real-time work activity sharing (repo@branch, busy/idle)~~ — privacy (2026-10-22)

