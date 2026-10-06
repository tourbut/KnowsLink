---
title: SAR-PUBLIC-AGENTS-001-REVIEW 독립 보안 delegate 리뷰
status: review
updated: 2026-10-06
owner: ops
tasks: [SAR-PUBLIC-AGENTS-001-REVIEW]
summary: 고정 후보 d1eef9b의 회원 agent 연결·키별 credential·관계·한도 변경을 별도 세션에서 검토한 결과를 기록한다
---

# SAR-PUBLIC-AGENTS-001-REVIEW 리뷰

- 검토자 / CLI / 모델: OPS Orca dispatch `ctx_c8eb7eba79f3`(task `task_cc3d65402f06`, Run `run_8ca8bc058ab7`) / Claude Code / `claude-opus-5-5` high. 리뷰 세션 ID는 `50b08fc6-fec2-44ef-91ec-b921315867f9`다. DEV 구현 세션 `01a10f52-ac0f-75a0-b253-9a926a8e5650`(Codex)과 다르다.
- base SHA / head SHA / merge-base: `d2f7ba5aeb6644fd2b27fe6ada5b61db5976933b` / `d1eef9bb90b9726149980320c42fb1fdbcaf584a` / `d2f7ba5aeb6644fd2b27fe6ada5b61db5976933b`.
- snapshot: `/tmp/knowslink-agents-review-d1eef9b`. detached `d1eef9b`이고 `git status --porcelain` 0줄이다. 읽기만 했다. 설치·검사는 별도 scratch clone(같은 `d1eef9b`, detached·clean)에서 실행했다.
- OCR 버전 / 적용 규칙: open-code-review v1.12.11 (a758d9c). `rule_sha256` `ab2116fb…2aa0`. OCR LLM 호출은 사용하지 않았다. 판정은 AI 검토자의 판단이며 OCR 자동 판정이 아니다.
- 적용 기준: `fullops-common-0.3.3`(rules/common README·coding-style·testing·security), FULLOPS.md, project.md, document-writing, SAR-PUBLIC-SERVICE PS-04–07·PS-11·개인정보 절, UX-04–05. 기준 문서는 snapshot `d1eef9b`에서 읽었다. 제품 구현이 늘린 규칙을 정본으로 간주하지 않았다.
- 요구사항·완료 기준 원천: `handovers/to_ops.md`(이 과제), `docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md`, `docs/planning/product-specs/SAR-PUBLIC-SERVICE.md`, `docs/design-docs/interface-design.md`.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 39 / 33 / 6 / 33 / 6.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: 0 / 8 / 0. head `d1eef9b`, lint.py exit 0. 등록 `product-lint`(make lint)·`product-test`(make test) 모두 passed, exit 0.
- SIZE-002: 실행 기록의 최초 예상은 12–18파일이다. 실제 추가는 1614줄(lint 대상 25파일)이다. DEV는 2ac91a7 기준 1584줄을 보고했다. 차이 30줄은 malformed rate 보완 커밋이다. 발급·사용·철회·HTTP·UI·클라이언트를 한 보안 경계로 검증해야 한다는 분할하지 않은 이유가 실행 기록에 있다. 수락한다.
- DEP-001: `adapters/package.json`의 test script에 `node dist/connect.test.js`만 추가했다. dependencies·devDependencies·version·lock 변경은 0이다. diff로 확인했다. Node stdlib(crypto·fs·http)만 사용한다. 수락한다.
- SEC-001: `adapters/src/connect.test.ts:20`의 `"a".repeat(43)`이다. mock relay fixture의 합성 token이며 실제 자격이 아니다. 수락한다.
- SIZE-001(5건): PLANS 807·http 571·member 312·store 458·통합 테스트 317줄이다. 기존 대형 파일의 소폭 증가와 새 통합 테스트다. 보안 회귀를 줄이지 않는 한 분할을 요구하지 않는다.
- UI 디자인: Go template 일반 CSS이며 theme·Tailwind·shadcn·디자인 전용 lint가 없다. 기존 `memberStyle`을 재사용했다. DESIGN 경고 0은 Go 문자열 CSS를 검사하지 않은 결과이며 시각 근거가 아니다. 테마 전환은 해당 없음이다. UX-04–05 직접 시각 검수는 designer 후속이다.

