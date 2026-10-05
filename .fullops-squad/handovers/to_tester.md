---
title: SAR-PUBLIC-IDENTITY-001-TESTER — 일반 이메일 신원과 세션 및 회원 경계를 독립 검증한다
status: draft
updated: 2026-10-05
owner: tester
tasks: [SAR-PUBLIC-IDENTITY-001-TESTER, SAR-PUBLIC-SERVICE-001]
summary: 고정 후보의 일반 이메일 신원과 세션 및 회원 분리 독립 QA를 인계한다
---

# SAR-PUBLIC-IDENTITY-001-TESTER — 일반 이메일 신원과 세션 및 회원 분리 QA

- 상태: waiting. identity DEV 완료 뒤 coor가 고정 통합 후보 SHA·별도 검증 환경·일반 QA 계정 권한을 기록한다.
- 소유: tester의 시나리오·QA 보고·자기 실행 기록. 제품 코드·인증 설정·배포·D03·PLANS/board는 수정하지 않는다.
- 복귀: coordinator `term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9`, Run `run_8ca8bc058ab7`. 새 QA dispatch의 preamble을 사용한다.

## 적용 기준과 예외

`fullops-common-0.3.2`, FULLOPS.md, orca-agents.md, project.md, document-writing.md와 coding-style/testing/security 규칙을 적용한다. 제품 기준 ref는 `6c0d132d26dac22693f1c1ffb986fc4f61bdbd09`이며 designer 준비 SHA는 `608fe06cf0d67d97a483bd088d9834c34049fb6f`다. coor는 실제 dispatch 전에 최신 main과 일반 서비스 제품 문서를 포함하는 고정 기준 SHA·후보 SHA를 이 인박스에 기록한다. 제품 문서의 review는 제품 구현 수락이 아니다. 사용자의 일반 서비스 완성 지시는 필요한 제품 구현의 승인이다. 기술 선택은 DEV/OPS가 맡는다. 실일정·결제·임의 업무 외부 발송·유료화·다닷 선행 연결은 제외한다. 비밀값·일반 이메일 전문·확인코드는 Git과 보고서에 기록하지 않는다.

## 먼저 읽을 문서

- `.fullops-squad/FULLOPS.md`, `.fullops-squad/project.md`, `.fullops-squad/rules/common/README.md`와 세 규칙.
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (PS-01–14, 이번 운영 기본값).
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (C1–C5, frozen 업무 계약).
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (UX-01–07).
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (현재 순서·당시 held 분리).
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (판단 근거·진행 OPS 조사·후속).
- `.fullops-squad/docs/operations/ops-guide.md` (현재 owner-only 운영과 readiness 후속).


탐색 근거: `docs/evaluations/jev/SAR-PUBLIC-IDENTITY-001-TESTER-find.json`, `SAR-PUBLIC-IDENTITY-001-TESTER-documents-find.json`, `SAR-PUBLIC-IDENTITY-001-TESTER-context.json` (HEAD `82a92f5`, 정상·fallback 없음). 문서와 코드는 별도 탐색했다.

keep 추가: `.fullops-squad/contexts/tester.md`, `.fullops-squad/docs/evaluations/scenarios/README.md`, `README.md`, `internal/relay/http.go`, `deploy/knowslink/verify.py`, `scripts/verify_mvp.py`, `scripts/verify_runtime.py`, `adapters/src/synthetic.ts`, `scripts/run_trial.py`, `internal/config/config.go`, `internal/relay/integration_test.go`.

필수 규칙과 제품 정본은 제외하지 않는다. 큰/민감 가능 passage로 원문을 미송신한 검증 스크립트도 keep이다.

### 지시 전제와 충돌 — 먼저 확인

기존 `/v1/owners`의 합성 owner 발급·owner UI Basic 인증·공개 owner-only Access를 일반 회원 신원으로 취급하지 않는다. 기존 DEC-03 수치는 당시 미승인 제안이다. 이번 D02의 새 기본값으로 구현한다. 실제 공개는 신원·자원·독립 수락 근거 전까지 차단한다. OPS의 `SAR-PUBLIC-SERVICE-OPS-READINESS`는 별도 조사로 진행 중이다. 기술 근거가 제품 수치 조정을 요구하면 coor를 통해 designer에게 전달한다. 독립 구현은 이메일 입력 대기로 중지하지 않는다. README.md/project.md의 이전 합성 범위와 현재 제품 확장 범위를 구분한다. deploy/knowslink/verify.py의 기존 공개 302 기대는 일반 이메일 positive 인증 증거가 아니다. 필요한 기술 안내/검사 갱신은 DEV/OPS가 맡는다. OPS 정정 근거 SHA는 `0313deae0dec9af813b70ee9685e1a6d0a2b84d7`다.

## 해야 할 일과 독립 판정

