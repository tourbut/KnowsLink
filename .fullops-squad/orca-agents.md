---
title: Orca 역할 배정
status: draft
updated: 2026-10-09
owner: coor
tasks: [SAR-SETUP-001, FULLOPS-UPDATE-099, SAR-SETUP-001-DEV, FULLOPS-UPDATE-0.9.10, FULLOPS-UPDATE-1.2.0]
summary: 기존 역할과 제품·기술 책임 분리 및 고정 기준 인계·완료 계약
---

<!-- fullops-mode:start -->
## 운영 책임

현재 모드·주 담당자·테스트 레벨은 fullops.json을 읽는다. FULLOPS.md의 운영 모드 절을 우선 적용한다.
dev 주 담당자는 직접 구현과 전문가 배정·통합을 맡는다. dispatched dev-worker는 기존 worker 권한과 worker_done 계약을 따른다.
<!-- fullops-mode:end -->


# Orca 역할 배정

역할별 브랜치는 고유하게 유지한다. 책임 전환에서는 기존 브랜치와 워크트리를 재사용할 수 있다. `fullops.json`이 역할·브랜치 정본이다.
Orca 계층은 `coor` 아래 `designer`, `dev`, `ops`, `tester`다. 기본 브랜치 세션은 초기 설정용이며 운영은 coor에서 진행한다.

| 역할 | 책임·담당 경로 | 브랜치 | CLI | 모델 | 기동 방식 |
|---|---|---|---|---|---|
| coor | 과제 분해·배정·Run·진행 관리; `PLANS.md`, `handovers/`, `board/` | `fullops/coor` | `codex` | 모델 후보 | 상설 워크트리; 첫 요청 때 coordinator 세션 시작 |
| designer | 제품 기획·요구사항·사용자 경험; `docs/planning/`, `docs/design-docs/mockups/` | `fullops/designer` | `codex` | 모델 후보 | 기획 과제 발생 시 `worker-start --run` |
| dev | 기술 설계·구현·직접 검증; `docs/design-docs/` 기술 문서, `cmd/`, `internal/`, `adapters/`, `db/`, `scripts/`, 루트 Go·sqlc·Make·Docker·Compose 설정 | `fullops/dev` | `claude` 또는 `codex` | 모델 후보 | 기술 설계·구현 과제 발생 시 `worker-start --run` |
| ops | 배포·통합·운영; `docs/operations/`, 향후 배포 설정 경로 | `fullops/ops` | `claude` | 모델 후보 | 운영 과제 발생 시 `worker-start --run` |
| tester | 재현·테스트·회귀 검증; `docs/evaluations/scenarios/`, `docs/evaluations/qa-reports/`, 향후 테스트 경로 | `fullops/tester` | `claude` | 모델 후보 | 구현 SHA 준비 후 `worker-start --run` |

문서 경로는 `.fullops-squad/` 기준이다. 제품 경로는 레포 루트 기준이다. SAR-SETUP-001-DEV의 초기 구성·직접 검증 파일은 dev 소유다. 독립 QA의 시나리오·보고서는 tester 소유다. 운영 배포의 파일 소유권은 후속 ops 지시서에서 정한다.
원격은 `origin`, 기준 브랜치는 `main`이다. 상설 워크트리만 구성하고 이번 setup에서는 에이전트 세션을 시작하지 않는다.

## Windows 워크트리 복원 — 2026-10-09

원본 체크아웃은 `D:/workspace/KnowsLink`다. 아래 상설 워크트리는 Orca에 등록했다.

| 역할 | 로컬 워크트리 | 원격 추적 브랜치 | Orca 부모 |
|---|---|---|---|
| coor | `C:/Users/shin/orca/workspaces/KnowsLink/coor` | `origin/fullops/coor` | 없음 |
| designer | `C:/Users/shin/orca/workspaces/KnowsLink/designer` | `origin/fullops/designer` | coor |
| dev | `C:/Users/shin/orca/workspaces/KnowsLink/dev` | `origin/fullops/dev` | coor |
| ops | `C:/Users/shin/orca/workspaces/KnowsLink/ops` | `origin/fullops/ops` | coor |
| tester | `C:/Users/shin/orca/workspaces/KnowsLink/tester` | `origin/fullops/tester` | coor |

각 체크아웃은 기존 `fullops/<역할>` 브랜치를 사용한다. 기존 역할·모델 후보·coor 모드·테스트 lite를 유지한다.
coor와 dev의 Google 연결 WIP 및 역할 인박스는 보존한다. 기능 개발·검증·배포는 사용자 재개 지시까지 중지한다.
새 에이전트 세션은 시작하지 않았다. 기본 PowerShell 터미널의 실제 작업 경로를 확인했다.
원본에는 비공개 `.fullops-squad/.env*`가 없다. `env_link.py --all`은 변경 없이 종료했다. 운영 자격과 DB는 별도 이전 대상이다.
Grok의 원본 폴더 신뢰는 사용자 승인 후 등록한다. 승인 전에는 tester 실행 시 신뢰 확인이 필요하다.
과거 Linux 경로·Run·터미널 핸들은 재사용하지 않는다. 다음 배정 전에 현재 경로와 새 런타임 주소를 확인한다.

