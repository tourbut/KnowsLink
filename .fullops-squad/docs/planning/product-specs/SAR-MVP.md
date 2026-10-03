---
id: D02
title: KnowsLink MVP 요구사항
status: review
updated: 2026-10-03
owner: designer
tasks: [SAR-PREP-002]
upstream: [D01]
downstream: [D03, D05, D06, D07, D09, D10]
summary: 최신 원천의 전체 MVP 규칙과 사용자 수락 조건 및 후속 검증 책임을 정의한다
---

# KnowsLink MVP 요구사항

## 정본과 범위

이 문서는 전체 MVP의 제품 요구사항이다. 원천은 service-design `7bc9ea190ea549fae8b047e850247a19322fc9c3`이며 [source.json](../sources/silent-agent-relay/source.json)에 고정돼 있다. [제품 결정](../sources/silent-agent-relay/product.md), [frozen 프로토콜](../sources/silent-agent-relay/protocol.md), [결정 로그](../sources/silent-agent-relay/decisions.md), [MVP 체크리스트](../sources/silent-agent-relay/mvp-checklist.md)를 따른다.

[SAR-SETUP-001 D02](SAR-SETUP-001.md)는 초기 구성과 lint의 이력이다. 초기 골격 수락을 전체 MVP 수락으로 소급하지 않는다. 업무 동작은 아직 미구현이다. 이 문서의 상태는 review다. 원천의 잠긴 결정과 이번 기획 정리는 별도 제품 승인 근거를 혼동하지 않는다.

제품 기획 담당은 designer다. 기술 계획·구조/API·구현·관련 회귀·기술 문서 갱신은 dev가 같은 기능 과제에서 맡는다. 적용 기준은 `fullops-common-0.3.2`, [project.md](../../../project.md), [문서 규칙](../../agents/document-writing.md)이다. 준비 과제 기준 ref는 `0dd08ec994771836c15d9d22a6a83393a71d7987`이다.

## 포함과 제외

포함 범위는 agent 가입·owner 바인딩·pubkey·rotate/revoke, invite·accept/deny·unpair, agent contacts, 서명·검증·인증·인가·ingest, 짧은 큐·멱등·receipt, pull lease·ACK·공유 실행 claim, 결과 반환, human-gate approve/deny다. 첫 Go 기능 구현은 human-gate까지 연결한다.

Go relay와 TypeScript pull-default 어댑터를 사용한다. human-gate UI는 Go `net/http` + `html/template`로 relay가 직접 서빙한다. 별도 TypeScript frontend는 두지 않는다. Postgres와 별도 SQL-only `cmd/migrate`, `pgx/v5` + `pgxpool`, goose, sqlc는 원천의 잠긴 선택이다. 기술 파일·함수·API·테이블 설계는 이 문서에서 정하지 않는다.

제외 범위는 장기 채팅 저장, 병원·폐쇄망, 벤더 코어 수정, 공식 cross-agent inbound API 가정, `schedule.commit` 실행 가능화, amend, evidence 자동 fetch/preview, optional webhook, latent KV/token-id handoff, 결제·구독 구현이다. Workers/DO는 현재 MVP 호스팅이 아니다. 운영 배포·Tunnel 연결·외부 발송은 별도 승인 전 실행하지 않는다.

## 사용자 흐름

1. 에이전트가 owner에 연결되고 공개키를 등록한다. owner 인증과 새 키 PoP를 확인한다.
2. A-agent 또는 A-owner가 B-agent를 초대한다. B-human이 처음 관계를 수락한다. 수락 전 요청을 전달하지 않는다.
3. A는 agent contacts에서 B-agent를 선택한다. A-agent는 typed body의 `relay.v1` 요청 M을 서명해 보낸다.
4. relay는 현재 권한을 확인하고 M과 receipt를 원자적으로 수락한다. B-agent는 pull lease를 받는다.
5. B-agent는 durable inbox에 저장하고 ACK 성공을 확인한다. 공유 실행 claim을 얻은 뒤 요청을 판단한다.
6. permission·judgment·risk가 모두 허용된 요청만 silent 처리한다. B-human에게 통지하지 않는다. 허용 정책이 없으면 정보 공개를 거부한다.
7. 인간 권한이나 판단이 필요하면 B-agent가 새 H=`relay.approval.request`를 보낸다. owner는 검증된 M과 정책을 보고 approve/deny를 결정한다.
8. B-agent는 원요청에 결속된 새 `relay.result` R을 보낸다. A는 transport 상태와 `done|denied|failed` 처리 결과를 구분한다.

