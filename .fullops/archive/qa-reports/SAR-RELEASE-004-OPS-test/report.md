---
title: "main 40aafe6f Vultr 운영 적용(migration 3·4, Redis 기동) 결과"
status: draft
updated: 2026-11-14
owner: ops
tasks: [SAR-RELEASE-004-OPS]
summary: "Vultr knowslink-relay-01을 c6dfc39a(goose 2)에서 40aafe6f(goose 4, Redis)로 D12 수동 절차로 옮긴 순서·시각·백업·격리 리허설·검증·실측 크기를 기록한다."
---

# SAR-RELEASE-004-OPS 결과

운영 제품을 `c6dfc39ac4f6dabddc5864ad611437b84ded383e`에서 `40aafe6ff396db2681a9423c3b344732934756e0`로 옮겼다. migration 3(`state` json·lz4, `human_model`)과 migration 4(메시지·멱등·Dots Inbox 삭제)를 적용했고 Redis를 처음 기동했다. 본 전환의 쓰기 중단은 약 28초다. 그 전에 도구 결함으로 1차 시도를 백업 전에 중단했고, 원래 relay를 다시 띄웠다(약 10초 중단, DB 변경 없음). 이행 전·후 dump가 각각 있다. 실패·rollback은 없다.

## 범위·기준

- 사용자 지시 "운영서버에 적용해"(결정 D056). Task `task_be315ae555b4`, Dispatch `ctx_faeb043e4f1e`. 지시서 `handovers/to_ops.md`(기준 `40aafe6f`, 커밋 `9a9e8467`). 하위 에이전트는 쓰지 않았다(직렬 운영).
- 절차: [D12](../../../operations/ops-guide.md)의 "SAR-VULTR-MIGRATE-001", "SAR-DB-WAIT-001 저장 형식 이행", "SAR-MSG-EPHEMERAL-001 메시지 휘발 저장 이행" 절. `beta.sh deploy`는 migration 변경을 거부하므로 쓰지 않았다.
- 바꾼 것: 서버 checkout `c6dfc39a`→`40aafe6f`(detached), relay·migrate 이미지 재빌드, Redis 컨테이너 생성, 운영 DB migration 3·4, 새 백업 2개, rollback 태그 2개, 서버 `deploy-history.log` 한 줄, 작업 폴더 `knowslink-state/release004-20261113/`(대상 소스 사본·스크립트·리허설 출력, 비밀값 없음).
- 바꾸지 않은 것: 운영 `.env`(읽거나 출력하지 않음), Tunnel 자격·`config.yml`·ingress, Cloudflare·DNS·Vultr 방화벽, 홈서버, 기존 백업·볼륨·이미지, 다른 project. 삭제·prune·`down`은 없다. 서버에서 지운 것은 이 과제가 만든 이미지 태그 `kl-release004:40aafe6f` 하나다.
- 시각은 서버 UTC다(KST = UTC+9). 서버 UTC 2026-11-13T15:39Z는 KST 2026-11-14 00:39다.
- 실계정 로그인·메시지 발송은 하지 않았다(tester 후속).
- 원시 로그: `D:/workspace/KnowsLink/.devtools/release-004/`(`01-precheck.log`~`10b-cloudflared-err.log`, 스크립트 `rehearse.sh`·`cutover.sh`·`post.sh`·`late.sh`·`smoke.sh`).

## SSH host key

이 PC의 `known_hosts`에는 `64.176.237.173` 항목이 없었다(다른 IP 항목만 있음). 첫 접속에서 `StrictHostKeyChecking=accept-new`로 등록했다. 등록한 키는 ED25519 `SHA256:3InoMmjs4Pcs7TEkhvW98vWAiQLh/FfQNzpGE+O5xSs`다. 이 지문을 Mac 관리 세션의 `known_hosts`나 Vultr 콘솔 값과 대조하지 않았다. 후속으로 한 번 대조한다. 키 인증은 서명이 세션에 묶이므로 가짜 서버가 있었다면 실제 서버에 접속할 수 없었다. 실제 서버 상태(아래 사전 확인)는 D12 기록과 같다.

## 사전 확인 (읽기 전용, 15:31:56Z)

