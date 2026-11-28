---
title: SAR-MCP-UPGRADE-001 MCP 새 기능 갭 분석과 우선순위
status: draft
updated: 2026-11-28
owner: dev
tasks: [SAR-MCP-UPGRADE-001]
summary: Skills·Tasks·Elicitation·File Uploads·2026-07-28 stateless를 원격 /mcp와 로컬 stdio 기준으로 비교하고 순위를 정했다. 1순위는 원격 /mcp의 2026-07-28 dual-era 적합성이다.
---

# SAR-MCP-UPGRADE-001 갭 분석

1차 출처만 썼다(2026-11-28 확인). 확인하지 못한 것은 '미확인'이다.

## 지금 KnowsLink 표면

| 표면 | 위치 | 프로토콜 | 비고 |
|---|---|---|---|
| 원격 `POST /mcp` | `internal/relay/remote_mcp.go` | 작업 전: legacy `initialize`(응답은 2025-11-25 또는 2026-07-28 고정) + `server/discover`가 2026-07-28을 광고. 세션·GET SSE·DELETE 없음 | OAuth(DCR·PKCE S256), 도구 4개, MCP Events(webhook). `outputSchema`·`structuredContent`·`annotations`·`title` 있음 |
| 로컬 stdio | `adapters/src/mcp.ts`, `orca-mcp.ts` | `@modelcontextprotocol/sdk` 1.32.0: 최신 2025-11-25, 지원 2024-10-07~2025-11-25. 2026-07-28 없음 | Codex 세션에서 form elicitation(일회 코드). `annotations` 있음, `structuredContent` 없음. SDK에 tasks는 experimental만 |

## 기능별 갭

