#!/bin/sh
D=/c/Users/shin/AppData/Local/KnowsLinkDevTools/autosetup-001-qa
cd $D/src
export PATH="/c/Program Files/Go/bin:$PATH"
TEST_SYNTHETIC_DATABASE=1 TEST_DATABASE_URL="postgres://knowslink:$(cat $D/pg.secret)@127.0.0.1:15481/knowslink?sslmode=disable" AUTOSETUP_FIXTURE_ADDR=127.0.0.1:18441 AUTOSETUP_FIXTURE_INSTALL_DIR="$(cygpath -w $D/rel2)" exec go test -tags "integration autosetupqa" -run TestAutosetupQA -count=1 -timeout 0 -v ./internal/relay
