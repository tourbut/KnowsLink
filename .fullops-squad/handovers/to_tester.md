---
title: SAR-PUBLIC-AGENTS-001-FIX-TESTER — 최신 AGENTS 후보 보존·rate·제품 조건 독립 좁은 QA
status: draft
updated: 2026-10-06
owner: tester
tasks: [SAR-PUBLIC-AGENTS-001-FIX-TESTER]
summary: 최신 AGENTS 후보 보존·rate·제품 조건 독립 좁은 QA
---

# SAR-PUBLIC-AGENTS-001-FIX-TESTER — 최신 AGENTS 후보 보존·rate·제품 조건 독립 좁은 QA

- 작성일: 2026-10-06
- From / To: coor / tester
- 상태: ready
- 승인된 범위: 로컬 비파괴 검수·격리 합성 fixture·필요한 기록 작성. 실제 메일·공개·배포·외부 계정·운영 데이터 정리·과금 금지.
- 담당 repo id: 818c78e5-d51c-4ff4-aa88-70e9ee185fbb; 워크트리 /home/shin/orca/workspaces/KnowsLink/fullops-tester, 브랜치 fullops/tester.
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

인박스는 현재 to_tester.md만 사용한다. 역할 기록만 수정하고 제품 코드는 수정하지 않는다. PLANS에 자기 결과를 append하며 원본을 삭제하지 않는다. work.py finish로 실제 지시서/전문을 보존·빈 인박스·커밋·공유 role push 뒤 전체 40자리 SHA로 worker_done 한 번 보낸다. 완료와 main 제품 수락은 구분한다. 실패도 원문·실행 exit·재실행 사유를 보존한다. 예정 규모 보고서/시나리오/필요 증거뿐이며 SIZE/DEP 경고를 설명한다. 이미지·로그는 비밀값 없이 판정에 필요한 범위만 만든다.
## 해야 할 일과 파일 소유권

- [ ] 별도 fixed 458798c2ee15c179edacfd6f94ebb9896d26f411 clone의 격리 Postgres/Compose fixture로 변경 영향만 독립 검증한다. 원본 d1 QA bcb06b8 시나리오/runner/실패기대값 정정 근거를 확인하고 의존성이 같은 항목만 원래 SHA로 재사용한다.
- [ ] POLICY 두 표: 같은/반대 방향 pending·active 반복의 수/세대/기한 불변·한도에서 중복 반복, 거절/만료/철회 뒤 수동 새세대/새수락, 양측철회·구자격·옛결정의 실패·경합·restart를 확인한다.
- [ ] M1 최소24h 이전 보존/실제 정리/legacy 최초철회시각·살아 있는 agent revoked kid 재사용금지, owner/key 기록 포화·기존자격 유지·철회가능·새agent/새pair 독립연결·24h cleanup/manual retry를 검증한다. 시간이동fixture는 실제24h wait와 구분한다.
- [ ] 거부 transaction의 부분변경 rollback과 rate 집계, invalid session GET anonymous rate, connect429 retry_at/Retry-After 실제 deadline 일치 및 restart 관련 좁은 회귀를 수행한다. UI 새 안내와 상태코드/next action 존재는 검사하고 직접 시각 판정은 designer에게 둔다.
- [ ] 새 시나리오와 report는 docs/evaluations/scenarios/SAR-PUBLIC-AGENTS-001-FIX-TESTER.md, docs/evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-FIX-TESTER.md. executable probe·원래exit·fixedhead·fixturecleanup/공유container불변·마스킹을 증거로 남긴다.

## 완료 기준과 산출물

관찰 조건 PASS 또는 정확한 FAIL/held 재현, critical/high0, 검사 본래 exit0 및 근거. clean committed 기록 HEAD에서 FullOps --from 458798c2ee15c179edacfd6f94ebb9896d26f411의 lint/test·문서strict·diffcheck를 한 번 실행한다. 최종 archive 후 head증거 누락은 명시한다. 전체 기존 QA 반복 금지; 변경 영향·새 실패·근거 결함에 필요한 검사만 한다. D10/D11 QA 소유 원천을 갱신한다. D12 운영 수락은 OPS 후속이며 실제메일/공개/운영자원/노우↔다닷 미검증. UI 테마 없음·캡처 불필요.
