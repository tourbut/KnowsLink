# 개발환경과 검증

이 문서는 제품 개발자를 위한 안내다. 웹 사용에는 개발 도구가 필요 없다.
Grok/Dots 사용자 설치는 [README](../README.md)를 따른다.
운영 계정·DB·Tunnel·키를 개발환경으로 복사하지 않는다.

## 도구와 버전

버전 정본은 [go.mod](../go.mod), [adapter package](../adapters/package.json), [Dockerfile](../Dockerfile), [Compose](../compose.yaml)다.

| 도구 | 프로젝트 기준 |
|---|---|
| Go | `1.27.1` |
| Node / npm | Node `>=22.22.2 <23`, npm `>=10.9.7 <11` |
| PostgreSQL | Compose의 digest 고정 `17-alpine` |
| Docker / Compose | Docker Engine과 Compose v2. 프로젝트의 `--wait`·health dependency가 동작해야 함 |
| Git·Make·sh·bash·Python 3 | Makefile·검증·배포 설정 검사의 실행 도구 |
| race 검사용 C toolchain | `go test -race`가 동작하는 호스트 환경 |
| sqlc | `make generate`가 `v1.30.0`을 지정하므로 전역 설치 불필요 |

Linux/POSIX shell을 기준으로 명령을 쓴다. Windows에서는 Go·Node·Python·Make·bash·C toolchain의 PATH를 각각 확인한다.
Windows에서 배포 스크립트 검사나 도구 부재로 실패하면 Linux의 같은 고정 SHA에서 검사하고 원 결과를 남긴다.
호스트 차이로 실행하지 못한 검사를 PASS로 쓰지 않는다.

## 새 개발 checkout

기존 폴더가 있으면 다른 이름을 쓴다. 승인된 대상 브랜치나 고정 SHA를 선택한 뒤 시작한다.

```sh
git clone https://github.com/tourbut/KnowsLink.git knowslink-dev-new
cd knowslink-dev-new
git status --short
make install
```

`make install`은 `go mod download`와 `npm ci --prefix adapters`를 실행한다.
의존성 잠금 파일을 유지한다. 버전 오류를 무시하는 설치 옵션을 추가하지 않는다.
Grok 전용 installer는 개발환경 전체 설치 명령이 아니다.

## .env와 로컬 Compose

`.env`가 아직 없을 때만 예제 파일을 복사한다. 기존 비밀 설정을 덮어쓰지 않는다.
실제 비밀값은 Git 미추적 `.env`나 호스트 비밀 저장소에 보관한다.

```sh
cp .env.example .env
chmod 600 .env
docker compose -p knowslink-dev-new --env-file .env up --build --wait relay
curl --fail http://127.0.0.1:8080/healthz
docker compose -p knowslink-dev-new --env-file .env ps
```

예제 비밀번호는 고립된 로컬 시험용이다. 실제 운영에 사용하지 않는다.
Compose는 private Postgres health를 기다리고, 별도 migrate의 성공 후 relay를 시작한다.
기본 relay 포트는 `127.0.0.1:8080`이다. 충돌하면 `.env`의 `RELAY_PORT`를 빈 포트로 바꾼다.
DB 포트는 호스트에 게시하지 않는다. `DATABASE_URL`의 `postgres`는 Compose 내부 서비스 이름이다.
기본 실행은 cloudflared를 시작하지 않는다. 실제 Tunnel token이나 `--profile tunnel`을 로컬 QA에 넣지 않는다.

멈출 때는 같은 project와 env 파일을 사용한다. 다음 명령은 컨테이너를 정지하고 볼륨을 보존한다.

```sh
docker compose -p knowslink-dev-new --env-file .env stop
```

`KNOWSLINK_SYNTHETIC_SIGNUP=1`은 고립된 합성 시험에서만 명시적으로 켠다.
실제 공개 환경에서 활성화하지 않는다. `KNOWSLINK_TEST_AGENTS`는 시험 allowlist이며 일반 회원 연결 설정이 아니다.

