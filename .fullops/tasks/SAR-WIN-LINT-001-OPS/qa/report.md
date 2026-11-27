---
title: SAR-WIN-LINT-001-OPS 검증 보고
status: draft
updated: 2026-11-27
owner: ops
summary: Windows에서 fullops lint의 lint/test 실행 근거와 병합 시 필요한 PATH 접두
---

# SAR-WIN-LINT-001-OPS 검증 보고

## coor가 할 일 (이 과제 병합은 한 번만 PATH 필요)

`fullops lint`는 merge-base 시점의 `lint.json`을 쓴다. 이 과제를 병합할 때의 lint는 아직 `make`를 직접 부른다. 병합하는 PowerShell 세션에서만 PATH를 붙인 뒤 `fullops merge`를 실행한다. 사용자 PATH는 바꾸지 않는다.

```powershell
$t = 'D:\workspace\KnowsLink\.devtools'
$env:PATH = "$t\node-v22.22.2-win-x64;C:\Program Files\Git\usr\bin;$env:PATH;$t\w64devkit\bin"
```

병합 뒤 기준부터(SAR-ORCA-BOT-001-DEV 병합 포함)는 `python3 scripts/make.py`가 적용돼 PATH 준비가 필요 없다. Docker Desktop은 켜져 있어야 한다.

## 변경과 이유

| 변경 | 이유 |
|---|---|
| `scripts/make.py` 추가, `lint.json`의 `make X` → `python3 scripts/make.py X` | 사용자 PATH를 바꾸지 않고 그 프로세스에만 `.devtools` 도구를 붙인다. Windows가 아니면 `make`를 그대로 부른다. w64devkit은 PATH 끝에 둔다(busybox sh는 Git sh를 가리지 않게). |
| `test_public_config.sh`가 Windows(MINGW/MSYS)에서 `python:3.12-slim` 컨테이너로 재실행 | NTFS(noacl)에서 `stat -c %a`가 644라 0600 검사가 성립하지 않는다. assert를 빼지 않고 같은 스크립트를 Linux에서 돌린다. checkout은 읽기 전용, network none. |
| 작업 파일 CRLF → LF | `.gitattributes` 이전 체크아웃 잔존. |
| `project.md` `## 작업 검증`, 역할 메모 | 실행 방법 기록. |

검사 항목: `Makefile` diff 0. `lint.json` diff는 두 명령의 `run`뿐이며 같은 target을 부른다(`git diff a09df330 -- Makefile .fullops/config/lint.json`).

## CRLF 정리 기록

`.devtools/win-lint-001/lf.py`: `i/lf` + `w/crlf|mixed` 파일만, LF 변환 결과의 blob 해시가 index와 같을 때만 쓴다. dirty tree면 쓰지 않는다. 이후 `git add -u`로 stat만 갱신했다(staged 0, HEAD 불변). `checkout-index` 재체크아웃 대신 택했다. 삭제가 없다.

| 체크아웃 | 전 status | 변환 | 후 status / CRLF |
|---|---|---|---|
| main(`D:/workspace/KnowsLink`) | 0 | 1365 | 0 / 0 |
| designer, ops, tester | 0 | 각 26 | 0 / 0 |
| dev | 0 | 13 | 0 / 0 |
| coor | 1 (손대지 않음) | 대상 0 | - |

## 실행 결과 (2026-11-27, 코드 HEAD `1eedce84`)

새 설정(병합 뒤 상태), 추가 PATH 없음: `fullops lint --from 0f1be95a`

```
[passed] product-lint: python3 scripts/make.py lint (18.7초)
[passed] product-test: python3 scripts/make.py test (194.4초)
통과: head 1eedce842b53 / ERROR 0 / WARNING 1 / 실행 불가 0
```

옛 설정(이 과제 병합 때), 프로세스 PATH만: `fullops lint --from a09df330`

```
[passed] product-lint: make lint (18.2초)
[passed] product-test: make test (195초)
통과: head 1eedce842b53 / ERROR 0 / WARNING 2 / 실행 불가 0
```

WARNING은 LINT-001(설정 변경은 병합 뒤 적용)과 LINT-002(test 120초 초과)다.

- `go test -race`: w64devkit gcc 16.2.0으로 전 패키지 `ok`.
- 실패 경계: `beta.sh`의 `umask 077`을 `022`로 바꾼 사본에서 `test_public_config.sh` exit 1. 래퍼에 없는 target을 주면 exit 2.
- 병합 후보: `git merge-tree HEAD 73f3ac3e`(충돌 없음) export에서 `make.py lint`·`test` 모두 exit 0.
- main 체크아웃(병합 전 코드, 래퍼는 ops 워크트리 것): `lint-go lint-adapters`와 `test` exit 0.

## 검증하지 못한 것

- main 체크아웃의 `fullops lint` 전체 통과. 변경이 아직 main에 없다. 병합 때 확인된다.
- Linux에서 이번 HEAD 재실행. 셸 변경은 Windows 분기뿐이고 `python3 scripts/make.py`는 Linux에서 `make`를 그대로 부른다.

## 남은 자료

- 로그·사본: `D:/workspace/KnowsLink/.devtools/win-lint-001/`(`neg/`, `orca-merge/` 포함). 삭제하지 않았다.
- 레포 밖: 실수로 `lf.py` 사본을 `%TEMP%\lf.py`에 하나 만들었다. 기존 테스트가 `%TEMP%`에 `kl-*`, `shell-guard-ws-*` 자료를 남긴다.
- 새 워크트리는 `npm ci --prefix adapters`(Node 22)가 먼저다. 컨테이너는 `--rm`이라 남지 않는다.

## 되돌리기

두 커밋 `0f1be95a`, `1eedce84`를 `git revert`한다. LF 작업 파일은 index와 같은 내용이라 되돌릴 필요가 없다.
