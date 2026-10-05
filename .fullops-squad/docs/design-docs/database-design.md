---
id: D07
title: 데이터베이스설계서
status: review
updated: 2026-10-05
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-PUBLIC-IDENTITY-001-DEV]
upstream: [D02]
summary: singleton Postgres 상태와 epoch CAS 및 처리량 한계를 정의한다
---

# KnowsLink DB 설계

[엔티티 정본](data-model.md)을 따른다.



실제 schema는 [생성 테이블 정의](../generated/db-schema.md)다.
`relay_state(singleton,epoch,clock,data)`의 boolean PK·CHECK는 singleton 행을 보장한다.
SQL-only migration 00001은 업무 행과 goose version table을 만든다.
sqlc의 LockRelay는 FOR UPDATE와 DB 현재 시각을 반환한다. SaveRelay는 epoch CAS의 영향 행 수를 반환한다.
성공한 CAS만 commit한다. 권한 작업과 상태 작업은 같은 lock을 사용한다.
JSON map 키는 transaction 내 uniqueness를 구현한다. 서로 다른 relay process도 공유 lock으로 경합을 직렬화한다.

단일 행의 전체 상태 재직렬화는 로컬 합성 규모의 의도적 한계다. 공개 규모에서 성능·한도 수락 없이 사용하지 않는다.
키 revoke metadata는 24h보다 오래 보관한다. kid 재할당을 영구 금지한다.
receipt·멱등·gate metadata는 24h 뒤 삭제한다. Envelope와 Inbox는 만료·철회·응답 완료 후 먼저 지운다.
1초 background 정리는 유휴 payload를 제거한다. 읽기·인가에서는 만료 원문을 즉시 사용할 수 없다.
행 삭제는 Postgres WAL/backup의 완전 삭제가 아니다.

SAR-PUBLIC-IDENTITY-001은 새 migration 없이 같은 JSON에 Member·Identity·Session·Challenge·Rate map을 추가했다. 이전 상태에는 이 key가 없으므로 빈 map으로 읽는다. 따라서 rollback은 코드 복귀만으로 가능하다. 단, 이전 코드는 새 key를 버리고 저장한다. 회원·세션은 사라진다.
세션 무활동 갱신은 회원 화면 요청마다 행을 다시 쓴다. 회원 100명 규모의 처리량은 미측정이다. 정규화와 부하 측정은 공개 수락 전 후속 과제다.
