---
title: 언어·설정(U3·U7) 독립 QA
status: review
updated: 2026-10-20
owner: tester
tasks: [SAR-UX-ISSUES-002-LANG-TESTER]
summary: "7704e6d(제품 2ff9e44) 격리 fixture 독립 QA: U3 언어 저장·우선순위·두 언어 지정 화면·상태 불변·3폭·키보드 PASS. critical·high 없음, medium 1건(회원 요청 한도 429 화면이 항상 한국어), low 관찰 4건"
---

# SAR-UX-ISSUES-002-LANG-TESTER 독립 QA

## 대상과 결론

대상은 fullops/dev `7704e6d3128831468bb474330210bc2a5099bce7`(제품 `2ff9e44`, 이후 테스트·보고만)다. 기준 main `0b29416`(지시서 base `08f2698`). fixture `http://127.0.0.1:18933`(coor SSH 포워딩, manifest `lang-ux002-dev/fixture-manifest.json`; manifest의 `origin` 필드 `39065`는 서버 쪽 포트). 서버가 이 SHA인지는 coor 지시와 manifest를 따랐고 tester가 독립 확인하지는 않았다. 제품 코드는 고치지 않았다.

**결론: 지시 check 모두 PASS. critical·high 결함 없음. medium 1건 D-LG-01, low 관찰 4건(O-LG-01~04).** 스크린리더·Chrome 설정의 실제 200% Page zoom·저장 실패(503)·remote 복귀 화면은 not_run이다. 전체 접근성 완료로 보고하지 않는다.

| 영역 | 실행 | 결과 |
|---|---|---|
| HTTP 읽기 전용 [checks.py](checks.py) → [http-result.json](http-result.json) | 232 | 232 PASS |
| 저장 HTTP(마지막) [mutate.py](mutate.py) → [mutate-result.json](mutate-result.json) | 64 | 64 PASS |
| Chrome 155.0.8059.85 layout en / ko ([en](browser-chrome-layout-en.json) · [ko](browser-chrome-layout-ko.json)) | 128 / 72 | 모두 PASS |
| Edge 155.0.4283.62 layout en / ko | 128 / 72 | 모두 PASS |
| Chrome·Edge keys en·ko(설정 Tab 순회·radio 방향키·Space·건너뛰기) | 15 × 4 | 모두 PASS |
| Chrome·Edge langflow(공개 페이지 키보드 전환·쿠키·remote 복귀·`?lang=`) | 17 × 2 | 모두 PASS |
| Chrome 쿠키 차단 프로필 [nocookie](browser-chrome-nocookie-nocookie.json) | 10 | 모두 PASS |
| Chrome reflow 등가 200%·400% [zoom](browser-chrome-zoom.json) | 16 | 모두 PASS |
| Chrome 실제 키 입력 저장 [save](browser-chrome-save.json)(en→ko→en, ko→en→ko, 소유 화면 한국어 3폭) | 73 | 모두 PASS |
| 한도 429 화면 언어 [rate429.py](rate429.py) | 4 계정·쿠키 조합 | **D-LG-01** |

실행 명령(종료코드 보존, 파이프 없음): `python checks.py`, `bash run_readonly.sh`(Chrome·Edge layout/keys/langflow/zoom), `python mutate.py`, `node browser.mjs 9333 chrome save`, `node browser.mjs 9335 chrome-nocookie nocookie`, `python rate429.py`. 브라우저는 `--headless=new --disable-extensions`, 비모바일 viewport 1440×900·390×844·320×640, 실제 키 입력은 `Input.dispatchKeyEvent`다. 시나리오는 [scenarios/SAR-UX-ISSUES-002-LANG-TESTER.md](../../scenarios/SAR-UX-ISSUES-002-LANG-TESTER.md)다.

## 결함

### D-LG-01 (medium) 회원 요청 한도 429 화면이 항상 한국어다

