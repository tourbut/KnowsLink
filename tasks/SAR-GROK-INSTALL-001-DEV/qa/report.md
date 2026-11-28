# SAR-GROK-INSTALL-001-DEV QA 보고

- 일자: 2026-11-28 / 구현 커밋: `98ce0565` / 테스트 레벨: lite
- 방향: lead 정정(5679d1ee)에 따라 curl 한 줄 installer 대신 Grok Bot 플러그인으로 배포한다. 코디네이터가 답하지 않아 선택지 1(팀 marketplace Git 저장소)을 dev가 정했다.

## 공식 형식 근거

| 사실 | 근거 |
|---|---|
| 개인의 Add a custom connector는 Remote HTTPS·Command만 지원한다 | https://docs.x.ai/grok-bot/computer-and-apps#connect-an-app |
| Team Bot setup이 plugins를 묻는다. Team Marketplace는 dashboard Plugins & MCPs에 있다 | https://docs.x.ai/grok-bot/create-a-team-bot, https://docs.x.ai/grok-bot/teams-and-enterprises |
| 팀 marketplace는 Teams/Enterprise 플랜에서 쓴다. Import from Repo가 GitHub URL을 받는다. 갱신은 Refresh/Auto Refresh다. 공개 marketplace는 수동 심사를 거친다 | https://cursor.com/docs/plugins |
| `.cursor-plugin/plugin.json`·`marketplace.json`·`mcp.json`을 쓰고 `${CURSOR_PLUGIN_ROOT}`를 확장한다 | https://cursor.com/docs/reference/plugins |

Grok Bot이 Cursor 플러그인 형식을 그대로 읽는지는 실제 계정에서 확인하기 전까지 추정이다. `https://docs.x.ai/grok-bot/team-bots`에는 `#set-up-what-the-bot-needs` 절이 없다. `grok.com/help/grok-bot/connect-plugins`는 본문 없는 페이지를 돌려줬다.

## 산출물 검증

- `python3 scripts/package_plugin.py --grok-release <새 폴더>`(Node 22.22.2, Windows)는 exit 0이었다. 추출한 standalone MCP 검사와 preparation 경계 검사가 PASS였다. 트리 구성: 루트 `.cursor-plugin/marketplace.json`·`.grok-plugin/marketplace.json`·`SHA256SUMS`, `knowslink/{.mcp.json, mcp.json, bin/knowslink-grok.sh, dist/plugin.js, …}`. 두 MCP 설정은 `sh ${GROK_PLUGIN_ROOT|CURSOR_PLUGIN_ROOT}/bin/knowslink-grok.sh`이며 env가 없다.
- `grok plugin validate <tree>/knowslink`(grok 1.0.92, 격리 HOME): `Plugin manifest is valid. … components: 1 skill dir(s) … MCP servers`.

## Linux 컨테이너 증명 (x86_64, `node:22` 이미지, PATH에서 이미지의 node를 숨김)

스크립트는 [qa.sh](qa.sh)와 [qa_offline.sh](qa_offline.sh)다. 트리는 `/plugin`에 읽기 전용으로 붙였고 `/workspace`는 새 볼륨이다. MCP는 앱처럼 `env -i`로 런처를 거쳐 실행했다.

```text
arch: x86_64
PASS 1 fresh install: env -i node --version=v22.22.2, plugin.js present, initialize ok, knowslink_status public_node_configured_unverified, env -i MCP boundary test exit 0, no stage left
PASS 2 rerun reuses Node v22.22.2 without download; old v22.0.0 made by the launcher is replaced by verified v22.22.2
PASS 3a checksum mismatch: exit 1, 'Node checksum mismatch; the existing /workspace/.knowslink/node is kept', existing node kept, no stage left
PASS 3c foreign node/: refused ('refusing to replace /workspace/.knowslink/node: not made by this installer (node/bin/node missing); move it away first'), kept
PASS 3b download failure (network none): exit 6, 'curl: (6) Could not resolve host: nodejs.org', existing node kept, no stage left
```

- 갱신: 플러그인 파일은 팀 marketplace Refresh가 바꾼다. 런처 쪽 갱신은 같은 prefix를 다시 쓰는 것과 Node 교체(2)로 확인했다.
- 개발자·대체 경로: `git archive HEAD`를 컨테이너에 풀어 `sh scripts/install_bot_mcp.sh`를 실행했다. 결과는 `install rc=0`, `env -i …/node --version` = `v22.22.2`, `knowslink/dist/plugin.js`와 `defaults.json` 존재였다. 재실행도 `rc=0`이었고 Node를 다시 받지 않았다("preparing Node" 첫 실행 1회, 재실행 0회).
- 웹(`internal/relay`): `go test ./internal/relay/ -count=1` ok(번역 누락·영어 화면 한국어 잔존 검사 포함), `gofmt -l` 출력 없음.

## 못 한 것

- aarch64 컨테이너는 돌리지 않았다. arm64 분기와 SHA-256 값은 기존 installer 값 그대로다.
- 실제 Grok Bot 계정에서 Import from Repo와 Add를 하지 않았다. 사람 단계이며 SAR-ORCA-BOT-001-OPS에서 한다.
- 배포 저장소 생성·push·tag는 승인 대상이라 하지 않았다.
- 브라우저로 화면을 직접 보지 않았다.

## 사람이 할 절차

1. 소유자가 `tourbut/knowslink-grok-plugin` 저장소를 만든다.
2. D12의 `Grok Bot 플러그인 배포` 절대로 `--grok-release` 트리를 commit하고 `v0.1.1` tag를 붙여 push한다.
3. 팀 관리자가 [Plugins & MCPs](https://cursor.com/dashboard/plugins)에서 Import from Repo로 저장소를 가져온 뒤 Add to Marketplace를 누른다.
4. Bot 소유자가 Setup → Plugins → Add → KnowsLink를 고른다.
5. `knowslink_connect` 동의를 진행한다.
6. 기존 Command 등록이 있던 Bot은 플러그인을 추가하고 agent ID를 대조한 뒤 기존 `knowslink` Command connector를 지운다.
