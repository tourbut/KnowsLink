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

설치 소유권은 Orca child Run `run_5b1823734740`에 분리했다. 하위 작업은 scripts의 패키징/설치, adapter manifest·skill·README 및 공식 근거 보고서만 수정한다. 부모는 Go·adapter 수신 소스/테스트·기술/운영 정본·커밋을 소유한다. 실제 호스트 설치·Google 수동 동의·배포는 후속이며 새 유료 API·Docker는 사용하지 않는다.


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

UI는 SSE 경로에서 변경하지 않는다. 별도 디자인 lint/테마는 프로젝트에 없으므로 적용하지 않는다. 원격 OAuth 동의 화면은 기존 template/CSS와 보안 정책을 재사용하며 직접 시각 검수는 후속이다. 운영 Google 계정·DB·키·관계·Tunnel·호스트 전역 설정은 변경하지 않는다.

SIZE-002: SSE source와 핵심 회귀를 하나의 기능 commit으로 분리했다. 구현의 연결 자원·권한·취소·동시성 경계와 실제 HTTP/DB 테스트가 필요해300줄/400추가 경고를 넘는다. 광범위 스타일 정리나 새 추상화는 하지 않는다. 설치 및 remote는 별도 기능 commit으로 보존한다.
