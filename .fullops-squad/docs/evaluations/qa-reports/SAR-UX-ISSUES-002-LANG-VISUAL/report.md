---
title: 언어·설정 후보 7704e6d의 두 언어 직접 시각 판정
status: review
updated: 2026-10-20
owner: designer
tasks: [SAR-UX-ISSUES-002-LANG-VISUAL]
summary: "후보 7704e6d의 설정·공개 첫페이지·주요 회원 화면을 한국어/English와 1440·390·320 캡처 90장으로 판정한다. 결론은 조건부 수락이며 로그인 회원의 언어 링크 무동작과 영어 1인칭 문구를 DEV 수정으로 남긴다"
---

# SAR-UX-ISSUES-002-LANG-VISUAL 직접 시각 판정

**판정은 조건부 수락이다.** 직접 본 화면은 U3 설정 순서·언어 선택 형태·공개 페이지 정보 구조와 U7 ‘설정’ 이름을 충족한다. 페이지 가로 넘침·44px 미만 동작·D04 밖 스타일은 0건이다. 조건은 V-LANG-01(medium)과 V-LANG-02(low, 영어 문구 수정 목록)다. critical·high는 없다.

대상은 고정 후보 `7704e6d3128831468bb474330210bc2a5099bce7`(fullops/dev, 제품 `2ff9e44`)다. 지시서 기준 main은 `08f2698`, 지시서 커밋은 `0b29416`이다. 제품 소스는 dev 워크트리(HEAD `7704e6d`)에서 읽기만 했다. 제품 코드는 고치지 않았다.

## 증거와 방법

- fixture: origin `http://127.0.0.1:18933`(coor 포워딩, 컨테이너 `kl-lang002-fixture-7704e6d31288-a8717b84`). ID는 manifest `C:/Users/shin/AppData/Local/KnowsLinkDevTools/lang-ux002-dev/fixture-manifest.json`에서 읽었다.
- 요청: 허용 GET과 세션용 `POST /fixture/as/{en|ko|none}`만 보냈다. 언어 저장·이름·친구·에이전트 POST는 보내지 않았다. 확인용 `GET /lang/ko?back=/` 1회는 응답 쿠키만 받고 버렸다(서버 저장값 변경 없음).
- 언어 쿠키: `__Host-kl_lang`은 Secure라 headless Chrome이 http 응답에서 버린다. `GET /lang/{code}`가 쓰는 것과 같은 이름·값·속성을 CDP로 넣었다(curl로 `/lang/en` 응답의 `Set-Cookie: __Host-kl_lang=en; Path=/; Max-Age=31536000; HttpOnly; Secure; SameSite=Strict`를 확인).
- 도구: [capture.mjs](capture.mjs). headless Chrome(비모바일 viewport, 배율 1)으로 1440·390·320 전체 높이를 캡처한다. 측정값은 페이지 가로 넘침, 44px 미만 동작, main 안의 가운데 정렬, 한글 잔존(사용자 데이터 제외 판정), `html lang`, 언어 링크 위치·`aria-current`, radio 상태, `select`·`script` 수다. 결과는 [captures/probe.json](captures/probe.json)(90장, `NO_OVERFLOW 90`, 종료코드 0, [capture.log](capture.log))이다. 1440 캡처마다 본문 텍스트(`*.txt`)를 남겼다.
- 조합: `en`(English 저장), `ko`(한국어 저장), `none`(저장값 없음), `none`+쿠키 en, `en`+쿠키 ko, 비로그인 쿠키 없음, 비로그인+쿠키 en.
- 허용 GET으로 열 수 없는 화면은 DEV 렌더 `C:/Users/shin/orca/workspaces/KnowsLink/dev/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-DEV-test/render/`(아래 `R/`)를 직접 열었다. `GET /auth/verify`는 보낸 코드가 없으면 첫페이지로 돌아가 fixture로 만들 수 없다.
- D04 스타일: 후보 diff(`1315448..2ff9e44`, 테스트·번역표 제외)의 새 CSS는 `.lang{display:flex;gap:8px;font-size:14px}`와 현재 언어 `font-weight:600`뿐이다. 새 색·그림자·그라디언트·모서리는 0건이다.
- 브라우저 프로필: 이 폴더의 `profile/`(`.gitignore`로 커밋 제외). 저장소 규칙에 따라 지우지 않는다. 정리는 coor가 정한다.
- 영어 문구: `internal/relay/i18n_en_member.go`·`i18n_en_human.go`(7704e6d) 전체를 읽고 캡처 본문과 대조했다.

## 화면별 판정

