---
title: 관계 초대 반복과 종료 후 재초대의 제품 판단 및 DEV·QA 인계
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-AGENTS-001-POLICY]
summary: L2의 제품 해석과 현재 권한·새 세대 및 독립 관찰 조건을 기록한다
---

# SAR-PUBLIC-AGENTS-001-POLICY — 관계 초대의 제품 판단과 기술 인계

## 범위와 판단

designer는 원본 AGENTS 보안 리뷰 L2의 모호한 제품 기준을 결정했다. 제품 규칙 결정은 완료했다. 구현 준수·독립 QA·직접 시각 수락·일반 서비스 공개 수락은 후속이다.

기준 ref는 `d1651784c4338efeb0d6141467d563c6b354e4a5`다. 기록 시작 HEAD는 `53aab5d7acf518ca02aad118e07bd1167ab8af8b`다. 원본 제품 관측 대상은 `d1eef9bb90b9726149980320c42fb1fdbcaf584a`다. 적용 기준은 `fullops-common-0.3.3`·FULLOPS·project·document-writing·designer context다. 예외는 없다. Jev 후보·분류는 [기존 context](../../evaluations/jev/SAR-PUBLIC-AGENTS-001-POLICY-context.json)를 재사용했다. D09/D10 추천은 DEV 소유여서 수정하지 않았다.

제품 답은 다음과 같다. 유효한 pending 또는 active의 반복 초대는 새 초대를 만들지 않는다. 거절·만료·양측 중 어느 쪽 철회 뒤에는 새로운 수동 초대를 허용한다. 새 수신 owner가 명시적으로 수락하기 전에는 메시지를 거부한다. 이전 종료 세대의 권한은 복구하지 않는다. 새 수락은 이전과 구별되는 관계 세대를 사용한다. 자동 재초대·자동 복구는 금지한다. 기존 제품 수치와 현재 권한 경계는 유지한다.

## 원문과 근거

