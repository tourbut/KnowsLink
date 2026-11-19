---
title: SAR-SERVICE-UI-001-DESIGN
status: draft
updated: 2026-10-11
owner: designer
tasks: [SAR-SERVICE-UI-001-DESIGN]
summary: 회원 웹 UX·비주얼 명세 작성의 확인·결정·한계·후속을 기록한다
---

# SAR-SERVICE-UI-001-DESIGN — 실행 기록

## 범위와 판단

- 기준 SHA: `7a8c2b9194edb77b52db41a14b2325af743f0d6e`. 기록 시작 HEAD: `0de25a5`.
- 범위: link.knowslog.com 회원 웹(로그인 → 내 에이전트 → 연결 관리)의 상용 UX·비주얼 명세와 DEV 인계. 제품 코드·dev 인박스·서버 상태·자격 변경 0건.
- 정본: [회원 웹 명세](../../design-docs/mockups/service-member-ui.md)(D04). D04 인덱스 원천에 이 문서를 추가했고 기존 [SAR-MVP-001-UI](../../design-docs/mockups/SAR-MVP-001-UI.md)에 링크 절을 추가했다.
- 공통 기준: fullops-common-0.3.3, FULLOPS/project, 문서 작성 규칙, designer 컨텍스트, 이전 UI 검수(F-UI-01–04, F-UI-MSG-01).

## 실제 확인

- Orca 내장 브라우저로 이미 로그인된 회원 홈 `/home`을 1440×900·390×844에서 직접 열었다. 에이전트 4개 상태다. 가로 overflow 없음, 문서 높이 6394px/7986px, 버튼 높이 52px 이상.
- 로그인 `/`는 로그아웃하지 않고 운영 응답 HTML을 받아 로컬 파일로 렌더했다. 운영은 Google 로그인만 노출한다. CSP 헤더를 확인했다.
- 캡처 PNG는 임시 폴더에서만 열람했고 Git에 넣지 않았다. 이메일·식별자·키 값은 문서에 쓰지 않았다.
- 대비 비율은 WCAG 공식으로 계산했다(명세 8.1).
- 세션 쿠키·로그인 상태·서버 데이터는 변경하지 않았다. 열었던 브라우저 탭은 닫았다.

## 결정 요약

- 홈은 "에이전트 목록 우선·행 접힘·위험 구역 분리·확인 필요 띠" 구조다. 상태는 서버가 아는 사실만 표시한다.
- 한국어 본문 줄바꿈을 `keep-all`+`break-word`로 바꾸고 식별자·지문에만 강제 줄바꿈을 쓴다.
- 새 CSS/JS 라이브러리·외부 폰트·이미지 없이 가능한 기준이다. CSP는 완화하지 않는다.
- CSRF·owner 분리·5분 재인증·파괴 동작 확인 방식은 변경하지 않는다. 별도 확인 모달은 추가하지 않는다.

## 검증과 한계

- 문서 변경만이다. 제품 코드 변경이 없어 `make lint`/`make test`는 해당 없음이다. 문서 front matter는 `deliverables.py --stamp`로 기록했다. 디자인 lint는 미구성(project.md)이라 해당 없음이다.
- 직접 보지 못한 것: 이메일 코드·오류 화면, 연결 관리·기기 승인·receipt·gate·refusal의 렌더 결과(템플릿 코드로만 확인), 에이전트 5개·긴 관계 목록, 다크 모드.
- 이번 결과는 디자인 인계 문서이며 UI 최종 수락이 아니다. 구현 후 직접 시각 검수는 DEV가 고정 SHA·검수 환경을 제공한 뒤 수행한다.
- `frontend-design` 스킬은 이 환경에서 찾지 못했다. 직접 관찰과 CSP·템플릿 제약 검토로 대신했다.
- 하위 에이전트는 쓰지 않았다(standard 레벨의 선택 예산 미사용).

## 후속

- DEV: 명세 12절의 R-01–R-11과 완료 기준으로 구현 계획을 세운다. dev 인박스는 coor가 정규 지시서로 배정한다.
- designer: 고정 SHA에서 1440/390 직접 시각 검수.
