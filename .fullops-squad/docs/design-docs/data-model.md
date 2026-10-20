---
id: D06
title: 엔티티정의서
status: review
updated: 2026-10-12
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-PUBLIC-IDENTITY-001-DEV, SAR-PUBLIC-AGENTS-001-DEV, SAR-PUBLIC-AGENTS-001-DEV-FIX, SAR-PUBLIC-MESSAGES-001-DEV, SAR-PUBLIC-MESSAGES-001-DEV-FIX, SAR-PUBLIC-MESSAGES-001-DEV-FIX-2, SAR-PUBLIC-MESSAGES-001-DEV-FIX-3, SAR-GOOGLE-CONNECT-001-DEV, SAR-SSE-INSTALL-001-DEV, SAR-HUMAN-PEER-001-DEV]
upstream: [D02]
summary: shared JSON 업무 엔티티와 권한 및 보존 경계를 정의한다
---

# KnowsLink 업무 엔티티·DB·CRUD

## SAR-HUMAN-PEER-001 DEV 후보

| 추가 필드 | 저장·권한 의미 |
|---|---|
| State.HumanModel | 현재1; 최초 transaction에서 legacy를 안전 변환. 알 수 없는 버전은503 |
| Member.Name/FriendCode/SelfConnect | 표시 이름, 무작위256bit 공유 코드, 자기 endpoint 범위 동의. issuer 기반 신원 불변 |
| Agent.Name/Client | 표시 이름·연결 클라이언트. 이름 중복 허용, ID/키와 별개 |
| State.Friends | 정렬한 회원쌍 → Friend. public random ID, A/B/Inviter, State, 당시 NameA/NameB, Generation, Exp, BlockA/BlockB, SuppressLegacy |
| Message.Scope | self:/friend: 출처 또는 빈 legacy. Generation과 함께 current 검사 |
| Message.Conversation | OwnerA/B, 당시 PersonA/B·AgentA/B, Text, Authorship(owner/agent). 신규 text만 수락부터24h |

친구 pending은24h, active 전체400/회원20, pending 전체200/발신10/수신10이다. 종료 tombstone은 deny/세대/차단 근거이므로 자동 부활하지 않는다. 친구 기록20000 보호 상한은 새 요청만 막는다. 숨긴 가짜 outgoing placeholder도 실제 pending 한도에 포함한다. 요청자에게 대기로 보이는 기한 전 `blocked-pending`(상대 차단·수신 한도 보류)도 pending 전체·발신 한도에 포함하고 수신 한도에는 넣지 않는다. 이전 코드가 차단·수신 한도 경로에서 만든 대리 행(같은 코드 2행)은 지우거나 합치지 않는다. 그 대리 행이 대기 중이면 같은 코드 재요청은 무동작이고, 끝난 뒤에는 쌍 행만 재사용한다. 차단 상태는 상대 화면에 공개하지 않는다. 친구 코드는 권한이 아닌 요청 수단이다.

옛 Member/Agent의 이름은 빈 값이고 관리 화면에서 설정한다. 기존 회원 SelfConnect=false, 기존 pending Pair=expired, 기존 active Pair는 정확한 쌍만 유지한다. 기존 원문·업무 메시지를 새 대화 사본으로 소급 복사하지 않는다. SyntheticSignup은 runtime-only(`json:"-"`)이며 production 인가에 저장되지 않는다.

## SAR-SSE-INSTALL-001 SSE 상태

`State.Streams`는 무작위 입장 ID를 키로 하는 `{Agent,Exp}` map이다. 원문·credential·키는 넣지 않는다. singleton lock에서 전체16/agent당2를 검사하고65초 기한을 부여한다. 연결은55초 이내에 끝나며 정상 취소/종료 때 해당 기록을 삭제한다. 정리 DB 장애는 Exp 이후 기존 sweep이 회수한다. 이전 JSON은 nil map으로 읽고 최초 입장 때만 생성한다. 기존 Agent/Key/Pair/Message/Inbox는 보존한다. 실제 table·migration·sqlc·D08은 바꾸지 않는다.


### Remote OAuth·Cloud 설치

선택 JSON `Remote`는 Clients/Attempts/Codes/Tokens/Installs/Subscriptions를 보유한다. token/code는 hash 인덱스다. 설치는 owner/별도agent/key/개인키/만료/private Inbox를 보유하며 기존 local private key를 복사하지 않는다. subscription은 installation binding, callback/secret/회전secret, 만료, pending/Seen, 안정 event ID, lease/attempt/next 시간을 보유한다. 철회·만료는 현재 owner/agent/key를 검사하고 설치 비밀·Inbox·구독을 제거한다. 관계 자동 생성은 없다.

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