90장 모두 페이지 가로 넘침 0, 44px 미만 동작 0이다. main의 가운데 정렬은 버튼 라벨과 공개 첫페이지의 기존 CTA 링크 버튼뿐이다. 영어 화면에서 남은 한글은 합성 사용자 이름 `긴한글이름…`과 언어명 `한국어`뿐이다. 모든 화면의 `html lang`이 표시 언어와 같다. `script`는 0개다. 입력란 안의 가로 스크롤(`INPUT`)은 긴 값이 입력란 안에서 넘기는 기본 동작이며 한국어판에도 같다. 캡처 경로는 `captures/` 기준이다.

| 화면 | 근거 캡처 | 판정 | 근거 |
|---|---|---|---|
| 설정(English) | `en-home-account-{1440,390,320}.png` | PASS | 순서가 My name → My friend code → Language → Sign-in and security다. Language는 radio ‘한국어’·‘English’와 ‘Save language’다. `select`·국기·JS 없음. 헤더 ‘Settings’가 `aria-current=page`다. 320에서 도움말이 줄바꿈되고 잘리지 않는다. |
| 설정(한국어) | `ko-home-account-*.png`, `none-home-account-*.png` | PASS | 내 이름 → 내 친구 코드 → 언어 → 로그인과 보안. ‘언어 저장’. 저장값 없음(`none`)은 한국어와 radio ‘한국어’ 선택이다. |
| 저장 성공 안내 | `en-home-account-n-language-*.png` | PASS | ‘Language saved. Screens from now on are shown in this language.’가 같은 설정 위치에 status로 보인다. 저장 POST 없이 `?n=language` GET으로 표시만 확인했다. |
| 회원 저장값 우선 | `en-cookie-ko-home-account-*.png` | PASS | 회원 en + 쿠키 ko에서 영어와 radio English다(U3 표 1행). |
| 비로그인 선호를 가진 저장값 없는 회원 | `none-cookie-en-home-account-*.png` | 조건부 | 영어로 보이고 radio English가 선택된다. 화면 적용은 U3 표 4행과 같다. 계정에 아직 저장되지 않았다는 표시가 없다. V-LANG-03(low). |
| 공개 첫페이지 두 언어 | `anon-ko-*.png`, `anon-en-*.png` | PASS | 제목 6개·본문 106줄의 구조가 같다. Grok/Dots 준비물·명령·README 링크(같은 대상, 한글 fragment는 percent-encoding)·가상 예시 표시가 같다. 영어판은 설치·로그인 성공을 새로 약속하지 않는다(‘These settings alone do not complete the connection’ 등). |
| ‘한국어 / English’ 링크 위치 | 같은 캡처, `anon-en-connect-*`, `en-remote-consent-*` | PASS | header 안의 작은 14px 링크다. 1440은 오른쪽 위 메뉴 뒤, 390·320은 메뉴 다음 줄 왼쪽이다. 각 링크 높이 44px, 현재 언어는 굵은 글자와 `aria-current`, 링크마다 `lang` 속성이 있다. 동의 화면은 320에서도 브랜드와 같은 줄이다. |
| 로그인 회원의 공개/동의 화면 링크 | `en-{1440,390,320}.png`, `en-remote-consent-*.png` | FAIL | 아래 V-LANG-01. |
| 대화 목록·대화 | `en-home-*.png`, `en-home-conversation-*.png` | PASS | ‘Today 7:49 PM’, ‘Written by me’, 만료 안내. 긴 한글·영문 이름이 320에서 줄바꿈되고 잘리지 않는다. 320에서 주 메뉴 ‘My agents’가 둘째 줄로 내려간다(가로 넘침 없음, 비차단). |
| 새 대화·내 이름 | `en-home-new-*.png`, `en-home-name-*.png` | PASS | 라벨·버튼·도움말 영어, 44px. |
| 친구 | `en-home-friends-*.png` | PASS | 받은 요청 Accept/Decline과 거절 앞 안내, 보낸 요청 ‘Awaiting reply / Entered code …zz99’, 친구 코드·요청 폼. 의미가 U4 한국어판과 같다. |
| 내 에이전트·Claude 선택 | `en-home-agents-*.png`, `en-home-agents-use-claude-start-first-*.png` | PASS(문구 제외) | 5선택지·처음/기존 연결·상태 행 재개 안내가 한국어판과 같다. ‘(my computer)’ 등 1인칭은 V-LANG-02. |
| 승인 | `en-home-gates-*.png` | PASS | Who to whom / What was requested / What approval allows / Deadline and status. KST 기한 문자열 그대로. 승인·거절이 글자와 색으로 구별된다. |
| 요청 상태 조회 | `en-home-receipts-*.png`, `en-home-receipts-id-*.png` | PASS(문구 제외) | 위치명 ‘My agents / Request status’, radio 라벨, 오류 ‘Select my agent.’(V-LANG-02). |
| legacy | `en-home-legacy-*.png` | PASS | ‘On this screen’ 앵커 목록, `aria-current` 없음. |
| 연결 단계 3종 | `en-home-connections-*-*.png` | PASS(문구 제외) | 단계·지문·권한·취소. ‘on my computer’는 V-LANG-02. |
| device·remote 동의 | `anon-{ko,en}-connect-*.png`, `en-remote-consent-*.png` | PASS(문구 제외) | 주의·지문·권한·Google 확인. remote 설명이 한국어와 같은 범위다. |
| 확인 코드·거부·Google 완료 | `R/en-verify-*`, `R/en-refusal-reauth-*`, `R/en-google-complete-home-*` | PASS | 오류 alert, 남은 시도, 재확인 버튼, 정확한 복귀. |
| remote 앱 복귀 | `R/en-remote-return-*` | not_run | DEV 렌더 3장이 `meta refresh` 뒤의 Chrome 오류 화면(`app.invalid`)이다. 화면을 볼 수 없어 템플릿·번역표 문구만 확인했다(‘Returning to the app’, ‘This screen does not mean the connection is complete’). |

