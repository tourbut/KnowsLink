---
title: KnowsLink
status: draft
updated: 2026-10-14
owner: dev
tasks: [SAR-SYSTEM-DIAGRAM-001-DEV, SAR-REFACTOR-002-OPS]
summary: 회원 가입·친구·Grok 및 Dots 설치와 기존 업데이트를 안내한다
---

# KnowsLink

사람끼리 친구가 되면, 서로 동의한 범위의 하위 에이전트들이 대화하는 서비스다.
웹에서 이름·에이전트·친구를 관리하고, 웹이나 연결한 클라이언트에서 메시지를 보낸다.
서비스 주소는 [link.knowslog.com](https://link.knowslog.com)이다.

## 먼저 읽을 문서

| 목적 | 문서 |
|---|---|
| 가입·친구·대화, Grok/Dots 설치와 업데이트 | 이 README |
| 독자별 전체 목차와 상세 정본 | [문서 목차](docs/README.md) |
| 구성·레포 지도·서명·Inbox·보관 기한 | [구조](docs/architecture.md) |
| 현재 시스템 구성도·Grok/Dots 전달 흐름·소스 근거 | [상세 설계](docs/system-design.md) |
| 개발 도구·로컬 실행·검증·오류 대처 | [개발](docs/development.md) |
| 배포·백업·writer 보호·수신 알림 점검 | [운영](docs/operations.md) |
| 다른 호스트 패키지와 합성/legacy 도구 | [adapter 안내](adapters/README.md) |

## 현재 기능과 확인 범위

2026-10-21 운영 제품은 `4787e85a0b6952970a861b6b100b3561b9cf30a0`다.
Google 로그인, 이름, 사람 친구 수락, 자기 에이전트 연결 동의, 웹/에이전트 연속 대화를 지원한다.
과거 제품 `26094e9`에서 기존 두 계정의 Grok Bot과 Dots 친구 수락·실제 대화 두 차례 왕복을 확인했다.
질문 없이 도착한 수신 알림도 양쪽 실제 화면에서 확인했다. 현재 배포는 해당 관계·키·Remote/Inbox의 상태 digest를 보존했다.
배포 뒤 새 실계정 송수신·화면 알림을 다시 시험한 결과는 아니다.
이는 해당 계정·시험 범위의 결과다. 모든 호스트·계정의 설치 성공을 보장하지 않는다.
근거는 [운영 기록](.fullops-squad/docs/evaluations/qa-reports/SAR-RELEASE-003-OPS/report.md),
[실제 대화 기록](.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-002-LIVE/report.md),
[수신 알림 기록](.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-003-NOTIFY/report.md)이다.

- 발송마다 사용자 승인이 필요하다. 친구 수락·수신·ACK는 자동 답장 승인이 아니다.
- 전달 큐 기한은 최대 180초다. 새 메시지는 독립 기한을 가지므로 대화를 계속 이어갈 수 있다.
- 참여 회원의 웹 대화 사본은 수락 후 최대 24시간 보관한다. 로컬 Inbox·WAL·백업 기한과 다르다.
- `queued`는 서버 수락, `delivered`는 수신 측 보관·ACK를 뜻한다. 읽음·업무 실행·답장 완료와 다르다.
- 캘린더 빈 시간 조회는 배포됐고 운영 키가 설정됐다. 실계정 연결·조회 왕복은 아직 확인하지 않았다. Google 앱이 Testing이면 연결이 7일 뒤 끊긴다.
- 외부 업무 도구 효과(일정 생성 등)·정보 공개 정책·positive silent done은 활성화하지 않았다.
- singleton DB의 공개 규모 성능, 장시간 알림, 앱 재시작 뒤 수신과 다음 Events 갱신은 보장하지 않는다.
- 합성 시험·legacy 정확한 두 agent 관계는 남아 있다. 현재 사람 친구 권한으로 자동 확대하지 않는다.

<a id="google-가입로그인-sar-google-login-001"></a>

## 웹에서 시작하기

웹 사용에는 Go·Node·Docker 설치가 필요 없다. 에이전트 클라이언트 연결은 아래 호스트 절차를 따른다.

1. [서비스](https://link.knowslog.com)를 열고 **Google로 계속**으로 로그인한다. 이메일 로그인이 설정된 환경에서는 확인 코드도 사용할 수 있다.
2. 친구에게 보일 내 이름을 정한다. 이름은 표시용이며 로그인 신원·agent ID를 바꾸지 않는다.
3. **내 에이전트**의 **어디에서 사용할까요?**에서 사용할 곳과 **처음 연결** 또는 **기존 연결 계속**을 고른다. 선택지는 웹에서 바로 사용, Grok Bot (내 컴퓨터), Dots (원격), Claude Code 작업 세션 (내 컴퓨터), 지원 확인이 필요한 호스트다. 화면은 고른 곳의 준비물·권한·다음 행동과 이 문서의 해당 절 링크를 보이며, 고르기만으로 에이전트나 연결 수단을 만들지 않는다. 지원 확인이 필요한 호스트(Codex, Hermes Agent 등)에는 설치 명령을 제공하지 않으므로 지금은 웹 대화를 사용한다. 에이전트 이름을 정하고 클라이언트를 연결한다. 기존 회원은 자기 에이전트 전체 연결에 한 번 명시적으로 동의한다.
4. **친구 추가**에서 상대가 직접 공유한 친구 코드를 입력한다. 이름이나 이메일로 회원을 검색하지 않는다.
5. 상대는 받은 요청을 수락한다. 양쪽 사람의 현재·앞으로 추가할 활성 에이전트가 연결된다는 안내를 확인한다.
6. **새 대화**에서 친구·내 발신 에이전트·상대 수신 에이전트를 고른다. 자기 에이전트끼리는 친구 요청이 필요 없다.
7. 대상과 본문을 확인하고 보낸다. 전달 상태와 웹 기록을 확인한다. 다음 발송도 새 사용자 승인과 독립 기한을 사용한다.

**내 에이전트와 연결 상태**는 서버가 확인한 단계만 보인다: 키 등록 전(준비 대기·지문 확인), 브라우저 동의 완료, 키 등록됨, 원격 앱 복귀·토큰 교환, 이전 시도 만료·취소, 미사용, 결과를 확인하지 못한 경우의 확인 필요. 앱 설치·도구 검색·Claude Code 세션 연결·온라인 여부는 서버가 관측하지 않으므로 앱에서 확인한다. 새로고침·뒤로가기·중복 클릭은 새 에이전트를 만들거나 기존 키를 바꾸지 않는다. 결과가 불확실하면 새로 등록하기 전에 시작한 앱의 저장 결과부터 확인한다.

친구 해제·차단·키 철회는 이후 전달을 막는다. 차단 해제만으로 친구가 복구되지 않는다.
남은 과거 대화 사본은 참여 회원에게 읽기 전용으로 보일 수 있다.
**기존 연결**은 원래 두 에이전트만 허용한다. 사람 전체 연결에는 새 친구 수락이 필요하다.
로그아웃은 브라우저 세션을 끝내며 에이전트 키를 철회하지 않는다.

## Grok 신규 설치

Grok Bot은 해당 대화 컴퓨터에서 **Command MCP**를 실행한다.
Grok CLI의 plugin 설치·doctor 성공은 앱 등록 성공과 별개다.
공식 [Team Bots 안내](https://docs.x.ai/grok-bot/team-bots#set-up-what-the-bot-needs)는 Setup의 Plugins에서 Add하는 경로를 설명한다.
개인 계정의 메뉴·Command 사용 권한은 계정에서 확인해야 한다. 메뉴가 없으면 지원되지 않는 도구 이름을 추측하지 않는다.

### 준비물과 bundle

Bot 컴퓨터의 Linux x86_64 또는 aarch64/arm64, Git, POSIX sh, curl, tar, gzip, sha256sum, Python 3가 필요하다.
installer가 Node `22.22.2`와 bundle을 준비하므로 Bot 컴퓨터에 Go·Make·xz는 필요 없다.
아래 명령은 신규 설치용이다. 기존 설치에는 [업데이트 절차](#grok-기존-설치-업데이트)를 적용한다.
installer는 자신이 관리하는 설치물을 교체하고 임시 stage를 정리한다. 실행 권한을 가진 소유자가 설치 범위를 확인한 뒤 실행한다.

1. 기존 작업·키 폴더와 겹치지 않는 **새 소스 폴더**를 고른다. 다음 폴더가 이미 있으면 다른 이름을 쓴다.

   ```sh
   git clone https://github.com/tourbut/KnowsLink.git /workspace/knowslink-source-new
   cd /workspace/knowslink-source-new
   git checkout --detach 2747fe98653d5c6946f41c219fc256395333d7fe
   git status --short
   ```

2. `git status --short`가 비어 있는지 확인한 뒤 installer를 실행한다.

   ```sh
   sh scripts/install_bot_mcp.sh
   ```

   검증 성공과 `Ready` 출력을 확인한다. 기본 설치 prefix는 `/workspace/.knowslink`다.
   다른 절대경로가 필요하면 실행 시 `KNOWSLINK_PREFIX`를 지정한다. 이후 등록에도 실제 출력 경로를 쓴다.
3. 다음 파일과 Node 버전을 확인한다.

   ```sh
   env -i /workspace/.knowslink/node/bin/node --version
   ls /workspace/.knowslink/knowslink/dist/plugin.js
   ```

   Node 결과는 `v22.22.2`여야 한다. 다운로드·checksum·빌드·검사가 실패하면 다음 단계로 진행하지 않는다.
   installer는 Node와 bundle을 순서대로 준비한다. 실패 시 전체 설치물이 원자적으로 복구된다고 가정하지 않는다.

### 앱 등록과 Google 동의

1. 실제 앱에서 custom MCP server의 **Command** 등록 화면을 연다. 기본값 또는 installer의 실제 출력을 입력한다.

   | 항목 | 기본값 |
   |---|---|
   | Name | `knowslink` |
   | Type | `Command` |
   | Command | `/workspace/.knowslink/node/bin/node` |
   | Arguments | `/workspace/.knowslink/knowslink/dist/plugin.js` |

2. 처음 환경 변수를 비우면 `held`다. `knowslink_status`로 도구 실행을 확인한다. held는 연결 완료가 아니다.
3. 신규 연결 승인을 받은 뒤 다음 **비밀이 아닌 설정**을 Command server 환경에 넣고 MCP를 재시작한다.

   ```text
   KNOWSLINK_MODE=public-node
   RELAY_URL=https://link.knowslog.com
   KNOWSLINK_AGENT_FOLDER=/workspace/.knowslink/connections/grok-new-01
   ```

   마지막 폴더 `grok-new-01`은 아직 없어야 한다. 상위 폴더는 자동 생성한다.
   호스트·인스턴스마다 다른 새 폴더를 쓴다. 교체되는 `/workspace/.knowslink/knowslink` 안에 자격 폴더를 두지 않는다.
4. 사용자에게 이 클라이언트 연결의 명시적 동의를 받은 뒤 `knowslink_connect`를 `{"confirmed":true}`로 호출한다.
5. 사용자가 반환 URL을 자기 브라우저에서 연다. Google 로그인 뒤 실행한 클라이언트와 화면의 지문을 비교한다.
   대상·이름·권한·기한을 확인하고 연결에 동의한다. 다른 사람이 보낸 연결 링크는 승인하지 않는다.
6. `knowslink_connect_status`의 `connected`를 확인한다. 요청은 10분 뒤 만료한다.
   `knowslink_status`의 `configured_unverified` 계열 상태만으로 연결 저장·전달 성공을 판단하지 않는다.
7. 웹의 에이전트 목록과 현재 친구/자기 연결 범위를 확인한다. 승인된 대상·비민감 본문으로만 별도 발송한다.

개인키 `private.pem`, 자격 `agent.json`, 연결 진행 파일은 해당 컴퓨터의 비공개 폴더에 저장한다.
POSIX 폴더는 0700, 파일은 0600을 사용한다. Windows 로컬 연결은 소유자·상속 차단 ACL도 검사한다.
키·credential·Google 비밀번호·연결 token을 채팅·argv·Git·캡처에 복사하지 않는다.
실패·취소·만료된 신규 연결은 기존 폴더를 덮어쓰지 않고 새 폴더에서 시작한다.
승인 뒤 완료 응답을 잃었으면 웹에서 미사용 agent를 확인하고 철회한다.

### 자동 수신과 화면 알림

연결된 `public-node` MCP는 기본 10초 polling으로 수신한다. 프로세스가 살아 있어야 한다.
`KNOWSLINK_RECEIVE_MODE=sse`를 선택하면 도착 신호로 기존 API pull을 깨우며 주기 polling으로 복구한다.
`KNOWSLINK_AUTO_RECEIVE=off`는 자동 수신을 끈다.
SSE·MCP log notice·Inbox 저장·ACK만으로 Grok 채팅에 새 턴이 생기지는 않는다.

상시 watcher가 필요하면 같은 연결 폴더와 기존 환경으로 다음 실행 형태를 사용한다.
이미 watcher가 있다면 새로 중복 실행하지 않는다.

```sh
KNOWSLINK_MODE=public-node RELAY_URL=https://link.knowslog.com \
nohup /workspace/.knowslink/node/bin/node \
  /workspace/.knowslink/knowslink/dist/plugin.js watch \
  /workspace/.knowslink/connections/grok-new-01 \
  >> /workspace/.knowslink/connections/grok-new-01/watch.log 2>&1 &
```

질문 없는 화면 알림에는 별도 webhook routine 설정과 수신 결과를 채팅에 표시하는 지시가 필요하다.
`KNOWSLINK_GROK_WEBHOOK_FILE`은 routine URL/key가 든 소유자 전용 비공개 파일의 절대경로다.
위 실행 환경에 이 경로를 제공한다. URL/key 자체를 명령 인자나 채팅에 넣지 않는다.
기존 검증 계정은 이 방식으로 실제 화면 알림을 확인했다. 새 설치에 자동으로 routine이 생기지는 않는다.
[설정 계약](adapters/README.md#선택형-grok-bot-webhook-routine-wake)과 [문제 판별](docs/operations.md#수신-알림-문제-판별)을 따른다.
미문서화 loopback gateway wake는 별도 선택 기능이며 기본 off다. 기존 webhook을 그 방식으로 바꾸지 않는다.

## Dots 신규 연결

Dots는 **Remote HTTPS MCP**를 사용한다. Grok의 로컬 ZIP·Node·연결 폴더를 설치하거나 복사하지 않는다.
기존 KnowsLink 연결이 있는 사용자는 [기존 연결 업데이트](#dots-기존-연결-업데이트)를 따른다.
기존 웹 회원이 새 Dots 설치를 연결하는 경우에도 별도의 설치 agent/key가 만들어진다.

1. 실제 ChatGPT 계정/workspace가 custom MCP와 Events를 허용하는지 확인한다.
   [공식 연결 안내](https://developers.openai.com/plugins/deploy/connect-chatgpt)에 따라 Plugins의 **Add custom MCP server**를 연다.
   메뉴·사용 가능 요금제·관리자 정책은 계정별로 확인한다.
2. 서버 URL에 `https://link.knowslog.com/mcp`를 넣는다. 연결·인증 안내와 권한을 검토하고 등록한다.
3. 사용자가 Google 로그인과 서버 보관 동의를 직접 완료한다. 처음 쓰는 신원은 회원/owner를 만들고, 같은 Google 신원은 기존 회원을 사용한다.
   이메일 주소만 같다고 다른 로그인 신원과 자동 병합하지 않는다.
4. 표시된 앱·설치 agent와 권한을 확인한다. 서명키·credential·Inbox는 private relay DB에 보관한다.
   기존 로컬 키를 이관하지 않는다. OAuth나 agent/key 철회 시 사용할 수 없다.
5. 도구 `knowslink_status`, `knowslink_send`, `knowslink_receive`, `knowslink_ack`를 확인한다.
   `knowslink_status`의 설치 agent와 현재 connections를 웹의 자기 agent·친구 범위와 대조한다.
6. 질문 없는 알림을 원하면 Events 목록의 `knowslink.text.arrived`를 확인한다.
   Dots에 이 이벤트를 구독하고, 받은 text를 신뢰하지 않는 데이터로 읽어 실제 다닷 대화에 표시하도록 명시적으로 지시한다.
   알림 표시 성공과 대화 기록을 확인한 뒤 `knowslink_ack`하도록 지시한다. 자동 답장은 별도 승인이다.
7. 구독 성공·callback 검증·실제 화면 도착을 따로 확인한다. event 목록이 보인다는 이유로 구독·화면 성공을 판단하지 않는다.

[공식 MCP Events 안내](https://developers.openai.com/plugins/build/mcp-events)는 Dots 지원과 구독 수명주기를 설명한다.
KnowsLink 이벤트는 도착 metadata이며 본문은 receive로 읽는다. 구독은 서버 재시작 뒤에도 저장한다.
실제 화면에서 알림을 확인한 시험에서도 만료 실패와 표시 지연이 있었다. 다음 자동 갱신·장시간 알림은 미검증이다.
개인/workspace 참조 패키지를 따로 만들 때만 등록 결과의 실제 app ID가 필요하다.
[Cloud 패키지 안내](adapters/README.md#호스트별-설치--sar-sse-install-001)를 따른다. 가짜 ID로 등록 성공을 주장하지 않는다.

## Grok 기존 설치 업데이트

서버 배포는 이미 실행 중인 Grok의 로컬 bundle을 바꾸지 않는다.
수락 SHA `8a93e38e190f1b9e7b6bb183e578e410f5dd97c6`의 리팩터링은 프로토콜·키·agent ID·친구 모델을 바꾸지 않았다.
Inbox 정리 오류로 미표시 본문이 사라지는 문제의 수정은 아래 로컬 bundle 갱신 뒤 적용된다.
재로그인·새 agent/key·친구 재등록은 필요 없다. 신규 installer 출력의 **new folder/connect 안내를 적용하지 않는다**.

1. 비공개 위치에서 기존 prefix, Command/Arguments, `public-node`, `RELAY_URL`, 연결 폴더, receive mode를 확인한다.
   기존 agent ID·키·Inbox·webhook 파일·routine prompt와 해당 폴더 watcher PID를 보존한다.
2. [신규 설치의 소스 준비](#준비물과-bundle)처럼 별도 깨끗한 clone에서 수락 SHA를 선택한다.
   기존 prefix로 `sh scripts/install_bot_mcp.sh`를 실행한다. 기본과 다른 prefix는 같은 `KNOWSLINK_PREFIX`를 지정한다.
   자격 폴더가 installer 교체 경로 밖에 있는지 먼저 확인한다.
3. 같은 앱 Command 등록의 MCP 프로세스를 재시작한다. 파일 교체만으로 기존 프로세스에 새 코드가 반영되지는 않는다.
4. `ps -eo pid,args`에서 **정확한 연결 폴더의 watcher PID**를 확인한다. 그 PID만 종료한다.
   넓은 `pkill -f`를 사용하지 않는다. 같은 기존 환경·폴더·webhook 설정으로 새 bundle의 watcher 하나를 시작한다.
5. `knowslink_status`의 autoReceive 상태를 확인한다. 웹/기존 기록과 agent ID가 같은지 대조한다.
   webhook 모드의 `plugin.js wake-check`는 설정 파일 검사이며 실제 발송은 하지 않는다.

설치·재시작 후 새 메시지 시험은 별도 사용자 승인을 받아 수행한다.
Reset/Update 뒤에는 Node·bundle·기존 자격 폴더·watcher를 확인한다. 파일이 없어졌다고 새 agent를 자동 등록하지 않는다.

## Dots 기존 연결 업데이트

이번 서버 구현은 이후 remote 요청부터 적용된다. 기존 endpoint/app·Google 연결·설치 agent/key·친구·Events 자동화를 유지한다.
이번 변경은 remote 도구/Events schema·OAuth metadata 변경이 없으므로 Refresh·재등록·재인증이 필수는 아니다.
이는 구현 diff와 운영 보존 검사에 따른 판단이며 Dots 내부 cache나 다음 갱신 성공을 새로 시험한 결과가 아니다.

향후 도구 이름·설명·schema·annotations·인증·UI resource가 바뀌면 서버 배포 후 **기존 연결에서 Refresh**한다.
Events 목록을 바꾸면 rescan하고 새 대화에서 변경 metadata를 확인한다.
기존 앱을 삭제하거나 Google 설치·친구를 새로 만들지 않는다.
인증 만료 등 별도 오류가 있을 때만 실제 오류에 맞게 복구한다.
근거는 [OpenAI Refresh metadata](https://developers.openai.com/plugins/deploy/connect-chatgpt#refresh-metadata)다.

## 개발환경 설치

개발자용 Go·Node·DB 설치와 실행은 [개발 안내](docs/development.md)를 따른다.
Grok 소유자는 위 installer, Dots 사용자는 remote 등록만으로 시작한다.
개발·격리 QA는 운영 계정·DB·Tunnel과 분리한다. 합성 시험 성공을 실제 계정·자동 답장·업무 실행 성공으로 해석하지 않는다.
