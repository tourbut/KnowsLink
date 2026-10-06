---
title: SAR-PUBLIC-MESSAGES-001-REVIEW-2 리뷰
status: review
updated: 2026-10-06
owner: ops
tasks: [SAR-PUBLIC-MESSAGES-001-REVIEW]
summary: 복원된 일반 회원 메시지 후보 09c523d의 독립 보안 리뷰 — high 1건으로 수락 불가
---

# SAR-PUBLIC-MESSAGES-001-REVIEW-2 리뷰

- 검토자 / CLI / 모델: ops / Claude Code CLI / claude-opus-5-5. 실제 세션 `9bc44cbf-9e5c-4a16-a098-ae2fe4cab6ce`.
- 구현자 세션: DEV `01a11106-b736-7db1-b518-a65b64dbc5fb`. 검토자 세션과 다르다.
- snapshot: `/tmp/knowslink-messages-review-09c523d`. detached·clean·HEAD `09c523da8a3407288d9f5d711e1834af12bc7808`. 읽기만 했다. 설치·빌드·재현은 별도 scratch clone에서 했다.
- base SHA / head SHA / merge-base: `7efbaa349a5857eb1eac859a955ec3a09c91f800` / `09c523da8a3407288d9f5d711e1834af12bc7808` / `7efbaa349a5857eb1eac859a955ec3a09c91f800`.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11 delegate. `.fullops-squad/review/rule.json` (sha256 `ab2116fb…`). OCR LLM은 사용하지 않았다. 판정은 검토자 AI의 판단이며 OCR 자동 판정이 아니다.
- 공통 규칙: `fullops-common-0.3.3` README·coding-style·testing·security, FULLOPS.md, project.md, document-writing.md, contexts/ops.md. 지시서 준비 커밋 `a7443d5`에서 읽었다. 예외 없음.
- 요구사항·완료 기준 원천: SAR-PUBLIC-SERVICE.md PS-04/06/07·PS-08–11과 운영 기본값 표, SAR-MVP.md C1–C5, UX06/07, DEV 실행 기록과 logs/2026-10-06_to_dev.md, COOR/dev-final-37f9a1e-lint.json.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 89 / 51 / 38 / 51 / 38.
- lint(`lint.json`): ERROR 0 / WARNING 11 / 실행 불가 0. scratch clone(09c523d, `make install` exit0, 설치 뒤 clean)에서 `lint.py --from 7efbaa3` exit0. product-lint exit0, product-test exit0.
- SIZE-002: 지시서 예상은 제품/자동 검사 12–18파일이다. 실제 추가 2095줄이고 제품 파일은 20개다. DEV는 중앙 HTTP budget·cleanup·CLI bundle 근본 원인 수정을 이유로 들었다. 같은 권한·용량 경로를 나누어 미완료로 릴리스하지 않는다는 판단에 동의한다.
- DEP-001: 해당 없음. 의존성·버전·lockfile 변경 없음을 diff로 확인했다.
- UI 디자인: Go template의 기존 memberStyle을 쓴다. 별도 디자인 lint·테마 전환은 미구성이다(project.md). DESIGN 경고 0. 직접 시각 검수는 designer UX06/07 후속이다. 이 리뷰는 template의 escape·form 경로만 확인했다.

## 검토 범위

result.json의 89개 (path, status)를 모두 기록했다. 제품 코드 20개는 diff와 관련 함수 전체를 읽었다. 특히 `current`, `sweep`, `leaseMatching`, `leased`, `operateAs`의 persist/ack/claim, `ingest`, `transaction`을 읽었다.
지시서의 keep 후보를 모두 읽었다: public_text.go, capacity.go, http.go, store.go, adapters/src/text.ts. 실제 receipt 경로 `internal/relay/member_receipt.go`도 읽었다. omit? 후보 adapters/README.md는 지원 인터페이스의 근거로 코드와 대조했다.
문서 16개와 Jev/증거 JSON·HTML은 열어 형식·SHA·비밀 부재를 확인했다. 제외 38개(.log/.exit/.txt/.png)는 rule.json exclude 증거다. 같은 폴더의 `.exit` 값과 head.txt·versions.txt를 직접 열어 확인했다.

검토 항목과 결과:

| 영역 | 확인 내용 | 결과 |
|---|---|---|
| 인가 | text 송신·lease·persist·ACK·답장마다 양측 활성 Owner+활성 Member, agent·key 미철회, active pair·세대 확인. 합성 owner(회원 없음)는 text 불가 | 결함 없음 |
| text/wire 분리 | `knowslink.text.v1` closed schema·별도 서명 접두사·raw 32768·text 4096·TTL 180. generic pull/operate/claim/parentRouting에서 text 제외. frozen relay.v1 필드 무변경 | 결함 없음 |
| 답장·receipt·멱등 | 원요청이 delivered+persisted·현재·같은 세대·부모 기한 안일 때 1회. 같은 key·digest는 기존 receipt만 반환. ACK·만료·철회에서 원문 삭제, metadata 24h | 결함 없음 |
| 회원 receipt 화면 | 세션 회원의 agent·endpoint만. 본문 없음. html/template escape | 결함 없음 |
| gate | GET·POST가 저장 서명·digest·ID·세대·현재 권한·기한 재검증. hint escape. CSRF HMAC 유지. 원문 부재·만료·철회에서 approve 불가 | 결함 없음 |
| 상한 | queue100·gate100·receipt20000·claim4를 frozen/trial/text 공통 ingest·claim에 적용. high priority 우회 없음 | 결함 없음 |
| HTTP 입장·동시성 | 로컬 채널 16/4·공유 State.HTTP·30s 회수·10s deadline·rate 1회 차감 | H-1, M-1, L-1 |
| 원실패·경고 | DEV 원본 실패 exit와 수정 근거 보존. SEC-001/SLOP-004 오탐 판단 | 수락 |

## 발견 사항

### H-1 — high, 미해결: 익명 slow body가 rate·인증 전에 HTTP 슬롯을 점유한다

- 위치: `internal/relay/capacity.go:147-155` (`boundedHTTP`).
- 원인: 로컬 채널 슬롯을 먼저 얻는다. 그 뒤 `io.ReadAll`로 본문을 읽는다. 인증과 rate 차감은 본문을 다 읽은 뒤 DB에서 한다.
- 재현 조건: 익명 클라이언트가 `Content-Length: 100`과 본문 1 byte를 보내고 연결을 유지한다. 슬롯은 read deadline 10s까지 점유된다.
- 재현 결과: 고정 09c523d의 scratch clone에서 `slot-repro_test.go.txt`를 `internal/relay/`에 두고 실행했다. exit0이다.
  - 익명 `/v1/ack` slow body 4개 뒤 `/v1/text/ack` → `429 {"error":"capacity"}`.
  - 익명 `/auth/start` slow body 16개 뒤 `/v1/text/send` → `429 {"error":"capacity"}`.
  - 로그 끝의 nil Pool panic은 시험 종료 때 연결을 닫아 DB 없는 Service가 트랜잭션에 들어간 결과다. 재현 판정과 무관하다.
- 영향: 비인증 행위자가 연결 4개로 한 relay 프로세스의 ACK·revoke·deny·unpair·logout을 모두 막는다. 연결 16개로 모든 신규 작업을 막는다. 거부는 DB rate를 소비하지 않으므로 공격 비용이 없다. PS-11의 "신규 수락 포화 중에도 유효한 ACK·deny·철회·unpair·로그아웃은 별도 budget으로 처리"를 무효화한다. 이전 base에는 전역 슬롯이 없어 이 경로가 없었다. 이 후보가 만든 회귀다.
- 미확인: 운영 Cloudflare Tunnel/edge가 요청 본문을 완충하는지 확인하지 않았다. 완충해도 제품 handler의 결함은 남는다.
- 수정 방향(DEV 판단): 슬롯을 얻기 전에 상한된 본문을 deadline 안에서 읽는다. 또는 본문 수신 대기를 슬롯 밖에 둔다. slow body 회귀 검사를 추가한다.

### M-1 — medium, 미해결: 정리 입장을 경로만으로 판정해 익명 요청도 정리 수용량을 쓴다

- 위치: `internal/relay/capacity.go:93-105` (`cleanupRequest`), `:142-168`.
- 원인: `cleanupRequest`는 principal을 보지 않는다. 익명 요청도 `/v1/ack`·`/home/gates/*/deny` 등이면 로컬 정리 채널과 공유 `State.HTTP`의 Clean 기록을 차지한다. rate는 anonymousRate로 차감된다.
- 영향: PS-11 "신규 작업이 정리 budget을 소비하지 못한다"와 어긋난다. H-1을 고쳐도 빠른 익명 요청이 정리 4개 슬롯을 경쟁한다.
- 수정 방향: 정리 입장은 인증된 principal에만 허용한다. 익명 요청은 신규 입장으로 분류한다.

### L-1 — low, 미해결: 요청마다 전역 lock 트랜잭션 2개가 추가된다

- 위치: `internal/relay/capacity.go:164-190`, `internal/relay/store.go:193`.
- 내용: 입장·해제 트랜잭션은 각각 전체 JSONB 상태 decode와 sweep을 한다. sweep의 `current()`는 text 메시지마다 `publicOwner`로 Members를 순회한다. 해제가 2s 안에 끝나지 않으면 입장 기록이 30s 남는다.
- 처리: 기존 ponytail 단일 row lock 상한 안의 위험이다. OPS 공개 수락 전 상태 크기·CPU·처리량 측정에 포함한다.

### 보존하는 기존 low

- L-A(삭제된 회원 agent ID의 합성 owner 재등록)·L-B(`/v1/invite-decision` Generation 미결속)는 이 후보에서 코드 변경이 없다. 미해결로 보존한다.
- text는 양측 `publicOwner`를 요구하므로 합성 owner가 text 경로로 L-A를 확대하지 못한다.
- 공개 전 `KNOWSLINK_SYNTHETIC_SIGNUP` 미설정과 운영 DB 합성 owner 0 확인은 계속 필수다.

