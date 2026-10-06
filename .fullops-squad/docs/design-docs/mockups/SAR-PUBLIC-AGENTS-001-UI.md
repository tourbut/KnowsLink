---
title: SAR-PUBLIC-AGENTS-001-UI — 일반 회원 연결·키·관계 직접 시각 검수
status: review
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-AGENTS-001-UI]
summary: 고정 d1eef9b의 실제 연결·키·관계 화면과 모바일·오류 UX 실패 및 미검증을 기록한다
---

# SAR-PUBLIC-AGENTS-001-UI — 일반 회원 연결·키·관계 직접 시각 검수

## 판정과 범위

검수 작업은 완료했다. UX04–05 전체 시각 수락은 **FAIL/보류**다. 연결·키·관계의 지정 상태는 실제 localhost 브라우저에서 확인했다. 모바일 지문의 가독성과 오류 뒤 다음 동작이 필수 UX 기준을 충족하지 않는다. 이번 시각 관측에서 critical/high는 발견하지 않았다. 이 결과는 보안 리뷰나 독립 TESTER QA를 대신하지 않는다.

대상 제품 SHA는 `d1eef9bb90b9726149980320c42fb1fdbcaf584a`다. 기준 ref는 `d2f7ba5aeb6644fd2b27fe6ada5b61db5976933b`다. 기록 시작 HEAD는 `d111fd4fde1538492c3c335e78ed510408718b2e`다. 제품 파일과 제품 규칙을 바꾸지 않았다. 최종 기록 HEAD는 이 과제의 worker_done에 고정한다.

정본은 [일반 서비스 PS04–07·PS11·개인정보](../../planning/product-specs/SAR-PUBLIC-SERVICE.md), [UX04–05](SAR-PUBLIC-SERVICE-UX.md), [D04](SAR-MVP-001-UI.md), [DEV 인계](../../exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md)다. DEV의 자동 PASS를 직접 시각 PASS로 바꾸지 않았다. Jev의 이 과제 code/documents-find와 context를 읽었다. omit 추천 `adapters/README.md`는 설치 위치 확인에 사용했다. 민감/크기 제외 `connect.ts`는 비밀값을 출력하지 않고 직접 읽었다.

## 환경과 실제 관측 방법

고정 SHA의 새 clone `/tmp/knowslink-agents-ui-d1eef9b`를 사용했다. Compose project는 `knowslink-agents-ui-d1eef9b-8398`다. 새 Postgres 17의 loopback port는 `57187`, host relay는 `59219`, 자기 SMTP sink는 `52591`이다. 실제 주소는 `http://localhost:59219/`다. SMTP는 합성 메일만 0700 폴더·0600 파일에 보관했다. 기존 .env·사용자 인증·운영 DB·Tunnel·서버·공유서비스를 사용하거나 변경하지 않았다. synthetic signup과 시험 allowlist도 사용하지 않았다.

일반 가입·재로그인·재확인은 실제 이메일 코드 경로와 Secure cookie로 실행했다. 회원 A는 별도 agent 두 개를 연결했다. 회원 B는 교차 회원 관계 확인용 agent 한 개를 연결했다. 합성 신원과 로컬 자격은 외부 계정 신원이나 노우·다닷 연결을 증명하지 않는다.

Orca의 새 page `5821c251-81fd-480c-86f3-996591415053`만 제어했다. 작업 시작의 자기 worktree 탭 목록은 비어 있었다. 다른 worktree의 사용자 Cloudflare·계정 탭에는 goto·입력·클릭·종료를 하지 않았다. Orca viewport 캡처는 1×1 PNG였다. 전체 캡처는 1265×843 PNG였지만 직접 열었을 때 배경만 보였다. 정상 DOM snapshot과 이 캡처 실패를 구분했다.

허용된 대안으로 Playwright `1.60.0`의 `chromium.launch`를 사용했다. 실제 engine은 Chromium `147.0.7727.15`다. 이미 설치된 cache를 읽기 전용으로 사용했다. 새 비영속 context 두 개로 같은 실제 relay에 접근했다. viewport는 1280×900과 390×844다. 독립 CDP 실행은 로딩 timeout으로 폐기했다. Playwright가 직접 시작한 engine의 정상 PNG만 UI 판정에 사용했다. 새 frontend·의존성·모의 HTML은 만들지 않았다.

