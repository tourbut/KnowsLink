---
id: D04
title: 화면설계서
status: review
updated: 2026-10-06
owner: designer
tasks: [SAR-MVP-001-UI, SAR-PUBLIC-AGENTS-001-UI, SAR-PUBLIC-AGENTS-001-UI-FIX, SAR-PUBLIC-MESSAGES-001-UI]
upstream: [D02]
summary: 합성 승인과 일반 회원 연결 화면의 기록 및 직접 시각 수락 경계를 보존한다
---

# SAR-MVP-001-UI — 합성 human-gate 화면 기록

## 일반회원 메시지·gate의 새 직접 기록

고정 `09c523da8a3407288d9f5d711e1834af12bc7808`의 일반회원 UX06/07은 [MESSAGES 직접 검수](SAR-PUBLIC-MESSAGES-001-UI.md)에 있다.
UX06은 좁은 로컬 시각 조건에서 PASS다. UX07은 Deny 뒤405의 F-UI-MSG-01 medium으로 FAIL/보류다.
기존 합성 gate·AGENTS 검수의 원래 SHA와 판정은 유지한다. 이 연결은 제품/UX 규칙을 바꾸지 않는다.

## 기준과 범위

이 D04는 구현된 합성 후보의 화면을 기록한다. 새로운 화면이나 제품 정책을 제안하지 않는다.
기준은 [D02](../../planning/product-specs/SAR-MVP.md)의 MVP-10·MVP-16과 UI 제품 방향이다.
제품 후보는 `a6a10c71977b7f3ec8274a1fb7c8a409f58e7c92`다.
원래 캡처와 독립 QA 기록은 `c59537b6fa0c7e008c4c6bdba0a251dd821d4ee8`에 포함된다.
검수 준비는 `661966f5646ecd44d108beddefed0630b65076c0`이다.
이번 D04의 review는 화면 기록의 문서 상태다. 운영 공개나 전체 MVP 승인을 뜻하지 않는다.

## 화면 구조

화면 제목은 `KnowsLink 요청 승인`이다. 밝은 배경 위의 단일 흰색 패널에 승인 판단 자료를 세로로 배치한다.
패널은 상태, 발신·대상 에이전트, 원요청 ID, Intent, 만료, 적용 정책, 검증된 typed body, 승인 의미, 결정 영역 순서다.
정책은 `schedule.query: disclosure policy absent; deny even after gate approve`로 표시된다.
typed body는 별도 제목과 JSON 블록으로 표시된다. `granularity_min`과 `window.start/end`를 읽을 수 있다.
Approve가 인간 게이트만 통과시키며 정보 공개나 일정 실행을 허용하지 않는다는 설명을 표시한다.
상태는 문자로 표시한다. 색만으로 상태를 구분하지 않는다.
검수한 화면에는 채팅 버블, 입력 composer, 장기 대화 timeline이 없다.
render.hint는 이 합성 화면에 표시되지 않는다. hint가 typed body를 가리지 않는다.

## 화면별 표시와 행동

| 항목 | 표시 | 결정 영역 | 시각 근거 |
|---|---|---|---|
| V-01 pending | 발신·대상, Intent, 만료, 정책, typed body | `Approve 승인`, `Deny 거절` 버튼 | [기존 pending](../../evaluations/qa-reports/SAR-MVP-001-TESTER-test/ui/v01-pending.png), [전체 pending](SAR-MVP-001-UI-v01-pending-full.png) |
| V-02 approved | `상태: approved`, 같은 판단 자료 | `승인·거절 버튼 비활성: approved` | [approved](../../evaluations/qa-reports/SAR-MVP-001-TESTER-test/ui/v02-approved.png) |
| V-02 denied | `상태: denied`, 같은 판단 자료 | `승인·거절 버튼 비활성: denied` | [denied](../../evaluations/qa-reports/SAR-MVP-001-TESTER-test/ui/v02-denied.png) |
| V-03 expired | `상태: expired`, 과거 만료 시각 | `승인·거절 버튼 비활성: expired` | [expired](../../evaluations/qa-reports/SAR-MVP-001-TESTER-test/ui/v03-expired.png) |
| V-03 revoked | `상태: revoked`, 빈 Intent, `원문 부재: 승인 불가` | `승인·거절 버튼 비활성: revoked` | [revoked](../../evaluations/qa-reports/SAR-MVP-001-TESTER-test/ui/v03-revoked.png) |
| V-04 unavailable | `상태: unavailable`, 빈 Intent, `원문 부재: 승인 불가` | `승인·거절 버튼 비활성: unavailable` | [unavailable](../../evaluations/qa-reports/SAR-MVP-001-TESTER-test/ui/v04-unavailable.png) |
| V-04 인증 실패 | 흰 화면에 `invalid_auth` | 요청 자료와 결정 버튼 없음 | [unauthorized](../../evaluations/qa-reports/SAR-MVP-001-TESTER-test/ui/v04-unauthorized.png) |

