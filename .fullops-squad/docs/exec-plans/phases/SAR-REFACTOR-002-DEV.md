---
title: 전반 감사와 기존 계약 보존 최적화
status: draft
updated: 2026-10-14
owner: dev
tasks: [SAR-REFACTOR-002-DEV]
summary: 실제 병목을 측정하고 보관·인가·수신 계약을 유지하는 최소 수정을 추적한다.
---

# SAR-REFACTOR-002-DEV

## 착수와 기술 계획

기준 main은 `8f9370bb62c9133c73e10ea2307b4cd696e10f7e`다. 착수 HEAD는 `9255852fa0e46cd40b7be1021eaa1f6c7f7d083b`이며 작업 트리는 깨끗했다. 규칙은 fullops-common-0.3.3, Test level lite, Subagent level standard다. 작성 부모 terminal은 `term_c12d7b8b-9a90-4fb9-9bf1-feaccde2d8ca`, dispatch는 `ctx_05c36bca2880`이다.

2026-10-14에 정식 `jev_packet.py --force --write-inbox`로 착수 HEAD의 packet을 갱신했다. producer partial/unknown과 오래된 후보 문서는 현재 소스와 날짜가 있는 QA 원문으로 보완한다. SAR-HUMAN-PEER-001의 후보 표현은 당시 기록이다. 현재 설명과 운영 관측은 docs/system-design.md의 서로 다른 기준 SHA를 따른다.

1. 진입점·HTTP admission·DB transaction·인가·보관을 부모가 추적한다. 코드·설정·검증 진입점과 실제 제품 상한을 감사 보고에 적는다.
2. 하위 Run `run_d02d27dce073`에서 Astra high 두 작업을 시작했다. adapter는 `adapters/**`와 패키징 두 스크립트만 소유한다. remote는 `remote_*.go`와 `text_events*.go`만 소유한다. 부모가 공유 DB·Git·문서를 통합한다. requested/effective 모두 codex/gpt-6-astra/high를 확인했다.
3. 부모의 확인된 병목은 만료된 보관 메시지마다 `current`가 회원 map을 여러 번 순회하는 것이다. 만료 처리가 이미 동일한 transport·claim 필드를 제거하므로 만료 이전에만 current를 호출한다. DB 저장 형식·권한 조회 자체·행 lock·CAS·writer guard를 바꾸지 않는다.
4. 수정 전 동일한 `BenchmarkRetainedHistory`로 100명/200 agent와 보관 메시지 100/1000/10000개의 sweep 및 transaction CPU를 측정한다. 수정 후 같은 fixture·명령을 반복한다. JSON 전체 비용과 DB 직렬화 한계는 별도로 기록한다.
5. 만료 직전/정각, 철회, delivered 보존, claim 제거, 24시간 metadata 제거 회귀를 유지한다. 최종 후보는 Linux exact SHA의 lint/test/build·race 및 격리 Postgres integration으로 검증한다.
6. 실제 영향이 있는 D03·D05–D10 원천과 현재 시스템 문서를 갱신한다. 제품 규칙·화면 외형·SQL schema가 같으면 D02·D04·D08은 no_change로 기록한다. 독립 QA·리뷰는 coordinator가 완료 SHA에서 별도 배정한다.

삭제·이동 정리, 실계정 발송, 운영 DB·배포, 의존성 변경은 실행하지 않는다. 검사 fixture와 로그는 보존한다.

## 구현·검증 결과

제품 후보는 `bcb1dd6629b508a30ec39a5d3141c8ce0a558597`다. sweep 최적화 외에 두 하위가 실제 RED로 재현한 Inbox 오류 후 본문 유실과 Remote 마지막 lease 뒤 새 도착 지연을 수정했다. 제품 3파일·회귀/벤치 4파일·테스트 등록 1파일이다. JSON 저장 형식·SQL·권한 출처/세대·wire·TTL·ACK·UI 외형은 그대로다. 전체 결함·호출자·측정과 미검사는 [감사와 QA](../../evaluations/qa-reports/SAR-REFACTOR-002-DEV/report.md)에 연결한다.

10000개 만료 이력의 Windows sweep 중앙값은 38.1ms→0.291ms다. 같은 조건 Linux 실제 DB transaction은 588.7ms→542.7ms다. 실제 병목의 일부만 줄었으며 전역 JSONB·row lock을 정규화하거나 운영 처리량을 보장하지 않는다. 같은 설치 4구독의 Remote pending 검색은 590438→162572ns/op이며 1구독 비용 증가는 하위 보고에 보존했다.

Linux 고정 제품 SHA의 install/build와 전체 실제 Postgres race integration(111.338초)은 native0이다. 최초 integration의 connect.js 부재는 build 준비 순서 실패이며 원 로그를 보존했다. 빌드를 선행 조건으로 고치고 고유 DB에서 재실행했다. 최종 문서 archive SHA의 등록 product-lint/product-test는 외부 final-lint.json과 worker_done으로 확정한다. 같은 제품 코드를 문서 변경마다 반복해 SQL·벤치하지 않는다.

startup packet은 원 history를 보존하고 정식 갱신했다. 전체 179 path/category outcomes와 producer mismatch·oversized·remaining 수동 점검을 정식 completion 검사로 확인했다. 과거 후보·high 키워드 검색은 현재 요구나 신규 QA 성공으로 승격하지 않았다. D03/D05/D07/D09/D10와 현재 시스템 문서를 갱신했고 D02/D04/D06/D08은 변경 불필요를 기록했다.

guard가 helper 생성의 삭제 API 부재 검사 문자열을 오인했다. coor가 원문을 검토한 `msg_0eb2c7322744`에 따라 비파괴 편집 도구로 helper를 작성했다. 거절된 명령을 재실행하거나 guard를 비활성화하지 않았다. 원문·판단·보존 환경 경로는 QA에 남겼다. nested Run 이후 parent check의 consumer_fenced는 child mailbox로 후속 전달을 유지했다. 원 task/dispatch는 바꾸지 않았다.

변경 줄 예산은 최소 회귀 3계열과 Windows/실제 DB 재현용 benchmark 및 필수 기술 원천 설명 때문에 400줄을 넘을 수 있다. 제품 수정 자체는 약 20줄이며 새 계층·의존성이 없다. 테스트를 억제하거나 권한/데이터 손실 회귀를 빼서 줄 수를 맞추지 않는다. 원천 문서는 서로 다른 DNN 책임과 현재 독자용 설명이므로 분할 과제로 넘기지 않았다.
