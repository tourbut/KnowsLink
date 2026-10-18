---
title: 이름 중심 표시·친구 코드·친구 요청 안내 독립 QA
status: review
updated: 2026-10-19
owner: tester
tasks: [SAR-UX-ISSUES-002-NAMES-TESTER]
summary: "고정 b2e683a의 U2·U4·U7과 1단위 조회 회귀를 격리 fixture에서 독립 확인했다. critical·high 없음, low 결함 1건과 designer 판단 관찰 7건"
---

# SAR-UX-ISSUES-002-NAMES-TESTER 독립 QA

대상은 고정 SHA `b2e683ae677c78cd7bdc47c75820bcb191b94838`(fullops/dev, 제품 `516bf62`)다. 격리 fixture는 `http://127.0.0.1:18930`이다. 서버가 이 SHA인지는 manifest·coor 지시를 따랐고 tester가 독립 확인하지는 않았다. 제품 코드는 고치지 않았다. 테스트 레벨은 lite다. DEV의 `make lint`·`make test`·integration(원격 Linux, EXIT=0)은 같은 SHA 결과를 재사용했고 다시 실행하지 않았다(Windows에 make 없음). 새 독립 실행으로 바꾸어 적지 않는다.

## 요약

- **critical·high 결함 없음.** 새 결함은 low 1건이다. D-NM-01: 이름 변경이 제어·방향 문자 때문에 거부될 때 오류 문구가 이유를 설명하지 않고, 보이지 않는 문자를 입력란에 그대로 되돌린다.
- U2·U4·U7과 1단위 회귀의 합격 기준은 모두 충족했다. 이름·종류·구분값이 1440·390·320px 모든 폭에서 가로 넘침·잘림 없이 보인다. 친구 요청 응답은 없는·유효·대소문자 변경·공백 패딩·이미 친구인 코드에서 같다.
- 1단위 판정 기준은 낮추지 않았다. D-RT-01의 닫힌 select 검사(열린 전체 문자열이 닫힌 폭에 들어가는가)를 radio 계약으로 바꿨다. 새 검사는 각 라벨의 이름·종류·구분값이 줄바꿈되어 모두 보이는지, 가로 넘침·잘림이 없는지, 44px 이상인지 본다.
- designer 판단 자료 7건은 관찰만 남겼다. 요청자 기록의 ‘거절’ 표시(DEV 판단 1)와 새 대화 select의 40자 이름 잘림(DEV 판단 2)은 둘 다 재현됐다.

## 검사 결과

PASS/FAIL/not_run은 실행한 파일 기준이다. 상세는 JSON 결과와 로그에 있다.

| 영역 | 파일 | 결과 |
|---|---|---|
| HTTP 읽기 전용 (A/B/C) | [checks.py](checks.py) → [http-result.json](http-result.json) | 19 PASS |
| 비가역 변경 | [mutate.py](mutate.py) → [mutate-result.json](mutate-result.json) | 10 PASS, 2 FAIL(아래 ‘첫 실행 실패 두 건’) |
| 후속 재판정 | [followup.py](followup.py) → [followup-result.json](followup-result.json) | 4 PASS, 1 FAIL(D-NM-01) |
| 정리 | [cleanup.py](cleanup.py) → [cleanup-result.json](cleanup-result.json) | 보낸 요청 10→1 취소, B의 대기 요청 유지 |
| Chrome 155 레이아웃 | [browser.mjs](browser.mjs) layout → [browser-chrome-layout.json](browser-chrome-layout.json) | 49 PASS |
| Edge 155 레이아웃 | 같은 파일 → [browser-edge-layout.json](browser-edge-layout.json) | 49 PASS |
| Chrome 키보드·접근성 트리·대비 | keys → [browser-chrome-keys.json](browser-chrome-keys.json) | 27 PASS |
| Edge 키보드·접근성 트리·대비 | keys → [browser-edge-keys.json](browser-edge-keys.json) | 27 PASS |
| Chrome 변경 뒤 화면 | after → [browser-chrome-after.json](browser-chrome-after.json) | 22 PASS |

