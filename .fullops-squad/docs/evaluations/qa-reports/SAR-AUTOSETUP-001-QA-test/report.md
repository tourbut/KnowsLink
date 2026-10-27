---
title: SAR-AUTOSETUP-001-QA — Claude Code·Codex 자동 설정 큰 틀 검수 보고
status: draft
updated: 2026-10-28
owner: tester
tasks: [SAR-AUTOSETUP-001-QA, SAR-AUTOSETUP-001-QA2, SAR-AUTOSETUP-001-QA3]
summary: "QA3(10-28)에서 0.1.1이 호스트 플러그인 3개 환경의 0.1.0→0.1.1 갱신과 같은 명령 재실행을 통과해 O-04가 해소됐다(미해결 high 0). 1차(10-23)는 N1–N6·갱신·S7·자격 경계를 확인했다. 2차(10-27, QA2)는 로그인한 격리 Claude Code 실세션으로 N7–N8·AC-02·AC-04를 PASS 확인했다(3분 24초). 갱신(AC-07)은 호스트 플러그인이 2개 이상이면 실패하는 결함 O-04(high)를 찾았다. Codex는 사용자 결정 D040으로 제외."
---

# SAR-AUTOSETUP-001-QA — 큰 틀 검수 보고

독립 동작 QA(lite, 새 사용자 흐름 1회)다. 제품 코드는 고치지 않았다. 운영 relay·실제 Google 계정·실제 상대·사용자 실제 `~/.claude`·`~/.codex`는 쓰지 않았다.

**최신 결론(QA3, 2026-10-28): O-04 해소, 미해결 high 0.** 아래 `QA3 재검증` 절이 최신 판정이다. 이전 결론: **(QA2 반영, 2026-10-27): Claude 쪽 AC-01·02·04 PASS, AC-07(갱신)은 O-04 때문에 조건부.** 아래 `QA2` 절이 최신 판정이다. 1차(10-23) 절은 기록으로 남기고 해당 행에 QA2 판정을 덧붙였다. 미해결 high **1건(O-04)**, critical 0. Codex 항목은 사용자 결정 D040으로 이번에 제외했다.

## QA3 재검증 (SAR-AUTOSETUP-001-QA3, 2026-10-28)

독립 동작 QA(lite)다. 제품 코드는 고치지 않았다. 운영 relay·운영 agent·실제 상대·실제 `~/.claude`·`~/.codex`는 쓰지 않았다. 로그인 자격 파일은 열지 않았다(크기·수정시각만: 559바이트, 시작·끝 동일, [전](ev/qa3-cred-stat-before.txt)/[후](ev/qa3-cred-stat-after.txt)). 삭제·`git worktree remove`·이동 정리는 하지 않았다.

**결론: O-04 해소. 다중 플러그인 갱신·0.1.0→0.1.1 갱신·같은 명령 재실행·재시작 뒤 세션 연결 PASS, `[F-10]` 안내 PASS. 미해결 high 0, critical 0.** 새 관찰 O-05(low–medium), O-06(info)를 아래에 적었다.

### 환경

| 항목 | 값 |
|---|---|
| 제품 SHA | main `0349ca73`(0.1.1). 0.1.0 산출물은 `66407760`에서 빌드(Claude ZIP `2d2c1884…`, QA2의 0.1.0과 같은 해시). 둘 다 격리 `git worktree add --detach`에서 빌드, 제품 파일 변경 0 |
| 호스트 | Windows 11 Pro 26300, Windows PowerShell 5.1, Claude Code 2.1.311, Node 22.22.2 |
| 로그인 | QA2와 같은 사용자 로그인 격리 `CLAUDE_CONFIG_DIR`(`…utosetup-001-qa\claude-config`). 시작 전 호스트 CLI(`claude plugin uninstall`/`marketplace remove`)로 QA2의 knowslink 항목만 해제했다 |
| relay | 격리 Postgres 17 컨테이너(15483) + 시험 fixture(합성 회원 a·b, Google 콜백만 대행) 포트 18443. 한 relay가 `/install/0.1.0/`과 `/install/0.1.1/`을 함께 낸다 |
| 호스트 플러그인 | `knowslink`(설치기가 등록) + **실제 둘째 플러그인** `qa3-other@qa3-other-mkt`(격리 폴더 안에 만든 로컬 마켓플레이스, 호스트 CLI로 설치·해제) + Orca가 주입한 inline `orca-agent-view@inline` = `plugin list --json` 항목 **3개** |
| 도구 | [qa3-install_run.py](tools/qa3-install_run.py)(`QA_VERSION`으로 설치 버전 지정), [qa3-fixture.sh](tools/qa3-fixture.sh), 기존 `qa.py`·`snap.py`(포트만 18443) |

