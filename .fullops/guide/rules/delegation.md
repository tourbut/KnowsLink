---
title: 하위 에이전트 활용 기준
status: draft
updated: 2026-11-19
owner: coor
tasks: [FULLOPS-UPDATE-1.3.0, FULLOPS-UPDATE-1.24.0]
summary: 선택형 위임의 범위와 부모의 완료 책임을 정한다.
---

# 병렬 진행

나눌 수 있는 일은 하위 에이전트로 동시에 진행해 구현을 빨리 끝낸다. 이 문서는 lead와 worker가 함께 읽는 병렬 진행의 정본이다. 방침은 `.fullops/fullops.toml`의 `parallel` 하나이고 `fullops status`의 첫머리와 배정 spec의 `Parallel:`에 같은 값이 나온다.

## 방침

| 값 | 뜻 | 호스트의 하위 에이전트 도구 | 선택형 Orca 위임 |
|---|---|---|---|
| `off` | 하위 에이전트와 하위 worker를 쓰지 않고 직접 처리한다 | hook이 거부한다 | `fullops dispatch`가 거절한다 |
| `read` | 읽기 전용 조사·리뷰에만 쓴다 | 지시에 `Purpose: research` 또는 `Purpose: review`와 `read only`가 있어야 통과한다 | 조사와 리뷰만 되고 읽기 전용으로 간다 |
| `on`(기본) | 나눌 수 있는 구현과 조사를 하위 에이전트로 동시에 진행하도록 권한다 | 통과한다 | 통과한다 |
| `ultra` | `on`에 더해 Claude Code에서 울트라코드(Workflow)를 쓰게 한다 | 통과한다 | 통과한다 |

- 방침을 바꾸는 것은 사람이다. `fullops.toml`의 `parallel`을 고쳐 커밋한다. 에이전트가 스스로 올리지 않는다.
- 한 과제만 끄려면 lead가 지시서의 `제약·협업·후속`에 적는다.
- 동시에 돌리는 수는 한 번에 4개 이하로 시작한다. 실제 상한은 호스트가 정한다.
- `on`과 `ultra`는 사용량이 빨리 준다. 읽기 전용만 허용하려면 `read`, 끄려면 `off`로 고친다.

## 수단과 고르는 기준

| 일의 모양 | 수단 | 변경이 쌓이는 곳 | 결과를 받는 쪽 |
|---|---|---|---|
| 한 세션 안에서 끝나는 탐색·조사, 서로 겹치지 않는 파일 묶음의 구현 | 호스트의 하위 에이전트 도구: Claude Code `Agent`, Codex `spawn_agent`, grok `spawn_subagent`, agy `invoke_subagent` | 부모의 워크트리 | 부모 세션 |
| 같은 모양인데 조각이 많거나 결과를 서로 검증해야 하는 일. Claude Code이고 울트라코드가 켜져 있을 때만 | 울트라코드: Workflow가 여러 하위 에이전트를 조율한다 | 부모의 워크트리 | 부모 세션 |
| 다른 모델·호스트가 맡아야 하는 부분 작업, 하위 에이전트 도구를 쓸 수 없는 세션의 부분 작업 | Orca 하위 배정: worker가 `fullops dispatch <KEY> --brief "<맡길 파일 묶음과 완료 조건>"` | 부모의 워크트리 | 부모 세션(자기 Run의 `worker_done`) |
| 서로 다른 워크트리가 필요한 독립 과제 | lead의 역할 단위 배정: `PLANS.md` 행 → `fullops task new` → `fullops dispatch` | 역할마다 자기 워크트리와 브랜치 | lead(완료 보고 → `fullops merge`) |

- 하위 에이전트에게는 묶음마다 맡는 파일과 완료 조건을 적어 준다. 파일 묶음은 겹치지 않게 나눈다.
- 호스트의 워크트리 격리(하위 에이전트의 격리 옵션, Workflow의 격리된 사본)는 쓰지 않는다. 격리된 사본의 변경은 부모의 워크트리에 없어서 lint와 `fullops done`에 잡히지 않는다. 서로 다른 워크트리가 필요하면 역할 단위 배정이다.
- 같은 종류의 일을 동시에 여럿 돌리려면 역할(워크트리)을 늘린다(예: `dev1`, `dev2`). 역할당 인박스는 하나이고 과제 이름으로 된 지시서 파일은 만들지 않는다.
- 나눌 수 없는 일은 직접 한다.
- 리뷰와 테스트는 이 수단으로 늘리지 않는다. 독립 리뷰는 리뷰 판정이 요구할 때 lead가 한 번 배정하고, 테스트 범위는 테스트 레벨(`guide/rules/common/testing.md`)이 정한다.

