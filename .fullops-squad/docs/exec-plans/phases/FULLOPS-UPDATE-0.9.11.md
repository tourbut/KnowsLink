---
title: FullOps 0.9.11 업데이트 적용 기록
status: draft
updated: 2026-10-03
owner: coor
tasks: [FULLOPS-UPDATE-0.9.11]
summary: Codex 최신 설치와 글로벌 설치기 릴리스의 레포 적용 및 기존 작업 보존을 기록한다
---

# FULLOPS-UPDATE-0.9.11 — 레포 적용

## 기준과 설치

- 기준 ref: `815628f9b8ca942ca8bfe5f5f1b187a4e4acccac`.
- 현재 CLI는 Codex다. 업데이트 전 실제 설치 버전은 0.9.11이다. 레포 적용 버전은 0.9.10이다.
- `plugin.json`과 `codex plugin add`의 설치 경로로 실제 버전을 확인했다.
- `codex plugin marketplace upgrade fullops-squad`는 종료코드 0이며 이미 최신이다.
- `codex plugin add fullops-squad@fullops-squad`는 종료코드 0이다. 실제 설치 버전은 0.9.11이다.
- 설치 경로는 `/home/shin/.codex/plugins/cache/fullops-squad/fullops-squad/0.9.11`이다. 새 설치의 update-fullops를 다시 읽었다.
- `deps.py --check`는 종료코드 0이며 필수 CLI가 모두 있다. 추가 의존성 설치는 해당 없음이다.
- `update.py --repo . --from 0.9.11`은 종료코드 0이다. 레포 적용 기준 0.9.10 이후 릴리스 0.9.11을 확인했다.

## 릴리스 판정과 반영

0.9.11은 글로벌 npm 설치 편의 기능이다. 릴리스의 기존 레포 적용 절에 따라 역할·제품 코드·검증 기준·setup 변경은 해당 없음이다. 설치기를 별도로 전역 설치할 필요는 없다. 이전 0.9.10 기록의 필수 적용은 완료됐다.

`fullops.json`의 적용 버전을 0.9.11로 갱신하고 `PLANS.md`와 `board/board.json`에 같은 상태를 기록했다. 역할·모델·원격·기존 실패와 held·제품 중지 지시를 유지한다.

## 보존과 후속 담당

기존 미추적 SAR-MVP-001-INTEGRATION-review 로그와 JSON 증거는 사용자 작업으로 보존한다. 이번 업데이트 커밋에는 포함하지 않는다. 진행 중 worker를 중단하거나 신규 제품 worker를 시작하지 않았다.

필수 레포 적용의 보류는 없다. 역할 워크트리 동기화는 전달 보류다. 담당자는 coor이며 실제 유휴·깨끗한 상태를 확인한 뒤 재개한다. 이전 기록의 0.9.10 역할 동기화는 완료됐으므로 다시 적용하지 않는다. Claude Code·grok 및 다른 CLI 홈의 설치는 이번 Codex 업데이트 범위 밖이다.

다음 coordinator는 새 세션에서 기존 Run과 작업 현황을 확인한다. 기존 제품 지시서의 기준 ref와 완료 SHA를 변경하지 않는다.

## 검증

제품 코드 변경이 없으므로 제품 동작 테스트·빌드·시각 검수는 해당 없음이다. Git 공백 검사와 JSON 파싱·기존 역할 및 원격 보존을 확인한다. 운영 문서와 설정만 지정해 커밋한다.

## 검증 결과 복구 — 0.9.12 업데이트

준비 SHA `5b94be5e0cb1`의 깨끗한 detached snapshot에서 lockfile 기준 npm ci와 FullOps lint를 실행했다. 종료코드 0, product-lint passed, ERROR 0, WARNING 0, 실행 불가 0이었다. 이전 0.9.11 캐시 삭제로 후속 기록 명령이 hook에서 차단되어 이번 기록으로 복구한다.
