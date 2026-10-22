---
title: main 4fdd149b 운영 배포(개인정보처리방침 14일 백업 문장) 결과
status: draft
updated: 2026-10-22
owner: ops
tasks: [SAR-PUBLISH-001-P03-DEPLOY]
summary: 개인정보처리방침 KO/EN의 14일 백업 문장 제품 4fdd149b를 운영 relay에 배포한 백업·격리 복원·전환·공개 검사·보존 결과를 기록한다.
---

# SAR-PUBLISH-001-P03-DEPLOY 결과

제품 `4fdd149b27d3ff251a73128120dce9f8bb136620`을 운영 relay에 배포했다. relay 중단은 약17초이고 healthy·restart 0이다. live `/privacy`(한국어)와 영어판에 ‘14일’ 문장이 보이고 옛 ‘교체’ 문장은 없다. `/terms`·`/`는 200이다. 실패·rollback은 없다.

## 범위·기준

- 사용자 결정 D020·D021. Task `task_b05e4e6b2ff4`, Dispatch `ctx_28970340f119`. 준비 HEAD는 브랜치 `fullops/ops`의 `5cbdb24e`이고 제품 파일은 main `4fdd149b`와 같다. 하위 에이전트는 쓰지 않았다.
- 바꾼 것: 서버 checkout `7c6ef373`→`4fdd149b`, relay 이미지(`policy/*.html` 포함 재빌드), 새 백업 1개, 서버 `deploy-history.log` 한 줄.
- 바꾸지 않은 것: 운영 env, Tunnel ingress(cloudflared 재시작 없음), DB schema, 컨테이너(relay·migrate 제외), Cloud Console, 기존 백업. 삭제·prune·`down`은 없다.
- 시각은 서버 UTC다(KST = UTC+9). 서버 UTC 2026-10-21T23:57은 KST 2026-10-22 08:57이다.
- `beta.sh deploy/backup`과 `verify.py regression/public`은 실행하지 않았다(직전 OPS와 같은 이유). `verify.py local`은 baseline, 전환 직후, 최종에 통과했다.

| 단계 | 시각 | 결과 |
|---|---|---|
| baseline | 2026-10-21T23:55Z | 서버 HEAD `7c6ef373`·clean, relay healthy, 백업 11개, prune cron 1줄 |
| relay 전환 | STOP 23:57:12Z → UP 23:57:29Z | 약17초. migrate EXIT0 |
| 공개 검사 | 23:57:43Z, 재실행 23:58Z | 7/7 통과 |
| 최종 관측 | 23:58:15Z | relay running·healthy restart 0, cloudflared restart 0 |

## 배포 입력 ([build.log](build.log))

- bundle SHA256 `e14f3f175bd13cf4a322ed2a6afd72dd3a7541e77fc1c33180c0a40be8800798`(113207bytes). 서버에서 hash 대조·`bundle verify`·`4fdd149b` 존재·`7c6ef373`→`4fdd149b` 조상 관계를 확인했다. bundle 끝은 origin/main `5cbdb24e`이고 체크아웃은 고정 SHA `4fdd149b`로만 옮겼다.
- `db Dockerfile go.mod go.sum Makefile sqlc.yaml compose.yaml` diff는 0이다. 제품 diff는 3파일 +25/−2(`privacy.ko.html`·`privacy.en.html` 각 1줄, `deploy/knowslink/prune-backups.sh` 신규 23줄)다. `prune-backups.sh`는 RETENTION 과제가 이미 서버 `~/.local/bin`에 설치했고 relay 이미지와 무관하다.
- build 전에 현재 이미지 `sha256:69f62b72483b…e95a2`에 `knowslink-rollback:7c6ef373-p03-20261022`를 붙였다. inspect ID가 같고 새 이미지 `sha256:95852e24153d…a6be25a0e`(manifest list)와 다르다. 실제 rollback은 실행하지 않았다.

## 백업·격리 복원 ([deploy.log](deploy.log))

| 항목 | 실제 결과 |
|---|---|
| writer 고정 | relay stop 뒤 epoch `1697089`이 2초 간격으로 같았다. 다른 DB client 0 |
| 새 백업 | `backups/7c6ef373-pre-p03deploy-20261021T235716Z.dump`. 0600·11226bytes·TOC 25줄 |
| 백업 SHA256 | `8eed17521aaac56dbe54fb3a412112455a9fab75e94f956f1ea946cbd87b6ca5` |
| 복원 | network none·host port 없음·postgres 사용자로 복원. 전체 `data::text` md5가 frozen 값과 같다 |
| 복원 검증 | goose2·guard1·rows1·HumanModel1, 격리 DB에서 HumanModel 제거 UPDATE 거부 |
| 재개 | `up -d --no-build --wait relay`. migrate EXIT0·"SQL migrations completed"·goose2 유지 |

