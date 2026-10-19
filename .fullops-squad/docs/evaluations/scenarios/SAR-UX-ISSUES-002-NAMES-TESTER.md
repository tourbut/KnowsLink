---
title: SAR-UX-ISSUES-002-NAMES-TESTER — 이름 중심 표시·친구 코드·친구 요청 안내 독립 QA 시나리오
status: draft
updated: 2026-10-19
owner: tester
tasks: [SAR-UX-ISSUES-002-NAMES-TESTER]
summary: "격리 names fixture에서 U2 이름 표시, 내 친구 코드, U4 친구 요청, U7 보조 화면 접근성, 1단위 조회 회귀를 HTTP와 실제 브라우저 키 입력으로 확인하는 시나리오와 통과 조건"
---

# SAR-UX-ISSUES-002-NAMES-TESTER — 독립 QA 시나리오

대상은 고정 SHA `b2e683ae677c78cd7bdc47c75820bcb191b94838`(fullops/dev, 제품 `516bf62`)다. 결과는 [QA 보고서](../qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/report.md)에 있다. Jev 웹 조작 도구는 쓰지 않는다. 이유: 비가역 POST 순서 통제, 실제 키 입력(CDP `Input.dispatchKeyEvent`), 비모바일 viewport 측정이 필요하다. 격리 fixture의 일반 조작은 이 시나리오의 허용 동작 안에서만 한다.

## 공통 조건

- origin `http://127.0.0.1:18930`만 쓴다(coor SSH 포워딩, 합성 DB). 운영 서비스·실계정·실메일·외부 발송은 쓰지 않는다.
- 허용 GET: `/`, `/home`, `/home/friends`, `/home/friends/{id}`, `/home/account`, `/home/agents`, `/home/legacy`, `/home/receipts`, `/home/new`. 허용 POST: `/home/friends/request`, `/home/friends/decision`, `/home/friend-code`, `/home/agent-name`, `/fixture/as/{a|b|c}`. 그 밖은 403이 정상이다.
- 합성 계정 A(40자 ‘가’ 이름, 같은 39자 이름 agent 둘, legacy pair), B(이모지 40자 agent 하나, A와 친구), C(agent 0개, A에게 대기 요청 보냄). 계정 전환은 `POST /fixture/as/{a|b|c}`.
- 비가역 동작(친구 결정·코드 교체·이름 변경)은 읽기 전용 검사와 캡처가 끝난 뒤 정해진 순서로 한 번 실행한다.
- 관측 수준: HTTP 응답 본문·헤더(`checks.py`, `mutate.py`), Chrome 155·Edge 155 CDP의 실제 키 입력·레이아웃 측정·접근성 트리·캡처(`browser.mjs`). 스크린리더 음성·실제 200%·보조 브라우저의 모든 조합은 범위 밖(not_run)이다.

## covers와 독립 checks

