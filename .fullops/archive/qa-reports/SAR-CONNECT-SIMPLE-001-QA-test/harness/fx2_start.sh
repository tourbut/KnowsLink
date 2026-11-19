#!/bin/sh
L=C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/local/SAR-CONNECT-SIMPLE-001-QA
cd D:/workspace/KnowsLink/.devtools/connect-simple-001-qa/src
export PATH="/c/Program Files/Go/bin:$PATH" TEST_SYNTHETIC_DATABASE=1 TEST_DATABASE_URL="postgres://knowslink:$(cat $L/pg.secret)@127.0.0.1:15531/knowslink2?sslmode=disable" CONNECT_FIXTURE_MANIFEST=$L/fixture2-manifest.json
exec go test -tags "integration connectfixture" -run TestConnectBrowserFixture -count=1 -timeout 0 -v ./internal/relay
