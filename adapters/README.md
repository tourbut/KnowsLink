# KnowsLink Grok Bot 플러그인

공식 Grok Bot에 연결할 MCP 플러그인 패키지다. 현재 버전은 `0.1.0`이며 실제 계정 연결은 held다. 설치만으로 relay를 호출하지 않는다. 명시적으로 선택한 로컬 합성 모드만 작동한다. 업무 발송·실데이터·calendar effect·유료 inference는 제공하지 않는다.

## 패키지 만들기

레포 루트에서 `make install` 뒤 `make plugin`을 실행한다. `build/knowslink-grok-bot-plugin.zip`에는 standalone marketplace와 `knowslink/` 플러그인이 있다. manifest·MCP 설정·skill·설치 문서·bundle만 포함한다. credential·원문 fixture·node_modules는 포함하지 않는다. bundle에 고정 SDK가 들어 있으므로 설치 대상에서는 Node `>=22.22.2 <23`만 필요하다. 생성 명령이 ZIP의 SHA256을 출력한다.

ZIP을 빈 준비 폴더에 풀면 다음 구조를 얻는다.

```text
.cursor-plugin/marketplace.json
knowslink/
  .cursor-plugin/plugin.json
  mcp.json
  package.json
  dist/plugin.js
  skills/knowslink/SKILL.md
  README.md
```

## 공식 Grok Bot 설치 경로

Grok Bot은 Cursor 계정을 사용하며 계정의 여러 Bot이 플러그인과 컴퓨터를 공유한다. [공식 Bot 플러그인 연결](https://cursor.com/help/grok-bot/connect-plugins)과 [Bot connector policy](https://docs.x.ai/grok-bot/teams-and-enterprises)를 따른다. 실제 앱 빌드와 hosted Node는 이번 과제에서 확인하지 않았다.

1. owner/admin은 검토된 ZIP과 SHA256을 받는다. 별도 승인 후 준비된 marketplace 폴더를 팀이 접근 가능한 Git 저장소에 등록한다. 이 과제는 새 저장소 공개나 marketplace 발행을 수행하지 않는다.
2. admin은 Cursor Dashboard의 Plugins & MCPs에서 Team Marketplaces의 Import from Repo로 준비 저장소를 가져온다. KnowsLink는 Default Off로 둔다. [공식 team marketplace 절차](https://cursor.com/docs/plugins#add-a-team-marketplace)를 참고한다.
3. owner는 Grok Bot의 Marketplace 또는 Plugins에서 KnowsLink를 추가한다. Installed 목록과 `knowslink_status`, `knowslink_pull_once` 검색을 확인한다. Disabled by team admin이면 admin이 해당 플러그인 정책을 확인한다.
4. 먼저 `knowslink_status`를 호출한다. 기본 결과 `held`가 정상이다. 설치 성공·도구 검색·실제 relay 연결은 각각 별도로 기록한다.

팀 marketplace가 없는 개인 계정은 공식 marketplace 배포와 검토가 선행한다. 이 저장소의 ZIP을 앱에 바로 업로드하는 경로는 공식 자료에서 확인하지 못했다. 로컬 준비 파일을 Installed 성공이라고 표시하지 않는다. Bot 앱이 packaged stdio 서버를 지원하지 않거나 hosted Node를 제공하지 않으면 해당 앱 버전과 도구 검색 실패를 기록한다. 공식 remote MCP 배포는 network·인증·실행 승인을 받은 후속 과제다. 기존 베타를 임의로 remote MCP로 노출하지 않는다.

Cursor IDE의 로컬 개발 검사는 [공식 로컬 플러그인 절차](https://cursor.com/docs/plugins#test-plugins-locally)에 따라 `knowslink/`를 `~/.cursor/plugins/local/knowslink`에 복사하고 창을 reload한다. 이는 Cursor IDE 검사이며 Grok Bot 설치 증거가 아니다. 이번 과제는 전역 설치 폴더를 수정하지 않았다.

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

공식 package 구조는 [Cursor plugin reference](https://cursor.com/docs/reference/plugins)를 따른다. `${CURSOR_PLUGIN_ROOT}`는 설치 경로다. MCP wire는 공식 TypeScript SDK `1.32.0`이 처리하며 relay.v1 봉투와 owner 권한은 기존 KnowsLink Adapter와 relay가 집행한다.
