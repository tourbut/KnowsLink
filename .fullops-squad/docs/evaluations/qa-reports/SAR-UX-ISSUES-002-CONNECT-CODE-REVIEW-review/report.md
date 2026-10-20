---
title: SAR-UX-ISSUES-002-CONNECT-CODE-REVIEW 리뷰
status: review
updated: 2026-10-20
owner: reviewer
tasks: [SAR-UX-ISSUES-002-CONNECT-CODE-REVIEW, SAR-UX-ISSUES-002-CONNECT-DEV]
summary: "8d8bad4cd75a 병합 전 독립 리뷰 보고서. critical/high 없음, medium 1건(차단된 대기 요청의 재요청 side-channel)과 low 3건."
---

# SAR-UX-ISSUES-002-CONNECT-CODE-REVIEW 리뷰

- 검토자 / CLI / 모델: reviewer 세션 `connect002-code-review-opus` / Claude Code / `claude-opus-5-5`. 실제 host 세션 ID: `d1124be7-5ae0-457b-a01c-8e2acc715b0b`. 구현자 세션: `ctx_81b0bfd36b71`.
- base SHA / head SHA / merge-base: `022c748e6fd153dea024785ebb045d28fee5a948` / `8d8bad4cd75a37546ba0f803e3c0c8290ddb6a87` / `022c748e6fd153dea024785ebb045d28fee5a948`.
- snapshot: `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-f8ba9e0ac0af4b3eaa6a027621fe664b` (detached, HEAD `8d8bad4`, `git status` clean, 수정하지 않음).
- OCR 버전 / 적용 규칙: open-code-review v1.13.1 (8eb8334) delegate 모드 / `.fullops-squad/review/rule.json` (rule_sha256 `f1061481…94de8`). 공통 규칙 `.fullops-squad/rules/common/README.md`, `coding-style.md`, `security.md`, `testing.md` (snapshot `8d8bad4` 기준).
- 요구사항·완료 기준 원천: `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md`의 U8(호스트 표, 선택지, 상태 표, U8-A), U7, U4 O-NM-03과 U4-A. DEV 지시·완료 기록 `.fullops-squad/handovers/logs/2026-10-20_to_dev.md`. DEV 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-DEV-test/report.md`.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 74 / 36 / 38 / 27 / 47. 지시서 범위(연결 UI/handler, remote 동의/복귀, README/adapters, D05/D10/D11/D13, O-NM-03)와 실제 diff가 일치한다. 범위 밖 제품 파일 변경은 없다.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: ERROR 0 / WARNING 7 (SIZE-001 6, SIZE-002 1) / 실행 불가 2 (`product-lint`, `product-test`: Windows 호스트에 make 없음, 원 status 보존). coor가 보완한 reason의 로그를 확인했다. 아래 '검증 및 남은 제약'을 본다.
- SIZE-002: 지시서는 U8+나머지 U7+O-NM-03+low 3을 한 단위로 배정했다. 실제 lint 대상 추가 894줄. 그중 테스트가 약 570줄(통합 330, 단위 97, fixture 124, 기타)이고 제품 코드는 약 300줄이다. 화면·handler·상태가 한 흐름이라 분할하지 않은 것은 타당하다.
- SIZE-001: `human.go`, `human_http.go`, `member_agents.go`, 두 테스트 파일이 상한을 넘었다. 기존부터 상한을 넘은 파일이 대부분이고 증가 폭이 작다. 새 기능은 `connect_guide.go`로 분리했다. 수락한다.
- DEP-001: 의존성 변경 없음(`go.mod`, `package.json` 변경 없음). `remote.py`의 paramiko는 저장소 밖 개인 도구 경로의 기존 DEV 도구이며 제품 의존성이 아니다.
- UI 디자인: 기존 `memberStyle`·`human-head`·`badge`·`kv`·`choice` 클래스를 그대로 쓴다. 새 색·토큰은 없다. DEV는 1440·390·320 렌더 27장에서 가로 넘침 0을 기록했다(`render/probe.json`). 디자인 lint와 테마 전환 검증은 이 저장소에 등록된 검사가 없어 해당 없음이다. designer 시각 수락과 U7-A 키보드·200%·스크린리더는 not_run이며 tester/designer 범위다.

## 검토 범위

result.json의 74개 `(path, status)` 모두에 상태와 사유를 기록했다.

- reviewed 27: 제품 Go 코드 10, Go 테스트 6, README 2, 운영·설계 문서 4, PLANS·board·contexts 3, DEV 보고와 실행 도구(`report.md`, `remote.py`, `render/probe.mjs`).
- skipped 47: prepare가 미리 채운 하네스 자동 생성 기록이다. `docs/evaluations/jev/` 7, `handovers/logs/` 1, 원시 로그 10, `raw-log-sha256.txt`, `render/probe.json`, PNG 27이다. 원시 로그는 직접 열어 `EXIT=` 줄과 대상 테스트 결과를 확인했다. manifest는 아래 low finding을 본다. PNG는 `probe.json`의 27개 항목(overflow false)으로 존재를 확인했고 시각 판단은 designer 범위다. 핸드오버 로그는 지시 정본으로 읽었다.

중점 항목별 확인 결과는 다음과 같다.

| 중점 | 결과 | 근거 |
|---|---|---|
| 상태 전이·재개가 서버 상태만 쓰는가 | 확인 | `connect_guide.go:agentResume`은 `Connections`·`Keys`·`Remote.Installs/Codes/Tokens`만 읽는다. `use`·`start` query는 안내 문구만 바꾸며 저장하지 않는다(`human_http.go:203-212`). 방문·details 열림을 기록하지 않는다. |
| 새로고침·중복 클릭·뒤로가기의 agent 중복 생성, 키 덮어쓰기 | 확인 | 폼 `attempt`(randomToken 43자)를 `Idempotency["#agent-create/owner/attempt"]`에 기록한다. 재제출은 같은 owner·미철회 agent로 303한다(`member_agents.go:242-258`). 길이 16–96 밖이면 멱등 없이 기존 동작이다. agentPattern `^[a-z]`라 `#` 키는 메시지 키(`from/key`)와 겹치지 않는다. agent 삭제 sweep이 해당 항목만 지운다(`connections.go:235-242`). 연결 재시도·취소·만료는 기존 키를 바꾸지 않으며 통합 검사가 단언한다. |
| 5분 재확인 창 | 확인 | `create`·`connect`·`confirm`·`device-confirm`·`key-revoke`·`agent-revoke`는 멱등 판정 전에 `reauthWindow`를 검사한다(`member_agents.go:204-208`). 변경 없음, 우회 경로 없음. |
| remote 복귀의 임의 return URL 금지 | 확인 | `memberBack`은 고정 두 경로만 반환한다. remote 거부 화면 Back은 `/home/agents` 고정이다. `remoteReturn`은 기존 검증된 redirect 값만 쓰며 로직 변경이 없다. confirm 복귀의 `connection` 값은 소유·존재를 확인한 맵 키다. |
| O-NM-03 CodeTail 요청자 한정·존재 비노출 | 대부분 확인, medium 1 | 저장은 모든 결과에서 같은 `codeTail` 함수다. 노출은 `!incoming`이고 pending/denied/cancelled/expired일 때만이다(`human_http.go:144-146`). 받은 행·친구 행·상대 화면에는 없다. 전체 코드는 저장·URL·로그에 없다. 단, 차단된 대기 요청의 취소/만료 뒤 재요청과 대기 상한에서 차단 추론이 가능하다(아래 F1). |
| 상태 JSON 필드 추가의 기존 데이터 호환 | 확인 | `Friend.CodeTail` 누락 시 빈 문자열이며 '이전 요청'으로 보인다. 새 state 값 `blocked-pending`은 Block 플래그와 함께 저장되어, 이전 코드로 되돌려도 수신자는 차단, 요청자는 해제로 읽는다. migration 없음. |
| 대상 외 호스트 미표시·가짜 명령 없음 | 확인 | 선택지 다섯 개의 순서·라벨·보조 문구가 설계 표와 글자 단위로 같다. 5번 안내에는 code·form·button·외부 링크가 없다. OpenClaw·Meta Muse 문자열은 제품 코드에 없다. Claude 안내의 `/knowslink-session on|status`와 `knowslink_text_receive`는 adapters 문서에 있는 명령이다. |
| README/adapters 안내 정확성 | 확인 | 화면 링크 앵커 `#grok-신규-설치`, `#grok-기존-설치-업데이트`, `#dots-신규-연결`, `#dots-기존-연결-업데이트`, `#claude-code-작업-세션-연결--sar-session-peer-001-dev`가 실제 제목과 맞는다. Grok 준비물(Linux x86_64/aarch64, Git, sh, curl, tar, gzip, sha256sum, Python 3), Node 22.22.2가 README와 같다. |
| HTML escape | 확인 | 모든 템플릿은 `html/template`이다. 사용자 값(이름, CodeLabel, Link)은 문맥 escape된다. CodeTail은 `[A-Za-z0-9_-]`만 그대로 보이고 그 밖은 '?'로 바꾼다. 알 수 없는 `use` 값은 안내를 그리지 않는다. |
| fixture가 테스트 태그 밖으로 새는가 | 확인 | `connect_browser_fixture_test.go`는 `_test.go`이고 `//go:build integration && connectfixture`다. `TEST_SYNTHETIC_DATABASE=1`과 manifest env가 없으면 실패한다. 허용 경로 밖은 403이고 manifest는 0600이다. |

