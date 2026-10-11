set -u
work=/tmp/knowslink-serviceui-tester-c7828b2-s1fe0133c-r2
cd "$work/repo"
go vet -tags 'integration browser' ./internal/relay; echo vet-integration-browser exit=$?
go test -tags 'integration browser' -c -o /dev/null ./internal/relay; echo compile-integration exit=$?
git status --porcelain | wc -l; git rev-parse HEAD
