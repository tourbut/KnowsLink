---
title: 공개 첫페이지 직접 시각 수락과 접근성 후속
status: review
updated: 2026-10-14
owner: designer
tasks: [SAR-LANDING-GUIDE-001-VISUAL]
summary: 고정 공개 화면 직접 시각 수락과 기존 회원 대비 실패 및 native 입력 후속
---

# 공개 첫페이지 직접 시각 검수

**판정: 공개 화면의 시각 PASS, low 1건. 회원 대표 화면에는 기존 대비 기준 FAIL 1건(medium)이 있다.** 전체 관찰 범위의 critical/high는0건이다. D04 전체 접근성·독립 인증 QA·운영 수락은 이 판정에 포함하지 않는다. native 입력의 tester 확인 범위와 미완료200%는 아래에 구분한다.

## 고정 대상과 독립성

- 제품 후보: `fe2387a7676f7997aea76a2606f4cc2dbcd4b56b`.
- 실제 handler/합성 DB fixture: `4795d619513fde4fbcd4a55b94ec91b056ebbaf0`.
- 착수/원천 HEAD: `297e6d3679f6d866171a35775b615685221ffd47`.
- designer 세션: `01a13ac6-2571-7631-b463-91367eb88529`, terminal `term_a631a267-d383-468b-8e0e-c209f4e1e476`.
- 구현자 세션: `01a13a96-02a2-7863-a346-67d01b79d394`. 구현자와 다른 세션이다.
- Task `task_e4422a1d293e`, Dispatch `ctx_d880f510c28f`, attempt `93adb3b5b2124aed99b9c02f45e39281`.
- 실행: codex/gpt-6-astra/high, 주 스킬 frontend-design, test lite. 새 의존성·보조 디자인 도구·하위 worker는 사용하지 않았다.
- 기준: fullops-common-0.3.3, project.md, D02 LG-01–08, D04 공개 첫페이지 절, 문서 규칙. D04와 제품 기준은 변경하지 않았다.

[source-identity.json](source-identity.json)은 fixture→제품 후보→착수 HEAD 사이 제품·fixture·의존성·빌드 입력의 `git diff --exit-code` 종료0을 기록한다. fixture의 SHA·origin·고유 DB는 [DEV manifest](../SAR-LANDING-GUIDE-001-DEV-test/browser-fixture.json)와 대조했다. 자체 hero 3장은 DEV hero와 SHA256까지 같다. 런타임 정체성은 이 제공 manifest와 입력 동일성에 근거하며 서버 내부 프로세스를 새로 감사한 것은 아니다.

## 실제 실행과 이미지

Orca CLI `1.4.224`의 내장 Chromium 페이지를 사용했다. browser engine 상세 버전은 별도 조회하지 않았다. isolated profile은 `023bac7a-dec0-48eb-89fe-bb09db4e1cb7`, page는 `87f40b80-f6f6-49f0-ab06-13bec9ac312e`다. DEV/tester의 페이지와 쿠키를 조작하지 않았다.

허용 origin은 `http://127.0.0.1:18781`(Google), `:18782`(both), `:18783`(email), `:18784`(none)이다. 이번 요청은 공개 GET `/`, `/?n=expired`, `/auth/google/callback?error=access_denied`뿐이다. 접기와 앵커는 현재 DOM의 native 요소 `.click()`으로 조작했다. 실제 POST·설치·외부 링크 이동·실계정·실메일·송신은 실행하지 않았다. 직접 callback GET은 안전한 오류 표시 관찰이며 실제 OAuth 취소 왕복이 아니다.

새 캡처16장과 DEV 원 PNG17장을 모두 직접 열었다. [자체 manifest](captures.json)는 실제 innerWidth/innerHeight·DPR·CSS zoom·scroll·origin·page·SHA256을 보존한다. [재사용 manifest](reused-captures.json)는 DEV 원17장의 hash 검증과 용도를 기록한다. DEV 실패 타일1장은 수락 근거에서 제외했고 CSS 확대1장은 native 확대 증거로 쓰지 않았다. 원 실패 이미지를 삭제하거나 덮어쓰지 않았다.