### 판정

| 항목 | 판정 | 근거 |
|---|---|---|
| 0.1.0 신규 설치(기준 상태) | **PASS**(항목 2개 호스트) | 플러그인 2개(inline+knowslink)에서 설치 완료, agent 1, 친구 1 |
| 갱신 0.1.0→0.1.1, 호스트 3개 항목(AC-07·S5) | **PASS** | 같은 화면 명령 URL만 0.1.1로: `작업: 기존 설치 갱신 (현재 0.1.0 → 새 0.1.1)` → `갱신 완료. 에이전트·키·친구는 그대로입니다.` 5.3초, 종료 0, 브라우저 단계 0([로그](ev/qa3-claude-update-010-to-011-multi.log.txt)). 전후 agent ID·`agent.json`·`private.pem` 포함 private 폴더 해시 `d24dfb72…` 동일·웹 agent 수 1·친구 수 1 동일([전](ev/qa3-snap-before-update.txt)/[후](ev/qa3-snap-after-update.txt)). 이전 설치 `.previous-*` 보존 |
| 0.1.1→0.1.1 같은 명령 재실행(3개 항목) | **PASS** | `갱신 완료`, 5.3초, 종료 0, 해시·agent 수·친구 수 동일([로그](ev/qa3-claude-rerun-011-multi.log.txt), [후](ev/qa3-snap-after-rerun.txt)) |
| 갱신 뒤 Claude Code 다시 열기 → `/knowslink-session on` | **PASS** | 새 `claude` 프로세스에서 `KnowsLink 연결됨: 세션 92366785 · 내 agent agent_745dc7… · 대기 0건`. 같은 agent. 한계: 이번에는 헤드리스(`claude -p`)로 열었다(QA2는 대화형 터미널). 발송 승인 대화상자는 다시 검증하지 않았다 |
| 예상 밖 오류 `[F-10]` | **PASS** | 합성 실패: `LOCALAPPDATA`를 비워 동의 전에 실패시켰다(3개 항목 호스트). `[F-10] 예상하지 못한 오류로 멈췄습니다: …` / `이미 바뀐 것: 없음` / `다음 행동: 같은 명령을 다시 실행하세요…`, 종료 1, 설치·키 변화 없음([로그](ev/qa3-claude-s6-f10-synthetic.log.txt)). 동의 뒤 단계의 F-10(`바뀐 것: 설치 기록 위치`)은 실행하지 못했다 |
| 자격 경계 | **PASS** | 자격 파일 크기·수정시각 동일. 설치기가 바꾼 호스트 설정은 `settings.json`의 knowslink 두 항목뿐(설치 로그의 계획 문구와 일치). 출력·ev에서 지문 전체·토큰·개인키 0건(지문 앞 8자, agent ID 앞 6자만 남김) |

### O-04 판정: **해소**

QA2 재현 조건(호스트 `plugin list --json` 항목이 2개 이상이고 knowslink가 그중 하나)을 더 강하게(항목 3개, 실제 둘째 플러그인 포함) 만들었는데 0.1.1 갱신·재실행이 통과했다. 같은 조건의 0.1.0 첫 설치는 아래 O-05처럼 실패했다([실패 로그](ev/qa3-claude-010-firstinstall-3plugins-fail.log.txt)). 즉 같은 환경에서 0.1.0은 실패하고 0.1.1은 통과한다.

### 새 관찰

