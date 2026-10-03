---
title: SAR-BETA-002-TESTER — 본인 전용 합성 베타의 실제 브라우저 동작 검증
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-BETA-002-TESTER]
summary: 본인 전용 합성 베타의 실제 브라우저 동작 검증
---

# SAR-BETA-002-TESTER — 본인 전용 합성 베타의 실제 브라우저 동작 검증

- From / To: coor / tester. 상태 ready.
- 기록 checkout: /home/shin/orca/workspaces/KnowsLink/fullops-tester, fullops/tester.
- 복귀: run_8ca8bc058ab7, coor term_9afa8217-862c-404d-9a43-2122427113fc. 새 preamble task/dispatch를 따른다.
- 승인: 사용자 요청에 따른 실제 합성 테스트·자기 QA 코드/증거 커밋. 실제데이터·외부벤더·Access OTP 사용자 대행·제품/배포 소스 변경·공유서비스 변경은 범위 밖이다.

## 대상과 적용 기준

실제 배포 /home/shin/deploy/knowslink의 고정28bd1bb와 제품78b1d92다. 최신main2a3bc0b에 최종 운영/QA 기록이 있다. 기존 독립QA c993d59(공개 Access negative),1762b43(로컬 런타임),제품QA659f4b0와 리뷰7ba9df0/311381f를 원래 실행조건으로 재사용한다. 단순 API 전범위·백업복원·최신gate의 통과 검사를 반복하지 않는다.

사용자는 OTP/Basic 로그인 뒤 요청 브라우저에서 sender_not_allowed를 겪었다. coor 직접HTTP에서 정확한owner GET200·approve/deny POST303·버튼비활성·중복409·다른owner403 sender_not_allowed를 재현했다. 사용자 브라우저의 원인은 아직 관측하지 않았으므로 단정하지 않는다. 이번 새 검증은 실제 브라우저 DOM/버튼 조작·인증컨텍스트에 한정한다.

FULLOPS, fullops-common-0.3.2 README/코딩/테스트/보안, project 정본, 문서 작성 규칙과 fullops-test를 적용한다. 예외 없음. 사용자 직접테스트 요청은 합성 seed 실행 및 인증값의 메모리내 사용을 허용한다. owner-login 명령으로 값을 출력하거나 trace/HAR/스크린샷/로그에 인증값을 포함하지 않는다.

## 먼저 읽을 문서

Jev 결과 docs/evaluations/jev/SAR-BETA-002-TESTER-*의 keep 전부: internal/relay/http.go, adapters/src/synthetic.ts, deploy/knowslink/beta.sh, SAR-BETA-001-TESTER-PUBLIC.md, SAR-BETA-001-TESTER.md, SAR-MVP-001-TESTER-FINAL.md, user-guide, transition, project, 공통 README/세규칙, FULLOPS, 현재 인박스. 충돌/주의 추천 없음. find의 absent는 신규제품구현부재 추천이며 대상은 고정SHA에서 존재한다.

## 해야 할 일

