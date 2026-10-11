---
title: SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS — 독립 검수·main 수락 후 OAuth 브라우저 복귀 수정을 기존 서버에 보존 배포
status: draft
updated: 2026-10-11
owner: ops
tasks: [SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS]
summary: 독립 검수·main 수락 후 OAuth 브라우저 복귀 수정을 기존 서버에 보존 배포
attempt: f5a42297abfc45748fe3746cb2fc975e
base: 8cfebf665520e3e8e5dae275442f51f289adf1aa
subagent_level: standard
test_level: lite
---
# SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS — 검수된 OAuth 복귀 수정의 기존 서버 보존 배포

- Purpose implementation(운영). 상태: DEV 최종 SHA·독립 tester 수락·main push 대기. 선행 조건 충족 전 착수하지 않는다.
- ops C:/Users/shin/orca/workspaces/KnowsLink/ops, fullops/ops. fresh Claude Sonnet5.5 high(Jev). test lite/subagent standard.
- coor 복귀 term_e61d3e14-29e9-4954-943a-4a75707c82de / run_dc9b4e6dd18a. task/dispatch는 preamble을 따른다. 기록 main 병합·push는 coor 담당이다.

## 승인 범위와 보존

기존 .env.server Linux 서버에서 승인된 Google/Dots 연결의 OAuth 브라우저 복귀 오류만 보완 적용한다. 사용자가 기존 서버 배포와 테스트를 승인했다. 추가 승인을 반복 요청하지 않는다. 기존 비용 안에서 진행한다.

로컬 Docker·새 유료 API/서비스/billing·삭제/이동/cleanup/복구·force reset/clean은 금지한다. F: 사고 자료·private backup·DB·키·agent·pair·Google 설정·Tunnel credential·공유 서비스를 보존한다. 실계정 승인·외부 메시지·새 설치 생성은 부모가 배포 후 수행한다. 이번 worker는 수행하지 않는다.

## 원천과 적용 기준

fullops-common-0.3.3과 FULLOPS, common README+세 규칙, project, orca-agents, delegation, document-writing을 읽는다. 원천은 현재 ops-guide/transition, SAR-SSE-INSTALL-001-OPS-LIVE, 이전 OPS 실행 기록과 logs/2026-10-11_to_ops.md, DEV 최종 OAuth 복귀 실행 기록 및 tester 수락 보고다. 고정 제품 SHA와 정확한 gate 경로는 선행 수락 뒤 coor가 연결한다. 기존 Context7/공식 CLI 근거는 동일 명령에 재사용한다.

운영 배포 /home/shin/deploy/knowslink, state /home/shin/deploy/knowslink-state, 도메인 https://link.knowslog.com. 마지막 배포 SHA 121cde923798c745af1f952347fcec1cae40b853. Tunnel config SHA256 f5abe860a019d5379cbbd5499eb309a645c4ae1e8f0b60fab6fea22326458512. 실제 착수 시 다시 읽어 현재값과 대조한다. 변경된 상태를 덮어쓰지 않는다.

SSH helper D:/workspace/KnowsLink/.git/fullops-gate/coor_remote.py와 기존 known_hosts RejectPolicy를 재사용한다. .env.server 값·code/token·메시지·DB 내용을 출력하거나 Git에 넣지 않는다. PowerShell literal stdin과 구조 인자로 전달한다. 원 helper·플러그인 캐시를 수정하지 않는다.

## 해야 할 일·소유권·완료 기준

1. 서버 clean·HEAD·서비스·env/credential/Tunnel hash·DB 집계 보존 기준을 읽기 전용으로 확보한다. 운영 설치1은 실제 사용자 승인 결과이므로 삭제하지 않는다. private 경로에는 원문, Git에는 비식별 요약만 보존한다.
2. 기존 state/backups에 고유 private DB backup을 생성하고 권한과 pg_restore --list를 확인한다. 이전 backup은 유지한다. 복구 실행은 하지 않는다.
3. 검수된 고정 SHA가 실제 origin/main 조상인지 확인한다. 서버 origin의 사라진 로컬 경로를 맹목적으로 fetch하지 않는다. 이전에 검증된 bundle/SFTP 방식으로 정확한 SHA를 전달하고 clean 배포에 안전하게 반영한다. migration이 있거나 서버 상태가 예상과 다르면 ask한다. 원 상태를 강제 덮어쓰지 않는다.
4. 기존 Compose project knowslink, compose.yaml+deploy/knowslink/compose.ops.yaml, state/.env와 UID/GID/STATE_DIR 계약을 유지한다. secret 출력 없는 config --quiet로 검증한다. relay만 최소 build/recreate한다. DB·Tunnel·다른 서비스·volume은 보존한다. 수정에 새 ingress 경로가 필요하지 않으면 config를 바꾸지 않는다. 제품 코드 수정이 필요하면 DEV에 반환한다.
5. local/public health, OAuth metadata, 무인증 MCP/SSE·polling 거부와 private 경로 404, 기존 공유 서비스의 이전 상태를 비교한다. 제품 전체 QA는 정확한 고정 SHA 증거를 재사용하고 반복하지 않는다. 실제 Google/Dots OAuth와 인증된 SSE/Events는 배포 후 부모 검증으로 남긴다. smoke 실패나 service degradation은 현재 상태·원인·담당·재개 조건을 보고한다. DB restore나 무조건 rollback은 하지 않는다.
6. 자신의 exec plan·QA 요약·contexts/ops와 D12 ops-guide/D13 transition만 실제 변경 영향에 맞게 갱신한다. 제품 source와 다른 역할 증거는 수정하지 않는다. 새 script/config가 없으면 원 exact-SHA product gate와 Windows 실행 불가를 정직하게 구분한다. 새 실행 파일/설정은 해당 registered gate와 고정 SHA 리뷰를 받는다. 예상 문서 변경은 5파일 이내다.
7. 모든 packet 비선택 path/category와 partial/unknown을 수동 근거로 outcomes 처리한다. work.py finish로 원 지시서와 완료 보고를 보존한다. worker_done은 최종 기록40SHA 한 개와 운영40SHA·검증 결과·미실행을 분리해 보고한다.

## 지시 전제와 충돌 — 먼저 확인

과거 ops-guide/transition의 held·Access·미적용 설명은 현재 OPS-LIVE와 구분한다. 최근 실제 배포는 위 121cde9이며 사용자 승인 이후 설치1/token0이 관찰됐다. 과거 수락을 새 OAuth 성공으로 간주하지 않는다. verify.py regression baseline 결함과 medium SSE 헤더 결함은 이미 보존된 잔여이며 이번 성공으로 바꾸지 않는다. 새 결함은 coor에 보고한다.

## 완료 보고

완료 전문과 최종 기록 SHA는 worker_done 전에 보존한다.