| ID | 심각도 | 내용 | 재현·영향 |
|---|---|---|---|
| O-05 | low–medium(0.1.0 한정, 0.1.1에서 해소) | 공개 0.1.0은 호스트 항목이 3개일 때 **첫 설치도** 실패한다. 호스트 등록(`marketplace add`·`plugin install`)까지 마친 뒤 등록 확인에서 `[ERROR] "1"개의 인수가 있는 "GetFullPath"을(를) 호출하는 동안 예외가 발생했습니다`로 종료(종료 코드 1). 키·agent는 아직 만들기 전이라 데이터 손실은 없으나 호스트 설정에는 knowslink가 등록된 채 남는다. QA2는 항목 2개라 첫 설치가 통과했다 | 재현: 항목 3개 호스트에서 0.1.0 명령. 원인은 O-04와 같은 `ConvertFrom-Json` 배열 미펼침(등록 확인 쪽). 0.1.1은 같은 확인 단계를 3개 항목에서 통과했다(갱신 경로가 같은 확인을 지난다). 0.1.1 **신규** 설치는 이번에 따로 실행하지 않았다(아래 한계). 조치: 0.1.1 게시로 신규 사용자는 해소 |
| O-06 | info–low | 0.1.1 갱신 뒤에도 호스트의 플러그인 캐시(`plugins/cache/…/knowslink/0.1.0`)는 0.1.0 그대로이고 `plugin list`의 version도 0.1.0, `folderVersion`만 0.1.1이다. 설치기가 `plugin install`을 다시 실행하지만 호스트는 이미 설치된 것으로 보고 캐시를 갱신하지 않는다. MCP 서버는 설치 폴더의 절대 경로 `plugin.js`(0.1.1)로 실행되고, 훅·스킬은 캐시(0.1.0)에서 로드된다. 0.1.0과 0.1.1은 훅·스킬이 같아 동작 차이가 없다 | 재현: 갱신 뒤 `claude plugin list --json`의 `installPath`·`version` 확인. 이후 버전이 훅·스킬을 바꾸면 갱신해도 반영되지 않을 수 있다 → dev 판단(설치기가 갱신 때 `plugin update` 또는 제거 후 재설치 필요 여부) |

### QA3 한계

- Google 로그인·동의는 fixture 대행이다(1·2차와 같음). 사용자 실계정 비밀번호 입력 0.
- 0.1.1 신규 설치(항목 3개 호스트)는 직접 실행하지 않았다. 0.1.0→0.1.1 갱신이 같은 등록 확인 단계를 통과한 것으로 갈음했다.
- 세션 재연결은 헤드리스 `claude -p`로 확인했다. 대화형 화면·발송 승인·도착 알림은 QA2 결과에 의존한다.
- Codex는 D040으로 제외했다. 호스트 3개 항목은 inline 1 + 실제 로컬 마켓플레이스 플러그인 1 + knowslink다. 외부 마켓플레이스에서 받은 플러그인은 쓰지 않았다.
- 남긴 임시 경로(삭제 안 함): `D:\workspace\KnowsLink\.devtoolsutosetup-001-qa3\`(워크트리 `src-new`·`src-old`, 상자 `work\claude-A·claude-C`, `rel`, `other-mkt`, 로그), Docker 컨테이너 `kl-autosetup-qa3-pg`(실행 중), 포트 18443의 fixture 프로세스. 격리 로그인 폴더에는 knowslink 0.1.1과 `qa3-other` 플러그인이 설치된 채 남았다.

## QA2 실세션 검수 (SAR-AUTOSETUP-001-QA2, 2026-10-27)

독립 동작 QA(lite)다. 제품 코드는 고치지 않았다. 운영 relay·운영 agent·실제 상대·실제 `~/.claude`·`~/.codex`는 쓰지 않았다. 격리 폴더의 로그인 자격 파일은 열지도 복사도 하지 않았다(크기·수정시각만 확인).

### 환경

| 항목 | 값 |
|---|---|
| 제품 SHA | `66407760`(main), 격리 `git worktree add --detach`에서 빌드. 제품 파일 변경 0 |
| 호스트 | Windows 11 Pro 26300, Windows PowerShell 5.1, Claude Code **2.1.311**(검증 버전 2.1.302와 다름), Node 22.22.2(빌드·실행) |
| 로그인 | 사용자가 마친 격리 `CLAUDE_CONFIG_DIR`=`D:\workspace\KnowsLink\.devtools\autosetup-001-qa\claude-config`. `claude auth status` loggedIn true(coor 확인, tester 재확인). 자격 파일은 설치·세션 전후 크기 559바이트·수정시각 동일 |
| 배포물 | `package_plugin.py --release … --relay-url http://127.0.0.1:18442`, Claude ZIP `2d2c1884…`, Codex ZIP `45a09512…`(합성 상대 B 설치용) |
| relay | 격리 Postgres 17(Compose 프로젝트 `knowslink-autosetup-qa2`, 15482) + 1차 시험 fixture(합성 회원 a·b 친구 1쌍, Google 콜백만 대행) |
| 합성 상대 B | 회원 b 이름으로 설치한 Codex 설치기 산출물의 `plugin.js`를 스크립트 MCP 클라이언트로 구동([qa2-b_step.mjs](tools/qa2-b_step.mjs)). **Codex 호스트 세션은 쓰지 않았다**(상대 역할만) |
| 실세션 | Orca 터미널에서 `CLAUDE_CONFIG_DIR`을 위 폴더로 두고 `claude` 실행([화면 1](ev/qa2-claude-session-1.screen.txt), [재시작 후](ev/qa2-claude-session-restart.screen.txt)) |
| 설치 기록 | [Claude 설치(N1–N6, 17.8초)](ev/qa2-claude-install.log.txt), [합성 상대 B 설치](ev/qa2-peer-B-install.log.txt) |
| 도구 | [qa2-install_run.py](tools/qa2-install_run.py)(`QA_CLAUDE_CONFIG`로 로그인 폴더 지정), [qa2-fixture.sh](tools/qa2-fixture.sh), [qa2-b_step.mjs](tools/qa2-b_step.mjs), [qa2-repro.ps1](tools/qa2-repro.ps1) |

