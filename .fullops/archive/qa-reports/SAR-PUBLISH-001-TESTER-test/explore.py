"""Run a remote python3 script against this task's loopback fixture. Usage: python explore.py <script.py>; prints stdout."""
import sys, remote
sha = remote.sha_head(); base = (remote.HERE / f'{sha}.base').read_text().strip()
c, v = remote.connect()
try:
    i, o, e = c.exec_command(f'BASE={base} python3 -', timeout=300)
    i.write(open(sys.argv[1], encoding='utf-8').read()); i.channel.shutdown_write()
    out = o.read().decode('utf-8', 'replace') + e.read().decode('utf-8', 'replace'); code = o.channel.recv_exit_status()
finally:
    c.close()
for x in v.values():
    if x: out = out.replace(x, '[redacted]')
sys.stdout.buffer.write((out + f'\nEXIT={code}\n').encode('utf-8'))
