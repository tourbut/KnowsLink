---
title: main c6dfc39a 운영 배포(자동 설정 0.1.1 게시·서버 override 대체) 결과
status: draft
updated: 2026-10-28
owner: ops
tasks: [SAR-AUTOSETUP-001-OPS]
summary: 자동 설정 0.1.1을 서버에 게시하고 relay를 c6dfc39a로 전환한 뒤 서버 로컬 compose.install.yaml을 repo compose.ops.yaml로 대체한 순서·시각·검사·보존 결과를 기록한다.
---

# SAR-AUTOSETUP-001-OPS 결과

제품 `c6dfc39ac4f6dabddc5864ad611437b84ded383e`를 운영 relay에 배포했다. 자동 설정 0.1.1(설치기·Claude·Codex ZIP)을 `https://link.knowslog.com/install/0.1.1/`에 공개했고 0.1.0은 그대로 200이다. 서버 로컬 override `compose.install.yaml`은 `compose.install.yaml.retired`로 이름만 바꿔 보관했고 install env·마운트는 repo `compose.ops.yaml`에서 온다. relay 중단은 약 18초이고 healthy·restart 0·migrate 0이다. 실패·rollback은 없다.

## 범위·기준

- 사용자 결정 D044('운영 교체 진행', 2026-10-28). Task `task_d0ab813e4088`, Dispatch `ctx_ecf85d36b31f`. 지시서 기준 main `c6dfc39a`(브랜치 `fullops/ops`). 하위 에이전트는 쓰지 않았다(직렬 운영).
- 바꾼 것: 서버 checkout `bf67dfcb`→`c6dfc39a`, relay 이미지 재빌드, 새 백업 1개, 설치 폴더 `knowslink-state/install/0.1.1/`(파일 3개), `compose.install.yaml` 개명 보관, 서버 `deploy-history.log` 두 줄.
- 바꾸지 않은 것: 운영 `.env`(해시·모드가 전 단계 같다), 공개 Tunnel `config.yml`·ingress(임의 semver를 이미 허용), DB schema·migration, 다른 서비스·볼륨, 기존 백업, `install/0.1.0/` 파일. 삭제·prune·`down`은 없다.
- 시각은 서버 UTC다(KST = UTC+9). 서버 UTC 2026-10-28T13:19Z는 KST 22:19다.
- `beta.sh deploy/backup`과 `verify.py regression/public`은 실행하지 않았다(실패 시 삭제를 포함하는 직전 OPS와 같은 이유). `verify.py local`은 baseline, 전환 직후, override 대체 직후, 최종에 통과했다.
- 실계정 로그인·발송, 설치기를 사용자 PC에서 실행하는 시험은 하지 않았다.

## 단계

| 단계 | 시각(UTC) | 결과 |
|---|---|---|
| baseline | 2026-10-28T13:15:40Z | 서버 HEAD `bf67dfcb`·clean, relay healthy·restart 0, live agent 32·owner 24(경고 1000·거절 2000보다 한참 아래), 백업 7개, `install/0.1.0/` 3파일 SHA 기록 |
| 0.1.1 게시 | 13:15Z | `install/0.1.1/` 3파일 SHA-256 일치, 0.1.0 SHA 그대로, 로컬 200 3개 |
| 게시 후 공개 검사 | 13:16Z 이후 | 26/26 통과([probe-staged.json](probe-staged.json)) |
| build·rollback 태그 | 13:17~13:18Z | diff 확인(migration 없음), 태그 고정, `build relay migrate` EXIT0 |
| relay 전환 | STOP 13:19:00Z → UP_DONE 13:19:18Z | 약 18초. migrate EXIT0, "SQL migrations completed" |
| 전환 후 공개 검사 | 13:19Z 이후 | 26/26 통과([probe-after.json](probe-after.json)) |
| override 대체 | 13:20:40Z~13:20:43Z | `config` 비교 SAME, 개명, `up -d --no-build --wait relay`. relay 컨테이너 재생성 없음 |
| 최종 관측(전환 5분 뒤) | 13:24:40Z, 13:24:53Z | healthy·restart 0, error/panic 0 ([final.log](final.log)) |
| 최종 공개 검사 | 13:25Z 이후 | 26/26 통과([probe-final.json](probe-final.json)) |

