---
title: SAR-MOTION-INTRO-001-DEV — KnowsLink 서비스 소개 15초 모션그래픽을 HyperFrames로 제작·렌더·검수한다
status: draft
updated: 2026-10-17
owner: dev
tasks: [SAR-MOTION-INTRO-001-DEV]
summary: KnowsLink 서비스 소개 15초 모션그래픽을 HyperFrames로 제작·렌더·검수한다
attempt: 863239b7e0bd448c84856fc3489e966b
base: c94b84303617564797f8bc21801dcd026811cccb
subagent_level: standard
test_level: lite
---

# SAR-MOTION-INTRO-001-DEV — KnowsLink 서비스 소개 15초 모션그래픽을 HyperFrames로 제작·렌더·검수한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-17
- From / To: coor → dev
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위(사용자 2026-10-17): `videos/service-intro/`에 `hyperframes@latest`를 npm으로 로컬 설치하고 HyperFrames 공식 스킬(heygen-com/hyperframes)을 받는다. 로컬 preview·check·render만 한다.
- 추가 승인이 필요한 행위: 클라우드 렌더, 외부 TTS·생성형 이미지·스톡 자산, usage/feedback/publish 업로드, Remotion 사용, 음원 추가, 전역 도구 설치. 필요하면 `ask`로 멈춘다.
- 담당 repo id / 워크트리 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/dev / fullops/dev
- 병합 책임자 / 기본 브랜치: coor / main·origin/main
- 복귀: coor / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_e3a3e5f6-e163-420f-bba6-51b428f7671f / run·task·dispatch는 실제 preamble 값

## 현재 상황과 확인 근거

- 사용자 요청: "해당 서비스 소개 모션 그래픽을 15초 짜리로 만들자"(2026-10-17). 직전에 운영 배포한 link.knowslog.com의 KnowsLink를 소개한다.
- 확정 규격: 16:9 1920×1080, 30fps, 15초(450프레임), H.264 MP4, 알파 없음, 무음(오디오 스트림 없음), 한글 문구.
- 로컬 환경: Node v24.13.0, FFmpeg 8.0.1(gyan full build)이 PATH에 있다. HyperFrames CLI와 공식 스킬은 설치돼 있지 않다. `videos/` 폴더는 없다.
- 메시지 정본은 공개 첫페이지다. 새 제품 메시지·기능 주장을 만들지 않는다.

## 적용 기준과 예외

공통 기준은 `.fullops-squad/rules/common/README.md` 및 연결 규칙이다. 기준 커밋 c94b843. fullops-motion 스킬의 toolchain.md·delivery.md 절차를 따른다. 예외 없음.

## 먼저 읽을 문서

- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` — `## 공개 첫페이지 정본 — SAR-LANDING-GUIDE-001`의 카피·관계 도식·시각 토큰 절만.
- fullops-motion 스킬(toolchain.md, delivery.md)과 설치한 HyperFrames 공식 스킬의 `hyperframes` → `motion-graphics`, `hyperframes-core`, `hyperframes-animation`, `hyperframes-cli`.
- `.fullops-squad/project.md`, `AGENTS.md`의 Windows 파일 작업 안전.

## 해야 할 일과 파일 소유권

- [ ] `videos/service-intro/`에 package.json·lock과 `hyperframes@latest`를 로컬 설치한다(`DO_NOT_TRACK=1`, `HYPERFRAMES_SKIP_SKILLS=1`). 실제 버전·스킬 ref를 기록한다. `node_modules/`는 커밋하지 않는다(기존 .gitignore).
- [ ] 장면표를 먼저 쓴다(`videos/service-intro/storyboard.md`): 장면별 시작/끝 초·프레임, 시각 내용, 전환, 실제 문자열. 아래 방향을 따른다.
  - 0–3초: 워드마크 `KnowsLink` 등장.
  - 3–11초: 첫페이지 관계 도식 재현 — 사람 '나'와 '친구', 그 아래 '내 에이전트'·'친구의 에이전트'가 들여쓰기로 나타나고, '서로 수락한 친구'로 사람끼리 연결된 뒤 에이전트 사이로 메시지가 오가는 움직임. 문구 '사람끼리 친구, 에이전트끼리 대화'.
  - 11–15초: h1 '친구와 연결하고, 에이전트로 대화하세요'와 `link.knowslog.com`. 마지막 1초 이상 정지해 읽을 수 있게 한다.
  - 문구는 정본 문자열만 쓴다(필요하면 정본 안에서 고른다). 가짜 수치·회원 이름·아바타·온라인 표시·없는 기능은 넣지 않는다.
