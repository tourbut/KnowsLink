# relay.v1 — 프로토콜 (frozen 2026-10-02; ops/seed round-2 2026-10-02; security C1–C5 2026-10-02)

구조화된 JSON 봉투. **자연어가 와이어 포맷이 아니다.**  
인간에게 보이는 NL은 이 봉투의 **렌더**다 (`render.hint`).

영감: ASMTP형 mailbox + A2A Task 의미론. MVP는 얇은 `relay.v1`만.

> **Freeze:** intent + typed body · signed · short-TTL · idempotent mailbox.  
> agent↔agent 채팅 프로토콜이 **아님**.

## 불변식

- `relay state ≠ agent processing status`
- `deliver` = `agent` | `human` only (`both` = invalid → `422 invalid_deliver`)
- 서명된 요청 봉투는 **immutable** (릴레이 mutate 금지)
- 같은 `idempotency_key` + 다른 semantic digest → `409` + **새 key 필수**
- NL wire / canonical `body.text` 금지
- **parent-link 전용 top-level 필드 없음.** 상관·escalation은 `reply_to` + typed body digest만 사용

## 봉투 필드 (frozen)

| 필드 | 타입 | 필수 | 규칙 | 서명 커버 |
|------|------|:----:|------|:--------:|
| `v` | string | yes | 정확히 `"relay.v1"` | yes |
| `id` | string | yes | UUIDv7 메시지 인스턴스 ID | yes |
| `from` | AgentID | yes | 등록 에이전트; lowercase ASCII 정규화 | yes |
| `to` | AgentID | yes | 등록 목적지 에이전트만. 사람/이메일 ID 금지 | yes |
| `intent` | string | yes | 등록 intent (예: `schedule.query`); max 64 bytes | yes |
| `body` | object | yes | intent별 closed typed object; max 16 KiB JCS | yes |
| `deliver` | string | yes (wire) | `"agent"` \| `"human"` only. 앱 기본값은 agent이나 **와이어에는 명시** | yes |
| `exp` | RFC3339 UTC | yes | 초 정밀도; 미래; max = 수신시각 + `MAX_TTL` | yes |
| `sig` | object | yes | `{alg:"Ed25519", kid, value}`; value는 base64url Ed25519 | value만 제외 |
| `idempotency_key` | string | yes | caller 생성, 16–128 ASCII | yes |
| `priority` | string | no | `low` \| `normal` \| `high`; default `normal` | yes |
| `evidence` | array | no | default `[]`; max 8 immutable refs | yes |
| `reply_to` | UUIDv7 | conditional | 이전 메시지 상관. `relay.result` / `relay.approval.request`에서는 **필수** | yes |
| `trace` | object | no | W3C `traceparent`, optional `tracestate`; 관측만 | yes |
| `render` | object | no | `{hint, lang?}`; **표시용만**, 권한/액션 의미 금지 | yes |
| `ext` | object | no | reverse-domain 확장 맵만 | yes |

### Nested rules

- Discriminator는 `intent`다. `body`에 generic `type`/`op` discriminator를 두지 않는다.
- Intent 스키마는 closed. 목적지 intent registry가 모르는 body 필드는 거부.
- `evidence[]` 항목: `{ref, media_type?, sha256?}`. `ref`는 `https:` 또는 `urn:`. max 8, ref 중복 금지, bytewise 정렬. **MVP: Relay·agent·UI의 evidence 자동 dereference/preview 비활성** (C5). 허용된 ref 수신 ≠ 네트워크 접근. `sha256`가 있어도 SSRF를 막지 못한다.
- `render.hint`: NFC, max 1024 UTF-8 bytes. 인가·액션 의미를 실으면 안 된다. **owner 승인 증명으로 쓰지 않는다.**
- 알 수 없는 top-level / nested protocol 필드 → `422 invalid_schema`. 확장은 `ext` 아래만.
- `ext`는 인가·deliver 의미를 바꾸지 못한다. 실행에 필수면 새 intent 또는 프로토콜 버전.
- **새 top-level parent-link 필드를 추가하지 않는다.**

### 폐기 (요청 봉투에서)

