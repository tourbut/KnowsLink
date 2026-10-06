# KnowsLink

승인된 에이전트 사이의 선택적 인간개입 전달 제품이다.
이번 후보는 합성 데이터의 로컬 relay.v1·owner 가입/키·pairing·queue/receipt·shared inbox/claim·Go approve/deny UI를 연결한다.
실제 벤더·calendar·도구 효과·정보 공개·positive silent done과 운영 공개는 후속이다.

## 설치와 검증

Go 1.27.1, Node 22.22.2, npm 10.9.7, Python 3, GNU Make, Docker Engine/Compose를 사용한다.

```sh
make install
make lint
make test
make build
make generate
make schema
make verify
make verify-runtime
make verify-mvp
```

`make test`는 unit/race 회귀를 수행한다. Postgres 검사는 integration build tag로 별도 실행한다. DB 환경이 없으면 해당 실행은 실패한다.
`make verify-mvp`는 실제 별도 Compose project에서 SQL migration·Postgres 경합·Go UI·TypeScript stub을 끝까지 검사한다.
Go integration 검사 동안 그 project의 relay를 멈춘다. relay의 1초 정리 sweep은 자기 시험 allowlist로 검사용 trial lease를 회수한다. TypeScript 검사 전에 relay를 다시 시작한다.
시험 전용 override는 자기 Postgres만 loopback 임시 포트에 연결한다. 제품 Postgres는 호스트 포트를 게시하지 않는다.
검사 종료는 자기 project의 컨테이너·볼륨만 회수한다. 다른 컨테이너·기존 Tunnel을 바꾸지 않는다.
`make verify`는 등록 product-lint에 실제 위반을 주입하고 실패 전파를 확인한다.
`make verify-runtime`은 migration·readiness·누락 설정·DB 불가·차단된 Tunnel·미설정 adapter를 검사한다.

sqlc는 v1.30.0이며 pgx/v5를 생성한다. `make generate` 뒤 `git diff --exit-code -- internal/database`로 재현을 확인한다.
`make schema`는 실제 migration SQL에서 D08을 생성한다. API 기동은 migration을 자동 적용하지 않는다.

## 독립 QA와 owner UI 실행

고유 project 이름과 빈 loopback 포트를 사용한다. 예시 비밀번호는 합성용 더미다.

```sh
POSTGRES_PASSWORD=example-local-only \
DATABASE_URL='postgres://knowslink:example-local-only@postgres:5432/knowslink?sslmode=disable' \
RELAY_PORT=18081 COMPOSE_PROFILES= KNOWSLINK_SYNTHETIC_SIGNUP=1 \
docker compose -p knowslink-qa-local --env-file .env.example up --build --wait relay
node adapters/dist/synthetic.js http://127.0.0.1:18081 --seed
```

seed는 두 owner/agent·active pair·정책 없음 denied 결과·별도 pending gate를 만든다.
`.env.example`은 합성 가입을 닫는다. seed가 쓰는 `/v1/owners`는 위처럼 셸에서 `KNOWSLINK_SYNTHETIC_SIGNUP=1`을 준 격리 실행에서만 열린다. 값이 없으면 seed는 403으로 실패한다.
`build/qa-fixture.json`은 0600 권한의 Git 미추적 파일이다. owner credential·agent credential·합성 private key가 있으므로 외부로 보내거나 로그에 출력하지 않는다.
브라우저에서 seed가 출력한 `/owner/gates/<id>`를 연다. HTTP Basic username은 fixture의 `b.owner.owner`, password는 `b.owner.credential`이다.
원요청은 생성 후 180초 만료다. 화면을 늦게 열면 seed를 다시 실행한다.

V-01은 pending의 verified typed body·policy·발신/대상·만료를 확인한다.
V-02는 POST Approve 또는 Deny 후 approved/denied 화면을 확인한다. GET 방문은 결정하지 않는다.
V-03은 만료까지 기다리거나 owner 인증으로 `/v1/key-revoke` 또는 `/v1/unpair`를 호출한다.
V-04의 원문 부재·권한 불명은 `TestGateFailureStates`가 별도 DB 상태로 재현한다. 운영 API는 원문 삭제용 test hook을 제공하지 않는다.
독립 시각 검수자는 고정 후보에서 해당 상태를 캡처한다. DEV 자동 HTML 검사는 직접 검수를 대체하지 않는다.

```sh
POSTGRES_PASSWORD=example-local-only \
DATABASE_URL='postgres://knowslink:example-local-only@postgres:5432/knowslink?sslmode=disable' \
RELAY_PORT=18081 COMPOSE_PROFILES= \
docker compose -p knowslink-qa-local --env-file .env.example down
```

`down`은 자기 QA volume을 보존한다. `--volumes`는 해당 합성 project를 폐기할 때만 추가한다.

## 일반 이메일 로그인 확인 (SAR-PUBLIC-IDENTITY-001)

