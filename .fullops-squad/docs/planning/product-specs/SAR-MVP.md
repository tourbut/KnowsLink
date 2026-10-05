---
id: D02
title: KnowsLink MVP 요구사항
status: review
updated: 2026-10-05
owner: designer
tasks: [SAR-PREP-002, SAR-MVP-PUBLIC-POLICY-001, SAR-PUBLIC-SERVICE-001]
upstream: [D01]
downstream: [D03, D05, D06, D07, D09, D10]
summary: frozen MVP 규칙과 일반 서비스 확장 정본 및 과거 파일럿 경계를 연결한다
---

# KnowsLink MVP 요구사항

## 정본과 범위

이 문서는 전체 MVP의 제품 요구사항이다. 원천은 service-design `7bc9ea190ea549fae8b047e850247a19322fc9c3`이며 [source.json](../sources/silent-agent-relay/source.json)에 고정돼 있다. [제품 결정](../sources/silent-agent-relay/product.md), [frozen 프로토콜](../sources/silent-agent-relay/protocol.md), [결정 로그](../sources/silent-agent-relay/decisions.md), [MVP 체크리스트](../sources/silent-agent-relay/mvp-checklist.md)를 따른다.

[SAR-SETUP-001 D02](SAR-SETUP-001.md)는 초기 구성과 lint의 이력이다. 초기 골격 수락을 전체 MVP 수락으로 소급하지 않는다. 합성 MVP·owner-only 운영과 실제 수동 CLI 왕복 이력은 보존한다. 일반 회원 서비스는 미완료다. 이 문서의 상태는 review다. 원천의 잠긴 결정과 이번 기획 정리는 별도 제품 승인 근거를 혼동하지 않는다.

제품 기획 담당은 designer다. 기술 계획·구조/API·구현·관련 회귀·기술 문서 갱신은 dev가 같은 기능 과제에서 맡는다. 적용 기준은 `fullops-common-0.3.2`, [project.md](../../../project.md), [문서 규칙](../../agents/document-writing.md)이다. 준비 과제 기준 ref는 `0dd08ec994771836c15d9d22a6a83393a71d7987`이다.

## 현재 일반 서비스 확장 정본

2026-10-05 사용자 지시는 일반 이메일 서비스 완성 뒤 OpenAI dot “다닷” 연결이다. 현재 확장 요구와 수락 기준은 [일반 서비스 D02](SAR-PUBLIC-SERVICE.md)다. PS-01–14는 아래 MVP-01–16에 추가된다. 실제 일반 신원·연결 확인 메시지·새 운영 기본값은 이 확장 정본을 따른다. 실일정 공개·결제·임의 업무 외부 발송·유료화의 보류와 frozen wire는 유지한다.

아래 DEC-03 파일럿 절과 수치 제안은 2026-10-03 당시 결정 기록이다. 이번 새 운영 기본값이 해당 공개 한도를 대체한다. 과거 제안이 승인됐다고 소급하지 않는다. 어댑터 우선순위는 일반 서비스 수락 뒤 다닷을 우선 연결하는 사용자 지시로 갱신한다. 이전 QA와 시험 자격은 실제 일반 회원 서비스 수락 증거가 아니다.

## 포함과 제외

포함 범위는 agent 가입·owner 바인딩·pubkey·rotate/revoke, invite·accept/deny·unpair, agent contacts, 서명·검증·인증·인가·ingest, 짧은 큐·멱등·receipt, pull lease·ACK·공유 실행 claim, 결과 반환, human-gate approve/deny다. 첫 Go 기능 구현은 human-gate까지 연결한다.

Go relay와 TypeScript pull-default 어댑터를 사용한다. human-gate UI는 Go `net/http` + `html/template`로 relay가 직접 서빙한다. 별도 TypeScript frontend는 두지 않는다. Postgres와 별도 SQL-only `cmd/migrate`, `pgx/v5` + `pgxpool`, goose, sqlc는 원천의 잠긴 선택이다. 기술 파일·함수·API·테이블 설계는 이 문서에서 정하지 않는다.

