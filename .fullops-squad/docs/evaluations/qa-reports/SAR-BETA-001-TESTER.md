---
title: SAR-BETA-001-TESTER — 베타 로컬 런타임 좁은 독립 QA 보고서
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-BETA-001-TESTER]
summary: "실행 배포 f824015의 loopback, 인증 음성, 백업, expose 차단과 공개 held를 기록한다"
---

# SAR-BETA-001-TESTER — 베타 로컬 런타임 좁은 독립 QA 보고서

## 판정

판정 실행의 배포 체크아웃은 `f824015314c66bcab42940cfe3db2edabb22e1dd`다. 시작 HEAD와 끝 HEAD가 같다.
이 커밋의 부모는 `8a7ad363b03851f1fea08b3fcfc06ab135bd38a4`다. `437f1432a158670a485413c1aba159debc3759e5`는 조상이다.
제품 경로와 `78b1d92c8aa626245d3349ffaf7367d28f1dd3ef`의 diff는 비어 있다.
loopback 8080, Postgres 비게시, 비밀 파일 0600, 재시작과 자원 제한, 합성 인증 음성, 백업, 격리 복원, Access 미기록 때 expose 차단, 공유 서비스 회귀는 통과했다.
이 판정은 공개 보호 연결의 통과가 아니다. 공개 검사는 held다. 제품 코드와 배포 설정은 수정하지 않았다.
시나리오: [SAR-BETA-001-TESTER.md](../scenarios/SAR-BETA-001-TESTER.md). 로그: [SAR-BETA-001-TESTER-test/](SAR-BETA-001-TESTER-test/).

## 배포 SHA

처음 읽은 체크아웃은 `437f1432a158670a485413c1aba159debc3759e5`였다. 작업 트리는 깨끗했다. `deploy/knowslink` 다섯 파일의 blob은 그 커밋과 같았다.
이 QA는 `git checkout`을 실행하지 않았다. 배포 reflog가 체크아웃을 바꿨다. 기록은 [reflog.log](SAR-BETA-001-TESTER-test/reflog.log)다.
19:56:04에 체크아웃이 `437f143`과 `8a7ad36` 사이를 이동했다. 19:56:23에 `437f143`에서 `8a7ad36`으로 이동했다. 19:57:02에 `8a7ad36`에서 `f824015`로 이동했다.
`f824015`와 `8a7ad36`의 차이는 OPS 문서 2개다. `deploy/knowslink` blob은 같다.
`8a7ad36`과 `437f143`에서 `beta.sh`와 `access_apply.py`가 다르다. `verify.py`, `compose.ops.yaml`, tunnel 템플릿은 같다.
backup, restore-verify, expose, selftest는 체크아웃 스크립트를 실행했다. 그 바이트는 `8a7ad36`과 같다.
relay 컨테이너는 `8a7ad36`으로 옮긴 시각에 만들어졌다. 문서 커밋 `f824015`는 relay를 다시 만들지 않았다.
HEAD가 `8a7ad36`인 동안의 첫 실행은 [during-move/](SAR-BETA-001-TESTER-test/during-move/)다. 판정은 그 다음 실행이다. `run.py` 종료코드는 0이다.

## 환경

날짜는 2026-10-03이다. 브랜치는 `fullops/tester`다. 검사 시작 워크트리 HEAD는 `97235aaaa9a4b0d761082d58db14573eb718c750`다.
공통 기준은 `fullops-common-0.3.2`다. Go 1.27.1, Node v22.22.2, Python 3.12.3, Docker 29.4.3, Compose v5.1.3이다. 기록은 [versions.log](SAR-BETA-001-TESTER-test/versions.log)다.
credential, CSRF, API 토큰, 사용자 이메일은 로그에 쓰지 않았다. `owner-login`과 `seed`는 실행하지 않았다.
Jev find/context는 `SAR-BETA-001-OPS-context.json`을 재사용했다. ops-guide와 contexts/ops의 예전 배포 금지는 지시서가 역사 기록으로 둔 내용이다.

