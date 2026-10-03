---
title: FullOps 0.9.12 업데이트 적용 기록
status: draft
updated: 2026-10-03
owner: coor
tasks: [FULLOPS-UPDATE-0.9.12]
summary: 완료 결과의 기본 브랜치 통합 정책과 역할 인박스 수명주기 및 기존 미통합 결과를 기록한다
---

# FULLOPS-UPDATE-0.9.12 — 운영 레포 적용

## 기준·설치·범위

coor 기준 ref는 `5b94be5e0cb1757911489909b7c591c2bff8af90`이다. 운영 준비·main 검증 기준 ref는 `83aec58fcae69bed349fcc7f21057e1947123406`이다. 설치 전 실제 plugin.json 버전은 0.9.12다. coor 레포 적용 버전은 0.9.11이고 main은 0.9.10이다.

Codex marketplace upgrade는 이미 최신이며 종료코드 0이다. plugin add는 `/home/shin/.codex/plugins/cache/fullops-squad/fullops-squad/0.9.12` 설치에 성공했다. 실제 설치 후 버전은 0.9.12다. 새 update-fullops를 다시 읽었다. deps.py --check는 종료코드 0이며 필수 CLI가 모두 있다. update.py --repo . --from 0.9.12로 coor 적용 기준 이후 릴리스 0.9.12를 확인했다. main의 미적용 0.9.11은 기존 기록과 릴리스 본문을 함께 회수한다.

## 릴리스 판정·변경 파일

0.9.11의 글로벌 npm 설치 편의 기능은 제품 코드·역할·검증 기준 변경 해당 없음이다. 이전 준비 SHA의 실제 lint 통과 결과를 0.9.11 기록에 복구한다.

0.9.12 정책은 적용한다. FULLOPS.md·orca-agents.md·handovers/_TEMPLATE.md에 즉시 main 병합·일반 push·조상 확인·역할별 동기화·명시적 hold와 역할 정규 인박스 규칙을 반영한다. PLANS.md에 미통합 완료 SHA·메시지·담당·재개 조건과 역할 동기화 예약을 기록한다. board/board.json도 운영 상태를 맞춘다.

같은 coor/designer/dev/ops/tester 역할·origin·main으로 setup dry-run을 실행했다. 종료코드 0이며 생성 예정 0개다. 신규 파일이 없으므로 setup 실행은 해당 없음이다. 기존 운영 문서에 필요한 절만 통합한다. 필수 정책 검증을 통과한 뒤 fullops.json의 버전을 0.9.12로 갱신한다.

## 기존 작업 보존·보류

PLANS의 미통합 결과 표를 정본으로 사용한다. DEV·QA·공개 기획 완료 메시지 세 건을 원래 ID와 SHA로 integration 기록에 회수했다. 독립 리뷰·UI 또는 별도 기획 검토 미완료에 따라 사유·coor 담당·재개 조건을 가진 hold를 기록했다. QA 31 pass/8 held/0 fail과 이전 실패·제품 중지 지시는 보존한다. worker-list에서 리뷰와 OPS는 failed/exited, UI는 failed/unverifiable이다. 사용자 자료와 실패를 성공으로 바꾸지 않는다.

기존 완료 작업명 핸드오버 두 건은 이미 logs에 있다. SAR-MVP-001-REVIEW는 실패한 미완료 지시서이므로 이동하지 않는다. 미배정 후속 작업은 PLANS에서 추적한다. 활성·미완료 인박스는 덮어쓰지 않는다. OPS·designer의 미커밋 증거와 coor의 미추적 리뷰 증거는 보존한다.

제품 코드를 포함한 누적 coor 브랜치 대신 main 기준 운영 준비 브랜치 fullops-update-0.9.12에서 운영 정책만 작성한다. 필수 검증 후 실제 main에 merge·일반 push하고 원격 조상 관계를 확인한다. 동기화는 실제 유휴·깨끗함을 확인한 역할만 수행한다. 상태 불명·변경 있음·진행 중 역할은 coor가 안전한 다음 배정 전에 재개한다. 신규 제품 worker·배포·정책 수치 확정은 수행하지 않는다.

## 검증

제품 코드 변경이 없으므로 제품 동작 테스트·시각 검수는 해당 없음이다. 운영 문서 메타데이터·Git 공백·JSON·역할 및 원격 보존·링크와 main 기준 lint를 확인한다. 준비 결과와 실제 통합·push·역할 동기화 결과는 아래에 이어 기록한다.

준비 SHA `e6aed25b7eb6`의 main 기준 FullOps lint는 종료코드 0, product-lint passed, ERROR 0, WARNING 0, 실행 불가 0이다. JSON·역할/원격 보존·변경 문서 메타데이터·상대 링크·Git 공백 검사를 통과했다. 필수 정책 적용과 검증을 완료해 적용 버전을 0.9.12로 갱신한다. 이 버전은 제품 완료나 hold 해제를 뜻하지 않는다.

## 실제 통합·원격·역할 동기화 결과

운영 준비 SHA `e6aed25b7eb6`와 적용 버전 확정 SHA `66f7ffcaed4c230424962fe05779aac16d5cbe90`를 실제 main 체크아웃 `/home/shin/Workspace/KnowsLink`에서 fast-forward로 통합했다. origin/main 일반 push는 종료코드 0이다. fetch 뒤 origin/main이 66f7ffc이며 운영 준비 SHA의 로컬·원격 조상 관계를 확인했다. 미수락 제품 후보 a6a10c7·QA c59537b·기획 69dbec4는 main에 포함하지 않았다.

coor는 main을 merge해 `c9b2a8e`로 동기화했다. 충돌은 최신 역할/모델·제품 상태·PLANS·미추적 증거를 유지하며 새 운영 정책을 통합해 해결했다. dev는 Orca terminal list가 0개이고 작업 트리가 깨끗함을 확인했다. main을 merge해 `25f03ad81d438b2cd816168770a4fe384fa1f9f8`로 동기화했다. 두 역할은 기존 HEAD와 운영 main을 모두 조상으로 보존한다.

designer·ops는 미커밋 자료 때문에, tester는 liveness unverifiable 때문에 동기화를 예약했다. 임시 designer-pilot·dev-mvp도 실제 상태 불명과 제품 수락 대기 때문에 보존한다. 역할별 예약의 담당은 coor이며 실제 유휴·깨끗함을 확인한 다음 dispatch 전에 최신 main을 반영한다. 역할 브랜치 원격도 일반 push로 운영 준비를 공유한다. 실제 제품 hold는 유지한다.

최종 기록 커밋에도 같은 main 기준 lint를 실행한다. coor는 사용자 미추적 자료를 유지하므로 같은 고정 HEAD의 깨끗한 검증 snapshot에서 검사한다. 새 coordinator 세션은 0.9.12 hook을 사용해야 한다. 이 세션이 참조하는 삭제된 0.9.11 hook 캐시는 수정하거나 우회하지 않는다.
