<!-- fullops-squad:start -->
## FullOps Squad

이 레포는 `.fullops-squad/fullops.json`이 있을 때만 FullOps 하네스를 사용한다.
작업 시작 시 `.fullops-squad/FULLOPS.md`를 직접 열고 관련 규약을 따른다.
<!-- fullops-squad:end -->

## Windows 파일 작업 안전

Windows 파일 정리는 PowerShell의 `-LiteralPath`와 확인된 절대경로를 사용한다. 정리 대상의 실제 경로와 workspace 내부 여부를 먼저 확인한다. 다른 셸로 삭제 문자열을 전달하지 않는다.
재발 방지 실행 차단의 독립 검수가 끝날 때까지 이 저장소의 worker는 삭제·이동 정리 작업을 실행하지 않는다. 임시 검증 자료는 보존하고 완료 보고에 경로를 남긴다.

실행 전 차단은 `scripts/shell_guard.py` PreToolUse hook이 담당한다. Claude Code와 Grok은 `.claude/settings.json`을, Codex는 `.codex/hooks.json`을 읽는다. Grok은 폴더 신뢰, Codex는 hook 신뢰 승인이 있어야 hook을 실행한다.
Codex hook은 session cwd에서 실행되므로 `git rev-parse --show-toplevel`로 찾은 루트의 guard를 호출한다. Codex는 Windows에서 `%COMSPEC% /C "<command>"`를 raw 인자로 실행하고(셸을 따로 설정하면 그 셸), 이 명령은 cmd·PowerShell(`pwsh`/`powershell -Command`)·sh에서 같은 종료코드를 낸다. 끝의 `; exit $LASTEXITCODE`는 PowerShell이 native 종료코드 2를 1로 바꿔 차단이 풀리는 것을 막는다. 따옴표 없는 형태는 PowerShell 파서 오류(code 1)를 내고, argv 목록으로 `cmd /C`를 부르면 따옴표가 `\"`로 바뀌어 guard가 실행되지 않으므로 이 명령을 그렇게 호출하지 않는다. guard가 실행되지 않으면 Claude Code·Grok은 셸 호출을 그대로 실행한다. Codex는 종료코드 2만 차단하고 다른 실패는 비차단 오류로 처리한다.
hook은 다른 셸·평가기(`bash -c`, `cmd /c`, `pwsh -Command`, `wsl`, `ssh`, `Invoke-Expression` 등)로 전달한 삭제 문자열을 거부한다. 대상이 검증되지 않은 재귀 삭제, 변수·와일드카드·빈 값·루트 대상 삭제, `find -delete`, `git clean`(`-n` 제외), `robocopy /MIR`, 인라인 코드의 삭제 API도 거부한다.
허용하는 재귀 삭제는 단독 PowerShell 명령 `Remove-Item -LiteralPath '<workspace 안의 실제 절대경로>' -Recurse -Force` 하나다. 경로는 존재해야 한다. 경로는 workspace 자체·드라이브 루트·`.git`·링크·reparse 지점이 아니어야 한다.
hook은 명령 문자열 검사이며 완전한 보안 경계가 아니다. 스크립트 파일·별칭·생성 코드는 hook을 우회할 수 있다. 따라서 위 규칙을 계속 지킨다. 명령이 거부되면 다른 형태로 우회하지 않고 coordinator에게 보고한다.

## 작업 파일 위치

레포 작업에 쓰는 도구·임시 산출물·검증 로그는 기본 체크아웃의 `.devtools/`(gitignore)에 둔다. 경로는 `git rev-parse --path-format=absolute --git-common-dir`의 상위 폴더 아래 `.devtools/`다(현재 `D:/workspace/KnowsLink/.devtools`). 모든 워크트리가 이 폴더를 공유한다.
- 도구: `.devtools/node-v22.22.2-win-x64`, `.devtools/w64devkit`, `.devtools/ssh`, `.devtools/ssh-python`.
- 과제 로그·bundle·캡처: `.devtools/<과제 폴더>/`. 레포 증거로 남길 것은 기존 규칙대로 `qa-reports/`에 커밋한다.

레포 밖(`%LOCALAPPDATA%`, `%TEMP%`, 홈 디렉터리)에 이 레포 전용 파일을 만들지 않는다. 호스트 세션 scratchpad는 세션 한정 임시 파일에만 쓴다. 레포 밖이 꼭 필요하면 사용자 승인 뒤 경로·용도·정리 시점을 이 절에 기록한다.
예전 경로 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/<X>`는 `.devtools/<X>`로 옮겼다. 과거 기록의 예전 경로는 이 대응으로 읽는다.
