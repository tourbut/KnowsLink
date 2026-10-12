---
title: 회원 웹 UI 검토와 결함 최소 수정 기록
status: draft
updated: 2026-10-12
owner: dev
tasks: [SAR-FRONTEND-POLISH-001-DEV]
summary: 회원 웹 기존 화면을 합성 fixture로 검토하고 남은 low 결함(V-04·V-05·V-06·V-07·V-09·V-10·V-11·N-01)을 최소 수정한 계획·변경·검증·미검증
---

# SAR-FRONTEND-POLISH-001-DEV — 회원 웹 검토와 최소 수정

## 기준

- 지시서 `handovers/to_dev.md`, attempt `e5b3085323c84ed69a5d155b9a006ed6`, 기준 SHA `5c6485b4c7c1c573a14bf3e683701b2399284666`(작업 시작 HEAD `7afaf02`).
- 규칙 `fullops-common-0.3.3`, `project.md`. 디자인 정본 D04 `docs/design-docs/mockups/service-member-ui.md` 8절 토큰을 그대로 쓴다. 새 색·폰트·라이브러리·JS·기능은 없다.
- 테스트 레벨 lite. 하위 에이전트는 쓰지 않았다(파일 3개의 작은 직렬 수정).
- 호스트·모델·effort: claude-code / claude-opus-5-5 / medium. 주 스킬 impeccable(polish/audit). 보조 검토 없음.
- Context7 조회는 하지 않았다. 표준 CSS 속성과 Go `html/template`의 변수 재할당(`{{$x = ...}}`)만 썼다. 기존 테스트로 동작을 확인했다.

## 준비 확인

| 항목 | 결과 |
|---|---|
| `frontend.py --repo . --host claude-code --model claude-opus-5-5 --effort medium --primary impeccable --check` | exit 0, `status: disabled`. 레포 `docs/agents/frontend.json`의 `enabled: false` 때문이다. 설정 파일은 바꾸지 않았다. |
| 같은 모듈의 `selection(enabled=True, ...)`·`inspect` 메모리 실행 | `configuration_ready`, problems 0, 스킬 `~/.claude/skills/impeccable/SKILL.md` 확인. 인증·브라우저 QA·MCP는 `not_run`이며 이 결과는 실행 성공을 뜻하지 않는다. |
| `impeccable context --target internal/relay/member_pages.go` | exit 0. PRODUCT.md·DESIGN.md 없음. 기존 코드 범위 수정이므로 init 없이 D04와 기존 코드를 정본으로 진행했다. |
| impeccable craft-floor | UI 수정 직전에 읽었다. |

## 결함 목록과 판단

출발점은 [SAR-SERVICE-UI-001-VISUAL](../../evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md)의 남은 low V-02–V-12, [FIX-VISUAL](../../evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL/report.md)의 N-01, 이번 최초 검사(before)다.

| ID | 판단 | 처리 |
|---|---|---|
| V-04 320px에서 펼침 표시가 식별자 위 한 줄을 차지 | 실제 결함 | `summary h3`의 flex 기준폭 16em → 12em, `min-width:0` |
| V-05 버튼·기한 문구 끝 단어만 다음 줄 | 실제 결함 | 버튼 `text-wrap:balance`, "(최대 10분·1회 사용)"을 `nowrap` |
| V-06 선택할 에이전트가 없는데 관계 초대 form이 빈 select로 열림 | 실제 결함 | 연결 완료 에이전트가 있을 때만 form 표시. 관계가 있으면 기존 문구로 다음 동작 안내 |
| (신규) 에이전트 0개 홈의 요청 상태 form도 빈 select | V-06과 같은 결함 | 에이전트가 있을 때만 form 표시, 없으면 한 줄 안내 |
| V-07 fixture `home-empty`의 "활성 /5" | fixture 결함 | fixture에 `ActiveAgents: 0` |
| V-09 "로그인 다시 확인" 버튼 두 개가 같은 이름 | 실제 결함(A-09) | 에이전트 영역 버튼에 `aria-label="로그인 다시 확인 (키·연결 변경 전)"` |
| V-10 OAuth 동의·반환 화면에 header·main·건너뛰기 없음, 승인과 취소가 같은 위계 | 실제 결함(A-02, 보안 승인 화면) | 공통 header·`main.page`·`section`·건너뛰기 링크, 승인·Google 확인 버튼을 `primary` |
| V-11 확인 코드 입력이 일반 크기 | D04 5.2 미충족 | 높이 56px, 28px 모노 600, 자간 .3em |
| N-01 영역 마지막 링크 아래 여백이 위보다 큼 | 실제 결함 | `section>:last-child{margin-bottom:0}` |
| 브라우저 기본 표면 | craft-floor | `accent-color`, `::selection`을 기존 토큰으로 지정 |
| V-02 받은 초대 행의 "관계 철회" | 수정하지 않음 | `connections_test.go`가 pending(받은 초대 포함) 관계의 unpair 제공을 고정한다. 거절과 철회는 서버 동작이 다르므로 제품 판단이다. |
| V-03 만료 시각 짧은 날짜 | 수정하지 않음 | `clock()`이 여러 화면·정책 검사에서 공유된다. F-UI-04(KST)는 충족. 문구 형식 변경은 designer 판단이다. |
| V-08 gate·receipt의 영어 표현 | 범위 밖 | D04가 gate 문구 유지를 정했다. 별도 문구 과제. |
| V-12 960px 미만에서 섹션 이동이 확인 필요 띠 위 | 수정하지 않음 | 순서를 바꾸면 랜드마크·탭 순서가 함께 바뀐다. 390px 첫 화면에 띠가 보인다(기존 판정). |
| impeccable detect `side-tab`·`em-dash` | 수정하지 않음 | D04 8.5가 안내 띠의 6px 왼쪽 선과 철회 배지 "—"를 정본으로 정했다. 정본이 우선한다. |

