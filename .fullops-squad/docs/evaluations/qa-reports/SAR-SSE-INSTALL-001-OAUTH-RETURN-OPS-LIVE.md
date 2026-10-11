---
title: OAuth 브라우저 복귀 수정의 운영 적용과 보존 smoke
status: draft
updated: 2026-10-11
owner: ops
tasks: [SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS]
summary: "독립 수락된 eac603b를 기존 Linux 서버에 relay만 갱신해 적용한 실제 SHA, 백업, 보존 증거, 무계정 HTML401 smoke 결과와 부모가 남은 실계정 검증"
---

# 운영 적용 결과 (2026-10-11 UTC 08:13~08:17)

원문 로그·백업·bundle·credential은 서버 비공개 경로에만 있다. 이 문서에는 redacted 요약만 둔다. 서버 시계는 UTC다. 공개 smoke는 이 PC에서 `https://link.knowslog.com`을 대상으로 했다.

## 적용 요약

| 항목 | 결과 |
|---|---|
| 이전 배포 | `121cde923798c745af1f952347fcec1cae40b853`. 시작 시 체크아웃 clean, HEAD 일치, relay·Postgres healthy, cloudflared up, local health 200 |
| 적용 제품 SHA | `eac603b6a41d2ce002e698600ac66b8fc67ad762`. 배포 체크아웃 `git rev-parse HEAD`로 확인, detached, 변경 0건 |
| 수락 근거 | tester 독립 리뷰·QA(`96c7b823`), designer 시각 수락(`610b2b0`), coor fixed-SHA 리뷰. 이 SHA가 `origin/main`의 조상임을 `git merge-base --is-ancestor`로 확인했다 |
| 전달 방식 | 서버 origin 경로가 없다. `origin fetch`와 `beta.sh deploy`를 실행하지 않았다. `git bundle create <파일> origin/main --not 121cde9`로 만든 bundle(SHA256 `4ffa17b0614b6c05b7c08ce0fa838e24955faec2b965e0bded1dc9de6481a926`, 673416 bytes)을 SFTP로 0600 전달했다. 서버에서 SHA256 일치, `git bundle verify`, `git fetch <bundle>`, `121cde9→eac603b` 조상 관계를 확인했다. bundle에는 이 SHA 이후의 문서 commit이 포함되지만 체크아웃은 `eac603b`로만 옮겼다 |
| 제품 diff | `121cde9..eac603b`의 제품 파일은 `internal/relay/remote_auth.go`, 테스트 2개, `testdata/consent_browser.mjs` 네 개뿐이다. `db/migrations`·`compose.yaml`·`deploy/`·`Dockerfile`·`go.mod`·`go.sum`·`Makefile` diff는 0건이다 |
| DB migration | 없음. goose 버전 1 유지. migrate 컨테이너 exit 0 |
| 서버 변경 | relay 이미지 재빌드·재생성(`beta.sh up`, 기존 Compose project `knowslink`). `deploy-history.log`에 `121cde9 -> eac603b oauth-return bundle deploy` 한 줄 추가 |
| 변경하지 않은 것 | Tunnel `config.yml`, cloudflared 컨테이너, Postgres 컨테이너·볼륨, `.env`, Google 설정, compose 파일. 수정이 새 ingress 경로를 요구하지 않는다. `/remote/consent/<43자>`는 기존 allowlist에 있다 |

`beta.sh deploy`는 `git fetch origin`이 먼저 실패하므로 쓰지 않았다. 같은 순서(migration diff 검사→backup→history 기록→detached checkout→`beta.sh up`→`verify.py local`)를 수동으로 실행했다. 배포 체크아웃이 시작 때 clean이고 HEAD가 `121cde9`일 때만 진행하도록 명령에 조건을 넣었다. 강제 reset·clean은 쓰지 않았다.

## Private backup

