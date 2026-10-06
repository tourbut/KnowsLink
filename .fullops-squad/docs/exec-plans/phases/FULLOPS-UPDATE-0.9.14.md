---
title: FULLOPS-UPDATE-0.9.14
status: draft
updated: 2026-10-06
owner: coor
tasks: [FULLOPS-UPDATE-0.9.14]
summary: 설치와 레포 적용 차이 및 0.9.14 릴리스별 판정·검증
---

# FULLOPS-UPDATE-0.9.14

## 목표와 적용 기준

Codex 플러그인 설치와 현재 레포의 누락된 0.9.14 운영 변경을 적용한다. 업데이트 전 실제 설치는 0.9.14, 레포 적용은 0.9.13, 대상은 0.9.14다. 기준 ref는 `94533b207b456c0560800fe30a7c90b2b5887c6e`다. coor 후보 `c8856eb203b52e93b98ea7b65bf4195a4969f5cf`의 미수락 제품 조상을 회피해 main에서 준비한다. 역할·원격·인박스·제품 보류·기존 실패를 보존한다.

## 설치 확인

`codex plugin marketplace upgrade fullops-squad`는 이미 최신이며 종료코드 0이다. `codex plugin add fullops-squad@fullops-squad`는 실제 설치 0.9.14 경로를 반환하고 종료코드 0이다. 새 패키지 plugin.json과 update-fullops를 직접 읽었다. `deps.py --check`는 필수 CLI 모두 있음, 종료코드 0이다. `update.py --repo <coor> --from 0.9.14`는 더 오래된 레포 기준 0.9.13 이후의 릴리스 0.9.14를 출력했다.

## 릴리스별 판정

| 기존 레포 적용 항목 | 판정과 근거 |
|---|---|
| 글로벌 플러그인·기존 작업 보존 | 적용. 실제 Codex 0.9.14 확인. 새 제품 worker를 시작하지 않는다. |
| 기존 역할 setup·실제 검사 연결 | 적용. coor/designer/dev/ops/tester, origin/main으로 dry-run·실행 모두 종료코드 0, 생성 0개. 기존 product-lint의 make lint·cwd . 유지, product-test의 make test·kind test 추가. |
| 사용자 lint 규칙 보존·추가 | 적용. 기존 exclude·상한·timeout 보존. SIZE-002 400, SLOP-001–004와 DESIGN-001–003 WARNING 추가. 좁은 별도 예외 필요 없음. 설정은 merge-base에 따라 병합 후 적용한다. |
| Tailwind v4 shadcn | 해당 없음. Go template 일반 CSS와 TypeScript adapter다. 관련 프레임워크·Tailwind·shadcn은 없으므로 의존성을 추가하지 않는다. Go 문자열 CSS는 DESIGN 정규식 대상이 아니며 실제 UI 검수 기준은 유지한다. |
| 공통 규칙·프로젝트·핸드오버·리뷰·lint 안내 | 적용. fullops-common-0.3.3과 UI 정본·컴포넌트·예외·테마 검증 및 SIZE/DEP 보고 항목을 통합한다. 기존 책임 분리·완료 인박스 규약과 과거 기록을 보존한다. |
| 보드 문서 뷰어 | 적용. board.py로 설치 패키지의 고정 HTML과 데이터를 생성한다. title·단계·산출물·원천은 보존한다. 생성 board-data.js는 커밋하지 않는다. |

## 변경 규모와 완료 조건

운영 파일 약 12개, 검사 대상 추가 줄 약 350줄과 제외된 제공 보드 HTML을 예상한다. 의존성·제품 코드 변경은 없다. 변경 파일은 fullops.json, FULLOPS.md, PLANS.md, project.md, rules/common/README.md·coding-style.md, handovers/_TEMPLATE.md, review/_REPORT.md, lint/README.md·lint.json, board/board.json·index.html과 이 기록이다.

완료 조건은 설치 확인, 릴리스 적용, 문서 strict·공백·보드 데이터 검사, 고정 HEAD product-lint와 새 kind test 증거 통과다. 필수 검증 이후 적용 버전을 0.9.14로 갱신한다. 제품 UI 변경과 테마 전환은 해당 없음이다. 기존 서비스 디자인 lint는 미구성이므로 실행했다고 주장하지 않는다.

## 보류와 재개

필수 레포 적용과 역할 동기화를 구분한다. 역할 checkout에 진행 세션·변경·충돌이 있으면 coor가 최신 main SHA와 사유를 PLANS에 예약한다. 실제 idle·clean 확인 후 다음 dispatch 전에 동기화한다. 기존 invalid_lease 실패·제품 수락·사람 이메일·운영 공개·노우↔다닷 검증은 기존 담당과 재개 조건을 유지한다. 다음 제품 작업은 새 coordinator 세션에서 진행한다.
