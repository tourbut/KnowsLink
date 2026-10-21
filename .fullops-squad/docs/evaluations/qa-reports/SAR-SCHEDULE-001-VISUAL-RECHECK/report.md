---
title: SAR-SCHEDULE-001-VISUAL-RECHECK 직접 시각 재검수
status: review
updated: 2026-10-21
owner: designer
tasks: [SAR-SCHEDULE-001-VISUAL-RECHECK]
summary: 일정 조회 화면 수정 V1·V2·V3·D2의 main 336533e 반영을 KO/EN 3폭으로 다시 보고 시각 수락한다
---

# SAR-SCHEDULE-001-VISUAL-RECHECK 직접 시각 재검수

## 대상과 결론

- 대상 제품 코드: main `336533e6e15679d67118cbec5ab907768b80a70f`. 수정 항목 정의는 [SAR-SCHEDULE-001-VISUAL 보고](../SAR-SCHEDULE-001-VISUAL/report.md)다.
- 결론: **수락**. V1·V2·V3·D2 모두 PASS다. 이전 검수 PASS 화면에 회귀가 없다.
- 이 수락은 일정 조회 화면의 시각 수락이다. 전체 제품 수락이 아니다. 독립 동작 QA는 tester의 SAR-SCHEDULE-001-QA다.

## 방법

1. `git archive 336533e`를 scratchpad에 풀었다. 워크트리 제품 코드는 바꾸지 않았다.
2. `KNOWSLINK_UI_FIXTURE_DIR=<scratch>/fx go test ./internal/relay -run 'TestHumanAccountCalendarStates|TestMemberUIFixtures|TestEnglishScreensHaveNoKorean|TestZZVisualScratch' -count=1`을 실행했다. 결과는 `ok`다.
3. 이전 검수의 scratch 테스트 `TestZZVisualScratch`를 다시 썼다. 해제 뒤 안내 문구만 `human_http.go:105`의 현재 문구로 바꿨다. 영어판은 실제 `i18n_en_human.go` 번역을 거친다. 이 파일은 커밋하지 않는다(`raw/tools/`).
4. headless Chrome(비모바일 viewport)으로 1440·390·320을 렌더했다. 각 화면을 details 닫힘과 모두 열림 두 상태로 쟀다. 측정은 가로 넘침(`scrollWidth ≤ innerWidth`)과 summary 분리(자식 요소가 없는 summary의 문구 시작이 summary 왼쪽 20px 안)다. 도구는 `raw/tools/recheck.mjs`다.
5. 판정 캡처는 KO/EN 행, 3폭 열의 비교판으로 묶어 직접 열어 보았다(`raw/tools/compose.mjs`).
6. 문구는 렌더된 HTML에서 정본 문장을 고정 문자열로 찾았다.

## 수정 항목 판정

| 항목 | 화면 | 판정 | 근거 |
|---|---|---|---|
| V1 | gate 조회 unconnected KO/EN, 3폭 | PASS | [비교판](v1-gate-query-unconnected.png). '승인하면 허용되는 범위' 아래에 정본 한 줄(KO '이 요청은 승인할 수 없습니다. 거절하거나 기한이 지나면 조회 없이 끝납니다.', EN 'This request cannot be approved. If you decline it or it expires, it ends without a lookup.')만 있다. 이전 두 줄('정보 공개나 일정 실행을 허용하지 않습니다', '가능 시간 공개 정책이 없어')은 0건이다. '친구 동의나 대화 메시지는 …' 줄, 오류 alert, '캘린더 연결' 링크, 거절 버튼이 남아 있다. connected 화면에는 새 줄이 없다 |
| V2 | gate 조회 connected·unconnected·approved KO/EN, 3폭 | PASS | [비교판](v2-gate-query-summary.png). EN 320의 'Reference note — not grounds for approval'은 기호와 같은 줄에서 시작하고 넘친 'approval'만 다음 줄로 간다. 측정 summary 분리 0 |
| V3 | 설정 캘린더 reconnect KO/EN, 3폭 | PASS | [비교판](v3-calendar-reconnect.png). '캘린더 다시 연결' 주 버튼 아래에 위험 details '캘린더 연결 해제'가 있다. 열면 D2 해제 안내 뒤에 해제 버튼이 있다. 해제 form은 details 밖에 0개다 |
| D2 해제 안내 | 설정 캘린더 connected·reconnect(details 열림) KO/EN, 3폭 | PASS | [비교판](d2-calendar-connected-open.png). KO/EN 문장이 정본과 같다. 이전 문구 '대기 중인 조회 요청은 거절로'는 모든 fixture에서 0건이다 |
| D2 해제 뒤 안내 | 설정 상단 status KO/EN, 3폭 | PASS | [비교판](d2-off-notice.png). KO '…아직 전달되지 않은 조회 요청과 결과는 거절로 끝났습니다.', EN 'Calendar disconnected. The token KnowsLink stored was discarded, and availability requests and results not yet delivered ended as declined.'가 정본과 같다 |

