import json, os, subprocess, sys, tempfile
G = sys.argv[1]
WS = os.path.dirname(os.path.abspath(__file__))
def run(cmd):
    env = dict(os.environ, CLAUDE_PROJECT_DIR=WS)
    r = subprocess.run([sys.executable, G], input=json.dumps({'tool_name': 'PowerShell', 'tool_input': {'command': cmd}, 'cwd': WS}),
                       capture_output=True, text=True, env=env, timeout=15)
    return {2: 'DENY', 0: 'ALLOW'}.get(r.returncode, r.returncode)
cases = [
    ('recurse switch colon abbreviated -r:$true', 'Remove-Item -r:$true C:\\Users\\x\\Documents'),
    ('recurse switch -re:$true', 'Remove-Item C:\\Users\\x\\Documents -re:$true'),
    ('recurse switch -r:1 literal path', "Remove-Item -LiteralPath 'C:\\Users\\x\\Documents' -r:1"),
    ('recurse full -Recurse:$true', 'Remove-Item C:\\Users\\x\\Documents -Recurse:$true'),
    ('ps comment hides delete (trailing backslash in dq)', 'Get-ChildItem "C:\\d\\"; Remove-Item -Recurse C:\\x # "'),
    ('ps comment hides delete (no sep)', 'echo "a\\" ; rm -r C:\\x # "'),
    ('dot-source cmdlet', '. Remove-Item -Recurse C:\\x'),
    ('rm -r:$true', 'rm -r:$true C:\\x'),
    ('del -r:$true via alias', 'del C:\\x -r:$true'),
    ('rm -Recurse via splat var', 'Remove-Item C:\\x @{Recurse=$true}'),
    ('rm array of paths', 'Remove-Item C:\\x1,C:\\x2 -Recurse'),
    ('git -C clean', 'git -C . clean -fdx'),
    ('python os.system', "python -c \"import os; os.system('rm -rf /')\""),
    ('node rmdirSync', "node -e \"require('fs').rmdirSync('x',{recursive:true})\""),
    ('quoted tilde expansion path', "Remove-Item -LiteralPath 'C:\\ws\\a' -Recurse"),
]
for label, cmd in cases:
    print(f'{run(cmd):6} {label}')