### AC 판정

| AC | 판정 | 근거 |
|---|---|---|
| AC-01 사람 행동 | **PASS** | 아래 행동 표. 경로·환경 변수·설정 파일 편집 입력 0, 붙여 넣은 셸 명령 1 |
| AC-02 Claude 첫 발송 | **PASS** | `/knowslink-session on` 뒤 기존 발송 승인 대화상자(「본문 전체 29자 … 그대로 보낼까요?」)에서 1회 승인 → `queued`, 상대 B `receive`가 같은 본문·발신자를 1건 받음([B 수신](ev/qa2-peer-B-recv.jsonl.txt)). 설치 완료·세션 연결만으로는 발송이 없었다: 연결 직후 `대기 0건`, B의 첫 `receive`가 정확히 이 1건 |
| AC-02 역방향 | **PASS** | B(합성) 발송은 B 쪽 4자리 코드 승인 후 `queued`([B 발송](ev/qa2-peer-B-send.jsonl.txt)). 실세션 Claude 화면에 알림 「KnowsLink 외부 메시지 도착 · 발신 친구 합성 회원 b … 대기 1건」과 상태줄 `대기 1`이 떴고, `knowslink_text_receive`로 본문·발신자를 받았다(`untrusted: true`, 세션 수신 기록 저장) |
| AC-03 Codex | **사용자 결정 제외(D040)** | 미검증 |
| AC-04 시간 | **PASS** | 명령 붙여 넣기 23:45:08 → 첫 발송 승인 완료 23:48:32 = **3분 24초**(한도 10분). 구간: 설치기 17.8초(내려받기·Node·등록·브라우저 대행 포함), `claude` 열기·첫 실행 안내, `/knowslink-session on` 23:47:25(연결 확인), 발송 요청 23:47:53, 승인 23:48:32([시각 기록](ev/qa2-timeline.txt)). 한계: 자동 구동 시간이다. 사람이 화면을 읽고 지문을 비교하는 시간은 더 걸린다. 이 시험에서 첫 실행 안내 중 신뢰 폴더 선택을 잘못 눌러 Claude Code를 한 번 다시 열었다(약 1분 포함) |
| AC-07 갱신 | **조건부(O-04 참조)** | 호스트 플러그인이 knowslink 하나뿐인 환경에서는 PASS: 같은 명령이 `기존 설치 갱신`, 5.5초, agent ID·`agent.json`·`private.pem` 해시·웹 agent 수 1·친구 수 1 동일([전](ev/qa2-update-before.txt)/[후](ev/qa2-update-after.txt)), 이전 설치 `.previous-*` 보존. 갱신 뒤 Claude Code를 다시 열어 `/knowslink-session on`이 같은 agent로 새 세션에 연결됨([재시작 후 화면](ev/qa2-claude-session-restart.screen.txt)). **호스트 플러그인 목록에 다른 플러그인이 함께 있으면 실패한다**(O-04) |
| AC-05·06 자격 경계 | **PASS** | 실제 `~/.claude`·`~/.codex` 전후: 변하는 파일은 1차와 같다(마켓플레이스 자동 새로고침으로 `known_marketplaces.json` 해시만 바뀜, 크기 같음, `knowslink` 0건. 이 세션 프로필 `.claude.json`은 실행 중이라 제외)([전](ev/qa2-real-before.txt)/[후](ev/qa2-real-after.txt)). 격리 로그인 폴더 `settings.json`에는 KnowsLink 항목 2개(`extraKnownMarketplaces.knowslink-plugins`, `enabledPlugins.knowslink@knowslink-plugins`)와 첫 실행의 `theme`만 있다. 출력·설치 로그에서 토큰·개인키 패턴 검색 0건. 지문은 앞 8자만 보고서·증거에 남겼다 |

### 사람 행동 표(Claude 2.1.311, 실세션)

