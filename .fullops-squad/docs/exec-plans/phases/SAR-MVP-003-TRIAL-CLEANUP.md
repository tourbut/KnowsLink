---
title: SAR-MVP-003-TRIAL-CLEANUP 실행 기록
status: draft
updated: 2026-10-05
owner: ops
tasks: [SAR-MVP-003-TRIAL-CLEANUP]
summary: 실제 Grok 왕복 확정 뒤 승인된 D12 종료 절차로 시험 인증과 앱 및 정책을 정리하고 원점 복원과 공유 서비스 회귀를 확인한 결과를 기록한다
---

# SAR-MVP-003-TRIAL-CLEANUP — 시험 종료와 원점 복원

## 기준과 범위

착수 HEAD는 `ccc17efcf251b08b6a84df7f3965ab889810396b`다. 브랜치는 `fullops/ops`, Task는 task_460d530b9673, Dispatch는 ctx_fc0f0f1c6aa7이다. 시작 시 작업 트리는 깨끗했다. 실제 Codex→Grok→Codex 왕복은 [성공 판정](SAR-MVP-003-ACTUAL-TRIAL.md)으로 확정됐다. 그 원본 기록은 바꾸지 않았다. 이번 작업은 승인된 D12 13.4 종료 절차만 실행했다. 제품 소스·DB·Tunnel 생성·DNS·초기 expose는 바꾸지 않았다. 사용자 private 파일과 기존 owner 자원은 삭제하지 않았다. 값(secret)은 출력·argv·export·Git에 남기지 않았다. 관리 token은 Python 메모리에서만 읽었다.

## 1. 변경 전 baseline (읽기 전용 GET, 11:46Z)

| 항목 | 관측 |
|---|---|
| 관리 token | `/user/tokens/verify` active, 만료 2026-10-05T23:59:59Z. 계정 범위 한정이라 `/accounts`는 빈 목록이다. 계정 ID는 Tunnel 자격 파일의 `AccountTag`로 얻었다. |
| service token | 정확히 2개: `knowslink-trial-codex` `283b0bdf-…`, `knowslink-trial-grok` `647039da-…` (UUID가 `trial-access.json`과 일치) |
| 앱 | `KnowsLink trial messages` `84b33961-…`(`link.knowslog.com/v1/test/*`, policy 하나), `KnowsLink beta (owner-only)` `bd210310-…` |
| policy | `knowslink-trial-agents` `884305ff-…`(`non_identity`, include는 위 두 token만, exclude·require 비어 있음), `knowslink-beta-owner-only` `a0d1a950-…`, `knowslog-bot - Production` `952aa38b-…` |
| root 보존 기준 | owner 앱 1개와 policy 2개의 전체 JSON을 정렬해 SHA-256으로 기록했다(앞 12자: 앱 `0966f1ba81c6`, owner policy `c3db9e309cd3`, bot policy `ee96e0bce448`). |
| Tunnel config | 현재 live는 백업 `tunnel-bak-pre-trial/config.yml`에 `/v1/test/.*` trial AUD rule(9줄)만 더한 것이다. `diff`로 확인했다. 백업은 owner AUD hostname rule 하나와 404뿐이다. |
| 배포·서비스 | 체크아웃 `0911c2c…`, relay·postgres healthy, knowslink cloudflared Up. `.env`에 `KNOWSLINK_TEST_AGENTS`가 있었다. |
| 이름·type·path·UUID | 삭제 직전에 앱 이름·domain·policy 연결, policy 이름·decision, token 이름을 다시 GET으로 대조해 모두 일치했다. 불일치는 없었다. |

## 2. 실행한 정리

순서: 로컬 원점 복원 → allowlist 비움 → Cloudflare 삭제. 로컬을 먼저 닫아 삭제 사이에 trial 경로가 열리지 않게 했다.

1. **Tunnel config 복원.** 시험 중 live config를 `tunnel-trial-active-ended/config.yml`(0700 폴더, 0600)에 보존했다. 백업 `tunnel-bak-pre-trial/config.yml`을 `tunnel/config.yml`에 0600으로 복사했다. `cmp`가 같다. 복원 전에 pinned `cloudflare/cloudflared:2026.9.1@sha256:b269e8ab…`(호스트 UID/GID)로 `ingress validate`가 `OK`였고 `ingress rule`이 `/v1/test/registry`·`/owner`·`/v1/pair`를 모두 rule #0(owner AUD hostname rule)에 매칭했다. `docker restart knowslink-cloudflared-1`로 knowslink connector만 재기동했다. precheck는 통과했고 `hard_fail=false`다. 호스트 cloudflared(PID 506937)·`myportfolio-*`는 건드리지 않았다.
2. **allowlist 비움.** 상태 `.env`에서 `KNOWSLINK_TEST_AGENTS` 줄만 제거했다. 다른 세 줄(`POSTGRES_PASSWORD`·`DATABASE_URL`·`RELAY_PORT`)은 그대로이고 모드는 0600이다. `beta.sh`의 `dc`와 같은 인자로 `docker compose -p knowslink -f compose.yaml -f deploy/knowslink/compose.ops.yaml --env-file <state>/.env up -d --no-deps relay`를 실행했다. relay만 재생성됐고 8초 뒤 healthy다. 컨테이너 안 `KNOWSLINK_TEST_AGENTS`는 빈 값이다. 변경 전 `.env` 사본은 비밀 값이 들어 있어 확인 후 내가 삭제했다.
3. **Cloudflare 삭제.** `DELETE`는 앱 `84b33961-…`, policy `884305ff-…`, token 두 개에만 보냈다. 응답은 앱 202, policy 202, token 두 개 각 200이며 모두 `success: true`다. 순서는 앱 → policy → token이다. 다른 자원에는 쓰기 호출이 없다.

