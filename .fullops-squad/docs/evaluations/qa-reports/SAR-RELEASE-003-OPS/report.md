---
title: main 통합 제품 4787e85 운영 배포와 캘린더 키 설정 결과
status: draft
updated: 2026-10-21
owner: ops
tasks: [SAR-RELEASE-003-OPS]
summary: 일정 조회를 포함한 제품 4787e85를 기존 운영 relay에 배포하고 운영 env에 KNOWSLINK_CALENDAR_KEY를 새로 만든 백업·격리복원·보존·egress·공개 검사 결과를 기록한다.
---

# SAR-RELEASE-003-OPS 결과

제품 `4787e85a0b6952970a861b6b100b3561b9cf30a0`를 기존 link.knowslog.com 운영 relay에 배포했다. 운영 env 정본에 `KNOWSLINK_CALENDAR_KEY`를 서버에서 새로 만들어 추가했다(값 비기록). 새 private 백업과 격리 복원, 보존 assertion 40개, egress 두 주소와 공개 smoke를 통과했다. ingress·cloudflared·Cloud Console·실계정은 건드리지 않았다.

## 범위·기준·시간

- 사용자 결정 D016 '배포 진행'(2026-10-21)과 D011 답(OAuth 앱 Testing, 사용자는 테스트 사용자)이다. 준비 HEAD는 브랜치 `fullops/ops`의 `375b5c26`이고 제품 파일은 main `4787e85`와 같다.
- Task `task_da2868667888`, Dispatch `ctx_d0f27dc316dd`, worker `term_2352835c-2932-4e31-be1a-21bc4a3f80fd`, 복귀 coor `term_fec2befe-1fd1-4564-b0cc-cd9e6b157161`이다.
- 하위 에이전트는 쓰지 않았다. 직렬 운영 전환은 직접 수행했다.
- `beta.sh deploy/backup/restore-verify`와 `verify.py regression/public`은 실행하지 않았다(삭제 포함 절차, 과거 baseline 전제). `verify.py local`은 baseline과 전환 직후, 마지막 로그 확인에서 세 번 통과했다.

| 항목 | 값 |
|---|---|
| fresh baseline | 2026-10-21T04:47:02Z. 서버 HEAD `7b4f9bc`·clean. relay 이미지 `sha256:9646febaf412…85443f66`·healthy·restart0·시작 2026-10-20T13:50:56Z. 컨테이너 113·백업 18 |
| 전환 구간 | STOP_BEGIN `2026-10-21T04:48:42Z`(13:48:42 KST) → UP_DONE `2026-10-21T04:48:59Z`(13:48:59 KST). 약17초 |
| 최종 relay | 이미지 `sha256:0e0eb8b44e47…ae7fd49`(manifest list ID)·시작 04:48:53Z·healthy·restart0 |

## 배포 입력 확인

- bundle SHA256 `fb1fead7dcf09a520b747ec2c6117373d6917228327882aaed5a487e053e881a`(4846107bytes)다. 서버에서 hash 대조·`bundle verify`·`FETCH_HEAD`=`4787e85`·`7b4f9bc`→`4787e85` 조상 관계를 확인했다([build.log](build.log)).
- `db Dockerfile go.mod go.sum Makefile sqlc.yaml`의 diff는 0이다. migration과 새 의존성은 없다. `compose.yaml`은 `KNOWSLINK_CALENDAR_KEY: ${KNOWSLINK_CALENDAR_KEY:-}` 한 줄만 늘었다. `deploy scripts` diff는 `public-ingress.yml`·`check_public_ingress.py`(직전 배포에서 이미 live에 반영된 `/lang/(ko|en)`)뿐이다. 제품 diff는 44파일 +2138/−78이고 `internal adapters cmd`는 31파일 +2068/−57이다.
- 고정 clean checkout에서 relay/migrate build가 EXIT0이었다. 빌드 중 기존 relay는 healthy였다.
- 새 build 전에 현재 이미지에 `knowslink-rollback:7b4f9bc-release003-20261021`을 붙였다. inspect 결과 ID가 `sha256:9646febaf412…`로 같다. 기존 tag(`1dd6800-release002-20261020`·`8a93e38-release001-20261017`·`2747fe9-refactor002-20261014`)는 덮지 않았다.
- rollback 이미지는 network none·read-only 컨테이너 `knowslink-release003-rollback-inspect-20261021`(exited0 보존)에서 `/app/relay`·`/app/migrate` 실행 파일과 SHA256을 확인했다([rollback-verify.log](rollback-verify.log)). DB 연결이나 앱 기동 검사는 아니다.

