---
title: SAR-CODEX-HOOK-QUOTE-001-TESTER — Codex hook 수정의 고정 SHA 독립 QA와 delegate 리뷰
status: draft
updated: 2026-10-12
owner: tester
tasks: [SAR-CODEX-HOOK-QUOTE-001-TESTER]
summary: Codex hook 수정의 고정 SHA 독립 QA와 delegate 리뷰
attempt: 3551964c190c41388206a3cf6349c3fd
base: 67e7de9287c6c2ac6bb92beb5b31b523f7e64ec6
subagent_level: standard
test_level: lite
---

# SAR-CODEX-HOOK-QUOTE-001-TESTER — Codex hook 독립 QA와 고정 SHA 리뷰

- Purpose: review
- 상태: blocked. DEV 최종 SHA와 gate를 받기 전에는 배정하지 않는다.
- Test level: lite. Subagent level: standard. 별도 하위 구현은 필요 없다.
- From / To: coor / tester. 담당 C:/Users/shin/orca/workspaces/KnowsLink/tester, fullops/tester.
- 복귀: coor terminal term_2a0a2b18-70ad-43bd-b225-75886bd669c0, Run run_ea85dc04b5f3. 실제 task/dispatch는 주입 preamble을 따른다.
- 구현 고정 SHA / 실제 구현자 session: DEV 완료 후 coor가 연결한다.
- 리뷰 key: SAR-CODEX-HOOK-QUOTE-001-REVIEW. 기준 ref: 67e7de9287c6c2ac6bb92beb5b31b523f7e64ec6. check --task-key는 구현 과제 SAR-CODEX-HOOK-QUOTE-001-DEV다.

## 목표와 적용 기준

현재 Codex 세션의 Hook failed code1 수정과 기존 삭제 차단을 독립 판정한다. 제품·정책 변경은 없다.
fullops-common-0.3.3, FULLOPS.md, project.md, rules/common/README.md 및 coding-style/testing/security, AGENTS.md, DEV의 archived handover/phase/source/test/hook, fullops-review 및 open-code-review-delegate 스킬을 읽는다.
현재 coor가 initial config 후보를 썼으므로 coor 자기 리뷰로 대신하지 않는다. 구현자와 다른 실제 tester session에서 review.py snapshot으로 고정 SHA clean detached read_only snapshot을 만들고 결과는 tester 기록 checkout에 쓴다.

## 해야 할 일과 파일 소유권

- [ ] 고정 SHA·별도 세션·clean detached snapshot의 실제 독립성을 확인한다.
- [ ] review.py prepare의 모든 diff를 reviewed/skipped와 근거로 처리한다. 원 config의 PowerShell parser error1, 보완 후보의 cmd list 인용 회귀 expected2/actual0과 최종 수정 원인을 대조한다.
- [ ] 실제 지원 셸 실행 방식에서 root/scripts/adapters 허용0·거부2와 deny JSON을 독립 확인한다. 코드를 문자열 리터럴로 넘겨 실행 없이0이 된 것을 PASS로 취급하지 않는다. PowerShell의 native nonzero 변환과 실제 Codex host 실행 검증 범위를 구분한다.
- [ ] DEV exact SHA의 gate·등록 product-lint/test exit0·회귀 결과를 확인하고 관련 코드/환경 동일성을 근거로 재사용한다. 독립 회귀는 좁은 보안 경계만 수행한다.
- [ ] result.json/report.md/lint.json, 독립 QA 증거 및 tester phase·contexts/tester·packet outcomes·인박스를 작성하고 check 통과·finish·worker_done을 수행한다.

## 완료 기준과 검증

새 critical/high0, 관련 회귀 PASS, 같은 고정 SHA 리뷰 check PASS. check는 실제 기준/고정 head와 --task-key SAR-CODEX-HOOK-QUOTE-001-DEV를 사용한다. 검사 통과와 사용자 host에서 실제 오류 소멸은 구분한다. 원 host stderr가 확보되지 않으면 그 한계를 명시한다.
원 실패·미검증·잔여 한계는 보존한다. 전역 hook 신뢰 승인과 캐시 변경을 하지 않는다. hook을 끄거나 보호 규칙을 우회하지 않는다.
UI·의존성·D01–D13 변경 없음. 시각 검수 해당 없음.

## 제약과 기대 산출물

제품 source/config/test 수정은 금지한다. 오류 발견 시 coor에 ask로 수정 필요를 전달한다. 실제 삭제·이동·cleanup과 사고 명령 실행 금지. 삭제 문자열은 guard JSON 입력으로만 전달한다. 임시 테스트·snapshot은 보존하고 경로를 보고한다.
수정 가능: 이 과제 QA/리뷰 폴더·phase·contexts/tester·현재 tester 인박스·outcomes. PLANS/board와 main 통합·push/동기화는 coor 책임이다.
산출물: docs/evaluations/qa-reports/SAR-CODEX-HOOK-QUOTE-001-REVIEW-review/, docs/exec-plans/phases/SAR-CODEX-HOOK-QUOTE-001-TESTER.md.

## 완료 보고

실제 reviewer session, 고정 head/base·snapshot·check, 재사용 gate, 독립 회귀, findings와 host 미검증을 기록한다. work.py finish 뒤 preamble worker_done을 직접 보낸다.
