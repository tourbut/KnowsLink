---
title: SAR-PUBLIC-AGENTS-001-FIX-TESTER — 보존·rate·관계 조건 QA
status: draft
updated: 2026-10-06
owner: tester
tasks: [SAR-PUBLIC-AGENTS-001-FIX-TESTER]
summary: "후보 458798c의 보존, rate, POLICY 관찰 조건은 통과했고 원본 QA와 UI 실패는 그 SHA에 둔다"
---

# SAR-PUBLIC-AGENTS-001-FIX-TESTER — 보존·rate·관계 조건 QA

## 판정

판정 후보는 `458798c2ee15c179edacfd6f94ebb9896d26f411`다. 실행 위치는 `/tmp/knowslink-agents-fix-qa-458798c`의 detached checkout이다. 시작과 끝의 HEAD는 같고 추적 트리는 비어 있다.
좁은 변경 영향 QA는 통과다. 새 critical/high는 없다. 제품 코드는 수정하지 않았다.
최종 프로브 종료코드는 0이다. `TestQAAgentsFix`는 3.62초이고 policy, retention, rate가 모두 PASS다.
이 통과는 원본 QA와 원본 UI 실패를 PASS로 바꾸지 않는다. 실제 이메일, 운영 공개, 운영 자원, 노우↔다닷은 미실행이다. 벽시계 24시간은 기다리지 않았다.

## 기준

날짜는 2026-10-06이다. 공통 기준은 `fullops-common-0.3.3`과 coding-style, testing, security, `project.md`다.
관찰 조건은 [POLICY](../../exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md)의 두 표다. 과거 `d1eef9bb90b9726149980320c42fb1fdbcaf584a`가 그 표를 통과했다고 기록하지 않는다.
원본 QA는 `bcb06b89bcb36d69a99cbeef3d94e4a9ffe88361`다. 원본 UI는 `d1651784c4338efeb0d6141467d563c6b354e4a5`이고 UX04–05는 FAIL/보류다. F-UI-01과 F-UI-02는 그대로다.
도구는 Go 1.27.1과 격리 Postgres 17이다. 회원 메일은 메모리 inbox다. Workers Free 밖의 유료 호출은 하지 않았다.

## 제품 동일성

후보는 `d1eef9bb90b9726149980320c42fb1fdbcaf584a`, `70f26bc0799e65e4647731612a8d3a7c098a5fec`, `bcb06b89bcb36d69a99cbeef3d94e4a9ffe88361`, `d1651784c4338efeb0d6141467d563c6b354e4a5`, `4a1b80aec8fa6a06144d51f3a5609927a2644928`, `48d12fae2dce35d92606b264313148f0a635b64e`, `83e0bfb907085bade31a193ab91ca723feb7d7ad`, `6d016e5979757e51aa8d0f65f42f460e9f15ad6b`의 자손이다.
`6d016e5979757e51aa8d0f65f42f460e9f15ad6b`와 후보의 `internal`, `cmd`, `adapters` diff는 비어 있다. 기록 체크아웃과 후보의 같은 경로, Makefile, `go.mod`, `go.sum` diff도 비어 있다.
`d1eef9bb90b9726149980320c42fb1fdbcaf584a` 이후 제품 차이는 relay 9개 파일이다. `api_rate.go`, `connections.go`, `connections_integration_test.go`, `connections_test.go`, `http.go`, `member.go`, `member_agents.go`, `policy_test.go`, `store.go`다. 삽입 842줄, 삭제 89줄이다.
`adapters`, `cmd`, `identity.go`, protocol, compose, Makefile은 그 diff에 없다. 원본 CLI와 신원, trial 증거는 `d1eef9bb90b9726149980320c42fb1fdbcaf584a`와 `bcb06b89bcb36d69a99cbeef3d94e4a9ffe88361`에 둔다.

## 실행