## 운영 env 키

[env.log](env.log)의 결과다. 값은 출력·기록하지 않았고 길이와 존재만 확인했다.

- 서버 `/home/shin/deploy/knowslink-state/.env`는 mode 0600·450bytes였고 `KNOWSLINK_CALENDAR_KEY` 줄이 없었다. 이전 내용을 `release003-ops-20261021/env.pre-release003`(0600·450bytes)에 보관했다. 이전 `.env` SHA256은 `de221b94f2bddc66da85b9b320bc8d3742e607fb52eecdbf991c6aecfb8d8d99`다.
- `openssl rand -base64 32`로 서버에서 만든 값을 한 줄 추가했다. 결과는 0600·518bytes, 키 줄 1개, base64 디코드 32바이트(텍스트 44자)다. 보관본이 새 파일의 앞부분과 바이트 단위로 같다(다른 줄 불변).
- `docker compose config`의 relay 환경에 키가 들어오고(길이 44), 기동 뒤 relay 컨테이너 안에서도 길이 44임을 확인했다([deploy.log](deploy.log)).
- relay 기동 로그의 상태 줄(`relay starting; … google login configured=true; …`)에는 캘린더 상태가 없다. 잘못된 키는 기동 실패이므로 기동·healthy 자체가 유효한 키의 증거다. 추가로 인증 없는 `POST /home/calendar/connect`가 503(키 없음)이 아닌 401을 돌려준다([public-smoke-cal.json](public-smoke-cal.json)). 키가 있어야 `googleStart`가 회원 확인 단계까지 간다. 로그인 회원 화면의 캘린더 절은 실계정이 필요해 확인하지 않았다.
- 키 교체·분실 절차는 미정이다. 잃으면 저장된 모든 캘린더 연결이 열리지 않아 회원이 다시 연결해야 한다.

## 백업·격리 복원·재개

| 항목 | 실제 결과 |
|---|---|
| writer 고정 | relay stop 뒤 epoch `1627651`이 2초 간격으로 같았다. 다른 DB client 0 |
| 새 백업 | `backups/7b4f9bc-pre-release003-20261021T044846Z.dump`. 상위 0700·파일 0600·11280bytes |
| 백업 SHA256 | `e6b31b28967f5fa205f678e3eff171bef1d3f248b786ba7c3d9d69e0ad8f62be` |
| list·복원 | `pg_restore --list` TOC 25줄. network none·host port 없음·postgres 사용자로 복원 |
| 복원 검증 | 복원 `data::text` md5가 frozen 운영 값과 같다. goose2·guard1·rows1·epoch1627651·HumanModel1 |
| 거부 경계 | 격리 DB에서 HumanModel 제거 UPDATE가 거부됐다. live DB에는 거부용 쓰기를 하지 않았다 |
| 재개 | 미리 만든 이미지로 `up -d --no-build --wait relay`. migrate EXIT0·"SQL migrations completed"·goose2 유지 |
| 로그 | relay 로그 1줄(기동 줄)·error/panic 0·restart 0·Google configured=true |

