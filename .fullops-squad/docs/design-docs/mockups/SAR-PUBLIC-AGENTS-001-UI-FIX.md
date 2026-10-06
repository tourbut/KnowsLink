---
title: SAR-PUBLIC-AGENTS-001-UI-FIX — 수정 후보 직접 시각 재검수
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-AGENTS-001-UI-FIX]
summary: 수정 후보의 모바일·오류·관계·기록 포화 안내 직접 시각 판정과 원본 실패를 보존한다
---

# SAR-PUBLIC-AGENTS-001-UI-FIX — 수정 후보 직접 시각 재검수

## 판정과 기준

지시된 **좁은 UI 재검수는 PASS**다. F-UI-01–04의 수정 화면과 POLICY의 새 안내를 실제 로컬 브라우저에서 직접 확인했다. 이 판정은 고정 후보의 아래 시각 조건에만 적용한다. 보안 리뷰·독립 TESTER QA·main 통합·일반 서비스 공개 수락은 각각 별도다. 이번 시각 범위에서 새 critical/high 또는 미해결 medium을 발견하지 않았다.

제품 고정 SHA와 lint 기준 ref는 `458798c2ee15c179edacfd6f94ebb9896d26f411`이다. 기록 시작 HEAD는 `f0b69a97ba2fdcb070111ff95aca0a42ae0e3f2a`다. `6d016e5`와 고정 후보의 전체 diff는 `.fullops-squad/` 기록에만 있다. 제품 경로 diff는 exit 0이다. 새 detached clone의 추적 제품 트리는 실행 전후 동일하다. 제품 코드·정책 수치·기술 정본을 바꾸지 않았다.

적용 기준은 `fullops-common-0.3.3`의 README/coding-style/testing/security, FULLOPS, project, document-writing, review/rule.json, D02 PS-07/PS-07-I·기록 보호 조건, UX-04/05와 [POLICY의 두 관찰표](../../exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md)다. DEV-FIX·DEV-POLICY-FIX와 원본 OPS·QA·probe-failures·UI 보고서를 읽었다. Jev 제외 추천은 필수 원문을 제외하는 근거로 쓰지 않았다.

원본 제품 `d1eef9bb90b9726149980320c42fb1fdbcaf584a`의 [UI FAIL/보류](SAR-PUBLIC-AGENTS-001-UI.md), 원본 UI 기록 `d1651784c4338efeb0d6141467d563c6b354e4a5`, PNG·Orca blank/1px·QA 원실행·OPS finding을 그대로 보존한다. 관련 파일 65개의 전후 SHA256이 동일하다. [보존 대조](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/original-preserved.json)는 원본을 최신 후보의 PASS로 바꾸지 않았음을 확인한다.

## 격리 실행과 직접 확인

실행 clone은 `/tmp/knowslink-ui-fix-458798c-33ddf7/src`다. Compose project는 `knowslink-ui-fix-33ddf7`이며 새 Postgres 17의 loopback port는 `44603`이다. 자기 host relay는 `56787`, 자기 SMTP sink는 `57199`다. 실제 주소 `http://localhost:56787/`에서 검수했다. 운영 DB·공유 컨테이너·기존 .env·Tunnel·서버·외부 계정은 사용하지 않았다. synthetic signup과 시험 allowlist를 켜지 않았다.

일반 이메일 코드 경로로 합성 회원 A와 B를 각각 로그인했다. 실제 Secure cookie의 별도 비영속 context를 사용했다. A의 두 agent로 같은 owner 수락을 확인했다. B의 agent로 독립 owner의 받은/보낸 초대를 확인했다. A의 교체 agent는 실제 Node prepare·브라우저 지문 확인·승인·complete를 따로 수행했다. CLI prepare와 complete는 두 연결에서 각각 exit 0이었다. 개인키·credential은 자기 임시 폴더에만 보관했다.

자기 worktree의 처음 탭 목록은 비어 있었다. 새 Orca page `be15fe73-240b-4690-b4e3-e0006bcc1f16`만 만들었다. viewport PNG는 1×1이었다. 전체 PNG는 1280×772였지만 직접 열었을 때 배경뿐이었다. 두 장애 PNG를 새 폴더에 보존하고 시각 PASS 근거에서 제외했다. 원본 과제의 장애 PNG도 유지했다.

