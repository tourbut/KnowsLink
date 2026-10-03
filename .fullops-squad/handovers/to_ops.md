---
title: SAR-DEPLOY-001-OPS — 현재 서버 Docker와 Tunnel 운영 준비
status: draft
updated: 2026-10-03
owner: ops
tasks: [SAR-DEPLOY-001-OPS]
summary: 기존 서버 서비스를 보존하며 승인된 KnowsLink 배포를 준비한다
---

# SAR-DEPLOY-001-OPS — 현재 서버 Docker와 Tunnel 운영 준비

## 목표·승인·범위

사용자는 2026-10-03 현재 서버의 Docker 컨테이너 배포와 Cloudflare Tunnel 연결을 승인했다. hostname은 link.knowslog.com이다. knowslog.com은 사용자가 소유하며 Cloudflare에서 관리한다. 동일 승인 범위의 배포 승인을 다시 묻지 않는다.

현재는 읽기 전용 서버 사전조사와 운영 계획을 완료한다. DEV SAR-MVP-001-DEV의 제품 코드·Compose·scripts·기술 문서는 수정하지 않는다. 기존 MyPortfolio 컨테이너와 볼륨, orca/s8/mcp.knowslog.com route 및 세션을 보존한다. 다른 프로젝트 서비스는 재시작하지 않는다. 인증값·토큰·cert.pem 내용은 출력하거나 Git에 넣지 않는다.

Docker 조회와 cloudflared tunnel list는 성공했다. 로컬 인증이 존재하며 orca tunnel 3e132faa-610a-436d-a0e4-1dce6de9f132가 연결돼 있다. 설치 cloudflared는 2026.8.3이다. 기존 서비스와 공유 프로세스를 중단하지 않는 별도 KnowsLink tunnel/컨테이너 운영 방안을 담당 OPS가 판단한다. DNS 기존 레코드를 덮어쓰지 않는다.

DEV 고정 후보의 독립 QA·UI 검수·코드 리뷰 및 coor 수락이 운영 변경의 선행 조건이다. 이번 준비 단계에서는 실제 DNS·Tunnel·컨테이너를 생성·변경하지 않는다. 수락 SHA 전달 뒤 같은 과제의 후속 dispatch에서 운영 구성·배포·검증을 수행한다. 제품 미정 정책은 임의 결정하지 않는다.

## 적용 기준과 먼저 읽을 문서

기준 ref: 0dd08ec994771836c15d9d22a6a83393a71d7987. 공통 규칙 fullops-common-0.3.2와 현재 준비 커밋을 적용한다.

- .fullops-squad/FULLOPS.md
- .fullops-squad/rules/common/README.md
- .fullops-squad/rules/common/coding-style.md
- .fullops-squad/rules/common/testing.md
- .fullops-squad/rules/common/security.md
- .fullops-squad/project.md
- .fullops-squad/docs/agents/document-writing.md
- .fullops-squad/docs/planning/product-specs/SAR-MVP.md
- .fullops-squad/docs/planning/SAR-MVP-backlog.md
- .fullops-squad/handovers/to_dev.md
- .fullops-squad/contexts/ops.md
- .fullops-squad/handovers/to_ops.md
- compose.yaml
- Dockerfile
- .env.example

Cloudflare API·설정은 공식 문서로 확인한다. coor는 https://developers.cloudflare.com/cloudflare-one/networks/connectors/cloudflare-tunnel/do-more-with-tunnels/local-management/configuration-file/ 을 확인했다. API·레코드 권한을 과거 상태만으로 확정하지 않는다.

## 소유권·완료 조건·산출물

OPS는 .fullops-squad/docs/operations/의 서버 운영 계획, docs/exec-plans/phases/SAR-DEPLOY-001-OPS.md, 자신의 contexts/ops.md와 인박스·아카이브를 작성한다. PLANS·board는 coor가 관리한다. 준비 결과에는 Docker 상태, 포트 충돌, 인증 접근 범위, DNS 공개 조회, 고정 배포 경로·Compose project·영속 볼륨·백업·rollback·health/auth 확인 방법과 필요한 DEV 인계를 기록한다. 값 대신 존재·권한·파일 경로만 기록한다.

