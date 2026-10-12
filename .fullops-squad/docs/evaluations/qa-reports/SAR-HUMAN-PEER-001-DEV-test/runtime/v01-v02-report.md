---
title: 사람 대화 V01/V02/V03 수정 검증과 재검수 인계
status: review
updated: 2026-10-12
owner: dev
tasks: [SAR-HUMAN-PEER-001-DEV]
summary: 실제 합성 화면에서 문구 축약·두 줄 목록·첫 화면 보내기 접근과 한계를 기록한다.
---

# 독립 시각 발견 수정

새 attempt는 `6b359aa6693c40eb98fdb20f707d0d08`이며 제품 기준은 `254ff6a98b12b7e51f0835491300f7695425d264`다. 원구현 완료·이전 실패·254ff6a의 독립 QA 진행과 과거 PNG는 보존한다. 자기 구현/캡처 SHA는 `fe7a509bd959a4adbb48b881e5228bbb62e56602`다. 이후 완료 기록 커밋은 제품 소스를 바꾸지 않으며 최종40SHA/source byte 동일성/등록 검사 결과는 아래 외부 final report에 고정한다.

## 수정 범위

- V01: 이름만 기본 발신/수신/보내기 대상에 표시한다. 발신·수신 이름이 같으면 짧은 구분값을 유지한다. 선택 확인의 전체 라벨은 유지한다. client/구분값·4096 UTF-8 bytes·180초·Enter 안내는 닫힌 전송 안내·대상 상세로 이동한다. textarea의 aria-describedby 연결과 서버 검증은 유지한다.
- V02: .agent-line display:block이 .name-preview의 두 줄 표시를 덮어쓰는 CSS 충돌을 좁은 결합 선택자로 고쳤다. 대화 헤더/보존/메시지/작성 여백은 기존 4/8/16px 토큰을 사용한다. 입력란은96px이며 사용자 세로 크기 조절을 유지한다. 중첩 스크롤·고정 높이·JS는 없다.
- V03: 홈 h1 대화와 aside의 최근 대화 접근성 이름을 유지한다. 홈 목록의 중복 h2를 생략하고 열린 대화 목록에서는 최근 대화로 표시한다. 친구 코드 기본 설명을 짧게 바꿨으며 이름 공개·현재/향후 agent 동의는 유지한다.

원제품 기준 이후 source delta는 internal/relay/human_pages.go, member_pages.go, human_pages_test.go뿐이다. 권한/DB/adapter/handler/선택 확인/POST/hidden endpoint/attempt/동의/보존/오류 회복은 변경하지 않았다. tester/designer 파일과 fixture·coor PLANS/board는 수정하지 않았다. D02/D04 의미와 여섯 색·글꼴·44px 제어·3px 포커스·CSP·noJS/no새의존성은 유지한다.

## 실제 실행과 증거

실제 runtime은 `/tmp/knowslink-humanpeer001-fe7a509bd959-44b44bdd/qa-13a5db3f/repo`다. DB/내부망은 `qa-13a5db3f-44b44bdd-postgres`/`qa-13a5db3f-44b44bdd`이며 공개 DB 포트는 없다. remote relay127.0.0.1:39545/SMTP39135, Windows loopback `http://localhost:2809`, forward PID51460을 보존한다. relay PID51634, sink PID51633이다. prepare/install/build/migration/healthz/auth/fixture는 각각 native0이다. Go1.27.1/Node22.22.2/Postgres17.11/Orca1.4.224/Chrome150.0.7871.250이다.

기존 승인 helper remote.py/runtime/auth/human_fixture_remote.py를 자기 정확한 SHA로 재사용했다. 실제 .invalid SMTP OTP와 profile·agent connect·친구 요청/명시 동의·owner 발송·signed agent 답장/persist/ACK로 두 메시지를 만들었다. 다른 Origin403/정상 owner303을 재확인했다. private session·OTP·키·credential은 로그/보고에 출력하지 않았다. 실계정/실메일/외부 상대/운영 배포는 없다.

복구 후 정본 캡처는 [manifest](fe7a509bd959a4adbb48b881e5228bbb62e56602/captures-1791804944045033500/manifest.json)다. 홈·친구·두 메시지 대화를320×640/390×844 DPR1과 CSS1440×900 DPR0.75에서 캡처했다.9상태/14개 원본 viewport PNG이며 이미지 합성·변형은 없다. desktop PNG는1080×675로 native1440/DPR1 검수로 표시하지 않는다. 반복 타일을 피하는 기존 실행 조건을 유지했다. 홈320·대화390/desktop을 직접 열어 확인했다.

