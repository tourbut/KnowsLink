---
title: tester 컨텍스트
status: draft
updated: 2026-10-15
owner: tester
tasks: [SAR-REFACTOR-002-TESTER, SAR-LANDING-GUIDE-001-TESTER, SAR-UX-ISSUES-002-RECEIPT-TESTER]
summary: 독립QA의 원실패 보존 및 CSS수정 source동일성·gate 재사용 경계를 기록한다
---

# tester 컨텍스트

10KB를 넘은 이전 전문은 [2026-10-14 archive](archive/tester-2026-10-14-refactor002.md)에 복사했다. 파일 삭제·이동은 하지 않았다. archive의 원래 상대 링크는 이 contexts 디렉터리를 기준으로 읽는다. 과거 실패·held·수락은 각 보고서의 당시 SHA에 남으며 최신 QA로 덮지 않는다.

- tester는 제품을 고치지 않는다. 합성 fixture·테스트·시나리오·증거를 소유하고 실패는 재현 절차로 dev/coor에 전달한다.
- 고정 제품 SHA와 기록 HEAD를 구분한다. 제품·의존성 동일성을 확인한 DEV 증거만 재사용하며 새 독립 실행으로 바꾸어 기록하지 않는다.
- 실제 명령 종료코드를 보존한다. 환경·기대값·증거 결함의 실패도 남기고 원인이 바뀐 뒤에만 재검증한다.
- 기준 lint.json에 등록 명령이 없으면 lint 성공이 제품 실행 증거가 아니다. Windows 미지원 도구는 exact HEAD·clean을 검사하는 고유 Linux clone에서 실행한다.
- SSH는 기존 known_hosts/RejectPolicy·비밀값 가림을 유지한다. DEV HEAD를 올리는 make.cmd를 tester 고정 clone에 재사용하지 않는다.
- Git Bash 원격 경로 전달은 MSYS_NO_PATHCONV=1을 확인한다. test_public_config.sh의 stat -c는 Linux 전용이다.
- packet과 실제 task/role/attempt/SHA를 맞춘다. 가짜 task-key로 검사기를 우회하지 않는다. 과거 archive의 해당 우회 조언은 2026-10-11 coor 정정으로 폐기됐다.
- shell_guard 거부를 다른 셸·표현으로 우회하지 않는다. 삭제·이동 정리는 현재 AGENTS와 coordinator 권한을 따른다.
- Codex hook 검증은 실제 host 실행 방식과 종료코드2·deny JSON을 함께 확인한다. 문자열 리터럴/argv cmd 호출의 성공만으로 guard 실행을 주장하지 않는다.
- Orca 브라우저 viewport는 goto 뒤 적용하고 실제 innerWidth를 측정한다. 화면·CLI·실계정·호스트 표시 증거는 각각 구분한다.
- 웹 사본 보관과 활성 transport 인가는 별도다. 관계 철회 뒤 참여자 이력이 남아도 과거 전송·claim·세대 권한이 되살아나면 안 된다.

## 최근 결과

- 2026-10-13 SAR-REFACTOR-001-TESTER: exact cd4dae8의 Linux race8개·MCP/owner-text·독립 text15개 exit0. [보고](../docs/evaluations/qa-reports/SAR-REFACTOR-001-TESTER-test/report.md). 제품 변경0, 실계정·운영·SQL 이행은 미실행이다.
- 2026-10-14 SAR-REFACTOR-002-TESTER: exact bcb1dd66 Linux race16개·새 Go 경계10개·Inbox 오류 matrix·합성 DB race4개 exit0. 제품 변경0, 새 C/H 발견0. [보고](../docs/evaluations/qa-reports/SAR-REFACTOR-002-TESTER-test/report.md).
- ACK completion received는 TTL 뒤에도 유지한다. 첫 잘못된 기대값/native1을 보존하고 기존 계약으로 바로잡았다. sweep CPU 이득은 DB/HTTP 처리량이 아니다. 단일 구독 claim 증가와 local claim crash window는 남는다.

- 2026-10-14 SAR-LANDING-GUIDE-001-TESTER: fixture4795/구현fe2387a의 독립 HTTP34·브라우저22·assert7·native키보드·실제Chrome200%를 확인했다. [보고](../docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-TESTER-test/report.md). 기존 회원 CTA 대비1.846967은medium 실패이며 DEV 후속이다.
- Chrome settings의 Page zoom과 DPR2/CSS zoom1/visualViewport.scale1을 함께 확인한다. viewport/pinch/CSS 확대를 실제 browser200%로 바꾸지 않는다. memberStyle 동일성과 member_pages.go 파일 전체 동일성은 별개다.

