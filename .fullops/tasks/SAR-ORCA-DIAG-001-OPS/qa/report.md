---
title: SAR-ORCA-DIAG-001-OPS 진단 보고
status: review
updated: 2026-11-28
owner: ops
tasks: [SAR-ORCA-DIAG-001-OPS]
summary: 01a220ef decision이 Grok 받은편지함에 없던 원인은 로컬 발송 기한 170초 만료로 추정하며 운영 서버 기록은 확인하지 못했다
---

# SAR-ORCA-DIAG-001-OPS 진단 보고

## 결과 요약

**원인은 추정이다.** 로컬 발송 기한 170초(`adapters/src/text.ts:130`)가 지난 뒤 Grok이 receive를 부른 것으로 본다. 만료 뒤 receive는 항상 `messages: []`다. **운영 relay·Redis·Postgres 기록은 읽지 못했고, Grok의 receive 시각도 모른다.**

## 서버를 확인하지 못한 이유

- 읽기 전용 스크립트를 올려 실행하는 명령을 Claude Code auto mode 분류기가 거부했다(`[Unauthorized Access]`). 우회하지 않았고 서버에는 아무것도 올라가지 않았다.
- ask `msg_5dfc9b1bfca8`의 답 `msg_1c4925aed66b`: 서버 확인 없이 코드 근거로 추정 보고하고 닫는다. 서버 스크립트는 실행하지 않는다.
- [trace.sh](trace.sh)는 **미실행**이다. `bash -n` 문법 검사만 했다.
- 메시지 기록은 수락 1시간 뒤(약 07:55Z) 지워지므로 이 건은 나중에 서버에서 확인할 수 없다.

## 추적 결과(로컬 근거만)

| 항목 | 값 | 근거 |
|---|---|---|
| 발송 시각 | `2026-11-28T06:55:13.921Z`(15:55 KST) | id가 UUIDv7(`text.ts:176`)이라 앞 48bit가 발송 PC 시각 |
| 만료 | 약 `06:58:03Z` | 발송 + 170초 |
| receiver·kind | `agent_2b73608560c8ba6e8c3d06`, `orca.decision.v1` | 쿠어 보고 `msg_fd3bfd0579b5` |
| relay 접수 | 접수한 것으로 본다 | `{"state":"sent",…,"transport":"queued"}`. `queued`는 relay 접수 응답이고 grant 판정을 통과해야 나온다(`public_text.go` `ingestPublicText`) |
| Grok receive 시각 | 모름. `07:00:19Z` 이전 | 발송 보고 `msg_d8dfc994cb8c` 06:55:22Z("Grok 쪽 수신은 성철 님이 확인"), 빈 결과 보고 `msg_fd3bfd0579b5` 07:00:19Z |
| grant 사용 수, 배달 상태(대기·ACK·만료) | **확인 못 함** | 서버 읽기 필요 |
| Grok receive 호출 기록 | **확인 못 함** | 코드에서 relay의 요청별 access log를 찾지 못했다 |

## 원인과 근거

1. 로컬 발송(`orca-mcp.ts` → `orcaSend` → `sendKind`)은 기한이 170초다. relay 상한은 180초다(`public_text.go:17`).
2. 원격 `knowslink_receive`는 만료 전 메시지만 돌려주고, 서버 Inbox 사본도 만료 때 지운다(`remote_auth.go` `sweepRemote`).
3. 재현([expiry-repro_test.go.txt](expiry-repro_test.go.txt), 메모리 상태, DB·네트워크 없음, `go test -overlay`로 실행, PASS):

```text
first receive at 169s | at=2m49s messages=1 transport=delivered
first receive at 170s | at=2m50s messages=0 transport=failed:expired
receive at 60s, again at 169s and 171s, no ack | at=1m0s messages=1 transport=delivered
receive at 60s, again at 169s and 171s, no ack | at=2m49s messages=1 transport=delivered
receive at 60s, again at 169s and 171s, no ack | at=2m51s messages=0 transport=delivered
receive with kinds argument: invalid_schema
```

4. **item 13이 성공한 이유(같은 근거):** decision `01a22090-c646…` 발송 05:11:28.6Z, 만료 약 05:14:18Z. Grok의 answer `01a22093-8cca…` 발송은 05:14:30.5Z(182초 뒤)다. Grok이 내용을 알았으므로 만료 전에 받았고, answer는 자기 기한을 가진 새 메시지라 늦게 보내도 통과한 것으로 본다. Grok의 수신 시각 기록은 없다.

### 다른 후보

| 후보 | 판정 |
|---|---|
| 1. Redis 1시간 보존 만료 | 아님. 빈 결과 보고가 발송 5분 뒤다 |
| 2. `kinds` 없이 불러 걸러짐 | 아님. 원격 receive는 인자를 받지 않고, 주면 `invalid_schema`다 |
| 3. 서버 보관 Inbox로 배달 안 됨 | 가능성 낮음. 같은 경로로 item 13이 왕복했다. 이 건의 서버 상태는 미확인 |
| 4. 다른 receive가 소비·ACK | **배제 못 함.** 서버 상태 필요 |
| 5. 다른 receiver 경로·grant 판정 | 가능성 낮음. relay가 `queued`로 접수했다 |

Grok Bot을 다시 연결해 다른 agent로 receive했을 가능성도 서버 확인 없이는 배제하지 못한다.

## 추천

1. **지금(사용 안내만):** 보낸 뒤 **3분 안에 받는다**(정확히는 170초). 지났으면 새 `idempotency_key`로 다시 보낸다. 재전송은 이 과제에서 하지 않았다.
2. **제품 결정 필요:** `orca.*`는 사람이 읽고 답하는 형식이라 기한을 늘릴지는 제품 규칙으로 정한다. 늘리면 바꿀 곳은 `adapters/src/text.ts`(`sendKind` 170초), `internal/relay/public_text.go`(`publicTextTTL`), `internal/relay/remote_mcp.go`(원격 send의 `exp`)와 문서의 180초 문구다. 새 DEV 과제로 둔다.
3. **문구를 코드와 맞춤:** `internal/relay/remote_mcp.go:145`의 "kept at most 1 hour"는 상한으로는 틀리지 않지만, agent 발송 text(`orca.*` 포함)는 메시지 만료(최대 180초) 때 지워진다. 읽는 쪽은 1시간 안에 받으면 된다고 이해하기 쉽다. README의 "최대 1시간" 문구도 같이 본다. **원격 `/mcp` 문구는 dev의 `SAR-MCP-UPGRADE-001` 수정 범위와 겹친다.**
4. 다음에 재발하면 발송 1시간 안에 [trace.sh](trace.sh)로 상태를 읽어 `failed:expired`(못 받고 만료)와 `delivered`(받은 적 있음)를 가른다. 먼저 SSH 허용 규칙이 필요하다.

## 못 한 것

- 운영 relay 로그·Redis·Postgres 읽기, grant 사용 수, 배달 상태, Grok receive 호출 기록 확인.
- 이 PC 설치본 로그에서 이 id의 발송 기록 확인(`watch.log`는 0 bytes, 목록만 봤다).
- 운영 쓰기·재전송·배포·재시작은 하지 않았다. D12는 원인이 운영 절차가 아니라 고치지 않았다.
