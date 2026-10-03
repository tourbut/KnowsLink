#!/bin/sh
# Narrow independent QA for install_bot_mcp.sh at 423db6a. Not the DEV 19-check script.
set -eu

RECORDING=/home/shin/orca/workspaces/KnowsLink/fullops-tester
CANDIDATE=423db6a2a388ea63610462f9d3a5f4c619dd781b
BASELINE=8e46c5a846e6d190e484e48be40b3dc368001a2b
EXPECT_ZIP=d3037d2067c28bf278023a229797eb02111f8d8416bf200d23feff2bf250e609
NODE_SRC=/tmp/claude-1000/-home-shin-orca-workspaces-KnowsLink-fullops-dev/86d7f18f-acff-4924-9ba2-079e434dead3/scratchpad/p1/node
BASE=${1:-/tmp/sar-mvp-002-fix-qa/work}
PROBE=/tmp/sar-mvp-002-fix-qa/probe.mjs

rm -rf "$BASE"
mkdir -p "$BASE/logs" "$BASE/home" "$BASE/outside" "$BASE/shims" /tmp/sar-mvp-002-fix-qa/npm-cache
CLONE=$BASE/clone
LOG=$BASE/logs
QA_HOME=$BASE/home
NPM_CACHE=/tmp/sar-mvp-002-fix-qa/npm-cache
fail=0

check() {
  name=$1
  shift
  if "$@"; then
    echo "ok   $name"
  else
    echo "FAIL $name"
    fail=$((fail + 1))
  fi
}

rc_of() {
  cat "$1.rc"
}

no_ready() {
  ! grep -q 'Ready. Register' "$1.out" && ! grep -q 'Ready. Register' "$1.err"
}

no_match() {
  ! grep -q "$1" "$2"
}

no_stage() {
  ! stage_left "$1"
}

outside_clone() {
  case $1 in
    "$2"|"$2"/*) return 1 ;;
  esac
  return 0
}

stage_left() {
  find "$1" -maxdepth 1 -name '.stage.*' -print | grep -q .
}

run_install() {
  prefix=$1
  log=$2
  path_prefix=${3-}
  set +e
  (
    set -eu
    cd "$CLONE"
    export HOME="$QA_HOME"
    export npm_config_cache="$NPM_CACHE"
    unset GROK_HOME GROK_CONFIG GROK_CONFIG_PATH GROK_MANAGED_CONFIG_URL
    export KNOWSLINK_PREFIX="$prefix"
    if [ -n "$path_prefix" ]; then
      PATH="$path_prefix:$PATH"
      export PATH
    fi
    sh scripts/install_bot_mcp.sh
  ) >"$log.out" 2>"$log.err"
  status=$?
  set -e
  echo "$status" >"$log.rc"
}

hash_file() {
  if [ -f "$1" ]; then
    sha256sum "$1" | awk '{print $1}'
  else
    echo MISSING
  fi
}

write_py_shim() {
  real=$1
  dest=$2
  cat >"$dest" <<EOF
#!/bin/sh
if [ "\$1" = "-m" ] && [ "\$2" = "zipfile" ]; then
  echo "shim: extract failed" >&2
  exit 1
fi
exec $real "\$@"
EOF
  chmod +x "$dest"
}

write_env_shim() {
  real=$1
  dest=$2
  cat >"$dest" <<EOF
#!/bin/sh
if [ "\$1" = "-i" ]; then
  echo "shim: env -i check failed" >&2
  exit 1
fi
exec $real "\$@"
EOF
  chmod +x "$dest"
}

write_mv_fail_shim() {
  dest=$1
  cat >"$dest" <<'EOF'
#!/bin/sh
if [ "$#" -eq 2 ]; then
  case "$1" in
    */.stage.*/knowslink)
      echo "shim: swap mv failed" >&2
      exit 1
      ;;
  esac
fi
exec /bin/mv "$@"
EOF
  chmod +x "$dest"
}

write_mv_term_shim() {
  dest=$1
  cat >"$dest" <<'EOF'
#!/bin/sh
if [ "$#" -eq 2 ]; then
  case "$2" in
    */.stage.*/old-knowslink)
      /bin/mv "$1" "$2" || exit $?
      echo "shim: sent SIGTERM to $PPID after first swap rename" >&2
      kill -TERM "$PPID"
      sleep 3
      exit 0
      ;;
  esac
fi
exec /bin/mv "$@"
EOF
  chmod +x "$dest"
}

