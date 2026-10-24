<!-- fullops-squad:start -->
## FullOps Squad

이 레포는 `.fullops-squad/fullops.json`이 있을 때만 FullOps 하네스를 사용한다.
작업 시작 시 `.fullops-squad/FULLOPS.md`를 직접 열고 관련 규약을 따른다.
<!-- fullops-squad:end -->

## Windows 셸 안전

삭제·이동 정리는 [AGENTS.md](AGENTS.md)의 `Windows 파일 작업 안전` 절을 따른다. `.claude/settings.json`의 `scripts/shell_guard.py` hook이 교차 셸 삭제와 검증되지 않은 재귀 삭제를 실행 전에 거부한다.
도구·임시 산출물·검증 로그 위치는 [AGENTS.md](AGENTS.md)의 `작업 파일 위치` 절을 따른다(레포 밖 생성 금지).
