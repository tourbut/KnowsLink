---
title: SAR-MVP-002-BOT-CATALOG-DEV 실행 기록
status: draft
updated: 2026-10-04
owner: dev
tasks: [SAR-MVP-002-BOT-CATALOG-DEV]
summary: 실제 Grok Bot 앱 카탈로그 미노출의 원인과 Command server 등록 조치·검증 및 재시험 댓글 초안을 기록한다
---

# SAR-MVP-002-BOT-CATALOG-DEV — 실제 Grok Bot 도구 카탈로그 등록 실패 진단·조치

## 기준

- 착수 HEAD와 lint 기준: `8c95bde` (`8c95bde ops: prepare actual Bot catalogue failure diagnosis`). 지시서 기준 `4d6ccfd003ee0823928f044b8528723468fba6ba`를 포함한다.
- 적용 규칙: `fullops-common-0.3.2`, `.fullops-squad/project.md`, FULLOPS.md, ponytail full, diagnosing-bugs.
- 원천: [이슈1 원본](../../evaluations/jev/SAR-MVP-002-BOT-CATALOG-DEV-issue-1.json) 댓글 5971027648·5971034506, [Jev context](../../evaluations/jev/SAR-MVP-002-BOT-CATALOG-DEV-context.json)의 keep 문서. 충돌 후보 [INSTALL-FIX 기록](SAR-MVP-002-INSTALL-FIX-DEV.md)은 CLI 성공만 주장한다. 이번 재시험이 앱 노출 실패를 보였으므로 이슈 원본을 적용했다.
- 공식 근거: `https://docs.x.ai/llms.txt`의 Grok Bot 21개 문서를 2026-10-04에 받았다. 관련 발췌는 [official-docs-excerpts](../logs/SAR-MVP-002-BOT-CATALOG-DEV/official-docs-excerpts.txt)에 있다.

## 증상

재시험(Grok CLI 1.0.40, Linux x86_64, Node v22.22.2)에서 CLI 설치와 `grok mcp doctor knowslink`는 healthy, 도구 2개였다. 같은 Bot 채팅의 동적 도구 검색 `knowslink`는 0건이었다. 앱 MCP 14개(`user-*` 12개, `cursor-github`, `cursor-origin`)에 knowslink가 없었다. `knowslink_status`는 호출할 수 없었다.

## 진단 루프

실제 앱 증상의 루프는 사람이 실제 계정에서 실행해야 한다. 이 세션에는 Grok Bot 계정과 앱 접근이 없다. 그래서 루프를 두 부분으로 나눴다.

1. agent 실행 가능 부분: 앱이 Command server로 띄울 명령을 그대로 실행한다. 빈 환경(`env -i`), 레포 밖 cwd, 고정 Node 절대경로, 준비한 bundle을 사용한다. 기대 신호는 `tools/list` 2개와 `knowslink_status` held다. [command-probes](../logs/SAR-MVP-002-BOT-CATALOG-DEV/command-probes.txt).
2. 사람 실행 부분: 아래 재시험 댓글 초안이다. 단계별 기대값과 실패 시 회신 항목을 고정했다.

CLI doctor 반복이나 mock 카탈로그로 앱 증상 해결을 주장하지 않는다.

## 가설과 판정

| 순위 | 가설 | 예측 | 판정 |
|---|---|---|---|
| 1 | 앱 카탈로그는 계정에 등록된 connector만 읽는다. CLI의 `~/.grok/installed-plugins`는 읽지 않는다 | Custom MCP server(Command)로 등록하면 노출된다. CLI 설치 여부와 무관하다 | 채택. 공식 근거 아래 |
| 2 | 앱이 CLI plugin을 새 세션이나 재시작 때만 읽는다 | 새 대화에서 knowslink가 보인다 | 근거 없음. 재시험 0단계에서 확인 |
| 3 | 앱의 실행 PATH에 Node 22가 없어 서버 시작이 실패한다 | 목록에 knowslink가 실패 상태로 보인다 | 증상과 불일치(목록에 없음). Node 절대경로로 제거 |
| 4 | 팀 connector policy가 차단한다 | `Disabled by team admin` 표시 | 증상과 불일치 |

가설 1의 근거:

