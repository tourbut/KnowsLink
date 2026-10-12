---
title: 사람 친구·대화 UI 직접 시각 검수
status: review
updated: 2026-10-12
owner: designer
tasks: [SAR-HUMAN-PEER-001-VISUAL]
summary: 원래21 PNG와 최종 소스 동일성을 검수하고 두 medium 수정 및 실제 surface 검사 범위를 기록한다
---

# 사람 친구·대화 UI 직접 시각 검수

## 판정과 대상

시각 수락은 FAIL이다. V-01·V-02 medium 수정이 필요하다. 이번에 관찰한 화면의 critical/high는 0건이다. 미관찰 상태의 결함 부재를 뜻하지 않는다. 검수 과제 완료와 제품 수락을 구분한다.

- 최종 제품 후보: `254ff6a98b12b7e51f0835491300f7695425d264`.
- 착수 문서 HEAD: `7c29f87a8cba756457c3873badef4881ad8cdf03`. coor가 제공한 인박스·packet WIP를 보존했다.
- 역할/실행: designer, fresh Codex 세션, codex / gpt-6.1-sol / high, 주 스킬 frontend-design. 하위 위임은 사용하지 않았다.
- 정본: [D02 HP-01–16](../../../planning/product-specs/SAR-HUMAN-PEER.md), [현행 D04](../../../design-docs/mockups/service-member-ui.md), fullops-common-0.3.3와 project.md.
- 허용 범위: 합성 loopback의 읽기·스크롤·뷰포트·캡처·포커스 확인. 상태 변경 POST는 tester 소유다. 제품 코드·coor PLANS/board·tester 증거는 수정하지 않았다.

## 증거 시점과 직접 관찰

먼저 [DEV runtime 보고](../SAR-HUMAN-PEER-001-DEV-test/runtime/report.md)의 정상 manifest에 속한 원래 PNG 21장을 모두 직접 열었다. 출처 SHA는 `4eaefa6e0a5cf623f0d8d8080d918086f86030ec`다. 정적 HTML이 아닌 당시 실제 합성 OTP·auth/API/DB 화면이다. 최종 SHA의 새 실행으로 표시하지 않는다.

[21장 목록·해시·규격·직접 관찰](original-capture-review.json)은 원래 manifest의 21개 SHA-256과 PNG IHDR 크기를 다시 대조한 기록이다. 5개 화면×3폭의 15개 상태다. 원본 이미지는 변형하거나 합성하지 않았다.

