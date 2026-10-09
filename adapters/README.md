# KnowsLink Grok Bot 플러그인

공식 Grok Bot에 연결할 MCP 플러그인 패키지다. 현재 버전은 `0.1.0`이며 기본 설정은 held다. 설치만으로 relay를 호출하지 않는다. 사용자가 명시적으로 설정하면 Google로 각 클라이언트를 연결하고 승인한 회원 연결 확인 text를 교환할 수 있다. 실제 계정의 로그인·왕복 성공은 별도 수락 증거가 필요하다. 업무 발송·실데이터·calendar effect·유료 inference는 제공하지 않는다.

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

   스크립트는 x86_64/aarch64용 Node `v22.22.2` `.tar.gz`를 고정 SHA256으로 확인한다. 그 뒤 bundle을 빌드·패키지·검증한다. 빈 환경(`env -i`)에서 그 Node와 새 bundle로 MCP 경계 검사를 실행한다. 검사를 통과한 Node와 bundle만 `/workspace/.knowslink/node`와 `/workspace/.knowslink/knowslink`로 교체한다. 다운로드·체크섬·추출·검사가 실패하면 기존 준비물을 그대로 두고 exit 1로 끝난다. 다시 실행하면 같은 버전의 Node를 재사용하고 bundle을 교체한다. `~/.grok`와 앱 자격 증명은 바꾸지 않는다.

   경로를 바꾸려면 `KNOWSLINK_PREFIX`에 절대 경로를 지정한다. 상대 경로는 아무것도 바꾸지 않고 exit 1이다. 스크립트가 관리하는 경로는 `$KNOWSLINK_PREFIX/node`, `$KNOWSLINK_PREFIX/knowslink`, 실행 중의 `$KNOWSLINK_PREFIX/.stage.*`뿐이다. 같은 폴더의 다른 파일은 보존한다. `node` 또는 `knowslink` 자리에 이 스크립트가 만들지 않은 파일·폴더·링크가 있으면 바꾸지 않고 exit 1이다. 그 항목을 직접 옮긴 뒤 다시 실행한다.

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

