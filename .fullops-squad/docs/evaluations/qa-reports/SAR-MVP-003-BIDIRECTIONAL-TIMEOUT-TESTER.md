---
title: SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER — 기본 10초 본문 중단 재시험
status: draft
updated: 2026-10-04
owner: tester
tasks: [SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER]
summary: 711f253의 멈춘 부분 응답이 10500ms 안에 TimeoutError로 끝나는지 기록한다
---

# SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER — 기본 10초 본문 중단 재시험

## 판정

판정 후보는 `711f2532be423d1ca7707463a20fdc168f50bece`다.
기록 브랜치 `fullops/tester`의 착수 HEAD는 `a9b033d9460c78b06b60699ec9373e89127f829d`다.
실행 위치는 `/tmp/sar-mvp-003-timeout-qa`의 detached clone이다. clone HEAD는 시작과 끝이 같다. 추적 파일 diff는 없다.
기본 6인자 `TestTransport`의 멈춘 부분 응답 3회는 모두 `TimeoutError`다. 경과 시간은 10007ms, 10008ms, 10007ms다. 세 값은 10500ms보다 작다.
100ms 간격 강제 GC 1회도 `TimeoutError`다. 경과 시간은 10007ms다.
이전 후보 `cd60e7f87eb5ce137eca887980f232b3f67a18d0`에서 기록된 15001ms hung은 이 조건에서 재현되지 않았다. 제품 코드는 수정하지 않았다.
실제 Grok 계정, Cloudflare, 공개 Access, 운영 배포는 미검증이다.

## 기준

날짜는 2026-10-04다. 공통 기준은 `fullops-common-0.3.2`다. 정본은 `.fullops-squad/project.md`다.
적용 문서는 FULLOPS, project, 공통 README, coding-style, security, testing, 문서 작성 규칙, `fullops-test`다.
Jev keep을 확인했다. 충돌 후보는 원본 QA `fa16893870fcaf33e968065e88f2e250a845d09c`의 10초 실패와 DEV 기록의 수정 주장이다. 이 보고서는 그 주장을 clone에서 다시 실행해 확인한다.
Node는 v22.22.2다. `npm run build --prefix adapters`의 종료코드는 0이다.
이번 표의 수치는 그 clone에서 실행한 probe stdout이다.

## 불변 범위와 재사용

`cd60e7f87eb5ce137eca887980f232b3f67a18d0`와 판정 후보의 제품 diff는 세 파일이다.
파일은 `adapters/README.md`, `adapters/src/core.ts`, `adapters/src/trial-boundaries.test.ts`다.
`core.ts` numstat는 31삽입, 2삭제다. README numstat는 1삽입, 1삭제다. `trial-boundaries.test.ts` numstat는 38삽입, 1삭제다.
`git diff --exit-code`의 종료코드는 0이다. 대상은 `cmd`, `internal`, `db`, `scripts`, `Makefile`, `adapters/package.json`, `adapters/package-lock.json`, `adapters/src/mcp.ts`, `adapters/src/test-transport.ts`다.
SQL, Go race, 전체 CLI, 전체 QA는 다시 실행하지 않았다.
왕복 성공 증거는 원본 QA의 실행이다. 이번 SHA에서 그 ID를 다시 만들지 않았다.
MCP 송신과 Grok 수신 ID는 `01a104fe-2f22-7a8c-9ed8-8863e1d962b4`다. text는 `Codex independent ping`이다. from은 `trial_codex`다. to는 `trial_grok`다.
MCP 회신과 Codex 수신 ID는 `01a104fe-2fde-771f-9959-e91f07a3ae67`다.
Codex CLI stdin 송신 ID는 `01a104fe-3073-74f3-9c20-1d6799450e25`다. text는 `Codex cli ping`이다.
Grok launcher 회신 ID는 `01a104fe-31b1-712c-95ee-a4f661328655`다.
이번 clone에서 완료된 JSON 본문 1건은 `{"ok":true}`다. 이 결과는 위 왕복 ID의 재실행이 아니다.

## 재시험

대상은 기본 6인자 `TestTransport`다. 저장 `timeoutMs`는 10000이다. 생성자 length는 6이다.
서버는 `127.0.0.1`의 실제 HTTP다. 응답은 200 header와 첫 청크 `{"partial":true}`를 쓰고 본문을 끝내지 않는다.
10500ms 안에 settle하지 않으면 판정은 hung이다. hung는 수락으로 바꾸지 않는다.

| 실행 | 종료코드 | 결과 | ms |
|---|---|---|---|
| default-1 | 0 | TimeoutError | 10007 |
| default-2 | 0 | TimeoutError | 10008 |
| default-3 | 0 | TimeoutError | 10007 |
| gc 100ms 1회 | 0 | TimeoutError | 10007 |

강제 GC는 `node --expose-gc`와 100ms `setInterval`이다. fetch stub은 쓰지 않았다. 각 probe의 stderr 길이는 0이다.

## 관련 좁은 검사

100ms 명시 timeout의 멈춘 stream은 `TimeoutError`다. 경과 시간은 104ms다. 저장 `timeoutMs`는 100이다.
302 Location header는 `TypeError`다. 경과 시간은 13ms다. `short-100ms` 종료코드는 0이다.
cleanup의 기본 timeout은 `TimeoutError`이고 10017ms다. socket close 시각은 10023ms다. 남은 `Timeout` 자원은 0이다. unhandled rejection은 0이다. 다음 요청 본문은 `{"ok":true}`다.
그 시점의 다른 활성 자원은 `TCPServerWrap`과 `TCPSocketWrap`이다. close 이벤트는 10023ms에 발생했다. `cleanup` 종료코드는 0이다.
MCP `knowslink_test_receive`의 첫 호출은 `failed`이고 10026ms다. 겹친 호출은 `busy`다. 이후 호출 state는 `failed`다. HTTP hit는 2다. stderr에 private key는 없다. `mcp-busy` 종료코드는 0이다.
이 10026ms는 MCP 도구 호출의 경과 시간이다. 10500ms 직접 상한은 위 `TestTransport` 3회, GC 1회, cleanup 1회에 적용했다.
held 모드의 canary hit는 0이다. status, test_send, test_receive, pull_once의 state는 held다. `actualConnection`은 held다. webhook과 evidenceFetch는 false다. `held-network` 종료코드는 0이다.

## 실행 명령

1. clone에서 `npm run build --prefix adapters`를 실행한다. 종료코드는 0이다.
2. `python3 .fullops-squad/docs/evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER-test/run_narrow.py /tmp/sar-mvp-003-timeout-qa`를 실행한다. 종료코드는 0이다.
3. `logs/summary.json`의 `ok`는 true다. clone HEAD는 `711f2532be423d1ca7707463a20fdc168f50bece`다. Node는 v22.22.2다. porcelain은 비어 있다.

## 하지 않은 검증

원본 `held.mjs` 전체, SQL, `make test`, `make verify-mvp`, `make verify-grok-plugin`, 공개 왕복은 실행하지 않았다.
실제 Grok 계정, Cloudflare token, 공개 Access, 운영 배포, 사용자 credential은 사용하지 않았다.
원본 QA 보고서와 원본 리뷰는 수정하지 않았다.
증거 디렉터리는 [SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER-test/](SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-TESTER-test/)다.
