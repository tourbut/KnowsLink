---
title: SAR-PUBLIC-AGENTS-001-FIX-REVIEW 리뷰
status: draft
updated: 2026-10-06
owner: ops
tasks: [SAR-PUBLIC-AGENTS-001-FIX-REVIEW]
summary: 고정 458798c AGENTS 후보의 d1eef9b 대비 독립 보안 delta 리뷰 결과와 수락 판단
---

# SAR-PUBLIC-AGENTS-001-FIX-REVIEW 리뷰

- 검토자 / CLI / 모델: OPS Orca dispatch `ctx_c958c7ee3895`(task `task_663b9832129e`, Run `run_8ca8bc058ab7`) / Claude Code / `claude-opus-5-5`. 리뷰 세션 `7dc8e8e4-8768-433e-a3d1-c36e6155cc43`이다. DEV 실제 세션 `01a10f52-…`, `e1275961-…`, `2ca6140e-…` 세 개와 모두 다르다.
- base SHA / head SHA / merge-base: `d1eef9bb90b9726149980320c42fb1fdbcaf584a` / `458798c2ee15c179edacfd6f94ebb9896d26f411` / `d1eef9b…`(base와 같음).
- snapshot: `/tmp/knowslink-agents-review-458798c`. detached `458798c`이고 리뷰 전후 `git status --porcelain` 0줄이다. 읽기만 했다. 설치·검사·프로브는 별도 scratch clone(같은 SHA, detached·clean)에서 실행했다.
- 후보 구성: 원본 DEV d1eef9b, 원본 OPS 70f26bc, 원본 QA bcb06b8, 원본 UI d165178, DEV-FIX 4a1b80a, 제품 답 48d12fa, DEV-POLICY-FIX 83e0bfb가 모두 458798c의 조상이다(`merge-base --is-ancestor` exit 0). `git diff 6d016e5 458798c -- . ':!.fullops-squad'`는 0줄이다. 제품 코드는 DEV 최종 6d016e5와 같다.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11 (a758d9c), `.fullops-squad/review/rule.json` sha256 `ab2116fb…`. OCR LLM은 쓰지 않았다. AI 검토자의 판단이며 OCR 자동 판정이 아니다.
- 적용 기준: `fullops-common-0.3.3`(README·coding-style·testing·security), FULLOPS, project, document-writing, review/rule.json. 제품 정본 D02 PS-07·PS-07 상세·기록 보호 조건, UX-04/05, POLICY 48d12fa의 두 관찰표. 예외 없음.
- 요구사항·완료 기준 원천: `handovers/to_ops.md`(SAR-PUBLIC-AGENTS-001-FIX-REVIEW), `SAR-PUBLIC-SERVICE.md`, `SAR-PUBLIC-SERVICE-UX.md`, `SAR-PUBLIC-AGENTS-001-POLICY.md`, `-DEV-FIX.md`, `-DEV-POLICY-FIX.md`, 원본 리뷰 70f26bc.
- 전체 변경 / reviewed / skipped: 152 / 87 / 65. skipped는 원시 증거 32개(prepare가 채움)와 원본 UI PNG 33개다.
- lint(`lint.json`): exit 0. head `458798c`, ERROR 0·WARNING 8·실행 불가 0. product-lint·product-test 모두 exit 0이다.
- SIZE-002: 추가 1505줄(상한 400). 그중 제품 Go 코드는 243줄이고 테스트는 599줄이다. 나머지는 역할 기록·Jev·원본 증거 보존이다. 이 범위는 DEV-FIX·POLICY·DEV-POLICY-FIX와 원본 QA/UI/OPS 기록을 하나의 고정 후보로 묶은 결과다. 개별 과제는 이미 분할돼 있으므로 다시 분할하지 않는다. 수락한다.
- DEP-001: 없음. `go.mod`·`go.sum`·`adapters/package*.json` diff 0줄이다.
- UI 디자인: 테마·공용 컴포넌트·새 라이브러리 변경이 없다. 디자인 lint 대상(DESIGN 경고)도 없다. memberStyle에 `overflow-wrap:anywhere`와 select 스타일만 추가됐다. DEV mobile-width(f9af9bf, 390px, overflow 0)는 f9af9bf..6d016e5 비테스트 제품 diff가 README 4줄뿐이므로 재사용할 수 있다. 직접 시각 판정은 designer 책임이며 이 리뷰는 PASS를 선언하지 않는다.

