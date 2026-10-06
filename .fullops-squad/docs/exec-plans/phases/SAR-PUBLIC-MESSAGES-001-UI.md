---
title: SAR-PUBLIC-MESSAGES-001-UI — 직접 검수 실행과 인계
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-MESSAGES-001-UI]
summary: 실제 일반회원 UI 검수·초기 실패·자기 회수·고정 기록 검사와 인계를 기록한다
---

# SAR-PUBLIC-MESSAGES-001-UI — 직접 검수 실행과 인계

## 범위와 결과

Task는 `task_f37e47dda648`, Dispatch는 `ctx_b2c603e79536`이다.
정규 `handovers/to_designer.md`를 기준으로 fixed `09c523da8a3407288d9f5d711e1834af12bc7808`을 직접 검수했다.
적용 규약·제품/UX 정본·실제 화면·격리 fixture·증거·결함은 [보고서](../../design-docs/mockups/SAR-PUBLIC-MESSAGES-001-UI.md)에 있다.
UX06은 로컬 시각 조건 PASS다. UX07은 Deny 뒤405의 F-UI-MSG-01 medium으로 FAIL/보류다.
OPS H-1 high·M-1 medium 발견의 coor follow-up `msg_cfbbb3035700`을 읽고 원본 검수를 계속했다.
coor에게 `msg_6194f4825058`로 Deny 결과 결함과 DEV/재검수 후속을 전달했다.

## 자기 실행과 보존

별도 detached clone·자기 Postgres17·기존 SMTP sink·변경 없는 host relay·실제 Chrome CDP로 일반회원 두 명을 실제 코드 로그인했다.
회원 cookie로 실제 receipt/gate를 열고 PNG를 직접 확인했다. 제품 코드를 추가·변경하지 않았다.
자기 QA 기록 스크립트만 작성했다. 제품/UX 정본·기술 정본·타인 인박스·기존 실패는 변경하지 않았다.
Go fixture 작성 hook 거부는 원문을 보존했다. 기존 relay와 SQL fixture로 해결했으므로 코드 수정 권한 예외를 요청하지 않았다.
실행1–5의 fixture/도구 실패 exit1, 실행6의 exit0, supplement 첫 exit1/후속 exit0를 보존했다.
정규 expired의 원래 accepted_at을 보정한 후속 캡처와 음수TTL 초기fixture2개를 구분했다.
첫 증거 validation은 desktop의 scrollWidth1265와 viewport1280의 일치를 요구해 Python exit1이었다. 스크롤바 차이를 포함한 overflow 없음 조건으로 고쳤고 후속 exit0이다. 셸이 후속 strict/stamp를 계속해 exit0이 된 사실도 validation-initial.json에 기록했다.
DEV 증거45파일과 제품/규약 정본의 내용 동일성을 확인했다. 원본OPS finding과 c6/Jev 오기는 보존한다.

자원 회수 명령은 exit0이다. 자기 container/volume·listener·Chrome context/profile·메일·비밀 설정·clone을 회수했다.
공유8개 컨테이너의 ID/이름은 회수 전후 같다. 캡처49개 중47개가 유효 UI 조건이고2개는 초기 fixture 오류다.
영상은 필요하지 않았다. 시각 증거는 manifest의 조건과 고정 SHA로 연결한다.

## 기록 검사와 완료 절차

보고서·D04 연결·QA 증거·designer context·PLANS·현재 인박스와 완료 로그를 커밋한다.
깨끗한 기록 HEAD에서 FullOps lint.py `--from 09c523da8a3407288d9f5d711e1834af12bc7808`를 실행한다.
등록 명령 product-lint/product-test, strict, ref-to-HEAD diff 검사의 원래 종료코드를 각각 기록한다.
기존 PLANS의 SIZE 경고는 누적 이력 보존 때문이다. 자기 증거를 지워 경고를 숨기지 않는다.
추가 규모는 요구된 상태별 두viewport·보존로그·실행 재현·manifest다. 제품·의존성 추가는0이다.
별도 디자인 lint·theme·Tailwind/shadcn은 없다. Go 문자열 CSS의 직접 판정은 자동 정규식 검사로 대체하지 않았다.

완료 전문을 인박스에 쓴 뒤 work.py finish로 날짜별 로그에 보존하고 인박스를 비운다.
마지막 커밋의 필수검사·일반 origin/fullops/designer push·원격 SHA 일치를 확인한다.
마지막 검사 원본은 레포 밖 자기 검증 폴더에 고정하고 worker_done의 fullSHA·원본 경로와 연결한다.
coor는 이 마지막 원본을 통합 증거에 보존한다. 이전 기록 SHA의 검사로 마지막 SHA 통과를 주장하지 않는다.

## 후속

DEV는 F-UI-MSG-01과 별도OPS H-1/M-1의 수정 방법을 결정한다.
designer는 수정된 고정 후보의 Deny 직후 상태·홈 복귀와 변경 영향 UI만 새 과제에서 재검수한다.
현재 원본 검수의 기록 완료는 수정 대기 때문에 인박스를 계속 점유할 이유가 아니다.
실메일·운영 공개·실24h·실제 노우↔다닷·부하/복원·전체PS13/14는 기존 담당과 재개 조건을 유지한다.

최초 staged 공백 검사는 Deny405 텍스트 snapshot2개의 EOF 빈 줄로 exit2였다. 원문 로그·exit·원래 텍스트를 보존하고 EOF만 정규화했다. PNG·제품 동작·판정은 바꾸지 않았다.
