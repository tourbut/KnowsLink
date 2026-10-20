---
title: SAR-UX-ISSUES-002-LOWS-DEV — 잔여 UI 단위의 친구 보류 행·부분 언어·내부 헤더 low 5건을 정리한다
status: draft
updated: 2026-10-20
owner: dev
tasks: [SAR-UX-ISSUES-002-LOWS-DEV]
summary: 잔여 UI 단위의 친구 보류 행·부분 언어·내부 헤더 low 5건을 정리한다
attempt: 54a72f9a5f2746358a4bf684c8604529
base: c3b51926dbae01bdb325a9d3be01ed1ad09275c0
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-LOWS-DEV — 잔여 UI 단위의 친구 보류 행·부분 언어·내부 헤더 low 5건을 정리한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → dev
- 상태: done
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 아래 low 5건 수정과 관련 테스트·기술 문서. 상태 JSON의 기존 필드 확장(`Friend.Held`에 이전 방향 보존 등)은 기존 데이터 호환 범위에서 승인. DB migration·새 경로·새 의존성·권한 확대·실계정·운영 배포·삭제/이동·force push는 coor에 ask.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

잔여 UI 1–5단위가 main c3b5192에 통합됐다. 각 단위 리뷰·QA에서 '후속'으로 미룬 low 중 제품 규칙에 어긋나는 것을 한 과제로 정리한다.

- F-R1·F-R2·F-R3: `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-FIX-REVIEW-review/report.md`의 해당 절(파일·줄·재현·제안 포함). F-R2는 차단 비공개 규칙의 예외다.
- O-LG-05·O-LG-06: `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-TESTER-test/report.md`의 관찰 표.
- 언어 내부 표식 헤더 `X-KnowsLink-Lang-Saved`가 렌더하지 않는 회원 응답(예: 303)에 남는 문제: 5단위 DEV 보고 '남은 한계', tester O-LG-06 끝부분.

## 적용 기준과 예외

fullops-common-0.3.3, `.fullops-squad/project.md`, 제품 정본 SAR-UX-ISSUES-002.md U4(재요청 결정·차단 비공개)·U3, D13 기술 문서. 기준 커밋 c3b5192. Windows make 없음: 격리 원격 Linux clone에서 make lint/test/integration. lint.py는 깨끗한 임시 detached worktree에서 실행하고 force 없이 정리.

## 먼저 읽을 문서

위 두 보고의 해당 절, `internal/relay/human.go`(requestFriend·decideFriend block), `internal/relay/human_http.go`(humanView incoming·Held), 언어 헤더 처리(`memberLang`·렌더 경로), 설정 언어 오류 문구, D13·data-model·interface-design의 친구 보류 행 설명, `.fullops-squad/contexts/dev.md`.

## 해야 할 일과 파일 소유권

- [x] F-R1: 수신 한도 보류 행이 받은 사람에게 이전 모습 그대로 보이도록 이전 방향(Inviter·CodeTail)을 보존한다(또는 동등한 방법). 받은 사람 화면이 재요청 전후 같음을 테스트로 보인다.
- [x] F-R2: 받은 사람이 보류 행에 block POST를 해도 요청자의 살아 있는 대기 요청이 '해제'로 바뀌지 않게(보이는 대기 요청 차단과 같은 결과) 한다.
- [x] F-R3: 기존 쌍 행+열린 대리 행 데이터에서 같은 코드 재요청이 열린 행 2개를 만들지 않게 한다(없는 코드와 같은 무동작).
- [x] O-LG-05: 영어 오류 'Choose 한국어 or English.'의 '한국어'(및 한국어판의 'English')에 부분 언어 lang 속성.
- [x] 언어 내부 표식 헤더가 어떤 응답으로도 클라이언트에 나가지 않게 한다(렌더하지 않는 응답 포함).
- [x] 위 경계의 자동 검사(수정 전 실패 확인), 영향 문서(D13·data-model·interface-design) 갱신.
- [x] 원격 Linux make lint/test/integration, FullOps lint.py --from c3b5192, DEV 보고(`.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LOWS-DEV-test/report.md`), packet-outcomes, finish, 커밋, worker_done.

O-LG-06의 (1) 쿠키 없는 저장값 회원의 429 한국어, (2) Retry-After 없음, (3) `?lang=` 한 응답 적용은 designer 결정·설계 허용 범위라 바꾸지 않는다. O-CT-03·V-AP-01·O-LG-01·02·04도 현행 유지한다.

소유 파일: `internal/relay/` 친구·언어 관련 코드와 테스트, 기술 문서, own 기록·DEV 보고. 예상 100–250줄.

## 완료 기준과 검증

