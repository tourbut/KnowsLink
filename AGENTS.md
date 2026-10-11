<!-- fullops-squad:start -->
## FullOps Squad

이 레포는 `.fullops-squad/fullops.json`이 있을 때만 FullOps 하네스를 사용한다.
작업 시작 시 `.fullops-squad/FULLOPS.md`를 직접 열고 관련 규약을 따른다.
<!-- fullops-squad:end -->

## Windows 파일 작업 안전

Windows 파일 정리는 PowerShell의 `-LiteralPath`와 확인된 절대경로를 사용한다. 정리 대상의 실제 경로와 workspace 내부 여부를 먼저 확인한다. 다른 셸로 삭제 문자열을 전달하지 않는다.
재발 방지 실행 차단의 독립 검수가 끝날 때까지 이 저장소의 worker는 삭제·이동 정리 작업을 실행하지 않는다. 임시 검증 자료는 보존하고 완료 보고에 경로를 남긴다.
