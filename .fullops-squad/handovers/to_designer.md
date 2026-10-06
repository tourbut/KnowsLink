---
title: POLICY 추가 기록 — 기록 보호의 교체와 목록 정리 제품 조건
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-AGENTS-001-POLICY, SAR-PUBLIC-AGENTS-001-POLICY-SUPPLEMENT]
summary: 같은 POLICY의 추가 handoff와 제품 답 및 보존 예외를 기록한다
---

# SAR-PUBLIC-AGENTS-001-POLICY-SUPPLEMENT — 같은 POLICY의 기록 보호 사용자 조건을 추가 보존한다

- 원래 제품 과제: SAR-PUBLIC-AGENTS-001-POLICY. 새 제품 배정이 아니라 같은 Dispatch의 추가 기록이다.
- Task: task_9ff87558884b. Dispatch: ctx_94d86ca2989c. 복귀 terminal: term_6895aaf1-7b43-4fe0-a416-76f1255a5946.
- 기준 ref: d1651784c4338efeb0d6141467d563c6b354e4a5. 추가 판단 전 기록 SHA: 3aa073bff7e8cb0966aec721a4b604aa4d48640c.
- 원래 지시서·결과 전문은 `.fullops-squad/handovers/logs/2026-10-06_to_designer.md`의 원래 POLICY 항목에 보존했다. 원래 제품 답은 변경하지 않는다.
- 적용 기준: fullops-common-0.3.3·FULLOPS·project·document-writing·D02 PS-05/06/07/11·UX-04/05. 제품 판단과 자기 문서만 변경한다.
- coordinator는 preamble ask 답변으로 현재 역할 인박스의 추가분만 SUPPLEMENT 기록 키로 finish하는 방식을 허용했다. 원 Task/Dispatch는 유지한다.

## 추가 지시와 완료 조건


2026-10-06 최종 worker_done 전 확인에서 같은 과제의 추가 handoff를 받았다. 최종 DEV-FIX SHA는 `4a1b80aec8fa6a06144d51f3a5609927a2644928`이다. DEV 실행 기록은 `/home/shin/orca/workspaces/KnowsLink/fullops-dev/.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-AGENTS-001-DEV-FIX.md`에서 읽기만 한다. 원래 제품 답과 완료 기록을 유지하고 이 인박스에서 후속 판단을 추가한다.

- [x] owner agent 기록 포화와 살아있는 agent의 키 기록 포화에서 사용자에게 보이는 상태·다음 동작을 판단한다. 기술 보호값은 DEV 소유이며 상품 quota로 표시하지 않는다.
- [x] 새 agent 교체의 별도 연결·새 식별자·명시적 관계 수락과 철회 agent의 최소 24h 뒤 목록 정리 표시를 D02/UX 및 DEV/QA 관찰 조건에 반영한다.
- [ ] 원본 기록과 제품 답을 보존하고 추가 완료 전문·검증·인박스 비움의 보존 절차를 기록한다. worker_done은 아직 보내지 않았다.


## 완료 보고

### 추가 제품 답 전문 — 기록 보호의 사용자 동작


DEV-FIX의 기술 보호값 owner agent 기록 10개·agent 키 기록 20개는 새 상품 quota가 아니다. 기존 제품 활성 agent 5개·활성 키 3개를 그대로 적용한다. 기술값·API·정리 방법과 실제 CPU/DB/복원 측정은 DEV/OPS 책임이다. 입력 DEV 실행 기록이 최종 고정 4a1b80aec8fa6a06144d51f3a5609927a2644928의 파일과 byte 동일함을 확인했다. 이 확인은 제품 동작 PASS가 아니다.

키 기록 포화는 해당 agent의 새 키 연결 실패와 새 agent 교체를 안내한다. 살아 있는 agent의 철회 키 삭제나 기존 kid 재할당으로 우회하지 않는다. 키 철회나 대기만으로 새 공간이 생긴다고 안내하지 않는다. 새 생성 실패만으로 기존 활성 자격을 임의 철회하지 않는다. 사용자가 명시적으로 철회하면 해당 키/관계를 종료한다.

교체는 새 식별자·정상 연결·owner 확인을 요구한다. 기존 agent의 이름·키·관계·승인·receipt 권한을 승계하지 않는다. 새 agent와 각 상대 사이에 새 초대·수신 owner의 명시적 수락이 필요하다. 동일 owner의 두 agent도 새 수락이 필요하다. 새 관계 pending에서 메시지는 계속 거부한다.

새 agent 생성에 활성/기록 여유가 없으면 성공을 약속하지 않는다. 필요한 기존 agent를 owner가 선택해 철회할 수 있다. 철회 전에 키·관계 종료를 안내한다. 활성 slot 해제와 보존 기록 slot 해제는 구분한다. owner 기록 포화는 철회 기록의 최소 24h 보존·실제 정리 뒤 수동 재시도를 안내한다. 포화 중에도 철회·취소·거절·unpair는 기존 정리 budget 안에서 가능하다. 정리 rate·현재 권한 실패는 안전하게 거부한다.

철회 agent는 보존 중 철회 상태를 표시한다. 최소 24h 뒤 실제 정리되면 홈 목록에서 사라질 수 있음을 안내한다. 정확히 24h에 공간이 생긴다고 약속하지 않는다. 반복 철회·재로그인이 기간을 줄이거나 권한을 복구하지 않는다. 목록에서 사라짐을 철회 복구·일시 오류·영구 개인정보/백업 삭제로 표시하지 않는다. 새 연결·새 관계 수락 전 전달을 차단한다.

D02 PS-05/06/11 상세 절과 UX-04/05에 상태·관리 복귀·새 생성·수동 재시도·철회 안내를 썼다. 추가 DEV/QA 표는 키/agent 기록 경계·포화 중 정리·교체·최소 24h 전후·옛 권한/관계 차단·재시작·실제 오류와 다음 동작을 포함한다. designer 직접 재검수와 OPS 실제 자원·복원 보호는 고정 후보 후속이다. 기존 UI FAIL·원문·증거·첫 완료 로그는 유지한다. 새 답이 이미 구현됐거나 독립 PASS라고 소급하지 않는다.

### 보존 예외와 검증

첫 POLICY의 work.py finish는 exit 0이며 원래 인박스/지시서/완료 전문을 그대로 로그에 보존했다. 최종 check에서 추가 handoff를 받았다. 원 키로 두 번째 finish를 시도한 결과 exit 1, `이미 아카이브된 과제입니다. 기록을 확인하세요`였다. 이 중복 보호 거부는 원래 로그나 인박스를 삭제하지 않았다. coordinator의 ask 답변은 같은 Dispatch의 추가 기록에 SUPPLEMENT 키를 허용했다. 원래 지시서/제품 답/완료로그를 재작성하지 않고 추가 지시와 제품 답만 현재 역할 인박스로 기록했다.

문서 strict는 추가 판단 뒤에도 exit 0(13개·문제 0·경고 0)이다. 제품 코드·기술 정본·DEV 체크아웃은 변경하지 않았다. 새 판단의 독립 QA·직접 시각·운영 측정은 미실행이다. 추가 문서의 깨끗한 기록 HEAD lint와 최종 finish/빈 인박스 결과는 완료 전에 기록한다.