프로브 파일은 실행 중에만 clone의 `internal/relay/qa_agents_fix_probe_test.go`로 복사했다. 종료 후 그 파일을 지웠다. clone HEAD와 추적 트리는 실행 전과 같다.
Compose project 이름은 `knowslink-agentsfix-`와 임의의 8자리다. Postgres는 loopback의 빈 포트만 연다. `up`, migrate 빌드, migrate, `down --volumes`의 종료코드는 0이다.
Go 명령은 `go test -tags=integration -race -count=1 -timeout 12m -v -run TestQAAgentsFix ./internal/relay`다. 종료코드는 파이프로 가리지 않았다.
`docker ps -a`의 실행 전후 목록은 같다. 위 세 컨테이너 ID는 그대로다. 이 프로브의 project 컨테이너는 목록에 없다.
로그는 [SAR-PUBLIC-AGENTS-001-FIX-TESTER-test/](SAR-PUBLIC-AGENTS-001-FIX-TESTER-test/)에 있다. DB URL과 43자 이상의 토큰은 `[redacted]`다.

## 1차 프로브 실패

1차 `probe-go.log` 종료코드는 1이다. 그 로그는 `probe-go-attempt1.log`다. 제품 코드를 고치지 않았다.
policy는 발신 홈에서 `받은 초대입니다`를 찾았다. 그 화면의 안내 문장은 `받은 초대라면`이고, 관계 문장은 수신 owner의 수락을 기다리라는 문장이다. 수신 홈의 `받은 초대입니다`는 별도 요청에서 확인한다.
retention은 키 상한 거부가 연결 수를 0으로 만든다고 기대했다. 앞선 20회의 성공한 연결은 consumed로 남아 있다. 열린 연결이 늘지 않는지로 다시 확인한다.
2차 실행의 종료코드는 0이다. 1차 실패는 프로브 기대값이다.

## 확인한 동작

P1부터 P7은 `TestQAAgentsFix/policy`다. 같은 방향과 반대 방향의 pending 반복은 세대 1과 같은 Exp를 유지했다. `http:member:` 기록은 2개 이상이다.
한 시간 뒤의 in-process 초대도 세대와 Exp를 유지했다. 이 시각은 벽시계 대기가 아니다.
Exp를 DB 시각보다 1분 전으로 옮긴 수락은 403이고 상태는 expired다. 수동 재초대는 세대 2다. 옛 세대 수락은 403이다.
Exp를 2분 뒤로 둔 수락은 303이고 active다. 같은 영수증의 재전송은 200이고 메시지는 1개다. active 반복은 세대 2를 유지한다.
발신 unpair와 수신 unpair를 각각 실행했다. 철회 뒤 send는 403이다. 종료 홈에는 `연결이 끝났습니다`, `새 초대 보내기`, `상대가 새로 수락해야 합니다`가 있고 `action="/home/unpair"`는 없다.
새 수락 전 옛 영수증 재전송은 403이다. 새 idempotency key의 send는 200이다.
deny는 pending 수를 0으로 만든다. 늦은 수락은 403이다. 수락과 거절의 동시 요청은 303 하나와 403 하나다. 세대는 그대로다.
같은 pending의 수락 8회는 모두 303이고 active는 하나다. 같은 owner의 초대는 자동 수락되지 않는다. 다른 owner의 수락은 403이다.
송신 pending 10, 전체 pending 200, owner active 20, 전체 active 400의 다음 시도는 409이고 세대를 만들지 않는다. `http:new`가 찬 초대는 429다. 그 상태의 deny는 303이다. 정리 예산 20을 채운 unpair는 429이고 상태는 pending이다.
새 Service의 홈은 같은 세대를 보인다. 재로그인의 회원 ID는 같고 관계는 pending으로 남는다. 철회한 agent와 비활성 owner의 재초대는 403이고 세대는 늘지 않는다.

