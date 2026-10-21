---
id: D10
title: 프로그램설계서
status: review
updated: 2026-10-21
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-MVP-002-DEV, SAR-MVP-002-BOT-CATALOG-DEV, SAR-MVP-002-BOT-CATALOG-DEV-FIX, SAR-MVP-003-BIDIRECTIONAL, SAR-PUBLIC-IDENTITY-001-DEV, SAR-PUBLIC-IDENTITY-001-DEV-RATE-FIX, SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG, SAR-PUBLIC-AGENTS-001-DEV, SAR-PUBLIC-AGENTS-001-DEV-FIX, SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX, SAR-PUBLIC-AGENTS-001-FIX-TESTER, SAR-PUBLIC-MESSAGES-001-DEV, SAR-PUBLIC-MESSAGES-001-DEV-FIX, SAR-PUBLIC-MESSAGES-001-TESTER, SAR-PUBLIC-MESSAGES-001-DEV-FIX-2, SAR-PUBLIC-MESSAGES-001-DEV-FIX-3, SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER, SAR-GOOGLE-LOGIN-001-DEV, SAR-GOOGLE-CONNECT-001-DEV, SAR-GOOGLE-CONNECT-002-DEV, SAR-AUTO-RECEIVE-001-DEV, SAR-SSE-INSTALL-001-DEV, SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV, SAR-GROK-WEBHOOK-001-DEV, SAR-SERVICE-UI-001-DEV, SAR-FRONTEND-POLISH-001-DEV, SAR-HUMAN-PEER-001-DEV, SAR-REFACTOR-001-DEV, SAR-REFACTOR-002-DEV, SAR-LANDING-GUIDE-001-DEV, SAR-UX-ISSUES-002-RECEIPT-DEV, SAR-SESSION-PEER-001-DEV, SAR-UX-ISSUES-002-NAMES-DEV, SAR-UX-ISSUES-002-APPROVAL-DEV, SAR-UX-ISSUES-002-CONNECT-DEV, SAR-UX-ISSUES-002-LANG-DEV, SAR-SCHEDULE-001-DEV, SAR-SCHEDULE-001-LOWS-DEV]
upstream: [D02]
summary: 실제 프로그램 책임과 요구사항 및 검증을 연결한다
---

# KnowsLink 프로그램 책임과 검증

## SAR-SESSION-PEER-001-DEV Claude Code 세션 연결

| 모듈 | 책임 | 검증 |
|---|---|---|
| `adapters/src/session-lock.ts` | agent 폴더당 한 수신 세션: `wx` 배타 생성, 다른 살아 있는 세션은 busy, `--take`·죽은 pid·같은 세션은 교체, 자기 lock만 해제 | `session-lock.test.ts`: claim/busy/take/release/dead holder/corrupt/alias 검증 |
| `adapters/src/text.ts` `watch(folder, session?, take)` | 세션 결속 watcher: bound/busy 이벤트, 10초 lock 감시 후 `replaced` 종료, SIGTERM/SIGINT 해제. 출력은 metadata만 | `inbox.test.ts` 6(기존), 격리 relay smoke(`watcher-smoke.log`: bound→takeover→replaced, `fix/watch-busy.out`: busy exit 2) |
| `adapters/src/core.ts` `agentId` | 비밀 없는 agent 주소 getter | 타입 검사 |
| `adapters/hooks/register.ts` | 모드: `/knowslink-session` 명령, 자식 spawn·이벤트 파싱, toast·상태줄, 미연결 text 도구 deny, 같은 세션 subagent·teammate 호출(`tool.call`의 `agentId` 있음) text 도구 deny, 발송 승인 대화상자(친구 이름·agent 이름·ID, 본문 전체·글자 수·reply_to, 비가시 문자 `⟨U+XXXX⟩` 표시·개수), 질문 2000자 초과 deny(Claude Code 2.1.302 대화상자의 표시 한도), 대화상자 중 연결 해제 시 deny, 인자 누락 deny, `replaced` 즉시 해제, 재연결 시 inbox 미확인 수, relay 통신 불가·수신 오류 상태줄·toast, 수신 뒤 대기 수 갱신, reload 시 상태줄 초기화 | `claude plugin validate adapters`, `hooks/register.test.ts` 11건(기존 8건에 정확히 2000자 질문 허용·초과 deny·4096 deny, 비가시 문자, 대화상자 중 off deny, 재연결 대기 수·relay 통신 불가 추가), 두 Claude 세션 V01과 SAR-SESSION-PEER-002 관측 |
| `adapters/src/text.ts` `contact`·`seen`·`printContact`, `mcp.ts` receive·`contact` CLI | 친구 agent 표시 이름 조회(O-03), 표시한 수신 text의 세션 수신 확인 보고(SP-12), 모드 대화상자용 `plugin.js contact <folder> <agent>`. 이름·확인은 실패해도 로컬 사본을 그대로 보인다 | `inbox.test.ts` 2(sender·sessionReceipt), Go `TestPublicMessageHTTPBoundaries` |
| `adapters/src/core.ts` fetch 실패, `inbox.ts` `safeError`, `mcp.ts` `publicFailure`, `text.ts` watcher relay 상태 | HTTP 응답 없는 실패(연결 거부·시간 초과)를 `relay unreachable`로 구분하고 5xx와 함께 '키·친구 문제가 아님' 안내. watcher는 10초 lock 감시 주기에 `knowslink_relay` ok/unreachable/error 변화를 알린다 (O-05) | `inbox.test.ts` 4(500 안내·닫힌 포트 unreachable) |
| `scripts/package_plugin.py` | Claude ZIP에 hooks 3파일 포함, `--connect` 시 Claude에만 `KNOWSLINK_AUTO_RECEIVE=off` | `--host claude --verify` 추출 MCP 검사·준비 경계 |

V05(외부 본문 명령 미실행·승인 없는 답장 0)는 모드 hook 수준의 연결 fixture 검사와 실제 세션 관측으로(모델이 본문 지시를 따르지 않는지는 자동화하지 않는다. 근거는 DEV QA 보고), V06(4096/4097)은 기존 `mcp.test.ts` textBoundaries로 확인했다. `claude plugin test adapters`는 `src/*.test.ts`(node 검사)까지 찾아 로드 실패로 표시하므로 모드 검사 결과는 `hooks/register.test.ts` 행만 본다. 패키지 추출 폴더에서는 모드 검사만 실행된다.

SP-04 해석(사용자 결정 2026-10-17, SAR-SESSION-PEER-001-DEV-FIX-2): 연결된 세션 안의 같은 프로세스 subagent는 부모와 같은 모드 상태를 보므로 연결을 그대로 쓸 수 있었다. 모드는 `tool.call` 입력의 `agentId`(Claude Code 2.1.302: subagent·teammate 루프에서만 있고 메인 루프에서는 없음)로 이를 구별해 KnowsLink text 도구(send/receive/receipt)를 deny한다. 메인 대화는 기존대로 쓴다. 별도 프로세스 자식은 연결 상태가 없어 기존 미연결 deny를 그대로 받는다.

V04 해석(사용자 결정 2026-10-17): 세션 연결(`/knowslink-session on`의 lock·watcher)은 소유자의 세션 설정이다. 친구 해제·차단 동안에도 연결은 유지되고 재수락 뒤 다시 `on`할 필요가 없다. 친구 권한은 relay가 발송마다 다시 확인한다. tester 관측(SAR-SESSION-PEER-001-TESTER, fullops/tester a32d0ab)에서 해제 중 A→B·B→A 발송은 모두 403이었고, 재수락 뒤 발송은 새 세대로 성공했으며, 해제 전 메시지의 receipt 조회는 `403 sender_not_allowed`였다. 따라서 V04의 '재수락 후 옛 연결'은 해제 전 친구 권한(과거 세대)을 뜻하고 세션 연결 자체를 뜻하지 않는다. 코드는 바꾸지 않았다.

## SAR-UX-ISSUES-002 이름 표시·친구 코드·친구 요청 안내

2단위(U2/U4와 해당 U7)는 새 schema·API·권한 없이 기존 template과 상태를 사용한다.