echo "== host =="
uname -s | tee "$LOG/uname-s.out"
uname -m | tee "$LOG/uname-m.out"
"$NODE_SRC/bin/node" --version | tee "$LOG/reused-node-version.out"
"$NODE_SRC/bin/npm" --version | tee "$LOG/reused-npm-version.out"
sha256sum "$NODE_SRC/bin/node" >"$LOG/reused-node.sha256"
check "reused node is v22.22.2" grep -qx 'v22.22.2' "$LOG/reused-node-version.out"
check "reused npm is 10.9.7" grep -qx '10.9.7' "$LOG/reused-npm-version.out"

echo "== grok hashes before =="
sha256sum /home/shin/.grok/config.toml /home/shin/.grok/trusted_folders.toml >"$LOG/grok-before.sha256"
find /home/shin/.grok/installed-plugins -type f -print0 | sort -z | xargs -0 sha256sum >"$LOG/plugins-before.sha256"
if [ -e /workspace/.knowslink ]; then echo present; else echo absent; fi >"$LOG/workspace-before.txt"

echo "== reuse diff =="
set +e
git -C "$RECORDING" diff --exit-code "$BASELINE" "$CANDIDATE" -- cmd internal adapters/src adapters/package.json adapters/.npmrc >"$LOG/reuse-core.diff"
echo $? >"$LOG/reuse-core.rc"
git -C "$RECORDING" diff --name-status "$BASELINE" "$CANDIDATE" -- cmd internal adapters >"$LOG/reuse-adapters.name"
git -C "$RECORDING" diff --name-only "$BASELINE" "$CANDIDATE" -- '*.html' '*.css' '*.tsx' '*.jsx' >"$LOG/reuse-ui.name"
echo $? >"$LOG/reuse-ui.rc"
set -e
check "core paths unchanged vs 8e46c5a" test "$(cat "$LOG/reuse-core.rc")" = 0
check "ui name diff empty" test ! -s "$LOG/reuse-ui.name"

echo "== clone =="
git clone --local "$RECORDING" "$CLONE" >"$LOG/clone.out" 2>&1
git -C "$CLONE" checkout --detach "$CANDIDATE" >"$LOG/checkout.out" 2>&1
check "clone head is candidate" test "$(git -C "$CLONE" rev-parse HEAD)" = "$CANDIDATE"

mkdir -p "$BASE/node-runtime"
cp -a "$NODE_SRC/." "$BASE/node-runtime/"
check "copied runtime version" test "$("$BASE/node-runtime/bin/node" --version)" = "v22.22.2"

P1=$BASE/p1
rm -rf "$P1"
mkdir -p "$P1"
cp -a "$BASE/node-runtime" "$P1/node"
echo "== first install, reused node, no download =="
run_install "$P1" "$LOG/install-first"
check "first exit 0" test "$(rc_of "$LOG/install-first")" = 0
check "first did not checksum a new tarball" no_match 'tar.gz: OK' "$LOG/install-first.out"
check "first Ready" grep -q 'Ready. Register' "$LOG/install-first.out"
check "first command path" grep -qx "  Command:     $P1/node/bin/node" "$LOG/install-first.out"
check "first arguments path" grep -qx "  Arguments:   $P1/knowslink/dist/plugin.js" "$LOG/install-first.out"
check "command file is the prepared node" test -x "$P1/node/bin/node"
check "arguments file is the prepared plugin" test -f "$P1/knowslink/dist/plugin.js"
check "prepared node version" test "$("$P1/node/bin/node" --version)" = "v22.22.2"
check "first no stage" no_stage "$P1"
sha256sum "$CLONE/build/knowslink-grok-bot-plugin.zip" >"$LOG/zip.sha256"
echo $? >"$LOG/zip.rc"
check "zip command exit 0" test "$(cat "$LOG/zip.rc")" = 0
check "zip sha256" grep -q "^$EXPECT_ZIP  " "$LOG/zip.sha256"
python3 -m zipfile -l "$CLONE/build/knowslink-grok-bot-plugin.zip" >"$LOG/zip.list"
check "zip has no STALE entry" no_match 'STALE' "$LOG/zip.list"

echo "== outside env -i probe =="
set +e
/usr/bin/node "$PROBE" "$P1/node/bin/node" "$P1/knowslink/dist/plugin.js" "$BASE/outside" "$LOG/probe-first.json" >"$LOG/probe-first.out" 2>"$LOG/probe-first.err"
echo $? >"$LOG/probe-first.rc"
set -e
check "probe first exit 0" test "$(cat "$LOG/probe-first.rc")" = 0
check "probe cwd is outside clone" outside_clone "$BASE/outside" "$CLONE"

