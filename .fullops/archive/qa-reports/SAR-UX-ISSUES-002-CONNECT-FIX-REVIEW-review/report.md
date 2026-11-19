---
title: SAR-UX-ISSUES-002-CONNECT-FIX-REVIEW 리뷰
status: review
updated: 2026-10-20
owner: reviewer
summary: "abc219120509 병합 전 독립 리뷰 보고서. 4단위 2차 수정의 재요청 행·Held·차단 해제·D-CT-01을 검토했고 critical/high 0, 새 low 3건이다"
---

# SAR-UX-ISSUES-002-CONNECT-FIX-REVIEW 리뷰

- 검토자 / CLI / 모델: connect002-fix-review-opus / Claude Code(host session `e832cc01-d97d-445a-a955-1870eb6c8d3f`) / claude-opus-5-5. 구현자 세션 `ctx_8aabe4b5e6bf`와 다른 세션이다.
- base SHA / head SHA / merge-base: `8d8bad4cd75a37546ba0f803e3c0c8290ddb6a87` / `abc219120509027b143b0e4ad121511a5a18395b` / `8d8bad4cd75a37546ba0f803e3c0c8290ddb6a87`. 제품 코드는 `37febf9`→`d6369a4`이고 `d6369a4`→`abc2191`은 테스트 주석 한 줄과 기록 문서다.
- snapshot: `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-ddaab037916e4e77adccd13f5a7cf91b`(detached, clean, 수정 없음). 결과·보고서는 coor 기록 체크아웃에만 썼다.
- OCR 버전 / 적용 규칙: open-code-review v1.13.1(8eb8334) delegate 모드, `.fullops-squad/review/rule.json`(rule_sha256 `f1061481…`). 공통 규칙 `fullops-common-0.3.3`(`.fullops-squad/rules/common/` README·coding-style·testing·security, 마지막 변경 `fbbe5a6`, 이번 delta에서 변경 없음).
- 요구사항·완료 기준 원천: snapshot의 [SAR-UX-ISSUES-002.md U4](../../../design-docs/SAR-UX-ISSUES-002.md) 2026-10-20 CONNECT-VISUAL 재요청 결정과 완료 조건 U4-A, [designer CONNECT-VISUAL 판정](../SAR-UX-ISSUES-002-CONNECT-VISUAL/report.md)(V-CN-01·02), DEV 지시서 로그 `handovers/logs/2026-10-20_to_dev.md` 297–323행(F1–F4, D-CT-01, O-CT-01~04 출처: tester CONNECT-TESTER 보고), [DEV 보고 2차 수정 절](../SAR-UX-ISSUES-002-CONNECT-DEV-test/report.md), 이전 리뷰 [CONNECT-CODE-REVIEW](../SAR-UX-ISSUES-002-CONNECT-CODE-REVIEW-review/report.md)(F1–F4).
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 166 / 50 / 116 / 23 / 143. skipped 143은 prepare가 하네스 자동 생성 기록(`docs/evaluations/jev/`, `handovers/logs/`, `qa-reports/` 원시 로그·JSON)으로 미리 채운 항목이다. 범위 이탈은 없다. 소유 파일은 `internal/relay/` 친구·연결·gate·fixture와 기술 문서·DEV 보고다.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: 0 / 7(SIZE-001) / 2. coor가 snapshot에서 delta lint를 실행했다. `product-lint`·`product-test`는 Windows 호스트에 make가 없어 `unavailable`이며, reason이 Linux 격리 clone의 `d6369a4` `make lint`·`make test` EXIT=0 로그를 가리킨다. 리뷰어는 재실행하지 않았다.
- SIZE-002: 지시서 예상 규모는 F1·D-CT-01·low 5건 수정이다. 제품 코드 delta는 11파일 +308/−36이며 그중 새 테스트 파일이 190줄이다. 분할하지 않은 이유(기존 함수 안 분기, 1차 attempt와 같은 과제)는 DEV 보고 ‘FullOps lint (2차)’에 있고 타당하다.
- SIZE-001: `human.go` 341·`human_http.go` 504·`member_agents.go` 409·`http.go` 690과 테스트 2개, `PLANS.md`. 모두 이전부터 상한을 넘은 파일에 몇 줄을 더했다. 이번 리뷰의 수정 요청 대상이 아니다.
- DEP-001: `go.mod`·`go.sum`·`package.json` 변경 없음. 표준 라이브러리(`strings`, `net/http/httptest`, `reflect`)만 추가로 import한다.
- UI 디자인: 변경 화면은 연결 화면 Unsure 분기, legacy 목차 제목, 연결 거부 복귀, 재확인 문구다. 기존 `.note`·`.toc`와 `--muted` 토큰만 쓴다. DEV가 1440·390·320 렌더 9장에서 가로 넘침 0을 기록했다([render-fix/probe.json](../SAR-UX-ISSUES-002-CONNECT-DEV-test/render-fix/probe.json)). designer 시각 재검수는 V-CN-01 판정상 필요 없고, 새 문구 화면의 designer 수락은 아직 없다. 테마 전환은 새 색을 추가하지 않아 해당 없음이다.

