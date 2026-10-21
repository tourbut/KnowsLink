---
id: D09
title: CRUD정의서
status: review
updated: 2026-10-21
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-PUBLIC-IDENTITY-001-DEV, SAR-PUBLIC-AGENTS-001-DEV, SAR-PUBLIC-AGENTS-001-DEV-FIX, SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX, SAR-PUBLIC-MESSAGES-001-DEV, SAR-PUBLIC-MESSAGES-001-DEV-FIX, SAR-PUBLIC-MESSAGES-001-TESTER, SAR-PUBLIC-MESSAGES-001-DEV-FIX-2, SAR-PUBLIC-MESSAGES-001-DEV-FIX-3, SAR-GOOGLE-CONNECT-001-DEV, SAR-SSE-INSTALL-001-DEV, SAR-HUMAN-PEER-001-DEV, SAR-REFACTOR-002-DEV, SAR-SCHEDULE-001-DEV, SAR-PUBLISH-001-DEV]
upstream: [D02]
summary: owner와 pair 및 transport와 gate의 상태 변경을 정의한다
---

# KnowsLink CRUD 정의

## SAR-REFACTOR-002 처리 순서 보완

로컬 Inbox는 unread 본문 claim 전에 housekeeping 오류를 확인한다. 다른 receiver가 이미 없앤 임시 파일·marker의 stat ENOENT만 건너뛴다. 실패한 receive가 unread 본문을 지우지 않으며 성공한 read marker·transport ACK 의미는 그대로다.

Remote callback의 마지막 lease 만료 후 기존 event의 3회 예산을 종료한다. 그 사이 다른 pending이 생겼으면 즉시 다음 loop를 예약하고 새로운 event·예산을 만든다. 활성 lease를 건너뛰지 않으며 같은 event를 네 번째로 발송하지 않는다. 현재 권한·TTL·명시적 Inbox ACK와 구독 상태 저장 형식은 유지한다.

## SAR-HUMAN-PEER-001 DEV 후보

| 동작 | 원자 변경·거부 |
|---|---|
| 최초 새 writer transaction | HumanModel1, pending legacy 종료, 기존 회원 SelfConnect=false/코드 생성; active pair·키·ID 보존 |
| 코드 요청 | 명시 동의/이름/현재 신원 검사. 같은 또는 반대 pending 재제출은 세대·기한·수신자 불변. 자동 수락 없음 |
| accept/deny/cancel/expire | 정확한 참여자와 generation 검사. 수락은 양측 활성 이름·현재 한도 검사. 실패는 pending 유지 |
| remove/block/unblock | remove/block은 SuppressLegacy=true, current 전달·claim/gate 즉시 차단. unblock은 친구 복원 없음 |
| 재요청·재수락 | 새 generation. 과거 메시지/lease/claim/gate/idempotency는 새 권한을 얻지 못함 |
| 웹 owner text / signed agent root text | 현재 출처·세대·키·신원 검사, Message+멱등+Conversation 원자 저장. 각 발송 독립180s |
| 대화 기록 GET | 유효 세션 owner가 사본의 참여자인지 검사. 당시 이름/본문만24h. 현재 친구 정보·키·발송 권한을 재생성하지 않음 |
| ACK/TTL/철회/24h | transport Envelope/Inbox는 기존 즉시삭제; Conversation은24h까지, 이후 Message+멱등과 함께 sweep |

friend deny/cancel/remove/block/unblock은 기존 검증된 자기 정리 입장 예산을 사용한다. 로그아웃·키/agent 철회·gate deny·ACK의 기존 안전 정리는 유지한다. 일반회원 legacy invite 신규 생성은 거부하며 테스트 synthetic API의 명시적 fixture 예외는 production에서 닫혀 있다.

## SAR-SSE-INSTALL-001 SSE 입장

