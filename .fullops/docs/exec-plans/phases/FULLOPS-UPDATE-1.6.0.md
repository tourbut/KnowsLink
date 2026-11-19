---
title: FULLOPS-UPDATE-1.6.0
status: draft
updated: 2026-10-17
owner: coor
tasks: [FULLOPS-UPDATE-1.6.0]
summary: 1.6.0의 세션 이어 하기 실측 기준을 스킬 업데이트로 받고 레포 파일은 바꾸지 않는다
---

# FULLOPS-UPDATE-1.6.0

## 버전

업데이트 전 실제 설치 버전은 1.5.0이다. 사용자가 `/plugin`으로 1.6.0을 설치했고 캐시 `fullops-squad/1.6.0`과 고정 hook 런타임 버전 1.6.0을 확인했다. 레포 적용 버전은 1.5.0에서 1.6.0으로 갱신한다.

## 설치

새 패키지 `deps.py --host claude-code` 종료0, `--check --host claude-code` 종료0(호스트 의존성 검증 완료). 이번 세션은 운영 업데이트만 수행하고 새 제품 worker는 시작하지 않는다.

## 릴리스 판정

| 항목 | 판정·근거 |
|---|---|
|1.6.0 적용1 deps|적용. 위 설치 항목.|
|1.6.0 적용2 레포 파일|해당 없음. 릴리스에 레포 파일 변경이 없다. 세션 이어 하기 실측 기준(`session_cache.py`)은 `fullops-orca` 스킬 업데이트로 적용된다.|
|setup dry-run|기존 역할·origin·main·coor·lite·standard로 생성 예정 0개. 실행할 변경 없음.|
|fullops-meter|해당 없음. 별도 플러그인이며 사용자가 이미 `/fullops`로 사용 중이다.|

## 검증

`update.py --repo . --from 1.5.0`이 릴리스 1.6.0만 출력했다. 제품 코드 변경 없음.
