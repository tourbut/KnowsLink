---
title: SAR-HUMAN-PEER-001-TESTER — 고정 DEV 후보의 사람 친구와 대화 및 기존 데이터 권한 이행 독립 QA
status: draft
updated: 2026-10-12
owner: tester
tasks: [SAR-HUMAN-PEER-001-TESTER]
summary: 고정 DEV 후보의 사람 친구와 대화 및 기존 데이터 권한 이행 독립 QA
attempt: aaa81f41cc094961803f3c01f36eaddf
base: f20b33a11addeb3debcb23e2e835be31c7dd1454
subagent_level: standard
test_level: lite
---

# SAR-HUMAN-PEER-001-TESTER — 고정 사람 친구·대화 후보의 독립 lite QA

- 상태: waiting-dev-sha. DEV 완료 SHA·기술 문서·검사/캡처 환경을 받을 때까지 실행하지 않는다.
- From / To: designer / tester. Purpose: review. Task key: SAR-HUMAN-PEER-001-TESTER.
- 담당: C:/Users/shin/orca/workspaces/KnowsLink/tester, fullops/tester. 시나리오·QA 증거·보고서만 수정한다. 제품 코드는 수정하지 않는다.
- 복귀: repo 0b08ec4c-9e3a-4613-8197-5a835e545335, coor 워크트리, term_47da6145-4e1e-4cc1-93ea-20bc0dd590b9, Run run_5534969da2d0. task/dispatch는 새 preamble이다.
- coor가 DEV 고정 SHA를 40자리로 인박스와 packet에 연결하고 실제 QA 기준을 고정한다. 현재 base는 기획 기준이며 구현 SHA로 오인하지 않는다. 실행 호스트/model/effort는 coor가 새 dispatch에서 고정한다.

## 해야 할 일과 파일 소유권

- [ ] docs/evaluations/scenarios/SAR-HUMAN-PEER-001-TESTER.md에 HP-01–16과 독립 assertions를 기록한다. 고정 DEV SHA와 환경/버전을 먼저 확인한다.
- [ ] 격리 회원 A/B/C, A의 agent 두 개·B의 agent 두 개·C의 agent 한 개로 실제 auth/API/DB 흐름을 검사한다. 합성 이메일/키만 쓰고 비밀값을 증거에서 제외한다.
- [ ] 이름 등록/변경/재로그인, 친구 요청·수락과 양방향 대화를 실제 UI에서 수행한다. 등록 이후 추가한 agent도 보이고 agent별 재수락 없이 연결됨을 확인한다.
- [ ] 미수락·거절·취소·만료·동시반대요청·중복 수락을 확인한다. 검색/코드만으로 타인 이름·agent·이메일이 노출되지 않고 제삼자 C의 직접 URL/API/발신 위조가 거부됨을 확인한다.
- [ ] 해제/차단/철회 직전 queued·leased·gate와 열린 composer를 만든다. 직후 신규 전달·ACK/claim/approve/result 및 legacy/SSE/remote 우회가 막히는지 확인한다. 재수락은 새 세대이며 옛 권한을 복구하지 않는다.
- [ ] 같은 발송 시도 응답 유실/중복 클릭은 하나만 전송됨을 확인한다. 변경 본문은 conflict이고 새 메시지는 새 key다. 최대 bytes와 TTL 경계, 24h 보존 전후와 본인 참여 복사본의 읽기 전용 경계를 검사한다.
- [ ] 실제 이전 형식 fixture로 member/agent/key/credential·정확한 pair 보존, pending 종료, 자동 친구 확대 금지, 같은 owner 전환 동의, 소유 불명 fail-closed·rollback/복원 차단을 검증한다.
- [ ] 320/390/1440, Tab/Enter/Space, 뒤로·새로고침·취소, 긴/동명 이름, 오류 뒤 복구, 승인/gate 회귀를 확인한다. 대화 대상이 바뀌거나 로그아웃 뒤 자동 발송하지 않아야 한다.
- [ ] qa-reports/SAR-HUMAN-PEER-001-TESTER-test/report.md에 범위별 PASS/FAIL/not_run과 명령 종료코드·SHA·요청/캡처 경로를 기록한다. designer에 실제 구현 캡처를 인계한다.

