---
title: SAR-PUBLIC-MESSAGES-001-REVIEW — 일반 회원 메시지 후보 독립 보안 리뷰
status: draft
updated: 2026-10-06
owner: coor
tasks: [SAR-PUBLIC-MESSAGES-001-REVIEW]
summary: 일반 회원 메시지 fixed 후보 독립 보안 리뷰 인계
---

# SAR-PUBLIC-MESSAGES-001-REVIEW — 일반 회원 메시지 후보 독립 보안 리뷰

- 상태: ready. fixed 09c523da8a3407288d9f5d711e1834af12bc7808
- 담당: ops / /home/shin/orca/workspaces/KnowsLink/fullops-ops / fullops/ops
- 복귀: repo 818c78e5-d51c-4ff4-aa88-70e9ee185fbb / fullops-coor / term_6895aaf1-7b43-4fe0-a416-76f1255a5946 / run_8ca8bc058ab7. Task/Dispatch는 preamble 기준.

## 적용 기준과 먼저 읽을 문서

fullops-common-0.3.3의 README·coding-style/testing/security, FULLOPS·project·document-writing·contexts/ops를 전달 준비 SHA에서 읽는다. 제품 기준은 docs/planning/product-specs/SAR-PUBLIC-SERVICE.md PS08–11·PS04/06/07와 SAR-MVP.md C1–C5 및 docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md UX06/07다. DEV 실행 기록 docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV.md와 logs/2026-10-06_to_dev.md의 최신 과제·DEV QA 증거·COOR/dev-final-37f9a1e-lint.json을 읽는다. Jev code/doc/context 결과도 이 과제 키로 준비한다.

## 해야 할 일과 완료 기준

- [ ] fullops-review와 open-code-review-delegate를 적용한다. base7efbaa349a5857eb1eac859a955ec3a09c91f800..fixed09c523d을 실제 diff와 관련 코드/검사/산출물로 직접 검토한다. preview/rules/result/report는 docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-REVIEW-2-review/에 coor가 준비한다.
- [ ] snapshot /tmp/knowslink-messages-review-09c523d은 깨끗한 detached 고정 SHA이며 읽기만 한다. 실제 구현 DEV session 01a11106-b736-7db1-b518-a65b64dbc5fb와 자신의 실제 session ID가 다름을 independence에 기록한다. 리뷰 결과는 OPS 기록 checkout에 쓰고 snapshot을 수정/설치/커밋하지 않는다.
- [ ] 회원·양측 agent/키·현재 관계/세대 인가, 별도 text와 frozen wire 분리, 1회 관련 회신·명시 consent·receipt/idempotency·ACK/TTL본문삭제·metadata24h, raw/queue/gate/receipt/HTTP/claim 상한·high priority·공유 동시성/재시작·process/DB timeout·안전 정리 budget, gate verified typed body/CSRF/XSS/현재 권한/만료/철회 경계를 검토한다. 기존 low L-A/L-B·공개전 합성가입unset/운영DB 합성owner0을 보존한다.
- [ ] result의 모든 path/status reviewed/skipped·사유·reviewer/conclusion·independence·findings와 report의 적용기준/재현/검사/커버리지/생략 영향·수락 결론을 채운다. SIZE-001/002·SLOP·SEC fixture·DEP0/DESIGN0 판단을 확인한다. 경고를 자동 수락하지 않는다. 새 critical/high 및 필수 실패는 수락 차단이다.
- [ ] fixed 09c523d의 FullOps --from7ef lint/test 증거를 scratch clone(같은 SHA)에서 make install 뒤 실행해 review/lint.json에 남긴다. snapshot은 읽기전용 유지. 기존 원실행은 원래 SHA/조건으로 연결해 재사용하며 review.py check --task-key SAR-PUBLIC-MESSAGES-001-DEV를 수행한다. records 최종 HEAD의 별도 필수 lint/test·strict·diff도 완료한다.
- [ ] 원본 실패·report/result를 보존하고 자신의 OPS 기록을 최종커밋·일반push한다. work.py finish로 완료 전문 아카이브/빈인박스·worker_done fullSHA/증거/수락결론을 보고한다.

## 소유권·제약·후속

OPS는 신규 리뷰 결과/증거·자기 실행 기록/context/PLANS·인박스/완료로그만 작성한다. 제품코드·기획/UX정본·frozen원천·타 역할 인박스는 수정하지 않는다. 로컬 fixture·실제 Node/MCP 프로세스 PASS와 실제 외부 플랫폼/운영PS08/13/14를 구분한다. 운영배포·실메일·실24h·외부계정·공유서버/Tunnel·운영자료삭제·유료전환/구독/초과과금은 수행하지 않는다. Workers Free 유지. 로컬 격리 검사·자기 자원 회수·문서커밋·일반push는 승인됐다. QA/UI는 별도 담당이다. 리뷰 결과는 기술 수락 판단이며 운영 공개 수락이 아니다.

## 갱신할 산출물

없음. 리뷰 기록/증거/실행 기록만 작성한다. front matter는 deliverables stamp를 쓴다.

## 완료 보고

브랜치/fullSHA·고정 대상/실제 독립session·수락/결함·재현/근거·명령exit/HEAD/경고·미검증·후속 담당/재개조건을 전문으로 작성한다.

추가 필수 확인: internal/relay/member_receipt.go가 실제 경로다. context의 member_receipts.go 거부는 오기이며 원본 JSON을 보존했다. keep 후보인 public_text.go/capacity.go/http.go/store.go/adapters/src/text.ts는 관련 검수에 읽고, omit? 후보 adapters/README.md도 실제 지원 인터페이스의 필수 근거로 읽는다. designer의 capacity/store omit?는 필요 시 확인이다. c6f0848 준비 리뷰는 coor PLANS 병합 오류 발견 뒤 미배정으로 남았으며 현재 새 fixed09c523d에서만 검수한다.
