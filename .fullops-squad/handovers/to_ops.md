---
title: SAR-PUBLIC-MESSAGES-001-FIX-3-REVIEW — 최종 후보의 독립 정적 코드 리뷰와 main 수락 근거
status: draft
updated: 2026-10-07
owner: ops
tasks: [SAR-PUBLIC-MESSAGES-001-FIX-3-REVIEW]
summary: 최종 후보의 독립 정적 코드 리뷰와 main 수락 근거
attempt: 1ca6bc4bed6140fe92d75ba371fdd315
base: 68b0d6a0c854fdaec6828a232dd3945814be1404
---

# SAR-PUBLIC-MESSAGES-001-FIX-3-REVIEW — 최종 후보의 독립 정적 코드 리뷰와 main 수락 근거

## 목적과 파일 소유권

최종 고정 후보의 독립 정적 코드·규약·정본 추적성을 검토한다. 결과는 본인 review 디렉터리와 PLANS/context/inbox에 쓴다. 제품 코드·기획 규칙·타인 inbox는 수정하지 않는다. 구현자와 다른 실제 세션·읽기 전용 detached snapshot을 사용한다.

## 해야 할 일과 검증

- [ ] fullops-review/open-code-review-delegate의 preview·rules·diff 전체 파일 coverage와 skipped 근거를 기록한다. main 이후 미수락 제품 변경의 실제 코드와 관련 호출자를 읽고 최종 FIX3의 상태 갱신·정리 처리·DB 원자성·동시성 및 기존 권한·CSRF·현재 자격·lease·rate/용량 계약 보존을 정적으로 판단한다.
- [ ] DEV의 원인/최소 수정/자동 회귀와 원 실패의 해소 근거를 대조한다. 기존 저장된 관측과 테스트를 이용하며 새 공격/부하 도구나 차단된 출력 재생성은 하지 않는다. 독립 실행 QA는 tester가 담당한다.
- [ ] 후보의 정확한 HEAD lint/test JSON·exit·warning·변경 규모·의존성·DESIGN 해당 여부를 확인한다. 이미 검증된 동일 SHA의 검사 근거는 재사용하며 결함/누락이 있을 때만 자기 scratch에서 필요한 검사를 보완한다.
- [ ] D01–13 인덱스에서 영향 정본의 현재 요구·결정 이유·구조·구현/미완료·검증·운영/복구·다음 담당을 대화 없이 확인한다. 새 공개 준비/D12/D13 연결은 기존 배포 검증이며 최종 공개 수락이 아님을 확인한다.
- [ ] 새 result/report의 actual independence·coverage·conclusion·findings를 작성하고 review.py check를 exact base/head 및 DEV 과제 키로 실행한다. 미해결 critical/high는 병합 차단이다. check 성공과 제품 수락은 구분한다.

## 먼저 읽을 문서와 갱신할 산출물

필수 공통 기준·제품/UX 정본·DEV-FIX-3 실행 기록/최종 검사·기존 TESTER 보고서·FIX-REVIEW report·FIX2 부분 기록·최신 UI 보고서를 읽는다. code/doc/context/packet은 최종 후보 고정 뒤 coordinator가 연결한다. D01–13 원천은 검토만 하며 잘못된 내용은 해당 담당에게 수정 요청한다. UI 직접 시각 판정은 designer가 담당한다. 기대 산출물은 이 키-review의 result.json/report.md/check 및 실제 세션/검사 provenance다.

## 대상·적용 기준·복귀

- 상태: blocked.DEV-FIX-3 최종 SHA 확정 전 착수하지 않는다. 후보 SHA와 snapshot은 coordinator가 dispatch 전에 고정한다.
- 공통 기준 fullops-common-0.3.3의 README/coding-style/testing/security, FULLOPS·project·document-writing·orca-agents·역할 context와 review/rule.json을 따른다. 기준 main은 68b0d6a0c854fdaec6828a232dd3945814be1404다. 외부 새 SDK/의존성은 없으며 기존 버전 근거를 재사용한다.
- 사용자 승인 모델: fresh Codex gpt-6.1-sol high. 과제별 지정이며 전역/역할 전체 설정·구독·추가 결제를 변경하지 않는다. 실제 세션 ID와 fixed SHA·실행 위치·명령별 종료코드를 남긴다.
- 복귀 repo 818c78e5-d51c-4ff4-aa88-70e9ee185fbb, coor /home/shin/orca/workspaces/KnowsLink/fullops-coor, terminal term_a8a1fa04-50ab-448d-94e7-11e8ee3c77f1, Run run_8ca8bc058ab7. 실제 Task/Dispatch 권한은 새 preamble을 따른다.
- 제품 기준은 docs/planning/product-specs/SAR-PUBLIC-SERVICE.md PS08–11/PS04·06·07과 docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md UX06·07이다. 원 TESTER09c 실패·OPS dfc H2 high·UI09c 실패·FIX2 리뷰 pending·플랫폼 차단은 원 SHA/시점으로 보존한다.

## 제약과 완료 기록

Workers Free·기존 서버/Tunnel을 유지한다. 운영 공개·실메일·외부 계정/플랫폼·사용자 자료 삭제·유료 전환은 범위 밖이다. 검사 실행은 자기 격리 scratch·fixture에서만 한다. 원본 snapshot은 detached clean read_only로 유지하고 기록은 역할 checkout에 쓴다. 원본 실패 파일은 수정하지 않는다. UI 화면 변경이나 제품 quota 판단은 designer에게, 제품 코드 결함은 같은 DEV 후속으로 coordinator에게 전달한다.

완료 보고 전문·실제 fixed SHA·세션·통과/실패·미실행·후속을 남긴다. work.py finish로 archive/빈 inbox, 마지막 기록 SHA에서 FullOps lint/test·strict·diff 검사, 역할 브랜치 일반 push와 실제 worker_done을 완료한다. 실패도 증거와 함께 보고하며 성공으로 바꾸지 않는다. synthetic/local 결과는 실메일/공개/실24h/노우↔다닷/운영 부하와 구분한다. 끝난 뒤 idle이며 다음 과제를 시작하지 않는다.

## 완료 보고

작업 종료 뒤 실제 결과 전문을 작성한다.
