---
title: SAR-PUBLIC-AGENTS-001-UI-FIX — AGENTS 수정 화면과 제품 답 직접 재검수
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-AGENTS-001-UI-FIX]
summary: AGENTS 수정 화면과 제품 답 직접 재검수
---

# SAR-PUBLIC-AGENTS-001-UI-FIX — AGENTS 수정 화면과 제품 답 직접 재검수

- 작성일: 2026-10-06
- From / To: coor / designer
- 상태: ready
- 승인된 범위: 로컬 비파괴 검수·격리 합성 fixture·필요한 기록 작성. 실제 메일·공개·배포·외부 계정·운영 데이터 정리·과금 금지.
- 담당 repo id: 818c78e5-d51c-4ff4-aa88-70e9ee185fbb; 워크트리 /home/shin/orca/workspaces/KnowsLink/fullops-designer, 브랜치 fullops/designer.
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀: coor /home/shin/orca/workspaces/KnowsLink/fullops-coor / term_6895aaf1-7b43-4fe0-a416-76f1255a5946 / run_8ca8bc058ab7. task/dispatch는 preamble 실제 값.

## 현재 상황과 확인 근거

고정 검수 후보 `458798c2ee15c179edacfd6f94ebb9896d26f411`. 원본 DEV d1eef9b, 원본 OPS 70f26bc, 원본 QA bcb06b8, 원본 UI d165178 FAIL, DEV-FIX 4a1b80a, 제품 답48d12fa, DEV-POLICY-FIX83e0bfb를 모두 포함한다. 제품 코드 최종6d016e5와 후보의 제품 동일성을 확인한다. 원본 실패·held·실행 SHA를 보존한다. 원본 QA/UI를 최신 SHA 실행으로 바꾸지 않는다.

## 적용 기준과 예외

fullops-common-0.3.3의 README/coding-style/testing/security, FULLOPS, project, document-writing, review/rule.json을 직접 읽는다. 제품 정본 D02 PS07·PS07-I 및 기록 보호 조건과 UX04–05, POLICY48의 두 관찰표를 적용한다. 기준은 위 fixed 후보, 예외 없음. Workers Free·기존 서버/Compose/Tunnel 유지. ponytail full. 기술 수정은 DEV, 제품 규칙 변경은 designer로 coor 경유 질문한다.

## 먼저 읽을 문서

- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md`
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md`
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` 전체 두 관찰표
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md`
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX.md`
- 원본 OPS report, TESTER report/probe-failures, UI 보고서의 관련 항목. Jev 추천은 탐색 보조이며 제외 권한이 아니다.

## 제약·완료 처리

인박스는 현재 to_designer.md만 사용한다. 역할 기록만 수정하고 제품 코드는 수정하지 않는다. PLANS에 자기 결과를 append하며 원본을 삭제하지 않는다. work.py finish로 실제 지시서/전문을 보존·빈 인박스·커밋·공유 role push 뒤 전체 40자리 SHA로 worker_done 한 번 보낸다. 완료와 main 제품 수락은 구분한다. 실패도 원문·실행 exit·재실행 사유를 보존한다. 예정 규모 보고서/시나리오/필요 증거뿐이며 SIZE/DEP 경고를 설명한다. 이미지·로그는 비밀값 없이 판정에 필요한 범위만 만든다.
## 해야 할 일과 파일 소유권

- [ ] 고정 458798c2ee15c179edacfd6f94ebb9896d26f411 자신의 격리 합성 relay/DB/smtp fixture·자신이 새로 만든 browser tab으로 실제 화면을 직접 본다. 기존 사용자 브라우저/공유 컨테이너는 건드리지 않는다.
- [ ] 원본 d165178 UI F-UI-01 mobile390 fingerprint 넘침, F-UI-02 refusal/reverify/unsupported/forbidden 복귀·다음행동, F-UI-03 종료 연결취소 버튼 부재, F-UI-04 KST deadline/select가독성을 desktop1280/mobile390에서 좁게 다시검수한다. 원본FAIL/PNG/blank Orca 실패는 보존한다.
- [ ] POLICY48 두 표와 새 안내: pending/active 반복notice와 받은/보낸결정위치, 종료관계 수동 새초대/새수락, 키 포화시 live revoked kid 공간미복구·새agent/새pair, owner기록포화 최소24h/실제cleanup/manualretry, revokedagent홈에서 사라질수있음·복구/영구삭제보장없음, revoke전경고를 직접 검수한다.
- [ ] 각 필수항목에 실제봤던 fixedSHA·fixture조건·PNG와 PASS/FAIL을 연결한다. 영상은 필요없다. bearer/session/code/private key를 이미지 전 마스킹하고 이미지들을 직접 열어 판정한다. 공개fingerprint는식별내용이므로보존. Orca screenshot blank/1px면 원래실패증거보존 후 기존처럼 설치된 Playwright/Chromium fallback을 사용한다.
- [ ] 보고서/PNG manifest/cleanup은 docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md 및 sibling folder, qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX에 남긴다. 기존 unchangedUI는 의존성같을때 d1 원래실행으로재사용한다. 기술수정 필요는 coor/DEV로 정확한항목 보고한다.

## 디자인 기준·완료 기준·산출물

D02/UX04–05/POLICY48이정본. memberStyle/Go template 재사용, 새theme·Tailwind·shadcn·공용theme전환·디자인lint없음. DESIGN자동검사미지원 영향과 직접검수근거를보고한다. 제품기획판단만담당하고코드수정금지. 필수UI FAIL해소및새안내제품조건준수여야수락가능. D03 UX원천/직접검수기록을갱신하되원본FAIL을PASS로덮지않는다. clean 기록HEAD FullOps --from 458798c2ee15c179edacfd6f94ebb9896d26f411·strict·diffcheck와원래exit를보존. 실제메일/공개/실24h운영/노우↔다닷미검증.
