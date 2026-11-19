---
title: SAR-SSE-INSTALL-001 원격 OAuth·MCP·도착 이벤트 구현
status: draft
updated: 2026-10-11
owner: dev
tasks: [SAR-SSE-INSTALL-001-DEV]
summary: Google 계정의 설치별 원격 agent와 검증된 텍스트 Inbox 및 서명된 도착 webhook을 구현한다.
---

# SAR-SSE-INSTALL-001 원격 경로

## 범위와 기술 계획

기준은 `7a6906851b872fc52310a0ac6623527a4b0adb2b`다. 적용 규약은 `fullops-common-0.3.3`, project.md, 문서 작성 규칙, delegation.md, Ponytail full이다. 테스트 레벨은 lite다. 이 worker는 추가 하위 작업을 배정하지 않았다. 부모가 Go 연동·설치 worker·커밋·최종 검사를 통합한다.

기존 Google 로그인과 회원·owner·agent·pair·signed text 상태 기계를 재사용한다. 원격 OAuth 설치에는 새 agent와 서명 키를 만든다. 원격 bearer는 기존 `/v1` credential로 사용할 수 없다. 수신은 서명 검증 뒤 서버 Inbox 저장과 기존 persist/ACK를 같은 트랜잭션으로 처리한다. 이벤트 전송은 빈 PostgreSQL NOTIFY를 LISTEN하고 구독별 lease를 획득한다. idle 상태에서는 DB를 폴링하지 않는다.

소유 파일은 `internal/relay/remote_auth.go`, `remote_mcp.go`, `remote_events.go`, `remote_test.go`, `remote_integration_test.go`, Google 로그인에 필요한 `google.go` 연동과 이 보고서다. 부모가 합의된 State·route·Cleanup lifecycle·admission·notification·ingress 연동을 소유한다. 부모 inbox를 보존한다.

## 구현 계약