복원 fixture `knowslink-release003-rehearse-20261021`은 exited0·network none·port 없음으로 보존했다. 운영 DB를 과거 dump로 덮지 않았다. Down·volume 삭제·prune·이미지 삭제는 없다. 서버 `deploy-history.log`에 한 줄을 추가했다.

## 보존

[preservation.json](preservation.json)의 assertion 40개가 통과했다. 비교 단계는 baseline·before·frozen·after·final이다. 값은 read-only 조회로 메모리에서 count와 digest만 만들었다.

- 회원·agent·owner·pair·신원·친구 등 상태 map digest와 count, Remote 설치·구독·Events 메타·Inbox 0이 5단계에서 같다. 직전 배포가 쓴 `Friend.CodeTail:""`이 이미 있어 이번에는 Friends 직렬화 차이도 없다. 새 `Calendars` map은 5단계 모두에 없다(연결한 회원이 없다). 기존 상태 JSON이 새 캘린더 필드(`Calendars`·`Gate.Lookup`·`GoogleAttempt.Calendar`) 없이 읽힌다는 증거다.
- schema SHA256·`tunnel.uuid`·`access.aud`·Tunnel credential·Tunnel config·파일 mode가 같다. `.env`는 mode 0600이고 해시는 승인된 키 한 줄 추가로 한 번만 바뀐다. 기존 백업18개의 hash/크기/mode가 같고 새 백업 1개만 추가됐다.
- 전환 구간에서 바뀐 컨테이너는 relay와 migrate뿐이다(cloudflared 재시작 없음). Postgres와 나머지 111개 컨테이너(공유 서비스 포함)는 ID·이미지·시작 시각이 같다. 새 컨테이너는 보존한 복원 fixture와 rollback 검사 컨테이너 둘이다.
- 호스트 cloudflared PID `2039731`, loopback 8080/5432 listener, orca/s8/mcp 응답이 같다.
- **Rates**는 요청 시각 기록과 window 정리로 변하는 volatile map이다(2→3건). 권한 변경으로 해석하지 않는다. 정확한 원인 요청은 집계로 확정할 수 없다.
- 준비 단계의 build가 서버 checkout을 새 SHA로 옮기므로 before부터 HEAD는 `4787e85`이고 baseline만 `7b4f9bc`다.

보존 assertion은 짧은 전환 중의 일치다. 정상 만료·cleanup·사용자 요청에 따른 이후 변경을 금지하는 보장은 아니다.

## egress

운영 relay 컨테이너(busybox `wget`)에서 인증 없이 호출했고 응답 코드만 남겼다([egress.log](egress.log), [egress2.log](egress2.log)).

| 주소 | GET | 빈 POST |
|---|---|---|
| `https://www.googleapis.com/calendar/v3/freeBusy` | 404 | 403 |
| `https://oauth2.googleapis.com/revoke` | 404 | 400 |

둘 다 Google이 응답했다. 인증·토큰을 쓰지 않았으므로 실제 조회·철회 성공은 아니다. Tunnel ingress는 이번에 바꾸지 않았고 live `path` 규칙은 서버 checkout의 repo 규칙과 같다([final.log](final.log)).

## 공개 smoke

모두 통과했다. 실제 로그인·송신·pull·ACK·token 소비는 없다. 호출은 1.5초 이상 간격으로 했고 429는 없었다. 응답 본문과 cookie 값은 보존하지 않았다.

[public-smoke.json](public-smoke.json)(직전 smoke 재실행, `all_pass`)
- 기존 안전 검사 21/21, 한국어 첫페이지 필수 문구·금지 문자열·Google 폼 1개·이메일 폼 0개·`<details>` 2개·`Set-Cookie` 없음, 새 API 무인증 경계(`/v1/contacts/{agent}`·`/v1/text/seen` 401)와 이웃 404 6개.

[public-smoke-lang.json](public-smoke-lang.json)(`all_pass`)
- `/lang/*` 6/6, 이웃 경로 8/8, 첫페이지 언어 5/5, 무인증 회원 경계 3/3.

