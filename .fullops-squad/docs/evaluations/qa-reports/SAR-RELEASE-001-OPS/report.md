---
title: main 통합 제품 운영 배포와 공개 ingress 보완 결과
status: draft
updated: 2026-10-17
owner: ops
tasks: [SAR-RELEASE-001-OPS]
summary: 제품 1dd6800을 기존 운영 relay에 배포한 백업·격리복원·보존·공개 검사 결과와 새 API 두 경로의 공개 ingress 보완을 기록한다.
---

# SAR-RELEASE-001-OPS 결과

제품 `1dd6800d5d4238e2060f99e98eaab437ea06682c`를 기존 link.knowslog.com 운영 relay에 배포했다. 새 private 백업과 격리 복원, 보존 38개 assertion, 공개 smoke를 통과했다. 새 API 두 경로가 공개 주소에서 404인 문제를 발견했다. 사용자 승인을 받아 repo와 live Tunnel config를 함께 보완했고 두 경로가 401이 됨을 확인했다.

## 범위·기준·시간

- 사용자 지시는 "운영배포진행해"(2026-10-17)다. 준비 HEAD는 `2ef95cb577f3afbf59a0102c9c5cdbce2006f5a4`이고 브랜치는 `fullops/ops`다. 준비 HEAD의 제품 파일은 `1dd6800`과 같다(`git diff 1dd6800 2ef95cb -- . ':!.fullops-squad'` 비어 있음).
- Task `task_b6c9974de87d`, Dispatch `ctx_b1e50f2d98e4`, worker `term_8c976e59-67d6-4be5-9c21-7546756ff1cc`, 복귀 coor `term_e3a3e5f6-e163-420f-bba6-51b428f7671f`, Run `run_0747179f36f3`이다.
- 흡수한 SAR-LANDING-GUIDE-001-OPS의 "QA 수락·main 통합 전 배포하지 않는다" 전제는 main 통합과 이번 지시로 충족됐다. 그 지시서의 운영·보존 계약만 이어받았다.
- 하위 에이전트는 쓰지 않았다. 직렬 운영 전환은 직접 수행했다. 독립 결과 리뷰는 coor 후속이다.
- `beta.sh deploy/backup/restore-verify`와 `verify.py regression/public`은 실행하지 않았다. 삭제 포함 절차와 과거 Access302·projects baseline 전제 때문이다. `verify.py local`은 세 번 실행해 모두 통과했다.

| 항목 | 값 |
|---|---|
| fresh baseline | 2026-10-17T09:36:12Z. 서버 HEAD `8a93e38`·clean. relay 이미지 `sha256:482d3197…381343d`·healthy·restart0·시작 2026-10-14T11:31:08Z |
| 전환 구간 | STOP_BEGIN `2026-10-17T09:37:48Z`(18:37:48 KST) → UP_DONE `2026-10-17T09:38:09Z`(18:38:09 KST). 약21초 |
| 최종 relay | 이미지 `sha256:4893d0f5…0098999`·시작 09:38:03Z·healthy·restart0 |
| Tunnel 재시작 | 2026-10-17T09:44:36Z. 약2초 끊김(아래 절) |

## 배포 입력 확인

- bundle SHA256 `f8e1b2de56887bcbc60286bfad351c7c6cdfabb8dc0fca61a94446aac1ba6c60`다. 서버에서 hash 대조·`bundle verify`·`FETCH_HEAD`·`8a93e38`→`1dd6800` 조상 관계를 확인했다.
- `db deploy compose.yaml Dockerfile go.mod go.sum Makefile sqlc.yaml`의 diff는 0이다. 새 migration은 없다. 제품 diff는 지시서와 같은 39파일 +2051/−121이다.
- 고정 clean checkout에서 relay/migrate build가 EXIT0이었다. 빌드 중 기존 relay는 healthy였다.
- 새 build 전에 현재 이미지에 `knowslink-rollback:8a93e38-release001-20261017`을 붙였다. inspect 결과 ID가 `sha256:482d3197…`로 같다. 기존 tag는 덮지 않았다. 이전 SHA 소스를 다시 빌드하지 않고 원 이미지를 보존했다.
- rollback 이미지는 network none·read-only 컨테이너에서 `/app/relay`·`/app/migrate` 실행 파일과 SHA256을 확인했다(`knowslink-release001-rollback-inspect-20261017` exited0 보존). DB 연결이나 앱 기동 검사는 아니다.

## 백업·격리 복원·재개