## 변경

| 파일 | 변경 |
|---|---|
| `internal/relay/member_pages.go` | 위 CSS 6건, 기한 `nowrap`, 재확인 버튼 `aria-label`, 관계 초대·요청 상태 form 조건부 표시 |
| `internal/relay/remote_auth.go` | OAuth 동의·반환 화면 랜드마크·건너뛰기·주 버튼. form method·action·hidden·button 값은 같다 |
| `internal/relay/member_ui_test.go` | fixture `ActiveAgents: 0`, `TestMemberInviteFormNeedsConnectedAgent` |
| `internal/relay/testdata/member_forms.golden` | `home-empty`의 빈 form 2행 삭제. 다른 행은 같다 |
| `docs/design-docs/module-design.md`(D10) | golden 의도 차이 설명에 이번 2행 삭제를 추가 |

CSRF·owner·5분 재확인·CSP·서버 상태와 API·문구 의미는 바꾸지 않았다. 새 안내 문장 두 개("에이전트를 연결한 뒤 상대 식별자로 초대하세요."는 D04 7절 문구의 일부, "에이전트를 만든 뒤 요청 상태를 확인할 수 있습니다.")만 추가했다.

## 검증

증거는 [QA 보고서](../../evaluations/qa-reports/SAR-FRONTEND-POLISH-001-DEV-test/report.md)에 있다.

| 검사 | 대상 | 종료코드·결과 |
|---|---|---|
| `go test ./internal/relay -count=1` | Windows, 구현 `f7b4786` 작업 트리 | 0 |
| `TestMemberInviteFormNeedsConnectedAgent` | 같은 위치 | PASS. 수정 전 템플릿은 항상 `from-agent`를 렌더하므로 "미연결"·"철회" 경우에 실패한다 |
| Orca 내장 브라우저 측정(after) | 합성 fixture 29개 × 1440×900·390×844·320×640 | 위반 6건 모두 승인 요청 UUID 링크의 range 줄 수 계산이다(이전 보고서와 같은 비결함, 1440 PNG에서 한 줄). overflow 0, 44px 미만 0, 랜드마크 누락 0, 중복 버튼 이름 0, 펼침 표시 단독 줄 0 |
| 키보드 | 390×844 home, remote-consent | Tab 순서가 화면 순서, 모든 포커스에 3px 실선, Enter로 에이전트 행이 열림 |
| FullOps lint G1 | HEAD `5c6ca7b`, `--from 5c6485b` | exit 0, ERROR 0, WARNING 1(PLANS.md SIZE-001). product-lint `make lint`·product-test `make test`는 Linux 격리 clone에서 passed |

- Windows `make lint`는 `deploy/knowslink/beta.sh`의 `cloudflared` 부재로 실패한다(exit 2, 기존 한계). 그래서 이전 과제와 같은 task-local SSH shim(`make.cmd` → 고유 `/tmp` clone, HEAD·clean 확인)을 썼다. 서버 자원과 이전 clone은 지우지 않았다.
- 최종 HEAD의 lint 결과는 완료 보고와 worker_done에 적는다.

## 미검증과 후속

- designer 직접 시각 검수, tester 독립 동작 QA, 고정 SHA 독립 리뷰는 후속이다. 이 기록은 시각 수락이 아니다.
- 실제 Google 로그인·운영 화면·실제 스크린 리더·200% 확대·다른 브라우저 엔진은 확인하지 않았다.
- Postgres integration·browser 태그 테스트와 `make verify-mvp`는 실행하지 않았다. 템플릿·CSS 변경이며 서버 로직·SQL을 바꾸지 않았다.
- 키보드 검사는 페이지 클릭 뒤 `h1`부터 시작했다. 건너뛰기 링크의 키보드 동작은 존재·위치만 측정했다.
