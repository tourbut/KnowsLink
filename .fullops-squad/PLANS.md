# KnowsLink 현재 계획

## 운영

- orchestration Run: `run_8ca8bc058ab7`
- Run 목적: FullOps 초기 구성과 첫 개발 요청 준비.
- 현재 Run은 setup 세션에 연결돼 있다. 첫 coor 세션에서 `orca-ide orchestration run-use --id run_8ca8bc058ab7 --json`으로 연결한다.
- 운영 구조: coor 아래 designer / dev / ops / tester. 에이전트 세션은 첫 요청 전 시작하지 않는다.
- 완료: FullOps 0.9.7 설치(Codex·Claude Code·grok), 다섯 역할 setup, 원격 역할 브랜치 생성, 설정·검증 기준 작성.
- 이번 구성: 기본 브랜치에 setup 커밋·push 후 각 역할 워크트리에 동일 커밋 반영. Orca에서 coor를 부모로 구성한다.
- 사용자 승인: 기존 미추적 `.gitignore`, `orca-agents.md`가 있는 상태에서 setup 진행 승인.

## 첫 요청 대기

- 제품 목적·요구사항·사용자 경험: designer가 첫 요청에서 기록.
- 기술 스택·기술 설계·제품 경로·lint/test/build 명령: dev가 첫 개발 과제에서 확정.
- 배포 대상·운영 환경: ops가 배포 과제에서 확정.
- 재현·테스트·회귀 기준: tester가 구현 완료 SHA를 대상으로 확정.
- 진행 중 제품 과제·worker Dispatch: 없음.
- 역할별 모델 후보는 사용자 참고 파일에서 유지했다. 실제 지원 여부는 과제 기동 시 확인한다.