- `s.remoteRoutes(mux)`는 `/mcp`, OAuth metadata, DCR·authorize·token·revoke와 `/remote/consent/{id}`를 등록한다. HTTPS origin은 설정된 Google callback URL에서만 얻는다. 요청 Host는 origin 근거로 사용하지 않는다. Google 설정이 없으면 원격 경로가 실패한다.
- OAuth는 public DCR client의 authorization-code·S256 PKCE를 지원한다. client·정확한 redirect URI·resource·scope를 묶는다. 60초 code는 한 번만 사용한다. access token은 1시간이다. refresh token은 30일이며 토큰 쌍을 원자적으로 회전한다. 토큰은 DB에 hash로만 저장한다. issuer identification을 metadata와 승인·거절 callback에 포함한다.
- 이미 회전한 refresh token을 재사용하면 현재 token family·설치 key·구독도 철회한다. opaque refresh token의 family 식별자는 installation·client·resource·nonce에 domain-separated HMAC을 적용한다. 검증된 이전 token만 family를 철회할 수 있다. 무한 token-history tombstone을 저장하지 않는다. refusal을 transaction 값으로 반환하여 철회가 commit되도록 한다. access token은 계속 43자 random token을 사용한다.
- Google 로그인은 remote attempt의 계정을 먼저 묶는다. 다른 Google 계정은 이미 묶인 요청을 바꾸지 못한다. 5분 이내 신원 확인과 session-bound CSRF가 있어야 명시적 연결 승인이 가능하다. 승인은 기존 agent·키·pair를 덮어쓰지 않는다. 새 agent의 pair는 사용자가 별도로 초대·수락한다.
- 원격 전용 서명 개인 키와 webhook secret은 서버 비공개 DB에 보관한다. 동의 화면이 이 보관 위치를 명시한다. 기존 로컬 개인 키를 가져오지 않는다. 별도 배포 secret 없이 background webhook을 수행해야 하므로 새 애플리케이션 암호화 키를 만들지 않았다. DB와 백업의 접근 제어가 이 비밀값의 보안 경계다. key·agent·owner 철회 또는 설치 만료 시 private key·토큰·구독 secret을 삭제한다.
- `knowslink_status`, `knowslink_send`, `knowslink_receive`, `knowslink_ack`를 제공한다. write tool에 approval 판단용 annotation을 제공한다. send에는 사용자의 명시적 요청이 필요하다. OAuth 연결과 arrival event가 자동 전송 권한을 부여하지 않는다.
- send는 기존 UUIDv7·UTC 초 단위 timestamp·Ed25519·idempotency·pair·관련 답장 검증을 재사용한다. receive는 발신 서명을 검증하고 설치별 Inbox에 저장한 뒤 relay ACK를 처리한다. HTTP 응답이 유실되어도 다음 receive가 같은 Inbox를 반환한다. 명시적 knowslink_ack만 읽은 항목을 삭제한다. 기존 메시지의 180초 TTL과 관계 철회가 Inbox에도 적용된다.
- MCP 2.0 `server/discover`는 `2026-07-28`과 events capability를 광고한다. 기존 클라이언트용 initialize와 stateless Streamable HTTP JSON 응답도 제공한다. 별도 GET SSE session은 제공하지 않는다.
- `knowslink.text.arrived` 구독의 filter는 빈 object이며 인증된 설치 자신의 agent로 고정된다. deterministic ID는 installation·callback URL·event name·빈 arguments를 묶는다. 구독·refresh·unsubscribe는 같은 인증 MCP endpoint를 사용한다. finite TTL은 기본·최대 1시간이며 최소 1분이다. 요청 TTL이 null이어도 finite TTL을 부여한다. replay는 지원하지 않으며 cursor는 null이다.
- callback 활성화 전에 signed verification challenge를 전송한다. fresh challenge의 2xx 응답과 constant-time echo 비교가 필요하다. 같은 설치·URL·secret의 검증 결과는 5분 동안 재사용한다. secret 교체 시 1분 동안 기존·새 secret 양쪽으로 서명한다.
- 모든 callback 연결은 HTTPS·443만 허용한다. 매 연결의 모든 DNS 응답이 public인지 확인한 뒤 검증한 IP로 연결한다. TLS hostname 검증은 원래 hostname을 사용한다. proxy·redirect를 사용하지 않는다. private·loopback·link-local·shared·예약 주소와 IPv6 전환 주소를 차단한다. verification에도 같은 검사를 적용한다.
- webhook payload는 `arrival: true`만 포함한다. 원문·메일·agent ID·credential·개인 키를 보내지 않는다. Standard Webhooks HMAC-SHA256은 event ID·Unix timestamp·정확한 body bytes를 서명한다. 재시도는 같은 event ID를 유지하고 timestamp·signature를 새로 만든다.
- `RemoteEvents(ctx)`는 프로세스당 한 LISTEN 연결을 사용한다. 시작·재연결·구독 변경·queue/auth NOTIFY에서 현재 상태를 평가한다. retry·lease·구독 만료 시각만 timer로 기다린다. 원자적 lease가 여러 프로세스의 같은 delivery 동시 claim을 막는다. 각 발송 직전에 현재 설치·key·owner·구독·pending을 다시 검증한다. context 취소가 LISTEN·HTTP·reconnect timer를 끝낸다.
- 전송은 최대 3회다. 재시도 지연은 1초·2초다. 프로세스 crash 뒤 lease 회수도 같은 시도 상한을 사용한다. 410·413에는 재시도하지 않고 구독을 종료한다. 성공한 동일 pending 집합은 재전송하지 않는다. 네트워크 응답 유실 시 동일 ID의 중복은 가능하다.
- 상태 상한은 client 1000개, authorize attempt 1000개, 전체 install 500개, 전체 subscription 128개, install별 subscription 4개, install별 Inbox 32개다. 기존 회원·agent·HTTP·message 한도를 추가로 유지한다. 키·본문·secret을 HTTP 오류나 로그에 넣지 않는다.

## 공식 문서 근거

Context7 `resolve-library-id`를 pgx v5 LISTEN API로 호출했다. 결과는 `Monthly quota exceeded`였다. 기존 pgx v5.10.0 호출 패턴과 아래 공식 문서를 직접 열어 대체했다. 새 dependency는 없다.

| 실제로 연 문서 | 확인한 절과 적용 |
|---|---|
| https://developers.openai.com/plugins/build/auth | Custom auth with OAuth 2.1, MCP authorization spec requirements, Protect callbacks with issuer identification, Redirect URL, Echo the resource parameter, Support the authorization-code flow. discovery·DCR·resource·PKCE·issuer binding에 적용했다. |
| https://developers.openai.com/plugins/build/mcp-events | Before you start, Advertise event support, Create a subscription, Verify the callback, Send an event, Sign the request, Handle delivery responses, Manage subscriptions. MCP2·webhook·challenge·서명·TTL·rotation·retry 계약에 적용했다. |
| https://developers.openai.com/plugins/build/mcp-server | Create the server, Define tools from user goals, Return useful results without UI, Security reminders. tool schema·annotation·structuredContent와 승인 경계에 적용했다. |
| https://developers.openai.com/plugins/deploy/connect-chatgpt | Prepare the endpoint, Add the MCP server, Refresh metadata, Test the complete plugin. 공개 HTTPS 연결·도구 scan·인증·설치 검증은 후속 절차로 구분했다. |
| https://developers.openai.com/plugins/build/plugins | Create a plugin manually, Configure authentication, Use automatic OAuth client registration. 패키지 worker가 등록된 실제 OpenAI connection ID를 사용한다. 임의 ID나 실사용 credential을 package에 넣지 않는다. |
| https://www.rfc-editor.org/rfc/rfc9700.html#section-4.14.2 | Refresh Token Protection의 Recommendations. public client의 refresh replay 탐지와 active family 철회를 적용했다. MAC으로 family 정보를 검증하여 history 저장을 늘리지 않았다. |

