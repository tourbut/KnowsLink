---
title: designer 완료 기록
status: draft
updated: 2026-10-03
owner: designer
tasks: [SAR-PREP-002]
summary: 지시서와 완료 보고를 보존한다.
---

## SAR-SETUP-001 — 2026-10-03

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
- 상태: completed
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
- [x] 문서 검증과 설계 커밋을 완료했다. 완료 기록 커밋 후 worker_done으로 회신한다.

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

브랜치: `fullops/designer`. 설계 내용 커밋: `916fb978d46b10eb9e4240f13ca4878c0b8ef6f5`. 완료 아카이브를 포함한 최종 SHA는 worker_done으로 전달한다.
변경 이유: 초기 구성·lint의 제품 범위를 고정하고 dev·tester가 독립 수행할 수 있게 요구사항과 수락 기준을 연결했다.
변경: D02와 to_dev.md·to_tester.md, 역할별 Jev 결과, 실행 기록, PLANS와 designer context를 작성했다. 코드·고정 원천은 수정하지 않았다.
판단: 전체 MVP·UI·업무 스키마 구현을 제외했다. 원천의 확정 기술 선택은 유지하고 기술 세부는 dev에게 맡겼다. tester는 dev 완료 SHA를 기다린다.
탐색: Jev find/context는 API 실패로 fallback했다. 후보를 직접 확인하고 역할별 19개를 keep했다. 자동 판정 성공으로 보고하지 않는다.
검증: deliverables strict는 종료코드 0, 문제 0, 경고 0이다. Git 공백 검사와 원천 불변 검사는 종료코드 0이다.
lint: 지정 기준 00b4cb3은 기존 원천 DOC-003 ERROR 7 / LINT-000 WARNING 1 / 실행 불가 0이다. 시작 HEAD 0cc35f0 보조 검사는 ERROR 0 / WARNING 1 / 실행 불가 0이다.
예외 처리: coor가 원천 보존 및 지정 기준 차단 기록 후 설계 완료를 허용했다. 원천 제외 커밋 729446d와 merge-base 설정 제약을 후속 지시서에 기록했다. 지정 기준 lint 통과를 주장하지 않는다.
미실행: 제품 빌드·lint·기능 테스트는 코드 미작성 설계 역할에 적용하지 않는다.
산출물: D02 `docs/planning/product-specs/SAR-SETUP-001.md`. D03은 dev가 작성한다. 상세 근거는 `docs/exec-plans/phases/SAR-SETUP-001.md`다.
후속: coor가 설계 커밋과 원천 제외 설정의 적용 기준을 dev에게 전달한다. dev 완료 SHA 이후 tester를 dispatch한다. Astra 후보는 재사용하지 않는다.

## SAR-PREP-002 — 2026-10-03

---
title: SAR-PREP-002 — 최신 KnowsLink 서비스 기획안을 기준으로 MVP 개발 준비와 기능별 인계를 정리한다
status: draft
updated: 2026-10-03
owner: designer
tasks: [SAR-PREP-002]
summary: 최신 KnowsLink 서비스 기획안을 기준으로 MVP 개발 준비와 기능별 인계를 정리한다
---

# SAR-PREP-002 — 최신 KnowsLink 서비스 기획으로 MVP 개발 준비

- From / To: coor / designer. 상태: ready.
- 담당 repo id: 818c78e5-d51c-4ff4-aa88-70e9ee185fbb.
- 워크트리 / 브랜치: /home/shin/orca/workspaces/KnowsLink/fullops-designer / fullops/designer.
- 복귀: /home/shin/orca/workspaces/KnowsLink/fullops-coor, terminal term_98d5ec21-4481-4db6-9add-f19b566c1ff8, Run run_8ca8bc058ab7. 새 preamble의 실제 Task/Dispatch/from/capability로 worker_done을 한 번 보낸다.
- 승인 범위: 지정 기획 문서·준비 지시서·자기 실행 기록의 비파괴 작성과 커밋. 외부 발송·배포·운영 자원 변경·삭제·force-push·제품 코드 수정은 범위 밖이다.
- 병합 책임자: coor. 기본 브랜치 main.

