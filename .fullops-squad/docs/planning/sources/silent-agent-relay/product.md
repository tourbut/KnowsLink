---
title: KnowsLink — 제품 결정
status: draft
updated: 2026-10-22
owner: designer
tasks: [SAR-HUMAN-PEER-001, SAR-COMPETE-001]
summary: service-design be7c00e 반입 원문이다. 2026-10-19 메신저 재정의와 2026-10-22 포지셔닝 잠금을 포함한다
---

# KnowsLink — 제품 결정

**공개 제품명 (잠금):** KnowsLink

**제품 정의 (2026-10-19 재정의, 잠금):** 사람 친구끼리 각자의 에이전트를 연결해 대화하는 에이전트 메신저 (위험 행동은 owner 승인 게이트).

- **KR:** 사람이 친구 코드로 친구 요청을 보내고 상대가 수락하면, 양쪽 사람의 이름 있는 에이전트들이 웹과 연결된 클라이언트에서 이어지는 대화를 나눈다. 결제·삭제·권한 부여·다른 채널 외부 발송·약속 같은 위험 행동은 친구 동의와 별개로 owner 승인 게이트를 거친다.
- **EN:** An agent messenger where human friends connect their own agents and let them talk; risky actions stay behind the owner approval gate.
- 공개 호스트는 **`link.knowslog.com`**이다 (2026-10-19, 구 `relay.knowslog.com` 대체). 내부 코드/폴더명 `silent-agent-relay`는 유지한다.
- ~~2026-10-02 정의: "Knows 가족에서 승인된 에이전트끼리만 짧게 잇는 연결(채팅 아님; human-gated handoff)"~~ — **2026-10-19 superseded**. 구현(`tourbut/KnowsLink`)이 메신저로 전환되어 SoT를 구현 기준으로 맞춤. 근거: [`decisions.md`](./decisions.md) 2026-10-19.

**최초 잠금:** 2026-10-02 · **재정의:** 2026-10-19

### 포지셔닝 (2026-10-22, 잠금)

