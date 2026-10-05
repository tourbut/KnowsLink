---
id: D10
title: 프로그램설계서
status: review
updated: 2026-10-05
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-MVP-002-DEV, SAR-MVP-002-BOT-CATALOG-DEV, SAR-MVP-002-BOT-CATALOG-DEV-FIX, SAR-MVP-003-BIDIRECTIONAL, SAR-PUBLIC-IDENTITY-001-DEV, SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX]
upstream: [D02]
summary: 실제 프로그램 책임과 요구사항 및 검증을 연결한다
---

# KnowsLink 프로그램 책임과 검증

상위 요구사항은 [D02 SAR-MVP](../planning/product-specs/SAR-MVP.md)다.
API 정본은 [D05](interface-design.md), 저장 정본은 [D06/D07/D09](data-model.md)다.

| 모듈 | 책임·외부 경계 | 직접 검증 |
|---|---|---|
| protocol | Strict·Parse·SigningBytes·Digest, registry seed schema | TestFrozenParsingAndSigning |
| registry | compiled seed SHA256·manifest 공개 | protocol test·TS once |
| store | Postgres transaction·epoch·clock·current·sweep·ingest·lease | TestPostgresSafety |
| HTTP owner | credential·PoP·pairing·Go template·CSRF decision | PostgresSafety, GateFailureStates |
| HTTP agent | `deliver:agent` 전용 pull·shared persist·ACK·claim·authorize·result | PostgresSafety, ApprovalAndResultInstanceBinding |
| cleanup | 유휴 payload·metadata 회수 | 동일 sweep 경계·runtime 기동 |
| sqlc database | pgx/v5 LockRelay·SaveRelay 생성 | generate diff, 실제 Postgres |
| TypeScript adapter | registry·signature 재검증, persist·ACK·claim·gate·denied result | synthetic.ts 실제 HTTP 검사 |
| migrate·config·health | 별도 SQL-only Up·설정·readiness 실패 | 기존 unit/race·verify-runtime |
| Compose 검증 | 고유 project·private 제품 DB·loopback 시험 DB·Tunnel OFF | verify-mvp·verify-runtime |

`make test`는 외부 DB가 없어도 protocol과 기존 회귀를 수행한다. Postgres 검사는 integration build tag로 별도 실행하며 DB 환경이 없으면 실패한다.
`make verify-mvp`는 별도 DB에서 12회 ingest 경합·8회 claim 경합·3번째 lease·late ACK·pool 재시작을 검사한다.
권한 철회·세대 교체·TTL rollback·시계 이상·CSRF·중복 gate·M.id 재사용·잘못된 endpoint·optional 결과 거부를 검사한다.
`agent_cannot_process_human_delivery`는 직접 `deliver:human` send 403과 저장된 human 전달의 agent pull·persist·ACK·claim 거부를 검사한다.
`legacy_unrouted_claim_parent_boundaries`는 경로 미기록 claim의 HTTP authorize·H·R·gate-consume 403과 경로 기록 뒤 정상 처리를 검사한다.
`TestLegacyClaimsCannotReachParentBoundaries`는 a6a10c7 State 메서드로 만든 `testdata/legacy_claims.json`을 읽어 같은 경계와 새 agent 정상 경로를 검사한다.
TypeScript 검사는 Go 서버를 통해 policy 없음 deny와 gate approve 후 deny 및 최소 R 수신을 끝까지 수행한다.

독립 QA는 tester가 고정 후보 SHA에서 QA-01–11을 수행한다.
직접 시각 검수는 designer가 V-01–04를 수행한다. DEV의 자동 HTML 검사는 독립 시각 검수를 대체하지 않는다.
별도 fixed-SHA 코드 리뷰와 critical/high 차단은 coordinator가 담당한다.

## SAR-MVP-002 플러그인 모듈과 검사

