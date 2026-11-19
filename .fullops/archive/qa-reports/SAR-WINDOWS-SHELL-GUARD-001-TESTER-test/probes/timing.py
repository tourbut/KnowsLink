import json, subprocess, sys, time, os
G = sys.argv[1]
env = {k: v for k, v in os.environ.items()}
env['CLAUDE_PROJECT_DIR'] = os.path.dirname(os.path.abspath(__file__))
def t(cmd):
    s = time.time()
    try:
        r = subprocess.run([sys.executable, G], input=json.dumps({'tool_input': {'command': cmd}}), capture_output=True, text=True, env=env, timeout=60)
        return r.returncode, round(time.time() - s, 2)
    except subprocess.TimeoutExpired:
        return 'TIMEOUT>60', 60
for n in (1000, 10000, 50000, 100000, 200000, 400000):
    print('echo a*n', n, t('echo ' + 'a' * n))
for n in (10000, 50000, 100000):
    print('multi-line echo', n, t('\n'.join(['echo hello world'] * (n // 16))))
    print('git commit heredoc-like words', n, t('git commit -m "' + 'fix the thing and update docs ' * (n // 30) + '"'))