- `kind`, `body.op`, `body.args`
- 요청 봉투 안의 `status`
- canonical `body.text` / NL wire
- `deliver: "both"`
- 다이어그램 레거시 `v: "1"`
- parent-link 전용 top-level 필드

## Intent registry seed (round-2 locked)

### 이름 불변식

- Intent **문자열은 immutable**: 같은 이름을 새 의미로 remap하지 않는다.
- Registry revision은 **ADD only** (기존 intent 의미 변경 금지). breaking은 **새 intent 이름** 또는 프로토콜 버전.
- Receiver는 registry **subset만 구현**한다. 로컬에서 의미를 재정의하지 않는다.

### Seed intents

| Intent | `deliver` | `reply_to` | Body (closed) | 메모 |
|--------|-----------|------------|---------------|------|
| `schedule.query` | `agent` (전형) | optional | `{window:{start,end}, granularity_min}` | **기본 비공개** (C4). pair별 disclosure policy 없으면 silent 거부. 반환 필드·window·granularity·누적 한도는 **TBD** |
| `schedule.commit` | `human` (전형) | optional | `{slot:{start,end}, timezone, commitment}` | **stub**. 실제 effect / auth / idempotency 없이 실행 불가. MVP에서 executable 아님 |
| `relay.result` | `agent` | **required** | `{status: done\|denied\|failed, result?, error?}` | 수신 에이전트가 원요청에 대해 **새** 서명 메시지로 완료 보고 |
| `relay.approval.request` | **`human` only** | **required** | `{reason: permission_required\|judgment_required, request_digest}` | escalate 전용. **도구를 직접 실행하지 않음**. `request_digest` = 원요청 semantic digest |

`relay.approval.request`는 인간 게이트용 요청일 뿐, 승인 실행·도구 호출을 수행하지 않는다.

## 운영 상수 (frozen 2026-10-02)

| 상수 | 값 | 메모 |
|------|-----|------|
| `MAX_TTL` | **300s** | 일반 요청 권장 60–120s; 초과 → `ttl_too_long` |
| Idempotency / receipt retention | **24h** | 메시지 TTL보다 김. cache hit는 **receipt metadata** 반환 (원 봉투 전체 아님) |
| Pull lease window | **30s** | `lease_until = min(exp, lease_granted_at + 30s)` |
| Delivery attempts | **3** | 성공적 lease grant마다 +1. 3번째 lease도 자기 window 동안 usable. 소진 시 `failed:max_attempts` |
| Intent registry | **central, ADD-only revisions** | 이름 remap 금지; breaking → 새 intent/버전 |
| Key revoke metadata retention | **≥ idempotency retention (24h)** | revoke 이력 보존 |

Rate limit / 추가 자원·봉투 크기 상한 숫자(이미 명시한 body JCS·hint·evidence ref 한도 외)는 **TBD**이며 **발명하지 않는다.** 무제한으로 배포하지 않는다 (C5).

## 상태 ownership

### Relay transport (봉투 밖 저장)

```
queued → leased → delivered
queued/leased → failed:expired
leased → queued          (lease timeout, attempts 남음)
leased → failed:max_attempts
```

**`delivered` 정의 (round-2 + C3):** lease grant가 아니다.  
시퀀스: **lease → durable inbox persist → ACK 성공 확인 → 공유 실행 claim → process**.  
transport `delivered`는 ACK 성공 확인 이후다.

- attempts는 **성공적 lease grant**에서만 증가한다.
- `lease_until = min(exp, lease_granted_at + 30s)`.
- Lease token은 **수신 주체 · 메시지 · 현재 lease 세대**에 결속한다. ACK는 해당 lease 소유권·유효성을 원자적으로 확인한다.
- 늦은 ACK로 다른 lease를 완료하거나 attempts를 변경하지 않는다. 3번째(마지막) lease도 그 window 안에서는 usable하다.
- 수신자 전체 **공유 durable inbox + 실행 claim**. 공유 조정 없는 다중 adapter 동시 활성화 금지. `delivered` 조회만으로 실행권을 얻지 않는다.
- 만료된 요청으로 신규 실행을 시작하지 않는다. 외부 도구 exactly-once를 주장하지 않는다.
- 릴레이는 `delivered → done`으로 전이하지 **않는다**. `done`은 에이전트 처리 결과 (`relay.result`).