- `adapters/src/core.ts`: 기존 Adapter와 signing·canonical 함수를 재사용한다. import 시 CLI를 실행하지 않는다. localAdapter는 loopback 구성과 키 파일을 읽는다. redirect를 따라가지 않는다.
- `adapters/src/index.ts`: 기존 CLI 진입점을 유지하고 core export를 제공한다. unconfigured 출력과 합성 직접 실행의 기존 계약을 보존한다.
- `adapters/src/mcp.ts`: 공식 SDK의 McpServer·StdioServerTransport·registerTool을 사용한다. 기본 held와 synthetic-only 경계를 집행하며 stdout에는 MCP만 기록한다. gate 로그를 model/tool output으로 전달하지 않는다.
- `adapters/.cursor-plugin/plugin.json`, `mcp.json`, `skills/knowslink/SKILL.md`: 공식 Cursor plugin 구조로 MCP와 사용 안내를 연결한다. 설치 기본값은 held다.
- `scripts/package_plugin.py`: 허용 목록으로 standalone marketplace ZIP을 만든다. bundle과 manifest·skill·설치 문서만 포함한다. 고정 ZIP timestamp로 같은 내용의 SHA256을 유지한다.
- `adapters/src/mcp.test.ts`: 실제 bundle의 MCP handshake·discovery·held/no-network·잘못된 모드/URL·redirect 거부를 검사한다. 별도 압축 해제 artifact 경로도 검사할 수 있다.
- `adapters/src/synthetic.ts`: 기존 SQL 합성 흐름에서 owner gate 경로를 실제 MCP bundle 호출로 검증한다. PEM은 자기 임시 0700 폴더의 0600 파일에 두고 처리 후 제거한다.

[플러그인 설치 문서](../../../adapters/README.md)와 [공식 조사·실행 기록](../exec-plans/phases/SAR-MVP-002-DEV.md)에 버전·실패 수정·한계·후속 담당을 기록한다. TESTER 독립 QA와 fixed-SHA 독립 리뷰는 coor 후속이며 이번 자동 검사로 대체하지 않는다.

## SAR-MVP-002-BOT-CATALOG-DEV 앱 등록 경로

- 관측: 재시험에서 CLI 설치·doctor는 성공했고 앱 카탈로그에는 knowslink가 없었다. 앱이 CLI plugin(`~/.grok`)을 읽지 않는다는 원인은 미확정 가설이다.
- KnowsLink는 custom MCP server **Command** 등록을 시도한다. 공식 근거는 Team Bots 문서의 Setup → Plugins → Add·채팅 요청과 Plugins 표에 한정된다. 개인 계정 UI·승인 카드는 미확인이다. 이슈1에서 Bot은 `AddMcpServer`를 호출할 수 없었다.
- `scripts/install_bot_mcp.sh`: Linux x86_64/aarch64용 고정 SHA256 Node `v22.22.2` `.tar.gz`와 bundle을 `/workspace/.knowslink`에 준비한다. 새 `.stage.*` 폴더에서 체크섬·빈 환경 MCP 경계 검사를 통과한 뒤에만 `node`·`knowslink`를 교체한다. 실패하면 기존 준비물을 보존한다. 상대 `KNOWSLINK_PREFIX`와 스크립트가 만들지 않은 `node`·`knowslink` 항목은 변경 전에 exit 1로 거절한다. PREFIX의 다른 파일은 건드리지 않는다. 환경 변수가 없으므로 `mcp.ts`의 기본 held를 사용한다.
- `scripts/verify_grok_plugin.py`: CLI 설치 검사다. 앱 카탈로그 증거가 아니다. 실패 시 grok 출력을 표시하고 `GROK_CONFIG*` 변수를 제거한다.

근거·가설·검증은 [실행 기록](../exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md)과 리뷰 보완 [FIX 기록](../exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV-FIX.md)을 따른다. 원인 판정은 FIX 기록이 우선한다.

## SAR-MVP-003 시험 transport 모듈

