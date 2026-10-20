---
title: SAR-CI-002-DEV — CI 파생 로고(파비콘 단순형·반전판·한글 lockup)를 만든다
status: draft
updated: 2026-10-20
owner: dev
tasks: [SAR-CI-002-DEV]
summary: CI 파생 로고(파비콘 단순형·반전판·한글 lockup)를 만든다
attempt: 71e338c013c74ea7bd9f3535dd4c680f
base: 760dc1f4bad73df3edb03cd02f0e5cb7075d0341
subagent_level: standard
test_level: lite
---

# SAR-CI-002-DEV — CI 파생 로고(파비콘 단순형·반전판·한글 lockup)를 만든다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: `docs/brand/`의 SVG 파생 에셋 3종 추가와 brand README 갱신, 렌더 확인. 제품 코드·웹 파비콘 연결·운영 변경은 하지 않는다(파비콘을 제품에 붙이는 일은 후속 과제).
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

사용자 결정 D008(2026-10-20)로 'CI 파생 로고'를 진행한다. 원본은 `docs/brand/knowslink-lockup.svg`·`knowslink-symbol.svg`(B3안, SAR-CI-001). `docs/brand/README.md` '알려진 한계'가 세 가지를 남겼다: 32px 이하 파비콘 단순형 없음, 어두운 배경용 반전판 없음, 한글 표기 버전 없음.

## 적용 기준과 예외

fullops-common-0.3.3, 문서 규칙. `docs/brand/README.md`의 색(잉크 `#182B31`, 청록 `#005A62`만)·도형 의미·기하 규칙을 지킨다. 원본 SVG는 바꾸지 않는다. 손으로 쓴 SVG가 원본이며 래스터 이미지 생성 도구는 쓰지 않는다. 한글 글자는 윤곽선으로 넣어 폰트 없이 렌더되게 한다(Pretendard 등 저장소에 이미 쓴 글꼴 출처·라이선스를 확인하고 README에 적는다).

## 먼저 읽을 문서

`docs/brand/README.md`, 두 원본 SVG, PLANS.md '## CI 로고·보드 정리 — 2026-10-18' 절.

## 해야 할 일과 파일 소유권

- [ ] `knowslink-favicon.svg`: 32px·16px에서 읽히는 단순형(선 틈·점을 줄이거나 굵기를 키우되 두 칸·청록 겹침 의미 유지). 16·32·48px 렌더 확인.
- [ ] `knowslink-lockup-inverse.svg`(필요하면 `knowslink-symbol-inverse.svg`): 어두운 배경용 반전판. 잉크색 대체 색과 청록 대비를 정하고 WCAG 비텍스트 대비 3:1 이상을 수치로 확인.
- [ ] `knowslink-lockup-ko.svg`: 한글 '노우스링크' 표기 lockup(윤곽선).
- [ ] `docs/brand/README.md` 에셋 표·사용 규칙·알려진 한계 갱신(front matter).
- [ ] 렌더 확인 결과(PNG 캡처·대비 수치)를 `.fullops-squad/docs/evaluations/qa-reports/SAR-CI-002-DEV-test/`에 남긴다.
- [ ] FullOps lint.py(깨끗한 임시 detached worktree, force 없이 정리), packet-outcomes, 완료 보고, work.py finish, 커밋, worker_done.

소유 파일: `docs/brand/` 새 SVG·README, 위 QA 폴더, own 기록.

## 완료 기준과 검증

세 파생 에셋이 원본 기하·색 규칙과 맞고 렌더·대비 근거가 있다. lint ERROR 0. 시각 수락은 사용자 또는 designer가 한다(coor 조정).

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: `docs/brand/README.md`.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 원본 SVG 기하를 재사용, 파비콘 단순화는 README에 근거 기록.
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(임시 worktree).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 밝은·어두운 배경 렌더.

## 갱신할 산출물

없음(brand README).

탐색 패킷의 task/role/attempt·기준 SHA와 direct_edit/impact_check/document_read/document_update 경로·구간·근거·partial/unknown을 연결한다. 명시적 없음과 미확인을 구분한다. 필수·충돌·주의 후보와 budget 잔여 확인 목록을 보존한다.

## 기대 산출물

SVG 3–4개, 갱신된 brand README, 렌더·대비 근거, 완료 보고.

## 제약·협업·후속

미추적 `videos/service-intro/`·`qa-reports/SAR-MOTION-INTRO-001-DEV-test/` 하위 미추적 폴더는 건드리지 않는다. workspace 밖에 파일을 만들지 않는다(브라우저 프로필 등은 QA 폴더 아래). 다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
