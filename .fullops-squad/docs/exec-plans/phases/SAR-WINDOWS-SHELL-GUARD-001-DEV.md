---
title: Windows 삭제 명령의 실행 전 차단
status: draft
updated: 2026-10-11
owner: dev
tasks: [SAR-WINDOWS-SHELL-GUARD-001-DEV]
summary: Claude·Grok·Codex PreToolUse hook으로 교차 셸 삭제와 검증되지 않은 재귀 삭제를 실행 전에 거부하고 삭제 없는 회귀로 검증한다
---

# Windows 삭제 명령의 실행 전 차단

## 기술 계획과 범위

기준 ref는 a90a9e6, 준비 커밋은 0e4f638이다. fullops-common-0.3.3과 project.md를 적용한다. 예외는 없다. 테스트 레벨은 lite, 하위 위임은 사용하지 않았다(단일 파일 구현이라 위임 이득이 없다).

사고 원인은 PowerShell 큰따옴표 문자열 안의 Bash 명령이다. PowerShell이 `\$VAR`의 `$VAR`를 먼저 빈 값으로 확장했다. 그 결과 재귀 삭제 대상이 루트가 됐다. 지침은 이미 있었지만 실행을 막지 못했다. 그래서 이번 과제는 호스트가 실행 직전에 호출하는 PreToolUse hook으로 거부한다.

1. `scripts/shell_guard.py`: hook payload의 셸 명령 문자열을 읽고 거부 사유를 결정한다. 거부하면 PreToolUse deny JSON을 stdout에, 사유를 stderr에 쓰고 exit 2로 끝난다. 허용하면 아무것도 쓰지 않고 exit 0이다. payload를 읽지 못하면 거부한다.
2. `.claude/settings.json`: Claude Code와 Grok이 읽는 `PreToolUse` `Bash|PowerShell` hook이다.
3. `.codex/hooks.json`: Codex가 읽는 `PreToolUse` `Bash` hook이다.
4. `scripts/test_shell_guard.py`: 세 provider payload의 합성 deny/allow 회귀다. `make test`에 연결했다. `make lint`는 두 파일을 `py_compile`한다.
5. AGENTS.md·CLAUDE.md에 hook 위치·허용 형식·우회 한계를 적었다. Grok과 Codex는 AGENTS.md를 읽는다. Claude Code는 CLAUDE.md를 읽는다. 별도 GROK.md는 만들지 않았다.

제품 relay/adapter 코드, 플러그인 캐시, 사용자 전역 설정(`~/.grok`, `~/.codex`, `~/.claude`)은 바꾸지 않았다. D10 모듈 설계는 제품 모듈 변경이 없어 갱신하지 않는다. 새 의존성·비용·외부 서비스·권한 완화는 없다.

## 거부 규칙과 허용 경로

