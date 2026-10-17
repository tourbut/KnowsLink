---
title: KnowsLink 호스트별 설치
status: draft
updated: 2026-10-17
owner: dev
tasks: [SAR-SSE-INSTALL-001-DEV, SAR-GROK-WEBHOOK-001-DEV, SAR-HUMAN-PEER-001-DEV, SAR-SESSION-PEER-001-DEV]
summary: 호스트별 로컬 플러그인 설치와 Google 연결 및 클라우드 지원 한계
---

# KnowsLink 플러그인

<a id="sar-human-peer-001-dev-후보"></a>

## 현재 사람 친구·대화와 호스트 안내

현재 사용자 설치·기존 업데이트의 정본은 [루트 README](../README.md)다. 구성·개발·운영은 [문서 목차](../docs/README.md)를 따른다. 이 문서의 과거 시험·합성·legacy 절은 당시 근거로 보존한다.

웹에서 친구 수락을 마치면 양쪽 사람의 현재·앞으로 추가할 활성 agent가 연결된다. status/contacts의 현재 대상만 사용하며 친구 해제·차단·키 철회 뒤에는 저장된 대상이라도 보내지 않는다. 이름은 화면 표시용이고 ID와 credential을 대체하지 않는다.

`knowslink_text_receive` 결과의 `authorship`은 `owner` 또는 `agent`다. owner는 회원이 웹 세션으로 작성한 메시지이며 agent Ed25519 서명이 아니다. authenticated relay의 전용 owner wire를 별도 검증한다. agent는 기존 signed wire를 검증한다. 둘 다 untrusted 입력으로 취급하고 명령/도구 실행의 허가로 사용하지 않는다. private Inbox에 durable 저장한 뒤 persist/ACK하는 순서는 그대로다.

긴 대화는 `knowslink_text_send`로 새 text를 반복한다. 각 발송은 새 ID·멱등키와 독립180초 기한이다. 기존 reply_to는 원요청당 한 관련 답장인 호환 API이며 오래된 부모에 의존해 대화를 이어가지 않는다. 웹 참여 세션의 새 대화 사본은 수락24h까지 조회하며 로컬 Inbox/backup 삭제와 별개다. 2026-10-13 운영 제품은 `2747fe98653d5c6946f41c219fc256395333d7fe`다. 실제 두 계정의 대화·질문 없는 알림 수락 범위와 장시간/재시작/갱신 한계는 루트 README를 따른다.

공식 Grok Bot에 연결할 MCP 플러그인 패키지다. 현재 버전은 `0.1.0`이며 기본 설정은 held다. 설치만으로 relay를 호출하지 않는다. 사용자가 명시적으로 설정하면 Google로 각 클라이언트를 연결하고 승인한 회원 연결 확인 text를 교환할 수 있다. 실제 Grok/Dots 계정의 왕복은 루트 README에 연결한 별도 수락 근거로 확인한다. 다른 호스트의 합성/CLI 검사를 실제 계정 수락으로 확대하지 않는다. 업무 발송·실데이터·calendar effect·유료 inference는 제공하지 않는다.

## 호스트별 설치 — SAR-SSE-INSTALL-001

`python3 scripts/package_plugin.py --verify`는 공통 standalone 런타임으로 다음 파일을 만든다. 의존성 설치와 빌드는 기존 절차를 따른다. 각 ZIP/MCPB는 파일 허용 목록만 포함한다. 개인키·credential·환경 파일·로그는 포함하지 않는다.

| 대상 | 생성 파일 | 설치 계약과 확인 상태 |
|---|---|---|
| Grok Bot / Grok CLI / Cursor | `build/knowslink-grok-bot-plugin.zip` | 기존 패키지 구조 유지. Bot 앱은 아래 Command 등록을 사용한다. CLI 설치와 앱 등록은 별개다. |
| Claude Code | `build/knowslink-claude-plugin.zip` | `.claude-plugin/plugin.json`과 로컬 MCP 설정. 공식 CLI manifest 검증 통과; 실제 Claude 대화 호출은 미검증이다. |
| Claude Desktop | `build/knowslink-claude-desktop.mcpb` | MCPB 0.3의 `manifest.json`이 archive root에 있다. 앱 Extensions에서 설치한다. 실제 앱 설치·내장 Node 호환성은 미검증이다. |
| Codex / ChatGPT desktop 로컬 실행 | `build/knowslink-codex-plugin.zip` | `.codex-plugin/plugin.json`과 repo marketplace. 아래 설치 준비가 필수다. 격리된 Codex CLI 설치·MCP 도구 9개 discovery를 확인했다. |
| Dots / ChatGPT Work Cloud / web | `https://link.knowslog.com/mcp` + 등록 app 참조 ZIP | 기존 Dots 실제 연결·대화·화면 알림을 확인했다. 신규 계정은 별도 Google 동의·설치 agent/key·Events 구독이 필요하다. 장시간/다음 갱신은 미검증이다. |

Claude Code와 Codex는 다음처럼 **아직 없는 절대 설치 경로**로 준비한다. `--host`는 `claude`, `codex`, `grok` 중 하나다. `--connect`를 생략하면 held를 유지한다.

```sh
python3 scripts/package_plugin.py --host codex --install-dir /absolute/knowslink-codex-01 --connect --verify
```

이 명령은 새 설치 폴더만 만든다. 호스트 설정·로그인·relay 연결은 실행하지 않는다. Node는 현재 PATH에서 찾으므로 먼저 지원 버전 `22.22.2`를 선택한다. MCP command와 bundle 경로는 절대 경로로 저장한다. Codex `0.162.0`은 MCP args의 plugin-root 변수를 치환하지 않아 이 준비 단계가 필요하다. 준비한 폴더는 설치 뒤에도 유지한다. 다른 컴퓨터로 옮기면 새 절대 경로로 다시 준비한다.

`--connect`는 `public-node`, 서비스 URL, `<설치 경로>-private`를 설정한다. private 폴더는 설치 폴더 밖에 있으며 아직 생성하지 않는다. 이미 존재하는 설치/private 폴더를 재사용하거나 덮어쓰지 않는다. 기존 설치·키는 보존하고 호스트/인스턴스마다 새 경로를 지정한다. 아래 Google 연결 절차가 성공한 뒤에만 키와 agent가 생성된다.