## 로그인 모드

| 모드 | 설정·동작 |
|---|---|
| Google | `KNOWSLINK_GOOGLE_CLIENT_ID`, `KNOWSLINK_GOOGLE_CLIENT_SECRET`, `KNOWSLINK_GOOGLE_REDIRECT_URL`을 함께 설정 |
| 캘린더 연결 | Google 설정에 `KNOWSLINK_CALENDAR_KEY`(base64 32바이트)를 더하면 회원 설정에 ‘캘린더 연결’이 열린다. 없으면 연결 절이 준비 중으로 보인다. 봉인은 회원 ID에 묶여 다른 회원 행에서 열리지 않는다. 격리 검사는 합성 제공자를 쓰며 실계정 연결은 D011 뒤 운영 과제다 |
| 이메일 코드 | `KNOWSLINK_SMTP_URL`, `KNOWSLINK_MAIL_FROM`을 설정. SMTP가 없으면 코드 발송 불가 |
| 합성 로그인 | 격리 fixture에서만 `KNOWSLINK_SYNTHETIC_SIGNUP=1`. 실메일·실계정 성공 근거가 아님 |
| 둘 다 미설정 | relay health 성공과 사용자의 로그인 성공은 별개 |

Google redirect는 개발용 OAuth client에 등록한 callback과 정확히 같아야 한다.
relay의 callback 경로는 `/auth/google/callback`이다. 운영 OAuth secret을 빌려 쓰지 않는다.
이 구현은 Google redirect에 HTTPS를 요구한다. 로컬 시험에도 별도 HTTPS origin과 TLS proxy를 준비한다.
plain HTTP Compose 주소는 Google callback으로 사용하지 않는다.
설정이 일부만 있거나 유효하지 않으면 기동이 실패할 수 있다.
Google 계정·동의 화면·정책과 실제 브라우저 왕복은 합성 provider 검사와 별개다.

SMTP URL은 `smtps://user:token@host:465` 또는 격리 `smtp://host:port` 형식이다.
비밀 URL을 argv·채팅·Git·로그에 넣지 않는다. `.env`에서 제공한다.
로컬 QA의 [mail sink](../scripts/mail_sink.py)는 실제 발송 서비스가 아니다.

```sh
python3 scripts/mail_sink.py /absolute/private/knowslink-mail-new --listen 127.0.0.1 --port 2525
```

호스트 relay는 이 sink에 연결할 수 있다. Compose relay의 `127.0.0.1`은 해당 컨테이너이므로 호스트 sink 주소와 다르다.
Compose QA에서는 격리 bridge에서 도달 가능한 주소만 사용하고 SMTP 포트를 공개하지 않는다.
회원 cookie는 Secure다. 로컬 이메일 화면 시험에는 브라우저가 신뢰하는 `localhost` 또는 HTTPS origin을 사용한다.
`KNOWSLINK_CLIENT_IP_HEADER`는 기본 빈 값이다. Tunnel만 원점에 접근하는 환경에서만 `CF-Connecting-IP`를 신뢰한다.

## 호스트에서 직접 실행

별도로 준비한 **개발용** Postgres를 사용할 때 `DATABASE_URL`은 해당 호스트와 DB를 가리켜야 한다.
위 Compose 기본 Postgres는 private이므로 호스트의 Go process에서 바로 접근할 수 없다.
현재 shell에 비밀값을 출력하지 않는 방식으로 환경을 공급한다. Go 바이너리는 `.env`를 자동으로 읽지 않는다.

```sh
make build
./build/migrate
./build/relay
```