현재 agent Bearer가 SSE에 입장할 때 공유 Streams 기록을 만든다. 기존 HTTP/rate 입장은 handshake 뒤 반납한다. LISTEN 확정 후와 NOTIFY 때만 현재 상태를 읽으며 heartbeat는 DB를 조회하지 않는다. 정상 종료는 자기 기록만 삭제하고 crash는65초 만료 sweep으로 회수한다. 신호는 Message를 lease·persist·ACK하거나 실행 권한을 만들지 않는다. 본문은 기존 text API와 private Inbox 저장 순서를 유지한다.

[엔티티 정본](data-model.md)과 [DB 설계](database-design.md)를 따른다.



| 엔티티 | Create | Read | Update | Delete·회수 |
|---|---|---|---|---|
| Owner/Agent/Key | owner 가입·PoP 등록 | 현재 owner/agent 인증, key lookup | rotate/revoke | ID/kid 이력 유지 |
| Pair | invite pending | contacts·routing | B-owner accept/deny, owner unpair | 재수락은 generation 증가 |
| Message/Receipt/Idempotency | send의 동일 transaction | endpoint 권한·receipt-only | `deliver:agent`만 lease·persist·ACK·claim·authorize·gate-consume·H/R 부모; H는 owner 결정; R completion | payload exp/완료/철회; metadata 24h |
| Member/Identity | 첫 코드 확인 성공 | 세션 검증·홈 | 없음(이번 범위) | 비활성화·30일 정리는 후속 |
| Session | 코드 확인 성공 | 회원 요청마다 검증 | Seen 갱신 | 로그아웃·전체 로그아웃·절대/무활동 만료·재확인 시 교체 |
| Challenge | `/auth/start`·`/auth/reauth` | `/auth/verify` | 오답 Attempts 증가 | 성공·5회 오답·10분 만료·새 요청·발송 실패 |
| Rate | 요청·발송 | 한도 계산 | 시각 추가 | 1h 지난 시각 정리 |
| Gate | 검증된 H send | owner UI, B-agent metadata | 결정 CAS·consume 1회 | 만료/철회/원문 부재 차단; metadata 24h |

MVP-15의 최소 저장은 HTTP 응답·inbox·로그·adapter memory에도 적용한다.
프로그램 종료 시 claim을 다른 worker에 재발급하지 않는다. 자동 복구보다 중복 처리 차단을 우선한다.
실제 외부 효과와 positive silent done 재개는 DEC-02와 후속 업무 과제의 책임이다.


### Remote OAuth·Cloud 설치

Remote consent만 새 설치 agent/key와 일회용 code를 생성한다. token exchange는 code를 소모하고 refresh 시 기존 token pair를 회전한다. remote receive는 기존 lease→서명 검사→persist→private server Inbox 기록→ACK를 한 transaction에서 commit하며 응답 유실에도 Inbox가 남는다. explicit remote ack는 Inbox만 삭제한다. subscription refresh는 동일 설치·이벤트·URL identity를 갱신하고 replacement secret의 짧은 이중서명을 허용한다. 이벤트는 권한을 다시 검사한 뒤 atomic lease로 발송하며 안정 ID의 제한 재시도를 수행한다.

## SAR-PUBLIC-AGENTS-001 상태 변경

| 대상 | 생성·읽기 | 변경·회수 |
|---|---|---|
| Agent | 회원 세션으로 무작위 ID·미연결 상태 생성 | 자기 최근 재인증으로 전체 철회; 활성 slot 회수 |
| Connection | 자기 세션으로10분 grant, 유효 token으로 info/prepare | 지문 확인 후 approved; client PoP complete 한 번만 키 생성; cancel/expire는 미연결; 만료 뒤24h 삭제 |
| Key | prepared 새 공개키·PoP, owner 확인 뒤 key별 credential | register 활성 최대3, rotate 기존 키 전체 철회, selected revoke; kid 재사용 금지 |
| Pair | agent ID 초대 pending; owner 홈의 발신/수신 상태 | 수신 owner accept/deny; 양측 unpair; UI는 Generation 검사; 24h pending expire; 재초대 새 세대 |
| Rate | 회원·agent stable principal의 성공/거부 요청 | 신규와 정리 분리; window 뒤 회수; 재시작 유지 |

