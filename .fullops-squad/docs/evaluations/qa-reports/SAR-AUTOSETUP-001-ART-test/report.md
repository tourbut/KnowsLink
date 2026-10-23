---
title: SAR-AUTOSETUP-001-ART 회원 화면 자동 설정 안내 시각 검수
status: draft
updated: 2026-10-23
owner: designer
tasks: [SAR-AUTOSETUP-001-ART]
summary: "use=claude·use=codex·선택지 목록을 3폭×KO·EN으로 직접 검수했다. D04 카피는 48행 모두 일치하고 가로 넘침은 0이다. 320px에서 명령 입력란이 마지막 줄을 가린다(V-01, medium, dev 수정 필요)."
---

# SAR-AUTOSETUP-001-ART 시각 검수

검수 대상은 main `11d91648`(제품 코드는 DEV 병합 `8ec61ab3`와 같다). 제품 코드는 고치지 않았다. 환경: Windows 11, Chromium 1248 headless(CDP, DPR 1), 폭 1440×900·390×844·320×640. 격리 로컬 relay(`cmd/relay` 빌드본 `127.0.0.1:18431`)·격리 Postgres 17(`127.0.0.1:55432`)·loopback SMTP sink·합성 이메일 회원 1명(`art-<시각>@autosetup-art.invalid`). 운영 relay·실계정은 쓰지 않았다. 도구: [shoot.mjs](shoot.mjs), [smtp_sink.py](smtp_sink.py), [compare.mjs](compare.mjs). 측정 원본: [measure.json](measure.json). 검수 뒤 relay·DB 컨테이너(`kl-autosetup-art`)를 정리한다.

## 판정 (3화면 × 3폭 × 2언어)

공통 측정: 가로 스크롤 0(`scrollWidth == innerWidth`, 넘치는 요소 0), 명령 입력란에 개행 없음·`readonly`·`spellcheck=false`·전체 선택(Ctrl+A 동작과 같은 `select()`)이 값 전체를 선택, 도움말 대비 6.19:1(≥4.5), 명령·라벨 대비 14.7:1, EN 화면의 한글 잔존 0.

| 화면 | 1440 KO / EN | 390 KO / EN | 320 KO / EN |
|---|---|---|---|
| 선택지 목록 | PASS / PASS ([ko](png/list-1440-ko.png) [en](png/list-1440-en.png)) | PASS / PASS ([ko](png/list-390-ko.png) [en](png/list-390-en.png)) | PASS / PASS ([ko](png/list-320-ko.png) [en](png/list-320-en.png)) |
| `use=claude` | PASS / PASS ([ko](png/claude-1440-ko.png) [en](png/claude-1440-en.png)) | PASS* / PASS* ([ko](png/claude-390-ko.png) [en](png/claude-390-en.png)) | FAIL V-01 / FAIL V-01 ([ko](png/claude-320-ko.png) [en](png/claude-320-en.png)) |
| `use=codex` | PASS / PASS ([ko](png/codex-1440-ko.png) [en](png/codex-1440-en.png)) | PASS* / PASS* ([ko](png/codex-390-ko.png) [en](png/codex-390-en.png)) | FAIL V-01 / FAIL V-01 ([ko](png/codex-320-ko.png) [en](png/codex-320-en.png)) |

\* 390px는 명령이 3줄로 입력란(3행)에 딱 맞는다(`scrollHeight == clientHeight`). 여유가 없어 글꼴 확대 때 V-01과 같은 결함이 된다.

## D04 카피·순서 대조

`compare.mjs`가 D04 카피 표의 KO·EN 칸과 화면 행 텍스트를 정규화해 비교한다. claude·codex × KO·EN × 처음 연결·기존 연결 계속의 6개 행 48건 모두 일치한다(차이 0). 음성 대조(화면 문구를 일부러 바꾼 사본)는 차이로 잡혔다. 행 순서는 준비물·실행할 명령·자동으로 하는 일·직접 확인할 일·다음 행동·세션 연결이다. 선택지 6개의 라벨·보조 문구는 D04 선택지 표(KO·EN)와 같다. 6번은 'Hermes Agent, 그 밖의 호스트'만 있다. details 요약은 KO '수동 설치 문서', EN 'Manual install guide (in Korean)'이고 Codex 링크는 adapters Codex 절이다. 새 색·그림자·알약 버튼·복사 버튼·JS는 없다.

| AC-10 | 결과 | 근거 |
|---|---|---|
| 명령은 claude·codex에만 | PASS | web·grok·dots·other(6번)에는 `textarea`·버튼 없음(1440, KO·EN, [measure.json](measure.json)의 `extra-*`) |
| 문구가 D04와 같음 | PASS | 위 48건 |

## 차이 분류

| ID | 관찰 | 심각도 | 분류 |
|---|---|---|---|
| V-01 | 320px에서 명령이 4줄이 되는데 입력란은 3행(`rows=3`)이라 마지막 줄(`-App <호스트> -Lang <언어>`)이 입력란 안에 숨는다(`scrollHeight 106 > clientHeight 83`, KO·EN, claude·codex 4조합 [예](png/claude-320-ko.png)). 복사는 전체가 되지만 사용자는 명령 끝을 눈으로 확인할 수 없다. | medium | 수정 필요(dev). 기존 토큰 안에서: 입력란을 4행으로 늘리거나 좁은 폭에서만 `min-height`를 키운다. 새 스타일 요소 없음. |
| V-02 | 1440에서 '실행할 명령' 라벨(dt)이 입력란 라벨(dd 첫 줄)보다 약 16px 위에 놓이고, 한 줄 명령이 3행 높이 입력란에 빈 공간을 남긴다([예](png/claude-1440-ko.png)). | low | 수용 |
| V-03 | '안내 보기' 버튼과 안내 상자 제목 사이 간격이 거의 없다. 이번 변경 이전부터 있던 레이아웃이다. | low | 수용 |
| V-04 | 1440에서 `/knowslink-session on` 코드가 두 줄로 갈린다. | low | 수용 |

critical·high는 0건이다. 관찰하지 않은 상태의 결함 부재를 뜻하지 않는다.

## 못 한 것

실제 Ctrl+C 클립보드 내용, 200% 확대·키보드 포커스 순서·고대비 모드, `start=continue` 화면의 캡처(문구는 48건에 포함), 설치기 실행(DEV·tester 몫)은 확인하지 않았다. V-01 수정 뒤에는 320px 4조합만 다시 보면 된다.
