# KnowsLink — 제품 결정

**공개 제품명 (잠금):** KnowsLink

- **KR:** Knows 가족에서 승인된 에이전트끼리만 짧게 잇는 연결(채팅 아님; human-gated handoff).
- **EN:** Selective, human-gated handoff between paired agents in the Knows family—not a chat app.
- 내부 코드/폴더명은 `silent-agent-relay`를 유지하며, 기술 호스트 `relay.knowslog.com`도 유지한다.

**잠긴 날짜:** 2026-10-02

개인 퍼블릭넷용 선택적 인간개입 에이전트 딜리버리. 병원·폐쇄망은 범위 밖.

## 잠긴 결정 (1–7)

1. **에이전트 가입 + pubkey**  
   에이전트가 등록하고 공개키를 올린다. owner가 revoke / rotate 할 수 있다.

2. **초대: A→B, B-human 수락 필수**  
   첫 페어링은 B 쪽 인간이 수락해야 한다. 이미 페어링된 재초대만 auto 허용.  
   Pending invite는 상대 **active pair slot을 소진하지 않는다** (C5). Unpair 후 재수락 = 새 관계 세대 (C1).

3. **짧은 TTL 큐만**  
   mid-message를 장기 저장하지 않는다. 만료 시 A에게 `failed: expired`. 장기 채팅 히스토리 없음.

4. **Deny-by-default silent ops**  
   침묵 처리는 기본 거부. 다음을 수행하려면 인간 게이트 필수:  
   - payment (결제)  
   - delete (삭제)  
   - grants (권한 부여)  
   - outbound sends (외부 발송)  
   - commitments (약속·확약)

5. **코어 = HTTPS 릴레이**  
   플러그인은 **pull-default**. **MVP: optional webhook OFF** (C5).  
   `delivered`는 lease가 아니라 **durable inbox persist → ACK 성공 → 공유 실행 claim** 이후.

6. **멱등성**  
   동일 idempotency key + payload → cache hit (**receipt** 반환, 원 봉투 전체 아님).  
   key는 같은데 payload가 다르면 conflict → **새 key 필요**.

7. **Contacts = agents**  
   연락처 단위는 에이전트다. A는 대상 에이전트를 고른다. 사람은 owner로만 존재한다.

## 등록·페어링 모델

### Agent-centric registration

- 가입 주체는 **에이전트**다.
- 사람은 그 에이전트의 **owner**로만 바인딩된다.
- owner는 pubkey revoke/rotate, 페어링 수락/거절, 인간 게이트 응답을 담당한다.

### Pairing (친구 추가형)

1. A-agent(또는 A-owner)가 B-agent에게 invite.
2. B-human이 수락 (재초대·**현재 active pair**만 auto 가능; 과거 paired 사실로 철회 복구 금지).
3. Accept는 중복 관계·slot 점유를 원자 처리. Pending은 active slot을 소모하지 않음.
4. 상호 수락 후 양쪽 **allowlist contacts**에 상대 에이전트가 오른다.
5. 이후 릴레이는 allowlist에 있는 에이전트에게만 허용. Agent cred ≠ owner accept/approve.

## Decision pipeline (수신 측)

### Relay ingest (send) — locked order (C3)

**이 순서만** 구현 기준. 옛 `validate(exp) → idempotency → auth` 서술은 **폐기**.

1. **structure / signature / auth / routing** — strict JSON·스키마·서명·발신 허용·allowlist (**auth 먼저**)
2. **semantic digest**
3. **atomic idempotency** — key+digest; conflict → 새 key; TTL/id 실패 시 reservation rollback
4. **exp / MAX_TTL / id collision**
5. **queue + receipt** same txn — cache hit → **receipt only** (재enqueue·TTL 연장·새 실행 허가 없음)

### Agent-side judgment (after shared exec claim)

수신 에이전트가 durable persist → ACK 성공 → **공유 실행 claim** 이후 아래를 판단한다.

1. **permission** — 요청 intent가 silent 허용 범위인지 (현재 key/pair/owner 재검사)
2. **judgment** — 도메인 판단이 에이전트만으로 가능한지
3. **risk** — deny-by-default 위험 액션(결제·삭제·grant·outbound·commitment) 여부
4. **execute or escalate**
   - 가능하면 자동 처리 → `relay.result` `status: done` (B-human **미통지**) — C4 출력 allowlist 적용
   - 인간 권한·판단 필요 → `relay.approval.request` (`deliver: human`) 으로 escalate

## Silent resolve vs escalate

| 경로 | 조건 | B-human | 결과 |
|------|------|---------|------|
| Silent | permission·judgment·risk 모두 에이전트 단독 가능 | 통지 없음 | `relay.result` `status: done` |
| Escalate | 위 중 하나라도 인간 필요 | 인간에게 렌더된 NL로 전달 | `relay.approval.request` → approve/deny (amend MVP 밖) → `relay.result` |

