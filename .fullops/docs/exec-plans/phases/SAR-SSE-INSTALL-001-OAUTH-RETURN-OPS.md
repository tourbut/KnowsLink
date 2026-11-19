---
title: SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS 운영 배포 실행 기록
status: draft
updated: 2026-10-11
owner: ops
tasks: [SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS]
summary: 수락된 OAuth 복귀 수정 eac603b를 기존 서버에 relay만 갱신해 적용한 계획·순서·결과와 부모 후속
---

# SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS

## 범위와 기준

배포 제품 SHA는 `eac603b6a41d2ce002e698600ac66b8fc67ad762`다. 지시서 attempt는 `f5a42297abfc45748fe3746cb2fc975e`다. 규약은 `fullops-common-0.3.3`, `project.md`, `rules/delegation.md`, 문서 작성 규칙이다. 테스트 레벨 lite, 하위 위임 standard다. 서버 작업은 순서 의존이 있어 하위 에이전트를 쓰지 않았다. 제품 코드·DB migration·compose·ingress는 변경하지 않았다. 새 라이브러리·클라우드 CLI 구문을 쓰지 않아 Context7 조회가 필요하지 않았다. 기존 `beta.sh`와 이전 OPS 기록의 bundle 방식을 재사용했다.

| 구분 | 근거 |
|---|---|
| 이전 운영 | `121cde9`, [이전 LIVE 기록](../../evaluations/qa-reports/SAR-SSE-INSTALL-001-OPS-LIVE.md) |
| 독립 QA·리뷰 | [TESTER 리뷰](../../evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-review/report.md), [TESTER 테스트](../../evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-test/report.md) |
| 시각 수락 | [VISUAL](../../evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL/report.md) |
| 이번 운영 적용 | [운영 LIVE 기록](../../evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS-LIVE.md) |

## 순서와 결과

1. 읽기 전용 baseline: HEAD `121cde9`, 체크아웃 clean, state 권한, Tunnel config SHA256 `f5abe860…8512`(기준과 같음), relay_state 집계, 컨테이너·공유 서비스 상태.
2. `121cde9..eac603b`의 제품 diff가 `remote_auth.go`와 테스트·드라이버뿐이고 migration·compose·deploy diff가 0건임을 확인했다. ingress 변경은 필요하지 않다.
3. private pg_dump 백업(0600)과 `pg_restore --list` 통과. 복구는 하지 않았다.
4. bundle 전달(SHA256 일치), 조상 관계 확인, `deploy-history.log` 기록, detached checkout `eac603b`, `beta.sh up`으로 relay만 재생성.
5. 보존 대조(해시·집계·컨테이너·공유 서비스)와 local·public smoke. 합성 43자 consent 경로가 새 HTML 401과 보안 헤더를 반환했다. 결과는 LIVE 기록을 따른다.

## 남은 일

1. 부모가 실제 Google 계정으로 원격 연결을 재시도해 승인 뒤 복귀와 token 발급, 인증 SSE·Events를 확인한다.
2. F1 medium(원인별 실패 안내), SSE 401 보안 헤더 medium, `verify.py regression` baseline 결함은 기존 잔여다. 담당은 dev다.
3. 미완료 install1은 삭제하지 않았다. 처리 판단은 사용자·부모 몫이다.
4. 장애 시 복구 순서는 LIVE 기록의 「실패 시 안전 조치」를 따른다.
