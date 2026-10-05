---
title: SAR-PUBLIC-SERVICE-OPS-READINESS 실행 기록
status: draft
updated: 2026-10-05
owner: ops
tasks: [SAR-PUBLIC-SERVICE-OPS-READINESS]
summary: 일반 이메일 서비스의 실제 신원 보호와 서버 자원 및 백업 상태를 읽기 전용으로 관측하고 공개 전 선행 조건과 rollback 순서를 제안한다
---

# SAR-PUBLIC-SERVICE-OPS-READINESS — 일반 이메일 서비스 운영 준비 근거

## 기준과 범위

착수 HEAD는 `6865988`(브랜치 `fullops/ops`)다. 기준 main은 `6c0d132d26dac22693f1c1ffb986fc4f61bdbd09`다. 서버 배포 체크아웃은 `0911c2c73468f8684260a277d4940a74d26bcf7d`다(`git rev-parse HEAD`, 작업 트리 깨끗). Task는 task_22ac345acfc3, Dispatch는 ctx_fe22da7866ab다. 제품 기준 담당은 designer ctx_000841343790다.

적용 규칙은 `fullops-common-0.3.2`, FULLOPS.md, project.md, orca-agents.md, 문서 작성·coding-style/testing/security 규칙이다.

이번 작업은 읽기 전용이다. 서버·Cloudflare 설정 쓰기, 새 token 발급, 인증 해제, agent secret 발급, 메일 발송, 제품 수정은 하지 않았다. 관리 API token(`state/cf-service-token-api.env`, 0600)은 Python 메모리에서 GET 요청에만 썼다. 이메일·secret·사용자 데이터는 출력하지 않았다. 관측 시각은 2026-10-05T12:10Z–12:30Z다. 임의 성능 보장과 제품 수치는 만들지 않았다.

## 1. 신원 보호: 현재 무엇이 있는가

| 항목 | 관측 (읽기 전용) |
|---|---|
| Access 앱 | `KnowsLink beta (owner-only)` 한 개. `link.knowslog.com` 전체 hostname을 덮는다(`self_hosted`, session 24h, 자동 IdP 이동). `options_preflight_bypass=false`. `oauth_configuration`은 비어 있다(Managed OAuth 꺼짐). |
| 앱 정책 | `knowslink-beta-owner-only`: `allow`, include 규칙 종류는 `email` 하나(값은 출력하지 않음), exclude·require 없음. 계정에는 정책이 두 개뿐이다(다른 하나는 `knowslog-bot - Production`, 이 앱에 연결되지 않음). |
| IdP | 앱의 `allowed_idps`는 1개다. 이름·type 목록은 확인하지 못했다. 현재 token은 IdP 읽기 권한이 없어 목록이 빈 배열로 돌아왔다. 이전 기록(11.2)은 One-time PIN 한 개다. |
| service token | 0개(`GET service_tokens` 빈 목록). trial 앱·policy 없음. |
| Tunnel | `knowslink`: 연결 4개(icn01 1개, icn06 2개, icn07 1개), connector 2026.9.1, 시작 2026-10-05T11:47Z. 상태 `config.yml`은 `link.knowslog.com` → `http://relay:8080` 한 규칙과 404뿐이다. `originRequest.access.required: true`, `teamName: scshin88`, `audTag` 한 개다. |
| 공개 probe | 인증 없는 `GET` `/healthz`, `/owner`, `/v1/signup`, `/v1/registry`, `/v1/test/registry`, `/mcp`, `/.well-known/oauth-authorization-server`가 모두 302 → `scshin88.cloudflareaccess.com`이다. 외부에서 relay 응답을 받는 경로는 없다. |
| 도구 | `verify.py local`·`regression`·`public` 모두 exit 0. `access_apply.py selftest` exit 0. 공유 서비스는 orca 200·s8 200·mcp 401, 호스트 cloudflared PID 506937 그대로다. |
| 시험 흔적 | relay 컨테이너의 `KNOWSLINK_TEST_AGENTS`는 빈 값이고 상태 `.env`에도 줄이 없다. trial token·앱·policy는 없다. |

### 1.1 relay가 아는 신원과 Access의 신원

