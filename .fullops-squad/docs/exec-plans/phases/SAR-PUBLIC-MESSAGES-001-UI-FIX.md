---
title: SAR-PUBLIC-MESSAGES-001-UI-FIX — 실행 기록
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-MESSAGES-001-UI-FIX]
summary: 고정 dfc 직접 검수·원실패 보존·자기 회수와 최종 기록 검사를 연결한다
---

# SAR-PUBLIC-MESSAGES-001-UI-FIX — 실행 기록

## 범위와 판단

- 기준 ref/실제 제품 SHA: `dfc70caa748a90614b02d48c78b4651345938339`.
- 기록 시작 HEAD: `359653a`. Task `task_2e87829b57d2`, Dispatch `ctx_6468a05f27bf`, worker `term_218a079f-e0cd-434c-b297-2386eef8bbbf`.
- 복귀: coor `term_6895aaf1-7b43-4fe0-a416-76f1255a5946` / `run_8ca8bc058ab7`.
- 공통 기준: fullops-common-0.3.3·FULLOPS/project·문서작성규칙·designer context·PS04/06/07/08–11·C1–C5·UX06/07.
- Jev context의 cleanup_admission 미전송과 원UI conflict는 직접 읽고 원래09c FAIL임을 구분했다. 필수 원천을 제외하지 않았다.
- 제품 코드·제품/UX 규칙·기술 정본·타인 인박스 변경은0이다. 제품과 코드 수정은 요청 범위 밖이다.

## 실제 검수와 보존

[직접 보고서](../../design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI-FIX.md)와 [QA 근거](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-UI-FIX/)가 결과 정본이다.
별도 detached clean clone·자기 Postgres17·SMTP·실제 일반회원 로그인·자기 서명 부모/claim/H로 수행했다.
Chrome CDP는 실제 button.click·Enter keyDown/keyUp·Tab을 보낸다. 화면은 PNG를 직접 열어 판정했다.
정상 Deny 클릭/Enter 네 조건의 즉시 canonical GET200·denied 저장·비활성·실제 홈 복귀는 PASS다.
owner Deny는 desktop GET200·denied·비활성·실제 owner 복귀의 좁은 PASS다.
42PNG 중20mobile은 실제 viewport/scrollWidth390이다. 원본09c Deny 오류 viewport980를 이 결과로 덮어쓰지 않는다.

CSRF/타 회원403·신규16/정리4 포화429·cleanup rate429·receipt rate429는 안전 오류와 다음 동작을 확인했다.
신규 rate200 포화의 Deny는 denied 저장 뒤 결과 GET429다. L-UI-FIX-1로 immediate200 PASS에서 제외한다.
포화 fixture를 해제한 뒤 실제 홈 클릭·denied 재조회가 가능했다. 실제60s 대기 검사가 아니다.

첫 run은 관측18개/PNG40개 저장 뒤 Chrome profile cleanup race로 exit1이다. 원log/exit를 유지한다.
별도 회수는 프로세스 argv 판별의 초기 오류를 좁혀 자기 profile/listener 회수를 확인했다.
receipt HTTP 후속은 exit0이며2PNG·실제429·홈 클릭을 기록했다. 전체 첫 run을 exit0으로 바꾸지 않았다.
영상은 정지 상태·최종 응답으로 판단 가능해 만들지 않았다. 제품/QA 스크립트를 레포에서 변경하지 않았다.
기존 browser.py/run.py를 임시 QA 폴더에서 경로·SHA만 바꿔 재사용하고 실행 절차를 recheck.txt로 남겼다.

원UI/DEV-FIX/OPS 원천216파일 hash는 동일하다. 기존 designer archive의 전문 prefix도 보존한다.
자기container/volume·relay/SMTP·browser context/profile·메일·비밀설정·clone/binary를 회수했다. 공유8container ID/이름은 같다.
이메일/code/CSRF/연결 입력은 가렸으며 개인키·token·메일은 로그/Git에 남기지 않았다.

## UX06 재사용과 남은 담당

