---
title: SAR-PUBLIC-IDENTITY-001-FIX-TESTER — rate 격리와 trial 검사 격리 시나리오
status: draft
updated: 2026-10-06
owner: tester
tasks: [SAR-PUBLIC-IDENTITY-001-FIX-TESTER]
summary: 고정 후보의 rate 격리와 trial lease 회수 시나리오와 미실행 운영 확인을 구분한다
---

# SAR-PUBLIC-IDENTITY-001-FIX-TESTER — rate 격리와 trial 검사 격리 시나리오

대상은 고정 후보 `eb2e34b93fe8d20fa1cd9166f73ff68d14bf17de`다. 실행 위치는 `/tmp/knowslink-identity-final-qa-eb2e34b`의 detached checkout이다. 제품 코드는 수정하지 않는다.
Jev 웹 시나리오는 쓰지 않는다. 화면 문구는 바꾸지 않았고 시각 판정은 `cf0ab09`에 둔다. 판정은 Go 검사와 `make verify-mvp`의 종료코드다.
결과는 [QA 보고서](../qa-reports/SAR-PUBLIC-IDENTITY-001-FIX-TESTER.md)와 `../qa-reports/SAR-PUBLIC-IDENTITY-001-FIX-TESTER-test/`에 있다.

## 공통 조건

1. 시작과 끝의 `git status --porcelain`은 비어 있다. HEAD는 위 후보다.
2. `make verify-mvp`는 고유 Compose project만 만들고 끝에서 그 project의 볼륨을 지운다. 기존 `knowslink-relay-1`, `knowslink-cloudflared-1`, `knowslink-postgres-1`은 바꾸지 않는다.
3. 출력에 운영 인증값과 DB URL을 남기지 않는다. fixture 주소는 검사 코드의 문서용 주소다.
4. 원본 157 fixture와 UI PNG는 다시 실행하거나 고치지 않는다.
5. 실제 이메일, 운영 공개, 최종 노우↔다닷 통과로 기록하지 않는다.

## 시나리오

| ID | 덮는 기준 | 단계 | 통과 조건 |
|---|---|---|---|
| R1 자기 한도 | RATE-FIX F1 | `TestRatePrincipalIsolation`을 `-race -count=1`로 실행한다. | 한 IP는 30, 한 회원 신규는 40, 한 회원 정리는 20에서 멈춘다. 다른 principal은 허용된다. |
| R2 공유 합 | RATE-FIX F1 공유 한도 | 같은 검사가 독립 principal을 더한다. | 신규 합은 200에서 허용되고 다음은 거부된다. 정리 합은 100에서 같다. 60초 뒤 회복된다. |
| R3 회전 상한 | RATE-FIX cleanup 상태 | `TestRateStateBoundedUnderRotation`을 실행한다. | 새 source 10000개 거부가 key 수를 늘리지 않는다. 공유 포화가 연장되지 않는다. |
| R4 발송 한도 | 기존 `take` | `TestRollingWindowBoundary`와 `TestSendLimits`를 같이 실행한다. | 두 검사가 PASS한다. |
| R5 Postgres 격리 | 거부 principal과 다른 회원 | `make verify-mvp` 안의 `TestEmailIdentity/refused_principal_does_not_spend_shared_rate`를 본다. | 자기 verify, `/home`, logout-all은 429다. 재시작 뒤에도 429다. 다른 IP start는 303이다. 다른 회원 home은 200이고 logout은 303이다. |
| T1 회수 규칙 | Trial 원인 | `TestTrialForeignAllowlistRevokesLease`를 격리 DB에서 실행한다. | 외국 allowlist의 빈 transaction 뒤 persist는 409 `invalid_lease`다. |
| T2 기존 trial | 검사 완화 금지 | 같은 Go 실행의 `TestTrialHTTP`를 본다. | PASS한다. 기존 409 단정이 남아 있다. |
| T3 검사 순서 | verify-mvp 격리 | `make verify-mvp` 로그의 명령 순서를 본다. | stop relay, Go integration, up --wait relay, synthetic, seed, trial-check, down --volumes의 종료코드가 모두 0이다. |
| L1 제품 검사 | product-lint, product-test | 후보에서 `make lint`와 `make test`를 실행한다. | 종료코드 0. 손상된 설치의 첫 `make test`는 제품 판정이 아니다. |
| H1 사람 확인 | 운영 수락 | 실행하지 않는다. | 미실행. 로컬 통과를 실제 이메일·공개·노우↔다닷 통과로 쓰지 않는다. |
