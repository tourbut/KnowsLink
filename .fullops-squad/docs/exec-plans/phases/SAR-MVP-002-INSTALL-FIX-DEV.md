---
title: SAR-MVP-002-INSTALL-FIX-DEV 실행 기록
status: draft
updated: 2026-10-04
owner: dev
tasks: [SAR-MVP-002-INSTALL-FIX-DEV]
summary: GitHub 이슈1의 Grok 설치 실패 원인과 수정·검증 및 재시험 댓글 초안을 기록한다
---

# SAR-MVP-002-INSTALL-FIX-DEV — Grok 설치 실패 수정

## 기준

- 착수 HEAD와 lint 기준: `376981441d1b0650897847273e1773f8473949e5`. 제품 기준 SHA: `0b2d5c6b6d69cba1ac40cfd8b6a9d6a2ef2eb849`.
- 패키지 수정 커밋: `e9a9c043fe765e7ce3376e5a51d2c4a3d0514c47` (브랜치 `fullops/dev`).
- 적용 규칙: `fullops-common-0.3.2`, `.fullops-squad/project.md`, FULLOPS.md, ponytail full, diagnosing-bugs.
- 원천: [이슈1 원본](../../evaluations/jev/SAR-MVP-002-INSTALL-FIX-DEV-issue-1.json), [Jev context](../../evaluations/jev/SAR-MVP-002-INSTALL-FIX-DEV-context.json). keep 문서 10개와 충돌 후보 project.md·[이전 DEV 기록](SAR-MVP-002-DEV.md)·[이전 QA](../../evaluations/qa-reports/SAR-MVP-002-DEV-TESTER.md)를 먼저 읽었다.

## 충돌 후보의 판정

이전 DEV 기록과 README는 Grok Bot이 Cursor 형식 `.cursor-plugin/plugin.json`과 `mcp.json`을 쓴다고 가정했다. 이전 QA는 ZIP 구조와 압축 해제 stdio만 검사했고 Grok 설치는 실행하지 않았다. 이슈1은 Bot 컴퓨터에 `grok` CLI가 있고 `grok plugin list`·`grok mcp list`가 동작함을 보였다. 현재 사용자 수정 요청과 실제 CLI 근거를 적용했다. 제품 규칙 충돌은 없으므로 ask하지 않았다. Cursor 파일은 Cursor IDE 검사 경로로 유지했다.

## 원인

