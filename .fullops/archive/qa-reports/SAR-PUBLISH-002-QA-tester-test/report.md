---
title: SAR-PUBLISH-002-QA 운영 KO/EN 정책·키보드·200%·언어 전환 QA
status: review
updated: 2026-10-22
owner: tester
tasks: [SAR-PUBLISH-002-QA]
summary: "운영 /·/privacy·/terms의 KO/EN 언어 전환, Tab 순회, 640px 재배치, V-01 회귀를 읽기 전용으로 확인하고 저장 언어를 격리 fixture에서 확인했다. 5개 항목 중 4개 PASS, 운영 Google 로그인 scope 실측은 승인 필요로 NOT RUN, 제품 결함 0"
---

# SAR-PUBLISH-002-QA 운영 KO/EN 정책·키보드·200%·언어 전환 QA

대상 운영 주소: `https://link.knowslog.com/`·`/privacy`·`/terms`. 제품 코드는 main `4fdd149b`다. 이 워크트리 HEAD `24198e40`은 `.fullops-squad/`만 바꾼다(`git diff --stat 4fdd149b HEAD -- . ':!.fullops-squad'` 출력 없음). 테스트 레벨 lite. 실행일 2026-10-22.

운영에는 GET, 링크 이동(`/lang/…`), Tab 키, 스크롤만 보냈다. 폼 제출·로그인·쓰기 요청은 없다. 언어 쿠키만 브라우저 쪽에 생겼고 실행마다 비웠다. 로그인 회원 확인은 원격 Linux 격리 fixture(합성 DB, loopback)에서 했다.

## 결과 요약

| 항목 | 판정 | 근거 한 줄 |
|---|---|---|
| 1. 언어 전환 링크 (`/`, `/privacy`, `/terms` × KO/EN) | PASS | 6개 왕복 모두 같은 경로로 돌아오고(`/privacy?lang=en` 형태) `<html lang>`·h1·title이 바뀐다 |
| 2. Tab 순회·포커스·`aria-current` | PASS | 6개 화면 모두 Tab 순서가 시각 순서와 같고 모든 포커스에 3px 윤곽이 보인다. 현재 정책 링크에 `aria-current="page"` |
| 3. 200%(640px)·320px 재배치, V-01 회귀 | PASS | 6개 화면 × 2폭(12 측정)에서 가로 넘침·잘림 0. 문장 속 링크 줄 간격이 본문 줄 높이와 같다 |
| 4. 로그인 회원 저장 언어 | PASS | 격리 fixture에서 EN 저장 회원은 `/privacy`·`/terms`·`/`를 영어로, KO 회원은 한국어로 본다 |
| 5. Google 로그인 scope에 calendar 없음 | 격리 PASS / 운영 NOT RUN | 운영 실측은 폼 제출이라 승인 범위 밖이다. 격리 `TestGoogleHTTP`가 `scope == "openid email"`을 단언하고 PASS |

## 1. 언어 전환 링크 — PASS

운영에서 쿠키를 비운 새 탭마다 페이지를 열고, 화면의 전환 링크(`/lang/en?back=…`, `/lang/ko?back=…`)로 이동했다. 증거: [prod-result.json](prod-result.json)의 `lang`.

| 시작 | `/lang/en?back=` 뒤 | `/lang/ko?back=` 뒤 |
|---|---|---|
| `/` (ko) | `/?lang=en`, lang `en`, h1 "Connect with friends, converse through agents" | `/?lang=ko`, lang `ko`, h1 "친구와 연결하고, 에이전트로 대화하세요" |
| `/privacy` (ko) | `/privacy?lang=en`, lang `en`, h1 "Privacy Policy" | `/privacy?lang=ko`, lang `ko`, h1 "개인정보처리방침" |
| `/terms` (ko) | `/terms?lang=en`, lang `en`, h1 "Terms of Service" | `/terms?lang=ko`, lang `ko`, h1 "서비스 약관" |

- 전환 뒤 현재 언어 링크에 `aria-current="true"`가 옮겨 간다.
- 관찰(결함 아님): 돌아온 주소에 `?lang=en|ko` 쿼리가 붙는다. 언어 쿠키가 지워진 브라우저에서도 같은 언어가 유지되도록 서버가 넣는 것으로 보인다. 쿼리의 목적은 코드로 확인하지 않았다.