## 제품 판단

- **admission 429/503 거부 화면의 언어: 허용한다.** 회원 확인 전에 부하를 막는 일시 화면이라 브라우저 언어 쿠키(없으면 한국어)를 따라도 된다. 과부하 중에 회원 저장값 조회를 추가하지 않는다. 조건은 문구와 `lang`·`Content-Language` 일치, 재시도 시각과 안전한 복귀, 저장 언어를 적용했다고 말하지 않음이다. 소스(`capacity.go` `refused` → `refusedTo(…, "/home")`, `limited(retry)`)가 이 조건을 충족한다. [U3](../../../design-docs/SAR-UX-ISSUES-002.md)의 ‘2026-10-20 LANG-VISUAL 결정’에 반영했다.
- **로그인 회원의 언어 링크(V-LANG-01):** 저장값이 있는 회원에게는 동작하지 않는 언어 전환 링크를 보이지 않는다. 대신 현재 언어와 설정 언어 위치로 가는 링크를 보인다. 아래 수락 조건을 따른다.
- **번역 용어:** 메뉴·제목의 고유 이름(My agents, My name, My friend code)은 유지한다. 문장·선택지·안내·오류에서 사용자를 가리키는 말은 ‘your’로 쓴다.

## DEV 수정 필요

| ID | 등급 | 내용 | 수락 조건 |
|---|---|---|---|
| V-LANG-01 | medium | 저장값 en 회원이 공개 첫페이지·remote 동의 화면에서 ‘한국어’를 누르면 `/lang/ko` → `/?lang=ko`로 가지만 회원 저장값이 이겨 영어가 그대로다(curl: `html lang="en"`, `aria-current` English). 링크가 아무 일도 안 한 것처럼 보이고, 브라우저 쿠키만 바뀌어 로그아웃 뒤에 다른 언어가 나온다. | 회원 저장값이 화면 언어를 정할 때 이 두 링크 대신 ‘English · Change in Settings’/‘한국어 · 설정에서 변경’(→ `/home/account#language`)을 보인다. 저장값 없는 회원과 비로그인은 지금 링크를 유지한다. 자동 검사로 세 경우의 링크를 단언한다. 문구·위치만 바뀌므로 tester 확인으로 충분하다. |
| V-LANG-02 | low | 영어 문장에서 ‘내 …’를 1인칭 ‘my’로 옮겼다. 아래 목록. | 목록 문구를 고친다. `TestEnglishScreensHaveNoKorean` 같은 렌더에서 메뉴·제목 고유 이름 외의 `\bmy\b`가 0임을 확인한다. designer 재검수는 필요 없다. |
| V-LANG-03 | low | 저장값 없는 회원이 쿠키 언어로 설정을 열면 radio가 쿠키 언어로 선택돼 이미 저장된 것처럼 보인다. | 저장값이 없을 때 Language 절에 ‘This browser’s choice is shown. It is not saved to your account yet.’/‘이 브라우저의 선택을 표시합니다. 아직 계정에 저장되지 않았습니다.’를 보인다. 다음 단위와 묶어도 된다. |

### V-LANG-02 영어 문구 수정 목록

키는 한국어 원문, 파일은 7704e6d 기준이다.

