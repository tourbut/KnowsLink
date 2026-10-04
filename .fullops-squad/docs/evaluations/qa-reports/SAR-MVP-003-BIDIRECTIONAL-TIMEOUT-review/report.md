---
title: SAR-MVP-003-BIDIRECTIONAL-TIMEOUT 리뷰
status: review
updated: 2026-10-04
owner: ops
tasks: [SAR-MVP-003-BIDIRECTIONAL-TIMEOUT-REVIEW]
summary: 기본 10초 body timeout 정정 delta를 독립 검토하고 원본 QA medium 해소를 판정한다
---

# SAR-MVP-003-BIDIRECTIONAL-TIMEOUT 리뷰

- 검토자 / CLI / 모델: ops 독립 검토 세션 `c4c411b8-d66d-4356-981b-9e2279440b02` / Claude Code / `claude-opus-5-5` high. 원본 리뷰와 같은 reviewer다. 최신 구현자 DEV 세션은 `91bb4020-55f1-4789-85f4-b53d2b3ce838`이다. 이 ID는 fullops-dev 세션 기록에서 확인했다. 두 세션은 서로 다르다.
- base SHA / head SHA / merge-base: `cd60e7f87eb5ce137eca887980f232b3f67a18d0` / `711f2532be423d1ca7707463a20fdc168f50bece` / `cd60e7f87eb5ce137eca887980f232b3f67a18d0`.
- snapshot: `/tmp/knowslink-timeout-review-711f253`. detached·clean·읽기 전용으로 유지했다. 실행은 별도 scratch clone에서 했다.
- OCR 버전 / 적용 규칙: open-code-review delegate, `review/rule.json`. 공통 규칙 `fullops-common-0.3.2`와 `project.md`는 원본 리뷰와 같은 버전이다. 예외 없음.
- 요구사항·완료 기준 원천: `handovers/to_ops.md`, [TIMEOUT 실행 기록](../../exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-TIMEOUT.md), [원본 QA](../SAR-MVP-003-BIDIRECTIONAL-TESTER.md) medium 판정. Jev TIMEOUT find/documents-find/context keep을 확인했다. 충돌은 지시서에 적힌 대로 독립 재검증 대상으로 처리했다.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 75 / 47 / 28 / 75 / 0.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: 0 / 1 / 0. WARNING은 `.fullops-squad/PLANS.md` 603→613줄 SIZE-001이다. coor 진행 기록이며 제품 코드가 아니다.

## 검토 범위

제품 delta는 `adapters/src/core.ts`(33줄), `adapters/src/trial-boundaries.test.ts`, `adapters/README.md` 세 파일이다. `git diff cd60e7f 711f253 -- . ':!.fullops-squad'`로 확인했다. Go relay·SQL·MCP 도구·CLI·launcher는 `cd60e7f`와 같다. 따라서 [원본 리뷰](../SAR-MVP-003-BIDIRECTIONAL-review/report.md)(기록 SHA `0c36730106025be95709236b63f080a510a4aebe`)를 재사용한다. 원본 리뷰 폴더는 `0c36730`과 diff 0이다. 원본 QA 폴더와 보고서는 `fa16893870fcaf33e968065e88f2e250a845d09c`와 diff 0이다. 전체 Go/SQL/Chrome/CLI 시험은 반복하지 않았다.

나머지 72개는 `.fullops-squad` 기록이다. 각 reason을 result.json에 적었다. 제외 28개(실행 로그 `.txt/.log/.err/.out/.exit`)도 직접 열람했다. 비밀값 패턴은 없었다.

### core.ts delta 정밀 검토

| 항목 | 확인 | 판정 |
|---|---|---|
| timer | `setTimeout` 콜백이 `AbortController`를 강하게 참조한다. GC가 timeout signal을 회수해 abort가 사라지는 원인을 제거한다. reason은 `TimeoutError` DOMException이다. | 적절 |
| cleanup | `request`의 `finally`가 성공·HTTP 오류·redirect 오류·timeout 모두에서 `clearTimeout`한다. cleanup 재실행에서 남은 Timeout 0이다. | 적절 |
| abort promise | listener 등록 전 `signal.aborted`를 먼저 검사한다. `{once:true}`로 등록한다. 빈 `catch`가 unhandled rejection을 막는다. 재실행 unhandled 0이다. | 적절 |
| race된 read | 각 `reader.read()`를 abort promise와 race한다. undici가 WeakRef로 `Response`를 찾지 못해도 abort가 read를 끝낸다. 남은 read promise는 race가 handler를 붙이므로 unhandled가 되지 않는다. | 적절 |
| reader cancel | 기존 `finally`의 `await reader.cancel()`이 socket을 닫는다. 재실행에서 timeout 10016 ms 뒤 서버 socket이 10024 ms에 닫혔다. 다음 요청도 정상이다. | 적절 |
| busy | MCP 도구의 `busy`는 `finally`에서 해제된다. request가 이제 약 10초에 끝나므로 busy가 고착되지 않는다. DEV `mcp-busy.mjs` 결과(첫 호출 10022 ms failed, 동시 호출 busy, 이후 비 busy)를 코드와 대조했다. | 적절 |
| 비밀값·기존 경계 | redirect `error`, 64 KiB 누적 상한, header·credential 처리는 바뀌지 않았다. 오류 메시지에 credential을 넣지 않는다. | 유지 |

