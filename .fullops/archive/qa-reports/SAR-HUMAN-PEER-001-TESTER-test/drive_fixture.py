"""Drive the recorded final-SHA helper auth/fixture without using the moved local DEV HEAD.

The helper's parent-HEAD assert is skipped because the dev worktree moved after the
recorded environment was created. Remote HEAD, a clean remote tree, known_hosts, and
RejectPolicy are still required. Secret values are not printed.
"""
import importlib.util
import json
import pathlib
import shlex
import sys

HELPER = pathlib.Path(r"C:\Users\shin\AppData\Local\KnowsLinkDevTools\humanpeer001-runtime-ctx9ded")
SHA = "254ff6a98b12b7e51f0835491300f7695425d264"
BASE = "/tmp/knowslink-humanpeer001-254ff6a98b12-c30163b6"


def load():
    spec = importlib.util.spec_from_file_location("hp_remote", HELPER / "remote.py")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def main():
    action = sys.argv[1]
    if action not in ("origin", "auth", "fixture"):
        raise SystemExit("usage: drive_fixture.py origin|auth|fixture")
    record = json.loads((HELPER / f"{SHA}.json").read_text(encoding="utf-8"))
    if record.get("sha") != SHA or record.get("base") != BASE:
        raise SystemExit("recorded environment does not match the final SHA")
    remote = load()
    client, values = remote.connect()
    try:
        prefix = (
            f"cd {shlex.quote(BASE)}/repo && test \"$(git rev-parse HEAD)\" = {SHA} "
            "&& test -z \"$(git status --porcelain)\" && "
        )
        if action == "origin":
            command = prefix + (
                "python3 -c \"import json,pathlib,urllib.request; "
                f"base=pathlib.Path({BASE!r}); "
                "rec=max(base.glob('qa-*/runtime.json'), key=lambda p: p.stat().st_mtime_ns); "
                "meta=json.loads(rec.read_text()); "
                "assert meta['sha']==%r; "
                "opener=urllib.request.urlopen(meta['origin']+'/healthz', timeout=5); "
                "print('REMOTE_SHA='+meta['sha']); print('ORIGIN='+meta['origin']); "
                "print('HEALTH='+str(opener.status))\""
            ) % SHA
            return remote.run(client, values, command, f"{SHA}-origin")
        script = {"auth": "auth_remote.py", "fixture": "human_fixture_remote.py"}[action]
        remote_script = BASE + "/" + script
        with client.open_sftp() as sftp:
            sftp.put(str(HELPER / script), remote_script)
        return remote.run(client, values, prefix + f"python3 {shlex.quote(remote_script)} {shlex.quote(BASE)} {SHA}", f"{SHA}-{action}")
    finally:
        client.close()


if __name__ == "__main__":
    try:
        sys.exit(main())
    except (OSError, AssertionError) as error:
        print(f"FAILED: {type(error).__name__}; sensitive details suppressed", file=sys.stderr)
        sys.exit(1)
