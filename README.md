# KnowsLink

승인된 에이전트 사이의 선택적 인간개입 전달 제품이다. 현재 코드는 초기 개발 구성이다. 업무 API와 human-gate는 후속 기능 과제다.

Go relay는 DB ping과 `GET /healthz`만 제공한다. 별도 `cmd/migrate`는 goose SQL-only 경로다. TypeScript adapter는 벤더 연결 없이 `unimplemented` 상태를 출력하고 종료한다.

## 개발 도구와 설치

Go 1.27.1, Node 22.22.2, npm 10.9.7, Python 3, GNU Make, Docker Engine와 Compose가 필요하다. Go 1.21 이상의 `go`는 `go.mod`에 지정한 도구 체인을 자동으로 내려받을 수 있다. Node 버전은 `.nvmrc`에 있다. Go 의존성은 `go.mod`·`go.sum`, adapter 의존성은 `adapters/package-lock.json`으로 고정한다.

```sh
make install
make lint
make test
make build
make verify
```

`make lint`는 Go 포맷·vet·module 무결성, TypeScript 포맷·ESLint·타입, YAML 포맷·Compose 구조, 검증 스크립트 문법을 검사한다. 자동 수정하지 않는다. Go 검사는 `cmd/`와 `internal/`로 한정한다. `node_modules`의 외부 Go 코드는 검사하지 않는다.

`make verify`는 임시 복제본에서 Go 포맷과 TypeScript 포맷·lint·타입 위반을 주입한다. 등록된 `make lint`가 실패하고 원복 후 통과하는지 확인한다. 제품 파일을 바꾸지 않는다.

빌드 결과는 `build/relay`, `build/migrate`, `adapters/dist/`에 생긴다. 생성물과 `node_modules/`는 Git에서 제외한다.

## 로컬 실행과 중지

예시 값은 로컬 검증용 더미다. 실제 키나 Tunnel token은 없다. 기본 profile은 cloudflared를 시작하지 않는다.

```sh
docker compose --env-file .env.example up --build --wait relay
curl --fail http://127.0.0.1:8080/healthz
docker compose --env-file .env.example logs migrate
docker compose --env-file .env.example down
```

Compose는 Postgres 상태 확인 뒤 one-shot migrate를 실행한다. migrate 성공 뒤 relay를 시작한다. Postgres 포트는 게시하지 않는다. relay 포트는 `127.0.0.1:8080`에만 게시한다. `down`은 DB volume을 보존한다.

실제 migration 파일이 없으므로 migrate는 DB ping 후 `no SQL migrations; no-op`을 출력한다. 업무 테이블이나 goose version 테이블을 만들지 않는다. relay는 migration을 자동 실행하지 않는다.

```sh
make verify-runtime
npm run start --prefix adapters
```

`make verify-runtime`은 고유 Compose 프로젝트와 임시 DB 저장소로 기동·no-op·상태 확인·필수 설정 누락·DB 오류를 검증한다. 검증 컨테이너를 종료한다. cloudflared는 네트워크를 차단한 별도 컨테이너에서 token 누락 실패만 확인한다. 실제 Tunnel은 연결하지 않는다.

설정 이름:

- `DATABASE_URL`: 필수 Postgres URL. Go 프로세스는 `.env` 파일을 자동으로 읽지 않는다. Compose가 값을 전달한다.
- `POSTGRES_PASSWORD`: Compose Postgres 초기화에 필수다. `DATABASE_URL`의 비밀번호와 일치시킨다.
- `RELAY_ADDR`: Go listen 주소다. 직접 실행 기본값은 `127.0.0.1:8080`이다. Compose 내부는 `0.0.0.0:8080`이다.
- `RELAY_PORT`: Compose가 loopback에 게시하는 포트다. 기본값은 `8080`이다.
- `MIGRATIONS_DIR`: 별도 migrate 프로세스의 SQL 디렉터리다. 기본값은 `db/migrations`다.
- `TUNNEL_TOKEN`: 선택 `tunnel` profile의 필수 자격이다. 예시에서는 비어 있다. token 없는 cloudflared 실행은 실패한다.

실제 설정은 Git에서 제외한 `.env.local` 등에 보관한다. shell 환경 변수는 Compose의 env 파일보다 우선한다. 값을 로그나 커밋에 넣지 않는다. 이번 구성으로 운영 Tunnel이나 배포를 활성화하지 않는다.

## 후속 범위와 근거

`db/migrations`와 `db/queries`는 비어 있다. `sqlc.yaml`은 Postgres·`sql_package: pgx/v5`로 준비했다. 업무 SQL이 없어 SQL 적용·sqlc 생성 검사는 미적용이다. 검사 실행을 위한 가짜 스키마는 만들지 않았다.

등록·페어링·서명·ingest·queue·lease·ACK·receipt·공유 실행 claim·approve/deny·실제 어댑터 연결은 미구현이다. 첫 기능 MVP의 human-gate는 Go `net/http` + `html/template`로 구현한다. 별도 TypeScript frontend는 추가하지 않는다. webhook과 evidence fetch는 비활성이다.

요구사항은 [D02](.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md)다. 기술 정본은 [D03 아키텍처](.fullops-squad/docs/design-docs/architecture.md)와 [D03 스택](.fullops-squad/docs/design-docs/tech-stack.md)이다. 직접 검증 결과는 [실행 기록](.fullops-squad/docs/exec-plans/phases/SAR-SETUP-001-DEV.md)에 있다. 독립 QA와 코드 리뷰는 완료 SHA를 대상으로 후속 수행한다.