## 검토 범위

result.json의 모든 (path, status)에 검토 상태와 사유를 기록했다.

- 제품 코드 6개(`connections.go`, `http.go`, `store.go`, `api_rate.go`, `member.go`, `member_agents.go`)는 diff 전체와 호출 경로를 읽었다. `st.Agents[...]` 역참조 전부를 검색해 agent 삭제 뒤 nil 안전성을 확인했다.
- 테스트 3개(`connections_test.go`, `connections_integration_test.go`, `policy_test.go`)는 POLICY 표의 행과 대조했다.
- 기술 문서(D05–D10 해당 절)·D02·UX·exec-plans는 코드 상수·sweep·안내 문구와 대조했다. 제품 수치 변경은 없다.
- 역할 기록(PLANS·board·contexts·handover logs)은 삭제 줄을 확인했다. 삭제는 front matter 갱신뿐이다. handover logs는 삭제 0줄이다.
- 원본 OPS 리뷰 디렉터리는 70f26bc와 diff 0이다. 원본 QA(TESTER)는 bcb06b8과, 원본 UI는 d165178과 diff 0이다. 과거 결과를 바꾸지 않았다.
- Jev JSON은 구조와 과제 키만 확인했다. 탐색 보조이며 제외 근거로 쓰지 않았다.

### OCR 제외 증거의 무결성

- 원시 증거 32개는 파일마다 열지 않았다. 같은 폴더의 요약과 대조했다. DEV-FIX `head.txt` 85fb40e·`exits.txt` 전부 0. DEV-POLICY-FIX `head.txt` 6d016e5·run 2 lint/test/verify-mvp 0(Go integration race PASS 55·FAIL 0). TESTER `clone-before/after` 모두 d1eef9b, `probe-go.log` PASS. UI `strict.txt`·`artifact-check.txt`·`product-unchanged.txt` exit 0.
- 원본 UI PNG 33개는 `evidence.json` manifest의 sha256과 모두 일치한다(불일치 0·누락 0). `cleanup.json`은 자기 smtp·relay·Playwright 프로세스의 종료와 기존 사용자 자원 미변경을 기록한다.
- 증거 텍스트에서 password·secret·token=·개인키 패턴 검색 결과는 설명 문구 3건뿐이다. 비밀값은 없다.
- 생략 영향: 판정에 필요한 값은 요약 파일로 모두 확인했다. 영향 없음.

## 원본 발견 사항의 해소 판정

| 원본 | 판정 | 근거 |
|---|---|---|
| M1 medium 철회 기록 무한 증가 | 해소 | 철회 agent는 `Changed`부터 24h 뒤 sweep이 Key·Pair와 함께 삭제한다. owner 기록 10·agent 키 기록 20 상한이 생성만 막고 철회는 막지 않는다. 기존 철회 agent는 첫 sweep에서 `Changed=now`로 24h를 새로 시작한다. 반복 철회는 최초 시각을 유지한다. 살아 있는 agent의 철회 키는 삭제하지 않아 kid 재할당이 없다(C1). 살아 있는 두 agent 사이의 종료 pair는 남지만 pairID가 고정이어서 최대 C(200,2)로 유한하다. |
| L1 low 무효 세션 GET rate rollback | 해소 | `memberHit`가 무효 세션을 익명 IP budget에 집계한다. 거부 transaction은 원본을 재로드·sweep한 뒤 `Rates`만 보존한다. 부분 operation 상태는 남지 않는다. |
| L2 low 재초대 제품 판단 | POLICY대로 구현 | pending/active 반복은 기존 Pair를 반환하고 rate만 소비한다. 종료 뒤 수동 초대는 `Generation+1`과 새 24h 기한이다. 회원 화면 결정·철회는 같은 transaction에서 Generation 일치를 검사한다. 상한 초과는 세대 증가 전 409다. |
| L3 low connect 429 retry_at | 해소 | `rateLimited`가 실제 retry 시각을 초 단위 올림해 `Retry-After`와 `retry_at`에 같은 값을 준다. scratch 프로브로 경계(1ns 초과→다음 초, 정각→그대로)를 확인했다. |

