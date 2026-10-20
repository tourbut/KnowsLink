---
title: "SAR-UX-ISSUES-002-LANG-DEV — 설정 화면과 한국어·English 언어 저장(U3), 헤더 '설정' 이름(U7)을 구현한다"
status: draft
updated: 2026-10-20
owner: dev
tasks: [SAR-UX-ISSUES-002-LANG-DEV]
summary: "설정 화면과 한국어·English 언어 저장(U3), 헤더 '설정' 이름(U7)을 구현한다"
attempt: 37a2e33e9f2c4687b71accf1eb36cd46
base: 131544873193795e77a2886dcdd81e8a2658c0f6
subagent_level: standard
test_level: lite
---

# SAR-UX-ISSUES-002-LANG-DEV — 설정 화면과 한국어·English 언어 저장(U3), 헤더 '설정' 이름(U7)을 구현한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: U3 설정·언어 저장과 지정 웹 화면 번역, U7 헤더 '설정' 이름 변경, 관련 검사·fixture·기술 문서. 회원 언어 선호 저장을 위한 기존 상태 JSON 필드 추가와 비로그인 선호용 쿠키 1개는 승인됨(기존 보안 쿠키 속성 수준 유지). DB migration·새 공개 API·새 의존성(번역 라이브러리 포함)·권한 확대·실계정·운영 배포·삭제/이동·force push는 coor에 ask.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

사용자 지시 "잔여작업 모두 수행해"의 순차 인계 마지막 5단위다. 1–4단위(U6·U2/U4·U5·U8+U7)가 main 1315448에 통합됐다. 정본 순서상 앞 단위의 한국어 문구가 확정된 뒤 한 번 번역한다. #3은 P1이다.

헤더의 '내 계정' 진입을 '설정'으로 바꾸는 U7 항목을 이 단위에서 함께 한다(4단위에서 미룸).

## 적용 기준과 예외

- 규칙: fullops-common-0.3.3, `.fullops-squad/project.md`, 문서 규칙. 기준 커밋 1315448.
- 제품 정본: `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` U3(설정 순서, 언어 선택·저장, 선호 우선순위 표, 비로그인 진입 링크, 번역 범위, lang·KST, 누락 번역 fallback, U3-A)·U7(설정 이름), D04, 공개 첫페이지 규칙(정보 구조·설치 사실 불변).
- 문구·상태 관리 방식은 DEV가 정한다. 새 번역 의존성을 도입하지 않는 쪽을 우선한다(표준 라이브러리·기존 html/template).
- frontend.json codex host는 Codex 제외 지시로 적용하지 않는다. claude/claude-fable-5-1/high(user, 전 템플릿에 걸친 큰 기능).
- Windows make 없음: 이전 단위와 같이 격리 원격 Linux clone(이전 단위 remote.py 방식, private 도구 디렉터리는 새 과제용으로 분리)에서 make lint/test/integration. fixture는 실제 기동까지 확인한다. lint.py는 미추적 사용자 자료 때문에 깨끗한 임시 detached worktree에서 실행하고 끝나면 force 없이 정리한다.
- 번역하지 않는 것: 프로토콜 ID·키·명령·사용자 이름·메시지·검증 원문 JSON, 외부 호스트 UI·Google 화면, README 전체.

## 먼저 읽을 문서

SAR-UX-ISSUES-002.md U3·U7, D04, 공개 첫페이지 정본(D04의 SAR-LANDING-GUIDE-001 절), 회원 identity/세션 코드(`internal/relay/identity.go`, `member.go`), 모든 템플릿 파일(`internal/relay/*_pages.go`, `landing_pages.go`, `remote_auth.go`의 화면), 시각 표시 함수(`clock` 등), 기존 렌더 fixture 테스트, `.fullops-squad/contexts/dev.md`, 4단위 완료 기록(`.fullops-squad/handovers/logs/2026-10-20_to_dev.md`의 SAR-UX-ISSUES-002-CONNECT-DEV 절).

## 해야 할 일과 파일 소유권

