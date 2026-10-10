---
title: SAR-SSE-INSTALL-001-DEV — 폴링 보존·SSE 도착 신호 및 공통 Google 인증 플러그인 설치를 구현하고 기존 연결 수신 코드를 최소 정리한다
status: draft
updated: 2026-10-11
owner: dev
tasks: [SAR-SSE-INSTALL-001-DEV]
summary: 폴링 보존·SSE 도착 신호 및 공통 Google 인증 플러그인 설치를 구현하고 기존 연결 수신 코드를 최소 정리한다
attempt: 21a4cee08d734fa8a668bb0d785c0400
base: 7a6906851b872fc52310a0ac6623527a4b0adb2b
subagent_level: standard
test_level: lite
---

# SAR-SSE-INSTALL-001-DEV — 폴링 보존·SSE 도착 신호 및 공통 Google 인증 플러그인 설치를 구현하고 기존 연결 수신 코드를 최소 정리한다

- 작성일: 2026-10-11
- From / To: coor / dev
- 상태: ready
- 테스트 레벨 / 하위 위임: lite / standard(선택형 병렬 작업 최대2, 파일 소유권을 분리할 때만).
- 담당: repo 0b08ec4c-9e3a-4613-8197-5a835e545335, C:/Users/shin/orca/workspaces/KnowsLink/dev, fullops/dev.
- 복귀: 같은 repo/coor C:/Users/shin/orca/workspaces/KnowsLink/coor, terminal term_e61d3e14-29e9-4954-943a-4a75707c82de, run은 tracked worker-start preamble 정본. main 통합은 coor 담당.

## 현재 상황과 확인 근거

사용자 요청: 기존 polling은 유지하고 SSE/API 도착 신호 push를 추가한다. 구현을 리팩토링해 다듬고 플러그인 설치방식으로 Grok Bot·OpenAI Dots·Claude·Codex에 설치하며 Google 로그인으로 등록·인증한다. Dots 제품 링크는 https://openai.com/ko-KR/index/introducing-dots/ 이다. 해당 제품의 공식 지원 기능을 실제 확인한다.

기준7a6906851b872fc52310a0ac6623527a4b0adb2b에는 public-node MCP의 자동poll10s, private Inbox, Google beginLogin/finishLogin, signed member transport, Grok/Cursor plugin 패키지·Bot Command 설치기가 있다. 실제 노우 1차 자동 수신·관련답장 증거는 SAR-AUTO-RECEIVE-001-LIVE/report.md다. 기존 Google 계정·키·agent 관계는 보존한다. 현재 session용 SDK 수신은 process 생존 중만 동작했고 Codex 대화 자체의 자동 wake는 미구현이다.

## 적용 기준과 예외

fullops-common-0.3.3 README/coding-style/testing/security, project.md, docs/agents/document-writing.md, rules/delegation.md, Ponytail full. 기준7a69068 및 이 준비 commit의 지시서·packet. 예외 없음. Context7 resolve Cloudflare는 monthly quota exceeded였다. 공식 문서 및 설치된 SDK1.32.0으로 대체하며 버전/API 근거를 기록한다. OpenAI Docs를 적용한다.

## 먼저 읽을 문서

이 아래 목록과 packet의 mandatory부터 시작하고 전체 repo를 읽지 않는다. find/context 추천을 확인해 필요한 후보를 좁힌다.
- internal/relay/public_text.go
- internal/relay/connections.go
- internal/relay/http.go
- adapters/src/inbox.ts
- adapters/src/text.ts
- adapters/src/mcp.ts
- adapters/src/login.ts
- adapters/README.md
- scripts/package_plugin.py
- scripts/install_bot_mcp.sh
- scripts/check_public_ingress.py
- .fullops-squad/docs/design-docs/architecture.md
- .fullops-squad/docs/design-docs/interface-design.md
- .fullops-squad/docs/design-docs/module-design.md
- .fullops-squad/docs/operations/ops-guide.md
- .fullops-squad/docs/operations/user-guide.md
- .fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-LIVE/report.md

## 해야 할 일과 파일 소유권

