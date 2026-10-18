---
title: SAR-UX-ISSUES-002-NAMES-DEV DEV 검증 보고
status: draft
updated: 2026-10-19
owner: dev
tasks: [SAR-UX-ISSUES-002-NAMES-DEV]
summary: "이름 중심 표시·친구 코드 공유·친구 요청 복구(U2/U4+U7)의 DEV 검사, 렌더 확인, fixture와 미실행 항목"
---

# SAR-UX-ISSUES-002-NAMES-DEV DEV 검증 보고

DEV 자동 검사와 DEV 렌더 확인의 기록이다. 독립 tester QA·designer 시각 수락·고정 SHA 리뷰·main 통합은 아니다.

## 대상과 기준

- 브랜치 `fullops/dev`. 첫 후보 `3c1567fcff8c0dcb4693387359ad0fa92fa69bff`, 리뷰 low 수정 뒤 최종 제품 후보 `516bf6247f06843c97b0292d0db01961843713ea`. 기준 `8d0db4bf59e905c9f398123c4e2e5bc13ed4c3fa`(main), 지시서 커밋 `6275431`.
- 516bf62 이후 커밋은 이 보고·packet-outcomes·완료 기록뿐이며 제품 코드는 같다. 최종 HEAD는 완료 보고에 적는다.
- 규칙 fullops-common-0.3.3, 제품 정본 SAR-UX-ISSUES-002.md U2/U4/U7, D02 HP-02–04·13·14, D04.

## 등록 검사 (원격 Linux 격리 clone)

[remote.py](remote.py)가 exact HEAD를 bundle로 고유 `/tmp` clone에 옮긴다. known_hosts RejectPolicy를 쓰고 자격 값은 출력하지 않는다. 원 로그 SHA256은 [raw-log-sha256.txt](raw-log-sha256.txt)다. 아래 로그는 줄 끝 공백만 정리한 사본이다.

| 검사 | 명령 | 516bf62 (최종) | 3c1567f (첫 후보) |
|---|---|---|---|
| product-lint | `make lint` | EXIT=0, [로그](516bf6247f06843c97b0292d0db01961843713ea-lint.log) | EXIT=0, [로그](3c1567fcff8c0dcb4693387359ad0fa92fa69bff-lint.log) |
| product-test | `make test` (go test -race, npm, python) | EXIT=0, [로그](516bf6247f06843c97b0292d0db01961843713ea-test.log) | EXIT=0, [로그](3c1567fcff8c0dcb4693387359ad0fa92fa69bff-test.log) |
| integration | 격리 Postgres 17 컨테이너, `go test -race -tags=integration ./internal/relay -skip TestRemoteConsentBrowserReturn -count=1 -v` | EXIT=0, 122.4초, [로그](516bf6247f06843c97b0292d0db01961843713ea-integration.log), 정지 EXIT=0 | EXIT=0, 117.9초, [로그](3c1567fcff8c0dcb4693387359ad0fa92fa69bff-integration.log), 정지 EXIT=0 |

lint 로그의 `ERROR: invalid tunnel UUID`는 셸 스크립트 테스트의 기대 출력이다. 각 SHA는 새 고유 clone에서 따로 실행했다.

`TestRemoteConsentBrowserReturn`은 `TEST_BROWSER`(Chromium 실행 파일)가 필요해 제외했다(not_run). 이번 변경은 remote 동의 화면을 바꾸지 않았다.

Windows 로컬: `gofmt -l` 출력 없음, `go vet`(기본·`integration receiptfixture`·`integration namesfixture` 태그) 통과, `go test ./internal/relay/` 통과.

## 완료 조건과 실패 경계

