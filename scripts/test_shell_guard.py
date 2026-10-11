#!/usr/bin/env python3
"""Synthetic deny/allow regression for scripts/shell_guard.py with Claude, Grok and Codex PreToolUse payloads.

Nothing is executed except the guard itself and nothing is deleted. The temporary workspace stays on disk
as evidence (repo AGENTS.md forbids cleanup until independent review); its path is printed.
"""
import json
import os
import subprocess
import sys
import tempfile
import unittest

GUARD = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'shell_guard.py')
WS = os.path.realpath(tempfile.mkdtemp(prefix='shell-guard-ws-'))
OUTSIDE = os.path.realpath(tempfile.mkdtemp(prefix='shell-guard-outside-'))
TARGET = os.path.join(WS, 'scratch', 'old')
os.makedirs(TARGET)
os.makedirs(os.path.join(WS, '.git'))
os.makedirs(os.path.join(OUTSIDE, 'inner'))  # reached only through LINK: exercises the link-ancestor check
LINK = os.path.join(WS, 'link')
if os.name == 'nt':
    import _winapi
    _winapi.CreateJunction(OUTSIDE, LINK)
else:
    os.symlink(OUTSIDE, LINK)


def literal(path, extra='-Recurse -Force'):
    return f"Remove-Item -LiteralPath '{path}' {extra}"


DENY = [
    # Synthetic shape of the incident: PowerShell double quotes hand bash a backslash-escaped variable.
    'bash -c "TMP_DIR=\\$(mktemp -d); rm -rf \\$TMP_DIR/"',
    'bash -lc "rm -rf $env:TEMP_X/"',
    'rm -rf /',
    'rm -rf "$TMP_DIR/"',
    'rm -rf ""/',
    'rm -fr $HOME',
    'rm -r -f node_modules',
    'rm -f "$TMP"/*',
    'rm -f ""',
    'Remove-Item -Recurse -Force C:\\',
    'Remove-Item -Recurse $env:TEMP\\x',
    'Get-ChildItem build | Remove-Item -Recurse -Force',
    'ls | xargs rm',
    'env -i rm -rf build',
    'timeout 5 rm -rf build',
    'rm --rec build',
    'Remove-Item @p',
    'rd /s /q %TEMP%\\x',
    'cmd /c rd /s /q C:\\x',
    'wsl rm -rf /mnt/c/x',
    'pwsh -Command "Remove-Item -Recurse x"',
    'ssh host "rm -rf /tmp/clone"',
    'Invoke-Expression "rm -r x"',
    'powershell -EncodedCommand cgBtACAALQByAGYAIAAvAA==',
    'find . -name "*.log" -delete',
    'find . -type d -exec rm -rf {} +',
    'git clean -fdx',
    'robocopy empty C:\\x /MIR',
    'python -c "import shutil; shutil.rmtree(\'x\')"',
    '[System.IO.Directory]::Delete("C:\\x", $true)',
    'rm -rf "unbalanced',
    literal('C:\\'),
    literal(WS),
    literal(os.path.join(WS, 'scratch', '..', '..')),
    literal(os.path.join(WS, '.git')),
    literal(LINK),
    literal(os.path.join(LINK, 'inner')),
    literal(OUTSIDE),
    literal(os.path.join(WS, 'missing')),
    literal(os.path.join(WS, 'scratch', '*')),
    f'Remove-Item -LiteralPath "{TARGET}" -Recurse',
    f"Remove-Item -Path '{TARGET}' -Recurse",
    literal(TARGET, '-Recurse -Exclude x'),
    # F1: switch value after a colon.
    'Remove-Item -r:$true C:\\Users\\x\\Documents',
    'Remove-Item -Recurse:$true -Force C:\\Users\\x',
    # F2: inline interpreter running a shell or fs deletion.
    'python -c "import os; os.system(\'rm -rf x\')"',
    'python3 -c "import subprocess; subprocess.run([\'rm\', \'-rf\', \'x\'])"',
    'node -e "require(\'child_process\').execSync(\'rm -rf x\')"',
    'node -e "require(\'fs\').rmdirSync(\'x\')"',
    'node -e "require(\'fs\').unlinkSync(\'x\')"',
    # F3: git global options before clean.
    'git -C . clean -fdx',
    'git --no-pager clean -fdx',
    'git -c core.x=1 -C sub clean -f',
]

