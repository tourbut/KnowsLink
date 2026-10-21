---
title: designer 완료 기록
status: draft
updated: 2026-10-22
owner: designer
summary: 지시서와 완료 보고를 보존한다.
---

## SAR-PUBLISH-001-VISUAL — 2026-10-22
<!-- fullops-attempt: SAR-PUBLISH-001-VISUAL 7b6d0140a0f6428b80aeb53f76be5647 -->

---
title: "SAR-PUBLISH-001-VISUAL — C:/Program Files/Git/privacy·/terms·첫페이지·설정 KO/EN 3폭 직접 시각 검수(PC-06)"
status: draft
updated: 2026-10-22
owner: designer
tasks: [SAR-PUBLISH-001-VISUAL]
summary: "C:/Program Files/Git/privacy·/terms·첫페이지·설정 KO/EN 3폭 직접 시각 검수(PC-06)"
attempt: 7b6d0140a0f6428b80aeb53f76be5647
base: 17a360b2051d96e39646e1b048c931f6f0ede88b
subagent_level: standard
test_level: lite
---

# SAR-PUBLISH-001-VISUAL — C:/Program Files/Git/privacy·/terms·첫페이지·설정 KO/EN 3폭 직접 시각 검수(PC-06)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-22
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: DEV 보고의 격리 fixture(`TestPublishBrowserFixture`, 빌드 태그 `integration,publishfixture`, `python3 remote.py fixture`, 수명 24시간) 또는 DEV 렌더를 보고 직접 시각 검수 기록·캡처를 커밋 / 운영 서버·운영 env·live Tunnel·Cloud Console·실계정 Google·실메일 사용, 제품 코드·정책 문안 수정, force-push
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

main 17a360b2에 SAR-PUBLISH-001-DEV(dev 722feed9, 제품 e20c199)가 병합됐다. 서버 Google Calendar 연결 제거(D018), 시작 시 일회 정리(저장 토큰 revoke 시도 후 무조건 삭제, 미소비 조회 gate denied, 개수만 기록), KO/EN `/privacy`·`/terms`(D017)와 첫페이지 하단·`#start`·설정 링크. 독립 코드 리뷰 수락(critical/high 0, low 4: `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-REVIEW-review/report.md`). DEV 검사·격리 fixture·3폭 렌더: `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md`. 운영 배포 전이다.

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md`와 연결 규칙(testing.md `담당과 반복 범위`)이다. 기준 커밋 main `17a360b2051d96e39646e1b048c931f6f0ede88b`. 정책 정본·완료 조건 PC-01–06·RC-01–05: `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`. 화면 정본 D04: `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`. 결정 D017·D018은 `decision.py list`. 산출물 ID D01~D13: `docs/deliverables/README.md`. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`(PC·RC 표와 KO/EN 문안)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md`(격리 fixture 절·화면 근거·검증하지 못한 것)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-REVIEW-review/report.md`(low 4건)
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`(정책 페이지·첫페이지 하단·설정 절)

## 해야 할 일과 파일 소유권

- [x] `/privacy`·`/terms`·첫페이지(하단·`#start`)·설정의 KO/EN을 1440·390·320px에서 직접 본다. DEV 렌더(`render/`)는 참고이며 대체가 아니다. 가능하면 DEV 보고의 격리 fixture(`TestPublishBrowserFixture`, 빌드 태그 `integration,publishfixture`, `python3 remote.py fixture`, 수명 24시간)에서 실제 화면을 연다.
- [x] D04·SAR-PUBLISH-001 화면 기준(기존 공개 화면 스타일 재사용, 고유 title, `aria-current`, 가로 넘침 없음)과 대조한다. 정책 페이지 하단에서 '홈' 링크 위치가 D04와 맞는지 판정한다(DEV 보고 미결 항목).
- [x] 피할 것이 없는지 확인: 크림 배경, 알약 버튼, 모노스페이스 라벨, 장식 그라디언트, 이모지 아이콘, 영어판의 한국어 잔존.
- [x] 판정(PASS/FAIL, 항목별 심각도)과 캡처를 `qa-reports/SAR-PUBLISH-001-VISUAL-designer-test/`에 남긴다.

소유 파일: 위 QA 폴더, 이 인박스. 제품 코드·D04 수정은 하지 않는다(FAIL이면 수정 지시안을 보고에 적는다).

## 완료 기준과 검증

