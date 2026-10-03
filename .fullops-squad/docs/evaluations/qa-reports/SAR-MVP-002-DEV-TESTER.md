---
title: SAR-MVP-002-DEV-TESTER — Grok Bot 플러그인 로컬 독립 QA
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-MVP-002-DEV-TESTER]
summary: "후보 552586b의 ZIP, stdio MCP 경계, 격리 SQL 연동 독립 검증을 기록한다"
---

# SAR-MVP-002-DEV-TESTER — Grok Bot 플러그인 로컬 독립 QA

## 판정

판정 후보는 `552586b6e886f95bffa9a000a031ea03070afedb`다. 실행 위치는 기록 체크아웃과 분리된 detached clone이다. clone HEAD는 시작과 끝이 같다.
ZIP SHA256은 `0e671d1a89c141d896034fff31619b9cd2148b73b567adbc3a97126031989117`이다. DEV 기록과 같다.
stdio MCP initialize/discovery, held, loopback, redirect, persist/ACK/claim 실패 중단, 격리 SQL의 gate와 최소 denied result, 동시 호출 경계의 자기 종료코드는 0이다.
새 결함은 없다. 새 critical/high는 없다. 실제 Grok Bot 계정, marketplace, hosted runtime, 유료 API는 held다.
UI와 Go 파일은 Chrome QA `9584aafcbb5fee88dcc6d618caf660884f6a527d`와 같다. 이번 실행은 새 브라우저 QA가 아니다.
제품 코드는 수정하지 않았다. 증거는 [SAR-MVP-002-DEV-TESTER-test/](SAR-MVP-002-DEV-TESTER-test/)다.

## 기준

날짜는 2026-10-03이다. 기록 브랜치는 `fullops/tester`다. 시작 HEAD는 `e8cf6cf531142bc1b440d46d527a7f6720dc317f`다.
공통 기준은 `fullops-common-0.3.2`다. 정본은 `.fullops-squad/project.md`다. 문서 작성 규칙과 `fullops-test`를 적용했다. 예외는 없다.
D03, D05, D10과 `adapters/README.md`의 loopback, redirect 거부, owner gate, 무정책 denied, 동시 pull `busy`, actual connection held를 판정 기준으로 사용했다.
Jev find, documents-find, context, route를 로컬에서 읽었다. route는 implementation/tester다. context의 conflict_ids와 caution_ids는 비어 있다. 민감 원문은 외부로 보내지 않았다.
Jev 웹 시나리오는 만들지 않았다. 통과 조건은 명령 종료코드와 도구 결과로 판정된다.

## 환경

| 항목 | 값 |
|---|---|
| Node | v22.22.2 |
| Go | go1.27.1 |
| Docker | 29.4.3 |
| clone | 레포 밖 detached `552586b6e886f95bffa9a000a031ea03070afedb` |
| Compose project | `knowslink-mvp-c564765cd4`. 종료 뒤 volume 삭제 |
| 공유 fixture | `/home/shin/deploy/knowslink/build/qa-fixture.json`, mode 600 |

첫 `make install`은 기록 체크아웃에서 실행됐다. 추적 파일 diff는 없었다. 판정 설치는 clone의 `make install`이다.
clone의 `--seed`는 clone `build/qa-fixture.json`을 mode 600으로 만들었다. 로그 확인 뒤 그 파일을 삭제했다. 공유 fixture의 SHA256은 실행 전과 후가 `d5c6e55306a5be50c961018400315ce1be2a56fc3dac07e83e7cf276eed7d3db`다.
기존 컨테이너 이름 8개의 집합은 실행 전과 같다. 격리 project만 만들고 제거했다.

## 패키지

`make plugin` 종료코드는 0이다. archive는 8개 항목이다. 항목은 marketplace manifest, plugin manifest, `mcp.json`, README, skill, `dist/plugin.js`, package.json, `THIRD_PARTY_NOTICES.txt`다.
`node_modules`, fixture, `.env`, PEM, secret, credential 경로는 없다. inventory 종료코드는 0이다.
ZIP timestamp는 1980-01-01이다. 같은 패키징 입력과 고정 timestamp에서 DEV hash가 재현됐다.
압축을 해제한 `knowslink/dist/plugin.js`에 대한 `mcp.test.js` 종료코드는 0이다.

## MCP 경계

자동 경계 검사가 다음을 통과로 판정했다.

| 경계 | 관찰 |
|---|---|
| initialize/discovery | 도구 이름은 `knowslink_pull_once`, `knowslink_status` |
| 기본 held | 자격과 relay URL이 있어도 mode 없으면 `held`, `isError=true`. relay 요청은 0 |
| production mode | `held`. relay 요청은 0 |
| 설정 누락 | `synthetic-loopback`만 있으면 `unconfigured` |
| loopback | `https`, `file`, path, query, userinfo는 `failed`. 이 호출의 relay 요청은 0 |
| redirect | loopback 302를 따라가지 않는다. 요청 수는 1 |
| persist/ACK/claim 실패 | mock HTTP 403에서 실패 경로 뒤의 판단과 result 요청이 없다 |
| 도구 출력 | state, transport, `actualConnection=held`, webhook false, evidenceFetch false. stderr는 빈 값 |

이 실패 전파는 mock fetch다. 실제 SQL 실패 전파는 아래 Go 통합 검사가 맡는다.
같은 프로세스에서 첫 pull이 loopback `/v1/registry`에서 대기하는 동안 둘째 pull은 1ms에 `busy`였다. 둘째 호출의 relay 요청은 없다. 대기 중 status는 `synthetic_only`다. 첫 호출이 `failed`로 끝난 뒤의 pull은 `failed`이며 relay에 도달했다.
probe 종료코드는 0이다. 출력에 probe 자격과 key 파일 내용은 없다.

