---
title: SAR-SCHEDULE-001-QA 독립 QA
status: review
updated: 2026-10-21
owner: tester
tasks: [SAR-SCHEDULE-001-QA]
summary: main 336533e의 일정 조회를 합성 제공자 fixture에서 SC-01~SC-08과 리뷰 low L1·L2·L3·L5로 독립 QA했다
---

# SAR-SCHEDULE-001-QA 독립 QA

## 대상과 결론

대상 제품은 main `336533e6e15679d67118cbec5ab907768b80a70f`다. 원격 clone은 그 SHA를 checkout했고 `git rev-parse HEAD`가 같은 값이다. 이 워크트리의 제품 트리도 그 SHA와 같다. 제품 코드는 바꾸지 않았다. 실제 Google 계정은 쓰지 않았다.

독립 QA이며 제품 수락이 아니다. 실행한 SC-01~SC-08과 L1·L2·L3·L5는 PASS다. 아래 not_run은 실패가 아니다. 실계정 1경로(SC-09, D011)는 실행하지 않았다.

## 방법

격리 clone `/tmp/knowslink-schedqa-336533e6e156-941429b7`에서 `TestCalendarBrowserFixture`(태그 `integration,calendarfixture`)와 과제 전용 Postgres를 기동했다. 로컬 관측은 `http://127.0.0.1:18961`이다. 판정 실행은 새 fixture에서 `walk.py` 한 번과 headless Chrome(`browser-early.json`, `browser-reconnect.json`, `browser-off.json`)이다. 상태 변경 시각은 [walk-result.json](walk-result.json)의 `mutations`에 있다. 2026-10-21T01:31:14Z~01:31:17Z.

브라우저 fixture의 `/v1/*`는 403이다. agent API는 같은 clone의 integration·단위 검사로 실행했다. 페이지 CSP가 `default-src 'none'`이라 화면 안 `fetch` 로그인은 쓰지 않았다. Node가 `POST /fixture/as/{a|b|c}`를 하고 CDP로 세션 쿠키를 넣었다. 쿠키 값은 기록하지 않았다.

첫 fixture 실행 4건은 검사 스크립트 판정이었다. 영문 apostrophe가 `&#39;`로 저장돼 raw 검색이 실패했고, 결정 뒤 페이지에는 csrf가 없어 재전송이 403이었고, 장애 JSON의 공백 검사가 어긋났고, 고정 디버그 포트가 다른 브라우저에 붙었다. 제품 결함으로 남기지 않았다. 그 기록은 저장소 밖 `schedule001-qa/first-walk`에 있다. 스크립트를 고친 뒤 새 fixture로 다시 실행한 결과가 이 보고의 판정이다.

## 판정

| 항목 | 판정 | 근거 |
|---|---|---|
| SC-01 연결·해제 화면 | PASS | [walk-result.json](walk-result.json) `sc01-*`. scope는 `openid email`과 freebusy, `access_type=offline`. 일부 동의 403. 해제 뒤 연결 버튼만 남음. 로그인 재확인 scope는 `openid email`. |
| SC-01 토큰 없음·revoke | PASS | [integration.txt](integration.txt) `TestCalendarConnection`. 해제 뒤 봉인 토큰이 없고 합성 revoke 목록에 그 refresh token이 있다. |
| SC-02 승인 1회·재전송·재접근 | PASS | walk `sc02-*` 409와 승인 버튼 없음. integration이 같은 csrf 재전송 409, 프로세스 안 handler 재기동 뒤에도 조회 1회, consume 2회째 409. |
| SC-03 공개 최소화 | PASS | walk `sc03-*`. disclosure 키는 window·granularity_min·busy뿐. 제목·참석자·위치 문자열이 화면과 disclosure에 없다. integration도 로그와 전달 본문에 `secret`이 없음을 본다. |
| SC-04 범위·한도 | PASS | integration `TestCalendarScopeEnforcedAtAccept`와 로컬 [local-unit.txt](local-unit.txt) `TestScheduleScopeBoundaries`. 단위 20·과거·14일 초과·6회째가 422 `schedule_scope`. |
| SC-05 미연결·장애·재연결 | PASS | walk `sc05-*`. 미연결 승인 409, lookup 없음. 장애는 `failed`이고 연결 유지. grant 철회는 `failed`와 `다시 연결 필요`. `done`이 아니다. |
| SC-06 미전달 결과 | PASS | walk `sc06-*`. consume 전 disclosure가 있던 gate가 해제 뒤 `agent_sees=denied`이고 disclosure가 없다. 다른 대기 gate도 denied. |
| SC-06 unpair | PASS | integration `relationship_revoked_after_approval_discloses_nothing`. |
| SC-07 문구 KO/EN, V1·V3·D2 | PASS | walk `sc07-*`와 브라우저 3파일. 연결·승인 범위·미연결 한 줄·reconnect 해제 안내·해제 뒤 안내. 미연결 화면에 승인 버튼과 옛 승인 두 줄이 없다. 1440·390·320 가로 넘침 0. |
| SC-08 commit U5 | PASS | integration `TestCalendarUnconnectedAndCommitStub`. U5 문구, 승인 시 조회 0, disclosure 없음, done 403. |
| L1 같은 key done 재전송 | PASS | integration `TestCalendarApprovalRunsOneLookup`. 같은 key·같은 내용이 기존 receipt. |
| L2 봉인 AAD | PASS | [local-unit.txt](local-unit.txt) `TestSealedTokenRoundTripAndTamper`. |
| L3 미저장 refresh revoke | PASS | integration `TestCalendarConnection`의 다른 계정 동의. 저장하지 않은 refresh token 1건이 revoke된다. |
| L5 owner-revoke | PASS | integration `owner_revoke_disconnects_and_revokes`. |

