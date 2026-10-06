---
title: SAR-PUBLIC-MESSAGES-001-UI-FIX — 수정 gate Deny 직후 결과·홈 복귀와 영향 UI 직접 재검수
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-MESSAGES-001-UI-FIX]
summary: 수정 gate Deny 직후 결과·홈 복귀와 영향 UI 직접 재검수
---

# SAR-PUBLIC-MESSAGES-001-UI-FIX — 수정 gate 결과·홈 복귀와 영향 UI 직접 재검수

- 상태: ready. fixed dfc70caa748a90614b02d48c78b4651345938339. 원본09c UX07 FAIL/UI기록f154·원OPS H1/M1/checkexit1은불변이다.
- 담당 designer / /home/shin/orca/workspaces/KnowsLink/fullops-designer / fullops/designer.
- 복귀 coor / term_6895aaf1-7b43-4fe0-a416-76f1255a5946 / run_8ca8bc058ab7 / repo818c78e5-d51c-4ff4-aa88-70e9ee185fbb. Task/Dispatch는preamble기준.

## 적용 기준과 먼저 읽을 문서

fullops-common-0.3.3 네규칙·FULLOPS/project/document-writing/contexts/designer·제품SAR-PUBLIC-SERVICE PS04/06/07/08–11·MVP C1–C5·UX06/07와 memberStyle 기준을읽는다. DEV-FIX phase/QA/완료로그·원UI보고/manifest/supplement와 UI-RECORD-REVIEW report를읽는다. 새capacity/cleanup_admission/http/store/member/member_receipt 관련영향과UI정본·Jev결과를읽는다.

## 해야 할 일과 완료 기준

- [ ] 정확한fixed를별도격리clone/자기fixture에서실제일반회원로그인·자기gate로직접검수한다. Deny클릭/keyboardEnter직후canonical결과200·denied저장·결정버튼비활성·홈복귀실제클릭을desktop1280/mobile390에서확인한다. owner쪽deny변경도좁게확인한다. 기존CSRF/권한/만료정책을바꾸지않는다.
- [ ] 원기록low R-UI-1은Deny405 mobile화면effectiveviewport980이다. 수정된실제390px을확인하고폭/긴ID/typedbody/keyboard/안전오류·원문부재·철회/만료 approve비활성을영향범위로재검수한다. HTTP입장수정이CSRF/429/cleanup결과화면에준영향을좁게확인한다. snapshot폭자동검사를직접시각PASS로삼지않는다.
- [ ] 변경없는UX06/의존성의diff0를확인하면원fixed09c의원래조건/근거만연결해재사용한다. 새실행PASS로바꾸지않는다. API전체권한/경합/rateQA는tester이며자기fixture화면PASS와구분한다.
- [ ] 지정시각항목PNG/manifest와실제봤던fixedSHA/조건/직접PASSFAIL·초기장애/원실패·마스킹/자기cleanup을보존한다. 영상은정지화면으로판정불가할때만. 보고서docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md 및같은키QA/phase/context/PLANS를작성하고필요D04연결만갱신한다. 원UI/OPS결과·제품규칙/코드/기술정본·타인인박스는변경하지않는다.
- [ ] 최종기록HEAD의FullOps --fromdfc70 lint/test·strict/diff exit/경고/HEAD·work.py finish전문/빈인박스·일반push뒤worker_done fullSHA/근거/미검증/후속을보고한다. 마지막JSON레포밖이면coor영속화경로를보고한다.

## 제약·검수 경계

theme/designlint/Tailwind/shadcn미구성이다. 기존Go template/memberStyle/UX정본을유지한다. 자유composer/채팅버블/장기timeline/queued성공오인표시를추가하지않는다. 직접시각검수만수행하고제품/기술판단은대신하지않는다. Workers Free·서버/Tunnel유지·유료전환/구독/초과과금/운영배포/실메일/외부계정/운영자료삭제금지. 자기격리로컬검수/자기자원회수/문서커밋·일반push는승인됐다. 실메일/PS08/13/14실운영·실24h/노우↔다닷·부하/복원·전체API/스크린리더는미검증으로유지한다. 실패는그대로보고하고정상원본을재작성하지않는다.

## 지시 전제와 충돌 — 먼저 확인

원UI보고의UX07 FAIL은09c 대상이며원본그대로유지한다. 이번dfc 후보의DEV자동통과는직접시각수락이아니므로실제새화면으로재검수한다. 원FAIL/모바일viewport980를새PASS로덮어쓰지않는다. sensitive/oversized code passage는Jev미전송이므로필수keep으로직접읽는다.