## SAR-PUBLIC-AGENTS-001-DEV-FIX 철회 기록 보존

- Agent.Changed: agent 철회 시각이다. 반복 철회는 이 값을 바꾸지 않는다. 이 필드가 없는 기존 철회 agent는 첫 sweep 시각을 철회 시각으로 기록한다. 따라서 배포 직후 삭제되지 않고 24h를 보존한다.
- 철회 agent는 Changed부터 24h 뒤 sweep에서 삭제된다. 그 agent의 Key와 그 agent를 포함한 Pair도 같은 sweep에서 삭제된다. 회원 agent ID는 무작위이며 다시 발급되지 않는다. 따라서 삭제로 `(from,kid)`가 재할당되지 않는다(C1).
- 살아 있는 agent의 철회 Key는 삭제하지 않는다. kid 영구 재할당 금지를 지킨다. Key.Changed는 첫 철회 시각이며 이후 회전·반복 철회로 바뀌지 않는다.
- 살아 있는 두 agent 사이의 거절·만료·철회 Pair는 유지한다. 재초대는 계속 이전 Generation보다 큰 세대를 받는다.
- 기술 보존 상한: owner당 agent 기록 10개(활성 5 포함), agent당 키 기록 20개(활성 3 포함)다. 제품 quota가 아니다. 상한은 새 agent·키만 거부하고 철회는 허용한다. 키 기록이 찬 agent는 새 agent로 교체한다.

## SAR-PUBLIC-MESSAGES-001 엔티티·보존

Message는 frozen 업무 또는 별도 knowslink.text 원문을 저장한다. Message.Parent는 text 원요청 ID, ReplyID는 수락된 관련 답장 ID다. text에는 claim을 만들지 않는다. ACK 후 completion=received이고 답장 ACK 후 원요청 completion=reply_received다. 원문은 ACK·TTL·철회·lease 실패에 지운다. receipt/멱등은 수락부터24h이며 개인정보·text·credential을 metadata에 넣지 않는다.
State.HTTP는 무작위 입장 token을 키로 하는 `{Clean,Exp,Owner}` map이다. Owner는 Clean 기록의 공정성 단위다. owner 자신의 제어는 owner ID, 그 owner의 agent ACK는 `<owner>/agents`이며 단위당 Clean 기록은 동시 1개다(SAR-PUBLIC-MESSAGES-001-DEV-FIX-2, DEV-FIX-3). 신규 기록과 이전 기록은 Owner를 생략한다. HTTP 신규16·정리4의 공유 수용량을 센다. Clean=true는 검증된 principal의 자기 기록 정리 요청만 기록한다(SAR-PUBLIC-MESSAGES-001-DEV-FIX). 로컬 채널 선택용 커밋 snapshot((commit 시각, epoch) 순서)·자격 색인·단위별 로컬 정리 표시·마지막 재조회 시작 시각·종료 실패 token은 프로세스 메모리의 파생값이며 저장하지 않는다. 정상 종료 때 삭제하고 Exp(30s) 뒤 sweep한다. 기존 직렬화 상태는 빈 map/ReplyID로 복원한다. Member·Owner·Agent·Pair와 기존 key 보존 규칙은 바꾸지 않는다.

## Google 연결 상태 — SAR-GOOGLE-CONNECT-001-DEV

기존 relay state JSON의 Connection에 `Device bool`을 추가하고 GoogleAttempt에 `Connection string`을 추가한다. requested에서는 Owner/Agent가 비어 있다. Google callback은 Owner를 바인딩하고 prepared로 전환한다. 명시적 동의는 새 Agent를 생성하고 approved로 전환한다. 기존 complete가 키·credential을 저장하고 consumed로 전환한다. 취소·만료·소비 상태는 재사용하지 않는다. 같은 회원의 두 요청은 별도 agent와 키를 갖는다. 관계는 자동 생성하지 않는다. 승인 뒤 완료 응답이 유실되면 미사용 agent가 남을 수 있으므로 홈에서 폐기 후 새 요청을 시작한다. JSON 추가 필드이며 SQL migration은 없다. D07 물리 설계와 D08 생성 테이블 정의는 변경하지 않는다.
