---
title: 회원 웹 글자 링크 44px 수정 직접 재측정·시각 수락 보고
status: review
updated: 2026-10-11
owner: designer
tasks: [SAR-SERVICE-UI-001-FIX-VISUAL]
summary: "고정 후보 295417a의 V-01 글자 링크 44px 수정을 Orca 내장 브라우저 320·390·1440에서 직접 재측정하고 화면을 열어 확인한 시각 수락 판정, 근거 구분, 미검증"
---

# SAR-SERVICE-UI-001-FIX-VISUAL — V-01 재측정·시각 수락

## 판정

**V-01 resolved (시각 기준). 제품 시각 수락: PASS 유지.** 글자 링크의 44px 미달은 해소됐다. 새 critical·high·medium 발견은 0건이다.

이 판정은 합성 fixture 29개 화면, Orca 내장 브라우저(Chromium 계열), 폭 3개에 한정한다. 실제 Google 로그인 왕복·실계정·운영 서버·실제 스크린 리더·200% 확대·다른 엔진은 하지 않았고 미검증으로 남는다. 원 보고서의 발견 V-01–V-12와 판정은 지우지 않았다. [원 보고서](../SAR-SERVICE-UI-001-VISUAL/report.md)에는 후속 기록 절만 추가했다.

## 대상과 방법

| 항목 | 값 |
|---|---|
| 고정 후보 | `295417acfdcd415ad0c768d164bd71dc07214ea9` |
| 구현자 세션 | `98481ab6-9f92-4bdd-ba02-c9881c19312c` (DEV, `ctx_4567096a3f21`) |
| 검수자 세션 | `059f2e8a-d31f-442f-a5af-ce1ac3255fe7` (등록 Claude Sonnet 5.5, 새 세션). 구현자와 다르다. |
| 기록 체크아웃 HEAD | `111cce03183b6e29e205a1aaa94371c295dcdcda` |
| 제품 바이트 | `git diff --stat 295417a HEAD -- internal adapters scripts go.mod`가 비었다. 후보와 같다. |
| 제품 변경 | DEV 제품 SHA `609fcfd`. 공통 CSS `a{display:inline-flex;align-items:center;min-height:44px;color:var(--accent)}` 한 규칙과 `.toc a` 중복 규칙 삭제. fixture HTML에서 이 규칙이 렌더됨을 직접 확인했다. |
| 정본 | D04 [service-member-ui.md](../../../design-docs/mockups/service-member-ui.md) 9절 "모든 상호작용 요소 높이 44px 이상" |
| 환경 | Orca 내장 브라우저 탭 1개, `file://` 합성 HTML. 기존 탭을 재사용했고 다른 세션 증거는 건드리지 않았다. |

1. 후보를 `git archive 295417a`로 scratchpad에 풀고 `go test ./internal/relay -run TestMemberUIFixtures -count=1`로 합성 HTML 29개를 새로 만들었다(Go 1.27.1, 통과). 리뷰 snapshot과 이 체크아웃에는 생성물을 쓰지 않았다.
2. [measure.py](measure.py)가 화면 29개 × 폭 3개(1440×900·390×844·320×640)를 열고 `innerWidth`가 목표 폭인지 assert한 뒤, 모든 `details`를 펼쳐 `a[href]`·버튼·입력·select·summary·textarea 전부의 높이·줄 수·겹침·인접 간격·가로 overflow를 쟀다. 87개 키.
3. 대표 PNG 26개(합계는 manifest 참조)를 저장하고 일부를 직접 열었다. base64는 출력하지 않았다.

## 증거와 구분

| 구분 | 내용 |
|---|---|
| 직접 DOM 측정 | [measurements.json](measurements.json): 87개 키, 링크 198건(건너뛰기 81 + 그 외 117), 상호작용 요소 전체 420건. |
| 직접 픽셀 관찰 | [manifest.json](manifest.json)의 PNG 25개 + `png/home-320-id-link.png`(sha256 `5d75066c53e78e56b05954793061348929847146f035981c74cb16d3c62b873d`, 추가 확인용). 열어 본 화면은 아래 "직접 관찰"에 적었다. |
| DEV 보고 | [FIX-DEV 기록](../../../exec-plans/phases/SAR-SERVICE-UI-001-FIX-DEV.md)의 87키·198링크·미달 0 주장은 이 측정과 일치한다. DEV 결과를 근거로 판정하지 않고 같은 측정을 새로 했다. |

## 측정 결과 (DOM, 87개 키)