GOOD_PLUGIN=$(hash_file "$P1/knowslink/dist/plugin.js")
GOOD_NODE=$(hash_file "$P1/node/bin/node")
echo "$GOOD_PLUGIN" >"$LOG/good-plugin.sha256"
echo "$GOOD_NODE" >"$LOG/good-node.sha256"

echo "== rerun keeps unrelated files and drops stale bundle markers =="
PR=$BASE/rerun
rm -rf "$PR"
mkdir -p "$PR"
cp -a "$P1/node" "$P1/knowslink" "$PR/"
mkdir -p "$PR/package/knowslink"
printf 'old-package\n' >"$PR/package/knowslink/STALE"
printf 'old-bundle\n' >"$PR/knowslink/STALE"
printf 'keep\n' >"$PR/notes.txt"
run_install "$PR" "$LOG/install-rerun"
check "rerun exit 0" test "$(rc_of "$LOG/install-rerun")" = 0
check "rerun Ready" grep -q 'Ready. Register' "$LOG/install-rerun.out"
check "rerun no tarball download" no_match 'tar.gz: OK' "$LOG/install-rerun.out"
check "rerun no STALE in new bundle" test -z "$(find "$PR/knowslink" -name STALE -print)"
check "old package STALE kept" test "$(cat "$PR/package/knowslink/STALE")" = "old-package"
check "notes kept" test "$(cat "$PR/notes.txt")" = "keep"
check "rerun no stage" no_stage "$PR"
check "rerun command still prepared node" grep -qx "  Command:     $PR/node/bin/node" "$LOG/install-rerun.out"
check "rerun arguments still prepared plugin" test -f "$PR/knowslink/dist/plugin.js"
set +e
/usr/bin/node "$PROBE" "$PR/node/bin/node" "$PR/knowslink/dist/plugin.js" "$BASE/outside" "$LOG/probe-rerun.json" >"$LOG/probe-rerun.out" 2>"$LOG/probe-rerun.err"
echo $? >"$LOG/probe-rerun.rc"
set -e
check "probe rerun exit 0" test "$(cat "$LOG/probe-rerun.rc")" = 0

echo "== relative prefix rejected before change =="
rm -rf "$CLONE/rel"
run_install "rel/pfx" "$LOG/relative"
check "relative exit 1" test "$(rc_of "$LOG/relative")" = 1
check "relative signal" grep -q 'must be an absolute path' "$LOG/relative.err"
check "relative dir absent" test ! -e "$CLONE/rel"
check "relative no Ready" no_ready "$LOG/relative"

echo "== node user file refused =="
PN=$BASE/nodefile
rm -rf "$PN"
mkdir -p "$PN"
printf 'mine\n' >"$PN/node"
run_install "$PN" "$LOG/nodefile"
check "nodefile exit 1" test "$(rc_of "$LOG/nodefile")" = 1
check "nodefile kept" test "$(cat "$PN/node")" = "mine"
check "nodefile only that file" test "$(ls -A "$PN")" = "node"
check "nodefile no Ready" no_ready "$LOG/nodefile"
check "nodefile signal" grep -q 'refusing to replace' "$LOG/nodefile.err"

echo "== node symlink refused =="
PL=$BASE/nodelink
rm -rf "$PL"
mkdir -p "$PL"
printf 'mine\n' >"$PL/node-target"
ln -s node-target "$PL/node"
run_install "$PL" "$LOG/nodelink"
check "nodelink exit 1" test "$(rc_of "$LOG/nodelink")" = 1
check "nodelink still a symlink" test -L "$PL/node"
check "nodelink target kept" test "$(cat "$PL/node-target")" = "mine"
check "nodelink no knowslink" test ! -e "$PL/knowslink"
check "nodelink no stage" no_stage "$PL"
check "nodelink no Ready" no_ready "$LOG/nodelink"
check "nodelink signal" grep -q 'refusing to replace' "$LOG/nodelink.err"

echo "== knowslink non-owned directory refused =="
PK=$BASE/kdir
rm -rf "$PK"
mkdir -p "$PK"
cp -a "$P1/node" "$PK/node"
rm -rf "$PK/knowslink"
mkdir -p "$PK/knowslink"
printf 'mine\n' >"$PK/knowslink/x"
node_before=$(hash_file "$PK/node/bin/node")
run_install "$PK" "$LOG/kdir"
check "kdir exit 1" test "$(rc_of "$LOG/kdir")" = 1
check "kdir content kept" test "$(cat "$PK/knowslink/x")" = "mine"
check "kdir plugin absent" test ! -e "$PK/knowslink/dist/plugin.js"
check "kdir node unchanged" test "$(hash_file "$PK/node/bin/node")" = "$node_before"
check "kdir no stage" no_stage "$PK"
check "kdir no Ready" no_ready "$LOG/kdir"
check "kdir signal" grep -q 'refusing to replace' "$LOG/kdir.err"

