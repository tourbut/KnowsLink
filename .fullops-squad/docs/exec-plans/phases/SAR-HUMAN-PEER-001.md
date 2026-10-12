---
title: SAR-HUMAN-PEER-001 — 제품 결정과 UI/UX 전면 개편 인계
status: draft
updated: 2026-10-12
owner: designer
tasks: [SAR-HUMAN-PEER-001]
summary: 사람 친구와 대화의 제품 결정·전면 디자인·탐색과 DEV QA 인계를 기록한다
---

# SAR-HUMAN-PEER-001 — 제품 결정과 UI/UX 전면 개편 인계

## 범위와 현재 상태

기획·디자인 기준은 `f20b33a11addeb3debcb23e2e835be31c7dd1454`다. 준비 커밋은 `dd50978`이며, 재개 시 HEAD는 `e521a0d`였다. 역할은 designer이고 실행 선택은 Codex / gpt-6.1-sol / high다. 테스트 레벨은 lite다. 하위 위임 standard는 선택 예산이며 이번 문서 완성에는 하위 worker를 사용하지 않았다.

최초 designer dispatch는 모델 용량 오류로 중단됐다. coor의 `msg_59d0bbe24839`가 같은 Task 재시도와 기존 초안의 소유권을 확인했다. 현재 dispatch는 `ctx_76b78b46e9c3`, Task는 `task_491ede5000f9`, Run은 `run_5534969da2d0`다. 기존 초안을 보존하고 충돌·필수 기록·역할별 패킷을 보완한다. 제품 코드는 수정하지 않는다.

현재 결과는 구현 지시가 가능한 제품·디자인 설계다. 실제 기능·운영·직접 구현 시각 수락은 아직 아니다. coor가 고정 결과의 독립 문서 리뷰와 main 통합을 마친 뒤 DEV를 배정한다.

## 사용자 근거와 책임

사용자는 자기 실제 이메일 아래 에이전트 등록과 이름 지정, 사람끼리 친구 수락, 친구의 하위 에이전트 목록과 대화 송수신을 요청했다. 후속 원문은 “UI/UX도 한번 다시한번 디자인스킬써서 전면개편해보고”다. 전체 원문은 [designer 역할 로그](../../../handovers/logs/2026-10-12_to_designer.md)에 완료 때 보존한다.

이 요청은 Contacts=agents·agent별 첫 수락·composer 금지를 대체하는 명시적 제품 근거다. 기술 구조·API·DB·이행 절차는 DEV가 결정한다. 승인되지 않은 운영 배포나 실계정 외부 발송은 이번 범위에 넣지 않는다.

적용 규칙은 `fullops-common-0.3.3`의 README·coding-style·testing·security, FULLOPS.md, orca-agents.md, project.md, docs/agents/document-writing.md다. 주 디자인 스킬은 frontend-design이다. fullops-frontend의 toolchain.md·delivery.md와 fullops-work·fullops-deliverables를 적용했다. 보안·필수 검사·독립 리뷰의 예외는 없다. 삭제·이동 정리도 실행하지 않는다.

## 확정 결정과 추적

제품 정본은 [SAR-HUMAN-PEER D02](../../planning/product-specs/SAR-HUMAN-PEER.md)다. 화면 정본은 [service-member-ui D04의 현행 절](../../design-docs/mockups/service-member-ui.md)이다.