| 검사 | 결과 |
|---|---|
| 높이 44px 미만 상호작용 요소 | **0건**. 최솟값 44px. 건너뛰기 링크(`position:absolute`)도 44px다. |
| V-01 글자 링크 종류 | 11종 모두 측정됐고 모두 44–44.8px다. "홈으로 돌아가기" 33건, "홈으로 이동" 15, "상태 새로고침" 6, "작업 화면으로 돌아가기" 6, "다른 주소로 시작" 3, "로그인 화면으로 이동" 3, 원격 동의 "내 agent 관리" 6, 확인 필요 링크 2종, 승인 요청 ID 링크 2종. |
| 가로 overflow | 폭 밖으로 나간 요소 0개. `scrollWidth - innerWidth` 최대 0. 320px 포함. |
| 인접 상호작용 박스 겹침 | **0건**. 같은 열의 위아래 이웃과 최소 간격은 0px(섹션 이동·확인 필요 목록처럼 박스가 맞닿음)이며 겹치지 않는다. |
| 글자 링크 줄바꿈 | 짧은 링크는 모두 1줄이다. 승인 요청 ID 링크(UUID `<code>`)는 하이픈 위치에서 2줄로 접힌다. 박스 높이는 44.8px, 폭은 320px 안(239px)이다. |
| 기존 버튼 | 버튼·입력·select·summary 높이는 44px 이상이다. 주 버튼 44.4px. 변경 전과 같은 크기다. |

측정 스크립트의 줄바꿈 표시 6건(`home@*` ID 링크)은 위 줄바꿈 항목의 접힘이다. UUID가 길어 생기며 결함이 아니다. 처음에는 `<code>` 내부 조각이 3줄로 읽혔다. 브라우저에서 링크 박스를 직접 재측정해 한 박스·한 폭 안의 하이픈 접힘임을 확인했다.

## 직접 관찰 (PNG를 열어 본 화면)

- `connection-prepared` 390: "연결 취소" 버튼과 "홈으로 돌아가기" 링크가 위아래로 있다. 두 요소 박스 간격은 4px이고 겹치지 않는다. 링크 글자는 버튼 아래 약 28px에 있고 밑줄 길이·글꼴 크기는 수정 전과 같다. 오조작을 부를 만큼 붙어 보이지 않는다. V-01이 지적한 "취소 아래 8px에 붙은 21px 링크"는 해소됐다.
- `refusal-limited` 320·390: 오류 띠 아래 카드 안에 "홈으로 돌아가기" 한 개뿐이다. 줄바꿈·잘림 없다. 링크 아래 여백이 위보다 약간 크다(아래 N-01).
- `verify` 390: 주 버튼·"코드 다시 받기"·"다른 주소로 시작"이 차례로 있고 링크가 버튼과 겹치지 않는다. 오류 띠·입력 폭은 수정 전과 같다.
- `remote-consent` 390: "이 설치 연결 승인"·"취소" 버튼과 "내 agent 관리" 링크가 겹치지 않는다. 이전 마크업 문제(V-10)는 그대로다.
- `home` 390: 섹션 이동 5개가 한 줄에 들어간다. "확인 필요" 목록 두 항목이 각각 44px 행으로 보여 간격이 커졌다. 읽기·탭에 문제 없다. 로그아웃은 `button`이라 변화 없다.
- `home` 320: 승인 요청 ID 링크가 두 줄로 접히고 대기·거절 배지가 아래에 붙는다. 카드 밖으로 나가지 않는다.

## 새 관찰

| ID | 등급 | 내용 | 판단 |
|---|---|---|---|
| N-01 | low | 카드 안 마지막 글자 링크가 44px 박스 + 문단 아래 여백 때문에 위쪽 여백보다 아래쪽이 조금 크다(`refusal-limited-320`에서 위 약 37px·아래 약 51px). 한 줄 카드에서 눈에 띈다. | 기능·접근성 영향 없음. 비차단. 후속 선택 개선(문단 아래 여백 조정)은 DEV 판단이다. |

## 변경하지 않은 기존 발견

V-02–V-12(low)는 이 과제 범위 밖이며 그대로 남는다. 특히 V-12(960px 미만 섹션 이동 위치)는 이번 CSS 변경으로 바뀌지 않았다.

## 미검증

- 실제 Google 로그인 왕복·실계정 데이터·운영 서버 렌더·운영 CSP 헤더.
- 실제 스크린 리더, 200% 확대, 다크 모드, Chromium 외 엔진.
- 키보드 Tab·포커스 외곽선은 이번에 다시 입력하지 않았다. 링크가 `inline-flex`가 되어 포커스 외곽선 크기가 바뀔 수 있으나 직접 확인하지 않았다. 원 보고서의 키보드 결과는 수정 전 후보 기준이다.
- 이 측정은 fixture 29개 화면에 한정한다. 실제 서버가 만드는 마크업이 fixture와 같다는 가정은 DEV 테스트(`member_forms.golden`)에 근거하며 이번에 직접 보지 않았다.
- 디자인 lint·다크 테마는 `project.md`가 미구성으로 정했으므로 해당 없음이고 실제 측정으로 보완했다.

## 재현

`KNOWSLINK_UI_FIXTURE_DIR=<새 폴더> go test ./internal/relay -run TestMemberUIFixtures -count=1` 후 `python3 -I measure.py <fixture 폴더> <출력 폴더>`. 도구 교훈은 원 보고서와 같다(`goto`가 뷰포트를 초기화하므로 설정 뒤 폭을 재확인한다).