브라우저는 headless Chrome 155.0.8059.80과 Edge 155.0.4283.62이다. 키 입력은 CDP `Input.dispatchKeyEvent`, viewport는 비모바일 `Emulation.setDeviceMetricsOverride`(1440×900, 390×844, 320×640)이고 `innerWidth`를 측정했다. 페이지 요청은 모두 fixture origin의 GET이다. 예외는 keys의 친구 요청 POST 1건(없는 코드)이다.

### U2 이름 먼저 + 종류 + 구분값

| 조건 | 근거 | 결과 |
|---|---|---|
| 조회 radio·내 에이전트·기존 연결·새 대화·친구 상세가 `이름 / 종류 / 구분값` 순서 | `name_first_kind_discriminator_on_every_surface` | PASS |
| 같은 39자 이름 둘이 모든 화면에서 서로 다른 라벨 | `colliding_names_are_distinguishable_on_each_surface` | PASS |
| 비친구는 ‘확인할 수 없는 에이전트’, ID는 상세 안에만 | `non_friend_agent_shows_unverifiable_label_and_only_allowed_id` | PASS |
| 타 계정·비친구에게 A의 에이전트 이름 비노출, 없는 ID와 타인 ID 응답 동일 | 같은 check, `error_pages_use_alert_and_do_not_leak` | PASS |
| 40자 한글·이모지 전체 표시, 말줄임 없음, 3폭 넘침 0 | `forty_char_names_*`, `layout-*` | PASS |
| 이름 변경 뒤 ID·링크·관계·세대 동일, 재시도로 에이전트가 늘지 않음 | `agent_rename_valid_40_chars_*` | PASS |
| 빈 이름·공백·BEL·NUL·방향 문자·줄바꿈·41자 거절, 이름 불변, 타 계정 변경 거절 | `rename_with_unsafe_name_is_refused_*` | PASS |
| 내 친구 코드: 보이는 라벨, 읽기 전용, 폼 밖, Ctrl+A 전체 선택, Ctrl+C 전체 복사 | `friend_code_*`, `ctrl-a-*`, `ctrl-c-*` | PASS |
| 코드 교체 뒤 옛 코드 비노출, 기존 친구·받은 요청 유지, 타인 코드 불변 | `code_rotation_*` | PASS |
| 옛 코드로 새 요청 불가(없는 코드와 같은 응답), 새 코드는 가능 | `old_code_after_cancel_*` | PASS |

### U4 친구 요청

| 조건 | 근거 | 결과 |
|---|---|---|
| 제출 응답 동일(없는·유효 대기 중복·공백 패딩·대소문자 변경·새 유효 코드·이미 친구) | `request_submit_same_status_and_text_*` | PASS |
| 제어·방향 문자·500/2000자 코드도 같은 일반 응답, 되돌려 쓰지 않음 | `odd_friend_codes_*` | PASS |
| 반복 요청이 24h 기한을 늘리지 않고 A에게 중복 요청이 생기지 않음 | 같은 check | PASS |
| 빈 코드·공백 코드·동의 누락 422, 입력 보존, alert·`aria-describedby`·`aria-invalid` | `friend_request_invalid_inputs_422_*` | PASS |
| 받은 대기: 요청자 이름·범위·절대 만료·수락/거절(이름+ID 이름) | `pending_received_and_sent_states` | PASS |
| 보낸 대기: 상대 정보 없음, 취소만 가능, 유효·없는 코드 행 모양 동일 | `sent_pending_view_reveals_nothing_*` | PASS |
| 거절 뒤 구세대 수락 403, 재요청은 새 세대, 중복 수락은 한 번만 반영 | `deny_*`, `re_request_after_deny_*` | PASS |
| agent 0개 C: 요청 가능, 빈 선택기 없음, 수락 뒤 ‘첫 에이전트를 추가하세요’, A 쪽은 ‘연결된 에이전트가 없습니다’와 대화 링크 없음 | `zero_agent_*`, `friends_page_request_form_*` | PASS |
| GET·쿼리 위조·뒤로가기로 관계 생성·수락 없음 | `get_requests_never_create_*`, `history-back-*` | PASS |
| 보낸 요청이 많아도(14건 추가) 새 없는 코드의 응답 동일 | `sender_with_many_pending_*` | PASS |