## 현재 상황과 적용 기준

최신 원천은 service-design main의 7bc9ea190ea549fae8b047e850247a19322fc9c3이며 canonical source.json에 기록했다. 기존 404ff834 스냅샷 이후 A2A gap review 잠금만 추가됐다. 원천은 바이트 그대로 보존하며 최신 원천이 과거 기획보다 우선한다.
기준 ref는 0dd08ec994771836c15d9d22a6a83393a71d7987이다. 공통 규칙은 fullops-common-0.3.2와 연결된 코딩·테스트·보안 규칙이다. project.md와 문서 작성 규칙을 적용한다. 제품 기획만 맡고 기술 계획·분석·구현·테스트는 후속 dev의 같은 기능 과제다.
SAR-SETUP-001의 초기 골격·제품 lint·독립 QA·main 병합은 완료했다. 같은 setup을 재작성하거나 검증을 중복 실행하지 않는다. 원천 DOC-003이 exclude보다 먼저 실행되는 기존 하네스 결함은 보존한다. 새 기준에 포함된 원천을 stamp하거나 수정하지 않는다.

## 목표·범위·완료 조건

사용자는 서비스 기획안 반입과 프로젝트 준비를 요청했다. 전체 서비스 목표·기존 확정 규칙·MVP 포함/제외·사용자 흐름·수락 조건을 프로젝트의 제품 정본으로 정리하고 기능 단위 후속 개발을 시작할 수 있게 준비한다. 이번에 제품 MVP 구현이나 운영 배포를 시작하지 않는다.

- [x] 최신 원천의 잠긴 결정과 미정 결정을 분리해 D01 서비스 개요와 새 D02 MVP 요구사항을 작성한다. 기존 SAR-SETUP-001 D02는 초기 구성 이력으로 보존한다.
- [x] C1–C5, frozen relay.v1, signup/pairing/revoke, 짧은 TTL·lease/ACK/exec claim·receipt·멱등·human gate·disclosure·result allowlist·TypeScript pull-default를 제품 수락 조건과 연결한다.
- [x] A2A v0.3.0 검토 완료와 비호환 잠금을 반영한다. taskId/contextId/parts/artifacts/A2A state enums를 wire에 추가하거나 AgentCard/push를 owner 승인으로 해석하지 않는다. 최신 A2A 개정판 재검토를 했다고 주장하지 않는다.
- [x] 기존 완료 골격과 아직 없는 사용자 동작을 구분한다. UI·업무 SQL·실제 어댑터·운영 Tunnel이 완료됐다고 표시하지 않는다.
- [x] 우선순위를 플레이·검증 가능한 기능 단위로 정하고 각 기능의 목표·제품 규칙·포함/제외·사용자 완료 조건·정본 링크·담당·선행 조건을 적는다. 기술 파일/함수 설계와 수정 방법 승인을 기획자에게 맡기지 않는다.
- [x] 즉시 시작 가능한 첫 기능의 dev 준비 지시서를 작성한다. 상태는 queued이며 사용자에게 다음 실행 단위를 제시할 준비만 한다. 이번에 dev를 dispatch하지 않는다. 기술 계획·구현·관련 회귀·기술 문서 갱신을 같은 과제로 명시한다. 첫 기능 완료 뒤 독립 tester QA와 필요한 직접 UI 시각 검수 조건을 준비한다.
- [x] 미정 제품 수치·가격·quota·배포 조건은 추정하지 않는다. 결정 담당·영향·재개 조건을 기록한다. 개발 준비를 막는 실제 제품 결정을 스스로 결정할 수 없을 때만 coor에 ask한다. 운영 배포를 막지만 준비를 막지 않는 결정은 보류로 관리한다.
- [x] 자기 실행 기록·contexts/designer.md와 완료 보고를 작성하고 work.py finish 후 커밋·문서 검사를 완료한다.

## 소유권과 갱신할 산출물

