---
title: 회원 웹 고정 후보 직접 시각 수락 보고
status: review
updated: 2026-10-11
owner: designer
tasks: [SAR-SERVICE-UI-001-VISUAL]
summary: "고정 후보 c7828b2의 합성 화면 29개를 Orca 내장 브라우저 1440·390·320에서 직접 보고 D04 수락 기준과 대조한 시각 수락 판정, 발견 사항, 한계"
---

# SAR-SERVICE-UI-001-VISUAL — 직접 시각 수락 보고

## 판정

**제품 시각 수락: PASS (비차단 잔여 있음).** critical·high 발견은 0건이다. medium 1건(V-01)과 low 10건(V-02–V-11)을 남긴다. 기능 QA와 코드 리뷰는 tester 담당이며 이 판정에 포함하지 않는다.

이 PASS는 합성 fixture 29개 화면, Chromium 계열 Orca 내장 브라우저, 3개 폭에 한정한다. 실제 Google 로그인 왕복, 실계정 데이터, 운영 서버, 실제 스크린 리더, 200% 확대, 다른 브라우저 엔진은 판정하지 않았다.

V-01은 D04 9절 "모든 상호작용 요소 높이 44px 이상" 기준을 글자 링크가 채우지 못한 결함이다. 오조작 위험이 낮고 같은 화면에 44px 주 동작이 있어 차단으로 올리지 않았다. 수정 여부는 coor가 정한다. 좁은 DEV 수정을 권고한다.

## 대상과 방법

| 항목 | 값 |
|---|---|
| 고정 후보 | `c7828b227b5715f7bdf02bff46bf9cf102a758d9` |
| 구현자 세션 | `5953ccd3-7f4f-4c5e-b63d-c02d36ee35ac` (DEV, dispatch `ctx_d554b501df04`) |
| 기록 체크아웃 HEAD | `030c5bd53c7645f06f5ac0b2ebbd11c4b31f5308` |
| 제품 바이트 | `git diff --stat c7828b2 HEAD -- internal adapters scripts go.mod`가 비었다. 후보와 같다. |
| 기준 main | `8539ef5`. `8539ef5..c7828b2` 제품 diff는 11개 파일(`.fullops-squad` 제외)이다. |
| 정본 | D04 [service-member-ui.md](../../../design-docs/mockups/service-member-ui.md) 필수 R-01–R-09, 권장 R-10·R-11 |
| 스킬 | `C:/Users/shin/.agents/skills/frontend-design/SKILL.md`를 읽었다. D04가 토큰·글꼴·간격을 고정하므로 정본을 따랐다. |
| 환경 | Orca 내장 브라우저 탭 1개. `file://` 합성 HTML. 응답 CSP 헤더는 정적 파일에 없다. |

1. 후보 파일을 `git archive c7828b2`로 scratchpad에 풀었다. 스냅샷·이 체크아웃에는 생성물을 쓰지 않았다.
2. `KNOWSLINK_UI_FIXTURE_DIR=<새 폴더> go test ./internal/relay -run TestMemberUIFixtures -count=1`이 통과했고 HTML 29개가 생겼다(Go 1.27.1, Windows).
3. [capture.py](capture.py)가 화면 29개 × 폭 3개(1440×900, 390×844, 320×640)를 열어 DOM을 측정하고 PNG를 저장했다. PNG는 파일로만 저장했고 base64를 출력하지 않았다.
4. [keyboard.py](keyboard.py)가 실제 키 입력(Tab·Enter·Space)과 렌더된 색 대비를 측정했다.
5. 직접 열어 본 PNG는 아래 "직접 관찰"에 적었다. 모든 PNG를 열어 보지는 않았다.

### 도구 한계와 폐기한 시도

