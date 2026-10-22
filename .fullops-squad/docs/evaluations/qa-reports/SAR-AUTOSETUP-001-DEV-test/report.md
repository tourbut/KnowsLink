---
title: SAR-AUTOSETUP-001-DEV 자체 검사 보고
status: draft
updated: 2026-10-23
owner: dev
tasks: [SAR-AUTOSETUP-001-DEV]
summary: "Windows 11 자동 설정 설치기를 격리 relay·합성 회원·격리 호스트 설정에서 실측했다. AC-05–AC-10과 F-01–F-09 KO·EN을 확인했고 AC-01 N7–N8·AC-02–04는 tester 후속이다."
---

# SAR-AUTOSETUP-001-DEV 자체 검사 보고

DEV 자체 검사다. 기술 계획은 [실행 기록](../../../exec-plans/phases/SAR-AUTOSETUP-001-DEV.md)이다. 운영 relay·실계정·실제 상대·사용자 실제 `~/.claude`·`~/.codex`는 쓰지 않았다.

- 증거 SHA: 제품 `4bc19f9b`(기준 main `3c8007e1`). EN·pick은 fixture만 고친 `bf2cb34e`에서 같은 배포물로 실행했다(`git diff 4bc19f9b bf2cb34e`는 `autosetup_fixture_test.go`뿐).
- 환경: Windows 11 Pro 26300, Windows PowerShell 5.1로 실제 명령 형태 `& ([scriptblock]::Create((irm <fixture>/install/0.1.0/knowslink-setup.ps1))) -App … -Lang …` 실행. Claude Code `2.1.308`(검증 버전과 달라 매번 AS-09 동의), Codex CLI `0.163.0`.
- relay: `internal/relay/autosetup_fixture_test.go` 격리 fixture(합성 Google 회원·친구 1명, 브라우저 동의·취소·응답 유실 제어). 배포물은 `package_plugin.py --release --relay-url http://127.0.0.1:18411`(설치기 `c702ca30…`, Claude ZIP `9361d4f1…`, Codex ZIP `dbea5ece…`).
- 시나리오마다 `LOCALAPPDATA`·`CLAUDE_CONFIG_DIR`·`CODEX_HOME`·`TEMP`를 새 폴더로 둔다. 도구: [drive.py](drive.py), [scenarios.py](scenarios.py). 판정 결과: [results-4bc19f9b.out](results-4bc19f9b.out)(main·KO 42 PASS, 뒤의 EN 부분은 아래 한계로 중단·FAIL 1), [results-bf2cb34e.out](results-bf2cb34e.out)(EN·pick, 18 PASS, exit 0).

| 기준 | 결과 | 근거 |
|---|---|---|
| AC-01 | 부분 PASS | N1–N6 사람 입력: 명령 붙여 넣기 1, (2.1.308만) 계속 `y` 1, 계획 동의 `y` 1, 브라우저 동의 1, 호스트 재시작 1. 경로·환경 변수·설정 편집 입력 0. N7–N8은 tester |
| AC-02·AC-03 | 미실행(tester) | 실제 세션 발송 승인은 사람 입력 필요. 등록된 MCP command·args·env가 `install()` 계약과 같음을 `claude plugin list`·`codex mcp list`로 확인 |
| AC-04 | 미실행(tester) | 설치기 구간만: Claude 신규 17.3초(Node 실제 내려받기 포함), Codex 신규 7.6초, 갱신 4.6·2.3초 |
| AC-05 | PASS | 신규·갱신 전후 다른 agent 폴더·Grok 설정·호스트 인증 파일·다른 MCP 항목·settings/config.toml(KnowsLink 항목 제외) 해시 동일, 두 호스트 |
| AC-06 | PASS | 고아 `<host>-01-private`를 두면 `<host>-02` 사용·원본 불변. credential·개인키가 출력·설치 기록에 없음 |
| AC-07 | PASS | 같은 명령 재실행이 Google 없이 갱신. agent 폴더 바이트 동일(ID·키), 웹 agent 수·친구 수 불변, 이전 설치 `.previous-*` 보존 |
| AC-08 | PASS | F-01–F-09 각각 KO·EN에서 `[F-xx]`·이미 바뀐 것·다음 행동과 '바뀐 것' 판정(아래) |
| AC-09 | PASS | F-03–F-08 뒤 반대 언어로 같은 명령 재실행이 정리 없이 완료. F-08 뒤 이전 설치 바이트 동일·등록 유지(실제 세션 수신은 tester) |
| AC-10 | PASS | `autosetup_guide_test.go`: 명령은 claude·codex에만, `-Lang`은 화면 언어, 버전은 plugin.json과 같음, D04 KO/EN 문장 대조, 6번 명령 없음 |

| 실패 | 주입 | 확인한 '바뀐 것' |
|---|---|---|
| F-01 | PATH에 호스트 없음 | 설치 폴더 없음 |
| F-02 / 계획 거절 | 미검증 버전에 `n` / 계획에 `n` | 설치 폴더 없음 |
| F-03 / F-04 | 배포 ZIP 404 / 다른 ZIP으로 교체 | 설치 폴더 없음, 임시 폴더 비어 있음 |
| F-05 | `plugin install` 실패 shim | 설치 폴더만 있음, 재실행은 등록부터(이전 파일 교체 없음) |
| F-06 | 브라우저 취소 | 설치·등록 있음, `agent.json` 없음, 재실행은 새 연결 URL |
| F-07 | 승인 뒤 complete 응답 1회 유실 | 이름·ID와 '확인 필요', 새 시도 없음. 재실행은 재확인 뒤 `new` 입력 때만 새 연결 |
| F-08 | 갱신 중 설치 폴더 파일 점유 | 이전 설치 바이트 동일·등록 유지 |
| F-09 | `approval_policy = "never"` | 설치 완료, config.toml 불변, YOLO 안내 |

제품 검사(Linux 격리 clone, `4bc19f9b`): `make lint` EXIT=0, `make test` EXIT=0, 통합 `TestGoogleDeviceHTTP|TestConnectionResumeAndCodeTail|…` EXIT=0. FullOps `lint.py --from 3c8007e1`: ERROR 0, WARNING 9(규모·DEP-001·SLOP, 완료 보고에 기록). Windows: `npm run check`·adapters 단위 테스트·`go test ./internal/relay`(DB 없음) 통과, 설치기 PS 5.1·7 구문 오류 0.

한계: 최초 EN 실행의 F-03 재실행은 fixture의 회원당 agent 기록 상한(합성 데이터 누적)으로 승인 409가 나서 중단됐다. fixture가 시나리오마다 합성 기록을 지우게 고친 뒤 EN 전체를 다시 실행했다. 깨끗한 PC 기준선 시간 측정과 운영 게시는 하지 않았다.
