---
id: D12
title: 운영자설명서
status: draft
updated: 2026-10-10
owner: ops
tasks: [SAR-DEPLOY-001-OPS, SAR-BETA-001-OPS, SAR-MVP-002-BOT-CATALOG-DEV, SAR-MVP-002-BOT-CATALOG-DEV-FIX, SAR-MVP-003-BIDIRECTIONAL, SAR-MVP-003-BIDIRECTIONAL-OPS, SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW, SAR-PUBLIC-SERVICE-OPS-READINESS, SAR-PUBLIC-SERVICE-OPEN-PREP, SAR-GOOGLE-CONNECT-001-DEV, SAR-GOOGLE-CONNECT-002-DEV]
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
2. owner는 Bot info pane의 **Setup → Plugins → Add**에서 custom MCP server를 추가한다. 그 화면이 없으면 Bot 채팅에서 추가를 요청한다. Type **Command**와 빈 환경 변수를 확인한 뒤 승인한다. 결과: Installed 목록에 knowslink가 보인다. 이 화면과 승인 카드는 공식 근거가 Team Bots에 한정되며 실제 계정에서 미확인이다.
3. Command 등록 수단이 없거나 Bot이 `AddMcpServer` 같은 도구가 없다고 답하면 owner는 등록을 멈춘다. 실제 메뉴 항목·Bot 응답·도구 목록을 회신한다. 이슈1에서 같은 Bot은 `AddMcpServer`를 호출할 수 없었다.
4. owner는 새 대화에서 `knowslink_status`만 호출한다. 결과: `held`.