## 검토 범위

result.json의 23개 reviewed 파일은 diff와 주변 코드를 직접 읽었다. 제품 코드는 `requestFriend`·`decideFriend`·`humanView`·`friendByID`·친구 권한 경로(`human.go` 118–135, `friendsNow`)·친구 템플릿(`human_pages.go` 19–23)·`connectionPage`·`approveDevice`·`sweepConnections`·`gateSummary`까지 따라갔다.

skipped 143건은 다음 방법으로 존재와 무결성을 확인했다.

- DEV 원시 로그: `raw-log-sha256.txt` 끝 14줄을 `git show abc2191:<path> | sha256sum`과 비교해 14개 모두 일치했다(F4 해소).
- `d6369a4` 로그 내용: lint EXIT=0, test `ok …internal/relay`·EXIT=0, integration `TestFriendRerequestReusesOneRow`·`TestFriendBlockedRequestCountsTowardSentLimit`·`TestFriendLegacyPlaceholderAndStaleUnblock` 등 PASS, `ok … 147.174s`·EXIT=0, fixture 재요청 왕복 `FAILURES 0`·EXIT=0.
- jev packet·route·handover 로그·이전 리뷰 폴더는 하네스 기록이다. 판단에 쓴 부분(DEV 지시서 297–323행, 이전 리뷰 F1–F4)만 직접 읽었다.

## 이전 리뷰 F1–F4와 연계 지적의 처리

| ID | 처리 확인 |
|---|---|
| F1 = V-CN-01 = O-CT-01 (medium) | 해소. 상대 차단·수신 한도는 대리 행 대신 쌍 행을 `blocked-pending`으로 재사용한다(`human.go:184`·`232–237`). 대리 행은 없는 코드·자기 코드·요청자 본인 차단에만 쓴다. 다섯 경우의 요청자 행 수·순서(ID 정렬, ID 유지)·상태·끝 4자·새 24h 기한·세대 +1이 같다. 발신 상한 `mine`과 전체 `pending`은 기한 전 `blocked-pending`을 센다(`open`). 20000 보호 상한은 첫 요청에만 걸리고 재사용 행에는 걸리지 않아 다섯 경우가 같다. |
| F2 (low) | 해소. `view["Active"] = active && summary != nil`. `gateSummary`는 실패 때 nil을 반환한다. 서버 POST 재검사는 바뀌지 않았다(현재 도달 경로 없음). |
| F3 (low) | 해소. ‘기존 키 등록됨 · 새 키 등록 전 · {단계}’로 ‘키 키’ 중복이 없다. |
| F4 (low) | 해소. 위 hash 비교. |
| D-CT-01 (medium) | 해소. 연결 화면은 `Device && expired && agent 존재 && 활성 키 0`일 때 카드와 같은 `deviceUnsure` 문구·링크만 보이고 ‘새로 시작하세요’와 단계 목록이 없다. agent는 `approveDevice`에서만 만들어지므로 승인 전 만료는 기존 만료 문구다. 카드 조건과 모순이 없다. |
| V-CN-02, O-CT-02, O-CT-04 (low) | 해소. 거부 복귀 `/home/agents`, 목차 `aria-labelledby` 제목, 재확인 문구. |

## 중점 확인 결과

