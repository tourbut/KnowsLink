---
title: SAR-PUBLIC-MESSAGES-001-UI-FIX — QA 기록
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-MESSAGES-001-UI-FIX]
summary: 자기 fixture의 시각 관측·원실패·42PNG·재사용 경계와 후속 QA를 기록한다
---

# SAR-PUBLIC-MESSAGES-001-UI-FIX — QA 기록

## 결과와 재현 범위

[직접 시각 판정](../../design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md)이 제품 UI 판단 정본이다.
대상은 exact `dfc70caa748a90614b02d48c78b4651345938339`의 별도 clean detached clone이다.
일반회원 실제 로그인·자기SQL/서명 gate·localhost SMTP로42PNG를 얻고 모두 직접 확인했다.
정상 Deny 클릭/Enter desktop/mobile 네 조건과 owner desktop의 canonical GET200·denied 저장·비활성·실제 복귀를 확인했다.
CSRF/타 회원403·신규/정리 동시 한도429·cleanup/receipt rate429는 안전 오류를 보인다.
신규 rate 포화 뒤 Deny 저장은 성공하지만 결과 GET429다. 이 즉시 결과 조건은 LIMITED로 분리한다.

## 파일과 종료코드

| 파일 | 의미 |
|---|---|
| manifest.json |42PNG의조건·exactSHA·viewport·크기·hash·직접 판정 |
| observations.json | 첫 브라우저의18관측·40캡처·자기fixture변경·실제응답/저장상태/홈 클릭 |
| receipt-boundary.json | 후속 receipt429·실제 홈 클릭·2캡처 |
| browser-run-first.log/.exit | 원run exit1. 시나리오 저장 후 profile 회수의 Errno39이며 제품 assertion 실패는 아니다 |
| teardown-recovery.json/.exit | 자기 listener/profile 후속 회수. 최종 exit0. 최초 자기검사 문자열 오분류를 명시 |
| receipt-boundary.log/.exit | 실제 로그인·receipt429·홈 클릭 후속 exit0 |
| build-relay/.exit·build-migrate/.exit·migrate/.exit | 변경 없는 제품 빌드/migration, 각각 exit0 |
| source-identity.json | 변경 없는 UX06 표시/wire 의존성의 원09c/dfc hash; absent 경로는 근거 제외 |
| snapshot-identity.json | exactSHA·detached·statusclean·제품diff0 |
| original-preserved.json·original-hashes.json | 원UI/DEV-FIX/OPS216파일 byte 동일 |
| cleanup.json·cleanup-container.log/.exit | 자기container/volume·clone/비밀·listener/context 회수와 공유8container동일 |
| archive-prefix.json | 이전designer완료전문의prefixhash·bytes |
| recheck.txt | 기존QAhelper를사용한이번시나리오 절차. 비밀CFG/메일은 보존하지 않음 |

첫 run exit1을 전체run PASS로 바꾸지 않는다. 화면 판단은 저장된 실관측/PNG와 별도 회수 결과에 한정한다.
명령 자신의 exit를 저장했다. 로그 tail의 종료코드를 검사 통과로 사용하지 않았다.
최종 archive SHA의 FullOps/strict/diff/push 결과는 레포밖 기록과 worker_done으로 고정한다.

## 한계와 후속

원09c UX06 PASS는 원래 SHA·조건·manifest로만 재사용한다. 새 전체UX06 실행은 아니다.
HTTP/store 공통경계는 변경됐으므로 전체의존성동일을 주장하지 않는다. 전체기능 QA는 TESTER다.
원09c UX07 FAIL·405·effectiveviewport980·OPS H1/M1/checkexit1은 그대로다.
운영/실메일/실24h/노우↔다닷/부하/복원/전체API/스크린리더·사람로그인은 미검증이다.
coor는 고정dfc의독립TESTER/OPS와좁은UI결과·마지막SHA검사를확인한뒤main판정을진행한다.

증거 검사 최초 exit1은 desktop의 scrollWidth1265를 viewport1280과 같다고 요구한 검사 오류다. 수직 scrollbar 폭을 반영해 `scrollWidth <= viewport`로 확인했다. mobile은 실제390/390을 계속 요구했다. `integrity-initial.json`과 후속 `prearchive-integrity.json`을 보존하며 제품 실패로 바꾸지 않는다.

## Archive 뒤 추가 인계

coordinator-followup.json의msg_23b8de8f85ac는OPS msg_d5687c98688e/c219의H2high를전달한다.
유효credential의타대상/no-lease cleanup flood와유효정리18/18의429는OPS관측이며designer재현이아니다.
이번dfc시각조건PASS는유지하되제품전체수락/main통합은차단한다. 동일DEV후속수정과새fixed독립QA/OPS가필요하다.
첫archive SHA3b31b8f의lint/test/strict/diff/productdiff exit0·ERROR0/WARNING2를record-3b31b8f-*로보존한다.
추가인계기록의마지막SHA도검사한다. 원래전문archive와빈인박스는유지하며다시finish하지않는다.
