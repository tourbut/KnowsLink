---
title: SAR-MVP-001-UI — 직접 시각 검수
status: draft
updated: 2026-10-03
owner: designer
tasks: [SAR-MVP-001-UI]
summary: 고정 합성 후보의 PNG 직접 관찰과 시각 판정 및 인증 high의 미해결 경계를 기록한다
---

# SAR-MVP-001-UI — 직접 시각 검수

## 판정과 기준

V-01–04의 캡처된 합성 조건은 PASS다. 발견한 미해결 critical/high 시각 결함은 없다.
이 판정은 직접 시각 검수의 완료다. 독립 코드 리뷰·병합·원격 공유·전체 MVP 수락·운영 공개를 대신하지 않는다.
제품 후보는 `a6a10c71977b7f3ec8274a1fb7c8a409f58e7c92`다.
독립 QA 증거 커밋은 `c59537b6fa0c7e008c4c6bdba0a251dd821d4ee8`이다.
준비 HEAD는 `661966f5646ecd44d108beddefed0630b65076c0`이다.
공개 정책 문서 `69dbec44c0193266f8f6c8499f22493e1e3c1722`도 준비 HEAD에 포함된다.
공통 규칙 `fullops-common-0.3.2`, [project.md](../../../project.md), [D02](../../planning/product-specs/SAR-MVP.md)의 MVP-10·MVP-16을 적용했다.
기존 DEV/TESTER 탐색 근거를 재사용했다. 제품 경로의 후보 대비 diff는 비어 있다.

## 실제 시각 도구와 증거

2026-10-03 designer는 아래 PNG 7개를 `view_image`로 직접 열었다.
HTML 내용이나 tester의 판정만으로 시각 수락을 선언하지 않았다.
모든 기존 PNG는 1280×900이다. 원래 tester 실행 SHA·조건을 유지한다.

| ID | PNG | 직접 관찰과 판정 |
|---|---|---|
| V-01 | [pending](SAR-MVP-001-TESTER-test/ui/v01-pending.png) | PASS. 발신·대상, 원요청 ID, `schedule.query`, 만료, 정책과 JSON이 분리돼 읽힌다. hint가 판단 자료를 가리지 않는다. |
| V-02 | [approved](SAR-MVP-001-TESTER-test/ui/v02-approved.png) | PASS. 제목의 `approved`와 하단 비활성 문구가 일치한다. 재승인 버튼이 없다. |
| V-02 | [denied](SAR-MVP-001-TESTER-test/ui/v02-denied.png) | PASS. `denied`와 비활성 문구가 일치한다. approved와 문자로 구분한다. |
| V-03 | [expired](SAR-MVP-001-TESTER-test/ui/v03-expired.png) | PASS. `expired`, 과거 만료 시각, 비활성 문구가 보인다. 활성 승인 버튼이 없다. |
| V-03 | [revoked](SAR-MVP-001-TESTER-test/ui/v03-revoked.png) | PASS. `revoked`, `원문 부재: 승인 불가`, 빈 Intent와 비활성 문구가 보인다. |
| V-04 | [unavailable](SAR-MVP-001-TESTER-test/ui/v04-unavailable.png) | PASS. `unavailable`와 원문 부재를 함께 표시한다. typed body를 임의 복원하지 않으며 승인 버튼이 없다. |
| V-04 | [unauthorized](SAR-MVP-001-TESTER-test/ui/v04-unauthorized.png) | PASS, 인증 실패 차단 범위. `invalid_auth`만 보이며 원문·정책·승인 버튼이 없다. 성공 상태로 표시하지 않는다. |

V-01의 기존 캡처는 화면 아래 두 버튼의 일부를 자른다.
이 캡처만으로 버튼의 전체 문구와 하단 배치를 확인할 수 없어 보완했다.
기존 [pending HTML](SAR-MVP-001-TESTER-test/ui/v01-pending.html)을 Chrome에서 1280×1100으로 렌더했다.
명령은 `google-chrome --headless=new --disable-gpu --no-sandbox --hide-scrollbars --window-size=1280,1100 --user-data-dir=<임시 경로> --screenshot=<보완 PNG> file://<기존 pending HTML>`이다.
프로세스 종료코드는 0이다. 기존 HTML이나 PNG를 수정하지 않았다.
[보완 PNG](../../design-docs/mockups/SAR-MVP-001-UI-v01-pending-full.png)도 `view_image`로 직접 확인했다.
`Approve 승인`과 `Deny 거절`의 전체 문구·버튼 경계·간격, 하단 `Owner 작업 화면` 링크가 보인다.
보완은 저장 HTML의 정지 렌더다. 새 서버 동작 QA나 새 후보 실행으로 기록하지 않는다.

