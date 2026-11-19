---
title: 기존 계약 보존 리팩터링의 독립 코드 검토
status: review
updated: 2026-10-14
owner: coor
tasks: [SAR-REFACTOR-002-DEV, SAR-REFACTOR-002-CODE-REVIEW]
summary: 세 제품 수정과 연결 계약을 검토했으며 최종 실행 게이트와 문서 인계 수락은 대기한다.
---

# SAR-REFACTOR-002-CODE-REVIEW

## 대상과 독립성

- base/merge-base: `9255852fa0e46cd40b7be1021eaa1f6c7f7d083b`.
- head: `bcb1dd6629b508a30ec39a5d3141c8ce0a558597`.
- 검토자: Codex / GPT-6 Astra high. 실제 `CODEX_THREAD_ID`는 `01a139ab-4840-7d71-b94b-e033b4b85efb`다. 환경 변수에서 직접 확인했다.
- 작성 부모 `01a1399a-8b4d-78b0-8318-61a5539394a6`, 자식 `01a1399d-3c8c-71d2-81c8-7e133abba2ce`, `01a1399d-cf0a-7f22-8d83-f0401f2c195c`와 모두 다르다.
- snapshot: `C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-a795871780674536a63697a5b542213c`. 설치된 FullOps 1.4.2 review.py snapshot으로 생성했다. detached HEAD와 clean 상태를 확인했다. 제품 파일·테스트·설정을 수정하지 않았다.
- 기록 checkout은 coor다. 결과와 해시는 이 정본 QA 폴더에만 쓴다. snapshot은 보존한다.
- OCR v1.13.0 delegate가 파일·규칙을 준비했다. 의미 판단은 검토자가 수행했으며 OCR 자동 판정으로 표시하지 않는다.

## 기준과 커버리지

고정 후보의 AGENTS.md, FULLOPS.md, rules/common/README.md와 coding-style/testing/security, project.md, docs/agents/document-writing.md를 적용했다. 공통 규칙은 fullops-common-0.3.3이다. 설치된 FullOps 1.4.2 fullops-review와 open-code-review-delegate를 읽었다. test_level은 lite다. 인박스의 현재 권한·서명·수신/ACK·TTL·재시도·writer guard 보존 기준을 적용했다. 프로젝트 정본과 안전 지시는 외부 일반 규칙보다 우선한다.

전체 25개 `(path,status)`, OCR 대상 16개, 제외 9개다. 직접 내용 검토 22개, 생성 탐색 자료의 의미 전수 검토 생략 3개다. 처리율은 100%, 직접 내용 검토율은 88%다. result.json에 모든 항목의 이유를 기록했다.

생략 자료는 최신 packet, history packet, packet-review-index다. JSON 파싱, 과제·attempt·기준 SHA, partial/unknown 상태를 확인했다. history는 base의 이전 packet과 JSON 의미가 같다. 최신 packet 173항목과 index 173줄을 대조했다. 모든 변경 파일의 SHA256·Git blob·크기는 evidence-manifest.json에 보존했다. 원시 benchmark·RED/GREEN·unit·guard 기록은 제외 규칙과 무관하게 직접 읽었다. 생성 추천의 모든 과거 문서를 다시 감사한 것은 아니다.

연결 문맥은 store.go의 transaction/current/sweep, human.go의 authority와 관계 변경, text_events.go, remote_auth.go, remote_mcp.go, remote_events.go, 기존 remote 테스트, adapters/src/inbox.ts·text.ts·mcp.ts와 신규 테스트다. docs/system-design.md·architecture.md·development.md·operations.md 및 D03/D05/D07/D10 관련 절을 대조했다. 라이브러리 API를 새로 설계하지 않는 내부 코드 리뷰이므로 외부 문서 조회는 추가하지 않았다.

## 코드 판단

### 보관 이력 sweep

`now.Before(Exp)`일 때만 current를 실행한다. `now == Exp`는 기존 만료 분기로 간다. 그 분기는 envelope·공유 Inbox·lease/claim token을 비우고 미완료 claimed 메시지 completion을 expired로 만든다. delivered receipt와 24시간 ConversationCopy·멱등 metadata는 보존한다. 기존 코드도 만료 분기가 revoked 결과를 덮었으므로 최종 상태는 같다. current 자체는 권한 상태를 변경하지 않는다.

미만료 철회, Gate 별도 검사, Remote 설치 Inbox의 current/expiry, 실제 receive·persist·ACK·답장·멱등 재전송의 권한 검사는 유지된다. 웹 사본의 참여자 열람 정책을 넓히지 않는다. 새 회귀는 queued/leased/delivered/revoked, revoked true/false, 정확한 Exp와 1ns 전, 24h 경계를 검사한다.

### 로컬 Inbox

pending 조회 뒤 prune을 완료하고 rename claim을 시작한다. housekeeping 오류가 본문을 비운 뒤 take를 실패시키던 경로를 제거한다. stat ENOENT만 다른 receiver의 정상 경합으로 허용하고 다른 오류는 전파한다. 테스트는 실제 임시 Inbox의 특정 marker stat에 EACCES를 주입한다. unread 본문 보존, ENOENT 후 수신, 빈 상태, 읽은 ID 중복 억제를 확인한다. 기존 구현은 unread 보존 assertion에 실패한다.

