---
title: SAR-MSG-EPHEMERAL-001-REVIEW 리뷰
status: draft
updated: 2026-11-12
owner: reviewer
summary: 9a5a07bd2239 병합 전 독립 리뷰 보고서
---

# SAR-MSG-EPHEMERAL-001-REVIEW 리뷰

- 검토자 / CLI / 모델: reviewer 세션 `reviewer-claude-opus-5-5-SAR-MSG-EPHEMERAL-001-REVIEW` / Claude Code / `claude-opus-5-5`. 구현자 세션 `db614234-2581-49da-a302-f5dcbd711da9`와 다르다. 읽기 전용 detached snapshot `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-f5a78d43818847f8b93d24a053d98c40`(HEAD `9a5a07bd`)만 읽었다.
- base SHA / head SHA / merge-base: `17acfa5e0b9800522854163b9d5df67b9870a159` / `9a5a07bd2239d10ff794149329270c411c4bdf23` / `17acfa5e0b9800522854163b9d5df67b9870a159`.
- OCR 버전 / 적용 규칙: open-code-review v1.13.7 delegate. `rules.json`의 6개 그룹(문서 정합, 설정 정합, Go, YAML 키, 일반 정확성, 스크립트)을 그룹당 한 번 읽었다. 공통 규칙은 `.fullops-squad/rules/common/README.md`·`coding-style.md`·`security.md`·`testing.md`(snapshot `9a5a07bd`)다. 지시서 예외는 Redis client 의존성 추가 허용 하나다.
- 요구사항·완료 기준 원천: `.fullops-squad/handovers/logs/2026-11-12_to_dev.md`의 SAR-MSG-EPHEMERAL-001 지시서(attempt `233358ae…`)와 결정 D049~D052. 기준별 판정은 아래 `완료 기준 대조`에 적었다.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 80 / 78 / 2(`go.mod`·`go.sum`, unsupported_ext) / 73 / 7. skipped 7건은 prepare가 채운 하네스 기록(jev packet·route 6건, `handovers/logs` 1건)이다. 지시서는 그 logs 파일에서 직접 읽었다. `go.mod`·`go.sum`은 diff를 직접 읽었다(`go-redis/v9 v9.23.0`과 간접 의존성 3개, `x/sys` 0.48.0). 지시서의 파일 소유권 밖 변경은 main 병합(`43ce01a7`)으로 들어온 SAR-VULTR-MIGRATE-001 문서뿐이다.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: ERROR 0 / WARNING 9 / 실행 불가 2(`make lint`·`make test`, Windows PATH에 `make` 없음). 사유는 coor가 보완했다. 제품 소스 `6fe42d32`를 Linux 컨테이너에서 실행한 결과가 exit 0이다. 로그는 sha256 `f27a317f…`·`14534b18…`로 기록되어 있다. `6fe42d32..9a5a07bd`는 보고서·완료 기록 커밋뿐임을 `git log`로 확인했다. 리뷰어는 재실행하지 않았다.
- SIZE-002: 지시서 예상 Go 600~1,200줄 / lint 기준 추가 909줄 / 범위 안이다. 시험·문서가 절반 이상이다.
- DEP-001: `go.mod`·`go.sum`만 바뀌었다. 표준 라이브러리에 Redis client가 없다. 버전 고정과 근거는 D03 `tech-stack.md`의 DEP-001 절에 있다. 수락한다.
- UI 디자인: 해당 없음. 지시서가 화면 구조 변경 없이 기존 안내 문장만 바꾼다고 정했다. 실제 diff도 템플릿 문장·i18n 문자열만 바뀌었다. 디자인 lint·테마 전환은 실행 대상이 아니다.

## 검토 범위

- 직접 읽은 코드: `ephemeral.go` 전문, `store.go` `transaction`·`sweep`·`State`, `text_events.go` `readEvents`, `cleanup_admission.go` `refresh`, `cleanup.go`(1초 sweep), `remote_auth.go` Inbox sweep, `remote_mcp.go` receive/ack, `capacity.go` 상한, `cmd/relay/main.go`, `internal/config/config.go`, migration 4.
- 직접 읽은 시험: `ephemeral_integration_test.go`, `ephemeral_test.go`, `state_migration_integration_test.go`, `policy_pages_test.go`, 나머지 시험의 `Redis: testRedis()` 주입·24시간→1시간 기대값 변경.
- 설정·스크립트: `compose.yaml`(Redis `--save ""`·`--appendonly no`·`noeviction`·`database` internal network), `compose.ops.yaml`(`mem_limit 192m`·`read_only`·`cap_drop`), `beta.sh` `message_state_rows`, `check_compose.py`·`verify_runtime.py`·`verify_mvp.py`.
- 문서: 정책 HTML 4개, 회원·랜딩 안내, README, `docs/*.md`, D03·D05·D06·D07·D08·D10·D11·D12.
- 제외 증거의 확인 방법: jev packet·route는 하네스 기록이라 열지 않았다. dev 보고서와 lint.json reason의 경로·hash로 존재를 확인했다. 원시 로그(`D:/workspace/KnowsLink/.devtools/msg-ephemeral/`)는 열지 않았다.