복원 fixture `knowslink-p03deploy-rehearse-20261022`는 exited0·network none으로 보존했다. 운영 DB를 과거 dump로 덮지 않았다. 새 relay 로그는 기동 1줄뿐이고 error/panic 0이다.

## 공개 검사 ([smoke.json](smoke.json), [smoke.py](smoke.py))

`all_pass`. 읽기 전용 GET 7회(두 번 실행)를 1.5초 이상 간격으로 했고 429는 없었다. 로그인·송신은 없다. 응답 본문과 cookie 값은 보존하지 않았다.

- `GET /privacy`(기본 한국어): 200, ‘백업은 만든 날부터 14일이 지나면 삭제합니다.’ 있음, 옛 ‘정기적으로 교체’ 없음.
- `GET /lang/en?back=/privacy`: 303 `Location: /privacy?lang=en`·언어 cookie 설정. 그 cookie로 `GET /privacy`: 200, ‘We delete backups 14 days after they are made.’ 있음, 옛 ‘replaced regularly’ 없음.
- `GET /terms` 200, `GET /` 200. 이웃 `/privacy/`·`/privacyx`는 404다.
- 첫 시도에서 `/lang/en` 리다이렉트 기대값을 `/privacy`로 잘못 두어 1건이 실패로 찍혔다. 제품 결함이 아니라 검사 기대값 오류였다(실제 `/privacy?lang=en`). 기대값을 고쳐 전체를 다시 실행해 7/7을 얻었다. 저장된 `smoke.json`은 재실행본이다.

## 보존 (로컬 비교, baseline·before·frozen·after·final)

- 회원·agent·owner·pair·친구·메시지·세션·신원 등 상태 map count와 digest가 5단계에서 같다(Rates 제외). Remote의 Clients·Attempts·Codes·Tokens·Installs·Subscriptions digest, 구독 메타, Inbox 수가 같다. 이번엔 Remote.Tokens도 같았다.
- **Rates**는 요청 시각 기록이라 달라졌다(volatile). 권한 변경으로 해석하지 않는다.
- schema SHA256·`.env`·`tunnel.uuid`·`access.aud`·`tunnel/config.yml`·Tunnel credential의 hash/mode가 같다. 기존 백업 11개의 hash/크기/mode가 같고 새 백업 1개만 늘었다(12).
- relay·migrate 외 컨테이너(공유 서비스, cloudflared 포함)의 ID·시작 시각이 같다. 호스트 cloudflared PID, orca/s8/mcp 응답이 같다. loopback 8080 listener는 relay 중지 구간(frozen)에서만 없었다(예상).
- 이 비교는 짧은 전환 중의 일치다. 이후 정상 만료·cleanup·사용자 요청에 따른 변경을 금지하는 보장은 아니다.

## 한계·미실행

- 미실행: 브라우저 렌더링·시각 확인, 실계정 로그인·발송, 실제 rollback, Cloud Console·Access·DNS 변경, 부하·장시간 검사, 다른 브라우저. 영어판은 cookie를 단 HTTP GET 문자열로 확인했다.
- `policyEffective`(`2026-10-22`)와 정책 문안의 법률 검토(P-05)는 이번 범위가 아니다.
- 새 백업은 2026-11-05 이후 14일 정리 대상이 된다(가장 최근 1개는 남는다). 더 오래 둘 복구점은 다른 곳에 복사해야 한다. cron 첫 자동 실행(2026-10-22 18:17Z)은 이 과제 시점에 아직 일어나지 않았다.
- rollback: `knowslink-rollback:7c6ef373-p03-20261022` 이미지로 되돌린다. schema·migration 변경이 없어 코드 전용 복귀다.

## 보존 위치

- 로컬 원본: `C:/Users/shin/AppData/Local/KnowsLinkDevTools/p03deploy-ops/`(bundle·스크립트·JSON 관측).
- 서버 private 작업물: `/home/shin/deploy/knowslink-state/p03deploy-20261022/`. 서버 백업: 위 경로 `backups/` 아래 0600.
- Git 증거: 이 폴더의 redacted log와 [commands.txt](commands.txt)(실행 스크립트 전문). credential·token·cookie 값·DB row는 포함하지 않는다.

## 문서 HEAD의 FullOps 검사

완료 보고의 검증 항목에 적는다(인박스 완료 보고와 같다).