D2 서버 동작(미전달 결과 삭제)은 이 시각 검수의 범위가 아니다. DEV 테스트 `disconnect_drops_undelivered_result`·`disconnect_keeps_consumed_result`가 있고 tester QA가 독립 확인한다.

## 회귀 확인

| 확인 | 범위 | 결과 |
|---|---|---|
| 이전 PASS 화면 | 설정 캘린더 none-google·email-member·not-ready, gate 조회 connected·approved. KO/EN, 3폭, details 닫힘·열림 | PASS. 배치·문구가 이전 검수와 같다. 넘침 0, summary 분리 0 |
| V2 공용 CSS 영향 | `summary:not(:has(*))`는 모든 화면의 text-only summary에 적용된다. summary가 있는 fixture 73화면 × 3폭, details 모두 열림 | 219개 모두 넘침 0, summary 분리 0 |
| 영어판 한국어 잔존 | 영어판 캘린더 6화면 | 언어 선택지의 `lang="ko"` '한국어'만 있다. 의도된 표기다. `TestEnglishScreensHaveNoKorean` 통과 |
| 대비 | `18d74dd..336533e`의 CSS 변경은 summary 배치 한 줄이다. 색 토큰은 바뀌지 않았다 | 이전 검수 값(최저 버튼 테두리 3.79:1, 본문 텍스트 6.19:1 이상)을 그대로 쓴다. PASS |
| 피할 항목 | 크림 배경·알약 버튼·모노스페이스 라벨·장식 그라디언트·이모지 아이콘 | 없다. 모노스페이스는 기술 상세의 ID·정책에만 쓴다 |

## 관찰(수정 요구 아님)

- `:has()`를 지원하지 않는 브라우저는 V2 이전 배치(기호 단독 줄)로 보인다. 기능 손실이 없어 수락한다(DEV 보고와 같다).
- unconnected 조회의 '기술 상세 > 적용 정책'은 정책 정의 문자열(`schedule.query: approve runs one free/busy lookup…`)을 그대로 보인다. 개발자용 정책 ID 표시이고 이전 검수에서도 PASS였다. 바꾸지 않는다.

## 미실행과 한계

- 정적 template 렌더다. 실제 인증·DB·격리 fixture 왕복은 실행하지 않았다.
- 해제 뒤 안내 KO 문장은 scratch fixture가 `human_http.go:105` 문자열을 그대로 넣어 렌더했다. 실제 `?n=calendar-off` 경로는 `calendar_integration_test.go`가 확인한다.
- 키보드·200% 확대·스크린리더는 not_run이다(tester 몫).
- 실제 Google 계정 1경로는 D011 뒤 coor 몫이다(SC-09).
- 판정에 쓰지 않은 캡처(닫힘·열림 × 3폭)·측정 JSON·scratch 도구는 `raw/`에 있고 커밋하지 않는다.
