---
title: SAR-SETUP-001 — 설계 실행 기록
status: draft
updated: 2026-10-03
owner: designer
tasks: [SAR-SETUP-001]
summary: "원천 검토와 역할별 탐색, 설계 범위 및 검증 근거를 기록한다"
---

# SAR-SETUP-001 — 설계 실행 기록

## 대상과 결정

- 기준 ref: `00b4cb34ae6e9f9fbc0b733ecaa3a2095fbc88eb`; 시작 HEAD: `0cc35f0b2c83f781f57a001bd64cd41d6e31bab4`.
- 원천: `.fullops-squad/docs/planning/sources/silent-agent-relay/`, 외부 SHA `404ff834c0607055d63d2053bf7771d2f46ad3ae`.
- 적용: `fullops-common-0.3.1`, project.md, document-writing.md, fullops-work, fullops-deliverables, fullops-orca의 설계 회신 절차. 범위 최소화는 Ponytail full 기준을 따른다.
- 산출물: [D02](../../planning/product-specs/SAR-SETUP-001.md), [dev 지시서](../../../handovers/to_dev.md), [tester 지시서](../../../handovers/to_tester.md).
- 전체 MVP 대신 초기 개발 골격·설정·lint를 확정했다. dev가 잠기지 않은 버전·파일 배치·도구를 정한다.
- 원천의 기술 선택은 그대로 유지했다. 첫 기능 MVP의 human-gate 포함 결정도 유지했다.
- 이번 산출물은 D02다. D03은 dev가 작성하며 D01·D04–D13은 이번에 작성하지 않는다.
- 역할 지시서 키를 SAR-SETUP-001-DEV와 SAR-SETUP-001-TESTER로 나눴다. tester는 dev 완료 SHA 이후에만 착수한다.

## 탐색과 전제 확인

원천 README·product·architecture·protocol·decisions·mvp-checklist·business-model과 source.json을 로컬에서 읽었다.
요청의 상대 경로를 실제 하네스 내부 경로로 풀었다. 원천 스냅샷은 수정하지 않았다.
project.md의 스택 미정 상태는 원천 반영 전 상태다. dev 지시서에서 원천의 확정 선택을 우선하고 project.md를 갱신하도록 했다.
원천에 남은 과거 webhook 허용 메모보다 현재 protocol.md의 C1–C5 및 후반 결정이 우선한다. MVP webhook/evidence fetch는 OFF다.
원천 이미지 assets는 이 스냅샷에 없다. 문서 텍스트만 근거로 사용했으며 그림을 확인했다고 주장하지 않는다.

역할마다 다음 명령을 실행했다. `<플러그인>`은 세션에 설치된 FullOps 0.9.8 경로다.

```text
python3 <플러그인>/scripts/jev_find.py find --repo . --role dev --key SAR-SETUP-001-DEV --scope code
python3 <플러그인>/scripts/jev_find.py find --repo . --role dev --key SAR-SETUP-001-DEV --scope documents
python3 <플러그인>/scripts/jev_find.py find --repo . --role tester --key SAR-SETUP-001-TESTER --scope code
python3 <플러그인>/scripts/jev_find.py find --repo . --role tester --key SAR-SETUP-001-TESTER --scope documents
python3 <플러그인>/scripts/jev_context.py --repo . --role dev --key SAR-SETUP-001-DEV --paths <지시서의 후보 18개>
python3 <플러그인>/scripts/jev_context.py --repo . --role tester --key SAR-SETUP-001-TESTER --paths <지시서의 후보 18개>
```

모든 명령의 종료코드는 0이지만 Jev API 판정은 실패했다. find 오류는 `API or response validation failed: ValueError`, context 오류는 `API or response validation failed`다.
기존 route는 API 키 부재 fallback이다. 이번 find/context의 구체적인 실패 원인을 키 부재로 단정하지 않는다.
`rg --files --hidden .fullops-squad`와 Git 추적 목록으로 후보를 보완했다. context는 자동 필수 지시서까지 역할당 19개를 전부 keep했다.
architecture.md·mvp-checklist.md는 `sensitive or oversized passage`로 전송하지 않고 keep했다. 로컬 원문으로 확인했다.
API 실패를 자동 관련성·충돌 검증 통과로 보고하지 않는다. 코드 존재 판정도 null이며 absent가 아니다.
결과 JSON은 `docs/evaluations/jev/SAR-SETUP-001-{DEV,TESTER}-{find,documents-find,context}.json`에 역할별로 보존했다.

## 검증

- `deliverables.py --repo . --strict`: 종료코드 0. 검사 13, 미작성 12, 문제 0, 경고 0.
- `git diff --check`: 종료코드 0.
- `git diff --exit-code <기준 ref> HEAD -- .fullops-squad/rules/common`: 종료코드 0. 기준 ref와 적용 공통 규칙이 같다.
- 위 검사는 설계 작업본을 대상으로 수행했다. 커밋 후 FullOps lint 결과는 아래에 추가한다.
제품 코드는 수정하지 않았다. 제품 빌드·lint·기능 테스트는 dev와 tester의 후속 과제이므로 실행하지 않는다.
FullOps lint는 문서 검사 확인용으로 커밋 후 실행한다. 기준 ref의 제품 commands는 비어 있으므로 LINT-000은 후속 dev 과제에서 처리한다.

## 후속

coor는 설계 커밋을 dev 워크트리에 전달하고 새 dispatch 복귀 정보를 기록한다. dev 완료 SHA를 tester에게 전달한다.
dev는 실제 제품 검사 명령을 직접 실행해야 한다. merge-base의 빈 lint commands를 사용하는 FullOps 통과만으로 제품 lint 완료를 주장하지 않는다.
설계 커밋의 최종 SHA와 lint 결과는 worker_done에 함께 전달한다.