### U7 보조 화면

| 조건 | 근거 | 결과 |
|---|---|---|
| 고유 title, `lang=ko`, h1 하나, main, 건너뛰기 링크(실제 Enter로 `#main`) | `landmarks_*`, `skip-link-*` | PASS |
| 기존 연결 화면은 현재 위치를 거짓 표시하지 않음, title이 ‘내 에이전트’와 다름 | `legacy_title_differs_*` | PASS |
| 정확한 내부 복귀(조회→내 에이전트·관리 화면, 친구 상세→친구 목록, 새 대화→대화 목록), 외부 return 파라미터 무시 | `return_links_*` | PASS |
| 결과는 status, 오류는 alert | 키보드 제출 결과, 422·403 화면 | PASS |
| 반복 버튼의 접근성 이름에 대상+구분값(`에이전트 이름 저장: 이름 / agent_a2`, `친구 요청 수락: 이름 / ID`) | `repeated_buttons_*`, AX 트리 | PASS |
| Tab 이동 전체 순환·함정 없음, 모든 정지점에 접근성 이름과 보이는 focus(3px) | `*-tab-walk-*` (계정·친구 A·친구 C·조회) | PASS (Chrome·Edge) |
| 44px 이상 동작(체크박스·radio는 라벨 영역 기준), 3폭 가로 넘침 0 | `layout-*` 48건/브라우저 | PASS |
| 대비: 본문 14.70, 보조 6.19, 링크·버튼 7.96, 입력 테두리 3.79(≥4.5·≥3) | `contrast-*` | PASS |

### 1단위 회귀

조회 radio는 이름/종류/구분값 라벨, 첫 화면 미선택, `?agent=`로 선택, 422·403에서 선택과 입력 보존, 없는·다른 에이전트 ID 응답 동일을 확인했다(`receipt_radio_selector_*`). 키보드로 radio 진입·방향키 선택·ID 입력·Enter 제출 뒤 선택과 ID가 남는다. 익명 429는 31번째 요청에서 나오고 ‘조회할 내 에이전트가 없습니다’ 문구가 없다(`anonymous_429_*`).

## 결함과 판단 자료

### D-NM-01 (low) 이름 거부 사유가 보이지 않는다 — FAIL

- 재현: A로 `POST /home/agent-name` `agent=agent_a2&name=ab‮cd`(방향 재정의), `ab\x07cd`(BEL), `ab⁦cd`, `ab\x00cd`, `ab\ncd`.
- 기대: U2는 제어·방향 제어 문자를 금지한다. 거부될 때 오류가 그 규칙을 알리고 보이지 않는 문자를 입력란에 되돌리지 않는다(조회 화면은 같은 문자를 버린다).
- 실제: 422와 `이름은 1–40자로 입력하세요`만 나온다. 입력란 `value`에 BEL·U+202E·U+2066이 그대로 들어가고 NUL은 U+FFFD로 바뀐다. 5자 이름이 왜 거부되는지 알 수 없다. 이름은 바뀌지 않았다(안전).
- 증거: [followup-result.json](followup-result.json) `echo_failures`, `rename_error_messages`. 고정 SHA b2e683a, fixture 18930. fixture가 `/home/profile`을 막아 ‘내 이름’의 같은 경로는 확인하지 못했다.
- 영향: 사용성. 데이터 노출·권한 문제 없음. 다른 사용자에게는 보이지 않는다.

### 첫 실행 실패 두 건 (mutate.py, 원본 보존)

