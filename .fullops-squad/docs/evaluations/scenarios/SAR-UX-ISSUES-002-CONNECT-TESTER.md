---
title: SAR-UX-ISSUES-002-CONNECT-TESTER — 연결 안내·재개·보낸 요청 코드 끝 4자 독립 QA 시나리오
status: draft
updated: 2026-10-20
owner: tester
tasks: [SAR-UX-ISSUES-002-CONNECT-TESTER]
summary: "격리 connect fixture에서 U8 호스트 선택지·재개 행, U7 연결 화면 접근성·복귀, O-NM-03 보낸 요청 코드 끝 4자, low 3건 회귀를 HTTP와 실제 브라우저 키 입력으로 확인하는 시나리오와 통과 조건"
---

# SAR-UX-ISSUES-002-CONNECT-TESTER — 독립 QA 시나리오

대상은 DEV `8d8bad4cd75a37546ba0f803e3c0c8290ddb6a87`(fullops/dev, 제품 `37febf9`)다. 결과는 [QA 보고서](../qa-reports/SAR-UX-ISSUES-002-CONNECT-TESTER-test/report.md)에 있다. Jev 웹 조작 도구는 쓰지 않는다. 이유: 비가역 생성의 순서 통제, 실제 키 입력(CDP `Input.dispatchKeyEvent`), 비모바일 viewport 측정이 필요하다. 테스트 레벨은 lite다.

## 공통 조건

- origin `http://127.0.0.1:18932`만 쓴다(coor SSH 포워딩, 합성 DB, 수명 12시간). 운영 서비스·실계정·실메일·외부 발송·실제 호스트 설치는 쓰지 않는다.
- 허용 GET: `/`, `/home`, `/home/agents`(query 포함), `/home/friends`, `/home/account`, `/home/legacy`, `/home/new`, `/home/receipts`, `/home/connections/*`, `/connect/*`. 허용 POST: `/home/agents`, `/home/connect`, `/home/confirm`, `/home/cancel`, `/home/friends/request`, `/home/friends/decision`, `/fixture/as/{a|b|c}`. 그 밖은 403이 정상이다.
- 합성 계정 a(연결 시도 waiting·prepared·approved 각 1, 키 등록 agent 1, 남은 1칸, 보낸 요청 5행, c에게서 받은 요청 1), b(만료·취소·device 승인 뒤 미저장·remote 앱 복귀 대기·교환 유실·교환 완료 agent, a에게서 받은 요청 1), c(agent 0개). 계정 전환은 `POST /fixture/as/{a|b|c}`.
- designer가 같은 fixture를 GET으로 함께 쓴다. 상태를 바꾸는 POST는 마지막 `mutate.py` 한 번(agent 1개 생성)과 거절 기대인 `stale.py`·여섯째 생성뿐이다. 이 생성은 되돌릴 수 없고 a의 남은 1칸을 쓴다.
- 이 fixture는 `/home/*` 전체에 회원당 분당 40회 한도가 있다. 러너는 429를 기다린다(클릭 이동은 `enterNavRL`로 reload).
- 관측 수준: HTTP 응답 본문·헤더(`checks.py`, `mutate.py`, `stale.py`), Chrome 155·Edge 155 headless CDP의 실제 키 입력·레이아웃 측정·접근성 트리·캡처(`browser.mjs`). 소스·Go 테스트 근거는 DEV SHA에서 읽기만 한다.

## covers와 독립 checks