- checkout `c6dfc39a`, 깨끗함. 컨테이너 relay·postgres healthy, cloudflared Up(27시간), migrate Exited 0.
- 디스크 `/` 52G 중 12G 사용(24%). 메모리 1962 MB 중 available 1400 MB, swap 5399 MB.
- goose 2. `relay_state` 1행. `pg_column_size(data)` 10746 bytes, 텍스트 21622 bytes. `pg_total_relation_size` 901120 bytes, DB 8722099 bytes.
- 메시지 부분: `Messages` 2 bytes(`{}`), `Idempotency` 2 bytes(`{}`), 설치 Inbox 1곳 2 bytes(빈 값). 이전 버전은 메시지를 24시간 보관했으므로 지난 24시간 동안 남은 메시지 기록이 0이다.
- 백업 1개(`c6dfc39a-pre-vultr-20261112T124516Z.dump`), prune cron 1줄.

## 무중단 준비 (02-prepare.log)

- `git fetch` 뒤 `40aafe6f`가 `origin/main`에 있고 `c6dfc39a`가 그 조상임을 확인했다.
- `git archive`로 대상 소스를 작업 폴더에 풀었다. 대상 `beta.sh`를 같은 폴더에 복사했다(이행 전 백업의 `message_state_rows` 출력에 필요).
- rollback 태그: `knowslink-rollback:c6dfc39a-relay-release004-20261113`(`sha256:ac0be3e45268…f468a`, 실행 중 relay와 같은 ID), `knowslink-rollback:c6dfc39a-migrate-release004-20261113`(`sha256:e75442331ada…9e37`). migration 3 뒤에는 코드 전용 복귀가 불가능하므로 기록용이다.
- Redis 이미지를 digest로 pull했다. 리허설용 이미지 `kl-release004:40aafe6f`를 미리 빌드했다(171초). 본 전환의 compose 빌드는 같은 캐시를 써서 9초에 끝났다.
- 리허설 스크립트를 기존 11-12 dump로 먼저 시험했다(운영 DB 미접촉). 첫 실행은 검사 SQL의 형 변환 오류로 멈췄고, 고친 뒤 모든 검사가 통과했다(03-rehearse-dryrun-old-dump.log).
- 이행 전 공개 smoke: 아래 표와 같고 `/privacy`에 새 문구 0건(04-smoke-before.log).

## 1차 시도 중단 (05-cutover.log)

| 시각(UTC) | 내용 |
|---|---|
| 15:38:44.465 | cloudflared·relay 정지 시작 |
| 15:38:45.226 | 정지 완료 |
| 15:38:45 | 복사한 `beta.sh`가 git 밖이라 `SRC` 계산(`git rev-parse`)이 실패. dump 파일을 만들기 전 |
| 15:38:54.600 | 스크립트의 실패 분기가 원래 relay·cloudflared를 다시 시작. migrate(`c6dfc39a`)가 다시 돌았지만 goose 2 그대로 |

DB·checkout은 바뀌지 않았다(05b-after-abort.log: goose 2, `c6dfc39a`, 부분 dump 없음). `KNOWSLINK_SRC`를 지정해 고친 뒤 다시 실행했다.

## 본 전환 (06-cutover.log)

| 단계 | 시각(UTC) | 결과 |
|---|---|---|
| 쓰기 중단 시작 | 15:39:19.536 | `cloudflared`·`relay`만 `stop`. Postgres 유지 |
| 정지 완료 | 15:39:20.213 | |
| 정지 DB 측정 | 15:39:20 | goose 2, 1행, 10773 bytes, 텍스트 21851 bytes, `md5(data::text)` `19a9edce353671f5486b24fbf01de03f`. 메시지 2·멱등 2·Inbox 2 bytes |
| 이행 전 백업 | 15:39:21 | 아래 "백업" 표 |
| restore-verify | 15:39:2x | `tables=2 relay_state_rows=1 state_type=jsonb human_model=none` |
| 격리 리허설 | 15:39:29 | 아래 "격리 리허설" |
| 리허설 통과 | 15:39:31.058 | |
| checkout | 15:39:31 | detached `40aafe6f`, `deploy-history.log`에 한 줄 |
| `beta.sh up` | 15:39:31~15:39:40.983 | 빌드 전 단계 캐시. Redis healthy 뒤 migrate가 15:39:33에 "SQL migrations completed", relay healthy. postgres는 재생성하지 않음 |
| 검증 | 15:39:41~ | 아래 "검증" |
| 트래픽 재개 | 15:39:47.250 | `cloudflared` 시작. Tunnel 연결 4개 등록(icn06 2, icn05 1, icn01 1) |

