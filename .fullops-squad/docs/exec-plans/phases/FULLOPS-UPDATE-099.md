---
title: FullOps 0.9.9 갱신과 coordinator 인계
status: draft
updated: 2026-10-03
owner: coor
tasks: [FULLOPS-UPDATE-099]
summary: 제품 기획과 기술 계획의 책임 분리 및 기존 SAR 작업을 보존한 동기화 계획을 기록한다
---

# FULLOPS-UPDATE-099 — 준비 갱신

## 기준과 적용

- 기준 ref: `00b4cb34ae6e9f9fbc0b733ecaa3a2095fbc88eb`.
- 기본 브랜치: 원격 HEAD `main`. 원격: `origin`.
- 현재 Codex 세션의 설치된 plugin.json과 로드된 스킬 경로에서 0.9.9를 확인했다. 현재 CLI 업데이트는 필요하지 않다.
- 설치된 `deps.py --check`는 종료코드 0이며 필수 CLI를 모두 확인했다.
- Claude Code 설치 기록은 0.9.7이다. 해당 호스트의 첫 후속 세션 전에 플러그인을 갱신한다.
- README의 Claude 업데이트 명령: `claude plugin marketplace update fullops-squad && claude plugin update fullops-squad@fullops-squad`.
- 이번 작업은 setup 운영 갱신이다. Jev 개발 요청 분류와 worker dispatch는 적용하지 않았다.
- 설치된 setup-fullops, fullops-orca, fullops-review, orca-cli 규약을 확인했다.
- 동일한 다섯 역할로 remote setup dry-run과 실행을 수행했다. 두 명령은 종료코드 0이다. 생성 브랜치와 파일은 각각 0개다.
- 템플릿 41개를 비교했다. 사용자 문서·모델 후보·규칙·인박스·컨텍스트·기록을 자동 교체하지 않았다.
- 역할·브랜치 매핑은 유지했다. 역할 브랜치 다섯 개에 중복이 없다.
- product_roles 판정은 designer/dev다. coordinator는 coor, 설계 역할은 designer, tester는 tester다.
- 기존 orca-agents.md에서 이미 합의된 제품/기술 책임 분리를 활성화했다. 최신 coor 브랜치의 Astra 제외 결정을 함께 보존했다.
- 신규 독립 리뷰는 별도 검토자 세션과 고정 SHA의 깨끗한 detached snapshot을 사용한다. 기존 리뷰 기록은 보존한다.

## 기존 진행 자료와 동기화

모든 워크트리의 작업 트리는 조사 시 깨끗했다. 깨끗한 작업 트리만으로 쉬는 역할을 판정하지 않았다.
Orca Run `run_8ca8bc058ab7`의 task-list, worker-list, terminal list를 함께 확인했다.

- coor: `3fde6f26130b35b76c6b7488f66a335d2b270ad0`. SAR 진행 조정과 열린 인계가 있다. 이번 준비 동기화는 보류한다.
- designer: `481d8ac8ac80a0c59bf35a0857eb9599ba6c50d7`. 설계·범위 질문의 완료 SHA와 retained 기록을 보존한다. 미통합 SAR 인계가 있어 동기화를 보류한다.
- dev: `c99fe83ca1d3155fb8f87498564dd138a00e3e6d`. SAR-SETUP-001-DEV가 ready다. 이전 worker는 failed이며 user_takeover로 retained됐다. 후속 구현 완료 후 동기화한다.
- ops와 tester: `00b4cb34ae6e9f9fbc0b733ecaa3a2095fbc88eb`. 활성 세션과 진행 인박스가 없으며 작업 트리가 깨끗하다. 준비 커밋을 fast-forward로 동기화한다.

인박스·컨텍스트·기획 원천·검증 기록은 변경 전 SHA와 파일 해시로 보존을 확인한다.
기존 제품 과제의 기준 ref `729446d8da57`와 외부 원천 SHA를 변경하지 않는다.
main에는 제품 과제의 수락 전 산출물을 병합하지 않는다. main 현황판은 역할 브랜치의 SAR 진행과 동기화 보류를 안내한다.

## 새 coordinator 세션

1. 새 coor 세션을 연다. 현재 업데이트 세션을 coordinator로 재사용하지 않는다.
2. `orca orchestration run-use --id run_8ca8bc058ab7 --json`으로 기존 Run을 연결한다. 실제 환경의 Orca 실행 파일을 사용한다.
3. task-list와 worker-list를 조회하고 failed/user_takeover의 실제 작업 상태를 확인한다. 이전 세션·과제를 복제하지 않는다.
4. designer 완료 SHA와 coor의 지시서 차이를 확인한다. SAR-SETUP-001-DEV의 후속 작업과 모델 후보를 기존 사용자 결정에 맞춰 재개한다.
5. coor·designer·dev의 진행 인계가 끝나고 역할이 쉬며 깨끗할 때 준비 커밋을 병합한다. 충돌이 있으면 원천·검증 근거와 최신 모델 결정을 보존한다.
6. 구현 완료 후 독립 리뷰와 tester 동작 QA를 수행한다. 필요하면 직접 시각 검수를 수행한다. 미해결 critical/high를 차단한다.

준비 커밋은 이 문서의 Git 이력으로 식별한다. 검증 결과와 동기화 영수증은 후속 운영 기록으로 연결한다.
제품 코드 변경은 없다. 제품 테스트·빌드는 현재 main에 실행기가 없어 적용하지 않는다.

## 적용 결과

- 준비 커밋: `bb311695d020ee99a58c0112261eed4b545e48c2`. 아래 실제 SHA 기록이 정본이다.
- main과 ops·tester에 준비 커밋을 push했다. ops·tester는 fast-forward로 반영했다.
- coor·designer·dev의 HEAD와 진행 자료 해시는 변경 전 값과 같다. 모든 기존 인박스·컨텍스트·원천·검증 파일의 해시가 일치한다.
- 준비 커밋의 지정 기준 FullOps lint: 종료코드 0, ERROR 0, WARNING 1, 실행 불가 0.
- WARNING은 LINT-000이다. main에는 제품 검사 실행기가 없어 기존 commands를 비워 두었다. 진행 중 dev가 제품 검사 명령을 구성한다.
- 산출물 strict: 검사 13, 미작성 13, 문제 0, 경고 0. 이는 main의 상태이며 역할 브랜치의 진행 산출물을 부정하지 않는다.
- Git 공백 검사와 최종 setup dry-run을 통과했다. 역할 판정과 브랜치 고유성 검사를 통과했다.
- 검증 근거: `docs/evaluations/qa-reports/FULLOPS-UPDATE-099/`의 lint.json, preservation-before.json, template-comparison.txt.
- 후속 검증 기록 커밋도 main과 쉬는 ops·tester에 반영한다. coor·designer·dev에는 두 커밋 모두 완료 후 반영한다.
