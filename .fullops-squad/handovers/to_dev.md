---
title: "SAR-MVP-001-DEV — 기존 합성 MVP의 deliver:human 인증 경계 high 결함 수정"
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-MVP-001-DEV]
summary: "기존 합성 MVP의 deliver:human 인증 경계 high 결함 수정"
---

# SAR-MVP-001-DEV — 기존 합성 MVP의 deliver:human 인증 경계 high 결함 수정

- 상태: ready. 기존 SAR-MVP-001-DEV의 리뷰 후속이며 신규 기능이 아니다.
- 담당: dev / /home/shin/orca/workspaces/KnowsLink/fullops-dev / fullops/dev.
- 복귀: coor / /home/shin/orca/workspaces/KnowsLink/fullops-coor / Run run_8ca8bc058ab7. 새 preamble의 Task/Dispatch/capability만 사용한다.
- 기준 ref: e554fe1. 규칙 fullops-common-0.3.2와 project.md의 현재 제품 정본을 적용한다. 예외 없음.

## 목표·확정 규칙·범위

기존 제품 a6a10c7의 C1 인증 경계 high 결함을 고친다. reviewer 메시지 msg_4fbcac80f76c는 deliver:human schedule.query를 B-agent credential로 pull·persist·ACK·claim 처리해 transport=delivered를 만들었고 owner gate가 0개임을 보고했다. 실제 Postgres/HTTP targeted.log는 종료코드 0으로 이를 재현했다. 정본 protocol.md C1은 agent credential의 deliver:human 처리와 owner 승인 대체를 금지한다. 기술 계획·원인 분석·수정·회귀·관련 기술 문서는 dev가 같은 과제로 수행한다. 제품 규칙 변경·공개 정책 수치 확정·벤더·배포는 범위 밖이다.

## 먼저 읽을 문서

FULLOPS.md, rules/common/README.md와 coding-style.md/testing.md/security.md, project.md, docs/agents/document-writing.md, contexts/dev.md, docs/planning/product-specs/SAR-MVP.md, docs/planning/sources/silent-agent-relay/protocol.md C1–C5, docs/design-docs/architecture.md, docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md, handovers/logs/2026-10-03_to_dev.md를 읽는다. 기존 SAR-MVP-001-DEV의 find/documents-find/context 근거는 같은 코드 범위이므로 재사용하고 변경 탐색 결과를 이어 기록한다.

재현 증거 정본은 coor 체크아웃의 docs/evaluations/qa-reports/SAR-MVP-001-INTEGRATION-review/targeted.log다. coor의 미추적 reviewer 증거는 이동하거나 수정하지 않는다. 원래 고정 리뷰 head는 31405e736a16be9d77239c6cdc6fdeb56892436f이고 제품 코드는 a6a10c7과 같다. stale epoch CAS targeted 검사는 별도로 zero-row·HTTP 503 rollback을 통과했으나 기존 QA 원문은 보존한다.

## 해야 할 일과 소유권

- [ ] 실제 요청과 정본으로 결함을 재현하고 짧은 기술 계획을 같은 실행 기록에 적는다.
- [ ] agent와 human delivery 인증 경계를 바로잡고 합법적인 owner 경로·agent 전달·approval/result 흐름을 유지한다.
- [ ] 실제 Postgres/HTTP에서 agent 자격으로 human pull·persist·ACK·claim이 수락·완료되지 않음을 회귀 검증한다.
- [ ] 관련 C1–C5·lease/ACK/claim·철회·current-auth·human gate 자동 회귀와 product-lint를 수행한다.
- [ ] D10과 영향받는 D03/D05–D09 기술 정본만 실제 변경에 맞춰 갱신하고 실행 증거·자기 context·완료 아카이브를 커밋한다.

제품 코드·직접 자동 테스트와 DEV 소유 기술 문서만 수정한다. PLANS·board·기획·tester QA·독립 리뷰는 수정하지 않는다. 임시 dev-mvp 작업은 보존한다. 소유 범위의 비파괴 작업은 완료까지 진행한다. 제품 규칙 변경·범위 확대·소유권 충돌·외부 권한 문제만 ask한다.

## 완료 기준·검증·후속

agent credential로 deliver:human을 처리하거나 owner 승인을 대체할 수 없어야 한다. 실제 경계 테스트 실패와 수정 후 통과, 기존 agent/owner 정상 동작을 같은 조건에서 확인한다. make lint·관련 race/unit/실제 Postgres 회귀와 필요한 합성 실행을 수행한다. 전체 QA는 독립 tester에게 후속으로 맡기며 변경 없는 시각 자료는 재사용 조건을 확인한다. 로그에 명령 자신의 종료코드를 남기고 실패·held를 보존한다.

커밋 후 설치된 0.9.12 scripts/lint.py --repo <자기 체크아웃> --from e554fe1을 실행한다. ERROR 0·product-lint passed가 필수다. 완료한 역할 인박스에 보고 전문을 적고 work.py finish로 logs에 보존·인박스 비우기·소유파일 커밋·새 worker_done을 한 번 수행한다. 결과에 수정 SHA·재현·회귀·남은 held·후속 독립 QA/고정 SHA 리뷰를 적는다. DEV 완료는 제품 최종 수락이 아니다. coor가 독립 QA·UI·리뷰·필수 조건 확인 뒤 main 병합·일반 push·임시 워크트리 정리를 수행한다.
