---
title: SAR-MCP-UPGRADE-001 검증 보고
status: draft
updated: 2026-11-28
owner: dev
tasks: [SAR-MCP-UPGRADE-001]
summary: 원격 /mcp 2026-07-28 dual-era의 헤더 검사 단위 테스트와 격리 Postgres·Redis integration으로 legacy 하위 호환과 modern 결과를 확인했다.
---

# SAR-MCP-UPGRADE-001 검증 보고

테스트 레벨 lite. 분석과 순위는 [analysis.md](../analysis.md)에 있다.

## 실행

| 검사 | 명령 | 결과 |
|---|---|---|
| 헤더·버전 검사(DB 없이) | `go test -count=1 -run 'TestRemoteModernHeaderChecks\|TestRemoteUnavailable' ./internal/relay/` | ok |
| relay 단위 전체 | `go test -count=1 ./internal/relay/` | ok |
| 원격 integration | 임시 `postgres:17-alpine`·`redis:8.10.2-alpine` 컨테이너, `go run ./cmd/migrate` 후 `go test -tags=integration -count=1 -run 'Remote\|Orca' ./internal/relay/` | ok. 컨테이너와 새로 받은 redis 이미지는 실행 직후 지웠다 |

`TestRemoteModernHeaderChecks`는 통과한 요청이 죽은 DB에 닿아 503이 되고, 거부되는 요청은 DB 전에 400이 되는 것으로 구분한다. 9개 경우: legacy 무헤더·legacy 헤더·modern 일치·Base64 `Mcp-Name`은 503(통과). 버전 헤더≠`_meta`, `Mcp-Method`≠method, `Mcp-Name`≠tool은 `-32020`. `_meta`의 legacy 버전, 모르는 헤더 버전은 `-32022`.

## 하위 호환 발췌 (`go test -v -run TestRemoteOAuthHTTP$`)

새 기능을 모르는 클라이언트(`initialize`, `_meta` 없음)는 이전과 같은 모양을 받고 새 필드가 없다(`resultType`·`ttlMs`·`cacheScope` 부재를 단언한다).

```text
 {"method":"initialize","params":{"protocolVersion"... -> 200 {"result":{"capabilities":{"tools":{}},"protocolVersion":"2025-03-26","serverInfo":{"name":"knowslink","version":"1.0.0"}}}
 {"method":"initialize",...2025-06-18 -> 200 {"result":{"capabilities":{"tools":{}},"protocolVersion":"2025-06-18",...}}
 {"method":"initialize",...2025-11-25 -> 200 {"result":{...,"protocolVersion":"2025-11-25",...}}
 {"method":"initialize",...2024-11-05 -> 200 {"result":{...,"protocolVersion":"2025-11-25",...}}
2025-06-18 {"method":"notifications/initialized"} -> 202
2025-06-18 {"method":"tools/list"} -> 200 {"result":{"tools":[{"annotations":{...},"description...
 {"method":"resources/list"} -> 200 {"error":{"code":-32602,"data":{"reason":"unsupported_method"},"message":"Request refused"}}
```

같은 테스트가 이어서 legacy `tools/call knowslink_status`와 `knowslink_send`(4096바이트 escape 본문)를 이전 단언 그대로 통과한다. receive·ack 흐름은 바뀐 코드 밖(`remoteTool`)이며 `TestRemoteGoogleBindingAndSignedInbox`·`TestRemoteOrcaFormsUseOwnerGrant`가 통과한다.

## modern 발췌

```text
2026-07-28 server/discover -> 200 {"result":{"_meta":{"io.modelcontextprotocol/serverInfo":{...}},"cacheScope":"public","capabilities":{"e...
2026-07-28 tools/list      -> 200 {"result":{"_meta":{...},"cacheScope":"public","resultType":"comp...
2026-07-28 tools/call knowslink_status -> 200 {"result":{"_meta":{...},"content":[{"text":"{\"agent\":...
2026-07-28 tools/call knowslink_nope   -> 200 {"error":{"code":-32602,"data":{"reason":"unsupported_method"},...}}
2026-07-28 resources/list -> 404 {"error":{"code":-32601,"message":"Method not found"}}
2026-07-28 initialize     -> 404 {"error":{"code":-32601,"message":"Method not found"}}
```

## 못 한 것

- 실제 Claude Code v2·Dots·Grok 계정으로 운영 `/mcp`에 붙는 왕복은 하지 않았다(원격 배포 전, 승인 필요).
- `fullops lint`의 Windows 한계(make 부재)는 기존 기록과 같다. 결과는 완료 보고에 적는다.
