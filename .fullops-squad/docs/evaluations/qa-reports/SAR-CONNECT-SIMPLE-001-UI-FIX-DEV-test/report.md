---
title: SAR-CONNECT-SIMPLE-001-UI-FIX-DEV 자체 검사 보고
status: draft
updated: 2026-10-27
owner: dev
tasks: [SAR-CONNECT-SIMPLE-001-UI-FIX-DEV]
summary: "ART·QA 후속 V-CS-01–06(EN Google 제목, 승인 재방문 문장, 시도·이름 버튼 aria-label, agent_rate 복귀, device 클라이언트 줄, 시도 최신순)을 고치고 Linux 격리 make lint·make test·integration EXIT=0으로 검사했다."
---

# SAR-CONNECT-SIMPLE-001-UI-FIX-DEV 자체 검사 보고

[ART 직접 시각 검수](../SAR-CONNECT-SIMPLE-001-ART-test/report.md)의 V-CS-01(medium)·V-CS-02–06(low)과 [QA](../SAR-CONNECT-SIMPLE-001-QA-test/report.md) 결함 1(V-CS-02와 같은 조사 문제)의 DEV 수정과 검사다. 테스트 레벨은 lite다. 운영 relay·운영 agent·실계정·실제 발송은 쓰지 않았다.

- 브랜치 `fullops/dev`, 기준 `e426f03be373`. 코드 커밋 `659380f1`(수정)·`9f279b4f`(Linux gofmt 정렬과 D04·D10 front matter).
- 원시 로그: `D:/workspace/KnowsLink/.devtools/connect-simple-001-uifix-dev/`(gitignore 공용 폴더). 파일 이름은 `<SHA>-lint.log`·`-test.log`·`-integration.log`·`fullops-lint.log`다.

## 결함별 수정과 검사

| ID | 수정 | 검사 | 결과 |
|---|---|---|---|
| V-CS-01 | `i18n_en_member.go`에 `"Google 로그인": "Google sign-in"` | `TestConnectFailureMessages`: EN refusal 렌더의 `<title>`·`<h1>`에 한글이 없고 'Google sign-in'이 있다 | 통과 |
| V-CS-02 | `member_pages.go` device `approved` KO 본문을 designer 결정 문장("…'내 에이전트'에서 `<이름>`의 상태가 '키 등록됨'으로 바뀝니다…")으로 바꿨다. EN은 그대로이고 사전 키만 새 KO 문장으로 바꿨다. D04 표의 같은 칸도 같은 문장이다 | `TestDeviceConsentTemplate`(approved 본문), i18n 번역 검사 | 통과 |
| V-CS-03 | `human_pages.go` 두 절 버튼에 `aria-label` '이 시도 정리: `<클라이언트>` · `<HH:MM>` 시작', '이름 저장: 이름 없는 에이전트 · `<구분값>`'. EN 'Remove this attempt: … · started …', 'Save name: Unnamed agent · …'. 보이는 문구는 그대로다 | `TestConnectScreens`(KO aria-label), i18n 번역 검사(EN 누락 0) | 통과 |
| V-CS-04 | `device.go` `devicePage`: `refused` 분기만 복귀를 `/home/agents`로 보낸다(`refusedTo`). 다른 거절은 현행 `/home`이다 | integration `TestGoogleDeviceHTTP`: `agent_rate` 거절 재방문이 429, `href="/home/agents"`, '내 에이전트로 돌아가기' | 통과 |
| V-CS-05 | `member_agents.go` `clientName(host)`(`hostNames`, 모르면 'Node 로컬 클라이언트')을 `ownAttempts` 제목과 device 화면 클라이언트 줄(로그인 전·동의 details)이 함께 쓴다 | `TestDeviceConsentTemplate`(`<dd>Claude Code</dd>`), `TestMemberUIFixtures` | 통과 |
| V-CS-06 | `ownAttempts`를 시작 시각 최신순, 같으면 ID순으로 정렬한다. 모든 시도의 TTL이 같아 `Exp` 순서가 시작 순서다 | `TestAttemptsAreNotAgents`: 세 시도의 순서와 만료 시도가 마지막 | 통과 |

D10 `module-design.md`의 A-05·M, W-01–W-04, W-02 행에 위 책임과 테스트를 더했다.

## 필수 검사

| 검사 | 환경 | 결과 |
|---|---|---|
| `make lint` | Linux 격리 clone, `9f279b4f` | EXIT=0 |
| `make test` | 같은 clone | EXIT=0 |
| `go test -race -tags=integration ./internal/relay` | Linux 격리 clone `659380f1`, 격리 Postgres 17 컨테이너(검사 뒤 정지) | EXIT=0, `ok … 161.809s` |
| `go test ./internal/relay` | Windows | 통과 |
| FullOps `lint.py --from e426f03be373` | 깨끗한 detached worktree | ERROR 0, WARNING 3(SIZE-001 `PLANS.md`·`i18n_en_member.go`·`member_agents.go` 기존 상한 초과), 실행 불가 2(`make`가 Windows에 없음. 위 Linux 결과로 대신한다) |

- `659380f1`의 Linux `make lint`는 내가 더한 EN 사전 세 줄의 정렬(gofmt) 때문에 EXIT=2였다. Windows gofmt는 기존 파일에도 다른 정렬을 요구해 판정에 쓰지 않았다. Linux gofmt 출력대로 고친 것이 `9f279b4f`다. `659380f1`→`9f279b4f` 차이는 사전 공백과 문서 front matter뿐이라 integration은 `659380f1` 결과를 쓴다.
- dev 작업 트리에 다른 과제(SAR-MOTION-INTRO)의 미추적 파일이 있다. 그래서 FullOps lint를 깨끗한 worktree에서 돌렸다. 그 파일은 옮기거나 지우지 않았다.

## 미실행과 후속

- designer EN 1폭 재검수(V-CS-01 `google-*-en`)는 coor가 이 SHA 뒤 배정한다.
- 보고용 캡처는 만들지 않았다(테스트 레벨 lite).
- 이전 CONNECT 시나리오 기대 갱신은 tester 몫이다(지시서 범위 밖).
- 정리하지 않은 임시 자료(AGENTS.md가 worker의 삭제·이동 정리를 막는다): FullOps lint용 detached worktree `D:/workspace/KnowsLink/.devtools/connect-simple-001-uifix-dev/wt`, `deliverables.py --stamp --path`에 `.fullops-squad/`를 붙여 잘못 생긴 보고서 사본 `.fullops-squad/.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-UI-FIX-DEV-test/report.md`(git exclude 대상), Linux 격리 clone `/tmp/knowslink-uifix-*`. Postgres 컨테이너는 검사 뒤 정지했다.
