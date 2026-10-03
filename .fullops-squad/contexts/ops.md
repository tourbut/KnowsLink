---
title: ops 컨텍스트
status: draft
updated: 2026-10-03
owner: ops
tasks: [SAR-DEPLOY-001-OPS, SAR-BETA-001-OPS]
summary: 운영 준비의 수락 경계와 서버 관찰 재사용 원칙을 보존한다
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
