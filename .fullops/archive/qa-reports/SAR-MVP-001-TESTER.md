---
title: SAR-MVP-001-TESTER — 첫 안전 전달 기능의 독립 QA 보고서
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-MVP-001-TESTER]
summary: "고정 합성 후보 a6a10c7의 QA-01–11 판정, 종료코드, held 경계, UI 캡처 인계를 기록한다"
---

# SAR-MVP-001-TESTER — 첫 안전 전달 기능의 독립 QA 보고서

## 판정

- DEV 완료 SHA와 통합 후보는 모두 `a6a10c71977b7f3ec8274a1fb7c8a409f58e7c92`다. 준비 HEAD `707298e`와 이 검증 시작 HEAD `69870db66618202b8981eee55805f0020f421be8`는 지시서와 route만 추가한다. `git diff --stat a6a10c7 HEAD --` 제품 경로 출력은 비어 있다.
- QA-01–11의 실행 항목은 최종 `isolated.py`에서 통과했다. 미해결 critical/high 제품 결함은 없다. 제품 코드는 수정하지 않았다.
- DEC-02, DEC-03, Free N, 실adapter, A2A 현행 검토, WAL/backup 삭제, 고의 stale epoch, designer 시각 판정은 held다. 합성 통과로 바꾸지 않는다.
- 이 판정은 병합 승인이 아니다. D11–D13은 완료로 표시하지 않는다. 공개 운영과 실데이터는 실행하지 않았다.
- 시나리오: [SAR-MVP-001-TESTER.md](../scenarios/SAR-MVP-001-TESTER.md). 로그: [SAR-MVP-001-TESTER-test/](SAR-MVP-001-TESTER-test/).
- DEV `build/evidence` 로그는 읽지 않았고 판정 근거로 쓰지 않았다.

## 환경

- 실행 날짜 2026-10-03, Linux amd64. 브랜치 `fullops/tester`. 워크트리 `/home/shin/orca/workspaces/KnowsLink/fullops-tester`.
- Go 1.27.1, Node v22.22.2, npm 10.9.7, Python 3.12.3, GNU Make 4.3, Docker 29.4.3, Compose v5.1.3.
- 명령은 `run.py`가 셸 없이 실행하고 `[exit N]`과 프로세스 종료코드에 같은 값을 남긴다. `| tail`로 종료코드를 가리지 않았다.
- 합성 credential, private key, `POSTGRES_PASSWORD`는 로그·HTML·PNG에 남기지 않았다. 저장 HTML의 CSRF 값은 `redacted`다. Tunnel profile은 켜지 않았다. 각 Compose project는 종료 시 그 project만 `down --volumes` 했다.

## 수락 기준별 결과

최종 probe 한 줄은 `PROBE fail 0 held 8 pass 31`이다. 근거는 [probe.log](SAR-MVP-001-TESTER-test/probe.log)의 마지막 실행과 [probe-results.json](SAR-MVP-001-TESTER-test/probe-results.json)이다.

| ID | 결과 | 관찰 |
|---|---|---|
| QA-01 | 통과 | 가입 200. agent credential의 owner 경로 401 네 건. URL kid 422 `invalid_schema`. 잘못된 PoP 422, rotate 200, 이전 kid 재등록 409 `key_exists`, 이전 key 서명 401 |
| QA-02 | 통과, Free N held | 수락 전 403 `human_invite_required`, contacts 빈 배열. deny 후 403. active 재초대는 같은 세대. pending contacts 0. 동시 accept 8건 모두 200, contacts 1건 |
| QA-03 | 통과 | 인증 없는 재전송 401, digest 미노출. 중복 키·deliver:both·unknown·taskId·role:user·잘못된 Unicode·대문자 from은 모두 422 |
| QA-04 | 통과 | 같은 key+digest는 receipt만 200. 다른 digest 409 `idempotency_conflict`. 301초 422 `ttl_too_long` 뒤 같은 key 200. id 충돌 409 `id_collision` |
| QA-05 | 통과 | lease 간격 30.0초, attempts 1. ACK-before-persist 409. exp와 lease_until 차이 0.0초. attempts 2/3, 늦은 ACK 409, delivered, claim 승자 1, 재시작 후 claim 409, `failed:max_attempts` |
| QA-06 | 통과, epoch held | revoke 후 pull 200이고 token 없음. 교체 key 뒤에도 기존 대기는 임대되지 않음. 새 전송 200 후 임대. unpair 중 ACK 409, 재초대 세대 2, 이전 replay 403, 새 전송 200. 시계 +1시간 503 `abnormal_clock`, 복구 후 health 200. 없는 credential 401 |
| QA-07 | 통과 | query authorize executable false, disclosure false. ACK 전 transport는 delivered가 아님. commit 정책에 non-executable, done 403. H는 claim 없으면 403, 결속 200, 중복 409 `duplicate_gate`, 재귀 403, parent exp 초과 422 `expired` |
| QA-08 | 통과 | typed body와 정책 표시. hint 문구는 화면에 없음. GET은 pending 유지. agent 401, Basic 200, 잘못된 CSRF 403, approve 303, 재결정 409, consume 200이며 executable false, disclosure false |
| QA-09 | 통과 | 잘못된 방향 403. optional result 422. stack 422. denied 200. 두 번째 result 403 `sender_not_allowed`. completion은 `denied`. 결과 pull의 intent는 `relay.result`, persist/ACK 200, 다음 pull에 token 없음 |
| QA-10 | 통과, 실연결 held | evidence 수신 200. high가 먼저 보낸 메시지를 넘지 않음. webhook 404. registry 본문에 taskId, AgentCard, webhook, evidenceFetch 없음. 미설정 adapter는 `webhook` false, `evidenceFetch` false |
| QA-11 | 통과, WAL held | exp 제어 후 transport `failed:expired`, completion 빈 값, Envelope 없음, 행은 남음. accepted_at 2020은 다음 인증 요청의 sweep으로 행 삭제. 23시간은 Envelope 유지. 벽시계 24시간을 기다리지 않음 |
| V-01–04 | 캡처 완료, 시각 판정 held | pending, approved, denied, expired, revoked, unavailable, invalid_auth. Chrome 스크린샷 종료코드 0. 영상은 만들지 않음 |