## 모델 후보

사용자가 제공한 후보의 약한 것부터 강한 순서를 유지한다. dev CLI는 선택한 후보에 따른다.
`worker-start --agent <CLI> --model <모델> --effort <effort> --run <run id>` 결과의 `launch.effective`로 실제 적용을 확인한다.

- `coor` `codex` `gpt-6-luna` `xhigh`: 과제 라우팅·Run 운영. 사용자 확정 배정
- `designer` `codex` `gpt-6.1-sol` `medium`: 단일 기능 기획, 단순 요구사항 명확화
- `designer` `codex` `gpt-6.1-sol` `high`: 여러 사용자 흐름 기획, 요구사항 수락 검토
- `dev` `claude` `claude-sonnet-5-5` `medium`: 한 줄 문구·수치·설정 변경
- `dev` `claude` `claude-sonnet-5-5` `high`: 기존 패턴을 따르는 단일 파일 소규모 구현·버그 수정
- `dev` `claude` `claude-opus-5-5` `medium`: 기존 패턴을 따르는 일반 기능 구현·검사 추가
- `dev` `codex` `gpt-6.1-sol` `medium`: 여러 스크립트에 걸친 기능 구현, 모듈 연결
- `dev` `claude` `claude-opus-5-5` `high`: 원인 추적이 필요한 버그, 여러 모듈 구현
- `dev` `codex` `gpt-6.1-sol` `high`: 공유 계약을 따르는 큰 기능 구현·리팩터링
- `ops` `claude` `claude-sonnet-5-5` `high`: 기록 검사·형식 검증·여러 SHA의 충돌 없는 main 통합 조정

- `tester` `grok` `grok-4.7` `high`: 독립 QA·시나리오·회귀 검증. 사용자 2026-10-03 지정

## 라우팅 기준

- coordinator 역할: `coor`
- 설계 역할: `designer`
- tester 역할: `tester`
- 제품 기획 역할: `designer`
- 기술 계획 역할: `dev`

- `coor`: 요청 접수, 과제 분해·배정, Run·진행·복귀 주소 관리. 제품 판단은 designer, 기술 판단은 dev에게 전달한다.
- `designer`: 제품 기획, 요구사항, 사용자 경험, 제품 범위·수락 기준을 맡는다. 기술 설계 요청은 dev에 넘길 수 있도록 제품 요구사항을 명확히 한다.
- `dev`: 개발에 대한 설계와 구현을 함께 맡는다. 기술 스택·아키텍처·API·데이터 모델·모듈 설계·제품 코드·직접 검증이 담당 범위다.
- `ops`: 배포·운영·수락된 SHA의 기계적 통합·동기화. 제품 의미 충돌은 designer, 기술 충돌은 dev로 돌린다.
- `tester`: 재현 조건·테스트 설계·테스트 코드·회귀 근거. 제품 코드 수정이 필요하면 coor를 통해 dev에 요청한다.

FullOps의 설계 역할 표시는 제품 기획 전용 designer를 뜻한다.
`product`는 제품 목표·규칙·수치·화면과 아트 방향·우선순위·완료 조건의 결정이다. designer가 맡는다.
`implementation`은 기존 제품 요구 안의 기술 계획·구조/API·버그 분석·구현·테스트·기술 문서 갱신이다. dev가 같은 과제에서 맡는다.
기술 난도가 높아도 기술 계획만을 이유로 designer에게 배정하지 않는다. designer는 상시 기술 승인자가 아니다.
`unresolved`는 책임 확인 전 보류한다. coor는 실제 소유 책임을 확인하고 재선정 근거를 라우팅 기록에 남긴다.
제품 규칙 변경·범위 확대·공유 제품 기준 질문은 coor를 통해 designer에게 전달한다. 기술 질문은 담당 구현자에게 전달한다.
ops는 배포·통합·운영의 기술 계획과 검증을 맡는다. tester는 독립 동작 QA를 맡고 제품 코드는 수정하지 않는다.
새 독립 코드 리뷰는 구현자와 다른 검토자의 별도 세션에서 수행한다. 검토자는 고정 SHA의 깨끗한 detached snapshot을 읽기 전용으로 사용한다.
리뷰 결과와 보고서는 기록 체크아웃에 작성한다. 별도 세션 ID·snapshot 경로·head·read_only를 기록한다.
필요한 직접 시각 검수와 독립 동작 QA를 수행한다. 미해결 critical/high가 있으면 수락·병합을 차단한다.
전환 전 리뷰 기록은 원본을 보존한다. 구현자 자기 리뷰로 독립 리뷰를 대체하지 않는다.
모델 후보는 최신 coor 브랜치의 사용자 결정(Astra 제외)을 보존했다. 실제 CLI 지원 여부는 기동 때 확인한다.

