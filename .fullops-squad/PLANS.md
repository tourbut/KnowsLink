---
title: KnowsLink 현재 계획
status: draft
updated: 2026-10-03
owner: coor
tasks: [FULLOPS-UPDATE-098, SAR-SETUP-001, SAR-SETUP-001-DEV]
summary: 초기 구성과 lint 설계 결과 및 dev와 tester의 선행 관계를 관리한다
---

# KnowsLink 현재 계획

## 운영

- orchestration Run: `run_8ca8bc058ab7`
- Run 목적: FullOps 초기 구성과 첫 개발 요청 준비.
- 현재 Run은 setup 세션에 연결돼 있다. 첫 coor 세션에서 `orca-ide orchestration run-use --id run_8ca8bc058ab7 --json`으로 연결한다.
- 운영 구조: coor 아래 designer / dev / ops / tester. 에이전트 세션은 첫 요청 전 시작하지 않는다.
- 완료: FullOps 0.9.7 설치(Codex·Claude Code·grok), 다섯 역할 setup, 원격 역할 브랜치 생성, 설정·검증 기준 작성.
- 완료: setup 커밋 `0c8da8d`를 main과 다섯 역할 브랜치에 push. Orca에서 coor를 부모로 네 worker 워크트리를 구성하고 동일 setup을 반영했다.
- 사용자 승인: 기존 미추적 `.gitignore`, `orca-agents.md`가 있는 상태에서 setup 진행 승인.

## 첫 요청 대기

- 제품 목적·요구사항: SAR-SETUP-001 D02에 초기 구성과 lint 범위를 기록했다. 사용자 기능은 후속 과제다.
- 기술 스택: 고정 원천의 Go·TypeScript·Postgres·Compose 선택을 따른다. 실제 버전·경로·검사 명령은 dev가 SAR-SETUP-001-DEV에서 확정한다.
- 배포 대상·운영 환경: ops가 배포 과제에서 확정.
- 재현·테스트·회귀 기준: tester가 구현 완료 SHA를 대상으로 확정.
- 현재 제품 과제: SAR-SETUP-001 설계 완료 후 SAR-SETUP-001-DEV, SAR-SETUP-001-TESTER 순서로 진행한다.
- 역할별 모델 후보는 사용자 참고 파일에서 유지했다. 실제 지원 여부는 과제 기동 시 확인한다.

## 워크트리·검증 근거

| 역할 | 워크트리 | 브랜치 | Orca 부모 |
|---|---|---|---|
| coor | `/home/shin/orca/workspaces/KnowsLink/fullops-coor` | `fullops/coor` | 없음 |
| designer | `/home/shin/orca/workspaces/KnowsLink/fullops-designer` | `fullops/designer` | coor |
| dev | `/home/shin/orca/workspaces/KnowsLink/fullops-dev` | `fullops/dev` | coor |
| ops | `/home/shin/orca/workspaces/KnowsLink/fullops-ops` | `fullops/ops` | coor |
| tester | `/home/shin/orca/workspaces/KnowsLink/fullops-tester` | `fullops/tester` | coor |

- Git 공백 검사와 JSON·역할·인박스·컨텍스트·모델 후보·라우팅·Run 정합성 검사 통과.
- `deliverables.py --strict`: 검사 13 / 미작성 13 / 문제 0 / 경고 0.
- setup 커밋 lint: ERROR 0 / WARNING 2 / 실행 불가 0. 신규 lint 설정(LINT-001), 기존 프로젝트 lint 도구 없음(LINT-000).
- 모든 워크트리의 브랜치·upstream·동일 HEAD·깨끗한 상태와 .env 링크 확인. 에이전트 세션은 기동하지 않았다.
- Jev OPENROUTER_API_KEY 존재 확인. 값은 출력하지 않았으며 실제 API 호출은 이번 setup에서 수행하지 않았다.
- 현황판: `.fullops-squad/board/index.html`. 생성 데이터는 Git 미추적이며 각 워크트리에서 board.py로 생성한다.

## FullOps 0.9.8 업데이트 — FULLOPS-UPDATE-098

