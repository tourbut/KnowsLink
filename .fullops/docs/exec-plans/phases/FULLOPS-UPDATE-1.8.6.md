---
title: FULLOPS-UPDATE-1.8.6
status: draft
updated: 2026-10-24
owner: coor
tasks: [FULLOPS-UPDATE-1.8.6]
summary: 1.8.6 필수 적용 4건(작업 파일 위치·local 제외·검증 하네스·지시서 위치)을 반영하고 .devtools를 예외로 기록한다
---

# FULLOPS-UPDATE-1.8.6

## 버전

- 업데이트 전 설치: 1.8.6(호스트 자동 업데이트). `claude plugin update`가 "already at the latest version (1.8.6)"을 확인했다. fullops-meter 0.10.0 최신.
- 레포 적용 버전: 1.8.5 → 1.8.6.
- 대상 버전: 1.8.6.

## 설치

`deps.py --host claude-code`는 exit 0이었다. 첫 `--check`는 "hook runtime missing or outdated"로 실패했다. `deps.py --runtime-only`로 hook 실행 사본을 갱신한 뒤 `--check`가 `호스트 의존성 검증 완료`를 출력했다. `frontend.json`이 없어 선택형 프런트 도구는 해당 없음.

## 릴리스 판정

| 항목 | 판정·근거 |
|---|---|
|필수: 작업 파일 위치 규칙|적용. `rules/common/README.md`를 새 자산으로 교체했다(0.3.5 → 0.3.6, 레포 고유 수정 없음).|
|필수: `local/` 제외|적용. `.fullops-squad/.gitignore`에 주석과 `local/`을 더했다. `tools/frontend/`는 선택 도구가 없어 넣지 않았다.|
|필수: 수락 기준 우선·검증 하네스 위치|적용. `rules/common/testing.md`의 두 문장을 새 판으로 바꾸고 레포 front matter를 유지했다.|
|필수: 지시서 템플릿 작업 파일 위치|적용. `handovers/_TEMPLATE.md`의 하네스 제약 문장에 로그·임시 산출물 위치를 더했다. 레포 고유 절은 보존했다.|
|project.md 규칙 묶음|적용. `fullops-common-0.3.6`.|
|레포 밖 작업 폴더|적용. 이전 작업으로 `KnowsLinkDevTools`를 기본 체크아웃 `/.devtools/`로 복사했다. 릴리스가 허용한 대로 옮기지 않고 `project.md`에 예외로 기록했다. 원본 삭제는 사용자 확인 대기.|
|lint 재사용·알림·오류 문구|레포 작업 없음. 새 세션부터 적용된다.|
|setup dry-run|생성 예정 0개.|

## 검증

- `update.py --repo . --from 1.8.6 --verify`: 필수 적용 4/4 충족, exit 0.
- 제품 코드 변경 없음.

## 보류와 동기화

- 보류: 원본 `C:/Users/shin/AppData/Local/KnowsLinkDevTools` 삭제(사용자).
- dev(V-01 진행 중)·tester(D040 대기)는 동기화를 예약한다. 유휴 designer·ops는 최신 main으로 동기화한다.
- Claude Code 세션은 `/reload-plugins` 뒤 같은 세션에서 이어간다.
