---
title: 잔여 UI 이슈 설계 독립 리뷰
status: review
updated: 2026-10-15
owner: coor
tasks: [SAR-UX-ISSUES-002-DESIGN-REVIEW]
summary: 고정 문서 범위의 요구·안전 경계·부분 탐색·미실행과 순차 인계를 검토한다
---

# SAR-UX-ISSUES-002-DESIGN 독립 설계 리뷰

문서 범위 수락이다. Critical 0, High 0, Medium 0, Low 0이다. 이 판단은 검토자의 원문·코드 대조 결과이며 OCR 자동 판단이 아니다. 제품 구현, 전체 UI 수락, 신규 설치 지원, GitHub 이슈 종료, main 통합, 운영 배포를 수락한 결과가 아니다.

## 고정 범위와 독립성

- 기준/merge-base: `76ec89fe66cbeea5c2ffb33ff07015fd722029b3`.
- 대상: `11ea298ef63cb90e86881d08fa83c44c9e3f75c5`.
- 작성자 실제 세션: `01a13f95-f910-77b3-a343-090f579de67f`.
- 검토자 실제 `session_meta.id`: `01a13faa-5f90-79e3-8b5e-d23b1641f491`. 자기 rollout 첫 session_meta에서 확인했다. 부모 session_id와 구분했다. CLI는 Codex다.
- 공식 snapshot: `C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-0283e5629d4d4b18b8b6766f1203f8d3`. 고정 SHA의 별도 detached checkout이며 source는 읽기 전용으로 유지했다. lint는 Git metadata와 외부 결과에만 기록했다. OS 쓰기 격리라고 주장하지 않는다.
- 기록 소유권은 이 리뷰 폴더다. snapshot 정리 담당은 coor다. 삭제·이동·commit·GitHub·실계정·배포는 수행하지 않았다.

공통 기준은 `fullops-common-0.3.3`이다. FULLOPS.md, rules/common/README.md와 coding-style/testing/security, project.md, document-writing.md, review/rule.json을 읽었다. coor에서 읽은 이 기준과 대상 SHA 사이 diff는 0이다. FullOps 1.4.4 fullops-review와 open-code-review-delegate를 적용했다. OCR v1.13.0 preview/rule만 사용했으며 LLM endpoint를 호출하지 않았다.

원 지시 base458eb565와 착수HEAD76ec89fe를 구분했다. 이 리뷰는 착수 뒤 독립 문서 변경22파일만 다룬다. 제품·의존성·검사 설정·배포 입력은 이 범위에서 변경되지 않았다. 기존 CSS 결함과 후속DEV fc4ad58 해결은 별도 리뷰다.

## 요구와 설계 판단

