"""Start an exact-SHA Linux clone and the synthetic calendar fixture for SAR-SCHEDULE-001-QA.

The clone and its Postgres container belong to this task. SSH values stay in memory and are
removed from saved output. No operational host, database, or Google account is used.
"""
import json
import pathlib
import secrets as random
import subprocess
import sys

sys.path.insert(0, 'C:/Users/shin/AppData/Local/KnowsLinkDevTools/ssh-python')
import paramiko

ROOT = pathlib.Path('C:/Users/shin/orca/workspaces/KnowsLink/tester')
HERE = pathlib.Path('C:/Users/shin/AppData/Local/KnowsLinkDevTools/schedule001-qa')
SHA = '336533e6e15679d67118cbec5ab907768b80a70f'
IMAGE = 'postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24'
RUN = 'TestCalendarConnection|TestCalendarApprovalRunsOneLookup|TestCalendarUnconnectedAndCommitStub|TestCalendarScopeEnforcedAtAccept|TestCalendarRevocationAndFailureClasses'


def connect():
    values = {}
    for line in (ROOT / '.env.server').read_text(encoding='utf-8').splitlines():
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
    client, values = connect()

    def run(command, label, timeout=1800):
        stdin, stdout, stderr = client.exec_command(command, timeout=timeout)
        stdin.close()
        output = stdout.read().decode('utf-8', 'replace') + stderr.read().decode('utf-8', 'replace')
        code = stdout.channel.recv_exit_status()
        for value in values.values():
            if value:
                output = output.replace(value, '[redacted]')
        (HERE / f'{SHA}-{label}.log').write_text(output + f'\nEXIT={code}\n', encoding='utf-8')
        print(f'{label}: EXIT={code}', flush=True)
        return code

    try:
        action = sys.argv[1]
        base_file = HERE / f'{SHA}.base'
        if action == 'prepare':
            base = f'/tmp/knowslink-schedqa-{SHA[:12]}-{random.token_hex(4)}'
            bundle = HERE / f'{SHA}.bundle'
            # A raw SHA is not a ref, so git refuses an empty bundle. The branch contains that commit.
            subprocess.run(['git', '-C', str(ROOT), 'bundle', 'create', str(bundle), 'refs/heads/fullops/tester'], check=True)
            if run(f'mkdir -m 700 {base}', 'mkdir') != 0:
                return 1
            with client.open_sftp() as sftp:
                sftp.put(str(bundle), f'{base}/head.bundle')
            code = run(
                f'set -e; cd {base}; git clone -q head.bundle repo; cd repo; git checkout -q {SHA}; '
                f'test "$(git rev-parse HEAD)" = {SHA}; make install',
                'install',
                timeout=2400,
            )
            if code == 0:
                base_file.write_text(base, encoding='utf-8')
            return code
        base = base_file.read_text(encoding='utf-8').strip()
        if not base.startswith(f'/tmp/knowslink-schedqa-{SHA[:12]}-') or not all(c.isalnum() or c in '/-' for c in base):
            raise ValueError('unexpected task clone')
        clean = f'cd {base}/repo && test "$(git rev-parse HEAD)" = {SHA} && test -z "$(git status --porcelain)"'
        if action == 'make':
            target = sys.argv[2]
            if target not in ('lint', 'test'):
                raise ValueError('unsupported target')
            return run(f'{clean} && make {target}', target, timeout=2400)
        if action == 'integration':
            name = 'kl-schedqa-integration-' + SHA[:12] + '-' + random.token_hex(4)
            (HERE / f'{SHA}-integration-database.txt').write_text(name, encoding='utf-8')
            prepare = (
                f'set -e; {clean}; '
                f'docker run -d --name {name} --cpus=1 --memory=512m -e POSTGRES_PASSWORD=synthetic-schedule-test '
                f'-e POSTGRES_USER=knowslink -e POSTGRES_DB=knowslink -p 127.0.0.1::5432 {IMAGE}; '
                f'for i in $(seq 1 40); do if docker exec {name} pg_isready -h 127.0.0.1 -U knowslink -d knowslink >/dev/null; then break; fi; sleep 1; done; '
                f'port=$(docker port {name} 5432/tcp); port=${{port##*:}}; '
                'export TEST_DATABASE_URL=postgres://knowslink:synthetic-schedule-test@127.0.0.1:$port/knowslink?sslmode=disable; '
                'export TEST_SYNTHETIC_DATABASE=1; export DATABASE_URL="$TEST_DATABASE_URL"; go run ./cmd/migrate; '
            )
            try:
                return run(
                    prepare + f'go test -race -tags=integration -v ./internal/relay -count=1 -timeout 40m -run "^({RUN})$"',
                    'integration',
                    timeout=2400,
                )
            finally:
                run(f'docker stop {name}', 'integration-stop')
        if action == 'fixture':
            name = 'kl-schedqa-fixture-' + SHA[:12] + '-' + random.token_hex(4)
            (HERE / f'{SHA}-fixture-database.txt').write_text(name, encoding='utf-8')
            command = (
                f'set -e; {clean}; '
                f'docker run -d --name {name} --cpus=1 --memory=512m -e POSTGRES_PASSWORD=synthetic-schedule-test '
                f'-e POSTGRES_USER=knowslink -e POSTGRES_DB=knowslink -p 127.0.0.1::5432 {IMAGE}; '
                f'for i in $(seq 1 40); do if docker exec {name} pg_isready -h 127.0.0.1 -U knowslink -d knowslink >/dev/null; then break; fi; sleep 1; done; '
                f'port=$(docker port {name} 5432/tcp); port=${{port##*:}}; '
                'export TEST_DATABASE_URL=postgres://knowslink:synthetic-schedule-test@127.0.0.1:$port/knowslink?sslmode=disable; '
                'export TEST_SYNTHETIC_DATABASE=1; export DATABASE_URL="$TEST_DATABASE_URL"; go run ./cmd/migrate; '
                f'go test -c -tags=integration,calendarfixture -o {base}/calendar.test ./internal/relay; '
                f'export CALENDAR_FIXTURE_MANIFEST={base}/fixture.json; '
                f'nohup {base}/calendar.test -test.run=^TestCalendarBrowserFixture$ -test.timeout=25h -test.v > {base}/fixture.log 2>&1 < /dev/null & echo $! > {base}/fixture.pid'
            )
            code = run(command, 'fixture-start', timeout=2400)
            if code == 0:
                (HERE / 'fixture-location.json').write_text(json.dumps({'sha': SHA, 'base': base, 'container': name}), encoding='utf-8')
            return code
        if action == 'wait':
            command = (
                f'set -e; for i in $(seq 1 90); do if test -s {base}/fixture.json && grep -q "calendar fixture ready" {base}/fixture.log; then break; fi; sleep 2; done; '
                f'test -s {base}/fixture.json; kill -0 $(cat {base}/fixture.pid); tail -n 5 {base}/fixture.log'
            )
            code = run(command, 'fixture-wait', timeout=300)
            if code == 0:
                with client.open_sftp() as sftp:
                    sftp.get(f'{base}/fixture.json', str(HERE / 'fixture.json'))
            return code
        if action == 'cleanup':
            location = json.loads((HERE / 'fixture-location.json').read_text(encoding='utf-8'))
            name = location['container']
            if not name.startswith('kl-schedqa-fixture-'):
                raise ValueError('unexpected container')
            return run(
                f'set +e; if test -f {base}/fixture.pid; then kill $(cat {base}/fixture.pid); fi; docker stop {name}; echo stopped',
                'cleanup',
                timeout=180,
            )
        raise ValueError('usage: prepare|make lint|test|integration|fixture|wait|cleanup')
    finally:
        client.close()


if __name__ == '__main__':
    raise SystemExit(main())
