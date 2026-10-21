---
title: SAR-PUBLISH-001-DEV DEV 보고
status: review
updated: 2026-10-21
owner: dev
tasks: [SAR-PUBLISH-001-DEV]
summary: "서버 캘린더 연결 제거(D018)와 일회 정리, /privacy·/terms 공개 페이지(D017)의 기술 계획·구현·원격 Linux 검사·격리 fixture·3폭 렌더 근거"
---

# SAR-PUBLISH-001-DEV DEV 보고

DEV 검사 기록이다. 독립 리뷰·tester QA·designer 시각 수락(PC-06)이 아니다. 정책 문안은 법률 검토(P-05) 전 초안이며 이 보고는 문안을 법적 최종본으로 보지 않는다.

## 기술 계획

- 제거 범위: SAR-SCHEDULE-001 이전 커밋 `67284e3`부터 기준 `dcf8825`까지 `internal/`·`cmd/`·`.env.example`·`compose.yaml`의 변경은 모두 캘린더 과제의 커밋이다. 그래서 해당 파일을 `67284e3` 판으로 되돌렸다. 예외는 `member_pages.go`의 `summary:not(:has(*))` 정렬 CSS 한 줄이다. 캘린더와 무관한 시각 수정이라 유지했다. 캘린더 구현·테스트·browser fixture 6개 파일(1,725줄)을 삭제했다.
- 되돌린 동작: 설정 캘린더 절, `/home/calendar/*`, 증분 동의(scope·offline), freebusy 호출, 토큰 봉인, gate lookup lease, gate-consume·`relay.result`의 disclosure, D010 범위·횟수 검사, KO/EN 캘린더 문구. Google 로그인 scope는 `openid email`이다. `schedule.query`는 U5 기준(승인은 gate만 통과, `relay.result` done은 `disclosure_denied`)이다. D010 wire 밖 검사(SR-06)는 기존 Parse 형식 검사만 남기고 범위·횟수 검사는 지웠다.
- 일회 정리: `calendar_retire.go`의 `Service.RetireCalendars`를 `cmd/relay/main.go`에서 HTTP 서비스 시작 전에 실행한다. 한 transaction에서 `State.Calendars`를 지우고, 모든 `Gate.Lookup`을 지우고, lookup이 있고 consume되지 않은 pending·approved gate를 `denied`로 끝낸다. 진행 중 캘린더 동의 attempt도 지운다. commit 뒤 키가 유효하면 각 토큰을 한 번 revoke 요청한다. 키가 없거나 잘못되면 요청 없이 삭제만 한다. 기록은 `calendar cleanup: connections=N revoked=M failed=K` 한 줄이다(N>0일 때만).
- 정리 시점 선택 이유: 필드를 바로 지우면 첫 transaction이 JSON을 다시 쓰며 토큰을 철회 없이 버린다. 그래서 legacy 필드 3개(`State.Calendars`의 봉인 값, `Gate.Lookup` raw JSON, `googleAttempt.Calendar`)를 읽기 전용으로 남기고 시작 시 정리 한 곳에서만 쓴다. 코드에 `shortcut:` 주석으로 제거 조건(운영 상태에 연결 0)을 남겼다. DB migration은 없다(상태는 `relay_state` JSON 한 행).
- 키 처리: `KNOWSLINK_CALENDAR_KEY`는 선택값이다. 시작 검증을 지워 키가 없어도 기동한다. `.env.example`·`compose.yaml`에는 'D018 일회 철회 전용'으로 남겼다. 운영 키가 있는 상태로 배포해야 철회 요청을 보낸다.
- rollback 영향: 정리 뒤 `4787e85`로 되돌리면 연결 0인 상태에서 시작한다. 토큰은 이미 철회·삭제됐으므로 재연결이 필요하다. 정리 전 rollback은 영향이 없다(정리는 제거 빌드 기동 때만 실행).
- 정책 페이지: 두 언어 문안은 run 단위 번역이 아니라 언어별 정본이다. 그래서 본문은 정적 HTML 4개(`internal/relay/policy/*.html`, embed)로 두고 h1·시행일·헤더·하단은 기존 i18n 표로 번역한다. 시행일은 `policyEffective` 상수 하나다(2026-10-22). 첫페이지의 `<style>`·하단을 `landing-style`·`landing-footer` template으로 나눠 재사용했다. 로그인 회원은 첫페이지와 같이 저장 언어를 쓴다.
- `{ }` 처리: P-01 2026-10-22, P-02 'KnowsLink 운영자'/'KnowsLink operator', P-03 표의 대체 문장, P-04 이메일 발송 업체 행 제외(운영은 SMTP 미설정으로 이메일 코드 로그인이 꺼져 있다. ops-guide·transition의 'Google 로그인 우선' 기록). 이메일 코드 로그인을 켜면 그 행을 designer·ops가 다시 정해야 한다.
- 하단 해석: 첫페이지 하단은 '개인정보처리방침 · 서비스 약관 · 설치·사용 문서 · 맨 위로'다(D04). 정책 페이지는 같은 하단에 '홈'을 서비스 약관 뒤에 더했다(SAR-PUBLISH-001 화면 기준의 '홈'). designer 시각 검수에서 확인이 필요하다.