브라우저 클릭·입력은 검수자가 한 단계씩 지정했다. 준비·완료만 실제 `adapters/dist/connect.js`로 실행했다. CLI prepare·complete는 모두 exit 0이었다. 지문을 CLI와 실제 화면에서 대조한 뒤 브라우저 승인 버튼을 눌렀다. 지정한 PNG 31개를 `view_image`로 직접 열어 판정했다. Orca 장애 PNG 2개는 PASS 근거에서 제외했다.

만료는 자기 relay를 중지하고 전용 DB의 해당 waiting `Connections.Exp`만 현재보다 1초 전으로 바꾼 뒤 관측했다. 재확인 경계도 자기 relay를 중지하고 전용 fixture의 `Sessions.Verified`만 6분 전으로 바꿨다. SQL은 exit 0이었고 relay를 다시 시작했다. 제품 TTL·재인증 수치·소스는 바꾸지 않았다. 실제 10분 경과와 장시간 타이머 품질은 미검증이다.

## 상태별 직접 판정

PNG 경로는 모두 [캡처 폴더](SAR-PUBLIC-AGENTS-001-UI/) 기준이다. 기계 판독 manifest는 [evidence.json](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI/evidence.json)에 있다.

| 항목 | 실제 관측·판정 | PNG |
|---|---|---|
| 일반 가입·자기 홈 | PASS. 확인 전 대기, 일반 코드 로그인, 자기 목록만 표시한다. 기본 정보 구조를 유지한다. | 01-verify, 02-home-empty, 03-home-two-unconnected |
| 연결 대기·실행 위치·권한 | PASS. 대상, Node 로컬, 자기 클라이언트 실행 위치, 10분·1회, agent 송수신 한정, owner/admin 권한 부재를 표시한다. 개인키·관리자 token·공유 서버 SSH를 요구하지 않는다. | 04-connection-waiting |
| 키 준비·지문 승인 | PASS. CLI와 화면 지문이 일치한다. prepared에 승인·취소가 있다. 승인 전에는 키가 없다. | 05-prepared-fingerprint |
| 승인과 실제 완료 | PASS. approved는 새 키 미연결과 complete 명령을 안내한다. CLI 완료 뒤 consumed와 연결 완료를 표시한다. 두 agent가 각각 활성 키를 가진다. | 06-approved-not-connected, 07-consumed-complete, 08-home-two-connected |
| 취소·만료 | PASS, fixture 경계 한정. cancelled·expired를 문자로 구분하고 새 키 미연결을 설명한다. 기존 키는 유지됐다. expired에는 홈 복귀가 있다. | 13-connection-cancelled, 14-connection-expired |
| 미지원·권한 거부 | 안전한 오류 문구는 PASS. 전용 DOM에 미지원 option을 추가해 제출했고, 별도 제출에는 존재하지 않는 대상 ID를 넣었다. 정상 UI에 미지원 option이 있다는 뜻은 아니다. 오류 뒤 복귀 동작은 F-UI-02 FAIL이다. | 15-unsupported-error, 16-permission-error |
| 회전·선택 철회·agent 철회 | PASS. 회전 후 이전 키 철회·새 키 활성이다. 선택 철회 후 A는 미연결이며 키 기록은 철회다. agent 전체 철회 후 B는 철회이고 연결 동작은 사라진다. | 17-rotation-keys, 27-selected-key-revoked, 28-agent-revoked |
| 재확인 경계·오답·복구 | 경계·오류·복구 PASS. 오래된 Verified에서 키 철회가 거부된다. 실제 60초 제한은 재시도 시각을 표시한다. 오답에는 남은 시도 4회가 보인다. 올바른 코드 재확인 뒤 같은 회원으로 돌아와 키 철회를 완료한다. 복귀 동작은 F-UI-02 FAIL이다. | 25-reauth-required, 26a-reauth-rate, 26-reauth-wrong-code, 27-selected-key-revoked |
| 같은 owner의 두 agent 관계 | PASS. 첫 초대는 pending이고 명시 수락 뒤 active다. 철회는 revoked, 재초대 뒤 거절은 denied다. 가입 초대와 구분하는 설명이 있다. | 09-pair-pending, 10-pair-active, 11-pair-revoked, 12-pair-denied |
| 교차 owner의 초대와 양측 철회 | PASS. 발신 화면에는 수락 버튼이 없고 수신 화면에만 수락·거절이 있다. 수락 뒤 발신 owner가 active 관계를 철회했다. 새 세대 수락 뒤 수신 owner도 철회했다. 상대 이메일은 관계 영역에 없다. | 20-cross-sender-pending, 21-cross-recipient-pending, 22-cross-active-sender, 23-cross-sender-revoked, 24-cross-recipient-revoked |
| 수락 전 메시지 상태 | UI PASS. pending 설명은 수락 전 전송을 거부한다고 표시한다. 메시지 composer·전송 버튼은 없다. 실제 API의 전송 차단은 독립 TESTER의 담당이며 이번 시각검수로 증명하지 않는다. | 09-pair-pending, 20-cross-sender-pending |
| label·키보드·긴 ID 선택 | PASS, 관측한 동작 한정. label이 입력과 연결돼 있다. 이메일 입력 뒤 Tab으로 코드 받기 버튼, Enter로 요청한다. readonly agent ID에서 Control+A는 28자 전체를 선택한다. Tab으로 키 철회 버튼에 이동하고 focus ring을 확인했다. OS clipboard 붙여넣기·스크린리더 전체 검수는 미검증이다. | 19-mobile-agent 및 manifest의 selection 측정 |
| 모바일 지문 가독성 | FAIL. 390px viewport에서 문서 scrollWidth가 596px이다. 지문이 카드 밖으로 넘치며 카드 캡처에서 잘린다. | 18-mobile-home, 19-mobile-agent |
| 상태 문자·세션 종료 구분 | PASS. 미연결/연결 완료/활성/철회와 관계 state를 문자로 표시한다. 로그아웃 안내는 agent 자격 유지다. 회원 B의 실제 로그아웃·재로그인 뒤 같은 agent와 활성 키가 남는다. | 29-logout-preserves-agent-note, 30-relogin-key-still-active |