| 단계 | 사람 행동 | 입력 종류 |
|---|---|---|
| N1 | 화면 명령 복사→PowerShell 붙여 넣기 | 붙여 넣기 1회 |
| N2 | 검증하지 않은 환경(2.1.311)에서 계속 `y` | 한 글자(O-01) |
| N3 | 설치 계획 동의 `y` | 한 글자 |
| N5 | 브라우저 Google 로그인·지문 비교·이름·동의 | 브라우저(fixture 대행) |
| 재시작 | 설치 때 Claude Code가 실행 중이 아니었으므로 없음 | — |
| N7 | `/knowslink-session on` | 명령 1 → `KnowsLink 연결됨: 세션 3e3cb6a2 · 내 agent agent_784f… · 대기 0건` |
| N8 | 발송 승인 대화상자 `1. 보내기` | 선택 1 |

합계 6(호스트 재시작 없음). 검증 버전 2.1.302면 5다. S1 한도 6에 닿는다(O-01). 별개로 새 `CLAUDE_CONFIG_DIR`에서는 Claude Code 자체의 첫 실행 안내(테마 선택·보안 안내·폴더 신뢰)가 따로 있다. 이미 Claude Code를 쓰는 사용자에게는 없는 단계라 S1 행동으로 세지 않았다.

### S1–S7 최종 판정

| ID | 판정 | 비고 |
|---|---|---|
| S1 Claude 행동 ≤6 | **PASS(경계)** | 6/6. 검증 버전 밖(2.1.311)이라 `y` 1개가 더해져 여유 0(O-01) |
| S2 Codex | **사용자 결정 제외(D040)** | 1차에서 N1–N6만 확인 |
| S3 ≤10분 | **PASS(Claude)** | 3분 24초(AC-04). Codex 미측정(D040) |
| S4 자격 경계 | **PASS** | 1차 + 실세션 후 재확인 |
| S5 갱신 | **조건부** | 단일 플러그인 환경 PASS, 다중 플러그인 환경 FAIL(O-04) |
| S6 실패 안내 | 1차 판정 유지 | 이번 범위 밖 |
| S7 화면 | **PASS** | 1차 판정 유지 |

### O-04 (high) — 호스트에 플러그인이 둘 이상이면 갱신·재실행이 실패한다

- 재현: 호스트 `claude plugin list --json`이 항목 2개 이상이고 그중 하나가 `knowslink@knowslink-plugins`일 때 같은 명령을 다시 실행한다. 이 시험에서는 Orca가 주입한 `CLAUDE_CODE_PLUGIN_DIRS`의 inline 플러그인(`orca-agent-view@inline`)이 둘째 항목이었다. 같은 명령을 그 환경변수만 지우고 실행하면 갱신이 통과한다([실패](ev/qa2-claude-update-fail-plugin-dirs.log.txt), [통과](ev/qa2-claude-update-ok.log.txt)).
- 증상: 버전 확인 `y` 직후 `[ERROR] 'Path' 매개 변수가 null이므로 인수를 해당 매개 변수에 바인딩할 수 없습니다.`로 종료(종료 코드 1). 실패 단계·바뀐 것·다음 행동 안내가 없다(S6 형식 미적용).
- 원인(코드 읽기와 합성 입력 재현으로 확인): `scripts/autosetup/knowslink-setup.ps1:243-246` `Get-Registered`가 `@(Get-HostJson …) | Where-Object { $_.id -eq $Plugin }`로 항목을 고른다. Windows PowerShell 5.1의 `ConvertFrom-Json`은 JSON 배열을 한 개의 배열 객체로 돌려주므로 `@()`가 이를 펼치지 못하고, 항목이 2개 이상이면 `Where-Object`가 배열 전체를 통과시켜 `$entry.readFromFolder`가 배열(null 포함)이 된다. 항목이 1개면 우연히 통과한다([합성 재현](ev/qa2-o04-repro.txt), [qa2-repro.ps1](tools/qa2-repro.ps1)).
- 영향: 다른 플러그인을 하나라도 설치한 Claude Code 사용자는 S5·AC-07(같은 명령 재실행 갱신)이 막힌다. 설치 도중 상태 변경 전에 종료하므로 기존 설치·키는 그대로다(전후 해시 동일). 첫 설치는 영향이 없다(이번 첫 설치가 같은 환경에서 성공).
- 한계: 실제 둘째 플러그인을 설치해 본 것이 아니라 Orca의 inline 플러그인과 합성 JSON으로 재현했다. 조치(수정)는 dev 판단이다(제품 코드는 고치지 않음).

### QA2 한계

