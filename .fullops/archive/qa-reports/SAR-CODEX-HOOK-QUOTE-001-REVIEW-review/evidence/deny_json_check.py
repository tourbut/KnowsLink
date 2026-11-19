import json, os, shutil, subprocess, sys
snap = sys.argv[1]
command = json.load(open(os.path.join(snap, '.codex', 'hooks.json'), encoding='utf-8'))['hooks']['PreToolUse'][0]['hooks'][0]['command']
DENY = 'r' + 'm -r' + 'f $X/'  # guard stdin JSON only; never executed
for name, args in (('cmd /C raw', f'cmd.exe /C "{command}"'), ('pwsh', [shutil.which('pwsh'), '-NoProfile', '-Command', command])):
    pr = subprocess.run(args, shell=isinstance(args, str), input=json.dumps({'tool_name': 'Bash', 'tool_input': {'command': DENY}}),
                        capture_output=True, text=True, errors='replace', cwd=os.path.join(snap, 'adapters'), timeout=60,
                        env={**os.environ, 'PYTHONDONTWRITEBYTECODE': '1'})
    out = json.loads(pr.stdout)['hookSpecificOutput']
    assert pr.returncode == 2 and out['hookEventName'] == 'PreToolUse' and out['permissionDecision'] == 'deny', (name, pr.returncode, out)
    print(name, 'exit', pr.returncode, 'deny JSON ok; reason:', out['permissionDecisionReason'][:120].replace('\n', ' '))