| 모듈 | 책임 | 검증 |
|---|---|---|
| `human_http.go` `distinguish` | 함께 보이는 agent 목록에서 짧은 구분값이 겹치면 전체 ID로 넓힌다. `humanAgents`·조회 선택기·legacy 목록(`ownAgents`)이 같은 규칙을 쓴다 | `TestReceiptSelectorDistinguishesCollidingNames`, `TestLegacyAgentLabelsWidenCollidingDiscriminators` |
| `human_http.go` `visibleAgent`·`friendsNow` | 자기 agent와 현재 활성·비차단 친구의 agent만 이름·종류를 보인다. 그 밖은 ‘확인할 수 없는 에이전트’와 전체 ID다. 상대 프로필을 새로 조회하지 않는다 | `TestVisibleAgentNamesOnlyOwnerAndCurrentFriend` |
| `http.go` `gateView`, `member.go` legacy, `member_agents.go` `ownAgents`·연결 화면 | 승인·legacy agent/관계·연결 대상에 이름 먼저 + 종류 + 구분값을 보이고 전체 ID는 상세에 둔다. legacy 제목은 ‘기존 연결’이다 | `TestGateNamesEndpointsWithinAuthority`, `member_forms.golden`, 렌더 fixture |
| `human_http.go` `friends/request` | 코드의 앞뒤 공백만 제거한다. 빈 코드·동의 누락은 422로 같은 친구 화면을 다시 보이고 입력 코드와 오류 필드 focus를 유지한다. 대소문자·Unicode 정규화는 하지 않는다 | `TestHumanNamesAndFriendCodePostgres`, `TestFriendCodeRotationKeepsReceivedRequests` |
| `human_pages.go` 친구·설정(옛 내 계정) | ‘내 친구 코드’를 보이는 라벨의 읽기 전용 `textarea`로 보인다(긴 코드도 줄바꿈). 친구 추가는 보내는 사람 이름·상대 코드·연결 범위·명시 동의·‘친구 요청 보내기’ 순서다. 복사 JS·CSP 변경은 없다 | `TestFriendPagesShareOwnCodeAndKeepRequestInput` |
| `member.go` `refusedTo` | 친구 동작·친구 상세의 거부 화면은 ‘친구 목록으로 돌아가기’로, 회원 gate 거부는 ‘승인 목록으로 돌아가기’(`/home/legacy#gates`)로 정확히 복귀한다 | 렌더 fixture |
| `gate_summary.go` `gateState`·`gateSummary`·`gateRows` | 표시 상태(만료·확인 불가·철회 포함), 지원 intent별 검증 필드→사용자 문구(필드 모양이 다르면 단언 대신 요약 없음으로 닫는다), 승인 목록 행(이름·구분값·유형·기한·상태)을 만든다. 이름 권한이 없으면 ‘확인할 수 없는 에이전트’와 ID 구분값만 쓴다. 결정 권한은 `gateDecision`의 기존 재검사만 판단한다 | `TestApprovalSummaryBoundaries`, `TestHumanGateKeepsVerifiedBodyAndCSRF`, 렌더 fixture |
| `connect_guide.go` `useChoices`·`agentResume`·`remoteResume` | U8 선택지 다섯 개(호스트 표 순서·문구, 보류·제외 호스트 이름 없음)와 agent별 연결 상태·다음 행동 하나를 만든다. 최신 연결 시도(waiting/prepared/approved)·활성 키·remote install의 미교환 code/token·device 승인 뒤 미저장만 본다. 앱 설치·세션 연결·온라인은 관측하지 않으며 불확실하면 ‘확인 필요’다. 표시 전용이며 상태를 바꾸지 않는다 | `TestConnectionResumeAndCodeTail`, `TestHumanPagesRenderAndEscape` |
| `member_agents.go` `create`·`memberBack` | 생성 `attempt`를 기존 `Idempotency` 맵에 `#agent-create/owner/attempt → agent`(`agentCreateKey`, agent ID가 `#`로 시작할 수 없어 메시지 키와 겹치지 않음)로 기록해 반복 제출이 새 agent를 만들지 않게 한다. 항목은 agent 삭제 sweep에서 함께 지운다. 거부 화면은 내 에이전트(legacy 관계는 기존 연결 목록)로 복귀한다 | `TestConnectionResumeAndCodeTail` |
| `human.go` `codeTail`·`humanFriend.CodeLabel` | O-NM-03 표시용 끝 4자를 요청 시점 입력에서 만든다. 요청자의 보낸·종료 행만 받는다. `decideFriend` block은 받은 대기 요청과 받은 사람에게 보류된 기한 전 요청을 `blocked-pending`으로 남겨 요청자 화면에서 차단이 드러나지 않게 한다. `requestFriend`는 코드마다 요청자 행 하나를 재사용하고, 상대 차단·수신 한도이면 같은 행을 `blocked-pending`(상대 비표시)으로 되살린다. 수신 한도 보류는 받은 사람이 보던 상태·방향·끝 4자(`Held`·`HeldFrom`·`HeldTail`)를 남기고, 같은 코드의 대리 행이 대기 중이면 쌍 행이 있어도 무동작이다 | `TestConnectionResumeAndCodeTail`, `TestFriendRerequestReusesOneRow`, `TestFriendBlockedRequestCountsTowardSentLimit`, `TestFriendLegacyPlaceholderAndStaleUnblock`, `TestFriendBlockOnHeldRequestKeepsRequesterPending`, `TestFriendRepeatWithOpenLegacyPlaceholderIsNoop` |
| `http.go` `verifiedGate` | 서명·digest·ID·endpoint·세대 결속에 더해 `gateKinds`에 없는 부모 intent를 검증 실패로 닫는다 | `TestApprovalSummaryBoundaries` |
| `human.go` `nameRefusal` | 이름 거부 문구를 고른다. 앞뒤 공백을 뺀 이름에 control/Bidi_Control 문자가 있으면 금지 문자 문구와 그 문자를 뺀 입력을, 아니면 ‘이름은 1–40자로 입력하세요’와 원 입력을 돌려준다. 내 이름·에이전트 이름·추가·device·remote 이름 422가 같이 쓰며 저장하지 않는다 | `TestHumanNamesAndIdentity`, `TestHumanNamesAndFriendCodePostgres` |
| `human_pages.go` 새 대화·내 에이전트·친구 요청 | 새 대화 발신·수신은 조회와 같은 fieldset/radio 라벨(`이름 / 종류 / 구분값`)이다. 내 에이전트 상세의 ‘요청 상태 조회’·‘연결·키 관리’는 각자 줄이다. 받은 요청 거절 앞에 거절 표시 안내가 있다 | `TestHumanEmptyChoicesOfferRecovery`, 렌더 fixture |
| `member_receipt.go` | 익명 한도 429 화면은 `Unavailable`로 ‘에이전트 없음’ 오안내를 하지 않는다(1단위 리뷰 low) | 코드 검토 |

동명·40자 한글/이모지·제어문자 거절·이름 변경 뒤 ID/관계 불변·코드 교체·GET만으로 관계 생성 없음·agent 0개 요청·24h 만료 뒤 새 세대는 `human_names_test.go`와 integration 검사로 확인한다. device/remote 생성 경로의 이름 결속은 기존 integration 검사에 단언을 추가했다. 검증 기록은 [DEV 보고](../evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-DEV-test/report.md)를 따른다.

## SAR-UX-ISSUES-002 설정·언어 모듈 (5단위)

한국어 원문이 키이고 영어는 표다. 템플릿 본문은 바꾸지 않고 소스 단위로 영어 템플릿 집합을 하나 더 파싱한다. 새 의존성은 없다.

| 모듈 | 책임 | 검증 |
|---|---|---|
| `i18n.go` `pageLang`·`browserLang`·`memberLang`·`setLang` | 회원 저장값 → `__Host-kl_lang` 쿠키 → 한 응답용 `?lang=` → 한국어 순서로 언어를 정하고 `Content-Language` 헤더에 기록한다. `browserLang`은 `boundedHTTP` 진입부에서 실행되므로 admission의 429/503 거부 화면은 회원 조회 없이 브라우저 쿠키 언어를 따른다(U3 LANG-VISUAL 결정). `memberHit`·공개 첫페이지·remote 동의가 회원을 확인하면 `memberLang`이 헤더를 회원 값으로 바꾸고, 저장값이 언어를 정했으면 내부 표식을 남긴다. `render`가 표식을 지우고 `LangSaved`로 넘기면 `lang-nav`는 브라우저 언어 링크 대신 ‘설정에서 변경’(`/home/account#language`)을 보인다. `boundedHTTP`의 `langMarks` 응답 래퍼가 헤더 전송 직전에 두 내부 표식을 지우므로 redirect·JSON처럼 render하지 않는 응답에도 나가지 않는다. 오류 문장(`intro`의 Problem)은 `langPart`가 다른 언어 이름에 `lang` span을 붙인다(O-LG-05). `/lang/{code}`는 쿠키 하나와 안전한 내부 복귀만 처리한다 | `TestLanguageChoiceOrderAndLink`, `TestSavedLanguageReplacesBrowserLinks`, `TestAdmissionRefusalsFollowBrowserLanguage`, `TestLanguagePartsAndMarkers`, `TestLanguagePreferencePostgres` |
| `i18n.go` `templateSegments`·`translateTemplate`·`parsePages` | 템플릿 소스를 태그 사이 본문·속성값 단위(제어 action에서 분리, 값 action은 포함)로 나누어 한국어 구간을 표로 바꾼 뒤 영어 집합(`memberPagesEN`)을 만든다. 누락 구간은 한국어로 남고 검사가 나열한다 | `TestTranslationCoversTemplates` |
| `i18n.go` `tr`·패턴 | 핸들러 문구의 정확 일치, 없으면 `%s`(다시 번역)·`%v`(그대로) 패턴. `render`가 Title·Problem·Notice에, 템플릿 `tr`/`name`/`badge` 함수가 상태·안내·에이전트 fallback 이름·배지에 적용한다. 사용자 이름·메시지는 `name`의 두 서버 라벨 외에는 거치지 않는다 | `TestTranslationCoversGoText`, `TestComposedHandlerTextTranslates` |
| `i18n_en_member.go`·`i18n_en_human.go` | 회원·공개·remote 화면과 사람 화면의 영어 표(파일 소유 분리). 값은 키와 같은 `{{…}}` action을 가져야 한다. 메뉴·제목 고유 이름(My agents·My name·My friend code) 밖에서는 사용자를 ‘your’로 부른다. 한국어 README 링크 문구에는 ‘(in Korean)’을 붙이고 템플릿 링크에 `hreflang="ko"`를 둔다 | `TestTranslationCoversTemplates`, `TestEnglishScreensHaveNoKorean`(지정 화면 영어 렌더에 언어명 외 한글 0, 고유 이름 밖 `my` 0, README 링크의 언어 표시) |
| `human_pages.go` `human-account`, `human_http.go` `language` | ‘설정’ 제목·헤더, 내 이름·내 친구 코드·언어(radio 한국어/English, ‘언어 저장’)·로그인과 보안 순서. 저장값이 없으면 언어 절에 ‘아직 계정에 저장되지 않았습니다’ 안내(`LangUnsaved`)를 radio 설명으로 연결한다. 저장은 `Member.Locale`만 바꾸고 선택만으로는 제출하지 않는다 | `TestLanguagePreferencePostgres`, `TestHumanPagesRenderAndEscape` |
| `human_http.go` `humanClock(lang)` | 짧은 대화 시각의 단어만 언어별로 바꾸고 KST 시각은 같다 | `TestHumanClockKeepsInstantAcrossLanguages` |
| `remote_auth.go` `remoteTemplates` | remote 동의·복귀 화면을 같은 템플릿 집합의 `remote`·`remote-return`으로 옮겨 언어 선택을 따른다 | `TestMemberUIFixtures`, `TestEnglishScreensHaveNoKorean` |
| `lang_browser_fixture_test.go` | opt-in `langfixture` 격리 fixture: 저장 English·저장 한국어·미저장 회원, 각 지정 화면에 닿는 합성 상태, loopback origin, 허용 GET/POST | DEV 보고의 실제 기동·probe |