허가된 fallback은 기존 설치 Playwright `1.60.0`과 Chromium `147.0.7727.15`다. 새 package나 frontend를 설치하지 않았다. 직접 지정한 클릭·입력·선택으로 실제 제품 페이지를 조작했다. 1280×900과 390×844에서 정상 PNG 62개를 확보했다. 이 PNG 62개와 장애 PNG 2개를 모두 직접 열어 판정했다. 정지 상태이므로 영상은 필요하지 않았다.

시간 경계와 기록 포화는 자기 relay를 멈추고 자기 DB만 바꾼 fixture다. waiting Exp·pending Exp는 1초 전, Verified는 6분 전으로 설정했다. 키 기록 20개는 기존 키를 유지한 채 25h 이상 된 철회 공개키 fixture를 추가했다. owner 기록 10개는 1h 된 철회 agent fixture를 추가했다. 정리 후 관측은 철회 agent의 Changed만 25h 전으로 보내고 실제 제품 sweep을 실행했다. SQL 명령은 모두 exit 0이다. 실제 10분/24h 경과를 검증한 것은 아니다. [fixture 변경 전문](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/fixture-mutations.txt)과 [비밀값 없는 상태 관찰](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/state-observations.json)을 연결한다.

## F-UI-01–04 재검수

아래 모든 새 관측의 제품 SHA는 위 고정 40자리 SHA다. PNG 번호는 [이번 캡처 폴더](SAR-PUBLIC-AGENTS-001-UI-FIX/)의 파일명 앞 두 자리다. 각 실제 파일·SHA256·크기·viewport·fixture·직접 열람·PASS/캡처 FAIL은 [manifest](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/evidence.json)에 있다.

| 항목 | 이번 실제 관측과 판정 | PNG |
|---|---|---|
| F-UI-01 지문·긴 식별자 | PASS. prepared의 공개 지문은 CLI와 일치한다. 모바일의 홈·키·prepared 지문이 카드 안에서 줄바꿈된다. document scrollWidth/viewport는 390/390이다. 활성/철회 키 20개인 카드에서도 안내가 읽힌다. | 01–05, 16, 40, 43–44 |
| F-UI-02 미지원·권한 거부 | PASS. 자기 form에 unsupported option 또는 없는 agent ID를 넣어 안전한 실제 거부를 확인했다. 각 화면의 홈 링크를 클릭해 관리 화면으로 복귀했다. 정상 제품 select에 미지원 선택지가 있는 것은 아니다. | 10–13 |
| F-UI-02 재확인·로그인 복귀 | PASS. 오래된 Verified에서 거부와 재확인 버튼이 함께 보인다. 버튼으로 실제 이메일 코드를 받아 확인한 뒤 같은 자기 홈으로 돌아왔다. 세션 없는 별도 context에서는 로그인 링크를 제공했다. Tab의 focus와 실제 로그인 화면 복귀도 확인했다. | 14–15, 61–62 |
| F-UI-03 종료 연결 취소 버튼 | PASS. cancelled·expired·consumed에는 연결 취소 버튼이 없다. 홈 복귀는 남는다. prepared에는 승인·취소가 있어 가능한 동작을 구분한다. | 01–03, 06–09 |
| F-UI-04 KST·select | PASS. 관계 pending·종료 기한은 연결 기한과 같은 KST 형식이다. raw UTC 문자열이 없다. select는 본문과 같은 18px이며 390px에서 카드 안에 들어간다. | 05, 16, 18, 20, 23, 29, 34 |

이 표는 원본 FAIL 기록을 수정하지 않는다. 수정 후보에서 F-UI-01/02 medium의 시각 원인이 해소됐다는 새 관측이다. 서버의 전체 인가·rate·CSRF·모든 API 차단을 이미지로 증명하지 않는다.

## POLICY 두 표의 UI 관찰 연결

POLICY 첫 표의 상태·방향·다음 동작과 두 번째 표의 기록 보호·교체·정리 안내를 직접 검수했다. 이 표의 PASS는 UI 관찰이다. API 전송·세대 경합·receipt-only replay·모든 rate/용량 수치·재시작 권한의 독립 증명은 TESTER/OPS 담당이다.