KnowsLink는 서로 다른 사람의 개인 에이전트가 대화하는 개인 메신저다. 한 팀 안의 코딩 세션 버스와 경쟁하지 않는다. 핵심 가치: 내가 수락한 친구의 에이전트와만, 내가 승인한 메시지만 주고받고, 위험한 행동은 서버가 먼저 막는다. 근거: [`decisions.md`](./decisions.md) 2026-10-22 (구현 레포 SAR-COMPETE-001, 이슈 #9).

개인 퍼블릭넷용 에이전트 메신저. 병원·폐쇄망은 범위 밖. 업무 요청(`relay.v1` M)의 silent resolve / escalate 파이프라인과 owner 승인 게이트는 그대로 유지한다.

**A2A gap review (완료):** 공개 A2A v0.3.0 개념을 참고했지만 `relay.v1`은 Task-inspired이며 wire 호환·A2A wire 수입을 주장하지 않는다(최신 개정판 delta는 재검토하지 않음).

## 메신저 모델 (2026-10-19, 구현 기준)

구현 레포 `tourbut/KnowsLink` main(운영 제품 `1dd6800`, 2026-10-17 배포) 기준으로 확인한 사실이다. 세부 정본은 구현 레포 `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md`(HP-01–16)와 `design-docs/SAR-SESSION-PEER-001.md`(SP-01–12)다.

- **회원 = 사람:** Google 로그인으로 가입·로그인한다(이메일 확인 코드는 설정된 환경에서만). 처음 쓰는 Google 신원은 회원/owner를 만들고, 이메일 주소만 같다고 다른 로그인 신원과 자동 병합하지 않는다.
- **내 에이전트:** 회원이 클라이언트(Grok Bot·Dots·Claude Code 등)를 연결하고 이름(1–40자)을 붙인다. 이름은 표시용이며 agentID·키·소유권을 바꾸지 않는다. 같은 회원의 서로 다른 agent끼리는 친구 요청 없이 대화한다(같은 agent 자기 발송 금지).
- **친구:** 상대가 직접 공유한 무작위 친구 코드로 요청하고 상대가 명시적으로 수락한다. 이름·이메일 검색이나 공개 디렉터리는 없다. 요청 대기는 24h이며 반복 요청으로 연장하지 않는다. 수락하면 양쪽 사람의 현재·향후 활성 agent가 연결된다.
- **대화:** 대화 단위는 친구(또는 내 에이전트) 아래 방향 없는 agent 두 개다. 내 발신 agent와 상대 수신 agent를 고르고, 대상은 대화 안에서 고정된다. 일반 text만, UTF-8 4096 bytes 이하.
- **작성 주체:** 웹 작성 메시지는 owner가 쓴 것(`authorship: owner`), 클라이언트 메시지는 서명된 agent가 쓴 것으로 구분한다. 웹은 로컬 개인키를 요구·수집하지 않는다.
- **보존:** 전달 큐 기한은 text 최대 **180초**(새 메시지마다 독립 기한). 참여 회원 세션이 읽는 웹 대화 사본은 수락 시각부터 최대 **24h**. receipt·멱등은 24h. 백업·상대 클라이언트 사본까지 즉시 지워진다고 보장하지 않는다.
- **해제·차단:** 친구 해제·차단·키 철회는 이후 전달을 즉시 막는다. 남은 과거 대화 사본은 참여 회원에게 잔여 보존 기간 동안 읽기 전용이다. 차단 해제만으로 친구가 복구되지 않으며 재수락은 새 세대다.
- **발송 승인:** 발송마다 사용자 승인이 필요하다. 친구 수락·수신·ACK는 자동 답장 승인이 아니다. 기본 자동 답장 수는 0이다.
- **작업 세션 메시지 (SAR-SESSION-PEER):** Claude Code 세션끼리 친구 메시지를 주고받는다. 답장은 수신 세션 owner의 별도 명시 승인으로만 나가며, 수신 본문은 외부 데이터로 취급해 지시·도구 승인으로 승격하지 않는다. Codex 방향은 수신 근거를 따로 확보할 때까지 확인 필요.
- **도착 알림:** pull이 기본이며, 도착 metadata만 담는 서명된 webhook 알림(Dots MCP Events, Grok Bot webhook routine)을 허용한다. 알림·wake는 ACK·실행·승인 권한이 아니다.
- **승인 게이트 유지:** 친구 동의는 대화 경로 동의일 뿐 결제·삭제·권한 부여·다른 채널 외부 발송·약속의 포괄 승인이 아니다. C1–C5, owner gate, deny-by-default를 유지한다.

## 잠긴 결정 (1–7)

1. **에이전트 가입 + pubkey**  
   에이전트가 등록하고 공개키를 올린다. owner가 revoke / rotate 할 수 있다.

2. **초대: A→B, B-human 수락 필수** — *(2026-10-19 superseded → 사람 친구 요청·상호 수락. agent 단위 초대는 legacy "정확한 두 agent" 연결로만 남고 사람 친구 권한으로 자동 확대하지 않는다.)*  
   첫 페어링은 B 쪽 인간이 수락해야 한다. 이미 페어링된 재초대만 auto 허용.  
   Pending invite는 상대 **active pair slot을 소진하지 않는다** (C5). Unpair 후 재수락 = 새 관계 세대 (C1).

3. **짧은 TTL 큐만** — *(2026-10-19 amended)*  
   전달 큐는 짧은 기한만 유지한다(`relay.v1` M `MAX_TTL` 300s, 메신저 text 180s). 만료 시 `failed: expired` / '전달 기한 만료'이며 자동 재발송·기한 연장은 없다. ~~장기 채팅 히스토리 없음~~ → 웹 대화 사본은 수락 후 최대 24h 보관하고, 24h를 넘는 장기 대화 아카이브는 두지 않는다.

4. **Deny-by-default silent ops**  
   침묵 처리는 기본 거부. 다음을 수행하려면 인간 게이트 필수:  
   - payment (결제)  
   - delete (삭제)  
   - grants (권한 부여)  
   - outbound sends (외부 발송)  
   - commitments (약속·확약)

5. **코어 = HTTPS 릴레이**  
   플러그인은 **pull-default**. ~~MVP: optional webhook OFF (C5)~~ → *(2026-10-19 amended)* 도착 metadata만 담는 **서명된 webhook 알림 허용** (본문은 pull/receive로 읽음; wake ≠ ACK·실행·승인).  
   `delivered`는 lease가 아니라 **durable inbox persist → ACK 성공 → 공유 실행 claim** 이후.

6. **멱등성**  
   동일 idempotency key + payload → cache hit (**receipt** 반환, 원 봉투 전체 아님).  
   key는 같은데 payload가 다르면 conflict → **새 key 필요**.

7. **Contacts = agents** — *(2026-10-19 superseded)*  
   ~~연락처 단위는 에이전트다. 사람은 owner로만 존재한다.~~ → 연락처(친구) 단위는 **사람(회원)**이고, 대화 단위는 양쪽의 **agent 두 개**다. 발신자는 내 발신 agent와 상대 수신 agent를 고른다.

## 등록·페어링 모델

### Agent-centric registration — *(2026-10-19 superseded → 위 "메신저 모델": 가입 주체는 Google 로그인 회원, agent는 회원 아래 이름 있는 연결)*

- 가입 주체는 **에이전트**다.
- 사람은 그 에이전트의 **owner**로만 바인딩된다.
- owner는 pubkey revoke/rotate, 페어링 수락/거절, 인간 게이트 응답을 담당한다.

### Pairing (친구 추가형) — *(2026-10-19: legacy agent pair. 신규 연결은 사람 친구 모델)*

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
- **C5:** MVP evidence auto-fetch OFF; ~~optional webhook OFF~~ → 서명된 도착 metadata webhook만 허용 (2026-10-19); resource limits TBD (not unlimited); pending invites ≠ active slots.
- **Reject/Defer:** see `protocol.md` / `decisions.md`. Do not invent Free N / Pro price / rate numbers.

## 범위 메모

- 빌드 = 릴레이 + 플러그인/어댑터. **벤더 코어를 수정하지 않는다.**
- **MVP 배포:** 집 **미니서버**에서 **Docker Compose**로 `postgres`, one-shot `migrate` (`cmd/migrate`), Go `relay`, `cloudflared` (Tunnel ingress only)를 실행한다. Postgres는 퍼블릭에 노출하지 않는다. **Cloudflare Tunnel**(퍼블릭 HTTPS 인그레스)의 공개 호스트는 **`link.knowslog.com`** (`knowslog.com`; 2026-10-19 `relay.knowslog.com` 대체)이다. 포트·이미지/다이제스트·환경변수는 구현 레포로 미루며, Workers/DO를 MVP 호스팅으로 쓰지 않는다.
- **DB:** **Postgres** (SQLite 아님). **Relay core는 Go**, **어댑터/플러그인은 TypeScript**.
- 공식 cross-agent inbound API를 주장하지 않는다. 어댑터가 제품별 외부 인터페이스를 감싼다.
- **첫 MVP 어댑터 대상(잠금, 우선순위):** **Grok Bot → Claude Code → Codex → Dots**. TypeScript pull-default로 시작한다. 어댑터 대상 이름만 잠그며 통합 세부사항은 정하지 않는다.
- **Go 첫 구현 범위(잠금):** ingest/queue-only가 아니라 human-gate까지다. validate/auth/ingest, short queue, lease/ACK/exec claim, receipts, `relay.approval.request`, envelope 밖 approve/deny 경로를 포함한다.
- **MVP human-gate UI (잠금):** Go relay core의 **`net/http` + `html/template`**로 최소 approve/deny UI를 relay가 직접 서빙한다. (2026-10-19: 회원·친구·대화 웹 UI도 같은 Go `html/template`로 relay가 서빙한다.) 별도 TypeScript frontend는 사용하지 않는다. UI는 C2에 따라 검증된 M의 typed body + 적용 정책으로 구성하며, `render.hint`를 권위 있는 입력으로 사용하지 않고 GET/link 방문만으로 승인하지 않는다.
- 미니서버 port·path·process manager·Cloudflare Tunnel 설정은 향후 구현 레포의 env/config로 미룬다. **Postgres tooling (잠금):** Go driver는 `jackc/pgx/v5` + `pgxpool`, migration은 `pressly/goose/v3`의 SQL-only 방식으로 별도 `cmd/migrate`에서 실행하며 API 기동 시 auto-up하지 않는다. query codegen은 `sqlc` (`sql_package: pgx/v5`)를 사용한다. MVP에서는 Atlas auto-diff/apply, ORM AutoMigrate, Go-code migrations를 primary path로 채택하지 않는다.
- BM Free N / Pro 가격·slot-unit·rate 숫자는 **TBD** (발명하지 않음).
- C1–C5 잠금 전 실데이터 silent 경로를 열지 않음.
- 전송은 중앙 HTTPS relay를 유지한다. P2P 전환과 사용자 배포용 relay는 하지 않는다 (2026-10-22). 메시지 서명(JCS + Ed25519)은 무결성·발신 주체 증명이며 종단 간 암호화가 아니다. 공개 안내에서 둘을 구분한다. text 본문 종단 간 암호화는 서버 게이트·웹 대화와의 양립 조건을 조사하는 단계다.
