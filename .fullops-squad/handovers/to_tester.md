---
title: SAR-GOOGLE-CONNECT-001-TESTER — 고정 인증 변경의 독립 코드 리뷰와 운영 서버의 좁은 Google 연결 QA
status: draft
updated: 2026-10-10
owner: tester
tasks: [SAR-GOOGLE-CONNECT-001-TESTER]
summary: 고정 인증 변경의 독립 코드 리뷰와 운영 서버의 좁은 Google 연결 QA
attempt: 8cc35b1d14d84a9ca2e9cf175b0cb6d5
base: f9f7675be6c9502bd6ab2f810bd42ff62a26a174
subagent_level: off
test_level: lite
---

# SAR-GOOGLE-CONNECT-001-TESTER — 고정 인증 변경의 독립 코드 리뷰와 운영 서버 lite QA

- From / To: coor / tester; 상태: ready.
- Task key: SAR-GOOGLE-CONNECT-001-TESTER; Purpose: review; Test level: lite; Subagent level: off.
- 검토 SHA: `152217f63cced85f620695961955d680ce27864e`; 기준 ref: `f9f7675be6c9502bd6ab2f810bd42ff62a26a174`.
- 구현자 실제 CODEX_THREAD_ID: `01a12123-fd98-7b80-880e-6c2802149a6d`. reviewer는 자기 실제 별도 host 세션 ID를 기록한다.
- 기록 체크아웃: C:/Users/shin/orca/workspaces/KnowsLink/tester. 복귀: Run run_86e0e674b5a0, term_e61d3e14-29e9-4954-943a-4a75707c82de. 실제 task/dispatch는 preamble 정본이다.

## 목표·범위·정본

Google 로그인·명시적 동의·클라이언트별 agent/키·로컬 자격 자동 저장의 인증 변경을 독립 검토한다. 제품 코드는 수정하지 않는다. 고정 SHA의 clean detached snapshot을 review.py snapshot으로 생성하고 읽기 전용으로 유지한다. 실제 리뷰 결과는 tester 기록 체크아웃에 쓴다. fullops-review·open-code-review-delegate·공통 fullops-common-0.3.3 세 규칙·project.md·document-writing.md를 직접 읽는다.

기존 제품 규칙은 고정 SHA의 제품 명세 SAR-PUBLIC-SERVICE와 DEV 완료 실행 기록이다. 먼저 고정 SHA의 `.fullops-squad/docs/exec-plans/phases/SAR-GOOGLE-CONNECT-001-DEV.md`, `docs/evaluations/qa-reports/SAR-GOOGLE-CONNECT-001-DEV/report.md`, D12 ops-guide의 Google 연결 적용·복귀 절을 읽는다. 코드는 merge-base diff의 변경 파일을 전부 검토하며 관련 호출자·검사만 좁게 확인한다. 새 의존성 없음·SQL schema 변경 없음·기존 동의/회원/키 분리·held 호환을 확인한다.

## 독립 리뷰와 검증

리뷰 key는 SAR-GOOGLE-CONNECT-001-TESTER다. review.py prepare는 기준/검토 SHA로 실행한다. result의 모든 path/status를 reviewed/skipped 및 근거로 채운다. report에는 고위험 인증·Origin/CSRF·만료/replay·회원/클라이언트·개인키 저장/권한·리다이렉트 경계·사용자 화면/동의·운영 ingress의 실제 판정을 적는다. 미해결 critical/high는 수락하지 않는다. 수정은 coor를 통해 DEV로 돌린다.

최종 SHA의 실제 서버 clean-clone lint 증거는 `D:/workspace/KnowsLink/.git/worktrees/dev/fullops-gate/google-connect-final-lint.json`이다. ERROR0/WARNING8/실행 불가0, 등록 make lint/test exit0다. 일치하는 head/base/config와 경고 근거를 검토해 같은 리뷰 폴더 lint.json에 보존한다. 변경 없는 등록 전체 검사를 반복하지 않는다. 필요할 때만 서버 clean clone에서 다시 실행한다.

별도 운영 서버 임시 공간·격리 DB/loopback 포트에서 핵심 연결 성공·명시적 동의·만료/replay·다른 회원/클라이언트 분리와 필요한 기존 로그인 회귀를 짧게 검증한다. DEV self-check 재사용과 직접 독립 실행 결과를 구분한다. 실제 Google/외부 Bot 성공을 합성 검사로 주장하지 않는다. 로컬 Docker를 기동하지 않는다. 운영 서비스·DB·키·Tunnel·다른 서비스를 보존한다.

서버 접근은 자기 루트 `.env.server`의 main 정본 링크를 비공개로 읽는다. IP/USER/PW를 대화·Git·로그·셸 인자로 출력하지 않는다. 서버 소스·도구·재실행 스크립트 정본은 `D:/workspace/KnowsLink/.git/worktrees/dev/fullops-gate/google-connect-server-handoff.json`이다. 새 고유 임시 clone/DB만 사용하고 자기 임시 자원만 정리한다. DEV 접속 도구 Paramiko5는 Python313에 설치됐다. 관련 라이브러리 근거가 필요하면 Context7, 실패 시 공식 문서를 사용한다.

## 완료 조건

고정 SHA 독립 리뷰·좁은 QA의 명령/종료코드·수락/실패·생략 근거를 기록한다. fullops-review check를 같은 기준/검토 SHA와 task-key로 실행한다. 별도 결과·report·lint와 QA 보고서를 커밋한다. packet outcomes와 work.py finish로 지시서·결과 전문을 보존한다. authentic worker_done으로 결과 SHA·구현 리뷰 SHA·발견 사항·검사·snapshot 경로·정리 조건을 보고한다. 검토자의 임시 snapshot 정리는 부모가 실제 release 이후 처리한다.
