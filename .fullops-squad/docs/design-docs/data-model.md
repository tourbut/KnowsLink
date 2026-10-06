---
id: D06
title: 엔티티정의서
status: review
updated: 2026-10-06
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-PUBLIC-IDENTITY-001-DEV, SAR-PUBLIC-AGENTS-001-DEV]
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
| Message | ID, frozen raw Envelope, shared Inbox, receipt, kid, deliver 경로(값 없는 이전 상태는 agent 처리 불가), pair generation, attempts, lease token/window/generation, claim token, completion | MVP-03–07/11/15 |
| Receipt | id/from/to/intent/digest/exp/accepted_at/transport | MVP-04/05 |
| Idempotency | from/key와 Message ID 연결 | MVP-04 |
| Gate | H.id, owner, M.id, M.digest, endpoints, generation, policy, exp, state, consumed | MVP-09/10/16 |
| Member | 회원 ID `mem_…`, Owner, 정규화 Email, Issuer, Active, Created | PS-01/02/04 |
| Identity | `issuer|email` → 회원 ID. 같은 신원의 회원은 하나 | PS-02 |
| Session | token SHA256 키, 회원 ID, Created·Seen·Verified. 절대 12h·무활동 60분 | PS-03 |
| Challenge | pending token SHA256 키, Email, `SHA256(pending:code)`, Exp(10분), Attempts(최대 5) | PS-01, PS-11 |
| Rate | bucket 키(이메일은 SHA256)와 시각 목록. 1h 뒤 정리 | PS-11 |

owner와 agent credential 원문은 발급 응답과 로컬 QA fixture에만 존재한다. DB는 hash만 보관한다.
lease·claim token은 짧은 유효 권한이며 해당 message 안에 결속한다. 이를 로그에 남기지 않는다.
Inbox는 원문과 같은 Postgres transaction에서 저장한다. ACK는 Inbox commit 뒤에만 성공한다.


회원 owner는 bearer credential이 없다. 회원 이메일은 로그인 매핑에만 쓴다. agent ID·contacts·receipt·rate key·로그에 넣지 않는다. 미확인 이메일은 Challenge의 최대 10분 동안만 남는다. 확인 코드와 세션 token 원문은 DB에 저장하지 않는다. 계정 비활성화·30일 정리는 후속 기능이다.

DB 정본은 [D07](database-design.md), CRUD 정본은 [D09](crud-design.md)다.

## SAR-PUBLIC-AGENTS-001 연결·권한 필드

- Agent.Revoked: 자기 agent 전체 철회. owner 비활성 또는 revoked agent는 인증·현재 메시지·pair 대상에서 제외한다. agent 생성 자체는 활성 등록 slot을 소비하며 키가 없으면 미연결이다.
- Key.Credential: 공개 회원 연결의 key별 bearer SHA256. 기존 Key.Public/Revoked/Changed는 유지한다. 일반 회원 Agent.Credential은 비어 있다. 키 철회 기록은 kid 재할당 금지와 최소24h 철회 근거로 보존한다.
- Connection map key: random grant token의 SHA256. 값은 Owner/Agent/Client/Mode/Kid/Public/State/Exp다. private key와 원 token은 저장하지 않는다. waiting/prepared/approved/consumed/cancelled/expired 상태를 쓴다. token 만료는 기존 연결 키의 수명을 바꾸지 않는다. 만료 뒤24h에는 연결 상태 기록을 제거한다.
- Pair.Exp: pending 수명24h. 기존 Exp가 없는 pending은 첫 sweep에서24h를 부여한다. 만료/거절/철회 뒤 재초대는 Generation 증가다. 활성 pair의 Exp는 초대 당시 표시용이다.
- Rate: 같은 stable 회원/agent principal의 key가 회전·재시작으로 초기화되지 않는다. 거절·철회는 cleanup bucket을 쓴다. 이메일은 agent ID·grant·pair·contacts에 없다.
회원 비활성화 UI와 최대30일 계정 매핑 삭제는 이번 과제에서 새로 구현하지 않았다. 기존 owner 비활성 경계와 새 agent/key/pair 철회 경계는 유지한다. 운영 보존량 상한·복원 뒤 철회 검증은 OPS/공개 수락 후속이다.
