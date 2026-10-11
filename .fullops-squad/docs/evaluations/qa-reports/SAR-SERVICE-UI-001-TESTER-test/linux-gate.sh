set -u
umask 077
work=/tmp/knowslink-serviceui-tester-c7828b2-s1fe0133c-r2
cd "$work/repo"
test -z "$(git status --porcelain)" && echo clean-before
for t in "product-lint make lint" "product-test make test"; do
  set -- $t; label=$1; shift
  "$@" > "$work/$label.log" 2>&1; echo "$label exit=$?" | tee -a "$work/exits.txt"
done
test -z "$(git status --porcelain)" && echo clean-after
tail -5 "$work/product-test.log"
git rev-parse HEAD
