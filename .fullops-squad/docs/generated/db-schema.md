---
id: D08
title: 테이블정의서
status: review
updated: 2026-11-07
owner: dev
tasks: [SAR-MVP-001-DEV, SAR-HUMAN-PEER-001-DEV, SAR-DB-WAIT-001]
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

원천: `db/migrations/00002_human_writer.sql`

SHA256: `83cfa1aab7ddb803a0705434187af01386bbd8ccb60ab37a4cf5e9220c661f36`

```sql
-- +goose StatementBegin
CREATE FUNCTION guard_human_writer() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF COALESCE((OLD.data->>'HumanModel')::integer, 0) >= 1
       AND COALESCE((NEW.data->>'HumanModel')::integer, 0) < (OLD.data->>'HumanModel')::integer THEN
        RAISE EXCEPTION 'human authorization writer downgrade refused';
    END IF;
    RETURN NEW;
END;
$$;
-- +goose StatementEnd
CREATE TRIGGER human_writer_guard BEFORE UPDATE OF data ON relay_state
FOR EACH ROW EXECUTE FUNCTION guard_human_writer();
```

원천: `db/migrations/00003_relay_state_json.sql`

SHA256: `6e44e4b8e661f41be4ea0a51ae8095243b07b8e1ade49acd2f534148fbccc9e0`

```sql
DROP TRIGGER human_writer_guard ON relay_state;
ALTER TABLE relay_state DROP CONSTRAINT relay_state_data_check;
ALTER TABLE relay_state ADD COLUMN human_model integer NOT NULL DEFAULT 0;
UPDATE relay_state SET human_model = COALESCE((data->>'HumanModel')::integer, 0);
ALTER TABLE relay_state
    ALTER COLUMN data DROP DEFAULT,
    ALTER COLUMN data SET COMPRESSION lz4,
    ALTER COLUMN data TYPE json USING data::json,
    ALTER COLUMN data SET DEFAULT '{}'::json;
ALTER TABLE relay_state RENAME COLUMN data TO state;
ALTER TABLE relay_state ADD CONSTRAINT relay_state_state_check CHECK (json_typeof(state) = 'object');
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION guard_human_writer() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    RAISE EXCEPTION 'human authorization writer downgrade refused';
END;
$$;
-- +goose StatementEnd
CREATE TRIGGER human_writer_guard BEFORE UPDATE OF human_model ON relay_state
FOR EACH ROW WHEN (NEW.human_model < OLD.human_model) EXECUTE FUNCTION guard_human_writer();
```

JSON 업무 엔티티와 CRUD 정본은 [데이터 모델](../design-docs/data-model.md)이다.
DB 행 삭제는 WAL·backup의 완전 삭제를 뜻하지 않는다.
