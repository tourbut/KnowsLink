---
title: SAR-MVP-003-BIDIRECTIONAL 리뷰
status: review
updated: 2026-10-04
owner: ops
tasks: [SAR-MVP-003-BIDIRECTIONAL-REVIEW]
summary: 실제 시험 transport와 운영 인증 준비의 고정 SHA 독립 리뷰 결과 critical/high 없이 수락 가능하다
---

# SAR-MVP-003-BIDIRECTIONAL 리뷰

- 검토자 / CLI / 모델: ops 독립 검토 세션 `c4c411b8-d66d-4356-981b-9e2279440b02` / Claude Code / `claude-opus-5-5` high. 구현자 Codex 세션 `01a104d1-8cc5-7430-bcf1-1a2732c5183f`와 다르다.
- base SHA / head SHA / merge-base: `f2849486ed48295e239714700e651d30c32f1c2c` / `cd60e7f87eb5ce137eca887980f232b3f67a18d0` / `f2849486ed48295e239714700e651d30c32f1c2c`.
- snapshot: `/tmp/knowslink-bidirectional-review-cd60e7f` detached·clean·읽기 전용으로 유지했다. 실행 검증은 별도 scratch clone에서 수행했다.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11 (a758d9c) delegate. `review/rule.json` sha256 `ab2116fb…`. 공통 규칙 `fullops-common-0.3.2`(README·coding-style·testing·security), `project.md`. 두 문서는 대상 SHA와 기록 체크아웃에서 차이가 없다. 예외 없음.
- 요구사항·완료 기준 원천: `handovers/to_ops.md`, [DEV 실행 기록](../../exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md), D03/D05/D10, [D12 13장](../../operations/ops-guide.md#13-승인된-양방향-시험-sar-mvp-003). Jev REVIEW find/documents-find/context keep을 읽었다. 충돌은 없다.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 47 / 38 / 9 / 47 / 0. 커버리지 100%.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: 0 / 2 / 0. WARNING은 기존 `http.go` 487→510줄, `store.go` 420→438줄의 SIZE-001이다.

## 검토 범위

result.json의 47개 (path, status)를 모두 reviewed로 기록했다. 제외 9개는 OCR `unsupported_ext` 실행 증거 `.txt`다. 직접 열어 비밀값 패턴이 없음과 exit 0 기록이 실행 기록 표와 일치함을 확인했다. 같은 명령을 리뷰어가 재실행해 대체 검증했다.

지시서의 중점 항목별 결과는 다음과 같다.

| 항목 | 확인 내용 | 결과 |
|---|---|---|
| 인증·권한 | `/v1/test/*`는 agent 인증 후 서버 allowlist를 요구한다. send/pull/persist/ack/claim과 GET registry/keys/receipts만 허용한다. signup·owner·pairing·키변경·업무 send·authorize·gate-consume은 403이다. | 통과 |
| 경로 우회 | `..`·중복 slash·잘못된 method·allowlist 밖 key 대상 등 13개 probe가 모두 403이었다. Go mux는 비정규 경로를 handler 전에 redirect한다. | 통과 |
| 두 agent 권한 | allowlist는 정확히 서로 다른 두 AgentID만 허용한다. 파싱 오류 시 relay 기동이 실패한다. ingest·current·operate가 양 endpoint를 모두 검사한다. | 통과 |
| 데이터 변경 | `TestAgents`는 DB로 직렬화하지 않는다. 오류 rollback에도 다시 적용한다. allowlist 제거 시 trial은 `failed:revoked`와 원문 삭제로 fail-closed였다(probe). | 통과 |
| 본문 수명 | claim 즉시 Envelope·Inbox를 지운다. 만료·철회·max attempts에도 지운다. metadata는 24시간 뒤 삭제한다. trial은 H/R·authorize·gate-consume 부모가 될 수 없다. | 통과 |
| remote origin·SSRF | `test-remote`는 `https://link.knowslog.com` root만, `test-loopback`은 HTTP loopback root만 허용한다. recipient는 env 고정이다. redirect는 error다. | 통과 |
| 10s body timeout·64 KiB | `AbortSignal.timeout`이 header 이후 body 읽기에도 적용된다. 스트림 누적 65536 bytes 초과 시 거부한다. 서버 입력도 `MaxBytesReader` 64 KiB다. 100ms body timeout 시험을 재실행했다. | 통과 |
| Access plan | 두 distinct service token, `non_identity` policy, `/v1/test/*` path 앱, trial rule 선행과 distinct AUD를 확인했다. 기존 원점 rule은 owner AUD만 받으므로 trial JWT를 거부한다. trial 경로는 owner identity를 받지 않는다. | 통과, low 1 |
| 기존 원점의 trial JWT 수용 후보 | edge가 trial로 분류해도 cloudflared가 trial rule을 고르지 않으면 owner AUD 검사로 거부된다. trial rule을 통과한 요청은 relay prefix 판정 뒤 testHandler 또는 비정규 경로 redirect로 끝난다. owner/signup 도달 경로를 찾지 못했다. | 통과 |
| 댓글·실행 정합 | 댓글 초안의 `install_bot_mcp.sh`는 `npm run build`로 `dist/trial-cli.js`와 `plugin.js`를 만든다. `run_trial.py`는 config의 mode를 상속 env보다 우선한다. idempotency key 예시는 regex를 만족한다. 비밀값을 argv·댓글에 넣지 않는다. | 통과 |

## 발견 사항

critical·high·medium은 없다. 다음 low 5건은 모두 미해결이며 수락을 차단하지 않는다.

1. low / `internal/relay/store.go:298-310` — 업무 `/v1/pull`은 `relay.test.message`를 제외하지 않는다. 재현: trial 송신 뒤 같은 수신 agent가 업무 `/v1/pull`을 호출하면 trial lease를 받는다(probe 로그 `OBSERVED business /v1/pull leased trial intent`). 업무 `once()`가 이를 claim하면 원문을 지운 뒤 approval 단계에서 실패한다. 업무 경로는 owner Access 뒤에 있고 synthetic 모드는 loopback 전용이다. 권한 상승·원문 노출은 없고 운영자 오설정 시 시험 메시지 유실만 생긴다. 업무 pull에서 시험 intent 제외를 권장한다.
2. low / `internal/relay/store.go:274-279` — trial claim 응답의 `policy`가 `schedule.commit` 라벨이다. TS 수신은 이 값을 쓰지 않는다. 진단 오해만 생긴다.
3. low / `deploy/knowslink/access_trial_plan.py:36-37` — cloudflared `path: /v1/test/.*`는 고정되지 않은 regex다. 경로 중간에 `/v1/test/`가 있으면 trial AUD를 요구한다. relay의 HasPrefix 판정과 범위가 다르지만 결과는 owner 요청 거부(fail-closed)다. `^/v1/test/` 고정을 권장한다.
4. low / `scripts/run_trial.py:17-20` — config는 0600 소유 regular file로 검사하지만 `AGENT_KEY_FILE` 권한은 검사하지 않는다. lstat와 read 사이에 교체 TOCTOU가 있다. 같은 UID 로컬 공격자를 전제로 하며 D12 13.3의 0700/0600 확인으로 완화된다.
5. low / `adapters/src/core.ts:76` — 비정상 HTTP 응답에서 본문을 cancel하지 않고 throw한다. 소켓이 GC까지 남을 수 있다. 수동 단발 호출이라 영향이 낮다.

## 검증 및 남은 제약

리뷰어가 scratch clone(HEAD `cd60e7f`, 깨끗한 작업 트리)에서 실행한 명령과 결과다. 로그는 [evidence](evidence/)에 있다. 종료코드는 명령별로 직접 기록했다.

| 명령 | 결과 |
|---|---|
| `npm ci --prefix adapters` | exit 0, 작업 트리 변경 없음 |
| `lint.py --repo <scratch> --from f284948… --out lint.json` | exit 0, ERROR 0 / WARNING 2 / 실행 불가 0 |
| `make test` | exit 0, Go race unit과 MCP·trial-boundaries |
| `make verify-mvp` | exit 0, 격리 Postgres·TrialHTTP·두 독립 MCP 왕복 `01a104f1-3899-…`/`01a104f1-397b-…` |
| `access_trial_plan.py selftest` | exit 0, render-only |
| 리뷰어 probe(`TestReviewProbe`) + `make verify-mvp` | exit 0, 우회 13건 403·업무 pull 교차 관측·allowlist 제거 fail-closed |

probe 파일은 제품에 넣지 않았다. [원문](evidence/review_probe_integration_test.go.txt)만 보존한다.

실행하지 못한 검증은 다음과 같다. 실제 Grok Bot 계정 왕복, 실제 Cloudflare service token·trial path 앱·trial AUD로 인증된 원격 호출, 운영 relay 후보 배포다. 이 리뷰는 이를 성공으로 확장하지 않는다(actualGrok 미검증). 실제 Access 적용 뒤 D12 13.2-7의 negative/positive 검사가 별도로 필요하다. Cloudflare edge의 path 정규화·대소문자 처리는 공식 문서 확인 없이 fail-closed 구조로만 판단했다. tester의 독립 동작 QA와 직접 시각 검수는 이 리뷰 범위 밖이다.

## 검토 결론

고정 head `cd60e7f87eb5ce137eca887980f232b3f67a18d0`은 수락 가능하다. 미해결 critical/high가 없고 target lint ERROR가 0이다. 필수 테스트를 재실행해 통과했다. low 5건은 후속 개선 후보로 남긴다. 게시·배포 가능 범위는 승인된 Codex↔Grok 시험 text와 D12 13장 절차에 한정한다. 업무 effect·자동 wake·실제 Grok 왕복 성공 주장은 포함하지 않는다. check 통과는 기록 검사이며 AI 검토 내용과 테스트 성공을 자동 보증하지 않는다.
