---
title: SAR-PUBLIC-AGENTS-001-TESTER — 일반 회원 agent 연결·키·관계 QA
status: draft
updated: 2026-10-06
owner: tester
tasks: [SAR-PUBLIC-AGENTS-001-TESTER]
summary: 후보 d1eef9b의 연결·키·관계 독립 QA 통과와 미실행 운영 확인을 기록한다
---

# SAR-PUBLIC-AGENTS-001-TESTER — 일반 회원 agent 연결·키·관계 QA

## 판정

판정 후보는 `d1eef9bb90b9726149980320c42fb1fdbcaf584a`다. 기준은 `d2f7ba5aeb6644fd2b27fe6ada5b61db5976933b`다. 실행 위치는 `/tmp/knowslink-agents-qa-d1eef9b`의 detached checkout이다. 시작과 끝의 HEAD는 같고 추적 트리는 비어 있다.
독립 QA는 통과다. 새 critical/high는 없다. 제품 코드, D10/D12 정본, board, 운영 컨테이너는 수정하지 않았다.
`make lint` 종료코드는 0이다. `make test` 종료코드는 0이다. `make verify-mvp` 종료코드는 0이다. 별도 Postgres 프로브 `TestQAAgentsIndependent` 종료코드는 0이다.
실제 이메일, 운영 공개, 플랫폼, 최종 노우↔다닷은 미실행이다. 화면 캡처는 designer 범위라 하지 않았다.

## 기준

날짜는 2026-10-06이다. 공통 기준은 fullops-common-0.3.3, FULLOPS, project, document-writing, PS04–07, 해당 PS11, 개인정보, UX04–05다. 제품 구현이 늘린 규칙을 정본으로 쓰지 않았다.
도구는 Go 1.27.1, Node v22.22.2, Python 3.12.3이다. Workers Free만 사용했다. 유료 전환, Cloudflare 쓰기, 실메일, 외부 연결은 하지 않았다.
DEV가 보고한 `TestPublicAgentHTTP` 결과를 복사하지 않았다. 같은 후보는 `make verify-mvp` 안에서 그 검사도 실행했고, 교차 계정·한도·rate·CLI는 별도 프로브로 다시 실행했다.

## 실행한 명령

명령은 후보 루트에서 실행했다. 종료코드는 명령 자신의 값이다. 로그는 [SAR-PUBLIC-AGENTS-001-TESTER-test/](SAR-PUBLIC-AGENTS-001-TESTER-test/)에 있다.

| 명령 | 종료코드 | 결과 |
|---|---|---|
| `make lint` | 0 | gofmt, vet, `go mod verify`, adapter format/lint/typecheck, compose prettier, `check_compose.py` |
| `make test` | 0 | `go test -race ./cmd/... ./internal/...`와 adapter MCP, trial boundary, connect PASS |
| 격리 프로브 최종 | 0 | `TestQAAgentsIndependent` 3.98s, 하위 6개 PASS. Postgres up, migrate, down 종료코드 0 |
| `make verify-mvp` | 0 | 프로젝트 `knowslink-mvp-c03730a729`. 아래 8개 하위 명령이 모두 0 |

`make verify-mvp` 순서는 후보의 `scripts/verify_mvp.py` 그대로다.

1. `compose up --build --wait relay` 종료코드 0.
2. `compose stop relay` 종료코드 0.
3. `go test -tags=integration -race -count=1 -v ./internal/relay` 종료코드 0.
4. `compose up --wait relay` 종료코드 0.
5. `node adapters/dist/synthetic.js` 종료코드 0.
6. `node adapters/dist/synthetic.js --seed` 종료코드 0.
7. `node adapters/dist/trial-check.js` 종료코드 0.
8. `compose down --volumes` 종료코드 0.

Go 통합 로그의 관련 PASS는 `TestPublicAgentHTTP` 5개 하위, `TestEmailIdentity`, `TestGateFailureStates`, `TestPostgresSafety`, `TestTrialHTTP`, `TestTrialForeignAllowlistRevokesLease`다. 로그 끝은 `PASS: isolated Compose migration, Postgres races, TS adapter and Go owner UI; Tunnel unused`다. `verify-mvp.log`의 Compose 진행 줄 끝 공백만 공백 검사에 맞게 제거했다. 명령과 종료코드 줄은 그대로다.

프로브의 중간 실패 원문과 재실행 이유는 [probe-failures.md](SAR-PUBLIC-AGENTS-001-TESTER-test/probe-failures.md)에 있다. 실패는 프로브 기대값이다. 제품 결함으로 승격하지 않았다.