ALLOW = [
    'git status',
    'git log --oneline -5',
    'git rm --cached notes.txt',
    'git clean -n',
    'git -C . clean -n',
    'git --no-pager log -1',
    'git -C . status',
    'python3 -c "print(1)"',
    'node -e "console.log(1)"',
    'Get-ChildItem -r:$true C:\\x',
    'go test -race ./cmd/... ./internal/...',
    'go build -trimpath -o build/relay ./cmd/relay',
    'make lint',
    'npm test --prefix adapters',
    'python3 scripts/check_compose.py',
    'bash deploy/knowslink/test_public_config.sh',
    'sh -n scripts/install_bot_mcp.sh',
    'grep -rn "rm -rf" scripts',
    'Get-ChildItem -LiteralPath C:\\x -Recurse',
    'rm build/relay.tmp',
    'Remove-Item -LiteralPath C:\\ws\\out.txt',
    literal(TARGET),
    literal(TARGET, '-Recurse'),
]


def run(payload, env=None):
    proc = subprocess.run([sys.executable, GUARD], input=json.dumps(payload), capture_output=True, text=True,
                          env={**{k: v for k, v in os.environ.items() if k != 'CLAUDE_PROJECT_DIR'}, **(env or {})})
    return proc.returncode, proc.stdout


PROVIDERS = {
    'claude-bash': (lambda c: {'hook_event_name': 'PreToolUse', 'tool_name': 'Bash', 'tool_input': {'command': c},
                               'cwd': os.path.join(WS, 'scratch')}, {'CLAUDE_PROJECT_DIR': WS}),
    'claude-powershell': (lambda c: {'hook_event_name': 'PreToolUse', 'tool_name': 'PowerShell',
                                     'tool_input': {'command': c}, 'cwd': WS}, {'CLAUDE_PROJECT_DIR': WS}),
    'grok': (lambda c: {'hookEventName': 'pre_tool_use', 'hook_event_name': 'PreToolUse',
                        'toolName': 'run_terminal_command', 'toolInput': {'command': c},
                        'workspaceRoot': WS, 'cwd': WS}, {'CLAUDE_PROJECT_DIR': WS}),
    'codex': (lambda c: {'hook_event_name': 'PreToolUse', 'tool_name': 'Bash', 'tool_input': {'command': c},
                         'cwd': WS, 'session_id': 's', 'turn_id': 't', 'tool_use_id': 'u', 'model': 'm',
                         'permission_mode': 'default', 'transcript_path': None}, None),
}


class ShellGuardTest(unittest.TestCase):
    def test_providers_deny_and_allow(self):
        for name, (build, env) in PROVIDERS.items():
            for command in DENY:
                with self.subTest(provider=name, deny=command):
                    code, out = run(build(command), env)
                    self.assertEqual(code, 2)
                    self.assertEqual(json.loads(out)['hookSpecificOutput']['permissionDecision'], 'deny')
            for command in ALLOW:
                with self.subTest(provider=name, allow=command):
                    self.assertEqual(run(build(command), env), (0, ''))

    def test_payload_edges(self):
        self.assertEqual(run({'tool_name': 'Read', 'tool_input': {'file_path': 'x'}}), (0, ''))
        self.assertEqual(run({'tool_name': 'Bash', 'tool_input': {'command': ['bash', '-lc', 'rm -rf /']}})[0], 2)
        proc = subprocess.run([sys.executable, GUARD], input='not json', capture_output=True, text=True)
        self.assertEqual(proc.returncode, 2)

    @unittest.skipUnless(os.name == 'nt', 'the Codex hook command is the Windows cmd /C form')
    def test_codex_hook_from_subfolder(self):
        # F4: Codex runs the hook with the session cwd, so the command must find the guard from a subfolder.
        repo = os.path.dirname(os.path.dirname(GUARD))
        with open(os.path.join(repo, '.codex', 'hooks.json'), encoding='utf-8') as f:
            command = json.load(f)['hooks']['PreToolUse'][0]['hooks'][0]['command']
        for text, code in (('git status', 0), ('rm -rf $X/', 2)):
            for args in (command, ['cmd.exe', '/C', command]):
                proc = subprocess.run(args, shell=isinstance(args, str), input=json.dumps(
                    {'tool_name': 'Bash', 'tool_input': {'command': text}}), capture_output=True, text=True,
                    cwd=os.path.join(repo, 'scripts'))
                self.assertEqual(proc.returncode, code, (text, proc.stderr))

    def test_nothing_was_deleted(self):
        for path in (TARGET, LINK, OUTSIDE, os.path.join(WS, '.git')):
            self.assertTrue(os.path.exists(path), path)


if __name__ == '__main__':
    print(f'shell guard evidence workspace (kept): {WS}', file=sys.stderr)
    unittest.main()
