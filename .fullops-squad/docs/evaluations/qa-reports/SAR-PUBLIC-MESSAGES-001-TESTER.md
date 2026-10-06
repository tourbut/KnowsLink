---
title: SAR-PUBLIC-MESSAGES-001-TESTER — PS08–11 독립 QA
status: draft
updated: 2026-10-06
owner: tester
tasks: [SAR-PUBLIC-MESSAGES-001-TESTER]
summary: 후보 09c523d의 PS08–11은 H-1 high로 수락이 차단되고 M-1과 HTTP 슬롯 잔류도 재현됐다
---

# SAR-PUBLIC-MESSAGES-001-TESTER — PS08–11 독립 QA

## 판정

판정 후보는 `09c523da8a3407288d9f5d711e1834af12bc7808`다. 실행 위치는 `/tmp/knowslink-messages-qa-09c523d`의 detached checkout이다. 시작과 끝의 HEAD는 같고 추적 트리는 비어 있다.
독립 QA는 완료다. 제품 수락은 차단이다. 미해결 H-1은 high다.
최종 프로브 `TestQAMessagesIndependent`의 종료코드는 1이다. 17차 실행은 23.20초다. 실패 줄은 보존한 `t.Error` 네 개다. 그 뒤에 `t.Fatal`은 없다.
`make lint`, `make test`, `make verify-mvp`의 종료코드는 0이다. 이 세 명령은 제품 스위트다. 운영 PS08, PS13, PS14의 통과로 쓰지 않는다.
제품 코드는 수정하지 않았다. 로컬 fixture의 왕복은 격리된 Node, MCP, Postgres 결과다. 실메일, 공개 배포, 벽시계 24시간, 운영 부하와 복원, 외부 플랫폼, 최종 노우↔다닷은 미실행이다.

## 기준

날짜는 2026-10-06이다. 공통 기준은 `fullops-common-0.3.3`과 coding-style, testing, security, `project.md`다.
제품 기준은 `docs/planning/product-specs/SAR-PUBLIC-SERVICE.md`의 PS08–11과 PS04, PS06, PS07, 그리고 `docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md`의 UX06, UX07이다.
DEV 원본 실패는 [실행 기록](../../exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV.md)에 둔다. 그 기록의 최종 자동 검사 종료코드 0을 이 QA의 통과로 바꾸지 않는다.
도구는 Go 1.27.1, Node v22.22.2, Postgres 17, Docker 29.4.3이다. 회원 메일은 메모리 inbox다. Workers Free 밖의 유료 호출은 하지 않았다.
실제 경로는 `internal/relay/member_receipt.go`다. context의 `member_receipts.go` 거부는 오기다. 그 JSON은 수정하지 않았다.
읽은 제품 파일은 `public_text.go`, `capacity.go`, `http.go`, `store.go`, `member_receipt.go`, `member.go`, `adapters/src/text.ts`, `adapters/README.md`다.

## 실행

프로브 원본은 [probe_test.go.src](SAR-PUBLIC-MESSAGES-001-TESTER-test/probe_test.go.src)다. 실행 중에만 clone의 `internal/relay/qa_messages_probe_test.go`로 복사했다. 종료 후 그 파일을 지웠다.
Go 명령은 `go test -tags=integration -race -count=1 -timeout 15m -v -run TestQAMessagesIndependent ./internal/relay`다. 종료코드는 파이프로 가리지 않았다.
Compose project 이름은 `knowslink-messagesqa-`와 임의의 8자리다. Postgres는 loopback의 빈 포트만 연다. 17차의 `up`, migrate 빌드, migrate, `down --volumes` 종료코드는 0이다.
자격은 이메일 코드 로그인으로 만든 일반 회원이다. 시험 allowlist와 공유 ServiceAuth를 일반 자격으로 쓰지 않았다. 회원마다 agent와 key credential이 다르다.
Node 왕복은 clone의 `adapters`에서 `connect.js`와 MCP를 실제 자식 프로세스로 띄운다. 17차 ID는 요청 `01a11167-8255-7ac7-991d-67b6c0a95d9f`, 회신 `01a11167-836e-7c3c-86c9-7ec4d996ac20`, 4096바이트 `01a11167-842d-78ac-a430-b3b4cd39b4a8`다. `actualGrok`, `actualDadot`, `realMail`은 `unverified`다.
`make verify-mvp`는 별도 project `knowslink-mvp-2862a643cf`를 만들고 끝에서 volume과 network를 지웠다. 로그의 마지막 줄은 `PASS: isolated Compose migration, Postgres races, TS adapter and Go owner UI; Tunnel unused`다.
공유 컨테이너 `knowslink-relay-1` `dfcd9d187117`, `knowslink-cloudflared-1` `07077b9ef5e4`, `knowslink-postgres-1` `bc3482dc52f2`는 멈추지 않았다. 같은 시각에 보인 `sar-messages-ui-fix-vnsxs7mb`와 `knowslink-mvp-70696d19ab-*`는 이 QA가 만든 project가 아니다. 그 컨테이너는 중지하지 않았다.
DB URL과 PEM은 프로브 로그에서 `[redacted]`다. 증거는 [SAR-PUBLIC-MESSAGES-001-TESTER-test/](SAR-PUBLIC-MESSAGES-001-TESTER-test/)에 있다.

