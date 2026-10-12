"""Upload one local script and run it on the verified final-SHA remote clone.

The helper parent-HEAD assert is not used. Remote HEAD, a clean tree, known_hosts,
and RejectPolicy still apply. Secret values are not printed.
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
    if len(sys.argv) < 2:
        raise SystemExit("usage: drive_upload.py <local-script> [remote-args...]")
    local = pathlib.Path(sys.argv[1]).resolve()
    record = json.loads((HELPER / f"{SHA}.json").read_text(encoding="utf-8"))
    if record.get("sha") != SHA or record.get("base") != BASE:
        raise SystemExit("recorded environment does not match the final SHA")
    remote = load()
    client, values = remote.connect()
    try:
        remote_script = BASE + "/" + local.name
        with client.open_sftp() as sftp:
            sftp.put(str(local), remote_script)
        prefix = (
            f"cd {shlex.quote(BASE)}/repo && test \"$(git rev-parse HEAD)\" = {SHA} "
            "&& test -z \"$(git status --porcelain)\" && "
        )
        args = " ".join(shlex.quote(item) for item in sys.argv[2:])
        command = prefix + f"python3 {shlex.quote(remote_script)} {shlex.quote(BASE)} {SHA}"
        if args:
            command += " " + args
        return remote.run(client, values, command, f"{SHA}-{local.stem}")
    finally:
        client.close()


if __name__ == "__main__":
    try:
        sys.exit(main())
    except (OSError, AssertionError) as error:
        print(f"FAILED: {type(error).__name__}; sensitive details suppressed", file=sys.stderr)
        sys.exit(1)
