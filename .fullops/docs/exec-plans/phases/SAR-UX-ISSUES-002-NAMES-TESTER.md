---
title: 이름 중심 표시 독립 QA 기록
status: draft
updated: 2026-10-19
owner: tester
tasks: [SAR-UX-ISSUES-002-NAMES-TESTER]
summary: 고정 b2e683a의 U2·U4·U7 독립 QA 범위·순서·결과 링크
---

# 이름 중심 표시 독립 QA 기록

고정 SHA b2e683ae677c78cd7bdc47c75820bcb191b94838, 기준 main 7a66f71, lite. 순서는 읽기 전용 HTTP → 브라우저 layout·keys → 비가역 변경(요청 검증·동일 응답·이름 변경·코드 교체·취소·거절·수락) → 후속 재판정 → 정리 → 변경 뒤 캡처다. 비가역 동작은 읽기 전용 검사와 캡처 뒤에 한 번만 실행했다.

시나리오는 [시나리오](../../evaluations/scenarios/SAR-UX-ISSUES-002-NAMES-TESTER.md), 결과는 [QA 보고](../../evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/report.md)에 있다. 제품 코드·정본 문서는 바꾸지 않았다. 새 결함은 low D-NM-01 하나다.
