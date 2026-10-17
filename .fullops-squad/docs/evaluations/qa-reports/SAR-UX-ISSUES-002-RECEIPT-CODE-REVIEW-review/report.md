---
title: SAR-UX-ISSUES-002-RECEIPT-CODE-REVIEW 리뷰
status: draft
updated: 2026-10-17
owner: reviewer
summary: 71dcac7822b8 병합 전 독립 리뷰 보고서
---

# SAR-UX-ISSUES-002-RECEIPT-CODE-REVIEW 리뷰

- 검토자 / CLI / 모델: Claude Code, 세션 `58c84abf-a77e-48ff-9609-1e46b3b05d1c`, `claude-opus-5-5`. 구현자 세션 `ctx_d3598f925a06`과 다르다. AI 검토 판단이며 OCR 자동 판정이 아니다.
- base SHA / head SHA / merge-base: `7f341230a26abc4afa97850687e624f2a583715c` / `71dcac7822b84acc2ae458a533dda96af1e94fe4` / `7f341230a26abc4afa97850687e624f2a583715c`.
- snapshot: `C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-c64b790767fa4dcd86d177a1fd3a04cf`, detached·clean. snapshot은 수정하지 않았다. 결과는 기록 체크아웃 `coor`에만 썼다.
- OCR 버전 / 적용 규칙: open-code-review v1.13.0 (1672852), `rules.json` rule_sha256 `f1061481…94de8`. 공통 규칙 `.fullops-squad/rules/common/README.md`·`coding-style.md`·`security.md`·`testing.md`는 기준(7f34123)과 head에서 같다(마지막 변경 `d4ba3ca`). 지시서의 예외는 없다.
- 요구사항·완료 기준 원천: `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md`의 U2·U6, `.fullops-squad/handovers/logs/2026-10-15_to_dev.md`(RECEIPT-DEV)와 `2026-10-17_to_dev.md`(RECEIPT-DEV-FIX), tester 보고서 `fullops/tester` 8ff3506의 `SAR-UX-ISSUES-002-RECEIPT-TESTER-test/report.md` D-RT-01 절.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 60 / 48 / 12 / 21 / 39. skipped는 하네스 자동 생성 기록(jev, handovers/logs, qa-reports JSON·로그)이다. 지시서 범위와 실제 diff 차이는 없다. 제품 변경은 `internal/relay` 11개 파일이고 문서는 D05·D10·D11과 기록이다.
- lint(`lint.json`) ERROR 0 / WARNING 2 / 실행 불가 2: product-lint·product-test가 `unavailable`이다. Windows 호스트에 make가 없다. 처음 받은 lint.json은 reason이 비어 check가 실패했다. coor가 리뷰 중 두 명령에 reason을 기록했다(Windows make 없음, da53e96 Linux clone EXIT=0 로그 경로, 71dcac7 직접 실행 아님). 다른 필드는 그대로다(finding 3).
- SIZE-002: 추가 651줄 > 상한 400줄. 이 diff는 RECEIPT-DEV(조회 복구 전체)와 RECEIPT-DEV-FIX(옵션 한 줄)를 함께 담는다. FIX 지시서의 예상 규모는 "코드 수십 줄 이내"이고 실제 FIX 제품 변경은 5줄(9ae29ec..da53e96)이다. 나머지는 이미 tester가 검증한 RECEIPT-DEV 범위다. 기능 하나의 화면·복귀·admission 경로가 묶여 있어 분할하지 않은 것이 타당하다. SIZE-001 `member.go` 382줄(이전 377)은 5줄 추가이며 수락한다.
- DEP-001: `go.mod`·`package.json` 변경 없음. 표준 라이브러리(`encoding/base64`, `encoding/json`, `unicode`)만 사용한다.
- UI 디자인: 정본 D04 토큰·44px·기존 글꼴, 새 CSS·JS 없음. 디자인 lint는 `frontend.json` tools가 비어 해당 없음. 테마 전환은 토큰 변경이 없어 해당 없음. 실제 320/390/1440 렌더는 tester(9ae29ec 기준)가 수행했고 D-RT-01 수정 후 재검증과 designer 시각 수락은 아직 없다.

