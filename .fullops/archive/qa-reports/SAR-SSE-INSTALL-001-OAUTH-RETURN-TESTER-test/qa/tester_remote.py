"""Tester-owned SSH helper (SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER). Existing host-key verification; credentials stay in memory."""
import pathlib
import sys
sys.path.insert(0, 'C:/Users/shin/AppData/Local/KnowsLinkDevTools/ssh-python')
import paramiko

ENV = pathlib.Path('D:/workspace/KnowsLink/.env.server')
secrets = {}
for line in ENV.read_text().splitlines():
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

def scrub(text):
    for v in secrets.values():
        if v:
            text = text.replace(v, '[redacted]')
    return text

def run(command, log=None, tail=6000):
    stdin, stdout, stderr = client.exec_command(command, timeout=1500)
    stdin.close()
    out = scrub(stdout.read().decode('utf-8', 'replace') + stderr.read().decode('utf-8', 'replace'))
    code = stdout.channel.recv_exit_status()
    if log:
        pathlib.Path(log).write_text(out + f'\nEXIT={code}\n', encoding='utf-8')
    print(out[-tail:])
    print('EXIT=', code)
    return code

try:
    action = sys.argv[1]
    if action == 'run':
        sys.exit(run(pathlib.Path(sys.argv[2]).read_text(), sys.argv[3] if len(sys.argv) > 3 else None))
    elif action == 'put':
        with client.open_sftp() as sftp:
            sftp.put(sys.argv[2], sys.argv[3])
    elif action == 'get':
        with client.open_sftp() as sftp:
            sftp.get(sys.argv[2], sys.argv[3])
finally:
    client.close()
