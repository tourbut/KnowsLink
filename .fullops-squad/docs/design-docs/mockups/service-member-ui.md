---
id: D04
title: 사람 친구와 에이전트 대화의 화면설계서
status: review
updated: 2026-10-27
owner: designer
tasks: [SAR-SERVICE-UI-001-DESIGN, SAR-SERVICE-UI-001-VISUAL, SAR-FRONTEND-POLISH-001-VISUAL, SAR-HUMAN-PEER-001, SAR-HUMAN-PEER-001-VISUAL, SAR-LANDING-GUIDE-001, SAR-LANDING-GUIDE-001-DEV, SAR-SCHEDULE-001-VISUAL, SAR-PUBLISH-001, SAR-PUBLISH-001-D04, SAR-PUBLISH-002-DOCS, SAR-AUTOSETUP-001, SAR-AUTOSETUP-001-ART, SAR-CONNECT-SIMPLE-001, SAR-CONNECT-SIMPLE-001-UI-FIX-DEV]
summary: 이름과 대화 중심의 정보 구조·토큰·반응형 화면·상태 및 후속 수락 기준
---

# 회원 웹 상용 UX·비주얼 명세와 개발 인계

## 공개 첫페이지 정본 — SAR-LANDING-GUIDE-001

2026-10-14 최신 사용자 요청은 첫페이지에 서비스 소개·설치·사용 순서·실제 기능에 맞는 예시와 로그인 시작을 요구한다. 이 절은 공개 진입 화면에 한해 아래 기존 최소 로그인 설계보다 우선한다. 이전 '세 줄 서비스 소개 삭제'와 빈 회원 홈의 '가짜 예시 대화 금지'는 원문을 보존한다. 공개 안내에 명시된 설명용 예시는 허용하지만 회원 대화 목록에는 넣지 않는다. 회원의 사람·친구·agent·대화 구조와 권한은 바꾸지 않는다.

기준 제품·문서는 착수 HEAD `93271bf0d2aa8b6677a1f3d3135dc5dd708cddf8`에서 읽었다. 현재 `startPage`는 유효 세션을 `/home`으로 보낸다. 새 설계에서는 로그인 상태와 무관하게 `/`에서 안내를 읽고, 유효 세션이면 '대화로 이동'으로 회원 공간에 들어간다. 이 변경은 의도된 제품 동작이다. 로그인 완료·이름 미설정·연결 동의 복귀는 기존 흐름을 유지한다. 상세 기술 방법은 DEV 책임이다.

### 목표·정보 구조·진입

대상은 자기 에이전트를 친구의 에이전트와 연결하려는 첫 방문자다. 첫 화면에서 서비스 목적과 로그인 다음 행동을 알 수 있어야 한다. 설치를 준비하는 사람과 이미 연결한 회원도 같은 안내를 다시 읽는다.

| 순서 | 섹션 / 앵커 | 방문자가 얻는 답 | 행동 |
|---|---|---|---|
| 1 | 소개 / `#main` | 누구와 무엇을 주고받는 서비스인가 | 로그인하고 시작하기, 시작 순서 보기 |
| 2 | 시작 순서 / `#how` | 로그인부터 친구 대화까지 어떤 순서인가 | 에이전트 연결 안내로 이동 |
| 3 | 설치와 연결 / `#install` | 웹·Grok·Dots에서 무엇을 준비하는가 | 호스트별 절차 펼치기, 상세 설치 문서 |
| 4 | 대화 예시 / `#examples` | 실제로 어떤 메시지를 보낼 수 있는가 | 설명용 예시 두 개 읽기 |
| 5 | 시작하기 / `#start` | 지금 어떻게 로그인하는가 | 활성 로그인 수단 또는 대화로 이동 |

헤더는 KnowsLink, '시작 순서', '설치와 연결', '대화 예시', '로그인'을 제공한다. 로그인 상태에서는 마지막 항목만 '대화로 이동'으로 바뀐다. 브랜드는 `/`로 이동한다. 앵커 링크는 해당 제목으로 이동하며 로그인이나 API 호출을 시작하지 않는다. 하단에는 '개인정보처리방침', '서비스 약관', '설치·사용 문서', '맨 위로'를 이 순서로 둔다(2026-10-21 SAR-PUBLISH-001, 정책 링크 추가). 회원 공간의 '내 계정'에 '서비스 안내' 링크를 추가한다. 회원의 주 이동 메뉴 세 개는 유지한다.

### 실제 카피와 정보의 깊이

다음 따옴표 안 문구는 구현 카피다. 괄호 속 구현 설명·과제 ID·SHA·검수 상태는 사용자 화면에 출력하지 않는다.

**소개**

- h1: '친구와 연결하고, 에이전트로 대화하세요'
- 본문: 'KnowsLink는 서로 친구가 된 사람들의 에이전트를 연결합니다. 웹에서 직접 메시지를 쓰거나, 연결한 Grok·Dots에서 메시지를 주고받으세요.'
- 보조: '웹은 설치 없이 시작할 수 있습니다. 에이전트를 연결하고 친구가 수락하면 대화를 시작할 수 있어요.'
- 비로그인 주 행동: '로그인하고 시작하기' (`#start`). 보조 행동: '시작 순서 보기' (`#how`).
- 로그인 주 행동: '대화로 이동' (`/home`). 보조 행동은 같다.
- 오른쪽 관계 도식 제목: '사람끼리 친구, 에이전트끼리 대화'. 사람 '나' 아래 '내 에이전트', 사람 '친구' 아래 '친구의 에이전트'를 들여쓴다. 가운데에는 '서로 수락한 친구'를 적는다. 온라인 표시·가짜 수치·실제 회원 아바타는 없다.

**시작 순서**

h2는 '이 순서로 시작하세요'다. 실제 순서이므로 번호를 사용한다. 카드 네 개 대신 하나의 세로 목록으로 이어 놓는다.

1. '로그인하고 내 이름 정하기' — '친구에게 보일 이름을 정하세요. 이름이 같아도 같은 계정은 아닙니다.'
2. '내 에이전트 연결하기' — '내 에이전트에서 연결할 클라이언트를 확인하고 이름을 정하세요. Grok와 Dots는 각각 연결합니다.' 하단 링크는 '설치와 연결 방법'이다.
3. '친구 요청하고 수락하기' — '상대가 공유한 친구 코드로 요청하세요. 서로 수락하면 두 사람의 현재·앞으로 추가할 활성 에이전트가 연결됩니다.'
4. '대상을 고르고 메시지 보내기' — '새 대화에서 친구, 내 발신 에이전트, 상대 수신 에이전트를 고르세요. 대상과 본문을 확인한 뒤 보내세요.'

목록 뒤에 '내 에이전트끼리는 친구 요청 없이 대화할 수 있습니다. 기존 회원은 내 에이전트 전체 연결에 한 번 동의해야 합니다.'를 표시한다. 친구 이름·이메일 검색이나 공개 회원 목록이 있는 것처럼 안내하지 않는다.

**설치와 연결**

h2는 '사용할 곳에 맞게 연결하세요'다. 첫 문장은 '웹에서 이름·친구·대화를 관리할 때는 별도 설치가 필요 없습니다. 에이전트를 연결할 때 아래에서 사용하는 클라이언트를 고르세요.'다. 웹만 열어도 연결된 agent 없이 즉시 송신할 수 있다는 뜻으로 표현하지 않는다.

Grok와 Dots는 각각 한 개의 독립 접기 영역이다. 둘 다 처음에는 닫혀 있다. 선택형 탭·자동 호스트 감지·중첩 접기는 없다. summary와 바로 아래 한 문장 요약은 접힌 상태에서도 읽힌다. 두 영역을 동시에 열 수 있다.

| 항상 보이는 제목 | 항상 보이는 요약 | 펼친 뒤 순서와 카피 |
|---|---|---|
| 'Grok에 연결하기' | 'Bot 컴퓨터에 설치하고 Command MCP로 연결합니다.' | ① 'Command MCP를 등록할 수 있는 Grok 환경인지 확인하세요.' ② 'Bot 컴퓨터에서 설치 문서에 따라 bundle을 준비하세요. Linux x86_64 또는 arm64 환경이 필요합니다.' ③ '설치 결과의 Command와 Arguments를 앱에 등록하세요.' ④ 'Grok에서 KnowsLink 연결을 시작하고, 내 브라우저에서 Google 로그인과 지문 확인 후 동의하세요.' ⑤ '연결 완료와 웹의 내 에이전트를 확인하세요.' |
| 'Dots에 연결하기' | '로컬 설치 없이 Remote HTTPS MCP로 연결합니다.' | ① '사용 중인 ChatGPT 계정이나 workspace에서 custom MCP를 허용하는지 확인하세요.' ② 'custom MCP 서버 주소에 아래 URL을 등록하세요.' ③ 'Google 로그인 후 표시된 설치 에이전트와 권한을 확인하고 서버 보관에 동의하세요.' ④ '연결한 에이전트와 웹의 내 에이전트가 맞는지 확인하세요.' |

Dots의 펼친 영역에는 선택 가능한 `https://link.knowslog.com/mcp`를 보인다. 보조 문구는 'Dots의 연결 키는 서버에 보관됩니다. Grok의 로컬 키나 연결 폴더를 복사하지 마세요.'다. 복사 버튼은 추가하지 않고 기본 텍스트 선택을 사용한다.

Grok의 펼친 영역에는 다음 최소 실행 안내를 제공한다. README 링크만으로 설치방법을 대신하지 않는다. 이 정보는 신규 설치 전용이며 기존 연결에 적용하지 않는다.

- '준비물: Bot 컴퓨터의 Linux x86_64 또는 arm64, Git, POSIX sh, curl, tar, gzip, sha256sum, Python 3. installer가 Node를 준비합니다.'
- '신규 설치 문서의 소스 준비 절차로 기존 작업과 겹치지 않는 새 폴더를 준비하세요. 검증된 소스 폴더에서 다음 명령을 실행하세요.'

```sh
sh scripts/install_bot_mcp.sh
```

- '검사가 끝나고 Ready가 표시돼야 앱에 등록합니다. 실패하면 다음 단계로 넘어가지 마세요.'
- '앱의 Command MCP 등록에 아래 기본값을 넣으세요. 설치 경로를 바꿨다면 installer가 출력한 실제 경로를 사용하세요.'

| 등록 항목 | 기본값 |
|---|---|
| Name | `knowslink` |
| Type | `Command` |
| Command | `/workspace/.knowslink/node/bin/node` |
| Arguments | `/workspace/.knowslink/knowslink/dist/plugin.js` |

'새 연결을 시작할 때만 다음 비밀이 아닌 설정을 등록하고 MCP를 재시작하세요. 마지막 연결 폴더는 아직 없는 새 경로여야 합니다.'

```text
KNOWSLINK_MODE=public-node
RELAY_URL=https://link.knowslog.com
KNOWSLINK_AGENT_FOLDER=/workspace/.knowslink/connections/grok-new-01
```

'이 설정만으로 연결이 끝나지 않습니다. Grok에 연결을 요청하고 본인이 시작한 브라우저 연결의 지문과 권한을 확인해 동의하세요. 연결 완료를 확인한 뒤 사용하세요.'를 덧붙인다. 키·credential·연결 token을 웹 페이지나 채팅에 붙이는 입력란은 없다. 기존 사용자는 이 신규 폴더·연결 단계를 반복하지 않고 아래 업데이트 링크를 따른다.