## 완료 기준과 검증

현행 D02의 '대화의 연속성과 과거 열람 경계'를 독립 확인한다. 세 번 이상 왕복하고 앞선 메시지 TTL 뒤 새 발송을 검사한다. 해제 뒤 참여 회원은 잔여 보존 원문만 읽으며 철회 키·만료 세션·제삼자는 읽지 못한다. 친구 코드 교체·없는 코드·차단 코드가 대상 존재를 노출하지 않아야 한다.

핵심 순서는 HP-16이다. 실제 동작과 정적 fixture를 구분하고 과거 수락 결과를 새 결과로 재명명하지 않는다. lite 범위로 핵심과 필요한 실패 경계만 검사하되 보안·데이터 손실·필수 프로젝트 검사는 유지한다. DEV와 같은 코드/환경의 재사용 근거가 명확한 변경 없는 검사만 원 SHA·조건을 표시해 재사용한다. 핵심 신규 친구/대화/이행은 독립 실행한다.

시각 증거는 로그인/이름, 대화 홈·본문/발신수신, 친구 요청/목록, 내 agent, 해제 읽기전용·오류의 390/1440 캡처다. 320폭은 overflow/동작 측정과 필요한 결함 캡처를 남긴다. 200% 확대를 짧게 확인하고 미실행이면 별도 표시한다. 실제 Google/노우/다닷/운영은 이 합성 QA의 성공으로 대체하지 않는다.

제품 전체 수락은 독립 QA뿐 아니라 designer 직접 구현 시각 판정과 별도 세션 fixed-SHA 코드 리뷰가 필요하다. critical/high는 coor에 전달하고 수락을 차단한다. 새 제품 수정은 DEV가 수행한다. 작성자와 다른 리뷰 세션을 확보하지 않았으면 독립 코드 리뷰 PASS를 쓰지 않는다.

## 갱신할 산출물

QA 시나리오·보고서·packet-outcomes와 자기 phase만 작성한다. D03/D05–D10/D12/D13의 내용 일치와 영향은 검사하되 기술 문서/제품 코드는 수정하지 않는다. D07은 DB 설계, 사용자설명서는 D11이다. 제품 의미 충돌은 coor를 통해 designer에 전달한다.

## 적용 기준과 예외

fullops-common-0.3.3의 README·coding-style·testing·security, FULLOPS.md, orca-agents.md, project.md, docs/agents/document-writing.md를 직접 읽는다. 기준 ref는 f20b33a11addeb3debcb23e2e835be31c7dd1454다. 테스트 lite·하위 위임 standard를 유지한다. 제품 변경 근거는 2026-10-12 사람 peer·이름·대화 및 frontend-design 전면 개편 사용자 요청이다. 보안·필수 검사·독립 리뷰 예외는 없다.

## 먼저 읽을 문서

필수: docs/planning/product-specs/SAR-HUMAN-PEER.md, docs/design-docs/mockups/service-member-ui.md의 현행 절, docs/planning/product-specs/SAR-PUBLIC-SERVICE.md의 현행 안내, docs/exec-plans/phases/SAR-HUMAN-PEER-001.md, docs/operations/user-guide.md의 새 흐름, project.md와 공통 규칙. 경로는 .fullops-squad/ 기준이다.

코드 시작점: internal/relay/member_agents.go, member_pages.go, identity.go, public_text.go, test_messages.go, store.go, member.go. 필요에 따라 현재 인가의 모든 호출자, adapter local/remote·SSE·gate 경로를 일반 검색으로 보완한다. test_messages.go는 시험 경로다. public_text.go가 일반 signed text 경로다. 존재하지 않는 member_messages.go를 전제로 삼지 않는다. 라이브러리 API/설정 확인은 Context7 resolve-library-id 다음 query-docs를 사용한다. 문서·패킷 partial/unknown을 부재로 해석하지 않는다.

## 지시 전제와 충돌 — 먼저 확인

