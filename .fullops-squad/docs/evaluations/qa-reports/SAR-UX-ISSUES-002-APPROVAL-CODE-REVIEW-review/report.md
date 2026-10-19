---
title: SAR-UX-ISSUES-002-APPROVAL-CODE-REVIEW 리뷰
status: review
updated: 2026-10-20
owner: reviewer
summary: 9d13a05ba93c 병합 전 독립 리뷰 보고서. critical/high 0, low 2(비차단), 수락
---

# SAR-UX-ISSUES-002-APPROVAL-CODE-REVIEW 리뷰

- 검토자 / CLI / 모델: Orca reviewer `approval002-code-review-opus` / Claude Code / claude-opus-5-5. 실제 host 세션 ID `a36fe212-a23f-4d69-9b09-c767707929c1`. 구현자 세션 `ctx_0974e0d7a12d`(fixture 수정 `ctx_61bc49c73c11`).
- base SHA / head SHA / merge-base: `688305054dc4f3570f0c9f0d1e7917c35a19db3a` / `9d13a05ba93c7c085034213590ce661305606010` / `688305054dc4f3570f0c9f0d1e7917c35a19db3a`.
- snapshot: `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-3b2961ffec0247909dfa434e663bd4f9`, detached HEAD 9d13a05, `git status --short` 빈 출력. snapshot은 읽기만 했다.
- OCR 버전 / 적용 규칙: open-code-review v1.13.1, `.fullops-squad/review/rule.json`(rule_sha256 f1061481…). 공통 규칙 fullops-common-0.3.3(`.fullops-squad/rules/common/` README·coding-style·security·testing, 기준 6883050). 이번 범위에서 rules·review·fullops.json 변경 없음.
- 요구사항·완료 기준 원천: `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` U5·U7, `.fullops-squad/handovers/logs/2026-10-19_to_dev.md`의 `SAR-UX-ISSUES-002-APPROVAL-DEV` 두 절(본 과제·fixture 수정), DEV 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-DEV-test/report.md`.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 81 / 37 / 44 / 27 / 54. 지시서 소유 파일(`internal/relay/` 승인·회원 화면·테스트·fixture, 기술 문서, DEV 보고)과 실제 diff가 일치한다. 범위 밖 제품 파일 변경 없음.
- lint(`lint.json`) ERROR 0 / WARNING 3 / 실행 불가 2. 실행 불가는 Windows에 make가 없어서다. coor가 보완한 reason의 로그를 확인했다: `ea94128…-lint.log`·`ea94128…-test.log` 마지막 줄 `EXIT=0`, `f43d05c…-integration.log` `PASS`·`EXIT=0`. `git diff ea94128..9d13a05`의 제품·테스트 변경은 없다. `f43d05c..9d13a05`는 build tag가 붙은 `approval_browser_fixture_test.go` 6줄뿐이며 `ea94128…-fixture-compile.log`(`go vet`·`go test -c -tags=integration,approvalfixture`) `EXIT=0`이다. 원 status는 보존돼 있다.
- WARNING: SIZE-001 `PLANS.md`(1983줄)·`http.go`(688줄)는 기존 누적 크기다. 이번 변경은 각 +2줄·+15줄이고 요약 로직은 새 파일로 분리했다. 수락한다.
- SIZE-002: 지시서 예상 300–700줄 / 실제 코드 추가 511줄 / 예상 범위 안. 분할 불필요.
- DEP-001: 해당 없음. go.mod·package.json·lock 변경 없음.
- UI 디자인: 정본 D04·SAR-UX-ISSUES-002.md, 공용 `human-agent-label`·`dl.kv`·native details 재사용, 새 테마·폰트·JS 없음, CSS는 `.deny-help` 한 규칙. 테마 전환 없음(해당 없음). DEV 렌더 24장 가로 넘침 0(`render/probe.json`). designer 시각 수락·tester QA는 별도 과제다.

## 검토 범위

result.json의 81개 `(path,status)`에 모두 기록했다. 제품 Go 파일 15개, 테스트·golden, 기술 문서 7개, DEV 보고·보조 스크립트 3개, 요구 정본인 handover 로그를 직접 읽었다. jev 패킷·route JSON, 원시 로그, png는 prepare가 skipped로 채웠다. 원시 로그는 직접 `EXIT=` 줄과 sha256을 확인했다. png는 `render/probe.json`과 DEV 보고 표로 존재와 결과를 확인했다.

중점 확인 결과:

| 중점 | 결과 | 근거 |
|---|---|---|
| 요약이 검증된 envelope 필드만 사용 | 확인 | `gateView`는 `gateState`→`verifiedGate`가 돌려준 envelope만 `gateSummary`에 넘긴다. `verifiedGate`는 Parse(intent별 closed schema), 서명·kid, digest·ID, From/To, 세대, gate digest를 검사한다. `gateSummary`가 읽는 window·granularity_min·slot·timezone·commitment는 `protocol.go` Parse가 검증하는 필드와 같다. DEV 대응표와 일치한다 |
| hint·자유문·표시 이름의 승격 | 없음 | commitment는 `요청자가 적은 내용`, hint는 `참고 문구 — 승인 근거 아님` details에 텍스트 노드로만 나온다. scope 문구는 서버 상수 `gateScopes`다. 이름은 `visibleAgent` 권한을 그대로 쓰고 결정 경로는 이름을 읽지 않는다 |
| 서버 재검사 유지 | 유지 | `gateDecision`(http.go)은 변경되지 않았다: owner, CSRF(세션 token+gate ID HMAC), `State=="pending"`, 기한, `verifiedGate` active(current·parent receipt 기한 포함) |
| 낡은 화면·교차 gate token·중복 클릭 | 안전 | 낡은 제출은 pending·verified 재검사에서 409. CSRF가 path의 gate ID에 묶여 다른 gate에 403. 두 번째 클릭은 상태가 pending이 아니어서 409. 모두 `TestApprovalSummaryBoundaries`에 있다 |
| 원문 부재·서명 불일치·미지원 유형 | 안전 fallback | `e==nil`이면 Summary·Body·Hint를 넘기지 않는다. 템플릿은 alert와 기한·상태·복귀만 보인다. `gateKinds` 필터는 gate 부모가 `parentRouting`상 schedule.*뿐이라 기존 동작을 바꾸지 않는다 |
| HTML escape | 확인 | `html/template`. aria-label 안의 `human-agent-label`은 태그가 없는 텍스트라 attribute 문맥에서도 escape된다. 테스트가 `<script>`·`<a href>` 미출력을 단언한다 |
| O-NM-04 익명 경계 | 확인 | `unverifiedReceipt`는 세션 cookie 존재만 본다. 익명이면 폼·하단 링크·restore 문구를 숨긴다. 조회 권한 판단에는 쓰지 않는다 |
| fixture 누출 | 없음 | `//go:build integration && approvalfixture`. Makefile·CI에 태그 참조가 없고 `remote.py`만 명시적으로 쓴다. 합성 계정, loopback httptest, 허용 GET/POST 목록, Origin·Sec-Fetch-Site 검사, 실행 환경변수 2개 필수 |

