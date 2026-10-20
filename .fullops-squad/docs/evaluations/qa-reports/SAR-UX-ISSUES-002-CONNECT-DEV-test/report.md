---
title: SAR-UX-ISSUES-002-CONNECT-DEV DEV 보고
status: review
updated: 2026-10-20
owner: dev
tasks: [SAR-UX-ISSUES-002-CONNECT-DEV]
summary: "회원 연결 경로 선택·재개(U8), 나머지 U7, O-NM-03, low 3건의 구현·검사·렌더·fixture 기동 근거와 미실행 항목"
---

# SAR-UX-ISSUES-002-CONNECT-DEV DEV 보고

DEV 검사 기록이다. 독립 리뷰·tester QA·designer 시각 수락이 아니다. 실제 호스트 설치·지원 완료 표시는 범위 밖이며 다섯 호스트를 지원 완료로 표시하지 않았다.

## 기술 계획 (실행함)

1. 경로 선택: `/home/agents?use=&start=` GET 폼. `use`는 web·grok·dots·claude·other, `start`는 first·continue만 받는다. 표시 전용이며 저장하지 않는다. 선택지 순서·문구는 `connect_guide.go:useChoices` 한 곳에 둔다.
2. 재개 상태: `agentResume`가 agent마다 저장된 상태만으로 단계와 다음 행동 하나를 계산한다. 근거는 최신 연결 시도(waiting·prepared·approved, 만료·취소 tombstone 24h), 활성 키, remote install의 미교환 code(1분)·token, device 승인 뒤 미저장이다. 불확실하면 ‘확인 필요’다. 온라인·앱 설치·도구 검색은 관측하지 않는다.
3. 세션 연결 관측: relay는 `/knowslink-session on`을 관측하지 않는다(SSE 구독은 세션 화면 성공의 근거가 아님, 설계 U8). 따라서 Claude Code 안내는 ‘앱에서 확인 필요’와 `/knowslink-session status`로 남긴다.
4. O-NM-03: 상태 JSON `Friend.CodeTail`(승인된 필드 추가). 요청 시점 trim 입력의 끝 4자, 4자 미만·`[A-Za-z0-9_-]` 밖 문자는 `?`, 이전 기록은 빈 값. 있는/없는/차단/자기 코드 모두 같은 함수로 저장한다. 요청자의 보낸 대기·거절·취소·만료 행에만 보인다. 전체 코드는 저장·URL·로그에 넣지 않는다.
5. 중복 클릭: 생성 폼의 기존 `Attempt` 값을 기존 `Idempotency` 맵에 `owner/agent-create:attempt → agent`로 기록한다. 새 필드·migration은 없다. agent 삭제 sweep이 해당 항목을 지운다.

## 변경 요약

- U8: 내 에이전트에 ‘어디에서 사용할까요?’(다섯 선택지, 보조 문구 포함)·‘처음 연결/기존 연결 계속’·고른 곳의 준비물·권한·다음 행동, ‘내 에이전트와 연결 상태’(agent별 단계·다음 행동·보낸 친구 요청 대기 수·요청 상태 조회 안내). 5번은 설치 명령·버튼·배지·도구 이름 없이 웹 사용으로 안내한다. OpenClaw·Meta Muse 이름과 ‘온라인’ 배지는 없다.
- U7: 연결 화면 consumed는 ‘키 등록됨’으로 바꾸고 단계 `aria-current=step`은 서버 상태의 한 단계만 둔다. 연결 화면·취소는 `/home/agents#agent-{id}`, confirm은 해당 연결 화면, 생성은 새 agent 위치로 복귀한다. 회원 연결 동작 거부는 내 에이전트(legacy 관계는 기존 연결 목록)로 복귀한다. 호출자 return URL은 받지 않는다. remote 앱 복귀 화면은 완료가 아님과 만료 시 앱에서 다시 시작을 안내한다. legacy 앵커 목록은 이미 ‘이 화면에서 찾기’였고 그대로 두었다.
- low 3건: 익명 조회 nav의 ‘내 에이전트’ 링크 제거, `gate_summary.go` comma-ok, 승인 목록 ‘요청 검토’ 접근성 이름에 받은 내 에이전트 추가.
- 문서: README·adapters 연결 절, D05·D10·D11·D13. D02/D04·D12는 바꾸지 않았다(D12는 새 운영 절차 없음).

## 검사

| 검사 | SHA | 결과 | 근거 |
|---|---|---|---|
| 원격 `make lint` | FINAL | EXIT=0 | FINAL-lint.log |
| 원격 `make test` | FINAL | EXIT=0 | FINAL-test.log |
| 원격 integration 전체(`-skip TestRemoteConsentBrowserReturn`) | FINAL | EXIT=0 | FINAL-integration.log |
| fixture 컴파일·실제 기동·probe | FINAL | 아래 | fixture 절 |
| 첫 실행 `TestConnectionResumeAndCodeTail` | bec551f | EXIT=1 | 429가 공용 admission에서 먼저 나 기대한 복귀 링크가 없었다(테스트 기대 문제). 403 2건은 첫 실패로 agent ID가 비어 생긴 연쇄 |
| 수정 뒤 `TestConnectionResumeAndCodeTail` | 137a7f5 | EXIT=0 | 하위 6건 PASS |

