---
title: FULLOPS-UPDATE-1.8.2
status: draft
updated: 2026-10-22
owner: coor
tasks: [FULLOPS-UPDATE-1.8.2]
summary: 1.8.1·1.8.2를 적용해 reviewer 모델 후보와 ego 실행기 문구를 레포에 옮긴다
---

# FULLOPS-UPDATE-1.8.2

## 버전

- 업데이트 전 설치: 플러그인 캐시 1.8.2(호스트 자동 업데이트, `installed_plugins.json` lastUpdated 2026-10-21T23:48Z). 고정 hook 런타임 `current.json`은 1.8.2(3e676806…)다.
- 레포 적용 버전: 1.8.0 → 1.8.2.
- 대상 버전: 1.8.2. `npx fullops-squad@latest update --hosts claude-code`가 "1.8.2 이미 목표 버전"을 확인했다.

## 설치

같은 설치기가 deps 설치와 `호스트 의존성 검증 완료`를 출력했다(ponytail 5.3.2, mattpocock-skills 1.3.1 최신). `.fullops-squad/frontend.json`이 없어 선택형 프런트 도구는 해당 없음.

## 릴리스 판정

| 항목 | 판정·근거 |
|---|---|
|1.8.1 적용1 reviewer 후보|적용. `orca-agents.md` `## 모델 후보`에 `reviewer` `claude` `claude-opus-5-5` `high`를 추가했다. 지금까지 리뷰어로 쓴 tuple과 같다.|
|1.8.1 적용2 후보 없을 때|해당 없음. 후보를 추가했다.|
|1.8.2 scenarios README|적용. 실행기 문단에 "선택 Playwright·ego", "요청 차단으로 제한할 수 없는 Orca와 ego는"을 반영했다.|
|setup dry-run|기존 역할 coor·designer·dev·tester·ops, origin·main, coor·lite·standard로 생성 예정 0개.|

## 검증

- `update.py --repo . --from 1.8.2`가 릴리스 1.8.1·1.8.2를 출력했다(적용 기준 1.8.0).
- FullOps lint를 변경 범위에 실행했다(커밋 기록 참조). 제품 코드 변경 없음.

## 보류와 동기화

- 보류 없음. 1.8.0의 진행 중 패킷 재생성 보류는 진행 중 인박스가 없어 해소됐다.
- 활성 worker가 없으므로 모든 역할 워크트리를 최신 main으로 동기화한다.
- 같은 세션 계속: 이 세션의 스킬은 이미 1.8.2 경로에서 로드됐다. hook 런타임도 1.8.2다.
