---
title: SAR-MVP-002-BOT-CATALOG-DEV-FIX 실행 기록
status: draft
updated: 2026-10-04
owner: dev
tasks: [SAR-MVP-002-BOT-CATALOG-DEV-FIX, SAR-MVP-002-BOT-CATALOG-DEV]
summary: 조건부 리뷰 F-01~F-08에 따라 앱 등록 안내의 원인 단정을 정정하고 installer 보존 경계를 보완한 조치·검증 및 새 재시험 댓글 초안을 기록한다
---

# SAR-MVP-002-BOT-CATALOG-DEV-FIX — 등록 안내 정정과 installer 경계 보완

## 기준

- 착수 HEAD와 lint 기준: `76443df275fb75b23d198bfde260a39096eecb28`. 기준 제품 `8e46c5a846e6d190e484e48be40b3dc368001a2b`, 리뷰 `96d9673e1a6dd80db1e5d9c5166efb64aafa0346`을 포함한다.
- 적용 규칙: `fullops-common-0.3.2`, `.fullops-squad/project.md`, FULLOPS.md, 문서 작성 규칙, ponytail full. 예외 없음.
- 원천: [조건부 리뷰](../../evaluations/qa-reports/SAR-MVP-002-BOT-CATALOG-DEV-review/report.md)와 result.json, [이슈1 원본](../../evaluations/jev/SAR-MVP-002-BOT-CATALOG-DEV-issue-1.json), Jev `SAR-MVP-002-BOT-CATALOG-DEV-FIX-*` keep 문서. Jev context의 충돌 후보는 [이전 실행 기록](SAR-MVP-002-BOT-CATALOG-DEV.md)이다. 그 원인 단정은 리뷰의 관측·미확정 구분으로 이 문서에서 정정한다. 이전 기록의 전문은 보존한다.
- 공식 근거: 2026-10-04에 `https://docs.x.ai/grok-bot/{team-bots,computer-and-apps,settings-and-notifications}.md`를 다시 받았다. 아래 인용은 그 원문이다.

## 관측 사실과 미확정 판단의 정정

| 항목 | 이전 기록·README의 표현 | 현재 판단 |
|---|---|---|
| 카탈로그 미노출 원인 (F-02) | 앱은 계정에 등록된 connector만 읽는다. CLI plugin은 앱에 등록되지 않는다 | **미확정 가설**. 관측 사실은 CLI 설치·doctor healthy(도구 2개)와 앱 MCP 14개에 knowslink 없음뿐이다. 공식 문서는 `Installed connectors are account-wide`만 쓴다 |
| 채팅 요청 → **Add MCP Server** 카드 (F-01) | 요청하면 카드에서 확인·승인한다 | **공식 근거 없음**. 제3자 글뿐이다. 이슈1 첫 보고에서 같은 Bot은 `AddMcpServer`·`InstallPlugin` 등을 호출할 수 없었고 검색도 0건이었다 |
| Command 종류의 근거 (F-03) | 일반 사실로 기술 | **Team Bots 문서에 한정**. `The Bot's info pane shows a Setup section with Plugins … Choose Add on any row, or ask the Bot in chat.`와 Plugins 표의 `Custom MCP server, Command`. 개인 계정의 같은 UI는 미확인이다 |
| 확인 위치 (F-04) | Marketplace → Your plugins(또는 Settings → Plugins) | `Settings → Plugins` 삭제. 원문: `Plugins are not a settings section.` |
| `/workspace` 유지 (F-05) | update·recovery 뒤에도 유지된다 | **보장하지 않음**. 원문: `designed to survive`, `Treat … manually installed packages … as replaceable`, Reset은 `very recent changes may be lost`. update·recover·Reset 뒤 파일을 확인하고 실패하면 installer를 다시 실행한다 |

사용자 질문(앱 Manage plugins and skills의 실제 추가 메뉴: URL MCP/Git plugin/없음)은 아직 답이 없다. 그래서 문서는 특정 메뉴의 존재를 가정하지 않는다. 등록 절차는 Setup → Plugins → Add, 채팅 요청, 중단·회신의 순서로 조건부 분기한다. 회신 항목에 실제 메뉴 이름을 넣었다.

## 조치

- `adapters/README.md` 앱 등록 절: 관측 사실·미확정 가설·공식 근거 범위를 분리했다. 등록 2단계를 조건부 분기로 바꿨다. `Settings → Plugins`를 지웠다. `/workspace` 확인 명령과 재설치 복구를 추가했다. installer가 관리하는 경로와 거절 조건을 적었다. CLI 절의 단정 문장 1개를 관측 사실로 낮췄다.
- `scripts/install_bot_mcp.sh` (F-06~F-08 최소 보완):
  - `KNOWSLINK_PREFIX`가 절대 경로가 아니면 변경 전에 exit 1.
  - `$PREFIX/node`·`$PREFIX/knowslink`가 링크이거나 이전 실행의 형태(`bin/node`, `dist/plugin.js`)가 아니면 변경 전에 exit 1. 사용자 파일을 지우지 않는다.
  - 실행마다 새 `$PREFIX/.stage.XXXXXX`를 만든다. Node archive·추출과 bundle 추출·`env -i` 경계 검사를 그 안에서 한다. 통과한 뒤에만 기존 `node`·`knowslink`와 교체한다. 실패하면 기존 준비물을 보존하고 stage와 archive는 종료 때 지운다.
  - PREFIX의 다른 파일과 예전 버전이 남긴 `$PREFIX/package`는 건드리지 않는다.
