---
title: SAR-SCHEDULE-001-DEV DEV 보고
status: review
updated: 2026-10-21
owner: dev
tasks: [SAR-SCHEDULE-001-DEV]
summary: "Google Calendar 빈 시간 조회 연동(SC-01~08)의 기술 계획, 외부 API 근거, 자동 검사·렌더·fixture 결과와 미실행 항목"
---

# SAR-SCHEDULE-001-DEV DEV 보고

## 대상과 결론

대상 후보는 fullops/dev `dfd5500`(첫 구현 `4b3cad6`, gofmt 보정 `c0eb0bc`, 통합 테스트 보정·문서 `7ad11be`, 테스트 기대 보정 `dfd5500`)다. 기준은 지시서 base `826ea80`이다. 실제 Google 계정·Cloud Console·운영 env·운영 배포는 건드리지 않았다(D011 열림). SC-01~SC-08은 합성 제공자 위의 자동 검사와 격리 fixture로 확인했다. DEV 완료이며 제품 수락이 아니다. 리뷰·tester 독립 QA·designer 직접 시각 검수·실계정 1경로(D011 뒤)·통합은 coor가 진행한다.

## 기술 계획

| 항목 | 결정 |
|---|---|
| 토큰 암호화·키 관리 | refresh token만 AES-256-GCM(stdlib `crypto/aes`)으로 봉인해 `State.Calendars[member].Sealed`에 둔다. 키는 env `KNOWSLINK_CALENDAR_KEY`(base64 32바이트)로만 프로세스 메모리에 있다. 키가 없으면 캘린더 절은 ‘준비 중’, 잘못되면 기동 실패. 평문·키는 상태·로그·화면에 없다. |
| 연결 상태 저장 | `Calendar{Sealed, State connected\|reconnect, Changed}`. JSON 필드 추가만이며 SQL migration 없음. |
| 증분 동의 흐름 | 기존 Google attempt 재사용(`googleAttempt.Calendar`). 요청 scope `openid email calendar.freebusy`, `access_type=offline`, `include_granted_scopes=true`, 같은 sub만 허용. 토큰 응답 `scope`에 freebusy가 없거나 refresh token이 없으면 저장하지 않는다. 로그인·재확인 scope는 그대로다. |
| freebusy 호출·반올림 | 승인 handler가 lock 밖에서 refresh→`POST freeBusy`(`items:[{id:"primary"}]`, 요청 window). `busy[].start/end`만 읽어 단위로 바깥 반올림·window 클립·병합 후 canonical `{window, granularity_min, busy}` 저장. |
| 승인 1회 실행 보장 | approve transaction이 연결 확인 뒤 `Gate.Lookup{Started}` lease를 쓰고, 두 번째 transaction이 결과를 1회 기록(`g.Lookup.State==""`일 때만). 중복 제출은 pending이 아니라 409. 재시작은 lease 1분 만료를 sweep이 `failed`로 닫고 재조회하지 않는다. agent에는 lookup이 끝나기 전까지 `pending`. |
| 친구당 횟수 집계 | 수락 시 같은 owner 쌍의 24시간 내 `schedule.query` 수락 수(Messages 24h 보존 재사용)가 5면 거부. 단위·기간 검사와 함께 `scheduleAllowed`(`schedule_scope` 422). 수치는 `scheduleScope` 한 곳. |
| 실패 분류 | 미연결·reconnect: 승인 버튼 없음·POST 409(거절은 가능). provider refresh 거부(400/401)→`reconnect`+`failed`. 일시 오류→8초 예산 안 최대 3회→`failed`, 연결 유지. `done`은 서버 disclosure echo만 수락하므로 실패를 `done`으로 보낼 수 없다. |
| 철회 처리 | 해제: 토큰 삭제·대기/미조회 승인 gate `denied`·provider revoke 요청(실패는 상세 없는 로그). 친구 해제·차단·unpair: 기존 sweep `revoked`가 consume·result를 막고 disclosure 삭제. 재연결(덮어쓰기)은 revoke하지 않는다(같은 project의 revoke는 모든 grant를 끊음). |
| 합성 제공자 구조 | `calendar_provider_test.go`: `/fixture/google/{auth,token,freebusy,revoke,toggle/*}`. code/PKCE/refresh/revoke, 비밀 필드(`summary`·`attendees`·`location`)가 섞인 busy 응답, 철회·장애 토글. integration·fixture가 공유. |
| adapter | `core.ts`가 consume 응답 disclosure를 그대로 `done`으로, `lookup=failed`는 `failed`, 그 밖은 `denied`로 보낸다. |