## 발견 사항

1. low / `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-DEV-test/raw-log-sha256.txt` 1–18행 / 재현: `sha256sum *.log`와 `git show HEAD:<path> | sha256sum` 값이 기록값과 18건 모두 다르다. `sed 's/$/\r/' <log> | sha256sum`은 일치한다(예: `ea94128…-lint.log` → `7c93aac3…`). 영향: 기록 해시가 CRLF 작업 사본 기준이라 커밋 blob으로 재현할 수 없다. 로그 내용과 EXIT는 직접 확인해 수락 판단에는 영향이 없다. 수정 상태: 미해결, 비차단. 다음 증거부터 blob 기준 해시를 쓰거나 계산 방식을 함께 기록한다.
2. low / `internal/relay/gate_summary.go` 44–58행 / 재현 조건: Parse schema를 바꾸거나 `gateKinds`에 새 intent를 추가하고 요약 분기를 맞추지 않는 경우. 영향: 검사 없는 타입 단언이 승인 화면 요청에서 panic할 수 있다. 현재 코드에서는 Parse가 같은 필드를 보장해 발생하지 않는다. 수정 상태: 미해결, 비차단. intent 추가 시 Parse와 함께 갱신하라는 주석이나 comma-ok 단언을 권장한다.

critical/high/medium은 없다.

