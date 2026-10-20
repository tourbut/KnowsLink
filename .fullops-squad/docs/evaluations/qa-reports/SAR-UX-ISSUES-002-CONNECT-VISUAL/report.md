---
title: 연결 안내 4단위 후보의 직접 시각 판정과 재요청 행 결정
status: review
updated: 2026-10-20
owner: designer
tasks: [SAR-UX-ISSUES-002-CONNECT-VISUAL]
summary: "후보 8d8bad4의 연결 경로 선택·재개·보낸 요청 끝 4자·연결 화면을 1440·390·320 캡처로 판정한다. 결론은 조건부 수락이며 재요청 행 차이는 '기존 행 재사용'으로 정한다"
---

# SAR-UX-ISSUES-002-CONNECT-VISUAL 직접 시각 판정

**판정은 조건부 수락이다.** 직접 본 모든 화면은 U8·나머지 U7·O-NM-03 범위와 D04를 충족한다(시각 PASS). 조건은 V-CN-01 하나다. 끝난 요청을 같은 코드로 다시 보내면 차단·수신 한도 경우에만 요청자 행이 하나 더 생긴다. 이 차이는 U4의 존재·차단 비노출 규칙을 어긴다. 제품 판단은 ‘기존 행 재사용’이며 [U4](../../../design-docs/SAR-UX-ISSUES-002.md)에 반영했다. critical·high는 없다.

대상은 고정 후보 `8d8bad4cd75a37546ba0f803e3c0c8290ddb6a87`(fullops/dev, 제품 `37febf9`)이다. 지시서 기준 main은 `d3e676a`, 지시서 커밋은 `4d15c77`이다. 제품 소스는 dev 워크트리(HEAD `8d8bad4`)에서 읽기만 했다. 제품 코드는 고치지 않았다.

## 증거와 방법

- fixture: origin `http://127.0.0.1:18932`(coor 포워딩, 컨테이너 `kl-connect002-fixture-8d8bad4cd75a-a27ac964`). 연결 ID는 manifest `C:/Users/shin/AppData/Local/KnowsLinkDevTools/connect-ux002-dev/fixture-manifest.json`에서 읽었다.
- 요청: 허용 GET(`/home/agents`(query 포함), `/home/friends`, `/home/legacy`, `/home/account`, `/home/connections/*`, `/connect/*`)과 세션용 `POST /fixture/as/{a|b|c}`만 보냈다. 상태를 바꾸는 POST(생성·연결·취소·친구 요청·결정)는 보내지 않았다.
- 도구: [capture.mjs](capture.mjs). headless Chrome(비모바일 viewport, 배율 1)으로 1440·390·320 전체 높이 캡처를 만든다. 측정값은 페이지 가로 넘침, 44px 미만 동작, main 안의 가운데 정렬, ‘온라인’ 문구와 배지, 대상 외 호스트 이름, `aria-current`다. 결과는 [captures/probe.json](captures/probe.json)(51장, `NO_OVERFLOW 51`)과 [captures/probe-connections.json](captures/probe-connections.json)(18장, `NO_OVERFLOW 18`)이다. 두 실행 모두 종료코드 0이다. 1440 캡처마다 본문 텍스트(`*.txt`)를 함께 남겼다.
- remote 앱 복귀 화면은 동의 POST 뒤에만 열려 허용 GET으로 만들 수 없다. DEV 렌더 `C:/Users/shin/orca/workspaces/KnowsLink/dev/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-DEV-test/render/`(아래 `R/`)의 `remote-return-{1440,390,320}.png`를 직접 열었다.
- D04 스타일: 후보 diff(`d6f9d7f..8d8bad4`, 테스트 제외)에 새 색·그림자·그라디언트·모서리 규칙이 0건이다. 새 표시는 기존 `badge` 함수(글자+기호)와 `nowrap` 묶음만 쓴다.
- 운영 기록: 첫 실행은 세션 쿠키(`__Host-`, Secure)를 headless Chrome이 http 응답에서 저장하지 않아 비로그인 GET이 반복됐고 fixture 요청 제한(429)에 걸렸다. 이 실행의 PNG는 모두 덮어썼다. 둘째 실행부터 쿠키를 CDP로 넣고 4초 간격·429 뒤 30초 대기로 바꿨다. tester가 같은 계정을 쓰는 동안 회원별 40회/분 제한을 나눠 썼다. 같은 시간대 tester 429의 일부는 이 원인일 수 있다.
- 브라우저 프로필: `D:/workspace/KnowsLink/.git/worktrees/designer/connect-visual-profile`(이 워크트리의 Git 디렉터리). 저장소 규칙(삭제 작업 금지)에 따라 지우지 않고 보존한다. 정리는 coor가 정한다.

