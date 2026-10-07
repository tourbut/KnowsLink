---
title: 일반 이메일 서비스 최종 main 수락과 공개 준비
status: draft
updated: 2026-10-07
owner: coor
tasks: [SAR-PUBLIC-SERVICE-OPEN-PREP]
summary: 기존 보호된 운영 배포의 백업·격리 복원과 공개 선행 조건을 보존한다
---

# 일반 이메일 서비스 최종 main 수락과 공개 준비

## 승인 범위와 현재 상태

사용자는 최종 코드의 main 병합·origin/main 반영과 공개 준비를 요청했다. Sol 독립 QA도 승인했다. 이 준비는 운영 공개 완료가 아니다. 기존 서버/Tunnel·Workers Free를 유지한다. 유료 플랜/구독/초과과금은 허용하지 않는다.

## 2026-10-07 운영 준비 실행 근거

- 운영 checkout의 실제 HEAD는0911c2c73468f8684260a277d4940a74d26bcf7d이며 clean이다. 현재 개발 main 수락 후보와 다르다. 운영 상태파일에서 SMTP_URL/MAIL_FROM 설정을 확인하지 못했다. 값/자격/메일주소를 출력하지 않았다.
- 실제운영DB를 pg_dump 형식으로 백업했다. 0911c2c-20261007T130747Z.dump,7873bytes,0600,TOC23,relay_state1행이다. 실제사용자 데이터는 이 문서에 쓰지 않는다. 운영 볼륨은 변경하지 않았다.
- 사전에 임시복원container 이름이없음을확인하고기존 restore-verify를네트워크none의자기container에서실행했다. exit0·tables2/relay_state1행이다. 임시container는회수됐고운영DB로restore하지않았다.
- local verify0,shared regression0이다. orca200/s8200/mcp401와호스트cloudflared PID506937은이전baseline과같다. public 보호상태확인도exit0이다. 이 검사는기존배포의보호유지이며새일반서비스공개성공이아니다.
- 증거는 docs/evaluations/qa-reports/COOR/open-readiness의local-state/backup-info/backup.log/restore.log/exit/operational-checks/public-closed.log이다. 비밀값과dump본문은Git에넣지않는다.

## 최종 코드 수락과 공개 적용 순서

1. DEV-FIX-3의최종고정SHA·별도정적리뷰·Sol독립QA·designer직접UI영향확인을완료한다. 모든기존실패는원본SHA로보존하며critical/high와필수실패가없어야한다.
2. coor가main에병합하고origin/main을일반push한다. main/origin동일과완료SHA의조상관계를확인한다. 실제유휴clean상설역할만sync한다. 공개후보는이원격main전체SHA로고정한다.
3. 공개적용직전운영DB의새backup/격리restore를확인한다. 현재검사는준비시점의근거이며새배포직전backup을대체하지않는다. beta.sh deploy <수락main전체SHA>는보호된배포후속의실행명령이다. 지금실행하지않았다.
4. 무료로허용된실제SMTP와발신주소·도메인설정을확인한다. KNOWSLINK_SMTP_URL/KNOWSLINK_MAIL_FROM은Git미추적상태.env에만둔다. 합성가입/시험allowlist는비어있어야한다. 실제가입메일수신을운영수락근거로확인한다.
5. 기존root Access를임의해제하지않는다. 제품의일반이메일/agent경로에맞는후보설정·정상회원접근/권한거부·공유서비스회귀를검증한후공개노출을마지막에적용한다. 이전관리token은만료기록이있어필요권한의현재유효성을적용직전에확인한다. 필요없는token재발급/구독은없다.
6. 실제노우↔다닷·외부플랫폼·실24h·운영부하/상태크기/복원·모니터링은각실제근거로따로판정한다. fixture PASS를그결과로사용하지않는다.

## 복귀 절차

공개변경실패시새경로설정만원래보호상태로복원한다. 기존root/공유Tunnel은유지한다. migration호환을확인한뒤 beta.sh deploy 0911c2c73468f8684260a277d4940a74d26bcf7d로이전코드를적용한다. 운영DBrestore는별도확인된복구판단과백업을사용한다. 공유서비스·사용자자료·운영볼륨을삭제하지않는다.

## 남은 입력과 완료 경계

실제SMTP 설정/무료제공자·발신도메인검증·실수신,필요Cloudflare설정권한·공개수락시험은미완료다. 최종코드main수락과공개준비문서/백업검증을완료하더라도이항목이없는동안현재보호된서버를공개완료로표시하지않는다.
