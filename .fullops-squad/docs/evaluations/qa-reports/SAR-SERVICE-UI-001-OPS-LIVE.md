---
title: 수락한 회원 웹·플러그인 UI의 상태 보존 운영 배포
status: draft
updated: 2026-10-11
owner: ops
tasks: [SAR-SERVICE-UI-001-OPS]
summary: "수락 제품 95df83c를 기존 Linux 서버에 relay만 갱신해 적용한 baseline·after·명령·종료코드·보존 증거·플러그인 package 해시·복귀 방법과 미검증 항목"
---

# 운영 적용 결과 (2026-10-11 UTC 14:26~14:30)

원문 로그·bundle·package 산출물·credential은 서버 비공개 경로에만 있다. 이 문서에는 redacted 요약만 둔다. 서버 시계는 UTC다. 공개 smoke는 이 PC에서 `https://link.knowslog.com`을 대상으로 했다. 서버 접속은 기존 `coor_remote.py` helper(RejectPolicy, 값 redaction)를 읽기 전용으로 재사용했다.

## 적용 요약

| 항목 | 결과 |
|---|---|
| 이전 배포 | `eac603b6a41d2ce002e698600ac66b8fc67ad762`. 시작 시 detached, 변경 0건, relay·Postgres healthy, KnowsLink cloudflared up, local health 200 |
| 적용 제품 SHA | `95df83c22c1559453a7c73504aa248088a8e4ed0`. 배포 체크아웃 `git rev-parse HEAD`로 확인, detached, `status --porcelain` 0줄 |
| 준비 main | `0045089f2a2c3efe3e210eedc0c5cd4868487e42`. `95df83c`와 `.fullops-squad` 밖 diff 0건(제품 동일). 배포 SHA와 구분한다 |
| 수락 근거 | 지시서 「수락 고정 SHA와 검수」: tester 독립 리뷰·QA, designer 직접 시각 결과, coor fixed-SHA 리뷰. `95df83c`는 `origin/main` 조상이며 `0045089`의 조상임을 `git merge-base --is-ancestor`로 확인했다 |
| 전달 방식 | 서버 origin 경로가 없다. `git fetch origin`과 `beta.sh deploy`를 실행하지 않았다. `git bundle create <파일> HEAD ^eac603b`로 만든 bundle(SHA256 `132262f16674b23cde05dc7a19504c320ce7753065274b2edbdab4c36aae2b59`, 5774360 bytes)을 SFTP로 0600 전달했다. 서버에서 SHA256 일치, `git bundle verify`, `git fetch <bundle> HEAD`, `eac603b→95df83c` 조상 관계를 확인했다. bundle에는 이 SHA 이후의 문서 commit이 포함되지만 체크아웃은 `95df83c`로만 옮겼다 |
| 제품 diff | `eac603b..95df83c`의 제품 파일은 15개다. `internal/relay`의 member·http·테스트·golden, `adapters/src/grok-wake*`, `adapters/package.json`(test 스크립트에 `grok-wake` 추가), README 2개다. `db/migrations`·`compose.yaml`·`deploy/`·`Dockerfile`·`go.mod`·`go.sum`·`Makefile` diff는 0건이다 |
| DB migration | 없음. goose 버전 1 유지. migrate 컨테이너 exit 0 |
| 서버 변경 | relay 이미지 재빌드·재생성(`beta.sh up`, 기존 Compose project `knowslink`). migrate 컨테이너 재생성(기존 `up`의 정상 동작). `deploy-history.log`에 `eac603b -> 95df83c service-ui bundle deploy` 한 줄 추가. `knowslink-state/serviceui-20261011/`(0700)에 bundle과 package 산출물 보존 |
| 변경하지 않은 것 | Tunnel `config.yml`, cloudflared 컨테이너, Postgres 컨테이너·볼륨, `.env`, Google 설정, compose 파일, 다른 역할·프로젝트의 컨테이너 |