- 업데이트 전 기준 커밋: `3c1a185`. 기본 브랜치는 main이며 기존 변경은 없었다.
- Codex 설치 버전은 0.9.8이다. `deps.py --check`에서 필수 CLI를 모두 확인했다.
- 진행 중 worker는 없었다. 다섯 역할의 인박스는 비어 있었다. 역할 워크트리는 모두 깨끗했다.
- 새 문서 작성 규칙을 추가했다. 하네스 지도에 규칙 링크를 추가했다. lint 안내에 DOC-003을 등록했다.
- 이번에 수정한 일반 문서는 `deliverables.py --stamp --path`로 등록했다. 기존 역할·모델·인박스·외부 공통 규칙은 유지했다.
- 기존 문서는 다음 수정 때 메타데이터를 등록한다. 진행 중 지시서의 일괄 변환은 보류한다.
- 동기화 대상은 coor, designer, dev, ops, tester다. 이 역할들은 쉬고 있으며 작업 트리가 깨끗하다. 예약된 동기화는 없다.
- 이 세션은 업데이트 전용이다. 새 coor 세션에서 Run `run_8ca8bc058ab7`을 연결한다.
- Claude Code와 grok은 이전 설치 버전 0.9.7이다. 해당 CLI로 worker를 시작하기 전에 README의 업데이트 명령을 실행한다.

## SAR-SETUP-001 — 프로젝트 초기 구성과 lint

- 사용자 요청의 GitHub 기획 문서를 커밋 고정 스냅샷으로 확보했다.
- Jev API 키가 없어 design → designer, codex gpt-6-astra high 폴백을 기록했다.
- designer는 제품 요구사항과 설정 범위를 확인하고 dev 및 tester 지시서를 작성한다.
- 새 과제이므로 designer 새 세션을 시작한다. 기준 ref는 `00b4cb34ae6e9f9fbc0b733ecaa3a2095fbc88eb`이다.
- designer 배정: Run `run_8ca8bc058ab7`, Task `task_49c6e00e6470`, Dispatch `ctx_1f186b4db5b5`, terminal `term_6935611a-04a4-4427-ba3b-d92dd4cf247e`. codex gpt-6-astra high 적용과 착수를 확인했다.

- 사용자 요청으로 모델 후보에서 gpt-6-astra를 전부 제외했다. 이후 배정은 갱신된 후보만 사용한다. 이미 착수한 designer Dispatch는 완료를 기다리고 후속 배정에서 변경한다.
- 원천 스냅샷은 외부 원본 보존 규칙에 따라 front matter를 추가하지 않는다. 프로젝트 lint에서 해당 스냅샷 경로만 제외하고 작성 문서는 계속 검사한다.

### SAR-SETUP-001 설계 결과 — 2026-10-03

- D02: `docs/planning/product-specs/SAR-SETUP-001.md`. 초기 개발 골격과 실제 lint 실패 검출만 수락 범위로 정했다.
- D03은 dev 책임이다. 잠긴 Go·TypeScript·Postgres·Compose·DB tooling·Go UI 선택을 유지한다.
- dev 지시서: `handovers/to_dev.md` (`SAR-SETUP-001-DEV`, ready).
- tester 지시서: `handovers/to_tester.md` (`SAR-SETUP-001-TESTER`, dev 완료 SHA 대기).
- 역할별 Jev find/code·documents 및 context는 API 실패로 fallback했다. 문서를 수동으로 좁히고 모두 keep했다. 원천은 수정하지 않았다.
- coor 후속: 설계 커밋을 역할 워크트리에 반영하고 실제 dispatch 복귀 정보를 기록한다. dev 완료 후 tester를 dispatch한다.
- 상세 근거와 검증: `docs/exec-plans/phases/SAR-SETUP-001.md`. 제품 기능·배포와 제품 테스트는 수행하지 않았다.
- 지정 기준 lint는 기존 원천 DOC-003 7건으로 차단됐다. coor가 원천 보존 및 차단 기록 후 설계 완료를 허용했다. 시작 HEAD 기준 설계 검사 ERROR 0, WARNING 1이다.
- coor 회신: 원천 전용 제외는 `729446d`에 기록했으며 기준 ref의 설정에는 아직 적용되지 않는다. dev는 제품 lint 등록과 함께 이 제약을 확인한다.
- 후속 모델 배정: 사용자 요청으로 제거한 Astra 후보를 재사용하지 않는다. coor가 갱신한 후보를 사용한다.


- 구현과 QA의 검사 기준은 원천 보존 설정을 포함한 준비 커밋 `729446d8da57`으로 갱신했다. 설계 단계의 원래 검사 실패 기록은 보존한다.
- dev 새 세션 배정: Run `run_8ca8bc058ab7`, Task `task_a7fd5d1b8806`, Dispatch `ctx_ab5236bf8b40`, terminal `term_5c300e96-6aca-49ad-ad46-441c378b7cee`. 갱신된 후보에서 codex gpt-6.1-sol high를 선택했고 착수를 확인했다.

## FullOps 0.9.9 업데이트 — FULLOPS-UPDATE-099

