# SAR-ORCA-MCP-SEND-001-DEV QA 보고

- 작성: 2026-11-28, dev
- 테스트 레벨: lite
- 결론: 원인은 **이 PC Orca 설치본이 옛 빌드라 고정 양식 발송 TTL을 170초로 요청한 것**이다. 로컬 MCP 경로와 CLI 경로는 같은 코드(`orcaSend` → `TextTransport.sendKind`)이고 차이가 없다. 수정 코드는 이미 HEAD에 있다(e48c4206, SAR-MSG-TTL-1H-001). 이번 과제는 재발 방지 표적 테스트를 붙이고 설치본을 갱신했다. 수정 위치: adapter(설치본 갱신). relay 변경·운영 배포는 필요 없다.

## 근거

### 1. 세 발송의 relay receipt (쿠어 agent의 `knowslink_text_receipt`, 본문 없음, 18:06 KST 조회)

| id | 경로 | accepted_at (UTC) | exp (UTC) | 수명 | transport |
|---|---|---|---|---|---|
| `01a220ef-…` | MCP(15:57) | — | — | — | 조회 불가(`sender_not_allowed`, 만료 뒤 정리) |
| `01a2215e-…` | MCP(17:56) | 08:56:31.006 | 08:59:20 | **약 169초** | `failed:expired` |
| `01a22164-…` | CLI(18:02) | 09:02:12.443 | 09:05:02 | **약 170초** | `delivered`, completion `received` |

- 2차와 3차의 수명이 같다. 3차는 Grok이 170초 안에 receive해서 성공했다. 1·2차는 그 안에 받지 못해 만료됐다.
- 그러니 "로컬 MCP만 실패"는 경로 차이가 아니라 수신 시점 차이다. 14:11 item 13 왕복도 같은 이유로 성공했다.

### 2. 설치본 번들의 발송 TTL

`%LOCALAPPDATA%\KnowsLink\orca\dist\plugin.js`(15:06 빌드, 8272fb91 기준)의 `sendKind`:

```js
return this.post(peer, text, key, Date.now() + 17e4, { kind });
```

- HEAD의 `adapters/src/text.ts`는 `Date.now() + SEND_TTL_MS`(3,590,000ms)다. relay d3670f2f는 자기 agent 사이 요청을 1시간까지 받아 준다(`deliveryTTL` → `messageRetention`).
- relay는 서명된 exp를 늘리지 않고 줄이기만 한다. 그래서 운영 relay를 d3670f2f로 올려도 옛 adapter가 170초를 요청하면 170초로 남는다. 이것이 1차 진단(SAR-ORCA-DIAG-001-OPS의 170초)이 맞았는데도 d3670f2f 배포 뒤 2차가 실패한 이유다.

### 3. 로컬 상태 (읽기 전용)

- 쿠어 agent 폴더 `orca-agent-20261128`: 하나뿐이다. MCP와 CLI 모두 같은 폴더·같은 relay(`https://link.knowslog.com`)를 쓴다. 발송은 로컬에 흔적을 남기지 않는다(설계대로). inbox에는 3차 answer 수신 흔적만 있다.
- 1·2차 MCP 발송은 쿠어가 `mcp-send.mjs`로 설치본 `plugin.js`를 새 stdio 프로세스로 띄워 `knowslink_send`를 부른 것이다. to `agent_2b73608560c8ba6e8c3d06`, kind `orca.decision.v1`, 고유 idempotency_key. CLI 3차와 수신자·kind가 같다.
- 별도 관찰: 쿠어 채팅에 붙은 상시 MCP 프로세스(14:40 시작)는 14:38 빌드라 `knowslink_send`·`knowslink_receive`가 목록에 없다(세션 기록의 도구 목록 기준). 쿠어 채팅의 MCP를 다시 시작해야 고정 양식 도구가 보인다.

## 재현 단언 (표적 테스트)

`adapters/src/mcp.test.ts`의 `textBoundaries`에 `sendKind`가 보낸 envelope의 exp가 3,580,000ms 초과·`SEND_TTL_MS` 이하인지 단언을 더했다.

- 재현 전(`sendKind`만 `Date.now() + 170000`으로 바꿔 빌드): `node dist/mcp.test.js` → `MCP synthetic boundary check failed`, exit 1
- 재현 후(HEAD 코드): `PASS: MCP initialize/discovery, …`, exit 0

격리 relay 왕복 비교는 하지 않았다. 운영 receipt가 두 경로의 수명이 같음을 직접 보여 줬고, relay 쪽 1시간 clamp는 SAR-MSG-TTL-1H-001의 기존 integration이 덮는다.

## 설치본 갱신

(아래 절에 install 결과를 적는다.)
