---
id: D03
title: KnowsLink 초기 기술 스택
status: review
updated: 2026-10-03
owner: dev
tasks: [SAR-SETUP-001-DEV, SAR-MVP-001-DEV]
upstream: [D02]
summary: 고정 도구 버전과 실제 SQL 및 JCS API 근거를 기록한다
---

# KnowsLink 기술 스택

## 잠긴 선택과 실제 버전

[D02](../planning/product-specs/SAR-MVP.md)와 고정 원천의 Go·TypeScript·Postgres·Compose 선택을 따른다. Go UI는 `net/http` + `html/template`이며 별도 TypeScript frontend는 없다.

| 항목 | 구현 버전과 고정 방식 |
|---|---|
| Go | `go.mod`의 1.27.1, Docker build 이미지 1.27.1-alpine와 digest |
| pgx | `github.com/jackc/pgx/v5` 5.10.0, relay는 pgxpool, migrate는 stdlib driver |
| goose | `github.com/pressly/goose/v3` 3.28.0, Provider API, SQL-only |
| Node·npm | Node 22.22.2 (`.nvmrc`), npm 10.9.7, package engines |
| TypeScript | 5.9.3, strict·NodeNext, npm lock |
| ESLint | 10.12.0, `@eslint/js` 10.0.1, typescript-eslint 8.71.0, flat config |
| Prettier | 3.6.2, 명시된 TS·JSON·설정·YAML 대상만 검사 |
| Postgres | 검증 버전 17.11, Compose의 17-alpine와 multi-platform digest, 호스트 포트 없음 |
| cloudflared | 2026.9.1와 multi-platform digest, 선택 profile |
| 이미지 런타임 | Alpine 3.23와 digest, non-root UID/GID 65532 |
| sqlc | config version 2, `sql_package: pgx/v5`; CLI v1.30.0; LockRelay·SaveRelay 실제 생성 |
| 보조 도구 | Python 3 표준 라이브러리, GNU Make, Docker Compose |

의존성 검증과 생성물은 `go.sum`, `adapters/package-lock.json`으로 고정한다. goose 의존성 검사 과정에서 다운로드된 SQLite 테스트 의존성은 제품 저장소로 사용하지 않는다. 런타임은 pgx Postgres만 연결한다.

초기 조회한 ESLint 9.39.3은 npm의 지원 종료 경고 때문에 채택하지 않았다. 최종 ESLint 10.12.0은 Node 22.22.2와 typescript-eslint 8.71.0의 peer 범위를 만족한다. npm 설치·실행으로 호환성을 확인한다.

## 명령과 검사 대상

루트의 `make install`, `make lint`, `make test`, `make build`, `make verify`, `make verify-runtime`, `make verify-mvp`, `make generate`, `make schema`를 사용한다. 상세 절차는 [README](../../../README.md)에 있다.

Go 포맷은 현재 module 도구 체인의 gofmt를 사용한다. vet와 race test는 `./cmd/... ./internal/...`만 검사한다. `go mod verify`는 의존성 무결성을 검사한다. module 아래의 `node_modules`에 포함된 외부 Go 코드를 검사하지 않는다.

TypeScript는 `npm run check --prefix adapters`로 Prettier·ESLint·tsc를 순서대로 실행한다. 각 명령이 실패하면 통합 명령도 실패한다. build는 `noEmitOnError`를 유지한다. YAML은 설치된 Prettier로 검사한다. Compose 구조 검사는 더미 값을 명시적으로 덮어써 실제 로컬 비밀값을 출력하지 않는다.

sqlc v1.30.0으로 실제 singleton schema와 업무 transaction 쿼리를 생성한다. make generate 뒤 internal/database diff를 확인한다. make schema는 실제 migration에서 D08을 생성한다.

## 라이브러리 API 근거 — 2026-10-03

Context7 `resolve-library-id` 후 `query-docs`로 확인했다. 비공개 코드와 비밀값은 조회하지 않았다.

- `/jackc/pgx/v5.10.0`: pgxpool.New·Ping·Close와 stdlib의 `sql.Open("pgx", ...)`. [고정 버전 소스](https://github.com/jackc/pgx/tree/v5.10.0). 실제 설치 버전도 5.10.0이다.
- `/pressly/goose`: NewProvider·DialectPostgres·Up·ErrNoMigrations. Context7은 버전 없는 main 근거를 반환했다. 실제 3.28.0의 `provider.go`, `provider_options.go`, `provider_errors.go`를 로컬 module cache에서 대조했다. [3.28.0 공식 소스](https://github.com/pressly/goose/blob/v3.28.0/provider.go)와 `WithDisableGlobalRegistry(true)`의 실제 설치 API를 확인했다.
- `/typescript-eslint/typescript-eslint`: recommended flat config 배열. 버전별 Context7 근거가 없어 설치한 8.71.0의 peer 범위와 ESLint 10의 실제 lint 실행을 대조했다. [공식 flat config 문서](https://eslint.org/docs/latest/use/configure/configuration-files)도 확인했다.
- `/websites/sqlc_dev_en`: config version 2·Postgres·pgx/v5와 schema·query 필요 조건을 확인했다. [공식 configuration](https://docs.sqlc.dev/en/latest/reference/config.html)을 따른다. 이번 CLI는 v1.30.0이며 WithTx·:execrows·FOR UPDATE와 clock_timestamp timestamptz cast를 실제 생성·DB에서 검증한다.
- [Go 공식 배포 목록](https://go.dev/dl/?mode=json)에서 1.27.1을 확인했다. 실제 도구 체인과 Docker build 버전이 일치한다.
- [Compose profile 공식 문서](https://docs.docker.com/compose/how-tos/profiles/)에 따라 cloudflared를 선택 profile로 둔다. 실제 Compose config와 runtime을 별도로 검사한다.

Context7의 main 근거를 고정 버전 문서라고 해석하지 않는다. 필요한 API는 설치한 실제 버전의 컴파일·실행으로 대조한다.

## 검증 한계

운영 Tunnel은 연결하지 않는다. 현재 직접 검증은 개발 구성이며 기능 MVP의 수락을 뜻하지 않는다. 실제 도구 버전·종료코드·한계는 [실행 기록](../exec-plans/phases/SAR-SETUP-001-DEV.md)에 남긴다.

## SAR-MVP-001 추가 근거

Context7 `/cyberphone/json-canonicalization`의 Go Transform·RFC8785·duplicate key 근거를 확인했다.
실제 Go module은 v0.0.0-20241213102144-19d51d7fe467이다. 공개 공식 소스의 lone surrogate 거부를 테스트로 대조했다.
TypeScript는 Node 내장 crypto Ed25519와 UTF-16 key sort 및 JSON.stringify를 사용한다. 실제 Go signature와 TS 재검증이 HTTP 합성 검사에서 일치한다.
API 문서를 고정 버전 근거로 확대 해석하지 않는다. 설치 소스·컴파일·실행을 함께 사용한다.
