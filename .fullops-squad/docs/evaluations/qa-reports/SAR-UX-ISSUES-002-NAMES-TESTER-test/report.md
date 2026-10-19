---
title: 이름 중심 표시·친구 코드·친구 요청 안내 독립 QA
status: review
updated: 2026-10-19
owner: tester
tasks: [SAR-UX-ISSUES-002-NAMES-TESTER]
summary: "5fa3333 재확인: 수정 4건 PASS, 새 결함 D-NM-02(medium, 40자 이름 상대와의 새 대화 제목이 390·320px에서 넘침), critical·high 없음. 이전 b2e683a 결과는 아래에 보존"
---

# SAR-UX-ISSUES-002-NAMES-TESTER 독립 QA

## 재확인 (DEV 수정 5fa3333, 제품 97a8899)

대상은 fullops/dev `5fa3333b5948fd36792662b202915bc6243cac14`(제품 `97a8899`)다. fixture는 같은 origin `http://127.0.0.1:18930`이다. 서버가 이 SHA인지는 manifest·coor 지시를 따랐고 tester가 독립 확인하지는 않았다. fixture는 재기동 직후 상태(A↔B 친구, C→A 대기 요청 1건)에서 시작했다. 제품 코드는 고치지 않았다. 테스트 레벨 lite. DEV의 `make lint`·`make test`·integration(원격 Linux, EXIT=0, SHA 97a8899)은 같은 SHA 결과를 재사용했고 다시 실행하지 않았다. 아래 첫 절의 b2e683a 결과는 이전 attempt 기록이며 덮어쓰지 않았다(새 결과는 `-5fa3333`·`chrome-fix-*`·`edge-fix-*` 파일).

### 재확인 요약

- **critical·high 없음. 새 결함 1건(D-NM-02, medium): 40자 이름의 상대와 새 대화를 열면(B→A) 수신 에이전트 그룹 제목이 390·320px에서 페이지 밖으로 넘친다.** 수정 4건 자체는 모두 동작한다.
- 수정 4건: ① PASS(아래 not_run 한계 표기), ② PASS, ③ PASS, ④ PASS. 영향 회귀 PASS.
- D-NM-01은 해결됐다. 금지 문자 거부는 전용 문구를 보이고, 입력란에는 금지 문자를 뺀 값이 남고, 이름은 저장되지 않는다.
- 상태를 바꾸는 POST는 마지막의 요청 비노출 회귀 1건(C의 친구 요청 4회, 새 없는 코드 2개로 보낸 요청 행 +2)과 거절되는 이름 저장뿐이다. 친구 결정·코드 교체·이름 저장은 하지 않았다. designer가 같은 fixture를 쓸 수 있다.

| 수정 | 판정 | 근거 |
|---|---|---|
| ① 새 대화 radio | **PASS** (한계 2건 아래) | 자기·친구·B 시점 모두 select 없음, `fieldset role=radiogroup`+legend, 라벨 `이름 / 종류 / 구분값`, 동명 둘의 라벨이 다름. 기본 선택 발신 첫째·수신 둘째(자기)/첫째(친구). 390·320에서 라벨이 줄바꿈되어 구분값까지 보임. Tab이 그룹의 선택된 radio에 들어가고 ↓·↑로 선택이 바뀐다. Enter가 고른 값으로 `/home/conversation?from=agent_a2&to=agent_a1`을 요청한다. 접근성 트리: 그룹 이름=legend, radio 이름=라벨. 다만 D-NM-02 |
| ② 두 링크 분리 | **PASS** | 각 `<p>`에 하나씩, 간격 16px, 높이 44px, 3폭 모두. 키보드 Enter로 `요청 상태 조회`→`/home/receipts?agent=<id>#receipt-form`(해당 radio 선택), `연결·키 관리`→`/home/legacy#agents`(앵커 존재) |
| ③ D-NM-01 | **PASS** | 아래 표 |
| ④ 거절 앞 안내 | **PASS** | 수락 → 안내(`p#deny-help-<id>`) → 거절 순서. 거절 버튼의 접근성 설명이 안내 문구와 같고 수락에는 설명 없음. Tab 순서는 수락→거절이며 안내는 정지점이 아니다. 요청자(C) 화면에는 안내가 없다. 문구는 designer 결정 1과 글자까지 같다 |