| 요구 묶음 | 결정 | DEV/QA 관찰 |
|---|---|---|
| HP-01 이름·신원 | 검증된 기존 issuer/sub·owner·agentID·키를 유지한다. 사람/agent 이름은 1–40자다. | 재로그인·이름 변경 뒤 동일 소유권과 키, 중복 이름 구분 |
| HP-02–04 친구 | 공유 코드로 요청하고 수신자가 명시 수락한다. 코드 조회만으로 존재·이메일·agent를 공개하지 않는다. | 없는/차단 코드 동일 안내, 반대 요청 자동 수락 금지, 중복·만료 경계 |
| HP-05–07 범위 | 수락하면 현재·향후 활성 agent가 보인다. agent별 재등록은 없다. 기존 자기 범위 확대에는 한 번 동의한다. | 새 agent 보임, 제삼자/철회 거부, 자기 두 agent 흐름 |
| HP-08–11 대화 | 친구 아래 두 agent를 고정 선택한다. owner 작성과 agent 작성은 구분한다. 매 발송은 독립 메시지다. | 세 번 이상 왕복, 앞선 TTL 뒤 새 발송, 응답 유실·중복·충돌, 상태 의미 |
| HP-12 보존 | 새 대화 원문만 수락 후 최대 24h다. 전달 TTL은 최대 180s다. 업무 원문은 기존 보존을 유지한다. | 제어 시각 경계, 서버 미리보기/원문 제거, 소급 복원 없음 |
| HP-13–15 철회 | 해제·차단은 모든 신규 전달·실행·gate·result·legacy 경로보다 우선한다. | 열린 composer, queued/leased, 옛 세대, SSE/remote/adapter 우회 거부 |
| HP-13 과거 열람 | 현재 회원 세션의 본인 참여 사본만 잔여 보존 기간 읽는다. | 상대 최신 프로필/새 agent 비조회, 철회 키/만료 세션/C 거부 |
| HP-16 수락 | 이름→친구 수락→agent 노출→양방향 대화→해제 차단이 한 고정 후보에서 동작한다. | DEV 필수 검사, 독립 QA, 직접 구현 시각 검수, 별도 fixed-SHA 리뷰 |

친구 수락은 결제·삭제·권한 부여·다른 채널 발송·약속의 승인이 아니다. C1–C5·검증 typed gate·현재 인가·CSRF·CSP·rate·deny-by-default를 유지한다. 웹 작성은 개인키 이전이나 agent 서명 위조 없이 실제 작성 주체를 표시해야 한다.

친구 한도는 전체 active 400·회원당 active 20, 전체 pending 200·회원별 송신/수신 pending 각각 10·24h다. agent 5/회원·키 3/agent·기존 큐/게이트/rate/receipt 한도는 유지한다. 대화도 24h 20000 메시지 상한을 공유한다. 가격·Free N·상품 슬롯은 결정하지 않는다. 상세 경계와 한도는 D02를 재사용한다.

## 기존 데이터의 권한 이행

기존 active pair는 정확히 그 두 agent의 범위만 유지한다. 이를 사람 친구로 자동 승격하거나 새 agent 목록을 공개하지 않는다. cross-owner 전체 연결에는 양쪽 검증 회원의 새 요청/수락이 필요하다. legacy pending은 종료하고 새 친구 요청으로 안내한다. denied/revoked는 이력이며 재로그인·이름 변경·코드 교체로 복구하지 않는다.

기존 member·신원·owner·agent·key·credential·설치 권한은 보존한다. 무명 agent는 '이름 없는 에이전트'와 구분값으로 표시한다. 소유가 불명확하거나 합성인 owner는 새 친구 모델에 귀속하지 않는다. 이미 삭제한 메시지 원문을 복원하거나 업무 M/H/R을 새 대화에 복사하지 않는다.

새 해제/차단은 legacy pair보다 우선한다. rollback/복원이 새 권한·철회 상태를 증명하지 못하면 전달을 중지한다. 예전 JSON writer가 필드를 버리고 권한을 되살리는 복귀는 허용하지 않는다. 구체적인 배포 순서·백업/복원·버전 경계는 DEV가 기술 문서에 쓰고 OPS가 수락 뒤 실행한다. 이 문서 과제는 DB를 변경하지 않는다.

## 디자인과 실제 관찰

기존 polish after의 `png/home-1440-s1of7.png`와 `png/home-390-s1of10.png`를 직접 열었다. 출처는 [manifest](../../evaluations/qa-reports/SAR-FRONTEND-POLISH-001-DEV-test/after/manifest.json)다. 합성 template 렌더의 화면이며 현재 운영 계정·새 구조 동작 증거가 아니다. 첫 화면은 ID·키·명령과 관리 카드가 차지한다. 이 관찰이 대화·친구·내 에이전트로 정보 구조를 바꾼 이유다.

