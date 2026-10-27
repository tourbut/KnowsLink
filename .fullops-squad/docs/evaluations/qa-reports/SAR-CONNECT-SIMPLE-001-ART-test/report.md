---
title: SAR-CONNECT-SIMPLE-001-ART 직접 시각 검수 (CS-AC-17)
status: review
updated: 2026-10-27
owner: designer
tasks: [SAR-CONNECT-SIMPLE-001-ART]
summary: "연결 단순화 회원 화면 KO·EN × 1440·390·320 84장 직접 검수. 가로 넘침 0·한도 문구 0·키보드 도달은 통과, EN Google 실패 제목 미번역 V-CS-01 medium으로 FAIL, low 5"
---

# SAR-CONNECT-SIMPLE-001-ART 직접 시각 검수 (CS-AC-17)

**판정: CS-AC-17은 FAIL이다(medium 1).** EN Google 실패 화면(W-02)의 제목이 한국어 'Google 로그인'으로 남는다(V-CS-01). 나머지 화면은 KO·EN × 1440·390·320에서 D04 문구와 배치가 같고, 가로 넘침이 0이며, 키보드로 모든 버튼에 닿는다. low 결함은 5건이다. critical·high는 없다. V-CS-01만 고치면 수락할 수 있다. 수정 확인은 DEV 테스트와 EN 1폭 재캡처로 충분하다.

