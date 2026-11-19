"""Run exact-HEAD Linux checks in a new task-owned clone and preserve synthetic resources.

Known host keys are required; private SSH values stay in memory and never enter argv or logs.
Copied from SAR-UX-ISSUES-002-CONNECT-DEV-test/remote.py with a separate private tool directory and the language fixture names.
"""
import json
import pathlib
import secrets as random
import subprocess
import sys

sys.path.insert(0, 'C:/Users/shin/AppData/Local/KnowsLinkDevTools/ssh-python')
import paramiko

ROOT = pathlib.Path('C:/Users/shin/orca/workspaces/KnowsLink/dev')
HERE = pathlib.Path('C:/Users/shin/AppData/Local/KnowsLinkDevTools/lang-ux002-dev')
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
            base = f'/tmp/knowslink-lang002-{sha[:12]}-{random.token_hex(4)}'
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
        if not base.startswith(f'/tmp/knowslink-lang002-{sha[:12]}-') or not all(c.isalnum() or c in '/-' for c in base):
            raise ValueError('unexpected task clone')
        clean = f'cd {base}/repo && test "$(git rev-parse HEAD)" = {sha} && test -z "$(git status --porcelain)"'
        if action == 'make':
            target = sys.argv[2]
            if target not in ('lint', 'test', 'build'):
                raise ValueError('unsupported target')
            return run(f'{clean} && make {target}', target)
        if action == 'fixture-compile':
            return run(f'{clean} && go vet -tags=integration,langfixture ./internal/relay && go test -c -tags=integration,langfixture -o {base}/lang.test ./internal/relay && test -x {base}/lang.test', 'fixture-compile')
        if action == 'integration' or action == 'fixture':
            name = 'kl-lang002-' + action + '-' + sha[:12] + '-' + random.token_hex(4)
            (HERE / f'{sha}-{action}-database.txt').write_text(name, encoding='utf-8')
            prepare = (f'set -e; {clean}; '
                       f'docker run -d --name {name} --cpus=1 --memory=512m -e POSTGRES_PASSWORD=synthetic-lang-test '
                       f'-e POSTGRES_USER=knowslink -e POSTGRES_DB=knowslink -p 127.0.0.1::5432 {IMAGE}; '
                       f'for i in $(seq 1 40); do if docker exec {name} pg_isready -h 127.0.0.1 -U knowslink -d knowslink >/dev/null; then break; fi; sleep 1; done; '
                       f'port=$(docker port {name} 5432/tcp); port=${{port##*:}}; '
                       'export TEST_DATABASE_URL=postgres://knowslink:synthetic-lang-test@127.0.0.1:$port/knowslink?sslmode=disable; '
                       'export TEST_SYNTHETIC_DATABASE=1; export DATABASE_URL="$TEST_DATABASE_URL"; go run ./cmd/migrate; ')
            if action == 'integration':
                only = sys.argv[2] if len(sys.argv) > 2 else ''
                if only and not all(c.isalnum() or c in '^$|' for c in only):
                    raise ValueError('unsupported test filter')
                selection = f'-run "{only}"' if only else '-skip "TestRemoteConsentBrowserReturn"'
                label = 'integration-' + only.strip('^$') if only else 'integration'
                try:
                    return run(prepare + f'go test -race -tags=integration ./internal/relay {selection} -count=1 -v', label)
                finally:
                    run(f'docker stop {name}', label + '-stop')
            command = (prepare + f'go test -c -tags=integration,langfixture -o {base}/lang.test ./internal/relay; '
                       f'export LANG_FIXTURE_MANIFEST={base}/fixture.json; '
                       f'nohup {base}/lang.test -test.run=^TestLanguageBrowserFixture$ -test.timeout=25h -test.v > {base}/fixture.log 2>&1 < /dev/null & echo $! > {base}/fixture.pid')
            code = run(command, 'fixture-start')
            if code == 0:
                (HERE / 'fixture-location.json').write_text(json.dumps({'sha': sha, 'base': base, 'container': name}), encoding='utf-8')
            return code
        if action == 'probe':
            # Loopback-only checks of the running fixture: manifest, sign-in as the English member, English pages, the
            # anonymous language link and a refused path.
            probe = r'''set -e; m=BASE/fixture.json; test -s $m; stat -c %a $m; tail -n 3 BASE/fixture.log; kill -0 $(cat BASE/fixture.pid)
o=$(python3 -c "import json,sys;print(json.load(open(sys.argv[1]))['origin'])" $m)
case "$o" in http://127.0.0.1:*) ;; *) exit 9;; esac
j=BASE/probe-cookies.txt
lang() { grep -o '<html lang="[a-z]*"' | head -n 1; }
echo "GET / $(curl -s -o /dev/null -w %{http_code} $o/) $(curl -s $o/ | lang)"
echo "GET /lang/en $(curl -s -c $j -o /dev/null -w '%{http_code} %{redirect_url}' "$o/lang/en?back=%2F")"
echo "GET / with cookie $(curl -s -b $j $o/ | lang)"
echo "POST /fixture/as/en $(curl -s -c $j -o /dev/null -w %{http_code} -X POST $o/fixture/as/en)"
for p in /home/account /home/friends /home/agents /home/legacy /home/receipts /home/new; do echo "GET $p $(curl -s -b $j -o /dev/null -w %{http_code} "$o$p") $(curl -s -b $j "$o$p" | lang)"; done
echo "POST /fixture/as/ko $(curl -s -c $j -o /dev/null -w %{http_code} -X POST $o/fixture/as/ko)"
echo "GET /home/account as ko $(curl -s -b $j $o/home/account | lang)"
echo "POST /fixture/as/none $(curl -s -c $j -o /dev/null -w %{http_code} -X POST $o/fixture/as/none)"
echo "GET /home as none $(curl -s -b $j $o/home | lang)"
echo "GET /v1/pull $(curl -s -o /dev/null -w %{http_code} $o/v1/pull)"
'''.replace('BASE', base)
            return run(probe, 'fixture-probe')
        raise ValueError('usage: prepare|make lint|test|build|integration [filter]|fixture-compile|fixture|probe')
    finally:
        client.close()


if __name__ == '__main__':
    raise SystemExit(main())
