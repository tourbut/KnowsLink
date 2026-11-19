#!/bin/sh
# Independent QA for install_bot_mcp.sh at 8e46c5a. Scratch only.
set -eu
ROOT=/tmp/knowslink-bot-catalog-qa
SRC=/home/shin/orca/workspaces/KnowsLink/fullops-tester
CANDIDATE=8e46c5a846e6d190e484e48be40b3dc368001a2b
PRIOR_INSTALL=5506d646c47c2d64b35d1254ddfdec5e2003084d
PRIOR_PLUGIN=552586b6e886f95bffa9a000a031ea03070afedb
PRIOR_UI=9584aafcbb5fee88dcc6d618caf660884f6a527d
PRIOR_GO=c4ebbecd3ef91be10ecbb517fe451e02769358bc
EXPECTED_ZIP=b7882df74537ad0bd32bdde45f9dd01677431ff74dda3312ef6c8fa650c00cad
X64_SUM=978978a635eef872fa68beae09f0aad0bbbae6757e444da80b570964a97e62a3
ARM_SUM=b2f3a96f31486bfc365192ad65ced14833ad2a3c2e1bcefec4846902f264fa28
LOG=$ROOT/logs
CLONE=$ROOT/clone
PREFIX=$ROOT/prefix
OUTSIDE=$ROOT/outside-cwd
export npm_config_cache=$ROOT/npm-cache
export npm_config_update_notifier=false
mkdir -p "$LOG" "$OUTSIDE" "$ROOT/bin"
rm -f "$ROOT/bin/curl" "$ROOT/bin/uname"
rm -rf "$PREFIX" "$ROOT/clone-build-fail" "$ROOT/unsupported-prefix" "$ROOT/download-fail-prefix" "$ROOT/checksum-fail-prefix"
: > "$LOG/asserts.txt"
fail=0

note() { printf '%s\n' "$1" >> "$LOG/asserts.txt"; }
ok() { note "PASS $1"; }
bad() { note "FAIL $1"; fail=1; }

record() {
  name=$1
  shift
  set +e
  "$@" >"$LOG/$name.out" 2>"$LOG/$name.err"
  code=$?
  set -e
  printf '%s %s\n' "$code" "$name" >> "$LOG/exits.txt"
  printf '%s\n' "$code"
}

{
  printf 'date %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  uname -s
  uname -m
  command -v node
  node --version
  npm --version
  python3 --version
  command -v xz
  command -v curl
  command -v sha256sum
  git --version
} > "$LOG/versions.out" 2>"$LOG/versions.err"
printf '0 versions\n' >> "$LOG/exits.txt"

python3 - "$HOME" "$LOG/grok-before.txt" <<'PY'
import hashlib, sys
from pathlib import Path
home, dest = Path(sys.argv[1]), Path(sys.argv[2])
grok = home / ".grok"
paths = [grok / "config.toml", grok / "trusted_folders.toml"]
plugins = grok / "installed-plugins"
if plugins.exists():
    paths.extend(sorted(path for path in plugins.rglob("*") if path.is_file() and not path.is_symlink()))
rows = []
for path in paths:
    if not path.is_file():
        rows.append(f"MISSING {path.relative_to(grok)}")
        continue
    digest = hashlib.sha256(path.read_bytes()).hexdigest()
    rows.append(f"{digest} {path.stat().st_size} {path.relative_to(grok)}")
dest.write_text("\n".join(rows) + ("\n" if rows else ""))
PY

if [ -e /workspace/.knowslink ]; then
  bad "workspace prefix absent before"
else
  ok "workspace prefix absent before"
fi

if [ -d "$CLONE/.git" ]; then
  ok "clone already present"
else
  code=$(record git_clone git clone --local "$SRC" "$CLONE")
  [ "$code" -eq 0 ] && ok "git clone exit 0" || bad "git clone exit $code"
fi
code=$(record git_detach git -C "$CLONE" checkout --detach "$CANDIDATE")
[ "$code" -eq 0 ] && ok "detach exit 0" || bad "detach exit $code"
head=$(git -C "$CLONE" rev-parse HEAD)
[ "$head" = "$CANDIDATE" ] && ok "clone HEAD $CANDIDATE" || bad "clone HEAD $head"
if [ -z "$(git -C "$CLONE" status --porcelain)" ]; then
  ok "clone clean after detach"
