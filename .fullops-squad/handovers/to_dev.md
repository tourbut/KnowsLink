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
- 상태: ready
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

- [ ] F-R1: 수신 한도 보류 행이 받은 사람에게 이전 모습 그대로 보이도록 이전 방향(Inviter·CodeTail)을 보존한다(또는 동등한 방법). 받은 사람 화면이 재요청 전후 같음을 테스트로 보인다.
- [ ] F-R2: 받은 사람이 보류 행에 block POST를 해도 요청자의 살아 있는 대기 요청이 '해제'로 바뀌지 않게(보이는 대기 요청 차단과 같은 결과) 한다.
- [ ] F-R3: 기존 쌍 행+열린 대리 행 데이터에서 같은 코드 재요청이 열린 행 2개를 만들지 않게 한다(없는 코드와 같은 무동작).
- [ ] O-LG-05: 영어 오류 'Choose 한국어 or English.'의 '한국어'(및 한국어판의 'English')에 부분 언어 lang 속성.
- [ ] 언어 내부 표식 헤더가 어떤 응답으로도 클라이언트에 나가지 않게 한다(렌더하지 않는 응답 포함).
- [ ] 위 경계의 자동 검사(수정 전 실패 확인), 영향 문서(D13·data-model·interface-design) 갱신.
- [ ] 원격 Linux make lint/test/integration, FullOps lint.py --from c3b5192, DEV 보고(`.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LOWS-DEV-test/report.md`), packet-outcomes, finish, 커밋, worker_done.

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

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
