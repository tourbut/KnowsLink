---
id: D02
title: 경쟁 서비스 Reason Network 대비 KnowsLink 포지셔닝과 개선 우선순위
status: review
updated: 2026-10-22
owner: designer
tasks: [SAR-COMPETE-001]
upstream: [D01]
summary: "이슈 #9의 Reason Network 비교로 KnowsLink 차별점·핵심 가치 문장·개선 후보 우선순위·non-goal·SoT 반영안·사용자 결정 질문을 정한다. 사용자 결정 전 제안이다"
---

# 경쟁 서비스 비교와 KnowsLink 포지셔닝 — SAR-COMPETE-001

## 상태와 근거

이 문서는 사용자 결정 전 제안이다. 제품 규칙을 바꾸지 않는다. 사용자가 아래 '사용자 결정 질문'에 답한 뒤 제품 수락과 SoT 반영을 진행한다.

요청 원문은 [이슈 #9 원문](../SAR-COMPETE-001-request.md)이다. 기준 커밋은 main `30ceb2607fffd20a50d7154e5d5ca53d2643e67a`다. 반입 SoT는 service-design `7bc9ea190ea549fae8b047e850247a19322fc9c3`([source.json](../sources/silent-agent-relay/source.json))이다.

근거 표기 규칙은 다음과 같다.

- **사실(이슈)**: 이슈 #9 본문에 사용자가 정리한 내용이다. 이번 과제에서 웹 조사는 하지 않았다. 경쟁사 수치는 이슈 본문의 값만 쓴다.
- **사실(KnowsLink)**: 이 레포의 현행 D02·D03에서 확인한 내용이다. 링크를 함께 적는다.
- **판단**: designer의 제안이다. 사용자 결정 전에는 확정이 아니다.

이슈가 '확인하지 못한 것'으로 둔 세 항목은 미확인으로 유지한다. 세 항목은 Reason 호스팅 에이전트 메시지의 서버 보관 기간, Reason Machines의 투자 규모, HN·Product Hunt 반응이다. 이 문서는 세 항목에 기대는 결론을 만들지 않는다.

## 현행 KnowsLink 사실 요약

이슈의 비교표는 2026-10-22 현재 KnowsLink를 대체로 정확하게 적었다. 비교에 쓰는 현행 사실은 다음과 같다.

| 항목 | 현행 KnowsLink | 정본 |
|---|---|---|
| 관계 | 사람끼리 친구 코드로 요청하고 상대가 수락한다. 수락은 양쪽의 현재·향후 활성 agent 연결 동의다. 공개 회원 디렉터리는 없다. | [HP-02–06](SAR-HUMAN-PEER.md) |
| 해제·차단 | 해제는 즉시 새 전달을 막는다. 차단은 양방향 새 요청·전달을 막고 상대에게 차단 여부를 알리지 않는다. | [HP-13–14](SAR-HUMAN-PEER.md) |
| 발송 승인 | 웹은 '보내기'를 명시적으로 실행해야 보낸다. agent 자동 답장에는 별도 발송 승인이 필요하다. Claude Code 세션은 발송마다 소유자 대화상자를 거친다. | [HP-10](SAR-HUMAN-PEER.md), [SP-05–06](../../design-docs/SAR-SESSION-PEER-001.md), [D03](../../design-docs/architecture.md) |
| 위험 행동 | 결제·삭제·권한 부여·외부 발송·약속은 기본 거부하고 owner 게이트를 거친다. | [product.md 결정 4](../sources/silent-agent-relay/product.md) |
| 메시지 | 일반 text는 UTF-8 4096 bytes 이하다. 업무 요청은 서명된 `relay.v1` typed intent다. 첨부는 없다. | [HP-10](SAR-HUMAN-PEER.md), [protocol.md](../sources/silent-agent-relay/protocol.md) |
| 보존 | text 전달 기한은 180초다. 웹 대화 원문은 최대 24시간 보존한다. 백업은 별도이며 웹 대화의 24시간 뒤 사라짐은 백업 영구 삭제 약속이 아니다. | [HP-11–12](SAR-HUMAN-PEER.md) |
| 클라이언트 | 웹, Grok Bot(로컬 Command MCP), Dots(Remote HTTPS MCP), Claude Code 세션 연결이 있다. Codex 세션 수신은 미구현이다. | [LG-04](SAR-HUMAN-PEER.md), [D03 SP-12](../../design-docs/architecture.md) |
| 전송 | 집 미니서버의 중앙 HTTPS relay와 Cloudflare Tunnel이다. 종단 간 암호화는 없다. | [product.md 범위 메모](../sources/silent-agent-relay/product.md) |
| 상품 | Free N·Pro 가격·슬롯 단위는 미정이다. | [business-model.md](../sources/silent-agent-relay/business-model.md) |

## 1. KnowsLink만의 차별점

### 포지셔닝

**판단:** 두 제품은 같은 큰 개념을 다른 사용자에게 판다.

- Reason Network는 **한 팀 안의 코딩 에이전트 세션을 빠르게 잇는 개발 도구**다. 팀 전체 가시성, 실시간 busy/idle 목록, 세션 주소, P2P 전송이 가치의 중심이다(사실(이슈)).
- KnowsLink는 **서로 다른 사람의 개인 에이전트가 대화하는 개인 메신저**다. 사람 사이의 수락, 발송마다의 승인, 서버가 강제하는 위험 행동 차단, 짧은 보존이 가치의 중심이다.

KnowsLink는 개발 팀 협업 도구 시장에서 Reason Network와 기능 수로 경쟁하지 않는다. 신뢰가 먼저 필요한 사람 사이의 관계를 주 시장으로 둔다.

### 타깃 사용자

| 구분 | 사용자 | 이유 |
|---|---|---|
| 주 타깃 | 개인 AI 에이전트(Grok Bot, Dots, Claude Code 등)를 쓰는 개인. 친구·가족·지인의 에이전트와 짧은 요청·전달을 주고받으려는 사람. | 같은 팀 workspace가 없다. 상대를 개별로 수락하고 차단할 수 있어야 한다. |
| 보조 타깃 | 서로 다른 회사·팀에 속한 개인 사이의 일회성·소규모 협업. 예: 외부 협력자 한 명과의 에이전트 연락. | 팀 전체 공개형 초대가 맞지 않는다. 1:1 관계와 짧은 보존이 부담을 줄인다. |
| 비타깃 | 한 팀 내부에서 여러 코딩 세션을 실시간으로 조율하려는 개발 팀. | Reason Network의 팀 목록·세션 주소·P2P가 더 맞는다. KnowsLink는 이 수요를 따라가지 않는다(NG-02, NG-06). |

### 핵심 가치 문장

- **KR:** KnowsLink는 내가 수락한 친구의 에이전트와만, 내가 승인한 메시지만 주고받게 하고, 위험한 행동은 서버가 먼저 막는 개인 에이전트 메신저다.
- **EN:** KnowsLink lets your personal agent talk only with agents of friends you accepted, sends only what you approve, and has the server block risky actions by default.

이 문장은 SoT의 잠긴 공개 설명("Selective, human-gated handoff between paired agents")을 사람 친구 모델(HP-01–16)에 맞게 다시 쓴 것이다. 채택 여부는 Q1이다.

### 차별점과 근거

| ID | 차별점 | 사실 근거 | 판단 |
|---|---|---|---|
| DF-01 서버 강제 게이트 | 위험 행동을 서버에서 기본 거부한다. 승인 기록은 봉투 밖에 두고 GET 링크 승인을 금지한다. | 사실(KnowsLink): product.md 결정 4, C2. 사실(이슈): Reason의 `ask_user`는 문서 스스로 "speed bump, not a guarantee"라고 적었고 `--full-access` 모드가 있다. | 가장 강한 차별점이다. 승인 생략 모드를 만들지 않는다(NG-03). |
| DF-02 사람 단위 관계 | 연락처는 사람이다. 친구 코드 요청과 상대 수락, 해제·차단이 있다. | 사실(KnowsLink): HP-02–06, HP-13–14. 사실(이슈): Reason은 다시 쓸 수 있는 팀 초대 링크로 팀 전체가 서로 보인다. | 개인 관계에 맞는 동의 단위다. 팀 전체 공개형 초대를 추가하지 않는다(NG-02). |
| DF-03 증명 가능한 전달 | 메시지마다 JCS + Ed25519 서명, 멱등 처리, receipt가 있다. | 사실(KnowsLink): protocol.md, decisions.md. | 전달 사실과 발신 주체를 나중에 확인할 수 있다. 다만 서명은 암호화가 아니다. 공개 문구에서 둘을 섞지 않는다(C-08). |
| DF-04 설치 없는 웹과 호스팅 agent | 웹에서 바로 대화한다. 사용자 PC 없이 호스팅 agent(Dots 등)끼리 대화한다. | 사실(이슈) 'KnowsLink에만 있는 것'. 사실(KnowsLink): LG-03–04. | 개발자가 아닌 사용자의 진입 장벽이 낮다. 주 타깃과 맞는다. |
| DF-05 짧은 보존 | text 전달 180초, 웹 사본 최대 24시간이다. 장기 아카이브·검색·내보내기가 없다. | 사실(KnowsLink): HP-11–12. 사실(이슈): Reason은 받는 기기 SQLite에 저장한다. | 사적인 대화의 흔적을 줄인다. 다만 백업과 상대 클라이언트 사본은 별도라는 한계를 계속 정확히 쓴다. |

### 약점(사실)과 해석

| 약점 | 사실 근거 | 해석(판단) |
|---|---|---|
| 지원 클라이언트 수 | Reason은 10종이 넘고 자동 설정을 제공한다(이슈). KnowsLink는 웹·Grok·Dots·Claude Code이며 Codex 세션 수신은 미구현이다. | 수보다 '검증된 수신 수준'을 표시한다(NG-05). Codex는 C-01로 보강한다. |
| 종단 간 암호화 없음 | 이슈 비교표. 중앙 relay가 본문을 처리한다. | 서버 게이트와 웹 대화가 본문 접근에 기댄다. 현재는 투명 고지가 먼저다(C-08, Q2). |
| 첨부·그룹 없음 | 이슈 비교표. HP의 범위 생략 목록. | 주 타깃의 핵심 요구로 확인되지 않았다. 보류한다(C-06, C-07). |
| 가격 | Reason은 무료다(이슈). KnowsLink는 미정이다. | 숫자를 발명하지 않는다. 무료 경쟁 위험을 BM open note로 기록한다(Q6). |
| 개발 속도 | Reason은 하루에도 여러 번 릴리스한다(이슈). | 속도 경쟁은 목표가 아니다. 신뢰 경계의 회귀가 없는 것을 우선한다. |

## 2. 개선·보강 후보와 우선순위

### 기준

- 우선순위: **P0** 다음 과제로 바로 진행한다(코드 없음 또는 문서 정정). **P1** 다음 개발 후보다. **P2** 사용자 결정 뒤 조사·설계한다. **P3** 보류한다.
- 규모: **S** 문서·설정 한 과제 이하다. **M** 단일 기능 DEV 과제와 QA다. **L** 여러 과제 또는 구조·전송·보존 변경이다. 규모는 designer의 대략 추정이며 기술 계획은 DEV가 정한다.
- 원칙 충돌 열은 owner 게이트·서명·짧은 보존 세 원칙을 기준으로 쓴다. '없음'은 세 원칙을 바꾸지 않는다는 뜻이다.

### 후보 표

| ID | 후보 | 사용자 가치 | 원칙 충돌 | 규모 | 우선 | 담당 후보 | 근거 |
|---|---|---|---|---|---|---|---|
| C-01 | Codex 지원 확인과 Codex 세션 수신 | Codex 사용자가 친구 세션과 메시지를 주고받는다. 이슈가 지원 확인을 요청했다. | 없음. SP-05–06의 발송 승인·자동 답장 0을 그대로 적용한다. | M | P1 | dev(기술 판정·구현), tester | 설계 [SP-12·V02](../../design-docs/SAR-SESSION-PEER-001.md)가 이미 있다. Codex 세션 수신은 미구현이다(D03). 새 제품 결정 없이 진행할 수 있다. |
| C-02 | 루프 방지와 전송량 제한 규칙의 명문화 | 에이전트끼리 무한 답장으로 비용·소음이 생기지 않는다. | 없음. 오히려 게이트를 강화한다. | S | P0 | designer | 현행은 자동 답장 0(SP-06)과 발송마다 승인(HP-10), rate 한도(PS-11)로 루프를 구조적으로 막는다. 이 불변식이 SoT에 없다. 향후 자동 답장 허용 논의 전에 "승인 없는 연속 답장 0"을 잠근다. |
| C-03 | 클라이언트 자동 설정 확대 | 설치·연결 실패가 줄어든다. 신규 사용자가 빨리 첫 메시지를 보낸다. | 낮음. 키 복사·자격 공유·소유자 동의 없는 연결을 금지해야 한다(LG-04). | M | P1 | designer(흐름·카피), dev(설치기) | Reason의 10종 자동 설정이 이슈의 주요 비교점이다. 현재 Grok installer·Dots Remote MCP가 있다. 대상은 Claude Code와 C-01 뒤의 Codex로 한정한다. |
| C-04 | 상태 표시(수신 가능 여부) | 상대 agent가 지금 받을 수 있는지 알고 보낸다. 180초 만료 실패가 줄어든다. | 중간. SP-11은 실시간 작업 활동을 공유하지 않는다. HP-11은 온라인을 추측하지 않는다. 실시간 busy/idle은 충돌한다. | M | P2 | designer(공개 범위), dev | 현행은 HP-05의 '등록됨/연결 필요'뿐이다. 소유자 opt-in의 거친 '수신 가능' 한 값만 검토한다. 작업 경로·repo·branch는 공유하지 않는다(NG-06). 결정 Q3. |
| C-05 | 실행 중 세션에 즉시 표시 | 작업 중 세션이 다음 조회를 기다리지 않고 메시지를 본다. | 중간. 외부 본문을 소유자 지시로 모델에 넣으면 게이트를 우회한다. | M | P2 | dev | Claude Code는 도착 toast와 상태줄까지 있고 본문 모델 표시는 '다음 조회'다(D03 SP-12). untrusted 데이터 경계를 증명한 경로만 검토한다. |
| C-06 | 파일 첨부 | 짧은 문서·이미지를 함께 전한다. | 높음. 4096 bytes text 규칙, 24시간 보존, 미니서버 저장량, 악성 파일·자동 열기 위험과 충돌한다. | L | P3 | designer(범위), dev | 주 타깃의 핵심 요구 근거가 없다. 재개하면 작은 크기, 같은 24시간 보존, 자동 열기 금지, 발송 승인 화면에 파일 표시를 조건으로 한다. 결정 Q4. |
| C-07 | 그룹 대화(소수 친구) | 세 명 이상이 한 대화를 쓴다. | 높음. 1:1 친구 수락 모델, 발송마다 승인, 차단 비공개와 충돌한다. 그룹 안 비친구 노출 문제가 생긴다. | L | P3 | designer | 팀 전체 공개형 그룹은 non-goal이다(NG-02). 서로 모두 친구인 소수 그룹만 장기 후보로 둔다. 결정 Q4. |
| C-08 | 종단 간 암호화 여부와 투명 고지 | 사용자가 운영자가 무엇을 볼 수 있는지 정확히 안다. | 전체 E2E는 높음. 서버 게이트의 intent 검사와 설치 없는 웹 대화가 서버의 본문 처리에 기댄다. 고지만 하는 것은 충돌 없음. | S(고지) / L(E2E) | P1(고지) / P2(E2E 조사) | designer(고지 문구), dev(E2E 조사) | 서명은 무결성·발신 주체 증명이며 기밀성이 아니다. 공개 안내·정책 문구가 이 차이를 정확히 쓰는지 확인한다. E2E는 text 본문만 대상으로 하는 조사부터 시작한다. 결정 Q2. |
| C-09 | 경쟁 비교 문구의 공개 안내 반영 | 방문자가 KnowsLink를 고르는 이유를 바로 안다. | 없음. | S | P2 | designer | Q1 채택 뒤 D04 첫페이지 카피에 핵심 가치 문장만 반영한다. 경쟁사 이름은 공개 화면에 쓰지 않는다. |
| C-10 | SoT 포지셔닝·non-goal 반영 | 이후 제품 결정이 같은 기준을 쓴다. | 없음. | S | P0 | coor(사용자 승인 뒤 외부 SoT 반영) | 아래 '4. SoT 반영안'. 외부 레포 변경은 사용자 승인이 필요하다. |

### 우선순위 근거 요약

- P0(C-02, C-10)는 코드 없이 현재 규칙을 잠근다. 차별점 DF-01이 앞으로의 기능 추가에서 무너지지 않게 한다.
- P1(C-01, C-03, C-08 고지)은 이슈가 지적한 실제 격차 중 원칙과 충돌하지 않는 것이다. C-01은 기존 설계가 있어 가장 먼저 진행한다.
- P2(C-04, C-05, C-08 E2E, C-09)는 원칙과 일부 충돌하거나 사용자 결정이 먼저 필요하다.
- P3(C-06, C-07)는 원칙과 크게 충돌하고 주 타깃 수요 근거가 없다.

## 3. 하지 않을 것(non-goal)

| ID | 하지 않을 것 | 이유 |
|---|---|---|
| NG-01 | P2P 전송으로 전환하지 않는다. | 중앙 relay가 owner 게이트·친구 인가·차단·180초/24시간 보존을 강제하는 지점이다. P2P는 이 강제를 각 기기에 맡긴다. SoT의 HTTPS relay 코어 결정과 충돌한다. |
| NG-02 | 다시 쓸 수 있는 팀 초대 링크, 팀 전체가 서로 보이는 디렉터리, 와일드카드 주소(`alice/*` 형태)를 만들지 않는다. | HP-02는 공개 회원 디렉터리를 금지한다. 친구는 사람 1:1 수락으로만 생긴다. 와일드카드 주소는 '대상을 명시적으로 고른다'(HP-08)는 규칙과 충돌한다. |
| NG-03 | 승인을 생략하는 전체 허용 모드나 무인 headless 작업자를 만들지 않는다. | 발송마다 승인과 서버 기본 거부가 핵심 차별점(DF-01)이다. Reason도 자기 승인을 "speed bump"라고 적었다. 같은 약점을 들이지 않는다. |
| NG-04 | 장기 보관·검색·내보내기·대용량 첨부 저장소를 만들지 않는다. | 짧은 보존(DF-05)과 SoT 결정 3(장기 채팅 이력 없음)을 유지한다. 집 미니서버의 저장 부담도 피한다. |
| NG-05 | 검증되지 않은 클라이언트를 지원 목록에 올리지 않는다. 클라이언트 수로 경쟁하지 않는다. | 지원 표시는 실제 수신 수준(SP-12)의 근거가 있어야 한다. 기존 Grok/Dots 관측을 다른 클라이언트 지원 근거로 확대하지 않는다. |
| NG-06 | 친구에게 실시간 작업 활동(repo@branch, 작업 경로, busy/idle)을 공유하지 않는다. | SP-11은 작업 경로·제목·프롬프트·실시간 활동을 공유하지 않는다. 사적 관계에서 활동 노출은 개인정보 위험이다. |
| NG-07 | 사용자가 직접 띄우는 relay 배포판을 제공하지 않는다. | SoT는 병원·폐쇄망 self-host를 범위 밖으로 두었다. 직접 운영하는 중앙 서버가 게이트의 신뢰 근거다. |

## 4. SoT 반영안

### 적용 원칙

- 외부 레포 `tourbut/service-design`의 `silent-agent-relay/`에 이 문서가 직접 쓰지 않는다. coor가 사용자 승인 뒤 반영한다.
- 아래 제안은 고정 SHA `7bc9ea19` 기준의 절 이름이다. 반영 전에 coor가 외부 SoT의 현재 HEAD를 확인한다. 그 사이 절 이름이 바뀌었으면 같은 의미의 절에 반영한다.
- 반입 사본 `sources/silent-agent-relay/`는 수정하지 않는다. 외부 반영 뒤 반입을 새로 고치는 것은 별도 과제다.
- Q1·Q2·Q3에 묶인 문장은 사용자 답에 맞게 고친 뒤 반영한다. 묶인 질문을 각 행에 적었다.

### 파일별 제안

| 파일 | 바꿀 절 | 제안 문장 | 묶인 질문 |
|---|---|---|---|
| `product.md` | `**공개 제품명 (잠금):** KnowsLink` 아래 KR/EN 설명 다음에 새 소절 `### 포지셔닝 (2026-10-22)` 추가 | "KnowsLink는 서로 다른 사람의 개인 에이전트가 대화하는 개인 메신저다. 한 팀 안의 코딩 세션 버스와 경쟁하지 않는다. 핵심 가치: 내가 수락한 친구의 에이전트와만, 내가 승인한 메시지만 주고받고, 위험한 행동은 서버가 먼저 막는다." | Q1 |
| `product.md` | `## 범위 메모` 끝에 항목 추가 | "전송은 중앙 HTTPS relay를 유지한다. P2P 전환과 사용자 배포용 relay는 하지 않는다. 메시지 서명(JCS + Ed25519)은 무결성·발신 주체 증명이며 종단 간 암호화가 아니다. 공개 안내에서 둘을 구분한다." | Q2 |
| `decisions.md` | 문서 끝에 새 날짜 절 `## 2026-10-22` → `### Competitive positioning vs Reason Network — locked` 추가 | "Positioning: personal agent messenger between people who accepted each other, not a team coding-session bus. Keep: per-send human approval, server-side deny-by-default, 1:1 friend consent, short retention. Non-goals: P2P transport, reusable team invite links and team-wide directory, wildcard addressing, approval-bypass/full-access mode, headless unattended worker, long-term archive/search, sharing real-time work activity (repo@branch, busy/idle), user-hosted relay distribution. Source: SAR-COMPETE-001 (issue #9)." | Q1 |
| `decisions.md` | 같은 절에 `### Reply-loop invariant — locked` 추가 | "Default automatic replies = 0. Every agent-drafted reply needs a separate explicit owner send approval. A message, result, or approval request must not trigger another send without that approval. Existing HTTP/queue rate limits stay; do not invent new numbers here." | 없음(C-02) |
| `business-model.md` | `## Open notes — 리스크 가설` 목록에 항목 추가 | "**Free competitor:** 유사 기능의 무료 개발자용 서비스(Reason Network, 2026-10 공개)가 있다. 유료 가치는 개발 팀 기능이 아니라 사람 사이 신뢰 경계(수락·승인·서버 차단·짧은 보존)에서 찾는다. N·가격은 계속 TBD다." | Q6 |
| `mvp-checklist.md` | `## Explicit non-goals (MVP)` 목록에 항목 추가 | "~~P2P transport / user-hosted relay~~ — central relay enforces gate", "~~Team invite links / team-wide directory / wildcard address~~ — 1:1 friend consent", "~~Approval-bypass full-access mode / headless unattended worker~~ — per-send approval", "~~Real-time work activity sharing (repo@branch, busy/idle)~~ — privacy" | Q1 |
| `mvp-checklist.md` | `## Open / TBD` 목록에 항목 추가 | "[ ] Codex 세션 수신 지원 판정 — SP-12 수준별 근거", "[ ] text 본문 종단 간 암호화 조사 — 서버 게이트·웹 대화와의 양립 조건", "[ ] 거친 '수신 가능' 상태의 opt-in 공개 범위 — 실시간 활동 공유 금지" | Q2, Q3 |
| `protocol.md` | 변경 없음 | `relay.v1`은 freeze 상태다. 이번 결론은 wire를 바꾸지 않는다. 첨부·그룹 주소 필드를 추가하지 않는다. | 없음 |
| `README.md` | 변경 없음 | 포지셔닝은 `product.md`에 둔다. 중복 서술을 만들지 않는다. | 없음 |

`product.md`와 `decisions.md`의 2026-10-12 현행 안내(사람 peer가 Contacts=agents를 대체)는 이 레포의 반입 사본에만 있다. 외부 SoT에 같은 안내가 아직 없으면 coor가 사람 친구 모델 반영 여부를 함께 사용자에게 묻는다. 이 문서는 그 반영을 제안 범위에 넣지 않는다.

## 5. 사용자 결정 질문

| ID | 질문 | 선택지 | 추천 | 영향 |
|---|---|---|---|---|
| Q1 | 포지셔닝과 핵심 가치 문장을 채택하는가? | A. 제안 문장 그대로 채택 / B. 문장 수정 후 채택 / C. 개발자 작업 세션도 주 타깃에 넣는다 | A | A·B면 SoT `product.md`·`decisions.md` 반영과 C-09 진행. C면 NG-02·NG-06을 다시 검토해야 하고 Reason과 정면 경쟁이 된다. |
| Q2 | 종단 간 암호화를 어떻게 다루는가? | A. 현행 유지 + 서명≠암호화 투명 고지 / B. A + text 본문 E2E 조사 착수 / C. 전체 E2E 도입 | B | A는 C-08 고지만 한다. B는 고지와 DEV 조사(L 규모 판단 근거 확보)를 한다. C는 서버 게이트·웹 대화·24시간 웹 사본과 충돌해 핵심 규칙 변경이 필요하다. |
| Q3 | 상태 표시를 얼마나 공개하는가? | A. 현행 '등록됨/연결 필요'만 / B. 소유자 opt-in의 거친 '수신 가능' 한 값 추가 / C. 실시간 busy/idle 공유 | B(P2) | A는 변경 없음. B는 C-04 설계 과제를 연다. C는 SP-11·NG-06과 충돌한다. |
| Q4 | 첨부와 그룹 대화를 어떻게 하는가? | A. 둘 다 보류(P3) / B. 소형 첨부만 P2로 올린다 / C. 소수 친구 그룹만 P2로 올린다 | A | A는 현행 HP 범위 유지. B는 보존·저장량·악성 파일 규칙 설계가 필요하다. C는 그룹 동의·차단 비공개 규칙 설계가 필요하다. |
| Q5 | 다음 개발 순서는? | A. C-01 Codex 세션 수신 → C-03 자동 설정 / B. C-03 자동 설정 먼저 / C. 둘 다 보류 | A | A는 기존 설계(V02)를 바로 DEV에 배정한다. B는 Codex 지원 판정 전 설치 흐름을 먼저 만든다. |
| Q6 | 무료 경쟁 서비스에 대한 상품 방향은? | A. 개인 무료 유지, N·가격은 계속 미정 / B. 지금 전면 무료를 선언 / C. 유료화 시점을 지금 정한다 | A | A는 BM open note만 추가한다. B·C는 SoT BM 잠금 방향(Tailscale형 freemium)의 변경 결정이 필요하다. |
| Q7 | 특허 FTO(자유실시 검토)를 언제 하는가? | A. 유료화·상용 결정 전에 전문가 FTO / B. 지금 바로 / C. 하지 않음 | A | A는 SAR-PUBLISH-001의 P-05 법률 검토와 함께 일정을 잡는다. B는 비용이 먼저 든다. C는 상용화 위험을 남긴다. |
| Q8 | SoT 반영안을 외부 레포에 적용하는가? | A. Q1–Q6 답 반영 후 적용 / B. 적용하지 않고 이 레포 D01·D02에만 둔다 | A | A는 coor가 외부 SoT를 바꾸고 반입을 새로 고친다. B는 이후 SoT 기준 판단에서 이 결론이 빠진다. |

### 특허 메모

이 메모는 법률 자문이 아니다. 이슈 본문에 따르면 넓은 개념 자체를 막는 특허는 찾지 못했다. 관찰 대상은 US20250373574A1(OpenAI, 심사 중)과 US12671588B1이다. 이 과제에서 특허를 추가 조사하지 않았다. 상용화 전 전문가 FTO 여부와 시점은 Q7의 사용자 결정이다. SAR-PUBLISH-001의 [P-05 법률 검토](SAR-PUBLISH-001.md)와 함께 진행할 수 있다.

## 완료 조건과 인계

- designer 완료는 이 문서, [D01 경쟁·포지셔닝 절](../business-plan.md#경쟁과-포지셔닝--2026-10-22), lint 통과까지다.
- 제품 수락은 사용자가 Q1–Q8에 답한 뒤다. 답에 따라 이 문서의 우선순위·SoT 문장을 고친다.
- SoT 반영과 이슈 #9 댓글·종료는 coor가 사용자 승인 뒤 처리한다.
- C-01·C-03 등 개발 후보는 사용자 결정 뒤 별도 과제로 배정한다. 기술 계획·규모 확정은 DEV 책임이다.

## 개정 이력

- 2026-10-22: SAR-COMPETE-001에서 이슈 #9 비교를 바탕으로 초안을 작성했다. 사용자 결정 전이다.