- relay 코드(`internal/relay/http.go`)는 `Authorization: Bearer`(agent/owner API)와 HTTP Basic(owner UI)만 읽는다. `Cf-Access-Jwt-Assertion`, 이메일, `X-Forwarded-*`를 읽는 코드는 없다(`grep` 결과 0건).
- 그러므로 현재 이메일 신원은 Cloudflare 가장자리에서만 존재한다. relay 안의 owner는 합성 owner credential이다. 합성 owner token과 실제 이메일 신원은 연결되어 있지 않다.
- 지금 Access 정책은 사용자 본인 이메일 한 개만 통과시킨다. 일반 사용자는 로그인 단계에서 거부된다.
- `originRequest.access.required: true`는 Access가 붙인 JWT를 cloudflared가 검증하게 한다. 이것은 "Access를 거쳤다"는 증명이다. 누가 로그인했는지를 relay에 전달하지 않는다.

### 1.2 일반 이메일 지원 방식의 선택지 (Cloudflare 공식 문서 근거)

- 선택 A. 가장자리에서 Access로 이메일 신원을 받는다. 정책에서 모든 One-time PIN 이메일을 허용하는 방식이다. 공식 문서는 `Include → Everyone`과 `Include → Login Methods = One-time PIN`을 "누구나 접근 가능한 구성"으로 분류한다. 그러므로 이 구성은 root 앱을 모두 공개로 바꾼다. relay가 JWT의 이메일 claim을 읽어 owner에 묶는 구현이 필요하다. 현재 이 구현은 없다. 해당 구현은 DEV 범위다.
- 선택 B. relay가 직접 이메일 신원을 확인한다(자체 OTP 또는 링크). 이 방식은 relay가 메일을 발송해야 한다. 발송 경로(Cloudflare Email Service 또는 외부)·도메인 SPF/DKIM/DMARC·발송 한도가 모두 미확정이다. 이번 작업은 실메일을 발송하지 않았고 발송 자원을 조사하지 않았다(zone 읽기 권한 없음).
- 선택 C. Cloudflare를 IdP로 쓴다. 2026-05-19 changelog에서 Cloudflare 계정 로그인을 지원한다. 일반 사용자는 Cloudflare 계정이 있어야 한다. 일반 서비스 후보로는 맞지 않을 가능성이 높다. 이 판단은 designer 몫이다.
- 어느 선택이든 Access 앱은 root에서 경로별로 나눠야 한다. Cloudflare 공식 문서(Application paths)는 더 구체적인 path가 우선한다고 설명한다. 이전 시험(D12 13장)에서 `/v1/test/*` 별도 앱으로 root owner 앱을 바꾸지 않고 경로를 분리한 실적이 있다.
- 이메일 로그인 사용자 수 제한(Zero Trust Free 요금제의 좌석 한도)은 이번에 공식 문서 검색으로 확인하지 못했다. 미확인이다. 공개 전에 계정의 실제 요금제와 좌석 한도를 사용자가 대시보드에서 확인해야 한다.

### 1.3 agent·원격 MCP 연결

- 현재 service token은 0개다. 이전 시험의 관리 token 범위는 만료일 2026-10-05T23:59:59Z까지다. 토큰 확인 결과 `active`다.
- 비브라우저 client(agent, CLI, MCP)는 지금 Access의 302를 받는다. 로그인 페이지를 완료하지 못한다.
- 선택 1. service token(`CF-Access-Client-Id/Secret`)을 별도 경로 앱의 `non_identity` 정책으로 허용한다. 시험에서 검증한 방식이다. 사용자별 발급·철회를 Cloudflare API가 관리한다. 일반 사용자 수만큼 token을 발급하는 구조는 맞지 않는다. 시험은 24h token 두 개였다.
- 선택 2. Access Managed OAuth를 켠다(공식 changelog 2026-03-20). 비브라우저 client가 401과 `WWW-Authenticate`로 OAuth 발견 endpoint(RFC 8414·9728)를 받는다. 사용자 브라우저 로그인으로 token을 얻는다. 같은 정책이 적용된다. client는 RFC 8707을 지원해야 한다. 문서는 self-hosted 앱에서 opt-in이고 MCP server는 Access JWT를 검증해야 한다고 설명한다. 이 방식이 실제 Grok·OpenAI dot client와 호환되는지는 이번에 시험하지 않았다. 미확인이다.
- 선택 3. relay가 agent 전용 Bearer와 PoP를 이미 가진다. 가장자리의 Access를 agent 경로에서만 service token으로 두고 relay 인증에 맡기는 방식이다. 이전 시험의 구조다. 이 경우 `any_valid_service_token`, bypass, everyone은 쓰지 않는다(D12 13.2의 금지 유지).
- 어느 선택이든 ops는 이번에 설정을 바꾸지 않았다. 선택은 designer·DEV 결정 뒤 ops가 적용한다.

