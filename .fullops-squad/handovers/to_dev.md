---
title: SAR-SYSTEM-DIAGRAM-001-DEV — 현재 시스템 구성도와 상세 기술 설명을 docs에 추가
status: draft
updated: 2026-10-14
owner: dev
tasks: [SAR-SYSTEM-DIAGRAM-001-DEV]
summary: 현재 시스템 구성도와 상세 기술 설명을 docs에 추가
attempt: 6c817804bcfa4eeaad94261ab79a8b28
base: 8f255f9991787e14b38052029d334b9536160fdb
subagent_level: standard
test_level: lite
---

# SAR-SYSTEM-DIAGRAM-001-DEV — 현재 시스템 구성도와 상세 기술 설명

- 2026-10-14 coor→dev, Purpose documentation. 기준8f255f9991787e14b38052029d334b9536160fdb, fullops-common-0.3.3, lite/standard. Codex/gpt-6.1-sol/high 고정.
- 사용자 승인: 현재 시스템 구성 다이어그램·상세 기술 설명을 문서에 넣는다. docs 전용 작업이며 제품 코드/설정/의존성/DB/배포/계정/메시지/삭제/이동은 변경하지 않는다.
- 담당 repo0b08ec4c-9e3a-4613-8197-5a835e545335, C:/Users/shin/orca/workspaces/KnowsLink/dev, fullops/dev. 병합coor, origin/main.
- 복귀 Run run_5534969da2d0, 현재 coor term_38ae16dc-0ca4-4beb-ab00-d1e1b641f061. 실제 dispatch preamble 주소로 완료한다. 기존 dev 사용자 터미널은 보존하며 fresh 문서 worker로 착수한다.

## 결과와 범위

현재 구현과 최신 확인된 배포를 근거로 한국어 상세 문서를 docs/system-design.md에 작성한다. docs/architecture.md는 입문 설명을 유지하고 상세 문서로 연결한다. docs/README.md와 루트 README에도 바로 찾을 수 있는 링크를 넣는다. 설치 본문을 재복제하지 않는다. 문서 기준SHA8f255f9와 마지막 운영 관측2026-10-13/제품2747fe98653d5c6946f41c219fc256395333d7fe를 구분한다. 운영 실시간 재확인이나 새 실계정 시험으로 주장하지 않는다.

구성요소와 물리/논리 경계를 상세히 설명한다: 웹 회원·owner·여러agent·사람 친구 수락과 legacy pair, Grok 로컬 Command MCP/watch/private Inbox/선택 webhook routine, Dots Remote HTTPS MCP/OAuth/설치별 서버 키와 Inbox/Events callback, Cloudflare ingress·cloudflared·Go relay·Postgres·별도 migration, Google 로그인 신원, 선택 SMTP 경로, 개발/운영의 차이. 현재 구현에 없는 queue broker·Redis·microservice·자동답장·업무 실행·무조건 알림 보장은 만들지 않는다.

실제 요청 경로·인증/권한/Ed25519 및 owner wire·송신수락/수신보관/ACK/화면표시·멱등/기한/재시도·관계 해제/철회/generation·singleton JSONB/row lock/epoch CAS/LISTEN NOTIFY·180초 전달/24h 웹사본/로컬Inbox/백업WAL·schema/writer guard·운영격리/백업/제약을 설명한다. 구성요소/소스 경로/프로토콜·책임·저장위치 표와 실제 코드 상대링크를 제공한다. 구현되지 않은 확장 제안은 범위에 넣지 않는다. D03/D12 원천을 복제하거나 대규모 과거기록을 다시 쓰지 않고 현행 상세문서에 연결한다.

## 다이어그램 형식

기존 docs/architecture.md와 같은 Mermaid를 docs/system-design.md 안에 직접 넣는다. GitHub Markdown에서 바로 표시되고 원문을 수정할 수 있어야 한다. 전체구성1개와 Grok/Dots의 서로 다른 전달·수신 흐름1~2개(총2~3개)를 만든다. 구조/배치·신뢰경계에는 flowchart, 시점순 송수신에는 sequenceDiagram을 사용한다. 관계 모델은 본문/표로 설명하며 필요 없는 추가그림은 만들지 않는다. 한그림은 architecture 최대9노드/12화살표, sequence 최대5lifeline/12message 정도로 읽기쉽게 나누며 의미상필수관계를 생략하지 않는다. 공개/비공개영역·선택형/비동기·방향·신호와본문/ACK를 구분한다. Mermaid를 데이터구조나서비스실행으로오인하지 않게 실제와도식의단순화를설명한다.

