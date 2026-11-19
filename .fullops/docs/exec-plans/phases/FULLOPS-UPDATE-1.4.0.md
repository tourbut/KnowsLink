---
title: FullOps 1.4.0 저장소 적용
status: draft
updated: 2026-10-12
owner: coor
tasks: [FULLOPS-UPDATE-1.4.0]
summary: 공통 프런트엔드 준비와 웹 시나리오 안전 계약 및 통합 보류 기록
---

# FULLOPS-UPDATE-1.4.0

## 기준과 완료 조건

업데이트 전 실제 Codex 설치는 1.4.0이며 새 설치도 1.4.0이다. 레포 적용은 1.3.0이다. 기준 ref는 a29a7460b7226429357ce90e12bb56ce585f97ea다. fullops-common-0.3.3, FULLOPS.md, project.md, 문서 작성 규칙을 적용한다. 기존 역할·모델·lite/standard·제품 보류·진행 중 인박스와 사용자 변경을 보존한다. 제품 코드 변경은 없다.

## 릴리스별 판정

1.4.0을 적용한다. marketplace upgrade와 plugin add는 종료코드 0이며 실제 1.4.0 경로를 확인하고 새 스킬을 다시 읽었다. deps dry-run으로 세 호스트의 공통 스킬 범위를 확인했다. deps 설치와 Codex host check는 종료코드 0이다. setup은 기존 다섯 역할로 dry-run 후 실행했고 frontend.json 한 개만 추가했다.

프런트엔드 정본과 토큰은 기존 service-member-ui.md로 연결한다. 이번 과제는 UI 변경이 아니므로 enabled=false이며 host/model/effort는 다음 UI 과제에서 기존 배정으로 고정한다. 최초 frontend check는 미정 tuple·정본으로 needs_attention/1이었다. 비활성 설정과 정본 연결 후 다시 검사한다. 선택형 Playwright·Impeccable 엔진·shadcn 설치, React 전환, 브라우저 다운로드와 hook 신뢰 변경은 해당 없음이다.

실행 가능한 기존 Jev 웹 JSON 시나리오는 없다. 과거 Markdown QA 시나리오와 증거는 보존한다. 시나리오 안내를 새 템플릿과 통합하여 covers·독립 checks·격리 fixture·정확한 권한·기존 탭·assert-only와 예산 판정을 반영했다. 원격 효과 있는 Orca 조작은 차단하며 탐색·shadow·CUAWright는 OFF다. 실제 웹 실행은 UI 변경이 없어 해당 없음이다.

## 검증과 통합 보류

관련 설정·문서 검증 완료 후 레포 적용 버전을 갱신한다. 제품 필수 명령의 결과는 별도로 기록하고 미실행을 성공으로 표시하지 않는다.

기존 main 체크아웃 D:/workspace/KnowsLink에는 미커밋 .codex/hooks.json이 있다. 현재 coor에는 독립 tester 수락 전 hook 수정 후보가 포함돼 있다. 업데이트와 함께 이를 main에 수락하지 않는다. 업데이트 커밋의 독립 리뷰 및 기본 브랜치 병합·원격 push는 담당 coor가 이어간다. 재개 조건은 기존 hook 과제의 독립 수락과 main 사용자 변경 해소 또는 깨끗한 업데이트 전용 통합 경계 확보다. 역할 동기화는 해당 main 통합 뒤 실제 유휴·clean 확인 후 진행한다. 진행 중 worker는 시작하거나 중단하지 않는다. 다음 개발은 새 coordinator 세션에서 진행한다.

설정 검사 frontend.py --check는 disabled·problems 0·종료코드 0이다. 문서 strict는 검사 13·문제 0·경고 0이며 git diff --check는 종료코드 0이다. 관련 운영 적용을 완료하여 plugin_version을 1.4.0으로 갱신했다. 인증·브라우저 QA·MCP 실행 검증은 수행하지 않았다.

### 커밋 후 필수 검사 결과

적용 커밋은 ddb41325a33f3ee418671595b2ff002a203dc652다. lint.py --from a29a7460b7226429357ce90e12bb56ce585f97ea는 종료코드 1이다. ERROR 0, 기존 PLANS 크기 WARNING 1, make 부재로 product-lint/product-test unavailable 2다. 제품 검사 통과가 아니며 제품 코드 변경이 없는 운영 적용 검증과 구분한다. 원본 로그는 C:/Users/shin/AppData/Local/Temp/fullops-update-1.4.0-lint.json에 보존하며 리뷰 폴더에도 복사한다. 현황판 생성은 종료코드 0이다. 독립 리뷰 snapshot은 삭제하지 않고 보존한다.

독립 검토자 thread 01a12eae-291b-7590-b4f0-3def2d888542가 적용 SHA의 7파일을 검토했다. findings 0이며 review.py check는 reviewed 7/skipped 0/WARNING 1로 통과했다. 증거는 docs/evaluations/qa-reports/FULLOPS-UPDATE-1.4.0-review/다. snapshot C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-20ffc8d847b142b9bd4608c730b66cfa는 보존한다. 별도 Orca dispatch가 없는 독립 검토 세션이므로 자동 cleanup release 계약을 충족하지 않으며 담당 coor가 안전한 정리 경계에서 확인한다. 이 후속 커밋은 검증 결과 보존 전용이며 적용 코드를 바꾸지 않는다. main 통합·push·역할 동기화 hold는 유지한다.