## 화면별 판정

69장 모두 페이지 가로 넘침 0, 44px 미만 동작 0이다. main의 가운데 정렬은 버튼 안의 라벨뿐이다. ‘온라인’ 배지는 0건이다. ‘온라인’이라는 단어는 ‘온라인 여부를 확인하지 않습니다’ 같은 설명 문장에만 있다. OpenClaw·Meta Muse·Claude Desktop·Cursor·Grok CLI 문자열은 모든 화면에 0건이다. 캡처 경로는 `captures/` 기준이다.

| 화면 | 근거 캡처 | 판정 | 근거 |
|---|---|---|---|
| ‘어디에서 사용할까요?’ 5선택지 | `a-home-agents-{1440,390,320}.png` | PASS | 순서·문구·보조 문구가 U8 선택지 표와 같다. radio 라벨은 44px 행이다. 설명 ‘고르기만으로 에이전트나 연결 수단을 만들지 않습니다’가 있다. 320에서 4번 보조 문구가 두 줄로 줄바꿈되고 잘리지 않는다. |
| 1 웹에서 바로 사용 | `a-home-agents-use-web-start-first-*.png` | PASS | ‘설치 필요 없음’과 발신 이름용 에이전트 하나, 친구 수락·매 발송 확인을 보인다. 다음 행동은 ‘친구 확인·추가’, ‘새 대화’다. |
| 2 Grok Bot | `a-home-agents-use-grok-start-first-*.png` | PASS | 준비물(OS·도구·Command MCP 등록 권한·앱 버전), 권한, ‘Grok 신규 설치 문서’ 링크. 설치 성공 표시 없음. |
| 3 Dots | `a-home-agents-use-dots-start-first-*.png`, `b-home-agents-use-dots-start-continue-*.png` | PASS | 계정/workspace 권한 확인, 로컬 설치 파일 없음, 서버 보관 범위 안내, ‘Dots에서 연결을 시작하세요’. 준비물 문구는 README의 기존 Dots 절과 같다. |
| 4 Claude Code 작업 세션 | `a-home-agents-use-claude-start-first-*.png` | PASS | 저장소·Node 22.22.2·설치 권한, 시험 환경 2.1.302·Windows, `/knowslink-session on`. ‘세션 연결: 앱에서 확인 필요’와 `/knowslink-session status`, 도착 알림은 알림이며 본문은 `knowslink_text_receive`라는 설명이 U8 규칙과 같다. |
| 5 지원 확인이 필요한 호스트 | `a-home-agents-use-other-start-first-*.png` | PASS | ‘아직 연결 방법을 확인하는 중입니다. 지금은 웹에서 대화를 이용하세요.’와 ‘웹에서 바로 사용’ 링크만 있다. 설치 버튼·명령·배지·도구 이름 0건. |
| 처음 / 기존 연결 계속 | `a-home-agents-use-grok-start-continue-*.png` | PASS | ‘기존 연결 계속’은 현재 에이전트와 키를 먼저 확인하라고 안내하고 새 에이전트·설치 폴더를 만들지 않는다고 밝힌다. ‘연결 상태로 이동’이 같은 화면의 상태 절로 간다. |
| 상태별 재개 행(local) | `a-home-agents-*.png` | PASS | waiting ‘키 등록 전 · 내 컴퓨터의 준비 대기’, prepared ‘지문 확인 필요’, approved ‘브라우저 동의 완료 · 클라이언트 저장 확인 필요’, consumed ‘키 등록됨’과 저장 실패 시 덮어쓰지 말라는 안내. 행마다 다음 행동이 하나다. 보낸 친구 요청 대기 수와 요청 상태 조회 안내가 있다. |
| 상태별 재개 행(만료·취소·remote) | `b-home-agents-{1440,390,320}.png` | PASS | 만료·취소 ‘키 없음’과 새로 시작, device 승인 뒤 미저장 ‘확인 필요’와 자동 철회 없음, remote 앱 복귀 대기·교환 유실(‘확인 필요’)·교환 완료가 구별된다. 앱 이동만으로 연결되지 않는다고 밝힌다. |
| 에이전트 0개 | `c-home-agents-*.png` | PASS | ‘첫 에이전트를 추가하세요’와 추가 폼. 빈 선택기 없음. |
| 연결 화면 단계 | `a-home-connections-{OARY…,e4IB…,jNyK…,0rwi…}-*.png` | PASS | `aria-current=step`이 한 단계에만 있다(waiting 1, prepared 2, approved 3 ‘연결 마무리’, consumed 3 ‘키 등록됨’). consumed는 ‘✓ 등록됨 키 등록됨. 앱에서 연결 상태를 확인하세요…온라인 여부를 확인하지 않습니다.’다. 복귀 링크는 ‘내 에이전트로 돌아가기’다. 선택 단계는 accent 왼쪽 선과 굵은 글자로도 구별된다. |
| 연결 화면 만료(b) | `b-home-connections-*-*.png` | PASS | 만료와 재사용 불가, 내 에이전트에서 새로 시작. `aria-current` 없음(확인된 진행 단계가 없음). |
| 보낸 요청 ‘입력한 코드 …abcd’ | `a-home-friends-{1440,390,320}.png` | PASS | 대기 행 ‘입력한 코드 …zz99’·‘…XXK2’, `ab` ‘입력한 코드 형식 다름’, 이전 요청 ‘이전 요청’, 취소 기록 ‘— 취소 / 입력한 코드 …ca11’. 있는 코드와 없는 코드 행의 모양이 같다. 320에서 라벨이 `nowrap`으로 한 덩어리로 줄바꿈된다. 받은 요청(C) 행에는 코드 일부가 없다. 취소 버튼 접근성 이름은 템플릿상 ‘{라벨} 친구 요청 취소, 기한 …’이다. |
| remote 앱 복귀·실패 안내 | `R/remote-return-{1440,390,320}.png` | PASS | ‘이 화면은 연결 완료를 뜻하지 않습니다’, 1분 안의 검증된 앱 링크, 만료 시 앱에서 다시 시작. 페이지 가로 넘침 0(DEV probe). |
| legacy ‘이 화면에서 찾기’ | `a-home-legacy-{1440,390,320}.png` | PASS | 앵커 목록이 `nav aria-label="이 화면에서 찾기"`이며 `aria-current`가 0건이다. 선택선 같은 현재 위치 표시가 없어 사람 주 메뉴와 구별된다. |
| 헤더 ‘내 계정’ | 모든 회원 화면 | 범위 밖 | ‘설정’ 이름 변경은 5단위(U3)다. 지시서 to_dev와 같다. |