- [ ] 같은 과제에서 짧은 기술 계획을 작성하고 구현한다. 기존 polling 경로·선택을 유지한다. SSE는 도착 신호만 보내며 본문 수신·서명 검증·private Inbox 기록·persist/ACK를 기존 API로 재사용한다. 수신 모드와 지원 상태를 명확히 한다.
- [ ] SSE 인증·자기 agent만 수신·현재키/관계/철회·TTL·원문/비밀 비노출·연결 자원 상한·heartbeat/timeout/취소·재접속/신호누락 복구·중복 신호 coalesce·직렬 pull·장애 fallback을 최소 구조로 다룬다. idle 상태에서 매번 DB polling하는 SSE 흉내는 피한다. 현 DB/relay singleton 모델에 맞는 가장 작은 안전한 구현을 선택한다. 신호는 수신권한이나 자동 실행/답장 권한이 아니다.
- [ ] 기존 연결·수신·설치의 중복/오류 경계를 필요한 부분만 정리한다. 전체 repo 재작성이나 광범위 스타일정리는 하지 않는다. 유효한 에러처리·검증·기존caller/tests/fixture/exports/config는 유지한다.
- [ ] Google 로그인 기존 flow를 재사용해 각 호스트/설치 인스턴스별 private credential과 별도 agent/key를 생성하도록 설치 경험을 정리한다. 기존 키/관계 덮어쓰기·credential을 plugin archive/env예시에 넣는 것은 금지한다.
- [ ] Grok Bot, Claude Code/Claude 앱, Codex, OpenAI Dots의 실제 공식 플러그인/MCP 설치계약을 확인한다. 최소 공통 bundle과 호스트별 필요한manifest/명령·설정만 제공한다. 네이티브 plugin 미지원이면 공식 Custom MCP 경로와 한계를 정확히 제공한다. Dots 미지원/공식문서부재라면 해당 항목만 명시해 coor에 질문하고 나머지 구현을 진행한다. 같은 설치물을 네 환경 모두 실제 설치했다고 주장하지 않는다. marketplace 공개 제출은 범위 밖이다.
- [ ] 자동 host turn은 호스트 기능과 별개임을 표시한다. 기존 Grok local gateway는 default-off/loopback/metadata-only를 유지하고 세션 wake를 보편 지원으로 표시하지 않는다.

제품 소유권: internal/relay/의 필요한연동과 tests, adapters/소스·manifest·skills·tests, scripts/설치·패키지·검사, 필요deploy ingress allowlist 및 기술·운영 정본. 외부호스트 전역 설정·사용자상설credential·DB/키·다른서비스·FullOps하네스규칙은 소유권 밖이다. 예상 수백줄~약1000줄; 실제 영향이 크면 기능단위 커밋으로 분리하고 SIZE-002 근거를 남긴다. 새 의존성은 최소화하고 기존/표준 라이브러리 대안 근거를 남긴다.

## 완료 기준과 검증

DEV는 최종 exact SHA에서 등록 product-lint/product-test(make lint/make test)를 실행한다. Windows make/cloudflared 제한 시 기존 Linux 운영 서버의 고유 /tmp clean clone으로 격리 검증한다. .env.server는 로컬 링크로 이미 제공한다. 값 출력 금지, 기존 host key verification 유지. local Docker는 실행하지 않는다. 운영 서비스/DB/Tunnel은 DEV 검사에서 변경하지 않는다. 필요 isolated test DB container는 기존 tester 정책과 동일하게 자기 고유 자원만 사용·회수하고 운영 DB와 혼용하지 않는다.

- SSE 신호에 원문·인증값이 없고 다른 agent/철회/잘못된키 접근이 거부되는 테스트.
- 기존polling 회귀, SSE연결/단절/재연결·pending catch-up·중복/동시pull·폐기/timeout cleanup의 핵심 자동 테스트.
- 설치아카이브는 standalone 실제 MCP initialize/discovery/held와 Google connect 계약을 검증한다. 호스트별manifest 및 공식 CLI검증은 가능한 환경에서 실행하고 미실행은 사유/영향을 남긴다. 기존zip payload·schema도 보존검증한다.
- 최종 고정 SHA의 독립 tester QA·읽기전용snapshot 리뷰는 coor가 별도 배정한다. DEV가 이를 대신하지 않는다.
- 실제 운영 도메인 SSE+Bot 왕복/Google 수동로그인/호스트설치는 후보수락 뒤 coor/ops가 진행한다. worker 완료와 제품 최종수락을 구분한다.

