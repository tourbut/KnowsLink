---
title: SAR-SCHEDULE-001-DESIGN — 실제 일정 API 연동의 제품 범위 결정 선택지와 규칙 초안을 정리한다
status: draft
updated: 2026-10-20
owner: designer
tasks: [SAR-SCHEDULE-001-DESIGN]
summary: 실제 일정 API 연동의 제품 범위 결정 선택지와 규칙 초안을 정리한다
attempt: eaa1057e93ce4166a4cff3bd9b02d93a
base: 760dc1f4bad73df3edb03cd02f0e5cb7075d0341
subagent_level: standard
test_level: lite
---

# SAR-SCHEDULE-001-DESIGN — 실제 일정 API 연동의 제품 범위 결정 선택지와 규칙 초안을 정리한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 실제 일정 API 연동의 제품 기획(목표·사용자 흐름·범위·동의·실패 규칙·완료 조건)과 사용자 결정이 필요한 선택지 정리, 제품 정본 문서 작성, designer 브랜치 커밋. 외부 캘린더 계정 연결·실데이터 조회·API 키 발급·제품 코드 수정·운영 변경은 하지 않는다.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

사용자 결정 D008(2026-10-20)로 '실제 일정 API 연동'을 진행한다. 현재 제품은 에이전트끼리 `schedule.query`·`schedule.commit` 요청을 서명·전달하고 사람 승인 게이트를 거치지만, 승인 뒤에도 실제 캘린더 조회·일정 생성은 실행하지 않는다(정책상 disclosure·executable false). 보드의 '실제 어댑터·일정 조회' 단계가 이 후속을 남겨 두었다. 제품 범위(어떤 캘린더, 읽기만인지 쓰기까지인지, 동의·공개 범위)는 정해지지 않았다.

- 기존 정본: `.fullops-squad/docs/planning/product-specs/SAR-MVP.md`(첫 기능 범위·실제 calendar 제외 근거), `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md`, `.fullops-squad/docs/design-docs/architecture.md`(gate approve 뒤 executable/disclosure false), `interface-design.md`(gate 화면·요약), `SAR-UX-ISSUES-002.md` U5(승인 요약·'일정 실행을 허용하지 않습니다' 고정 문구).
- 로그인은 Google이다(운영 1dd6800). Google 연결 범위 기록: `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md`, `SAR-GOOGLE-CONNECT-002-DEV.md`.

## 적용 기준과 예외

fullops-common-0.3.3, 문서 규칙 `.fullops-squad/docs/agents/document-writing.md`, D01/D02 제품 정본. 기준 main(이 지시서 커밋). 외부 API의 실제 기능·권한 범위는 Context7 또는 공식 문서로 확인하고 출처를 남긴다. 확인하지 못한 API 동작을 확정 사실로 쓰지 않는다.

## 먼저 읽을 문서

위 정본들, `.fullops-squad/contexts/designer.md`, PLANS.md의 '실제 어댑터·일정 조회' 관련 절.

## 해야 할 일과 파일 소유권

- [ ] 사용자 결정이 필요한 선택지를 근거와 함께 정리: 대상 캘린더(예: Google Calendar만/그 밖), 기능 범위(빈 시간 조회만 / 일정 생성까지), 공개 범위(바쁨·한가함만 / 제목 등), 계정 연결·동의·철회 방식, 승인 게이트와의 관계(승인 = 그 요청 한 번의 실행 허용인지).
- [ ] 각 선택지의 위험·비용·운영 영향(외부 OAuth scope 심사 필요 여부, 토큰 보관, 실패·재시도, 감사 기록)을 짧게 비교하고 designer 추천안을 적는다.
- [ ] 추천안 기준의 제품 규칙 초안: 사용자 흐름, 화면 문구 원칙, 실패·만료·철회, 완료 조건(DEV/QA가 쓸 수락 기준), 범위 밖.
- [ ] 정본 `.fullops-squad/docs/planning/product-specs/SAR-SCHEDULE-001.md`(D02 성격) 작성, 결정 근거 `.fullops-squad/docs/evaluations/qa-reports/SAR-SCHEDULE-001-DESIGN/report.md`.
- [ ] 사용자 결정이 필요한 항목은 완료 보고에 질문 목록(선택지·추천)으로 남긴다. 결정 전까지 DEV 인계는 하지 않는다.
- [ ] 완료 보고 → work.py finish → 커밋 → worker_done.

소유 파일: 위 두 문서, 필요 시 D01/D02 인덱스 연결, own route/inbox/log.

## 완료 기준과 검증

coor가 사용자에게 바로 물을 수 있는 결정 목록(각 2–4 선택지·추천·영향)과 추천안 기준 제품 규칙 초안이 있다. 외부 API 근거 출처가 있다. 문서 lint ERROR 0.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04(화면 원칙만 참조).
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 해당 없음(기획).
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 해당 없음.

## 갱신할 산출물

D01, D02(새 제품 정본 연결).

## 기대 산출물

SAR-SCHEDULE-001.md, 결정 근거 report, 완료 보고(결정 질문 목록 포함).

## 제약·협업·후속

기술 구현 방법은 정하지 않는다(DEV 담당). 다른 역할 인박스·PLANS·board·GitHub·제품 코드는 수정하지 않는다. workspace 밖에 파일을 만들지 않는다.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