## 프로브 기대값 수정

1차부터 16차의 중간 실패 중 제품 결함이 아닌 항목은 프로브 기대값이다. 제품 코드는 고치지 않았다. 원본 로그는 `probe-go-attempt8.log`부터 `probe-go-attempt16.log`다. 요약은 [probe-failures.md](SAR-PUBLIC-MESSAGES-001-TESTER-test/probe-failures.md)다.
DB `clock_timestamp()`로 TTL 180초와 181초를 나눈 뒤 200과 422가 맞았다. filler의 `Accepted`가 0이면 다음 transaction이 24시간 sweep로 지운다. `Accepted`를 현재 시각으로 둔 뒤 receipt 20000의 다음은 409다.
ACK는 유효한 lease와 앞선 persist가 필요하다. persist를 신규 슬롯이 비어 있을 때 실행한 뒤 슬롯을 채우면 ACK는 200이다.
slow body는 `Read`가 release까지 멈추는 reader와 100ms 대기가 있어야 프로세스 슬롯을 먼저 점유한다. 429 본문이 `rate_limited`이면 슬롯 점유로 기록하지 않는다. `capacity`만 결함으로 남긴다.
unpair 뒤 sweep는 현재 세대가 아닌 gate를 `revoked`로 바꾼다. deny 직후에 gate가 `denied`인지 확인하고, 그 다음 unpair를 본다.
신규 DB 슬롯 16개가 남아 있는 GET `/home/receipts`는 인증 전에 429다. 화면 문장은 `동시 처리 한도에 도달했습니다`다. 명세가 포화 중에도 요구하는 동작은 ACK, deny, 철회, unpair, 로그아웃이다. GET을 그 예산의 예외로 두지 않았다.

## 확인한 동작