## 실행하지 않은 것 (지시서 금지 항목)

`beta.sh`의 `prepare`·`deploy`·`backup`·`restore-verify`·`seed`·`owner-login`, Tunnel create/render/expose, raw `pg_dump`·SQL dump, 파일·컨테이너 삭제·정리는 실행하지 않았다. 따라서 **이번 배포 직전의 새 DB backup은 없다**. 복귀 기준 DB 증거는 이전 backup 13개(최신 `ba3754d-pre-sse-20261011T033643Z.dump`)와 아래 집계 digest다. 이번 변경은 migration 0이며 이전 binary가 같은 DB를 읽는다. 원문 state는 수집하지 않았다.

## 명령과 종료코드

| 단계 | 명령(요약) | 종료코드 |
|---|---|---|
| baseline | 읽기 전용 `git`·`stat`·`sha256sum`·`docker ps`·`psql select` 집계 | 스크립트 종료 1(마지막 `grep -c`가 0건일 때 1. 실패 아님. 출력 전 항목은 정상) |
| 디렉터리·전달 | `mkdir -m 0700`, SFTP put, `chmod 0600`, `sha256sum` | 0 |
| 배포 | precheck(HEAD=eac603b·clean)→bundle SHA·verify→`git fetch <bundle> HEAD`→조상·infra diff→`git checkout --detach 95df83c`→history→`beta.sh up`→`verify.py local` | 전체 0. `beta.sh up` 0, `verify.py local` 0 |
| after | baseline 스크립트 재실행 | 1(위와 같은 `grep -c` 사유) |
| package | `git archive 95df83c`→`npm ci`→`npm run build`→`package_plugin.py --host all --verify` | 전체 0. `npm ci` 0, build 0, package 0 |
| 공개 smoke | `curl` 무인증 | 아래 표 |
| `verify.py regression` | 서버 | 1(기존 baseline 결함. 아래) |

## 보존 증거 (baseline 14:26:37Z → after 14:28:18Z)

| 항목 | 전 → 후 |
|---|---|
| `.env`, `tunnel.uuid`, `access.aud`, Tunnel credential JSON, `tunnel/config.yml` SHA256 | 모두 동일. config `f5abe860a019d5379cbbd5499eb309a645c4ae1e8f0b60fab6fea22326458512`는 이전 LIVE 기준과도 같다 |
| compose·`compose.ops.yaml`·`beta.sh` SHA256 | 동일(infra diff 0) |
| 권한 | state 0700, `.env` 0600, tunnel 0700, config 0600, backups 0700. 동일 |
| backups | 13개, 파일 변경 없음 |
| relay_state 키별 개수와 키 digest | Agents 28, Challenges 0, Connections 0, Gates 0, GoogleAttempts 0, HTTP 0, Idempotency 30, Identities 2, Members 2, Messages 30, Owners 24, Pairs 13, Remote 6, Sessions 1이 전·후 동일. goose 버전 1, relay_state 행 1 |
| 콘텐츠 md5 | Agents·Idempotency·Identities·Members·Messages·Owners·Pairs·Remote·Sessions·빈 맵의 md5가 전·후 동일(원문·값은 출력하지 않았다) |
| 변화 | `epoch` 723366→723482, `Rates` 6→7(키 digest·내용 md5 변경). relay가 살아 있는 동안의 rate 계수 갱신이다. 이 과제는 DB에 직접 쓰지 않았다 |
| Postgres 컨테이너·볼륨 | 같은 컨테이너(생성 2026-10-03, Up 8일 유지), 볼륨 `knowslink_postgres-data` 1개 그대로, 볼륨 총 48개 동일 |
| cloudflared | KnowsLink 컨테이너 Up 유지(재시작 없음), host cloudflared PID 506937 동일 |
| 공유·다른 서비스 | myportfolio 컨테이너 6개 상태, orca 200 / s8 200 / mcp 401, jevpet-demo 컨테이너 변화 없음 |
| 다른 역할 자원 | `knowslink-tester-oauth001-d7e2c059`(Up), `knowslink-dev-serviceui-*` 3개(Exited 0)를 만들지 않았고 변경·삭제하지 않았다 |
| 기존 포트 | 목록 동일(127.0.0.1:8080 등) |