추가 실행: 지시 목록 밖의 L4(body close)는 로컬 `TestRevokeClosesBodyOnEveryStatus` PASS다. 판정 표에는 넣지 않는다.

## not_run

- SC-04 친구 아님. 심각도 없음. fixture가 `/v1`을 403으로 막아 새 `schedule.query`를 보낼 수 없다. calendar integration은 기본 친구 쌍만 쓰고, 친구 아닌 상대의 `schedule.query` 거부 테스트가 없다. 소스 읽기로 PASS를 주지 않았다. 기대는 거부, 실제는 미실행.
- SC-06 친구 해제·차단. 심각도 없음. fixture가 `/home/friends/decision`을 403으로 막는다. unpair PASS는 이 두 조작과 별도다. 기대는 대기·승인 조회의 즉시 중단, 실제는 미실행.
- SC-08 브라우저 commit 화면. fixture가 commit gate를 심지 않는다. U5는 integration PASS.
- L3 브라우저 계정 전환. 일부 동의는 refresh token을 발급하지 않는다. 다른 계정 동의는 integration의 `sub-other` grant다. 브라우저 계정 선택 화면은 없다.
- L3 이미 연결된 자기 계정 token을 revoke하지 않는 분기. DEV 보고와 같이 자동 검사가 없다. 이번 QA도 그 분기를 실행하지 않았다. PASS로 세지 않는다.
- 이메일 회원 안내 문구의 브라우저. fixture 계정은 모두 Google issuer다. integration `TestCalendarConnection`이 그 한 줄을 본다.
- 브라우저에서 6회째 한도. 대기 gate가 3개이고 `/v1`이 막혀 있다. 6회째는 integration PASS.
- OS 서비스 재시작, 스크린리더, 브라우저 200% 확대. 미실행. integration의 handler 재기동과 consume 2회째는 SC-02 행에 포함된다. 가로 넘침 측정은 200% 확대가 아니다.
- SC-09 실제 Google 계정. D011이 열려 있다. 합성 통과는 실계정 수락이 아니다.
- `KNOWSLINK_CALENDAR_KEY` 교체·분실 복구. [transition.md](../../../operations/transition.md)가 절차 미정이라고 적는다. 운영 env는 보지 않았다.

## 검사

Windows에는 `make`가 없다. 아래는 원격 Linux의 exact SHA clone이다.

- `make lint` EXIT=0. [make-lint.txt](make-lint.txt). 원격 원 로그의 stderr에는 public-config 음수 검사의 `invalid tunnel UUID`가 stdout 뒤에 붙어 있다. 그 검사의 stdout은 `validation failures propagated`이고 종료코드는 0이다.
- `make test` EXIT=0. [make-test.txt](make-test.txt). integration 태그는 이 명령에 없다.
- calendar integration `-race -tags=integration -v` EXIT=0. [integration.txt](integration.txt). 전용 DB는 검사 뒤 `docker stop`했다.
- 로컬 단위 3개 EXIT=0. [local-unit.txt](local-unit.txt). 제품 트리가 336533e와 같은 이 워크트리에서 실행했다.

## fixture

판정에 쓴 컨테이너는 `kl-schedqa-fixture-336533e6e156-9a58544e`다. 수명은 기동 후 12시간이다. 판정 로그를 남긴 뒤 tester가 `remote.py cleanup`으로 fixture 프로세스와 그 컨테이너를 멈췄다. clone 디렉터리는 지우지 않았다. 로컬 임시 기록은 `C:\Users\shin\AppData\Local\KnowsLinkDevTools\schedule001-qa`와 Chrome 프로필 `schedqa-chrome-*`다. 삭제하지 않았다.

시각 수락 캡처는 만들지 않았다. 화면은 문구와 동작만 판정했다. 시각 수락은 designer SAR-SCHEDULE-001-VISUAL-RECHECK 몫이다.