## 2. 자원 관측 (측정 시점의 값이며 보장값이 아니다)

| 항목 | 관측 |
|---|---|
| CPU | 4코어. load average 2.93 / 1.51 / 1.09 (1·5·15분). |
| RAM | 총 13,906 MiB, 사용 7,903, 가용 6,002, 여유(free) 586. |
| swap | 4,095 MiB 중 4,074 MiB 사용(약 99%). 측정 시점의 memory PSI `some`/`full` avg10·60·300은 0.00이다. `pswpout` 누적 1,902,757 페이지다. |
| 메모리 상위 | 개발 도구 프로세스(`orca-ide`, `claude`)가 각 약 290–450 MiB를 쓴다. 서비스 컨테이너가 아니다. |
| 디스크 | `/` 457 GB 중 123 GB 사용, 가용 312 GB(29%). |
| Docker | 이미지 38.88 GB(회수 가능 11.3 GB), 빌드 캐시 27.37 GB(25.62 GB 회수 가능, 공유 자원이므로 prune하지 않음), 볼륨 584.6 MB. |
| KnowsLink 컨테이너 | relay 7.8 MiB/256 MiB, cloudflared 15.3 MiB/128 MiB, postgres 37.3 MiB/512 MiB. CPU는 각 3% 이하, pids 8·7·8. |
| 설정된 한도 | `compose.ops.yaml`: postgres 512m·1 CPU·200 pids, relay 256m·0.5 CPU·128 pids, cloudflared 128m·0.5 CPU·100 pids, 로그 10m×3, restart `unless-stopped`. 값은 기술 설정이며 상품 한도가 아니다. |
| 다른 서비스 | `myportfolio-{worker,beat,viewer,postgres,redis}-1` Up 4 weeks. worker 76.5 MiB, viewer 20 MiB, postgres 22 MiB. 이 컨테이너는 메모리·CPU 한도가 없다(한도 표시 13.58 GiB = 호스트 전체). |
| 포트 | KnowsLink는 `127.0.0.1:8080`(relay)만 게시한다. Postgres는 게시하지 않는다. 호스트 `127.0.0.1:5432`는 myportfolio postgres다. |

해석 (판단이 아닌 사실):

- 서비스 컨테이너의 현재 사용량은 설정 한도의 약 3–7%다. 그러나 이 값은 유휴 상태다. 부하 시험을 하지 않았다.
- 호스트 swap이 거의 가득 차 있다. 여유 메모리가 낮다. 이 압력은 KnowsLink가 아니라 개발 도구 프로세스가 만든다. 공개 서비스가 같은 호스트에서 돌 때 다른 서비스(`myportfolio`, 호스트 cloudflared)에 영향을 줄 수 있다. 측정 시점에는 PSI가 0이었다.
- `myportfolio` 컨테이너에는 한도가 없다. KnowsLink의 한도가 다른 서비스를 보호하는 방향이다. 반대 방향 보호(다른 서비스가 KnowsLink를 압박)는 현재 구성에 없다.

## 3. DB 크기와 구조적 제약