- `goto`는 뷰포트 에뮬레이션을 초기화한다. 첫 실행은 `set viewport`를 `goto` 앞에 두어 세 폭이 모두 1316px로 측정되었다. 그 결과는 폐기했다. 최종 스크립트는 `goto` 뒤에 설정하고 `innerWidth == 목표 폭`을 assert한다.
- `full-screenshot`과 긴 뷰포트 캡처는 에뮬레이션에서 첫 화면을 반복해 이어 붙인다. 두 방식의 PNG는 폐기했다. 최종 PNG는 뷰포트 높이씩 스크롤하며 찍은 조각(`-sKofN`)이다.
- `keypress`는 `document.hasFocus()`가 false이면 무시된다. 키보드 프로브는 `h1`을 클릭해 포커스를 준 뒤 `tabindex=-1` 표식을 맨 앞에 두고 Tab을 시작한다.
- 폐기한 PNG는 이 과제가 직접 만든 임시 결과이며 폐기 시 해당 폴더만 지웠다.

## 증거

- [manifest.json](manifest.json): PNG 149개의 경로·뷰포트·상태·SHA256·크기. 합계 약 4.5MB.
- [measurements.json](measurements.json): 99개 키(화면 29개 × 폭 3개와 일부 details 펼침 상태)의 DOM 측정.
- [keyboard.json](keyboard.json): Tab 순서·포커스 외곽선·summary 토글·색 대비.
- PNG는 [png/](png/)에 있다. 파일명 규칙은 `<화면>-<폭>[-open][-sKofN].png`이다. `-open`은 모든 `details`를 연 상태다. `focus-*.png` 4개는 키보드 포커스 상태다.
- 모든 값은 합성 fixture다. 이메일은 `m***@example.test`, 토큰은 `synthetic-connection-token-not-a-secret-0000`이다. live 개인 정보·비밀값·base64 로그는 없다.

## 측정 요약 (DOM 측정, 99개 키)

| 검사 | 결과 |
|---|---|
| 가로 overflow (`scrollWidth <= innerWidth`) | 99개 모두 통과. 폭 밖으로 나간 요소 0개. 320px에서도 같다. |
| 외부 자원 (`src=`·`<link`·`url(`·`<script`·`@import`) | 0건 |
| `lang="ko"`·`h1` 1개·제목 단계 건너뜀 | 전체 통과 |
| 입력 라벨 | 보이는 입력 전체에 라벨이 있다(라벨 없는 입력 0개). |
| 본문 줄바꿈 | `word-break:keep-all`, `overflow-wrap:break-word`. `code`·`pre`는 `anywhere`. |
| 본문 글자 | 16px, 줄 높이 26.4px(1.65). 가장 작은 글자 14px. |
| 높이 44px 미만 상호작용 요소 | 글자 링크 11종만 해당(V-01). 버튼·입력·select·summary·이동 링크는 0개. |
| 색 대비(렌더된 계산값) | 본문 14.78, 보조 8.05, 주 버튼 8.07, 위험 버튼 8.01, 연결 완료 배지 5.74, 대기 배지 6.76, 철회 배지 7.43, 링크 6.98, 입력 테두리 4.43. D04 8.1 표와 같다. |
| 랜드마크·건너뛰기 링크 | 회원 화면은 `header`·`nav`(홈)·`main`·건너뛰기 링크가 있다. OAuth 원격 동의 화면 2개는 없다(V-10). |

폭 3개에서 `scrollWidth`가 `innerWidth`보다 15px 작은 것은 에뮬레이션의 세로 스크롤바 폭이다.

## 키보드와 포커스 (직접 입력)

- 첫 Tab이 "본문으로 건너뛰기"에 간다. 3px 파란 외곽선이 보인다. Enter가 `#main`으로 이동한다.
- 홈 Tab 순서: 건너뛰기 → 로그아웃 → 섹션 이동 5개 → 확인 필요 링크 2개 → 새 에이전트 만들기 → 에이전트 `summary`. 390px에서 시각 순서와 같다. 1440px는 두 열이며 같은 순서로 읽힌다.
- 에이전트 `summary`는 Enter로 열리고 Space로 닫힌다. 마우스 없이 위험 구역 `summary`와 "이 키 철회" 버튼에 도달한다. 위험 구역 안의 외곽선은 `--danger`(`rgb(160,22,22)`) 3px, 오프셋 2px이다([PNG](png/focus-danger-button-390.png)).
- 포커스된 요소 높이는 44px다(글자 링크 제외).
- 건너뛰기 링크가 보이는 동안 390px에서 제품명과 이메일 일부를 가린다. 포커스 때만 나타나므로 결함으로 올리지 않는다([PNG](png/focus-home-390-skiplink.png)).