로컬 Node와 MCP가 text를 보내고, 상대가 pull, persist, ACK한 뒤 관련 회신을 보냈다. 요청 ID와 회신 ID가 서로 연결됐다. receipt 화면은 `자동 wake는 없습니다`, `queued는 상대 수신 성공이 아닙니다`, `수동 receive`, `180`을 보인다. 다른 회원의 receipt는 403이다. composer는 없다.
TTL은 DB 시각 기준 180초가 200이고, 181초와 이미 지난 시각은 422다. text 4097바이트는 422다. 4096바이트 왕복은 위 Node ID다.
같은 key와 같은 내용의 동시 송신 8회는 모두 200이고 메시지는 1개다. 같은 key의 다른 내용은 409다. 관련 회신은 1회가 200이고 두 번째는 409다. text 요청의 claim은 403이다.
기한을 지난 항목의 receipt는 `failed:expired`, `오프라인`, `새 key로 명시 송신`이다. 반복 lease 실패는 `failed:max_attempts`다. 원문 ACK 뒤 envelope는 남지 않는다. 이 기한은 DB에 넣은 시각이다. 벽시계 24시간을 기다리지 않았다.
trial 경로 `/v1/send`와 `/v1/test/send`는 일반 자격에서 403이다. 다른 회원의 수락은 403이다. 그 회원의 deny는 303이다. 옛 세대 수락은 403이다.
`schedule.query`의 claim policy는 deny를 포함한다. 잘못된 CSRF와 cross-site approve는 403이다. 정상 approve는 303이고 consume의 `executable`과 `disclosure`는 false다. authorize 뒤 `relay.result`는 403이다. gate 힌트는 `&lt;script&gt;`로 이스케이프된다. 원문이 없으면 `원문 부재`이고 approve는 409다.
`schedule.commit`의 claim policy는 `non-executable`이다. authorize의 `executable`은 false다.
queue는 99개인 상태에서 high 8개의 동시 송신이 200을 1개만 만든다. 나머지는 409다. 그 다음 high도 409다. 업무 ACK로 한 자리를 비우면 다음 송신은 200이다.
claim 5개의 동시 요청은 200이 4개다. 나머지 하나의 claim은 200이다. 이미 claim된 항목은 기한 이동 뒤 `failed:expired`가 되고 claim token이 비는다.
raw 32769바이트는 200이 아니다. raw 32768바이트의 `{` 본문도 메시지로 200이 아니다. receipt 20000개의 다음 text는 409이고 filler는 그 뒤에 남아 있다. pending gate 100개의 다음 approval request는 409다.
두 handler가 같은 DB를 쓴다. 점유 중인 신규 슬롯을 뺀 나머지만 입장한다. 신규 pull은 cleanup 채널을 쓰지 않고, cleanup ACK는 신규 채널을 쓰지 않는다. 상한의 다음은 429다. 해제한 뒤 남은 HTTP 행 수는 검사 전에 세어 둔 수와 같다.
프로세스 재시작 뒤 만료 전 DB 슬롯 16개는 pull을 429로 막는다. 시각을 1초 전으로 옮기면 pull은 200이다. 그 lease가 있으면 persist와 ACK도 200이다. 이 이동은 벽시계 대기가 아니다.
회원 신규 요청 40회는 422다. 41번째는 429다. key 철회 뒤 pull은 401이고 홈은 그 kid와 `철회`를 보인다.
신규 DB 슬롯 16개를 심은 뒤 새 text send는 429다. 그 전에 persist한 ACK는 200이고 receipt는 `delivered`다. 경로 `/home/gates/{id}/deny`는 303이고 gate는 `denied`다. unpair는 303이고 pair는 `revoked`다. logout은 303이고 Location은 `/?n=logout`이다.

## 재현한 제품 결함

H-1은 high다. 별도 Service의 프로세스 채널에서 익명 slow `POST /v1/text/send` 16개가 `io.ReadAll` 전에 신규 슬롯을 점유한다. 다음 요청은 HTTP 429이고 본문은 `capacity`다. 점유는 본문 기한 10초까지 남을 수 있다. 자격 검사와 DB rate는 본문을 다 읽은 뒤에 실행된다. 17차 관측은 send 429다.
같은 방식으로 slow `POST /v1/text/ack` 4개가 정리 슬롯을 점유한다. 다음 cleanup 요청은 429다. 17차 관측은 ack 429다. 이 두 줄은 [admission-defects.json](SAR-PUBLIC-MESSAGES-001-TESTER-test/admission-defects.json)에 있다.
M-1의 등급은 medium이다. `POST /home/gates/{id}`의 `decision=deny`는 본문 파싱 전에 신규 작업으로 분류된다. 신규 슬롯이 가득하면 그 형식은 429다. 정리 예산이 그 형식을 받지 않는다. 17차 관측은 formDeny 429다.
같은 조건의 `POST /home/gates/{id}/deny`는 잘못된 CSRF에서 403이다. 그 경로는 정리 작업으로 분류된 뒤 권한에서 거부된다. `/deny`가 429였다는 영향 문장은 추가하지 않았다.
추가로, receipt 20000 거절 409 뒤에 만료되지 않은 신규 `st.HTTP` 행이 1개 남았다. cleanup 행은 0개다. 11차부터 17차까지 같은 수다. 큰 snapshot의 해제 transaction이 2초 안에 끝나지 않으면 그 행은 30초 만료까지 남는다. 영향은 핸들러 15개만 진행 중일 때도 16번째 신규 요청이 429가 될 수 있다는 점이다. 등급은 medium이다. 동시성 검사는 남은 15개와 정리 4개만 통과시켰고, 해제 뒤에 추가 누수는 없었다.
이 세 결함의 기대값을 통과로 바꾸지 않았다. 결함이 다시 나타나지 않으면 프로브는 그 줄을 실패로 강제하지 않는다. 이번 17차에서는 세 결함이 모두 나타났다.

