---
title: SAR-AUTOSETUP-001-O04-DEV 수정 검증 보고
status: draft
updated: 2026-10-28
owner: dev
tasks: [SAR-AUTOSETUP-001-O04-DEV]
summary: "O-04(호스트 플러그인 2개 이상에서 갱신 실패)를 Get-HostJson 한 곳에서 고쳤다. 0·1·2개 이상 합성 JSON 회귀 테스트가 PS5.1·PS7에서 통과하고, 수정 전 코드는 같은 오류로 실패한다. 예상 밖 오류에도 S6 안내를 더했고 버전을 0.1.1로 올렸다."
---

# SAR-AUTOSETUP-001-O04-DEV 수정 검증 보고

DEV 자체 검사다. 운영 서버·실제 `~/.claude`·`~/.codex`는 건드리지 않았다. 호스트 CLI는 stub이고 합성 JSON만 썼다. tester 재검증(실제 다중 플러그인 환경)과 운영 교체는 후속이다.

## 원인과 수정

- 원인: Windows PowerShell 5.1 `ConvertFrom-Json`은 JSON 배열을 배열 객체 한 개로 돌려준다. `Get-Registered`의 `@(Get-HostJson …) | Where-Object { $_.id -eq $Plugin }`는 이를 펼치지 못해 항목이 2개 이상이면 배열 전체가 통과하고 `Split-Path -Parent $entry.readFromFolder`가 `Path` null로 실패한다([QA2 O-04](../SAR-AUTOSETUP-001-QA-test/report.md)).
- 수정(`scripts/autosetup/knowslink-setup.ps1`): `Get-HostJson`이 `ConvertFrom-Json | ForEach-Object { $_ }`로 항목을 펼친다. 같은 패턴을 쓰는 곳은 모두 이 함수를 지난다: `Get-Registered`(Claude `plugin list`, Codex `plugin list --available`), 등록 확인(Claude `plugin list`, Codex `mcp list`). 호출부는 바꾸지 않았다. PS7은 원래 펼치므로 영향이 없다.
- S6: 예상 밖 오류는 `[ERROR]` 한 줄로 끝났다. 이제 `[F-10]` 실패 안내와 `이미 바뀐 것`·`다음 행동`을 보이고 설치 기록(`F-10`)에 단계를 남긴다. 동의 전이면 `바뀐 것: 없음`, 이후면 설치 기록 위치를 알린다. KO·EN 문구 `F10`·`F10n`·`cMaybe`를 더했다(D04 표의 F 형식 그대로, 화면 문구 변경은 아님).

## 회귀 테스트

`scripts/autosetup/test-get-registered.ps1`이 설치기의 `Norm`·`Get-HostJson`·`Get-Registered`를 AST로 꺼내 stub 호스트(.ps1)에 대해 실행한다. `scripts/autosetup/test_get_registered.py`가 `powershell.exe`·`pwsh`에서 돌리고 `make test`에 연결했다(셸이 없으면 SKIP).

| 단계 | 명령 | 종료코드 | 결과 |
|---|---|---|---|
| 수정 전(해당 한 줄만 원복) | `python3 scripts/autosetup/test_get_registered.py` | 1 | 0·1개 PASS, 2개 이상에서 `Split-Path : Cannot bind argument to parameter 'Path' because it is null.`(O-04 재현) |
| 수정 후 | 같은 명령 | 0 | `powershell.exe` 5.1.26100 10/10 PASS, `pwsh` 7.6.8 10/10 PASS |

케이스(Claude): 0개, knowslink 1개, 다른 플러그인 1개, 2개(knowslink 앞/뒤), 3개(`readFromFolder` null 포함), knowslink 없는 2개. (Codex): 0개, installed 1개, 여러 항목. 기대값은 knowslink 설치 폴더 경로 또는 `$null`이다.

## S6 확인

`LOCALAPPDATA`를 비워 동의 전에 예상 밖 오류를 만들고 설치기를 두 셸에서 실행했다(stub `claude.cmd`).