## 검토 범위

`preview.json`의 merge_base..head diff를 읽었다. 제품 코드는 diff와 연결 코드를 직접 읽었다: `identity.go:session`, `member.go:memberHit/setCookie/render/expired`, `capacity.go:requestBuckets/capacityResponse` 호출부, `store.go:transaction`, `human_http.go:humanAgentView/humanAgents`, `human.go:shortID`, `protocol.go:uuidPattern`, `http.go`의 라우트 등록, `google.go`의 callback 복귀 순서.

- `member_receipt.go`: 모든 경로가 `Status` int를 설정하므로 마지막 `view["Status"].(int)`가 panic하지 않는다. 없는·타인·보존 종료 ID는 같은 `receiptUnavailable` 문구와 403이다. UUIDv7 형식 검사 후 조회한다. 철회 agent는 재선택을 요구한다.
- `receipt_context.go`: 쿠키 1024자 상한, base64url·JSON 실패, 만료와 만료 상한(now+70분), 128바이트·제어/Cf 문자 검사가 있다. 쿠키 Member가 로그인 회원과 다르면 버린다. 소유하지 않은 agent는 선택·ID를 반사하지 않는다. 쿠키는 서명되지 않은 입력이며 권한 근거로 쓰지 않는다.
- `capacity.go`: admission transaction 안에서 인증한 view만 실패 화면에 쓴다. `requestBuckets`가 이미 `st.session`을 호출하므로 session idle 갱신이라는 새 부작용이 없다. DB 인증 전 실패와 로컬 동시 한도는 `Unavailable` 화면이며 사적 입력을 보이지 않는다.
- `member.go`/`google.go`: 복귀 경로는 고정 내부 경로다. Google은 Connection·Remote 복귀가 우선한다. `logout`·`logout-all`이 같은 함수로 쿠키를 지운다.
- 템플릿: html/template escape, 응답당 autofocus 한 곳, `aria-invalid`·`role=alert`·`role=status`를 확인했다.
- 테스트: `receipt_context_test.go`(단위), `receipt_integration_test.go`(integration 태그, 실제 로그인·Postgres), `receipt_browser_fixture_test.go`(integration && receiptfixture 태그와 `TEST_SYNTHETIC_DATABASE=1` 필수). `member_ui_test.go`와 `member_forms.golden`의 변경은 두 계약 변경만 반영하며 기대값 완화가 아니다.
- 증거 helper `remote.py`·`forward.py`: SSH는 RejectPolicy와 known_hosts를 쓴다. 자격 값은 추적되지 않는 `.env.server`에서 읽고 출력하지 않는다. 하드코딩 비밀은 없다.
- skipped 확인 방법: jev·handovers/logs·qa-reports JSON·로그는 같은 폴더의 `report.md`와 DEV-FIX 보고의 링크로 존재를 확인했다. lint·test 로그 두 개는 첫·끝 줄과 `EXIT=0`을 직접 읽었다.

## 발견 사항

critical·high·medium은 없다. low 3건이다. 1·2는 미해결·비차단이고 3은 해결했다.

1. low — `internal/relay/member_pages.go:102` (U2 / D-RT-01). 옵션 표시가 `구분값 / 이름 / 종류`다. 동명 agent를 닫힌 선택기에서 구별하는 D-RT-01의 핵심 영향은 해결된다. 그러나 U2는 "이름을 먼저 보이고"라고 정한다. 39자 이름이면 320/390px에서 종류가 여전히 잘린다. 지시서가 "구분값을 앞에 두는" 방식을 허용했으므로 차단하지 않는다. 후속: designer 시각 수락에서 U2 순서 예외를 기록하고 tester가 선택기 세 건을 재검증한다.
2. low — `internal/relay/member_receipt.go:27-31`. 세션이 없거나 무효이고 익명 IP 한도가 찬 GET `/home/receipts`는 429 조회 화면을 만든다. view에 `Own`·`Unavailable`이 없어 템플릿이 "조회할 내 에이전트가 없습니다"와 `/home/agents` 링크를 보인다. 로그인하지 않은 사용자에게 사실과 다른 안내다. 재현: 세션 쿠키 없이 같은 IP로 익명 한도를 채운 뒤 조회한다. 데이터 노출은 없다. 수정안: 이 분기에 `"Unavailable": true`를 넣는다.
3. low, 해결 — `lint.json` 기록. 71dcac7의 product-lint·product-test가 reason 없이 `unavailable`이었다. coor가 reason을 기록해 해결했다. 증거 한계는 남는다: 등록 검사는 71dcac7에서 직접 실행되지 않았고 같은 제품 코드의 da53e96 Linux 결과로 대신한다(아래).

