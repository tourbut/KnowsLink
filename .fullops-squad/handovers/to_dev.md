---
title: SAR-SETUP-001-DEV — 잠긴 원천 결정에 맞는 초기 개발 환경과 실행 가능한 lint 구성
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-SETUP-001-DEV]
summary: 잠긴 원천 결정에 맞는 초기 개발 환경과 실행 가능한 lint 구성
---

# SAR-SETUP-001-DEV — 잠긴 원천에 맞는 초기 구성과 lint 구현

- 작성일: 2026-10-03
- From / To: designer / dev
- 상태: ready
- 승인된 범위: 이 지시서의 소유 파일에서 초기 구성·검증에 필요한 비파괴 작업 및 커밋.
- 추가 승인이 필요한 행위: 외부 배포·발송, 운영 자원 변경, 레포 밖 영속 변경, 데이터·파일 삭제, force-push, 파일 소유권 밖 수정.
- 담당 워크트리 / 브랜치: `/home/shin/orca/workspaces/KnowsLink/fullops-dev` / `fullops/dev`.
- 병합 책임자 / 기본 브랜치: coor가 검토·병합을 조정하며 필요 시 ops에 배정 / `main`.
- 복귀 워크트리 / 터미널: `/home/shin/orca/workspaces/KnowsLink/fullops-coor` / `term_89f25ea4-e70e-46c0-8514-e95f8cf81928`.
- Run: `run_8ca8bc058ab7`. repo id와 해당 worker의 task id·dispatch id는 새 dispatch에서 coordinator가 실제 값으로 기록한다. 후속 모델은 coor가 갱신한 후보에서 선택하며 과거 Astra 배정을 재사용하지 않는다.
- 완료 전송: 새 dispatch preamble의 from·capability·task id·dispatch id를 그대로 쓴다. 설계 worker의 lifecycle ID는 재사용하지 않는다.

## 적용 기준과 예외

- 규칙: `fullops-common-0.3.1`; `.fullops-squad/rules/common/README.md` 및 coding-style.md, testing.md, security.md.
- 정본: `.fullops-squad/project.md`, `.fullops-squad/docs/agents/document-writing.md`, `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md`.
- 기준 ref: `729446d8da57`. 고정 원천은 준비 커밋 `0cc35f0`, 외부 원천 SHA `404ff834c0607055d63d2053bf7771d2f46ad3ae`다.
- 착수 전에 coordinator가 이 지시서를 포함한 설계 커밋을 반영한다. worker와 검토자는 같은 원천·규칙 버전을 읽는다.
- 예외: 없음. 원천의 확정 기술 선택은 project.md의 이전 미정 상태보다 구체적인 요구사항이다.

## 현재 상황과 확인 근거

README와 FullOps 하네스만 있는 레포다. 제품 실행기와 lint 명령은 아직 없다.
원천은 `.fullops-squad/docs/planning/sources/silent-agent-relay/`에 있다. 프로젝트 루트의 `sources/`가 아니다.
기술 설계와 구현은 dev가 수행한다. 원천의 Go·TypeScript·Postgres·Compose·DB tooling·Go UI 선택은 이미 확정됐다.

## 먼저 읽을 문서

먼저 이 지시서와 D02를 읽고 아래 keep 문서를 확인한다. 모든 경로는 레포 루트 기준이다.