## 완료 기준 대조

| 기준 | 판정 | 근거 |
|---|---|---|
| Postgres에 본문·envelope·사본 없음 | 충족(본문 기준). metadata 예외 있음 | `Messages`·`Idempotency`·`remoteInstall.Inbox`가 `json:"-"`이다. `TestMessagesLiveOnlyInRedis`가 public 전체 테이블을 텍스트로 훑는다. 다만 Gate metadata(digest·대상)와 Dots 구독의 메시지 ID는 Postgres에 남는다(finding 1). |
| 사본 1시간, 전달 180초 | 충족 | `messageRetention = time.Hour`. 1초 cleanup transaction의 sweep이 수락 1시간 뒤 지운다. Redis TTL 1시간은 보조 상한이다. 180초는 기존 시험이 그대로다. |
| Redis 재시작·장애 | 충족 | 키 없음은 빈 저장소, 오류·미설정은 `unavailable`, epoch 불변 시험이 있다. `verify-runtime`은 stop 중 `/healthz` 503과 재시작 뒤 probe 키 0을 확인한다. 장애 범위 문장은 일부 부정확하다(finding 3). |
| 철회·친구 해제·키 철회 뒤 미전달 | 충족 | 전달 경로는 그대로 `st.current(m)`를 검사한다. Dots Inbox는 sweep이 `current`가 거짓이면 지운다(`remote_auth.go:129-133`). 기존 시험이 Redis 경유 상태로 통과한다. |
| 여러 process 원자성 | 충족 | 행 lock 안에서 GET, 새 키 SET, SaveRelay, commit 뒤 이전 키 제거 순서다. commit 전 실패는 새 키를 지우고, 모호한 commit 실패는 TTL로 묶는다. 잠금 없는 두 읽기 경로는 교체된 키를 다시 읽거나 이전 snapshot을 유지한다. |
| migration 4 데이터 보존 | 충족 | 세 값만 지우고 나머지 JSON·`human_model`·lz4를 보존한다. 1,001건 legacy 행 시험이 있다. Down은 실패 고정이다. 남은 Gate는 다음 sweep이 지운다(D052 범위). |
| 페이지 문구 ko·en | 대부분 충족 | `TestPolicyPages`가 핵심 문장을 확인한다. "승인 요청"의 Redis 전용 표기는 구현과 다르다(finding 1). |
| 검증 명령·FullOps lint | 충족 | 위 lint 항목. |

## 발견 사항

1. [준수] medium / `internal/relay/policy/privacy.ko.html:29`, `privacy.en.html:29`
   - 재현: 승인이 필요한 메시지를 보낸 뒤 1시간 안에 `beta.sh backup`을 실행한다. dump의 `relay_state.state.Gates`에 ID·Owner·Parent·Digest·From·To·Policy·State가 남는다.
   - 영향: 방침은 "승인 요청"을 Redis에만 두고 영구 저장하지 않는다고 적는다. 실제로는 본문 없는 metadata가 DB에 최대 약 1시간, 백업에 최대 14일 남는다. D07은 이를 바르게 적어 정책과 정본이 어긋난다. 본문 유출은 아니다.
   - 수정 상태: 미해결. 문장 분리 또는 Gate의 Redis 이동이 필요하다.
2. [버그] medium / `internal/relay/ephemeral.go:101`
   - 재현: Redis `used_memory`가 128mb에 닿게 메시지를 쌓는다. 그 뒤 ACK·sweep 삭제도 새 전체 사본 SET이 필요해 OOM으로 거부된다.
   - 영향: 저장소를 줄이는 쓰기가 막혀 현재 키 TTL(최대 1시간)까지 메시지 경로가 `unavailable`이다. 그 뒤 기록이 한꺼번에 사라진다. 고아 사본도 1시간 동안 한도를 쓴다. 128mb·192m 근거와 행 lock 안 전체 사본 처리 시간은 측정되지 않았다. 현재 베타 규모에서는 도달 가능성이 낮다.
   - 수정 상태: 미해결. 운영 이행 전 측정이 필요하다.