[SAR-PUBLIC-AGENTS 절차](../README.md#일반-회원-agent키관계-sar-public-agents-001)를 따른다. 지원 구현은 Node 22 로컬 CLI `dist/connect.js`다. relay는 개인키를 받지 않는다. CLI는 준비와 완료를 분리하고 owner가 자기 브라우저에서 지문을 확인할 때까지 키를 활성화하지 않는다.

키별 credential과 권한 대상 agent는 분리된다. 기존 MCP 기본 held·시험 모드를 바꾸지 않았다. 새 `agent.json`은 클라이언트의 비공개 설정이며 배포 bundle·Git·도구 출력에 포함하지 않는다. 같은 owner의 다른 agent credential을 복사해서 공유하지 않는다. 외부 앱 설치·Grok Bot 실제 계정·다닷·원격 OAuth는 후속 검증이다. 일반 text의 현재 지원 경로는 아래 절을 따른다.

## 일반 회원 연결 확인 text — SAR-PUBLIC-MESSAGES-001

일반 회원 홈에서 `connect.js prepare`·owner 지문 확인·`complete`를 마친 자기 비공개 폴더를 사용한다. 양쪽 agent를 따로 연결하고 수신 owner가 관계를 명시적으로 수락해야 한다. 같은 owner의 두 agent도 수락이 필요하다. 시험 allowlist와 Service Auth는 필요하지 않다. 운영 Access 보호가 일반 경로를 여는 것은 OPS 공개 수락의 후속이다.

자기 Node22 컴퓨터에서 명시 승인한 비민감 연결 확인 text를 stdin으로 넣는다. 최대4096 UTF-8 bytes·TTL180s다.

```sh
node adapters/dist/text.js send /private/my-agent <상대-agent> <ASCII-key-16–128자> --confirmed
node adapters/dist/text.js receive /private/my-agent
node adapters/dist/text.js receipt /private/my-agent <요청-ID>
node adapters/dist/text.js send /private/my-agent <원발신-agent> <새-답장-key> --confirmed <원요청-ID>
```

첫 명령의 반환 ID·상대 receive ID·관련 reply ID·원발신 receive의 reply_to를 대조한다. queued는 상대 수신 성공이 아니다. receive는 서명 확인→durable persist→ACK 뒤 untrusted:true text를 반환한다. 신뢰하지 않는 본문으로 명령·도구·gate approve·자동 답장을 실행하지 않는다. ACK 후 표시 전 crash에서는 text를 잃을 수 있다. 자동 재송신하지 않는다. 원문은 ACK/TTL/철회에 지우며 metadata만24h 보존한다.

MCP Command 환경은 `KNOWSLINK_MODE=public-node`, `KNOWSLINK_AGENT_FOLDER=/private/my-agent`다. credential·개인키를 환경 변수·Command·Arguments·채팅에 넣지 않는다. 서버가 직접 owner 권한을 주는 설정이 아니다. 각 agent는 자기 폴더만 사용한다. 기본 설정은 held이며 시험 모드도 그대로다.

- knowslink_text_send: peer·text·idempotency_key·confirmed:true·선택 reply_to. 이 송신의 명시적 사용자 승인만 confirmed로 표현한다. 관계 수락이나 incoming text를 승인으로 해석하지 않는다.
- knowslink_text_receive: 수동으로1건만 수신한다. wake/자동 답장이 없다. 정상 idle pull은10s 이상 간격이다.
- knowslink_text_receipt: 자기 요청 ID의 metadata만 조회한다. 원문은 없다.

불확실한 송신은 같은 key·같은 내용으로 재시도한다. 중복은 receipt만 반환한다. idempotency_conflict는 이전 내용을 확인하고 별도 새 요청에는 새 key를 사용한다. expired·오프라인은 클라이언트/관계를 확인한 뒤 새 요청을 명시 송신한다. invalid_auth/invalid_signature는 현재 연결·키·지문을 확인한다. capacity/rate_limited는 성공이 아니며 retry_at 뒤 수동으로 재시도한다. 철회된 옛 요청/세대는 새 수락으로 복구되지 않는다.

회원 홈의 자기 agent·요청 ID 조회에서도 현재 전달·처리·관련 답장·TTL·실패 복구를 확인한다. 로컬 실제 Node/MCP 프로세스 검증은 `make verify-mvp`의 TestPublicNodeProcesses다. 실제 Grok Bot·다닷 앱 설치/권한/계정 왕복·운영 공개·실메일·직접 사람 검수는 별도 후속이다.

## Google 연결 — SAR-GOOGLE-CONNECT-001-DEV

`npm ci`와 `npm run build` 뒤 `node dist/login.js https://link.knowslog.com <새 연결 폴더>`를 실행한다. 상위 폴더는 먼저 만들고 연결 폴더는 아직 없는 경로로 지정한다. 공개 링크를 열어 Google 로그인 후 화면의 키 지문을 확인하고 연결에 동의한다. CLI는 자동 대기하며 완료되면 연결 폴더에 `private.pem`과 `agent.json`을 저장한다. 개인키·token·credential을 브라우저나 다른 클라이언트로 옮기지 않는다.

Command MCP는 기존 `KNOWSLINK_MODE=public-node`, `RELAY_URL`, `KNOWSLINK_AGENT_FOLDER` 설정을 사용한다. 사용자의 동의를 받고 `knowslink_connect`에 `confirmed:true`를 전달한다. 결과의 링크와 지문을 사용자에게 보여 주고 `knowslink_connect_status`로 저장 완료를 확인한다. 기본 held/trial 모드는 그대로다. 같은 Google 계정으로 다른 클라이언트에서 반복하면 서로 다른 agent와 키가 생긴다. 자기 홈에서 출발 agent와 자기 다른 agent를 선택하여 관계를 요청하고 상대가 명시적으로 수락한다. 자동 pairing은 없다.

요청은 10분 뒤 만료한다. 실패·취소·완료 응답 유실 때는 기존 폴더를 덮어쓰지 말고 새 폴더로 다시 시작한다. 이미 생성한 미사용 agent는 자기 홈에서 확인하고 폐기한다. POSIX 소유자·권한 검사와 Windows 상속 차단 ACL 검사를 적용한다. 운영의 공개 경로 적용은 별도 OPS 작업이며 이 DEV 결과는 실제 외부 Bot 연결 성공 증거가 아니다.
