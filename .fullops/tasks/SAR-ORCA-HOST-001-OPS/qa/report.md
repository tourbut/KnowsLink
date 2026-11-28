---
title: SAR-ORCA-HOST-001-OPS 검증 보고
status: review
updated: 2026-11-28
owner: ops
tasks: [SAR-ORCA-HOST-001-OPS]
summary: 8272fb91을 운영에 배포해 verify를 통과했고 agent·grant 유지와 이 PC Orca 설치본이 8272fb91 빌드와 같음을 확인했다
---

# SAR-ORCA-HOST-001-OPS 검증 보고

## 결과 요약

**완료.** 배포했고 복귀하지 않았다. 설치본은 갱신할 것이 없었다.

| 완료 기준 | 값 |
|---|---|
| 배포 SHA | `8272fb911e1a0cf660e80390b7e6ef4ba32e2e3c`(직전 `0899e5bd…`), `deploy-history.log` 06:36:36Z |
| 백업 | `backups/0899e5bd-20261128T063635Z.dump`, 0600, 12788 bytes, TOC 25줄, SHA256 `2619ad7b45bdf4e821f100e6baa2e3685fa462c52bb9329974d732e646e90bec` |
| rollback 태그 | `knowslink-rollback:0899e5bd-relay-orcahost001`(`sha256:8f5a8baa45ec…`), `…-migrate-orcahost001`(`sha256:6a13e672ab2b…`) |
| `verify.py local`·`public` | 통과·종료 0 |
| relay | healthy, `error` 줄 0(기동 직후와 60초 뒤) |

스크립트: [deploy-8272fb91.sh](deploy-8272fb91.sh), [state-check.sh](state-check.sh)(읽기 전용).

## 배포 전 확인

`git diff 0899e5bd 8272fb91`: `db/`·`deploy/`·`compose.yaml`·`Dockerfile`·`tunnel/` 변경 없음. 서버 코드는 `human_pages.go`·`i18n_en_human.go`의 문구와 `connect_guide.go`의 선택지 1줄(인수 내용에는 앞의 두 파일만 있었다), 나머지는 테스트·`adapters/`·README다.

## verify 발췌

```text
local checks ok: health 200, unknown 404, owner/api 401, loopback-only relay, private Postgres, limits, 0600/0700
unauthenticated /healthz False -> 404
unauthenticated /owner False -> 404
unauthenticated /v1/registry False -> 404
relay_error_lines=0
head=8272fb911e1a0cf660e80390b7e6ef4ba32e2e3c
```

공개 응답: `/`·`/privacy`·`/terms`·`/.well-known/oauth-authorization-server` 200, `/home` 303.

## 웹 'Orca 봇 (내 컴퓨터)' 항목

`/home/agents?use=orca`(`&lang=en` 포함)는 로그인 없이는 303이라 curl로 화면을 볼 수 없다. 두 가지로 확인했다.

1. 실행 중인 relay container의 `/app/relay` binary 문자열 개수(`docker cp … -`, 읽기): `Orca 봇 (내 컴퓨터)` 2, `Orca bot (your computer)` 1, `수동 설치 · Windows 11 · 도어벨 자동 시작` 1, `Manual install · Windows 11 · Doorbell starts automatically` 1.
2. 같은 코드(`internal/`이 `8272fb91`과 동일)의 `human-agents` template을 합성 데이터로 렌더(`go test -overlay`, 임시 테스트는 `.devtools`, 커밋하지 않음):

```text
ko <label class="choice"><input type="radio" name="use" value="orca" required checked><span>Orca 봇 (내 컴퓨터)<br><span class="muted">수동 설치 · Windows 11 · 도어벨 자동 시작</span></span></label>
en <label class="choice"><input type="radio" name="use" value="orca" required checked><span>Orca bot (your computer)<br><span class="muted">Manual install · Windows 11 · Doorbell starts automatically</span></span></label>
```

**한계:** 로그인한 운영 화면은 보지 않았다. 주인 계정 세션을 쓰지 않았다.

## agent·grant(운영 읽기 전용, 배포 전과 후 동일)

```text
agent agent_2b73608560c8ba6e8c3d06 revoked=false   (Grok Bot)
agent agent_5513ec51367fed3748903d revoked=false   (쿠어)
grant Grok Bot -> 쿠어  orca.task.v1      만료 2026-12-28T05:08:40Z
grant Grok Bot -> 쿠어  orca.answer.v1    만료 2026-12-28T05:08:54Z
grant 쿠어 -> Grok Bot  orca.decision.v1  만료 2026-12-28T05:08:56Z
```

grant 줄의 agent ID는 이름으로 바꿔 적었다.

## 이 PC 설치본(`%LOCALAPPDATA%\KnowsLink\orca\`)

| 항목 | 값 |
|---|---|
| `dist\plugin.js` SHA256 | `eafb815ab686cb4f80626c4d7fc6ac03ccebf257d4e2bb8a3f36f072c692bd8f`. 이 worktree에서 Node 22.22.2로 새로 빌드한 것과 동일 |
| 버전 | plugin `0.1.1`, `node.exe` v22.22.2 |
| `node dist/orca.test.js` | PASS(새 빌드) |
| watch 프로세스 | 1개, pid 105428(`orca-watch.lock`의 pid와 같음), 폴더 `orca-agent-20261128`, `--bot 쿠어`. `watch.log`는 0 bytes |
| 다른 설치본 프로세스 | pid 94184, 인자 없는 `plugin.js`(봇 MCP로 본다). watch가 아니다 |
| HKCU Run `KnowsLinkOrcaWatch` | `…\orca
ode.exe …\orca\dist\plugin.js orca start` |

**선택:** bundle이 같아서 `orca install`을 실행하지 않았다(실행 중인 watch를 끊지 않는다). 갱신이 필요해지면 같은 명령을 실행한다.

## 못 한 것

- 로그인한 운영 웹 화면 확인, rollback 실행, 재부팅 뒤 자동 시작 확인(사람 단계), Grok으로의 decision 실발송.
- `npm ci`가 보고한 audit 항목은 확인하지 않았다(범위 밖).
