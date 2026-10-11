---
title: SSE·Remote 수락 후보의 운영 적용과 보존 smoke
status: draft
updated: 2026-10-11
owner: ops
tasks: [SAR-SSE-INSTALL-001-OPS]
summary: "수락 제품 121cde9를 기존 Linux 서버에 최소 적용한 실제 SHA, 백업, 보존 증거, ingress 변경과 smoke 결과 및 남은 검증"
---

# 운영 적용 결과 (2026-10-11 KST 12:36~12:45)

원문 로그·백업·credential은 서버 비공개 경로에만 있다. 이 문서에는 redacted 요약만 둔다. 서버 시계는 UTC이며 아래 시각은 UTC다.

## 적용 요약

| 항목 | 결과 |
|---|---|
| 이전 배포 | `ba3754d2c41d3ba8258eae76fa0cf0cfad5b27c4`. 시작 시 체크아웃 clean, relay·Postgres healthy, cloudflared up, local health 200 |
| 적용 제품 SHA | `121cde923798c745af1f952347fcec1cae40b853`. 배포 체크아웃 `git rev-parse HEAD`로 확인, detached, 변경 없음 |
| 전달 방식 | 기존 origin 경로가 서버에 없다. `origin fetch`와 `beta.sh deploy`를 실행하지 않았다. 로컬 `^ba3754d HEAD` bundle(SHA256 `6364b8da61d3fa14909378b229eb8956684bc20d79a7dd6686fd2908834c76be`)을 SFTP로 0600 전달했다. `git bundle verify`, `git fetch <bundle> HEAD`, ba3754d→121cde9 fast-forward 조상 관계를 확인했다 |
| DB migration | `db/migrations` diff 없음. goose 버전 1 유지 |
| 서버 변경 | relay 이미지 재빌드·재생성(`beta.sh up`, 기존 compose project `knowslink`). Tunnel `config.yml`의 ingress path 한 줄. cloudflared 컨테이너 재시작. 다른 변경 없음 |
| `deploy-history.log` | `ba3754d -> 121cde9 sse-install bundle deploy` 한 줄 추가 |

`beta.sh deploy`는 `git fetch origin`이 먼저 실패하므로 쓰지 않았다. 같은 순서(migration diff 검사→backup→history 기록→detached checkout→`beta.sh up`→`verify.py local`)를 수동으로 실행했다. `beta.sh backup`은 실패 시 파일을 지우는 동작이 있어 같은 pg_dump·pg_restore --list 명령을 직접 실행했다.

## Private backup

- 파일: `knowslink-state/backups/ba3754d-pre-sse-20261011T033643Z.dump`. 모드 0600, 10315 bytes.
- `pg_restore --list`(고정 postgres 이미지) 종료코드 0, TOC 23줄. 복구와 DB 덮어쓰기는 하지 않았다.
- backup은 서버 밖으로 복사하지 않았다. 기존 backup 파일을 변경하지 않았다.
- 전달 bundle과 pre-snapshot은 `knowslink-state/sse-install-20261011/`(0700)에 보존했다.

## 보존 증거

| 항목 | 전 → 후 |
|---|---|
| `.env`, Tunnel credential JSON, `tunnel.uuid`, `access.aud` SHA256 | 모두 동일 |
| relay_state 집계 | Pairs 12, Agents 25, Owners 23, Members 1, Messages 10, Identities 1, Connections 2, Idempotency 10, GoogleAttempts 0 동일. Rates 3→4는 정상 사용 증가 |
| Postgres 컨테이너·볼륨 | 같은 컨테이너(시작 2026-10-03), 볼륨 `knowslink_postgres-data` 그대로 |
| 상태 디렉터리·`.env`·tunnel 권한 | 0700 / 0600 / 0700 / config 0600 |
| 다른 서비스 | orca 200, s8 200, mcp 401. host cloudflared PID 506937. myportfolio 컨테이너 6개 상태와 Compose project 목록 동일 |
| Google env 키 | `.env`에 기존 Google 키 3개 존재. 값은 읽지 않았다. 새 환경 변수는 추가하지 않았다 |

`verify.py regression`은 baseline JSON의 `projects`에 `knowslink`가 남아 있어 실패한다. 2026-10-10 LIVE 기록과 같은 기존 비교 결함이다. 출력에서 codes·PID·containers 세 항목이 baseline과 같고 projects만 `knowslink` 하나 차이임을 직접 대조했다. 스크립트는 수정하지 않았다.

## Ingress 변경

- 원본 config SHA256 `3b6f569e5603489496cda27802626043ef5716dae9954c0e457c95da4c7c8cd3`을 `tunnel/config.pre-sse-20261011.yml`(0600)에 보존하고 `cmp`로 동일을 확인했다.
- 후보는 `deploy/knowslink/tunnel/public-ingress.yml`(수락 SHA)로 만든 `sse-install-20261011/config.sse.candidate.yml`(0600)이다. 기존 config 대비 차이는 `path` 정규식 한 줄이다. 추가 경로는 `/v1/text/events`, `/mcp`, `/.well-known/oauth-protected-resource(/mcp)`, `/.well-known/oauth-authorization-server`, `/remote/oauth/(register|authorize|token|revoke)`, `/remote/consent/<43자>`다.
- `scripts/check_public_ingress.py` 통과. 서버 host cloudflared 2026.8.3의 `ingress validate` OK. `ingress rule`로 위 신규 경로·`/home`은 규칙 #0(relay), `/owner`·`/healthz`·`/v1/test/pull`은 규칙 #1(404)임을 확인했다.
- 적용 뒤 `config.yml`을 후보와 같은 내용으로 덮고(0600) KnowsLink cloudflared 서비스만 재시작했다. 새 config SHA256 `f5abe860a019d5379cbbd5499eb309a645c4ae1e8f0b60fab6fea22326458512`. 4개 연결이 icn06/icn07에 등록됐다. 종료 때의 `Connection terminated` ERR는 재시작 로그다.
- hostname, `required:false`, 404 fallback과 owner/admin/test 비공개 경계를 유지했다. Access 앱·DNS·credential은 변경하지 않았다.