- 실제 Google 로그인·동의는 하지 않았다(1차와 같은 fixture 대행). 사용자 실계정 비밀번호 입력 0.
- 상대는 합성 B의 스크립트 MCP 클라이언트다. 상대 쪽 호스트 UI는 검증하지 않았다.
- 발송 승인 대화상자와 도착 알림은 실세션 터미널 화면에서 직접 읽었다. 화면 캡처 이미지는 만들지 않았다(토스트는 텍스트 스트림에서 확인했고 [화면 1](ev/qa2-claude-session-1.screen.txt)에는 대화상자·알림 전 상태가 남는다).
- Claude Code 2.1.311이라 검증 버전(2.1.302)과 다르다. 화면은 계속 검증 버전을 말한다.

## 1차 검수 기록 (2026-10-23)

**1차 결론: 부분 완료.** 통과 조건(AC-01–AC-04 PASS)은 충족하지 못했다. 이유는 제품 결함이 아니라 검증 환경이다. 격리 `CLAUDE_CONFIG_DIR`·`CODEX_HOME`의 실제 호스트 로그인은 브라우저 OAuth 동의가 필요한데, 자동화 권한 분류기가 거부했다(우회하지 않음). 사용자가 로그인해 주면 남은 항목을 이어서 검증한다(`남은 일`).

## 고정 후보와 환경

| 항목 | 값 |
|---|---|
| 제품 SHA | `11d91648b30f13aefd545b19d0508fff9eb5b667` (fullops/tester = main `4fefdf68` + 문서 1커밋). 격리 `git worktree add --detach`에서 빌드, 제품 파일 변경 0 |
| 호스트 | Windows 11 Pro 26300, Windows PowerShell 5.1, Claude Code **2.1.309**(검증 버전 2.1.302와 다름), Codex CLI 0.163.0(검증 버전) |
| 배포물 | `package_plugin.py --release … --relay-url http://127.0.0.1:18441`, Claude ZIP `9361d4f1…`, Codex ZIP `dbea5ece…` |
| relay | 격리 Postgres 17(Compose 프로젝트 `knowslink-autosetup-qa`, 15481) + 합성 회원 A(Claude)·B(Codex) 친구 1쌍을 만드는 시험 fixture(`tools/autosetup_qa_fixture_test.go.txt`, 워크트리 밖 격리 복사본에서만 실행, 미커밋) |
| fixture가 대신한 것 | Google OAuth 콜백 한 가지(연결 시도에 회원 binding). 연결 화면(`/connect/<id>`)의 지문·이름 입력·`/home/device-confirm`은 실제 relay 코드로 실행했다. 실제 Google 로그인·동의는 하지 않았다 |
| 명령 | 회원 화면 `use=claude`·`use=codex`가 보인 한 줄을 그대로 실행. 바꾼 것은 운영 URL `https://link.knowslog.com`→격리 `http://127.0.0.1:18441` 뿐 |
| 도구 | [tools/install_run.py](tools/install_run.py), [qa.py](tools/qa.py), [snap.py](tools/snap.py), [realhash.py](tools/realhash.py), [mcp_roundtrip.mjs](tools/mcp_roundtrip.mjs) |

## S1–S7 판정

| ID | 판정 | 근거 |
|---|---|---|
| S1 Claude 행동 ≤6 | **부분 확인(N1–N6)**, N7–N8 미확인 | 아래 행동 표. 사람 입력 경로·환경 변수·설정 편집 0, 붙여 넣는 명령 1. 현재 호스트가 2.1.309라 '검증하지 않은 환경에서 계속' `y`가 한 번 더 있다 → 관찰 O-01 |
| S2 Codex 행동 | **부분 확인(N1–N6)**, 세션 연결·코드 입력 미확인 | `CODEX_HOME`·marketplace·plugin add 수동 실행 0. 검증 버전이라 추가 `y` 없음 |
| S3 ≤10분 | **미확인** | 설치기 구간만 실측: Claude 19.7초, Codex 21.3초(Node 22.22.2 실제 내려받기 포함, 브라우저 단계는 fixture라 1.5초). 세션 연결까지는 실호스트 로그인 필요 |
| S4 자격 경계 | **PASS(설치기 구간)** | 격리 호스트 설정에는 KnowsLink 항목 2개만 생김(ev/ 아래 증거). 출력·설치 기록에 credential·개인키 0건(검색). 호스트 인증 파일을 읽은 흔적 0(상자에 인증 파일 자체가 없음). 사용자 실제 파일: 아래 해시 표 |
| S5 갱신 | **PASS** | 같은 명령 재실행이 Google 없이 `기존 설치 갱신`. 호스트별 전후 agent ID·`agent.json`·`private.pem`을 포함한 private 폴더 해시·웹 agent 수(1)·친구 수(1) 동일 ([claude 전](ev/snap-claude-before-update.txt)/[후](ev/snap-claude-after-update.txt), [codex 전](ev/snap-codex-before-update.txt)/[후](ev/snap-codex-after-update.txt)). 갱신 시간 3.6초·2.1초. 이전 설치 `.previous-*` 보존 |
| S6 실패 안내 | **미확인(이번 실측 범위 밖)** | F-01–F-09는 DEV가 KO·EN PASS(주입 fixture). 이번에는 EN에서 계획 거절(F-02 후속)만 확인: `Not started. Nothing was changed.`, 설치 폴더 생성 0 |
| S7 화면 | **PASS** | 회원 화면 선택지 6개(web·grok·dots·claude·codex·other) 중 명령은 claude·codex에만 있다. KO·EN 모두 '시험 환경 확인 (Claude Code 2.1.302, Windows 11)', Codex는 '도착 알림 없음'. 실제 설치 호스트가 2.1.309여도 화면은 검증 버전을 말한다 ([claude 화면](ev/screen-claude.html.txt), [codex 화면](ev/screen-codex.html.txt)) |

