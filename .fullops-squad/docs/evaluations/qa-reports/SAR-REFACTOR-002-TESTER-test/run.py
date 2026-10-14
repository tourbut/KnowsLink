"""Run fixed-product independent QA in a unique Linux clone; retain every fixture and native exit code."""
import json
import pathlib
import secrets
import subprocess
import sys

sys.path.insert(0, 'C:/Users/shin/AppData/Local/KnowsLinkDevTools/ssh-python')
import paramiko

ROOT = pathlib.Path(__file__).resolve().parents[5]
OUT = pathlib.Path(__file__).resolve().parent
SHA = 'bcb1dd6629b508a30ec39a5d3141c8ce0a558597'
BASE = '8f9370bb62c9133c73e10ea2307b4cd696e10f7e'
IMAGE = 'postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24'
LOCAL = pathlib.Path('C:/Users/shin/AppData/Local/KnowsLinkDevTools/qa-refactor002')
LOCAL.mkdir(exist_ok=True)
auth = {}
for line in (ROOT / '.env.server').read_text().splitlines():
    if '=' in line and not line.lstrip().startswith('#'):
        key, value = line.split('=', 1)
        auth[key.strip()] = value.strip().strip('"').strip("'")
client = paramiko.SSHClient()
client.load_system_host_keys()
client.set_missing_host_key_policy(paramiko.RejectPolicy())
client.connect(auth['IP'], username=auth['USER'], password=auth['PW'], look_for_keys=False, allow_agent=False, timeout=20)


def run(name, command):
    stdin, stdout, stderr = client.exec_command(command, timeout=900)
    stdin.close()
    output = stdout.read().decode('utf-8', 'replace') + stderr.read().decode('utf-8', 'replace')
    code = stdout.channel.recv_exit_status()
    for value in auth.values():
        if value:
            output = output.replace(value, '[redacted]')
    (OUT / (name + '.log')).write_text(output + f'\nEXIT={code}\n', encoding='utf-8')
    with (OUT / 'commands.jsonl').open('a', encoding='utf-8') as log:
        log.write(json.dumps({'name': name, 'product_sha': SHA, 'command': command, 'exit': code}) + '\n')
    print(f'{name}: EXIT={code}', flush=True)
    if code:
        raise SystemExit(code)


