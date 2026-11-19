---
title: SAR-PUBLIC-MESSAGES-001-TESTER — PS08–11 시나리오
status: draft
updated: 2026-10-06
owner: tester
tasks: [SAR-PUBLIC-MESSAGES-001-TESTER]
summary: "후보 09c523d의 text, receipt, gate, 한도, 정리 예산을 격리 프로세스로 확인한다"
---

# SAR-PUBLIC-MESSAGES-001-TESTER — PS08–11 시나리오

대상은 고정 후보 `09c523da8a3407288d9f5d711e1834af12bc7808`다. 실행 위치는 `/tmp/knowslink-messages-qa-09c523d`의 detached checkout이다. 제품 코드는 수정하지 않는다.
이 시나리오는 일반 회원 text, 관련 회신, receipt, gate, 한도, 동시성, 재시작, 정리 예산을 다룬다. 화면 조작은 Jev 웹 시나리오로 돌리지 않는다. 직접 시각 판정은 designer 범위다.
결과는 [QA 보고서](../qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER.md)와 `../qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER-test/`에 있다.

## 공통 조건

1. 시작과 끝의 clone HEAD는 위 후보다. `git status --porcelain`은 비어 있다.
2. 프로브는 고유 Compose project만 만들고 끝에서 그 project의 볼륨을 지운다. `knowslink-relay-1` `dfcd9d187117`, `knowslink-cloudflared-1` `07077b9ef5e4`, `knowslink-postgres-1` `bc3482dc52f2`는 바꾸지 않는다.
3. 회원은 이메일 코드로 로그인한다. 시험 allowlist와 공유 ServiceAuth를 일반 자격으로 쓰지 않는다. 메일은 메모리 inbox다.
4. 출력의 DB URL과 PEM은 가린다. 자격과 개인키는 증거에 남기지 않는다.
5. 실메일, 공개 서버, 벽시계 24시간, 운영 부하와 복원, PS13, PS14, 최종 노우↔다닷은 실행하지 않는다.
6. H-1, M-1, receipt 상한 뒤의 HTTP 슬롯 잔류는 관측될 때 `t.Error`로 남긴다. 기대값을 통과로 바꾸지 않는다.

## 시나리오

| ID | 조건 | 단계 | 확인 |
|---|---|---|---|
| S1 왕복 | 회원 둘, 별도 key, 활성 관계 | 로컬 Node와 MCP로 text를 보내고 상대가 pull, persist, ACK한다. 관련 회신을 같은 순서로 받는다. | 요청 ID와 수신 ID가 같다. 회신 ID가 연결된다. 4096바이트도 같은 경로다. 실제 Grok과 다닷 계정은 호출하지 않는다. |
| S2 receipt | 소유 화면 | 송신자가 자기 receipt를 연다. 다른 회원이 같은 주소를 연다. | 화면은 수동 receive, 180, queued가 수신 성공이 아니라는 문장을 보인다. 다른 회원은 403이다. composer와 본문은 없다. |
| S3 TTL과 길이 | DB 시각 | 180초, 181초, 지난 시각, 4097바이트를 보낸다. | 180초는 200이다. 나머지는 422다. |
| S4 멱등 | 같은 key | 같은 내용을 8개 동시에 보낸다. 이어서 다른 내용을 보낸다. | 8개는 200이고 메시지는 1개다. 다른 내용은 409다. |
| S5 관련 회신 | 전달된 원문 하나 | 회신 한 번 뒤에 두 번째 회신을 보낸다. text ID로 claim한다. | 첫 회신은 200이다. 두 번째는 409다. claim은 403이다. |
| S6 실패 | 기한과 반복 | 기한을 옮기고, lease를 반복한 뒤 ACK하지 않는다. | receipt는 `failed:expired` 또는 `failed:max_attempts`다. 안내는 오프라인과 새 key 송신을 말한다. |
| S7 인가 | 다른 회원과 trial | trial 경로, 다른 회원의 수락, 옛 세대를 시도한다. | 일반 자격의 trial과 다른 회원 수락, 옛 세대는 403이다. 수신 회원의 deny는 303이다. |
| S8 query | 검증 본문 | claim, 잘못된 CSRF, cross-site, 정상 approve, authorize, result를 순서대로 본다. | policy는 deny다. 잘못된 요청은 403이다. approve는 303이고 executable은 false다. result는 403이다. |
| S9 commit | 실행 불가 | commit을 전달하고 claim, authorize를 본다. | policy는 non-executable이다. authorize의 executable은 false다. |
| S10 gate 거부 | 힌트와 원문 없음 | 힌트에 script를 넣고, 원문을 없앤 뒤 approve를 누른다. | script는 이스케이프된다. 원문 부재가 보인다. approve는 409다. |
| S11 queue | 99와 high 8개 | 동시에 보내고, 한 건 ACK한 뒤 다시 보낸다. | 200은 1개다. 나머지는 409다. 자리를 비우면 다음은 200이다. |
| S12 claim | 5개 동시 | claim하고, 남은 하나와 기한 이동을 본다. | 200은 4개다. 남은 claim은 200이다. 이미 claim된 항목은 만료되고 token이 빈다. |
| S13 raw와 저장 | 32768, 20000, gate 100 | 경계를 넘긴다. | 초과는 200이 아니다. receipt와 gate의 다음은 409다. filler는 sweep로 바로 사라지지 않는다. |
| S14 공유 HTTP | 두 handler, 신규 16, 정리 4 | 남은 슬롯만 동시에 연다. | 채널이 바뀌지 않는다. 상한의 다음은 429다. 해제 뒤 행 수는 검사 전과 같다. |
| S15 재시작 | DB 슬롯과 만료 | 새 Service로 pull하고, 시각을 1초 전으로 옮긴다. | 만료 전에는 429다. 만료 후 pull은 200이다. lease가 있으면 persist와 ACK도 200이다. |
| S16 rate와 철회 | 40회와 key 철회 | `{` 본문 41회와 철회 뒤 pull, 홈을 본다. | 40회는 422다. 41번째는 429다. pull은 401이다. 홈은 `철회`를 보인다. |
| S17 정리 예산 | 신규 DB 슬롯 16 | persist 후 슬롯을 심고 ACK, 경로 deny, unpair, logout, GET receipt를 본다. | ACK, deny, unpair, logout은 303 또는 200이다. GET receipt는 429다. |
| S18 입장 결함 | 익명 slow body | 신규 16개와 정리 4개를 본문 완료 전에 붙잡는다. form deny와 경로 deny를 비교한다. | 신규와 정리의 다음은 `capacity` 429다. form deny는 429다. 경로 deny의 잘못된 CSRF는 403이다. |

## 시나리오 밖

`make lint`, `make test`, `make verify-mvp`는 제품 스위트다. 종료코드 0은 위 표의 운영 통과가 아니다.
designer UI와 OPS의 운영 측정은 각 역할이 한다. D12는 이 시나리오가 갱신하지 않는다.