| 항목 | 실제 결과 |
|---|---|
| writer 고정 | relay stop 뒤 epoch `1297596`이 2초 간격으로 같았다. 다른 DB client 0 |
| 새 백업 | `backups/8a93e38-pre-release001-20261017T093752Z.dump`. 상위 0700·파일 0600·11104bytes |
| 백업 SHA256 | `8498f8ebfce0f57dab69b9ad2ac64198f978176eafefe4beb23d7cca89ea3452` |
| list·복원 | `pg_restore --list` TOC 25줄. network none·host port 없음·postgres 사용자로 복원 |
| 복원 검증 | 복원 `data::text` md5가 frozen 운영 값과 같다. goose2·guard1·rows1·epoch1297596·HumanModel1 |
| 거부 경계 | 격리 DB에서 HumanModel 제거 UPDATE가 거부됐다. live DB에는 거부용 쓰기를 하지 않았다 |
| 재개 | 미리 만든 이미지로 `up -d --no-build --wait relay`. migrate EXIT0·"SQL migrations completed"·goose2 유지 |
| 로그 | relay error/panic 0·Google configured=true·새 Tunnel 프로세스 error 0·연결 4개 등록 |

복원 fixture `knowslink-release001-rehearse-20261017`은 exited0·network none·port 없음으로 보존했다. 운영 DB를 과거 dump로 덮지 않았다. Down·volume 삭제·prune·이미지 삭제는 없다. 서버 `deploy-history.log`에 한 줄을 추가했다.

## 보존

[preservation.json](preservation.json)의 assertion 38개가 통과했다. 비교 단계는 baseline·before·frozen·after·final이다. 값은 read-only 조회로 메모리에서 count와 digest만 만들었다.

- 회원2·agent28·owner24·친구1·pair13·신원2 등 상태 map digest가 5단계에서 같다. Remote 설치·구독·Events 메타·Inbox 0도 같다.
- Messages는 baseline부터 0이다. 10월 14일의 11건은 이번 작업 전에 보존 기간이 끝난 값이다. 이번 작업은 이를 바꾸지 않았다.
- Rates만 달랐다. 요청 시각 기록과 window 정리에 따라 변하는 volatile map이다. 권한 변경으로 해석하지 않는다. 정확한 원인 요청은 집계로 확정할 수 없다.
- schema SHA256·`.env`·`tunnel.uuid`·`access.aud`·Tunnel credential·파일 mode가 같다. 기존 백업16개의 hash/크기/mode가 같고 새 백업 1개만 추가됐다.
- 전환 구간에서 바뀐 컨테이너는 relay와 migrate뿐이다. Postgres와 나머지 컨테이너(공유 서비스 포함)는 ID·이미지·시작 시각이 같다. 이후 Tunnel 보완으로 `knowslink-cloudflared-1`도 재시작됐다. 새 컨테이너는 보존한 복원 fixture 하나다.
- 호스트 cloudflared PID `3637725`, loopback 8080/5432 listener, orca/s8/mcp 응답이 같다.

보존 assertion은 짧은 전환 중의 일치다. 정상 만료·cleanup·사용자 요청에 따른 이후 변경을 금지하는 보장은 아니다.

## 공개 smoke

[public-smoke.json](public-smoke.json)은 모두 통과했다(`all_pass=true`). 실제 로그인·송신·pull·ACK·token 소비는 없다. 응답 본문과 cookie는 보존하지 않았다.

- 기존 안전 검사 21/21: `/` 200과 no-store/CSP/no-referrer/nosniff, 무세션 `/home*` 303과 expired redirect, 보호 경로 404, SSE·MCP 무인증 401, OAuth well-known 200, consent 401, connect 405, 공유 orca/s8/mcp 200/200/401.
- 새 첫페이지: 200·no-store·보안 헤더·`Set-Cookie` 없음. 소개·시작 순서·설치와 연결(Grok·Dots `<details>` 2개)·가상 예시·로그인 폼(Google 1개, 이메일 0개)·기한/보관 안내 문구가 모두 있다. 구현 경로·`<img>`·`<details open>` 같은 금지 문자열은 없다.
- 새 API 무인증 경계: `GET /v1/contacts/{agent}`·`GET /v1/contacts/{%3A 인코딩 agent}`·`POST /v1/text/seen`이 401이다.
- 이웃 경로 노출 없음: `/v1/contacts`, `/v1/contacts/`, `/v1/contacts/a/b`, `/v1/contacts/Agent`, `/v1/text/seen/extra`, `/v1/text/unknown`이 404다.
- 관찰: 공개 응답에는 `<script type="module" src="https://static.cloudflareinsights.com/beacon.min.js/…">` 하나가 있다. Cloudflare edge가 붙인 Web Analytics 주입이다. origin(loopback) 응답에는 `<script`가 0건이다. 제품 코드의 문제로 보지 않는다. 브라우저에서 CSP가 이 script를 막는지는 확인하지 않았다.

