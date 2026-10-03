---
id: D10
title: 프로그램설계서
status: review
updated: 2026-10-03
owner: dev
tasks: [SAR-MVP-001-DEV]
upstream: [D02]
summary: 실제 프로그램 책임과 요구사항 및 검증을 연결한다
---

# KnowsLink 프로그램 책임과 검증

상위 요구사항은 [D02 SAR-MVP](../planning/product-specs/SAR-MVP.md)다.
API 정본은 [D05](interface-design.md), 저장 정본은 [D06/D07/D09](data-model.md)다.

| 모듈 | 책임·외부 경계 | 직접 검증 |
|---|---|---|
| protocol | Strict·Parse·SigningBytes·Digest, registry seed schema | TestFrozenParsingAndSigning |
| registry | compiled seed SHA256·manifest 공개 | protocol test·TS once |
| store | Postgres transaction·epoch·clock·current·sweep·ingest·lease | TestPostgresSafety |
| HTTP owner | credential·PoP·pairing·Go template·CSRF decision | PostgresSafety, GateFailureStates |
| HTTP agent | pull·shared persist·ACK·claim·authorize·result | PostgresSafety, ApprovalAndResultInstanceBinding |
| cleanup | 유휴 payload·metadata 회수 | 동일 sweep 경계·runtime 기동 |
| sqlc database | pgx/v5 LockRelay·SaveRelay 생성 | generate diff, 실제 Postgres |
| TypeScript adapter | registry·signature 재검증, persist·ACK·claim·gate·denied result | synthetic.ts 실제 HTTP 검사 |
| migrate·config·health | 별도 SQL-only Up·설정·readiness 실패 | 기존 unit/race·verify-runtime |
| Compose 검증 | 고유 project·private 제품 DB·loopback 시험 DB·Tunnel OFF | verify-mvp·verify-runtime |

`make test`는 외부 DB가 없어도 protocol과 기존 회귀를 수행한다. Postgres 검사는 integration build tag로 별도 실행하며 DB 환경이 없으면 실패한다.
`make verify-mvp`는 별도 DB에서 12회 ingest 경합·8회 claim 경합·3번째 lease·late ACK·pool 재시작을 검사한다.
권한 철회·세대 교체·TTL rollback·시계 이상·CSRF·중복 gate·M.id 재사용·잘못된 endpoint·optional 결과 거부를 검사한다.
TypeScript 검사는 Go 서버를 통해 policy 없음 deny와 gate approve 후 deny 및 최소 R 수신을 끝까지 수행한다.

독립 QA는 tester가 고정 후보 SHA에서 QA-01–11을 수행한다.
직접 시각 검수는 designer가 V-01–04를 수행한다. DEV의 자동 HTML 검사는 독립 시각 검수를 대체하지 않는다.
별도 fixed-SHA 코드 리뷰와 critical/high 차단은 coordinator가 담당한다.