**이전 기록과 다른 수치.** 이전 LIVE(`eac603b`, 2026-10-11 08:17Z)는 Agents 26, Members 1, Messages 10, Owners 23, Pairs 12, Connections 2였다. 이 과제 시작 시점의 baseline은 Agents 28, Members 2, Messages 30, Owners 24, Pairs 13, Connections 0이다. 변화는 이 과제 시작 전에 이미 있었고(그 사이 6시간 동안의 사용자·coor Bot 왕복 및 시험) 원인은 확인하지 않았다. 이전 수치에 맞추려고 DB를 바꾸지 않았다. 비교 기준은 이 과제의 baseline이다. Connections 2→0은 만료·정리일 수 있으나 확인하지 않았고, 이 과제의 영향은 아니다(baseline이 배포 전이다).

## Plugin package 검증

`95df83c`의 `git archive`를 `knowslink-state/serviceui-20261011/pkg/`(0700)에 풀어 Node 22.22.2·npm 10.9.7·Python 3.12.3으로 만들었다. 배포 체크아웃은 건드리지 않았다. `npm ci`·`npm run build`·`package_plugin.py --host all --verify` 종료코드 모두 0이다. 각 host는 추출한 standalone MCP 테스트와 준비·private 폴더 경계 검사를 통과했다.

| 산출물 | SHA256 |
|---|---|
| `knowslink-grok-bot-plugin.zip` | `762411df0e49dd0568c38b118d67c947910b39f56df5208afd897b5dc434cf94` |
| `knowslink-claude-plugin.zip` | `5388516edcfb992f777ef13f1c748df6b52a703fdc7200e391676da3e3ba8b9f` |
| `knowslink-codex-plugin.zip` | `d79abee4dc1639e8e83a12b10063a2fb63a04ba3f6492b3db62a53dab82c3410` |
| `knowslink-claude-desktop.mcpb` | `bfd13777aa21dadabbf46a7f0439a6d8a9a04c7a90b0bff2401cb7aabb5aacb2` |

Grok zip의 `knowslink/dist/plugin.js`(1641155 bytes)에서 `wake`·`wakeCheck`·`wakePrompt` 심볼을 확인했다. 즉 webhook wake 코드(`adapters/src/grok-wake.ts`)가 `95df83c`에서 빌드돼 들어 있다. 이 hash는 서버 Linux 빌드의 값이다. Windows 빌드와 같을 것이라고 가정하지 않는다. Bot runtime·env·MCP 재설치와 watcher는 coor 후속이며 건드리지 않았다. 산출물은 서버의 `/home/shin/deploy/knowslink-state/serviceui-20261011/pkg/build/`에 있고 서버 밖으로 복사하지 않았다. 새 설치용 배포는 coor가 정한다.

## 공개 smoke (이 PC → link.knowslog.com, 무인증)