## 실행한 명령과 종료코드

| 명령 | 종료코드 | 로그 |
|---|---|---|
| 버전·제품 diff | 0 | [versions.log](SAR-MVP-001-TESTER-test/versions.log) |
| `npm ci --prefix adapters` | 0 | [install.log](SAR-MVP-001-TESTER-test/install.log) |
| `make lint` 첫 실행 | 2 | [lint.log](SAR-MVP-001-TESTER-test/lint.log). `@types/node`가 없어 `node:crypto` 등을 찾지 못했다 |
| `make lint` (`npm ci` 이후) | 0 | 같은 lint.log의 두 번째 실행 |
| `make test` | 0 | [unit.log](SAR-MVP-001-TESTER-test/unit.log). integration tag는 포함하지 않는다 |
| `make build` | 0 | [build.log](SAR-MVP-001-TESTER-test/build.log) |
| `make verify-mvp` | 0 | [verify-mvp.log](SAR-MVP-001-TESTER-test/verify-mvp.log). 고유 project의 migration, Postgres race, TS stub, Go UI seed. 마지막 줄은 PASS |
| `isolated.py` 최종 실행 | 0 | probe.log 마지막 `[exit 0]`. 그 앞의 `[exit 1]`은 아래 관찰의 probe 수정 전 실행이다 |

`make verify`와 `make verify-runtime`은 실행하지 않았다. 이번 공식 통합 명령은 `make verify-mvp`다. 초기 골격의 404, 빈 migration, adapter unimplemented를 기대값으로 쓰지 않았다.

## 결함

제품 결함으로 인계할 항목은 없다. 아래는 관찰이다.

1. 첫 `make lint` 종료코드 2는 작업 트리의 `adapters/node_modules/@types/node` 부재다. `npm ci --prefix adapters` 종료코드 0 뒤 `make lint`는 0이다. 제품 소스를 고치지 않았다.
2. probe.log의 앞선 `[exit 1]`은 probe 기대값 오류다. 폐기 확정된 기존 대기를 교체 key가 복구한다고 기대한 점, `FROM relay_state` 없는 SQL, 자식 exp가 parent exp보다 1초 늦은 점이다. 제품을 고치지 않고 probe를 고친 뒤 최종 종료코드는 0이다.
3. 같은 parent의 두 번째 result는 409 `duplicate_result`가 아니라 403 `sender_not_allowed`다. 첫 result가 claim token을 지운 뒤의 fail-closed다. completion은 `denied`로 남고 두 번째 결과는 적용되지 않는다. QA-09 통과로 기록한다.
4. key revoke가 커밋된 뒤 기존 대기는 `failed:revoked`로 남고, 새 key를 등록해도 그 대기는 임대되지 않는다. 새 전송만 임대된다. 철회 확정 전과 후를 이렇게 구분한다.
5. revoked 캡처는 상태 `revoked`와 함께 `원문 부재: 승인 불가`를 보여 준다. unpair sweep이 parent envelope를 지우기 때문이다. Intent 칸은 비어 있다.
6. 동시 accept 8건과 claim 8건은 승자가 하나다. 트랜잭션 밖에서 stale epoch 503을 일부러 만들지 않았다. 이 항목은 held다.
7. V-04 권한 불명 캡처는 승인 화면이 아니다. 인증 없는 GET의 본문 `invalid_auth`다. 상태코드는 401이다.

## held와 미실행