## 제품 완료 조건 대조

검증된 typed body는 별도 제목 아래 JSON 전체를 표시한다. 정책은 그 위에 보인다.
승인 설명은 정보 공개나 일정 실행을 허용하지 않는다고 명시한다.
확인한 모든 화면은 문자 상태를 사용한다. 색에만 의존하지 않는다.
채팅 composer·버블·장기 timeline·hint 강조로 숨긴 본문은 없다.
[기존 QA-08](SAR-MVP-001-TESTER.md#수락-기준별-결과)은 GET pending 유지, POST approve 303, 재결정 409를 기록한다.
pending 저장 HTML의 POST form과도 일치한다. 화면의 탐색 링크를 GET 승인 링크로 판단하지 않는다.
중복 승인 불가의 서버 동작 증거는 기존 QA다. 이번 시각 검수에서 서버 동작을 재실행하지 않았다.

unauthorized는 인증 없는 GET의 HTTP 401 본문을 감싼 정지 캡처다.
이 화면은 인증된 owner에게 보이는 별도 `권한 불명` gate 화면이 아니다.
기존 완료 조건은 권한 확인 실패 시 활성 승인을 보이지 않는 것이다.
이번 인증 실패 캡처는 자료 미노출·승인 불가·성공 오인 방지 조건을 충족한다.
현재 기준에 없는 오류 안내 디자인이나 로그인 경로를 추가하지 않는다.
인증된 owner의 권한 저장소 장애·고의 stale epoch의 동작은 이 PNG로 검증하지 않았다.
그 미실행을 현재 권한 전체 검증 PASS로 확장하지 않는다.

## 범위 밖·held·인계

DEC-02, DEC-03, Free N, 실adapter, A2A 현행 검토, WAL/backup 삭제, 고의 stale epoch는 기존 held다.
시각 판정의 held만 이번 직접 검수로 해소했다. 공개 정책·수치·실데이터 조건은 확정하지 않았다.
모바일, 키보드, 스크린리더, 장시간 동작은 검사하지 않았다. 정지 화면으로 판단 가능하므로 영상은 만들지 않았다.
D04는 [화면 기록](../../design-docs/mockups/SAR-MVP-001-UI.md)에 작성했다.
제품 코드·기획 정책·PLANS·board·기존 QA 증거는 보존했다.
검증과 자기 기록은 [실행 기록](../../exec-plans/phases/SAR-MVP-001-UI.md)에 연결한다.
coor는 고정 SHA의 독립 리뷰와 필수 검사에 이 직접 검수 결과를 합쳐 병합·원격 공유를 처리한다.
사용자 지시에 따라 현재 과제의 필수 검수·병합 후 중지한다. 후속 기능·배포를 시작하지 않는다.

## 이번 재개 검수와 독립 리뷰 경계

2026-10-03 Dispatch `ctx_f6be2c2d0ce4`에서 PNG 8개를 `view_image`로 다시 직접 확인했다.
기존 관찰과 보완 PNG 생성 기록은 이전 Dispatch의 기록으로 유지한다.
이번 재개는 저장된 증거의 직접 시각 판정이며 새 서버 실행이나 새 캡처가 아니다.
캡처된 합성 조건의 V-01–04는 PASS다. `invalid_auth`는 인증 실패 차단 범위에 한정한다.
reviewer `msg_4fbcac80f76c`가 보고한 `deliver:human` 인증 경계 high는 미해결이다.
이 시각 PASS는 해당 high를 해소하지 않는다. DEV 수정과 독립 재검증 전 제품 수락·병합을 차단한다.
coor에게 문서 결과의 통합과 제품 수락의 구분을 인계한다. 현재 과제 검수·병합 후 중지 지시를 유지한다.