- [원본 리뷰 L2](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-REVIEW-review/report.md#l2--low-제품-판단-필요-거절만료-뒤-같은-대상에게-즉시-재초대할-수-있다)는 B 거절 뒤 A의 새 pending·Generation 증가·B 차단 수단 부재를 관측했다. PS의 “반복 초대는 새 초대를 만들지 않는다”의 범위가 불명확했다.
- [일반 서비스 D02](../../planning/product-specs/SAR-PUBLIC-SERVICE.md)의 PS-07은 수신 owner 결정·양측 철회·재수락의 새 세대를 요구한다. PS-11은 송신 pending·rate·전체 자원 제한과 안전 정리 분리를 요구한다. 거절을 영구 차단으로 정한 기존 요구는 없다.
- [frozen 원천](../../planning/sources/silent-agent-relay/protocol.md)의 Security 절은 현재 key/pair/owner 재검사와 unpair 후 새 관계 세대를 요구한다. C5의 `Auto reinvite = current active pair only`는 현재 active 반복에 한정한다. 최신 일반 서비스 D02가 종료 뒤 수동 재초대의 제품 해석을 정한다. frozen의 보안 경계·wire·보존 조건은 대체하지 않는다.
- [원본 UI 관측](../../design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI.md)은 UX04–05 FAIL/보류다. F-UI-01 모바일 지문 overflow와 F-UI-02 오류 뒤 복귀 부재는 미해결 medium이다. low F-UI-03/04·blank PNG·metadata 검사 원실패·보정 근거·최종 통과 기록도 보존한다.

거절은 해당 초대에 대한 거부다. 영구 차단과 같은 추가 기능으로 확대하지 않는다. 기한 만료는 관계 권한을 부여하지 않는 종료다. 철회는 현재 권한을 즉시 끝내며 새 초대를 통한 새 동의까지 금지한다는 뜻은 아니다. 재초대를 금지하면 정상적인 재연결까지 막으므로 기존 수락 흐름과 한도로 처리한다.

수신 owner가 반복 초대를 받을 위험은 남는다. 기존 제한은 자원 증가율을 줄이지만 수신 owner별 괴롭힘 차단을 보장하지 않는다. 이번 L2의 모호성만 확정하며 남용이 해결됐다고 선언하지 않는다. 차단·쿨다운·수신 pending 한도·새 상품 수치는 추가하지 않는다. PS-11/12의 공개 전 남용·보존량·공유 자원 검증을 유지한다. 실제 남용 또는 기존 방어 실패 근거가 생기면 coor가 designer에 별도 제품 변경을 인계한다. 별도 수치를 미확정 필수 조건으로 새로 만들지 않는다.

## 변경과 소유권

D02의 PS-07 상세 규칙과 관계·초대 기본값 행을 명확하게 했다. 수치는 바꾸지 않았다. UX-05에는 반복 대기·종료·수동 재초대·새 수락 대기·안전한 오류의 다음 동작을 반영했다. 원천 front matter는 `deliverables.py --stamp`로 갱신한다. 문서 상태는 기존 review/draft를 유지한다.

DEV-FIX는 기존 M1·L1·L3와 UI medium 두 건의 기술 구현을 계속 맡는다. 이 답은 coor가 DEV-FIX에 그대로 전달한다. API·함수·DB·세대 발급 시점·거부 표현·정리 구현은 DEV가 정한다. D03/D05–10·코드·board·타 역할 기록은 수정하지 않는다.

## DEV/QA 관찰 조건

아래 기대값은 새 제품 기준이다. 과거 d1eef9b가 통과했다고 주장하지 않는다. DEV는 수정 후보에서 자동 검사를 남긴다. TESTER는 안정된 고정 후보에서 독립 상태·인가 QA를 수행한다. 후보 SHA·실제 상태·허용/거부·각 명령 종료코드를 남긴다. 시간 경계는 격리 fixture의 제어된 시각으로 확인할 수 있다. 실제 24h 대기나 운영 계정 사용은 요구하지 않는다.

| 조건 | 허용/거부와 관찰값 |
|---|---|
| 유효 pending의 같은 방향 반복 | 기존 pending만 유지한다. pending 수·세대·기한이 바뀌지 않는다. 메시지는 계속 거부한다. 반복 요청도 기존 신규 rate에 집계한다. |
| 유효 pending의 반대 방향 초대 | 별도 pending·세대·기한 연장·자동 수락이 없다. 현재 수신 owner만 해당 초대를 결정한다. |
| active의 반복 초대 | 현재 active·세대·active 수를 유지한다. 추가 pending·권한을 만들지 않는다. 현재 권한이 유효한 기존 메시지 경로만 허용한다. |
| 수신 owner 거절 뒤 | 기존 pending이 종료되고 pending slot이 해제된다. 해당 초대 수락과 메시지는 거부한다. 기존 한도를 만족한 수동 재초대는 새 pending이다. 새 수락 전 메시지는 거부한다. |
| pending 24h 만료 경계 | 기한 전에는 현재 수신 owner가 수락할 수 있다. 기한 도달 이후에는 옛 초대 수락을 거부한다. 만료 pending은 용량을 계속 점유하지 않는다. 새 수동 초대만 새 기한의 pending을 만든다. |
| 발신 owner 철회 / 수신 owner 철회 | 두 방향을 각각 검증한다. 현재 pending 또는 active를 종료한다. 기존 메시지 권한과 옛 수락은 재사용할 수 없다. 수동 재초대 뒤 새 수신 owner의 수락이 필요하다. |
| 새 수동 초대·재수락 | 현재 agent·owner 권한과 모든 기존 한도를 검사한다. active 전환은 한 관계만 만든다. 종료 전 세대와 구별되는 새 세대다. 옛 결정·승인·메시지 권한이 새 세대를 활성화하거나 새 전달 권한으로 사용되지 않는다. |
| 중복·동시 결정·늦게 도착한 결정 | 같은 pending의 동시 수락은 한 active만 만든다. 종료된 초대의 늦은 수락/거절은 새 pending/active를 변경하지 못한다. 거절·철회와 수락의 경합은 실제 확정 순서와 현재 권한으로 판정한다. |
| 동일 owner의 두 agent / 독립 owner | 두 경우 모두 첫 수락과 재수락을 명시적으로 수행한다. 타 owner가 초대 결정을 바꿔치기하면 거부한다. owner가 같다는 이유로 active를 자동 생성하지 않는다. |
| 한도 경계·포화 | 송신 pending 10·전체 pending 200·owner active 20·전체 active 400·기존 rate 등 필요한 모든 한도를 집행한다. 상한의 다음 시도는 새 자원·권한·세대를 만들지 않는다. 신규 포화 중 거절·철회는 분리된 정리 budget 안에서 가능하다. 정리 budget 초과도 성공으로 표시하지 않는다. |
| 재시작·재로그인·옛 자격 | 종료 상태와 새 세대 경계가 재시작 뒤 유지된다. 재로그인은 자동 복구하지 않는다. agent 자체가 철회되었거나 owner가 비활성이면 관계 재초대로 자격을 되살릴 수 없다. |

receipt-only replay는 frozen의 기존 receipt 반환 규칙을 유지한다. 위 옛 메시지 권한 차단은 기존 receipt 반환까지 새 전달로 오해하게 만들지 않는다. 새 enqueue·lease/ACK·승인 소비·실행·결과 공개에는 현재 권한을 다시 검사한다. 이미 확정된 공개를 소급 취소한다고 약속하지 않는다. 철회 metadata 최소 24h와 기존 보존량 보호·정리 책임을 유지한다.

designer는 DEV 수정 후보의 실제 관계 화면을 직접 확인한다. pending 기한과 반복 안내·active 유지·종료 안내·수동 재초대 뒤 새 수락 대기·오류 뒤 다음 동작이 관찰 대상이다. 내부 세대의 기술 표시는 요구하지 않는다. 원본 F-UI-01/02 수정과 관련 low의 변경 영향만 재검수한다. 변경 없는 증거는 원래 SHA·조건과 의존성 동일성을 확인한 뒤 재사용한다.

## 미결정과 후속

이번 제품 질문에 미결정 항목은 없다. 새 차단·쿨다운·수신 한도는 미구현 필수 기능으로 남기지 않는다. 기술 설계·후보 준수 여부·독립 QA·직접 화면 재검수는 DEV/TESTER/designer 후속이다. coor는 고정 후보와 필수 근거가 준비되면 각 역할의 빈 인박스에서 재개한다.

원본 UI FAIL/보류·M1 공개 차단·미해결 critical/high 차단과 PS-13 전체 수락은 유지한다. 실제 일반 이메일·운영 공개·노우↔다닷·실메시지·Workers Free와 기존 서버/Tunnel 제한은 확대하지 않는다. 실제메일·배포·외부 계정·과금·서버/Tunnel 변경은 수행하지 않았다.

## 검증 기록

문서 검사·깨끗한 기록 HEAD의 FullOps lint 결과는 아래에 완료 시 기록한다. 정책 동작 QA와 실제 화면 재검수는 이번 문서 과제에서 실행하지 않는다.

- `git diff --check`: exit 0. 문서 변경의 공백을 검사했다.
- `deliverables.py --repo . --strict`: exit 0. 검사 13·미작성 0·문제 0·경고 0이다.
- 정책 동작·실제 브라우저·실제 이메일·운영 검증은 미실행이다. 문서 변경 과제이며 DEV 구현과 고정 후보의 독립 검증이 선행해야 한다. 기존 관측을 새 PASS로 바꾸지 않는다.

- 깨끗한 기록 HEAD `1de5de129697d138b9c5fb59432c30e8b99bfb4d`에서 `lint.py --repo . --from d1651784c4338efeb0d6141467d563c6b354e4a5 --out <레포 밖 JSON>`은 exit 0이다. ERROR 0·WARNING 1·실행 불가 0이며 등록 product-lint(`make lint`)·product-test(`make test`) 모두 exit 0이다. 원본 JSON은 [record-lint.json](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-POLICY/record-lint.json)에 보존했다.
- WARNING은 SIZE-001의 기존 누적 PLANS 842줄(기준 ref 819줄, 상한 500줄)이다. 이번 designer는 자기 결과만 추가했다. 기존 통합·실패 기록을 임의 삭제하거나 타 역할 소유 문서를 분할하지 않는다. 경고는 coor의 별도 정리 판단 대상으로 유지한다. SIZE-002·DEP-001은 없다. 제품 코드·의존성 변경은 0이다.
- 기록 시작 HEAD 대비 문서 파일만 바뀌었다. 제품 기본값 표의 값 열은 전후 동일하다. 원본 리뷰·UI 보고서 byte 동일과 UI 증거 폴더 diff 0을 확인했다. 새 원천 문서의 로컬 링크도 존재한다. 이 정적 대조는 exit 0이다.
- 완료 전문은 `work.py finish --repo . --role designer --key SAR-PUBLIC-AGENTS-001-POLICY`로 보존한다. finish 뒤 빈 인박스와 로그 전문 일치를 확인하고 커밋한다. 최종 깨끗한 HEAD의 동일 lint 재검사 결과와 고정 SHA는 worker_done에 기록한다. 새 정책 동작이 검증됐다는 뜻은 아니다.
- `work.py finish`는 exit 0이다. [완료 로그](../../../handovers/logs/2026-10-06_to_designer.md)에 지시서·제품 답·완료 보고 전문을 보존했다. 보존 전 인박스 전문과 로그 추가분의 일치 및 빈 인박스를 확인했다(exit 0).
