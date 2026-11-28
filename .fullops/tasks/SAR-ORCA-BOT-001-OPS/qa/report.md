---
title: SAR-ORCA-BOT-001-OPS 검증 보고
status: draft
updated: 2026-11-28
owner: ops
tasks: [SAR-ORCA-BOT-001-OPS]
summary: "미완료. 2차에서 b4c76b64를 운영에 배포했고, Orca agent를 연결했다. Grok adapter 갱신 확인이 없어 grant·왕복은 하지 않았다"
---

# SAR-ORCA-BOT-001-OPS 검증 보고

현재 상태는 [2차](#2차-2026-11-28-운영-배포) 절이다. 그 앞은 1차 기록이다.

## 1차 (SSH 거부)

**결과: 미완료.** 운영 배포, 운영 agent 연결, grant 생성, Grok 왕복을 하지 않았다. 운영 서버는 바꾸지 않았다.

## 막힌 이유

- Claude Code auto mode 분류기가 로컬 diff 확인과 운영 서버 SSH 읽기를 묶은 명령을 거부했다(`[Unauthorized Access]`). SSH 부분이 원인이라는 것은 추정이다. 우회하지 않았다.
- 그 전의 SSH 읽기 1회는 실행됐다: HEAD `40aafe6ff396db2681a9423c3b344732934756e0`, relay·redis·cloudflared·postgres Up(health 표시가 있는 셋은 healthy), 백업 `40aafe6f-20261113T154024Z.dump`·`c6dfc39a-20261113T153921Z.dump`.
- 코디네이터에게 ask `msg_fcfa12c87f7b`, escalation `msg_28b3df13b05f`를 보냈다. 코디네이터는 답 없이 멈췄고 판단을 worker에게 넘겼다(`msg_3fe65d3cd2c7`).
- **선택:** 분류기 거부를 유지하고, 서버 명령은 주인이나 lead가 직접 실행하게 한다. 되돌리려면 주인이 이 worktree에 ssh/scp 권한 규칙을 추가한 뒤 같은 과제를 다시 배정한다.

## 확인한 것(로컬)

| 항목 | 결과 |
|---|---|
| `40aafe6f..b4c76b64`의 `db/`·`compose.yaml`·`Dockerfile` | 변경 없음. `deploy/`는 `test_public_config.sh`만 변경 |
| 새 웹 경로 `/home/grant-create`·`grant-revoke` | `tunnel/public-ingress.yml`의 `/home(/.*)?`가 덮는다. Tunnel 변경 불필요(live 설정 파일은 대조하지 않음) |
| 코드 복귀 | 상태에 `omitempty` 필드 `Grants`만 늘었다. 구 코드로 복귀하면 다음 저장에서 grant가 사라질 것으로 본다(구조체 diff로 본 추정, 실행 안 함) |
| adapter 빌드 | b4c76b64와 같은 `adapters/` 소스, Node 22.22.2, `npm run build` 통과, `node dist/orca.test.js` 통과 |
| `adapters/dist/plugin.js` SHA256 | `22cab9edb2a82259d7c6f9830b6dcb605e2ad250b2d3d43973c1d6c46469fc76` (이 ops worktree) |

## 배포 실행 방법(주인 또는 lead)

[deploy-b4c76b64.sh](deploy-b4c76b64.sh)는 `bash -n` 문법 검사만 했고 서버에서 실행한 적이 없다. HEAD·clean 확인, rollback 태그, `beta.sh deploy`(백업·checkout·build·`verify.py local`), 백업 해시, `verify.py public`, relay error 줄 수를 차례로 한다. `knowslink-migrate` 이미지 이름이 다르면 태그 단계에서 아무것도 바꾸기 전에 멈춘다.

이 worktree 루트의 PowerShell에서 실행한다.

```powershell
scp .fullops/tasks/SAR-ORCA-BOT-001-OPS/qa/deploy-b4c76b64.sh shin@64.176.237.173:/home/shin/deploy/knowslink-state/
ssh shin@64.176.237.173 'bash /home/shin/deploy/knowslink-state/deploy-b4c76b64.sh'
```

push, 서버 DB·`.env`·Tunnel·DNS·방화벽 변경은 하지 않았다.

## 남은 단계

1. 위 배포와 출력 기록.
2. 받는 쪽 adapter 갱신. **Grok 역할 agent의 설치 위치를 확인하지 못했다.** `%LOCALAPPDATA%\KnowsLink`에는 `claude-link-test-01-private`, `codex-google-it-*`만 있다. 양쪽을 갱신하기 전에는 grant를 만들지 않는다.
3. `orca setup-connect`로 `쿠어 (Orca)` 연결(Google 동의·지문 비교·이름은 사람 단계).
4. 웹 "내 에이전트"에서 grant 3개(사람 단계).
5. `orca watch <폴더> --bot 쿠어`와 왕복 1회, 도어벨 argv·`receive` 출력·`--item` 기록.
6. D12에 이번 적용 절 추가. 적용한 것이 없어 이번에는 고치지 않았다.

## 2차 (2026-11-28, 운영 배포)

**결과: 미완료.** 운영 배포와 Orca agent 연결은 끝났다. grant 생성과 Grok 왕복은 하지 않았다.

### 배포

| 항목 | 결과 |
|---|---|
| 배포 SHA | `b4c76b646f315178852ce56650685ca64a4a343e` (직전 `40aafe6f`), 01:47:20Z, `deploy-history.log`에 1줄 |
| 마이그레이션 | `40aafe6f..b4c76b64`의 `db/`·`compose.yaml`·`Dockerfile` diff 0 |
| 백업 | `backups/40aafe6f-20261128T014719Z.dump`, 0600, 11300 bytes, TOC 25줄, `message_state_rows=0`, SHA256 `0d3d88b125440d43a5fb09224fcbbe365ad9bbb1cb030db4870734dda7100f8b` |
| rollback 태그 | `knowslink-rollback:40aafe6f-relay-orcabot001` (`sha256:be0ff2162c71…`), `…-migrate-orcabot001` (`sha256:af8c10bf68bb…`) |
| `verify.py local` | `health 200, unknown 404, owner/api 401, loopback-only relay, private Postgres, limits, 0600/0700` |
| `verify.py public` | 종료 0. 미인증 `/healthz`·`/owner`·`/v1/registry` 404 |
| relay | healthy, error 0줄(기동 직후와 약 1분 뒤) |
| 공개 확인 | `/`·`/privacy`·`/terms`·OAuth metadata 200, `/home` 303, 미인증 POST `/home/grant-create`·`grant-revoke` 303, 없는 경로 404 |

복귀는 하지 않았다. 필요하면 서버에서 `bash /home/shin/deploy/knowslink/deploy/knowslink/beta.sh deploy 40aafe6ff396db2681a9423c3b344732934756e0`. grant를 만든 뒤 복귀하면 grant가 사라질 것으로 본다(1차 추정).

### 연결과 사람 단계

- 이 PC adapter: `npm run build`, `node dist/orca.test.js` PASS, `adapters/dist/plugin.js` SHA256 `22cab9edb2a82259d7c6f9830b6dcb605e2ad250b2d3d43973c1d6c46469fc76`.
- `orca setup-connect https://link.knowslog.com %LOCALAPPDATA%\KnowsLink\orca-agent --host orca`를 실행했다. 만료(02:01:02Z) 전에 `connected`·종료 0이 나왔다. agent `agent_b865fba924692fafbb27b8`, 이름 `Orca Bot`(요청한 `쿠어 (Orca)`와 다름), 지문 `SHA256:Im2nPfxt-GppujUkek5dXm2FPU4WSA_9nRw6dEw8i3Q`(ask에 적은 값과 같음). 누가 어떤 화면에서 동의했는지는 직접 확인하지 못했다.
- ask `msg_e71dac575ae2`(동의, 맥북 Grok adapter 갱신). 코디네이터는 답 없이 멈추고 판단을 넘겼다(`msg_7721fdcdee4c`).
- **선택:** 사람 단계를 자동화하지 않는다. Grok 갱신 확인 전이라 grant를 만들지 않았다. 연결을 되돌리려면 주인이 웹 "내 에이전트"에서 이 agent를 해제한다.

### 남은 단계

1. 주인이 `Orca Bot` agent가 본인 계정에 있는지 확인한다. 이름을 `쿠어 (Orca)`로 할지 정한다.
2. 맥북 Grok(`/workspace/.knowslink`) adapter를 `b4c76b64` 빌드로 갱신하고 agent 이름·`plugin.js` 경로·연결 폴더를 받는다.
3. 주인이 grant 3개를 만든다.
4. `orca watch <폴더> --bot 쿠어`, 왕복 1회, 도어벨 argv·`receive`·`--item` 기록.