## 검증 및 남은 제약

- 리뷰어는 lint·테스트를 재실행하지 않았다(지시). 재사용 증거: DEV 원격 Linux 격리 clone 로그(`remote.py`, 자체 Postgres 컨테이너). f43d05c(제품 최종) lint·test·integration·fixture-compile EXIT=0, ea94128(fixture 최종) lint·test·fixture-compile·fixture-start·fixture-probe EXIT=0. fixture-probe는 manifest 7개 gate, `GET / -> 200`, 범위 밖 `/v1/pull -> 403`, A 로그인 204·gate 200을 보인다.
- 47117a4 integration EXIT=1은 DEV 보고에 원인(housekeeping이 원문 부재를 저장 상태로 바꿈, 버튼 문구 기대)이 기록돼 있다. 바뀐 단언은 ‘pending 유지’에서 ‘결정 미기록(approved/denied 아님)’이다. 낡은 제출이 결정을 만들지 못한다는 보안 경계는 그대로라서 정당한 보정으로 판단한다.
- 기존 테스트 변경은 문구·복귀 href 기대 교체다. 결정·서명·XSS 검사는 삭제·완화되지 않았다.
- integration은 기존과 같이 `-skip TestRemoteConsentBrowserReturn`이다. 이번 범위 밖이다.
- not_run: 실제 200% 확대·스크린리더·두 번째 브라우저(U7-A 일부)는 DEV 범위에서 미실행이고 tester·designer 과제로 남는다.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과(확인/미확인/해당 없음) | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `design-docs/SAR-UX-ISSUES-002.md` U5·U7, `contexts/dev.md` 2026-10-19 APPROVAL-DEV | 확인 | 없음 |
| 구조와 구현/미완료 상태 | `design-docs/module-design.md` gate_summary·verifiedGate 행, `interface-design.md` gate 절, DEV 보고 ‘기술 계획’ | 확인 | 없음 |
| 실행·검증 방법과 증거 | DEV 보고 ‘검사 로그’·‘tester/designer용 격리 fixture’, `remote.py` | 확인 | raw-log-sha256 기준 불일치(finding 1) |
| 운영·복구(해당 시) | `operations/transition.md` DEV 후보 인계(2026-10-19 APPROVAL-DEV) | 확인 | 운영 미반영, 기존 relay 교체 절차 |
| 다음 작업·담당·재개 조건 | `PLANS.md` 3단위 행·대기 행, DEV 보고 ‘대상과 결론’ | 확인 | tester QA·designer 검수·main 통합은 coor |
| 로컬 링크·절 접근/지원 한계 | DEV 보고의 `remote.py`·`raw-log-sha256.txt`·`render/probe.*` 상대 링크 | 확인 | 링크 대상 존재 |
| snapshot 정리 후 정본 접근 | 리뷰 폴더는 coor 체크아웃 `docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-CODE-REVIEW-review/` | 확인 | snapshot 밖에 보존 |

## 검토 결론

수락한다. critical/high 미해결이 없고, U5 요약·hint 분리·서버 재검사·낡은 제출·교차 token·중복 클릭·안전 fallback·escape·O-NM-04·fixture 격리가 코드와 테스트·로그 증거로 확인됐다. low 2건은 후속 정리 항목이다. 이 판단은 AI 리뷰어의 판단이며 OCR 자동 판정이 아니다. check 통과는 기록 검사일 뿐 테스트 성공을 자동 보증하지 않는다.