else
  bad "clone dirty after detach"
fi
cd "$CLONE"

code=$(record sh_n sh -n "$CLONE/scripts/install_bot_mcp.sh")
[ "$code" -eq 0 ] && ok "sh -n exit 0" || bad "sh -n exit $code"

code=$(record shasums curl -fsSL -o "$LOG/SHASUMS256.txt" https://nodejs.org/dist/v22.22.2/SHASUMS256.txt)
[ "$code" -eq 0 ] && ok "SHASUMS256 download" || bad "SHASUMS256 download $code"
x64_line=$(grep 'node-v22.22.2-linux-x64.tar.gz$' "$LOG/SHASUMS256.txt" || true)
arm_line=$(grep 'node-v22.22.2-linux-arm64.tar.gz$' "$LOG/SHASUMS256.txt" || true)
printf '%s\n' "$x64_line" "$arm_line" > "$LOG/shasums-pins.txt"
[ "$x64_line" = "$X64_SUM  node-v22.22.2-linux-x64.tar.gz" ] && ok "x64 pin matches official" || bad "x64 pin $x64_line"
[ "$arm_line" = "$ARM_SUM  node-v22.22.2-linux-arm64.tar.gz" ] && ok "arm64 pin matches official" || bad "arm64 pin $arm_line"

{
  git -C "$SRC" diff --exit-code "$PRIOR_INSTALL" "$CANDIDATE" -- adapters/package.json adapters/.npmrc adapters/src
  printf 'pkg_src %s\n' $?
  git -C "$SRC" diff --exit-code "$PRIOR_PLUGIN" "$CANDIDATE" -- cmd internal adapters/src
  printf 'go_adapter %s\n' $?
  git -C "$SRC" diff --exit-code "$PRIOR_GO" "$CANDIDATE" -- cmd internal
  printf 'go_c4ebbec %s\n' $?
  git -C "$SRC" diff --name-only "$PRIOR_UI" "$CANDIDATE" -- '*.html' '*.css' '*.tsx' '*.jsx' cmd internal
  printf 'ui_names_done %s\n' $?
  git -C "$SRC" diff --name-status "$PRIOR_INSTALL" "$CANDIDATE" -- Makefile adapters/README.md scripts/install_bot_mcp.sh scripts/verify_grok_plugin.py adapters/src cmd internal
} > "$LOG/reuse-diff.out" 2>"$LOG/reuse-diff.err"
printf '0 reuse_diff\n' >> "$LOG/exits.txt"

# xz is first on PATH and fails if any child runs it.
cat > "$ROOT/bin/xz" <<'EOF'
#!/bin/sh
echo xz-called >> /tmp/knowslink-bot-catalog-qa/logs/xz-sentinel.txt
exit 1
EOF
chmod 755 "$ROOT/bin/xz"
: > "$LOG/xz-sentinel.txt"

mkdir -p "$PREFIX/unrelated" "$PREFIX/user-settings"
printf 'unrelated-marker\n' > "$PREFIX/unrelated/keep.txt"
printf 'user-setting-marker\n' > "$PREFIX/user-settings/keep.txt"
marker_before=$(sha256sum "$PREFIX/unrelated/keep.txt" "$PREFIX/user-settings/keep.txt")

code=$(record install_first env PATH="$ROOT/bin:$PATH" KNOWSLINK_PREFIX="$PREFIX" sh "$CLONE/scripts/install_bot_mcp.sh")
if [ "$code" -eq 0 ]; then ok "install first exit 0"; else bad "install first exit $code"; fi
if [ ! -s "$LOG/xz-sentinel.txt" ]; then ok "xz not executed"; else bad "xz executed"; fi
if grep -q 'Ready. Register in the Grok Bot app' "$LOG/install_first.out"; then
  ok "first output has Ready"
else
  bad "first output missing Ready"
fi
if grep -q 'node-v22.22.2-linux-x64.tar.gz: OK' "$LOG/install_first.out"; then
  ok "tarball checksum OK line"
else
  bad "tarball checksum line missing"
fi
ready_count=$(grep -c 'PASS: MCP initialize/discovery' "$LOG/install_first.out" || true)
[ "$ready_count" -ge 2 ] && ok "two MCP PASS lines" || bad "MCP PASS count $ready_count"
if grep -q 'extracted standalone MCP test exit=0' "$LOG/install_first.out"; then
  ok "package verify exit line 0"
else
  bad "package verify line"
fi

python3 - "$LOG/install_first.out" "$LOG/registration-first.txt" <<'PY'
import re, sys
text = open(sys.argv[1], encoding="utf-8").read()
def field(label):
    match = re.search(rf"^  {label}:\s+(.*)$", text, re.M)
    return match.group(1).strip() if match else ""
open(sys.argv[2], "w", encoding="utf-8").write(
    "\n".join([field("Name"), field("Type"), field("Command"), field("Arguments"), field("Environment")]) + "\n"
)
PY
name=$(sed -n '1p' "$LOG/registration-first.txt")
type=$(sed -n '2p' "$LOG/registration-first.txt")
command=$(sed -n '3p' "$LOG/registration-first.txt")
arguments=$(sed -n '4p' "$LOG/registration-first.txt")
environment=$(sed -n '5p' "$LOG/registration-first.txt")
[ "$name" = "knowslink" ] && ok "name knowslink" || bad "name $name"
[ "$type" = "Command" ] && ok "type Command" || bad "type $type"
[ "$environment" = "none (held is the default)" ] && ok "environment none" || bad "environment $environment"
case "$command" in
  /*) ok "command absolute" ;;
  *) bad "command not absolute" ;;
esac
case "$arguments" in
  /*) ok "arguments absolute" ;;
  *) bad "arguments not absolute" ;;
esac
[ "$command" = "$PREFIX/node/bin/node" ] && ok "command is staged node" || bad "command path $command"
[ "$arguments" = "$PREFIX/knowslink/dist/plugin.js" ] && ok "arguments is staged bundle" || bad "arguments path $arguments"
[ -x "$command" ] && ok "staged node executable" || bad "staged node missing"
[ -f "$arguments" ] && ok "staged bundle file" || bad "staged bundle missing"
if [ -x "$command" ]; then
  node_version=$("$command" --version)
  npm_version=$("$PREFIX/node/bin/npm" --version)
else
  node_version=missing
  npm_version=missing
fi
printf '%s\n%s\n' "$node_version" "$npm_version" > "$LOG/staged-versions.txt"
[ "$node_version" = "v22.22.2" ] && ok "staged node v22.22.2" || bad "staged node $node_version"
[ "$npm_version" = "10.9.7" ] && ok "staged npm 10.9.7" || bad "staged npm $npm_version"
code=$(record zip_hash sha256sum "$CLONE/build/knowslink-grok-bot-plugin.zip")
[ "$code" -eq 0 ] && ok "zip sha256sum exit 0" || bad "zip sha256sum exit $code"
zip_hash=$(awk '{print $1}' "$LOG/zip_hash.out")
printf '%s\n' "$zip_hash" > "$LOG/zip-hash.txt"
grep -q "$zip_hash" "$LOG/install_first.out" && ok "script printed same zip hash" || bad "script zip hash mismatch"
[ "$zip_hash" = "$EXPECTED_ZIP" ] && ok "zip hash expected" || bad "zip hash $zip_hash"
marker_after=$(sha256sum "$PREFIX/unrelated/keep.txt" "$PREFIX/user-settings/keep.txt")
[ "$marker_before" = "$marker_after" ] && ok "markers preserved after first" || bad "markers changed after first"
bundle_mtime=$(stat -c %Y "$arguments")

code=$(record tarball_again curl -fsSL -o "$LOG/node-v22.22.2-linux-x64.tar.gz" https://nodejs.org/dist/v22.22.2/node-v22.22.2-linux-x64.tar.gz)
[ "$code" -eq 0 ] && ok "independent tarball download" || bad "independent tarball download $code"
printf '%s  %s\n' "$X64_SUM" "$LOG/node-v22.22.2-linux-x64.tar.gz" > "$LOG/tarball.sha"
code=$(record tarball_check sha256sum -c "$LOG/tarball.sha")
[ "$code" -eq 0 ] && ok "independent tarball checksum" || bad "independent tarball checksum $code"

code=$(record probe node "$ROOT/probe.mjs" "$command" "$arguments" "$OUTSIDE" "$LOG/probe.json")
[ "$code" -eq 0 ] && ok "independent probe exit 0" || bad "independent probe exit $code"

# Rerun must not download. A curl shim fails the run if curl starts.
cat > "$ROOT/bin/curl" <<'EOF'
#!/bin/sh
echo curl-called >> /tmp/knowslink-bot-catalog-qa/logs/curl-sentinel.txt
exit 99
EOF
chmod 755 "$ROOT/bin/curl"
: > "$LOG/curl-sentinel.txt"
code=$(record install_rerun env PATH="$ROOT/bin:$PATH" KNOWSLINK_PREFIX="$PREFIX" sh "$CLONE/scripts/install_bot_mcp.sh")
[ "$code" -eq 0 ] && ok "install rerun exit 0" || bad "install rerun exit $code"
if [ ! -s "$LOG/curl-sentinel.txt" ]; then ok "rerun did not call curl"; else bad "rerun called curl"; fi
if grep -q 'node-v22.22.2-linux-x64.tar.gz: OK' "$LOG/install_rerun.out"; then
  bad "rerun repeated checksum line"
else
  ok "rerun skipped checksum line"
fi
if grep -q 'Ready. Register in the Grok Bot app' "$LOG/install_rerun.out"; then
  ok "rerun Ready"
else
  bad "rerun missing Ready"
fi
python3 - "$LOG/install_rerun.out" "$LOG/registration-rerun.txt" <<'PY'
import re, sys
text = open(sys.argv[1], encoding="utf-8").read()
def field(label):
    match = re.search(rf"^  {label}:\s+(.*)$", text, re.M)
    return match.group(1).strip() if match else ""
open(sys.argv[2], "w", encoding="utf-8").write(
    "\n".join([field("Command"), field("Arguments")]) + "\n"
)
PY
[ "$(sed -n '1p' "$LOG/registration-rerun.txt")" = "$command" ] && ok "rerun command same" || bad "rerun command changed"
[ "$(sed -n '2p' "$LOG/registration-rerun.txt")" = "$arguments" ] && ok "rerun arguments same" || bad "rerun arguments changed"
[ -x "$command" ] && [ -f "$arguments" ] && ok "rerun targets still exist" || bad "rerun targets missing"
marker_rerun=$(sha256sum "$PREFIX/unrelated/keep.txt" "$PREFIX/user-settings/keep.txt")
[ "$marker_before" = "$marker_rerun" ] && ok "markers preserved after rerun" || bad "markers changed after rerun"
bundle_mtime_rerun=$(stat -c %Y "$arguments")
[ "$bundle_mtime_rerun" != "$bundle_mtime" ] && ok "bundle mtime replaced" || bad "bundle mtime unchanged"
printf '%s %s\n' "$bundle_mtime" "$bundle_mtime_rerun" > "$LOG/bundle-mtime.txt"

# Unsupported uname. Prefix must stay absent and Ready must stay absent.
cat > "$ROOT/bin/uname" <<'EOF'
#!/bin/sh
if [ "$1" = "-m" ]; then
  echo arm64
else
  echo Darwin
fi
EOF
chmod 755 "$ROOT/bin/uname"
bad_prefix=$ROOT/unsupported-prefix
rm -rf "$bad_prefix"
code=$(record install_unsupported env PATH="$ROOT/bin:$PATH" KNOWSLINK_PREFIX="$bad_prefix" sh "$CLONE/scripts/install_bot_mcp.sh")
[ "$code" -eq 1 ] && ok "unsupported exit 1" || bad "unsupported exit $code"
if [ -e "$bad_prefix" ]; then bad "unsupported created prefix"; else ok "unsupported created no prefix"; fi
if grep -q 'unsupported platform: Darwin arm64' "$LOG/install_unsupported.err"; then
  ok "unsupported message"
else
  bad "unsupported message missing"
fi
if grep -q 'Ready. Register in the Grok Bot app' "$LOG/install_unsupported.out"; then
  bad "unsupported printed Ready"
else
  ok "unsupported no Ready"
fi

# Download failure.
cat > "$ROOT/bin/uname" <<'EOF'
#!/bin/sh
if [ "$1" = "-m" ]; then echo x86_64; else echo Linux; fi
EOF
cat > "$ROOT/bin/curl" <<'EOF'
#!/bin/sh
echo curl-download-fail >&2
exit 1
EOF
chmod 755 "$ROOT/bin/uname" "$ROOT/bin/curl"
dl_prefix=$ROOT/download-fail-prefix
rm -rf "$dl_prefix"
code=$(record install_download_fail env PATH="$ROOT/bin:$PATH" KNOWSLINK_PREFIX="$dl_prefix" sh "$CLONE/scripts/install_bot_mcp.sh")
[ "$code" -eq 1 ] && ok "download fail exit 1" || bad "download fail exit $code"
if grep -q 'Ready. Register in the Grok Bot app' "$LOG/install_download_fail.out"; then
  bad "download fail printed Ready"
else
  ok "download fail no Ready"
fi
if [ -e "$dl_prefix/node" ]; then bad "download fail created node"; else ok "download fail no node"; fi

# Checksum failure. curl writes the -o file and exits 0.
cat > "$ROOT/bin/curl" <<'EOF'
#!/bin/sh
out=
prev=
for arg in "$@"; do
  if [ "$prev" = "-o" ]; then out=$arg; fi
  prev=$arg
done
printf 'not-a-node-tarball' > "$out"
exit 0
EOF
chmod 755 "$ROOT/bin/curl"
sum_prefix=$ROOT/checksum-fail-prefix
rm -rf "$sum_prefix"
code=$(record install_checksum_fail env PATH="$ROOT/bin:$PATH" KNOWSLINK_PREFIX="$sum_prefix" sh "$CLONE/scripts/install_bot_mcp.sh")
[ "$code" -eq 1 ] && ok "checksum fail exit 1" || bad "checksum fail exit $code"
if grep -q 'Ready. Register in the Grok Bot app' "$LOG/install_checksum_fail.out"; then
  bad "checksum fail printed Ready"
else
  ok "checksum fail no Ready"
fi
if grep -q 'FAILED' "$LOG/install_checksum_fail.out" || grep -q 'FAILED' "$LOG/install_checksum_fail.err"; then
  ok "checksum FAILED visible"
else
  bad "checksum FAILED not visible"
fi

# Build failure in a disposable copy. Node is already staged, so curl is not required.
rm -rf "$ROOT/clone-build-fail"
cp -a "$CLONE" "$ROOT/clone-build-fail"
printf '\nexport const broken =\n' >> "$ROOT/clone-build-fail/adapters/src/mcp.ts"
code=$(record install_build_fail env KNOWSLINK_PREFIX="$PREFIX" sh -c 'cd "$1" && exec sh scripts/install_bot_mcp.sh' sh "$ROOT/clone-build-fail")
[ "$code" -eq 1 ] && ok "build fail exit 1" || bad "build fail exit $code"
if grep -q 'Ready. Register in the Grok Bot app' "$LOG/install_build_fail.out"; then
  bad "build fail printed Ready"
else
  ok "build fail no Ready"
fi

# verify_grok_plugin.py failure shim: show grok output and drop GROK_CONFIG*.
mkdir -p "$ROOT/grok-shim"
cat > "$ROOT/grok-shim/grok" <<'EOF'
#!/bin/sh
env | sort > /tmp/knowslink-bot-catalog-qa/logs/grok-shim-env.txt
echo SHIM_STDOUT_MARKER
echo SHIM_STDERR_MARKER >&2
exit 7
EOF
chmod 755 "$ROOT/grok-shim/grok"
code=$(record verify_shim env \
  PATH="$ROOT/grok-shim:$PATH" \
  GROK_HOME=canary-home \
  GROK_CONFIG=canary-config \
  GROK_CONFIG_PATH=canary-config-path \
  GROK_MANAGED_CONFIG_URL=canary-managed-url \
  python3 "$CLONE/scripts/verify_grok_plugin.py")
[ "$code" -eq 1 ] && ok "verify shim exit 1" || bad "verify shim exit $code"
if grep -q 'SHIM_STDOUT_MARKER' "$LOG/verify_shim.err" && grep -q 'SHIM_STDERR_MARKER' "$LOG/verify_shim.err"; then
  ok "verify shows stdout and stderr"
else
  bad "verify hid grok output"
fi
if grep -q 'exit=7' "$LOG/verify_shim.err"; then ok "verify shows exit 7"; else bad "verify exit text"; fi
if grep -q 'PASS: grok plugin validate/install' "$LOG/verify_shim.out" || grep -q 'PASS: grok plugin validate/install' "$LOG/verify_shim.err"; then
  bad "verify failure printed PASS"
else
  ok "verify failure has no PASS"
fi
if grep -E '^(GROK_HOME|GROK_CONFIG|GROK_CONFIG_PATH|GROK_MANAGED_CONFIG_URL)=' "$LOG/grok-shim-env.txt"; then
  bad "GROK_CONFIG family reached grok"
else
  ok "GROK_CONFIG family removed"
fi
if grep -q 'installed bundle held under this node' "$CLONE/scripts/verify_grok_plugin.py"; then
  ok "new PASS wording present"
else
  bad "new PASS wording absent"
fi
if grep -q 'installed copy held"' "$CLONE/scripts/verify_grok_plugin.py"; then
  bad "old PASS wording remains"
else
  ok "old PASS wording absent"
fi

python3 - "$HOME" "$LOG/grok-before.txt" "$LOG/grok-after.txt" "$LOG/grok-diff.txt" <<'PY'
import hashlib, sys
from pathlib import Path
home, before_path, after_path, diff_path = map(Path, sys.argv[1:])
grok = home / ".grok"
paths = [grok / "config.toml", grok / "trusted_folders.toml"]
plugins = grok / "installed-plugins"
if plugins.exists():
    paths.extend(sorted(path for path in plugins.rglob("*") if path.is_file() and not path.is_symlink()))
rows = []
for path in paths:
    if not path.is_file():
        rows.append(f"MISSING {path.relative_to(grok)}")
        continue
    digest = hashlib.sha256(path.read_bytes()).hexdigest()
    rows.append(f"{digest} {path.stat().st_size} {path.relative_to(grok)}")
after_path.write_text("\n".join(rows) + ("\n" if rows else ""))
before = set(before_path.read_text().splitlines())
changed = sorted(set(rows).symmetric_difference(before))
diff_path.write_text("\n".join(changed) + ("\n" if changed else ""))
PY
if [ -s "$LOG/grok-diff.txt" ]; then
  bad "user grok config or installed-plugins changed"
else
  ok "user grok config and installed-plugins unchanged"
fi
if [ -e "$SRC/build/knowslink-grok-bot-plugin.zip" ]; then
  bad "recording checkout zip present"
else
  ok "recording checkout zip absent"
fi
fixture_mode=$(stat -c %a "$SRC/build/qa-fixture.json")
[ "$fixture_mode" = "600" ] && ok "qa-fixture mode 600" || bad "qa-fixture mode $fixture_mode"
if [ -z "$(git -C "$SRC" status --porcelain)" ]; then
  ok "recording checkout tracked clean"
else
  bad "recording checkout tracked dirty"
fi
if [ -e /workspace/.knowslink ]; then
  bad "workspace prefix created"
else
  ok "workspace prefix still absent"
fi
status=$(git -C "$CLONE" status --porcelain)
printf '%s\n' "$status" > "$LOG/clone-status.txt"
if printf '%s\n' "$status" | grep -q .; then
  # Generated install outputs may be ignored. Tracked changes are a problem.
  tracked=$(git -C "$CLONE" status --porcelain --untracked-files=no)
  printf '%s\n' "$tracked" > "$LOG/clone-tracked.txt"
  if [ -n "$tracked" ]; then bad "clone tracked changes"; else ok "clone tracked clean"; fi
else
  ok "clone porcelain clean"
  : > "$LOG/clone-tracked.txt"
fi
git -C "$CLONE" rev-parse HEAD > "$LOG/clone-head-end.txt"
[ "$(cat "$LOG/clone-head-end.txt")" = "$CANDIDATE" ] && ok "clone HEAD unchanged" || bad "clone HEAD moved"

printf '%s\n' "$fail" > "$LOG/fail-flag.txt"
exit "$fail"
