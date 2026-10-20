---
title: SAR-UX-ISSUES-002-APPROVAL-TESTER — 승인 요약·승인 화면 접근성·O-NM-04·O-NM-08 독립 QA 시나리오
status: draft
updated: 2026-10-20
owner: tester
tasks: [SAR-UX-ISSUES-002-APPROVAL-TESTER]
summary: "격리 approval fixture에서 U5 승인 요약, 승인·거절·중복 클릭, U7 승인 화면 접근성, O-NM-04, O-NM-08을 HTTP와 실제 브라우저 키 입력으로 확인하는 시나리오와 통과 조건"
---

# SAR-UX-ISSUES-002-APPROVAL-TESTER — 독립 QA 시나리오

대상은 고정 SHA `9d13a05ba93c7c085034213590ce661305606010`(fullops/dev, 제품 `f43d05c`, fixture 수정 `ea94128`)다. 결과는 [QA 보고서](../qa-reports/SAR-UX-ISSUES-002-APPROVAL-TESTER-test/report.md)에 있다. Jev 웹 조작 도구는 쓰지 않는다. 이유: 비가역 POST 순서 통제, 실제 키 입력(CDP `Input.dispatchKeyEvent`), 비모바일 viewport 측정이 필요하다.

## 공통 조건

- origin `http://127.0.0.1:18931`만 쓴다(coor SSH 포워딩, 합성 DB, manifest `approval-ux002-dev/fixture-manifest.json`). 운영 서비스·실계정·실메일·외부 발송은 쓰지 않는다.
- 허용 GET: `/`, `/home`, `/home/legacy`, `/home/friends`, `/home/account`, `/home/gates/*`. 허용 POST: `/home/gates/*`와 `/fixture/as/{a|b|c}`. 그 밖은 403이 정상이다.
- 합성 A(승인자, 이름 있는 40자 한글 에이전트 `agent_b` 소유), B(요청자 `agent_a` 소유), C(gate 없음, A에게 친구 요청 1건). A와 B는 legacy pair라 요청자는 ‘확인할 수 없는 에이전트 / 구분값’으로 보인다. 이름 있는 요청자 화면은 fixture로 열 수 없어 DEV 렌더 fixture에서만 확인된다(not_run).
- 승인·거절 POST는 `pending_decide_1`·`pending_decide_2`에만, 읽기 전용 검사와 designer 캡처가 끝난 뒤 실행한다. designer가 같은 fixture를 GET으로 함께 쓴다.
- 관측 수준: HTTP 응답 본문·헤더(`checks.py`, `mutate.py`), Chrome 155·Edge 155 headless CDP의 실제 키 입력·레이아웃 측정·접근성 트리·대비·캡처(`browser.mjs`). 스크린리더 음성·실제 브라우저 200%·보조 브라우저 전 조합은 범위 밖(not_run)이다.

## covers와 독립 checks

| covers | 독립 checks (파일) | 통과 조건 |
|---|---|---|
| U5 빈 목록 | `checks.py` list_c_empty_text_exact | C의 승인 요청 절이 설계 문구 그대로이고 행·‘요청 검토’가 없다 |
| U5 목록 행 | `checks.py` list_* | 행마다 요청한·받은 에이전트(허용된 이름·구분값), 요청 유형, 절대 기한, 상태, ‘요청 검토’ 링크(설명하는 aria-label). 원요청 ID는 목록에 없고 gate ID는 href에만 있다. 상태 대기 4·만료·승인·확인할 수 없음. 대화 화면 ‘승인 N건’은 대기만 센다 |
| U5 읽기 순서·내용 | `checks.py` *_reading_order*, *_summary_matches_signed_body, *_fixed_scope_sentence | 누가 누구에게 → 요청한 내용 → 승인하면 허용되는 범위 → 기한·상태 → 행동 → 기술 상세 → 복귀. 요약 값이 서명된 원문 JSON과 KST 환산까지 일치. 고정 설명 문구 |
| U5 hint·자유문·HTML | `checks.py` hint_*, *_free_text_* | hint는 ‘참고 문구 — 승인 근거 아님’ 접힘 영역 1곳에 escape. 자유문은 ‘요청자가 적은 내용’에만 있고 링크·태그 없음, 범위 절에 오르지 않음. 페이지에 script 없음 |
| U5 닫힌 상태 | `checks.py` expired_*, processed_*, original_missing_*; `browser.mjs keys` | 버튼·form·CSRF 없음, status 안내, 원문 부재는 ‘승인할 내용을 확인할 수 없습니다’ alert와 복귀만 |
| U5 권한·재검사 | `checks.py` other_member_*, anonymous_*; `mutate.py 1` | 타 회원 GET·POST 403에 내용 비노출, 없는 gate도 같은 응답, 익명 303. 교차 gate token·타 세션 token·잘못된·빈·없는 CSRF·알 수 없는 decision 403/4xx에 상태 불변. GET 쿼리로 결정되지 않음 |
| U5 승인·거절·중복 클릭 | `mutate.py 1·2`, `browser.mjs flow` | 승인 303 → 상세 승인·버튼 없음·목록 승인·대기 수 감소. 같은 세션 중복 승인·거절 409와 ‘이번 제출은 적용되지 않았습니다’·복귀, 첫 결정 유지. 실제 키 입력 거절 1회 POST, 뒤로가기는 서버 최신 상태(폼 없음) |
| U7 승인 화면 | `checks.py` *_status200_title*, `browser.mjs keys` | 고유 title·h1 1·main·skip·위치 nav(`aria-current=page`는 요청 승인만)·승인 목록 복귀(nav·하단 링크), Tab 순서, 보이는 focus, 44px, Space·Enter로 details, Shift+Tab, 목록 ↔ 상세 왕복 |
| U7 레이아웃 | `browser.mjs layout·after` | 1440·390·320 가로 넘침·잘림·말줄임 0, 44px 미만 0. 640px 폭 reflow(200% 등가) |
| U7 대비 | `browser.mjs contrast` | 본문 텍스트 4.5:1(큰 글자 3:1) |
| O-NM-04 | `browser.mjs onm04` | 익명 429 화면에 선택기·폼 없음, 재시도 시각, 실제 로그인 링크(`/`), 빈 href 없음, 복원 문구 없음. fixture가 `/home/receipts`를 막아 고정 SHA 템플릿을 렌더해 연다(실제 429 아님) |
| O-NM-08 | `checks.py` onm08_*, `browser.mjs keys` | 수락 → 안내 → 거절 순서, 안내가 Tab 정지점이 아님, 거절의 접근성 설명 = 안내, 위 16px |

## 시나리오 밖 (인계)

- 이름 있는 요청자 gate, 실제 익명 429 응답, 이름 변경 뒤 같은 원요청 식별: fixture가 막거나 데이터가 없어 DEV Go 테스트·렌더 fixture 근거이며 not_run으로 남긴다.
- 스크린리더, 실제 브라우저 200%, Firefox·Safari·모바일 실기. 시각 수락은 designer, 고정 SHA 코드 리뷰와 통합은 coor가 조정한다.
