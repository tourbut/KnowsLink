import json, os, subprocess, sys
G = sys.argv[1]
WS = os.path.dirname(os.path.abspath(__file__))


def run(cmd):
    env = dict(os.environ, CLAUDE_PROJECT_DIR=WS)
    r = subprocess.run([sys.executable, G], input=json.dumps({'tool_name': 'PowerShell', 'tool_input': {'command': cmd}, 'cwd': WS}),
                       capture_output=True, text=True, env=env, timeout=15)
    return {2: 'DENY', 0: 'ALLOW'}.get(r.returncode, r.returncode)


for label, cmd in [
    ('-r:$true + variable target', 'Remove-Item -r:$true $env:X\\sub'),
    ('-r:$true + wildcard', 'Remove-Item -r:$true C:\\Users\\x\\*'),
    ('-r:$true + drive root', 'Remove-Item -r:$true C:\\'),
    ('-r:$true + dq var', 'Remove-Item -r:$true "$env:X\\"'),
    ('-r:$true + literal outside path', 'Remove-Item -r:$true C:\\Users\\x\\Documents'),
    ('git clean -n -fd', 'git clean -n -fd'),
]:
    print(f'{run(cmd):6} {label}')
