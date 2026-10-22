---
title: FULLOPS-UPDATE-1.8.3
status: draft
updated: 2026-10-22
owner: coor
tasks: [FULLOPS-UPDATE-1.8.3]
summary: "1.8.3을 적용해 공통 규칙 0.3.5, 지시서 상위 의도·증명 방법 칸, intents 템플릿을 레포에 옮기고 빈 cps 폴더를 지운다"
---

# FULLOPS-UPDATE-1.8.3

## 버전

- 업데이트 전 설치: 1.8.3(호스트 자동 업데이트, `installed_plugins.json` lastUpdated 2026-10-22T14:19Z).
- 레포 적용 버전: 1.8.2 → 1.8.3.
- 대상 버전: 1.8.3. `claude plugin marketplace update` 뒤 `claude plugin update`가 "already at the latest version (1.8.3)"을 확인했다.

## 설치

`deps.py --host claude-code`로 의존성을 갱신했고 `deps.py --check --host claude-code`가 `호스트 의존성 검증 완료`를 출력했다(ponytail 5.3.2, mattpocock-skills 1.3.1 최신). `.fullops-squad/frontend.json`이 없어 선택형 프런트 도구는 해당 없음.

## 릴리스 판정

| 항목 | 판정·근거 |
|---|---|
|필수: 공통 테스트 규칙 0.3.5|적용. `rules/common/testing.md` 본문을 새 패키지 자산으로 교체하고 레포 front matter를 유지했다. 레포 본문은 1.8.2 자산과 같아 통합할 로컬 수정이 없었다. `rules/common/README.md`도 교체했다(0.3.4 → 0.3.5). coding-style·security는 새 자산과 같다.|
|필수: 지시서 상위 의도·증명 방법|적용. `handovers/_TEMPLATE.md`에 `상위 의도` 줄과 `증명 방법` 문장을 넣고 `완료 기준과 검증`·`UI 작업의 디자인 기준과 검증`·보고서 경로 문장을 새 판으로 바꿨다. 레포 고유 절(통합·탐색 패킷)은 보존했다.|
|필수: intents 템플릿|적용. setup을 기존 역할로 dry-run(생성 예정 1개) 뒤 실행해 `docs/planning/intents/_TEMPLATE.md`를 추가했다.|
|필수: cps 폴더 정리|적용. `docs/planning/cps/`에 `.gitkeep`만 있음을 확인하고 지웠다.|
|project.md 규칙 묶음|적용. `fullops-common-0.3.3` → `fullops-common-0.3.5`.|
|현황판|할 일 없음. `board/index.html`은 board.py가 갱신한다.|
|hook(링크 안 된 증거·차수 폴더 차단)|레포 작업 없음. 이후 과제부터 적용된다. 진행 중 인박스는 없다.|

## 검증

- `update.py --repo . --from 1.8.3 --verify`: 필수 적용 5/5 충족, exit 0.
- FullOps lint를 이 커밋 범위에 실행했다(커밋 기록 참조). 제품 코드 변경 없음.

## 보류와 동기화

- 보류 없음.
- 활성 worker가 없으므로 모든 역할 워크트리를 최신 main으로 동기화한다.
- Claude Code 세션은 `/reload-plugins` 뒤 같은 세션에서 이어간다.
