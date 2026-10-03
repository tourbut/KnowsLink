# KnowsLink Grok Bot 플러그인

공식 Grok Bot에 연결할 MCP 플러그인 패키지다. 현재 버전은 `0.1.0`이며 실제 계정 연결은 held다. 설치만으로 relay를 호출하지 않는다. 명시적으로 선택한 로컬 합성 모드만 작동한다. 업무 발송·실데이터·calendar effect·유료 inference는 제공하지 않는다.

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

Grok Bot 앱의 도구 카탈로그는 계정에 등록된 connector만 읽는다. 등록 경로는 Marketplace plugin과 custom MCP server 두 가지다. custom MCP server의 종류는 **Remote HTTPS**와 **Command**다. **Command** server는 대화가 사용하는 Bot 컴퓨터에서 실행된다. 근거는 [Team Bots — Plugins](https://docs.x.ai/grok-bot/team-bots#plugins)와 [Connect an app](https://docs.x.ai/grok-bot/computer-and-apps#connect-an-app)이다. 아래의 `grok plugin install`은 Bot 컴퓨터의 Grok CLI(`~/.grok`)에만 설치한다. 앱 카탈로그에는 등록하지 않는다. 이슈1 재시험의 결과는 CLI doctor healthy, 앱 MCP 14개 중 knowslink 0개였다.

KnowsLink는 새 공용 서비스 없이 **Command** server로 등록한다. Remote HTTPS는 공개 HTTPS endpoint와 별도 승인이 필요하다. Marketplace 등록은 팀 marketplace 또는 공개 발행과 별도 승인이 필요하다.

1. Bot 컴퓨터의 레포 루트에서 같은 셸로 실행한다. `make`·Go·xz는 필요 없다. `curl`·`tar`·`gzip`·`sha256sum`·`python3`이 필요하다.

   ```sh
   sh scripts/install_bot_mcp.sh
   ```

   스크립트는 x86_64/aarch64용 Node `v22.22.2` `.tar.gz`를 고정 SHA256으로 확인한다. 그 뒤 `/workspace/.knowslink/node`에 풀고 bundle을 빌드·패키지·검증한다. bundle은 `/workspace/.knowslink/knowslink`에 둔다. 마지막으로 빈 환경(`env -i`)에서 그 Node와 bundle로 MCP 경계 검사를 실행한다. `/workspace`는 컴퓨터 update·recovery 뒤에도 유지되는 공유 작업 폴더다. 경로를 바꾸려면 `KNOWSLINK_PREFIX`를 지정한다. `~/.grok`와 앱 자격 증명은 바꾸지 않는다. 다시 실행하면 같은 Node를 재사용하고 bundle을 교체한다.

2. 출력의 값으로 앱에 custom MCP server를 추가한다. 같은 Bot 채팅에서 다음과 같이 요청한다.

   ```text
   Add a custom MCP server called knowslink that runs: /workspace/.knowslink/node/bin/node /workspace/.knowslink/knowslink/dist/plugin.js
   No environment variables, no headers.
   ```

   **Add MCP Server** 승인 카드에서 Name `knowslink`, Type **Command**, 위 command와 argument, 빈 환경 변수를 확인한 뒤 승인한다. 환경 변수가 없으므로 기본 `held`로 실행된다. 환경 변수가 있는 Command server는 owner 자신의 채팅에서만 실행된다. 비밀값은 command·argument·채팅에 넣지 않는다.

3. **Marketplace → Your plugins**의 Installed 목록(앱 버전에 따라 **Settings → Plugins**)에서 knowslink를 확인한다. 새 대화에서 `@knowslink`로 connector를 붙이고 `knowslink_status`만 호출한다. 기대 결과는 `{"state":"held","transport":"pull","actualConnection":"held","webhook":false,"evidenceFetch":false}`다.

제거는 같은 목록에서 knowslink를 삭제한다. 실제 앱 등록·호출은 실제 계정 재시험 전까지 미확인이다. 승인 카드의 정확한 필드명은 공식 문서에 없다. 카드가 Command 종류를 제공하지 않으면 카드 화면과 오류를 그대로 회신한다.

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
| doctor | `"healthy": true`, `command found`, `handshake OK`, `2 tools discovered` |

install은 파일을 `~/.grok/installed-plugins/`로 복사하므로 준비 폴더는 지워도 된다. `grok mcp list`는 config.toml 서버만 보여 주므로 plugin 서버 확인에는 `grok mcp doctor` 또는 `grok inspect --json`을 사용한다. 새 Grok 세션에서 `knowslink_status`를 호출하면 `held`가 정상이다. 제거는 `grok plugin uninstall knowslink`다.

로컬 marketplace로도 설치할 수 있다. `grok plugin marketplace add /tmp/knowslink-plugin` 뒤 `grok plugin install knowslink --trust`를 실행한다. 새 Git 저장소 공개나 marketplace 발행은 별도 승인 대상이다.

CLI 설치는 Bot 앱의 도구 카탈로그에 등록하지 않는다. 앱 사용은 위 [앱 등록](#grok-bot-앱-등록)을 따른다. CLI doctor 성공을 앱 Installed 성공으로 표시하지 않는다.

Cursor IDE의 로컬 개발 검사는 [공식 로컬 플러그인 절차](https://cursor.com/docs/plugins#test-plugins-locally)에 따라 `knowslink/`를 `~/.cursor/plugins/local/knowslink`에 복사하고 창을 reload한다. 이는 Cursor 검사이며 Grok 설치 증거가 아니다.

## 로컬 합성 연결 검사

레포의 `npm test --prefix adapters`는 실제 MCP stdio initialize·discovery·held·URL 차단을 검사한다. `make verify-mvp`는 격리 Compose의 실제 Postgres를 사용한다. 합성 두 agent의 pairing·서명·persist·ACK·claim·owner gate approve·최소 denied R을 bundle의 MCP 호출로 검증한다. 자기 시험 자원만 회수한다.

수동 합성 검사에서는 준비 폴더의 Node 프로세스 환경에 다음 값을 제공한다. 비밀값을 chat·도구 인자로 전달하지 않는다.

- `KNOWSLINK_MODE=synthetic-loopback`
- `RELAY_URL=http://127.0.0.1:<시험 포트>`
- `AGENT_CREDENTIAL`, `AGENT_ID`, `AGENT_KID`, `AGENT_KEY_FILE`

`AGENT_KEY_FILE`은 시험용 Ed25519 PEM 파일이다. owner credential은 서버에 전달하지 않는다. 플러그인 `mcp.json`의 기본 `KNOWSLINK_MODE=held`를 유지한다. 해당 합성 검사만 허가됐으면 시험 설정의 모드를 변경한다. hosted `127.0.0.1`은 사용자의 개발 호스트나 운영 relay가 아니다. production·remote 모드는 현재 코드에서 활성화할 수 없다. URL은 HTTP loopback root만 허용하고 userinfo·path·query·fragment와 redirect를 거부한다.

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
