---
title: SAR-PUBLIC-AGENTS-001-DEV-FIX — AGENTS 보존량·rate 발견의 원인 수정과 관련 검증
status: draft
updated: 2026-10-06
owner: dev
tasks: [SAR-PUBLIC-AGENTS-001-DEV-FIX]
summary: AGENTS 보존량·rate 발견의 원인 수정과 관련 검증
---

# SAR-PUBLIC-AGENTS-001-DEV-FIX — AGENTS 보존량·rate 발견의 원인 수정과 관련 검증

- 작성일: 2026-10-06
- From / To: coor / dev
- 상태: ready
- 담당: repo 818c78e5-d51c-4ff4-aa88-70e9ee185fbb, /home/shin/orca/workspaces/KnowsLink/fullops-dev, fullops/dev
- 복귀: /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_6895aaf1-7b43-4fe0-a416-76f1255a5946, run_8ca8bc058ab7. Task/Dispatch는 영수증과 preamble이다.
- 기준 ref: d1eef9bb90b9726149980320c42fb1fdbcaf584a. main/origin dc60fbf의 조상 관계를 준비 커밋에서 확인한다.

## 현재 상황과 적용 기준

원본 d1의 독립 리뷰70f26bc(별도actual50b08fc6 vs DEV01a10f52)는 critical/high0, check/lint/test/verify-mvp0다. M1 보존량은 공개 전 차단이다. 이번 일반 서비스 기능의 후속에서 M1 및 L1/L3을 해결한다. 원본 QA/UI는 현재 같은 fixed d1에서 진행 중이므로 결과를 변경하지 않는다. 수정 후 별도 delta 리뷰와 관련 좁은 QA/UI 영향 검수를 진행한다.

fullops-common-0.3.3의 README·코딩/테스트/보안 규칙, FULLOPS·project·document-writing, 제품 SAR-PUBLIC-SERVICE PS04–07/PS11/개인정보와 UX04–05를 적용한다. 기술 분석·계획·구현·회귀를 같은 DEV 과제에서 처리한다. 제품 수치를 새로 결정하거나 기존 보존 최소24h·권한 차단을 낮추지 않는다.

## 먼저 읽을 문서

Jev code/documents-find와 context는 docs/evaluations/jev/SAR-PUBLIC-AGENTS-001-DEV-FIX-*.json이다. 다음 keep 후보와 필수 정본에서 시작한다. 민감값은 외부조회/로그에 넣지 않는다.

- FULLOPS.md, project.md, contexts/dev.md, rules/common/README.md 및 세 규칙(.fullops-squad 기준).
- .fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md의 M1/L1/L3 원문과 result/lint.
- .fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md.
- .fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md.
- internal/relay/connections.go, member_agents.go, member.go, api_rate.go, http.go, store.go와 관련 호출/테스트.

## 해야 할 일과 완료 기준

- [ ] 짧은 기술 계획과 예상 변경 규모를 실행 기록에 남긴다. M1의 신규/철회 agent·key·pair 누적과 nil 역참조/옛 세대 재사용·receipt/lease/claim 의존 경계를 확인한다.
- [ ] 기존 최소24h 보존과 권한 철회 안전성을 지키며 만료 철회 기록의 정리 및 보존 기록 포화 보호를 구현한다. 시간 동안 쌓이는 상한·한도 포화 중 신규 수락 차단/안전 정리·동시성과 재시작을 검증한다. 보호값은 기존 제품 quota를 바꾸지 않는 기술 자원 보호로 근거를 기록한다. OPS의 실제 CPU/DB/복원 검증은 후속이며 성공으로 표시하지 않는다. 운영 데이터에는 정리를 실행하지 않는다.
- [ ] L1의 무효 세션 GET 연결 화면과 동형 home GET에서 rate 기록 rollback 원인을 고친다. 실제 실패/거부 요청이 budget을 소비하도록 관련 호출을 함께 확인한다.
- [ ] L3의 connect 429에서 기존 rateAPI와 일관된 실제 retry_at/Retry-After를 제공하고 CLI/안내 영향과 경계를 검증한다.
- [ ] L2 거절/만료 후 재초대 규칙은 designer 결정 대기다. 임의 변경하지 않는다. 답이 오면 같은 과제 후속으로 처리하도록 기록한다.
- [ ] 결정적 회귀·make lint/test/verify-mvp·최종 깨끗한 고정 HEAD FullOps lint --from d1eef9b·git diff --check·문서 strict를 완료한다. 자체 exit와 실패/재실행 이유를 보존한다. 원본 증거를 덮어쓰지 않는다.

제품/테스트/기술 문서 및 자기 인박스/contexts/완료로그를 소유한다. PLANS는 결과만 추가하고 board/제품 규칙/다른 역할 결과는 수정하지 않는다. D10 및 실제 영향 D05/D06/D07/D09 원천을 갱신하고 D12 운영 경계는 구현 기술 인계로 제공한다. 로그는 docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md다.

## 검증과 UI 후속

DEV 자동 검증 뒤 최신 후보의 독립 OPS delta 리뷰·TESTER 보존/철회/rate 좁은 QA가 필요하다. UI 변화가 있으면 designer에게 해당 변화만 인계하고 기존 fixed d1 시각 결과는 원래 SHA로 재사용한다. 기존 Go template/memberStyle을 쓰고 별도 theme/Tailwind/shadcn/디자인 lint 없음 근거를 남긴다. DESIGN/SIZE/DEP 경고의 실제 영향과 예외를 보고한다. 영상이나 전체 시각 검사를 중복하지 않는다.

## 제약과 완료 보고

