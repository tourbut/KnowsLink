---
title: SAR-MVP-003-BIDIRECTIONAL-TESTER — 고정 양방향 시험 QA
status: draft
updated: 2026-10-04
owner: tester
tasks: [SAR-MVP-003-BIDIRECTIONAL-TESTER]
summary: "cd60e7f의 격리 Postgres 왕복, held 경계, 10초 body timeout 실패를 기록한다"
---

# SAR-MVP-003-BIDIRECTIONAL-TESTER — 고정 양방향 시험 QA

## 판정

판정 후보는 `cd60e7f87eb5ce137eca887980f232b3f67a18d0`다. 변경 기준은 `f2849486ed48295e239714700e651d30c32f1c2c`다.
기록 브랜치 `fullops/tester`의 착수 HEAD는 `080c2df21ecf96d5338f6dba5884a723c9a60ae3`다.
실행 위치는 `/tmp/sar-mvp-003-qa`의 detached clone이다. clone HEAD는 시작과 끝이 같다. 추적 파일 diff는 없다.
독립 관찰 스크립트의 종료코드는 1이다. live 단계와 closed 단계의 종료코드는 0이다. held 단계의 종료코드는 1이다.
두 독립 MCP와 Codex stdin, Grok `--node` `--bundle` launcher의 왕복 ID·from·to·text가 일치했다.
격리 Postgres 집계는 owner 2, `relay.test.message` 4, 업무 intent 0, gate 0, envelope/inbox payload 0이다. held 전후 집계는 모두 0이다.
빈 allowlist로 relay를 다시 만든 뒤 서명된 trial send는 `403 sender_not_allowed`다. trial 행 수는 4로 유지됐다.
기본 10초 body timeout은 첫 청크 뒤 멈춘 응답에서 request를 끝내지 못했다. 이 실패의 등급은 medium이다. 제품 코드는 수정하지 않았다.
`make test`, `make verify-mvp`, `make verify-grok-plugin`, 레포 밖 `env -i` 네 도구 프로브의 종료코드는 0이다.
실제 Grok 계정, Cloudflare, 사용자 credential, 공개 왕복은 미검증이다.
증거는 [SAR-MVP-003-BIDIRECTIONAL-TESTER-test/](SAR-MVP-003-BIDIRECTIONAL-TESTER-test/)다.

## 기준

날짜는 2026-10-04다. 공통 기준은 `fullops-common-0.3.2`다. 정본은 `.fullops-squad/project.md`다. `fullops-test`와 문서 작성 규칙을 적용했다.
Jev keep은 FULLOPS, project, 공통 README, ops-guide, adapters README, 핸드오버, coding-style, security, testing이다.
DEV 실행 기록은 로컬 준비 성공과 actual 계정 차단을 함께 적는다. DEV의 로컬 ID를 이번 Grok 계정 성공으로 쓰지 않는다.
사용자 시험 승인은 text 왕복이다. 기존 owner-only 보호는 유지한다. FullOps 업데이트는 범위 밖이다.
관찰 기대값은 실행 전에 고정했다. held 뒤 owner, message, gate는 0이다. live 뒤 owner는 2, trial metadata는 4, 업무 intent는 0, gate는 0, envelope와 inbox payload는 0이다. 빈 allowlist 재기동 뒤에도 trial 수는 4이고 서명 send는 403이다.

## 환경

| 항목 | 값 |
|---|---|
| 호스트 | Linux x86_64 |
| Node | v22.22.2 |
| Go | 1.27.1 |
| Postgres | compose `postgres:17`, 프로젝트 `knowslink-qa003-ad81fe984d`, loopback 포트 45887 |
| clone | `/tmp/sar-mvp-003-qa`, detached `cd60e7f87eb5ce137eca887980f232b3f67a18d0` |
| 프로토콜 | MCP `2025-11-25` |

`make install`과 `make build`는 clone에서 종료코드 0이다. cloudflared profile과 사용자 `~/.grok` 설정은 이 관찰에 쓰지 않았다.
기록 체크아웃에서 시작한 설치는 즉시 중단했다. 그 시점의 추적 diff는 없다. 판정 빌드는 clone에서만 남긴다.

## 독립 왕복

관찰 순서는 compose up, held, live, 빈 `KNOWSLINK_TEST_AGENTS`로 relay recreate, closed, compose down `--volumes`다.
up 종료코드는 0이다. recreate 종료코드는 0이다. down 종료코드는 0이다.
held 전 집계와 held 후 집계는 owner 0, gate 0, trial 0, business 0, payload 0이다.
live 후와 closed 후 집계는 owner 2, gate 0, trial 4, business 0, payload 0이다.

