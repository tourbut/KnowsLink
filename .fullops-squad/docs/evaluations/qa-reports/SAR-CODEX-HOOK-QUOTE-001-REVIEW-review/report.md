# SAR-CODEX-HOOK-QUOTE-001-REVIEW 리뷰

- 검토자 / CLI / 모델: tester / Claude Code / Sonnet 5.5. 실제 reviewer session `70c06551-37d6-4603-908a-6956e4a2ea2e`. 구현 session `cf470ffe-7c1a-4b79-8138-cecfc8b9744a`와 다르다.
- base SHA / head SHA / merge-base: `67e7de9287c6c2ac6bb92beb5b31b523f7e64ec6` / `d7d6d744d5b0010a384a7389ec8f64ec7bb93336` / base와 같다(gate pass.json 확인).
- OCR 버전 / 적용 규칙: open-code-review v1.12.13 delegate(preview·rules만 사용, LLM 호출 없음). `review/rule.json`의 sha256 `f1061481…4de8`은 prepare 시점과 같다. 공통 기준 `fullops-common-0.3.3`(`rules/common/README.md`, coding-style, testing, security)과 `project.md`, `AGENTS.md`를 읽었다. 예외 없음. 읽은 규칙 버전은 tester 체크아웃(HEAD `079ee84`, packet 기준 `b63e877`)의 파일이다.
- 요구사항·완료 기준 원천: `.fullops-squad/handovers/to_tester.md`(attempt `3551964c…`), DEV 지시서·phase `SAR-CODEX-HOOK-QUOTE-001-DEV`.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 16 / 16 / 0 / 10 / 6. 6건 skipped는 모두 생성된 Jev 탐색 증거(DEV context·documents-find·find·packet·route·route 이력)이며 아래에 사유를 적었다.
- snapshot: `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-d5965114daa146b4add9d1dae27ef08a` — HEAD `d7d6d744…`, detached, 작업 트리 clean. `review.py snapshot`으로 생성했다. 테스트 실행 때 `PYTHONDONTWRITEBYTECODE=1`을 사용했고 실행 뒤에도 clean이다.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: ERROR 0 / WARNING 1(`PLANS.md` SIZE-001, 1600→1607줄, coor 관리 파일) / 실행 불가 0. `product-lint`·`product-test`는 DEV의 exact-SHA 결과를 복사·재사용했다(아래 gate 절).
- SIZE-002: 해당 없음(추가 80줄 < 상한 400).
- DEP-001: 의존성 선언 파일 변경 없음.
- UI 디자인: 해당 없음. UI·테마·토큰·공용 컴포넌트 변경이 없고 디자인 lint 대상 파일이 없다.

## 검토 범위

result.json의 16개 `(path, status)`를 모두 처리했다. 제품 변경은 `.codex/hooks.json`, `scripts/test_shell_guard.py`, `AGENTS.md` 3개뿐이다. 나머지는 phase·인계·상태 문서와 생성된 Jev 증거다. 구현 커밋 `501d9d72`와 최종 기록 `d7d6d744`의 차이는 `.fullops-squad` 아래 3개 파일(phase 문구·로그 보존·인박스 비움)뿐이다. 제품 경로 diff는 비어 있다(`git diff --stat 501d9d72 d7d6d744`).

제외 증거(`qa-reports/`)는 이번 diff에 없다. 이 리뷰의 독립 QA 증거 `evidence/`는 위 폴더에 둔다.

## 발견 사항

| # | 심각도 | 파일 | 내용 | 상태 |
|---|---|---|---|---|
| 1 | low | `.codex/hooks.json` | argv 목록으로 `cmd.exe /C <command>`를 호출하면 guard가 실행되지 않고 허용0·거부0이다. DEV가 문서화한 잔여 한계다. 현재 Codex 경로는 raw 인자라 영향이 확인되지 않았다. | 미해결(low), 수정 요청 없음 |

새 critical/high는 없다.

### 원 실패·미수락 후보·최종 구현의 대조 (packet conflict 판정)

packet은 보존된 DEV 로그를 `conflict: contradicts task`로 표시했다. 세 단계를 구분하면 모순이 아니다. 같은 입력(`git status` 허용, 삭제 문자열은 guard stdin JSON으로만 제공한 거부)을 snapshot의 `scripts/`에서 직접 실행했다(`evidence/indep_hist.out`).