## AC-01 사람 행동 표(실측 N1–N6)

| 단계 | 사람 행동 | 입력 종류 | Claude | Codex | 기준선(intent 9단계) |
|---|---|---|---|---|---|
| N1 | 화면 명령 복사→PowerShell 붙여 넣기 | 붙여 넣기 1회 | ○ | ○ | 5개 이상 명령 + 절대 경로 2번 이상 |
| N2/AS-09 | 검증하지 않은 환경에서 계속 `y` | 한 글자 | ○ (2.1.309) | 없음 | 해당 없음 |
| N3 | 설치 계획 동의 `y` | 한 글자 | ○ | ○ | 해당 없음 |
| N4 | (자동) 내려받기·checksum·Node·설치·등록·등록 확인 | 없음 | 완료 줄 4개 | 완료 줄 4개 | 단계 1–5 |
| N5 | 브라우저: 로그인·지문 비교·이름·동의 | 브라우저 | fixture 대행 | fixture 대행 | 단계 7 |
| N6 | 완료 화면 읽기(재시작 안내) | 없음 | ○ | ○ | 단계 8 |
| 재시작 | Claude Code/Codex를 닫고 다시 열기(실행 중이었을 때만) | 없음 | 조건부 | 조건부 | 해당 없음 |
| N7 | `/knowslink-session on` / `knowslink_session on`+Codex 도구 승인+4자리 코드 | 명령·코드 | **미확인** | **미확인** | 단계 9 |
| N8 | 메시지 요청·발송 승인 | 승인 | **미확인** | **미확인** | 단계 10 |

경로·환경 변수·설정 파일 편집 입력은 모두 0이었다. 셸에 붙여 넣은 명령은 호스트당 1개다. 터미널에 보인 지문과 연결 화면의 지문은 두 호스트 모두 같았다(`SHA256:OXwJi6Vo…`, `SHA256:Mo9CCGun…`). N1–N6의 호스트별 사람 입력은 Claude(2.1.309) 3개(붙여 넣기·`y`·`y`)+브라우저 1, Codex 2개(붙여 넣기·`y`)+브라우저 1이다. N7·N8을 더하면 Claude 총 6(재시작 없음)~7(재시작 있음), 검증 버전의 Claude Code 2.1.302면 5~6이다.

## AC-02/03 첫 발송 — 실호스트 미확인, MCP 수준 대체 확인

실호스트 세션이 없어 AC-02/03은 **판정하지 않는다.** 대신 설치기가 호스트 설정에 등록한 MCP 서버 설정(`config/local.mcp.json`, `.mcp.json`, 절대 경로 node·plugin.js·agent 폴더 포함)을 그대로 읽어 스크립트 클라이언트로 띄웠다 ([mcp-roundtrip.jsonl](ev/mcp-roundtrip.jsonl)).

| 방향 | 결과 |
|---|---|
| Codex B `knowslink_session on` | 확인 코드 창(세션 구분값·내 agent ID 표시) → 코드 입력 → `KnowsLink 연결됨: Codex 세션 … · 대기 0건` |
| Claude A → Codex B | `transport: queued` → B `receive`로 본문·발신자 이름 일치(`untrusted: true`, `sessionReceipt: recorded`) |
| Codex B → Claude A | 발송 창에 본문 전체 22자 표시 → 코드 입력 → `queued` → A `receive` 일치 |

이것은 설치기 산출물이 실제로 동작하는 MCP 설정을 만든다는 증거다. 호스트의 도구 discovery, `/knowslink-session on` mod, Claude 발송 승인 대화상자, Codex 도구 실행 승인은 검증하지 않았다.