| covers | 독립 checks (파일) | 통과 조건 |
|---|---|---|
| U8 선택지 순서·문구 | `checks.py` choices_*, guide_* | 5개 선택지가 설계표 순서·문구·보조 문구와 같다. 선택 전에는 안내와 미리 선택이 없다. |
| U8 선택 뒤 안내 | `checks.py` guide_*, `browser.mjs` keys | 선택마다 준비물·권한·다음 행동과 README 앵커. 4번은 `/knowslink-session on`, 세션 연결은 ‘앱에서 확인 필요’. 5번은 고정 문장과 웹 링크. |
| U8 금지 표시 | `checks.py` hidden_hosts_*, guide_other_*, no_online_* | OpenClaw·Meta Muse 미표시, 5번에 버튼·명령·배지·도구 이름 없음, 온라인 배지·주장 없음. |
| U8 GET은 무변경 | `checks.py` selection_creates_*, `browser.mjs` | 선택·새로고침·뒤로가기로 agent 수가 같다. 위조 query는 무시된다. |
| U8 상태별 재개 행 | `checks.py` resume_* | local waiting/prepared/approved/키 등록, remote 앱 복귀·유실·완료, 만료·취소, device 승인 뒤 미저장이 서버 상태와 맞는 설명과 다음 행동 하나를 가진다. 완료·온라인 단정이 없다. |
| U8 중복 생성·키 보존 | `mutate.py` | 같은 attempt 동시 2회 제출로 agent 1개, 이후 재전송도 같은 위치, 새 attempt는 한도 거절, 기존 키 카드 불변. |
| U8 5분 재확인 | `stale.py` | 로그인 5분 뒤 생성 POST가 거부되고 agent가 늘지 않는다. |
| O-NM-03 | `checks.py` onm03_*, `browser.mjs` friends | 보낸 행 ‘입력한 코드 …abcd’(대소문자 유지), 형식 다름, 이전 요청, 종료 행 유지, 있는/없는 코드 행의 같은 모양, 받는 쪽·URL에 미노출, 취소 버튼 이름에 끝 4자. |
| U7 연결 화면 | `checks.py` connection_*, `browser.mjs` | 확인된 단계만 `aria-current=step`(만료·취소는 없음), 키 등록됨 표시, 정확한 복귀(`/home/agents#agent-{id}` 도착·스크롤), title·h1·main·skip, 반복 이름. |
| U7 키보드·3폭 | `browser.mjs` layout·keys | Tab 순회에 함정 없음·보이는 focus·44px, 선택기 화살표·Enter, 1440·390·320 넘침 0. |
| low 3건 | 소스·Go 테스트 근거(DEV SHA), `checks.py` | 익명 429 nav에 내 에이전트 링크 없음(템플릿), 승인 목록 ‘요청 검토’ 이름에 받은 내 에이전트(템플릿), gate_summary comma-ok. fixture가 `/gate`·`/home/receipts`를 막아 화면 확인은 not_run. |

## 시나리오 밖 (인계)

- remote 앱 복귀 화면(`/connect/remote/return` 계열)과 remote 동의는 fixture에서 열 수 없다. 소스 템플릿과 DEV 렌더로만 확인하고 not_run으로 남긴다.
- 차단 요청 만료 뒤 재요청 행 차이는 fixture에 차단 경로가 없어 소스 근거로만 기록한다.
- 스크린리더 음성, 실제 브라우저 200% 확대, 보조 브라우저의 모든 조합은 not_run이다. 시각 수락은 designer, 고정 SHA 코드 리뷰와 통합은 coor가 조정한다.

## 재확인 attempt (DEV `abc2191`, 제품 `d6369a4`)

같은 fixture origin·허용 목록에 계정 d·e가 추가되었다(`POST /fixture/as/{a..e}`). manifest `rerequest_codes`(`blocked`·`receive_limit`·`unknown_cancelled`)로 재요청을 확인한다. 이 fixture는 초기 상태에서 한 번만 쓸 수 있는 비가역 경로이므로 순서는 `checks.py`(읽기) → `rerequest.py` → `cap.py` → `stale.py`(5분) → `mutate.py`(마지막 생성) → `browser.mjs`다.

| covers | 독립 checks (파일) | 통과 조건 |
|---|---|---|
| F1·V-CN-01·O-CT-01 재요청 동일화 | `rerequest.py` | A가 세 코드를 다시 보낸 뒤 보낸 행 수(7)·순서·상태·끝 4자·기한 형식·POST 응답이 세 경로에서 같다. D·E(와 대조 B·C) 화면이 바이트 단위로 같다. 대기 중 반복은 무동작, 취소 뒤 재요청도 같다. |
| mine 대기 상한 | `cap.py` | 대기 보낸 행이 10에서 멈추고(보류 행 포함 계산) 11번째는 조용히 버려진다. |
| D-CT-01 | `checks.py` device_lost_* | 연결 화면이 카드와 같은 ‘확인 필요’이고 ‘만료되었습니다’·‘새 키가 연결되지 않았습니다’·‘새로 시작’이 없다. 기존 기대를 낮추지 않고 단정 부재 검사를 더했다. |
| V-CN-02, O-CT-02, O-CT-04, F3 | `checks.py` rejected_connection_*, legacy_*, `stale.py` o_ct_04_wording, 소스 | 거부 복귀 ‘내 에이전트’, 보이는 ‘이 화면에서 찾기’(aria-current 없음), 재확인 안내, ‘키 키’ 없음. F2는 소스·DEV 테스트 근거. |
| 회귀 | `checks.py`, `mutate.py`, `browser.mjs` | 이전 PASS 핵심(5선택지·재개 행·중복 생성 방지·O-NM-03·3폭 넘침 0). |

`checks.py`의 O-NM-03 기대값은 fixture 인스턴스마다 바뀌는 임의 코드 접미를 하드코딩하고 있었다(`XXK2`·`Z-3N`). 새 fixture에서 매번 같은 판정을 하도록 b·a의 친구 코드에서 끝 4자를 읽게 바꿨다. 판정 기준은 같다.
