---
title: SAR-MVP-001-DEV 실행 기록
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-MVP-001-DEV]
summary: 로컬 합성 MVP의 기술 계획과 구현 및 검증 증거를 기록한다
---

# SAR-MVP-001-DEV 실행 기록

## 기준과 기술 계획

착수 SHA는 f72c5b4다. lint 기준은 0dd08ec994771836c15d9d22a6a83393a71d7987이다.
원천은 service-design 7bc9ea1이며 fullops-common-0.3.2와 D02 SAR-MVP를 적용한다.
현재 Dispatch는 ctx_66989e4a879d이며 Task는 task_491be61b82eb다.

기존 Go relay와 별도 goose migration을 유지한다. Postgres 단일 행의 JSON 업무 상태와 FOR UPDATE 및 epoch CAS로 모든 권한 확정·멱등·gate·claim을 직렬화한다.
이는 로컬 합성 MVP의 명시적 처리량 한계다. 공개 운영과 다중 사용자 규모 확장은 정규화·제품 자원 제한 확정 뒤 별도 수락한다.
owner는 서버가 발급한 별도 opaque credential로 인증한다. agent credential은 owner 작업에 사용할 수 없다.
키 등록은 owner 인증과 AgentID/kid/pubkey에 묶인 PoP를 사용한다. RFC8785는 검증된 라이브러리를 재사용한다.
한 TypeScript stub이 공유 DB inbox를 저장하고 ACK·claim을 확인한 뒤 deny-by-default 또는 gate를 수행한다.
외부 도구·정보 공개·벤더 연결·Tunnel 변경은 수행하지 않는다.

## API 근거

Context7 /cyberphone/json-canonicalization에서 Go Transform의 RFC8785 숫자·중복 키 처리를 확인했다.
Context7 /websites/sqlc_dev_en에서 pgx/v5 WithTx, FOR UPDATE, :execrows 계약을 확인했다.
실제 의존성은 go.mod와 go.sum에 고정한다. sqlc CLI는 v1.30.0으로 고정한다.

## 검증과 후속

구현 후보를 완료했다. 독립 QA와 직접 UI 검수는 고정 SHA로 후속 담당에게 인계한다.

| 검사 | 실제 결과·근거 |
|---|---|
| make lint | 종료코드 0; build/evidence/lint.log, 완료 커밋의 FullOps product-lint로 다시 고정한다 |
| make test | 종료코드 0; build/evidence/unit-final.log, DB 환경 없음의 skip을 실제 DB 검사와 구분한다 |
| make build | verify-runtime와 verify-mvp의 선행 target, 종료코드 0 |
| make generate | 종료코드 0; sqlc v1.30.0 실제 schema와 쿼리 생성, build/evidence/generate.log |
| make schema | 종료코드 0; 실제 migration SQL의 D08 생성 |
| make verify | 종료코드 0; format/lint/type 위반 실패와 원복 통과, build/evidence/verify.log |
| make verify-runtime | 종료코드 0; actual migration·2개 테이블·readiness·설정/DB 실패·네트워크 차단 Tunnel, build/evidence/runtime.log |
| make verify-mvp | 종료코드 0의 후보 검사; Postgres race/claim/ACK/gate·TS HTTP·UI seed, build/evidence/mvp-candidate.log; 최종 추가 field/pair 검사는 mvp-fixed-code.log에 고정한다 |
| deliverables.py --strict | 검사 13, 미작성 4, 문제 0, 경고 0; D04/D11–13은 후속 담당 |
| git diff --check | 종료코드 0 |

제품 코드 검사는 작업 중 snapshot에서 수행했다. 완료 SHA는 최종 회신으로 고정한다.
최종 문서·handovers 기록 커밋은 제품 코드가 바뀌지 않으면 위 동작 증거를 그대로 재사용한다.
로그는 이 체크아웃의 Git 미추적 build/evidence 아래에 있다. 영구 기록은 이 문서의 결과표와 완료 아카이브다.

## 구현과 요구사항별 증거

