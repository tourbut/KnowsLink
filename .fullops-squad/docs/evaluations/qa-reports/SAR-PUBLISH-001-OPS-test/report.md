---
title: main 7c6ef373 운영 배포·캘린더 키 제거·정책 페이지 공개 결과
status: draft
updated: 2026-10-22
owner: ops
tasks: [SAR-PUBLISH-001-OPS]
summary: 캘린더 제거·정책 페이지 제품 7c6ef373을 운영 relay에 배포하고 운영 env의 캘린더 키를 지운 뒤 live ingress에 /privacy·/terms를 더한 백업·정리 확인·보존·공개 검사 결과를 기록한다.
---

# SAR-PUBLISH-001-OPS 결과

제품 `7c6ef3734a38a8f273714aa2915dc2f14f83a7e4`를 기존 link.knowslog.com 운영 relay에 배포했다. 저장된 캘린더 연결은 처음부터 0개라 정리 줄은 나오지 않았고, 운영 env에서 `KNOWSLINK_CALENDAR_KEY`를 지운 뒤 relay를 다시 띄웠다. live 공개 ingress에 `/privacy`·`/terms`만 더했고 두 경로는 비로그인 200이다. 실패·rollback은 없었다.

## 범위·기준·시간

- 사용자 결정 D020 '배포 진행'(2026-10-22)이다. 준비 HEAD는 브랜치 `fullops/ops`의 `82c4b81c`이고 제품 파일은 main `7c6ef373`과 같다.
- Task `task_1ecf590bc634`, Dispatch `ctx_b642becdd74d`, worker `term_0abbfbe0-12e7-4e20-8489-ada1889321f2`, 복귀 coor `term_24c75449-b954-447a-9f01-d385ce1999f0`이다.
- 하위 에이전트는 쓰지 않았다. 직렬 운영 전환은 직접 수행했다.
- `beta.sh deploy/backup/restore-verify`와 `verify.py regression/public`은 실행하지 않았다(직전 배포와 같은 이유). `verify.py local`은 baseline, relay 전환 직후, 키 제거 재기동 직후, 최종에 통과했다.
- 시각은 서버 UTC다(KST = UTC+9).

| 단계 | 시각 | 결과 |
|---|---|---|
| baseline | 2026-10-21T22:39:35Z | 서버 HEAD `4787e85`·clean, relay healthy, 캘린더 연결 0 |
| relay 전환 | STOP 22:45:36Z → UP 22:45:53Z (07:45:36~07:45:53 KST) | 약17초. migrate EXIT0 |
| 키 제거 재기동 | 22:46:32Z → 22:46:40Z | relay만 재생성, healthy |
| ingress 반영 | cloudflared 재시작 22:47:14Z | 연결 4개 재등록 |
| 최종 관측 | 22:50:33Z | relay·cloudflared running, restart 0 |

## 배포 전 상태 ([baseline.log](baseline.log), [baseline.json](baseline.json))

- 서버 HEAD `4787e85`·clean. relay healthy, `healthz` 로컬 200. 연결 수: 회원 2·agent 28·pair 13·친구 1·Remote 설치 2. **캘린더 연결 0**(상태 JSON에 `Calendars` 키 없음). relay 로그의 캘린더 줄 0. 운영 env 모드 0600·518bytes, 캘린더 키 줄 1개.
- 컨테이너 124·백업 파일 19.

## 배포 입력

