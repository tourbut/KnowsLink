---
title: SAR-WINDOWS-SHELL-GUARD-001-REVIEW 리뷰
status: draft
updated: 2026-10-11
owner: tester
tasks: [SAR-WINDOWS-SHELL-GUARD-001-TESTER]
summary: Windows 삭제 실행 차단 4e82b4f의 독립 고정 SHA 리뷰와 미해결 medium 4건
---

# SAR-WINDOWS-SHELL-GUARD-001-REVIEW 리뷰

- 검토자 / CLI / 모델: tester / claude / claude-sonnet-5-5 high. 세션 `0ef9c64c-15ec-4fcc-a3be-c11022475dc1`.
- base SHA / head SHA / merge-base: `a90a9e60ed1f009627aa39031fbf3c228d0d7599` / `4e82b4f4cd1bc31a59d0d230349fd34d95c7f0c7` / `a90a9e60ed1f009627aa39031fbf3c228d0d7599`. 지시서 base `0e4f638`은 이 범위 안의 준비 커밋이다. DEV `lint.json`이 a90a9e6 기준이라 같은 기준으로 맞췄다.
- OCR 버전 / 적용 규칙: `result.json`의 `ocr_version`·`rule_sha256`. 규칙은 `.fullops-squad/review/rule.json`이며 변경하지 않았다.
- 요구사항·완료 기준 원천: `.fullops-squad/handovers/to_tester.md`(보존 후 logs), `SAR-WINDOWS-SHELL-GUARD-001-DEV.md`, DEV 지시서 사본(`logs/2026-10-11_to_dev.md`).
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 19 파일 / 모두 `result.json` 기록 / reviewed 12 · skipped 7(생성 JSON 탐색 산출물 6, 무관한 coor lint 증거 1).
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: 0 / 5 / 0. DEV 보존 `final2-lint.json`을 복사했다(head 4e82b4f, attempt가 DEV `pass.json`과 일치, `product-lint`·`product-test` exit 0). tester가 같은 명령을 다시 실행하지 않았다(변경·근거 결함 없음, 로컬 Docker 금지, Linux 증거 재사용).
- SIZE-002: 지시서 예상 규모 기록 없음 / 추가 508줄(상한 400) / exec plan·packet-outcomes 문서와 회귀 테스트 포함. 코드는 guard 220줄·테스트 146줄로 각각 300줄 이내다. 분할하지 않은 이유는 단일 파일 구현이다. 수용.
- DEP-001: 의존성 선언 파일 변경 없음(`package.json`·`go.mod` 불변). 새 의존성 없음.
- UI 디자인: 해당 없음. UI 변경이 없다.

## 검토 범위

`result.json`의 모든 (path, status)에 reviewed 또는 skipped와 사유를 기록했다. 제품 코드·hook 설정·테스트·안내 문서는 snapshot에서 직접 읽었다. 생성 JSON 탐색 산출물은 `packet`·`packet-outcomes` 검사로 갈음했다. 사고 원문 명령·비밀값·사고 USN/VSS 자료는 diff에 없다(설명 문장과 비공개 폴더 이름만 있다). 이 리뷰에서도 원문을 읽지 않았다.

## 발견 사항

critical 0 · high 0 · medium 4 · low 5. 모두 미해결(`resolved: false`)이며 수정은 DEV 후속 담당이다. tester는 제품 코드를 수정하지 않았다. 재현 입력은 guard 프로세스의 stdin에 보낸 합성 payload이고 어떤 셸에도 전달하지 않았다. 증거: `SAR-WINDOWS-SHELL-GUARD-001-TESTER-test/probes/`.