| covers | 독립 checks (파일) | 통과 조건 |
|---|---|---|
| U2 이름 먼저+종류+구분값 | `checks.py` name_first_*, `browser.mjs` | 조회 radio·내 에이전트·기존 연결·새 대화·친구 상세가 이름, 종류, 구분값 순서로 보인다. 동명 둘의 라벨이 서로 다르다. |
| U2 40자 한글·이모지 | `browser.mjs` layout-* | 1440·390·320에서 가로 넘침 0, 이름이 잘리거나 `text-overflow:ellipsis`가 아니다. |
| U2 비친구 이름 | `checks.py` non_friend_* | ‘확인할 수 없는 에이전트’와 허용된 ID만 보이고 이름은 없다. 타인 계정에도 이름이 새지 않는다. |
| U2 이름 변경 | `mutate.py` rename_* | 변경 뒤 ID·링크·친구 관계가 같다. 제어문자·방향문자·41자·빈 이름은 422이고 이름이 바뀌지 않는다. |
| U2 내 친구 코드 | `checks.py` friend_code_*, `browser.mjs` code-select-all | 보이는 라벨, `readonly`, 키보드 전체 선택. 코드 교체 뒤 옛 코드로 새 요청이 생기지 않는다. |
| U4 요청 제출 | `mutate.py` request_* | 없는·유효·대소문자 변경·공백 패딩 코드의 응답 상태와 안내문이 같다. 빈 코드·동의 누락은 422이며 입력을 보존한다. |
| U4 대기 표시·행동 | `checks.py` pending_*, `mutate.py` | 받은 대기(수락·거절, 절대 만료 시각), 보낸 대기(요청 취소, 상대 정보 없음), 반복 요청이 기한을 늘리지 않는다. |
| U4 에이전트 0개 | `checks.py` zero_agent_*, `mutate.py` | C는 요청 가능하고 빈 선택기가 없다. 수락 뒤 C의 친구 상세에 대화 링크가 없다. |
| U4 GET 불변 | `checks.py` get_no_side_effect, `browser.mjs` | GET·뒤로가기·쿼리 위조로 관계가 생기거나 수락되지 않는다. |
| U7 접근성 | `checks.py` landmarks_*, labels_*, `browser.mjs` | 고유 title·단일 h1·main·skip, 정확한 내부 복귀, status/alert, 반복 버튼 이름에 대상과 구분값, Tab 이동·보이는 focus, 44px. |
| 1단위 회귀 | `checks.py` receipt_* | radio 선택기 값·U6 입력 보존·없는/타인 동일 문구, 익명 429 화면. |

## 시나리오 밖 (인계)

- 승인 화면(`/gate`)은 fixture가 허용하지 않아 열 수 없다. DEV Go 테스트가 근거이며 not_run으로 남긴다.
- 조회 성공 화면은 fixture에 receipt 데이터가 없어 열 수 없다(not_run). 실제 인증·DB 화면은 범위 밖이다.
- 시각 수락은 designer, 고정 SHA 코드 리뷰와 통합은 coor가 조정한다.

## 재확인 5fa3333 (DEV 수정 4건)

대상 `5fa3333b5948fd36792662b202915bc6243cac14`(제품 `97a8899`), 같은 fixture(재기동 상태). 상태를 바꾸는 POST는 마지막 요청 비노출 회귀와 거절되는 이름 저장만 쓴다(designer와 fixture 공유). 결과는 [보고서](../qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/report.md)의 ‘재확인’ 절에 있다.

| covers | 독립 checks (파일) | 통과 조건 |
|---|---|---|
| ① 새 대화 radio | `recheck.py` new_conversation_*, `browser.mjs recheck` new-* | select 없음, 그룹 legend와 `이름 / 종류 / 구분값` 라벨, 기본 선택 유지, 방향키 선택, Enter 요청의 `from`·`to`가 고른 값, 쿼리 위조로 선택지가 바뀌지 않음, 3폭 가로 넘침 0·44px |
| ② 링크 분리 | `recheck.py` my_agents_two_links_*, `browser.mjs recheck` agents-* | 두 링크가 다른 줄(별도 `p`), 각 44px, 정확한 이동 대상 |
| ③ D-NM-01 | `mutate.py` agent_rename_rejects_unsafe_names_without_change(갱신), `followup.py`, `browser.mjs recheck` rename-refusal-* | 금지 문자 거부 시 전용 문구·입력란에 제거값·이름 불변, 길이·빈 이름은 기존 문구·입력값 보존 |
| ④ 거절 안내 | `recheck.py` received_request_has_refusal_notice_*, `browser.mjs recheck` deny-* | 수락→안내→거절 순서, 거절의 `aria-describedby`와 접근성 설명 일치 |
| 영향 회귀 | `checks.py`, `regress_last.py`, `followup.py` 선택 실행 | 이전 PASS 핵심이 유지됨 |

시나리오 밖: `/home/conversation`·`/home/profile`·`/home/agents`는 fixture가 막아 서버 결과·내 이름·에이전트 추가 문구는 not_run이다.