pending의 두 버튼은 POST form에 속한다. GET 승인 링크가 아니다.
`Owner 작업 화면`은 탐색 링크다. 승인 결정 링크가 아니다.
결정 완료·만료·철회·원문 부재 상태에는 승인 버튼이 없다.
POST 방식과 중복 결정 차단은 이미지와 별개로 기존 QA-08의 실행 결과를 연결한다.
이미지 자체로 서버의 원자성이나 CSRF 방어를 증명하지 않는다.

## 시각 수락과 한계

[직접 시각 판정](../../evaluations/qa-reports/SAR-MVP-001-UI.md)의 V-01–04는 캡처된 합성 조건에서 PASS다.
designer는 기존 PNG 7개와 보완 PNG 1개를 `view_image`로 직접 확인했다.
기존 1280×900 pending은 아래 버튼 일부가 잘린다. 저장 HTML을 1280×1100으로 다시 렌더해 두 버튼을 확인했다.
보완 PNG는 기존 HTML의 정지 렌더다. 현재 서버에서 새 요청을 실행한 독립 QA가 아니다.
영상·모바일·키보드·스크린리더 검사는 이번 범위에 없다.
`invalid_auth`는 인증 없는 GET의 401 본문이다. 인증된 owner의 권한 저장소 장애 화면은 이 캡처로 검증하지 않는다.
현재 기준은 권한 확인 실패 시 승인 불가다. 별도 오류 안내 문구나 로그인 화면을 새 수락 조건으로 추가하지 않는다.

## 인계와 기존 기록 보존

이 원천은 검수 과제의 명시적 D04 작성 지시에 따라 추가했다.
D02의 준비 단계 `D04 미작성` 문장과 coor 소유 산출물 인덱스는 이번 소유 범위 밖이므로 보존했다.
coor는 병합 시 D04 원천이 review 상태로 추가됐다는 사실을 기록과 대조한다.
기획 정책·원천 잠금·PLANS·board·기존 QA 증거를 수정하지 않았다.
DEC-02·DEC-03·Free N·실adapter·A2A 현행 검토·WAL/backup 삭제·고의 stale epoch의 held를 유지한다.
후속 기능이나 배포는 시작하지 않는다. 최종 검토와 병합 후 중지는 coor에게 인계한다.

## 이번 재개 검수와 독립 리뷰 경계

2026-10-03 Dispatch `ctx_f6be2c2d0ce4`에서 PNG 8개를 `view_image`로 다시 직접 확인했다.
기존 관찰과 보완 PNG 생성 기록은 이전 Dispatch의 기록으로 유지한다.
이번 재개는 저장된 증거의 직접 시각 판정이며 새 서버 실행이나 새 캡처가 아니다.
캡처된 합성 조건의 V-01–04는 PASS다. `invalid_auth`는 인증 실패 차단 범위에 한정한다.
reviewer `msg_4fbcac80f76c`가 보고한 `deliver:human` 인증 경계 high는 미해결이다.
이 시각 PASS는 해당 high를 해소하지 않는다. DEV 수정과 독립 재검증 전 제품 수락·병합을 차단한다.
coor에게 문서 결과의 통합과 제품 수락의 구분을 인계한다. 현재 과제 검수·병합 후 중지 지시를 유지한다.

## SAR-PUBLIC-AGENTS-001-UI의 직접 검수 결과

2026-10-06 제품 고정 `d1eef9bb90b9726149980320c42fb1fdbcaf584a`의 일반 회원 UX04–05를 실제 격리 localhost에서 직접 검수했다. [과제별 관측과 판정](SAR-PUBLIC-AGENTS-001-UI.md)이 이번 결과 정본이다. 연결·키·관계 상태는 확인했으며 모바일 지문과 오류 뒤 복귀의 medium 두 건으로 전체 시각 수락은 FAIL/보류다. 기존 합성 gate와 제품 규칙 및 운영 공개의 보류 조건은 유지한다.

## SAR-PUBLIC-AGENTS-001-UI-FIX의 직접 재검수

고정 `458798c2ee15c179edacfd6f94ebb9896d26f411`의 F-UI-01–04와 POLICY 새 안내는 [직접 재검수 보고](SAR-PUBLIC-AGENTS-001-UI-FIX.md)의 좁은 시각 조건에서 PASS다. 원본 d1/d165 UI FAIL과 원본 QA/OPS 실행 기록은 그대로 보존했다. 독립 TESTER/OPS·main 통합·일반 서비스 공개 수락은 별도다. 제품 규칙과 디자인 방향은 바꾸지 않았다.
