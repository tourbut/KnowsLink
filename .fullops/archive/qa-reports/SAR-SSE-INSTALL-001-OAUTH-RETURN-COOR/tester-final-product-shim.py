"""Tester-owned `make` shim: runs the registered make target in the isolated Linux clone of the exact local HEAD."""
import pathlib, subprocess, sys
sys.path.insert(0, 'C:/Users/shin/AppData/Local/KnowsLinkDevTools/ssh-python')
import paramiko
target = sys.argv[1]
assert target in ('lint', 'test')
head = subprocess.check_output(['git', 'rev-parse', 'HEAD'], text=True).strip()
# Product bytes (everything outside .fullops-squad) must equal the fixed review SHA that the Linux clone holds.
subprocess.check_call(['git', 'diff', '--quiet', 'eac603b6a41d2ce002e698600ac66b8fc67ad762', 'HEAD', '--', '.', ':!.fullops-squad'])
head = 'eac603b6a41d2ce002e698600ac66b8fc67ad762'
assert subprocess.check_output(['git', 'status', '--porcelain'], text=True).strip() == ''
secrets = {}
for line in pathlib.Path('D:/workspace/KnowsLink/.env.server').read_text().splitlines():
    if '=' in line and not line.lstrip().startswith('#'):
        k, v = line.split('=', 1)
        secrets[k.strip()] = v.strip().strip('"').strip("'")
c = paramiko.SSHClient()
c.load_system_host_keys()
c.set_missing_host_key_policy(paramiko.RejectPolicy())
c.connect(secrets['IP'], username=secrets['USER'], password=secrets['PW'], look_for_keys=False, allow_agent=False, timeout=20)
B = '/tmp/knowslink-oauth001-tester-d7e2c059/lint'
si, so, se = c.exec_command(f'cd {B} && test "$(git rev-parse HEAD)" = {head} && test -z "$(git status --porcelain)" && make {target}', timeout=900)
si.close()
out = so.read().decode('utf-8', 'replace') + se.read().decode('utf-8', 'replace')
code = so.channel.recv_exit_status()
for v in secrets.values():
    if v:
        out = out.replace(v, '[redacted]')
pathlib.Path(f'{sys.argv[2] if len(sys.argv) > 2 else "."}').exists()
print(out[-4000:])
print('EXIT=', code)
c.close()
sys.exit(code)