D04에는 두 차례 계획·자기 비평·수정, 여섯 색 토큰·글꼴·간격·형태, 1440 desktop 및 390/320 mobile 와이어프레임, 빈 상태·오류·승인·해제·보존 안내와 접근성을 기록했다. 특징은 사람 이름 아래 agent 이름이 이어지는 목록이다. 새 대화에서만 발신/수신을 선택하고 이후 헤더와 작성란에 고정한다. 기술 정보는 연결/계정 상세로 옮긴다. 장식 카드 묶음·그라디언트·알약 버튼·상시 보안 설명·가짜 대화는 제외한다.

`frontend.py --repo . --host codex --model gpt-6.1-sol --effort high --check`는 exit 0, `configuration_ready`였다. frontend.json의 활성 tuple/정본을 재사용했다. 설치·설정 준비와 로그인·MCP·새 브라우저 QA 성공은 구분한다. 인증·MCP·새 제품 브라우저 QA는 not_run이다. 기존 Go template/CSS를 사용하며 새 프레임워크·폰트·이미지·유료 도구를 설치하지 않는다.

## 탐색 원문·partial 처리

[designer packet](../../evaluations/jev/SAR-HUMAN-PEER-001-packet.json)의 원문을 읽었다. instruction excerpt의 partial은 역할 인박스 전체 읽기로 보완했다. 필수 path/category와 remaining_context_paths를 일반 검색·원문/관련 절 확인으로 판정한다. 낮은 점수·unknown은 파일이나 기능의 부재 근거가 아니다.

`member_agents.go`의 ownAgents/memberAction, `member_pages.go`의 memberStyle/memberPages, `identity.go`의 signIn/session, `test_messages.go`의 시험 allowlist/handler, `public_text.go`의 signed text ingest를 확인했다. 사람 이름·친구·웹 composer는 현재 모델에 없다. `public_text.go`의 관련 답장은 부모 기한 내 한 번이며 새 연속 대화 규칙과 호환 검토가 필요하다. `test_messages.go`는 일반 회원 구현을 대신하지 않는다. 존재하지 않는 `member_messages.go`의 과거 오지정은 다시 쓰지 않는다.

패킷의 stale find/documents-find/context producer는 권위 있는 완전 탐색으로 취급하지 않는다. 원문 지시·현재 코드·수동 문서 검색으로 보완하고 새 DEV/tester의 find/context/packet을 역할별 키로 작성한다. source SHA와 instruction/attempt는 dispatch 직전 coor가 실제 준비 SHA로 갱신해야 한다. 이 설계 후보를 제품 구현 SHA로 표시하지 않는다.

PLANS·ops-guide·옛 역할 로그는 크기 때문에 의미 탐색에서 제외됐을 수 있다. PLANS의 현재 과제와 작업 순서, D12/D13의 현행 운영·복귀 경계, 해당 로그의 과제 제목과 보존 사실을 직접 확인한다. 과거 검증 원문은 수정하지 않는다. 동적 호출 그래프·운영 데이터 전체 분석은 수행하지 않았으며 DEV의 공유 인가 호출자/legacy fixture 검증으로 인계한다. 모든 비선택 path/category 결과는 [outcomes](../../evaluations/jev/SAR-HUMAN-PEER-001-packet-outcomes.json)에 기록한다.

## 산출물 영향과 다음 담당

| ID/문서 | 이번 결과 또는 영향 | 다음 담당·시점 |
|---|---|---|
| D01 | 사업 가치·사람 중심 범위 갱신 | 문서 독립 리뷰 |
| D02 | HP-01–16 새 정본, 기존 PS/MVP/setup 및 원천 우선순위 안내 | DEV 구현·tester 추적 |
| D04 | 현행 UI 전면 개편, 과거 설계/시각 이력 보존 | DEV 구현 뒤 designer 실제 캡처 판정 |
| D03·D05 | 구조·인가·웹 작성/친구/대화 API·직접 링크 영향 | DEV가 같은 구현 과제에서 갱신 |
| D06·D07 | Member/Agent 이름·친구 세대·과거 참여/보존·legacy 이행 영향 | DEV가 실제 모델/DB 설계 갱신 |
| D08 | 실제 table/SQL 변경 여부는 기술 계획 후 결정 | DEV가 실제 schema 생성 또는 변경 불필요 근거 |
| D09·D10 | 상태 전이·안전 정리·공유 인가·화면/adapter 모듈 영향 | DEV가 실제 CRUD/모듈 갱신 |
| D11 | 새 사용자 제품 흐름 추가, 과거 베타 절차 구분 | DEV가 실제 경로·설치 절차 갱신 |
| D12·D13 | 백업·복원·버전/rollback 및 수락 뒤 배포 영향 | DEV 기술 인계 후 OPS 실행 |