### Agent completion — `relay.result`

B-agent는 **새 서명 메시지**를 보낸다:

- `intent`: `relay.result`
- `reply_to` = 원요청 `id` (**필수**)
- `body.status`: `done` \| `denied` \| `failed`
- `body.result` / `body.error`: optional structured — **부모 intent별 출력 allowlist**로 제한 (C4)

권위 있는 원요청 **M**에 대해 결과 **R** 바인딩 (C4):

```text
R.reply_to == M.id
R.from     == M.to
R.to       == M.from
```

이후에도 현재 관계·공개 권한을 재확인하고, 일정 제목·참석자·위치·원본 객체·stack trace 등을 자동 포함하지 않는다.  
결과 수신을 새 도구 실행 명령으로 해석하거나 제어 메시지에 자동 재응답하지 않는다. 허용된 결과 반환 ≠ 일반 outbound send 권한.

**릴레이 transport failure를 B-서명 `relay.result`로 위조하지 않는다.**  
만료·max_attempts·라우팅 실패 등은 릴레이 측 transport 상태/에러로만 보고한다.

### Escalate — `relay.approval.request`

권한·판단이 필요하면 **새** 서명 메시지. 원요청 `deliver: "agent"`를 릴레이가 mutate하지 않는다.

Escalation 메시지 **H** 규칙 (원요청 **M**):

| 규칙 | 값 |
|------|-----|
| `intent` | `relay.approval.request` |
| `deliver` | `human` |
| `reply_to` | `M.id` |
| `from` / `to` | 둘 다 `M.to` (B-agent → B-agent 목적지, 인간 게이트로 전달) |
| `body.reason` | `permission_required` \| `judgment_required` |
| `body.request_digest` | `semantic_digest(M)` |
| `H.exp` | `≤ M.exp` |

추가 규칙 (C2):

- `relay.result` / `relay.approval.request`에서 **재귀 escalate 금지**.
- B→B 경로(위 `from=to=M.to`)는 **좁은 예외**: **권위 있는 부모 기록** + B의 실제 수신 확인 + 검증된 부모 A↔B 페어링 + **현재** B-owner accept. revoke 후 pairing restore로 쓰지 않는다. 부모별 pending gate 중복을 억제한다.
- B-human 초대 미수락 → `403 human_invite_required`.
- Human gate UI는 `render.hint`가 아니라 **검증된 원 서명 메시지 M의 typed body + 적용 정책**으로 구성한다. 원문이 없으면 승인 불가이며 hint로 대체하지 않는다 (원문을 24h 보관해서도 안 됨).
- **승인 기록은 봉투 밖**: 최소 owner, `M.id`, `semantic_digest(M)`, `M.from`/`M.to`, pair 세대, 적용 정책, 만료, 결정·소비 상태. 유효기간 ≤ `min(M.exp, H.exp)`. 같은 digest라도 다른 `M.id`에 재사용 금지. 결정·소비는 원자적.
- 승인 endpoint는 owner 인증·CSRF 방어. **GET이나 링크 방문만으로 승인하지 않는다.**
- agent `sig` / `render.hint` / `ext` / 모델 confidence는 owner 승인 증명이 **아니다**. Hint는 비권위 보조 텍스트만 (HTML·외부 preview·기만적 표시 차단).
- **Approve = 해당 요청의 gate 통과일 뿐.** `schedule.commit` stub을 실행 가능하게 만들지 않는다.
- **amend는 MVP 밖** (approve / deny만).

## Ingest 순서 (round-2 + C3 locked)

수신(send) 경로 처리 순서 — **이 순서만** 구현 기준 (`product.md`의 옛 validate→idempotency→auth 서술은 폐기):

1. **Structure / signature / auth / routing** — strict JSON·스키마·서명·발신 허용·allowlist 라우팅 (**auth 먼저**)
2. **Semantic digest** — `I` 계산
3. **Atomic idempotency** — `(from, idempotency_key)` + digest 비교
4. **`exp` / `MAX_TTL` / `id` collision** — TTL·만료·메시지 ID 충돌
5. **Queue + receipt 동일 트랜잭션** — enqueue와 receipt 기록 atomic

