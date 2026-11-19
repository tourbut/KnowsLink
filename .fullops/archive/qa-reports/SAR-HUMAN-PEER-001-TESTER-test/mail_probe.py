"""Count recent sink files without printing addresses, codes, or bodies."""
import json
import pathlib
import sys
import time

base = pathlib.Path(sys.argv[1])
qa = max(base.glob("qa-*/runtime.json"), key=lambda path: path.stat().st_mtime_ns).parent
now = time.time()
rows = []
for message_file in sorted((qa / "mail").glob("*.eml"), key=lambda path: path.stat().st_mtime, reverse=True)[:5]:
    raw = message_file.read_bytes()
    rows.append({
        "age_s": round(now - message_file.stat().st_mtime, 1),
        "size": len(raw),
        "has_domain": b"@humanpeer.invalid" in raw,
        "has_to": b"To:" in raw,
        "has_qp": b"quoted-printable" in raw,
        "has_subject_marker": b"Subject:" in raw,
    })
print(json.dumps({"qa": qa.name, "recent": rows}))
