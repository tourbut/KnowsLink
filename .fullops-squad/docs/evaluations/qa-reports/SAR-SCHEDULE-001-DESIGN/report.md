---
title: 실제 일정 API 연동 결정 근거
status: draft
updated: 2026-10-20
owner: designer
tasks: [SAR-SCHEDULE-001-DESIGN]
summary: 캘린더 연동 선택지의 외부 API 근거와 위험·비용 비교 및 추천 이유와 미확인 항목을 기록한다
---

# SAR-SCHEDULE-001-DESIGN 결정 근거

## 판정 범위

제품 범위 선택지와 규칙 초안의 문서 작업이다. 정본은 [SAR-SCHEDULE-001](../../../planning/product-specs/SAR-SCHEDULE-001.md)이다. 외부 캘린더 계정 연결·실데이터 조회·API 키 발급·Cloud Console 확인·제품 코드 수정을 하지 않았다. 새 동작 테스트는 실행하지 않았다. 구현 대상 동작이 없기 때문이다.

## 읽은 내부 근거

| 근거 | 확인한 내용 |
|---|---|
| [D02 MVP](../../../planning/product-specs/SAR-MVP.md) | MVP-08/11, DEC-02 보류, approve는 gate 통과만이며 `schedule.commit`은 stub |
| [백로그](../../../planning/SAR-MVP-backlog.md) | SAR-MVP-003 held, DEC-02 재개 조건(필드·수치·권한·누적 정책과 승인 근거) |
| 원천 [product](../../../planning/sources/silent-agent-relay/product.md)·[protocol](../../../planning/sources/silent-agent-relay/protocol.md) | C4 기본 비공개, `schedule.query` body `{window, granularity_min}`, 수치 TBD, 수치 발명 금지 |
| [사람 친구 D02](../../../planning/product-specs/SAR-HUMAN-PEER.md) | HP-13/14 해제·차단, HP-15 친구 수락은 대화 경로 동의 |
| [D03](../../../design-docs/architecture.md)·[interface-design](../../../design-docs/interface-design.md) | gate-consume 효과·공개 false, gate 요약 필드, 실제 calendar는 DEC-02 이후 |
| [U5](../../../design-docs/SAR-UX-ISSUES-002.md#u5--사람이-판단하는-승인-요약) | 승인 직전 고정 문구, hint 비권한, 요약은 검증된 구조화 필드만 |
| [Google 로그인](../../../exec-plans/phases/SAR-GOOGLE-LOGIN-001-DEV.md)·[연결 보완](../../../exec-plans/phases/SAR-GOOGLE-CONNECT-002-DEV.md) | 현재 OAuth Web client는 `openid email` scope만 쓴다. Testing이면 Test users 등록 필요 |

## 외부 API 근거

Context7은 monthly quota 초과로 쓰지 못했다. 하위 조사 에이전트(Purpose: research, 읽기 전용)가 2026-10-20에 공식 문서를 WebFetch로 직접 읽었다. WebFetch는 발췌 요약을 돌려주므로 인용은 원문과 다를 수 있다.

| ID | 사실 | 출처 |
|---|---|---|
| E1 | `freebusy.query` 입력은 `timeMin`·`timeMax`·`items[].id`다. 출력은 캘린더별 `busy[]` start/end와 errors뿐이며 제목 필드가 없다. | [freebusy.query](https://developers.google.com/workspace/calendar/api/v3/reference/freebusy/query) |
| E2 | 빈 시간 조회에 쓸 수 있는 최소 권한은 `calendar.freebusy`(View your availability)와 `calendar.events.freebusy`다. | [Calendar auth](https://developers.google.com/workspace/calendar/api/auth), E1 |
| E3 | 일정 생성은 `calendar.events`·`calendar.events.owned`·`calendar.app.created` 중 하나가 필요하다. `sendUpdates=none`이어도 일부 메일이 갈 수 있다고 문서가 경고한다. | [events.insert](https://developers.google.com/workspace/calendar/api/v3/reference/events/insert) |
| E4 | Calendar 권한은 restricted 목록에 없다. 캘린더 일정 읽기는 sensitive의 예다. sensitive 심사는 브랜드 검증·사유·데모 영상이 필요하며 보통 3–5 영업일이다. | [restricted 목록](https://support.google.com/cloud/answer/13464325), [sensitive 심사](https://developers.google.com/identity/protocols/oauth2/production-readiness/sensitive-scope-verification) |
| E5 | 미검증 앱이 sensitive 권한을 쓰면 경고 화면과 신규 사용자 누적 100명 상한이 있다. | [unverified apps](https://support.google.com/cloud/answer/7454865) |
| E6 | External·Testing 상태 앱의 refresh token은 7일 뒤 만료된다. refresh token은 사용자 철회·6개월 미사용 등으로도 무효화된다. | [OAuth 2.0 개요](https://developers.google.com/identity/protocols/oauth2) |
| E7 | `include_granted_scopes=true`로 증분 동의를 받는다. 사용자는 일부 권한만 허용할 수 있다. 앱은 `oauth2.googleapis.com/revoke`로 철회하며 사용자는 myaccount 권한 화면에서 철회한다. | [Web server OAuth](https://developers.google.com/identity/protocols/oauth2/web-server) |
| E8 | 쿼터는 프로젝트당 분당 10,000, 사용자당 분당 600 요청이다. 무료 상한은 하루 1,000,000 요청이며 2026년 이후 초과 과금 예정 문구가 있다. 한도 오류는 403/429이며 truncated exponential backoff를 권장한다. | [quota](https://developers.google.com/workspace/calendar/api/guides/quota), [errors](https://developers.google.com/workspace/calendar/api/guides/errors) |
| E9 | Microsoft getSchedule은 개인 Microsoft 계정을 지원하지 않는다. 최소 위임 권한은 Calendars.ReadBasic이다. 응답 예시에 subject·location이 있다. | [getSchedule](https://learn.microsoft.com/en-us/graph/api/calendar-getschedule), [permissions](https://learn.microsoft.com/en-us/graph/permissions-reference) |
| E10 | iCloud 서드파티 앱은 Apple 계정 승인 또는 앱 암호를 쓴다. 앱 암호는 2단계 인증이 필요하다. | [Apple 121539](https://support.apple.com/en-us/121539), [Apple 102654](https://support.apple.com/en-us/102654) |

### 미확인

1. `calendar.freebusy`·`calendar.events.freebusy`·`calendar.events.owned`의 sensitive 여부. Cloud Console Data Access 화면에서 확인해야 한다.
2. freebusy 한 번의 최대 조회 기간.
3. `sendUpdates` 생략 시 기본값과 일정 생성·초대 수 한도.
4. 현재 운영 OAuth 앱의 게시 상태(Testing/In production)와 user type.
5. iCloud 공식 개발자 API 경로, Microsoft 테넌트 동의 제한의 영향.
6. Grok·Dots 등 연결 agent의 자체 캘린더 도구 지원(SD-02 B의 전제).

## 선택지 비교와 추천 이유

| 결정 | 추천 | 위험·비용·운영 영향 요약 |
|---|---|---|
| SD-01 캘린더 | Google만 | 로그인과 같은 제공자다(E7). 응답에 제목이 없다(E1). Outlook은 개인 계정 미지원·제목 포함 응답(E9)으로 비용 대비 이득이 작다. iCloud는 앱 암호 수집 위험(E10). |
| SD-02 실행 위치 | 서버 | 서버가 공개 필드를 강제한다(C4). 비용은 토큰 암호화 보관·유출 영향이다. agent 측 실행은 각 agent 도구 지원이 미확인이다(미확인 6). |
| SD-03 기능 | 조회만 | 권한 하나(E2)로 끝난다. 생성은 쓰기 권한·중복·취소 규칙이 필요하다(E3). 초대는 `none`이어도 메일이 갈 수 있어 HP-15 외부 발송 위험이 있다(E3). |
| SD-04 공개 | 바쁨 구간 반올림 | 제공자 응답 자체가 바쁨 구간뿐이다(E1). 제목 공개는 상세 읽기 권한(sensitive, E4)과 개인정보 위험을 늘린다. |
| SD-05 승인 | 매 요청 승인 1회 | 실패·오남용 영향을 요청 하나로 묶는다. 쿼터(E8) 소비가 사람 결정 속도로 제한된다. 친구별 자동 허용은 수치·화면 결정이 더 필요해 후속으로 미룬다. |
| SD-06 연결 | 별도 증분 동의 | 로그인 동의 범위를 유지한다(E7). 권한 일부 거부를 화면이 처리해야 한다. 철회 경로가 두 곳(KnowsLink·Google)이다(E7). |

공통 운영 영향: 출시 전 OAuth 앱을 In production으로 두어야 한다(E6). 권한이 sensitive이면 심사 3–5 영업일과 심사 전 100명 상한이 있다(E4, E5). 토큰 보관·철회·감사 기록 규칙은 정본의 '기록과 보존'에 있다. 재시도는 승인한 요청 하나의 기한 안으로 묶는다(E8).

## 소유권과 보존

새 파일은 정본과 이 보고 두 개다. D01·D02 인덱스에 연결 한 줄씩을 추가했다. 다른 역할 인박스·PLANS·board·제품 코드·미추적 사용자 자료는 수정하지 않았다. FullOps lint 결과는 커밋 뒤 완료 보고에 기록한다.

## 사용자 결정 기록 (2026-10-20)

coor가 사용자에게 물은 답을 두 번째 attempt에서 정본에 기록했다.

- D009: SD-01–06 추천안(A)을 모두 채택했다. 위 비교의 추천 이유가 그대로 확정 근거다.
- D010: 수치 제안(최대 14일, 15·30·60분과 기본 30분, 친구당 24시간 5회, 기존 gate 기한)을 승인했다.
- D011(열림): OAuth 앱 게시 상태와 `calendar.freebusy` 심사 여부는 미확인 4번과 같은 항목이다. 사용자가 나중에 확인한다. 개발과 격리 QA는 진행하고 운영 연결 전에 다시 묻는다. Testing이면 7일 만료(E6), In production이면 심사 가능(E4, E5)이다.

이번 attempt도 외부 계정·실데이터·제품 코드를 건드리지 않았고 새 동작 테스트는 없다. 정본 status는 `approved`이며 근거는 D009·D010이다.
