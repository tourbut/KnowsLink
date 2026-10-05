---
title: SAR-PUBLIC-IDENTITY-001-DEV — 일반 이메일 가입과 세션 및 자기 owner 화면을 구현한다
status: draft
updated: 2026-10-05
owner: dev
tasks: [SAR-PUBLIC-IDENTITY-001-DEV, SAR-PUBLIC-SERVICE-001]
summary: 일반 이메일 신원과 세션 및 자기 owner UI의 구현과 검증을 인계한다
---

# SAR-PUBLIC-IDENTITY-001-DEV — 일반 이메일 가입과 세션 및 자기 owner 화면

- 상태: ready. coor가 고정 기준 SHA와 새 dispatch 복귀 정보를 채운 뒤 배정한다.
- 소유: dev 역할 체크아웃의 제품 코드·검사·D03 및 해당 D05–D10·D11 사용자 안내. OPS 설정과 PLANS/board는 수정하지 않는다.
- 복귀: coordinator `term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9`, Run `run_8ca8bc058ab7`. 이 제품 과제의 기존 Task/Dispatch를 재사용하지 않는다. 실행 worker는 새 preamble의 Task/Dispatch/capability로 회신한다.

## 적용 기준과 예외

`fullops-common-0.3.2`, FULLOPS.md, orca-agents.md, project.md, document-writing.md와 coding-style/testing/security 규칙을 적용한다. 제품 기준 ref는 `6c0d132d26dac22693f1c1ffb986fc4f61bdbd09`이며 designer 준비 SHA는 `608fe06cf0d67d97a483bd088d9834c34049fb6f`다. coor는 실제 dispatch 전에 최신 main과 일반 서비스 제품 문서를 포함하는 고정 기준 SHA·후보 SHA를 이 인박스에 기록한다. 제품 문서의 review는 제품 구현 수락이 아니다. 사용자의 일반 서비스 완성 지시는 필요한 제품 구현의 승인이다. 기술 선택은 DEV/OPS가 맡는다. 실일정·결제·임의 업무 외부 발송·유료화·다닷 선행 연결은 제외한다. 최종 운영 E2E는 동일 이메일의 Grok Bot “노우”↔다닷이며 사용자에게 두 이메일을 요구하지 않는다. 다음 단계의 두 agent 식별자·키·credential은 분리한다. 비밀값·일반 이메일 전문·확인코드는 Git과 보고서에 기록하지 않는다.

## 먼저 읽을 문서

- `.fullops-squad/FULLOPS.md`, `.fullops-squad/project.md`, `.fullops-squad/rules/common/README.md`와 세 규칙.
- `.fullops-squad/docs/planning/product-specs/SAR-PUBLIC-SERVICE.md` (PS-01–14, 이번 운영 기본값).
- `.fullops-squad/docs/planning/product-specs/SAR-MVP.md` (C1–C5, frozen 업무 계약).
- `.fullops-squad/docs/design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md` (UX-01–07).
- `.fullops-squad/docs/planning/SAR-MVP-backlog.md` (현재 순서·당시 held 분리).
- `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-001.md` (판단 근거·진행 OPS 조사·후속).
- `.fullops-squad/docs/operations/ops-guide.md` (현재 owner-only 운영과 readiness 후속).


탐색 근거: `docs/evaluations/jev/SAR-PUBLIC-IDENTITY-001-DEV-find.json`, `SAR-PUBLIC-IDENTITY-001-DEV-documents-find.json`, `SAR-PUBLIC-IDENTITY-001-DEV-context.json` (HEAD `82a92f5`, 정상·fallback 없음). 문서와 코드는 별도 탐색했다.

keep 추가: `.fullops-squad/contexts/dev.md`, `README.md`, `Makefile`, `internal/relay/http.go`, `internal/config/config.go`, `internal/relay/store.go`, `cmd/relay/main.go`, `scripts/verify_mvp.py`.

필요 시 확인: `adapters/README.md`는 context의 omit? 추천이다. 실제 일반 클라이언트 연결에 필요하면 읽고 완료 보고에 기록한다.

필수 규칙과 제품 정본은 제외하지 않는다. 큰/민감 가능 passage로 원문을 미송신한 검증 스크립트도 keep이다.

### 지시 전제와 충돌 — 먼저 확인

기존 `/v1/owners`의 합성 owner 발급·owner UI Basic 인증·공개 owner-only Access를 일반 회원 신원으로 취급하지 않는다. 기존 DEC-03 수치는 당시 미승인 제안이다. 이번 D02의 새 기본값으로 구현한다. 실제 공개는 신원·자원·독립 수락 근거 전까지 차단한다. OPS의 `SAR-PUBLIC-SERVICE-OPS-READINESS`는 별도 조사로 진행 중이다. 기술 근거가 제품 수치 조정을 요구하면 coor를 통해 designer에게 전달한다. 독립 구현은 이메일 입력 대기로 중지하지 않는다. README.md/project.md의 이전 합성 범위와 현재 제품 확장 범위를 구분한다. deploy/knowslink/verify.py의 기존 공개 302 기대는 일반 이메일 positive 인증 증거가 아니다. 필요한 기술 안내/검사 갱신은 DEV/OPS가 맡는다. OPS 정정 근거 SHA는 `0313deae0dec9af813b70ee9685e1a6d0a2b84d7`다.

