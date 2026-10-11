---
title: OAuth 승인 반환·오류 화면 직접 시각 수락
status: review
updated: 2026-10-11
owner: designer
tasks: [SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL]
summary: tester의 같은 후보 캡처 5장을 직접 열어 반환 화면과 재시작 안내 화면의 사용성을 판정하고 F1·F2의 시각 수락 여부를 기록한다.
---

# SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL — 직접 시각 수락

## 판정

반환 화면 desktop·mobile과 만료·CSRF·계정 binding 안내 화면은 **시각 PASS**다. F2(low, 링크 24px)는 **시각 수락 가능**하며 차단이 아니다. F1(medium, 원인별 문구 없음)은 비차단 잔여로 그대로 둔다. 새 시각 finding은 낮은 심각도 1건(V1)이다.

이 PASS는 5장의 합성 캡처와 tester 근거에 한정한다. 실제 Dots·ChatGPT 호스트 복귀, 실제 Google 계정, 운영 환경, Chromium 외 엔진은 판정하지 않았다.

## 대상과 원천

| 항목 | 값 |
|---|---|
| 제품 SHA | `eac603b6a41d2ce002e698600ac66b8fc67ad762` |
| tester 증거 SHA | `96c7b8236f8335b4e0c7344e09c00e3cb36c0e8b` |
| 기록 체크아웃 시작 HEAD | `16b802e80f87759f39f143af62e7da5809428306` |
| 캡처 유형 | 합성 Chromium 148, desktop 1280×800, mobile 390×844 |
| 입력 목록 | [visual-inputs.json](../SAR-SSE-INSTALL-001-OAUTH-RETURN-COOR/visual-inputs.json) |
| 근거 보고 | [tester report](../SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-test/report.md) |

이 체크아웃의 제품 바이트(`.fullops-squad` 제외)는 `git diff eac603b6… HEAD`가 비어 고정 SHA와 같다. 새 캡처·계정 조작·서버 접속은 하지 않았다. `visual-inputs.json`의 SHA256 5개를 실제 파일과 대조했고 모두 일치했다.

## 캡처가 새 코드에서 나왔는지

- `앱으로 돌아가는 중`은 `eac603b6`의 `remote_auth.go`에 1건 있고, 기준 `c399a5e`에는 0건이다. 반환 화면은 새 코드에만 있다.
- 반환 화면 문구(`앱으로 돌아가는 중`, `자동으로 이동하지 않으면 1분 안에 아래 링크를 누르세요.`, `연결을 시작한 앱으로 돌아가기`)는 `remoteReturn` template과 글자 단위로 같다.
- 안내 화면 문구는 `remoteConsent`의 `refusal` 인자(`원격 연결을 계속할 수 없습니다`, `문제: …`, `안내: …`)와 `member.go`의 `refusal`·`head` 정의 및 `자기 홈으로 돌아가기`와 같다.
- `1분`은 `approveRemote`의 code 만료(`now.Add(time.Minute)`)와 맞다.
- 기준 handler의 RED와 새 SHA의 GREEN은 tester가 실행했다. 이 보고는 그 결과를 연결할 뿐 다시 실행하지 않았다.

## 직접 관찰

| # | 캡처 (SHA256 앞 12자) | 크기 | 관찰 | 판정 |
|---|---|---|---|---|
| 1 | `approve-link-desktop-return.png` (`68f34815d1f4`) | 1280×800 | 제목 `앱으로 돌아가는 중`, 안내문, 파란 밑줄 링크가 왼쪽 본문 열(x=260)에 세로로 놓인다. 글자 넘침·가림·가로 잘림이 없다. 링크는 밑줄과 파란색으로 클릭 대상임을 알 수 있다. 화면에 callback 주소·code·state가 보이지 않는다. | PASS |
| 2 | `approve-link-mobile-return.png` (`504ce720d470`) | 780×1688 (390×844, DPR 2) | 같은 구성이며 좌우 여백이 있다. 안내문이 두 줄로 줄바꿈되고 링크는 한 줄에 온전히 보인다. 가로 스크롤이 없다. 비밀값이 보이지 않는다. | PASS (V1) |
| 3 | `expired-attempt-desktop-refusal.png` (`086157f79cbd`) | 1280×800 | 제목 `원격 연결을 계속할 수 없습니다`, 빨간 띠 `문제:`, 파란 띠 `안내: … 앱에서 KnowsLink 연결을 다시 시작하세요.`, 카드 안 `자기 홈으로 돌아가기` 링크. 재시작 경로가 문장으로 분명하다. 비밀값이 보이지 않는다. | PASS |
| 4 | `csrf-tamper-mobile-refusal.png` (`4c1ef7daf7d0`) | 780×1688 (390×844, DPR 2) | 같은 안내 구조가 모바일에서 줄바꿈되며 넘침·가림이 없다. 문제·안내·링크 순서가 읽힌다. | PASS (F1, V1) |
| 5 | `owner-binding-mobile-refusal.png` (`4c1ef7daf7d0`) | 780×1688 (390×844, DPR 2) | 4번과 같은 파일이다. | PASS (F1, V1) |

4번과 5번은 SHA256이 같고 픽셀이 같다. 이 한 장은 CSRF 변조와 계정 binding 두 조건에서 같은 화면이 나온다는 증거다. 두 조건을 이미지만으로 구분할 수 없다. 조건 구분은 tester 보고의 관찰(403·401, callback 0, install 0)에 의존한다.