## 재요청 행 차이의 제품 결정

DEV가 남긴 질문은 차단된 요청이 만료된 뒤 같은 코드로 다시 보내면 새 행이 생기고, 없는 코드는 기존 행을 다시 쓰는 차이다. 소스(`internal/relay/human.go:requestFriend`, `8d8bad4`)를 읽어 같은 차이가 두 경우 더 있음을 확인했다.

| 이전 상태 → 같은 코드 재요청 | 저장 행의 동작 | 요청자 화면 |
|---|---|---|
| 있는 코드, 만료·취소·거절 | 같은 쌍 행을 다시 대기로 바꾼다(세대 +1) | 행 1개 |
| 없는 코드·자기 코드, 만료·취소 | 같은 대리 행을 다시 대기로 바꾼다 | 행 1개 |
| 받은 사람이 대기 중 차단, 기한 뒤(또는 요청자가 취소한 뒤) | 차단 행은 남고 대리 행을 새로 만든다(세대 1) | 행 2개(옛 만료/취소 행 + 새 대기 행) |
| 상대 수신 대기 10건 이상, 이전 쌍 행이 끝남 | 쌍 행은 남고 대리 행을 새로 만든다 | 행 2개 |

행 수와 새 행의 위치·세대 값으로 ‘이 코드는 실제 회원이며 나를 차단했다’ 또는 ‘수신 한도에 걸렸다’가 드러난다. U4는 차단 여부를 상대에게 알리지 않고 존재 여부로 결과를 구분하지 않는다고 정한다.

