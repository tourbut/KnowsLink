---
title: SAR-LANDING-GUIDE-001-OPS — 독립 수락한 공개 첫페이지를 기존 운영에 배포하고 상태 보존·공개 안내를 확인한다
status: draft
updated: 2026-10-14
owner: ops
tasks: [SAR-LANDING-GUIDE-001-OPS]
summary: 독립 수락한 공개 첫페이지를 기존 운영에 배포하고 상태 보존·공개 안내를 확인한다
attempt: 552ff51c9ba34fcd90790bccbe37a104
base: fe2387a7676f7997aea76a2606f4cc2dbcd4b56b
subagent_level: standard
test_level: lite
---

# SAR-LANDING-GUIDE-001-OPS — 수락한 공개 첫페이지 운영 반영

- 작성일: 2026-10-14 / From coor → ops / 상태: queued.
- 선행: fe2387a 구현의 독립 코드 리뷰·tester 기능 QA·designer 시각 수락과 main/origin 통합. 지금 배포하지 않는다. coor가 정확한 수락 runtime SHA와 준비 HEAD를 고정하고 route/find/context/packet을 갱신한 뒤 dispatch한다.
- 실행: codex/gpt-6-astra/high, 최신 사용자 지정. test_level=lite/subagent_level=standard.
- 담당: repo0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/ops / fullops/ops.
- 복귀: coor / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_38ae16dc-0ca4-4beb-ab00-d1e1b641f061 / run_5534969da2d0 / 실제 새 task·dispatch preamble.
- 병합 책임: coor / main·origin/main. 결과 문서도 별도 세션 고정 SHA 리뷰 후 통합한다.

## 적용 기준과 승인 범위

AGENTS/FULLOPS·fullops-common-0.3.3·project.md·문서 규칙·D07/D12/D13을 유지한다. 사용자의 기존 '배포해'와 후속 현재 사이트 첫페이지 개편 요청을 기존 link.knowslog.com의 동일 서비스 반영으로 이어서 수행한다. 별도 신규 서비스 공개나 DNS/Access/계정 정책 변경 권한은 없다. 설계·DEV·QA 인박스에서 제외한 운영 행동은 이 OPS 역할에서만 수행한다.

서버의 기존 KnowsLink checkout·relay/migrate만 수락 코드로 전환할 수 있다. 지정 서버의 과제 고유 private backup·restore fixture·작업 디렉터리는 이 배포에 필요한 승인 범위다. 원 DB/키/친구/인증/Inbox/Events·다른 공유 서비스·기존 이미지/백업/사용자 터미널을 보존한다. 삭제·이동·prune·volume 삭제·실계정 로그인/발송·pull/ACK/token 소비·Grok 로컬 bundle 변경은 하지 않는다. 문서의 삭제 포함 beta.sh 절차는 실행 지시가 아니며 무삭제 기존 helper를 먼저 읽고 재사용한다.

설치 캐시는1.4.3으로 바뀌었지만 레포 하네스 정본1.4.2를 보존한다. 필요한 현재 스킬/help를 읽고 호환 도구를 사용한다. 이번에 update-fullops·설정/레벨 완화는 수행하지 않는다.

## 먼저 읽을 문서

- .fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/report.md.
- .fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CODE-REVIEW-review/report.md.
- .fullops-squad/docs/design-docs/mockups/service-member-ui.md — 공개 안내와 로그인 계약.
- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-OPS/report.md — 직전 배포·원 이미지 실패·교훈.
- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-OPS/commands.txt.
- .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-OPS/public-smoke.json.
- .fullops-squad/docs/design-docs/database-design.md.
- .fullops-squad/docs/operations/transition.md.
- .fullops-squad/docs/operations/ops-guide.md.
- README.md, docs/operations.md.
- deploy/knowslink/compose.ops.yaml, deploy/knowslink/verify.py, deploy/knowslink/beta.sh.
- .fullops-squad/contexts/ops.md.
- .fullops-squad/project.md, .fullops-squad/lint/lint.json.

경로는 실제 존재를 확인하고 후보 불일치는 숨기지 않는다. 직전 private helper는 C:/Users/shin/AppData/Local/KnowsLinkDevTools/refactor002-ops/ 및 humanpeer001-ops/에 있다. known_hosts/RejectPolicy·자격 메모리 사용·비밀 비출력 계약을 유지한다. 이전 관측/백업을 현재 상태로 대신하지 않는다.