## 제품 diff 확인 ([build.log](build.log))

- `bf67dfcb..c6dfc39a` 중 `db Dockerfile go.mod go.sum sqlc.yaml compose.yaml` diff는 0이다. migration·Go 의존성 변경은 없다.
- `.fullops-squad` 밖 변경 파일은 18개(+119/−42): relay 코드는 `internal/relay/human_pages.go` 한 줄(화면 명령 0.1.0→0.1.1), `deploy/knowslink/compose.ops.yaml`(relay install env·`:ro` 마운트 4줄), 나머지는 어댑터 버전 상향·zod 4.6.5→4.6.7 patch(어댑터 의존성, relay 이미지에는 들어가지 않는다)·설치기 O-04 수정과 그 테스트다.
- 서버가 쓴 incremental bundle(`^bf67dfcb origin/main`) SHA256 `21464122ad3e2d42db84971846efc1d1be9e429117d92cf2bcd204c6bc609020`(261846bytes). 서버에서 hash 대조·`bundle verify`·`bf67dfcb`→`c6dfc39a` 조상 관계를 확인했다. bundle 끝은 origin/main(`a1cd96f4`, 문서 1파일만 `c6dfc39a`보다 앞섬)이고 체크아웃은 고정 SHA `c6dfc39a`로만 옮겼다.
- rollback 태그 `knowslink-rollback:bf67dfcb-autosetup011-20261028`은 이전 relay 이미지 `sha256:ff7b47e0a997…e954d630`다. build 전에 붙였고 inspect ID가 같다. 새 이미지 `sha256:dfc3ffeb259e…7eefa3`(manifest list)와 다르다. 실제 rollback은 실행하지 않았다.

## 설치 산출물 0.1.1

- 만든 곳: `git archive c6dfc39a`로 푼 깨끗한 사본(Windows, Node v22.22.2, npm 10.9.7, `npm ci` → `npm run build` → `scripts/package_plugin.py --release`, 운영 relay URL 기본값 `https://link.knowslog.com`). 서버에는 파일만 올렸다.
- `knowslink-setup.ps1` SHA256 `c4d6825dd5519cbd6166e588123e83d0624ab6a45db6239e1d2e47b1b91a4dfe`(35823bytes, LF, BOM 없음).
- `knowslink-claude-plugin.zip` SHA256 `06b64bc2c2ed7584e711f33c10ccb40dac26607f62dbfe92914f3ebfb48bdb0b`(337538bytes).
- `knowslink-codex-plugin.zip` SHA256 `0eacd83ba3dc2235c3abb7a42882bc10dcc491dd0a71970776e6fbefddfe14ae`(328676bytes).
- 서버: `/home/shin/deploy/knowslink-state/install/0.1.1/`, 폴더 0755·파일 0644. 같은 소스로 dev가 앞서 만든 `o04-release/0.1.1`과 세 파일이 바이트 단위로 같다(빌드가 재현된다).
- 설치기 안의 두 ZIP SHA-256(`$Sums`)은 위 ZIP과 같고 `$Version = '0.1.1'`·`$Base = …/install/0.1.1`이다(probe 단언).
- `install/0.1.0/`(`33b3ba37…16fe`, `767c0e87…403e`, `498644fe…1b3c`)는 baseline·게시 후·override 대체 후 모두 `sha256sum -c` OK다. 폴더는 지우지 않았다.

## 백업·격리 복원 ([deploy.log](deploy.log))

| 항목 | 실제 결과 |
|---|---|
| writer 고정 | relay stop 뒤 epoch `2267807`이 2초 간격으로 같았다. 다른 DB client 0 |
| 새 백업 | `backups/bf67dfcb-pre-autosetup011-20261028T131905Z.dump`. 0600·12571bytes·TOC 25줄 |
| 백업 SHA256 | `000989324dcad1d8937ba91084bab9b443212e4f3cd1564f89df02ee209aa233` |
| 복원 | network none·host port 없음·postgres 사용자로 복원. 전체 `data::text` md5가 frozen 값과 같다 |
| 복원 검증 | goose2·guard1·rows1·HumanModel1, 격리 DB에서 HumanModel 제거 UPDATE 거부 |
| 재개 | `up -d --no-build --wait relay`. migrate EXIT0·goose2 유지 |

