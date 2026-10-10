---
title: SAR-SSE-INSTALL-001 호스트 설치와 공식 계약 검증
status: draft
updated: 2026-10-11
owner: dev
tasks: [SAR-SSE-INSTALL-001-DEV]
summary: 로컬 플러그인 준비와 공식 host 지원 범위 및 remote MCP Events 구현 격차
---

# SAR-SSE-INSTALL-001 — 설치 구현과 공식 근거

## 범위와 기술 계획

기준 ref는 `7a6906851b872fc52310a0ac6623527a4b0adb2b`다. 실행 지시서는 `handovers/to_dev.md` 전체를 확인했다. `fullops-common-0.3.3`, `project.md`, `rules/delegation.md`, 문서 작성 규칙을 적용했다. Test level은 lite다. 이 하위 작업에서 재위임하지 않았다. 설치 파일 소유권이 좁고 부모가 소스·테스트를 병렬 구현하기 때문이다.

기존 `scripts/package_plugin.py`와 `install_bot_mcp.sh`, `verify_grok_plugin.py`, adapter manifest·skill·README 및 connect 도구의 환경 계약을 확인했다. 변경 영향은 package caller, host별 manifest/config, archive 경로, 기존 MCP smoke, 설치 문서다. 소스·테스트·package.json·기술/운영 정본·Git 작업은 부모 소유다.

계획은 기존 standalone bundle과 Google begin/finish를 재사용하는 것이다. Grok/Cursor의 기존 ZIP 이름과 파일 배치를 보존한다. Claude Code, Codex, Claude Desktop은 다른 설치 계약이므로 별도 archive를 생성한다. 준비 단계는 새 절대 폴더와 private sibling만 설정한다. 키·credential·사용자 설정을 복사하거나 덮어쓰지 않는다. 실제 호스트 등록과 Google 승인은 설치 준비와 분리한다.

새 의존성은 없다. 표준 Python ZIP/JSON/path/subprocess와 기존 SDK `1.32.0` bundle을 사용한다. UI 변경과 디자인 검사는 해당 없다.

## 공식 문서 확인

2026-10-11 KST에 아래 URL을 실제 열고 해당 절을 읽었다. OpenAI Docs 스킬을 적용했다. Context7의 `resolve-library-id`로 Claude Code와 OpenAI Codex 설치 계약을 요청했지만 `Monthly quota exceeded`가 반환됐다. library ID를 얻지 못했으므로 query-docs 성공으로 기록하지 않는다. 공식 페이지와 실제 설치 CLI로 대체했다.

