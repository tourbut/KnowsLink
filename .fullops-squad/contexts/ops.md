---
title: ops 컨텍스트
status: draft
updated: 2026-10-06
owner: ops
tasks: [SAR-DEPLOY-001-OPS, SAR-BETA-001-OPS, SAR-PUBLIC-AGENTS-001-REVIEW, SAR-PUBLIC-MESSAGES-001-REVIEW]
summary: 운영 준비의 수락 경계·서버 관찰 재사용·리뷰 재현 원칙을 보존한다
---

# ops 컨텍스트

결정·교훈을 항목당 3줄 이내로 기록한다.

## 2026-10-03 — SAR-DEPLOY-001-OPS

D12는 계획이며 실제 배포는 미실행이다. 사용자 중지 지시 이후 수락 SHA만으로 운영을 재개하지 않는다.
기존 a6a10c7에는 업무 SQL·인증·gate UI가 있다. 초기 골격 관찰과 최신 제품 정본을 구분한다.
상세 근거는 [실행 기록](../docs/exec-plans/phases/SAR-DEPLOY-001-OPS.md)에 보존한다. proxied DNS는 CNAME 질의만으로 실패 판정하지 않는다.

## 2026-10-03 — 기록 마무리 재개

coordinator msg_d89612ea9341은 기존 아카이브·빈 인박스 보존과 검증·커밋만 지시했다. 중복 finish와 서버 재조사는 수행하지 않는다.
a6a10c7은 deliver:human C1 high로 수락·배포 금지다. DEV 수정과 독립 QA·리뷰 후 새 수락 SHA가 필요하다.
현재 Dispatch는 ctx_3c54fe7d2043이다. 역할 기록 완료는 제품 수락과 구분한다.

## 2026-10-03 — SAR-BETA-001-OPS

수락 후보 557ebc3/제품 78b1d92로 a6a10c7 금지는 역사적 기록이 됐다. 구성은 `deploy/knowslink/`, 상태·비밀은 `/home/shin/deploy/knowslink-state`(0700)다.
Cloudflare API MCP는 읽기 전용이다. `cert.pem`은 Tunnel 생성 가능, Access 쓰기는 별도 토큰이 필요하다. Access 앱 확인 전에는 DNS·connector를 열지 않는다.
reusable policy는 새로 만든다. 기존 `knowslog-bot - Production`은 다른 이메일을 허용하므로 재사용하지 않는다. OTP는 이메일 제한과 함께만 쓴다.

## 2026-10-03 — SAR-BETA-001-OPS 연결 적용

구성 28bd1bb로 Access 앱·정책·DNS·connector를 적용했다. 관리 쓰기는 Codex file-store OAuth(공식 MCP)를 일회성 bridge로 사용했다. 로컬 resolver의 NXDOMAIN 음성 캐시는 `dig @1.1.1.1`+`curl --resolve`로 우회한다.
사용자 이메일 로그인은 인간 검사로 남는다. 코드 변경 SHA에서 `expose`를 다시 하기 전 독립 리뷰와 `access_apply.py check`가 필요하다.

## 2026-10-03 — SAR-MVP-002-DEV-REVIEW

고정 SHA 552586b를 구현자 Codex 세션과 다른 세션에서 독립 검토했다. snapshot은 읽기 전용이므로 재현 build는 `git archive` 임시 사본에서 한다. `node_modules`를 symlink하면 bundle 주석 경로가 달라져 ZIP SHA가 바뀐다. 경로 정규화 후 동일성으로 판정하고, SHA 일치를 주장하려면 실제 디렉터리에서 `npm ci`를 쓴다.
공식 Grok Bot connect 문서는 stdio MCP·ZIP 업로드·Node runtime을 언급하지 않는다. 설치 지원을 주장하지 않는다. 결과는 critical/high 0이며 후속 F-01 호스트 tool timeout 확인이 실제 연결 재개 조건이다.

리뷰 findings의 줄 번호는 파일별로 `cat -n`/`nl`을 따로 실행하거나 `sed -n`으로 확정한다. 여러 파일을 한 번에 출력하면 번호가 누적돼 존재하지 않는 줄이 기록된다(SAR-MVP-002-DEV-REVIEW F-01 보정). 기록 전에 `wc -l`로 범위를 대조한다.

## 2026-10-04 — SAR-MVP-003-BIDIRECTIONAL-OPS

배포 `0911c2c`와 시험 allowlist 적용, 실제 key loopback 왕복은 완료했다. Cloudflare service token 쓰기 권한은 저장된 어떤 grant에도 없다. 세션 cloudflare MCP는 읽기 전용(쓰기 1010)이고 Codex file-store OAuth에는 `access-service-token.*` scope가 없다. 빈 본문 POST는 자원을 만들지 않고 권한만 확인한다. 만료된 OAuth는 refresh하지 않는다(다른 도구 credential 회전). 시험 relay의 idempotency key는 16~128자 ASCII다. 상세는 phases/SAR-MVP-003-BIDIRECTIONAL-OPS.md.

## 2026-10-06 — SAR-PUBLIC-AGENTS-001-REVIEW

단일 JSONB 상태에서는 활성 한도와 별도로 철회 기록의 보존량을 확인한다. 활성 개수 한도만으로 상태 크기가 제한되지 않는다. 핵심 인가 테스트가 통합 전용이면 reviewer scratch에서 verify-mvp를 직접 재실행한다.

## 2026-10-06 — SAR-PUBLIC-MESSAGES-001-REVIEW-2

전역 동시 슬롯은 획득 시점이 인증·rate 앞이면 비인증 slow body가 비용 없이 슬롯을 묶는다. 슬롯을 잡은 뒤 본문을 읽는 순서를 먼저 본다. DB 없는 `(&Service{}).Handler()`와 raw TCP 연결로 재현할 수 있다. 정리 경로 분류는 경로뿐 아니라 인증 principal 여부도 확인한다.
review.py check는 미해결 high에서 첫 실패로 멈춘다. 나머지 기록 조건은 result.json 사본에서 high만 임시 resolved로 바꾼 probe로 확인하고 바로 복원한다. probe 결과를 수락 근거로 쓰지 않는다.
