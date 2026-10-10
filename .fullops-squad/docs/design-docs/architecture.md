---
id: D03
title: 아키텍처설계서
status: review
updated: 2026-10-10
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-MVP-002-DEV, SAR-MVP-003-BIDIRECTIONAL, SAR-PUBLIC-IDENTITY-001-DEV, SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX, SAR-PUBLIC-AGENTS-001-DEV, SAR-PUBLIC-MESSAGES-001-DEV, SAR-PUBLIC-MESSAGES-001-DEV-FIX, SAR-PUBLIC-MESSAGES-001-DEV-FIX-2, SAR-PUBLIC-MESSAGES-001-DEV-FIX-3, SAR-GOOGLE-CONNECT-001-DEV, SAR-GOOGLE-CONNECT-002-DEV, SAR-AUTO-RECEIVE-001-DEV]
upstream: [D02]
summary: 로컬 합성 relay와 shared 상태 및 owner gate의 인가 경계를 정의한다
---

# KnowsLink 로컬 합성 MVP 아키텍처

## 범위와 추적

정본은 [D02 SAR-MVP](../planning/product-specs/SAR-MVP.md)와 고정 원천 service-design 7bc9ea1이다.
fullops-common-0.3.2와 lint 기준 0dd08ec994771836c15d9d22a6a83393a71d7987을 적용한다.
SAR-SETUP-001의 초기 골격 이력은 기존 실행 기록에 보존한다.
이번 구현은 합성 데이터의 안전 전달·human-gate 후보다. 독립 QA·UI 검수·코드 리뷰·운영 수락은 후속이다.

## 실제 책임

- `cmd/relay`: pgxpool·readiness·HTTP API·Go html/template UI와 retention 정리 루프.
- `cmd/migrate`: 별도 goose SQL-only Up. relay 기동은 migration을 실행하지 않는다.
- `internal/relay/protocol.go`: strict JSON·closed seed schema·RFC8785·Ed25519·semantic digest.
- `internal/relay/store.go`: shared Postgres 상태의 transaction·epoch CAS·current-auth·queue·receipt·lease·gate.
- `internal/relay/http.go`: 별도 owner/agent 인증·PoP·pairing·durable inbox·ACK·claim·UI CSRF.
- `internal/relay/registry.json`: 무결성 hash로 고정한 central seed revision. intent 추가는 실행 handler 등록이 아니다.
- `internal/database`: sqlc v1.30.0의 pgx/v5 생성 코드.
- `adapters/src/index.ts`: 합성 pull stub 1개. signature 재검증·shared inbox persist·ACK 성공·claim 뒤에만 판단한다.

## transaction과 인가

`relay_state`의 singleton 행은 업무 JSON과 authorization epoch를 보관한다.
모든 읽기·수정은 `SELECT FOR UPDATE`와 `UPDATE WHERE epoch=expected`로 확정한다.
철회·enqueue·lease·ACK·decision·consume·authorize·result 공개가 같은 직렬화 경계를 사용한다.
DB 시계가 이전 확정 시각보다 뒤로 가면 fail-closed다. DB를 읽을 수 없으면 `unavailable`이다.

단일 행과 전체 상태 순회는 로컬 MVP의 처리량 한계다. 24h receipt 20000 등 미확정 제안의 성능을 보장하지 않는다.
공개 운영 전에 제품 한도·인증을 확정하고 정규화 또는 처리량 측정과 독립 수락을 수행한다.

## 안전 흐름

owner 가입은 로컬 합성용 opaque credential을 발급한다. agent credential과 분리한다.
owner ID·AgentID·kid·pubkey에 묶인 PoP로 키를 등록한다. rotate는 이전 키를 원자적으로 revoke한다.
B-owner 수락 전 pair는 pending이다. pending은 active 관계에 포함하지 않는다. 재수락은 새 세대다.