M1부터 M3는 `TestQAAgentsFix/retention`다. 회전으로 철회된 key1의 prepare는 409 `key_exists`다. 키 기록은 20개다.
다음 rotate는 409다. 화면에는 `기다려도 이 agent에 새 키 공간은 생기지 않습니다`, `action="/home/agents"`, `href="/home"`가 있다. 키 수는 그대로이고 열린 연결은 0이다.
홈에는 `새 키를 더 연결할 수 없습니다`와 `새 agent를 만들어 따로 연결`이 있다. `결제`, `업그레이드`, `요금`, `quota`는 없다.
key2 철회는 303이다. 키 수는 20개로 남는다. Changed를 30일 전으로 옮겨도 key1은 남는다.
활성 agent 5개의 다음 생성은 409 `회원당 5개`다. 결제 문구는 없다.
Changed가 비어 있는 철회 agent는 다음 GET에서 남고 Changed가 채워진다. 홈은 `최소 24시간 보존`, `목록에서 사라질 수 있습니다`, `권한이 돌아오거나 백업까지 영구 삭제된 것은 아닙니다`를 보인다.
Changed를 25시간 전으로 옮기면 그 agent는 상태와 홈에서 사라진다. 이 이동은 벽시계 24시간이 아니다.
연결된 agent를 철회하면 pull은 401이다. owner 기록이 10개일 때 생성은 409 `철회한 agent 기록을 보존하는 중`, `최소 24시간`, `자기 홈으로 돌아가기`다. 기록 수는 그대로다.
그 상한에서 철회는 303이고 기록 수는 그대로다. 아직 철회하지 않은 agent의 pull은 200이다.
철회 시각을 25시간 전으로 옮긴 GET은 그 기록을 지운다. 지운 agent의 pull은 401이다. 수동 생성은 303이다.
새 agent의 pair는 없다. 수락 전 send는 403이다. 수신 owner의 수락은 303이고 그 다음 send는 200이다. 기존 관계의 세대와 상태는 바뀌지 않는다. 옛 key1 pull은 401이고 마지막 키의 pull은 200이다.

R1부터 R3는 `TestQAAgentsFix/rate`다. `invalid_schema`를 반환하는 transaction은 phantom agent와 바꾼 owner를 저장하지 않는다. 기한이 지난 철회 agent는 지워진다. 기한이 남은 철회 agent는 남고, 시각이 비어 있던 기록은 Changed가 채워진다. 익명 rate 기록은 1개다.
`%` 본문 40회는 422다. agent 수는 그대로다. 회원 rate 기록은 40개 이상이다. 41번째는 429 `이후 다시 시도`다. 다른 회원의 생성은 303이다. 새 Service에서 같은 회원의 생성은 429다.
잘못된 세션의 GET 30회는 303 `/?n=expired`다. `/home`과 `/home/connections/unknown`을 섞었다. 31번째와 새 Service의 GET, logout, reauth는 429다. 다른 IP의 첫 GET은 303이다.
`{` prepare 30회는 422다. 31번째는 429 `rate_limited`다. `http.ParseTime`으로 읽은 `Retry-After`와 RFC3339 `retry_at`은 같다. 그 시각은 현재 근처의 미래이고 90초 안이다. 바로 다음 요청은 429다. 다른 IP의 `{`는 422다.

## 안내 존재

위 HTTP 응답에서 반복 안내, 종료 뒤 수동 초대, 키 기록 포화의 새 agent 안내, owner 기록 포화의 24시간 안내, 철회 목록 안내, 거부 화면의 `자기 홈으로 돌아가기`와 `action="/home/agents"`를 확인했다.
상태코드는 용량과 `key_exists` 409, 권한 실패 403, 잘못된 본문 422, rate 429다. 이 확인은 문자열과 상태코드다. 모바일 배치와 F-UI-01, F-UI-02의 시각 판정은 designer가 한다.

## 재사용한 증거

원본 프로브의 최종 종료코드 0과 `probe-failures.md`의 기대값 불일치는 `d1eef9bb90b9726149980320c42fb1fdbcaf584a`에 둔다. 그 프로브를 다시 실행하지 않았다.
`adapters`와 `cmd`가 후보 diff에 없으므로 원본 Node CLI 결과를 다시 실행하지 않았다. 연결, 관계, 회원 HTTP, rate는 코드가 바뀌어 이 프로브로 다시 확인했다.
원본 `make lint`, `make test`, `make verify-mvp`의 종료코드 0은 `bcb06b89bcb36d69a99cbeef3d94e4a9ffe88361`의 결과다. 이 과제에서 `make verify-mvp`는 실행하지 않았다.

## 미실행

D12 운영 수락, 실제 메일, 공개 서버, 운영 데이터 정리, 복원, 노우↔다닷은 미실행이다. 계정 비활성화는 홈에서 `준비 중`이다.
후보 동작이 `link.knowslog.com`에 배포되었다고 기록하지 않는다. 시각 테마 검사와 화면 캡처는 만들지 않았다.

## 기록 검사

LINT_PENDING
