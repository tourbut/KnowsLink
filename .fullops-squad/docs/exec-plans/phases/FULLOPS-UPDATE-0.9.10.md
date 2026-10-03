---
title: FullOps 0.9.10 업데이트 적용 기록
status: draft
updated: 2026-10-03
owner: coor
tasks: [FULLOPS-UPDATE-0.9.10]
summary: 설치와 레포 적용 버전을 구분하고 검증 정책 0.3.2의 적용 및 동기화 보류를 기록한다
---

# FULLOPS-UPDATE-0.9.10 — 레포 적용

## 기준과 설치

- 기준 ref: `3d91754313028cb5350b2c0f0e2505d2828088ae`. 시작 시 main 작업 트리는 깨끗했다.
- 현재 CLI는 Codex다. 실제 설치 버전은 갱신 전 0.9.10, 갱신 후 0.9.10이다. 레포 적용 버전은 0.9.9다.
- `codex plugin list --json`과 설치 경로로 실제 버전을 확인했다.
- `codex plugin marketplace upgrade fullops-squad`는 이미 최신임을 확인했다. `codex plugin add fullops-squad@fullops-squad`는 성공했다.
- 새 설치의 update-fullops를 다시 읽었다. `deps.py --check`는 종료코드 0이며 필수 CLI가 모두 있다.
- `update.py --repo . --from 0.9.10`으로 레포 적용 기준 0.9.9 이후 릴리스 0.9.10을 확인했다.
- 최초 설치 목록 후보 `~/.codex/plugins/installed_plugins.json`은 없었다. 실제 설치 확인은 Codex CLI로 완료했다.

## 릴리스 0.9.10 판정

1. 적용: 공통 README와 testing을 비교하고 0.3.2 변경을 통합했다. project.md의 적용 기준을 갱신했다. 기존 추가 기준은 보존했다.
2. 적용: 핸드오버 템플릿과 시나리오 안내의 변경 절을 통합했다. 검사 담당·대상·시점·조건과 DEV 완료·최종 수락을 구분한다.
3. 적용: FULLOPS.md에 update-fullops를 연결했다. project.md·orca-agents.md·PLANS.md에 관련 회귀·독립 전체 QA·직접 시각 검수와 증거 재사용·보류 책임을 반영했다.
4. 적용: 캡처·영상과 반복 범위 정책을 반영했다. 해당 없음: 이번 운영 변경에는 제품 실행·시각 자료를 새로 만들 필요가 없다.
5. 적용: 기존 실패·held·미해결 critical/high와 제품 정지·최종 수락 기준을 유지한다. 진행 중 지시서와 완료 보고는 수정하지 않았다.
6. 적용: 설정·문서 링크·변경 범위를 검증한다. 필수 검증 후 레포 적용 버전을 0.9.10으로 갱신한다.

동일한 coor·designer·dev·ops·tester 역할로 setup dry-run과 실행을 수행했다. 종료코드 0이며 신규 파일과 브랜치는 0개다.
변경 파일은 위 운영 문서, fullops.json, board/board.json 및 이 기록이다. 문서 front matter는 deliverables.py로 등록한다.

## 보존과 후속 담당

역할·모델·원격·기존 제품 자료·인박스·컨텍스트·QA 기록을 유지한다. 제품 worker를 새로 시작하지 않았다.
모든 역할 동기화와 push는 이번 작업에서 수행하지 않는다. coor가 실제 활성 상태와 열린 인계를 확인한 뒤 안전한 체크아웃에 준비 커밋을 전달한다. 이는 레포 필수 적용의 보류가 아닌 전달 보류다.
coor·designer·dev는 기존 SAR 인계가 있다. tester는 이전 업데이트 이후 별도 HEAD를 갖는다. 깨끗한 작업 트리만으로 유휴 상태를 판단하지 않는다.
Claude Code·grok의 설치 갱신은 해당 호스트의 다음 세션 담당자가 수행한다. 이번 요청의 현재 CLI는 Codex다.
후속 coordinator는 새 세션에서 기존 Run `run_8ca8bc058ab7`을 연결하고 진행 상태를 확인한다. 기존 제품 과제의 기준 ref와 원천 SHA는 변경하지 않는다.

## 검증

제품 코드 변경은 없다. main에는 제품 테스트 실행기가 없으므로 제품 테스트·빌드는 해당 없음이다.
Git 공백 검사, JSON·역할 보존, 변경 Markdown 상대 링크, 산출물 strict, 현황판 생성과 기준 ref 이후 FullOps lint를 수행한다. 실제 결과는 아래에 기록한다.

- 준비 정책 커밋: `36b164d49693`. 기준 ref 이후 lint는 종료코드 0, ERROR 0, WARNING 1, 실행 불가 0이다.
- WARNING LINT-000은 main에 프로젝트 lint 실행기가 없는 기존 상태다. 제품 검사 통과로 취급하지 않는다.
- 산출물 strict는 검사 13, 미작성 13, 문제 0, 경고 0이다. Git 공백·상대 링크·역할·원격 보존 검사와 현황판 생성은 통과했다.
- 필수 적용과 검증을 완료해 fullops.json의 plugin_version을 0.9.10으로 갱신했다. 최종 커밋에도 같은 기준 ref의 lint를 실행한다.
- 준비 전달 담당은 coor다. 재개 조건은 역할의 실제 활성 상태·열린 인계·작업 트리 확인과 안전한 동기화 범위 확정이다. 원격 push와 역할 체크아웃 변경은 하지 않았다.

## 사용자 요청에 따른 하위 워크트리 적용

2026-10-03 사용자가 하위 워크트리 적용을 요청했다. 기존 보류를 해제하고 모든 역할에 준비 커밋 e825dc7을 병합했다.

- coor: `f5efbc5`. 미추적 SAR-SETUP-001-DEV-099-review 폴더는 그대로 유지했다.
- designer: `78ac098`. 기존 제품 기획 자료와 완료 기록을 유지했다.
- dev: `830131a`. 제품 코드와 프로젝트 검증 명령을 유지했다.
- ops: `45acdec`.
- tester: `a3d654b`. 기존 구현 결과와 독립 QA 인계를 유지했다.

충돌은 기존 역할별 기술 기준·SAR 계획·제품 단계와 새 운영 규약을 함께 보존해 해결했다. 모든 역할의 기존 HEAD와 준비 커밋은 새 HEAD의 조상이다. 모든 역할의 plugin_version은 0.9.10이다.
역할별 lint는 변경 전 HEAD를 기준으로 실행했다. designer·ops는 ERROR 0, WARNING 1이며 기존 LINT-000이다. dev·tester는 product-lint를 포함해 ERROR 0, WARNING 0이다.
coor는 미추적 사용자 자료 때문에 원본 작업 트리 lint가 거부됐다. 동일 고정 SHA의 깨끗한 detached snapshot에서 검사했다.
tester와 coor snapshot은 prettier 의존성 누락으로 첫 product-lint가 실패했다. lockfile 기준 npm ci 후 재검증했다. 실패는 제품 결함이나 QA 통과로 바꾸어 기록하지 않는다.
원격 push는 수행하지 않았다. 이미 열린 에이전트 세션은 이전 hook을 유지할 수 있으므로 다음 세션에서 갱신 규약을 읽는다.

coor 고정 SHA snapshot의 최종 product-lint와 FullOps lint는 종료코드 0, ERROR 0, WARNING 0이다. 원본 미추적 자료는 이동·커밋·삭제하지 않았다.
