---
title: SAR-PUBLIC-AGENTS-001-TESTER — 일반 회원 연결·키·관계 시나리오
status: draft
updated: 2026-10-06
owner: tester
tasks: [SAR-PUBLIC-AGENTS-001-TESTER]
summary: 고정 후보의 연결·키·관계·한도 시나리오와 미실행 운영 확인을 구분한다
---

# SAR-PUBLIC-AGENTS-001-TESTER — 일반 회원 연결·키·관계 시나리오

대상은 고정 후보 `d1eef9bb90b9726149980320c42fb1fdbcaf584a`다. 실행 위치는 `/tmp/knowslink-agents-qa-d1eef9b`의 detached checkout이다. 제품 코드는 수정하지 않는다.
Jev 웹 시나리오는 쓰지 않는다. 화면 캡처와 시각 판정은 designer 범위다. 판정은 격리 Postgres 프로브, `make lint`, `make test`, `make verify-mvp`의 종료코드다.
결과는 [QA 보고서](../qa-reports/SAR-PUBLIC-AGENTS-001-TESTER.md)와 `../qa-reports/SAR-PUBLIC-AGENTS-001-TESTER-test/`에 있다.

## 공통 조건

1. 시작과 끝의 `git status --porcelain`은 비어 있다. HEAD는 위 후보다.
2. 프로브와 `make verify-mvp`는 각자 고유 Compose project만 만들고 끝에서 그 project의 볼륨을 지운다. 기존 `knowslink-relay-1`, `knowslink-cloudflared-1`, `knowslink-postgres-1`은 바꾸지 않는다.
3. 출력에 운영 인증값과 DB URL을 남기지 않는다. fixture 주소는 `example.test`와 문서용 IP다. 메일은 메모리 inbox다.
4. 실제 이메일, 운영 공개, 플랫폼, 최종 노우↔다닷 통과로 기록하지 않는다.
5. enqueue는 `POST /v1/send`다. exec는 claim과 authorize, result는 `relay.result` send다. `/v1/enqueue`는 없다.

## 시나리오

| ID | 덮는 기준 | 단계 | 통과 조건 |
|---|---|---|---|
| C1 교차 owner | 연결 발급과 조회 | Bob이 Alice agent로 connect하고 Alice 연결을 연다. | connect와 조회가 거부된다. Alice 발급 화면에는 Node prepare 안내가 있고 개인키는 없다. |
| C2 PoP 치환 | owner, agent, client, public | prepare 서명과 client를 바꾼다. | owner, agent, public 치환은 422, client 치환은 403이다. 키는 0개다. |
| C3 저장 grant | 다른 owner의 agent와 owner 필드 | 저장 연결의 agent를 다른 owner로 바꾸고, 이어서 owner만 Bob으로 바꾼다. | info는 401이다. Bob confirm은 303일 수 있으나 complete는 401이고 새 키는 없다. |
| C4 1회와 동시 완료 | complete | 승인 뒤 complete를 8회 동시에 호출한다. | 200은 1회다. 재호출은 401이다. |
| C5 재인증과 정리 | 5분, cancel, 만료 | 검증 시각을 5분 과거로 보낸 뒤 confirm, connect, 철회, 생성을 호출하고 cancel과 만료를 본다. | 재인증 대상은 422 `5분`이다. cancel은 303, 만료 info는 401이다. 기존 키는 유지된다. |
| K1 키 3개 | 상한과 선택 철회 | key1–key3를 연결하고 key2만 철회한다. | 네 번째 register는 409다. key2 pull은 401, key1과 key3 pull은 200이다. |
| K2 실패 보존 | 잘못된 proof와 cancel | prepare proof를 망가뜨리고 cancel한 뒤 기존 키로 pull한다. | info는 401이고 기존 pull은 200이다. 활성 키 수는 그대로다. |
| K3 회전과 옛 자격 | enqueue, pull, ACK, exec, result | rotate 뒤 새 Service에서 옛 credential로 pull, persist, ack, claim, authorize, gate-consume, send를 호출한다. | 모두 401이다. 새 credential pull은 200이다. 활성 키는 1개다. |
| R1 타회원 초대 | 명시적 수락 | Alice가 Bob에게 초대한다. | Alice 홈에는 수락 버튼이 없다. Bob 홈에는 수락과 거절이 있다. Alice 수락은 403이다. pending send는 403이다. |
| R2 거절과 세대 | 거절, 재초대, 양측 철회 | deny, 옛 세대 수락, 새 세대 수락, 양쪽 unpair를 순서대로 한다. | 옛 세대는 403, 새 세대 수락 뒤 send는 200, unpair 뒤 send는 403이다. |
| R3 같은 owner | 첫 수락과 동시 수락 | 같은 owner의 두 agent를 초대하고 수락 8회를 동시에 보낸다. | 초대 직후는 pending이다. 수락은 모두 303이고 쌍은 1개, 세대는 1이다. 세대 0은 403이다. |
| P1 agent 한도 | 전역 200, owner 5, 재시작 | 외국 200, 철회 한 칸, owner 5, 외국 owner 활성 전환을 본다. | 경계의 다음 생성은 409다. 새 Service도 전역 한도를 유지한다. 한 칸을 비우면 한 번 303이다. |
| P2 관계 한도 | pending 200과 송신 10, active 400과 owner 20, 24시간 | 더미 쌍으로 경계를 채운 뒤 초대와 수락을 한다. | 경계는 409다. 만료 뒤 새 초대는 다음 세대다. 옛 세대 수락은 403이고 상태는 pending으로 남는다. |
| P3 안전 정리 | 공유 예산이 찬 뒤 deny | `http:new`를 201로 채우고 생성과 deny를 호출한다. | 생성은 429다. deny는 303이다. 회원 버킷은 늘지 않는다. |
| M1 깨진 요청 | 회원 40, 익명 30, connect prepare | 로그인 form `%`, 무세션 생성, prepare `{`를 한도 너머까지 보낸다. | 한도 안은 422 또는 401이고 다음은 429다. 깨진 form은 agent를 만들지 않는다. 다른 principal은 자기 한도를 쓴다. |
| M2 포화 불연장 | 공유 200 | `http:new`가 찬 새 IP로 깨진 prepare를 보낸다. | 429다. 그 IP 버킷은 비어 있고 `http:new` 길이는 201이다. |
| M3 회전 예산 | credential 교체 | 같은 agent로 pull 40회 뒤 rotate한다. | 버킷 길이는 40으로 남는다. 새 credential pull은 429이고 재시작 뒤에도 429다. 옛 credential은 401이다. |
| N1 CLI | prepare, confirm, complete | 로컬 `node adapters/dist/connect.js`로 잘못된 URL, 기존 폴더, prepare, 조기 complete, confirm, complete, 재완료를 한다. | 모드 0700/0600, 안내 문장, 키 출력 금지, 파일 보존, 1회 완료, pull 200이다. |
| V1 회귀 | lint, test, verify-mvp | 후보에서 세 명령을 한 번씩 실행한다. | 종료코드 0. verify-mvp는 stop relay, Go integration, up relay 순서를 유지한다. 신원과 gate, trial lease 검사가 PASS다. |
| H1 사람 확인 | 운영 수락 | 실행하지 않는다. | 미실행. 로컬 통과를 실제 이메일·공개·노우↔다닷 통과로 쓰지 않는다. |