| 근거 | 결과 |
|---|---|
| `grok plugin --help`, `grok plugin install --help` (grok 1.0.46) | `install <SOURCE>`는 git URL·GitHub shorthand·로컬 경로를 받는다. `--trust`가 MCP 활성화 조건이다 |
| CLI 동봉 안내서 `~/.grok/docs/user-guide/09-plugins.md` | manifest는 `.grok-plugin/` 또는 `.claude-plugin/`, MCP는 `.mcp.json`, 변수는 `GROK_PLUGIN_ROOT`다 |
| 바이너리 문자열 `crates/codegen/xai-grok-agent/src/plugins/manifest.rs` | `.grok-plugin/plugin.json`, `.claude-plugin/plugin.json`, `.mcp.json`만 있고 `.cursor-plugin`·`CURSOR_PLUGIN_ROOT`는 없다 |
| [Grok Build MCP 문서](https://docs.x.ai/build/features/mcp-servers) | 프로젝트 `.mcp.json`을 MCP 설정으로 읽는다 |
| 0b2d5c6 ZIP의 `grok plugin validate` | `No plugin.json found` |
| 0b2d5c6 ZIP을 임시 HOME에 설치한 `grok inspect --json` | 이름 `knowslink-e45c424f`, version null, `mcpServers: 0`, `grok mcp doctor knowslink`는 `not found` |

따라서 이전 ZIP을 `grok plugin install`해도 KnowsLink MCP 도구는 노출되지 않았다. 이슈의 `grok mcp list` 빈 결과는 별도 이유도 있다. 이 명령은 config.toml 서버만 보여 주고 plugin 서버는 보여 주지 않는다. Bot 앱 도구 목록에 InstallPlugin 등이 없던 것은 앱의 동적 카탈로그 문제이며 CLI 설치와 다른 경로다. 이번 수정은 CLI 설치 경로를 고쳤고 앱 카탈로그 노출은 주장하지 않는다.

## 수정

- `adapters/.grok-plugin/plugin.json`: 이름 `knowslink`, version `0.1.0`. 구성요소는 기본 디렉터리에서 찾는다.
- `adapters/.mcp.json`: `node ${GROK_PLUGIN_ROOT}/dist/plugin.js`, `KNOWSLINK_MODE=held`.
- `scripts/package_plugin.py`: 두 파일과 `.grok-plugin/marketplace.json`을 ZIP에 추가한다. Cursor 파일은 유지한다.
- `adapters/.npmrc`: `engine-strict=true`. engines `>=22.22.2 <23`은 그대로다.
- `scripts/verify_grok_plugin.py`, `make verify-grok-plugin`: 실제 grok CLI로 임시 HOME에 설치하고 doctor·held를 검사한다.
- `adapters/README.md`: Grok 설치 명령·기대 결과·Node 22 준비 절차. `project.md`: 검증 명령 추가.

## Node 20 경고

공식 Node 20.19.2 배포본(SHASUMS256 확인)으로 `npm ci`를 실행해 이슈의 `EBADENGINE` 경고를 재현했다(exit 0, 경고 후 계속). 같은 Node 20으로 설치본 bundle의 MCP 검사도 통과했다. 그러나 [Node 공식 일정](https://github.com/nodejs/Release/blob/main/schedule.json)상 Node 20 지원은 2026-04-30에 끝났다. 그래서 engines를 넓히지 않았다. 대신 `engine-strict`로 Node 20의 `npm ci`가 `EBADENGINE` 오류(exit 1)로 멈추게 했다. Node 22.22.2(npm 10.9.7, 2026-03-24 배포)의 공식 tarball 준비 절차를 README에 넣었다. MCP 서버는 Grok 프로세스 PATH의 `node`로 실행된다. `grok mcp doctor`의 `command found` 경로로 확인한다.

## 검증

실행 위치는 레포 루트, Node v22.22.2, npm 10.9.7, grok 1.0.46이다. grok 검사는 모두 임시 HOME에서 실행했다. 사용자 `~/.grok/installed-plugins`는 바꾸지 않았다.

| 명령 | 대상 | 결과 | 로그 |
|---|---|---|---|
| `scripts/verify_grok_plugin.py` | 0b2d5c6 소스 ZIP | exit 1, `knowslink-2f0467f4`, version None (RED) | [red](../logs/SAR-MVP-002-INSTALL-FIX-DEV/red-0b2d5c6-package.txt) |
| `grok inspect --json` 전후 | 0b2d5c6 / 수정 ZIP | mcpServers 0 → 1 | [inspect](../logs/SAR-MVP-002-INSTALL-FIX-DEV/grok-inspect-before-after.txt) |
| `make verify-grok-plugin` | e9a9c04 | exit 0; validate·install·list·doctor healthy·2 tools·설치본 held | [verify](../logs/SAR-MVP-002-INSTALL-FIX-DEV/verify-grok-plugin.txt) |
| `grok plugin marketplace add <폴더>` → `grok plugin install knowslink --trust` | e9a9c04 ZIP | exit 0, doctor healthy 1 | 이 기록 |
| `make lint` | e9a9c04 | exit 0 | [lint](../logs/SAR-MVP-002-INSTALL-FIX-DEV/product-lint.txt) |
| `make test` | e9a9c04 | exit 0 | [test](../logs/SAR-MVP-002-INSTALL-FIX-DEV/product-test.txt) |
| Node 20.19.2 `npm ci` + engine-strict | e9a9c04 | exit 1, `EBADENGINE` | [node20](../logs/SAR-MVP-002-INSTALL-FIX-DEV/node20-npm-ci-engine-strict.txt) |
| Node 22.22.2 `npm ci --prefix adapters` | e9a9c04 | exit 0 | 이 기록 |

`make plugin`의 ZIP SHA256은 `fb745c66b4e786f5267228099c3794763632381451bd37cf66a82111f9b14a75`이며 두 번 생성해도 같았다. ZIP 항목은 11개이며 credential·fixture·node_modules·.env는 없다. 도구 출력에 비밀값이 없다는 검사는 MCP 시험의 PASS 줄로 확인했다.

미검증 항목은 다음과 같다. 실제 Bot 계정에서의 설치, Bot 세션의 `knowslink_status` 호출, 앱 Marketplace UI 노출, hosted Node 버전은 원격 재시험 전까지 미확정이다. Grok 세션에서 모델을 통한 도구 호출은 유료 inference이므로 실행하지 않았다. 실제 relay·DEC-02·calendar는 held다.

## GitHub 이슈 댓글 초안

coor가 리뷰·QA·push 뒤 게시한다. `<SHA>`는 병합된 main SHA 또는 `e9a9c043fe765e7ce3376e5a51d2c4a3d0514c47`이다.

````markdown
## 원인과 수정

`0b2d5c6`의 ZIP은 Cursor 형식(`.cursor-plugin/plugin.json`, `mcp.json`)만 담았다. Grok CLI 1.0.46은 `.grok-plugin/plugin.json`과 `.mcp.json`만 읽는다. 그래서 설치해도 이름은 `knowslink-<hash>`, MCP 서버는 0개였다. `grok mcp list`는 plugin 서버를 보여 주지 않는다. 확인에는 `grok mcp doctor`를 쓴다.

`<SHA>`에서 Grok manifest와 `.mcp.json`을 추가했다. Node 20은 2026-04-30 EOL이므로 빌드에 Node 22.22.2를 요구한다. Node 20의 `npm ci`는 이제 `EBADENGINE` 오류로 멈춘다.

## 재시험 (KNOWSLINK_MODE=held 유지, relay 미사용)

0. Node 22.22.2 준비 (이미 v22.22.2면 생략)

```sh
cd /workspace
curl -fsSLO https://nodejs.org/dist/v22.22.2/node-v22.22.2-linux-x64.tar.xz
curl -fsSL https://nodejs.org/dist/v22.22.2/SHASUMS256.txt | grep ' node-v22.22.2-linux-x64.tar.xz$' | sha256sum -c -
mkdir -p "$HOME/.local/node" && tar -xJf node-v22.22.2-linux-x64.tar.xz -C "$HOME/.local/node"
export PATH="$HOME/.local/node/node-v22.22.2-linux-x64/bin:$PATH"
node --version && npm --version
```

기대: `node-v22.22.2-linux-x64.tar.xz: OK`, `v22.22.2`, `10.9.7`.

1. 다운로드·빌드

```sh
cd /workspace/KnowsLink
git fetch origin && git checkout --detach <SHA>
make install && make plugin
sha256sum build/knowslink-grok-bot-plugin.zip
```

기대: 둘 다 exit 0, `extracted standalone MCP test exit=0`. 같은 도구 버전이면 SHA256 `fb745c66b4e786f5267228099c3794763632381451bd37cf66a82111f9b14a75`. 다르면 값을 그대로 적는다.

2. 설치

```sh
rm -rf /tmp/knowslink-plugin
python3 -m zipfile -e build/knowslink-grok-bot-plugin.zip /tmp/knowslink-plugin
grok plugin validate /tmp/knowslink-plugin/knowslink
grok plugin install /tmp/knowslink-plugin/knowslink --trust
```

기대: `Plugin manifest is valid.`, `name: knowslink`, `version: 0.1.0`, `MCP servers`. 이어서 `Installed 1 plugin(s) from ...: knowslink`.

3. 카탈로그

```sh
grok plugin list --json
grok mcp doctor knowslink --json
```

기대: list에 `"name": "knowslink"`, `"version": "0.1.0"`. doctor에 `"healthy": true`, `command found`(경로의 `node --version`이 v22.22.2), `handshake OK`, `2 tools discovered`.

4. held 상태

새 Grok Bot 세션을 시작하고 `knowslink_status`를 호출한다. `knowslink_pull_once`는 호출하지 않는다.

기대: `{"state":"held","transport":"pull","actualConnection":"held","webhook":false,"evidenceFetch":false}`. 도구가 보이지 않으면 Bot 앱 버전, 세션 도구 목록, 3단계 출력을 그대로 적는다.

제거: `grok plugin uninstall knowslink`.
````

## 남은 일

coor가 고정 SHA 독립 리뷰와 TESTER QA를 배정한다. 리뷰·QA·push 뒤 위 댓글을 게시하고 원격 재시험 결과를 받는다. Bot 앱 동적 카탈로그 노출과 hosted Node는 재시험 결과로 판정한다.