회원 화면은 `/`(시작)·`/auth/verify`(확인)·`/home`(자기 owner 홈)이다. relay가 6자리 코드를 SMTP로 보낸다. 계약은 [D05](.fullops-squad/docs/design-docs/interface-design.md#sar-public-identity-001-회원-화면과-세션)다.
로컬 화면 검수는 `scripts/mail_sink.py`의 QA 전용 SMTP sink를 사용한다. 실제 메일을 보내지 않는다.

```sh
python3 scripts/mail_sink.py build/qa-mail --port 2525
DATABASE_URL='postgres://knowslink:<password>@127.0.0.1:<port>/knowslink?sslmode=disable' \
RELAY_ADDR=127.0.0.1:18082 KNOWSLINK_SMTP_URL=smtp://127.0.0.1:2525 KNOWSLINK_MAIL_FROM=noreply@knowslog.com \
./build/relay
```

브라우저로 `http://localhost:18082/`를 연다. cookie가 `Secure`이므로 `localhost` 또는 HTTPS 주소를 사용한다. 받은 메일은 `build/qa-mail/*.eml`(0600)에 있다. 코드는 캡처·로그에 남기지 않는다.
`KNOWSLINK_SYNTHETIC_SIGNUP=1`은 격리 로컬 합성 fixture(`/v1/owners`)에서 셸로만 준다. `.env.example`과 공개 후보는 이 값을 비운다. 운영 SMTP 값은 Git 미추적 `.env`로만 제공한다.

## API와 adapter

계약은 [D05](.fullops-squad/docs/design-docs/interface-design.md)다.
API Bearer credential은 owner/agent 역할을 구분한다. `/v1/owners` 가입은 합성용이며 실사용자 신원 인증이 아니다. 실제 회원은 이메일 코드 확인으로만 만든다.
key PoP는 owner ID·AgentID·kid·public에 묶인 Ed25519다. rotation은 이전 키를 revoke한다.
accept/approve/revoke는 owner만 수행한다. pending invite는 active pair slot이 아니다.

TypeScript 실행에는 `RELAY_URL`, `AGENT_CREDENTIAL`, `AGENT_ID`, `AGENT_KID`, `AGENT_KEY_FILE`이 필요하다.
PEM 키 파일을 사용한다. `ADAPTER_GATE=1`은 합성 judgment gate를 만든다. URL은 loopback만 허용한다.
adapter는 registry·signature를 재검증한 뒤 shared persist·ACK·claim을 수행한다.
하나의 claim만 성공한다. 재시작 후 claim 재발급을 하지 않으므로 duplicate effect는 만들지 않는다.
승인 후에도 무정책 query는 denied이며 commit은 실행 불가능하다. result 수신은 새 도구 명령이 아니다.
설정이 없으면 `unconfigured`를 출력하고 네트워크 연결 없이 종료한다.

## 운영 경계와 저장

기존 `DATABASE_URL`, `POSTGRES_PASSWORD`, `RELAY_ADDR`, `RELAY_PORT`, `MIGRATIONS_DIR`, `TUNNEL_TOKEN` 설정을 유지한다. 신원 설정은 `KNOWSLINK_SMTP_URL`, `KNOWSLINK_MAIL_FROM`, `KNOWSLINK_CLIENT_IP_HEADER`, `KNOWSLINK_SYNTHETIC_SIGNUP`이다.
기본 cloudflared profile은 OFF다. Postgres는 private network에 있고 relay는 loopback에만 게시한다.
공개 hostname 연결은 수락된 후보 이후 OPS 과제다. 공개 한도·신원 인증·독립 QA/UI/리뷰 없이 이 합성 후보를 공개하지 않는다.

singleton Postgres 행의 global lock은 로컬 합성 처리량의 의도적 한계다. 공개 규모의 성능 보장은 없다.
원문·inbox는 exp·철회·응답 완료에 삭제한다. 유휴 상태도 1초 정리 루프가 동작한다.
receipt metadata만 24h 보관한다. DB 삭제는 WAL/backup 완전 삭제가 아니다.
webhook·evidence fetch/preview·실제 벤더 연결은 OFF다.

[D02](.fullops-squad/docs/planning/product-specs/SAR-MVP.md), [D03](.fullops-squad/docs/design-docs/architecture.md), [실행 기록](.fullops-squad/docs/exec-plans/phases/SAR-MVP-001-DEV.md)을 따른다.
초기 골격 이력은 [SAR-SETUP-001-DEV](.fullops-squad/docs/exec-plans/phases/SAR-SETUP-001-DEV.md)에 보존한다.

## 일반 회원 agent·키·관계 (SAR-PUBLIC-AGENTS-001)

이메일 확인 뒤 `/home`에서 새 agent를 만든다. 자기 agent만 관리할 수 있다. 회원당 활성 agent는 5개다. 연결 권한 변경에는 5분 안의 이메일 재확인이 필요하다.

1. 홈에서 Node 22 로컬 클라이언트와 등록 또는 회전을 선택한다. 회전은 완료 시 기존 키를 모두 철회한다.
2. **자기 클라이언트 컴퓨터**에 저장소를 준비하고 `npm ci --prefix adapters`, `npm run build --prefix adapters`를 실행한다. 관리자 계정·공유 서버 SSH·서버 파일 배치는 필요하지 않다.
3. `node adapters/dist/connect.js prepare <서비스 URL> <새 비공개 폴더>`를 실행한다. 연결 수단은 argv·대화·로그 대신 비공개 표준입력으로 넣고 입력을 끝낸다. 터미널 입력 반향은 먼저 끈다 (`stty -echo`; 입력 후 `stty echo`). URL은 HTTPS root 또는 loopback HTTP root만 허용한다. 개인키는 로컬에서 생성한다.
4. CLI가 표시한 자기 브라우저 연결 확인 주소를 연다. 대상·클라이언트·권한·공개키 지문·기한을 확인하고 승인한다. 지문이 다르면 취소한다.
5. `node adapters/dist/connect.js complete <새 비공개 폴더>`를 실행한다. 연결 수단은 최대 10분·1회다. 완료 전에는 새 키와 agent credential이 활성화되지 않는다.

로컬 폴더는 0700이며 `private.pem`, `pending.json`, `agent.json`은 0600이다. `agent.json`의 relay·agent·kid·credential과 `private.pem`은 자기 adapter의 설정으로만 사용한다. 로컬 `RELAY_URL`, `AGENT_ID`, `AGENT_KID`, `AGENT_CREDENTIAL`, `AGENT_KEY_FILE`에 대응한다. 개인키·credential을 모델 대화·Git·캡처에 넣지 않는다. 기존 합성 adapter 모드와 시험 allowlist를 일반 메시지 성공으로 사용하지 않는다. 일반 text·실제 왕복은 MESSAGES 후속이다.

실패·취소·만료는 새 키를 연결하지 않는다. 완료 응답을 잃으면 자격을 다시 조회할 수 없다. 홈에서 새 연결을 **회전**으로 진행해 잃은 키를 철회한다. 홈에서 키별 지문·활성/철회를 확인하고 선택 철회할 수 있다. 키당 credential을 분리하므로 철회 키의 송신·pull·ACK·claim·authorize·result는 거부된다. 로그아웃은 키를 철회하지 않는다. 정상 idle pull 간격은 10초 이상이다.

관계 초대는 상대 agent 식별자만 사용한다. 상대 이메일과 회원 디렉터리는 없다. 수신 owner만 수락·거절하며 양측이 철회할 수 있다. 같은 owner의 두 agent도 명시적으로 수락해야 한다. pending은 24시간 뒤 만료되며 재초대는 새 세대다. 옛 화면의 결정은 새 세대에 적용되지 않는다. 포화 중에도 거절·키/agent 철회·관계 철회는 별도 정리 budget을 사용한다.

같은 상대에게 다시 초대하면 기존 대기 초대나 연결된 관계를 그대로 유지한다. 기한도 바뀌지 않는다. 거절·만료·철회 뒤에는 홈의 `새 초대 보내기`로 직접 다시 초대하며 상대가 새로 수락해야 한다. 자동 재연결은 없다.

agent의 키 기록(활성+철회)이 20개면 그 agent에는 새 키를 연결할 수 없다. 키 철회나 대기로 공간이 생기지 않는다. 새 agent를 만들어 따로 연결하고 각 상대와 새로 수락한다. 철회 agent 기록까지 포함해 owner당 10개면 새 agent를 만들 수 없다. 철회 기록은 최소 24시간 보존된 뒤 정리되며 홈 목록에서 사라진다. 그 뒤 다시 시도한다. 사라짐은 권한 복구나 백업 영구 삭제가 아니다.

지원 구현은 위 Node 로컬 CLI다. Grok Bot 컴퓨터에서 안내를 실행할 수 있는지는 해당 계정의 설치·실행 권한으로 별도 확인한다. 앱 카탈로그·다닷·OAuth·실제 외부 계정 연결은 이 후보의 검증 결과가 아니다. 운영 공개·실메일·유료 설정은 실행하지 않았다.

## 일반 회원 연결 확인·receipt — SAR-PUBLIC-MESSAGES-001

[Node CLI/MCP 연결 확인 절차](adapters/README.md#일반-회원-연결-확인-text--sar-public-messages-001)를 따른다. 서로 다른 두 자기 agent도 명시적 관계 수락 뒤에만 비민감 text·관련 답장을 주고받는다. 회원 홈에서 자기 agent와 요청 ID로 전달·처리·답장 ID·TTL·실패 복구를 확인한다. queued는 상대 수신 성공이 아니다. 정상 idle pull은10초 이상 간격이며 자동 wake/답장이 없다.

`make verify-mvp`는 일반 신원 HTTP·실제 Node prepare/complete·두 MCP 프로세스의 send/receive/reply/receipt와 shared Postgres 권한·한도·재시작을 검사한다. 합성 SMTP sink/fixture는 실메일 증거가 아니다. 실제 Grok Bot·다닷 운영 연결·독립 QA·직접 시각 검수·운영 공개는 후속이다. Workers Free·기존 서버/Tunnel 보호를 유지한다.
