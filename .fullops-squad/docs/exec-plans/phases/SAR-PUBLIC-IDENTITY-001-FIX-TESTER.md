---
title: SAR-PUBLIC-IDENTITY-001-FIX-TESTER — 좁은 독립 QA 실행 기록
status: draft
updated: 2026-10-06
owner: tester
tasks: [SAR-PUBLIC-IDENTITY-001-FIX-TESTER]
summary: 후보 eb2e34b에서 실행한 좁은 QA 명령과 종료코드와 한계를 기록한다
---

# SAR-PUBLIC-IDENTITY-001-FIX-TESTER — 좁은 독립 QA 실행 기록

## 기준

- 후보 `eb2e34b93fe8d20fa1cd9166f73ff68d14bf17de`. 실행 checkout `/tmp/knowslink-identity-final-qa-eb2e34b`, detached. 기록 checkout은 `fullops/tester`.
- lint 기준 `9c915dc71e2a872243ffec294126d4668b4d32a4`. 공통 기준 `fullops-common-0.3.3`, `project.md`, coding-style, testing, security. 예외 없음.
- Task `task_42132cbbc82a`, Dispatch `ctx_80045eb5c01a`, Run `run_8ca8bc058ab7`.
- 원본 리뷰 `25b110f`, RATE 리뷰 `77dd464`, 원본 TESTER `9e2654d`, UI `cf0ab09`를 재사용한다. 원본 기록을 고치지 않았다.

## 실행 전후

시작 HEAD와 끝 HEAD는 `eb2e34b93fe8d20fa1cd9166f73ff68d14bf17de`다. `git status --porcelain`은 전후 모두 비어 있다.
기존 `knowslink-relay-1` `dfcd9d187117`, `knowslink-cloudflared-1` `07077b9ef5e4`, `knowslink-postgres-1` `bc3482dc52f2`의 ID는 유지됐다.
이 실행의 Compose project는 `knowslink-mvp-697ca71a30`이다. `down --volumes` 뒤 남은 컨테이너는 없다.

## 명령과 종료코드

| 명령 | 종료코드 |
|---|---|
| `npm ci --prefix adapters` 1회 | 0 |
| rate 단위 네 검사 `-race -count=1` | 0 |
| `make lint` | 0 |
| `make test` 1회 | 2 |
| `npm ci --prefix adapters` 2회 | 0 |
| `make test` 2회 | 0 |
| `make verify-mvp` | 0 |

1회차 `make test`의 2는 GNU make 종료코드다. 첫 `npm ci`가 zod `compat.js` 주석에 비ASCII 9바이트를 넣었고 esbuild가 그 파일을 거부했다. Go 패키지 줄은 그 전에 ok였다. registry tarball은 lock sha512와 일치했다. 2회차 설치 파일은 tarball member와 바이트가 같았다. product-test 판정은 2회차 0이다.
`make verify-mvp` 하위 명령 8개의 종료코드는 모두 0이다. 순서는 up --build --wait, stop relay, Go integration, up --wait, synthetic, seed, trial-check, down --volumes다.

## 판정과 한계

좁은 QA는 통과다. 새 critical/high는 없다. 상세 판정은 [QA 보고서](../../evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-FIX-TESTER.md)다.
`TestTrialHTTP`를 완화하지 않았다. 새 검사는 409 `invalid_lease`를 기대한다. 이 1회를 옛 relay-up 순서의 간헐 실패 해소로 기록하지 않는다.
실제 이메일, 운영 공개, 최종 노우↔다닷, 157 fixture 재실행, 새 시각 검수는 미실행이다. F2, F3, 다수 IP의 D02 공유 한도, `beta.sh` 합성 가입은 기존 후속이다.
FullOps lint HEAD, ERROR, WARNING, 실행 불가는 커밋 뒤 이 기록과 QA 보고서에 덧붙인다.
