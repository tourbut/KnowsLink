"""Task-local exact-HEAD SSH checks; synthetic DB and clone artifacts remain for coordinator cleanup.

The coordinator approved this non-destructive file edit after reviewing the blocked inline generator.
Host keys must already be known. Credentials remain in memory and are redacted from logs.
"""
import pathlib
import secrets as rnd
import subprocess
import sys

sys.path.insert(0, 'C:/Users/shin/AppData/Local/KnowsLinkDevTools/ssh-python')
import paramiko

ROOT = pathlib.Path('C:/Users/shin/orca/workspaces/KnowsLink/dev')
HERE = pathlib.Path('C:/Users/shin/AppData/Local/KnowsLinkDevTools/refactor002')
BASELINE = '9255852fa0e46cd40b7be1021eaa1f6c7f7d083b'
IMAGE = 'postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24'


def self_check():
    source = pathlib.Path(__file__).read_text(encoding='utf-8')
    compile(source, str(__file__), 'exec')
    assert 'RejectPolicy' in source and 'git status --porcelain' in source
    checked = '\n'.join(line for line in source.splitlines() if not line.lstrip().startswith('assert '))
    assert 'docker rm' not in checked
    assert 'rmtree' not in checked
    print('helper syntax and security checks passed')


if sys.argv[1:] == ['self-check']:
    self_check()
    raise SystemExit(0)

HERE.mkdir(parents=True, exist_ok=True)
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
    print(f'Log: {log.name}; EXIT={code}')
    return code


def head():
    sha = subprocess.check_output(['git', '-C', str(ROOT), 'rev-parse', 'HEAD'], text=True).strip()
    assert len(sha) == 40 and all(c in '0123456789abcdef' for c in sha)
    return sha


try:
    action, sha = sys.argv[1], head()
    base_file = HERE / f'{sha}.base'
    if action == 'probe':
        sys.exit(run('set -e; go version; node --version; npm --version; docker --version', HERE / 'environment.log'))
    if action == 'prepare':
        base = f'/tmp/knowslink-refactor002-{sha[:12]}-{rnd.token_hex(4)}'
        bundle = HERE / f'{sha}.bundle'
        subprocess.run(['git', '-C', str(ROOT), 'bundle', 'create', str(bundle), 'HEAD'], check=True)
        if run(f'mkdir -m 700 {base}', HERE / f'{sha}-mkdir.log') != 0:
            sys.exit(1)
        with client.open_sftp() as sftp:
            sftp.put(str(bundle), f'{base}/head.bundle')
        code = run(f'set -e; cd {base}; git clone -q head.bundle repo; cd repo; git checkout -q {sha}; '
                   f'test "$(git rev-parse HEAD)" = {sha}; make install', HERE / f'{sha}-install.log')
        if code == 0:
            base_file.write_text(base)
        sys.exit(code)
    base = base_file.read_text().strip()
    assert base.startswith(f'/tmp/knowslink-refactor002-{sha[:12]}-') and all(c.isalnum() or c in '/-' for c in base)
    clean = f'cd {base}/repo && test "$(git rev-parse HEAD)" = {sha} && test -z "$(git status --porcelain)"'
    if action == 'environment':
        sys.exit(run(f'{clean} && go version && node --version && npm --version', HERE / f'{sha}-environment.log'))
    if action == 'make':
        target = sys.argv[2]
        if target not in ('lint', 'test', 'build'):
            raise SystemExit('unsupported target')
        sys.exit(run(f'{clean} && make {target}', HERE / f'{sha}-{target}.log'))
    if action == 'integration':
        name = 'kl-refactor002-' + sha[:12] + '-' + rnd.token_hex(4)
        (HERE / f'{name}-database.txt').write_text(name)
        command = (f'set -e; {clean}; test -f adapters/dist/connect.js; '
                   f'docker run -d --name {name} --cpus=1 --memory=512m -e POSTGRES_PASSWORD=synthetic-refactor-test '
                   f'-e POSTGRES_USER=knowslink -e POSTGRES_DB=knowslink -p 127.0.0.1::5432 {IMAGE}; '
                   f'for i in $(seq 1 40); do if docker exec {name} pg_isready -U knowslink -d knowslink >/dev/null; then break; fi; sleep 1; done; '
                   f'port=$(docker port {name} 5432/tcp); port=${{port##*:}}; '
                   'export TEST_DATABASE_URL=postgres://knowslink:synthetic-refactor-test@127.0.0.1:$port/knowslink?sslmode=disable; '
                   'export TEST_SYNTHETIC_DATABASE=1; export DATABASE_URL="$TEST_DATABASE_URL"; '
                   'go run ./cmd/migrate; go test -race -tags=integration ./internal/relay -count=1; '
                   f'cd {base}; git clone -q head.bundle baseline; cd baseline; git checkout -q {BASELINE}; '
                   f'cp {base}/repo/internal/relay/store_retention*_test.go internal/relay/; '
                   'go test -tags=integration ./internal/relay -run "^$" -bench "^BenchmarkRetained(Postgres|History)$" -benchmem -benchtime=5x -count=3 > ../before-linux.txt; '
                   f'cd {base}/repo; '
                   'go test -tags=integration ./internal/relay -run "^$" -bench "^BenchmarkRetained(Postgres|History)$" -benchmem -benchtime=5x -count=3 > ../after-linux.txt; '
                   f'docker stop {name}')
        code = run(command, HERE / f'{name}-integration.log')
        with client.open_sftp() as sftp:
            for label in ('before', 'after'):
                try:
                    sftp.get(f'{base}/{label}-linux.txt', str(HERE / f'{sha}-{label}-linux.txt'))
                except FileNotFoundError:
                    if code == 0:
                        raise
        sys.exit(code)
    raise SystemExit('usage: remote.py self-check|probe|prepare|integration|make lint|test|build')
finally:
    client.close()