## 발견한 문제와 승인된 보완: 새 API 두 경로의 공개 ingress

첫 공개 smoke에서 `GET /v1/contacts/{agent}`와 `POST /v1/text/seen`이 `https://link.knowslog.com`에서 404였다. loopback에서는 둘 다 401이다. 원인은 공개 Tunnel ingress의 `path` 정규식에 두 경로가 없는 것이다. repo의 `deploy/knowslink/tunnel/public-ingress.yml`과 live `config.yml`이 같은 규칙이었다. 이 상태로는 Claude Code 세션 메시지의 새 API 두 개가 공개 주소로 도달하지 못한다.

배포 입력에 `deploy/` 변경이 없어 사전 diff로는 발견할 수 없었다. 두 경로를 여는 것은 외부 노출 변경이라 coordinator에게 `ask`로 멈췄다. 사용자 결정(2026-10-17)은 "repo+live 함께 반영"이었고 범위는 두 경로로 한정됐다.

변경 내용은 정규식 대안 2개 추가다. `/v1/text/(send|pull|persist|ack|events)`를 `…|events|seen)`으로 바꿨다. `/v1/contacts/[a-z][a-z0-9_:-]*`를 추가했다. 이 agent 패턴은 `/v1/keys/…`와 relay의 `agentPattern`과 같다. 다른 경로는 추가하지 않았다. 목록 경로 `/v1/contacts`는 계속 막힌다.

- repo: `deploy/knowslink/tunnel/public-ingress.yml`과 `scripts/check_public_ingress.py`(허용 2개·차단 이웃 경로 5개 기대값 추가). 이 `deploy/` 수정은 이번 과제에 한해 ops 소유로 인정받았다. 서버 checkout은 배포 SHA `1dd6800` 그대로이므로 이 수정은 다음 main 반영 때 서버에 도착한다.
- repo 검사: 두 검사는 `make lint`의 `lint-config`에 속한다. 로컬 `check_public_ingress.py` 통과. Windows에서 `test_public_config.sh`는 파일 mode 644 ≠ 600로 실패했다(NTFS에는 0600이 없다. 변경과 무관한 환경 한계). 같은 두 파일을 `git archive 1dd6800` 사본 위에 덮고 서버 Linux에서 실행해 두 검사 모두 EXIT0이다([ingress-test.log](ingress-test.log)).
- live 원본 백업: `release001-ops-20261017/config.yml.pre-release001`, 0600, SHA256 `f5abe860a019d5379cbbd5499eb309a645c4ae1e8f0b60fab6fea22326458512`.
- 후보 검증: 변경 줄 2개. `cloudflared tunnel ingress validate` OK. 경로별 `ingress rule`에서 새 두 경로와 `/`는 규칙 #0, 이웃·보호 경로는 404 fallback #1이다.
- 첫 시도는 `ingress rule | grep -m1`의 SIGPIPE로 종료 141이었고, live 변경 전이었다([attempt1](ingress-apply.attempt1.log)). 같은 후보로 grep을 고쳐 다시 실행했다.
- 적용: 2026-10-17T09:44:36Z에 live config를 제자리 갱신하고(mode 0600 유지) `docker restart knowslink-cloudflared-1`을 실행했다. 서버에서 0.4초 간격으로 공개 `/`를 관찰했다([restart-poll.txt](restart-poll.txt)). 502 → 530 → 200이 09:44:36.8~09:44:38.7이다. 끊김은 약2초다. 새 프로세스는 연결 4개를 등록했고 error가 0이다. 종료 때의 ERR 5줄(09:44:36Z)은 의도한 재시작의 종료 기록이다.
- 동일성: live의 `path` 규칙이 repo 규칙(SHA256 `f97e0fd95a61164afe854564f1d8314eaf1dd535767fd7e18bb46caa3d648cb7`)과 같다. 추가된 대안은 `seen`과 `/v1/contacts/[a-z][a-z0-9_:-]*` 둘뿐이다([final.log](final.log)).
- 되돌림은 필요하지 않았다. 필요하면 백업 config를 복사하고 cloudflared를 재시작한다.

## 429 관찰

Tunnel 재시작 직후 내 서버 측 관찰 루프가 `/`를 약 2.5회/초로 90번 호출했다. 09:44:51Z부터 relay를 거치는 공개 경로가 429가 됐다. 같은 시각 로컬 PC의 smoke도 429였다. 09:47:05Z에는 200으로 돌아왔다. 요청 제한의 단위(IP별인지 공유인지)는 확인하지 않았다. 서버와 PC가 같은 출구 IP일 가능성이 있다. 다른 사용자가 영향을 받았는지는 알 수 없다. 이후 공개 검사는 간격을 두고 다시 실행해 통과했다.