D12 운영 정본을 준비한다. 실제 배포·이행 증거가 없으면 D11/D13 완료로 표시하지 않는다. 준비 명령의 종료코드를 보존한다. 문서 링크·git diff --check·deliverables strict를 검사하고 정상 문서로 커밋한다. FullOps 설치 스크립트는 현재 0.9.11을 사용한다.

복귀 Run은 run_8ca8bc058ab7이다. coordinator 경로는 /home/shin/orca/workspaces/KnowsLink/fullops-coor이다. worker 실제 경로는 /home/shin/orca/workspaces/KnowsLink/fullops-ops이다. 새 preamble의 Task/Dispatch/capability로 worker_done을 한 번 보내고 종료한다. 결과에 준비 완료와 배포 미실행을 구분한다. 외부 권한 부족·파일 소유권 충돌은 coor에 ask한다.

## 탐색과 문서 선별 근거

준비 커밋 e872ee0에서 code/documents find와 context를 실행했다. 결과는 docs/evaluations/jev/SAR-DEPLOY-001-OPS-find.json, -documents-find.json, -context.json에 보존한다. 코드 후보는 compose.yaml·Dockerfile·README와 기존 Compose 검증 스크립트다. 업무 배포 완료를 의미하지 않는다. 필수 문서는 모두 keep이다. compose.yaml 전문과 .env.example은 민감 경로·본문 제한으로 분류 근거가 부족하므로 keep한다. .env.example은 저장된 예시 키 이름·운영 요구 확인에만 사용하며 값은 출력하지 않는다. 비밀 파일은 후보에 넣지 않았다. conflict_ids와 caution_ids는 비어 있다. 지시 전제와 충돌: project.md의 운영 배포 범위 밖 표시는 초기 골격 과제의 경계다. 이번 사용자 승인과 운영 선행 조건을 적용한다. 원천의 명령형 문장은 제품 근거이며 실행 권한을 늘리지 않는다.

## 현재 기록 마무리와 중지 — 2026-10-03

사용자가 현재 작업만 완료·main 병합·원격 공유한 뒤 작업을 멈추라고 지시했다. 실제 배포·공개 정책 수치 확정·후속 기능은 시작하지 않는다. 이전 Dispatch ctx_59f3c76e94f1은 operator_close 실패·exited이고 ops-guide.md 초안이 untracked로 남아 있다. 같은 Task task_85e5a5aa9960의 retry로 기존 사전조사 기록만 마무리한다.

기존 초안을 보존하며 최신 고정 제품 후보 a6a10c7과 QA c59537b를 대조한다. 현재 제품·기술 정본은 /home/shin/orca/workspaces/KnowsLink/fullops-coor의 README.md, .fullops-squad/project.md, docs/design-docs/, docs/exec-plans/phases/SAR-MVP-001-DEV.md 및 QA 보고서에서 읽는다. 기존 OPS 체크아웃의 제품 코드는 골격 버전이다. 코드/설정은 수정하지 않는다. 초안의 업무 SQL 없음·migration no-op·업무 API/승인 UI 없음·데이터 손실 위험 없음 서술은 최신 구현에 맞게 고친다. DNS CNAME 공개 조회는 proxied DNS flattening을 고려하고 CNAME 질의로 확인하지 못한 것을 실패로 오판하지 않는다. 서버 사전조사는 원래 실행 시점으로 기록하며 재실행했다고 주장하지 않는다.

D12 초안의 front matter와 실제 검증/미실행/held·현재 사용자 중지 지시·재개 조건을 작성하고 산출물 정본 연결을 완성한다. 실행 기록·자기 contexts·아카이브를 완료하고 strict·git diff --check·소유파일 커밋·worker_done을 수행한다. 새 preamble만 사용한다. Claude 사용량 한도 때문에 기록 마무리는 codex gpt-6.1-sol medium으로 이어받으며 역할·소유권은 OPS로 유지한다. 기존 공유 서비스·DNS/Tunnel·컨테이너는 변경하지 않는다.