- PC-06의 designer 몫 판정이 report.md에 있다. FAIL은 화면·폭·언어·재현 위치·수정 제안과 함께 적는다.
- 기록 커밋 뒤 `lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/designer --from 17a360b2051d96e39646e1b048c931f6f0ede88b` ERROR 0(스크립트 경로는 tester 지시서와 같다). 캡처는 판정 항목에 연결한 것만 커밋한다.
- `work.py finish` 후 `worker_done`.

### UI 작업의 디자인 기준과 검증

- 디자인 문서: `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`, `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` 화면 기준 절.
- 테마 전환·디자인 lint: 이 제품은 단일 테마이고 디자인 lint가 없다. 해당 없음으로 적는다.

## 갱신할 산출물

없음(검수 기록만).

## 기대 산출물

`qa-reports/SAR-PUBLISH-001-VISUAL-designer-test/report.md`와 판정 캡처.

## 제약·협업·후속

tester QA(SAR-PUBLISH-001-TESTER)와 병렬로 진행한다. 파일 소유권이 겹치지 않는다. 결과 뒤 coor가 ops 배포를 배정한다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-PUBLISH-001-VISUAL-packet.json` / SHA `17a360b2051d96e39646e1b048c931f6f0ede88b` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md` (document_read) · 줄 7, 7, 70, 70, 75, 75 · inferred · 필수
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_read) · 줄 9, 18, 18, 23, 24, 47, 47, 129, 217, 241, 321, 335, 359, 366, 372, 394, 394 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/ops.md` (document_read) · 줄 70, 70 · inferred
- `.fullops-squad/docs/deliverables/README.md` (document_read) · 줄 37, 37 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 277, 277 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read) · 줄 127, 127 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 335, 335 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.4.md` (document_read) · 줄 44, 68 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 15, 15, 20, 20, 21, 21, 136 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 57, 57 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 전체/미확인 · inferred
- `README.md` (document_read) · 줄 45, 45 · inferred
- `deploy/knowslink/tunnel/public-ingress.yml` (impact_check) · 줄 4, 4, 7, 7 · inferred
- `docs/operations.md` (document_read) · 줄 27, 27 · inferred
- `internal/relay/human_pages.go` (impact_check) · 줄 46, 46 · inferred
- `internal/relay/landing_integration_test.go` (impact_check) · 줄 92, 93, 95, 100, 104 · inferred
- `internal/relay/landing_pages.go` (impact_check) · 줄 70, 70, 78, 78 · inferred
- `internal/relay/policy/terms.en.html` (impact_check) · 줄 1, 3, 16, 37, 42, 45 · inferred
- `internal/relay/policy_pages_test.go` (impact_check) · 줄 31, 31, 32, 32, 56, 56, 57, 57, 72, 72, 76, 76, 87, 87, 93, 93 · inferred
- `scripts/check_public_ingress.py` (impact_check) · 줄 9, 9, 18, 18 · inferred
나머지 참조 후보 16건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-<역할>-test/report.md` 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

브랜치 `fullops/designer` / 완료 SHA `040dfccbd782a02b05503d2a4cec1641e6b9722f`(기준 main `17a360b2`) — 기록 커밋이며 제품 코드는 바꾸지 않았다.
- 변경 이유: PC-06 designer 직접 시각 검수. 결과 **PASS**, critical·high 0, medium 1(V-01), low 3(V-02~V-04), 문서 1(D-01). 상세: [report.md](../docs/evaluations/qa-reports/SAR-PUBLISH-001-VISUAL-designer-test/report.md)
- 검토 필요: V-01(전역 `a{display:inline-flex;min-height:44px}`가 문장 속 링크 줄 간격을 약 44px로 벌림 — `#start` 한 줄·개인정보처리방침 §6)은 운영 배포 전 DEV 수정을 권한다. 배포를 막지는 않는다.
- 홈 링크 판정: 정책 페이지 하단 '개인정보처리방침 · 서비스 약관 · 홈 · 설치·사용 문서 · 맨 위로'를 승인한다(첫페이지 하단에는 '홈' 없음). D04 배치도에 '홈'이 빠져 D-01로 보고한다(이 과제는 D04를 바꾸지 않았다).
- 검증한 것: Chrome headless=new + CDP로 1440·390·320px 고정, 8개 화면(`/privacy`·`/terms`·첫페이지·설정 × KO/EN) 24장을 직접 확인, 가로 넘침 0, 영어판 한국어 잔존 0, 금지 스타일 0, 대비 6.19:1·7.96:1. lint `--from 17a360b2` HEAD 040dfccb: ERROR 0, WARNING 0, 실행 불가 1(Windows에 make가 없는 product-lint, DEV와 같은 사유). 이 과제는 코드를 바꾸지 않아 의존성 변경 없음.
- 못 한 것: 원격 격리 fixture 서버 경유 화면(언어 쿠키 전환·로그인 회원 저장 언어)은 SSH 자격이 dev 워크트리에 있어 쓰지 못했다. 렌더 HTML은 단위 테스트 출력이다. 키보드 순회·200% 확대 미실행. Windows `go test ./internal/relay`의 `TestFrozenParsingAndSigning`(registry integrity failure)은 이 과제와 무관한 기존 실패로 보인다.
- 산출물·로그: report.md와 `captures/` 6장, `measure.json`, `cdp_shots.mjs`, `docs/evaluations/jev/SAR-PUBLISH-001-VISUAL-packet-outcomes.json`.
- 남은 일: dev — V-01(필수 권장)·V-02·V-03 수정. coor — D04 '홈' 배치도 반영 배정(D-01). tester — 독립 QA(SAR-PUBLISH-001-TESTER)와 서버 경유 화면 확인. ops — 배포 전 V-01 수정 반영 여부 확인.

