---
title: SAR-SSE-INSTALL-001-TESTER — SSE 도착 신호·Google 인증·호스트별 설치의 독립 lite QA와 고정 SHA 리뷰
status: draft
updated: 2026-10-11
owner: tester
tasks: [SAR-SSE-INSTALL-001-TESTER]
summary: SSE 도착 신호·Google 인증·호스트별 설치의 독립 lite QA와 고정 SHA 리뷰
attempt: c3903c515d954261b24465fe35d7abef
base: c442bcaed63b8a2e01fb858a99c2852e23ff9962
subagent_level: standard
test_level: lite
---

# SAR-SSE-INSTALL-001-TESTER — SSE 도착 신호·Google 인증·호스트별 설치의 독립 lite QA와 고정 SHA 리뷰

- 작성일: 2026-10-11. From / To: coor / tester.
- 상태: blocked — DEV 최종 exact SHA와 실제 구현자 session 확정 뒤에만 dispatch한다.
- 테스트 / 하위 위임: lite / standard(선택형 동시2). 제품 소유권은 없다.
- 담당: repo 0b08ec4c-9e3a-4613-8197-5a835e545335, C:/Users/shin/orca/workspaces/KnowsLink/tester, fullops/tester.
- 복귀: C:/Users/shin/orca/workspaces/KnowsLink/coor, terminal term_e61d3e14-29e9-4954-943a-4a75707c82de, run_dc9b4e6dd18a. task/dispatch는 실제 preamble 정본.
- 병합 책임자 coor, main/origin. 운영 적용·실제 Google 로그인과 외부호스트 설치는 후보 수락 뒤 coor 담당이다.

## 현재 상황과 확인 근거

사용자가 기존 polling 보존, SSE/API 도착 신호 추가, 최소 리팩토링, Grok Bot·Claude Code/앱·Codex·OpenAI Dots 플러그인 설치와 Google 인증을 요청했다. DEV SAR-SSE-INSTALL-001-DEV는 진행 중이다. 구현 최종 SHA·실제 구현자 session·정본 기술 로그는 완료 뒤 여기 확정한다. 기존 자동 수신 수락은 이번 SSE/OAuth 지원을 대신 증명하지 않는다.

## 적용 기준과 예외

fullops-common-0.3.3 README/coding-style/testing/security, project.md, docs/agents/document-writing.md, rules/delegation.md, review/rule.json. 기준 규약 c442bca 및 후보 exact SHA. 보안·필수 검사·독립 QA/리뷰 생략 없음. Context7 quota exceeded 이력은 공식 문서로 대체한다. 이미 성공한 DEV 검사는 exact SHA·조건·의존성 동일성을 확인해 재사용하고, 독립 QA는 새 증거로 구분한다.

## 먼저 읽을 문서

- .fullops-squad/FULLOPS.md, project.md, rules/common/README.md 및 연결 세 규칙.
- docs/agents/document-writing.md, rules/delegation.md, review/rule.json, docs/deliverables/README.md.
- SAR-SSE-INSTALL-001-DEV 최종 archived 지시서·기술 계획·검증 요약(경로는 후보 확정 때 기록).
- 관련 정본 D03/D05/D10/D11/D12/D13 중 변경 절과 실제 구현 코드·호출자·검사.
- 탐색 find/context/packet은 후보 확정 뒤 생성·연결한다. 현재 템플릿은 배정 불가다.

## 해야 할 일과 파일 소유권

