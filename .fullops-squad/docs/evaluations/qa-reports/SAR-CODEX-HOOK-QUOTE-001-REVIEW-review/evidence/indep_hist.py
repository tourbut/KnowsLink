import json, os, shutil, subprocess, sys
snap, repo = sys.argv[1], sys.argv[2]
DENY = 'r' + 'm -r' + 'f $X/'  # guard stdin JSON only; never executed
def cmd_at(ref):
    t = subprocess.run(['git', '-C', repo, 'show', f'{ref}:.codex/hooks.json'], capture_output=True, text=True, check=True).stdout
    return json.loads(t)['hooks']['PreToolUse'][0]['hooks'][0]['command']
for ref in ('67e7de9287c6c2ac6bb92beb5b31b523f7e64ec6', '7cae8d7', 'd7d6d744d5b0010a384a7389ec8f64ec7bb93336'):
    try:
        command = cmd_at(ref)
    except Exception as e:
        print(ref, 'unavailable', e); continue
    launch = {'cmd /C raw': f'cmd.exe /C "{command}"', 'pwsh': [shutil.which('pwsh'), '-NoProfile', '-Command', command],
              'powershell': [shutil.which('powershell'), '-NoProfile', '-Command', command], 'argv cmd /C': ['cmd.exe', '/C', command]}
    for name, args in launch.items():
        r = []
        for text in ('git status', DENY):
            pr = subprocess.run(args, shell=isinstance(args, str), input=json.dumps({'tool_name': 'Bash', 'tool_input': {'command': text}}),
                                capture_output=True, text=True, errors='replace', cwd=os.path.join(snap, 'scripts'), timeout=60,
                                env={**os.environ, 'PYTHONDONTWRITEBYTECODE': '1'})
            r.append((pr.returncode, ((pr.stderr or '').strip().splitlines() or [''])[0][:90]))
        print(ref[:8], name, 'allow=%s deny=%s' % (r[0][0], r[1][0]), '|', r[0][1], '|', r[1][1])
