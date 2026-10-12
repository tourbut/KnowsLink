---
title: 회원 웹 UI 검토·최소 수정 DEV 검증 보고
status: draft
updated: 2026-10-12
owner: dev
tasks: [SAR-FRONTEND-POLISH-001-DEV]
summary: "합성 fixture 기반 before/after 측정·캡처·키보드 검사와 FullOps lint 결과, 비평과 미실행 사유"
---

# SAR-FRONTEND-POLISH-001-DEV 검증 보고

DEV 자체 검증이다. 시각 수락·독립 QA가 아니다. 계획과 결함 판단은 [실행 기록](../../../exec-plans/phases/SAR-FRONTEND-POLISH-001-DEV.md)에 있다.

## 자료와 환경

- 모든 화면은 `TestMemberUIFixtures`의 합성 fixture HTML이다(`KNOWSLINK_UI_FIXTURE_DIR`). 실제 계정·키·관계·메시지·운영 서버는 쓰지 않았다. 비밀값을 출력하거나 캡처하지 않았다.
- before는 기준 코드 `7afaf02`, after는 구현 `f7b4786`의 렌더다. 원본 HTML은 레포 밖 임시 폴더에 두었다.
- 브라우저: Orca 내장 브라우저(Chromium 계열), 이 worktree에 새로 만든 탭 하나.

## 명령과 종료코드

| 명령 | 종료코드 | 결과 |
|---|---|---|
| `python3 -I measure.py <before> before` | 0 | 87개 키, 링크 198, PNG 27, 위반 18 |
| `python3 -I measure.py <after> after` | 0 | 87개 키, 링크 204, PNG 26, 위반 6 |
| `python3 -I keyboard.py <after> after/keyboard.json` | 0 | failures 0 |
| `impeccable detect --json home.html remote-consent.html verify.html` | 2(발견 있음) | side-tab 3, em-dash advisory 1. 정본 D04 8.5에 따라 유지 |
| `go test ./internal/relay -count=1` (Windows) | 0 | ok |
| `make lint` (Windows) | 2 | `cloudflared: command not found`(기존 환경 한계). 통과 근거로 쓰지 않음 |
| `lint.py --repo . --from 5c6485b… --out <레포 밖>` (G1, HEAD `5c6ca7b`) | 0 | ERROR 0, WARNING 1, product-lint·product-test passed. [gate/lint-g1.out](gate/lint-g1.out) |
| Linux 격리 clone `make lint`, `make test` (같은 HEAD, clean 확인) | 0, 0 | [gate/linux-make-lint.log](gate/linux-make-lint.log), [gate/linux-make-test.log](gate/linux-make-test.log) |

측정 스크립트는 [measure.py](measure.py)(이전 FIX-VISUAL 스크립트를 이번 폭·검사로 수정), [keyboard.py](keyboard.py), SSH shim 사본은 [gate/remote.py](gate/remote.py)다. 로그는 출력 전 접속값을 가린다.

## before → after

| 항목 | before | after |
|---|---|---|
| OAuth 동의·로그인 화면 `main`·건너뛰기 | 6건 누락(2화면 × 3폭) | 0 |
| `home-empty` 같은 이름 버튼 "로그인 다시 확인" | 3폭 모두 | 0 |
| 에이전트 행 펼침 표시 단독 줄 | before 측정식 오류로 수치 없음(패딩 미보정). 원 보고 V-04 PNG가 근거 | 0(보정한 측정식) |
| 빈 select 관계 초대·요청 상태 form | [home-empty-390 s3](before/png/home-empty-390-s3of5.png) | [home-empty-390 s2](after/png/home-empty-390-s2of4.png)에 없음 |
| 확인 코드 입력 | [verify-390](before/png/verify-390.png) 240px 일반 입력 | [verify-390](after/png/verify-390.png) 56px·28px 모노 |
| 원격 동의 주 동작 | [remote-consent-390](before/png/remote-consent-390.png) 승인=취소 위계 | [remote-consent-390](after/png/remote-consent-390.png) 승인 주 버튼, header·영역 |
| 연결 기한 문구 | "·"가 줄 앞 | [connection-token-390 s1](after/png/connection-token-390-s1of2.png) "(최대 10분·1회 사용)" 한 덩어리 |
| 가로 overflow, 44px 미만 | 0, 0 | 0, 0 |

남은 위반 6건은 `home`의 승인 요청 UUID 링크 줄 수다. range 사각형 계산이 하이픈 위치를 따로 세는 결과이며 [home-1440 s6](after/png/home-1440-s6of7.png)에서 한 줄이다. 이전 FIX-VISUAL 보고와 같은 비결함이다.

## 키보드·포커스

390×844 `home`: 받은 관계 초대 링크 → 승인 요청 링크 → 새 에이전트 만들기 → 첫 에이전트 `summary` 순서이고 모두 3px 실선 포커스다. Enter로 행이 열렸다. `remote-consent`: 승인 → 취소 → 내 agent 관리 순서, 3px 포커스. Orca 내장 브라우저는 페이지를 한 번 클릭해야 키를 전달하므로 `h1`을 클릭한 뒤 시작했다. 건너뛰기 링크 키보드 이동은 확인하지 않았다(존재·위치만 측정).

## 비평 (impeccable polish/audit, Operate 모드)

- 위계: 홈은 확인 필요 띠 → 에이전트 → 관계 순서가 유지되고, 빈 상태에서 쓸 수 없는 form이 사라져 다음 동작이 하나로 분명해졌다.
- 보안 승인 화면: OAuth 동의가 회원 화면과 같은 머리말·영역·주 버튼을 써서 승인 동작이 구분된다. 본문 문구는 바꾸지 않았다(영문 "agent" 혼용은 문구 과제 후보).
- 남은 시각 부채: 명령 `pre`의 `--confirmed`가 390px에서 `--`와 나뉜다(복사 값은 같다). V-02·V-03·V-08·V-12는 실행 기록의 사유대로 남긴다.

## 해당 없음

- 테마 전환: 밝은 테마 하나뿐이다(D04 1절, 다크 모드 후속). 전환 검사 해당 없음.
- 별도 디자인 lint: 구성되어 있지 않다(project.md). FullOps DESIGN 규칙은 Go 문자열 CSS를 검사하지 않는다. 대신 impeccable detect를 렌더된 HTML에 한 번 실행했다.
- 영상: 정지 화면으로 판정 가능하므로 만들지 않았다.