과거 원천의 Contacts=agents·agent별 첫 수락·composer 금지는 현행 HP/D04가 대체한다. 업무 relay.v1·서명·현재 인가·인간 게이트는 유지한다. 새 대화 text의 24h 보존만 기존 완료 뒤 원문 삭제를 바꾼다. 과거 메시지·업무 M/H/R을 소급 보존하지 않는다. 새 사람 친구 권한이 기존 pair를 자동 확대하지 않도록 HP 이행 표를 적용한다. 과거 리뷰/QA PASS는 새 기능의 증거가 아니다.

## 협업과 실행 제약

삭제·이동 정리는 금지한다. 임시 검증 자료는 보존하고 경로를 보고한다. 실제 이메일/Google 계정·외부 상대 발송·운영 배포·DNS/신뢰 설정·비밀값 출력은 이 지시의 실행 범위가 아니다. 실계정 확인이 필요한 항목은 로컬 합성 검증과 분리해 인계한다. 이미 위임된 제품 세부 기준은 반복 승인받지 않는다. 기술 판단은 DEV가 처리하며 설명되지 않은 제품 충돌만 coor에게 ask한다.

선택형 하위 작업은 rules/delegation.md를 읽고 Purpose: research|review|implementation, 테스트 lite와 파일 소유권을 선언한다. 최대 두 개이며 부모가 회수한다. 독립 리뷰를 자기 하위 자기 리뷰로 대체하지 않는다.

Windows make 부재 참고: coordinator msg_907ca28147c9가 C:/Users/shin/AppData/Local/KnowsLinkDevTools/frontpolish001-coor/remote.py 및 make.cmd의 기존 Linux 검증 경로를 안내했다. ROOT가 coor이므로 그대로 실행하지 않는다. 자기 대상 SHA/ROOT·known hosts·비밀값 미출력과 실제 동기화 범위를 확인한 후 기술 담당자가 재사용 여부를 결정한다. 과거 helper 존재는 새 후보의 검사 통과가 아니다.

## UI 선택과 검증 계약

제품 디자인 정본은 docs/design-docs/mockups/service-member-ui.md의 현행 SAR-HUMAN-PEER-001 절이다. 사용자 선택 frontend-design, codex / gpt-6.1-sol / high를 유지한다. 다른 호스트는 coor가 실제 지원 tuple을 명시하기 전 자동 전환하지 않는다. 새 도구 설치 없음. 기존 Go templates/CSS·Orca 브라우저를 재사용한다. memberStyle은 기존 공용 스타일 시작점이며 현행 D04 토큰으로 바꾼다. DESIGN 검사가 Go 문자열 CSS를 판정하지 못하는 한계를 기록한다. 다크 테마는 범위 밖이다.

장식 카드 묶음·긴 기술 설명·ID 위주 홈·상시 보안 배너·그라디언트·알약 버튼·가짜 대화/온라인을 만들지 않는다. 보안/보존/실패 안내는 실제 동작 옆에 유지한다. 320/390/1440폭, 44px 동작, 포커스, 키보드, 이름 중복/긴 이름, 한국어 줄바꿈을 확인한다.

웹 QA는 자기 loopback origin·고유 포트/fixture/DB·정확한 허용 GET/명시 POST·독립 assertions를 기록한다. 운영 origin/외부 발송은 금지한다. 권장 예산은 핵심 동작 시나리오 120 steps·20분이며 예산 중단은 not_run이다. 필수 보안 검사는 예산으로 생략하지 않는다. E2E/Orca 실행기 버전, 요청/시나리오 지문과 실제 완료 범위를 보고한다. 탐색·shadow·CUAWright는 OFF다. 캡처는 시각 항목만 만든다.

## 완료와 복귀

진행은 자기 인박스 한 곳에 체크한다. 상세 로그와 packet-outcomes에 모든 비선택 path/category의 completed/no_change와 구체적 이유를 쓴다. partial 잔여를 수동 검토한다. 결과·검증·미검증·후속을 완료 보고에 전문으로 기록하고 커밋·work.py finish를 수행한다. task/dispatch/terminal은 실제 주입 preamble을 따른다. worker_done은 3문장과 최종 SHA 하나로 정확히 한 번 보낸다. coor가 고정 후보 리뷰·main 통합·origin push를 조정한다.

## 완료 보고

미완료. 고정 SHA·검증·미실행·후속을 실제 완료 때 작성한다.
