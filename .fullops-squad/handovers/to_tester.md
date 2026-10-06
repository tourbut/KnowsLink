---
title: SAR-PUBLIC-IDENTITY-001-FIX-TESTER — 최종 신원 후보의 rate 격리 및 trial 검사 격리 변경 영향 좁은 독립 QA
status: draft
updated: 2026-10-06
owner: tester
tasks: [SAR-PUBLIC-IDENTITY-001-FIX-TESTER]
summary: 최종 신원 후보의 rate 격리 및 trial 검사 격리 변경 영향 좁은 독립 QA
---

# SAR-PUBLIC-IDENTITY-001-FIX-TESTER — 최종 신원 후보의 rate 격리 및 trial 검사 격리 변경 영향 좁은 독립 QA

- From / To: coor / tester. 상태: ready.
- 기록 checkout: `/home/shin/orca/workspaces/KnowsLink/fullops-tester`, `fullops/tester`, repo `818c78e5-d51c-4ff4-aa88-70e9ee185fbb`.
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

- [ ] 별도 detached 고정 후보 `/tmp/knowslink-identity-final-qa-eb2e34b`를 실행하고 검증 전후 추적 트리·HEAD가 동일한지 확인한다.
- [ ] RATE-FIX 변경 영향만 독립 검증한다. 거부된 동일 principal이 공유 예산을 고갈하지 않는지, 다른 principal·기존 회원·cleanup 접근과 shared 한도·회전 시 상태 상한을 확인한다. 기존 테스트는 직접 실행하며 구현자 보고만 복사하지 않는다.
- [ ] TestTrialHTTP 원인과 변경된 verify-mvp 격리 순서를 확인한다. 자기 격리 DB에서 새 결정적 foreign-allowlist 회수 검사와 기존 Trial·관련 회귀를 실행한다. 테스트 완화·우연한 재실행 PASS로 간헐 실패 해소를 선언하지 않는다.
- [ ] make verify-mvp의 stop relay → Go integration → up --wait relay → TS 검사와 정리를 확인한다. 고정 후보 product-lint/product-test 증거를 남긴다. 명령 자신의 종료코드를 보존한다.
- [ ] 원본 QA/UI의 변경 없는 의존성은 동일성을 확인해 원래59 SHA로 재사용한다. 전체157개 fixture 검사를 새 과제로 반복하지 않는다. 실제 이메일/공개/사람 확인은 미실행으로 유지한다.
- [ ] 시나리오·QA 보고서·상세 실행 기록과 필요한 최소 probe/증거만 기록하고 finish·커밋·worker_done으로 복귀한다.

TESTER는 자기 시나리오·QA 보고서/최소 검사 스크립트·실행 기록·인박스·컨텍스트만 수정한다. 제품·원본 증거는 수정하지 않는다. 검증할 API/한도 수치는 원본 D02와 현재 구현의 확정 규칙을 따른다.
완료 조건은 좁은 변경 영향 QA의 독립 통과/실패 판정과 fixed `eb2e34b93fe8d20fa1cd9166f73ff68d14bf17de`의 직접 실행 증거다. 필수 실패와 critical/high는 수락을 차단한다. 사용자 실제 이메일/공개/최종 노우↔다닷은 후속 운영 수락이다.
갱신할 D01–D13 추천은 없음이다. 상세 기록은 `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-FIX-TESTER.md`다.

고정 eb2e34b의 coordinator lint/product-lint/product-test 통과 증거는 `.fullops-squad/docs/evaluations/SAR-PUBLIC-IDENTITY-001-final/candidate-lint.json`에 있다. 등록 명령은 같은 HEAD에서 실행했고 종료코드는 모두0, ERROR0/WARNING1(기존 PLANS 길이)/실행불가0이다.