다섯 항목이 자동 검사로 확인되고 기존 회귀 없음. Linux make lint/test/integration EXIT=0, lint ERROR 0. 이후 coor가 diff 확인과 필요 시 delta 리뷰·tester 확인을 조정한다.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 문구 lang 속성만, 레이아웃 변경 없음.
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(임시 worktree).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 해당 없음(시각 변화 없음).

## 갱신할 산출물

D13과 실제 변경된 기술 문서.

탐색 패킷의 task/role/attempt·기준 SHA와 direct_edit/impact_check/document_read/document_update 경로·구간·근거·partial/unknown을 연결한다. 명시적 없음과 미확인을 구분한다. 필수·충돌·주의 후보와 budget 잔여 확인 목록을 보존한다.

## 기대 산출물

코드·테스트, 기술 문서, DEV 보고, 완료 보고.

## 제약·협업·후속

하위 작업은 [위임 기준](../rules/delegation.md)을 따른다. 미추적 `videos/service-intro/`·`qa-reports/SAR-MOTION-INTRO-001-DEV-test/` 하위 미추적 폴더는 건드리지 않는다. workspace 밖에 파일을 만들지 않는다(private 도구 디렉터리 제외). 다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-UX-ISSUES-002-LOWS-DEV-packet.json` / SHA `aa276b414fceba0aca37e69268322583db1a35a1` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` (document_read) · 줄 6, 10, 14, 26, 40, 64, 76, 161 · inferred · 필수
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 6, 6, 14, 14, 31, 32, 37, 51, 51 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/designer.md` (document_read) · 줄 6, 39, 41 · inferred
- `.fullops-squad/contexts/dev.md` (document_read) · 줄 6, 58, 60, 62, 64 · inferred
- `.fullops-squad/contexts/tester.md` (document_read) · 줄 6, 38, 46, 50 · inferred
- `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md` (document_read) · 줄 86, 111 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 7, 14, 20, 32 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 전체/미확인 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 7, 32, 38, 47, 49, 53, 55, 66, 70, 72, 76 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 15, 21, 42, 59, 65 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-DEV-test/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-TESTER-test/report.md` (document_read) · 줄 6, 10, 25 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-VISUAL/report.md` (document_read) · 줄 6, 10, 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 6, 10, 16, 39, 71, 98 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-DEV-test/report.md` (document_read) · 줄 2, 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-FIX-REVIEW-review/report.md` (document_read) · 줄 2, 9, 15, 21, 88 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-TESTER-test/report.md` (document_read) · 줄 6, 10, 92 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-CONNECT-VISUAL/report.md` (document_read) · 줄 6, 10, 12, 21, 75 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN-REVIEW-review/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN/report.md` (document_read) · 줄 6, 10, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-HOSTS-DESIGN/report.md` (document_read) · 줄 6, 10, 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 15, 21, 65, 84, 90 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-DEV-test/report.md` (document_read) · 줄 2, 6, 10, 20, 27, 66 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-TESTER-test/report.md` (document_read) · 줄 6, 10, 30, 80 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-VISUAL/report.md` (document_read) · 줄 6, 10, 23, 56, 107 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 14, 26, 26, 26, 44, 47, 61, 62, 63 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-DEV-test/report.md` (document_read) · 줄 2, 6, 10, 18, 96, 102, 108, 108 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-TESTER-test/report.md` (document_read) · 줄 6, 10, 73 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-VISUAL/report.md` (document_read) · 줄 6, 10, 18, 89, 111 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-CODE-REVIEW-review/report.md` (document_read) · 줄 2, 9, 15, 24, 24, 58, 59, 60 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-FIX-test/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-DEV-test/report.md` (document_read) · 줄 6 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-RETEST-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-TESTER-test/report.md` (document_read) · 줄 6, 10 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-APPROVAL-TESTER.md` (document_read) · 줄 2, 6, 10, 12 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-CONNECT-TESTER.md` (document_read) · 줄 2, 6, 10, 12 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-LANG-TESTER.md` (document_read) · 줄 2, 6, 10, 12 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-UX-ISSUES-002-NAMES-TESTER.md` (document_read) · 줄 2, 6, 10, 12, 46 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-HUMAN-PEER-001-DEV.md` (document_read) · 줄 33 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-DESIGN.md` (document_read) · 줄 6, 10, 14, 38, 42 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-NAMES-DEV.md` (document_read) · 줄 6, 7, 16, 30 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-NAMES-TESTER.md` (document_read) · 줄 6, 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-UX-ISSUES-002-RECEIPT-DEV.md` (document_read) · 줄 6, 20 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 7, 13, 15, 19, 21, 25, 27, 29, 31 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 7, 35, 45, 57, 67 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_dev.md` (document_read) · 줄 136, 137, 140, 144, 152, 156, 185, 199 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_designer.md` (document_read) · 줄 69, 90 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_ops.md` (document_read) · 줄 123, 207, 317, 318, 319, 320, 321, 396 · inferred
- `.fullops-squad/handovers/logs/2026-10-18_to_dev.md` (document_read) · 줄 198, 199 · inferred
- `.fullops-squad/handovers/logs/2026-10-19_to_designer.md` (document_read) · 줄 9, 10, 13, 17, 25, 42, 43, 48, 52, 61, 64, 92, 94, 137, 138, 139, 140, 141, 142, 143, 158, 159, 203, 206, 209, 210, 213, 217, 225, 244, 254, 258, 263, 266, 281, 300, 302, 340, 341, 342, 343, 344, 345, 346, 347, 348, 362, 363, 392, 395 · inferred
- `docs/system-design.md` (document_read) · 줄 70, 98 · inferred
- `internal/relay/connect_guide_test.go` (impact_check) · 줄 15, 16 · inferred
- `internal/relay/friend_rerequest_test.go` (impact_check) · 줄 169 · inferred
- `internal/relay/gate_summary.go` (impact_check) · 줄 80 · inferred
- `internal/relay/human.go` (direct_edit) · 줄 32, 43, 48, 61, 248, 373 · inferred
- `internal/relay/human_http.go` (direct_edit) · 줄 23, 50, 80, 316, 329 · inferred
- `internal/relay/human_names_integration_test.go` (impact_check) · 줄 96 · inferred
- `internal/relay/human_names_test.go` (impact_check) · 줄 27, 40 · inferred
- `internal/relay/i18n_test.go` (impact_check) · 줄 73, 316, 319 · inferred
- `internal/relay/member.go` (direct_edit) · 줄 전체/미확인 · unknown
- `internal/relay/member_agents.go` (impact_check) · 줄 77, 100, 215, 222, 232, 238, 268, 268, 378, 378 · inferred
- `internal/relay/receipt_context.go` (impact_check) · 줄 104, 107 · inferred
- `internal/relay/receipt_context_test.go` (impact_check) · 줄 83 · inferred
미확인 25건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.