## 울트라코드

- 울트라코드는 사용자가 켠다(`/effort ultracode`나 사용자가 직접 친 요청). 에이전트가 지시서나 spec에 키워드를 넣어 켜지 않는다.
- lead 세션: `parallel = "ultra"`이고 울트라코드가 켜져 있으면 조각이 많은 일을 Workflow로 조율한다. 켜져 있지 않으면 `Agent` 도구를 쓴다.
- Claude Code worker: `parallel = "ultra"`인 레포에서 `fullops dispatch`가 새 터미널로 띄우는 worker는 울트라코드로 뜬다. 있는 터미널에 다시 배정한 worker와 다른 호스트의 worker는 그 호스트의 하위 에이전트 도구를 쓴다.
- worker의 권한 모드가 auto나 bypass가 아니면 Workflow 실행마다 승인 창에서 멈춘다. 그 PC에서는 `ultra`가 사람 없이 이어지지 않는다. 권한 모드는 사용자가 정하고 fullops는 바꾸지 않는다.

## 하위 배정

배정받은 worker가 자기 과제의 일부를 다른 Orca worker에게 맡기는 길이다.

1. worker가 `fullops dispatch <자기 KEY> --brief "<맡길 파일 묶음과 완료 조건>"`을 실행한다. 다른 모델이 맡아야 하면 `--agent <CLI> --model <모델> --effort <effort>`를 함께 준다. 하위 worker는 부모와 같은 워크트리에서 뜬다. 인박스와 준비 커밋은 만들지 않는다.
2. 깊이는 한 단계다. 하위 worker는 다시 하위 배정을 하지 못한다. 다른 과제 키와 다른 역할은 하위 배정으로 맡기지 못한다. 그런 일은 `fullops done --followup "<역할>: <글>"`로 lead에게 넘긴다.
3. 하위 worker는 맡은 파일만 고친다. 커밋·lint·`fullops done`은 하지 않는다. 끝나면 바꾼 파일 목록과 확인한 것을 본문에 적어 Orca preamble의 `worker_done`을 한 번 보낸다.
4. 부모는 턴을 끝낸다. 하위 worker가 남아 있는 동안에는 커밋·lint·완료 보고 없이 턴을 끝낼 수 있다. Orca가 깨우면 `fullops wait --run <자기 Run>`으로 보고를 받는다.
5. 부모가 결과를 직접 확인해 커밋하고 `fullops wait --run <자기 Run> --ack <메시지 ID> --no-wait`로 처리 표시를 한다.
6. 부모가 lint를 통과시키고 `fullops done`을 한다. 하위 worker가 남았거나 받지 않은 보고가 있으면 `fullops done`이 거절한다. lead가 받는 완료 보고 하나가 하위 worker의 변경까지 담고, 통합은 `fullops merge` 한 번이다.

## 게이트는 그대로다

- 하위 에이전트와 하위 worker의 변경은 부모의 워크트리에 쌓인다. 부모가 커밋하고 `fullops lint`를 통과시켜야 끝낼 수 있다.
- 완료 보고는 부모 하나가 한다. 병합 전 검사(리뷰 판정, lint 증거)는 보고 SHA 하나에 대해 그대로 한다. 병렬로 만든 변경이 리뷰 기준을 넘으면 리뷰 판정에 걸린다.
- 하위 에이전트의 도구 호출도 부모와 같은 hook을 탄다. 날것 배정, 인박스 직접 생성, 설계 역할의 코드 수정은 하위 에이전트에서도 거부된다.
- 같은 워크트리에서 여러 에이전트가 동시에 고치므로 파일 소유는 지시로만 나뉜다. 겹치면 부모가 lint와 직접 확인으로 잡는다.
