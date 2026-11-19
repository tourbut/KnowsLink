---
title: 수락한 사람 peer·UI의 운영 적용과 기존 계정 보존 LIVE 기록
status: draft
updated: 2026-10-12
owner: ops
tasks: [SAR-HUMAN-PEER-001-OPS]
summary: 26094e9 운영 적용의 backup·격리 복원·migration·보존 집계·공개 smoke·미검증과 복구 방침
---

# 수락한 사람 peer·UI의 운영 적용과 기존 계정 보존 (LIVE 기록)

원문 로그는 같은 폴더의 `*.log`다. 값·원문은 출력하지 않았다. 서버 시계는 UTC다. 서버 접속은 task 전용 helper(`AppData/Local/KnowsLinkDevTools/humanpeer001-ops/remote.py`, 알려진 host key 검증, 자격은 메모리만)로 했다. 공개 smoke는 이 PC에서 `https://link.knowslog.com`을 대상으로 했다.

## 적용 요약

| 항목 | 결과 |
|---|---|
| 이전 배포 | `95df83c22c1559453a7c73504aa248088a8e4ed0`. 시작 시 detached, clean, relay·Postgres healthy, goose 버전 1 |
| 적용 제품 SHA | `26094e90648aa71e8d07b381e9891acb37438730`. 서버 `git rev-parse HEAD` 일치, detached, `status --porcelain` 0줄 |
| 전달 | 서버에 origin 경로가 없다. `git bundle create <파일> main ^95df83c`(SHA256 `7a92999f9a8bb426d82ac3997abcdac1aee43db83bd3b2a2f5a671fb425c6dfe`, 3556773 bytes)를 SFTP 0600으로 전달했다. 서버에서 SHA256 일치·`bundle verify`·`FETCH_HEAD==26094e9`·`95df83c` 조상 관계를 확인했다 |
| 인프라 diff | `95df83c..26094e9`에서 `db/migrations/00002_human_writer.sql` 한 파일만 다르다. compose·deploy·Dockerfile·go.mod·Makefile·sqlc는 0건이다 |
| 쓰기 중단 | `docker compose stop relay` 2026-10-12T12:43:44Z. relay 재기동 healthy 약 12:44:34Z(약 50초). 이 동안 Tunnel 요청이 origin 오류를 받았다 |
| 새 backup | `knowslink-state/backups/95df83c-pre-humanpeer-20261012T124350Z.dump`, 0600, 11098 bytes, SHA256 `76b78824cbd0699e24ff7afd17555fe70b2e2b3cab513465760946021daad206`. 쓰기 중단 뒤 epoch 813716 고정 상태에서 `pg_dump -Fc`를 실행했다 |
| backup 검증 | `pg_restore --list` 종료 0, TOC 23줄. 격리 복원 리허설 통과(아래) |
| migration | `beta.sh up`(종료 0)의 migrate 컨테이너가 `00002`를 적용했다. goose 버전 1→2, `human_writer_guard` trigger 1개 |
| relay | 새 이미지, healthy, restart 0, 5분 error/panic 0건, `google login configured=true`. `verify.py local` 종료 0 |

## 사용하지 않은 것 (지시서 금지·충돌 처리)

`beta.sh`의 `deploy`·`backup`·`restore-verify`·`prepare`·`seed`·`owner-login`과 Tunnel 명령은 실행하지 않았다. `deploy`와 `backup`은 실패 시 `rm`을 포함하고, `deploy`는 서버에 없는 origin을 fetch한다. 대신 같은 순서를 수동 실행했다: 전달→detached checkout→이미지 build→relay 정지→삭제 없는 `pg_dump`(noclobber, 새 고유 파일)→격리 복원→`beta.sh up`. 기존 backup 13개는 변경하지 않았다(14개로 증가). 로컬 Docker는 쓰지 않았다. seed·owner-login·DB 복원·구 writer 복귀·삭제·이동은 없다.

## 격리 복원 리허설 (복사본에서만)

