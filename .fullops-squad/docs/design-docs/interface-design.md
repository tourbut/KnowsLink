---
id: D05
title: 인터페이스설계서
status: review
updated: 2026-10-03
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-MVP-002-DEV]
upstream: [D02]
summary: owner와 agent HTTP 계약 및 gate와 adapter 흐름을 정의한다
---

# KnowsLink HTTP·adapter 계약

## 인증과 오류

상위 요구사항은 [D02 SAR-MVP](../planning/product-specs/SAR-MVP.md)다.
API는 `Authorization: Bearer <credential>`을 사용한다. owner와 agent credential은 서로 대체할 수 없다.
UI는 HTTP Basic의 password에 owner credential을 넣는다. username은 owner ID를 사용한다.
모든 응답은 no-store다. HTML은 CSP·Referrer-Policy·escaping을 사용한다.
성공은 HTTP 200이며 gate 결정은 303이다. 인증은 401, 권한은 403, 경합은 409, schema는 422, DB/시계 불명은 503이다.
JSON 오류는 `{error: code}`다. 내부 경로·DB 오류·비밀값을 노출하지 않는다.

## owner·등록·pairing — MVP-01/02/07

| 메서드·경로 | principal | 입력·결과 |
|---|---|---|
| POST /v1/owners | 로컬 합성 가입 | `{}`; 새 owner ID·credential |
| POST /v1/agents | owner | `{agent,kid,public,proof}`; agent credential |
| POST /v1/keys | owner | 같은 key 입력; rotate, 이전 키 revoke |
| POST /v1/key-revoke | owner | `{agent,kid}`; 자기 키 revoke |
| POST /v1/owner-revoke | owner | `{}`; owner 권한 revoke |
| POST /v1/invites | A-agent 또는 A-owner | `{agent,target}`; pending 또는 현재 active pair |
| POST /v1/invite-decision | B-owner | `{agent,target,decision:accept|deny}` |
| POST /v1/unpair | pair의 owner | `{agent,target}`; revoke |
| GET /v1/contacts | agent | active 상대 AgentID 목록 |
| GET /v1/keys/{agent}/{kid} | 자기 또는 active pair agent | 등록 공개키; URL/path kid lookup 없음 |

`public`와 `proof`는 base64url-no-pad다. PoP bytes는 다음 UTF-8 문자열이며 끝 NUL은 없다.
`KNOWSLINK-KEY-POP\0<owner ID>\0<AgentID>\0<kid>\0<public>`
`proof`는 새 private key의 Ed25519 서명이다. 같은 kid는 revoke 뒤에도 재할당하지 않는다.
agent ID는 lowercase ASCII다. 실제 사용자 신원 확인·회복·credential 재발급은 후속 운영 인증 과제다.

## transport·실행 조정 — MVP-03–07/11

| 메서드·경로 | principal | 계약 |
|---|---|---|
| GET /v1/registry | static | manifest·SHA256; compiled seed hash 검증 |
| POST /v1/send | envelope.from agent | frozen relay.v1; receipt-only 결과. H 외 `deliver:human`은 403 sender_not_allowed |
| POST /v1/pull | recipient agent | `{}`; envelope·lease_token·lease_until·generation·attempts 또는 null |
| POST /v1/persist | recipient agent | `{id,token}`; shared durable inbox commit |
| POST /v1/ack | recipient agent | `{id,token}`; persisted·유효 lease 확인 후 delivered |
| POST /v1/claim | recipient agent | `{id}`; 1회 claim token |
| POST /v1/authorize | claim 보유 agent | `{id,claim}`; current-auth 재검사·executable/disclosure false |
| GET /v1/receipts/{id} | 현재 endpoint agent | receipt metadata와 별도 completion |

pull·persist·ACK·claim은 저장된 `deliver:agent` 메시지만 처리한다. human 전달은 owner gate 결정만 delivered로 바꾼다 (C1).
직접 human inbox가 없으므로 H 외 `deliver:human` 요청은 수락하지 않는다. 배포 전 저장된 경로 미기록 메시지도 agent가 처리하지 못한다.
이미 발급된 human 또는 경로 미기록 claim은 authorize·gate-consume·H·R 부모로 사용할 수 없으며 403 sender_not_allowed다.
H/R의 새 send는 wire 밖 `X-Execution-Claim` header에 부모 claim token을 요구한다.
동일 서명 재전송은 인증·routing 뒤 receipt만 반환한다. 새로운 승인·실행권을 만들지 않는다.
receipt transport는 queued/leased/delivered/failed:expired/failed:max_attempts/failed:revoked다.
processing 결과는 R의 body.status다. relay transport 실패를 B 결과로 서명하지 않는다.
R에 result/error optional 정보가 있으면 closed schema에서 거부한다. done은 disclosure_denied다.