## 발견 사항

### F1 — medium, 미해결: 차단된 대기 요청의 재요청 동작으로 차단을 추론할 수 있음

- 파일·줄: `internal/relay/human.go:170-210` (`requestFriend`), `human.go:299` (`decideFriend` cancel).
- 재현 조건 1(취소): A가 C의 코드로 요청한다. C가 받은 요청을 차단한다(`blocked-pending`, BlockC). A가 그 요청을 취소한다(`cancelled`, Block 유지). A가 같은 코드로 다시 요청한다. Block 플래그 때문에 placeholder 키로 새 행이 생기고 옛 행은 '취소'로 남는다. 차단이 없으면 같은 행이 다시 pending이 되어 행이 하나다.
- 재현 조건 2(만료): 위와 같되 24시간이 지난 뒤 재요청한다. 결과는 같다. DEV 보고는 이 경우만 '남은 한계'로 적었다.
- 재현 조건 3(대기 상한): 대기 상한 `mine >= 10`은 `State == "pending"`만 센다. 보낸 대기 10건 중 하나가 `blocked-pending`이면 11번째 요청이 새 행으로 생긴다. 차단이 없으면 생기지 않는다.
- 영향: U4 '차단: 상대에게 차단 여부를 공개하지 않는다', U8 상태 표 '상대 존재/차단을 누설하지 않음'에 어긋난다. 요청자가 행 수를 비교해야 하므로 즉시 노출은 아니다. 활성 친구 차단 뒤 재요청에도 같은 기존 패턴이 있다.
- 검증 근거: 코드 경로 추적. `TestConnectionResumeAndCodeTail/sent_request_tail_is_existence_blind`는 차단 직후 같은 코드 반복(무동작)만 단언하고 취소·만료·상한 경로는 단언하지 않는다.
- 제안: 요청자가 같은 코드를 다시 보낼 때 대상 행이 상대 차단 행이면 placeholder로 옮기지 않고 같은 행을 `blocked-pending`으로 갱신한다(세대·기한·CodeTail). `mine` 계산에 요청자의 `blocked-pending` 행을 포함한다. 각 경로에 행 수 단언을 추가한다. 수용하려면 coor/기획 결정으로 기록한다.

