"""Generate D08 directly from the committed SQL migration, preserving stamped metadata."""

import hashlib
from pathlib import Path

root = Path(__file__).resolve().parents[1]
source = root / "db/migrations/00001_relay.sql"
target = root / ".fullops-squad/docs/generated/db-schema.md"
raw = source.read_text()
sql = raw.split("-- +goose Up\n", 1)[1].split("-- +goose Down", 1)[0].strip()
metadata = ""
if target.exists() and target.read_text().startswith("---\n"):
    metadata = "---\n" + target.read_text().split("---\n", 2)[1] + "---\n"
target.parent.mkdir(parents=True, exist_ok=True)
target.write_text(metadata + "\n# KnowsLink 업무 테이블 정의\n\n"
                  "이 문서는 `python3 scripts/schema.py`가 실제 migration SQL에서 생성한다.\n"
                  "실제 Postgres 적용·경합 검증은 `make verify-mvp`가 수행한다.\n\n"
                  f"원천: `db/migrations/00001_relay.sql`\n\nSHA256: `{hashlib.sha256(source.read_bytes()).hexdigest()}`\n\n"
                  f"```sql\n{sql}\n```\n\n"
                  "JSON 업무 엔티티와 CRUD 정본은 [데이터 모델](../design-docs/data-model.md)이다.\n"
                  "DB 행 삭제는 WAL·backup의 완전 삭제를 뜻하지 않는다.\n")
print(target.relative_to(root))
