---
title: SAR-AUTOSETUP-001-DEV — Windows 11 Claude Code·Codex 자동 설정 기술 계획
status: draft
updated: 2026-10-23
owner: dev
tasks: [SAR-AUTOSETUP-001-DEV]
summary: 한 줄 PowerShell 설치기·고정 배포물·호스트 등록·Google 연결 재개·갱신 교체의 기술 선택 근거와 한계
---

# SAR-AUTOSETUP-001-DEV — Windows 11 Claude Code·Codex 자동 설정 기술 계획

승인된 [intent SAR-AUTOSETUP-001](../../planning/intents/SAR-AUTOSETUP-001.md)(D031–D037)과 [D02](../../planning/product-specs/SAR-AUTOSETUP-001.md) AS-01–13을 구현한다. 기준 커밋은 main `3c8007e1`이다(지시서 커밋 `29305b76`). 실측 결과는 [DEV QA 보고](../../evaluations/qa-reports/SAR-AUTOSETUP-001-DEV-test/report.md)에 있다.

## 확인한 근거

- 이 PC의 호스트는 Claude Code `2.1.308`, Codex CLI `0.163.0`이다. Context7은 월간 한도 초과로 조회하지 못했다. 대신 설치본의 `claude plugin --help`·`claude plugin marketplace add|install|list --help`, `codex plugin --help`·`codex plugin marketplace --help`를 읽고 격리 `CLAUDE_CONFIG_DIR`·`CODEX_HOME`에서 실제 동작을 확인했다.
- Claude Code: `plugin marketplace add <dir> --json` 뒤 `plugin install knowslink@knowslink-plugins --scope user --json`. 같은 명령의 재실행은 exit 0이다. 다른 경로로 다시 add하면 같은 이름의 marketplace가 새 경로를 가리킨다. 디렉터리 marketplace의 plugin은 설치 폴더에서 제자리 로드된다("loads in place"). 설정은 `settings.json`의 `extraKnownMarketplaces.knowslink-plugins`·`enabledPlugins`에만 쓰인다.
- Claude Code `2.1.308`의 `plugin marketplace list`는 사용자 `model` 별칭을 바꾼다(`opus`→`opus[1m]`). 새 설정 폴더에서는 첫 plugin 명령이 같은 1회성 이주를 한다. 설치기는 `marketplace list`를 쓰지 않고 `plugin list --json`의 `readFromFolder`로 등록 경로를 찾는다.
- Codex: `plugin marketplace add <dir>`·`plugin add knowslink@knowslink-plugins --json`. 같은 source의 재실행은 exit 0이고 cache를 새로 복사한다. 다른 source로 add하면 실패하므로 먼저 `plugin marketplace remove`가 필요하다. 설정은 `config.toml`의 `[marketplaces.knowslink-plugins]`·`[plugins."knowslink@knowslink-plugins"]`에 쓰인다.
- npm으로 설치한 `codex`는 `codex.ps1` shim이 먼저 잡힌다. 실행 정책이 Restricted면 막히므로 설치기는 `Application`(`.cmd`/`.exe`)만 실행한다.
- Node `22.22.2` 공식 zip SHA-256: win-x64 `7c93e9d9…c57c`, win-arm64 `380d375c…6ab3`(nodejs.org `SHASUMS256.txt`).
- relay의 device 연결 `complete`는 1회만 credential을 준다. 응답을 잃으면 다시 받을 수 없다(F-07).

## 수동 절차 기준선

깨끗한 Windows 11 PC가 없어 실측 시간 대신 문서 기준 수를 기록한다. 시간 실측은 tester 큰 틀 검수(AC-04)에서 한다.

| 항목 | 현재 수동 절차(Claude Code / Codex) | 자동 설정 |
|---|---|---|
| 붙여 넣는 셸 명령 | 5개 이상 / 7개 이상(clone·checkout·npm ci·build·package·marketplace·install, Codex는 `CODEX_HOME`) | 1개 |
| 직접 입력하는 절대 경로 | 2개 이상 | 0개 |
| 사람 행동 | clone부터 세션 연결까지 9단계 | 명령 실행, 설치 계획 동의, 브라우저 Google 동의, 호스트 재시작, 세션 연결 |

## 기술 선택