추가 확인 결과는 다음과 같다.

- 권한: `connectLimit`은 `beginConnection`의 소유자·철회 검사 뒤에만 실행된다. 타인 agent의 키 기록 수를 노출하지 않는다. `agentLimit`의 순서(total→records→active)는 기존 `agentCapacity` 결과와 같다.
- 정보 노출: 거부 템플릿은 html/template 자동 escape를 쓴다. `Back`은 `/home`과 `/` 두 고정값이다. 홈의 `Own/Other`는 이미 표시되는 자기 관계 식별자다.
- 경합·재시작: 모든 변경은 기존 단일 row lock transaction 안에 있다. integration은 동시 생성 8건 중 4건만 허용, 재시작 뒤 상한 유지, 24h 경과 agent만 삭제를 실제 Postgres에서 검사한다.
- 삭제 뒤 경로: `prepareConnection`은 10분 기한의 `connection()` 검사를 먼저 통과해야 하며 삭제는 24h 뒤다. `current`·`ingest`·`/v1/keys`·gate 화면은 nil을 검사한다. 철회 agent의 envelope은 sweep에서 즉시 비워진다. gate 화면은 `원문 부재`로 표시된다.

## 이번 delta의 발견 사항

critical 0·high 0·medium 0·low 2다. 두 low 모두 미해결이며 수락을 막지 않는다.

### L-A — low: 삭제된 회원 agent ID를 합성 owner 경로가 다시 등록할 수 있다

- 위치: `internal/relay/connections.go:216-225`(sweep 삭제), `internal/relay/http.go:259-296`(`agents`).
- 재현: scratch 프로브 `TestOpsProbeDeletedIDReuse`. 회원 agent `agent_a`(철회 키 `k1`)를 철회하고 24h sweep으로 삭제한다. 합성 owner가 `operateAs("agents")`로 같은 ID·같은 kid를 등록한다. 결과는 `err=<nil>`이고 새 소유자는 합성 owner다([ops-probe.log](ops-probe.log), [원문](ops-probe_test.go.src)).
- 영향: 코드 주석의 "Member agent IDs are random and never issued again"은 합성 경로에서 강제되지 않는다. 합성 가입을 켜거나 legacy 합성 owner가 남은 DB에서는 C1의 `(from,kid)` 재사용 금지가 깨진다. 상대는 같은 식별자를 보고 새 초대를 수락할 수 있다. 옛 envelope과 pair는 이미 없으므로 옛 메시지 위조·관계 승계는 없다.
- 공개 경로: 회원 owner는 `Credential`이 비어 bearer가 없다. 공개 후보는 `KNOWSLINK_SYNTHETIC_SIGNUP`을 비운다. 따라서 공개 배포에서는 재현되지 않는다.
- 요청: OPS는 공개 전 DB의 합성 owner(Credential 비어 있지 않은 Owner) 0개를 확인한다. DEV는 후속에서 `agent_` 접두사 예약 또는 삭제 ID tombstone을 검토한다.

### L-B — low: `/v1/invite-decision`은 Generation을 결속하지 않는다

- 위치: `internal/relay/http.go:332-350`.
- 재현: scratch 프로브 `TestOpsProbeAPIDecisionGeneration`. 초대→API 거절→재초대(세대 2) 뒤 Generation 없는 API accept가 `state=active generation=2`를 만든다.
- DEV 판단 평가: 동의한다. owner API는 owner bearer가 필요하다. bearer는 합성 owner에게만 있다. 회원의 유일한 결정 경로 `/home/invite-decision`은 같은 transaction에서 Generation 일치를 검사한다. 공개 경로 제한으로 PS-07의 "옛 화면의 결정은 새 세대에 적용되지 않는다"를 충족한다.
- 조건: 이 충족은 L-A와 같은 배포 조건(합성 가입 꺼짐·합성 owner 부재)에 의존한다. owner API를 회원에게 공개하려면 Generation 필드를 먼저 추가해야 한다.

