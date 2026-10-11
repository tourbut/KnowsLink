---
title: SAR-SERVICE-UI-001-FIX-VISUAL 실행 기록
status: review
updated: 2026-10-11
owner: designer
tasks: [SAR-SERVICE-UI-001-FIX-VISUAL]
summary: 고정 후보 295417a의 V-01 글자 링크 44px 수정을 Orca 내장 브라우저에서 직접 재측정한 실행 단계·결정·제한 기록
---

# SAR-SERVICE-UI-001-FIX-VISUAL 실행 기록

## 목적과 범위

DEV가 고친 V-01(글자 링크 21px, D04 9절 44px 미달)을 고정 후보 `295417acfdcd415ad0c768d164bd71dc07214ea9`에서 직접 재측정한다. 제품 코드는 바꾸지 않았다. 실행자는 등록된 Claude Sonnet 5.5 새 세션(`059f2e8a-d31f-442f-a5af-ce1ac3255fe7`)이다. 테스트 레벨 lite, 하위 에이전트 standard이며 작은 직렬 검수라 하위 에이전트를 쓰지 않았다. Codex Trust 미승인으로 Codex를 쓰지 않았다.

## 단계

1. 지시서·FULLOPS.md·D04·원 시각 보고 V-01·FIX-DEV 기록을 읽었다. 제품 바이트가 후보와 같음을 `git diff`로 확인했다.
2. 후보를 `git archive`로 scratchpad에 풀고 `TestMemberUIFixtures`로 합성 HTML 29개를 새로 만들었다.
3. `measure.py`로 29개 화면 × 3개 폭의 모든 상호작용 요소를 측정하고 대표 PNG를 저장했다.
4. PNG를 열어 연결 취소 버튼 옆 링크·오류 화면·확인 코드·원격 동의·홈의 시각 부작용을 확인했다.
5. [재측정 보고](../../evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL/report.md)를 쓰고 원 보고서·D04 14절에 후속 기록을 덧붙였다.

## 결정

- V-01 resolved(시각 기준). 제품 시각 PASS 유지. 새 critical·high·medium 없음. low N-01(마지막 링크 아래 여백)만 비차단 관찰로 남긴다.
- 원 보고서 발견 V-01–V-12는 삭제하지 않았다. V-02–V-12는 범위 밖이라 그대로다.

## 제한과 재개 조건

- 합성 fixture·`file://`·Chromium 한 엔진이다. 실제 Google 왕복·실계정·운영 렌더·스크린 리더·200%·다크 모드·키보드 포커스 재입력은 미검증이다.
- 첫 실행이 viewport 폭 불일치로 중단됐다. 폭을 다시 설정하고 확인하는 재시도를 넣어 고쳤다. 이 중단 실행의 PNG 2개는 같은 이름으로 다음 실행이 덮어썼다.
- 재개 조건: 운영 배포 뒤 실제 렌더 확인은 OPS·사용자 세션이 필요하다. main 통합·배포는 coor가 정한다.
