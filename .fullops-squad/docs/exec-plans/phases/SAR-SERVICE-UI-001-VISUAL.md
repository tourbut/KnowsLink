---
title: SAR-SERVICE-UI-001-VISUAL 실행 기록
status: review
updated: 2026-10-11
owner: designer
tasks: [SAR-SERVICE-UI-001-VISUAL]
summary: 회원 웹 고정 후보 c7828b2의 합성 화면을 Orca 내장 브라우저에서 직접 검수한 실행 단계·결정·제한 기록
---

# SAR-SERVICE-UI-001-VISUAL 실행 기록

## 목적과 범위

DEV가 구현한 회원 웹(고정 후보 `c7828b227b5715f7bdf02bff46bf9cf102a758d9`)이 D04 [service-member-ui.md](../../design-docs/mockups/service-member-ui.md)의 시각 수락 기준을 채우는지 직접 확인한다. 기능 QA·코드 리뷰는 tester 담당이다. 제품 코드·DEV·tester 공간은 바꾸지 않았다.

실행자는 등록된 Claude Sonnet 5.5 high 새 세션이다. Codex hook 사용자 Trust가 없어 Codex를 쓰지 않았고 Trust를 승인하지 않았다. 테스트 레벨 lite, 하위 에이전트 standard이며 하위 에이전트는 쓰지 않았다.

## 단계

1. 지시서·D04·DEV 실행 기록·`frontend-design` 스킬·FULLOPS.md·문서 작성 규칙·공통 테스트 규칙을 읽었다. 제품 바이트가 후보와 같음을 `git diff`로 확인했다.
2. 후보를 `git archive`로 scratchpad에 풀고 `TestMemberUIFixtures`로 합성 HTML 29개를 만들었다.
3. `capture.py`로 29개 화면 × 3개 폭을 열어 DOM을 측정하고 PNG 149개(포커스 4개 포함)와 manifest를 저장했다.
4. `keyboard.py`로 Tab·Enter·Space 입력과 렌더 색 대비를 측정했다.
5. 홈·연결 관리·기기 승인·receipt·오류·로그인·확인 코드·gate·원격 동의 PNG를 직접 열어 판정했다.
6. [시각 수락 보고](../../evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md)에 발견 11건과 D04 대조를 기록했다. D04에 수락 상태 한 절(14)을 추가했다.

## 결정

- 제품 시각 수락: PASS. critical·high 0건. medium 1건(V-01 글자 링크 21px)·low 10건은 비차단 잔여다.
- V-01은 좁은 DEV 수정을 권고한다. 수용 여부는 coor가 정한다.
- 디자인 lint·테마 전환은 `project.md`가 미구성으로 정해 해당 없음이다.

## 제한과 재개 조건

- 합성 fixture·정적 `file://`·Chromium 한 엔진이다. 실제 Google 왕복·실계정·운영 CSP 헤더·200% 확대·스크린 리더는 미검증이다.
- 도구 한계로 첫 두 캡처 방식(뷰포트 설정 순서, 전체 화면 캡처)의 결과를 폐기하고 스크롤 조각 방식으로 다시 찍었다. 상세는 보고의 "도구 한계"에 있다.
- 재개 조건: DEV가 V-01을 고치면 링크 11종만 다시 측정한다. 담당은 designer(재측정)와 DEV(수정)이며 전달은 coor가 한다.