## gate 재사용·미실행·한계

- 제품 코드 gate는 각 기능의 수락 증거(독립 리뷰·tester QA, main 통합 기록)를 재사용한다. 이번 운영 검사와 구분한다. SAR-SESSION-PEER-002 코드 리뷰는 사용자가 생략을 선택했다(waive). 새 API 권한 경계는 tester의 직접 HTTP 20건 기록을 따른다. 이번 공개 무인증 401은 그 증거를 대체하지 않는다.
- 이 작업에서 제품 코드·테스트·의존성은 바꾸지 않았다. 바꾼 파일은 위 ingress 두 파일과 ops 문서다. lint는 결과 문서 HEAD에서 실행해 worker_done으로 전달한다.
- 미실행: 실계정 로그인·메시지 왕복, 인증된 `GET /v1/contacts/{agent}`·`POST /v1/text/seen`의 성공 경로, pull/ACK/token 소비, 실제 rollback, Access·DNS 변경, 부하·장시간 검사.
- 한계: 브라우저 시각 확인은 하지 않았다(UI 코드 변경 없음, 지시서 기준). 공개 페이지는 HTTPS 내용·헤더만 확인했다. rollback 이미지는 파일 존재·hash까지 확인했고 앱 기동은 하지 않았다. 새 relay는 manifest list ID로 보인다. Docker 이미지 저장소가 보고하는 ID 형식이다.
- D10(`module-design.md`)은 갱신하지 않았다. 공개 경로 목록이 없고 `render-public-config`를 일반 설명만 하기 때문이다.
- 오래된 포인터: `docs/README.md`가 "2026-10-14 배포 제품 8a93e38"을 말한다. ops 소유 파일이 아니라 수정하지 않았다. coor 또는 DEV가 갱신한다.
- 후속(DEV): 새 공개 경로를 추가하는 변경은 `public-ingress.yml`·`check_public_ingress.py`를 함께 바꿔야 한다. 이번 누락은 배포 전 단위 gate로 잡히지 않았다. Grok 로컬 bundle과 Claude Code hook 패키지는 서버 배포로 갱신되지 않는다.

## 문서 HEAD의 FullOps 검사

- 기준 `1dd6800`, 대상 HEAD `13f9e3c`(보고 보완 전). 기본 PATH에서 `lint.py`는 ERROR0·WARNING0·실행 불가2(`make lint`, `make test`), 종료1이다. 직전 OPS와 같은 Windows `make` 부재다. 새 PASS로 보고하지 않는다.
- w64devkit의 `make`를 PATH 앞에 둔 보조 실행은 ERROR2였다. `core.autocrlf=true` 작업 사본이 CRLF라서 `gofmt -l`이 Go 파일을 나열하고, `TestFrozenParsingAndSigning`이 registry 무결성 오류를 냈다. 이번 작업은 Go 파일을 바꾸지 않았다. 환경 한계이며 제품 결함으로 판정하지 않았다.
- `make lint`의 ingress 두 검사는 서버 Linux에서 통과했다(위 절). 산출물 strict는 13개·문제0·경고0이다. `git diff --check` 통과.
- 원 결과: 로컬 `release001-ops/final-lint-default.json`, `lint-with-make.json`.

## 보존 위치

- 로컬 원본: `C:/Users/shin/AppData/Local/KnowsLinkDevTools/release001-ops/`(bundle·스크립트·원 로그).
- 서버 private 작업물: `/home/shin/deploy/knowslink-state/release001-ops-20261017/`(bundle·관측·backup.list·config 백업·후보·ingress-check 사본).
- 서버 백업: `/home/shin/deploy/knowslink-state/backups/8a93e38-pre-release001-20261017T093752Z.dump`(0600).
- 보존 컨테이너: `knowslink-release001-rehearse-20261017`(복원 fixture), `knowslink-release001-rollback-inspect-20261017`(rollback artifact 검사). 이미지 tag `knowslink-rollback:8a93e38-release001-20261017`.
- Git 증거: 이 폴더의 redacted JSON/log와 [commands.txt](commands.txt). 원 credential·token·cookie·실메시지·DB row는 포함하지 않는다. [post-check2.log](post-check2.log)의 assertion 실패는 의도한 cloudflared 재시작의 종료 ERR 5줄을 "전환 이후 Tunnel error"로 센 결과다. 새 프로세스 기준 0건은 [tun2.log](tun2.log)에 있다.
