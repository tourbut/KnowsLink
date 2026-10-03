---
title: KnowsLink 프로젝트 기준
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-SETUP-001-DEV]
summary: 잠긴 기술 선택과 초기 제품 구성의 경로·검증 명령을 정의한다
---

# KnowsLink 프로젝트 기준

| 항목 | 값 |
|---|---|
| 제품 목적 | 승인된 에이전트 사이의 선택적 인간개입 전달; SAR-SETUP-001 D02 |
| 기본 브랜치·원격 | `main`, `origin` (`https://github.com/tourbut/KnowsLink.git`) |
| 현재 상태 | 초기 Go relay·별도 migrate·TypeScript adapter·Compose·제품 lint 구성; 업무 MVP는 후속 |
| 기술 스택 | Go 1.27.1, pgx/v5 5.10.0, goose/v3 3.28.0, TypeScript 5.9.3, Node 22.22.2, Postgres 17, Compose |
| 기술 설계 정본 | `docs/design-docs/architecture.md`, `docs/design-docs/tech-stack.md` — D03, dev 담당 |
| 기획 정본 | `.fullops-squad/docs/planning/`, 사용자 경험은 `docs/design-docs/mockups/` — designer 담당 |
| 공통 개발 기준 | [rules/common/README.md](rules/common/README.md), `fullops-common-0.3.1`; Ponytail full |
| 보안·코딩 규칙 | [코딩](rules/common/coding-style.md), [테스트](rules/common/testing.md), [보안](rules/common/security.md) |
| 문서 언어 | 한국어 |
| 이슈 트래커 | 로컬 `PLANS.md`·역할 인박스. GitHub Issues/Projects 연동은 미설정 |

## 기술 기준

개발 설계와 구현의 책임자는 dev다. 실제 요구사항과 기존 패턴을 확인하고 표준 라이브러리·플랫폼 기능·설치된 의존성부터 사용한다.
제품 경로는 루트 Go module, `cmd/`, `internal/`, `adapters/`, `db/`, `sqlc.yaml`, `Dockerfile`, `compose.yaml`, `Makefile`, `scripts/`다. 초기 구성과 직접 검증은 dev 소유다. 독립 QA는 tester 소유다. 운영 배포는 후속 ops 과제다.
인증값은 Git 미추적 `.env`로 제공하며 값·전문을 문서와 로그에 기록하지 않는다. FullOps 운영용 `.fullops-squad/.env`는 루트 `.env`의 로컬 링크다.
제품 비즈니스 판단은 designer, 기술 판단은 dev, 통합·배포는 ops, 독립 검증은 tester가 맡는다.

## 검증 명령

레포 루트에서 실행한다. `<플러그인>`은 해당 세션에 설치된 FullOps 패키지 경로를 조회해 사용하며 캐시 절대경로를 레포에 기록하지 않는다.

- 공백·패치 검사: `git diff --check` (스테이징 후 `git diff --cached --check`).
- FullOps lint: 깨끗한 커밋 상태에서 `python3 <플러그인>/scripts/lint.py --repo . --from <기준 SHA> --out <레포 밖 검증 결과 경로>`.
- 산출물 검사: `python3 <플러그인>/scripts/deliverables.py --repo . --strict`.
- 현황판 생성: `python3 <플러그인>/scripts/board.py --repo .`.
- 설치: `make install` (Go module download, npm ci).
- 제품 lint·format·타입·설정 검사: `make lint` (`cmd/`, `internal/`, adapter 소스·설정, Compose·sqlc YAML).
- unit/race test: `make test` (`./cmd/... ./internal/...`, 외부 node_modules의 Go 코드는 제외).
- Go·TypeScript 빌드: `make build`.
- lint 위반과 등록 명령의 실패 전파: `make verify` (임시 복제본).
- 로컬 Compose·DB·빈 migration no-op·누락 설정·차단된 Tunnel·adapter: `make verify-runtime`.
- 개발 실행·중지: `docker compose --env-file .env.example up --build --wait relay`, `docker compose --env-file .env.example down`.
- 운영 배포·실제 Tunnel 연결·업무 SQL 적용·sqlc 생성: 범위 밖 또는 업무 SQL 부재로 미적용.

## 공통 기준의 적용과 예외

기존 프로젝트 규칙은 없으므로 공통 규칙을 기본값으로 적용한다. 이후 기술 정본이 생기면 연결하며 보안·권한·리뷰 수락 기준은 낮추지 않는다.
변경한 동작과 실패·경계 조건을 검증한다. 초기 구성의 직접 테스트와 독립 기능 QA를 구분한다. `/healthz` 외 업무 API와 human-gate는 미구현이다. 실제 SQL이 없어 migration 적용과 sqlc 생성 성공을 주장하지 않는다.
작업 지시서에는 적용 문서와 기준 SHA를 남기고 worker와 검토자가 같은 버전을 읽도록 한다.