coor 회신 뒤 tester의 같은 fixture 회원 PNG2장을 추가로 직접 열었다. [tester manifest](tester-captures.json)는 원 경로·SHA256·제품 SHA·Chrome154.0.8037.98·target을 연결한다. 두 원본을 own QA의 `tester-member-*.png`로 복사해 보존했다. 실제 CSS1440×900/DPR1/browserZoom1이며 이번 designer가 인증하거나 생성한 화면은 아니다. 총35 PNG를 직접 열었다.

| 증거 묶음 | 실제 조건 | 판정에 사용한 범위 |
|---|---|---|
| 자체 `hero-1440/390/320.png` | CSS1440×900, 390×844, 320×844 | 첫 화면 CTA, 제목·관계 도식·위계 |
| DEV `browser-layout/{install,examples,start}-{1440,390,320}.png` | 각 이름의 실제 CSS 폭; DEV viewport 보정 근거 연결 | 설치 상단, 가상 표시·예시 두 개, both 폼 |
| 자체 `install-closed-320.png`, `grok-values-320.png`, `dots-open-320.png` | CSS320×844 | 접힌 요약, 긴 명령/설정, Dots URL과 안내 |
| 자체 `google-only-start-390.png`, `email-only-start-390.png`, `none-start-390.png`, `expired-*-390.png`, `google-error-*-390.png` | **파일명과 달리 실제 CSS1000×841** | provider/안내 상태의 보조 시각 증거. 모바일390 증거로 사용하지 않음 |
| 자체 `google-error-top-actual390.png`, `google-error-form-actual390.png`, `none-actual390.png` | CSS390×844 | 보정 후 오류 발견·재시도와 none 상태 |
| DEV `browser-states/authed.png`, `google-error-*.png` | DEV 상태 실행의 CSS1000 계열; PNG1096×921 | 로그인된 공개 hero와 오류 상단/폼의 기존 실제 화면 |
| tester `member-home-1440.png`, `member-account-1440.png` | CSS1440×900/DPR1/zoom100%, 별도 Chrome CDP | 빈 회원 대화·계정/서비스 안내. 실제 이미지 직접 판정 |

Orca의 inherited 배율은 DPR 기준 약1.095445다. CSS1440 캡처에서는 device scale0.65를 사용해 DPR약0.712039, PNG1025×641이다. mobile은 DPR약1.095445이며 PNG427×925 또는351×925다. PNG 너비를 CSS viewport로 해석하지 않는다. 모든 자체 화면에서 `body.style.zoom`은 비어 있다. native browser zoom의 정확한 percent는 확인하지 못했다.

[run-metadata](run-metadata.json)는29개 의도적 조작과40개 상한을 기록한다. 첫 자체 hero부터 마지막 none390까지 캡처 구간은172.1초다. 관측/캡처 호출은 조작 수에서 제외했다. 모델 없는 결정적 CLI 실행이며 시각 판정은 이 세션에서 직접 수행했다.

실행 도중 navigation이 viewport를 CSS1000×841로 되돌렸다. 원16장 manifest의 실제 측정을 보존했고 오류·none만 올바른390으로 추가 캡처했다. 이 증거 결함 때문에 수행한 재캡처이며 제품 결함이 아니다. `exec --command 'help set'`은 `Unknown command: help`였고 브라우저 기능 PASS로 쓰지 않았다.

## D04 항목별 판정

