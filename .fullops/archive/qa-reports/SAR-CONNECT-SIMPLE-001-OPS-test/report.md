---
title: main bf67dfcb 운영 배포(연결 단순화·자동 설정 0.1.0 공개) 결과
status: draft
updated: 2026-10-27
owner: ops
tasks: [SAR-CONNECT-SIMPLE-001-OPS]
summary: 연결 단순화 제품 bf67dfcb를 운영 relay에 배포하고 자동 설정 0.1.0 설치 파일과 /install ingress를 공개한 백업·격리 복원·전환·공개 검사·보존 결과를 기록한다.
---

# SAR-CONNECT-SIMPLE-001-OPS 결과

제품 `bf67dfcbb4f7164149901d847edc0aa0e5783dc7`을 운영 relay에 배포하고 자동 설정 0.1.0(설치기·Claude·Codex ZIP)을 `https://link.knowslog.com/install/0.1.0/`에 공개했다. relay 중단은 약 18초이고 healthy·restart 0·migrate 0이다. 설치 파일 3개는 공개 주소에서 200이며 로컬 산출물과 SHA-256이 같다. 실패·rollback은 없다.

## 범위·기준

- 사용자 결정 D042(배포 진행)·D043(자동 설정 0.1.0 함께 공개). Task `task_60054516e2a4`, Dispatch `ctx_626d8f3c3217`. 지시서 기준 `bf67dfcb`, 브랜치 `fullops/ops`. 하위 에이전트는 쓰지 않았다.
- 바꾼 것: 서버 checkout `4fdd149b`→`bf67dfcb`, relay 이미지 재빌드, 새 백업 1개, 설치 폴더 `knowslink-state/install/0.1.0/`(파일 3개), 서버 로컬 compose override `knowslink-state/compose.install.yaml`, 공개 Tunnel ingress `path` 한 대안(`/install/...`), 서버 `deploy-history.log` 한 줄.
- 바꾸지 않은 것: 운영 env 파일(`.env`는 해시·모드가 전 단계 같다), DB schema·migration, 다른 서비스·볼륨, 기존 백업, Cloud Console. 삭제·prune·`down`은 없다.
- 시각은 서버 UTC다(KST = UTC+9). 서버 UTC 2026-10-27T13:27은 KST 22:27이다.
- `beta.sh deploy/backup`과 `verify.py regression/public`은 실행하지 않았다(직전 OPS와 같은 이유: 실패 시 삭제를 포함). `verify.py local`은 baseline, 전환 직후, 최종에 통과했다.

## 지시서와 달라진 점 (coor 승인)

지시서는 "운영 env에 `KNOWSLINK_INSTALL_DIR`만 추가"라고 했다. 그러나 `bf67dfcb`의 `compose.yaml`은 relay 컨테이너에 이 변수를 넘기지 않고 호스트 폴더를 마운트하지도 않는다(relay는 read-only 컨테이너). env 파일에 한 줄만 더하면 relay가 읽지 못해 `/install`이 404로 남는다. coor에게 질문했고 **옵션 A 승인**을 받았다.

- 서버 로컬 override `knowslink-state/compose.install.yaml`(repo 밖, 0644, 비밀 없음)이 relay에 `KNOWSLINK_INSTALL_DIR=/install`과 호스트 `knowslink-state/install` → `/install` **읽기 전용** 바인드 마운트만 더한다.
- 반영 전후 `docker compose config` diff는 정확히 그 env 한 줄과 마운트 한 항목뿐이다([stage.log](stage.log)). 운영 `.env`는 바꾸지 않았다.
- 후속: dev가 이 두 줄을 `compose.ops.yaml`로 옮겨 repo가 추적하게 한다. 그 전에는 모든 compose 호출에 같은 `-f`를 써야 한다(D12 절).

## 단계