### F2 — low, 미해결: 요약을 확인할 수 없는 승인 요청에 승인 버튼이 남음

- 파일·줄: `internal/relay/gate_summary.go:46-74`, 화면 `human_pages.go:56`, 호출 `http.go:575`.
- 내용: `gateSummary`가 모양 불일치로 nil을 반환하면 화면은 '승인할 내용을 확인할 수 없습니다'를 보인다. 그러나 `{{if .Active}}` 승인 버튼은 요약과 무관하게 남는다. 예전에는 같은 상황에서 panic(500)이었다. Parse가 모양을 보장하므로 현재 도달 가능성은 낮다.
- 제안: 요약이 nil이면 `Active`를 false로 둔다.

### F3 — low, 미해결: 키가 있는 agent의 새 시도 상태 문구 중복

- 파일·줄: `internal/relay/connect_guide.go:49-51`.
- 내용: 접두를 붙이면 '기존 키 등록됨 · 새 키 키 등록 전 · 지문 확인 필요'처럼 '키'가 겹친다. approved 분기도 '새 키 브라우저 동의 완료 …'로 어색하다. 의미는 맞으므로 low다. designer 문구 검수에서 함께 다듬는다.

### F4 — low, 미해결: 원시 로그 manifest의 hash 기준이 커밋 바이트와 다름

- 파일: `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-DEV-test/raw-log-sha256.txt`.
- 내용: 10개 hash가 커밋 blob(LF, `.gitattributes` eol=lf)과 맞지 않는다. 리뷰어가 `git show 8d8bad4:<path> | sed 's/$/\r/' | sha256sum`로 CRLF 변환 뒤 비교해 10개 모두 일치를 확인했다. Windows 작업 사본에서 계산한 값이라 어떤 checkout에서도 `sha256sum -c`가 실패한다.
- 제안: 다음부터 커밋 blob 기준으로 계산하거나 계산 방식을 manifest에 적는다. 로그 내용의 신뢰에는 영향이 없다.

DEV 하위 리뷰의 high 2·medium 2·low 2 처리(`d4a1026`, `37febf9`)는 코드에서 반영을 확인했다.

## 검증 및 남은 제약

리뷰어는 lint·테스트를 다시 실행하지 않았다. snapshot과 커밋된 증거만 읽었다.

