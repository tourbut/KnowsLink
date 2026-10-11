---
title: SAR-WINDOWS-SHELL-GUARD-001-DEV — Windows 삭제 명령 사고의 실행 전 차단 및 회귀 검증
status: draft
updated: 2026-10-11
owner: dev
tasks: [SAR-WINDOWS-SHELL-GUARD-001-DEV]
summary: Windows 삭제 명령 사고의 실행 전 차단 및 회귀 검증
attempt: 2844ae38d8d64434889e5d5d1f385da2
base: a90a9e60ed1f009627aa39031fbf3c228d0d7599
subagent_level: standard
test_level: lite
---

# SAR-WINDOWS-SHELL-GUARD-001-DEV — Windows 삭제 실행 전 차단

- Task key: SAR-WINDOWS-SHELL-GUARD-001-DEV; Purpose: implementation.
- From/To: coor/dev. fresh Claude Opus5.5 high. 테스트 lite, 하위 위임 standard.
- 역할 체크아웃: C:/Users/shin/orca/workspaces/KnowsLink/dev, fullops/dev.
- 병합 책임자 coor, 기본 main/origin. 현재 수락된 SSE 제품 febd1282는 coor a90a9e6에 보존됐다.
- 복귀: run_dc9b4e6dd18a, coor terminal term_e61d3e14-29e9-4954-943a-4a75707c82de. 실제 task/dispatch는 worker-start preamble을 따른다.

## 확정 요구와 사고 근거

사용자는 파일복구를 중단하고 재발 방지만 한 뒤 개발을 재개하라고 지시했다. PowerShell의 큰따옴표 문자열로 Bash 명령을 전달하면서 Bash 변수에 쓴 backslash escape가 PowerShell 확장을 막지 못했다. 임시 디렉터리 변수가 사라져 약64초간 Windows C: 루트에 재귀 삭제가 실행됐다. 실제 앱 자료 삭제가 확인됐다. 원문·USN/VSS 자료는 F: 비공개 KnowsLink-Incident 폴더에 보존돼 있으며 제품 Git에 넣지 않는다. 사고 명령을 실제 실행하여 재현하지 않는다.

## 완료 조건과 파일 소유권

1. 실제 installed CLI와 repo hook 범위를 확인하고 Windows의 교차 셸 삭제·검증되지 않은 재귀 삭제를 실행 전에 차단한다. 지침만으로 강제 차단 완료를 주장하지 않는다. 먼저 사고 당사자인 Grok worker와 다음 Claude worker에 적용 가능한 최소 경로를 검증하고, Codex 지원 범위도 실제 근거로 구분한다.
2. 임시 삭제가 꼭 필요한 정상 경로는 단일 PowerShell의 LiteralPath, workspace 안의 검증된 실제 절대경로, reparse/드라이브 루트/빈 변수 방어를 유지한다. 검수 전 모든 삭제·이동 정리를 금지한 AGENTS 규칙을 지킨다. 정상 읽기·빌드·테스트·허가된 git 명령은 보존한다. 우회 가능한 일반 regex를 완전한 보안 경계로 주장하지 않는다.
3. provider별 실제 hook payload의 합성 deny/allow 회귀를 작성한다. 사고 인용 입력·빈/루트/부모탈출/교차 셸·reparse 및 정상 비파괴 작업을 검사한다. 실제 차단 hook 설치를 직접 증명하되 실제 삭제는 하지 않는다. hook 미지원/재시작 필요는 정확히 보고하고 안전한 실행 경로를 제시한다.
4. 새 비용/외부 서비스/권한 완화 없이 저장소 로컬 코드·설정·에이전트 안내와 실행 계획만 수정한다. 플러그인 캐시·vendor core·사용자 전역 설정은 수정하지 않는다. 실제 계정·서버·Google·DB·Tunnel·복원·임의 파일 삭제도 범위 밖이다. 기존 설정이 있으면 병합하며 다른 hook을 제거하지 않는다.
5. scripts/의 최소 검사기·작은 테스트·지원되는 repo별 CLI 설정·AGENTS/CLAUDE/GROK 안내와 .fullops-squad/docs/exec-plans/phases/SAR-WINDOWS-SHELL-GUARD-001-DEV.md가 소유 범위다. 제품 relay/adapter 코드는 바꾸지 않는다. D10 추천은 제품 모듈 변경이 없어 불필요하면 사유를 기록한다.
6. 고정 최종 SHA에서 등록 lint/test 및 직접 guard 회귀를 실행하고 packet outcomes/work finish/archive/commit을 완료한다. Windows make 대신 격리 Linux clone이 필요하면 기존 helper를 참고하되 DEV 이전 final 고정 shim으로 다른 HEAD를 검사했다고 하지 않는다. 고유 임시 clone과 실제 HEAD/clean 확인, 종료코드 전달을 유지한다. 정리 삭제는 하지 않는다.

## 먼저 읽을 문서 / 적용 기준

AGENTS.md, .fullops-squad/FULLOPS.md, .fullops-squad/project.md, .fullops-squad/rules/common/README.md와 coding-style/testing/security, .fullops-squad/rules/delegation.md, .fullops-squad/docs/agents/document-writing.md, .fullops-squad/orca-agents.md, .fullops-squad/fullops.json, .fullops-squad/docs/evaluations/jev/SAR-WINDOWS-SHELL-GUARD-001-DEV-route.json.
기준 fullops-common-0.3.3. writing-for-agents를 적용한다. CLI/API 근거는 installed help/code와 Context7을 사용한다. Context7 quota exceeded면 공식 문서로 대체하고 출처를 남긴다. Grok/Claude/Codex 실제 지원을 추정하지 않는다. 거대한 과거 제품 로그를 전부 읽지 않는다.
기존 안전 규칙이 있었는데 실행이 차단되지 않은 사고이므로 Markdown 추가만으로 끝내지 않는다. 최소 검증 가능한 강제 방어를 목표로 기술 계획·구현·검사를 같은 과제에서 완료한다.

## 독립 검수와 복귀

### 지시 전제와 충돌 — 먼저 확인

context가 기존 SSE 리뷰·project·orca-agents의 충돌을 추천했다. 기존 SSE 기능 수락과 Windows 삭제 사고 예방은 다른 검증이다. SSE 수락을 사고 방지 완료로 사용하지 않는다. coor 역할과 독립 리뷰 규칙을 유지하며 이번 기술 판단은 DEV가 한다. code find의 존재 unclear는 현재 저장소에 호스트 실행 차단 모듈이 확인되지 않았다는 뜻이다. adapter 플러그인 manifest는 실제 호스트의 shell hook 설정이 아니므로 혼동하지 않는다. 제시된 역사 archive는 optional 참조이며 이번 실행 차단에 영향 없는 사유를 outcomes에 남긴다.

구현자와 다른 실제 세션의 clean detached fixed-SHA 독립 리뷰가 후속이다. 미해결 critical/high는 수락하지 않는다. 호스트 hook 활성화와 단위 parser PASS를 구분한다. 현재 세션/실제 명령 지원·적용된 role 경로/재시작 필요/알려진 우회 범위/실행하지 않은 검사를 완료 보고에 적는다. 완료 메시지는 과제 키·최종 40자리 SHA 하나를 명확히 포함해 worker_done으로 보낸다.