- [ ] QA-P01: 일반 이메일 하나의 실제 확인·가입·로그인과 자기 owner 바인딩을 검사한다. 미확인·만료·재사용·오답·발송 실패·위조 신원을 거부한다.
- [ ] QA-P02: 재로그인과 동시 첫 가입의 동일 회원 연속성을 검사한다. 다른 발급자 신원의 이메일 문자열 일치만으로 병합하지 않는다.
- [ ] QA-P03: 현재/전체 로그아웃·만료·브라우저 뒤로 가기·옛 세션 재사용·이메일 재확인 복구를 검사한다. agent 연결은 별도 수명임을 확인한다.
- [ ] QA-P04: 독립 fixture의 A/B owner 식별자를 바꾸어 조회·수정·gate 결정·관리자 화면 접근을 시도한다. agent credential로 owner 승인 경계를 넘지 못하는지 확인한다.
- [ ] QA-P05: PS-11의 확인 발송/오답/세션 기본값을 이하·경계·초과·동시 요청·재시작에서 검사한다. 기존 owner/합성 업무/gate 영향 회귀를 수행한다.
- [ ] QA-P06: 실제 공개 후보의 signup 우회와 현재 owner-only 경계를 검사한다. OPS 공개 준비 전에는 미실행이며 이 후보의 운영 PASS를 선언하지 않는다.
- [ ] QA-P07: UX-01–03의 실패·다음 동작 표시를 확인하고 designer의 직접 검수와 사용자 일반 이메일 사람 확인에 후보 SHA·필요 캡처를 인계한다.

## 완료 기준과 검증

DEV와 별도 세션에서 안정된 고정 후보를 검증한다. 구현자 로그는 참고이며 직접 실행한 QA와 구분한다. 명령·후보 SHA·환경/설정 식별자·기대/실제·종료코드·결함 심각도를 보고한다. 이메일·코드·token·개인키를 출력하지 않는다. 사용자에게 두 다른 이메일·인간 계정을 요구하지 않는다. 최종 운영 E2E는 동일 이메일의 공식 xAI Grok Bot “노우”↔다닷으로 별도 agent 자격·명시적 관계·앱 도구·양쪽 송수신 ID를 검증한다. 실제 이메일/계정/공개 환경이 없으면 그 항목만 미실행으로 기록하고 필요한 재개 조건을 coor에 보낸다.

QA-P01–05의 필수 실패 또는 미해결 critical/high는 identity 기능 수락을 차단한다. QA-P06의 공개 경계와 QA-P07의 사람 확인이 미완료면 일반 서비스 수락도 미완료다. 실제 메시지·키·관계·복구·다닷은 각각 후속 후보에서 PS-05–14를 추가 검증한다. 테스트 목록을 축소하여 전체 일반 서비스 PASS를 만들지 않는다.

coor는 별도 fixed-SHA 리뷰를 맡고 designer는 직접 시각 검수를 맡는다. 공유 서비스 배포 회귀는 OPS와 최종 QA 후보에서 한 번 확인한다. 변경 없는 증거는 의존성 동일성과 원래 SHA/조건을 확인해 재사용한다. 새 실패·변경 영향·증거 결함 때만 재검증한다. 불필요한 장시간 반복과 개인정보 캡처는 피한다.

## 산출물·후속

`docs/evaluations/scenarios/SAR-PUBLIC-IDENTITY-001-TESTER.md`, `docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-TESTER.md`와 자기 실행 기록을 남긴다. D02 제품 수치를 수정하지 않는다. 실패는 coor를 통해 DEV에 재현과 함께 인계한다. 완료 보고 전문을 현재 인박스에 쓰고 `work.py finish`로 보존한다. 고정 SHA·빈 인박스·QA/직접 UI/운영 후속을 새 worker_done으로 회신한다.

## 완료 보고

실행 worker가 독립 관측과 PASS/FAIL/BLOCKED·후보 SHA·남은 사람 확인 및 운영 검증을 작성한다.

## coor 고정 후보·실행 확정

- 상태 ready. DEV 결과59b66ada8b36802484cc6d7e22523257b50572cc, 제품 코드a446d89ff288c4243ad6d7f8780a778517154584. 기준94533b207b456c0560800fe30a7c90b2b5887c6e. main 수락 전 독립 QA 후보이다.
- 실제 제품 검증 checkout은 /tmp/knowslink-public-identity-qa-59b66ad, 깨끗한 detached59다. tester 상설 checkout은 지시서/기록 전용 준비 커밋이다. 코드 후보와 기록 SHA를 혼동하지 않는다. 임시 PostgreSQL·로컬 SMTP sink·loopback 서비스를 독립적으로 시작·검증·정리하며 기존 운영 자원을 변경하지 않는다.
- 새 Grok4.7high 세션이 QA-P01–07 중 자동/fixture 범위를 끝낸다. 운영 SMTP·실제 사용자 이메일·공개 후보가 아직 없어 실제 일반 이메일 사람 확인/공개 검증만 미실행으로 남긴다. fixture PASS를 실제 사람 PASS로 표시하지 않는다.
- 먼저 읽을 코드 추가: internal/relay/identity.go, member.go, mail.go, identity_test.go, identity_integration_test.go, scripts/mail_sink.py. README 실제 sink 절차와 DEV 실행 기록을 읽고 독립 검사를 설계한다. 기존 context 지도 뒤 추가된 신원 구현 파일이라는 근거다.
- 원본 DEV 검사 모두exit0은 참고다. 자신의 실제 명령/종료코드/기대·관측/실행 SHA를 보존한다. 제품 코드를 고치지 않는다. fixture 네트워크·SMTP·DB로 재시작/동시성/경계·권한을 검사한다. 실제 주소·코드·cookie를 로그나 캡처에 남기지 않는다.
