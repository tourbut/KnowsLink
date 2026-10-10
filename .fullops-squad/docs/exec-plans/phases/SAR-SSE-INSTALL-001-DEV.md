---
title: SSE 신호와 Google 인증 호스트 설치 구현
status: draft
updated: 2026-10-11
owner: dev
tasks: [SAR-SSE-INSTALL-001-DEV]
summary: 기존 pull을 유지한 SSE 및 호스트별 설치와 Dots 원격 인증 경로의 구현·검증 기록
---

# SAR-SSE-INSTALL-001-DEV

## 기준과 책임

base는 `7a6906851b872fc52310a0ac6623527a4b0adb2b`다. 지시서 준비 commit은 `4786a00`이다. fullops-common-0.3.3, project.md, 문서 작성 규칙, delegation.md, Ponytail full을 적용한다. 테스트는 lite, 하위 위임은 standard다. 실제 부모 구현 세션은 `01a12678-994d-7680-a15c-83577943b52d`, task는 `task_c76a5ef82349`, dispatch는 `ctx_e4490c0d8bef`다.

부모는 SSE/adapter 수신과 기술·운영 문서를 소유한다. Orca child Run `run_5b1823734740`에서 installer는 scripts·manifest·skill·README만, remote worker는 `remote_*.go`·필요 Google callback 연동만 소유한다. Git 커밋은 부모만 수행한다. 독립 고정 SHA 리뷰와 tester QA는 coor가 후속 배정한다.

## 기술 계획과 구현


기존 singleton transaction에서 큐·현재 agent 자격 변화만 감지하고 빈 PostgreSQL NOTIFY를 commit과 함께 보낸다. SSE는 LISTEN 뒤 현재 상태를 확인하며 자기 agent에 대한 `ready`/`arrival` 빈 신호만 전송한다. 공유 JSON의 단기 입장 기록으로 전체16·agent당2 연결을 제한하고 일반 HTTP 입장 슬롯은 handshake 뒤 반납한다. heartbeat15초에는 DB를 읽지 않고55초 연결 수명과 취소로 자원을 회수한다. API pull·서명·Inbox·persist/ACK는 그대로 재사용한다.

Adapter는 기본 polling10초를 보존하고 `KNOWSLINK_RECEIVE_MODE=sse`에서 신호·reconnect catch-up과 동일 polling 복구를 사용한다. 자동·수동 pull은 기존 serial을 공유하고 신호는 coalesce한다. Grok wake 기본 off·metadata-only 경계를 유지한다.

설치 소유권은 Orca child Run `run_5b1823734740`에 분리했다. 하위 작업은 scripts의 패키징/설치, adapter manifest·skill·README 및 공식 근거 보고서만 수정한다. 부모는 Go·adapter 수신 소스/테스트·기술/운영 정본·커밋을 소유한다. 실제 호스트 설치·Google 수동 동의·배포는 후속이며 새 유료 API와 local Docker는 사용하지 않는다. Linux 고유 임시 DB container만 격리 검증에 사용한다.


SSE 기능 commit은 `9f4c3b6`이다. singleton transaction은 현재 public-agent 자격과 유효 queued text ID의 hash가 달라질 때만 빈 `pg_notify('knowslink_text','')`를 같은 transaction에 넣는다. PostgreSQL은 commit 후 알림을 전달한다. 별도 queue·본문·자격을 NOTIFY payload에 넣지 않는다. rate·HTTP admission·세션 Seen·stream bookkeeping은 알림을 만들지 않는다.

`GET /v1/text/events`는 기존 Bearer/현재 회원 권한으로 자기 agent만 구독한다. 기존 boundedHTTP의 rate/handshake 슬롯을 사용한 뒤 일반 요청 슬롯을 반납한다. shared JSON Streams의 전체16·agent당2 한도가 모든 relay process를 제한한다. 각 stream은 pool 밖의 LISTEN connection 하나를 갖는다. 연결 수명55초, 입장 기록65초, setup/read/write3초, heartbeat15초다. heartbeat에는 DB read/transaction이 없다. LISTEN을 먼저 확정한 뒤 상태를 읽어 연결과 send 사이의 누락을 복구한다. 신호 때는 DB 현재 시계·키·회원·관계·세대·TTL을 다시 검사한다. 잘못된키·자기 agent 철회는 연결을 거부하거나 닫는다. 취소·timeout은 connection과 입장 기록을 정리하며 DB 장애 때65초 만료로 회수한다.