## 해야 할 일과 파일 소유권

- [ ] PS-01–04의 일반 이메일 가입·로그인·자기 owner 바인딩과 홈을 구현한다. 신원 확인 전 owner를 발급하지 않는다.
- [ ] 재로그인·동시 첫 가입·다른 발급자/동일 이메일·미확인 이메일의 회원 연속성과 분리를 검증한다.
- [ ] 현재/전체 로그아웃·세션 만료·이메일 재확인 복구를 구현한다. 기존 agent credential과 브라우저 세션의 수명을 분리한다.
- [ ] 신원/세션의 운영 기본값을 적용하고 발송 실패·추측·만료·재사용·NAT·재시작·한도 경계를 검사한다. 신원 제공자의 더 강한 제한은 근거와 함께 기록한다.
- [ ] 사용자 화면이 관리 오너 아이디·서버 SSH·공유 Service Auth를 요구하지 않도록 한다. 기존 owner-only 운영과 검증된 업무 gate의 회귀를 수행한다.
- [ ] 같은 회원의 다음 agent 연결과 향후 사용자별 원격 MCP/OAuth 연결을 지원할 신원 경계를 기술 계획에 적는다. OAuth 또는 신규 IdP의 무조건 선구현은 요구하지 않는다. 기존 검증된 제공자 지원을 먼저 확인한다.
- [ ] 구조/API/데이터·세션·CSRF·신뢰할 인증 근거·신원 매핑·공개 signup 우회 차단을 DEV가 정한다. 기존 synthetic 가입이 외부에서 회원 발급 우회가 되지 않는지 검증한다.
- [ ] D03과 해당 D05–D10·D11을 실제 동작에 맞춰 갱신한다. 미작성/미실행 기능을 완료로 표시하지 않는다.

## 완료 기준과 검증

DEV 완료 범위는 PS-01–04와 해당 PS-11 및 UX-01–03이다. 회원 비활성화·agent 키/관계·실메시지·다닷 연결은 다음 기능의 수락 기준으로 남긴다. 실제 허가된 일반 이메일 하나로 확인·로그인·로그아웃·재로그인을 검증한다. 다른 회원 신원의 음성 QA는 독립 fixture를 사용할 수 있다. 실제 이메일이나 외부 제공자 실행 권한이 없으면 자동 검사를 완료하고 해당 실제 관측만 blocked로 기록한다. fixture 또는 owner 계정으로 대체 PASS를 만들지 않는다.

DEV는 변경 동작 자동 검사와 인증/권한/세션·기존 frozen 업무/gate 회귀를 맡는다. project.md의 실제 `make lint`, `make test`, `make build`와 변경에 필요한 runtime/MVP 검사를 실행한다. Go/TypeScript/DB 작업 범위에 해당하지 않는 명령은 이유를 적는다. `git diff --check`, deliverables strict와 커밋 뒤 FullOps lint를 지정 기준 ref로 실행한다. 명령 자신의 종료코드와 HEAD를 보존하며 파이프로 가리지 않는다.

TESTER는 DEV 고정 후보에서 별도로 독립 QA를 수행한다. designer는 같은 후보 UX-01–03과 기존 gate 영향 화면을 직접 확인한다. 일반 이메일 사람 로그인 확인을 자동 테스트와 분리한다. coor는 다른 검토자의 별도 세션·깨끗한 read-only detached snapshot·고정 SHA 리뷰를 준비한다. 미해결 critical/high·신원 우회·필수 실패는 해당 기능 수락을 차단한다. 기능 수락이 전체 PS-13 수락이나 공개 배포 승인은 아니다.

## 갱신할 산출물과 기대 산출물

D03과 영향 있는 D05–D10, 일반 사용자 실제 안내 D11. 기술 정본·동작 검사·민감정보 없는 결과·`docs/exec-plans/phases/SAR-PUBLIC-IDENTITY-001-DEV.md`를 남긴다. OPS에 필요한 인증 제공자/발송 설정·보호 경계·지원 권한을 비밀값 없이 인계한다. 인박스에 완료 보고 전문을 쓰고 `work.py finish`로 보존한 뒤 빈 인박스·커밋 SHA를 회신한다.

## 제약·협업·후속

별도 frontend·장식 에셋·채팅 UI를 추가하지 않는다. DEV는 같은 과제에서 기술 계획과 구현을 끝낸다. 기술 계획 승인만을 이유로 작업을 멈추지 않는다. 실제 신규 비용·권한 밖 외부 발송/배포·제품 규칙 변경·설명되지 않는 검사 실패만 coor에 ask한다. 이미 승인된 실행은 재승인받지 않는다. 다음 과제는 SAR-PUBLIC-AGENTS-001-DEV이며 현재 인박스를 선점하지 않는다.

## 완료 보고

실행 worker가 실제 결과·SHA·검증·미실행·OPS/QA/UI 인계를 작성한다.