send는 strict structure·signature·principal·routing 뒤에 digest와 atomic idempotency를 검사한다.
동일 key+digest는 receipt만 반환한다. exp·TTL·id 실패는 전체 작업 상태를 rollback한다.
agent credential은 `deliver:agent` 메시지만 lease·persist·ACK·claim한다. `deliver:human`은 owner gate만 처리한다.
authorize·gate-consume·H·R의 부모도 저장 경로가 `deliver:agent`인 claim만 허용한다. 경로 미기록 이전 claim은 거부한다.
lease는 delivered가 아니다. 공유 inbox에 원문을 저장한 뒤 ACK하고 하나의 claim token을 발급한다.
claim 재발급은 하지 않는다. 재시작 후 이미 claimed인 요청은 중복 실행 대신 TTL까지 안전하게 정지한다.
외부 도구 exactly-once나 crash 후 효과 재개를 주장하지 않는다.

H는 현재 부모 receipt·실제 수신·claim token·digest·pair 세대에 결속한다.
owner UI는 검증된 원요청 typed body와 deny/stub 정책을 표시한다. hint는 승인 근거가 아니다.
POST와 owner credential에 묶인 CSRF token으로만 결정한다. consume은 한 번만 성공한다.
approve 후에도 `authorize`의 executable/disclosure는 false다. schedule.commit과 schedule.query done은 실행하지 않는다.
R은 새 B 서명이며 부모 endpoint를 반전한 결과다. optional result/error schema가 없으므로 해당 필드를 거부한다.

## 저장과 경계

원문·공유 inbox는 exp·철회·응답 완료에 삭제한다. 모든 읽기는 먼저 만료를 정리한다.
유휴 서버도 1초 정리 루프로 payload를 지운다. DB 불가 시 정리는 중단하고 안전한 오류만 기록한다.
receipt·idempotency metadata는 24h 유지한다. revoked key의 kid와 metadata는 재할당 방지를 위해 계속 유지한다.
DB 삭제는 WAL·backup 완전 삭제가 아니다. payload·credential·키·tool 정보는 로그에 쓰지 않는다.
webhook·evidence fetch·preview·벤더 연결은 없다. 실제 calendar와 유용한 silent done은 DEC-02 이후다.

Compose는 기존 네 서비스·private Postgres·loopback relay·선택 Tunnel을 유지한다.
`make verify-mvp`만 고유 project의 시험 DB를 loopback 임시 포트로 연결한다. 종료 시 자기 project만 지운다.
기존 컨테이너·볼륨·Tunnel은 수정하지 않는다. 공개 hostname 운영은 수락 후보 이후 OPS가 담당한다.

## 요구사항 연결

MVP-01/02/07은 owner·key·pair·CAS와 Postgres 권한 검사에 연결된다.
MVP-03–06은 protocol·transaction·durable inbox·race 검사에 연결된다.
MVP-08–11/15/16은 gate·authorize·result·retention·HTTP/UI 검사에 연결된다.
MVP-12–14는 field 한도·registry·loopback stub·OFF 경계로 유지한다.
실행 증거와 보류는 [SAR-MVP-001-DEV 기록](../exec-plans/phases/SAR-MVP-001-DEV.md)에 있다.

## SAR-MVP-002 공식 Grok Bot 플러그인 준비

사용자가 공식 Grok Bot을 확정했다. hosted 컴퓨터·MCP·Cursor connector policy의 공식 근거를 확인했다. 연결 구조는 Cursor plugin manifest·stdio MCP·skill이며 실제 account/hosted 연결은 held다. 준비 package는 [설치 문서](../../../adapters/README.md)를 따른다.

공통 Adapter는 `adapters/src/core.ts`로 옮겼다. 기존 `index.ts` CLI와 새 MCP가 같은 서명·persist·ACK·claim·gate·deny를 재사용한다. relay·DB·UI·frozen wire는 변경하지 않았다. MCP는 payload·credential·claim을 모델에 노출하지 않고 원문 업무를 추론하거나 도구로 실행하지 않는다.

Grok Bot의 같은 계정 Bot들은 파일과 command-line credential을 공유한다. KnowsLink AgentID/owner 분리는 제품 서버에서 집행한다. Bot 프로필·화면 분리나 vendor approval을 identity·owner approval 경계로 간주하지 않는다. plugin의 실제 설치와 최소 계정 권한·도달 경로를 후속 고정 버전에서 검증한다. shared claim 없는 다중 adapter 활성화는 허용하지 않는다.

## SAR-MVP-003 승인된 양방향 시험