- 대상: main `3f1da403086e4449a7b3df214c441cd1b5972db6`(DEV `0fff4ca0` + FIX-DEV `987152d4`). 지시서 커밋은 fullops/designer `5c750943`이다.
- 기준: [D02](../../../planning/product-specs/SAR-CONNECT-SIMPLE-001.md) M·W·E 표, [D04 연결 단순화 절](../../../design-docs/mockups/service-member-ui.md#연결-단순화와-내-에이전트-정리--sar-connect-simple-001)의 문구 표와 '자기 비평과 수락'.
- 제품 코드·D02·D04는 고치지 않았다. 운영 relay·운영 agent·실계정·실제 발송은 쓰지 않았다.

## 방법

- fixture: `3f1da403`을 `git archive`로 `.fullops-squad/local/SAR-CONNECT-SIMPLE-001-ART/src/`(gitignore)에 풀었다. 그 사본에만 fixture 테스트 `csart_fixture_test.go`(build tag `integration && csartfixture`)를 더했다. 이 파일은 커밋하지 않는다. 기존 `connect_browser_fixture_test.go`는 키 저장 전 agent를 만드는 옛 모델이라 새 '진행 중인 연결' 상태를 만들지 못한다.
- DB: 격리 Postgres 17 컨테이너 `kl-csart-pg`, `127.0.0.1:15541`. tester의 `kl-csqa-pg`(15531)와 포트·DB가 다르다. 스키마는 `db/migrations`의 goose Up 부분만 적용했다.
- 합성 회원: a(KO), e(EN)는 같은 상태를 가진다. b(KO), f(EN)는 다른 회원이다. 이메일은 `csart-*@example.test`다.
  - 키 있는 agent 4개: `Claude Code`, 40자 한글 이름, 이름 없음 2개.
  - 연결 시도: prepared(Grok Bot), 동의 뒤 키 미저장(Codex, agent 예약만), 만료(호스트 모름), 완료(Claude Code), 취소, `agent_rate` 거절.
  - 요청마다 새로 만드는 링크: 동의 화면용 prepared와 로그아웃 첫 화면용 requested.
- W-02: fixture 경로 `/fixture/google/{원인}`이 `googleCallback`과 같은 `connectFailure` + `render(403, "refusal")`를 부른다. OIDC 교환은 흉내 내지 않았다. '그 밖의 실패'는 현행 `googleFailure` 문구다.
- 도구: 일회성 스크립트 `raw/capture.mjs`(미커밋). headless Chrome, CDP 고정 폭(배율 1, 비모바일)으로 28화면 × 3폭 = 84장을 찍고 측정했다.
  - 측정값: 페이지 가로 넘침, 화면 밖 동작, 44px 미만 동작, 한도 문구('회원당', '5개', '먼저 철회', 'per member', 'revoke first'), role.
  - 1440px에서는 Tab 키 이벤트로 초점 순서와 초점 outline을 기록했다.
  - 결과: [captures/measure.json](captures/measure.json). 실행 종료코드 0, `views 84 overflow 0 limitText 0`.
  - 84장 전체와 본문 텍스트는 `.fullops-squad/local/SAR-CONNECT-SIMPLE-001-ART/captures/`에 있다. 통과 항목 캡처는 `raw/captures/`(미커밋)에 두었다. 결함 근거로 링크한 7장만 커밋했다.

## 공통 측정 결과

- 84장 모두 `scrollWidth`가 viewport 폭과 같다(가로 넘침 0). 화면 밖으로 나간 동작은 0이다.
- 44px 미만 동작은 '그 밖의 실패'(현행 첫 페이지)의 본문 속 약관 링크 2개뿐이다. 이 과제 범위 밖이며 기존 관찰(SAR-PUBLISH-001-VISUAL)과 같다.
- 한도 문구는 84장 모두 0건이다(M-03, L-06).
- 1440 초점 순서에서 모든 초점에 `3px solid` outline이 있다. 순서는 D04 수락 문장과 같다.
  - 내 에이전트: '안내 보기' → '이 시도 정리' ×3 → (이름 입력 → '이름 저장') ×2 → agent 카드.
  - 동의 화면: 이름 입력 → '지문이 같습니다. 연결합니다' → '취소' → details 요약 → '내 에이전트로 이동'. 지문은 정적 텍스트라 초점을 받지 않는다.
- 완료·동의 직후 재방문은 `p.note role=status`다. 만료·취소·거절은 `p.problem role=alert`다(D04와 같음).
- 테마: 제품에 테마 전환과 `prefers-color-scheme` 규칙이 없다. 라이트·다크 확인은 not_run이다.

## 화면 × 언어 × 폭 판정

폭 세 개(1440·390·320)는 측정을 모두 마쳤다. 표의 '캡처'에서 링크는 결함 근거로 커밋한 파일이다. 코드 표기는 `raw/captures/`의 미커밋 캡처다.

| 화면 | KO 1440·390·320 | EN 1440·390·320 | 캡처 | 근거 |
|---|---|---|---|---|
| 내 에이전트: 진행 중인 연결 | PASS ×3 | PASS ×3 | `agents-ko-1440`, [ko-390](captures/agents-ko-390.png), `agents-ko-320`, `agents-en-1440`, [en-390](captures/agents-en-390.png), `agents-en-320` | M-01 순서(선택지 → 진행 중인 연결 → 이름을 정할 에이전트 → 목록). 상태 문구 세 개, 정리 버튼, 정리 도움말이 D04 KO·EN 표와 같다. 제목은 `<클라이언트> · HH:MM 시작` / `· started HH:MM`이다. 호스트를 모르면 'Node 로컬 클라이언트'다. 시도 순서가 KO·EN에서 다르다(V-CS-06). |
| 내 에이전트: 이름을 정할 에이전트 | PASS ×3 | PASS ×3 | 위와 같음 | 접지 않은 이름 입력과 저장이 있다(A-05). 320에서 입력과 버튼이 위아래로 쌓이고 가려지지 않는다. KO 도움말은 D04와 같다. EN 도움말은 기존 번역 'Friends also see this name. Enter 1–40 characters.'로 D04 EN과 말만 다르다. 기존 문자열을 유지하기로 결정했다(남은 일). 버튼 접근성 이름이 겹친다(V-CS-03). |
| 연결 링크: 대기(로그인 전) | PASS ×3 | PASS ×3 | [ko-390](captures/login-ko-390.png), `login-en-320` | 현행 주의 문구, 지문, 'Google로 로그인하고 연결 계속'이 있다. 클라이언트 줄이 호스트와 상관없이 'Node 로컬 클라이언트'다(V-CS-05). |
| 연결 링크: 대기(동의 화면) | PASS ×3 | PASS ×3 | `consent-ko-1440`, `consent-ko-320`, `consent-en-1440`, `consent-en-390` | W-03 순서: 1. 지문 비교 → 2. 에이전트 이름(기본값 `Claude Code 2`, A-04 번호) → 연결 버튼 하나 → 취소 → 접힌 details. details 요약에 기한이 있다. 320에서 지문이 두 줄로 줄바꿈되고 잘리지 않는다. |
| 재방문: 동의 직후·키 저장 대기 | PASS ×3 (low) | PASS ×3 | [ko-390](captures/revisit-approved-ko-390.png), `revisit-approved-en-1440` | 제목·본문·링크가 D04와 같다. '연결됨'이라는 말이 없다(W-04). KO 본문의 'Codex이'는 조사 오류다(V-CS-02). |
| 재방문: 완료 | PASS ×3 | PASS ×3 | `revisit-consumed-ko-390`, `revisit-consumed-en-320` | '이 연결은 완료됐습니다'와 agent 이름. 오류 스타일이 아닌 `p.note`다(W-01, C-06 해소). |
| 재방문: 만료 | PASS ×3 | PASS ×3 | `revisit-expired-ko-390` | D04 문구와 같다. `p.problem`. |
| 재방문: 취소 | PASS ×3 | PASS ×3 | `revisit-cancelled-ko-390` | D04 문구와 같다. `p.problem`. |
| 재방문: `agent_rate` 거절 | PASS ×3 (low) | PASS ×3 (low) | [ko-390](captures/revisit-refused-ko-390.png) | D02 E 표 `agent_rate` 문구와 다시 가능한 시각이 있다. 복귀 링크가 '홈으로 돌아가기'다(V-CS-04). |
| 재방문: 다른 회원 | PASS ×3 | PASS ×3 | `revisit-other-ko-390` | 현행 401 일반 문구다. 상태와 agent 이름이 없다(W-01). |
| 재방문: 로그아웃 | PASS ×3 | PASS ×3 | `revisit-loggedout-en-390` | 위와 같다. |
| Google 실패: 다른 계정이 시작함 | PASS ×3 | **FAIL ×3** | `google-invalid_auth-ko-390`, [en-1440](captures/google-invalid_auth-en-1440.png) | 본문과 '다른 계정 선택'은 D04와 같다. 다른 계정의 정보가 없다. EN 제목·`<title>`이 'Google 로그인'이다(V-CS-01). |
| Google 실패: 새 가입 닫힘 | PASS ×3 | **FAIL ×3** | `google-capacity-ko-390`, [en-320](captures/google-capacity-en-320.png) | 같다. |
| Google 실패: 기한 지남 | PASS ×3 | **FAIL ×3** | `google-expired-ko-390` | 행동 버튼이 없다(D04 '없음'과 같음). EN 제목은 V-CS-01. |
| Google 실패: 그 밖의 실패 | PASS ×3 | PASS ×3 | measure.json `google-other-*` | 현행 첫 페이지 문구다(D04 '현행'). |

EN 로그아웃·Google 화면의 'This browser did not keep the language choice' 안내는 fixture가 쿠키 없이 `?lang=en`으로 연 결과다. 결함이 아니다.

## 결함과 수정 지시안

| ID | 심각도 | 위치 | 재현 | 수정 지시안 |
|---|---|---|---|---|
| V-CS-01 | medium | `internal/relay/google.go:210` view `Title: "Google 로그인"`. EN 사전에 이 키가 없다(`i18n_en_member.go`). `e7f4826f`에서 생겼다. | EN 브라우저로 연결 중 Google 실패 세 원인 중 하나를 본다. h1과 `<title>`이 'Google 로그인 · KnowsLink'다. | EN 사전에 `"Google 로그인": "Google sign-in"`을 더한다. EN refusal 렌더의 h1·title에 한글이 없음을 단언하는 테스트를 `TestConnectFailureMessages`에 더한다. 재검수는 EN 1폭 캡처 1장으로 충분하다. |
| V-CS-02 | low | `member_pages.go:75` device outcome `approved`, KO `{{name .Name}}이` | 'Codex', 'Grok Bot'처럼 받침 없는 이름으로 동의 직후 링크를 다시 연다. 'Codex이'가 된다. | 제품 결정: KO 본문을 "시작한 앱으로 돌아가세요. 앱이 키를 저장하면 '내 에이전트'에서 `<이름>`의 상태가 '키 등록됨'으로 바뀝니다. 승인만으로는 아직 연결되지 않았습니다."로 바꾼다. EN은 그대로다. D04 표도 같은 문장으로 바꾼다(남은 일). |
| V-CS-03 | low | `human_pages.go:32–33` 새 두 절의 버튼 | 진행 중인 시도가 셋이면 접근성 이름 '이 시도 정리'인 버튼이 세 개다. '이름 저장'도 겹친다. 기존 카드는 '에이전트 이름 저장: <이름> / <구분값>' 형식이다. | 기존 형식대로 `aria-label`을 단다. '이 시도 정리: `<클라이언트>` · `<HH:MM>` 시작', '이름 저장: 이름 없는 에이전트 · `<구분값>`'. EN도 같은 형식이다. 보이는 문구는 바꾸지 않는다. |
| V-CS-04 | low | `devicePage` `refused` 분기의 refusal `Back` | `agent_rate`로 거절된 연결 링크를 시작 회원이 다시 연다. 복귀가 '홈으로 돌아가기'다. | 키 연결 문맥이므로 '내 에이전트로 돌아가기'(`/home/agents`)로 바꾼다. 다른 재방문 화면과 같다. |
| V-CS-05 | low | `member_pages.go:77, 84` device 템플릿의 `<dd>Node 로컬 클라이언트</dd>` 고정값 | 호스트 `claude`의 링크를 연다. 로그인 전 화면과 동의 화면의 details가 'Node 로컬 클라이언트'다. '내 에이전트'의 같은 시도는 'Claude Code'로 보인다. | 클라이언트 줄을 `hostNames` 표시 이름(모르면 현행 문구)으로 보인다. '내 에이전트' 시도 제목과 같은 값이다. |
| V-CS-06 | low | [agents-ko-390](captures/agents-ko-390.png)·[agents-en-390](captures/agents-en-390.png), `member_agents.go:191` `ownAttempts`가 무작위 연결 ID 순서로 정렬한다 | 같은 시도 3개가 KO 회원과 EN 회원에서 다른 순서로 보인다. 새로 고쳐도 시각과 상관없는 순서다. | 시작 시각의 최신순으로 정렬한다. 같은 시각이면 ID 순서다. |

비차단 관찰:

- 시작 회원이 연 재방문 화면(완료·동의 직후·만료·취소)의 머리글에 언어 줄('한국어 · 설정에서 변경')이 없다. 동의 화면에는 있다. 읽기와 동작에는 영향이 없다.
- D04 시도 제목 예시는 구분값이 4자(`a1b2`)다. 실제 화면은 8자다. 기존 카드의 구분값과 같아 그대로 둔다.

## not_run

- 실제 Google OAuth 왕복. W-02는 같은 판정 함수와 렌더 호출로 그린 화면이다. 실제 Google 화면은 tester(CS-AC-13) 범위다.
- 상태를 바꾸는 POST(이 시도 정리, 이름 저장, 동의 제출, 취소)의 화면 전이. fixture가 이 POST를 막는다. 동작은 DEV 테스트와 tester QA 범위다.
- 브라우저 실제 200% 확대와 스크린리더. 320 CSS px 재배치 결과(넘침 0)를 확대 대체 근거로 쓴다.
- 라이트·다크 테마. 제품에 테마 전환이 없다.

## 남은 일

- DEV: V-CS-01(수락 차단)을 고친다. V-CS-02–06은 같은 수정에 묶거나 PLANS.md에 대기시킨다.
- designer 후속(D04): V-CS-02의 KO 문장을 반영한다. 이름 도움말 EN을 기존 문자열 'Friends also see this name. Enter 1–40 characters.'로 맞춘다.
- 재검수: V-CS-01 수정 뒤 `google-*-en` 1폭 캡처와 측정만 다시 한다. 이 보고서에 절을 덧붙인다.
- 정리: 컨테이너 `kl-csart-pg`와 fixture 프로세스는 보고 뒤 정지한다. `.fullops-squad/local/SAR-CONNECT-SIMPLE-001-ART/`는 gitignore 대상으로 보존한다.
