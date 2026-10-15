---
title: 잔여 UI 이슈 설계 검증 기록
status: review
updated: 2026-10-15
owner: designer
tasks: [SAR-UX-ISSUES-002-DESIGN]
summary: 문서 검사와 보호 원천·미실행·후속 검증 범위를 구분한다
---

# SAR-UX-ISSUES-002-DESIGN 검증 기록

## 판정 범위

제품 규칙·화면 UX·수락 조건의 문서 검토다. [설계](../../../design-docs/SAR-UX-ISSUES-002.md)에 #2–#8 원문 대조와 U2–U8, DEV/QA 순서를 연결했다. 제품·DB·운영·계정·호스트 설치를 실행하지 않았다. 직접 브라우저·접근성·새 UI 시각 검수는 not_run이다. 예전 실행 성공은 날짜·대상 SHA·한계를 유지한 참고이며 새 수락이 아니다.

## 근거와 확인

| 항목 | 결과·증거 |
|---|---|
| 이슈 원문 | [issues-open.json](issues-open.json), #2–#8 7건의 본문과 추가 요구를 읽음 |
| 현행 규칙 | 사람 친구 D02 HP-01–16과 D04 공개/사람 현행 절 직접 확인. 옛 개별 초대 요구의 대체 근거를 설계에 명시 |
| 현행 구현 | 이름 세 경로·친구 상태·receipt refusal·gate 검증·remote landmark 직접 읽음. 구현 확인과 실행 PASS를 구분 |
| frontend 준비 | codex/gpt-6-astra/high, frontend-design, tools=[]; `configuration_ready`, exit 0. authentication/browser_qa/mcp_connection은 not_run |
| 탐색 | [packet-screening.json](packet-screening.json): 386경로 존재/지문/제목 선별. unknown19·remaining349·optional0. 관련 본문과 범위 제외를 구분 |
| 보호 경로 | [preserved-inputs.json](preserved-inputs.json): 제품·D04·PLANS·board·타 역할 인박스 152경로의 착수 상태 지문 |
| 의존성 | 코드·manifest·lock 변경 없음. 신규 의존성·서드파티 설치 없음 |
| 신규 동작 테스트 | 미실행. 문서 전용 과제이며 새 제품 동작을 구현하지 않았음 |

자동 문서 검사와 현재 커밋 FullOps lint의 원본 출력·종료코드는 이 폴더의 검사 기록을 따른다. 실패·unavailable을 PASS로 바꾸지 않는다. 기록 검사는 실제 제품 수락이나 독립 리뷰가 아니다.

## 수동 의미 점검

- 이름과 불변 ID·친구 코드를 구분했다. handle/링크 미도입은 미구현을 지원 완료로 표현하지 않는 최소 결정이다.
- 수락 전 상대 미리보기 요구를 거절했다. 코드 교체·차단·세대·same owner·legacy 예외를 추적했다.
- 조회 오류의 복구와 비공개 비노출을 함께 정의했다. 승인 요약은 검증된 원요청만 사용한다.
- 언어 저장/fallback/계정 교대·원문 불변·시간대·누락 번역 판정을 명시했다.
- 다섯 호스트의 근거 상태와 설치부터 재개까지의 실패 경계를 분리했다. 실제 앱 버전 미기록은 미기록으로 남긴다.
- 실제 screen reader·200%·다른 브라우저 수락을 후속 QA로 명시했다. 현재 remote landmark가 이미 있다는 사실을 반영했다.

## 남은 일과 보존

DEV 구현·tester 독립 동작 QA·designer 실제 시각 검수·fixed-SHA 리뷰·main 통합은 후속이다. 실제 호스트 신규 설치와 계정 발송은 별도 권한이 필요하다. coor가 기존 인박스가 비면 기능 단위로 배정한다.

임시/검증 자료는 이 QA 폴더와 작업 트리별 Git 디렉터리의 FullOps lint 원본에 보존한다. 삭제·이동 정리는 하지 않았다. 원본 issues/route/find/context/packet은 그대로 커밋한다. 명령 로그는 실패를 포함해 보존하며 종료코드를 파이프로 가리지 않는다.

## 실행 결과

- `deliverables.py --repo . --strict`: 종료 0, 검사 13·문제 0·경고 0. [원본](deliverables.log).
- `git diff --check`와 staged diff 검사: 종료 0. [원본](diff-check.log).
- 로컬 문서 링크 17개·이슈 7개 추적·보호 경로 152개 지문: 통과. 착수 HEAD 대비 제품 diff0. [확인](document-check.json).
- packet outcomes: 403 path/category 쌍. 직접 읽은 근거와 범위 밖 보존 후보를 구분했다. 선택형 하위 위임 없음.
- 첫 고정 문서 SHA는 `34169e7f21b6d9658c5552fe8e6f1e4ff96daccd`다. 이 SHA의 lint 원본을 보존했다.
- 기준 `458eb565c4a1551ec58c70f14aa164d9024241d3`의 [원상태 lint](lint-initial.json)는 종료 1, ERROR 1·WARNING 3·unavailable 2다. 등록 `make lint`/`make test`는 현재 PATH에 make가 없어 실행 불가다. PASS로 기록하지 않는다.
- DOC-003은 타 과제 `SAR-LANDING-GUIDE-001-CODE-REVIEW-review/author-final-handoff.md`의 front matter 누락이다. 착수 HEAD에도 같은 내용이고 own diff0이다. [기준 대조](lint-baseline.json). 최신 지시의 소유권 제한으로 수정하지 않았다.
- SIZE-001 두 건은 기존 PLANS 1880줄과 D04 798줄이다. 둘 다 own diff0이며 축약하지 않았다. SIZE-002는 기준 ref 전체의 추가 742줄이다. 새 설계는 7개 이슈의 충돌·공통 경계·순차 인계를 함께 읽도록 유지했다. 원본 탐색·이슈 스냅샷은 증거 보존을 위해 축약하지 않았다.
- coor에 `msg_f8b96091c3e7`로 기존 lint 결함과 실행환경 처리를 인계했다. 검사 설정·심각도·PATH를 바꿔 통과시키지 않았다.

완료 기록 커밋 뒤 같은 기준의 최종 lint를 다시 실행한다. 원본은 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/ux-issues002-designer/lint-final.json`과 작업 트리별 Git 검사 기록에 남긴다. 최종 SHA·종료코드는 worker_done에서 회신한다. 최초 결과를 최종 실행으로 바꿔 쓰지 않는다.