- 재현: `python rate429.py`(합성 `en` 계정으로 `GET /home/account`를 분당 한도까지 반복, 약 40번째 응답). 계정별로 저장값·쿠키를 바꿔 반복했다.
- 기대: DEV 보고는 ‘admission 429/503 화면은 회원 저장값이 아니라 브라우저 쿠키 언어를 따른다’고 적었다. 따라서 `__Host-kl_lang=en` 쿠키가 있으면 영어, 설계 U3은 거부·오류 화면도 번역 범위로 둔다.
- 실제: `en` 저장·쿠키 없음, `en` 저장·쿠키 `ko`, `none`·쿠키 `en`, `ko` 저장·쿠키 `en` 모두 `<html lang="ko">`, 제목 ‘요청 제한 · KnowsLink’, 본문 ‘요청이 많습니다. … 이후 다시 시도하세요.’다. 쿠키 `en`도 무시된다. `Content-Language` 헤더가 없다. 원본: [rate429-result.json](rate429-result.json), [rate429-result-2.json](rate429-result-2.json).
- 영향: 영어를 고른 사용자가 한도에 걸리면 한국어 오류만 본다. 데이터·권한 영향 없음. 번역 범위의 오류 화면 누락이며 DEV 문서의 알려진 한계 설명과 실제 동작이 다르다.
- SHA: `7704e6d`. 수정 여부는 dev/coor 판단이다. 이 QA는 수정하지 않았다.

## 확인한 내용 (지시서 항목별)

**설정 화면 (U3·U7).** `en`·`ko`·`none` 모두 title ‘설정 · KnowsLink’/‘Settings · KnowsLink’, h1, 헤더 링크(`aria-current=page`)가 ‘설정’/‘Settings’이고 ‘내 계정’은 어디에도 없다. 섹션 순서는 내 이름 → 내 친구 코드 → 언어 → 로그인과 보안이다. 언어 섹션은 `fieldset role=radiogroup`·legend ‘표시 언어’/‘Display language’, radio ‘한국어’(`lang=ko`)·‘English’(`lang=en`), 현재 저장값이 checked이고 `none`은 한국어다. 저장 버튼은 ‘언어 저장’/‘Save language’다(O-LG-02). 방향키·Space로 선택해도 요청이 없고 새로고침하면 저장값으로 돌아온다. 선택 뒤 Tab으로 저장 버튼에 가서 Enter를 눌러야만 POST가 한 번 나간다.

**저장 성공·거부 (U3).** `en`→ko 저장은 303 `/home/account?n=language#language`, 응답이 언어 쿠키를 쓰지 않고, 도착 화면이 선택 언어·같은 위치·상태 안내(‘언어를 저장했습니다. 다음 화면부터 이 언어로 표시합니다.’ / ‘Language saved. Screens from now on are shown in this language.’)다. 열려 있던 다른 세션은 강제 새로고침되지 않고 다음 요청부터 바뀐다. 새로고침·새 로그인 세션·다른 화면·새 브라우저에서 유지되고 다른 계정과 익명은 영향이 없다. 같은 값 재저장도 정상이다. 잘못된 값(`xx`·빈 값·누락·`EN`·NUL)은 422·`role=alert`·이전 선택 유지·저장값 불변이고 영어 화면은 영어 안내다. 저장 뒤 `en`·`ko`는 원래 값으로 되돌렸다(마지막에 세 계정을 다시 읽어 확인).

**선호 우선순위·손상.** 저장값 > 비로그인 쿠키 > `?lang=` 한 응답 > 한국어를 회원 17조합과 비로그인 4화면×5쿠키 상태로 확인했다. 손상·빈·대문자 쿠키는 한국어다. 쿠키 없는 `?lang=en`은 그 응답만 영어이고 ‘이 브라우저에 언어를 저장하지 못했습니다…’ 안내(`role=status`)가 붙으며 링크·폼으로 전파되지 않는다.

**계정 교대·첫 로그인.** 같은 브라우저에서 en→ko는 한국어, en→none·ko→none은 이전 선호를 복사하지 않는다. `none`이 비로그인 `en` 쿠키로 로그인하면 화면만 영어이고 radio도 English지만, 쿠키 없는 새 브라우저의 `none`은 한국어다(저장되지 않음). 저장값 있는 계정은 반대 쿠키를 이긴다. 로그아웃 뒤 공개 페이지는 저장값을 이어받지 않는다.

