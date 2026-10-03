---
title: KnowsLink 현재 계획
status: draft
updated: 2026-10-03
owner: coor
tasks: [FULLOPS-UPDATE-098, FULLOPS-UPDATE-0.9.10]
summary: 운영 계획과 FullOps 적용 및 역할 동기화 보류를 관리한다
---

# KnowsLink 현재 계획

## 운영

- orchestration Run: `run_8ca8bc058ab7`
- Run 목적: FullOps 초기 구성과 첫 개발 요청 준비.
- 현재 Run은 setup 세션에 연결돼 있다. 첫 coor 세션에서 `orca-ide orchestration run-use --id run_8ca8bc058ab7 --json`으로 연결한다.
- 운영 구조: coor 아래 designer / dev / ops / tester. 에이전트 세션은 첫 요청 전 시작하지 않는다.
- 완료: FullOps 0.9.7 설치(Codex·Claude Code·grok), 다섯 역할 setup, 원격 역할 브랜치 생성, 설정·검증 기준 작성.
- 완료: setup 커밋 `0c8da8d`를 main과 다섯 역할 브랜치에 push. Orca에서 coor를 부모로 네 worker 워크트리를 구성하고 동일 setup을 반영했다.
- 사용자 승인: 기존 미추적 `.gitignore`, `orca-agents.md`가 있는 상태에서 setup 진행 승인.

## 첫 요청 대기

- 제품 목적·요구사항·사용자 경험: designer가 첫 요청에서 기록.
- 기술 스택·기술 설계·제품 경로·lint/test/build 명령: dev가 첫 개발 과제에서 확정.
- 배포 대상·운영 환경: ops가 배포 과제에서 확정.
- 재현·테스트·회귀 기준: tester가 구현 완료 SHA를 대상으로 확정.
- 진행 중 제품 과제·worker Dispatch: 없음.
- 역할별 모델 후보는 사용자 참고 파일에서 유지했다. 실제 지원 여부는 과제 기동 시 확인한다.

## 워크트리·검증 근거

| 역할 | 워크트리 | 브랜치 | Orca 부모 |
|---|---|---|---|
| coor | `/home/shin/orca/workspaces/KnowsLink/fullops-coor` | `fullops/coor` | 없음 |
| designer | `/home/shin/orca/workspaces/KnowsLink/fullops-designer` | `fullops/designer` | coor |
| dev | `/home/shin/orca/workspaces/KnowsLink/fullops-dev` | `fullops/dev` | coor |
| ops | `/home/shin/orca/workspaces/KnowsLink/fullops-ops` | `fullops/ops` | coor |
| tester | `/home/shin/orca/workspaces/KnowsLink/fullops-tester` | `fullops/tester` | coor |

- Git 공백 검사와 JSON·역할·인박스·컨텍스트·모델 후보·라우팅·Run 정합성 검사 통과.
- `deliverables.py --strict`: 검사 13 / 미작성 13 / 문제 0 / 경고 0.
- setup 커밋 lint: ERROR 0 / WARNING 2 / 실행 불가 0. 신규 lint 설정(LINT-001), 기존 프로젝트 lint 도구 없음(LINT-000).
- 모든 워크트리의 브랜치·upstream·동일 HEAD·깨끗한 상태와 .env 링크 확인. 에이전트 세션은 기동하지 않았다.
- Jev OPENROUTER_API_KEY 존재 확인. 값은 출력하지 않았으며 실제 API 호출은 이번 setup에서 수행하지 않았다.
- 현황판: `.fullops-squad/board/index.html`. 생성 데이터는 Git 미추적이며 각 워크트리에서 board.py로 생성한다.

## FullOps 0.9.8 업데이트 — FULLOPS-UPDATE-098

- 업데이트 전 기준 커밋: `3c1a185`. 기본 브랜치는 main이며 기존 변경은 없었다.
- Codex 설치 버전은 0.9.8이다. `deps.py --check`에서 필수 CLI를 모두 확인했다.
- 진행 중 worker는 없었다. 다섯 역할의 인박스는 비어 있었다. 역할 워크트리는 모두 깨끗했다.
- 새 문서 작성 규칙을 추가했다. 하네스 지도에 규칙 링크를 추가했다. lint 안내에 DOC-003을 등록했다.
- 이번에 수정한 일반 문서는 `deliverables.py --stamp --path`로 등록했다. 기존 역할·모델·인박스·외부 공통 규칙은 유지했다.
- 기존 문서는 다음 수정 때 메타데이터를 등록한다. 진행 중 지시서의 일괄 변환은 보류한다.
- 동기화 대상은 coor, designer, dev, ops, tester다. 이 역할들은 쉬고 있으며 작업 트리가 깨끗하다. 예약된 동기화는 없다.
- 이 세션은 업데이트 전용이다. 새 coor 세션에서 Run `run_8ca8bc058ab7`을 연결한다.
- Claude Code와 grok은 이전 설치 버전 0.9.7이다. 해당 CLI로 worker를 시작하기 전에 README의 업데이트 명령을 실행한다.

## FullOps 0.9.9 업데이트 — FULLOPS-UPDATE-099

- 준비 기준은 `00b4cb34ae6e9f9fbc0b733ecaa3a2095fbc88eb`이다. 현재 Codex 플러그인은 0.9.9이며 필수 의존성을 확인했다.
- 제품 기획 designer / 기술 계획 dev의 합의된 책임 분리를 marker로 활성화했다. coor·설계·tester와 역할 브랜치 매핑은 유지했다.
- 이전 첫 요청 대기 기록은 당시 이력이다. 현재 SAR-SETUP-001 진행 자료는 coor·designer·dev 브랜치에 있다.
- ops·tester만 쉬고 깨끗한 동기화 대상이다. coor·designer·dev는 기존 SAR 인계와 작업을 보존하기 위해 완료 후 동기화한다.
- 새 coordinator는 [갱신과 인계 기록](docs/exec-plans/phases/FULLOPS-UPDATE-099.md)을 읽고 기존 Run을 연결한다.

## FullOps 0.9.10 업데이트 — FULLOPS-UPDATE-0.9.10

- 설치 버전은 갱신 전후 0.9.10이다. 레포 적용 버전 0.9.9의 미적용 릴리스를 회수한다.
- 공통 기준 0.3.2, 핸드오버와 시나리오 기준, 업데이트 연결과 역할별 검증 인계를 반영한다.
- DEV 관련 회귀와 최종 수락을 구분한다. tester는 안정된 통합 후보의 독립 전체 QA를 맡고 designer는 직접 시각 검수를 맡는다.
- 증거는 원래 SHA·조건과 의존성 동일성을 확인해 재사용한다. 기존 실패와 제품 정지는 유지한다.
- 모든 역할 워크트리의 동기화는 보류한다. coor가 진행 과제·활성 세션·깨끗한 상태를 확인한 뒤 준비 커밋을 전달한다.
- 상세 적용·검증·보류 담당과 재개 조건은 [업데이트 기록](docs/exec-plans/phases/FULLOPS-UPDATE-0.9.10.md)을 따른다.

### 하위 워크트리 적용 — 2026-10-03

사용자 요청으로 FullOps 0.9.10 준비 커밋을 이 역할 브랜치에 병합했다. 위 동기화 보류는 당시 기록이다.
기존 제품 자료·지시서·실패 기록과 역할별 프로젝트 기준을 보존했다. 실행 중 세션은 다음 시작 시 새 플러그인 규약을 읽는다.