| 명령 | cmd /C raw | pwsh | powershell | argv 목록 `cmd /C` |
|---|---|---|---|---|
| 원 config 67e7de9(따옴표 없음) | 0 / 2 | **parser error, 허용 1 / 거부 1** | **parser error, 1 / 1** | 0 / 2 |
| 미수락 후보 7cae8d7(따옴표만) | 0 / 2 | 0 / **1** | 0 / **1** | 0 / **0**(guard 미실행) |
| 최종 d7d6d744(`; exit $LASTEXITCODE` 추가) | 0 / 2 | 0 / 2 | 0 / 2 | 0 / **0**(guard 미실행) |

- 원 실패 `Hook failed code 1`과 일치하는 현상은 원 명령의 PowerShell parser error(종료 1)다. 호스트가 낸 원 stderr는 확보하지 못했다. tester도 확보하지 않았고 추정하지 않는다.
- 후보의 `expected 2 / actual 0` 회귀는 argv 목록 호출에서 재현된다. 후보는 pwsh에서 native 종료코드 2가 1로 바뀌어 Codex가 차단하지 않는 별도 결함도 있다(거부 1).
- 최종 수정은 두 결함을 막는다. 다만 argv 목록 호출은 여전히 0/0이며 이것이 발견 사항 1이다.
- 이전 phase `SAR-WINDOWS-SHELL-GUARD-001-DEV`의 "`cmd /C` 전용" 전제는 PowerShell parser error 재현과 충돌했고, DEV가 정정 주석으로 고쳤다. 원문은 보존됐다.

## 검증 및 남은 제약

### 독립 확인(좁은 보안 경계, lite)

실제 지원 셸 실행 방식으로 snapshot에서 실행했다. 코드를 문자열 리터럴로 넘겨 guard가 실행되지 않은 채 0이 되는 경우는 PASS로 취급하지 않았다. 허용 칸이 0이어도 거부 칸이 2이고 deny JSON을 파싱해 guard 실행을 증명했다.

- `evidence/indep_hook.py` → `indep_hook.out`: cwd root·`scripts`·`adapters` × {cmd `shell=True`, cmd `/C` raw, pwsh 7.6.6, Windows PowerShell 5.1}: **허용 0 · 거부 2**, stdout에 deny JSON과 stderr `Blocked by scripts/shell_guard.py` 확인. 12개 조합 모두 같다.
- `evidence/deny_json_check.py`: `cmd /C` raw와 pwsh에서 종료 2 + `hookSpecificOutput.permissionDecision == "deny"`, `hookEventName == "PreToolUse"` 단정 통과.
- Git Bash `sh -c`: 허용 0 · 거부 2.
- 삭제 문자열은 JSON 입력으로만 전달했고 실제 삭제·이동·정리는 실행하지 않았다. 임시 workspace `%TEMP%\shell-guard-ws-*`는 테스트가 만든 증거로 보존했다.
- `python3 -B scripts/test_shell_guard.py`를 snapshot에서 실행: 4 tests OK, 종료 0(약 17초).
- PowerShell의 native nonzero 변환은 후보 7cae8d7의 pwsh 거부 1로 직접 재현했고, 최종의 `; exit $LASTEXITCODE`로 2가 유지되는 것을 확인했다.

### gate 재사용(DEV exact SHA)

- gate 정본 `D:/workspace/KnowsLink/.git/worktrees/dev/fullops-gate/pass.json`: head `d7d6d744…`, base `67e7de9…`, merge_base 같음, attempt `82a809fa…`, `config_sha256`가 이 리뷰의 lint 설정 digest와 일치한다. `lint-attempt.json`의 attempt와도 같다.
- Linux exact SHA 로그 `KnowsLinkDevTools/wsg001/d7d6d744…-lint.log`, `-test.log`: 둘 다 끝에 `EXIT=0`. lint 로그의 `ERROR: …` 줄은 설정 렌더 음성 검사의 기대 출력이며 종료 0이다. test 로그는 `test_shell_guard.py` 4 tests OK를 포함한다.
- 재사용 근거: 제품 경로 diff가 `501d9d72`와 `d7d6d744` 사이에 없고(위), `.fullops-squad/lint/lint.json`과 `fullops.json`의 merge-base 설정이 gate의 `config_sha256`과 같다. `lint.json`의 `product-lint`·`product-test`는 로그 경로·sha256·gate attempt를 `reason`에 적어 복사했다. tester가 새로 실행한 결과가 아니다.
- 재실행하지 않은 이유: 로컬에 `make`가 없고, 등록 `make.cmd`는 DEV worktree HEAD를 원격 SSH 서버에 올리는 도구라 이 snapshot을 검증하지 않으며 자격 정보가 필요하다. 외부 발송이므로 실행하지 않았다.
- `lint.json`의 violation·SIZE·DEP 결과는 snapshot에서 `lint.py`를 직접 실행해 얻었다(명령 2건만 실행 불가였다가 위 증거로 채움).

