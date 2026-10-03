---
title: SAR-MVP-002-DEV-REVIEW — F-01 리뷰 위치 참조를 실제 고정 후보에 맞게 보완한다
status: draft
updated: 2026-10-03
owner: ops
tasks: [SAR-MVP-002-DEV-REVIEW]
summary: F-01 리뷰 위치 참조를 실제 고정 후보에 맞게 보완한다
---

# SAR-MVP-002-DEV-REVIEW — F-01 위치 참조 보완

- 상태: ready
- From / To: coor / ops
- 기준 ref: dbdd70086971285b790683f362702e5a9ff55acd
- 제품 후보: 552586b6e886f95bffa9a000a031ea03070afedb
- 복귀: Run run_8ca8bc058ab7. 실제 Task/Dispatch는 새 preamble을 따른다.

## 적용 기준과 예외

기존 SAR-MVP-002-DEV-REVIEW의 규칙·공식 근거·same read-only snapshot /tmp/knowslink-plugin-review-552586b와 세션 독립성을 재사용한다. fullops-common-0.3.2, 문서 규칙 및 fullops-review/open-code-review-delegate를 따른다. 제품과 기존 판정 범위는 변경하지 않는다.

## 먼저 읽을 문서

- docs/evaluations/qa-reports/SAR-MVP-002-DEV-review/result.json, report.md
- snapshot의 adapters/src/mcp.ts, core.ts
- docs/exec-plans/phases/SAR-MVP-002-DEV-REVIEW.md
- handovers/logs/2026-10-03_to_ops.md의 직전 완료 전문

## 해야 할 일과 파일 소유권

- [ ] F-01은 mcp.ts:268-283로 표시됐지만 실제 파일은 67줄이다. mcp.ts:47-60의 adapter.once await와 core.ts:170-182의 gate 대기를 대조해 report/result의 정확한 주 경로·줄과 필요한 연결 위치를 수정한다.
- [ ] 다른 findings 줄도 고정 snapshot에서 유효한지 확인한다. 기록-only 보완이며 전체 리뷰/빌드/QA를 반복하지 않는다.
- [ ] 정확한 기존 refs로 review.py check를 통과시키고 실행 기록에 보완 근거를 남긴다.
- [ ] 현재 인박스 완료 전문, work.py finish·커밋·기준 lint·역할 push 후 새 worker_done을 보낸다.

소유는 기존 리뷰 result/report, ops 상세 실행·컨텍스트와 정규 인박스/완료로그다. 제품 코드·PLANS/board·다른 역할 기록 변경 금지.

## 완료 기준과 검증

F-01 위치 참조가 실제 고정 후보에 유효하고 check 통과 및 모든 줄 범위가 확인돼야 한다. 제품 SHA와 severity·수락 결론의 변경 여부를 명시한다. 코드·제품 불변 증거를 재사용한다.

## 갱신할 산출물

없음.

## 완료 보고

첫 줄 [완료] SAR-MVP-002-DEV-REVIEW | SHA <전체 보완 SHA>. 정확한 위치·제품 불변·check/lint 종료코드·결론을 보고한다.