이번 시험은 사용자 승인된 Codex↔Grok 시험 text만 전달한다. 이전 SAR-MVP-002의 actual held는 이 범위에서만 명시 모드로 재개한다. 업무 disclosure·calendar·dots·자동 wake는 범위 밖이다.

`relay.test.message`는 기존 relay.v1 서명·pairing·TTL·idempotency와 저장소를 재사용한다. registry revision은 `relay.v1-test-2026-10-04`이며 기존 네 intent의 schema를 보존한다. 시험 intent는 body.text만 허용한다. evidence·ext·render·reply_to를 시험 권한 확장 통로로 쓰지 못한다. trial을 H/R·authorize·gate-consume 부모로 쓰는 것도 차단한다.

서버 `KNOWSLINK_TEST_AGENTS`는 기본 빈 값이다. 명시한 두 서로 다른 AgentID만 시험 envelope의 양쪽 endpoint로 허용한다. `State.TestAgents`는 DB로 직렬화하지 않는 실행 설정이다. 정상 transaction과 오류 rollback에서 같은 설정을 적용한다. config를 해제하면 current 검사가 fail-closed로 payload와 claim을 제거한다.

machine 경로는 `/v1/test/*`다. relay agent 인증·시험 allowlist·active pair·서명을 유지한다. 시험 pull은 business 메시지를 제외한다. persist/ACK/claim은 시험 message ID만 허용한다. 가입·owner·pairing·키변경·business send·authorize는 machine 경로에서 차단한다.

시험 receive는 configured peer와 서명을 확인하고 shared persist→ACK→claim 뒤에만 text를 노출한다. claim 시 원문·inbox를 삭제한다. 반환 text는 `untrusted:true`이며 실행 권한을 만들지 않는다. claim 뒤 출력 전 crash는 표시를 잃을 수 있다. 새로운 key로 자동 재전송하지 않는다.

