---
title: 철회 기록 보존·rate·회원 화면 후속 수정 기록
status: draft
updated: 2026-10-06
owner: dev
tasks: [SAR-PUBLIC-AGENTS-001-DEV-FIX]
summary: "철회 agent 24h 정리와 기록 상한, 거부 요청 rate 저장, connect 429 재시도 시각, 회원 화면 F-UI-01–04 수정과 검증을 기록한다"
---

# SAR-PUBLIC-AGENTS-001-DEV-FIX — 철회 기록 보존·rate·회원 화면 후속 수정 기록

## 기준과 기술 계획

기준 ref는 `d1eef9bb90b9726149980320c42fb1fdbcaf584a`다. 준비 HEAD는 `b5df23b`다. d1과 main/origin `dc60fbf`가 준비 HEAD의 조상임을 확인했다. Task는 `task_356146f610c0`, Dispatch는 `ctx_3bec7b292d75`다. 적용 기준은 fullops-common-0.3.3·FULLOPS·project·document-writing·PS-04–07/PS-11/개인정보·UX-04–05·MVP C1/MVP-05다. 정본 인박스는 `handovers/to_dev.md`다.

입력 발견은 [독립 리뷰](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md)의 M1·L1·L3이다. 작업 중 coor가 같은 과제에 designer 직접 검수 F-UI-01–04를 추가했다. 원천은 designer 체크아웃의 `SAR-PUBLIC-AGENTS-001-UI.md`(기록 `f9c5eded`)이며 읽기만 했다. 원본 리뷰·QA·UI 기록과 PNG는 바꾸지 않았다.

계획: 기존 singleton sweep·transaction·memberStyle·Go template을 재사용한다. 새 라이브러리·migration·wire 필드는 없다. 예상 규모는 Go 6파일과 테스트 2파일, 기술 문서 6–8개였다. 실제 제품 변경은 8파일 +425/−77줄이다(테스트 +267). Context7 조회는 새 라이브러리가 없어 해당 없음이다.

## 원인과 수정

### M1 철회 기록 무한 보존

원인: sweep은 Connection만 삭제했다. 철회 agent·키·관계 기록은 삭제하지 않았다. create/agent-revoke 반복과 회전 반복이 단일 JSONB와 `principal()` 순회를 계속 키웠다.

C1 경계: MVP-01/C1은 `(from,kid)`의 kid 재할당을 영구 금지한다. 살아 있는 agent에서 철회 키를 지우면 같은 kid를 다른 키에 다시 등록할 수 있다. 따라서 살아 있는 agent의 철회 키는 지우지 않는다. 철회 agent 전체는 지운다. 회원 agent ID는 무작위 28자이며 다시 발급되지 않으므로 `(from,kid)` 재할당이 생기지 않는다.

수정:

- Agent.Changed에 철회 시각을 기록한다. 반복 철회·회전은 Agent·Key의 첫 철회 시각을 바꾸지 않는다(`revokeKeys`).
- sweep은 철회 시각부터 24h가 지난 agent를 키와 함께 삭제한다. 그 agent를 포함한 Pair도 같은 sweep에서 삭제한다. Changed가 없는 기존 철회 agent는 첫 sweep 시각부터 24h를 보존한다.
- 살아 있는 두 agent 사이의 거절·만료·철회 Pair는 유지한다. 재초대 Generation이 계속 증가하므로 옛 세대 화면을 재사용할 수 없다.
- `invite-decision`·`unpair`는 없는 agent를 가리키는 pair를 nil 역참조 없이 거부한다.
- 기술 보존 상한: owner당 agent 기록 10개, agent당 키 기록 20개다. 상한이면 회원 create·connect 발급·complete와 합성 `/v1/agents`·`/v1/keys`가 `capacity`(409)다. 철회·취소·거절·unpair는 기록을 늘리지 않으므로 포화 중에도 정리 budget으로 처리한다.

receipt·lease·claim·gate 의존: 메시지는 수락 24h 뒤 삭제된다. 철회 agent는 수락이 불가능하므로 그 agent의 메시지는 철회 전에 수락됐다. 따라서 agent 삭제 시점에 그 agent의 메시지·gate·멱등 기록은 이미 없다. `current()`·gate 화면·contacts·keys 조회는 없는 agent를 이미 거부한다. Connection은 `connection()`에서 없는 agent를 invalid_auth로 거부한다.

