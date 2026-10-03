---
id: D08
title: 테이블정의서
status: review
updated: 2026-10-03
owner: dev
tasks: [SAR-MVP-001-DEV]
upstream: [D02]
summary: 실제 migration SQL에서 생성한 업무 테이블을 기록한다
---

# KnowsLink 업무 테이블 정의

이 문서는 `python3 scripts/schema.py`가 실제 migration SQL에서 생성한다.
실제 Postgres 적용·경합 검증은 `make verify-mvp`가 수행한다.

원천: `db/migrations/00001_relay.sql`

SHA256: `31aa978fb6e7d26eb4be7f65cd58c6abc286c0675bacd1856368221d095c7507`

```sql
CREATE TABLE relay_state (
    singleton boolean PRIMARY KEY DEFAULT true CHECK (singleton),
    epoch bigint NOT NULL DEFAULT 0,
    clock timestamptz NOT NULL DEFAULT clock_timestamp(),
    data jsonb NOT NULL DEFAULT '{}'::jsonb CHECK (jsonb_typeof(data) = 'object')
);
INSERT INTO relay_state(singleton) VALUES (true);
```

JSON 업무 엔티티와 CRUD 정본은 [데이터 모델](../design-docs/data-model.md)이다.
DB 행 삭제는 WAL·backup의 완전 삭제를 뜻하지 않는다.