## 검증 및 남은 제약

리뷰어는 lint.py와 전체 테스트를 재실행하지 않았다(지시). 재사용한 증거:

- `da53e961…-lint.log`·`da53e961…-test.log`: exact-HEAD da53e96의 Linux clone에서 `make lint` EXIT=0, `make test`(go test -race ./cmd/... ./internal/..., npm, python) EXIT=0. lint 로그의 `ERROR:` 줄은 셸 스크립트 테스트의 기대 출력이다.
- `git diff --stat da53e96 71dcac7`: 문서·로그·인박스만 바뀌었다. 제품 코드는 같다.
- `git diff --stat 9ae29ec da53e96 -- internal/`: `member_pages.go` 1줄, `receipt_context_test.go` 4줄. 9ae29ec의 제품 source는 1e49033이고 그 SHA의 integration race 검사(`1e49033…-integration.log`)가 종료 0이다. 옵션 문자열 변경은 integration 기대값(`value="agent_a" selected`)에 영향이 없다.
- tester 8ff3506: U6 격리·복구·키보드는 통과했다. D-RT-01만 남았다.

미실행·남은 일: D-RT-01 수정 후 실제 브라우저 폭 판정(tester 선택기 세 건), designer 시각 수락, 실제 200% 확대·스크린리더·Firefox·Safari, 명시적 로그아웃 브라우저 확인. 보존 종료와 만료 receipt는 integration 검사가 근거다.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과(확인/미확인/해당 없음) | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `design-docs/SAR-UX-ISSUES-002.md` U2·U6 | 확인 | U2 이름 선두와 D-RT-01 구분값 선두의 차이를 정본에 예외로 기록해야 한다(finding 1). |
| 구조와 구현/미완료 상태 | `design-docs/module-design.md` "SAR-UX-ISSUES-002 조회 입력과 복귀", `interface-design.md` "SAR-UX-ISSUES-002 요청 상태 조회" | 확인 | 코드와 일치한다. |
| 실행·검증 방법과 증거 | `qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md`, `…-DEV-FIX-test/report.md` | 확인 | 로그 링크가 해석된다. fixture는 철회 상태라 재검증에 재준비가 필요하다. |
| 운영·복구(해당 시) | `operations/user-guide.md` "요청 상태 조회" | 확인 | "DEV 후보" 표시가 남아 있다. 수락 후 갱신한다. |
| 다음 작업·담당·재개 조건 | DEV-FIX 보고 "못 한 것", `handovers/logs/2026-10-17_to_dev.md` 완료 보고 | 확인 | tester 재검증·designer 수락·main 통합은 coor 담당이다. |
| 로컬 링크·절 접근/지원 한계 | 위 보고서의 상대 링크 | 확인 | 원 fixture·private 로그는 레포 밖 경로라 접근을 확인하지 않았다. |
| snapshot 정리 후 정본 접근 | 리뷰 폴더 `coor` 기록 체크아웃 | 확인 | 결과는 snapshot 밖에 있다. |

## 검토 결론

수락한다. critical·high가 없다. 조회 복구의 계정 격리·입력 검증·실패 경로가 요구 U6과 일치한다. D-RT-01의 코드 수정을 확인했다. low 3건은 병합을 막지 않는다. 다만 제품 최종 수락에는 tester 선택기 재검증과 designer 시각 수락이 따로 필요하다. `review.py check` 통과는 기록 검사이며 이 판단이나 테스트 성공을 자동 보증하지 않는다.
