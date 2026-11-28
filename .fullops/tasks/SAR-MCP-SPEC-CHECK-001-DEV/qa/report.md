---
title: SAR-MCP-SPEC-CHECK-001-DEV 원격 /mcp의 MCP 2026-07-28 필수 요구 점검
status: draft
updated: 2026-11-28
owner: dev
tasks: [SAR-MCP-SPEC-CHECK-001-DEV, SAR-MCP-UPGRADE-001]
summary: 배포 대기 묶음(ce4109b4)을 2026-07-28 공식 스펙과 대조했다. 서버 MUST 누락 2건(필수 헤더 부재 거부, clientCapabilities 부재 거부)과 DCR application_type 거부(RFC 7591 MUST 위반)를 작은 수정으로 고쳤다. 판정은 배포 전 수정 필요이며, 수정은 이 커밋에 들어 있다.
---

# SAR-MCP-SPEC-CHECK-001-DEV 점검 보고

테스트 레벨은 lite다. 근거는 2026-11-28에 직접 읽은 1차 출처로만 잡았다. 스펙 본문은 `/docs/2026-07-28`가 아니라 `https://modelcontextprotocol.io/specification/2026-07-28/...`에 있다. 아래 표에서 쓰는 약어는 다음과 같다.

- SH: [basic/transports/streamable-http](https://modelcontextprotocol.io/specification/2026-07-28/basic/transports/streamable-http)
- BI: [basic/index](https://modelcontextprotocol.io/specification/2026-07-28/basic/index)
- VER: [basic/versioning](https://modelcontextprotocol.io/specification/2026-07-28/basic/versioning)
- MRTR: [basic/patterns/mrtr](https://modelcontextprotocol.io/specification/2026-07-28/basic/patterns/mrtr)
- CACHE: [server/utilities/caching](https://modelcontextprotocol.io/specification/2026-07-28/server/utilities/caching)
- AUTH: [basic/authorization](https://modelcontextprotocol.io/specification/2026-07-28/basic/authorization)
- REG: [authorization/client-registration](https://modelcontextprotocol.io/specification/2026-07-28/basic/authorization/client-registration)
- CL: [changelog](https://modelcontextprotocol.io/specification/2026-07-28/changelog)

SEP 페이지와 RFC 9207·9728·8414 본문은 직접 읽지 않았다. 스펙이 인용한 문구만 근거로 썼다. RFC 7591은 §2를 직접 읽었다.

## 1. 서버 쪽 MUST 표

상태 값은 충족, 누락→수정, 해당 없음, 미확인이다. 줄 번호는 이 커밋 기준이다.

### 헤더

| 요구 사항 | 근거 | 상태 | 파일:줄 |
|---|---|---|---|
| POST를 받는 단일 엔드포인트 | SH | 충족 | `internal/relay/remote_auth.go:194` |
| `Origin`이 있고 틀리면 403 | SH#security | 충족(stdlib CrossOriginProtection) | `internal/relay/http.go:213` |
| 응답 Content-Type은 `application/json` 또는 SSE | SH#sending-messages | 충족 | `internal/relay/http.go:39` |
| 알림 POST는 202(본문 없음) 또는 4xx | SH#sending-messages | 충족. `notifications/initialized`는 202이고 다른 알림은 400이다 | `internal/relay/remote_mcp.go:60-66` |
| `MCP-Protocol-Version`과 `_meta` 버전이 다르면 400 `-32020` | SH#protocol-version-header | 충족 | `internal/relay/remote_era.go:52-56` |
| modern 요청에 `Mcp-Method`가 없거나 다르면 400 `-32020`. `Mcp-Name`은 `tools/call`·`resources/read`·`prompts/get`에 필요 | SH#standard-request-headers, #server-validation("These headers are REQUIRED") | **누락→수정**. 이전에는 헤더가 없어도 받았다(`shortcut:`). 이제 modern 요청은 `Mcp-Method`가 있어야 하고 `tools/call`은 `Mcp-Name`도 있어야 한다. `resources/read`와 `prompts/get`은 404라서 비교하지 않는다 | `internal/relay/remote_era.go:50-56` |
| Base64 sentinel `Mcp-Name`은 디코드한 뒤 비교 | SH#value-encoding | 충족 | `internal/relay/remote_era.go:26-35` |
| `Mcp-Param-*` 검증 | SH#server-behavior-for-custom-headers | 해당 없음(`x-mcp-header`를 쓰는 도구 없음) | — |
| GET·DELETE는 405, `Mcp-Session-Id`는 무시 | SH#earlier-streamable-http-revisions(SHOULD) | 충족(ServeMux 405) | `internal/relay/remote_auth.go:193-194` |
| SSE 스트림으로 독립 요청을 보내지 않음 | SH#receiving-messages(MUST NOT) | 충족(SSE 응답 없음) | — |

### 버전 협상과 오류 코드

| 요구 사항 | 근거 | 상태 | 파일:줄 |
|---|---|---|---|
| `server/discover` 구현 | VER, [discover](https://modelcontextprotocol.io/specification/2026-07-28/server/discover) | 충족 | `internal/relay/remote_mcp.go:98-99` |
| 미지원 버전은 400 `-32022`, `data:{supported, requested}` | VER#protocol-version-negotiation | 충족 | `internal/relay/remote_era.go:61-64` |
| `_meta`에 `protocolVersion`과 `clientCapabilities`가 없으면 400 `-32602` | BI#_meta | **누락→수정**. `protocolVersion`은 modern 판정 기준이라 이미 있었다. `clientCapabilities`는 검사하지 않았다 | `internal/relay/remote_era.go:65-69` |
| 모르는 메서드는 404 `-32601` | SH#protocol-version-header | 충족(modern만) | `internal/relay/remote_mcp.go:68-72` |
| 선언하지 않은 클라이언트 capability가 필요하면 `-32021` | BI#_meta | 해당 없음(클라이언트 capability가 필요한 기능 없음) | — |
| 모든 결과에 `resultType` | BI#resulttype | 충족(modern 성공 응답 전부, `events/*` 포함) | `internal/relay/remote_era.go:86-101` |
| `server/discover`와 `*/list` complete 결과에 `ttlMs ≥ 0`과 `cacheScope` | CACHE | 충족(300000, `public`) | `internal/relay/remote_era.go:90-92` |
| `serverInfo`를 결과 `_meta`에 | BI#_meta(SHOULD) | 충족 | `internal/relay/remote_era.go:89` |
| 이전 요청에 기대지 않음(stateless) | BI#statelessness | 충족(매 요청이 bearer를 다시 검증) | `internal/relay/remote_mcp.go:44-49` |
| 정의되지 않은 `-32020..-32099` 코드, 제거된 `-32002`·`-32042` 금지 | BI#error-codes | 충족. 쓰는 코드는 `-32020`·`-32022`·`-32601`·`-32602`·`-32015`(Events extension 콜백 검증)다. `-32015`는 금지 범위 밖이다 | `internal/relay/remote_events.go:147` |
| 2025-xx `initialize`도 함께 서비스(dual-era) | VER#backward-compatibility(MAY) | 충족 | `internal/relay/remote_mcp.go:83-97` |

### 인증(서버·AS)

| 요구 사항 | 근거 | 상태 | 파일:줄 |
|---|---|---|---|
| Protected Resource Metadata(RFC 9728), `authorization_servers` 1개 이상 | AUTH/authorization-server-discovery | 충족 | `internal/relay/remote_auth.go:184-185,204` |
| 401 `WWW-Authenticate`에 `resource_metadata`(+`scope` SHOULD) | 같은 곳 | 충족 | `internal/relay/remote_mcp.go:55-58` |
| AS metadata(RFC 8414), `code_challenge_methods_supported` | AUTH, security-considerations | 충족(`S256`) | `internal/relay/remote_auth.go:206` |
| AS 엔드포인트는 HTTPS | security-considerations | 충족(`remoteOrigin`은 https만) | `internal/relay/remote_auth.go:146-155` |
| redirect URI는 정확히 일치. localhost 아니면 HTTPS | security-considerations | 충족 | `internal/relay/remote_auth.go:177-181,272` |
| 공개 클라이언트의 refresh token 회전 | security-considerations | 충족(재사용 시 family 폐기) | `internal/relay/remote_auth.go:464-485` |
| 토큰 audience 검증, 만료·무효면 401 | AUTH | 충족(`Resource` 일치) | `internal/relay/remote_auth.go:81-91` |
| 토큰 passthrough 금지 | security-considerations(MUST NOT) | 충족(외부로 토큰을 보내지 않음) | — |
| `iss`를 내면 `authorization_response_iss_parameter_supported: true`를 광고 | AUTH(MUST) | 충족 | `internal/relay/remote_auth.go:206,337` |
| DCR 요청의 모르는 메타데이터는 무시 | [RFC 7591 §2](https://www.rfc-editor.org/rfc/rfc7591#section-2) "The authorization server MUST ignore any client metadata sent by the client that it does not understand" | **누락→부분 수정**. 2026-07-28 클라이언트가 보내야 하는(REG, MUST) `application_type`을 엄격 디코더가 400으로 거부했다. 이 필드는 받게 했다. 다른 모르는 필드는 여전히 거부한다(아래 남은 편차) | `internal/relay/remote_auth.go:229-230` |

### 서버 주도 요청과 deprecated 기능

| 요구 사항 | 근거 | 상태 | 파일:줄 |
|---|---|---|---|
| `roots/list`·`sampling/createMessage`·`elicitation/create`는 MRTR로만 | MRTR(MUST) | 충족(서버 주도 요청 없음, 아래 3절) | — |
| `InputRequiredResult`는 `tools/call`·`resources/read`·`prompts/get`에만, `requestState`는 무결성 보호 | MRTR(MUST) | 해당 없음(MRTR을 쓰지 않음) | — |
| Roots·Sampling·Logging은 deprecated. 새 구현은 채택하지 않는다(SHOULD NOT) | [deprecated](https://modelcontextprotocol.io/specification/2026-07-28/deprecated) | 충족(셋 다 구현 없음) | — |
| `notifications/message`는 요청 `_meta.logLevel`이 있을 때만 | server/utilities/logging(MUST NOT) | 해당 없음 | — |
| `ping`·`logging/setLevel` 제거 | CL major #5 | 참고. modern 경로가 `ping`에 `{}`로 답한다. 금지 조항이 없어 그대로 두었다 | `internal/relay/remote_era.go:82` |

## 2. 인증 판정

| 질문 | 판정 | 근거 |
|---|---|---|
| CIMD 없이 2026-07-28 클라이언트가 연결에 실패하는가 | **스펙상 실패하지 않을 수 있다.** CIMD는 AS·클라이언트 모두 SHOULD다. DCR은 deprecated(최소 12개월 유지)이고 MAY다. 클라이언트의 DCR fallback도 MAY다. CIMD만 하는 클라이언트도 적합하므로, 그런 클라이언트는 DCR만 하는 KnowsLink에 붙지 못한다 | REG "SHOULD use the following priority order… 3. Dynamic Client Registration as a fallback", "MAY fall back to Dynamic Client Registration"; CL "remains available for backwards compatibility" |
| DCR fallback은 MUST인가 SHOULD인가 | 둘 다 아니다. 클라이언트 우선순위가 SHOULD이고 fallback은 MAY다 | REG |
| 인가 응답의 `iss`는 서버 MUST인가 | **SHOULD**(오류 응답 포함). 내면 metadata 광고가 MUST다. 클라이언트는 RFC 9207 §2.4 검증이 MUST다(광고했는데 `iss`가 없으면 거부). 다음 개정에서 MUST로 올릴 예정이라고 적혀 있다 | AUTH "SHOULD include the iss parameter in authorization responses, including error responses" |
| 지금 응답에 `iss`가 있는가 | **있다.** 승인과 거부(`access_denied`) 리다이렉트 모두에 `iss`가 붙는다. metadata는 `true`를 광고한다. 기존 integration 테스트(`remote_integration_test.go:126`, `remote_browser_integration_test.go:98`)가 단언한다 | `internal/relay/remote_auth.go:206,337` |
| Claude Code | CIMD와 DCR을 둘 다 지원한다. v2.1.81 changelog: "support Client ID Metadata Document (CIMD / SEP-991) for servers without Dynamic Client Registration". 그래서 DCR이 있는 KnowsLink에는 DCR로 붙는다고 본다. 둘 다 광고할 때의 선택은 미확인. v2 runtime(2026-07-28)은 예상과 다른 issuer를 거부한다 | [code.claude.com/docs/en/mcp](https://code.claude.com/docs/en/mcp), [CHANGELOG](https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md) |
| Claude Code v2의 DCR 본문 | MCP TS SDK client `auth.ts`의 `resolveClientMetadata`가 `application_type: clientMetadata.application_type ?? deriveApplicationType(redirect_uris)`를 넣는다. 요약 발췌라 원문 줄은 미확인이다. 그래서 수정 전 KnowsLink는 이 등록을 400으로 거부했을 가능성이 높다(위 표) | [typescript-sdk auth.ts](https://github.com/modelcontextprotocol/typescript-sdk/blob/main/packages/client/src/client/auth.ts) |
| Grok | **미확인.** xAI 문서는 "complete any required authentication"과 API `authorization` 토큰 파라미터만 적는다. DCR·CIMD 언급이 없다. 운영 관찰(loopback redirect DCR)은 `remote_auth.go:175-176` 주석에 있다 | [docs.x.ai connectors](https://docs.x.ai/grok/connectors), [remote-mcp](https://docs.x.ai/developers/tools/remote-mcp) |
| Dots | **미확인.** 공식 문서에서 Dots라는 MCP 클라이언트를 찾지 못했다. 이름이 같은 docs.usedots.com은 MCP 서버다 | — |

## 3. MRTR 판정

- 원격 `/mcp`는 서버가 먼저 거는 요청이 없다. `elicitation`·`sampling`·`roots/list`·`notifications/message`를 보내는 코드가 없다(`grep`). 응답은 모두 단일 `application/json`이고 SSE 스트림이 없다. 그래서 MRTR MUST("Servers MUST send server-to-client requests… using the MRTR pattern")를 위반하지 않는다.
- 로컬 stdio(`adapters/src/mcp.ts`, SDK 1.32, 2025-11-25)는 Codex 세션에서 form elicitation을 서버 주도 요청으로 보낸다. 2025-11-25 era에서는 정상이다. 로컬 SDK를 2.x(2026-07-28)로 올리면 이 흐름을 `tools/call`의 `InputRequiredResult`로 바꿔야 한다. `requestState`를 쓰면 무결성 보호(HMAC/AEAD)도 MUST다. 이 내용을 analysis.md 3순위에 붙였다.
- 5순위 Elicitation 판단 수정: "원격에 없음"은 MUST 누락이 아니다. 원격에 넣을 때는 MRTR만 허용되고 `inputRequests`는 클라이언트가 선언한 capability에만 보낸다. 순위는 그대로 5다.

## 4. 고친 것

| 수정 | 파일 | 기존 클라이언트 영향 | 되돌리기 |
|---|---|---|---|
| modern 요청에 `Mcp-Method`가 없거나 `tools/call`에 `Mcp-Name`이 없으면 400 `-32020` | `internal/relay/remote_era.go` | legacy(`_meta` 버전 없음)는 영향 없다. 스펙 적합 modern 클라이언트는 이 헤더를 REQUIRED로 보낸다 | `remote_era.go`의 `missing` 조건을 지운다 |
| modern 요청에 `_meta.clientCapabilities`가 없으면 400 `-32602` | `internal/relay/remote_era.go` | legacy는 영향 없다 | 해당 `if` 블록을 지운다 |
| DCR이 `application_type`을 받는다(값은 무시) | `internal/relay/remote_auth.go` | 받는 범위만 넓어진다. 거부되던 등록이 통과한다. 인가 흐름은 바뀌지 않는다 | 필드 한 줄을 지운다 |

### 검증

| 검사 | 명령 | 결과 |
|---|---|---|
| 헤더·capability·DCR 단위(DB 없이) | `go test -count=1 -run 'TestRemoteModernHeaderChecks\|TestRemoteUnavailable\|TestRemoteRegisterApplicationType' ./internal/relay/` | ok |
| relay 단위 전체와 vet | `go vet ./internal/relay/ && go test -count=1 ./internal/relay/` | ok |

`TestRemoteModernHeaderChecks`에 세 경우를 더했다. `Mcp-Method` 없음, `tools/call`의 `Mcp-Name` 없음은 `-32020`이다. `clientCapabilities` 없음은 `-32602`다. 기존 modern 픽스처에는 `clientCapabilities`를 넣었다. `TestRemoteRegisterApplicationType`은 `native`·`web` 등록이 검증을 통과해 죽은 DB에서 503이 되는 것을 본다(거부되면 400). integration(`-tags=integration`)은 lite라서 다시 돌리지 않았다. 그 테스트의 modern 호출은 이미 `Mcp-Method`·`Mcp-Name`·`clientCapabilities`를 모두 보낸다(`remote_integration_test.go:247-249`).

## 5. 결론

**배포 판정: 배포 전 수정 필요.** 필요한 수정은 모두 이 커밋에 있다. 이 커밋을 병합한 뒤 묶음을 배포하면 된다.

- 배포 대기 묶음(`ce4109b4`)에는 서버 MUST 누락 2건이 있었다. modern 필수 헤더가 없는 요청을 받았고, `clientCapabilities` 부재를 검사하지 않았다. 둘 다 고쳤다.
- DCR의 `application_type` 거부는 묶음 이전부터 운영에 있던 문제다. 2026-07-28 SDK 클라이언트(Claude Code v2 runtime)의 등록이 400으로 실패했을 가능성이 높다. 받게 고쳤다. 실제 Claude Code v2 왕복은 하지 않았다.
- 인증의 나머지 MUST(PRM, AS metadata, PKCE S256, audience, refresh 회전, `iss` 광고)는 충족한다.

### 남은 편차와 새 과제 후보(lead 결정)

1. **CIMD 지원(SHOULD).** DCR이 deprecated다. CIMD만 하는 클라이언트는 붙지 못한다. AS metadata의 `client_id_metadata_document_supported`, 문서 fetch·SSRF 방어·캐시, `client_id` 일치 검증, redirect 검증, 동의 화면의 redirect 호스트 표시가 필요하다. OAuth 흐름 변경이라 승인이 필요하다. analysis.md 4순위에 넣었다.
2. **DCR의 다른 모르는 메타데이터 거부(RFC 7591 MUST ignore).** `jwks_uri`·`tos_uri`·`policy_uri` 등도 지금은 400이다. 엄격 디코더를 바꾸는 OAuth 입력 변경이라 이번에는 `application_type`만 열었다. 작은 후속 과제 후보다.
3. **`MCP-Protocol-Version: 2026-07-28` 헤더만 있고 `_meta`가 없는 요청**은 legacy 모양으로 처리된다. `initialize`가 `2026-07-28`을 요청하면 그대로 `2026-07-28`로 답한다(`remote_mcp.go:93-96`). 스펙 위반인지는 미확인이다. 정리 후보다.
4. modern `ping`을 404로 바꿀지(제거된 메서드). 금지 조항이 없어 그대로 두었다.

## 미확인

- SEP 본문(#2858, #2468, SEP-837, SEP-991)과 RFC 9207·9728·8414 원문.
- MCP TS SDK `auth.ts`의 `application_type` 줄 원문(요약 발췌만 확인).
- 둘 다 광고하는 AS에서 Claude Code가 CIMD와 DCR 중 무엇을 고르는지. Grok·Dots의 등록 방식과 `iss` 검증 여부.
- `subscriptions/listen`·pagination·progress 페이지(쓰지 않는 기능).
- 운영 `/mcp`와 실제 Claude Code v2·Grok·Dots 계정의 왕복(배포 전).