[geometry.json](fe7a509bd959a4adbb48b881e5228bbb62e56602/captures-1791804944045033500/geometry.json)은 profile/source·8개 viewport 상태·두 줄 이름·닫힌 상세·desktop 보내기·좁은 화면 스크롤 접근·PNG hash/크기를 검사한다. 최대 두 줄 이름은 실제 computed height≤2×line-height다.1440×900의 textarea top618.109375/bottom714.109375, 보내기 bottom833.296875다.9개 캡처 상태의 가로 넘침0/44px 미만 제어0이다. 전체 이름은 대화 헤더와 선택/상세에서 유지한다. native details를 열어4096bytes/180초/전체 대상이 존재함을 확인했다.

focused `go test ./internal/relay -run TestHuman(Pages|Friend|Conversation|Empty|Gate) -count=1`은 native0이다. 동명/다른 이름 대상 구분·closed details·escape·폼 고정/읽기 전용/응답 유실·빈 선택·gate를 검증한다. 테스트 추가는 DEV 기존 human_pages_test.go의 좁은 회귀 하나뿐이다. 전체 SQL은 UI-only source delta 근거로 이전 정확한254ff6a의 native0 결과를 재사용한다. 정본은 기존 외부254ff6a-final-evidence/integration-result.json이며 새 SHA SQL 실행으로 주장하지 않는다.

## 실행기 한계와 격리 복구

첫 캡처 `captures-1791804826681751400`은 새 탭의 default profile에 합성 cookie를 설치한 상태였다. 격리 누락을 확인하고 자료를 보존했다. 같은 page `1240bd76-9c65-4de0-b889-bda238f274eb`를 새 isolated profile `695dbff1-0b06-4b1a-bb4b-84c628f3d5d3`으로 전환했다. 실제 QA Alice 세션을 메모리에서 재설치하고 위 정본 캡처와 geometry를 다시 실행했다. profile 검사와 합성 앨리스/밥의 실제 두 메시지 UI를 대조했다. 원격 runtime과 합성 actor는 같은 source/DB다. 첫 default 설치/복구 사실은 coor msg_19936e7c0702로 전달했다. coor msg_c8067ab02027은 기존4359b1d6 isolated page에 영향 없음을 확인하고 자료 보존을 지시했다. default cookie 삭제/정리는 하지 않았다.

실제 Tab 입력 후 document.hasFocus=false/activeElement BODY다. Control++ 뒤 innerWidth1440/DPR0.75가 그대로라 실제200% 확대는 not_run이다.720×450은 확대 대응 reflow 근사 검사일 뿐 브라우저200% PASS가 아니다. 키보드/포커스-visible 실행·실200%·mobile 가상키보드/스크린리더는 독립 designer/tester가 실제 보이는 surface에서 확인해야 한다. 자동 포커스 이동은 추가하지 않았다. 외부 geometry helper의 첫 self-check는 CLI eval이 boolean을 문자열로 반환하여 exit1이었다. 원 DOM 보내기 위치는 정상이며 반환형 해석을 고친 뒤 동일 assertion은 native0이다. 제품 실패로 재분류하지 않는다.

## 최종 게이트와 인계

완료 기록 커밋 후 `lint.py --from 254ff6a98b12b7e51f0835491300f7695425d264`를 현재 최종 HEAD에서 실행한다. 등록 product-lint=make lint/product-test=make test를 변경·생략하지 않는다. 기존 bridge가 exact clean Linux snapshot에서 원 make 종료코드를 전달한다. 최종 정본은 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded/v01-v02-final-lint.json`과 `v01-v02-final-report.md`다. 이 문서는 실행 전에 작성했고 결과를 선행 PASS로 표시하지 않는다.

39개 packet path/category와 unknown15/remaining12 수동 범위를 outcomes에 기록했다. 이전 outcomes/packet은 history와 기존 Git 이력에 보존했다. phase에는 원완료를 유지하고 새 수정절만 추가했다. work.py finish는 동일키 새 attempt의 지시/완료 전문을 날짜별 로그에 append한다. coor가 새 최종 SHA의 독립 시각 재검수·고정 SHA 코드 리뷰·main 수락/통합을 조정한다. 자기 개선 판정을 독립 수락으로 표시하지 않는다. 모든 임시자료/환경/DB/container/process는 보존한다.
