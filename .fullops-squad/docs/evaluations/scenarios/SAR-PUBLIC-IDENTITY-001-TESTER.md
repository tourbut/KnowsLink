---
title: SAR-PUBLIC-IDENTITY-001-TESTER — 이메일 신원 fixture 시나리오
status: draft
updated: 2026-10-05
owner: tester
tasks: [SAR-PUBLIC-IDENTITY-001-TESTER]
summary: 고정 후보의 fixture 신원 시나리오와 미실행 공개·사람 확인을 구분한다
---

# SAR-PUBLIC-IDENTITY-001-TESTER — 이메일 신원 fixture 시나리오

대상은 고정 후보 `59b66ada8b36802484cc6d7e22523257b50572cc`다. 실행 위치는 `/tmp/knowslink-public-identity-qa-59b66ad`의 detached checkout이다. 제품 코드는 수정하지 않는다.
확인 코드를 시나리오 `values`에 넣지 않는다. Jev 웹 시나리오는 쓰지 않는다. 판정은 `probe.py`의 check id와 `f1_observe.py`의 상태 집계다.
결과는 [QA 보고서](../qa-reports/SAR-PUBLIC-IDENTITY-001-TESTER.md)와 `../qa-reports/SAR-PUBLIC-IDENTITY-001-TESTER-test/`에 있다.

## 공통 조건

1. 시작과 끝의 `git status --porcelain`은 비어 있다. HEAD는 위 후보다.
2. 임시 PostgreSQL, `scripts/mail_sink.py`, loopback relay를 이 실행에서만 띄운다. `127.0.0.1:5432`와 `127.0.0.1:8080`은 쓰지 않는다.
3. 서로 다른 source가 필요하면 `KNOWSLINK_CLIENT_IP_HEADER=CF-Connecting-IP`와 문서용 주소 `192.0.2.0/24`, `198.51.100.0/24`, `203.0.113.0/24`를 쓴다.
4. 출력과 `result.json`에 주소, 코드, cookie, token, DB URL을 남기지 않는다.
5. fixture 통과를 실제 사용자 이메일 통과나 운영 공개 통과로 기록하지 않는다.

## 시나리오

| ID | 덮는 기준 | 단계 | 통과 조건 |
|---|---|---|---|
| P01 가입과 거부 | QA-P01 | `probe.py`가 sink로 확인, 가입, 로그인, 오답, 만료, 재사용, 형식 오류, 발송 실패, 미설정 mail, 비-loopback 평문 SMTP, cross-site를 실행한다. | 해당 `p01-*`가 통과한다. 확인 전 회원은 없다. 발송 실패는 503이다. |
| P02 동일 회원 | QA-P02 | 재로그인, 다른 발급자, 동시 첫 가입, 점·plus 별칭을 실행한다. | 같은 이메일은 회원 하나다. 다른 발급자와 점·plus 별칭은 별도다. 대소문자는 접힌다. |
| P03 세션 | QA-P03 | 현재 로그아웃, 전체 로그아웃, 재시작, 유휴·절대 만료, 재확인, agent 수명을 실행한다. | 옛 세션은 거부된다. 전체 로그아웃은 최근 확인 창에서만 성공한다. agent는 로그아웃 뒤에도 남는다. |
| P04 회원 경계 | QA-P04 | A/B fixture의 조회, 결정, 관리 화면, synthetic, agent credential을 실행한다. | 다른 회원 조회는 403이다. 다른 회원 결정은 409다. 회원 cookie는 `/owner`와 agent API가 아니다. synthetic 기본값은 403이다. |
| P05 한도 | QA-P05 | 재발송 60초, 이메일·IP·전역 발송, 익명 30, 회원 40, 정리 budget, 회원 100, 재시작을 실행한다. | 경계 안은 허용이고 초과는 429 또는 503이다. 한도는 재시작 뒤 남는다. `make verify-mvp` 종료코드는 0이다. |
| F1 공유 예산 | QA-P05 추가 관측 | `f1_observe.py`가 잘못된 형식 200회를 한 source에 보낸 뒤 다른 source 1회를 보낸다. 이어서 fixture 로그인 뒤 같은 포화를 만들고 `/home`과 logout을 본다. | 처음 30회는 422다. 다음 170회는 429다. 그 시점 `http:new` 길이는 200이고 다른 source 버킷은 0인데 첫 요청은 429다. 기존 세션 `/home`은 429다. logout은 303이다. |
| P06 공개 경계 | QA-P06 | 운영 공개 후보와 현재 hostname이 이 Run에 없다. | 미실행. 로컬 synthetic 403을 공개 통과로 쓰지 않는다. |
| P07 사람 확인 | QA-P07 | HTML 문구 `p07-ux01-start`, `p07-ux02-verify`, `p07-ux03-home`만 자동으로 본다. | 문구 검사는 fixture 근거다. designer 화면 검수와 실제 사용자 이메일은 미실행이다. |