- 파일: `knowslink-state/backups/121cde9-pre-oauthret-20261011T081432Z.dump`. 모드 0600, 10594 bytes, SHA256 앞 16자 `16d6cc9d6d4081bf`.
- `pg_restore --list`(고정 postgres 이미지) 종료코드 0, TOC 23줄. 복구와 DB 덮어쓰기는 하지 않았다.
- backup은 서버 밖으로 복사하지 않았다. 이전 backup 12개를 변경하지 않았다.
- 전달 bundle은 `knowslink-state/oauth-return-20261011/`(0700)에 보존했다.

## 보존 증거

| 항목 | 전 → 후 |
|---|---|
| `.env`, Tunnel credential JSON, `tunnel.uuid`, `access.aud` SHA256 | 모두 동일 |
| Tunnel `config.yml` SHA256 | `f5abe860a019d5379cbbd5499eb309a645c4ae1e8f0b60fab6fea22326458512` 동일. 지시서 기준값과도 같다 |
| relay_state 집계 | Agents 26, Challenges 0, Connections 2, Gates 0, GoogleAttempts 0, HTTP 0, Idempotency 10, Identities 1, Members 1, Messages 10, Owners 23, Pairs 12, Rates 3, Remote 6, Sessions 0이 전·후 동일. relay_state 행 1, goose 버전 1 |
| Postgres 컨테이너·볼륨 | 같은 컨테이너(생성 2026-10-03), 볼륨 `knowslink_postgres-data` 1개 그대로 |
| 권한 | state 0700, `.env` 0600, tunnel 0700, config 0600, backups 0700 |
| cloudflared | KnowsLink 컨테이너 Up(5시간 유지, 재시작 없음). host cloudflared PID 506937 동일 |
| 다른 서비스 | myportfolio 컨테이너 5개 상태와 Compose project 목록 동일. orca 200 / s8 200 / mcp 401 |
| 다른 역할 자원 | `knowslink-tester-oauth001-d7e2c059` 컨테이너가 실행 중이다. 이 과제가 만들지 않았고 변경하지 않았다 |

시작 시점 Agents 26과 Remote 6은 이전 기록(Agents 25)보다 늘었다. 이 과제 시작 전에 이미 있던 변화다. 원인은 사용자가 승인한 연결 시도로 추정하며 확인하지 않았다. 이 과제는 DB에 쓰지 않았다. 집계 비교는 이 과제 시작 시점이 기준이다. 미완료 설치(install1)는 삭제하지 않았다.

## Smoke 결과

| 검사 | 결과 |
|---|---|
| `verify.py local` | 통과. health 200, unknown 404, owner/api 401, loopback 8080, Postgres 비공개, 제한, 0600/0700 |
| `docker compose config --quiet` | 통과. 기존 UID/GID/STATE_DIR 계약을 유지 |
| 컨테이너 | relay healthy(재생성 직후), migrate exit 0, postgres healthy 유지. relay 로그 10분 error/panic 0건, `google login configured=true` |
| 합성 consent 경로(핵심) | `GET /remote/consent/<합성 43자>`(영숫자 43자, 존재하지 않는 attempt) → **401**, `Content-Type: text/html; charset=utf-8`. 본문은 새 안내 화면(`원격 연결을 계속할 수 없습니다`, 문제·안내 문구, `자기 홈으로 돌아가기` 링크)이다. 헤더에 `Cache-Control: no-store`, `Referrer-Policy: no-referrer`, `X-Content-Type-Options: nosniff`, 기존 CSP(`default-src 'none'; style-src 'unsafe-inline'; form-action 'self' https://accounts.google.com; frame-ancestors 'none'; base-uri 'none'`)가 있다. 사용자 원 링크와 실제 attempt·code·token은 쓰지 않았다. 새 설치를 만들지 않았다 |
| 공개 `/`, `/home` | 200 / 303 (기존 동작) |
| 공개 `/healthz`, `/owner`, `/v1/test/pull` | 404 |
| OAuth discovery | `/.well-known/oauth-protected-resource`, `…/mcp`, `oauth-authorization-server` 모두 200 |
| 무인증 remote | `POST /mcp` initialize → 401, `WWW-Authenticate: Bearer resource_metadata=…`. `/remote/oauth/authorize` 빈 요청 400. `/remote/oauth/token` 빈 요청 400 |
| SSE 무인증 | `GET /v1/text/events` 401. `?token=x` query도 401 |
| polling 기존 경로 | `POST /v1/text/pull`·`ack`: JSON 본문 무인증 401 `invalid_auth`, 빈 본문 422 `invalid_json` (기존 계약) |
| 공유 서비스 | orca 200 / s8 200 / mcp 401, 변경 없음 |