- 준비 기준은 `00b4cb34ae6e9f9fbc0b733ecaa3a2095fbc88eb`이다. 현재 Codex 플러그인은 0.9.9이며 필수 의존성을 확인했다.
- 제품 기획 designer / 기술 계획 dev의 합의된 책임 분리를 marker로 활성화했다. coor·설계·tester와 역할 브랜치 매핑은 유지했다.
- 이전 첫 요청 대기 기록은 당시 이력이다. 현재 SAR-SETUP-001 진행 자료는 coor·designer·dev 브랜치에 있다.
- ops·tester만 쉬고 깨끗한 동기화 대상이다. coor·designer·dev는 기존 SAR 인계와 작업을 보존하기 위해 완료 후 동기화한다.
- 새 coordinator는 [갱신과 인계 기록](docs/exec-plans/phases/FULLOPS-UPDATE-099.md)을 읽고 기존 Run을 연결한다.

## FullOps 0.9.9 coordinator 인계 — 2026-10-03

- 운영 동기화 기준: coor `3fde6f2`, main 준비 `bb31169`, 검증 `3d91754`. 병합 충돌은 기존 SAR 기록과 최신 역할 marker를 모두 보존해 해결했다. 제품 산출물은 main에 병합하지 않았다.
- Run `run_8ca8bc058ab7`을 현재 coordinator에 연결했다. designer 결과 `481d8ac8ac80a0c59bf35a0857eb9599ba6c50d7`의 지시서와 답 원문을 반영했다. 기존 승인·QA·원천·검사 기준 `729446d8da57`은 유지한다.
- 이전 dev Dispatch `ctx_ab5236bf8b40`는 failed이며 worker-show의 exactWorker observation `exited`, terminal exitCause `operator_close`를 확인했다. 작업 트리는 깨끗하며 완료 코드 SHA는 없다. 이전 user_takeover 기록은 유지한다.
- SAR-SETUP-001-DEV를 같은 키로 force 재선정했다. D02와 범위 답변에 근거해 override-role dev를 기록했다. Jev API 키 부재로 implementation/dev와 codex gpt-6.1-sol high 폴백을 기록했다. 기존 결과는 prior 파일로 보존했다. 갱신 산출물은 지시서의 D03을 유지한다.
- 같은 Task `task_a7fd5d1b8806`를 retry-of로 새 세션에서 재개한다. 이전 세션 종료 및 0.9.8 캐시 경로 오류 때문에 세션을 재사용하지 않는다. dev가 기술 계획·구현·테스트·문서 갱신을 수행한다.
- designer 후속 세션은 병합 검토까지 retain한다. tester는 성공 완료 SHA 이후 배정한다. 병합 전 서로 다른 구현·리뷰 세션 ID와 깨끗한 read-only detached snapshot, fixed-SHA 리뷰·lint·테스트 게이트를 확인한다.

- 재개 영수증: Task `task_a7fd5d1b8806`, Dispatch `ctx_67f98ed4cd42`, terminal `term_cf9a131a-84f0-40fa-91a2-c5bd7e0fb28f`. effective codex gpt-6.1-sol high, turn_started 확인. 기존 ready Task는 retry-of를 거절하므로 dispatch-show로 이전 실패를 확인한 후 같은 Task를 ready 상태에서 시작했다. 과제는 복제하지 않았다.
- 이전 질문 reply는 dispatch_inactive로 거절됐다. 답 원문은 designer 커밋과 최신 dev 지시서에 보존하여 새 세션에 전달했다. dev 워크트리는 종료와 깨끗한 상태를 확인한 후 준비 HEAD `dbe0b40`까지 fast-forward했다. designer 진행 브랜치는 변경하지 않았다.
- 운영 변경의 lint: HEAD `dbe0b40`, 기준 `729446d8da57`, 종료코드 0, ERROR 0, WARNING 1(LINT-000), 실행 불가 0. 제품 코드 변경은 없으므로 제품 테스트는 적용하지 않았다.
## SAR-SETUP-001-DEV 구현 결과 — 2026-10-03

- 초기 Go relay·별도 SQL-only migrate·TypeScript adapter·Compose·제품 lint를 구현했다. 업무 MVP와 운영 배포는 후속이다.
- 코드 체크포인트 `929832aa0ecd`에서 직접 검사·위반 주입·깨끗한 clone 재현·로컬 DB runtime을 통과했다. FullOps 기준 `729446d8da57`의 lint는 ERROR 0, WARNING 2, 실행 불가 0이다.
- D03은 review다. 독립 QA·고정 SHA 독립 코드 리뷰와 병합 수락은 대기한다. 최종 완료 SHA는 현재 Dispatch의 worker_done으로 coordinator에게 전달한다.
- 상세 근거와 한계: [실행 기록](docs/exec-plans/phases/SAR-SETUP-001-DEV.md).