## 실행한 경계

| 검사 | 결과 | 관찰 |
|---|---|---|
| deploy-lineage | 통과 | `f824015`. `437f143`의 자손. 같은 커밋은 아님 |
| product-unchanged | 통과 | 제품 diff 비어 있음 |
| loopback-8080 | 통과 | relay `127.0.0.1:8080`. Postgres `PortBindings` 비어 있음 |
| ss-loopback | 통과 | `127.0.0.1:8080`만 있음 |
| restart-limits | 통과 | postgres·relay `unless-stopped`. migrate `no`. 메모리 512m·256m·128m |
| log-limit | 통과 | json-file 10m, 파일 3개 |
| secret-mode | 통과 | `.env` 0600. tunnel 자격 json 1개 0600. 상태 디렉터리 0700 |
| verify-local | 통과 | 종료코드 0. health 200. 미지 경로 404. owner·API 401 |
| verify-regression | 통과 | 실행 전후 종료코드 0. 200/200/401. 호스트 cloudflared PID 506937 |
| access-selftest | 통과 | 종료코드 0. `access.aud` 없음 |
| NEG-no-auth | 통과 | `GET /owner` 401 `invalid_auth` |
| NEG-agent-owner | 통과 | agent credential의 `GET /owner` 401 `invalid_auth` |
| NEG-agent-revoke | 통과 | agent credential의 `POST /v1/owner-revoke` 401 `invalid_auth`. 이후 owner `GET /owner` 200 |
| NEG-bad-csrf | 통과 | owner credential과 `csrf=bad`는 403 `invalid_csrf`. gate 상태 불변 |
| backup | 통과 | `f824015-20261003T110011Z.dump`. 0600. toc 23행. `relay_state_rows=1` |
| prior-backup-kept | 통과 | 기존 dump를 지우지 않음 |
| restore-isolated | 통과 | 종료코드 0. `tables=2 relay_state_rows=1`. 임시 컨테이너 제거. live health 200. 볼륨 유지 |
| expose-blocked | 통과 | 종료코드 1. `Access app not recorded`. DNS 비어 있음. knowslink cloudflared 없음 |
| deploy-stable | 통과 | 판정 실행 동안 HEAD 불변 |
| deploy-untouched | 통과 | 체크아웃 clean. `.env` 크기 불변. baseline 해시 불변 |

relay는 `cap_drop ALL`, `read_only`, `no-new-privileges`다. postgres 호스트 포트 매핑은 없다. 값은 [limits.json](SAR-BETA-001-TESTER-test/limits.json)이다.
로컬 health 200은 공개 수락이 아니다. 공개 negative의 통과 조건은 Access가 미인증을 막는 것이다. `verify.py public`은 실행하지 않았다.

## 명령과 종료코드

| 명령 | 종료코드 | 로그 |
|---|---|---|
| `run.py` 판정 실행 | 0 | [results.json](SAR-BETA-001-TESTER-test/results.json) |
| `verify.py local` | 0 | [verify-local.log](SAR-BETA-001-TESTER-test/verify-local.log) |
| `verify.py regression` 전 | 0 | [verify-regression-before.log](SAR-BETA-001-TESTER-test/verify-regression-before.log) |
| `verify.py regression` 후 | 0 | [verify-regression-after.log](SAR-BETA-001-TESTER-test/verify-regression-after.log) |
| `access_apply.py selftest` | 0 | [access-selftest.log](SAR-BETA-001-TESTER-test/access-selftest.log) |
| `beta.sh backup` | 0 | [backup.log](SAR-BETA-001-TESTER-test/backup.log) |
| `beta.sh restore-verify` | 0 | [restore.log](SAR-BETA-001-TESTER-test/restore.log) |
| `beta.sh expose` | 1 | [expose.log](SAR-BETA-001-TESTER-test/expose.log). 기대된 차단 |