## 검토 범위

result.json의 모든 `(path,status)`에 검토 상태와 사유를 기록했다. 33개 reviewed, 6개 skipped다.
skipped 6개는 `qa-reports/SAR-PUBLIC-AGENTS-001-DEV/*.log` 원시 로그다. exits.json의 1차·보완 종료코드 0과 각 로그 끝의 명령·PASS 줄을 대조했다. 로그에서 credential·token 값 패턴을 검색했다. 일치는 테스트 이름과 Docker image digest뿐이다. 후행 공백 정규화 사유는 exits.json에 있다.

보안 검토 항목과 확인 결과는 다음과 같다.

| 항목 | 확인 결과 |
|---|---|
| owner·session 결정 | 회원 경로는 세션에서 owner를 정한다. confirm·cancel·connection 조회·connect·agent-revoke는 `c.Owner`/`a.Owner` 일치를 검사한다. 다른 회원 ID 바꿔치기는 403/401이다. |
| 최근 재인증 | create·connect·confirm·key-revoke·agent-revoke에 5분 재인증이 있다. cancel·unpair·deny는 cleanup 경로이며 재인증이 없다. 실행 기록과 같다. |
| grant 결속·PoP·재사용 | token은 32바이트 난수이고 SHA256 해시로만 저장한다. PoP 바이트열은 token·owner·agent·client·mode·kid·public을 묶는다. prepare는 waiting에서 한 번만 가능하다. 준비 뒤 공개키를 바꿀 수 없다. complete는 approved에서 PoP를 다시 검증하고 consumed로 바꾼다. 같은 kid는 철회 기록을 포함해 재사용할 수 없다. 새 grant 발급은 같은 agent의 미완료 grant를 취소한다. |
| key credential·옛 key 차단 | `principal()`은 철회되지 않은 키의 credential만 인정한다. 공개 연결 완료는 agent-wide credential을 비운다. revoked agent·inactive owner는 인증에서 제외된다. 통합 테스트가 옛 credential의 pull·persist·ack·claim·authorize·gate-consume·send(query·result)를 재시작 뒤 401로 확인한다. |
| agent·pair 용량 | agent 200/owner 5, key 3, active pair 400/owner 20, pending 200/송신 owner 10, pending 24h를 상태 변경 지점에서 검사한다. 단일 row lock transaction으로 직렬화한다. 재시작 뒤에도 같은 JSONB 상태에서 집행한다. |
| rate·CSRF·입력 실패 | 회원 POST는 `CrossOriginProtection`과 SameSite=Strict·Secure·HttpOnly 세션 쿠키를 쓴다. 입력 실패·권한 거부는 값으로 반환해 rate 기록을 commit한다. `/v1/connect/*`는 8192 bytes strict JSON이다. 무효 grant는 source IP, 유효 grant는 회원 principal로 집계한다. `/v1/*`는 stable agent ID로 집계해 키 교체로 budget이 초기화되지 않는다. |
| 보존·화면 개인정보 | agent ID·owner ID는 이메일과 무관한 난수다. 화면은 html/template로 자동 이스케이프한다. token은 연결 대기 화면에 한 번 표시되고 `Cache-Control: no-store`다. credential·개인키는 화면에 표시되지 않는다. grant 기록은 만료 24h 뒤 삭제한다. 철회 agent·key·pair 기록은 삭제하지 않는다(M1). |
| Node 파일·URL 경계 | HTTPS root 또는 loopback HTTP root만 허용한다. userinfo·path·query·hash를 거부한다. fetch는 `redirect: "error"`다. stdin 입력은 128자 상한이다. 새 0700 폴더와 `wx` 0600 파일만 쓴다. complete는 폴더·파일의 소유자·권한·symlink를 lstat으로 검사한다. CLI 확인 URL의 해시 형식은 서버 `hashToken`(SHA256 base64url)과 같다. |