복원 fixture `knowslink-autosetup011-rehearse-20261028`은 exited0·network none으로 보존했다. 운영 DB를 과거 dump로 덮지 않았다. 새 백업은 2026-11-11 이후 14일 자동 정리 대상이 된다(가장 최근 1개는 남는다).

## 화면 명령 0.1.1

실계정으로 로그인하지 않아 회원 화면 렌더링은 열지 않았다. 대신 실행 중인 relay 바이너리를 확인했다: `docker exec knowslink-relay-1 grep -a -c … /app/relay`가 `install/0.1.1/knowslink-setup.ps1`에서 1, `install/0.1.0/knowslink-setup.ps1`에서 0이다([deploy.log](deploy.log) 끝). 템플릿 `human-setup-command` 문자열이 0.1.1만 담는다. 화면 렌더링 결과는 로그인 이후라 tester 실세션 QA의 몫이다.

## override 대체 ([override.log](override.log))

- 서버 checkout `c6dfc39a`의 `compose.ops.yaml`에 `KNOWSLINK_INSTALL_DIR`가 있다(1줄). 비밀이 섞인 `config` 전문은 파일에 저장하지 않고 `diff <(…) <(…)`로만 비교했다: override 없는 `config relay`와 override를 얹은 `config relay`가 `SAME_relay`, 전체 서비스 `config`도 `SAME_all_services`다. `KNOWSLINK_INSTALL_DIR: /install`, `target: /install`, `read_only: true`가 override 없는 config에 있다.
- 13:20:40Z `compose.install.yaml`을 `compose.install.yaml.retired`로 개명했다(`mv -n`, 삭제 안 함, 0644·277bytes 그대로). 그다음 override 없이 `up -d --no-build --wait relay`를 실행했다. relay 컨테이너 ID `81452ca51dfc…`와 StartedAt `13:19:12.65Z`가 전후 같다(재생성·끊김 없음). `migrate`만 한 번 다시 돌았고 exit 0이다.
- 이후 relay 마운트: `bind …/knowslink-state/install -> /install rw=false`, 환경 `KNOWSLINK_INSTALL_DIR=/install`, 루트 파일시스템 read-only. 로컬 `127.0.0.1:8080/install/{0.1.0,0.1.1}/*` 6개 모두 200이다.

## 공개 검사 ([probe-staged.json](probe-staged.json), [probe-after.json](probe-after.json), [probe-final.json](probe-final.json), [probe.py](probe.py))

읽기 전용 요청을 1.5초 이상 간격으로 보냈다. 429는 없었다. 실계정 로그인·발송, 계정 상태를 바꾸는 요청은 없다. 응답 본문과 cookie 값은 보존하지 않았다. 세 번(게시 직후·relay 전환 후·최종) 모두 26/26 통과했다.

- `/install/0.1.1/` 세 파일: 200, `text/plain; charset=utf-8`·`application/zip`, 본문이 로컬 산출물과 바이트 단위로 같고 SHA-256이 위와 같다. `/install/0.1.0/` 세 파일도 200이고 로컬 0.1.0 산출물과 같다.
- 설치기 안의 두 ZIP SHA-256·버전·기준 URL이 각 버전의 ZIP·폴더와 맞다.
- 이웃 경로 9개는 404: 다른 이름 `other.txt`·`knowslink-evil.zip`, 없는 버전 `0.1.2`, 짧은 버전 `0.1`, `/install/`, `/install/0.1.1/`, 하위 경로, 대문자 `.PS1`, 인코딩된 `..%2f`.
- `/` 200, `/privacy`·`/terms` 200, `/home/agents` 비로그인 303. `/v1/connect/start`에 올바른 형태의 토큰·`host`·가짜 proof를 보내면 `invalid_signature`다(아무것도 저장하지 않는다). `/owner`·`/v1/owners`·`/healthz`는 404.
- 첫 기동 로그는 1줄(`relay starting; … google login configured=true`)뿐이다. error/panic 0, `agent capacity warning` 0. 전환 뒤 Tunnel `ERR` 0, Postgres `ERROR|FATAL|PANIC` 0([final.log](final.log)).