## Smoke 결과

로컬은 서버 loopback, 공개는 이 PC에서 `https://link.knowslog.com` 대상이다. 응답 본문은 공개 metadata와 오류 코드만 확인했다.

| 검사 | 결과 |
|---|---|
| `verify.py local` (배포 전·후) | 통과. health 200, unknown 404, owner/api 401, loopback 8080, Postgres 비공개, 제한, 0600/0700 |
| 컨테이너 | relay healthy(2026-10-11T03:37Z 재생성), migrate exit 0, postgres healthy 유지. relay 로그 error/panic 0건, `google login configured=true`, 신규 환경 변수 없음 |
| 공개 `/`, `/home` | 200 / 303→`/?n=expired` (기존 동작) |
| 공개 `/healthz`, `/owner`, `/v1/test/pull` | 404 |
| OAuth discovery | `/.well-known/oauth-protected-resource`(+`/mcp`)와 `oauth-authorization-server` 200 JSON. resource `https://link.knowslog.com/mcp`, PKCE S256, bearer header만 지원 |
| 무인증 remote | JSON-RPC `initialize`를 `/mcp`에 POST → 401, `WWW-Authenticate: Bearer resource_metadata=…`. `/remote/oauth/authorize` 빈 요청 400. `/remote/oauth/token` 빈 grant 400 `invalid_grant` |
| SSE 무인증 | `GET /v1/text/events` 401 `invalid_auth`. `?token=x` query credential도 401. Authorization 없이 본문·stream 없음 |
| polling 기존 경로 | `POST /v1/text/pull`·`ack` 무인증 401. `send`·`connect/start` 빈 본문 422 `invalid_schema` (기존 계약) |
| 공유 서비스 | orca 200 / s8 200 / mcp 401, 변경 없음 |

## 실행하지 않은 검사와 한계

- 인증된 SSE stream(ready·arrival·본문 없음·heartbeat·55초 수명)은 운영 도메인에서 실행하지 않았다. 이 경로는 Google 회원 agent 자격이 필요하다. 운영 DB에 합성 owner/agent를 만들지 않았다. 격리 스택도 만들지 않았다. 제품 코드 `relay`·`adapters`는 tester가 `febd128`에서 실행한 Linux `make lint`·`make test`(종료코드 0)와 바이트 동일하다. `febd128`→`121cde9` 사이 변경은 `.claude`·`.codex` hook, `AGENTS.md`, `CLAUDE.md`, `Makefile`, `scripts/shell_guard.py`와 테스트뿐이다. tunnel 경유 SSE 버퍼링은 실제 계정으로 확인해야 한다.
- 운영 도메인의 무인증 `GET /v1/text/events` 401 응답에는 `Cache-Control: no-store`·CSP·`Referrer-Policy`·`nosniff` 헤더가 없다. 같은 호스트의 `POST /v1/text/pull`에는 있다. 이것은 tester가 `febd128`에서 기록한 medium(`http.go`, SSE 라우트가 보안 헤더 middleware 밖)이며 coor가 수락 가능한 잔여로 기록했다. 제품 코드가 바이트 동일하므로 새 결함이 아니다. 본문은 `{"error":"invalid_auth"}`뿐이다. 이 과제에서 제품 코드를 고치지 않았다. 담당은 dev다.
- 실제 Google consent, OAuth client 등록(`/remote/oauth/register`), Dots·ChatGPT 계정 연결, MCP Events 구독·webhook, 실제 agent 발신은 실행하지 않았다.
- 새 실행 파일·스크립트·compose 변경은 없다. 등록 검증 게이트와 추가 코드 리뷰 대상이 없다.

## 남은 사용자·부모 단계

1. Google 계정으로 자기 host 설치(Claude/Codex/Grok/Dots)를 실행하고 `KNOWSLINK_RECEIVE_MODE=sse`로 인증 SSE를 실제 관찰한다.
2. Remote OAuth 동의와 Cloud app ID 등록은 사용자·부모 후속이다.
3. 실제 SSE를 확인할 때 로그에 Authorization과 본문이 남지 않는지 확인한다.

## 실패 시 안전 조치

1. 공개 오류 시 `config.pre-sse-20261011.yml`을 `config.yml`에 같은 방식으로 덮고 KnowsLink cloudflared만 재시작한다. 신규 경로는 다시 404가 된다.
2. 코드 문제 시 새 SSE 연결과 local receive mode를 poll로 돌린다. 이전 코드로 복귀할 때는 `ba3754d` bundle 객체가 서버에 남아 있다. remote 사용 이후에는 이전 binary가 shared JSON의 선택 field를 잃을 수 있으므로 사용 중지와 private backup 확인 뒤 `git checkout --detach ba3754d`와 `beta.sh up`을 쓴다.
3. 새 backup으로 DB를 덮지 않는다. 기존 키·철회·관계 상태를 보존한다.