제외 범위는 장기 채팅 저장, 병원·폐쇄망, 벤더 코어 수정, 공식 cross-agent inbound API 가정, `schedule.commit` 실행 가능화, amend, evidence 자동 fetch/preview, optional webhook, latent KV/token-id handoff, 결제·구독 구현이다. Workers/DO는 현재 MVP 호스팅이 아니다. 2026-10-03 사용자는 현재 서버 Docker·Cloudflare Tunnel의 첫 파일럿 배포를 승인했다. 공개 hostname은 `link.knowslog.com`이다. 승인 근거는 [공개 기준 결정 기록](../../exec-plans/phases/SAR-MVP-PUBLIC-POLICY-001.md)와 [OPS 지시서](../../../handovers/to_ops.md)다. 제품 공개 조건과 수락 SHA가 충족되면 같은 배포 승인을 다시 요청하지 않는다. 실제 외부 업무 발송·실데이터 연결은 이번 승인에 포함하지 않는다.

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
| MVP-13 | 어댑터; P 범위, D stack | TypeScript pull stub 1개로 수신·persist·ACK·claim·gate·결과 흐름을 로컬 검증한다. 초기 원천 순서는 Grok Bot, Claude Code, Codex, Dots다. 현재 후속 우선순위는 일반 서비스 D02를 따른다. 공식 inbound API·벤더 코어 패치를 전제로 하지 않는다. 실제 제품 연결을 stub 성공으로 보고하지 않는다. |
| MVP-14 | A2A 잠금; D A2A gap review | 검토 완료 범위는 공개 A2A v0.3.0 개념이다. 최신 개정판 delta를 재검토했다고 하지 않는다. wire 비호환이며 `taskId/contextId/parts/artifacts/A2A state enums`를 추가하지 않는다. `role:user`·AgentCard·push는 owner 승인이 아니다. push로 webhook을 활성화하지 않는다. |
| MVP-15 | 저장·기밀 최소화; W C4, architecture 저장 | exp 또는 전달·응답 완료 뒤 원문을 클리어한다. receipt 24h는 원문 24h 저장이 아니다. inbox/log/trace/model context를 최소화한다. digest를 익명화 데이터로 간주하지 않는다. DB 행 삭제를 WAL/backup 완전 삭제로 표시하지 않는다. |
| MVP-16 | UI 사용자 판정; P Human-gate UI | owner가 발신·대상 에이전트, intent·typed body, 적용 정책, 만료·유효 상태, approve/deny 결과를 확인한다. pending/approved/denied/expired/revoked/unavailable 상태를 혼동하지 않는다. hint와 실제 요청이 달라도 실제 검증 본문이 판단 근거다. |

## 전달 상태 용어

원천 README/product의 delivered 설명은 persist·ACK·claim까지의 처리 선행 조건을 축약한다. 세부 상태 정본은 W의 상태 ownership이다. transport `delivered`는 ACK 성공 이후이고, 실행은 공유 claim 이후다. `done`은 새 `relay.result`의 처리 결과다. 이 구분은 기술 계약 변경이 아니다.

멱등 digest에서 `id/exp/sig/trace/render`는 제외된다. 포함 필드와 서명 prefix·RFC8785 JCS·Ed25519 규칙은 W를 그대로 따른다. natural-language mid-hop, canonical `body.text`, `deliver:both`, top-level parent-link, 요청 봉투 `status`, `kind/body.op/body.args`를 추가하지 않는다.

## UI 제품 방향

최소 owner 작업 화면과 approve/deny 화면을 사용한다. 메시지 대화창·채팅 버블·입력 composer·장기 대화 타임라인을 만들지 않는다. render.hint를 요청 제목으로만 크게 보여주고 검증 본문을 숨기지 않는다. 승인 버튼을 GET 링크로 만들지 않는다. 색만으로 승인·거절·만료를 구분하지 않는다. 만료나 권한 확인 실패 상태에서 활성 승인 버튼을 보이지 않는다.

