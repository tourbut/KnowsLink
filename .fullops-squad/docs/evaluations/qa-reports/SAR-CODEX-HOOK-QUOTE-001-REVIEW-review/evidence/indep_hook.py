import json, os, shutil, subprocess, sys
snap = sys.argv[1]
command = json.load(open(os.path.join(snap, '.codex', 'hooks.json'), encoding='utf-8'))['hooks']['PreToolUse'][0]['hooks'][0]['command']
print('command:', command)
launch = {'cmd shell=True': command, 'cmd /C raw': f'cmd.exe /C "{command}"'}
for exe in ('pwsh', 'powershell'):
    p = shutil.which(exe)
    print(exe, '->', p)
    if p:
        launch[exe] = [p, '-NoProfile', '-Command', command]
launch['argv list cmd /C (known non-supported)'] = ['cmd.exe', '/C', command]
DENY = 'r' + 'm -r' + 'f $X/'  # deletion text only goes to guard stdin as JSON; never executed
for cwd in ('', 'scripts', 'adapters'):
    for name, args in launch.items():
        out = []
        for text in ('git status', DENY):
            pr = subprocess.run(args, shell=isinstance(args, str), input=json.dumps({'tool_name': 'Bash', 'tool_input': {'command': text}}),
                                capture_output=True, text=True, cwd=os.path.join(snap, cwd), timeout=60,
                                env={**os.environ, 'PYTHONDONTWRITEBYTECODE': '1'})
            out.append((pr.returncode, pr.stdout.strip()[:200], pr.stderr.strip()[:160]))
        print(f'[{cwd or "root"}] {name}: allow={out[0][0]} deny={out[1][0]} | deny stdout={out[1][1]!r} stderr={out[1][2]!r} | allow stderr={out[0][2]!r}')
# deny JSON path: exit 0 + JSON? check guard stdout for the exit-2 case; report keys