| 항목 | 관측 |
|---|---|
| DB 크기 | `pg_database_size` 8,484,531 byte(8.3 MB). 데이터 디렉터리 64 MB. |
| 테이블 | `relay_state` 696,320 byte(총합), `goose_db_version` 24,576 byte. `relay_state`는 행 1개다. |
| 행 | 컬럼: `singleton` boolean, `epoch` bigint, `clock` timestamptz, `data` jsonb. 행 크기(`pg_column_size`) 7,202 byte. |
| Postgres 설정 | `max_connections` 100, `shared_buffers` 128MB, `wal_level` replica, `archive_mode` off, `fsync` on, `synchronous_commit` on, `max_wal_size` 1GB. 측정 시점 연결 8개. |

- 모든 업무 상태가 `relay_state` 한 행의 JSON에 있다. 아키텍처 문서(`architecture.md`)도 singleton 행과 authorization epoch를 명시한다. 요청마다 한 행을 CAS로 갱신한다.
- 현재 행 크기는 7.2 KB다. 가입자·agent·메시지가 늘면 행이 커지고 매 쓰기가 전체 JSON을 다시 쓴다. 처리량 상한과 행 크기 상한은 측정하지 않았다. 제품이 정한 owner 100명·agent 200개 제안(DEC-03, 미확정)이 이 구조에서 견디는지는 미측정이다. 부하 시험은 DEV 소유다. ops는 같은 구성에서 자원 측정을 지원할 수 있다.
- 데이터에는 아직 실제 사용자 정보가 없다(합성과 시험 identity).

## 4. 보존·백업·복구 준비

| 항목 | 관측 |
|---|---|
| 백업 파일 | `state/backups/`(0700)에 dump 6개, 모두 0600. 가장 최근 `28bd1bb-20261004T044744Z.dump`(13,920 byte). |
| 일정 | 자동 백업 없음. crontab에 knowslink 항목 없음, systemd timer 없음(`dpkg-db-backup`만 있음). 백업은 `beta.sh backup`·`beta.sh deploy` 때만 만들어진다. |
| 최신 상태와 격차 | 최신 dump는 `0911c2c` 배포 직전(2026-10-04T04:47Z)이다. 그 뒤 만든 trial identity·시험 메시지·정리 결과는 dump에 없다. 현재 DB의 신규 데이터는 합성·시험 데이터다. |
| 보관 위치 | live 볼륨과 같은 디스크(`/`)다. 다른 호스트·오프사이트 사본이 없다. 암호화·보관 기간·복구 목표(RPO/RTO)가 정해지지 않았다. `/home/shin/backups`는 없다. |
| 복원 검증 | `beta.sh restore-verify`가 네트워크 없는 임시 컨테이너에 복원하는 절차를 갖는다. 이전 Dispatch에서 실행했다. 이번에는 컨테이너를 만들지 않아 재실행하지 않았다. 최신 dump의 복원은 이번에 미검증이다. |
| 삭제 | 업무 행 삭제는 WAL·백업의 완전 삭제가 아니다. 이 보장은 held다(D12 11.4). |
| 복구 한계 | `deploy`는 `db/migrations` 차이가 있으면 중단한다. migration이 바뀌는 새 기능은 rollback이 코드 복귀만으로 끝나지 않는다. |

실제 사용자 이메일과 agent 연결을 받기 시작하면 위 백업 체계로는 불충분하다. 일정, 별도 보관 위치(다른 디스크 또는 호스트), 암호화, 보관 기간, 복원 시험 주기를 공개 전에 정해야 한다. 보관 기간·복구 목표는 제품 수치이므로 designer가 정한다. ops는 정해진 값으로 구현한다.

## 5. 공개 전 선행 조건과 순서 제안

구현 설계와 제품 결정은 DEV와 designer 몫이다. 아래는 ops 관점의 선행 조건이다.

### 5.1 선행 조건 (모두 충족되기 전에 공개 경로를 열지 않는다)

