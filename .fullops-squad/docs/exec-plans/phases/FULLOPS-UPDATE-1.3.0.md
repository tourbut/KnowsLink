---
title: FULLOPS-UPDATE-1.3.0
status: draft
updated: 2026-10-09
owner: coor
tasks: [FULLOPS-UPDATE-1.3.0]
summary: FullOps 1.3.0 적용과 환경 링크 및 tester 후보 변경의 검증 기록
---

# FULLOPS-UPDATE-1.3.0

## 목표와 적용 기준

업데이트 전후 실제 Codex 설치 버전은 1.3.0이다. 레포 적용 버전은 1.2.0이며 대상은 1.3.0이다. 기준 ref는 `f970905e428103dfc5dfef856aa0d00fc516696e`다. 공통 기준 `fullops-common-0.3.3`, FULLOPS.md, project.md, 문서 작성 규칙을 적용한다. 제품 코드와 의존성, 기존 인박스, 제품 중지와 QA 보류는 보존한다.

## 설치와 릴리스별 판정

- marketplace upgrade는 already up to date, plugin add는 설치 경로의 1.3.0을 반환했다. 두 명령과 deps.py --check --host codex는 종료코드 0이다. 출처 전환과 캐시 삭제는 없다.
- update.py --from 1.3.0은 레포 기준 1.2.0 이후의 1.2.1과 1.3.0을 반환했다. 새 설치 경로의 update-fullops를 다시 읽었다.
- 1.2.1: 해당 없음. 선택형 이슈 모드는 활성화 요청이 없으며 폴러나 lease를 시작하지 않는다. 새 세션에서 사용한다.
- 1.3.0: 적용. 기존 역할의 setup dry-run과 실행으로 delegation.md와 관리 블록을 추가했다. 테스트는 lite, 선택형 위임은 off를 유지했다. minimal/exhaustive 설명을 추가했고 기존 필수 명령의 level과 required는 변경하지 않았다.
- 새 delegation.md에 setup이 중복 생성한 front matter를 한 개로 정리했다. 기존 문서 본문은 보존했다.

## 사용자 추가 요청

main의 `.fullops-squad/.env`를 coor/designer/dev/ops/tester에 심볼릭 링크했다. env_link.py --all의 다섯 결과는 모두 link다. 기존 추적 `.env.example`은 보존했다. 비밀값은 읽거나 출력하지 않았다. 링크 대상과 Git ignored 상태를 검사한다.

tester에 기존 레포의 `claude-sonnet-5-5` 식별자로 medium/high 후보를 추가했다. 기존 Grok 4.7 high 후보와 다른 역할 모델은 보존했다. 실제 Claude 기동은 이번 운영 업데이트에서 수행하지 않는다.

## 검증과 통합 조건

설정 JSON·관리 블록·문서 strict·공백과 비밀 파일 ignored 상태를 검사한다. 제품 코드 변경은 없으므로 제품 동작 테스트는 해당 없음이다. 등록 lint 실행 결과는 성공·실행 불가를 구분해 기록한다. Windows에 make가 없으면 필수 제품 검사를 통과로 표시하지 않는다. 문서 운영 변경 검증과 독립 리뷰를 완료한 뒤 적용 버전을 갱신한다.

기존 coor와 dev에는 중지된 WIP가 있다. 업데이트의 작업 diff만 깨끗한 main에 적용하여 별도 커밋한다. 기존 coor와 dev WIP 커밋은 main에 병합하지 않는다. 검토한 main 커밋을 원격에 push하고 역할 브랜치에 병합한다. 유휴 clean 역할만 동기화하며 진행본은 보존한다.

문서 strict 검사 13개는 문제 0·경고 0이다. 초기 stamp는 신규 문서 summary 누락으로 한 번 실패했으며 summary를 지정해 복구했다. 작업 트리에서 lint 실행은 clean 조건으로 거절됐으며 커밋 후 다시 실행한다.
