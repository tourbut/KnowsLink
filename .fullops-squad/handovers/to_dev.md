---
title: SAR-PUBLIC-AGENTS-001-DEV — 일반 회원의 agent 연결·키 수명·관계 관리를 구현하고 검증한다
status: draft
updated: 2026-10-06
owner: dev
tasks: [SAR-PUBLIC-AGENTS-001-DEV]
summary: 일반 회원의 agent 연결·키 수명·관계 관리를 구현하고 검증한다
---

# SAR-PUBLIC-AGENTS-001-DEV — 일반 회원의 agent 연결·키 수명·관계 관리를 구현하고 검증한다

- 작성일: 2026-10-06
- From / To: coor / dev
- 상태: ready
- 담당: repo 818c78e5-d51c-4ff4-aa88-70e9ee185fbb, /home/shin/orca/workspaces/KnowsLink/fullops-dev, fullops/dev
- 병합 책임자: coor, main/origin
- 복귀: /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_6895aaf1-7b43-4fe0-a416-76f1255a5946, run_8ca8bc058ab7. Task/Dispatch는 착수 영수증과 preamble을 사용한다.

## 현재 상황과 적용 기준

사용자가 다음 작업 진행을 요청했다. 신원 로컬 코드·독립 리뷰/QA는 main/origin d2f7ba5aeb6644fd2b27fe6ada5b61db5976933b에 수락됐다. 실제 이메일·공개·최종 노우↔다닷은 미검증이다. 계정 화면/접근 대기는 독립 구현을 막지 않는 제품 정본에 따라 다음 기능을 진행한다. DEV 인박스와 현재 터미널은 비어 있고 작업 트리는 clean이다.

기준 SHA는 d2f7ba5aeb6644fd2b27fe6ada5b61db5976933b다. fullops-common-0.3.3의 README와 코딩·테스트·보안 규칙, project.md, FULLOPS.md를 적용한다. 제품 정본 SAR-PUBLIC-SERVICE PS-04–07와 해당 PS-11/개인정보 한도, UX-04–05가 범위다. route 요청의 UX-06 언급은 후속 MESSAGES의 receipt 범위이므로 이번 구현에 포함하지 않는다. 기술 분석·계획·구현·검증은 같은 DEV 과제에서 맡는다. 제품 규칙 변경만 coor를 통해 designer에게 묻는다.

## 먼저 읽을 문서

Jev find/code와 documents를 별도로 실행했다. 신규 일반 연결 기능은 absent 추천이며 기존 registry/pair/신원 기반을 확인해 재사용한다. context 결과는 docs/evaluations/jev/SAR-PUBLIC-AGENTS-001-DEV-context.json이다. 충돌/지시문 경고는 없다.

- .fullops-squad/FULLOPS.md
- .fullops-squad/project.md
- .fullops-squad/rules/common/README.md 및 연결된 세 규칙
- .fullops-squad/contexts/dev.md
- .fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md
- .fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md
- .fullops-squad/docs/design-docs/interface-design.md
- .fullops-squad/docs/design-docs/data-model.md
- .fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV.md
- README.md
- internal/relay/identity.go
- internal/relay/member.go
- internal/relay/store.go
- internal/relay/http.go

필요 시 adapters/README.md와 위 파일의 실제 호출 경로를 좁혀 확인한다. Jev omit 추천이나 지원 클라이언트 안내에 필요하면 사용 근거를 보고한다. 라이브러리 API 근거는 기존 동일 버전 근거를 재사용하고 부족하면 Context7→공식 문서로 확인한다.

## 해야 할 일과 파일 소유권

- [ ] 기존 구현을 확인하고 짧은 기술 계획·예상 변경 규모를 실행 기록에 남긴다.
- [ ] 일반 회원이 자기 agent를 생성하고 지원 클라이언트를 연결하게 한다. 관리자 자격·SSH·서버 수동 파일 배치를 요구하지 않는다. 연결 수단은 owner/agent/client에 결속하며 최대10분·1회 사용, 실패/취소/만료는 미연결이다.
- [ ] 5분 이내 재인증과 새 공개키 PoP, 지문·활성/철회 표시, 등록/회전/선택 철회, agent별 credential 및 권한 분리를 구현한다. 개인키를 수집하지 않는다. 철회한 옛 키는 enqueue/pull/ACK/exec/result에서 거부한다.
- [ ] B-agent 식별자 초대, B-owner 명시 수락/거절, 양측 철회, 중복/동시 수락·새 세대와 같은 owner 두 agent의 명시적 첫 수락을 구현한다. 이메일/회원 디렉터리를 노출하지 않는다.
- [ ] 활성 agent 전체200/owner5, 키3, active pair400/owner20, pending200/송신owner10/24h와 공통 rate를 집행한다. 신규 포화에도 deny/revoke/unpair를 독립 정리 budget으로 처리한다. 재시작·동시성과 기존 신원/비활성화 권한 경로의 영향을 검증한다.
- [ ] UX-04–05와 변경된 자기 홈 안내, 클라이언트 실행 위치·권한·기한·만료/미지원/철회·안전한 오류를 구현한다.
- [ ] 관련 회귀·원천 문서·실행 기록·QA/UI 후속 인계를 완성한다.