별도 아트 에셋이나 장식 목업은 현재 필요하지 않다. 기존 D04와 일반 서비스의 [화면 수락 기준](../../design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md)을 연결한다. DEV의 UI 후보가 준비되면 designer가 MVP-16을 직접 시각 검수한다. 정상 승인 대기, 결정 완료, 만료·철회·원문 부재 화면을 캡처한다. 시간 변화가 정지 화면으로 판정되지 않을 때만 짧은 영상을 만든다.

## 미정 결정과 기능 수락 경계

[백로그의 결정 보류표](../SAR-MVP-backlog.md#미정-결정과-재개-조건)가 담당·영향·재개 조건의 정본이다. Free N·가격·slot-unit, disclosure·결과 schema·window/granularity/누적 한도는 기존 held를 유지한다. 추가 resource/rate/size/concurrency 수치는 아래 DEC-03 공개 기준에 승인 근거와 함께 기록한다. 미확정 수치를 DEV/OPS나 합성 fixture가 대신 결정하지 않는다. 미니서버 설정과 어댑터 인터페이스는 담당 dev/ops의 기술 확인이다.

첫 기능은 합성 요청의 로컬 안전 전달과 human-gate다. 실제 calendar 조회·외부 부작용·실데이터 silent 성공은 수락 범위에 포함하지 않는다. `schedule.query`는 정책 없음의 거부를 검증한다. 승인과 정보 공개는 별개다. `schedule.commit`은 승인 후에도 실행하지 않는다. positive silent `done` 수락은 공개 정책과 output schema가 확정된 후속 기능에 남긴다.

## 검증 책임과 추적

[백로그](../SAR-MVP-backlog.md)는 요구사항별 기능·담당·선행 조건을 연결한다. [DEV 지시서](../../../handovers/to_dev.md)는 사용자 구현 시작 승인으로 ready다. 실제 배정·진행 상태는 coor가 관리한다. [TESTER 준비 지시서](../../../handovers/to_tester.md)는 고정 구현 후보를 기다린다. 이전 준비 과제의 queued 상태를 현재 시작 금지로 해석하지 않는다.

DEV는 같은 과제에서 기술 계획, 실제 기능 구현, 변경 동작 자동 검사·관련 회귀, 필요한 짧은 로컬 확인, D03과 필요한 D05–D10 갱신을 완료한다. tester는 DEV 완료 뒤 고정된 안정 통합 후보에서 독립 QA를 수행한다. designer는 같은 후보 UI를 직접 검수한다. 구현자와 다른 세션의 독립 fixed-SHA 코드 리뷰는 coor가 준비한다. 미해결 critical/high는 수락·병합을 차단한다.

제품 최종 수락은 DEV 완료, 독립 QA, 직접 UI 검수, 독립 코드 리뷰, 정책 보류 해결을 확인한 뒤 coor가 조정한다. 이번 준비 문서 검사는 해당 제품 수락 증거를 대신하지 않는다. 기존 골격의 검증은 원래 SHA·조건과 의존성 동일성을 확인해 재사용한다.

## 개정 이력

- 2026-10-03: SAR-PREP-002에서 새 MVP D02를 작성했다. A2A v0.3.0 잠금, C1–C5, 기능 수락 ID와 미정 결정의 경계를 연결했다. 기존 초기 구성 D02는 보존했다.

## 첫 인증 파일럿의 공개 제품 기준 — DEC-03

이 절은 SAR-MVP-PUBLIC-POLICY-001의 현재 제품 기준이다. 원천 `7bc9ea1`의 wire·상품·보안 잠금은 보존한다. 원천의 `relay.knowslog.com`은 이전 기술 호스트 기록이다. 이번 사용자 지정 공개 hostname `link.knowslog.com`을 적용하며 원천을 소급 수정하지 않는다.

사용자 결정에 따라 누구나 가입할 수 있는 공개 서비스다. 초대 전용 가입이나 사전 owner 명단으로 제한하지 않는다. agent 사이의 pairing 초대와 B-human 수락은 기존대로 필수이며 서비스 가입 초대와 구분한다. 첫 공개 단계는 등록·owner 바인딩·페어링·합성 안전 요청·human-gate를 검증하는 인증 파일럿이다. 공개 URL은 누구나 업무를 실행할 수 있다는 뜻이 아니다. 누구나 가입 시작 경로에 접근할 수 있다. owner 인증과 키 PoP 등 등록 요건을 충족한 뒤 agent를 바인딩한다. 인증되지 않은 주체는 등록 변경·초대/수락·연락처·메시지·receipt·gate·결정 기록에 접근하지 못한다. agent credential은 owner 작업을 대체하지 못한다. owner는 자신의 에이전트와 관계 및 gate만 다룬다.

실제 일정·개인정보·실벤더·외부 도구 효과를 연결하지 않는다. `schedule.query` 정책 없음은 deny다. approve는 해당 gate만 통과시키며 일정 공개나 `schedule.commit` 실행을 허용하지 않는다. `relay.result`의 기존 최소 status와 부모 binding을 유지한다. optional result/error 데이터와 positive silent done은 DEC-02 확정 전 보류한다. evidence 자동 fetch/preview·webhook은 OFF다.

### 제한과 거부의 공통 의미

기존 `body` JCS `16KiB`, intent `64 bytes`, hint `1024 UTF-8 bytes`, evidence `8`개, idempotency key `16–128 ASCII`를 유지한다. `MAX_TTL=300s`, lease `30s`, attempts `3`, receipt·멱등 `24h`, revoke metadata `≥24h`도 유지한다. 추가 전체 봉투 한도는 기존 body 한도를 대체하지 않는다. pending invite는 active pair 슬롯에 포함하지 않는다. 파일럿 자원 상한은 Free N·상품 slot-unit 결정이 아니다.

한도를 넘으면 새 작업을 수락하지 않는다. 거부된 작업을 queued·delivered·approved·done으로 표시하지 않는다. 거부는 새 전달·실행권·gate·관계를 만들지 않으며 기존 요청의 TTL이나 권한을 연장하지 않는다. 동시 요청과 `priority:high`도 같은 한도를 지킨다. 현재 인증·권한 검사와 frozen receipt-only replay를 우회하는 제한 검사는 금지한다. 신규 등록·초대·enqueue의 수용량을 기존 요청 정리에 재적용하지 않는다. 기존 유효한 ACK·deny·철회·unpair·receipt-only replay는 자원 수용량이 가득 찼다는 이유만으로 거부하지 않는다. 이 경로도 인증·현재 권한·lease·만료·CSRF·봉투 한도와 남용 방어를 지킨다. DEV/OPS는 안전 정리 경로의 별도 rate·동시성 budget과 경계 검사 근거를 인계한다. 정확한 API·오류 코드·제한 구현 방법은 DEV/OPS가 정한다.

### 초기 안전 한도 권장안 — 사용자 승인 전 제안값

아래 수치는 운영 측정 결과나 확정값이 아니다. coor를 통해 사용자에게 초기 정책 묶음의 승인 또는 조정 값만 질문했다. 공개 가입을 유지하면서 무제한 운영을 피하기 위한 보수적인 초기 권장안이다. 두 agent의 합성 흐름과 gate를 운영하고 작은 서비스 수용량에서 검증하는 목적이다. 성능 보장·가입자 목표·Free N으로 사용하지 않는다.

| 대상 | 제안값·단위 | 집계·경계와 거부 동작 |
|---|---|---|
| 등록 owner / agent | 서비스 전체 owner 100명 / agent 200개 | 삭제·철회만으로 필수 보존 중인 보안 기록을 없애지 않는다. 새 등록이 상한을 넘으면 등록을 거부한다. 누구나 같은 가입 조건을 적용받는다. |
| active pair / pending invite | 서비스 전체 각각 400개 / 200개 | active는 현재 유효 관계, pending은 아직 결정되지 않은 초대다. pending은 active에 더하지 않는다. 중복 active 관계나 같은 pending 초대의 재시도는 새 개체를 만들지 않는다. 초과 신규 초대/수락을 거부하며 deny/unpair는 허용된 안전 종료 경로를 유지한다. |
| 전달 대기 봉투 | 서비스 전체 queued+leased 합계 100개 | M/H/R 모두 해당 상태면 각각 한 개로 집계한다. delivered/terminal은 제외한다. 새 enqueue가 상한을 넘으면 거부한다. 유효한 receipt-only replay는 새 queue 자원을 소모하지 않는다. |
| pending gate | 서비스 전체 100개 | 권위 있는 부모 M별 pending gate 한 개다. approved/denied/expired/revoked는 제외한다. 초과 신규 gate를 거부하고 기존 gate의 deny/철회는 보존한다. |
| receipt·멱등 기록 | 서비스 전체 24h 보존 중인 수락 메시지 20000건 | 같은 메시지의 receipt와 멱등 기록은 함께 한 건이다. H/R도 새 수락 메시지면 집계한다. 초과 신규 수락을 거부한다. 기존 replay로 증가시키거나 보존을 24h 미만으로 줄이지 않는다. |
| 익명 HTTP 시도 | source IP당 rolling 60s에 30회 | 인증 실패·거부 시도도 집계한다. 새 source IP로 서비스 전체 상한을 우회하지 못한다. NAT 공유 영향과 신뢰할 source IP 판별은 DEV/OPS 근거가 필요하다. |
| 인증 HTTP 시도 | owner 또는 agent principal당 rolling 60s에 60회 | credential/키 교체로 같은 principal의 한도를 초기화하지 않는다. 신규 작업 40회와 기존 안전 정리 20회로 분리하는 보완안을 제안한다. replay·거부·pull/ACK·owner 작업도 해당 경로에 집계한다. |
| 전체 HTTP 시도 | 서비스 전체 rolling 60s에 300회 | 익명·인증 시도의 합이다. 신규 작업 200회와 기존 안전 정리 100회로 분리하는 보완안을 제안한다. 비공개 운영 health 확인만 제외한다. 사용자 API를 health 예외로 우회하지 못한다. |
| 추가 봉투 크기 | 봉투 전체 raw UTF-8 32KiB = 32768 bytes | JSON의 공백·sig·ext·trace·render·evidence를 포함한다. 기존 body JCS 16KiB는 별도로 지킨다. 초과 봉투를 수락하지 않는다. 압축·encoding 우회 방어는 DEV가 정한다. |
| 진행 중 HTTP / 실행 claim | 서비스 전체 각각 20개 / 4개 | HTTP는 처리 시작부터 종료까지 집계한다. HTTP 20개 중 신규 작업 16개·기존 안전 정리 4개를 별도 유지하는 보완안을 제안한다. claim은 실행권 획득 뒤 처리 완료/중단까지 집계한다. 추가 동시 작업은 수락하지 않는다. 숨은 무제한 대기열을 만들지 않는다. |

rolling 60s는 시각 t에서 `(t−60s,t]`에 시작한 시도를 뜻한다. 상한 값까지 허용하고 다음 시도를 거부한다. 적용 가능한 모든 한도를 만족해야 한다. 제한 시도도 rate에 포함하므로 재시도 폭주가 허용량을 늘리지 않는다. 시각·집계 상태가 불명확하면 새 작업은 fail-closed다. 새로운 process·credential·pair 세대가 서비스 전체 수용량이나 보존 중인 기록을 초기화하지 않는다.

coor의 검토 요청에 따라 신규 수용량과 기존 안전 정리의 경계를 보완했다. 위 rate·HTTP 동시성의 분리값도 승인 전 제안이다. 신규 작업이 안전 정리의 budget을 소진하지 못한다. 안전 정리는 인증된 해당 기록에 대한 ACK·deny·철회·unpair·receipt 조회/replay다. pull은 새 lease를 만드는 신규 작업으로 집계한다. gate approve·새 H/R은 새 인가/수락을 만들므로 안전 정리 예외를 자동 적용하지 않는다. 제한이 가득 찼을 때 H/R 수락과 부모 종료의 연계 또는 필요한 예약량은 DEV 근거를 받은 뒤 제품 결정으로 확정한다.

rate·HTTP 동시 한도와 안전 종료 경로의 충돌은 아직 기술 검증 전이다. DEV/OPS는 ACK·deny·철회·receipt 조회의 반복 실패가 기존 요청 종료나 권한 차단을 지연시키지 않는다는 근거를 제공한다. 필요하면 designer가 안전 종료 경로별 별도 한도를 사용자 결정에 포함한다. 이 검증 전에는 “모든 조건 확정”이나 “공개 가능”으로 표시하지 않는다.

CPU·메모리·디스크·DB 보존 자원의 실제 보호 상한은 OPS가 현재 서버와 기존 서비스 보존 조건을 측정해 근거를 제공한다. designer는 그 근거가 요청 수용량·거부 동작 변경을 요구하면 제품 기준을 갱신한다. 장기 key revoke metadata와 기록 증가, 정상 polling·ACK·owner 작업의 budget, 부모/결과 관계가 가득 찬 queue에서 안전 종료되는지를 DEV/OPS가 확인한다. 임의의 CPU/RAM/디스크 값을 제품 정책으로 확정하지 않는다.

### 사용자 완료 조건과 후속 담당

- DEV는 MVP-01–16의 합성 흐름과 승인된 한도의 적용 대상·경계·거부·동시성·재시작 동작을 구현하고 검증한다. 수치가 미정이면 기존 로컬 합성 구현은 계속하되 공개 경로를 열지 않는다. 제안값만으로 현재 DEV 과제에 공개 한도 구현을 추가하지 않는다. 확정 기준과 coor 후속 인계가 있어야 제한 집행 범위를 반영한다.
- tester는 coor가 지정한 고정 통합 후보에서 인증·owner/agent 권한 분리·공개 거부 경계·한도 이하/경계/초과·high priority·철회/만료·재시작을 독립 검증한다. 실제 데이터나 외부 효과 없이 수행한다.
- designer는 같은 후보의 verified typed body·정책·만료/철회·원문 부재·approve/deny 상태를 직접 검수한다. 인증·제한 거부를 성공으로 표시하지 않는지도 확인한다.
- OPS는 수락 SHA와 승인된 제품 한도를 적용한 인증 공개 후보를 검증한다. 기존 서버 서비스와 route를 보존한다. 설정·배포 방법·운영 자원 보호는 OPS의 기술 책임이다.
- coor는 독립 QA·직접 UI 검수·별도 세션의 fixed-SHA 독립 리뷰·제한 집행 근거를 확인한다. 미해결 critical/high와 공개 선행 조건 실패는 수락·공개를 차단한다.

제품 정책 문서 완료는 배포 성공이나 QA PASS가 아니다. 이 합성 파일럿 수락은 전체 MVP 완료가 아니다. 실벤더 연결과 실데이터 silent 업무 성공·유료화의 기존 held 및 재개 조건을 유지한다.

- 2026-10-03: SAR-MVP-PUBLIC-POLICY-001에서 누구나 가입 가능한 인증 합성 파일럿 범위를 반영했다. DEC-03 수치는 승인 전 제안이며 신규 수락과 기존 안전 정리의 경계를 분리했다. 실제 공개와 전체 MVP 수락은 선언하지 않았다.

- 2026-10-05: SAR-PUBLIC-SERVICE-001 일반 서비스 확장 정본을 연결했다. 기존 DEC-03 제안은 당시 기록으로 보존한다.
