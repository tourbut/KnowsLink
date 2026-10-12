---
id: D07
title: 데이터베이스설계서
status: review
updated: 2026-10-12
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-PUBLIC-IDENTITY-001-DEV, SAR-PUBLIC-AGENTS-001-DEV, SAR-PUBLIC-AGENTS-001-DEV-FIX, SAR-PUBLIC-MESSAGES-001-DEV, SAR-PUBLIC-MESSAGES-001-DEV-FIX-2, SAR-PUBLIC-MESSAGES-001-DEV-FIX-3, SAR-SSE-INSTALL-001-DEV, SAR-HUMAN-PEER-001-DEV]
upstream: [D02]
summary: singleton Postgres 상태와 epoch CAS 및 처리량 한계를 정의한다
---

# KnowsLink DB 설계

## SAR-HUMAN-PEER-001 DEV 후보

기존 singleton 행·sqlc 질의·epoch CAS·DB clock을 유지한다. 새 실제 migration `00002_human_writer.sql`은 `guard_human_writer()`와 BEFORE UPDATE OF data trigger를 설치한다. HumanModel이1 이상인 상태를 더 낮은 버전/필드 없는 JSON으로 쓰면 SQL 오류로 rollback된다. 서비스도 알 수 없는 버전을 거부한다. [D08](../generated/db-schema.md)은 두 migration 전체 Up SQL/hash를 생성한다.

OPS 이행 순서: 쓰기 중단 → private 최신 backup과 복원 리허설 → migration Up → 새 reader/writer만 실행 → legacy ID/키/철회/정확한 pair 및 새 friend deny 검증 → 트래픽 재개. 서비스 transaction의 최초 HumanModel 변환은 idempotent하다. mixed old writer는 차단되므로 구버전과 rolling 쓰기를 병행하지 않는다.

Down은 의도적으로 오류를 내며 trigger를 제거하지 않는다. 코드만 과거로 되돌리는 방식은 허용하지 않는다. 복구는 writer 중단 상태에서 현재 schema·friend/block/suppression·generation·철회·키를 보존하는 검토된 reader/writer로 진행한다. 과거 backup 전체 복원은 DB trigger로 감지할 수 없으므로 권한이 불명확하면 read-only/전달 차단 상태에서 재인증·재연결·명시적 새 수락으로 복구한다. 이번 DEV는 운영 migration/복원·WAL 완전 삭제를 실행하지 않았다.

새 Conversation도 기존 message20000 보호 상한 안에 저장하고 수락24h 뒤 sweep한다. ACK/TTL/철회는 transport 원문을 지우되 사본은 남은24h 동안 참여 세션에만 공개한다. singleton 성능·backup 보존·실서비스 보존 정책 수락은 OPS 후속이다.

## SAR-SSE-INSTALL-001 SSE 저장

기존 JSONB State에 Streams map만 추가한다. 각 기록은 자기 agent와65초 기한이며 전체16/agent당2를 row lock에서 제한한다. 새 table·SQL migration·sqlc 생성물은 없다. NOTIFY는 빈 payload이고 업무 변경 transaction과 함께 commit한다. LISTEN 연결과 signal용 현재 상태 조회는 행을 수정하지 않는다. 원문 수신은 기존 row lock과 서명/현재 권한 검사를 유지한다. 이전 코드로 복귀하면 단기 stream 기록이 없어지므로 기존 SSE 연결을 먼저 종료한다. 영속 키·관계·철회 기록을 과거 백업으로 덮지 않는다.

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


### Remote OAuth·Cloud 설치

Remote와 Streams는 기존 relay_state JSON의 additive optional field다. 신규 SQL migration·새 DB/계정·운영 schema 변경은 없다. 과거 JSON은 nil map을 허용하고 최초 실제 remote consent가 설치만 추가한다. 기존 owner/agent/pair/message·generation은 유지한다. 새 서버 보관 개인키·callback secret은 DB 접근권한과 백업 비공개를 전제로 하며 Google token을 저장하지 않는다. 구독 상한128·설치당4, 설치 상한500을 둔다. 세부 만료/철회와 재시도는 REMOTE 기록을 따른다.

## SAR-PUBLIC-AGENTS-001 저장 변경

