---
title: SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW 실행 기록
status: draft
updated: 2026-10-05
owner: ops
tasks: [SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW]
summary: 만료된 시험 service token 두 개를 24h 재갱신하고 공개 검증과 Grok 전용 private 묶음을 준비한 결과 및 Grok 안내 댓글 초안을 기록한다
---

# SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW — 시험 인증 24h 갱신

## 기준과 결과 요약

착수 HEAD는 `4132354`(기준 main `31cf1a5b6ab0e9ac607f86ad06d7ad797ceeab11` 위 지시서 커밋)다. 시작 시 작업 트리는 깨끗했다. 브랜치는 `fullops/ops`, Task는 task_55c56687078d, Dispatch는 ctx_44832730a4a5다. 제품 소스는 `0911c2c`와 같다(`git diff 0911c2c --stat -- adapters internal cmd` 비어 있음). 변경한 제품 소스·배포·DB·Tunnel·DNS·owner 정책은 없다. 이전 기록(`SAR-MVP-003-BIDIRECTIONAL-OPS.md`)은 보존한다.

| 단계 | 결과 |
|---|---|
| 관리 token 확인 | `/user/tokens/verify` active, 만료 **2026-10-05T23:59:59Z**. 파일 `cf-service-token-api.env`(0600)를 Python 메모리에서만 읽었다. 값은 출력·argv·env export·Git에 남기지 않았다. |
| baseline (변경 전 GET) | service token 2개 만료(`2026-10-05T06:42:09Z`), 앱 2개(owner-only·trial), policy 3개. 이름 충돌 없음. |
| service token 갱신 | 완료. 아래 「갱신 방법」. |
| agent별 private 설정 | 변경 불필요. UUID와 secret이 그대로라 `environment.json`의 CF 두 항목이 계속 유효하다. 메타 `trial-access.json`의 만료만 갱신했다. |
| 공개 negative·positive | 통과. 아래 표. |
| Grok private 묶음 | 완료. 아래 「private archive」. |

## 갱신 방법 (CF API 근거)

Cloudflare OpenAPI(search)에는 `POST …/service_tokens/{id}/refresh`(만료를 **1년** 연장, 공식 문서 기준)와 `PUT …/service_tokens/{id}`(`duration` 지정), `POST …/rotate`(secret 교체)가 있다. 이번 범위는 24h이므로 refresh는 쓰지 않았다(1년으로 늘어나기 때문). 공식 문서의 「Renew service tokens」대로 `PUT` 본문 `{"name":…,"duration":"24h"}`를 각 token에 한 번씩 보냈다. 문서는 만료를 갱신 시각 기준으로 재설정한다고 설명한다. 실제 응답도 같다. 삭제·재생성은 필요 없었다.

| token | UUID (불변) | `client_secret_version` | 새 `expires_at` |
|---|---|---|---|
| `knowslink-trial-codex` | `283b0bdf-3c2b-4c99-b9ce-96516d226cd9` | 1 (불변) | **2026-10-06T08:09:43Z** |
| `knowslink-trial-grok` | `647039da-10bb-4459-a595-6e0caa545adb` | 1 (불변) | **2026-10-06T08:09:50Z** |

갱신 뒤 GET으로 확인했다: trial policy `884305ff-26ef-4fc5-828a-8145e7a5706c`(`non_identity`, include 두 UUID만)와 앱 `84b33961-6f78-42d3-b27b-06e064a05c2b`(`link.knowslog.com/v1/test/*`, AUD `57fc4acf…`, policy 하나, 24h)는 바뀌지 않았다. root 앱·policy·Tunnel 설정은 건드리지 않았다. 호출은 token 두 개의 PUT 두 번과 GET뿐이다.

**관리 token과 service token의 만료 차이.** 관리 token은 2026-10-05T23:59:59Z에, service token은 2026-10-06T08:09Z에 만료된다. 관리 token이 약 8시간 먼저 끝난다. 종료 시 token revoke와 앱·policy 삭제에는 그 전에 정리하거나 새 Access Edit 권한(또는 대시보드)이 필요하다.

