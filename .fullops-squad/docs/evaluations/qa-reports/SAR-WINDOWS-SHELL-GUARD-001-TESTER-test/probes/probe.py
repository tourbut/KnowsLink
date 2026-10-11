import json, os, subprocess, sys, tempfile
G = sys.argv[1]
SP = os.path.dirname(os.path.abspath(__file__))
WS = os.path.realpath(tempfile.mkdtemp(prefix='tester-ws-', dir=SP))
OUT = os.path.realpath(tempfile.mkdtemp(prefix='tester-out-', dir=SP))
os.makedirs(os.path.join(WS, 'a', 'b'))
os.makedirs(os.path.join(OUT, 'inner'))
open(os.path.join(WS, 'a', 'f.txt'), 'w').close()
LINK = os.path.join(WS, 'lnk')
import _winapi
_winapi.CreateJunction(OUT, LINK)


def run(cmd, env_root=WS, cwd=None, raw=None):
    p = raw if raw is not None else {'hook_event_name': 'PreToolUse', 'tool_name': 'Bash',
                                     'tool_input': {'command': cmd}, 'cwd': cwd or WS}
    env = {k: v for k, v in os.environ.items() if k != 'CLAUDE_PROJECT_DIR'}
    if env_root:
        env['CLAUDE_PROJECT_DIR'] = env_root
    data = p if isinstance(p, str) else json.dumps(p)
    try:
        r = subprocess.run([sys.executable, G], input=data, capture_output=True, text=True, env=env, timeout=15)
    except subprocess.TimeoutExpired:
        return 'TIMEOUT'
    return r.returncode


def L(p, x='-Recurse -Force'):
    return f"Remove-Item -LiteralPath '{p}' {x}"


