---
id: D12
title: 운영자설명서
status: draft
updated: 2026-10-03
owner: ops
tasks: [SAR-DEPLOY-001-OPS]
upstream: [D02, D03]
summary: 서버 관찰 이력과 배포·백업 계획 및 C1 high 배포 차단을 기록한다
---

# KnowsLink 운영자설명서 (D12) — 현재 서버 배포 계획

이 문서는 현재 서버의 Docker와 `link.knowslog.com` Cloudflare Tunnel 배포 계획이다. 상태는 **계획**이다. 이전 OPS Dispatch가 2026-10-03에 수행한 읽기 전용 조사 초안을 보존한다. 이번 retry Dispatch는 서버 조사를 재실행하지 않았다. 원본 명령 로그와 개별 종료코드가 초안에 없으므로 아래 관찰을 새 검증 통과로 해석하지 않는다. DNS·Tunnel·컨테이너를 만들거나 바꾸지 않았다.
배포·이행 증거가 없으므로 D11(사용자설명서)과 D13(인수인계서)은 미작성으로 유지한다. 이 문서도 배포 완료 근거가 아니다.
비밀값·토큰·`cert.pem` 내용은 기록하지 않는다. 존재·권한·경로만 기록한다.

현재 Dispatch `ctx_3c54fe7d2043`는 이전 완료 기록의 검증·커밋만 수행한다. coordinator 메시지 `msg_d89612ea9341`에 따라 `a6a10c7`은 수락·배포 후보가 아니다. 독립 reviewer가 agent credential로 `deliver:human`을 처리하는 C1 high 결함을 보고했다. DEV 수정과 새 고정 SHA의 독립 QA·리뷰·coor 수락 전에는 배포하지 않는다. 이전 QA 통과와 이번 D12 기록 완료는 이 결함의 해결이나 제품 수락을 의미하지 않는다.

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

## 6. 적용 절차 (사용자 재개 지시·선행 조건 충족·후속 Dispatch 뒤)

현재는 모든 단계가 미실행이다. 각 단계의 종료코드를 보존한다. 실패하면 중단하고 보고한다.
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
| 이전 서버 관찰 보존·운영 계획(D12 초안) | 기록 완료. 새 서버 검증은 미실행 |
| DNS·Tunnel·컨테이너 변경 | **미실행** |
| 수락 SHA 배포·공개 health 검증 | 미실행·held. 기존 a6a10c7은 C1 high로 수락·배포 금지. 사용자 재개 지시와 후속 Dispatch가 필요 |
| D11 사용자설명서·D13 인수인계서 | 미작성 |

## 10. 개정 근거와 검증 범위

2026-10-03 retry에서 이전 미추적 초안을 보존하여 수정했다. 지시서 기준 ref는 `0dd08ec994771836c15d9d22a6a83393a71d7987`이며 공통 규칙은 `fullops-common-0.3.2`다.
OPS 체크아웃의 제품 코드는 초기 골격이다. 기존 후보의 제품 정본은 coordinator 체크아웃 `/home/shin/orca/workspaces/KnowsLink/fullops-coor`에서 읽었다. 대상은 README.md, .fullops-squad/project.md, docs/design-docs/architecture.md·interface-design.md·database-design.md, docs/exec-plans/phases/SAR-MVP-001-DEV.md, docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md다. 문서 경로는 `.fullops-squad/` 기준이다. 제품 SHA는 a6a10c7이고 QA 기록 SHA는 c59537b다. 원천 파일은 수정하지 않았다.

이 단계는 D12 계획 문서 검사만 수행한다. 실제 서버·DNS·Tunnel·복원·health/auth 검증을 재실행하지 않았다. 원래 관찰에는 새 PASS나 종료코드를 붙이지 않았다. D11/D13은 미작성이다. 현재 사용자 지시에 따라 이번 기록 완료 뒤 운영·후속 기능을 시작하지 않는다.

관련 원천: [D02](../planning/product-specs/SAR-MVP.md), [기능·held](../planning/SAR-MVP-backlog.md), [실행 기록](../exec-plans/phases/SAR-DEPLOY-001-OPS.md), [산출물 인덱스](../deliverables/README.md).