TTL·ID 검사 실패 시 임시 idempotency reservation까지 **rollback**.  
인증 전 idempotency/cache probe 금지.

Cache hit (auth 이후, 같은 key+digest):

- **원 봉투 전체가 아니라 RECEIPT**를 반환한다 (receipt-only replay).
- Receipt metadata (24h retention): `id`, `from`, `to`, `intent`, `digest`, `exp`, `accepted_at`, transport state.
- 24h 내 인증된 동일 key+digest 재전송은 만료됐더라도 **기존 receipt만** 반환. 재enqueue·TTL 연장·새 실행 허가로 해석하지 **않는다**.
- 24h 경과 후 원 서명 replay는 `exp`로 차단. 새 서명 요청까지 영구 업무 멱등성을 주장하지 않는다.

## 서명 (frozen)

- RFC 8785 JCS + **Ed25519 only** (v1)
- `E = envelope_without_sig`
- Signing bytes:

```
"SILENT-AGENT-RELAY\0relay.v1\0Ed25519\0" + sig.kid + "\0" + UTF8(JCS(E))
```

- `sig.value` = base64url-no-pad Ed25519
- `sig.value`를 제외한 모든 존재 필드 커버 (`kid`는 prefix에 포함 → key substitution 방지)
- E2E 서명 ≠ TLS hop auth. 릴레이 검증 후 B-agent 재검증
- 키 유효/revoke는 **수신 시각** 기준 (수신 시점 서명 판정 유지)
- 검증 키는 등록된 **`(from, kid)`** 에서만 조회. `kid`를 URL·파일 경로로 해석하거나 기존 `kid`를 다른 키에 재할당하지 않는다 (C1)
- 키 등록·rotation = **owner 인증 + 새 키 PoP**. Agent credential로 owner invite accept·approve·revoke를 행사할 수 없다
- Agent 서명은 **발신 무결성**만 증명한다. owner 인간 승인·amend 증명으로 확장하지 않는다 (승인 기록은 봉투 밖)
- 수신 서명 판정과 별개로 enqueue·lease/ACK·승인 소비·도구 실행·결과 공개에서 **현재** key/pair/owner 권한을 재검사한다. 서버 상태 확정은 authorization epoch에 대한 **원자적 CAS**. unpair 후 재수락 = **새 관계 세대**. 권한 상태를 확인할 수 없으면 차단
- 보장 범위: “철회가 해당 인가 확정 지점보다 먼저 반영됐으면 거부”. 이미 확정된 정보 공개를 소급 취소하지는 않음

## 멱등성 (frozen)

- DB unique scope: `(from, idempotency_key)`
- Semantic fingerprint `I`:

```
v, from, to, intent, body, deliver,
priority ?? "normal",
evidence ?? [],
reply_to ?? null,
ext ?? {}
```

- Digest: `BASE64URL_NO_PAD(SHA256(UTF8(JCS(I))))`
- **포함:** 위 `I` 필드
- **제외:** `id`, `idempotency_key`, `exp`, `sig`, `trace`, `render`
- 같은 key + 같은 digest → **receipt 반환** (재enqueue/재deliver 없음). 저장된 `exp`는 연장하지 않음
- 같은 key + 다른 digest → `409 idempotency_conflict` (overwrite/patch/LWW 금지)
- `to`/`intent`/`body`/`deliver`/priority/evidence/`reply_to`/`ext` 변경 → **새 key**
- 만료 후 재실행 → **새 key**
- `trace`/`render`만 달라도 conflict 아님, 원 저장본 mutate 금지

## 에러 코드 (권장 계약)

`invalid_json`, `invalid_auth`, `invalid_signature`, `sender_not_allowed`, `human_invite_required`, `idempotency_conflict`, `expired`, `invalid_schema`, `invalid_deliver`, `unsupported_intent`, `ttl_too_long`

## 샘플 — silent `schedule.query`