컴퓨터 update·recover·Reset 뒤 `/workspace/.knowslink`의 Node 또는 bundle이 없으면 1단계를 다시 실행한다. `/workspace` 유지는 보장되지 않는다. 앱 등록 값은 같다. 중단은 Installed 목록에서 knowslink를 삭제한다. 실제 relay·환경 변수·비밀값은 등록하지 않는다. 상세 절차는 [플러그인 문서](../../../adapters/README.md#grok-bot-앱-등록)를 따른다.

## 13. 승인된 양방향 시험 (SAR-MVP-003)

이번 사용자 요청은 Codex↔Grok 시험 text와 기존 `https://link.knowslog.com` 재사용을 승인했다. 11.4의 실제 벤더 held 중 이 시험 범위만 재개한다. 업무 effect·disclosure·dots·FullOps 갱신·공개 가입 제품 정책은 바꾸지 않는다. DEV 준비물과 실제 적용·계정 QA는 분리한다.

### 13.1 확인된 현재 상태

2026-10-04 DEV가 기존 배포 HEAD `28bd1bb2bdd4b90d0c6f5d2a8d4ba3e00f630ff2`와 relay/postgres healthy·knowslink Tunnel Up을 확인했다. 공개 registry는 Access 302다. Cloudflare GET은 기존 `KnowsLink beta (owner-only)` 앱 한 개와 owner allow policy 한 개를 반환했다. service token은 0개이며 trial path 앱은 없다. 실제 배포는 새 trial API를 포함하지 않는다.

DEV는 기존 로컬 beta에서 `trial_codex`와 `trial_grok`의 신규 시험 identity·Ed25519 PoP·pair를 준비했다. `/home/shin/deploy/knowslink-state/trial-SAR-MVP-003-BIDIRECTIONAL`은 0700이다. 두 agent의 하위 `environment.json`과 `key.pem`은 0600이다. `owners.json`은 운영자에게만 둔다. 원격 agent에 owner 기록을 보내지 않는다. Access ID/secret은 아직 없다. 이 등록은 합성 beta owner 등록이며 실제 Grok 계정 신원 검증이 아니다.

### 13.2 고정 SHA 배포와 machine Access 준비

1. coor는 전체 40자리 `REVIEWED_SHA`의 독립 리뷰·QA·필수 gate를 확인한다. 기존 11.3의 `beta.sh deploy "$REVIEWED_SHA"`로 백업 후 기존 knowslink 배포를 이동한다. DB migration은 변경하지 않았다. 공유 서비스와 볼륨은 보존한다.
2. OPS는 상태 `.env`에 `KNOWSLINK_TEST_AGENTS=trial_codex,trial_grok`만 추가한다. 0600을 유지한다. 기존 Compose 인자로 relay를 갱신한다. 명시 allowlist가 없으면 trial을 거부하는 것이 정상이다.
3. 기존 계정에서 `knowslink-trial-codex`, `knowslink-trial-grok` 이름의 서로 다른 service token을 24h 유효기간으로 만든다. 생성 본문은 `python3 deploy/knowslink/access_trial_plan.py tokens`로 확인한다. API 경로는 `/accounts/{account_id}/access/service_tokens`다. 생성 응답의 client secret은 한 번만 안전한 로컬 파일에 저장한다. 응답 전문을 터미널·issue·chat·Git에 기록하지 않는다. 필요한 계정 권한은 Access service-token 쓰기와 Apps/Policies 쓰기다. 새 서비스·DNS·Tunnel·유료 자원은 만들지 않는다.
4. 실제 두 token UUID로 `python3 deploy/knowslink/access_trial_plan.py policy --codex-token "$CODEX_TOKEN_UUID" --grok-token "$GROK_TOKEN_UUID"`를 실행한다. reusable policy 생성 경로는 `/accounts/{account_id}/access/policies`다. decision은 `non_identity`(Service Auth)이며 include는 이 두 token만 사용한다. `any_valid_service_token`, bypass, everyone은 사용하지 않는다.
5. 생성 policy UUID로 `python3 deploy/knowslink/access_trial_plan.py app --policy "$TRIAL_POLICY_UUID"`를 실행한다. `/accounts/{account_id}/access/apps`에 별도 self_hosted 앱을 생성한다. domain과 public destination은 `link.knowslog.com/v1/test/*`다. 기존 root owner 앱·이메일 policy·IdP·team 조직 설정을 변경하지 않는다. 생성 뒤 GET으로 destination·policy ID·decision·두 token UUID·trial AUD를 대조한다. 기존 root 앱과 trial 외 service-token 거부도 확인한다.
6. `python3 deploy/knowslink/access_trial_plan.py tunnel --uuid "$KNOWSLINK_TUNNEL_UUID" --team scshin88 --owner-aud "$OWNER_AUD" --trial-aud "$TRIAL_AUD"`로 후보를 생성한다. 기존 Tunnel을 재사용한다. `/v1/test/.*` 첫 rule은 trial AUD만, 나머지 hostname rule은 기존 owner AUD만 검증한다. 둘 다 `access.required: true`다. 후보를 0600에 두고 현재 cloudflared UID로 pinned image의 `ingress validate`와 `ingress rule`을 검증한다. 상태 config를 백업 후 교체하고 knowslink connector만 재기동한다. 기존 `beta.sh render-config`는 시험 entry를 만들지 않으므로 활성 시험 중 사용하지 않는다. 기존 `expose`는 DNS 최초 생성용이므로 재실행하지 않는다. 이번 path 후보는 11.2의 단일 AUD 검사로 대체 검증하지 않는다.
7. 공개 negative는 credential 없이·잘못된 service token으로 시험 API를 거부하는지 확인한다. 실제 token+잘못된 relay agent token도 거부해야 한다. 실제 두 credential로 registry/keys/send/pull을 검증한다. token만으로 owner/signup/business API를 열 수 없어야 한다. 원격 성공은 이 검증과 실제 Grok 왕복 ID가 확보된 뒤에만 기록한다.

### 13.3 안전한 credential 전달과 실행

운영자는 각 `environment.json`에 그 agent 전용 `CF_ACCESS_CLIENT_ID`와 `CF_ACCESS_CLIENT_SECRET`을 넣는다. Codex와 Grok credential을 섞지 않는다. 서버 관리자가 승인한 secret 전달 채널로 Grok의 자기 key/config만 전달한다. Grok 파일 경로는 `/workspace/.knowslink-trial/trial_grok/key.pem` 같은 실제 절대경로로 바꾸고 config의 `AGENT_KEY_FILE`도 맞춘다. parent 폴더 0700·파일 0600을 확인한다. private host 파일을 issue 첨부로 보내지 않는다.

Codex 수신:

```sh
python3 scripts/run_trial.py --config /home/shin/deploy/knowslink-state/trial-SAR-MVP-003-BIDIRECTIONAL/trial_codex/environment.json receive
```

Codex 송신은 표준 입력에 시험 text만 주고 `send sar-mvp-003-codex-round-1`을 사용한다. uncertain retry는 같은 text와 key만 사용한다. Grok Command 등록은 Python launcher·config 경로·고정 Node 경로·standalone bundle 경로만 argv에 넣는다. 비밀값은 argv에 넣지 않는다. owner/admin이 기존 `user-knowslink` 등록의 Command·환경 지원을 확인한 뒤 교체한다. 미지원이면 같은 launcher의 CLI 수동 pull을 사용한다. Grok 부모 컴퓨터의 hostname/loopback을 서버로 오인하지 않는다.

Grok 수신 ID는 Codex send ID와 같아야 한다. Grok 회신 text에는 첫 ID를 포함한다. Grok reply ID는 Codex receive ID와 같아야 한다. 각 단계는 180초 TTL 안에 수행한다. `empty`는 네트워크 성공이나 유실 확정이 아니다. 기본 모드를 held로 바꾸면 송수신 도구가 멈춘다. 네트워크 queue가 대화를 자동 wake하지 않는다.

### 13.4 중단·복귀와 차단 조건

시험 종료 시 machine token을 disable/revoke하고 trial path 앱·policy만 제거한다. 기존 owner Tunnel config 백업을 복원한다. 상태 `.env`의 시험 allowlist를 비워 relay를 갱신한다. 필요하면 `owners.json`의 자기 owner credential을 운영자 로컬에서만 사용해 시험 owner를 revoke한다. pair·키·claim 권한의 철회를 확인한다. 기존 owner Access·Tunnel·DNS·업무 DB·공유 서비스를 삭제하지 않는다.

현재 차단은 trial service token과 path 앱·trial AUD의 부재, 후보 미배포, Grok parent의 환경 변경/secure-file 지원 회신 미도착이다. secret 안전 전달 채널도 아직 확인하지 않았다. DEV는 이 값을 추측하거나 Access를 해제하지 않았다. 신규 비용은 발생하지 않았다. 코드·로컬 왕복·배포 절차·댓글 초안은 [실행 기록](../exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md)에 연결한다.

근거: [Cloudflare Service tokens](https://developers.cloudflare.com/cloudflare-one/access-controls/service-credentials/service-tokens/)의 두 header·Service Auth, [Application paths](https://developers.cloudflare.com/cloudflare-one/access-controls/policies/app-paths/)의 구체 path 우선순위, [Grok Team Bots](https://docs.x.ai/grok-bot/team-bots)의 Command 실행 위치·owner chat secrets 경계를 확인했다. API 본문은 2026-10-04 Cloudflare OpenAPI search로 대조했다. 문서 지원을 실제 계정 지원으로 확대하지 않는다.

### 13.5 실제 운영 결과 (2026-10-04, SAR-MVP-003-BIDIRECTIONAL-OPS)

- 완료: 기존 배포를 `0911c2c73468f8684260a277d4940a74d26bcf7d`로 이동했다(`beta.sh deploy` exit 0, DB 백업 `28bd1bb-20261004T044744Z.dump`, migration 동일). 상태 `.env`에 `KNOWSLINK_TEST_AGENTS=trial_codex,trial_grok`만 추가했다. 실제 두 시험 key로 loopback 왕복을 확인했다. 공유 서비스 회귀와 공개 무인증 probe(302)는 불변이다.
- 차단: 13.2의 3~7단계. 저장된 Cloudflare 권한 어디에도 `access-service-token.write`가 없다. 세션 MCP는 읽기 전용이고 Codex OAuth grant에는 service-token 범위가 없다. root owner 앱과 정책은 변경하지 않았다.
- 준비: 요청 본문·원점 후보 구조 검증·Grok 전용 private 파일(`grok-export/`, CF 항목 제외)·재개와 종료 절차·Grok 최종 댓글 초안이 있다.
- 미검증: 공개 HTTPS 성공·공개 negative/positive·실제 Grok 왕복. 상세와 증거는 [실행 기록](../exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md)이다.

### 13.6 차단 해소와 공개 검증 (2026-10-04, 재개 Dispatch)

13.5의 차단은 해소됐다. 13.4·13.5의 당시 차단 기록은 보존한다. 범위 한정 관리 token(Service Tokens Edit·Apps and Policies Edit, 계정 하나, 만료 2026-10-04T23:59:59Z)으로 13.2의 3~7단계를 수행했다.

- 적용: service token 2개(`knowslink-trial-codex`·`knowslink-trial-grok`, 24h, 만료 2026-10-05T06:42:09Z), reusable policy `knowslink-trial-agents`(`non_identity`, 두 token만), self_hosted 앱 `KnowsLink trial messages`(`link.knowslog.com/v1/test/*`). 원점 Tunnel은 `/v1/test/.*` trial AUD rule을 앞에 둔 후보로 교체했고 백업 `tunnel-bak-pre-trial`을 유지한다. 앱·policy 변경은 root owner 앱·정책·IdP·team에 영향이 없다.
- 공개 negative: 무인증·잘못된 CF 403, 유효 CF+잘못된 relay credential 401, token만으로 owner·pair·business 경로 302. 공개 positive: 두 local client가 registry/keys/send/pull로 합성 text를 왕복했다. 실제 Grok은 아니다.
- 종료: token 두 개만 revoke하고 앱·policy만 삭제한다. 원점 config를 복원하고 allowlist를 비운다. 관리 token이 먼저 만료되므로 삭제에는 새 권한이나 대시보드가 필요하다.
- 현재 차단 중 남은 것: Grok 전용 파일의 안전 전달과 Grok 설치 회신이다. 증거·ID·절차는 [실행 기록](../exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md)의 「재개 결과」를 따른다.

### 13.7 만료 갱신 (2026-10-05, SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW)

13.6의 service token 두 개는 2026-10-05T06:42:09Z에 만료됐다. 새 범위 한정 관리 token(Service Tokens Edit·Apps and Policies Edit, 계정 하나, 만료 2026-10-05T23:59:59Z)으로 두 token에 `PUT /accounts/{account_id}/access/service_tokens/{id}`를 `duration: "24h"`로 보냈다. UUID와 client secret(`client_secret_version` 1)은 바뀌지 않았다. 새 만료는 codex 2026-10-06T08:09:43Z, grok 2026-10-06T08:09:50Z다. 공식 `refresh`는 만료를 1년 연장하므로 24h 범위에서 쓰지 않았다. 앱·policy·Tunnel·`environment.json`의 CF 항목은 변경하지 않았다.

- 공개 negative: 무인증·잘못된 CF 403, 유효 CF+잘못된 relay credential 또는 없음 401, owner·pair·signup·업무 경로 302. positive: registry 200(두 agent), 두 local client의 HTTPS 왕복 통과(실제 Grok 아님).
- Grok 전달 묶음: `/home/shin/deploy/knowslink-state/trial-SAR-MVP-003-BIDIRECTIONAL/knowslink-grok-trial-20261005.tar.gz`(0600). `key.pem`·`environment.json` 두 개만 담는다. 외부 전달은 하지 않았다.
- 관리 token은 service token보다 약 8시간 먼저 만료된다. 종료 정리(token revoke·앱·policy 삭제)는 그 전에 하거나 새 권한을 받는다. 증거와 Grok 댓글 초안은 [실행 기록](../exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS-RENEW.md)이다.

### 13.8 시험 종료 (2026-10-05, SAR-MVP-003-TRIAL-CLEANUP)

실제 Grok 왕복이 확정된 뒤 13.4의 종료 절차를 실행했다. service token 두 개를 revoke하고 trial 앱·policy만 삭제했다. 원점 Tunnel config를 백업에서 복원(0600)하고 knowslink connector만 재기동했다. 상태 `.env`의 `KNOWSLINK_TEST_AGENTS` 줄을 제거하고 relay를 같은 Compose 인자로 갱신했다. root owner 앱·policy 3건은 전체 JSON 해시가 변경 전과 같다. 공유 서비스 회귀(`verify.py regression`·`public`·`local`)는 통과했다. 사용자 private 시험 파일은 삭제하지 않았고 안의 CF 자격은 더는 연결되지 않는다. 관리 token은 2026-10-05T23:59:59Z에 만료된다. 증거는 [실행 기록](../exec-plans/phases/SAR-MVP-003-TRIAL-CLEANUP.md)이다.

## 14. 일반 이메일 서비스 운영 준비 근거 (SAR-PUBLIC-SERVICE-OPS-READINESS)

이 장은 읽기 전용 관측과 제안이다. 서버·Cloudflare 설정·제품 코드는 바꾸지 않았다. 서버 읽기 전용 관측은 2026-10-05T12:11:29Z부터 약 12:14Z까지, 배포 체크아웃 `0911c2c73468f8684260a277d4940a74d26bcf7d` 기준이다. 증거와 전체 표는 [실행 기록](../exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md)이다. 13장의 시험 종료 상태는 유지된다. 이 장은 공개 수락 근거가 아니며 D13을 수락으로 바꾸지 않는다.

### 14.1 실제 신원 보호 (관측)

- Access 앱은 `KnowsLink beta (owner-only)` 한 개다. `link.knowslog.com` 전체를 덮고 정책은 사용자 이메일 한 개만 허용한다. 무인증 경로는 모두 302 → Access 로그인이다. service token과 trial 앱·policy는 없다.
- relay는 Access JWT·이메일을 읽지 않는다. Bearer(API)와 Basic(owner UI)만 읽는다. 그래서 relay는 이메일 신원을 쓰지 않고 합성 owner credential과 실제 이메일이 연결되지 않는다. 공식 문서상 Access는 `Cf-Access-Jwt-Assertion` 헤더로 JWT를 보내고 cloudflared가 `access.required`에서 검증한다. claim은 relay로 전달될 수 있다. 이 서버에서 헤더 도달은 미확인이다([근거 URL](../exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md)).
- 그러므로 root를 모든 이메일 허용으로 넓혀도 이메일 신원이 relay 인증에 반영되지 않고 익명 사용자가 Bearer/Basic 앞단까지 도달한다. relay가 claim을 읽어 owner에 묶는 코드(DEV)가 배포되기 전에는 넓히지 않는다.
- Managed OAuth는 꺼져 있다(`oauth_configuration` 비어 있음). Cloudflare 공식 문서는 비브라우저 client용 Managed OAuth(2026-03-20)를 설명한다. 실제 Grok·OpenAI dot client와의 호환은 미확인이다.

### 14.2 자원·DB·복구 (관측)

- 호스트 4 CPU, RAM 13,906 MiB(가용 6,002), swap 4,095 MiB 중 4,074 사용, `/` 가용 312 GB. swap 압력은 개발 도구 프로세스에서 온다. KnowsLink 컨테이너는 한도의 약 3–7%를 쓴다(유휴). `myportfolio` 컨테이너에는 한도가 없다.
- DB는 8.3 MB다. 모든 업무 상태가 `relay_state` 한 행(7,202 byte)의 jsonb다. 처리량·행 크기 상한은 미측정이다. 부하 시험은 DEV 소유다.
- 자동 백업이 없다. 최신 dump는 `0911c2c` 배포 직전(2026-10-04T04:47Z)이다. 같은 디스크에만 있다. 보관 기간·암호화·복구 목표는 정해지지 않았다.

### 14.3 공개 전 선행 조건과 rollback 순서

선행 조건: 신원 방식과 agent·MCP 연결 방식의 제품 결정, relay의 JWT 신원 연결, 가장자리 rate limit·bot 보호 결정, 백업 체계(일정·별도 위치·암호화·복원 시험), DEC-03 한도 구현 검증, 독립 QA·리뷰. 관리 token은 현재 `active`이고 2026-10-05T23:59:59Z에 만료 예정이다. 새 발급은 지금 필요하지 않다. 만료 뒤 작업이나 부족한 권한(IdP·조직·zone 읽기)이 확인된 때만 요청한다.

### 14.4 최종 운영 E2E 조건

최종 운영 E2E는 사용자 본인의 동일한 일반 이메일 한 개로 가입한 한 owner 아래에서 Grok Bot 노우와 OpenAI dot 다닷 두 agent를 연결해 수행한다. 이메일 두 개는 필수가 아니다. 이전 trial token과 allowlist는 되살리지 않는다. 공식 문서 URL·조회일(2026-10-05)과 미확인 항목은 실행 기록 8장이다.

적용 순서: 새 백업과 복원 검증 → `beta.sh deploy <SHA>` → 새 경로 Access 앱·정책(기존 root 앱 불변) → 후보 `config.yml` 검증·교체 → negative/positive·회귀 검증 → 노출 확대는 마지막.

rollback 순서(역순): 새 앱·정책 삭제 → `config.yml` 백업 복원과 connector 재기동 → 필요하면 `beta.sh unexpose` → 코드 `beta.sh deploy <이전 SHA>`(migration 동일일 때만, 아니면 5장 3항) → 공유 서비스는 건드리지 않는다.

## 2026-10-07 최종 코드 수락과 공개 준비

사용자 요청으로 최종 main 수락과 공개 준비를 진행한다. [공개 준비 실행 기록](../exec-plans/phases/SAR-PUBLIC-SERVICE-OPEN-PREP.md)의 고정 SHA·실제 백업/격리 복원·현재 보호상태·배포/복귀 절차를 따른다. 기존 배포0911c2c는유지됐고새서비스/실메일은아직미검증이다. 과거시험의미완료/완료는원래시점으로보존한다. 이 연결은 운영 공개 수락이 아니다.

## Google 로그인 우선 출시 — 2026-10-07

최신 사용자 결정은 SMTP 대신 Google 로그인으로 핵심 기능을 먼저 출시하는 것이다. 과거 이메일 전용 공개 준비·실메일 미검증은 원래 기록으로 유지한다. 새 프로젝트 knowslink-auth의 Web OAuth client를 등록했고 지정 계정은 testing 사용자다. 결제 설정을 활성화하지 않았다.

운영 env 정본은 `/home/shin/deploy/knowslink-state/.env`0600이다. `KNOWSLINK_GOOGLE_CLIENT_ID`, `KNOWSLINK_GOOGLE_CLIENT_SECRET`, `KNOWSLINK_GOOGLE_REDIRECT_URL`을 Compose가 relay에 전달한다. 값은 Git/대화/로그에 남기지 않는다. callback은 `https://link.knowslog.com/auth/google/callback`이다. Google discovery/token/JWK HTTPS egress가 필요하다. 기존 Tunnel·Access와 Postgres volume은 유지한다.

제품 후보017bf456의 독립 인증 delta 리뷰와 기존 lint/test를 확인했다. main 통합 뒤 `beta.sh deploy <수락 SHA>`로 사전 DB backup·migration diff 거부·기동 확인을 수행한다. 이번 변경은 DB migration이 없다. 이전 제품 배포d08903a55c3638128827010400e66e9d45b61d7c로 코드 복귀할 때는 새 Google env 세 개를 비공개 환경에서 비활성화한다. 실제 사용자 로그인은 배포 뒤 `/`의 Google로 계속 → Google 화면의 선택/동의 → 자기 홈 링크로 확인한다. 일반 공개·Grok Bot 실연결 완료는 별도 기록한다.

## Google 연결 공개 경로 적용 — SAR-GOOGLE-CONNECT-001-DEV

이 절은 OPS 실행 계획이다. DEV는 운영 계정·환경·Access·Tunnel을 변경하지 않았다. 이전 owner-only `access_apply.py`나 beta Access 초기화 도구를 그대로 재실행하지 않는다. 먼저 현재 Access application ID·aud·destinations·정책·IdP와 Tunnel config를 비밀값 없이 백업하고 shared/trial 호스트 기준 상태를 기록한다. Google redirect는 기존 `/auth/google/callback`이며 새 scope·secret은 없다.

1. 기존 호스트 전체 owner Access application의 보호 대상을 `link.knowslog.com/owner`와 `link.knowslog.com/v1`로 좁힌다. 같은 application ID·aud·owner 정책·IdP를 유지하고 실제 API 응답에서도 보존을 확인한다. 이전 도메인이 남아 있으면 `/owner` 보호를 함께 유지한다. 호스트 전체 Everyone 허용은 만들지 않는다.
2. 더 구체적인 공개 경로 `/v1/connect`, `/v1/text`, `/v1/keys/*`, `/v1/receipts/*`에 Bypass Everyone application을 구성한다. 정확한 관리자 POST `/v1/keys`는 기존 owner 보호 아래 둔다. 기존 `/v1/test/*` trial application과 aud는 보존한다. 경로 우선순위를 실제 Access 설정과 요청으로 확인한다.
3. `/`, `/auth/*`, `/home`, `/home/*`, `/connect/*`가 다른 광범위 application에 잡히지 않게 한다. 서비스는 자체 Google 세션·Origin·최근 인증·회원 소유권 검사를 계속 적용한다. 세션 없는 홈은 서비스 로그인으로 이동해야 한다.
4. `deploy/knowslink/tunnel/public-ingress.yml` fragment를 기존 owner 보호 fallback 앞에 삽입한다. 기존 tunnel UUID·credentials 경로·trial 및 다른 호스트와 `required:true` fallback을 보존한다. origin loopback, 비공개 Postgres, 관리자 Basic/Bearer 보호도 유지한다.
5. 변경 파일에 `cloudflared tunnel ingress validate`와 `cloudflared tunnel ingress rule <URL>` 검사를 수행한다. 공개 경로와 owner/admin/test/미등록 경로를 각각 확인한 뒤 두 계층을 함께 적용한다. fragment 자체는 edge Access application을 바꾸지 않는다.

Cloudflare의 [Access 경로 우선순위](https://developers.cloudflare.com/cloudflare-one/access-controls/policies/app-paths/)와 [Tunnel 첫 일치 규칙](https://developers.cloudflare.com/cloudflare-one/networks/connectors/cloudflare-tunnel/do-more-with-tunnels/local-management/configuration-file/)을 근거로 한다. Context7 조회는 월간 quota 초과여서 공식 문서를 확인했다. DEV의 regex self-check는 실제 edge 검사를 대신하지 않는다.

좁은 운영 검사는 익명 root/connect의 Cloudflare 로그인 redirect 부재, 실제 Google callback·동의·로컬 저장·자기 홈 한 흐름, 세션 없는 홈의 서비스 로그인, 잘못된 token/proof의 401/422 JSON을 확인한다. `/owner`, `/v1/owners`, 정확한 관리자 `/v1/keys`, `/v1/authorize`의 보호와 위조 JWT 거부를 확인한다. 기존 trial/shared 호스트 회귀를 확인한다. 실제 외부 Bot 설치와 같은 Google 계정의 별도 두 키·관계 수락·왕복 전달은 후속 수락 증거로 남긴다.

복구할 때 먼저 기존 보호 ingress와 Access 대상을 함께 복원하여 새 공개 연결을 중단한다. 이전 수락 코드로 되돌리되 현재 DB·키·credential·폐기 기록은 보존한다. 과거 DB 백업을 덮어쓰지 않는다. 이전 코드가 새 JSON 필드를 버릴 수 있으므로 대기 요청은 새로 시작하고, 완료된 키와 폐기 상태를 별도 확인한다.
## 비용 없는 Tunnel 적용 — SAR-GOOGLE-CONNECT-002-DEV

이 절이 새 Google 연결의 운영 적용 정본이다. 위 SAR-GOOGLE-CONNECT-001의 Access 앱 생성·Bypass 계획은 Access를 유지하는 별도 선택지의 이력이다. 이번 경로는 Access 가입·유료 요금제·초과 자동 과금 동의 없이 기존 Tunnel과 도메인을 사용한다. 공개 고정 IP를 구매하지 않는다. 기존 Google callback과 운영 .env·DB·키·credential을 보존한다. DEV는 후보만 제공하며 실제 적용은 coor가 수행한다.

coor msg_3319725732d2에서 로그인된 dashboard API GET은 success였다. 현재 대상은 `fc81b205-d4d1-445e-a2bd-384a2ed82f62` member agent API와 `bd210310-fd5e-4e6e-8cb4-d36d618cbebd` owner-only 앱이다. root owner 앱이 Google/home을 막는다. 과거 trial 앱 ID는404이며 새 Access 정책을 만들 필요가 없다. 앱 UI의 요금제 gate는 무료 활성화가 필수라는 근거로 사용하지 않는다. API 변경도 현재 세션 권한·응답을 다시 확인하고 이 두 KnowsLink 앱에만 한정한다.

1. 고정 SHA 독립 리뷰·좁은 QA 후 현재 배포 코드와 Tunnel config 및 두 앱의 full app/policies를 비공개 위치에 백업한다. coor가 확보한 원본 config SHA256은 `40ce65ea1450329ac73c9b3188f12966b99ea562fae31c334e06dd92b06cf283`이다. 적용 직전 일치 여부를 확인한다. 앱 백업 SHA256은 coor 비공개 증거를 따른다. 다른 호스트·서비스의 기준 상태도 기록한다.
2. `KNOWSLINK_DEPLOY=<수락 checkout> KNOWSLINK_STATE_DIR=<기존 비공개 상태> bash deploy/knowslink/beta.sh render-public-config`를 실행한다. `tunnel/config.public.yml`0600만 새로 만들고 기존 파일이 있으면 거부한다. 기존 config·DNS·컨테이너·Access는 변경하지 않는다. 검증 실패 후보는 적용하지 않는다. 서버 PATH에 cloudflared가 없으면 기존 pinned cloudflared 이미지의 동일 바이너리를 비공개 임시 도구 경로에 준비하거나 후보를 같은 이미지의 `tunnel --config <후보> ingress validate`로 검사한다. 새 유료 도구는 필요 없다.
3. 후보는 member allowlist 하나와 `http_status:404` catch-all이다. `/owner`, 모든 관리자 POST, `/v1/test/*`, `/healthz`와 미허용 경로는 relay로 전달하지 않는다. origin은 계속 loopback이며 공개 Postgres 포트는 없다. 이 과제에서 trial API 코드·설정 계약은 그대로지만 현재 운영 trial 앱은 부재다. 비활성 trial을 공개 경로로 복구하지 않는다. 나중에 trial을 재활성화하면 기존 독립 인증을 검증한 구체적 경로만 404 앞에 추가한다. shared host가 같은 config에 있으면 전체 후보로 덮지 말고 KnowsLink 호스트에만 member rule과 명시적 hostname 404 fallback을 병합한다. 다른 호스트 규칙과 마지막 catch-all을 보존한다.
4. 실제 cloudflared `ingress validate`와 `ingress rule <URL>`로 `/`, `/auth/google/callback`, `/home`, `/connect/<43자 ID>`, `/v1/connect/start`, `/v1/text/send`, `/v1/keys/agent/key`, `/v1/receipts/<UUID>`의 relay 선택을 확인한다. `/owner`, `/v1/owners`, 정확한 `/v1/keys`, `/v1/authorize`, `/v1/test/pull`, `/healthz`, 미등록 경로는404 선택을 확인한다. allowlist rule의 `required:false`만 바꿔서는 edge Access가 해제되지 않는다.
5. 먼저 검증된 ingress를 적용해 owner/admin/test의404 차단을 확인한다. 그 뒤 백업한 두 KnowsLink Access 앱만 제거한다. 정책·IdP·다른 앱은 그대로 둔다. 광범위 Everyone 허용·새 Bypass 앱·Access 청구 동의를 만들지 않는다. API가 권한·요금제 제한으로 거부하면 추가 결제 없이 blocker로 기록한다. 경로 제한 없이 Access만 먼저 제거하지 않는다. 기존 `beta.sh expose`는 owner-only bootstrap용이므로 이 전환에 사용하지 않는다. 기존 DNS는 유지하고 검증된 config로 해당 cloudflared 서비스만 재기동한다.
6. 외부 익명 root/connect에서 Cloudflare 로그인 redirect가 사라졌는지 확인한다. 세션 없는 home은 서비스 로그인을 요구하고 잘못된 token/proof는401/422 JSON이어야 한다. owner/admin/test/미허용 경로는404여야 한다. 실제 Google 로그인·동의·로컬 저장, 두 클라이언트의 별도 키·관계 수락·명시 승인 text 왕복을 coor·사용자가 확인한다. 다른 공유 호스트의 기존 응답도 확인한다. 합성 검사는 이 수락을 대신하지 않는다.

공식 근거는 [Tunnel 첫 일치 및 catch-all](https://developers.cloudflare.com/cloudflare-one/networks/connectors/cloudflare-tunnel/do-more-with-tunnels/local-management/configuration-file/)과 [Access application 경로](https://developers.cloudflare.com/cloudflare-one/access-controls/policies/app-paths/)다. Context7 Cloudflare 조회는 quota 초과여서 공식 문서로 대체했다. UI는 변경하지 않았으므로 직접 시각 수락은 기존 Google/동의 흐름에 한정된다.

복구는 먼저 cloudflared를 중지해 외부 노출을 닫는다. 기존 config와 앱 백업을 복원하고 실제 앱 응답의 새/기존 ID·aud와 origin 검사를 대조한 뒤 재노출한다. 삭제한 앱을 재생성하면 aud가 달라질 수 있으므로 과거 aud를 그대로 신뢰하지 않는다. Access 복구가 막히면 Tunnel을 중지한 채 유지하거나 member allowlist+404 상태로만 복구한다. owner/admin을 무인증 fallback으로 열지 않는다. 코드 복귀는 `beta.sh deploy <직전 수락 코드>`의 backup/migration diff 검사를 사용한다. 과거 DB를 덮지 않고 기존 키·철회 상태를 보존한다.

익명 시작의 2000개 상한은 만료 뒤24h tombstone을 포함한다. 공격자가 새 Device 연결을 포화시킬 수 있는 medium 한계는 남는다. 상태 상한과 replay 보존을 유지하며 조기 삭제·IP 정책 변경은 이 최소 과제에서 하지 않는다. 기존 키의 메시지와 기존 회원 로그인·철회는 이 cap에 묶이지 않는다. 신규 연결이 막히면 요청 rate·capacity와 보존 기간을 확인하고 한도 상향이나 DB 삭제로 우회하지 않는다.
