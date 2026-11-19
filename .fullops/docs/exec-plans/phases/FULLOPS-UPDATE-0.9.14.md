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

## 검증 결과

- 초기 main의 make lint와 make test는 종료코드 2다. adapter node_modules의 MCP SDK·zod 누락으로 TypeScript TS2307이 발생했다. 제품 코드를 수정하지 않았다. 기존 lockfile의 `npm ci --prefix adapters`로 188개 패키지를 설치했으며 종료코드 0이다. lockfile·의존성 선언 변경은 없다.
- 준비 SHA `cf55ea1d05a5`의 실제 전체 값은 아래 보존 JSON의 head에서 확인한다. `lint.py --repo <main> --from HEAD`는 새 설정 기준으로 product-lint와 kind test의 product-test를 모두 실행했다. 종료코드 0, ERROR 0, WARNING 0, 실행 불가 0이다. 파일 차이 검사는 최종 SHA에서 최초 기준 94533b2로 별도 수행한다.
- `deliverables.py --strict`는 13개, 문제 0, 경고 0이다. git diff --check는 종료코드 0이다. board.py 종료코드 0, HTML은 설치 패키지 제공본과 byte 동일하다. 데이터 JSON·13개 산출물과 생성 JS·inline JS 문법을 확인했다. 실제 브라우저 조작은 이번 레포에서 재실행하지 않았다. 제공 패키지의 뷰어 검증과 구분한다.
- 고정 HEAD 검사 JSON은 docs/evaluations/FULLOPS-UPDATE-0.9.14/preflight-lint.json에 보존한다. command의 kind·exit_code·원출력으로 확인한다. 파이프라인으로 명령 종료코드를 가리지 않았다.
- 운영 파일만 변경했으며 DEP-001 대상 의존성 변경은 없다. 코드 설계 판단·제품 UI·공개 정책 변경은 없다. 필수 적용·관련 검사 통과 후 plugin_version을 0.9.14로 갱신했다. 기존 제품 후보와 실패 수락은 그대로다.

## main 통합과 역할 동기화 확인

운영 변경은 main `08cf165e7e2e3208cd14b9ea1bff9da134d1fd51`에 커밋하고 origin/main에 일반 push했다. fetch와 ls-remote가 같은 SHA를 반환하며 준비 cf55ea1의 조상 관계를 확인했다. 최초 기준 94533b2의 최종 lint는 product-lint 통과, ERROR 0, WARNING 2, 실행 불가 0이다. 경고는 LINT-001 설정 변경 안내와 누적 PLANS SIZE-001이다. 설정 도입 이후 기준 cf55ea1의 같은 HEAD lint는 product-lint·product-test 모두 통과, ERROR 0, WARNING 1, 실행 불가 0이다. 경고는 누적 PLANS 길이다. 두 결과는 같은 평가 폴더의 final-diff-lint.json·final-command-lint.json에 보존한다.

Orca 역할 terminal 조회에서 designer/dev/ops/tester는 0개이며 Git clean이었다. coor의 다른 터미널은 셸 프롬프트 상태이며 현재 업데이트 세션만 작업 중이었다. 다섯 역할에 main을 반영했다. coor/dev/tester의 PLANS·project 충돌은 각 역할의 기존 제품 기록·과제 키와 새 운영 절을 모두 보존해 해결했다. 제품 파일은 main으로 가져오지 않았다. 역할의 동기화 결과는 다음과 같다.

- coor `d86712ae8c15812972af40a15e39ef6706af8142`
- designer `08cf165e7e2e3208cd14b9ea1bff9da134d1fd51`
- dev `b80ba9857296038e23251886fceef9297665dbcc`
- ops `08cf165e7e2e3208cd14b9ea1bff9da134d1fd51`
- tester `e12ef97fc5e74c2788d32854610170d4b4e0bc14`

이 기록의 후속 main 커밋도 모든 역할에 반영한다. 진행 중 worker 동기화 예약은 없다. 제품 후보의 기존 실패·검수 보류는 유지하며 해당 후보를 이번 업데이트로 수락하지 않는다. 다음 개발은 새 coordinator 세션에서 진행한다.
