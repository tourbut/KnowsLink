---
title: 재초대·기록 포화 제품 답의 구현 대조와 회원 화면 안내 수정 기록
status: draft
updated: 2026-10-06
owner: dev
tasks: [SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX]
summary: POLICY 48d12fa 관찰 조건과 4a 구현을 대조하고 반복 초대·수동 재초대·기록 포화 복구·24h 목록 안내와 회귀 검사를 기록한다
---

# SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX — 재초대·기록 포화 제품 답의 구현 반영

## 기준과 기술 계획

기준 ref는 `4a1b80aec8fa6a06144d51f3a5609927a2644928`이다. 준비 HEAD는 `0a83bbb`다. 제품 답은 `48d12fae2dce35d92606b264313148f0a635b64e`의 [POLICY](SAR-PUBLIC-AGENTS-001-POLICY.md) 원문과 두 관찰 조건표다. 4a와 48d12fa가 준비 HEAD의 조상이며, 준비 HEAD의 제품 코드가 4a와 같음을 `git diff 4a1b80a HEAD -- internal cmd adapters`(출력 0)로 확인했다. 적용 기준은 fullops-common-0.3.3·FULLOPS·project·document-writing·D02 PS-07과 기록 보호 조건·UX-04/05다. 예외는 없다.

계획: 관찰 조건을 4a 코드에 한 행씩 대조한다. 이미 일치하는 동작은 재구현하지 않고 결정적 검사로 고정한다. 차이는 회원 화면의 상태·다음 동작 안내에만 있었다. 기존 Go template·memberStyle·`refusal`·notice 구조를 재사용한다. 새 라이브러리·migration·wire 필드·상태 필드는 없다. 예상 규모는 Go 3파일과 테스트 2파일, 기술 문서 5개였다. 실제 제품 변경은 5파일 +428/−33줄이다(테스트 +328/−10). Context7 조회는 새 라이브러리가 없어 해당 없음이다.

## 관찰 조건과 4a 구현의 대조

| 조건 | 4a 구현 | 이번 처리 |
|---|---|---|
| 유효 pending 같은 방향 반복 | `invites`가 pending·active면 기존 Pair를 그대로 반환한다. 회원 rate는 거부 요청까지 소비한다 | 일치. 반복 결과 안내(`invite-pending`)만 추가 |
| 유효 pending 반대 방향 | pairID가 대칭이라 같은 Pair를 반환한다. 수신자·기한 불변, 결정은 Recipient owner만 | 일치. 안내가 받은 초대의 결정 위치를 알린다 |
| active 반복 | 기존 active 반환. 세대·수 불변 | 일치. `invite-active` 안내 추가 |
| 거절 뒤 | denied는 pending 한도에서 빠진다. 수락은 403. 재초대는 `Generation+1` pending | 일치. 홈에 종료 안내와 `새 초대 보내기` 추가 |
| 24h 만료 경계 | transaction이 먼저 sweep한다. `now >= Exp`이면 expired | 일치. 기한 직전 수락·기한 도달 거부 검사 추가 |
| 양측 철회 | unpair는 양쪽 owner 허용. 옛 메시지는 `current()`의 세대 비교로 거부 | 일치. 양방향 각각 검사 추가 |
| 새 수동 초대·재수락 | 한도·agent/owner 재검사, 새 세대, active 한 관계 | 일치 |
| 동시·늦은 결정 | singleton lock에서 한 active. 회원 화면은 Generation 일치를 검사 | 일치. 기존 integration 유지 |
| 동일 owner·독립 owner | 자동 수락 없음. 타 owner 결정·unpair 403 | 일치 |
| 한도 경계 | 송신 pending 10·전체 200·active 20/400. 상한이면 세대 증가 전 409 | 일치. 상한 재초대가 세대를 만들지 않음을 검사 |
| 재시작·옛 자격 | 상태는 Postgres JSON. 철회 agent·비활성 owner 재초대 거부 | 일치 |
| 키 기록 보호 상한 | connect·complete·합성 keys 409. 기존 자격 유지, 철회 허용 | 의미 일치. 409 안내가 일반 문구여서 교체 안내로 변경 |
| 포화 agent의 키 철회·대기 | 살아 있는 agent의 철회 키는 지우지 않는다 | 일치. 안내에 공간이 생기지 않음을 명시 |
| 교체할 새 agent | 새 무작위 ID·별도 연결·관계 승계 없음 | 일치. 거부 화면에 생성 가능 시 `새 agent 만들기` 버튼 |
| 활성 한도·owner 기록 상한 | 둘 다 409 `capacity`, 같은 일반 문구 | 한도별 안내로 분리(`agentLimit`) |
| 포화 중 정리 | 정리 budget, 기록 증가 없음 | 일치 |
| 철회 직후·24h 전·정리 뒤 | 철회 즉시 거부, 24h 뒤 agent·키·pair 삭제 | 일치. 홈 안내 추가, 보존 중·정리 뒤 옛 credential 거부 검사 추가 |
| 오류·정리 대기 화면 | 일반 운영 한도 문구와 홈 링크 | 원인·다음 동작·관리 복귀 안내로 변경 |

## 변경

