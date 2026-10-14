---
title: Astra 리팩터링 독립 QA 실행 기록
status: review
updated: 2026-10-14
owner: tester
tasks: [SAR-REFACTOR-002-TESTER]
summary: 제품 고정 SHA·독립 경계·실패 원인·성능 범위와 완료 인계를 기록한다.
---

# SAR-REFACTOR-002-TESTER

고정 제품은 bcb1dd6629b508a30ec39a5d3141c8ce0a558597, 착수 HEAD는 8c08a1b0e69307c05f6d55a9674ac91925f036af다. 제품 경로는 동일하다. task_677998bdba45/ctx_2b4b7f958935의 별도 tester 세션에서 동작 QA만 수행했다. 제품·DEV 인박스·coor PLANS/board는 수정하지 않았다.

fullops-test의 확정 검사부터 수행했다. 현재 지시서 전문과 공통 0.3.3·project·시스템 설계·DEV 근거를 읽고 find/context/packet을 정식 갱신했다. route bind 실패는 이전 producer SHA 불일치이며 그대로 보존했다. packet의 220개 후보와 217개 비선택 쌍은 로컬 원문·검색 구간·독립 실행으로 수동 보완했다.

새 Go overlay는 queued/delivered × active/revoked 4조합의 TTL·24h·열람 권한·ACK와 callback 마지막 lease의 6개 경계를 검사한다. Node matrix는 2개 unread 본문의 반복 EACCES/EIO/EPERM 보존과 ENOENT·읽음 marker·중복·빈 Inbox를 검사한다. Linux race 최상위 16개와 Node 세 실행은 exit0이다. 실제 합성 Postgres의 관련 4개 race 회귀도 exit0이며 DB는 stop만 했다.

첫 Go 시도의 delivered completion 기대값 오류는 native1 로그와 소스를 보존했다. ACK가 received를 설정하고 sweep은 이미 완료된 상태를 유지하는 원문에 따라 기대값만 수정했다. 같은 제품 재검사는 exit0이다. 제품 결함 수정으로 기록하지 않는다.

동일 Linux CPU 벤치 3회×3표본에서 10000개 sweep은 98.0%, JSON 포함 CPU는 24.8%, 4구독 claim은 68.6% 감소했다. 1구독 claim은 8.5% 증가했다. DEV의 DB 7.8% 관측과 실행자를 구분한다. 전체 JSONB·row lock·HTTP 처리량 한계는 남는다.

상세 결과·정확한 명령·종료코드·원본 로그·미실행 범위·보존 경로는 [QA 보고](../../evaluations/qa-reports/SAR-REFACTOR-002-TESTER-test/report.md)에 있다. 새 QA 코드의 mandatory done-gate 때문에 최종 기록 HEAD에서 lint.py 등록 lint/test를 실행하며 외부 qa-refactor002/final-lint.json에 고정한다. DEV 전체 DB integration/build는 반복하지 않는다.

work.py finish로 지시서와 완료 보고 전문을 archive하고 인박스를 비운 뒤 커밋한다. 후속은 coordinator의 별도 코드 리뷰·DEV 최종 gate·통합이다. 실계정·운영·시각·장시간·백업 복원은 이번 범위 밖이다. 하위 에이전트는 사용하지 않았다.