## 2. 키보드 Tab 순회 — PASS

Chrome headless(1280px)에 CDP `Input.dispatchKeyEvent` Tab을 보내 `document.activeElement`를 읽었다. KO·EN × `/`·`/privacy`·`/terms` 6개 화면, 처음부터 `<body>`로 돌아올 때까지. 증거: [prod-result.json](prod-result.json)의 `tab`.

- 첫 포커스는 '본문으로 건너뛰기'/'Skip to content'. 그다음 브랜드·언어 링크 순이다.
- 포커스가 가는 순서의 세로 위치(top)가 6개 화면 모두 줄어들지 않는다. 즉 시각 순서와 같다.
- 모든 포커스에 `outline: solid 3px rgb(0, 90, 98)`(강조색 `#005A62`)이 있고 `:focus-visible`이 참이다. 포커스 요소는 모두 화면 안에 들어온다.
- 첫페이지 `#start` 정책 링크: 서비스 약관 → 개인정보처리방침 순으로 'Google로 계속' 바로 뒤에 온다(EN은 Terms of Service → Privacy Policy). 이어서 하단이 개인정보처리방침 → 서비스 약관 → 설치·사용 문서 → 맨 위로(EN은 Back to top)다. 첫페이지 하단에 '홈'은 없다(D04 일치).
- 정책 페이지 하단: 개인정보처리방침 → 서비스 약관 → **홈** → 설치·사용 문서 → 맨 위로, KO·EN 모두 같은 순서다.
- `aria-current="page"`: `/privacy`는 하단 '개인정보처리방침', `/terms`는 하단 '서비스 약관'에만 있다. 계산 스타일은 `font-weight 600`, `text-decoration none`이고 나머지 하단 링크는 400·밑줄이라 눈으로도 구별된다([prod_weight.mjs](prod_weight.mjs) 출력). designer V-02 수정이 운영에 반영됐다.
- 관찰: 정책 페이지의 표(보관 기간·위탁 업체)는 가로 스크롤 영역이라 Chrome이 자동으로 Tab 대상에 넣는다. 포커스 윤곽은 보이고 이상은 없다.
- 관찰(low, 결함 아님): `#start` 줄은 '서비스 약관 · 개인정보처리방침', 하단은 '개인정보처리방침 · 서비스 약관' 순이다. 같은 화면에서 두 링크 쌍의 순서가 반대다. designer가 정한 문안 순서인지는 확인하지 않았다.

## 3. 200% 확대·재배치 — PASS

200% 확대는 1280px 창에서 CSS 폭 640px과 같다. 실제 브라우저 확대 대신 CDP `Emulation.setDeviceMetricsOverride`로 폭 640px(deviceScaleFactor 2)과 320px(1)을 고정했다. 6개 화면(KO·EN × 3페이지) × 2폭. 증거: [prod-result.json](prod-result.json)의 `reflow`, [prod-inline.json](prod-inline.json).

- 가로 넘침: 12개 측정(KO·EN × 3페이지 × 640·320px) 모두 `scrollWidth == innerWidth`, 폭을 넘는 요소 0, 잘린 텍스트 0.
- V-01 회귀(문장 속 링크 줄 간격): 문장에 박힌 링크(`#start` 안내 줄, 개인정보처리방침 §6 Google API Services User Data Policy, 이메일, 정책 상호 링크)가 모두 `display: inline`, `min-height: 0`이다. 이 링크가 든 문단의 줄 간격은 25–26px(본문 줄 높이 25.6px), `#start` 줄은 22px(22.4px)로 정상이다. 가장 긴 문단(§6, 320px에서 8–9줄)도 같다. 측정: 1280·640·320px × KO·EN × 3페이지.
- 단독 줄 링크(하단·설정·언어 전환)는 44px 터치 대상을 유지한다(`inline-flex`, `min-height: 44px`).

## 4. 로그인 회원 저장 언어 — PASS

원격 Linux의 `TestPublishBrowserFixture`(합성 DB, loopback, 외부 호출 없음)에서 `POST /fixture/as/en|ko`로 합성 회원을 저장 언어 `en`·`ko`로 로그인시킨 뒤 GET으로 확인했다. 증거: [member_lang.py](member_lang.py), [member-lang-result.txt](member-lang-result.txt).