순서는 DB 준비·migration Up·relay다. `MIGRATIONS_DIR` 기본값은 `db/migrations`다.
`RELAY_ADDR` 기본값은 `127.0.0.1:8080`이다. Compose의 `RELAY_PORT`와 구분한다.
운영 DB를 대상으로 migration을 시험하지 않는다. 복구 기준은 [운영 안내](operations.md#schema와-writer-보호)를 따른다.

## 빌드·검사

레포 루트에서 실행한다. 전체 검사 목록의 정본은 [Makefile](../Makefile)이다.

| 명령 | 검사 범위 |
|---|---|
| `git diff --check` | 패치 공백 오류 |
| `make lint` | Go format/vet/module, adapter format/ESLint/type, YAML·Compose·ingress·shell·Python 정적 검사 |
| `make test` | `cmd`/`internal` Go race, adapter 검사, shell guard 회귀 |
| `make build` | relay·migrate 바이너리와 TypeScript/standalone plugin bundle |
| `make generate` | pinned sqlc 생성. 이후 `git diff --exit-code -- internal/database`로 drift 확인 |
| `make schema` | 실제 migration 기반 테이블 문서 생성 |
| `make plugin` | 빌드 후 호스트별 패키지와 허용 파일 검사 |

문서만 바꾸면 링크·fragment·예제·정적 검사를 수행한다. 제품 테스트 재실행 여부는 실제 영향으로 정한다.
검사 명령의 종료코드를 그대로 남긴다. `| tail`로 앞 명령 실패를 가린 로그를 성공 근거로 쓰지 않는다.

## 격리 QA

`make test`만으로 실제 Postgres integration·브라우저·실계정 검증이 끝나지 않는다.
아래 도구는 고유 fixture를 만들고 자신의 자원을 정리할 수 있다. 삭제 금지 세션에서는 실행하지 않는다.
운영 DB·공유 Compose project·기존 키를 대상에 넣지 않는다.

| 명령·입력 | 목적 |
|---|---|
| `make verify` | 임시 복제본에서 lint 위반과 실패 전파 확인 |
| `make verify-runtime` | Compose readiness·실제 SQL migration·누락 설정·Tunnel OFF·미설정 adapter |
| `make verify-mvp` | 고유 Compose DB·실제 HTTP·서명·claim 경합·TS/Go 왕복·integration |
| `make verify-grok-plugin` | 임시 HOME의 Grok CLI 설치·doctor. 실제 Bot 앱 등록 증거가 아님 |
| integration build tag | 실제 격리 Postgres·`TEST_SYNTHETIC_DATABASE=1`·`TEST_DATABASE_URL` 필요. 누락하면 실패 |
| `integration browser` tags | 해당 consent 브라우저 검사에 `TEST_BROWSER` 필요 |

격리 통과·실제 Grok/Dots 대화·독립 QA·사람의 직접 화면 수락은 서로 다른 근거다.
본문·키·쿠키·callback secret을 테스트 로그로 복사하지 않는다.

## 실패 대처

| 증상 | 먼저 확인할 것 |
|---|---|
| `DATABASE_URL is required`·DB ping 실패 | 환경 공급 여부, DB host/네트워크/health. 비밀번호를 출력하지 않음 |
| migrate 실패 | migration 경로·SQL-only 계약·버전·권한. Down이나 구 writer로 우회하지 않음 |
| 로그인 버튼/메일 없음 | Google 세 설정·SMTP 설정 여부. health를 로그인 성공으로 해석하지 않음 |
| Google callback 실패 | 등록 redirect의 정확한 일치와 개발용 계정 정책 |
| Node engine·Make·race 실패 | 도구 버전·PATH·C toolchain·POSIX 환경 |
| MCP `held` | 정상 기본값. 신규 연결의 mode·URL·새 private 폴더와 명시적 동의를 확인 |
| 수락됐으나 표시 안 됨 | receipt·Inbox·ACK·host 턴을 [운영 표](operations.md#수신-알림-문제-판별) 순서로 구분 |

실패 로그와 fixture는 승인된 보존 규칙을 따른다. 다른 프로젝트 자원이나 기존 자격 폴더를 지워서 재시도하지 않는다.
