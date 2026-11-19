---
title: Google 인증 변경 독립 검토
status: review
updated: 2026-10-07
owner: coor
tasks: [SAR-GOOGLE-LOGIN-001-DEV, SAR-GOOGLE-LOGIN-001-REVIEW]
summary: 고정 Google 인증 변경을 수락하고 실제 로그인 확인을 기존 서버 배포 뒤 진행한다
---

# Google 인증 변경 독립 검토

검토자는 Codex GPT-6.1-Sol medium coordinator 세션 `01a10ed4-154b-7572-9160-7630e9ac42bc`다. 구현자는 실제 DEV 세션 `01a116ca-e458-7550-a8f7-3a87fd65ef1c`다. 별도 managed clean detached snapshot `017bf456c3ab0f7367b1e3980b48f949ce250856`을 읽기 전용으로 확인했다. base/merge-base는 `abf2de1e3f0501e56b41f9d909514f591b656d4c`다. OCR v1.12.11은 대상·규칙 계산만 수행했다. AI 검토이며 자동 코드 판정이 아니다.

## 범위와 판단

전체 26개 (path,status)를 reviewed로 확인했다. skipped 0이며 coverage 100%다. `rules.json`의 Markdown/JSON/Go/YAML 규칙, `fullops-common-0.3.3`의 코딩·테스트·보안, project.md, archive 지시서의 범위·최신 사용자 Google 우선 출시 결정을 적용했다. 고정 SHA와 현재 공통 규칙은 동일하다.

code exchange/issuer/audience/RS256/expiry 검증, nonce·최근 iat·검증 이메일, 브라우저 state cookie와 PKCE 및 단일 소비를 확인했다. 소비를 value로 반환하여 실패/취소도 transaction commit되는 경계를 확인했다. token exchange는 SQL lock 밖에서 수행하고 timeout을 적용한다. verifier의 생명주기는 startup context 취소 뒤 유지하되 HTTP timeout을 보존한다. issuer/sub 신원 분리, 이메일 자동 병합 거부, 같은 Google 회원·시작 세션에 제한된 재확인, owner 활성/회원 한도, 기존 Strict session·CSRF·logout·TTL과 sweep의 JSON 호환을 확인했다. 새로운 DB migration은 없다.

새 critical/high 발견은 없다. 기존 M-UTF8-01과 전체 QA/UI 미완료를 변경하거나 PASS로 처리하지 않는다. OAuth testing 및 기존 Access 보호 때문에 일반 공개 수락은 아니다. 실제 Google 로그인과 callback 접근은 배포 뒤 사용자가 직접 확인한다.

## 검증 근거 재사용

`lint.json`은 DEV Git directory의 최종 clean 017bf456 결과를 그대로 복사했다. base/head/config hash를 확인했다. product-lint(kind lint) 및 product-test(kind test) 모두 exit 0, ERROR0/WARNING6/실행불가0다. 새로운 SHA 실행이나 독립 기능 QA로 표시하지 않는다. DEV 좁은 signed-provider/Postgres 인증 및 이메일 회귀·빌드 exit0은 실행 기록과 authentic 완료 전문으로 확인했다. 추가 전체 검사·캡처·부하 시험은 사용자 지시에 따라 실행하지 않았다.

SIZE-002 실제 추가664줄은 인증 거부 및 실제 서명/Postgres 회귀 약300줄과 설정·인계가 포함된 규모다. 범위는 단일 Google 인증이다. 기존 큰 파일 분할 리팩터링은 생략한다. DEP-001의 coreos/go-oidc v3.17.0/oauth2 v0.36.0/go-jose v4.1.3은 기존 OIDC 검증기가 없는 코드에서 사용자 token 검증 재발명을 피한다. 공식 문서와 고정 버전 소스 근거는 DEV 기록을 재사용했다. SEC-001의 local-fixture는 합성 provider 더미이며 운영 자격이 아니다. 나머지 경고는 기존 큰 파일의 좁은 변경으로 수락한다.

UI는 기존 memberStyle/template을 재사용한다. 새 색·테마·컴포넌트가 없고 전용 디자인 lint/테마 전환은 프로젝트에 미구성이다. Strict cookie를 유지하는 완료 화면 홈 링크는 실제 브라우저 확인이 남았다. Google 로그인 버튼·문구 외 전체 UI 검수는 재개하지 않는다.

## 대화 미참조 인계 점검

산출물 README에서 D02/D10/D12/D13 정본을 찾고 README Google 절과 실행 기록·archive를 대조했다. 현재 요구/구현/검증/다음 담당은 확인했다. 과거 D02 이메일 전용 설명 및 D12/D13 Google env 안내는 최신 결정 보완이 필요하다. coor가 배포 기록에 최신 우선순위·설정 경로·복구 SHA를 추가한다. README→D10/실행 기록 로컬 경로는 존재한다. Google 절의 정확한 callback과 env 이름이 실제 등록 메타데이터와 같다. snapshot 정리 뒤에도 result/report/lint/preview/rules와 exact SHA를 정본 폴더에 보존한다.

## 결론

고정 인증 후보를 수락한다. 필수 정적 기록 확인 뒤 main/origin/main에 통합하고 기존 서버에 배포한다. Google 테스트 사용자 knowslog01@gmail.com의 실제 로그인·자기 홈은 사용자와 확인한다. 공개 Access 전환·Google publishing·Grok Bot 실연결은 이 정적 수락의 PASS가 아니다.
