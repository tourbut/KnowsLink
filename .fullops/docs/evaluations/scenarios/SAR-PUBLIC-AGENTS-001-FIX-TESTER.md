---
title: SAR-PUBLIC-AGENTS-001-FIX-TESTER — 보존·rate·관계 조건 시나리오
status: draft
updated: 2026-10-06
owner: tester
tasks: [SAR-PUBLIC-AGENTS-001-FIX-TESTER]
summary: "후보 458798c의 보존, rate, POLICY 관찰 조건 시나리오를 원본 QA와 구분한다"
---

# SAR-PUBLIC-AGENTS-001-FIX-TESTER — 보존·rate·관계 조건 시나리오

대상은 고정 후보 `458798c2ee15c179edacfd6f94ebb9896d26f411`다. 실행 위치는 `/tmp/knowslink-agents-fix-qa-458798c`의 detached checkout이다. 제품 코드는 수정하지 않는다.
이 시나리오는 보존, rate, POLICY 관찰 조건만 다룬다. 원본 `d1eef9bb90b9726149980320c42fb1fdbcaf584a` QA의 전체 시나리오를 다시 실행하지 않는다.
화면 캡처와 시각 판정은 designer 범위다. 시간 이동은 DB에 저장한 시각이다. 벽시계 24시간을 기다리지 않는다.
결과는 [QA 보고서](../qa-reports/SAR-PUBLIC-AGENTS-001-FIX-TESTER.md)와 `../qa-reports/SAR-PUBLIC-AGENTS-001-FIX-TESTER-test/`에 있다.

## 공통 조건

1. 시작과 끝의 clone HEAD는 위 후보다. `git status --porcelain`은 비어 있다.
2. 프로브는 고유 Compose project만 만들고 끝에서 그 project의 볼륨을 지운다. `knowslink-relay-1` `dfcd9d187117`, `knowslink-cloudflared-1` `07077b9ef5e4`, `knowslink-postgres-1` `bc3482dc52f2`는 바꾸지 않는다.
3. fixture 주소는 `example.test`와 문서용 IP다. 메일은 메모리 inbox다. 출력의 DB URL과 긴 토큰은 가린다.
4. 실제 이메일, 운영 공개, 운영 자원, 노우↔다닷은 실행하지 않는다.
5. 원본 QA `bcb06b89bcb36d69a99cbeef3d94e4a9ffe88361`와 UI `d1651784c4338efeb0d6141467d563c6b354e4a5`의 실패와 held는 그 SHA의 결과로 둔다.

## 시나리오

| ID | 조건 | 단계 | 통과 조건 |
|---|---|---|---|
| P1 pending 반복 | 같은 방향과 반대 방향 | 초대를 세 번 보낸다. 한 번은 수신 쪽에서 보낸다. | pending은 1개다. 세대와 기한은 그대로다. 회원 rate 기록이 늘어난다. 발신 홈에는 수락 버튼이 없다. pending send는 403이다. |
| P2 만료 | 기한 전과 기한 도달 | Exp를 2분 뒤로 둔 뒤 수락한다. 다른 실행에서는 Exp를 1분 전으로 둔다. | 기한 전 수락은 303이고 active다. 기한 도달 수락은 403이고 pending 수는 0이다. 수동 재초대는 다음 세대다. |
| P3 양측 철회 | 발신 철회와 수신 철회 | 각 방향으로 unpair한 뒤 수동 재초대와 새 수락을 한다. | 철회 뒤 send는 403이다. 옛 세대 수락은 403이다. 새 수락 뒤에만 send가 200이다. 같은 영수증의 재전송은 403이다. |
| P4 거절과 경합 | 늦은 결정과 동시 결정 | deny 뒤 옛 수락, 수락 대 거절 동시 요청, 수락 8회를 본다. | 늦은 결정은 403이다. 경합은 303 하나와 403 하나다. 수락 8회는 active 하나와 같은 세대다. |
| P5 같은 owner | 자동 수락 없음 | 같은 owner의 두 agent를 초대하고 그 owner가 수락한다. | 초대 직후는 pending이다. 옛 세대 수락은 403이다. 다른 owner의 수락은 403이다. |
| P6 한도와 정리 | pending 10·200, active 20·400, rate | 더미 쌍과 공유 예산을 채운 뒤 초대, 수락, 거절, unpair를 한다. | 상한의 다음은 409 또는 429다. 세대는 늘지 않는다. 신규 예산이 꽉 찬 deny는 303이다. 정리 예산이 꽉 찬 unpair는 429이고 상태는 pending이다. |
| P7 재시작 | 프로세스와 재로그인 | 새 Service로 홈을 연다. 같은 이메일로 다시 로그인한다. | 세대가 유지된다. 재로그인이 관계를 복구하지 않는다. 철회 agent와 비활성 owner의 재초대는 403이다. |
| M1 키 기록 | 상한 20과 kid | key1을 철회한 뒤 같은 kid로 prepare한다. 키를 20개까지 회전한다. | prepare는 409 `key_exists`다. 21번째 연결은 409이고 열린 연결은 늘지 않는다. 선택 철회는 303이다. 30일 시각 이동은 키를 지우지 않는다. |
| M2 agent 보존 | 24시간 경계 | 철회 시각이 비어 있는 기록을 한 번 읽고, 25시간 전으로 옮긴 뒤 다시 읽는다. | 첫 읽기는 기록을 남기고 시각을 채운다. 25시간 전에는 홈에서 사라진다. 안내는 권한 복구와 백업 영구 삭제를 말하지 않는다. |
| M3 기록 상한 | owner 기록 10 | 기록을 10개로 채운 뒤 생성과 철회를 한다. 철회 시각을 25시간 전으로 옮기고 다시 만든다. | 생성은 409다. 결제 문구는 없다. 철회는 303이고 기록 수는 그대로다. 정리 뒤 수동 생성은 303이다. 새 agent는 옛 pair를 받지 않는다. |
| R1 거부 경로 | rollback과 집계 | 오류를 반환하는 transaction과 `%` 본문 40회를 구분한다. | 오류 경로는 phantom을 버리고 rate와 정리 결과를 남긴다. `%`는 422이고 agent는 늘지 않는다. 41번째는 429다. 다른 회원 생성은 303이다. |
| R2 무효 세션 | 익명 30 | 잘못된 세션으로 `/home`과 연결 조회를 섞어 30회 호출한다. | 30회는 303 `/?n=expired`다. 31번째와 재시작한 GET, logout, reauth는 429다. 다른 IP의 첫 GET은 303이다. |
| R3 connect 429 | 실제 시각 | `{` 본문을 31회 보낸다. | 30회는 422다. 31번째는 429다. `retry_at`과 `Retry-After`는 같은 시각이다. 다른 IP는 422다. |