첫 실행의 `run.py` 종료코드는 1이다. HEAD가 이미 `8a7ad36`이라 고정 SHA 비교와 dump 이름 비교가 실패했다. revoke 응답의 JSON `invalid_auth`를 본문 문자열과 다르게 비교했다. owner 화면은 200으로 남았다. 제품 결함으로 기록하지 않는다.
`make test`, `make verify-mvp`, `make verify-runtime`은 실행하지 않았다. 제품 트리가 `78b1d92`와 같기 때문이다.

## 재사용

`78b1d92` 제품 QA는 [SAR-MVP-001-TESTER-FINAL.md](SAR-MVP-001-TESTER-FINAL.md)다. 보고 커밋은 `659f4b06dfda1bd59997f8f2b06cba23b48026a6`다. 그 SHA와 조건으로 재사용한다.
최종 리뷰 `311381f`와 직접 UI `e238777`도 원래 SHA와 조건으로 재사용한다. 새 캡처는 없다.
main `557ebc3`은 그 기록의 기준이다. 이번 실행의 통과로 바꾸지 않는다.
`4262d02`와 `a6a10c7`의 QA, pass, fail, held도 그 파일에 남긴다.

## held

공개 DNS, Tunnel connector, Access 앱은 없다. 공개 negative와 본인 이메일 로그인은 held다. 인간 검사와 구분한다.
DEC-02 실데이터, DEC-03 공개 한도와 신원, 실adapter, 실제 벤더, WAL 또는 backup 삭제, 공개 정책은 held다. 이번 통과로 옮기지 않는다.
`tunnel.uuid` 모드는 0664다. 자격 json과 `.env`는 0600이다. 이 파일은 비밀 생성값이 아니다. high로 올리지 않는다.

## 추가 게이트

coor가 실행 배포 SHA를 `f824015314c66bcab42940cfe3db2edabb22e1dd`로 확인했다. 이 절은 그 확인 뒤에 임시 상태와 임시 clone에서만 돌렸다. 라이브 체크아웃은 옮기지 않았다.
`gates.py` 종료코드는 0이다. 기록은 [gates.json](SAR-BETA-001-TESTER-test/gates.json)과 [gates.log](SAR-BETA-001-TESTER-test/gates.log)다.

| 검사 | 결과 | 관찰 |
|---|---|---|
| access-selftest | 통과 | 종료코드 0. 깨진 app·tunnel 설정은 거절 |
| missing-proof | 통과 | `access.json`이 없으면 `check`가 끝난다 |
| stale-proof | 통과 | 10분이 지난 snapshot은 `SystemExit` |
| mismatched-aud | 통과 | 기록된 aud와 app aud가 다르면 거절 |
| foreign-email-rejected | 통과 | 다른 이메일 정책은 거절. 이메일 값은 로그에 없음 |
| extra-idp | 통과 | IdP가 둘이면 거절 |
| expose-missing-aud | 통과 | 종료코드 1. DNS 이전 |
| expose-missing-jwt | 통과 | `required: true`가 없으면 종료코드 1 |
| expose-stale-before-dns | 통과 | 오래된 증명은 `tunnel route` 전에 종료코드 1 |
| deploy-unknown-sha | 통과 | 없는 커밋은 종료코드 1 |
| deploy-migration-blocked | 통과 | migration이 다르면 backup과 checkout 전에 종료코드 1 |

임시 clone의 migration 커밋은 삭제했다. dump는 만들지 않았다. 라이브 HEAD, relay ID, DNS는 검사 전과 같다.
라이브 rollback은 다시 실행하지 않았다. reflog의 `437f143`과 `8a7ad36` 이동은 OPS 기록이다.
관리 Access 권한과 재리뷰는 대기다. 외부 DNS와 connector는 없다.

## 후속

coor가 설정 리뷰와 제품 수락을 이어간다. 공개 좁은 QA는 보호 연결 뒤의 배포 SHA, URL, Access ID가 있을 때 시작한다.
이 경계에서 새 high는 없다. 상태 디렉터리에 이번 QA의 dump 2개가 추가되었다. 기존 dump와 볼륨은 지우지 않았다.
