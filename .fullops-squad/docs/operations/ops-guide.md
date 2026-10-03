---
id: D12
title: 운영자설명서
status: draft
updated: 2026-10-04
owner: ops
tasks: [SAR-DEPLOY-001-OPS, SAR-BETA-001-OPS, SAR-MVP-002-BOT-CATALOG-DEV]
upstream: [D02, D03]
summary: 서버 관찰 이력과 본인 전용 합성 베타 배포 구성·검증·복귀 절차 및 held 항목을 기록한다
---

# KnowsLink 운영자설명서 (D12) — 현재 서버 배포 계획

이 문서는 현재 서버의 Docker와 `link.knowslog.com` Cloudflare Tunnel 배포 정본이다. 현재 상태는 11장과 [D13](transition.md)을 따른다. 1~10장은 2026-10-03 이전 OPS Dispatch의 **역사적 계획·관찰 기록**이다. 그 Dispatch가 2026-10-03에 수행한 읽기 전용 조사 초안을 보존한다. 당시 retry Dispatch는 서버 조사를 재실행하지 않았다. 원본 명령 로그와 개별 종료코드가 초안에 없으므로 2장 관찰을 새 검증 통과로 해석하지 않는다. 당시에는 DNS·Tunnel·컨테이너를 만들거나 바꾸지 않았다. 이후 SAR-BETA-001-OPS가 로컬 스택과 `knowslink` Tunnel을 만들었다(11장). D11·D13은 작성됐다. 이 문서도 공개 배포 완료 근거가 아니다.
비밀값·토큰·`cert.pem` 내용은 기록하지 않는다. 존재·권한·경로만 기록한다.

**역사적 기록**: 이전 Dispatch는 `a6a10c7`을 C1 high 때문에 배포 금지로 기록했다. 이 결정은 그 당시 후보에만 해당한다. C1/RF-01 high는 해소됐다. 수락 후보는 main `557ebc3`, 제품 `78b1d92`, QA `659f4b0`, 최종 리뷰 `311381f`, 직접 UI 검수 `e238777`이다. 본인 전용 합성 베타의 구성과 검증은 11장에 기록한다. 전체 공개 한도·실제 신원·실데이터 held는 유지한다.

## 1. 승인 범위와 선행 조건

- 사용자 승인: 현재 서버에서 Docker 컨테이너로 배포하고 Cloudflare Tunnel로 `link.knowslog.com`에 연결한다. `knowslog.com`은 사용자 소유이며 Cloudflare가 관리한다.
- 선행 조건: DEV 고정 후보의 독립 QA, 직접 UI 검수, 코드 리뷰, coor 수락. 사용자가 현재 기록을 완료·병합·공유한 뒤 중지하라고 지시했다. 수락 SHA 전달만으로 재개하지 않는다. 사용자 재개 지시와 후속 Dispatch를 받은 뒤 6장을 실행한다. 공개 정책과 실제 신원 인증이 held이면 공개 연결은 계속 보류한다.
- 보존 대상: `myportfolio`·`myportfolio-master` Compose 프로젝트와 볼륨, `orca` Tunnel과 `orca/s8/mcp.knowslog.com` route, 기존 `~/.cloudflared/config.yml`과 실행 중인 `cloudflared` 프로세스, 다른 프로젝트 서비스.

## 2. 보존한 서버 사전조사 관찰 (2026-10-03, 이전 Dispatch)