## 공개 HTTPS 검증 (link.knowslog.com, 값 없음)

| 구분 | 요청 | 결과 |
|---|---|---|
| negative | CF header 없음, `/v1/test/registry` | 403 |
| negative | 잘못된 CF secret | 403 |
| negative | 유효 CF + 잘못된 relay credential | 401 |
| negative | 유효 CF + relay credential 없음 | 401 |
| negative | 유효 CF + 유효 relay credential으로 `/owner`·`/v1/pair`·`/v1/owners`·`/v1/signup` | 302 → Access 로그인 |
| negative | 유효 CF + 유효 relay credential으로 업무 `/v1/registry`·`/v1/pull` | 302 |
| negative | `/v1/test/../owner`·`/v1/test%2f..%2fowner` | 302 |
| positive | 유효 CF + 유효 relay credential, `/v1/test/registry` (codex·grok 각각) | 200 · 200 |

만료됐던 token이 refresh 뒤 같은 secret으로 다시 통과함을 이 결과가 확인한다. 예상한 403/401/owner 302와 모두 일치한다.

두 local client의 합성 text 왕복이다. **실제 Grok이 아니다.** nonce는 `ebadd5413d3c`다. 업무 pull은 하지 않았다.

| 단계 | 결과 |
|---|---|
| 사전 두 쪽 receive | 둘 다 `message:null` |
| Codex send | id `01a10b1d-0da4-742b-be98-9328a0afdb2f`, to=trial_grok, transport=queued |
| Grok receive | 같은 id, from=trial_codex, text에 nonce, `untrusted:true` |
| Grok reply send | id `01a10b1d-1b60-7308-9cea-bd237b541df5`, text에 첫 id 포함 |
| Codex receive | reply id 일치, nonce 일치 |
| 사후 두 쪽 receive | 둘 다 `message:null` |

TTL 180초와 업무 pull 금지를 지켰다. 제품 전체 QA·배포·loopback은 재실행하지 않았다. 공유 서비스는 불변이다: orca 200·s8 200·mcp 401, 호스트 cloudflared PID 506937, `myportfolio-*` 컨테이너 Up, knowslink relay·postgres healthy, knowslink cloudflared Up.

## private 파일과 archive (값 없음)

- 폴더 0700, 모든 `key.pem`·`environment.json`·`owners.json`·`trial-access.json` 0600, 소유자 shin, regular file, 심링크 0건.
- `grok-export/environment.json`의 `AGENT_KEY_FILE`은 `/workspace/.knowslink-trial/trial_grok/key.pem`이다.
- archive: `/home/shin/deploy/knowslink-state/trial-SAR-MVP-003-BIDIRECTIONAL/knowslink-grok-trial-20261005.tar.gz` (0600, 539 bytes). 목록은 `key.pem`·`environment.json` 상대 basename 두 개뿐이다(`tar -tvz` 확인). 항목은 uid/gid 0, 모드 0600이다. Codex key·`owners.json`·관리 token은 들어 있지 않다. 외부로 올리지 않았다. 전달은 coor·사용자 안전 채널이 담당한다.
- 안의 CF 자격은 **2026-10-06T08:09:50Z**(grok token) 이후 쓸 수 없다.

## 시험 종료 절차 (변경 없음)

1. `knowslink-trial-codex`·`knowslink-trial-grok` 두 token만 revoke(삭제)한다. 앱 `84b33961-6f78-42d3-b27b-06e064a05c2b`와 policy `884305ff-26ef-4fc5-828a-8145e7a5706c`만 삭제한다. 관리 token은 위 시각에 먼저 만료되므로 그 전에 하거나 새 권한을 쓴다.
2. 원점 config를 `tunnel-bak-pre-trial/config.yml`에서 `tunnel/config.yml`로 복원(0600)하고 knowslink connector만 재기동한다.
3. 상태 `.env`의 `KNOWSLINK_TEST_AGENTS`를 비우고 relay를 같은 Compose 인자로 갱신한다.
4. archive와 `grok-export/`, 상태 폴더의 시험 파일은 사용자 확인 뒤 삭제한다. 사용자 자료·기존 owner·공유 자원은 보존한다.

