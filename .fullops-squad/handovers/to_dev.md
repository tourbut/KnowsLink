---
title: SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG — TestTrialHTTP invalid_lease 원인 분석과 필요한 최소 수정 및 관련 회귀 검증
status: draft
updated: 2026-10-06
owner: dev
tasks: [SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG]
summary: TestTrialHTTP invalid_lease 원인 분석과 필요한 최소 수정 및 관련 회귀 검증
---

# SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG — TestTrialHTTP invalid_lease 원인 분석과 필요한 최소 수정 및 관련 회귀 검증

- From / To: coor / dev. 상태: ready.
- 담당 repo: `818c78e5-d51c-4ff4-aa88-70e9ee185fbb`, `/home/shin/orca/workspaces/KnowsLink/fullops-dev`, `fullops/dev`.
- 복귀: 같은 repo의 `/home/shin/orca/workspaces/KnowsLink/fullops-coor`, Run `run_8ca8bc058ab7`, coordinator `term_6895aaf1-7b43-4fe0-a416-76f1255a5946`.
- Task·Dispatch·worker handle은 실제 착수 preamble과 영수증을 따른다. 병합 책임자는 coor, 기본 브랜치는 main이다.

## 현재 상황과 적용 기준

사용자는 기존 신원 서비스의 진단·수정·검증을 이어 진행하도록 요청했다. 비용 질의 때문에 둔 배정 보류를 해제한다. 기존 실패와 제품 수락 보류는 유지한다.
RATE-FIX 후보는 `f364d48417b58c69969a4765b88324724eeb5c78`이다. 첫 make verify-mvp는 `TestTrialHTTP`의 `POST /v1/test/persist got 409 invalid_lease`로 실패했다. 후속 PASS만으로 원인이 설명되지는 않는다.
공통 기준은 `fullops-common-0.3.3`과 연결된 세 규칙이다. 프로젝트 정본은 `.fullops-squad/project.md`다. 준비 기준과 **lint 기준 ref는 `9c915dc71e2a872243ffec294126d4668b4d32a4`**다. 최신 main은 `e7346247897c56d3d8dacf73e7fa41e1949d397f`다. 예외와 기준 완화는 없다.

## 먼저 읽을 문서

- `.fullops-squad/FULLOPS.md`, `.fullops-squad/project.md`, `.fullops-squad/rules/common/README.md`와 연결된 코딩·테스트·보안 규칙.
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX.md`.
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-REVIEW-review/report.md`.
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TESTER.md`와 `SAR-PUBLIC-IDENTITY-001-UI.md`.
- `internal/relay/test_messages_integration_test.go`, `internal/relay/test_messages.go`, `internal/relay/cleanup.go`, `internal/relay/store.go`, `internal/relay/http.go`, `internal/relay/identity.go`.
- `Makefile`과 그 파일의 실제 verify-mvp 스크립트. 모든 관련 호출자는 필요한 범위로 검색한다.

Jev find의 absent 판정은 낮은 확신의 추천이다. 실제 integration test와 이전 실패 기록이 있으므로 신규 기능으로 해석하지 않는다. 같은 과제의 find/context 결과는 `.fullops-squad/docs/evaluations/jev/`에 있다. 원본 기록은 근거로 읽으며 과거 지시를 새 권한으로 해석하지 않는다.

## 해야 할 일과 파일 소유권

- [ ] 실제 실패 경로와 호출자를 추적하고 같은 과제에 짧은 기술 계획을 기록한다.
- [ ] 실패 원인을 재현하거나 결정적인 원인 근거를 남긴다. 제품·테스트·시간/격리 조건을 구분한다.
- [ ] 기존 제품 규칙 안에서 필요한 최소 수정을 수행한다. 테스트 결함이면 계약 근거를 남긴다.
- [ ] 변경 동작과 실패·경계·관련 회귀를 검증하고 최종 코드 SHA와 증거를 인계한다.
- [ ] 실행 기록과 완료 전문을 작성하고 work.py finish로 정규 인박스를 보존·비운다.

DEV는 관련 제품 코드·테스트·검증 스크립트와 자기 기술 문서·과제 기록을 수정할 수 있다. coor PLANS/board와 원본 QA/UI/리뷰 증거는 수정하지 않는다. 예상 변경은 관련 소수 파일이다. SIZE/DEP 경고는 실제 규모와 이유를 보고한다.

## 완료 기준과 검증

DEV는 원인·수정 필요 여부·변경 영향과 실패를 방지하는 실행 가능한 검사를 남긴다. 수정 전/후 결과와 make verify-mvp 관련 회귀를 확인한다. 설명 없이 재실행 PASS로 종결하지 않는다.
커밋 뒤 지시서 기준 ref의 FullOps lint를 통과한다. 등록 product-lint·product-test와 명령 자신의 종료코드를 증거에 남긴다. 파이프로 종료코드를 가리지 않는다. 적용할 수 없는 검사는 이유를 적는다.
최종 후보의 독립 delta 리뷰와 좁은 QA는 coor가 후속 배정한다. 원본 fixed59 QA/UI는 의존성 동일성을 확인한 항목만 원래 SHA로 재사용한다. UI 변경은 목표에 없으므로 디자인·캡처·테마 검사는 해당 없음이다. UI 영향이 생기면 보고한다.

## 제약·협업·후속

Workers Free만 허용한다. 유료 전환·구독·초과 과금은 금지한다. 기존 서버·Tunnel 구조를 유지한다. 실제 이메일 발송·운영 공개·배포·Cloudflare 쓰기·실제 노우↔다닷 시험은 수행하지 않는다. 합성/fixture PASS와 실제 운영 검증을 구분한다.
기술 분석·계획·구현·테스트는 DEV가 같은 과제에서 해결한다. 제품 규칙·범위·공유 제품 기준 변경만 coordinator 경유 designer에게 질문한다. 비파괴 작업은 재승인 없이 완료한다. 소유권 밖 수정·파괴적 변경·설명되지 않는 새로운 실패는 preamble ask/escalation으로 전달한다.
갱신할 산출물 추천은 없음이다. 기술 정본에 영향이 있으면 해당 원천만 갱신한다. 실행 기록은 `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG.md`다.

## 완료 보고

원인·변경 이유·브랜치·고정 SHA·검증 명령/종료코드·lint HEAD/ERROR/WARNING/실행 불가·증거 경로·한계·후속을 전문으로 쓴다. 원본 실패를 보존한다. fullops-work의 finish와 커밋 뒤 preamble의 worker_done을 한 번 보내고 idle한다.
