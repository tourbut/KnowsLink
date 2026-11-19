---
title: SAR-SSE-INSTALL-001-TESTER lite QA
status: draft
updated: 2026-10-11
owner: tester
tasks: [SAR-SSE-INSTALL-001-TESTER]
summary: febd128 Linux 고유 클론의 make lint와 make test 증거
---

# SAR-SSE-INSTALL-001-TESTER lite QA

대상 SHA는 `febd1282df2e84a29fa0a53834569e7b9cc028e2`다. 클론은 `/tmp/knowslink-sse001-tester-ctx70e952fbc5c8`다. 소스는 기존 `/tmp/knowslink-sse001-ctxe4490c0d/final`의 로컬 clone이다. 그 final 트리는 수정하지 않았다. dirty 수는 0이다. `make install` 종료코드는 0이다.

## 명령

1. `make lint`를 그 디렉터리에서 실행한다. 종료코드는 0이다. `go vet`, `go mod verify`, adapter format/lint/typecheck, compose, public ingress, `test_public_config.sh`, `sh -n`, `py_compile`이 포함된다.
2. `make test`를 이어서 실행한다. 종료코드는 0이다. `go test -race ./cmd/... ./internal/...`에서 `internal/relay`는 ok다. `npm test`의 `events.test.js`는 `PASS: SSE initial/push/reconnect/coalesce/serial/fallback/cancel and strict metadata frames`다.

`test_public_config.sh` stderr에 `cannot overwrite existing file`과 `invalid tunnel UUID`가 있다. helper가 stdout 뒤에 stderr를 이어 붙인다. 스크립트는 라이브 설정 덮어쓰기를 거부한 뒤 종료코드 0으로 끝난다. 이 문구를 lint 실패로 세지 않는다.

## 경계

polling 기본과 sse coalesce, Bearer SSE, refresh family 철회, webhook SSRF/DNS/TLS/서명/구독 상한, private folder 거부, cloud app ID 형식은 리뷰 보고서의 코드 판정과 이 테스트 통과로 확인했다. Postgres integration tag와 운영 도메인 SSE는 실행하지 않았다.

## 미실행

실제 Google 로그인, 외부 앱 등록, Claude Desktop MCPB 런타임, Dots 계정 Events callback, 운영 SSE, 동의 화면 시각 검수는 실행하지 않았다. 시각 검수는 designer dispatch `ctx_4283f224b494` 소유다. 이 항목을 PASS로 두지 않는다.

## 재현

같은 서버에서 final 트리를 고유 경로로 clone하고 `git checkout --detach febd1282df2e84a29fa0a53834569e7b9cc028e2` 뒤에 `make install`, `make lint`, `make test`를 실행한다. 호스트 키 거부와 비밀값 출력 금지를 유지한다. 운영 DB, Tunnel, 기존 final 클론은 대상으로 쓰지 않는다.