- 존재·차단 추론(요청자 쪽): 차단·수신 한도·legacy 대리 행 경로 모두 요청자 행 하나가 같은 모양으로 바뀐다. 대기 중 반복은 모든 경우 무동작이다. 본인 차단 상대의 대리 행 경로는 요청자가 이미 아는 정보다. F-R3의 기존 데이터 예외만 남는다.
- `Friend.Held`와 수신자 화면: `-`이면 받은 사람에게 행이 없고, 그 외에는 이전 상태를 보인다. 취소·만료·반복에도 바뀌지 않는다. 받은 사람이 예전 발신자였던 행은 방향이 바뀐다(F-R1).
- `decideFriend` unblock 세대 생략: 조건은 `f == nil || (f.Generation != generation && decision != "unblock")`이다. 다른 결정(수락·거절·취소·해제·차단)은 여전히 세대가 맞아야 한다. `friendByID`가 참여자만 찾고, unblock은 호출자 자신의 플래그(`id == f.A`이면 `BlockA`, 아니면 `BlockB`)만 확인한다. 상대가 건 차단이나 다른 회원의 행은 풀 수 없다. 테스트가 요청자 쪽 unblock 거부를 단언한다. 차단은 원래 세대를 올리지 않았으므로 stale 폼 재전송 위험은 이전과 같다.
- 차단 행 Generation 0: 차단자 화면의 차단 행은 비대기 템플릿이라 이름·기한·방향을 렌더하지 않는다. 세대 0만 해제 폼에 실린다. 차단 행은 `Friends` 목록에 들어가지 않는다. 친구 상세(`/home/friends/{id}`)는 차단·비활성 행을 거부한다. 따라서 차단된 사람의 재요청이 차단자 화면에 드러나지 않는다.
- 권한: `blocked-pending`은 `active`가 아니고 차단 플래그나 `SuppressLegacy` 검사가 그대로라 전달 권한을 만들지 않는다. 수락은 `pending`만 받으므로 보류 행을 수락할 수 없다.
- 기존 저장 데이터 호환: `Held`는 `omitempty`이며 없으면 양쪽 표시다. migration은 없다. 끝난 대리 행은 쌍 키로 옮겨 ID·세대를 유지한다(테스트 있음). 되돌리기 한계는 data-model·transition에 기록돼 있다.
- gate Active: 위 F2.

## 발견 사항

### F-R1 — low, 미해결: 수신 한도 보류 행이 받은 사람에게 방향이 바뀌어 보임

- 파일·줄: `internal/relay/human.go:214-237`, `internal/relay/human_http.go:134`, `:152`.
- 재현 조건: T가 A에게 보낸 요청이 거절·취소·만료로 끝났다. T의 수신 대기가 10건 이상이다. A가 T의 코드로 요청한다. 행은 `blocked-pending`, `Held`=이전 상태가 되지만 `Inviter`가 A로 바뀐다. T 화면에서 이 행이 incoming이 되어 ‘입력한 코드 …abcd’(또는 ‘이전 요청’) 라벨이 사라진다.
- 영향: 정본 interface-design·data-model의 ‘수신 한도 보류 행은 받은 사람에게 이전 모습 그대로’와 다르다. T는 A가 보류된 재요청을 보냈음을 추론할 수 있다. 요청자 쪽 존재·차단 추론은 아니다. DEV 보고 ‘남은 한계’에만 적혀 있고 테스트 `receive_limit`은 A가 처음 발신자인 경우만 본다.
- 제안: `Held`에 이전 방향(이전 Inviter·CodeTail)을 함께 보존하거나, 이 예외를 interface-design·data-model에 기록한다.

### F-R2 — low, 미해결: 보류 행을 받은 사람이 직접 차단 POST하면 요청자 화면이 ‘해제’로 바뀜

- 파일·줄: `internal/relay/human.go:332-340`.
- 재현 조건: F-R1처럼 받은 사람 T에게 이전 기록(`Held`가 `-` 아님)으로 보이는 보류 행이 있다. T가 화면의 행 ID와 추측한 세대로 `decision=block`을 직접 POST한다. `pending` 판정이 `State == "pending"`만 보므로 행이 `removed`+차단이 된다. 요청자 A의 살아 있는 대기 요청이 ‘해제’로 보인다. 보이는 대기 요청을 차단하면 `blocked-pending`으로 대기가 유지되므로 결과가 다르다.
- 영향: UI 버튼이 없고 세대 추측이 필요한 조작 경로다. 차단 비공개 규칙의 예외다.
- 제안: block에서 상대가 Inviter인 기한 전 `blocked-pending`도 `pending`처럼 `blocked-pending`으로 남긴다.

### F-R3 — low, 미해결: 쌍 행과 열린 대리 행이 함께 있는 기존 데이터의 재요청