| 조건 | 근거 |
|---|---|
| U2 이름 먼저 + 종류 + 구분값 | 조회 radio(`TestReceiptFormRecoveryFixtures`), 승인(`TestGateNamesEndpointsWithinAuthority`), legacy·연결(golden·렌더 fixture) |
| 구분값 충돌 확장 | `TestReceiptSelectorDistinguishesCollidingNames`, `TestLegacyAgentLabelsWidenCollidingDiscriminators` |
| 권한 밖 이름 비노출 | `TestVisibleAgentNamesOnlyOwnerAndCurrentFriend`: 타인·없는 agent·차단 뒤 친구는 ‘확인할 수 없는 에이전트’+ID |
| 40자 한글·이모지 허용, 제어문자 거절, 오류 재시도로 agent 중복 없음 | `TestHumanNamesAndFriendCodePostgres`, 기존 `TestHumanNamesAndIdentity` |
| 이름 변경 뒤 ID·관계·세대 불변 | `TestHumanNamesAndFriendCodePostgres` |
| 세 생성 경로 이름 결속 | 수동: 위 integration. device: `TestGoogleDeviceHTTP` 추가 단언(이름=클라이언트 agent, 422 재시도로 중복 없음). remote: `TestRemoteOAuthHTTP` 추가 단언 |
| 내 친구 코드 읽기 전용·라벨·미제출 | `TestFriendPagesShareOwnCodeAndKeepRequestInput`, integration의 `/home/account`·`/home/friends` |
| 앞뒤 공백만 제거, 빈 입력·동의 누락 422와 입력 유지 | integration(공백 코드 422·동의 누락 422·패딩 코드 성공), `TestFriendCodeRotationKeepsReceivedRequests`(매처는 정확 비교) |
| 코드 교체: 옛 코드 발견 불가, 기존 친구·받은 요청 유지 | integration, `TestFriendCodeRotationKeepsReceivedRequests` |
| GET만으로 관계 생성 없음 | `TestFriendViewsNeverCreateRelationships`, integration |
| agent 0개 요청·수락, 빈 발신 선택기 없음 | `TestFriendViewsNeverCreateRelationships`, integration(친구 상세에 대화 링크 없음·‘첫 에이전트를 추가하세요’) |
| 24h 경합·옛 세대 결정 | `TestFriendViewsNeverCreateRelationships`(기한 시각 수락 거부·재요청 새 세대), 기존 `TestHumanFriendTransitionsAndDiscovery` |
| 차단·존재 비노출·반복 요청 기한 불변 | 기존 `TestHumanFriendTransitionsAndDiscovery`, `TestHumanFriendPrivacyAndCapacity` |
| 익명 429 오안내 제거 | `member_receipt.go`에 `Unavailable` 추가. 별도 자동 검사 없음(한 줄 변경) |

## 리뷰 하위 에이전트

읽기 전용 리뷰 하나(같은 세션의 하위 에이전트, 독립 리뷰 아님)가 3c1567f를 검토했다. critical·high·medium은 없었다. low 4건을 516bf62에서 고쳤다: legacy 관계 행의 충돌 구분값 재사용, 친구 상세 429의 친구 목록 복귀, 동의 체크박스의 `aria-describedby`에 오류 연결, 이름을 볼 수 없는 agent 라벨의 ID 중복 제거. 남은 지적은 gate/legacy가 `visibleAgent`를 부르는 배선의 직접 테스트 부재와 라벨당 회원 순회 비용(MVP 규모에서 수용)이다.

## 렌더 확인 (DEV, 합성 fixture)

`KNOWSLINK_UI_FIXTURE_DIR`로 만든 HTML을 headless Chrome 155(CDP, 비모바일 viewport)로 1440·390·320px에서 측정했다. 스크립트 [render/probe.mjs](render/probe.mjs), 결과 [render/probe.json](render/probe.json), 캡처 `render/*.png`.