`docker run --rm -i --network none --user postgres`(고유 이름, 같은 pinned `postgres:17-alpine` digest)로 빈 클러스터를 만들어 backup을 복원했다. 컨테이너는 종료 때 자동 소멸했고 잔류 0개다. 첫 시도는 migration SQL 파일 권한(0600) 때문에 Up이 적용되지 않았다. 이때 "DOWNGRADE NOT BLOCKED"는 trigger가 없는 복사본의 결과이며 실패 로그(`backup.log`)로 보존한다. 파일을 0644로 바꾼 재실행(`rehearse.log`)은 다음을 확인했다.

| 검사 | 결과 |
|---|---|
| 복원 | goose 1, relay_state 1행, epoch 813716. 복원한 `Agents`의 md5가 라이브와 같다(`49372e05…`) |
| migration Up 적용 | trigger 1개 |
| 첫 HumanModel=1 쓰기 | 성공 |
| HumanModel 필드를 지우는 쓰기 | `human authorization writer downgrade refused`로 거부 |

## 보존 대조 (집계·md5만, 값 미출력)

`fingerprint-before.log`(12:41:46Z, 쓰기 중단 전)와 `fingerprint-after.log`(12:44:44Z)의 비교다.

| 항목 | 전 → 후 |
|---|---|
| 상태 파일 `.env`·`tunnel.uuid`·`access.aud`·`tunnel/config.yml` SHA256 앞 16자 | 동일 (`de221b94…`·`25ec7513…`·`cd6a969f…`·`f5abe860…`) |
| 권한 | state 0700, `.env` 0600, tunnel 0700, config 0600, backups 0700. 동일 |
| relay_state 키별 개수 | Agents 28, Owners 24, Members 2, Identities 2, Pairs 13, Sessions 1, Remote 6, Messages 9, Idempotency 9 동일. 추가: `HumanModel` 0→1, `Friends` 0개 |
| Owners·Identities·Sessions·Remote md5 | 동일 |
| Agents의 기존 필드(Changed·Credential·Keys·Owner·Revoked) md5 | 동일. `Keys`(키 목록·공개키) md5 `e2c705a9…` 동일 |
| Agents 전체 md5 | 변경(`49372e05…`→`b2be8b01…`). 원인은 모든 28개 agent에 빈 문자열 `Name`·`Client` 필드가 추가된 것이다. 값은 28개 모두 빈 문자열이다. `agentdiff.log`가 backup 복원본과 필드별로 대조했다 |
| Members | 기존 필드 md5 동일. 두 회원 모두 `FriendCode` 생성, `SelfConnect=false`, 이름 미설정, 활성 |
| Pairs | 13개 모두 active·세대1, 기존 필드 md5 동일. pending 0건이어서 만료로 바뀐 pair 없음 |
| 회원별 | 회원1: agent 4(키 있는 agent 3), 회원2: agent 2(키 있는 agent 1). 활성 pair는 회원1 쪽 2개, 회원2 쪽 1개다(`m2.log`) |
| agent 철회 | 0 → 0 |
| 변화한 값 | `epoch`와 `Rates`(relay가 살아 있는 동안의 rate 계수). 이 과제는 DB에 직접 쓰지 않았다 |

`HumanModel=1` 변환은 relay 기동 직후 자동으로 저장됐다(12:44:44Z 시점 이미 1). 변환 코드(`migrateHuman`)는 pending pair 만료·회원별 `SelfConnect=false`·`FriendCode` 생성만 한다. 이번 상태에는 pending이 없어 만료가 0건이다.

## 공개 smoke (이 PC → link.knowslog.com, 무인증)

`public-smoke.log` 12:45:41Z.

