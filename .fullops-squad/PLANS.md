---
title: KnowsLink 현재 계획
status: draft
updated: 2026-10-03
owner: coor
tasks: [FULLOPS-UPDATE-098]
summary: "역할 워크트리와 Run, 초기 구성 및 플러그인 업데이트 기록을 관리한다"
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

## SAR-SETUP-001 — 프로젝트 초기 구성과 lint

- 사용자 요청의 GitHub 기획 문서를 커밋 고정 스냅샷으로 확보했다.
- Jev API 키가 없어 design → designer, codex gpt-6-astra high 폴백을 기록했다.
- designer는 제품 요구사항과 설정 범위를 확인하고 dev 및 tester 지시서를 작성한다.
- 새 과제이므로 designer 새 세션을 시작한다. 기준 ref는 `00b4cb34ae6e9f9fbc0b733ecaa3a2095fbc88eb`이다.
- designer 배정: Run `run_8ca8bc058ab7`, Task `task_49c6e00e6470`, Dispatch `ctx_1f186b4db5b5`, terminal `term_6935611a-04a4-4427-ba3b-d92dd4cf247e`. codex gpt-6-astra high 적용과 착수를 확인했다.

- 사용자 요청으로 모델 후보에서 gpt-6-astra를 전부 제외했다. 이후 배정은 갱신된 후보만 사용한다. 이미 착수한 designer Dispatch는 완료를 기다리고 후속 배정에서 변경한다.
- 원천 스냅샷은 외부 원본 보존 규칙에 따라 front matter를 추가하지 않는다. 프로젝트 lint에서 해당 스냅샷 경로만 제외하고 작성 문서는 계속 검사한다.