### 수락한 lint 경고

- SIZE-001 8건·SIZE-002: 위 SIZE-002 판단과 같다. 후속 분할은 DEV가 다음 변경에서 판단한다.
- SEC-001 `internal/relay/public_text_test.go:213`: `ClaimToken = "active"` 시험 fixture다. 비밀이 아니다.
- SLOP-004 `scripts/verify_grok_plugin.py:39`: 사용자에게 보이는 PASS 출력이다.

## 검증 및 남은 제약

| 검사 | 대상 HEAD | 결과 | 증거 |
|---|---|---|---|
| `make install` | 09c523d scratch clone | exit0, 설치 뒤 clean | install.exit |
| FullOps `lint.py --from 7efbaa3` (product-lint, product-test 포함) | 09c523d | exit0, ERROR0/WARNING11/실행불가0, `make lint`·`make test` 각 exit0 | lint.json, lint.log, lint.exit |
| `make verify-mvp` (격리 Compose·Postgres·실제 Node/MCP) | 09c523d | exit0, `--- PASS` 39, FAIL/SKIP 0. TestPublicNodeProcesses·HTTP 경계·gate 안전·공유 용량·claim·두 인스턴스 동시성 포함. 자기 volume/network 회수 | verify-mvp.log/.exit/.head |
| H-1 재현 | 09c523d + 리뷰 시험 파일 1개 | exit0 (결함 재현) | slot-repro_test.go.txt, slot-repro.log/.exit/.head |

로그는 각 명령의 stdout/stderr다. `git diff --check`를 위해 후행 공백만 제거했다. 종료코드는 `.exit`에 명령이 반환한 값 그대로 저장했다.

- 제품 코드는 DEV 원실행 SHA `ef5c571`과 09c523d 사이에 차이가 없다(`git diff --name-only ef5c571 09c523d`의 비 `.fullops-squad` 파일 0개). DEV의 lint/test/verify-mvp/plugin 원실행은 ef5c571 조건으로 유효하다. 이 리뷰는 같은 검사를 09c523d에서 다시 실행했다.
- DEV 원본 실패(lint-first exit2, test-first exit2, mvp-first/second exit2)는 기록 그대로 보존됐다.
- 실행하지 않음: `make verify-grok-plugin`(DEV ef5c571 원실행 exit0을 재사용한다. 이 리뷰 범위의 보안 판단에 필요하지 않다), 실제 Grok Bot·다닷 계정 왕복, 실메일, 운영 배포, 실24h 관측, Tunnel/edge 본문 완충, 실부하·백업/복원.
- skipped 38개는 원시 로그·종료코드·이미지 증거다. 판단에는 `.exit` 값과 요약 JSON을 사용했다. 생략 영향은 없다.
- `review.py check --key SAR-PUBLIC-MESSAGES-001-REVIEW-2 --from 7efbaa3 --to 09c523d --task-key SAR-PUBLIC-MESSAGES-001-DEV`는 exit1이다. 사유는 "critical/high 미해결 사항이 있습니다"이며 H-1 때문이다. 수락 차단이 정상 동작이다(review-check.log/.exit).
- 나머지 기록 조건을 확인하려고 result.json 사본에서 H-1만 임시 resolved로 바꾸어 같은 check를 실행했다. exit0, reviewed 51·skipped 38·total 89·lint WARNING 11이다. 실행 뒤 원본 result.json을 바로 복원했다(review-check-probe.log). 이 probe는 수락 근거가 아니다. jev-find-score.json(recall 0.5·precision 0.667)은 이 probe 실행이 생성했다. 수락 판단에는 쓰지 않는다.
- records 최종 HEAD의 lint/test·deliverables strict·diff 검사는 완료 보고에 기록한다. 기록 자체에 자기 SHA를 순환 기록하지 않는다.

## 검토 결론

수락 불가다. 미해결 high H-1이 main 수락을 차단한다. H-1은 이 후보의 HTTP 입장 구조가 만든 회귀다. 비인증 행위자가 PS-11의 정리 budget 보장을 무효화할 수 있다.
인가·text/receipt·답장·gate·queue/gate/receipt/claim 경계에서는 다른 critical/high를 찾지 않았다. M-1은 H-1과 같은 수정 범위에서 함께 처리하는 것을 권한다. L-1과 기존 L-A/L-B는 OPS 공개 전 조건으로 보존한다.
재개 조건: DEV가 H-1(필요 시 M-1)을 고치고 slow body 회귀 검사를 추가한 새 fixed SHA를 낸다. 그 뒤 OPS가 별도 세션의 새 리뷰 키로 재검토한다.
이 결론은 기술 리뷰 판단이다. TESTER 독립 QA·designer 직접 시각 검수·운영 공개 수락을 대신하지 않는다.
