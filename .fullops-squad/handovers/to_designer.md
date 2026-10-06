---
title: SAR-PUBLIC-MESSAGES-001-UI — 일반 회원 메시지·gate 직접 UI 검수
status: draft
updated: 2026-10-06
owner: coor
tasks: [SAR-PUBLIC-MESSAGES-001-UI]
summary: 일반 회원 메시지·gate fixed 후보 직접 UI 검수 인계
---

# SAR-PUBLIC-MESSAGES-001-UI — 일반 회원 메시지·gate 직접 UI 검수

- 상태: ready. fixed 09c523da8a3407288d9f5d711e1834af12bc7808
- 담당: designer / /home/shin/orca/workspaces/KnowsLink/fullops-designer / fullops/designer
- 복귀: repo 818c78e5-d51c-4ff4-aa88-70e9ee185fbb / fullops-coor / term_6895aaf1-7b43-4fe0-a416-76f1255a5946 / run_8ca8bc058ab7. Task/Dispatch는 preamble 기준.

## 적용 기준과 먼저 읽을 문서

fullops-common-0.3.3의 README·coding-style/testing/security, FULLOPS·project·document-writing·contexts/designer를 전달 SHA에서 읽는다. 제품 기준은 docs/planning/product-specs/SAR-PUBLIC-SERVICE.md PS08–11·PS04/06/07와 SAR-MVP.md C1–C5 및 docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md UX06/07다. DEV 실행 기록 docs/exec-plans/phases/SAR-PUBLIC-MESSAGES-001-DEV.md와 logs/2026-10-06_to_dev.md의 최신 과제·DEV QA 증거/mobile-width.py·README/adapters README를 읽는다. Jev code/doc/context도 이 키로 준비한다.

## 해야 할 일과 완료 기준

- [ ] fixed 09c523d을 별도격리clone·자기fixture에서 실제 브라우저로 직접 검수한다. DEV 자동검사/390px 측정을 직접시각 PASS로 대신하지 않는다. UX06/07의 실제 UI가 단순공개fixture와 달라질 수 있어 일반회원 세션/자기receipt/검증gate를 대상으로 한다. 비밀값·이메일·코드·key/credential은 출력/캡처 전에 가린다.
- [ ] UX06: 선택agent·명시송신 도구안내·요청ID/관련답장ID·queued/실제수신/회신·TTL/manual pull·offline/expired/rate/revoked/invalid-key/conflict 상태와 실제 가능한 다음동작을 직접 확인한다. 자유 composer/채팅버블/장기timeline/queued성공오인 배너를 추가하지 않는다.
- [ ] UX07: 검증 typed-body·정책·발신/대상·기한이 판단근거이고 hint/HTML/명령이 실행되지 않음을 확인한다. 원문부재/만료/철회 approve비활성·deny·오류 뒤 복구안내, 회원권한분리를 실제 화면에서 확인한다. 긴ID·390px모바일·키보드이동·읽을 수 있는 오류를 확인한다.
- [ ] 지정 시각조건의 캡처·manifest·고정SHA·실제장애/초기실패·정상fixture·자기cleanup을 보존한다. 시간변화가 정지화면으로 판정되지 않는 경우만 영상이다. 변경없는 UI는 관련의존성 동일성과 원래SHA/조건을 연결해 재사용한다. 새 UI는 기존 AGENTS PASS로 대체하지 않는다.
- [ ] 보고서 docs/design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md와 필요한 시각증거·자기실행기록을 작성한다. final 기록HEAD FullOps --from09c523d lint/test·strict·diff exit/HEAD/경고를 남긴다. work.py finish로 전문아카이브/빈인박스·최종커밋·일반push 뒤 worker_done fullSHA/증거/수락·미검증·후속을 보고한다.

## 소유권·제약·후속

designer는 새직접검수보고/시각증거·자기context/PLANS·인박스/완료로그만 작성한다. 코드·기술정본·타인인박스·기존실패/리뷰는 수정하지 않는다. 새제품판단이 필요한 때만 coor로 질문한다. theme/designlint/Tailwind/shadcn은 현스택에 없으며 기존Go template/memberStyle·UX정본을 기준으로 직접판정한다. OPS/TESTER 검사와 직접시각PASS는 분리한다. Workers Free 유지, 과금/유료전환/구독/운영배포/실메일/외부계정/공유서버/Tunnel/운영자료삭제는 수행하지 않는다. 로컬격리검사·자기fixture회수·문서커밋·일반push는 승인됐다. 실제운영 PS08/13/14·실메일/공개/실24h/노우↔다닷은 미검증으로 기록한다.

## 갱신할 산출물과 완료 보고

D04 직접검수기록 연결만 실제영향에 따라 갱신한다. 제품/UX규칙은 변경하지 않는다. final fullSHA·고정대상·직접판정/실패·캡처manifest·명령exit/경고·미검증·후속담당/재개조건을 전문으로 작성한다. front matter는 deliverables stamp를 쓴다.

추가 필수 확인: internal/relay/member_receipt.go가 실제 경로다. context의 member_receipts.go 거부는 오기이며 원본 JSON을 보존했다. keep 후보인 public_text.go/capacity.go/http.go/store.go/adapters/src/text.ts는 관련 검수에 읽고, omit? 후보 adapters/README.md도 실제 지원 인터페이스의 필수 근거로 읽는다. designer의 capacity/store omit?는 필요 시 확인이다. c6f0848 준비 리뷰는 coor PLANS 병합 오류 발견 뒤 미배정으로 남았으며 현재 새 fixed09c523d에서만 검수한다.