- D10·D12·D13: 원인을 가설로 낮추고 근거 범위·`AddMcpServer` 부재·조건부 등록·재설치 복구를 반영했다. D13은 이 문서의 댓글 초안을 게시 대상으로 바꿨다.
- [이전 실행 기록](SAR-MVP-002-BOT-CATALOG-DEV.md): 본문을 보존하고 제목 아래에 정정 안내 1줄만 추가했다. 이전 리뷰·원본 QA·완료 로그는 수정하지 않았다.

## 검증

실행 위치는 이 체크아웃의 scratch clone(`76443df` + 작업 트리의 installer·README)이다. Node v22.22.2, Linux x86_64. 사용자 `/workspace`와 `~/.grok`는 사용하지 않았다.

| 명령·조건 | 결과 | 로그 |
|---|---|---|
| 이전 installer: 상대 prefix, `node` 파일, `package/knowslink/STALE` | exit 0, 상대 등록 값 출력, `node` 파일이 디렉터리로 교체, STALE이 bundle에 유입(RED) | [red-old-installer](../logs/SAR-MVP-002-BOT-CATALOG-DEV-FIX/red-old-installer.txt) |
| 새 prefix 첫 실행 | exit 0, tar.gz checksum OK, PASS 2줄, 등록 값 출력 | [install-first](../logs/SAR-MVP-002-BOT-CATALOG-DEV-FIX/install-first.txt) |
| 경계 회귀 19개 | 모두 ok, `fail=0` | [boundary-regress](../logs/SAR-MVP-002-BOT-CATALOG-DEV-FIX/boundary-regress.txt), [스크립트](../logs/SAR-MVP-002-BOT-CATALOG-DEV-FIX/boundary-regress.sh.txt) |
| └ 재실행 + 이전 `package/knowslink/STALE`·기존 bundle의 STALE·무관 `notes.txt` | exit 0, Node 재사용, 새 bundle에 STALE 없음, `notes.txt`·`package/` 보존, stage 없음 | 같은 로그 |
| └ 상대 prefix `rel/pfx` | exit 1 `must be an absolute path`, 폴더 미생성 | 같은 로그 |
| └ `node`가 사용자 파일 / `knowslink`가 bundle 없는 폴더 | exit 1, 내용 보존, 다른 항목 미생성 | 같은 로그 |
| └ 잘못된 Node archive(curl shim) | exit 1 `FAILED`, 기존 node·bundle 보존, archive·stage 없음 | 같은 로그 |
| └ bundle 추출 실패(python3 shim) / `env -i` 검사 실패(env shim) | exit 1, 기존 bundle 보존, stage 없음 | 같은 로그 |
| 등록 command를 `/`에서 `env -i`로 실행 | tools `knowslink_pull_once`·`knowslink_status`, status held, stderr 0 byte | [command-probe](../logs/SAR-MVP-002-BOT-CATALOG-DEV-FIX/command-probe.txt) |
| 최종 README·installer로 재실행 | exit 0, PASS 2줄, ZIP SHA256 `d3037d2067c28bf278023a229797eb02111f8d8416bf200d23feff2bf250e609` | [install-final-rerun](../logs/SAR-MVP-002-BOT-CATALOG-DEV-FIX/install-final-rerun.txt) |
| `make lint`(product-lint)·FullOps `lint.py --from 76443df` | 완료 보고에 결과를 기록 | handovers/logs |

ZIP 해시는 ZIP에 들어가는 README가 바뀌어 `b7882df7…`에서 `d3037d20…`으로 바뀌었다. `adapters/src`·Go 코어·UI·Makefile은 바꾸지 않았다. Go test·Node 20 거절·UI·`make verify-grok-plugin`·전체 CLI QA는 반복하지 않았다. 해당 증거는 8e46c5a의 DEV·리뷰 기록을 재사용한다.

미검증: 실제 Grok Bot 계정의 Command 등록, 개인 계정의 Setup 화면·추가 메뉴·승인 카드, `AddMcpServer` 제공 여부, 앱 카탈로그 노출, Bot 세션의 `knowslink_status`, 실제 aarch64 호스트, update·Reset 뒤 `/workspace` 유지. 실제 relay·`knowslink_pull_once`·유료 inference·FullOps 업데이트는 실행하지 않았다.