Grok의 상세 링크는 'Grok 신규 설치 문서', Dots의 링크는 'Dots 신규 연결 문서'다. 각각 [README Grok](../../../../README.md#grok-신규-설치), [README Dots](../../../../README.md#dots-신규-연결)가 실행 내용 정본이다. 제품 링크는 공개 저장소 `https://github.com/tourbut/KnowsLink/blob/main/README.md`의 같은 앵커를 사용한다. 내부 `.fullops-squad` 문서·worker 로컬 경로를 노출하지 않는다. 위 Bot 설치 경로는 사용자가 실제 입력하는 값이므로 예외다. 긴 소스 준비·고정 버전 선택·검사·실패 복구 절차는 README로 연결한다. DEV는 표시 값과 현재 README/installer의 일치를 확인한다.

두 접기 영역 바로 아래 '이미 연결했다면 새로 등록하지 마세요'를 보인다. 본문은 '기존 에이전트와 키를 유지한 채 업데이트 안내를 따르세요.'다. [Grok 업데이트](../../../../README.md#grok-기존-설치-업데이트)와 [Dots 업데이트](../../../../README.md#dots-기존-연결-업데이트)도 위 공개 README의 해당 앵커로 연결한다. 새 설치와 업데이트를 혼동해 새 계정·키·친구를 만들도록 안내하지 않는다.

**수신·알림 설명**

설치 영역의 마지막에는 접지 않은 소제목 '받는 것과 자동 답장은 다릅니다'와 다음 세 문장을 둔다. 호스트별 조건은 해당 접기 영역에도 짧게 표시한다.

- 공통: '메시지를 받았다고 자동으로 답장하지 않습니다. 새 메시지를 보낼 때마다 대상과 본문에 대한 사용자 승인이 필요합니다.'
- Grok: '자동 수신에는 실행 중인 MCP 또는 watcher가 필요합니다. 질문 없이 채팅에 알림을 띄우려면 별도 webhook routine 설정이 필요합니다.'
- Dots: '질문 없는 알림에는 Events 지원과 구독, 받은 내용을 대화에 표시하는 지시가 필요합니다. 메뉴와 사용 권한은 계정마다 확인하세요.'

마지막 보조 문구는 '연결 완료나 수신 상태만으로 화면 표시·답장 성공을 보장하지 않습니다.'다. 내부 polling 초수·loopback wake·검수 실패 이력은 첫페이지에 나열하지 않는다. 장시간·재시작 후 자동화, 일정 조회·예약·결제 등 외부 업무 실행을 보장하지 않는다.

**대화 예시**

h2는 '이렇게 메시지를 주고받을 수 있어요'다. 바로 아래에 '기능을 설명하기 위한 가상 예시입니다. 실제 회원의 대화나 운영 성공 사례가 아닙니다.'를 항상 보인다. 예시는 정적 문서이며 입력창·보내기 버튼·실제 receipt·시각·온라인 표시는 없다.

| 예시 | 표시할 내용 | 반드시 붙일 구분 |
|---|---|---|
| '웹에서 직접 쓰기' | '나 / 메모 도우미 · 내가 작성' 다음에 '오늘 정리한 질문은 두 가지예요. 첫째, 소개 문구는 어떤가요? 둘째, 시작 순서가 이해되나요?'를 둔다. 수신 대상은 '친구 민지 / 검토 도우미'다. | '웹에서 선택한 내 에이전트 명의로 직접 작성한 메시지입니다.' |
| '연결한 에이전트에 보내기 요청' | 사용자 요청: '민지의 검토 도우미에게 다음 문장을 보내 줘. “소개 문구 초안을 확인해 주세요.”' 다음 줄은 '대상과 본문 확인·사용자 승인 후 발송'이다. | '연결과 친구 수락이 먼저 필요합니다. 수신이나 친구 수락은 답장 승인이 아닙니다.' |

예시 이름은 합성이며 실제 운영의 노우·다닷 이름과 실제 본문을 재사용하지 않는다. 상대 답변을 자동 생성해 보여주지 않는다. 캘린더 접근·첨부·장기 검색·업무 완료 예시는 넣지 않는다.

예시 다음 '대화에서 알아둘 점'에 세 문장을 보인다. '대기는 서버가 메시지를 받은 상태입니다. 전달됨은 수신 측 보관·확인이 끝난 상태이며 읽음이나 답장 완료를 뜻하지 않습니다.' '전달 기한은 최대 180초입니다. 기한이 지나면 대상을 확인하고 새 메시지로 보내세요.' '웹 대화 사본은 수락 후 최대 24시간 보관됩니다. 백업이나 상대 클라이언트의 사본까지 즉시 지워지는 것은 아닙니다.'

**시작하기**

h2는 '내 에이전트로 대화 시작하기'다. 비로그인 본문은 '로그인한 뒤 내 이름과 에이전트를 확인하세요.'다. Google 활성 시 'Google로 계속'과 '같은 Google 계정으로 내 에이전트를 계속 사용합니다.'를 둔다. 이메일은 실제 활성 환경에서만 아래 상태표대로 제공한다. 가입 전 별도 요금제·결제·대기 신청·광고 추적은 없다.

### 로그인 상태와 실패 복구

소개·순서·설치·예시는 아래 모든 상태에서 공개로 읽을 수 있다. 인증 관련 안내·오류는 긴 소개 아래에 묻히지 않도록 헤더 다음에도 요약한다. 로그인 폼 근처에는 전체 오류와 다음 행동을 제공한다. 정보 읽기와 인증 POST는 별개다.

| 상태 | `#start`의 내용·행동 | 진입·복구 계약 |
|---|---|---|
| anonymous + Google-only | 'Google로 계속' | 기존 `POST /auth/google`; 이메일 입력은 없음 |
| anonymous + Google+email | Google 주 버튼, '또는', '이메일로 가입·로그인', 라벨 '이메일 주소', '확인 코드 받기' | 기존 `POST /auth/start`; 6자리·10분·한 번 사용 안내를 유지 |
| anonymous + email-only | 이메일 폼을 주 수단으로 표시 | 존재하지 않는 Google 버튼을 보이지 않음 |
| 로그인 수단 none | '로그인을 준비하고 있습니다. 잠시 뒤 다시 시도하세요.'와 '안내 다시 불러오기' | 버튼처럼 보이는 가짜 로그인 없음. 헤더/hero CTA는 '로그인 상태 확인'으로 `#start` 이동 |
| 유효한 로그인 | '로그인되어 있습니다.'와 '대화로 이동' | `/` 공개 안내 유지, `/home` 진입. 이름 미설정은 기존 이름 설정으로 이동 |
| expired 또는 로그아웃 | 기존 안전 안내와 활성 로그인 폼 | '세션이 끝났습니다. 다시 로그인하세요.' 등 원 의미 유지. 재로그인 뒤 자동 발송하지 않음 |
| 이메일 형식·발송 실패·요청 제한 | 필드에 연결된 오류, 실제 재시도 가능 시각, 다시 시도할 폼 | 오류를 성공으로 표시하지 않음. 현재 응답의 안전한 입력만 유지. 회원 존재 노출 없음 |
| Google 실패·취소 | 기존 안전 오류와 '로그인 다시 시도' 경로 | 기존 오류 응답과 인증 허용 범위를 유지. 원 redirect/연결 복귀를 일반 홈으로 덮지 않음 |

폼은 `#start` 한 곳에만 둔다. 헤더·hero는 앵커 또는 `/home` 링크다. 앵커만으로 인증을 시작하지 않는다. 오류 응답은 문서 처음에 '로그인 문제를 확인하세요' 링크를 보여 폼으로 바로 이동하게 한다. error focus·fragment 처리 방법은 DEV가 기존 서버 렌더 방식에서 정한다. 브라우저 뒤로가기나 재시도 때 폼·오류를 찾을 수 있어야 한다.

이메일 코드 확인은 기존 별도 화면을 유지한다. Google 완료·Remote 동의·클라이언트 연결·재인증 화면은 이번 공개 페이지로 통합하지 않는다. 유효 세션의 존재만 노출하고 이메일·회원 이름·agent 목록·친구·본문·키는 공개 안내에 렌더하지 않는다. 사용자별 상태 응답의 캐시·CSP·CSRF·redirect 검증은 기존 보안 기준을 유지한다.

### 시각 계획·토큰·반응형 와이어프레임

기억할 요소는 사람 아래에 에이전트를 들여쓴 관계 도식이다. 회원 앱의 이름 목록과 같은 문법을 사용한다. 소개는 두 열, 시작 순서는 세로 목록, 설치는 두 접기 행, 예시는 두 개의 읽기 블록으로 구성한다. 반복 카드·장식 통계·배지·큰 그라디언트·애니메이션·설치 명령으로 채운 첫 viewport는 사용하지 않는다.

| 요소 | 공개 화면 값 | 기존 정본과 관계 |
|---|---|---|
| 색 | canvas `#FFFFFF`, ink `#182B31`, muted `#52646B`, accent `#005A62`, selected `#E8F3F3`, edge `#75868B` | 기존 여섯 토큰 그대로. 오류는 기존 danger 사용 |
| 글꼴 | `"Malgun Gothic", "Apple SD Gothic Neo", "Segoe UI", sans-serif` | 기존 한글 가독성·외부 요청 없음 유지 |
| 글자 | hero 40px/1.3/700, mobile hero 30px/1.3/700, h2 24px/1.3/700, h3 18px/600, 본문 16px/1.6, 보조 14px/1.6 | 공개 hero/h2만 확대. 회원 화면 제목 크기는 유지 |
| 폭·정렬 | 공개 콘텐츠 최대 1120px, 본문 최대 64ch, 모두 왼쪽 정렬 | 회원 앱의 1376px 열 구조와 별도 배치 |
| 간격 | 기존 4/8/12/16/24/32px; 공개 섹션 사이 64px, mobile 48px | 공개 큰 구획은 기존 32/24 간격을 두 번 사용 |
| 형태·접근 | 최소 동작 높이 44px, 4px 모서리, 3px accent focus와 2px offset | 기존 버튼·입력·summary 재사용, 새 테마 없음 |

1440px에서는 콘텐츠 양쪽 여백 160px, hero는 3:2 비율과 32px 간격이다. 페이지 헤더는 문서 흐름에 있고 고정하지 않는다. hero 높이는 콘텐츠에 맞추며 고정 100vh를 쓰지 않는다. 설치 명령과 예시 전체를 hero 안에 넣지 않는다.

```text
1440 × 900
KnowsLink           시작 순서   설치와 연결   대화 예시   로그인

친구와 연결하고,                   사람끼리 친구, 에이전트끼리 대화
에이전트로 대화하세요              나                   친구
설명 두 문장                        내 에이전트          친구의 에이전트
웹 설치 불필요 안내                   서로 수락한 친구
[로그인하고 시작하기]  시작 순서 보기

이 순서로 시작하세요
1 로그인하고 내 이름 정하기
2 내 에이전트 연결하기           설치와 연결 방법
3 친구 요청하고 수락하기
4 대상을 고르고 메시지 보내기

이하: 설치와 연결 / 가상 대화 예시 / 실제 로그인 폼 / 문서 링크
```

768px 미만에서는 한 열과 좌우 16px 여백을 사용한다. 768px 이상에서도 두 열이 이름과 본문을 좁히면 한 열로 감싼다. 390·320px에서 헤더 링크는 자연스럽게 두 줄 이상으로 감싼다. 햄버거 메뉴와 가로 스크롤 탭은 만들지 않는다. 순서는 제목·설명·CTA·관계 도식이며 첫 viewport에 CTA가 들어가야 한다. 글자 200%에서는 같은 높이 목표를 강제하지 않고 전체 내용에 접근하게 한다.

```text
390 × 844 / 320px
KnowsLink                    로그인
시작 순서   설치와 연결   대화 예시

친구와 연결하고,
에이전트로 대화하세요
서비스 설명
웹 설치 불필요·연결 필요 안내
[     로그인하고 시작하기      ]
시작 순서 보기
사람끼리 친구, 에이전트끼리 대화
나                   친구
  내 에이전트          친구의 에이전트
서로 수락한 친구

이하 같은 읽기 순서, 설치 details는 한 열
```

도식은 텍스트 읽기 순서로 같은 관계를 설명한다. 선·위치만으로 의미를 전달하지 않는다. 좁은 폭에서는 도식도 세로로 쌓는다. URL은 줄바꿈하며 페이지 전체 가로 넘침을 만들지 않는다. 영문 URL 줄바꿈이 실제 선택 문자열을 바꾸면 안 된다.

### 자기 비평과 수락 기준

초안은 hero 옆에 가짜 채팅과 설치 버튼 두 개를 놓았다. 채팅은 실제 성공 사례로 오인될 수 있고 설치가 로그인보다 먼저 필요해 보였다. 수정안은 hero에 관계 도식만 두고, 설치를 별도 구획으로 내렸다. 예시에는 가상 표시와 작성 주체를 붙였다.

두 번째 비평은 긴 안내 하단의 로그인과 오류 발견성이었다. header·hero 앵커로 폼에 바로 도달하게 하고 오류 요약을 상단에 둔다. 인증 폼은 한 곳만 유지한다. 첫 방문자는 읽을 수 있고 반복 방문자는 대화로 바로 들어갈 수 있다.

coor의 설치 완결성 점검(msg_6593a1e11cb6)을 받고 링크만으로는 신규 사용자가 앱에 무엇을 입력할지 알기 어렵다고 판단했다. Grok 접기 안에 준비물·조건부 installer 명령·Command/Arguments·신규 연결 설정을 추가했다. 기존 키/agent를 유지하는 업데이트 경로는 분리했다. 전체 README를 복제하지 않고 오류·버전별 상세만 링크로 남긴다. 이 비평은 자기 검토이며 실제 렌더 수락이 아니다.

LG-01–08 제품 기준은 [D02](../../planning/product-specs/SAR-HUMAN-PEER.md)의 공개 안내 절을 따른다. DEV는 위 카피·상태·반응형을 구현하고 핵심 auth/template 회귀와 필수 make lint/test를 수행한다. tester는 같은 고정 SHA에서 anonymous·Google-only·both·email-only·none·error·expired·authed, 실제 POST·redirect·정보 비노출을 독립 검증한다. designer는 1440·390·320의 hero·설치 열림·예시·오류 실제 캡처를 확인한다.

키보드는 Tab/Shift+Tab, 앵커 Enter, summary Enter/Space, 폼 Enter, 오류 뒤 재시도를 확인한다. 200% 확대·44px hit area·focus·contrast·긴 URL·overflow를 검증한다. 장식 모션은 없어 reduced-motion에서 동작 변화가 없어야 한다. 실제 Google·운영 메일·실상대 메시지는 실행하지 않는다. 독립 QA·직접 시각 검수·별도 세션 고정 SHA 리뷰 전 구현 수락으로 보고하지 않는다.

이번 실행은 codex / gpt-6-astra / high, 주 스킬 frontend-design이다. 최신 사용자 지정이 기존 Sol tuple보다 우선한다. frontend check는 configuration_ready이며 인증·브라우저 QA·MCP는 not_run이다. tools=[]와 기존 설정을 유지한다. [설계 검증 기록](../../evaluations/qa-reports/SAR-LANDING-GUIDE-001-test/report.md)에 원천 사실·미실행·후속을 구분한다.

구현 연결(2026-10-14): 공개 template·인증 view·내 계정 복귀 링크와 관련 회귀는 [DEV 구현 기록](../../exec-plans/phases/SAR-LANDING-GUIDE-001-DEV.md)에 연결한다. 위 설계 검사 기록과 과거 실패는 보존한다. DEV 자체 검증은 독립 QA·designer 직접 시각 수락·코드 리뷰·운영 배포를 대신하지 않는다.

## 현재 디자인 정본 — SAR-HUMAN-PEER-001

2026-10-12 사용자 요청에 따라 관리 중심 화면을 사람·이름·대화 중심으로 교체한다. 이 절부터 '이전 설계와 검수 기록' 앞까지가 현행 D04다. 아래 이전 1–14절은 과거 화면과 수락 기록이며 새 구현 지시가 아니다. 현재 제품 규칙은 [사람 중심 D02](../../planning/product-specs/SAR-HUMAN-PEER.md)의 HP-01–16이다. 토큰·레이아웃·문구가 충돌하면 현행 절을 적용한다. 이 문서는 설계 완료이며 구현·직접 시각 수락은 대기다.

### 관찰과 전면 개편 목표

기준 f20b33a의 `member_pages.go`는 에이전트·관계·승인·명령·계정이 긴 홈에 함께 있다. `member_agents.go`는 ID 위주 목록이며 이름 필드가 없다. `public_text.go`는 signed text와 최대 180s TTL을 사용한다. 웹 composer와 사람 친구 모델은 없다. `test_messages.go`는 시험 전용 제한 경로이며 일반 회원 흐름을 대신하지 못한다.

이번에 기존 after PNG `home-1440-s1of7.png`와 `home-390-s1of10.png`를 직접 열었다. 출처는 [polish DEV after](../../evaluations/qa-reports/SAR-FRONTEND-POLISH-001-DEV-test/after/manifest.json)의 f7b4786 합성 렌더다. 실제 계정/현재 운영 브라우저의 관찰은 아니다. 두 화면은 에이전트 ID·키·연결 설명이 첫 화면을 차지한다. 친구 이름과 대화는 없다. 기존 polish의 좁은 PASS는 이번 구조의 수락 근거가 아니다.

제품의 주된 화면은 '누구의 어느 에이전트와 이야기하는가'를 보여주는 대화 공간이다. 기억할 시각 요소는 사람 이름 아래 에이전트 이름이 들여쓰기되어 이어지는 목록이다. 카드·통계·기술 설명이 첫인상을 만들지 않게 한다.

### 디자인 계획 1과 비평 1

초안은 좌측 184px 이동 열, 중앙 친구별 agent 목록, 우측 대화였다. 흰 바탕과 청록색 선택선, 한국어 시스템 글꼴을 정했다. 내 agent와 친구 agent를 항상 두 목록으로 놓고 대화 위에 선택기 두 개를 배치했다.

비평: 세 열을 모두 선택기로 쓰면 관리 콘솔과 비슷하다. 두 선택기가 열린 상태로 남으면 메시지 대상을 바꾸기 쉽다. 390px에서는 폼과 목록이 본문보다 먼저 나온다. 수정: 가운데 열은 최근 대화 목록으로 바꾼다. '새 대화'에서만 친구/내 발신/상대 수신을 선택한다. 선택 후 대화 헤더에 대상을 고정한다. 친구 화면에서 사람을 열 때만 그 아래 agent를 보인다.

### 디자인 계획 2와 비평 2

수정안은 '대화 / 친구 / 내 에이전트' 세 기능으로 나눈다. 계정은 헤더의 '내 계정'으로 이동한다. 승인 요청은 대화 목록 상단의 조건부 '승인 N건' 링크로 발견하고 별도 상세에서 결정한다. 본문은 말풍선 대신 작성자·시각·본문의 조용한 줄 단위로 배치한다. 내 메시지에는 좁은 청록 선택선과 '내가 작성' 표시를 쓴다.

비평: 메신저의 말풍선을 제거하면 양쪽 agent의 이름이 같은 경우 발신 구분이 약해진다. 보존 설명을 숨기면 대화가 계속 남는다고 오인할 수 있다. 수정: 매 메시지에 사람 이름과 agent 이름을 함께 표시한다. 내가 직접 작성했으면 '내가 작성'을 추가한다. 본문 위에는 '메시지는 24시간 안에 사라집니다' 한 문장만 둔다. 영구 삭제·백업 한계는 '보존 안내'에 넣는다. 빈 홈은 대화 제목과 바로 쓸 수 있는 다음 동작 하나로 끝낸다. 불필요한 예시 대화나 가짜 받은 메시지는 넣지 않는다.

두 비평은 designer의 자기 비평이다. 별도 독립 문서 리뷰와 구현 화면 검수는 coor가 후속으로 배정한다.

### 현행 토큰과 형태

| 토큰 | 값 | 역할 |
|---|---|---|
| `--canvas` | `#FFFFFF` | 전체 바탕·입력 |
| `--ink` | `#182B31` | 본문·이름 |
| `--muted` | `#52646B` | 시각·보조 설명 |
| `--accent` | `#005A62` | 선택선·주 버튼·포커스 |
| `--selected` | `#E8F3F3` | 선택된 대화·내 작성 표식 |
| `--edge` | `#75868B` | 입력 테두리·기능상 구분선 |

위 여섯 값이 기본 팔레트다. 오류·위험 동작은 기존 `--danger:#A01616`을 재사용한다. 성공·대기·실패는 색을 추가하지 않고 텍스트로 표시한다. 토큰 대비 계산 7개 조합은 통과했다. 본문/보조/주 버튼의 흰 바탕 대비는 14.70/6.19/7.96:1, 입력 테두리는 3.79:1이다. [계산 근거](../../evaluations/jev/SAR-HUMAN-PEER-001-design-check.json)는 명세 토큰만 검사하며 실제 구현 픽셀의 적합성은 후속 확인한다. 본문 4.5:1, 입력/포커스 3:1 이상을 요구한다.

- 한국어 UI는 `"Malgun Gothic", "Apple SD Gothic Neo", "Segoe UI", sans-serif`를 사용한다. 첫 가족은 Windows의 한글 획과 이름 목록 가독성을 기준으로 골랐다. 플랫폼에 없으면 다음 가족으로 간다. 외부 폰트·이미지는 없다.
- 제품명 20px/700, 페이지 제목 28px/700, 친구 이름 18px/600, agent/본문 16px/400–600, 시각·보조 14px/400이다. 본문 줄높이 1.6, 제목 1.3이다. 모든 내용은 왼쪽 정렬이다. 본문은 최대 64ch다.
- 간격은 4/8/12/16/24/32px다. desktop 바깥 여백 32px, mobile 16px다. 버튼·입력은 최소 44px이고 모서리 4px다. 카드 그림자·그라디언트·장식 배지·큰 통계·알약 버튼은 없다.
- 구조 구분은 열 경계와 목록 선택선으로만 한다. 각 메시지·각 agent를 테두리 카드로 감싸지 않는다. 들여쓰기 16px는 사람 아래 agent라는 의미를 갖는다.
- 긴 이름은 최대 두 줄로 표시한다. 줄임표만으로 대상을 구분하지 않는다. 전체 이름은 대화 헤더·선택 확인·상세에서 확인한다. 식별자만 `overflow-wrap:anywhere`와 모노스페이스를 쓴다. 한글 본문은 `keep-all`과 긴 단어 fallback을 쓴다.
- 포커스는 3px `--accent`, offset 2px다. hover만으로 동작을 찾게 하지 않는다. 장식 모션은 없고 reduced-motion에도 즉시 동작한다. 모바일 가상 키보드가 작성/보내기/오류를 가리지 않아야 한다.

### 정보 구조와 화면

기능 이름은 고정한다. URL·handler·저장 구조는 DEV가 정하고 기술 문서에 연결한다. `/home`은 대화 진입점으로 바꾼다. 기존 연결 승인·gate·receipt 직접 링크는 해당 권한을 유지하며 적절한 상세로 복귀한다.

| 화면 | 기본 정보와 행동 | 감추거나 옮길 내용 |
|---|---|---|
| 로그인 | KnowsLink, '친구의 에이전트와 대화하세요', 'Google로 계속' | 세 줄 서비스 소개·내부 진행 상태 삭제. 이메일 로그인은 실제 활성 설정일 때만 표시 |
| 첫 이름 설정 | '내 이름', 저장 후 대화로 이동 | 이메일 입력을 다시 요구하지 않음. 이름은 친구에게 보인다는 한 줄 안내 |
| 대화 홈 | 최근 대화, 새 대화, 조건부 승인 링크, 선택한 대화 | agent ID·키·설치 명령·계정 만료·상시 보안 배너 없음 |
| 새 대화 | 친구 또는 내 에이전트, 내 발신 agent, 상대 수신 agent, '대화 열기' | 선택 대상 없는 비어 있는 폼 없음. 친구가 없으면 친구 추가, 내 agent가 없으면 agent 추가 |
| 친구 | 받은 요청, 보낸 요청, 친구 이름 목록, 친구 추가 | 일반 디렉터리·이메일 검색 없음. 받은/보낸 요청은 해당 동작만 표시 |
| 친구 상세 | 사람 이름, 아래 active agent 이름·클라이언트·상태, '대화하기' | 키·지문·credential·상대 친구 목록 없음. 해제·차단은 '친구 관리' details 안 |
| 내 에이전트 | 이름·클라이언트·등록됨/연결 필요, '에이전트 추가' | 키 수·ID는 상세. 이름 변경·연결·키 관리·철회를 상세에서 분리 |
| 연결 승인 | 대상 이름, 지문 대조, 실제 허용 범위, 승인·취소 | 실행 코드와 연결 수단은 필요한 연결 화면에만 표시 |
| 승인 상세 | 검증된 원요청과 실제 정책·기한, 승인/거절 | 대화 text나 hint를 승인 증명으로 사용 금지 |
| 내 계정 | 자기 마스킹 이메일·이름·친구 코드·재확인·로그아웃 | 홈 상시 노출 없음. 계정 비활성화는 실제 기능 없으면 표시하지 않음 |
| 기존 연결 | 이전 두 agent의 이름/구분값, 기존 범위 안내, 친구 전환 | 친구로 이미 연결됐다는 표시 금지. 소유 확인 불가면 지원 안내 |

### Desktop 와이어프레임 — 1440×900

최대 폭 1376px, 좌측 이동 184px, 대화 목록 280px, 나머지 대화 본문이다. 열 간격 24px다. 첫 화면에서 이름·최근 대화·본문·작성 위치를 볼 수 있다. 본문이 길면 문서 흐름으로 스크롤하고 작성란은 본문 다음에 둔다. 높이를 고정해 키보드 포커스를 가두지 않는다.

```text
KnowsLink                                                   내 계정
──────────────────────────────────────────────────────────────────
대화               대화            새 대화   지민 / 일정 도우미
친구 (받은 요청 1)  승인 1건                  내 발신: 노우
내 에이전트                                   메시지는 24시간 안에 사라집니다
                   지민
                     노우 / 일정 도우미       지민 · 일정 도우미  오후 2:10
                   내 에이전트                [수신 text]
                     노우 / 다닷              나 · 노우 · 내가 작성  오후 2:12
                                             [송신 text]  대기

                                             메시지
                                             [여러 줄 입력            ]
                                             노우 → 지민의 일정 도우미
                                             [보내기]
```

와이어프레임의 이름·text는 설명용 합성 예시다. 구현의 빈 계정에는 주입하지 않는다. 목록의 최근 표시는 24h 안 참여한 대화만 사용한다. 미열람 수·읽음·접속 추정은 없다.

### Mobile 와이어프레임 — 390×844, 320px

한 화면에 목록 또는 대화 하나만 보인다. 고정 하단 메뉴와 중첩 스크롤은 쓰지 않는다. 목록은 문서 흐름이며 대화로 이동하면 위에 '대화 목록' 복귀 링크가 있다. 768px 미만은 이 구조다. 768–1199px은 상단 이동과 목록/본문 두 열, 1200px 이상은 desktop 세 열이다.

```text
KnowsLink                         내 계정
대화       친구       내 에이전트
대화                            새 대화
지민
  노우 / 일정 도우미                 대기
내 에이전트
  노우 / 다닷

별도 대화 화면:
대화 목록
지민 / 일정 도우미
내 발신: 노우
메시지는 24시간 안에 사라집니다
지민 · 일정 도우미          오후 2:10
[수신 text]
나 · 노우 · 내가 작성        오후 2:12
[송신 text]                     대기
메시지
[여러 줄 입력                       ]
노우 → 지민의 일정 도우미
[              보내기              ]
```

320px에서도 세 이동 링크는 두 줄 이내로 감싸고 본문 가로 스크롤은 없다. 긴 사람/agent 이름은 줄바꿈한다. 키보드/200% 확대에서도 대상·입력·보내기 접근을 유지한다.

### 상태·문구·복구

| 상황 | 보이는 문구/동작 | 완료 관찰 |
|---|---|---|
| 새 계정 | '첫 에이전트를 추가하세요' / 에이전트 추가 | 대화 제목은 보이고 가짜 대화는 없음 |
| agent 있고 친구 없음 | '친구를 추가하거나 내 에이전트끼리 대화하세요' / 새 대화 | 내 agent가 하나면 두 번째 연결로 안내 |
| 받은 요청 | '{이름}의 친구 요청' / 수락, 거절 | 수락 전 '서로의 현재·앞으로 추가할 에이전트가 연결됩니다' 표시 |
| 보낸 요청 | '요청 대기' / 요청 취소 | 상대 agent·이메일·존재 검증 결과 없음 |
| 요청 제출 | '요청을 보냈습니다. 가능한 요청은 상대가 확인할 수 있습니다' | 없는/차단 코드에도 같은 안내. 친구 성공으로 표시하지 않음 |
| 친구 agent 0개 | '연결된 에이전트가 없습니다' / 친구 목록 | 빈 select/보내기 없음 |
| 보내는 중/응답 유실 | '전송 결과를 확인하세요' / 상태 확인 | 같은 멱등키로 확인/재시도, 무조건 새 발송 금지 |
| 대기 | '대기' | 수신 ACK 전 전달됨/읽음 표시 없음 |
| 전달 완료 | '전달됨' | 승인·처리 완료와 구분 |
| TTL 만료 | '전달 기한이 지났습니다' / 새 메시지 | 자동 재발송 없음. 같은 시도 TTL 연장 없음 |
| 해제/차단/agent 철회 | '지금은 메시지를 보낼 수 없습니다' / 대화 목록 | 과거 복사본은 읽기 전용. 상대에게 차단 사유 노출 없음 |
| 24h 경과 | '보존 기간이 지난 대화입니다' / 대화 목록 | 서버 원문/미리보기 비노출 |
| 재확인 필요 | '로그인을 다시 확인하세요' / 로그인 다시 확인 | 민감 본문은 URL에 넣지 않음. 재로그인 뒤 자동 발송 없음 |
| 한도/일시 장애 | 실제 한도/실패와 재시도 가능 시각 / 돌아가기 | 실패 입력은 비공개 상태에서 복구 가능. 성공 토스트 금지 |
| 이름 저장 오류 | '이름은 1–40자로 입력하세요' | 필드 라벨과 오류 연결, 입력한 안전한 값 유지 |

친구 해제·차단, agent/key 철회는 관리 details 안에서 결과 문장을 먼저 보여준다. 기존 '별도 모달 없이 결과 안내 후 단일 POST' 방향을 유지한다. 친구 해제 문구는 '새 대화가 막힙니다. 남은 대화는 보존 기간까지 볼 수 있습니다'다. agent 철회 문구는 '이 에이전트의 모든 키와 새 전달 권한이 끝납니다'다. 철회 agent의 최소 24h 기록 보존과 정확한 정리 시각 미보장 안내는 상세에 유지한다.

### 접근성·구현 경계

`lang=ko`, 고유 title, h1 하나, header/nav/main, 본문 건너뛰기와 보이는 라벨을 사용한다. 선택한 기능/대화는 `aria-current`와 텍스트로 식별한다. 오류는 alert, 저장/발송 결과는 status로 알린다. 메시지 영역은 읽기 순서대로 작성자·시각·본문·상태를 제공한다. 자동 전체 화면 낭독이나 자동 포커스 이동을 하지 않는다. 이름이 같은 버튼은 대상 이름과 구분값으로 접근성 이름을 구분한다.

Go template/CSS·기존 네이티브 폼·details로 구현한다. CSP와 인증/서명/CSRF/owner 검사를 완화하지 않는다. JS·외부 폰트·이미지·새 프레임워크·유료 에셋은 추가하지 않는다. 선택·보내기·새로고침은 서버 왕복이어도 된다. 자동 갱신/가짜 온라인은 필수가 아니다. '새로고침'으로 수신 상태와 본문을 확인하는 실제 경로를 제공한다. 기존 SSE/polling/remote 수신 경로는 친구 권한을 동일하게 집행한다. 웹 발송의 실제 주체와 서명 설계는 DEV가 결정한다.

### DEV/QA/직접 시각 수락

P0 전체 흐름: 로그인·이름, 내 agent 등록/이름 변경, 친구 요청/수락, 친구 agent 목록, 발신/수신 선택, 실제 text 송수신, 해제 뒤 새 전달 차단, legacy 안전 이행이다. UI의 간격/색만 바꾸고 완료로 표시하지 않는다.

| 담당/시점 | 범위 | 통과 조건 |
|---|---|---|
| DEV / 구현 고정 SHA 전 | 핵심 실제 handler/API/DB와 실패 경계, 필수 제품 lint/test | HP-01–16 추적, 서명/인가/세대/멱등/보존/이행 검증, D03/D05–D10 갱신 |
| tester / DEV SHA 고정 뒤 | 격리 A/B/C 회원·legacy fixture, 실제 조작 | 두 회원의 요청/수락/양방향 대화와 제삼자 거부, 미수락·차단·해제·철회·재수락·24h 경계 |
| tester / 같은 SHA | 320/390/1440와 키보드 | overflow 없음, 44px 동작, Tab/Enter/Space/뒤로/취소, 대화 대상 보존과 입력 복구 |
| designer / 같은 SHA | 실제 구현의 로그인·홈·친구·대화·내 agent·해제·gate 오류 캡처 | 관리 카드 대신 이름/대화가 첫 화면을 차지함, 두 이름 구분, 문구 최소화, 보안 안내가 실제 동작 옆에 있음 |
| 별도 reviewer / 고정 SHA | 코드·제품 문서 독립 리뷰 | critical/high 0, 현행 정본과 구현 일치 |

브라우저 QA는 기존 Orca 엔진을 재사용한다. tester는 자기 loopback origin과 포트, fixture/DB 식별자, 허용 GET/명시 POST, 독립 assertions와 요청/시간 예산을 기록한다. 운영 origin·실메일·실상대 발송은 허용하지 않는다. 탐색/shadow/CUAWright는 OFF다. 정적 fixture 캡처와 실제 인증·DB 동작을 구분한다. 실제 Google 왕복·외부 앱·운영은 별도 기록한다.

과거 V-02는 새 요청 상태별 동작 분리로 대체한다. V-03 날짜는 짧은 KST '오늘 오후 2:10', 이전 날짜는 '10월 11일 오후 2:10', 연도가 다르면 연도를 붙인다. 정확한 기한은 상세에서 제공한다. V-08 gate/receipt도 같은 한국어 상태를 적용하되 원 의미를 유지한다. V-12 이동 순서와 N-02 mobile 주 버튼 폭은 위 구조를 적용한다. 명령은 연결 상세로 이동하며 복사 값과 줄바꿈을 확인한다. 옛 결함을 새 구현에서 해결됐다고 미리 판정하지 않는다.

실행 선택은 codex / gpt-6.1-sol / high, 주 스킬 frontend-design이다. `frontend.py --check`는 configuration_ready, exit 0이었다. frontend.json의 이미 활성화된 tuple/정본을 재사용하며 설정을 다시 쓰지 않는다. 인증·MCP·새 브라우저 QA는 not_run이다. 현재 작업은 제품 코드가 없는 설계이며 구현 캡처·동작 수락은 후속이다.

### 2026-10-12 구현 후보 직접 시각 검수 — SAR-HUMAN-PEER-001-VISUAL

최종 제품 후보 `254ff6a98b12b7e51f0835491300f7695425d264`의 시각 수락은 FAIL이다. V-01 작성란 기술 설명/반복 대상 문장과 V-02 긴 목록 이름/desktop 작성 위치가 medium 수정 대상이다. V-03 홈의 반복 제목은 low다. 관찰 범위의 critical/high는 0건이며 미관찰 상태의 결함 부재를 뜻하지 않는다. 기존 제품/디자인 요구는 완화하지 않았다.

먼저 실제 합성 auth/API/DB의 `4eaefa6e0a5cf623f0d8d8080d918086f86030ec` 원래 PNG 21장을 모두 직접 열었다. 최종 후보와 memberStyle 및 관련 대화 template의 동일성을 확인했다. 과거 이미지를 최종 SHA 실행으로 표시하지 않는다. CSS1440/DPR0.75·320/390 DPR1 조건, 실제 surface·200%·focus·미포함 상태의 검사 범위와 not_run, 항목별 수정 수락 조건은 [이번 시각 보고](../../evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/report.md)에 있다. DEV 수정 후 새 고정 SHA의 영향 화면 검수와 tester/별도 리뷰가 필요하다.

### 2026-10-12 수정 후보 직접 시각 재검수 — 12dc0c7

원후보254ff6a FAIL은 위 절과 원보고에 보존한다. 후보 `12dc0c733d4d0bcdf0094fdebda7309da4c3fdd5`의 제공된 V01–V03 수정 화면은 시각 PASS다. 실제14PNG를 모두 직접 열고 sourcefe7a509와 최종 제품 동일성을 확인했다. 별도 isolated page의15상태에서 두 줄 목록·가로 overflow0·44px 제어·desktop 첫 화면 textarea/보내기를 확인했다. 실제 키보드/200%·로그인/이름/pending/해제/gate의 직접 시각 검사는 not_run이며 전체 제품 수락을 뜻하지 않는다. 범위·보정·후속 정본: [수정 후보 직접 재검수](../../evaluations/qa-reports/SAR-HUMAN-PEER-001-VISUAL/final-report.md).

### 서버 캘린더 연결 제거 — SAR-PUBLISH-001 (D018)

2026-10-21 사용자 결정 D018로 서버 캘린더 연결을 제거한다. 제품 규칙의 정본은 [SAR-SCHEDULE-001 '철회 — D018'](../../planning/product-specs/SAR-SCHEDULE-001.md#철회--d018-2026-10-21)이다. 이전 배치표(SAR-SCHEDULE-001-VISUAL)는 비현행이며 [시각 검수 보고](../../evaluations/qa-reports/SAR-SCHEDULE-001-VISUAL/report.md)에 기록으로 남는다. 새 스타일·토큰·JS는 없다.

| 화면·상태 | 제거 뒤 기대 화면 |
|---|---|
| 설정 | '캘린더 연결' 절 전체가 없다. 연결 전 안내, 연결·다시 연결 버튼, 이메일 회원 안내, 준비 중 안내, 연결 해제 details가 모두 없다. 나머지 절의 순서와 간격은 바뀌지 않는다. |
| 해제 뒤 status | 없다. 상단 status에 캘린더 문구가 나오지 않는다. |
| 승인 상세 `schedule.query` | U5 기본 상세로 돌아간다. 요약은 누가 누구에게, 요청 유형, 검증된 구조화 필드(기간·단위가 wire에 있으면 표시), 기한·상태다. '공개되는 내용' 행이 없다. 고정 설명은 U5 원문 '승인은 이 요청의 인간 게이트만 통과시킵니다. 정보 공개나 일정 실행을 허용하지 않습니다.'다. '이 요청 승인'·'거절'을 모든 회원에게 보인다. |
| 승인 상세 `schedule.commit` | 바뀌지 않는다(U5). |
| 처리됨·만료 | 기존 U5 note와 배지를 쓴다. |
| 기존 연결 회원 | 화면 안내를 두지 않는다. 대상이 테스트 사용자뿐이라 coor가 직접 알린다. |

KO/EN 모두 '캘린더'·'calendar'·'바쁜 시간'·'busy' 문구가 회원 화면에 남지 않는다. 첫페이지의 기존 '일정 조회·예약·결제 등 외부 업무 실행을 보장하지 않습니다' 문장은 유지한다.

### 개인정보처리방침·서비스 약관 페이지 — SAR-PUBLISH-001 (D017)

문안과 완료 조건의 정본은 [SAR-PUBLISH-001](../../planning/product-specs/SAR-PUBLISH-001.md)이다. 이 절은 배치만 정한다.

```text
Desktop 1440 / Mobile 390·320 (같은 한 열, 읽기 폭 약 72ch)
┌───────────────────────────────────────────────┐
│ KnowsLink                          한국어 | EN │  ← 공개 첫페이지 헤더(앵커 메뉴 없음)
├───────────────────────────────────────────────┤
│ 개인정보처리방침                    (h1)       │
│ 시행일: 2026-MM-DD                  (muted)    │
│ 서비스 약관도 함께 확인하세요.      (링크)     │
│                                               │
│ 1. 이 방침의 범위                   (h2)       │
│ 본문 문단·목록·표                              │
│ …                                             │
│ 12. 연락처  knowslog01@gmail.com (mailto)      │
├───────────────────────────────────────────────┤
│ 개인정보처리방침 · 서비스 약관 · 홈 · 설치·사용 문서 · 맨 위로 │ ← 현재 페이지 aria-current
└───────────────────────────────────────────────┘
```

- 스타일은 첫페이지의 `memberStyle`·`body.landing` 토큰을 재사용한다. 표는 기존 표 스타일을 쓰고 320px에서는 표 영역만 가로 스크롤한다. 페이지 전체 가로 넘침은 없다.
- 정책 페이지 하단에만 '홈'(`/`)을 둔다. 첫페이지 하단에는 '홈'이 없다(2026-10-22 SAR-PUBLISH-001-D04, 시각 검수 D-01).
- 첫페이지 `#start` 로그인 버튼 아래 `muted` 한 줄로 두 페이지를 링크한다. 회원 설정은 '서비스 안내' 옆에 두 링크를 둔다.
- 링크 모양(시각 검수 V-01~V-03): 문장 속 링크는 줄 높이를 키우지 않고 inline으로 둔다. 하단·설정처럼 단독 줄 링크만 44px 터치 대상을 유지한다. 하단의 현재 페이지 링크(`aria-current`)는 굵게(600) 표시하고 새 색·글꼴은 쓰지 않는다. 설정의 세 링크 사이는 약 16px 가로 간격을 둔다.
- 피할 것: 크림 배경, 알약 버튼, 모노스페이스 라벨, 장식 그라디언트, 이모지 아이콘, 접는 상자, 영어판의 한국어 잔존.

### Claude Code·Codex 자동 설정 안내 — SAR-AUTOSETUP-001

현행 화면이다. [intent SAR-AUTOSETUP-001](../../planning/intents/SAR-AUTOSETUP-001.md)의 자동 설정이 main `8ec61ab3`에 구현됐고, 2026-10-23 직접 시각 검수([보고](../../evaluations/qa-reports/SAR-AUTOSETUP-001-ART-test/report.md))에서 아래 U8 선택지 표·안내 상자 카피와 일치함을 확인했다(AC-10 PASS). 미해결은 아래 '시각 검수 결과'의 V-01이다. 제품 규칙은 [D02 SAR-AUTOSETUP-001](../../planning/product-specs/SAR-AUTOSETUP-001.md)이다. 바꾸는 곳은 '내 에이전트 → 어디에서 사용할까요?'의 선택지와 Claude Code·Codex 안내 상자뿐이다. 다른 선택지·연결 상태 절·'기존 연결 계속' 안내 문장은 그대로다.

#### 선택지

| 순서 | 선택지 KO / EN | 보조 문구 KO / EN | 바뀐 점 |
|---|---|---|---|
| 1 | 웹에서 바로 사용 / Use directly on the web | 현행 | 없음 |
| 2 | Grok Bot (내 컴퓨터) | 현행 | 없음 |
| 3 | Dots (원격) | 현행 | 없음 |
| 4 | Claude Code 작업 세션 (내 컴퓨터) / Claude Code work session (your computer) | 자동 설정 · 시험 환경 확인 (Claude Code 2.1.302, Windows 11) / Automatic setup · Checked in a test environment (Claude Code 2.1.302, Windows 11) | '연결 안내 있음'을 '자동 설정'으로 바꾼다 |
| 5 | Codex 작업 세션 (내 컴퓨터) / Codex work session (your computer) | 자동 설정 · 시험 환경 확인 (Codex CLI 0.163.0, Windows 11) · 도착 알림 없음 / Automatic setup · Checked in a test environment (Codex CLI 0.163.0, Windows 11) · No arrival notice | 새 선택지. `use=codex` |
| 6 | 지원 확인이 필요한 호스트 / Hosts pending support confirmation (현행) | Hermes Agent, 그 밖의 호스트 / Hermes Agent and other hosts | 보조 문구에서 Codex를 뺀다 |

- 버전·OS 문구는 D02 범위 표의 근거가 바뀔 때만 바꾼다(NG-05). Q4가 B로 결정되기 전에는 macOS·Linux를 쓰지 않는다.
- 6번에는 명령·성공 배지·설치 버튼이 없다(현행 유지).

#### 안내 상자 구성

기존 `dl.kv` 안내 상자를 그대로 쓴다. 행 순서: 준비물, 실행할 명령, 자동으로 하는 일, 직접 확인할 일, 다음 행동, 세션 연결. 명령은 친구 코드와 같은 `textarea.share`(readonly, `spellcheck="false"`)와 `muted` 도움말로 보인다. 복사 버튼·새 JS·새 토큰·새 컴포넌트는 쓰지 않는다. 수동 설치 문서 링크는 맨 아래 `details`로 접는다. 명령 문자열은 DEV가 정한다. 아래 `<명령>`은 그 자리를 뜻한다.

```text
Claude Code 작업 세션 (내 컴퓨터)                (h3)
준비물        Windows 11, 로그인한 Claude Code …
실행할 명령   [<명령>                         ]   (textarea.share)
              입력란에서 전체 선택(Ctrl+A)한 뒤 …  (muted)
자동으로 하는 일  내려받기와 확인, …
직접 확인할 일    설치 계획 동의, …
다음 행동     …
세션 연결     앱에서 확인 필요. …
▸ 수동 설치 문서                                 (details)
```

1440px에서는 기존 읽기 폭 안의 `dl.kv` 두 열이다. 390·320px에서는 기존처럼 한 열이다. 명령 입력란은 폭 100%이고 줄바꿈된다. 페이지 전체 가로 넘침은 없다.

#### Claude Code 카피 (`use=claude`)

| 행 | KO | EN |
|---|---|---|
| 준비물 | Windows 11, 이 컴퓨터에 설치·로그인한 Claude Code, 인터넷 연결. 저장소·Node·설치 경로는 준비하지 않아도 됩니다. 시험 환경은 Claude Code 2.1.302입니다. | Windows 11, Claude Code installed and signed in on this computer, and an internet connection. You don't need the repository, Node, or an install path. Tested with Claude Code 2.1.302. |
| 실행할 명령 | 라벨: 이 컴퓨터의 PowerShell에서 실행할 명령 / 도움말: 입력란에서 전체 선택(Ctrl+A)한 뒤 복사(Ctrl+C)해 PowerShell에 붙여 넣으세요. 설치기는 바꾸기 전에 설치 계획을 보이고 동의를 묻습니다. | Label: Command to run in PowerShell on this computer / Help: Select all in the box (Ctrl+A), copy (Ctrl+C), and paste it into PowerShell. The installer shows its plan and asks you before it changes anything. |
| 자동으로 하는 일 | 내려받은 파일 확인, 새 설치 폴더 만들기, Claude Code에 KnowsLink 등록, Google 연결 시작. 다른 앱의 키나 Claude Code 로그인 정보는 읽거나 복사하지 않습니다. | Checks the download, creates a new install folder, adds KnowsLink to Claude Code, and starts Google connection. It never reads or copies other apps' keys or your Claude Code sign-in. |
| 직접 확인할 일 | 설치 계획에 동의, 브라우저에서 Google 로그인·지문 비교·에이전트 이름·연결 동의, 메시지를 받을 세션에서 `/knowslink-session on`, 보낼 때마다 발송 승인. | Agree to the install plan; in your browser, sign in with Google, compare the fingerprint, name the agent, and agree; run `/knowslink-session on` in the session that should receive; approve every send. |
| 다음 행동 (처음 연결) | 명령을 실행하고 설치기 안내를 따르세요. 끝나면 Claude Code를 다시 열고 메시지를 받을 세션에서 `/knowslink-session on`을 실행하세요. | Run the command and follow the installer. When it finishes, reopen Claude Code and run `/knowslink-session on` in the session that should receive messages. |
| 다음 행동 (기존 연결 계속) | 같은 명령을 실행하면 기존 설치를 갱신합니다. 에이전트·키·친구는 그대로입니다. 갱신 뒤 Claude Code를 다시 열고 `/knowslink-session on`을 다시 실행하세요. | Running the same command updates your existing install. Your agent, key, and friends stay the same. After the update, reopen Claude Code and run `/knowslink-session on` again. |
| 세션 연결 | 현행 문구 유지. | Keep the current text. |
| details 요약 | 수동 설치 문서 | Manual install guide (in Korean) |

#### Codex 카피 (`use=codex`)

| 행 | KO | EN |
|---|---|---|
| 준비물 | Windows 11, 이 컴퓨터에 설치·로그인한 Codex CLI, 인터넷 연결. 저장소·Node·설치 경로는 준비하지 않아도 됩니다. 시험 환경은 Codex CLI 0.163.0입니다. | Windows 11, Codex CLI installed and signed in on this computer, and an internet connection. You don't need the repository, Node, or an install path. Tested with Codex CLI 0.163.0. |
| 실행할 명령 | Claude Code와 같은 라벨·도움말. | Same label and help as Claude Code. |
| 자동으로 하는 일 | 내려받은 파일 확인, 새 설치 폴더 만들기, Codex에 KnowsLink 등록, Google 연결 시작. 다른 앱의 키나 Codex 로그인 정보는 읽거나 복사하지 않습니다. Codex 승인 모드는 바꾸지 않습니다. | Checks the download, creates a new install folder, adds KnowsLink to Codex, and starts Google connection. It never reads or copies other apps' keys or your Codex sign-in. It does not change Codex's approval mode. |
| 직접 확인할 일 | 설치 계획에 동의, 브라우저에서 Google 로그인·지문 비교·에이전트 이름·연결 동의, 메시지를 받을 세션에서 `knowslink_session`을 `on`으로 실행하게 한 뒤 Codex 도구 승인과 4자리 확인 코드 입력, 보낼 때마다 확인 코드로 발송 승인. | Agree to the install plan; in your browser, sign in with Google, compare the fingerprint, name the agent, and agree; in the receiving session, ask Codex to run `knowslink_session` with `on`, approve the Codex tool, and enter the 4-digit code; approve every send with its code. |
| 다음 행동 (처음 연결) | 명령을 실행하고 설치기 안내를 따르세요. 끝나면 승인을 묻는 모드(예: `approval_policy = "on-request"`)로 Codex를 열고 세션을 연결하세요. YOLO(전체 자동) 모드에서는 연결되지 않습니다. | Run the command and follow the installer. When it finishes, open Codex in a mode that asks for approval (for example `approval_policy = "on-request"`) and connect the session. Connection does not work in YOLO (full auto) mode. |
| 다음 행동 (기존 연결 계속) | 같은 명령을 실행하면 기존 설치를 갱신합니다. 에이전트·키·친구는 그대로입니다. 갱신 뒤 Codex를 다시 열고 세션을 다시 연결하세요. | Running the same command updates your existing install. Your agent, key, and friends stay the same. After the update, reopen Codex and connect the session again. |
| 세션 연결 | 앱에서 확인 필요. 서버는 세션 연결을 확인하지 않으므로 Codex에서 `knowslink_session`을 `status`로 확인하세요. Codex 화면에는 도착 알림이 없습니다. 새 메시지는 받기를 요청해야 보입니다. | Check in the app. The server can't see session connection, so ask Codex to run `knowslink_session` with `status`. Codex shows no arrival notice; ask it to receive to see new messages. |
| details 요약 | 수동 설치 문서 | Manual install guide (in Korean) |

Codex details 링크는 [adapters Codex 절](../../../../adapters/README.md#codex-작업-세션-연결--sar-codex-recv-001-dev)이다. Claude Code는 현행 링크를 유지한다.

#### 자기 비평과 수락

- 첫 구상은 명령 옆에 복사 버튼을 두었다. 새 JS가 필요하고 D04의 'JS 없음' 원칙과 어긋나므로 친구 코드와 같은 선택·복사 안내로 바꿨다.
- 준비물 첫머리에 OS를 둔다. 검증 안 된 OS 사용자가 명령을 실행하기 전에 알 수 있다.
- '직접 확인할 일'을 '자동으로 하는 일' 바로 뒤에 둔다. 자동 설정이 승인을 줄인다고 오해하지 않게 한다.
- 수락: KO·EN 문구가 위 표와 같다. 1440·390·320px에서 명령 입력란이 넘치지 않는다. 키보드로 선택지→안내 보기→명령 입력란→details 순서로 이동한다. 200% 확대에서 읽힌다. 6번 선택지에 명령이 없다.

#### 시각 검수 결과 — SAR-AUTOSETUP-001-ART (main `11d91648`, 2026-10-23)

| 항목 | 결과 |
|---|---|
| U8 선택지 표(1–6번 KO·EN) | 화면과 일치. 4·5번은 '자동 설정', 5번은 새 선택지(`use=codex`), 6번 보조 문구에서 Codex를 뺐다. |
| 안내 상자 카피·순서(Claude Code·Codex × KO·EN × 처음·기존) | 48행 모두 일치. 새 색·버튼·JS 없음. |
| 명령 범위(AC-10) | 명령 입력란은 `use=claude`·`use=codex`에만 있다. web·grok·dots·6번에는 없다. |
| 가로 넘침·대비 | 1440·390·320 모두 가로 스크롤 0. 도움말 대비 6.19:1. |
| 열린 결함 | V-01(medium, dev 수정 필요): 320px에서 명령이 4줄인데 입력란이 3행이라 마지막 줄이 가려진다. 390px는 3줄에 딱 맞는다. 수정은 입력란을 4행으로 늘리거나 좁은 폭에서 `min-height`를 키우는 방식이며 새 토큰은 쓰지 않는다. V-02–V-04는 low로 수용했다. |

위 '수락'의 "320px에서 명령 입력란이 넘치지 않는다"는 마지막 줄이 보이는 것까지 포함한다. V-01 수정 뒤 320px 4조합(claude·codex × KO·EN)을 다시 확인한다.

### 연결 단순화와 내 에이전트 정리 — SAR-CONNECT-SIMPLE-001

설계안이다. 구현 전이다. 제품 규칙은 [D02 SAR-CONNECT-SIMPLE-001](../../planning/product-specs/SAR-CONNECT-SIMPLE-001.md)이다. 바꾸는 화면은 '내 에이전트' 목록 위쪽, 연결 링크 화면(`/connect/{id}`), 연결 중 Google 실패 화면이다. 다른 화면·'어디에서 사용할까요?' 선택지·자동 설정 안내 상자는 그대로다.

토큰·컴포넌트는 위 '현행 토큰과 형태'를 그대로 쓴다. 새 색·새 토큰·새 JS·복사 버튼·카드 그림자·알약 버튼·상태 색 배지를 쓰지 않는다. 상태는 텍스트로 보인다. 기존 `dl.kv`, `p.note`, `p.problem`, `details`, `button.primary`, `code.fp.big`, 폼 입력을 재사용한다.

#### 내 에이전트 목록 위쪽

```text
내 에이전트                                          (h1, 현행)
[어디에서 사용할까요?]                                (현행)

진행 중인 연결                                       (h2, 시도가 있을 때만)
  Node 로컬 클라이언트 · 10:42 시작                   (h3)
  동의 완료 · 키 저장 확인 필요                        (p)
  시작한 앱에서 저장 결과를 확인하세요. …              (p.muted)
  [이 시도 정리]                                     (button, 키 없는 시도만)

이름을 정할 에이전트                                  (h2, 이름 없는 agent가 있을 때만)
  이름 없는 에이전트 · a1b2                            (h3)
  에이전트 이름 [            ] [이름 저장]             (form, 접지 않음)

에이전트 목록                                        (현행 카드)
```

- 빈 절은 제목도 보이지 않는다(M-01).
- 진행 중인 연결의 상태 문구는 세 가지다: '지문 확인 대기', '동의 완료 · 키 저장 확인 필요', '만료'. 만료 시도는 A-02 정리 전까지만 보인다.
- '이 시도 정리'는 확인 단계 없이 바로 지운다. 키가 없는 시도라 잃는 자격이 없다. 키가 있는 agent의 철회는 현행처럼 카드의 접힌 `details.danger` 안에 둔다.
- 한도 안내 문구, '쓰지 않는 agent를 먼저 철회하세요' 거절 화면, '새 에이전트 만들기' 대체 버튼(활성 한도 경로)을 없앤다(M-03). 키 기록 포화(agent당 20개) 안내는 현행 유지다.
- 1440px은 기존 읽기 폭, 390·320px은 한 열이다. 이름 입력과 저장 버튼은 320px에서 위아래로 쌓인다.

| 요소 | KO | EN |
|---|---|---|
| 절 제목 | 진행 중인 연결 | Connections in progress |
| 시도 제목 | `<클라이언트 종류>` · `<HH:MM>` 시작 | `<client type>` · started `<HH:MM>` |
| 상태: 대기 | 지문 확인 대기. 연결 링크에서 지문을 비교하고 승인하세요. | Waiting for fingerprint check. Compare the fingerprint on the connection link and approve. |
| 상태: 동의 완료 | 동의 완료 · 키 저장 확인 필요. 시작한 앱에서 저장 결과를 확인하세요. 앱이 저장하면 이 시도는 에이전트 목록으로 옮겨집니다. | Approved · Key storage not confirmed yet. Check the result in the app you started from. When the app saves the key, this attempt moves to your agent list. |
| 상태: 만료 | 만료. 키가 연결되지 않았습니다. 시작한 앱에서 다시 연결하세요. | Expired. No key was connected. Connect again from the app you started from. |
| 정리 버튼 | 이 시도 정리 | Remove this attempt |
| 정리 도움말 | 키가 저장되지 않은 시도라 지워도 잃는 것이 없습니다. 정리하지 않아도 하루 뒤 자동으로 정리됩니다. | No key was saved, so nothing is lost. If you leave it, it's removed automatically after a day. |
| 이름 절 제목 | 이름을 정할 에이전트 | Agents that need a name |
| 이름 도움말 | 친구에게도 이 이름이 보입니다. 이름은 1–40자로 입력하세요. | Friends see this name too. Use 1–40 characters. |

#### 연결 링크 화면 `/connect/{id}`

```text
내 클라이언트 연결                                    (h1)
[주의] 내 컴퓨터에서 시작한 요청만 승인하세요 …         (p.note.caution, 현행)
1. 지문 비교                                          (h2)
   내 컴퓨터에 표시된 지문과 같은지 확인하세요.
   SHA256:xxxx…                                      (code.fp.big)
   다르면 취소하세요.
2. 에이전트 이름                                      (h2)
   [Claude Code 2              ]                      (기본값 채움, A-04)
   친구에게도 이 이름이 보입니다. …                     (p.muted)
[지문이 같습니다. 연결합니다]                          (button.primary)
[취소]                                               (button)
▸ 승인하면 허용되는 것 · 기한 10:52까지                 (details: 현행 '승인하면' 목록과 클라이언트·기한)
```

- 순서는 지문, 이름, 연결 버튼 하나다(W-03). 현행 '승인하면' 목록은 지우지 않고 `details` 안으로 옮긴다. 요약 줄에 기한을 함께 보여 접힌 상태에서도 기한을 안다.
- 연결 버튼 문구가 지문 확인을 담는다. 지문 비교와 버튼 누름은 생략하지 않는다.
- 기본 이름은 서버가 정한 값이다(A-04). 클라이언트가 보낸 호스트 정보는 기본 이름 계산에만 쓰고 사용자가 고칠 수 있다.

| 요소 | KO | EN |
|---|---|---|
| 1단계 제목 | 1. 지문 비교 | 1. Compare the fingerprint |
| 2단계 제목 | 2. 에이전트 이름 | 2. Agent name |
| 연결 버튼 | 지문이 같습니다. 연결합니다 | The fingerprint matches. Connect |
| details 요약 | 승인하면 허용되는 것 · 기한 `<시각>`까지 | What approval allows · until `<time>` |

링크를 다시 열 때(W-01). 연결을 시작한 회원이 로그인한 경우에만 아래 상태별 화면을 보인다. 다른 회원·로그아웃 상태는 현행 401 일반 문구와 로그인 링크다.

| 상태 | 제목 | KO 본문 | EN 본문 | 링크 |
|---|---|---|---|---|
| 동의 직후·키 저장 대기 | 연결을 승인했습니다 | 시작한 앱으로 돌아가세요. 앱이 키를 저장하면 '내 에이전트'에서 `<이름>`의 상태가 '키 등록됨'으로 바뀝니다. 승인만으로는 아직 연결되지 않았습니다. | Go back to the app you started from. When it saves the key, `<name>` shows as "Key registered" in My agents. Approval alone doesn't connect it yet. | 내 에이전트로 이동 |
| 완료(키 저장됨) | 이 연결은 완료됐습니다 | `<이름>` 에이전트에 키가 등록됐습니다. 이 링크는 다시 쓸 수 없습니다. 메시지를 받을 세션은 앱에서 연결하세요. | The key for `<name>` is registered. This link can't be used again. Connect the receiving session in the app. | 내 에이전트로 이동 |
| 만료 | 이 연결은 만료됐습니다 | 키가 연결되지 않았습니다. 시작한 앱에서 다시 연결하면 새 링크를 만듭니다. | No key was connected. Connect again from the app you started from to get a new link. | 내 에이전트로 이동 |
| 취소 | 이 연결은 취소됐습니다 | 키가 연결되지 않았습니다. 시작한 앱에서 다시 연결하면 새 링크를 만듭니다. | No key was connected. Connect again from the app you started from to get a new link. | 내 에이전트로 이동 |

완료 화면은 오류 스타일(`p.problem`, `role="alert"`)을 쓰지 않는다. `p.note`와 `role="status"`를 쓴다. 만료·취소는 `p.problem`이다.

#### 연결 중 Google 실패 화면

현행 403 문구 하나를 원인별 세 문구로 나눈다(W-02). 다른 계정의 이메일·이름·존재는 보이지 않는다. '다른 계정 선택'은 Google 계정 선택 화면부터 다시 시작한다.

| 원인 | KO | EN | 행동 |
|---|---|---|---|
| 이 연결을 다른 Google 계정이 이미 시작함 | 이 연결은 다른 Google 계정으로 시작됐습니다. 처음 고른 Google 계정으로 다시 선택하세요. | This connection was started with a different Google account. Choose the account you used first. | [다른 계정 선택] |
| 새 가입이 닫힘 | 이 Google 계정은 아직 KnowsLink 회원이 아니며 지금은 새 가입을 받지 않습니다. 이미 쓰고 있는 Google 계정으로 다시 선택하세요. | This Google account isn't a KnowsLink member yet, and sign-ups are closed right now. Choose the Google account you already use. | [다른 계정 선택] |
| 로그인 기한 지남·연결 만료 | 로그인 시간이 지났습니다. 시작한 앱에서 연결 링크를 다시 여세요. | Sign-in timed out. Open the connection link again from the app you started from. | 없음 |
| 그 밖의 실패 | 현행 문구 | 현행 문구 | 현행 |

#### 자기 비평과 수락

- 첫 구상은 진행 중인 시도를 agent 카드 안의 상태로 섞었다. 실패 시도가 agent처럼 보여 C-02를 되풀이하므로 별도 절로 분리했다.
- 첫 구상은 '이 시도 정리'에 확인 단계를 두었다. 잃는 것이 없는 동작에 확인을 두면 정리가 번거로워 C-01의 쌓임이 남으므로 뺐다. 키 있는 agent 철회의 확인 위치는 그대로다.
- 연결 링크 화면의 권한 목록을 지우지 않고 접었다. 사람 확인 범위(AS-04)는 그대로이고 첫 화면의 결정만 줄인다.
- 수락: 위 표의 KO·EN 문구가 같다. 1440·390·320px에서 가로 스크롤 0, 이름 입력과 버튼이 가려지지 않는다. 키보드로 진행 중인 연결의 정리 버튼, 이름 저장, 연결 링크의 지문→이름→연결→취소→details 순서로 이동한다. 200% 확대에서 읽힌다. 한도 문구('회원당', '5개', '먼저 철회')가 회원 화면에 없다.
- 시각 검수: DEV 후보에서 designer가 직접 검수한다(ART, [tester 지시서 초안](../../../handovers/pending/SAR-CONNECT-SIMPLE-001-tester.md)의 인계 조건).

## 이전 설계와 검수 기록 — 비현행

다음 1–14절은 2026-10-11–12의 이전 관리 UI 명세와 검수 이력이다. 기존 판정·제약·미확인 원문을 보존한다. 사람 친구·대화·현행 토큰/레이아웃/노출 범위에는 위 현행 D04와 HP 요구사항을 적용한다.

## 1. 목적과 책임

link.knowslog.com의 회원 웹(로그인 → 내 에이전트 → 연결 관리)을 상용 서비스 수준으로 개선하는 제품 요구를 정의한다.
designer는 제품 방향·정보 구조·문구·시각 기준·수락 기준을 결정한다. DEV는 구현 방법과 파일별 기술 계획을 결정한다.
이 문서는 디자인 인계 문서다. 구현 후 직접 시각 검수는 후속 고정 SHA에서 별도로 수행한다. 이 문서는 UI 최종 수락이 아니다.

이번 범위가 아닌 것:
- 서버 상태·권한·재인증·관계 규칙·rate 정책·키/agent 한도 변경. 기존 규칙은 그대로 유지한다.
- 새 CSS/JS 라이브러리, 외부 폰트, 유료 에셋, 새 프런트엔드 프레임워크.
- 채팅 화면, 장기 메시지 타임라인, 사용량·결제·팀 기능.
- 다크 모드(후속 후보). 이번 색상 토큰은 밝은 테마만 정의한다.

## 2. 현재 화면 관찰 (근거)

Orca 내장 브라우저로 관찰했다.
- 회원 홈 `/home`: 이미 로그인된 브라우저 세션으로 1440×900과 390×844 뷰포트를 직접 확인했다. 에이전트 4개 상태다.
- 로그인 `/`: 로그아웃하지 않고 운영 응답 HTML을 받아 로컬 파일로 렌더링했다. 같은 CSS이므로 같은 화면으로 판정한다. 운영 `/`에는 Google 로그인만 보였다. 이메일 코드 로그인 입력은 운영 설정에서 꺼져 있어 보지 못했다.
- 이메일·식별자·키 값은 이 문서에 쓰지 않았다. 캡처 PNG는 Git에 넣지 않았다(개인정보 최소화).

| ID | 관찰 | 영향 |
|---|---|---|
| O-01 | 홈 문서 높이가 1440px 폭에서 6394px, 390px 폭에서 7986px이다(에이전트 4개). 섹션 6개(내 신원·내 agent·관계·연결 확인·승인 요청·로그아웃)가 같은 흰 박스로 이어진다. | 첫 화면에서 "내 에이전트 목록"과 "지금 할 일"이 보이지 않는다. 정보 위계가 없다. |
| O-02 | 첫 화면의 첫 박스가 "회원 식별자" 원문(`mem_…`)과 세션 만료 시각이다. | 일반 사용자에게 가치가 낮은 값이 가장 높은 위치를 차지한다. |
| O-03 | 에이전트 한 개가 ID 제목·복사용 입력·키 행·키 철회·클라이언트 선택·연결 방식·연결 수단 발급·agent 철회까지 항상 펼쳐진다. 4개면 화면 대부분이 폼이다. | 목록 훑기와 에이전트 비교가 불가능하다. |
| O-04 | 위험 동작(키 철회, agent 철회)이 일반 동작 버튼과 같은 회색 기본 버튼이다. "연결 수단 발급" 바로 아래에 "이 agent와 모든 키·관계 철회"가 있다. | 오조작 위험이 있다. 기능 발견성은 높지만 위험 구분이 없다. |
| O-05 | `overflow-wrap:anywhere`가 body에 적용되어 390px에서 한국어 단어가 중간에서 끊긴다("사용합/니다", "끝납니/다"). | 본문 가독성이 낮다. 긴 ID·지문에만 강제 줄바꿈이 필요하다. |
| O-06 | 390px에서 복사용 readonly 입력의 ID가 잘린다(`agent_077c666294c4eb2…`). select 안의 문구도 잘린다("새 키 등록 (활성 최대 3가…"). | 값 확인이 어렵다. |
| O-07 | 키 ID·SHA256 지문이 본문 폰트 크기의 긴 한 줄이다. 지문 대조는 연결 승인의 핵심 보안 동작이다. | 지문 대조 오류 위험이 있다. |
| O-08 | 상태가 "상태: 연결 완료" 문장 한 줄이다. 색·배지가 없다. | 상태 스캔이 어렵다. |
| O-09 | 개발자 명령(`node adapters/dist/...`)이 홈 중간에 상시 노출된다. "다닷·외부 앱 실제 연결은 후속 검증입니다", "관계 세대: N", "계정 비활성화 상태: 준비 중" 같은 내부 말투가 있다. | 상용 서비스 첫인상에 맞지 않는다. 일부 값은 사용자가 쓸 일이 없다. |
| O-10 | 버튼은 브라우저 기본 스타일이다. 높이 52px로 터치 크기는 충분하다. 1440px에서 본문 폭은 약 760px이고 가운데 정렬이다. | 터치 크기는 유지한다. 시각 위계와 브랜드가 없다. |
| O-11 | 로그인 화면은 제목 "KnowsLink 시작", 박스 한 개, 설명 한 문단, "Google로 계속" 버튼 한 개다. 서비스가 무엇인지 설명이 없다. | 처음 방문자가 서비스 목적을 알 수 없다. |
| O-12 | 가로 overflow는 없다(1440·390 모두 `scrollWidth`가 뷰포트 폭 이내). | 이 상태를 유지한다. |

관찰하지 않은 것: 이메일 로그인 입력·확인 대기·오류 화면, 연결 관리(`/home/connections/{id}`)·기기 승인(`/connect/{id}`)·receipt·gate 화면의 렌더 결과. 이 화면들은 템플릿 코드(`internal/relay/member.go`, `member_agents.go`, `http.go`)로 문구와 상태 분기만 확인했다. 후속 시각 검수에서 직접 확인한다.

## 3. 제품 방향

한 문장: **"내 에이전트가 연결되어 있는지, 지금 내가 할 일이 있는지를 10초 안에 알려 주는 조용한 관리 화면."**

사용자 목표 순서:
1. Google로 가입·로그인한다.
2. 내 에이전트 목록과 각 상태를 본다.
3. 새 에이전트를 만들고 내 클라이언트 컴퓨터에 연결한다.
4. 다른 에이전트와 관계를 초대·수락한다.
5. 필요하면 키·에이전트를 철회하고 로그아웃한다.

디자인 원칙:
- 목록은 한눈에, 동작은 필요할 때. 에이전트 한 개는 접힌 한 행(식별자·상태·키 수)이 기본이다.
- 위험은 분리한다. 되돌릴 수 없는 동작은 시각적으로 따로 두고, 결과를 먼저 말한다.
- 있는 사실만 표시한다. 서버가 알지 못하는 상태("온라인", "최근 접속", "보안 검증됨")는 만들지 않는다.
- 장식은 쓰지 않는다. 무작위 그라디언트·장식용 통계·모든 정보의 동일한 카드·알약(pill) 버튼을 피한다. 색은 의미(상태·위험·주 동작)에만 쓴다.
- 전문 용어를 줄인다. 필요한 곳에는 한 문장 설명을 붙인다.

## 4. 정보 구조

### 4.1 화면 지도 (기존 경로 유지)

| 경로 | 화면 | 변경 |
|---|---|---|
| `/` | 로그인 | 서비스 설명 추가, 버튼 위계 정리 |
| `/auth/verify` | 확인 코드 입력(이메일 로그인 설정 시) | 시각 토큰 적용, 문구 정리 |
| `/home` | 내 에이전트 홈 | 구조 재배치(4.2) |
| `/home/connections/{id}` | 연결 관리 | 3단계 진행 표시(5.4) |
| `/connect/{id}` | 클라이언트 연결 승인(기기 승인) | 지문 대조 강조(5.5) |
| `/home/receipts` | 연결 확인 receipt | 시각 토큰·정의 목록 정리 |
| `/home/gates/{id}` | 승인 요청(gate) | 시각 토큰만 공유. 동작·문구 의미는 유지 |
| 거절 화면(refusal) | 오류와 다음 동작 | 오류 규칙(7) 적용 |

새 경로·JS 없이 가능하다. DEV가 구조상 필요하다고 판단하면 새 경로를 추가할 수 있으나, 권한·owner 검사와 CSP는 기존과 같아야 한다.

### 4.2 `/home` 섹션 순서

위에서 아래로 다음 순서다. 번호는 화면 순서다.

1. **헤더 영역**: 제품명 "KnowsLink", 확인된 이메일(마스킹된 값 그대로), "로그아웃" 링크형 버튼. 회원 식별자와 세션 만료 시각은 헤더에서 뺀다(8단계 "계정"으로 이동).
2. **확인 필요 띠** (조건부): 받은 관계 초대 N건, 승인 요청 N건, 재인증 필요 등 사용자가 해야 할 일이 있을 때만 표시한다. 각 항목은 해당 섹션 앵커 링크다. 없으면 표시하지 않는다.
3. **섹션 이동 링크**: "에이전트 · 관계 · 연결 확인 · 승인 요청 · 계정". 같은 페이지 앵커다. 960px 이상에서는 왼쪽 고정 열(sticky), 미만에서는 상단 가로 줄이다. JS가 필요 없다.
4. **내 에이전트** (주 영역): 제목 옆에 "활성 N/5" 텍스트와 "새 에이전트 만들기" 주 버튼. 에이전트 행 목록(5.3).
5. **관계**: 받은 초대 → 보낸 초대 → 연결됨 → 종료된 관계 순서. 새 초대 폼은 접힌 `details`("새 관계 초대")로 둔다. 받은 초대가 있으면 그 항목은 펼쳐진 상태다.
6. **승인 요청**: 목록만. 없으면 한 줄 빈 상태.
7. **연결 확인(receipt)**: 접힌 `details`("연결 확인 방법"). 안의 명령은 코드 블록이다.
8. **계정과 보안**: 회원 식별자, 세션 만료, 로그인 신원 다시 확인, 모든 브라우저 로그아웃. 계정 비활성화 "준비 중" 문구는 제거한다(기능이 생길 때 추가).

## 5. 화면별 명세

### 5.1 로그인 `/`

구성:
- 제품명 "KnowsLink"와 한 줄 설명: "내 AI 에이전트끼리 안전하게 메시지를 주고받는 연결 서비스". (최종 한 줄은 제품 소개 정본과 맞추되 의미를 바꾸지 않는다.)
- 설명 3줄(목록): "Google 계정으로 가입 없이 바로 시작합니다." / "연결한 에이전트는 계정마다 따로 관리됩니다." / "Google 비밀번호는 Google 화면에만 입력합니다."
- 주 버튼 "Google로 계속". 전체 폭 버튼(390px), 데스크톱은 폭 자동(최소 240px).
- 이메일 로그인 설정이 켜진 경우: "또는" 구분선 뒤 이메일 입력과 "확인 코드 받기". 설정이 꺼져 있으면 영역 자체를 숨긴다(현재 동작 유지).
- 둘 다 꺼진 경우 안내 문구는 "로그인을 준비하고 있습니다. 잠시 뒤 다시 시도하세요."
- 로그아웃·세션 종료 후 안내 띠(`role="status"`)는 현재처럼 상단에 표시한다.

금지: Google 로고 이미지·외부 이미지(CSP가 이미지를 막는다). "Google로 계속" 텍스트 버튼만 쓴다.

### 5.2 확인 대기 `/auth/verify` (이메일 로그인 설정 시)

- 제목 "확인 코드 입력". 첫 문장은 "{마스킹 이메일}로 코드를 보냈습니다. 아직 로그인되지 않았습니다."
- 코드 입력은 큰 숫자 입력 1개(`inputmode=numeric`, `autocomplete=one-time-code`), 주 버튼 "확인하고 로그인".
- 유효 시각·남은 시도 횟수·재발송 가능 시점을 입력 아래 보조 텍스트로 둔다. 기존 값(10분, 5회, 60초)을 바꾸지 않는다.
- "코드 다시 받기"는 보조 버튼, "다른 주소로 시작"은 링크.

### 5.3 내 에이전트 목록 (홈 4단계)

에이전트 한 개는 접힘 행(`details/summary`)이다.

접힌 행에 보이는 것(한 줄, 390px에서는 두 줄):
- 식별자(`code`, 12자 이상 줄임 없이 줄바꿈 허용), 상태 배지, "키 N개 활성".
- 상태 배지는 기존 4개 상태를 그대로 쓴다: 미연결 / 연결 대기 / 연결 완료 / 철회. 색과 텍스트를 함께 쓴다.

펼쳤을 때 보이는 것 (위에서 아래로):
1. 식별자 전체(선택하기 쉬운 `user-select:all` 코드). JS 복사 버튼은 쓰지 않는다(CSP가 스크립트를 막는다). "식별자를 눌러 전체 선택 후 복사하세요" 보조 문구를 한 번만 섹션 머리말에 둔다.
2. 키 목록: 키 ID, 지문, 상태. 지문은 값을 바꾸지 않고 그대로 쓴다. 모노스페이스 14px, 줄바꿈 허용, `letter-spacing:.02em`.
3. **연결 추가** 영역: 지원 클라이언트·연결 방식 선택과 "연결 수단 발급". 연결 방식 선택지 문구는 짧게: "새 키 추가(최대 3개)" / "키 교체(기존 키 모두 철회)". 키 기록이 가득 찬 경우는 기존 안내문(키 기록 보호 상한)을 이 자리에 `note` 스타일로 보여 준다.
4. **위험 구역** 영역(접힌 `details`, 제목 "이 에이전트 철회·키 철회"): 키별 "이 키 철회", "이 에이전트와 모든 키·관계 철회". 영역 맨 위에 결과 문장 "철회하면 이 에이전트의 모든 키와 관계가 끝나며 되돌릴 수 없습니다."를 둔다. 위험 버튼은 위험 색 외곽선 버튼이다.

철회된 에이전트: 접힌 상태로 시작하고 회색 배지 "철회". 펼치면 기존 보존·정리 안내(최소 24시간)를 `note`로 보여 준다. 연결·철회 동작은 숨긴다(현재 동작 유지).

빈 상태(에이전트 0개): "아직 연결한 에이전트가 없습니다. 아래 버튼으로 첫 에이전트를 만드세요." + 주 버튼. 단계 1·2·3 짧은 안내("만들기 → 내 컴퓨터에서 연결 → 관계 초대")를 번호 목록으로 둔다.

한도: "활성 N/5". 5개에 도달하면 "새 에이전트 만들기"는 비활성이 아니라 그대로 두고, 누르면 서버의 기존 한도 안내(활성/기록 포화)가 나온다. 사유 없는 비활성 버튼은 만들지 않는다.

"연결 완료"의 뜻: 활성 키가 있다는 뜻이다. 클라이언트가 지금 켜져 있다는 뜻이 아니다. 헤더 섹션 머리말에 한 문장으로 적는다: "연결 완료는 키가 등록되었다는 뜻이며 지금 접속 중이라는 뜻은 아닙니다." 실제 접속 상태 데이터가 생기기 전에는 "온라인·최근 접속" 표시를 만들지 않는다.

### 5.4 연결 관리 `/home/connections/{id}`

3단계 진행 표시를 화면 위에 둔다. 현재 서버 상태와 대응한다.

| 단계 | 서버 상태 | 사용자에게 보이는 이름 |
|---|---|---|
| 1 | waiting | 연결 수단 발급됨: 내 컴퓨터에서 준비 |
| 2 | prepared | 지문 확인: 내 컴퓨터에 표시된 지문과 같은지 확인 |
| 3 | approved → consumed | 연결 마무리 → 연결 완료 |

종료 상태 expired·cancelled는 진행 표시 대신 상태 `note`("만료되었습니다. 새 키가 연결되지 않았습니다. 홈에서 다시 시작하세요.")와 "홈으로" 링크를 보여 준다.

규칙:
- 진행 표시는 숫자 + 텍스트다. 색만으로 구분하지 않는다. 현재 단계에 `aria-current="step"`을 둔다.
- 연결 수단(token)은 한 번만 보여 주는 값이다. 읽기 전용 입력 대신 `user-select:all` 코드 블록으로 보여 주고 "대화·로그·캡처에 붙이지 마세요"를 값 바로 위에 둔다.
- `node adapters/dist/connect.js ...` 명령은 코드 블록 두 개(준비, 마무리)로 보여 준다. 앞에 "내 컴퓨터(Node 22)에서 실행" 한 줄을 둔다.
- prepared 단계: 지문을 화면에서 가장 큰 모노스페이스 값으로 보여 준다(5.5와 같은 스타일). 주 버튼 "지문이 같습니다. 연결 승인". 보조 버튼 "취소".
- 만료 시각은 KST 표기를 유지한다. 형식 "10월 12일 07:55까지"처럼 짧은 날짜를 우선한다. 정확한 초 단위 값은 보조 텍스트로 둘 수 있다.
- "공개키 준비 상태 확인" 링크는 "상태 새로고침"으로 바꾸고 같은 URL로 이동한다. 자동 새로고침은 하지 않는다(JS 없음).

### 5.5 클라이언트 연결 승인 `/connect/{id}`

- 맨 위에 경고 띠: "내 컴퓨터에서 시작한 요청만 승인하세요. 다른 사람이 보낸 링크는 승인하지 마세요." (`note` 스타일, 위험 색 왼쪽 선).
- 지문은 큰 모노스페이스 블록이다. 위에 "내 컴퓨터에 표시된 지문과 같은지 확인하세요", 아래에 "다르면 취소하세요".
- 권한 설명은 목록: "이 에이전트의 보내기·받기만 허용" / "다른 에이전트와의 관계는 따로 초대·수락" / "개인키는 내 컴퓨터에만 보관". 기존 의미를 그대로 쓴다.
- 버튼 순서: 주 버튼 "지문 확인 후 연결 승인", 보조 버튼 "취소". 승인 후 화면은 `role="status"`로 "승인했습니다. 내 컴퓨터에서 저장 결과를 확인하세요."

### 5.6 관계

- 항목 한 개 = 한 행. 보이는 것: 내 에이전트 → 상대 에이전트(식별자), 방향 라벨(보낸 초대·받은 초대·연결됨·종료), 상태 배지, 초대 기한(대기 중일 때만).
- "관계 세대: N" 숫자는 화면에서 숨긴다. 폼의 hidden 값은 유지한다(서버 검증에 필요).
- 받은 초대: 주 버튼 "수락", 보조 위험 버튼 "거절". 두 버튼 위에 "수락하기 전에는 메시지를 주고받을 수 없습니다."
- 연결됨: "관계 철회"는 위험 색 외곽선 버튼이며 행 오른쪽 아래에 둔다.
- 종료됨: 기존 안내("메시지를 보낼 수 없습니다. 다시 연결하려면 새 초대를 보내고 상대가 새로 수락해야 합니다.")와 "새 초대 보내기" 보조 버튼.
- 새 초대 폼: 라벨 "내 에이전트"(select), "상대 에이전트 식별자"(입력). 도움말 "가입 초대가 아닙니다. 상대 소유자가 수락해야 메시지가 허용됩니다."
- 상대 이메일·타 회원 정보는 어떤 경우에도 표시하지 않는다(기존 규칙).

### 5.7 승인 요청과 receipt

- 승인 요청 목록: 행마다 원요청 ID(링크), 상태 배지. 없으면 "승인 요청이 없습니다."
- gate 상세(`/home/gates/{id}`): 이번에는 시각 토큰(글꼴·색·버튼 스타일)만 공유한다. 본문이 판단 근거이고 hint는 참고라는 구분, 비활성 상태의 버튼 숨김 규칙은 그대로다.
- receipt 화면: 정의 목록(`dl`)을 2열(라벨/값)로 정렬한다. 상태 값 옆에 배지를 쓰되 "queued는 상대 수신 성공이 아닙니다" 문장은 유지한다. 자동 wake·자동 답장 없음 설명은 유지한다.

### 5.8 계정과 보안

- 회원 식별자, 세션 최대 유지 시각(KST), 60분 무활동 만료.
- "로그인 신원 다시 확인"(보조 버튼)과 설명 "키·연결 권한 변경에는 5분 안의 재확인이 필요합니다."
- 재확인이 필요한 동작 전에는 해당 영역 위에 `note`로 재확인 안내를 먼저 보여 준다. 서버가 거부하고 나서야 알게 하지 않는다(현재 서버 동작은 변경하지 않는다. 화면은 이미 계산된 `Recent` 값을 쓴다).
- "로그아웃"(이 브라우저) 보조 버튼, "모든 브라우저에서 로그아웃"은 `Recent`일 때만 활성. 아닐 때는 재확인 안내를 보여 준다(현재 동작 유지).
- 문구 고정: "로그아웃은 브라우저 세션만 끝냅니다. 연결한 에이전트의 키와 자격은 철회되지 않습니다."

## 6. 한국어 문구 표

원칙: 합쇼체(~합니다), 한 문장 한 사실, 같은 개념에 같은 용어. 영문 `agent`는 화면에서 "에이전트"로 통일한다(식별자 값 `agent_…`은 원문).

| 위치 | 현재 | 제안 |
|---|---|---|
| 홈 제목 | 내 KnowsLink | 내 에이전트 (문서 제목: 내 에이전트 · KnowsLink) |
| 로그인 제목 | KnowsLink 시작 | KnowsLink에 로그인 |
| 로그인 버튼 | Google로 계속 | Google로 계속 (유지) |
| 에이전트 추가 | 새 agent 만들기 | 새 에이전트 만들기 |
| 연결 방식 | 새 키 등록 (활성 최대 3개) | 새 키 추가 (최대 3개) |
| 연결 방식 | 회전 (완료 시 기존 키 전체 철회) | 키 교체 (완료하면 기존 키가 모두 철회됨) |
| 키 철회 | 선택한 키 철회 | 이 키 철회 |
| 에이전트 철회 | 이 agent와 모든 키·관계 철회 | 이 에이전트 철회 (모든 키·관계 종료) |
| 지문 승인 | 대상·지문 확인 후 연결 승인 | 지문이 같습니다. 연결 승인 |
| 관계 | 수신 owner로 수락 / 수신 owner로 거절 | 수락 / 거절 |
| 재확인 | 로그인 신원 다시 확인 | 로그인 다시 확인 |
| 상태 안내 | 상태: 연결 완료 | 배지 "연결 완료" + 보조 "키 N개 활성" |
| 계정 비활성화 | 상태: 준비 중. 아직 요청할 수 없습니다. | (제거) |
| 내부 문구 | 다닷·외부 앱 실제 연결은 후속 검증입니다. | (제거. 지원 클라이언트 표에서 "Node 22 로컬 클라이언트"만 안내) |
| 오류 접두 | 문제: … / 안내: … | 접두 대신 아이콘 문자 + 굵은 라벨("오류", "안내")과 같은 문장 |

서버가 주는 오류·한도 문구(`safeProblem`, `agentLimits`, `keyRecordsFull`, `inviteNotice`)의 의미는 바꾸지 않는다. 용어 통일("agent" → "에이전트")과 문장 분리는 허용한다. 약속하지 않는 것(정확한 정리 시각, 영구 삭제, 권한 복구)은 그대로 약속하지 않는다.

## 7. 오류·빈 상태·제한 상태

| 상황 | 상태 코드(현재) | 표시 | 다음 동작 |
|---|---|---|---|
| 세션 만료 | 303 → `/?n=expired` | 로그인 화면 상단 `note` "세션이 끝났습니다. 다시 로그인하세요." | Google 로그인 |
| 요청 과다 | 429 | `problem`: "요청이 많습니다. {시각} 이후 다시 시도하세요." | "홈으로" 링크 |
| 재확인 필요 | 403 | `problem` + "로그인 다시 확인" 버튼 | 재확인 후 홈 복귀 |
| 한도·기록 포화 | 409 | 서버 안내문 그대로 + (가능 시) "새 에이전트 만들기" | 홈 복귀 또는 새 에이전트 |
| 연결 만료·취소 | 200 | `note` "새 키가 연결되지 않았습니다" | 홈에서 다시 시작 |
| 권한 없는 연결 조회 | 403 | `problem` "연결을 조회할 수 없습니다." | 홈 복귀 |
| 일시 오류 | 503 | `problem` "일시적으로 처리할 수 없습니다. 잠시 뒤 다시 시도하세요." | 홈 복귀 |
| 에이전트 0개 | 200 | 5.3 빈 상태 | "새 에이전트 만들기" |
| 관계 0개 | 200 | "관계가 없습니다. 에이전트를 연결한 뒤 상대 식별자로 초대하세요." | 새 관계 초대 열기 |
| 승인 요청 0개 | 200 | "승인 요청이 없습니다." | 없음 |

모든 오류 화면에는 실제로 가능한 다음 동작 링크가 하나 이상 있다. 오류를 성공 배너로 바꾸지 않는다. 에러 문구에 비밀값·상대 이메일·다른 회원 정보를 넣지 않는다.

## 8. 비주얼 시스템

CSS는 기존 `memberStyle` 방식(인라인 `<style>`)으로 쓸 수 있다. CSP가 `default-src 'none'; style-src 'unsafe-inline'`이므로 스크립트·이미지·외부 폰트·외부 스타일시트는 쓸 수 없다. 이 기준은 모두 CSS 변수와 기본 HTML 요소로 구현한다. 아이콘은 쓰지 않거나 유니코드 문자만 쓴다.

### 8.1 색상 토큰 (밝은 테마)

대비 비율은 WCAG 상대 휘도 공식으로 직접 계산했다.

| 토큰 | 값 | 용도 | 대비 |
|---|---|---|---|
| `--bg` | `#f4f6f8` | 페이지 배경 | 본문 `--text`와 14.78 |
| `--surface` | `#ffffff` | 영역 배경 | `--text`와 16.01 |
| `--text` | `#17232b` | 본문 | — |
| `--muted` | `#44525c` | 보조 텍스트 | 흰 배경 8.05 / `--bg` 7.43 |
| `--line` | `#c9d2d9` | 구분선·입력 테두리(장식용 선) | 입력 테두리는 `--line-strong`를 쓴다 |
| `--line-strong` | `#6b7a84` | 입력·버튼 테두리 | 흰 배경 4.43 (비텍스트 3:1 이상) |
| `--accent` | `#164da0` | 주 버튼 배경, 링크 | 흰 글자와 8.07 |
| `--accent-soft` | `#e8effb` | 안내 배경 | `--accent` 글자와 6.98 |
| `--danger` | `#a01616` | 위험 버튼 글자·테두리, 오류 | 흰 배경 8.01 |
| `--danger-soft` | `#fdecec` | 오류 배경 | `--danger`와 7.01 |
| `--ok` | `#146c43` | 연결 완료 배지 글자 | 흰 배경 6.45 / `#e7f5ee` 5.74 |
| `--ok-soft` | `#e7f5ee` | 연결 완료 배지 배경 | — |
| `--warn` | `#7a4b00` | 대기·주의 배지 글자 | `#fff4d6` 배경 6.76 |
| `--warn-soft` | `#fff4d6` | 대기·주의 배경 | — |

토큰 값은 DEV가 구현하며 동일한 대비 기준(텍스트 4.5:1 이상, 비텍스트 3:1 이상)을 지킨다. 값이 바뀌면 설계 의도를 유지하는 범위에서 조정할 수 있다. 조정 값은 검수 때 기록한다. 기존 색(`#17232b`, `#164da0`, `#a01616`, `#f5f7f8`)은 새 토큰과 같은 계열이다.

### 8.2 타이포그래피

- 글꼴: `system-ui, -apple-system, "Segoe UI", "Apple SD Gothic Neo", "Malgun Gothic", "Noto Sans KR", sans-serif`. 외부 폰트는 쓰지 않는다.
- 모노스페이스(식별자·지문·명령): `ui-monospace, SFMono-Regular, Consolas, "Courier New", monospace`.
- 크기: 본문 16px(현재 18px에서 축소, 줄 길이 개선), 보조 14px, 모노 14px(지문 값은 확인 화면에서 16px), 섹션 제목 20px/600, 페이지 제목 26px/700. 줄 높이 본문 1.65, 제목 1.3.
- 한국어 줄바꿈: body에 `word-break:keep-all; overflow-wrap:break-word;`. 식별자·지문·코드·URL에만 `overflow-wrap:anywhere`를 쓴다. 한국어 본문에 `anywhere`를 쓰지 않는다. 긴 한국어 단어가 넘치면 `break-word`가 처리한다.
- 숫자·시각은 본문 폰트에 `font-variant-numeric:tabular-nums`.

### 8.3 간격·형태

- 간격 척도: 4 / 8 / 12 / 16 / 24 / 32 / 48px.
- 페이지 좌우 여백: 390px 폭에서 16px. 960px 이상에서 컨테이너 최대 폭 1040px(섹션 이동 열 200px + 본문 800px).
- 영역(`section`)은 `--surface` 배경, 1px `--line` 테두리, 모서리 8px. 영역 안의 하위 항목(에이전트 행)은 테두리 없이 위쪽 구분선만 쓴다(카드 중첩 금지).
- 그림자·그라디언트는 쓰지 않는다.

### 8.4 버튼과 입력

- 모든 버튼 높이 최소 44px(현재 52px 이상이면 유지 가능). 모서리 6px. 알약 모양 금지.
- 주 버튼: `--accent` 배경, 흰 글자, 600. 화면(영역)마다 주 버튼은 한 개를 원칙으로 한다.
- 보조 버튼: 흰 배경, `--line-strong` 테두리, `--text` 글자.
- 위험 버튼: 흰 배경, `--danger` 테두리·글자. 위험 구역 안에서만 쓴다. 채운 위험 버튼은 쓰지 않는다(오조작 방지와 시각 소음 감소).
- 비활성은 의미 있는 경우에만 쓴다. 비활성 사유는 가까운 텍스트로 쓴다.
- 입력: 높이 44px 이상, 1px `--line-strong` 테두리, 라벨은 항상 보인다(placeholder로 대체하지 않는다). 포커스는 8.6의 규칙을 쓴다.
- 모바일에서 영역 안 주 버튼은 전체 폭이다. 데스크톱은 내용 폭이다.

### 8.5 배지와 안내

| 상태 | 배지 | 안내 띠 |
|---|---|---|
| 연결 완료·활성 | `--ok` 글자 + `--ok-soft` 배경 + 앞에 "✓" | — |
| 미연결·연결 대기·pending | `--warn` 글자 + `--warn-soft` 배경 + "…" | `note` (accent 왼쪽 6px 선) |
| 철회·종료·만료 | `--muted` 글자 + `--bg` 배경 + "—" | — |
| 오류·위험 | `--danger` 글자 + `--danger-soft` 배경 + "!" | `problem` (danger 왼쪽 6px 선, `role="alert"`) |

색과 텍스트·기호를 함께 쓴다(색만으로 상태를 구분하지 않는다). 배지는 알약이 아니라 모서리 4px 사각형이다. 사실 상태가 없는 배지("보안 검증됨", "온라인", "활성 사용자")는 만들지 않는다.

### 8.6 포커스·모션

- 모든 링크·버튼·입력·`summary`에 `:focus-visible { outline: 3px solid var(--accent); outline-offset: 2px; }`. 위험 구역 안에서는 `--danger`.
- 모션은 사용하지 않는다. `details` 펼침도 즉시 전환이다. `prefers-reduced-motion`에서 전환이 있으면 끈다.

## 9. 반응형 기준

| 폭 | 레이아웃 |
|---|---|
| 1440px | 가운데 정렬 컨테이너 최대 1040px. 왼쪽 고정 섹션 이동 열(200px) + 본문. 에이전트 행은 한 줄(식별자 · 배지 · 키 수 · 펼침 표시). |
| 768px | 이동 열이 본문 위 가로 줄로 바뀐다. 본문 한 열. |
| 390px | 좌우 여백 16px. 영역 안쪽 여백 16px. 에이전트 행은 두 줄(식별자 / 배지·키 수). 주 버튼 전체 폭. 지문·키 ID는 줄바꿈하여 가로 스크롤 없음. |
| 320px | 위와 같음. 가로 스크롤 없음. |

수락: 1440·390에서 `document.documentElement.scrollWidth <= innerWidth`(스크롤바 폭 제외). 모든 상호작용 요소 높이 44px 이상. 200% 확대에서 가로 스크롤 없이 사용 가능.

## 10. 접근성 수락 기준

- A-01 `lang="ko"`, 페이지마다 고유 `title`. 헤딩은 `h1` 하나, 섹션은 `h2`, 에이전트 행은 `h3`로 순서를 건너뛰지 않는다.
- A-02 랜드마크: `header`, `nav`(섹션 이동), `main`, 섹션마다 `aria-labelledby`. 첫 요소는 "본문으로 건너뛰기" 링크이며 포커스 때만 보인다.
- A-03 키보드: 모든 동작을 Tab·Enter·Space로 수행한다. 탭 순서가 시각 순서와 같다. `details/summary`는 키보드로 열고 닫는다.
- A-04 대비: 8.1 표의 쌍을 지킨다. 텍스트 4.5:1, 비텍스트 3:1 이상.
- A-05 상태는 색 외에 텍스트·기호로 구분한다(8.5).
- A-06 폼: 모든 입력에 보이는 라벨. 오류는 `role="alert"`로 입력 위쪽에 표시하고 해당 입력과 `aria-describedby`로 연결한다. 이미 있는 `role="alert"`·`role="status"`는 유지한다.
- A-07 읽기 전용 값(식별자·지문·연결 수단)은 `readonly input` 대신 `code`로 쓰되 보조기기가 값 이름을 알 수 있도록 라벨 텍스트를 함께 둔다. `user-select:all`로 한 번에 선택된다.
- A-08 확대: 200%에서 내용 손실 없음. 320px에서 가로 스크롤 없음. 텍스트 간격 확대(줄 높이 1.5, 단어 간격 0.16em) 후에도 내용 손실 없음.
- A-09 위험 동작: 위험 구역의 `summary` 이름은 결과를 포함한다("이 에이전트 철회·키 철회"). 버튼 이름은 대상을 포함한다("이 에이전트 철회"). 같은 화면에 같은 이름의 버튼이 여러 개면 에이전트 식별자를 `aria-label`로 구분한다.
- A-10 인쇄·스크린 리더 전용 숨김 텍스트는 필요할 때만 쓰며 `display:none`으로 숨기지 않는다.

## 11. 보존해야 하는 보안·상태 불변식

화면 개선이 다음을 바꾸지 않는다. 바뀌면 이 문서의 위반이다.

- I-01 CSRF/요청 위조 방어: 쿠키 `Secure; HttpOnly; SameSite=Strict`와 gate의 `csrf` hidden 값, 상태 변경은 `POST`만. 모든 hidden 필드(`agent`, `kid`, `connection`, `generation`, `csrf`)와 폼 `action`·`method`를 유지한다.
- I-02 CSP: `default-src 'none'; style-src 'unsafe-inline'; form-action 'self'[+ accounts.google.com]; frame-ancestors 'none'; base-uri 'none'`. 스크립트·이미지·외부 리소스를 넣지 않는다. CSP를 완화하지 않는다.
- I-03 owner 분리: 홈·연결·관계·gate는 세션 회원의 owner 것만 보여 준다. 다른 owner의 에이전트·이메일을 어떤 화면에도 노출하지 않는다. 최근 로그인한 개인 계정과 노우 계정은 서로 다른 owner다.
- I-04 재인증: 5분 안 재확인이 필요한 동작(create·connect·confirm·device-confirm·key-revoke·agent-revoke·logout-all)의 서버 검사를 그대로 둔다. 화면의 사전 안내는 서버 거부를 대체하지 않는다.
- I-05 파괴 동작: 별도 확인 모달이나 새 확인 화면을 추가하지 않는다(기존 결정). 위험 구역 `details`와 결과 문장은 단일 `POST` 폼을 유지한다. 한 번 눌러 되돌릴 수 없는 동작이므로 위험 구역을 펼친 뒤에만 버튼에 도달한다.
- I-06 사실 상태: 실제 서버 상태 없는 "연결됨", "온라인", "보안 검증됨" 표시를 만들지 않는다.
- I-07 비밀값: 연결 수단(token)·키·지문·이메일을 로그·캡처·Git·채팅에 쓰지 않는다. 마스킹된 이메일 표시를 유지한다.
- I-08 기존 안내문의 약속 범위: 철회 기록 최소 24시간 보존, 영구 삭제·권한 복구 불가 약속 없음, 반복 초대 규칙은 그대로다.

## 12. DEV 개발 인계

DEV는 구현 방법과 파일별 기술 계획을 결정한다. 제품 요구와 완료 기준은 아래와 같다. 제품 규칙을 바꾸어야 하는 상황이 생기면 designer에게 질문한다.

### 12.1 요구 (우선순위 순)

| ID | 우선 | 요구 | 근거 |
|---|---|---|---|
| R-01 | 필수 | 8절 색·타입·간격·버튼·배지 토큰을 회원 웹 전 화면(로그인, 확인 대기, 홈, 연결 관리, 기기 승인, receipt, gate, 오류 화면)에 공통 적용한다. 한 곳에서 정의하고 재사용한다. | O-05, O-08, O-10 |
| R-02 | 필수 | 한국어 줄바꿈을 `keep-all` + `break-word`로 바꾸고 식별자·지문·코드에만 `anywhere`를 적용한다. | O-05 |
| R-03 | 필수 | 4.2 순서로 홈을 재구성한다: 헤더·확인 필요 띠·섹션 이동·에이전트·관계·승인·연결 확인·계정. | O-01, O-02 |
| R-04 | 필수 | 에이전트를 접힌 행으로 바꾸고(5.3) 키 철회·에이전트 철회를 위험 구역에 둔다. | O-03, O-04 |
| R-05 | 필수 | 연결 관리 화면에 3단계 진행 표시와 지문 강조를 적용한다(5.4, 5.5). | O-07 |
| R-06 | 필수 | 문구 표(6절)를 적용하고 내부 문구·"계정 비활성화 준비 중"을 제거한다. | O-09 |
| R-07 | 필수 | 오류·빈 상태 표(7절)의 다음 동작 링크를 모든 오류 화면에 보장한다. | — |
| R-08 | 필수 | 접근성 기준 A-01–A-10과 11절 불변식을 지킨다. | — |
| R-09 | 필수 | 로그인 화면에 서비스 설명 3줄을 추가한다(5.1). | O-11 |
| R-10 | 권장 | "확인 필요" 띠를 받은 초대·승인 요청 건수로 구현한다. 건수 계산이 기존 상태에서 즉시 가능하다. | — |
| R-11 | 권장 | 관계 목록을 받은 초대 → 보낸 초대 → 연결됨 → 종료 순으로 정렬한다. | — |
| R-12 | 후속 후보 | 다크 모드, 실제 접속 상태 표시(서버가 사실 데이터를 제공할 때만), 클라이언트 복사 버튼(CSP 변경 필요 — 이번 범위 아님). | — |

### 12.2 완료 기준 (DEV 자체 확인)

1. 홈·연결 관리·기기 승인·오류 화면의 모든 `form`의 `method`·`action`·hidden 필드가 변경 전과 같다(자동 비교 또는 테스트).
2. 기존 Go 테스트가 통과하고, 응답 헤더 CSP가 변경 전과 같다.
3. 다른 owner의 에이전트·이메일이 렌더되지 않음을 확인하는 기존 테스트가 유지된다.
4. 1440px·390px에서 `scrollWidth <= innerWidth`.
5. 새 외부 리소스(스크립트·이미지·폰트·스타일시트) 참조가 0건이다. 응답 HTML에서 `src=`·`<link`·`url(`이 없다.
6. 접힌 행·위험 구역이 JS 없이 키보드로 동작한다.

### 12.2.1 이전 UI 결함의 회귀 방지

이전 직접 검수의 결함이 되살아나지 않아야 한다. 후속 검수에서 다시 확인한다.
- F-UI-01(AGENTS-UI): 390px에서 키 지문·긴 식별자가 카드 안에서 줄바꿈되고 가로 overflow가 없다.
- F-UI-02: 모든 거절·오류 화면에 홈 복귀 또는 재확인·재선택 같은 다음 동작이 있다.
- F-UI-03: 만료·취소된 연결 화면에는 "연결 취소" 버튼이 없다.
- F-UI-04: 관계 기한은 raw UTC 문자열이 아니라 KST 표기다(`memberPair.Deadline`의 현재 형식 유지).
- F-UI-MSG-01(MESSAGES-UI): gate 거절 저장 뒤 결과와 홈 복귀가 있다. 이번 작업은 gate의 동작을 바꾸지 않는다.

### 12.2.2 디자인 정본과 검사

이 문서의 8절이 이번 회원 웹의 디자인 토큰 정본이다. 현재 UI는 Go 템플릿의 일반 CSS이며 별도 테마·디자인 lint는 구성되어 있지 않다(project.md). 새 색 값(8.1)은 기존 색과 같은 계열을 정리한 것이고 임의 값이 아니다. DEV는 토큰의 정의 위치(공유 CSS 상수)와 적용 가능한 검사 또는 미실행 사유를 정규 지시서에 기록한다. FullOps DESIGN 규칙은 Go 문자열 안의 CSS를 검사하지 못하므로 시각 판정을 대신하지 않는다.

### 12.3 제안 범위(DEV가 확정)

- 영향 가능 코드: `internal/relay/member.go`(`memberStyle`, `memberPages` 템플릿), `internal/relay/member_agents.go`(뷰 데이터), `internal/relay/http.go`(gate 템플릿의 스타일 공유 여부). 이 목록은 제안이며 구조 결정은 DEV 책임이다.
- 구조 선택(템플릿 분리, CSS 상수 공유 등)과 테스트 추가 여부는 DEV가 정한다. 새 패키지·라이브러리는 추가하지 않는다.
- 서버 상태·API·rate 정책 변경은 이번 범위가 아니다. 필요하면 별도 과제로 요청한다.

### 12.4 후속 검수 (designer)

DEV가 고정 SHA와 검수 환경을 제공하면 designer가 직접 확인한다:
- 1440×900·390×844에서 로그인 → 홈 → 에이전트 펼침·위험 구역 → 연결 발급 → 연결 관리 → 관계 → 로그아웃.
- 정보 위계, 기능 발견성, overflow, 한국어 줄바꿈, 포커스 표시, 대비, 오류·빈 상태.
- 실제 Google 로그인·실계정 데이터는 사용자 소유다. 합성 fixture 또는 사용자 허가 세션으로만 확인하고 값은 가린다.
이 검수 전에는 UI PASS를 기록하지 않는다.

## 13. 미확인과 한계

- 운영 `/`에는 Google 로그인만 보였다. 이메일 확인 코드 화면과 오류 화면은 템플릿 코드로만 확인했고 렌더 결과는 보지 않았다.
- 연결 관리·기기 승인·receipt·gate·refusal 화면은 이번에 직접 캡처하지 않았다. 문구와 상태 분기만 코드로 확인했다.
- 홈은 에이전트 4개, 관계·승인 요청이 적은 계정 한 개 기준이다. 에이전트 5개·긴 관계 목록의 길이는 관찰하지 않았다.
- 대비 비율은 계산값이며 실제 렌더 픽셀 측정이 아니다. 다크 모드는 확인하지 않았다.
- 한국어 "keep-all" 줄바꿈이 매우 긴 한글 연속 문자열에서 넘치지 않는지는 구현 후 390px·320px에서 확인한다.
- `frontend-design` 이름의 스킬은 이 환경에 없었다. 같은 목적으로 직접 관찰과 CSP·템플릿 제약 검토를 수행했다.
- 이번 결과는 디자인 인계 문서다. 최종 UI 수락·운영 배포·실제 Google/이메일 로그인 왕복 검증은 후속이다.

## 14. 구현 후 직접 시각 수락 기록

고정 후보 `c7828b227b5715f7bdf02bff46bf9cf102a758d9`를 합성 fixture 29개 화면, Orca 내장 브라우저 1440·390·320에서 직접 검수했다. 판정은 시각 PASS이며 critical·high는 0건이다. D04 9절 44px 기준을 글자 링크가 채우지 못한 medium 1건(V-01)과 low 11건이 남는다. 실제 Google 로그인·실계정·200% 확대·스크린 리더는 미검증이다. 근거와 한계: [시각 수락 보고](../../evaluations/qa-reports/SAR-SERVICE-UI-001-VISUAL/report.md).

후속(2026-10-11): V-01은 고정 후보 `295417acfdcd415ad0c768d164bd71dc07214ea9`에서 글자 링크 11종이 44px 이상임을 Orca 내장 브라우저 320·390·1440으로 직접 재측정했다. 시각 기준으로 resolved이며 위 기록과 low 발견은 보존한다. 키보드 포커스 재입력·실제 Google·스크린 리더·200%는 미검증이다. 근거: [재측정 보고](../../evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL/report.md).


후속(2026-10-12): SAR-FRONTEND-POLISH-001-VISUAL은 후보 `9597bad351a2fbfbaa726444b320b8cf7000e1e2`와 구현 `f7b4786`의 제품 동일성 및 before/after manifest 해시를 확인했다. after 합성 PNG 26장(홈 1440×900·390×844와 변경 화면 mobile)을 직접 열어 제공된 변경 화면을 시각 PASS로 판정했다. V-05·V-06·V-07·V-10 동의 화면·V-11 입력 크기·N-01은 제공된 사례에서 개선을 확인했다. V-04 320px·V-09 접근성 이름·포커스는 DEV 근거 재사용이며 새 직접 입력 검증이 아니다. V-02·V-03·V-08·V-12, mobile 생성 버튼 폭 N-02와 명령 줄바꿈은 비차단 잔여로 보존한다. 실제 로그인·운영·200%·스크린 리더·전체 상태의 최종 수락은 미검증이다. 기존 판정과 미확인 기록은 유지한다. 근거와 항목별 범위 판단: [이번 직접 시각 검수](../../evaluations/qa-reports/SAR-FRONTEND-POLISH-001-VISUAL/report.md).