try:
    action = sys.argv[1]
    if action in ('gate-prepare', 'make'):
        OUT = LOCAL / 'gate'
        OUT.mkdir(exist_ok=True)
        current = subprocess.check_output(['git', '-C', str(ROOT), 'rev-parse', 'HEAD'], text=True).strip()
        assert len(current) == 40 and all(c in '0123456789abcdef' for c in current)
        SHA = current
        if action == 'gate-prepare':
            base = '/tmp/knowslink-qa-gate-' + secrets.token_hex(6)
            bundle = LOCAL / 'qa.bundle'
            subprocess.run(['git', '-C', str(ROOT), 'bundle', 'create', str(bundle), 'HEAD'], check=True)
            run('mkdir', f'mkdir -m 700 {base}')
            with client.open_sftp() as sftp:
                sftp.put(str(bundle), base + '/qa.bundle')
            (LOCAL / 'gate-remote.txt').write_text(base)
            run('prepare', f'set -e; cd {base}; git clone -q qa.bundle repo; cd repo; git checkout -q {current}; '
                f'test "$(git rev-parse HEAD)" = {current}; test -z "$(git status --porcelain)"; make install')
        else:
            base = (LOCAL / 'gate-remote.txt').read_text().strip()
            assert base.startswith('/tmp/knowslink-qa-gate-') and all(c.isalnum() or c in '/-' for c in base)
            target = sys.argv[2]
            assert target in ('lint', 'test')
            run(target, f'set -e; cd {base}/repo; test "$(git rev-parse HEAD)" = {current}; '
                f'test -z "$(git status --porcelain)"; make {target}')
        raise SystemExit(0)
    if action == 'prepare':
        base = '/tmp/knowslink-qa-refactor002-' + secrets.token_hex(6)
        bundle = LOCAL / 'product.bundle'
        subprocess.run(['git', '-C', str(ROOT), 'bundle', 'create', str(bundle), SHA, '--all'], check=True)
        run('mkdir', f'mkdir -m 700 {base}')
        with client.open_sftp() as sftp:
            sftp.put(str(bundle), base + '/product.bundle')
        (LOCAL / 'remote.txt').write_text(base)
        run('prepare', f'set -e; cd {base}; git clone -q product.bundle repo; cd repo; git checkout -q {SHA}; '
            f'test "$(git rev-parse HEAD)" = {SHA}; test -z "$(git status --porcelain)"; '
            'git rev-parse HEAD; go version; go env CGO_ENABLED; node --version; npm --version; '
            'npm ci --prefix adapters; npm run build --prefix adapters')
        raise SystemExit(0)
    base = (LOCAL / 'remote.txt').read_text().strip()
    assert base.startswith('/tmp/knowslink-qa-refactor002-') and all(c.isalnum() or c in '/-' for c in base)
    clean = f'cd {base}/repo; test "$(git rev-parse HEAD)" = {SHA}; test -z "$(git status --porcelain)"; '
    if action == 'qa':
        with client.open_sftp() as sftp:
            for name in ('independent_test.go', 'inbox-boundaries.mjs'):
                sftp.put(str(OUT / name), base + '/' + name)
            overlay = {'Replace': {base + '/repo/internal/relay/qa_independent_test.go': base + '/independent_test.go'}}
            with sftp.open(base + '/overlay.json', 'w') as handle:
                handle.write(json.dumps(overlay))
        run('qa-race', 'set -e; ' + clean + f'go test -race -overlay {base}/overlay.json ./internal/relay -count=1 -v '
            '-run "^(TestQA.*|TestSweepExpiryPreservesRetentionAndRevocation|TestRemote.*|TestTextEvent.*|TestHumanConversationLifetimeReplayAndRevocation|TestHumanRemoteOwnerAndAgentAuthorship)$"')
        run('qa-inbox', 'set -e; ' + clean + 'node adapters/dist/inbox-files.test.js; node adapters/dist/owner-text.test.js; '
            f'node {base}/inbox-boundaries.mjs {base}/repo')
    elif action == 'database':
        name = 'kl-qa-refactor002-' + secrets.token_hex(6)
        (LOCAL / 'container.txt').write_text(name)
        try:
            run('database', 'set -e; ' + clean +
                f'docker run -d --name {name} --cpus=1 --memory=512m -e POSTGRES_PASSWORD=synthetic-qa-test '
                f'-e POSTGRES_USER=knowslink -e POSTGRES_DB=knowslink -p 127.0.0.1::5432 {IMAGE}; '
                f'for i in $(seq 1 40); do if docker exec {name} pg_isready -U knowslink -d knowslink >/dev/null; then break; fi; sleep 1; done; '
                f'port=$(docker port {name} 5432/tcp); port=${{port##*:}}; '
                'export TEST_DATABASE_URL=postgres://knowslink:synthetic-qa-test@127.0.0.1:$port/knowslink?sslmode=disable; '
                'export TEST_SYNTHETIC_DATABASE=1; export DATABASE_URL="$TEST_DATABASE_URL"; '
                'go run ./cmd/migrate; go test -race -tags=integration ./internal/relay -count=1 -v '
                '-run "^(TestHumanPostgres|TestSharedExecutionClaimHTTP|TestHTTPConcurrencyAcrossInstances|TestRemoteOAuthHTTP)$"')
        finally:
            run('database-stop', f'docker stop {name}')
    elif action == 'benchmark':
        run('benchmark-prepare', f'set -e; cd {base}; git clone -q product.bundle baseline; cd baseline; git checkout -q {BASE}; '
            f'cp {base}/repo/internal/relay/store_retention_test.go internal/relay/; '
            f'cp {base}/repo/internal/relay/remote_events_test.go internal/relay/; git rev-parse HEAD')
        for label, folder in [('before', 'baseline'), ('after', 'repo')]:
            run('benchmark-' + label, f'set -e; cd {base}/{folder}; '
                'go test ./internal/relay -run "^$" -bench "^(BenchmarkRetainedHistory|BenchmarkRemoteClaimSubscriptions)$" '
                '-benchmem -benchtime=3x -count=3')
    elif action == 'verify':
        run('post-clean', 'set -e; ' + clean + 'git rev-parse HEAD; git status --porcelain')
    else:
        raise SystemExit('prepare|qa|database|benchmark|verify')
finally:
    client.close()
