---
title: SAR-AUTOSETUP-001-QA — Claude Code·Codex 자동 설정 큰 틀 검수 보고
status: draft
updated: 2026-10-23
owner: tester
tasks: [SAR-AUTOSETUP-001-QA]
summary: 실제 화면 명령을 격리 relay에서 실행해 N1–N6·갱신·화면(S7)·자격 경계를 확인했다. 실제 호스트 로그인이 없어 세션 연결·첫 발송(AC-02/03)과 AC-04 전체 시간은 미확인이다. 미해결 high 0, 새 관찰 3건.
---

# SAR-AUTOSETUP-001-QA — 큰 틀 검수 보고

독립 동작 QA(lite, 새 사용자 흐름 1회)다. 제품 코드는 고치지 않았다. 운영 relay·실제 Google 계정·실제 상대·사용자 실제 `~/.claude`·`~/.codex`는 쓰지 않았다.

**결론: 부분 완료.** 통과 조건(AC-01–AC-04 PASS)은 충족하지 못했다. 이유는 제품 결함이 아니라 검증 환경이다. 격리 `CLAUDE_CONFIG_DIR`·`CODEX_HOME`의 실제 호스트 로그인은 브라우저 OAuth 동의가 필요한데, 자동화 권한 분류기가 거부했다(우회하지 않음). 사용자가 로그인해 주면 남은 항목을 이어서 검증한다(`남은 일`).

## 고정 후보와 환경

| 항목 | 값 |
|---|---|
| 제품 SHA | `11d91648b30f13aefd545b19d0508fff9eb5b667` (fullops/tester = main `4fefdf68` + 문서 1커밋). 격리 `git worktree add --detach`에서 빌드, 제품 파일 변경 0 |
| 호스트 | Windows 11 Pro 26300, Windows PowerShell 5.1, Claude Code **2.1.309**(검증 버전 2.1.302와 다름), Codex CLI 0.163.0(검증 버전) |
| 배포물 | `package_plugin.py --release … --relay-url http://127.0.0.1:18441`, Claude ZIP `9361d4f1…`, Codex ZIP `dbea5ece…` |
| relay | 격리 Postgres 17(Compose 프로젝트 `knowslink-autosetup-qa`, 15481) + 합성 회원 A(Claude)·B(Codex) 친구 1쌍을 만드는 시험 fixture(`tools/autosetup_qa_fixture_test.go.txt`, 워크트리 밖 격리 복사본에서만 실행, 미커밋) |
| fixture가 대신한 것 | Google OAuth 콜백 한 가지(연결 시도에 회원 binding). 연결 화면(`/connect/<id>`)의 지문·이름 입력·`/home/device-confirm`은 실제 relay 코드로 실행했다. 실제 Google 로그인·동의는 하지 않았다 |
| 명령 | 회원 화면 `use=claude`·`use=codex`가 보인 한 줄을 그대로 실행. 바꾼 것은 운영 URL `https://link.knowslog.com`→격리 `http://127.0.0.1:18441` 뿐 |
| 도구 | [tools/install_run.py](tools/install_run.py), [qa.py](tools/qa.py), [snap.py](tools/snap.py), [realhash.py](tools/realhash.py), [mcp_roundtrip.mjs](tools/mcp_roundtrip.mjs) |

## S1–S7 판정

