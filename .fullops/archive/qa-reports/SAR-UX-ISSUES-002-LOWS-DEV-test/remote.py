"""Run exact-HEAD Linux checks in a new task-owned clone and preserve synthetic resources.

Known host keys are required; private SSH values stay in memory and never enter argv or logs.
Copied from SAR-UX-ISSUES-002-LANG-DEV-test/remote.py with a separate private tool directory; fixture actions dropped and a
`red` action added that runs the new checks against the base product files in a fresh copy of the clone.
"""
import json
import pathlib
import secrets as random
import subprocess
import sys

sys.path.insert(0, 'C:/Users/shin/AppData/Local/KnowsLinkDevTools/ssh-python')
import paramiko

ROOT = pathlib.Path('C:/Users/shin/orca/workspaces/KnowsLink/dev')
HERE = pathlib.Path('C:/Users/shin/AppData/Local/KnowsLinkDevTools/lows-ux002-dev')
BASE_SHA = 'c3b51926dbae01bdb325a9d3be01ed1ad09275c0'
IMAGE = 'postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24'


def sha_head():
    sha = subprocess.check_output(['git', '-C', str(ROOT), 'rev-parse', 'HEAD'], text=True).strip()
    if len(sha) != 40 or any(c not in '0123456789abcdef' for c in sha):
        raise ValueError('invalid HEAD')
    return sha


def connect():
    values = {}
    for line in (ROOT / '.env.server').read_text().splitlines():
        if '=' in line and not line.lstrip().startswith('#'):
            key, value = line.split('=', 1)
            values[key.strip()] = value.strip().strip('"').strip("'")
    client = paramiko.SSHClient()
    client.load_system_host_keys()
    client.set_missing_host_key_policy(paramiko.RejectPolicy())
    try:
        client.connect(values['IP'], username=values['USER'], password=values['PW'], look_for_keys=False, allow_agent=False, timeout=20)
    except (paramiko.SSHException, OSError):
        raise SystemExit('SSH connection/host-key verification failed; credentials not printed')
    return client, values


def main():
    HERE.mkdir(parents=True, exist_ok=True)
    sha = sha_head()
    client, values = connect()

    def run(command, label):
        stdin, stdout, stderr = client.exec_command(command, timeout=1800)
        stdin.close()
        output = stdout.read().decode('utf-8', 'replace') + stderr.read().decode('utf-8', 'replace')
        code = stdout.channel.recv_exit_status()
        for value in values.values():
            if value:
                output = output.replace(value, '[redacted]')
        (HERE / f'{sha}-{label}.log').write_text(output + f'\nEXIT={code}\n', encoding='utf-8')
        print(f'{label}: EXIT={code}', flush=True)
        return code

    try:
        action = sys.argv[1]
        base_file = HERE / f'{sha}.base'
        if action == 'prepare':
            base = f'/tmp/knowslink-lows002-{sha[:12]}-{random.token_hex(4)}'
            bundle = HERE / f'{sha}.bundle'
            subprocess.run(['git', '-C', str(ROOT), 'bundle', 'create', str(bundle), 'HEAD'], check=True)
            if run(f'mkdir -m 700 {base}', 'mkdir') != 0:
                return 1
            with client.open_sftp() as sftp:
                sftp.put(str(bundle), f'{base}/head.bundle')
            code = run(f'set -e; cd {base}; git clone -q head.bundle repo; cd repo; git checkout -q {sha}; test "$(git rev-parse HEAD)" = {sha}; make install', 'install')
            if code == 0:
                base_file.write_text(base, encoding='utf-8')
            return code
        base = base_file.read_text().strip()
        if not base.startswith(f'/tmp/knowslink-lows002-{sha[:12]}-') or not all(c.isalnum() or c in '/-' for c in base):
            raise ValueError('unexpected task clone')
        clean = f'cd {base}/repo && test "$(git rev-parse HEAD)" = {sha} && test -z "$(git status --porcelain)"'
        if action == 'make':
            target = sys.argv[2]
            if target not in ('lint', 'test', 'build'):
                raise ValueError('unsupported target')
            return run(f'{clean} && make {target}', target)
        if action in ('integration', 'red'):
            name = 'kl-lows002-' + action + '-' + sha[:12] + '-' + random.token_hex(4)
            (HERE / f'{sha}-{action}-database.txt').write_text(name, encoding='utf-8')
            prepare = (f'set -e; {clean}; '
                       f'docker run -d --name {name} --cpus=1 --memory=512m -e POSTGRES_PASSWORD=synthetic-lows-test '
                       f'-e POSTGRES_USER=knowslink -e POSTGRES_DB=knowslink -p 127.0.0.1::5432 {IMAGE}; '
                       f'for i in $(seq 1 40); do if docker exec {name} pg_isready -h 127.0.0.1 -U knowslink -d knowslink >/dev/null; then break; fi; sleep 1; done; '
                       f'port=$(docker port {name} 5432/tcp); port=${{port##*:}}; '
                       'export TEST_DATABASE_URL=postgres://knowslink:synthetic-lows-test@127.0.0.1:$port/knowslink?sslmode=disable; '
                       'export TEST_SYNTHETIC_DATABASE=1; export DATABASE_URL="$TEST_DATABASE_URL"; go run ./cmd/migrate; ')
            if action == 'integration':
                command = prepare + 'go test -race -tags=integration ./internal/relay -skip "TestRemoteConsentBrowserReturn" -count=1 -v'
            else:
                # Base product files with the new tests in a fresh copy; the unit test naming the new wrapper stays at base.
                files = ' '.join('internal/relay/' + f for f in ('capacity.go', 'human.go', 'human_http.go', 'i18n.go', 'member_pages.go', 'i18n_test.go'))
                command = (prepare + f'red=$(mktemp -d {base}/red.XXXXXX); git clone -q {base}/repo $red/repo; cd $red/repo; '
                           f'git checkout -q {BASE_SHA} -- {files}; git status --short; '
                           'set +e; go test -tags=integration ./internal/relay -count=1 -v '
                           '-run "^(TestLanguagePreferencePostgres|TestFriendRerequestReusesOneRow|TestFriendBlockOnHeldRequestKeepsRequesterPending|TestFriendRepeatWithOpenLegacyPlaceholderIsNoop)$"; '
                           'code=$?; echo GO_TEST_EXIT=$code; test $code -ne 0')
            try:
                return run(command, action)
            finally:
                run(f'docker stop {name}', action + '-stop')
        raise ValueError('usage: prepare|make lint|test|build|integration|red')
    finally:
        client.close()


if __name__ == '__main__':
    raise SystemExit(main())
