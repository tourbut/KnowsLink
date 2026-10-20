---
title: SAR-UX-ISSUES-002-LANG-TESTER — 언어·설정 독립 QA 시나리오
status: draft
updated: 2026-10-20
owner: tester
tasks: [SAR-UX-ISSUES-002-LANG-TESTER]
summary: "격리 lang fixture에서 U3 언어 저장·선호 우선순위, U7 헤더 '설정', 지정 화면 두 언어, 상태 불변, 3폭·키보드를 HTTP와 실제 브라우저 키 입력으로 확인하는 시나리오와 통과 조건"
---

# SAR-UX-ISSUES-002-LANG-TESTER — 독립 QA 시나리오

대상은 고정 SHA `7704e6d3128831468bb474330210bc2a5099bce7`(fullops/dev, 제품 `2ff9e44`)다. 결과는 [QA 보고서](../qa-reports/SAR-UX-ISSUES-002-LANG-TESTER-test/report.md)에 있다. Jev 웹 조작 도구는 쓰지 않는다. 이유: 비가역 저장 순서 통제, 실제 키 입력(CDP `Input.dispatchKeyEvent`), 쿠키 차단 프로필, 비모바일 viewport 측정이 필요하다.

## 공통 조건

- origin `http://127.0.0.1:18933`만 쓴다(coor SSH 포워딩, 합성 DB). 운영 서비스·실계정·실메일·외부 발송은 쓰지 않는다. manifest의 origin 필드(`39065`)는 서버 쪽 포트이며 tester는 포워딩 포트를 쓴다.
- 허용 GET: `/`, `/auth/verify`, `/home`, `/home/name`, `/home/agents`, `/home/friends`, `/home/account`, `/home/legacy`, `/home/new`, `/home/receipts`, `/home/conversation`, `/home/connections/*`, `/connect/*`, `/home/gates/*`, `/remote/consent/*`, `/lang/*`. 허용 POST: `/home/language`, `/auth/logout`, `/fixture/as/{en|ko|none}`. 친구 요청·결정·이름·에이전트 POST는 쓰지 않는다(designer와 공유).
- 합성 계정 `en`(English 저장, 모든 소유 화면의 주인), `ko`(한국어 저장, 대화 수신자), `none`(저장값 없음). `none`은 저장값을 지울 수 없어 저장 POST를 하지 않는다.
- 쿠키(`__Host-` Secure)는 jar 없이 직접 싣는다. 저장 POST는 마지막에 하고 `en`·`ko`를 원래 값으로 되돌린다.
- 관측 수준: HTTP 응답 본문·헤더(`checks.py`, `mutate.py`, `rate429.py`), Chrome 155·Edge 155 CDP의 실제 키 입력·레이아웃·접근성 트리·캡처(`browser.mjs`), 쿠키 차단 Chrome 프로필. 스크린리더 음성과 Chrome 설정의 실제 200% Page zoom은 범위 밖(not_run)이며 등가 reflow만 확인한다.

## covers와 독립 checks

| covers | 독립 checks (파일) | 통과 조건 |
|---|---|---|
| U7 헤더 '설정' | `checks.py` settings_*_h1_and_header_link, `browser.mjs keys` | 헤더·제목·title이 '설정'/'Settings', '내 계정' 없음, `aria-current=page`, 키보드로 열림 |
| U3 설정 순서·폼 | `checks.py` settings_* | 내 이름 → 내 친구 코드 → 언어 → 로그인과 보안, radio '한국어'/'English', 저장 버튼, 선택만으로 제출 없음(`browser.mjs keys`) |
| U3 저장 성공 | `mutate.py` save_*, `browser.mjs save` | 303 `/home/account?n=language#language`, 선택 언어로 같은 위치, 상태 안내, 새로고침·재로그인·다른 화면·새 브라우저 유지, 다른 열린 세션은 다음 요청부터 |
| U3 거절 | `mutate.py` reject_* | 잘못된 값 422·alert·이전 선택·저장값 유지 |
| 선호 우선순위·손상 | `checks.py` priority_*, anon_* | 저장값 > 비로그인 쿠키 > `?lang=` 한 응답 > 한국어, 손상 쿠키는 한국어 |
| 계정 교대·첫 로그인 | `checks.py` switch_*, first_login_* | 이전 계정 선호 미복사, 비로그인 선호는 화면만 적용하고 저장 안 됨 |
| 비로그인 진입·외부 back | `checks.py` langlink_*, langnav_*, `browser.mjs langflow` | 쿠키 속성, 내부 경로만 복귀, 키보드로 전환 |
| 쿠키 차단 | `browser.mjs nocookie` | 현재 응답만 적용·안내, 저장 안 됨 |
| 익명 저장 차단 | `checks.py` anon_post_language_* | 303 `/?n=expired`, 쿠키 미설정 |
| 지정 화면 두 언어 | `checks.py` screen_*, `browser.mjs layout` ax_names_* | HTML lang = Content-Language, 영어판 한글 잔존 없음(이름·언어명 제외), 접근성 이름 번역 |
| 원문·ID·기한 불변 | `checks.py` invariant_*, `mutate.py` invariant_owner_* | KST 문자열·ID·href·서명 원문 동일 |
| 상태 불변 | `checks.py` state_counts_*, `mutate.py` state_counts_* | 언어 동작 전후 친구·에이전트·연결·메시지·승인 폼 수 동일 |
| 키보드·focus·44px·3폭 | `browser.mjs` layout·keys·zoom | 넘침 0, 잘림 0, 44px, focus 보임, 함정 없음, 720px·320px 등가 reflow |

## 시나리오 밖 (인계)

- 저장 실패(503)는 DB 장애 주입 수단이 없어 열 수 없다(not_run).
- remote 복귀 화면은 fixture가 허용하지 않아 열 수 없다(not_run, DEV 렌더 근거).
- 실제 확인 메일·Google 인증이 필요한 `/auth/verify` 입력 화면은 열 수 없다. 만료 응답(410)만 확인한다.
- 시각 수락은 designer, 고정 SHA 코드 리뷰와 통합은 coor가 조정한다.