coor가 별도 HTML/SVG에 쓸 스타일을 물었지만 현재 문서완료에필요한입력은아니다. Mermaid기술구성도·상세본문은지금끝까지작성한다. 별도브랜드HTML/SVG와전역스타일설정/프로필저장은이과제범위에넣지않는다. 답변대기를이유로문서작업을멈추거나사용자에게같은질문을보내지않는다. diagram-design의구조/밀도원칙만참고하고HTML/SVG출력·최초스타일온보딩스킬은이네이티브Markdown출력에적용하지않는다.

## 먼저 읽을 근거

AGENTS.md, .fullops-squad/FULLOPS.md, rules/common/README.md와 coding-style/testing/security, project.md, docs/agents/document-writing.md, docs/deliverables/README.md, contexts/dev.md. 사용자문서: README.md, docs/README.md, docs/architecture.md, docs/operations.md. 기술근거는 다음18개에서 시작하고 필요한 관련 파일만 보완한다:
cmd/relay/main.go; cmd/migrate/main.go; compose.yaml; deploy/knowslink/compose.ops.yaml; deploy/knowslink/tunnel/public-ingress.yml; internal/relay/remote_mcp.go; internal/relay/remote_auth.go; internal/relay/remote_events.go; internal/relay/text_events.go; internal/relay/member_receipt.go; internal/relay/member_agents.go; internal/relay/google.go; adapters/src/mcp.ts; adapters/src/inbox.ts; db/migrations/00002_human_writer.sql; .fullops-squad/docs/design-docs/module-design.md; .fullops-squad/docs/design-docs/database-design.md; .fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS/report.md.
필요시 interface-design.md·db queries·adapter index/login/wake·registry·실계정 NOTIFY 기록·로컬 디자인 정본을 관련 절만 읽는다. 과거전문 전수재독은 요구하지 않는다. 비밀.env·private 계정 파일은 읽지 않는다. 공개예제/코드/이미 존재하는 검증보고만 근거로 사용한다. 새 외부API계약 설명을 추가해야 하면 Context7 resolve→query(월quota실패 시 공식문서)를 적용한다. OpenAI제품 현재계약은 openai-docs를 읽는다.

## 검증과 완료

문서/소스 링크·fragment/닫힌fence·Mermaid 문법과 노드/방향/참여자·소스 일치를 검토한다. 기존 실행 가능한 Mermaid 검증/렌더 도구가 있으면 재사용하고 새프로젝트의존성은추가하지않는다. 구조추출검사는실제렌더통과가아님을구분한다. 필요한브라우저확인은문서diagram범위에서하고미실행은솔직히기록한다. 제품 전체테스트를 재실행하거나 새상설도구/의존성을 추가하지 않는다. 기존 exactcd4 Linux gate는 제품 입력 동일성을 확인해 재사용하고 문서HEAD 신규제품PASS라고 쓰지 않는다. git diff --check, strict 및 문서 정적lint를 남기며 Windows make 부재는 unavailable/reason으로 보존한다.

예상 사용자문서 추가는 상세설명 약200~350줄과 Mermaid diagram2~3개다. 변경규모경고 때 필요한범위와분할판단을보고한다. 파일 소유권은 위 docs/Markdown, root README 연결, 본과제phase/QA/outcomes/context/완료 인박스뿐이다. D03/D12는read원천이며 내용충돌이 있으면 명확한현재절 링크만 제안하고 임의변경하지 않는다. 내부메일·실agent/app/subscription/dispatch/session·운영private주소/경로·비밀은 사용자문서에 추가하지 않는다.

phase/qa-reports/SAR-SYSTEM-DIAGRAM-001-DEV/report.md에 설계명세·코드근거·Mermaid형식판단·검증/한계·실제작성자세션을 기록한다. packet192 등 이전과제결과를복사하지않고 이번packet의모든쌍/unknown/partial를이번실제확인범위로 outcomes에남긴다. 본문범위변경없는진행/복귀/완료갱신만인박스에반영한다. work.py finish로원지시/완료전문을logs보존·인박스비움·commit후 [완료] SAR-SYSTEM-DIAGRAM-001-DEV | 브랜치 fullops/dev | SHA <40자리> 형식으로 worker_done한다. coor가별도세션fixedsnapshot 내용/시각리뷰·main일반통합/push를수행한다.