A = os.path.join(WS, 'a', 'b')
cases = [
    ('incident ps->bash var', 'bash -c "rm -rf \\$TMP/"', 'D'),
    ('incident full path git bash', '& "C:\\Program Files\\Git\\bin\\bash.exe" -c "rm -rf \\$X/"', 'D'),
    ('wsl dash e sh', 'wsl.exe -e sh -c "rm -rf $X/"', 'D'),
    ('call var program', '& $sh -c "rm -rf $X/"', 'D'),
    ('python os.system', "python -c \"import os; os.system('rm -rf /')\"", 'D'),
    ('python subprocess', "python3 -c \"import subprocess; subprocess.run(['rm','-rf','/'])\"", 'D'),
    ('node rmdirSync', "node -e \"require('fs').rmdirSync('x',{recursive:true})\"", 'D'),
    ('node unlinkSync', "node -e \"require('fs').unlinkSync('x')\"", 'D'),
    ('docker mount rm', 'docker run -v C:\\:/h alpine rm -rf /h/x', 'D'),
    ('busybox rm', 'busybox rm -rf /', 'D'),
    ('git -C clean', 'git -C . clean -fdx', 'D'),
    ('git --no-pager clean', 'git --no-pager clean -fdx', 'D'),
    ('git clean -fdn dry', 'git clean -fdn', 'A'),
    ('robocopy /MOVE', 'robocopy a b /MOVE', 'D'),
    ('robocopy /mov', 'robocopy a b /E /MOV', 'D'),
    ('backslash obfusc', 'r\\m -rf /', 'D'),
    ('empty-quote obfusc', 'r""m -rf /', 'D'),
    ('alias backslash', '\\rm -rf /', 'D'),
    ('Remove-Item -r drive', 'Remove-Item -r C:\\', 'D'),
    ('ri -r', 'ri -Recurse C:\\x', 'D'),
    ('del /s', 'del /s /q C:\\x', 'D'),
    ('rd /s via ps', 'rd /s /q C:\\x', 'D'),
    ('Remove-Item nonrec outside file', 'Remove-Item C:\\Users\\x\\f.docx', 'A?'),
    ('rm nonrec outside', 'rm /c/Users/x/f.docx', 'A?'),
    ('rm glob', 'rm -f /c/Users/x/*', 'D'),
    ('Move-Item drive', 'Move-Item -LiteralPath C:\\Users -Destination D:\\x', 'A?'),
    ('mv', 'mv /c/Users /tmp/x', 'A?'),
    ('scriptblock create', '& ([scriptblock]::Create("rm -r x"))', 'D'),
    ('ForEach remove-item var', 'gci | % { Remove-Item $_ -Recurse }', 'D'),
    ('Remove-Item (subexpr)', 'Remove-Item (Join-Path $env:X "a") -Recurse', 'D'),
    ('Remove-Item -Path in ws', f"Remove-Item -Path '{A}' -Recurse", 'D'),
    ('Remove-Item -LiteralPath $p', 'Remove-Item -LiteralPath $p -Recurse', 'D'),
    ('sanctioned ok', L(A), 'A'),
    ('sanctioned file ok', L(os.path.join(WS, 'a', 'f.txt'), '-Force'), 'A'),
    ('sanct -WhatIf', L(A, '-Recurse -WhatIf'), 'D?'),
    ('sanct chain ; outside', L(A) + "; Remove-Item -LiteralPath 'C:\\' -Recurse", 'D'),
    ('sanct && outside', L(A) + ' && ' + L(OUT), 'D'),
    ('sanct 2 paths', f"Remove-Item -LiteralPath '{A}' '{OUT}' -Recurse", 'D'),
    ('sanct comma paths', f"Remove-Item -LiteralPath '{A}','{OUT}' -Recurse", 'D'),
    ('sanct via junction', L(os.path.join(LINK, 'inner')), 'D'),
    ('sanct junction itself', L(LINK), 'D'),
    ('sanct outside', L(OUT), 'D'),
    ('sanct trailing sep', L(A + '\\'), '?'),
    ('sanct fwd slash mix', L(A.replace('\\', '/')), '?'),
    ('sanct upper-case path', L(A.upper()), '?'),
    ('sanct long-path prefix', L('\\\\?\\' + A), '?'),
    ('sanct ws parent', L(os.path.dirname(WS)), 'D'),
    ('sanct sibling prefix', L(WS + '-evil'), 'D'),
    ('sanct trailing dot', L(A + '.'), '?'),
    ('sanct ADS', L(os.path.join(WS, 'a', 'f.txt') + ':x'), 'D'),
    ('sanct tilde short', L(os.path.join(WS, 'A~1')), 'D'),
    ('commit msg false pos', 'git commit -m "remove-item via bash"', 'A?'),
    ('bash script file', 'bash cleanup.sh', 'A?'),
    ('make clean', 'make clean', 'A?'),
    ('python script', 'python cleanup.py', 'A?'),
    ('format-volume', 'Format-Volume -DriveLetter D', 'D?'),
    ('Clear-Disk', 'Clear-Disk -Number 1 -RemoveData', 'D?'),
    ('cipher /w', 'cipher /w:C:\\', 'D?'),
    ('cmd del dir', 'cmd /c del C:\\Users\\x', 'D'),
    ('pwsh nodelete', 'pwsh -Command "Get-Date"', 'A'),
    ('pwsh -File script', 'pwsh -File cleanup.ps1', 'A?'),
    ('env assignment', 'FOO=1 rm -rf $X', 'D'),
    ('rm with --', 'rm -rf -- "$X"', 'D'),
    ('rm -Rf', 'rm -Rf x', 'D'),
    ('rm -vrf', 'rm -vrf x', 'D'),
    ('rm -d drive', 'rm -d /c', 'D'),
    ('rm -f /c/', 'rm -f /c/', 'D'),
    ('rm -f ~', 'rm -f ~', 'D'),
    ('no-preserve-root', 'rm --no-preserve-root -r /', 'D'),
    ('subshell paren', '(rm -rf $X)', 'D'),
    ('brace group', '{ rm -rf $X; }', 'D'),
    ('background &', 'sleep 1 & rm -rf $X', 'D'),
    ('newline', 'echo a\nrm -rf $X', 'D'),
    ('CRLF', 'echo a\r\nrm -rf $X', 'D'),
    ('tab sep', 'rm\t-rf\t$X', 'D'),
    ('uppercase RM', 'RM -RF $X', 'D'),
    ('REMOVE-ITEM upper', 'REMOVE-ITEM -RECURSE C:\\x', 'D'),
    ('xargs rm', 'cat list | xargs rm -f', 'D'),
    ('xargs -I rm', 'cat list | xargs -I{} rm {}', '?'),
    ('xargs -0 rm', 'find . -print0 | xargs -0 rm', 'D'),
    ('find -exec sh', 'find . -exec sh -c "rm $0" {} \\;', 'D'),
    ('find delete', 'find . -delete', 'D'),
    ('find -exec rm', 'find . -exec rm {} +', 'D'),
    ('find -ok', 'find . -ok rm {} \\;', 'D'),
    ('rsync delete', 'rsync -a --delete empty/ target/', 'D'),
    ('git clean -x -d -f', 'git clean -x -d -f', 'D'),
    ('git rm -rf .', 'git rm -rf .', '?'),
    ('git reset hard', 'git reset --hard', '?'),
    ('unlink', 'unlink x', 'A'),
    ('shred', 'shred -u x', 'A?'),
    ('Remove-Item Env', 'Remove-Item Env:FOO', 'A'),
    ('RI -Recurse "ws"', f'Remove-Item -Recurse "{os.path.join(WS, "a")}"', 'D'),
    ('RI path -Recurse', f'Remove-Item {os.path.join(WS, "a")} -Recurse', 'D'),
    ('Recurse:$true outside', "Remove-Item -LiteralPath 'C:\\x' -Recurse:$true", 'D'),
    ('-LP alias', f"Remove-Item -LP '{A}' -Recurse", 'D'),
    ('-Literal abbreviated', f"Remove-Item -Literal '{A}' -Recurse", 'D'),
    ('-rec -fo abbreviated', f"Remove-Item -LiteralPath '{A}' -rec -fo", '?'),
    ('-Confirm:$false', L(A, '-Recurse -Force -Confirm:$false'), 'D?'),
    ('call op Remove-Item', '& Remove-Item -Recurse C:\\x', 'D'),
    ('call op quoted', '& "Remove-Item" -Recurse C:\\x', 'D'),
    ('dot source', '. rm -rf x', 'D'),
    ('Invoke-Command', 'Invoke-Command { Remove-Item C:\\x -Recurse }', 'D'),
    ('Start-Job', 'Start-Job { Remove-Item -Recurse C:\\x }', 'D'),
    ('pipe Remove-Item', 'Get-Item C:\\x | Remove-Item -Recurse', 'D'),
    ('Delete($true)', '(Get-Item C:\\x).Delete($true)', 'D'),
    ('Delete()', '(Get-Item C:\\x).Delete()', 'A?'),
    ('Directory.Delete', '[IO.Directory]::Delete("C:\\x",$true)', 'D'),
    ('Directory.Delete 1', "[System.IO.Directory]::Delete('C:\\x', 1)", 'D'),
    ('VB DeleteDirectory', '[Microsoft.VisualBasic.FileIO.FileSystem]::DeleteDirectory("C:\\x","DeleteAllContents")', 'D?'),
    ('FSO DeleteFolder', '(New-Object -ComObject Scripting.FileSystemObject).DeleteFolder("C:\\x")', 'D?'),
    ('cmd rmdir quoted', 'cmd.exe /c "rmdir /s /q C:\\x"', 'D'),
    ('powershell.exe nested', 'powershell.exe -NoProfile -Command "rm -r $env:X"', 'D'),
    ('ssh', 'ssh host rm -rf /', 'D'),
    ('eval', 'eval "rm -rf $X"', 'D'),
    ('node child_process', "node -e \"require('child_process').execSync('rm -rf x')\"", 'D'),
    ('perl unlink', 'perl -e "unlink glob q(*)"', 'D?'),
    ('ruby FileUtils', 'ruby -e "FileUtils.rm_rf(%q(x))"', 'D?'),
    ('python rmtree', "python -c \"import shutil; shutil.rmtree('x')\"", 'D'),
    ('python pathlib rmdir', "python -c \"import pathlib; pathlib.Path('x').rmdir()\"", 'D'),
    ('python removedirs', "python -c \"import os; os.removedirs('x')\"", 'D?'),
    ('python exec b64', "python3 -c \"exec(__import__('base64').b64decode('cm0gLXJm'))\"", 'D?'),
    ('go run', 'go run x.go', 'A?'),
]
for label, cmd, exp in cases:
    c = run(cmd)
    got = 'D' if c == 2 else ('A' if c == 0 else f'E{c}')
    flag = ''
    if exp == 'D' and got != 'D':
        flag = 'BYPASS'
    elif exp == 'A' and got != 'A':
        flag = 'FALSEPOS'
    print(f'{flag:9}{got:3}exp={exp:3} {label}')