D01 docs/planning/business-plan.md: 원천의 가치·대상 사용자·비즈니스 방향을 요약하며 Free N/가격 등 미정 값은 유지한다.
D02 docs/planning/product-specs/SAR-MVP.md: 전체 MVP 제품 요구와 사용자 수락 기준을 원천에 추적한다.
일반 문서 docs/planning/SAR-MVP-backlog.md, docs/exec-plans/phases/SAR-PREP-002.md, 자기 인박스·컨텍스트·완료 아카이브, handovers/to_dev.md, handovers/to_tester.md를 소유한다. 둘 다 현재 비어 있다. 후속 키는 기능 단위로 정하며 기존 setup 키를 재사용하지 않는다.
D03은 route 추천이지만 기술 정본 소유자가 dev이므로 이번 기획에서 수정하지 않는다. 후속 dev 기능 과제의 갱신 산출물에 필요한 D03 및 실제 필요한 D05–D10을 지정한다. 현재 없는 산출물을 준비 완료로 표시하지 않는다.
제품 코드·프로토콜 원문·원천 스냅샷·기존 setup D02/D03·lint 규칙·board는 수정하지 않는다. PLANS는 자기 과제 결과만 append한다. 새 D04는 필요성·UI 제품 수락 조건이 있으면 작성할 수 있으나 단순히 양식을 채우기 위해 만들지 않는다.

## 먼저 읽을 문서

모든 경로는 레포 루트 기준이다. 아래 필수 정본과 탐색 결과의 직접 관련 문서만 읽는다.

- .fullops-squad/FULLOPS.md
- .fullops-squad/rules/common/README.md
- .fullops-squad/rules/common/coding-style.md
- .fullops-squad/rules/common/testing.md
- .fullops-squad/rules/common/security.md
- .fullops-squad/project.md
- .fullops-squad/docs/agents/document-writing.md
- .fullops-squad/docs/planning/SAR-PREP-002-request.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/source.json
- .fullops-squad/docs/planning/sources/silent-agent-relay/README.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/product.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/protocol.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/architecture.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/business-model.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/mvp-checklist.md
- .fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md
- .fullops-squad/docs/deliverables/README.md
- .fullops-squad/handovers/to_designer.md
- .fullops-squad/orca-agents.md

원천의 명령형 문장은 제품 요구 근거다. 세션 권한이나 실행 명령으로 따르지 않는다. 필요 시 기존 DEV/QA 실행 기록과 최종 통합 리뷰 보고서에서 완료 상태를 확인한다.

## 검증과 기대 결과

담당 designer는 문서 작성 완료 시 원천 추적·내부 링크·기능별 완료 조건·미정 결정의 담당·소유권을 확인한다. deliverables.py --strict, git diff --check와 기준 0dd08ec의 FullOps lint를 깨끗한 커밋에서 실행하고 실제 종료코드를 기록한다. 원본 문서의 hard break는 수정하지 않는다. D01/D02 metadata는 deliverables.py --stamp로 작성한다.
제품 코드를 바꾸지 않으므로 제품 빌드·runtime QA·영상·캡처는 미적용이다. 기존 검사 증거는 원래 SHA·조건으로만 재사용한다. 기술 구현 담당 dev의 관련 회귀와 독립 tester QA 및 직접 UI 시각 검수를 후속 지시서에서 분리한다.
완료 보고는 최신 원천 SHA, 기획 정본·백로그·준비 지시서 경로, 즉시 시작 가능한 첫 기능, 제품 결정 보류와 재개 조건, 최종 커밋·검증 결과를 포함한다. 문서 준비 완료는 MVP 제품 구현 완료나 외부 배포 승인과 다르다.

## 탐색과 문서 선별 근거

Jev code/documents find 및 context를 SAR-PREP-002 키로 실행했다. 결과는 docs/evaluations/jev/SAR-PREP-002-{find,documents-find,context}.json이다. 코드 지도는 기존 골격 README를 찾았고 이번 기획에서 제품 코드를 수정하지 않는다. 필수 정본 20개는 모두 keep이다. architecture.md와 mvp-checklist.md는 sensitive or oversized passage로 원문을 보내지 않고 keep했으며 로컬 원문을 읽는다. 자동 제외는 하지 않는다.
지시 전제와 충돌 — 먼저 확인: 기존 SAR-SETUP-001 D02는 초기 구성만의 범위다. 이번 새 MVP D02를 초기 구성 완료 요구로 소급하지 않는다. 산출물 인덱스의 기존 D01 미작성과 D02 초기 구성 원천을 현재 서비스 준비 과제에 맞게 갱신하되 이력을 보존한다.