MCP/CLI 실패 전파와 공개 반환형은 같다. store sync/rename 후 transport persist/ACK 순서를 유지한다. local claim 뒤 process crash·payload read/write 실패·표시 전 crash는 기존 한계다. 이번 수정이 전체 crash 창을 해결했다고 해석하지 않는다.

### Remote callback

pendingByAgent는 한 claimRemote 호출의 지역 map이다. 각 subscription의 설치와 만료를 계속 검사한다. 함수는 메시지/인가 상태를 바꾸지 않으므로 같은 agent의 pending 재사용은 동일 transaction 안에서 안전하다. 다음 claim과 실제 deliverRemote 직전에는 현재 권한을 다시 검사한다.

활성 lease는 새 도착에도 유지된다. 마지막 lease가 끝나 attempts가 3이면 이전 Pending을 Seen으로 옮긴다. 현재 pending이 다를 때만 Next=now다. remoteListen의 deadline이 즉시 끝나 다음 transaction이 새 event ID와 시도 예산을 만든다. 같은 event의 3회 상한은 유지된다. finishRemote 성공·429/5xx·terminal 응답·lease/EventID 비교와 응답 유실 경로를 대조했다.

신규 회귀는 t+0/16/32 claim, t+40 새 도착, t+47 lease 만료를 구성한다. 기존 Next=구독만료 실패를 잡고 즉시 재평가를 검사한다. 별도 회귀는 빈 pending과 이전 권한을 다음 claim에 재사용하지 않는지 확인한다. 기존 remote 테스트는 같은 event 재시도·backoff·3회·410/413·철회를 검사한다.

## 성능 근거

Windows before/after는 같은 fixture와 컴파일러의 짧은 합성 측정이다. 10000 retained 메시지 sweep 중앙값은 38.100ms에서 0.291ms, transactionCPU 중앙값은 120.844ms에서 69.122ms다. CPU harness는 JSON decode, eventVersion, sweep, encode를 포함한다. SQL/lock/네트워크·전체 HTTP 처리량은 포함하지 않는다. 전체 JSONB와 singleton lock 한계는 남는다.

Remote 10000기록/4구독 중앙값은 590438ns에서 162572ns다. 1구독은 148226ns에서 156203ns로 증가했다. 보고서는 반대 결과를 숨기지 않는다. delivered fixture의 map 순회 측정이며 queued 인가·여러 설치·DB/TLS 부하를 대표하지 않는다.

Postgres 벤치는 명시된 synthetic DB의 실제 Service.transaction과 retained count를 검사한다. baseline에 동일 테스트만 주입하는 helper를 읽었다. 후속 전달된 원문과 linux-medians.json을 직접 대조했다. 10000개 실제 DB transaction 중앙값은 588.663216ms에서 542.726660ms로 약 7.8% 감소했다. sweep 중앙값은 57.910169ms에서 0.915579ms다. 같은 synthetic DB와 제한(CPU 1, 메모리 512MiB)에서 baseline→candidate 순서, 각 5회 연산을 3회 반복했다. 짧은 직렬 측정이며 통계적 유의성이나 전체 요청 처리량 개선은 증명하지 않는다. 운영 부하·최대 본문·동시 포화·장시간을 측정했다고 주장하지 않는다.

## 발견 사항과 검사 상태

새로 확정한 제품 critical/high/medium/low 결함은 없다. 전체 저장소 무결함 보증은 아니다.