```json
{
  "v": "relay.v1",
  "id": "0199a3f2-4c10-7a11-8b22-334455667788",
  "from": "agt_01K6S4M9V7Q2A5N8C3D1E0F6GH",
  "to": "agt_01K6S4P1Y8R4B6M2N7T9V3W5XZ",
  "intent": "schedule.query",
  "body": {
    "window": {
      "start": "2026-10-02T03:00:00Z",
      "end": "2026-10-02T04:00:00Z"
    },
    "granularity_min": 30
  },
  "deliver": "agent",
  "exp": "2026-10-02T02:06:00Z",
  "idempotency_key": "idem_01K6S5A0JH1Q4X7C9N2V8M3R6T",
  "sig": {
    "alg": "Ed25519",
    "kid": "ed25519:reJ4SZKc0Ji2iITvq8xvJg",
    "value": "l6u1UxcAUS_2WhUwppIiCJgt6W4Wd4O7WWPD205T6qW_jQOfJNaiR6L8I09TXm6y2Ve4fgVSRFbT18b75zsHAQ"
  }
}
```

## 샘플 — escalate `relay.approval.request`

원요청 M (`schedule.query` 등, `deliver: agent`)에 대해 B-agent가 권한·판단이 필요할 때:

```json
{
  "v": "relay.v1",
  "id": "0199a3f2-9e30-7c44-ad55-66778899aabb",
  "from": "agt_01K6S4P1Y8R4B6M2N7T9V3W5XZ",
  "to": "agt_01K6S4P1Y8R4B6M2N7T9V3W5XZ",
  "intent": "relay.approval.request",
  "body": {
    "reason": "permission_required",
    "request_digest": "BASE64URL_NO_PAD_SHA256_OF_SEMANTIC_DIGEST_M"
  },
  "deliver": "human",
  "reply_to": "0199a3f2-4c10-7a11-8b22-334455667788",
  "exp": "2026-10-02T02:06:00Z",
  "render": {
    "hint": "일정 조회 요청을 승인할까요? (표시용; 승인 증명 아님)",
    "lang": "ko-KR"
  },
  "idempotency_key": "idem_01K6S5D8WX3N5P9R2T6V0Y4Z7A",
  "sig": {
    "alg": "Ed25519",
    "kid": "ed25519:bAgentKidExample0001",
    "value": "PLACEHOLDER_BASE64URL_ED25519_SIGNATURE_VALUE_NOT_REAL_000000000000"
  }
}
```

(`request_digest`·`sig.value`는 illustrative placeholder. `from`=`to`=`M.to`, `H.exp`≤`M.exp`, `reply_to`=`M.id`.)

## 샘플 — stub `schedule.commit` (executable 아님)

```json
{
  "v": "relay.v1",
  "id": "0199a3f2-8d20-7b33-9c44-5566778899aa",
  "from": "agt_01K6S4M9V7Q2A5N8C3D1E0F6GH",
  "to": "agt_01K6S4P1Y8R4B6M2N7T9V3W5XZ",
  "intent": "schedule.commit",
  "body": {
    "slot": {
      "start": "2026-10-03T01:00:00Z",
      "end": "2026-10-03T01:30:00Z"
    },
    "timezone": "Asia/Seoul",
    "commitment": "hold"
  },
  "deliver": "human",
  "priority": "normal",
  "exp": "2026-10-02T02:06:30Z",
  "render": {
    "hint": "2026-10-03 10:00~10:30 일정 보류를 승인할까요?",
    "lang": "ko-KR"
  },
  "idempotency_key": "idem_01K6S5C4QZ8B2W7M1N9R3T6V5X",
  "sig": {
    "alg": "Ed25519",
    "kid": "ed25519:reJ4SZKc0Ji2iITvq8xvJg",
    "value": "7n4ypNsjF3aZW1sDgAZU2NBG0yIFTj-GWqg4Jv7uLWDNCjlSNTRv3-rYDNya6ICr_8s0h67kzfmpLOZSkqZnDQ"
  }
}
```

> **Stub only:** effect·auth·commit idempotency 미구현. MVP에서 실행 경로로 쓰지 않는다. escalate는 `relay.approval.request`를 쓴다.

