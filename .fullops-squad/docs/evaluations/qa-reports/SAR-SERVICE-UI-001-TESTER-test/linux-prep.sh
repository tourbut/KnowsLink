set -eu
umask 077
work=/tmp/knowslink-serviceui-tester-c7828b2-s1fe0133c-r2
test ! -e "$work"
mkdir "$work"
git clone --no-checkout /tmp/knowslink-serviceui-tester-s1fe0133c.bundle "$work/repo" > "$work/clone.log" 2>&1
git -C "$work/repo" checkout --detach c7828b227b5715f7bdf02bff46bf9cf102a758d9 >> "$work/clone.log" 2>&1
git -C "$work/repo" rev-parse HEAD
git -C "$work/repo" status --porcelain | wc -l
mkdir "$work/tools"
tar -xf /tmp/knowslink-serviceui-tester-tools-s1fe0133c.tar -C "$work/tools"
cd "$work/repo"
node --version; go version; make --version | head -n 1
npm ci --prefix adapters > "$work/npm-ci.log" 2>&1
echo npm-ci=$?