| 원천 조건 | fixture와 직접 관측·판정 | PNG |
|---|---|---|
| pending 같은 방향 반복 | PASS. A→B의 실제 첫 초대와 반복 제출이다. 새 초대·기한 변경 없음 notice가 보인다. 보낸 초대에는 수락/거절이 없고 수신 owner 대기를 안내한다. | 17–20 |
| pending 반대 방향 반복 | PASS. B가 같은 두 agent로 반대 방향을 제출했다. 같은 세대 1·같은 KST 기한이 유지된다. notice는 받은 초대의 결정 위치를 가리킨다. B 화면에만 수락/거절이 있다. | 21–23 |
| active 반복 | PASS. B의 실제 수락 뒤 A가 다시 제출했다. 이미 연결됨·새 초대 없음 notice와 같은 active 세대 1이 보인다. | 24–27 |
| 거절·만료·양측 철회 | PASS. 실제 수신 거절·송신 철회·수신 철회를 각각 실행했다. 만료는 제어된 pending Exp다. 종료와 메시지 불가·새 초대 버튼을 보인다. | 28–29, 31–35, 38–39 |
| 종료 뒤 수동 재초대·새 수락 | PASS. 종료 카드의 새 초대 버튼을 직접 눌렀다. 새 pending과 새 수신 owner 대기가 보인다. 다시 수락하기 전 active로 표시하지 않는다. | 30, 34, 36, 46 |
| 같은 owner의 두 agent | PASS. 실제 A의 두 agent도 첫 초대는 pending이었다. 수신 owner 수락 버튼을 누른 뒤 active로 바뀌었다. | 36–37 |
| live agent 키 기록 포화 | PASS. 20개 기록의 홈에서 새 키 연결 불가와 키 철회/대기만으로 공간이 생기지 않음을 안내한다. 실제 포화 connect 거부는 기존 관리와 새 agent 버튼을 제공한다. | 40–43 |
| 키 철회와 오래된 기록 | PASS, 제어 fixture 한정. 철회 key fixture는 이미 25h 이상이다. 실제 선택 키 철회 뒤에도 live agent 기록은 20개이고 포화 안내가 남는다. 새 kid 재할당이나 시간만 기다리는 우회를 안내하지 않는다. | 44 및 상태 관찰 |
| 새 agent 교체와 별도 관계 | PASS. 포화 오류의 새 agent 버튼으로 별도 ID를 만들었다. 정상 prepare·지문 승인·complete 뒤 연결 완료다. B와 새 관계는 별도 세대 1 pending이며 B의 새 수락 뒤에만 active다. 기존 관계를 승계하지 않는다. | 45–46, 56 및 상태 관찰 |
| 활성 한도와 owner 기록 포화 | PASS. 실제 활성 agent 5개 생성 뒤의 추가 생성 거부와, owner 기록 10개 fixture의 거부를 따로 확인했다. 상품 quota·결제 안내가 없다. 기존 관리 복귀가 있다. | 47–50, 57–58 |
| 키와 owner가 모두 포화 | PASS. 새 agent도 지금 생성 불가라고 안내한다. 불가능한 새 agent 생성 버튼은 없다. 최소 24h 보존·실제 홈 목록 정리 뒤 재시도를 안내한다. | 51–52 |
| 포화 중 철회·최소 보존 | PASS, UI fixture 한정. owner 기록 10개에서 실제 agent 철회는 됐다. 활성 slot이 줄어도 기록은 10개라 새 생성은 여전히 거부된다. 1h fixture와 방금 철회한 agent가 홈에 철회 상태로 남는다. | 53–55 및 상태 관찰 |
| 실제 정리 뒤 목록·수동 재시도 | PASS, 제어 fixture 한정. 철회 시각을 25h 전으로 보낸 뒤 실제 sweep이 옛 agent·관련 pair를 지웠다. 홈에서 옛 ID가 사라졌다. 사용자가 새 agent 버튼을 눌러 생성했다. 새 agent는 미연결이며 기존 교체 관계만 유지된다. | 56, 59–60 및 상태 관찰 |
| 철회 전 경고·사라짐의 의미 | PASS. 철회 버튼 앞에 모든 키·관계 종료와 복구 불가를 알린다. 철회 후에는 최소 24h 뒤 정리·목록에서 사라질 수 있음·권한 복구/백업 영구 삭제 아님·새 연결/새 수락을 안내한다. | 05, 16, 40, 43, 53–54, 60 |

