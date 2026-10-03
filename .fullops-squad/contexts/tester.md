---
title: tester 컨텍스트
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-SETUP-001-TESTER]
summary: 독립 QA의 FullOps lint 구분·위반 주입·복구 방식과 SAR-SETUP-001 결과를 기록한다
---

# tester 컨텍스트

결정·교훈을 항목당 3줄 이내로 기록한다.

- 기준 ref의 `lint.json` `commands`가 비어 있으면 FullOps lint는 제품 검사를 실행하지 않는다. 기준을 새 SHA로 바꾼 임시 복사본에 위반을 커밋해 등록 명령의 실패 전파(종료코드 1)를 따로 증명했다.
- 검증 명령은 셸 없이 실행하는 래퍼(`run.py`)로 실제 종료코드를 로그에 남긴다. 위반 주입은 clone의 복사본에서만 하고 `git checkout -- .`로 원복한다.
- SAR-SETUP-001: dev SHA `0cc10b0` 독립 QA는 결함 없이 통과했다. 업무 SQL·sqlc·UI는 코드가 없어 미적용으로 남긴다.
