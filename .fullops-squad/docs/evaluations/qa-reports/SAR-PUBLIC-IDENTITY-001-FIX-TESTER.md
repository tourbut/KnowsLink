---
title: SAR-PUBLIC-IDENTITY-001-FIX-TESTER — rate 격리와 trial 검사 격리 QA
status: draft
updated: 2026-10-06
owner: tester
tasks: [SAR-PUBLIC-IDENTITY-001-FIX-TESTER]
summary: 후보 eb2e34b의 rate 격리와 trial 검사 격리 QA 통과와 미실행 운영 확인을 기록한다
---

# SAR-PUBLIC-IDENTITY-001-FIX-TESTER — rate 격리와 trial 검사 격리 QA

## 판정

판정 후보는 `eb2e34b93fe8d20fa1cd9166f73ff68d14bf17de`다. 실행 위치는 `/tmp/knowslink-identity-final-qa-eb2e34b`의 detached checkout이다. 시작과 끝의 HEAD는 같고 추적 트리는 비어 있다.
좁은 변경 영향 QA는 통과다. 새 critical/high는 없다. 제품 코드, 원본 157 fixture, UI 증거, PLANS, board는 수정하지 않았다.
`make lint` 종료코드는 0이다. 검증된 `node_modules`에서 `make test` 종료코드는 0이다. `make verify-mvp` 종료코드는 0이다.
원본 `59b66ad`의 F1 medium 관측은 유지한다. 이 후보의 통과는 그 관측을 PASS로 다시 쓰지 않는다. 실제 이메일, 운영 공개, 최종 노우↔다닷은 미실행이다.

## 기준

날짜는 2026-10-06이다. 공통 기준은 후보의 `fullops-common-0.3.3`과 coding-style, testing, security, `project.md`다. lint 기준 ref는 `9c915dc71e2a872243ffec294126d4668b4d32a4`다.
원본 리뷰 `25b110f`와 RATE delta 리뷰 `77dd464`를 재사용한다. 원본 TESTER `9e2654d`와 UI `cf0ab09`는 `59b66ad` 결과로만 연결한다.
도구는 Go 1.27.1, Node 22.22.2, Python 3.12.3이다. Workers Free만 사용했다. 유료 전환과 Cloudflare 쓰기는 하지 않았다.

## 변경 범위

`59b66ad`와 이 후보의 제품 차이는 8개 파일이다. 비시험 Go 파일은 `internal/relay/identity.go`만 바뀌었다. 바뀐 부분은 `hit`과 `anonymousRate`, `memberRate`, `cleanupRate`다.
`internal/relay/member.go`와 `internal/relay/http.go`의 blob은 `59b66ad`와 같다. 회원 화면과 owner gate 템플릿은 같다. 새 시각 검수는 하지 않았다.
`9c915dc` 이후 제품 차이는 `scripts/verify_mvp.py`의 relay stop/up, `TestTrialForeignAllowlistRevokesLease`, README 한 문장이다. `TestTrialHTTP` 본문은 추가분만 있고 기존 단정은 유지된다.

## 실행한 명령

명령은 후보 루트에서 실행했다. 종료코드는 명령 자신의 값이다. 로그는 [SAR-PUBLIC-IDENTITY-001-FIX-TESTER-test/](SAR-PUBLIC-IDENTITY-001-FIX-TESTER-test/)에 있다.

| 명령 | 종료코드 | 결과 |
|---|---|---|
| `npm ci --prefix adapters` 1회 | 0 | 설치는 끝났고 `zod` `compat.js` 주석 9바이트가 손상됐다 |
| `go test -race -count=1 -v -run 'TestRatePrincipalIsolation\|TestRateStateBoundedUnderRotation\|TestRollingWindowBoundary\|TestSendLimits' ./internal/relay/` | 0 | 네 검사 PASS |
| `make lint` | 0 | `synthetic signup closed` 포함 |
| `make test` 1회 | 2 | GNU make 종료코드. esbuild가 손상된 zod 주석을 거부했다. Go 패키지는 그 전에 ok |
| `npm ci --prefix adapters` 2회 | 0 | registry tarball과 설치 파일이 바이트로 같다 |
| `make test` 2회 | 0 | product-test 판정. adapter MCP·trial boundary PASS |
| `make verify-mvp` | 0 | 아래 8개 하위 명령이 모두 0 |

손상된 첫 설치는 제품 단정 실패가 아니다. lock의 zod 4.6.5 sha512와 registry tarball이 일치했다. 두 번째 설치 파일은 그 tarball의 `package/v4/classic/compat.js`와 같았다. product-test 판정은 2회차 종료코드 0이다.

`make verify-mvp` 프로젝트는 `knowslink-mvp-697ca71a30`이다. 순서는 다음과 같다.

1. `compose up --build --wait relay` 종료코드 0.
2. `compose stop relay` 종료코드 0.
3. `go test -tags=integration -race -count=1 -v ./internal/relay` 종료코드 0.
4. `compose up --wait relay` 종료코드 0.
5. `node adapters/dist/synthetic.js` 종료코드 0.
6. `node adapters/dist/synthetic.js --seed` 종료코드 0.
7. `node adapters/dist/trial-check.js` 종료코드 0.
8. `compose down --volumes` 종료코드 0.

Go 통합 로그의 관련 PASS는 `TestEmailIdentity/refused_principal_does_not_spend_shared_rate`, `TestRatePrincipalIsolation`, `TestRateStateBoundedUnderRotation`, `TestRollingWindowBoundary`, `TestSendLimits`, `TestTrialHTTP`, `TestTrialForeignAllowlistRevokesLease`다. PostgresSafety, Gate, Trial 관련 회귀도 PASS다.

