---
id: D03
title: KnowsLink 초기 구성 아키텍처
status: review
updated: 2026-10-03
owner: dev
tasks: [SAR-SETUP-001-DEV]
upstream: [D02]
summary: 초기 relay·SQL migration·adapter·Compose의 경계와 요구사항 추적을 정의한다
---

# KnowsLink 초기 구성 아키텍처

## 범위와 원천

상위 요구사항은 [D02](../planning/product-specs/SAR-SETUP-001.md)다. 고정 원천은 [architecture.md](../planning/sources/silent-agent-relay/architecture.md), [decisions.md](../planning/sources/silent-agent-relay/decisions.md), [protocol.md](../planning/sources/silent-agent-relay/protocol.md)다. 검사 기준은 `729446d8da57`이며 규칙은 `fullops-common-0.3.1`이다.

이 설계는 초기 개발 구성만 설명한다. 전체 기능 MVP의 human-gate 포함 결정과 C1–C5는 유지한다. 업무 기능이 구현되었다고 주장하지 않는다.

## 실제 구성

| 경로 | 책임 | 현재 동작 |
|---|---|---|
| `cmd/relay` | Go 개발 relay | 설정 확인, pgxpool DB ping, `GET /healthz`, 종료 신호 처리 |
| `cmd/migrate` | 별도 one-shot migration | DB ping, goose SQL-only Up, 빈 SQL의 명시적 no-op |
| `internal/config` | 환경 변수 검증 | Postgres URL·listen 주소 확인, 값이 없는 안전한 오류 |
| `adapters/src` | TypeScript adapter 골격 | `unimplemented` 상태 출력 후 종료, 네트워크 연결 없음 |
| `db/migrations`, `db/queries` | 향후 업무 SQL | 빈 디렉터리, 가짜 스키마·쿼리 없음 |
| `sqlc.yaml` | 향후 query codegen | Postgres·pgx/v5, 출력 `internal/database` |
| `Dockerfile`, `compose.yaml` | 로컬 개발 패키징 | 네 서비스와 기동 의존 순서 |
| `Makefile`, `scripts/` | 직접 검증 | 제품 lint·test·build·위반 주입·로컬 runtime 검사 |

## 기동과 실패 경계

Compose는 `postgres`의 healthcheck 성공을 기다린다. `migrate`는 별도 `/app/migrate` 프로세스로 실행한다. `relay`는 migrate의 종료코드 0을 기다린다. relay는 API 기동 중 migration을 호출하지 않는다.

relay는 기동 시 DB ping에 실패하면 종료코드 1로 종료한다. `/healthz`는 요청마다 DB ping을 수행한다. DB가 응답하지 않으면 일반 문구와 HTTP 503을 반환한다. 잘못된 설정과 DB 오류에는 URL·비밀번호·내부 드라이버 오류를 노출하지 않는다.

relay는 업무 endpoint를 등록하지 않는다. ingest·approve·exec 요청은 HTTP 404다. TypeScript adapter는 pull 방향을 표시하되 polling이나 실행 권한을 제공하지 않는다.

goose provider는 전역 Go migration registry를 비활성화한다. migration 디렉터리의 Go 파일을 거부한다. 읽을 수 없는 디렉터리는 실패한다. 빈 SQL 디렉터리는 DB ping 뒤 명시적 no-op으로 종료한다. 업무 SQL을 추가한 뒤에는 SQL Up만 실행한다.

## 네트워크와 설정

Postgres는 internal `database` 네트워크에만 연결하고 호스트 포트를 게시하지 않는다. migrate도 이 네트워크에만 연결한다. relay는 `database`와 `ingress`에 연결한다. 호스트 게시 주소는 loopback이다.

cloudflared는 선택 `tunnel` profile이며 `ingress`에만 연결한다. 기본 실행에서는 비활성이다. token 없는 실행은 실패한다. 검증 시 네트워크를 차단해 실제 Tunnel을 연결하지 않는다. 운영 hostname `relay.knowslog.com`의 ingress 설정은 후속 운영 과제다.

비밀값은 실행 환경에서 전달한다. `.env.example`에는 더미 DB 값과 빈 token만 있다. Docker build context는 Go 코드·module 파일·migration 디렉터리로 제한한다. 실제 `.env`, FullOps 자료, TypeScript 의존성을 이미지에 보내지 않는다. 런타임은 UID/GID `65532`다.

HTTP timeout은 상태 확인 서버의 기술 기본값이다. frozen relay.v1의 TTL·lease·attempts·receipt 수치나 미결정 제품 rate/size 정책을 변경하지 않는다.

## 요구사항 추적

| D02 ID | 구현 | 직접 검증 |
|---|---|---|
| SETUP-01 | Go module·npm lock·고정 도구 버전 | install·build, 깨끗한 체크아웃 재현 |
| SETUP-02 | relay·별도 migrate·adapter·sqlc 설정 | unit test·로컬 DB 기동·no-op·adapter 실행 |
| SETUP-03 | Compose 네 서비스와 의존 순서 | `make lint-config`, `make verify-runtime` |
| SETUP-04 | 예시 값·필수 설정 실패·선택 Tunnel | config unit test·runtime 누락 검사·차단된 Tunnel 검사 |
| LINT-01 | `make lint` | 정상 종료코드 0 |
| LINT-02–03 | 등록된 통합 명령과 복제본 위반 주입 | `make verify`의 오류별 실패·원복 통과 |
| DOC-01 | README·project.md·D03 | 경로 대조·deliverables strict·FullOps lint |
| SCOPE-01 | health endpoint와 미구현 adapter | endpoint unit test·diff·원천 변경 없음 |
| QA-01 | coordinator가 tester에 완료 SHA 전달 | 후속 독립 QA, 직접 검증으로 대체하지 않음 |

## 후속과 검증 한계

업무 SQL이 없어 실제 SQL migration 적용과 sqlc 생성은 미적용이다. C1–C5·frozen relay.v1의 업무 경로, 네 어댑터 통합, human-gate UI, 운영 Tunnel과 배포는 후속이다. UI가 없어 직접 시각 검수 대상은 없다. UI 구현 때 직접 시각 검수를 수행한다.

검사 결과와 명령은 [실행 기록](../exec-plans/phases/SAR-SETUP-001-DEV.md)에 있다. 독립 리뷰·QA 및 미해결 critical/high 차단은 유지한다.
