---
title: SAR-MVP-003-BIDIRECTIONAL-TIMEOUT 정정 실행 기록
status: draft
updated: 2026-10-04
owner: dev
tasks: [SAR-MVP-003-BIDIRECTIONAL-TIMEOUT]
summary: 기본 10초 응답 본문 timeout이 GC 뒤 동작하지 않는 원인과 최소 수정 및 실제 10초 회귀를 기록한다
---

# SAR-MVP-003-BIDIRECTIONAL-TIMEOUT — 기본 body timeout 정정

## 기준

착수 HEAD는 `f19c49852b48089ccc3b1539221be32c542f7f93`다. 시작 시 작업 트리는 깨끗했다. 대상 제품 코드는 `cd60e7f87eb5ce137eca887980f232b3f67a18d0`다. `git diff cd60e7f HEAD -- . ':!.fullops-squad'`는 비어 있다. 따라서 착수 HEAD의 제품 코드는 대상과 같다. 근거 QA는 `fa16893870fcaf33e968065e88f2e250a845d09c`의 [독립 QA 보고](../../evaluations/qa-reports/SAR-MVP-003-BIDIRECTIONAL-TESTER.md)다. Node는 v22.22.2다.

fullops-common-0.3.2, FULLOPS.md, project.md, 문서 작성 규칙, ponytail full, diagnosing-bugs 절차를 적용했다. Jev find/documents-find/context의 keep을 확인했다. 충돌 후보는 원본 DEV 기록의 "10초 body 중단 성공" 주장이다. 이 기록은 아래 진단으로 정정한다.

## 재현

모든 재현은 실제 로컬 HTTP socket을 사용한다. fetch stub은 사용하지 않는다. 서버는 200 header와 첫 청크 `{"partial":true}`를 쓴 뒤 본문을 끝내지 않는다. 대상은 기본 6인자 `TestTransport`다. 저장 `timeoutMs`는 10000이다.

1. 수정 전 기본 생성자 3회는 모두 15010–15011 ms에 `hung`이다.
2. QA observer `held.mjs`도 수정 전 build에서 `body timeout rejects` hung, `about 10s` 15000으로 실패한다.
3. 같은 fixture의 timeoutMs 1000·3000·6000·8000은 통과한다. 9000·10000·12000은 hung이다.
4. `--trace-gc`는 약 8245 ms의 idle Mark-Compact를 기록한다. 이후 abort가 발생하지 않는다.
5. 100 ms 간격 강제 GC에서는 timeoutMs 1000·3000도 hung이다.

결론: 실패 원인은 관찰 fixture가 아니라 제품 코드다. 실패 여부는 timeout 전에 major GC가 실행되는지에 따라 달라진다. 이 사실이 QA의 "일부 10002 ms reject, 일부 12000 ms 초과" 관찰과 기존 100 ms 시험 통과를 설명한다.

## 원인

`Adapter.request`는 `signal: AbortSignal.timeout(this.timeoutMs)`를 fetch 인자에 inline으로 전달했다. 두 참조가 GC 대상이 된다.

1. Node의 timeout signal timer는 signal을 강하게 참조하지 않는다. fetch가 응답을 반환하면 코드에 signal 참조가 남지 않는다. GC가 signal을 회수하면 abort가 발생하지 않는다.
2. undici는 abort 시 `Response`를 WeakRef로 찾아 body stream을 오류로 끝낸다. 코드는 reader만 유지하고 `Response`는 유지하지 않는다. 그래서 abort가 발생해도 멈춘 `reader.read()`가 끝나지 않을 수 있다.

1번만 수정한 중간 build에서도 강제 GC 시 hung이 남았다. 그래서 두 원인을 모두 수정했다.

## 수정

`adapters/src/core.ts`만 수정했다.

- `request`는 `AbortController`와 `setTimeout`을 만든다. timer 콜백이 controller를 강하게 참조한다. `finally`의 `clearTimeout`이 성공·실패·HTTP 오류·redirect 오류 모두에서 timer를 제거한다.
- 기존 fetch·본문 읽기 코드는 private `fetchJson`으로 옮겼다. 동작은 바꾸지 않았다. 바뀐 부분은 signal 인자뿐이다.
- 각 `reader.read()`는 signal abort promise와 race한다. abort 시 `TimeoutError`로 reject한다. 기존 `finally`의 `reader.cancel()`이 socket을 닫는다. abort promise에는 빈 catch를 붙인다. 이 catch는 unhandled rejection을 막는다.

redirect `error`, 64 KiB 누적 상한, header, credential 처리는 바꾸지 않았다. 새 의존성은 없다.

## 검증

증거는 [실행 로그 폴더](../logs/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT/)에 있다. `red.txt`와 `green.txt`는 같은 fixture와 인자를 사용한다. 실행한 모든 결과를 기록했다. 성공 결과만 골라내지 않았다.

| 검사 | 수정 전 (`red.txt`) | 수정 후 (`green.txt`) |
|---|---|---|
| 기본 10000 ms, GC 강제 없음, 3회 | hung 15010–15011 ms | TimeoutError 10004–10007 ms |
| 기본 10000 ms, 100 ms 강제 GC | — | TimeoutError 10008 ms |
| 1000·3000 ms, 100 ms 강제 GC | hung | TimeoutError 1008·3007 ms |
| QA `held.mjs` observer | 실패, hung/15000 | 3회 모두 전체 통과, 10001·10003·10003 ms |
| socket·timer·rejection (`cleanup.mjs`) | — | 10015 ms TimeoutError, 서버 socket close 10022 ms, 남은 Timeout 0, unhandled 0, 다음 요청 정상 |
| MCP busy (`mcp-busy.mjs`) | — | 멈춘 receive는 10022 ms에 failed, 동시 호출은 busy, 이후 호출은 busy 아님 |

`trial-boundaries.test.ts`에 회귀 시험을 추가했다. 이 시험은 기본 6인자 `TestTransport`, 실제 10초, 100 ms 간격 강제 GC를 사용한다. 결과는 9900–12000 ms 안의 `TimeoutError`여야 한다. 13초 guard가 hung을 실패로 바꾼다. 수정 전 core.ts에서는 exit 1, 수정 후에는 통과했다. 기존 100 ms 시험도 유지한다.

`verify.txt`에 `npm run check --prefix adapters`, `make test`, `make lint` 결과가 있다. 모두 exit 0이다. `make test`의 adapter 시험은 약 10초 길어진다.

실행하지 않은 검사: `make verify-mvp`(실제 Postgres와 두 MCP 왕복). 이번 변경은 Go relay·SQL·MCP 도구 코드를 바꾸지 않았다. 원본 왕복 증거는 정상 본문 경로의 동작으로 재사용한다. 정상 본문 경로는 `npm test`의 MCP 시험, `cleanup.mjs`의 후속 요청, held observer로 확인했다. 실제 Cloudflare·Grok·private credential·배포는 사용하지 않았다.

## 문서 정정

[원본 DEV 실행 기록](SAR-MVP-003-BIDIRECTIONAL.md)의 body timeout 문장은 보존했다. 그 아래에 이 기록으로 가는 정정 링크를 추가했다. `adapters/README.md`는 "응답 본문 읽기까지 포함해 10초"로 실제 코드와 맞췄다. 리뷰·QA·handover 로그는 수정하지 않았다.

## 후속 담당

새 fixed SHA의 delta 독립 리뷰와 좁은 TESTER 재시험은 coor가 배정한다. 다른 low 5건과 실제 배포는 이번 범위에 포함하지 않았다.