## 직접 관찰 (PNG를 열어 본 화면)

| 화면 | 폭 | 관찰 |
|---|---|---|
| 로그인(`start-google`, `start-both-problem`, `start-none`) | 1440·390 | 제품명·한 줄 설명·3줄 목록·"Google로 계속" 주 버튼 한 개. 이메일 로그인은 "또는" 뒤에 있고 오류는 상단 `problem` 띠다. 로그인 불가 상태에 "로그인을 준비하고 있습니다" 안내가 나온다. R-09·5.1 충족. |
| 확인 코드(`verify`) | 390 | 오류 띠와 입력이 연결되어 있다(`aria-describedby="code-help problem"`, `inputmode=numeric`, `autocomplete=one-time-code`). 입력 폭은 240px다. D04 "큰 입력"보다 작다(V-11). |
| 홈 | 1440 | 헤더 → 확인 필요 띠 → 에이전트 → 관계 → 승인 요청 → 연결 확인 → 계정 순서. 왼쪽 섹션 이동 열이 고정되어 모든 조각에 보인다. 첫 화면에 에이전트 4개가 접힌 행으로 보인다. O-01·O-02 해소. 문서 높이는 1440px 3069px, 390px 3984px로 기준(6394·7986px) 대비 절반 안팎이다. |
| 홈 | 390 | 이동 링크가 가로 한 줄이다. 에이전트 행은 두 줄(식별자 / 배지·키 수)이다. 한국어가 단어 중간에서 끊기지 않는다(O-05 해소). 키 지문·식별자는 줄바꿈된다(F-UI-01). |
| 홈 펼침(`-open`) | 1440·390 | 키 목록·연결 추가·위험 구역이 한 에이전트 안에 차례로 나온다. 위험 구역은 위험 색 테두리 상자이고 결과 문장이 맨 위에 있다. 위험 버튼은 흰 배경 외곽선이다. 채운 위험 버튼은 없다. |
| 홈 빈 상태(`home-empty`) | 390·1440 | 재확인 안내 → "아직 연결한 에이전트가 없습니다" → 1·2·3 안내 → 주 버튼. 재확인 안내가 주 버튼보다 먼저 나온다(5.8). |
| 연결 관리 | 390 | 1·2·3 단계 표시. 현재 단계는 굵은 외곽선과 `aria-current="step"`이다. 대기 단계에 토큰 경고 문장과 코드 블록 2개가 있다. 준비 단계에서 지문이 큰 모노스페이스 블록으로 보이고 주 버튼은 "지문이 같습니다. 연결 승인"이다. |
| 연결 종료(`expired`) | 390 | 진행 표시 대신 안내 띠와 "홈으로 돌아가기"가 나온다. DOM에서 `expired`·`cancelled`에 버튼 0개, "연결 취소" 버튼 없음(F-UI-03). `consumed`는 3단계 "연결 완료"가 현재다. |
| 기기 승인(`device-prepared`) | 1440·390 | 맨 위 위험 색 왼쪽 선 경고 띠, 지문 블록, 권한 목록, 주 버튼과 "취소". D04 5.5 순서와 같다. |
| receipt | 1440·390 | 1440px에서 라벨/값 2열 정의 목록이다. 390px에서는 한 열로 쌓인다. "queued — … 상대 수신이 확인되지 않았습니다" 문장이 있다. |
| 오류(`refusal-limited`, `refusal-reauth`) | 390 | `!` 기호 + "오류" 굵은 라벨 + 문장. 다음 동작 링크가 있다(재확인 화면은 주 버튼까지). 오류를 성공 띠로 바꾼 곳은 없다. |
| gate(`gate-active`) | 390 | 토큰만 공유한다. 본문·hint 구분과 Approve/Deny는 그대로다. 영어 레이블이 남아 있다(범위 밖, V-08). |
| 원격 동의(`remote-consent`) | 390 | 새 토큰 글꼴·색만 적용된 이전 마크업이다(V-10). |

