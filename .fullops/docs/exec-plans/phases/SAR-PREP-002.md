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

## 준비 탐색과 1차 문서 검증

42adf86에서 DEV/TESTER 각각 code/documents find와 context를 실행했다. find는 초기 골격을 찾았으며 업무 기능 존재를 뜻하지 않는다. context는 각 20개 후보를 모두 keep했다. fallback·자동 conflict_ids·caution_ids는 없다. tester의 verify_runtime.py는 sensitive or oversized passage로 원문 미전송이며 로컬 keep이다. 결과 JSON 6개를 과제 키별로 보존한다.

수동 대조에서 기존 setup의 업무 404·DB 테이블 0·adapter unimplemented·UI 없음은 새 기능의 수락 조건과 범위가 다름을 확인했다. 원문은 보존하고 후속 지시서에 먼저 확인할 이력으로 표시했다. 과거 commands 빈 값과 QA 러너 최초 한계도 최신 보완 근거와 분리했다.

1차 검사 대상은 준비 문서 8개다. 로컬 링크 41개·요구사항 16개·metadata·queued 상태·필수 경로를 검사해 오류 0, 종료코드 0을 확인했다. 기준 ref 대비 원천·제품 경로·setup D02/D03·lint/board/fullops.json의 보존 검사는 각각 종료코드 0이다. deliverables.py --strict는 검사 13, 미작성 10, 문제 0, 경고 0, 종료코드 0이다. git diff --check와 스테이징 검사는 종료코드 0이다.

## 문서 준비 검사 결과

고정 준비 HEAD `06de8461a3c58658247184a79f88b7419e34d3be`에서 기준 `0dd08ec994771836c15d9d22a6a83393a71d7987`의 lint를 실행했다. 첫 실행은 product-lint 실패, ERROR 1, WARNING 0, 실행 불가 0, 종료코드 1이었다. 원인은 로컬 npm 의존성 부재이며 정확한 오류는 `sh: 1: prettier: not found`다. [첫 실패 기록](../logs/SAR-PREP-002/lint-prearchive.json)을 보존한다.

기존 lock을 사용해 `npm ci --prefix adapters`를 실행했고 종료코드 0을 확인했다. 제품 파일·lock은 바뀌지 않았다. 설치 상태가 달라진 뒤 같은 HEAD와 기준으로 재실행했다. product-lint make lint passed, ERROR 0, WARNING 0, 실행 불가 0, lint 종료코드 0이다. [재검증 기록](../logs/SAR-PREP-002/lint-prearchive-retry.json)을 보존한다. 제품 검사를 낮추거나 원천을 수정하지 않았다.

같은 준비 HEAD에서 `git diff --check 0dd08ec994771836c15d9d22a6a83393a71d7987 HEAD`와 `deliverables.py --repo . --strict`를 실행했다. 둘 다 종료코드 0이다. strict는 검사 13, 미작성 10, 문제 0, 경고 0이다. D01/D02/D03의 문서 작성 상태와 D04–D13 미작성 상태를 구분한다.

[문서·요구사항 검사](../logs/SAR-PREP-002/document-check.txt), [원천 보존](../logs/SAR-PREP-002/preserved-source.txt), [제품 보존](../logs/SAR-PREP-002/preserved-product.txt), [setup·하네스 보존](../logs/SAR-PREP-002/preserved-setup-and-harness.txt)을 기록했다. 초기 문서 8개·링크 41개·요구사항 행 16개의 오류는 0이며 각 보존 명령 종료코드는 0이다. 인계 탐색 링크 추가 뒤 최종 변경 문서 검사와 최종 HEAD lint는 아카이브·커밋 이후 수행하고 실제 SHA와 결과를 worker_done에 보낸다. 이 체크포인트를 최종 커밋 검사로 표시하지 않는다.

제품 코드·동작이 바뀌지 않아 build/test/runtime QA와 UI 캡처·영상은 미적용이다. setup 독립 QA를 반복하지 않았다. 이번 문서 검사와 등록 product-lint는 실제 실행했다. 문서 준비 수락 뒤 사용자 시작 지시, 기능 구현, 독립 QA·UI 직접 검수·독립 코드 리뷰와 운영 승인 조건은 후속에 남는다.

## 아카이브와 최종 인계 준비

work.py finish는 SAR-PREP-002 지시서와 완료 보고를 기존 날짜별 [designer 완료 아카이브](../../../handovers/logs/2026-10-03_to_designer.md)에 한 번 추가했다. designer 인박스는 비웠다. 기존 완료 기록은 보존했다. DEV/TESTER는 queued를 유지한다.

아카이브 이후 문서 8개와 로컬 링크 53개를 검사해 오류 0, 종료코드 0을 확인했다. 완료 로그의 단일 추가·빈 designer 인박스·queued DEV/TESTER도 확인했다. [검사 결과](../logs/SAR-PREP-002/document-final-before-commit.txt)를 보존한다. strict·공백·원천·제품 보존 명령의 각 종료코드 0을 별도 실행 로그에 기록했다. 이 검사는 최종 커밋 전 문서 상태의 결과다.

최종 커밋 이후 같은 기준 ref의 FullOps lint와 strict·공백·보존 검사를 실행한다. 최종 고정 SHA·실제 결과는 preamble worker_done으로 전달한다. 최종 lint JSON은 레포 밖 /tmp/SAR-PREP-002-validation/lint-final.json에 남기며, 이 문서의 체크포인트 JSON과 구분한다.

최종 staging 검사에서 자체 문서 검사 로그 2개의 EOF 빈 줄을 검출해 종료코드 2였다. 두 로그의 마지막 빈 줄을 정리했다. 원천 hard break는 수정하지 않았다. 해당 커밋이 진행되지 않은 상태의 lint 시도는 dirty tree로 종료코드 2였으며 통과 증거로 사용하지 않는다. [실패 기록](../logs/SAR-PREP-002/final-commit-check-failure.txt)을 보존하고 깨끗한 최종 커밋에서 다시 검사한다.