| 파일 | 원문(키) | 현재 | 수정 |
|---|---|---|---|
| member | 내 컴퓨터(Node 22)에서 저장소의 adapters를… | Install the repository's adapters on my computer (Node 22)… | …on your computer (Node 22)… |
| member | 지문 확인: 내 컴퓨터에 표시된… / 내 컴퓨터에 표시된 지문과… | …shown on my computer | …shown on your computer |
| member | 연결 수단 발급됨: 내 컴퓨터에서 준비 | Connection token issued: prepare on my computer | …prepare on your computer |
| member | 승인했습니다. 아직 새 키가… 내 컴퓨터에서 다음 명령을… | Run the following command on my computer. | …on your computer. |
| member | {{badge "등록됨"}} 키 등록됨… 내 컴퓨터에 보관합니다… | …credentials on my computer. | …on your computer. |
| member | 개인키는 내 컴퓨터에만 보관합니다(3곳) | The private key stays only on my computer. | …only on your computer. |
| member | 내 컴퓨터에서 시작한 요청만 승인하세요… | Approve only requests started from my computer. | …started from your computer. |
| member | 내 계정에 새 에이전트 하나를 연결합니다… | One new agent is connected to my account. | …to your account. |
| member | 승인했습니다. 내 컴퓨터에서 저장 결과를… | Check the save result on my computer. | …on your computer. |
| member | 연결 완료는 키가 등록되었다는… 내 컴퓨터에서 실행합니다. | Run the connection on my computer where Node 22 is installed. | …on your computer where… |
| member | 비민감 연결 확인 송신… 선택한 내 에이전트의 Node client… | …from the Node client of the selected my agent. | …of your selected agent. |
| member | 로그인하고 내 이름 정하기 / 내 에이전트 연결하기 | Sign in and choose my name / Connect my agent | Sign in and choose your name / Connect your agent |
| member | 새 대화에서 친구, 내 발신 에이전트… | …pick the friend, my sending agent and their receiving agent. | …your sending agent… |
| member | 내 에이전트끼리는 친구 요청 없이… | My own agents can talk to each other… | Your own agents can talk… |
| member | 웹에서 선택한 내 에이전트 명의로… | …in the name of the selected my agent. | …in the name of your selected agent. |
| member | 내 에이전트로 대화 시작하기 / 로그인한 뒤 내 이름과 에이전트를… / 같은 Google 계정으로 내 에이전트를… | Start a conversation with my agent / check my name and agents / Keep using my agents… | …with your agent / check your name and agents / Keep using your agents… |
| member | 내 클라이언트 연결(제목) | Connect my client | Connect your client |
| member | 내 에이전트를 선택하세요. / …내 에이전트와 요청 ID를 확인하세요. / …사용할 내 에이전트를 다시 선택하세요. / …내 에이전트를 다시 선택하세요. | Select my agent. / Check my agent and the request ID. / Select which of my agents to use again. / Select my agent again. | Select your agent. / Check your agent and the request ID. / Select your agent again. / Select your agent again. |
| member | 내 에이전트 관리 | Manage my agents | Manage your agents |
| human | Grok Bot (내 컴퓨터) / Claude Code 작업 세션 (내 컴퓨터) | Grok Bot (my computer) / … (my computer) | Grok Bot (your computer) / … (your computer) |
| human | 키 등록 전 · 내 컴퓨터의 준비 대기 | …waiting for my computer to prepare | …waiting for your computer to prepare |
| human | 표시 언어 도움말(저장하면 다음 화면부터…) | …Choosing alone changes nothing, and you are not signed in again. | …Selecting a language changes nothing until you save. You stay signed in. |
| human | 지원하지 않는 언어입니다… 한국어 또는 English를 선택하세요. | …Choose Korean or English. | …Choose 한국어 or English. (보이는 라벨과 같게) |

비차단 관찰(수정 선택): 영어 320에서 주 메뉴 ‘My agents’가 둘째 줄로 내려간다. legacy의 ‘Needs your check’는 ‘Needs attention’이 더 자연스럽다. ‘Approval only passes this request’s human gate’의 ‘human gate’는 내부 용어에 가깝다.

## 실행하지 못한 것 (not_run)

- 실제 언어 저장·실패(503)·재로그인·계정 교대 화면 전이. 지시서대로 저장 POST를 보내지 않았다. 저장 성공 안내는 `?n=language` GET 표시로만 봤다.
- 실제 `/lang/{code}` 클릭 흐름의 브라우저 쿠키 저장. Secure 쿠키라 http fixture에서 Chrome이 버린다. 같은 쿠키를 CDP로 넣었다.
- 실제 키보드 순회, 브라우저 실제 200% 확대, 스크린리더, 보조 브라우저(U7-A).
- remote 앱 복귀 화면의 시각 확인(DEV 렌더가 오류 화면). 확인 코드 화면은 DEV 렌더로 봤다.
- admission 429/503 화면 캡처. 소스로만 확인했다.
- tester 산출물 `C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-TESTER-test/`는 작성 중이라 이 판정에 쓰지 않았다.