공식 MCP Events 문서는 Work web·desktop Cloud·Dots 지원을 명시한다. 일반 relay SSE와 이 webhook 경로는 다른 프로토콜이다. 실제 Dots 계정 연결·설치·event-triggered task 실행은 아직 검증하지 않았다.

## 검증과 인계

Windows 작업 트리에서 다음 명령이 종료코드 0으로 통과했다. shell pipeline으로 테스트 종료코드를 숨기지 않았다.

- `go test ./internal/relay`: OAuth 설치 격리·code/refresh replay·client/resource/redirect/PKCE·current key revoke·Google binding·verified Inbox response-loss/restart·send idempotency·pair denial·Inbox TTL 검사를 포함한다.
- 같은 suite의 event 검사는 pending catch-up·동시 claim 거절·retry backoff·안정 event ID·terminal delivery status·crash retry bound·subscription identity/expiry·notification bookkeeping·SSRF IP·no redirect·정확한 Standard Webhooks 서명·callback challenge 실패를 확인한다.
- `go vet ./internal/relay`.
- `go test -tags integration ./internal/relay -run '^$'`: 통합 테스트 컴파일만 확인했다. DB 실행 통과로 표시하지 않는다.

`TestRemoteOAuthHTTP`는 isolated DB 전용 통합 테스트다. 기존 `TEST_SYNTHETIC_DATABASE=1`과 `TEST_DATABASE_URL` 보호 조건을 사용한다. 로컬 RSA/JWK Google provider로 실제 HTTP DCR·Google callback·consent·PKCE code exchange를 검증한다. 같은 Google 계정에서 설치별 agent가 다르고, 다른 Service 프로세스가 저장된 token을 인식하며, localhost callback이 실패하고 revoke가 접근을 끝내는지 확인한다. 부모가 Linux 격리 DB에서 이 테스트와 필수 검사를 실행한다.

초기 후보 `3c523fc`와 부모 문서 후보 `515eabb`의 Linux 격리 DB 검사에서 `TestRemoteOAuthHTTP`, `TestTextEventsHTTP`, `TestPublicMessageHTTPBoundaries`, `TestPublicNodeProcesses`가 race 모드로 통과했다고 부모가 보고했다. 이후 자체 보안 점검에서 refresh replay가 active family를 철회해야 하는 경계를 발견했다. 같은 과제에서 이 경계를 보강하고 unit 검사를 통과했다. HTTP 통합 테스트에도 새 Service 요청에서 철회가 유지되는 조건을 추가했다. 변경 후 고정 SHA와 재검사는 부모 완료 기록을 정본으로 사용한다.

부모가 커밋을 단독 소유한다. 이 worker는 git commit·branch·생산 DB·Tunnel·Google 설정·운영 key·유료 API·local Docker를 변경하지 않는다. 부모가 최종 SHA와 registered product-lint/product-test 결과를 부모 완료 보고에 남긴다. 새 동의 화면의 직접 시각 검수도 부모가 synthetic preview로 수행한다. 독립 fixed-SHA 리뷰와 tester 수락은 coordinator 후속이다.

부모는 synthetic consent를 390×844와 1280×800에서 직접 확인했다. 읽기·버튼 상태는 확인되었다. 긴 callback URI의 줄바꿈과 기존 화면 일관성을 위해 기존 `memberStyle`을 그대로 재사용했다. 새로운 CSS나 디자인 token은 추가하지 않았다. 스타일 적용 뒤 최종 확인은 부모가 수행한다.

## 한계와 남은 일

실제 public callback 성공·Dots 수동 OAuth consent·ChatGPT metadata scan·호스트 task wake는 미실행이다. 단위 테스트의 callback transport는 synthetic 응답을 사용한다. 실제 외부 webhook을 만들거나 보내지 않았다. 공개 배포 전 공식 client로 OAuth와 events 전체 lifecycle을 확인해야 한다.

이 이벤트는 replayable하지 않다. 3회 실패나 연결 중단으로 arrival 신호를 놓칠 수 있다. 본문은 기존 180초 TTL 안에서 명시적으로 receive해야 한다. 구독·이벤트는 사용자 대신 자동 답장하지 않는다. remote key의 서버 custody와 DB backup 보호는 운영 인계 항목이다.

변경 크기는 OAuth server·JSON-RPC tools·SSRF-safe signed webhook이라는 새 신뢰 경계에 따른다. 공유 기존 API의 권한 검사를 복제하지 않고 재사용했다. 새로운 범용 framework·SDK·cipher secret·배포 설정은 추가하지 않았다.
