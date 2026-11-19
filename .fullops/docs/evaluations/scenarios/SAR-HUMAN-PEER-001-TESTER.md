---
title: SAR-HUMAN-PEER-001 독립 QA 시나리오
status: draft
updated: 2026-10-12
owner: tester
tasks: [SAR-HUMAN-PEER-001-TESTER]
summary: HP-01–16을 고정 SHA 254ff6a의 합성 회원 A/B/C 단언과 미실행 이유로 연결한다
---

# SAR-HUMAN-PEER-001-TESTER 시나리오

대상 SHA는 `254ff6a98b12b7e51f0835491300f7695425d264`다. 환경 정본은 DEV `final-report.md`다. 원격 소스는 `/tmp/knowslink-humanpeer001-254ff6a98b12-c30163b6/qa-ccaa8cd7/repo`다. 릴레이 원격 loopback은 `41381`이다. 통과 실행의 증거 디렉터리는 `tester-independent-b0128515`다.

합성 회원은 A `검수회원가`, B `검수회원나`, C `검수회원다`다. A의 agent는 `검수에이전트가`와 `검수에이전트가둘`이다. B의 agent는 `검수에이전트나` 하나다. C의 agent는 `검수에이전트다` 하나다. B의 두 번째 agent는 이번 통과 실행에 없다.

비밀값, 친구 코드, 세션 쿠키는 시나리오와 로그 본문에 남기지 않는다. 요청 로그는 `/home/friends/request`와 `/home/friends/decision`까지 `{id}`로 가렸다. 판정은 check 이름과 HTTP 상태를 사용한다.

## HP-01 내 에이전트

41 rune 이름은 422다. 40 rune 이름은 303이다. A는 그 다음 `검수회원가`로 바뀐다. B는 A와 같은 표시 이름 303 뒤 `검수회원나`로 바뀐다. 41 rune agent 이름은 422다. A의 agent 두 개와 B·C의 agent 한 개가 생성 후 목록에 보인다.

재로그인으로 같은 agentID·키가 남는지 라이브로 확인하지 않았다. 결과는 not_run이다.

## HP-02 친구 발견

없는 코드는 303과 `n=requested`다. C의 친구 화면에 B 이름이 없다. 차단된 코드가 없는 코드와 같은 안내인지는 라이브로 제출하지 않았다. 결과는 not_run이다.

## HP-03 친구 요청

A가 B에게 요청한다. A의 대기 목록은 B 이름을 숨긴다. B의 수신 목록은 A 이름을 보인다. C와 요청자 A의 수락은 403이다. 다른 Origin의 수락은 403이다.

## HP-04 상태

없는 코드 요청을 취소한다. 반대 요청은 자동 수락되지 않고 대기만 남는다. A가 자기 요청을 취소한 뒤 C의 요청을 거절한다. 거절 뒤 C는 A의 활성 친구가 아니다. 같은 수락을 다시 보내면 303이다. 대기 24시간 벽시계와 자기 자신에게 보낸 요청은 not_run이다.

## HP-05 친구 목록

수락 뒤 A의 친구 상세는 B의 agent `검수에이전트나`를 보인다. C가 A-B 대화를 읽거나 그 대화로 보내면 403이다. B와 C의 친구 관계를 만들어 A에게 숨기는 경우는 not_run이다.

## HP-06 새 agent와 철회

친구 수락 뒤 A의 두 번째 agent 발신은 303이다. 새 수락을 요구하지 않았다. B의 두 번째 agent와 agent 철회 직후 목록 제외는 not_run이다.

## HP-07 내 agent끼리

A의 첫 agent에서 둘째 agent로 보내는 자기 연결은 친구 없이 303이다. 같은 agent에게 보내는 거부는 not_run이다.

## HP-08 대화

A가 고정 `from`/`to`로 연 대화에 owner 본문과 agent 답장이 함께 보인다. 브라우저 뒤로 가기와 새로고침 키 입력은 not_run이다. 서버 GET은 같은 쿼리의 두 본문을 다시 돌려준다.

## HP-09 발송 주체

owner 수신의 authorship는 `owner`다. agent 답장의 authorship는 `agent`다. 브라우저 세션으로 `GET /v1/contacts`를 호출하면 401 또는 403이다.

## HP-10 송수신

같은 attempt와 같은 본문의 재전송은 303이다. 같은 attempt의 다른 본문은 409다. 공백 본문은 422다. 4097 byte 본문은 422다. Enter만으로 보내는 키보드 경로는 not_run이다.

## HP-11 전달 상태

큐 상태 문구와 180초 벽시계 만료는 라이브에서 not_run이다. 코드의 큐 TTL 180초는 리뷰에서 확인했다. 라이브 미실행을 통과로 적지 않는다.

## HP-12 짧은 보존

공유 QA DB의 24시간 벽시계 이동은 not_run이다. designer가 같은 DB를 쓴다. `make test` EXIT 0과 `human_test.go`의 24시간 sweep은 별도 자동 검사다. 라이브 원문 소멸의 통과가 아니다.

## HP-13 해제·차단

B가 A를 차단하면 A는 기존 본문을 200으로 읽고 발신 폼이 없다. 차단 전 attempt의 재전송은 403이다. 차단 직전에 queued, leased, gate를 만들어 ACK와 claim을 보는 절차는 not_run이다. 차단과 별개인 해제만 하는 경로는 not_run이다.

## HP-14 차단과 재수락

B의 차단 해제는 303이다. 해제만으로 A의 활성 친구 링크가 돌아오지 않는다. 다시 요청하고 수락하면 generation이 이전보다 크다. 차단 전 attempt는 403으로 남는다. 새 attempt 발신은 303이다. A 화면에 차단이라는 단어가 없는지 별도 단언하지 않았다. 그 문구 비노출은 not_run이다.

## HP-15 승인과 안전

익명 `POST /v1/invites`는 401이다. 익명 홈은 303이다. 다른 Origin의 수락과 발신은 403이다. GET으로 수락하거나 보내지 않았다. CSP 헤더 측정과 rate flood는 not_run이다. 보안 검사의 실행한 항목을 예산으로 빼지 않았다.

## HP-16 실제 완료

같은 SHA에서 이름, 요청, 수락, 상대 agent 표시, 양방향 본문, 차단 후 신규 발신 거부가 통과했다. 이 결과는 독립 동작 QA의 실행 범위다. 제품 전체 수락은 designer 시각 판정과 이 시나리오의 not_run을 남긴 채 닫지 않는다. 정적 fixture와 DEV의 이전 PNG를 이 결과로 재명명하지 않는다.