클라이언트 기본은 기존 poll10초다. `KNOWSLINK_RECEIVE_MODE=sse`는 `ready`/`arrival`의 빈 JSON만 받아 동일 AutoReceiver를 깨운다. timer와 중복 신호를 합치고 auto/manual pull은 같은 serial을 쓴다. 재연결 첫 ready는 pending을 당겨오며 기존 poll/backoff/retry_at는 장애·누락 복구를 계속 담당한다. TextTransport의 원문 API·서명·private Inbox fsync/rename·persist/ACK는 보존한다. MCP close는 stream과 재접속 timer를 취소한다. 신호는 실행·답장 승인 또는 호스트 turn 증거가 아니다. Grok loopback wake는 기본 off와 metadata-only를 유지한다.

## Dots 경계 확인

공식 OpenAI Docs의 [MCP Events](https://developers.openai.com/plugins/build/mcp-events) `Before you start`를 직접 열었다. Dots/Work Cloud는 MCP2.0 `2026-07-28` 및 HTTPS webhook을 요구한다. relay SSE나 stdio logging notification을 그 계약으로 표시하지 않는다. [Package your plugin](https://developers.openai.com/plugins/build/plugins)의 호스트별 패키징도 열어 확인했다.

tracked ask의 coor 답변 결과에서 coor는 Dots를 미구현으로 남겨 전체 완료하지 말고 같은 과제에 최소 remote OAuth/MCP·신호 경로를 구현하도록 확인했다. 이 답변은 기존 승인 대상의 기술 차이를 해결하는 범위다. 새 유료 서비스·운영 변경 권한을 부여하지 않는다. 원격 구현과 공식 근거는 별도 REMOTE 보고서에 남긴다. 실제 계정 설치는 후보 수락 뒤 coor/사용자가 수행한다.

## API 근거와 의존성

Context7의 MCP SDK 및 pgx resolve는 monthly quota exceeded로 실패했다. 기존 SDK1.32.0의 protocol onclose 선언과 stdio close 코드를 직접 확인했다. pgx5.10.0을 유지한다. [pgx Conn.WaitForNotification](https://pkg.go.dev/github.com/jackc/pgx/v5#Conn.WaitForNotification), [PostgreSQL17 NOTIFY](https://www.postgresql.org/docs/17/sql-notify.html)를 직접 열어 확인했다. 새 제품 의존성은 없다.

## 검증과 후속

Windows Go unit 및 integration-tag 컴파일, Node22.22.2 check/build·events.test.js가 exit0이다. events 테스트는 실제 loopback HTTP와 서명된 envelope로 initial/push/reconnect/coalesce/직렬pull/fallback/cancel·metadata schema를 검사한다. 기존 inbox.test는 polling·복사-before-ACK·중복·재시작·오류·wake 경계를 보존한다. Linux 격리 Postgres와 필수 최종 SHA 검사 결과는 완료 전에 아래에 기록한다.

UI는 SSE 경로에서 변경하지 않는다. 별도 디자인 lint/테마는 프로젝트에 없으므로 적용하지 않는다. 원격 OAuth 동의 화면은 기존 memberStyle과 보안 정책을 재사용한다. 합성 화면의 DEV 직접 검수와 실제 Google/독립 UI 수락을 구분한다. 운영 Google 계정·DB·키·관계·Tunnel·호스트 전역 설정은 변경하지 않는다.

SIZE-002: SSE source와 핵심 회귀를 하나의 기능 commit으로 분리했다. 구현의 연결 자원·권한·취소·동시성 경계와 실제 HTTP/DB 테스트가 필요해300줄/400추가 경고를 넘는다. 광범위 스타일 정리나 새 추상화는 하지 않는다. 설치 및 remote는 별도 기능 commit으로 보존한다.


## 통합 검증 근거

검증 원문은 [DEV evidence](../../evaluations/dev-checks/SAR-SSE-INSTALL-001/evidence.json)에 보관한다. Linux는 SSH 기존 host key 검증을 유지하고 `/tmp/knowslink-sse001-ctxe4490c0d/final`의 clean clone을 사용한다. `.env.server` 값은 메모리에만 읽고 출력·commit하지 않았다. Postgres17 고정 image의 고유 container는 loopback 임시 port·256MiB·1CPU로 실행하고 회수했다. 실행 전후 기존 container 이름/ID가 일치했다. local Docker·운영 DB·서비스·Google 설정·키·관계·Tunnel·호스트 사용자 설정은 변경하지 않았다.

| 검사 | 결과와 실제 범위 |
|---|---|
| Windows Go unit / integration-tag compile | exit0. 실제 DB 실행과 구분한다. |
| pinned Node22.22.2 check/build/events | exit0. initial/push·중복100·직렬 auto/manual·재접속·503 fallback·cancel·부적합 body frame 거부. |
| 4개 local archive `package_plugin.py --host all --verify` | exit0. 추출 standalone MCP initialize/discovery/held·Google connect tools·경로/기존 private 폴더 거부. 실제 account 설치가 아니다. |
| Cloud package synthetic self-check | valid app ID 형식·app reference·portable/compatibility identity·30자 subtitle·stdio/MCP import 미포함 및 invalid ID 거부 exit0. 임시 합성 archive는 회수했고 실제 app archive는 등록 ID가 없어 생성하지 않았다. |
| Linux `make install`, `make build` | exit0. Node22.22.2·lockfile 의존성을 사용한다. |
| Linux 실제 격리 DB integration | `TestRemoteOAuthHTTP`, `TestTextEventsHTTP`, `TestPublicMessageHTTPBoundaries`, `TestPublicNodeProcesses`를 `go test -race -tags integration`으로 실행해 exit0. 515eabb의 제품 코드와 후속 security 수정의 재검사 SHA를 worker_done/gate 기록으로 구분한다. |
| 실제 cloudflared2026.9.1 ingress validate/rule | exit0. SSE/MCP/metadata/token/정확 consent는 relay, owner/admin/잘못된 MCP suffix는404. candidate render만 수행했다. |
| 등록 `lint.py --from 7a6906851b872fc52310a0ac6623527a4b0adb2b` | bd411fb에서 product-lint `make lint`0, product-test `make test`0, ERROR0/WARNING19/실행불가0. archive 뒤 최종 HEAD에서 다시 실행하며 최종 exact SHA·종료코드는 worker_done과 Git 디렉터리 `fullops-gate/pass.json`이 정본이다. |

Windows에는 Linux make 실행 환경이 없으므로 task-local `make.cmd`가 같은 HEAD와 clean 여부를 확인한 Linux clone에서 실제 make 명령을 실행하고 실제 종료코드를 전달한다. repo의 lint 설정·exclude·보안 수준은 바꾸지 않는다. 검사 출력의 tail은 요약 표시만 하며 command 종료코드는 별도로 보존한다.

### 실패와 수정

1. 최초 실제 SSE HTTP 검사는 정상 서버의 10초 read/write deadline 때문에15초 heartbeat 전에 연결이 닫혔다. `10168a4`에서 handshake 뒤 inherited socket deadline을 해제하고 stream write마다3초 deadline·55초 lifetime을 적용했다. 같은 HTTP/DB 검사가 통과했다.
2. Postgres initdb의 임시 Unix socket readiness가 DB ping보다 먼저 성공했다. 임시 harness를 TCP readiness로 고쳐 최종 service 준비를 기다렸다. 제품 코드 결함으로 표시하지 않는다.
3. 새 clean clone의 integration 실행 전에 adapter build를 빠뜨려 `dist/connect.js` 누락으로 실패했다. `make build` 뒤 전체 선택 integration을 재실행해 통과했다.
4. 최초 FullOps gate는 `Inbox.pending()`를 `pending()` 테스트 보류 표식으로 오인했다. gate를 완화하지 않고 private Inbox의 실제 저장 파일을 검사하는 assertion으로 변경했다. 제품 lint/test는 최초에도 통과했으며 gate 재실행은 ERROR0이다.
5. 최종 security 검토에서 refresh token 재사용 시 이전 토큰만 거부하는 부족한 처리를 찾았다. RFC9700 §4.14.2에 맞춰 검증된 기존 refresh가 설치 token family를 철회하도록 보완했다. 거절 응답에도 철회가 transaction에 commit되는 실제 HTTP 회귀를 추가했다. [OAuth 보안 BCP](https://www.rfc-editor.org/rfc/rfc9700.html#section-4.14.2)를 직접 열어 확인했다.

### UI 검수와 지원 상태

DEV는 Chrome에서 현재 Go template의 합성 값을 렌더링해 직접 확인했다. 최초 Google 확인/승인 화면을390×844·1280×800에서 확인했고 기존 memberStyle 재사용 뒤390×844에서 긴 callback URI의 줄바꿈·가로 overflow 없음(scrollWidth390)·승인/취소 버튼을 확인했다. 실제 Google 동의나 실제 relay 계정 승인은 실행하지 않았다. 최종 스타일의 desktop preview는 로컬 미리보기 연결 오류 뒤 background helper 자동승인 차단으로 미실행이다. 보안 보호를 완화하거나 우회하지 않았으며 독립 tester가 실제 consent 화면을 확인해야 한다. 별도 디자인 lint/테마는 프로젝트에 없어 기존 CSS 재사용 범위만 적용한다.

공식 지원 URL·절·CLI 증거는 [INSTALL](SAR-SSE-INSTALL-001-INSTALL.md), OAuth·Events의 세부 계약은 [REMOTE](SAR-SSE-INSTALL-001-REMOTE.md)가 정본이다. 설치 하위 보고의 당시 remote 미구현 항목은 부모의 remote 구현으로 해결했고 해당 이력은 덮어쓰지 않았다. Grok CLI·Claude validator·Codex app-server discovery는 실제 앱 계정·대화 호출을 대신하지 않는다. Claude Desktop MCPB app runtime과 actual Dots OAuth/Events callback 성공·task turn은 후속이다.

### 의존성·규모·경고 처리

새 제품 의존성은0이다. package.json은 테스트와 manifest 검사 목록만 바뀌었다(DEP-001). SDK1.32.0·pgx5.10.0과 Go/Python 표준 라이브러리를 유지했다. 사이즈 경고의 기존 core/text/mcp/http/store 파일은 공통 caller·기존 상태·수신 API에 필요한 좁은 연동만 추가했으며 전체 재구성은 하지 않았다. remote_auth/events는 새 OAuth·SSRF·callback·claim/retry 신뢰 경계를 기능별 파일로 분리했고 주요 거부·상태보존·HTTP 검증을 포함해 예상 약1000줄을 넘었다. 기능 commits로 SSE/timeoutfix/host installer/remote/security를 나눴다. 상세 auth/events 내부를 다시 분할해 검증·권한 흐름을 숨기지 않았다(SIZE-001/002). PLANS 크기 경고는 coor 준비 commit의5줄이며 DEV는 진행표를 수정하지 않았다.

SEC-001은 합성 unit의 `whsec_` protocol prefix+32 zero byte key다. 실제 credential이 아니다. SLOP-004는 생성 archive hash·설치 경로·실제 smoke 종료코드를 사용자가 확인하는 CLI 출력이며 debug 잔재가 아니다. 경고를 exclude/억제로 지우지 않았다. 미해결 critical/high를 면제하지 않는다. 독립 고정 SHA review·QA는 coor가 배정한다.

### Packet·이행

partial packet의153개 nonoptional path/category를 [outcomes](../../evaluations/jev/SAR-SSE-INSTALL-001-DEV-packet-outcomes.json)에 기록했고 identity·완료 검사가 통과했다. 오래된 QA/handovers·다른 역할 후보는 제목/범위/hash와 현재 원천·변경을 대조해 no_change 사유를 남겼다. 과거 모든 원문을 무의미하게 재독하지 않았다. 필수 live auto-receive QA, current source·common 규약·inbox는 직접 읽었다. oversized/sensitive 파일은 실제 소스에서 확인했고 producer mismatch를 의미 분석 성공으로 표시하지 않았다. D08 SQL은 migration 변화가 없어 재생성하지 않는다.

최종 제품 source SHA와 archive 후 exact SHA는 실제 worker_done에 보고한다. 완료 inbox는 work.py finish로 원문과 완료 보고를 archive한 뒤 비운다. 기존 production ba3754d는 보존한다. 후보 수락 이후 실제 host별 설치/Google 명시 동의/새 관계 승인/왕복/Events 활성화가 남는다. remote 개인키·webhook secret은 새 설치의 private DB custody이며 backup 접근 제어가 필요하다. webhook은3회 재시도·cursor null이므로 신호 유실이 가능하고 본문은180초 TTL 안에서 receive해야 한다.


최종 제품 코드 commit은 `1f7a23f595696e8b9dd84b9ddfe60714e251dbd0`다. replay 보완 이전 검증은 과거 후보 증거이며 변경 경계의 isolated HTTP 재실행·최종 gate로 대체한다. 설치물 runtime은7f36c7b 이후 변경되지 않았으며4개 Windows bundle SHA256은 evidence와 INSTALL 하위 hash를 구분해 아래에 기록한다.
- `knowslink-claude-desktop.mcpb`: `a645abb5554589296962fb77996adc59d0e24ab0833623dc115a7d163c284efd`
- `knowslink-claude-plugin.zip`: `2b60e09f4fb0f401ef9a5793d807a60608a8d63045fb4fcc9d0a6eacc8e86e7d`
- `knowslink-codex-plugin.zip`: `6925cfe68afcbad5a782023692b19f383106ec2d1a0048505772f4434396f432`
- `knowslink-grok-bot-plugin.zip`: `0d265b6704824e0cbe12b8cc6c399e71df31d678b41f09c836ea8942035800db`