- bundle SHA256 `ce1866100ce179832270699976e302197e744df5727256b7a219454408731ed0`(3319072bytes)다. 서버에서 hash 대조·`bundle verify`·`FETCH_HEAD`=`7c6ef373`·`4787e85`→`7c6ef373` 조상 관계를 확인했다([build.log](build.log)).
- `db Dockerfile go.mod go.sum Makefile sqlc.yaml`의 diff는 0이다. migration과 새 의존성은 없다. `compose.yaml`은 `KNOWSLINK_CALENDAR_KEY` 줄의 주석만 바뀌었다. `deploy scripts` diff는 `public-ingress.yml`·`check_public_ingress.py`뿐이다. 제품 diff는 51파일 +950/−2095다.
- `policyEffective`는 `2026-10-22`이고 배포일(KST 2026-10-22)과 같아 바꾸지 않았다.
- build 전에 현재 이미지(`sha256:0e0eb8b44e47…ae7fd49`)에 `knowslink-rollback:4787e85-publish001-20261022`를 붙였다. inspect ID가 같다. 기존 rollback 태그는 덮지 않았다. network none·read-only 컨테이너 `knowslink-publish001-rollback-inspect-20261022`에서 `/app/relay`·`/app/migrate` 실행 파일과 SHA256을 확인했다([rollback-verify.log](rollback-verify.log)). 앱 기동·실제 rollback은 하지 않았다.
- 새 relay 이미지는 manifest list `sha256:69f62b72483b…e95a2`다.

## 백업·격리 복원·재개 ([deploy.log](deploy.log))

| 항목 | 실제 결과 |
|---|---|
| writer 고정 | relay stop 뒤 epoch `1692641`이 2초 간격으로 같았다. 다른 DB client 0 |
| 새 백업 | `backups/4787e85-pre-publish001-20261021T224540Z.dump`. 파일 0600·11178bytes |
| 백업 SHA256 | `9d21d9ef0e8459ea9934ea3ce5feb3fbda16330c8c4c6991c7613c0050aed848` |
| list·복원 | `pg_restore --list` TOC 25줄. network none·host port 없음·postgres 사용자로 복원 |
| 복원 검증 | 복원 `data::text` md5가 frozen 운영 값과 같다. goose2·guard1·rows1·HumanModel1, 격리 DB에서 HumanModel 제거 UPDATE 거부 |
| 재개 | 미리 만든 이미지로 `up -d --no-build --wait relay`. migrate EXIT0·"SQL migrations completed"·goose2 유지 |

복원 fixture `knowslink-publish001-rehearse-20261022`는 exited0·network none·port 없음으로 보존했다. 운영 DB를 과거 dump로 덮지 않았다. Down·volume 삭제·prune·이미지 삭제는 없다. 서버 `deploy-history.log`에 두 줄을 더했다.

## 첫 기동 캘린더 정리 판정

- 첫 기동 relay 로그는 1줄(기동 줄)뿐이고 `calendar cleanup` 줄은 **0**이다. 키는 그대로 둔 채(`RELAY_ENV_CALENDAR_KEY_TEXT_LEN=44`) 기동했다.
- 리뷰 low(commit과 revoke 사이에서 멈추면 줄이 없을 수 있음)에 따라 줄이 아니라 상태로 판정했다. baseline과 키 제거 직전·직후의 상태 JSON 모두 `Calendars` 키가 없고 연결 수가 0이다([envkey.log](envkey.log), [baseline.log](baseline.log)). 지울 연결이 처음부터 없었으므로 정리 줄이 없는 것이 맞다. 중간에 멈춘 흔적이 아니다.
- 따라서 `connections=0 revoked=0 failed=0`에 해당한다. Google revoke 호출은 일어나지 않았다. `failed`가 0이므로 Google 계정에서 직접 접근을 지워 달라고 안내할 사용자는 없다.

## 운영 env 키 제거 ([envkey.log](envkey.log))

값은 출력·기록하지 않았고 줄 수·길이·모드만 확인했다.