| 단계 | 시각(UTC) | 결과 |
|---|---|---|
| baseline | 2026-10-27T13:23Z | 서버 HEAD `4fdd149b`·clean, relay healthy·restart 0, live agent 30(≤200), 백업 7개 |
| build·rollback 태그 | 13:24~13:26Z | diff 확인, 태그 고정, `build relay migrate` EXIT0 |
| 설치 파일 스테이징 | 13:26Z | `install/0.1.0/` 3파일 SHA-256 일치, override 검증 |
| relay 전환 | STOP 13:27:25Z → UP_DONE 13:27:43Z | 약 18초. migrate EXIT0 |
| ingress 반영 | 13:28:19Z~13:28:20Z | cloudflared 한 번 재시작, 연결 4개 등록 |
| 공개 검사 | 13:28Z 이후 | 21/21 통과 |
| 최종 관측 | 13:29:36Z | relay healthy·restart 0, 보존 39/39 |
| 늦은 관측([late.log](late.log)) | 13:32:56Z | relay·Postgres healthy·restart 0, relay 로그 1줄(error/panic 0, capacity warning 0), 재시작 뒤 Tunnel `ERR` 0 |

## 배포 입력 ([build.log](build.log), [baseline.log](baseline.log))

- 서버가 쓴 bundle SHA256 `e7a0ef307425588658e0e7afd89ca8a340d0eef68fe50e326e651ca9951b4224`(3698702bytes). 서버에서 hash 대조·`bundle verify`·`bf67dfcb` 존재·`4fdd149b`→`bf67dfcb` 조상 관계를 확인했다. bundle 끝은 origin/main이고 체크아웃은 고정 SHA `bf67dfcb`로만 옮겼다.
- `db Dockerfile go.mod go.sum Makefile sqlc.yaml compose.yaml` diff는 0이다. `deploy/`의 변경 파일은 `public-ingress.yml` 하나다. 제품 diff는 59파일 +3549/−336(`.fullops-squad` 제외)이다. migration·새 의존성은 없다.
- agent 용량 확인(L-03): live agent(철회 안 됨, owner 활성) 30개, owner 24개. 경고 1000·거절 2000은 현재의 33배·66배다.
- build 전에 현재 이미지 `sha256:95852e24153d…a6be25a0e`에 `knowslink-rollback:4fdd149b-connect001-20261027`을 붙였다. inspect ID가 같고 새 이미지 `sha256:ff7b47e0a9970f…cd7d630`(manifest list)와 다르다. 실제 rollback은 실행하지 않았다.

## 설치 산출물 0.1.0

- 만든 곳: `git archive bf67dfcb`로 푼 깨끗한 사본(Windows, Node v22.22.2, `npm ci` → `npm run build` → `scripts/package_plugin.py --release`, 운영 relay URL 기본값 `https://link.knowslog.com`). 서버 install 폴더에는 파일만 올렸다.
- `knowslink-setup.ps1` SHA256 `33b3ba372340edeef216c8f6fd4c6be548244faf4b782bd6cadde1b674fe16fe`(34875bytes, LF, BOM 없음).
- `knowslink-claude-plugin.zip` SHA256 `767c0e87ca6d7c151f7eccf08914137ec6518cbf00293a0fa1015ced09ce403e`(336625bytes).
- `knowslink-codex-plugin.zip` SHA256 `498644fe3110d357bfcfc41f68a6b1fcf62f0110a38198a50b263bcea3711b3c`(327763bytes).
- 서버: `/home/shin/deploy/knowslink-state/install/0.1.0/`, 폴더 0755·파일 0644(컨테이너 사용자 65532가 읽는다). 게시한 버전은 바꾸지 않는다(패키저가 기존 버전 폴더를 거부한다). 설치기 안의 ZIP SHA-256은 위 ZIP과 같은 빌드에서 채워졌다.

## 백업·격리 복원 ([deploy.log](deploy.log))