## 검증 및 남은 제약

reviewer scratch clone(`458798c`, detached·clean)에서 다음을 실행했다. 명령 자신의 종료코드를 기록했고 파이프로 가리지 않았다.

| 명령 | 결과 |
|---|---|
| `lint.py --repo <scratch> --from d1eef9b… --out <scratch>/lint.json`(설치 전 첫 실행) | exit 1. ERROR 2. product-test가 `sh: 1: tsc: not found`(make exit 127)로 실패했다. 새 clone에 `node_modules`가 없는 환경 원인이다. 원본은 [lint-initial-no-install.json](lint-initial-no-install.json)에 보존했다. |
| `make install` | exit 0. 리뷰 전과 같은 방법(원본 OPS 리뷰 70f26bc)이다. clone `status --porcelain` 0줄 유지. |
| 같은 `lint.py` 재실행 | exit 0. head `458798c2ee15c179edacfd6f94ebb9896d26f411`, ERROR 0·WARNING 8·실행 불가 0. product-lint·product-test exit 0. [lint.json](lint.json) |
| `go test ./internal/relay -run TestOpsProbe -v -count=1`(별도 scratch clone에 프로브 파일만 추가) | exit 0. 세 프로브 PASS. L-A·L-B 동작을 관측했다. |
| snapshot `git status --porcelain`·`rev-parse HEAD` | 리뷰 전후 0줄, `458798c` |
| `review.py check --key SAR-PUBLIC-AGENTS-001-FIX-REVIEW --from d1eef9b --to 458798c… --task-key SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX` | exit 0. reviewed 87·skipped 65·total 152·lint WARNING 8. Jev 적중률 recall 0.444·precision 0.333([jev-find-score.json](jev-find-score.json))은 수락 판단에 쓰지 않았다. |

- `make verify-mvp`는 다시 실행하지 않았다. DEV-POLICY-FIX run 2(6d016e5, exit 0, Go integration race PASS 55·FAIL 0)를 재사용했다. 6d016e5와 458798c의 제품 diff가 0줄이고 의존성 manifest diff도 0줄이기 때문이다.
- WARNING 8: SIZE-001 7개(PLANS 877줄, `connections_integration_test.go` 446, `connections_test.go` 389, `http.go` 573, `member.go` 339, `member_agents.go` 312, `store.go` 461)와 SIZE-002 1개다. `http.go`·`store.go`는 기존에도 상한을 넘었고 이번 증가는 각각 2·3줄이다. 테스트 파일 증가는 POLICY 관찰 행 고정 때문이다. PLANS 분할은 coor 판단 대상이다. 수락한다.
- 원본 QA bcb06b8·원본 UI d165178 FAIL/보류·원본 OPS 70f26bc는 원래 SHA 기록으로 유지한다. 최신 SHA 실행으로 바꾸지 않았다.
- 실행하지 않은 검증: UI 직접 시각 재검수(designer), 최신 후보의 독립 동작 QA(TESTER), 실제 JSONB 크기·CPU·복원 뒤 철회 보호(OPS 공개 수락 후속). D12/D13 인계: 이론 상한은 DEV 기록대로 키 20000개(회원 100×agent 기록 10×키 기록 20)다. 실제 운영 자원·복원 PASS는 이번 리뷰가 만들지 않는다. 롤백 시 `Agent.Changed`가 사라지고 재배포 시 24h를 다시 시작한다(보존이 짧아지지 않음).

## 검토 결론

수락 가능하다. 고정 458798c의 독립 delta에서 원본 M1·L1·L3이 해소됐고 L2는 POLICY 답대로 구현됐다. critical/high 0, medium 0이다. 미해결 low 2건(L-A, L-B)은 공개 배포에서 합성 가입이 꺼진 조건에서만 무해하다. 공개 전 OPS가 합성 owner 0개를 확인하는 조건으로 기록한다. lint ERROR 0, 고정 SHA 일치, 실행 불가 0이다. `review.py check` 통과는 기록 검사이며 이 판단과 테스트 결과를 자동 보증하지 않는다. 이 완료는 main 제품 수락이나 공개 수락이 아니다.