### 전역 설정(읽기 전용 확인)

`~/.codex/config.toml`에는 셸 override 키가 없다. hook 신뢰 항목은 원본 `D:\workspace\KnowsLink\.codex\hooks.json:pre_tool_use:0:0`만 있고 coor·tester 워크트리 경로 항목은 없다. 전역 신뢰 승인·캐시는 변경하지 않았다.

### 미검증·한계

- 실제 Codex host에서 hook 실행과 오류 소멸은 확인하지 않았다(hook 신뢰 승인 필요, 승인하지 않음). 검사 통과는 사용자 host의 오류 소멸을 의미하지 않는다.
- 원 host stderr와 현재 native 세션 소멸 원인은 확보·추정하지 않았다.
- Codex의 차단 의미(종료코드 2/deny JSON)와 셸 선택 규칙은 DEV가 인용한 Codex 원문에 의존한다. tester는 그 원문을 재확인하지 않았다. 로컬 codex-cli는 DEV 기록상 0.162.1이다.
- 셸 override가 PowerShell이 아닌 다른 방식으로 명령을 전달하는 경로는 검증하지 않았다.
- guard는 문자열 검사이며 완전한 보안 경계가 아니다(이전 phase와 같은 한계).
- 등록 `kind: test` 명령은 위 재사용 증거로 기록했다. `work.py finish`와 `review.py check`는 테스트 성공을 보증하지 않는다.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과(확인/미확인/해당 없음) | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `handovers/logs/2026-10-12_to_dev.md`, `docs/exec-plans/phases/SAR-CODEX-HOOK-QUOTE-001-DEV.md` "원인과 실행 계약" | 확인 | 누락 없음 |
| 구조와 구현/미완료 상태 | 위 phase "기술 계획과 범위"·"한계와 미검증", `AGENTS.md` Codex hook 절 | 확인 | argv 목록 호출·실제 host 실행은 미완료로 보존 |
| 실행·검증 방법과 증거 | 이 폴더 `evidence/`, `lint.json`, `docs/exec-plans/phases/SAR-CODEX-HOOK-QUOTE-001-TESTER.md` | 확인 | gate 로그는 레포 밖 `KnowsLinkDevTools`(보존 필요) |
| 운영·복구(해당 시) | hook 롤백은 `.codex/hooks.json` revert. 보호 규칙 끄기 없음 | 해당 없음 | 배포·운영 변경 없음 |
| 다음 작업·담당·재개 조건 | coor: main 병합·push·역할 동기화, 사용자 host에서 trust 승인 뒤 실제 확인 | 확인 | PLANS.md는 coor 책임 |
| 로컬 링크·절 접근/지원 한계 | phase의 상대 링크 `SAR-CODEX-HOOK-QUOTE-001-DEV.md`, AGENTS.md | 확인 | 깨진 링크 없음 |
| snapshot 정리 후 정본 접근 | review 폴더(result·report·lint·evidence)는 tester 체크아웃에 커밋. snapshot은 정리하지 않고 보존 | 확인 | cleanup은 coor가 reviewer release 뒤 `review.py cleanup` |

## 검토 결론

수락 가능하다. 새 critical/high는 없다. 원 실패(원 명령의 PowerShell parser error), 미수락 후보(argv 목록 회귀와 pwsh 종료코드 변환)와 최종 수정의 원인을 독립 재현으로 구분했고, 지원 셸 4종과 sh에서 허용 0·거부 2·deny JSON을 root·scripts·adapters에서 확인했다. DEV exact-SHA gate와 product-lint/test exit 0은 코드·설정 동일성을 근거로 재사용했다. low 1건과 실제 Codex host 미검증은 보존한다. check 통과는 AI 검토의 내용·테스트 성공을 자동 보증하지 않는다.
