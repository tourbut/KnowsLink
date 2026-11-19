---
title: SAR-PUBLISH-001-VISUAL designer 시각 검수
status: review
updated: 2026-10-22
owner: designer
tasks: [SAR-PUBLISH-001-VISUAL]
summary: "/privacy·/terms·첫페이지 하단·#start·설정의 KO/EN을 1440·390·320px에서 직접 본 PC-06 designer 판정. PASS, medium 1·low 3"
---

# SAR-PUBLISH-001-VISUAL designer 시각 검수

기준 main `17a360b2051d96e39646e1b048c931f6f0ede88b`, 제품 코드는 DEV 병합본 그대로다(HEAD `2ad0c6bb`에서 제품 파일 변경 없음). 정책 문안은 법률 검토(P-05) 전 초안이며 이 검수는 문안 내용을 평가하지 않는다. tester 독립 QA(SAR-PUBLISH-001-TESTER)를 대체하지 않는다.

## 판정

**PASS** — PC-06의 designer 몫(KO/EN, 3폭 직접 시각 검수)을 충족한다. critical·high 0, medium 1, low 3. medium V-01은 운영 배포 전 수정을 권한다. 배포를 막는 결함은 아니다.

## 검수 방법과 범위

- 렌더: 현재 제품 코드의 `go test ./internal/relay`가 쓰는 `KNOWSLINK_UI_FIXTURE_DIR` 출력(`policy-{privacy,terms,start}-{ko,en}.html`, `human-account`, `en-human-account`)을 실제 Chrome(headless=new)에 연결했다. 폭은 Chrome DevTools Protocol의 `Emulation.setDeviceMetricsOverride`로 1440·390·320px CSS 뷰포트를 고정했다([cdp_shots.mjs](cdp_shots.mjs)). DEV render의 iframe 크롭이 아니다. 8개 화면 × 3폭 = 24장을 직접 읽고 DOM 측정값([measure.json](measure.json))과 대조했다.
- 원격 격리 fixture는 쓰지 않았다. 접속 수단이 dev 워크트리의 SSH 자격(`.env.server`)이라 designer 범위 밖이다. 그래서 서버 경유 동작(언어 쿠키 전환, 로그인 상태 리다이렉트)은 보지 않았다. 아래 `검증하지 못한 것`에 적는다.
- 보지 않은 화면: 승인 상세(`gate-active`)는 이 과제 지시의 대상이 아니다.

## 항목별 판정

| 항목 | 결과 | 근거 |
|---|---|---|
| 스타일 재사용(PC-04) | PASS | 배경 `#FFFFFF`, 글자 `#182B31`, 강조 `#005A62`, 글꼴 첫페이지와 같다. 정책 페이지에 크림 배경·알약 버튼(radius ≥20px 요소 0)·장식 그라디언트(0)·이모지(0)·모노스페이스(0)·`details`·`script`·`img` 없음. 첫페이지·설정의 모노스페이스는 기존 설치 명령 `code`/`pre`와 친구 코드 입력란뿐이며 이번 변경이 아니다 |
| 고유 title·h1·`aria-current` | PASS | 4개 title: `개인정보처리방침 · KnowsLink`, `서비스 약관 · KnowsLink`, `Privacy Policy · KnowsLink`, `Terms of Service · KnowsLink`. 모두 h1 1개. `aria-current`는 현재 정책 링크와 현재 언어 링크 각 1개(measure.json `current`) |
| 가로 넘침(PC-04) | PASS | 8개 화면 × 3폭 모두 `scrollWidth == innerWidth`. 페이지 전체 넘침 0. 320px 영어 개인정보처리방침의 위탁 표만 표 영역 안에서 스크롤(304 > 288) — D04가 허용한 동작이다 → V-04 |
| 영어판의 한국어 잔존 | PASS | 영어 `/privacy`·`/terms`·`#start`의 한국어 텍스트는 언어 전환 이름 '한국어' 하나뿐이다. 하단 '설치·사용 문서'는 'Install and usage guide (in Korean)'로 표시된다 |
| 본문 폭·위계 | PASS | 1440에서 본문이 약 72ch 한 열이다([capture](captures/policy-privacy-ko-1440.png)). h1→시행일(`muted`)→상호 링크→번호 h2 순서, 표·목록 가독성 양호. 390·320도 한 열 |
| 대비·접근성 | PASS | `muted #52646B`/흰색 6.19:1, 강조 `#005A62`/흰색 7.96:1(4.5:1 이상). 본문 건너뛰기 링크(`.skip`)·`:focus-visible` 3px 윤곽 CSS 확인. 실제 키보드 순회는 보지 않았다 |
| 첫페이지 하단 순서 | PASS | KO·EN 모두 개인정보처리방침 · 서비스 약관 · 설치·사용 문서 · 맨 위로(D04). 첫페이지 하단에 '홈'은 없다 |
| `#start` 한 줄 | PASS | 로그인 버튼 아래 `muted` 한 줄에 서비스 약관·개인정보처리방침 링크(KO/EN 문안 일치, 체크박스 없음)([capture](captures/cmp-start-bottoms.png)). 줄바꿈 간격은 V-01 |
| 설정 링크 | PASS(low) | '서비스 안내' 옆에 '개인정보처리방침'·'서비스 약관'이 놓인다. 1440 KO, 320 EN에서 주 메뉴는 그대로다. 간격은 V-03 |
| 정책 페이지 하단 '홈' 위치 | **PASS(유지)** | DEV 미결 항목 판정은 아래 `홈 링크 판정` |