#### ③ D-NM-01 재판정 (`mutate.py` 갱신 검사 + `followup.py` + 브라우저)

`agent_rename_rejects_unsafe_names_without_change`를 새 기준(제거값 복구)에 맞춰 갱신했다. ‘저장 안 됨’(이름 불변)·‘타 계정 변경 거절’ 단언은 유지했고 응답 문구·입력값 단언을 추가했다. 첫 attempt의 제어 문자 echo 단언은 그대로 두고, 입력값이 ‘금지 문자를 뺀 값’인지로 바꿨다.

| 입력 | 응답 | 문구 | 입력란 값 | 이름 변경 |
|---|---|---|---|---|
| `ab`+BEL+`cd`, NUL, RLO(U+202E), 격리(U+2066), 줄바꿈 | 422 | 금지 문자 전용 문구 | `abcd` | 없음 |
| BEL만 | 422 | 전용 문구 | 빈 값 | 없음 |
| BEL + 41자 | 422 | 전용 문구(두 조건 동시) | 제거 뒤 41자 | 없음 |
| 빈 이름·공백만·41자·이모지 41자 | 422 | 기존 ‘이름은 1–40자로 입력하세요’ | 입력한 값 그대로 | 없음 |

전용 문구: ‘이름에 줄바꿈·제어 문자·방향 제어 문자는 쓸 수 없습니다. 보이는 글자로 1–40자를 입력하세요.’ 모든 거부 화면에 `role=alert`·`aria-invalid=true`가 있고 응답에 제어·대체(U+FFFD) 문자가 없다. 브라우저에서 입력란에 `ab`+BEL+RLO+`cd`를 넣고 Enter로 제출해도 같은 결과(입력란 `abcd`)이며 요청은 POST 1건뿐이다. 다시 연 내 에이전트 화면의 이름은 그대로다(39자). 3폭 가로 넘침·44px 위반 없음.
**경로 한계(not_run)**: fixture가 `/home/profile`(내 이름)·`/home/agents`(에이전트 추가) POST를 막는다(둘 다 403 ‘outside names fixture scope’). device·remote 동의 화면은 열 수 없다. 이 네 경로의 제품 문구는 이 QA가 독립 확인하지 못했다. 같은 `nameRefusal` 함수 사용은 DEV 보고와 `TestHumanNamesAndIdentity`(DEV 결과 재사용)가 근거다. 에이전트 이름 저장 한 경로만 독립 확인했다.

#### ① 한계 (not_run)

- **대화 열기 결과가 고른 대상과 일치**: fixture가 `/home/conversation`을 막는다(GET 모두 403). 브라우저가 보낸 요청의 `from`·`to`(자기: a2→a1, 친구: a2→b1)가 고른 radio와 일치함까지만 확인했다. 서버가 그 대화를 열었는지는 확인하지 못했다.
- **권한 밖 agent 값 조작 거부**: `/home/new`의 쿼리 `from`·`to`는 선택에 영향이 없고(다른 계정 agent·존재하지 않는 값·스크립트 문자열 9종 모두 기본 선택 유지, 선택지에 타 계정 agent 없음) 이것은 PASS다. 열기 요청의 서버 거부는 fixture가 막아 확인하지 못했다. 선택지 출처와 `/home/conversation` 인가 검사는 DEV가 바꾸지 않았다고 보고했다(마크업만 변경).
- **명시 From/To 우선 선택**: `/home/new`는 실패한 열기 뒤 다시 그릴 때만 From/To를 받으므로 이 fixture로 재현할 수 없다.

