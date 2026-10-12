---
title: Codex hook 실행 오류와 셸 인용 회귀 수정
status: draft
updated: 2026-10-12
owner: dev
tasks: [SAR-CODEX-HOOK-QUOTE-001-DEV]
summary: Codex PreToolUse hook 명령을 cmd·PowerShell·sh에서 같은 종료코드로 실행되게 고치고 실제 실행 방식의 회귀로 검증한다
---

# Codex hook 실행 오류와 셸 인용 회귀 수정

## 기술 계획과 범위

기준 ref는 67e7de9, 준비 커밋은 7cae8d7이다. fullops-common-0.3.3과 project.md를 적용한다. 예외는 없다. 테스트 레벨 lite, 하위 위임은 사용하지 않았다(소규모 직렬 수정). 제품 코드·전역 설정·플러그인 캐시는 바꾸지 않았다. `scripts/shell_guard.py`는 변경 필요가 입증되지 않아 수정하지 않았다.

| 파일 | 변경 |
|---|---|
| `.codex/hooks.json` | 따옴표 후보 뒤에 `; exit $LASTEXITCODE`를 붙였다. |
| `scripts/test_shell_guard.py` | `test_codex_hook_from_subfolder`가 실제 Codex 실행 방식(`cmd.exe /C "<command>"` raw 인자)과 PowerShell을 검증한다. 비Windows는 sh로 실행한다. |
| `AGENTS.md` | "Windows `cmd /C` 전용" 설명을 실제 보장 범위로 바꿨다. |

## 원인과 실행 계약

이전 phase(SAR-WINDOWS-SHELL-GUARD-001-DEV)는 hook이 Windows `cmd /C` 전용이라고 전제했다. 현재 재현은 이 전제와 충돌한다. 아래는 같은 입력(`git status` 허용 0, 삭제 문자열은 JSON 입력으로만 제공한 거부 2)을 `scripts/` 하위 폴더 cwd에서 실행한 실측이다.

| 실행 방식 | 기준 67e7de9 (따옴표 없음) | 후보 7cae8d7 (따옴표만) | 수정 |
|---|---|---|---|
| `cmd.exe /C "<command>"` raw 인자(Codex 기본) | 0 / 2 | 0 / 2 | 0 / 2 |
| `shell=True`(cmd) | 0 / 2 | 0 / 2 | 0 / 2 |
| `pwsh -Command` | parser error, 종료 1 | 0 / **1** | 0 / 2 |
| `powershell -Command` | parser error, 종료 1 | 0 / **1** | 0 / 2 |
| argv 목록 `['cmd.exe','/C',command]` | — | 0 / **0**(guard 미실행) | 0 / **0** |

- 원 실패 `Hook failed code 1`은 기준 명령을 PowerShell이 파싱할 때의 parser error와 일치한다. 호스트가 낸 원 stderr는 직접 확보하지 못했다.
- 후보는 PowerShell 파서 오류를 없앴지만 `-Command`가 native 종료코드 2를 1로 바꿔 Codex가 차단하지 않았다(Codex는 2만 차단). 끝의 `exit $LASTEXITCODE`가 이를 막는다. cmd에서는 `;`가 인자 구분자이고 python의 `sys.argv`로만 넘어가 영향이 없다. sh에서는 `LASTEXITCODE`가 비어 직전 종료코드를 그대로 쓴다.
- 기존 회귀 실패(expected 2 / actual 0)의 원인은 테스트의 argv 목록 호출이다. Python이 목록의 따옴표를 `\"`로 바꿔 cmd가 python 코드를 문자열 리터럴로 넘겼다. 이는 Codex의 실제 실행 방식이 아니다.
- Codex 실행 계약(openai/codex `codex-rs/hooks/src/engine/command_runner.rs`, main 브랜치 원문): Windows는 셸 미설정이면 `%COMSPEC%`(없으면 `cmd.exe`) `/C`로 `"<command>"`를 raw 인자로 전달한다. 설정된 셸이 `/c` 인자를 가지면 같은 raw 방식이고 그 밖에는 일반 인자다. 비Windows는 `$SHELL`(없으면 `/bin/sh`) `-lc`다. 종료코드는 `output.status.code()`이고, 차단은 종료코드 2 또는 종료코드 0과 deny JSON이다(공식 hooks 문서). 공식 문서는 셸 종류를 명시하지 않는다. 로컬 설치는 codex-cli 0.162.1이며 이 원문이 그 버전과 같은지는 확인하지 않았다. Context7은 월 한도 초과로 조회하지 못해 공식 문서와 원문으로 대체했다.
- argv 목록 방식은 계속 막히지 않는다(guard 미실행). 현재 Codex 코드는 셸을 설정해도 `/c`가 있으면 raw로 보내므로 해당 경로를 확인하지 못했다. 이를 호출하는 래퍼가 있으면 보장하지 않는다.

## 검증 결과

삭제 문자열은 JSON 입력으로만 제공했고 실제 삭제·이동·정리는 실행하지 않았다.

| 검사 | 결과 |
|---|---|
| 셸별 실측(위 표) | 수정 명령은 cmd 2종·pwsh·powershell에서 허용 0·거부 2, 거부 사유 stderr 확인 |
| `python3 scripts/test_shell_guard.py` | 4 tests OK(exit 0). codex subfolder 회귀는 cmd shell=True·cmd /C raw·pwsh·powershell을 검사 |
| FullOps lint `--from 67e7de9` (구현 SHA 501d9d7, `make.cmd`→Linux 격리 clone) | exit 0, ERROR 0, WARNING 1(PLANS.md SIZE-001), product-lint·product-test passed |
| 같은 SHA Linux clone `make lint`·`make test` | 둘 다 exit 0 (HEAD 일치·clean 검사 포함) |

## 한계와 미검증

- 실제 Codex 세션의 hook 실행은 hook 신뢰 승인이 필요해 실행하지 않았다. 사용자 전역 `~/.codex` 설정(셸 override 여부)은 읽거나 바꾸지 않았다.
- 전역 Orca `codex-hook.cmd`는 변경하지 않았다. 플러그인 캐시와 이슈 오토모드도 건드리지 않았다.
- guard는 문자열 검사이며 완전한 보안 경계가 아니다(이전 phase의 한계와 같다).