| 항목 | 결과 |
|---|---|
| Docker | Engine 29.4.3, Compose v5.1.3. 사용자 `shin`이 `docker` 그룹이다. 조회 성공. |
| 실행 중 컨테이너 | `myportfolio` 프로젝트 5개(worker·beat·viewer·postgres·redis) 정상. `myportfolio-master-viewer-1`은 Created. `hermes-9ae8f587`은 Exited(255). |
| 호스트 게시 포트 | `0.0.0.0:8765`(myportfolio viewer), `127.0.0.1:5432`(myportfolio postgres). 호스트 프로세스: `0.0.0.0:6768`, `127.0.0.1:8766`, `127.0.0.1:9119`, SSH 22. |
| 포트 충돌 | 후보 게시 포트 `127.0.0.1:8080`은 비어 있다. Postgres 포트는 게시하지 않는다. 원래 조사 시점에 충돌을 관찰하지 않았다. 재개 전에 다시 확인한다. |
| Docker 네트워크·볼륨 | `myportfolio_default`, `myportfolio-master_default`, `workagent-stack`. 볼륨 5개 모두 `myportfolio*`. `knowslink*` 이름은 없다. |
| 자원 | RAM 13 GiB 중 가용 약 7.0 GiB, CPU 4개, 디스크 `/` 가용 325 GB. 이미지 34.6 GB 중 빌드 캐시 22 GB는 회수 가능하나 공유 자원이므로 정리하지 않는다. |
| 로컬 이미지 | `postgres:17-alpine`은 있다. DEV가 고정한 `golang`·`alpine`·`cloudflared 2026.9.1` digest 이미지는 첫 빌드 때 pull한다. |
| cloudflared | 호스트 바이너리 `~/.local/bin/cloudflared` 2026.8.3(2026.9.3 권고 경고). 프로세스 1개가 `orca` Tunnel(`3e132faa-…`)을 `~/.cloudflared/config.yml`로 20일째 실행한다. 연결 3개(icn). |
| 호스트 cloudflared 설정 | ingress: `orca.knowslog.com`·`s8.knowslog.com` → `127.0.0.1:6768`, `mcp.knowslog.com` → `127.0.0.1:8766`, 기본 `http_status:404`. 이 파일을 수정하지 않는다. |
| 인증 접근 범위 | `~/.cloudflared/cert.pem`(0600), `orca` 자격 파일(0400) 존재. `cloudflared tunnel list/info` 성공. Cloudflare API(MCP) 읽기: `knowslog.com` zone `active`, DNS 레코드 목록 200. **쓰기 권한은 시험하지 않았다.** `cert.pem`의 DNS 쓰기 범위는 읽기 전용으로 확인할 수 없다(후속 Dispatch의 실제 Tunnel·DNS 쓰기 결과로만 확인). |
| DNS 공개 조회 | `link.knowslog.com`: A·AAAA·CNAME 없음, 상태 NXDOMAIN. zone에는 `orca`·`s8`·`mcp` CNAME(proxied) 3개만 있다. 와일드카드 레코드 없음. apex 레코드 없음. |
| 공개 HTTPS | `https://link.knowslog.com/` 연결 실패(000). 정상이다. 아직 배포 전이다. |

## 3. 결정: 별도 KnowsLink Tunnel과 컨테이너

기존 `orca` Tunnel과 호스트 `cloudflared` 프로세스는 다른 서비스가 공유한다. 재시작하거나 ingress를 추가하지 않는다. 대신 다음을 사용한다.