approve는 해당 요청의 인간 게이트 통과만 뜻한다. 일정 공개 정책이나 실행 불가능한 stub을 활성화하지 않는다. 실데이터 silent 성공 흐름은 미정 공개 정책이 확정되고 C1–C5가 구현·검증된 뒤 수락한다.

## 요구사항과 사용자 완료 조건

다음 ID는 백로그·DEV 검사·독립 QA에서 공통으로 사용한다. 근거 P는 [product.md](../sources/silent-agent-relay/product.md), W는 [protocol.md](../sources/silent-agent-relay/protocol.md), D는 [decisions.md](../sources/silent-agent-relay/decisions.md)다.

| ID | 제품 규칙·근거 | 사용자 완료 조건·검증 증거 |
|---|---|---|
| MVP-01 | agent signup·owner·키·PoP, C1; P 등록 모델, W 서명/C1 | owner가 가입·rotate/revoke를 수행한다. agent credential로 accept/approve/revoke할 수 없다. `(from,kid)`만 키를 찾으며 URL/path kid와 재할당을 거부한다. |
| MVP-02 | pairing·contacts, C1/C5; P Pairing, W C5 | B-human 수락 전 전달이 거부된다. pending invite는 active slot을 소모하지 않는다. 동시 accept는 중복 관계를 만들지 않는다. auto 재초대는 현재 active pair에만 허용한다. unpair 재수락은 새 세대다. |
| MVP-03 | frozen `relay.v1`, C5; W 봉투/서명 | strict JSON의 중복 키·잘못된 Unicode·알 수 없는 필드·잘못된 서명을 거부한다. verify와 exec는 같은 파싱 객체를 쓴다. 등록 AgentID의 lowercase ASCII 일치를 검증한다. `ext`로 권한을 확대하지 않는다. |
| MVP-04 | auth-first ingest·멱등, C3; W Ingest/멱등 | structure/signature/auth/routing, digest, atomic idempotency, exp/TTL/id, queue+receipt 순서를 지킨다. 인증 전 cache probe가 불가능하다. 동일 `(from,key)`+digest는 receipt만 반환한다. 다른 digest는 `409 idempotency_conflict`다. TTL/id 실패는 reservation을 rollback한다. |
| MVP-05 | 짧은 TTL·receipt; W 운영 상수/상태 | `MAX_TTL=300s`, lease `30s`, 성공 lease grant만 attempts 증가, attempts `3`, receipt·멱등 `24h`, revoke metadata `≥24h`를 지킨다. `lease_until=min(exp,lease_granted_at+30s)`다. 마지막 lease도 유효 window에서 쓸 수 있다. 만료는 `failed:expired`, 소진은 `failed:max_attempts`로 알려준다. replay로 TTL·실행권을 연장하지 않는다. |
| MVP-06 | durable persist·ACK·공유 claim, C3; W 상태 | lease 자체를 delivered로 표시하지 않는다. transport delivered는 ACK 성공 이후다. processing은 durable persist, ACK 성공 확인, 공유 실행 claim 이후다. lease token은 recipient/msg/generation에 결속된다. 늦은 ACK·다중 adapter·재시작으로 중복 실행권을 얻지 못한다. 외부 도구 exactly-once는 주장하지 않는다. |
| MVP-07 | current-auth, C1; W 서명/C1 | enqueue·lease·ACK·approval consume·exec·result 공개에서 현재 key/pair/owner를 재검사한다. epoch CAS로 확정한다. 철회가 인가 확정 전에 반영되면 거부한다. 이전 receipt나 pair 세대로 복구하지 않는다. 권한 불명·비정상 시계는 차단한다. |
| MVP-08 | deny-by-default·seed; P 판단, W registry | payment/delete/grants/outbound sends/commitments는 인간 게이트가 필수다. `schedule.query`는 정책 없으면 deny다. `schedule.commit`은 stub이다. registry는 central·ADD-only며 subset 지원이다. intent 추가가 자동 실행 등록을 뜻하지 않는다. |
| MVP-09 | escalation H, C2; W Escalate | H는 새 서명 메시지다. `deliver=human`, `reply_to=M.id`, `from=to=M.to`, `request_digest=semantic_digest(M)`, `H.exp≤M.exp`다. 권위 있는 부모·B 실제 수신·현재 부모 pair와 B-owner accept를 확인한다. 부모별 pending gate 중복과 result/approval 재귀 escalate를 막는다. |
| MVP-10 | approve/deny, C2; W Escalate/C2 | UI는 검증된 M typed body+정책을 표시한다. 원문 없으면 승인 불가다. owner 인증·CSRF 방어를 적용한다. GET/링크만으로 승인하지 않는다. 봉투 밖 기록은 owner, M.id, digest, endpoints, pair 세대, 정책, 만료, 결정·소비를 결속한다. 결정·소비는 원자적이며 수명은 `≤min(M.exp,H.exp)`다. 같은 digest의 다른 M.id에는 재사용하지 않는다. |
| MVP-11 | result binding·최소화, C4; W Agent completion/C4 | `R.reply_to=M.id`, `R.from=M.to`, `R.to=M.from`과 현재 공개 권한·부모 intent별 output allowlist를 확인한다. 허용되지 않은 일정 제목·참석자·위치·원본 객체·stack trace를 반환하지 않는다. 허용 schema 확정 전 optional result/error 데이터를 열지 않는다. 결과를 도구 명령으로 실행하지 않는다. transport 실패를 B-서명 결과로 위조하지 않는다. |
| MVP-12 | OFF 경계·자원, C5; W Nested/C5 | evidence ref 수신이 자동 fetch/preview를 일으키지 않는다. webhook은 OFF다. high priority도 제한을 우회하지 않는다. 추가 rate/size/concurrency 수치는 TBD며 무제한 운영 배포는 차단한다. body JCS `16KiB`, intent `64bytes`, hint `1024 UTF-8 bytes`, evidence `8`, key `16–128 ASCII`의 기존 한도는 유지한다. |
| MVP-13 | 어댑터; P 범위, D stack | TypeScript pull stub 1개로 수신·persist·ACK·claim·gate·결과 흐름을 로컬 검증한다. 실제 대상은 Grok Bot, Claude Code, Codex, Dots 순서다. 공식 inbound API·벤더 코어 패치를 전제로 하지 않는다. 실제 제품 연결을 stub 성공으로 보고하지 않는다. |
| MVP-14 | A2A 잠금; D A2A gap review | 검토 완료 범위는 공개 A2A v0.3.0 개념이다. 최신 개정판 delta를 재검토했다고 하지 않는다. wire 비호환이며 `taskId/contextId/parts/artifacts/A2A state enums`를 추가하지 않는다. `role:user`·AgentCard·push는 owner 승인이 아니다. push로 webhook을 활성화하지 않는다. |
| MVP-15 | 저장·기밀 최소화; W C4, architecture 저장 | exp 또는 전달·응답 완료 뒤 원문을 클리어한다. receipt 24h는 원문 24h 저장이 아니다. inbox/log/trace/model context를 최소화한다. digest를 익명화 데이터로 간주하지 않는다. DB 행 삭제를 WAL/backup 완전 삭제로 표시하지 않는다. |
| MVP-16 | UI 사용자 판정; P Human-gate UI | owner가 발신·대상 에이전트, intent·typed body, 적용 정책, 만료·유효 상태, approve/deny 결과를 확인한다. pending/approved/denied/expired/revoked/unavailable 상태를 혼동하지 않는다. hint와 실제 요청이 달라도 실제 검증 본문이 판단 근거다. |

