#!/usr/bin/env python3
"""PreToolUse shell guard: block cross-shell and unverified recursive deletion before it runs.

Claude Code and Grok load it from .claude/settings.json, Codex from .codex/hooks.json. Each sends the
shell command as tool_input.command (Grok: toolInput.command). A blocked call prints a PreToolUse deny
and exits 2; an allowed call prints nothing and exits 0.

The only allowed recursive deletion is one standalone PowerShell command:
    Remove-Item -LiteralPath '<absolute path>' -Recurse [-Force]
for an existing path inside the workspace that is not the workspace, a drive root, .git, or a link.

This is a pre-execution text check, not a complete security boundary: script files, aliases,
generated code and deliberate obfuscation bypass it. Claude Code and Grok let the call run when
this script cannot start (python3 unresolved, timeout); Codex blocks on exit 2 and treats any other failure as
a non-blocking error. The hook command must keep resolving python3 and this file.
"""
import json
import os
import re
import stat
import sys

DELETE_WORD = re.compile(
    r'(?<![\w.-])(?:rm|rmdir|rd|del|erase|ri|remove-item|unlink|shred|rimraf|robocopy|rmtree|rmsync|removeall)'
    r'(?![\w-])|git\s+clean|\]::delete\(|-delete(?![\w-])', re.I)
LAUNCHER = re.compile(
    r'(?<![\w.-])(?:[\w.:~-]*[/\\])?(?:bash|sh|zsh|dash|ksh|wsl|cmd|pwsh|powershell|ssh)(?:\.exe)?(?=[\s"\']|$)'
    r'|(?<![\w-])(?:invoke-expression|iex|start-process|saps|eval|invoke-command|icm)(?=[\s"\'(]|$)', re.I)
ENCODED = re.compile(r'(?<![\w-])-(?:e|ec|enc|encodedcommand)\s+[A-Za-z0-9+/=]{16,}', re.I)
INTERPRETER = re.compile(r'(?<![\w.-])(?:python3?|py|node|deno|bun|ruby|perl)(?:\.exe)?(?=\s)', re.I)
PROGRAM_DELETE = re.compile(r'rmtree|rmsync|rmdirsync|unlinksync|rimraf|removeall|os\.remove|os\.unlink|\.unlink\(|rmdir\(|fs\.rm\b', re.I)
DOTNET_DELETE = re.compile(r'\]::delete\(|\.delete\(\s*\$true', re.I)
VERBS = {'rm', 'rmdir', 'rd', 'del', 'erase', 'ri', 'remove-item', 'unlink', 'shred', 'rimraf'}
PREFIXES = {'sudo', 'doas', 'command', 'builtin', 'exec', 'nohup', 'time', 'env', 'xargs', 'nice', 'timeout',
            'stdbuf', 'ionice', 'chrt', 'unbuffer', 'strace', 'watch'}
UNVERIFIED = re.compile(r'[$%`*?\[{~@]')
SAFE_HINT = ("Safe path: one PowerShell command `Remove-Item -LiteralPath '<absolute path inside the workspace>' "
             "-Recurse -Force` after checking the path; repo AGENTS.md forbids delete/move cleanup until the "
             "guard's independent review.")


def tokenize(command):
    """Split into simple commands of (value, single_quoted) tokens; None when quotes are unbalanced."""
    commands, words, word, quoted, mode, i = [], [], None, False, None, 0

    def end_word():
        nonlocal word, quoted
        if word is not None:
            words.append((word, quoted))
        word, quoted = None, False

    while i < len(command):
        c = command[i]
        if mode == "'":
            if c == "'":
                mode = None
            else:
                word += c
        elif mode == '"':
            if c in '\\`' and i + 1 < len(command) and command[i + 1] in '"$`\\':
                word += command[i:i + 2]
                i += 1
            elif c == '"':
                mode = None
            else:
                word += c
        elif c in '\'"':
            quoted = c == "'" and word is None  # stays true only for a word that is one '...' segment
            word, mode = (word or ''), c
        elif c.isspace() and c != '\n':
            end_word()
        elif c in ';|&\n(){}':
            end_word()
            if words:
                commands.append(words)
            words = []
        else:
            word = (word or '') + c
            quoted = False
        i += 1
    if mode:
        return None
    end_word()
    if words:
        commands.append(words)
    return commands


def verb_of(words):
    """Skip env assignments, wrapper commands and their options/numbers; return (verb, index of the verb word)."""
    for i, (value, _) in enumerate(words):
        name = os.path.basename(value.replace('\\', '/')).lower().removesuffix('.exe')
        if re.fullmatch(r'\w+=.*', value) or name in PREFIXES or (i and (value.startswith('-') or re.fullmatch(r'[\d.]+[smhd]?', value))):
            continue
        return name, i
    return '', len(words)


def recursive(verb, args):
    for value, _ in args:
        arg = value.lower().split(':')[0]  # PowerShell switch value: -r:$true
        if arg in ('-r', '/s') or arg.startswith(('-rec', '--r')):
            return True
        if verb in ('rm', 'rmdir', 'rd', 'del', 'erase', 'ri', 'remove-item', 'shred') and re.fullmatch(r'-[a-z]{1,4}', arg) and 'r' in arg:
            return True
    return False


def targets(args):
    out, skip = [], False
    for value, single in args:
        if skip:
            out.append((value, single))
            skip = False
        elif value.lower() in ('-path', '-literalpath', '-lp', '-pspath'):
            skip = True
        elif not value.startswith('-') and not re.fullmatch(r'/[A-Za-z]', value):
            out.append((value, single))
    return out