- [ ] 짧은 기술 계획(완료 보고에 기록): 문구 관리 구조, 회원 선호 저장 필드, 비로그인 선호 쿠키(이름·속성·수명), 요청별 언어 결정 순서, 누락 번역 감지 방법.
- [ ] 설정 화면: '내 계정'을 '설정'으로(헤더·제목·복귀 링크), 순서 '내 이름'·'내 친구 코드'·'언어'·'로그인과 보안', '한국어 / English'와 '언어 저장 / Save language', 저장 성공·실패 표시(정본 규칙).
- [ ] 언어 결정 우선순위 표 그대로: 회원 저장값 → 비로그인 명시 선호 → 한국어, 계정 교대 시 이전 계정 선호 미복사, 손상 locale은 한국어 fallback. 첫 로그인 시 비로그인 선호는 화면 적용만 하고 저장은 명시 저장으로 확정.
- [ ] 비로그인 진입: 공개 페이지·로그인·동의 화면에 작은 '한국어 / English' 링크. 공개 페이지 정보 구조·설치 사실 불변. 쿠키 저장이 막히면 현재 응답만 적용하고 한계 안내.
- [ ] 번역 범위의 모든 지정 화면(공개/로그인/이름/대화/친구/내 에이전트/설정/legacy/receipt/승인/local 연결/remote 동의·복귀)의 제목·메뉴·라벨·상태·빈 상태·오류·검증·접근성 이름·도움말 영어판. HTML lang 일치, 시각은 KST(UTC+09:00) 유지, 기한·정렬·timestamp 불변.
- [ ] 누락 번역은 한국어 안전 fallback이며 내부 키·빈 버튼 미표시. 지정 화면의 누락이 0임을 자동 검사로 보인다.
- [ ] U3-A 핵심 실패 경계 자동 검사: 저장·실패·새로고침·재로그인·계정 교대·익명 저장 차단·잘못된 locale, 언어 전환이 친구 요청·발송·승인을 실행하거나 입력을 잃지 않음, 원문/ID/서명값 불변, 같은 시각. 문구 복제 테스트는 쓰지 않는다.
- [ ] tester/designer용 격리 fixture(두 언어 회원, 비로그인, 각 지정 화면에 닿는 합성 상태) 준비·실제 기동 확인.
- [ ] 영향받는 D05/D10/D11/D13 갱신, 원격 Linux make lint/test/integration, FullOps lint.py --from 1315448, packet-outcomes, 완료 보고, work.py finish, 커밋, worker_done.

소유 파일: `internal/relay/` 템플릿·화면 handler·identity/세션 선호·문구 관리와 테스트·fixture, 기술 문서, own 기록·DEV 보고(`.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-LANG-DEV-test/`). 예상 규모가 크다(전 템플릿). SIZE 경고는 분할하지 않은 이유와 함께 보고한다.

## 완료 기준과 검증

U3-A가 자동 검사와 실제 기동 fixture로 확인되고, 지정 화면 두 언어의 누락 번역 0, 기존 기능 회귀 없음. 원격 Linux make lint/test/integration EXIT=0, FullOps lint ERROR 0. worker 완료는 제품 수락이 아니다(리뷰 판단·tester·designer·통합은 coor).

### UI 작업의 디자인 기준과 검증

- 디자인 문서 / 테마·토큰 정본 경로: D04, SAR-UX-ISSUES-002.md '화면 배치와 자기 비평'.
- 재사용할 공용 컴포넌트 경로 / variant·허용 변경 / 의도적 예외와 근거: 기존 Go html/template·radio fieldset·native form. 새 테마·폰트·JS 없음.
- 변경 후 다시 실행할 디자인 lint 명령·cwd / DESIGN 경고 처리: FullOps lint.py(레포 루트).
- 테마 전환 후 확인할 화면·컴포넌트와 토큰 전파 / 검증 결과·미실행 사유: 테마 전환 없음. 두 언어로 1440·390·320 렌더해 긴 영문 넘침 0 확인.

## 갱신할 산출물

D10, D13(route 추천)과 실제 변경에 따라 D05·D11. D02/D04 제품 정본은 수정하지 않는다.

탐색 패킷의 task/role/attempt·기준 SHA와 direct_edit/impact_check/document_read/document_update 경로·구간·근거·partial/unknown을 연결한다. 명시적 없음과 미확인을 구분한다. 필수·충돌·주의 후보와 budget 잔여 확인 목록을 보존한다.

## 기대 산출물

제품 코드·테스트·fixture, 기술 문서, DEV 보고(기술 계획·누락 번역 검사·렌더·fixture 기동·not_run), packet-outcomes, 완료 보고.

## 제약·협업·후속

하위 작업은 [위임 기준](../rules/delegation.md)을 따른다(standard, 최대 2; 번역처럼 파일 소유를 나눌 수 있는 작업은 하위 위임 가능하되 부모가 취합·검증). 미추적 `videos/service-intro/`·`qa-reports/SAR-MOTION-INTRO-001-DEV-test/` 하위 미추적 폴더는 건드리지 않는다. workspace 밖에 파일을 만들지 않는다(private 도구 디렉터리 제외). 4단위 delta 리뷰 low F-R1–F-R3 등 다른 후속은 섞지 않는다(coor가 별도 정리 과제로 배정). 다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다. 피할 스타일: 국기 아이콘·드롭다운 JS·상태 카드·알약 버튼·새 색 토큰.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
lint DEP-001이 있으면 파일별 의존성 변경 여부·필요성·표준 라이브러리 대안을 적는다. 의존성 변경이 없으면 그 근거를 적는다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