- Escalate는 **새** 메시지: `intent=relay.approval.request`, `reply_to`=원요청 id, typed `request_digest`. parent-link top-level 필드 없음.
- 승인 UI = 검증된 **M typed body** (not `render.hint`). 승인 기록은 **봉투 밖**. No GET-link approve. Approve ≠ enable commit stub (C2).
- agent sig / `render.hint` ≠ owner 승인 증명.
- `deliver` 앱 기본값은 `agent`이나 **와이어에는 명시**. **`both`는 사용하지 않는다** (2026-10-02 폐기).
- 스키마·ops·security SoT: [`protocol.md`](./protocol.md) frozen 2026-10-02 (+ round-2 ops/seed + **security C1–C5**).

## Intent seed (요약)

| Intent | MVP 역할 |
|--------|----------|
| `schedule.query` | **기본 비공개** (C4); disclosure policy 없으면 deny; 범위 수치 TBD |
| `schedule.commit` | stub only (executable 아님) |
| `relay.result` | 완료 보고 (`reply_to` 필수) |
| `relay.approval.request` | 인간 게이트 escalate (도구 미실행) |

## Security notes (C1–C5 locked)

- **C1:** `(from,kid)` key lookup only; PoP on register/rotate; current-auth re-check at enqueue/lease/ACK/approval/exec/result; authorization epoch CAS.
- **C2:** approval bound to verified M instance; record outside envelope.
- **C3:** ingest auth→idempotency (above); shared exec claim; no uncoordinated multi-adapter.
- **C4:** calendar disclosure default deny; `relay.result` endpoint binding + output allowlist; log minimization.
- **C5:** MVP evidence auto-fetch OFF; optional webhook OFF; resource limits TBD (not unlimited); pending invites ≠ active slots.
- **Reject/Defer:** see `protocol.md` / `decisions.md`. Do not invent Free N / Pro price / rate numbers.

## 범위 메모

- 빌드 = 릴레이 + 플러그인/어댑터. **벤더 코어를 수정하지 않는다.**
- **MVP 배포:** 집 **미니서버**에서 **Docker Compose**로 `postgres`, one-shot `migrate` (`cmd/migrate`), Go `relay`, `cloudflared` (Tunnel ingress only)를 실행한다. Postgres는 퍼블릭에 노출하지 않는다. **Cloudflare Tunnel**(퍼블릭 HTTPS 인그레스)의 공개 호스트는 **`relay.knowslog.com`** (`knowslog.com`)이다. 포트·이미지/다이제스트·환경변수는 구현 레포로 미루며, Workers/DO를 MVP 호스팅으로 쓰지 않는다.
- **DB:** **Postgres** (SQLite 아님). **Relay core는 Go**, **어댑터/플러그인은 TypeScript**.
- 공식 cross-agent inbound API를 주장하지 않는다. 어댑터가 제품별 외부 인터페이스를 감싼다.
- **첫 MVP 어댑터 대상(잠금, 우선순위):** **Grok Bot → Claude Code → Codex → Dots**. TypeScript pull-default로 시작한다. 어댑터 대상 이름만 잠그며 통합 세부사항은 정하지 않는다.
- **Go 첫 구현 범위(잠금):** ingest/queue-only가 아니라 human-gate까지다. validate/auth/ingest, short queue, lease/ACK/exec claim, receipts, `relay.approval.request`, envelope 밖 approve/deny 경로를 포함한다.
- **MVP human-gate UI (잠금):** Go relay core의 **`net/http` + `html/template`**로 최소 approve/deny UI를 relay가 직접 서빙한다. 별도 TypeScript frontend는 사용하지 않는다. UI는 C2에 따라 검증된 M의 typed body + 적용 정책으로 구성하며, `render.hint`를 권위 있는 입력으로 사용하지 않고 GET/link 방문만으로 승인하지 않는다.
- 미니서버 port·path·process manager·Cloudflare Tunnel 설정은 향후 구현 레포의 env/config로 미룬다. **Postgres tooling (잠금):** Go driver는 `jackc/pgx/v5` + `pgxpool`, migration은 `pressly/goose/v3`의 SQL-only 방식으로 별도 `cmd/migrate`에서 실행하며 API 기동 시 auto-up하지 않는다. query codegen은 `sqlc` (`sql_package: pgx/v5`)를 사용한다. MVP에서는 Atlas auto-diff/apply, ORM AutoMigrate, Go-code migrations를 primary path로 채택하지 않는다.
- BM Free N / Pro 가격·slot-unit·rate 숫자는 **TBD** (발명하지 않음).
- C1–C5 잠금 전 실데이터 silent 경로를 열지 않음.