## human-gate — MVP-08–11/16

| 메서드·경로 | principal | 계약 |
|---|---|---|
| GET /owner | owner | 자기 agent와 gate 목록 |
| GET /owner/gates/{id} | gate owner | verified typed body·policy·상태·만료 |
| POST /owner/gates/{id} | gate owner | form csrf·decision=approve|deny; 원자적 결정 |
| GET /v1/gates/{id} | B-agent | state·consumed·parent metadata |
| POST /v1/gate-consume | B-agent claim | `{id,claim}`; approved 1회 consume, 효과·공개 false |

owner POST decision은 H의 인증된 human ACK를 함께 확정하고 H 원문을 지운다. GET 방문은 결정하지 않는다. CSRF는 owner credential과 gate ID에 묶인 HMAC이다.
UI는 pending/approved/denied/expired/revoked/unavailable을 텍스트로 표시한다.
원문 없음·철회·만료·이미 결정됨은 활성 버튼이 없다. render.hint·ext·agent sig는 owner 증명이 아니다.

## TypeScript 합성 경로

`RELAY_URL`, `AGENT_CREDENTIAL`, `AGENT_ID`, `AGENT_KID`, `AGENT_KEY_FILE`이 필요하다.
URL은 loopback만 허용한다. PEM은 로컬 파일에서 읽으며 로그에 쓰지 않는다.
`ADAPTER_GATE=1`은 합성 judgment gate를 만든다. gate 조회 간격은 2초다.
registry·B가 받은 signature를 검증한 뒤 persist·ACK·claim을 수행한다.
query는 무정책 denied다. 승인된 query도 denied다. commit은 non-executable stub이다.
실제 벤더 inbound API·MCP 연결 성공은 주장하지 않는다. evidence URL은 읽지 않는다.

## SAR-MVP-002 Grok Bot MCP 계약

사용자가 공식 `docs.x.ai/grok-bot` 제품을 확정했다. 실제 연결은 held다. [설치 문서](../../../adapters/README.md)와 [조사·검증 기록](../exec-plans/phases/SAR-MVP-002-DEV.md)을 따른다.

Cursor plugin manifest·stdio MCP·skill을 패키지에 포함한다. SDK는 MCP wire만 처리한다. frozen relay.v1·owner/agent credential·shared durable inbox·ACK/claim·gate·무정책 deny는 기존 Go/TypeScript 계약을 유지한다. 임의 inbound endpoint·natural-language wire를 추가하지 않는다.

`knowslink_status`와 `knowslink_pull_once`는 입력 없는 도구다. status는 held/synthetic_only만 보고한다. pull은 기본 held(isError=true)이며 `KNOWSLINK_MODE=synthetic-loopback`만 허용한다. 설정은 서버 환경으로 제공한다. HTTP loopback root만 수락하고 URL userinfo·path·query·fragment와 redirect를 거부한다.

합성 pull은 한 delivery의 검증·persist·ACK·claim 후 owner gate를 만든다. 승인 뒤에도 최소 denied R을 보낸다. control result는 ACK까지만 수행하고 재응답하지 않는다. tool 출력은 state·transport·actualConnection·webhook·evidenceFetch만 포함하며 원문·claim·lease·gate ID·credential을 숨긴다. processing busy/failure/empty는 업무 done과 구분한다. 같은 프로세스의 동시 pull은 busy로 거부한다. 공유 claim은 계속 relay가 집행한다.

Grok Bot Auto Review/Allow once는 KnowsLink owner approve를 대신하지 않는다. 공식 Bot 앱의 실제 도구 검색·hosted Node·network·credential은 후속 확인 대상이다. Cursor IDE 로컬 plugin loading 검사를 Bot 설치 성공으로 표시하지 않는다.
