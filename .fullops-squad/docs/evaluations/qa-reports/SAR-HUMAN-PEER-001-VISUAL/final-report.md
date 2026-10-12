---
title: 사람 친구·대화 수정 후보 직접 시각 재검수
status: review
updated: 2026-10-12
owner: designer
tasks: [SAR-HUMAN-PEER-001-VISUAL]
summary: 12dc0c7의 V01–V03 개선과 실제 격리 화면 재측정 및 미검사 경계를 기록한다
---

# 사람 친구·대화 수정 후보 직접 시각 재검수

## 판정과 원기록 보존

후보 `12dc0c733d4d0bcdf0094fdebda7309da4c3fdd5`의 **제공된 V01–V03 수정 화면은 시각 PASS**다. V-01·V-02 medium과 V-03 low는 아래 관찰 범위에서 resolved다. 관찰 범위의 신규 critical/high와 미해결 blocking 시각 결함은 0건이다. 미관찰 상태나 전체 제품의 수락을 뜻하지 않는다.

[원후보 보고](report.md)의 `254ff6a98b12b7e51f0835491300f7695425d264` FAIL·21 PNG·미실행 기록은 그대로 보존한다. 이번 attempt는 `161c5363d3ac4a84bdcd43032cd69f51`이다. 착수 HEAD는 `a04d7c7b35adfaca0543e57970619e2ba469fa71`이며 coordinator의 정식 packet/inbox WIP를 사용했다. 원FAIL을 새 후보 PASS로 소급하지 않는다.

Purpose review, designer, codex / gpt-6.1-sol / high, frontend-design/fullops-frontend, test lite/subagent standard다. 작은 직접 재검수이므로 하위 위임은 사용하지 않았다. [D02 HP-01–16](../../../planning/product-specs/SAR-HUMAN-PEER.md), [현행 D04](../../../design-docs/mockups/service-member-ui.md), project.md와 fullops-common-0.3.3을 적용했다. 요구·토큰·보안 기준은 완화하지 않았다.

## 출처와 실제 실행 경계

[DEV manifest](../SAR-HUMAN-PEER-001-DEV-test/runtime/fe7a509bd959a4adbb48b881e5228bbb62e56602/captures-1791804944045033500/manifest.json)의 원본 14 PNG를 모두 view_image로 직접 열었다. 실제 합성 Alice/Bob 인증·API·DB의 홈/친구/두 메시지 대화, 9상태다. 파일 해시 14개와 PNG IHDR 크기를 독립 assert하여 exit0을 확인했다. 원본은 변형하지 않았다.

[제품 동일성·14 PNG 검수 기록](attempt-161c5363/source-and-original-images.json)은 Git으로 `fe7a509bd959a4adbb48b881e5228bbb62e56602`와 최종 `12dc0c733d4d0bcdf0094fdebda7309da4c3fdd5`의 차이가 FullOps 문서뿐임을 확인한다. human_pages.go/member_pages.go/human_pages_test.go의 blob과 SHA-256도 일치한다. 기존 캡처의 실행 SHA를 최종 SHA로 바꿔 표시하지 않는다.

실제 origin은 `http://localhost:2809`다. designer 별도 page `0e37999e-887b-4800-af8e-08ebf8770e35`를 isolated profile `695dbff1-0b06-4b1a-bb4b-84c628f3d5d3`에 만들고 `tab profile show`로 확인했다. DEV page `1240bd76-9c65-4de0-b889-bda238f274eb`와 원래 profile의 설정·쿠키를 변경하지 않았다. 실제 표시는 합성 앨리스 작성과 합성 밥 agent 답장, 원본과 같은 두 메시지·이름·구분값이다. runtime source 결속은 DEV 실행 보고와 위 제품 동일성에 근거한다. designer가 실행 바이너리를 새 빌드하거나 DB를 검사한 것은 아니다.

GET `/home`, `/home/friends`, `/home/agents`, `/home/new`와 실제 홈 DOM에서 확인한 conversation 링크만 열었다. 읽기·뷰포트·스크롤·캡처·textbox focus만 실행했다. 클릭/입력/POST/상태 변경·쿠키 조회/설치는 하지 않았다. 탐색/shadow/CUAWright는 OFF다. Orca 1.4.224의 버전 맞춤 browser/orchestration 안내를 사용했다. 새 도구·호스트·의존성 설치는 없다.

## 직접 이미지 관찰

| 화면 | 직접 연 원본 PNG | 판정 근거 |
|---|---|---|
| 홈 | `00-home-{320,390,1440}-s1of1.png` 3장 | 사람 아래 agent 들여쓰기와 최대 두 줄 이름이 보인다. 중복 목록 제목이 사라졌고 h1 대화는 유지된다. |
| 친구 | `01-home-friends-320-s{1,2}of2.png`, `01-home-friends-{390,1440}-s1of1.png` 4장 | 코드 안내가 한 문장이다. 이름 공개와 현재·향후 agent 연결 동의는 보인다. 320은 문서 스크롤로 버튼에 접근한다. |
| 대화 | `02-home-conversation-320-s{1,2,3}of3.png`, `02-home-conversation-{390,1440}-s{1,2}of2.png` 7장 | 고정 대상·전체 긴 이름·owner/agent 작성 구분·전달됨·24h 안내·닫힌 상세가 보인다. desktop 첫 조각에서 textarea와 보내기가 보인다. |