| 검사 | 결과 |
|---|---|
| `GET /` | 200, no-store, CSP(`form-action 'self' https://accounts.google.com`, `frame-ancestors 'none'`), no-referrer, nosniff. 새 UI 증거: `--accent:#005A62`·`--canvas:#FFFFFF`·`--selected:#E8F3F3`(이전 `#164da0` 없음), `<h1>친구의 에이전트와 대화하세요`, `<title>로그인 · KnowsLink`, 건너뛰기 링크 |
| 새 경로 무세션 GET | `/home`, `/home/name`, `/home/friends`, `/home/friends/x`, `/home/agents`, `/home/account`, `/home/new`, `/home/conversation`, `/home/legacy`, `/home/receipts` 모두 303 → `/?n=expired` |
| 새 POST 무세션 | `/home/profile`, `/home/agent-name`, `/home/friends/request`, `/home/friends/decision`, `/home/friend-code`, `/home/self-connect`, `/home/conversation/send` 모두 303 → `/?n=expired`. 상태 변경 없음 |
| public deny | `/healthz`, `/owner`, `/v1/test/pull`, `/_probe`, `/v1/owners`, `/v1/authorize` 404 |
| 무인증 API·SSE | `/v1/text/events`, `?token=x` 401 JSON. `POST /mcp`(initialize) 401, 빈 본문 400 |
| OAuth | well-known 두 경로 200 JSON. 합성 43자 `/remote/consent/…` 401 HTML. `GET /v1/connect/start` 405 |
| 공유 서비스 | orca 200 / s8 200 / mcp 401 (기준과 동일) |

## 다른 서비스 영향

`verify.py regression`은 종료 1이다. 원인은 이전 LIVE와 같은 baseline 결함이다. 저장된 baseline의 `projects`에 `knowslink`가 있고, 현재는 그것을 걸러 내면 `jevpet-demo, myportfolio`다. `codes`(200·200·401)·`host_cloudflared_pids`(506937)·myportfolio 컨테이너 6개 상태는 baseline과 같다. 스크립트는 수정하지 않았고 이 실패를 PASS로 바꾸지 않았다. 추가 직접 대조(`after.log`)도 같다: myportfolio·jevpet-demo 컨테이너 목록 md5 동일, compose project 3개 running, 호스트 cloudflared PID 506937 유지, 5432·8080 listener는 127.0.0.1만, KnowsLink cloudflared 컨테이너 재시작 없음(Up 33시간 유지). KnowsLink cloudflared 로그의 origin 오류 12건은 모두 relay 정지 구간(12:43:53~12:44:15Z)에 있고 이후 2분은 0건이다. `verify.py public`은 이전 Access 302 전제의 검사여서 현재 공개 ingress에 맞지 않아 실행하지 않았다.

## 미검증 (자동 PASS로 기록하지 않음)

- 실제 Google 계정의 새 `/home`(내 이름 입력→대화)·모바일·키보드 포커스·200% 확대는 실행하지 않았다. 담당은 coor다.
- 인증된 SSE·실제 메시지 전송·친구 요청·자기 연결 동의·legacy exact pair 전송은 실행하지 않았다. 기존 키로 서명한 전송을 만들지 않았다. 실제 계정의 가림 없는 legacy 전송 가능 여부는 coor 후속이다.
- 새 relay 이미지로 backup 복사본에 대한 변환 리허설은 하지 않았다. 변환은 라이브에서 한 번 실행됐고 위 보존 대조가 그 결과다.
- 독립 QA·리뷰의 low 3건 공백, 실제 keyboard/200% 검사, scenario not_run은 제품 수락 기록 그대로다.
- 서버 `make test`와 제품 전체 QA는 반복하지 않았다. 고정 SHA의 tester·visual 증거를 재사용한다.

## 복구 방침

코드만 이전으로 되돌리지 않는다. 이미 `HumanModel=1`이고 trigger가 구 writer를 거부한다. 문제가 생기면 relay를 정지하고 현재 DB를 새 backup한다. 현재 schema·친구·차단·세대·철회·키를 보존하는 검토된 reader/writer로 전진 복구한다. 위 backup(`95df83c-pre-humanpeer-20261012T124350Z.dump`)은 복원 대상이 아니다. 복원하면 이후의 이름·친구·철회 상태를 잃고 권한이 되살아날 수 있다. 권한이 불명확하면 전달을 차단하고 재인증·재연결·새 수락으로 복구한다. 이 절차는 coor의 별도 결정이 필요하다.