소유권: cmd/internal/adapters/db/scripts와 필요 Make/Compose 설정, 기술 정본·README·자기 contexts/인박스/완료로그. 공유 PLANS는 자신의 결과를 추가하고 기존 기록을 보존한다. board는 coor 소유다. 제품 기획 정본은 변경하지 않고 실제 요구와 충돌하면 질문한다. 변경 규모가 크면 SIZE-002의 실제 규모와 한 과제에 남긴 이유를 보고한다.

## 완료 기준과 검증

DEV는 현재 변경의 자동 검사와 실패·경계 회귀를 끝낸다. 자기/타 회원 ID 바꿔치기·agent 자격 owner권한 우회·재인증/PoP 실패·만료/재사용/취소·철회 전후·동시 첫 연결/수락·한도 경계/초과/재시작·공통budget/NAT·정리 가능성을 검증한다. make lint, make test, make verify-mvp와 변경 영향에 필요한 검사를 명령 자체 exit code로 남긴다. corrected relay-stop/Go/relay-up 순서와 기존 invalid_lease 실패 근거를 유지한다. git diff --check, strict 산출물 및 깨끗한 커밋 HEAD의 FullOps lint --from 기준SHA를 완료한다. ERROR/실행불가와 설명되지 않는 실패는 원인 해결 또는 구체적 보류를 남긴다.

DEV 완료는 최종 수락이 아니다. 완료 고정 SHA 뒤 별도 OPS 세션의 독립 보안 delta 리뷰, tester의 교차 계정·연결/회전/철회·관계/한도 QA, designer의 UX-04–05 직접 시각 검수를 coor가 배정한다. 기존 신원 증거는 의존성 동일성 확인 시 원래SHA로만 재사용한다. 실제 사람 이메일/공개/실클라이언트는 fixture PASS와 분리한다. 안정 후보에서 독립 QA와 UI 검수를 한 번 수행한다. 검수용 fixture 실행 안내와 필요한 마스킹 캡처 항목을 인계한다. 영상은 정지 화면으로 판단 불가능한 시간 변화에만 필요하다.

## UI 기준

Go member template·기존 owner UI의 정보 구조/스타일을 재사용한다. UX 정본은 위 mockups 문서다. theme/Tailwind/shadcn/디자인 전용 lint는 현재 없으므로 해당 없음 근거와 적용되는 make lint/FullOps DESIGN 경고 처리를 보고한다. 새 frontend·장식 에셋·채팅버블·자유 composer·긴 타임라인·무조건 성공배너를 추가하지 않는다. 상태를 색만으로 구분하지 않으며 키보드/label/오류/긴ID 복사를 확인한다. 개인키/인증값/연결수단을 캡처하지 않는다.

## 갱신할 산출물과 기대 결과

Jev D10/D09/D05 및 실제 영향 D03/D06/D07/D08을 확인해 원천 문서를 갱신한다. 과제 로그는 .fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md다. 기술 문서·코드·실행 가능한 테스트·명령/exit증거·QA/UI 인계·변경/미검증/후속을 제공한다. front matter는 document-writing/deliverables 규약을 따른다.

## 제약과 완료 보고

Workers Free만 허용한다. 유료 플랜/구독/초과 과금·기존 서버/Tunnel 변경·운영배포·실메일 발송·실제 계정 외부 연결은 실행하지 않는다. 로컬 격리 fixture·비파괴 코드 구현·검증·커밋·일반 역할push는 허가됐다. 메시지 text 일반화·실제 송수신·receipt UX-06은 MESSAGES 후속이며 frozen relay.v1 wire를 임의 변경하지 않는다. 실제 플랫폼 지원을 문서 존재나 CLI PASS로 주장하지 않는다. 수집/삭제/권한/비용 등 범위 밖 행위와 제품 기준 불명확성만 질문하고 나머지는 완료까지 진행한다.

완료 보고에 고정 SHA/변경/검증한 것과 못한 것/의존성 변경 이유/SIZE·DESIGN 경고/산출물/후속을 적는다. work.py finish로 인박스와 결과 전문을 아카이브하고 빈 인박스를 확인해 커밋한다. preamble의 worker_done을 한 번 보내고 idle로 둔다.
