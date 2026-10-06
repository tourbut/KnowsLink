---
title: SAR-PUBLIC-MESSAGES-001-TESTER — 일반 회원 메시지 후보 독립 동작 QA
status: draft
updated: 2026-10-06
owner: coor
tasks: [SAR-PUBLIC-MESSAGES-001-TESTER]
summary: 일반 회원 메시지 fixed 후보 독립 QA 인계
---

# SAR-PUBLIC-MESSAGES-001-TESTER — 일반 회원 메시지 후보 독립 동작 QA

- 상태: ready. fixed 09c523da8a3407288d9f5d711e1834af12bc7808
- 담당: tester / /home/shin/orca/workspaces/KnowsLink/fullops-tester / fullops/tester
- 복귀: repo 818c78e5-d51c-4ff4-aa88-70e9ee185fbb / fullops-coor / term_6895aaf1-7b43-4fe0-a416-76f1255a5946 / run_8ca8bc058ab7. Task/Dispatch는 preamble 기준.

## 적용 기준과 먼저 읽을 문서

fullops-common-0.3.3의 README·coding-style/testing/security, FULLOPS·project·document-writing·contexts/tester를 같은 전달 SHA에서 읽는다. 제품 기준은 docs/planning/product-specs/SAR-PUBLIC-SERVICE.md PS08–11·PS04/06/07와 SAR-MVP.md C1–C5 및 docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md UX06/07다. DEV 실행 기록 docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV.md와 logs/2026-10-06_to_dev.md의 최신 과제·DEV QA 증거·COOR/dev-final-37f9a1e-lint.json을 읽는다. Jev code/doc/context 결과도 이 과제 키로 준비한다.

## 해야 할 일과 완료 기준

- [ ] fullops-test를 적용한다. fixed 09c523d을 별도 clone/자기 격리 자원에서 검증한다. 회원 두 agent·별도 key/credential·명시적 활성관계의 text 송신/수신/관련회신ID를 실제 로컬 Node/MCP 프로세스로 대조한다. 시험 allowlist/공유ServiceAuth를 일반 자격으로 쓰지 않는다. 로컬 fixture 결과는 전체 운영 PS08 PASS나 실제 노우↔다닷으로 보고하지 않는다.
- [ ] PS09 TTL180s/4096UTF8bytes·queued/leased/delivered/처리상태·수동pull·원문ACK/만료삭제·멱등동일key내용/충돌·관련회신1회/consent·키철회/관계세대/타회원거부·offline/expired/rate/revoked/invalid-key 실패를 독립 검사한다. PS10 verified typed-body·정책·CSRF/XSS/권한/gate approve/deny와 누락원문/만료/철회 거부, schedule.query 기본deny·commit불가도 확인한다.
- [ ] PS11 raw/queue/gate/receipt/HTTP신규16/정리4/claim4의 이하·경계·초과·동시성·두relay 공유/재시작/crash회수/high priority를 검사한다. 신규포화 중 유효ACK/deny/철회/unpair/로그아웃이 별도budget으로 가능한지 확인한다. 변경 없는 신원/AGENTS 증거는 원래SHA/조건과 동일성을 연결해 재사용한다.
- [ ] DEV 자체검사만으로 독립 QA를 대체하지 않는다. 새 시나리오·독립 재현 프로브·원본 stdout/exit·판정/초기실패를 보존한다. make lint/test/verify-mvp와 영향 자동검사를 담당 범위에서 실행한다. 제품결함은 coor경유DEV에 재현으로 보고하며 코드/기대조건을 바꿔 숨기지 않는다.
- [ ] 보고서 docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-TESTER.md와 시나리오/증거를 작성한다. 최종 기록HEAD FullOps --from09c523d lint/test·strict·diff의 exit/HEAD/ERROR/WARNING/실행불가를 남긴다. work.py finish로 전문아카이브/빈인박스·최종커밋·일반push 뒤 fullSHA/증거/실패·미검증·후속을 worker_done으로 보고한다.

## 소유권·제약·후속

소유는 QA 시나리오/프로브/증거·자기 실행 기록/context/PLANS·인박스/완료로그다. 제품코드·기획/UX정본·frozen원천·타 역할 인박스·기존리뷰는 수정하지 않는다. 필수실패/critical/high는 main 수락 차단이다. 실메일·공개·실24h·운영부하/복원·외부플랫폼/최종노우↔다닷은 별도미검증으로 기록한다. Workers Free만허용하고 유료전환/구독/초과과금·운영배포/실메일/외부계정/공유서버/Tunnel/운영자료삭제는 수행하지 않는다. 자기격리fixture/자원만 사용·회수한다. 문서/검사커밋·일반push는 승인됐다. OPS독립리뷰/designer직접UI는 별도이다.

## 갱신할 산출물과 완료 보고

원천QA/시나리오·D09/D10 관련 연결만 실제 영향에 따라 기록한다. final fullSHA·고정대상·검증/실패/미검증·명령exit/경고·증거·후속담당/재개조건을 전문으로 쓴다. front matter는 deliverables stamp를 쓴다.

추가 필수 확인: internal/relay/member_receipt.go가 실제 경로다. context의 member_receipts.go 거부는 오기이며 원본 JSON을 보존했다. keep 후보인 public_text.go/capacity.go/http.go/store.go/adapters/src/text.ts는 관련 검수에 읽고, omit? 후보 adapters/README.md도 실제 지원 인터페이스의 필수 근거로 읽는다. designer의 capacity/store omit?는 필요 시 확인이다. c6f0848 준비 리뷰는 coor PLANS 병합 오류 발견 뒤 미배정으로 남았으며 현재 새 fixed09c523d에서만 검수한다.