**비로그인 진입·외부 back·익명 저장.** `/lang/{ko|en}`은 303이며 쿠키가 `__Host-`·Secure·HttpOnly·SameSite=Strict·Path=/·365일이다. `back`에 외부 URL·`//`·`/\`·`javascript:`·절대 URL·상대 경로를 넣어도 내부 경로로만 돌아간다(`/%2f%2fevil.example`은 같은 호스트의 리터럴 경로). 알 수 없는 코드·POST는 거부된다. 공개 첫페이지·확인 코드(`/auth/verify`)·device 연결·remote 동의의 ‘한국어 / English’ 링크는 현재 언어에 `aria-current`가 있고 키보드 Tab 정지점이며 높이 44px 이상이다. 익명 `POST /home/language`는 값과 무관하게 303 `/?n=expired`이고 언어 쿠키를 쓰지 않으며 교차 출처 POST는 회원도 403이다. 쿠키를 막은 Chrome 프로필에서 세 화면 모두 영어 링크가 현재 응답만 영어로 바꾸고 안내를 보이며 다음 방문은 한국어다(저장 안 됨).

**지정 화면 두 언어.** 공개·로그인(410)·이름·대화·친구·내 에이전트·설정·legacy·조회·승인·local 연결 3종·device 연결·remote 동의를 en/ko로 열었다. 모두 `<html lang>`=`Content-Language`, 고유 title(홈과 대화는 설계상 같은 ‘대화’), h1 1개, 건너뛰기 링크·main이다. 영어판 HTML 전체(속성·aria-label 포함)에서 한글은 사용자 이름 데이터와 언어명 ‘한국어’뿐이다. 접근성 트리의 영어 화면 이름에도 한글이 없다. 한국어판에는 사용자 데이터·기술 용어 외 영어 문구가 없다. 소유 화면이 `en` 계정에 있어 한국어판은 `en`을 잠시 한국어로 저장해 3폭으로 캡처했다. 403 거부 화면(비소유자)도 두 언어다. 만료 코드 410 안내는 영어로 번역된다.

**원문·ID·기한 불변.** 한국어 저장 중과 영어 저장 중의 소유 화면(승인·대화·legacy·에이전트·연결 3종·친구) 비교에서 href·form action·hidden 값(세션별 csrf·attempt 제외)·ID·KST 기한 문자열·서명 원문 JSON이 같다(명령 안 `<자리표시자>` 문구만 언어별). device·remote 화면 KST·지문·URL도 같다. 정확한 기한은 ‘2026-10-21 07:49:06 KST’ 그대로다.

**상태 불변.** 언어 링크·`?lang=` 반복·저장 전후 친구 3·에이전트 8·legacy 4·메시지 1·승인 폼 수와 모든 KST 문자열이 같다. 언어 동작이 친구 요청·발송·승인을 실행하지 않는다.

**키보드·focus·44px·3폭.** 설정 화면 Tab 순서는 건너뛰기 → 브랜드 → 설정 → 주 메뉴 3 → 서비스 안내 → 이름 → 저장 → 친구 코드(교체) → 언어 radio(한 정지점) → 언어 저장 → 보안 버튼 → 기존 연결이며 한 바퀴 후 함정이 없다. 모든 정지점에 보이는 focus 2px 이상, 높이 44px 이상(건너뛰기 제외)이다. 1440·390·320에서 두 언어 모두 가로 넘침 0, 잘림·ellipsis 0, 44px 미만 0이다(긴 영문·긴 한글 이름 포함). Chrome 설정의 실제 200%가 아니라 720px(DPR2)·320px(DPR4) reflow로 등가 확인했다.

## 관찰 (결함 아님, 모두 low)

| ID | 내용 | 재현 | 비고 |
|---|---|---|---|
| O-LG-01 | 영어를 저장한 로그인 회원이 `/auth/verify`(만료 코드 410)를 열면 한국어다. 이 경로는 쿠키만 본다. | `en` 세션으로 `GET /auth/verify` → `Content-Language: ko` | 비로그인 화면이라 회원을 조회하지 않는 설계로 보인다. 설계 문서는 로그인 전 화면의 언어를 쿠키로 둔다. designer 판단 |
| O-LG-02 | 저장 버튼 문구가 단일 언어(‘언어 저장’ 또는 ‘Save language’)다. 설계 U3 본문은 ‘언어 저장 / Save language’ 병기를 적고 와이어프레임은 ‘[언어 저장]’이다. | 설정 화면 | 설계 두 곳이 어긋나 해석 차이다. 영어 화면에서 한국어를 읽지 못하는 사용자에게는 현재 구현이 낫다. designer 확인 |
| O-LG-03 | `details` 요약의 접근성 이름에 장식 문자가 붙는다(‘▸ Friend code (replace)’). CSS `::before` 내용이 이름에 들어간다. | 설정 키 정지점의 AX 이름 | 이번 변경 이전부터 있던 공용 스타일이다. 두 언어 동일 |
| O-LG-04 | 한국어 legacy 화면에 ‘wake’·‘idle’·‘TTL’·‘bytes’·‘text’ 같은 영어 기술 용어가 남는다. | 한국어 `/home/legacy` 설명문 | 명령·프로토콜 용어로 원문이다. 번역 누락은 아니다 |

시각 판정(390px 공개 첫페이지의 언어 링크가 다음 줄로 내려감 등)은 tester 범위 밖이며 designer 몫이다.

## 실패·보정 기록 (원 결과 보존)

- `checks.py` 첫 실행 [227 PASS / 8 FAIL](http-result-run1-expectation-errors.json)([로그](checks-run1.log)): 모두 **내 기대값 오류**다. ko가 대화 수신자라 403이 아니라 200, 홈·대화가 같은 title, 이름 화면에 aria-label 없음, 깨진 쿠키 예(`en;ko`)가 쿠키 두 개로 파싱됨, `/lang`의 복귀 경로가 `?lang=` 힌트를 붙이는 설계, `/%2f%2fevil.example`이 같은 호스트 경로, 세션 만료 시각이 로그인마다 다름. 기대를 고쳐 232/232 PASS.
- Chrome layout 첫 실행 1·2건 FAIL: 접근성 트리의 `InlineTextBox` 조각 ‘름’(긴 한국어 사용자 이름의 일부)을 한글로 오인한 **검사 결함**. 조각 노드를 제외해 PASS. [첫 결과](browser-chrome-layout-en-first-run-ax-textbox-fragment.json).
- Chrome keys 첫 실행 각 1건 FAIL: 요약 AX 이름의 `▸ ` 접두를 몰라 기대 순서가 어긋난 **내 기대값 오류**(O-LG-03 관찰의 원인). [en](browser-chrome-keys-en-first-run-expectation-error.json) · [ko](browser-chrome-keys-ko-first-run-expectation-error.json).
- Chrome langflow 첫 실행 8건 FAIL: 두 번째 Tab 순회가 마지막 focus 지점에서 이어져 영어 링크에 못 간 **스크립트 오류**. 영어 링크에 직접 focus하도록 바꿔 PASS. [첫 결과](browser-chrome-langflow-first-run-script-focus-error.json) · [로그](chrome-langflow-first-run.log).
- `mutate.py` 1·2차 실행: 저장 뒤 landing에서 상태 안내가 없다고 나온 것은 내 http 클라이언트가 URL 조각(`#language`)을 요청 경로에 실어 보낸 **검사 결함**이다(브라우저는 조각을 보내지 않는다). 조각을 제거한 뒤 안내가 나온다. 그 밖에 세션별 csrf·attempt 토큰, 세션 만료 시각, 명령 자리표시자, 영어 기술 용어 허용 목록을 기대값 오류로 보정했다. [1차](mutate-result-run1-client-fragment-and-expectation-errors.json) · [2차](mutate-result-run2-session-expiry-expectation-error.json). 3차 64/64 PASS. 두 번의 중간 실행도 같은 복구 절차로 `en`·`ko`를 원래 값으로 되돌렸다.
- 셸 가드가 `rm`을 막아 한 번 명령이 취소됐다. 삭제는 하지 않았고 스크립트를 다시 만들어 실행했다.