- `request_validation_422_keeps_input_and_creates_nothing` FAIL: 공백만 입력한 코드의 422에서 `value=""`를 기대했다. 제품은 입력한 공백을 보존한다. 계약(입력 보존)에 맞는 동작이라 기대가 틀렸다. [followup.py](followup.py)가 계약대로 다시 판정했고 PASS다.
- `agent_rename_rejects_unsafe_names_without_change` FAIL: 첫 번째 거부 이름(BEL)에서 ‘되돌려 쓰지 않음’ 단언이 실패해 멈췄다. 제품 동작이 맞는 부분(이름 불변)은 followup `rename_with_unsafe_name_is_refused_*`가 전부 PASS로 확인했다. 되돌려 쓰는 부분이 D-NM-01이다. 검사를 완화하지 않았다.

### 관찰 (결함 아님, designer·coor 판단)

| ID | 관찰 | 근거 |
|---|---|---|
| O-NM-01 | **요청자에게 ‘거절’ 표시(DEV 판단 1).** A가 C의 요청을 거절한 뒤 C의 기록에 `친구 요청 기록 / — 거절 … 새 친구 요청과 수락이 필요합니다`가 보인다. 없는 코드 요청은 만료까지 ‘요청 대기’로 남는다. 결정 전에는 구분할 수 없지만 거절되면 코드가 실제 회원이었음이 드러난다. | mutate-result `requester_view_after_deny`, `requester_sees_declined_word: true` |
| O-NM-02 | **새 대화 select 잘림(DEV 판단 2).** 가용 폭 대비 선택 옵션 전체 폭: 1440에서 911/1118px(40자 이모지 상대는 1167/1118px 넘침), 390에서 911/308px, 320에서 911/223px. 1440 외에서는 종류·구분값이 닫힌 상태에 보이지 않는다. 열린 목록과 줄바꿈 라벨은 아니다. | browser-chrome-layout `selects` |
| O-NM-03 | 보낸 요청 행이 서로 같다(‘보낸 친구 요청 / 요청 대기 / 기한 / 요청 취소’). 만료 시각만 다르다. 접근성 이름에는 요청 ID가 있다. 상대 정보 비공개 규칙 때문이지만 사용자는 어느 코드로 보낸 요청인지 알 수 없다. | 같은 C 화면에서 10행 |
| O-NM-04 | 익명 429 화면: 빈 `radiogroup`(범례만, 선택지 없음)과 `<a href="">로그인 확인 후 다시 조회</a>`, 저장된 입력이 없는 익명에게 ‘저장된 입력은 … 복원됩니다’가 보인다. 1단위의 같은 관찰이 남아 있다. 오안내 문구는 사라졌다. | `anonymous_429_*` evidence |
| O-NM-05 | 값이 입력 안에서 가로 스크롤된다: 40자 내 이름·39자 에이전트 이름 입력(390·320), 조회 ID 입력(320의 36자 UUID, 64자 입력). DEV 보고와 같은 기존 동작이다. 친구 코드는 `textarea rows=2`라 320에서도 두 줄에 모두 보이고 잘리지 않는다(D04 O-06의 잘림은 해소). | browser-chrome-layout `inputScroll` |
| O-NM-06 | 회원당 `/home` 요청 한도가 모든 화면에 적용된다. 같은 회원이 1분에 41번째 요청(GET, 40개째까지 200)에서 429 ‘요청이 많습니다. … 이후 다시 시도하세요’를 받았다. 이 QA는 429를 기다리며 진행했다(HTTP 읽기 검사 16회, 비가역 검사 8회 대기). 사람의 일반 사용으로는 드물지만 구현과 문구 범위가 조회 한도와 다르다. | 측정 41번째 429, 응답 재시도 시각 약 1분 |
| O-NM-07 | 요청 제출 뒤 포커스는 `body`에 남고 status는 페이지 읽기 순서에 있다. 낭독 도구가 새 페이지의 status를 읽는지는 확인하지 못했다. 입력란은 오류 때 `autofocus`로 이동한다. | `cRequestResult.active: BODY` |

## 캡처 (designer 시각 검수용)

