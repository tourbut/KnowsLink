---
title: SAR-SCHEDULE-001-QA — 일정 조회 합성 fixture 독립 QA 시나리오
status: draft
updated: 2026-10-21
owner: tester
tasks: [SAR-SCHEDULE-001-QA]
summary: "main 336533e의 일정 조회를 합성 제공자 fixture에서 SC-01~SC-08과 리뷰 low L1·L2·L3·L5로 나누고, 브라우저 표면과 agent API 표면의 통과 조건을 적는다"
---

# SAR-SCHEDULE-001-QA — 독립 QA 시나리오

대상 제품은 main `336533e6e15679d67118cbec5ab907768b80a70f`다. 결과는 [QA 보고서](../qa-reports/SAR-SCHEDULE-001-QA-tester-test/report.md)에 있다. Jev 웹 조작 도구는 쓰지 않는다. 이유: 합성 동의·철회·장애 토글의 순서와 `GET /fixture/gate/{id}` 관측을 고정해야 한다. 테스트 레벨은 lite다. 제품 코드는 수정하지 않는다. 실제 Google 계정은 쓰지 않는다.

## 공통 조건

- 격리 clone은 이 과제 전용이다. `remote.py`가 `336533e`를 checkout하고 `TestCalendarBrowserFixture`(태그 `integration && calendarfixture`)와 자체 Postgres 컨테이너를 기동한다.
- 로컬 접근은 `forward.py`의 loopback 포트 `18961`만 쓴다. 운영 서비스·운영 DB·실계정·외부 발송은 쓰지 않는다.
- 계정 a는 요청자, b는 연결 회원, c는 미연결 회원이다. 세션은 `POST /fixture/as/{a|b|c}`다.
- 브라우저 fixture가 허용하지 않는 `/v1/*`는 403이 정상이다. agent API 왕복(L1·L2·L3·L5, 범위 초과, commit, 관계 철회)은 같은 clone의 기존 integration·단위 검사로 실행한다.
- 상태를 바꾸는 판정 실행은 새 fixture에서 `walk.py` 한 번이다. 시각은 `walk-result.json`의 `mutations`에 남긴다. 그 전 실행은 검사 스크립트 판정 오류라 저장소 밖에 두었고 판정이 아니다.
- 화면 문구 판정은 서버 HTML(`walk.py`)과 headless Chrome 렌더 텍스트(`browser.mjs`)다. 시각 수락 캡처는 만들지 않는다.

## covers와 독립 checks

| covers | 표면 | 통과 조건 |
|---|---|---|
| SC-01 연결 | `walk.py` | a의 동의 URL scope가 `openid email`과 `calendar.freebusy`이고 `access_type=offline`, `include_granted_scopes=true`다. 일부 동의는 403이고 연결되지 않는다. 허용 뒤 ‘연결됨 — 바쁜 시간만 확인합니다’. 해제 뒤 해제 안내가 있고 해제 폼은 없다. b의 로그인 재확인 scope는 `openid email`만이다. |
| SC-01 토큰·revoke | integration `TestCalendarConnection` | 해제 뒤 봉인 토큰이 없고 합성 revoke 목록에 그 refresh token이 있다. |
| SC-02 승인 1회 | `walk.py` + integration `TestCalendarApprovalRunsOneLookup` | 승인 전 lookup이 없고 승인 뒤 `done` 1회. 재접근에 승인 버튼이 없다. 같은 폼 재전송은 409이고 disclosure는 같다. 만료·이미 거절 gate에 승인 버튼이 없다. |
| SC-03 공개 최소화 | `walk.py` | disclosure 키는 `window`·`granularity_min`·`busy`뿐이고 busy 항목은 `start`·`end`뿐이다. `secret-title`·참석자·위치 문자열이 disclosure와 승인 화면에 없다. |
| SC-04 범위 | integration `TestCalendarScopeEnforcedAtAccept`, 단위 `TestScheduleScopeBoundaries` | 단위 20·과거 시작·14일 초과·6회째가 422 `schedule_scope`다. |
| SC-04 친구 아님 | 브라우저 fixture | 새 `schedule.query`를 보낼 수 없다(`/v1` 403). 이 항목은 not_run으로 남긴다. |
| SC-05 실패 구분 | `walk.py` | 미연결 승인 POST는 409이고 lookup이 없다. 장애 토글 중 승인은 `failed`이고 연결은 유지되며 `done`이 아니다. grant 철회 뒤 승인은 `failed`이고 화면은 ‘다시 연결 필요’다. |
| SC-06 연결 해제 | `walk.py` | 승인돼 disclosure가 있는 gate를 consume하기 전에 해제하면 `agent_sees`가 `denied`이고 disclosure가 없다. 다른 대기 gate도 `denied`다. |
| SC-06 친구 해제·차단 | integration `relationship_revoked_after_approval_discloses_nothing` | unpair 뒤 gate가 `revoked`이고 disclosure가 없다. 친구 해제·차단 POST는 fixture가 막으므로 그 두 조작은 not_run이다. |
| SC-07 문구 | `walk.py`, `browser.mjs` | KO/EN 연결 문구, 승인 고정 문구, V1 미연결 한 줄, V3 reconnect 해제 안내, D2 해제 안내·해제 뒤 안내가 있다. 미연결 화면에 승인 버튼과 옛 ‘승인하면’ 두 줄이 없다. |
| SC-08 commit | integration `TestCalendarUnconnectedAndCommitStub` | commit 화면에 U5 문구가 있고 승인에 조회가 없다. 브라우저 fixture에는 commit gate가 없어 화면 왕복은 not_run이다. |
| L1 done 재전송 | integration `TestCalendarApprovalRunsOneLookup` | 같은 key·같은 내용의 재전송이 기존 receipt다. |
| L2 봉인 AAD | 단위 `TestSealedTokenRoundTripAndTamper` | 다른 member ID로는 봉인을 열 수 없다. |
| L3 미저장 refresh revoke | integration `TestCalendarConnection` | 다른 계정 동의로 저장하지 않은 refresh token이 revoke된다. 합성 동의 화면은 계정 전환이 없어 브라우저 왕복은 not_run이다. |
| L5 owner-revoke | integration `owner_revoke_disconnects_and_revokes` | owner-revoke가 연결을 지우고 provider revoke를 호출한다. fixture는 이 경로를 열지 않는다. |

## 시나리오 밖

- 실제 Google 계정 1경로는 D011 뒤 coor 몫이다(SC-09). 합성 통과는 실계정 수락이 아니다.
- 이메일 회원 한 줄은 이 fixture에 이메일 회원이 없다. integration `TestCalendarConnection`이 그 문구를 본다.
- 친구당 6회째, 서비스 재시작, consume 2회째, 스크린리더, 실제 브라우저 200%는 브라우저 fixture에서 실행하지 않는다.
- 시각 수락은 designer SAR-SCHEDULE-001-VISUAL-RECHECK가 맡는다.
