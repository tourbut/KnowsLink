---
title: dev 컨텍스트
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-SETUP-001-DEV, SAR-MVP-001-DEV, SAR-MVP-002-DEV]
summary: 합성 MVP의 권한 경계와 검증 및 후속 수락 조건을 기록한다
---

# dev 컨텍스트

결정·교훈을 항목당 3줄 이내로 기록한다.

- 2026-10-03 SAR-SETUP-001-DEV: Go 1.27.1·pgx/v5·goose SQL-only·TypeScript adapter·Postgres 17.11·Compose의 초기 구성을 구현했다. 업무 MVP와 운영 연결은 후속이다.
- `make lint`와 `make verify`는 제품 위반과 실패 전파를 검출한다. Go 검사 대상은 `cmd/`, `internal/`로 한정해 node_modules의 외부 Go 코드를 제외한다.
- 코드 체크포인트 `929832aa0ecd`의 깨끗한 clone과 실제 로컬 runtime 검증을 통과했다. 빈 SQL no-op은 migration 적용 성공과 구분한다. 독립 QA·리뷰는 coordinator가 후속 배정한다.
- 상세 근거: [실행 기록](../docs/exec-plans/phases/SAR-SETUP-001-DEV.md). FullOps 기준 commands의 부재와 제품 직접 검사 결과를 분리해서 보고한다.

- 2026-10-03 SAR-MVP-001-DEV: shared Postgres singleton lock·epoch CAS로 합성 등록/pairing/relay/gate/result를 연결했다. 공개 규모의 성능과 실제 신원 인증은 후속이다.
- claim은 한 번만 발급한다. 재시작·철회·세대 교체 후 receipt나 과거 gate로 실행권을 복구하지 않는다. 승인 후에도 disclosure/stub은 차단한다.
- 검증·API·QA/UI 경로와 보류: [SAR-MVP-001-DEV 실행 기록](../docs/exec-plans/phases/SAR-MVP-001-DEV.md). PLANS/board와 독립 QA/UI/리뷰는 coor 후속이다.
- 2026-10-03 SAR-MVP-001-DEV 보안 후속: 메시지에 `deliver` 경로를 저장하고 agent transport는 `deliver:agent`만 처리한다. 직접 human inbox가 없으므로 H 외 `deliver:human` send는 403이다.
- 2026-10-03 RF-01: 경로 검사는 transport뿐 아니라 이미 발급된 claim을 다시 쓰는 부모 경계(`parentRouting`)에도 둔다. 저장 형식이 바뀌면 이전 메서드로 만든 직렬화 상태를 fixture로 남겨 회귀한다.
- 봉투 필드로 처리 주체가 갈리면 저장 상태에 경로를 남기고 lease·persist·ACK·claim 각 확정 지점에서 검사한다. intent 예외 목록에 기대지 않는다.
- 2026-10-03 SAR-BETA-001-REVIEW-FINAL: 보안 게이트를 Python `assert`로 쓰면 `PYTHONOPTIMIZE`에서 모두 사라진다. 운영 게이트 리뷰에서는 `-O` 실행과 증거 파일의 미래 mtime·ID 결합을 직접 시험한다.
- 2026-10-03 SAR-BETA-001-REVIEW-N1: 수정 검증은 negative 행렬을 일반·`-O`·`-OO`·`PYTHONOPTIMIZE=1`로 돌리고, 이전 SHA 파일을 대조군으로 같은 행렬에 넣어 검출력을 먼저 증명한다. 게이트만 고치면 같은 `assert` 패턴의 사후 증명(`verify.py`)이 남는다.

- 2026-10-03 SAR-MVP-002-DEV: 사용자가 공식 Grok Bot을 확정했다. Cursor manifest·stdio MCP·skill·standalone bundle을 준비하며 actual connection은 held다.
- 공통 Adapter와 direct-run CLI를 분리해야 bundle의 stdout에 CLI 출력이 섞이지 않는다. 기본 held와 synthetic-loopback·redirect 차단을 실제 MCP handshake로 검사한다.
- 지원 근거·package·검증·owner/admin 설치·독립 리뷰/QA 후속은 [실행 기록](../docs/exec-plans/phases/SAR-MVP-002-DEV.md)에 보존한다.