- 대상: 조회(첫·철회·성공), 40자 이름·43자 코드로 바꾼 친구·내 계정·친구 상세·내 에이전트·새 대화, legacy, 승인, 연결, 거부.
- 페이지 가로 넘침 0, 44px 미만 동작 0. 조회 radio 라벨은 39자 이름·종류·구분값이 줄바꿈되어 모두 보인다(D-RT-01 해소).
- 처음 실행에서 320px 친구 동의 라벨의 40자 공백 없는 이름이 넘쳤다. `.consent span`에 `overflow-wrap:anywhere`를 더해 고쳤다. 모바일 emulation은 넘친 만큼 viewport를 넓혀 이 넘침을 숨겼으므로 비모바일 viewport로 다시 측정했다.
- 캡처는 3c1567f 기준이다. 516bf62는 이름을 볼 수 없는 agent 라벨에서 ‘/ ID’를 빼고 상세의 ID만 남겼다. CSS는 같다.
- 남은 표시: 편집 입력(내 이름 40자, 요청 ID UUID)은 입력 안에서 가로 스크롤된다. 기존 동작이다.

not_run: 실제 200% 확대, 스크린리더, Tab/Shift+Tab 실제 키 입력, 두 번째 브라우저, 실제 인증·DB 화면 캡처. 이들은 tester/designer 후속이다.

## tester/designer용 격리 fixture

`internal/relay/names_browser_fixture_test.go`(build tag `integration && namesfixture`). 실행: `remote.py fixture`(prepare 뒤). 내부는 `go test -c -tags=integration,namesfixture`와 `NAMES_FIXTURE_MANIFEST=<private path> <binary> -test.run=^TestNamesBrowserFixture$ -test.timeout=25h -test.v`다. `TEST_SYNTHETIC_DATABASE=1`과 고유 Postgres DB가 필수다. 운영 DB에는 실행하지 않는다.

- 합성 계정: A(이름 40자 ‘가’, 같은 39자 이름 agent 두 개, legacy agent와의 정확한 pair), B(이모지 40자 이름 agent 하나, A와 친구), C(agent 0개, A에게 보낸 대기 요청).
- origin: `httptest` loopback(`http://127.0.0.1:<임의 포트>`), manifest에 `origin`·`friend_ab`·`request_ca`를 0600으로 기록.
- 계정 전환: `POST /fixture/as/{a|b|c}`(동일 출처·non-cross-site만).
- 허용 GET: `/`, `/home`, `/home/friends`, `/home/friends/{id}`, `/home/account`, `/home/agents`, `/home/legacy`, `/home/receipts`, `/home/new`. 허용 POST: `/home/friends/request`, `/home/friends/decision`, `/home/friend-code`, `/home/agent-name`. 그 밖은 403.
- 수명: 24시간 또는 프로세스 종료. 메시지 발송·외부 인증·실메일은 없다.
- 이번 DEV는 fixture를 컴파일(`go vet` 태그)만 했고 기동하지 않았다(not_run). tester가 고정 SHA에서 기동한다.

## 1단위 잔여 처리

- 조회 선택기 종류 잘림·순서: select를 radio fieldset으로 바꾸고 `이름 / 종류 / 구분값`으로 통일했다.
- tester `checks.py`의 구 순서 단언(`opt.endswith("agent_a")`)과 선택기 폭 검사는 tester 소유라 수정하지 않았다. 새 계약은 `<input type="radio" name="agent" value="…" required checked>`와 라벨 `<span>이름 / 종류 / 구분값</span>`이다. integration Go 단언 5건은 새 계약으로 갱신했다.
- 429 문구: 위 표.

## 남은 위험과 판단 필요

- 보낸 요청이 거절되면 요청자 기록에 ‘거절’이 보인다. 없는 코드의 placeholder는 만료만 된다. 따라서 거절 표시는 코드가 실제 회원임을 사후에 알릴 수 있다. HP-04의 상태 구분 요구와 U4 비노출 사이의 기존 동작이며 이번 범위에서 바꾸지 않았다. designer 판단 대상이다.
- 새 대화의 발신/수신 select는 이름이 먼저이지만 40자 이름이면 닫힌 상태에서 종류·구분값이 잘릴 수 있다. 이번 지시 범위(조회 선택기)만 radio로 바꿨다.
- legacy 화면 제목은 ‘내 에이전트’에서 ‘기존 연결’로 바뀌었다(고유 title, U7).
