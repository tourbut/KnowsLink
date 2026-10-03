---
title: SAR-MVP-001-UI — 직접 시각 검수 실행 기록
status: draft
updated: 2026-10-03
owner: designer
tasks: [SAR-MVP-001-UI]
summary: 실패한 UI 검수 Dispatch의 기록을 보존하고 재개 검증과 인계를 기록한다
---

# SAR-MVP-001-UI — 직접 시각 검수 실행 기록

## 이전 Dispatch의 범위와 적용 기준

착수 HEAD는 `661966f5646ecd44d108beddefed0630b65076c0`, 브랜치는 `fullops/designer`다.
제품 후보는 `a6a10c71977b7f3ec8274a1fb7c8a409f58e7c92`다.
QA 원래 증거 커밋 `c59537b6fa0c7e008c4c6bdba0a251dd821d4ee8`와 정책 커밋 `69dbec44c0193266f8f6c8499f22493e1e3c1722`의 포함을 확인했다.
공통 규칙 `fullops-common-0.3.2`, FULLOPS.md, project.md, 문서 작성 규칙, D02, tester 보고서와 designer 컨텍스트를 읽었다.
기준 ref는 `0dd08ec994771836c15d9d22a6a83393a71d7987`이다. 기존 탐색 근거를 재사용했다.
이전 Dispatch의 기록에는 fullops-work·fullops-deliverables·fullops-test, orca-cli·orchestration, Ponytail full·caveman full 적용이 남아 있다. 이번 재개는 fullops-work·fullops-deliverables·orchestration과 caveman full을 적용한다.
문서 내용은 일반 한국어 STE로 작성했다.

## 결정과 증거

기존 V-01–04 PNG 7개를 실제 `view_image`로 직접 확인했다.
pending의 기존 900px 캡처는 버튼 일부가 잘려 같은 저장 HTML의 1100px 렌더를 한 번 보완했다.
Chrome 프로세스 종료코드 0을 직접 확인했고 보완 PNG도 `view_image`로 확인했다.
기존 증거는 수정하지 않았다. 시간 변화 검수나 새로운 동작 QA가 아니므로 영상과 전체 QA 반복은 하지 않았다.
[D04](../../design-docs/mockups/SAR-MVP-001-UI.md)와 [시각 판정](../../evaluations/qa-reports/SAR-MVP-001-UI.md)에 경로와 원래 SHA를 기록했다.
캡처된 합성 V-01–04는 PASS다. 미해결 critical/high 시각 결함은 발견하지 않았다.
401 `invalid_auth`는 인증 실패 차단의 증거다. 인증된 owner의 권한 저장소 장애 화면으로 오인하지 않는다.
중복 승인·GET 비결정·CSRF·현재 권한의 서버 동작은 기존 QA-08과 원래 held를 구분해 재사용한다.

## 파일 소유권과 D04 메타데이터

변경 범위는 D04 원천과 보완 PNG, 시각 판정, 실행 기록, designer 컨텍스트와 자기 인박스·아카이브다.
제품 코드·기획 정책·원천 스냅샷·PLANS·board·기존 QA 증거·coor 소유 산출물 인덱스는 보존한다.
D04 작성은 이번 지시서의 명시적 허용이다. D02의 이전 미작성 문장을 이번 작업에서 수정하지 않는다.
D04 front matter는 deliverables.py --stamp로 생성한다. 자동으로 바뀐 인덱스는 원래 내용으로 보존한다.
공유 인덱스의 D04 행은 미작성 이력이 남는다. coor가 병합 때 새 D04 review 원천과 대조할 항목이다.
전체 strict의 D04 미작성 수치를 새 원천 검사의 통과 근거로 사용하지 않는다.
D04 메타데이터와 정규 front matter는 별도로 검사한다.

## 검증 범위

문서 strict, D04 메타데이터, 변경 문서의 로컬 링크, git diff --check, 파일 소유권과 빈 인박스를 확인한다.
검사 명령의 직접 종료코드를 남긴다. 파이프로 종료코드를 숨기지 않는다.
최종 깨끗한 커밋에서 FullOps lint를 기준 ref로 실행하고 고정 SHA·결과를 worker_done으로 전달한다.
검증 로그는 레포 밖 `/tmp/SAR-MVP-001-UI-validation/`에 보존한다.
제품 코드 변경이 없으므로 개발 build/test/runtime과 새 독립 동작 QA는 반복하지 않는다.
제품 코드 변경 시 done-gate의 의무는 이번 문서·이미지 작업에 미적용이다.
이번 기록 검증은 제품 정책 held나 독립 리뷰 수락을 대체하지 않는다.

## 완료와 중지 조건