1. **신원 방식 결정(designer)**: 1.2의 A/B/C 중 하나. 결정 전에는 Access 정책을 넓히지 않는다.
2. **relay의 신원 연결(DEV)**: 선택 A라면 relay가 Access JWT를 서명·`aud`·`exp`로 검증하고 이메일 claim을 owner로 매핑해야 한다. 선택 B라면 메일 발송과 OTP 검증이 있어야 한다. 현재 두 가지 모두 없다.
3. **agent·MCP 연결 방식 결정(designer·DEV)**: 1.3의 선택 1–3 중 하나. Managed OAuth를 쓰려면 실제 대상 client(Grok, OpenAI dot)가 RFC 8707·9728 흐름을 지원하는지 DEV가 먼저 시험한다.
4. **필요 권한 확인(사용자)**: 아래 5.3 표.
5. **서버 보호 한도**: 익명 접근이 가능한 hostname에 대한 가장자리 rate limit·bot 보호가 필요한지 결정한다. 이번에 zone 읽기 권한이 없어 WAF·rate limit 규칙의 현재 상태를 확인하지 못했다. 미확인이다.
6. **백업 체계(ops, 값은 designer)**: 4장의 격차를 닫는다. 일정·보관 위치·보관 기간·복원 시험.
7. **DEV 검사**: 공개 한도 구현(DEC-03)의 경계·동시성·재시작·안전 종료 경로 충돌 검증. 이 검증 전에는 "공개 가능"으로 표시하지 않는다.
8. **독립 QA·리뷰**: tester의 고정 SHA 독립 QA와 직접 시각 검수. 미해결 critical/high는 차단이다.
9. **원점 검증**: 새 경로 규칙마다 `cloudflared tunnel ingress rule`로 어느 규칙·AUD가 매칭되는지 확인한다. 경로 앱마다 AUD가 다르면 규칙을 앞에서부터 나눈다. 경로 규칙에 `access.required`와 `audTag`가 없으면 그 경로는 원점 JWT 검증을 받지 않는다.

### 5.2 적용 순서와 rollback

1. 준비(쓰기 없음): 위 선행 조건 1–3, 5, 6 결정과 DEV 구현·검사 통과. 고정 SHA 확정.
2. 백업: `beta.sh backup`으로 새 dump를 만들고 `restore-verify`를 통과시킨다. dump 경로·SHA를 기록한다.
3. 배포: `beta.sh deploy <SHA>`. migration이 다르면 중단하고 D12 5장 절차를 따른다. 새 앱 설정은 아직 열지 않는다.
4. 보호 먼저: 새 경로 Access 앱·정책을 만든다. 기존 root owner 앱·정책은 바꾸지 않는다. `access_apply.py check`와 같은 live 대조(앱 type·destinations·정책 ID·decision·AUD)를 새 앱에도 적용한다.
5. 원점 설정: 후보 `config.yml`을 0600으로 만들고 pinned `cloudflared`로 `ingress validate`와 `ingress rule`을 확인한다. 현재 `config.yml`을 백업한 뒤 교체하고 knowslink connector만 재기동한다.
6. 검증: 무인증·잘못된 자격 negative, 정상 자격 positive, 기존 root 경로가 그대로 owner 보호인지 확인한다. `verify.py regression`으로 공유 서비스 불변을 확인한다.
7. 노출 확대(예: root를 일반 이메일로 확대)는 마지막이다. 이 단계 전에 negative·rollback 시험을 통과해야 한다.

rollback 순서 (실패 시, 역순):

1. 새 경로 Access 앱·정책을 삭제하거나 비활성화한다(적용한 ID만).
2. 직전 `config.yml` 백업으로 복원하고 knowslink connector만 재기동한다(`cmp`로 일치 확인).
3. 연결을 즉시 끊어야 하면 `beta.sh unexpose`로 connector만 중지한다. 앱·DB·볼륨은 보존한다.
4. 코드 문제면 `beta.sh deploy <이전 SHA>`(migration 동일일 때만). 비호환이면 D12 5장 3항: 현재 DB 추가 백업 → 격리 복원 검증 → 데이터 손실 범위를 확인한 뒤 적용.
5. 다른 project(`myportfolio`)·호스트 cloudflared·`orca` Tunnel은 건드리지 않는다. `down -v`와 공유 자원 prune은 금지다.

기존 root를 `everyone`·bypass·모든 OTP 허용으로 바꾸는 조사는 여기서 하지 않았고 적용도 하지 않았다. 그 변경은 relay가 JWT 신원을 검증하는 코드가 배포되기 전에는 실행하면 안 된다. 현재 relay는 Access 신원을 보지 않으므로 root를 열면 익명 사용자가 Bearer/Basic 인증 앞단까지 도달한다.