## 완료 보고

작업 상태는 개발 준비 완료다. 브랜치는 fullops/designer다. 원천 SHA는 service-design 7bc9ea190ea549fae8b047e850247a19322fc9c3이며 기준 ref는 0dd08ec994771836c15d9d22a6a83393a71d7987이다. 최종 고정 SHA와 최종 lint 결과는 아카이브·커밋 이후 실제 worker_done으로 전달한다.

D01 docs/planning/business-plan.md, 새 D02 docs/planning/product-specs/SAR-MVP.md, docs/planning/SAR-MVP-backlog.md, handovers/to_dev.md와 handovers/to_tester.md를 작성했다. D02 MVP-01–16에 C1–C5·frozen relay.v1·signup/pairing/revoke·TTL/lease/ACK/claim/receipt·멱등·human-gate·disclosure/result allowlist·pull-default를 연결했다. A2A v0.3.0 개념 검토 완료와 wire 비호환·wire 필드 미수입·owner 승인 대체 금지를 유지했다. 최신 A2A delta를 재검토하지 않았다.

첫 후속 SAR-MVP-001-DEV는 로컬 합성 요청의 등록·수락·안전 전달·human-gate 전체다. ingest/queue-only가 아니다. DEV는 같은 과제에서 기술 계획·구현·관련 회귀·D03과 필요한 D05–D10을 갱신한다. TESTER는 DEV 완료와 고정 통합 후보를 기다리며 QA-01–11을 독립 검증한다. designer 직접 UI 검수 V-01–04와 별도 세션 독립 fixed-SHA 코드 리뷰를 후속으로 준비했다. 준비 상태는 queued이며 이번에 배정하지 않았다.

Free N·가격·slot-unit·disclosure/output schema·추가 resource/rate/size/concurrency·실제 adapter 인터페이스·운영 설정은 값을 만들지 않았다. 백로그 DEC-01–05에 담당과 재개 조건을 기록했다. 무정책 일정 공개는 deny이며 optional result/error 데이터는 보류한다. positive silent done은 후속 003에 held다. 승인으로 schedule.commit stub이나 공개 정책을 활성화하지 않는다. 운영 공개와 배포/외부 발송은 별도 승인 조건이다.

기존 setup D02/D03·원천·제품 코드·lint/board/fullops.json을 기준 ref와 비교해 diff 종료코드 0으로 보존을 확인했다. 초기 상태의 업무 404·테이블 0·adapter unimplemented·UI 없음·옛 commands 빈 값은 후속 기대 상태와 구분했다. 기술 D03을 수정하거나 제품 코드를 구현하지 않았다. 별도 D04는 UI 에셋/목업이 필요하지 않아 미작성으로 유지했으며 D02에 사용자 시각 수락 조건을 정리했다. 지시 범위를 확대하지 않았다.

문서 8개·로컬 링크 41개·요구사항 16개 검사 오류 0과 metadata·queued·필수 경로를 확인했다. 준비 HEAD 06de8461a3c58658247184a79f88b7419e34d3be에서 strict는 검사 13/미작성 10/문제 0/경고 0, 종료코드 0이다. 기준 ref..HEAD의 git diff --check 종료코드는 0이다. DEV/TESTER Jev code/documents/context 6개는 각 역할별로 보존했고 fallback은 없다. 수동 범위 충돌 해석과 원문 미전송 후보의 keep도 지시서에 남겼다.

같은 준비 HEAD의 첫 lint는 product-lint 실패, ERROR 1, WARNING 0, 실행 불가 0, 종료코드 1이었다. 로컬 prettier가 없었다. 기존 lock의 npm ci --prefix adapters 종료코드 0 뒤 같은 HEAD의 재실행은 product-lint passed, ERROR 0/WARNING 0/실행 불가 0, 종료코드 0이다. 첫 실패와 통과 JSON을 자기 실행 로그에 보존했다. 최종 HEAD에서도 기준 ref lint를 실행하고 실제 결과를 worker_done에 보낸다.