## 검증 담당과 후속 인계

DEV는 변경 동작의 자동 검사·관련 회귀·필요한 짧은 실행 확인을 완료한다. DEV 완료와 제품 최종 수락은 구분한다.
독립 전체 QA는 tester가 안정된 고정 통합 후보에서 수행한다. 직접 시각 검수는 designer가 담당한다. 별도 ART 역할은 구성하지 않는다.
coor는 검사별 담당·대상 SHA·실행 시점·통과 조건과 후속 인계 조건을 지시서에 기록한다.
캡처는 지정 시각 항목에만 만든다. 영상은 정지 화면으로 판정할 수 없는 항목에만 만든다.
변경 없는 증거는 관련 의존성의 동일성을 확인하고 원래 실행 SHA·조건을 연결해 재사용한다. 새 SHA에서 실행한 결과로 표시하지 않는다.
재검증은 변경 영향·새 실패·증거 결함·미충족 조건이 있을 때 수행한다. 기존 실패·held·미해결 critical/high·제품 정지·최종 플랫폼과 사람 평가 기준은 유지한다.
보류 항목에는 담당과 재개 조건을 남긴다. 상세 반복 범위는 [공통 테스트 기준](rules/common/testing.md)을 따른다.

## tester Grok 전환 — 2026-10-03

사용자가 tester를 Grok 4.7 high로 지정했다. `grok --model grok-4.7 --reasoning-effort high`로 실행한다. Orca가 grok의 model 인자를 직접 받지 않으므로 기존 tester 체크아웃에서 해당 명령으로 터미널을 만들고 준비를 확인한 뒤 worker-start --terminal로 감독한다. Claude 사용량 한도 대기 과제는 같은 Task의 retry로 이어받는다.

## 완료 작업의 즉시 통합과 원격 공유 — 2026-10-03 사용자 지시

사용자는 역할별 하위 워크트리에 개발을 배정하고, 완료 작업을 main에 즉시 병합해 origin에 올리도록 승인했다. 완료 보고를 받으면 coordinator는 해당 고정 SHA의 필수 검사·독립 QA·직접 시각 검수·독립 리뷰를 확인한다. 미해결 critical/high와 필수 실패가 없고 과제 수락 조건이 충족되면 main에 병합하고 origin/main에 push한다. 같은 범위의 병합·push 승인을 다시 묻지 않는다. worker_done만으로 필수 검증을 통과한 것으로 간주하지 않는다.

작업 완료 커밋과 준비·운영 기록은 해당 원격 역할 브랜치에도 공유한다. 원격 변경을 확인하고 force-push 없이 반영한다. idle이며 작업 트리가 깨끗한 상설 역할 워크트리는 최신 main으로 동기화한다. 진행 중인 worker의 체크아웃은 변경하지 않고 다음 자연스러운 착수 시점에 동기화한다. 다른 역할이 최신 코드를 받을 때 origin/main과 관련 고정 SHA를 전달한다.

상설 체크아웃 main과 coor/designer/dev/ops/tester를 우선 사용한다. 임시 워크트리는 결과·증거 병합과 원격 push를 확인한 뒤 Orca CLI로 정리한다. 미커밋 변경·진행 세션·필요한 검증 자료가 있으면 보존하고 재개 조건을 기록한다. 강제 reset·force-push·진행 작업 삭제는 하지 않는다.

## 완료 결과의 통합과 역할 인박스 — FullOps 0.9.12

실행 지시서는 역할별 `handovers/to_<role>.md` 하나로 고정한다. 다른 과제로 사용 중인 인박스를 덮어쓰지 않는다. 다음 과제는 PLANS.md에 대기시킨다. 과제명 파일·pending·logs는 참조 자료이며 배정 지시서로 사용하지 않는다. 역할 작업 완료 시 `work.py finish`로 지시서와 완료 보고 전문을 logs에 보존하고 인박스를 비운다.

coor는 worker_done을 받으면 고정 SHA의 필수 검토·검증을 확인하고 실제 기본 브랜치 main에 병합한다. 원격 origin/main에 일반 push한 뒤 완료 SHA의 로컬·원격 조상 관계를 확인한다. 역할 브랜치 push만으로 통합을 완료하지 않는다. 전체 제품 수락과 개별 문서 결과의 통합은 구분한다. 미해결 critical/high와 필수 검증 실패는 계속 차단한다.