edges = {
    'no command key': run(None, raw={'tool_name': 'Bash', 'tool_input': {}}),
    'command null': run(None, raw={'tool_name': 'Bash', 'tool_input': {'command': None}}),
    'command int': run(None, raw={'tool_name': 'Bash', 'tool_input': {'command': 5}}),
    'command dict': run(None, raw={'tool_name': 'Bash', 'tool_input': {'command': {'a': 'rm -rf /'}}}),
    'tool_input str': run(None, raw={'tool_name': 'Bash', 'tool_input': 'rm -rf /'}),
    'payload list': run(None, raw=['rm -rf /']),
    'payload null': run(None, raw='null'),
    'payload empty': run(None, raw=''),
    'toolInput cmd': run(None, raw={'toolInput': {'command': 'rm -rf /'}}),
    'alt key cmd': run(None, raw={'tool_input': {'cmd': 'rm -rf /'}}),
    'alt key script': run(None, raw={'tool_input': {'script': 'rm -rf /'}}),
    'list cmd': run(None, raw={'tool_input': {'command': ['bash', '-lc', 'rm -rf /']}}),
    'empty cmd': run('  '),
    'huge cmd': run('echo ' + 'a' * 2_000_000),
    'spaces redos': run('bash ' + ' ' * 50000 + 'x'),
    'nul byte': run('rm\x00 -rf /'),
    'fullwidth rm': run('ｒｍ -rf /'),
    'sanct no env no cwd': run(None, env_root=None, raw={'tool_input': {'command': L(A)}}),
    'sanct no env cwd=subdir (root shrinks)': run(None, env_root=None, raw={'tool_input': {'command': L(os.path.join(WS, 'a'))}, 'cwd': A}),
}
for k, v in edges.items():
    print('EDGE', v, k)
print('exists after:', all(os.path.exists(p) for p in (WS, OUT, A, os.path.join(WS, 'a', 'f.txt'), LINK)))