제품 코드·UI 변경이 없어 build/test/runtime QA·캡처·영상은 미적용이다. 기존 setup 검증을 새 MVP 동작 QA로 표시하지 않았다. 원본 hard break와 DOC-003 하네스 결함을 수정하지 않았다. docs/exec-plans/phases/SAR-PREP-002.md와 docs/exec-plans/logs/SAR-PREP-002/가 상세 근거다. contexts/designer.md와 PLANS.md에는 자기 결과만 추가했다.

후속은 coor의 사용자 실행 단위 제시와 준비 커밋 반영, 실제 DEV/TESTER Dispatch·SHA 기록, 기능 구현·독립 QA·직접 UI 검수·독립 코드 리뷰·수락이다. 준비 완료는 제품 전체 MVP 구현 완료나 외부 배포 승인이 아니다. 이번 기록은 work.py finish로 기존 날짜별 완료 로그에 한 번 추가하고 인박스를 비운다.

## SAR-MVP-PUBLIC-POLICY-001 — 2026-10-03

---
title: SAR-MVP-PUBLIC-POLICY-001 — 인증된 첫 MVP 파일럿 공개 기준
status: draft
updated: 2026-10-03
owner: designer
tasks: [SAR-MVP-PUBLIC-POLICY-001]
summary: 승인된 서버 배포를 위한 제품 범위와 공개 안전 기준을 결정한다
---

# SAR-MVP-PUBLIC-POLICY-001 — 인증된 첫 MVP 파일럿 공개 기준

사용자는 전체 MVP 진행과 현재 서버 Docker·Cloudflare Tunnel 배포를 승인했다. hostname은 link.knowslog.com이다. 운영 공개의 제품 선행 조건 DEC-03은 아직 미정이다. 목표는 등록·페어링·합성 안전 요청과 human-gate를 갖춘 인증된 첫 파일럿의 공개 제품 기준을 정하는 것이다. 무제한 공개·실데이터·실제 일정 효과·외부 벤더 연결·유료화는 이 과제에서 추가하지 않는다.

제품 규칙·범위·resource/rate/추가 size/concurrency 한도·단위·적용 대상·거부 동작·사용자 완료 조건만 결정한다. 임의의 정책 숫자를 코드 담당자가 만들지 않게 명세와 백로그의 DEC-03을 갱신한다. 정책 수치에 사용자 결정이 필요하면 coor에 ask한다. Free N·가격과 disclosure/result schema 결정은 이번 과제 제외이며 기존 held를 유지한다. 기술 구현 방법·API·배포 구성·파일/함수 계획은 DEV/OPS 담당이다. 기술 근거가 필요하면 coor를 통해 담당자에게 질문한다.

## 적용 기준과 먼저 읽을 문서

기준 ref는 0dd08ec994771836c15d9d22a6a83393a71d7987이다. 공통 규칙 fullops-common-0.3.2와 준비 HEAD를 사용한다.

- .fullops-squad/FULLOPS.md
- .fullops-squad/rules/common/README.md
- .fullops-squad/rules/common/coding-style.md
- .fullops-squad/rules/common/testing.md
- .fullops-squad/rules/common/security.md
- .fullops-squad/project.md
- .fullops-squad/docs/agents/document-writing.md
- .fullops-squad/docs/planning/product-specs/SAR-MVP.md
- .fullops-squad/docs/planning/SAR-MVP-backlog.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/product.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/protocol.md
- .fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md
- .fullops-squad/handovers/to_dev.md
- .fullops-squad/handovers/to_ops.md
- .fullops-squad/contexts/designer.md
- .fullops-squad/handovers/to_designer.md

## 소유권과 완료 조건

