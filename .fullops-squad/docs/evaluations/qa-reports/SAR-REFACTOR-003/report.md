---
title: 전반 리팩터링·코드 정비·불필요 파일 정리·로직 점검 결과
status: draft
updated: 2026-11-05
owner: coor
tasks: [SAR-REFACTOR-003]
summary: "점검 162건 중 확인 159건을 영역별로 반영·보류한 내역, 독립 리뷰, Windows·Linux 고정 SHA 검증 결과를 기록한다."
---

# SAR-REFACTOR-003 결과

기준 main `952f6eaa`에서 coor가 사용자 지시로 직접 수행했다. 실행 기록과 보류 사유는 [실행 계획](../../../exec-plans/phases/SAR-REFACTOR-003.md)에 있다.

- **점검:** 하위 에이전트 240개가 영역 8곳을 점검하고 지적마다 반박 검증을 거쳤다. 162건 중 확인 159건, 이견 1건(F003), 반박 2건(F033, F175)이다.
- **반영:** 확인 159건 가운데 148건을 반영하고 11건을 보류했다. 반영 항목 일부는 검증자가 선택이라고 한 부분을 뺐다. 보류는 사용자·기획 판단이 필요하거나 게시 파일을 바꾸는 항목, 효과가 작은 항목이다. 이견 1건(F003)도 보류했다.
- **변경 규모:** 120개 파일, 1,107줄 추가, 2,023줄 삭제다. 제품 코드만 보면 97개 파일, 1,030줄 추가, 1,958줄 삭제다.
- **운영:** 운영 배포, 실계정 시험, 게시 파일 교체는 하지 않았다.

## 검증

| SHA | 범위 | install | make lint | make test | build | plugin verify | migrate | race integration |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `952f6eaa` | 기준선 | 0 | 0(재실행, 첫 실행은 컨테이너 docker CLI 부재) | 0 | 0 | 0 | 0 | `abnormal_clock` 2건 간헐 실패, 3회 재실행 0 |
| `29ce3d47` | 영역 병합 | 0 | 0 | 0 | 0 | 0 | 0 | 0 (186 s) |
| `3140acde` | 리뷰 SHA | 0 | 0 | 0 | 0 | 0 | 0 | 0 (150 s) |
| `97195a76` | 리뷰 수정 | 0 | 0 | 0 | 0 | 0 | 0 | 0 (115 s) |
| `02dd8c54` | 최종 후보 | 0 | 0 | 0 | 0 | 0 | 0 | 0 (102 s) |

- **Linux 검사 환경:** 로컬 Docker 컨테이너(node 22.23, Go 1.27.1, gcc, docker CLI)와 격리 Postgres 17(compose와 같은 digest)이다. 운영 서버는 쓰지 않았다. 매 실행마다 해당 SHA의 bundle을 새로 clone했다.
- **Windows 검사:**
  - toolchain gofmt, `go vet`(기본과 fixture 태그 전체), Go 단위 시험이 통과했다.
  - `npm run check`, `npm test`(adapters), `package_plugin.py --verify`가 통과했다.
  - `claude plugin test adapters`의 hook 시험 13건이 통과했다. `src/*.test.ts` 13건의 "did not load"는 node 시험을 이 도구가 불러오지 못한다는 기존 표시다.
  - shell guard 시험과 Get-Registered 시험이 통과했다.
- **FullOps lint `--from 952f6eaa`:** ERROR 0이다.
  - 실행 불가 2는 등록 명령 `make lint`·`make test`다. Windows PATH에 make가 없고, public-config 시험은 Linux 파일 모드 600이 필요하다. 위 Linux 결과로 대체했다.
  - WARNING은 SIZE-001(시험 추가로 파일 증가), ANTI-005(캘린더 정리 시험·일회성 fixture 삭제, RC-01 시험은 `integration_test.go`로 이동), SLOP-004(기존 결과 출력 문구 변경), DEP-001(모션 npm script 인자만 수정, 의존성 변경 없음), SIZE-002(전반 정비라 400줄 초과)다.
- **독립 리뷰 `3140acde`:** 영역 5곳 리뷰어와 보안 리뷰어를 두고, 지적마다 반박 검증을 거쳤다.
  - 지적 3건 중 확인은 low 1건이다. 수정은 `97195a76`이며 red/green 확인을 거쳤다.
  - 변경분 `3140acde..97195a76`을 별도 리뷰어가 다시 보았다. 새 결함은 없고, 남은 low 1건(미연결 watcher 미정지)은 `e4ba2edc`에서 고쳤다.
  - critical·high·medium은 0건이다.
- **환경 정리:**
  - CRLF로 체크아웃돼 있던 2,346개 파일을 LF로 맞췄다. 내용 해시는 index blob과 같다.
  - PATH의 오래된 gofmt 대신 toolchain gofmt를 썼다.
  - 다른 작업이 Docker 이미지를 정리해 Linux 검사 이미지를 다시 빌드했다.

## 영역별 반영 항목

### go-core (영역 커밋 `bbca5802`)