## 발견 사항

| ID | 등급 | 내용 | 근거 | 권고 |
|---|---|---|---|---|
| V-01 | medium | 글자 링크 11종이 높이 21px이다. D04 9절은 모든 상호작용 요소 44px 이상을 요구한다. 대상: "홈으로 돌아가기"(11개 화면), "홈으로 이동"(4), "상태 새로고침"(2), "작업 화면으로 돌아가기"(2), "다른 주소로 시작"(verify), "로그인 화면으로 이동", 확인 필요 링크 2개, 승인 요청 ID 링크 2개, 원격 동의 "내 agent 관리"(2). 모든 오류 화면의 유일한 다음 동작이 이 링크일 수 있다. 390px에서 "취소" 버튼 아래 8px 간격에 붙어 있어 오조작 가능성이 있다. | DOM 측정 99개 키 중 해당 화면 전체. 예: [connection-prepared-390](png/connection-prepared-390-s2of2.png), [refusal-limited-390](png/refusal-limited-390.png). DEV 검증표는 "버튼·입력·select·summary·이동 링크"만 셌다. 글자 링크는 포함하지 않았다. | `a`에 `display:inline-block; min-height:44px; padding-block:12px` 같은 터치 영역. 시각 크기는 유지한다. |
| V-02 | low | 받은 초대 행에 "수락·거절"과 "관계 철회"가 함께 있다. "관계 철회"는 연결된 관계의 동작이다(D04 5.6). 두 위험 버튼이 같은 행에 있어 뜻이 겹친다. | [home-390 s2](png/home-390-s2of5.png), [home-1440 s2](png/home-1440-s2of4.png). DEV가 form 계약 보존을 이유로 밝혔다. | 위치 조정 또는 받은 초대에서 철회 숨김. 서버 form 계약은 DEV가 판단한다. |
| V-03 | low | 만료 시각이 `2026-10-12 07:55:00 KST`다. D04 5.4는 "10월 12일 07:55까지" 같은 짧은 날짜를 우선하라고 했다. | PNG 여러 장. DEV가 미적용으로 밝혔다. | 후속 개선 후보. 기능 영향 없음. |
| V-04 | low | 320px에서 에이전트 `summary`의 펼침 표시(▸)가 식별자 위에 혼자 한 줄을 차지한다. 390px에서는 식별자와 같은 줄이다. | [home-320 s2](png/home-320-s2of8.png) | 마커를 식별자와 같은 줄에 두는 레이아웃. |
| V-05 | low | 기한 문장 "(최대 10분 ·1회 사용)"이 390px에서 줄이 바뀌며 "·"가 줄 앞에 온다. 위험 버튼 "이 에이전트 철회 (모든 키·관계 종료)"가 390px에서 "종료)"만 새 줄에 둔다. | [connection-token-390 s1](png/connection-token-390-s1of2.png), [home-390-open s2](png/home-390-open-s2of9.png) | 문구 줄이기 또는 `text-wrap:balance`. |
| V-06 | low | 에이전트가 0개이면 "새 관계 초대" 폼이 열린 채 빈 select와 "관계 초대" 버튼을 보인다. 선택할 에이전트가 없다. | [home-empty-390 s2](png/home-empty-390-s2of4.png) | 에이전트가 없으면 폼을 숨기고 한 줄 안내만 둔다. |
| V-07 | low (fixture) | `home-empty` 화면에서 "활성 /5"가 숫자 없이 나온다. fixture 입력에 `ActiveAgents`가 없기 때문이다. 실제 서버 `member.go` 297행은 정수를 항상 넣는다. 제품 결함이 아니다. 실제 서버 렌더는 직접 보지 않았다. | [home-empty-390 s1](png/home-empty-390-s1of4.png), `member_ui_test.go` 71행, `member.go` 297행 | tester가 실제 서버에서 "활성 0/5"를 확인한다. fixture는 DEV가 필요하면 고친다. |
| V-08 | low (범위 밖) | gate·receipt에 영어 표현이 남아 있다(`Approve 승인`, `Intent`, `typed body`, `queued`, `public-node`, `inbox`). D04는 gate에 시각 토큰만 공유하고 receipt 상태 문장을 유지하라고 했다. | [gate-active-390](png/gate-active-390-s2of2.png), [receipt-1440](png/receipt-1440.png) | 별도 문구 과제 후보. |
| V-09 | low | 재확인이 필요한 홈에서 "로그인 다시 확인" 버튼이 두 개(에이전트 안내, 계정 섹션)다. 이름이 같다. A-09는 같은 이름의 버튼을 구분하라고 한다. 두 버튼은 같은 동작이다. | `measurements.json`의 `home-empty@*` `dupButtonNames` | 한쪽에 `aria-label` 또는 한 곳으로 합친다. |
| V-10 | low | OAuth 원격 동의·로그인 화면 2개(`remote-consent`, `remote-login`)는 이전 마크업이다. `header`·`main`·건너뛰기 링크가 없다. "이 설치 연결 승인"이 주 버튼 스타일이 아니고 "취소"와 같은 보조 스타일이다. 본문에 "agent"가 섞인다. D04 4.1 화면 지도에 없는 화면이며 DEV가 마크업을 바꾸지 않았다고 밝혔다. | [remote-consent-390](png/remote-consent-390.png), `measurements.json`의 `landmarks`·`skipLink` | 후속 과제에서 같은 템플릿으로 정리. 보안 승인 화면이므로 주 동작 구분을 권한다. |
| V-11 | low | 확인 코드 입력이 폭 240px 일반 입력이다. D04 5.2는 "큰 숫자 입력 1개"를 말했다. 속성(`inputmode`, `autocomplete`)은 맞다. | [verify-390](png/verify-390.png), 입력 속성 DOM 조회 | 기능에 영향 없음. 선택 개선. |