PS와 구현의 차이: PS-04·PS-05·PS-06·PS-07의 수락 조건은 구현과 테스트에서 확인했다. PS-11은 아래 L1·L3의 작은 차이가 있다. PS 본문 27·78행의 보존 기록 포화 보호는 구현되지 않았고 실행 기록이 공개 수락 후속으로 명시했다(M1).

## 발견 사항

### M1 — medium, 미해결: 철회된 agent·키·pair 기록이 삭제되지 않아 회원 행동으로 상태가 무한히 커진다

- 위치: `internal/relay/connections.go:158-183`(sweep은 Connection만 삭제), `internal/relay/member_agents.go:114-120`(create), `:140-151`(agent-revoke), `internal/relay/connections.go:144-151`(rotate가 철회 키를 남김).
- 재현 조건: 최근 재인증한 회원 한 명이 create(회원 신규 budget 40/60s)와 agent-revoke(cleanup budget 20/60s)를 반복한다. 철회 agent가 분당 최대 약 20개씩 쌓인다. Node CLI의 rotate 반복도 철회 키를 agent 아래에 계속 남긴다. 코드에 `delete(st.Agents…)`, `delete(a.Keys…)`, `delete(st.Pairs…)`가 없다.
- 영향: 모든 요청이 단일 JSONB 상태 전체를 global row lock 안에서 읽고 쓴다. `principal()`은 `/v1` 요청마다 모든 agent·key를 순회한다. 회원 100명이 반복하면 처리 지연과 공유 서비스 가용성 저하로 이어진다. 활성 개수 한도는 지키므로 권한 상승은 없다.
- 판단: PS는 철회 metadata 최소 24h 보존과 공개 전 OPS의 보존량 보호값 증명을 요구한다. DEV 기록은 보존량 검증을 공개 수락 후속으로 명시했다. 로컬 후보의 main 수락은 막지 않는다. 공개 전 차단 조건이다.
- 권고: 최소 24h 뒤 철회 agent·key와 관련 pair를 정리하거나 owner별 보존 철회 기록 상한을 둔다. 정리 시 `invite-decision`·`unpair`의 `st.Agents[...]` 역참조가 nil이 되지 않도록 pair도 함께 정리한다. DEV가 보존 기록 포화 테스트를 추가한다.

### L1 — low, 미해결: 무효 세션의 GET 연결 화면이 익명 rate에 집계되지 않는다

- 위치: `internal/relay/member_agents.go:187-191`. 같은 형태가 기존 `member.go` homePage에도 있다.
- 재현 조건: 무효 쿠키로 `GET /home/connections/{id}`를 반복한다. 오류를 반환하므로 transaction이 rate 기록 없이 rollback된다.
- 영향: PS-11의 "실패·거부도 집계"와 다르다. 상태를 바꾸지 않는 GET이며 POST 경로는 익명 budget을 소비하므로 영향은 작다.

### L2 — low, 제품 판단 필요: 거절·만료 뒤 같은 대상에게 즉시 재초대할 수 있다

- 위치: `internal/relay/http.go` operateAs `invites`.
- 재현 조건: B-owner가 거절한 뒤 A가 같은 대상에게 다시 초대한다. Generation이 늘어난 새 pending이 생긴다. B에는 차단 수단이 없다.
- 영향: 송신 owner pending 10·rate 40/60s로 제한된다. PS 67행 "반복 초대는 새 초대를 만들지 않는다"가 pending 중 반복만 뜻하는지 거절 뒤 반복도 포함하는지 불명확하다. designer 판단 대상이다.