## 사용자 실제 설정 파일 전후 (S4)

[전](ev/real-before.txt)·[후](ev/real-after.txt) 비교. 후 해시는 `realhash.py`를 같은 조건으로 다시 실행한 결과다. 값이 달라진 파일은 `~/.claude/plugins/known_marketplaces.json`(및 같은 파일을 가리키는 이 세션 프로필)이다. 크기 1751바이트 동일, `knowslink` 문자열 0건, 모든 항목의 `lastUpdated`가 `2026-10-23T01:20:46Z`로 갱신됐다. 호스트의 마켓플레이스 자동 새로고침으로 보이며 KnowsLink 항목은 생기지 않았다. 설치기는 항상 격리 `CLAUDE_CONFIG_DIR`·`CODEX_HOME`로 실행했다(설치기가 만든 항목은 격리 상자 안에서만 확인). `~/.codex/config.toml`·`auth.json`·`~/.claude/settings.json`·`installed_plugins.json`은 전후 동일하다. 이 세션 프로필의 `.claude.json`은 실행 중인 이 세션이 계속 쓰므로 비교 대상에서 뺀다.

## 관찰·결함

1차 시점 심각도 high 0, 미해결 critical 0. QA2에서 high O-04가 추가됐고 QA3에서 해소됐다(위 QA3 절).

| ID | 심각도 | 내용 | 재현·영향 |
|---|---|---|---|
| O-01 | low–medium(S1 여유 소진) | 현재 설치된 Claude Code는 2.1.309인데 검증 버전은 2.1.302다. 새 사용자는 호스트가 자동 갱신되므로 대부분 같은 `y`를 한 번 더 입력한다. S1 행동 수가 한도 6에 닿고, 호스트 재시작이 필요하면 7이 된다 | 재현: 2.1.309에서 화면 명령 실행 → `검증하지 않은 환경에서 계속할까요? (y/N)`. 조치는 designer 판단(S1 목록에 포함 또는 검증 버전 갱신) |
| O-02 | low | 같은 고정 포트(18431)에 다른 팀의 relay가 떠 있으면 격리 fixture가 조용히 실패하고 다른 서비스의 404가 응답한다. 제품 문제가 아니라 시험 환경 충돌. 포트를 18441로 옮겼다 | 시험 환경 기록 |
| O-03 | info | 완료 화면의 `시험 환경 확인` 문구가 Codex(검증 버전)에만 나오고 Claude 2.1.309에는 없다. AS-09 의도대로다 | 화면 대조 |

## 한계 (정직한 범위)

- 실제 Google 로그인·동의는 하지 않았다(fixture 대행). 사용자 실계정 비밀번호 입력 0.
- 실제 Claude Code·Codex 로그인과 세션은 쓰지 못했다(위 결론). 세션 연결·첫 발송·소요 시간 AC-02/03/04는 미확인이다.
- Claude 첫 설치 원시 로그는 도구가 같은 이름으로 덮어써 잃었다(결과는 이 보고서에 옮겼다). 이후 실행은 시각 접미사를 붙인다. Codex 첫 설치·갱신 로그는 [codex-ko1-install.txt](ev/codex-ko1-install.txt), [codex-ko1-update.txt](ev/codex-ko1-update.txt).
- F-01–F-09 전체와 EN 전체는 이번에 다시 실행하지 않았다(DEV 결과 신뢰, 근거 `SAR-AUTOSETUP-001-DEV-test/report.md`).

## 남은 일

1. dev: O-04 수정(`Get-Registered`가 `claude plugin list --json` 항목이 여러 개여도 `knowslink@knowslink-plugins` 한 개만 고르게 한다). 수정 뒤 호스트 플러그인이 2개 이상인 환경에서 같은 명령 재실행(갱신)을 다시 검증한다(tester, 약 10분). 함께 갱신·재실행 실패 시 S6 형식(실패 단계·바뀐 것·다음 행동)이 보이는지 확인한다.
2. designer: O-01(S1 행동 수 6/6 여유 0과 검증 버전 문구) 판단. D02·D04 변경은 이번 검수에서 필요하지 않았다.
3. Codex 실세션 N7–N8·AC-03·S2·S3는 사용자 결정 D040으로 제외했다. Codex를 검수 범위에 되돌리면 격리 `CODEX_HOME` 로그인 뒤 재개한다.
4. 사용자: 격리 로그인 폴더 `D:\workspace\KnowsLink\.devtools\autosetup-001-qa\claude-config`는 지우지 않고 뒀다(사용자 소유, 로그인 유지).