## 확인한 동작

격리 Compose 프로젝트의 Postgres에 마이그레이션을 적용한 뒤 `Service`가 그 DB를 사용했다. 메일은 메모리 inbox다. fixture 주소는 `example.test`와 문서용 IP다.

교차 owner와 연결:

- 다른 회원이 내 agent의 연결 발급을 호출하면 403이다. 미지원 클라이언트는 422이고 문장은 `미지원`이다.
- 발급 화면은 `연결 대기`, `개인키는 자기 클라이언트에만`, `node adapters/dist/connect.js prepare`, `발급만으로 연결되지 않습니다`를 보여 준다. `BEGIN `은 없다.
- 확인 GET은 `상태: waiting`, `공개키 준비 전입니다`, `아직 연결되지 않았습니다`다. 다른 회원은 403 `연결을 조회할 수 없습니다`다.
- 클라이언트가 서명에 owner, agent, public을 바꾸면 prepare는 422다. client를 바꾸면 403이다. 이 호출은 키를 만들지 않는다.
- 승인 전 complete는 403이다. `Sec-Fetch-Site: cross-site` confirm은 403이다. 다른 회원의 confirm은 401이다. owner confirm은 303이고 complete는 한 번 200, 다시 호출하면 401이다.
- agent credential로 `/v1/agents`, `/v1/keys`, `/v1/key-revoke`, `/v1/owner-revoke`를 호출하면 401이다.
- 저장 grant의 agent를 다른 owner의 agent로 바꾸면 `/v1/connect/info`는 401이고 키는 늘지 않는다.
- 저장 grant의 owner만 다른 회원으로 바꾸면 그 회원의 confirm은 303이다. complete는 401이고 어느 agent에도 새 키가 생기지 않는다. 자격은 발급되지 않는다.
- 승인된 complete 8회 동시는 승자 1명이다. 그 시점의 활성 키는 2개다.
- 검증 시각을 5분 뒤로 보내면 confirm, connect, key-revoke, agent-revoke, agent 생성이 422이고 문장은 `5분`이다. cancel은 303이고 info는 401이다. 기존 키 수는 유지된다.
- 만료된 연결의 확인 화면은 `만료`이고 info는 401이다. 키는 늘지 않는다.

키:

- register로 키 3개가 생기고 각 credential의 pull은 200이다. 네 번째 register 발급은 409 `운영 한도`다.
- key2만 철회하면 그 credential의 pull은 401이고 key1과 key3은 200이다. 홈은 key2 `철회`, key1 `활성`이다. 활성 키는 2개다.
- 잘못된 proof의 prepare는 422다. 그 연결을 취소하면 info는 401이고 기존 pull은 200이다. 활성 키 수는 그대로다.
- rotate 완료 뒤 활성 키는 1개다. 이전 세 credential은 새 Service에서도 `/v1/pull`, `/v1/persist`, `/v1/ack`, `/v1/claim`, `/v1/authorize`, `/v1/gate-consume`이 401이다.
- 제품에 `/v1/enqueue`는 없다. enqueue는 `POST /v1/send`로 확인했다. 옛 키의 `schedule.query`와 `relay.result` send는 401이다. 새 credential의 pull은 200이고 옛 key1 pull은 401이다.
- exec는 별도 경로가 없다. claim과 authorize를 exec 경계로, `relay.result` send를 result 경계로 확인했다.

관계:

- 다른 회원 초대 뒤 발신 홈은 `발신`, `수신`, `상태: pending`, `관계 세대: 1`이고 `수신 owner로 수락`은 없다. 수신 홈에는 `수신 owner로 수락`과 `수신 owner로 거절`이 있다.
- 발신 owner의 수락은 403이고 상태는 pending, 세대 1이다. pending에서 send는 403이다.
- 수신 거절은 303이고 상태는 denied다. 다시 초대한 뒤 세대 1 수락은 403, 세대 2 수락은 303이다. 수신 agent에 키가 있는 send는 200이다.
- 발신 unpair와 수신 unpair 모두 303이다. unpair 뒤 send는 403이다. 수신 철회 상태는 revoked, 세대 3이다.
- 같은 owner의 두 agent 초대는 pending으로 남고 홈에 `수신 owner로 수락`이 있다. 자동 수락은 없다.
- 같은 수락 8회 동시는 모두 303이고 그 쌍은 1개, 상태는 active, 세대는 1이다. 세대 0 수락은 403이다.

한도와 정리:

- 활성 외국 agent 200이면 생성은 409다. 새 Service에서도 409다. 하나 철회하면 한 번 303이고 다음은 409다.
- 외국 owner를 비활성으로 두면 자기 agent 5개까지 303이고 6번째는 409다. 하나를 철회하면 다시 303이다.
- 외국 owner가 비활성이면 다른 회원 생성은 303이다. 다시 활성이면 그 생성은 409다.
- 자기 pending 10개면 초대는 409다. 만료 뒤 초대는 303이고 세대 1, 기한은 약 24시간이다.
- 활성 owner의 pending 200이면 다음 초대는 409다. 만료와 철회 뒤 다시 초대하면 세대 2다. 세대 1 수락은 403이고 상태는 pending이다.
- 자기 active 20개면 수락은 409이고 대상 쌍은 pending이다. 전역 active 400이면 수락은 409이고 상태는 pending이다.
- `http:new` 201개면 생성은 429 `이후 다시 시도`다. deny는 303이다. 회원 버킷 길이는 0으로 남고 쌍은 denied다. 없는 키의 key-revoke는 403이다.

rate:

- 로그인한 깨진 form `%`는 40회 422, 41회째 429다. agent는 생기지 않는다. 다른 회원 생성은 303이다.
- 세션 없는 생성은 같은 IP에서 30회 401, 31회째 429다. 다른 IP는 401이다.
- `/v1/connect/prepare`의 `{`는 30회 422, 31회째 429다. 다른 IP는 422다.
- `http:new`가 201이고 principal 버킷이 비어 있으면 그 IP의 깨진 prepare는 429다. 그 IP 버킷 길이는 0이고 `http:new`는 201로 남는다.
- 같은 agent id로 pull 40회는 200이다. rotate로 credential이 바뀌어도 `http:member:<agentId>` 길이는 40이다. 새 credential의 pull은 429이고 새 Service에서도 429다. 거부가 기록되면 길이는 41이다. 옛 credential은 401이다.

Node CLI:

- 잘못된 URL과 이미 있는 폴더는 종료코드 1, stdout은 비어 있다. 기존 파일 내용은 유지된다.
- prepare 성공 stdout은 `상태: 공개키 준비 완료`, `지문: SHA256:`, `자기 브라우저에서 확인:`으로 시작하고 줄은 3개다. 토큰과 `PRIVATE`는 stdout과 stderr에 없다.
- 폴더는 0700, `private.pem`과 `pending.json`은 0600이다. 다시 prepare하면 종료코드 1이고 pem 해시는 같다.
- 승인 전 complete는 종료코드 1이고 안내 문장은 `연결 실패. 아직 사용할 새 자격이 없습니다. 브라우저 상태·기한을 확인하세요. 완료 응답을 잃었다면 새 연결·회전으로 복구하세요.`다. `agent.json`은 없다.
- owner confirm 뒤 complete stdout은 `상태: 연결 완료. agent.json·private.pem은 이 클라이언트에만 보관하세요. idle pull은 10초 이상 간격입니다.`다. `agent.json`은 0600이고 credential은 pending.json에 없다.
- 두 번째 complete는 종료코드 1이고 `agent.json` 해시는 같다. 저장된 credential의 pull은 200이다. 같은 토큰 complete는 401이다. 요청 본문에 `BEGIN PRIVATE`는 없다. 활성 키는 1개다.

## 관측

저장 연결의 owner 필드를 다른 회원으로 바꾼 뒤 그 회원이 confirm하면 303이다. confirm은 저장 owner와 세션 owner를 비교한다. 이어서 complete는 401이고 키는 생기지 않는다. 외부에서 owner, agent, client, public 서명을 바꾼 prepare는 거절된다.
홈의 agent 목록과 관계 표시가 같은 `<code>agent_` 형태다. 프로브는 자기 agent를 `<h3><code>`로만 골랐다. 이것은 검사 조회의 경계이고 제품 결함으로 기록하지 않는다.

## 기존 서비스

프로브 전후와 verify-mvp 뒤의 컨테이너 ID는 같다. `knowslink-relay-1`은 `dfcd9d187117`, `knowslink-cloudflared-1`은 `07077b9ef5e4`, `knowslink-postgres-1`은 `bc3482dc52f2`, `knowslink-migrate-1`은 `b4ef4f742466`이다.
프로브 프로젝트는 `knowslink-agentsqa-<8 hex>`이고 실행 끝에서 `down --volumes` 했다. verify-mvp 프로젝트 `knowslink-mvp-c03730a729`도 `down --volumes` 뒤 남아 있지 않다.
후보 clone과 `/tmp/knowslink-agents-review-d1eef9b`의 HEAD는 `d1eef9bb90b9726149980320c42fb1fdbcaf584a`이고 porcelain은 비어 있다. 마감 확인 시각 2026-10-06T13:52:47+09:00에 coor는 `df112d10b24bf9ab389ee9d7246632849a605af9`, designer는 `48d12fae2dce35d92606b264313148f0a635b64e`, dev는 `0a83bbb68aabe9f1e6d8d146de82d2509ac073af`, ops는 `70f26bc0799e65e4647731612a8d3a7c098a5fec`이다. 네 체크아웃의 porcelain은 비어 있다. 이 QA는 그 체크아웃을 수정하지 않았다.

