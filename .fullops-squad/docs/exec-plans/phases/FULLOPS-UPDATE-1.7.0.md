---
title: FULLOPS-UPDATE-1.7.0
status: draft
updated: 2026-10-18
owner: coor
tasks: [FULLOPS-UPDATE-1.7.0]
summary: 1.7.0의 사용자 결정 기록과 템플릿 기반 모션 경로를 스킬 업데이트로 받고 레포 파일은 바꾸지 않는다
---

# FULLOPS-UPDATE-1.7.0

## 버전

업데이트 요청 시점의 실제 설치 버전은 이미 1.7.0이다(`claude plugin list`, 캐시 `fullops-squad/1.7.0`). `claude plugin marketplace update fullops-squad`와 `claude plugin update fullops-squad@fullops-squad`가 최신 1.7.0을 확인했다. 레포 적용 버전은 1.6.0에서 1.7.0으로 갱신한다. 이 coordinator 세션은 1.6.0 스킬로 시작됐다.

## 설치

새 패키지 `deps.py --host claude-code` 종료0(ponytail 5.3.2, mattpocock-skills 1.3.1 최신 확인), `--check --host claude-code` 종료0(호스트 의존성 검증 완료). 고정 hook 런타임 `current.json`이 1.7.0 패키지(4c24ea22…)를 가리킨다. 이번 세션은 운영 업데이트만 수행하고 새 제품 worker는 시작하지 않는다.

## 릴리스 판정

| 항목 | 판정·근거 |
|---|---|
|1.7.0 적용1 deps|적용. 위 설치 항목.|
|1.7.0 적용2 레포 파일|해당 없음. 결정 기록(`decision.py`)은 Git 공용 디렉터리 `fullops-decisions/`에 처음 기록할 때 생기며 작업 트리를 바꾸지 않는다. 다음 coordinator 세션부터 `fullops-orca` 규칙으로 사용한다.|
|1.7.0 적용3 fullops-meter|적용 확인. `claude plugin list`에서 fullops-meter 0.7.0.|
|1.7.0 적용4 모션 렌더 도구|해당 없음(선택). 기존 `videos/knowslink-motion/`이 자체 package.json으로 playwright를 쓰고 FFmpeg 8.0.1이 PATH에 있다. 공통 설치에 추가하지 않는다.|
|setup dry-run|기존 역할·origin·main·coor·lite·standard로 생성 예정 0개. 실행할 변경 없음.|

## 검증

`update.py --repo . --from 1.7.0`이 릴리스 1.7.0만 출력했다. 제품 코드 변경 없음.

## 보류

없음. 열린 사용자 결정(모션 영상 청취·시각 수락)은 새 세션에서 `decision.py add`로 기록할 수 있다.