## fixture 상태 변경 알림 (designer)

- `en`·`ko` 계정의 저장 언어를 `mutate.py`·`browser.mjs save`에서 여러 번 바꿨다가 **모두 원래 값(`en`=English, `ko`=한국어)으로 되돌렸다**. 마지막 확인에서 `en`→English, `ko`→한국어, `none`→저장값 없음(한국어)이다. 바뀐 동안(합산 약 수 분) `en`의 소유 화면이 한국어로 보였을 수 있다. `none`은 저장하지 않았다.
- `rate429.py`는 회원 요청 한도를 일부러 소진했다(`ko`·`en`·`none` 각 약 1분). 그 사이 designer가 같은 계정으로 GET하면 429 화면이 보였을 수 있다. 친구 요청·에이전트·이름·승인 상태는 바꾸지 않았다.

## designer용 캡처

Chrome, 파일명 `chrome-<세션>-<언어>-<화면>-<폭>.png`(이 폴더, 230장). 세션 `en`/`ko`는 해당 합성 계정, `anon`은 비로그인(쿠키 언어가 이름에 있음). 폭은 1440·390·320이고 `-open`은 details를 펼친 상태다.
- 설정: `chrome-en-en-account-*`, `chrome-ko-ko-account-*`, 저장 직후 `chrome-en-ko-account-after-save-*`·`chrome-ko-en-account-after-save-*`(상태 안내 포함), reflow `-200pct`·`-400pct`.
- 지정 화면 영어: `chrome-en-en-<home|name|conversation|friends|agents|legacy|new|receipts|conn-*|gate|connect-device|remote>-*`, 한국어(소유 화면 포함): `chrome-ko-ko-*`와 `chrome-en-ko-<conversation|gate|conn-*|legacy|agents|friends|home>-*-ownerko.png`.
- 비로그인: `chrome-anon-<en|ko>-anon-<public|verify|connect-device|remote>-*`, 쿠키 차단 안내 `chrome-anon-en-public-query-notice-1440.png`·`chrome-anon-en-*-cookies-blocked-1440.png`(nocookie 프로필 캡처는 `chrome-nocookie-…` 이름 접두를 쓴다).

