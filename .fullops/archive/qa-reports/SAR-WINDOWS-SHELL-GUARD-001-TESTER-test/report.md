---
title: SAR-WINDOWS-SHELL-GUARD-001-TESTER lite QA
status: draft
updated: 2026-10-11
owner: tester
tasks: [SAR-WINDOWS-SHELL-GUARD-001-TESTER]
summary: 4e82b4f Windows 삭제 실행 차단의 독립 경계 QA와 실제 hook 결선 증거
---

# SAR-WINDOWS-SHELL-GUARD-001-TESTER lite QA

대상 SHA는 `4e82b4f4cd1bc31a59d0d230349fd34d95c7f0c7`다. 읽기 전용 snapshot은 `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-4ff90f54195a4728bc150bd64f431020`다. snapshot은 detached HEAD이며 `git status --porcelain`이 비어 있다. 검수 세션은 `0ef9c64c-15ec-4fcc-a3be-c11022475dc1`이고 구현 세션 `bacd5d1f-de07-43b3-9efa-0df5636729c7`과 다르다. 테스트 레벨은 lite다.

모든 검사에서 실제 삭제·이동·복구·cleanup을 실행하지 않았다. 합성 payload만 guard 프로세스의 stdin으로 보냈고 probe 문자열은 어떤 셸에도 전달하지 않았다. 임시 workspace(`scratchpad\tester-ws-*`, `tester-out-*`)와 DEV 회귀의 `%TEMP%\shell-guard-ws-*`는 보존한다. 모든 probe 뒤에 대상 경로가 그대로 있음을 확인했다(`exists after: True`).

## 근거 재사용 — DEV 등록 lint/test

| 확인 항목 | 값 |
|---|---|
| DEV 체크아웃 HEAD | `4e82b4f4cd1bc31a59d0d230349fd34d95c7f0c7`, `git status` 깨끗함 |
| DEV git-dir `fullops-gate/pass.json` | head 4e82b4f, base a90a9e6, attempt `4ad09a11…`(lint-attempt.json과 일치). pass.json은 ERROR 0·unavailable 0일 때만 기록된다(lint.py `stamp`) |
| DEV 보존 `final2-lint.json`(sha256 `0edf4a1a…29fe`) | head 4e82b4f, attempt 일치, `product-lint` passed exit 0, `product-test` passed exit 0, ERROR 0 / WARNING 5 / unavailable 0 |
| 같은 SHA Linux 로그 | `…\KnowsLinkDevTools\wsg001\4e82b4f…-test.log` 끝에 `python3 scripts/test_shell_guard.py` `Ran 3 tests OK`, `EXIT=0` |

위 근거로 제품 `make lint`·`make test`는 새로 실행하지 않았다(변경·근거 결함 없음, Docker 로컬 실행 금지). 근거는 DEV 실행이며 tester 실행이 아니다. 리뷰 기록의 `lint.json`은 이 보존 파일의 복사본이다.

## 직접 실행한 검사

| 검사 | 명령 또는 방법 | 결과 |
|---|---|---|
| DEV 회귀 재실행 | snapshot에서 `python3 -B scripts/test_shell_guard.py` | exit 0, 3 tests OK(14.2초). 삭제 없음 확인 테스트 통과. snapshot 변경 없음 |
| 실제 Claude Bash 도구 결선 | 이 세션(tester 체크아웃)에서 `bash -c "echo rm …"` | 실행 전 `PreToolUse:Bash hook error: Blocked by scripts/shell_guard.py` |
| 실제 Claude PowerShell 도구 결선 | `cmd /c echo rd …` | 실행 전 `PreToolUse:PowerShell hook error: Blocked by scripts/shell_guard.py` |
| 실제 Claude 오탐 관찰 | probe 스크립트를 heredoc으로 쓰려던 Bash 호출 2회 | 둘 다 차단됨(삭제 단어와 `sh` 단어가 한 명령에 있음). 알려진 오탐이며 실제 작업을 막는다 |
| Grok 발견 | tester 체크아웃에서 `grok inspect --json`(모델 호출 없음) | `pre_tool_use` project hook `python3 "${CLAUDE_PROJECT_DIR}/scripts/shell_guard.py"`, matcher `Bash|PowerShell`, `compatibilityStatus: enabled`, `projectTrusted: true` |
| Grok 종료코드 의미 | 설치된 `~/.grok/docs/user-guide/10-hooks.md` | exit 2는 PreToolUse 명시 거부. timeout·crash·잘못된 출력은 fail-open |
| Codex 기능 상태 | `codex features list` | `hooks stable true`. codex-cli 0.162.0 |
| Codex hook 의미 | 공식 문서 `learn.chatgpt.com/docs/hooks`(Context7은 월 한도 초과로 실패) | exit 2는 차단. 명령은 session cwd에서 실행. 미신뢰·변경된 hook은 신뢰 전까지 건너뜀. shell 명령 matcher는 `Bash`. `apply_patch`·MCP 도구도 PreToolUse 대상이나 이 matcher에는 포함되지 않음 |
| 경계 matrix | `probes/probe.py`(약 120건)·`probe2.py`·`probe3.py`, 결과 `probe-out.txt` 등 | 아래 판정 |
| 입력 크기 | `probes/timing.py` | 400KB 단일 단어 1.76초, 100KB 다줄·커밋 메시지형 0.2초 이하. 2MB 단일 단어는 15초 안에 끝나지 않음(hook timeout 10초 초과 = fail-open) |
| 종료코드 | `python3 nonexistent.py` | exit 2(파일 없음) |

