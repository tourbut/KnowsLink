---
title: SAR-MVP-003-BIDIRECTIONAL — 실제 시험 메시지 양방향 연결
status: draft
updated: 2026-10-04
owner: dev
tasks: [SAR-MVP-003-BIDIRECTIONAL]
summary: 기존 relay와 MCP를 재사용해 Codex와 Grok Bot의 시험 메시지 송수신을 준비한다
---

# SAR-MVP-003-BIDIRECTIONAL — 실제 시험 메시지 양방향 연결

- 작성: coor→dev, ready. 기준 main 5c813a7979c3c2e230b546260221d1d18dd947d9.
- 담당 /home/shin/orca/workspaces/KnowsLink/fullops-dev, fullops/dev. 병합 coor/main/origin.
- 복귀 /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9, run_8ca8bc058ab7. task/dispatch는 preamble.
- 사용자 원문: 실제 운영 시험. 여기서 보내면 Grok Bot이 받고 Grok Bot 회신을 여기서 받는지 본다. 나중 OpenAI dots 연결. Grok 작업은 이슈 댓글로 전달.
- 승인: 시험 메시지 실제 네트워크 송수신 및 필요한 구현/테스트/일반GitHub push·issue댓글. 기존 실제 relay held 제한은 이 좁은 시험에 한해 새 요청이 우선. 업무발송/calendar/실데이터 조회/FullOps 업데이트/dots 연결 제외. 자격 증명 공개 금지.

## 현재 근거
issue1 최신5975152081: 부모 Grok의 AddMcpServer 등록 성공 user-knowslink connected/tools2, statusheld. 기존 mcp.ts는 held/합성loopback만 있고 임의메시지send 없음. core의 signed relay.v1와 Go relay pairing/ACK/claim 등을 재사용한다.
Grok 사전점검 댓글 https://github.com/tourbut/KnowsLink/issues/1#issuecomment-5975907733 게시됨. 부모 HTTPS/환경변경/자동실행 또는 수동pull 지원 회신 대기. 사용자는 공개relay 보유여부 async질문 답대기. 없는 것으로 단정하지 말고 구현 독립진행.

## 적용 기준·먼저 읽기
FULLOPS.md, rules/common/README.md와 연결3규칙(fullops-common-0.3.2), project.md, 문서작성규칙, fullops-work/ponytail full. 이 지시서는 과거 실제relay금지 기록보다 이번 사용자 시험승인을 우선하되 보안게이트는 낮추지 않는다.
관련 docs/D03/D05/D10 및 relay/MCP/auth만 좁게 읽는다. Jev code/doc/context 결과 keep/충돌 확인. SDK 새 API는 Context7 또는 공식문서, Grok공식 docs 확인. dots정체는 이번 조사/설계범위 아님.

## 해야 할 일
- [ ] 기존요구 내 짧은 기술계획과 최소 양방향 send/receive 구현. sender/recipient 식별·인증·서명·TTL·중복·원문 노출경계 정의. 일반 메시지 text를 허용된 당사자에게 전달하는 시험과 기존 disclosure/calendar intent를 구분한다. 무권한 라우팅/owner자격증명전달/SSRF·redirect/무제한응답을 방지.
- [ ] remote 실제 HTTPS relay 시험 경로, agent별 credential 준비/안전전달 절차와 Codex측 실행도구/수신 CLI 또는 MCP 준비. 기본 held는 보존. 승인된 명시모드만 시험메시지네트워크를 사용. timeout/비차단/자동wake와 수동pull 차이 검증. hostedloopback으로 remote성공 주장 금지.
- [ ] 로컬두독립agent 왕복을 실제 relay 저장소/MCP로 검증. send ID와 Grok수신·Grok회신·Codex수신을 대조할 구체절차. 실제Grok협업 전 로컬성공만 표기.
- [ ] 배포가능한 준비물/명령을 만들고 여기서 가능한 서버현황·credential유무를 비밀출력 없이 점검. 실제공개서비스에 신규비용/계정권한 필요하면 구체준비후 ask. 승인된서버정보가없으면 임의외부배포 대신 정확한차단조건과 대안을 보고.
- [ ] D03/D05/D10 및 필요한D12 구현정합갱신. 실행기록에 Grok이 실행할 이슈댓글초안(고정SHA자리/설치/안전env/호출/회신/실패절차). 원본기록 보존.
- [ ] 관련 회귀·auth거부·중복/TTL·crossagent·timeout 검증, lint --from 착수HEAD exit0, Jev score, work.pyfinish/archive/commit. 자기 inbox·code·기술문서 소유, 타인박스/PLANS/board 수정금지.

## 완료·후속
DEV완료는 독립fixedSHA OPS리뷰·GrokTESTER QA 후coor통합. actual Grok↔Codex 라운드트립은 별도 실제증거까지 미검증이다. 외부endpoint/계정 없다고 구현준비까지 멈추지 않는다. 제품규칙불명확은 ask, 기술해결은 dev책임. 완료 worker_done에 final SHA 전체40자리와 테스트결과/남은접속조건/댓글초안경로를 명시.

## 주소 확정·기존 배포 재사용
사용자 확인: link.knowslog.com으로 시험한다. 현재HTTPS HEAD는302 Cloudflare Access 로그인이다. 기존 SAR-BETA-001-OPS 실행기록과 deploy/knowslink 및 D12를 먼저 읽는다. 기존 knowslink tunnel/server/owner-only Access를 재사용하고 별도서비스 구축하지 않는다. agent용 machine 인증경로/최소권한을 준비하며 Access를 우회·해제하지 않는다. 기존배포현황은 직접확인하고 변경범위·검증을 기록한다. Grok댓글에도 주소확정추가 예정.

## Jev 근거·충돌 확인
find/documents-find/context: docs/evaluations/jev/SAR-MVP-003-BIDIRECTIONAL-*.json. keep 전부 읽는다. 지시 전제와 충돌 — 먼저 확인: SAR-BETA-001-OPS는 과거 owner-only 공개 구성이다. 이번 시험 승인으로 메시지 시험범위는 확대하나 인증·Access보호는 유지한다. 불명확한 제품권한변경만 ask한다.