## 재시험 댓글 초안 (coor 게시)

````markdown
## 앱 카탈로그 미노출 후속과 재시험 (검증 대상: `<병합 SHA>`)

관측: 지난 재시험에서 CLI 설치와 `grok mcp doctor knowslink`는 healthy(도구 2개)였지만 앱 MCP 14개에 knowslink가 없었습니다. 첫 보고에서는 같은 Bot이 `AddMcpServer`·`InstallPlugin` 등을 호출할 수 없었습니다.

원인은 **아직 미확정**입니다. 가장 유력한 가설은 앱이 계정에 등록한 connector만 카탈로그에 넣고 `~/.grok`의 CLI plugin은 읽지 않는다는 것입니다. 공식 문서는 이를 직접 말하지 않습니다.

공식 문서의 custom MCP server **Command** 등록은 Team Bots 문서에 있습니다(Bot info pane → **Setup** → **Plugins** → **Add**, 또는 채팅 요청). https://docs.x.ai/grok-bot/team-bots#set-up-what-the-bot-needs . 개인 계정에도 같은 화면이 있는지는 문서에 없어 이번 재시험으로 확인합니다.

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

기대: `...tar.gz: OK`(첫 설치), `extracted standalone MCP test exit=0`, PASS 줄 2개, 마지막에 아래 등록 값. ZIP SHA256 `d3037d2067c28bf278023a229797eb02111f8d8416bf200d23feff2bf250e609`(다르면 그대로 보고). `~/.grok`와 `/workspace/.knowslink`의 다른 파일은 바꾸지 않습니다. 실패하면 기존 준비물을 그대로 둡니다.

```text
Name:        knowslink
Type:        Command
Command:     /workspace/.knowslink/node/bin/node
Arguments:   /workspace/.knowslink/knowslink/dist/plugin.js
Environment: none (held is the default)
```

### 2. 앱 등록 (owner 직접, 화면 미확인)

1. Bot info pane에 **Setup → Plugins → Add**가 있으면 custom MCP server, Type **Command**로 위 값을 넣고 환경 변수는 비웁니다.
2. 그 화면이 없으면 같은 Bot 채팅에 보냅니다. 카드나 입력 화면이 나오면 Name·Type **Command**·command/argument·빈 환경 변수를 확인하고 승인합니다.

   ```text
   Add a custom MCP server called knowslink that runs: /workspace/.knowslink/node/bin/node /workspace/.knowslink/knowslink/dist/plugin.js
   No environment variables, no headers.
   ```

3. 두 경로 모두 Command 등록 수단이 없거나 Bot이 `AddMcpServer` 같은 도구가 없다고 답하면 **여기서 멈추고** 아래 항목을 회신해 주세요. 다른 등록 방식은 별도 승인 대상입니다.

### 3. 확인

1. **Marketplace → Your plugins**의 Installed 목록 또는 Setup → Plugins에서 knowslink를 확인합니다.
2. **새 대화**에서 `@knowslink`를 붙이고 `knowslink_status`만 호출합니다.
3. 기대: `{"state":"held","transport":"pull","actualConnection":"held","webhook":false,"evidenceFetch":false}`.

`knowslink_pull_once`·실제 relay·업무 발송은 실행하지 않습니다. 컴퓨터 update·Reset 뒤 Node나 bundle이 없으면 1단계를 다시 실행합니다(등록 값은 같음).

### 회신 항목 (실패 또는 화면이 다를 때, 비밀값 제외)

- `uname`, 1단계 마지막 30줄과 exit code.
- 앱의 plugin 추가 메뉴에 실제로 보이는 항목 이름(URL MCP, Command, Git plugin, 없음 등)과 Setup 화면 유무.
- 채팅 요청에 대한 Bot 응답 원문과 Bot이 밝힌 사용 가능 도구 목록.
- 카드·입력 화면과 Type 선택지, Installed 목록의 knowslink 상태 문구, 앱 버전.
- 새 대화의 동적 도구 검색 `knowslink` 결과와 오류 원문.

검증: 로컬에서 같은 command를 빈 환경으로 실행해 도구 2개·held를 확인했습니다. 실제 계정 등록·앱 카탈로그 노출·Bot 세션 호출은 **미검증**이며 이 재시험으로 판정합니다. 이슈는 열어 둡니다.
````

## 후속

- coor: 이 수정 고정 SHA의 delta 독립 리뷰와 필요한 좁은 후속 QA를 배정한다. 위 초안을 게시한다.
- owner: 실제 계정 재시험. 성공 조건은 Installed 목록의 knowslink와 새 대화의 `knowslink_status` held다.
- Command 등록 수단이 개인 Bot에 없으면 Remote HTTPS(공개 endpoint·인증 승인) 또는 팀 marketplace를 coor가 designer·owner에게 ask한다.
