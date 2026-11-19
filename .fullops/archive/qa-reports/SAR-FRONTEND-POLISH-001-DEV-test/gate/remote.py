"""SAR-FRONTEND-POLISH-001-DEV task-local (copy of wsg001fix) SSH helper: unique /tmp clone per exact HEAD, HEAD/clean check, exit code passthrough.

Host keys must already be known; credentials stay in memory and are redacted from logs. Nothing is deleted locally or remotely.
"""
import pathlib
import secrets as rnd
import subprocess
import sys
sys.path.insert(0, 'C:/Users/shin/AppData/Local/KnowsLinkDevTools/ssh-python')
import paramiko

ROOT = pathlib.Path('C:/Users/shin/orca/workspaces/KnowsLink/dev')
HERE = pathlib.Path(__file__).parent
secrets = {}
for line in (ROOT / '.env.server').read_text().splitlines():
    if '=' in line and not line.lstrip().startswith('#'):
        key, value = line.split('=', 1)
        secrets[key.strip()] = value.strip().strip('"').strip("'")
client = paramiko.SSHClient()
client.load_system_host_keys()
client.set_missing_host_key_policy(paramiko.RejectPolicy())
try:
    client.connect(secrets['IP'], username=secrets['USER'], password=secrets['PW'], look_for_keys=False, allow_agent=False, timeout=20)
except (paramiko.SSHException, OSError):
    raise SystemExit('SSH connection/host-key verification failed; credentials not printed')


def run(command, log):
    stdin, stdout, stderr = client.exec_command(command, timeout=1800)
    stdin.close()
    output = stdout.read().decode('utf-8', 'replace') + stderr.read().decode('utf-8', 'replace')
    code = stdout.channel.recv_exit_status()
    for secret in secrets.values():
        if secret:
            output = output.replace(secret, '[redacted]')
    log.write_text(output + f'\nEXIT={code}\n', encoding='utf-8')
    print(output[-4000:])
    print('EXIT=', code)
    return code


def head():
    sha = subprocess.check_output(['git', '-C', str(ROOT), 'rev-parse', 'HEAD'], text=True).strip()
    assert len(sha) == 40 and all(c in '0123456789abcdef' for c in sha)
    return sha


try:
    action, sha = sys.argv[1], head()
    base_file = HERE / f'{sha}.base'
    if action == 'prepare':
        base = f'/tmp/knowslink-frontpolish001-{sha[:12]}-{rnd.token_hex(4)}'
        bundle = HERE / f'{sha}.bundle'
        subprocess.run(['git', '-C', str(ROOT), 'bundle', 'create', str(bundle), 'HEAD'], check=True)
        run(f'mkdir -m 700 {base}', HERE / f'{sha}-mkdir.log') == 0 or sys.exit(1)
        with client.open_sftp() as sftp:
            sftp.put(str(bundle), f'{base}/head.bundle')
        code = run(f'set -e; cd {base}; git clone -q head.bundle repo; cd repo; git checkout -q {sha}; '
                   f'test "$(git rev-parse HEAD)" = {sha}; make install', HERE / f'{sha}-install.log')
        if code == 0:
            base_file.write_text(base)
        sys.exit(code)
    elif action == 'make':
        target = sys.argv[2]
        if target not in ('lint', 'test'):
            raise SystemExit('unsupported target')
        base = base_file.read_text().strip()
        sys.exit(run(f'cd {base}/repo && test "$(git rev-parse HEAD)" = {sha} && test -z "$(git status --porcelain)" && make {target}',
                     HERE / f'{sha}-{target}.log'))
    else:
        raise SystemExit('usage: remote.py prepare | make lint|test')
finally:
    client.close()