## SAR-UX-ISSUES-002 조회 입력과 복귀

`member_receipt.go`는 최초·성공·실패의 공통 조회 화면을 만들고 기존 소유권/참여 endpoint를 검사한다. `receipt_context.go`는 길이·제어문자·쿠키 수명·회원 경계를 검사하고 소유 agent 선택기를 구성한다. 선택기는 2단위(SAR-UX-ISSUES-002-NAMES-DEV)부터 `fieldset`의 radio 목록이다. 각 라벨은 `이름 / 종류 / 안정 구분값`이며 줄바꿈되므로 닫힌 select처럼 잘리지 않는다(D-RT-01·U2). 1단위의 `구분값 / 이름 / 종류` 옵션은 대체됐다. 복귀 쿠키는 신뢰하지 않는 입력이다. DB 엔티티나 migration은 추가하지 않는다. `capacity.go`는 원래 admission transaction 안에서만 조회 폼을 인증하며 기존 예산·동시 처리 거부를 유지한다. DB 인증 전 실패는 사적 입력 없는 같은 조회 화면이다.

`member_pages.go`의 네이티브 GET 폼은 안전한 입력·선택을 유지한다. 각 응답은 오류 필드/오류 요약/성공 결과 중 한 곳에만 autofocus를 둔다. `human_pages.go`의 소유 agent 상세와 legacy 관리 영역에서 같은 화면으로 진입한다. D04 CSS·주 메뉴·CSP는 바꾸지 않는다. `member.go`/`google.go`는 기존 인증 성공 뒤 고정 조회 복귀와 로그아웃 초기화만 추가한다. Google의 외부 callback에서 Strict 쿠키를 사용할 수 있다는 가정은 하지 않는다.

자동 검사는 `receipt_context_test.go`, `receipt_integration_test.go`와 기존 회원/Google/인증 회귀를 사용한다. member_forms.golden의 변경은 legacy 인라인 조회를 링크로 대체하고 hidden 성공 폼을 편집 폼으로 바꾼 두 계약이다. `receipt_browser_fixture_test.go`는 integration,receiptfixture 명시 태그에서만 합성 A/B/C·429·503·만료·철회 환경을 제공한다. 실제 검증 SHA·원 실패·종료코드·후속 담당은 [DEV 보고](../evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md)를 따른다.


## SAR-LANDING-GUIDE-001 공개 안내

`landing_pages.go`는 D04의 공개 소개·시작 순서·설치와 연결·가상 예시·단일 로그인 영역을 제공한다. `member_pages.go`는 기존 memberStyle과 함께 이 템플릿을 등록한다. 공개 전용 CSS는 회원 화면의 제목·배치를 바꾸지 않는다.

`member.go`의 `startPage`는 유효 세션도 공개 안내에 머물게 하고 존재 여부만 전달한다. 이름·이메일·친구·agent·키는 안내에 넣지 않는다. `/home`의 이름 설정과 회원 인가는 유지한다. `Service.startView`는 Google/Mail 설정을 모든 start 오류 화면에 반영한다. `google.go`의 오류 호출자도 이 view를 사용한다. Google callback·재확인·클라이언트/Remote 복귀 계약은 바꾸지 않는다.

`human_pages.go`의 설정(5단위 전 ‘내 계정’)에는 `/`의 서비스 안내 링크가 있다. 기존 회원 주 메뉴·인증 폼 계약은 유지한다. `landing_test.go`와 `landing_integration_test.go`는 로그인 수단·공개 카피·오류 escape·실제 DB 세션·이름 설정·만료·개인정보 비노출을 검증한다. 기존 member form golden과 Google/device/Remote 회귀도 적용한다. opt-in `landingfixture` 태그의 fixture는 별도 합성 DB와 loopback 서버만 사용한다. 실제 검증 SHA·명령·결과는 [DEV 기록](../exec-plans/phases/SAR-LANDING-GUIDE-001-DEV.md)에 연결한다.

## SAR-REFACTOR-002 감사·개선

| 모듈 | 변경·회귀 |
|---|---|
| store.go sweep | 만료 이력의 current 탐색 생략. store_retention_test의 만료 정각·철회·delivered·claim·24h 보관과 규모별 벤치 |
| remote_events.go claimRemote | 호출 내 agent별 pending 재사용, 마지막 lease 소진 후 새 pending의 즉시 재예약. remote_events_test의 crash 상태 재현과 다음 호출의 권한 재검증 |
| adapters/src/inbox.ts take/prune | housekeeping을 읽음 claim 전에 수행하고 stat ENOENT만 허용. inbox-files.test의 EACCES 본문 보존·ENOENT 수신·중복 거부 |

새 계층·라이브러리·저장 schema·UI 변경은 없다. 테스트 등록만 package.json에 추가했다. [감사](../evaluations/qa-reports/SAR-REFACTOR-002-DEV/audit.md)와 [실행 보고](../evaluations/qa-reports/SAR-REFACTOR-002-DEV/report.md)는 과거 QA와 이번 합성 실행을 구분한다. 전역 JSONB 처리량·외부 host 표시·운영 부하를 보장하지 않는다.

## SAR-REFACTOR-001-DEV 기능 보존 정리

`human_http.go`의 `humanAgents`는 자기·친구 화면의 동일한 owner별 활성 agent 선택과 ID 정렬을 공유한다. `TestHumanAgentListsKeepOwnerAndFriendBoundaries`가 다른 owner·철회 agent 제외, 키 없는 agent 표시, 차단 뒤 접근 거부를 확인한다. template·문구·권한 판정은 유지한다.

`device.go`의 `connectionResponse`는 device API와 connect API의 동일 응답 처리를 공유한다. `TestConnectionResponseKeepsRetryAndErrorPriority`가 retry 우선순위·Retry-After·명령 오류·transaction 오류·정상 응답을 확인한다. route·transaction·rate 집행은 유지한다.

adapter는 `core.ts`의 `validText`로 일반·시험 text의 공백·UTF-8 4096bytes·고립 surrogate 검사를 공유한다. `mcp.ts`는 기존 `trialResult`로 동일한 MCP text 응답 포장을 재사용한다. wire·서명·held/trial 상태와 수신 순서는 유지한다. 실행 근거와 보존 후보는 [리팩터링 기록](../exec-plans/phases/SAR-REFACTOR-001-DEV.md)과 [QA 보고](../evaluations/qa-reports/SAR-REFACTOR-001-DEV-test/report.md)에 기록한다.

## SAR-HUMAN-PEER-001 DEV 후보