## 발견 사항과 DEV 인계

| ID·등급 | 위치·재현·영향 | 판정·후속 |
|---|---|---|
| F-UI-01 medium | `internal/relay/member.go:17,24`. 두 agent와 활성/철회 키가 있는 자기 홈을 390×844로 연다. viewport 390에 scrollWidth 596이다. readonly 입력은 250px이며 키 지문 텍스트는 카드 밖으로 넘친다. 비교할 지문을 한 화면에서 읽기 어렵다. | UX04의 읽을 수 있는 긴 식별자·지문 기준 FAIL. DEV가 같은 상태의 모바일 지문 가독성을 수정한 뒤 designer가 직접 재검수한다. 기술 수정 방법은 DEV 책임이다. |
| F-UI-02 medium | `internal/relay/member_agents.go:182–189`, `internal/relay/member.go:20`. 재확인 필요·미지원·권한 거부를 제출한다. 결과는 제목과 alert뿐이며 홈 링크·재확인·다시 선택할 컨트롤이 없다. 문구는 홈이나 선택을 요구하지만 화면 안에서는 진행할 수 없다. 브라우저 Back 또는 주소 입력으로만 복귀했다. | UX04의 실패 뒤 다음 동작 기준 FAIL. 실제 안전 거부는 유지됐다. DEV가 오류 화면에서 안내한 다음 동작을 제공한 뒤 designer가 재검수한다. |
| F-UI-03 low | `internal/relay/member.go:23`. cancelled·expired 화면에도 연결 취소 버튼이 남아 있다. 만료/취소 완료 상태에서 할 수 없는 동작으로 읽힌다. 서버의 반복 취소 결과는 이번에 실행하지 않았다. | 13·14 PNG의 버튼 관측만 기록한다. DEV가 완료 상태의 가능한 동작을 검토한다. |
| F-UI-04 low | `internal/relay/member.go:24`. 관계 기한이 `2026-10-07 04:19:59.162233 +0000 UTC`처럼 raw 문자열이다. 연결 기한의 KST 표현과 달라 일반 회원이 비교하기 어렵다. native select는 13.3333px, 본문은 18px이다. | 관계 기한·컨트롤 가독성 후속. 정책 수치나 권한 규칙을 변경하지 않는다. |