## 외부 API 근거

Context7 `resolve-library-id`는 monthly quota 초과로 두 번 모두 실패했다(2026-10-20). 지시서대로 공식 문서를 WebFetch로 직접 확인했다(발췌 요약이므로 원문과 다를 수 있다).

| 사실 | 출처 |
|---|---|
| `POST https://www.googleapis.com/calendar/v3/freeBusy`, 요청 `timeMin`·`timeMax`(RFC3339)·`items[].id`, 응답 `calendars.<id>.busy[].start/end`(start 포함, end 제외)·`errors[]`. 제목 필드 없음. 허용 scope에 `calendar.freebusy` 포함 | [freebusy.query](https://developers.google.com/workspace/calendar/api/v3/reference/freebusy/query) |
| `https://www.googleapis.com/auth/calendar.freebusy` = ‘내 캘린더의 가능 여부 확인’. 문서는 민감 여부를 분류하지 않음(D011) | [Calendar auth](https://developers.google.com/workspace/calendar/api/auth) |
| 증분 동의 `include_granted_scopes=true`, `access_type=offline`이 첫 교환에서 refresh token 반환, `prompt=consent`, 토큰 응답 `scope`는 공백 구분 목록이며 일부만 허용될 수 있음 | [Web server OAuth](https://developers.google.com/identity/protocols/oauth2/web-server) |
| refresh: `POST https://oauth2.googleapis.com/token` `grant_type=refresh_token`; revoke: `POST https://oauth2.googleapis.com/revoke` `token=`, 200 성공·400 실패; revoke는 project의 모든 grant를 끊음 | 같은 문서 |
| 라이브러리: `golang.org/x/oauth2` v0.36.0(기존)의 `Config.TokenSource`·`AccessTypeOffline`·`SetAuthURLParam`, `RetrieveError.Response.StatusCode`. 새 의존성 없음 | go.mod, 설치된 모듈 소스 |

## SC 대응과 자동 검사

원격 Linux 격리 clone(자기 Postgres 컨테이너)에서 실행한 integration 결과는 아래 ‘검사 로그’ 표다.

| SC | 검사 | 확인 내용 |
|---|---|---|
| SC-01 연결 | `TestCalendarConnection` | 동의 URL scope·offline·include_granted_scopes, 로그인/재확인 scope `openid email` 불변, 일부 동의 403·다른 계정 403·미저장, 봉인 저장(평문 없음)·복호 일치, 해제 뒤 토큰 없음·provider revoke 토큰 일치, 이메일 회원 한 줄 안내 |
| SC-02 승인 1회 | `TestCalendarApprovalRunsOneLookup` | 승인 1회에 freebusy 1회·요청 window 일치, 중복 제출 409, 재접근 폼 없음, 재시작 Service에서 조회 수 불변·`lookup: done`, consume 2회째 409 |
| SC-03 공개 최소화 | 같은 검사 | 저장 disclosure가 반올림 기대값과 canonical 동일, 상태 JSON·결과 envelope·로그에 `secret`/refresh/`attendees`/`location` 없음, 변조·누락 disclosure의 `done` 403, 올바른 echo 200, 두 번째 결과 403, 수락 뒤 disclosure 삭제 |
| SC-04 범위 집행 | `TestCalendarScopeEnforcedAtAccept`, 단위 `TestScheduleScopeBoundaries` | 단위 20·과거 시작·14일 초과 422 `schedule_scope`, 5회 뒤 6회째 422, 24h 지나면 수락 |
| SC-05 실패 구분 | `TestCalendarUnconnectedAndCommitStub`, `TestCalendarRevocationAndFailureClasses` | 미연결·reconnect 승인 불가(조회 0), provider 철회→`reconnect`·`failed`·`done` 403·`failed` 200, 장애→3회 시도 뒤 `failed`·연결 유지, lease 만료→`failed` |
| SC-06 철회 | 같은 검사 | 해제 뒤 대기 gate `denied`·승인 409·consume 409·조회 0, 승인 뒤 unpair→`revoked`·consume 409·disclosure 삭제 |
| SC-07 문구 | 위 검사 + `TestHumanAccountCalendarStates`, `TestMemberUIFixtures`, `TestEnglishScreensHaveNoKorean` | 연결 회원 고정 문구·‘공개되는 내용: 바쁜 시간대만’·승인 버튼, 미연결 회원 안내·링크·거절만, 설정 절 5상태, 영어판 |
| SC-08 commit 유지 | `TestCalendarUnconnectedAndCommitStub` | commit 화면 U5 문구 유지, 승인에 조회 0, consume `disclosure:false`, `done` 403 |

단위 검사(`calendar_test.go`): 반올림 7사례·disclosure JSON, 봉인 왕복·변조·잘못된 키·`ConfigureGoogle` 키 거부(값 비노출), 범위 경계 8사례+친구당 횟수, `relay.result` disclosure shape 7사례, lease sweep·disclosure 정리. 기존 회귀는 전체 `go test ./internal/relay`(Windows)와 원격 `make test`·integration이다.

기존 테스트 보정: `queryBody()`를 현재 시각 기준 window로 바꿨다(범위 집행). `TestSharedMessageCapacityHTTP`·`owner_gate_csrf_binding_and_no_effect`는 `schedule.commit`으로 seed한다(합성 owner는 캘린더가 없어 query를 승인할 수 없음). `member_gate_isolation_and_decision`·`TestApprovalSummaryBoundaries`는 회원 캘린더를 연결 상태로 둔다(제공자 없음→`failed`).

## 검사 로그

원격 Linux 격리 clone(`/tmp/knowslink-schedule001-<sha12>-<rand>`, 자기 Postgres 컨테이너, [remote.py](remote.py))에서 실행했다. 로그는 이 폴더에 있고 SHA256은 [raw-log-sha256.txt](raw-log-sha256.txt)다.

| SHA | make lint | make test | integration | fixture-compile | fixture 기동 |
|---|---|---|---|---|---|
| dfd5500 | EXIT=0 | EXIT=0 | EXIT=0 (`-skip TestRemoteConsentBrowserReturn`) | EXIT=0 | EXIT=0 (실제 기동·loopback probe) |
| 7ad11be | EXIT=0 | EXIT=0 | EXIT=1(테스트 기대 2건: 이름 없는 이메일 회원의 303, approved 화면 문구) | 미실행 | 미실행 |
| c0eb0bc | EXIT=0 | EXIT=0 | EXIT=1(아래 보정) | 미실행 | 미실행 |
| 4b3cad6 | EXIT=2(PATH의 다른 Go 버전 gofmt 정렬 차이) | 미실행 | 미실행 | 미실행 | 미실행 |

c0eb0bc 실패 7건의 원인과 보정: (1) Postgres JSONB가 저장 disclosure의 키 순서·공백을 바꿔 바이트 비교가 실패→`canonicalRaw`로 재정규화해 비교·반환. (2) 미연결 회원에게 결정 폼 전체를 숨겨 기존 거절 검사와 거절 UX가 깨짐→거절 전용 폼 유지, 승인 버튼만 제거. (3) 합성 owner·이메일 회원의 query 승인 기대→commit seed 또는 연결 상태로 보정. (4) 멱등키 16자 미만·해제 안내 query 누락은 테스트 오류. 제품 수치·규칙은 바꾸지 않았다.

로컬 Windows: `go vet`(태그 없음·integration·calendarfixture), `go test ./internal/relay`, adapters `npm run check`(prettier·eslint·tsc)·`npm test`, `prettier --check compose.yaml`, `check_compose.py`, `git diff --check` 모두 통과.

## 렌더 확인

`KNOWSLINK_UI_FIXTURE_DIR` 합성 HTML을 headless Chrome(비모바일 viewport)으로 1440·390·320에서 렌더했다. 스크립트 [render/probe.mjs](render/probe.mjs), 결과 [render/probe.json](render/probe.json), 캡처 `render/*.png` 39장. 대상 13화면: gate-query-connected/unconnected/approved, human-account-calendar-{none-google, connected, reconnect, email-member, not-ready}, 영어판 en-gate-query-connected/unconnected·en-human-account-calendar-{none, connected, reconnect}. 39개 모두 가로 넘침 0. DEV 확인이며 designer 시각 수락이 아니다.

designer 확인 필요: 이메일 코드 회원의 설정 절 한 줄 ‘Google 계정으로 로그인한 회원만 캘린더를 연결할 수 있습니다.’(coor 답변으로 기본값 채택, msg 2026-10-20). 그 밖의 운영 문구(연결 완료 안내, 해제 안내, ‘다시 연결 필요’ 설명, 준비 중 안내, 일부 동의 거부 403)는 정본의 고정 문구 밖에서 DEV가 기존 Google 로그인 문구 양식으로 썼다.

## tester/designer용 격리 fixture

`internal/relay/calendar_browser_fixture_test.go` `TestCalendarBrowserFixture`(태그 `integration && calendarfixture`). 기동은 `remote.py fixture`(자기 Postgres 컨테이너, `CALENDAR_FIXTURE_MANIFEST`).

- 계정: a 요청자(agent_a), b 연결 회원(agent_b, 봉인 토큰 저장), c 미연결 회원(agent_c). 모두 Google issuer 합성 회원이며 a–b, a–c 친구. `POST /fixture/as/{a|b|c}`로 세션(동일 origin만).
- gate(manifest `gates`): `connected_pending_2h_30`, `connected_pending_4h_15`, `connected_pending_8h_60`, `connected_expired`, `connected_processed_denied`, `unconnected_pending`. 바쁨 데이터는 window 시작+5~10분, 끝-25~-20분(비정렬)이며 제목·참석자 필드가 응답에 섞여 있다.
- 합성 Google: `/fixture/google/auth`가 ‘허용/일부만 허용/거절’ 링크를 보이는 동의 화면이다(계정은 `/fixture/as` 쿠키). `POST /fixture/google/toggle/revoke-grant`·`/outage`로 철회·장애를 켠다. `GET /fixture/gate/{id}`는 요청자 agent가 받을 state·lookup·disclosure를 보인다.
- 허용 GET: `/`, `/home`, `/home/legacy`, `/home/friends`, `/home/account`, `/home/gates/*`, `/auth/google/callback`, `/fixture/google/auth`. 허용 POST: `/home/gates/*`, `/home/calendar/connect|disconnect`, `/auth/reauth`, `/fixture/*`. 그 밖은 403. 수명 12시간. 외부 업무·메일·실계정 없음.
- 기동 근거(원격 clone `/tmp/knowslink-schedule001-dfd55004df90-47ccf514`, 컨테이너 `kl-schedule001-fixture-dfd55004df90-95e18efa`): [fixture-start](dfd55004df9067bdb8ba0941e1d5f8d2b18ac9ad-fixture-start.log)에 `calendar fixture ready; synthetic provider and accounts only`, [probe](dfd55004df9067bdb8ba0941e1d5f8d2b18ac9ad-fixture-probe.log)에 프로세스 생존, `GET /` 200, 비로그인 `/home` 303, 범위 밖 `/v1/pull` 403, `POST /fixture/as/b` 204 뒤 연결 회원 gate 200(고정 문구 1회)·`/home/account` 200(‘연결됨’ 1회)·`/fixture/gate/{id}` 200. 프로세스는 남겨 두었다(수명 12시간, coor·tester가 사용하거나 새 SHA로 다시 기동).

## 미실행(not_run)

- 실제 Google 계정 연결·실데이터 조회·Cloud Console scope 등록·운영 배포: D011 뒤 coor/ops. 합성 제공자 통과는 실계정 수락이 아니다(SC-09).
- 브라우저 키보드·200% 확대·스크린리더·실제 fixture 왕복: tester/designer 몫.
- `TestRemoteConsentBrowserReturn`: 기존 원격 실행과 같이 제외(remote consent 브라우저 검사, 무관).
- `KNOWSLINK_CALENDAR_KEY` 교체·분실 복구 절차: 미정(transition.md에 기록).
- SD-05 B(친구별 자동 허용)·commit 실행·다른 캘린더: 범위 밖.

## FullOps lint

RESULT_FULLOPS_LINT

## 의존성

go.mod·package.json 변경 없음. 봉인은 stdlib `crypto/aes`·`crypto/cipher`, 토큰 갱신은 기존 `golang.org/x/oauth2`다. DEP-001 해당 없음.

## 규모와 분할하지 않은 이유

제품 코드 약 +646/-37(새 `calendar.go` 427줄 포함), 테스트 약 +1,138. 연결·승인 조회·공개·범위·철회가 한 상태 모델과 한 승인 handler를 공유하므로 과제를 나누면 중간 상태(승인됐지만 조회 없음)가 운영에 노출된다. 지시서가 큰 규모를 예상했고 SIZE 경고는 수용했다. `http.go`는 기존 큰 파일에 좁은 델타(gate 상태·consume·approve 분기)만 더했고 새 로직은 `calendar.go`에 분리했다.