| 기능 | 스펙 상태 | 원격 지금 | 로컬 지금 | 클라이언트 지원 | 하위 호환 영향 | KnowsLink 이득 | 보안 영향 |
|---|---|---|---|---|---|---|---|
| 2026-07-28 stateless | 확정(2026-07-28 정식). SEP-2567·2575·2322·2243·2549 | 부분. 이미 modern 광고 중인데 헤더 검증(-32020), 미지원 버전(-32022), 404 `-32601`, `resultType`·`ttlMs`·`cacheScope`, `_meta` serverInfo가 없음. legacy `initialize`가 2025-06-18 요청에 2025-11-25로 답함 | 없음(SDK 1.32). Claude Code v2는 stdio에 `server/discover` probe 후 legacy로 fallback | Claude Code: v2 runtime이 HTTP·stdio에서 2026-07-28을 먼저 묻는다([docs](https://code.claude.com/docs/en/mcp)). Dots: MCP Events가 `initialize`·2026-07-28 둘 다 대상([OpenAI](https://developers.openai.com/plugins/build/mcp-events)), Dots 자체 문서는 미확인. Grok: Streamable HTTP·SSE만 명시, 버전 미확인([xAI](https://docs.x.ai/docs/guides/tools/remote-mcp-tools)). Orca: 자체 클라이언트 아님, 띄운 agent를 따름([manifest](https://github.com/stablyai/orca/blob/main/docs/reference/plugin-manifest.md)) | dual-era로 하면 깨지지 않는다. legacy 응답 모양을 그대로 두고 modern 요청에만 필드를 더한다 | 이미 광고한 버전과 실제 동작을 맞춘다. modern 클라이언트가 캐시·버전 오류를 스펙대로 처리 | 헤더·body 불일치 거부는 스펙 MUST(경로마다 다른 값을 믿는 혼동 방지) |
| Skills | SEP-2640 Final, extension `io.modelcontextprotocol/skills`(beta 문서) | 없음 | 없음 | Claude Code: public beta, v2.1.331+([docs](https://code.claude.com/docs/en/mcp)). Grok·Dots·Orca 미확인 | extension 미선언 클라이언트는 영향 없음 | untrusted·발송 규칙을 skill로 배포할 수 있다. 같은 내용의 짧은 판은 이번에 `server/discover` `instructions`로 넣었다 | skill 내용은 host가 digest·size를 검증. 승인은 manifest에 묶임 |
| Tasks | 2025-11-25 core experimental(SEP-1686) → 2026-07-28에서 extension SEP-2663 Final, `io.modelcontextprotocol/tasks`. 두 판은 와이어 비호환 | 없음 | 없음(SDK experimental만) | Claude Code: 2026-07-28 연결, v2.1.304+, `tasks/get` polling. 다른 클라이언트 미확인 | 클라이언트가 extension을 선언하지 않으면 task를 돌려주면 안 됨(MUST NOT) → 영향 없음 | 낮다. 도구 4개 모두 짧은 동기 호출이다. 발송 승인 대기(`approval_required`)를 task로 만들 수 있으나 승인 규칙 판단이 필요 | taskId 엔트로피·요청별 인가 MUST. 승인 대기를 task로 열면 승인 경계를 다시 봐야 함 |
| Elicitation | 확정. form 2025-06-18, URL mode 2025-11-25(SEP-1036), 2026-07-28에선 MRTR(`InputRequiredResult`) | 없음 | form 사용(Codex 세션 bind·발송, 일회 코드 입력 요구) | Claude Code: form·URL 지원. OpenAI 플러그인 문서에 elicitation 안내(모드 미확인). Grok·Dots 미확인 | capability 미선언 클라이언트에는 보내지 않으면 영향 없음 | 낮다. 발송 승인은 웹 사전 승인·발송별 사람 승인이 정본이다 | 사람 승인을 대신하면 안 된다(YOLO 자동 수락 사례). 민감정보는 form 금지, URL mode는 서버 인가 대체 금지(스펙 MUST NOT) |
| File Uploads | 스펙·확정 SEP 없음. SEP-2356·2631 닫힘. `experimental-ext-file-uploads` 초안(breaking 예상) | 없음 | 없음 | 미확인(초안이라 대상 없음) | 해당 없음 | 없음. 제품은 4096바이트 text만 다룬다 | 도입 시 크기 상한·untrusted 처리·업로드 대상 정보를 model context에 넣지 않기(초안 MUST NOT) |
| 도구 메타(annotations·outputSchema·structuredContent·title·icons) | 확정. annotations 2025-03-26, outputSchema·structuredContent·title 2025-06-18, icons 2025-11-25 | icons 외 있음 | annotations만 | Claude Code는 readOnly·idempotent hint 사용, structuredContent 처리 미확인 | 추가 필드라 영향 없음 | 로컬 structuredContent는 작다(텍스트 JSON 이미 반환) | 없음 |

출처: [2026-07-28 changelog](https://modelcontextprotocol.io/specification/2026-07-28/changelog), [streamable HTTP](https://modelcontextprotocol.io/specification/2026-07-28/basic/transports/streamable-http), [versioning](https://modelcontextprotocol.io/specification/2026-07-28/basic/versioning), [discover](https://modelcontextprotocol.io/specification/2026-07-28/server/discover), [MRTR](https://modelcontextprotocol.io/specification/2026-07-28/basic/patterns/mrtr), [SEP-2663](https://modelcontextprotocol.io/seps/2663-tasks-extension), [2025-11-25 tasks](https://modelcontextprotocol.io/specification/2025-11-25/basic/utilities/tasks), [elicitation](https://modelcontextprotocol.io/specification/2026-07-28/client/elicitation), [SEP-2640](https://modelcontextprotocol.io/seps/2640-skills-extension), [file uploads WG](https://modelcontextprotocol.io/community/working-groups/file-uploads), [ext 초안](https://github.com/modelcontextprotocol/experimental-ext-file-uploads), [client matrix(비공식)](https://modelcontextprotocol.io/extensions/client-matrix). 헤더·오류 코드·`CacheableResult`는 `schema/2026-07-28/schema.ts` 원문으로 다시 확인했다.

## 우선순위 (추천, 잠정 진행)

1. **원격 `/mcp` 2026-07-28 dual-era 적합성.** 이미 `server/discover`로 2026-07-28을 광고해 modern 클라이언트가 이 경로로 온다. 빠진 MUST(헤더 불일치 거부, 미지원 버전 오류, `resultType`, `ttlMs`·`cacheScope`)를 채운다. legacy 응답은 바꾸지 않는다. 보안 이득(헤더 혼동 거부)이 있고 기존 클라이언트를 깨지 않는다. → 이번에 구현.
2. **Skills extension.** Claude Code beta가 있고 untrusted·발송 규칙을 한 곳에서 배포할 수 있다. `resources` capability와 `skills/list`·`skills/get`·`resources/read`가 필요하다. modern 경로가 먼저 맞아야 한다.
3. **로컬 stdio SDK 1.32 → 2.x(2026-07-28).** Claude Code v2가 stdio에서도 modern을 먼저 묻는다. 지금은 fallback으로 동작한다. 의존성 대폭 변경이라 단독 과제로 둔다. 로컬 `structuredContent`도 여기서 함께 한다.
4. **Tasks extension.** 도구가 짧아 이득이 작다. 승인 대기를 task로 열 경우 승인 규칙 판단이 먼저다.
5. **Elicitation(원격 MRTR).** 사람 승인을 대신할 수 없고 웹 승인이 정본이라 이득이 작다.
6. **File Uploads.** 스펙·SEP가 없고(실험 초안) 텍스트 제품이라 보류.

## 1순위 구현 요약

- 위치: `internal/relay/remote_mcp.go`(`remoteEra`, `remoteError`, `remoteHeader`, 결과 장식), 인터페이스 정본 [D05](../../docs/design-docs/interface-design.md)의 dual-era 절.
- 선택: modern 요청에 `Mcp-Method`·`Mcp-Name`이 없어도 받는다(헤더로 라우팅하는 계층이 없음). 값이 있고 다르면 era와 관계없이 거부한다. 되돌리기: `remoteEra`의 `method != ""`·`name != ""` 조건을 modern일 때 필수로 바꾼다(코드의 `shortcut:` 주석).
- 선택: 캐시는 `ttlMs:300000`, `cacheScope:"public"`(도구 목록·discover는 설치와 무관하게 같다). 배포로 바뀔 수 있어 5분으로 짧게 둔다.
- 하지 않은 것: `events/subscribe` 결과의 `resultType`(클라이언트는 없으면 `complete`로 본다), `subscriptions/listen`, MRTR, Base64 sentinel의 `Mcp-Param-*`(`x-mcp-header`를 쓰는 도구가 없다).