- [ ] 최종 정확한 SHA와 구현자와 다른 실제 reviewer session을 기록한다. fullops-review/open-code-review-delegate를 적용해 read-only clean detached snapshot에서 리뷰하고 결과는 tester 체크아웃에 쓴다.
- [ ] polling 회귀 및 SSE 신호→기존 API→서명 검증→private Inbox→persist/ACK 경로를 독립 검사한다. 원문·인증값 비노출, 다른 agent·잘못된키·철회 거부, TTL, 연결 종료/취소 cleanup, reconnect/catch-up, 중복/동시 pull·fallback의 위험 경계를 확인한다.
- [ ] 설치 bundle standalone 실제 MCP initialize/discovery와 Google 연결 계약을 검사한다. 호스트 manifest/공식 CLI 검증은 사용 가능한 환경에서 실행한다. 원 키·private folder·관계 보존을 확인하고 실제 호스트 설치 미실행을 PASS로 바꾸지 않는다.
- [ ] Google 기기 등록과 원격 OAuth2.1/PKCE/MCP 인증 계약의 차이를 검사한다. Dots 지원 주장을 공식 문서와 대조한다. MCP Events가 구현됐다면 구독 소유권·지속성/만료/철회·callback 검증·서명·SSRF/redirect/DNS rebinding·bounded retry·본문 미포함을 검사한다. 미구현이면 사용자 요구의 남은 범위를 명시한다.
- [ ] 공개 ingress가 필요한 경로만 여는지 확인하고 owner/admin/unknown 404 정책 보존을 검사한다. 운영 도메인 실제 SSE는 coor 후속 검사다.
- [ ] 관련 문서만 검토한다. 산출물 인덱스에서 대화 없이 설치/인증/수신/운영/복구/지원 한계를 찾을 수 있는지 의미적으로 확인한다. packet outcomes를 path/category별 작성한다.

제품 코드·하네스 규약·전역 호스트 설정·운영 DB/키/관계는 수정하지 않는다. QA 검사/시나리오/증거·리뷰 기록만 소유한다. 결함은 재현 근거로 보고하며 구현자에게 수정을 돌린다.

## 완료 기준과 검증

고정 후보의 required product-lint/product-test(make lint/make test), 실제 명령 종료코드·SHA·증거를 확인한다. 환경 한계로 Windows 미실행이면 기존 Linux 서버 고유 임시 clean checkout에서 격리 확인한다. local Docker/신규 과금 금지. 운영 서비스·DB·Tunnel 변경 금지. .env.server 값 출력 금지 및 기존 host-key 검증 유지. 필요한 임시 검사 자원만 만들고 회수한다.

독립 QA와 delegate 리뷰를 구분해 기록한다. 미해결 critical/high·필수 검사 실패·pending 파일이 있으면 수락 불가다. review.py check에는 실제 기준/후보 SHA, 원 DEV task-key를 넣고 모든 변경 파일을 reviewed/skipped와 이유로 채운다. 새 변경은 새 정확한 SHA 리뷰가 필요하다. 증거 밖 실제 호스트 wake·장시간·재부팅 검증을 주장하지 않는다. Dots/Claude 계정 권한이 없는 검사는 미실행으로 남기고 구현 지원 자체와 구분한다.

UI 변경은 기존 공용 템플릿/토큰을 대조한다. 해당 없음 또는 시각 검수 미실행 이유를 기록한다. 전체 디자인 재작성은 범위 밖이다.

## 갱신할 산출물과 기대 결과

D12 검증/운영 인계 근거만 필요한 절에 연결한다. 제품 문서 결함은 finding으로 보고한다. qa-reports/SAR-SSE-INSTALL-001-TESTER-test/와 별도 review 키의 보고서·result/lint, 필요한 scenarios를 작성한다. 새 QA 코드는 작은 검증 가능한 검사만 추가한다.

## 제약·완료 보고

컴퓨터 유즈는 Grok Bot 일회 설치 때만 coor가 사용한다. 메시지 송신 뒤 Grok UI로 수신을 확인하거나 유도하지 않는다. 임의 외부 메시지·새 유료 API·로컬 Docker·marketplace 제출 금지.

검증한 것/실패/미실행, 명령 exit code·대상 SHA, implementer/reviewer 실제 session, snapshot, finding severity·수락 여부, packet outcomes, 로그/산출물을 보고한다. work.py finish로 원문과 완료 보고를 보존하고 인박스를 비운 뒤 실제 preamble의 worker_done으로 복귀한다.