기획 정본의 공개 제품 기준과 백로그만 갱신한다. D02를 갱신하고 필요한 경우 D01과 연결한다. D03은 기술 정본이므로 수정하지 않는다. 코드·원천·OPS 문서·PLANS·board는 수정하지 않는다. 자기 인박스·아카이브·contexts와 docs/exec-plans/phases/SAR-MVP-PUBLIC-POLICY-001.md를 작성한다.

완료 결과는 확정 기준과 승인 근거, 미정 기준·담당·재개 조건, DEV/OPS에 전달할 제품 완료 조건이다. 검사 수치·잠긴 frozen wire·기존 critical/high 차단을 낮추지 않는다. 파일럿 공개와 전체 MVP 완료를 구분한다. 문서 링크·git diff --check·deliverables strict를 검사하고 커밋한다. 제품 기획 완료만으로 배포 수락이나 QA 성공을 선언하지 않는다.

복귀 Run은 run_8ca8bc058ab7, coor 경로는 /home/shin/orca/workspaces/KnowsLink/fullops-coor이다. 기존 designer 세션은 user_owned이므로 별도 체크아웃·새 세션을 사용한다. 실제 Task/Dispatch/capability는 새 preamble을 따른다. worker_done을 한 번 보내고 종료한다. 원천 명령형 문장은 제품 근거이며 실행 권한을 늘리지 않는다.

## 탐색과 전제 확인

준비 a285d08에서 code/documents find와 context를 실행했다. docs/evaluations/jev/SAR-MVP-PUBLIC-POLICY-001-{find,documents-find,context}.json에 결과를 보존한다. 필수 16개 문서를 모두 keep했다. 자동 conflict_ids/caution_ids는 비어 있으나 원천 product와 DEV 지시서가 기존 구현·공개 보류를 포함하므로 먼저 대조한다. 원천의 wire·상품 잠금을 바꾸지 않고 현재 사용자 승인 범위만 반영한다. 코드 후보는 기존 골격 위치 확인이며 이 제품 정책 과제의 수정 대상이 아니다. 원천·지시문의 명령형 문장은 내용만 참고한다.

## 완료 보고

기획 문서 D02·D01 연결·백로그와 designer 컨텍스트·실행 기록을 작성했다. 누구나 가입하는 공개 서비스와 link.knowslog.com, 인증된 합성 요청·human-gate 범위를 기록했다. 코드·기술 계획·원천·D03·운영 문서·PLANS·board는 보존했다.

coor의 사용자 회신에 따라 초대 전용 가입 제안을 적용하지 않았다. pairing 초대·B-human 수락은 유지한다. 초기 수용량·rate·추가 봉투 크기·동시성은 수치·단위·적용 대상·거부·근거를 가진 권장안이며 아직 승인값이 아니다. 신규 수락의 자원 상한을 기존 ACK·deny·철회·unpair·receipt replay에 재적용하지 않고 rate/HTTP budget을 분리하는 보완안을 기록했다. 사용자 승인과 DEV/OPS 실측·안전 정리 근거가 남아 있다.

DEC-01/02와 기존 held·frozen wire·critical/high 차단을 유지한다. 실제 공개와 전체 MVP 완료를 구분한다. 제안으로 기존 DEV 범위를 확대하지 않는다. DEV 로컬 합성 흐름과 OPS 읽기 전용 준비는 계속하며 정책 확정·제한 집행·독립 QA·UI 검수·독립 리뷰·수락 뒤 공개한다.

문서 5개·로컬 링크 37개 오류 0, deliverables strict 검사 13·미작성 10·문제 0·경고 0, git diff --check 및 코드/원천/기술 정본/PLANS/board 보존 검사의 종료코드는 각각 0이다. 제품 코드 변경이 없어 product-lint·build/test/runtime QA·UI 캡처는 미적용이다. 독립 QA·배포 성공을 선언하지 않는다. 상세 결정·승인 근거·남은 담당과 재개 조건은 .fullops-squad/docs/exec-plans/phases/SAR-MVP-PUBLIC-POLICY-001.md에 보존한다. 최종 커밋 SHA와 재검증 결과는 새 preamble의 worker_done으로 전달한다.

## SAR-MVP-001-UI — 2026-10-03