1. 새 locally-managed Tunnel `knowslink`를 `cloudflared tunnel create knowslink`로 만든다. 기존 `cert.pem`을 사용한다. 자격 파일은 새 UUID 이름으로 생성되며 `orca` 파일을 건드리지 않는다.
2. DEV 정본 `compose.yaml`의 `cloudflared` 서비스는 `TUNNEL_TOKEN`(remotely-managed Tunnel 전용 `--token`)을 쓴다. 이전 조사에서는 `cert.pem`의 존재와 읽기 접근만 확인했다. 별도 Tunnel 편집 토큰의 존재·쓰기 권한은 미확정이다. 따라서 `compose.yaml`을 고치지 않고 ops 소유 Compose override로 서비스 명령을 `tunnel --no-autoupdate --config /etc/cloudflared/config.yml run`으로 바꾸고 설정·자격 파일을 읽기 전용으로 마운트한다.
3. `postgres`·`relay`·`cloudflared`에만 `restart: unless-stopped`를 추가하는 override를 계획한다. 일회성 `migrate`는 `restart: "no"`를 유지한다. DEV `compose.yaml`의 상시 서비스에는 restart 정책이 없고 migrate는 `restart: "no"`다. 서버 재부팅 뒤 복구하려면 필요하다.
4. 근거: Cloudflare 공식 문서의 locally-managed Tunnel 절차(`tunnel create` → `config.yml` → `tunnel route dns` → `tunnel run`)와 run parameter `--token`이 remotely-managed 전용이라는 설명. 기록 마무리에서 2026-10-03 공식 문서를 재확인했다: [설정 파일](https://developers.cloudflare.com/cloudflare-one/networks/connectors/cloudflare-tunnel/do-more-with-tunnels/local-management/configuration-file/), [run parameters](https://developers.cloudflare.com/cloudflare-one/networks/connectors/cloudflare-tunnel/configure-tunnels/run-parameters/).

Tunnel ingress 설정(값 아님, 구조):

```yaml
tunnel: <knowslink UUID>
credentials-file: /etc/cloudflared/<knowslink UUID>.json
ingress:
  - hostname: link.knowslog.com
    service: http://relay:8080
  - service: http_status:404
```

`relay`는 Compose `ingress` 네트워크에서 서비스 이름으로 접근한다. 호스트 포트를 거치지 않는다.

## 4. 고정 배포 경로와 구성

| 항목 | 값 |
|---|---|
| 배포 체크아웃 | `/home/shin/deploy/knowslink` — 수락 SHA의 detached 체크아웃. 레포 작업 트리와 분리한다. |
| Compose project | `knowslink` (`docker compose -p knowslink`). `myportfolio*` 이름을 쓰지 않는다. |
| 환경 파일 | `/home/shin/deploy/knowslink/.env`, 권한 0600, Git 미추적. 배포 체크아웃의 ignore 규칙을 후속에서 확인한다. `POSTGRES_PASSWORD`·`DATABASE_URL`·`RELAY_PORT=8080`을 둔다. 값은 서버에서 새로 생성하고 `.env.example`의 더미 값을 쓰지 않는다. |
| Tunnel 자격·설정 | `/home/shin/deploy/knowslink/tunnel/{config.yml,<UUID>.json}`. 자격 파일은 0600을 유지한다. 전용 UID 65532 또는 제한된 ACL로 컨테이너 읽기를 허용하고 상위 디렉터리는 0700으로 제한한다. 컨테이너 UID·부모 디렉터리 탐색 권한은 후속에서 확인한다. 공유 자격 파일 권한은 바꾸지 않는다. Git 미추적. |
| 영속 볼륨 | `knowslink_postgres-data` (Compose가 생성). `down -v`를 쓰지 않는다. |
| 게시 포트 | `127.0.0.1:8080` → relay만. Postgres·cloudflared는 게시하지 않는다. |
| 네트워크 | `knowslink_database`(internal), `knowslink_ingress`. 기존 네트워크에 연결하지 않는다. |
| ops 소유 파일(후속 생성) | `deploy/knowslink/compose.ops.yaml`(override), `deploy/knowslink/tunnel/config.yml.tmpl`, 검증 스크립트. 루트 `deploy/` 경로의 소유를 coor가 후속 지시서에서 확인한다. |

## 5. 백업과 rollback

기존 제품 후보 `a6a10c71977b7f3ec8274a1fb7c8a409f58e7c92`에는 실제 SQL migration 00001과 `relay_state` 업무 저장이 있다. goose version table도 생성한다. 빈 migration no-op과 데이터 손실 위험 없음이라는 이전 초안의 설명을 폐기한다. owner/agent credential, gate, inbox, receipt를 포함한 업무 데이터가 있으므로 백업도 민감정보로 취급한다.

1. 배포마다 해당 project의 DB dump를 만든다. 명령은 6장의 고정 Compose 인자를 그대로 사용하여 `exec -T postgres pg_dump -U knowslink -Fc knowslink`를 실행한다. 출력 파일은 `/home/shin/backups/knowslink/<SHA>-<UTC시각>.dump`다. 디렉터리는 0700, 파일은 0600으로 제한한다. dump 종료코드와 `pg_restore --list` 종료코드를 확인한다. 실제 격리 DB 복원 검증은 후속 담당 OPS가 수행한다. 첫 설치에서 DB가 없으면 그 이유를 기록한다.
2. 이전 이미지 digest·SHA·Compose 인자·migration 버전과 백업 경로를 기록한다. schema가 이전 앱과 호환될 때만 이전 SHA의 앱으로 복귀한다. `up`은 migrate를 다시 실행할 수 있으므로 코드만 되돌렸다고 DB rollback 성공으로 기록하지 않는다.
3. schema가 비호환이면 먼저 KnowsLink Tunnel과 relay만 중지한다. 현재 DB를 추가 백업한다. 격리 DB에 dump를 복원해 검증한다. 복원은 백업 이후 변경을 잃을 수 있으므로 후속 복귀 지시의 데이터 복구 범위를 확인한 뒤 KnowsLink DB에 적용한다. 다른 project는 중지하지 않는다.
4. 공개 노출을 중단할 때는 KnowsLink `cloudflared` 서비스만 중지한다. 앱·DB·볼륨은 보존한다. Cloudflare 오류의 정확한 HTTP 코드를 고정하지 않는다.
5. 전체 철회가 명시적으로 지시되면 KnowsLink connector를 중지하고 `link.knowslog.com`의 해당 레코드와 새 Tunnel만 제거한다. 기존 `orca` Tunnel·route는 유지한다. `down -v`와 공유 자원 prune은 금지한다.

백업 보관 기간·암호화 저장소·복구 목표는 운영 재개 전 확정할 항목이다. 임의로 최근 7개 정책을 확정하지 않는다. 업무 행 삭제는 WAL/backup 완전 삭제가 아니다. 백업 삭제 보장을 제품 완료로 표시하지 않는다.

## 6. 적용 절차 (역사적 계획 — 실제 절차는 11장)

당시 계획이며 이 순서로 실행하지 않았다. 11장의 `beta.sh`가 대체한다. 각 단계의 종료코드를 보존한다. 실패하면 중단하고 보고한다.
배포 루트에서 모든 명령에 `docker compose -p knowslink -f compose.yaml -f deploy/knowslink/compose.ops.yaml --env-file /home/shin/deploy/knowslink/.env`를 공통 적용한다. override가 실제로 존재하고 소유권이 확정된 경우에만 사용한다. `config --quiet`로 병합 구성을 검사한다. `config` 전문은 비밀값을 출력하므로 저장하지 않는다. 준비 문서의 명령은 실행 증거가 아니다.

1. 사전 확인: 8080 비어 있음, `link.knowslog.com` 레코드 충돌 없음(Cloudflare DNS API 목록·공개 A/AAAA 조회), `docker compose ls`의 기존 프로젝트 상태를 기록한다. 레코드가 이미 있으면 덮어쓰지 않고 coor에 ask한다.
2. 체크아웃: 수락 SHA를 `/home/shin/deploy/knowslink`에 detached로 둔다. SHA를 기록한다.
3. 환경 파일과 백업 디렉터리를 만든다(권한 0600/0700).
4. DB·앱 기동: 공통 Compose 인자에 `up -d --build --wait relay`를 붙인다. 기대: postgres healthy, migrate 종료코드 0, relay healthy.
5. 로컬 health: `curl -fsS http://127.0.0.1:8080/healthz` → 200.
6. Tunnel 생성: `cloudflared tunnel create knowslink`. 실패 원인이 권한이면 중단하고 coor에 ask한다. 이 단계가 `cert.pem`의 계정 쓰기 권한 증거다.
7. 설정 파일 작성 후 DNS: `cloudflared tunnel route dns knowslink link.knowslog.com`. `--overwrite-dns`를 쓰지 않는다. 기존 레코드가 있으면 실패해야 정상이다.
8. Tunnel 기동: override를 포함해 `--profile tunnel up -d cloudflared`. `docker logs`에서 연결 등록 성공과 연결 상태를 확인한다(비밀값 제외). 정확히 4개 연결을 필수 성공 조건으로 고정하지 않는다.
9. 공개 검증(아래 7장).
10. 기존 서비스 회귀 확인: `myportfolio` 컨테이너 상태·`orca/s8/mcp` HTTPS 응답·호스트 `cloudflared` PID 불변.

## 7. health·인증 확인 방법

| 확인 | 명령 | 기대 |
|---|---|---|
| 컨테이너 | 6장의 공통 Compose 인자에 `--profile tunnel ps -a`를 붙인다 | postgres·relay healthy, migrate Exited(0), cloudflared Up |
| 로컬 health | `curl -fsS -o /dev/null -w '%{http_code}' http://127.0.0.1:8080/healthz` | 200 |
| 공개 health | `curl -fsS -o /dev/null -w '%{http_code}' https://link.knowslog.com/healthz` | 200 |
| DNS | 공개 A/AAAA 조회와 Cloudflare DNS API의 레코드 유형·대상·proxied 속성 확인 | proxied CNAME은 flattening으로 공개 CNAME 응답이 없을 수 있다. CNAME 질의만으로 실패 판정하지 않는다 |
| 미정의 경로 | `curl -s -o /dev/null -w '%{http_code}' https://link.knowslog.com/_probe` | 404 |
| DB 비노출 | `ss -ltn` | 5432 외 새 Postgres 게시 없음. `knowslink` DB 포트 미게시 |
| 인증·권한 | 수락 후보의 D05와 QA-01/08을 기준으로 합성 검사 | 업무 API Bearer 인증, owner UI Basic 인증. 인증 없는 보호 경로 401, agent의 owner 경로 401, 잘못된 CSRF 403. 공개 가입 신원·한도는 held |

## 8. DEV·coor 인계 요청

1. **현재 구현**: 후보 a6a10c7에 owner/agent 역할별 Bearer 인증, owner gate UI의 HTTP Basic·CSRF, 실제 업무 API·SQL 저장이 있다. 합성 가입은 실사용자 신원 증명이 아니다. 공개 신원·가입 보호·DEC-03 한도는 designer와 DEV가 확정·구현·검증해야 한다. Cloudflare Access 도입 여부를 OPS가 상품 규칙으로 확정하지 않는다.
2. **수락 경계**: coordinator는 reviewer `msg_4fbcac80f76c`의 C1 high를 확인했다. agent credential로 human delivery의 pull·persist·ACK·claim을 처리하는 우회가 있다. 담당 DEV의 수정과 독립 tester·reviewer 검증 후 새 SHA를 수락해야 한다.  QA 기록 c59537b의 QA-01–11 실행 항목은 통과했다. 직접 시각 판정·독립 리뷰·coor 수락은 별도 근거가 필요하다. QA 통과를 공개 수락으로 확대하지 않는다. 최신 결과는 후속 Dispatch에서 다시 연결한다.
3. **Compose 인계**: root Compose·Dockerfile·scripts는 DEV 소유다. OPS override의 경로 소유권과 수락 후보의 service/network/image 차이를 후속에서 확인한다. migrate는 일회성으로 유지하고 DB·relay·Tunnel만 restart 정책을 보완한다.
4. **외부 이미지·권한**: 첫 pull과 DNS/Tunnel 쓰기는 미검증이다. 새 Tunnel 생성·DNS 연결 실패 시 기존 레코드를 덮어쓰거나 공유 connector를 수정하지 않는다. coor에게 필요한 범위의 권한 확인을 요청한다. 인증값은 채팅·로그·Git으로 받지 않는다.
5. **보류**: DEC-02 실제 disclosure와 positive silent done, DEC-03 공개 resource/rate/concurrency·신원 정책, 실제 벤더 연결, singleton 처리량, WAL/backup 삭제 보장은 held다. OPS 기술 계획으로 제품 수치를 만들지 않는다.
6. **공유 프로세스**: 호스트 cloudflared 업그레이드·재시작은 이번 범위가 아니다. 운영 컨테이너는 수락 후보의 고정 digest를 사용한다.

DNS 판정 근거: [Cloudflare CNAME flattening](https://developers.cloudflare.com/dns/cname-flattening/). 공개 응답은 대상 IP로 평탄화될 수 있다. 원래 NXDOMAIN 관찰과 후속의 proxied 레코드 존재 판정을 구분한다.

## 9. 완료 상태

| 항목 | 상태 |
|---|---|
| 이전 서버 관찰 보존·운영 계획(D12 초안) | 역사적 기록으로 보존 |
| DNS·Tunnel·컨테이너 변경 | 로컬 컨테이너, `knowslink` Tunnel, Access 앱·정책, `link` proxied CNAME, connector 컨테이너가 생성됐다(11.2 적용 완료). 기존 자원은 불변 |
| 수락 SHA 배포·공개 health 검증 | 로컬 배포는 11장 참조. 공개 연결은 적용됐다(보호 확인 후). 사용자 이메일 로그인 인간 검사는 미실행이며 공개 수락은 held. 기존 a6a10c7 금지는 역사적 기록 |
| D11 사용자설명서·D13 인수인계서 | 작성됨(draft). 공개 연결 전 상태를 반영 |

## 10. 개정 근거와 검증 범위

2026-10-03 retry에서 이전 미추적 초안을 보존하여 수정했다. 지시서 기준 ref는 `0dd08ec994771836c15d9d22a6a83393a71d7987`이며 공통 규칙은 `fullops-common-0.3.2`다.
OPS 체크아웃의 제품 코드는 초기 골격이다. 기존 후보의 제품 정본은 coordinator 체크아웃 `/home/shin/orca/workspaces/KnowsLink/fullops-coor`에서 읽었다. 대상은 README.md, .fullops-squad/project.md, docs/design-docs/architecture.md·interface-design.md·database-design.md, docs/exec-plans/phases/SAR-MVP-001-DEV.md, docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md다. 문서 경로는 `.fullops-squad/` 기준이다. 제품 SHA는 a6a10c7이고 QA 기록 SHA는 c59537b다. 원천 파일은 수정하지 않았다.

당시 단계는 D12 계획 문서 검사만 수행했다. 원래 관찰에는 새 PASS나 종료코드를 붙이지 않았다. 이번 SAR-BETA-001-OPS의 실행 증거는 [실행 기록](../exec-plans/phases/SAR-BETA-001-OPS.md)에 있다.

관련 원천: [D02](../planning/product-specs/SAR-MVP.md), [기능·held](../planning/SAR-MVP-backlog.md), [실행 기록](../exec-plans/phases/SAR-DEPLOY-001-OPS.md), [산출물 인덱스](../deliverables/README.md).

## 11. 본인 전용 합성 베타 (SAR-BETA-001-OPS)

범위는 합성 데이터와 사용자 본인 이메일 한 개다. 이메일은 Git 미추적 `/tmp/knowslink-beta-owner-email`(0600)에서만 읽는다. 문서·로그·Git에 쓰지 않는다. 전체 가입 제품 규칙은 바꾸지 않는다. 이번 접근 집단만 제한한다.

### 11.1 구성 (고정 파일: `deploy/knowslink/`)

| 파일 | 역할 |
|---|---|
| `compose.ops.yaml` | restart `unless-stopped`, 자원 제한(postgres 512m·1 CPU, relay 256m·0.5 CPU, cloudflared 128m), 로그 10m×3, relay/migrate/cloudflared의 cap_drop ALL·no-new-privileges·read_only. 값은 기술 설정이며 상품 정책이 아니다 |
| `tunnel/config.yml.tmpl` | `link.knowslog.com` → `http://relay:8080`, 나머지 404. `originRequest.access.required: true`와 teamName·audTag로 원점에서 Access JWT를 검증한다 |
| `beta.sh` | prepare·up·stop·unexpose·backup·restore-verify·seed·owner-login·tunnel-create·render-config·expose |
| `access_apply.py` | 사용자 이메일만 allow인 reusable policy와 self-hosted app을 만든다. 기존 앱·정책·IdP를 덮어쓰지 않는다 |
| `verify.py` | baseline·regression(공유 서비스)·local·public(미인증 negative) |

배포 체크아웃은 `/home/shin/deploy/knowslink`(detached)다. 상태는 `/home/shin/deploy/knowslink-state`(0700)에 둔다: `.env`(0600), `tunnel/`, `backups/`, `access.aud`. 체크아웃은 항상 깨끗하다. Compose project는 `knowslink`다.

### 11.2 노출 순서 (보호 먼저)

1. 로컬 기동·검증·백업·격리 복원을 끝낸다. 공개 DNS는 없다.
2. `CF_API_TOKEN_FILE=<0600> python3 access_apply.py apply`로 Access 앱을 만든다. `access.aud`가 저장된다. 필요한 권한은 `Access: Apps and Policies Edit`(쓰기)와 `Access: Organizations, Identity Providers, and Groups Read`(IdP 목록 읽기)다. 토큰은 계정 하나만 볼 수 있어야 한다. 관리 MCP로 쓰는 경우 같은 두 권한이 필요하다.
3. `beta.sh render-config`로 AUD를 넣은 `config.yml`을 만든다. `cloudflared tunnel ingress validate`가 통과해야 한다.
4. `beta.sh expose`: `access_apply.py check`가 live 앱·정책을 읽어 도메인·destinations·정책 1개(사용자 이메일 단독 allow)·IdP 1개·우회 옵션 없음·`aud`==기록값==Tunnel `audTag`·`teamName`을 검증한다. 토큰이 없으면 10분 이내의 읽기 전용 MCP 응답을 `access.live.json`으로 저장해 같은 검증을 쓴다. 이어서 relay health와 DNS 부재(`dig` 실패도 중단)를 확인한 뒤에만 `tunnel route dns`(덮어쓰기 없음)와 connector를 시작한다.
5. `verify.py public`과 `verify.py regression`을 실행한다. 사용자가 마지막 이메일 로그인을 직접 확인한다(인간 검사).

### 11.3 중단·복귀

- 구성 갱신·코드 rollback: `beta.sh deploy <sha>`. 먼저 백업(0600)하고 `deploy-history.log`에 이전→새 SHA를 남긴 뒤 detached 체크아웃을 옮겨 재빌드하고 `verify.py local`을 실행한다. `db/migrations`가 현재와 다르면 중단한다. 이 경우 본 문서 5장 3항(격리 DB 검증 뒤 복원)을 따른다. rollback도 같은 명령에 이전 SHA를 준다. 현재 SHA의 정본은 `git -C /home/shin/deploy/knowslink rev-parse HEAD`다. `deploy-history.log`는 `beta.sh deploy`로 옮긴 기록만 담고 수동 체크아웃은 담지 않는다. `8a7ad36` 이전 SHA에는 `access_apply.py check`가 없다. 그 SHA의 `beta.sh expose`는 live 검증이 없는 이전 게이트다. 그러므로 `expose`는 `check` 명령이 있는 SHA에서만 실행한다. 이전 SHA의 `beta.sh`에 `deploy`가 없으면 한계가 있다. 이 경우 `git -C /home/shin/deploy/knowslink checkout --detach <sha>`로 직접 옮기고 `beta.sh up`을 실행한다. 2026-10-03에 `437f143`으로의 rollback과 앞으로의 이동을 실제 실행했다(둘 다 exit 0, 백업 생성, 이력 기록).
- 노출 중단: `beta.sh unexpose`(connector만 중지). 앱·DB·볼륨을 보존한다.
- 전체 중지: `beta.sh stop`. `down -v`와 공유 자원 prune은 금지한다.
- 전체 철회: `access_apply.py remove`, Cloudflare 대시보드에서 `link` CNAME과 `knowslink` Tunnel을 삭제한다. 기존 `orca` Tunnel·호스트 cloudflared는 건드리지 않는다.
- 데이터 복구: `beta.sh backup`(0600 dump)과 `beta.sh restore-verify <dump>`(네트워크 없는 임시 컨테이너에 복원)로 검증한다. 라이브 DB에는 적용하지 않는다.

### 11.4 held

공개 한도(DEC-03)·실제 신원·실데이터·실제 벤더·무제한 공개·WAL/backup 삭제 보장은 held다. 자동 API 검사에 Access 토큰이 필요하면 사용자 정책을 넓히지 않고 단기 service token 승인을 별도로 요청한다.

### 11.5 비밀 취급 주의

`beta.sh owner-login`은 합성 Basic 값을 터미널에 출력한다. 사용자 본인만 자기 터미널에서 실행한다. agent·자동화는 실행하지 않는다. 출력이 세션 기록에 남는다. `access_apply.py`는 토큰·이메일·응답 본문을 출력하지 않는다.

### 11.6 적용 결과 (2026-10-03, 구성 `28bd1bb`)

11.2의 1~5단계를 순서대로 실행했다. 2단계는 `access_apply.py`와 같은 본문 함수를 공식 Cloudflare MCP로 호출해 수행했다(`access.live.json`은 읽기 GET 응답). 결과·증거는 [실행 기록](../exec-plans/phases/SAR-BETA-001-OPS.md)에 있다. 로컬 resolver가 NXDOMAIN을 캐시하면(SOA TTL) `verify.py public`이 `Name or service not known`으로 실패할 수 있다. 이때는 `dig @1.1.1.1`로 엣지 IP를 얻어 `curl --resolve`를 쓴다.

## 12. Grok Bot 앱 connector 등록 (SAR-MVP-002-BOT-CATALOG-DEV)

1. owner는 Bot 컴퓨터의 레포 루트에서 `sh scripts/install_bot_mcp.sh`를 실행한다. 결과: 마지막에 Name·Type·Command·Arguments 등록 값이 출력된다.
2. owner는 Bot 채팅에서 custom MCP server 추가를 요청하고 **Add MCP Server** 카드의 Type **Command**와 빈 환경 변수를 확인한 뒤 승인한다.
3. owner는 새 대화에서 `knowslink_status`만 호출한다. 결과: `held`.

중단은 Installed 목록에서 knowslink를 삭제한다. 실제 relay·환경 변수·비밀값은 등록하지 않는다. 상세 절차는 [플러그인 문서](../../../adapters/README.md#grok-bot-앱-등록)를 따른다.
