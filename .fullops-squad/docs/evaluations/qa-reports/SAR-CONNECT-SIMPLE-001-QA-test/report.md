---
title: SAR-CONNECT-SIMPLE-001-QA 독립 QA 보고
status: draft
updated: 2026-10-27
owner: tester
tasks: [SAR-CONNECT-SIMPLE-001-QA]
summary: "에이전트 한도 제거·연결 시도 분리·기본 이름·재시도 폴더(C:·D:)·실패 사유·링크 재방문·Google 실패 구분·Grok release 기본값을 격리 relay와 실제 클라이언트로 독립 확인했다. 결함 critical/high/medium 0건, low 1건. 미확인은 시계 주입·Dots·실제 Google·실제 Linux 설치기."
---

# SAR-CONNECT-SIMPLE-001-QA 독립 QA 보고

1단계(격리 fixture) 결과다. 2단계 실연결(CS-AC-18)은 이 배정 범위 밖이고 실행하지 않았다. 시나리오는 [SAR-CONNECT-SIMPLE-001-TESTER](../../scenarios/SAR-CONNECT-SIMPLE-001-TESTER.md)다. 테스트 레벨은 lite다.

- 검증 브랜치 `fullops/tester` `5c750943`. 제품 코드는 main `3f1da403086e4449a7b3df214c441cd1b5972db6`와 같다(`git diff 3f1da403 HEAD -- internal adapters scripts cmd db` 비어 있음).
- 운영 relay·운영 agent·실계정·실제 발송·사용자 `~/.claude`·`~/.codex`는 쓰지 않았다. 제품 코드는 고치지 않았다.
- 환경: Windows 11, Node 22.22.2(`.devtools`), Go 1.27.1, 격리 `postgres:17`(검사 뒤 정지), Chromium 계열 ego 브라우저(화면 일부).
- 원시 로그는 [logs/](logs/), 실행 코드는 [harness/](harness/), fixture 소스는 [connect_simple_qa_fixture_test.go](connect_simple_qa_fixture_test.go)다. 실행은 이 파일을 `git archive` 사본의 `internal/relay/`에 복사해 `go test -tags "integration connectsimpleqa"`로 했다. 작업 트리의 제품 폴더는 건드리지 않았다.
- 최종 판정 표는 [results-consolidated.txt](logs/results-consolidated.txt)다(95 PASS, 1 not_run, 하네스 결함으로 대체된 4행). 여러 번 나눠 돌린 원본은 [results-all-runs.log](logs/results-all-runs.log)다.

## 결과 요약