| 항목 | 직접 대조와 판정 |
|---|---|
| #2 이름·공유ID | 원문은 handle 지원 여부의 제품 결정을 요구한다. D02 HP01/02/03과 최신 지시가 사람 코드 발견을 우선하므로 handle·초대링크 미도입은 허용된 최소 결정이다. 이름/코드/wire ID, 동명 안정 구분값, Unicode1–40·제어문자 거절, 이름 권한 fallback을 구분한다. 비공개 상대 프로필 추가 조회를 금지한다. |
| #4 친구 | HP03/07·legacy 이행에 맞게 사람 수락·현재/향후 활성agent 범위를 유지한다. 옛 개별agent 초대·대상 미리보기·same-owner 매번 수락을 대체하는 이유가 있다. 유효/없는/차단/상대 한도를 노출하지 않는다. 빈agent·중복/기한/세대/코드교체 검증이 있다. |
| #6 조회 | memberReceipt의 refusal과 hidden 재조회 한계를 소스로 확인했다. 성공/오류 모두 같은 편집폼, 비공개 실패 안내, 과대입력 거절·이스케이프·계정교대 복원금지, 429/503/만료·철회 조건이 구체적이다. 조회를 온라인/발송 기능으로 확대하지 않는다. |
| #5 승인 | verifiedGate/gateView와 human-gate를 읽었다. 요약은 검증된 원요청 구조화 필드만 사용한다. 자유문/hint는 권한이 아니며 원문·digest·현재owner/key/세대·기한을 보존한다. 인간게이트와 공개/업무실행을 분리하고 무원문/미지원 유형은 거절한다. |
| #7 구조·접근성 | remote 동의/return의 현재 skip/header/main 존재를 반영한다. 정확한 복귀·현재위치·오류/상태를 요구한다. native200%/키보드/낭독/다른브라우저 미실행을 완료로 바꾸지 않는다. D04토큰·44px·대비·320/390/1440 기준을 유지한다. |
| #3 언어 | Member에 locale이 없는 현재 상태와 구분한다. ko/en 계정저장·익명명시선택·B계정 우선·A선호 복사금지·저장실패/차단/손상 fallback이 있다. 누락번역은 완료차단이다. 이름·원문·ID·명령·기한을 번역하지 않는다. KST와 html lang 정책이 있다. |
| #8 연결·재개 | connections의 waiting/prepared/approved/consumed를 대조했다. 외부설치/검색/실제화면 상태는 확인필요로 남긴다. 응답유실 시 기존agent/key와 앱저장부터 대조하며 새등록·덮어쓰기·자동재발송을 요구하지 않는다. 구체 저장/API는 DEV책임이다. |
| 지원근거 | 10-12 LIVE,10-13 NOTIFY,10-14 OPS 원보고를 읽었다. Grok/Dots 기존왕복·Dots 초기만료/지연·후속보존을 깨끗한 신규설치로 승격하지 않는다. Meta Muse/Hermes/OpenClaw는 제품/버전/인증/권한 미확인이다. 새 API/CLI문법을 제안하지 않아 외부docs 조회는 해당없다. |
| 순차완료 | 조회→이름/친구→승인→연결→언어 순서로 각 단위DEV·tester·designer·fixed SHA리뷰를 요구한다. #3 P1은 유지하며 중간단계로 닫지 않는다. #8 안내완료와 지원완료를 구분한다. |

불필요한 새 라이브러리·프레임워크·JS·공개검색·범용온보딩 엔진은 추가하지 않는다. 좁은 제품결정은 designer권한이다. D02/D04/D11 직접 갱신은 최신 지시가 제한했으므로 후속담당·변경원천을 새 설계에 남긴 방식은 타당하다.

## 커버리지와 탐색

[result.json](result.json)에 모든(path,status)의 개별 근거가 있다. 총22파일, reviewed14, skipped8, 직접검토율63.64%다. [coverage.json](coverage.json)과 [verification.json](verification.json)에 목록·원본byte hash를 남겼다. skipped는 context/find/documents-find/packet 원시탐색4파일, packet-screening, preserved-inputs, OCR제외 원시로그2파일이다. 구조·기준·요약·존재는 확인했지만 모든 후보본문이나 원실행을 재검증했다고 표시하지 않았다. 필수설계·이슈원문·QA·phase·context변경·route/outcomes·완료아카이브 실질지시는 직접 검토했다.

packet은386경로, unknown19, remaining349, optional0이다. outcomes403 path/category가 전체packet과 정확히 대응하며 누락/중복0이다. completed29/no_change374다. no_change357은 본문전체 미검토를 명시하고 document_update17은 소유제한이다. D04전송실패는 author의 직접 현행절읽기로 보완했으며 설계의 실제 적용내용을 대조했다. 코드후보 밖 verifiedGate·Member·connections도 보완했다. 349개 과거문서를 모두 의미검토했다는 주장은 없다.

## 대화 미참조 인계 점검

산출물 인덱스 docs/deliverables/README.md에서 시작했다.