agent 자격은 owner·관계 수락·gate 승인 권한을 만들지 않는다. 기존 owner 상태 기계의 operateAs를 회원 세션 경로가 재사용한다. 이메일/회원 디렉터리 공개 조회는 없다. fixture와 외부 계정 설치 성공을 구분한다.

### SAR-PUBLIC-AGENTS-001-DEV-FIX 보존·삭제

| 대상 | 생성 거부 | 삭제 |
|---|---|---|
| Agent | owner의 agent 기록(활성+철회) 10개면 create·합성 agents 409 | 철회 시각부터 24h 뒤 sweep이 삭제 |
| Key | agent의 키 기록 20개면 connect 발급·complete·합성 keys 409 | 살아 있는 agent에서는 삭제하지 않음(kid 재할당 금지). agent 삭제 때 함께 삭제 |
| Pair | 기존 pair 한도 유지 | 삭제되는 agent를 포함한 pair만 같은 sweep에서 삭제. invite-decision·unpair는 없는 agent를 가리키는 pair를 403으로 거부 |

철회·취소·거절·unpair는 기록을 늘리지 않는다. 상한 포화 중에도 정리 budget으로 처리한다.

### SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX 관계 반복과 재초대

| 동작 | 상태 변경 |
|---|---|
| 유효 pending·active에 같은/반대 방향 초대 | 없음. 기존 Pair를 반환한다(수·세대·기한·수신자 불변). 회원 rate는 소비한다 |
| 거절·만료·철회 Pair에 수동 초대 | 한도·현재 agent/owner 검사 뒤 같은 Pair 키에 `Generation+1`, 새 24h 기한의 pending. 초대자는 새 제출자 |
| 상한 초과 재초대 | 없음(409). 세대 증가 없음 |
| 종료된 초대의 늦은 결정 | 없음(403). 회원 화면은 Generation 일치도 검사한다 |

기록 보호 포화의 생성 거부·철회 허용·24h 삭제 규칙은 DEV-FIX 표와 같다. 이번 과제는 거부 안내만 바꿨다.

## SAR-PUBLIC-MESSAGES-001 — PS08–11 CRUD

| 동작 | 조회·변경·삭제 |
|---|---|
| 일반 text 송신 | 현재 양측 Member/Owner/Agent/Key/Pair 확인. 새 Message+Idempotency를 원자 생성. 동일 key/digest는 receipt만 조회 |
| 관련 답장 | 현재 전달된 원요청·세대·기한 확인. 답장 Message 생성과 Parent.ReplyID를 원자 확정. 두 번째 답장은 거부 |
| text receive | text 전용 lease·durable persist·ACK. ACK에 Envelope/Inbox 삭제, completion received. 관련 부모는 reply_received |
| 회원 receipt | 세션·선택 agent 소유권·endpoint 확인. metadata만 조회. 타 회원/임의 agent는 거부 |
| 세션 수신 확인 | 수신 agent만 delivered·현재 세대 text의 Message.SessionSeen을 처음 한 번 기록(U). 발신 receipt는 `session_received_at`으로 조회(R). 기존 relay 상태 JSON 필드이며 DB 스키마 변경 없음 |
| 친구 agent 표시 이름 | 친구·자기 연결·active pair 권한이 있는 agent만 상대 Agent.Name·Client와 Member.Name을 조회(R). 변경 없음 |
| gate GET/approve/deny | 본문 서명·digest·현재 권한·세대·기한 확인. GET은 결정 없음. POST csrf. 자기 gate deny는 정리 budget. 결정 뒤 정식 gate 화면 |
| 한도·안전 정리 | 포화는 신규 자원 생성 없음. ACK·deny·철회·unpair·로그아웃은 별도 입장/rate. 부모 TTL에서 claim 종료 |
| HTTP 입장 | 본문 수신 뒤 입장. 신원별 request budget 한 번 차감. 공유 HTTP map 입장/종료 삭제·30s crash 정리. 신규16/정리4. 정리는 검증된 자기 기록만, owner 제어·agent ACK 단위별 동시 1개(로컬·공유). persist 전 ACK는 신규. 신규 포화 중 snapshot이 증명 못한 정리 자격은 lock 없는 1회 재조회(DEV-FIX-3). 종료 실패 기록은 다음 commit이 삭제 |