- `lint.json`: base `022c748`, head `8d8bad4`, ERROR 0. `product-lint`·`product-test`는 `unavailable`(Windows에 make 없음)이며 원 status를 보존했다. coor가 reason에 Linux 실행 로그를 보완했다.
- 보완 로그 확인: 제품 소스 `37febf9`와 `8d8bad4`의 차이는 `.fullops-squad/` 아래 16개 파일뿐이다(`git diff --name-only 37febf9..8d8bad4`에 제품 경로 없음). 따라서 `37febf9` 실행 결과를 `8d8bad4` 제품 코드에 적용할 수 있다.
  - `37febf9…-lint.log`: `make lint`, 39줄, `EXIT=0`. 로그의 'ERROR: candidate exists' 줄은 `test_public_config.sh`의 의도된 실패 경계 출력이며 종료코드 0이다.
  - `37febf9…-test.log`: `make test`(`go test -race` 포함), `EXIT=0`.
  - `37febf9…-integration.log`: `-race -tags=integration` 전체(`-skip TestRemoteConsentBrowserReturn`), `ok … internal/relay 131.908s`, `EXIT=0`. `TestConnectionResumeAndCodeTail` 하위 6건 PASS.
  - fixture compile·start·probe 로그 `EXIT=0`. probe에서 manifest 권한 600, 범위 밖 `/v1/pull` 403.
  - `4d0fe06` integration `EXIT=1`과 `bec551f` 첫 실행 실패는 DEV 보고의 원인 설명(NUL 접두사 JSONB 거부, 테스트 기대 문제)과 일치하며 뒤 SHA에서 고쳐졌다.
- 환경 동일성: 실행은 task 소유 격리 clone `/tmp/knowslink-connect002-37febf953bcb-*`와 `postgres:17-alpine@sha256:b0f9…` 컨테이너다(`remote.py`). 자동 판정은 하지 않았고 로그 첫 줄·HEAD 검증 명령으로 확인했다.
- 테스트 감소·완화: 기대값 변경 3곳은 복귀 문구 변경을 반영한 것이며 단언을 약하게 하지 않았다. 새 단언(aria-current 개수, '연결 완료' 부재)이 추가됐다.
- 미실행: `TestRemoteConsentBrowserReturn`(실제 Chromium 필요), U7-A 실제 키보드·200%·스크린리더·보조 브라우저, 실제 호스트 설치·왕복. 모두 DEV 보고에 not_run으로 기록됐고 tester/designer/별도 승인 범위다.

## 대화 미참조 인계 점검

이전 대화 없이 PLANS.md와 설계 정본에서 시작했다.

| 확인 항목 | 정본 경로/절 | 결과(확인/미확인/해당 없음) | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `design-docs/SAR-UX-ISSUES-002.md` U4·U7·U8, O-NM-03 결정 기록 링크 | 확인 | 없음 |
| 구조와 구현/미완료 상태 | `module-design.md`(connect_guide·member_agents·human), `interface-design.md`(GET /home/agents, attempt, 복귀) | 확인 | F1 미해결 상태는 PLANS에 남겨야 함 |
| 실행·검증 방법과 증거 | DEV 보고 '검사'·'fixture' 절, `remote.py` | 확인 | F4 manifest 기준 |
| 운영·복구(해당 시) | `operations/transition.md` DEV 후보 인계 절 | 확인 | migration 없음, 기존 relay 교체 절차 |
| 다음 작업·담당·재개 조건 | `PLANS.md` 4단위 항목, transition 절 '남은 일' | 확인 | tester QA·designer 시각 검수·main 통합 대기. 5단위 U3 |
| 로컬 링크·절 접근/지원 한계 | README·adapters 앵커, DEV 보고 상대 링크 | 확인 | 깨진 링크 없음 |
| snapshot 정리 후 정본 접근 | 이 폴더(coor 기록 체크아웃)와 `8d8bad4` 커밋 증거 | 확인 | snapshot 경로에만 있는 정본 없음 |

## 검토 결론

수락 가능하다. 단, medium F1의 처리를 병합 전에 결정해야 한다. critical/high는 없다. 연결 재개는 저장 상태만 쓴다. 생성 멱등, 5분 재확인 순서, 고정 복귀 경로, CodeTail의 요청자 한정 노출, 상태 JSON 호환, 대상 외 호스트 미표시, 문서 정확성, escape, fixture 격리는 요구와 일치한다. F1(차단된 대기 요청의 취소/만료 뒤 재요청과 대기 상한에서 차단 추론)은 U4 차단 비공개 규칙에 어긋난다. 수정 과제로 보내거나 coor/기획의 수용 결정을 기록한다. F2–F4는 후속 low다. 이 결론은 AI 검토 판단이며 OCR 자동 판정이 아니다. `review.py check` 통과는 기록 형식 검사이며 테스트 성공을 보증하지 않는다.
