---
id: D09
title: CRUD정의서
status: review
updated: 2026-10-05
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-PUBLIC-IDENTITY-001-DEV]
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