| ID | 판정 | 근거 |
|---|---|---|
| S1 Claude 행동 ≤6 | **부분 확인(N1–N6)**, N7–N8 미확인 | 아래 행동 표. 사람 입력 경로·환경 변수·설정 편집 0, 붙여 넣는 명령 1. 현재 호스트가 2.1.309라 '검증하지 않은 환경에서 계속' `y`가 한 번 더 있다 → 관찰 O-01 |
| S2 Codex 행동 | **부분 확인(N1–N6)**, 세션 연결·코드 입력 미확인 | `CODEX_HOME`·marketplace·plugin add 수동 실행 0. 검증 버전이라 추가 `y` 없음 |
| S3 ≤10분 | **미확인** | 설치기 구간만 실측: Claude 19.7초, Codex 21.3초(Node 22.22.2 실제 내려받기 포함, 브라우저 단계는 fixture라 1.5초). 세션 연결까지는 실호스트 로그인 필요 |
| S4 자격 경계 | **PASS(설치기 구간)** | 격리 호스트 설정에는 KnowsLink 항목 2개만 생김(ev/ 아래 증거). 출력·설치 기록에 credential·개인키 0건(검색). 호스트 인증 파일을 읽은 흔적 0(상자에 인증 파일 자체가 없음). 사용자 실제 파일: 아래 해시 표 |
| S5 갱신 | **PASS** | 같은 명령 재실행이 Google 없이 `기존 설치 갱신`. 호스트별 전후 agent ID·`agent.json`·`private.pem`을 포함한 private 폴더 해시·웹 agent 수(1)·친구 수(1) 동일 ([claude 전](ev/snap-claude-before-update.txt)/[후](ev/snap-claude-after-update.txt), [codex 전](ev/snap-codex-before-update.txt)/[후](ev/snap-codex-after-update.txt)). 갱신 시간 3.6초·2.1초. 이전 설치 `.previous-*` 보존 |
| S6 실패 안내 | **미확인(이번 실측 범위 밖)** | F-01–F-09는 DEV가 KO·EN PASS(주입 fixture). 이번에는 EN에서 계획 거절(F-02 후속)만 확인: `Not started. Nothing was changed.`, 설치 폴더 생성 0 |
| S7 화면 | **PASS** | 회원 화면 선택지 6개(web·grok·dots·claude·codex·other) 중 명령은 claude·codex에만 있다. KO·EN 모두 '시험 환경 확인 (Claude Code 2.1.302, Windows 11)', Codex는 '도착 알림 없음'. 실제 설치 호스트가 2.1.309여도 화면은 검증 버전을 말한다 ([claude 화면](ev/screen-claude.html.txt), [codex 화면](ev/screen-codex.html.txt)) |

## AC-01 사람 행동 표(실측 N1–N6)

| 단계 | 사람 행동 | 입력 종류 | Claude | Codex | 기준선(intent 9단계) |
|---|---|---|---|---|---|
| N1 | 화면 명령 복사→PowerShell 붙여 넣기 | 붙여 넣기 1회 | ○ | ○ | 5개 이상 명령 + 절대 경로 2번 이상 |
| N2/AS-09 | 검증하지 않은 환경에서 계속 `y` | 한 글자 | ○ (2.1.309) | 없음 | 해당 없음 |
| N3 | 설치 계획 동의 `y` | 한 글자 | ○ | ○ | 해당 없음 |
| N4 | (자동) 내려받기·checksum·Node·설치·등록·등록 확인 | 없음 | 완료 줄 4개 | 완료 줄 4개 | 단계 1–5 |
| N5 | 브라우저: 로그인·지문 비교·이름·동의 | 브라우저 | fixture 대행 | fixture 대행 | 단계 7 |
| N6 | 완료 화면 읽기(재시작 안내) | 없음 | ○ | ○ | 단계 8 |
| 재시작 | Claude Code/Codex를 닫고 다시 열기(실행 중이었을 때만) | 없음 | 조건부 | 조건부 | 해당 없음 |
| N7 | `/knowslink-session on` / `knowslink_session on`+Codex 도구 승인+4자리 코드 | 명령·코드 | **미확인** | **미확인** | 단계 9 |
| N8 | 메시지 요청·발송 승인 | 승인 | **미확인** | **미확인** | 단계 10 |

경로·환경 변수·설정 파일 편집 입력은 모두 0이었다. 셸에 붙여 넣은 명령은 호스트당 1개다. 터미널에 보인 지문과 연결 화면의 지문은 두 호스트 모두 같았다(`SHA256:OXwJi6Vo…`, `SHA256:Mo9CCGun…`). N1–N6의 호스트별 사람 입력은 Claude(2.1.309) 3개(붙여 넣기·`y`·`y`)+브라우저 1, Codex 2개(붙여 넣기·`y`)+브라우저 1이다. N7·N8을 더하면 Claude 총 6(재시작 없음)~7(재시작 있음), 검증 버전의 Claude Code 2.1.302면 5~6이다.

## AC-02/03 첫 발송 — 실호스트 미확인, MCP 수준 대체 확인

실호스트 세션이 없어 AC-02/03은 **판정하지 않는다.** 대신 설치기가 호스트 설정에 등록한 MCP 서버 설정(`config/local.mcp.json`, `.mcp.json`, 절대 경로 node·plugin.js·agent 폴더 포함)을 그대로 읽어 스크립트 클라이언트로 띄웠다 ([mcp-roundtrip.jsonl](ev/mcp-roundtrip.jsonl)).

| 방향 | 결과 |
|---|---|
| Codex B `knowslink_session on` | 확인 코드 창(세션 구분값·내 agent ID 표시) → 코드 입력 → `KnowsLink 연결됨: Codex 세션 … · 대기 0건` |
| Claude A → Codex B | `transport: queued` → B `receive`로 본문·발신자 이름 일치(`untrusted: true`, `sessionReceipt: recorded`) |
| Codex B → Claude A | 발송 창에 본문 전체 22자 표시 → 코드 입력 → `queued` → A `receive` 일치 |

