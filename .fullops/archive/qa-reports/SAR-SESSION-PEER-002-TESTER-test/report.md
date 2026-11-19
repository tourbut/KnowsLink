---
title: SAR-SESSION-PEER-002-TESTER — 세션 메시지 후속 8건 후보 da97dba 독립 재확인
status: draft
updated: 2026-10-17
owner: tester
tasks: [SAR-SESSION-PEER-002-TESTER]
summary: "da97dba의 후속 항목 1·2·3·5·6·7·8, 새 API 권한 경계(제삼자 contacts 거부, 비수신자 seen 거부), V01 왕복, 미연결·subagent deny가 모두 통과했다. 결함 0건. 관찰 3건(low/info)."
---

# SAR-SESSION-PEER-002-TESTER — 독립 재확인

독립 동작 QA(lite)다. 사용자가 이 후보의 독립 코드 리뷰를 생략해서 이 기록이 유일한 독립 검증이다. 그래서 새 서버 계약의 권한 경계는 직접 HTTP로 확인했다. 제품 코드와 사용자 전역 `~/.claude`는 수정하지 않았다. 실제 회원·운영 relay·실제 상대에게 보내지 않았다. 항목 4(지시서에 없음)는 번호를 건너뛴 그대로다.

## 고정 후보와 환경

| 항목 | 값 |
|---|---|
| 제품 SHA | `da97dbaa86f777fd1f80e6f856c7282303a45790` (`fullops/dev`, 코드 9b7b222·3500e1b). `git archive`로 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/session-peer-002-qa/src`에 추출 |
| 호스트 | Windows 11, Claude Code 2.1.302, Node 22.22.2, Docker Desktop, Go |
| 설치 | `npm ci`·`npm run build --prefix adapters`·`npm test --prefix adapters` exit 0 ([npm-ci.log](npm-ci.log), [npm-build.log](npm-build.log), [npm-test.log](npm-test.log)). `package_plugin.py --host claude --install-dir … --connect --verify` A·B exit 0 ([pkg-a.log](pkg-a.log), [pkg-b.log](pkg-b.log)). 설치본의 `config/local.mcp.json`은 RELAY_URL과 에이전트 폴더만 격리 값으로 바꿨다 |
| relay | **같은 SHA에서 빌드한** 새 Compose 프로젝트 `knowslink-sp002-qa`, `http://127.0.0.1:18261`, `KNOWSLINK_SYNTHETIC_SIGNUP=1`. 메일 sink는 relay 네트워크 사이드카. DEV 프로젝트 `knowslink-sp002-dev`는 쓰지 않았다 |
| 회원 | `*@sessionpeer002.invalid` A·B(친구 수락)·C(친구 아님), 각 `node-local` agent 1개 ([fixture.py](fixture.py), [fixture-agents.json](fixture-agents.json)) |
| 세션 | Orca 터미널의 `claude --plugin-dir <설치>/knowslink`, `auto mode on`, 합성 소유자 승인은 tester가 Orca `terminal send`로 눌렀다 ([launch.sh](launch.sh), [t.sh](t.sh)) |

## 판정 요약

| # | 검사 | 판정 | 핵심 근거 |
|---|---|---|---|
| 1 | 승인 창 길이 한도 | PASS | 본문 1849자 → 창이 `…TAILEND」 을 그대로 보낼까요?`까지 끝까지 표시. 1850자·1879자 → 대화상자 없이 deny, "약 1849자 이하로 나눠" 안내 |
| 2 | 승인 창 중 연결 상실 | PASS (`--take`) | 창을 연 채 CLI watcher `--take` → 승인 → "승인하는 동안 연결이 해제돼 발송을 막았습니다", B 수신 `empty`. `off` 변형은 직접 못 했다(아래) |
| 3 | 비가시 문자 | PASS | `visible⟨U+FE0F⟩hidden⟨U+200B⟩end⟨U+202E⟩rev`, `보이지 않는 문자 3개` 표시 |
| 5 | 재연결 대기 수 | PASS | B가 2건을 받은 뒤 `off`→`on`: 연결 문구 `대기 2건`, 상태줄 `대기 2`, `status` `대기 2건`, inbox 미확인 `.json` 2개 |
| 6 | 친구·agent 이름 | PASS | 승인 창·toast·receive 결과 모두 표시 |
| 7 | relay 통신 불가 | PASS | 정지 약 45초 뒤 두 세션 상태줄 `통신 불가 · 재시도 중`. 도구 결과는 `relay unreachable`와 "키·친구 문제가 아님" 안내. 재시작 21초 뒤 상태줄 `연결` 복귀 |
| 8 | 세션 수신 확인 | PASS | `transport:"delivered"`는 수신 세션이 `receive`로 표시하기 전에 이미 있고, `session_received_at`은 표시 뒤에야 별도로 붙는다 |
| 경계 | 제삼자 contacts·비수신자 seen | PASS | 직접 HTTP 20개 검사(+info 3줄) 전부 통과 ([sec_http.mjs](sec_http.mjs), [sec-http.log](sec-http.log)) |
| 회귀 | V01 왕복, 미연결 deny, subagent deny | PASS | 아래 |

