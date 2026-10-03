---
title: tester 컨텍스트
status: draft
updated: 2026-10-03
owner: tester
tasks: [SAR-SETUP-001-TESTER, SAR-MVP-001-TESTER, SAR-MVP-001-TESTER-FIX]
summary: "독립 QA의 lint 구분과 SAR-SETUP-001, SAR-MVP-001, 수정 후보 QA 결과를 기록한다"
---

# tester 컨텍스트

결정·교훈을 항목당 3줄 이내로 기록한다.

- 기준 ref의 `lint.json` `commands`가 비어 있으면 FullOps lint는 제품 검사를 실행하지 않는다. 기준을 새 SHA로 바꾼 임시 복사본에 위반을 커밋해 등록 명령의 실패 전파(종료코드 1)를 따로 증명했다.
- 검증 명령은 셸 없이 실행하는 래퍼(`run.py`)로 실제 종료코드를 로그에 남긴다. 위반 주입은 clone의 복사본에서만 하고 `git checkout -- .`로 원복한다.
- SAR-SETUP-001: dev SHA `0cc10b0` 독립 QA는 결함 없이 통과했다. 업무 SQL·sqlc·UI는 코드가 없어 미적용으로 남긴다.
- SAR-MVP-001: 후보 `a6a10c7`의 QA-01–11 실행 항목은 최종 probe 종료코드 0이다. DEC-02·DEC-03·Free N·실adapter·WAL·고의 epoch·designer 시각 판정은 held다.
- 보고서: [SAR-MVP-001-TESTER.md](../docs/evaluations/qa-reports/SAR-MVP-001-TESTER.md). 제품 코드는 수정하지 않았다. 초기 골격의 404·빈 migration 기대값은 쓰지 않았다.
- 2026-10-03 SAR-MVP-001-TESTER-FIX: 제품 `4262d02`의 신규 human transport 차단과 관련 회귀는 통과했다. legacy claim의 authorize·result 수락은 high이며 제품 수락을 차단한다.
- 판정·증거: [SAR-MVP-001-TESTER-FIX.md](../docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX.md). 원래 held와 `a6a10c7` QA는 유지했다. 제품 코드는 수정하지 않았다.