병합 뒤 coor를 포함한 모든 등록 역할의 실제 worker 상태와 작업 트리를 확인한다. 쉬고 있으며 깨끗한 역할만 최신 main으로 동기화한다. 진행 중·미커밋 변경·상태 확인 불가인 역할은 PLANS.md에 최신 기본 SHA와 예약을 남긴다. 다음 dispatch 전에 최신 로컬·원격 main 포함을 확인한다. 진행 중 체크아웃과 사용자 자료를 덮어쓰지 않는다.

완료 메시지는 Git 공용 디렉터리의 fullops-integration에 보존된다. `integration.py --repo <레포> status`로 미통합 결과를 확인한다. 검수 대기·실패·충돌·원격 오류·사용자 제한은 PLANS.md에 메시지 ID·SHA·사유·담당·재개 조건을 기록한다. 같은 내용을 `integration.py hold`로 남긴다. 조건 충족 시 resume하고 병합·push를 이어간다. 성공 결과를 통합한 뒤 release·ack하고 다음 독립 과제를 배정한다. 사용자의 현재 과제 완료 뒤 중지 지시는 유지한다.

기존 활성 과제명 지시서는 작업 중 이동하지 않는다. 해당 과제 완료 뒤 logs에 보존하고 다음 과제부터 정규 역할 인박스를 사용한다.

## FullOps 1.0.0 이후의 인계·완료 계약

새 과제는 fresh 세션을 기본으로 한다. 같은 역할·관련 과제·작은 기존 컨텍스트·동일 모델/effort이고 직전 완료 후 몇 분 이내인 짧은 후속만 retain한다. 오래 기다린 후속, 큰 컨텍스트, 붙이기 실패는 새 세션을 사용한다. 새 세션 착수 확인 후 이전 완료 세션을 release한다. 독립 리뷰는 항상 작성자와 다른 세션을 사용한다.

1. 새 지시서는 `work.py new --base <40자리 고정 SHA>`로 만든다. 같은 과제 재작업은 `work.py reopen --base <고정 SHA>`로 attempt를 갱신하고 최초 보고·SHA를 보존한다.
2. spec에는 `Task key`와 `Purpose`를 명시한다. 지시서 작성 후 `jev_route.py --repo <루트> --key <과제 키> --bind-inbox --role <역할>`로 현재 인박스·attempt에 연결한다.
3. 탐색 패킷의 direct_edit/impact_check/document_read/document_update 경로·근거·partial/unknown과 남은 확인 목록을 worker에게 전달한다. 같은 입력의 탐색은 재사용한다. worker SHA·원천·지시가 바뀌면 find/context/packet을 force 갱신하고 이전 결과를 보존한다.
4. worker는 완료 전에 path/category별 completed/no_change와 사유를 packet-outcomes.json에 기록한다. unknown·누락은 완료로 처리하지 않는다. 동적 호출과 지원하지 않는 Markdown anchor는 직접 확인한다.
5. lint는 같은 checkout에서 직렬 실행한다. 중단 후 lock이 남으면 실제 lint 종료를 확인하고 해당 checkout의 Git 디렉터리 안 fullops-lint.lock만 제거한 뒤 재검사한다. 현재 HEAD와 지시서 기준 SHA의 통과 근거를 확인한다.
6. 완료 전송은 현재 task/dispatch로 `orchestration send --json`을 사용한다. 실제 성공 전달 receipt를 확인한다. failed 완료와 전송 실패를 구분한다. ask는 완료가 아니다. 전송 실패는 같은 dispatch에서 재시도한다.
7. 신규 리뷰는 `review.py snapshot`으로 명시 과제·40자리 SHA·별도 세션·clean detached 공간을 만든다. 결과·보고서·lint는 snapshot 밖에 쓴다. 검토·증거 보존·reviewer release 후 `review.py cleanup`으로 관리된 snapshot만 정리한다. historical check는 신규 수락 근거가 아니다.
8. 리뷰 보고서는 이전 대화를 참조하지 않고 정본 인덱스에서 현재 요구·결정 이유·구조·구현/미완료·실행/검증·운영/복구·다음 작업을 찾은 경로와 점검 결과를 기록한다. 누락·오래된 정보·깨진 링크·미지원 anchor를 구분한다.

제품 기획과 기술 계획의 기존 책임 분리 및 역할 인박스·main 통합 규약을 유지한다. 선택형 이슈 모드는 사용자 요청 때만 활성화한다. 자동 이슈 과제 `GH-<저장소 ID>-<번호>-A<attempt>`는 draft PR과 integration hold까지 처리하며 main 병합은 사용자 판단을 기다린다.