(시그니처 value는 illustrative placeholder.)

## Security MUST — C1–C5 (locked 2026-10-02)

실데이터 silent 경로를 열기 전에 아래를 잠근다. 등록·페어링·유효 서명은 요청 안전성이나 owner 승인을 증명하지 **않는다**.

### C1 Identity + current auth

- 키 lookup = 등록 `(from, kid)` only. PoP on register/rotate. Agent cred ≠ owner approve/accept/revoke.
- API principal: `send`→허용 `from`; `pull`/ACK→해당 `to`+delivery 경로; receipt→해당 기록 접근권. Agent cred가 `deliver:human` 처리·owner 승인 대체 금지.
- 수신 서명 판정 유지 + enqueue/lease/ACK/approval/exec/result에서 **현재** key/pair/owner 재검사.
- Authorization epoch **CAS**. Unpair 후 재수락 = **새 관계 세대**. 권한 상태 불명 → 차단.
- 정책은 adapter/tool gateway에서 집행 (모델 prompt가 아님).

### C2 Approval binding

- UI = 검증된 **M typed body + 정책** (not `render.hint`). 원문 없으면 승인 불가.
- 승인 기록 = 봉투 밖 (owner, `M.id`, digest, endpoints, pair 세대, 정책, 만료, 결정·소비).
- **No GET-link approve.** Self B→B만 권위 있는 부모+실제 수신 확인 시.
- **Approve ≠ enable `schedule.commit` stub.**

### C3 Atomic ingest + shared exec claim

- Order: `structure/sig/auth/routing → digest → atomic idempotency → exp/TTL/id → queue+receipt`.
- Receipt-only replay. Lease token bound to recipient/msg/generation. Durable persist → ACK success → exec claim → process.
- No multi-adapter without shared coordination. Fail-closed on abnormal clock.

### C4 Disclosure + result allowlist

- `schedule.query` **default private**; no policy = deny. Scope numbers **TBD**.
- `relay.result`: `R.reply_to==M.id`, `R.from==M.to`, `R.to==M.from` + output allowlist.
- Data minimization on durable inbox · logs · traces · model context. DB delete ≠ WAL/backup wipe.

### C5 Boundaries

- Strict JSON; **same parse object** for verify + exec. Reject duplicate keys / bad Unicode / non-canonical protocol values. Registered AgentID lowercase ASCII match; do not lowercase-then-verify arbitrary envelope.
- Registry: integrity-checked revision + implemented subset; ADD intent ≠ auto-exec register.
- **MVP OFF:** evidence auto-fetch/preview; optional webhook.
- Resource / rate / size / concurrency limits: values **TBD**, never unlimited. `priority:high` does not bypass.
- Pending invites do **not** consume active pair slots. Auto reinvite = current active pair only.

### Reject / Defer (constraints)

**Defer (MVP inactive):** evidence remote fetch; optional webhook; `schedule.commit` executable + amend.

**Reject:** NL mid-hop / canonical `body.text` / `deliver:both` / new top-level parent-link / request-envelope `status` / remap intent meaning; using `sig`·hint·`ext`·confidence as owner approval; self-message pairing restore; tool exec from `relay.approval.request`; mutate original to escalate; pre-auth idempotency probe; full-envelope cache hit / re-enqueue / TTL extend; lease grant as `delivered`; forged B-signed `relay.result`; "HTTPS ref is safe"; execute on pre-revoke receipt alone; ACK ⇒ external exactly-once; invent Free **N** / Pro price / rate numbers / extra size caps.

## Non-goals (MVP)


- Latent KV / token-id LLM handoff — research only
- Auth tokens in message body — forbidden
- 장기 채팅 로그, `deliver: both`, 병원·폐쇄망 전제 프로토콜
- amend (승인 수정) — out of MVP
- parent-link top-level 필드
- 릴레이가 B-서명 `relay.result`를 위조하는 transport failure 보고
- rate limit / 추가 봉투 크기 / Free N / Pro 가격 숫자 발명
- MVP evidence 자동 fetch/preview · optional webhook (C5 OFF)
- `schedule.commit` 실행 경로 (stub only)