동작·실패·경합·재시작 근거는 [실행 기록](../exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV.md)에 연결한다. 운영 자료 삭제는 수행하지 않았다.

독립 QA `09c523da8a3407288d9f5d711e1834af12bc7808`는 위 표의 ACK, 경로 deny, unpair, 로그아웃이 신규 DB 슬롯 포화 중에도 동작함을 확인했다. `POST /home/gates/{id}` form deny는 그 정리 예산에 들어가지 않고 신규 슬롯이 가득하면 429다. 상세는 [QA 보고서](../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER.md)다.

## Google 연결 CRUD — SAR-GOOGLE-CONNECT-001-DEV

start는 Device 요청만 생성한다. Google callback은 요청의 Owner와 prepared 상태를 갱신한다. 최근 인증된 동일 회원의 device-confirm만 독립 agent를 생성한다. 서명된 poll은 승인 전 회원·agent를 반환하지 않는다. 기존 complete는 키·credential을 생성하고 요청을 소비한다. 기존 cancel·revoke·cleanup을 재사용하며 관계 요청/수락은 변경하지 않는다. 타 회원 재바인딩, 타 키 조회, 만료·재실행은 거부한다. 상세 계약은 D05와 D06을 따른다.

## 캘린더 연결 정리 CRUD — SAR-PUBLISH-001-DEV

D018로 아래 '캘린더 연결 CRUD' 절의 연결 동의·승인 조회·결과 공개·연결 해제·lookup sweep을 모두 제거했다. 남은 캘린더 CRUD는 relay 기동 때 한 번 실행하는 정리 삭제뿐이다.

| 동작 | 생성 | 갱신 | 삭제 |
|---|---|---|---|
| 기동 정리(`RetireCalendars`) | — | lookup이 있던 미consume `pending`·`approved` gate → `denied`, 모든 `Gate.Lookup` → 빈 값 | `State.Calendars` 전체, calendar 표시 진행 중 Google attempt. commit 뒤 키가 유효하면 토큰마다 provider revoke 요청 1회 |
| 그 뒤 기동 | — | — | 대상 없음 |

revoke 실패는 삭제를 되돌리지 않는다. consume된 gate는 이미 결과를 넘겼으므로 `approved`로 남는다. 정책 페이지(`/privacy`, `/terms`)는 읽기 전용 정적 본문이며 CRUD 대상 데이터가 없다.

## 캘린더 연결 CRUD — SAR-SCHEDULE-001-DEV

> 비현행: 이 절은 D018 철회 전 기록이다. 현행은 위 '캘린더 연결 정리 CRUD — SAR-PUBLISH-001-DEV' 절이다.

| 동작 | 생성 | 갱신 | 삭제 |
|---|---|---|---|
| 연결 동의 callback | `Calendars[member]`(봉인 토큰, connected) | 기존 연결 덮어쓰기 | — |
| 승인(approve, schedule.query) | `Gate.Lookup{Started}` | lookup 결과(done+disclosure 또는 failed), 거부 시 `Calendars[member].State=reconnect` | — |
| 결과 수락(`relay.result done`) | R 메시지 | 부모 completion=done | `Gate.Lookup.Disclosure` |
| 연결 해제 | — | 대기·미조회 승인 gate → denied | `Calendars[member]`, provider revoke 요청 |
| sweep | — | lease 1분 초과 lookup → failed | approved가 아니거나 부모 만료·완료된 gate의 disclosure |
