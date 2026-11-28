---
title: SAR-MCP-DEPLOY-001-OPS 검증 보고
status: review
updated: 2026-11-28
owner: ops
tasks: [SAR-MCP-DEPLOY-001-OPS]
summary: d3670f2f를 운영에 배포해 verify를 통과했고 agent·grant 유지를 확인했다. Redis 128 MiB는 이론 최악값보다 작아 보고만 한다
---

# SAR-MCP-DEPLOY-001-OPS 검증 보고

## 결과 요약

**완료.** 배포했고 복귀하지 않았다.

| 완료 기준 | 값 |
|---|---|
| 배포 SHA | `d3670f2fa3c85a8d320542055188296948aa783b`(직전 `8272fb91…`), `deploy-history.log` 2026-11-28T08:50:17Z |
| 백업 | `backups/8272fb91-20261128T085016Z.dump`, 0600, 12714 bytes, SHA256 `57d7b7dc07935bece249f70b5f29e91a2c162eab47b40b367a7a18c26f064c43` |
| rollback 태그 | `knowslink-rollback:8272fb91-relay-mcpdeploy001`(`sha256:4c642a897ade…`), `…-migrate-mcpdeploy001`(`sha256:b92f1cbc8125…`) |
| `verify.py local`·`public` | 통과·종료 0 |
| relay | healthy, restart 0, `error` 줄 0(기동 직후와 약 78초 뒤) |
| `policyEffective` | 2026-11-28 = 배포일(KST·UTC 모두) |

스크립트: [deploy-d3670f2f.sh](deploy-d3670f2f.sh), [state-check.sh](state-check.sh)(읽기 전용).

## 배포 전 확인

`git diff 8272fb91 d3670f2f -- db deploy compose.yaml Dockerfile tunnel`: 변경 없음. 운영 checkout은 `8272fb91`, dirty 0줄. `d3670f2f..origin/main(92cdda44)`은 `.fullops/` 4개 파일뿐이다.

## verify 발췌

```text
local checks ok: health 200, unknown 404, owner/api 401, loopback-only relay, private Postgres, limits, 0600/0700
unauthenticated /healthz False -> 404
unauthenticated /owner False -> 404
unauthenticated /v1/registry False -> 404
unauthenticated /healthz True -> 404
unauthenticated /healthz True -> 404
relay_error_lines=0
head=d3670f2fa3c85a8d320542055188296948aa783b
```

## 운영 Redis(읽기 전용, 배포 전·후)

| 항목 | 배포 전 08:50Z | 배포 후 08:54Z |
|---|---|---|
| `maxmemory` | 134217728(128 MiB) | 동일 |
| `maxmemory-policy` | `noeviction` | 동일 |
| `used_memory` | 1572120(1.50M), peak 1.51M | 1572088(1.50M), peak 1.51M |
| `dbsize` | 0 | 0 |
| 컨테이너 memory 한도 | 192 MiB | 동일 |

[R3 용량 어림](../../SAR-MSG-TTL-1H-001-DEV/qa/report.md)의 이론 최악값(자기 agent 사이 원문, persist 후 미ACK)과 비교:

| 조건 | 원문 최악 | 128 MiB 대비 |
|---|---|---|
| member 1명 2400건/시간, ≈4.5 KiB | 약 21 MiB | 안 |
| member 1명 2400건/시간, 32 KiB | 약 150 MiB | **초과** |
| 전체 12000건/시간, ≈4.5 KiB | 약 105 MiB | 기록·사본 약 60 MiB를 더하면 **초과** |
| 전체 12000건/시간, 32 KiB | 약 750 MiB | **초과**(약 5.9배) |

**판정:** 현재 사용량(1.5 MiB)에는 여유가 크다. 이론 최악은 `maxmemory`를 넘는다. 넘으면 `noeviction`이라 relay가 `unavailable`을 답한다. Redis 설정은 바꾸지 않았다(승인 범위 밖). 후속 판단은 lead 몫이다.

