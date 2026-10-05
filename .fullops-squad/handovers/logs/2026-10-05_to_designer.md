---
title: designer 완료 기록
status: draft
updated: 2026-10-05
owner: designer
summary: 지시서와 완료 보고를 보존한다.
---

## SAR-PUBLIC-SERVICE-001 — 2026-10-05

---
title: 일반 이메일 서비스 제품 기준 확정
status: draft
updated: 2026-10-05
owner: coor
tasks: [SAR-PUBLIC-SERVICE-001]
summary: 일반 사용자 가입과 실제 에이전트 메시지 서비스의 수락 기준 및 기능별 인계를 정한다
---

# SAR-PUBLIC-SERVICE-001 — 일반 이메일 사용자 서비스와 후속 dots 연결

## 목표와 원문

사용자: “이제 해당 테스트를 오너 아이디 말고 일반적인 메일 아이디로 진행해보자.” 이어 “일반서비스가 가능한 수준까지 하고 openai dots 연결하는거야.” 일반 사용자가 이메일 신원을 확인하여 가입하고 자기 에이전트를 연결하며 관계 수락 후 실제 메시지를 송수신하는 서비스까지 완성한다. 이를 수락한 뒤 OpenAI dots 연결로 이어간다. 기존 합성 가입과 시험용 Service Auth를 일반 사용자 로그인 완료로 표시하지 않는다.

## 현재 근거

기준 SHA `6c0d132d26dac22693f1c1ffb986fc4f61bdbd09`. 공개 `link.knowslog.com`은 현재 오너 전용 Cloudflare Access로 보호된다. `internal/relay/http.go`의 POST /v1/owners는 합성 가입이고 이메일 신원 검증이 아니다. 기존 제품 D02에는 누구나 가입 가능한 인증 파일럿 방향이 있지만 DEC-03 제안 수치와 실제 신원 정책이 미확정이다. 실제 Codex↔공식 xAI Grok Bot 수동 CLI 시험은 성공했고 시험 전용 token2·앱·policy는 삭제했다. Grok도 시험 파일 삭제 완료를 보고했다. 이전 기록과 held는 보존한다.

## 해야 할 일

- [x] 기존 제품 원천·D01/D02·백로그·운영 기록을 읽고 일반 서비스 목표와 기존 미구현/held의 차이를 확인한다.
- [x] 이메일 가입·로그인·로그아웃·재로그인, 회원별 에이전트 소유권, 키 등록/회전/철회, 관계 초대·수락·철회, 실제 양방향 메시지 및 실패 상태의 사용자 완료 조건을 정한다. 일반 사용자에게 관리 오너 권한이나 서버 파일 배치를 요구하지 않는 연결 흐름을 정한다.
- [x] 일반 서비스에 필요한 신원·데이터 분리·최소 남용/자원 한도·운영 복구·상태 표시 기준을 제품 책임 범위에서 정한다. 기존 사용자 승인 전 제안은 승인됐다고 소급하지 않는다. 서비스 목표에 필요한 합리적 기본안을 근거와 함께 정하고, 반드시 사용자만 결정할 항목만 구분한다. 기술 API/DB/구현 선택은 DEV/OPS에 맡긴다.
- [x] 실제로 검증 가능한 기능 단위 순서를 정하고 DEV/OPS/TESTER 인박스를 준비한다. 후보 SHA 독립 QA·읽기 전용 fixed-SHA 별도 세션 리뷰·일반 이메일 직접 로그인 UI 검수·기존 오너/공유 서비스 회귀를 포함한다. 다른 역할 인박스가 사용 중이면 PLANS 대기용 제안에 남기며 덮어쓰지 않는다.
- [x] dots는 일반 서비스 수락 뒤 연결한다. 공식 문서 https://learn.chatgpt.com/docs/dots 및 computers-and-apps를 확인하고 플러그인/앱 연결 또는 지원되는 컴퓨터 경로의 확인 항목을 남긴다. 사용자는 자기 OpenAI dot의 이름이 “다닷”이라고 확인했다. 해당 dot 접속·플러그인 설치 화면은 아직 미확인이다. 공식 inbound API·자동 wake를 추정하지 않는다.
- [x] D01/D02·백로그·실행 기록을 갱신한다. D03 기술 설계는 직접 수정하지 않고 DEV 인계에 기록한다.

