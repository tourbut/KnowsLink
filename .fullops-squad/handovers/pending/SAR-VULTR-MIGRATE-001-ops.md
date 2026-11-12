---
title: SAR-VULTR-MIGRATE-001 — 운영 relay를 홈서버에서 Vultr로 이전(대기 초안)
status: done
updated: 2026-11-12
owner: ops
tasks: [SAR-VULTR-MIGRATE-001]
summary: 사용자 지시에 따라 운영 relay를 Vultr 서울 새 서버로 옮기고, 홈서버 스택을 내리고, 도메인을 옮기는 인계 초안과 사전 조사 결과
---

# SAR-VULTR-MIGRATE-001 — 운영 relay Vultr 이전

**완료(2026-11-12):** Windows에서는 Vultr에 접속할 수 없다. 그래서 사용자 지시대로 Mac ops 세션이 직접 실행했다. 결과는 [D12 이전 기록](../../docs/operations/ops-guide.md#sar-vultr-migrate-001-운영-서버-vultr-이전-2026-11-12)에 있다. 아래는 실행 전 인계 초안이다. 이 파일을 dispatch 지시서로 쓰지 않는다.

## 사용자 지시 (2026-11-12)

- "현재 서버를 로컬에 있는 서버 말고 vultr에 만들어서 서비스하게 하자"
- "현재 작업중인거 windows에서 작업중인 coor에게 전달해 로컬 서버올린것도 내리고 작업하고 도메인도 옮겨"

범위는 세 가지다.

1. 운영 relay를 홈서버에서 Vultr 새 서버로 옮긴다.
2. 이전을 확인한 뒤 홈서버의 `knowslink` Compose 스택을 내린다.
3. `link.knowslog.com`을 Vultr 서버로 옮긴다.

## 사용자 결정 (Mac 세션 질문의 답)

- **서버:** Vultr 새 인스턴스 `vc2-1c-2gb`, 리전 `icn`(서울), Ubuntu 24.04(os 2284), 월 $10. 기존 `nonili-app-01`·`aitop100-bot`에는 올리지 않는다.
- **버전:** 현재 운영 `c6dfc39a`를 그대로 옮긴다. main의 migration 00003(되돌릴 수 없음)은 이번에 적용하지 않고 별도 승인을 받는다. 진행 중인 SAR-MSG-EPHEMERAL-001(Redis)도 이번 이전과 분리한다.

Mac 세션은 조회만 했다. Vultr 인스턴스를 만들지 않았고 홈서버도 바꾸지 않았다.

## 사전 조사 결과 (읽기 전용)

### Vultr

- Mac의 `vultr-cli` v3.11.0(`~/.vultr-cli.yaml`)이 인증되어 있다. Windows에 CLI나 API 키가 없으면 Mac 세션에 실행을 요청하거나 사용자에게 받는다.
- 기존 인스턴스는 2개(`nonili-app-01`, `aitop100-bot`, 둘 다 icn 1c/2GB)다.
- SSH 키 1개가 등록되어 있다(id `88ac30e7-67bc-4f79-8630-cabdb085ac9c`). Mac `~/.ssh/id_ed25519`와 같은 키다.
- firewall group은 0개다.

### 홈서버 (`ssh homeserver`, 192.168.0.74, user shin)

- docker 29.4.3, compose v5.1.3.
- Compose project `knowslink` = postgres(17-alpine) + relay + cloudflared. 실행 SHA는 `c6dfc39a`다.
- 사용량: relay 15MiB, postgres 30MiB, cloudflared 19MiB. DB 8.6MB, goose version 2.
- 구성 파일: `~/deploy/knowslink/compose.yaml`, `deploy/knowslink/compose.ops.yaml`, `~/deploy/knowslink-state/compose.install.yaml`. 마지막 파일은 state 폴더에 `.retired` 사본도 있으니 실제 내용을 확인한다.
- state: `~/deploy/knowslink-state`(0700). `.env` 키 이름은 `POSTGRES_PASSWORD`, `DATABASE_URL`, `RELAY_PORT`, `KNOWSLINK_GOOGLE_CLIENT_ID/SECRET/REDIRECT_URL`이다. 실행 중 relay의 환경변수 이름에는 `KNOWSLINK_CALENDAR_KEY` 등도 있다(값은 확인하지 않았다).
- relay는 `state/install`(0.1.0, 0.1.1의 `knowslink-setup.ps1`·plugin zip)을 `/install`에 ro로 마운트한다(`KNOWSLINK_INSTALL_DIR`).
- 공개 경로는 Cloudflare named tunnel이다(`state/tunnel/config.yml` + `<uuid>.json`, 0600). cloudflared 컨테이너는 user 1000:1000으로 실행한다.
- 백업 cron: `17 3 * * * ~/.local/bin/knowslink-prune-backups.sh ...`. Vultr로 함께 옮긴다.
- deploy checkout의 origin은 홈서버 로컬 경로다. Vultr에서는 GitHub origin에서 `c6dfc39a`를 clone한다.
- 홈서버에는 공유 서비스(jevpet-*, myportfolio, orca/s8/mcp.knowslog.com의 별도 cloudflared)가 있다. `knowslink` project 밖은 건드리지 않는다.

### 도메인

DNS는 tunnel CNAME이다. 같은 tunnel 자격을 Vultr로 옮기면 DNS를 바꾸지 않아도 도메인이 따라온다. 홈 cloudflared를 먼저 멈추고 Vultr cloudflared를 기동한다. 두 커넥터를 동시에 실행하면 트래픽이 양쪽으로 갈린다.

### 문서·도구의 홈서버 전제

`docs/operations.md`와 `deploy/knowslink/verify.py`는 홈서버 경로(`/home/shin/...`)와 공유 서비스 회귀 검사를 전제로 한다. Vultr용으로 경로와 검사를 조정한다.

## 권장 순서

D12 원칙(쓰기 중단 → 새 고유 백업 → 격리 복원 검증 → 기동 → 검증)을 유지한다.

1. Vultr 인스턴스를 만든다(label `knowslink-relay-01`, 위 SSH 키). firewall group은 인바운드 22만 허용하고 관리 IP로 제한한다. tunnel은 outbound라 공개 80/443이 필요 없다.
2. docker를 설치하고 비root 사용자를 만든다. `c6dfc39a`를 clone해 이미지를 빌드하고 postgres만 먼저 기동한다.
3. 홈서버에서 relay를 멈춰 쓰기를 중단한다. `pg_dump -Fc`로 새 고유 파일(0600)을 만들고 scp로 Vultr에 보낸다. Vultr에서 network-none 컨테이너로 복원을 검증한 뒤 실제 DB에 `pg_restore`한다.
4. state(`.env`, `tunnel/`, `install/`)를 0700/0600 권한 그대로 복사한다. 비밀값은 출력하지 않는다.
5. 홈 cloudflared를 멈춘다. Vultr에서 migrate·relay·cloudflared를 기동한다. local health와 공개 smoke를 확인한다. 공개 확인은 1초 이상 간격으로 한다(429 주의).
6. 확인 뒤 홈서버 `knowslink` project를 stop한다. 볼륨·이미지·state는 rollback용으로 보존하고 삭제하지 않는다. rollback은 Vultr cloudflared 정지 → 홈 스택 재기동이다.
7. prune cron을 옮기고 `ops-guide.md`·`docs/operations.md`·PLANS를 갱신한다.

## 승인

사용자가 이전·로컬 서버 내림·도메인 이전을 직접 지시했다. 실행 직전 다운타임(수 분 쓰기 중단) 확인은 coor가 판단해 받는다. 진행 중인 SAR-MSG-EPHEMERAL-001 dev 작업은 중단하지 않는다.
