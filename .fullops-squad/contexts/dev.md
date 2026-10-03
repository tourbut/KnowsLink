---
title: dev 컨텍스트
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-SETUP-001-DEV]
summary: 초기 구성의 기술 결정과 재현 검증·후속 수락 경계를 기록한다
---

# dev 컨텍스트

결정·교훈을 항목당 3줄 이내로 기록한다.

- 2026-10-03 SAR-SETUP-001-DEV: Go 1.27.1·pgx/v5·goose SQL-only·TypeScript adapter·Postgres 17.11·Compose의 초기 구성을 구현했다. 업무 MVP와 운영 연결은 후속이다.
- `make lint`와 `make verify`는 제품 위반과 실패 전파를 검출한다. Go 검사 대상은 `cmd/`, `internal/`로 한정해 node_modules의 외부 Go 코드를 제외한다.
- 코드 체크포인트 `929832aa0ecd`의 깨끗한 clone과 실제 로컬 runtime 검증을 통과했다. 빈 SQL no-op은 migration 적용 성공과 구분한다. 독립 QA·리뷰는 coordinator가 후속 배정한다.
- 상세 근거: [실행 기록](../docs/exec-plans/phases/SAR-SETUP-001-DEV.md). FullOps 기준 commands의 부재와 제품 직접 검사 결과를 분리해서 보고한다.