- 제품 `git diff --check BASE HEAD -- adapters internal`: 원 process exit 0.
- 전체 `git diff --check BASE HEAD`: 원 process exit 2. raw benchmark CPU 행과 생성 packet-index의 trailing whitespace다. patch-check.txt에 원문을 보존했다. 제품 오류로 보지 않으며 원시 실행 기록을 편집해 성공으로 만들지 않는다.
- 하위 focused 실행·RED/GREEN은 해당 시점의 미커밋 작성자 근거다. 현재 후보 독립 실행으로 바꾸지 않는다.
- 현재 Linux 원문 직접 확인: bcb1dd6 build exit0, 새 합성 DB integration exit0/111.338s. helper의 exact HEAD·clean 검사와 build 산출물 선행 조건을 확인했다. 최초 TestPublicNodeProcesses의 connect.js 미생성 실패(exit1/113.075s)는 별도 보존했다. verified-*와 external-execution-manifest.json에 전달 원문과 출처/해시를 보존했다. 저장소 밖 초기 probe는 Go1.22.2였다. 추가 module cwd 원문에서 Go1.27.1 linux/amd64·Node22.22.2·npm10.9.7과 exit0을 확인했다. 작성자 실행이며 reviewer 독립 재실행으로 표시하지 않는다.
- installed lint.py는 내부 pass/lock 파일 unlink를 수행한다. snapshot 읽기 전용과 worker 정리 금지에 따라 reviewer가 실행하지 않았다. 로컬 make도 없다. root와 합의해 DEV의 exact bcb1dd6 installed FullOps 결과를 받았다. base/head·merge-base config SHA256·등록 명령과 실제 Linux lint/test 원문을 대조했다. product-lint/test는 실제 exit0, unavailable0이다. FullOps 자체는 ANTI-004 ERROR1이며 lint.json에 원문 그대로 보존했다.
- review.py check를 고정 base/head와 실제 task-key로 실행했다. lint.json 미도착으로 원 exit1이다. check-pending.txt를 보존했다. 기록 검사를 통과로 표시하지 않는다.
- tester 독립 QA 결과를 기다린다. 비용이 큰 전체 회귀를 중복 실행하지 않았다.
- DEP-001 대상 package.json은 test 체인만 바꾼다. 새 의존성이나 lockfile 변경이 없다. 표준/기존 도구로 충분하다.
- 예상 순증 400줄은 제품/테스트 범위다. 큰 전체 diff는 생성 packet 보존·갱신과 QA 자료다. 세 수정의 제품 코드는 작고 테스트/벤치가 대부분이다. 실제 lint added_lines=367이다. SIZE-002는 없고 remote_events.go SIZE-001 WARNING(471줄, 이전463)이 있다. 신규 추상화 없이 같은 claim 함수의 캐시/재예약만 추가했으므로 별도 파일 분할을 요구하지 않는다.
- UI 외형·테마·공용 컴포넌트를 바꾸지 않아 시각/테마 검증은 해당 없다. 디자인 lint 성공을 대신 주장하지 않는다.

## 대화 미참조 인계 점검

산출물 인덱스 `.fullops-squad/docs/deliverables/README.md`에서 시작했다.

| 항목 | 정본과 확인 | 결과·후속 |
|---|---|---|
| 현재 요구·결정 이유 | D02 사람 친구, to_dev, audit.md | 확인. 제품 규칙 변경 없음. 세 최소 수정 근거 존재 |
| 구조·구현 상태 | D03/D05/D07/D10, docs/system-design.md | 일부 확인. 기존 계약 일치. 새 pruning/expiry/cache/wake 정본 반영은 DOC 후속 대기 |
| 실행·검증 | project.md, docs/development.md, 하위 보고서 | 일부 확인. 현재 exact Linux/FullOps·독립 QA 원문 대기 |
| 운영·복구 | docs/operations.md, D07 | 확인. schema/writer/권한 복구 계약 유지. 이번 운영 배포·복원 미실행 |
| 다음 작업·담당 | to_dev, phase, 하위 보고서 | 확인. DEV 최종 report/packet-outcomes/정본/완료 archive, coor 독립 수락·통합 |
| 링크·절 | 인덱스와 실제 정본·보고 링크 | 관련 경로 직접 열람. 모든 과거 anchor 전수 렌더 미실행 |
| snapshot 이후 정본 | coor QA result/report/manifest | 확인. snapshot 유지. 정리 후에도 canonical 증거 보존 필요 |

후보에는 최종 부모 report, packet-outcomes, 빈 인박스/완료 archive가 아직 없다. 진행 중 상태로 기록한다. 최종 DOC SHA가 도착하면 별도 delta 리뷰로 문서 의미·링크·인계 완료를 확인해야 한다.

## 후속 lint와 테스트 수정 판단

실제 FullOps ANTI-004는 inbox-files.test.ts:39의 `inbox.pending()`를 skip/pending 테스트 표식으로 오인했다. 설치된 lint의 `\bpending\s*\(` 정규식과 실제 Inbox.pending 정의/호출을 대조했다. 테스트를 skip/only 처리한 코드가 아니다. 그렇더라도 실제 lint ERROR를 PASS로 수정하지 않는다.

DEV가 해당 pending 목록 assertion 한 줄을 제거한 미커밋 delta를 읽었다. 바로 다음 `readFile(<id>.json)`와 본문 equality가 원본 unread 파일/본문을 더 직접적으로 검사한다. 기존 제품 구현은 rename 뒤 marker를 비우고 prune에서 EACCES를 내므로 `.json` readFile에서 ENOENT로 실패한다. finally는 mock 복원만 하고 오류를 삼키지 않는다. 따라서 원래 본문 유실에 대한 회귀 보호는 유지된다. 후속 ENOENT 수신·빈 상태·중복 억제도 유지된다. 이 판단은 아직 새 고정 SHA/실행 통과를 대신하지 않으며 final delta에서 다시 대조한다.

## 결론

코드 의미 검토는 완료했다. bcb1dd6 FullOps 자체 실패를 유지한다. 최종 수락은 수정된 테스트와 최종 DOC 고정 SHA의 새 gate·독립 QA·인계 검토까지 보류한다. review.py check는 기록 일치이며 실제 검사 성공을 대신하지 않는다.

생략: 운영 계정/발송·배포·장시간 호스트/부하·전체 과거 문서 감사·중복 전체 테스트.
위험: singleton JSONB 처리량과 기존 local claim 뒤 crash 창이 남는다.
