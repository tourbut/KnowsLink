---
title: SAR-FRONTEND-POLISH-001-TESTER — 회원 UI 변경의 고정 SHA 독립 리뷰와 lite 동작 QA
status: draft
updated: 2026-10-12
owner: tester
tasks: [SAR-FRONTEND-POLISH-001-TESTER]
summary: 회원 UI 변경의 고정 SHA 독립 리뷰와 lite 동작 QA
attempt: 7e7e92f513ca440e9fa6253ec2ba265a
base: 7afaf02c6a652ef77cd3d330afc0b0ecbbceb1bf
subagent_level: standard
test_level: lite
---

# SAR-FRONTEND-POLISH-001-TESTER — 회원 UI 고정 후보 독립 리뷰와 lite QA

- Purpose: review. From/To: coor/tester. 상태: DEV 완료 SHA 대기. Test level: lite. Subagent level: standard.
- tuple: claude / claude-sonnet-5-5 / high (jev), fresh 별도 세션.
- 작업 경로 C:/Users/shin/orca/workspaces/KnowsLink/tester, fullops/tester. 복귀 run_2d134e8802eb. task/dispatch는 preamble.

## 목표·범위·규칙

사용자 요청은 기존 프런트엔드 검토·완성도 개선이다. SAR-FRONTEND-POLISH-001-DEV 완료 후 coor가 고정 SHA·구현자 실제 session ID를 추가한다. 그 전 실행하지 않는다.
구현자와 다른 실제 session ID로 review.py snapshot을 만들고 깨끗한 detached 고정 SHA를 읽기 전용으로 검토한다. 보고서·result·lint·QA 근거는 tester 기록 체크아웃에 쓴다.
제품 코드는 수정하지 않는다. 운영 계정·키·관계·메시지·로그아웃·외부 발송·배포·파일 삭제/이동은 금지한다. 기존 합성 fixture와 정확한 로컬 origin만 사용한다.
기존 디자인·브랜드·제품 규칙·서버 권한·CSRF·form 계약을 유지하는지 확인한다. 새 critical/high 또는 필수 검사 실패는 수락 불가다.

## 먼저 읽을 문서·적용 기준

.fullops-squad/FULLOPS.md, project.md, contexts/tester.md, rules/common/README.md와 coding-style.md·testing.md·security.md (fullops-common-0.3.3).
.fullops-squad/docs/design-docs/mockups/service-member-ui.md, DEV phase와 완료 로그·QA 보고서.
fullops-review, open-code-review-delegate, fullops-test, fullops-frontend를 적용한다. 주 디자인 스킬 impeccable. CLI/라이브러리 문서는 Context7을 사용한다.

## 해야 할 일·완료 기준

- [ ] 고정 SHA의 diff·호출자·검사·디자인 정본을 읽고 OCR delegate 리뷰한다. snapshot 독립성·read_only를 실제 session ID로 기록한다.
- [ ] 변경 핵심 경로·필요한 실패 경계만 독립 동작 QA한다. 1440×900/390×844, 키보드 포커스·overflow·오류/빈 상태 중 변경 영향이 있는 항목을 합성 데이터로 검증한다.
- [ ] DEV의 현재 HEAD product-lint/product-test 종료0·FullOps lint ERROR0를 확인한다. 동일 SHA·config·대상이면 정확한 근거로 재사용하고 새 실행으로 표시하지 않는다.
- [ ] 독립 review result의 모든 path/status를 reviewed/skipped와 구체적 이유로 처리한다. 대화 미참조 인계 점검과 디자인 lint·테마 해당 없음 이유를 적는다.
- [ ] review.py check 통과와 QA 결과·잔여 결함·미검사를 정본 qa-reports/SAR-FRONTEND-POLISH-001-TESTER-review/ 및 -test/에 기록한다.
- [ ] packet-outcomes·work.py finish로 인박스와 보고 전문을 보존한다. 커밋 SHA·리뷰 결론·독립 QA·증거·후속을 worker_done으로 보고한다.

## 후속·산출물

갱신 D01–D13 없음. designer 직접 시각 검수는 coor 별도 배정이다. snapshot과 임시 증거는 삭제하지 않는다.

## 완료 보고

완료 때 전문을 기록한다.