실행하지 않은 검사: 실제 Grok·Codex 모델 세션의 tool call 차단, Codex 사용자 hook 신뢰, 실제 삭제·reparse 하위 지점 삭제 동작, Linux pwsh 경로, 전역·managed 설정 강제.

## 판정 요약

차단 확인(deny): 사고 형태(PowerShell 이중 따옴표 안 `bash -c "rm -rf \$VAR/"`, `bash.exe` 전체 경로, `wsl -e sh -c`, `cmd /c rd /s`, `powershell.exe -Command`, `ssh`, `eval`, `Invoke-Expression`, 인코딩 명령), 변수·와일드카드·빈 값·루트 대상, `rm -rf`의 모든 조합 표기·대소문자·탭·CRLF·서브셸·백그라운드, `find -delete/-exec/-ok`, `xargs`, 파이프의 `Remove-Item`, 허용 형식 밖의 `Remove-Item`(`-Path`·변수·다중 경로·쉼표 경로·`-Confirm`), junction·상위 junction·워크스페이스 밖·상위·형제 접두·ADS·8.3 이름, `.git`, 드라이브 루트. 입력 경계: 잘못된 JSON·`null`·빈 입력·리스트는 exit 2. 허용 확인(allow): 허용 형식의 workspace 내부 경로, `git status`, `go test`, `make lint`, `grep "rm -rf"`, 읽기 명령. `Remove-Item Env:FOO`도 통과.

미해결 발견 9건. critical·high는 없다. 상세와 재현은 `SAR-WINDOWS-SHELL-GUARD-001-REVIEW-review/report.md`다.

| 번호 | 심각도 | 요약 | 재현 입력(실행하지 않음, guard 판정만) |
|---|---|---|---|
| F1 | medium | PowerShell 콜론 스위치 약어가 재귀 판정을 우회한다 | `Remove-Item -r:$true C:\Users\x\Documents` → allow. `Get-Command Remove-Item`에서 `R*` 매개변수는 `Recurse` 하나라 `-r`은 유일 접두사다. 변수·와일드카드·루트 대상은 같은 약어에서도 deny다 |
| F2 | medium | 인라인 인터프리터의 셸 호출·`rmdirSync`가 감지되지 않는다 | `python -c "import os; os.system('rm -rf /')"`, `python3 -c "subprocess…['rm','-rf']"`, `node -e "…rmdirSync(…recursive…)"`, `node -e "…execSync('rm -rf x')"` → allow |
| F3 | medium | `git clean`이 `-C`·전역 옵션 뒤에서 감지되지 않는다 | `git -C . clean -fdx`, `git --no-pager clean -fdx` → allow |
| F4 | medium | Codex hook이 session cwd 상대경로라 하위 폴더 세션에서 모든 셸 호출을 막는다. 문서의 fail-open 서술이 틀렸다 | 누락 스크립트의 python 종료코드 2 실측 + Codex 문서의 exit 2 = block |
| F5 | low | PowerShell 문법과 tokenizer 불일치로 명령이 문자열에 숨는다 | `Get-ChildItem "C:\d\"; Remove-Item -Recurse C:\x # "` → allow |
| F6 | low | 같은 부류의 미감지 형태 | `robocopy a b /MOVE`, `rsync --delete`, `docker run -v … rm -rf`, `busybox rm -rf`, `. Remove-Item -Recurse …`, `& ([scriptblock]::Create("…"))`, `r\m` |
| F7 | low | 오탐이 실제 작업을 막는다 | heredoc·`git commit -m` 본문의 삭제 단어+`sh`, `-WhatIf`·`-Confirm:$false` 붙은 허용 형식, `git clean -fdn` |
| F8 | low | 아주 큰 명령은 10초 timeout을 넘겨 fail-open | 2MB 단일 단어 |
| F9 | low | 범위 한계(설명·결정 필요) | 아래 |

F9 범위 한계: 비재귀 삭제는 workspace 밖 절대경로도 허용한다(`Remove-Item C:\Users\x\f.docx`). `Move-Item`·`mv`는 검사하지 않는다. 허용한 디렉터리 안쪽의 중첩 reparse 지점과 TOCTOU는 검사하지 않는다(실제 삭제 금지로 미검증). hook은 프로젝트 체크아웃 안에서만 적용된다. 다른 저장소·전역 셸·Gemini·agy 등 hook 없는 클라이언트·`apply_patch`·MCP 삭제는 보호하지 않는다. 스크립트 파일·별칭·`make`는 DEV 문서가 밝힌 대로 우회한다.

## 결선 판정

- Claude Code: 실제 활성 확인(두 도구 차단). 새 세션이 시작할 때 설정을 읽는다. 다른 역할 체크아웃은 main 병합·동기화 뒤 적용된다(designer dc21d6c, ops 7a69068에는 아직 없다).
- Grok: 발견·신뢰·활성은 확인했다. 실제 모델 tool call 차단은 미확인이다. 실행 중인 세션은 재시작 또는 `/hooks` Reload가 필요하다.
- Codex: 설정 파일·기능만 확인했다. 사용자 hook 신뢰 전에는 실행되지 않으며 적용 완료로 표시하지 않는다. 하위 폴더 세션 결함은 F4다.
