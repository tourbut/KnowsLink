---
title: FULLOPS-UPDATE-1.9.0
status: draft
updated: 2026-10-25
owner: coor
tasks: [FULLOPS-UPDATE-1.9.0]
summary: 1.8.7–1.9.0을 확인하고 지시서 worker_done 형식과 lint 안내(paths·seconds·로그·--level·LINT-002)를 반영한다
---

# FULLOPS-UPDATE-1.9.0

## 버전

- 업데이트 전 설치: 1.9.0(호스트 자동 업데이트). `claude plugin update`가 "already at the latest version (1.9.0)"을 확인했다. fullops-meter 0.10.1 최신.
- 레포 적용 버전: 1.8.6 → 1.9.0.
- 대상 버전: 1.9.0.

## 설치

`deps.py --host claude-code` exit 0, `deps.py --check --host claude-code`가 `호스트 의존성 검증 완료`를 출력했다. 로그는 기본 체크아웃 `.devtools/deps-1.9.0.log`. `frontend.json`이 없어 선택형 프런트 도구는 해당 없음.

## 릴리스 판정

| 릴리스 | 판정·근거 |
|---|---|
|1.8.7 설치기 반복·스킬 백업·grok|해당 없음. 레포 파일 변경 없음. 쌓인 스킬 백업 정리는 선택이라 하지 않았다.|
|1.8.8 설치기 GitHub 실행|해당 없음. 이후 통합 설치기는 `npx --yes github:tourbut/Fullops-Squad-plugin update`를 쓴다.|
|1.8.9 지시서 worker_done 형식|적용. `handovers/_TEMPLATE.md` `완료 보고`에 한 줄 형식 문장을 넣었다. 레포 고유 절은 보존했다. 나머지 hook은 새 세션에서 적용된다.|
|1.9.0 lint 안내|적용. `lint/README.md`에 LINT-002 행과 non-required 시험 skip·`paths`·`seconds`·전체 출력 로그·`--level` 설명을 통합했다. 레포 고유 절(SAR-SETUP-001-DEV·0.9.14 기록, hook 확인 절)은 보존했다.|
|1.9.0 LINT-002 범위 조정|해당 없음. `lint.json`에 `required: true` 시험 명령이 없다.|
|setup dry-run|생성 예정 0개.|

## 검증

- `update.py --repo . --from 1.9.0 --verify`: 필수 적용 0/0 충족, exit 0.
- 제품 코드 변경 없음.

## 보류와 동기화

- 보류 없음(원본 `KnowsLinkDevTools` 삭제는 1.8.6 기록의 사용자 보류를 그대로 둔다).
- 진행 중 dev·tester는 동기화를 예약한다. 유휴 designer·ops는 최신 main으로 동기화한다.
- Claude Code는 1.8.9의 새 hook 이벤트(PostToolUseFailure)를 읽으려면 세션을 다시 시작해야 한다.