| 모듈 | 책임·검증 |
|---|---|
| `human.go` | 이름, 신원 자격, legacy 변환, authority, Friend 상태/한도/세대, 당시 사본. `human_test.go` |
| `human_messages.go` / `public_text.go` | owner 작성 wire와 signed agent text, 독립 TTL·24h 이력·멱등. owner를 서명된 agent로 표시하지 않음 |
| `human_http.go` / `human_pages.go` / `member_pages.go` | 새 전체 제품 화면·세션 폼·D04 공용 토큰. raw 입력 escape, 오류 회복·필드 labels·44px 동작 |
| `store.go` / `http.go` / events·remote | 출처+세대 current 검사, 기존 gate/claim/CSRF/CSP/rate/키 유지 |
| `00002_human_writer.sql` / `scripts/schema.py` | 구버전 JSON writer 차단, 전체 migration D08 생성 |
| `adapters/src/text.ts` / `owner-text.test.ts` | owner wire schema/대상/기한 검증, 서명 분리·untrusted 결과, 저장 실패 시 ACK 차단 |
| `human_integration_test.go` | 실제 합성 OTP·profile·명명·connect proof·A/B/C·친구·연속6회·signed·철회·멱등·CSRF·DB downgrade |

기존 fixture는 유지된 legacy 화면/정확한 관계와 새 제품 의미를 명시적으로 구분한다. 독립 tester/designer/review의 최종 수락은 coordinator 후속이며 DEV 자체 검사를 대신하지 않는다. [DEV 실행 기록](../exec-plans/phases/SAR-HUMAN-PEER-001-DEV.md)을 따른다.

## SAR-SSE-INSTALL-001 SSE 모듈

| 모듈 | 책임 | 핵심 검사 |
|---|---|---|
| `text_events.go` | current-auth/pending 선택·shared stream 입장·전용 LISTEN·빈 SSE·heartbeat·취소 | TextEventSelectionAndAuthorization, TextEventSharedBudgetAndExpiry, TextEventsHTTP |
| `store.go` transaction | public 자격/queue 변화 hash 비교·commit과 빈 NOTIFY 결속·Streams 만료 | cross-process HTTP 도착·철회·heartbeat 무변경 |
| `core.ts` openEvents·`text.ts` listen | redirect 거부·Bearer header·UTF8/크기 제한·빈 frame schema·abort | events.test.ts의 실제 HTTP·분할 frame·payload frame 거부 |
| `inbox.ts` AutoReceiver | 선택 SSE·coalesce·기존 serial·poll fallback·reconnect catch-up·stop | events.test.ts와 기존 inbox.test.ts |
| `mcp.ts` | 기존 도구9개·자동 수신 상태·stdio close cleanup | 기존 MCP initialize/discovery/held와 Inbox 회귀 |