## 확인한 동작

`TestRatePrincipalIsolation`은 한 IP 1000회 중 30회, 한 회원 신규 40회, 한 회원 정리 20회까지 허용한다. 그 다음 요청은 거부된다. 다른 IP, 다른 회원 신규, 다른 회원 정리는 허용된다. 독립 principal의 합은 신규 200에서 허용되고 201번째에서 거부된다. 정리는 100에서 같다. 60초 뒤 세 경로는 다시 허용된다.
`TestRateStateBoundedUnderRotation`은 공유 예산이 찬 뒤 새 IPv6 source 10000개의 거부가 rate key 수를 늘리지 않는다. `http:new`는 201개, 기존 IP 버킷은 31개로 유지된다. 거부가 포화를 연장하지 않아 60초 뒤 다른 source가 허용된다.
`refused_principal_does_not_spend_shared_rate`는 한 IP의 verify 250회, 한 회원의 `/home` 250회, `logout-all` 150회 뒤에도 자기 요청이 429다. 같은 DB의 새 Service에서도 429가 유지된다. 다른 IP의 `/auth/start`는 303이다. 다른 회원의 `/home`은 200이다. 그 회원의 logout은 303이다.
`TestTrialForeignAllowlistRevokesLease`는 allowlist `agent_a,agent_b`로 send와 pull을 한 뒤, allowlist `trial_codex,trial_grok`인 두 번째 Service의 빈 transaction을 한 번 실행한다. 이어서 persist는 409이고 `error`는 `invalid_lease`다. `TestTrialHTTP`는 그 검사와 따로 PASS했다. 기존 persist 409 단정은 남아 있다.

## 원인 근거

`anonymousRate`, `memberRate`, `cleanupRate`는 principal bucket을 공유 bucket보다 앞에 둔다. `hit`은 첫 bucket이 가득 차면 그 bucket만 기록하고 공유 bucket을 바꾸지 않는다. 공유 bucket만 가득 차고 principal에 창 안 기록이 없으면 어느 bucket에도 기록하지 않는다.
`verify_mvp.py`는 Go integration 전에 같은 Compose project의 relay를 멈춘다. 그 relay의 allowlist는 `trial_codex,trial_grok`다. 검사 Service의 allowlist는 `agent_a,agent_b`다. 실행 중 relay의 sweep이 검사 lease를 회수하는 경로를 이 순서가 피한다. TypeScript 검사 전에 relay를 다시 시작한다.
이 1회 통과를 옛 순서의 간헐 실패가 사라졌다는 선언으로 쓰지 않는다. 옛 순서는 이 후보의 `make verify-mvp`가 실행하지 않는다. 회수 규칙 자체는 새 검사가 409 `invalid_lease`로 고정한다.

## 재사용한 증거

UI 템플릿 blob이 `59b66ad`와 같으므로 `cf0ab09`의 로컬 fixture 시각 검수를 그 SHA의 결과로 연결한다. 이 후보의 새 시각 검수로 기록하지 않는다.
원본 QA `9e2654d`의 QA-P01–P04와 발송 한도 검사는 같은 비-rate 코드에 의존한다. 그 157개 fixture를 이 과제에서 다시 실행하지 않았다. QA-P05의 "다른 source도 429" 관측은 F1이며 `59b66ad`에만 둔다.
디자인 lint와 공용 테마 전환은 미구성이다. 이번 UI 파일 변화가 없어 해당 없음이다.

## 기존 서비스

검사 전후 같은 컨테이너 ID가 유지됐다. `knowslink-relay-1`은 `dfcd9d187117`, `knowslink-cloudflared-1`은 `07077b9ef5e4`, `knowslink-postgres-1`은 `bc3482dc52f2`다.
이 실행의 Compose 명령은 `knowslink-mvp-697ca71a30`만 지정했다. `down --volumes` 뒤 그 프로젝트 컨테이너는 없다.
스냅샷 시작 때 `knowslink-mvp-a3723f3371`이 10초·30초 동안 떠 있었다. 이 실행의 명령은 그 프로젝트 이름을 쓰지 않았다. 끝 스냅샷에는 그 컨테이너가 없다.

## 미실행과 남은 항목

실제 운영 SMTP, 일반 사용자 이메일, 운영 공개, Cloudflare 쓰기, 최종 노우↔다닷은 미실행이다. 일반 서비스 수락은 이 QA로 완료하지 않는다.
F2 재발송 무효화와 누적 추측, F3 회원 gate 경로 rate, 독립 IP 다수가 공유 200을 채우는 D02 한도, `beta.sh`의 합성 가입 부재는 이 후보에서 새로 고치지 않았다. 새 결함으로 승격하지 않는다.
전체 157 fixture 재실행과 통계 40회 재현은 하지 않았다. `make generate`, `make schema`, `make verify-grok-plugin`, `make verify-runtime`은 이 변경의 SQL·adapter·런타임 설정 차이가 없어 실행하지 않았다.

## FullOps lint

내용 커밋 `40f8f3d69fd9abe0ddf32ac530a9227fc1a819c4`에서 lint.py 종료코드는 0이다. 기준은 `9c915dc71e2a872243ffec294126d4668b4d32a4`다. ERROR 0, WARNING 1, 실행 불가 0이다.
등록 명령 product-lint와 product-test의 종료코드는 모두 0이다. 검사 파일 8개, 추가 228줄이다.
WARNING은 SIZE-001이다. `.fullops-squad/PLANS.md`는 767줄이고 기준 blob은 744줄이며 상한은 500줄이다. 이 QA 커밋은 PLANS.md를 수정하지 않았다. DEP-001은 없다.