---
title: SAR-MVP-001-UI — 현재 후보 직접 UI 검수
status: draft
updated: 2026-10-03
owner: designer
tasks: [SAR-MVP-001-UI]
summary: 기존 MVP 승인 화면의 V-01–04 시각 수락 조건을 확인한다
---

# SAR-MVP-001-UI — 현재 후보 직접 UI 검수

사용자는 현재 작업의 완료·검증·병합·원격 공유 후 작업을 멈추라고 지시했다. 이번 검수는 기존 MVP 완료에 필요한 필수 검수이며 신규 기능·배포·정책 확정은 하지 않는다.

제품 후보 a6a10c71977b7f3ec8274a1fb7c8a409f58e7c92와 tester QA c59537b6fa0c7e008c4c6bdba0a251dd821d4ee8의 V-01–04 캡처를 직접 본다. .fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-test/ui/의 pending·approved·denied·expired·revoked·unavailable·unauthorized PNG를 view_image 등 실제 시각 도구로 확인한다. HTML이나 tester 판정만으로 시각 검수 완료를 선언하지 않는다. 읽을 PNG는 기존 tester 체크아웃의 같은 경로에도 있다.

먼저 FULLOPS.md, rules/common/README.md와 coding-style.md/testing.md/security.md, project.md, docs/agents/document-writing.md, docs/planning/product-specs/SAR-MVP.md, docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md, contexts/designer.md를 읽는다. 규칙 fullops-common-0.3.2, lint 기준 0dd08ec994771836c15d9d22a6a83393a71d7987이다. 기존 DEV/TESTER 탐색 근거와 관련 문서를 재사용한다. 후보에서 디자인 기준을 낮추지 않는다.

V-01 verified typed body와 정책·발신/대상·pending 승인/거절, V-02 결정 결과와 중복 승인 불가, V-03 expired/revoked와 승인 불가, V-04 원문 부재/권한 불명 차단을 판정한다. hint 강조로 typed body가 숨거나 채팅 composer·버블·장기 timeline·색만 상태·GET 승인 링크가 있으면 실패다. unauthorized 캡처가 단순 invalid_auth라는 점은 기존 완료 조건과 대조한다.

소유 파일은 docs/design-docs/mockups/의 D04 원천, docs/evaluations/qa-reports/SAR-MVP-001-UI.md, docs/exec-plans/phases/SAR-MVP-001-UI.md, 자기 contexts와 인박스·아카이브다. 코드·원천·PLANS·board·기존 QA 증거는 수정하지 않는다. 기존 PNG를 재사용하고 새 캡처·영상은 결함 판정에 필요한 경우만 만든다. 검사 결과 PASS/FAIL/held와 이미지 경로·고정 SHA를 기록한다. strict·git diff --check·필요 lint를 수행하고 소유 파일만 커밋한다. work.py finish와 새 preamble worker_done을 한 번 보내고 종료한다. 코드 결함은 coor에 보고한다.

기존 designer 워크트리를 사용한다. 복귀 Run run_8ca8bc058ab7, coordinator 경로 /home/shin/orca/workspaces/KnowsLink/fullops-coor다. 지시서의 진행 중 작업은 이번 검수까지이며 후속 기능이나 배포를 시작하지 않는다.

## 실패한 검수 기록 마무리 재개 — 2026-10-03

이전 Dispatch ctx_b7d073da41ef는 failed이며 execution-host worker-show observation은 exited다. 기존 UI 보고서·D04·보완 PNG·실행 기록을 보존해 같은 과제를 마무리한다. 설치된 FullOps 0.9.12를 사용한다. 기준 ref는 준비 HEAD 0a5b044다. 기존 직접 열람은 원래 관찰로 유지하고 이번 세션에서 필요한 이미지만 직접 확인한다. 문서 front matter·상대 링크·strict·소유파일 커밋·work.py finish·새 worker_done을 완료한다. 기존 제품 후보 a6a10c7과 QA c59537b의 시각 검수이며 신규 코드 실행이 아니다. reviewer msg_4fbcac80f76c의 deliver:human 인증 경계 high는 DEV 수정·독립 재검증 대상이다. 이를 UI PASS로 해소하지 않는다. DEV 수정 후보가 UI에 영향을 주면 coor가 같은 과제 후속을 전달한다. 제품/기획/PLANS/board는 수정하지 않는다.