MCP 송신과 Grok 수신 ID는 `01a104fe-2f22-7a8c-9ed8-8863e1d962b4`다. text는 `Codex independent ping`이다. from은 `trial_codex`다. to는 `trial_grok`다. `untrusted`는 true다.
MCP 회신과 Codex 수신 ID는 `01a104fe-2fde-771f-9959-e91f07a3ae67`다.
Codex CLI stdin 송신 ID는 `01a104fe-3073-74f3-9c20-1d6799450e25`다. text는 `Codex cli ping`이다.
Grok launcher 회신 ID는 `01a104fe-31b1-712c-95ee-a4f661328655`다. Codex receive가 같은 ID·from·to·text를 반환했다.
두 MCP PID는 `551682`와 `551689`다. 송신자 큐는 비어 있다. claim 뒤 양쪽 큐는 비어 있다. trial 중 `knowslink_pull_once`는 held다.
같은 idempotency key와 같은 text는 같은 ID다. 다른 text는 실패다. HTTP 재전송은 200과 원래 ID다. 다른 text는 `409 idempotency_conflict`다.
손상 서명은 `422 invalid_signature`다. exp `+310000` ms는 `422 ttl_too_long`다. 과거 exp는 `422 expired`다.
다른 agent와 업무 intent `schedule.query`는 `403 sender_not_allowed`다. Authorization이 없으면 `401 invalid_auth`다.
`/v1/test/owners`, `agents`, `invites`, `authorize`, `gate-consume`은 agent token으로 `403 sender_not_allowed`다.
명령 stdout와 stderr에 credential이 없다. `actualGrok`은 unverified다.

## held, origin, launcher

held MCP는 도구 4개를 반환한다. 이름은 `knowslink_pull_once`, `knowslink_status`, `knowslink_test_receive`, `knowslink_test_send`다.
status는 `state` held, `actualConnection` held, `webhook` false, `evidenceFetch` false다. 세 call 도구는 held다. canary 연결 수는 0이다.
`testTransport()`는 held에서 null이다. Access secret이 없는 `test-remote`는 null이고 key 파일을 읽지 않는다.
승인 origin `https://link.knowslog.com`은 구성되고 fetch 수는 0이다. `sendText`는 호출하지 않았다.
거절 origin은 `http://link.knowslog.com`, `https://example.com`, `https://127.0.0.1`, path, query, userinfo, port 444다. loopback 모드는 공개 origin을 거절한다.
65537바이트 응답은 `too large`다. `302` Location `http://127.0.0.1:1/leak`는 거절된다.
mode `0644` config와 symlink config는 `run_trial.py` 종료코드 1이다. 출력에 sentinel이 없다.
`0600` config와 `--node`, `--bundle`은 부모 `AGENT_CREDENTIAL`을 config 값으로 바꾼다. `CF_ACCESS_CLIENT_SECRET`은 빈 값이다. 종료코드는 0이다. 부모 secret은 출력에 없다.

## 바디 타임아웃

held의 `body timeout rejects`는 실패다. detail은 `hung`이다. `body timeout is about 10s`는 실패다. detail은 `15001`이다.
대상은 기본 6인자 `TestTransport`다. 저장 `timeoutMs`는 10000이다. 서버는 부분 JSON을 쓰고 본문을 끝내지 않는다.
후속 측정에서 `AbortSignal`은 10001 ms에 `TimeoutError`로 abort한다. `reader.read()`는 그 시점에 reject하지 않는다. 12010 ms에도 request promise가 남아 있다.
같은 기본 생성자의 다른 측정은 12508 ms에 hung이다. 서버를 닫은 뒤에야 read가 `TypeError`로 끝난다.
명시 `timeoutMs` 100은 103 ms에 `TimeoutError`다. `make test`의 `trial-boundaries`는 이 100 ms 경로를 통과한다. 출력에 `100ms timeout` PASS가 있다.
저장된 10000 ms 시험 중 일부는 약 10002 ms에 `TimeoutError`다. 일부는 12000 ms를 넘긴다. 기록된 관찰 실패는 유지한다.
영향은 멈춘 응답 본문이 문서의 10초 상한 뒤에도 `Adapter.request`를 붙잡는 것이다. 인증 우회, secret 출력, trial 행 추가는 없다. 64 KiB 상한과 redirect 거절은 통과했다.

