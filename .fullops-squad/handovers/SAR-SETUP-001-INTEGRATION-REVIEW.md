---
title: 초기 통합 후보의 최종 독립 리뷰
status: draft
updated: 2026-10-03
owner: dev
tasks: [SAR-SETUP-001-INTEGRATION-REVIEW]
summary: main 기준의 안정된 후보에서 기존 QA 증거와 신규 기록 및 러너를 검토한다
---

# SAR-SETUP-001-INTEGRATION-REVIEW — 안정된 통합 후보의 최종 리뷰

목표는 main 수락 전 실제 병합 refs의 독립 리뷰 게이트 완성이다. 제품 규칙이나 제품 기능을 바꾸지 않는다.
고정 base는 f94510f48a050c08254ad028ac66a14a96946b3a이며 main 현재 HEAD다. 고정 head는 59be02d다.
읽기 전용 detached snapshot은 /tmp/SAR-SETUP-001-integration-review-59be02d다. 기록 체크아웃은 /home/shin/orca/workspaces/KnowsLink/fullops-coor다.
이 지시서는 기록 체크아웃의 절대경로로 제공한다. 후보 head 이후의 운영 인계 파일이며 제품 후보를 변경하지 않는다.

먼저 설치된 fullops-review 및 open-code-review-delegate 스킬을 읽는다. snapshot의 .fullops-squad/FULLOPS.md, rules/common/README.md와 연결된 세 규칙, project.md, docs/agents/document-writing.md, D02, D03, dev 완료 지시서와 실행 기록을 읽는다. 적용 규칙은 fullops-common-0.3.2다. 기존 구현 시점의 0.3.1 기록은 보존한다.
기존 리뷰는 snapshot의 docs/evaluations/qa-reports/SAR-SETUP-001-DEV-099-review/다. 독립 QA와 종료코드 보완은 SAR-SETUP-001-TESTER.md 및 SAR-SETUP-001-TESTER-test/의 요약과 보완 로그로 확인한다. 경로는 .fullops-squad/ 기준이다.

소유권은 기록 체크아웃의 .fullops-squad/docs/evaluations/qa-reports/SAR-SETUP-001-INTEGRATION-review/뿐이다. 준비된 preview/rules/result/report를 이어서 완성한다. 보고서에 front matter를 등록한다. 제품 코드·snapshot·PLANS·board·다른 지시서는 수정하지 않는다. git add -A를 쓰지 않고 소유 파일만 커밋한다.

main 기준 전체 diff를 preview에 따라 검토한다. 제품 경로가 기존 독립 리뷰 head 0cc10b083771be9b3423833b222c57d426315333과 동일함을 확인하면 기존 직접 리뷰를 원래 SHA와 함께 재사용한다. 독립 QA 대상과 코드 동일성을 확인하여 기존 QA를 재사용한다. 전체 제품 테스트를 다시 실행하지 않는다. 신규 QA run.py와 inject.py는 직접 읽고 종료코드 보완·원복·비밀값 경계 및 실제 로그와 대조한다. 기획·원천·인계·운영·문서 diff는 기존 결정·원본 보존과 실제 내용의 정합성을 확인한다. 원시 로그는 무결성과 요약·명령 종료코드로 확인하고 구체적인 skipped 이유를 기록한다. OCR 제외 파일을 누락하지 않는다.

result.json의 모든 파일과 발견 사항을 기록한다. 실제 구현자 Codex 세션 01a0ffa6-bcbc-7383-a8e5-f14521f0dfa3는 기존 리뷰의 rollout 직접 근거로 확인한다. QA 러너 작성 세션은 tester Dispatch ctx_8bc7450abd67 / ctx_a0b3241794f3의 실제 Claude 세션을 추가 확인한다. 새 검토자의 실제 세션은 둘과 달라야 한다. snapshot head/detached/clean/read_only를 기록한다.

설치된 lint.py를 snapshot에서 --from f94510f48a050c08254ad028ac66a14a96946b3a로 실행해 lint.json에 기록한다. 지시서의 기존 완료 검사 기준 729446d8da57 결과는 이미 존재하며 재사용한다. main 기준 commands가 비어 있어 제품 lint를 실행하지 않는 한계는 명시한다. review.py check의 key는 SAR-SETUP-001-INTEGRATION, from/to는 위 고정 refs다. 실제 명령 종료코드를 보존한다. critical/high나 필수 증거 결함이 있으면 수락을 차단한다.

완료 조건은 모든 파일의 검토 또는 근거 있는 생략, 실제 독립성, lint와 check 결과, main 병합 후보에 대한 수락 또는 차단 결론이다. 결과 SHA·refs·심각도·로그/보고서·남은 수락 조건을 새 preamble의 worker_done으로 한 번 보고한다. main 병합과 외부 발송·배포는 하지 않는다. Run은 run_8ca8bc058ab7이다.

## Coordinator 실제 병합 refs 갱신

main 준비 기준은 fa971df3a36e다. 기존 승인된 외부 원천 전용 lint 제외만 반영했으며 준비 lint 종료코드 0, ERROR 0, WARNING 2다. 최종 후보는 ec190f98848fbe21cc4df6d78e7d01dca0d58129이고 이 준비 커밋을 조상으로 갖는다. 새 읽기 전용 snapshot은 /tmp/SAR-SETUP-001-integration-review-ec190f9다. 최종 review key는 SAR-SETUP-001-INTEGRATION-FINAL이며 대응 디렉터리를 소유권에 추가한다. 이전 INTEGRATION 기록과 당시 DOC-003 실패는 보존한다. 검토한 제품·QA 증거는 동일성을 확인해 재사용하고 준비·병합 운영 diff만 추가 검토한다. lint와 check의 실제 refs는 fa971df3a36e..ec190f98848fbe21cc4df6d78e7d01dca0d58129를 사용한다.

## Coordinator 원본 원천 준비와 최종 refs

0.9.10 DOC-003이 exclude보다 먼저 실행되어 기존 제외는 원천 front matter 오류를 막지 못한다. 앞선 INTEGRATION과 FINAL의 ERROR 7을 보존한다. 기존 외부 스냅샷 원문을 main 준비 커밋 3eb7647938111c9f13ad523760aeb7e90c7fa7f3에 그대로 반입했다. 원천 반입의 lint ERROR 7 및 원문 hard break의 whitespace 실패는 별도 기록하고 통과로 표시하지 않는다. 원문과의 무결성 diff는 0이다. 새 최종 head는 9c96232e6230e00319c33ff77637c80923e2438b, snapshot은 /tmp/SAR-SETUP-001-integration-review-9c96232, 최종 key는 SAR-SETUP-001-INTEGRATION-SOURCE다. 대응 리뷰 디렉터리를 소유권에 추가한다. 이 실제 base/head로 lint와 check를 완료한다. 기존 제품·QA 검토는 동일성 확인 뒤 재사용한다.