| 항목 | 선택 | 이유 |
|---|---|---|
| 설치기 형태 | Windows PowerShell 5.1 호환 스크립트 `scripts/autosetup/knowslink-setup.ps1`, 명령 `& ([scriptblock]::Create((irm <base>/knowslink-setup.ps1))) -App claude|codex -Lang ko|en` | Windows 11 기본 셸이다. script block은 실행 정책 파일 검사를 받지 않는다. 인자로 호스트·언어를 넘긴다. `exit`는 사용자 창을 닫으므로 쓰지 않고 내부 예외로 멈춘다. |
| 배포 위치 | `https://link.knowslog.com/install/<version>/`(relay `KNOWSLINK_INSTALL_DIR`, 허용 파일 3개만) | KnowsLink 도메인 HTTPS(Q1). 설치기는 `text/plain; charset=utf-8`로 보내 `irm`이 한국어를 바르게 읽는다. 공개 ingress 규칙에 같은 경로만 더했다. |
| 고정 버전·checksum | `package_plugin.py --release <dir>`가 설치기에 버전·URL·두 ZIP SHA-256을 채운다. Node zip SHA-256은 설치기에 고정한다 | 내용물이 다르면 F-04로 멈추고 임시 폴더를 지운다. ZIP은 기존 결정적 패키지(고정 시각)다. |
| Node 준비 | `%LOCALAPPDATA%\KnowsLink\runtime\node-v22.22.2-win-<arch>\node.exe`를 공식 zip에서 `node.exe`만 꺼내 둔다 | 사용자 Node·PATH와 무관하다. 설치 폴더 밖이라 갱신 교체가 실행 중 파일에 막히지 않는다. |
| 설치 폴더 | `%LOCALAPPDATA%\KnowsLink\<host>-NN`과 `<host>-NN-private`. 둘 다 없는 첫 번호 | 기존 설치·agent 폴더를 덮지 않는다(AS-03, AC-06). 기존 `install()` 보호 조건과 같은 구조다. |
| 호스트 설정 값 | MCP command는 runtime Node, args는 설치 폴더의 `dist/plugin.js`. env는 `install()`과 같다(Claude `KNOWSLINK_AUTO_RECEIVE=off`, Codex `KNOWSLINK_SESSION_HOST=codex`) | 기존 세션 연결 계약(SP-01)을 그대로 쓴다. |
| 기존 설치 찾기 | 호스트 등록 경로 + 기본 폴더의 설치 기록. 하나면 그것, 여럿이면 목록과 번호 선택, 불명확하면 `[확인 필요]` | AS-10. 기존 설정 파일의 `KNOWSLINK_AGENT_FOLDER`를 그대로 쓴다. agent 폴더는 `agent.json` 존재만 본다. |
| 원자적 교체 | 새 내용물을 `<설치>.new-<시각>`에 풀고 설정·검사한 뒤 `<설치>`→`<설치>.previous-<시각>`, 새 폴더→`<설치>`로 이름을 바꾼다. 두 번째가 실패하면 되돌린다 | U3·F-08. 이전 설치를 보존한다(AS-03). |
| Google 연결 | `plugin.js setup-connect <relay> <agent 폴더> [--fresh]`가 JSON 줄로 공개 상태만 알린다. 설치기가 링크를 열고 지문을 보인다 | 기존 `beginLogin`·`finishLogin`·`complete`를 재사용한다. 대기 중 시도는 이어서 기다리고, 취소·만료 시도는 옆으로 옮긴 뒤 새로 시작한다(AS-11). 승인 뒤 응답 유실은 새 시도 전에 저장 결과를 다시 확인하고 사람이 `new`를 입력해야 새로 시작한다(F-07). |
| agent 이름 | relay device poll의 approved 응답에 이름을 더하고 `agent.json`에 표시용으로 보존한다 | relay `/v1/contacts`는 자기 조회를 403으로 막는 기존 계약이 있다. 이 계약은 바꾸지 않았다. |
| 설치 기록 | `<설치>\knowslink-install.json`: 버전·시각·단계 결과·agent ID/이름/공개 지문·호스트 항목·해제 명령 | AS-12·AS-13. 키·토큰은 없다. 실패 전 교체 단계(F-03/04/08)에서는 기존 폴더에 쓰지 않는다. |
| YOLO 확인 | `config.toml` 최상위 `approval_policy = "never"`면 완료 뒤 F-09를 보인다 | AS-08. 설치기는 승인 모드를 바꾸지 않는다. |
| 회원 화면 | `use=claude` 안내를 D04 카피로 바꾸고 `use=codex`를 더한다. 명령은 `textarea.share`, 언어는 `{{lang}}`. 문장 안 명령은 새 `code` 템플릿 함수로 넣는다 | 번역은 문장 조각 단위다. `<code>` 태그가 문장을 나누면 영어 어순을 맞출 수 없다. |

## 범위 밖과 한계

- 실제 호스트 세션의 `/knowslink-session on`·발송 승인 대화상자(AC-02·AC-03)는 사람 입력이 필요하다. DEV는 등록된 MCP 설정을 확인했고 실제 세션 발송은 tester 큰 틀 검수에서 확인한다. 세션·발송 코드는 바꾸지 않았다.
- 운영 relay 배포·`/install/0.1.0/` 게시는 ops 후속이다(사용자 승인). 게시 전에는 회원 화면 명령이 404다.
- 이 PC의 Claude Code는 `2.1.308`이라 모든 실행이 '검증하지 않은 환경' 동의를 거쳤다. 완료 화면의 '시험 환경 확인' 문구는 `2.1.302`에서만 나온다.
- macOS·Linux·ARM Windows 실측은 하지 않았다. ARM은 Node arm64 zip을 쓰도록만 했다.