| 항목 | 판정과 직접 관찰 |
|---|---|
| 소개·CTA·관계 | PASS. 1440에서 두 열과 왼쪽 정렬을 유지한다. 390 제목은2줄, 320은3줄이며 CTA는 세 폭 모두 첫 viewport 안이다. 320 도식은 세로로 쌓여 사람 아래 agent 관계를 읽을 수 있다. |
| 시작 순서 | PASS. 단일 세로 번호 목록이며 이름→agent→친구 수락→대상/본문 확인 순서를 유지한다. 카드나 설치 명령이 hero를 밀어내지 않는다. |
| Grok/Dots | PASS. 처음 두 details가 닫히고 요약은 보인다. 320에서 둘을 동시에 펼쳐도 페이지 가로 넘침이 없다. Command/Arguments·신규 설정은 줄바꿈하며 실제 textContent는 원 문자열을 보존한다. Dots URL은 선택 가능한 텍스트다. 신규/업데이트 안내가 구분된다. |
| 수신·예시 | PASS. 가상 표시, 웹 직접 작성과 agent 발송 요청, 매 발송 승인, 수신/표시/답장 차이가 읽힌다. 1440은 두 예시 열, mobile은 세로 순서다. 실제 성공 사례나 입력창처럼 보이지 않는다. 180초·24h 설명을 유지한다. |
| provider·오류·none | PASS(표시). Google-only/email-only/none에는 해당 상태의 행동만 보인다. both의 라벨·코드 안내·버튼은320에서도 잘리지 않는다. 오류는 상단 링크와 폼 근처 전체 문구로 찾을 수 있다. none은 가짜 버튼 없이 안내·다시 불러오기 링크를 제공한다. |
| 로그인된 공개 화면 | PASS(DEV 원 이미지 직접 판정). hero의 '대화로 이동'과 공개 안내를 확인했다. 전체 개인정보 비노출/실제 인증은 tester의 기능 검사 범위다. |
| overflow·44px | PASS(측정 범위). 자체 hero3폭와 열린320 및 상태 화면에서 가로 넘침이 없다. 열린320의 보이는 동작 최소 높이는43.988964 CSS px로44px 반올림 오차0.02px 이내다. [layout-checks](layout-checks.json), [details-checks](details-checks.json), [states-checks](states-checks.json). |
| 대비·focus | 대비 PASS(동일 토큰 DEV 계산 재사용). 본문14.70, 보조6.19, 링크/주 버튼7.96, 오류8.01, 선택면12.98, 경계3.79다. [DEV 실제 색 계산](../SAR-LANDING-GUIDE-001-DEV-test/browser-audit/contrast.json). focus3px/offset2px 규칙과 DEV 오류 이미지 outline은 확인했으나 실제 키보드 탐색 PASS는 아니다. |
| 회원 대표 화면 | 배치 보존 PASS, 기존 대비 FAIL. 대화/친구/내 에이전트 이동·대화 두 열·계정 입력·서비스 안내는 공개 hero 스타일의 영향을 받지 않는다. 빈 대화의 주 CTA 대비는 V-LG-02를 참조한다. |

본문은 keep-all로 한글 단어를 유지하며 긴 URL/경로만 강제로 감싼다. 320의 Command 끝 `node`가 줄 경계에서 갈리는 것은 값 손실이나 overflow가 아니며 D04가 허용한 줄바꿈 범위다. 줄바꿈이 설치 명령의 원 문자열을 바꾸지 않았음을 확인했다.

## 구체적 발견

**V-LG-01 / low / 비차단 / DEV 후속:** 예시 두 블록 다음 '대화에서 알아둘 점' 제목이 마지막 예시 경계선 바로 아래에 붙는다. 1440·390·320 DEV 예시 이미지와 자체 none390 이미지에서 보인다. 열린320 DOM에서 `.examples` 아래와 `#examples > h3` 위의 간격은 반올림 기준0px다. 같은 예시 안의 후속 문구처럼 보일 수 있지만 제목·본문은 읽히고 가려지지 않는다. D04 기존 간격 토큰을 적용해 예시와 보존 안내가 독립 구획임을 보이게 하면 개선된다. 구현 방법은 DEV 책임이며 이번에 정본이나 코드를 바꾸지 않았다. 수정하면 해당 예시 경계만 재검수한다.

**V-LG-02 / medium / 기존 회원 대비 FAIL / DEV 후속:** agent가 없는 합성 회원의 `/home`에서 '에이전트 추가' CTA 글자가 청록 배경에 어둡게 표시된다. [직접 확인한 회원 화면](tester-member-home-1440.png)과 [tester computed 대비](tester-contrast-results-snapshot.json)의 `.conversation-list .primary`는 foreground `rgb(24, 43, 49)`, background `rgb(0, 90, 98)`, 16px/600, 대비 **1.846967:1**이다. D04 본문 최소4.5:1에 미달한다. 위치는 desktop1440의 왼쪽 대화 목록 빈 상태, 약 x240–487/y301–351이다. 첫 agent 연결 행동을 읽기 어렵게 한다.