- 이전 env를 `publish001-ops-20261022/env.pre-publish001`(0600·518bytes, SHA256 `131e0798…50bb8`)에 보관했다. 키 줄 1개를 지운 새 `.env`는 0600·450bytes, 키 줄 0개다. 보관본과 달라진 줄은 1개(키 줄)뿐이다. 이 450bytes는 키를 추가하기 전 RELEASE-003의 `.env` 크기와 같다.
- 기록: `docker compose … up -d --no-build --wait relay`가 env 변경으로 relay만 재생성했다. 재기동 relay의 `KNOWSLINK_CALENDAR_KEY` 길이는 0이다.
- idempotent 확인: 재기동 로그는 1줄(기동 줄)이고 `calendar cleanup` 줄이 0이다. error/panic 0. 상태 JSON의 캘린더 연결은 0 그대로다. 키가 없는데도 기동 정리가 아무 일도 하지 않고 기동에 실패하지 않았다.
- 서버 `compose.yaml`의 전달 줄은 선택 값이라 그대로 뒀다.

## live ingress ([ingress-apply.log](ingress-apply.log))

- 승인 범위(D020) 안에서 `/privacy|/terms`만 더했다. live 규칙이 `4787e85`의 repo 규칙과 같음을 확인한 뒤 후보를 만들었고 후보 규칙은 `7c6ef373`의 repo 규칙과 문자열이 같다. 바뀐 줄은 2줄(`path`의 한 줄)이다.
- 원본은 `publish001-ops-20261022/config.yml.pre-publish001`(0600·1035bytes, SHA256 `a9e24b10…bf6fac`)에 보관했다.
- `cloudflared tunnel --config <후보> ingress validate` OK. 후보 `ingress rule`: `/privacy`·`/terms`·`/`·`/home`·`/home/calendar/connect`·`/lang/ko`·`/v1/text/seen`은 relay 규칙(#0), `/privacy/`·`/privacy/x`·`/privacyx`·`/terms/`·`/terms.html`·`/owner`·`/healthz`·`/v1/owners`는 404 규칙(#1)이다.
- `python3 scripts/check_public_ingress.py`는 반영 전(로컬 repo·서버 checkout)과 반영 뒤(서버 checkout) 모두 통과했다. 반영 뒤 live 규칙과 repo 규칙이 같다.
- `knowslink-cloudflared-1`을 한 번 재시작했다(22:47:14Z). 연결 4개가 다시 등록됐다. `ERR` 5줄은 재시작 때의 `Connection terminated` 4줄과 `no more connections active and exiting` 1줄이다([cferr.log](cferr.log)). 이전 교훈에 따라 재시작 관찰 반복 호출은 하지 않았다.

## 보존 ([preservation.json](preservation.json))

assertion 41개가 통과했다. 비교 단계는 baseline·before·frozen·after·final이다. 값은 read-only 조회로 메모리에서 count와 digest만 만들었다.

- 회원·agent·owner·pair·신원·친구 등 상태 map digest와 count가 5단계에서 같다(Rates 제외). `Calendars` 키는 5단계 모두에 없다. Remote의 Clients·Attempts·Codes·Installs·Subscriptions digest와 Inbox·구독 메타가 같다.
- **Rates**와 **Remote.Tokens**는 volatile이다. Rates는 요청 시각 기록으로 변하고, Remote.Tokens는 final에서만 2→1로 줄었다. 그 사이 이 작업이 하지 않은 정리나 만료일 수 있으나 집계만으로 원인을 확정하지 못한다. 권한 변경으로 해석하지 않으며 Clients·Installs·Subscriptions는 같다.
- schema SHA256·`tunnel.uuid`·`access.aud`·Tunnel credential·파일 mode가 같다. `.env`는 승인된 키 줄 삭제로 한 번, `tunnel/config.yml`은 승인된 ingress 반영으로 한 번 바뀐다. 기존 백업19개의 hash/크기/mode가 같고 새 백업 1개만 추가됐다(총 20).
- relay 전환 구간에서 바뀐 컨테이너는 relay와 migrate뿐이다. final에서 추가로 바뀐 것은 KnowsLink cloudflared(승인된 재시작)뿐이다. Postgres를 포함한 나머지 121개 컨테이너(공유 서비스 포함)는 ID·이미지·시작 시각이 같다. 새 컨테이너는 보존한 복원 fixture와 rollback 검사 컨테이너 둘이다.
- 호스트 cloudflared PID, orca/s8/mcp 응답, loopback 8080/5432 listener가 같다.
- 보존 assertion은 짧은 전환 중의 일치다. 정상 만료·cleanup·사용자 요청에 따른 이후 변경을 금지하는 보장은 아니다.

## 공개 검사

모두 통과했다. 실제 로그인·송신·pull·ACK·token 소비는 없다. 호출은 1.5초 이상 간격으로 했고 429는 없었다. 응답 본문과 cookie 값은 보존하지 않았다.

[public-smoke-policy.json](public-smoke-policy.json)(이번 변경 검사, `all_pass`)
- `GET /privacy`·`/terms`는 한국어(기본)·영어(쿠키 `__Host-kl_lang=en`) 모두 비로그인 200, `text/html`, `Set-Cookie` 없음, 시행일 `2026-10-22` 포함, `<h1>` 있음이다.
- 이웃 경로 `/privacy/`·`/privacy/x`·`/privacyx`·`/terms/`·`/terms.html`·`/Privacy`·`/policy`는 모두 404다(다른 경로 노출이 늘지 않았다).
- `GET`·`POST /home/calendar/connect`, `POST /home/calendar/disconnect`, `GET /home/calendar`는 모두 404다(RELEASE-003 때의 401/303/405에서 바뀜. D018 제거와 일치).
- 첫페이지 `/`: `#start` 영역의 Google 로그인 버튼 아래 ‘로그인하기 전에 서비스 약관과 개인정보처리방침을 확인하세요’ 줄(링크 `/terms`·`/privacy`)과 푸터의 두 링크가 있다. Google 폼 1개이고 `Set-Cookie`가 없다. 첫페이지 본문에 calendar·freebusy·캘린더·일정 조회 단어가 없다. 이 `#start` 줄은 운영 Google 로그인이 켜져 있어 보이는 줄이다. 브라우저 렌더링 확인은 하지 않았다(HTML 문자열로 확인).

[public-smoke.json](public-smoke.json)(직전 smoke 재실행, `all_pass`): 기존 안전 검사 21/21, 한국어 첫페이지 필수 문구·금지 문자열·Google 폼 1개·`<details>` 2개·`Set-Cookie` 없음, 새 API 무인증 경계(`/v1/contacts/{agent}`·`/v1/text/seen` 401)와 이웃 404 6개.
[public-smoke-lang.json](public-smoke-lang.json)(`all_pass`): `/lang/*` 6/6, 이웃 경로 8/8, 첫페이지 언어 5/5, 무인증 회원 경계 3/3.

기존 친구·메시지 동작은 실계정 발송 없이 확인 가능한 범위로 본다. 친구 1·pair 13·agent 28·Remote 설치 2의 digest가 전환 전후 같고, 무인증 API 경계와 `verify.py local`이 통과했다. 실제 송수신은 하지 않았다.

## P-03·P-04 실제 운영 값 ([retention.log](retention.log))

- **P-03 백업 보존**: 문안은 ‘백업은 정기적으로 교체하며 교체할 때 함께 삭제합니다’다. 실제로는 **자동 교체·삭제가 없다**. 서버에 cron·timer·`beta.sh`의 정리 단계가 없고(`dpkg-db-backup.timer`는 DB 백업이 아니다), Postgres WAL archive는 꺼져 있다(`archive_mode=off`). `knowslink-state/backups/`에 개인정보가 담긴 dump가 2026-10-03부터 2026-10-22까지 20개 쌓여 있고 14일보다 오래된 것이 6개다. 지금까지 정리는 수동이며 이번 작업도 기존 백업을 지우지 않았다(승인 범위 밖). **문안과 실제가 다르다.** 문안은 고치지 않았다. designer가 (a) 실제 주기를 정해 문안을 쓰거나 (b) ops가 정기 정리를 만들고 문안을 맞출지 정해야 한다. 어느 쪽이든 법률 검토(P-05) 전 초안이다.
- **P-04 이메일 업체**: 운영 relay 기동 줄이 `email login configured=false`다(SMTP 미설정). 이메일 발송 업체 행은 생략을 유지한다. 문안 변경은 필요 없다.

## 한계·미실행

- 미실행: 실계정 로그인·발송, 브라우저 렌더링·시각 확인, 키보드·스크린리더, 실제 rollback, Cloud Console·Access·DNS 변경, 부하·장시간 검사, 다른 브라우저. Cloud Console(calendar.freebusy scope 제거, Calendar API 비활성화, 브랜딩 URL, 앱 게시)은 coor와 사용자의 몫이다.
- 정책 문안은 법률 검토(P-05) 전 초안이다.
- 제품 코드 gate는 SAR-PUBLISH-001 각 단위의 수락 증거(독립 리뷰·tester QA 62/62·designer 시각 검수)를 재사용한다. 이번 운영 검사와 구분한다. 이 작업에서 제품 코드·테스트·의존성·migration은 바꾸지 않았다.
- rollback: `knowslink-rollback:4787e85-publish001-20261022` 이미지로 되돌릴 수 있다. 연결이 0이었으므로 되돌려도 잃는 캘린더 연결은 없다. 키는 이제 운영 env에 없고 `env.pre-publish001`에 보관돼 있다. 정리 뒤에는 상태 JSON을 되돌릴 일이 없다.

## 보존 위치

- 로컬 원본: `C:/Users/shin/AppData/Local/KnowsLinkDevTools/publish001-ops/`(bundle·스크립트·원 로그).
- 서버 private 작업물: `/home/shin/deploy/knowslink-state/publish001-ops-20261022/`(bundle·관측·env 보관본·ingress 원본).
- 서버 백업: `/home/shin/deploy/knowslink-state/backups/4787e85-pre-publish001-20261021T224540Z.dump`(0600).
- 보존 컨테이너: `knowslink-publish001-rehearse-20261022`(복원 fixture), `knowslink-publish001-rollback-inspect-20261022`. 이미지 tag `knowslink-rollback:4787e85-publish001-20261022`.
- Git 증거: 이 폴더의 redacted JSON/log와 [commands.txt](commands.txt)(실행 스크립트 전문). 원 credential·token·cookie 값·캘린더 키 값·실메시지·DB row는 포함하지 않는다.

## 문서 HEAD의 FullOps 검사

- 기준 `7c6ef373`, 대상 HEAD `258df551`(결과 문서 커밋)에서 `lint.py`는 ERROR0·WARNING1·실행 불가1(`make lint`)이고 `make test`는 skipped, 종료1이다. 실행 불가는 Windows `make` 부재이며 직전 OPS와 같다. 새 PASS로 보고하지 않는다. 제품 lint·테스트는 SAR-PUBLISH-001 각 단위의 수락 증거를 재사용한다. 이 작업은 제품 코드·의존성을 바꾸지 않았다.
- WARNING 1건은 `ops-guide.md`의 SIZE-001(513줄, 이전 509)이었다. 새 절의 bullet 4개를 합쳐 509줄로 되돌렸다. 이미 상한 500을 넘던 기존 문서이며 분할은 이번 소유 범위 밖이다.
- `lint.py`는 작업 트리가 깨끗해야 한다. `deliverables.py --stamp --path`가 잘못된 경로 `.fullops-squad/.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-OPS-test/report.md`에 보고서 복사본을 만들었다(원본과 같은 front matter). 이 저장소는 worker의 삭제·이동 정리를 금지하므로 지우지 않고 보존했다. Git에는 넣지 않고 공유 `.git/info/exclude`에 `/.fullops-squad/.fullops-squad/` 한 줄을 더해 lint가 돌게 했다. coor가 폴더를 지운 뒤 그 줄을 빼면 된다.
