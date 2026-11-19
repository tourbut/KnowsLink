#!/bin/sh
Q=/d/workspace/KnowsLink/.devtools/autosetup-001-qa3
cd $Q/src-new
export PATH="/c/Program Files/Go/bin:$PATH"
TEST_SYNTHETIC_DATABASE=1 TEST_DATABASE_URL="postgres://knowslink:$(cat $Q/pg.secret)@127.0.0.1:15483/knowslink?sslmode=disable" AUTOSETUP_FIXTURE_ADDR=127.0.0.1:18443 AUTOSETUP_FIXTURE_INSTALL_DIR="D:/workspace/KnowsLink/.devtools/autosetup-001-qa3/rel" exec go test -tags "integration autosetupqa" -run TestAutosetupQA -count=1 -timeout 0 -v ./internal/relay
