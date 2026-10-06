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