Connection map·Key.Credential·Agent.Revoked·Pair.Exp를 같은 JSONB에 추가한다. 실제 table·SQL·sqlc 계약과 migration은 바뀌지 않았다. 따라서 D08 테이블 정의 생성은 해당 없음이다. 기존 상태는 빈 Connection map과 false Revoked로 읽는다. 기존 pending Exp는 첫 sweep에서24h로 보강한다.
grant 완료·key credential 등록·기존 key 철회는 한 transaction이다. 동시 완료는 한 번만 credential을 반환한다. agent200/owner5·키3·active pair400/owner20·pending200/송신owner10은 같은 lock에서 집행한다. pair 수락 실패는 pending을 바꾸지 않는다. 재시작에도 한도·grant 사용·철회·rate를 유지한다.
이전 코드로 rollback하면 추가 JSON 필드를 버릴 수 있다. 특히 key별 credential·agent 철회 필드의 손실은 보안 rollback이 아니다. 공개 전에 OPS는 해당 SHA의 전체 스키마를 보존하고 철회 상태가 불명확하면 재연결·재수락 전 차단해야 한다. 레코드 보존량과 전체 처리량 검증은 공개 수락 후속이다.

## SAR-PUBLIC-AGENTS-001-DEV-FIX 저장 변경

같은 JSONB에 Agent.Changed를 추가한다. table·SQL·sqlc·migration은 바뀌지 않았다. D08 생성은 해당 없음이다. 기존 철회 agent는 첫 sweep에서 Changed를 받고 그 시각부터 24h를 보존한다.
sweep은 모든 transaction 시작과 1초 cleanup loop에서 실행한다. 철회 24h가 지난 agent와 그 Key·Pair를 같은 transaction에서 삭제한다. 거부된 요청의 transaction도 sweep 결과와 사용한 rate 기록만 저장한다. 부분 업무 상태는 저장하지 않는다.
상태 크기의 이론 상한은 회원 100명 × agent 기록 10개 × 키 기록 20개 = 키 20000개다. 이 값은 코드 상한에서 계산한 최대치다. 실제 JSONB 크기·CPU·처리량·복원 측정은 OPS 공개 수락 후속이며 이번에 측정하지 않았다.
이 SHA에서 이전 SHA로 rollback하면 Agent.Changed를 잃는다. 다시 이 SHA로 올리면 기존 철회 agent는 그 시점부터 24h를 다시 보존한다. 보존 기간이 짧아지지 않는다.

## SAR-PUBLIC-MESSAGES-001 저장과 복구

기존 relay_state JSONB에 Message.ReplyID와 State.HTTP를 추가했다. table·SQL·migration·sqlc 생성물은 바꾸지 않았다. D08 실제 테이블 정의도 변경 없음이다. 원문은 text ACK/만료/철회에 삭제하며 receipt·멱등은24h 보존한다. HTTP 입장 기록은 정상 종료 시 삭제하고30s 만료로 crash를 회수한다. SAR-PUBLIC-MESSAGES-001-DEV-FIX-2는 같은 JSONB 기록에 선택 필드 Owner를 추가했다. table·SQL·migration·sqlc는 바꾸지 않았다. 종료 transaction이 실패한 기록은 같은 프로세스의 다음 commit이 지운다. SAR-PUBLIC-MESSAGES-001-DEV-FIX-3은 읽기 전용 sqlc 질의 `ReadRelay`(`SELECT epoch, clock, data`, lock 없음)를 추가했다. 정리 입장 snapshot 재조회만 사용한다. table·migration·D08은 바꾸지 않았다. Owner 값에 agent ACK 단위 `<owner>/agents`가 추가됐다.
queue100·gate100·receipt20000·claim4의 판단과 수락·claim 확정은 동일 row lock·DB 시계를 사용한다. HTTP 신규16/정리4도 동일 상태에 저장하므로 프로세스 교체·다중 relay로 초기화하지 않는다. 입장 기록이 남은 재시작은 보수적으로 거부한다. 현재1초 sweep·권한 철회·세대·epoch CAS를 유지한다. H/R 포화 시 부모는 기존 TTL에서 failed:expired가 되고 claim을 지운다.
실제 크기·CPU·처리량·백업/WAL 완전 삭제·운영 복원은 이번 로컬 회귀로 보장하지 않는다. OPS 공개 수락에서 측정한다. 이전 SHA로 rollback한 후보는 일반 text 도구를 제공하지 않는다. 권한 상태가 증명되지 않으면 재인증·새 연결·새 관계 수락 전 전달을 열지 않는다.