```
[F-10] 예상하지 못한 오류로 멈췄습니다: Cannot bind argument to parameter 'Path' because it is an empty string.
이미 바뀐 것: 없음
다음 행동: 같은 명령을 다시 실행하세요. 계속되면 위 오류 줄을 알려 주세요.
```

`powershell.exe`·`pwsh` 모두 `LASTEXITCODE=1`. 동의 뒤 단계(`cMaybe` 경로)는 코드 읽기로만 확인했고 실행하지 않았다.

## 버전 판단: 0.1.1

공개된 `/install/0.1.0/`은 URL·SHA-256이 고정되어 바꾸지 않는다(`package_plugin.py --release`는 기존 버전 폴더가 있으면 멈춘다). 그래서 수정본은 새 버전 0.1.1로 낸다. 설치기 버전은 플러그인 매니페스트 버전을 따르므로 함께 올렸다.

바꾼 곳: `adapters/.claude-plugin`·`.codex-plugin`·`.cursor-plugin`·`.grok-plugin`의 `plugin.json`, `adapters/manifest.json`, `adapters/src/mcp.ts`, `scripts/package_plugin.py`(Codex 매니페스트·패키지 package.json), `scripts/verify_grok_plugin.py`, `internal/relay/human_pages.go`(화면 명령 URL 숫자만, 문구 변경 없음), `README.md`, `adapters/README.md`. 화면 테스트(`autosetup_guide_test.go`)는 `plugin.json` 버전을 읽어 그대로 통과한다. 바꾸지 않은 곳: `install_files_test.go`·`check_public_ingress.py`의 `0.1.0`은 임의 semver를 쓰는 fixture/정규식 검사다. 0.1.0 사용자는 0.1.1 명령을 한 번 실행하면 갱신되고 에이전트·키·친구는 유지된다.

운영 교체 절차 후보는 [D12 ops-guide](../../../operations/ops-guide.md)의 '자동 설정 0.1.1 교체 후보' 절에 썼다(서버 접속 안 함).

## 제품 검사

| 명령 | 종료코드 |
|---|---|
| `python3 scripts/package_plugin.py --verify` | 0 |
| `python3 scripts/package_plugin.py --release <.devtools 작업 폴더>` (0.1.1 산출, 설치기 PS5.1·PS7 구문 오류 0) | 0 |
| `npm run check --prefix adapters` | 0 |
| `npm test --prefix adapters` | 0 |
| `go test ./cmd/... ./internal/...` (DB 없는 환경) | 0 |
| `go vet ./cmd/... ./internal/...`, `gofmt -l internal/relay/human_pages.go` | 0, 출력 없음 |
| `python3 scripts/check_compose.py`, `check_public_ingress.py`, `test_shell_guard.py` | 0, 0, 0 |

Windows에는 `make`가 없어 `make lint`·`make test`는 위 구성 명령으로 대신했다. `go test -race`는 쓰지 않았고(C 컴파일러 없음), `prettier --check`·`bash deploy/knowslink/test_public_config.sh`는 실행하지 않았다. `gofmt -l`이 건드리지 않은 `i18n_en_human.go`·`i18n_en_member.go`를 나열하는 것은 변경 전부터이며 이번 diff와 무관하다(Windows CRLF).

## 한계

- 실제 둘째 플러그인을 설치한 호스트에서 갱신을 실행한 것이 아니다. 합성 JSON과 stub이다. 실환경 재검증은 tester 몫이다.
- 테스트는 임시 폴더를 `%TEMP%`에 만들고 끝나면 지운다(검증 로그·도구는 `.devtools/`).

## 남은 일

1. tester: 호스트 플러그인이 2개 이상인 환경에서 같은 명령 재실행(갱신, AC-07·S5)과 실패 안내(S6) 재검증.
2. ops(사용자 승인 필요): D12 절 순서로 `/install/0.1.1/` 게시, 확인 뒤 회원 화면 명령을 0.1.1로 바꾸는 relay 배포.
