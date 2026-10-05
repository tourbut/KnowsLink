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