## 해야 할 일과 완료 기준

- [ ] fresh 운영 SHA·image ID·healthy·현재 요청/DB/공유 서비스 baseline을 읽기 전용으로 확인한다. 마지막 관측8a93e38은 참고값이다.
- [ ] 수락 SHA bundle/hash·서버 clean checkout·source ancestor·runtime/DB/배포 입력 diff를 확인한다. 변경이 UI/auth-view 범위를 넘으면 보고한다.
- [ ] 새 build 전에 현재 실행 이미지에 과제 고유 rollback 태그를 붙이고 실제 image inspect/가용성을 확인한다. 직전 원 이미지220edb 조회 실패의 재발을 막는다. 기존 tag를 덮지 않는다.
- [ ] relay가 정상일 때 수락 이미지를 먼저 build한다. 전환 직전 단일 writer stop·frozen epoch 확인·새 private dump/list·network none/host port 없는 고유 격리 복원과 frozen 전체 digest/guard를 검증한다. 기존 backup/volume은 보존한다.
- [ ] 미리 만든 이미지로 relay/migrate만 재개하고 health·restart·로그·현재 schema/guard·상태 digest·공유 컨테이너 보존을 확인한다. 만료/volatile map 차이는 값과 원인을 구분한다. 실패하면 같은 승인 범위의 검증된 이전 이미지로 서비스 복구하고 원 실패를 남긴다. 운영 DB를 과거 dump로 덮지 않는다.
- [ ] 기존 안전 public smoke와 새 공개 홈페이지 내용을 실제 HTTPS에서 확인한다. /소개·설치 details·가상 예시·로그인 provider 및 무인증 회원 경계·보안 헤더를 확인한다. 공개 DOM/캡처로만 검사하고 실로그인/메시지/비밀 URL을 수집하지 않는다.
- [ ] 결과 D12/D13·README 현재 운영 기준·docs/operations·own phase/QA/context를 갱신한다. 제품 source/설계/다른 역할 문서는 수정하지 않는다.
- [ ] 재사용한 exact 제품 gate와 이번 운영 검사를 구분한다. 실행 명령 자신의 종료코드, UTC/KST 전환 시각·중단 시간·새 백업 mode/hash·이미지·보존 assertions·원 실패·미실행·복구 한계를 남긴다.
- [ ] packet outcomes·strict문서·diff검사를 마치고 work.py finish·고정 결과 SHA·실제 worker_done으로 회신한다.

## 검사·증거·소유권

직접 운영 전환은 직렬 처리한다. UI 기능/코드 gate는 coor가 제공할 exact 수락 증거와 runtime 입력 동일성으로 재사용하고 문서만 바뀐 HEAD에서 전체 검사를 반복하지 않는다. 등록 lint 증거와 명령 계약은 바꾸지 않는다. Docker CLI의 현재 버전을 확인하고 필요 시 Context7→연결 실패 시 공식 docs를 쓴다. 직전 Compose5.1.3 up/stop 공식 문서 근거는 같은 버전·같은 옵션이면 재사용 가능하다.

새 비공개 자원은 AppData/Local/KnowsLinkDevTools/landing-guide001-ops 및 서버 knowslink-state 아래 과제 고유 디렉터리에 보존한다. 원 credential/token/cookie/실메시지/DB row는 Git·터미널·보고서에 출력하지 않는다. Git에는 redacted count/digest/status와 재현 명령만 둔다. private dump를 업로드하지 않는다. restore fixture는 stop 후 보존한다. 기존 사용자 소유 OPS 터미널을 종료하거나 새 작업을 주입하지 않는다.

이번 역할 소유권은 own phase/QA/context/inbox/outcomes와 D12/D13·README/docs운영 기준이다. runtime 제품 변경이 필요하면 DEV에 인계한다. main 병합·push·현황판은 coor 담당이다.

## 완료 보고

미착수. 독립 QA·시각 수락과 main 반영을 기다린다. 실제 운영 반영/보존/검사와 실계정·Grok 로컬 업데이트 미실행을 구분한다.