3. [준수] low / `.fullops-squad/docs/design-docs/interface-design.md:16`
   - "Redis에 닿지 못하면 모든 relay transaction이 unavailable"은 행에 Ephemeral 포인터가 있을 때만 맞다. 포인터가 비면 메시지를 쓰지 않는 transaction은 commit된다. 동작은 안전하고 문장만 고친다. dev 보고서의 같은 문장도 같다.
4. [nit] low / `.fullops-squad/board/board.json:357`: 완료 뒤에도 note가 "dev 구현 중"이다. nit은 이 1건이다.

보안 점검: Redis는 비밀번호 없이 `--protected-mode no`이다. 그러나 `database` network가 `internal: true`이고, 그 network에는 postgres·migrate·relay·redis만 있다. 포트도 게시하지 않는다. 지시서 요구와 맞다. `REDIS_URL` 오류 문구는 값을 출력하지 않는다(`TestRedisURL`). 로그에 본문을 쓰는 경로는 찾지 못했다.

## 검증 및 남은 제약

- 리뷰어는 lint·테스트를 다시 실행하지 않았다(coor 지시). 재사용 증거는 `lint.json`(head `9a5a07bd`, ERROR 0)과 dev 보고서의 Linux 실행 기록(`6fe42d32`)이다. 환경 동일성은 hash 기록으로만 확인했고, 원시 로그는 열지 않았다.
- 시험 삭제·완화: `state_migration_integration_test.go`의 envelope 정규화 비교를 map 비교로 바꿨다. migration 3 뒤 서명 검증은 유지했다. 24시간 기대값은 1시간으로 바뀌었다. D049에 따른 정당한 변경이다.
- 미검증: 실제 Redis 메모리 사용량·처리량, 운영 이행, 독립 QA. 모두 지시서가 후속으로 둔 범위다.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과(확인/미확인/해당 없음) | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `.fullops-squad/PLANS.md` "메시지 휘발성 Redis 보관", 지시서 D049~D052 | 확인 | 없음 |
| 구조와 구현/미완료 상태 | D07 `database-design.md` 과제 절, `docs/architecture.md` | 확인 | board note가 오래됨(nit) |
| 실행·검증 방법과 증거 | `qa-reports/SAR-MSG-EPHEMERAL-001-test/report.md` | 확인 | 원시 로그는 레포 밖 `.devtools` |
| 운영·복구(해당 시) | D12 `ops-guide.md` "SAR-MSG-EPHEMERAL-001 메시지 휘발 저장 이행 (운영 미적용)" | 확인 | relay 정지 뒤 migrate, 이행 직후 새 백업이 명시됨 |
| 다음 작업·담당·재개 조건 | `PLANS.md` 후속 2건(ops 이행, tester QA), dev 보고서 "남은 일" | 확인 | 이 리뷰의 medium 2건 처분 추가 필요 |
| 로컬 링크·절 접근/지원 한계 | dev 보고서의 D12 anchor | 확인 | 한글 anchor는 GitHub 렌더 기준 |
| snapshot 정리 후 정본 접근 | 이 폴더(coor 체크아웃) | 확인 | 없음 |

## 검토 결론

조건부 수락한다. critical/high는 없다. Postgres에 본문·envelope·사본을 쓰지 않는 구조, 행 lock 기반 원자성, Redis 장애 시 대체 저장 없음, migration 4의 보존 범위는 완료 기준과 맞다. `review.py check`의 통과는 기록 검사이며 AI 판단을 자동 보증하지 않는다.

미해결 medium 처분:

- finding 1(개인정보 문구): 병합 전 수정을 권장한다. 문장 한 줄과 en 대응 문장, `TestPolicyPages` 기대값만 바뀐다. 병합을 먼저 하면 운영 이행(사용자 승인 대상) 전까지는 반드시 고친다. 공개 페이지가 사실과 다른 상태로 배포되지 않게 한다.
- finding 2(Redis 한도 도달 시 축소 불가): 후속 과제로 둔다. ops 운영 이행 전에 크기·lock 시간을 측정한다. 축소 경로 설계(relay 쪽 사본 크기 상한 또는 메시지별 키)는 dev 후속으로 PLANS.md에 등재한다.

low 2건(finding 3·4)은 PLANS.md 후속 또는 병합 때 정리한다.