## 홈 링크 판정

정책 페이지 하단은 '개인정보처리방침 · 서비스 약관 · **홈** · 설치·사용 문서 · 맨 위로'다. 이 위치(서비스 약관 바로 뒤, 외부 문서 앞)를 승인한다. 이유: 같은 사이트 안 이동 링크가 한 묶음이고 외부 문서·맨 위로가 뒤에 와서 묶음 규칙이 깨지지 않는다. 헤더에 앵커 메뉴를 두지 않는 정책 페이지에서 브랜드 링크 말고 본문 끝에서 첫페이지로 돌아갈 길이 생긴다. 제품 정본(SAR-PUBLISH-001 화면 기준)도 '홈'을 요구한다. 첫페이지 하단에는 '홈'이 없어야 한다(현재 상태가 맞다).

D04 `service-member-ui.md`의 정책 페이지 ASCII 배치(`개인정보처리방침 · 서비스 약관 · 설치·사용 문서 · 맨 위로`)에는 '홈'이 없어 제품 정본과 어긋난다. 제품 코드의 결함이 아니라 D04 문서의 누락이다. 수정 지시안은 `수정 제안`의 D-01이다.

## 발견 사항

| ID | 심각도 | 화면·폭·언어 | 재현 위치 | 설명 | 수정 제안 |
|---|---|---|---|---|---|
| V-01 | medium | 전 폭·KO/EN. 줄바꿈이 있는 문단에 인라인 링크가 있을 때 | 첫페이지 `#start`의 '로그인하기 전에 …' 한 줄(390·320, [capture](captures/cmp-start-bottoms.png)), 개인정보처리방침 §6 'Google API Services User Data Policy' 문단(1440 KO, [capture](captures/policy-privacy-ko-1440.png)), 정책 페이지 상단 '… 함께 확인하세요' 한 줄이 줄바꿈될 때 | 전역 규칙 `a{display:inline-flex;min-height:44px}`(`member_pages.go:11`)가 문장 속 링크를 44px 높이 상자로 만든다. 링크가 든 줄만 줄 간격이 약 44px로 벌어져(정상 약 24px) 문장이 끊겨 읽힌다. 공개 화면의 새 인라인 링크가 처음 드러낸 결함이다 | `p a, li a`(본문 안 링크)에 `display:inline;min-height:0`을 적용한다. 단독 줄 링크(하단·설정 링크 줄)는 44px 터치 대상을 유지한다. 수정 뒤 390·320 `#start`와 개인정보처리방침 §6을 다시 본다 |
| V-02 | low | 전 폭·KO/EN, 정책 페이지 하단 | 하단 링크 줄([capture](captures/cmp-policy-mobile.png)) | 현재 페이지 링크가 `aria-current="page"`지만 시각 표시가 없다. 스크린리더에는 전달되나 눈으로는 같은 링크로 보인다 | 하단의 `[aria-current]`에 `font-weight:600`(또는 밑줄 제거)을 준다. 새 색·글꼴 추가 금지 |
| V-03 | low | 설정 1440 KO / 320 EN | 설정 상단 링크 줄([1440](captures/human-account-1440.png), [320](captures/en-human-account-320.png)) | '서비스 안내 개인정보처리방침 서비스 약관' 세 링크가 간격 없이 이어져 한 문장처럼 읽힌다. 320 EN에서는 'Terms of Service'만 다음 줄로 내려간다 | 링크 사이 가로 간격 16px 정도(`gap`)를 둔다. 기능 변경 아님 |
| V-04 | low | 320 EN 개인정보처리방침 | 위탁 표([capture](captures/policy-privacy-en-320.png)) | 표가 표 영역 안에서만 가로 스크롤한다(304 > 288). D04가 허용한 동작이며 페이지 넘침은 없다. 스크롤 가능하다는 시각 단서가 없다 | 수정 불필요. 필요하면 표 영역에 오른쪽 그림자 대신 열 폭을 줄여 스크롤을 없앤다 |
| D-01 | 문서 | — | `service-member-ui.md` 정책 페이지 배치도 | 하단 배치도에 '홈'이 없다 | 배치도 하단 줄을 `개인정보처리방침 · 서비스 약관 · 홈 · 설치·사용 문서 · 맨 위로`로 고치고, 첫페이지 하단에는 '홈'이 없다고 한 줄 적는다. coor가 D04 변경을 배정한다 |