## 전달 상태 용어

원천 README/product의 delivered 설명은 persist·ACK·claim까지의 처리 선행 조건을 축약한다. 세부 상태 정본은 W의 상태 ownership이다. transport `delivered`는 ACK 성공 이후이고, 실행은 공유 claim 이후다. `done`은 새 `relay.result`의 처리 결과다. 이 구분은 기술 계약 변경이 아니다.

멱등 digest에서 `id/exp/sig/trace/render`는 제외된다. 포함 필드와 서명 prefix·RFC8785 JCS·Ed25519 규칙은 W를 그대로 따른다. natural-language mid-hop, canonical `body.text`, `deliver:both`, top-level parent-link, 요청 봉투 `status`, `kind/body.op/body.args`를 추가하지 않는다.

## UI 제품 방향

최소 owner 작업 화면과 approve/deny 화면을 사용한다. 메시지 대화창·채팅 버블·입력 composer·장기 대화 타임라인을 만들지 않는다. render.hint를 요청 제목으로만 크게 보여주고 검증 본문을 숨기지 않는다. 승인 버튼을 GET 링크로 만들지 않는다. 색만으로 승인·거절·만료를 구분하지 않는다. 만료나 권한 확인 실패 상태에서 활성 승인 버튼을 보이지 않는다.