| 항목 | 실제 결과 |
|---|---|
| writer 고정 | relay stop 뒤 epoch `2180035`이 2초 간격으로 같았다. 다른 DB client 0 |
| 새 백업 | `backups/4fdd149b-pre-connect001-20261027T132730Z.dump`. 0600·11464bytes·TOC 25줄 |
| 백업 SHA256 | `c6662ff7532b7c652d49d777486e458dd0100cb5cce56e42fb1d5b743df2ea0d` |
| 복원 | network none·host port 없음·postgres 사용자로 복원. 전체 `data::text` md5가 frozen 값과 같다 |
| 복원 검증 | goose2·guard1·rows1·HumanModel1, 격리 DB에서 HumanModel 제거 UPDATE 거부 |
| 재개 | `up -d --no-build --wait relay`. migrate EXIT0·"SQL migrations completed"·goose2 유지 |

복원 fixture `knowslink-connect001-rehearse-20261027`은 exited0·network none으로 보존했다. 운영 DB를 과거 dump로 덮지 않았다. 새 백업은 2026-11-10 이후 14일 자동 정리 대상이 된다(가장 최근 1개는 남는다).

## 공개 ingress ([ingress.log](ingress.log))

- 반영 전 `check_public_ingress.py`(서버 checkout `bf67dfcb`) 통과. live `path` 규칙이 `4fdd149b`의 repo 규칙과 같음을 확인했고, 새 규칙이 그 규칙에 `|/install/[0-9]+\.[0-9]+\.[0-9]+/knowslink-[a-z0-9-]+\.(ps1|zip)` 한 대안만 더한 것인지 단언했다. 바뀐 줄은 2줄(`<`,`>`)이다. `cloudflared tunnel ingress validate` OK, 후보 규칙으로 경로별 매칭을 확인했다.
- 원본 `knowslink-state/connect001-ops-20261027/config.yml.pre-connect001`(0600)을 보관했다. 13:28:19Z에 cloudflared를 한 번 재시작했고 약 1초 끊겼다. 연결 4개가 등록됐다. `TUNNEL_ERR_LINES=5`는 재시작 직전 옛 프로세스의 `Connection terminated` ×4와 `no more connections active and exiting` ×1이다([cferr.log](cferr.log)). 이전 ingress 변경과 같은 종료 로그다.
- 반영 뒤 `check_public_ingress.py` 통과, live 규칙 == repo 규칙(`bf67dfcb`) `True`.

## 공개 검사 ([probe-before.json](probe-before.json), [probe-after.json](probe-after.json), [probe.py](probe.py))

읽기 전용 요청을 1.5초 이상 간격으로 보냈다. 429는 없었다. 실계정 로그인·발송, 계정 상태를 바꾸는 요청은 없다. 응답 본문과 cookie 값은 보존하지 않았다.

- 배포 전(`before`): `/install/0.1.0/*` 3개 404, `/v1/connect/start`에 `host`를 보내면 `invalid_schema`(422, 이전 relay가 알 수 없는 필드로 거절).
- 배포 뒤(`after`) 21/21 통과:
  - `/install/0.1.0/knowslink-setup.ps1` 200 `text/plain; charset=utf-8`, ZIP 2개 200 `application/zip`. 세 본문 SHA-256이 로컬 산출물과 같다.
  - 이웃 경로 9개는 404: 다른 이름 `other.txt`·`knowslink-evil.zip`, 다른 버전 `0.1.1`, 짧은 버전 `0.1`, `/install/`, `/install/0.1.0/`, 하위 경로, 대문자 `.PS1`, 인코딩된 `..%2f`.
  - `/` 200, `/privacy`·`/terms` 200, `/home/agents` 비로그인 303 `/?n=expired`.
  - `/v1/connect/start`에 올바른 형태의 토큰·`host`를 보내면 `invalid_signature`다. `host` 필드를 받아들였고(422 `invalid_schema` 아님) 가짜 proof만 거절한다. 아무것도 저장하지 않는다.
  - `/owner`·`/v1/owners`·`/healthz`는 계속 404.