35번 PNG는 B가 새 초대 발신자가 된 뒤 B가 철회한 화면이다. 이를 현재 수신측 철회로 오기하지 않도록 파일명을 `35-reinviter-revoked-mobile.png`로 바로잡았다. 실제 현재 수신 B 철회는 A→B 세대 5의 38–39번이다. 47번의 최초 캡처 viewport는 390×844여서 파일명을 `47-active-agent-limit-mobile-initial.png`로 바로잡았다. 별도 실제 desktop 활성 한도 화면은 57번이다. 두 수정은 새 증거의 이름 보정이며 원본 실패를 삭제한 것이 아니다.

## 마스킹·정리·변경 없는 범위

이메일·code 입력·43자 연결 수단은 PNG 생성 전에 불투명 mask로 가렸다. 화면의 이메일 문단과 신원 dd도 가렸다. 공개키 지문과 합성 agent ID는 식별·비교에 필요한 공개값이라 유지했다. PNG를 후처리하거나 HTML을 모의 렌더하지 않았다. CLI 토큰 입력과 이메일 code 입력은 로컬 비공개 경로로만 처리했다.

자기 Orca page·비영속 Chromium context·SMTP·relay와 자기 Compose container/volume/network를 회수했다. 자기 합성 메일·클라이언트 개인키/credential·fixture DB password 파일도 제거했다. 기존 KnowsLink relay/cloudflared/postgres의 ID는 시작 관측과 같다. [cleanup.json](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/cleanup.json)과 원 명령 exit 0을 보존했다. 종료 뒤 자기 worktree 탭은 다시 비어 있다.

범위 밖 가입·로그아웃·회전·human-gate·신원 전체 UI는 이번 SHA에서 다시 PASS로 선언하지 않는다. memberStyle이 바뀌어 시각 의존성이 동일하지 않으므로 원본 d1의 PNG를 새 후보 시각 근거로 재사용하지 않았다. 기존 증거는 원래 SHA의 원실행으로 유지한다. 소스가 같은 Node CLI 정상 연결은 이번 좁은 흐름의 실제 준비 수단으로 실행했다.

## 검증·후속·한계

별도 theme·Tailwind·shadcn·디자인 lint는 없으며 해당 없음이다. FullOps DESIGN은 Go 문자열 CSS를 검사하지 못한다. 자동 DESIGN 통과를 모바일·시각 수락으로 해석하지 않았다. 기존 memberStyle/Go template 화면만 검수했다. 새 디자인 자산·제품 코드·의존성 변경은 0이다.

지시서의 D03 UX 표기는 실제 매핑과 다르다. D03은 DEV 소유 아키텍처이므로 수정하지 않았다. designer 소유 UX 원천과 실제 D04 화면 원천에 이번 근거 링크를 추가했다. 제품 요구와 디자인 방향은 바꾸지 않았다.

strict·공백·제품 무변경·원본 65개 보존·PNG/JSON/로컬 링크 검사와 깨끗한 기록 HEAD의 FullOps `--from 458798c2ee15c179edacfd6f94ebb9896d26f411` 결과는 이번 [검증 폴더](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/)와 [실행 기록](../../exec-plans/phases/SAR-PUBLIC-AGENTS-001-UI-FIX.md)에 보존한다. 등록 product-lint/product-test와 lint 원 명령의 종료코드를 따로 남긴다. 파이프로 종료코드를 가리지 않는다.

실제 이메일·운영 공개/배포·실제 24h 운영 경과·운영 데이터 정리·부하/복원·노우↔다닷·앱/OAuth·실메시지·스크린리더 전체·OS clipboard는 미검증이다. 기존 PS-13·운영 공개 held와 원본 QA/UI의 실행 SHA를 유지한다. 남은 담당은 고정 후보 독립 TESTER QA·OPS delta 리뷰/운영 조건 및 coor 통합이다. 새 UI 코드 수정 요청은 없다.