미해결 critical/high: 없음.

## D04 수락 기준 대조

| 기준 | 판정 | 근거 |
|---|---|---|
| R-01 토큰 공통 적용 | PASS | 29개 화면이 같은 색·글꼴·버튼·배지를 쓴다. 대비 계산값이 8.1 표와 같다. OAuth 원격 동의는 토큰만 적용(V-10). |
| R-02 줄바꿈 | PASS | `keep-all`+`break-word`, 식별자·지문·코드만 `anywhere`. 390·320px에서 한글 단어가 끊기지 않았다. |
| R-03 홈 구조 | PASS | 4.2 순서와 같다. |
| R-04 접힌 행·위험 구역 | PASS | 접힌 행, 위험 구역 `details`, 결과 문장, 위험 외곽선 버튼. |
| R-05 진행 표시·지문 강조 | PASS | 1·2·3 단계, `aria-current="step"`, 큰 모노 지문. |
| R-06 문구 표 | PASS (부분 V-03) | 계정 비활성화 "준비 중"·내부 문구가 화면에 없다. 짧은 날짜는 미적용(V-03). DEV가 서버 안내문의 "agent" 용어를 유지했다고 밝혔다. |
| R-07 오류 다음 동작 | PASS | 확인한 오류 화면마다 링크 또는 주 버튼이 있다(V-01: 링크 높이). |
| R-08 접근성 A-01–A-10 | 부분 PASS | A-01·A-03·A-04·A-05·A-06·A-09(V-09 제외)는 측정·관찰로 확인했다. A-02는 원격 동의 제외(V-10). A-07은 `code`+라벨 확인(식별자 라벨 "식별자"). A-08은 아래 미검증. |
| R-09 로그인 설명 3줄 | PASS | 직접 관찰. |
| R-10 확인 필요 띠 | PASS | 홈에 받은 초대 1건·승인 요청 1건. |
| R-11 관계 정렬 | PASS | 받은 초대 → 보낸 초대 → 연결됨 → 종료됨 순서. |
| 9절 overflow | PASS | 99개 키 모두 `scrollWidth <= innerWidth`. |
| 9절 44px | medium 1건 | V-01. |
| F-UI-01 | PASS | 390·320 줄바꿈, overflow 0. |
| F-UI-02 | PASS | 오류 화면 링크 존재. |
| F-UI-03 | PASS | 만료·취소 화면에 "연결 취소" 버튼 없음(DOM). |
| F-UI-04 | PASS | 직접 본 화면의 모든 시각이 KST다. |
| 11절 불변식 | 해당 없음(시각 검수) | form·CSP·owner 분리는 tester 담당. DEV golden이 form 계약을 고정했다. |