- `internal/relay/test_messages.go`: 두 시험 identity allowlist와 machine prefix 허용 경로를 집행한다. normal mux에 넘기는 request copy에서 owner API 진입을 차단한다.
- `protocol.go`/`registry.json`: 시험 text closed schema와 별도 registry revision을 추가한다. 기존 업무 intent는 유지한다.
- `store.go`/`http.go`: current-auth·서명·active pair·TTL·중복을 재사용한다. 시험 전용 pull·persist·ACK·claim을 집행하고 claim 시 payload를 삭제한다. rollback도 동일 allowlist를 적용한다.
- `adapters/src/core.ts`: 서명과 registry 검증을 synthetic/test receive가 공유한다. redirect 차단·요청 timeout·응답 스트림 상한을 공통 request에 적용한다.
- `test-transport.ts`/`mcp.ts`: configured peer로만 trial send한다. 수신 text는 claim 뒤 untrusted data로 반환한다. 기존 업무 pull 도구는 synthetic 모드를 유지한다.
- `trial-cli.ts`/`scripts/run_trial.py`: Codex 송신 stdin·수신 CLI와 private env 로딩을 제공한다. mode·credential은 tool/argv로 받지 않는다. launcher는 0600 소유 regular config를 요구하고 이전 agent env를 상속하지 않는다.
- `trial-setup.ts`: local beta operator만 두 신규 시험 owner/agent·Ed25519 PoP·pair를 등록한다. private 목적지 재사용을 거부한다. agent별 키와 config를 분리하고 owner 기록은 운영자에게만 둔다.
- `deploy/knowslink/access_trial_plan.py`: 단기 distinct service token·reusable non_identity policy·path 앱 본문과 보호된 Tunnel 후보를 생성한다. live mutation은 수행하지 않는다.

검증은 TestTrialMessageSafety·TestTrialSchemaAndConfiguration·TestTrialHTTP, trial-boundaries.test.ts, trial-check.ts로 연결된다. `make verify-mvp`는 기존 업무/gate 회귀와 두 독립 MCP 프로세스의 실제 Postgres 왕복을 함께 검사한다. `make verify-grok-plugin`은 네 tool을 기대한다. default held·실제 계정 수락의 분리는 계속 유지한다. [실행 기록](../exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md)에 ID·명령·exit code·남은 실제 조건을 기록한다.

## SAR-PUBLIC-IDENTITY-001 신원 모듈

- `internal/relay/identity.go`: 이메일 정규화·마스킹, rolling 한도(`take`/`hit`), 코드 발급·검증, 회원·owner 바인딩, 세션 수명, 정리.
- `internal/relay/member.go`: 시작·확인·홈·재확인·로그아웃 화면과 HTTP 상태. 메일 발송은 lock 밖에서 수행한다.
- `internal/relay/mail.go`: 표준 `net/smtp` submission. implicit TLS와 STARTTLS를 지원한다. PLAIN 인증은 TLS 또는 localhost에서만 보낸다.
- `internal/relay/http.go`: `http.CrossOriginProtection`, 합성 가입 설정, Basic owner와 회원 세션이 공유하는 gate 화면.
- `cmd/relay/main.go`: SMTP·발신 주소·client IP header·합성 가입 설정을 읽는다. 값은 출력하지 않는다.
- `scripts/mail_sink.py`: 로컬 QA 전용 SMTP sink. 받은 메일을 0600 파일로 저장한다. 실제 발송 서비스가 아니다.

| 검사 | 내용 |
|---|---|
| `TestNormalizeEmail`, `TestRollingWindowBoundary`, `TestSendLimits` | 주소 형식·별칭·rolling 경계·all-or-nothing·거부 집계·60s/5/20/100 한도·NAT 분리·rate key 원문 미보관 |
| `TestRatePrincipalIsolation` | 한 IP·회원의 반복 거부가 다른 IP·회원의 신규·정리를 막지 않음, 자기 한도 30/40/20 이하·상한·다음, 독립 principal 합의 전체 200·정리 100 포화, rolling 회복 |
| `TestCodeVerificationAndMembers`, `TestSessionLifetime` | 오답 5회·만료·재사용·확인 전 owner 미생성·재로그인 연속성·issuer 분리·회원 100 수용량·절대/무활동 수명 |
| `TestSMTPMailer` | 설정 오류의 비밀값 미노출·이름 표기 거부·실제 SMTP 대화·연결 실패 |
| `TestEmailIdentity`(integration) | 실제 Postgres HTTP 흐름: 가입·홈·재로그인·재시작·동시 첫 가입·1회 코드·로그아웃 재사용 차단·전체 로그아웃 재확인·무활동 만료·agent credential 유지·발송 실패·429·client IP header·cross-site 403·회원 gate 격리와 결정·합성 가입 차단·한 IP/회원 flood(신규 250·정리 150) 뒤 재시작에도 자기 429 유지와 다른 source 가입 시작·다른 회원 홈·로그아웃 허용 |