미해결 medium 두 건 때문에 해당 필수 UX 항목을 PASS로 표시하지 않았다. DEV의 기능 검사 PASS와 이 시각 FAIL을 함께 인계한다. 독립 보안 finding·TESTER 결과·운영 공개 held는 별도로 유지한다. 새 high가 없다는 관측으로 다른 검토의 high를 닫지 않는다.

## 마스킹과 보존

화면 캡처 전에 Playwright screenshot의 불투명 `mask`로 이메일 입력·확인코드 입력·연결 수단 입력을 가렸다. 이메일이 들어 있는 문단과 홈의 신원 dd도 가렸다. 연결 수단의 실제 값은 Node prepare의 비공개 stdin으로만 보냈다. 코드·메일·cookie·credential·개인키를 터미널이나 Git에 출력하지 않았다. private key와 agent.json은 0700/0600의 자기 fixture에만 있었다.

공개키 지문은 비교의 필수 공개값이라 PNG에 유지했다. agent ID·회원 ID는 이메일에서 파생되지 않은 합성 식별자다. 관계 section과 agent section 캡처에는 이메일·코드·grant가 없다. PNG는 실제 브라우저가 생성했으며 후처리로 UI를 그리거나 합성하지 않았다. 임시 마스킹은 제품 소스와 CSS 정본을 변경하지 않았다.

사용자 기존 탭과 자격은 보존했다. 자기 Orca page·Playwright 프로세스·SMTP·relay와 자기 Compose DB/volume/network만 회수했다. 합성 메일·클라이언트 키·fixture 자격 파일도 자기 임시 경로에서 제거했다. 정리 결과는 [cleanup.json](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI/cleanup.json)에 있다.

## 검증과 미검증

[증거 manifest](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI/evidence.json)는 PNG 크기·SHA256·engine·viewport·fixture를 기록한다. [fixture-mutations](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI/fixture-mutations.txt)는 시간 경계의 설정을 기록한다. [제품 무변경](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI/product-unchanged.txt)은 고정 후보 대비 제품 경로 diff exit 0이다. 실제 strict·공백·FullOps 결과는 [검증 폴더](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI/)에 보존한다.

별도 theme·Tailwind·shadcn·디자인 전용 lint는 없다. 해당 없음이다. FullOps DESIGN 규칙은 Go 문자열 안의 CSS를 검사하지 않으므로 모바일 F-UI-01을 탐지하지 못한다. 실제 PNG 관측을 생략할 수 없다. 새 frontend와 디자인 의존성은 추가하지 않았다. 채팅 버블·composer·장기 타임라인·무조건 성공 배너는 없다.

기존 UX01–03 증거는 이번 실행 PASS로 재사용하지 않았다. `memberStyle`·start/verify template은 기준 diff에서 동일하지만 home template과 데이터는 변경됐다. 변경된 홈과 이번 흐름의 가입·오답·재로그인은 직접 관측했다. 과거 59b66ad의 신원 전체 QA/시각 근거는 원래 SHA의 [기존 기록](../../exec-plans/phases/SAR-PUBLIC-IDENTITY-001-UI.md)으로만 유지한다.

실제 일반 이메일 로그인·공개 서비스·운영 배포·Grok Bot 노우와 OpenAI dot 다닷·앱 설치/권한·OAuth·실메시지 왕복·자동 wake는 미검증이다. 24시간 초대 만료·실제 10분 타이머·오프라인·큰 수용량·동시성·모든 API 철회 경계·모든 세션 종료·스크린리더·OS clipboard는 이번 직접 시각 범위에서 미검증이다. 독립 TESTER·OPS와 사용자 실제 확인 조건은 그대로 남는다.