### 5.3 필요 권한·외부 입력

| 항목 | 현재 | 필요 |
|---|---|---|
| 관리 API token | 읽기: Access 앱·정책 GET 성공, service tokens GET 성공 | Access Apps and Policies Edit (새 경로 앱·정책). Managed OAuth는 앱 `PUT`(Apps and Policies Edit). 만료 2026-10-05T23:59:59Z이므로 새 token이 필요하다. |
| IdP | 목록 읽기 불가(빈 응답) | Access: Organizations, Identity Providers, and Groups Read (확인용). 새 IdP 추가에는 Write. |
| Zero Trust 조직 | `GET access/organizations` 403 | 같은 Organizations 권한(세션·로그인 설정 확인). |
| Tunnel 목록 API | `GET cfd_tunnel` 결과 빈 배열(권한 부족) | 필요 없음. `cloudflared tunnel list/info`(cert.pem)가 대체 관측이다. |
| zone | DNS·SSL·ruleset·firewall·rate limit·email routing 모두 403 | Zone: DNS Read, Zone Settings Read, Firewall Services/WAF Read, Email Routing Read(읽기 확인용). 쓰기는 결정 후 별도 승인. |
| 계정 요금제 | subscriptions 403 | 대시보드 확인 또는 Billing Read. 좌석 한도 확인에 필요. |
| 메일 발송 | 미조사 | 선택 B일 때 Cloudflare Email Service 또는 외부 발송 계정, 도메인 SPF/DKIM/DMARC, 발송 한도. 사용자 입력이 필요하다. |
| 사용자 이메일 | owner 이메일은 `/tmp/knowslink-beta-owner-email`(0600) | 일반 서비스에서는 필요 없다. |

## 6. 미측정·미확인 목록

- 부하 시험 결과(요청 처리량, 지연, 동시 접속). 측정하지 않았다.
- `relay_state` 행 크기 증가에 따른 쓰기 시간.
- 최신 dump의 격리 복원(이번에 실행하지 않음).
- IdP 이름·type 목록, Zero Trust 조직 설정, 좌석 한도.
- zone의 DNS·WAF·rate limit·SSL 모드·email routing 상태.
- 선택 A/B/C 각각의 실제 가입 흐름(실메일 발송 없이 불가). Managed OAuth와 Grok·OpenAI dot client의 호환.
- 호스트 swap이 가득 찬 원인의 장기 추이. 한 시점 관측이다.
- 호스트 cloudflared 2026.8.3은 2026.9.3 업그레이드 권고 경고가 있다. 공유 프로세스이므로 이번에 바꾸지 않았다.

## 7. 실행 증거 요약

| 검사 | 결과 |
|---|---|
| `python3 deploy/knowslink/verify.py local` | exit 0 (health 200, unknown 404, owner/api 401, loopback-only, private Postgres, 한도, 0600/0700) |
| `verify.py regression` | exit 0 (orca 200·s8 200·mcp 401, PID 506937) |
| `verify.py public` | exit 0 (무인증 5개 probe 모두 302) |
| `access_apply.py selftest` | exit 0 |
| `cloudflared tunnel list/info` | knowslink 연결 4개, connector 2026.9.1 |
| Cloudflare GET | token verify `active`(만료 2026-10-05T23:59:59Z), 앱·정책·service token 목록 성공, 나머지는 5.3의 403 |
| 공식 문서 확인 | Managed OAuth(2026-03-20), Tunnel `access` originRequest, Access 정책 순서·오설정 주의, Cloudflare IdP(2026-05-19), Application token(JWT claim) |
| 변경 | 서버·Cloudflare·제품 코드 변경 없음. 이 기록과 D12 14장, 인박스·logs만 바뀐다. |

## 완료 판정

지시서의 세 항목(신원 방식 확인, 자원·복구 관측, 배포 계약 제안)을 읽기 전용으로 수행했다. 막힌 단계는 없다. 권한 부족으로 확인하지 못한 항목은 5.3과 6장에 구분했다. D13은 일반 서비스 검증 전이므로 수락으로 갱신하지 않았다.
