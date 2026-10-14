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
