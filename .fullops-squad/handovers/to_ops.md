---
title: SAR-SERVICE-UI-001-OPS — 수락한 회원 웹·플러그인의 상태 보존 운영 배포
status: draft
updated: 2026-10-11
owner: ops
tasks: [SAR-SERVICE-UI-001-OPS]
summary: 수락한 회원 웹·플러그인의 상태 보존 운영 배포
attempt: f34e8cdec6d3408eacc74ae2b5170e08
base: 8539ef5b6d3e7ed99259762c62d164b8ba40fab0
subagent_level: standard
test_level: lite
---

# SAR-SERVICE-UI-001-OPS — 수락한 회원 웹·플러그인 보존 배포

- 상태: queued. UI 후보는 main 미수락이며 tester/visual 대기. coor가 수락 main 고정 SHA와 두 review 근거를 추가한 뒤 dispatch한다. 지금 배포하지 않는다.
- 역할: ops, C:/Users/shin/orca/workspaces/KnowsLink/ops, fullops/ops. 복귀 run_7176a2ad6ef1, coor term_e61d3e14-29e9-4954-943a-4a75707c82de. 실제 task/dispatch는 preamble.
- 공통 기준: fullops-common-0.3.3와 세 규칙, FULLOPS.md/project.md/document-writing. Ponytail full; lite/subagent standard. 등록 Claude Sonnet5.5 high fresh. Codex Trust 임의 승인 금지.
- 사용자 승인: 기존 .env.server의 Linux 서버에서 직접 시험·진행, link.knowslog.com 상용 UI 개선·기존 플러그인 구현의 운영 반영. 추가 유료서비스·로컬 Docker·파일/컨테이너 삭제·이동·광범위 정리 금지. F: 사고 증거 보존.

## 보존과 배포

운영 위치 /home/shin/deploy/knowslink, 상태 /home/shin/deploy/knowslink-state. 기존 .env.server 자격은 메모리에서 읽고 출력하지 않는다. 기존 coor_remote.py SSH helper와 OPS-LIVE/D12 runbook을 직접 읽어 재사용한다. 이전 운영 제품 eac603b6a41d2ce002e698600ac66b8fc67ad762를 복귀 기준으로 확인한다. 새 기준은 실제 작업 직전 수집한다. 이전 수치에 맞추려고 DB를 바꾸지 않는다.

배포 전후 relay DB·member/owner/agent/key/pair·기존 메시지/구독·Google/Tunnel/상태 자격과 config 해시·볼륨·공유 서비스 Orca200/s8200/mcp401/myportfolio·기존 포트·Tunnel ingress를 비교한다. 민감 데이터 원문·SQL dump는 보존/출력하지 않고 집계/비밀값을 제외한 digest만 기록한다. 새 integration 시험에서 생긴 변화는 원인과 구분한다. 기존 state·DB·keys·cert·Google/Tunnel credential·shared resources를 재생성/철회/삭제하지 않는다.

수락 SHA에 실제 Linux 배포 checkout이 clean detached인지 확인한다. 필요 artifact와 binary/container만 기존 runbook대로 갱신한다. 외부 의존성/유료 서비스를 추가하지 않는다. DB migration 변경0이면 기존 schema·state 유지하고 초기화하지 않는다. 실패하면 확인된 기존 배포로 복귀하고 사유와 증거를 보존한다. 새 설치 plugin zip이 webhook wake 코드와 정확한 배포 SHA에서 생성됐는지 package verify/hash로 확인한다. Bot runtime/env는 coor 후속이며 건드리지 않는다.

## 공개 검증과 완료

공개 / 및 /home·member 경로의 상태/CSP/로그인 경로, readiness·public deny·기존 회귀를 기존 도구로 확인한다. 새 Google 실제계정 로그인은 coor Orca browser 후속이며 자동 PASS로 기록하지 않는다. 웹 조작 필요 시 Orca 내장 브라우저만. Grok Bot UI 접근 금지. secret input/token이 든 페이지 캡처 금지.

본인 OPS-LIVE에 baseline/after/명령·종료코드·SHA·artifact hash·불변식·복귀 방법·미검증을 보존한다. 필요한 D12/D13 상태만 실제 결과로 갱신한다. phase/contexts/packet outcomes, 대화 미참조 인계 점검과 lint 근거를 기록한다. work.py finish로 전문 아카이브·인박스 비움·커밋 후 [완료] SAR-SERVICE-UI-001-OPS | 최종 SHA <실제 결과40자리>로 보고한다. 결과 SHA·운영 SHA를 구분한다. main 결과 통합은 coor 담당이다.

## 완료 보고
현재는 수락 고정 SHA와 검수 근거를 기다린다.

## 지시 전제와 충돌 — 먼저 확인
Jev context는 D12/D13/이전 OPS 로그를 conflict로 표시했다. 이 문서들을 먼저 읽는다. D12의 1–10장은 역사적 조사·계획이고 SSE 운영 후보 절의 미적용 문장은 상단 적용 결과가 이전 사실로 명시했다. D13/최근 OPS 로그의 실제 운영 제품은 eac603b다. 신규 UI는 독립 검수 후 수락 SHA를 추가할 때까지 미배포다. 이전 held/fail·미검증을 새 PASS로 바꾸지 않는다. 배포 직전 실제 서버가 정본과 다르면 coor에게 보고하고 무작정 checkout/up하지 않는다. 기존 승인된 공개 Google/member 서비스와 Tunnel allowlist를 유지하는 제품 갱신이며 Access/새인증 정책 변경은 없다.