| ID | 등급 | 종류 | 반영 내용 |
| --- | --- | --- | --- |
| F058 | medium | logic-bug | Trial (relay.test.message) claims no longer keep a stored ClaimToken, so they hold none of the 4 service-wide claim slots. The claim response still carries a claim value and Claimed stays true. Regression test: 5 more trial claims in a row all succeed and claimCapacity() is still true. |
| F025 | medium | logic-bug | boundedHTTP admits form bodies for POST /home/conversation/send up to 16384 bytes. Other form routes and /v1/connect/ stay at 8192. Integration test: the 1366x'한' case now also checks that the composer re-renders (action="/home/conversation/send"), and a new 1365x'한'+'x' case expects 303. interface-design.md:322 updated. |
| F060 | low | error-handling | Used verifier 2's form: when ParseForm fails in boundedHTTP, r.Form and r.PostForm are reset to nil so the handler parses again and returns its own refusal. capacity.go:287 is unchanged. Integration test: POST /home/conversation/send?x=%zz with a valid form expects 422. |
| F061 | low | logic-bug | New optional field Service.Stopping (context). When it is cancelled, textEvents cancels its stream: defer context.AfterFunc(s.Stopping, cancel)(). cmd/relay sets Stopping: ctx, the signal context. New integration test TestTextEventsEndOnStopping: the stream ends and st.Streams is empty. |
| F062 | low | consistency | respond() now writes errorBody(err), so the 429 agent_rate refusal on /v1/agents includes retry_at. Every other error body is byte-identical. |
| F063/F051/F135 | medium | unnecessary-file | Deleted (git rm) calendar_retire.go, calendar_retire_integration_test.go and the one-off publish_browser_fixture_test.go. Removed Gate.Lookup, State.Calendars, googleAttempt.Calendar, the RetireCalendars call in main.go, KNOWSLINK_CALENDAR_KEY in compose.yaml, and .env.example lines 19-21. Moved TestCalendarRoutesRetired (RC-01) into integration_test.go. Added a one-line SAR-REFACTOR-003 note to t |
| F064 | low | dead-code | Removed the duplicate_result branch in ingest, which could not be reached. The statusFor mapping stays because public_text.go still returns that code. |
| F065 | low | simplification | Removed leaseMessageFor. http.go now calls leaseMatching(principal, now, testOnly, false). leaseMessage stays for the legacy test, as verifier 1 refined. |
| F067 | low | stale-doc | Rewrote the header comments of sqlc.yaml, cmd/relay/main.go and config.go. Deleted the orphaned gate-page comment in http.go. |
| F070 | low | dead-code | ownerAuth now returns (token, back). gateOwner calls st.principal(token, "owner") directly. Updated the three destructurings. |
| F072 | low | simplification | Merged sweep's revoked and expired branches into one guard. The expired check still comes first, so st.current is not called on expired messages. |
| F074 | low | duplication | The four repeated pool setups in integration_test.go now call messagePool(t). TestPostgresSafety keeps databaseURL for its restart subtest. The optional extra files were not converted (other areas). |
| F075 | low | simplification | All 8 items, (a) to (h): gate Digest set in the literal, object(sig) check, package-level sha256Pattern, redundant credential clause removed, capacity codes moved into the 409 case, strings.CutPrefix, trailing continue removed, unused m and nil checks removed from gateDecision. |
| F161 | medium | logic-bug | Added _ "time/tzdata" to cmd/relay so LoadLocation works in the alpine image. |
| F162 | low | error-handling | Removed the ErrNoMigrations no-op success branch in cmd/migrate. Dropped ErrNoMigrations from the API list in tech-stack.md, which the finding marked optional. |
| F163 | low | dead-code | Removed seedLegacy, the RELAY_LEGACY_SEED branch, and the crypto/rand and base64 imports. The header now gives the fixture's origin (a6a10c7, seed in 3c8dc733). |
| F164 | low | dead-code | The invites guard is now !st.SyntheticSignup. policy_test sets SyntheticSignup = true after connectionFixture() and again after the JSON reload. |
| F165 | low | duplication | transaction() loads state through one load closure, used on both the normal path and the rejected-operation path. |
| F166 | low | unnecessary-file | Moved rateLimited into capacity.go and deleted api_rate.go with git rm. Updated the three module-design.md rows to name capacity.go. |
| F167 | low | stale-doc | Removed the /ingest, /approve and /exec rows from main_test.go, renamed the test to TestReadiness, and fixed its header comment. |
| F044/F066 | low | duplication | Removed the duplicate Cache-Control, Referrer-Policy and X-Content-Type-Options sets in Handler(). Kept the CSP override, with a one-line comment. |
| F143 | low | unnecessary-file | Deleted db/queries/.gitkeep with git rm and kept db/migrations/.gitkeep. Replaced the sqlc.yaml header with an accurate line. |

건너뜀:

- F026/F059: The task says not to do this here (dead request-budget/rate removal is a later step). requestHit, memberHit and requestBudgetKey are unchanged.

### go-identity-member (영역 커밋 `9927c4d6`)

