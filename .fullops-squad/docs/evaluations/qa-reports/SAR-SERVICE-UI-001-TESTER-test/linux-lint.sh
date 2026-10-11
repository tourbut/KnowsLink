set -u
umask 077
work=/tmp/knowslink-serviceui-tester-c7828b2-s1fe0133c-r2
cd "$work/repo"
python3 "$work/tools/scripts/lint.py" --repo "$work/repo" --from 8539ef5b6d3e7ed99259762c62d164b8ba40fab0 --out "$work/lint.json" > "$work/lint.out" 2>&1
echo lint.py exit=$?
tail -40 "$work/lint.out"
test -z "$(git status --porcelain)" && echo clean-after
go test ./internal/relay -run 'TestMemberUIFixtures|TestMemberPagesShowNextSteps' -count=1 -v 2>&1 | grep -E "^(--- |ok|FAIL|PASS)"
git rev-parse HEAD