## Grok 안내 댓글 초안 — coor 검토 후 게시

이 Dispatch는 issue에 게시하지 않았다. 이전 댓글 5977433403의 만료 안내(2026-10-05T06:42:09Z)는 아래로 대체한다. 실제 Grok 검증은 없다. 자동 wake도 없다.

> ## 시험 인증 갱신 — 새 만료 2026-10-06T08:09:50Z
>
> Access service token 두 개를 24시간 갱신했습니다. 앞선 파일의 값은 그대로 유효합니다. 이미 받은 파일이 있으면 다시 받을 필요가 없습니다. 아직 받지 않았다면 운영자가 사용자 소유 안전 채널로 `knowslink-grok-trial-20261005.tar.gz` 한 개를 전달합니다. 이슈·댓글·chat에 올리지 않습니다. 안에는 `key.pem`과 `environment.json` 두 파일만 있습니다.
>
> ```sh
> install -d -m 0700 /workspace/.knowslink-trial/trial_grok
> tar -xzf knowslink-grok-trial-20261005.tar.gz -C /workspace/.knowslink-trial/trial_grok --no-same-owner --no-same-permissions
> chown "$(id -u):$(id -g)" /workspace/.knowslink-trial/trial_grok/key.pem /workspace/.knowslink-trial/trial_grok/environment.json
> chmod 0600 /workspace/.knowslink-trial/trial_grok/key.pem /workspace/.knowslink-trial/trial_grok/environment.json
> stat -c '%a %U %n' /workspace/.knowslink-trial/trial_grok /workspace/.knowslink-trial/trial_grok/*
> rm -f knowslink-grok-trial-20261005.tar.gz
> ```
>
> **준비 회신.** 위 `stat` 결과(모드·소유자만)와 `python3 --version`·`node` 고정 경로 존재 여부를 회신해 주세요. 값은 출력하지 마세요. 이 회신이 오면 coor가 Codex 송신을 시작합니다. 시험 메시지 TTL은 180초입니다. 준비 전에는 송신하지 않습니다. 수신·회신 절차와 실패 보고 형식은 앞선 댓글 5976760103을 따릅니다. 새 회신 key는 아래를 씁니다.
>
> ```sh
> printf '%s' 'Grok trial reply to <Codex send ID>' | python3 scripts/run_trial.py --config /workspace/.knowslink-trial/trial_grok/environment.json --node /workspace/.knowslink/node/bin/node send sar-mvp-003-grok-renew-round-1
> ```
>
> 앞선 댓글의 `sar-mvp-003-grok-actual-round-1`은 아직 쓰이지 않았지만 이번 round는 위 새 key 하나만 씁니다. 불명확한 send는 같은 key와 text로만 재확인하고 새 key로 자동 재전송하지 않습니다. 시험 종료 후 `/workspace/.knowslink-trial/trial_grok/`를 삭제해 주세요.

## 미검증과 후속

- 실제 Grok 왕복은 미검증이다. 파일 전달(coor·사용자)과 Grok 준비 회신이 선행 조건이다. TTL 180초 안에서 coor가 Codex 송신을 협업한다.
- 관리 token(2026-10-05T23:59:59Z)이 service token(2026-10-06T08:09Z)보다 먼저 만료된다. 종료 정리를 그 전에 하지 못하면 새 권한이 필요하다. 만료된 service token은 trial 경로를 거부로 되돌리지만 앱·policy는 남는다.
- 변경한 제품 소스는 없다. 변경 파일은 `.fullops-squad` 문서뿐이다.
