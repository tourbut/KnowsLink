---
title: FULLOPS-UPDATE-1.14.0
status: draft
updated: 2026-11-02
owner: coor
tasks: [FULLOPS-UPDATE-1.14.0]
summary: 1.11.0–1.14.0을 확인하고 lint 안내(LINT-002 범위·lint_tune)를 반영한다
---

# FULLOPS-UPDATE-1.14.0

## 버전

- 업데이트 전 설치: 1.14.0(호스트 자동 업데이트). `claude plugin update`가 "already at the latest version (1.14.0)"을 확인했다.
- 레포 적용 버전: 1.10.0 → 1.14.0.
- 대상 버전: 1.14.0.

## 설치

`deps.py --host claude-code`는 "의존성 설치 완료"로 끝났지만 첫 `--check`가 hook runtime 미갱신을 알렸다. `deps.py --runtime-only`로 `~/.fullops-squad/runtime/fullops-squad`를 갱신한 뒤 `--check --host claude-code`가 `호스트 의존성 검증 완료`를 출력했다. `frontend.json`이 없어 선택형 프런트 도구는 해당 없음.

## 릴리스 판정

| 릴리스 | 판정·근거 |
|---|---|
|1.11.0 이슈 자동 모드 lease·조회 오탐|보류(사용자). 자동 모드 상태는 enabled·paused다. 다시 쓰려면 새 coor 세션에서 `enable`하고 이전 running·claimed 작업은 `reconcile`로 인계한다.|
|1.11.1 lease 링크 제거|해당 없음. 기본 체크아웃과 다섯 워크트리에 lease 파일 링크가 없다.|
|1.12.0 슬래시 명령·사이드바|해당 없음. 다음 자동 모드 상태 변경 때 `status.json`이 갱신된다. 사이드바는 fullops-meter reload 뒤 보인다.|
|1.12.1 held 이슈 질문|해당 없음. 자동 질문은 없고 held 이슈는 coor가 필요할 때 `question`으로 묻는다.|
|1.13.0 done-gate·PR 후속·prune|해당 없음. 설정 변경 없음.|
|1.14.0 `lint/README.md` LINT-002 행·`rules` 설명|적용. LINT-002 행을 새 문구(범위 없는 명령 포함)로 바꾸고 `lint_tune.py` 안내 줄을 `rules` 설명 뒤에 넣었다. 레포 고유 내용은 보존했다.|
|1.14.0 `lint_tune.py` 실행|적용. 리뷰 98건·지적 184건. 사용자 `rules` 조정 후보 없음. 경고가 많은 내장 규칙은 SIZE-001 208건(지적 0)이다. 범위 없는 명령 product-lint·product-test는 Windows 실행 불가로 시간 0초라 120초 후보가 아니다. 설정은 바꾸지 않았다.|
|setup dry-run|생성 예정 0개.|

## 검증

- `update.py --repo . --from 1.14.0 --verify`: 필수 적용 0/0 충족, exit 0.
- 제품 코드 변경 없음.

## 보류와 동기화

- 이슈 자동 모드 재활성화는 사용자가 원할 때 새 coor 세션에서 한다.
- 이 세션의 hook은 `/reload-plugins` 뒤 1.14.0 경로로 바뀐다.