원본09c의 receipt/wire 표시 관련 member/member_receipt/public_text/identity/member_agents/text.ts·의존성은diff0이다.
원래 UX06 시각 관측은 원SHA·조건·원manifest로 재사용한다. dfc의 새로운 전체UX06 PASS가 아니다.
공통 capacity/store/http는 변경됐다. 전체 의존성 동일성을 주장하지 않으며 영향 오류 화면만 새 실행했다.
전체 API권한·경합·rate·다중 인스턴스는 TESTER, H1/M1 delta 보안 수락은 OPS다.
제품 전체 수락·main 병합/공개는 coor가 필수 결과를 모아 판단한다. 원09c FAIL·review check exit1은 유지한다.
실메일·운영PS08/13/14·실24h·노우↔다닷·부하/복원·전체스크린리더·schedule.commit/OS clipboard는 미검증이다.

## 최종 기록 게이트

제품/코드 변경은0이지만 지시서의 명시 조건에 따라 마지막 archive HEAD도 검사한다.
`work.py finish`는 지시서·완료 전문을 보존하고 designer 인박스를 비운다. 이전 archive prefix를 대조한다.
문서/PNG/QA/완료 전문을 커밋한 clean HEAD에서 아래 검사를 실행한다.

1. FullOps lint `--from dfc70caa748a90614b02d48c78b4651345938339`; 등록 product-lint와product-test의 종료코드·kind·HEAD를 보존한다.
2. deliverables `--strict`; 원래stdout·종료코드를 보존한다.
3. `git diff --check dfc70caa748a90614b02d48c78b4651345938339 HEAD`; 패치 검사 종료코드를 보존한다.
4. 일반 `git push origin fullops/designer`; `ls-remote`의 exact SHA와 clean status를 확인한다.

마지막 검사 결과는 커밋밖 JSON/log/exit로 보존한다. 실행 전 성공을 선언하지 않는다.
worker_done에서 최종 fullSHA·실제 검사결과·경고·원본 경로와 coor 영속화 경로를 전달한다.
SIZE 경고가 있으면 누적 PLANS·필수 상태별 검수 기록의 증가 때문인지 확인한다. 원실패를 삭제해 줄이지 않는다.
DEP0·새 의존성0·제품 CSS변경0이다. theme/design lint가 미구성이므로 별도 디자인 검사 PASS를 주장하지 않는다.

증거 검사 최초 exit1은 desktop의 scrollWidth1265를 viewport1280과 같다고 요구한 검사 오류다. 수직 scrollbar 폭을 반영해 `scrollWidth <= viewport`로 확인했다. mobile은 실제390/390을 계속 요구했다. `integrity-initial.json`과 후속 `prearchive-integrity.json`을 보존하며 제품 실패로 바꾸지 않는다.

## Archive 뒤 coor 후속과 마지막 검사

첫 archive HEAD `3b31b8f2463d11f94932fee5609f4b984b7271cd`의 FullOps --fromdfc는exit0이다.
product-lint/product-test 각각exit0·ERROR0/WARNING2/unavailable0이다. strict13종·문제/경고0·diff/productdiff exit0이다.
원본을 QA의 `record-3b31b8f-*`에 보존했다. SIZE-001은 누적PLANS993줄, SIZE-002는추가2669줄이다.
추가줄은 요구된42상태PNG의동반텍스트·보고서·기록이다. 제품코드0·DEP0이며필수상태/원실패삭제나규칙완화로숨기지않는다.

그 뒤 체크포인트에서coor `msg_23b8de8f85ac`를받았다. OPS `msg_d5687c98688e`/c219 결과의H2high가dfc에남아있다.
유효credential의타대상/no-lease cleanup flood가DB입장중로컬정리슬롯을점유해유효정리18/18이429라는인계다.
이수치는coor/OPS원관측이다. designer의자기시각fixture결과나새QA PASS로사용하지않는다.
동일DEV 후속수정·새fixed의독립QA/OPS가필요하며제품전체수락/main통합은차단한다.
현재dfc의정상Deny직접판정과캡처는보존한다. 다음후보는UI의존성동일성/영향만후속확인한다.
원완료전문/빈인박스를다시작성하거나finish하지않고phase/QA/PLANS에추가인계를연결했다.
이추가문서의최종cleanHEAD도필수검사를한번실행한뒤일반push하고worker_done으로회신한다.
