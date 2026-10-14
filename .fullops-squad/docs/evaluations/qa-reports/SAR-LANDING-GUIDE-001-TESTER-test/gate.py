"""Run registered gates on an exact QA commit, reusing installed dependencies without installation or cleanup."""
import importlib.util
import pathlib
import secrets
import subprocess
import sys

ROOT = pathlib.Path(__file__).resolve().parents[5]
HERE = ROOT / '.fullops-squad/tmp/landing-tester-gate'
SOURCE = '/tmp/knowslink-landing001-4795d619513f-725b85fc/repo'
spec = importlib.util.spec_from_file_location('dev_gate', ROOT / '.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/gate/remote.py')
helper = importlib.util.module_from_spec(spec)
spec.loader.exec_module(helper)
helper.ROOT = ROOT


def main():
    HERE.mkdir(parents=True, exist_ok=True)
    sha = subprocess.check_output(['git', '-C', str(ROOT), 'rev-parse', 'HEAD'], text=True).strip()
    assert len(sha) == 40 and all(c in '0123456789abcdef' for c in sha)
    client, private = helper.connect()

    def run(command, label):
        stdin, stdout, stderr = client.exec_command(command, timeout=1800)
        stdin.close()
        output = stdout.read().decode('utf-8', 'replace') + stderr.read().decode('utf-8', 'replace')
        code = stdout.channel.recv_exit_status()
        for value in private.values():
            if value:
                output = output.replace(value, '[redacted]')
        (HERE / f'{sha}-{label}.log').write_text(output + f'\nEXIT={code}\n', encoding='utf-8')
        print(f'{sha} {label}: EXIT={code}', flush=True)
        return code

    try:
        action = sys.argv[1]
        record = HERE / (sha + '.base')
        if action == 'prepare':
            remote = f'/tmp/knowslink-landing001-tester-{sha[:12]}-{secrets.token_hex(4)}'
            bundle = HERE / (sha + '.bundle')
            subprocess.run(['git', '-C', str(ROOT), 'bundle', 'create', str(bundle), 'HEAD'], check=True)
            if run(f'mkdir -m 700 {remote}', 'mkdir'):
                return 1
            with client.open_sftp() as sftp:
                sftp.put(str(bundle), remote + '/head.bundle')
            command = (f'git clone -q {remote}/head.bundle {remote}/repo && cd {remote}/repo && '
                       f'git checkout -q {sha} && test "$(git rev-parse HEAD)" = {sha} && '
                       f'test -d {SOURCE}/adapters/node_modules && '
                       f'cmp adapters/package-lock.json {SOURCE}/adapters/package-lock.json && '
                       f'ln -s {SOURCE}/adapters/node_modules adapters/node_modules')
            code = run(command, 'prepare')
            if code == 0:
                record.write_text(remote, encoding='utf-8')
            return code
        target = sys.argv[2]
        assert action == 'make' and target in ('lint', 'test')
        remote = record.read_text().strip()
        assert remote.startswith(f'/tmp/knowslink-landing001-tester-{sha[:12]}-') and all(c.isalnum() or c in '/-' for c in remote)
        return run(f'cd {remote}/repo && test "$(git rev-parse HEAD)" = {sha} && test -z "$(git status --porcelain)" && make {target}', target)
    finally:
        client.close()


if __name__ == '__main__':
    raise SystemExit(main())