결함 0건. 미해결 critical/high 없음.

## 보안 경계 (직접 HTTP)

실행: `node sec_http.mjs <dist> fixture-agents.json`, 서명된 SDK 대신 각 agent의 bearer 자격으로 `fetch`를 직접 호출했다. 모든 응답 본문은 [sec-http.log](sec-http.log)에 있다.

**GET /v1/contacts/{agent}**

| 호출 | 결과 |
|---|---|
| 친구 B → A | 200. 키는 정확히 `agent`·`client`·`friend`·`name` 4개. 이메일·키·kid·자격 없음 |
| 친구 A → B | 200 |
| **제삼자 C → A** | **403 `sender_not_allowed`**, 본문에 정보 없음 |
| **제삼자 C → B** | **403** |
| C → C(자기), A → A(자기) | 403. 자기 agent 조회는 거부된다(DEV 보고와 같다. 지시서의 '자기연결'은 별도 연결 관계를 뜻한다고 읽었다) |
| C → 존재하지 않는 agent | C → A와 상태·본문이 같다(존재 여부 노출 없음) |
| 자격 없음 | 401 `invalid_auth` |
| `GET /v1/contacts` (C의 목록) | `[]`, A·B 미포함 |
| 인코딩된 슬래시 우회 | 403 |
| 친구 해제 뒤 B ↔ A | 양방향 403 ([contacts-after-remove.log](contacts-after-remove.log)). 재수락 뒤 200 ([contacts-after-readd.log](contacts-after-readd.log)) |

**POST /v1/text/seen**

| 호출 | 결과 |
|---|---|
| 수신자 B, 아직 `queued`(전달 전)인 메시지 | 403 |
| **발신자 A**, 전달된 메시지 | **403** |
| **제삼자 C**, 전달된 메시지·다른 메시지 | **403** |
| 수신자 B, 존재하지 않는 id | 403 |
| 자격 없음 | 401 |
| 수신자 B, 전달된 메시지 | 200, `session_received_at` |
| 수신자 B, 재호출 | 200, 시각 불변(첫 시각 유지) |
| 발신자 A의 receipt | `transport:"delivered"`와 별도 `session_received_at`. seen 전 receipt에는 필드 없음 |
| 제삼자 C의 receipt | 403 |

## 항목별 관찰

### 1. 승인 창 길이

- 한도는 본문이 아니라 질문 전체 2000자다. 이 세션의 질문 부가 문구(친구 이름·agent 이름·ID 포함)는 151자라 본문 한도는 1849자다. deny 문구가 이 값(`약 1849자 이하`)을 정확히 계산한다.
- 1849자: [ev-dialog-1849-a.txt](ev-dialog-1849-a.txt)에서 `w001`부터 `--TAILEND」 을 그대로 보낼까요?`까지 줄이어 보였다. 취소해 발송 0.
- 1850자: [ev-deny-1850-a.txt](ev-deny-1850-a.txt). 1879자: [ev-deny-1879-a.txt](ev-deny-1879-a.txt). 둘 다 대화상자가 열리지 않았다.
- 관찰 O-T2-01: DEV 보고는 영문 본문 1879자가 한도라고 적었다. 이 환경에서는 같은 1879자가 deny됐다. 이름이 질문에 들어가 부가 문구가 길기 때문으로 보인다. 제품 동작은 맞다(안내 값이 실제 한도). 문서의 수치만 어긋난다.

### 2. 승인 창 중 연결 상실