원 지시의 'D07 사용자 안내'는 실제 인덱스와 달랐다. D07은 데이터베이스설계서이며 사용자설명서는 D11이다. D07을 사용자 안내로 잘못 stamp하지 않는다. PLANS/board의 현재 상태는 coor가 갱신한다.

현재 후속 지시서는 [DEV](../../../handovers/to_dev.md)와 [tester](../../../handovers/to_tester.md)다. DEV는 설계 통합 전 대기이고 tester는 고정 DEV SHA 전 대기다. 새 역할 dispatch는 이 Task/dispatch ID를 재사용하지 않는다. DEV는 짧은 기술 계획부터 전체 구현·필수 회귀·기술 문서를 같은 과제에서 끝낸다. tester는 독립 A/B/C·legacy fixture와 실제 loopback UI/API/DB를 확인한다. 실제 이메일·실계정·외부 상대 발송·운영 이행은 별도 실행으로 남긴다.

## 검증 결과와 미실행

문서 검증은 front matter stamp, 변경 패치 공백, 로컬 링크/산출물 strict, packet identity/outcomes, 토큰 대비 계산과 사용자 요구·와이어프레임 대조다. 명령 결과와 최종 고정 SHA는 완료 보고에 기록한다. 코드 변경이 없으므로 제품 동작 시험을 새로 만들거나 과거 통과로 대체하지 않는다.

새 실제 기능, 브라우저 조작/키보드/200%/스크린 리더, 실 Google 계정·노우/다닷 왕복, 운영 이행·복원은 not_run이다. 이는 DEV·tester·designer·OPS 후속으로 남는다. 새 다크 모드·공개 검색·첨부·장기 아카이브·읽음·agent별 공유 토글·상품 기능은 생략했다.

### 이번 실행 근거

- `deliverables.py --repo . --strict`: exit 0, 13종 검사·문제 0·경고 0. D01/D02/D04/D11과 관련 원천·phase/context의 메타데이터는 stamp로 기록했다. 처음 메타데이터가 없던 원천은 owner/summary를 명시해 보완했다.
- `git diff --check`: 초안 와이어프레임의 trailing whitespace 한 줄을 발견해 수정했다. 수정 후 exit 0이다.
- 토큰 대비: [design-check.json](../../evaluations/jev/SAR-HUMAN-PEER-001-design-check.json)의 7개 조합이 4.5:1 본문·3:1 입력/포커스 기준을 통과했다. 실제 화면 픽셀 검사는 아니다.
- 로컬 링크 존재: 변경 Markdown과 후속 인박스에서 169개를 확인했고 누락 0건이다. fragment·원격 링크의 의미 검증은 포함하지 않는다.
- 역사 보존: 기준 f20b33a의 product.md·decisions.md 원본문이 새 안내/메타데이터 아래 그대로 포함됨을 assert했다. D04의 이전 1절 이후 본문도 그대로다. 새 안내를 포함한 전체 파일의 바이트 동일성은 주장하지 않는다.
- designer completion packet 검사: exit 0. 52개 비선택 path/category 결과와 미확인 15건의 수동 검토를 outcomes에 기록했다. 원 패킷의 producer mismatch·큰 문서·remaining 11개를 미검토로 숨기지 않았다.

### 후속 탐색 패킷과 대기 경계

