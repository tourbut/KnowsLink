#!/bin/sh
# Restart the QA fixture from the current fixture source (state is reset on start).
L=C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/local/SAR-CONNECT-SIMPLE-001-QA
S=D:/workspace/KnowsLink/.devtools/connect-simple-001-qa/src
cp C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-QA-test/connect_simple_qa_fixture_test.go $S/internal/relay/
cd $S
export PATH="/c/Program Files/Go/bin:$PATH" TEST_SYNTHETIC_DATABASE=1 TEST_DATABASE_URL="postgres://knowslink:$(cat $L/pg.secret)@127.0.0.1:15531/knowslink?sslmode=disable" CSQA_ADDR=127.0.0.1:18951
exec go test -tags "integration connectsimpleqa" -run TestConnectSimpleQAFixture -count=1 -timeout 0 -v ./internal/relay
