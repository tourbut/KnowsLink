---
title: SAR-ORCA-BOT-001-OPS 검증 보고
status: draft
updated: 2026-11-28
owner: ops
tasks: [SAR-ORCA-BOT-001-OPS]
summary: "미완료. 0899e5bd까지 운영에 배포했다. 에이전트 전체 철회 뒤 쿠어를 새 ID로 다시 연결했고, Grok 재연결과 grant 3개를 확인했고 왕복은 하지 않았다"
---

# SAR-ORCA-BOT-001-OPS 검증 보고

현재 상태는 마지막 [4차](#4차-2026-11-28-전체-철회-뒤-재연결) 절이다. 그 앞은 이전 기록이며, 거기 적힌 agent ID·수신기·남은 단계는 유효하지 않다.

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
- 두 번째 ask `msg_84a6c1b61bcb`(agent 확인, Grok 갱신)도 코디네이터 중지로 답이 없다(`msg_cb3c95224786`). **선택:** agent 이름은 그대로 두고, grant·`watch`·왕복은 보류한 채 미완료로 보고한다. 인박스를 닫지 않으려고 `fullops done`은 실행하지 않았다.

### Grok CLI 설치 준비 — 미사용 (계획 변경 `msg_d2da8acc2304`, 취소 `msg_4d8e0eb7e520`)

**미사용.** Grok 쪽은 이 PC의 grok CLI가 아니라 Grok Bot(자기 `/workspace`, `install_bot_mcp.sh`)이다. 아래 설치 폴더와 `grok-agent` 연결 시도(`login.json`·`private.pem`, 02:05:03Z 만료, `agent.json` 없음)는 지시대로 지우지 않았다. `grok plugin install`은 실행되지 않았고(`grok plugin list`에 knowslink 없음) agent는 만들어지지 않았다. 아래의 '리뷰 low 1의 실제 값' 중 Grok 쪽 경로 두 개는 쓰지 않는다.

- Grok은 맥북이 아니라 이 PC의 grok CLI 1.0.92다. 위 `msg_84a6c1b61bcb`의 보류 선택은 이 답으로 바뀌었다.
- 주인이 웹에서 `agent_b865fba924692fafbb27b8`의 이름을 바꿨다(`msg_b40681b4a7c0`). 새 이름은 전달받지 못했다. 기록은 ID를 쓴다.
- `python scripts/package_plugin.py --host grok --install-dir %LOCALAPPDATA%\KnowsLink\grok-plugin-b4c76b64 --connect --verify`(Node 22.22.2): MCP 경계 검사 exit 0, `grok plugin validate` 통과. `plugin.js` SHA256은 위 Orca 빌드와 같다.
- **선택:** 생성된 `.mcp.json`·`mcp.json`의 agent 폴더를 `%LOCALAPPDATA%\KnowsLink\grok-agent`로 바꿨다. `node.exe` 22.22.2를 설치 폴더에 복사해 `command`로 지정했다(`.devtools` 정리에 영향받지 않게). 되돌리려면 `grok plugin uninstall knowslink` 뒤 설치 폴더를 지운다. `~/.grok` 설정과 다른 plugin은 바꾸지 않았다.
- 리뷰 low 1의 실제 값: `<plugin.js>` = `%LOCALAPPDATA%\KnowsLink\grok-plugin-b4c76b64\knowslink\dist\plugin.js`, `<this agent folder>` = `%LOCALAPPDATA%\KnowsLink\grok-agent`, `<Orca agent>` = `agent_b865fba924692fafbb27b8`.
- 사람 단계 ask `msg_ba0c13a412f6`: `grok plugin install <설치 폴더>\knowslink --trust`, Google 동의(링크 만료 02:05:03Z).

### Grok 경로 확정: Remote HTTPS (D066, `msg_88e00f83b5dd`)

- Grok Bot은 custom connector로 `https://link.knowslog.com/mcp`(서버 보관 키)를 추가한다. 위의 Grok CLI 경로와 `install_bot_mcp.sh` 경로는 쓰지 않는다. **미사용.**
- background `setup-connect`는 둘 다 이미 끝났다(Orca는 `connected`, Grok CLI는 만료·종료 2). 남은 `setup-connect`·`orca watch` node 프로세스는 없다. 폴더는 지우지 않았다.
- dev의 `SAR-GROK-INSTALL-001-DEV`(Remote MCP의 orca 형식 송수신·grant)가 운영에 배포된 뒤 grant·왕복을 잇는다.

## 3차 (2026-11-28, `0899e5bd` 배포와 수신기)

승인: `msg_6a4182c97040`, 결정 D067(독립 리뷰 1회 뒤 push+운영 배포). 리뷰 `ba339edd34dd`는 accept·high 0, low 2(안내 문구), nit 1이다.

| 항목 | 결과 |
|---|---|
| 배포 SHA | `0899e5bd29cc11ba4ee8e70419185c97ffa536d1` (직전 `b4c76b64`), [deploy-0899e5bd.sh](deploy-0899e5bd.sh) |
| 마이그레이션 | `b4c76b64..0899e5bd`의 `db/`·`compose.yaml`·`Dockerfile`·`deploy/` diff 0 |
| 백업 | `backups/b4c76b64-20261128T024511Z.dump`, 0600, 11826 bytes, TOC 25줄, `message_state_rows=0`, SHA256 `9cf79c755427380c57c0c2957d7e119a0996555bf4d0277a891fac0ef3ca7dc5` |
| rollback 태그 | `knowslink-rollback:b4c76b64-relay-grokinstall001` (`sha256:a2c2d26ce456…`), `…-migrate-grokinstall001` (`sha256:ca941727b109…`) |
| 검증 | `verify.py local` 통과, `public` 종료 0, relay healthy·error 0줄 |
| 공개 확인 | `/`·`/privacy`·`/terms`·OAuth metadata 2종 200, 무인증 POST `/mcp` 401 |

복귀는 하지 않았다(`beta.sh deploy b4c76b646f315178852ce56650685ca64a4a343e`).

- 수신기: Orca 터미널 `term_2d169c23-9260-43e4-af27-b0283d969a64`(`--keep-open`, 제목 `KnowsLink orca watch 쿠어`)에서 `node <ops worktree>/adapters/dist/plugin.js orca watch %LOCALAPPDATA%\KnowsLink\orca-agent --bot 쿠어`. process argv에 `--bot 쿠어` 확인. `orca receive`는 `null`·종료 0.
- adapter: `plugin.js` SHA256 `22cab9ed…`. `b4c76b64..0899e5bd`의 `adapters/` 변경은 문서 2개뿐이라 다시 빌드하지 않았다.
- **선택:** 수신기는 이 ops worktree의 build를 쓴다. worktree를 지우면 수신기와 도어벨의 `receive` 명령 경로가 깨진다. 멈추려면 그 터미널을 닫는다.

### Grok 등록과 grant 확인 (읽기 전용, 본문·키·토큰 미조회)

- Grok agent 등록 확인: `agent_bce13b6103af001d83b842`(`Grok bot`, 원격 호스트, 2026-11-28T03:32:45Z), 쿠어 `agent_b865fba924692fafbb27b8`와 같은 주인. 안내 `msg_cf9b0e7cc211`.
- `msg_d1695a29f5d3`은 grant 3개와 task 발송을 알렸으나 03:55:18Z 운영 상태는 **grant 0개**(`Grants` 키 없음)다.
- task: 메시지 기록이 없다(Redis `DBSIZE` 0, 상태의 메시지 pointer 없음, 보존 1시간). accepted된 발송은 없다. Grok 쪽 `knowslink_send` 결과는 서버에 남지 않아 확인하지 못했다. grant가 없는 고정 형식 발송은 `approval_required`가 된다(코드 확인).
- 웹: 이 계정의 세션은 1개이고 마지막 요청이 03:32:45Z다. 그 뒤 이 relay에 로그인한 웹 요청은 없는 것으로 본다. `grant-create`는 Google 인증 후 5분이 지나면 재인증을 요구하고 저장하지 않는다.
- 수신기: process 실행 중, 약 10초 간격 요청이 relay에 도달한다. 받을 메시지가 없어 도어벨은 시도하지 않았다. relay error 0줄.

## 4차 (2026-11-28, 전체 철회 뒤 재연결)

- `msg_67bb5c936f6a`: 쿠어가 이 계정 agent 6개를 웹에서 철회했다고 전달됐다. 05:00:52Z 운영 읽기 확인: 6개 모두 `revoked=true`(`agent_b865…`, `agent_bce1…` 포함), grant 0개.
- 이전 수신기 터미널 `term_2d169c23-…`를 닫았다. 이전 폴더 `orca-agent`는 지우지 않았다.
- 새 폴더 `%LOCALAPPDATA%\KnowsLink\orca-agent-20261128`에서 `orca setup-connect https://link.knowslog.com <폴더>`, `connected`·종료 0. 동의는 쿠어가 했다(D065, `msg_dcc62ec31835`: 화면 지문 일치, 이름 `쿠어`). 이 worker는 동의 화면을 보지 않았다.
- 새 agent: `쿠어` `agent_5513ec51367fed3748903d`, 생성 05:02:27Z, 지문 `SHA256:Ne3GA9QHNdrLCorD-EDhSi3x9nk9NkTfXHdTB5vCCbk`. 운영 읽기 확인: 키 1개, `revoked=false`, 이전과 같은 주인.
- 수신기: 터미널 `term_fa736957-2a41-4fe4-89b6-550d9ff0ca6a`(`--keep-open`), `node <ops worktree>/adapters/dist/plugin.js orca watch %LOCALAPPDATA%\KnowsLink\orca-agent-20261128 --bot 쿠어`. 실행 중이고 오류 출력은 없다.
- grant 0개의 원인 추정(`msg_67bb5c936f6a`): grant 생성·철회는 5분 재인증이 필요하고, 재인증 없이 제출하면 '오류' 페이지만 뜨고 목록이 그대로다. 코드의 5분 제한은 확인했으나 이번에 실제로 그랬는지는 확인하지 못했다. 안내 개선은 `SAR-ORCA-GRANT-UX-001-DEV`.

### grant 확인 (05:09:29Z, 읽기 전용, `msg_d5f64993411f`)

Grok Bot `agent_2b73608560c8ba6e8c3d06`(원격 호스트, 생성 05:07:35Z)와 쿠어 `agent_5513ec51367fed3748903d`는 `revoked=false`이고 같은 주인이다. 전달받은 지문 앞자리 `38576d71`은 대조하지 않았다.

| from → to | kind | 생성(Z) | 만료 |
|---|---|---|---|
| Grok Bot → 쿠어 | `orca.task.v1` | 05:08:40 | 2026-12-28T05:08Z |
| Grok Bot → 쿠어 | `orca.answer.v1` | 05:08:54 | 2026-12-28T05:08Z |
| 쿠어 → Grok Bot | `orca.decision.v1` | 05:08:56 | 2026-12-28T05:08Z |

grant는 쿠어가 웹에서 만들었다(전달받은 내용). 시험 task 문구는 `msg_a396773da071`로 보냈다.

### 남은 단계

1. Grok 시험 task 발송(주인 요청), 도어벨 argv·`receive` 출력 기록.
2. decision 발송, answer 도어벨의 `--item` 기록.
3. D12 문구 갱신.
