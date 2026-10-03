---
title: Orca 역할 배정
status: draft
updated: 2026-10-03
owner: coor
tasks: [SAR-SETUP-001, FULLOPS-UPDATE-099]
summary: 역할별 책임과 Astra를 제외한 모델 후보 및 라우팅 기준을 정의한다
---

# Orca 역할 배정

역할별 브랜치는 고유하게 유지한다. 책임 전환에서는 기존 브랜치와 워크트리를 재사용할 수 있다. `fullops.json`이 역할·브랜치 정본이다.
Orca 계층은 `coor` 아래 `designer`, `dev`, `ops`, `tester`다. 기본 브랜치 세션은 초기 설정용이며 운영은 coor에서 진행한다.

| 역할 | 책임·담당 경로 | 브랜치 | CLI | 모델 | 기동 방식 |
|---|---|---|---|---|---|
| coor | 과제 분해·배정·Run·진행 관리; `PLANS.md`, `handovers/`, `board/` | `fullops/coor` | `codex` | 모델 후보 | 상설 워크트리; 첫 요청 때 coordinator 세션 시작 |
| designer | 제품 기획·요구사항·사용자 경험; `docs/planning/`, `docs/design-docs/mockups/` | `fullops/designer` | `codex` | 모델 후보 | 기획 과제 발생 시 `worker-start --run` |
| dev | 기술 설계·구현·직접 검증; `docs/design-docs/`의 기술 문서, 향후 제품 코드·에셋 경로 | `fullops/dev` | `claude` 또는 `codex` | 모델 후보 | 기술 설계·구현 과제 발생 시 `worker-start --run` |
| ops | 배포·통합·운영; `docs/operations/`, 향후 배포 설정 경로 | `fullops/ops` | `claude` | 모델 후보 | 운영 과제 발생 시 `worker-start --run` |
| tester | 재현·테스트·회귀 검증; `docs/evaluations/scenarios/`, `docs/evaluations/qa-reports/`, 향후 테스트 경로 | `fullops/tester` | `claude` | 모델 후보 | 구현 SHA 준비 후 `worker-start --run` |

문서 경로는 `.fullops-squad/` 기준이다. 실제 제품 코드·에셋·배포·테스트 경로는 첫 개발 과제에서 dev가 기술 스택과 함께 정한다.
원격은 `origin`, 기준 브랜치는 `main`이다. 상설 워크트리만 구성하고 이번 setup에서는 에이전트 세션을 시작하지 않는다.

## 모델 후보

사용자가 제공한 후보의 약한 것부터 강한 순서를 유지한다. dev CLI는 선택한 후보에 따른다.
`worker-start --agent <CLI> --model <모델> --effort <effort> --run <run id>` 결과의 `launch.effective`로 실제 적용을 확인한다.

- `coor` `codex` `gpt-6-luna` `xhigh`: 과제 라우팅·Run 운영. 사용자 확정 배정
- `designer` `codex` `gpt-6.1-sol` `medium`: 단일 기능 기획, 단순 요구사항 명확화
- `designer` `codex` `gpt-6.1-sol` `high`: 여러 사용자 흐름 기획, 요구사항 수락 검토
- `dev` `claude` `claude-sonnet-5-5` `medium`: 한 줄 문구·수치·설정 변경
- `dev` `claude` `claude-sonnet-5-5` `high`: 기존 패턴을 따르는 단일 파일 소규모 구현·버그 수정
- `dev` `claude` `claude-opus-5-5` `medium`: 기존 패턴을 따르는 일반 기능 구현·검사 추가
- `dev` `codex` `gpt-6.1-sol` `medium`: 여러 스크립트에 걸친 기능 구현, 모듈 연결
- `dev` `claude` `claude-opus-5-5` `high`: 원인 추적이 필요한 버그, 여러 모듈 구현
- `dev` `codex` `gpt-6.1-sol` `high`: 공유 계약을 따르는 큰 기능 구현·리팩터링
- `ops` `claude` `claude-sonnet-5-5` `high`: 기록 검사·형식 검증·여러 SHA의 충돌 없는 main 통합 조정
- `tester` `claude` `claude-sonnet-5-5` `medium`: 단순 재현 단계 검증·짧은 수동 테스트
- `tester` `claude` `claude-sonnet-5-5` `high`: 여러 시나리오의 테스트 코드 작성·회귀 검증
- `tester` `claude` `claude-opus-5-5` `medium`: 모호한 재현 조건 분석·여러 모듈의 테스트 설계

## 라우팅 기준

- coordinator 역할: `coor`
- 설계 역할: `designer`
- tester 역할: `tester`
- 제품 기획 역할: `designer`
- 기술 계획 역할: `dev`

- `coor`: 요청 접수, 과제 분해·배정, Run·진행·복귀 주소 관리. 제품 판단은 designer, 기술 판단은 dev에게 전달한다.
- `designer`: 제품 기획, 요구사항, 사용자 경험, 제품 범위·수락 기준을 맡는다. 기술 설계 요청은 dev에 넘길 수 있도록 제품 요구사항을 명확히 한다.
- `dev`: 개발에 대한 설계와 구현을 함께 맡는다. 기술 스택·아키텍처·API·데이터 모델·모듈 설계·제품 코드·직접 검증이 담당 범위다.
- `ops`: 배포·운영·수락된 SHA의 기계적 통합·동기화. 제품 의미 충돌은 designer, 기술 충돌은 dev로 돌린다.
- `tester`: 재현 조건·테스트 설계·테스트 코드·회귀 근거. 제품 코드 수정이 필요하면 coor를 통해 dev에 요청한다.

FullOps의 설계 역할 표시는 제품 기획 전용 designer를 뜻한다.
`product`는 제품 목표·규칙·수치·화면과 아트 방향·우선순위·완료 조건의 결정이다. designer가 맡는다.
`implementation`은 기존 제품 요구 안의 기술 계획·구조/API·버그 분석·구현·테스트·기술 문서 갱신이다. dev가 같은 과제에서 맡는다.
기술 난도가 높아도 기술 계획만을 이유로 designer에게 배정하지 않는다. designer는 상시 기술 승인자가 아니다.
`unresolved`는 책임 확인 전 보류한다. coor는 실제 소유 책임을 확인하고 재선정 근거를 라우팅 기록에 남긴다.
제품 규칙 변경·범위 확대·공유 제품 기준 질문은 coor를 통해 designer에게 전달한다. 기술 질문은 담당 구현자에게 전달한다.
ops는 배포·통합·운영의 기술 계획과 검증을 맡는다. tester는 독립 동작 QA를 맡고 제품 코드는 수정하지 않는다.
새 독립 코드 리뷰는 구현자와 다른 검토자의 별도 세션에서 수행한다. 검토자는 고정 SHA의 깨끗한 detached snapshot을 읽기 전용으로 사용한다.
리뷰 결과와 보고서는 기록 체크아웃에 작성한다. 별도 세션 ID·snapshot 경로·head·read_only를 기록한다.
필요한 직접 시각 검수와 독립 동작 QA를 수행한다. 미해결 critical/high가 있으면 수락·병합을 차단한다.
전환 전 리뷰 기록은 원본을 보존한다. 구현자 자기 리뷰로 독립 리뷰를 대체하지 않는다.
모델 후보는 최신 coor 브랜치의 사용자 결정(Astra 제외)을 보존했다. 실제 CLI 지원 여부는 기동 때 확인한다.