- 파일·줄: `internal/relay/human.go:184-193`.
- 재현 조건: 이전 코드의 차단 경로가 쌍 행(종료, 상대 차단)과 대리 행(대기)을 남긴 상태에서 A가 같은 코드를 다시 보낸다. `old`가 쌍 행이라 `legacy` 검사가 꺼지고 쌍 행이 `blocked-pending`으로 살아난다. 같은 코드의 열린 요청이 2개가 되고 발신 상한을 하나 더 쓴다. 없는 코드는 같은 상황에서 무동작이다.
- 영향: 배포 시점에 대기 중이던 기존 대리 행(최대 24시간)에만 해당한다. data-model은 ‘합치지 않는다’만 적고 이 경우를 다루지 않는다.
- 제안: `hidden`이고 같은 코드의 대리 행이 열려 있으면 무동작으로 돌려주거나 정본에 한계로 적는다.

## 검증 및 남은 제약

- 리뷰어는 lint.py·테스트를 재실행하지 않았다(리뷰 규칙). 재사용 증거는 DEV 보고의 `d6369a4` Linux 격리 clone 로그(lint·test·integration·fixture compile·기동·재요청 왕복, 모두 EXIT=0)와 coor의 snapshot delta `lint.json`(ERROR 0)이다. 로그 바이트는 커밋 blob hash로 확인했다. 실행 환경 동일성(같은 원격 호스트·Go·Postgres 컨테이너)은 DEV 보고 서술로만 확인했다.
- 테스트 변경: `human_test.go`의 차단 재요청 단언은 새 규칙(행 하나·`blocked-pending`)으로 바뀌었고 이전보다 강하다. 세대 재취득은 재요청이 세대를 올리는 새 동작 때문이다. 기대값 완화가 아니다.
- tester 재확인(fixture 재기동 필요)과 designer의 새 문구 화면 수락은 아직 없다. fixture 왕복이 A의 세 행을 이미 대기로 바꿨다는 DEV 기록을 확인했다.
- 남은 일: F-R1–F-R3 후속 또는 정본 예외 기록, tester D-CT-01·재요청 재확인, main 통합.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과(확인/미확인/해당 없음) | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `design-docs/SAR-UX-ISSUES-002.md` U4 ‘2026-10-20 CONNECT-VISUAL 결정’, CONNECT-VISUAL report V-CN-01 | 확인 | 수신 한도 보류의 방향 예외(F-R1)가 정본에 없다 |
| 구조와 구현/미완료 상태 | `design-docs/module-design.md` human.go 행, `data-model.md` 친구 문단, DEV 보고 ‘남은 한계 (2차)’ | 확인 | F-R3 경우가 data-model에 없다 |
| 실행·검증 방법과 증거 | DEV 보고 ‘검사 (2차)’·‘fixture 추가 경로 (2차)’, `raw-log-sha256.txt` | 확인 | 없음 |
| 운영·복구(해당 시) | `operations/transition.md` ‘2차 수정’ 문단 | 확인 | 되돌리기 시 보류 행 노출이 기록돼 있다 |
| 다음 작업·담당·재개 조건 | `PLANS.md` 4단위 줄, `board/board.json` note, transition ‘남은 일’ | 확인 | board note가 ‘F1 수정 예정’ 시점 문구다. 통합 때 갱신 필요 |
| 로컬 링크·절 접근/지원 한계 | DEV 보고·VISUAL report 상대 링크 | 확인 | tester 보고는 tester 워크트리 절대경로라 main에서 열리지 않는다 |
| snapshot 정리 후 정본 접근 | 이 폴더의 result.json·report.md·lint.json, coor 체크아웃 | 확인 | 판단 근거는 커밋된 경로와 SHA로만 적었다 |

## 검토 결론

조건부 수락한다. critical·high는 없다. 이전 리뷰 F1–F4와 designer V-CN-01·02, tester D-CT-01·O-CT-02·04는 코드와 증거로 해소를 확인했다. 새 발견은 low 3건(F-R1 수신 한도 보류 행의 방향 변화, F-R2 보류 행 직접 차단 POST, F-R3 기존 2행 데이터의 재요청)이며 병합을 막지 않는다. 다음 단위에서 수정하거나 정본에 예외로 기록한다. 이 결론은 AI 검토자의 판단이며 OCR의 자동 판정이 아니다. `review.py check` 통과는 기록 형식 검사이며 테스트 성공을 보증하지 않는다.