- A의 승인 창을 연 채([ev-release-during-dialog-before.txt](ev-release-during-dialog-before.txt)) A의 에이전트 폴더에 CLI watcher를 `--session cli-take --take`로 시작했다([ev-watch-take.log](ev-watch-take.log): `knowslink_session_bound … cli-take`). 락이 바뀐 것을 확인하고 `보내기`를 눌렀다.
- 결과: `Send blocked. The session's KnowsLink connection dropped during approval, so nothing was delivered.` ([ev-release-during-dialog-after.txt](ev-release-during-dialog-after.txt)). A 상태줄 표시가 사라졌다. B `receive` → `empty`.
- 한계: `/knowslink-session off` 변형은 승인 창이 입력을 막아서 같은 세션에서 직접 만들 수 없었다. 제품 코드는 둘을 같은 연결 비교로 처리한다고 DEV가 보고했고 단위 시험이 있다. 이 QA는 `--take`만 직접 확인했다.

### 3. 비가시 문자

- [ev-dialog-invisible-a.txt](ev-dialog-invisible-a.txt): `본문 전체 22자 · 보이지 않는 문자 3개를 ⟨U+…⟩로 표시 「visible⟨U+FE0F⟩hidden⟨U+200B⟩end⟨U+202E⟩rev」`. 파일 내용(`xxd`로 모델이 확인)과 일치한다. 이 승인 창은 2번 시험에서 취소 대신 연결 상실 deny로 끝났다.

### 5. 재연결 대기 수

[ev-reconnect-pending-b.txt](ev-reconnect-pending-b.txt): B가 수신한 2건(inbox `.json` 2개) → `off`→`on` → `대기 2건`, 상태줄 `대기 2`, `status` `대기 2건`. 이전 QA(O-02)는 0이었다. 연결 직후 relay 쪽에 이미 대기 중이던 메시지(처음 B 연결 때)는 inbox 저장 전이라 `대기 0건`으로 시작했고 곧 `대기 1`로 올랐다. inbox 미확인 수와 일치한다.

### 6. 이름 표시

| 위치 | 증거 |
|---|---|
| 발송 승인 창 | `→ 친구 합성 소유자 B · agent QA 세션 B (agent_dd3bb38614cfd813a83dcf)` ([ev-dialog-names-a.txt](ev-dialog-names-a.txt)). 답장 창은 `답장(reply_to …)` 포함 ([ev-reply-dialog-b.txt](ev-reply-dialog-b.txt)) |
| toast | `KnowsLink 외부 메시지 도착 · 발신 친구 합성 소유자 A · agent QA 세션 A (agent_22bb187ebb5358d648b6e7) · 대기 …` ([ev-toast-names-b.txt](ev-toast-names-b.txt)). 좁은 toast 폭에서 끝(`대기 N건`)이 `…`로 잘린다 |
| `knowslink_text_receive` 결과 | `"sender":{"friend":"합성 소유자 A","name":"QA 세션 A","client":"node-local"}` ([ev-receive-names-seen-b.txt](ev-receive-names-seen-b.txt), [ev-receive-b-2.txt](ev-receive-b-2.txt)) |

### 7. relay 통신 불가

- 상태줄: 두 세션 모두 `통신 불가 · 재시도 중`([ev-outage-a.txt](ev-outage-a.txt), [ev-outage-b.txt](ev-outage-b.txt)). `연결`로 남지 않았다.
- `/knowslink-session status`: `연결됨 … 대기 0건 · relay 통신 불가(자동 재시도 중)`.
- 도구 결과([ev-outage-tool-b.txt](ev-outage-tool-b.txt)): `{"state":"failed","error":"relay unreachable","received":false,"next":"Relay is unreachable or unavailable now; this is not a key or friend problem. Retry later; retry uncertain sends with the same key/content…"}`. 이전의 "key, active pair, TTL 확인"과 다르다. 발송도 같은 분류다([ev-outage-send-a.txt](ev-outage-send-a.txt)).
- 재시작(07:49:40Z) 뒤 07:50:01Z에 상태줄 `연결`([ev-outage-recovered-b.txt](ev-outage-recovered-b.txt)).
- 확인하지 못한 것: 통신 불가 전환 때의 toast 글자(상태줄·status·도구 결과로만 확인).

### 8. 세션 수신 확인

