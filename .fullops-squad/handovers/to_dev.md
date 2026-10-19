---
title: SAR-UX-ISSUES-002-APPROVAL-DEV — 검증된 승인 요청을 사람이 판단할 수 있는 요약과 기술 상세로 나눈다(U5+U7)
status: draft
updated: 2026-10-19
owner: dev
tasks: [SAR-UX-ISSUES-002-APPROVAL-DEV]
summary: 검증된 승인 요청을 사람이 판단할 수 있는 요약과 기술 상세로 나눈다(U5+U7)
attempt: b38259160bfe4a0b9656ae49784d5a45
base: 688305054dc4f3570f0c9f0d1e7917c35a19db3a
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-APPROVAL-DEV — 검증된 승인 요청을 사람이 판단할 수 있는 요약과 기술 상세로 나눈다(U5+U7)

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-19
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 아래 U5+해당 U7 구현·검사·기술 문서 갱신과 소규모 후속 2건. 새 DB 스키마·새 공개 API·권한 확대·정책 활성화·실계정/실메일/상대 발송·운영 배포·삭제/이동·force push는 coor에 ask.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

사용자 2026-10-18 지시 "잔여작업 모두 수행해"의 순차 인계 3단위다. 2단위(U2/U4 이름·친구)는 main 6883050에 통합됐다(DEV f799250, 리뷰·tester·designer 수락). 기준은 main 6883050이다.

현재 승인 화면은 `internal/relay/human_pages.go`의 `human-gate` 템플릿과 `internal/relay/http.go`의 `verifiedGate`(541행 부근)·`gateView`(556행 부근)다. 상세는 ID·Intent·typed body·정책을 먼저 보이고, 승인 목록은 parent ID 위주다(legacy-home·member_pages.go). 2단위에서 발신·대상 에이전트 이름 라벨은 이미 들어갔다.

소규모 후속(2단위 designer 분류):
- O-NM-04(low, 1단위 잔여): 익명 429 화면에 빈 radiogroup·빈 `href` 링크와 익명에게 '저장된 입력 복원' 문구가 보인다. 결정: 익명 429는 선택기 없이 재시도 시각과 실제 로그인 경로 링크만, 저장 입력 복원 문구는 로그인 회원에게만, 빈 href 금지(`.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-VISUAL/report.md` 'O-NM-03~07 분류').
- O-NM-08(low): 받은 친구 요청의 거절 안내가 수락 버튼 바로 아래에 붙어 있다. 안내가 거절 버튼과 함께 읽히도록 배치를 정리한다(문구 불변).

## 적용 기준과 예외

- 규칙: fullops-common-0.3.3, `.fullops-squad/project.md`, 문서 규칙. 기준 커밋 6883050.
- 제품 정본: `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` U5·U7·순차 인계 표 3행, D04, D02 SAR-HUMAN-PEER.md.
- frontend.json의 codex host는 사용자 Codex 제외 지시로 적용하지 않는다. claude/claude-opus-5-5/high(user).
- Windows make 없음: 2단위와 같이 격리 원격 Linux clone(`.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-DEV-test/remote.py` 방식, private 도구 디렉터리는 새 과제용으로 분리)에서 make lint/test/integration.
- 보안 경계: 표시 이름·요약은 읽기 보조다. 서명·ID·digest·관계 세대 결속, 서버의 owner·키·세대·기한 재검사를 그대로 유지한다. hint·자유문을 승인 범위로 올리거나 AI 요약하지 않는다.

## 먼저 읽을 문서

SAR-UX-ISSUES-002.md U5·U7, D04, `internal/relay/http.go`(verifiedGate·gateView·gate 처리), `internal/relay/human_pages.go`(human-gate), `internal/relay/member_pages.go`(승인 목록·legacy), 관련 gate/policy/schema 테스트, `internal/relay/member_receipt.go`(O-NM-04), `.fullops-squad/contexts/dev.md`, 2단위 완료 기록 `.fullops-squad/handovers/logs/2026-10-19_to_dev.md`.

## 해야 할 일과 파일 소유권

