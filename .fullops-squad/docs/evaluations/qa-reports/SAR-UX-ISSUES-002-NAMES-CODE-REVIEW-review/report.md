---
title: SAR-UX-ISSUES-002-NAMES-CODE-REVIEW 리뷰
status: review
updated: 2026-10-19
owner: reviewer
summary: b2e683ae677c 병합 전 독립 리뷰 보고서. critical·high·medium 없음, low 1건(D-NM-01) 미해결로 수락
---

# SAR-UX-ISSUES-002-NAMES-CODE-REVIEW 리뷰

- 검토자 / CLI / 모델: Orca worker(term_fb9aff74) / Claude Code / claude-opus-5-5. reviewer_session `names002-code-review-opus`, 실제 host 세션 ID `c5a56fd2-5d78-4607-adc3-3df833496e3d`. 구현자 세션 `ctx_bc77cf9072cc`와 다르다.
- base SHA / head SHA / merge-base: `8d0db4bf59e905c9f398123c4e2e5bc13ed4c3fa` / `b2e683ae677c78cd7bdc47c75820bcb191b94838` / `8d0db4bf59e905c9f398123c4e2e5bc13ed4c3fa`. snapshot `C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-188ae82872b94838bf6069b8b6f73f48`(detached, clean, 수정하지 않음).
- OCR 버전 / 적용 규칙: open-code-review v1.13.1 (8eb8334) delegate preview·rule(prepare가 생성). rule_sha256 `f1061481…94de8`. 공통 규칙 `fullops-common-0.3.3`(`.fullops-squad/rules/common/README.md`, coding-style·testing·security), `.fullops-squad/project.md`, 문서 규칙 `docs/agents/document-writing.md`. 기준 커밋 8d0db4b이며 base..head에서 rules 파일 변경이 없다. 지시서 `적용 기준과 예외`의 예외(frontend host 대신 claude, Windows make 부재로 격리 Linux clone 실행)를 확인했다.
- 요구사항·완료 기준 원천: `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` U2·U4·U7, `.fullops-squad/handovers/logs/2026-10-19_to_dev.md`(지시·완료 기록), DEV 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-DEV-test/report.md`, 참고 tester 보고 `C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/report.md`.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 62 / 37 / 25 / 29 / 33(OCR 대상 중 prepare 자동 skipped 8: jev 7·handover 로그 1, probe.json 1 포함). 커버리지(reviewed+skipped) 62/62 = 100%. 지시서 소유 범위(`internal/relay/` 회원·친구 화면/handler/템플릿·테스트, D05/D10/D11/D13, own jev·QA 보고) 밖의 변경은 PLANS.md·contexts/dev.md뿐이며 coor·dev 기록 관례다.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: ERROR 0 / WARNING 5(SIZE-001 4, SIZE-002 1) / 실행 불가 2(product-lint, product-test). 원 status `unavailable`·exit null은 보존됐다. 사유: Windows 호스트에 make가 없다. coor가 보완한 Linux 로그를 확인했다(아래 검증 절).
- SIZE-002: 지시서 예상 규모 `internal/` 500–900줄·10–15파일 / lint 집계 추가 573줄(25파일) / 완료 기록 `internal/` 19파일 +489/−58. 파일 수 초과 4개는 select→radio·legacy 라벨 변경에 따른 integration 기대 문자열 갱신이다. 한 템플릿 파일을 공유해 분할 이득이 없다. 수락한다.
- SIZE-001: `human_http.go` 464줄(+52), `member.go` 394줄(+12), `member_agents.go` 380줄(+10), `PLANS.md` 1969줄. 모두 기존 초과 파일의 소폭 증가다. 분할은 후속 정리 대상이며 차단하지 않는다.
- DEP-001: 해당 없음. `go.mod`·`go.sum`·`adapters/package*.json` 변경이 없고 lint에 DEP-001이 없다. 새 import는 표준 `strings`뿐이다.
- UI 디자인: 정본 D04 `docs/design-docs/mockups/service-member-ui.md`, 공용 `memberStyle`·`human-agent-label` 템플릿을 사용한다. 예외 근거: D04 A-07(읽기 전용 값은 `code`)과 U2(읽기 전용 입력)가 충돌해 더 새 결정인 U2를 따라 `textarea.share`를 골랐다(exec-plan 결정 절). 디자인 lint는 이 레포에 등록돼 있지 않다(해당 없음). 테마 전환: CSS 토큰 변수만 사용하고 새 색을 추가하지 않았다. 실측은 DEV CDP probe(1440/390/320 넘침 0)와 tester Chrome·Edge 레이아웃 49 PASS·키보드/대비 27 PASS다. 실제 200% 확대·스크린리더는 DEV·tester 모두 not_run이며 designer 시각 검수 단계의 몫이다.

## 검토 범위

result.json 62개 `(path,status)` 모두 `review_status`와 `reason`이 있다.

- 제품 코드 9개(`http.go`, `human_http.go`, `human_pages.go`, `member.go`, `member_agents.go`, `member_pages.go`, `member_receipt.go`, `receipt_context.go`, `member_forms.golden`): base..head diff 전체와 호출 문맥(`requestFriend`, `humanMember`, `friendByID`, `humanAgentView`, `agentName`, `shortID`, `humanPage`/`humanAction`, `connectionPage`)을 읽었다.
- 테스트 10개: 새 테스트 3개는 전체를, 수정 테스트 7개는 diff를 읽었다. 기대값 변경은 select→radio, legacy 요약의 ID 위치 이동에 맞춘 것이다. 검증 강도를 낮춘 곳은 없다.
- 문서 9개와 DEV 도구 2개(`remote.py`, `render/probe.mjs`): 전체 또는 diff를 읽고 코드와 대조했다.
- skipped 33개는 prepare의 공통 사유(하네스 자동 생성·검증 증거)를 따른다. 무결성 확인 방법: `raw-log-sha256.txt`의 516bf62/3c1567f lint·test·integration·integration-stop 로그 sha256을 `git show HEAD:<path> | sha256sum`으로 대조해 일치를 확인했다. 작업 트리 파일은 autocrlf로 hash가 다르다. 각 로그의 마지막 `EXIT=0`을 확인했다. PNG 15개는 DEV 보고와 probe.json에 연결된 시각 증거이며 열람하지 않았다. 영향은 없다. tester가 같은 SHA에서 독립 측정했다.

## 발견 사항

### D-NM-01 (low, 미해결) 이름 거부 사유가 보이지 않고 보이지 않는 문자가 입력란에 되돌아온다

- 파일·줄: `internal/relay/human_http.go:367-385`(profile·agent-name). 같은 문구가 `member_agents.go:222,237`, `remote_auth.go:358`에도 있다.
- 재현: tester 보고 D-NM-01. `POST /home/agent-name`에 `name=ab\x07cd`, U+202E, U+2066을 넣으면 422와 "이름은 1–40자로 입력하세요"만 나온다. 입력값 `r.FormValue("name")`이 escape된 채 `value`에 그대로 돌아온다. 코드에서도 같은 경로를 확인했다.
- 영향: 사용성. 이름은 바뀌지 않고, html/template escape로 XSS는 아니다. 다른 사용자에게 보이지 않는다. U2-A의 "잘못된 제어문자" 확인 기준의 안내 품질 잔여다.
- 판단: 이번 diff가 만든 결함이 아니다(기존 코드). 수정 제안: 문구에 제어·방향 제어 문자 금지를 추가하고, 되돌릴 값에서 해당 문자를 제거한다. 후속 과제로 PLANS.md에 남긴다. 병합을 차단하지 않는다.

### 중점 검토 결과 (결함 없음)

- 친구 요청 존재 비노출·차단 비공개: 새 422 분기는 빈 코드(공백 제거 후)와 `consent` 누락만 본다. 코드 존재·차단·한도를 조회하기 전에 판단하므로 응답이 대상 존재와 무관하다. 형식이 맞으면 `requestFriend`(기존 placeholder 방어 유지) 뒤 항상 `303 /home/friends?n=requested`다. 공백만 제거하고 대소문자·Unicode 정규화는 없다. 매처는 `hmac.Equal` 정확 비교를 유지한다. 친구 상세 거부는 없음·대기·차단 모두 `sender_not_allowed` 한 가지로 친구 목록에 복귀한다.
- 이름 표시 권한 경계: `visibleAgent`는 자기 소유 또는 `friendsNow`(양쪽 humanMember 존재, 친구 state `active`, `BlockA`·`BlockB` 모두 false)일 때만 이름을 돌려준다. 그 밖에는 "확인할 수 없는 에이전트"와 ID만 둔다. Client·Short는 비운다. 상대 프로필 추가 조회는 없다. 승인 화면(`gateView`, 기존 owner 검사 뒤), legacy pair 라벨, 연결 화면(기존 `c.Owner` 검사 뒤)이 이 경계를 쓴다. 친구 상세 Title의 상대 이름은 이미 active·비차단 검사를 통과한 분기에서만 설정된다. 테스트: `TestVisibleAgentNamesOnlyOwnerAndCurrentFriend`(대기·제삼자·차단), fixture의 legacy 비회원 agent.
- 관계 세대·GET 무변경: 새 코드는 GET 경로에서 상태를 쓰지 않는다(`visibleAgent`·`distinguish`는 값 복사본 수정). 422 분기는 `requestFriend`를 호출하지 않는다. 테스트: `TestFriendViewsNeverCreateRelationships`, integration의 GET·잘못된 POST 뒤 `len(Friends)==0`, 이름 변경 뒤 Generation 불변, 만료 뒤 재요청 Generation+1.
- 입력 escape: `RequestCode`, `Name`(동의 문구·보내는 사람), Title, aria-label 속성은 모두 `html/template` 문맥 escape다. `TestFriendPagesShareOwnCodeAndKeepRequestInput`이 `"><b>`의 escape를 단언한다.
- 계정 격리: 친구 화면·계정 화면의 코드는 세션 회원 자신의 `m.FriendCode`다. 계정 화면 textarea에는 `name` 속성이 없어 폼으로 제출되지 않는다. receipts 선택기는 기존 소유 agent 목록만 radio로 바꿨다. legacy `MineAgent` 구분값은 자기 `ownAgents` 확장값만 재사용한다.
- 관찰(결함 아님): `distinguish`는 한 소유자 목록 안에서만 충돌을 넓힌다. 대화 화면의 내 agent와 친구 agent 사이 32bit 구분값 충돌은 확장하지 않는다. 확률이 무시할 수준이고 전체 ID가 상세에 있다.

