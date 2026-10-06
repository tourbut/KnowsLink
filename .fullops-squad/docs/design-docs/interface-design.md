---
id: D05
title: 인터페이스설계서
status: review
updated: 2026-10-06
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-MVP-002-DEV, SAR-MVP-003-BIDIRECTIONAL, SAR-PUBLIC-IDENTITY-001-DEV, SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX, SAR-PUBLIC-AGENTS-001-DEV, SAR-PUBLIC-AGENTS-001-DEV-FIX, SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX]
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

## SAR-MVP-003 시험 메시지 계약

이번 사용자 승인 범위에서 `relay.test.message`를 추가했다. `body`는 `{text:string}`만 허용하며 UTF-8 1–4096 bytes의 비공백 text가 필요하다. deliver는 agent다. reply_to·evidence·ext·render는 허용하지 않는다. envelope 나머지 서명·UUIDv7·idempotency key·TTL 계약은 기존 relay.v1과 같다. registry hash는 변경된 manifest의 compiled SHA256이다.

서버는 명시한 `KNOWSLINK_TEST_AGENTS` 두 agent와 active pair를 검사한다. 시험 envelope는 기존 `/v1/send`에서도 allowlist를 요구한다. machine prefix는 다음 동작만 제공한다. 모든 경로에 agent credential과 시험 allowlist가 필요하다.

| 경로 | 계약 |
|---|---|
| GET /v1/test/registry | 시험 revision manifest·SHA256 |
| GET /v1/test/keys/{agent}/{kid} | 시험 allowlist 상대의 등록 공개키, 기존 pair 검사 |
| POST /v1/test/send | 서명된 relay.test.message만 허용, receipt-only 결과 |
| POST /v1/test/pull | 시험 intent의 lease만 반환 |
| POST /v1/test/persist·ack·claim | 시험 ID만 허용, 기존 lease/recipient/TTL/current 검사 |
| GET /v1/test/receipts/{id} | 자기 endpoint의 시험 receipt metadata |

HTTP 입력은 64 KiB로 제한한다. 시험 payload는 claim 직후 삭제한다. lease/persist/ACK만으로 text를 model에 노출하지 않는다. 시험 부모로 gate나 업무 result를 만들 수 없다. 수동 회신은 별도 trial send이며 첫 ID를 text에 넣어 대조한다.

MCP `knowslink_test_send` 입력은 `{text,idempotency_key}`다. recipient·URL·credential 입력은 없다. `knowslink_test_receive`는 입력이 없다. 성공 출력은 `{state:received,message:{id,from,to,text,exp,untrusted:true}}`이며 빈 queue는 `{state:empty,message:null}`다. 원문 노출 예외는 이 승인된 시험 도구에만 적용한다.

`test-loopback`은 HTTP loopback root만 허용한다. `test-remote`는 `https://link.knowslog.com` root와 별도 CF Access header 두 개를 요구한다. relay request는 `/v1/test/` prefix로 변환한다. redirect와 다른 URL을 거부한다. timeout은 body 스트림까지 요청별 10초이며 응답 상한은 64 KiB다. 인증 실패 상세·lease·claim·credential을 tool error에 반환하지 않는다.

