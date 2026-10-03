---
title: SAR-PREP-002 — MVP 개발 준비 실행 기록
status: draft
updated: 2026-10-03
owner: designer
tasks: [SAR-PREP-002]
summary: 최신 원천의 MVP 개발 준비 결정과 문서 검증 및 후속 책임을 기록한다
---

# SAR-PREP-002 — MVP 개발 준비 실행 기록

## 범위와 기준

작업자는 designer다. 작업 브랜치는 fullops/designer다. 기준 ref는 `0dd08ec994771836c15d9d22a6a83393a71d7987`이고 착수 HEAD는 `392a21d`다. 원천은 service-design `7bc9ea190ea549fae8b047e850247a19322fc9c3`이다. [원천 고정 기록](../../planning/sources/silent-agent-relay/source.json), [요청](../../planning/SAR-PREP-002-request.md), 받은 지시서를 읽었다.

fullops-common-0.3.2, FULLOPS.md, project.md, 문서 규칙을 적용했다. fullops-work/fullops-deliverables로 기획·인계·기록을 관리한다. orchestration과 fullops-orca의 worker 회신 규약을 적용한다. 제품 코드는 수정하지 않는다. 기술 계획은 후속 DEV의 같은 기능 과제다.

## 기획 결정과 산출물

- [D01](../../planning/business-plan.md)에 제품 가치·대상 사용자·freemium 방향과 미정 결정을 정리했다.
- [새 MVP D02](../../planning/product-specs/SAR-MVP.md)에 MVP-01–16, C1–C5, frozen wire, 고정 수치, 사용자 흐름과 UI 수락 조건을 연결했다. 기존 setup D02는 수정하지 않는다.
- [기능 백로그](../../planning/SAR-MVP-backlog.md)에 SAR-MVP-001–007을 정리했다. 첫 기능 001은 등록·수락·안전 전달·human-gate의 로컬 합성 흐름이다. ingest/queue-only로 줄이지 않는다.
- [DEV](../../../handovers/to_dev.md)와 [TESTER](../../../handovers/to_tester.md)를 queued로 준비했다. DEV는 같은 과제에서 기술 계획·구현·관련 회귀·기술 문서 갱신을 한다. TESTER는 DEV 완료 뒤 고정 통합 후보를 기다린다. 이번에는 배정하지 않는다.
- 현재 없는 실제 업무 SQL·UI·실제 제품 어댑터·운영 Tunnel을 완료로 표시하지 않는다. 초기 골격 완료와 MVP 기능 완료를 구분한다.

공개 A2A v0.3.0 개념 검토 완료·wire 비호환·wire 필드 미수입·owner 승인 대체 금지를 반영했다. 최신 A2A delta를 재검토하지 않았다. wire와 원천 원문은 보존한다.

transport delivered는 ACK 성공 뒤이며 processing은 공유 exec claim 뒤라는 protocol 정본을 유지한다. README/product의 축약 문장을 실행권 부여로 해석하지 않는다. 같은 digest의 다른 M.id로 승인 재사용을 허용하지 않는다.

미정 disclosure와 result schema는 정책 없음 deny·optional 데이터 반환 보류로 유지한다. 첫 로컬 기능은 실데이터 silent 성공을 요구하지 않는다. positive silent done은 백로그 003에 held로 남긴다. Free N·가격·slot-unit·추가 자원 제한·운영 설정은 값을 만들지 않고 DEC-01–05에 담당과 재개 조건을 기록했다.

D04는 별도 에셋·장식 목업이 필요하지 않아 미작성으로 유지한다. UI 제품 조건은 D02에 있고 구현 후보의 V-01–04 직접 검수는 designer 책임이다. 후속 D03/D05–D10은 DEV 책임이며 지금 완성한 것으로 표시하지 않는다.

## 검증 계획과 실제 결과 기록 위치

designer는 원천 추적·요구사항/백로그/QA 연결·로컬 링크·front matter·파일 소유권을 직접 검토한다. D01/D02와 일반 문서 metadata는 deliverables.py --stamp로 작성한다. strict 산출물 검사와 git diff --check, 기준 ref의 깨끗한 커밋 FullOps lint를 실행한다. 실제 종료코드는 실행 후 아래에 추가한다.

실행 로그는 레포 밖 `/tmp/SAR-PREP-002-validation/`에 보존한다. 필요한 요약은 이 기록과 완료 아카이브에 쓴다. 레포 안에 검증용 제품 코드나 QA 러너를 만들지 않는다. 제품 build/test/runtime QA·영상·캡처는 제품 코드·UI가 없는 이번 준비에서 미적용이다. 등록 product-lint는 지시서가 요구한 FullOps lint 실행 결과로 확인한다.

기존 setup 증거는 원래 대상 `0cc10b083771be9b3423833b222c57d426315333`, 최종 통합 후보 `9c96232e6230e00319c33ff77637c80923e2438b`의 초기 골격 근거다. [최종 통합 리뷰](../../evaluations/qa-reports/SAR-SETUP-001-INTEGRATION-SOURCE-review/report.md)와 기존 PLANS의 main 수락 기록을 확인했다. 새 MVP 동작 QA로 바꾸어 보고하지 않는다.

## 후속 담당과 수락 경계

coor는 사용자에게 첫 실행 단위 SAR-MVP-001-DEV를 제시한다. 구현 시작 지시 뒤 준비 커밋·원천·규칙을 worker에 반영하고 실제 Dispatch와 시작 HEAD를 기록한다. DEV 완료 뒤 tester 독립 QA·designer 직접 UI 검수·별도 세션 독립 fixed-SHA 코드 리뷰를 조정한다.

designer는 DEC-01–03의 제품 정책을 담당한다. dev는 DEC-04의 인터페이스와 기술 구현을 맡는다. dev/ops와 coor는 DEC-05의 설정과 실행 승인을 맡는다. 미해결 critical/high·held·운영 공개 제한은 유지한다. 준비 문서 완료는 MVP 제품 수락이나 외부 배포 승인이 아니다.