| ID | 등급 | 종류 | 반영 내용 |
| --- | --- | --- | --- |
| F018 | low | error-handling | logout, authReauth (both transactions) call expired() only on invalid_auth, otherwise render the existing 503 refusal; reauth's first-transaction fault returns before the email path; startPage clears the cookie/sets expired notice only on invalid_auth. |
| F019 | low | security | device-confirm invalid-name branch now requires c.Owner == owner (else invalid_auth refusal) and adds Client to the 422 view. |
| F155 | low | logic-bug | Same branch: owner check plus "Client": clientName(c.Host) and "Here": "/connect/"+id so the 422 page shows the client line and lang-nav. Integration regression: owner posts name=+++ -> 422 with client dd and /lang/en?back=%2fconnect%2f<id>; foreign owner -> 401 without fingerprint. |
| F030 | low | logic-bug | Dropped unread "State" key from the device-confirm 422 view and from devicePage's view (Client/Here added as above). |
| F020 | low | logic-bug | maskEmail keeps the whole first UTF-8 rune (utf8.DecodeRuneInString); unit assert for 한글@example.com. |
| F021 | low | error-handling | sendOwnerText rejects bodies containing NUL as invalid_schema; agent-create treats an attempt containing NUL like a bad-length attempt (no idempotency). Unit regression: "a\x00b" added to invalid-body list. |
| F022 | low | dead-code | Removed unreachable invalid_auth -> expired() branch in memberAction's error path; only the 503 refusedTo remains. |
| F023 | low | dead-code | Removed connectionPage unsure block/Unsure key (State passed directly), the {{if .Unsure}} template branch, deviceUnsure var and agentResume's deviceUnsure return; removed hand-built tests (connect_guide_test unsure render, i18n_test deviceUnsure.State, integration device_lost part with subtest renamed to remote_uncertain_results_stay_check_needed, connectfixture device_lost seed); removed the D-C |
| F087 | low | dead-code | Same deviceUnsure cleanup as F023 (merged). |
| F159 | low | dead-code | Removed agentResume's unreachable `if last.Device { page = "/connect/"+lastID }`. |
| F024 | low | simplification | finishGoogleLogin resolves the owner once from Identities after a successful signIn and binds Connection/Remote without the two unreachable session lookups. |
| F027 | medium | logic-bug | Removed the legacy-home '활성 {{.ActiveAgents}}/5' span, the active counter loop and ActiveAgents view key; dropped ActiveAgents from member_ui_test fixtures. |
| F078 | medium | consistency | Same counter removal as F027 plus the unused legacy agents slice/loop/sort and "Agents" key. |
| F158 | low | dead-code | Removed unread view keys: legacy-home Agents (+loop/sort), connect view Owner, device State (both maps). |
| F028 | low | dead-code | Removed unreachable invites redirect (`before`, inviteNotice), the invited/invite-pending/invite-active notices, memberPair.Own/Other and their side loop, connect view Owner key; updated member_ui_test pair helper and home-empty notice (now notices["logout"]), connections_test own field; marked interface-design.md invite notice/re-invite lines as superseded by HP migration; removed inviteNotice/me |
| F148 | low | dead-code | Deleted the test-only {{define "home"}} alias; tests now render "legacy-home" (connections_test x5, member_ui_test fixtures and invite test, public_text_test). Device State key removal as in F158. |
| F096 | low | dead-code | Gate fixtures use tmpl "human-gate"; renderFixture reduced to renderHumanPage(t, f.tmpl, f.data) (bytes import dropped); removed the gate mapping in i18n_test. Golden unchanged (keyed by fixture name). |
| F029 | low | logic-bug | Read-only conversation header: history sorted before the canSend branch, names taken from the newest message, Peer set only when OwnerA != OwnerB. New unit test TestHumanReadOnlyConversationHeader (self conversation, renamed agent, revoked target; 20 renders). |
| F156 | low | logic-bug | Self read-only conversation no longer sets Peer (covered by F029 change and test). |
| F081 | low | logic-bug | Self-conversation list label now tr(lang, "내 에이전트") (EN key already exists). |
| F031 | low | simplification | humanClock uses t.Format("3:04 PM") / "3:04" instead of hand-rolled 12-hour math; '오전 '/'오후 ' literals unchanged. |
| F034 | low | simplification | humanPage error path converts a non-auth error into the refusal value handled by the existing switch; conversation/send error view sets FormTemplate "human-conversation" and the special-case render branch is removed. |
| F036 | low | dead-code | Removed dead `back = "/home"` in devicePage closure; migrateHuman drops unused now parameter (4 callers updated incl. store.go); googleFailure is one line; connectAPI uses errorBody(parseErr). |
| F144 | medium | logic-bug | sweepIdentity uses a per-key cutoff: http:* and cleanup* keys pruned after 1 minute, everything else (send buckets, future buckets) after 1 hour (verifier 2's fail-safe allowlist). Unit assert: after +2m http:ip:ip1 and http:new are gone, send:email:<hash> stays. |
| F145 | low | error-handling | verifyPage renders the existing 503 start view on a transaction error instead of the 410 'code expired' page; pending cookie kept. |
| F146 | low | error-handling | deliverCode logs the constant mailer error: log.Printf("verification mail not sent: %v", err). |

건너뜀:

- F147: Needs a product/design decision: humanAgentView also feeds GET /v1/contacts/{agent} "client" and MCP sender.client, so mapping supportedClient to 'Node 22 로컬 클라이언트' changes an agent-facing API value (and puts Hangul in EN agent output) besides the visible label; the finding itself asks for a design 

### go-pages-i18n (영역 커밋 `c3d71e5d`)

| ID | 등급 | 종류 | 반영 내용 |
| --- | --- | --- | --- |
| F077 | medium | consistency | human-retention notice now says backups keep messages up to 14 days. The Korean template run and the matching enHuman key were changed in one edit, and the EN value now says 'for up to 14 days'. The optional SAR-HUMAN-PEER/SAR-PUBLIC-SERVICE doc notes were not made. This is user-visible copy, so the coordinator or user signs off on the wording. |
| F079 | medium | error-handling | startPage now clears the session cookie and shows the expired notice only when the error is invalid_auth; any other DB error renders the anonymous landing page with status 200 and keeps the cookie. authReauth now calls expired() only on invalid_auth and otherwise answers 503 through refused("내 에이전트", ...). Added the regression test TestSessionLookupFailureKeepsCookie, which uses an unreachable pgx |
| F080 | medium | logic-bug | Deleted enHuman's "요청 대기": "Awaiting reply". The English 429 concurrency refusal now uses enMember's 'Request pending'. The '요청 대기 /' key is kept. |
| F082 | low | logic-bug | memberLang now deletes langUnsavedHeader when a member's saved locale decides the language. Added a regression check to TestLanguageChoiceOrderAndLink: with /?lang=en and saved locale 'ko', the page is ko and shows no notice. |
| F083 | low | logic-bug | The policy template now has the {{if .Notice}} note after the 시행일 line. Added TestPolicyPageUnsavedLanguageNotice: /privacy?lang=en with no cookie shows the notice, and with the en cookie it does not. |
| F084 | low | consistency | The refusal back-link chain now labels /home/account#language as 설정/Settings. Added a ko/en assertion in TestLanguagePartsAndMarkers. |
| F085 | low | dead-code | Deleted dead EN entries that no tr() or template lookup can reach. From enMember: the 6 'plain part' fragments and their comment, the 2 agent_rate fragments, 승인하면, and the status words 연결 대기, 대기, 연결됨, 승인, 거절, 만료, 전달됨, 처리 완료, 실패, 해제, 차단 and 확인할 수 없음. From enHuman: 을 다시 실행하세요., 을 실행하세요., 대기/전달됨/만료/실패, and the duplicate 연결 필요/등록됨/철회. Kept 연결 완료, 활성, 미연결, 철회, 등록됨, 연결 필요 and 없음 in enMember. Updated the |
| F086 | low | duplication | Deleted every enMember entry that is also defined in enHuman (enHuman already won the merge, so the en map output is unchanged). This includes the 5 copies whose text differs and never showed. Added a guard at the top of TestTranslationCoversTemplates that fails when a key is in both tables. |
| F092 | low | duplication | Same change as F086 (one combined change). All 49 shared keys besides 요청 대기 are gone from enMember. 연결 필요, 등록됨 and 철회 now live only in enMember: their identical enHuman copies were deleted instead (see F085). The guard test is in place. |
| F035 | low | dead-code | Deleted the dead enMember keys '지문 확인 후 연결 승인' and '승인했습니다. 내 컴퓨터에서 저장 결과를…', and the enHuman fixture key '합성 다음 행동' with its comment. The human_pages_test fixture Next now uses the real string '내 컴퓨터에 표시된 지문과 같은지 확인하세요. 다르면 취소하세요.'. |
| F088 | low | dead-code | The friends list now renders only the active link, dropping the unreachable non-active and blocked branch. human-new no longer has the always-true active check. Deleted the now-unused key '차단 해제: {{name .Name}} / {{.ID}}'. |
| F094 | low | test-gap | In TestTranslationCoversGoText, the substring fallback now accepts only pattern keys (keys containing %). The test passes without any new EN entries. |
| F095 | low | dead-code | Deleted the never-produced badge states completed, registered and unconnected. Moved the orphaned badge doc comment above func badge, and removed 'status badges' from the member_pages.go header comment. The optional queued/leased removal was not done. |
| F097 | low | dead-code | Deleted the stale KNOWSLINK_UI_ARTIFACT_DIR block in TestMemberReceiptView, along with the os and path/filepath imports it used. |
| F098 | low | dead-code | Removed only selectors that match nothing anywhere (checked templates, Go strings and JS): .who, .inline, .sechead h1, button:disabled{…}, .note p:last-child/.problem p:last-child{…}, and .people .agents. |
| F100 | low | simplification | enPatterns is now built once at package init (strings.IndexByte) in place of the sync.Once wrapper, and tr reads it directly. The en map is built with maps.Copy(enMember) then maps.Copy(enHuman). The sync import is dropped and maps is added. |
| F154 | low | logic-bug | human-alerts now renders {{langPart .Problem}}. TestLanguagePartsAndMarkers now renders the real human-account page through humanPageFixture. Updated the module-design.md:106 wording to name human-alerts. |
| F157 | low | dead-code | The intro h1 is now {{.Title}}; .Heading has had no setter since 4795d619. |

건너뜀:

- F098 (partial): The assignment says to remove only selectors with zero matching class, id or element use. These remaining F098 items match real elements, so they stay: .ok,.warn,.off{color:var(--muted)} (badge() emits these classes), the overridden declarations (.conversation-header{margin-bottom:24px}, .retention{
- F095 (optional part): The optional removal of the queued/leased badge entries and the fixture changes from queued to pending were not done. The minimal change was enough.
- F077 (optional docs / friend-code dedup): The verifier said not to do the friend-code template dedup. The SAR-HUMAN-PEER/SAR-PUBLIC-SERVICE D021 notes are optional historical spec docs that this area does not own.

### go-remote (영역 커밋 `7dd69675`)

| ID | 등급 | 종류 | 반영 내용 |
| --- | --- | --- | --- |
| F002 | medium | error-handling | remoteMCP: right after the auth transaction, a 503-class error (unavailable/abnormal_clock) now goes to respond(w,nil,err), so the reply is 503 {"error":"unavailable"} with no WWW-Authenticate challenge. Real invalid_auth and o=="" still get the 401 challenge. Regression unit test TestRemoteUnavailableIsNotAuthRefusal uses an unreachable pgxpool and calls remoteMCP and remoteToken directly. |
| F006 | medium | error-handling | remoteToken: before the invalid_grant branch, a 503-class error now returns 503 temporarily_unavailable (same body as register/revoke). Client refusals and remoteGrantRefusal still return 400 invalid_grant. Covered by the same unit test as F002. |
| F041 | low | logic-bug | remoteMCP now reads the body with io.ReadAll(http.MaxBytesReader(w, r.Body, rawEnvelopeLimit)) plus Strict, the same pattern as /v1/send, instead of decodeCommand's 8 KiB cap; the 400 invalid_request response is unchanged. testSend's 64*1024 cap is now rawEnvelopeLimit. interface-design.md:227 now says 32 KiB (rawEnvelopeLimit 32768 bytes), as Verifier 2 named. New integration case in TestRemoteOA |
| F004 | low | consistency | Same fix as F041 (the /mcp body cap is now rawEnvelopeLimit). I used the explicit MaxBytesReader read rather than the httpBody variant, so the handler does not depend on running behind boundedHTTP. |
| F005 | low | dead-code | Removed the unreachable `!now.Before(i.Exp)` check from the subscribe lease transaction. remotePrincipal's install result is now discarded with `_`, with a one-line note that remoteInstall already refuses expired installs. |
| F007 | low | dead-code | Together with F005, moved the per-install count loop and the 4/128 cap check inside the existing `if sub == nil` block, so the count only runs when a subscription is created. |
| F008 | low | simplification | Deleted the hand-written GET /mcp 405 handler, so ServeMux now answers every other method with 405 and Allow: POST. Reworded the stateless Streamable HTTP comment above the POST /mcp registration. |
| F009 | low | simplification | remoteAuthorize: the redirect allowlist check is now slices.Contains, folded into the client nil/expiry check. remoteRegister: two slices.ContainsFunc checks, with the redirect check still first, then grants/responses; error codes unchanged. Added the slices import. |
| F010 | low | test-gap | Test-only, in TestRemoteOAuthHTTP (integration tag). Four authorize refusals: unregistered_redirect with the real client_id, unknown_client, plain_pkce and wrong_resource. Each copies q and expects 400 with no Location. Three register refusals: non-https redirect, implicit grant, token response_type. These probes run from their own client IP (198.51.100.7), because otherwise they used up the 30/mi |
| F011 | low | duplication | Deleted remoteAdmission. In requestBuckets (capacity.go), the /mcp case now uses st.remotePrincipal(bearer(r), s.remoteOrigin()+"/mcp", now) and keeps the 'never by the local credential API' comment. TestRemoteAdmissionAndEventVersions now asserts on remotePrincipal for the same three cases. Intended behavior change: a token with a stale resource spends the IP's anonymous bucket. |
| F012 | low | duplication | Added Service.listen(ctx) (*pgx.Conn, func(), error) to text_events.go. It sets up under a 3s ctx, LISTENs on knowslink_text and returns a release func that closes the connection with a 1s background ctx. textEvents and remoteListen now call it with a deferred release. Dropped the pgx import from remote_events.go. |
| F013 | low | simplification | (a) TTL clamp is now ttl = max(time.Duration(min(*p.TTL, 3600000))*time.Millisecond, time.Minute). (b) RemoteEvents waits with a select on ctx.Done and time.After(5s). (c) work.Delivery = &remoteDelivery{id, *sub}, which removes the shadowing of the builtin copy. (d) Not applied, per the verifier: the min cap stays. |
| F014 | low | duplication | remoteConsent: reads id and token once before the transaction and builds one view map {Name, Redirect, ID, CSRF}. GET adds Login and the default AgentName. The invalid-name and agentRefused paths set Problem and AgentName on that map. Login is still set only on GET, so the refusal maps keep the same keys as before. |
| F015 | low | stale-doc | remoteCredentials no longer sets exp on remote honored entries, so the honored doc comment holds again. |
| F016 | low | consistency | Added const publicTextTTL = 180*time.Second in public_text.go and used it in ingestPublicText, knowslink_send and owner text (human_messages.go). Added const remoteAccessLifetime = time.Hour for the access token expiry and expires_in. The Attempts>=3 sites, the deliverRemote check, store.go's 300s cap and the i18n copy are unchanged. |
| F119 | low | logic-bug | Verifier 2 version: finishRemote(job, status, now) dropped its failed parameter and now judges delivery by HTTP status only. deliverRemote ignores remotePost's error. Removed the false argument at the 5 test call sites. No new test: the change is in the DB-backed deliverRemote, and the existing finishRemote tests still cover the status rules. |
| F120 | low | error-handling | approveRemote at the 500-install cap now returns agentRefused{Code: "service_capacity"}, so the consent form re-renders with status 422 and the reason, and the attempt is kept. Regression unit test: TestRemoteInstallCapIsCapacityRefusal. |
| F121 | low | simplification | deliverRemote no longer re-checks the old-secret window with time.Now(). It passes sub.OldSecret, which sweepRemote already cleared on DB time. |
| F122 | low | simplification | remoteTools: dropped the read and idempotent fields. Annotations now use readOnlyHint: s.name == "knowslink_status" and idempotentHint: true, which gives the same output. |
| F123 | low | simplification | Removed the remoteNotify wake after a successful events/unsubscribe. Subscribe still sends it. |

건너뜀:

- F013(d): Verifier 1 said to skip it: the min(sub.Attempts-1, 3) cap stays.

### ts-adapters (영역 커밋 `7bbc680a`)

| ID | 등급 | 종류 | 반영 내용 |
| --- | --- | --- | --- |
| F037 | high | logic-bug | TextTransport.sendText now signs exp = now + 170000 (was 180000), staying under the relay's strict 180 s cap when the client clock runs slightly fast. Replies still clamp to the parent's exp. No Go change. |
| F101 | high | logic-bug | Same one-line change as F037. Also changed the README wording '독립180초 기한' to '독립된 최대 180초 기한', the optional tweak the verifier suggested. |
| F039 | medium | logic-bug | fetchJson keeps the raw body. On !response.ok it tries JSON.parse inside a try/catch and always throws 'relay HTTP <status>[ code][ retry_at]'. Successful responses are still parsed strictly. The text.js CLI timeout gap is fixed through F106 (safeError, which maps TimeoutError to 'relay unreachable'). Regression test: inbox.test step 4, where a 502 text/html pull gives error 'relay HTTP 502' and t |
| F102 | medium | error-handling | Same core.ts fix as F039 (duplicate finding). Not covered: a connection reset in the middle of the body (TypeError: terminated), as the verifier noted. |
| F103 | medium | error-handling | AutoReceiver.tick: `fresh` is declared before the try, and the notify block runs after try/catch, so a locally stored message is announced even when persist or ACK fails. Dedupe on re-lease is unchanged. The Inbox.store comment is updated. Regression test inbox.test step 4b: the ACK fails once, one notice still arrives, the state is backoff and lastReceivedId is set. |
| F109 | medium | error-handling | approvedAttempt parses inside a single try/catch, so a corrupt or empty pending.json returns null (a dead leftover, R-01) instead of throwing. Regression test in setup-connect.test: with an empty pending.json, setupConnect resolves 'connected', the first event is 'waiting', and the folder is kept aside as broken-private.unused-*. The optional mcp.ts .catch was skipped. |
| F040 | low | logic-bug | transient() now also matches 'HTTP 429 capacity', the relay's 1 s admission refusal; the comment names it. This follows the narrow token from F104 and F040 verifier 2, not the broader `HTTP 429(?! agent_rate)`, so a code-less edge 429 stays terminal. Test: connect-simple.test adds a one-shot injectOnce map, and a single 429 capacity on /v1/connect/complete still ends 'connected'. |
| F104 | low | error-handling | Same regex change as F040 (`HTTP 429 capacity`). |
| F042 | low | consistency | validText uses /[^\p{White_Space}]/u (the same character set as Go strings.TrimSpace) instead of text.trim(). Test in owner-text.test: '\u0085' is rejected as an invalid envelope and validText('﻿') is true. |
| F105 | low | error-handling | transient() moved word for word from setup-connect.ts into login.ts and exported, with F040's capacity token. setup-connect imports it. waitForLogin(folder, expires) treats transient errors as waiting until expires + 60 s, and main passes login.expires. The stderr text is unchanged. |
| F106 | low | consistency | The text.js CLI catch now uses safeError(e) from inbox.js instead of its own copied whitelist. Observable CLI changes: a timeout prints the outage text, the generic fallback word is 'failed', and invalid text envelope/signature are named. mcp.ts publicFailure is left alone, as verifier 1 advised. |
| F107 | low | concurrency | startAuto resolves memberTransport and Inbox.open first, then returns if `auto` is already set, before it creates the AutoReceiver. This is verifier 2's approach: no new flag, and a failed attempt can still be retried. |
| F110 | low | logic-bug | Added process.stdin.setEncoding('utf8') before the send stdin loop, so a UTF-8 sequence split across chunks is no longer corrupted. No test (it would need a child process). |
| F112 | low | duplication | localAdapter, setupTrial and testTransport now use relayBase() plus an `http:` check for loopback, and test-remote compares the origin exactly. The dead `url.port` clause is gone. Added a one-line comment on relayBase about this coupling. The existing trial-boundaries and mcp URL tests pass. |
| F113 | low | duplication | Added a protected Adapter.signedBy(from, sig, bytes) for the key fetch, SPKI wrap and verify. verifyMessage and TextTransport.receive call it; the error strings are unchanged. text.ts no longer imports createPublicKey or verify. The optional shared abort-race helper was not done. |
| F114 | low | simplification | The setup-connect poll loop now calls only finishLogin, which already handles the approved-attempt branch and scrubs login.json. The direct complete() call and the ./connect.js import are removed. |
| F115 | low | simplification | Split out status(state) from result(). knowslink_status spreads status(...) instead of re-parsing result's JSON text. Output is unchanged. |
| F116 | low | dead-code | Removed the CodexSession `bound` getter. visible is no longer exported, CodexSession.stop is private, text.ts watch is no longer exported, and reasons.ts failure is no longer exported. The connect.ts relayBase re-export is deleted and connect.test.ts imports relayBase from core.js. The optional PrivateFolderError.path removal was skipped (churn with no gain). |
| F152 | low | dead-code | Covered by the F116 edits: the connect.ts re-export is removed and failure is no longer exported. The optional index.ts:4 and synthetic.ts change was skipped, as verifier 1 advised. |
| F117 | low | consistency | core.ts exports AGENT and MESSAGE_ID (no g/y flags) with a comment pointing to protocol.go. They replace the regex literals in text.ts, connect.ts, login.ts and test-transport.ts, and inbox.ts's private ID. |
| F149 | low | concurrency | register.ts keeps the in-flight attempt in `starting`. A repeat `on` returns that promise and spawns nothing. `starting` is cleared (only if it is still the same attempt) on bound, on busy, in the reader's finally, and when config or session-id reading fails. The guard is set before the readConfig/session.id awaits, because two overlapping `on` calls otherwise both pass the check (the new test cau |
| F151 | low | logic-bug | GrokWake.request treats ECONNREFUSED, ENOTFOUND, EAI_AGAIN, ENETUNREACH and EHOSTUNREACH as retry. grok-wake.test loops the unreachable stub over ENOTFOUND and EAI_AGAIN. |
| F153 | low | duplication | One stat-based `exists` is exported from private-files.ts (the verifier chose stat over access). The three local copies in inbox.ts, mcp.ts and setup-connect.ts are deleted. |
| F118 | low | stale-doc | README Grok doctor row now says '9 tools discovered'. The mcp.ts header comment describes trial/synthetic/member modes and auto-receive unless KNOWSLINK_AUTO_RECEIVE=off. The text.ts header says 'manual or always-on (watch) receive'. |
| F049 | low | stale-doc | adapters/README.md: line 16 now says the screen shows an auto-setup command for Claude Code and Codex and lists Hermes Agent and others as needing a support check, keeping the Codex section link. Deleted the stale 2026-10-13 SHA sentence (:22) and '운영 게시는 ops 과제다.' (:60). Line 243 says 9 tools. No headings or anchors changed. |
| F150 | low | stale-doc | SKILL.md gains a '## Codex work session' paragraph (denied until owner `knowslink_session on`, take, off, status, YOLO, no arrival display). The autoReceive sentence now says it describes only the MCP process's own receiver: `off` in Claude Code, `stopped` in Codex, and which status command to use instead. |

건너뜀:

- F149 (doc part): Verifier 2's optional test-count update in module-design.md:66 was not made. The file is .fullops-squad/docs/design-docs/module-design.md, and .fullops-squad/** is off-limits unless a finding names the file; this one only names module-design.md:66 as optional. The coordinator may update it: register

### scripts-deploy (영역 커밋 `30add9f6`)

| ID | 등급 | 종류 | 반영 내용 |
| --- | --- | --- | --- |
| F124 | medium | logic-bug | verify_of now remembers that a wrapper prefix was seen and keeps scanning past option values until it finds a delete/find/git/robocopy verb, falling back to the first candidate. Added 3 DENY cases (timeout -s KILL 60, sudo -u root, env -u HOME rm -rf "$TMP_DIR/") and ALLOW 'timeout 60 go test ./...'. env -S is still not covered. |
| F141 | low | simplification | recursive() uses verb != 'unlink' instead of a copied verb tuple. targets() returns plain values and check() iterates them directly. |
| F125 | medium | logic-bug | regression() drops the knowslink project from both the saved baseline and the live state before comparing. |
| F128 | low | logic-bug | public() now requires 404 with no 'relay' in the body for every probe (D12 http_status:404 fallback). Removed the 302/cloudflareaccess branch, the urllib.parse import and the word 'Access' from the docstring. Chosen over F133 because F168 deletes the Access rollback template. |
| F133 | medium | logic-bug | Superseded by F128, which makes the same change in a stricter form: public now passes on the current 404 edge. |
| F136 | low | consistency | Added need(), copied from access_apply. Every assert in verify.py is now an explicit gate that raises SystemExit, with string messages, so the checks still run under -O. |
| F126 | medium | error-handling | beta.sh deploy: a failed 'git fetch origin' now only warns on stderr; cat-file still rejects an unknown SHA. D12 11.3 and contexts/ops.md now give the bundle-fetch step first. |
| F129 | low | error-handling | restore_verify cleanup trap changed from RETURN to EXIT, with '\|\| true', so the restore container is removed when set -e aborts the function. |
| F134 | medium | logic-bug | test_public_config.sh finds the repo root from its own path instead of git rev-parse, and exports KNOWSLINK_SRC so beta.sh never calls git. |
| F137 | low | consistency | Makefile lint-config replaces the hand-kept py_compile list with 'python3 -m compileall -q -f scripts deploy'. The optional selftest targets were not added because both selftest scripts are deleted (F168/F170). |
| F130 | low | test-gap | lint-go adds a second go vet pass with all fixture tags: integration, approvalfixture, autosetupfixture, browser, connectfixture, connectsimpleqa (missing from the verifier's list, added), landingfixture, langfixture, namesfixture, publishfixture, receiptfixture. The py_compile part is covered by F137's compileall. |
| F140 | low | consistency | package_plugin.py reads VERSION once from adapters/.claude-plugin/plugin.json and uses it for the bundled package.json, the cloud manifest and release(). verify_grok_plugin.py reads the version from the source adapters/.grok-plugin/plugin.json for its assert and its message. |
| F142 | low | simplification | Applied (a), (b) and (c); skipped (d). (a) FILES removed; grok, claude and claude-desktop maps are built from lists (order checked identical to HEAD). (b) other-private is created once before the verify loop. (c) PRODUCTION_RELAY moved next to ROOT and used for RELAY_URL. |
| F127 | low | error-handling | release() builds both ZIPs, fills the SHAs and runs the placeholder check before target.mkdir(exist_ok=False). Only then are the ZIPs copied and the installer written. A failed release no longer leaves an empty version folder (checked in scratch). |
| F172 | low | error-handling | Licence and JSON reads use encoding='utf-8'. The release installer is written with newline='\n'. |
| F168 | medium | unnecessary-file | Deleted access_apply.py and tunnel/config.yml.tmpl with git rm. In beta.sh removed HOST, TEAM, render_config, expose and their case/usage entries; kept tunnel_create. Removed the 'older Access-protected deployment' comment from public-ingress.yml. D12 section 11 gets a removal note, the 11.1 file table is updated, the access_apply/expose sentences are dropped from 11.3, and step 6 of the no-cost T |
| F170 | low | dead-code | Removed beta.sh seed and owner_login with their case and usage entries, and deleted access_trial_plan.py with git rm (together with the template from F168). Kept run_trial.py and the trial MCP tools (Grok tool count stays at 9). Dropped the module-design.md bullet, replaced the ops-guide 11.5 owner-login paragraph, and added a 'removed, historical' line to the user-guide trial procedure. ops-guide |
| F169 | medium | logic-bug | schema.py reads and writes UTF-8, writes LF, reads the target once, and hashes the LF-normalized text. Regenerated D08: only line 19 changed (43e0ce1e... to 31aa978f..., which equals the sha256 of git show HEAD:db/migrations/00001_relay.sql). Ran it without PYTHONUTF8. The front matter 'updated:' date was not bumped. |
| F173 | low | test-gap | check_compose.py also renders compose.yaml with deploy/knowslink/compose.ops.yaml using dummy STATE_DIR/UID/GID. It asserts the same services, Postgres without ports, read_only on migrate/relay/cloudflared, a read-only /install mount on relay, and no TUNNEL_TOKEN on cloudflared (a probe showed !reset clears it even when TUNNEL_TOKEN is set). |
| F174 | low | test-gap | test_get_registered.py passes -ExecutionPolicy Bypass before -File, which affects the process scope only. |
| F045 | low | test-gap | Added /v1/text/pull, /v1/text/persist, /v1/text/ack, /v1/connect/info, /v1/connect/prepare and /install/0.1.0/knowslink-codex-plugin.zip to the allowed paths in check_public_ingress.py. |
| F171 | low | test-gap | Added /auth/start, /auth/verify, /auth/logout and /auth/logout-all to the allowed tuple; the other F171 paths came in with F045. The regex is unchanged. |

건너뜀:

- F142(d): The verifier said to skip it: it only moves config_files to module level.
- F137-optional-selftests: access_apply.py and access_trial_plan.py are deleted by F168/F170, so there are no selftests to wire into make test.

### docs-files (영역 커밋 `e3087698`)

| ID | 등급 | 종류 | 반영 내용 |
| --- | --- | --- | --- |
| F038 | medium | stale-doc | README Grok pins changed. The new-install checkout (2747fe98) and the update 수락 SHA (8a93e38e) now both point to c6dfc39ac4f6dabddc5864ad611437b84ded383e, which contains ddda8706. I rewrote the update sentence: this SHA adds defaults.json defaults, failure-reason codes and retry folders; protocol, keys, agent ID and the friend model are unchanged; env vars set in the app still win. Line 123 is unc |
| F047 | medium | stale-doc | Same edit as F038. Both Grok headings and their anchors are unchanged. The Inbox-fix sentence is kept. |
| F046 | low | stale-doc | Current production lines now give c6dfc39a (2026-10-28) as a short SHA with a link to the D12 line (ops-guide.md#sar-autosetup-001-ops-운영-적용-결과-2026-10-28), in README.md, docs/README.md and docs/operations.md. README autosetup now says 0.1.1 is published and 0.1.0 or multi-plugin installs rerun the 0.1.1 command once (O-05). F-01–F-09 is now F-01–F-10. The docs/development.md make test row now inc |
| F048 | medium | stale-doc | README:34 now cites CS-AC-18 on bf67dfcb. README evidence links now point to SAR-AUTOSETUP-001-OPS-test and SAR-CONNECT-SIMPLE-001-LIVE-test, and the two HUMAN-PEER links are kept. README:44-45 became one bullet: calendar removed 2026-10-22 (7c6ef373), /privacy and /terms published in the same deploy. In docs/operations.md: front matter updated and the 3 OPS tasks added; the current-scope section  |
| F050 | low | logic-bug | fetch-audio.sh now runs mkdir -p "$(dirname "$0")/assets/audio" before the cd, so a fresh clone no longer fails. |
| F052 | low | logic-bug | package.json "beats" changed from 'node render.mjs --beats' to 'node render.mjs beats' (verifier 2's conservative one-line fix; the render script is kept). The README video examples (lines 35 and 49) now use a placeholder for a new version (e.g. v04) instead of the existing v02. |
| F053 | low | unnecessary-file | git rm of the root .npmrc, which duplicated adapters/.npmrc. Only archived .fullops-squad/docs/evaluations records refer to it (SAR-SETUP-001-* preview/result/scope.log, SAR-PUBLIC-MESSAGES-001-FIX-2-REVIEW snapshot-before.json, SAR-MVP-002-INSTALL-FIX-DEV-review report). Those were left as they are. |
| F054 | low | unnecessary-file | Partial. sqlc.yaml:1 comment changed to '# sqlc generates internal/database from db/migrations (schema) and db/queries.'. db/migrations/.gitkeep is kept on purpose (beta.sh:111 migrations-diff guard). db/queries/.gitkeep was left to the other agent, as instructed. |
| F089 | low | stale-doc | The GET / row in interface-design.md now says: always 200; login screen when there is no session; a valid session gets '로그인되어 있습니다.' and a /home link with no redirect; an invalid or expired cookie is deleted and the expired notice is shown. Front matter updated is bumped. |

건너뜀:

- F056: Needs a product/owner decision. Deleting renders v01/v02 reverses the documented keep decision (motion README:38, the 2026-10-18 handover, a2/report.md:15), and PLANS.md:213 shows the user has not yet accepted the video. Both verifiers say to get owner sign-off first. The files and README:38 are unc

### phase-b (커밋 `51b3eb52`)

| ID | 등급 | 종류 | 반영 내용 |
| --- | --- | --- | --- |
| F026 | medium | dead-code | Handler-side rate charging does nothing behind boundedHTTP: every in-handler 429 branch and bucket choice is unreachable |
| F059 | medium | dead-code | The per-handler rate budget (requestHit and memberHit's retry branches) never runs; about 100 lines of dead rate code |

## 독립 리뷰 지적

| 범위 | 등급 | 판정 | 내용 | 처리 |
| --- | --- | --- | --- | --- |
| ts-adapters | low | confirmed | Hooks turnOn: a watcher that has not bound within 20 s locks out every later 'on', and 'off' cannot stop it (adapters/hooks/register.ts:117) | `97195a76`에서 수정·회귀 시험 |
| scripts-deploy | low | refuted | transition.md still tells the operator to run the removed expose and access_apply.py check (.fullops-squad/docs/operations/transition.md:187) | 반박됨, 문서 일관성만 `97195a76`에서 보완 |
| scripts-deploy | low | refuted | ops-guide still lists the deleted tunnel/config.yml.tmpl as an OPS-owned file (.fullops-squad/docs/operations/ops-guide.md:314) | 반박됨, 문서 일관성만 `97195a76`에서 보완 |