### UI 디자인

현재 Go template/login 화면과 기존 CSS만 필요 시 재사용한다. 별도 디자인 lint/테마전환 미구성, 적용불가 사유를 기록한다. 새 UI스타일이나 웹페이지는 불필요하다. 화면 변경 시 해당 정본·재사용·실제시각검수 대상을 완료보고에 명시한다.

## 갱신할 산출물

D03 architecture, D05 interface-design(신호계약), D10 module-design, D11 user-guide(설치·Google 등록), D12 ops-guide, D13 transition은 영향있는절만 갱신한다. packet document_update는 검토 뒤 실제변경/불필요근거를 기록한다. 새산출물13종 전체생성은 하지 않는다.

## 제약·완료 보고

추가 과금·유료 API·결제동의·local Docker·vendor core수정 금지. 기존 production deploy는 ba3754d로 보존중이며 후보수락 전에 바꾸지 않는다. 기존정상Google연결2개와 관계generation1·키를보존한다. 별도host 설치에는 새privatefolder를 사용하며 로그인은 사용자가 수행한다.

하위위임standard는 파일소유권이 겹치지 않을 때 최대2로 사용한다. 부모는 결과를 취합하며 Git커밋/브랜치조작을 단일 담당으로 유지한다. 선택형 위임 여부와 이유를 기록한다.

실제preamble로 worker_done을 보낸다. 도구별 공식근거(열어확인한페이지)·설치지원을 표로 기록한다. 정확한 최종SHA·필수검사·실패/미실행·새dependency·SIZE근거·packet path/category outcomes를 남긴다. work.py finish로원문을archive하고인박스를비운다. 보안검사나 미해결critical/high를 완화하지 않는다. 정상작업은 중간승인없이 완료한다.


공식 확인 근거: Dots 소개페이지에서 OpenAI 플러그인 생태계 연결을 확인했다(https://openai.com/ko-KR/index/introducing-dots/). OpenAI 설치/인증 계약은 https://developers.openai.com/plugins/build/plugins 및 https://developers.openai.com/plugins/build/auth 를 열어 확인했다. Claude manifest는 https://code.claude.com/docs/en/plugins-reference 를 열어 확인했다. 상세 schema/호스트별 제한은 worker가 해당 절을 확인해 구현한다. 기존 단순 Google 등록과 remote MCP OAuth2.1은 같은 계약이 아니므로 필요한 연결방식을 정확히 구분한다. context의 login.ts는 sensitive/oversized 로컬 확인 대상으로 유지하고, ops/user 문서 충돌 추천은 역사적 제한과 새 사용자 요청을 대조하며 규칙을 임의완화하지 않는다. server.go 부재는 실제 http.go로 보완한다.

추가 공식 근거: https://developers.openai.com/plugins/build/mcp-events 를 열어 확인했다. MCP Events는 ChatGPT Work/cloud desktop 및 Dots용 공식 이벤트 연동 문서다. 기존 모든 호스트 wake 불가라는 이력은 시점 한정이며 현재 지원을 다시 확인한다. 이 계약이 요청한 Dots 설치/도착신호에 적합하면 필요한 최소경로로 반영한다. 단순 relay SSE와 OpenAI MCP Events webhook subscription은 같은 것으로 표시하지 않는다. 구독/외부발송/실제호스트 이벤트 활성화는 사용자 승인범위와 추가과금 없는 조건 안에서만 진행한다. https://learn.chatgpt.com/docs/dots 와 https://developers.openai.com/plugins/deploy/connect-chatgpt 도 열어 확인했다. 문서 근거는 독립검토자가 재확인할 수 있게 URL·해당절·지원/미검증 상태를 남긴다.