### 새 결함

#### D-NM-02 (medium) 상대 이름이 40자이면 새 대화 수신 그룹 제목이 좁은 화면에서 넘친다 — FAIL

- 재현: fixture 계정 B로 `GET /home/new?peer=wXUsX1ctF16LR8h1bj2V6M`(A와의 새 대화). A의 이름은 `가`×40(공백 없음). 폭 390 또는 320에서 연다.
- 기대: designer 결정 2 수락 기준 ‘390·320에서 페이지 가로 넘침 0’. 제목 `<이름>의 수신 에이전트`가 줄바꿈된다.
- 실제: `fieldset#to`의 `legend`가 줄바꿈되지 않아 페이지가 390에서 282px, 320에서 352px 넘친다. 이름 끝과 ‘의’가 화면 밖에서 잘린다(캡처 `chrome-fix-new-peer-b-390.png`: 제목 둘째 줄이 ‘가가…’로 끝나고 ‘수신 에이전트’만 다음 줄). Chrome 155·Edge 155 동일. 1440px은 정상. A의 시점(상대 이름 ‘합성 회원 b’)은 정상이다.
- 원인(추정, 코드 미수정): legend에 긴 문자열 줄바꿈 규칙이 없다. 다른 이름 표시는 `.wrap`/`.name-preview`로 처리되어 통과한다. 수정은 DEV의 일(예: legend에 `overflow-wrap:anywhere`).
- 영향: 사용성(상대 이름이 잘려 누구에게 보내는지 확인이 어렵고 가로 스크롤). 데이터·권한 문제 없음. 증거: [browser-chrome-fix-recheck.json](browser-chrome-fix-recheck.json)·[browser-edge-fix-recheck.json](browser-edge-fix-recheck.json) `layout-390-new-peer-b`·`layout-320-new-peer-b`, SHA 5fa3333, fixture 18930. 검사를 약화하지 않았고 FAIL로 남겼다.
- 40자 이름은 U2가 허용하는 최대값이다. 이 화면의 legend는 이번 수정(select→radio fieldset)에서 생겼다.

### 영향 회귀 (이전 PASS 핵심)

| 영역 | 결과 | 근거 |
|---|---|---|
| 읽기 전용 HTTP 전체(A/B/C) | 19 PASS | [checks.py](checks.py) → [http-result-5fa3333.json](http-result-5fa3333.json). 새 대화·거절 안내 단언을 새 계약으로 갱신. 나머지 단언은 그대로 |
| 친구 요청 비노출 | 판정 PASS (결과 파일은 FAIL로 남음) | [regress_last.py](regress_last.py) → [regress-last-5fa3333.json](regress-last-5fa3333.json). 없는 코드·실제 코드(이미 대기)·대소문자 변경·공백 패딩 모두 200과 같은 안내문. A의 받은 요청은 1건 그대로. 결과 파일의 FAIL은 내 기대 오류다: 보낸 요청 행이 최대 +1이라고 단언했으나 대소문자 변경 코드는 두 번째 없는 코드라 +2(1→3)가 맞다. 단언을 +2로 고쳤지만 fixture를 더 바꾸지 않으려고 다시 실행하지 않았다. 원본 결과는 `regress-last-5fa3333.json`과 같은 내용의 [regress-last-5fa3333-first-run.json](regress-last-5fa3333-first-run.json)이다 |
| 422 입력 보존 | PASS | `followup.py` 선택 실행(빈 코드·공백·동의 누락, 입력 보존, alert·describedby, 요청 생성 없음) |
| 조회 radio·U6 입력 보존·429 | PASS | `checks.py` 안 `receipt_radio_selector_*`, `anonymous_429_*` |
| 친구 코드(읽기 전용·라벨·Ctrl+A) | PASS | `checks.py`의 읽기 전용 단언. 브라우저의 Ctrl+A/Ctrl+C는 이 재확인에서 다시 돌리지 않았다(`keys` 단계 not_run, 변경 없는 화면) |
| 레이아웃 3폭 | 41 PASS, 2 FAIL(D-NM-02) | Chrome·Edge 각 43건 중 같은 2건. [browser-chrome-fix-recheck.json](browser-chrome-fix-recheck.json), [browser-edge-fix-recheck.json](browser-edge-fix-recheck.json) |

