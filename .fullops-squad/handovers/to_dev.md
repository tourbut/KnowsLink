---
title: SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX — 확정된 재초대·기록 포화 제품 답을 구현과 검증에 반영한다
status: draft
updated: 2026-10-06
owner: dev
tasks: [SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX]
summary: 확정된 재초대·기록 포화 제품 답을 구현과 검증에 반영한다
---

# SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX — 확정된 재초대·기록 포화 제품 답을 구현과 검증에 반영한다

- 상태: ready. 작성일: 2026-10-06. From/To: coor/dev.
- 담당: repo818c78e5-d51c-4ff4-aa88-70e9ee185fbb, /home/shin/orca/workspaces/KnowsLink/fullops-dev, fullops/dev.
- 복귀: /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_6895aaf1-7b43-4fe0-a416-76f1255a5946, run_8ca8bc058ab7. Task/Dispatch는새preamble/영수증이다.
- 기준 ref: 4a1b80aec8fa6a06144d51f3a5609927a2644928. 이전 제품85fb40e와원래d1/UI FAIL을보존한다. 제품답은48d12fae2dce35d92606b264313148f0a635b64e다.

## 현재 상황과 적용 기준

기술 M1/L1/L3·F-UI01–04 수정후보4a는 clean·lint/test/verify-mvp/strict0을보고했다. 아직독립delta검수는없다. designer POLICY48은D02/UX04–05와제품답을확정했으나새동작PASS를주장하지않았다. 원본d1 QA는진행중이며원본UI FAIL·보안리뷰와검사실패기록을유지한다.

fullops-common-0.3.3 README·코딩/테스트/보안 규칙·FULLOPS·project·document-writing을적용한다. 제품 책임은designer, 기술분석·계획·구현·검증은이DEV과제다. 새판단이나상품quota/숫자를추가하지않는다.

## 먼저 읽을 문서

Jev code/documents-find 및context는 docs/evaluations/jev/SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX-*.json이다. 충돌이있으면제품원천과기술원천의책임을구분해coor로질문한다.

- .fullops-squad/FULLOPS.md, project.md, contexts/dev.md, rules/common/README.md 및세규칙.
- .fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md의제품답원문과두관찰조건표전체. coor가줄이거나해석하지않으며이원문을그대로전달한다.
- .fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md의PS07-I·기록보호조건, docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md의UX04–05.
- .fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md.
- internal/relay/http.go, member_agents.go, member.go, connections.go, connections_integration_test.go 및필요한호출/테스트.

## 해야 할 일과 완료 기준

- [ ] 48의관찰조건과4a의실제구현차이를확인하고짧은기술계획·예상규모를남긴다. 이미일치하는기능을재구현하지않는다.
- [ ] 유효pending 같은/반대방향및active반복의수/세대/기한불변과권한유지를맞춘다. 종료후수동재초대는새세대·새수락이며옛결정/승인/메시지권한과자동복구를차단한다. 모든기존한도/rate/owner인가와교차owner·same-owner를유지한다.
- [ ] 기록보호상한의안전한오류와가능한다음동작/관리복귀를구현한다. 키기록포화는새agent별도연결·각상대새관계수락을안내한다. owner기록포화는필요철회와최소24h보존/실제정리뒤수동재시도를안내한다. 키철회/대기만으로살아있는agent키공간이생긴다거나24h에정확히정리된다거나결제가필요하다고안내하지않는다.
- [ ] 철회목록의최소24h보존/실제정리후홈에서사라질수있음과권한복구/백업영구삭제아님을안내한다. 보존중철회상태와정리뒤옛key/pair거부를검증한다. UX단순정보구조를유지하며새frontend/장식/대화composer를추가하지않는다.
- [ ] 최신변경의결정적회귀와make lint/test/verify-mvp·최종깨끗한HEAD FullOps --from4a·git diff --check·strict를완료한다. 명령자체exit/실패원문과재실행이유를보존한다. 원본검수와4a검증을새SHA실행으로재작성하지않는다.
- [ ] 기술원천D10/실제영향D05/D09/README와실행기록 docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX.md를갱신한다. Jev D12/D13 추천은OPS소유라구현기술인계/미검증으로제공하고운영수락을쓰지않는다. 자기contexts·완료로그·PLANS결과만추가한다.

## 검증·UI 인계와 제약

최종후보는OPS의독립delta리뷰·TESTER의보존/rate/제품조건좁은QA·designer의F-UI01–04/반복초대/기록포화복구안내직접재검수뒤수락한다. 안정후보에서지정담당이한번검사하며기존증거는변경의존성동일성확인시원래SHA로만재사용한다. Go memberStyle/template을재사용하고별도theme/Tailwind/shadcn/디자인lint없음의근거와실제변경영향을남긴다. SIZE/DEP/DESIGN경고는실제규모/이유를보고한다.

Workers Free·서버/Tunnel을유지한다. 실제메일/공개/배포/외부계정/운영데이터정리/과금은없다. 격리localfixture·코드검증·커밋·일반역할push는허가됐다. 제품규칙변경/범위밖행위/설명되지않는실패만묻고나머지는완료까지진행한다. board/제품정본/타역할결과는수정하지않는다.

## 완료 보고

고정SHA·제품원문준수/구현차이·검증/미검증·경고/의존성·기술문서/QA UI후속을전문으로쓴다. work.py finish로정규인박스와결과를아카이브하고빈인박스확인후커밋한다. 새preamble의Task/Dispatch/worker_done을한번사용하고idle로둔다.