- [Connect an app](https://docs.x.ai/grok-bot/computer-and-apps#connect-an-app): connector는 Marketplace에서 plugin으로 설치한다. 설치한 connector는 계정 전체에 적용된다.
- [Team Bots — Plugins](https://docs.x.ai/grok-bot/team-bots#plugins): plugin 종류에 `Custom MCP server, Remote HTTPS`와 `Custom MCP server, Command`가 있다. Command server는 대화가 사용하는 컴퓨터에서 실행된다. 환경 변수가 필요한 Command server는 owner 자신의 채팅에서만 실행된다.
- [Teams and enterprises](https://docs.x.ai/grok-bot/teams-and-enterprises#connector-policy): connector는 앱에서 plugin으로 보이고 Cursor connector policy를 따른다. Tool I/O 절은 `stdio servers on the Bot's computer`를 별도로 언급한다.
- 공식 문서 어디에도 Bot 컴퓨터의 Grok CLI plugin이 앱 카탈로그에 들어간다는 내용은 없다.
- 보조 정황(계약 아님): 로컬 grok 1.0.46 바이너리의 `grok cursor-worker`는 Cursor private worker다. 그 crate에는 exec·fs·skills·subagents·plugin artifact 설치가 있고 MCP 모듈은 없다.
- 제3자 안내(보조): [Composio](https://composio.dev/content/how-to-add-mcp-servers-to-grok-bot)는 채팅에서 Bot에게 custom MCP server 추가를 요청하는 방법과 stdio command 예를 제시한다. [Scrapeless](https://www.scrapeless.com/en/blog/grok-bot-scrapeless-connector)는 **Add MCP Server** 승인 카드를 설명한다.

원인: 이전 지침은 CLI 설치만 안내했다. 앱 계정에 MCP server를 등록하는 단계가 없었다. 그래서 앱 카탈로그에 knowslink가 없었다.

## 조치

제품 목표(실제 Bot이 KnowsLink MCP 도구를 사용)는 바꾸지 않았다. 대안 skill·terminal 호출로 대체하지 않았다. Remote HTTPS는 공개 endpoint 승인이 필요하다. Marketplace는 팀 marketplace 또는 공개 발행 승인이 필요하다. 그래서 승인 없이 가능한 **Command** server를 선택했다.

- `scripts/install_bot_mcp.sh` (신규): Linux x86_64/aarch64에 맞는 Node `v22.22.2` `.tar.gz`를 받는다. 고정 SHA256으로 확인하고 `/workspace/.knowslink/node`에 둔다. 그 Node로 bundle을 빌드·패키지·검증한다. `/workspace/.knowslink/knowslink`에 bundle을 두고 빈 환경에서 MCP 경계 검사를 실행한다. 마지막으로 앱 등록 값을 출력한다. `~/.grok`와 앱 자격 증명은 바꾸지 않는다. 환경 변수가 없으므로 기본 held다.
- `adapters/README.md`: 앱 등록 절차를 추가했다. CLI 절은 "앱 카탈로그와 별개"로 표시했다. Node 준비를 `.tar.gz`(xz 불필요)와 아키텍처 선택으로 바꿨다.
- `scripts/verify_grok_plugin.py`: 직전 리뷰 F-01(실패 시 grok 출력 표시), F-02(PASS 문구 하향), F-03(`GROK_CONFIG*` 변수 제거)를 반영했다.
- INSTALL-FIX 기록: F-04의 단정 문장을 추정으로 낮췄다. F-05는 스크립트와 README의 아키텍처 선택·gzip으로 해결했다.
- `Makefile` `lint-config`: `sh -n scripts/install_bot_mcp.sh`.
- `project.md`, D10·D12·D13: 앱 등록 경로와 검증 명령을 반영했다.

## 검증

실행 위치는 레포 루트, Node v22.22.2, npm 10.9.7, grok 1.0.46이다.

| 명령 | 결과 | 로그 |
|---|---|---|
| `KNOWSLINK_PREFIX=<scratch>/bot sh scripts/install_bot_mcp.sh` | exit 0; Node tar.gz checksum OK; package verify PASS; 준비 bundle의 `env -i` 경계 검사 PASS; 등록 값 출력 | [first](../logs/SAR-MVP-002-BOT-CATALOG-DEV/install-bot-mcp-first.txt) |
| 같은 명령 재실행 | exit 0; Node 다운로드 생략, bundle 교체 | [rerun](../logs/SAR-MVP-002-BOT-CATALOG-DEV/install-bot-mcp-rerun.txt) |
| 등록할 command를 `/`에서 `env -i`로 실행 | tools 2개, `knowslink_status` held, stderr 0 byte | [probes](../logs/SAR-MVP-002-BOT-CATALOG-DEV/command-probes.txt) |
| `uname` shim Darwin/arm64 | exit 1 `unsupported platform`, 폴더 미생성 | probes |
| `grok` 실패 shim으로 `verify_grok_plugin.py` | exit 1, grok stdout·stderr 표시 | probes |
| `make lint` | exit 0 | [lint](../logs/SAR-MVP-002-BOT-CATALOG-DEV/product-lint.txt) |
| `make verify-grok-plugin` | exit 0; ZIP SHA256 `b7882df74537ad0bd32bdde45f9dd01677431ff74dda3312ef6c8fa650c00cad` | [verify](../logs/SAR-MVP-002-BOT-CATALOG-DEV/verify-grok-plugin.txt) |

ZIP 해시는 README가 ZIP에 들어가므로 이전 `fb745c66…`에서 바뀌었다. `adapters/src`·Go 코어·UI는 바꾸지 않았다. Go test와 `make verify-mvp`는 이전 증거를 재사용한다.

미검증: 실제 Grok Bot 계정의 custom MCP 등록, 승인 카드의 Command 종류 제공 여부, 앱 카탈로그 노출, Bot 세션의 `knowslink_status` 호출, 실제 aarch64 호스트, `/workspace` 유지 동작. 실제 relay·`knowslink_pull_once`·유료 inference는 held다.

## 재시험 댓글 초안 (coor 게시)

````markdown
## 앱 카탈로그 미노출 원인과 재시험 (검증 대상: `<병합 SHA>`)

원인: Grok Bot 앱의 도구 카탈로그는 **계정에 등록된 connector**만 읽습니다. `grok plugin install`은 Bot 컴퓨터의 Grok CLI(`~/.grok`)에만 설치하므로 앱에는 보이지 않습니다. 공식 문서의 등록 방법은 Marketplace plugin 또는 custom MCP server(**Remote HTTPS**/**Command**)입니다. 근거: https://docs.x.ai/grok-bot/team-bots#plugins , https://docs.x.ai/grok-bot/computer-and-apps#connect-an-app

수정: 새 공용 서비스 없이 **Command** server로 등록하는 설치 스크립트를 추가했습니다. Node는 `.tar.gz`(xz 불필요), x86_64/aarch64를 지원합니다.

### 0. 등록 전 확인

새 Bot 대화에서 동적 도구 검색 `knowslink` 결과를 남겨 주세요. 기대: 0건(CLI 설치만으로는 앱에 보이지 않음).

### 1. 준비 (Bot 컴퓨터, 같은 bash 세션)

```bash
set -euo pipefail
uname -s; uname -m
cd /workspace/KnowsLink
git fetch origin
git checkout --detach <병합 SHA>
sh scripts/install_bot_mcp.sh
sha256sum build/knowslink-grok-bot-plugin.zip
```

기대: `...tar.gz: OK`, `extracted standalone MCP test exit=0`, PASS 줄 2개, 마지막에 아래 등록 값. ZIP SHA256 `b7882df74537ad0bd32bdde45f9dd01677431ff74dda3312ef6c8fa650c00cad`(다르면 그대로 보고). `~/.grok`는 변경하지 않습니다. CLI plugin은 제거하지 않아도 됩니다.

```text
Name:        knowslink
Type:        Command
Command:     /workspace/.knowslink/node/bin/node
Arguments:   /workspace/.knowslink/knowslink/dist/plugin.js
Environment: none (held is the default)
```

### 2. 앱 등록 (owner 직접)

같은 Bot 채팅에 보냅니다.

```text
Add a custom MCP server called knowslink that runs: /workspace/.knowslink/node/bin/node /workspace/.knowslink/knowslink/dist/plugin.js
No environment variables, no headers.
```

**Add MCP Server** 승인 카드에서 Name `knowslink`, Type **Command**, 위 command/argument, 빈 환경 변수를 확인하고 승인합니다. 비밀값은 넣지 않습니다.

### 3. 확인

1. Marketplace → Your plugins(또는 Settings → Plugins)의 Installed 목록에 knowslink가 있는지 확인합니다.
2. **새 대화**에서 `@knowslink`를 붙이고 `knowslink_status`만 호출합니다.
3. 기대: `{"state":"held","transport":"pull","actualConnection":"held","webhook":false,"evidenceFetch":false}`.

`knowslink_pull_once`·실제 relay·업무 발송은 실행하지 않습니다.

### 실패 시 회신 항목

- 0단계 결과, `uname`, 1단계 마지막 30줄과 exit code.
- 승인 카드 화면(값은 그대로, 비밀 없음)과 Type 선택지. Command 종류가 없으면 그 사실.
- Installed 목록의 knowslink 상태 문구, 앱 버전.
- 새 대화의 동적 도구 검색 `knowslink` 결과와 오류 원문.

검증: 로컬에서 같은 command를 빈 환경으로 실행해 도구 2개·held를 확인했습니다. 실제 계정 등록·앱 카탈로그 노출·Bot 세션 호출은 **미검증**이며 이 재시험으로 판정합니다. 이슈는 열어 둡니다.
````

## 후속

- 실제 계정 재시험: owner. 성공 조건은 Installed 목록과 새 대화의 `knowslink_status` held다.
- Command 종류가 개인 Bot에 없으면 Remote HTTPS(공개 endpoint·인증 승인 필요) 또는 팀 marketplace를 coor가 designer·owner에게 ask한다.
- 독립 코드 리뷰와 TESTER QA: coor가 고정 SHA로 배정한다.