echo "== bundle extract failure keeps old bundle =="
PE=$BASE/badextract
rm -rf "$PE"
mkdir -p "$PE"
cp -a "$P1/node" "$P1/knowslink" "$PE/"
printf 'old-bundle\n' >"$PE/knowslink/OLD-BUNDLE"
printf 'keep\n' >"$PE/notes.txt"
extract_before=$(hash_file "$PE/knowslink/dist/plugin.js")
mkdir -p "$BASE/shims/extract"
write_py_shim "$(command -v python3)" "$BASE/shims/extract/python3"
run_install "$PE" "$LOG/badextract" "$BASE/shims/extract"
check "badextract exit 1" test "$(rc_of "$LOG/badextract")" = 1
check "badextract signal" grep -q 'shim: extract failed' "$LOG/badextract.err"
check "badextract no Ready" no_ready "$LOG/badextract"
check "badextract old marker kept" test "$(cat "$PE/knowslink/OLD-BUNDLE")" = "old-bundle"
check "badextract plugin hash kept" test "$(hash_file "$PE/knowslink/dist/plugin.js")" = "$extract_before"
check "badextract notes kept" test "$(cat "$PE/notes.txt")" = "keep"
check "badextract no stage" no_stage "$PE"

echo "== env -i verification failure keeps old bundle =="
PV=$BASE/badcheck
rm -rf "$PV"
mkdir -p "$PV"
cp -a "$P1/node" "$P1/knowslink" "$PV/"
printf 'old-bundle\n' >"$PV/knowslink/OLD-BUNDLE"
check_before=$(hash_file "$PV/knowslink/dist/plugin.js")
mkdir -p "$BASE/shims/env"
write_env_shim "$(command -v env)" "$BASE/shims/env/env"
run_install "$PV" "$LOG/badcheck" "$BASE/shims/env"
check "badcheck exit 1" test "$(rc_of "$LOG/badcheck")" = 1
check "badcheck signal" grep -q 'shim: env -i check failed' "$LOG/badcheck.err"
check "badcheck no Ready" no_ready "$LOG/badcheck"
check "badcheck old marker kept" test "$(cat "$PV/knowslink/OLD-BUNDLE")" = "old-bundle"
check "badcheck plugin hash kept" test "$(hash_file "$PV/knowslink/dist/plugin.js")" = "$check_before"
check "badcheck no stage" no_stage "$PV"

echo "== swap second mv failure: preservation is observed, not assumed =="
PS=$BASE/swapfail
rm -rf "$PS"
mkdir -p "$PS"
cp -a "$P1/node" "$P1/knowslink" "$PS/"
printf 'old-bundle\n' >"$PS/knowslink/OLD-BUNDLE"
mkdir -p "$PS/package/knowslink"
printf 'old-package\n' >"$PS/package/knowslink/STALE"
printf 'keep\n' >"$PS/notes.txt"
swap_before=$(hash_file "$PS/knowslink/dist/plugin.js")
swap_node=$(hash_file "$PS/node/bin/node")
mkdir -p "$BASE/shims/mvfail"
write_mv_fail_shim "$BASE/shims/mvfail/mv"
run_install "$PS" "$LOG/swapfail" "$BASE/shims/mvfail"
{
  echo "rc=$(rc_of "$LOG/swapfail")"
  echo "plugin_before=$swap_before"
  echo "plugin_after=$(hash_file "$PS/knowslink/dist/plugin.js")"
  echo "node_before=$swap_node"
  echo "node_after=$(hash_file "$PS/node/bin/node")"
  if [ -d "$PS/knowslink" ]; then echo "knowslink=present"; else echo "knowslink=absent"; fi
  if [ -f "$PS/knowslink/OLD-BUNDLE" ]; then echo "old_marker=present"; else echo "old_marker=absent"; fi
  if stage_left "$PS"; then echo "stage=present"; else echo "stage=absent"; fi
  echo "notes=$(cat "$PS/notes.txt" 2>/dev/null || echo MISSING)"
  echo "package_stale=$(cat "$PS/package/knowslink/STALE" 2>/dev/null || echo MISSING)"
} >"$LOG/swapfail.observe"
check "swapfail exit 1" test "$(rc_of "$LOG/swapfail")" = 1
check "swapfail signal" grep -q 'shim: swap mv failed' "$LOG/swapfail.err"
check "swapfail no Ready" no_ready "$LOG/swapfail"
check "swapfail no stage" no_stage "$PS"
check "swapfail unrelated notes kept" test "$(cat "$PS/notes.txt")" = "keep"
check "swapfail old package kept" test "$(cat "$PS/package/knowslink/STALE")" = "old-package"
check "swapfail node kept" test "$(hash_file "$PS/node/bin/node")" = "$swap_node"
# Document claim: a failure keeps the previous bundle. Record the result either way.
check "swapfail previous bundle kept" test "$(hash_file "$PS/knowslink/dist/plugin.js")" = "$swap_before"
check "swapfail previous marker kept" test "$(cat "$PS/knowslink/OLD-BUNDLE")" = "old-bundle"