최종 실행/환경 정본은 외부 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/final-report.md`다. 이 문서와 과거 runtime 보고의 다른 SHA·검사 결과는 서로 다른 실행 시점이다. 과거 실패/제약을 지우지 않았고 최신 검사로 과거 캡처를 다시 이름 붙이지 않았다.

[source-delta.json](source-delta.json)은 두 SHA의 memberStyle 전체 및 human-home/conversation-list/friends/new/conversation template를 Git blob에서 비교한다. 모두 바이트 동일하다. 기존 21 PNG가 보여주는 레이아웃·문구 결함은 최종 소스에도 남는다. 이름 maxlength 40→80 UTF-16 unit, 이름 검증, agent 제한 안내, unavailable 번역, remote authorship metadata는 기존 이미지의 증명 범위 밖이다.

| 직접 연 화면 | 원본 PNG 이름 | 조건·관찰 |
|---|---|---|
| 대화 홈 | `00-home-{320,390,1440}-s1of1.png` | 세 메뉴와 사람 아래 agent 들여쓰기가 보인다. 대화 제목이 연속해서 반복된다. 긴 agent 줄은 320/1440에서 3줄 이상이다. |
| 친구 | `01-home-friends-320-s{1,2}of2.png`, `01-home-friends-{390,1440}-s1of1.png` | 친구 이름과 코드 입력이 분리됐다. 이름 공개·현재/향후 agent 동의가 명시됐다. 받은/보낸 요청은 이 fixture에 없다. |
| 내 에이전트 | `02-home-agents-320-s{1,2}of2.png`, `02-home-agents-{390,1440}-s1of1.png` | 같은 이름 2개를 client/짧은 구분값으로 식별한다. 키/전체 ID는 닫힌 상세 뒤다. 추가·연결 완료를 구분하는 안내가 있다. |
| 새 대화 | `03-home-new-{320,390,1440}-s1of1.png` | 친구/내 에이전트 상대 선택 진입이 보인다. 발신/수신 select 단계는 미포함이다. |
| 대화·작성 | `04-home-conversation-320-s{1,2,3}of3.png`, `04-home-conversation-{390,1440}-s{1,2}of2.png` | 고정 발신/수신, 사람+agent, 내가 작성/에이전트가 작성, 전달됨, 보존 안내, 작성란을 확인했다. V-01·V-02를 관찰했다. |

원본 경로는 [manifest](../SAR-HUMAN-PEER-001-DEV-test/runtime/4eaefa6e0a5cf623f0d8d8080d918086f86030ec/captures-1791801092715284800/manifest.json) 기준이다. 320×640·390×844는 DPR1이다. desktop은 CSS1440×900, DPR0.75, PNG1080×675다. native 1440px/DPR1 이미지 검수로 표시하지 않는다. 최초 반복 타일 full-page 증거는 부적합 기록으로 보존됐으며 사용하지 않았다.

## 발견과 수정 수락 조건

### V-01 — 작성란의 기술 설명과 반복 대상 문장, medium

`04-home-conversation-320-s3of3.png`와 `390-s2of2.png`에서 textarea 다음에 `최대 4096 UTF-8 bytes`, `최대 180초`, `Enter는 줄바꿈`이 상시 표시된다. 그 아래 대상 문장은 client와 짧은 ID를 다시 포함한다. 320폭에서는 안내 3줄과 대상 5줄이 보내기 앞을 차지한다. 헤더·메시지에도 같은 대상 정보가 이미 있다. 이는 D04의 이름·대화 중심 방향과 사용자 심플화 요청에 맞지 않는다.

DEV 수정 수락 조건:

1. 기본 작성 면에서는 친구·내 발신·상대 수신 이름을 명확히 유지한다. client/구분값 반복은 동명 선택의 구분에 필요할 때만 쓴다.
2. bytes·전달 기한 등 기술 상세는 접힌 안내 또는 해당 실패 안내로 옮긴다. 값·실제 제한·Enter가 자동 발송하지 않는 동작은 유지한다.
3. `메시지는 24시간 안에 사라집니다`, 보존 한계, 동의, 해제/차단/철회 결과, 오류·복구 안내는 관련 위치에 유지한다.
4. 320/390/1440에서 대상 확인과 보내기 접근을 다시 직접 확인한다. 서버 검증·보안 정책은 변경하지 않는다.

친구 추가의 `코드만으로 상대의 이름이나 에이전트는 확인할 수 없습니다`는 다음 동작을 돕지 않는 설명이다. 기본 면에서는 `친구가 공유한 코드를 입력하세요`로 충분하다. 이름 공개와 현재·향후 agent 연결 동의는 그대로 유지한다. 이 문구 정리는 V-01의 같은 수정 범위다.

### V-02 — 긴 이름 목록과 desktop 작성 위치, medium

`00-home-320-s1of1.png`의 agent 이름은 3줄이다. `00-home-1440-s1of1.png` 및 `04-home-conversation-1440-s1of2.png`의 목록도 3줄이다. 현행 D04의 목록 이름 최대 두 줄 기준을 충족하지 않는다. 전체 이름을 상세에서 확인하는 요구도 함께 유지해야 한다.

두 메시지뿐인 CSS1440×900 대화에서 textarea는 첫 화면 맨 아래에서 시작하고 보내기는 `s2of2`에 있다. 긴 헤더/중복 대상/보존 영역과 메시지 여백이 작성 위치를 밀어낸다. 긴 이력이 있는 대화의 모든 내용을 첫 화면에 넣으라는 요구는 아니다. 지정된 짧은 fixture에서도 D04의 첫 화면 작성 위치 목표가 충분히 드러나지 않는다.

DEV 수정 수락 조건:

1. 긴 agent 목록 이름을 320/390/1440에서 최대 두 줄로 표시한다. 전체 이름·동명 구분은 대화 헤더/선택 확인/상세에서 접근 가능해야 한다.
2. 두 메시지 fixture의 1440×900에서 고정 대상·본문·작성 위치·보내기가 첫 화면에 접근 가능하게 한다. 중첩 스크롤·포커스 가두기는 만들지 않는다.
3. 320/390에서는 문서 흐름과 전체 이름 줄바꿈을 유지한다. 200%에서도 대상·입력·보내기를 접근할 수 있어야 한다.

### V-03 — 홈의 반복 제목, low

모바일 홈에서 세 메뉴의 `대화`, h1 `대화`, 목록 제목 `대화`가 연속한다. h1과 목록 제목의 반복을 줄이면 실제 사람 이름이 더 빨리 보인다. 페이지 h1과 목록의 접근성 이름은 유지한다. 비차단 개선이며 V-01·V-02 수정 후 별도 장식이나 새 기능을 추가할 이유는 없다.

## 검사 범위와 한계

| 검사 | 판정 | 근거·한계 |
|---|---|---|
| 원래 21 PNG 직접 시각 검수 | completed | 모두 view_image로 직접 열었다. 새 최종 SHA 실행이 아니다. |
| PNG 해시/규격 21개, manifest 15상태 | PASS | Python hashlib/struct assert, exit0. |
| 세 메뉴·이름·동명 구분·owner/agent 표시 | 부분 PASS | 제공된 정상 사례에 한정한다. V-01·V-02는 FAIL이다. |
| 토큰/형태 | 부분 PASS | 흰 바탕·청록 선택선·텍스트 상태·4px 형태·장식 없는 구조를 관찰했다. 최종 memberStyle의 6색 값은 D04와 같다. 픽셀 대비 전수 검사는 하지 않았다. |
| 320/390 가로 넘침·44px | 근거 재사용 | DEV manifest의 overflow0·작은 제어0과 이미지 가로 잘림 부재를 대조했다. 독립 새 DOM 측정은 아니다. |
| CSS1440 반응형 | 근거 재사용 | DPR0.75 조건을 보존한다. native1440/DPR1은 not_run이다. |
| 최종 SHA 합성 loopback 직접 화면 | not_run | tester가 같은 Run에서 독립 QA/fixture를 준비 중이며 URL/profile/page를 아직 전달하지 않았다. coor 회신에 따라 원후보 기록을 마무리한다. |
| 200% 브라우저 확대 | not_run | 아직 최종 인증 surface를 받지 못했다. 뷰포트 축소나 CSS zoom을 실제 브라우저 확대 PASS로 바꾸지 않는다. |
| 실제 Tab/Enter/Space·focus-visible | not_run | DEV의 document.hasFocus=false·Tab BODY는 실행기 제한이다. 이미지/DOM focus를 실제 키보드 PASS로 바꾸지 않는다. |
| 로그인/이름·요청 pending·친구 상세·해제 읽기 전용·gate/오류 | not_run | 원래 5화면 manifest에 없는 상태다. 최종 fixture에서 제공받은 범위만 추가할 수 있다. |
| 실 Google·실계정·실메일·외부 플랫폼 왕복·운영 | not_run | 이번 합성 UI 검수 범위 밖이다. |
| 가상 키보드·스크린 리더·다른 엔진·다크 모드 | not_run | 실행하지 않았다. 다크 모드는 범위 밖이다. |

## 절차와 담당

`frontend.py --repo . --host codex --model gpt-6.1-sol --effort high --check`는 exit0, configuration_ready, problems0이다. 설치 스킬은 직접 읽었다. 인증·브라우저 QA·MCP 성공을 뜻하지 않는다. 도구·호스트·모델을 바꾸거나 설치하지 않았다. Orca 1.4.224의 버전 맞춤 orchestration/browser 안내를 읽었다. 이 designer 체크아웃의 tab list는 비어 있었으며 실제 제품 탐색·POST는 실행하지 않았다.

착수 `jev_packet.py --force --write-inbox` 뒤 `jev_packet.check(... required=True)`는 identity/HEAD/source PASS였다. packet의 producer mismatch 4건은 과제별 route/find/context가 없거나 다른 시도라는 탐색 한계다. 필수 정본과 원본 이미지·소스 delta를 직접 확인했다. [packet-source-scope.json](packet-source-scope.json)은 62개 문서의 메타데이터와 관련 구간 분류다. 역사적 문서의 전체 재독이나 제품 기능 QA를 주장하지 않는다. unknown 11개 큰 문서는 존재/제목/현행 과제 포함 여부를 별도 확인했다. coor PLANS/board는 읽기 경계만 확인하고 수정하지 않았다.

V-01·V-02와 수락 조건은 coor에 즉시 전달했다. DEV가 제품 수정을 맡는다. tester는 독립 기능 QA와 fixture를 소유한다. 수정 후 coor가 새 고정 SHA의 영향받은 화면 검수·독립 리뷰·필수 검사·통합을 조정한다. 이번 시각 FAIL을 과거 polish PASS나 DEV 자체 검사로 해소하지 않는다.

coor는 질문 `msg_ef0f8a0d5800`에 회신 `msg_7d224a4c95a0`으로 원후보 추가 surface를 기다리지 않고 완료하도록 지시했다. 당시 tester는 원후보 QA/코드리뷰 중이며 URL/profile/page가 미전달이었다. 실제 surface·200%·focus·빠진 상태는 DEV 수정 완료와 tester fixture 전달 뒤 같은 키의 새 attempt에서 검수한다. 이 실행 한계로 이번 제품 FAIL이나 미충족 조건을 PASS로 바꾸지 않는다.

62개 packet path/category 결과는 completed14/no_change48로 보존했다. unknown/누락을 완료로 위장하지 않았으며 partial 수동 확인 범위와 역사 자료 재독 한계를 명시했다. `deliverables.py --repo . --strict`는 13종·문제0·경고0, exit0이었다. `git diff --check`도 exit0이었다. 완료 전문·packet 검사와 `work.py finish`를 거쳐 자기 인박스를 아카이브한다. 최종 기록 SHA는 worker_done에서 고정한다.

최종 커밋의 FullOps lint·문서 검사·패치/링크/제품 파일 경계 검사 원본은 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-visual/final-validation.json`과 같은 폴더 `final-lint.json`에 보존한다. lint 기준은 지시서 base `254ff6a98b12b7e51f0835491300f7695425d264`다. Windows PATH의 make 부재로 등록 제품 명령을 실행할 수 없으면 unavailable로 기록한다. 제품 코드/설정/의존성을 바꾸지 않았으므로 새 제품 시험은 만들지 않으며 별도 make shim·다른 후보 검사로 이 문서 HEAD의 성공을 꾸미지 않는다. DEV 최종 후보의 실제 make lint/test 통과 근거는 외부 final-report.md·final-hook-base-lint.json에 있으며 이번 독립 제품 실행이 아니다.

삭제·이동을 하지 않았으며 모든 임시 자료와 기존 증거를 보존한다. 제품 시각 FAIL 해소는 DEV 수정·후속 고정 SHA 검수의 책임이다.