- 초대 제출 결과: 기존 pending 유지, active 유지, 새 pending을 notice로 구분한다. 반복은 새 초대·기한 변경이 없음을 알린다.
- 관계 목록: 받은/보낸 pending의 다음 동작을 표시한다. 거절·만료·철회 관계는 연결 종료·메시지 불가·새 초대와 상대의 새 수락 필요를 표시한다. 자기 살아 있는 agent가 있으면 같은 상대에 대한 수동 `새 초대 보내기` 버튼을 둔다. 종료 관계에는 의미 없는 관계 철회 버튼을 표시하지 않는다.
- 키 기록 포화: 철회·대기로 공간이 생기지 않음, 새 agent 별도 연결, 각 상대와 새 초대·수락을 안내한다. 새 agent를 만들 수 있으면 생성 버튼을, 없으면 그 실제 한도를 덧붙인다. 홈 카드도 연결 수단 발급 대신 같은 안내를 표시한다.
- owner 기록 포화: 필요한 철회 선택, 철회 기록의 최소 24시간 보존, 정리되어 홈에서 사라진 뒤 수동 재시도를 안내한다. 활성 5개는 철회 선택과 복구 불가를 안내한다. 활성 키 3개 등록은 키 철회 또는 회전을 안내한다.
- 철회 agent 카드: 최소 24시간 보존 뒤 목록에서 사라질 수 있음, 권한 복구·백업 영구 삭제가 아님, 새 연결과 새 수락 필요를 표시한다. agent 철회 버튼 앞에 모든 키·관계 종료와 복구 불가를 표시한다.
- 상품 quota·결제·정확한 정리 시각·키 공간 생성은 안내하지 않는다. `/v1/*` 응답은 바꾸지 않았다.

## 자동 검증 증거

레포 루트에서 실행했다. 명령의 원래 종료코드를 저장했고 파이프로 가리지 않았다. 로그는 [QA 증거](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-DEV-POLICY-FIX/)이며 대상 SHA는 `head.txt`의 제품 커밋 `f9af9bf`다. 로그 후행 공백만 정규화했다. 4a와 원본 d1 검사 기록은 다시 쓰지 않았다.

| 명령 | 결과·대상 |
|---|---|
| make lint | exit 0 |
| make test | exit 0. Go race 전체와 adapter 검사 |
| make verify-mvp | exit 0. 고유 Compose project의 격리 Postgres·migration·Go integration race(PASS 55·FAIL 0)·TS 왕복. 그 project의 relay는 integration 동안 정지한다. 기존 `knowslink-*` 컨테이너·Tunnel은 건드리지 않았다 |
| 390×844 폭 측정 | `mobile-width.json`. 새 안내 상태의 홈과 두 거부 화면이 scrollWidth 390, 넘침 요소 0 |

새 검사와 검출력:

- `TestRelationshipPolicy`: POLICY 첫 표의 State 단위 행렬. 반복 시 기한 연장·재초대 세대 재사용 두 변형을 임시로 넣어 각각 실패함을 확인한 뒤 되돌렸다.
- `TestSaturationGuidance`: 두 번째 표의 안내·생성 가능 여부·24h 전후·옛 credential 거부.
- `TestMemberPagesShowNextSteps` 확장: 관계 상태별 안내·수동 새 초대·종료 관계의 철회 버튼 부재(변형 시 실패 확인)·철회 agent 보존 안내·키 포화 카드·거부 화면의 생성 버튼.
- `TestPublicAgentHTTP` 두 하위 검사 확장: 실제 Postgres HTTP의 notice redirect 3종과 반대 방향 반복의 기한 불변, 종료 관계 화면, owner 기록 포화 409 안내, 키 기록 포화 connect 409 안내·생성 버튼·연결 미생성.

FullOps lint·strict·`git diff --check`·고정 SHA 결과는 아래 완료 기록에 남긴다.

## 기술 판단과 미변경

- `/v1/invite-decision`은 Generation을 받지 않는다. 같은 방향의 새 pending에 늦게 도착한 owner API 수락은 현재 수신 owner의 새 결정으로 처리된다. 회원 owner는 bearer가 없고 공개 후보는 합성 가입을 끄므로 공개 경로는 Generation을 검사하는 회원 화면뿐이다. wire 필드 추가는 범위 밖이라 바꾸지 않았다. OPS delta 리뷰의 확인 대상이다.
- 관계 화면의 기존 `관계 세대` 표시는 유지했다. 새 기술값 노출은 추가하지 않았다.
- 차단·쿨다운·수신 pending 한도·새 수치는 추가하지 않았다.

## QA·UI·운영 인계

- OPS: 최종 fixed SHA의 독립 delta 리뷰. `agentLimit` 분리가 한도 집행을 바꾸지 않았는지, 회원 거부 안내가 정보 노출을 늘리지 않았는지, 위 `/v1/invite-decision` 판단을 확인한다.
- TESTER: POLICY 두 표의 좁은 독립 QA. 반복 초대 notice와 상태 불변, 기한 경계, 양측 철회 뒤 재초대·새 수락, 키/owner 기록 포화 안내와 24h 정리 전후, 기존 보존·rate 회귀.
- designer: F-UI-01–04 영향과 새 안내(반복 초대 notice, 종료 관계의 새 초대, 키/owner 기록 포화 거부와 생성 버튼, 철회 agent 보존 안내, 철회 전 경고)를 실제 화면에서 직접 재검수한다. DEV 폭 측정은 시각 PASS가 아니다.
- D12/D13(OPS 소유, 구현 기술 인계·미검증): 이번 변경은 HTML 안내만이며 상태 형식·보존 규칙·운영 절차를 바꾸지 않는다. DEV-FIX의 D12 경계(24h 삭제·rollback 시 보존 재시작·합성 가입 금지)는 그대로다. 운영 수락·실제 자원 측정은 하지 않았다.

## 미검증·범위

실메일·공개·운영 배포·외부 계정·운영 데이터 정리·과금·실제 24h 경과·부하는 실행하지 않았다. 원본 d1 UI FAIL/보류·리뷰·4a 검증 기록은 바꾸지 않았다. 독립 QA·직접 시각 판정·운영 수락은 각 담당 후속이다.
