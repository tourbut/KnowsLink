---
id: D09
title: CRUD정의서
status: review
updated: 2026-10-06
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-PUBLIC-IDENTITY-001-DEV, SAR-PUBLIC-AGENTS-001-DEV, SAR-PUBLIC-AGENTS-001-DEV-FIX, SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX, SAR-PUBLIC-MESSAGES-001-DEV, SAR-PUBLIC-MESSAGES-001-DEV-FIX, SAR-PUBLIC-MESSAGES-001-TESTER]
upstream: [D02]
summary: owner와 pair 및 transport와 gate의 상태 변경을 정의한다
---

# KnowsLink CRUD 정의

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
| gate GET/approve/deny | 본문 서명·digest·현재 권한·세대·기한 확인. GET은 결정 없음. POST csrf. 자기 gate deny는 정리 budget. 결정 뒤 정식 gate 화면 |
| 한도·안전 정리 | 포화는 신규 자원 생성 없음. ACK·deny·철회·unpair·로그아웃은 별도 입장/rate. 부모 TTL에서 claim 종료 |
| HTTP 입장 | 본문 수신 뒤 입장. 신원별 request budget 한 번 차감. 공유 HTTP map 입장/종료 삭제·30s crash 정리. 신규16/정리4. 정리는 검증된 자기 기록만 |

동작·실패·경합·재시작 근거는 [실행 기록](../exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV.md)에 연결한다. 운영 자료 삭제는 수행하지 않았다.

독립 QA `09c523da8a3407288d9f5d711e1834af12bc7808`는 위 표의 ACK, 경로 deny, unpair, 로그아웃이 신규 DB 슬롯 포화 중에도 동작함을 확인했다. `POST /home/gates/{id}` form deny는 그 정리 예산에 들어가지 않고 신규 슬롯이 가득하면 429다. 상세는 [QA 보고서](../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER.md)다.