## 완료 조건 대조(designer 몫)

| ID | 판정 | 근거 |
|---|---|---|
| PC-02 | 해당 없음 | 문안 일치는 DEV [parity.txt](../SAR-PUBLISH-001-DEV-test/parity.txt)가 맡는다. 이 검수는 영어판 한국어 잔존만 봤다 |
| PC-03 | PASS | 첫페이지 하단 순서·`#start` 한 줄·설정 링크·두 정책 페이지 상호 링크를 화면에서 확인 |
| PC-04 | PASS | 위 표. V-01은 스타일 재사용의 부작용이며 PC-04 문구(스타일 재사용·금지 요소 없음·가로 넘침 없음·고유 title·`aria-current`)는 충족 |
| PC-06 | **PASS(designer)** | KO/EN 3폭 직접 검수. tester 독립 확인은 별도 |

## 검증하지 못한 것

- 서버를 경유한 실제 화면(언어 쿠키 전환, `/lang/en?back=/privacy` 왕복, 로그인 회원의 저장 언어, 로그인 상태 리다이렉트). 렌더는 단위 테스트가 쓴 HTML이며 합성 값(회원 '합성 회원', 친구 코드 `synthetic-friend-code`)이 들어 있다. 설정 화면의 title('합성 대화'·'ascii-text')은 fixture 값이다.
- 실제 키보드 순회, 200% 확대, 다크 모드(이 제품은 단일 테마라 해당 없음).
- 테마 전환·디자인 lint: 단일 테마·디자인 lint 없음으로 해당 없음.
- Windows `go test ./internal/relay` 전체는 `TestFrozenParsingAndSigning`(`protocol_test.go:46 registry integrity failure`)이 실패한다. 이 과제와 무관한 기존 Windows 환경 문제로 보인다(DEV 보고는 원격 Linux에서 통과). 정책·설정 렌더 테스트(`TestPolicyPages`·`TestPolicyLinks`·`TestMemberUI` 계열)는 통과해 HTML이 생성됐다. 원인은 조사하지 않았다.

## 검사 기록

| 검사 | 결과 |
|---|---|
| `KNOWSLINK_UI_FIXTURE_DIR=… go test ./internal/relay -count=1` (HTML 생성) | 정책·설정 테스트 PASS, 전체 FAIL 1건(위 Windows 기존 실패) |
| `node cdp_shots.mjs <fixture 폴더> <출력 폴더>` (Chrome headless=new, CDP 1440·390·320) | EXIT=0, 24장 |
| 대비 계산(`--muted #52646B`, `--accent #005A62`, 흰 배경) | 6.19:1, 7.96:1 |

커밋하는 캡처는 위 발견 항목에 연결한 6장뿐이다. 나머지 18장은 로컬 임시 폴더에만 있다.