| 번호 | 심각도 | 파일·줄 | 요약 | 수정 제안 |
|---|---|---|---|---|
| F1 | medium | `scripts/shell_guard.py:98` | `Remove-Item -r:$true C:\Users\x\Documents`가 allow다. `-r`은 `Remove-Item`의 유일한 `R*` 매개변수 접두사이고 `-r:$true`는 `recurse()`의 `-r`·`-rec*`·`-[a-z]{1,4}` 판정에 모두 맞지 않는다. 변수·와일드카드·루트 대상은 같은 약어에서도 deny라 사고 형태는 막힌다 | 콜론 값을 허용하는 정규식으로 `-r` 계열을 판정하고 회귀 deny 추가 |
| F2 | medium | `scripts/shell_guard.py:30` | `python -c "…os.system('rm -rf …')"`, `subprocess`, `node -e` `execSync`·`fs.rmdirSync`·`unlinkSync`가 allow다. `PROGRAM_DELETE`에 없다 | `INTERPRETER`와 `DELETE_WORD` 동시 존재 시 거부 |
| F3 | medium | `scripts/shell_guard.py:179` | `git -C . clean -fdx`, `git --no-pager clean -fdx`가 allow다. `args[0]`만 본다 | 옵션을 건너뛴 하위 명령 판정 |
| F4 | medium | `.codex/hooks.json:9` | 상대경로 hook은 하위 폴더 Codex 세션에서 python 종료코드 2(실측)로 모든 셸 호출을 막는다. Codex 문서는 exit 2를 차단으로 정의한다. DEV 문서·AGENTS의 "fail-open" 서술이 사실과 다르다. 삭제 우회가 아니라 가용성 결함이다 | git root 기준 경로를 쓰는 래퍼와 서술 정정 |
| F5 | low | `scripts/shell_guard.py:41` | tokenizer가 bash 이스케이프로 PowerShell을 읽는다. `Get-ChildItem "C:\d\"; Remove-Item -Recurse C:\x # "`가 allow다(고의적 구성) | PowerShell 도구 전용 규칙 또는 따옴표 안 삭제 단어 거부 |
| F6 | low | `scripts/shell_guard.py:20` | `robocopy /MOVE`, `rsync --delete`, `docker run … rm`, `busybox rm`, `. Remove-Item`, scriptblock Create, `r\m` 미감지 | 필요 시 목록 추가 |
| F7 | low | `scripts/shell_guard.py:166` | 오탐. 삭제 단어와 `sh` 단어가 한 명령에 있으면 거부한다. 검수 중 heredoc probe 작성 호출 2회가 실제로 막혔다. 허용 형식에 `-WhatIf`·`-Confirm:$false`가 있으면 거부, `git clean -fdn` 거부 | `-WhatIf` 허용·결합 `-n` 판정 |
| F8 | low | `scripts/shell_guard.py:199` | 2MB 단일 단어 명령은 10초 timeout을 넘겨 fail-open이다(400KB는 1.76초) | 누적 방식 개선 또는 상한 |
| F9 | low | `AGENTS.md:13` | 서술되지 않은 한계: workspace 밖 비재귀 삭제 허용, `Move-Item`·`mv` 미검사, 허용 디렉터리 안 중첩 reparse·TOCTOU 미검증, 프로젝트 체크아웃 밖·hook 없는 클라이언트(Gemini·agy)·`apply_patch`·MCP 삭제 미보호 | 사용자 결정: 전역·managed 강제, 이동 금지 반영 |

알려진 위험의 판정: ① fail-open(python3 미해결·timeout)은 세 호스트 공통 플랫폼 한계이며 DEV가 밝혔다. 이 세션에서 `python3`이 해석되고 hook이 실제로 차단하는 것을 확인했다. 결함이 아니라 수락된 한계다(F8은 입력 크기 부분). ② Codex 하위 cwd는 DEV 서술과 달리 fail-closed다(F4, 가용성 결함). ③ 스크립트 파일·별칭·`make` 우회는 문자열 검사의 구조적 한계이며 문서가 밝혔다. 수락된 한계다. ④ TOCTOU와 중첩 reparse는 실제 삭제 금지로 검증하지 못했고 F9에 남겼다. ⑤ 오탐은 F7이다.

## 검증 및 남은 제약

