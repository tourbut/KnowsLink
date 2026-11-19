---
title: SAR-SCHEDULE-001-VISUAL 직접 시각 검수
status: review
updated: 2026-10-21
owner: designer
tasks: [SAR-SCHEDULE-001-VISUAL]
summary: 일정 조회 설정·승인 화면 KO/EN 3폭 직접 시각 판정과 이메일 회원 안내·해제 시 미전달 결과 처리의 확정 내용
---

# SAR-SCHEDULE-001-VISUAL 직접 시각 검수

## 대상과 결론

- 대상 제품 코드: main `18d74dd72132311c9abb1e2aa6fdf6eb062d38e9`. 검수 체크아웃 HEAD `27da4a28`은 `18d74dd..27da4a28` 사이에 `.fullops-squad/` 밖 변경이 없다(`git diff --name-only`로 확인). 따라서 제품 코드는 18d74dd와 같다.
- 결론: 넘침·대비·영어판 잔존 한국어는 문제가 없다. 수정 필요 4건(모두 low, 화면 문구·배치)과 결정 반영 1건이 있다. 아래 표의 PASS 화면만 시각 수락한다. 전체 제품 수락이 아니다.
- 판단 1(이메일 회원 안내): 현행 문구와 버튼 없음을 수락한다.
- 판단 2(해제 때 미전달 결과): 지운다. 상대에게 공개하지 않는다. 해제 안내 문구를 바꾼다.
- 결정의 정본: [SAR-SCHEDULE-001](../../../planning/product-specs/SAR-SCHEDULE-001.md) '화면 문구 원칙'·'실패·만료·철회'·SC-06/07, [D04 캘린더 절](../../../design-docs/mockups/service-member-ui.md#캘린더-연결과-일정-조회-승인--sar-schedule-001).

## 방법

1. 체크아웃 HEAD를 `git archive`로 scratchpad에 풀었다. 제품 워크트리의 코드는 바꾸지 않았다.
2. `KNOWSLINK_UI_FIXTURE_DIR=<scratch>/fx go test ./internal/relay -run 'TestHumanAccountCalendarStates|TestMemberUIFixtures|TestEnglishScreensHaveNoKorean|TestZZVisualScratch' -count=1`을 실행했다. 결과는 `ok`다(Go 1.27.1, Windows).
3. 기존 검사에 없는 상태는 scratch 전용 테스트 `TestZZVisualScratch`로 같은 template과 같은 fixture에서 렌더했다. 대상은 영어판 email-member·not-ready와 해제 뒤 안내(KO/EN)다. 이 파일은 커밋하지 않았다(`raw/tools/`).
4. DEV의 `probe.mjs`를 복사해 headless Chrome(비모바일 viewport)으로 1440·390·320을 렌더했다. 18화면×3폭 = 54장이다. 캘린더 절만 잘라낸 캡처와 해제 details를 연 캡처를 따로 만들었다.
5. 모든 캡처를 KO/EN 행, 3폭 열의 비교판으로 묶어 직접 열어 보았다.

색 대비는 `:root` 토큰 값으로 계산했다. 흰 글자/accent 버튼 7.96:1, accent 링크/흰 배경 7.96:1, danger 글자/흰 배경 8.01:1, muted/흰 배경 6.19:1, 버튼 테두리 edge/흰 배경 3.79:1이다. 모두 AA(본문 4.5:1, 비텍스트 3:1)를 넘는다.

## 화면별 판정

넘침은 54장 모두 0이다(`scrollWidth ≤ innerWidth`). 영어판 6화면에 한국어가 없다(`TestEnglishScreensHaveNoKorean` 통과, scratch 영어판 3화면 grep 0). 피할 항목(크림 배경·알약 버튼·모노스페이스 라벨·장식 그라디언트·이모지 아이콘)은 없다. 버튼 radius는 4px이고 모노스페이스는 기술 상세의 ID에만 쓴다.

| 화면 | KO | EN | 판정 |
|---|---|---|---|
| 설정 캘린더 절 none-google | 3폭 | 3폭 | PASS |
| 설정 캘린더 절 connected | 3폭 | 3폭 | 배치 PASS. 해제 안내 문구는 판단 2로 바뀐다(D2) |
| 설정 캘린더 절 reconnect | 3폭 | 3폭 | 수정 필요 V3 |
| 설정 캘린더 절 email-member | 3폭 | 3폭 | PASS. 판단 1로 수락 |
| 설정 캘린더 절 not-ready | 3폭 | 3폭 | PASS |
| 해제 뒤 안내(status) | 3폭 | 3폭 | 배치 PASS. 문구는 판단 2로 바뀐다(D2) |
| gate 조회 connected | 3폭 | 3폭 | PASS. EN 320은 V2 |
| gate 조회 unconnected | 3폭 | 3폭 | 수정 필요 V1. EN 320은 V2 |
| gate 조회 approved | 3폭 | 3폭 | PASS. EN 320은 V2 |

fixture의 `ascii-text`·`확인할 수 없는 에이전트` 이름은 검사용 자리표시다. 결함이 아니다.

## 수정 필요 항목

DEV가 바로 고칠 수 있도록 현재와 기대를 적는다. 문구는 정본(SAR-SCHEDULE-001 '화면 문구 원칙')과 같다.

### V1 — low — 미연결 회원 조회 상세의 '승인하면 허용되는 범위'가 모순된다

- 화면·폭: gate-query-unconnected, KO/EN, 1440·390·320 모두. 캡처 [KO](v1-gate-query-unconnected-ko.png), [EN](v1-gate-query-unconnected-en.png).
- 현재: 승인 버튼이 없는데 '승인은 이 요청의 인간 게이트만 통과시킵니다. 정보 공개나 일정 실행을 허용하지 않습니다.'와 '가능 시간 공개 정책이 없어 승인해도 상대에게 가능 시간이 공개되지 않습니다.'를 보인다. 같은 화면 요약은 '공개되는 내용: 바쁜 시간대만'이다. 캘린더 기능이 생겨 '정책이 없어'도 사실이 아니다.
- 기대: 캘린더 미연결 회원(이메일 회원 포함)의 `schedule.query`에서 위 두 줄 대신 한 줄을 보인다.
  - KO: '이 요청은 승인할 수 없습니다. 거절하거나 기한이 지나면 조회 없이 끝납니다.'
  - EN: 'This request cannot be approved. If you decline it or it expires, it ends without a lookup.'
- 유지: '친구 동의나 대화 메시지는 …' 줄, 오류 alert '캘린더가 연결되지 않아 이 요청에 답할 수 없습니다.', '캘린더 연결' 링크, 거절 버튼. `schedule.commit`의 U5 문구는 바꾸지 않는다.

### V2 — low — EN 320 '참고 문구' 접기의 표시 기호가 혼자 줄을 차지한다

- 화면·폭: gate-query-connected/unconnected/approved, EN, 320. 캡처 [EN 320](v2-en-gate-query-connected-320.png). KO 320과 EN 390·1440은 정상이다.
- 현재: 'Reference note — not grounds for approval'의 ▸ 기호가 첫 줄에 혼자 있고 문구가 다음 줄에서 시작한다.
- 기대: 기호와 문구가 같은 줄에서 시작하고, 넘치는 문구만 다음 줄로 이어진다. 기존 공용 `details > summary` 스타일이므로 다른 gate 화면(예: commit hint)도 같이 확인한다.

### V3 — low — reconnect 상태의 해제 버튼에 결과 안내가 없다

- 화면·폭: 설정 캘린더 절 reconnect, KO/EN, 3폭. 캡처 [KO/EN 3폭](v3-calendar-reconnect.png).
- 현재: '캘린더 연결 해제' 버튼이 결과 설명 없이 바로 보인다. connected 상태는 위험 details 안에 결과 문장을 먼저 보인다.
- 기대: reconnect 상태도 connected와 같은 위험 details('캘린더 연결 해제' summary → D2의 해제 안내 → 해제 버튼)를 쓴다. '캘린더 다시 연결' 주 버튼은 그대로 둔다. D04의 '관리 details 안에서 결과 문장을 먼저 보여준다' 원칙과 같다.

### D2 — 결정 반영 — 해제 안내와 해제 뒤 안내 문구

- 화면·폭: 설정 캘린더 절 connected(details 열림), 해제 뒤 안내. KO/EN, 3폭. 현재 캡처 [KO/EN 3폭](d2-calendar-disconnect-copy.png).
- 현재 해제 안내: '해제하면 KnowsLink에 보관한 토큰을 즉시 폐기하고 Google에도 권한 철회를 요청합니다. 대기 중인 조회 요청은 거절로 끝납니다.'
- 기대 해제 안내:
  - KO: '해제하면 KnowsLink에 보관한 토큰을 즉시 폐기하고 Google에도 권한 철회를 요청합니다. 아직 상대에게 전달되지 않은 조회 요청과 결과는 거절로 끝나고 공유되지 않습니다. 이미 전달된 결과는 되돌릴 수 없습니다.'
  - EN: 'Disconnecting discards the token KnowsLink stores at once and asks Google to revoke access. Availability requests and results not yet delivered end as declined and are not shared. Results already delivered cannot be taken back.'
- 현재 해제 뒤 안내: '캘린더 연결을 해제했습니다. KnowsLink에 보관한 토큰을 폐기했고 대기 중이던 조회 요청은 거절로 끝났습니다.'
- 기대 해제 뒤 안내:
  - KO: '캘린더 연결을 해제했습니다. KnowsLink에 보관한 토큰을 폐기했고 아직 전달되지 않은 조회 요청과 결과는 거절로 끝났습니다.'
  - EN: 'Calendar disconnected. The token KnowsLink stored was discarded, and availability requests and results not yet delivered ended as declined.'
- 서버 동작(리뷰 L2): 해제 transaction에서 승인·조회가 끝났지만 아직 consume되지 않은 gate의 disclosure를 지운다. 요청자 agent는 대기 요청과 같은 `denied`를 받는다. 해제 사실을 따로 알리지 않는다. 이미 consume된 결과는 회수하지 않는다. 상태 전이와 구현 방법은 DEV가 정한다.

## 판단 1 — 이메일 회원 안내

현행을 수락한다. 이메일 코드로 로그인한 회원(Google 신원 없음)은 설정 '캘린더 연결' 절에서 공통 안내와 아래 한 줄만 본다. 연결 버튼은 없다.

- KO: 'Google 계정으로 로그인한 회원만 캘린더를 연결할 수 있습니다.'
- EN: 'Only members signed in with a Google account can connect a calendar.'

근거는 다음과 같다.

- SD-06 C(로그인 계정과 다른 Google 계정 연결)는 D009에서 선택하지 않았다. 이메일 계정에 Google 신원을 붙이는 계정 연결도 범위 밖이다.
- 버튼을 두면 Google 동의 뒤 같은 신원 확인에서 실패하는 막다른 길이 된다. 'Google로 로그인하세요' 안내는 별도의 새 계정을 만들게 할 수 있어 쓰지 않는다.
- 조회 상세의 '캘린더 연결' 링크는 이메일 회원에게도 이 설명 절로 이동한다. 절 제목이 '캘린더 연결'이므로 링크 이름을 바꾸지 않는다.

## 판단 2 — 해제 때 미전달 결과

지운다. 상대에게 공개하지 않는다. 문구와 서버 동작은 위 D2다.

근거는 다음과 같다.

- 해제는 회원의 동의 철회다. 사용자는 해제 뒤 더는 공유되지 않는다고 기대한다. 화면도 지금 '거절로 끝납니다'라고 말한다.
- 친구 해제·차단(HP-13/14)은 이미 승인된 요청의 결과 공개를 막고 disclosure를 지운다. 캘린더 해제도 같은 원칙을 따른다.
- 요청자는 결과를 받기 전이었으므로 잃는 정보가 없다. 다시 요청할 수 있다.
- 이미 전달된 결과는 회수할 수 없다. 문구가 이를 그대로 말한다.

## 미실행과 한계

- 정적 template 렌더다. 실제 인증·DB·격리 fixture(`remote.py fixture`) 왕복은 실행하지 않았다. 화면 구조는 같은 template·같은 fixture 데이터다.
- 브라우저 키보드·200% 확대·스크린리더는 not_run이다(tester 몫).
- 실제 Google 계정·동의 화면은 D011 뒤 coor 몫이다.
- 해제 뒤 approved gate의 회원 쪽 상태 표시는 DEV의 상태 전이 결정에 따른다. 어떤 상태든 '공유했다'고 표시하지 않는다.
- 판정에 쓰지 않은 캡처 54장·비교판·scratch 도구는 `raw/`에 있고 커밋하지 않는다.