`reader.cancel()` 자체는 timeout으로 제한되지 않는다. undici는 cancel 시 연결을 즉시 abort하며 재실행에서도 8 ms 안에 끝났다. 결함으로 보지 않는다. 원본 low 5번(비정상 HTTP 응답 본문 미해제)은 이 delta로 바뀌지 않았다.

### 시험과 문서

새 회귀 시험은 기본 6인자 `TestTransport`, 실제 10초, 100 ms 간격 강제 GC를 사용한다. 9900–12000 ms 안의 `TimeoutError`를 요구하고 13초 guard로 hung을 실패로 바꾼다. 리뷰어가 `cd60e7f`의 core.ts에 이 시험을 적용하자 exit 1이었다. 새 core에서는 통과했다. 시험이 결함을 실제로 검출한다. README 문구 "응답 본문 읽기까지 포함해 10초"는 새 동작과 일치한다. 원본 DEV 기록은 원문을 보존하고 정정 인용만 추가했다.

## 발견 사항

이 delta의 신규 발견 사항은 없다(critical/high/medium/low 0). 원본 리뷰의 low 5건은 미해결로 그대로 남는다.

원본 리뷰 정정: 원본 리뷰(`0c36730`)는 100 ms 시험 재실행을 근거로 10초 body timeout을 통과로 판정했다. 이 판단은 GC 의존 결함을 놓쳤다. QA `fa16893`의 medium 판정이 맞았다. 원본 리뷰 기록은 보존하고 이 보고서로 정정한다.

## 원본 QA medium 해소 판정

해소됐다. 근거는 다음과 같다.

1. 원본 QA observer `held.mjs`를 `711f253` scratch build에 1회 실행했다. `body timeout rejects`는 `TimeoutError`, `body timeout is about 10s`는 10003 ms로 통과했다. 원본 QA는 hung/15001이었다.
2. 새 회귀 시험은 실제 10초와 강제 GC에서 통과한다. 이전 core에서는 실패한다.
3. cleanup 시험은 10016 ms `TimeoutError`, socket close, 남은 timer 0, unhandled 0, 후속 요청 정상이다.

## 검증 및 남은 제약

리뷰어가 scratch clone(HEAD `711f253`, 깨끗한 작업 트리)에서 실행했다. 로그는 [evidence](evidence/)에 있다.

| 명령 | 결과 |
|---|---|
| `npm ci --prefix adapters` | exit 0 |
| `lint.py --repo <scratch> --from cd60e7f… --out lint.json` | exit 0, ERROR 0 / WARNING 1 / 실행 불가 0 |
| `make test` | exit 0, Go race unit, MCP, trial-boundaries의 100 ms와 기본 10초 GC 시험 |
| 원본 QA `observe.mjs held` | exit 0, 10003 ms `TimeoutError` |
| DEV `cleanup.mjs` | exit 0, 10016 ms, socket 10024 ms, timer 0, unhandled 0 |
| 이전 core + 새 시험 | exit 1(기대한 red) |

실행하지 않은 검사: `make verify-mvp`, Chrome, CLI 왕복. 해당 코드가 바뀌지 않아 원본 리뷰·QA 증거를 재사용한다. 반복 QA는 tester 담당이다. 실제 Grok 왕복, 실제 Access 인증 호출, 운영 배포는 미검증이다.

## 검토 결론

고정 head `711f2532be423d1ca7707463a20fdc168f50bece`는 수락 가능하다. 원본 QA medium이 해소됐고 미해결 critical/high가 없다. target lint ERROR는 0이다. 수락 범위는 원본 리뷰와 같다. 승인된 시험 text와 D12 13장 절차에 한정하며 실제 Grok 성공은 주장하지 않는다. check 통과는 기록 검사이며 AI 검토와 시험 결과를 자동 보증하지 않는다.
