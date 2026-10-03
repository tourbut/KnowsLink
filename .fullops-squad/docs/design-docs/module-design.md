---
id: D10
title: 프로그램설계서
status: review
updated: 2026-10-04
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-MVP-002-DEV, SAR-MVP-002-BOT-CATALOG-DEV]
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

- Grok Bot 앱 도구 카탈로그는 계정에 등록된 connector만 읽는다. Bot 컴퓨터의 Grok CLI plugin(`~/.grok`)은 앱에 등록되지 않는다. KnowsLink는 custom MCP server **Command**로 등록한다.
- `scripts/install_bot_mcp.sh`: Linux x86_64/aarch64용 고정 SHA256 Node `v22.22.2` `.tar.gz`와 bundle을 `/workspace/.knowslink`에 준비한다. 빈 환경에서 기존 MCP 경계 검사를 실행하고 등록 값을 출력한다. 환경 변수가 없으므로 `mcp.ts`의 기본 held를 사용한다.
- `scripts/verify_grok_plugin.py`: CLI 설치 검사다. 앱 카탈로그 증거가 아니다. 실패 시 grok 출력을 표시하고 `GROK_CONFIG*` 변수를 제거한다.

근거·가설·검증은 [실행 기록](../exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md)을 따른다.