보호값 근거: 제품 활성 한도(agent 5, 키 3)를 그대로 포함한다. agent 기록 10개는 24h 안에 활성 5개를 한 번 더 교체할 수 있는 값이다. 키 기록 20개는 활성 3개와 회전 기록 17개다. 회원 100명 기준 이론 최대는 agent 1000개·키 20000개다. 이 값은 제품 quota를 바꾸지 않는다. 키 기록이 찬 agent는 회전할 수 없으며 새 agent로 교체한다. 이 사용자 영향은 designer 판단 후속으로 인계한다. 실제 JSONB 크기·CPU·DB·복원은 측정하지 않았다(OPS 후속).

### L1 무효 세션 요청의 rate 미집계

원인: `GET /home/connections/{id}`와 `GET /home`은 세션 오류를 바로 반환했다. transaction은 오류 시 상태 전체를 되돌리므로 rate 기록도 저장하지 않았다. 같은 형태가 `POST /auth/reauth`·`/auth/logout`·`/auth/logout-all`에도 있었다.

수정: transaction은 거부된 요청에서도 sweep 결과와 사용한 rate 기록만 저장한다. 부분 업무 상태는 계속 버린다. 회원 세션 처리기는 `memberHit` 하나로 세션 확인과 budget 소비를 수행한다. 유효 세션은 회원 신규 또는 정리 budget을, 무효 세션은 source IP 익명 budget 30/60s를 쓴다. 무효 세션 GET은 기존처럼 `/?n=expired`로 이동하고 budget 초과 시 429다.

확인한 미적용 경로: `GET /`, `GET /auth/verify`, 회원 `GET·POST /home/gates/{id}`, owner `/owner/*`에는 원래 HTTP rate가 없다. 이번 L1 형태(세션 실패의 rollback)가 아니다. gate approve·deny의 budget 분류는 PS-11 MESSAGES 범위이므로 바꾸지 않았다. coor에 후속 후보로 보고한다.

### L3 connect 429 재시도 시각

원인: `/v1/connect/*`는 `Retry-After: 60` 고정값과 retry_at 없는 body를 반환했다.

수정: `/v1/*`와 같은 `rateLimited`를 사용한다. body는 `retry_at`(RFC3339 UTC), header `Retry-After`는 같은 시각의 HTTP-date다. 두 값은 실제 시각을 초 단위로 올림한다. 내림하면 그 시각의 재시도가 다시 거부·집계되기 때문이다. Node CLI는 실패 시 공통 안내만 출력하며 retry 값을 읽지 않는다. CLI 동작과 안내 문구는 바꾸지 않았다.

### F-UI-01–04 회원 화면

| ID | 원인 | 수정 |
|---|---|---|
| F-UI-01 medium | 홈의 키 지문이 `<code>` 밖 일반 텍스트라 줄바꿈되지 않았다 | 지문을 `<code>`로 표시하고 body에 `overflow-wrap:anywhere`를 적용했다. select도 전체 폭 block이다 |
| F-UI-02 medium | 거부 화면이 `head` template만 렌더링해 다음 동작 컨트롤이 없었다 | `refusal` template과 `refused`를 추가했다. 홈 링크(세션 무효면 로그인 링크)를 표시한다. 재확인 필요와 전체 로그아웃 재확인은 `POST /auth/reauth` 버튼도 표시한다. 안전 문구·상태 코드는 같다 |
| F-UI-03 low | 취소 버튼 조건이 `consumed`만 제외했다 | waiting·prepared·approved에서만 표시한다 |
| F-UI-04 low | 관계 기한이 time.Time 기본 문자열이었다. select가 브라우저 기본 13.33px였다 | `memberPair.Deadline`이 연결 기한과 같은 KST 형식을 쓴다. select에 본문 글꼴을 상속한다 |

기존 Go template·memberStyle만 수정했다. 별도 theme·Tailwind·shadcn·디자인 전용 lint는 해당 없음이다. DESIGN 규칙은 Go 문자열 CSS를 검사하지 않으므로 아래 좁은 폭 측정을 남겼다. 이 측정은 designer 직접 시각 PASS가 아니다.