[public-smoke-cal.json](public-smoke-cal.json)(이번 변경 검사, 통과)
- 무세션 `POST /home/calendar/connect`는 401, `POST /home/calendar/disconnect`는 303 `/?n=expired`(세션 쿠키 삭제만), `GET /home/calendar/connect`는 405, `GET /home/calendar`·`/home/calendar/x`는 404다. 세션 상태나 캘린더 정보가 노출되지 않는다. 공개 ingress에 새 경로를 더하지 않았다(`/home(/.*)?` 안).
- `Calendar`가 포함된 로그인 회원 화면은 실계정이 필요해 확인하지 않았다.

## gate 재사용·미실행·한계

- 제품 코드 gate는 일정 조회 각 단위의 수락 증거(main 통합 기록: 리뷰 SAR-SCHEDULE-001-CODE-REVIEW, QA SAR-SCHEDULE-001-QA, VISUAL-RECHECK)를 재사용한다. 이번 운영 검사와 구분한다. 이 작업에서 제품 코드·테스트·의존성·migration은 바꾸지 않았다.
- 미실행: 실계정 로그인·캘린더 연결·조회 왕복(SC-09), 로그인 상태 화면, pull/ACK/token 소비, 실제 rollback, Cloud Console·Access·DNS·ingress 변경, 부하·장시간 검사. 이 항목은 coor와 사용자의 후속 과제다.
- 한계: 브라우저 시각 확인은 하지 않았다(UI 코드 변경 없음). rollback 이미지는 파일 존재·hash까지 확인했고 앱 기동은 하지 않았다. 새 relay는 manifest list ID로 보인다. rollback 시 새 `Calendars` 봉인 토큰은 이전 코드가 읽지 않는다.
- Windows `make`가 없어 `lint.py`의 `make lint`·`make test`는 실행 불가로 남는다. 같은 한계는 직전 배포와 같다. 아래 절의 결과를 따른다.

## coor 준비 상태 (SC-09)

- 배포 SHA: `4787e85a0b6952970a861b6b100b3561b9cf30a0`. relay healthy·restart0. 운영 env의 `KNOWSLINK_CALENDAR_KEY` 설정 완료(값 비노출). egress 두 주소 도달.
- 사용자가 먼저 Cloud Console Data Access에 `calendar.freebusy` scope를 등록해야 캘린더 동의 화면이 동작한다(이번 작업은 Cloud Console을 바꾸지 않았다). OAuth 앱은 Testing이고 사용자는 테스트 사용자다. refresh token은 7일 뒤 만료하므로 연결이 주마다 끊길 수 있다(E6).
- 그 뒤 coor가 사용자와 실계정 1경로(연결·승인 1회·해제)를 진행하고 로그에 비밀값·바쁨 원본·이메일이 없음을 확인한다.

## 문서 HEAD의 FullOps 검사

아래 절은 결과 문서 커밋 뒤 채운다.

## 보존 위치

- 로컬 원본: `C:/Users/shin/AppData/Local/KnowsLinkDevTools/release003-ops/`(bundle·스크립트·원 로그).
- 서버 private 작업물: `/home/shin/deploy/knowslink-state/release003-ops-20261021/`(bundle·관측·backup.list·env 보관본).
- 서버 백업: `/home/shin/deploy/knowslink-state/backups/7b4f9bc-pre-release003-20261021T044846Z.dump`(0600).
- 보존 컨테이너: `knowslink-release003-rehearse-20261021`(복원 fixture), `knowslink-release003-rollback-inspect-20261021`(rollback artifact 검사). 이미지 tag `knowslink-rollback:7b4f9bc-release003-20261021`.
- Git 증거: 이 폴더의 redacted JSON/log와 [commands.txt](commands.txt)(실행 스크립트 전문). 원 credential·token·cookie 값·캘린더 키 값·실메시지·DB row는 포함하지 않는다.