| 실제 열린 공식 URL | 확인한 절 | 구현에 적용한 사실 |
|---|---|---|
| [OpenAI Package your plugin](https://developers.openai.com/plugins/build/plugins) | Install a local plugin manually; Plugin structure; Bundled MCP servers | Codex compatibility manifest와 repo marketplace를 사용한다. public submission은 local ZIP 준비와 별개다. |
| [OpenAI Authentication](https://developers.openai.com/plugins/build/auth) | Custom auth with OAuth 2.1; MCP authorization spec requirements; Support the authorization-code flow | Remote host OAuth는 resource/authorization metadata, host client 등록, PKCE와 token 검증 계약이다. 로컬 Google agent 등록으로 대체할 수 없다. |
| [OpenAI MCP Events](https://developers.openai.com/plugins/build/mcp-events) | Before you start; Advertise event support; Create a subscription; Sign the request; Stop a subscription | MCP 2.0 `2026-07-28`, event methods, persistent subscription, callback 검증과 Standard Webhooks 서명이 필요하다. Work web, desktop Cloud, Dots 지원을 명시한다. |
| [OpenAI Connect and test](https://developers.openai.com/plugins/deploy/connect-chatgpt) | Prepare the endpoint; Add the MCP server; Test the complete plugin | Cloud host 연결은 public HTTPS 또는 Secure MCP Tunnel로 별도 연결한다. local package 설치 성공과 account 실행 검증을 구분한다. |
| [Meet dots](https://learn.chatgpt.com/docs/dots) | Make it your own; Connected apps and computers; Tasks and conversations | Dot은 cloud에서 실행한다. 연결한 컴퓨터에서 local task를 만들 수 있다. 로컬 파일 접근 권한은 별도다. |
| [Dots computers and apps](https://learn.chatgpt.com/docs/dots/computers-and-apps) | Connect your computer; Connect apps | supported plugin 접근은 실행 환경과 권한에 달려 있다. local skills에는 연결한 컴퓨터가 필요하다. 컴퓨터는 online이며 ChatGPT app이 열려 있어야 한다. |
| [OpenAI Plugin management](https://learn.chatgpt.com/docs/enterprise/plugin-management) | Desktop-only plugins; Reference an existing app with .app.json | bundled MCP import는 HTTPS여도 Desktop only다. cloud 등록 app mapping은 실제 등록 ID가 있어야 하며 이번 작업에서 만들지 않았다. |
| [OpenAI MCP](https://learn.chatgpt.com/docs/extend/mcp) | Supported MCP features; Plugin-provided MCP servers | 로컬 Codex host는 stdio MCP를 지원한다. 설정 목록과 실제 MCP handshake는 다른 검사다. |
| [OpenAI App Server](https://learn.chatgpt.com/docs/app-server) | Initialization; API overview | initialize/initialized 뒤 mcpServerStatus/list로 inference 없이 실제 MCP discovery를 검사했다. turn/start는 호출하지 않았다. |
| [Claude plugin reference](https://code.claude.com/docs/en/plugins-reference) | mcpServers; Environment variables | Claude Code는 plugin MCP args의 CLAUDE_PLUGIN_ROOT를 지원한다. 업데이트 시 경로가 바뀌므로 private state를 plugin root에 저장하지 않는다. |
| [Claude plugins](https://code.claude.com/docs/en/plugins) | plugin installation/local loading | local plugin 개발·설치와 별도 session 로딩 경로를 확인했다. |
| [Claude Desktop local MCP](https://support.claude.com/en/articles/10949351-getting-started-with-local-mcp-servers-on-claude-desktop) | Install desktop extensions; Developer FAQ | Desktop extension은 MCPB이며 built-in Node를 제공한다. 실제 app runtime 버전은 별도 확인해야 한다. |
| [MCPB manifest spec](https://github.com/modelcontextprotocol/mcpb/blob/main/MANIFEST.md) | Manifest Schema; server; user_config; compatibility | manifest 0.3, Node entry_point, __dirname 치환과 비밀값 없는 사용자 설정을 적용했다. |
| [Grok Bot settings](https://docs.x.ai/grok-bot/settings-and-notifications#plugins) | Plugins | Bot app plugin 설정과 CLI 설치를 구분했다. |
| [Grok Team Bots](https://docs.x.ai/grok-bot/team-bots#plugins) | Plugins; Where the work runs | Command는 해당 대화의 computer에서 실행한다. 환경 변수가 필요한 Command는 owner 대화로 제한될 수 있다. command/args에 credential을 넣지 않는다. |

## 지원 표

| 호스트/환경 | 공식 지원 | 이번 결과 | 실제 미검증/미구현 |
|---|---|---|---|
| Grok Bot | custom MCP Command와 Remote HTTPS | 기존 Command installer 보존, Google 별도 폴더 안내 | Bot 앱 등록·실제 Google consent와 왕복 |
| Grok CLI / Cursor | 기존 host plugin manifest | Grok isolated plugin validate/install/list/doctor 검사 | Cursor 실제 계정 설치; CLI 성공을 Bot 앱 설치로 표시하지 않음 |
| Claude Code | local plugin + stdio MCP | 전용 ZIP/manifest/config와 설치 준비, 공식 CLI validate | Claude 대화에서의 도구 호출·Google consent |
| Claude Desktop | local MCPB extension | root manifest 0.3의 .mcpb와 non-secret 설정 | 실제 앱 설치, 내장 Node 22 지원, 계정 왕복 |
| Codex local CLI | local plugin/marketplace + stdio MCP | 독립 CODEX_HOME 설치, 절대 경로로 도구 9개 discovery | 사용자 계정 대화 turn·Google consent |
| ChatGPT desktop local | local plugin + stdio host | Codex 설치물 후보를 재사용 | desktop UI 설치/도구 호출 미검증 |
| Dots connected-computer local task | local task와 supported plugin, local skill 권한 | 로컬 Codex 작업 후보 경로 문서화 | Dots에서 이 custom local MCP가 직접 사용되는지 미검증; cloud 직접 설치로 주장하지 않음 |
| Dots / Work Cloud / ChatGPT web | registered remote MCP/eligible plugins; MCP Events | 공식 계약과 제한 확인 | Remote MCP transport, OAuth 2.1, app 등록, Events subscription/webhook 모두 미구현 |

공식 근거가 있으므로 Dots 자체를 미지원 제품으로 표시하지 않는다. 현재 KnowsLink remote 기능의 부재를 부모에게 status로 전달했다. 기존 relay SSE 신호와 OpenAI Events webhook을 혼동하지 않는다. Event 연동은 protocol·인증·구독 저장·외부 HTTPS 보안 설계가 필요하므로 이 설치 하위 범위에서 추가하지 않는다.

## 변경 결과

- `package_plugin.py --host all`은 기존 Grok/Cursor ZIP, Claude Code ZIP, Codex ZIP, Claude Desktop MCPB를 생성한다. allowlist·고정 timestamp·라이선스 고지를 유지한다.
- `--install-dir`는 `grok|claude|codex` 한 호스트를 새 절대 경로에 준비한다. command와 bundle args를 실제 절대 경로로 쓴다. 기존 폴더와 private 폴더가 있으면 거부한다.
- `--connect`는 `public-node`, 공개 서비스 URL과 `<설치경로>-private`를 설정한다. 로그인은 시작하지 않으며 private folder도 만들지 않는다. 각 설치 인스턴스는 서로 다른 키/agent를 기존 begin/finish 도구로 생성한다.
- 원본 Codex ZIP은 설치 template다. Codex `0.162.0`은 MCP args의 `CLAUDE_PLUGIN_ROOT`, `PLUGIN_ROOT`, `pluginRoot`를 치환하지 않았다. CLI add/list는 성공했지만 discovery가 실패했다. 절대 경로 baseline은 도구 9개를 발견했다. 따라서 Codex 설치에는 위 준비 단계를 필수로 문서화했다.
- `verify_grok_plugin.py`는 사용자 HOME/config/env를 상속하지 않는 임시 환경으로 설치한다. 기존 잘못된 7 tools 기대값을 현재 9 tools로 고쳤다.
- 기본 poll10s와 `KNOWSLINK_RECEIVE_MODE=sse` opt-in을 안내한다. SSE는 arrival 신호이며 본문은 기존 API로 받는다. log notice·수신·host turn은 서로 다른 상태다. Grok wake는 default-off/loopback/metadata-only를 유지한다.

## 검증 기록

이 하위 검사는 Windows, Python 3.13, 로컬 Node `v24.13.0`에서 실행했다. Node 24 실행은 추가 smoke이며 제품 지원 버전을 넓히는 근거가 아니다. 배포/빌드는 기존 Node `22.22.2` 고정을 유지한다. 부모가 pinned Node와 전체 필수 검사를 담당한다.

| 검사 | 결과/한계 |
|---|---|
| `python3 scripts/package_plugin.py --verify` | exit 0. 4개 추출 bundle의 실제 MCP initialize/discovery, held/실패 경계 smoke 통과. CLI 설치 준비의 새 private folder 미생성·기존 설치/private 경로 거부 통과. |
| `python3 scripts/verify_grok_plugin.py` | isolated Grok CLI validate/install/list/doctor와 installed bundle smoke 통과. 도구 9개 발견. 실제 Bot 앱 등록 아님. |
| `claude plugin validate <extracted-claude-plugin> --json` | exit 0, errors/warnings 없음. 실제 tool start 증거는 아님. |
| isolated `codex plugin marketplace add`, `plugin add`, `mcp list --json` | exit 0. 초기 CODEX_HOME 폴더 누락 실패를 고쳐 새 임시 폴더 생성 후 재검사했다. temp helper aliases 경고는 발생했으며 호스트 기본 설정은 건드리지 않았다. |
| isolated Codex app-server `initialize`, `initialized`, `mcpServerStatus/list` | plugin-root 변수 3종은 handshake 실패. 절대 경로 baseline과 최종 installer가 준비한 실제 설치에서 각각 knowslink 도구 9개와 toolsError 없음. connect/connect_status 포함을 확인했다. inference turn은 호출하지 않았다. |
| MCPB 구조/manifest | archive root와 entrypoint, 기본 held 사용자 설정, __dirname args, credential 제외 정적 확인 통과. 공식 mcpb CLI schema validator와 실제 Desktop 설치는 미실행이다. |
| package allowlist/security | 비밀값·env·credential·node_modules 제외. 설치는 archive와 별도 private sibling을 사용한다. |

최종 installer의 Codex isolated native install/discovery와 Claude validator 재확인은 exit 0이다. `git diff --check`와 4개 신규 JSON의 Prettier 검사도 통과했다. 초기 Codex temp directory 정리는 Windows process의 잠깐 남은 파일 잠금으로 한 번 실패했다. 이후 직접 executable 실행과 종료 대기를 사용했으며 실패를 제품 검사 통과로 바꾸지 않았다.

최종 생성물 SHA256은 다음과 같다. 부모가 bundle이나 문서를 다시 빌드하면 값이 달라지므로 최종 후보의 hash는 부모가 다시 기록한다.

| 생성물 | SHA256 |
|---|---|
| knowslink-grok-bot-plugin.zip | `cb900501f26fb745c3154df9baabf472cb22eb6c6baea739d51ad49b397f59ed` |
| knowslink-claude-plugin.zip | `9bc65ad13e34fe097833b8a1d4e963323490d8742330ca6474129353c5140f75` |
| knowslink-codex-plugin.zip | `5905ce44d6c5d4727ad9560547a55749831a5173a84cb03f1b22072bc3bb680a` |
| knowslink-claude-desktop.mcpb | `ade17c83b40d60829970c247c05a5ec7fe5995fcec526e02c26998d4e303d53e` |

생성물은 `build/` 아래 재생성 가능한 산출물이며 커밋하지 않는다. 전체 `make lint`, `make test`, FullOps exact-SHA lint·커밋·독립 리뷰는 부모가 실행한다. 이 하위 작업은 shared inbox를 archive하거나 harness를 수정하지 않는다.

## 남은 수락

부모가 소스·테스트·문서 통합과 pinned Node 검사를 수행한다. coor/ops는 후보 수락 뒤 실제 호스트 설치·Google 수동 로그인·도구 discovery/왕복을 검증한다. 운영 DB, production 배포, 키/관계, 전역 설정, local Docker와 유료 API는 변경하지 않았다. Remote Dots/ChatGPT MCP와 MCP Events는 별도 미구현 상태로 남는다.


## 부모 통합에서 해결한 remote 격차

위 미구현 표는 설치 하위 작업 단독 결과다. 부모의 같은 과제에서 remote OAuth2.1/MCP2·Events 구현과 app ID 기반 Cloud 참조 패키지를 추가했으므로 최종 지원 상태는 DEV·REMOTE 보고서 및 adapters/README.md가 정본이다. 실제 계정 등록·배포·Google consent·Dots turn은 여전히 미실행이다. 하위 작업의 당시 실패와 hash는 보존하며 최종 bundle hash는 부모가 다시 기록한다.
