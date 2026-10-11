---
title: Windows 실행 차단 F1~F4 독립 재검수
status: review
updated: 2026-10-11
owner: coor
tasks: [SAR-WINDOWS-SHELL-GUARD-001-FIX-DEV]
summary: 고정 c2a0ce4의 최소 보완을 다른 세션에서 검수하고 Windows 회귀와 실행기 경로를 확인한다.
---

# Windows 실행 차단 F1~F4 독립 재검수

수락한다. 구현자 d370df55-6e39-4ca7-9b30-c7c8de2ea612과 다른 coor 세션 01a1203d-35f3-7ac2-9d7d-a2d8bd9823ed에서 검수했다. base 972ccd6f30025260bcfaeb2dec91c1e68f8a5f1c, head c2a0ce4ab673116d0d18be368032ac26e1327a9c다. snapshot은 result.json에 기록된 clean detached 경로이며 읽기 전용으로 유지했다.

OCR 규칙 그룹과 15개 diff 항목을 확인했다. 제품·지시·완료 기록 8개는 reviewed, 생성 탐색 JSON 7개는 skipped 사유와 packet 검증으로 보존했다. fullops-common-0.3.3과 project.md의 코딩·보안·테스트 기준을 적용했다. 사용자 로컬 Docker 금지와 복구 중단이 과거 일반 실행 명령보다 우선한다. UI·의존성·D01~D13 변경은 없다. 제품 relay/adapter/DB/cmd는 원 수락 후보 이후 diff 0이다.

F1의 콜론 옵션, F2의 interpreter 삭제, F3의 Git global 옵션 감지와 정상 allow를 코드 및 회귀로 확인했다. F4는 Windows cmd용 hook이 git top-level을 찾도록 변경됐다. 고정 snapshot에서 Windows 회귀 4 tests OK를 직접 확인했다. coor-codex-runner.json은 root/scripts/adapters 각각 정상 payload exit0와 거부 payload exit2 총 6건의 독립 실행 근거다. 삭제 문자열은 guard stdin에만 전달했고 실제 삭제·이동·복원·cleanup은 실행하지 않았다. temp는 보존했다.

lint.json은 DEV exact SHA의 실제 실행 근거를 복사해 재사용했다. pass.json의 head/base/attempt와 일치하며 ERROR0/WARNING3/unavailable0, product-lint/test exit0다. Linux task-local 로그의 같은 SHA·clean 실행 근거를 확인했다. tester/coor가 제품 전체 검사를 다시 실행했다고 표시하지 않는다. SIZE-001은 기존 PLANS 누적 크기, SLOP-004 두 건은 합성 테스트 문자열이다. 제품 추가 61줄이며 SIZE-002·DEP-001은 없다.

새 critical/high는 없다. 원 tester의 F1~F4는 이번 변경으로 처리됐다. 원 리뷰의 F5~F9(다른 구문·별칭·스크립트·이동·대형 입력 timeout·reparse/TOCTOU 등 문자열 검사의 한계)는 남는다. 실제 Grok/Codex 모델 tool call·Codex 사용자 신뢰·전역 강제는 미확인이다. Codex 명령은 Windows cmd 전용이며 Linux 보호로 표시하지 않는다. 사용자 hook 신뢰를 자동 승인하지 않았다.

대화 미참조 인계는 archived to_dev의 요구·범위, DEV exec plan의 기술 계획/지원 근거/한계/활성화, source/tests/hook 설정, 원 tester report, PLANS에서 확인했다. 수정 뒤 문서와 코드가 일치하며 CLAUDE.md→AGENTS.md 링크는 유효하다. 복구 원시 자료는 비공개라 읽지 않았다. 원 tester·DEV 실패 및 조건부 수락 기록은 수정하지 않는다. 신규 고정 SHA의 review.py check 후 main 병합/push와 역할 동기화는 coor가 담당한다. snapshot cleanup은 사용자 자료 보존 방침으로 수행하지 않는다.
