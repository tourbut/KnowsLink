---
title: 회원 UI 변경의 독립 lite 동작 QA 보고
status: draft
updated: 2026-10-12
owner: tester
tasks: [SAR-FRONTEND-POLISH-001-TESTER]
summary: "고정 SHA 9597bad의 빈 form·OAuth·접근성·반응형 변경을 합성 fixture 29개 x 3폭으로 독립 검증한 결과와 미검사 범위"
---

# SAR-FRONTEND-POLISH-001-TESTER 독립 동작 QA

테스트 레벨 lite. 대상은 고정 SHA `9597bad351a2fbfbaa726444b320b8cf7000e1e2`의 제품 코드다. 리뷰 보고서는 [review/report.md](../SAR-FRONTEND-POLISH-001-TESTER-review/report.md)에 있다. 이 QA는 designer의 직접 시각 수락이 아니다.

## 환경·격리

- tester 체크아웃 `471e28a`는 9597bad와 `.fullops-squad` 밖 변경이 0건이다(`git diff --name-only 9597bad HEAD`로 확인).
- fixture는 `KNOWSLINK_UI_FIXTURE_DIR=<scratchpad>/fx go test ./internal/relay -run TestMemberUIFixtures`가 합성 데이터로 쓴 29개 HTML이다. `python3 -I -m http.server 18765 --bind 127.0.0.1`로 `http://127.0.0.1:18765`에서만 제공했다. 운영 origin·계정·키·관계·메시지·로그아웃·외부 발송은 쓰지 않았다.
- 브라우저: Orca 내장 브라우저(이 worktree의 새 탭 하나). 파일 삭제·이동은 없다.

## 명령과 종료코드

| 명령 | 위치 | 종료코드 | 결과 |
|---|---|---|---|
| `go test ./internal/relay -count=1 -run 'TestMember\|TestRemote' -v` | tester 체크아웃 | 0 | TestMemberPagesShowNextSteps, TestMemberUIFixtures, TestMemberInviteFormNeedsConnectedAgent, TestMemberReceiptView, TestRemote* 11개 PASS |
| `go test ./internal/relay -count=1` (전체 패키지) | 같은 SHA의 clean detached snapshot | 0 | ok |
| `go test ./internal/relay -count=1` (전체 패키지) | tester 체크아웃 | 1 | `TestFrozenParsingAndSigning` "registry integrity failure" 한 건. 환경 원인(아래) |
| `go vet ./internal/relay` | tester 체크아웃 | 0 | 출력 없음 |
| `python3 -I qa.py http://127.0.0.1:18765 result.json` | 레포 루트 | 0 | fixtures 29, viewports 3, failures 0 |

tester 체크아웃의 전체 실패 원인: `core.autocrlf=true`로 이 체크아웃의 작업 파일이 CRLF(`git ls-files --eol`이 `w/crlf`)이고 registry 무결성 해시가 달라진다. 변경 diff와 무관하다. 같은 SHA의 snapshot은 `w/lf`이며 같은 테스트가 통과한다. 제품 결함이 아니고 tester 환경 한계다. 이번 변경의 판정은 snapshot 결과를 따른다.

## 독립 QA 항목 (`qa.py`, `result.json`)

fixture 29개 전체를 1440×900·390×844·320×640에서 모든 `details`를 열고 측정했다(87개 조합).

| 항목 | 결과 |
|---|---|
| 가로 overflow | 0건. `scrollWidth`≤폭, 폭 밖 요소 0 |
| 44px 미만 조작 요소(skip 링크 제외) | 0건 |
| 랜드마크 | 모든 화면에 `main` 1·`h1` 1·본문 건너뛰기 링크 1. OAuth 동의·로그인(`remote-consent`, `remote-login`) 포함 |
| 같은 이름 버튼 중복 | 0건 |
| 빈 `select` | 0건 |
| `home-empty` | `/home/invites`·`/home/receipts` form 없음(3폭). 안내 "에이전트를 만든 뒤 요청 상태를 확인할 수 있습니다." 표시(390·1440) |
| `home`(연결 완료 에이전트 있음) | `post /home/invites`, `get /home/receipts` form 유지(3폭). 같은 SHA의 `TestMemberInviteFormNeedsConnectedAgent`가 연결 완료·미연결·철회 경계를 확인 |
| OAuth 승인 화면 | `remote-consent`의 primary는 "이 설치 연결 승인" 하나, "취소"는 보조. `remote-login`의 primary는 "Google 계정 확인" 하나(3폭) |
| 재확인 버튼 이름 | 홈의 `aria-label`은 "로그인 다시 확인 (키·연결 변경 전)"이며 보이는 문구 "로그인 다시 확인"으로 시작한다(WCAG 2.5.3 label in name). 같은 이름 중복 없음 |
| 확인 코드 입력 | 높이 56px, 28px 글꼴. `000000` 입력 시 390폭 `scrollWidth 250 ≤ clientWidth 250`, 320폭 `237 ≤ 237`로 6자리가 잘리지 않음 |
| 연결 기한 문구 "(최대 10분·1회 사용)" | 3폭 모두 한 줄 |
| 탭 순서(DOM) | `remote-consent`: 건너뛰기 → 이 설치 연결 승인 → 취소 → 내 agent 관리. `remote-login`: 건너뛰기 → Google 계정 확인 → 내 agent 관리. 건너뛰기 링크 `href="#main"`, `main#main` 존재 |

