---
title: SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW — 최종 lease 검사 격리 수정과 결과 기록의 독립 delta 리뷰
status: draft
updated: 2026-10-06
owner: ops
tasks: [SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW]
summary: 최종 lease 검사 격리 수정과 결과 기록의 독립 delta 리뷰
---

# SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW — 최종 lease 검사 격리 수정과 결과 기록의 독립 delta 리뷰

- From / To: coor / ops. 상태: ready.
- 기록 checkout: `/home/shin/orca/workspaces/KnowsLink/fullops-ops`, `fullops/ops`, repo `818c78e5-d51c-4ff4-aa88-70e9ee185fbb`.
- 복귀 Run `run_8ca8bc058ab7`, coordinator `term_6895aaf1-7b43-4fe0-a416-76f1255a5946`. Task/Dispatch/worker handle은 실제 착수 preamble이 정본이다.

## 현재 상황과 적용 기준

고정 최종 후보는 `eb2e34b93fe8d20fa1cd9166f73ff68d14bf17de`다. lint 기준은 `9c915dc71e2a872243ffec294126d4668b4d32a4`다. 공통 기준은 후보의 fullops-common-0.3.3과 연결된 세 규칙 및 project.md다. 요구사항과 보안·검증 기준을 완화하지 않는다.
원본59 구현 독립 리뷰25b110f와 fixed9c915dc RATE/QA/UI 독립 리뷰77dd464를 재사용한다. 원본 TESTER9e2654d/UIcf0ab09는 fixed59 결과로만 연결한다.
DEV00a1384는 verify-mvp의 실행 relay Cleanup이 서로 다른 allowlist의 Go integration trial lease를 회수하는 원인을 규명했다. 제품 비시험 코드는 변경하지 않고 Go 검사 중 relay stop, 뒤 up --wait로 격리했다. 코드2111ff4와 완료00a1384의 제품 내용은 동일하다. 원래 invalid_lease 실패를 보존한다.

## 먼저 읽을 문서

- `.fullops-squad/FULLOPS.md`, `project.md`, `rules/common/README.md`와 연결된 세 규칙.
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG.md`와 `SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX.md`.
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW-review/report.md`, result.json, test-run.json.
- 원본 `SAR-PUBLIC-IDENTITY-001-TESTER.md`, `SAR-PUBLIC-IDENTITY-001-UI.md`와 필요한 요약 증거.
- `scripts/verify_mvp.py`, `internal/relay/test_messages_integration_test.go`, `internal/relay/identity.go`와 관련 rate 검사, `Makefile`, `README.md`의 실제 실행 설명.

## 제약·협업·후속

Workers Free만 허용한다. 유료 전환·구독·초과 과금은 금지한다. 기존 서버·Tunnel·사용자 자원은 변경하지 않는다. 실제 이메일·운영 공개·배포·Cloudflare 쓰기·최종 노우↔다닷 시험은 수행하지 않는다. 실제 확인은 전체 일반 서비스 수락의 운영 후속이며 로컬 코드 병합의 추가 선행 조건으로 확대하지 않는다.
UI 템플릿·상태 문구는 변경되지 않았다. 원본 직접 시각 검수는 의존성 동일성을 확인해 원래 SHA로 연결한다. 새 시각 검수로 표시하지 않는다. 디자인 lint와 공용 테마 전환은 미구성이며 이번 코드/UI 변화가 없어 해당 없음이다.
coor PLANS/board·원본 QA/UI/리뷰·제품 코드·운영 설정은 수정하지 않는다. 필요한 수정은 finding이나 preamble ask/escalation으로 보고한다. 지시서 범위의 비파괴 검증은 재승인 없이 완료한다. 원본 실패와 미실행을 PASS로 재작성하지 않는다.

## 완료 보고

고정 SHA·기준·범위·검증 명령/종료코드·실패/미실행·lint HEAD/ERROR/WARNING/실행 불가·산출물·남은 low/후속을 전문으로 쓴다. SIZE/DEP 경고는 근거를 남긴다. finish로 지시서와 결과 전문을 logs에 보존하고 인박스를 비운다. 커밋 뒤 preamble의 worker_done을 한 번 보내고 idle한다.

## 해야 할 일과 파일 소유권

- [ ] fullops-review/open-code-review-delegate로 리뷰 key SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW, from `9c915dc71e2a872243ffec294126d4668b4d32a4`, to `eb2e34b93fe8d20fa1cd9166f73ff68d14bf17de`를 prepare한다.
- [ ] 별도 read-only detached snapshot `/tmp/knowslink-identity-trial-review-eb2e34b`에서 delta 전체를 검토한다. snapshot에서 설치·빌드·테스트·파일 수정을 하지 않는다. 실행 검증은 별도 scratch clone을 사용한다.
- [ ] 검사 격리 변경·새 결정적 회귀·README/D10·진단 증거·RATE 리뷰 결과 기록과 최신 운영 인계의 정합성을 검토한다. 원본 범위는 재리뷰하지 않는다.
- [ ] 구현자와 다른 실제 reviewer 세션 ID 및 snapshot head/read_only를 independence에 기록한다. 실제 DEV 세션은 worker-show나 이번 진단 provider 기록에서 확인한다.
- [ ] 모든 preview 항목을 reviewed/skipped로 기록하고 증거·lint/test·경고·check를 확인한다. 고정 후보의 기존 lint/test가 HEAD·조건에 맞으면 재사용하고 재실행하지 않은 검사를 PASS로 바꾸지 않는다.
- [ ] 보고서와 실행 기록, 완료 전문을 작성하고 finish·커밋·worker_done으로 복귀한다.

OPS는 자기 리뷰 디렉터리와 실행 기록·인박스·컨텍스트만 수정한다. 고위험 lease/동시성 경로의 검증이므로 Jev Sonnet 추천 대신 fullops-review에 따라 새 Opus5.5 high 세션을 쓴다.
완료 조건은 전체 커버리지·독립성·고정 refs check와 latest 후보 delta의 수락 판정이다. critical/high·필수 실패·누락된 증거는 차단한다. narrow QA는 별도 tester 결과와 결합한다.
갱신할 D01–D13은 없음이다. 상세 기록은 `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW.md`다.

고정 eb2e34b의 coordinator lint/product-lint/product-test 통과 증거는 `.fullops-squad/docs/evaluations/SAR-PUBLIC-IDENTITY-001-final/candidate-lint.json`에 있다. 등록 명령은 같은 HEAD에서 실행했고 종료코드는 모두0, ERROR0/WARNING1(기존 PLANS 길이)/실행불가0이다.