`verify.py regression`은 baseline JSON의 `projects`에 `knowslink`가 남아 있어 실패한다. 이전 LIVE 기록과 같은 기존 비교 결함이다. 출력에서 codes·PID·containers 세 항목이 baseline과 같고 projects만 `knowslink` 하나 차이임을 직접 대조했다. 스크립트는 수정하지 않았다.

## 실행하지 않은 검사와 한계

- 실제 Google 계정 consent, Dots·ChatGPT의 실제 callback 복귀, 인증된 SSE stream·Events는 실행하지 않았다. 이 과제는 계정과 외부 메시지를 만들지 않는다. 부모가 배포 뒤에 수행한다.
- 실제 승인 뒤 반환 화면(`앱으로 돌아가는 중`) 성공 경로는 운영에서 확인하지 못했다. 이 화면은 유효한 attempt·로그인 세션이 필요하다. 근거는 고정 SHA의 tester Chromium QA와 designer 시각 수락이다. 이 smoke는 같은 handler의 실패 경로가 새 코드로 응답함을 보인다. 운영에서 새 코드가 실행 중이라는 직접 증거이며, 성공 경로의 증거는 아니다.
- 제품 전체 QA는 반복하지 않았다. 고정 SHA 증거를 재사용했다. 새 실행 파일·스크립트·compose 변경은 없다. 등록 검증 게이트와 추가 코드 리뷰 대상이 없다. 원 exact-SHA product gate는 tester 기록에 있다. Windows 로컬에서는 Linux Go·Chromium gate를 다시 실행할 수 없다.
- 운영 도메인의 무인증 `GET /v1/text/events` 401 응답에는 `no-store`·CSP·`Referrer-Policy` 헤더가 없다. 이전 LIVE 기록의 medium(SSE 라우트가 보안 헤더 middleware 밖)이며 이 변경으로 바뀌지 않았다. 담당은 dev다.
- reviewer가 남긴 F1 medium(원인별 실패 안내 부족)은 비차단 잔여로 그대로 둔다. 새 consent 화면은 한도 초과·CSRF·로그인 불일치도 같은 만료 안내로 표시한다. 담당은 dev다.
- 새 결함은 발견하지 않았다.

## 남은 부모 단계

1. 실제 Google 계정으로 Dots(또는 ChatGPT) 원격 연결을 다시 시도하고 승인 뒤 앱으로 돌아가는지 확인한다.
2. 연결 뒤 token 발급과 인증 SSE·Events를 관찰한다. 로그에 Authorization과 본문이 남지 않는지 확인한다.
3. 기존 미완료 install1의 처리는 사용자·부모가 판단한다. 이 과제는 삭제하지 않았다.

## 실패 시 안전 조치

1. 코드 문제 시 이전 코드로 복귀한다. 서버에 `121cde9` 객체가 있다. DB 상태를 보존한 채 `git checkout --detach 121cde923798c745af1f952347fcec1cae40b853`와 `beta.sh up`을 쓴다. 이전 binary는 같은 DB를 읽는다. migration이 없다.
2. DB를 새 backup으로 덮지 않는다. 기존 키·철회·관계 상태를 보존한다. 위 backup은 사고 대비용이다.
3. ingress는 바꾸지 않았다. Tunnel 복구는 필요 없다. 이전 복구 절차는 [이전 LIVE 기록](SAR-SSE-INSTALL-001-OPS-LIVE.md)을 따른다.