완료 보고를 자기 인박스에 쓰고 work.py finish로 이번 과제를 한 번 아카이브한다.
소유 파일만 커밋한 뒤 최종 SHA·시각 판정·검사 결과·인계 링크를 새 Dispatch의 worker_done으로 보낸다.
Run은 `run_8ca8bc058ab7`, Task는 `task_911c61d88587`, Dispatch는 `ctx_b7d073da41ef`다.
coor가 기존 필수 독립 리뷰·검사를 확인하고 현재 과제를 병합·원격 공유한다.
후속 기능·배포·제품 수치 확정은 시작하지 않는다. 기존 held와 사용자 중지 지시는 유지한다.

## 실패한 Dispatch의 기록 마무리 재개

이전 Dispatch `ctx_b7d073da41ef`는 failed다. 위 직접 열람과 보완 PNG 생성은 이전 관찰로 보존한다.
이번 착수 HEAD는 `1067ccd`다. 현재 지시서의 재개 기준 ref는 `0a5b044`다.
새 Dispatch는 `ctx_f6be2c2d0ce4`, Task는 `task_911c61d88587`다.
기존 미추적 D04·시각 판정·실행 기록·보완 PNG를 보존해 완성한다.
이번 세션에서도 기존 PNG 7개와 보완 pending PNG 1개를 `view_image`로 직접 열었다.
발신·대상, typed body, 정책, 상태 문자, 결정 차단과 pending의 두 버튼을 직접 확인했다.
합성 캡처의 V-01–04 시각 판정은 PASS다. 새 캡처나 서버 실행은 수행하지 않았다.
현재 HEAD와 제품 후보의 제품 경로 diff는 비어 있다. QA·기획 커밋의 준비 SHA 포함도 확인했다.

reviewer `msg_4fbcac80f76c`의 `deliver:human` 인증 경계 high는 미해결로 유지한다.
이는 DEV 수정과 독립 재검증 대상이다. 기존 tester의 결함 없음 판정이나 이번 UI PASS로 해소하지 않는다.
coor는 high 해결 전 제품 수락·병합을 차단한다. 문서 결과의 통합은 제품 수락과 구분한다.
DEV 수정이 화면에 영향을 주면 coor가 같은 과제의 후속 검수 범위를 전달한다.
공유 산출물 인덱스의 D04 상태는 소유 범위 밖이므로 보존한다.
D04 자체의 정규 메타데이터·로컬 링크를 별도로 검사해 전체 strict의 미작성 집계 한계를 보완한다.
이번 작업에는 제품 코드 변경이 없다. 신규 기능·배포·기획 정책·PLANS·board 수정은 수행하지 않는다.

## 재개 문서 검증 결과

`deliverables.py --repo . --strict` 종료코드는 0이다. 검사 13, 미작성 4, 문제 0, 경고 0이다.
미작성 4는 공유 인덱스의 D04·D11·D12·D13이다. 전체 strict는 새 D04 자체를 검사하지 않는다.
별도 검사에서 D04의 `deliverables.problems(..., doc_id="D04")` 결과는 빈 목록이다.
변경 문서와 현재 인박스의 정규 front matter·로컬 링크도 통과했다.
기존 PNG 7개는 QA `c59537b`의 Git blob과 바이트가 일치한다.
`c59537b`와 `69dbec4`의 `661966f` 조상 검사 종료코드는 각각 0이다.
제품 경로의 `git diff --exit-code a6a10c7 -- ...`와 `git diff --check` 종료코드는 각각 0이다.
명령의 직접 종료코드와 PNG SHA-256은 `/tmp/SAR-MVP-001-UI-validation/ctx_f6be2c2d0ce4/`에 보존한다.
최종 FullOps lint는 깨끗한 커밋에 수행하고 해당 SHA와 결과를 완료 보고·worker_done으로 전달한다.

## 커밋된 문서의 lint와 환경 복구

문서 커밋은 `26e601d8d69394602945d9333acf8bd32a22c536`이다.
`lint.py --repo . --from 0a5b044` 첫 실행 종료코드는 1이다.
`product-lint`의 TypeScript 검사가 `@types/node` 부재로 실패했다. ERROR 1, WARNING 0, 실행 불가 0이다.
이는 tester의 기존 의존성 부재 관찰과 같은 환경 조건이다. 제품 소스나 잠금 파일을 변경하지 않았다.
잠금 파일 기반 `npm ci --prefix adapters` 종료코드는 0이다. 설치 대상은 Git 제외 `adapters/node_modules`다.
같은 커밋에서 lint를 다시 실행했다. 종료코드 0, `product-lint` passed, ERROR 0, WARNING 0, 실행 불가 0이다.
실패·복구 결과는 `lint-initial.json`·`install.log`·`lint-restored.json`에 보존한다.
인박스 아카이브 후 최종 커밋에도 동일 기준 lint를 실행한다. 최종 SHA의 결과는 worker_done으로 전달한다.