상세 설정과 실패 경계는 [D05](interface-design.md#sar-sse-install-001-sse-계약), 실제 검사 결과는 [DEV 기록](../exec-plans/phases/SAR-SSE-INSTALL-001-DEV.md)을 따른다. 운영 서버 적용과 무인증 거부·metadata smoke는 [운영 LIVE 기록](../evaluations/qa-reports/SAR-SSE-INSTALL-001-OPS-LIVE.md)이다. 인증 SSE의 운영 도메인 확인은 실제 계정 후속이다.

상위 요구사항은 [D02 SAR-MVP](../planning/product-specs/SAR-MVP.md)다.
API 정본은 [D05](interface-design.md), 저장 정본은 [D06/D07/D09](data-model.md)다.

| 모듈 | 책임·외부 경계 | 직접 검증 |
|---|---|---|
| protocol | Strict·Parse·SigningBytes·Digest, registry seed schema | TestFrozenParsingAndSigning |
| registry | compiled seed SHA256·manifest 공개 | protocol test·TS once |
| store | Postgres transaction·epoch·clock·current·sweep·ingest·lease | TestPostgresSafety |
| HTTP owner | credential·PoP·pairing·Go template·CSRF decision | PostgresSafety, GateFailureStates |
| HTTP agent | `deliver:agent` 전용 pull·shared persist·ACK·claim·authorize·result | PostgresSafety, ApprovalAndResultInstanceBinding |
| cleanup | 유휴 payload·metadata 회수 | 동일 sweep 경계·runtime 기동 |
| sqlc database | pgx/v5 LockRelay·SaveRelay 생성 | generate diff, 실제 Postgres |
| TypeScript adapter | registry·signature 재검증, persist·ACK·claim·gate·denied result | synthetic.ts 실제 HTTP 검사 |
| migrate·config·health | 별도 SQL-only Up·설정·readiness 실패 | 기존 unit/race·verify-runtime |
| Compose 검증 | 고유 project·private 제품 DB·loopback 시험 DB·Tunnel OFF | verify-mvp·verify-runtime |

`make test`는 외부 DB가 없어도 protocol과 기존 회귀를 수행한다. Postgres 검사는 integration build tag로 별도 실행하며 DB 환경이 없으면 실패한다.
`make verify-mvp`는 별도 DB에서 12회 ingest 경합·8회 claim 경합·3번째 lease·late ACK·pool 재시작을 검사한다. Go integration 검사 동안 같은 project의 relay를 멈춘다. relay 정리 sweep의 시험 allowlist가 검사용 trial lease를 회수하기 때문이다(TestTrialForeignAllowlistRevokesLease).
권한 철회·세대 교체·TTL rollback·시계 이상·CSRF·중복 gate·M.id 재사용·잘못된 endpoint·optional 결과 거부를 검사한다.
`agent_cannot_process_human_delivery`는 직접 `deliver:human` send 403과 저장된 human 전달의 agent pull·persist·ACK·claim 거부를 검사한다.
`legacy_unrouted_claim_parent_boundaries`는 경로 미기록 claim의 HTTP authorize·H·R·gate-consume 403과 경로 기록 뒤 정상 처리를 검사한다.
`TestLegacyClaimsCannotReachParentBoundaries`는 a6a10c7 State 메서드로 만든 `testdata/legacy_claims.json`을 읽어 같은 경계와 새 agent 정상 경로를 검사한다.
TypeScript 검사는 Go 서버를 통해 policy 없음 deny와 gate approve 후 deny 및 최소 R 수신을 끝까지 수행한다.

독립 QA는 tester가 고정 후보 SHA에서 QA-01–11을 수행한다.
직접 시각 검수는 designer가 V-01–04를 수행한다. DEV의 자동 HTML 검사는 독립 시각 검수를 대체하지 않는다.
별도 fixed-SHA 코드 리뷰와 critical/high 차단은 coordinator가 담당한다.


### Remote OAuth·Cloud 설치

`remote_auth.go`는 metadata·DCR·PKCE·Google consent·token rotation/revocation·설치 권한을 담당한다. `remote_mcp.go`는 MCP discovery와 네 도구를 기존 ingest/lease/서명/persist/ACK에 연결한다. `remote_events.go`는 지속 구독·서명 callback·claim/retry·LISTEN 수명주기를 맡는다. `capacity.go`는 /mcp 설치별 member budget과 consent session budget을 적용한다. remote token은 기존 /v1 credential이나 reserved cleanup 권한이 아니다. `cmd/relay`가 Cleanup과 RemoteEvents의 공통 종료 context를 관리한다. consent 승인·거절은 `remoteReturn` 동일 출처 화면으로 callback에 이동한다. 공통 CSP의 form-action은 넓히지 않는다. 실제 Chromium 경계는 `-tags "integration browser"`와 `TEST_BROWSER`로 `TestRemoteConsentBrowserReturn`이 검사한다.

## SAR-MVP-002 플러그인 모듈과 검사

- `adapters/src/core.ts`: 기존 Adapter와 signing·canonical 함수를 재사용한다. import 시 CLI를 실행하지 않는다. localAdapter는 loopback 구성과 키 파일을 읽는다. redirect를 따라가지 않는다.
- `adapters/src/index.ts`: 기존 CLI 진입점을 유지하고 core export를 제공한다. unconfigured 출력과 합성 직접 실행의 기존 계약을 보존한다.
- `adapters/src/mcp.ts`: 공식 SDK의 McpServer·StdioServerTransport·registerTool을 사용한다. 기본 held와 synthetic-only 경계를 집행하며 stdout에는 MCP만 기록한다. gate 로그를 model/tool output으로 전달하지 않는다.
- `adapters/.cursor-plugin/plugin.json`, `mcp.json`, `skills/knowslink/SKILL.md`: 공식 Cursor plugin 구조로 MCP와 사용 안내를 연결한다. 설치 기본값은 held다.
- `scripts/package_plugin.py`: 허용 목록으로 standalone marketplace ZIP을 만든다. bundle과 manifest·skill·설치 문서만 포함한다. 고정 ZIP timestamp로 같은 내용의 SHA256을 유지한다.
- `adapters/src/mcp.test.ts`: 실제 bundle의 MCP handshake·discovery·held/no-network·잘못된 모드/URL·redirect 거부를 검사한다. 별도 압축 해제 artifact 경로도 검사할 수 있다.
- `adapters/src/synthetic.ts`: 기존 SQL 합성 흐름에서 owner gate 경로를 실제 MCP bundle 호출로 검증한다. PEM은 자기 임시 0700 폴더의 0600 파일에 두고 처리 후 제거한다.

[플러그인 설치 문서](../../../adapters/README.md)와 [공식 조사·실행 기록](../exec-plans/phases/SAR-MVP-002-DEV.md)에 버전·실패 수정·한계·후속 담당을 기록한다. TESTER 독립 QA와 fixed-SHA 독립 리뷰는 coor 후속이며 이번 자동 검사로 대체하지 않는다.

## SAR-MVP-002-BOT-CATALOG-DEV 앱 등록 경로

- 관측: 재시험에서 CLI 설치·doctor는 성공했고 앱 카탈로그에는 knowslink가 없었다. 앱이 CLI plugin(`~/.grok`)을 읽지 않는다는 원인은 미확정 가설이다.
- KnowsLink는 custom MCP server **Command** 등록을 시도한다. 공식 근거는 Team Bots 문서의 Setup → Plugins → Add·채팅 요청과 Plugins 표에 한정된다. 개인 계정 UI·승인 카드는 미확인이다. 이슈1에서 Bot은 `AddMcpServer`를 호출할 수 없었다.
- `scripts/install_bot_mcp.sh`: Linux x86_64/aarch64용 고정 SHA256 Node `v22.22.2` `.tar.gz`와 bundle을 `/workspace/.knowslink`에 준비한다. 새 `.stage.*` 폴더에서 체크섬·빈 환경 MCP 경계 검사를 통과한 뒤에만 `node`·`knowslink`를 교체한다. 실패하면 기존 준비물을 보존한다. 상대 `KNOWSLINK_PREFIX`와 스크립트가 만들지 않은 `node`·`knowslink` 항목은 변경 전에 exit 1로 거절한다. PREFIX의 다른 파일은 건드리지 않는다. 환경 변수가 없으므로 `mcp.ts`의 기본 held를 사용한다.
- `scripts/verify_grok_plugin.py`: CLI 설치 검사다. 앱 카탈로그 증거가 아니다. 실패 시 grok 출력을 표시하고 `GROK_CONFIG*` 변수를 제거한다.

근거·가설·검증은 [실행 기록](../exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md)과 리뷰 보완 [FIX 기록](../exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV-FIX.md)을 따른다. 원인 판정은 FIX 기록이 우선한다.

## SAR-MVP-003 시험 transport 모듈

- `internal/relay/test_messages.go`: 두 시험 identity allowlist와 machine prefix 허용 경로를 집행한다. normal mux에 넘기는 request copy에서 owner API 진입을 차단한다.
- `protocol.go`/`registry.json`: 시험 text closed schema와 별도 registry revision을 추가한다. 기존 업무 intent는 유지한다.
- `store.go`/`http.go`: current-auth·서명·active pair·TTL·중복을 재사용한다. 시험 전용 pull·persist·ACK·claim을 집행하고 claim 시 payload를 삭제한다. rollback도 동일 allowlist를 적용한다.
- `adapters/src/core.ts`: 서명과 registry 검증을 synthetic/test receive가 공유한다. redirect 차단·요청 timeout·응답 스트림 상한을 공통 request에 적용한다.
- `test-transport.ts`/`mcp.ts`: configured peer로만 trial send한다. 수신 text는 claim 뒤 untrusted data로 반환한다. 기존 업무 pull 도구는 synthetic 모드를 유지한다.
- `trial-cli.ts`/`scripts/run_trial.py`: Codex 송신 stdin·수신 CLI와 private env 로딩을 제공한다. mode·credential은 tool/argv로 받지 않는다. launcher는 0600 소유 regular config를 요구하고 이전 agent env를 상속하지 않는다.
- `trial-setup.ts`: local beta operator만 두 신규 시험 owner/agent·Ed25519 PoP·pair를 등록한다. private 목적지 재사용을 거부한다. agent별 키와 config를 분리하고 owner 기록은 운영자에게만 둔다.
- `deploy/knowslink/access_trial_plan.py`: 단기 distinct service token·reusable non_identity policy·path 앱 본문과 보호된 Tunnel 후보를 생성한다. live mutation은 수행하지 않는다.

검증은 TestTrialMessageSafety·TestTrialSchemaAndConfiguration·TestTrialHTTP, trial-boundaries.test.ts, trial-check.ts로 연결된다. `make verify-mvp`는 기존 업무/gate 회귀와 두 독립 MCP 프로세스의 실제 Postgres 왕복을 함께 검사한다. `make verify-grok-plugin`은 네 tool을 기대한다. default held·실제 계정 수락의 분리는 계속 유지한다. [실행 기록](../exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md)에 ID·명령·exit code·남은 실제 조건을 기록한다.

## SAR-PUBLIC-IDENTITY-001 신원 모듈

- `internal/relay/identity.go`: 이메일 정규화·마스킹, rolling 한도(`take`/`hit`), 코드 발급·검증, 회원·owner 바인딩, 세션 수명, 정리.
- `internal/relay/member.go`: 시작·확인·홈·재확인·로그아웃 화면과 HTTP 상태. 메일 발송은 lock 밖에서 수행한다.
- `internal/relay/mail.go`: 표준 `net/smtp` submission. implicit TLS와 STARTTLS를 지원한다. PLAIN 인증은 TLS 또는 localhost에서만 보낸다.
- `internal/relay/http.go`: `http.CrossOriginProtection`, 합성 가입 설정, Basic owner와 회원 세션이 공유하는 gate 화면.
- `cmd/relay/main.go`: SMTP·발신 주소·client IP header·합성 가입 설정을 읽는다. 값은 출력하지 않는다.
- `scripts/mail_sink.py`: 로컬 QA 전용 SMTP sink. 받은 메일을 0600 파일로 저장한다. 실제 발송 서비스가 아니다.

| 검사 | 내용 |
|---|---|
| `TestNormalizeEmail`, `TestRollingWindowBoundary`, `TestSendLimits` | 주소 형식·별칭·rolling 경계·all-or-nothing·거부 집계·60s/5/20/100 한도·NAT 분리·rate key 원문 미보관 |
| `TestRatePrincipalIsolation` | 한 IP·회원의 반복 거부가 다른 IP·회원의 신규·정리를 막지 않음, 자기 한도 30/40/20 이하·상한·다음, 독립 principal 합의 전체 200·정리 100 포화, rolling 회복 |
| `TestRateStateBoundedUnderRotation` | 공유 포화 뒤 새 source 10000개의 거부가 rate key·공유 bucket을 늘리지 않음, 기존 principal은 limit+1 보관, 회전이 포화를 연장하지 않음, 재시도 시각 |
| `TestCodeVerificationAndMembers`, `TestSessionLifetime` | 오답 5회·만료·재사용·확인 전 owner 미생성·재로그인 연속성·issuer 분리·회원 100 수용량·절대/무활동 수명 |
| `TestSMTPMailer` | 설정 오류의 비밀값 미노출·이름 표기 거부·실제 SMTP 대화·연결 실패 |
| `TestEmailIdentity`(integration) | 실제 Postgres HTTP 흐름: 가입·홈·재로그인·재시작·동시 첫 가입·1회 코드·로그아웃 재사용 차단·전체 로그아웃 재확인·무활동 만료·agent credential 유지·발송 실패·429·client IP header·cross-site 403·회원 gate 격리와 결정·합성 가입 차단·한 IP/회원 flood(신규 250·정리 150) 뒤 재시작에도 자기 429 유지와 다른 source 가입 시작·다른 회원 홈·로그아웃 허용 |

## SAR-PUBLIC-AGENTS-001 모듈·검증

| 파일 | 책임 | 검사 |
|---|---|---|
| connections.go | owner/agent/client bound PoP·10분 grant·key credential·quota·pending expiry | TestConnectionApprovalAndKeyCredentials, TestConnectionFailureExpiryAndCancellation, TestAgentAndPairCapacity |
| member_agents.go·member.go | 자기 agent/지문·지원 client·확인/취소 UI. 관계 세대는 hidden 값으로만 전달한다(SAR-SERVICE-UI-001-DEV 이후 화면 숫자 표시 없음) | TestPublicAgentHTTP; UX-04–05 직접 검수 후속 |
| api_rate.go | agent API stable principal·별도 cleanup budget·DB 재시작 유지 | 기존 TestEmailIdentity, TestPublicAgentHTTP·TestPostgresSafety |
| http.go·store.go | 공용 owner 동작과 키별 current-auth·철회·pair cap | 기존 업무/게이트/시험 회귀와 새 key credential 검사 |
| adapters/src/connect.ts | Node local prepare/complete·로컬 private key·0700/0600 설정 | connect.test.ts 실제 HTTP·PoP·파일 권한·URL·1회 완료 |

TestPublicAgentHTTP는 실제 Postgres에서 타 회원 ID 바꿔치기·agent로 owner권한 우회·cross-origin 확인·재인증·회전·철회된 credential의 send/pull/persist/ACK/claim/authorize/gate-consume/result·same-owner 명시 수락·동시 수락·새 세대와 옛 화면 거부·pending 만료·동시 agent cap·정리 budget·재시작·동시1회 완료·취소/만료를 검사한다. make verify-mvp의 relay-stop → Go → relay-up 순서를 유지했다. 과거 invalid_lease 원인 근거는 기존 실행 기록을 유지한다.
Node 검사는 mock relay와 실제 crypto/파일 I/O를 사용한다. Go 검사는 actual Postgres HTTP다. 실메일·실제 플랫폼·공개·독립 QA·직접 시각 검수와 구분한다. 새 라이브러리는 없다. 기존 버전의 stdlib Ed25519·crypto·filesystem·Adapter.request 사용 근거를 재사용했다.

## SAR-PUBLIC-AGENTS-001-DEV-FIX 모듈·검증

| 파일 | 변경 | 검사 |
|---|---|---|
| connections.go | 보존 상수, agent·키 기록 상한, `revokeKeys`, 철회 agent 24h 삭제와 관련 pair 삭제 | TestRevokedRecordRetention, TestPublicAgentHTTP/revoked_record_saturation_concurrency_restart_and_retention |
| store.go | Agent.Changed. 거부 transaction도 rate 기록을 저장 | TestPublicAgentHTTP/invalid_session_pages_spend_anonymous_budget |
| member.go·member_agents.go | `memberHit`: 세션 확인과 회원/익명 budget 소비를 한 곳에서 수행. connect 429 응답 공용화. agent 철회 시각 기록 | 같은 integration, malformed_connection_requests_spend_anonymous_budget |
| api_rate.go | `rateLimited`: 실제 retry 시각을 초 단위 올림으로 header·body에 표시 | malformed_connection_requests_spend_anonymous_budget, 기존 TestEmailIdentity |
| http.go | 합성 keys 경로 키 기록 상한, 반복 key-revoke의 철회 시각 유지, 없는 agent pair 거부 | TestRevokedRecordRetention |
| member.go·member_agents.go(UI) | `refusal` template·`refused`: 거부 화면의 홈/로그인 링크와 재확인 버튼. 지문 `<code>`, 줄바꿈·select 스타일, 열린 연결만 취소, `memberPair.Deadline` KST | TestMemberPagesShowNextSteps, 기존 TestPublicAgentHTTP·TestEmailIdentity 문구 검사; designer 좁은 재검수 후속 |

TestRevokedRecordRetention은 State 단위 검사다. 회전 반복의 키 기록 상한, 포화 중 철회 허용, 반복 철회 시각 유지, 30일 뒤에도 살아 있는 agent의 철회 kid 유지, owner agent 기록 상한, JSON 재시작 뒤 24h 삭제와 관련 pair 삭제, 살아 있는 pair 세대 유지, 기존 철회 agent의 보존 시작, 없는 agent pair의 nil 역참조 방지, 합성 키 경로 상한을 확인한다.
integration은 실제 Postgres HTTP에서 동시 생성 8건 중 기록 상한이 정확히 4건만 허용하는지, 재시작 뒤 상한 유지, 포화 중 agent 철회 303, 24h 경과 agent만 삭제되는지를 확인한다. 무효 세션 GET 30회 뒤 재시작한 relay에서 GET·logout·reauth가 429인지 확인한다. connect 429의 retry_at과 Retry-After가 같고 60초 안의 미래 시각인지 확인한다.

## SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX 모듈·검증

POLICY 48d12fa의 관계 반복·재초대 의미는 기준 4a 코드와 이미 일치했다. 코드 의미는 바꾸지 않았다. 관찰 조건을 결정적 검사로 고정하고 회원 화면의 상태·다음 동작 안내만 추가했다.

| 파일 | 변경 | 검사 |
|---|---|---|
| connections.go | `agentLimit`: 새 agent를 막는 한도 이름(total·records·active)을 반환한다. `agentCapacity`는 이 함수를 사용한다. 한도 값·순서 집행은 같다 | TestSaturationGuidance, 기존 TestAgentAndPairCapacity·TestRevokedRecordRetention |
| member_agents.go | `agentLimits`·`keyRecordsFull`·`connectLimit`·`agentProblem`: create·connect의 409를 실제 한도별 안내로 바꾼다. `inviteNotice`: 초대 제출 결과를 새 pending·기존 pending 유지·active 유지로 구분한다. `memberAgent.KeyFull`, `memberPair.Own/Other` | TestSaturationGuidance, TestMemberPagesShowNextSteps, TestPublicAgentHTTP 두 하위 검사 |
| member.go | `refusal.create`와 거부 화면의 새 에이전트 만들기 버튼(당시 문구 "새 agent 만들기"), 초대 notice 3종, 홈의 철회 agent 보존 안내·철회 전 경고·키 기록 포화 안내, 관계 상태별 안내·종료 관계의 수동 새 초대 버튼, 종료 관계의 관계 철회 버튼 제거 | TestMemberPagesShowNextSteps, TestPublicAgentHTTP |

TestRelationshipPolicy는 State 단위 관계 행렬이다. 같은/반대 방향 pending 반복의 수·세대·기한·수신자 불변, pending 메시지 거부, 발신 owner 결정 거부, 기한 직전 수락과 기한 도달 거부, 만료 뒤 새 기한·새 세대, active 반복 불변, 양측 각각의 unpair와 옛 메시지·옛 결정 거부, 양방향 재초대의 새 수락과 새 세대, 거절 뒤 늦은 수락 거부와 pending slot 해제, 타 owner 결정·unpair 거부, 같은 owner 자동 수락 없음, 송신 pending 상한의 다음 재초대가 세대를 만들지 않음, 재시작 뒤 상태 유지, 비활성 owner·철회 agent 재초대 거부를 확인한다.
TestSaturationGuidance는 키 기록 포화의 교체 안내와 생성 가능 여부, owner 기록 포화의 24h 보존 전후, 보존 중·정리 뒤 옛 키 credential 거부, 활성 키 3개·활성 agent 5개 안내, 결제·정확한 시각 문구 부재를 확인한다.
integration은 실제 Postgres HTTP에서 초대 notice redirect 3종, 반대 방향 반복의 기한 불변, 종료 관계의 새 초대 버튼, owner 기록 포화 409 안내, 키 기록 포화 connect 409의 교체 안내·생성 버튼·연결 미생성을 확인한다.

## SAR-PUBLIC-AGENTS-001-FIX-TESTER 독립 QA

독립 QA는 후보 `458798c2ee15c179edacfd6f94ebb9896d26f411`의 detached clone에서 회원 HTTP와 Postgres 행을 확인했다. DEV의 `TestRelationshipPolicy` 통과를 이 QA의 통과로 쓰지 않았다.
보존, rate, POLICY 두 관찰 표의 프로브 종료코드는 0이다. 새 critical/high는 없다. 상세는 [QA 보고서](../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-FIX-TESTER.md)다.
원본 d1 QA와 UI FAIL은 원래 SHA에 둔다. 이 절은 그 판정을 바꾸지 않는다. 운영 공개와 실제 메일은 확인하지 않았다.

## SAR-PUBLIC-MESSAGES-001 모듈과 검증

| 파일·모듈 | 책임 | 자동 검사 |
|---|---|---|
| public_text.go | 별도 closed wire·서명·digest·양측 회원/현재 권한·1회 관련 답장·text 전용 HTTP | TestPublicTextWire, TestPublicTextRoundtripAndCurrentAuth, TestPublicTextTTLAndLeaseFailures, TestPublicMessageHTTPBoundaries |
| capacity.go·store.go | 공통 queue100/gate100/receipt20000/claim4·shared HTTP16/4·owner당 Clean 1·본문 선수신 입장·TTL 정리·현재 세대·(commit 시각, epoch) 순서 snapshot(remember)·종료 실패 기록 회수(reclaim)·owner 제어/agent ACK 공정성 단위(agentsUnit) | TestSharedCapacitiesAndCleanupClassification, TestSharedMessageCapacityHTTP, TestSharedExecutionClaimHTTP, TestHTTPConcurrencyAcrossInstances, TestHTTPSlowBodyAndCleanupAdmission, TestValidCredentialCleanupFlood, TestFailedFinishReclaimedAtNextCommit, TestOwnerControlNotHeldByOwnAgentAcks |
| cleanup_admission.go | 정리 입장 판정: 무상태 자격·같은 출처·gate CSRF·본문(cleanupCredential), 자기 기록(cleanupTarget), 커밋 snapshot의 자격 색인(kind·principal)으로 DB 전 자기 기록 대조와 공정성 단위 결정(cleanupOwner), 신규 포화 중 증명 못한 정리 자격의 lock 없는 단일 재조회(refresh, ReadRelay) | TestCleanupAdmissionNeedsVerifiedOwnRecord, TestRememberKeepsNewestCommit, TestCleanupSnapshotAcrossProcessesAndRestart, TestCleanupProvenWhileSnapshotStaleAndNewFull, TestValidCredentialCleanupFlood |
| member_receipt.go·member.go | 자기 agent/요청별 metadata·queued/수신/답장/실패 구분·수동 안내 | TestMemberReceiptView, TestPublicNodeProcesses, TestPublicMessageHTTPBoundaries |
| http.go·member_agents.go·api_rate.go | 공통 rate 차감·gate 검증본문/서명·deny 입장·현재 owner/CSRF·결정 뒤 정식 gate 화면 | TestPublicHTTPAdmissionAndGateSafety, TestEmailIdentity/member_gate_deny_form_returns_to_result와 기존 신원/gate 회귀 |
| adapters text.ts·core.ts·connect.ts | private 연결 폴더·명시 text CLI·서명 확인·persist/ACK·안전 오류 코드·relayBase 재사용 | TestPublicNodeProcesses의 실제 Node CLI/MCP + 기존 adapter 검사 |
| mcp.ts·public-check.ts | public-node opt-in·명시 승인 send/관련 reply·수동 receive/receipt·untrusted 데이터 | 실제 독립 MCP 두 프로세스의 ID 대조와 mcp.test.ts held/discovery |

기본 MCP 도구는7개이며 기능 활성은 각 모드의 명시 설정에 따른다. bundle에 CLI 진입점이 실행되지 않도록 text.js 이름도 확인한다. connect.ts의 relayBase는 core.ts에 옮겨 bundle에서 연결 CLI가 실행되는 import 부작용을 제거했다. 원본 실패와 교정 증거를 실행 기록에 보존한다.
로컬 실제 프로세스 왕복을 실제 Grok Bot/다닷 계정·실메일·운영 공개·사람 직접 UI PASS로 표시하지 않는다. OPS 독립 보안 리뷰·TESTER QA·designer UX06/07 직접 검수는 동일 고정 후보의 후속이다.

## SAR-PUBLIC-MESSAGES-001-TESTER 독립 QA

독립 QA는 후보 `09c523da8a3407288d9f5d711e1834af12bc7808`의 detached clone에서 일반 회원 HTTP, 로컬 Node/MCP, 격리 Postgres를 확인했다. DEV 자동 검사의 종료코드 0을 이 QA의 통과로 쓰지 않았다.
프로브 종료코드는 1이다. H-1 high는 익명 slow body가 인증 전에 신규 16과 정리 4 프로세스 슬롯을 점유한다. M-1 medium은 `POST /home/gates/{id}`의 deny가 본문 파싱 전에 신규 작업으로 분류된다. receipt 20000 거절 뒤에는 신규 HTTP 행 1개가 30초 만료까지 남을 수 있다.
`make lint`, `make test`, `make verify-mvp`의 종료코드는 0이다. 이 결과는 격리 fixture다. 상세는 [QA 보고서](../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER.md)다. 실메일, 공개, 노우↔다닷은 확인하지 않았다. H-1은 main 수락을 차단한다.

## SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER 독립 QA 중단 기록

고정 d089의 [독립 보고서](../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER.md)와 [시나리오](../evaluations/scenarios/SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER.md)를 연결한다. 기본 lint·unit/race·build는 exit 0이다. runtime 재빌드와 tester QA 컴파일이 실패했고 사용자 지시로 추가 검증을 중단했다. 전체 PS08–11과 정상 cleanup/snapshot의 독립 기능 판정은 미검증이다. 기존 실패·운영/vendor 제한은 유지한다.

## SAR-GOOGLE-LOGIN-001 인증 모듈

최신 사용자 Google 우선 출시 지시를 적용한다. 원 이메일 기획의 과거 수락 결과는 바꾸지 않는다.

- `google.go`: OIDC discovery와 `oauth2` code 교환·PKCE를 사용한다. `coreos/go-oidc/v3` v3.17.0이 RS256 서명·Google issuer·client audience·expiry를 검증한다. 서버는 nonce·최근 iat·검증 이메일·비어 있지 않은 sub를 추가 확인한다. 동일 issuer의 Google 표준 issuer 변형은 라이브러리가 검증하고 저장 키는 `https://accounts.google.com|sub`로 고정한다.
- `State.GoogleAttempts`: state 해시로 10분짜리 nonce·PKCE verifier·현재 세션 해시·재확인 회원을 저장한다. callback은 별도 `__Host-kl_google` Secure/HttpOnly/SameSite=Lax cookie를 대조한다. 시도는 교환 전에 transaction으로 한 번 소비하고 만료 기록은 기존 sweep에서 제거한다. raw token·OAuth 자격은 state에 저장하지 않는다.
- `identity.go/signIn`: 이메일 흐름과 Google 흐름이 회원 한도·owner 생성·세션 발급을 공유한다. Google 신원은 issuer/sub로만 매핑한다. 이메일 회원과 자동 병합하지 않는다. 재확인은 시작한 기존 세션과 같은 회원의 Google 신원만 허용한다. callback 전에 logout·세션 만료·owner 비활성이 발생하면 재확인을 거부한다.
- `member.go`: Google-only에서 이메일 폼을 숨긴다. Google 회원의 `/auth/reauth`는 Google로 연결한다. 로그인·재확인 모두 명시 계정 선택·동의와 새 nonce에 연결된 최근 발급 ID token으로 신원을 확인한다. Google 비밀번호 재입력을 보장하지 않는다. 세션 발급 뒤 완료 화면에서 홈 링크를 누르게 하여 Strict 세션 cookie를 유지한다. 기존 CSRF/rate/기간·철회·owner/agent 권한을 재사용한다. CSP는 POST 로그인 redirect에 필요한 `https://accounts.google.com`만 form-action에 추가한다.
- `cmd/relay/main.go`, `compose.yaml`: 세 Google env를 relay에 전달한다. 일부 설정·HTTP 또는 잘못된 callback URL·discovery 실패는 기동 실패다. 정확한 env·콘솔 URI와 로컬 HTTPS 조건은 [README](../../../README.md#google-가입로그인-sar-google-login-001)를 따른다.

`TestGoogleTokenVerification`은 실제 RSA/JWK로 정상·위조 서명·issuer/audience/expiry/nonce/iat/이메일 거부를 확인한다. `TestGoogleIdentityContinuityAndReauth`는 이메일 미병합·동일 신원·회원 한도·owner·세션 회전과 재확인 거부를 확인한다. `TestGoogleAttemptAndConfiguration`은 state/cookie/재사용/만료·설정 거부를 확인한다. `TestGoogleHTTP`는 격리 Postgres와 signed local provider에서 SMTP 없는 전체 인증 delta를 확인한다. 실제 Google 계정과 직접 브라우저 확인·독립 코드 리뷰는 coor 후속이다. 상세 실행은 [DEV 기록](../exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md)에 남긴다.

## Google 연결 모듈 — SAR-GOOGLE-CONNECT-001-DEV

`internal/relay/device.go`는 시작·서명 조회·동의 페이지와 agent 생성을 담당한다. `google.go`는 기존 callback 검증 후 연결 회원을 고정한다. `member_agents.go`와 `member.go`는 기존 최근 인증·Origin 검사 및 자기 홈 선택 UI를 재사용한다. `adapters/src/login.ts`는 로컬 키 생성·승인 대기·기존 complete·자동 저장을 담당한다. `private-files.ts`는 connect/text의 파일 검사를 공유하고 Windows .NET ACL로 소유자와 접근자를 검사한다. 새 패키지는 없다.

Device 단위 검사와 GoogleDeviceHTTP 합성 OAuth/DB 검사가 회원·키 결합, 명시적 동의, 만료, replay, 별도 agent를 확인한다. Node login 검사는 실제 HTTP 서명과 자동 저장·권한을 확인한다. `.gitattributes`의 LF는 Windows gofmt 및 동결 registry digest를 동일하게 유지한다. `check_public_ingress.py`는 공개 경로와 관리자 경로 분리를 assert로 검사한다. 실계정·외부 Bot 검증은 부모 QA/OPS 인계다.
## Google 연결 최소 보완 — SAR-GOOGLE-CONNECT-002-DEV

- `login.ts/beginLogin`: 상위 폴더가 없으면 mkdir recursive로 생성한다. 마지막 키 폴더는 기존 privateDirectory exclusive 생성과 Windows ACL/POSIX 검사를 유지한다. 기존 폴더 덮어쓰기를 거부한다.
- `beta.sh/render-public-config`: Access 파일 없이 기존 UUID와 member ingress로 config.public.yml0600을 새로 생성하고 검증한다. live config를 변경하지 않으며 재생성·잘못된 UUID·검증 실패를 거부한다. `test_public_config.sh`가 make lint에 연결됐다.
- 패키지 manifest·marketplace·skill·README는 Google 연결과 별도 명시적 text 승인을 안내한다. MCP 도구9개와 기본 held, trial 계약은 그대로다. SDK API 변경은 없다.
- `TestDeviceCapacityRetention`은2000번째 허용, 다음 요청 거부, 회원 연결 제외, 만료 뒤24h 직전 보존과 경계 회복을 확인한다. `TestDeviceGoogleConnection`은 기존 다른 회원의 callback이 회원·세션·연결 소유권을 바꾸지 못하는 회귀를 포함한다. 신규 Device 포화의 medium 한계는 [D12](../operations/ops-guide.md)에 남긴다.

## 자동 수신 모듈 — SAR-AUTO-RECEIVE-001-DEV

| 모듈 | 책임 | 검증 |
|---|---|---|
| `adapters/src/inbox.ts` `Inbox` | private inbox 생성·확인, fsync·rename 저장, ID 중복 제거, rename claim 기반 1회 조회, 읽음 표시 24h 정리 | `inbox.test.ts` 2–3·6 |
| `adapters/src/inbox.ts` `AutoReceiver` | process당 직렬 pull, 10s idle, 300s 상한 backoff·`retry_at`, drain, 상태·알림 결과 기록 | `inbox.test.ts` 1·4·병렬 |
| `adapters/src/text.ts` `receive(store)` | 서명 검증 뒤 persist/ACK 전에 store 호출. store 실패 시 ACK 없음 | ACK 시점 inbox 파일 확인, store를 ACK 뒤로 옮긴 변이에서 실패 |
| `adapters/src/text.ts` CLI `watch` | MCP와 같은 inbox로 상시 수신, metadata만 출력 | `inbox.test.ts` 6 |
| `adapters/src/mcp.ts` | logging capability, 연결 뒤·로그인 완료 뒤 loop 시작, status `autoReceive`, 수동 receive의 inbox 우선 | 실제 SDK stdio Client |
| `adapters/src/grok-wake.ts` `GrokWake`, `inbox.ts` wake marker·`AutoReceiver.ring` | 선택형 loopback `sendPrompt` doorbell(고정 문장·ID만), 요청마다 gateway token 읽기, retry/rejected/uncertain 구분, ID별 `.wake` 배타 claim·결과 기록, `wake-check` 읽기 전용 `listAgents` | `inbox.test.ts` 7: 거절→재시도, 수락, 중복·재시작 무재전송, 기본 off, 401 무재시도, 연결 reset uncertain, 잘못된 대상, token·text 비노출. retry 해제를 끈 변이에서 실패 |
| `grok-wake.ts` webhook 모드(`KNOWSLINK_GROK_WEBHOOK_FILE`) | 비공개 `{url,key}` 파일을 요청마다 읽음, 고정 host URL만 허용, Bearer·`X-Automation-Key` POST 8s 단일 시도, 메타데이터만 전송, `sendPrompt`와 배타, `wake-check` 무전송 파일 검사 | `grok-wake.test.ts`: 2xx·3xx·4xx·5xx·DNS 실패·timeout, 불신 URL·짧은 key·0644 파일 무전송, 비UUID ID 제거, 결과에 key·URL ID 없음. URL host 완화·redirect follow·ENOTFOUND 분류 제거 변이에서 실패 |

`inbox.test.ts`는 번들 `dist/plugin.js`를 SDK `StdioClientTransport`로 실행하고 로컬 relay 대역을 사용한다. 실제 운영 relay·Grok Bot 호출은 포함하지 않는다.

## 회원 웹 상용 UI — SAR-SERVICE-UI-001-DEV

디자인 정본은 D04 [service-member-ui.md](mockups/service-member-ui.md)다. 서버 상태·권한·rate·CSP는 바꾸지 않았다.

| 파일·모듈 | 책임 | 자동 검사 |
|---|---|---|
| `member_pages.go` `memberStyle` | 회원 웹 디자인 토큰의 단일 정의(D04 8절 색·글꼴·간격·버튼·배지·포커스). gate 화면(`http.go` `page`)과 OAuth 동의 화면(`remote_auth.go`)이 같은 상수를 쓴다 | `TestMemberUIFixtures`: 외부 자원(`src=`·`<link`·`url(`·`<script`·`@import`) 0건 |
| `member_pages.go` `badge`·`badges` | 상태를 글자·기호·색으로 함께 표시한다. 없는 상태(온라인·최근 접속)는 만들지 않는다 | 같은 검사의 전 상태 렌더 |
| `member_pages.go` `memberPages` | 로그인·확인 코드·홈·연결 관리·기기 승인·receipt·거부 화면. 홈은 헤더·확인 필요 띠·섹션 이동·에이전트(접힘 행·위험 구역)·관계·승인 요청·연결 확인·계정 순서다 | `TestMemberPagesShowNextSteps`, integration 문구 검사 |
| `member.go` `homePage` | 관계 정렬(받은 초대→보낸 초대→연결됨→종료), 활성 에이전트 수, 받은 초대·대기 승인 요청 건수 | integration `TestPublicAgentHTTP`·`TestEmailIdentity` |
| `member_agents.go` `memberAgent.ActiveKeys`, `memberPair.Mine/Peer` | 접힌 행의 활성 키 수, 관계 행의 내 에이전트→상대 방향 | 위 검사 |
| `testdata/member_forms.golden` | 재설계 전 템플릿에서 만든 form 계약(method·action·hidden 값·필드 이름·button 라우팅). 의도한 차이는 홈 헤더 로그아웃 form 1개 추가뿐이다. SAR-FRONTEND-POLISH-001-DEV에서 선택할 에이전트가 없을 때 빈 select만 남던 홈의 관계 초대·요청 상태 form 2개를 숨겼다(`TestMemberInviteFormNeedsConnectedAgent`) | `TestMemberUIFixtures` |

시각 검수 자료: `KNOWSLINK_UI_FIXTURE_DIR=<빈 폴더> go test ./internal/relay -run TestMemberUIFixtures`가 합성 데이터의 29개 화면 HTML을 쓴다. 실제 계정·비밀값은 없다. 방법은 [실행 기록](../exec-plans/phases/SAR-SERVICE-UI-001-DEV.md)을 따른다.
FullOps DESIGN 규칙은 Go 문자열 안의 CSS를 검사하지 않는다. 대신 위 렌더 검사와 Orca 내장 브라우저의 1440·390·320 폭 측정을 사용했다.

## 캘린더 연동 모듈 — SAR-SCHEDULE-001-DEV

- `calendar.go`: D010 설정 `scheduleScope`, `Calendar`·`Lookup` 타입, `sealToken/openToken`(AES-256-GCM), `scheduleAllowed`(수락 시 범위·친구당 횟수), `roundBusy`·`disclosureJSON`(바깥 반올림·클립·병합·canonical), `startLookup`(승인 transaction 안의 lease), `GoogleLogin.freeBusy`(refresh→`freeBusy` POST, 400/401은 `errReconnect`), `Service.runLookup`(lock 밖 최대 3회·8초, 결과 1회 기록), `lookupVisible`, `sweepLookups`, `disconnectCalendar`, `Service.revoke`, `/home/calendar/{connect,disconnect}`.
- `google.go`: `ConfigureGoogle`에 `calendarKey` 인자. `googleStart(reauth, calendar)`가 calendar attempt에 freebusy scope·offline·include_granted_scopes를 더한다. callback은 scope·refresh token을 확인해 봉인 저장하고 ‘캘린더 연결 완료’를 보인다. `GoogleLogin`의 `freeBusyURL`·`revokeURL`은 합성 제공자가 바꾼다.
- `store.go`: `State.Calendars`, `Gate.Lookup`, `State.CalendarReady`(런타임 플래그), `policy("schedule.query")` 문구, `sweep→sweepLookups`, `ingest`의 `scheduleAllowed`와 `done` 결과의 disclosure 대조, 결과 수락 뒤 disclosure 삭제.
- `http.go`: `GET /v1/gates/{id}`의 `lookupVisible`·`lookup`, `gate-consume`의 `lookup`·`disclosure`, `gateView`의 `Calendar` 상태와 문구 선택, `gateDecision` approve의 `startLookup`→commit→`runLookup`.
- `protocol.go`: `relay.result`의 선택 `disclosure`와 `disclosureShape`. `gate_summary.go`: ‘공개되는 내용’ 행과 `calendarGateScope`. `human_http.go`/`human_pages.go`/`member_pages.go`: 설정 ‘캘린더 연결’ 절, gate 화면 분기, 연결 완료 복귀, `?n=calendar-off` 안내. `i18n_en_*.go`: 영어판.
- `adapters/src/core.ts`: consume 응답의 disclosure를 `done` 결과로 그대로 보내고 `lookup=failed`는 `failed`. `synthetic.ts`: D010 유효 window.
- 검사: `calendar_test.go`(반올림·봉인·범위·결과 shape·lease sweep), `calendar_provider_test.go`(합성 Google: code/refresh/revoke/freebusy, 비밀 필드 포함 응답, 철회·장애 토글), `calendar_integration_test.go`(SC-01~08 실제 Postgres), `calendar_browser_fixture_test.go`(태그 `integration && calendarfixture`, tester/designer 격리 fixture).
- SAR-SCHEDULE-001-LOWS-DEV 보정: `sealToken/openToken`은 member ID를 GCM additional data로 묶어 다른 회원 행으로 옮긴 `Sealed`가 열리지 않는다. `Service.revoke(member, sealed)`는 봉인을 연 뒤 `revokeToken`을 부르고, `revokeToken`은 상태 코드와 관계없이 응답 body를 닫는다. callback은 저장하지 않는 refresh token(scope 부족·다른 계정·만료 등)을 revoke한다. 단, 회원 자신의 연결된 계정 token은 같은 grant라 revoke하지 않는다. `owner-revoke`는 소유자 비활성화 전에 회원을 찾아 `disconnectCalendar` 후 lock 밖에서 revoke한다(`revokedOwner`의 비공개 필드는 응답에 나가지 않는다). `ingest`의 `done` disclosure 대조는 idempotency 분기 뒤라 같은 key·같은 내용 재전송은 기존 receipt를 받는다. 검사: `TestSealedTokenRoundTripAndTamper`, `TestRevokeClosesBodyOnEveryStatus`, `TestCalendarConnection`(다른 계정 revoke), `TestCalendarApprovalRunsOneLookup`(done 재전송), `TestCalendarRevocationAndFailureClasses/owner_revoke_disconnects_and_revokes`.
