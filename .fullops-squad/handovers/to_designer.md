---
title: SAR-UX-ISSUES-002-HOSTS-DESIGN — U8 호스트 지원 표를 최신 사용자 범위와 Claude Code 세션 연결로 갱신하고 O-NM-03을 결정한다
status: draft
updated: 2026-10-20
owner: designer
tasks: [SAR-UX-ISSUES-002-HOSTS-DESIGN]
summary: U8 호스트 지원 표를 최신 사용자 범위와 Claude Code 세션 연결로 갱신하고 O-NM-03을 결정한다
attempt: b5c7601c40374611a78963fdb76dab25
base: 79864592804e595810706f6550ff255ee65407ac
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-HOSTS-DESIGN — U8 호스트 지원 표를 최신 사용자 범위와 Claude Code 세션 연결로 갱신하고 O-NM-03을 결정한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → designer
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 제품 정본(SAR-UX-ISSUES-002.md U8·U4 관련 절) 갱신과 근거 기록, designer 브랜치 커밋. 호스트 설치·실계정 연결·외부 발송·운영 변경·제품 코드 수정은 하지 않는다.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

잔여 UI 1–3단위가 main 7986459에 통합됐다. 다음은 4단위 U8(회원 연결 경로·재개, 나머지 U7)이다. 그런데 U8 '근거에 따른 호스트 표시' 표(2026-10-15)가 이후 사실과 맞지 않는다.

- 사용자 호스트 범위(PLANS.md 2026-10-15 '최신 사용자 호스트 범위' 줄): Meta Muse는 대상에서 제외, OpenClaw는 사용자 지시로 보류, Hermes는 사용자가 지정한 서버의 Nous Research Hermes Agent 0.21.3을 확인(설치·왕복 수락 아님).
- 2026-10-17 Claude Code 세션 연결이 구현·수락·main 통합됐다: `/knowslink-session on|off|status [--take]`, 도착 toast, 승인 대화상자(SAR-SESSION-PEER-001/002, 정본 `.fullops-squad/docs/design-docs/SAR-SESSION-PEER-001.md`, 기술 `.fullops-squad/docs/design-docs/module-design.md`·`architecture.md`, QA `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-00*-*`). Codex는 V02 미실행.
- 운영 배포 1dd6800(2026-10-17)에 세션 API가 반영됐다(PLANS '운영 배포 완료·통합').
- 2단위 관찰 O-NM-03(후속 U4 개선 후보): 보낸 친구 요청 행이 서로 같아 요청자가 어느 코드로 보냈는지 알 수 없다. 이전 판정은 '요청자가 입력한 코드의 일부(예: 끝 4자) 표시 검토'(`.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-NAMES-VISUAL/report.md` 'O-NM-03~07 분류').

## 적용 기준과 예외

fullops-common-0.3.3, 문서 규칙, 제품 정본 SAR-UX-ISSUES-002.md(U8·U4), D02 SAR-HUMAN-PEER.md, D04. 기준 main 7986459. 근거가 없는 호스트를 지원 완료로 올리지 않는 기존 원칙을 유지한다.

## 먼저 읽을 문서

SAR-UX-ISSUES-002.md U4·U8, SAR-SESSION-PEER-001.md, PLANS.md의 2026-10-15 호스트 범위 줄과 2026-10-17 세션 통합 줄, README.md·adapters/README.md의 현재 설치 안내, `.fullops-squad/contexts/designer.md`.

## 해야 할 일과 파일 소유권

- [ ] U8 호스트 표 갱신: 대상 호스트 목록(Meta Muse 제외 명시, OpenClaw 보류 표시 방식, Hermes 상태, Claude Code 세션 연결 추가)과 각 근거 날짜·SHA·화면의 지원 상태 문구·다음 행동. 소스/패키지 근거와 실제 앱 설치·왕복 근거를 구분한다. Codex 세션은 미검증으로 둔다.
- [ ] '어디에서 사용할까요?' 선택지 목록을 위 표와 맞춘다(웹·Grok·Dots·Claude Code·지원 확인 필요).
- [ ] O-NM-03 결정: 보낸 요청 행에 요청자 자기 입력 코드의 일부를 보일지, 보인다면 형식(자리 수·마스킹)과 상대 존재 비노출 경계. 보이지 않기로 하면 근거.
- [ ] 변경을 SAR-UX-ISSUES-002.md에 반영(front matter 갱신), 결정 근거를 `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-HOSTS-DESIGN/report.md`에 짧게 기록.
- [ ] 완료 보고 → work.py finish → 커밋 → worker_done(DEV 4단위 구현 기준 요약 포함).

소유 파일: SAR-UX-ISSUES-002.md, 위 report, own route/inbox/log.

## 완료 기준과 검증

DEV가 추측 없이 4단위를 구현할 수 있을 만큼 호스트 표·선택지·O-NM-03이 확정됐다. 문서 lint ERROR 0.

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04, SAR-UX-ISSUES-002.md.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 해당 없음(제품 기준 문서).
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 해당 없음.

## 갱신할 산출물

SAR-UX-ISSUES-002.md(필요 시 D02).

## 기대 산출물

갱신된 정본, 결정 근거 report, 완료 보고.

## 제약·협업·후속

기술 구현 방법은 정하지 않는다(DEV 담당). 다른 역할 인박스·PLANS·board·GitHub·제품 코드는 수정하지 않는다.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