Chrome 캡처 `chrome-<화면>-<폭>.png`, 폭 1440·390·320(전체 높이). `-after`가 없는 파일은 비가역 변경 전 상태이고 `-after`는 코드 교체·수락 뒤 상태다.

| 항목 | 파일 |
|---|---|
| 친구 빈 상태(C, 친구 없음+보낸 대기) | `chrome-friends-c-*` |
| 친구 추가·내 코드(A, 받은 요청+내 코드+요청 폼) | `chrome-friends-a-*` |
| 내 계정 코드 | `chrome-account-a-*`, 교체 뒤 `chrome-account-a-after-*` |
| 조회 radio(첫·거절) | `chrome-receipts-a-*`, `chrome-receipts-refusal-a-*` |
| 내 에이전트(상세 펼침, A·C) | `chrome-agents-a-*`, `chrome-agents-c-*` |
| 연결·기존 연결(상세 펼침) | `chrome-legacy-a-*` |
| 친구 상세(A→B 이모지, B→A 동명 둘) | `chrome-friend-ab-a-*`, `chrome-friend-ab-b-*` |
| 새 대화 선택기 | `chrome-new-self-a-*`, `chrome-new-peer-a-*`, `chrome-new-c-*` |
| 요청 제출 결과(390·1440) | `chrome-c-request-result-*` |
| 수락 뒤 친구 상세(0 에이전트 C) | `chrome-friend-ac-a-after-*`, `chrome-friend-ac-c-after-*`, `chrome-friends-a-after-*`, `chrome-friends-c-after-*` |

`edge-*.png`는 Edge 155의 같은 화면이다. 승인 화면 캡처는 없다(not_run).

## not_run

- **승인(gate) 화면.** fixture가 허용하지 않는다(`/home/gate` 403). 이름 표시는 DEV Go 테스트(`TestGateNamesEndpointsWithinAuthority`)를 근거로 하며 이 QA는 독립 확인하지 못했다. 기존 연결 화면의 ‘승인 요청’ 구역은 요청 없음 상태만 봤다.
- **조회 성공 화면.** fixture에 receipt 데이터가 없다. 거절·형식 오류 화면만 확인했다.
- **구분값 충돌 확장.** fixture 에이전트는 구분값이 서로 달라 충돌 확장이 일어나지 않는다. DEV Go 테스트(`TestReceiptSelectorDistinguishesCollidingNames`, `TestLegacyAgentLabelsWidenCollidingDiscriminators`)를 근거로 하며 이 QA는 독립 확인하지 못했다.
- **‘내 이름’ 변경.** fixture가 `/home/profile` POST를 막는다.
- **24h 경합·시간 경과 만료·차단·친구 해제.** 시간과 비가역 동작이 필요해 하지 않았다.
- **스크린리더 음성, 실제 브라우저 200% 확대, Firefox·Safari·모바일 실기.** 하지 않았다. CDP 비모바일 viewport로 3폭만 확인했다. 전체 접근성 완료를 주장하지 않는다.
- **실제 인증·메일·운영 화면, 메시지 발송.** 범위 밖이다.
- **등록 `make lint`·`make test`·integration.** DEV 결과 재사용, 재실행 없음.

## 환경과 한계

- fixture는 비가역이다. 이 QA 뒤 상태: A↔B 친구, A↔C 친구(C는 agent 0개), C→B 대기 요청 1건, A의 친구 코드 교체. 이름은 원래대로 되돌렸다. 재검증하려면 coor가 fixture를 다시 준비해야 한다.
- fixture의 허용 목록 때문에 POST 전용 경로의 GET 위조(`/home/friends/decision` 등)는 제품이 아니라 fixture가 403을 돌려준다. GET 불변 증거는 허용 화면의 쿼리 위조와 뒤로가기에 한정된다.
- 변경 전 캡처는 키보드 검사(없는 코드 요청 1건)보다 앞서 찍었다. C의 보낸 요청이 한 행뿐이다.
- 브라우저 검사의 canvas 폭 측정은 추정이다(select). 레이아웃·키보드 값은 실제 DOM·CDP 값이다.