별도 아트 에셋이나 장식 목업은 현재 필요하지 않다. D04는 미작성 상태를 유지한다. DEV의 UI 후보가 준비되면 designer가 MVP-16을 직접 시각 검수한다. 정상 승인 대기, 결정 완료, 만료·철회·원문 부재 화면을 캡처한다. 시간 변화가 정지 화면으로 판정되지 않을 때만 짧은 영상을 만든다.

## 미정 결정과 기능 수락 경계

[백로그의 결정 보류표](../SAR-MVP-backlog.md#미정-결정과-재개-조건)가 담당·영향·재개 조건의 정본이다. Free N·가격·slot-unit, disclosure·결과 schema·window/granularity/누적 한도, 추가 resource/rate/size/concurrency 수치를 만들지 않는다. 미니서버 설정과 어댑터 인터페이스는 담당 dev/ops의 기술 확인이다.

첫 기능은 합성 요청의 로컬 안전 전달과 human-gate다. 실제 calendar 조회·외부 부작용·실데이터 silent 성공은 수락 범위에 포함하지 않는다. `schedule.query`는 정책 없음의 거부를 검증한다. 승인과 정보 공개는 별개다. `schedule.commit`은 승인 후에도 실행하지 않는다. positive silent `done` 수락은 공개 정책과 output schema가 확정된 후속 기능에 남긴다.

## 검증 책임과 추적

[백로그](../SAR-MVP-backlog.md)는 요구사항별 기능·담당·선행 조건을 연결한다. [DEV 준비 지시서](../../../handovers/to_dev.md)와 [TESTER 준비 지시서](../../../handovers/to_tester.md)는 queued다. 지금 dispatch하지 않는다.

DEV는 같은 과제에서 기술 계획, 실제 기능 구현, 변경 동작 자동 검사·관련 회귀, 필요한 짧은 로컬 확인, D03과 필요한 D05–D10 갱신을 완료한다. tester는 DEV 완료 뒤 고정된 안정 통합 후보에서 독립 QA를 수행한다. designer는 같은 후보 UI를 직접 검수한다. 구현자와 다른 세션의 독립 fixed-SHA 코드 리뷰는 coor가 준비한다. 미해결 critical/high는 수락·병합을 차단한다.

제품 최종 수락은 DEV 완료, 독립 QA, 직접 UI 검수, 독립 코드 리뷰, 정책 보류 해결을 확인한 뒤 coor가 조정한다. 이번 준비 문서 검사는 해당 제품 수락 증거를 대신하지 않는다. 기존 골격의 검증은 원래 SHA·조건과 의존성 동일성을 확인해 재사용한다.

## 개정 이력

- 2026-10-03: SAR-PREP-002에서 새 MVP D02를 작성했다. A2A v0.3.0 잠금, C1–C5, 기능 수락 ID와 미정 결정의 경계를 연결했다. 기존 초기 구성 D02는 보존했다.