- 쓰기 중단: 15:39:19.536Z~15:39:47.250Z(2026-11-14 00:39:19~00:39:47 KST), 약 28초.
- 1차 시도까지 합한 공개 중단: 약 38초.
- `dc start cloudflared`가 의존성으로 migrate를 한 번 더 돌렸다. 이미 goose 4라 바뀐 것은 없다.

## 백업

| dump | 모드 | bytes | TOC | SHA256 | relay_state_rows | message_state_rows | restore-verify |
|---|---|---|---|---|---|---|---|
| `c6dfc39a-20261113T153921Z.dump`(이행 전) | 0600 | 11411 | 25 | `eef7815f179c0d19d281580fa7b7362e75ec1c4673e1cf31ffc375bd92a32e03` | 1 | 1 | `state_type=jsonb human_model=none` |
| `40aafe6f-20261113T154024Z.dump`(이행 직후) | 0600 | 11454 | 25 | `2f0f8a5b8f801cfdcfd6827d5becc385c951b948a88d08fe65ce93d2bd3974a3` | 1 | 0 | `state_type=json human_model=1` |

- 이행 전 dump의 `message_state_rows=1`은 `Messages`·`Idempotency` 키가 있다는 뜻이다. 두 값은 빈 객체여서 메시지 내용은 들어 있지 않다.
- 이행 직후 dump는 첫 실행 스크립트가 stdin 소비로 조기 종료한 뒤 `beta.sh backup`을 단독으로 실행해 만들었다(08-post.log, 08b-backup-debug.log). 원인은 `bash -s` heredoc을 `docker compose exec`가 stdin으로 읽은 것이다. 남은 부분 파일은 없다.
- 이제 가장 최근 dump가 이행 직후 dump이므로 이행 전 dump 2개는 14일 정리 대상이다. `prune-backups.sh --list`는 현재 대상 0개다(14일 미경과).

## 격리 리허설

이행 전 dump를 network-none Postgres(`kl-release004-rehearsal`)에 `--role=knowslink`로 복원했다. 대상 이미지의 `/app/migrate`를 같은 network namespace(`--network container:`)에서 실행해 loopback으로만 접속했다. 운영 DB에는 닿지 않는다. 결과(06-cutover.log):

- 전: goose 2, `jsonb`. 개수: Agents 31, Owners 24, Pairs 13, Remote 6, Members 2, Identities 2, Friends 1, HumanModel 1, Rates 3, 나머지 0. Inbox가 있는 설치 1곳.
- 후: goose 4, `state_type=json`, `attcompression=l`, `pg_column_compression=lz4`, `human_model=1`(전 `HumanModel` 1), guard trigger 1, 메시지 키 행 0, Inbox 있는 설치 0.
- 보존: `Messages`·`Idempotency`·`Remote`를 뺀 14개 키가 `jsonb`로 모두 같다. 키 집합은 두 키만 빠졌다. `Remote`는 각 설치의 `Inbox`를 뺀 값과 같다.
- `human_model`을 낮추는 UPDATE가 거부됐다.
- 리허설 컨테이너는 스크립트 종료 때 지웠다.

## 검증