색은 기존 `memberStyle`이다. 본문 `#17232b`, 배경 `#f5f7f8`, 링크 `#164da0`, `문제`는 빨간 띠, `안내`는 파란 띠다. 링크와 본문의 대비는 읽는 데 충분하다. 문제와 안내는 색만이 아니라 `문제:`·`안내:` 글자로도 구분된다. 새 CSS·테마·애니메이션·장식은 없다.

## 발견 사항

| ID | 심각도 | 내용 | 판정 |
|---|---|---|---|
| F1 | medium (tester) | 모든 consent 실패가 같은 문구를 쓴다. 4번·5번 이미지가 그 증거다. CSRF 불일치와 계정 불일치에도 `이미 처리되었습니다. 이전 요청은 다시 사용할 수 없습니다`가 나온다. agent 한도 초과(409)는 앱에서 다시 시작해도 풀리지 않는다. | 비차단 잔여로 보존. 이번 시각 수락의 범위가 아니다. 후속은 DEV가 오류 종류별 문구를 정한다. |
| F2 | low (tester) | 반환 링크가 한 줄 24px 높이(229×24)이며 카드가 없다. | 시각 수락 가능. 아래 근거. |
| V1 | low (새 기록) | `overflow-wrap:anywhere` 때문에 모바일에서 한국어 단어가 줄 중간에서 끊긴다. 2번 이미지는 `링크/를`, 4번·5번은 `처리되/었습니다`, `수/없습니다`, `연/결을`로 갈라진다. 글자는 잘리지 않으며 읽을 수 있다. | 비차단. 기존 `memberStyle`의 전역 규칙이며 이번 변경이 만든 것이 아니다. 원천 CSS를 바꾸는 별도 판단이 필요하다. 새 CSS는 추가하지 않았다. |

### F2 수락 근거

- 링크는 밑줄·파란색이고 본문 바로 아래 한 줄이라 찾기 쉽다. 화면에 다른 선택지가 없어 오조작 위험이 작다.
- tester가 같은 후보에서 실제 마우스와 터치 클릭으로 callback 도착과 PKCE 교환 200을 확인했다. 링크 상자는 desktop `[260,207,229,24]`, mobile `[20,232,229,24]`로 보고되었다.
- 링크 크기는 24×24 CSS px 이상이다. 접근성 최소 기준(WCAG 2.2 SC 2.5.8, 24px)에 맞는다. 권장 44px에는 못 미친다.
- 이 화면은 자동 이동(`meta refresh 0`)이 실패할 때만 사람이 보는 보조 경로다. 그래서 낮은 심각도 유지가 맞다.
- 더 큰 터치 대상은 선택 개선이다. 필요하면 DEV가 기존 `section`·버튼 형식을 재사용하는 별도 과제로 다룬다. 이번 수락 조건이 아니다.

## 확인하지 못한 것과 한계

- 캡처는 viewport만 담는다. 주소창·링크 hover 상태·`href`는 보이지 않는다. 링크 목적지는 tester의 코드·DOM 근거를 따른다. 승인 반환 화면의 HTML 소스에는 code가 있다. 화면에는 보이지 않는다. 이는 DEV·tester가 기록한 Location header와 같은 노출이다.
- 거절(deny) 반환 화면은 5장에 없다. 같은 template을 쓰므로 문구 `아래 링크를 누르세요`는 같다고 추정한다. 이 추정은 시각 PASS에 넣지 않는다.
- 로그인 없음·agent 한도 초과 안내 화면은 지정 5장에 없다. 별도 판정하지 않았다.
- mobile은 390×844 DPR 2의 Chromium이며 실제 휴대폰, 터치 기기, 스크린리더, Firefox, Safari, 앱 WebView가 아니다.
- 글꼴은 Linux Chrome의 `system-ui` 대체 글꼴이다. 사용자 기기 글꼴에서는 줄바꿈 위치가 다를 수 있다.
- 자동 복귀·token 교환 성공은 tester 근거로만 연결했다. 실제 Dots·ChatGPT 호스트에서는 확인되지 않았다.
- 기존 시각 PASS(`SAR-SSE-INSTALL-001-VISUAL`)는 동의 화면에 대한 것이며 이 새 화면의 수락으로 쓰지 않았다.

## 등록 lint와 제품 검증

제품 바이트가 고정 SHA와 같아 원 exact-SHA 제품 검증을 재사용한다. 근거는 [dev-eac603b6-final-gate.json](../SAR-SSE-INSTALL-001-OAUTH-RETURN-COOR/dev-eac603b6-final-gate.json)(`product-lint`·`product-test` 종료코드 0)과 [tester-96c7b823-final-gate.json](../SAR-SSE-INSTALL-001-OAUTH-RETURN-COOR/tester-96c7b823-final-gate.json)(같은 두 명령 passed, ERROR 0, WARNING 3)이다. 이는 이번 기록 HEAD의 새 제품 실행이 아니다.

문서 결과의 FullOps lint 실행 결과는 아래 완료 기록에 추가한다.

## 인계

- 시각 수락: 반환 화면과 안내 화면의 desktop·mobile 사용성 PASS. 미해결 critical/high 없음.
- 잔여 비차단: F1 medium(원인별 문구, DEV 후속), F2 low(수락), V1 low(기존 전역 줄바꿈 규칙).
- 부모 후속: main 통합, OPS 배포, 실제 Dots·ChatGPT 호스트와 사용자 계정에서의 복귀 확인.
- 제품 코드·tester 증거·원천 문서·`visual-inputs.json`은 변경하지 않았다.
