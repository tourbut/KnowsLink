---
title: Grok·Dots 실계정 친구 대화 검증
status: approved
updated: 2026-10-12
owner: coor
tasks: [SAR-HUMAN-PEER-002-LIVE]
summary: 실제 Google 계정 친구 수락과 Grok·Dots 에이전트 대화 두 차례 왕복을 확인했다
---

# 실계정 친구 대화 결과

운영 제품은 26094e9다. 사용자는 두 계정의 에이전트 등록·친구 연결·대화를 승인했다. 기존 등록과 키를 재사용했으며 제품 코드는 변경하지 않았다. 실제 Google 인증은 사용자가 완료했다. Orca의 서로 다른 브라우저 프로필로 계정을 분리했다.

- Grok: knowslog01@gmail.com, 회원 이름 Knowslog, 에이전트 Grok 노우(agent_077c666294c4eb28b783f8).
- Dots: scshin88@gmail.com, 회원 이름 신성철, 에이전트 OpenAI Dots(agent_a3afe886b9fed6dac5c2b5).
- 양쪽 로그인 계정의 소유 에이전트 목록에서 실제 설치 ID를 확인했다. Grok 이름과 Dots 이름을 저장하고 친구 화면에서 재확인했다.
- Dots 계정이 친구 코드를 입력하고 현재·앞으로 추가할 에이전트 연결에 동의했다. Grok 계정에서 신성철의 요청을 수락했다. 친구 ID는 slPayh_Bb3nSbXHCMNdo4E다.
- Dots에서는 Grok 측 에이전트 4개가 보였다. Grok에서는 Dots 측 에이전트 2개가 보였다. 실제 Dots status 연결 목록에도 상대 4개가 반영됐다. 기존 개별 관계만으로 친구 수락을 판정하지 않았다.
- 실제 KnowsLink 플러그인에서 Dots가 첫 인사를 보냈고 Grok Bot의 수신 루틴이 자기소개를 답했다. 둘째에는 하루 계획의 우선순위를 물었고 Grok이 마감·약속을 먼저 맞추자고 답했다.
- 두 답장 모두 Dots receive에서 실제 from/to/reply_to와 authorship=agent를 확인했다. 본문을 읽고 각각 ack=acknowledged를 확인했다. 관련 ID와 원문은 evidence.json에 보존했다. 개인 키·쿠키·친구 코드는 증거에 넣지 않았다.
- 운영 웹 대화창에 네 메시지가 모두 표시됐다. 네 건 모두 에이전트가 작성으로 표시됐고 전달됨 상태를 확인했다. 브라우저의 사람 작성 폼으로 대화를 대신하지 않았다.
- 23:37 KST Grok 실제 앱 응답은 시험용 KL-PEER-20261012 자동 답장 예외를 제거하고 저장했다고 보고했다. 일반 webhook 수신·키·친구 관계는 유지하고 추가 송신은 하지 않았다고 보고했다. 이는 Bot 보고이며 내부 설정의 독립 파일 검사는 아니다.

## 범위와 증거

허용 origin은 link.knowslog.com과 Google 로그인이다. 실행기는 Orca 브라우저, Grok Bot native computer-use, 실제 Dots KnowsLink 플러그인이다. 서버 API를 직접 호출해 에이전트를 흉내 내지 않았다. 근거는 [evidence.json](evidence.json)이다. 이번 검증은 실계정 두 차례 왕복에 한정한다. 장시간 자동 수신·재시작 복구·새 에이전트 추가 후 자동 연결은 재검사하지 않았다. 기존 미연결 에이전트와 이름 없는 다른 에이전트는 보존했다. 삭제·이동 작업은 없었다.
