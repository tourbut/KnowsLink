---
title: SAR-SSE-INSTALL-001-OPS 운영 이행 실행 기록
status: draft
updated: 2026-10-11
owner: ops
tasks: [SAR-SSE-INSTALL-001-OPS]
summary: 수락 제품 121cde9를 기존 Linux 서버에 bundle로 최소 적용한 계획·순서·결과와 부모 후속
---

# SAR-SSE-INSTALL-001-OPS

## 범위와 기준

기준 SHA는 지시서 front matter의 `121cde923798c745af1f952347fcec1cae40b853`이다. 규약은 `fullops-common-0.3.3`, `project.md`, `rules/delegation.md`, 문서 작성 규칙이다. 테스트 레벨 lite, 하위 위임 standard다. 서버 작업은 순서 의존이 있고 소유 파일이 겹쳐 하위 에이전트를 쓰지 않았다. 제품 코드·DB migration은 변경하지 않았다. 새 라이브러리·클라우드 CLI 구문은 사용하지 않아 Context7 조회가 필요하지 않았다. 기존 `beta.sh`·`cloudflared ingress` 명령과 D12 근거를 재사용했다.

과거 문서의 held·미구현 표시는 당시 시점의 기록이다. 현재 근거는 아래와 같이 구분한다.

| 구분 | 근거 |
|---|---|
| Google 운영 수락(2026-10-10) | `ba3754d`, [LIVE 기록](../../evaluations/qa-reports/SAR-GOOGLE-CONNECT-002-LIVE.md) |
| SSE·Remote 독립 리뷰·tester | `febd128`, [REVIEW](../../evaluations/qa-reports/SAR-SSE-INSTALL-001-REVIEW-review/report.md), [TESTER](../../evaluations/qa-reports/SAR-SSE-INSTALL-001-TESTER-test/report.md) |
| 동의 화면 시각 | [VISUAL](../../evaluations/qa-reports/SAR-SSE-INSTALL-001-VISUAL/report.md) |
| main 수락 | `121cde9` (origin/main 반영) |
| 운영 이행 | 이 과제. [운영 LIVE 기록](../../evaluations/qa-reports/SAR-SSE-INSTALL-001-OPS-LIVE.md) |

## 순서와 결과

1. 읽기 전용 preflight: HEAD `ba3754d`, 체크아웃 clean, 기존 origin 경로 부재, state 권한, config SHA256 `3b6f569e…8cd3`, compose config 통과, relay_state 집계 확보.
2. migration diff 없음. 제품 diff는 relay·adapters·package·문서 중심이고 compose·Dockerfile·`beta.sh`·`verify.py` 변경 없음.
3. bundle 전달과 검증, private pg_dump 백업(0600)과 `pg_restore --list` 통과.
4. detached checkout `121cde9`, `beta.sh up`으로 relay만 재생성. DB·cloudflared 컨테이너 유지.
5. `public-ingress.yml` 후보를 `check_public_ingress.py`, cloudflared validate·rule로 검사한 뒤 config를 갱신하고 cloudflared만 재시작.
6. local·public smoke와 보존 대조. 결과는 LIVE 기록을 따른다.

## 남은 일

1. 실제 Google 계정으로 호스트 설치와 인증 SSE·Remote OAuth·Dots Events를 확인한다. 담당은 부모와 사용자다.
2. `verify.py regression`의 baseline 비교 결함은 기존 상태다. 스크립트 수정은 이 과제 범위 밖이다. 후속 담당은 dev다.
3. 장애 시 복구 순서는 LIVE 기록의 「실패 시 안전 조치」를 따른다.
