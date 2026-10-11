---
title: SAR-WINDOWS-SHELL-GUARD-001-TESTER — Windows 삭제 실행 차단의 독립 fixed SHA 검수
status: draft
updated: 2026-10-11
owner: tester
tasks: [SAR-WINDOWS-SHELL-GUARD-001-TESTER]
summary: Windows 삭제 실행 차단의 독립 fixed SHA 검수
attempt: 4c29bb1b670045eba09ab8c9c9a49fb4
base: 0e4f638d2f99622a320cb38473b0768604af2bff
subagent_level: standard
test_level: lite
---

# SAR-WINDOWS-SHELL-GUARD-001-TESTER — Windows 실행 차단 독립 검수

- Task key: SAR-WINDOWS-SHELL-GUARD-001-TESTER; Purpose: review.
- 상태: ready. 후보 4e82b4f4cd1bc31a59d0d230349fd34d95c7f0c7.
- Review SHA: 4e82b4f4cd1bc31a59d0d230349fd34d95c7f0c7; Implementer session: bacd5d1f-de07-43b3-9efa-0df5636729c7. 실제 reviewer session은 착수 후 기록한다.
- 역할: tester, fresh Claude Sonnet5.5 high. lite/standard.
- 경로: C:/Users/shin/orca/workspaces/KnowsLink/tester, fullops/tester.
- 복귀 run_dc9b4e6dd18a/coor term_e61d3e14-29e9-4954-943a-4a75707c82de. 실제 task/dispatch는 preamble을 따른다.
- review 키 SAR-WINDOWS-SHELL-GUARD-001-REVIEW, 원 구현 task-key SAR-WINDOWS-SHELL-GUARD-001-DEV.

## 요구와 완료 조건

사용자는 Windows 드라이브 삭제 사고의 복구를 중단하고 재발 방지 후 개발 재개를 지시했다. 제품 SSE는 기존 febd1282 독립 검수 자료를 재사용하며 새 실행으로 바꾸지 않는다. 이번 검수는 새 실행 차단과 관련 agent 안내·설정·테스트만 소유한다.

DEV actual session과 다른 실제 session에서 clean detached fixed-SHA snapshot을 생성한다. snapshot은 읽기 전용이다. 리뷰 결과·증거·필요한 작은 테스트만 tester 체크아웃에 작성한다. 제품 코드·hook 정책·전역 설정·원본 사고 자료는 수정하지 않는다.

실제 지원 host hook 종류/설정 위치/적용 role/재시작 필요를 installed 근거로 확인한다. parser 단위 PASS와 실제 CLI 활성화를 구분한다. 안전 검사 자체가 우회되는 한계가 있으면 정확히 보고한다. Windows 교차 셸·인용·빈값·드라이브 루트·부모 탈출·reparse 경계와 malformed 입력, 정상 비파괴 작업을 검증한다. 실제 삭제/이동/복원/cleanup 명령은 어떤 테스트에서도 실행하지 않는다. 실제 hook 결선을 시험할 때도 미작동 시 무해한 명령만 사용한다.

최종 exact SHA의 등록 lint/test는 DEV 실제 근거의 HEAD/clean/종료코드를 확인하여 재사용한다. 변경·근거 결함이 없으면 전체 제품 테스트를 반복하지 않는다. 독립 guard 경계 QA는 직접 실행한다. 미해결 critical/high·필수 검사 실패·실제 결선 미확인을 PASS로 바꾸지 않는다. source 지원 범위 밖 클라이언트를 적용 완료로 표시하지 않는다.

review.py prepare/check, 원 task-key 연결, 모든 diff 대상 reviewed/skipped 사유, packet outcomes, fullops-common-0.3.3 및 프로젝트 정본 일치, session 독립성을 확인한다. finish/archive/빈 인박스와 실제 최종 40자리 SHA를 포함해 worker_done으로 보고한다. 임시 snapshot/자료는 정리하지 않는다.

## 먼저 읽을 정본

AGENTS.md, .fullops-squad/FULLOPS.md, .fullops-squad/project.md, .fullops-squad/rules/common/README.md와 coding-style/testing/security, .fullops-squad/rules/delegation.md, .fullops-squad/docs/agents/document-writing.md, .fullops-squad/orca-agents.md, 해당 DEV 실행 계획과 완료 보고/guard 테스트/host hook 설정. 구현 최종 경로·SHA·session은 배정 전에 확정한다.

구현 정본은 scripts/shell_guard.py, scripts/test_shell_guard.py, .claude/settings.json, .codex/hooks.json, AGENTS.md, CLAUDE.md, Makefile, .fullops-squad/docs/exec-plans/phases/SAR-WINDOWS-SHELL-GUARD-001-DEV.md다. DEV 최종 exact SHA pass.json(DEV git-dir fullops-gate)은 coor가 확인했다. product-lint/test0, ERROR0/WARNING5/unavailable0. provider 4종 합성 deny43/allow17이며 실제 삭제는 없었다. 실제 Claude 두 도구 차단과 Grok inspect enabled/trusted, Grok/Codex runner 재현은 서로 다른 증거다. 실제 Grok/Codex 모델 tool call과 Codex 사용자 신뢰는 미실행이다.

### 알려진 지원 한계 — 반드시 판정

DEV 계획의 fail-open·script/alias 우회·Codex session cwd 상대경로 실패·TOCTOU·오탐을 직접 평가한다. 플랫폼 한계와 구현 결함을 구분한다. 현재 위험 재발 방지 목적을 충족하지 못하는 경계는 finding으로 올리며 설명만으로 통과시키지 않는다. 기존 세션 재시작·checkout별 사용자 신뢰는 수락 후 사용자 단계로 남긴다. 원 고정 snapshot에서만 제품을 읽고 reviewer 기록 체크아웃을 수정한다.

## 산출물과 후속

qa-reports/SAR-WINDOWS-SHELL-GUARD-001-TESTER-test와 SAR-WINDOWS-SHELL-GUARD-001-REVIEW-review에 실제 증거를 보존한다. 대화 원문·사용자 private 파일·사고 원시 USN/VSS는 Git에 넣지 않는다. 수락 후 실제 main 통합과 SSE 운영/계정 테스트 재개는 coor/ops 담당이다.

### 지시 전제와 충돌 — 먼저 확인

Context가 project.md를 충돌 후보로 표시했다. coor 원문 확인 결과 과거 제품 후보 상태와 Docker 검증 명령이 남아 있다. 현재 과제는 제품 변경 없이 Windows hook 경계를 검수한다. 사용자 지시인 로컬 Docker 금지·복구 중단·실제 삭제 없는 검증이 우선이다. 등록 lint/test는 DEV exact SHA의 실제 Linux 증거를 확인해 재사용한다. 공통 보안·독립성·critical/high 차단은 유지한다. Hook JSON 두 파일은 탐색기의 민감 경로 제외로 API에 보내지 않았으며 reviewer가 snapshot에서 직접 읽는다. partial 후보는 일반 검색과 outcomes로 보완한다.