Claude Code는 `--host claude`로 준비한 뒤 세션 한정으로 `claude --plugin-dir /absolute/knowslink-claude-01/knowslink`를 사용한다. marketplace 방식은 `/plugin marketplace add /absolute/knowslink-claude-01` 후 `/plugin install knowslink@knowslink-plugins`다. 이 작업에서는 실제 사용자 설정을 바꾸지 않았다. [Claude 설치 문서](https://code.claude.com/docs/en/plugins), [manifest와 MCP 경로](https://code.claude.com/docs/en/plugins-reference).

### Claude Code 작업 세션 연결 — SAR-SESSION-PEER-001-DEV

Claude Code 패키지에는 MCP 서버와 함께 함수 hook 모드(`hooks/register.ts`)가 들어 있다. 설치만으로는 어떤 세션도 relay에 연결하지 않는다(SP-01). 준비 단계의 `--connect`는 Claude 설치에 `KNOWSLINK_AUTO_RECEIVE=off`를 함께 쓰므로 MCP 프로세스가 스스로 수신하지 않는다. 검증한 호스트는 Claude Code `2.1.302`(Windows 11)다.

1. `--host claude --install-dir <새 절대경로> --connect`로 준비하고 그 세션에서 `knowslink_connect`로 Google 연결을 마친다. agent 이름은 동의 화면에서 정한다.
2. 연결할 작업 세션에서 `/knowslink-session on`을 실행한다. 모드는 `config/local.mcp.json`의 연결 폴더로 `plugin.js watch <폴더> --session <세션 구분값>` 자식 프로세스를 시작한다. 구분값은 세션 ID 앞 8자다. 출력: `KnowsLink 연결됨: 세션 <구분값> · 내 agent <agent>`.
3. 같은 agent 폴더는 한 번에 한 세션만 받는다(SP-02). 다른 세션이 이미 받고 있으면 `연결 보류`와 그 세션 구분값·시작 시각을 보인다. 소유자가 전환하려면 `/knowslink-session on --take`를 실행한다. 옛 세션은 10초 안에 `연결 해제` toast를 보이고 그 즉시 연결이 풀린다. 이후 그 세션의 text 도구는 모두 거부된다(SP-09). 종료된 세션의 lock은 pid 확인으로 자동 회수한다.
4. 메시지가 도착하면 `KnowsLink 외부 메시지 도착 · 발신 agent <id> · 대기 N건 · 확인: knowslink_text_receive` toast와 상태줄 `KnowsLink <구분값> 대기 N`이 보인다. toast는 도착 알림이며 모델은 본문을 보지 않는다. 본문은 사용자가 요청할 때 `knowslink_text_receive`가 untrusted 데이터로 보여 준다. 작업 중인 turn은 중단되지 않는다(SP-08).
5. `knowslink_text_send`는 호출마다 모드가 소유자에게 `KnowsLink 발송 승인: 내 agent → <peer>. 본문 「…」` 대화상자를 띄운다. 대화상자는 새 메시지인지 답장(`reply_to`)인지, 본문 전체 글자 수와 본문 전체를 잘라내지 않고 보인다. 4096자를 넘는 본문은 전체를 확인할 수 없으므로 묻지 않고 거부한다. `보내기`를 고른 한 번의 시도만 보낸다(SP-05/06). `Other`에 입력한 글은 승인이 아니다. 취소하거나 `-p` 실행이면 발송을 거부한다. 연결하지 않은 세션의 text 도구 호출은 모두 거부한다.
6. `/knowslink-session off`는 자식 프로세스를 끝내고 이 세션의 수신·발송을 막는다. `/knowslink-session status`는 구분값·agent·대기 수를 보인다.

한계: 모드 파일이 바뀌어 플러그인이 다시 로드되면 자식 프로세스가 함께 끝나므로 `/knowslink-session on`을 다시 실행한다. Windows에서 `off`나 세션 종료는 lock 파일을 지우지 않고 pid만 죽는다. 다음 `on`이 죽은 pid의 lock을 회수한다. 슬래시 명령 이름이 `knowslink-session`인 이유는 `/knowslink`가 패키지 skill의 명령이기 때문이다. 지원 수준은 '기존 세션 수신'(도착 toast·세션 결속 로컬 사본, 작업 중·대기 양쪽 관측)이며 본문 표시는 '다음 조회'다. 실제 왕복 근거는 [DEV QA 보고](../.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-DEV-test/report.md)다. 플러그인 자체 검사는 `claude plugin validate adapters`와 `claude plugin test <설치 폴더>/knowslink`다.

Codex는 준비한 폴더를 독립 작업 폴더로 열면 `.agents/plugins/marketplace.json`을 찾을 수 있다. 명령으로 확인하려면 별도 `CODEX_HOME` 디렉터리를 먼저 만들고 각 명령에 같은 경로를 지정한다. 사용자 기본 설정과 인증 파일은 복사하지 않는다. 실제 대화 사용에 필요한 호스트 로그인은 사용자가 수행한다.

```sh
CODEX_HOME=/absolute/knowslink-codex-home codex plugin marketplace add /absolute/knowslink-codex-01
CODEX_HOME=/absolute/knowslink-codex-home codex plugin add knowslink@knowslink-plugins
CODEX_HOME=/absolute/knowslink-codex-home codex mcp list --json
```

`mcp list`는 설정 확인이며 MCP 실행 성공 증거가 아니다. 도구 discovery와 실제 대화 호출을 따로 확인한다. 기본 사용자 설정으로 설치하려는 경우에도 먼저 설치 경로와 권한을 검토한다. [OpenAI 설치 문서](https://developers.openai.com/plugins/build/plugins), [로컬 MCP 지원](https://learn.chatgpt.com/docs/extend/mcp).

Claude Desktop에서는 `.mcpb`를 열거나 **Settings → Extensions**에서 설치한다. 초기 mode는 held다. 연결하려면 extension 설정의 mode를 `public-node`로 바꾸고 새 절대 private folder를 지정한다. 폴더를 extension 설치 디렉터리 안에 두지 않는다. 사용자가 아래 Google 동의 절차를 완료해야 한다. 웹 Claude의 remote connector와 호환성은 이 local extension 검사로 입증하지 않는다. [Desktop 공식 문서](https://support.claude.com/en/articles/10949351-getting-started-with-local-mcp-servers-on-claude-desktop).

Dots Cloud의 신규 연결은 **Plugins → Add custom MCP server**에서 `https://link.knowslog.com/mcp`를 등록하고 연결 검토·Create as a plugin을 거친다. Google 로그인 뒤 표시된 앱·별도 설치 agent/key의 서버 보관에 명시적으로 동의한다. 같은 회원의 자기 agent는 기본 연결 또는 기존 회원의 명시적 범위 동의를 따른다. 다른 사람과는 친구 수락을 별도로 받는다. 기존 로컬 agent/key를 가져오지 않는다. 서버 private DB에 해당 설치의 서명키와 Inbox를 보관하며 OAuth 철회 또는 agent/key 철회 시 사용할 수 없다. 기존 로컬 Google begin/finish는 이 host OAuth 계약과 별개다. 기존 연결의 업데이트에는 신규 등록 절차를 적용하지 않는다. [신규/기존 안내](../README.md#dots-신규-연결)를 따른다. [remote 연결 계약](https://developers.openai.com/plugins/deploy/connect-chatgpt), [인증 계약](https://developers.openai.com/plugins/build/auth).

Cloud 개인/workspace 플러그인은 등록 결과의 실제 app ID를 `.app.json`으로 참조한다. 앱 URL의 `plugin_` 접두사는 제거하고 실제 `asdk_app_…`, `connector_…`, `templated_apps_…` ID를 사용한다. `python3 scripts/package_plugin.py --host openai-cloud --app-id "$KNOWSLINK_REGISTERED_APP_ID" --verify`로 `build/knowslink-openai-cloud-plugin.zip`을 만든다. ID 형식만 검사하며 계정 등록 여부는 host에서 확인해야 한다. 실 ID가 없는 현재 검증은 임시 합성 ID의 구조 self-check만 수행한다. Cloud ZIP에는 stdio process·bundled MCP·credential을 넣지 않는다. bundled MCP import는 HTTPS여도 Desktop only이기 때문이다. 공개 marketplace 제출은 별도 범위다. [패키지 계약](https://developers.openai.com/plugins/build/plugins), [app 참조와 Desktop only 제한](https://learn.chatgpt.com/docs/enterprise/plugin-management), [Dots 컴퓨터·앱](https://learn.chatgpt.com/docs/dots/computers-and-apps).

기본 수신은 `KNOWSLINK_RECEIVE_MODE=poll`이며 10초 polling을 유지한다. `sse`는 `/v1/text/events`의 도착 신호로 기존 API pull을 깨우고 주기적 복구를 유지한다. 본문 수신·검증·private Inbox·ACK는 기존 경로를 사용한다. process가 살아 있을 때만 수신한다. relay SSE와 MCP log notice는 호스트 대화 시작을 보장하지 않는다. Grok의 선택형 gateway wake는 아래 default-off·loopback·metadata-only 규칙을 유지한다.

OpenAI MCP Events는 relay SSE와 별도다. 원격 서버는 MCP 2.0 `2026-07-28` discovery·events methods와 지속 구독·HTTPS callback 검증·서명 webhook을 구현한다. 도착 이벤트는 본문을 포함하지 않으며 `knowslink_receive`가 기존 서명·lease·persist/ACK를 재사용해 설치별 durable Inbox로 받는다. `knowslink_ack`는 읽은 Inbox 항목만 삭제한다. 이벤트는 자동 답장·업무 실행 허가가 아니다. 실제 Dots turn과 이벤트 활성화는 계정·host 정책에 따라 별도로 검증한다. [MCP Events 공식 문서](https://developers.openai.com/plugins/build/mcp-events).

## 도메인으로 Google 연결 시작

앱 Command server 등록 후 다음 환경을 설정하고 MCP를 다시 시작한다. ZIP manifest의 기본 `KNOWSLINK_MODE=held`도 이 값으로 바꾼다. 자격 값은 도구 인자나 채팅에 넣지 않는다.

```text
KNOWSLINK_MODE=public-node
RELAY_URL=https://link.knowslog.com
KNOWSLINK_AGENT_FOLDER=/workspace/.knowslink-connect/nou/new-google-client
```

상위 폴더는 자동 생성하며 마지막 연결 폴더는 아직 없어야 한다. 각 클라이언트는 다른 새 폴더를 사용한다. 사용자 허가 뒤 `knowslink_connect`에 `confirmed:true`를 전달한다. 반환된 URL을 사용자가 자기 브라우저에서 열고 Google 로그인한다. 실행한 클라이언트의 지문과 화면 지문을 비교한 뒤 명시적으로 동의한다. 다른 사람이 보낸 연결 링크는 승인하지 않는다. `knowslink_connect_status`의 `connected`를 확인한다. 개인키·credential은 해당 컴퓨터에만 자동 저장된다. 이 설정에는 Cloudflare Access token이 필요 없다.

root/connect가 Cloudflare 로그인이나 요금제 화면으로 이동하면 서비스 연결이 준비되지 않은 것이다. 관리자에게 [D12의 비용 없는 적용 절차](../.fullops-squad/docs/operations/ops-guide.md#비용-없는-tunnel-적용--sar-google-connect-002-dev)를 전달한다. 플러그인에서 결제·Access 가입을 진행하지 않는다. 관계 수락과 exact text 발송 동의는 연결 뒤에도 별도로 필요하다.

## 패키지 만들기

레포 루트에서 `make install` 뒤 `make plugin`을 실행한다. 빌드 호스트에는 `.nvmrc`의 Node `22.22.2`와 npm `10.9.7`이 필요하다. `adapters/.npmrc`의 `engine-strict=true` 때문에 다른 버전에서는 `npm ci`가 `EBADENGINE`으로 실패한다. Node 20은 2026-04-30에 지원이 끝났으므로 허용 범위를 넓히지 않는다. Grok Bot 컴퓨터에서는 아래 [앱 등록](#grok-bot-앱-등록)의 스크립트가 Node를 준비한다. 다른 Linux 빌드 호스트에 Node 22가 없으면 다음 절차로 공식 배포본을 사용자 영역에 준비한다. `.tar.gz`는 gzip만 필요하다. `xz`가 없는 호스트에서도 동작한다. x86_64는 `x64`, aarch64는 `arm64` 배포본을 사용한다.

```sh
arch=$(uname -m | sed 's/x86_64/x64/; s/aarch64/arm64/')
curl -fsSLO "https://nodejs.org/dist/v22.22.2/node-v22.22.2-linux-$arch.tar.gz"
curl -fsSL https://nodejs.org/dist/v22.22.2/SHASUMS256.txt | grep " node-v22.22.2-linux-$arch.tar.gz\$" | sha256sum -c -
mkdir -p "$HOME/.local/node" && tar -xzf "node-v22.22.2-linux-$arch.tar.gz" -C "$HOME/.local/node"
export PATH="$HOME/.local/node/node-v22.22.2-linux-$arch/bin:$PATH"
```

`build/knowslink-grok-bot-plugin.zip`에는 Grok과 Cursor용 manifest·MCP 설정·skill·설치 문서·bundle·제3자 라이선스 고지가 있다. credential·원문 fixture·node_modules는 포함하지 않는다. bundle에 고정 SDK가 들어 있으므로 설치 대상에는 npm 설치가 필요 없다. MCP 서버는 Grok 프로세스 PATH의 `node`로 실행되므로 그 `node`도 `22.22.2`여야 한다. 생성 명령이 ZIP의 SHA256을 출력하고, 압축 해제한 standalone bundle의 MCP 경계 검사를 실행한다.

ZIP을 빈 준비 폴더에 풀면 다음 구조를 얻는다.

```text
.grok-plugin/marketplace.json     Grok marketplace 색인
.cursor-plugin/marketplace.json   Cursor marketplace 색인
knowslink/
  .grok-plugin/plugin.json        Grok manifest
  .mcp.json                       Grok MCP 설정 (${GROK_PLUGIN_ROOT})
  .cursor-plugin/plugin.json      Cursor manifest
  mcp.json                        Cursor MCP 설정 (${CURSOR_PLUGIN_ROOT})
  package.json
  dist/plugin.js
  skills/knowslink/SKILL.md
  README.md
  THIRD_PARTY_NOTICES.txt
```

## Grok Bot 앱 등록

관측 사실: 이슈1 재시험에서 Bot 컴퓨터의 CLI 설치와 `grok mcp doctor knowslink`는 healthy, 도구 2개였다. 같은 Bot의 앱 MCP 14개에는 knowslink가 없었다. 이슈1 첫 보고에서 같은 Bot은 `AddMcpServer`·`InstallPlugin` 등을 호출할 수 없었다. 도구 검색 결과도 0건이었다.

미확정 가설: 앱 도구 카탈로그는 앱 계정에 등록한 connector를 읽는다. `grok plugin install`이 만든 `~/.grok` 설치는 읽지 않는다. 공식 문서는 Marketplace에서 설치한 connector가 계정 전체에 적용된다고만 쓴다. CLI plugin을 읽지 않는다는 문장은 없다. 실제 계정 재시험 전까지 이 원인은 확정하지 않는다.

공식 근거의 범위는 [Team Bots](https://docs.x.ai/grok-bot/team-bots#set-up-what-the-bot-needs)에 한정된다.

- Team Bot owner는 Bot info pane의 **Setup** 절 **Plugins** 줄에서 **Add**를 고르거나 채팅으로 Bot에게 요청한다.
- Plugins 표의 custom MCP server 종류는 **Remote HTTPS**와 **Command**다. **Command** server는 대화가 사용하는 컴퓨터에서 실행된다.
- 개인 계정 Bot에 같은 Setup 화면이나 Command 종류가 있는지는 공식 문서에 없다.
- 채팅 요청 뒤 **Add MCP Server** 승인 카드가 나온다는 설명은 제3자 글에만 있다. 공식 계약이 아니다.
- 설치한 plugin은 **Marketplace → Your plugins**에서 확인한다. Plugins는 Settings의 절이 아니다([근거](https://docs.x.ai/grok-bot/settings-and-notifications#plugins)).

KnowsLink는 새 공용 서비스 없이 **Command** server 등록을 시도한다. Remote HTTPS는 공개 HTTPS endpoint와 별도 승인이 필요하다. Marketplace 등록은 팀 marketplace 또는 공개 발행과 별도 승인이 필요하다.

1. Bot 컴퓨터의 레포 루트에서 같은 셸로 실행한다. `make`·Go·xz는 필요 없다. `curl`·`tar`·`gzip`·`sha256sum`·`python3`이 필요하다.

   ```sh
   sh scripts/install_bot_mcp.sh
   ```

   스크립트는 x86_64/aarch64용 Node `v22.22.2` `.tar.gz`를 고정 SHA256으로 확인한다. 그 뒤 bundle을 빌드·패키지·검증한다. 빈 환경(`env -i`)에서 그 Node와 새 bundle로 MCP 경계 검사를 실행한다. 검사를 통과한 Node와 bundle만 `/workspace/.knowslink/node`와 `/workspace/.knowslink/knowslink`로 교체한다. 실패하면 다음 등록 단계로 진행하지 않는다. Node 교체는 bundle 검사보다 먼저 일어날 수 있으므로 전체 설치물이 원자적으로 복구된다고 가정하지 않는다. 다시 실행하면 같은 버전의 Node를 재사용하고 bundle을 교체한다. `~/.grok`와 앱 자격 증명은 바꾸지 않는다.

   경로를 바꾸려면 `KNOWSLINK_PREFIX`에 절대 경로를 지정한다. 상대 경로는 아무것도 바꾸지 않고 exit 1이다. 스크립트가 관리하는 경로는 `$KNOWSLINK_PREFIX/node`, `$KNOWSLINK_PREFIX/knowslink`, 실행 중의 `$KNOWSLINK_PREFIX/.stage.*`뿐이다. 같은 폴더의 다른 파일은 보존한다. `node` 또는 `knowslink` 자리에 이 스크립트가 만들지 않은 파일·폴더·링크가 있으면 바꾸지 않고 exit 1이다. 기존 항목을 보존하고 충돌하지 않는 다른 prefix를 선택한다.

   `/workspace` 유지는 보장하지 않는다. 공식 문서는 `/workspace` 파일이 일반 update·recovery 뒤에도 유지되도록 설계됐다고 쓴다. 수동 설치한 패키지는 교체 가능한 것으로 취급하라고 쓴다. **Reset**은 최근 변경을 잃을 수 있다([근거](https://docs.x.ai/grok-bot/computer-and-apps#work-with-files)). update·recover·Reset 뒤에는 다음을 확인한다. 하나라도 실패하면 1단계를 다시 실행한다. 등록 값은 같으므로 앱 등록은 다시 하지 않는다.

   ```sh
   env -i /workspace/.knowslink/node/bin/node --version   # v22.22.2
   ls /workspace/.knowslink/knowslink/dist/plugin.js
   ```

2. 출력의 등록 값으로 앱에 custom MCP server를 추가한다. 아래 화면과 카드는 실제 계정에서 미확인이다. 순서대로 시도한다.
   1. Bot info pane에 **Setup → Plugins → Add**가 있으면 custom MCP server를 고른다. Type **Command**, Command, Arguments를 입력한다. 환경 변수는 비운다.
   2. 그 화면이 없으면 같은 Bot 채팅에서 다음과 같이 요청한다. 승인 카드나 입력 화면이 나오면 Name `knowslink`, Type **Command**, command와 argument, 빈 환경 변수를 확인한 뒤 승인한다.

      ```text
      Add a custom MCP server called knowslink that runs: /workspace/.knowslink/node/bin/node /workspace/.knowslink/knowslink/dist/plugin.js
      No environment variables, no headers.
      ```

   3. 두 경로 모두 Command 등록 수단이 없으면 등록을 멈춘다. Bot이 `AddMcpServer` 같은 도구가 없다고 답할 때도 멈춘다. 아래 회신 항목을 남긴다. Remote HTTPS·팀 marketplace는 별도 승인 대상이다.

   환경 변수가 없으므로 기본 `held`로 실행된다. 환경 변수가 있는 Command server는 owner 자신의 채팅에서만 실행된다. 비밀값은 command·argument·채팅에 넣지 않는다.

3. **Marketplace → Your plugins**의 Installed 목록 또는 Bot info pane의 Setup → Plugins에서 knowslink를 확인한다. 실제 표시 위치는 미확인이다. 새 대화에서 `@knowslink`로 connector를 붙이고 `knowslink_status`만 호출한다. 기대 결과는 `{"state":"held","transport":"pull","actualConnection":"held","webhook":false,"evidenceFetch":false}`다.

등록하지 못했거나 화면이 위와 다르면 다음을 회신한다. 비밀값은 넣지 않는다.

- 앱의 plugin 추가 메뉴에 실제로 보이는 항목 이름(예: URL MCP, Command, Git plugin, 없음)과 Setup 화면 유무.
- 채팅 요청에 대한 Bot 응답 원문과 Bot이 밝힌 사용 가능 도구 목록.
- 승인 카드 또는 입력 화면과 Type 선택지.
- Installed 목록의 knowslink 상태 문구, 앱 버전, 새 대화의 도구 검색 `knowslink` 결과.

제거는 등록한 위치에서 knowslink를 삭제한다. 준비물까지 지우려면 그 뒤 `/workspace/.knowslink/node`와 `/workspace/.knowslink/knowslink`를 삭제한다. 실제 앱 등록·카탈로그 노출·호출은 실제 계정 재시험 전까지 미확인이다.

## Grok CLI 설치 경로 (앱 카탈로그와 별개)

Grok Bot 컴퓨터의 `grok` CLI(1.0.46에서 확인)는 `.grok-plugin/` 또는 `.claude-plugin/` manifest와 `.mcp.json`만 읽는다. `.cursor-plugin/`과 `mcp.json`은 읽지 않는다. 이 둘만 있던 커밋 `0b2d5c6`의 ZIP은 이름 `knowslink-<hash>`, MCP 서버 0개로 설치됐다. 근거는 `grok plugin --help`와 CLI가 함께 설치한 `~/.grok/docs/user-guide/09-plugins.md`다.

레포 루트에서 다음을 실행한다. `grok`과 `node`는 같은 셸의 PATH를 사용한다.

```sh
python3 -m zipfile -e build/knowslink-grok-bot-plugin.zip /tmp/knowslink-plugin
grok plugin validate /tmp/knowslink-plugin/knowslink
grok plugin install /tmp/knowslink-plugin/knowslink --trust
grok plugin list --json
grok mcp doctor knowslink --json
```

| 단계 | 기대 결과 |
|---|---|
| validate | `Plugin manifest is valid.`, `name: knowslink`, `version: 0.1.0`, `MCP servers` |
| install | `Installed 1 plugin(s) from ...: knowslink` |
| list | `"name": "knowslink"`, `"version": "0.1.0"` |
| doctor | `"healthy": true`, `command found`, `handshake OK`, `4 tools discovered` |

install은 파일을 `~/.grok/installed-plugins/`로 복사하므로 준비 폴더는 지워도 된다. `grok mcp list`는 config.toml 서버만 보여 주므로 plugin 서버 확인에는 `grok mcp doctor` 또는 `grok inspect --json`을 사용한다. 새 Grok 세션에서 `knowslink_status`를 호출하면 `held`가 정상이다. 제거는 `grok plugin uninstall knowslink`다.

로컬 marketplace로도 설치할 수 있다. `grok plugin marketplace add /tmp/knowslink-plugin` 뒤 `grok plugin install knowslink --trust`를 실행한다. 새 Git 저장소 공개나 marketplace 발행은 별도 승인 대상이다.

이슈1 재시험에서 CLI 설치만으로는 Bot 앱의 도구 카탈로그에 knowslink가 보이지 않았다. 앱 사용은 위 [앱 등록](#grok-bot-앱-등록)을 따른다. CLI doctor 성공을 앱 Installed 성공으로 표시하지 않는다.

Cursor IDE의 로컬 개발 검사는 [공식 로컬 플러그인 절차](https://cursor.com/docs/plugins#test-plugins-locally)에 따라 `knowslink/`를 `~/.cursor/plugins/local/knowslink`에 복사하고 창을 reload한다. 이는 Cursor 검사이며 Grok 설치 증거가 아니다.

## 로컬 합성 연결 검사

레포의 `npm test --prefix adapters`는 실제 MCP stdio initialize·discovery·held·URL 차단을 검사한다. `make verify-mvp`는 격리 Compose의 실제 Postgres를 사용한다. 합성 두 agent의 pairing·서명·persist·ACK·claim·owner gate approve·최소 denied R을 bundle의 MCP 호출로 검증한다. 자기 시험 자원만 회수한다.

수동 합성 검사에서는 준비 폴더의 Node 프로세스 환경에 다음 값을 제공한다. 비밀값을 chat·도구 인자로 전달하지 않는다.

- `KNOWSLINK_MODE=synthetic-loopback`
- `RELAY_URL=http://127.0.0.1:<시험 포트>`
- `AGENT_CREDENTIAL`, `AGENT_ID`, `AGENT_KID`, `AGENT_KEY_FILE`

`AGENT_KEY_FILE`은 시험용 Ed25519 PEM 파일이다. owner credential은 서버에 전달하지 않는다. 플러그인 `mcp.json`의 기본 `KNOWSLINK_MODE=held`를 유지한다. 해당 합성 검사만 허가됐으면 시험 설정의 모드를 변경한다. hosted `127.0.0.1`은 사용자의 개발 호스트나 운영 relay가 아니다. 일반 production 모드는 활성화할 수 없다. 명시한 시험 remote 모드는 아래 절을 따른다. URL은 HTTP loopback root만 허용하고 userinfo·path·query·fragment와 redirect를 거부한다.

## 도구와 결과

`knowslink_status`는 비밀 파일을 읽거나 네트워크에 접속하지 않는다. `held` 또는 `synthetic_only`를 반환한다. 설정 검증이나 실제 연결 성공을 뜻하지 않는다.

`knowslink_pull_once`는 모델 입력 없이 한 delivery를 처리한다. 동시 호출은 `busy`다. gate는 KnowsLink의 인증된 owner UI에서 결정한다. Grok Bot의 Allow once·chat reply·skill은 owner 승인 증명이 아니다. 승인돼도 disclosure policy와 calendar stub은 유지되며 결과는 denied다. control result는 소비만 하고 재응답하지 않는다. 도구는 원문·gate ID·lease/claim·credential을 모델에 반환하지 않는다.

| state | 의미 |
|---|---|
| held | 실제 연결 비활성; pull은 isError=true |
| unconfigured | 합성 relay/agent 설정 누락 |
| busy | 같은 프로세스의 pull 처리 중 |
| failed | 검증·권한·transport 등 실패; 상세 비밀은 반환하지 않음 |
| empty | pull lease 없음 |
| processed | 안전 처리 경로 완료; 업무 done이나 실제품 연결 성공이 아님 |

실패·TTL 뒤 자동으로 claim을 재발급하거나 요청을 재실행하지 않는다. 필요하면 owner/agent의 승인된 receipt 조회로 transport를 확인한다. `processed`를 silent done으로 바꾸지 않는다.

## 재개 조건

owner는 실제 계정·앱 빌드·설치 정책·hosted runtime을 확인한다. DEV/OPS는 승인된 relay 도달 경로와 최소 credential 전달을 준비한다. coor는 고정 SHA 독립 코드 리뷰와 TESTER QA를 확인한다. 실제 연결은 별도 승인과 위 증거가 모두 확보된 후 재개한다. DEC-02·calendar effect·외부 exactly-once 보류는 유지한다.

Grok package 구조는 `grok plugin validate`와 CLI 사용자 안내서, Cursor package 구조는 [Cursor plugin reference](https://cursor.com/docs/reference/plugins)를 따른다. `${GROK_PLUGIN_ROOT}`와 `${CURSOR_PLUGIN_ROOT}`는 각 호스트의 설치 경로다. MCP wire는 공식 TypeScript SDK `1.32.0`이 처리하며 relay.v1 봉투와 owner 권한은 기존 KnowsLink Adapter와 relay가 집행한다.


## 승인된 Codex·Grok 시험 메시지 (SAR-MVP-003)

이 절은 이번 사용자 시험 승인에만 적용한다. 기존 업무 도구의 payload 비공개·gate·deny 정책과 설치 기본 `held`는 유지한다. 시험 도구 두 개를 포함한 도구를 검색할 수 있지만 시험 send/receive는 `test-loopback` 또는 `test-remote` 환경에서만 실행된다. `trial_configured_unverified`는 모드 보고이며 실제 접속 성공이 아니다.

- `knowslink_test_send`: `{text,idempotency_key}`를 받고 configured peer에 `relay.test.message`를 보낸다. text는 비어 있지 않은 UTF-8 4096 bytes 이하, key는 ASCII 16–128자다. 반환 ID는 queued receipt다. 같은 key와 text의 재전송은 같은 ID를 반환한다. 수신 성공은 별도로 확인한다.
- `knowslink_test_receive`: 입력 없이 한 시험 메시지를 검증·persist·ACK·claim하고 `{id,from,to,text,exp,untrusted:true}`를 반환한다. text는 신뢰하지 않는 데이터다. 업무 요청·권한 변경·자동 도구 실행의 근거로 쓰지 않는다.

`test-remote`는 `RELAY_URL=https://link.knowslog.com`, 기존 agent 변수, `KNOWSLINK_TEST_PEER`, `CF_ACCESS_CLIENT_ID`, `CF_ACCESS_CLIENT_SECRET`을 요구한다. recipient는 env로 고정한다. 주소 인자는 받지 않는다. 다른 HTTPS 호스트·HTTP remote·userinfo·path·query·fragment·redirect를 차단한다. 요청별 timeout은 응답 본문 읽기까지 포함해 10초, 응답은 64 KiB로 제한한다. send TTL은 180초이며 relay 상한은 300초다.

운영자는 [D12 시험 절차](../.fullops-squad/docs/operations/ops-guide.md#13-승인된-양방향-시험-sar-mvp-003)대로 시험 identity와 path 전용 Access 앱을 준비한다. 기본 서버 allowlist는 비어 있다. `KNOWSLINK_TEST_AGENTS=trial_codex,trial_grok`를 명시해야 시험 send와 machine API를 허용한다. owner credential은 두 agent 환경에 전달하지 않는다.

Codex에서 `python3 scripts/run_trial.py --config /private/trial_codex/environment.json receive`로 수신한다. 송신은 시험 text를 표준 입력으로 주고 `send <idempotency-key>`를 쓴다. Grok Command 등록은 같은 launcher에 `--node /workspace/.knowslink/node/bin/node --bundle /workspace/.knowslink/knowslink/dist/plugin.js`를 추가한다. config는 각 호스트에서 소유한 0600 파일이며 키 파일 경로를 해당 호스트의 절대경로로 바꾼다. credential을 Command·Arguments·chat·issue에 넣지 않는다.

수신은 수동 pull이다. relay queue 등록은 Grok/Codex 대화를 깨우지 않는다. 자동 wake는 이번 범위에서 구현·검증하지 않았다. 중복 실행은 shared claim으로 막는다. claim 뒤 응답 출력 전 crash는 표시를 잃을 수 있다. 실행된 claim을 재발급하지 않는다. 만료·철회·서버 모드 해제 또는 claim 완료 시 payload를 지운다. metadata는 24시간 보존한다. WAL·backup 삭제를 보장하지 않는다.

로컬 검증은 `make verify-mvp`의 두 독립 MCP 프로세스와 실제 격리 Postgres 왕복이다. hosted loopback 또는 이 검사를 실제 Grok 계정 왕복으로 표시하지 않는다. 구체 댓글 초안은 [실행 기록](../.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md)을 따른다.

## 일반 회원의 로컬 연결

[현재 클라이언트 연결 안내](../README.md#웹에서-시작하기)를 따른다. 아래 `dist/connect.js`는 수동 prepare/complete와 키 회전의 기존 로컬 경로다. relay는 개인키를 받지 않는다. CLI는 준비와 완료를 분리하고 owner가 자기 브라우저에서 지문을 확인할 때까지 키를 활성화하지 않는다.

키별 credential과 권한 대상 agent는 분리된다. 기존 MCP 기본 held·시험 모드를 바꾸지 않았다. 새 `agent.json`은 클라이언트의 비공개 설정이며 배포 bundle·Git·도구 출력에 포함하지 않는다. 같은 owner의 다른 agent credential을 복사해서 공유하지 않는다. Grok/Dots 실제 연결 근거와 신규 호스트 검증 한계는 루트 README를 따른다. 일반 text의 현재 지원 경로는 아래 절을 따른다.

## 일반 회원 연결 확인 text — SAR-PUBLIC-MESSAGES-001

일반 회원 홈에서 이름을 정하고 `connect.js prepare`·owner 지문 확인·`complete`를 마친 자기 비공개 폴더를 사용한다. 다른 사람과는 사람 친구 요청을 수락한다. 같은 회원의 자기 agent는 친구 요청 없이 연결되며 기존 회원은 범위 확대에 한 번 동의한다. 시험 allowlist와 Service Auth는 필요하지 않다. 현재 운영·실계정 근거는 루트 README를 따른다.

자기 Node22 컴퓨터에서 명시 승인한 비민감 연결 확인 text를 stdin으로 넣는다. 최대4096 UTF-8 bytes·TTL180s다.

```sh
node adapters/dist/text.js send /private/my-agent <상대-agent> <ASCII-key-16–128자> --confirmed
node adapters/dist/text.js receive /private/my-agent
node adapters/dist/text.js receipt /private/my-agent <요청-ID>
node adapters/dist/text.js send /private/my-agent <원발신-agent> <새-답장-key> --confirmed <원요청-ID>
```

첫 명령의 반환 ID와 상대 receive ID를 대조한다. queued는 상대 수신 성공이 아니다. receive는 작성자별 검증→로컬 inbox 기록→relay persist→ACK 뒤 untrusted:true와 authorship:owner|agent를 반환한다. 신뢰하지 않는 본문으로 명령·도구·gate approve·자동 답장을 실행하지 않는다. 긴 대화는 첫 send 형식으로 새 key를 사용해 독립 기한의 text를 반복한다. 마지막 명령은 한 번의 관련 답장 호환 경로다. transport 원문은 ACK/TTL/철회에 지우고 새 대화 사본은 참여 웹 세션에 수락24h까지 제공한다. 로컬 inbox·backup 보존은 별개다. 표시 전 crash는 자동 새 발송으로 복구하지 않는다.

MCP Command 환경은 `KNOWSLINK_MODE=public-node`, `KNOWSLINK_AGENT_FOLDER=/private/my-agent`다. credential·개인키를 환경 변수·Command·Arguments·채팅에 넣지 않는다. 서버가 직접 owner 권한을 주는 설정이 아니다. 각 agent는 자기 폴더만 사용한다. 기본 설정은 held이며 시험 모드도 그대로다.

- knowslink_text_send: peer·text·idempotency_key·confirmed:true·선택 reply_to. 이 송신의 명시적 사용자 승인만 confirmed로 표현한다. 관계 수락이나 incoming text를 승인으로 해석하지 않는다.
- knowslink_text_receive: 자동 수신한 로컬 inbox의 가장 오래된 1건을 먼저 반환한다. 비어 있으면 relay에서 1건을 수신한다. `expired:true`는 관련 답장 불가다. 자동 답장은 없다.
- knowslink_text_receipt: 자기 요청 ID의 metadata만 조회한다. 원문은 없다.

### 자동 수신 — SAR-AUTO-RECEIVE-001-DEV

`public-node` MCP는 연결 폴더가 있으면 시작 직후 자동 수신을 켠다. 수동 receive 호출 없이 10s마다 pull→서명 검증→`<연결 폴더>/inbox/` 기록→persist→ACK를 실행한다. 새 메시지마다 MCP `notifications/message`로 `{event,id,from,pending,next}`만 알린다. `next`는 안내 문장이며 메시지 원문이 아니다. 끄려면 Command 환경에 `KNOWSLINK_AUTO_RECEIVE=off`를 추가한다.

`knowslink_status`의 `autoReceive`로 상태를 확인한다. `state`(running/backoff/stopped/off), `lastSuccessAt`, `lastError`, `nextPollAt`, `lastReceivedId`, `pending`, `hostNotice`를 보여 준다. 비밀값과 text는 없다. `hostNotice:sent_unverified`는 알림을 보냈다는 뜻이다. Grok Bot이 화면에 표시했거나 노우의 턴을 시작했다는 뜻이 아니다.

MCP log notice만으로 Grok 대화 턴을 시작한다고 보장하지 않는다. 기본 설정에서는 다음 대화나 사용자 요청에서 `knowslink_text_receive`로 보관된 메시지를 확인한다. 질문 없는 화면 알림은 아래 webhook routine과 채팅 표시 지시를 설정한 실제 계정에서 확인했다. 신규 설치에 routine이 자동으로 생기지는 않는다. 선택형 [loopback wake](#선택형-grok-bot-loopback-wake--미문서화-gateway)를 켜면 새 메시지 ID를 노우 대화에 알릴 수 있다. 근거와 한계는 [D03 자동 수신 절](../.fullops-squad/docs/design-docs/architecture.md#자동-수신과-호스트-알림--sar-auto-receive-001-dev)을 따른다.

Command MCP process가 대화 사이에도 살아 있는지는 공식 문서에 없다. 대화가 없을 때도 수신하려면 Bot 컴퓨터에서 상시 watcher를 실행한다. MCP와 같은 inbox를 사용한다.

```sh
nohup /workspace/.knowslink/node/bin/node /workspace/.knowslink/knowslink/dist/plugin.js watch <연결 폴더> >> <연결 폴더>/watch.log 2>&1 &
```

watcher 출력에는 ID·발신자·대기 수만 있다. 같은 폴더에 MCP 자동 수신과 watcher가 함께 돌면 각 process가 10s마다 pull한다. 중복 표시는 ID로 막는다. 재시작 뒤에도 미확인 메시지는 inbox에 남는다. 읽은 메시지는 원문 없는 표시만 24h 남는다.

### 선택형 Grok Bot loopback wake — 미문서화 gateway

> 경고: 이 기능은 Grok Bot의 미문서화 내부 API를 사용한다. 업데이트 뒤 route·파일 경로가 바뀔 수 있다. port `1340`을 공개하거나 SSH tunnel·Tailscale로 외부에 열지 않는다. gateway token을 출력·복사·채팅·issue에 넣지 않는다.

지원 상태는 다음과 같다(2026-10-10 확인).

| 항목 | 상태 | 근거 |
|---|---|---|
| 외부/컴퓨터에서 노우 대화 wake | 대화 wake는 공식 지원 없음. webhook trigger routine은 공식 스킬 문서로 확인한 후보(아래 webhook 절) | Cursor staff Colin(Community Support Engineer)의 [2026-08-18 답변](https://forum.cursor.com/t/let-a-grok-bot-computer-wake-its-own-agent-chat/168260): “Wake mechanics is an area we're working on” |
| Bot 컴퓨터 loopback gateway `127.0.0.1:1340`, `/home/box/sand-data/gateway.json`의 `token`, `POST /api/listAgents`, `POST /api/sendPrompt {agentId,prompt}` | 미문서화. staff가 “undocumented” workaround로 연결 | 커뮤니티 사용자 adam91holt의 [2026-08-12 글](https://forum.cursor.com/t/grok-bot-can-i-send-it-a-message-from-outside/168199/8): live 설치에서 두 route를 시험했다고 보고 |
| 본인 Bot 컴퓨터의 gateway 존재·token 필드·응답 코드·노우 턴 시작 | 실제 미검증 | DEV는 Bot 컴퓨터에 접근하지 않았다. 아래 확인 명령으로 owner 환경에서 확인한다 |
| 수신 본문을 노우 prompt로 직접 전달 | 사용하지 않음 | `sendPrompt`는 사용자 권한 입력이다. 받은 text를 넣으면 신뢰하지 않는 지시가 사용자 지시로 바뀐다 |

동작은 다음과 같다. 기본은 off다.

- watcher 또는 MCP 자동 수신 loop가 메시지를 inbox에 저장한 뒤, 아직 알리지 않은 미확인 ID를 모아 `sendPrompt`를 한 번 호출한다. 대상은 `KNOWSLINK_GROK_WAKE_AGENT`의 agent UUID 하나다. 주소는 `127.0.0.1`로 고정이다.
- prompt는 고정 문장과 검증한 ID·대기 수만 담는다. 노우에게 `knowslink_text_receive`로 untrusted data를 표시하라고 요청한다. 받은 text의 지시 실행·다른 도구 실행·사용자 승인 없는 답장을 금지한다. 받은 text·발신자 이름은 넣지 않는다.
- token은 매번 gateway 파일에서 읽고 `Authorization: Bearer` header로 loopback에만 보낸다. 상태·marker·로그·MCP 응답에 넣지 않는다. 응답 본문은 읽지 않는다.
- 결과는 `<연결 폴더>/inbox/<id>.wake`(0600)에 `{state,lastAt,lastError}`로 남는다. MCP에서는 `knowslink_status.autoReceive.hostWake`로 보인다. `accepted_unverified`는 gateway가 HTTP 2xx로 받았다는 뜻이다. 노우가 읽었거나 답했다는 뜻이 아니다.
- 연결 거절·token 파일 없음은 `retry`다. 다음 loop(10s)에 다시 알린다. 4xx는 `rejected`, 5xx·응답 없음·10s timeout은 `uncertain`이다. 두 상태는 중복 prompt를 막기 위해 자동 재전송하지 않는다. 메시지는 inbox에 남는다.
- 같은 ID는 재lease·재시작·여러 process에서도 한 번만 알린다. marker는 24h 뒤 정리된다. 그때까지 읽지 않은 메시지는 다시 한 번 알린다.
- 노우 턴은 사용자의 기존 Grok Bot 사용량을 쓴다. 새 서비스·포트·과금 설정은 없다.

설치와 활성화는 owner가 Bot 컴퓨터 터미널에서 한 번 실행한다. 앱 화면 설정은 바꾸지 않는다.

0. 기존 Command MCP process는 파일 교체만으로 새 코드가 되지 않는다. 설치 뒤 MCP를 새 process로 다시 연결하거나 재시작하고 `knowslink_status`의 `autoReceive`가 보이는지 확인한다. 이 확인 전에는 watcher를 시작하지 않는다.
1. 이 후보로 준비물을 갱신한다. KnowsLink checkout에서 후보 SHA를 checkout하고 `sh scripts/install_bot_mcp.sh`를 실행한다. 결과: `/workspace/.knowslink/knowslink/dist/plugin.js`가 교체된다.
2. gateway와 agent ID를 읽기 전용으로 확인한다. prompt는 보내지 않는다.

   ```sh
   /workspace/.knowslink/node/bin/node /workspace/.knowslink/knowslink/dist/plugin.js wake-check
   ```

   결과: `{"state":"ok","agentIds":[...]}`. 응답의 UUID만 출력한다. 이름·token은 출력하지 않는다. `gateway_unreachable`·`gateway_token_unavailable`·`gateway HTTP 401`이면 이 기능을 사용할 수 없다. 그대로 보고한다.
3. 노우 대화의 agent ID를 고른다. 같은 명령에 `KNOWSLINK_GROK_WAKE_AGENT=<UUID>`를 붙여 `"configuredListed":true`를 확인한다.
4. 이 연결 폴더의 기존 watcher만 멈추고 wake 설정으로 다시 시작한다. 먼저 `ps -eo pid,args | grep 'plugin.js watch <연결 폴더>'`로 정확한 PID를 확인한다. 다른 연결 폴더의 watcher를 끝내지 않는다. `pkill -f`처럼 넓은 일치는 쓰지 않는다.

   ```sh
   kill <확인한 PID> ; KNOWSLINK_GROK_WAKE_AGENT=<UUID> nohup /workspace/.knowslink/node/bin/node /workspace/.knowslink/knowslink/dist/plugin.js watch <연결 폴더> >> <연결 폴더>/watch.log 2>&1 &
   ```

5. 상대 agent가 연결 확인 text를 보낸다. 그 송신 ID로 `<연결 폴더>/inbox/<ID>.wake`의 `state`를 확인한다. 노우 대화에 doorbell prompt가 나타나고 노우가 같은 ID를 표시했는지를 따로 확인한다.

port나 파일이 다르면 `KNOWSLINK_GROK_GATEWAY_PORT`, `KNOWSLINK_GROK_GATEWAY_FILE`(절대경로)로 바꾼다. 형식이 틀린 설정은 `invalid_config`로 wake만 멈추고 자동 수신은 계속한다. 중지는 해당 연결 폴더 watcher의 정확한 PID만 확인해 끝내고 `KNOWSLINK_GROK_WAKE_AGENT` 없이 다시 시작한다. MCP Command 환경에도 같은 변수를 넣을 수 있지만 watcher 하나로 충분하다.

한계: gateway가 꺼진 동안 같은 컴퓨터의 다른 process가 port를 점유하면 token을 받을 수 있다. Bot 컴퓨터를 owner 단독으로 사용하는 경우에만 켠다. 컴퓨터 Update/Reset은 watcher와 준비물을 지울 수 있다. 그 뒤 1–4단계를 다시 실행한다.

<a id="선택형-grok-bot-webhook-routine-wake--후보-실제-미검증"></a>

### 선택형 Grok Bot webhook routine wake

`sendPrompt` doorbell은 노우 대화에 사용자 메시지로 보인다. 공식 [make-bot-ui 스킬](https://github.com/cursor/plugins/blob/main/pstack/skills/make-bot-ui/SKILL.md)의 webhook trigger routine은 `<webhook_event>` routine 턴을 시작한다. 이 경로는 기존 실제 Grok 계정에서 질문 없는 화면 알림을 확인했다. [알림 근거](../.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-003-NOTIFY/report.md)의 시험 범위만 적용한다. 장시간·재시작 성공은 미검증이다. 기존 webhook 설정을 loopback 방식으로 바꾸지 않는다.

- 비공개 파일 하나에 routine 화면의 URL과 sender key를 넣는다. 파일은 소유자 전용(0600)이어야 한다. 링크·그룹 읽기 권한이 있으면 사용하지 않는다.

  ```json
  { "url": "https://api2.cursor.sh/automations/webhook/<routine 화면의 ID>", "key": "<sender key>" }
  ```

- watcher를 `KNOWSLINK_GROK_WEBHOOK_FILE=<절대경로>`로 시작한다. 이 변수가 있으면 그 process는 webhook만 울린다. `KNOWSLINK_GROK_WAKE_AGENT`가 함께 있어도 `sendPrompt`는 호출하지 않는다. webhook 실패 뒤에도 `sendPrompt`로 넘어가지 않는다.
- URL은 위 형식만 허용한다. 다른 host·http·query·port·userinfo는 거부한다. redirect는 따르지 않는다. 요청은 `Authorization: Bearer <key>`와 `X-Automation-Key: <key>`를 담은 POST JSON 한 번이며 timeout은 8s다.
- body는 `{"event":"knowslink_text_arrived","ids":[UUID…],"pending":N,"next":"…"}`뿐이다. 받은 text·발신자는 보내지 않는다. routine prompt는 `body`를 untrusted data로 읽고 `knowslink_text_receive`를 호출하도록 owner가 작성한다.
- 결과·marker·재시도 규칙은 loopback과 같다. 2xx는 `accepted_unverified`, 연결 거절·DNS 실패·설정 파일 불량은 `retry`, 3xx·4xx는 `rejected`, 5xx·timeout은 `uncertain`이다. URL·key는 상태·marker·로그에 남지 않는다.
- `plugin.js wake-check`는 webhook 모드에서 파일만 검사한다. 결과: `{"state":"ok","mode":"webhook"}`. 아무것도 보내지 않는다. 파일 검사와 실제 routine 턴·화면 표시는 별개다. 신규 시험 발송은 별도 사용자 승인을 받아 확인한다.
- 중지는 해당 watcher의 정확한 PID를 끝내고 변수 없이 다시 시작하는 것이다. routine 사용량은 기존 Bot 사용량을 쓴다. 새 유료 서비스는 없다.

불확실한 송신은 같은 key·같은 내용으로 재시도한다. 중복은 receipt만 반환한다. idempotency_conflict는 이전 내용을 확인하고 별도 새 요청에는 새 key를 사용한다. expired·오프라인은 클라이언트/관계를 확인한 뒤 새 요청을 명시 송신한다. invalid_auth/invalid_signature는 현재 연결·키·지문을 확인한다. capacity/rate_limited는 성공이 아니며 retry_at 뒤 수동으로 재시도한다. 철회된 옛 요청/세대는 새 수락으로 복구되지 않는다.

회원 홈의 자기 agent·요청 ID 조회에서도 현재 전달·처리·관련 답장·TTL·실패 복구를 확인한다. 로컬 실제 Node/MCP 프로세스 검증은 `make verify-mvp`의 TestPublicNodeProcesses다. 실제 Grok/Dots 왕복·화면 알림은 루트 README의 별도 근거를 따른다. 이 합성 검사는 다른 호스트 계정·실메일·장시간·직접 사람 검수를 대신하지 않는다.

## Google 연결 — SAR-GOOGLE-CONNECT-001-DEV

`npm ci`와 `npm run build` 뒤 `node dist/login.js https://link.knowslog.com <새 연결 폴더>`를 실행한다. 상위 폴더는 먼저 만들고 연결 폴더는 아직 없는 경로로 지정한다. 공개 링크를 열어 Google 로그인 후 화면의 키 지문을 확인하고 연결에 동의한다. CLI는 자동 대기하며 완료되면 연결 폴더에 `private.pem`과 `agent.json`을 저장한다. 개인키·token·credential을 브라우저나 다른 클라이언트로 옮기지 않는다.

Command MCP는 기존 `KNOWSLINK_MODE=public-node`, `RELAY_URL`, `KNOWSLINK_AGENT_FOLDER` 설정을 사용한다. 사용자의 동의를 받고 `knowslink_connect`에 `confirmed:true`를 전달한다. 결과의 링크와 지문을 사용자에게 보여 주고 `knowslink_connect_status`로 저장 완료를 확인한다. 기본 held/trial 모드는 그대로다. 같은 Google 계정으로 다른 클라이언트에서 반복하면 이름을 정한 별도 agent와 키가 생긴다. 자기 agent는 신규 회원 기본 연결 또는 기존 회원의 명시적 범위 동의를 따르고, 다른 사람은 친구 수락을 별도로 받는다.

요청은 10분 뒤 만료한다. 실패·취소·완료 응답 유실 때는 기존 폴더를 덮어쓰지 말고 새 폴더로 다시 시작한다. 이미 생성한 미사용 agent는 자기 홈에서 확인하고 폐기한다. POSIX 소유자·권한 검사와 Windows 상속 차단 ACL 검사를 적용한다. 현재 공개 경로와 실제 외부 Bot 연결 결과는 루트 README의 운영·실계정 근거를 따른다. 이 로컬 연결 구현 검사만으로 실제 계정 성공을 주장하지 않는다.