- [ ] 짧은 기술 계획(완료 보고에 기록): 지원 요청 유형별 검증된 필드 → 사용자 문구 대응표를 DEV 보고에 남긴다.
- [ ] 승인 목록: 빈 목록 문구, 행마다 요청 주체/상대의 허용된 이름·구분값, 요청 유형, 만료 시각, 현재 상태, '요청 검토'. 원요청 ID는 기술 상세로.
- [ ] 승인 상세 읽기 순서: '누가 누구에게' → '요청한 내용'(검증된 구조화 필드만) → '승인하면 허용되는 범위' → '기한·상태' → 행동('이 요청 승인'·'거절') → 기술 상세(원문 JSON·내부 ID·정책). 고정 설명 문구와 hint '참고 문구 — 승인 근거 아님'. 자유문은 '요청자가 적은 내용'으로 분리, 링크 자동 열기 없음.
- [ ] 만료·처리됨·원문 부재·서명/바인딩 불일치·미지원 유형: 중복 행동 없음, '승인할 내용을 확인할 수 없습니다'와 안전한 복귀. 서버 재검사 유지, 낡은 화면 제출이 새 요청에 적용되지 않음.
- [ ] U7(해당 화면): 고유 title·h1·main·skip, 승인 목록으로 정확한 복귀, status/alert, 반복 버튼 접근성 이름에 대상·구분값.
- [ ] O-NM-04, O-NM-08 반영.
- [ ] U5-A 핵심 실패 경계 자동 검사: 요약/원문 일치, 이름 변경 뒤 같은 원요청 식별, 승인/거절/만료/재접근/중복 클릭, 악성 hint·HTML·바뀐 digest·다른 owner·원문 부재·관계 철회에서 권한 확대 없음. 문구 복제 테스트는 쓰지 않는다. 실제 외부 업무 실행을 켜지 않는다.
- [ ] tester/designer용 격리 fixture(승인 대기·만료·처리됨·원문 부재·악성 hint 요청, 합성 A/B/C, loopback origin, 허용 GET/POST, 수명) 준비·컴파일. 기동은 coor가 한다.
- [ ] 영향받는 D05/D10/D11/D13 갱신, 원격 Linux make lint/test/integration, FullOps lint.py --from 688305054dc4, packet-outcomes, 완료 보고, work.py finish, 커밋, worker_done.

소유 파일: `internal/relay/` 승인·회원 화면/handler/템플릿과 테스트·fixture, 기술 문서, own 기록·DEV 보고(`.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-APPROVAL-DEV-test/`). 예상 300–700줄.

## 완료 기준과 검증

U5-A가 자동 검사 또는 fixture로 확인되고 실패 경계가 기록됐다. 기존 서명·권한·세대 검사 회귀 없음. 원격 Linux make lint/test/integration EXIT=0, FullOps lint ERROR 0. worker 완료는 제품 수락이 아니다(후속 리뷰·tester·designer·통합은 coor).

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04, SAR-UX-ISSUES-002.md '화면 배치와 자기 비평'.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 기존 Go html/template, `human-agent-label`, `dl.kv`, native details. 새 테마·폰트·JS 없음.
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 테마 전환 없음. 1440·390·320 렌더로 긴 이름·긴 자유문·긴 JSON 넘침 0 확인.

## 갱신할 산출물

D10, D13, D12(route 추천)와 실제 변경에 따라 D05·D11. D02/D04 제품 정본은 수정하지 않는다.

탐색 패킷의 task/role/attempt·기준 SHA와 direct_edit/impact_check/document_read/document_update 경로·구간·근거·partial/unknown을 연결한다. 명시적 없음과 미확인을 구분한다. 필수·충돌·주의 후보와 budget 잔여 확인 목록을 보존한다.

## 기대 산출물

제품 코드·테스트·fixture, 기술 문서, DEV 보고(대응표·검사 로그·렌더·fixture·not_run), packet-outcomes, 완료 보고.

## 제약·협업·후속

하위 작업은 [위임 기준](../rules/delegation.md)을 따른다(standard, 최대 2). 미추적 `videos/service-intro/`·`qa-reports/SAR-MOTION-INTRO-001-DEV-test/` 하위 미추적 폴더는 건드리지 않는다. 연결 온보딩(U8)·언어(U3)·O-NM-03은 섞지 않는다. 다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다. 피할 스타일: 상태 카드·통계 배지·그라디언트·알약 버튼·새 색 토큰·가운데 정렬 본문.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