## 디자인 lint와 테마

`project.md`는 서비스의 별도 디자인 lint와 공용 테마 전환이 구성되어 있지 않다고 정한다. 다크 모드는 D04 R-12 후속이다. 해당 없음이다. 등록 lint는 완료 전에 실행하고 결과를 완료 보고에 적는다.

## 직접 관찰·DOM 측정·worker 보고의 구분

- 직접 픽셀 관찰: "직접 관찰" 표의 PNG와 포커스 PNG 4개.
- DOM·계산 측정: `measurements.json`(overflow·44px·대비·줄바꿈 속성·랜드마크), `keyboard.json`(Tab 순서·외곽선).
- DEV 보고(재검증하지 않음): Linux gate·`go test`·integration·golden. DEV가 밝힌 사항(짧은 날짜 미적용, 서버 안내문 "agent" 유지, 200%·실스크린리더 미검증)은 위 V-03·R-06·미검증에 반영했다.
- DEV의 "44px 미만 0개"는 버튼·입력·select·summary·이동 링크 기준으로 내 측정과 같다. 글자 링크는 DEV 표에 없다.

## 미검증과 한계

- 실제 Google 인증 왕복, 실계정 데이터, 운영 서버 렌더는 보지 않았다. 정적 `file://` HTML이므로 응답 CSP 헤더와 실제 서버 값(예: 활성 에이전트 수 0)은 미확인이다.
- 200% 확대와 텍스트 간격 확대(A-08)는 측정하지 않았다. 320px overflow 0으로 간접 확인만 했다.
- 실제 스크린 리더 낭독은 확인하지 않았다. 랜드마크·라벨·`aria-current`는 DOM으로만 확인했다.
- Orca 내장 브라우저(Chromium 계열) 한 엔진만 봤다. Safari·Firefox는 미확인이다.
- 에이전트 5개와 긴 관계 목록은 fixture에 없다(에이전트 4개, 관계 5개).
- PNG는 뷰포트 조각이다. 한 장짜리 전체 화면이 아니다. 조각 사이 겹침은 마지막 조각의 스크롤 한계에서 생긴다.
- 29개 화면의 모든 PNG를 열어 보지는 않았다. 열지 않은 화면(예: `connection-cancelled`, `device-login`, `start-logout`)은 DOM 측정과 같은 템플릿 분기로 판단했다.

## 후속과 인계

1. coor: V-01 medium을 DEV 좁은 수정으로 보낼지 비차단 잔여로 수용할지 정한다. V-02–V-06·V-09–V-11은 선택 개선이다.
2. tester: V-07의 실제 서버 "활성 0/5" 확인. 기능 QA·고정 SHA 코드 리뷰는 tester 담당이다.
3. 재개 조건: DEV가 V-01을 고치면 해당 화면(링크 11종)만 다시 측정한다. 전체 재촬영은 필요 없다.
4. 배포 뒤 실제 Google 로그인·실계정 화면 확인은 사용자 소유 세션이 필요하며 OPS 후속이다.
5. D04는 수정하지 않았다. 수락 상태만 14절에 한 줄 기록했다.