| 규칙 | 거부 대상 예 |
|---|---|
| 교차 셸·평가기 삭제 | 삭제 단어와 `bash -c`·`sh`·`cmd /c`·`pwsh -Command`·`powershell`·`wsl`·`ssh`·`Invoke-Expression`·`Start-Process`·`eval`이 한 명령에 있음 |
| 인코딩 명령 | 셸 실행과 `-EncodedCommand` |
| 재귀 삭제 | `rm -r/-rf/-fr/--recursive`, `Remove-Item -Recurse`(약어 포함), `rd /s`, `del /s`, `rimraf` — 허용 형식 이외 전부 |
| 대상 미검증 삭제 | 변수(`$`,`%VAR%`)·와일드카드·빈 값·루트(`/`, `C:\`, `/c`, `/mnt/c`) 대상, 파이프·`xargs` 입력 |
| 기타 | `find -delete/-exec rm`, `git clean`(`-n` 제외), `robocopy /MIR·/PURGE`, 인라인 `python -c`/`node -e`의 `rmtree` 등, `[IO.Directory]::Delete` |

허용하는 재귀 삭제는 한 단순 명령 `Remove-Item -LiteralPath '<절대경로>' -Recurse [-Force]`다. 경로는 작은따옴표 하나로 감싸야 한다. 경로에는 변수·와일드카드·`.`/`..` 구간이 없어야 한다. 경로는 존재해야 한다. 경로는 링크·junction·reparse 지점이 아니어야 하고 상위 경로도 링크를 지나면 안 된다(`realpath` 일치). 경로는 workspace 안이어야 하고 workspace 자체·드라이브 루트·`.git`이면 안 된다. workspace는 `CLAUDE_PROJECT_DIR`, 없으면 payload `workspaceRoot`, 없으면 `cwd`다.

검증되지 않은 대상이 없는 비재귀 삭제(`rm build/x.tmp`)와 읽기·빌드·테스트·git 명령(`git rm --cached`, `git clean -n` 포함)은 허용한다.

## 실제 지원 근거

| CLI | 설치 버전 | hook 위치·조건 | 근거 |
|---|---|---|---|
| Claude Code | 2.1.296 | 프로젝트 `.claude/settings.json`. Windows는 `Bash`와 `PowerShell` 도구를 모두 match해야 한다. exit 2는 차단이다. | Context7 `/websites/code_claude` hooks·tools-reference |
| Grok | 1.0.50 | 프로젝트 `.claude/settings.json`을 Claude 호환으로 읽는다(기본 사용). 폴더 신뢰가 필요하다. `Bash`는 `run_terminal_command` 별칭이다. `${CLAUDE_PROJECT_DIR}`는 runner가 확장한다. stdout의 deny JSON은 종료코드와 관계없이 차단한다. 실패한 hook은 fail-open이다. | 설치된 `~/.grok/docs/user-guide/10-hooks.md` |
| Codex | codex-cli 0.162.0, feature `hooks` stable | 프로젝트 `.codex/hooks.json`. 프로젝트 신뢰와 hook별 신뢰(`hooks.state.<key>.trusted_hash`)가 모두 필요하다. 신뢰 전에는 실행하지 않는다. Windows는 `cmd /C`로 session cwd에서 실행하고 프로젝트 hook에는 env 치환이 없다. 셸 도구 이름은 `Bash`, 입력은 `{"command": cmd}`다. | Context7 `/openai/codex` hooks discovery·command_runner·exec_command |

## 검증 결과

모든 검사에서 실제 삭제는 실행하지 않았다. 사고 원문 명령은 F: 비공개 자료에만 있으며 실행·복사하지 않았다. 회귀는 같은 형태의 합성 입력을 쓴다.

| 검사 | 대상 | 결과 |
|---|---|---|
| `python3 scripts/test_shell_guard.py` | Windows 로컬, 최종 작업 트리 | 3 tests OK. provider 4종(Claude Bash·Claude PowerShell·Grok·Codex) × deny 43·allow 17, payload 경계 3종, 삭제 없음 확인 |
| 같은 회귀(사전 검사) | Linux 고유 clone `/tmp/knowslink-wsg001-124f4af3b3b8-86b21d59` | OK (symlink 경로) |
| `make lint`·`make test` | 같은 Linux clone, HEAD 124f4af·clean 확인 | 둘 다 exit 0 |
| Claude Code 실제 hook | 이 세션의 Bash 도구 `bash -c "echo rm"`, PowerShell 도구 `cmd /c echo rd` | 둘 다 실행 전 `PreToolUse ... hook error: Blocked by scripts/shell_guard.py` |
| Grok 실제 발견 | `grok inspect --json` (모델 호출 없음) | `pre_tool_use` project hook, source `.../dev/.claude`, matcher `Bash|PowerShell`, vendor claude, `compatibilityStatus: enabled`, `projectTrusted: true` |
| Grok runner 셸 재현 | `pwsh -NoProfile -Command 'python3 "$env:CLAUDE_PROJECT_DIR/scripts/shell_guard.py"'`에 Grok payload 파일 | deny: stdout deny JSON(pwsh exit 1, Grok은 JSON deny를 따른다). allow: exit 0·출력 없음 |
| Codex runner 셸 재현 | `cmd.exe /C "python3 scripts/shell_guard.py"`에 Codex payload 파일 | exit 2, deny JSON |

자체 리뷰에서 wrapper 옵션(`env -i`, `timeout 5`), `rm --rec`, PowerShell splatting(`@p`) 우회를 찾아 고치고 deny 4건을 추가했다. 중간 커밋 6e3881a의 lint.py gate는 product-lint·product-test passed, ERROR 0, WARNING 5다. SLOP-004 3건은 hook 프로토콜의 deny JSON·stderr 사유와 회귀 증거 경로 출력이며 디버그 잔재가 아니다. SIZE-001은 기존 PLANS.md 누적 크기, SIZE-002는 exec plan·packet-outcomes 문서를 포함한 추가량이다. guard·회귀 코드는 각각 300줄 이내다. 최종 HEAD의 등록 lint/test와 FullOps gate 결과는 worker_done과 Git 디렉터리 `fullops-gate/pass.json`이 정본이다.

## 활성화 상태와 재시작

- Claude Code: 이 세션에서 설정 저장 직후 차단을 확인했다. 새 세션은 시작할 때 읽는다. 병합 전에는 dev 체크아웃에만 있다. 다른 역할 체크아웃은 main 병합과 동기화 뒤 적용된다.
- Grok: 새 세션이 시작할 때 읽는다. 실행 중인 Grok 세션은 `/hooks`의 Reload 또는 재시작이 필요하다. 워크트리 신뢰는 원본 `D:\workspace\KnowsLink` 신뢰로 충족된다. 실제 Grok 모델 세션의 tool call 차단은 비용 때문에 실행하지 않았다.
- Codex: hook 신뢰 승인 전에는 실행되지 않는다. 사용자가 이 체크아웃에서 Codex를 시작해 hook 검토 화면의 Trust를 선택해야 한다. 신뢰 키에 파일 경로가 들어가므로 체크아웃마다 한 번 승인한다. 전역 설정은 수정하지 않았다.
- hook 실행 실패(python3 미해결, timeout 10초 초과)는 세 CLI 모두 fail-open이다. python3는 Git Bash·pwsh·Windows PowerShell PATH에서 `Python313\python3.exe`로 해결됨을 확인했다.

## 알려진 우회와 한계

- 문자열 검사이며 완전한 보안 경계가 아니다. 스크립트 파일 실행(`python cleanup.py`, `make <target>`), 별칭·함수 정의, 문자열 조합·인코딩 변형, 파일 쓰기 도구로 만든 스크립트는 우회한다. 셸 도구 이외(MCP·편집 도구)는 검사하지 않는다.
- 경로 검사와 실제 실행 사이에는 시간 차가 있다(TOCTOU).
- 오탐: 삭제 단어와 셸 실행 단어가 한 명령에 함께 있으면 거부한다. 예: 커밋 메시지·heredoc·`Remove-Item Env:X`와 `pwsh`의 조합. 커밋 메시지는 파일(`git commit -F`)로 전달한다.
- Codex 프로젝트 hook은 session cwd 기준 상대경로라 하위 폴더에서 시작한 Codex에서는 실패하고 fail-open이다.

## 실행하지 않은 검사

실제 삭제, 사고 원문 재실행, 실제 Grok·Codex 모델 세션의 tool call(비용), Codex hook 신뢰 승인(사용자 결정), 전역 설정 기반 강제(managed/requirements), Linux pwsh 경로. 독립 고정 SHA 리뷰는 구현자와 다른 세션의 후속이다.

## 보존한 임시 증거

삭제·이동 정리를 하지 않았다. 아래 경로를 보존한다.

- Windows 회귀 workspace: `C:\Users\shin\AppData\Local\Temp\shell-guard-ws-*`, `shell-guard-outside-*` (실행마다 1쌍).
- Linux helper·로그: `C:\Users\shin\AppData\Local\KnowsLinkDevTools\wsg001\` (`remote.py`, `make.cmd`, `<SHA>-install/lint/test.log`, bundle).
- 서버 clone: `/tmp/knowslink-wsg001-<SHA12>-<random>/`, 서버 회귀 workspace `/tmp/shell-guard-ws-*`.
- provider payload 파일: 세션 scratchpad `grok-deny.json`, `grok-allow.json`, `codex-deny.json`.