**결정: 기존 행 재사용.** 같은 코드로 다시 보내면 그 코드의 기존 요청자 행 하나를 다시 대기로 바꾸고 새 행을 추가하지 않는다. 있는/없는/차단/수신 한도/자기 코드에서 행 수·순서·상태·끝 4자·새 24h 만료·취소 폼 숨은 값이 같아야 한다. ‘모두 새 행’ 안은 기존 있는/없는 코드 동작까지 바꾸고 종료 행이 쌓여 채택하지 않았다. 세부 규칙은 U4의 ‘2026-10-20 CONNECT-VISUAL 결정’ 절이다. 저장 방식은 DEV가 정한다.

## DEV 수정 필요

| ID | 등급 | 내용 | 수락 조건 |
|---|---|---|---|
| V-CN-01 | medium | 위 표의 차단(기한 뒤·취소 뒤)과 수신 한도 경로에서 재요청이 새 행을 만든다. U4 결정대로 요청자 화면을 기존 행 재사용과 같게 한다. 보낸 요청 상한·전체 저장 상한 판단도 존재·차단·수신 한도로 달라지지 않게 한다. | 자동 검사로 다섯 경우의 요청자 행 수·상태·라벨·세대 증가가 같음을 단언한다. 차단한 사람 화면에 새 요청이 없다. 화면 모양은 바뀌지 않으므로 designer 시각 재검수는 필요 없다. tester 확인으로 충분하다. |
| V-CN-02 | low | `GET /home/connections/{없는 ID}`의 거부 화면 복귀가 ‘홈으로 돌아가기’다(`a-home-connections-agent-fx-registered-*.png`). U7은 진입 문맥이 없으면 해당 안전 목록(키 연결은 내 에이전트)으로 돌아가게 한다. | 거부 복귀를 ‘내 에이전트로 돌아가기’로 바꾼다. 다음 단위와 묶어도 된다. |

비차단 관찰(수정 선택): 1440에서 ‘안내 보기’ 버튼 바로 아래에 선택한 곳의 제목이 붙어 위아래 간격이 좁다(`a-home-agents-use-grok-start-first-1440.png`). 읽기 순서는 바르다.

## 실행하지 못한 것 (not_run)

- 실제 키보드 순회, 브라우저 실제 200% 확대, 스크린리더, 보조 브라우저(U7-A). tester 산출물 `C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-TESTER-test/`에 Chrome·Edge 키보드·레이아웃 기록이 생기고 있으나 이 판정에 쓰지 않았다.
- 상태를 바꾸는 흐름(생성·중복 클릭·취소·친구 요청)의 화면 전이. 지시서대로 POST를 보내지 않았다. 재요청 행 차이는 소스로만 확인했다.
- remote 앱 복귀 화면의 fixture 캡처(허용 GET 아님). DEV 렌더로 판정했다.
- 실제 호스트 설치·인증·왕복. 범위 밖이다.