Workers Free·기존 서버/Tunnel을 유지한다. 유료 설정/실메일/운영 데이터 정리/공개/배포/외부 계정 연결은 실행하지 않는다. 격리 로컬 fixture와 비파괴 코드 구현·검증·커밋·일반 역할push는 허가됐다. 제품 규칙 변경·범위 밖 행위·설명되지 않는 검증 실패만 묻고 나머지는 완료까지 수행한다.
완료 보고에 고정 SHA·원인/수정·검증/미검증·규모/의존성·산출물·QA/UI 영향과 후속을 적는다. work.py finish로 전문 보존과 빈 인박스를 확인해 커밋하고 preamble의 worker_done을 한 번 보낸다.

## 추가 지시 (coor msg_6a55391817b0, 2026-10-06)

- [x] F-UI-01 medium: 390×844 홈의 지문 넘침(scrollWidth 596)을 고친다.
- [x] F-UI-02 medium: 재확인 필요·미지원·권한 거부 결과에 홈·재확인·재선택 동작을 제공한다. 안전 거부는 유지한다.
- [x] F-UI-03 low·F-UI-04 low: 완료 상태의 취소 버튼, 관계 raw UTC 기한·select 가독성을 기술 판단으로 처리한다.
- 원천은 designer의 `SAR-PUBLIC-AGENTS-001-UI.md`(f9c5eded)이며 읽기만 한다. 기존 시각 FAIL/PNG는 보존한다. 변경 UI는 designer 좁은 재검수로 인계한다.

## 착수 기술 계획

기준 d1eef9b·준비 b5df23b·main dc60fbf의 조상 관계를 확인했다. sweep·transaction·memberStyle을 재사용하고 새 의존성·migration·wire 변경은 없다. M1은 C1 kid 재할당 금지를 지키기 위해 철회 agent 전체만 24h 뒤 삭제하고, 신규 기록을 owner·agent 상한으로 막는다. L1은 transaction의 rate 저장과 세션 helper 하나로 고친다. L3은 /v1 rate 응답을 공용화한다. 예상 규모는 Go 6파일·테스트 2파일·기술 문서 6–8개였다.

## 완료 보고

SAR-PUBLIC-AGENTS-001-DEV-FIX의 분석·계획·구현·회귀·기술 문서를 완료했다. 제품 코드 고정 SHA는 `85fb40e643aa1f08e1c7b2c41eac199853517cba`다(보존 수정 `1f261b3`, UI 수정 `85fb40e`). 기록 HEAD는 work.py finish 커밋과 worker_done에 고정한다.

원인과 수정: M1은 sweep이 Connection만 지운 것이 원인이었다. 철회 agent는 철회 시각(Agent.Changed)부터 24h 뒤 키·pair와 함께 삭제한다. 살아 있는 agent의 철회 키는 C1 kid 재할당 금지 때문에 지우지 않는다. owner agent 기록 10·agent 키 기록 20의 기술 상한은 신규 생성·연결만 409로 거부하고 철회는 허용한다. 없는 agent를 가리키는 pair는 nil 역참조 없이 거부한다. L1은 transaction rollback이 rate 기록까지 버린 것이 원인이었다. 거부 요청의 rate를 저장하고 `memberHit`이 home·연결 화면·reauth·logout·회원 POST의 무효 세션을 익명 budget에 집계한다. L3은 connect 429가 /v1과 같은 실제 retry_at·Retry-After(초 올림)를 준다. F-UI-01–04는 지문 `<code>`·줄바꿈·select 스타일, 거부 화면의 홈/로그인 링크와 재확인 버튼, 열린 연결만 취소, KST 관계 기한으로 고쳤다.

검증: 85fb40e에서 make lint 0·make test 0·make verify-mvp 0(Go PASS 53·FAIL 0)이다. 321881a에서 FullOps lint --from d1eef9b exit 0(ERROR 0·WARNING 7·product-lint/test passed), deliverables --strict 0(13/0/0/0), git diff --cached --check 0이다. 390×844 렌더 측정은 같은 홈 상태의 scrollWidth를 d1 595에서 390으로 줄였다. 새 검사는 TestRevokedRecordRetention·TestMemberPagesShowNextSteps와 integration 세 하위 검사다. 실패 뒤 재실행 사례: 첫 단위 검사는 연결 TTL 10분을 넘는 테스트 시각 때문에 실패했고 시각 간격을 줄여 통과했다. 제품 결함이 아닌 테스트 설정 오류다.

경고 처리: SIZE-001 6건은 기존 대형 파일(PLANS·http·member·store)의 소폭 증가와 회귀 테스트 증가다. SIZE-002 추가 598줄은 제품 425줄(테스트 267 포함)과 기술 문서다. 보존·rate·화면이 같은 회원 연결 경계라 같은 후보로 검증하려고 분할하지 않았다. DEP 변경 0, DESIGN 경고 0이다. DESIGN은 Go 문자열 CSS를 검사하지 않으므로 폭 측정을 남겼다.

미검증·후속: 실제 24h 경과·부하·JSONB 크기·CPU·복원·운영 데이터 정리·공개·실메일·외부 계정은 실행하지 않았다. OPS delta 리뷰, TESTER 보존/철회/rate 좁은 QA, designer 변경 화면 좁은 재검수가 필요하다. 키 기록 포화 시 agent 교체 안내와 24h 뒤 철회 agent가 홈에서 사라지는 동작은 designer 제품 판단으로 인계한다. L2는 designer 결정 대기이며 코드를 바꾸지 않았다. `GET /`·`/auth/verify`·`/home/gates/*`·`/owner/*`의 HTTP rate 미적용은 MESSAGES 후속 후보다. 상세는 `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md`, 증거는 `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-DEV-FIX/`다.
