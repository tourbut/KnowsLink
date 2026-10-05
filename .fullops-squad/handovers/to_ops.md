---
title: SAR-PUBLIC-SERVICE-OPS-READINESS — 일반 이메일 서비스의 운영 신원·자원·복구 근거를 읽기 전용으로 확인한다
status: draft
updated: 2026-10-05
owner: ops
tasks: [SAR-PUBLIC-SERVICE-OPS-READINESS]
summary: 일반 이메일 서비스의 운영 신원·자원·복구 근거를 읽기 전용으로 확인한다
---

# SAR-PUBLIC-SERVICE-OPS-READINESS — 일반 이메일 서비스 운영 준비 근거

## 현재 상황과 기준

사용자는 일반 이메일 가입부터 실제 에이전트 양방향 메시지까지 일반 서비스 수준으로 완성한 뒤 OpenAI dot “다닷”을 연결하라고 지시했다. 제품 기준은 designer task_c21892048d9c/ctx_000841343790에서 준비 중이다. 지금은 제품 결정에 필요한 독립 운영 근거를 읽기 전용으로 제공한다. 기준 main `6c0d132d26dac22693f1c1ffb986fc4f61bdbd09`, 현재 제품 배포 `0911c2c`. trial 인증은 종료됐다.

## 적용 기준과 예외

fullops-common-0.3.2, FULLOPS.md, project.md, orca-agents.md, 문서 작성·coding-style/testing/security 규칙을 적용한다. FullOps 업데이트 제외. 실제 서버 `/home/shin/deploy/knowslink`·state, CF API/브라우저의 읽기 전용 조사만 승인 범위다. secret·이메일·사용자 데이터는 출력하지 않고 Python 메모리 헤더로만 인증한다. 현재 관리 API token 파일은 state/cf-service-token-api.env(0600), 자동 만료 2026-10-05T23:59:59Z다. 기존 구성을 변경하거나 새 token을 발급하지 않는다. GET 권한 부족은 필요 권한과 대체 관측을 기록한다.

## 해야 할 일과 소유권

- [ ] OPS 가이드·CLEANUP·beta/access/verify 도구를 읽어 현재 root owner 보호·Tunnel·IdP의 일반 이메일 지원 방식과 제한을 확인한다. 합성 owner token과 실제 이메일 신원을 구분한다. CF 설정 쓰기·인증 해제·agent secret 발급·메일 발송은 실행하지 않는다.
- [ ] CPU RAM 디스크·현재 knowslink/공유 서비스 자원·DB 크기와 보존/백업/복구 준비를 읽기 전용으로 관측한다. 다른 서비스 보호를 포함한 운영 한도 근거·미측정 항목·공개 전 DEV/OPS 검사 조건을 제공한다. 임의 성능 보장이나 제품 수치 확정은 하지 않는다.
- [ ] 일반 사용자 브라우저 로그인과 agent/원격 MCP 연결을 위한 배포 계약·원점 검증·경로별 보호·필요 권한·rollback 선행과 순서를 짧게 제안한다. 구현 설계/제품 결정은 DEV/designer에 넘긴다. 기존 root를 everyone/bypass로 바꾸는 조사는 적용까지 진행하지 않는다.
- [ ] 자기 실행 기록 `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md` 및 D12 운영 가이드의 별도 readiness 절만 갱신한다. designer와 제품 문서/PLANS/board 소유권은 겹치지 않는다.

## 먼저 읽을 문서

Jev code/documents-find·context 정상 실행, 후보와 필수 문서는 모두 keep이다. beta.sh는 민감/큰 구간을 외부 분류로 보내지 않고 필수로 유지했다. 필수 `.fullops-squad/docs/operations/ops-guide.md`, `.fullops-squad/docs/exec-plans/phases/SAR-MVP-003-TRIAL-CLEANUP.md`, `.fullops-squad/docs/planning/product-specs/SAR-MVP.md`, `deploy/knowslink/access_apply.py`, `deploy/knowslink/verify.py`, `deploy/knowslink/beta.sh`. 미추적 비밀 파일은 검색/외부 Jev 후보에 넣지 않는다.

## 완료 기준과 산출물

갱신 D12; D13은 현재 일반 서비스 검증 전이라 수락으로 갱신하지 않는다. 실제 관측·미확인·다음 필요권한/외부입력·배포 rollback 조건을 구분한 기술 근거가 있으면 과제 완료다. 서버 변경이 없으므로 제품 전체 QA를 재실행하지 않는다. 변경 문서 링크·diff check·clean commit 후 lint를 검증한다. 받은 인박스에 전문 작성→work.py finish로 archive/빈 인박스→고정 SHA와 reportPath 및 짧은 관측을 run_8ca8bc058ab7에 worker_done으로 직접 보낸다. 새로운 제품 정책·외부 수정은 수행하지 않는다. 다음 실제 DEV/OPS/독립QA는 coor가 배정한다.

## 완료 보고

worker가 전문을 적는다.
