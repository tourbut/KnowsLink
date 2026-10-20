---
title: FULLOPS-UPDATE-1.8.0
status: draft
updated: 2026-10-21
owner: coor
tasks: [FULLOPS-UPDATE-1.8.0]
summary: 1.7.1·1.8.0을 적용하고 핸드오버 템플릿·증거 커밋 범위 규칙·raw 제외를 레포에 옮긴다
---

# FULLOPS-UPDATE-1.8.0

## 버전

- 업데이트 전 설치: 플러그인 캐시 1.8.0(호스트 자동 업데이트, `installed_plugins.json` lastUpdated 2026-10-20), 고정 hook 런타임 1.7.1(`current.json` b1b289ce…). 런타임 기준으로 `--from 1.7.1`을 썼다.
- 레포 적용 버전: 1.7.0 → 1.8.0.
- 대상 버전: 1.8.0. `claude plugin marketplace update fullops-squad`, `claude plugin update fullops-squad@fullops-squad`가 최신 1.8.0을 확인했다.

## 설치

새 패키지 `deps.py --host claude-code` 종료0(ponytail 5.3.2, mattpocock-skills 1.3.1 최신), `--check --host claude-code` 종료0(호스트 의존성 검증 완료). 고정 hook 런타임 `current.json`이 1.8.0 패키지(61dfa11e…)로 전환됐다. `.fullops-squad/frontend.json`이 없어 선택형 프런트 도구는 해당 없음.

## 릴리스 판정

| 항목 | 판정·근거 |
|---|---|
|1.7.1 적용1 deps|적용. 위 설치 항목.|
|1.7.1 적용2 레포 파일|해당 없음. 같은 세션 reload 안내와 non-FF 병합 통합 lint는 스킬 절차라 업데이트된 `fullops-orca`·`update-fullops`를 읽으면 적용된다.|
|1.8.0 적용1 deps|적용. 위 설치 항목.|
|1.8.0 적용2 핸드오버 템플릿|적용. `handovers/_TEMPLATE.md`에 상태 구분 설명, ID 접두어 정본 경로, 하네스 제약(lint 경로·cwd·`--from`, 증거 형식, 재검토 base, 격리 환경 스크립트), 완료 보고 상세도·브랜치/SHA·후속 등재 안내를 옮겼다. 레포 고유 절(0.9.12 통합·고정 기준과 탐색 패킷)은 보존했다.|
|1.8.0 적용3 진행 중 패킷 재생성|보류(선택). 진행 중인 dev(SAR-SCHEDULE-001-DEV)·tester(SAR-HOST-VERIFY-001-LIVE) 인박스를 덮어쓰지 않는다. 담당 coor. 재개 조건: 다음 dispatch부터 새 `jev_packet.py`가 자동 적용되므로 별도 조치 없음.|
|1.8.0 적용4 증거 커밋 범위|적용. `rules/common/testing.md` `담당과 반복 범위`에 증거 커밋 범위 문단 추가, `rules/common/README.md` 규칙 묶음 `fullops-common-0.3.4`, `.fullops-squad/.gitignore`에 `docs/evaluations/qa-reports/**/raw/` 추가. 이미 커밋된 qa-reports 증거는 이력 보존을 위해 지우지 않는다. 용량 정리는 필요 시 별도 과제.|
|setup dry-run|기존 역할·origin·main·coor·lite·standard로 생성 예정 0개.|

## 검증

- `update.py --repo . --from 1.7.1`이 릴리스 1.7.1·1.8.0을 출력했다.
- `git ls-files -ci --exclude-standard` 결과 없음: 새 raw/ 제외 규칙이 기존 추적 파일과 겹치지 않는다.
- FullOps lint를 변경 범위에 실행했다(커밋 기록 참조). 제품 코드 변경 없음.

## 보류와 동기화

- 진행 중 패킷 재생성: 위 표.
- 역할 워크트리 동기화: 진행 중인 dev·tester는 과제 완료 뒤 동기화를 예약한다. 깨끗한 유휴 역할만 최신 main으로 동기화한다.
- 같은 세션 계속: Claude Code에서 사용자가 `/reload-plugins`를 실행하면 같은 coordinator 세션에서 이어간다.
