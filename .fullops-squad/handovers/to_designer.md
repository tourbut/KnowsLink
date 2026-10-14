---
title: SAR-LANDING-GUIDE-001-VISUAL — 고정 구현 후보의 공개 첫페이지를 실제 브라우저에서 시각 검수하고 D04 수락 여부를 기록한다
status: draft
updated: 2026-10-14
owner: designer
tasks: [SAR-LANDING-GUIDE-001-VISUAL]
summary: 고정 구현 후보의 공개 첫페이지를 실제 브라우저에서 시각 검수하고 D04 수락 여부를 기록한다
attempt: 93adb3b5b2124aed99b9c02f45e39281
base: 4795d619513fde4fbcd4a55b94ec91b056ebbaf0
subagent_level: standard
test_level: lite
---

# SAR-LANDING-GUIDE-001-VISUAL — 고정 구현 후보의 공개 첫페이지 직접 시각 검수

- 작성일: 2026-10-14 / From coor → designer / 상태: queued.
- 선행: DEV 완료 고정 SHA와 실제 fixture 제공 후 coor가 기준·route/find/context/packet을 갱신해 dispatch한다. 현재 base는 예비 후보이며 아직 실행하지 않는다.
- 실행: codex/gpt-6-astra/high, 사용자 최신 지정. test_level=lite/subagent_level=standard.
- 담당: repo0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/designer / fullops/designer.
- 복귀: coor / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_38ae16dc-0ca4-4beb-ab00-d1e1b641f061 / run_5534969da2d0. 실제 새 task/dispatch preamble을 사용한다.
- 병합 책임: coor / main·origin/main. 제품 코드 변경·배포·실계정·실메일·외부 발송은 범위 밖이다.

## 적용 기준과 예외

AGENTS.md, FULLOPS.md, fullops-common-0.3.3 공통 코딩·테스트·보안, project.md, document-writing.md, fullops-frontend/delivery.md를 적용한다. frontend.json 기존 주 스킬 frontend-design과 D04 공개 첫페이지 정본을 사용한다. 고정 실행 tuple은 사용자 지정이 우선하며 새 도구 설치·유료 키·이미지 에셋은 추가하지 않는다. 과거 최소 로그인 절은 새 공개 첫페이지 정본에서 대체한 범위만 달라진다. 실제 QA와 미실행을 구분하고 기존 실패를 보존한다.

## 먼저 읽을 문서

- .fullops-squad/docs/design-docs/mockups/service-member-ui.md — 공개 첫페이지 정본.
- .fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md — LG-01~08.
- .fullops-squad/docs/agents/frontend.json.
- .fullops-squad/docs/operations/user-guide.md.
- .fullops-squad/docs/exec-plans/phases/SAR-LANDING-GUIDE-001.md.
- README.md — 실제 설치·수신 조건.
- internal/relay/landing_pages.go, internal/relay/member_pages.go — 고정 후보만 읽는다.
- .fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DEV-test/report.md — 실제 DEV evidence와 한계.
- .fullops-squad/handovers/to_tester.md — 검사 담당 분리.

## 해야 할 일과 소유권

- [ ] 실제 고정 후보의 loopback 페이지와 SHA를 대조하고 desktop1440/mobile390/320을 직접 캡처·열어 확인한다.
- [ ] hero·관계 설명·시작 순서·Grok/Dots 펼침·가상 예시·로그인/오류/none을 D04와 비교한다. 기존 회원 CSS에 영향을 주지 않았는지 대표 화면을 확인한다.
- [ ] 텍스트 위계·여백·한글 줄바꿈·긴 명령 overflow·44px 대상·focus/대비를 판정한다. 키보드/200% 실제 조작은 tester 결과와 연결하며 사진만으로 실행 PASS를 주장하지 않는다.
- [ ] 결함은 재현·viewport·증거·심각도와 수정 담당 DEV를 기록한다. 제품 코드를 수정하지 않는다.
- [ ] own phase/QA/manifest/context/packet outcomes를 작성하고 work.py finish 및 실제 worker_done으로 고정 기록 SHA를 회신한다.

파일 소유권은 own 인박스/context, docs/exec-plans/phases/SAR-LANDING-GUIDE-001-VISUAL.md 및 docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/와 own 탐색 기록이다. D04 변경·다른 역할 결과·제품 코드는 읽기 전용이다.

## 실행·완료 기준

coor 제공 고유 fixture origin에서 GET·앵커·native details 펼침과 viewport만 허용한다. 합성 오류는 제공 fixture의 테스트 경로만 사용한다. 실제 Google/운영/메일/설치 실행·계정 생성·발송은 금지한다. Orca CLI와 frontend 기존 도구를 사용하고 페이지별 browserPageId로 tester와 충돌을 피한다. 최대40개 의도적 조작, 20분/시나리오로 제한한다. capture 기존 evidence는 exactSHA·viewport·상태 확인 후 재사용하고 실제 이미지를 열어 판정한다. 빠진 상태만 새 캡처한다. 영상은 필요 없다.

보고서에 고정 제품 SHA, author와 다른 실제 session, 실행기·origin·page·viewport·zoom·파일 경로, 수락/FAIL/not_run, D04 대비 차이와 미해결 C/H를 기록한다. static wireframe·브라우저 설치·CLI 성공으로 실제 시각 PASS를 대신하지 않는다. 디자인 미해결 critical/high는 수락을 차단한다. 독립 tester 기능 QA·코드 리뷰·운영 배포 성공을 이번 시각 검수로 주장하지 않는다.

Windows 삭제·이동 정리는 실행하지 않는다. 검증 자료를 보존하고 경로·fixture 유지 담당을 남긴다. PLANS/board와 main 통합은 coor 담당이다.

## 완료 보고

미착수. DEV 고정 최종 SHA·fixture와 실제 dispatch를 기다린다.