def git_subcommand(args):
    """First git argument that is not a global option or an option's value."""
    skip = False
    for value, _ in args:
        if skip:
            skip = False
        elif value in ('-C', '-c', '--git-dir', '--work-tree', '--namespace', '--exec-path'):
            skip = True
        elif not value.startswith('-'):
            return value
    return None


def same(a, b):
    return os.path.normcase(os.path.normpath(a)) == os.path.normcase(os.path.normpath(b))


def sanctioned(words, root):
    """Return None when words are the allowed LiteralPath form for a verified path, else the reason."""
    values = [v.lower() for v, _ in words]
    if not values or values[0] != 'remove-item' or '-literalpath' not in values:
        return 'recursive deletion is only allowed as Remove-Item -LiteralPath with a single-quoted absolute path'
    at = values.index('-literalpath') + 1
    rest = set(values[1:at - 1] + values[at + 1:])
    if at >= len(words) or not words[at][1] or not rest <= {'-recurse', '-force'}:
        return 'Remove-Item must take one single-quoted -LiteralPath plus only -Recurse/-Force'
    path = words[at][0]
    parts = re.split(r'[/\\]+', path)
    if not os.path.isabs(path) or UNVERIFIED.search(path) or '..' in parts or '.' in parts[1:]:
        return 'path must be absolute, without variables, wildcards or ./.. segments'
    if os.path.dirname(os.path.normpath(path)) == os.path.normpath(path) or re.fullmatch(r'[A-Za-z]:[/\\]?', path):
        return 'drive or filesystem root'
    if '.git' in (p.lower() for p in parts):
        return '.git is never a cleanup target'
    if not root:
        return 'workspace root unknown'
    workspace = os.path.realpath(root)
    try:
        info = os.lstat(path)
    except OSError:
        return 'target does not exist'
    if stat.S_ISLNK(info.st_mode) or getattr(info, 'st_file_attributes', 0) & 0x400:
        return 'target is a symlink, junction or other reparse point'
    real = os.path.realpath(path)
    if not same(real, path):
        return 'path passes through a link; real path differs'
    if same(real, workspace) or os.path.normcase(os.path.commonpath([workspace, real])) != os.path.normcase(workspace):
        return 'path must be inside the workspace and not the workspace itself'
    return None


def check(command, root):
    """Return a deny reason for the command, or None to allow."""
    text = command.strip()
    if not text:
        return None
    if ENCODED.search(text) and LAUNCHER.search(text):
        return 'encoded shell command cannot be verified'
    if LAUNCHER.search(text) and DELETE_WORD.search(text):
        return ('cross-shell or indirect deletion: deletion text is passed to another shell or evaluator, '
                'where the outer shell can expand variables to empty paths')
    if DOTNET_DELETE.search(text) or (INTERPRETER.search(text) and (PROGRAM_DELETE.search(text) or DELETE_WORD.search(text))):
        return 'programmatic deletion inside an inline command cannot be verified'
    commands = tokenize(text)
    if commands is None:
        return 'unbalanced quotes around a deletion command' if DELETE_WORD.search(text) else None
    for words in commands:
        verb, at = verb_of(words)
        args = words[at + 1:]
        if verb == 'find' and any(v.lower() in ('-delete', '-exec', '-execdir', '-ok') for v, _ in args) and DELETE_WORD.search(' '.join(v for v, _ in args)):
            return 'find with -delete/-exec deletion is an unverified recursive deletion'
        if verb == 'git' and git_subcommand(args) == 'clean' and not any(v in ('-n', '--dry-run') for v, _ in args):
            return 'git clean deletes untracked files and evidence; use git clean -n to inspect'
        if verb == 'robocopy' and any(v.lower() in ('/mir', '/purge') for v, _ in args):
            return 'robocopy /MIR or /PURGE deletes destination files'
        if verb not in VERBS:
            continue
        if verb == 'rimraf' or recursive(verb, args):
            reason = sanctioned(words[at:], root)
            if reason:
                return f'unverified recursive deletion: {reason}'
            continue
        found = targets(args)
        if not found:
            return 'deletion without explicit targets (pipeline or xargs input) cannot be verified'
        for value, _ in found:
            if not value or UNVERIFIED.search(value) or re.fullmatch(r'[/\\]+|[A-Za-z]:[/\\]?|/[a-z]|/mnt/[a-z]/?', value):
                return f'deletion target {value!r} is empty, a root, or contains a variable or wildcard'
    return None


def main():
    try:
        payload = json.load(sys.stdin)
        tool_input = payload.get('tool_input') or payload.get('toolInput') or {}
        command = tool_input.get('command') if isinstance(tool_input, dict) else None
        if isinstance(command, list):
            command = ' '.join(map(str, command))
        root = os.environ.get('CLAUDE_PROJECT_DIR') or payload.get('workspaceRoot') or payload.get('cwd')
        reason = check(command, root) if isinstance(command, str) else None
    except Exception as error:  # fail closed: a shell call the guard cannot read is not run
        reason = f'shell guard could not read the hook payload ({type(error).__name__})'
    if not reason:
        return 0
    message = f'Blocked by scripts/shell_guard.py: {reason}. {SAFE_HINT}'
    print(json.dumps({'hookSpecificOutput': {'hookEventName': 'PreToolUse', 'permissionDecision': 'deny',
                                             'permissionDecisionReason': message}}))
    print(message, file=sys.stderr)
    return 2


if __name__ == '__main__':
    sys.exit(main())