| ID | 판정 | 근거 |
|---|---|---|
| CS-AC-01 | PASS | A가 웹에서 활성 agent 8개 생성, 모두 303. KO·EN 화면 7쪽 × 금지 문구 11개 검색 0건(`s_limits.py`). |
| CS-AC-02 | PASS | B가 10개 생성·철회, 철회 기록 10개 보존, 새 agent 생성 303. B 화면 한도 문구 0건. |
| CS-AC-03 | not_run | 전체 201개·임계 주입은 fixture로 불가. DEV `TestAgentAndPairCapacity` 근거를 재사용한다(재실행 안 함). |
| CS-AC-04 | PASS(보조) | B가 24시간 안 20개를 채우자 21번째 웹 생성은 429 `agent_rate` 문구와 시각, 연결 동의도 `agent_rate`+`retry_at`으로 거절됐다. 기존 20개(활성 10·철회 10)는 같다. 키 저장 전 시도 비집계와 시계 주입은 DEV 근거. |
| CS-AC-05 | PASS | C가 취소·relay 만료·앱 종료·로그인만·동의 뒤 미저장 5종을 만든 뒤 '내 에이전트'와 서버 agent 0, 친구 E가 보는 C의 agent 0. 이어서 정상 연결 성공(이름 `Claude Code`), 양쪽 1. '진행 중인 연결' 절에 시도가 보인다. |
| CS-AC-06 | PASS | '이 시도 정리' 즉시 삭제(`before=True after=False`), 정리 뒤 앱 폴링 거절. 동의 뒤 저장 유실 시도는 기한을 2시간 넘긴 뒤에도 정상 완료(키 1, 선택한 이름). 기한 +23시간 유지, +25시간 정리, 정리 뒤 완료 불가. 키 있는 agent 불변. |
| CS-AC-07 | PASS(Dots 제외) | 호스트 `claude` 기본 `Claude Code`, 두 번째 `Claude Code 2`, `codex` `Codex`, 호스트 없음 `내 에이전트 1`, 웹 추가 칸 `내 에이전트 N`(만든 뒤 N+1). 브라우저 흐름(로그인한 링크 [shot-02](shot-02-connect-link-signed-in-before-google.png) → Google 대역 계정 선택 [shot-03](shot-03-fixture-google-account-chooser.png) → 동의 화면)도 같은 기본 이름이다([shot-04](shot-04-consent-default-name-claude-code.png)). 동의 화면 지문은 클라이언트 출력 지문과 같았다. Dots 원격은 not_run. |
| CS-AC-08 | PASS | `#unnamed` 절이 '내 에이전트' 목록 앞에 있고 이름 입력·'이름 저장'이 `<details>` 밖에 보인다. 저장 즉시 이름이 정해지고 마지막 하나를 정하면 절이 사라진다. 실제 브라우저로도 확인([shot-01](shot-01-agents-unnamed-section.png)). |
| CS-AC-09 | PASS(부분) | Grok 패키지 `knowslink_status`·`knowslink_connect_status`에 `me` 이름이 있고 키·토큰은 없다. Claude Code·Codex·보낸·받은 출력은 DEV 스냅숏 근거. |
| CS-AC-10 | PASS | C:·D: 모두 남은 시도 파일 폴더(재시도 연결 성공, `.unused-<시각>`에 `login.json`·`private.pem` 보존), 대기 시도 폴더(같은 링크 재개, 새 start 0건), 미리 만든 빈 폴더(이동 없음), 새 폴더, 연결 완료 폴더(재사용)가 모두 성공. 권한은 소유자 `SHINPC\shin`, 상속 차단, 허용 `SHINPC\shin` FullControl 하나. 상위가 파일인 경로는 `private_folder`. |
| CS-AC-11 | PASS(`unconfirmed` 라이브 제외) | 사유 코드·문구·다음 행동을 KO·EN으로 확인: `expired`(실제 10분 4초 대기, KO), `cancelled`, `private_folder`(C:·D:), `unreachable`, `rate_limited`(시각 포함), `agent_rate`(시각·`retry_at`·'기존 에이전트'). 모든 실패 JSON에 `reason`이 있고 사유 없는 `failed`는 0건. `unconfirmed`는 서버에 저장 응답을 보내지 않는 주입에서 클라이언트가 일시 오류로 보고 다시 폴링해 라이브 출력을 못 만났다. 같은 폴더 재시도는 그 시도를 먼저 확인해 선택한 이름으로 연결됐다. `unconfirmed` 문구·`service_capacity`는 DEV 근거. `signed_out`은 웹 전용. |
| CS-AC-12 | PASS | 시작 회원 A: 완료 링크 '이 연결은 완료됐습니다'+이름([shot-06](shot-06-link-revisit-completed.png)), 동의 뒤 저장 전 '연결을 승인했습니다…승인만으로는 아직 연결되지 않았습니다'(D04 원문, [shot-05](shot-05-after-consent-not-yet-connected.png)), 취소·만료 '이 연결은 취소/만료됐습니다'+'시작한 앱에서 다시 연결'. 다른 회원 B와 로그아웃은 401 일반 문구이고 이름·상태 없음([shot-07](shot-07-link-revisit-logged-out.png)). |
| CS-AC-13 | PASS(대역) | 다른 계정(403 '다른 Google 계정으로 시작됐습니다', 계정 선택 폼), 새 가입 닫힘(403 '아직 KnowsLink 회원이 아니며 지금은 새 가입을 받지 않습니다'), 기한 지남('로그인 시간이 지났습니다. 시작한 앱에서…')이 서로 다르다. 응답에 이메일·이름·ID 없음. 실제 Google 화면은 아니다(Google 대역). |
| CS-AC-14 | PASS(Windows 재현) | 이 SHA 패키지 zip + 설치기 `defaults.json` 형태, 환경 변수 없음: status·connect_status·확인 없는/`false` connect에서 relay 요청 0건. `confirmed:true` 뒤 start 1건, 브라우저 동의 뒤 연결, 폴더 `<prefix>/connections/grok-agent`. 환경 변수 설정이 잘못된 defaults를 이긴다. defaults 없는 패키지는 `held`. 설치기 줄 68의 기본값은 소스로 확인. 실제 Linux Bot의 `install_bot_mcp.sh`는 실행하지 않았다. |
| CS-AC-15 | PASS | `python scripts/package_plugin.py --host grok --verify` 종료코드 0([package-grok.log](logs/package-grok.log)). DEV `--verify` EXIT=0은 같은 코드이므로 재사용했다. |
| CS-AC-16 | PASS | 디스크의 비밀 문자열 12–39개가 캡처·로그 23개에 0건, 패턴 검색 0건, A·B·D 화면에 비친구의 이메일·agent ID 0건. 친구 C의 agent는 친구 화면에만 나온다. |
| 이전 시나리오 회귀 | PASS(기대 변경 4건 별도) | [아래](#이전-시나리오-회귀)와 같다. |
| CS-AC-17 | 이 과제 아님 | designer ART. |
| CS-AC-18 | 이 과제 아님 | 2단계, 운영 배포 뒤 coor가 재배정. |

DEV의 `make lint`·`make test`·`go test -race -tags=integration`은 같은 제품 코드이고 새 실패를 만나지 않아 재실행하지 않았다(Windows에는 `make`가 없다). 직접 돌린 것은 위 표의 동작 검사와 `package_plugin.py --verify`다.

## 이전 시나리오 회귀

`SAR-UX-ISSUES-002-CONNECT-TESTER`의 `checks.py`(읽기)와 `mutate.py`를 새 fixture(`connect_browser_fixture_test.go`, 별도 DB)에서 다시 돌렸다([checks](logs/regress-ux002-checks.log), [mutate](logs/regress-ux002-mutate.log)). 92 PASS, 3 FAIL, mutate 1 FAIL. 재개·만료·취소·device 승인 뒤 미저장·remote 행 24개 항목은 모두 PASS다.

| 실패 | 원인 | 판정 |
|---|---|---|
| `choices_five_in_design_order_and_text`, `agents_c_zero_agents_first_agent_prompt_and_selector`, `guide_claude_adapters_link_and_session_on` | 선택지가 6개로 늘었고 Codex 항목과 Claude Code '자동 설정' 문구가 생겼다. SAR-AUTOSETUP-001 변경이다. | 이 과제와 무관한 기대 변경. 시나리오 기대를 갱신해야 한다. |
| `mutate.py sixth_agent_with_new_attempt_refused` | 한도 제거(L-01) 때문에 여섯 번째 agent가 303으로 만들어진다. | 의도된 변경. CS-AC-01 PASS가 같은 사실의 새 근거다. |

## 결함과 관찰

| # | 심각도 | 내용 | 재현 |
|---|---|---|---|
| 1 | low | 동의 뒤 저장 전 화면의 조사가 이름 끝소리에 맞지 않는다. D04 원문이 `<이름>이` 고정이라 이름이 모음으로 끝나면 '…에 저장 대기이 …'가 된다. 제품 규칙 변경이라 보고만 한다. | fixture에서 동의만 하고 저장하지 않은 시도의 이름을 `저장 대기`로 두고 시작 회원이 `/connect/<id>`를 연다([s_browser.py](harness/s_browser.py)). |
| — | 관찰 | relay에서 먼저 만료된 시도는 클라이언트 로컬 기한이 안 지났으면 `expired`가 아니라 `cancelled`로 보고된다(둘 다 relay 401). 실제 10분 대기에서는 `expired`가 나온다. 문구·다음 행동은 같다. | `s_attempts.py`의 relay 만료 항목. |
| — | 관찰 | 로그아웃 상태에서 완료된 링크를 열면 일반 문구 끝에 '아직 새 키가 연결되지 않았습니다'가 붙는다. D04·D02가 현행 문구 유지로 정했으므로 맞다. 완료 링크에는 어색하다. | [shot-07](shot-07-link-revisit-logged-out.png). |
| — | 관찰 | 이름 입력의 `maxlength`는 80인데 안내·서버 규칙은 1–40자다(41자 이상은 서버가 거절). 이전부터 그랬다. | 소스 `human.go validName`, `human_pages.go`. |

critical·high·medium 결함은 없다.

## 하네스 정정 기록

한 번에 통과하지 않은 실행을 숨기지 않는다. 아래는 모두 제품이 아니라 검사 코드의 문제였고 고쳐 다시 돌렸다. 대체된 4행은 results-consolidated.txt에 `SUPERSEDED(harness)`로 표시했다.

- 회원 승인 스레드가 429 대기 한도를 넘겨 '호스트 없음' 기본 이름 확인이 한 번 비었다 → 대기 한도를 늘리고 `s_nohost.py`로 다시 PASS.
- 첫 `CS-AC-08` 실행은 `<body` 태그가 없는 HTML을 잘라 판정이 틀렸다. 이후 같은 회원에 이름 없는 agent가 남아 있어 '절이 사라짐' 한 행이 FAIL로 남았고, 모두 이름을 정한 `s_names2.py`에서 PASS했다. 수동 확인에서도 저장은 즉시 됐다.
- W-01 '승인 뒤 저장 전' 기대를 D02 줄임말로 적어 FAIL이었다. D04 원문으로 맞춰 PASS. 만료 링크 항목도 같은 정정이다.
- `CS-AC-10` 첫 실행은 Windows PowerShell 5.1에서 `Get-Acl` 모듈 로드 실패와 한글 출력 인코딩 때문에 중단됐다 → `pwsh`로 바꿔 다시 실행. 같은 회원 A가 24시간 20개 한도(`agent_rate`)에 닿아 폴더 실행이 한 번 실패했다. 이것은 제품이 의도대로 동작한 것이다. 회원 D로 바꿔 전부 PASS.
- 주 fixture는 검사 중 여러 번 재시작했다(상태 초기화). `results-all-runs.log`가 모든 시도를 담는다.

## 미확인과 이유

- CS-AC-03·04의 시계·임계 주입, `service_capacity`: DEV 테스트 근거. fixture로 2000개 agent를 만들지 않았다.
- Dots 원격 기본 이름: fixture에 원격 OAuth 클라이언트가 없다.
- `unconfirmed` 라이브 출력(종료코드 3): 위 CS-AC-11. 동작 확인은 같은 폴더 재시도로 대신했다.
- 실제 Google 로그인·계정 선택 화면, 실제 Grok Bot 앱, Linux `install_bot_mcp.sh`, Claude Code·Codex 호스트 안의 도구 출력.
- 1440·390·320px·키보드: designer ART.

## 남은 일

- 사용자·coor: 결함 1(조사)을 고칠지 결정한다. D04 문구 표의 `<이름>이`를 이름 뒤 조사 처리 규칙으로 바꾸면 된다.
- designer/DEV: `SAR-UX-ISSUES-002-CONNECT-TESTER` 시나리오의 선택지 5개 기대와 한도 거절 기대를 현재 규칙으로 갱신한다.
- coor: 운영 배포 뒤 2단계(CS-AC-18)를 같은 과제로 재배정한다. relay를 클라이언트·설치기보다 먼저 배포한다(D12).
- D02·D04·D11 변경이 필요한 항목은 위 결함 1뿐이다.