- MVP-01/02: owner/agent credential 분리, owner-bound PoP, atomic rotate/revoke, kid 재할당 금지, pending/active 분리와 8회 accept 경합.
- MVP-03: strict JSON·duplicate/Unicode/unknown/container 거부, Ed25519 signed fields·semantic digest·JCS 숫자, hint/key/evidence/body 한도.
- MVP-04/05: auth-first, 12회 atomic replay, digest conflict, TTL/id rollback, receipt-only 만료 replay, 300s/30s/3/24h, 마지막 lease와 late ACK.
- MVP-06/07: shared Inbox commit 후 ACK, 8회 claim 경합에서 1회 성공, 새 pool 재시작 중복 claim 차단, key/pair/owner 재검사, epoch CAS와 비정상 DB 시계 차단.
- MVP-08–11: no-policy deny·commit non-executable, H 부모·claim·digest·generation, duplicate gate/CSRF/owner 권한/원자적 결정·consume, 다른 M claim 재사용·wrong R endpoint·optional 결과·done 거부, 원요청 key rotation 후 결과 공개 차단.
- MVP-12–14: central hash-locked seed와 구현 subset, webhook/evidence fetch OFF, loopback TypeScript stub 1개. 실제 벤더/A2A wire 호환은 주장하지 않는다.
- MVP-15/16: exp/revoke/result 후 raw·inbox 제거, background sweep, receipt 분리, escaped Go UI의 pending/approved/denied/expired/revoked/unavailable·verified body·policy·disabled actions.

## 실패와 해결 근거

처음 시험 DB는 internal network의 publish 포트에서 연결되지 않았다. test-only override에 자기 ingress network를 추가해 해결했다. 제품 private Postgres 설정은 유지했다.
원문을 null로 저장한 RawMessage가 다음 decode에서 길이 4의 null bytes가 되었다. payload 필드 omitempty로 원문 부재 상태를 보존했다.
duplicate_gate의 기대 HTTP 코드를 구현의 conflict 409로 교정했다. duplicate gate의 거부 자체를 완화하지 않았다.
결과 공개에서 부모 key의 current-auth를 함께 검사하도록 보강하고 rotation 후 공개 거부를 실제 DB로 검증했다.

## 기술 판단과 범위 차이

업무 상태는 정규화 테이블 대신 singleton JSON+row lock을 사용한다. SQL·sqlc·Postgres transaction은 실제 업무 저장과 CAS를 수행한다.
이는 합성 MVP의 처리량 한계이며 공개 성능을 보장하지 않는다. 코드 ponytail 주석·D03/D07에 후속 정규화 조건을 기록했다.
D06/D07/D09는 front matter의 단일 id 검사 때문에 각각 data-model/database-design/crud-design으로 매핑했다.
새 검사 명령은 제품 Makefile과 project.md에 추가했다. lint 규칙이나 원천 문서는 수정하지 않았다.
PLANS와 board는 지시서의 coor 소유권에 따라 직접 수정하지 않는다. 완료 회신으로 coor가 함께 갱신한다.

## QA/UI 실행 인계와 보류

독립 QA는 README의 make install/build/verify-mvp를 고정 후보에서 실행한다. QA-01–11을 tester에게 인계한다.
직접 UI는 README의 고유 knowslink-qa-local Compose와 synthetic.js --seed를 실행한다. owner 인증은 build/qa-fixture.json의 B-owner credential이다.
V-01–04는 designer가 같은 후보에서 캡처한다. 원문 부재·권한 불명은 TestGateFailureStates의 DB fixture로 재현하며 제품 test hook은 없다.
생성 원요청 TTL은 180초다. 만료 화면 외 검수는 seed 후 180초 안에 수행하고 필요 시 새 seed를 만든다.
DEV 자동 HTML 검사·합성 seed는 독립 시각 수락을 대신하지 않는다.

DEC-01 quota/가격·DEC-02 disclosure/output schema·DEC-03 운영 한도는 designer 소유이며 값을 만들지 않았다.
coor의 DEC-03 제안은 운영 측정·poll/ACK·철회 처리·claim 회수 조건으로 기술 검토했으며 수치는 적용하지 않았다.
DEC-04 실제 벤더 연결은 후속 dev다. DEC-05 Docker/Tunnel은 수락 후보 이후 OPS다.
hostname 사용자 결정 link.knowslog.com을 수신했지만 기존 Tunnel과 다른 컨테이너는 변경하지 않았다.
고정 SHA 독립 리뷰·tester QA·designer UI·critical/high 차단과 최종 수락은 coor가 후속 배정한다.