원본 320×640·390×844는 DPR1이다. desktop은 CSS1440×900/DPR0.75이며 PNG1080×675다. native1440/DPR1이나 실200% 검수로 표시하지 않는다.

실제 별도 페이지에서도 [홈320](attempt-161c5363/actual-home-320.png), [홈1440](attempt-161c5363/actual-home-1440.png), [대화320](attempt-161c5363/actual-conversation-320.png), [대화1440](attempt-161c5363/actual-conversation-1440.png), [내 agent390](attempt-161c5363/actual-agents-390.png), [새 대화390](attempt-161c5363/actual-new-390.png), [대화390 보내기](attempt-161c5363/actual-conversation-390-verified-bottom.png)를 직접 열었다. 동일한 실제 내용과 문서 흐름을 확인했다. 이 PNG는 이번 designer 캡처이며 DEV 원본과 구분한다.

## V01–V03 결과

| 항목 | 판정 | 실제 근거와 한계 |
|---|---|---|
| V-01 기본 작성란의 기술 안내·반복 대상 | resolved | 화면의 기본 면에 4096 bytes/180초/Enter 안내가 없다. 닫힌 전송 안내·대상 상세의 DOM에는 값과 전체 대상이 남아 있다. 친구·내 발신·수신 이름과 24h 안내는 기본 면에 남는다. |
| V-01 동명·동의·보존 | 제공 범위 PASS | 실제 내 agent 화면의 같은 이름 두 개를 node-local/53467732·58746875로 구분한다. 두 메시지는 사람·agent·구분값과 내가 작성/에이전트가 작성으로 구분된다. 친구 동의와 보존 상세의 문구를 확인했다. 오류 발생·동명 발신/수신 쌍의 별도 실행은 이번 범위에 없다. |
| V-02 긴 이름 두 줄·전체 이름 | resolved | 홈의 agent 이름은 320/390/1440에서 computed clamp2, 높이≤2×line-height다. 줄임표 뒤 전체 이름은 실제 대화 h1/수신 문장과 목록 접근성 이름에 남는다. |
| V-02 두 메시지 desktop 작성/보내기 | resolved | DEV geometry textarea top618.109375/send bottom833.296875를 원PNG와 대조했다. 별도 실제 CSS1440×900에서는 textarea top617.927979/send bottom833.094555다. 둘 다 첫 화면 안이다. 약0.2px 차이는 별도 surface의 스케일/반올림 조건이며 DEV 숫자를 새 측정값으로 복사하지 않았다. |
| V-02 좁은 화면·스크롤 | 제공 범위 PASS | 320/390에서 긴 헤더와 전체 대상은 문서 흐름으로 줄바꿈한다. 보내기는 스크롤 후 화면 안에 온다. 15상태에서 가로 overflow0·높이44px 미달 가시 제어0·중첩 세로 스크롤0을 직접 측정했다. 실제 키보드 가두기 여부는 not_run이다. |
| V-03 제목·목록 이름 | resolved | 홈 h1 대화는 하나이며 중복 h2가 없다. snapshot의 complementary 이름 최근 대화와 링크의 전체 대상 이름은 유지된다. 열린 desktop 대화의 목록 제목은 최근 대화다. |

[실제 surface 15상태 측정](attempt-161c5363/actual-surface.json)은 다섯 경로×세 폭의 관측·독립 assertions를 보존한다. 흰 바탕·청록 선택선·텍스트 상태·장식 없는 구조는 기존 방향과 같다. 픽셀 대비 전수 검사는 하지 않았다.

## 실행기 보정·미실행·후속

별도 page에서 첫 `set viewport 320 640 1`은 실제 CSS292×584/DPR1.095445를 반환했다. 원시 요청 폭 assertion과 첫 보정 높이 assertion은 exit1이었다. 이 실행기 조건을 제품 overflow 실패로 해석하지 않는다. 요청 크기·scale만 조정하여 실제 CSS320×640/DPR약1, 390×844/DPR약1, 1440×900/DPR약0.75를 각각 assert했다. 요청/실제 치수와 보정 probe를 actual-surface.json에 보존한다. 일부 부동소수 값은 근사다. 원 실패·partial 자료도 보존한다. viewport 보정은 실제 브라우저 200%가 아니다.

