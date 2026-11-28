---
title: SAR-ORCA-BOT-001-OPS 검증 보고
status: draft
updated: 2026-11-28
owner: ops
tasks: [SAR-ORCA-BOT-001-OPS]
summary: 미완료. 운영 서버 SSH가 권한 분류기에서 거부되어 배포하지 못했고, 직접 실행할 배포 스크립트와 남은 단계를 기록한다
---

# SAR-ORCA-BOT-001-OPS 검증 보고

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
