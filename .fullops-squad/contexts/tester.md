---
title: tester 컨텍스트
status: draft
updated: 2026-10-14
owner: tester
tasks: [SAR-REFACTOR-002-TESTER]
summary: 독립 QA의 SHA·종료코드·재사용 경계와 최근 리팩터링 교훈을 보존한다
---

# tester 컨텍스트

10KB를 넘은 이전 전문은 [2026-10-14 archive](archive/tester-2026-10-14-refactor002.md)에 복사했다. 파일 삭제·이동은 하지 않았다. archive의 원래 상대 링크는 이 contexts 디렉터리를 기준으로 읽는다. 과거 실패·held·수락은 각 보고서의 당시 SHA에 남으며 최신 QA로 덮지 않는다.

- tester는 제품을 고치지 않는다. 합성 fixture·테스트·시나리오·증거를 소유하고 실패는 재현 절차로 dev/coor에 전달한다.
- 고정 제품 SHA와 기록 HEAD를 구분한다. 제품·의존성 동일성을 확인한 DEV 증거만 재사용하며 새 독립 실행으로 바꾸어 기록하지 않는다.
- 실제 명령 종료코드를 보존한다. 환경·기대값·증거 결함의 실패도 남기고 원인이 바뀐 뒤에만 재검증한다.
- 기준 lint.json에 등록 명령이 없으면 lint 성공이 제품 실행 증거가 아니다. Windows 미지원 도구는 exact HEAD·clean을 검사하는 고유 Linux clone에서 실행한다.
- SSH는 기존 known_hosts/RejectPolicy·비밀값 가림을 유지한다. DEV HEAD를 올리는 make.cmd를 tester 고정 clone에 재사용하지 않는다.
- Git Bash 원격 경로 전달은 MSYS_NO_PATHCONV=1을 확인한다. test_public_config.sh의 stat -c는 Linux 전용이다.
- packet과 실제 task/role/attempt/SHA를 맞춘다. 가짜 task-key로 검사기를 우회하지 않는다. 과거 archive의 해당 우회 조언은 2026-10-11 coor 정정으로 폐기됐다.
- shell_guard 거부를 다른 셸·표현으로 우회하지 않는다. 삭제·이동 정리는 현재 AGENTS와 coordinator 권한을 따른다.
- Codex hook 검증은 실제 host 실행 방식과 종료코드2·deny JSON을 함께 확인한다. 문자열 리터럴/argv cmd 호출의 성공만으로 guard 실행을 주장하지 않는다.
- Orca 브라우저 viewport는 goto 뒤 적용하고 실제 innerWidth를 측정한다. 화면·CLI·실계정·호스트 표시 증거는 각각 구분한다.
- 웹 사본 보관과 활성 transport 인가는 별도다. 관계 철회 뒤 참여자 이력이 남아도 과거 전송·claim·세대 권한이 되살아나면 안 된다.

## 최근 결과

- 2026-10-13 SAR-REFACTOR-001-TESTER: exact cd4dae8의 Linux race8개·MCP/owner-text·독립 text15개 exit0. [보고](../docs/evaluations/qa-reports/SAR-REFACTOR-001-TESTER-test/report.md). 제품 변경0, 실계정·운영·SQL 이행은 미실행이다.
- 2026-10-14 SAR-REFACTOR-002-TESTER: exact bcb1dd66 Linux race16개·새 Go 경계10개·Inbox 오류 matrix·합성 DB race4개 exit0. 제품 변경0, 새 C/H 발견0. [보고](../docs/evaluations/qa-reports/SAR-REFACTOR-002-TESTER-test/report.md).
- ACK completion received는 TTL 뒤에도 유지한다. 첫 잘못된 기대값/native1을 보존하고 기존 계약으로 바로잡았다. sweep CPU 이득은 DB/HTTP 처리량이 아니다. 단일 구독 claim 증가와 local claim crash window는 남는다.