시각 확인은 390폭의 `remote-consent`·`verify`·`home-empty`를 직접 캡처해 읽었다. 승인 주 버튼이 분명하고 취소는 보조, 확인 코드 입력이 커졌으며 빈 상태 홈은 쓸 수 없는 form 없이 다음 동작이 "새 에이전트 만들기" 하나다. 캡처는 scratchpad의 판정용 임시 파일이며 레포에 저장하지 않았다.

## 새 결함

없다. critical·high·medium·low 0건.

## 미검사와 한계

- 실제 키 입력(Tab·Shift+Tab·Enter): Orca CLI `keypress`가 이 세션에서 포커스를 옮기지 못했다(클릭 뒤 `--page` 지정 포함 `activeElement`가 `body`). 첫 시도 기록은 [result-raw.json](result-raw.json)에 남겼다. 브라우저 뷰가 OS 포커스를 받지 않아 `:focus`·`:focus-visible` 스타일도 적용되지 않았다. 그래서 3px 포커스 링과 건너뛰기 링크가 포커스 때 보이는지는 검사하지 못했고, 탭 순서는 DOM 순서와 `focus()` 가능 여부로만 확인했다. 해당 CSS(`:focus-visible`, `.skip:focus`)는 이 변경에서 바뀌지 않았고 DEV 보고서에 390px 실제 Tab 기록이 있다(DEV 자체 검증, 독립 근거 아님).
- `remoteReturn`(앱 복귀 화면): fixture가 없어 렌더하지 않았다. 템플릿은 패키지 초기화의 `template.Must` 파싱과 위 테스트 통과로 구문만 확인했다.
- 실제 Google 로그인·운영 화면·실제 스크린 리더·200% 확대·다른 브라우저 엔진·Postgres integration·`make verify-mvp`는 실행하지 않았다.
- designer의 직접 시각 수락은 별도 배정이다.

## 재사용한 DEV 증거 (새 실행 아님)

| 증거 | 확인 |
|---|---|
| Git 공용 `worktrees/dev/fullops-gate/pass.json` | head 9597bad, base 5c6485b, attempt a6d3ba7956ff4f61a1a169ee0b828ac9, config_sha256 9c49bb2d… |
| `lint-attempt.json` | attempt a6d3ba79…로 pass.json과 일치 |
| lint 설정 동일성 | `git cat-file -p 9597bad:.fullops-squad/lint/lint.json`의 sha256이 9c49bb2d…로 pass.json과 같다. 로컬 파일은 CRLF 변환으로 해시가 다르다 |
| `%LOCALAPPDATA%/KnowsLinkDevTools/frontpolish001/9597bad…-lint.log`, `-test.log` | 둘 다 마지막 줄 `EXIT=0`. sha256 1fa38da2…, ce5e846b… |
| 이전 G1 `gate/lint-g1.json` | HEAD 5c6ca7b, ERROR 0, 두 product 명령 passed. 5c6ca7b..9597bad의 `.fullops-squad` 밖 변경 0건이라 제품 코드 동일 |

tester가 만든 `lint.json`(snapshot 9597bad, `--from 5c6485b`)은 ERROR 0, WARNING 1(PLANS.md SIZE-001, coor 관리 파일)이고 product-lint·product-test는 이 Windows 셸에 `make`가 없어 `unavailable`이다. 사유에 위 DEV 근거를 기록했다. 이는 tester의 독립 재실행이 아니다.
