---
id: D02
title: 연동 단순화와 에이전트 한도 제거의 제품 규칙
status: review
updated: 2026-10-26
owner: designer
tasks: [SAR-CONNECT-SIMPLE-001]
upstream: [D01]
downstream: [D04, D10]
summary: "회원당 agent 수 제한 제거와 남용 방지 속도·용량 경계, 키 저장 전 시도 분리, 기본 이름, 재시도 폴더, 실패 사유, 연결 화면 규칙과 수락 기준 CS-AC-01–18"
---

# 연동 단순화와 에이전트 한도 제거의 제품 규칙 — SAR-CONNECT-SIMPLE-001

## 상태와 근거

이 문서는 SAR-CONNECT-SIMPLE-001의 요구사항(D02)이다. 화면 기준은 [D04 연결 단순화 절](../../design-docs/mockups/service-member-ui.md#연결-단순화와-내-에이전트-정리--sar-connect-simple-001)이다. 모듈 책임은 [D10 절](../../design-docs/module-design.md#sar-connect-simple-001-연결-단순화-책임)이다.

사용자 요청 원문(2026-10-26):

1. "일단 연결 갯수 한도를 풀어놔야겟다"
2. "그리고 기본적으로 연동하는게 너무 복잡해 전체적으로 손좀 봐라"
3. "전체적인 구조를 다시한번 살피고 최대한 UI/UX를 편하게 수정하고 연동 처리도 깔끔하게 되도록해"
4. "에이전트 연결한도도 없애고"

사용자 결정: 회원당 에이전트 연결 한도를 없앤다. 그 밖의 보호 장치는 이 문서가 정한다(L-01–L-06).

기준 커밋은 main `b413464adb86fba56488cb715d9bae58f247e96e`다. designer가 이 커밋에서 직접 확인한 코드 사실은 다음과 같다. 수정 방법은 DEV가 정한다.

| 사실 | 위치 |
|---|---|
| 새 agent 거절 한도는 전체 활성 200(`total`), 회원별 기록 10(`records`, 철회 포함), 회원당 활성 5(`active`)다. | `internal/relay/connections.go` `agentLimit`, `ownerAgentRecords` |
| 서버 상태는 한 행을 잠그고 읽고 쓴다. 기록 상한은 이 행이 끝없이 커지지 않게 하려는 기술 보호다. | `internal/relay/store.go` `transaction`, `connections.go` 주석 |
| 브라우저 동의(`device-confirm`) 때 agent를 만든다. 클라이언트가 키를 저장하기 전이다. 저장에 실패하면 키 없는 agent가 남는다. | `internal/relay/device.go` `approveDevice` |
| Dots 원격 연결은 이름 없는 agent를 만든다. | `internal/relay/remote_auth.go` |
| 연결 링크 화면은 사용된·만료된·취소된 연결을 모두 같은 401 일반 문구로 보인다. | `internal/relay/device.go` `devicePage` |
| Google 로그인 실패는 원인과 상관없이 같은 403 문구를 보인다. | `internal/relay/google.go` |
| 비공개 폴더 생성은 이미 있는 폴더를 거부한다. Windows에서는 새 폴더의 소유자를 다시 설정한다. | `adapters/src/private-files.ts` `privateDirectory`, `windowsACL` |

coor 실측(운영 relay, 2026-10-26, [지시서](../../../handovers/to_designer.md))은 다음 증상을 보고했다: 실패 시도가 이름 없는 agent 3–4개로 쌓였다. 그래서 세 번째 Claude Code 연결이 '연결 에이전트 초과'로 거절됐다. `knowslink_connect`는 실패 사유 없이 `{"state":"failed"}`만 돌려줬다.

이 문서가 바꾸는 기존 규칙:

- [SAR-PUBLIC-SERVICE](SAR-PUBLIC-SERVICE.md) 수용량 표의 '회원·agent' 행 중 '활성 agent 200개, 회원당 활성 agent 5개'
- [SAR-HUMAN-PEER](SAR-HUMAN-PEER.md) '용량과 남용 방어'의 '회원당 활성 agent 5개'와 'agent 200개'
- [SAR-SESSION-PEER-001](../../design-docs/SAR-SESSION-PEER-001.md)의 '회원당 활성 agent 5개'
- [SAR-UX-ISSUES-002](../../design-docs/SAR-UX-ISSUES-002.md) U8 표의 '미사용 agent만 남음: 자동 철회·삭제 없음'(키가 한 번도 저장되지 않은 시도에만 A-02로 바꾼다)

이 문서가 바꾸지 않는 규칙: 지문 비교와 명시적 연결 동의, 비공개 폴더 권한 검사, 발송마다의 승인(HP-10, SP-05·06), 친구 동의 범위(HP-03·06), agent당 활성 키 3개와 키 기록 20개, 친구·legacy pair 한도, 전체 활성 회원 100명, HTTP rate 한도, SAR-AUTOSETUP-001의 AS-01–AS-13.

## 현재 흐름과 복잡한 지점

사람 행동 수는 사용자가 직접 입력하거나 고르거나 실행하는 단위다. Google 로그인은 이미 로그인된 브라우저에서 생략될 수 있어 1로 센다.

| 호스트 | 현재 처음 연결의 사람 행동 | 수 |
|---|---|---|
| Grok Bot | 새 소스 폴더 clone, 고정 SHA checkout, `git status` 확인, installer 실행, Node·파일 확인, Command 등록 4칸 입력, `knowslink_status` 확인, 환경 변수 3개 편집, MCP 재시작, 연결 동의, 브라우저 Google·지문·이름·동의, `knowslink_connect_status` 확인 | 12 |
| Grok Bot 두 번째 | 계정 단위 플러그인의 환경 변수가 첫 Bot 폴더를 가리킨다. custom MCP를 손으로 하나 더 등록하고 위 등록·환경 단계를 반복한다. | 약 8 |
| Dots | custom MCP URL 등록, Google·동의, 도구 확인 | 3 |
| Claude Code | 운영에 자동 설정이 게시되지 않았다. 레포 빌드, `scripts/package_plugin.py --host claude --install-dir … --connect`, `claude --plugin-dir …`, 브라우저 동의, 세션 연결 | 7 이상 |
| Codex | Claude Code와 같다. | 7 이상 |

복잡한 지점:

| ID | 지점 | 사용자가 겪는 일 |
|---|---|---|
| C-01 | 회원당 활성 5개·기록 10개 한도 | 실패 시도가 자리를 차지해 정상 연결이 거절된다. |
| C-02 | 동의 시점의 agent 생성 | 키 저장에 실패한 시도가 '이름 없는 에이전트'로 목록에 쌓인다. |
| C-03 | 실패 사유 없음 | `failed`만 보여 다음 행동을 알 수 없다. |
| C-04 | 재시도 폴더 | 남은 시도 파일, 미리 만든 빈 폴더, 비시스템 드라이브(D:)의 권한 설정이 재시도를 막는다. |
| C-05 | Google 403 일반 문구 | 다른 계정을 고른 것인지 새 가입이 막힌 것인지 구별할 수 없다. |
| C-06 | 완료 뒤 링크 재방문 | 성공한 연결도 '만료·취소·사용되었거나 유효하지 않습니다'로 보여 실패로 오해한다. |
| C-07 | 이름 없는 에이전트 | Dots 연결과 일부 경로가 이름 없는 agent를 만든다. Grok 도구 결과에 표시 이름이 없다. |
| C-08 | Grok 설정 단계 | 환경 변수 3개 편집이 필수다. 두 번째 Bot은 수동 등록이 필요하다. |
| C-09 | 자동 설정 미게시 | 화면 안내와 실제 운영 설치 방법이 다르다. 게시는 사용자 승인 대기다(이 과제 범위 밖). |
| C-10 | 패키지 검사 실패 | `package_plugin.py --verify`의 `adapters/dist/mcp.test.js` 실패 경계 검사가 main `b413464`에서 실패한다(`/v1/registry` 뒤 `/v1/pull` 미호출). 원인 미확인. |

## 제품 규칙

### 한도 — L

| ID | 규칙 | 지켜야 할 조건 |
|---|---|---|
| L-01 회원당 수 제한 없음 | 회원당 활성 agent 수 상한을 없앤다. | 회원 화면·오류·도구 출력·문서에 '회원당 N개', 남은 개수, 카운터를 쓰지 않는다. |
| L-02 기록 상한으로 거절하지 않음 | 회원별 agent 기록 수(철회 포함)로 새 agent를 거절하지 않는다. | 철회 agent의 최소 24시간 보존과 이후 정리는 현행대로 둔다. |
| L-03 서버 전체 보호 | 고정 전체 상한 200을 없앤다. 서버 상태 크기 보호는 운영 용량 경계로 바꾼다. | DEV가 현재 상태 저장 방식에 맞춰 경고 임계와 거절 임계를 정하고 D12에 적는다. 거절 임계는 운영의 현재 활성 agent 수의 10배 이상이어야 한다. 경고 임계에 닿으면 운영 로그에 남긴다. 거절 임계에 닿으면 새 agent만 거절하고 L-03 문구(E 표 `service_capacity`)를 보인다. 회원의 사용량 탓으로 쓰지 않는다. |
| L-04 남용 방지 속도 | 새 agent 생성 속도만 제한한다. 회원당 24시간 이동 창에서 키 저장까지 끝난 새 agent 20개다. | 같은 창에서 나중에 철회한 agent도 센다. 키 저장 전 시도·취소·만료는 세지 않는다. 넘으면 E 표 `agent_rate` 문구와 다시 가능한 시각을 보인다. 기존 agent·키·친구는 영향받지 않는다. |
| L-05 그대로 두는 한도 | agent당 활성 키 3개, 키 기록 20개, 친구·legacy pair 한도, 전체 회원 100명, HTTP rate, 연결 시도 보관 2000개는 이번에 바꾸지 않는다. | DEV는 agent 수가 늘 때 legacy pair 전체 한도(active 400, pending 200)와 사람 친구 전달이 영향을 받는지 확인하고 보고한다. 영향이 있으면 같은 원칙(L-03)으로 designer에게 되돌린다. |
| L-06 화면 | 한도를 설명하는 안내·버튼을 없앤다. | '쓰지 않는 agent를 먼저 철회하세요'를 새 agent 거절 사유로 쓰지 않는다. 철회는 사용자가 원할 때만 쓰는 기능으로 남긴다. |

L-04 근거: 새 agent마다 Google 로그인, 브라우저 지문 비교와 동의, 키 소유 증명이 필요하다. 그래서 사람 한 명의 정상 사용은 하루에 호스트 수(현재 4종) 근처다. 20개는 정상 사용과 연결 시험 반복(coor 실측의 하루 약 7회)을 넘는다. 이 값은 한 행 상태가 한 회원 때문에 하루에 20개보다 빨리 커지지 않게 한다. 사용자가 이 값도 한도로 본다면 coor를 통해 값만 바꾼다.

### 연결 시도와 에이전트 — A

| ID | 규칙 | 지켜야 할 조건 |
|---|---|---|
| A-01 키 저장 뒤에만 에이전트 | 클라이언트의 키 저장이 서버에서 확인된 뒤에만 '내 에이전트' 목록의 agent로 보인다. | 브라우저 동의만 받은 시도는 '진행 중인 연결'에 따로 보인다. 친구에게도 보이지 않는다. 생성 시점을 옮길지 상태로 구분할지는 DEV가 정한다. |
| A-02 키 없는 시도의 자동 정리 | 키가 한 번도 저장되지 않은 시도는 연결 기한 뒤 24시간이 지나면 서버가 정리한다. | 정리 전에 클라이언트가 저장 완료를 보내면 정상 완료한다(F-07 복구 유지). 사용자는 정리 전에 '이 시도 정리'로 바로 지울 수 있다. 자격이 없으므로 철회 기록을 남기지 않는다. |
| A-03 키 있는 agent | 키가 저장된 agent는 자동 철회·삭제하지 않는다. | 현행 철회 절차와 문구를 유지한다. |
| A-04 기본 이름 | 모든 생성 경로(내 컴퓨터 연결, Dots 원격, 웹 '에이전트 추가')의 agent는 이름을 가진다. | 동의 화면 이름 칸은 기본 이름으로 채운다. 기본 이름은 `<호스트 표시 이름>`이다. 같은 이름의 활성 agent가 있으면 비어 있는 가장 작은 번호를 붙인다(예: `Claude Code 2`). 호스트를 모르면 `내 에이전트 <번호>`다. 사용자는 1–40자 규칙 안에서 고친다. 컴퓨터 이름·OS 사용자 이름·이메일을 기본 이름에 넣지 않는다(친구에게 보인다). |
| A-05 기존 이름 없는 agent | 이미 있는 이름 없는 agent는 목록에서 바로 이름을 정할 수 있다. | 카드 첫 줄에 '이름 정하기' 입력과 저장 버튼을 접지 않고 보인다. |
| A-06 도구 결과의 이름 | 모든 호스트 도구 결과(상태, 연결 결과, 보낸·받은 메시지)는 내 agent와 상대 agent의 표시 이름을 포함한다. | 이름이 없으면 `이름 없는 에이전트 (<구분값>)`이다. 키·credential·토큰은 포함하지 않는다. |

### 재시도와 폴더 — R

| ID | 규칙 | 지켜야 할 조건 |
|---|---|---|
| R-01 정리 없는 재시도 | 같은 연결을 다시 시도할 때 사용자는 폴더를 지우거나 새 경로를 입력하지 않는다. | `agent.json`이 있으면 연결된 agent로 재사용한다(AS-10). 연결 전 시도 파일만 있으면 재개할 수 있는 시도는 재개한다. 재개할 수 없으면 폴더를 `<폴더>.unused-<시각>`으로 옮기고 새로 시작한다. 지우지 않는다. 비어 있는 폴더는 비공개 권한을 설정·확인할 수 있으면 그대로 쓴다. |
| R-02 비공개 폴더 권한 | 현재 사용자가 쓸 수 있는 로컬 경로면 드라이브와 상관없이 연결이 성공한다. | 권한 결과는 현행과 같다: 소유자는 현재 사용자, 허용은 현재 사용자·SYSTEM·Administrators뿐, 링크 거부. 권한을 확인할 수 없으면 E 표 `private_folder` 사유·경로·다음 행동을 보인다. 이 조건을 완화하지 않는다. |
| R-03 기본 위치 | 사용자가 경로를 입력하지 않아도 된다. | Windows는 `%LOCALAPPDATA%\KnowsLink` 아래, Grok Bot은 설치 prefix의 `connections/` 아래에 호스트별 새 이름을 자동으로 정한다. 사용자가 경로를 지정하면 그 경로에도 R-01·R-02를 적용한다. |

### 실패 사유 — E

클라이언트(설치기, `knowslink_connect`, `knowslink_connect_status`)와 웹 화면은 실패마다 사유 코드, 화면 언어 문구, 다음 행동 하나를 보인다. `{"state":"failed"}`처럼 사유 없는 결과를 내지 않는다. 사유 코드는 아래 이름을 쓴다. 출력 형식은 DEV가 정한다. 키·토큰·credential은 출력하지 않는다.

| 사유 | 상황 | KO 문구 | EN 문구 | 바뀐 것 |
|---|---|---|---|---|
| `expired` | 10분 안에 브라우저 동의가 없음 | 연결 시간이 지났습니다. 다시 연결하면 새 링크를 만듭니다. | The connection timed out. Connect again to get a new link. | 없음 |
| `cancelled` | 브라우저에서 취소함(지문 불일치 포함) | 연결을 취소했습니다. 다시 연결하면 새 링크를 만듭니다. | You cancelled the connection. Connect again to get a new link. | 없음 |
| `unconfirmed` | 동의 뒤 키 저장 확인 실패 | 브라우저 동의는 끝났지만 키 저장을 확인하지 못했습니다. 다시 연결하면 먼저 이 시도를 확인합니다. | You agreed in the browser, but key storage wasn't confirmed. Connect again; it checks this attempt first. | 진행 중인 연결 1개(A-02로 정리) |
| `private_folder` | 비공개 폴더 권한을 설정·확인할 수 없음 | 키를 보관할 비공개 폴더를 준비하지 못했습니다(<경로>). 기본 위치로 다시 연결하세요. | Couldn't prepare a private folder for the key (<path>). Connect again using the default location. | 없음 |
| `unreachable` | 서버에 연결할 수 없음 | KnowsLink 서버에 연결할 수 없습니다. 인터넷 연결을 확인하고 다시 연결하세요. | Can't reach the KnowsLink server. Check your internet connection and connect again. | 없음 |
| `rate_limited` | 기존 HTTP rate 초과 | 요청이 많아 잠시 멈췄습니다. <시각> 뒤 다시 연결하세요. | Too many requests. Connect again after <time>. | 없음 |
| `agent_rate` | L-04 초과 | 오늘 새 에이전트를 많이 만들었습니다. <시각> 뒤 다시 연결하세요. 기존 에이전트는 그대로 쓸 수 있습니다. | You've made many new agents today. Connect again after <time>. Your existing agents still work. | 없음 |
| `service_capacity` | L-03 거절 임계 | 지금은 서비스 용량 때문에 새 에이전트를 만들 수 없습니다. 기존 에이전트는 그대로 쓸 수 있습니다. 나중에 다시 연결하세요. | The service can't add new agents right now. Your existing agents still work. Try again later. | 없음 |
| `signed_out` | 회원 세션 없음·재인증 필요 | Google로 다시 로그인한 뒤 연결하세요. | Sign in with Google again, then connect. | 없음 |

### 브라우저 화면 — W

| ID | 규칙 | 지켜야 할 조건 |
|---|---|---|
| W-01 링크 재방문 | 연결을 시작한 회원이 로그인한 상태로 링크를 다시 열면 그 연결의 실제 상태를 보인다. | 완료: 완료 화면과 agent 이름. 동의 완료·키 저장 대기: '앱이 키를 저장하는 중'. 만료·취소: 사유와 '시작한 앱에서 다시 연결'. 다른 회원·로그아웃 상태는 현행 일반 문구를 유지한다(상태·이름 비노출). |
| W-02 Google 계정 문제 구분 | 연결 중 Google 실패는 원인별 문구를 보인다. | 연결을 다른 계정이 이미 시작함, 새 가입이 닫힘, 로그인 기한 지남을 구분한다. 다른 계정의 이메일·이름·존재는 보이지 않는다. '다른 계정 선택' 링크는 Google 계정 선택 화면부터 다시 시작한다. 문구는 D04를 따른다. |
| W-03 동의 화면 순서 | 지문 비교, 이름(기본값), 연결 버튼 하나 순서다. | 대상·권한·기한은 한 줄 요약과 접힌 상세로 둔다. 연결 버튼 문구는 동의를 뜻해야 한다. 지문 비교와 명시적 버튼 누름은 생략하지 않는다. |
| W-04 동의 뒤 화면 | 동의 뒤 '앱으로 돌아가세요'와 '내 에이전트' 링크를 보인다. | 키 저장 확인 전에는 '연결됨'이라고 쓰지 않는다(AS-05). |

### 내 에이전트 화면 — M

| ID | 규칙 | 지켜야 할 조건 |
|---|---|---|
| M-01 순서 | 진행 중인 연결(있을 때만), 이름이 필요한 agent(있을 때만), agent 목록 순서다. | 빈 절은 제목도 보이지 않는다. |
| M-02 진행 중인 연결 | 시도마다 클라이언트 종류, 시작 시각, 상태, '이 시도 정리'를 보인다. | 상태는 지문 확인 대기, 동의 완료·키 저장 확인 필요, 만료다. 정리 버튼은 키가 없는 시도에만 있다. |
| M-03 한도 문구 삭제 | L-01·L-06의 문구와 거절 화면을 지운다. | 새 agent 거절은 E 표의 `agent_rate`·`service_capacity`만 남는다. |

### 호스트별 목표 흐름 — T

| ID | 호스트 | 목표 사람 행동 | 수 |
|---|---|---|---|
| T-01 | Claude Code·Codex(자동 설정 게시 뒤) | 명령 복사·실행, 설치 계획 동의, Google 로그인, 지문 비교 뒤 연결 누름(이름은 기본값), 호스트 다시 열기, 세션 연결 | 6 |
| T-02 | Claude Code·Codex 플러그인의 `knowslink_connect` | 도구 실행 동의, 링크 열기, Google 로그인, 지문 비교 뒤 연결 누름 | 4 |
| T-03 | Grok Bot | installer 실행, Command 등록 2칸(Command, Arguments), 연결 동의, Google 로그인, 지문 비교 뒤 연결 누름 | 5 |
| T-04 | Grok Bot 두 번째(다른 컴퓨터·작업 공간) | 연결 동의, Google 로그인, 지문 비교 뒤 연결 누름. 추가 MCP 등록·환경 변수 편집이 없다. | 3 |
| T-05 | Dots | MCP URL 등록, Google 로그인과 동의(이름은 기본값 `Dots`) | 2 |

T-03 조건: 배포 패키지(release)는 환경 변수가 없을 때 `public-node`, `https://link.knowslog.com`, R-03 자동 폴더를 기본값으로 쓴다. 이 기본값으로도 `knowslink_connect`의 명시적 동의(`confirmed`) 없이는 서버에 연결 시도를 만들지 않는다. 개발 빌드의 `held` 기본값은 DEV가 유지 여부를 정한다. 기존 환경 변수가 있으면 그 값을 우선한다(기존 설치 보존).

T-04 조건: 같은 컴퓨터에서 Bot 두 개를 구별할 신호가 없으면 두 Bot은 같은 agent를 쓴다. 이때 `knowslink_status`가 그 사실과 agent 이름을 보인다. 구별 신호의 존재는 DEV가 확인한다.

Grok Bot의 한 줄 설치기(clone·checkout 없는 설치)는 이번 범위가 아니다. 자동 설정 게시(사용자 승인 대기) 뒤 후속 과제로 PLANS.md에 대기시킨다.

## 수락 기준

CS-AC-01–CS-AC-17은 격리 합성 fixture에서 측정한다. 운영 relay·운영 agent·실제 상대 발송은 쓰지 않는다. CS-AC-18은 아래 실연결 재검증으로 측정한다.

| ID | 수락 기준 | 규칙 | 증명 방법 |
|---|---|---|---|
| CS-AC-01 | 한 회원이 활성 agent 7개를 만든다. 웹·API·도구 어디에서도 수 때문에 거절되지 않는다. 한도 문구가 화면에 없다. | L-01, L-06, M-03 | tester fixture 기록과 화면 문구 검색 |
| CS-AC-02 | 철회 agent 10개를 보존 중인 회원이 새 agent를 만든다. | L-02 | DEV 테스트, tester 재현 |
| CS-AC-03 | 회원 여러 명의 활성 agent 합이 201개가 되어도 새 agent가 만들어진다. 거절 임계 주입 때는 `service_capacity` 문구만 보인다. | L-03 | DEV 테스트(임계 주입) |
| CS-AC-04 | 24시간 안 21번째 새 agent가 `agent_rate`와 다시 가능한 시각으로 거절된다. 기존 agent·키·친구는 같다. 키 저장 전 시도 5회는 수에 들어가지 않는다. | L-04 | DEV 테스트(시계 주입) |
| CS-AC-05 | 취소·만료·동의 뒤 미저장 시도 5회 뒤에도 '내 에이전트' 수와 친구 화면의 agent 수가 같다. 다음 정상 연결이 성공한다. | A-01, C-01, C-02 | tester fixture |
| CS-AC-06 | 키 없는 시도는 기한 + 24시간 뒤 정리된다. '이 시도 정리'는 즉시 지운다. 정리 전 저장 완료를 보내면 정상 완료된다. 키 있는 agent는 그대로다. | A-02, A-03 | DEV 테스트, tester fixture |
| CS-AC-07 | 내 컴퓨터·Dots·웹 추가 경로 모두 이름 있는 agent를 만든다. 동의 화면 이름 칸에 기본 이름이 있다. 같은 호스트 두 번째는 번호가 붙는다. | A-04 | DEV 테스트, tester 화면 |
| CS-AC-08 | 이름 없는 기존 agent의 카드에서 접힌 영역을 열지 않고 이름을 저장한다. | A-05 | tester 화면 |
| CS-AC-09 | Grok·Claude Code·Codex 도구의 상태·연결·보낸·받은 결과에 양쪽 표시 이름이 있다. | A-06 | DEV 테스트(출력 snapshot) |
| CS-AC-10 | 남은 시도 파일 폴더, 미리 만든 빈 폴더, Windows 비시스템 드라이브 폴더에서 다시 연결하면 사용자 정리 없이 성공한다. 옮긴 폴더는 지워지지 않는다. 권한을 확인할 수 없는 경로는 `private_folder`를 보인다. | R-01, R-02 | DEV Windows 실측, tester 재현 |
| CS-AC-11 | E 표의 사유마다 클라이언트 출력에 사유 코드·문구·다음 행동이 있다. 사유 없는 `failed`가 0건이다. KO·EN 모두 확인한다. | E | DEV 테스트, tester 실패 주입 |
| CS-AC-12 | 완료한 연결 링크를 시작한 회원이 다시 열면 완료 화면이 보인다. 다른 회원·로그아웃 상태는 현행 일반 문구다. | W-01 | tester 화면 |
| CS-AC-13 | 다른 계정이 시작한 연결, 새 가입 닫힘, 기한 지남이 서로 다른 문구다. 다른 계정의 정보가 응답에 없다. | W-02 | DEV 테스트, tester 화면 |
| CS-AC-14 | release 패키지의 Grok 설정에서 환경 변수 없이 `knowslink_connect` 동의만으로 연결된다. 동의 전 서버 연결 시도가 0건이다. 기존 환경 변수 설치는 그대로 동작한다. | T-03 | DEV 테스트, tester 패키지 검사 |
| CS-AC-15 | `scripts/package_plugin.py --verify`가 통과한다. `mcp.test.js` 실패 원인이 기록된다. | C-10 | DEV 로그 |
| CS-AC-16 | 지문 비교, 명시적 동의, 비공개 폴더 권한, 발송 승인, 친구 권한의 기존 테스트가 통과한다. 키·토큰·credential이 출력·로그·화면에 없다. | 보안 경계 | DEV `make test`, tester 출력 검색 |
| CS-AC-17 | 회원 화면 변경(M, W)이 D04 문구·배치와 같다. 1440·390·320px에서 가로 넘침이 없다. 키보드로 모든 버튼에 닿는다. | M, W | designer 시각 검수(ART) |

실연결 재검증(사용자 결정 2026-10-26): 격리 fixture 수락과 운영 배포 뒤 서로 다른 두 회원으로 다시 연결한다. Grok Bot 쪽은 `knowslog01@gmail.com` 회원, Claude Code 쪽은 `scshin88@gmail.com` 회원이다. 친구 요청·수락 흐름을 포함한다. Google 로그인·동의 등 브라우저 단계는 에이전트가 ego 브라우저에서 모두 진행한다. 비밀번호·2단계 인증 입력과 실제 메시지 발송은 별도 승인이다. 절차는 [tester 지시서 초안](../../../handovers/pending/SAR-CONNECT-SIMPLE-001-tester.md)의 '실연결 재검증' 절이다.

| ID | 수락 기준 | 규칙 | 증명 방법 |
|---|---|---|---|
| CS-AC-18 | 배포된 relay에서 두 회원이 각각 Grok Bot·Claude Code를 연결하고 친구가 된다. 실패 시도가 목록에 남지 않는다. 서로의 agent가 기본 이름(또는 정한 이름)으로 보인다. 사람이 꼭 해야 한 행동이 기록된다. | L-01, A-01, A-04, A-06, W-01, T-02, T-03 | tester 2단계 기록 |

## 범위 밖

- 운영 배포와 SAR-AUTOSETUP-001 자동 설정 게시(사용자 승인 대기).
- D040 실세션 QA.
- 운영의 기존 agent·이름 없는 agent·테스트 agent(`agent_c18976751d5a96d8adfbb0` 등) 정리. A-02는 배포 뒤 새 규칙으로만 적용되며 키가 있는 기존 agent는 건드리지 않는다.
- Grok Bot 한 줄 설치기, agent당 키 한도 변경, 친구 한도 변경.

## 개정 이력

- 2026-10-26: SAR-CONNECT-SIMPLE-001에서 작성했다. 사용자 결정(회원당 한도 제거)과 designer 결정(L-02–L-06, A, R, E, W, M, T)을 담았다.
- 2026-10-26: 사용자 결정(coor 전달)으로 실연결 재검증 조건과 CS-AC-18을 더했다.