이것은 설치기 산출물이 실제로 동작하는 MCP 설정을 만든다는 증거다. 호스트의 도구 discovery, `/knowslink-session on` mod, Claude 발송 승인 대화상자, Codex 도구 실행 승인은 검증하지 않았다.

## 사용자 실제 설정 파일 전후 (S4)

[전](ev/real-before.txt)·[후](ev/real-after.txt) 비교. 후 해시는 `realhash.py`를 같은 조건으로 다시 실행한 결과다. 값이 달라진 파일은 `~/.claude/plugins/known_marketplaces.json`(및 같은 파일을 가리키는 이 세션 프로필)이다. 크기 1751바이트 동일, `knowslink` 문자열 0건, 모든 항목의 `lastUpdated`가 `2026-10-23T01:20:46Z`로 갱신됐다. 호스트의 마켓플레이스 자동 새로고침으로 보이며 KnowsLink 항목은 생기지 않았다. 설치기는 항상 격리 `CLAUDE_CONFIG_DIR`·`CODEX_HOME`로 실행했다(설치기가 만든 항목은 격리 상자 안에서만 확인). `~/.codex/config.toml`·`auth.json`·`~/.claude/settings.json`·`installed_plugins.json`은 전후 동일하다. 이 세션 프로필의 `.claude.json`은 실행 중인 이 세션이 계속 쓰므로 비교 대상에서 뺀다.

## 관찰·결함

심각도 high 0, 미해결 critical 0.

| ID | 심각도 | 내용 | 재현·영향 |
|---|---|---|---|
| O-01 | low–medium(S1 여유 소진) | 현재 설치된 Claude Code는 2.1.309인데 검증 버전은 2.1.302다. 새 사용자는 호스트가 자동 갱신되므로 대부분 같은 `y`를 한 번 더 입력한다. S1 행동 수가 한도 6에 닿고, 호스트 재시작이 필요하면 7이 된다 | 재현: 2.1.309에서 화면 명령 실행 → `검증하지 않은 환경에서 계속할까요? (y/N)`. 조치는 designer 판단(S1 목록에 포함 또는 검증 버전 갱신) |
| O-02 | low | 같은 고정 포트(18431)에 다른 팀의 relay가 떠 있으면 격리 fixture가 조용히 실패하고 다른 서비스의 404가 응답한다. 제품 문제가 아니라 시험 환경 충돌. 포트를 18441로 옮겼다 | 시험 환경 기록 |
| O-03 | info | 완료 화면의 `시험 환경 확인` 문구가 Codex(검증 버전)에만 나오고 Claude 2.1.309에는 없다. AS-09 의도대로다 | 화면 대조 |

## 한계 (정직한 범위)

- 실제 Google 로그인·동의는 하지 않았다(fixture 대행). 사용자 실계정 비밀번호 입력 0.
- 실제 Claude Code·Codex 로그인과 세션은 쓰지 못했다(위 결론). 세션 연결·첫 발송·소요 시간 AC-02/03/04는 미확인이다.
- Claude 첫 설치 원시 로그는 도구가 같은 이름으로 덮어써 잃었다(결과는 이 보고서에 옮겼다). 이후 실행은 시각 접미사를 붙인다. Codex 첫 설치·갱신 로그는 [codex-ko1-install.txt](ev/codex-ko1-install.txt), [codex-ko1-update.txt](ev/codex-ko1-update.txt).
- F-01–F-09 전체와 EN 전체는 이번에 다시 실행하지 않았다(DEV 결과 신뢰, 근거 `SAR-AUTOSETUP-001-DEV-test/report.md`).

## 남은 일

1. 사용자가 격리 상자에서 두 호스트에 로그인한다. `$env:CLAUDE_CONFIG_DIR='C:\Users\shin\AppData\Local\KnowsLinkDevTools\autosetup-001-qa\work\claude-ko1\claude'; claude auth login` 과 `$env:CODEX_HOME='C:\Users\shin\AppData\Local\KnowsLinkDevTools\autosetup-001-qa\work\codex-ko1\codex'; codex login`.
2. 로그인 뒤 tester가 실세션으로 N7–N8, AC-02/03, AC-04 시간 측정, S1·S2·S3 최종 판정(약 20–30분).
3. designer: O-01(S1 행동 수와 검증 버전 문구) 판단.