## 보존 ([preservation.json](preservation.json), [preservation.py](preservation.py))

baseline·before·frozen·after·pre-override·after-override·final·final2 8단계 비교가 41/41 통과했다. 회원2·agent32·owner24·친구1·pair13·메시지1·멱등1·연결 시도3의 map digest, Remote(Clients·Attempts·Codes·Installs·Subscriptions·Tokens) digest, 구독 메타, schema SHA, goose2·guard1·HumanModel1, `.env`·credential·Tunnel config·tunnel uuid 해시, 기존 백업 7개의 해시/크기/모드, relay·migrate 외 컨테이너(Postgres·KnowsLink cloudflared·공유 서비스 포함) ID·시작 시각이 같다. 새 컨테이너는 보존한 복원 fixture 하나다.

- **Rates**는 요청 시각 기록이라 달라졌다(volatile, 단계별로 내 probe와 사용자 요청이 만든다). 권한 변경으로 해석하지 않는다. 전환 중(`before`→`frozen`)은 같다.
- **호스트 cloudflared PID**: `after-override`까지 `2039731`였고 `final`부터 `2354559`다. 이 프로세스는 KnowsLink 컨테이너가 아니라 사용자 서비스 `cloudflared-orca.service`(Orca 터널, `~/.cloudflared/config.yml`)다. 같은 파일의 `config.yml`이 22:21:07 KST(13:21:07Z)에 수정되고 `.bak-20261028`이 생겼으며 서비스가 22:21:38 KST에 다시 시작됐다. 이 과제는 `knowslink-state/`와 KnowsLink compose만 건드렸다. 다른 작업의 변경이며 공유 orca/s8/mcp 공개 상태(`shared`)는 모든 단계에서 같다. 의도하지 않은 영향으로 해석하지 않지만 보고에 남긴다.
- 서버 checkout은 `c6dfc39a`·clean이다.
- 이 비교는 짧은 전환 중의 일치다. 이후 정상 만료·cleanup·사용자 요청에 따른 변경을 금지하는 보장은 아니다.

## 한계·미실행

- 미실행: 실계정 로그인·발송, 회원 화면 렌더링(로그인 필요; 바이너리 문자열로 대체), 설치기를 사용자 PC에서 실행하는 실세션 QA(tester 과제: 플러그인 2개 이상 환경에서 갱신 재실행 AC-07·S5), 브라우저 시각 확인, 실제 rollback, Cloud Console·Access·DNS 변경, 부하·장시간 검사.
- 0.1.0을 설치한 사용자와 다단 플러그인 환경은 0.1.1 명령을 한 번 실행해야 한다(O-05). 옛 0.1.0 명령은 계속 공개되며 다단 환경에서 실패한다. 0.1.0 파일은 삭제하지 않는다.
- 설치 산출물은 Windows에서 빌드했다. 같은 소스를 Linux에서 다시 빌드하면 ZIP 해시가 달라질 수 있다. 게시본의 해시는 위에 고정했다.
- `compose.install.yaml.retired`는 서버에 남아 있다. 서버 checkout을 `bf67dfcb` 이하로 되돌리면 그 파일을 원래 이름으로 돌려 세 번째 `-f`로 얹어야 `/install`이 유지된다(D12).
- rollback: `knowslink-rollback:bf67dfcb-autosetup011-20261028`로 코드만 되돌린다(schema·migration 변경 없음). 되돌리면 회원 화면 명령이 0.1.0으로 돌아간다. `/install/0.1.1/` 폴더는 그대로 둔다.

## 보존 위치

- 로컬 원본: `D:/workspace/KnowsLink/.devtools/autosetup-001-ops/`(bundle·스크립트·관측 JSON `obs-*.json`·릴리스 산출물 `release/0.1.1/`·빌드 사본 `src/`).
- 서버 private 작업물: `/home/shin/deploy/knowslink-state/autosetup011-ops-20261028/`(0700). 서버 백업: `backups/` 아래 0600.
- Git 증거: 이 폴더의 redacted log와 [commands.txt](commands.txt)(실행 스크립트 전문). credential·token·cookie 값·DB row는 포함하지 않는다.

## 문서 HEAD의 FullOps 검사

완료 보고의 검증 항목에 적는다(인박스 완료 보고와 같다).
