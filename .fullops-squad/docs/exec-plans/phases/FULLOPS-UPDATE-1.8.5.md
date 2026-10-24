---
title: FULLOPS-UPDATE-1.8.5
status: draft
updated: 2026-10-24
owner: coor
tasks: [FULLOPS-UPDATE-1.8.5]
summary: "1.8.4·1.8.5를 확인했다. 두 릴리스 모두 레포 적용 항목이 없어 plugin_version만 1.8.5로 올린다"
---

# FULLOPS-UPDATE-1.8.5

## 버전

- 업데이트 전 설치: 1.8.5(호스트 자동 업데이트, `installed_plugins.json` lastUpdated 2026-10-24T04:26Z).
- 레포 적용 버전: 1.8.3 → 1.8.5.
- 대상 버전: 1.8.5. `claude plugin marketplace update` 뒤 `claude plugin update`가 "already at the latest version (1.8.5)"을 확인했다. fullops-meter는 0.10.0 최신.

## 설치

`deps.py --host claude-code`로 의존성을 갱신했다(exit 0). `deps.py --check --host claude-code`가 `호스트 의존성 검증 완료`를 출력했다(ponytail 5.3.2, mattpocock-skills 1.3.1 최신). `.fullops-squad/frontend.json`이 없어 선택형 프런트 도구는 해당 없음.

## 릴리스 판정

| 릴리스 | 판정·근거 |
|---|---|
|1.8.4 흐름 그래프(경계 알림·반복 예산·과제 상태·Flow 탭)|해당 없음. 릴리스가 "할 일 없음"으로 명시한다. hook·스크립트 변경이라 새 세션에서 적용된다. 선택 항목 `반복 예산: N`은 기존 합의가 없어 쓰지 않는다.|
|1.8.5 `cd`·`Set-Location` 뒤 병합·커밋 검사|해당 없음. 릴리스가 "할 일 없음"으로 명시한다.|
|setup dry-run|기존 역할로 실행, 생성 예정 0개.|

## 검증

- `update.py --repo . --from 1.8.5 --verify`: 필수 적용 0/0 충족, exit 0.
- 제품 코드 변경 없음.

## 보류와 동기화

- 보류 없음.
- dev 워크트리는 SAR-AUTOSETUP-001-V01 진행 중이라 동기화를 예약한다. 그 밖의 유휴 역할 워크트리는 최신 main으로 동기화한다.
- Claude Code 세션은 `/reload-plugins` 뒤 같은 세션에서 이어간다.
