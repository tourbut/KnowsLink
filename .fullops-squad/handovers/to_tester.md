---
title: SAR-MVP-002-INSTALL-FIX-DEV-TESTER — 설치 수정 후보 5506d64의 실제 Grok CLI 독립 설치 QA
status: draft
updated: 2026-10-04
owner: tester
tasks: [SAR-MVP-002-INSTALL-FIX-DEV-TESTER]
summary: 설치 수정 후보 5506d64의 실제 Grok CLI 독립 설치 QA
---

# SAR-MVP-002-INSTALL-FIX-DEV-TESTER — 설치 수정 고정 SHA 실제 Grok CLI 독립 설치 QA

- 작성일: 2026-10-04; From: coor; 상태: ready
- 복귀: /home/shin/orca/workspaces/KnowsLink/fullops-coor, term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9, run_8ca8bc058ab7. task/dispatch는 preamble.
- 병합 담당 coor, main/origin. 사용자 승인: 이슈1 진단·수정·검증·푸시·완료 후 댓글 게시. worker는 외부 댓글을 쓰지 않는다.
- 실제 Bot 계정/유료 inference/relay/DEC-02/calendar·실데이터는 held, FullOps 업데이트 제외.

## 적용 기준과 먼저 읽기

FULLOPS.md, project.md, rules/common/README.md 및 coding-style.md/testing.md/security.md, docs/agents/document-writing.md, docs/exec-plans/phases/SAR-MVP-002-INSTALL-FIX-DEV.md, 해당 DEV handovers/logs/2026-10-04_to_dev.md, adapters/README.md, scripts/package_plugin.py, scripts/verify_grok_plugin.py. 공통 0.3.2. 현재 준비 SHA에서 규칙을 읽고 제품 대상은 아래 고정 SHA를 쓴다. 이전 QA0fb32cd·리뷰0055a5b는 코어/Go/UI 변경 없는 부분의 증거로 재사용 가능하다. 설치 및 패키지 변경은 새로 확인한다.

- 역할/워크트리: tester, /home/shin/orca/workspaces/KnowsLink/fullops-tester, fullops/tester
- 고정 후보: 5506d646c47c2d64b35d1254ddfdec5e2003084d, 기준 376981441d1b0650897847273e1773f8473949e5.

Jev find/documents-find/context: docs/evaluations/jev/SAR-MVP-002-INSTALL-FIX-DEV-TESTER-*.json. keep 전부 읽는다. 지시 전제와 충돌 — 먼저 확인: 현재 DEV 실행 기록과 이전 QA(설치 미실행); current DEV 기록의 주장 자체를 독립 검증하며 이전 QA는 불변 코어에만 재사용한다. 근거가 틀리면 결함으로 보고하고 제품 규칙 충돌이면 ask한다.

## 해야 할 일과 소유권

- [ ] fullops-test를 적용하고 /tmp 아래 새 독립 clone의 detached5506d64로 검증한다. grok --version과 help 계약을 먼저 확인한다.
- [ ] 독립 clone에서 Node22.22.2/npm10.9.7로 make install·make plugin·make verify-grok-plugin을 각각 실행하고 자체 종료코드·ZIP hash를 기록한다. 예상 hash fb745c66b4e786f5267228099c3794763632381451bd37cf66a82111f9b14a75, 다르면 원인 판정한다.
- [ ] 임시 HOME/cwd에서 문서의 실제 validate·install --trust·plugin list·mcp doctor를 직접 재현한다. 서버1개·tools2개와 설치본 knowslink_status held, stdout 순수프로토콜을 확인한다. make verify-grok-plugin만 실행한 결과로 독립 의미 검증을 대체하지 않는다.
- [ ] 실제 Grok CLI 설치 경로, 설치 후 원본 폴더 제거 가능, marketplace 경로, 기존 HOME/사용자 ~/.grok 무변경을 확인한다. 새 CLI 기능의 trust 누락 시 활성화 결과도 확인한다. 사용자 Grok 설정과 계정은 수정하지 않는다.
- [ ] Node20.19.2 npm ci 실패1(EBADENGINE)와 Node22 성공0을 검증한다. README 공식 Node 준비/체크섬 명령과 필요한 Python/Make/Grok 조건이 재현 가능한지 확인한다. Node 다운로드는 공식 checksum 확인하고 시험 영역을 쓴다.
- [ ] 변경 없는 Go/UI/relay 증거는 기존 QA를 재사용한다. 새 설치 흐름의 수락 결과·명령·exit·재현·불가 항목을 docs/evaluations/qa-reports/SAR-MVP-002-INSTALL-FIX-DEV-TESTER.md와 -test/에 기록한다. Bot 계정 설치/동적 카탈로그/새 Bot 세션은 원격 재시험 전까지 미검증이다.

## 완료 기준과 결과

정상/실패/격리 증거가 있어야 한다. 유료 inference·실제 relay는 실행하지 않는다. 제품 코드/타 역할 인박스·준비 리뷰 파일/PLANS·board를 수정하지 않는다. 테스트·QA 증거·완료 로그만 소유한다. 결함은 수정하지 않고 재현과 판정으로 보고한다. strict metadata 문제를 발견하면 자기 문서만 수정하고 다른 역할은 coor에 전달한다. D01–D13 갱신 없음. work.py finish로 지시서/완료 전문 보존·인박스 비우기 후 커밋, lint --from착수 HEAD exit0을 확인하고 worker_done으로 SHA/결함/미검증을 보고한다.