- QA-02-free-n: Free N, 가격, slot 단위를 추정하지 않았다.
- QA-06-epoch-cas: 고의 stale epoch 거부를 외부에서 강제하지 않았다.
- QA-10-real-adapter, QA-10-a2a: loopback stub과 wire 필드 거부만 확인했다. 벤더 연결, 공개 성공, A2A 현행 개정 검토는 주장하지 않는다.
- QA-11-wal: 행 삭제는 WAL 또는 backup 삭제가 아니다.
- DEC-02: 실데이터 positive silent done은 held다. 합성 `status done` 403은 그 통과가 아니다.
- DEC-03: 공개 한도, 실제 신원 증명, singleton 처리량은 held다.
- UI-designer: 캡처 공유와 designer 시각 판정은 별개다. 이 보고서는 시각 판정을 통과로 기록하지 않는다.
- 다른 플랫폼, 빈 이미지 캐시의 첫 pull, `make verify-runtime`은 실행하지 않았다.

## UI 인계

designer는 같은 후보의 아래 파일을 직접 본다. tester는 상태 문구, 발신, 대상, 정책, typed body, 버튼 유무만 확인했다. 시간 변화는 SQL로 재현했고 영상은 없다.

| 화면 | HTML | PNG |
|---|---|---|
| V-01 pending | [v01-pending.html](SAR-MVP-001-TESTER-test/ui/v01-pending.html) | [v01-pending.png](SAR-MVP-001-TESTER-test/ui/v01-pending.png) |
| V-02 approved | [v02-approved.html](SAR-MVP-001-TESTER-test/ui/v02-approved.html) | [v02-approved.png](SAR-MVP-001-TESTER-test/ui/v02-approved.png) |
| V-02 denied | [v02-denied.html](SAR-MVP-001-TESTER-test/ui/v02-denied.html) | [v02-denied.png](SAR-MVP-001-TESTER-test/ui/v02-denied.png) |
| V-03 expired | [v03-expired.html](SAR-MVP-001-TESTER-test/ui/v03-expired.html) | [v03-expired.png](SAR-MVP-001-TESTER-test/ui/v03-expired.png) |
| V-03 revoked | [v03-revoked.html](SAR-MVP-001-TESTER-test/ui/v03-revoked.html) | [v03-revoked.png](SAR-MVP-001-TESTER-test/ui/v03-revoked.png) |
| V-04 unavailable | [v04-unavailable.html](SAR-MVP-001-TESTER-test/ui/v04-unavailable.html) | [v04-unavailable.png](SAR-MVP-001-TESTER-test/ui/v04-unavailable.png) |
| V-04 권한 불명 | [v04-unauthorized.html](SAR-MVP-001-TESTER-test/ui/v04-unauthorized.html) | [v04-unauthorized.png](SAR-MVP-001-TESTER-test/ui/v04-unauthorized.png) |

pending에는 Approve와 Deny가 있다. approved, denied, expired, revoked, unavailable에는 `승인·거절 버튼 비활성`이 있다. hint `이 힌트는 승인 근거가 아니다`는 pending 화면에 없다.

## 독립성과 증거 한계

- `make lint`, `make test`, `make build`, `make verify-mvp`, `isolated.py`는 이 워크트리에서 다시 실행했다.
- `make verify-mvp`의 합성 seed 경로는 재사용했다. DEV가 남긴 로그 파일은 재사용하지 않았다.
- `verify-mvp.log`는 `build/qa-fixture.json` 경로만 말한다. 그 파일은 gitignore 대상이며 커밋하지 않는다.
- 기준 ref `0dd08ec994771836c15d9d22a6a83393a71d7987`의 FullOps lint는 완료 커밋 `6bb6fad9f32a`의 깨끗한 트리에서 실행했다. 결과는 아래 절이다.

## tester 최종 검사

- 검사 대상 HEAD는 `6bb6fad9f32a`다. QA 기록 커밋은 `72941857c1fe6490fa7fe154b81b712b7212ec8d`다. 그 다음 커밋이 로그 끝의 빈 줄을 제거했다.
- `git diff --check 69870db66618202b8981eee55805f0020f421be8 HEAD` 종료코드는 0이다.
- `deliverables.py --strict` 종료코드는 0이다. 검사 13, 미작성 4, 문제 0, 경고 0이다. 미작성은 D04, D11, D12, D13이다. 이 QA는 그 상태를 완료로 바꾸지 않는다.
- `lint.py --from 0dd08ec994771836c15d9d22a6a83393a71d7987` 종료코드는 0이다. `product-lint`의 `make lint`는 passed다. ERROR 0, WARNING 3, 실행 불가 0이다.
- WARNING 3건은 모두 SIZE-001이다. `internal/relay/http.go` 487줄, `integration_test.go` 416줄, `store.go` 416줄이다. 기준 ref 이후의 제품 파일이며 이 QA가 만든 파일이 아니다. 제품 코드는 나누지 않았다.
