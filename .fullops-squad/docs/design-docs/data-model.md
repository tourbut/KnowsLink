---
id: D06
title: 엔티티정의서
status: review
updated: 2026-10-03
owner: dev
tasks: [SAR-MVP-001-DEV]
upstream: [D02]
summary: shared JSON 업무 엔티티와 권한 및 보존 경계를 정의한다
---

# KnowsLink 업무 엔티티·DB·CRUD

## 엔티티 — D06

상위 요구사항은 [D02 SAR-MVP](../planning/product-specs/SAR-MVP.md)다.
업무 값은 `relay_state.data`의 JSON에 저장한다. 구조 정본은 `internal/relay/store.go`다.

| 엔티티 | 키·내용 | 요구사항 |
|---|---|---|
| Owner | ID, credential SHA256, Active | MVP-01/07 |
| Agent | AgentID, Owner, credential SHA256, kid별 Keys | MVP-01 |
| Key | 등록 공개키, Revoked, Changed; kid 영구 재할당 금지 | MVP-01/07 |
| Pair | 정렬한 A/B ID, inviter/recipient, state, generation | MVP-02/07 |
| Message | ID, frozen raw Envelope, shared Inbox, receipt, kid, deliver 경로, pair generation, attempts, lease token/window/generation, claim token, completion | MVP-03–07/11/15 |
| Receipt | id/from/to/intent/digest/exp/accepted_at/transport | MVP-04/05 |
| Idempotency | from/key와 Message ID 연결 | MVP-04 |
| Gate | H.id, owner, M.id, M.digest, endpoints, generation, policy, exp, state, consumed | MVP-09/10/16 |

owner와 agent credential 원문은 발급 응답과 로컬 QA fixture에만 존재한다. DB는 hash만 보관한다.
lease·claim token은 짧은 유효 권한이며 해당 message 안에 결속한다. 이를 로그에 남기지 않는다.
Inbox는 원문과 같은 Postgres transaction에서 저장한다. ACK는 Inbox commit 뒤에만 성공한다.


DB 정본은 [D07](database-design.md), CRUD 정본은 [D09](crud-design.md)다.