## SAR-PUBLISH-001-D04 — 2026-10-22
<!-- fullops-attempt: SAR-PUBLISH-001-D04 b0cfa862fa4241b59e1971229fd3e6ef -->

---
title: "SAR-PUBLISH-001-D04 — D04 정책 페이지 하단 배치도에 '홈' 반영(D-01)"
status: draft
updated: 2026-10-22
owner: designer
tasks: [SAR-PUBLISH-001-D04]
summary: "D04 정책 페이지 하단 배치도에 '홈' 반영(D-01)"
attempt: b0cfa862fa4241b59e1971229fd3e6ef
base: 179bdfe8d1b0b9a70f8e45c979bd19206b656c86
subagent_level: standard
test_level: lite
---

# SAR-PUBLISH-001-D04 — D04 정책 페이지 하단 배치도에 '홈' 반영(D-01)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-22
- From / To: coor → designer
- 상태: done
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: D04 문서 수정과 커밋 / 제품 코드 수정, 운영 서버·운영 env·live Tunnel·Cloud Console·실계정 사용, 정책 문안 변경, force-push
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / (dispatch 후 기록) / run_76a9511dbb0d / (dispatch 후 기록) / (dispatch 후 기록)

`상태`는 과제 상태다. front matter `status`는 문서 상태(draft/review/approved, lint DOC-003 어휘)라 서로 다른 값이 맞다.

## 현재 상황과 확인 근거