## 회귀와 단독 도구

| 명령 | 종료코드 |
|---|---|
| `make test` | 0 |
| `make verify-mvp` | 0 |
| `make verify-grok-plugin` | 0 |
| 레포 밖 `env -i` 프로브 | 0 |

`make test`는 `go test -race`와 adapter `npm test`다. Go `TestTrialHTTP`와 `TestTrialMessageSafety`는 `verify-mvp` 로그에서 PASS다.
`verify-mvp`의 격리 왕복 ID는 이번 독립 관찰 ID와 다르다. `actualGrok`은 `unverified`다. Tunnel은 사용하지 않았다.
`verify-grok-plugin`은 임시 HOME에서 validate, install `--trust`, list `knowslink` `0.1.0`, doctor `4 tools discovered`를 통과한다.
사용자 `~/.grok/config.toml` SHA256은 실행 전후 `25230ebf41576694b9219616e17b839781e990d60edfd86cbcf7b50212ad791b`다.
`trusted_folders.toml` SHA256은 `1f4eb552cb7eefe711b3bd00de814e29cbc114a20cbedad83c4bb299f7c81820`다.
`installed-plugins` 내용 hash는 `57650efb533c2026165dd0a180b865d91b9d36b549b5d0e9d005dfc90d94cfa6`다. 세 값은 실행 전후가 같다.
레포 밖 cwd `/tmp`의 `env -i`는 도구 4개다. status 본문은 `{"state":"held","transport":"pull","actualConnection":"held","webhook":false,"evidenceFetch":false}`다. stderr 길이는 0이다.

## 문서와 댓글 절차

토큰 스캐너의 `documents_ok`는 false다. 실행 기록과 `adapters/README.md`는 요구 문자열을 모두 가진다.
ops-guide 13장은 `--node`, `--bundle`, 네 도구 이름, `idempotency_key`, `test-remote`, `untrusted` 문자열을 반복하지 않는다.
13.3의 Codex 수신 명령은 `python3 scripts/run_trial.py --config <0600 environment.json> receive`다. 송신은 stdin text와 `send <key>`다. 이번 live가 이 형태로 종료코드 0이다.
Grok Command 인자는 launcher, config, 고정 Node, standalone bundle이다. 댓글 초안의 Arguments는 `run_trial.py --config ... --node ... --bundle .../plugin.js`다. live의 Grok 프로세스가 이 형태로 receive와 reply를 수행했다.
`--bundle`을 생략하면 `run_trial.py`는 `adapters/dist/trial-cli.js`를 실행한다. Codex 경로가 그 분기다. 댓글의 CLI 대안은 `--node`만 주고 `receive`를 붙인다. 그 한 줄은 별도 프로세스로 반복하지 않았다.
댓글은 게시하지 않았다. `https://link.knowslog.com`에 `sendText`를 보내지 않았다. `REVIEWED_SHA`는 치환하지 않았다.

## allowlist 비교식

요약 JSON의 `allowlist_empty`는 false다. 스크립트는 `printenv` stdout가 빈 문자열인지를 본다.
값이 빈 변수에서 로컬 `printenv` stdout는 줄바꿈 하나이고 종료코드는 0이다. 컨테이너 stdout는 저장하지 않았다.
같은 실행의 closed 검사는 `403 sender_not_allowed`다. trial 수는 4다. 이 불리언을 제품 실패로 기록하지 않는다. 관찰 스크립트는 재실행하지 않았다.

## 재사용한 UI

owner dashboard 템플릿 리터럴은 `78b1d92c8aa626245d3349ffaf7367d28f1dd3ef`와 `cd60e7f87eb5ce137eca887980f232b3f67a18d0`에서 같다. 길이는 362다.
Chrome QA는 반복하지 않았다. `http.go`와 adapter 소스는 바뀌었다. Go와 adapter 검사는 이번 clone에서 다시 실행했다.

## 하지 않은 검증

실제 Grok 계정, Cloudflare token, trial path 앱, AUD, hosted relay, 공개 왕복, 이메일 OTP, 유료 API, 사용자 브라우저 UI는 실행하지 않았다.
실제 배포 HEAD는 문서의 `28bd1bb`다. 이번 QA는 그 배포를 이동하지 않았다.
`/home/shin/deploy/knowslink-state`는 읽지 않았다.
