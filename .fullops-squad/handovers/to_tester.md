---
title: SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER — 최종 후보 PS08–11과 수정된 정상 사용자 동작의 독립 QA
status: draft
updated: 2026-10-07
owner: tester
tasks: [SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER]
summary: 최종 후보 PS08–11과 수정된 정상 사용자 동작의 독립 QA
attempt: 17a764b1de5c430b9235b82643c1b00d
base: 68b0d6a0c854fdaec6828a232dd3945814be1404
---

# SAR-PUBLIC-MESSAGES-001-FIX-3-TESTER — 최종 후보 PS08–11과 수정된 정상 사용자 동작의 독립 QA

## 목적과 파일 소유권

최종 고정 후보의 PS08–11과 변경된 정상 사용자 흐름을 독립 검증한다. 구현자/최종 reviewer와 다른 실제 세션을 사용한다. 결과는 본인 QA 보고서·시나리오/테스트 증거·PLANS/context/inbox에 쓴다. 제품 코드는 수정하지 않는다.

## 해야 할 일과 검증

- [ ] 원본 TESTER09c 최종 보고서·실패 항목과 DEV-FIX-3 변경 영향을 연결한 기대 동작 표를 만든다. 이전 Grok 실패/중간 프로브 오류/미실행은 소급 변경하지 않는다.
- [ ] 이메일 fixture 일반 회원→서로 다른 owner·agent→text 요청/관련 회신→persist/ACK→receipt의 정상 흐름과 PS08–11의 TTL/멱등·권한/CSRF·상한·재시작/다중 인스턴스 계약을 기존 로컬 테스트/QA 도구로 검증한다. 공유 admission/store 경계의 영향 회귀를 포함한다. 소유 앱의 격리 fixture만 사용하고 새 외부 대상/공격·부하 도구는 만들지 않는다.
- [ ] DEV-FIX-3의 동일 DB 시각 snapshot 순서·다른 인스턴스 상태 갱신·신규 budget 포화 중 자기 유효 정리·각 cleanup 호출자의 정상 동작을 독립 판정한다. 관측 테스트 exit0만으로 기대값 만족을 주장하지 않고 관측값과 요구값을 함께 확인한다.
- [ ] 원 TESTER H1/M1/HTTP 슬롯 잔류 실패의 새 후보 해소 근거를 기존 테스트와 현재 결과로 확인한다. 회귀/추가 코드 결함은 증거·재현 조건·심각도로 DEV에게 인계한다. 실패 기대값을 근거 없이 낮추지 않는다.
- [ ] make lint/test 및 필요한 기존 integration/verify-mvp 검사를 자기 scratch에서 실행하고 명령 자신의 exit·candidate SHA·fixture/자식 프로세스 회수를 남긴다. 변경 없는 증거는 관련 코드/설정 동일성을 확인하고 원 실행 SHA로만 재사용한다. 전체 API/PS08–11 검수와 좁은 FIX3 검증을 구분한다.
- [ ] fixed 후보의 QA 판정·잔여 실패/미검증·독립성·원본 무결성을 보고서에 남긴다. designer 직접 시각 검수는 별도이며 값/상태는 코드로 판정하고 불필요한 캡처/영상은 만들지 않는다.

## 먼저 읽을 문서와 갱신할 산출물

필수 공통 기준·제품/UX 정본·DEV-FIX-3 실행 기록/최종 검사·원 TESTER 보고서/기존 시나리오·FIX-REVIEW 및 최신 UI 보고서를 읽는다. code/doc/context/packet은 최종 후보 고정 뒤 coordinator가 연결한다. D10의 이번 QA 시나리오/보고서 연결만 갱신하며 운영 공개 PASS는 쓰지 않는다. 기대 산출물은 이 키.md와 이 키-test의 명령/exit·판정·fixture 회수·원본 무결성이다.

## 대상·적용 기준·복귀

- 상태: blocked.DEV-FIX-3 최종 SHA 확정 전 착수하지 않는다. 후보 SHA와 snapshot은 coordinator가 dispatch 전에 고정한다.
- 공통 기준 fullops-common-0.3.3의 README/coding-style/testing/security, FULLOPS·project·document-writing·orca-agents·역할 context와 review/rule.json을 따른다. 기준 main은 68b0d6a0c854fdaec6828a232dd3945814be1404다. 외부 새 SDK/의존성은 없으며 기존 버전 근거를 재사용한다.
- 사용자 승인 모델: fresh Codex gpt-6.1-sol high. 과제별 지정이며 전역/역할 전체 설정·구독·추가 결제를 변경하지 않는다. 실제 세션 ID와 fixed SHA·실행 위치·명령별 종료코드를 남긴다.
- 복귀 repo 818c78e5-d51c-4ff4-aa88-70e9ee185fbb, coor /home/shin/orca/workspaces/KnowsLink/fullops-coor, terminal term_a8a1fa04-50ab-448d-94e7-11e8ee3c77f1, Run run_8ca8bc058ab7. 실제 Task/Dispatch 권한은 새 preamble을 따른다.
- 제품 기준은 docs/planning/product-specs/SAR-PUBLIC-SERVICE.md PS08–11/PS04·06·07과 docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md UX06·07이다. 원 TESTER09c 실패·OPS dfc H2 high·UI09c 실패·FIX2 리뷰 pending·플랫폼 차단은 원 SHA/시점으로 보존한다.

## 제약과 완료 기록

Workers Free·기존 서버/Tunnel을 유지한다. 운영 공개·실메일·외부 계정/플랫폼·사용자 자료 삭제·유료 전환은 범위 밖이다. 검사 실행은 자기 격리 scratch·fixture에서만 한다. 원본 snapshot은 detached clean read_only로 유지하고 기록은 역할 checkout에 쓴다. 원본 실패 파일은 수정하지 않는다. UI 화면 변경이나 제품 quota 판단은 designer에게, 제품 코드 결함은 같은 DEV 후속으로 coordinator에게 전달한다.

완료 보고 전문·실제 fixed SHA·세션·통과/실패·미실행·후속을 남긴다. work.py finish로 archive/빈 inbox, 마지막 기록 SHA에서 FullOps lint/test·strict·diff 검사, 역할 브랜치 일반 push와 실제 worker_done을 완료한다. 실패도 증거와 함께 보고하며 성공으로 바꾸지 않는다. synthetic/local 결과는 실메일/공개/실24h/노우↔다닷/운영 부하와 구분한다. 끝난 뒤 idle이며 다음 과제를 시작하지 않는다.

## 완료 보고

작업 종료 뒤 실제 결과 전문을 작성한다.