## 미실행과 남은 항목

실제 운영 SMTP, 일반 사용자 이메일, 운영 공개, Cloudflare 쓰기, 최종 노우↔다닷, Grok 계정, 다닷 실제 연결은 미실행이다. 일반 서비스 수락은 이 QA로 완료하지 않는다.
영상과 전체 UI 캡처는 designer 소유라 하지 않았다. HTML 문구만 동작 확인에 사용했다.
`make generate`, `make schema`, `make verify-grok-plugin`, `make verify-runtime`은 SQL과 플러그인 설치 차이가 없어 실행하지 않았다.

## FullOps lint

기록 체크아웃 `fullops/tester`에서 플러그인 `lint.py`를 깨끗한 HEAD로 실행했다. 각 실행의 product-lint(`make lint`)와 product-test(`make test`) 종료코드는 0이다.

보정 전 head `1cf8f573a632bb5050be7ef946f576eb463b2089`, 기준 `9c915dc71e2a872243ffec294126d4668b4d32a4`의 결과는 ERROR 1, WARNING 8, 실행 불가 0이다. `lint.py` 종료코드는 1이다. ERROR는 `docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md`의 DOC-003이다. 그 파일은 부모 `d111fd4fde1538492c3c335e78ed510408718b2e`에 이미 있던 빈 리뷰 양식이다. product-lint와 product-test는 통과했다.

coor 메시지 `msg_9161a46153fb`가 이 양식의 metadata만 보정하도록 허용했다. `deliverables.py --stamp --path docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md --owner ops --status draft --task SAR-PUBLIC-AGENTS-001-REVIEW --summary 미완료 리뷰 양식`을 실행했다. H1 `# SAR-PUBLIC-AGENTS-001-REVIEW 리뷰`부터 원문 30줄은 보정 전과 같다. front matter와 제목 앞의 빈 줄만 추가됐다. 제품 코드, 고정 후보, fixture, 실행 조건은 그대로다. ops 체크아웃 `70f26bc0799e65e4647731612a8d3a7c098a5fec`는 수정하지 않았다.

보정 커밋 `4c7593802f3b6091389c147c701f1af079b7dd1d`의 기준 `9c915dc71e2a872243ffec294126d4668b4d32a4` 결과는 ERROR 0, WARNING 8, 실행 불가 0이다. `lint.py` 종료코드는 0이다. 파일 31개, 추가 1928줄이다. DOC-003은 없다. WARNING은 `PLANS.md` SIZE-001 818줄(이전 744, 상한 500), `adapters/package.json` DEP-001, `adapters/src/connect.test.ts:20` SEC-001, `connections_integration_test.go` 317줄, `http.go` 571줄, `member.go` 312줄, `store.go` 458줄, SIZE-002 1928줄이다.

같은 head의 기준 `d111fd4fde1538492c3c335e78ed510408718b2e` 결과는 ERROR 0, WARNING 1, 실행 불가 0이다. `lint.py` 종료코드는 0이다. 파일 1개, 추가 5줄이다. WARNING은 `PLANS.md` SIZE-001 818줄(이전 813, 상한 500)이다.

`PLANS.md`는 coor 소유다. 이 QA는 결과 절만 추가했고 헤더를 깎지 않았다. DEP-001, SEC-001, 제품 파일의 SIZE-001, SIZE-002는 기준 `9c915dc71e2a872243ffec294126d4668b4d32a4` 이후의 제품 변경이다. 이 QA는 그 파일을 수정하지 않았다. SEC-001은 제품 테스트의 더미 의심이다.

산출물 strict는 이 문서를 고친 뒤 실행했다. 결과는 검사 13, 미작성 0, 문제 0, 경고 0, 종료코드 0이다.

## 산출물

- 이 보고서
- [시나리오](../scenarios/SAR-PUBLIC-AGENTS-001-TESTER.md)
- [증거](SAR-PUBLIC-AGENTS-001-TESTER-test/)