### 완료 보고 (2026-10-20, dev)

- 브랜치 / SHA: `fullops/dev` / 제품·문서 `e9da442e9b3d111c2073b9edd0c9d6f38e54713b`, 그 뒤 DEV 보고·packet-outcomes·dev 컨텍스트·완료 기록 커밋.
- 변경 이유: F-R1 수신 한도 보류 행에 이전 방향·끝 4자(`Friend.HeldFrom`·`HeldTail`, 생략 가능)를 보존해 받은 사람 화면이 재요청 전후 같다. F-R2 받은 사람의 보류 행 block이 `blocked-pending`을 유지해 요청자에게 대기로 보인다. F-R3 같은 코드의 대리 행이 대기 중이면 쌍 행이 있어도 재요청은 무동작이다. O-LG-05 오류 문장의 다른 언어 이름에 `lang` span. 내부 표식 `Lang-Saved`·`Lang-Unsaved`는 `boundedHTTP`의 `langMarks` 래퍼가 모든 응답에서 지운다.
- 검토 필요: `capacity.go` 응답 래퍼(`Unwrap`·`MaxBytesReader`에 원래 writer 사용), `requestFriend`의 legacy 검사 조건 변경, 상태 JSON 새 생략 가능 필드 2개. migration·새 경로·의존성 없음.
- 검증한 것: 새 검사 5건의 수정 전 실패(로컬, 원격 기준 파일 복원 사본 4건 FAIL), 원격 Linux 격리 clone `make install`·`make lint`·`make test`·integration(-race) EXIT=0, FullOps lint.py `--from c3b5192` 깨끗한 임시 detached worktree ERROR 0(WARNING SIZE-001 4, 실행 불가 2=Windows make 없음), worktree는 force 없이 제거. DEP-001 없음(`go.mod`·`go.sum`·`package.json` 변경 없음, 새 import는 표준 `fmt`).
- 못 한 것: 스크린리더 실제 낭독, 실제 브라우저 헤더 관찰, 독립 코드 리뷰·tester 확인(coor 조정).
- 산출물: [DEV 보고](../docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LOWS-DEV-test/report.md), [packet-outcomes](../docs/evaluations/jev/SAR-UX-ISSUES-002-LOWS-DEV-packet-outcomes.json), D13·data-model·interface-design·module-design 갱신.
- 남은 일: coor diff 확인·필요 시 delta 리뷰·tester 확인·main 통합.
