# KnowsLink 프로젝트 기준

| 항목 | 값 |
|---|---|
| 제품 목적 | 미정 — 첫 요청에서 designer가 요구사항으로 기록 |
| 기본 브랜치·원격 | `main`, `origin` (`https://github.com/tourbut/KnowsLink.git`) |
| 현재 상태 | README만 있는 초기 레포. 제품 코드·에셋·테스트·배포 설정 없음 |
| 기술 스택 | 미정 — dev가 첫 개발 과제에서 기술 설계와 함께 선택 |
| 기술 설계 정본 | `.fullops-squad/docs/design-docs/` — dev 담당, 실제 문서는 과제 때 작성 |
| 기획 정본 | `.fullops-squad/docs/planning/`, 사용자 경험은 `docs/design-docs/mockups/` — designer 담당 |
| 공통 개발 기준 | [rules/common/README.md](rules/common/README.md), `fullops-common-0.3.1`; Ponytail full |
| 보안·코딩 규칙 | [코딩](rules/common/coding-style.md), [테스트](rules/common/testing.md), [보안](rules/common/security.md) |
| 문서 언어 | 한국어 |
| 이슈 트래커 | 로컬 `PLANS.md`·역할 인박스. GitHub Issues/Projects 연동은 미설정 |

## 기술 기준

개발 설계와 구현의 책임자는 dev다. 실제 요구사항과 기존 패턴을 확인하고 표준 라이브러리·플랫폼 기능·설치된 의존성부터 사용한다.
제품 코드·에셋·배포·테스트 경로는 기술 스택이 정해질 때 배정표에 등록한다. Unity 등 이전 참고 파일의 제품 가정은 상속하지 않는다.
인증값은 Git 미추적 `.env`로 제공하며 값·전문을 문서와 로그에 기록하지 않는다. FullOps 운영용 `.fullops-squad/.env`는 루트 `.env`의 로컬 링크다.
제품 비즈니스 판단은 designer, 기술 판단은 dev, 통합·배포는 ops, 독립 검증은 tester가 맡는다.

## 검증 명령

레포 루트에서 실행한다. `<플러그인>`은 해당 세션에 설치된 FullOps 패키지 경로를 조회해 사용하며 캐시 절대경로를 레포에 기록하지 않는다.

- 공백·패치 검사: `git diff --check` (스테이징 후 `git diff --cached --check`).
- FullOps lint: 깨끗한 커밋 상태에서 `python3 <플러그인>/scripts/lint.py --repo . --from <기준 SHA> --out <레포 밖 검증 결과 경로>`.
- 산출물 검사: `python3 <플러그인>/scripts/deliverables.py --repo . --strict`.
- 현황판 생성: `python3 <플러그인>/scripts/board.py --repo .`.
- 제품 lint·format·타입 검사·테스트·빌드·배포 명령: 미정. 현재 실행기·설정이 없어 `lint/lint.json`의 `commands`는 비워 둔다. 도입 시 실제 명령을 등록하고 실행 근거를 남긴다.

## 공통 기준의 적용과 예외

기존 프로젝트 규칙은 없으므로 공통 규칙을 기본값으로 적용한다. 이후 기술 정본이 생기면 연결하며 보안·권한·리뷰 수락 기준은 낮추지 않는다.
변경한 동작과 실패·경계 조건을 검증한다. 제품 도구가 없는 현 단계에는 설정·Git·하네스 검증을 수행하고 제품 테스트 성공으로 보고하지 않는다.
작업 지시서에는 적용 문서와 기준 SHA를 남기고 worker와 검토자가 같은 버전을 읽도록 한다.
