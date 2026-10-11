---
title: SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL — tester의 같은 후보 캡처로 OAuth 반환·오류 화면 직접 시각 수락
status: draft
updated: 2026-10-11
owner: designer
tasks: [SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL]
summary: tester의 같은 후보 캡처로 OAuth 반환·오류 화면 직접 시각 수락
attempt: b61277b0559d44d992788a7d19d490cd
base: eac603b6a41d2ce002e698600ac66b8fc67ad762
subagent_level: standard
test_level: lite
---
# SAR-SSE-INSTALL-001-OAUTH-RETURN-VISUAL — OAuth 반환·오류 화면 직접 시각 수락

- Purpose visual acceptance. 상태: tester 합성 캡처와 독립 QA 전달 완료, ready.
- designer C:/Users/shin/orca/workspaces/KnowsLink/designer, fullops/designer. fresh Codex gpt-6.1-sol medium(Jev). test lite/subagent standard.
- 고정 제품 SHA eac603b6a41d2ce002e698600ac66b8fc67ad762. 실제 구현자 e20cc179-625e-43fa-84fd-b431d432b9c6. 소스·제품 기준을 바꾸지 않는다.
- 복귀 coor term_e61d3e14-29e9-4954-943a-4a75707c82de / run_dc9b4e6dd18a. task/dispatch는 preamble을 따른다.

## 요구와 검수 범위

OAuth consent POST303이 Chromium CSP에 막히는 오류를 DEV가 반환 화면으로 고쳤다. 기능 QA/fixed-SHA code review는 tester 담당이다. designer는 같은 통합 후보의 tester 캡처를 직접 보고 반환·재시작 안내가 사용 가능한지 판정한다. 새로운 기획·스타일·화면을 만들지 않는다.

기존 memberStyle/refusal 스타일과 한국어 안내를 유지한다. 별도 마케팅 화면·장식·새 CSS·테마·애니메이션을 추가하지 않는다. 반환 링크가 보이고 클릭 대상임을 알 수 있는지, desktop/mobile에 글자 넘침·가림이 없는지, expired/used 화면에 재시작 경로가 명확한지, 화면에 민감한 callback/code/state가 노출되지 않는지 확인한다. 자동 복귀·토큰 교환 성공은 tester 근거를 연결한다.

## 먼저 읽을 문서·적용 기준

지정 캡처·SHA256는 `docs/evaluations/qa-reports/SAR-SSE-INSTALL-001-OAUTH-RETURN-COOR/visual-inputs.json`에 있다. 반환 desktop/mobile, expired desktop, CSRF/owner-binding mobile의 5장을 직접 연다. 원 근거는 `SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER-test/report.md`와 같은 폴더 screenshots(총19장)다. tester 기록 SHA는 96c7b8236f8335b4e0c7344e09c00e3cb36c0e8b다. F1 medium은 오류 원인별 문구가 같다는 비차단 잔여로 남긴다. F2 low의 링크24px은 실제 캡처와 클릭 성공을 함께 보고 시각 수락 가능 여부를 판단한다.

제품 gate와 tester 최종 기록 gate는 같은 COOR 폴더의 JSON에 보존됐다. `tester-final-product-shim.py`는 tester가 실제 실행한 최종 shim 보존본이다. 자신의 clean HEAD와 고정 eac603b6의 제품 바이트 전체(.fullops-squad 제외)가 같은지 검사한 뒤, 기존 격리 Linux의 같은 eac603b6 clean clone에서 등록 make target만 실행한다. Windows make 부재 시 자기 task-local make.cmd에서 이 보존본을 호출해 사용할 수 있다. 제품 바이트와 lint 정책이 다르면 재사용하지 않는다. 원문을 고치거나 exact Git HEAD 실행이라고 오인하지 않는다. 검사를 실행하면 byte-identical product 검증임을 결과에 적는다.

FULLOPS, project, orca-agents, rules/common README+세 규칙, delegation, document-writing, interface-design, DEV OAuth 복귀 실행 기록을 읽는다. fullops-common-0.3.3, test lite를 적용한다. 제품 정본은 고정 SHA의 remote_auth.go와 기존 memberStyle/refusal이며 캡처 경로·hash는 tester 완료 후 연결한다. 기존 시각 PASS를 새 화면의 수락으로 해석하지 않는다.

## 소유권과 완료 기준

직접 이미지 도구로 desktop/mobile 지정 캡처를 확인하고 본인의 exec plan·QA 보고·contexts/designer만 작성한다. 제품 source·다른 역할 증거·운영/실계정/Google·Grok 조작·삭제/이동/cleanup은 금지한다. 부족한 캡처는 coor에 구체적으로 요청한다. 전체 QA와 이미 완료된 자동 검사를 다시 실행하지 않는다. 영상은 필요 없다.

관찰별 PASS/FAIL/미확인과 실제 대상 SHA·이미지 hash·한계를 기록한다. 이미지가 새 코드에서 실제 생성됐는지 tester manifest와 source를 대조한다. 문서 결과의 등록 lint를 실행하고 Windows make 부재는 원문을 보존해 사유를 적는다. 제품 바이트가 같으므로 원 exact-SHA product gate 재사용을 명시하며 통과를 꾸미지 않는다. packet의 모든 비선택 path/category를 outcomes로 처리하고 work.py finish로 원 지시서·보고를 보존한다. 최종40자리 기록 SHA 하나를 worker_done에 보고한다. main 통합·배포는 부모가 맡는다.

## 완료 보고

최종 기록과 남은 제한을 전문으로 보존한다.