## 적용 기준과 예외

`fullops-common-0.3.2`, `.fullops-squad/FULLOPS.md`, `orca-agents.md`, `project.md`, 문서 작성 규칙, coding-style/testing/security 세 규칙을 적용한다. 사용자 지시는 제품 범위 확대 및 완성 작업의 승인이다. FullOps 업데이트는 제외한다. 지금 과제는 제품 문서와 인계이며 코드·배포·CF 설정을 변경하지 않는다. 이메일은 질문 중이며 답을 기다리지 않고 제품 작업을 진행한다. 이메일·인증코드·키·secret을 Git·지시서·공개 이슈에 기록하지 않는다. 결제·일정 실데이터·임의 업무 외부 발송은 범위에 추가하지 않는다.

## 먼저 읽을 문서

Jev documents-find/context를 실행했다(정상·fallback 없음). keep과 필수 문서는 아래와 같다: `.fullops-squad/docs/planning/product-specs/SAR-MVP.md`, `.fullops-squad/docs/planning/SAR-MVP-backlog.md`, `.fullops-squad/docs/operations/ops-guide.md`, `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-ACTUAL-TRIAL.md`, `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-TRIAL-CLEANUP.md`, `.fullops-squad/docs/exec-plans/phases/SAR-MVP-PUBLIC-POLICY-001.md`.

### 지시 전제와 충돌 — 먼저 확인

PUBLIC-POLICY-001은 누구나 가입 공개 방향은 확정했으나 수치 승인은 대기였고 이전 실제 벤더 연결은 제외했다. product.md의 agent-centric 등록은 사용자를 에이전트 owner로 바인딩한다. 이번 사용자 지시는 일반 서비스 완성과 뒤이어 실제 dots 연결을 명시하므로 이 범위를 새 제품 결정으로 기록한다. 숫자·실데이터 공개·업무 효과의 과거 held는 승인됐다고 소급하지 않는다. 일반 사용자도 자기 에이전트 owner이며 운영 관리자와 구분한다.

### 필요 시 확인

Jev code-find 후보 중 `internal/relay/http.go`, `deploy/knowslink/access_apply.py`, `adapters/README.md`에서 현 상태의 구현 경계만 확인한다. 나머지 code-find 후보는 구현 과제에서 기술 담당자가 판단한다.

## 소유 범위와 산출물

`.fullops-squad/docs/planning/`, `docs/design-docs/mockups/` 제품 UX 문서, 자기 실행 기록·contexts/designer.md, 비어 있는 후속 역할 인박스. PLANS/board는 coor가 관리한다. 갱신 산출물 D01/D02; D03은 DEV 후속. 구현 파일과 기술 정본은 변경하지 않는다.

## 완료 기준과 복귀

단순 조사 보고로 끝내지 않는다. 일반 서비스 수락 기준·범위·필수 선행·기능 순서와 실행 가능한 다음 DEV/OPS 지시서를 확정한다. 미확정 외부 입력은 작업 가능한 항목과 분리한다. 문서 검증을 완료하고 `work.py finish`로 전문 보존·빈 인박스를 확인한다. 커밋 SHA·보고 경로·다음 담당을 `run_8ca8bc058ab7`에 worker_done으로 직접 보낸다. 준비 SHA에서 lint ERROR0을 확인한다. 독립 검토/구현/QA/배포는 coor가 후속 배정한다.

## 착수 뒤 전달된 사용자 기준과 운영 정정

`msg_e7183c748087`은 “내가 운영버전 테스트할꺼는 grok bot 과 openai dots를 연결해서 테스트하는거야 둘다 내 이메일로 가입할꺼야.”를 전달했다. `msg_bae19c774823`은 Grok Bot 이름 “노우”를 확정했다. 최종 운영 E2E는 동일 본인 일반 이메일의 노우↔다닷이다. 두 다른 이메일을 요구하지 않는다. 각 agentID·키·credential을 분리하고 같은 owner여도 관계 첫 수락을 유지한다. 일반 서비스 기반 수락 뒤 다닷을 연결한다.