`trial_configured_unverified`는 모드 보고일 뿐이다. 실제 remote 수락은 네 관측 ID와 인증된 호출 증거가 모두 있어야 한다. 세부 변수·Codex launcher·Grok 준비는 [adapter 문서](../../../adapters/README.md#승인된-codexgrok-시험-메시지-sar-mvp-003)를 따른다.

## SAR-PUBLIC-IDENTITY-001 회원 화면과 세션

상위 요구사항은 [D02 일반 이메일 서비스](../planning/product-specs/SAR-PUBLIC-SERVICE.md) PS-01–04다. 회원 화면은 Go `html/template` form이며 JSON API가 아니다. 모든 쓰기는 POST다. cross-origin 브라우저 쓰기는 403이다.

| 메서드·경로 | 인증 | 계약 |
|---|---|---|
| GET / | 없음 | 이메일 입력 화면. 유효 세션이면 303 `/home`. `n=logout|all|expired`는 고정 안내만 표시 |
| POST /auth/start | 없음 | form `email`. 유효하면 코드 메일 발송 후 303 `/auth/verify`와 `__Host-kl_pending`(600s). 형식 오류 422, 한도 429, 발송 실패·미설정 503 |
| GET /auth/verify | pending cookie | 확인 대기 화면. 코드가 없거나 만료면 410 |
| POST /auth/verify | pending cookie | form `code`. 성공 303 `/home`과 `__Host-kl_session`. 오답 401(남은 횟수), 만료·5회 오답·재사용 401, 회원 수용량 503, 비활성 403, 한도 429 |
| GET /home | session | 마스킹 이메일·회원 ID·자기 agent·관계·gate·세션 동작 |
| POST /auth/logout | session | 현재 세션 삭제, 303 `/?n=logout` |
| POST /auth/reauth | session | 회원 이메일로 새 코드 발송, 303 `/auth/verify` |
| POST /auth/logout-all | session | 최근 5분 확인이 있으면 회원의 모든 세션 삭제(303 `/?n=all`). 없으면 403 |
| GET·POST /home/gates/{id} | session | 기존 gate 화면과 결정. CSRF는 세션 token HMAC. 다른 회원 gate 조회 403, 결정 409 |

세션이 없거나 만료된 회원 경로는 303 `/?n=expired`다. HTTP Basic 창을 띄우지 않는다. 회원 세션 cookie는 `/v1` owner·agent API의 credential이 아니다(401). `POST /v1/owners`는 합성 가입 설정이 없으면 403이다.

확인 메일은 `text/plain; charset=UTF-8` quoted-printable이다. 제목은 `KnowsLink 확인 코드`다. 본문은 6자리 코드와 10분·1회 안내다. 링크는 넣지 않는다. 메일 보안 검사기가 링크를 열어 코드를 소모하는 문제를 피한다.

설정은 다음 환경 변수다. 값은 로그와 오류에 출력하지 않는다.

| 변수 | 의미 |
|---|---|
| `KNOWSLINK_SMTP_URL` | `smtps://user:secret@host:465` 또는 `smtp://[user:secret@]host:port`. 비우면 코드 발송은 503 |
| `KNOWSLINK_MAIL_FROM` | 발신 주소 하나. 이름 표기는 거부 |
| `KNOWSLINK_CLIENT_IP_HEADER` | 빈 값 또는 `CF-Connecting-IP`. Tunnel만 relay에 닿을 때만 설정 |
| `KNOWSLINK_SYNTHETIC_SIGNUP` | `1`일 때만 `/v1/owners` 합성 가입 허용. `.env.example`·공개 후보는 비움. 격리 로컬 fixture(`make verify-mvp`, README QA 실행)만 셸에서 `1`을 준다 |

## SAR-PUBLIC-AGENTS-001 회원 agent·키·관계

일반 회원은 owner bearer·Basic을 쓰지 않는다. 기존 Secure 세션과 CrossOriginProtection으로 회원 POST를 보호한다. 권한·한도 실패는 새 자원을 만들지 않으며 거부 요청 rate를 저장한다. 429는 재시도 시각을 표시한다. 정리 요청도 현재 세션·대상 소유권을 요구한다.

| 메서드·경로 | 입력·동작 |
|---|---|
| POST /home/agents | 새 무작위 ID의 미연결 자기 agent. 5분 재인증 |
| POST /home/connect | agent, client=node-local, mode=register 또는 rotate. 5분 재인증. 화면에서만 grant token 1회 표시 |
| GET /home/connections/{grant hash} | 자기 연결 상태·지문·권한·기한. 원 token 재조회 없음 |
| POST /home/confirm | connection=grant hash. 5분 재인증·prepared 상태 검사. 아직 키 미활성 |
| POST /home/cancel | 자기 grant 무효화. 별도 정리 budget |
| POST /home/key-revoke | agent,kid. 5분 재인증·선택 키 철회 |
| POST /home/agent-revoke | agent. 5분 재인증·모든 키/자격/관계 철회·slot 회수 |
| POST /home/invites | agent,target. 자기 발신·상대 active agent 검사. 반복은 현재 pair 반환 |
| POST /home/invite-decision | agent,target,decision=accept 또는 deny,generation. 수신 owner·현재 세대 검사 |
| POST /home/unpair | agent,target,generation. 양측 owner·현재 세대 검사 |
| POST /v1/connect/info | token,client. 유효 grant의 owner·agent·client·mode·기한·상태. 이메일 없음 |
| POST /v1/connect/prepare | token,client,kid,public,proof. 새 키 PoP 저장; 기존 kid 재사용 금지 |
| POST /v1/connect/complete | token,client,proof. approved·PoP·기한·cap 재검사 뒤 key-specific credential을 한 번 반환 |

연결 PoP는 UTF-8 `KNOWSLINK-CONNECT\0<token>\0<owner>\0<agent>\0<client>\0<mode>\0<kid>\0<public>`의 Ed25519 서명이다. public/proof는 base64url-no-pad다. private key 입력은 없다. grant는 10분·1회다. register는 활성 최대3개이며 rotate는 완료 때 기존 키 전체를 철회한다. 발급·준비 실패는 기존 활성 키를 철회하지 않는다.
`/v1/connect/*` 입력은 8192 bytes이며 strict JSON이다. 유효 grant는 해당 회원 principal 40/60s, 무효 grant는 source IP 30/60s를 적용한다. 다른 agent API는 stable agent principal을 쓴다. 신규 전체200·정리 전체100/60s와 principal 정리20/60s를 적용한다. 성공은 200, 회원 상태 변경은303, invalid_auth401·권한403·cap409·입력422·rate429다.
연결 완료 응답을 잃으면 token 재사용으로 credential을 복구하지 않는다. 새 회전으로 복구한다. grant 만료는 완료된 키를 철회하지 않는다. 일반 text·원격 MCP OAuth·외부 계정 성공은 별도 과제다.

### SAR-PUBLIC-AGENTS-001-DEV-FIX 보존 상한과 rate 응답

- `/v1/connect/*`의 429는 `/v1/*` rate와 같은 응답이다. body는 `{"error":"rate_limited","retry_at":"<RFC3339 UTC>"}`이다. `Retry-After`는 같은 시각의 HTTP-date다. 두 값은 실제 재시도 가능 시각을 초 단위로 올림한 값이다. 이전의 고정 60초 header와 retry_at 없는 body는 제거했다.
- 회원 세션 경로는 무효 세션 요청도 source IP 익명 budget 30/60s에 집계한다. 대상은 `GET /home`, `GET /home/connections/{id}`, `POST /auth/reauth`, `POST /auth/logout`, `POST /auth/logout-all`과 기존 회원 POST다. 무효 세션 GET은 기존처럼 `/?n=expired`로 303 이동한다. budget을 넘으면 429와 재시도 시각을 표시한다.
- agent 생성과 키 연결은 기술 보존 상한도 검사한다. owner당 agent 기록(활성+철회) 10개, agent당 키 기록(활성+철회) 20개다. 상한이면 `POST /home/agents`·`/home/connect`·`/v1/connect/complete`·합성 `/v1/agents`·`/v1/keys`가 409 `capacity`다. 철회 요청은 기록을 늘리지 않으므로 상한에서도 허용한다.
- 철회 agent는 철회 시각부터 24h 뒤 키·관계와 함께 삭제된다. 그 뒤 홈 목록과 관계 목록에서 사라진다. 삭제된 agent·pair를 가리키는 옛 화면 요청은 403이다.
- 회원 거부 화면은 같은 안전 문구와 함께 다음 동작을 제공한다. 세션 무효(401)는 로그인 화면 링크다. 그 밖의 거부·429·503은 자기 홈 링크다. 재확인 필요(키·연결 권한 변경, 전체 로그아웃)는 `POST /auth/reauth` 버튼도 표시한다.
- 연결 화면의 취소 버튼은 waiting·prepared·approved에서만 표시한다. 홈의 관계 초대 기한은 연결 기한과 같은 `YYYY-MM-DD HH:MM:SS KST`다. 기한이 없는 기존 관계는 `없음`이다.
- 키 지문은 `<code>`로 표시하고 화면 전체에 `overflow-wrap:anywhere`를 적용한다. select는 본문 글꼴 18px·전체 폭이다. 기존 memberStyle 안의 변경이며 새 theme·의존성은 없다.

### SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX 반복 초대·기록 포화 안내

- `/v1/*` wire·상태 코드·오류 코드는 바꾸지 않았다. 아래는 회원 HTML 경로의 안내 변경이다.
- `POST /home/invites` 성공은 `/home?n=invited`(새 pending), `/home?n=invite-pending`(같은/반대 방향의 기존 pending 유지), `/home?n=invite-active`(현재 active 유지)로 303 이동한다. 홈은 해당 안내를 `role="status"`로 표시한다. 반복 제출은 기한·세대를 바꾸지 않는다.
- 홈 관계 목록: 받은 pending은 수락 전 메시지 불가를, 보낸 pending은 반복해도 기한이 바뀌지 않음을 안내한다. 거절·만료·철회 관계는 연결 종료·메시지 불가·새 초대와 새 수락 필요를 안내한다. 자기 살아 있는 agent가 있으면 같은 상대에 대한 `새 초대 보내기` 버튼(기존 `POST /home/invites`)을 표시한다. 종료 관계에는 관계 철회 버튼을 표시하지 않는다. 자동 재연결은 없다.
- `POST /home/agents` 409: owner 기록 포화는 철회 기록의 최소 24시간 보존·정리 뒤 홈에서 사라진 다음 다시 시도하라고 안내한다. 활성 5개는 철회 선택과 복구 불가를 안내한다. 전체 200개는 기존 운영 한도 문구다.
- `POST /home/connect` 409: 키 기록 포화는 키 철회·대기로 공간이 생기지 않음과 새 agent 연결·각 상대 새 수락을 안내한다. 새 agent를 만들 수 있으면 거부 화면에 `새 agent 만들기` 버튼을 함께 표시한다. 만들 수 없으면 그 한도 안내를 덧붙인다. 활성 키 3개의 등록은 키 철회 또는 회전을 안내한다.
- 홈 agent 카드: 철회 agent는 최소 24시간 보존 뒤 정리되면 목록에서 사라질 수 있음과 권한 복구·백업 영구 삭제가 아님을 표시한다. 키 기록 포화 agent는 연결 수단 발급 대신 교체 안내를 표시한다. agent 철회 버튼 앞에 모든 키·관계 종료와 복구 불가를 표시한다.
- 상품 quota·결제·정확한 정리 시각은 표시하지 않는다. 상대 이메일·회원 존재·추가 기술값은 노출하지 않는다.