## 완료 보고

### 변경과 판정

기존 실패 Dispatch의 D04·시각 판정·실행 기록·보완 PNG를 보존해 완성했다.
이번 Dispatch `ctx_f6be2c2d0ce4`에서 기존 PNG 7개와 보완 PNG 1개를 `view_image`로 직접 확인했다.
제품 후보 `a6a10c71977b7f3ec8274a1fb7c8a409f58e7c92`, QA `c59537b6fa0c7e008c4c6bdba0a251dd821d4ee8`를 유지했다.
캡처된 합성 V-01–04는 PASS다. `invalid_auth`는 인증 실패 차단만 증명한다.
D04를 review 상태로 기록했다. 전체 MVP 수락이나 공개 승인이 아니다.
제품 코드·기획 정책·PLANS·board·기존 QA 증거는 수정하지 않았다.
소유 범위 밖 공유 산출물 인덱스의 D04 미작성 이력은 보존했다. coor에게 review 원천과의 대조를 인계한다.

### 검증과 한계

문서 커밋 `26e601d8d69394602945d9333acf8bd32a22c536`에서 `lint.py --from 0a5b044`를 수행했다.
첫 실행은 `@types/node` 부재로 종료코드 1, ERROR 1이다. 기존 tester와 같은 환경 문제다.
`npm ci --prefix adapters` 종료코드 0으로 Git 제외 의존성을 복구했다. 제품 소스·잠금 파일은 보존했다.
같은 SHA의 재실행은 종료코드 0, product-lint passed, ERROR 0, WARNING 0, 실행 불가 0이다.
전체 strict는 종료코드 0, 검사 13, 미작성 4, 문제 0, 경고 0이다.
공유 인덱스 상태 때문에 strict가 새 D04를 건너뛰므로 D04 자체의 정규 메타데이터를 별도로 검사했다.
변경 문서·인박스의 메타데이터와 로컬 링크, 원래 QA PNG 7개 바이트 일치 검사는 통과했다.
제품 후보 대비 제품 diff와 git diff --check 종료코드는 0이다.
검증 증거는 `/tmp/SAR-MVP-001-UI-validation/ctx_f6be2c2d0ce4/`다. 직접 종료코드를 유지했다.
서버 QA·build/test/runtime·새 캡처·영상·모바일·키보드·스크린리더 검사는 반복하지 않았다.
원래 캡처와 저장 HTML의 정지 렌더를 새 서버 실행 증거로 바꾸지 않았다.
완료 기록 아카이브 후 최종 커밋을 만들고 같은 lint·strict·공백 검사를 실행한다.
최종 고정 SHA와 최종 검사 결과는 새 preamble worker_done으로 전달한다.

### 미해결 경계와 인계

reviewer `msg_4fbcac80f76c`의 `deliver:human` 인증 경계 high는 미해결로 유지한다.
이번 시각 PASS는 해당 high를 해소하지 않는다. DEV 수정·독립 재검증 전 제품 수락·병합을 차단한다.
DEV 수정이 화면에 영향을 주면 coor가 같은 과제의 후속 검수를 배정한다.
DEC-02·DEC-03·Free N·실adapter·A2A 현행 검토·WAL/backup 삭제·고의 stale epoch의 held를 유지한다.
D04: [화면 기록](../../docs/design-docs/mockups/SAR-MVP-001-UI.md).
판정: [직접 시각 검수](../../docs/evaluations/qa-reports/SAR-MVP-001-UI.md).
상세 결정·검증: [실행 기록](../../docs/exec-plans/phases/SAR-MVP-001-UI.md).
coor가 문서 결과의 통합과 제품 수락을 구분하고 필수 검수·병합·원격 공유를 처리한다.
사용자의 현재 과제 완료 후 중지 지시를 유지한다. 후속 기능·배포를 시작하지 않는다.