main 179bdfe8: SAR-PUBLISH-001-DEV(캘린더 제거·/privacy·/terms) 병합, 독립 리뷰 수락, tester QA 62/62 PASS(`.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-TESTER-test/report.md`), designer 시각 검수 PASS(medium 1, low 3, 문서 D-01: `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-VISUAL-designer-test/report.md`). 운영 배포 전에 시각 지적을 고친다.

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md`와 연결 규칙(testing.md `담당과 반복 범위`)이다. 기준 커밋 main `179bdfe8d1b0b9a70f8e45c979bd19206b656c86`. 정책 정본·완료 조건 PC-01–06·RC-01–05: `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`. 화면 정본 D04: `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`. 결정 D017·D018은 `decision.py list`. 산출물 ID D01~D13: `docs/deliverables/README.md`. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md`(PC·RC 표와 KO/EN 문안)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-DEV-test/report.md`(격리 fixture 절·화면 근거·검증하지 못한 것)
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-VISUAL-designer-test/report.md`(V-01~V-04, D-01)
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`(정책 페이지·첫페이지 하단·설정 절)

## 해야 할 일과 파일 소유권

- [x] D-01: `.fullops-squad/docs/design-docs/mockups/service-member-ui.md`의 정책 페이지 하단 배치도 줄을 `개인정보처리방침 · 서비스 약관 · 홈 · 설치·사용 문서 · 맨 위로`로 고친다. 첫페이지 하단에는 '홈'이 없다고 한 줄 적는다.
- [x] V-01~V-03 수정 방향(본문 링크 inline, 하단 aria-current 굵게, 설정 링크 gap)이 D04 화면 기준과 다르면 D04에 한 줄로 반영한다. 같으면 바꾸지 않는다.

소유 파일: D04, 이 인박스.

## 완료 기준과 검증

- D04가 실제 화면(정책 하단 '홈' 위치)과 맞는다. front matter는 `deliverables.py --stamp --task SAR-PUBLISH-001-D04`로 갱신한다.
- 기록 커밋 뒤 `lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/designer --from 179bdfe8d1b0b9a70f8e45c979bd19206b656c86` ERROR 0(스크립트 경로는 dev 지시서와 같다). `work.py finish`, `worker_done`.

### UI 작업의 디자인 기준과 검증

해당 없음(문서 변경).

## 갱신할 산출물

D04(service-member-ui).

## 기대 산출물

갱신된 D04.

## 제약·협업·후속

dev의 V-01~V-03 수정(SAR-PUBLISH-001-FIX)과 병렬이다. dev 수정 뒤 시각 재확인은 coor가 별도로 정한다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-PUBLISH-001-D04-packet.json` / SHA `179bdfe8d1b0b9a70f8e45c979bd19206b656c86` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLISH-001-VISUAL-designer-test/report.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_update) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.4.md` (document_read) · 줄 44, 68 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-DEV.md` (document_read) · 줄 70 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-DEV.md` (document_read) · 줄 72 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001.md` (document_read) · 줄 114 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-001-DEV.md` (document_read) · 줄 101 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV.md` (document_read) · 줄 77 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-POLICY.md` (document_read) · 줄 67 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 184 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SESSION-PEER-001-DESIGN.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 136 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLISH-001.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md` (document_update) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_update) · 줄 전체/미확인 · inferred
나머지 참조 후보 4건은 정본 JSON의 remaining_context_paths에 있다. 완료 검사(packet-outcomes)는 위 목록만 요구한다.
미확인 4건: .fullops-squad/PLANS.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/mockups/service-member-ui.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/design-docs/module-design.md: missing, oversized, sensitive or invalid; manual inspection; .fullops-squad/docs/operations/ops-guide.md: missing, oversized, sensitive or invalid; manual inspection. 정본의 unknown/producer_status 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/<과제 키>-<역할>-test/report.md` 링크로 대신한다. 브랜치와 완료 SHA는 여기에 적는다. worker_done 메시지에만 적지 않는다. `남은 일`에 다른 역할의 후속(예: tester 독립 QA)을 쓰면 `work.py finish`가 PLANS.md `## 완료 보고 후속 요청`에 등재한다.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
- 브랜치 / SHA: fullops/designer / (finish 로그의 완료 SHA 참조)
- 변경 이유: 시각 검수 D-01. D04 정책 페이지 하단 배치도에 '홈'이 없어 제품 정본(SAR-PUBLISH-001 51행)과 어긋났다.
- 변경: D04 배치도 줄을 `개인정보처리방침 · 서비스 약관 · 홈 · 설치·사용 문서 · 맨 위로`로 고쳤다. 정책 페이지에만 '홈'이 있고 첫페이지 하단에는 없다고 적었다.
- 지시와 달라진 판단: V-01~V-03 방향(본문 링크 inline, 하단 aria-current 굵게, 설정 링크 gap 16px)은 D04에 없어 링크 모양 한 줄로 추가했다. 제품 정본 문안은 바꾸지 않았다.
- 검토 필요: 제품 정본 51행은 '첫페이지와 같은 하단에 … 홈'이라 읽혀 D04의 '첫페이지에 홈 없음'과 표현이 다르다. 실제 화면 기준으로 D04를 맞췄다. coor가 정본 문장 정리 여부를 정한다.
- 검증: `deliverables.py --stamp --id D04 --path docs/design-docs/mockups/service-member-ui.md --task SAR-PUBLISH-001-D04` 성공. lint 결과는 아래 finish 로그에 남는다. 제품 코드·화면은 바꾸지 않았고 시각 재확인은 하지 않았다.
- 산출물: D04(service-member-ui.md), `docs/evaluations/jev/SAR-PUBLISH-001-D04-packet-outcomes.json`.
- 의존성: 코드 변경이 없어 의존성 변경 없음.
- 남은 일: dev의 V-01~V-03 구현(SAR-PUBLISH-001-FIX) 뒤 시각 재확인은 coor가 별도로 정한다.

처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