## 완료 조건 대조

| ID | 근거 |
|---|---|
| RC-01 | 설정 절·문구 제거(unit `TestPolicyLinks`·i18n 검사, fixture probe `calendar_words=0`), `/home/calendar/*` 404(`TestCalendarRoutesRetired`, probe), Google 요청 scope `openid email`·offline 없음(`google_integration_test` 단언) |
| RC-02 | 캘린더 API 호출 코드 없음(`calendar.go` 삭제). 정리 뒤 상태 JSON에 `Calendars`·`Lookup`·`Disclosure`·`"Calendar"` 없음(`TestRetireCalendars`). 키 없이 기동(`no_key` 사례, main 시작 검증 제거) |
| RC-03 | 키 있음 성공·키 있음 철회 거부(500)·키 없음 3사례 모두 토큰 0, 미전달 결과 gate `denied`, 정리 뒤 gate-consume `invalid_gate`, 기록에 토큰·회원 ID 없음, 두 번째 시작은 요청 0(`TestRetireCalendars`). consume된 gate는 결정 유지(`TestRetireCalendarsKeepsConsumedGate`) |
| RC-04 | U5 고정 설명·승인 버튼(`TestApprovalSummaryBoundaries`, fixture probe KO/EN `u5_text=1 approve_button=1`), 승인 뒤 `disclosure:false`(`TestEmailIdentity`), `schedule.commit` 기존 동작(`TestPostgresSafety`) |
| RC-05 | 원격 integration 전체 PASS(친구·메시지·승인·원격 회귀 포함) |
| PC-01 | 비로그인 200·영어 cookie·저장 언어(`TestPolicyPagesPublic`, probe). 공개 ingress 후보에 정확 경로 추가(`check_public_ingress.py`). live 반영은 ops 몫 |
| PC-02 | 렌더 텍스트와 정본 KO/EN 4개 문안 일치([parity.txt](parity.txt), [parity.py](parity.py)), 영어판 한국어 0·`{ }` 0(`TestPolicyPages`) |
| PC-03 | 첫페이지 하단 순서·`#start` 한 줄·설정 링크·두 페이지 상호 링크(`TestPolicyLinks`, `TestPolicyPages`) |
| PC-04 | 고유 title·`aria-current` 1개·h1 1개·건너뛰기 링크·금지 요소 없음(script·img·details·code)(`TestPolicyPages`). 1440·390·320 가로 넘침 없음([widths.json](render/widths.json)). 320px 영어 위탁 표만 표 영역 안에서 스크롤(304>273) |
| PC-05 | 개인정보처리방침의 캘린더 언급은 '수집하지 않습니다' 1회뿐, 약관 0회(`TestPolicyPages`). 수집·보관 사실은 제거 뒤 코드와 같다(수집 항목표의 코드 근거는 변경 없음) |