| 시점 | A의 receipt (`knowslink_text_receipt`) |
|---|---|
| B가 `receive`하기 전 (relay는 이미 전달함) | `transport:"delivered"`, `session_received_at` 없음 ([ev-receipt-before-seen-a.txt](ev-receipt-before-seen-a.txt)) |
| B `receive`(결과 `sessionReceipt:"recorded"`) 뒤 | `transport:"delivered"` + `session_received_at:"2026-10-17T07:39:29.510049Z"` ([ev-receipt-after-seen-a.txt](ev-receipt-after-seen-a.txt)) |

관찰 O-T2-02: `completion:"received"`는 세션이 보기 전에도 receipt에 나타난다(relay ACK 시점의 기존 필드). 새 필드 `session_received_at`이 별도라 판별은 되지만, 이름만 보면 '받음'이 두 곳에 있다. 사용자에게 보이는 표기에서 `completion`을 '세션 수신 확인'으로 읽지 않게 안내하는 편이 낫다(low).

## 회귀

| 흐름 | 결과 |
|---|---|
| V01 왕복 | A→B 승인 창(`새 메시지 · 본문 전체 17자`) → queued → B toast·`대기 1` → B `receive` 일치 → B→A 답장 승인 창(`답장(reply_to 01a148cc-…)`) → A `receive` `QA2 pong 1 from B`, `reply_to` 일치 ([ev-v01-reply-received-a.txt](ev-v01-reply-received-a.txt)). 답장의 만료는 원 메시지 만료에 묶여 `expired:true`로 표시된다(이전 QA와 같다) |
| 미연결 deny | `off` 뒤 `receive`·`send` 모두 `이 세션은 연결되지 않았습니다. /knowslink-session on 으로 먼저 연결하세요.` ([ev-unconnected-deny-a.txt](ev-unconnected-deny-a.txt)). 승인 창 없음 |
| subagent deny | 같은 세션 subagent의 `receive`·`send` 모두 `하위 에이전트에서는 KnowsLink 메시지를 쓸 수 없습니다.` ([ev-subagent-deny-b.txt](ev-subagent-deny-b.txt)). 이번 subagent는 백그라운드로 실행됐다 |

## 관찰

| ID | 심각도 | 내용 |
|---|---|---|
| O-T2-01 | info | DEV 보고서의 '영문 1879자 한도'는 이름이 긴 환경에서 1849자다. 동작은 맞고 안내 값이 정확하다. 보고 수치만 환경 의존이다 |
| O-T2-02 | low | receipt의 기존 `completion:"received"`가 세션 확인 전에도 나타난다. 새 `session_received_at`과 이름이 겹쳐 오독할 수 있다 |
| O-T2-03 | info | toast 폭이 좁으면 끝의 `대기 N건`이 `…`로 잘린다. 이름·ID는 보인다 |

## 미실행과 한계

`/knowslink-session off` 변형의 승인 창 중 해제(위 2번), 통신 불가 toast 글자, 한글 4096바이트 최대 본문, 1850자 초과의 다른 길이(3000자 등), 운영 relay(새 엔드포인트 없음, 범위 밖), Codex·Hermes는 실행하지 않았다. 질문 부가 문구는 호스트(Claude Code 2.1.302)의 2000자 표시 한도에 의존한다. `make lint/test` 전체와 Linux race 시험은 실행하지 않았다(DEV가 보고). 격리 relay·Windows 한 대·단일 조합이다.

## 실행한 자동 검사

| 명령 | 결과 |
|---|---|
| `npm ci --prefix adapters` / `npm run build --prefix adapters` / `npm test --prefix adapters` | exit 0 |
| `package_plugin.py --host claude … --connect --verify` ×2 | exit 0 |
| `go test -count=1 ./internal/relay -run "Contact|Seen|SessionReceived|PublicText|OwnerText"` | `ok` ([go-text-test.log](go-text-test.log)) |
| `node sec_http.mjs …` | 20 PASS, 0 FAIL ([sec-http.log](sec-http.log)) |

## 정리 상태

세션 2개 `off`·`/exit`, 터미널 닫음. CLI watcher(pid 39996) 종료. Compose 컨테이너 `knowslink-sp002-qa-{relay-1,postgres-1}`와 sink `knowslink-sp002-qa-sink`를 정지했고 삭제하지 않았다. 도구·로그는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/session-peer-002-qa/`에 남아 있다. 사용자 전역 설정은 건드리지 않았다. 작업 폴더 `work-a`·`work-b`(Claude 폴더 신뢰 대화상자에서 신뢰)와 시험용 본문 파일을 그 폴더에 남겼다.