## 검증 및 남은 제약

- 리뷰어는 lint·테스트를 다시 실행하지 않았다. 지시에 따라 기존 증거만 읽었다.
- 재사용 증거: snapshot에 커밋된 DEV 로그 `516bf6247f06843c97b0292d0db01961843713ea-{lint,test,integration}.log`, 각각 `EXIT=0`이다. `make test`는 `go test -race ./cmd/... ./internal/...`와 adapters·shell guard 검사를 포함한다. integration은 `TestHumanNamesAndFriendCodePostgres` PASS를 포함한다(`TestRemoteConsentBrowserReturn`은 TEST_BROWSER 부재로 not_run, remote 화면 변경 없음). 실행 환경은 `remote.py`의 격리 Linux clone이다(HEAD 40자 검증, clean 확인, 고정 digest postgres:17-alpine).
- 환경 동일성 근거: `git diff --stat 516bf62 b2e683a -- internal adapters cmd Makefile go.mod`가 비어 있다. b2e683a는 516bf62 위에 문서·완료 기록만 더했다. 따라서 516bf62의 EXIT=0이 head 제품 소스에 적용된다. lint.json의 reason이 이 경로를 가리키고 로그 존재와 hash 일치를 확인했다.
- 미실행 범위: 실제 200% 확대, 스크린리더, 실제 Tab 키 입력의 일부, 테마 전환 실측은 DEV not_run이다. tester가 Chrome·Edge 키보드·대비를 일부 확인했다. designer 시각 검수가 남는다.
- 남은 일: D-NM-01 후속 수정, SIZE-001 대상 파일 분할 검토, designer 시각 검수 후 main 통합.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과(확인/미확인/해당 없음) | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `docs/design-docs/SAR-UX-ISSUES-002.md` U2·U4·U7, `docs/exec-plans/phases/SAR-UX-ISSUES-002-NAMES-DEV.md` 결정 절 | 확인 | D04 A-07과 U2 충돌의 선택 근거가 기록됐다 |
| 구조와 구현/미완료 상태 | `docs/design-docs/module-design.md` "SAR-UX-ISSUES-002 이름 표시·친구 코드·친구 요청 안내", `interface-design.md` 같은 절 | 확인 | 모듈표와 실제 함수·테스트 이름이 일치한다 |
| 실행·검증 방법과 증거 | DEV 보고 `SAR-UX-ISSUES-002-NAMES-DEV-test/report.md`, `remote.py`, `raw-log-sha256.txt` | 확인 | not_run 항목이 명시됐다 |
| 운영·복구(해당 시) | `docs/operations/transition.md` "DEV 후보 인계 (2026-10-19)" | 확인 | 새 schema·migration 없음, 기존 relay 교체 절차 |
| 다음 작업·담당·재개 조건 | `PLANS.md` "잔여 작업 일괄 수행 — 2026-10-18" | 확인 | D-NM-01 후속을 coor가 PLANS.md에 추가해야 한다 |
| 로컬 링크·절 접근/지원 한계 | contexts/dev.md·D10·D13·D11의 DEV 보고 링크 | 확인 | snapshot에서 링크 대상 파일 존재를 확인했다 |
| snapshot 정리 후 정본 접근 | 본 리뷰 폴더(coor 체크아웃)와 main 병합 뒤 같은 경로 | 확인 | 근거 경로가 모두 레포 상대 경로다 |

## 검토 결론

수락한다. critical·high·medium 결함이 없다. 미해결은 low 1건(D-NM-01, 기존 코드의 안내 품질)뿐이다. 중점 경계(친구 요청 존재·차단 비노출, 비친구 이름 미노출, GET 무변경·세대, 입력 escape, 계정 격리)는 코드 검토와 테스트·EXIT=0 로그로 확인했다. 이 결론은 AI 검토 판단이다. OCR의 자동 판정이 아니다. `review.py check` 통과는 기록 검사이며 테스트 성공을 자동 보증하지 않는다.
