---
title: KnowsLink 프로젝트 기준
status: review
updated: 2026-10-06
owner: dev
tasks: [SAR-SETUP-001-DEV, FULLOPS-UPDATE-0.9.10, SAR-MVP-001-DEV, SAR-MVP-002-INSTALL-FIX-DEV, SAR-MVP-002-BOT-CATALOG-DEV, SAR-PUBLIC-IDENTITY-001-DEV, FULLOPS-UPDATE-0.9.14]
summary: 프로젝트 정본과 실제 lint·테스트 및 UI 검사 한계
---

# KnowsLink 프로젝트 기준

| 항목 | 값 |
|---|---|
| 제품 목적 | 승인된 에이전트 사이의 선택적 인간개입 전달; SAR-SETUP-001 D02 |
| 기본 브랜치·원격 | `main`, `origin` (`https://github.com/tourbut/KnowsLink.git`) |
| 현재 상태 | 로컬 합성 relay.v1·owner gate·shared inbox·TypeScript stub 후보와 일반 이메일 신원·세션·자기 owner 홈 후보(SAR-PUBLIC-IDENTITY-001); 독립 수락·운영 공개는 후속 |
| 기술 스택 | Go 1.27.1, pgx/v5 5.10.0, goose/v3 3.28.0, TypeScript 5.9.3, Node 22.22.2, Postgres 17, Compose |
| 기술 설계 정본 | `docs/design-docs/architecture.md`, `docs/design-docs/tech-stack.md` — D03, dev 담당 |
| 기획 정본 | `.fullops-squad/docs/planning/`, 사용자 경험은 `docs/design-docs/mockups/` — designer 담당 |
| 공통 개발 기준 | [rules/common/README.md](rules/common/README.md), `fullops-common-0.3.3`; Ponytail full |
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
- 로컬 Compose·실제 SQL migration·누락 설정·차단된 Tunnel·미설정 adapter: `make verify-runtime`.
- 개발 실행·중지: `docker compose --env-file .env.example up --build --wait relay`, `docker compose --env-file .env.example down`.
- 합성 업무 DB·경합·TypeScript/Go UI 연동: `make verify-mvp` (고유 Compose project, 자기 자원만 회수).
- 실제 sqlc 생성: `make generate`; 생성물 diff는 `git diff --exit-code -- internal/database`.
- D08 테이블 정의 생성: `make schema`.
- Grok plugin 패키지·실제 CLI 설치 검사: `make verify-grok-plugin` (임시 HOME의 `grok plugin install`·`grok mcp doctor`, 사용자 `~/.grok` 미변경). CLI 검사이며 앱 카탈로그 증거가 아니다.
- Grok Bot 앱 Command server 준비: Bot 컴퓨터에서 `sh scripts/install_bot_mcp.sh` (`KNOWSLINK_PREFIX` 기본 `/workspace/.knowslink`). 앱 등록은 owner가 [README](../adapters/README.md#grok-bot-앱-등록)대로 실행한다.
- 운영 배포·실제 Tunnel 연결: 후보 수락 뒤 OPS 담당.

## 공통 기준의 적용과 예외

기존 프로젝트 규칙은 없으므로 공통 규칙을 기본값으로 적용한다. 이후 기술 정본이 생기면 연결하며 보안·권한·리뷰 수락 기준은 낮추지 않는다.
변경한 동작과 실패·경계 조건을 검증한다. DEV 자동 검증과 독립 기능 QA를 구분한다. 업무 API·Go owner UI·실제 singleton SQL·sqlc 생성이 후보에 포함된다. 일반 이메일 코드 신원·세션·신원 한도는 SAR-PUBLIC-IDENTITY-001 후보다. agent 연결·관계·실메시지 한도·DEC-02 공개 정책·벤더 연결은 후속이다.
작업 지시서에는 적용 문서와 기준 SHA를 남기고 worker와 검토자가 같은 버전을 읽도록 한다.

## 검증 담당과 후속 인계

DEV는 변경 동작의 자동 검사·관련 회귀·필요한 짧은 실행 확인을 완료한다. DEV 완료와 제품 최종 수락은 구분한다.
독립 전체 QA는 tester가 안정된 고정 통합 후보에서 수행한다. 직접 시각 검수는 designer가 담당한다. 별도 ART 역할은 구성하지 않는다.
coor는 검사별 담당·대상 SHA·실행 시점·통과 조건과 후속 인계 조건을 지시서에 기록한다.
캡처는 지정 시각 항목에만 만든다. 영상은 정지 화면으로 판정할 수 없는 항목에만 만든다.
변경 없는 증거는 관련 의존성의 동일성을 확인하고 원래 실행 SHA·조건을 연결해 재사용한다. 새 SHA에서 실행한 결과로 표시하지 않는다.
재검증은 변경 영향·새 실패·증거 결함·미충족 조건이 있을 때 수행한다. 기존 실패·held·미해결 critical/high·제품 정지·최종 플랫폼과 사람 평가 기준은 유지한다.
보류 항목에는 담당과 재개 조건을 남긴다. 상세 반복 범위는 [공통 테스트 기준](rules/common/testing.md)을 따른다.

## FullOps 0.9.14 검사 연결

검사 정본은 루트 Makefile과 adapters/package.json 및 adapters/eslint.config.mjs다. `product-lint`는 `make lint`, kind `lint`, cwd `.`다. Go 포맷·vet·module 무결성, adapter Prettier·ESLint·타입 검사와 설정 검사를 유지한다. ESLint recommended의 미사용 검사도 유지한다. `product-test`는 `make test`, kind `test`, cwd `.`다. Go race와 adapter 테스트를 같은 HEAD에서 실행하고 FullOps 결과에 종료코드를 남긴다.

현재 UI는 Go template의 일반 CSS이며 Tailwind v4·React·Vue·Svelte·shadcn 설정은 없다. shadcn 설치는 해당 없음이다. 서비스의 별도 디자인 lint·공용 테마 전환은 미구성이다. DESIGN 기본 규칙은 Go 문자열 안의 CSS를 검사하지 않으므로 실제 UI 시각 검수와 제품 기준을 대신하지 않는다. 후속 UI 변경 담당 dev가 정본·공용 컴포넌트·예외와 적용 가능한 검사 또는 미실행 사유를 정규 지시서에 기록한다. 이 운영 업데이트는 제품 UI를 변경하지 않는다.