## L2 대기

거절·만료 뒤 같은 대상 재초대 허용 범위는 designer 결정 대기다. 코드는 바꾸지 않았다. 결정이 오면 같은 과제의 후속으로 처리한다.

## 자동 검증 증거

레포 루트에서 실행했다. 명령의 원래 종료코드를 저장했고 파이프로 가리지 않았다. 로그는 [QA 증거](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-DEV-FIX/)이며 대상 SHA는 `head.txt`다. 로그 후행 공백만 정규화했다.

| 명령 | 결과·대상 |
|---|---|
| make lint | exit 0, 제품 HEAD `85fb40e` |
| make test | exit 0. Go race 전체와 adapter 검사 |
| make verify-mvp | exit 0. 격리 Postgres·migration·Go integration race(PASS 53·FAIL 0)·TS 왕복 |
| 390×844 폭 측정 | `mobile-width.json`. 같은 홈 상태의 scrollWidth가 d1에서 595, 85fb40e에서 390이다. select 13.33px → 18px |

새 검사:

- `TestRevokedRecordRetention`: 회전 반복의 키 기록 상한, 포화 중 철회, 반복 철회 시각 유지, 30일 뒤 살아 있는 agent의 철회 kid 유지, owner agent 기록 상한, JSON 재시작 뒤 24h 경계 삭제, 관련 pair 삭제와 살아 있는 pair 세대 유지, 기존 철회 agent 보존 시작, 없는 agent pair 거부, 합성 키 경로 상한.
- `TestMemberPagesShowNextSteps`: 재확인 버튼·홈/로그인 링크, 열린 상태만 취소, 지문 `<code>`, KST 기한, 스타일 규칙.
- `TestPublicAgentHTTP/revoked_record_saturation_concurrency_restart_and_retention`: 실제 Postgres에서 철회 6개 상태의 동시 생성 8건 중 기록 상한으로 정확히 4건만 성공한다. 재시작 뒤 상한이 유지된다. 포화 중 agent 철회는 303이다. 24h가 지난 agent만 삭제된다.
- `TestPublicAgentHTTP/invalid_session_pages_spend_anonymous_budget`: 무효 세션 GET 30회 뒤 재시작한 relay에서 두 GET·logout·reauth가 429다.
- `malformed_connection_requests_spend_anonymous_budget`: 31번째 429의 retry_at과 Retry-After가 같고 60초 안의 미래 시각이다.

## QA·UI·운영 인계

- OPS: 최종 fixed SHA의 독립 delta 리뷰가 필요하다. 보존 상한·C1 판단·rate 저장 범위 변경을 확인한다.
- TESTER: 철회 24h 경계·기록 포화 중 신규 거부와 정리 허용·무효 세션 rate·connect 429 재시도의 좁은 QA가 필요하다.
- designer: F-UI-01–04 변경 화면만 좁게 재검수한다. 기존 d1 시각 결과와 PNG는 원래 SHA 근거로 유지한다. 키 기록 20개 포화 시 agent 교체 안내 문구와 홈에서 24h 뒤 철회 agent가 사라지는 동작의 제품 판단도 함께 요청한다.
- D12 운영 경계(구현 기술 인계): 철회 agent는 24h 뒤 DB에서 삭제되며 운영자가 수동으로 복구하지 않는다. 이전 SHA로 rollback하면 Agent.Changed를 잃는다. 다시 올리면 보존을 다시 24h 시작하므로 짧아지지 않는다. 합성 가입(`KNOWSLINK_SYNTHETIC_SIGNUP=1`)은 agent ID를 고를 수 있어 삭제된 ID를 다시 쓸 수 있다. 공개 후보에서는 꺼야 한다. 실제 상태 크기·처리량·백업 복원 뒤 철회 유지는 OPS가 측정한다. 이번 DEV는 운영 데이터 정리를 실행하지 않았다.

## 미검증·범위

실메일·공개·운영 배포·외부 계정·운영 데이터 정리·실제 24h 경과·부하·실제 JSONB 크기와 CPU는 실행하지 않았다. 직접 시각 판정은 designer 담당이다. UX-06·일반 text·HTTP 동시 처리·gate rate는 MESSAGES/공개 수락 후속이다.
