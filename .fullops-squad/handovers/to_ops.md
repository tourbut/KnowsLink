---
title: SAR-PUBLIC-AGENTS-001-FIX-REVIEW — 최신 AGENTS 고정 후보 독립 보안 delta 리뷰
status: draft
updated: 2026-10-06
owner: ops
tasks: [SAR-PUBLIC-AGENTS-001-FIX-REVIEW]
summary: 최신 AGENTS 고정 후보 독립 보안 delta 리뷰
---

# SAR-PUBLIC-AGENTS-001-FIX-REVIEW — 최신 AGENTS 고정 후보 독립 보안 delta 리뷰

- 작성일: 2026-10-06
- From / To: coor / ops
- 상태: ready
- 승인된 범위: 로컬 비파괴 검수·격리 합성 fixture·필요한 기록 작성. 실제 메일·공개·배포·외부 계정·운영 데이터 정리·과금 금지.
- 담당 repo id: 818c78e5-d51c-4ff4-aa88-70e9ee185fbb; 워크트리 /home/shin/orca/workspaces/KnowsLink/fullops-ops, 브랜치 fullops/ops.
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀: coor /home/shin/orca/workspaces/KnowsLink/fullops-coor / term_6895aaf1-7b43-4fe0-a416-76f1255a5946 / run_8ca8bc058ab7. task/dispatch는 preamble 실제 값.

## 현재 상황과 확인 근거

고정 검수 후보 `458798c2ee15c179edacfd6f94ebb9896d26f411`. 원본 DEV d1eef9b, 원본 OPS 70f26bc, 원본 QA bcb06b8, 원본 UI d165178 FAIL, DEV-FIX 4a1b80a, 제품 답48d12fa, DEV-POLICY-FIX83e0bfb를 모두 포함한다. 제품 코드 최종6d016e5와 후보의 제품 동일성을 확인한다. 원본 실패·held·실행 SHA를 보존한다. 원본 QA/UI를 최신 SHA 실행으로 바꾸지 않는다.

## 적용 기준과 예외

fullops-common-0.3.3의 README/coding-style/testing/security, FULLOPS, project, document-writing, review/rule.json을 직접 읽는다. 제품 정본 D02 PS07·PS07-I 및 기록 보호 조건과 UX04–05, POLICY48의 두 관찰표를 적용한다. 기준은 위 fixed 후보, 예외 없음. Workers Free·기존 서버/Compose/Tunnel 유지. ponytail full. 기술 수정은 DEV, 제품 규칙 변경은 designer로 coor 경유 질문한다.

## 먼저 읽을 문서

- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md`
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md`
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` 전체 두 관찰표
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md`
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX.md`
- 원본 OPS report, TESTER report/probe-failures, UI 보고서의 관련 항목. Jev 추천은 탐색 보조이며 제외 권한이 아니다.

## 제약·완료 처리

인박스는 현재 to_ops.md만 사용한다. 역할 기록만 수정하고 제품 코드는 수정하지 않는다. PLANS에 자기 결과를 append하며 원본을 삭제하지 않는다. work.py finish로 실제 지시서/전문을 보존·빈 인박스·커밋·공유 role push 뒤 전체 40자리 SHA로 worker_done 한 번 보낸다. 완료와 main 제품 수락은 구분한다. 실패도 원문·실행 exit·재실행 사유를 보존한다. 예정 규모 보고서/시나리오/필요 증거뿐이며 SIZE/DEP 경고를 설명한다. 이미지·로그는 비밀값 없이 판정에 필요한 범위만 만든다.
## 해야 할 일과 파일 소유권

- [ ] fullops-review 및 open-code-review-delegate로 d1eef9bb90b9726149980320c42fb1fdbcaf584a..458798c2ee15c179edacfd6f94ebb9896d26f411 최신 후보의 독립 delta를 리뷰한다. 기존 d2..d1 OPS 리뷰를 원래 SHA로 재사용하고 신규 파일·변경은 모두 검토한다.
- [ ] 읽기 전용 detached clean snapshot `/tmp/knowslink-agents-review-458798c`만 읽고 결과는 현재 OPS 기록 checkout의 `SAR-PUBLIC-AGENTS-001-FIX-REVIEW-review`에 쓴다. schema independence actual reviewer session은 DEV actual sessions 01a10f52-ac0f-75a0-b253-9a926a8e5650, e1275961-b0bc-4c1d-ad9c-5b3682464a89, 2ca6140e-15e7-4c1a-8564-94f88ded0af7와 모두 달라야 한다.
- [ ] M1 최소24h 보존·정리/legacy first sweep·C1 살아 있는 agent revoked kid 재사용 금지·owner/key 기록 guards·부분 실패 rollback·rate/retry_at·restart·권한·경합을 확인한다. DEV의 `/v1/invite-decision` Generation 미결속 판단을 스스로 평가하고 공개 경로 제한만으로 요구를 충족하는지 판단한다.
- [ ] POLICY 두 표와 실제 화면 안내·원본 QA 독립 probe/exits/실패 정정 근거·UI raw PNG manifest/cleanup 증거를 확인한다. OCR 제외 증거도 무결성과 구체적 skip 영향을 기록한다. 디자인 warning/SIZE/DEP를 평가한다.
- [ ] 최신 fixed 후보의 FullOps lint --from d1을 별도 clean scratch에서 실행하고 output head가 458798c2ee15c179edacfd6f94ebb9896d26f411와 같게 기록한다. snapshot은 node_modules 생성/빌드 등으로 변경하지 않는다. 기존 make verify-mvp55 및 QA/UI의 해당 SHA 증거를 의존성 같을 때 재사용한다.
- [ ] review.py check --key SAR-PUBLIC-AGENTS-001-FIX-REVIEW --from d1eef9b --to 458798c2ee15c179edacfd6f94ebb9896d26f411 --task-key SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX 통과. report·result·lint·actual 독립성·skipped reason과 critical/high 판정을 기록한다.

## 완료 기준과 산출물

누락/pending0·미해결critical/high0·고정SHA일치·lintERROR0·실행불가0(예외는영향근거)여야 수락 가능. 실패 발견 시 coor로 정확한 재현·경로·심각도 보고한다. D12/D13 운영 부담·보존/복원 위험의 인계만 기록하며 실제 운영 자원/복원 PASS를 만들지 않는다. 제품 수치 변경 금지. 갱신할 산출물 없음(독립 리뷰 기록). UI 직접 시각 판정은 designer 책임. 새 라이브러리/테마 없음.

## 탐색·문서 분류 근거

`docs/evaluations/jev/SAR-PUBLIC-AGENTS-001-FIX-REVIEW-find.json`, `-documents-find.json`, `-context.json`을 사용한다. keep 목록은 먼저 읽을 문서와 공통 규약이다. 지시 전제와 충돌 — 먼저 확인: 원본 OPS/QA/UI 기록은 과거 d1의 M1/L1/L3·UI FAIL을 설명한다. 최신후보가그실패를해소했는지검수하며과거결과는바꾸지않는다. 실제내용이현재전제와달라진경우coor에ask한다.