### L3 — low, 미해결: `/v1/connect/*`의 429 응답이 실제 재시도 시각을 주지 않는다

- 위치: `internal/relay/member_agents.go:258-263`.
- 재현 조건: grant 경로 rate 초과. `Retry-After`는 항상 60이고 body는 `{"error":"rate_limited"}`만 반환한다. `/v1/*` rateAPI는 실제 시각과 `retry_at`을 준다.
- 영향: interface-design의 "429는 재시도 시각을 표시한다"와 다르다. 클라이언트가 실제보다 늦게 또는 이르게 재시도할 수 있다. 보안 영향은 없다.

### 관찰 — 수정 요청 아님

- `client=node-local`은 선언 문자열이다. 클라이언트 증명이 아니다. 실제 결속은 owner 지문 확인과 PoP다. 설계·문서와 일치한다.
- `/v1` rate는 agent ID별 principal이다. 회원 한 명이 agent 5개로 회원 budget 외에 agent별 40/60s를 더 쓴다. PS의 "인증 principal" 정의 안이다. 전체 200/60s가 상한이다.
- 전체 신규 200/60s는 PS 수치다. 익명 IP 7개 정도로 전체 budget을 채울 수 있다. 정리 budget은 분리돼 있다. 공개 전 OPS가 IP 기반 보호와 함께 판단한다.

## 검증 및 남은 제약

reviewer scratch clone(`d1eef9b`, detached·clean)에서 다음을 실행했다. 명령 자신의 종료코드를 기록했고 파이프로 가리지 않았다.

| 명령 | 결과 |
|---|---|
| `make install` | exit 0 |
| `python3 <플러그인>/scripts/lint.py --repo <scratch> --from d2f7ba5… --out <scratch>/lint.json` 뒤 리뷰 디렉터리로 원본 복사 | exit 0. head `d1eef9b`, ERROR 0·WARNING 8·실행 불가 0. product-lint·product-test passed |
| `make verify-mvp` | exit 0. HEAD `d1eef9b`. 고유 Compose project, 자기 자원 down --volumes. Go PASS 49줄·FAIL 0. TestPublicAgentHTTP 5개 하위·TestConnection*·TestAgentAndPairCapacity PASS |
| snapshot `git status --porcelain` | 리뷰 전후 0줄, HEAD `d1eef9b` |

DEV의 원래 증거는 exits.json의 1차(2ac91a7 당시)와 보완(d1eef9b) 결과를 구분한다. 이번 재실행은 d1eef9b의 독립 재현이다. DEV 증거를 새 SHA 결과로 대체하지 않는다.

실행하지 않은 검증: 실제 이메일 발송, 운영 공개·배포, Grok Bot·다닷·OAuth·외부 계정 연결, 실제 브라우저 시각 검수, 처리량·보존량 부하 시험. TESTER의 교차 계정·회전·철회·관계·한도 QA와 designer의 UX-04–05 직접 검수는 별도 과제다. fixture PASS를 실제 플랫폼 PASS로 표시하지 않는다.

`review.py check`는 기록 검사다. 이 문서의 서술과 테스트 결과의 수락 가능성을 자동 판정하지 않는다.

## 검토 결론

critical·high 발견은 없다. 연결 grant 결속·PoP·1회 완료, key별 credential과 철회 키의 전 API 차단, 소유권 분리, 용량, rate·CSRF, Node 파일·URL 경계는 정본과 일치한다. 고정 SHA `d1eef9b`의 보안 리뷰는 로컬 후보 main 수락 조건으로 수락 가능하다.
M1은 공개 전 차단 조건으로 OPS·DEV 후속에 남긴다. L1·L3은 DEV 후속 수정 후보다. L2는 coor를 거쳐 designer가 판단한다. 최종 main 수락은 TESTER QA와 designer UI 검수 결과를 함께 따른다.