| 저장 언어 | `/privacy` | `/terms` | `/` |
|---|---|---|---|
| en | 200, `Content-Language: en`, h1 "Privacy Policy" | 200, en, "Terms of Service" | 200, en |
| ko | 200, `ko`, "개인정보처리방침" | 200, ko, "서비스 약관" | 200, ko |

같은 규칙을 `TestPolicyPagesPublic`(저장 언어 `en` 회원이 `/privacy`를 영어로 봄)도 단언하며 PASS했다([google-policy-run.log](google-policy-run.log)).

## 5. Google 로그인 scope — 운영 NOT RUN, 격리 PASS

- 지시서는 로그인 버튼의 '이동 대상 URL'을 읽으라고 했지만 운영 첫페이지의 'Google로 계속'은 링크가 아니라 `<form method="post" action="/auth/google">` 버튼이다([prod_login.mjs](prod_login.mjs)로 확인). 이동 대상 URL은 이 POST의 303 응답에만 있다.
- `googleStart`는 POST마다 DB 상태에 로그인 시도(`googleAttempt`: nonce·verifier)를 쓴다(`internal/relay/google.go:85`). 승인된 범위의 '쿠키 외 상태 변경 없음'과 '폼 제출 금지'에 어긋나서 운영 POST를 보내지 않았다. **운영 live scope 실측은 하지 않았다.**
- 대신 격리 환경(원격 Linux, 합성 Postgres, 로컬 토큰 서버)에서 `TestGoogleHTTP`를 실행했다. 단언: 303 응답 URL의 `scope == "openid email"`, `access_type`·`include_granted_scopes` 없음, `prompt=select_account consent`, PKCE S256 — PASS([google-policy-run.log](google-policy-run.log)). 코드: `internal/relay/google.go:50` `Scopes: {openid, email}`.
- 운영 값 확인이 필요하면 승인 후 한 번의 POST로 `Location`의 `scope`만 읽는다(redirect를 따라가지 않고, Google 동의 화면에 들어가지 않는다). 승인 전에는 이 항목을 '코드·격리 근거'로만 인계한다.
- Google 쪽 scope 삭제(D022)는 Cloud Console 상태이며 이 QA가 보지 않았다.

## 실행한 명령과 종료코드

| 명령 | 종료코드 | 산출 |
|---|---|---|
| `node prod_qa.mjs .` (운영 읽기 전용: 언어·Tab·재배치) | 0 | [prod-result.json](prod-result.json) |
| `node prod_inline.mjs` (문장 속 링크 줄 간격) | 0 | [prod-inline.json](prod-inline.json) |
| `node prod_login.mjs` / `node prod_weight.mjs` (요소 조사만, 클릭 없음) | 0 / 0 | 본문 인용 |
| `python3 remote.py prepare` (고정 SHA `24198e40` clone·`make install`) | 0 | 원격 |
| `python3 remote.py fixture` + `python3 explore.py member_lang.py` | 0 / 0 | [member-lang-result.txt](member-lang-result.txt) |
| `python3 remote.py run 'TestGoogleHTTP\|TestPolicyPagesPublic'` | 0, 두 테스트 PASS | [google-policy-run.log](google-policy-run.log) |
| `python3 explore.py stop.py` / 컨테이너 `docker stop` | 204 / 0, 남은 `kl-pub002` 컨테이너 0 | 원격 |

## 검증하지 못한 것

- 운영 Google 로그인 URL의 `scope` 실측(위 5절). 승인 필요.
- 실제 브라우저의 확대(Ctrl + +) 조작. 640px 에뮬레이션으로 대체했다.
- 스크린리더 읽기. `aria-current`와 DOM 순서만 확인했다.
- 운영의 로그인 회원 화면(실계정 금지). 격리 fixture로 대체했다.
- 정책 문안의 법적 적합성(P-05), designer D-01·D11 문서 후속은 범위 밖이다.

## 결함

제품 결함 0. 제품 코드를 수정하지 않았다. 새 시나리오 파일은 만들지 않았다(조작이 읽기·키 입력뿐이라 CDP 스크립트가 재실행 도구다. `node prod_qa.mjs <출력 폴더>`).