U8-A 자동 검사 `TestConnectionResumeAndCodeTail`(Postgres): 선택지 순서·보류 호스트 미표시·5번 설치 경로 없음·읽기/새로고침으로 agent 0; 429 뒤 agent 0과 복구 뒤 같은 attempt로 정확히 1개, 중복 제출은 같은 위치; local waiting→prepared→approved 구분과 재개 링크, 응답 유실 뒤 같은 token 재완료 401·키 1개 유지, 취소·만료(rotate 포함) 뒤 기존 키 불변, 만료된 시도와 미사용 구분, 서비스 재시작 뒤 같은 상태; device 승인 뒤 미저장 ‘확인 필요’, remote 거절은 agent 0·재사용 불가, 승인 뒤 code 대기·유실(확인 필요)·교환 구분; O-NM-03 있는/없는/차단/자기 코드 행의 동일 모양, 형식 다름, 이전 요청, 상대 코드 교체 뒤 불변, 상대 화면 비노출, 전체 코드 비노출, 취소 행 유지; 친구·기존 키 보존. 단위 `TestCodeTailKeepsOnlyPlainTails`, `TestConnectGuidePagesRender`.

## 렌더

`TestConnectGuidePagesRender`가 `KNOWSLINK_UI_FIXTURE_DIR`로 만든 합성 HTML을 headless Chrome(비모바일 viewport) 1440·390·320에서 [render/probe.mjs](render/probe.mjs)로 렌더했다. 대상 9개(선택 없음·다섯 선택지 안내·보낸 요청 끝 4자·키 등록됨 연결 화면·remote 앱 복귀), 27장 모두 페이지 가로 넘침 0([probe.json](render/probe.json)). 320px에서 종료 행의 ‘입력한 코드 …abcd’가 줄 중간에서 끊겨 `nowrap`으로 묶었다. DEV 확인이며 designer 시각 수락이 아니다.

## fixture

FIXTURE

## 하위 리뷰

읽기 전용 하위 리뷰 1건(Purpose: review, 고정 `bec551f`)을 배정했다. coor의 독립 고정 SHA 리뷰를 대신하지 않는다. 결과와 처리(`d4a1026`):

| 등급 | 지적 | 처리 |
|---|---|---|
| high | `connections_integration_test.go:317`의 거부 복귀 기대가 옛 ‘홈으로 돌아가기’ | 새 복귀 문구로 고침 |
| high | 받은 사람이 대기 요청을 차단하면 요청자 화면에 ‘해제’로 보여 코드 존재·행동이 드러남(기존 동작, O-NM-03 동일 모양 위반) | 차단된 대기 요청을 `blocked-pending`으로 저장하고 요청자에게 기한까지 대기·뒤 만료로 보인다. 요청자 취소는 허용, 같은 코드 반복은 대기 중 반복처럼 무동작. 받은 사람은 ‘차단’, 해제 뒤에는 ‘해제’로 보이며 수락할 수 없다. 통합 검사에 같은 모양·행 수 단언 추가 |
| medium | 키가 있는 agent의 추가·교체 시도 중 ‘키 등록 전’ 표시 | ‘기존 키 등록됨 · 새 키 …’로 표시 |
| medium | 생성 재전송 키와 메시지 멱등 키의 이론적 충돌 | 키를 NUL로 시작하는 `agentCreateKey`로 분리(메시지 키는 출력 가능 문자만) |
| low | 매 sweep의 Idempotency 전체 순회 | 삭제된 agent가 있을 때만 순회 |
| low | 철회된 agent로 재전송 복귀, 검사 공백 | 철회된 agent는 재전송 대상에서 제외. ‘온라인’ 단언은 배지 마크업 기준으로 유지(본문 설명 문장에 단어가 있음). Inbox 보존은 이번 변경이 메시지·Inbox 경로를 바꾸지 않아 기존 통합 검사(`TestHumanPostgres` 등) 통과로 대신한다 |

남은 한계: 차단된 요청이 기한을 넘긴 뒤 같은 코드로 다시 보내면 새 행이 생기고, 없는 코드는 같은 행이 다시 대기가 된다. 행 수 차이가 남으며 coor/designer 판단 대상이다.

## 실행하지 못한 것 (not_run)

- 실제 키보드·브라우저 200% 확대·스크린리더·보조 브라우저(U7-A), fixture 브라우저 왕복(tester 범위).
- 실제 호스트(Grok Bot·Dots·Claude Code·Codex·Hermes) 설치·인증·왕복. 범위 밖이다.
- `TestRemoteConsentBrowserReturn`(실제 Chromium 필요, 이전 단위와 같이 skip).
- 공용 admission 429 화면의 복귀 링크는 기존 `/home`이다. 이번 범위에서 바꾸지 않았다.