OPS 진행 중 인박스 점유 메시지 `msg_71854ecb4417`과 후속 대기 유지 지시를 반영했다. 최초 OPS SHA `2267a4a`와 정정 `0313deae0dec9af813b70ee9685e1a6d0a2b84d7`를 git show/diff로 직접 읽었다. JWT/header 전달 가능성과 relay 미사용·실서버 도달 미확인, active 관리 token·현재 쓰기 권한 미확인, 실제 자원·처리량/백업 복구 미충족을 반영했다. OPS/IdP/API 기술 선택은 DEV/OPS 책임이다.

## 완료 보고

### 변경 이유와 결정

일반 이메일 신원 확인·가입·로그인·현재/전체 로그아웃·재로그인, 같은 회원 owner 연속성, 자기 agent 소유권·키/연결 권한·관계 수락/철회, 실제 송수신과 실패 상태를 PS-01–14로 정했다. 최초 실제 메시지는 명시적 비민감 연결 확인 text와 관련 답장이다. frozen 업무 wire와 구분하며 trial 자격을 일반 회원 인증으로 재사용하지 않는다. 같은 이메일의 노우와 다닷은 별도 agentID·키·credential을 갖는다. 같은 owner 안의 두 agent에도 첫 관계 수락을 요구한다.

이번 완성 지시의 위임 범위에서 초기 운영 기본값·보존·복구 목표를 새로 결정했다. 과거 DEC-03 제안에 사용자 승인이 있었다고 소급하지 않았다. Free N/가격/slot-unit·실일정 disclosure·결제·임의 업무 외부 발송 held는 유지했다. 공식 dots·computers-and-apps·plugin auth/events를 직접 열고 지원 앱/컴퓨터·사용자별 인증·이벤트 별도 구독 확인 항목을 남겼다. 앱 설치/CLI pull을 실제 수신·자동 wake 성공으로 취급하지 않는다.

### 산출물과 인계

D01 `docs/planning/business-plan.md`, 기존 D02 `docs/planning/product-specs/SAR-MVP.md`와 새 D02 `SAR-PUBLIC-SERVICE.md`, `docs/planning/SAR-MVP-backlog.md`, 일반 UX `docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md`, 자기 실행 기록 `docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md`와 contexts/designer.md를 갱신했다. D03·제품 코드·배포·CF·원천·PLANS/board는 수정하지 않았다.

빈 DEV/TESTER 인박스에서 work.py new로 SAR-PUBLIC-IDENTITY-001-DEV ready와 SAR-PUBLIC-IDENTITY-001-TESTER waiting을 준비했다. 각 역할의 code-find/documents-find/context를 정상 실행하고 여섯 JSON을 보존했다. 필수 keep/충돌/필요 시 확인을 인박스에 반영했다. OPS 착수 인박스가 다른 체크아웃에서 사용 중인 것을 확인하고 to_ops를 수정하지 않았다. coor 지시대로 OPS 완료 후에도 후속은 실행 기록의 대기 표에 남겼다. 다음 agents/messages/OPS/전체 수락/다닷은 coor가 PLANS에 등록하고 완료 아카이브 뒤 정규 역할 인박스로 옮긴다.

### 검증한 것과 못 한 것

준비 `608fe06cf0d67d97a483bd088d9834c34049fb6f`의 detached snapshot에서 기준 `6c0d132` lint exit 0·ERROR 0·WARNING 1·실행 불가 0과 product-lint 통과를 확인했다. 첫 clean-state symlink 오류(exit 2)와 의존성 누락(exit 1)은 로그로 보존했다. 잠금 파일대로 make install(exit 0) 후 같은 SHA 검사에서 통과했다. WARNING은 기존 coor PLANS SIZE-001이며 소유권 밖 파일을 수정하지 않았다.

체크포인트 `07ccd07405135efcad0b324737681f254ee70dab`의 lint도 같은 기준에서 exit 0·ERROR 0·WARNING 1·실행 불가 0이고 product-lint가 통과했다. deliverables strict는 검사 13·미작성 0·문제 0·경고 0·exit 0, git diff 608fe06 --check는 exit 0이다. 문서 audit는 로컬 링크 54개·PS-01–14·UX-01–08·파일 소유권/보존·인박스 과제 키·Jev 자료에서 오류 0·exit 0이다. 로그는 `/tmp/SAR-PUBLIC-SERVICE-001-validation/`에 명령 자신의 종료코드로 보존했다.

