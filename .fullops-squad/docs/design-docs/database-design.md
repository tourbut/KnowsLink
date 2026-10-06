---
id: D07
title: 데이터베이스설계서
status: review
updated: 2026-10-06
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-PUBLIC-IDENTITY-001-DEV, SAR-PUBLIC-AGENTS-001-DEV, SAR-PUBLIC-AGENTS-001-DEV-FIX]
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

## SAR-PUBLIC-AGENTS-001 저장 변경

Connection map·Key.Credential·Agent.Revoked·Pair.Exp를 같은 JSONB에 추가한다. 실제 table·SQL·sqlc 계약과 migration은 바뀌지 않았다. 따라서 D08 테이블 정의 생성은 해당 없음이다. 기존 상태는 빈 Connection map과 false Revoked로 읽는다. 기존 pending Exp는 첫 sweep에서24h로 보강한다.
grant 완료·key credential 등록·기존 key 철회는 한 transaction이다. 동시 완료는 한 번만 credential을 반환한다. agent200/owner5·키3·active pair400/owner20·pending200/송신owner10은 같은 lock에서 집행한다. pair 수락 실패는 pending을 바꾸지 않는다. 재시작에도 한도·grant 사용·철회·rate를 유지한다.
이전 코드로 rollback하면 추가 JSON 필드를 버릴 수 있다. 특히 key별 credential·agent 철회 필드의 손실은 보안 rollback이 아니다. 공개 전에 OPS는 해당 SHA의 전체 스키마를 보존하고 철회 상태가 불명확하면 재연결·재수락 전 차단해야 한다. 레코드 보존량과 전체 처리량 검증은 공개 수락 후속이다.

## SAR-PUBLIC-AGENTS-001-DEV-FIX 저장 변경

같은 JSONB에 Agent.Changed를 추가한다. table·SQL·sqlc·migration은 바뀌지 않았다. D08 생성은 해당 없음이다. 기존 철회 agent는 첫 sweep에서 Changed를 받고 그 시각부터 24h를 보존한다.
sweep은 모든 transaction 시작과 1초 cleanup loop에서 실행한다. 철회 24h가 지난 agent와 그 Key·Pair를 같은 transaction에서 삭제한다. 거부된 요청의 transaction도 sweep 결과와 사용한 rate 기록만 저장한다. 부분 업무 상태는 저장하지 않는다.
상태 크기의 이론 상한은 회원 100명 × agent 기록 10개 × 키 기록 20개 = 키 20000개다. 이 값은 코드 상한에서 계산한 최대치다. 실제 JSONB 크기·CPU·처리량·복원 측정은 OPS 공개 수락 후속이며 이번에 측정하지 않았다.
이 SHA에서 이전 SHA로 rollback하면 Agent.Changed를 잃는다. 다시 이 SHA로 올리면 기존 철회 agent는 그 시점부터 24h를 다시 보존한다. 보존 기간이 짧아지지 않는다.
