---
title: SAR-SERVICE-UI-001-OPS 운영 배포 실행 기록
status: draft
updated: 2026-10-11
owner: ops
tasks: [SAR-SERVICE-UI-001-OPS]
summary: 수락한 회원 웹 UI·Grok wake 플러그인 95df83c를 기존 서버에 금지 entrypoint 없이 relay만 갱신해 적용한 계획·순서·결과와 coor 후속
---

# SAR-SERVICE-UI-001-OPS

## 범위와 기준

배포 제품 SHA는 `95df83c22c1559453a7c73504aa248088a8e4ed0`, 준비 main은 `0045089f2a2c3efe3e210eedc0c5cd4868487e42`다. 지시서 attempt는 `f34e8cdec6d3408eacc74ae2b5170e08`이다. 규약은 `fullops-common-0.3.3`, `project.md`, 문서 작성 규칙이다. 테스트 레벨 lite, 하위 위임 standard다. 서버 작업은 순서 의존이 있어 하위 에이전트를 쓰지 않았다. 제품 코드·migration·compose·ingress는 바꾸지 않았다. 증거 정본은 [운영 LIVE 기록](../../evaluations/qa-reports/SAR-SERVICE-UI-001-OPS-LIVE.md)이다.

| 구분 | 근거 |
|---|---|
| 이전 운영 | `eac603b`, [OAuth 복귀 LIVE 기록](../../evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-OPS-LIVE.md) |
| 독립 QA·리뷰 | [FIX-TESTER 리뷰](../../evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-TESTER-review/report.md) |
| 시각 수락 | [FIX-VISUAL](../../evaluations/qa-reports/SAR-SERVICE-UI-001-FIX-VISUAL/report.md) |

## 지시서 전제와의 충돌 처리

지시서가 `beta.sh`의 `prepare`·`deploy`·`backup`·`restore-verify`·`seed`·`owner-login`과 삭제·정리를 금지했다. `deploy`는 backup(raw `pg_dump`, 실패 시 파일 삭제)을 호출하므로 쓰지 않았다. 이전 과제의 "backup 후 배포" 순서에서 backup을 뺀 같은 순서(migration diff 검사→bundle 전달→detached checkout→`beta.sh up`→`verify.py local`)를 수동 실행했다. 따라서 이번에는 새 DB backup이 없다. 복귀는 코드 전용이다.

## 순서와 결과

1. 읽기 전용 baseline(서버 HEAD `eac603b` clean, state 권한·해시, relay_state 키별 개수·digest, 컨테이너·포트·cloudflared·공유 서비스).
2. `eac603b..95df83c` 제품 diff 15개 파일. migration·compose·deploy·Dockerfile·go.mod 0건. `95df83c`와 준비 main `0045089`는 제품 동일.
3. bundle(SHA256 `132262f1…2b59`) 전달·검증·fetch, 조상 관계 확인, detached checkout, `deploy-history.log`, `beta.sh up`, `verify.py local`. 종료코드 전부 0.
4. 전·후 비교: 해시·집계 digest 동일. `epoch`와 `Rates`만 relay 동작으로 변함.
5. `git archive` 사본에서 package 4종 생성·`--verify`(종료 0). Grok zip에 wake 코드 확인.
6. 무인증 공개 smoke. `/`에서 새 디자인 토큰과 건너뛰기 링크를 확인해 새 코드가 실행 중임을 직접 확인.

## 남은 일

1. coor가 Orca 내장 브라우저로 실제 Google 로그인·`/home` 회원 목록·모바일 렌더·키보드 포커스를 확인한다. 정확한 운영 SHA는 `95df83c22c1559453a7c73504aa248088a8e4ed0`, URL은 `https://link.knowslog.com/`다.
2. 새 설치용 Grok plugin zip(서버 `knowslink-state/serviceui-20261011/pkg/build/`, SHA256 `762411df…cf94`)의 전달·설치는 coor가 정한다. Bot runtime·watcher는 건드리지 않았다.
3. `verify.py regression` baseline 결함(`projects`)과 SSE 401 보안 헤더 medium은 기존 잔여다. 담당은 dev다.
4. 새 DB backup이 필요하면 coor가 허용한 별도 지시가 필요하다.