- [ ] 시각: 정본 여섯 토큰(canvas #FFFFFF, ink #182B31, muted #52646B, accent #005A62, selected #E8F3F3, edge #75868B)과 4px 모서리를 쓴다. 큰 그라디언트·네온·파티클·3D·글리치는 쓰지 않는다. 움직임은 절제된 이동·페이드·선 그리기 위주.
- [ ] 글꼴: 외부 요청 없이 로컬 `C:/Windows/Fonts/malgun.ttf`(Malgun Gothic, 정본 첫 글꼴)를 프로젝트 자산으로 참조해 로드한다. 외부 GSAP/font URL은 로컬 자산으로 바꾼다. 한글 glyph 로드 성공을 확인한다. 글꼴 파일은 커밋하지 않고 경로와 라이선스 확인 결과를 기록한다(배포 가능 여부 unknown이면 그렇게 적는다).
- [ ] preview로 장면을 확인하고 `hyperframes check --snapshots --json`을 실행한다. browserSkipped·빈 audit는 통과로 보지 않는다.
- [ ] `renders/v01/service-intro.mp4`로 렌더한다. 실패·수정 렌더는 새 버전 폴더(v02…)를 쓰고 이전 파일을 덮지 않는다.
- [ ] ffprobe로 길이 15.0초·1920×1080·30fps·450프레임·영상 스트림만 있음을 확인한다. 프레임 0, 각 장면 대표·전환 직전/직후, 449를 PNG로 추출해 직접 보고 문구 잘림·겹침·glyph 누락을 판정한다.
- [ ] 깨끗한 다른 폴더에서 소스·lock·명령으로 재렌더해 규격과 대표 프레임을 대조한다(바이트 동일성은 요구하지 않는다).

소유권: `videos/service-intro/`(소스·storyboard·package/lock·README에 재렌더 명령), `.fullops-squad/docs/evaluations/qa-reports/SAR-MOTION-INTRO-001-DEV-test/`(명령·exit·ffprobe JSON·프레임·판정). 제품 source·첫페이지·설계 문서는 수정하지 않는다. 최종 MP4는 10MB 이하면 커밋하고, 넘으면 커밋하지 않고 경로·sha256을 보고한다.

## 완료 기준과 검증

- 위 규격의 MP4가 있고 ffprobe·프레임 판정 근거가 delivery.md 검수 표 형식(passed|failed|not_run|unknown)으로 남는다. 음성 관련 항목은 무음이라 해당 없음으로 적는다.
- 재렌더 명령만으로 같은 규격이 재현된다.
- worker 완료와 제품 수락은 별개다. 시각·내용 최종 수락은 사용자(또는 이후 designer 시각 검수)가 한다.
- lint.py는 결과 HEAD에서 ERROR 0이어야 한다.

### UI 작업의 디자인 기준과 검증

- 디자인 정본: service-member-ui.md 공개 첫페이지 절의 토큰·글꼴·관계 도식.
- 재사용 컴포넌트: 없음(영상 전용 HTML). 의도적 예외: 정본은 '애니메이션 금지'를 공개 페이지에 둔 것이며 별도 영상에는 적용하지 않는다.
- 디자인 lint: 해당 없음. HyperFrames check로 대신한다.
- 테마 전환: 해당 없음(단일 밝은 테마).

## 갱신할 산출물

없음(route 결과). 영상은 홍보 자산이며 D01–D13 원천을 바꾸지 않는다.

## 기대 산출물

`videos/service-intro/` 소스·storyboard·README, `renders/v01/service-intro.mp4`(또는 경로·sha256), `qa-reports/SAR-MOTION-INTRO-001-DEV-test/` 검수 기록.

## 제약·협업·후속

하위 위임은 하지 않는다(단일 직렬 제작). 피할 스타일: 크림 배경, 알약 버튼, 모노스페이스 라벨, 큰 그라디언트, 네온·글로우, 장식 통계·배지, 이모지. 사용자 판단이 필요한 표현(문구 변경·톤)은 렌더 후 보고에 선택지로 적는다.

<!-- fullops-packet:start -->
### 탐색 근거와 읽을 구간

정본: `.fullops-squad\docs\evaluations\jev\SAR-MOTION-INTRO-001-DEV-packet.json` / SHA `c94b84303617564797f8bc21801dcd026811cccb` / partial=True
- `.fullops-squad/FULLOPS.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/handovers/to_dev.md` (document_read) · 줄 2, 2, 2, 6, 6, 7, 14, 14, 14, 24, 26, 49, 60, 60, 82, 82 · inferred · 필수
- `.fullops-squad/project.md` (document_read) · 줄 2, 10, 15 · inferred · 필수 · {'relevant': 0.17, 'evidence': 0.4, 'contradicts': 0.58, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `.fullops-squad/rules/common/README.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/coding-style.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/security.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/rules/common/testing.md` (document_read) · 줄 전체/미확인 · inferred · 필수
- `.fullops-squad/contexts/archive/dev-2026-10-10.md` (document_read) · 줄 75 · inferred
- `.fullops-squad/docs/deliverables/README.md` (document_read) · 줄 33 · inferred
- `.fullops-squad/docs/design-docs/architecture.md` (document_read) · 줄 12, 119, 197, 209 · inferred
- `.fullops-squad/docs/design-docs/crud-design.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/design-docs/data-model.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/design-docs/database-design.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/design-docs/interface-design.md` (document_read) · 줄 12, 156, 199, 298 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-MVP-001-UI.md` (document_read) · 줄 39 · inferred
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md` (document_read) · 줄 75 · inferred
- `.fullops-squad/docs/design-docs/mockups/service-member-ui.md` (document_read) · 줄 전체/미확인 · unknown
- `.fullops-squad/docs/design-docs/module-design.md` (document_read) · 줄 12, 28 · inferred
- `.fullops-squad/docs/design-docs/tech-stack.md` (document_read) · 줄 3, 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/COOR/dev-fix-2-final/coordinator-handover-supplement.md` (document_read) · 줄 13 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.2.0-review/report.md` (document_read) · 줄 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.3.0-FINAL-review/report.md` (document_read) · 줄 28 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.4.0-review/report.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.4.2-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/FULLOPS-UPDATE-1.4.4-FINAL-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-COPY-REVIEW-review/report.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-AUTO-RECEIVE-001-REVIEW-review/report.md` (document_read) · 줄 14, 46, 50 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-FINAL-review/report.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-N1-review/report.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-BETA-001-REVIEW-review/report.md` (document_read) · 줄 13 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-CODEX-HOOK-QUOTE-001-REVIEW-review/report.md` (document_read) · 줄 17, 67, 75 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-FRONTEND-POLISH-001-TESTER-review/report.md` (document_read) · 줄 13 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-DEV-review/report.md` (document_read) · 줄 17, 39 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-LIVE.md` (document_read) · 줄 18, 20, 28 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-GROK-WEBHOOK-001-TESTER-REVIEW-review/report.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-DEV-test/runtime/report.md` (document_read) · 줄 30 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-TESTER-review/report.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-001-TESTER-test/report.md` (document_read) · 줄 16, 45, 55 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-002-LIVE/report.md` (document_read) · 줄 19, 26 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CODE-REVIEW-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-CSS-REVIEW-review/report.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-DESIGN-REVIEW-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-FINAL-RECORDS-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-LANDING-GUIDE-001-VISUAL/css-followup/dev-final-handoff.md` (document_read) · 줄 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FINAL-review/report.md` (document_read) · 줄 26 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-REVIEW-FIX-review/report.md` (document_read) · 줄 26 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX.md` (document_read) · 줄 24 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-DEV-review/report.md` (document_read) · 줄 13 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-TESTER.md` (document_read) · 줄 137 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW-review/report.md` (document_read) · 줄 26, 28 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-REVIEW-review/report.md` (document_read) · 줄 23, 110 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-UI.md` (document_read) · 줄 31 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX-2.md` (document_read) · 줄 23 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-CODE-review/report.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-FINAL-review/report.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS-FINAL-review/report.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS-RESULT-review/report.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-001-OPS/report.md` (document_read) · 줄 60, 74 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-CLEANUP/report.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-CODE-REVIEW-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/adapter-report.md` (document_read) · 줄 19 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-DEV/audit.md` (document_read) · 줄 2 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-FINAL-MAIN-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-OPS-REVIEW-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-OPS/report.md` (document_read) · 줄 54 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REFACTOR-002-QA-REVIEW-review/report.md` (document_read) · 줄 16 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-REPO-DOCS-001-REVIEW-review/report.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-TESTER-review/report.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-review/report.md` (document_read) · 줄 57 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SERVICE-UI-001-TESTER-test/report.md` (document_read) · 줄 100 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-CODE-REVIEW-2-review/report.md` (document_read) · 줄 13 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-CODE-REVIEW-3-review/report.md` (document_read) · 줄 13, 26, 27, 75 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-CODE-REVIEW-review/report.md` (document_read) · 줄 13 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-DEV-test/report.md` (document_read) · 줄 32, 38, 92 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-RETEST-test/report.md` (document_read) · 줄 52, 58 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-001-TESTER-test/report.md` (document_read) · 줄 53 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SESSION-PEER-002-TESTER-test/report.md` (document_read) · 줄 86, 102, 128 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-test/report.md` (document_read) · 줄 17 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL/report.md` (document_read) · 줄 45 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-REVIEW-review/report.md` (document_read) · 줄 22 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-SYSTEM-DIAGRAM-001-REVIEW-review/report.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-DESIGN-REVIEW-review/report.md` (document_read) · 줄 20 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-UX-ISSUES-002-RECEIPT-CODE-REVIEW-review/report.md` (document_read) · 줄 13 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-WINDOWS-SHELL-GUARD-001-REVIEW-review/report.md` (document_read) · 줄 58 · inferred
- `.fullops-squad/docs/evaluations/qa-reports/SAR-WINDOWS-SHELL-GUARD-001-TESTER-test/report.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER-FINAL.md` (document_read) · 줄 13 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER-FIX.md` (document_read) · 줄 13 · inferred
- `.fullops-squad/docs/evaluations/scenarios/SAR-MVP-001-TESTER.md` (document_read) · 줄 13 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-0.9.12.md` (document_read) · 줄 42 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.2.0.md` (document_read) · 줄 59 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.0.md` (document_read) · 줄 36 · inferred
- `.fullops-squad/docs/exec-plans/phases/FULLOPS-UPDATE-1.4.4.md` (document_read) · 줄 18, 31 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-CODEX-HOOK-QUOTE-001-TESTER.md` (document_read) · 줄 18 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-002-DEV.md` (document_read) · 줄 27 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-GROK-WEBHOOK-001-TESTER.md` (document_read) · 줄 36 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV-FIX.md` (document_read) · 줄 80 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-BOT-CATALOG-DEV.md` (document_read) · 줄 56, 101 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md` (document_read) · 줄 34, 97 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-ACTUAL-TRIAL.md` (document_read) · 줄 16, 38, 46, 48, 52 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md` (document_read) · 줄 34, 80, 143, 189 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-BIDIRECTIONAL.md` (document_read) · 줄 96 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-TRIAL-CLEANUP.md` (document_read) · 줄 22, 42, 47, 70 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-RATE-REVIEW.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-TRIAL-REVIEW.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` (document_read) · 줄 24, 67 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SERVICE-UI-001-TESTER.md` (document_read) · 줄 40 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-INSTALL.md` (document_read) · 줄 57 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-DEV.md` (document_read) · 줄 36 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER.md` (document_read) · 줄 34 · inferred
- `.fullops-squad/docs/exec-plans/phases/SAR-WINDOWS-SHELL-GUARD-001-DEV.md` (document_read) · 줄 67 · inferred
- `.fullops-squad/docs/generated/db-schema.md` (document_read) · 줄 12 · inferred
- `.fullops-squad/docs/operations/transition.md` (document_read) · 줄 11, 47, 85, 154 · inferred
- `.fullops-squad/docs/operations/user-guide.md` (document_read) · 줄 11, 27, 29, 30, 32, 147, 148 · inferred
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (document_read) · 줄 2, 10 · inferred
- `.fullops-squad/docs/planning/SAR-PREP-002-request.md` (document_read) · 줄 2 · inferred
- `.fullops-squad/docs/planning/business-plan.md` (document_read) · 줄 3, 12 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-HUMAN-PEER.md` (document_read) · 줄 11 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (document_read) · 줄 3, 13 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (document_read) · 줄 3, 11, 42, 140, 149 · inferred
- `.fullops-squad/docs/planning/product-specs/SAR-SETUP-001.md` (document_read) · 줄 3, 34 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/README.md` (document_read) · 줄 1, 3 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/decisions.md` (document_read) · 줄 180 · inferred
- `.fullops-squad/docs/planning/sources/silent-agent-relay/product.md` (document_read) · 줄 2, 12, 14 · inferred
- `.fullops-squad/handovers/SAR-MVP-001-REVIEW.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_designer.md` (document_read) · 줄 27, 29, 89, 94, 97, 101, 102, 235, 274 · inferred
- `.fullops-squad/handovers/logs/2026-10-03_to_ops.md` (document_read) · 줄 17, 28, 60, 70, 115, 219, 240 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_dev.md` (document_read) · 줄 26, 28, 33, 89, 90, 110, 146, 197, 198, 204, 275 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_ops.md` (document_read) · 줄 23, 31, 70, 78, 118, 160, 161, 197, 229, 230, 270, 271, 295 · inferred
- `.fullops-squad/handovers/logs/2026-10-04_to_tester.md` (document_read) · 줄 23, 31, 73, 81, 122, 123, 166, 167, 207 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_designer.md` (document_read) · 줄 109 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_dev.md` (document_read) · 줄 115, 132 · inferred
- `.fullops-squad/handovers/logs/2026-10-05_to_ops.md` (document_read) · 줄 22, 64, 134 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_designer.md` (document_read) · 줄 26, 50, 94, 95, 183, 234, 236, 301, 381 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_ops.md` (document_read) · 줄 23, 95, 167, 191, 231, 233, 299, 352 · inferred
- `.fullops-squad/handovers/logs/2026-10-06_to_tester.md` (document_read) · 줄 23, 91, 115, 161, 163, 223 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_designer.md` (document_read) · 줄 30 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_dev.md` (document_read) · 줄 86, 87 · inferred
- `.fullops-squad/handovers/logs/2026-10-07_to_tester.md` (document_read) · 줄 47, 243, 256 · inferred
- `.fullops-squad/handovers/logs/2026-10-12_to_ops.md` (document_read) · 줄 30, 31 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_dev.md` (document_read) · 줄 32, 33, 293, 298 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_ops.md` (document_read) · 줄 34, 36 · inferred
- `.fullops-squad/handovers/logs/2026-10-13_to_tester.md` (document_read) · 줄 30, 85 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_designer.md` (document_read) · 줄 31, 32, 217, 218 · inferred
- `.fullops-squad/handovers/logs/2026-10-14_to_dev.md` (document_read) · 줄 29, 102, 103, 378, 379 · inferred
- `.fullops-squad/handovers/logs/2026-10-15_to_dev.md` (document_read) · 줄 30, 31, 46, 47 · inferred
- `.fullops-squad/handovers/logs/2026-10-17_to_ops.md` (document_read) · 줄 32, 33, 40, 114, 116, 118 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-DEV-REVIEW.md` (document_read) · 줄 17 · inferred
- `.fullops-squad/handovers/logs/SAR-SETUP-001-INTEGRATION-REVIEW.md` (document_read) · 줄 14 · inferred
- `.fullops-squad/lint/README.md` (document_read) · 줄 2 · inferred
- `.fullops-squad/orca-agents.md` (document_read) · 줄 36, 40, 41, 42, 43, 44, 50 · inferred
- `.gitignore` (impact_check) · 줄 전체/미확인 · inferred
- `AGENTS.md` (document_read) · 줄 전체/미확인 · inferred · {'relevant': 0.06, 'evidence': 0.2, 'contradicts': 0.19, 'injection': 0.04, 'decision': 'keep', 'reason': None}
- `Makefile` (impact_check) · 줄 전체/미확인 · inferred
- `README.md` (document_read) · 줄 2, 10, 84, 171, 190 · inferred
- `adapters/.claude-plugin/plugin.json` (impact_check) · 줄 6 · inferred
- `adapters/.codex-plugin/plugin.json` (impact_check) · 줄 6, 11, 14 · inferred
- `adapters/.cursor-plugin/plugin.json` (impact_check) · 줄 6 · inferred
- `adapters/.grok-plugin/plugin.json` (impact_check) · 줄 6 · inferred
- `adapters/README.md` (document_read) · 줄 2, 10, 55, 57, 58, 60, 249, 338 · inferred
- `adapters/hooks/register.test.ts` (impact_check) · 줄 210, 211, 251 · inferred
- `adapters/hooks/register.ts` (impact_check) · 줄 1, 35, 37, 39, 40, 114, 119, 131, 165, 170, 177, 190, 192, 195, 207, 222, 225, 232, 237, 245, 257, 258, 268, 272, 278, 289, 291, 295, 304, 308, 320 · inferred
- `adapters/manifest.json` (impact_check) · 줄 6, 30 · inferred
- `adapters/package.json` (impact_check) · 줄 전체/미확인 · inferred
- `adapters/skills/knowslink/SKILL.md` (document_read) · 줄 3, 6, 22, 36, 50 · inferred
- `adapters/src/mcp.ts` (impact_check) · 줄 156, 187, 427, 432, 440 · inferred
- `adapters/src/trial-cli.ts` (impact_check) · 줄 35 · inferred
- `adapters/src/trial-setup.ts` (impact_check) · 줄 전체/미확인 · inferred
- `cmd/migrate/main.go` (impact_check) · 줄 15 · inferred
- `cmd/relay/main.go` (impact_check) · 줄 15, 16 · inferred
- `deploy/knowslink/access_trial_plan.py` (impact_check) · 줄 25 · inferred
- `deploy/knowslink/tunnel/config.yml.tmpl` (impact_check) · 줄 전체/미확인 · inferred
- `deploy/knowslink/tunnel/public-ingress.yml` (impact_check) · 줄 전체/미확인 · inferred
- `deploy/knowslink/verify.py` (impact_check) · 줄 전체/미확인 · inferred
- `docs/README.md` (document_read) · 줄 2, 10 · inferred
- `docs/architecture.md` (document_read) · 줄 전체/미확인 · inferred
- `docs/development.md` (document_read) · 줄 30 · inferred
- `docs/operations.md` (document_read) · 줄 전체/미확인 · inferred
- `docs/system-design.md` (document_read) · 줄 전체/미확인 · inferred
- `go.mod` (impact_check) · 줄 1 · inferred
- `internal/relay/human_http.go` (impact_check) · 줄 전체/미확인 · inferred
- `internal/relay/human_pages.go` (impact_check) · 줄 7, 8 · inferred
- `internal/relay/human_pages_test.go` (impact_check) · 줄 52 · inferred
- `internal/relay/landing_integration_test.go` (impact_check) · 줄 전체/미확인 · inferred
- `internal/relay/landing_pages.go` (impact_check) · 줄 5, 12, 43, 46, 52, 54, 75 · inferred
- `internal/relay/member_pages.go` (impact_check) · 줄 48, 50 · inferred
- `internal/relay/member_ui_test.go` (impact_check) · 줄 전체/미확인 · inferred
- `scripts/check_public_ingress.py` (impact_check) · 줄 전체/미확인 · inferred
- `scripts/install_bot_mcp.sh` (impact_check) · 줄 2 · inferred
- `scripts/package_plugin.py` (impact_check) · 줄 41, 137, 142, 145, 155 · inferred
- `scripts/schema.py` (impact_check) · 줄 20 · inferred
- `scripts/shell_guard.py` (impact_check) · 줄 전체/미확인 · inferred
- `scripts/test_shell_guard.py` (impact_check) · 줄 전체/미확인 · inferred
- `scripts/verify_grok_plugin.py` (impact_check) · 줄 전체/미확인 · inferred
- `scripts/verify_setup.py` (impact_check) · 줄 전체/미확인 · inferred
미확인 22건: 정본의 unknown/producer_status/remaining_context_paths/optional_context_paths 확인. bounded string/definition search; dynamic references and language server semantics unverified
<!-- fullops-packet:end -->

## 완료 보고

미착수.