## not_run

- 스크린리더 낭독(NVDA·VoiceOver 등), Chrome 설정 Page zoom의 실제 200%(reflow 등가만 확인), Firefox·Safari·모바일 실기.
- 저장 실패(503) 경로: DB 장애 주입 수단이 없다. DEV 코드 경로만 있고 자동 검사도 없다(DEV 보고와 동일).
- remote 복귀 화면(`/remote/…` 복귀): fixture가 열지 않는다. DEV 렌더 근거에 의존한다.
- `/auth/verify`의 실제 코드 입력 화면·확인 메일·Google 화면: 합성 메일 없음. 만료 410 응답만 확인했다.
- 공용 admission 503 화면과 익명 429 화면의 언어: 429(회원 한도)만 확인했다(D-LG-01).
- `make lint`·`make test`·integration: coor의 같은 SHA Linux 결과를 재사용했고 다시 실행하지 않았다. 영어 문장의 전체 문체 교정은 이 QA 범위 밖(표본 읽기만).

## 파일

`checks.py`·`checks_lib.py`·`client.py`·`dump.py`·`probe.py`(HTTP), `mutate.py`(저장), `rate429.py`, `browser.mjs`(CDP), `run_readonly.sh`, `pages/`(두 언어 화면 HTML 사본), 결과 JSON·로그·캡처. 브라우저 프로필은 `.browser-profiles/`에 두고 `.gitignore`로 커밋에서 제외했다. 서버는 건드리지 않았고 요청은 fixture origin뿐이다.