| 항목 | 원천/절·결과·한계 |
|---|---|
| 현재요구·결정이유 | D02폴더→SAR-HUMAN-PEER HP01–16 확인. 새설계 적용범위·추적표·U2/U4에서 최신요구 우선과 미도입이유 확인. |
| 구조·구현/미완료 | D04매핑→service-member-ui 현행사람절·새설계 화면배치 확인. 새설계 첫절/추적표·QA에서 source확인과 미구현/미실행 구분 확인. |
| 실행·검증 | project.md 등록검사·U2–U8 완료조건·순차표의 담당/조건 확인. author strict13·링크17·보호152 원보고/JSON을 재사용하며 반복하지 않음. |
| 운영·복구 | U8 기존자격보존/만료/응답유실 재개 확인. 실제운영절차 변경은 해당없음. OPS원기록을 새 수락근거로 쓰지 않음. |
| 다음작업 | 순차DEV/QA표·own phase/context 링크에서 확인. 기존D04/인덱스에서 새문서 직접링크는 소유제한으로 아직 없음. coor가 정규inbox에 새설계를 연결하고 후속designer 원천반영을 배정해야 함. |
| 링크·정본접근 | author 로컬링크17 검사결과와 핵심 D02/D04/보고 경로를 확인. 전체앵커를 재검사하지 않음. 고정커밋과 coor의 본리뷰로 snapshot 정리 뒤에도 증거에 접근 가능. 현재 정리 미실행. |

인덱스의 오래된 초기MVP 설명을 새완료 주장으로 읽지 않았다. 이번 추가설계는 아직 별도 후속문서이므로 인덱스만으로 자동발견할 수 있다는 완전한 인계 주장은 하지 않는다. 명시 소유제한·후속전달 조건이며 이 문서범위 수락을 막는 새결함은 아니다.

## lint·검증·남은 제약

첫 designer lint는 후속inbox 미커밋 때문에 clean gate exit1로 중단됐다. 현재역할 HEAD/inbox를 되돌리지 않았다. coor 후속허가 뒤 clean fixed snapshot에서 공식lint를 실행했다. 원실패는 verification.json에 보존했다.

정확한76ec→11ea lint는 exit1, ERROR0, WARNING0, unavailable2다. Windows PATH에 make가 없어 product-lint/product-test는 실행되지 않았다. [lint-raw.json](lint-raw.json)은 원본byte 사본이고 [lint.json](lint.json)은 unavailable reason만 추가한 prepare입력이다. reason 외 완전히 같다는 assertion을 통과했다. 명령·종료코드·설정·심각도·정책을 바꾸지 않았다. 공식prepare exit0이다.

작성자 final lint는458eb→11ea, ERROR0/WARNING3/unavailable2다. [author-lint-final-original.json](author-lint-final-original.json)에 보존했다. own-base와 달라 prepare입력으로 쓰지 않았다. 초기ERROR1 DOC003은 원QA에 보존돼 있다. 허가msg_c9b3e1c2d111에 따른 author-final-handoff front matter9줄만 추가됐고 본문불변이다. SIZE001은 착수에 이미 있던PLANS/D04, SIZE002는 넓은원base다. 231줄 설계가 이슈7개의 교차경계를 함께 정의하는 규모는 타당하다. DEP001 관련 의존성 추가없음이다.

기존 SAR-LANDING-GUIDE-001-CODE-REVIEW-review/lint-original.json의 fe2387a에는 make lint/test exit0이 있고 원로그 마지막EXIT=0을 읽었다. 제품/DB/browser fixture4795와 finalfe를 구분했다. fe2387a→11ea 및76ec→11ea의 .fullops-squad밖 diff는 각각0이다. 현재Windows환경이 과거실행환경과 같다고 주장하지 않으며 과거PASS를11ea 신규실행으로 쓰지 않는다. 문서전용 수락은 제품입력 불변·원증거 연결에 한정한다. 별도CSS후보의 gate·독립QA는 계속 필요하다.

브라우저·native200%·스크린리더·테마전환·실계정·설치실증은 수행하지 않았다. 기존 실제시각 실패는 별도과제로 처리 중이다. 디자인lint/테마시스템은 project.md상 미구성이며 토큰명세를 실제픽셀PASS로 간주하지 않는다. 새화면 정본과 기존공용 Go template 재사용·허용범위는 설계에 있다.

공식 check는 exit 0, reviewed 14 / skipped 8 / total 22로 통과했다. 마지막 snapshot도 fixed HEAD·detached·clean을 확인했다. 실제 결과는 [check-execution.json](check-execution.json)과 [check.log](check.log)에 남긴다. check는 기록일치·독립snapshot·critical/high·lint상태 검증이며 제품기능 수락이 아니다. 후속DEV는 U2–U8 실패경계와 구현영향에 맞춘 검증을 수행해야 한다.