- 서버 로컬: relay 컨테이너 안 `KNOWSLINK_INSTALL_DIR=/install`, 마운트 `bind /install rw=false`, rootfs read-only. `127.0.0.1:8080/install/0.1.0/*` 3개 200과 같은 content-type, 다른 이름·버전 404.
- 첫 기동 로그는 1줄(`relay starting; … google login configured=true`)뿐이다. error/panic 0, `agent capacity warning` 0([final.log](final.log)).

## 보존 ([preservation.json](preservation.json), [preservation.py](preservation.py))

baseline·before·frozen·after·final 5단계 비교가 39/39 통과했다. 회원2·agent30·owner24·친구1·pair13·메시지/멱등 0의 map digest, Remote(Clients·Attempts·Codes·Installs·Subscriptions·Tokens) digest, 구독 메타, schema SHA, goose2·guard1·HumanModel1, `.env`·credential·tunnel uuid 해시, 기존 백업 7개의 해시/크기/모드, relay·migrate·KnowsLink cloudflared 외 컨테이너(공유 서비스 포함) ID·시작 시각이 같다. 새 컨테이너는 보존한 복원 fixture 하나다.

- **Rates**는 요청 시각 기록이라 달라졌다(volatile). 권한 변경으로 해석하지 않는다.
- **Connections**: baseline에는 만료 대기 중인 기기 연결 시도 1건이 있었고 `before`(구 relay 가동 중, 전환 전)에는 0이었다. 전환 구간 안에서는 `before`→`frozen`→`after`→`final`이 모두 같다. 구 relay의 만료 정리로 본다. 이 과제의 어떤 요청도 연결 시도를 만들지 않았다.
- 서버 checkout은 `bf67dfcb`·clean이다(build 단계에서 옮김).
- 이 비교는 짧은 전환 중의 일치다. 이후 정상 만료·cleanup·사용자 요청에 따른 변경을 금지하는 보장은 아니다.

## 한계·미실행

- 미실행: 실계정 로그인·발송, 실제 연결 시험(CS-AC-18), 자동 설정 설치기를 사용자 PC에서 실행하는 실세션 QA(tester 과제), 브라우저 렌더링·시각 확인, 실제 rollback, Cloud Console·Access·DNS 변경, 부하·장시간 검사.
- 설치 산출물은 Windows에서 빌드했다. 같은 소스를 Linux에서 다시 빌드하면 ZIP 해시가 달라질 수 있다(경로·압축 구현). 게시본의 해시는 위에 고정했다.
- `compose.install.yaml`은 repo가 추적하지 않는다. `beta.sh`나 override 없는 compose 호출로 relay를 다시 만들면 install 마운트가 빠져 `/install`이 404가 된다(회원 기능은 영향 없음). D12 절의 명령을 쓴다.
- rollback: `knowslink-rollback:4fdd149b-connect001-20261027` 이미지로 되돌린다. schema·migration 변경이 없어 코드 전용 복귀다. 되돌리면 새 클라이언트의 `host` 필드를 이전 relay가 거절한다(D12 절의 순서와 영향).

## 보존 위치

- 로컬 원본: `D:/workspace/KnowsLink/.devtools/connect-simple-001-ops/`(bundle·스크립트·JSON 관측·릴리스 산출물 `release/0.1.0/`·빌드 사본 `src/`).
- 서버 private 작업물: `/home/shin/deploy/knowslink-state/connect001-ops-20261027/`(0700). 서버 백업: `backups/` 아래 0600.
- Git 증거: 이 폴더의 redacted log와 [commands.txt](commands.txt)(실행 스크립트 전문). credential·token·cookie 값·DB row는 포함하지 않는다.

## 문서 HEAD의 FullOps 검사

완료 보고의 검증 항목에 적는다(인박스 완료 보고와 같다).