## 3. 정리 후 검증

| 구분 | 결과 |
|---|---|
| service token 목록 | 빈 목록. 두 UUID `GET`은 404. |
| 앱 목록 | `KnowsLink beta (owner-only)` `bd210310-…` 하나뿐. trial 앱 `GET`은 404. |
| policy 목록 | `knowslink-beta-owner-only`·`knowslog-bot - Production` 둘뿐. trial policy `GET`은 404. |
| root 불변 | owner 앱·owner policy·bot policy의 JSON SHA-256이 baseline과 모두 같다(`unchanged` 3건). IdP·team·Tunnel·DNS는 건드리지 않았다. |
| Tunnel/config | live `tunnel/config.yml`과 백업 `cmp` 일치, 0600. |
| allowlist | 컨테이너 env 비어 있음, 상태 `.env`에 줄 없음. |
| 로컬 trial 경로 닫힘 | 무인증 loopback `/v1/test/registry`는 401이다. trial_codex key로 루프백 relay에 `receive`하면 실패한다(`KnowsLink trial failed`, exit 1). allowlist가 비어 시험 identity가 거부된다. |
| 공개 owner 경로 | CF header 없이 `/v1/test/registry`·`/owner`·`/v1/pair`·`/v1/registry`·`/v1/signup`이 모두 302 → `scshin88.cloudflareaccess.com` 로그인이다. 이전에는 `/v1/test/registry`가 403이었고 이제 owner 앱이 덮는다. |
| 이전 시험 CF header | 이전 `trial_codex`의 `CF_ACCESS_CLIENT_ID/SECRET`으로 `/v1/test/registry`·`/owner`·`/v1/pair`·`/v1/registry`를 호출하면 모두 403이다. 리다이렉트는 따르지 않았다. 값은 출력하지 않았다. |
| 공유 서비스 | `verify.py regression` exit 0: orca 200·s8 200·mcp 401, 호스트 cloudflared PID 506937. `verify.py public` exit 0(전부 302). `verify.py local` exit 0(health 200, owner/api 401, loopback-only, private Postgres). |
| 컨테이너 | `knowslink-relay-1` healthy, `knowslink-cloudflared-1` Up, `knowslink-postgres-1` healthy, `myportfolio-{worker,beat,viewer,postgres,redis}-1` Up 4 weeks 그대로. |
| 제품 변경 | 없음. 배포 체크아웃은 `0911c2c…` 그대로이고 깨끗하다. |

## 4. 사용자 private 파일과 남은 항목

- `/home/shin/deploy/knowslink-state/trial-SAR-MVP-003-BIDIRECTIONAL/`(0700)의 `trial_codex/`·`trial_grok/`·`grok-export/`·`owners.json`·`trial-access.json`·`knowslink-grok-trial-20261005.tar.gz`는 **삭제하지 않았다**. 모두 0600·소유자 shin·일반 파일이고 폴더는 0700이다. 안의 Cloudflare 자격은 이미 삭제된 token을 가리키므로 **더는 연결할 수 없다**. 정리 판단은 사용자 몫이다.
- `trial-access.json`의 token·앱·policy UUID는 이제 존재하지 않는 자원이다. 파일은 보존했고 바꾸지 않았다.
- Grok 쪽 `/workspace/.knowslink-trial/trial_grok/`와 받은 tar 삭제는 Bot에게 댓글로 요청됐다. 직접 원격 삭제하지 않았다. 이 Dispatch에는 그 삭제의 확인 근거가 없다.
- 관리 사용자 API token(`knowslink-trial-ops-24h`)은 이 token으로 무효화 endpoint를 호출할 권한이 없다. 값을 삭제하거나 새로 발급하지 않았다. **2026-10-05T23:59:59Z에 자동 만료**된다. 파일 `cf-service-token-api.env`(0600)는 사용자 파일이므로 남겼다.
- 시험 중 live config 사본 `tunnel-trial-active-ended/config.yml`(AUD 식별자만 있고 secret 없음, 0600)을 새로 만들어 남겼다. 필요 없으면 사용자가 지운다.
- 자동 wake·MCP 부모 도구 목록 갱신·dots 연결은 이번 범위가 아니다. 별도 후속으로 이슈에 열려 있다.
- 실제 왕복 성공 근거는 반복하지 않았다.

## 5. 완료 판정

trial 인증·앱·policy 정리와 기존 owner 보호 복귀, 공유 회귀 통과로 지시서 완료 기준을 충족했다. 실패하거나 건너뛴 단계는 없다. 변경한 파일은 이 기록과 인박스·logs뿐이다. issue 최종 댓글과 main 병합은 coor가 한다.

## coor 종료 확인

worker 보고 이후 [Grok 정리 댓글](https://github.com/tourbut/KnowsLink/issues/1#issuecomment-5993825445)에서 시험 파일 제거 보고를 확인했다. 위 4절의 미확인은 worker 보고 시점의 관측이다. coor는 고정 SHA `2182401680806f4938cf8ad46a9f04c644e60bf5`의 제품 변경 없음·완료 기록·자원 대조·회귀 결과를 검수했다. [최종 이슈 댓글](https://github.com/tourbut/KnowsLink/issues/1#issuecomment-5993872965)을 게시하고 원문을 대조했다.