- [x] 배포/product/설정 SHA와 기존 증거의 불변 경계를 확인한다.
- [x] 기존fixture(build/qa-fixture.json)를 값출력 없이 보존하고 finally 복원한다. 테스트 seed는 합성 DB 레코드만 추가한다. 기존DB/fixture/자격값을 삭제하거나 교체한 채 남기지 않는다.
- [x] 설치된 실제 브라우저/E2E 도구를 사용해 loopback 소유자 UI를 실행한다. deterministic DOM/버튼 판정이 가능하면 fullops-test 순서1을 적용한다. 기존 Playwright가 없으면 외부 임시 경로에 최소 도구를 준비하고 저장소 의존성을 바꾸지 않는다. Orca/Jev 브라우저로 진행해야 하는 경우 해당 스킬·시나리오 규칙을 따른다. 브라우저가 실행되지 않으면 HTTP 통과로 대신 표시하지 말고 재현 근거와 차단 조건을 보고한다.
- [x] 정확한 소유자 새브라우저컨텍스트에서 owner dashboard와 fresh gate 렌더, typed body/정책/만료표시, Approve 클릭·결정표시·버튼비활성·새로고침을 확인한다.
- [x] 다른fresh gate에서 Deny 클릭과 결정표시·새로고침을 확인한다. 승인해도 실제일정/정보공개 효과는 없다는 제품조건을 보존한다.
- [x] 다른owner 또는 앞선seed의owner 컨텍스트로 새gate 열기403 sender_not_allowed를 브라우저에서 재현한다. 정확한 새owner 컨텍스트/맞는gate로 바꾸면200복구되는지 확인한다. 실제사용자캐시를 관측한것과 재현시나리오를 구별한다.
- [x] 기존 만료gate가 있으면 같은owner의 읽기전용 상태로 만료·버튼비활성을 확인한다. 없으면 최신gate의180초만료를 실제 기다려 확인한다. 시간을조작하거나라이브DB수정하지 않는다. 캡처는 DOM으로 판정하지 못하는 항목만 만들고 기본trace/HAR를 남기지 않는다.
- [x] 새 qa-reports/SAR-BETA-002-TESTER.md, 관련test/와 scenarios, tester context/inbox/logs에 실행 명령/환경/시나리오/종료코드/한계를 기록한다. lint/strict/whitespace와 fixture복원·배포/제품diff 불변을 확인하고 finish/commit/worker_done한다.

## 완료 기준과 후속

브라우저 정상 approve/deny, 표시/새로고침, 만료 차단, wrongowner403 재현과 새context200복구를 검증한다. 모바일 실제기기·사용자OTP·공개인증후UI의 직접브라우저관측은 loopback브라우저와 구별한다. 기존사용자로그인성공 및 공개negativeQA를 연결하되 미실행을PASS로바꾸지 않는다. 새critical/high 또는 제품결함이면 원인·재현·영향을 보고하며 소스를 직접 수정하지 않는다. 갱신할 D01–D13은 없음(참조만). 새기능배정 없이 이 요청 테스트결과를 보고한다.

## 완료 보고

상태: 독립 브라우저 QA를 마쳤다. 배포 `28bd1bb2bdd4b90d0c6f5d2a8d4ba3e00f630ff2`와 제품 `78b1d92c8aa626245d3349ffaf7367d28f1dd3ef`는 읽기 전용으로 유지했다.
headless Google Chrome 148이 정확한 owner의 승인·거절, 표시와 새로고침, 기존 만료 gate의 버튼 비활성, 이전 owner와 다른 owner의 `403 sender_not_allowed`, 맞는 owner의 새 컨텍스트 200을 확인했다. 판정 명령 종료코드는 0이다.
사용자 브라우저 캐시, 이메일 OTP, 공개 로그인 뒤 UI는 실행하지 않았다. 기존 공개 negative, 로컬 런타임, 제품 QA와 리뷰는 원래 SHA로 재사용했다. 새 critical/high는 없다.
변경은 QA 보고서, 시나리오, `browser.mjs`와 결과, tester 컨텍스트, 이 인박스다. 제품 코드와 배포 소스와 D01–D13과 PLANS는 바꾸지 않았다.
지시와 다른 점은 두 가지다. 저장소에 Playwright가 없어 외부 `playwright-core` 1.63.0으로 Chrome을 실행했다. fixture owner의 expired gate가 이미 있어 180초 대기는 하지 않았다.
첫 SQL 집계는 브라우저 기동 전에 실패했다. 그 실행은 판정이 아니다. Jev 충돌·주의 후보는 없다. 제외 추천 문서를 추가로 읽어야 판정이 바뀌지는 않았다.
검증한 것은 위 브라우저 checks, fixture 복원, `relay_state` 1행 유지, 배포 clean이다. 검증하지 않은 것은 사용자 프로필, OTP, 모바일 실기기, 공개 로그인 뒤 UI, adapter completion `denied`의 재실행이다.
산출물은 [SAR-BETA-002-TESTER.md](../docs/evaluations/qa-reports/SAR-BETA-002-TESTER.md)다. 후속은 coor의 기록 통합이다. 새 기능 배정은 없다. 보고 커밋 SHA는 worker_done 본문에 적는다. 이 QA 완료는 전체 제품 수락이나 공개 운영 수락이 아니다.
