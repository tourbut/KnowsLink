---
title: SAR-PUBLIC-AGENTS-001-POLICY — 거절·만료 후 재초대의 제품 규칙과 검증 조건을 확정한다
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-AGENTS-001-POLICY]
summary: 거절·만료 후 재초대의 제품 규칙과 검증 조건을 확정한다
---

# SAR-PUBLIC-AGENTS-001-POLICY — 거절·만료 후 재초대의 제품 규칙과 검증 조건을 확정한다

- 상태: ready. 작성일: 2026-10-06. From/To: coor/designer.
- 담당: repo818c78e5-d51c-4ff4-aa88-70e9ee185fbb, /home/shin/orca/workspaces/KnowsLink/fullops-designer, fullops/designer.
- 복귀: /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_6895aaf1-7b43-4fe0-a416-76f1255a5946, run_8ca8bc058ab7. Task/Dispatch는 preamble/영수증.
- 기준 ref: d1651784c4338efeb0d6141467d563c6b354e4a5. 동일 fixed제품 d1eef9b의 원본review/UI를 참고한다.

## 현재 상황과 적용 기준

원본 AGENTS 보안 리뷰의 L2 제품 해석만 판단한다. DEV-FIX는 M1 보존량/L1·L3 rate 및 UI medium2를 구현 중이다. 기술 계획·구현·QA는 DEV에 유지한다. 현재 designer UI기록은 완료됐으며 UX FAIL/보류와 metadata검사실패는 보존됐다. 제품 범위/고정수치/실제운영 승인 상태를 확대하지 않는다.

적용 기준은 fullops-common-0.3.3·FULLOPS·project·document-writing 및 최신 SAR-PUBLIC-SERVICE PS07/PS11와 UX05다. 다음을 먼저 읽는다.

- .fullops-squad/FULLOPS.md, project.md, rules/common/README.md 및 세 규칙, contexts/designer.md.
- .fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md.
- .fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md.
- .fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md의 L2 원문.
- .fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md의 실제관측(원래FAIL불변).
- .fullops-squad/docs/planning/sources/silent-agent-relay/protocol.md는 기존관계세대/철회 보안 근거이며 최신 일반서비스 제품결정과 충돌하면 우선순위를 설명한다.

Jev code/documents-find와 context는 docs/evaluations/jev/SAR-PUBLIC-AGENTS-001-POLICY-*.json이다. 추천 D09/D10은 기술 소유이므로 직접 수정하지 않는다. D02와 필요 UX05 문구만 명확하게 한다.

## 제품 질문 원문

L2: B-owner가 거절한 뒤 A가 같은 대상에게 다시 초대하면 Generation이 늘어난 새 pending이 생기며 B에는 차단 수단이 없다. 송신 owner pending10·rate40/60s로 제한된다. PS의 “반복 초대는 새 초대를 만들지 않는다”가 pending 중 반복만인지 거절 뒤 반복도 포함하는지 불명확하다.

## 해야 할 일과 완료 기준

- [ ] 기존 목표·규칙·남용방어 안에서 pending반복/거절/만료/양측철회 이후 재초대·재수락/세대에 대해 명확한 제품 판단을 한다. 새 blocking기능/제품수치/조건이 필요하면 실제 필요성과 범위를 설명하고 기존 요구와의 관계를 남긴다. 질문을 사용자에게 넘기는 대신 담당 제품판단을 완료한다.
- [ ] D02 정본과 필요 UX05에 짧은 규칙을 반영한다. 과거 구현이 새 규칙을 이미 지켰다고 기록하지 않는다. DEV/TESTER가 관찰 가능한 상태·허용/거부 조건과 미결정 항목을 작성한다. 기술 API/함수 구현 방식을 강제하지 않는다.
- [ ] docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md에 판단·원문/근거·변경·기술 인계를 남긴다. D02는 deliverables.py --stamp로 메타데이터를 보존하며 product-specs 원천을 갱신한다.
- [ ] git diff --check·strict·깨끗한 기록HEAD의 FullOps lint --from d165178을 완료하고 명령exit/미실행/경고를 적는다. 제품코드는 바꾸지 않는다.

소유권은 제품기획 D02/필요 UX정본과 자기인박스/contexts/로그/과제실행기록이다. PLANS는 자기결과만 추가한다. board/기술정본/코드/타역할기록을 변경하지 않는다. 실제메일/배포/외부계정/과금/서버/Tunnel 변경 없음. Workers Free 제한과 원본 UI FAIL·일반 서비스 후속 검증을 유지한다.

## 완료 보고

정확한 제품 답·D02 변경/검증·고정SHA·미결정/후속을 전문으로 남긴다. work.py finish로 인박스와 전문을 보존하고 빈인박스 확인뒤커밋한다. preamble의 worker_done을한번보내고idle로둔다. coor는 이답을DEV-FIX에그대로전달하고구현변경과최종후보독립검증을진행한다.
