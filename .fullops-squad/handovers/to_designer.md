---
title: SAR-SETUP-001 — 프로젝트 초기 구성과 lint 범위 및 역할별 수락 기준 확정
status: draft
updated: 2026-10-03
owner: designer
tasks: [SAR-SETUP-001]
summary: 프로젝트 초기 구성과 lint 범위 및 역할별 수락 기준 확정
---

# SAR-SETUP-001 — 프로젝트 초기 구성과 lint 범위 및 역할별 수락 기준 확정

- 작성일: 2026-10-03
- From / To: coor / designer
- 상태: running
- 승인된 범위: 기획·역할별 지시서·탐색 및 완료 기록 작성과 커밋. 코드 수정 금지.
- 담당 워크트리 / 브랜치: `/home/shin/orca/workspaces/KnowsLink/fullops-designer` / `fullops/designer`.
- 병합 책임자 / 기본 브랜치: coor / main.
- 복귀: `/home/shin/orca/workspaces/KnowsLink/fullops-coor`, `term_89f25ea4-e70e-46c0-8514-e95f8cf81928`, Run `run_8ca8bc058ab7`.
- task id / dispatch id: `task_49c6e00e6470` / `ctx_1f186b4db5b5`. 전송 권한은 세션 preamble에서만 사용한다.

## 현재 상황과 확인 근거

요청은 프로젝트 초기 구성 및 필요한 lint다. 원천은 `.fullops-squad/docs/planning/sources/silent-agent-relay/`다.
준비 커밋 `0cc35f0`에서 원천 일곱 문서와 source.json을 읽었다. 기술 선택은 이미 잠긴 부분과 dev가 정할 부분으로 나뉜다.

## 적용 기준과 예외

`fullops-common-0.3.1`, project.md, 문서 작성 규칙을 적용한다. 기준 ref는 `00b4cb34ae6e9f9fbc0b733ecaa3a2095fbc88eb`이다.
외부 원천 SHA는 `404ff834c0607055d63d2053bf7771d2f46ad3ae`다. 코드·원천 수정과 기준 완화 예외는 없다.

## 먼저 읽을 문서

FULLOPS.md, 공통 규칙 네 문서, project.md, document-writing.md, contexts/designer.md, SAR-SETUP-001-request.md, 고정 원천 일곱 문서와 source.json을 읽었다.
dev/tester의 먼저 읽을 목록은 각 역할의 Jev code/documents find와 context 기록으로 분리했다. API 실패를 기록하고 전부 keep했다.

## 해야 할 일과 파일 소유권

- [x] 원천의 확정 결정과 초기 구성 범위를 구분한다.
- [x] D02 수락 기준과 산출물 선택 근거를 작성한다.
- [x] dev 지시서와 dev SHA 이후 tester 지시서를 작성한다.
- [x] 역할별 Jev 탐색·문서 탐색·context 호출과 대체 후보 근거를 남긴다.
- [ ] 문서 검증·완료 기록·커밋 후 worker_done을 전송한다.

## 완료 기준과 검증

D02와 두 지시서가 원천 범위·고정 ref·역할 책임·선행 조건·수락 기준을 연결한다. 문서 검증과 Git 공백 검사를 수행한다.
제품 코드 수정과 테스트는 이 역할에 적용하지 않는다. 문서 lint는 커밋 후 기준 ref로 실행해 한계를 기록한다.

## 갱신할 산출물

D02: docs/planning/product-specs/SAR-SETUP-001.md. D03은 dev에게 배정하며 이 세션에서 작성하지 않는다.

## 기대 산출물

D02, handovers/to_dev.md, handovers/to_tester.md, 역할별 Jev 기록, docs/exec-plans/phases/SAR-SETUP-001.md, 완료 아카이브.

## 제약·협업·후속

전체 MVP 구현으로 확대하지 않는다. 가격·정책 수치를 발명하지 않는다. coor가 설계 커밋을 dev에게 전달하고 dev 완료 SHA 이후 tester를 dispatch한다.

## 완료 보고

문서 검증 후 결과를 기록한다.