DEV와 tester의 code/documents find 및 context는 역할별 키로 실행해 원본 JSON을 보존했다. DEV route는 implementation/dev와 Codex gpt-6.1-sol high로 고정했다. tester route의 Grok 4.7 high는 추천이며 실제 새 dispatch의 실행 tuple은 coor가 확정한다. 각 route는 현재 역할 인박스에 bind했다. 처음 packet-decisions의 action 값이 스키마에 맞지 않아 exit 1이었다. 허용 값 `change`로 고쳐 재생성했고 두 packet 명령은 exit 0이었다.

준비 패킷 HEAD는 `5d5d26cd02e98a6a0978ca681a7565a6b6e5a29d`다. [DEV packet](../../evaluations/jev/SAR-HUMAN-PEER-001-DEV-packet.json)은 후보 65·partial true·unknown 17·remaining 32다. [tester packet](../../evaluations/jev/SAR-HUMAN-PEER-001-TESTER-packet.json)은 후보 203·partial true·unknown 17·remaining 166이다. tester direct_edit/document_update는 0이며 제품 코드 수정 권한은 없다. 큰 과거 문서·로그와 동적 호출은 추천 목록이지 자동 실행 명령이 아니다. HP/D04와 명시한 핵심 코드/권한 검증이 우선이다.

동시 작성한 후속 인박스는 다른 packet의 참조 원천에도 포함돼 source hash가 바뀌었다. 준비 packet의 `check`는 exit 1, `packet source changed; regenerate packet before worker start`였다. 이 결과를 PASS로 표시하지 않는다. 두 인박스는 이미 waiting-design-integration / waiting-dev-sha이며 실행을 시작하지 않았다. 최종 문서 리뷰·main 통합과 역할 인박스 확정 뒤 coor가 DEV의 실제 준비 HEAD에서 find/context/route/packet을 갱신한다. tester는 고정 DEV 구현 SHA에서 갱신한다. 그때 비선택 path/category·partial 잔여의 수동 검토를 각 worker가 자기 outcomes에 기록해야 한다. 이번 designer의 outcomes로 미래 구현/QA 완료를 대신하지 않는다.

원 제품·기술 경계를 바꾸지 않는 탐색/아카이브 과정의 실패는 원인·교정·대기 조건을 위에 보존했다. 최종 SHA와 최종 HEAD의 문서 lint 결과는 worker_done 및 외부 검증 파일로 고정한다. 제품 코드 변경이 없으므로 새 제품 테스트는 해당 없음이며 후속 DEV 필수 product-lint/product-test는 생략하지 않는다.

### FullOps lint와 규모 판단

`lint.py --repo <designer> --from f20b33a11addeb3debcb23e2e835be31c7dd1454`를 커밋 `d6a4ef55df44`에서 실행했다. 종료코드는 1, ERROR 0·WARNING 3·실행 불가 2다. 등록 `product-lint: make lint`와 `product-test: make test`는 Windows 세션에서 make를 찾지 못해 unavailable이다. 이 결과를 필수 제품 검사 PASS로 기록하지 않는다. 제품 코드·실행 설정을 바꾸지 않은 기획 문서 과제이므로 해당 제품 검사는 후속 DEV의 같은 고정 구현 HEAD에서 실제 실행한다. make shim이나 다른 체크아웃을 임의로 실행해 현재 문서 후보의 통과로 표시하지 않는다.

원 결과는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-designer/pre-finish-lint.json`에 보존했다. finish 뒤 최종 HEAD의 재검사 파일은 같은 폴더의 `final-lint.json`이다. 최종 worker_done은 실제 결과와 이 경로를 인계한다.

경고는 PLANS 1666줄(기준 1657), D04 579줄(기준 411), 적용 문서 추가 467줄/400줄 상한이다. PLANS 증가는 준비/기동 기록이며 coor 소유다. 이번 designer는 수정하지 않았다. D04는 현행 전면 명세와 과거 1–14절 판정의 원문을 보존하는 정본이다. 삭제·압축·분할로 과거 검수 문장을 바꾸거나 정본을 두 곳에 복제하지 않았다. 여러 제품 문서·역할 지시서를 예상한 원 지시 범위에서 HP/이행/UI/QA를 함께 전달하므로 이번 기능을 조각내지 않았다. 후속 문서 정리는 과거 기록 보존과 링크 검증을 함께 맡는 별도 소유 과제로 처리한다.