- 실행: DEV 회귀를 Windows snapshot에서 재실행(exit 0, 3 tests OK). 독립 probe 약 120건과 PowerShell 전용 probe 15건, 입력 크기 측정. 실제 Claude Bash·PowerShell 도구 차단을 이 세션에서 직접 확인. `grok inspect --json`으로 hook 발견·신뢰·활성 확인. 모든 probe 뒤 대상 경로가 그대로 있다.
- 재사용: 등록 `product-lint`·`product-test`는 DEV 같은 SHA 근거(pass.json·lint.json·Linux 로그, HEAD 일치·clean)를 재사용했다. tester 실행으로 기록하지 않는다.
- 미실행: 실제 Grok·Codex 모델 tool call, Codex 사용자 hook 신뢰, 실제 삭제·중첩 reparse 삭제, Linux pwsh 경로, 전역·managed hook. Context7은 월 한도 초과로 실패해 Codex 의미는 공식 문서로 확인했다.
- 사용자 단계로 남는 일: 기존 세션 재시작(Claude·Grok `/hooks` Reload), 체크아웃별 Codex hook 신뢰, 다른 역할 체크아웃의 main 동기화.

## 검토 결론

수락 가능(조건부). 미해결 critical·high는 없다. 사고 형태와 변수·빈 값·와일드카드·루트 대상의 재귀 삭제는 실제 Claude 두 도구에서 실행 전 차단이 확인됐다. F1~F3은 작은 수정으로 막을 수 있고 F4는 Codex 사용 전 수정이 필요하다. 병합을 막는 항목은 아니며 후속 DEV 과제를 권고한다. Grok·Codex의 실제 모델 tool call 차단은 확인하지 못했으므로 적용 완료로 표시하지 않는다. check 통과는 AI 검토의 내용을 보증하지 않는다.

## 독립성과 인계 점검

- 구현자 실제 세션 ID / 검토자 별도 실제 세션 ID: `bacd5d1f-de07-43b3-9efa-0df5636729c7` / `0ef9c64c-15ec-4fcc-a3be-c11022475dc1`.
- 관리된 clean detached snapshot 경로 / 고정 40자리 SHA / 읽기 전용 확인: `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-4ff90f54195a4728bc150bd64f431020` / `4e82b4f4cd1bc31a59d0d230349fd34d95c7f0c7` / detached HEAD, `git status --porcelain` 비어 있음, `review.py snapshot` 기록 `fullops-snapshots/SAR-WINDOWS-SHELL-GUARD-001-REVIEW.json`.
- 현재 attempt·packet의 path/category별 처리 근거 / partial·unknown 해소: 원 구현 packet(`…-DEV-packet.json`)과 `…-DEV-packet-outcomes.json`을 `review.py check --task-key SAR-WINDOWS-SHELL-GUARD-001-DEV`가 검증한다. tester 현재 packet(239항목, partial·unknown 8)은 `…-TESTER-packet-outcomes.json`에 항목별 처리와 `uncertainty_review`를 남겼다.
- 대화 미참조 점검: 현재 요구는 `handovers/to_tester.md`와 DEV exec plan의 "기술 계획과 범위", 결정 이유는 같은 문서의 사고 원인 절, 구조·구현은 `scripts/shell_guard.py`와 hook JSON 두 파일, 실행·검증은 exec plan의 "검증 결과"와 DEV 보존 lint·test 로그, 운영·복구는 "활성화 상태와 재시작"과 "보존한 임시 증거", 다음 작업은 PLANS.md 마지막 항목에서 찾았다. 모두 확인했다. 복구 자료 위치는 비공개라 해당 없음(읽지 않음). 깨진 링크 없음(CLAUDE.md→AGENTS.md "Windows 파일 작업 안전" 절 존재 확인). 오래된 정보는 Codex 하위 cwd 서술(F4)이다.
- 테스트 레벨 / 필수 검사와 skipped 사유 / 원래 SHA의 재사용 증거: lite. `product-lint`·`product-test`는 skipped가 아니라 DEV 실제 증거 재사용(위 lint 항목). 독립 guard 경계 QA는 직접 실행.
- 검토·증거 보존·reviewer release·관리된 snapshot cleanup 결과: 증거를 이 폴더와 `SAR-WINDOWS-SHELL-GUARD-001-TESTER-test/`에 보존했다. 지시에 따라 snapshot과 임시 자료를 정리하지 않았다. cleanup과 reviewer release는 부모(coor) 담당이다.

historical check는 보존 증거 확인이며 신규 SHA의 수락 리뷰를 대신하지 않는다.
