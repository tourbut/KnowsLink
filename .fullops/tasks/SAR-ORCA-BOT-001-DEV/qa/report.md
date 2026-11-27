---
title: SAR-ORCA-BOT-001-DEV 검증 보고
status: draft
updated: 2026-11-27
owner: dev
tasks: [SAR-ORCA-BOT-001-DEV]
summary: Orca 봇 사전 승인(grant)·adapter의 표적 테스트와 격리 relay 왕복(S1–S6) 결과
---

# SAR-ORCA-BOT-001-DEV 검증 보고

테스트 레벨 lite. 상위 의도는 [SAR-ORCA-BOT-001](../../../docs/planning/intents/SAR-ORCA-BOT-001.md)이다(커밋 e112c129에 있음).

## 표적 테스트

| 검사 | 명령 | 결과 |
|---|---|---|
| relay 단위 | 레포 루트 `go test ./cmd/... ./internal/...` (Windows, make·race 없음) | 통과. 새 테스트 `TestGrantAdmitsOnlyLiveOwnerForm`·`TestGrantEndsOnRevokeAndNeverRevives`·`TestGrantWireKeepsOldSignatures` 포함 |
| i18n | 같은 실행의 `TestTranslationCoversTemplates` 등 | 새 화면 문구 11개 영어 포함 통과 |
| adapter | Node 22.22.2, `npm run build` 뒤 `node dist/{orca,inbox,grok-wake,events,mcp,setup-connect,owner-text}.test.js` | 모두 0 |
| adapter 정적 | `npm run check`(prettier·eslint·tsc) | 통과 |

relay 단위 테스트가 처음 실행에서 실제 결함 하나를 찾았다. 두 agent의 kid 이름이 같으면(`key1`) 받는 쪽 키 철회를 보내는 쪽 키로 잘못 보아 grant가 남았다. 키 목록을 보내는 쪽·받는 쪽으로 나눠 고쳤다.

## S6 격리 relay 왕복

환경: Docker `kl-orca-bot-pg`(postgres:17-alpine, 고정 digest)·`kl-orca-bot-redis`(redis:8.10.2-alpine, 고정 digest), goose 마이그레이션, opt-in fixture [orca_roundtrip_fixture_test.go](../../../../internal/relay/orca_roundtrip_fixture_test.go)(합성 email 회원 1명, 자기 agent 2개, 실제 relay handler). adapter는 `adapters/dist/plugin.js orca …`를 Node 22로 실행했다. 도어벨은 argv만 기록하는 가짜 orca([fake-orca.go.txt](fake-orca.go.txt))로 받았다. 실제 Orca 봇에는 메일을 보내지 않았다. 실행 스크립트 [s6.py](s6.py), 결과 원문 [s6-result.json](s6-result.json).

| 기준 | 확인 | 결과 |
|---|---|---|
| S1 | grant 없음 → `orca send` | 403 `approval_required`, 종료 1 |
| S2 | agent Bearer + `confirmed=true`, 세션 없이 `POST /home/grant-create` | 303 `/?n=expired`(세션 없음). grant 생기지 않음 |
| 웹 | 회원 세션으로 grant 3개(task·answer: Grok→Orca, decision: Orca→Grok) | 303 `/home/agents#grants`, `GET /home/agents` 200에 행 3개 |
| S6 task | Grok `send orca.task.v1` → watcher 도어벨 → Orca `receive` | 도어벨 1통, `receive`에 `kind`·body·`untrusted: true` |
| S3 | 도어벨 argv에 본문·발신 agent ID 포함 여부 | 없음. ID·대기 수·receive 명령·고정 안내뿐 |
| S6 decision→answer | Orca `send orca.decision.v1`(item 3) → Grok `receive` → Grok `send orca.answer.v1` | Grok이 question·options를 받음. answer 도어벨에 `--item 3` |
| S4 | answer 도어벨 | `--item 3` 있음 |
| S3 재시도 | 가짜 orca를 실패로 둔 뒤 task 2 → 실패 기록 → 성공으로 바꿈 | 다음 polling에 같은 ID로 다시 울림. 성공 뒤 30초 동안 다시 울리지 않음(기록 4줄 유지) |
| S1 방향 | Orca→Grok `orca.task.v1`(grant 없는 방향) | 403 `approval_required` |
| 끄기 | 웹에서 grant 3개 끄기 → 다시 task 발송 | 303 뒤 403 `approval_required` |
| S5 | watcher stdout, relay 상태 행, Redis 값의 본문 검색 | stdout·Postgres `relay_state`·Redis 모두 본문 문자열 없음 |

## Linux clean clone 제품 검사

Windows `fullops lint`는 ERROR 0, WARNING 14, 실행 불가 2(`make lint`·`make test`, Windows에 make 없음)로 "미통과"다. 같은 HEAD 441d1a41을 `.devtools/SAR-ORCA-BOT-001-DEV/clone`에 clone해 Docker `node:22`(Node 22.23.5, npm 10.9.9, Debian python3·bash·make)에 Go 1.27.1 linux-amd64를 받아 실행했다.

| 명령 | 결과 |
|---|---|
| `make install` | 0 |
| `make test` | 0. `go test -race`: `ok cmd/migrate`, `ok cmd/relay`, `ok internal/config`, `ok internal/relay`. `npm test`: 모든 PASS 줄, 마지막 `PASS: orca schema, signing bytes without kind unchanged, doorbell metadata only, ring once, retry after failure, --item answer, untrusted receive`. shell_guard·autosetup 테스트 포함 |
| `make lint` | lint-go(gofmt·vet·fixture 태그 vet·`go mod verify`)와 lint-adapters(`npm run check`) 통과. lint-config의 `check_compose.py`에서 멈춤: 컨테이너에 `docker` 실행 파일이 없음(`FileNotFoundError: 'docker'`) |
| lint-config 나머지 | 같은 컨테이너에서 prettier(compose·sqlc), `check_public_ingress.py`, `test_public_config.sh`, `sh -n install_bot_mcp.sh`, `compileall` 통과(`set -e`). `test_public_config.sh`의 ERROR 출력 줄은 그 스크립트의 음성 검사다 |
| `python scripts/check_compose.py` | clone에서 Windows host docker로 실행해 0: `Compose: five services, one-shot migrate, private Postgres, volatile private Redis, loopback relay, optional Tunnel, synthetic signup closed, OPS override merges`. docker compose config가 필요해 host에서 돌렸다 |

`make lint`의 fixture vet 태그에 `orcafixture`를 넣었다(441d1a41).

## 못 한 것

- 실제 Orca 봇(쿠어) 메일 수신과 운영 relay·운영 계정 왕복: 승인 필요 범위라 하지 않았다.
- macOS·Linux의 `orca` 실행 파일 탐색은 PATH 이름만 본다(코드의 `shortcut:` 주석).

## 정리

컨테이너 2개(검사용 `kl-orca-bot-gate`는 `--rm`)와 이번에 받은 redis 이미지를 지웠다. postgres 이미지는 원래 있던 것이라 남겼다. 작업 파일은 `.devtools/SAR-ORCA-BOT-001-DEV/`(gitignore)에 있다.