| 검사 | 결과 | 영향·담당·재개 조건 |
|---|---|---|
| 실제 Tab/Enter/Space·focus-visible·키보드 가두기 | not_run | 실제 document.hasFocus=false다. 허용된 `orca focus`는 TEXTAREA를 active로 만들지만 focus-visible=false이며 실제 키 입력 PASS가 아니다. tester/designer가 native focus 가능한 허가 surface에서 확인해야 한다. [관측](attempt-161c5363/focus-observation-verified.json). |
| 실제 브라우저200% | not_run | focused native surface가 없다. viewport720×450나 CSS zoom을 실200%로 대체하지 않았다. 확대 후 대상/작성/보내기 접근 수락은 열려 있다. |
| 로그인·첫 이름·pending·친구 상세·해제 읽기 전용·gate/오류 | not_run | 현재 제공14PNG와 허용 GET fixture에 없다. 상태별 문구·복구·위험 동작의 직접 시각 수락은 열려 있다. coor에 msg_10c0551d7285로 요청했고 msg_e68ff88264d1은 tester에게 요청 후 미제공 상태를 명시하고 현재 범위를 완료하도록 지시했다. source/form/static 검사를 실제 상태 시각 PASS로 올리지 않는다. |
| 동명 두 endpoint 선택 확인·상세 펼치기 | 이번 직접 실행 not_run | 내 agent 동명 구분은 직접 봤지만 현재 열린 대화 endpoints는 다른 이름이다. 닫힌 상세의 원문 존재는 DOM 읽기로 확인했다. 펼침·선택 상태 실행은 허용된 읽기/focus 범위를 넘어 수행하지 않았다. DEV 회귀/geometry 근거는 독립 새 실행과 구분한다. |
| 실제 Google·실계정·실메일·외부 앱·운영 | not_run | 이번 합성 UI 범위 밖이다. |
| mobile 가상 키보드·스크린 리더·다른 엔진·다크 모드 | not_run | 실행하지 않았다. 기능 QA·접근성 전체 적합성을 주장하지 않는다. |

첫 focus 명령의 따옴표 없는 `@e11`은 PowerShell splat 때문에 local CLI invalid_argument였다. 따옴표로 고친 뒤 focus가 성공했다. 중간 `actual-conversation-390-bottom.png`는 실제391px 캡처이므로 최종390 근거에서 제외한다. 최종 근거는 actual-conversation-390-verified-bottom.png다. 이전 파일을 삭제하지 않았다.

착수 환경 전체 열거로 Orca hook/launch 세션 토큰 값이 도구 출력에 포함됐다. 제품 인증값은 조회하지 않았다. 값은 증거/보고/커밋에 복사하지 않았고 msg_ac159db37045로 coordinator에게 자격 교체 판단을 인계했다. 환경 전체 열거는 중단했다. 이 운영 도구 노출을 제품 UI 결함으로 재분류하지 않는다.

## 문서·게이트·인계

startup `jev_packet.check(... required=True)`는 identity/HEAD/source PASS다. packet의 65개 비선택 path/category를 확인하고 관련 현행 절·역사 자료 경계를 [수동 범위](attempt-161c5363/packet-manual-scope.json)에 기록했다. producer mismatch4는 의미 탐색 한계다. 12개 unknown 경로는 실제 존재·제목·현행 키를 [별도 범위](attempt-161c5363/unknown-manual-scope.json)로 확인했다. 역사 문서 전체 재독이나 기능 QA 완료를 주장하지 않는다. 현재 지시는 역할 인박스 하나이며 logs/pending은 참고만 했다. coor PLANS/board와 다른 역할 파일은 수정하지 않았다.

`frontend.py --check`는 exit0/configuration_ready/problems0이다. 인증·MCP 성공 검사가 아니다. 문서-only 변경이므로 새 제품 테스트를 만들지 않는다. 최종12dc의 DEV `make lint`/`make test` native0과 FullOps ERROR0/WARNING1/unavailable0은 외부 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/v01-v02-final-report.md` 및 v01-v02-final-lint.json의 기존 실행으로 재사용한다. 이번 designer 문서 HEAD의 새 제품 실행으로 표시하지 않는다.

자기 문서/링크/패치/제품파일 경계·packet completion·deliverables strict·finish 검사를 수행한다. 커밋 후 FullOps lint 기준은 지시서 base `12dc0c733d4d0bcdf0094fdebda7309da4c3fdd5`다. 최종 결과와 native 종료코드는 외부 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-visual-final-161c5363/final-validation.json`, final-lint.json에 보존한다. Windows PATH에 make가 없으면 등록 두 명령을 unavailable로 기록하며 제품 PASS로 바꾸지 않는다. 별도 shim이나 검사 설정 변경은 하지 않는다.

work.py finish로 같은 키 새 attempt의 지시/결과 전문을 자기 날짜별 로그에 추가한다. 최종40SHA는 worker_done에 고정한다. 모든 임시 자료와 별도 page를 보존한다. 삭제·이동·운영은 수행하지 않았다. 후속 독립 기능 QA/코드 리뷰와 미검사 상태·키보드·실200%의 수락 및 main 통합은 coordinator가 조정한다.