## 격리 SQL

`make verify-mvp` 종료코드는 0이다. 자기 종료코드는 로그의 `exit:` 값이다.

| 명령 | 종료코드 |
|---|---|
| compose up `--build --wait` relay | 0 |
| `go test -tags=integration -race -count=1 -v ./internal/relay` | 0 |
| `node adapters/dist/synthetic.js http://127.0.0.1:50975` | 0 |
| 같은 명령 `--seed` | 0 |
| compose down `--volumes` | 0 |

Go 통합에서 다음 하위 검사가 통과했다. `auth_first_atomic_replay_and_rollback`, `last_lease_ack_and_shared_claim_restart`, `owner_gate_csrf_binding_and_no_effect`, `unpair_generation_expiry_max_attempts_and_clock`, `agent_cannot_process_human_delivery`, `legacy_unrouted_claim_parent_boundaries`, `rotation_pop_and_key_reassignment`, gate `denied`/`expired`/`revoked`/`unavailable`, `TestApprovalAndResultInstanceBinding`, `TestPendingAndConcurrentAccept`, `TestLegacyClaimsCannotReachParentBoundaries`.
권한 실패는 403, 401, 409로 남고 성공 경로는 200 또는 303이다. 잘못된 claim, human delivery, legacy route는 parent 경계에 도달하지 않는다. gate consume은 executable과 disclosure를 false로 유지한다. `status=done` result는 403이고 최소 `denied` result는 200이다.
동시 send replay 12개는 200이다. 동시 claim의 승자는 1개다. 동시 invite accept 8개는 활성 관계를 1개로 유지한다.
MCP bundle의 synthetic pull은 owner approve 303 뒤 receipt completion `denied`를 확인했다. 로그 문장은 `PASS: TS verify/persist/ACK/claim, policy deny, owner approve and denied result; no effects`다.
`--seed`는 파괴된 격리 DB의 loopback gate URL만 로그에 남겼다. 자격 값은 로그에 없다.

## 재사용한 증거

| 증거 | SHA | 이번 실행 |
|---|---|---|
| `make test`의 Go race `./cmd/... ./internal/...` | `c4ebbecd3ef91be10ecbb517fe451e02769358bc` | 재사용. 대상 소스 diff 종료코드 0. [product-test.txt](../../exec-plans/logs/SAR-MVP-002-DEV/product-test.txt) |
| Chrome owner UI | `9584aafcbb5fee88dcc6d618caf660884f6a527d` | 재사용. Go, HTML, CSS, TSX, JSX diff 종료코드 0. 새 브라우저 QA로 기록하지 않음 |

`make lint`는 `scripts/package_plugin.py`가 `c4ebbec` 이후에 바뀌어 clone에서 다시 실행했다. 종료코드는 0이다.
`make build`는 `make plugin`과 `make verify-mvp`의 의존 대상으로 두 번 실행됐다. 두 상위 명령의 종료코드는 0이다.
실제 Grok Bot 설치, 도구 검색, hosted Node, marketplace 등록, 계정 로그인, 유료 호출, 운영 배포는 실행하지 않았다. 상태는 held다.

## 명령

| 명령 | 종료코드 | 로그 |
|---|---|---|
| `make install` | 0 | [install.log](SAR-MVP-002-DEV-TESTER-test/install.log) |
| `make plugin` | 0 | [plugin.log](SAR-MVP-002-DEV-TESTER-test/plugin.log) |
| 압축 해제본 `node adapters/dist/mcp.test.js` | 0 | [mcp-extracted.log](SAR-MVP-002-DEV-TESTER-test/mcp-extracted.log) |
| `node concurrent.mjs` | 0 | [concurrent.log](SAR-MVP-002-DEV-TESTER-test/concurrent.log) |
| `make verify-mvp` | 0 | [verify-mvp.log](SAR-MVP-002-DEV-TESTER-test/verify-mvp.log) |
| `make lint` | 0 | [lint.log](SAR-MVP-002-DEV-TESTER-test/lint.log) |
| ZIP inventory | 0 | [inventory.log](SAR-MVP-002-DEV-TESTER-test/inventory.log) |

`verify-mvp.log`는 Compose progress의 줄 끝 공백 50줄을 제거한 파일이다. 원본 SHA256은 `23aa83a711454c0001f80a1bbd23bbf611f0f781e94b3dbd8cfea274bca6485b`다. 정규화 SHA256은 `485d152365bb40fc64ebb034cfc7910ea0cc73b67d7d1d52c99a0a2c9c64387f`다. 명령과 종료코드와 테스트 이름은 유지했다.
기준 ref `dbdd70086971285b790683f362702e5a9ff55acd`의 FullOps lint와 `deliverables.py --strict`는 기록 커밋의 깨끗한 HEAD에서 실행한다. 그 종료코드는 같은 test 디렉터리에 남긴다.

## 한계

MCP의 persist/ACK/claim 실패 순서는 mock fetch다. 같은 실패를 실제 Postgres에 주입하는 별도 MCP probe는 실행하지 않았다. 실제 SQL의 실패 코드는 Go 통합 검사가 판정했다.
in-process `busy`는 지연되는 loopback HTTP로 관찰했다. 두 OS 프로세스의 동시 pull은 이번 probe 범위가 아니다. 공유 claim의 단일 승자는 SQL 검사가 판정했다.
새 브라우저 QA, 실제 계정 설치, hosted runtime은 실행하지 않았다. 미실행 항목을 통과로 기록하지 않는다.