제품 구현·실메일 발송·직접 로그인 UI·일반 신원 플랫폼 왕복·배포·독립 기능 QA와 리뷰는 이번 designer 문서 과제에서 실행하지 않았다. 전체 일반 서비스 PS-13과 최종 노우↔다닷 PS-14는 미완료다. 현재 자원 수용량·신원 바인딩·backup/격리 복구·미확인 권한/좌석/IdP/MCP 호환성은 DEV/OPS 후속 조건이다. 일반 이메일 사람 확인과 실제 다닷 연결 동의만 사용자 입력으로 분리했다.

### 완료 처리와 복귀

브랜치는 fullops/designer다. 이 전문과 전체 지시서를 work.py finish로 `handovers/logs/2026-10-05_to_designer.md`에 한 번 보존하고 자기 인박스를 비운다. 최종 완료 커밋의 SHA·lint/strict·전문 일치·빈 인박스·깨끗한 상태를 다시 확인해 현재 Run의 worker_done으로 보낸다. 보고 SHA는 worker_done의 고정 값이 정본이며 체크포인트를 최종 SHA로 바꾸어 쓰지 않는다. coor가 문서 검토·main 통합·다음 역할 배정을 수행한다. 이번 문서 완료는 제품 수락과 별개다.

## SAR-PUBLIC-IDENTITY-001-UI — 2026-10-05

---
title: SAR-PUBLIC-IDENTITY-001-UI — 고정 이메일 신원 후보의 UX01–03을 실제 브라우저에서 직접 검수한다
status: draft
updated: 2026-10-05
owner: designer
tasks: [SAR-PUBLIC-IDENTITY-001-UI]
summary: 고정 이메일 신원 후보의 UX01–03을 실제 브라우저에서 직접 검수한다
---

# SAR-PUBLIC-IDENTITY-001-UI — 일반 이메일 가입·세션 직접 시각 검수

- 상태 completed. 고정 제품 후보59b66ada8b36802484cc6d7e22523257b50572cc, 코드a446d89ff288c4243ad6d7f8780a778517154584, 기준94533b207b456c0560800fe30a7c90b2b5887c6e.
- 소유: 기록 checkout /home/shin/orca/workspaces/KnowsLink/fullops-coor의 이 designer 인박스·해당 logs 전문, docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-UI.md와 필요한 UI 증거, 자기 실행 기록뿐. 제품 코드·다른 기록·PLANS/board·다른 인박스·GitHub·CF·배포는 수정하지 않는다.
- 복귀 term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9, Run run_8ca8bc058ab7. 새 세션 dispatch preamble을 사용한다.

## 적용 기준과 예외

fullops-common-0.3.2, FULLOPS.md/project.md/document-writing.md/orca-agents.md, coding-style/testing/security. 제품 SAR-PUBLIC-SERVICE PS01–04·해당PS11·UX01–03과 frozen SAR-MVP C1–C5. 사용자 일반 서비스 구현 지시에 따른 직접 UI 검수다. 새로운 제품수치/화면/코드 구현이나 기술 승인을 맡지 않는다. 최종 시험은 같은 일반 이메일의 노우↔다닷이다. 이번 로컬 fixture 직접 UI를 실제 사람 이메일/운영 시험 PASS로 대체하지 않는다.

## 먼저 읽을 문서

필수 규칙과 제품/UX 정본, DEV 실행 기록 SAR-PUBLIC-IDENTITY-001-DEV.md, README의 일반 이메일 로그인 확인 절차, 사용자 가이드, internal/relay/member.go·identity.go·http.go·scripts/mail_sink.py. 탐색 근거는 coor가 아래 추가한다.

## 지시 전제와 충돌

기존 owner-only Basic 화면은 일반 회원 가입 화면이 아니다. 이번 UX03 범위는 신원·세션·자기 owner·로그아웃이다. agent/key/pair/회원 비활성화는 다음 기능이다. 준비 중 표시는 실패로 바꾸지 말고 기능이 아직 미구현이라는 표시와 실제 성공버튼 없음 여부를 검사한다. 실제 이메일·운영 SMTP 미설정은 실제 확인만 미실행이다. fixture의 로컬 SMTP 받은 코드로 직접 브라우저 시각/흐름 검사 자체를 수행한다. 자동 QA/보안 리뷰는 별도 담당이므로 전체를 중복하지 않는다.