| 검사 | 결과 |
|---|---|
| `GET /` | 200, `text/html`, no-store, CSP(`default-src 'none'; style-src 'unsafe-inline'; form-action 'self' https://accounts.google.com; frame-ancestors 'none'; base-uri 'none'`), `Referrer-Policy: no-referrer`, `nosniff`. 본문에 새 디자인 토큰(`--accent:#164da0`)과 `class="skip"`(건너뛰기 링크)이 있다. 새 코드가 실행 중이라는 직접 증거다 |
| `GET /home`, `/home/receipts` | 303 → `/?n=expired` (무세션, 기존 동작) |
| `GET /healthz`, `/owner`, `/v1/test/pull` | 404 (public deny 유지) |
| OAuth discovery | `/.well-known/oauth-protected-resource`, `oauth-authorization-server` 200 JSON |
| `POST /mcp` initialize | 401 |
| `/remote/oauth/authorize` 빈 GET | 400. `/remote/oauth/token` GET 405 |
| `/remote/consent/<합성 43자>` | 401 `text/html`, no-store, CSP |
| `GET /v1/text/events`, `?token=x` | 401 `application/json`. 보안 헤더 없음(이전 LIVE 기록의 기존 medium, dev 담당) |
| 공유 서비스 | orca 200 / s8 200 / mcp 401 |
| 서버 `verify.py local` | 통과(health 200, unknown 404, owner/api 401, loopback 8080, Postgres 비공개, 제한, 0600/0700) |
| 컨테이너 | relay healthy(재시작 0), migrate exit 0, postgres healthy 유지. 배포 뒤 relay 로그 5분 error/panic 0건, `google login configured=true` |

`verify.py regression`은 종료코드 1이다. 원인은 이전 LIVE와 같은 baseline 결함이다. baseline JSON의 `projects`가 `knowslink, myportfolio`인데 현재는 `knowslink`를 걸러 낸 뒤 `jevpet-demo, myportfolio`다. `jevpet-demo`는 이 과제 이전부터 있던 다른 프로젝트(배포 전 baseline의 compose ls에도 있다)이며 codes·PID·myportfolio 컨테이너 세 항목은 baseline과 같다. 스크립트는 수정하지 않았다.

## 미검증 (자동 PASS로 기록하지 않음)

- 새 Google 실제 계정 로그인·`/home` 회원 목록·모바일 렌더·키보드 포커스·새 플러그인의 실제 설치와 webhook 왕복은 실행하지 않았다. 담당은 coor(Orca 내장 브라우저)다. 필요한 정확한 운영 SHA는 `95df83c22c1559453a7c73504aa248088a8e4ed0`, URL은 `https://link.knowslog.com/`다.
- 인증 SSE·Events·Remote OAuth 실계정 흐름은 실행하지 않았다.
- 새 DB backup은 지시서 금지로 만들지 않았다. 미래 복귀가 DB 복원을 요구하면 coor 결정이 필요하다. 이번 코드는 DB 쓰기 형식을 바꾸지 않았다(제품 diff에 migration·schema 없음)고 diff로 확인했다. 다만 새 relay가 상태 JSON 필드를 추가로 기록하는지는 런타임에서 확인하지 않았다.
- 전체 제품 QA와 Linux Go·Chromium gate는 반복하지 않았다. 고정 SHA의 tester·visual 증거를 재사용한다. 서버에서의 `make test`는 실행하지 않았다.
- Bot watcher(PID 88404)·private webhook 설정은 확인·변경하지 않았다.

## 복귀 방법

1. 코드 문제 시 이전 코드로만 복귀한다. 서버에 `eac603b` 객체가 있다. `git -C /home/shin/deploy/knowslink checkout -q --detach eac603b6a41d2ce002e698600ac66b8fc67ad762` 뒤 `bash /home/shin/deploy/knowslink/deploy/knowslink/beta.sh up`과 `verify.py local`을 쓴다. 체크아웃이 clean이고 HEAD가 `95df83c`일 때만 진행한다. 이전 binary는 같은 DB를 읽는다(migration 없음).
2. DB를 backup으로 덮지 않는다. 기존 키·철회·관계 상태를 보존한다. `deploy-history.log`에 복귀 한 줄을 추가한다.
3. ingress·Tunnel은 바꾸지 않았으므로 복구가 필요 없다. 이전 Tunnel 복구는 [SSE LIVE 기록](SAR-SSE-INSTALL-001-OPS-LIVE.md)을 따른다.
4. 복귀하지 않은 이전 운영 제품은 `eac603b`([OAuth 복귀 LIVE 기록](SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS-LIVE.md))다.
