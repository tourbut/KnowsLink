---
id: D03
title: 아키텍처설계서
status: review
updated: 2026-10-06
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-MVP-002-DEV, SAR-MVP-003-BIDIRECTIONAL, SAR-PUBLIC-IDENTITY-001-DEV, SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX, SAR-PUBLIC-AGENTS-001-DEV]
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