실행 결과 요약

| 파일 | 결과 |
|---|---|
| [checks.py](checks.py) | 19 PASS |
| [recheck.py](recheck.py) → [recheck-result-5fa3333.json](recheck-result-5fa3333.json) | 8 PASS |
| [mutate.py](mutate.py) `ONLY_CHECKS`(이름 거부 1건) → [mutate-result-5fa3333.json](mutate-result-5fa3333.json) | 1 PASS. 나머지 비가역 검사는 다시 돌리지 않았다 |
| [followup.py](followup.py) `ONLY_CHECKS`(3건) → [followup-result-5fa3333.json](followup-result-5fa3333.json) | 3 PASS |
| [regress_last.py](regress_last.py) | 결과 파일 1 FAIL(내 기대 오류, 증거는 요건 충족. 위 표) |
| [browser.mjs](browser.mjs) `recheck` (Chrome 155.0.8059.80) | 43건 중 41 PASS, 2 FAIL(D-NM-02) |
| 같은 단계 (Edge 155.0.4283.62) | 43건 중 41 PASS, 2 FAIL(D-NM-02) |

브라우저 검사는 headless CDP, 비모바일 viewport 1440×900·390×844·320×640, 실제 키 이벤트(`Input.dispatchKeyEvent`)다. 제출 하나(거절되는 이름 저장 POST 1건)를 제외한 페이지 요청은 fixture origin의 GET이다. 폼 입력값 설정에는 `Input.insertText`를 썼다. 회원당 `/home` 요청 한도(O-NM-06)로 읽기 전용 검사에서 429를 16번 기다렸다.

### 재확인 캡처 (designer 시각 검수용)

Chrome 캡처 `chrome-fix-<화면>-<폭>.png`(1440·390·320). 변경 전 b2e683a 캡처 `chrome-*.png`는 그대로 있다.

| 항목 | 파일 |
|---|---|
| 새 대화(A 자기, 동명 둘) | `chrome-fix-new-self-a-*` |
| 새 대화(A→B, 이모지 40자) | `chrome-fix-new-peer-a-*` |
| 새 대화(B→A, 40자 이름 상대) — D-NM-02 | `chrome-fix-new-peer-b-*` |
| 내 에이전트(상세 펼침, 두 링크 분리) | `chrome-fix-agents-a-*` |
| 친구(받은 요청, 거절 앞 안내) | `chrome-fix-friends-a-*` |
| 이름 거부 오류 화면(금지 문자) | `chrome-fix-rename-refusal-a-*` |

`edge-fix-*.png`는 Edge 155의 같은 화면이다.

### 재확인 not_run

- 내 이름·에이전트 추가·device·remote 동의의 금지 문자 문구(fixture 차단). ‘대화 열기’의 서버 결과와 권한 밖 값 거부(fixture 차단). 명시 From/To 우선.
- 이전 not_run 그대로: 승인(gate) 화면, 조회 성공 화면, 구분값 충돌 확장, 시간 경과 만료·차단·친구 해제, 스크린리더, 실제 200% 확대, Firefox·Safari·모바일 실기.
- 등록 `make lint`·`make test`·integration: DEV 결과 재사용, 재실행 없음.
- fixture 상태: 이 재확인 뒤 보낸 요청 행이 C에 +2(없는 코드 2개). 친구 관계·코드·이름은 바뀌지 않았다.

---

## 이전 결과 (b2e683a, 이전 attempt)

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
