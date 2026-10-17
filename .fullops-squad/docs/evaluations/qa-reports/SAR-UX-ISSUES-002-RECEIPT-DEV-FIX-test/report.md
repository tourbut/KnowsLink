---
title: 조회 선택기 구분값 선두 표시 DEV 결과
status: draft
updated: 2026-10-17
owner: dev
tasks: [SAR-UX-ISSUES-002-RECEIPT-DEV-FIX]
summary: D-RT-01 수정과 원격 lint·test 결과
---

# 조회 선택기 구분값 선두 표시 DEV 결과

제품 source는 `da53e961e7a1399c9c974b6be7b35e8b51947922`다. 기준은 `9ae29ec`다.

- 변경: `member_pages.go`의 조회 선택기 옵션을 `구분값 / 이름 / 종류` 한 줄로 바꿨다. 구분값을 맨 앞에 두어 39자 이름이 닫힌 상태에서 잘려도 `agent_a`와 `agent_a2`가 보인다. 토큰·44px·글꼴은 그대로다. 공용 `human-agent-label`은 다른 화면이 쓰므로 바꾸지 않았다. 충돌 확장·철회 표시·autofocus 구조는 그대로다.
- 테스트: `receipt_context_test.go`의 폼 fixture 테스트가 여섯 상태 모두에서 옵션이 구분값으로 시작하는지 검사한다. `go test ./internal/relay -run 'TestReceipt|TestMemberReceiptView|TestMemberUIFixtures'`가 통과했다.
- D10: `module-design.md`에 표시 순서를 반영했다.
- 검사: Windows에 make가 없어 정본 `lint.py`는 product-lint/product-test를 실행 불가로 보고했다(원 결과). 같은 HEAD의 Linux clone에서 [make lint](da53e961e7a1399c9c974b6be7b35e8b51947922-lint.log) EXIT=0, [make test](da53e961e7a1399c9c974b6be7b35e8b51947922-test.log) EXIT=0이다. lint 로그의 `ERROR:` 줄은 스크립트 테스트의 기대 출력이며 종료코드는 0이다.
- 못 한 것: 320/390/1440 실제 렌더 판정은 하지 않았다. tester 선택기 세 건 재검증과 designer 시각 수락이 맡는다. 새 dependency는 없다.