- 2026-10-15 SAR-LANDING-GUIDE-001-TESTER 후속: CSS수정 f2b8cb6의 독립Chrome44·source18·예시간격 실제200% PASS. 원대비FAIL/원5149증거는 보존하고 V-LG-01/02의 수정후 해소만 판정했다. [보고](../docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-TESTER-test/css-followup/report.md).
- same-source Linux fc4 gate재사용과 현재QA HEAD Windows make unavailable을 구분한다. fixture HTML hash와 source style동일성은 인증·DB 성공을 뜻하지 않는다.

- 2026-10-17 SAR-UX-ISSUES-002-RECEIPT-TESTER: fixture 18928/DEV 9ae29ec의 HTTP16·Chrome/Edge CDP 각51·철회 후7 검사. U6 격리·복구 통과, medium D-RT-01(닫힌 선택기가 긴 이름의 구분값을 가림). [보고](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-TESTER-test/report.md).
- Orca 탭은 hasFocus false라 Tab 입력이 무효이고 페이지 CSP가 fetch POST를 막는다. 키보드는 Chrome/Edge CDP Input.dispatchKeyEvent를 쓴다. `__Host-` 쿠키는 `cookie set --domain`이 아니라 `exec "cookies set ... --url"`로 넣는다.
- fixture의 `revoke`는 되돌릴 수 없다. 모든 다른 검사 뒤에 실행하고 이후 재실행은 철회 상태를 기대값으로 따로 둔다. 회원당 40회 조회 한도는 기본 폼 GET도 센다. 검사 볼륨을 `recover`로 비운다.

- 2026-10-17 SAR-SESSION-PEER-001-TESTER: 고정 bdc6500의 V01·V03–V07 독립 QA. 승인 대화상자 본문 200자 잘림 medium D-SP-01과 관찰 6건. V02 미실행. [보고](../docs/evaluations/qa-reports/SAR-SESSION-PEER-001-TESTER-test/report.md).
- relay는 비-loopback 평문 SMTP를 거부한다. 합성 메일 sink는 relay 컨테이너 netns 사이드카(`--network container:`)로 띄운다. 코드 메일은 주소당 1분 1건이며 IP 시간당 20건이다. Git Bash는 `/knowslink-session`을 경로로 바꾸므로 MSYS_NO_PATHCONV=1이 필요하다.
- Claude 폴더 신뢰 대화상자는 빈 Enter가 `No, exit`다. 방향키 Down 뒤 Enter. 호스트 모델은 주입 본문의 작성을 거부할 수 있어 수신 쪽 검증은 서명 fixture 발송을 쓴다.

- 2026-10-19 SAR-UX-ISSUES-002-NAMES-TESTER: 고정 b2e683a의 U2·U4·U7·1단위 독립 QA. critical·high 없음, low D-NM-01(이름 거부 사유·보이지 않는 문자 되돌림). [보고](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/report.md).
- 이 fixture는 `/home/*` 전체에 회원당 분당 40회 한도가 있다(41번째 429). client·CDP 러너에서 429를 기다리게 한다. headless Edge는 Tab을 body로 빼지 않고 페이지 안에서 순환하니 첫 정지점 재도달을 한 바퀴로 본다.
- 없는 코드·제어문자·500자 코드도 일반 성공을 돌려주므로 거부 기대를 두지 않는다. 비가역 단계가 단언 하나로 멈추지 않게 mutate와 followup으로 나눴다.

- 2026-10-20 SAR-UX-ISSUES-002-APPROVAL-TESTER: 고정 9d13a05의 U5 승인 요약·U7·O-NM-04·O-NM-08 독립 QA. critical·high·medium 없음, low 관찰 3건. [보고](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-TESTER-test/report.md).
- 승인 상세 CSP(`default-src 'none'`)가 페이지 안 `fetch()`를 막는다. 상태 변경 뒤 읽기 전용 단계를 다시 돌리면 결정 전 캡처가 덮이니 캡처 단계는 변경 전에 끝낸다. CSRF token은 세션에 묶여 다른 세션 재제출은 409가 아니라 403이다. 셸 가드가 파일 삭제 명령을 막으니 조립용 임시 파일은 scratchpad에 둔다.
