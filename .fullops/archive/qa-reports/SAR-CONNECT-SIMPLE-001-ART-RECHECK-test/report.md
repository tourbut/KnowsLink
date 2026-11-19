---
title: SAR-CONNECT-SIMPLE-001-ART-RECHECK V-CS-01–06 재검수 (CS-AC-17)
status: review
updated: 2026-10-27
owner: designer
tasks: [SAR-CONNECT-SIMPLE-001-ART-RECHECK]
summary: main 06af3fe2에서 V-CS-01(EN 1440 Google 제목)과 V-CS-02–06(KO·EN 390 문장·aria 이름·복귀 링크·클라이언트 줄·최신순)이 모두 고쳐져 CS-AC-17 PASS
---

# SAR-CONNECT-SIMPLE-001-ART-RECHECK V-CS-01–06 재검수 (CS-AC-17)

**판정: CS-AC-17은 PASS다.** [이전 ART 검수](../SAR-CONNECT-SIMPLE-001-ART-test/report.md)의 V-CS-01(medium)과 V-CS-02–06(low)이 main `06af3fe2`에서 모두 고쳐졌다. 남은 critical·high·medium은 없다. 이전 검수의 나머지 항목(84장 측정)은 이번 수정 범위 밖이라 반복하지 않았다.

- 대상: main `06af3fe2318d1f6cf1e423e1f5257ae656c51895`(UI-FIX-DEV `061c4142`). 수정 내용은 [UI-FIX-DEV 보고](../SAR-CONNECT-SIMPLE-001-UI-FIX-DEV-test/report.md)에 있다.
- 기준: D02 SAR-CONNECT-SIMPLE-001 CS-AC-17, D04 연결 단순화 절, 이전 ART 보고의 결함 표와 수정 지시안.
- 제품 코드·D02·D04는 고치지 않았다. 운영 relay·운영 agent·실계정·실제 발송은 쓰지 않았다.

## 방법

- `06af3fe2`를 `git archive`로 `.fullops-squad/local/SAR-CONNECT-SIMPLE-001-ART-RECHECK/src/`(gitignore)에 풀었다. 이전 ART의 미커밋 fixture `csart_fixture_test.go`를 그 사본에만 복사했다. `3f1da403`→`06af3fe2`에서 `db/migrations`가 바뀌지 않아 이전 스키마(goose Up)를 그대로 썼다.
- DB: 이전 격리 컨테이너 `kl-csart-pg`(`127.0.0.1:15541`)에 새 DB `csrecheck`를 만들었다. 합성 회원 a(KO)·e(EN)·b·f와 상태는 이전 ART와 같다.
- W-02 fixture 경로는 `googleCallback`(`google.go:209–215`)과 같은 view·`render(403, "refusal")`를 부른다. 이 코드는 이번 수정에서 바뀌지 않았다.
- 도구: 일회성 스크립트 `recheck.mjs`(미커밋, 위 local 폴더). headless Chrome, CDP 고정 폭(배율 1, 비모바일)으로 13화면을 찍었다. 버튼 이름은 CDP `Accessibility.getFullAXTree`의 role=button 노드 이름이다.
- 결과: [captures/measure.json](captures/measure.json). 실행 종료코드 0, `views 13 overflow 0 navError 0`.

## 결함별 재검수

| ID | 폭·언어 | 확인 | 캡처 | 판정 |
|---|---|---|---|---|
| V-CS-01 | EN 1440, 세 원인 | `invalid_auth`·`capacity`·`expired` 모두 h1 'Google sign-in', `<title>` 'Google sign-in · KnowsLink'다. 본문의 한글은 언어 선택 '한국어'뿐이다. | [google-invalid_auth-en-1440](captures/google-invalid_auth-en-1440.png) | PASS |
| V-CS-02 | KO·EN 390 | KO: "시작한 앱으로 돌아가세요. 앱이 키를 저장하면 '내 에이전트'에서 Codex의 상태가 '키 등록됨'으로 바뀝니다. 승인만으로는 아직 연결되지 않았습니다." designer 결정 문장과 같다. 'Codex이'가 없다. EN은 이전 문장 그대로다. | [revisit-approved-ko-390](captures/revisit-approved-ko-390.png) | PASS |
| V-CS-03 | KO·EN 390 | 접근성 트리 버튼 이름: '이 시도 정리: Grok Bot · 11:27 시작', '… Codex · 11:16 시작', '… Node 로컬 클라이언트 · 11:15 시작', '이름 저장: 이름 없는 에이전트 · 34373932'·'… 72415f35'. EN은 'Remove this attempt: … · started HH:MM', 'Save name: Unnamed agent · …'. 겹치는 이름이 없다. 보이는 문구는 '이 시도 정리'·'이름 저장' 그대로다. | [agents-ko-390](captures/agents-ko-390.png), [agents-en-390](captures/agents-en-390.png) | PASS |
| V-CS-04 | KO·EN 390 | `agent_rate` 거절 재방문의 복귀 링크가 '내 에이전트로 돌아가기' / 'Back to My agents' → `/home/agents`다. | [revisit-refused-ko-390](captures/revisit-refused-ko-390.png) | PASS |
| V-CS-05 | KO·EN 390 | 호스트 `claude` 링크의 클라이언트 줄이 로그인 전 화면과 동의 화면 details 모두 'Claude Code'다(EN 'Client: Claude Code'). '내 에이전트' 시도 제목의 표시 이름과 같은 `hostNames` 값이다. | [login-ko-390](captures/login-ko-390.png) | PASS |
| V-CS-06 | KO·EN 390 | '진행 중인 연결' 순서가 KO·EN 모두 Grok Bot 11:27 → Codex 11:16 → Node 로컬 클라이언트 11:15로 같다. 시작 시각 최신순이다. | [agents-ko-390](captures/agents-ko-390.png), [agents-en-390](captures/agents-en-390.png) | PASS |

13화면 모두 `scrollWidth`가 viewport 폭과 같다(가로 넘침 0).

## not_run

- 이전 ART와 같다. 실제 Google OAuth 왕복, 상태를 바꾸는 POST의 화면 전이, 실제 200% 확대·스크린리더, 라이트·다크 테마(제품에 없음)는 실행하지 않았다.
- 320·1440의 V-CS-02–06 재캡처. 수정이 문구·속성·정렬뿐이고 배치를 바꾸지 않아 지시서대로 390 한 폭만 봤다.

## 남은 일

- designer 후속(D04, 이전 ART 남은 일): 이름 도움말 EN을 기존 문자열 'Friends also see this name. Enter 1–40 characters.'로 맞춘다. V-CS-02 KO 문장은 UI-FIX-DEV가 D04에 반영했다.
- 정리: fixture 프로세스와 컨테이너 `kl-csart-pg`는 보고 뒤 정지한다. `.fullops-squad/local/SAR-CONNECT-SIMPLE-001-ART-RECHECK/`는 gitignore 대상으로 보존한다.
