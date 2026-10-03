---
title: Orca 역할 배정
status: draft
updated: 2026-10-03
owner: coor
tasks: [SAR-SETUP-001, FULLOPS-UPDATE-099, SAR-SETUP-001-DEV, FULLOPS-UPDATE-0.9.10]
summary: 역할별 책임과 Astra를 제외한 모델 후보 및 라우팅 기준을 정의한다
---

# Orca 역할 배정

역할별 브랜치는 고유하게 유지한다. 책임 전환에서는 기존 브랜치와 워크트리를 재사용할 수 있다. `fullops.json`이 역할·브랜치 정본이다.
Orca 계층은 `coor` 아래 `designer`, `dev`, `ops`, `tester`다. 기본 브랜치 세션은 초기 설정용이며 운영은 coor에서 진행한다.

| 역할 | 책임·담당 경로 | 브랜치 | CLI | 모델 | 기동 방식 |
|---|---|---|---|---|---|
| coor | 과제 분해·배정·Run·진행 관리; `PLANS.md`, `handovers/`, `board/` | `fullops/coor` | `codex` | 모델 후보 | 상설 워크트리; 첫 요청 때 coordinator 세션 시작 |
| designer | 제품 기획·요구사항·사용자 경험; `docs/planning/`, `docs/design-docs/mockups/` | `fullops/designer` | `codex` | 모델 후보 | 기획 과제 발생 시 `worker-start --run` |
| dev | 기술 설계·구현·직접 검증; `docs/design-docs/` 기술 문서, `cmd/`, `internal/`, `adapters/`, `db/`, `scripts/`, 루트 Go·sqlc·Make·Docker·Compose 설정 | `fullops/dev` | `claude` 또는 `codex` | 모델 후보 | 기술 설계·구현 과제 발생 시 `worker-start --run` |
| ops | 배포·통합·운영; `docs/operations/`, 향후 배포 설정 경로 | `fullops/ops` | `claude` | 모델 후보 | 운영 과제 발생 시 `worker-start --run` |
| tester | 재현·테스트·회귀 검증; `docs/evaluations/scenarios/`, `docs/evaluations/qa-reports/`, 향후 테스트 경로 | `fullops/tester` | `claude` | 모델 후보 | 구현 SHA 준비 후 `worker-start --run` |

문서 경로는 `.fullops-squad/` 기준이다. 제품 경로는 레포 루트 기준이다. SAR-SETUP-001-DEV의 초기 구성·직접 검증 파일은 dev 소유다. 독립 QA의 시나리오·보고서는 tester 소유다. 운영 배포의 파일 소유권은 후속 ops 지시서에서 정한다.
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

## 검증 담당과 후속 인계

DEV는 변경 동작의 자동 검사·관련 회귀·필요한 짧은 실행 확인을 완료한다. DEV 완료와 제품 최종 수락은 구분한다.
독립 전체 QA는 tester가 안정된 고정 통합 후보에서 수행한다. 직접 시각 검수는 designer가 담당한다. 별도 ART 역할은 구성하지 않는다.
coor는 검사별 담당·대상 SHA·실행 시점·통과 조건과 후속 인계 조건을 지시서에 기록한다.
캡처는 지정 시각 항목에만 만든다. 영상은 정지 화면으로 판정할 수 없는 항목에만 만든다.
변경 없는 증거는 관련 의존성의 동일성을 확인하고 원래 실행 SHA·조건을 연결해 재사용한다. 새 SHA에서 실행한 결과로 표시하지 않는다.
재검증은 변경 영향·새 실패·증거 결함·미충족 조건이 있을 때 수행한다. 기존 실패·held·미해결 critical/high·제품 정지·최종 플랫폼과 사람 평가 기준은 유지한다.
보류 항목에는 담당과 재개 조건을 남긴다. 상세 반복 범위는 [공통 테스트 기준](rules/common/testing.md)을 따른다.

## 완료 결과의 통합과 역할 인박스 — FullOps 0.9.12

실행 지시서는 역할별 `handovers/to_<role>.md` 하나로 고정한다. 다른 과제로 사용 중인 인박스를 덮어쓰지 않는다. 다음 과제는 PLANS.md에 대기시킨다. 과제명 파일·pending·logs는 참조 자료이며 배정 지시서로 사용하지 않는다. 역할 작업 완료 시 `work.py finish`로 지시서와 완료 보고 전문을 logs에 보존하고 인박스를 비운다.

coor는 worker_done을 받으면 고정 SHA의 필수 검토·검증을 확인하고 실제 기본 브랜치 main에 병합한다. 원격 origin/main에 일반 push한 뒤 완료 SHA의 로컬·원격 조상 관계를 확인한다. 역할 브랜치 push만으로 통합을 완료하지 않는다. 전체 제품 수락과 개별 문서 결과의 통합은 구분한다. 미해결 critical/high와 필수 검증 실패는 계속 차단한다.

병합 뒤 coor를 포함한 모든 등록 역할의 실제 worker 상태와 작업 트리를 확인한다. 쉬고 있으며 깨끗한 역할만 최신 main으로 동기화한다. 진행 중·미커밋 변경·상태 확인 불가인 역할은 PLANS.md에 최신 기본 SHA와 예약을 남긴다. 다음 dispatch 전에 최신 로컬·원격 main 포함을 확인한다. 진행 중 체크아웃과 사용자 자료를 덮어쓰지 않는다.

완료 메시지는 Git 공용 디렉터리의 fullops-integration에 보존된다. `integration.py --repo <레포> status`로 미통합 결과를 확인한다. 검수 대기·실패·충돌·원격 오류·사용자 제한은 PLANS.md에 메시지 ID·SHA·사유·담당·재개 조건을 기록한다. 같은 내용을 `integration.py hold`로 남긴다. 조건 충족 시 resume하고 병합·push를 이어간다. 성공 결과를 통합한 뒤 release·ack하고 다음 독립 과제를 배정한다. 사용자의 현재 과제 완료 뒤 중지 지시는 유지한다.

기존 활성 과제명 지시서는 작업 중 이동하지 않는다. 해당 과제 완료 뒤 logs에 보존하고 다음 과제부터 정규 역할 인박스를 사용한다.
