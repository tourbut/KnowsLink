---
title: SAR-PUBLIC-MESSAGES-001-UI-FIX-2 — 최종 수정 후보의 정리 흐름과 오류·복귀 화면 직접 확인
status: draft
updated: 2026-10-07
owner: designer
tasks: [SAR-PUBLIC-MESSAGES-001-UI-FIX-2]
summary: 최종 수정 후보의 정리 흐름과 오류·복귀 화면 직접 확인
attempt: 57c4c753a28e4eb5ba1b5db794c2b545
base: 68b0d6a0c854fdaec6828a232dd3945814be1404
---

# SAR-PUBLIC-MESSAGES-001-UI-FIX-2 — 최종 수정 후보의 정리 흐름과 오류·복귀 화면 직접 확인

## 대상·기준·복귀

- 상태 blocked. DEV-FIX-3 최종 SHA와 안정된 검수 후보 확정 전 착수하지 않는다. 실제 후보·독립 세션·snapshot 경로는 coordinator가 dispatch 전에 고정한다. 이 과제는 이전에 배정하지 않았으며 기존 UI-FIX 완료 기록과 다르다.
- 담당 designer. 기존 route의 Codex gpt-6.1-sol high를 사용한다. 사용자 소유의 기존 터미널·모델 선택 화면은 변경하지 않고 fresh 세션에서 수행한다.
- 공통 fullops-common-0.3.3 README/coding-style/testing/security, FULLOPS·project·document-writing·orca-agents·contexts/designer와 SAR-PUBLIC-SERVICE.md PS11 및 mockups/SAR-PUBLIC-SERVICE-UX.md UX06/07을 따른다. 기본 main 68b0d6a0c854fdaec6828a232dd3945814be1404와 고정 후보의 같은 문서를 읽는다.
- 복귀 coor /home/shin/orca/workspaces/KnowsLink/fullops-coor, terminal term_a8a1fa04-50ab-448d-94e7-11e8ee3c77f1, Run run_8ca8bc058ab7. 실제 Task/Dispatch 권한은 새 preamble을 따른다.

## 목적과 소유권

최종 후보의 변경 영향에 필요한 직접 UI 검수와 제품 계약 적합성을 확인한다. 검수 보고서·자기 증거·D04 연결·PLANS/context/inbox만 수정한다. 구현과 기술 검증은 DEV/독립 tester/별도 reviewer 담당이다. 제품 조건과 코드가 어긋나면 근거와 함께 coordinator에게 질문하고 기존 명세를 임의 완화하지 않는다.

## 해야 할 일과 완료 기준

- [ ] 원 UI09c FAIL 및 UI-FIX dfc의 좁은 PASS/GET rate 제한 관측을 원 SHA·조건·무결성으로 보존한다. 코드/템플릿/에셋/설정의 변경 영향과 의존성 동일성을 먼저 확인하고 재사용 여부를 판단한다. 원 증거를 새 실행으로 표시하지 않는다.
- [ ] 새 fixed에서 정상 회원의 Deny 클릭/Enter 결과와 홈 복귀를 실제 브라우저의 desktop1280/mobile390에서 직접 확인한다. 유효 정리가 성공한 상태와 신규/정리 rate 제한의 오류 안내·재시도 조건을 구분한다. 포화 GET에 즉시200을 약속하지 않는다. 변경된 자기 철회/unpair/logout 흐름의 화면 영향도 필요한 최소 항목으로 직접 확인한다.
- [ ] 새 기술 보호가 PS11의 안전 정리·기존 권한·제품 한도와 일치하는지 명세·DEV 근거와 UI 관측으로 판단한다. 회귀나 불명확한 제품 조건을 숨기지 않고 담당과 재개 조건을 보고한다.
- [ ] 실제 viewport/overflow·문구·버튼/링크·키보드 조작·상태/Location을 evidence manifest와 지정 PNG에 연결한다. 자동 판정 가능한 값은 코드로 확인하고 시각 판정은 PNG를 직접 읽는다. 정지 화면으로 충분하므로 영상은 만들지 않는다.
- [ ] 마지막 기록 SHA의 FullOps lint/test·strict·diff와 원본 hash/자기 fixture·프로세스 회수를 확인한다. 지시서 완료 보고 전문·work.py finish archive/빈 inbox·일반 push·실제 worker_done을 완료한다.

## 디자인 기준·탐색·산출물

기존 UX 정본과 member HTML/CSS를 그대로 사용한다. 새 색·글꼴·알약 버튼·카드 배치·아이콘·테마·애니메이션을 추가하지 않는다. 이번 과제는 UI 코드 변경이 없으며 테마 전환/새 디자인 lint 도구는 해당 없음과 이유를 기록한다. 기존 DESIGN 경고는 실제 후보 결과로 판단한다.

먼저 읽을 문서는 공통 기준/제품·UX 정본·DEV-FIX-3 기록/최종 검사·원 UI/직전 UI-FIX 보고서/manifest·TESTER 보고서다. 최종 fixed 이후 code/doc/context/packet을 연결한다. D04는 이번 검수 보고서 연결만 갱신한다. 기대 산출물은 이 과제 QA 보고서·manifest·지정 직접 시각 증거·실제 세션/snapshot/검증/회수 provenance다.

## 제약과 후속

Workers Free·기존 서버/Tunnel을 유지한다. 유료 전환·실메일·운영 공개·실제 외부 계정·사용자 자료 삭제는 금지다. 자기 격리 fixture만 만들고 회수한다. 제품 코드는 수정하지 않는다. 실메일/운영 공개/실24h/노우↔다닷/부하·복원은 미검증으로 구분한다. 기존 플랫폼 차단 출력을 재생성하지 않는다. 필요한 코드 수정은 coordinator를 통해 같은 DEV 후속으로 인계한다. 완료 뒤 idle이며 다른 작업을 시작하지 않는다.

## 완료 보고

작업 종료 뒤 실제 결과 전문을 작성한다.