## DEV 원본 실패

DEV 표의 초기 실패는 그 실행 기록에 둔다. 이 QA가 그 실패를 통과로 바꾸지 않는다.
초기 `go test ./internal/relay`는 exit 1이다. 홈 전체의 `자동` 금지가 수동 안내의 `자동 wake`에 걸렸다.
첫 `make lint`는 exit 2다. `text.ts` Prettier 위반이다.
첫 `make test`는 exit 2다. MCP bundle이 CLI main을 실행했다.
첫 `make verify-mvp`는 exit 2다. 서명 helper의 KST 표기가 replay에서 422였다.
두 번째 `make verify-mvp`는 exit 2다. revoked 검사가 text envelope를 써서 schema 422가 먼저 났다.
DEV 최종 `make lint`, `make test`, `make verify-mvp`는 exit 0이다. 그 종료코드는 DEV 증거다. 이 QA는 같은 제품 명령을 후보 clone에서 다시 실행했고 종료코드는 0이다.

## 미실행

실메일 발송, 공개 서버, `link.knowslog.com` 배포, 벽시계 24시간, 운영 부하, 백업과 복원, 운영 자료 삭제, PS13, PS14, 외부 Grok Bot 계정, 외부 다닷 계정, 최종 노우↔다닷은 실행하지 않았다.
designer의 UX06, UX07 직접 시각 판정은 하지 않았다. 상태코드와 문구만 확인했다. OPS 독립 보안 리뷰의 H-1, M-1은 이 프로브가 같은 후보에서 다시 보았다. OPS 리뷰 문서 자체는 수정하지 않았다.
D12 운영 수락은 갱신하지 않았다. c6f0848 준비 리뷰는 미배정으로 둔다.

## 기록 검사

기록 커밋 `b8c5f2ea65b78871f8494985722b5dd1ff6d4e13`에서 `lint.py --from 09c523da8a3407288d9f5d711e1834af12bc7808`의 종료코드는 0이다.
merge_base는 기준 `09c523da8a3407288d9f5d711e1834af12bc7808`과 같다.
product-lint의 `make lint`는 passed다. product-test의 `make test`는 passed다.
ERROR는 0이다. WARNING은 1이다. 실행 불가는 0이다.
검사 파일은 3개다. 추가 줄은 17이다. 세 파일은 `.fullops-squad/PLANS.md`, `.fullops-squad/docs/design-docs/crud-design.md`, `.fullops-squad/docs/design-docs/module-design.md`다. 평가 기록, 인수인계, 역할 문맥은 lint 제외 경로다.
WARNING은 SIZE-001이다. `.fullops-squad/PLANS.md`는 947줄이다. 이전 blob은 940줄이다. 상한은 500줄이다.
늘어난 줄에는 이 QA의 결과 절이 들어 있다. PLANS는 coor 소유라 이 QA가 기존 기록을 줄이지 않았다.
SIZE-002는 없다. DEP 경고는 없다. 의존성 선언 파일은 변경하지 않았다.
설정 해시는 `9c49bb2dcd2b74d3b97b15b756d2ad1fdaa4e19d60d5e07cda5e215ef111c2d6`다.
JSON은 [fullops-lint-b8c5f2e.json](SAR-PUBLIC-MESSAGES-001-TESTER-test/fullops-lint-b8c5f2e.json)이다.
이 문단을 넣기 전 깨끗한 `b8c5f2ea65b78871f8494985722b5dd1ff6d4e13`에서 `deliverables.py --strict`는 검사 13, 미작성 0, 문제 0, 경고 0이다.
같은 트리의 `git diff --check` 종료코드는 0이다.
이 문단이 들어간 커밋은 위 lint의 HEAD가 아니다.