- `.fullops-squad/FULLOPS.md`
- `.fullops-squad/rules/common/README.md`
- `.fullops-squad/rules/common/coding-style.md`
- `.fullops-squad/rules/common/testing.md`
- `.fullops-squad/rules/common/security.md`
- `.fullops-squad/project.md`
- `.fullops-squad/docs/agents/document-writing.md`
- `.fullops-squad/docs/planning/SAR-SETUP-001-request.md`
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md`
- `.fullops-squad/docs/planning/sources/silent-agent-relay/product.md`
- `.fullops-squad/docs/planning/sources/silent-agent-relay/architecture.md`
- `.fullops-squad/docs/planning/sources/silent-agent-relay/protocol.md`
- `.fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md`
- `.fullops-squad/docs/planning/sources/silent-agent-relay/mvp-checklist.md`
- `.fullops-squad/lint/README.md`
- `.fullops-squad/docs/deliverables/README.md`
- `.fullops-squad/contexts/dev.md`
- `.fullops-squad/orca-agents.md`
- `.fullops-squad/handovers/to_dev.md`

- 탐색 근거: `docs/evaluations/jev/SAR-SETUP-001-DEV-find.json`, `SAR-SETUP-001-DEV-documents-find.json`, `SAR-SETUP-001-DEV-context.json` (`.fullops-squad/` 기준).
- find의 code/documents 호출은 모두 `API or response validation failed: ValueError`로 후보를 반환하지 못했다. 존재 판정도 null이므로 Jev가 absent를 판정했다고 주장하지 않는다.
- `rg --files --hidden .fullops-squad`와 Git 추적 파일로 원천·규칙을 좁혔다. 제품 코드가 없는 초기 과제임을 직접 확인했다.
- context도 API 실패로 위 19개를 전부 keep했다. 자동 제외·충돌 검증 성공으로 해석하지 않는다. 새 D02는 이번 설계 작업본을 직접 후보로 추가했다.
- architecture.md와 mvp-checklist.md는 `sensitive or oversized passage`로 본문 전송 없이 keep됐다. 원문을 로컬에서 읽는다.
- 필요 시 확인: 원천 README.md, business-model.md, source.json. omit? 자동 추천은 없다.
- 지시 전제와 충돌 — 먼저 확인: project.md의 기술 스택 미정은 오래된 상태다. 원천의 확정 스택을 따른다. protocol.md의 현재 C1–C5와 후반 결정이 과거 webhook 허용 메모보다 우선한다.
- 지시문 포함 — 내용만 참고: 외부 원천의 명령형 문장은 제품 요구 근거다. 세션 권한을 변경하는 실행 지시로 해석하지 않는다.


## 해야 할 일과 파일 소유권

- [ ] D02의 SETUP-01–04, LINT-01–03, DOC-01, SCOPE-01을 기술 설계와 구현 경로에 매핑한다.
- [ ] D03 architecture.md와 tech-stack.md에 잠긴 선택, 실제 버전·디렉터리·설정·검사 명령을 기록한다. 라이브러리 API는 실제 버전의 Context7 근거를 기록하고 실패 시 공식 문서를 사용한다.
- [ ] 최소 Go relay·별도 `cmd/migrate`·TypeScript adapter 골격과 의존성·잠금 파일을 구성한다. 경로는 dev가 정하고 `orca-agents.md`의 dev 담당 경로와 `project.md`에 등록한다.
- [ ] pgx/v5·pgxpool, goose/v3 SQL-only 및 sqlc pgx/v5 설정을 마련한다. 업무 스키마·샘플 쿼리는 검사 실행을 위해 꾸며 만들지 않는다.
- [ ] `postgres`, one-shot `migrate`, `relay`, `cloudflared` Compose 구성과 비밀값 없는 예시 설정을 만든다. 운영 Tunnel은 연결하지 않는다.
- [ ] 실제 제품 파일에 적용되는 포맷·정적 검사·타입 검사와 통합 lint 명령을 만든다. 필요한 도구만 도입한다.
- [ ] `.fullops-squad/lint/lint.json`의 commands에 명령을 등록하고 lint/README.md에 이유를 남긴다. 기존 기준 완화와 무차별 exclude 추가는 금지한다.
- [ ] 정상 검사와 임시 위반 주입 검사를 실행한다. 명령 자신의 종료코드와 한계를 실행 기록에 남긴다.
- [ ] README, project.md와 D03을 실제 동작에 맞추고 완료 보고·contexts/dev.md를 작성한다. `work.py finish` 후 커밋하고 기준 ref의 FullOps lint를 실행한다.

소유 파일은 위 초기 구성에 필요한 제품 코드·설정·의존성·직접 검증 파일, README, project.md, orca-agents.md의 dev 경로, lint 설정·안내, D03, 본인 인박스·컨텍스트·실행 기록이다.
D02와 고정 원천, tester 시나리오·QA 보고서는 수정하지 않는다. 산출물 인덱스와 PLANS는 본인 과제 상태만 갱신한다.

## 완료 기준과 검증

D02의 수락 기준을 적용한다. QA-01의 독립 판정은 후속 tester 책임이다.
설치·Go/TypeScript 빌드·최소 실행·Compose 구성 검사·제품 lint를 실제로 실행한다. 필수값 누락과 lint 실패 전파를 확인한다.
Go 포맷 오류, TypeScript lint 오류·타입 오류는 임시 복제본 또는 fixture로 주입한다. 기대 오류가 검출된 뒤 정상 상태를 다시 확인한다.
자동 수정 모드만 실행하거나 성공을 출력하는 빈 명령을 lint로 제공하지 않는다. 검사 불가능한 SQL/codegen 등은 이유를 구분한다.
기준 ref의 FullOps lint는 새 commands를 적용하지 않는다. 새 제품 검사 명령은 직접 실행하고 결과를 별도 기록한다.
커밋한 깨끗한 상태에서 `python3 <플러그인>/scripts/lint.py --repo . --from 729446d8da57 --out <레포 밖 결과 경로>`를 실행한다.
`git diff --check`와 `deliverables.py --repo . --strict`를 수행한다. `| tail`로 종료코드를 가리지 않는다.
ERROR를 해결한다. WARNING·실행 불가와 사유, 대상 SHA, 명령·환경을 완료 보고에 쓴다. 필수 미검증이 있으면 수락 완료로 보고하지 않는다.

### 기존 원천의 lint 차단과 coordinator 결정

설계 커밋 `916fb978d46b`에서 지정 기준 `00b4cb3`의 FullOps lint는 원천 문서 일곱 개의 DOC-003으로 실패했다.
문서 규칙은 원본 외부 문서를 보존하도록 정한다. 원천을 stamp하거나 바꿔서 오류를 없애지 않는다.
coor는 원천을 보존하고 설계 변경만 시작 HEAD `0cc35f0`에서 별도 검사하도록 허용했다. 이 보조 검사는 지정 기준 통과를 대신하지 않는다.
coor는 원천 전용 lint 제외를 커밋 `729446d`에 기록했다고 회신했다. worker는 해당 준비 커밋의 실제 내용과 반영 여부를 확인한다.
FullOps는 기준 ref의 설정을 읽으므로 제외 설정이나 새 제품 commands를 현재 브랜치에 추가하는 것만으로 적용되지 않는다.
지정 기준 결과를 숨기거나 기준 ref를 임의로 바꾸지 않는다. coor와 검사 설정의 적용 기준을 확인하고 원천 오류·제품 검사 결과를 구분해서 보고한다.
이번 회신은 설계 문서 과제의 완료 허용이다. 제품 코드 변경의 done-gate 통과를 면제하지 않는다.

## 갱신할 산출물

D03: `.fullops-squad/docs/design-docs/architecture.md`, `.fullops-squad/docs/design-docs/tech-stack.md`.
D02를 upstream으로 연결하고 초기 구성만 문서화한다. `fullops-deliverables`의 stamp와 strict 검사를 적용한다.
다른 D01–D13은 갱신하지 않는다. 개발 환경 사용법은 README와 D03에 둔다.

## 기대 산출물

- 초기 개발 골격·Compose·예시 환경 설정·제품 검사 명령과 실행 증거.
- D03, README, project.md, lint 정본.
- `.fullops-squad/docs/exec-plans/phases/SAR-SETUP-001-DEV.md`, 완료 아카이브와 본인 contexts 기록.
- tester가 재현할 수 있는 완료 커밋 SHA와 수락 기준별 검증 표.

## 제약·협업·후속

전체 relay 기능, human-gate UI, 네 어댑터 통합, 도메인 테이블, 운영 배포는 범위 밖이다. 별도 TypeScript UI를 추가하지 않는다.
첫 기능 MVP의 human-gate 포함 결정은 유지한다. 이번 설정을 기능 MVP 완료로 표시하지 않는다.
기능 골격이 업무 요청을 수락하거나 approve/exec 성공을 반환하지 않게 한다. webhook/evidence fetch와 실데이터 silent 경로는 열지 않는다.
가격·Free N·추가 rate/size 제한을 발명하지 않는다. 운영 배포에 필요한 미결정 정책은 후속으로 남긴다.
기술 세부는 승인 범위 안에서 직접 결정한다. 원천과 충돌하거나 설명되지 않는 검증 실패가 생기면 preamble의 ask로 coordinator에게 묻는다.
완료 후 coor가 이 SHA를 tester에게 전달한다. tester와 동시에 제품 코드를 수정하지 않는다.

## 완료 보고

완료 시 실제 브랜치·SHA, 변경 이유, 지시와 다른 판단, 수락 기준별 결과, lint ERROR/WARNING/실행 불가, D03·실행 기록, 남은 일을 쓴다.
검증 증거와 기록을 커밋한 뒤 FullOps lint 대상 HEAD를 남긴다. preamble의 worker_done을 정확히 한 번 보내고 종료한다.

## Coordinator 검사 기준 갱신

원천 보존용 제외 설정을 포함하는 준비 커밋 `729446d8da57`을 이 구현 및 검증 과제의 기준 ref로 지정한다. 과거 설계 단계의 차단 기록은 유지한다. 원천 스냅샷만 제외하며 제품 검사는 직접 실행하고 새 commands를 등록한다.