## 검사

원격 Linux 격리 clone(자기 Postgres 컨테이너, [remote.py](remote.py))에서 실행했다. 최종 제품 SHA는 `e20c199`(clone `/tmp/knowslink-publish001-e20c199a4ec2-caeb11f3`)다. fixture는 그 전 `d54aa68`에서 기동했다. 두 SHA 사이의 제품 차이는 adapter 되돌림(아래)과 env 변수로만 켜지는 테스트 HTML 출력뿐이며 Go relay 동작은 같다. 로그 SHA256은 [raw-log-sha256.txt](raw-log-sha256.txt)다(작업 트리의 CRLF 사본 기준).

| 검사 | 결과 | 로그 |
|---|---|---|
| make lint | `e20c199` EXIT=0 (`d54aa68` EXIT=0) | [lint](e20c199a4ec2d81c434962954614a8f0feba2b85-lint.log) |
| make test | `e20c199` EXIT=0 (`d54aa68` EXIT=0) | [test](e20c199a4ec2d81c434962954614a8f0feba2b85-test.log) |
| integration(`-race -tags=integration`, `TestRemoteConsentBrowserReturn` 제외: 기존과 같음) | `e20c199` EXIT=0, FAIL 0 (`d54aa68` EXIT=0) | [integration](e20c199a4ec2d81c434962954614a8f0feba2b85-integration.log), [d54aa68](d54aa68ab2e4e777cd73dbbf2faebddf9f4b0966-integration.log) |
| fixture 기동·probe | `d54aa68` EXIT=0 | [start](d54aa68ab2e4e777cd73dbbf2faebddf9f4b0966-fixture-start.log), [probe](d54aa68ab2e4e777cd73dbbf2faebddf9f4b0966-fixture-probe.log) |

실패 기록: 첫 integration(`cbe1c05`)은 새 정리 테스트의 idempotency key가 15자(최소 16자)라 실패했다. `d54aa68`에서 고쳤다. 같은 SHA의 integration을 make test 전에 돌린 실행은 `adapters/dist`가 없어 `TestPublicNodeProcesses`가 실패했다. make lint·make test 뒤 다시 실행해 EXIT=0이다(dev 컨텍스트의 기존 교훈). 문서 하위 작업의 대조에서 SAR-SCHEDULE-001이 `adapters/src/core.ts`·`synthetic.ts`도 바꾼 것을 찾았다. 처음 되돌림 범위(`internal`·`cmd`)에서 빠졌던 파일이다. `67284e3` 판으로 되돌려 adapter가 승인 뒤 다시 `denied`만 보내고 합성 조회는 고정 기간을 쓴다. 그 뒤 `e20c199`에서 전체 게이트를 다시 실행했다.

로컬 Windows: `go vet ./...`, `go vet -tags integration[,publishfixture]`, `go test ./internal/relay ./cmd/...` 통과. gofmt는 toolchain gofmt(go1.27.1)로 확인했다.

FullOps lint(`--from dcf8825`, 깨끗한 임시 detached worktree, HEAD `681bacf`): ERROR 0, WARNING 17, 실행 불가 2([로그](681bacf-fullops-lint.log)). 실행 불가는 Windows에 make가 없는 product-lint·product-test이며 같은 명령을 위 원격 Linux에서 EXIT=0으로 실행했다. 경고 처리:

- ANTI-005(캘린더 테스트 5개 파일·`human_pages_test.go` 1사례 삭제): 삭제된 기능의 테스트다. 대체는 `calendar_retire_integration_test.go`(정리·경로 404)와 U5로 되돌린 기존 승인 테스트다.
- SIZE-002(추가 1,552줄)·SIZE-001(`i18n_en_member.go`·`member.go`·`ops-guide.md`·`PLANS.md`·`service-member-ui.md`): 제거 단위와 정책 문안 4개(약 300줄) 때문이다. 기존 파일은 몇 줄만 늘었고 분할은 이번 범위 밖이다. `PLANS.md`·`service-member-ui.md`는 이 과제가 바꾸지 않았다.
- DOC-001(`policy/*.html`): 방문자에게 그대로 보내는 정적 문서라 헤더 주석을 넣지 않았다. 설명은 `policy_pages.go` 헤더에 있다.
- SLOP-004(`check_public_ingress.py:24`): 검사 결과를 알리는 기존 출력 줄이다.

## 격리 fixture (tester·designer)

`TestPublishBrowserFixture`(빌드 태그 `integration,publishfixture`)는 합성 DB에서 SAR-SCHEDULE-001 빌드가 저장한 모양의 정리 전 상태(연결 2개, 미전달 결과가 있는 approved gate 1개, pending gate 1개)로 시작한다. 실행은 `python3 remote.py fixture`이며 manifest(`<base>/fixture.json`)에 URL과 gate ID가 있다. 현재 기동 위치는 [fixture-location.json](fixture-location.json)이고 수명은 24시간이다.

- `GET /fixture/state`: 저장 상태 요약(연결 수, `Calendars`·`Lookup`·`Disclosure` 잔존 여부, 두 gate 상태, 로컬 revoke 요청 수). 토큰·이메일 값은 내지 않는다.
- `POST /fixture/retire?key=present|absent[&revoke=fail]`: 시작 정리와 같은 함수를 로컬 revoke 끝점으로 실행하고 정리 기록 한 줄을 돌려준다. `POST /fixture/seed`는 정리 전 상태를 다시 쓴다.
- `POST /fixture/as/ko|en`: 합성 회원 '밥'으로 로그인(저장 언어 지정) 뒤 `/home/account`로 이동한다. `/home/gates/{id}`에서 승인 상세를 본다. `POST /fixture/consume?gate=<id>`는 agent 쪽 gate-consume 응답을 보여 준다.
- 정책 페이지: `GET /privacy`, `GET /terms`, `/lang/en?back=/privacy`.
- 범위 밖 경로는 403, `POST /fixture/stop`으로 끝낸다. 외부 Google·메일 호출은 없다.

## 화면 근거

[render/](render/)에 KO/EN의 `/privacy`·`/terms`·첫페이지(`policy-start-*`)·설정(`human-account`)·승인 상세(`gate-active`)를 1440·390·320px로 찍은 PNG가 있다. 측정은 headless Chrome의 고정 폭 iframe이다(창 최소폭 500px 제약 회피, [measure.py](render/measure.py)·[shots.py](render/shots.py)). 원본은 단위 테스트의 `KNOWSLINK_UI_FIXTURE_DIR` 출력이다. 이 캡처는 designer 시각 수락이 아니다.

## 영향과 호환

- 크기: 제품 코드(테스트·정적 문안 제외) +232/−673, 전체 `internal`·`cmd` +898/−2,052. 제거 단위라 분할하지 않았다. SIZE 경고는 이 이유로 처리한다.
- 의존성 변경 없음. `golang.org/x/oauth2`·`go-oidc`는 로그인에 계속 쓴다. revoke는 기존 `calendar.go`의 표준 라이브러리 경로(`net/http` POST)를 재사용해 Context7 조회가 필요 없었다.
- `/home/calendar/*`를 호출하는 오래된 화면이 있으면 404를 받는다. 기존 클라이언트의 `schedule.query`·`schedule.commit` 전송은 그대로 받는다.

## 검증하지 못한 것

- 실제 Google revoke 끝점 호출(실계정 금지). 로컬 끝점으로 요청 형식·횟수만 확인했다.
- 운영 상태의 실제 정리 결과·live Tunnel 반영·Cloud Console 정리·앱 게시(ops·coor 몫).
- 실제 브라우저의 키보드·200% 확대 검사, designer 직접 시각 수락, tester 독립 QA.