echo "== SIGTERM after first swap rename =="
PT=$BASE/swapterm
rm -rf "$PT"
mkdir -p "$PT"
cp -a "$P1/node" "$P1/knowslink" "$PT/"
printf 'old-bundle\n' >"$PT/knowslink/OLD-BUNDLE"
printf 'keep\n' >"$PT/notes.txt"
term_before=$(hash_file "$PT/knowslink/dist/plugin.js")
term_node=$(hash_file "$PT/node/bin/node")
mkdir -p "$BASE/shims/mvterm"
write_mv_term_shim "$BASE/shims/mvterm/mv"
run_install "$PT" "$LOG/swapterm" "$BASE/shims/mvterm"
{
  echo "rc=$(rc_of "$LOG/swapterm")"
  echo "plugin_before=$term_before"
  echo "plugin_after=$(hash_file "$PT/knowslink/dist/plugin.js")"
  echo "node_before=$term_node"
  echo "node_after=$(hash_file "$PT/node/bin/node")"
  if [ -d "$PT/knowslink" ]; then echo "knowslink=present"; else echo "knowslink=absent"; fi
  if [ -f "$PT/knowslink/OLD-BUNDLE" ]; then echo "old_marker=present"; else echo "old_marker=absent"; fi
  if stage_left "$PT"; then echo "stage=present"; else echo "stage=absent"; fi
  echo "notes=$(cat "$PT/notes.txt" 2>/dev/null || echo MISSING)"
} >"$LOG/swapterm.observe"
check "swapterm signal" grep -q 'sent SIGTERM' "$LOG/swapterm.err"
check "swapterm no Ready" no_ready "$LOG/swapterm"
check "swapterm no stage" no_stage "$PT"
check "swapterm notes kept" test "$(cat "$PT/notes.txt")" = "keep"
check "swapterm node kept" test "$(hash_file "$PT/node/bin/node")" = "$term_node"
check "swapterm previous bundle kept" test "$(hash_file "$PT/knowslink/dist/plugin.js")" = "$term_before"

echo "== clone and recording cleanliness =="
git -C "$CLONE" rev-parse HEAD >"$LOG/clone-head.out"
git -C "$CLONE" diff --exit-code >"$LOG/clone-diff.out"
echo $? >"$LOG/clone-diff.rc"
git -C "$CLONE" status --porcelain >"$LOG/clone-status.out"
git -C "$RECORDING" status --porcelain >"$LOG/recording-status.out"
check "clone head unchanged" test "$(cat "$LOG/clone-head.out")" = "$CANDIDATE"
check "clone tracked diff empty" test "$(cat "$LOG/clone-diff.rc")" = 0
check "recording checkout untouched during run" test ! -s "$LOG/recording-status.out"

echo "== grok hashes after =="
sha256sum /home/shin/.grok/config.toml /home/shin/.grok/trusted_folders.toml >"$LOG/grok-after.sha256"
find /home/shin/.grok/installed-plugins -type f -print0 | sort -z | xargs -0 sha256sum >"$LOG/plugins-after.sha256"
if [ -e /workspace/.knowslink ]; then echo present; else echo absent; fi >"$LOG/workspace-after.txt"
check "user grok toml unchanged" cmp -s "$LOG/grok-before.sha256" "$LOG/grok-after.sha256"
check "user installed-plugins unchanged" cmp -s "$LOG/plugins-before.sha256" "$LOG/plugins-after.sha256"
check "workspace prefix still absent" grep -qx absent "$LOG/workspace-after.txt"

echo "qa_fail=$fail"
echo "$fail" >"$LOG/qa_fail.txt"
exit 0