공유 CSS의 `.conversation-list a` 색상 규칙이 `.primary`보다 우선하는 원천을 확인했다. 기술 수정 방법은 DEV가 결정한다. [member-style-identity](member-style-identity.json)는 공개 변경 전 기준eae9a8b와 fixture4795의 memberStyle SHA256이 같음을 증명한다. 이 문제는 이번 공개 CSS의 새 회귀로 분류하지 않지만 기존 문제라는 이유로 대비 FAIL을 PASS로 바꾸지 않는다. 공개 landing CTA는 흰 글자/청록 배경7.958:1로 별도 PASS다. 회원 CTA 수락은 DEV 수정 후 해당 빈 상태의 실제 대비·대표 화면을 재검수해야 한다.

## 별도 미수락과 보존

- designer native keyboard: **not_run/unavailable**. 자체 page도 `document.hasFocus=false`다. DEV의 실제 Enter/Space 실패와 성공 응답만 반환한 한계를 보존한다. 이번 DOM summary/anchor click을 native 마우스·키보드 성공으로 표시하지 않는다. 후속 tester는 별도 Chrome CDP에서22개 browser checks를 통과했고 `browser-email-enter`와 실제 verify Enter 실행을 기록했다. [제공 결과 사본](tester-browser-results-snapshot.json)을 연결하며 designer가 다시 실행한 것으로 쓰지 않는다. 이22개가 모두 키보드 assertion인 것은 아니다. Tab/Shift+Tab·앵커·summary Enter/Space의 최종 범위는 tester 보고를 따른다.
- native browser200%: 보고 시점 **not_run/후속 대기**. DEV CSS `body.zoom=2` 이미지는 별도 관찰일 뿐이다. tester가 실제 지원 실행기로 완료 증거를 남겨야 한다. 이 제한과 회원 대비 FAIL이 남아 있어 전체 LG-08 수락을 열지 않는다.
- email 형식/발송/제한 오류, 실제 인증·연결 복귀·CSRF/CSP와 데이터 비노출은 독립 tester 담당이다. 이 보고는 사진만으로 그 동작을 통과시키지 않는다.
- 실계정·실제 Google/메일·운영 배포·Grok/Dots 설치/Events·장시간 수신·다른 엔진·스크린리더는 실행하지 않았다. 과거 운영/인증/시각 FAIL과 not_run을 이번 좁은 PASS로 대체하지 않는다.
- 새 제품 코드·설정·검사 코드 파일을 만들거나 수정하지 않았다. 제품 gate는 입력 동일성 확인 후 [DEV exact 후보 gate](../SAR-LANDING-GUIDE-001-CODE-REVIEW-review/author-final-handoff.md)의 make lint/test 종료0을 재사용한다. 이번 Windows에서 새 제품 gate PASS를 주장하지 않는다. 별도 디자인 lint/다중테마는 project.md상 미구성이다.

공유 fixture와 forward를 유지했다. 유지·종료 책임은 coor/tester다. `POST /fixture/stop`, 프로세스 종료, DB stop, 파일 삭제·이동을 실행하지 않았다. 모든 캡처·실패·탐색 원본은 보존한다. 최종 기록 SHA와 완료 인계는 [phase](../../../exec-plans/phases/SAR-LANDING-GUIDE-001-VISUAL.md)와 worker_done에 연결한다.

문서 검증은 [validation.json](validation.json)에 있다. 공식 packet 완료63쌍, deliverables strict13종 문제0/경고0, 변경 서술 문서의 로컬 링크와 `git diff --check`를 확인했다. 이 기록 검사는 시각/접근성 기준 FAIL을 없애지 않는다.


## CSS 수정 후속 — 2026-10-15

원 V-LG-01/02의 당시 실패와 원 PNG는 위 본문 그대로 보존한다. source `f2b8cb622dc397034423582978e77560e0355771` / DEV final `fc4ad58b65471e2de298cb9f1aba5d2253ea8135`의 PNG 7장을 새 designer 세션에서 직접 열었다. 두 항목은 수정 후보에서 resolved / 시각 PASS다. [후속 상세 보고](css-followup/report.md)는 실제 읽힘·구획 분리와 DEV 계산치 24px·7.958:1을 구분하고 source/PNG 동일성, 캡처 한계, 기존 tester native keyboard/200% 근거, 미검사 범위를 기록한다. 완성 제품 수락과 main 통합은 coor가 독립 코드 리뷰·tester 결과와 취합한다.