- 상태: postgres·redis·relay healthy, restart 0. migrate exit 0. 운영 checkout `40aafe6f`, 깨끗함.
- DB: goose 4, guard trigger 1, `state` json, 압축 `l`/`lz4`, `human_model=1`, 메시지 키 행 0. `pg_column_size(state)` 8894 bytes, 텍스트 19971 bytes, `pg_total_relation_size` 155648 bytes(전 901120). `md5(state::jsonb::text)` `fe3abd0b9c02b929d67f415785779567`.
- 전환 2분 뒤 live 개수는 리허설 전 개수와 같다. Rates만 3→0이다(휘발성 값). `VACUUM FULL`은 실행하지 않았다(선택 단계, 표가 이미 작다).
- Redis: `CONFIG GET save`=빈 값, `appendonly`=`no`, `maxmemory-policy`=`noeviction`, `maxmemory`=134217728. `used_memory` 1497784 bytes(1.43 MiB), DBSIZE 0, 컨테이너 메모리 30.8 MiB/192 MiB. `docker port` 결과 없음(내부 `database` network의 6379/tcp만).
- `verify.py local` 통과: health 200, unknown 404, owner/api 401, loopback-only relay, private Postgres, limits, 0600/0700.
- `verify.py public` 통과: `/healthz`·`/owner`·`/v1/registry`와 Access 헤더 두 종류 모두 404.
- relay 로그의 panic·error·unavailable 0줄. cloudflared ERR 10줄은 두 번의 의도한 정지 시각(15:38:44Z, 15:39:19Z)의 연결 종료 로그다.
- 서버 메모리 available 1473 MB, 디스크 26%.

공개 smoke(이 PC, 요청 간격 1초, 07-smoke-after.log):

| 요청 | 이행 전 | 이행 후 |
|---|---|---|
| `GET /`, `/privacy`, `/terms` | 200 | 200 |
| `GET /.well-known/oauth-authorization-server`, `/.well-known/oauth-protected-resource`, `/.well-known/oauth-protected-resource/mcp` | 200 | 200 |
| `GET /install/0.1.1/knowslink-setup.ps1` | 200 | 200 |
| `GET /healthz`, `/owner`, `/_probe` | 404 | 404 |
| `POST /mcp` initialize(무인증) | 401 | 401 |
| `/privacy`의 "서버 메모리(Redis)" 문구 | 0 | 1 |

## Redis 한도 여유 (리뷰 medium 2)

- 실측 메시지 부분: 이행 직전 `Messages`·`Idempotency`·Inbox 각 2 bytes(빈 값). 이전 버전의 24시간 보관 창에서 남은 메시지가 0이다. 이행 뒤 Redis 키 0개, `used_memory` 1.43 MiB다.
- 여유: `maxmemory` 128 MiB에서 기본 사용 1.43 MiB를 빼면 약 126.6 MiB다. commit 중 사본 2개를 쥐므로 메시지 기록 전체는 약 63 MiB까지 담을 수 있다. 현재 사용은 0 bytes라 여유는 사실상 100%다.
- 판단: 현재 베타 규모에서는 한도 도달 가능성이 낮다. 축소 쓰기가 막히는 결함은 그대로이므로 dev 후속(사본 크기 상한 또는 메시지별 키)은 유지한다. 메시지가 없어 행 lock 안 전체 사본 처리 시간은 측정하지 못했다.
- 감시: `docker exec knowslink-redis-1 redis-cli INFO memory | grep '^used_memory:'`. 64 MiB에 다가가면 dev 후속을 앞당긴다.

## 복귀 경계

- 이전 SHA는 `data` 컬럼이 없어 SQL 오류로 멈추고 메시지를 다시 Postgres에 쓴다(D051 위반). Down은 실패 고정이다. 기본은 전진 수정이다.
- 이행 전 dump로 되돌리는 복원은 이행 뒤 쓰기를 잃으므로 사용자 승인 대상이다. 실제 rollback은 실행하지 않았다.
- Redis를 멈추면 relay는 `unavailable`, `/healthz` 503이다.

## 정리·보존

- 이 PC에는 Docker 자원을 만들지 않았다.
- 서버: 리허설·restore-verify 컨테이너 0개 남음. 과제 이미지 태그 `kl-release004:40aafe6f`를 지웠다(레이어는 운영 이미지와 공유). 이전 이미지(rollback 태그), 볼륨, 백업 3개는 그대로다.
- 작업 폴더 `knowslink-state/release004-20261113/`는 기록으로 남겼다(0700, 비밀값 없음).

## 남은 일

- tester: 실계정 QA(로그인·메시지 왕복·Redis 경로, 새 개인정보 문구 렌더링).
- dev: 리뷰 medium 2 축소 경로 설계.
- 사용자: 이전 백업 보존 판단(홈서버 `c6dfc39a-pre-vultr-…dump`와 그 이전 dump), 개인정보처리방침 시행일·고지.
- ops: 이 PC에 등록한 host key 지문을 Mac `known_hosts` 또는 Vultr 콘솔과 한 번 대조한다.
