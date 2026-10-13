---
title: FullOps 1.4.2 저장소 적용
status: draft
updated: 2026-10-14
owner: coor
tasks: [FULLOPS-UPDATE-1.4.2, SAR-REPO-DOCS-001-DEV]
summary: 최신 의존성 갱신과 1.4.1·1.4.2 적용 및 기존 문서 결과 통합
---

# FULLOPS-UPDATE-1.4.2

## 기준과 보존

업데이트 전 실제 Codex 설치는 1.4.2이며 레포 적용은 1.4.0이다. 설치 대상과 레포 목표는 1.4.2다. 사용자 Codex 홈과 Orca 전용 Codex 홈의 실제 plugin.json·설치기 결과를 확인했다. 사용자 요청으로 운영 업데이트 뒤 이미 완료된 README·docs 결과의 검토·main 통합을 이어간다. 새 제품 worker는 시작하지 않는다.

기준 coor HEAD는 a2961818ee33f3786aa6f949c5b411bd4b31aa1a다. DEV 완료 SHA ea80e98510af0e640557f0ea2b757c06e27a5779와 미완료 리뷰 폴더·읽기 전용 snapshot을 보존했다. 사용자 소유 dev/ops/designer 터미널은 종료·이동·덮어쓰기하지 않는다. origin/main, 다섯 역할, coor/lite/standard, Sol high, 인증·디자인 정본·제품 코드·설정은 유지한다.

## 릴리스별 적용

| 릴리스·항목 | 판정·근거 |
|---|---|
| 1.4.1 latest 외부 의존성 | 적용. 통합 설치기 update가 사용자·Orca Codex 홈별 실제 deps 설치·host check를 실행했고 종료코드 0이다. 기존 스킬은 각 홈의 fullops-skill-backups에 보존했다. |
| 1.4.1 선택형 프런트 도구 | 해당 없음. frontend.json tools=[]다. 설치된 공통 디자인 스킬은 갱신했고 기존 host/model/effort·정본은 보존했다. |
| 1.4.2 플러그인·설치기 | 적용. fullops-squad@latest update --hosts codex --repo로 설치·갱신했다. 두 홈 모두 실제 1.4.2다. 설치기 버전은 별도 검증 자료에 기록한다. |
| 1.4.2 macOS cwd | 해당 없음. Windows이며 lint.json 두 검사 cwd는 `.`다. macOS 실행 성공은 주장하지 않는다. 새 패키지 경로 정규화는 그대로 사용한다. |
| 1.4.2 Jev 입력·Enter·허용 동작 | 실행기 반영. 기존 시나리오 안내를 확인했고 이 문서 과제에 실행할 웹 시나리오는 없다. 입력 예산·위험·confidence·origin·독립 checks·CUAWright 기본값을 바꾸지 않았다. |
| 이전 업데이트 보류 | 보존. 1.4.0 당시 기록은 과거 사실이며 후속 main 반영 기록을 유지한다. 실제 hook trust·오류 소멸과 장시간 Events 등 별도 제품 미검증을 이번 업데이트로 완료 처리하지 않는다. |

## 적용과 검증

설치된 새 update-fullops와 INSTALL.md를 읽고 update.py --repo . --from 1.4.2로 레포 기준 이후 두 릴리스를 순서대로 확인했다. Context7 resolve는 Monthly quota exceeded였다. 설치된 공식 패키지의 INSTALL.md·releases를 근거로 사용했고 외부에 비공개 코드나 인증값을 전송하지 않았다.

기존 다섯 역할·origin/main·coor/lite/standard로 setup dry-run 후 실행했다. 기존 FULLOPS.md의 같은 모드 절과 fullops.json만 대상이었고 실질 변화는 mode=coor 명시뿐이었다. 기존 규칙·review rule·사용자 AGENTS와 Windows guard는 유지했다. FULLOPS.md에 latest 정책과 새 실행기·세션 경계를 추가했다.

통합 설치기 종료코드 0, 현재 Codex deps host check 종료코드 0, frontend check configuration_ready/problems 0이다. frontend 인증·브라우저 QA·MCP 연결은 not_run이다. 필요한 정적 검사와 통합 SHA는 아래 완료 기록에서 갱신한다. 필수 적용과 검증 확인 뒤 plugin_version만 1.4.2로 갱신한다.

## 통합·동기화와 다음 작업

README·docs 결과는 다른 실제 세션에서 고정 SHA snapshot으로 리뷰한다. 통합은 main 일반 push와 로컬·원격 조상 검증으로 완료한다. dev/ops/designer는 user_owned 또는 상태 확인 불가이므로 다음 배정 전 최신 main 동기화를 예약한다. tester는 실제 유휴·clean 확인 뒤에만 동기화한다.

리뷰 snapshot과 설치기 스킬 백업은 보존한다. 정리 담당 coor는 실제 세션 종료·소유권·증거 보존 확인 때 정리한다. 미검증 항목의 담당과 재개 조건은 기존 PLANS를 따른다. 제품 개발 재개는 최신 hook을 로드한 새 coordinator 세션에서 수행한다.