## 해야 할 일·완료 기준

1. 자신의 제품 실행 checkout /tmp/knowslink-public-identity-ui-59b66ad의 detached59를 확인하고 별도 임시 PostgreSQL·SMTP sink·loopback relay를 시작한다. 기존 운영 서비스/QA의 임시 환경을 변경하지 않는다. 코드를 수정하지 않는다. 필요한 실행/fixture 준비를 독립 수행한다.
2. Orca 실제 내장 브라우저에 새 전용 페이지를 열어 그 page ID를 명시적으로 제어한다. 사용자 기존 Cloudflare/다른 페이지를 탐색하거나 변경하지 않는다. orca-cli browser reference를 읽는다.
3. UX01 시작 화면의 일반 이메일 입력·로그인 절차, UX02 확인대기·마스킹·오답·만료·제한·발송 실패의 안내와 다음 동작, UX03 실제 회원 홈·빈 agent/관계 준비중·현재 로그아웃·전체 로그아웃/재확인·오래된 세션 화면을 직접 확인한다. 기존 회원 gate 영향 화면은 변경 영향에 필요한 범위만 확인한다.
4. 성공/실패의 문구·입력 레이블·가독성·레이아웃·키보드·모바일 폭을 판정한다. 테스트는 기능 판정에 필요한 최소 조작으로 한다. 실제 직접 화면 관측과 자동 근거를 구분한다. 스타일 변경을 새로 구현하지 않는다. 새 장식·별도 frontend·채팅 UI를 요구하지 않는다.
5. 필요한 화면만 캡처하고 이메일·코드·cookie·token은 캡처/보고 전에 숨긴다. dummy fixture도 값 전문을 보고하지 않는다. 보고에 환경·실제 후보SHA·브라우저page ID·경로/근거·PASS/FAIL/미실행·심각도/재현을 쓴다. 단계별 필수 실패/미해결critical/high가 있으면 수락 불가를 보고하고 임의 수정하지 않는다.
6. 자신의 임시 실행 자원을 정리하고 다른 사용자 자료를 보존한다. 이 인박스에 전문 완료 보고를 쓰고 work.py finish --role designer --key SAR-PUBLIC-IDENTITY-001-UI로 보존한다. 자신의 소유 파일만 commit하고 다른 root/reviewer 기록을 함께 add/commit하지 않는다. fixed59·결과SHA·독립관측/실제확인미실행·finding을 이 Run worker_done으로 회신한다.

## 완료 보고

실행 완료: 고정 제품 `59b66ada8b36802484cc6d7e22523257b50572cc`의 로컬 fixture UX01–03 직접 시각 검수는 PASS다. 구현은 `a446d89ff288c4243ad6d7f8780a778517154584`, 준비는 `1207bdf4542fac98d228f86de79aad9f4e126ce8`, 기준은 `94533b207b456c0560800fe30a7c90b2b5887c6e`다. 후보 checkout은 작업 전후 깨끗한 detached 상태였다.

- [x] 자기 PostgreSQL 17·SMTP sink·loopback relay를 시작했다. 후보에서 별도 바이너리를 빌드했고 운영·다른 QA 환경은 변경하지 않았다.
- [x] Orca page `56050624-50c1-4700-b313-a3ed78e71059`를 새로 만들고 모든 페이지 동작에 명시했다. 기존 Cloudflare 페이지를 조작하지 않았다.
- [x] 실제 브라우저에서 코드 발송·입력·오답·재발송 제한·만료·회원 홈·재확인·전체 종료·재로그인·현재 종료·메일 실패를 확인했다. 같은 fixture 회원 ID와 빈 목록의 연속성을 확인했다.
- [x] Orca 직접 정상 PNG 2개를 확인했다. 반복 blank를 escalation/ask로 보고했다. coordinator가 허용한 별도 실제 Chromium `147.0.7727.15`·Playwright core `1.63.0`의 새 context에서 같은 fixed59 서비스로 보완했다. 최종 PNG 15개를 직접 열어 지정 390/1280 폭·레이블·문구·오류·키보드·레이아웃을 검수했다. 보완 PNG를 Orca PNG로 표시하지 않았다.
- [x] 이메일·코드·cookie·token을 가린 보고와 증거를 작성했다. 실제 화면 관측과 자동 DOM/overflow/continuity 근거를 구분했다.
- [x] 자기 실행 자원·메일·DB anonymous volume·전용 page를 정리했다. 세 loopback port 종료를 확인했다. 제품 코드·제품 수치·기술 승인·PLANS/board·다른 인박스·GitHub·CF·배포를 변경하지 않았다.