## 공개 응답(`https://link.knowslog.com`)

`/`·`/privacy`·`/terms`·`/.well-known/oauth-authorization-server`·`/.well-known/oauth-protected-resource` 200, `/home` 303.

```text
oauth-authorization-server: issuer=https://link.knowslog.com, authorization/token/registration/revocation endpoint, code_challenge_methods_supported=["S256"]
POST /mcp initialize(2025-06-18), token 없음 -> 401 {"error":"invalid_token"}
  www-authenticate: Bearer resource_metadata="https://link.knowslog.com/.well-known/oauth-protected-resource", scope="knowslink.text", ...
POST /mcp, MCP-Protocol-Version: 1999-01-01 -> 400 -32022, supported=["2026-07-28","2025-11-25","2025-06-18","2025-03-26"]
시행일: 2026-11-28 / 2026-12-28부터 모든 메시지의 전달 대기가 최대 1시간으로 바뀝니다.
Effective date: 2026-11-28 / From 2026-12-28, pending delivery for every message will last at most 1 hour.
```

### legacy `initialize` 2025-06-18

**운영에서는 확인하지 못했다.** `/mcp`는 OAuth access token이 있어야 `initialize`에 답한다. token 발급은 운영 상태 쓰기와 Google 로그인이 필요해 하지 않았다.

대신 같은 코드(`.fullops/` 밖 `d3670f2f`와 diff 0)를 격리 Postgres·Redis에서 실행했다. `go test -tags=integration -v -run TestRemote ./internal/relay/` 통과. 첫 실행은 test DB에 migration을 하지 않아 `TestRemoteOAuthHTTP`가 `relation "relay_state" does not exist`로 실패했고, `cmd/migrate` 뒤 통과했다. 컨테이너·network는 지웠다.

```text
initialize {"protocolVersion":"2025-06-18",...} -> 200 {"result":{"capabilities":{"tools":{}},"protocolVersion":"2025-06-18","serverInfo":{"name":"knowslink","version":"1.0.0"}}}
2025-06-18 notifications/initialized -> 202
2025-06-18 tools/list -> 200
--- PASS: TestRemoteOAuthHTTP
```

## agent·grant(운영 읽기 전용, 배포 전과 후 동일)

```text
agent agent_2b73608560c8ba6e8c3d06 revoked=false   (Grok Bot)
agent agent_5513ec51367fed3748903d revoked=false   (쿠어)
grant Grok Bot -> 쿠어  orca.task.v1      만료 2026-12-28T05:08:40Z
grant Grok Bot -> 쿠어  orca.answer.v1    만료 2026-12-28T05:08:54Z
grant 쿠어 -> Grok Bot  orca.decision.v1  만료 2026-12-28T05:08:56Z
```

grant 줄의 agent ID는 이름으로 바꿔 적었다.

## rollback(미실행)

`bash /home/shin/deploy/knowslink/deploy/knowslink/beta.sh deploy 8272fb911e1a0cf660e80390b7e6ef4ba32e2e3c`

README의 배포 순서(relay 먼저, adapter 나중) 때문에, 3590초 adapter를 설치한 뒤에는 이전 relay가 `ttl_too_long`으로 거부한다. 지금은 adapter를 갱신하지 않았다.

## 못 한 것·남긴 것

- 운영 `/mcp`의 인증된 호출, 로그인한 운영 웹 화면, rollback 실행, Grok decision 재시험(lead가 쿠어에게 요청), 이 PC Orca 설치본 갱신(범위 밖).
- 서버 `/tmp/kl-mcpdeploy001-{deploy,precheck,postcheck}.sh`는 남겼다(비밀 없음). 원 로그는 기본 체크아웃 `.devtools/sar-mcp-deploy-001/`.
- 이 PC Docker에 검증용으로 받은 `redis:8.10.2-alpine`(compose가 고정한 digest) 이미지는 남겼다.