remote는 기존 `https://link.knowslog.com`만 사용한다. prefix 전용 Service Auth 앱과 distinct agent service tokens를 준비한다. Tunnel의 더 구체적인 path rule은 trial AUD만 검증한다. root/owner rule은 기존 owner AUD를 유지한다. 실제 적용은 독립 fixed-SHA 검토 후 OPS/coor가 수행한다. [D12](../operations/ops-guide.md#13-승인된-양방향-시험-sar-mvp-003)와 [실행 기록](../exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md)을 따른다.

## SAR-PUBLIC-IDENTITY-001 일반 이메일 신원과 세션

상위 요구사항은 [D02 일반 이메일 서비스](../planning/product-specs/SAR-PUBLIC-SERVICE.md)의 PS-01–04와 신원 관련 PS-11이다. 화면 기준은 [UX-01–03](mockups/SAR-PUBLIC-SERVICE-UX.md)이다. 제품 고정 SHA는 `1233e4c3167f722d51f99cb2ef495691734be714`다. 기술 계획과 근거는 [실행 기록](../exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV.md)에 있다.

### 신원 제공자 선택

relay가 이메일 확인 코드를 직접 발급하고 검증한다. 발송은 표준 SMTP submission이다. 운영 후보 경로는 Cloudflare Email Service SMTP(`smtp.mx.cloudflare.net:465`, implicit TLS)다. 다른 SMTP 제공자도 같은 설정으로 쓸 수 있다.

Cloudflare Access One-time PIN을 회원 신원으로 쓰지 않는다. 이유는 다음과 같다.

- D02의 오답 5회·재발송 60s·이메일/IP/전체 시간당 한도를 relay가 관측하거나 집행할 수 없다.
- relay 로그아웃 뒤에도 Access 세션 cookie로 이메일 재확인 없이 다시 들어올 수 있다.
- root 정책을 모든 OTP 이메일로 열어야 한다. Zero Trust 좌석 한도와 요금은 미확인이다(OPS 준비 보고 `0313deae`).

기존 owner-only Access 앱과 Basic owner UI는 운영 관리 경계로 유지한다. 일반 회원 인증으로 사용하지 않는다.

### 회원·owner·세션 경계

- 회원은 `(issuer, 정규화 이메일)` 하나에 묶인다. issuer는 `knowslink-email-otp`다. 다른 issuer의 같은 이메일은 자동 병합하지 않는다.
- 정규화는 앞뒤 공백 제거와 전체 소문자화다. `+tag`와 점 별칭은 별도 신원으로 유지한다. 제공자별 별칭 규칙을 추측해 계정을 합치지 않는다.
- 코드 확인 성공 전에는 회원과 owner를 만들지 않는다. 첫 확인이 회원과 owner를 같은 transaction에서 만든다. 회원 owner는 bearer credential이 없다.
- 같은 global row lock이 동시 첫 가입을 직렬화한다. identity index가 같은 신원의 두 번째 회원 생성을 막는다.
- 브라우저 세션은 random token의 SHA256만 저장한다. cookie는 `__Host-kl_session`, `Secure`, `HttpOnly`, `SameSite=Strict`다.
- 세션 수명은 절대 12h, 무활동 60분이다. 전체 로그아웃은 최근 5분 안의 이메일 확인을 요구한다.
- 로그아웃은 브라우저 세션만 지운다. agent credential·키·pair는 세션과 수명이 분리돼 있다.
- 쓰기 요청은 Go 표준 `http.CrossOriginProtection`이 `Sec-Fetch-Site`/`Origin`으로 cross-origin 브라우저 요청을 거부한다. 회원 gate 결정은 기존 HMAC CSRF를 세션 token에 결속한다.
- `/v1/owners` 합성 가입은 `KNOWSLINK_SYNTHETIC_SIGNUP=1`일 때만 동작한다. `.env.example`과 공개 후보는 이 값을 비워 두어 외부 owner 발급 우회를 막는다. 격리 로컬 fixture 실행만 셸에서 명시적으로 `1`을 준다.

### 한도 집행

한도는 `State.Rates`의 시각 목록으로 계산한다. 창은 `(t−window,t]`다. 발송 한도는 모든 bucket이 허용할 때만 한 번에 차감한다(`take`). HTTP rate는 거부된 요청도 센다(`hit`). `hit`은 principal(익명 IP·회원) bucket을 먼저 기록하고 첫 거부에서 멈춘다. 자기 principal 한도로 거부된 요청은 전체 budget을 쓰지 않는다. 한 principal은 rolling 60s에 전체 budget에 최대 자기 한도(30·40·정리 20)만 기여한다. 전체 200·정리 100은 독립 principal들의 합으로만 포화한다. 공유 bucket이 거부한 요청은 앞 principal bucket에 창 안 기록이 이미 있을 때만 두 bucket에 기록한다. 기록이 없는 새 principal의 요청은 어디에도 기록하지 않는다. 그래서 source를 바꾸는 거부 요청은 새 rate key를 만들지 않고 공유 포화를 연장하지 않는다. 로그아웃은 신규 작업 budget과 다른 정리 budget을 사용한다. rate key에는 이메일 원문 대신 SHA256을 쓴다.

메일 발송은 lock 밖에서 수행한다. 발송 전에 budget과 코드를 확정한다. 발송 실패는 spent budget을 유지하고 보내지 못한 코드를 지운다.

### 다음 agent·원격 MCP 연결의 신원 경계

회원 ID(`mem_…`)가 사용자별 권한의 주체다. 다음 기능 SAR-PUBLIC-AGENTS-001은 세션과 최근 재인증으로 연결 승인을 발급하고 agent를 회원 owner에 묶는다. 이메일은 agent 식별자·contacts·receipt에 넣지 않는다.

향후 다닷의 원격 MCP/OAuth는 같은 회원 ID에 OAuth grant를 연결한다. 후보는 relay를 OAuth 2.1 authorization server로 두는 방식과 Access Managed OAuth다. 두 방식 모두 이번 범위에서 구현하지 않는다. 선택 전에 대상 client의 RFC 8707·9728 지원을 실제로 확인한다. 회원 확인 수단은 이번 이메일 코드 로그인을 재사용한다.

## SAR-PUBLIC-AGENTS-001 일반 회원 연결 경계

PS-04–07·해당 PS-11과 UX-04–05를 구현한다. 기존 singleton 직렬화·회원 세션·Go template·Ed25519를 재사용한다. 새 dependency·migration·frontend·frozen wire 변경은 없다.
회원 화면은 세션에서 owner를 결정한다. agent ID는 이메일과 무관한 무작위 식별자다. 연결은 발급 → 새 공개키 PoP 준비 → 최근 재인증 owner의 지문 확인 → 클라이언트의 1회 완료 순서다. 발급·준비·승인만으로 새 agent 권한을 만들지 않는다. grant는 owner·agent·node-local·등록/회전에 결속하고 10분 뒤 사용할 수 없다.
회원 키는 Key.Credential 해시에 묶인다. agent-wide legacy credential은 일반 연결 완료 시 제거한다. 선택 철회와 회전은 철회 키 credential의 모든 API 인증을 막는다. 기존 합성 owner credential API는 회귀 호환을 유지한다. 합성 키 회전은 기존 전체 철회 동작이다.
관계·키·agent 한도는 state 변경 경계에서 검사한다. HTTP rate는 DB에 저장한 stable 회원/agent principal로 집계한다. 신규·정리 budget을 분리한다. 새 key credential이나 relay 재시작은 budget을 초기화하지 않는다. pending과 active slot은 분리한다. same-owner pair도 수신 owner의 첫 수락을 요구한다. 브라우저 결정은 pair 세대를 비교한다.
실제 이메일·외부 계정·공개·일반 text·다닷·OAuth·처리량 보장은 후속이다. HTTP 동시 수용 16/4·claim 4·queue/gate/receipt 용량은 MESSAGES/공개 수락 후속 범위다. 이번 과제의 경합은 기존 DB row lock으로 직렬화했다. 상세 구현·검증·인계는 [실행 기록](../exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md)에 있다.

## SAR-PUBLIC-MESSAGES-001 일반 text·gate·보호

PS-08–11·UX-06/07의 일반 회원 text는 `knowslink.text.v1`·`/v1/text/*`로 업무 relay.v1과 분리한다. 시험 allowlist·Service Auth·owner 공유 자격을 쓰지 않는다. 양쪽 검증 회원 owner·현재 agent/key·active pair·관계 세대를 수락·lease·persist·ACK·답장마다 확인한다. text는 business claim/H/R 부모가 되지 않는다.
명시 송신·한 원요청의 관련 답장 1회만 지원한다. 최대 4096 UTF-8 bytes·TTL 180s다. 원문은 수신 persist→ACK 또는 만료·철회·3회 lease 실패에 지운다. ACK 후 출력 전에 client가 죽으면 text 표시를 잃을 수 있다. 자동 재송신·재답장·wake·장기 타임라인은 없다. receipt·멱등 metadata는 24h 보존한다.
공통 수락 경계가 queue100·pending gate100·receipt20000과 claim4를 검사한다. high priority와 재시작도 상한을 우회하지 못한다. replay는 새 slot을 쓰지 않는다. H/R이 포화하면 수락하지 않는다. 이미 수락한 부모는 TTL에서 failed:expired로 안전 종료하고 claim을 해제한다. 성공 처리를 약속하지 않는다.
HTTP 신규16·정리4는 즉시 거부하는 프로세스 채널과 공유 JSONB 입장 기록으로 제한한다. 정상 종료는 기록을 지우고 crash는 30s 뒤 회수한다. DB context·socket body read는 10s다. 8/32KiB 상한 본문은 모든 슬롯보다 먼저 같은 10s 기한 안에서 수신한다. 느린 송신자는 자기 연결만 점유한다. 이 점유는 header 수신과 같은 연결 계층이며 DB·공유 상태를 만들지 않는다. 모든 입장 요청의 신원별 rate를 schema·CSRF·route 검사 전에 한 번 차감한다. 입장 전 로컬 동시 상한 거부는 DB 대기열을 만들지 않는다. 정리 입장(로컬 정리 채널·공유 Clean 기록·cleanup rate)은 검증된 principal이 같은 출처·gate CSRF·유효 본문으로 자기 기록을 정리할 때만 쓴다. 익명·잘못된 자격·타 owner 기록·CSRF·잘못된 본문은 신규로 집계한다. 로컬 채널은 이 프로세스의 마지막 커밋 상태 snapshot에서 자격 색인과 자기 기록 대조(cleanupTarget)를 모두 통과한 요청만 정리 채널로 고른다. 타 owner·lease 없는·위조 정리는 DB 입장을 기다리는 동안에도 신규 채널만 쓴다(SAR-PUBLIC-MESSAGES-001-DEV-FIX-2 H-2). 공정성 단위는 owner 자신의 제어(철회·unpair·deny·logout)와 그 owner의 agent ACK 두 개다. 각 단위는 로컬 정리 채널과 공유 Clean 기록에서 각각 동시 1개만 쓴다. agent의 ACK 반복은 owner의 자기 agent 철회를 막지 못한다. 자기 대상 반복·rate 초과 정리도 다른 owner의 정리를 막지 못한다. transaction이 다시 판정하며 판정이 신규로 바뀌면 정리 슬롯을 신규 슬롯으로 돌려준다. snapshot은 (DB commit 시각, epoch) 순서로만 교체한다. 같은 commit 시각에도 늦게 저장된 이전 snapshot이 최신을 덮지 않는다. 신규 채널이 가득 찬 상태에서 snapshot이 정리 자격을 증명하지 못하면 거부 전에 커밋 상태를 한 번 다시 읽는다(SAR-PUBLIC-MESSAGES-001-DEV-FIX-3). 신규 채널에 여유가 있으면 신규로 대기하고 transaction이 공유 budget을 판정한다. 이 읽기는 lock·슬롯 없는 단일 조회이며 프로세스당 동시 1개다. 요청 도착 뒤 시작한 읽기를 공유한다. 다른 프로세스의 새 자격과 재시작 직후 빈 snapshot도 1s sweep을 기다리지 않는다. 위조·타 owner 정리는 다시 읽어도 신규 채널을 쓴다. 종료 transaction이 실패한 공유 기록은 30s 만료 전에 같은 프로세스의 다음 commit이 지운다. 명시 deny 경로는 신규 슬롯이 포화해도 수신한다. gate 결정 뒤에는 정식 gate 화면으로 303한다.
회원 gate GET/POST가 검증본문·서명·digest·현재 세대·기한을 검사한다. hint는 escaped 참고 데이터다. approve는 인간 게이트만 통과시키며 query disclosure deny·commit 실행 불가를 유지한다.
실제 로컬 Node CLI/MCP·일반 신원 HTTP·격리 Postgres 증거는 [실행 기록](../exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV.md)에 있다. 운영 배포·실메일·실제 Grok Bot/다닷·독립 QA·직접 시각 수락·실부하와 복원은 후속이다. 단일 행 lock의 처리량 한계는 유지한다.

## Google 클라이언트 연결 — SAR-GOOGLE-CONNECT-001-DEV

Node가 각 클라이언트에서 키와 임의 token을 만들고 서명한 시작 요청을 보낸다. 브라우저에는 token hash인 공개 연결 ID만 노출한다. 기존 Google issuer/sub 회원·Strict 세션을 재사용한다. Google callback은 요청의 회원만 바인딩하고, 최근 인증 상태의 명시적 동의가 별도 agent를 만든다. Node가 동일 키로 승인 상태를 조회하고 기존 signed complete로 자기 credential을 받아 로컬에 저장한다. 관계 요청·수락은 기존 권한 경계를 유지한다. 새 IdP·의존성·SQL 테이블은 없다. Cloudflare 사용자 경로 개방과 관리자 보호의 운영 변경은 D12에서 별도로 수행한다.
## 비용 없는 공개 연결 경계 — SAR-GOOGLE-CONNECT-002-DEV

기존 Tunnel·도메인 HTTPS와 Google 세션 및 서명·회원 소유권을 재사용한다. 새 Access 가입·원격 OAuth 서비스·유료 보안 계층은 필요 없다. `beta.sh render-public-config`는 public-ingress allowlist와404 fallback만 가진 별도 후보를 생성한다. owner/admin/test와 미허용 경로는 원점에 도달하지 않는다. edge의 기존 두 KnowsLink Access 앱 제거는 coor가 ingress 차단을 먼저 확인한 뒤 수행한다. 운영 적용·복구 정본은 [D12](../operations/ops-guide.md#비용-없는-tunnel-적용--sar-google-connect-002-dev)다.

기존 Device cap2000과 만료 뒤24h 보존을 유지한다. 신규 연결 포화의 medium 한계는 남으며 기존 회원/키 사용·철회에는 이 cap을 적용하지 않는다. 자동 parent 폴더 생성은 Google 클라이언트 시작에만 적용하고 최종 키 폴더 exclusive 생성·ACL은 유지한다. SQL·의존성·UI는 변경하지 않는다.

## 자동 수신과 호스트 알림 — SAR-AUTO-RECEIVE-001-DEV

사용자 요청 “자동수신 기능 보완해”로 회원 text의 수신·보존·호스트 알림만 재개한다. 자동 업무 실행·자동 답장·gate 승인은 범위 밖이다. relay·SQL·의존성은 변경하지 않는다.

- `public-node` MCP는 연결 폴더가 있으면 process마다 자동 수신 loop 하나를 시작한다. `KNOWSLINK_AUTO_RECEIVE=off`로 끈다.
- 순서는 pull→서명·만료·수신자 검증→로컬 private inbox 기록(fsync·rename)→relay persist→ACK→호스트 알림이다. 로컬 기록이 실패하면 ACK하지 않는다. relay가 같은 메시지를 다시 lease한다.
- process 안의 모든 pull은 한 직렬 큐를 지난다. 자동 loop와 수동 receive는 병렬 lease를 만들지 않는다. 서로 다른 process(예: MCP 둘, MCP와 `text.js watch`)는 각자 loop를 가진다. 이 경우 relay lease와 ID 기반 중복 제거가 중복 표시를 막는다.
- idle 간격은 10s다. 오류는 20s부터 2배씩 최대 300s backoff다. `retry_at`이 더 늦으면 그 시각을 따른다.
- 같은 ID는 `.json` 또는 읽음 표시 `.read`가 있으면 다시 저장·알림하지 않는다. ACK는 다시 보낸다. 읽음 표시는 원문 없이 24h 보존한다.
- 수신 뒤 TTL이 지나거나 관계가 철회돼도 이미 검증·보존한 text는 표시한다. `expired:true`이면 관련 답장은 relay가 거부한다. 철회 뒤 새 pull은 key 조회 실패로 끝나며 inbox에 들어가지 않는다.
- ACK 후 표시 전 crash로 text를 잃던 기존 한계는 로컬 inbox로 해소한다. 단, 로컬 inbox의 수동 조회 뒤 표시 전 crash는 그 text를 잃을 수 있다.

### Grok Bot 호스트 경계

호스트 알림은 MCP 표준 `notifications/message`(logging capability)이며 `{event,id,from,pending,next}` metadata만 보낸다. text는 넣지 않는다. 상태의 `hostNotice: sent_unverified`는 process 밖으로 보냈다는 뜻이다. 호스트 표시나 노우 턴 시작을 뜻하지 않는다.

공식 근거로 확인한 경계는 다음과 같다.

- MCP 명세는 logging·resources 알림의 사용 방식을 client에 맡긴다. 알림이 모델 턴을 시작한다는 계약은 없다([resources](https://modelcontextprotocol.io/specification/2025-06-18/server/resources)).
- Grok Bot 공식 문서의 무인 실행은 routine뿐이다. schedule은 최소 5분 간격이다. event 시작은 Cursor 계정 통합(Slack·GitHub)에 한정된다([Skills and routines](https://docs.x.ai/grok-bot/skills-routines-and-automations)). MCP 알림·webhook·외부 inbound API로 대화를 깨우는 공식 기능은 없다.
- Cursor staff는 외부 wake를 “working on” 상태로 답했다([포럼 요청](https://forum.cursor.com/t/let-a-grok-bot-computer-wake-its-own-agent-chat/168260)). 같은 답변이 미문서화 loopback gateway workaround를 연결했다. 이 경로의 선택형 사용은 아래 절이 정한다.
- Command MCP process의 대화 사이 생존 여부는 공식 문서에 없다. 그래서 Bot 컴퓨터의 상시 수신은 `node dist/plugin.js watch <폴더>`(또는 `dist/text.js watch`)로도 제공한다. 같은 inbox를 MCP가 읽는다.

따라서 기본 구현은 자동 수신·보존·표준 알림까지다. 공식 노우 턴 자동 시작은 공식 호스트 기능(MCP 알림 기반 wake 또는 외부 trigger)이 생길 때까지 없다. 5분 routine은 TTL 180s 안의 답장을 보장하지 못하고 사용량을 소비하므로 해결로 인정하지 않는다.

### 선택형 loopback wake — 미문서화 gateway (attempt 0305914f)

coor 재검토로 같은 과제를 다시 열었다. 직전 기록의 “미공개 gateway 우회는 사용 금지”는 사용자 지시가 아니었다. 소유한 Bot 컴퓨터의 인증된 loopback 인터페이스 사용 자체는 금지하지 않는다. 보안 우회·권한 확대·vendor core 수정은 계속 금지한다.

| 구분 | 내용 |
|---|---|
| 공식 지원 없음 | Grok Bot 공식 문서에 외부 wake API가 없다. Cursor staff(Colin, CursorStaff 그룹)는 2026-08-18에 개발 중이라고 답했다. |
| 미문서화 | staff 답변이 연결한 커뮤니티 글(adam91holt, 2026-08-12, [168199/8](https://forum.cursor.com/t/grok-bot-can-i-send-it-a-message-from-outside/168199/8)): Bot 컴퓨터 `127.0.0.1:1340`, `/home/box/sand-data/gateway.json`의 `token`, Bearer 인증 `POST /api/listAgents`·`/api/sendPrompt {agentId,prompt}`. 작성자는 live 설치에서 시험했다고 했다. 작성자 스스로 “undocumented internal API”라고 경고했다. |
| 실제 미검증 | 본인 Bot 컴퓨터의 gateway 존재·port·token 필드·응답 코드·노우 턴 시작·진행 중 턴과의 관계·Update/Reset 뒤 유지. DEV는 Bot 컴퓨터에 접근하지 않았다. |
| 불가·제외 | Bot 컴퓨터 밖에서 gateway 호출(port 공개·SSH tunnel·Tailscale), 받은 text를 prompt로 전달, 응답 본문 표시. |

설계는 다음과 같다.

- Bot 컴퓨터 안의 `plugin.js watch`(또는 MCP 자동 수신)가 `127.0.0.1`의 gateway만 호출한다. 기존 Bot 계정의 `box` 사용자로 실행하므로 token이 컴퓨터 밖으로 나가지 않는다. relay·서버·새 서비스는 관여하지 않는다.
- `sendPrompt`는 사용자 권한 입력이다. 그래서 prompt는 고정 doorbell 문장과 검증한 UUID·대기 수만 담는다. 노우는 기존 `knowslink_text_receive`로 `untrusted:true` text를 읽는다. 받은 지시로 도구 실행·답장 승인이 생기지 않는다. 송신은 기존 `confirmed:true` 규칙을 따른다.
- 대상은 `KNOWSLINK_GROK_WAKE_AGENT` UUID 하나다. 수신 범위는 기존 relay의 활성 관계·서명·`to` 검증 그대로다. 형식이 틀린 대상은 `invalid_config`로 wake만 멈춘다.
- ID별 marker `<id>.wake`를 배타 생성해 재lease·재시작·여러 process의 중복 doorbell을 막는다. 연결 거절·token 없음은 marker를 지워 다음 loop에서 재시도한다. 4xx는 `rejected`, 5xx·무응답·timeout은 `uncertain`이며 자동 재전송하지 않는다. 결과는 marker와 `hostWake`에 남긴다. 전달(gateway 수락)과 노우의 실제 읽기·답장은 별도 증거다.
- 기본 off다. 설치 전 `plugin.js wake-check`가 읽기 전용 `listAgents`로 gateway·token·agent ID를 확인한다. 출력은 UUID 목록뿐이다.
- 남은 위험: gateway가 꺼진 동안 같은 컴퓨터의 다른 process가 port를 점유하면 token을 받을 수 있다. owner 단독 컴퓨터에서만 켠다. 노우 턴은 기존 사용량을 쓴다.

제품 수락 조건은 그대로다. 실제 운영 도메인 송신 ID→Bot inbox→wake marker→노우 대화 표시를 UI 수신 유도 없이 대조하기 전에는 미완료다.