변경 이유와 판단: 로컬 신원 기능 단계의 직접 UI 근거를 보완했다. Orca blank 캡처는 제품 결함으로 단정하지 않았다. 독점 표시 복구도 불안정해 이 Dispatch의 coordinator ask 답변에 따라 실제 별도 Chromium을 사용했다. 제품 렌더링을 합성하거나 수락 기준을 낮추지 않았다. 초기 DB bootstrap 감지와 SMTP 동일 port 재기동 실패는 자기 환경에서 복구했다. 최종 브라우저 명령은 exit 0이며 전체 과정의 미실행과 실패를 실행 기록에 남겼다.

Finding: 새로운 UI critical/high는 없다. low F-UI-01은 기존 회원 gate가 만료 시각을 Go 기본 UTC·중복 `+0000`으로 표시하는 가독성 문제다. 만료·원문 부재·승인 비활성·홈 복귀는 유지돼 이번 신원 UX를 차단하지 않는다. 후속 DEV gate 표시 과제에서 검토한다. 별도 보안 리뷰의 F1/F4/F2는 이 UI 결과로 해소하지 않는다.

검증과 미실행: 독립 후보 Go 빌드 exit 0, 최종 실제 Chromium 상태별 흐름 exit 0, 직접 PNG 검수 PASS다. challenge 기한·session Verified·자기 fixture budget은 상태 화면 검사에 맞춰 조정했다. 10분/12시간 실시간 대기와 전체 보안·한도 QA로 보고하지 않는다. 실제 사람 이메일·운영 SMTP·공개 배포·QA-P06 공개 확인·QA-P07 사람 확인·최종 동일 이메일 노우↔다닷은 미실행이다. 제품 코드 변경이 없으므로 product-lint 적용 대상은 없다. 보고·링크·증거 형식과 Git 공백 검사 결과는 자체 검증 근거에 남긴다.

수락 경계: 이 결과는 로컬 신원 코드 단계의 UX 근거다. coor가 독립 리뷰·QA와 연결해 코드 통합을 판단하고, OPS는 수락 후보로 공개 운영·SMTP·외부 확인을 수행한다. 운영 최종 수락이 완료돼야만 로컬 코드 단계 검수가 가능하다는 순환 조건으로 읽지 않는다. 전체 일반 서비스 PASS와 실제 이메일/노우↔다닷 수락은 그대로 후속이다. 새 보안 수정 후보의 UI delta는 별도 과제다.

산출물: `docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-UI.md`, 같은 위치의 `SAR-PUBLIC-IDENTITY-001-UI-evidence/`, `docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-UI.md`. 이 인박스와 결과 전문은 work.py finish로 보존한다. 소유 파일만 커밋하며 결과 SHA와 fixed59·finding·미실행은 이 Run worker_done으로 직접 회신한다. 기본 브랜치 병합·push·역할 동기화는 coordinator 담당이다.

## coor 탐색 근거

SAR-PUBLIC-IDENTITY-001-UI-{find,documents-find,context}.json은 코드/문서를 분리한20개 후보와 필수 규칙의 근거다. keep은 README·member.go·mail_sink.py·identity.go·http.go·실행준비 compose/main/verify_mvp, DEV 실행 기록·이번 UX/제품 정본·필수규칙이다. 이전UI 보고·로그와 담당 context는 과거 기준/실패의 참고다. 현재후보59를 직접 확인한다. context의 추천을 참고하되 필수정본·passage 미송신자료를 제외하지 않는다. 기존trial 조작 run_trial.py는 필요시 확인만 하며 종료된trial을실행하지 않는다.
