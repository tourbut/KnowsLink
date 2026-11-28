---
title: SAR-EVENTS-DEPLOY-001-OPS 검증 보고
status: review
updated: 2026-11-28
owner: ops
tasks: [SAR-EVENTS-DEPLOY-001-OPS]
summary: KNOWSLINK_SEAL_KEY를 운영 env에 넣고 65a752fe를 배포해 verify, 기동 로그, agent·grant 유지를 확인했다
---

# SAR-EVENTS-DEPLOY-001-OPS 검증 보고

## 결과 요약

**배포 완료, 복귀 없음.** 배포 1시간 뒤 평문 확인은 **미확인**이다. 직전 값(11:09:45Z)은 구독 0건, 평문 0건이다.

| 완료 기준 | 값 |
|---|---|
| 배포 SHA | `65a752fe2e88e0da1cbf1a844cb0000c51595295`(직전 `d3670f2f…`), `deploy-history.log` 2026-11-28T11:05:49Z |
| relay 기동 | 11:08:14Z, healthy, restart 0, `error` 줄 0 |
| 백업(배포 전) | `backups/d3670f2f-20261128T110546Z.dump`, 0600, 12846 bytes, SHA256 `691d66b292d8cba400dede667d9f5a19914254e576df3727a0565d76e976442b` |
| 백업(`beta.sh deploy`) | `backups/d3670f2f-20261128T110548Z.dump`, 0600, 12847 bytes, SHA256 `9ca7db5e5ac658530c98580ed26e7c6db657ed4448bef349ca352c4ab76cc016` |
| `.env` 백업 | `eventsdeploy001-20261128/env.pre-eventsdeploy001`, 0600, 450 bytes, SHA256 `de221b94f2bddc66da85b9b320bc8d3742e607fb52eecdbf991c6aecfb8d8d99` |
| rollback 태그 | `knowslink-rollback:d3670f2f-relay-eventsdeploy001`(`sha256:7dc071345a2d…`), `…-migrate-eventsdeploy001`(`sha256:7322832c39b9…`) |
| 키 주입 | `KNOWSLINK_SEAL_KEY` 1줄, 값 길이 44(base64 32바이트). `.env` 0600, 450→514 bytes, 새 줄을 뺀 내용은 백업과 byte 동일 |

스크립트: [deploy-65a752fe.sh](deploy-65a752fe.sh), [state-check.sh](state-check.sh)(읽기 전용). 키는 서버 스크립트의 변수에서 `printf` builtin으로 append했다. argv, 출력, 로그, 이 PC에 값이 없다.

## 배포 전 확인

`git diff d3670f2f 65a752fe`: `compose.yaml`은 `KNOWSLINK_SEAL_KEY` 한 줄, `deploy/knowslink/verify.py`는 redis 점검 추가, `db/`·`Dockerfile`·`tunnel/`은 변경 없음. 운영 checkout `d3670f2f`, dirty 0줄, 키 0줄.

## verify 발췌

```text
local checks ok: health 200, unknown 404, owner/api 401, loopback-only relay, private Postgres, limits, 0600/0700
redis used_memory 1.5 MiB / maxmemory 128 MiB (1%)
unauthenticated /healthz False -> 404
unauthenticated /owner False -> 404
unauthenticated /v1/registry False -> 404
unauthenticated /healthz True -> 404
unauthenticated /healthz True -> 404
relay_error_lines=0
relay_restarts=0
head=65a752fe2e88e0da1cbf1a844cb0000c51595295
```

공개(`https://link.knowslog.com`): `/`·`/privacy`·`/terms`·OAuth metadata 2개 200, `/home` 303, `/healthz`·`/owner` 404, token 없는 `POST /mcp` 401.

## 기동 로그

```text
2026/11/28 11:08:14 relay starting; vendor effects and disclosure disabled; email login configured=false; google login configured=true; synthetic signup=false; events webhook configured=true
```

## Redis(읽기 전용)

| 항목 | 배포 전 11:05Z | 배포 후 11:08Z |
|---|---|---|
| `maxmemory` / policy | 134217728 / `noeviction` | 동일 |
| `used_memory` | 1577816(1.50M), peak 1.53M | 1574776(1.50M), peak 1.53M |
| `dbsize` | 1 | 1 |

## agent·grant(배포 전과 후 동일)

```text
agent agent_2b73608560c8ba6e8c3d06 revoked=false   (Grok Bot)
agent agent_5513ec51367fed3748903d revoked=false   (쿠어)
grant Grok Bot -> 쿠어  orca.task.v1      만료 2026-12-28T05:08:40Z
grant Grok Bot -> 쿠어  orca.answer.v1    만료 2026-12-28T05:08:54Z
grant 쿠어 -> Grok Bot  orca.decision.v1  만료 2026-12-28T05:08:56Z
```

grant 줄의 agent ID는 이름으로 바꿔 적었다.

## 평문 secret

| 시각 | total | plain(`Secret`·`OldSecret`) | sealed | empty |
|---|---|---|---|---|
| 배포 전 11:05:26Z | 0 | 0 | 0 | 0 |
| 배포 후 11:08:45Z | 0 | 0 | 0 | 0 |
| 배포 후 11:09:45Z | 0 | 0 | 0 | 0 |
| 1시간 뒤(12:09Z 이후) | 미확인 | 미확인 | 미확인 | 미확인 |

배포 전부터 Events 구독이 0건이었다. 옮길 평문이 없었고, 봉인 저장(`sealed:v1:`)은 운영에서 아직 한 번도 일어나지 않았다.

### 1시간 뒤 재확인 절차(읽기 전용)

2026-11-28 12:09Z(21:09 KST) 이후에 이 폴더에서 Git Bash로 실행한다. 마지막 줄이 `plain=0`이면 통과다.

```bash
MSYS_NO_PATHCONV=1 scp state-check.sh shin@64.176.237.173:/tmp/kl-eventsdeploy001-state.sh
MSYS_NO_PATHCONV=1 ssh shin@64.176.237.173 'bash /tmp/kl-eventsdeploy001-state.sh </dev/null'
```

## rollback(미실행)

`bash /home/shin/deploy/knowslink/deploy/knowslink/beta.sh deploy d3670f2fa3c85a8d320542055188296948aa783b`. 키 줄은 남겨도 이전 코드가 읽지 않는다.

## 키 백업(추천, 미실행)

서버 밖 사본이 없다. 기존 관례는 서버 안 0600 env 보관뿐이다. 추천: 주인이 자기 터미널에서 그 줄만 읽어 이미 쓰는 비밀번호 관리자에 넣는다(절차는 D12). 지금 키를 잃으면 Events 구독이 지워지고 Dots의 다음 갱신까지(최대 1시간) 쉰다. 메시지는 잃지 않는다. 새 외부 저장소는 만들지 않았다.

## 못 한 것·남긴 것

- 배포 1시간 뒤 평문 확인(위 절차). 실제 Dots 구독 갱신과 배달(주인 단계), 인증된 `/mcp` `events/list`·`events/poll`, rollback 실행.
- 서버 `/tmp/kl-eventsdeploy001-{state,deploy}.sh`는 남겼다(비밀 없음). 원 로그는 기본 체크아웃 `.devtools/sar-events-deploy-001/`.
