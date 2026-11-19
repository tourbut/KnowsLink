---
title: SAR-PUBLIC-AGENTS-001-UI-FIX — 직접 시각 재검수 실행 기록
status: draft
updated: 2026-10-06
owner: designer
tasks: [SAR-PUBLIC-AGENTS-001-UI-FIX]
summary: 고정 후보의 격리 시각 검수·증거·정리·완료 인계를 기록한다
---

# SAR-PUBLIC-AGENTS-001-UI-FIX — 직접 시각 재검수 실행 기록

## 배정·범위

Task는 `task_33ddf712fa70`, Dispatch는 `ctx_f2dd9d1df05f`다. 현재 지시서는 `handovers/to_designer.md`다. 제품 fixed SHA와 lint 기준은 `458798c2ee15c179edacfd6f94ebb9896d26f411`이다. 기록 시작 HEAD는 `f0b69a97ba2fdcb070111ff95aca0a42ae0e3f2a`다. 적용 기준은 fullops-common-0.3.3·FULLOPS·project·D02 PS-07/PS-07-I·UX-04/05·POLICY 두 표다. 기술 수정과 제품 코드 변경은 0이다.

## 결정·실행

[직접 보고서](../../design-docs/mockups/SAR-PUBLIC-AGENTS-001-UI-FIX.md)에 고정 후보·fixture·필수항목·PNG·판정·마스킹·후속을 기록했다. 좁은 UI 재검수는 PASS다. 원본 d1/d165 UI FAIL과 원본 OPS/TESTER 실행·중간 실패를 그대로 보존했다. 제품 최종 수락이나 공개 수락을 선언하지 않는다.

새 detached clone·Compose Postgres·SMTP sink·host relay에서 실제 이메일 코드 가입과 비영속 두 회원 context를 사용했다. node CLI 두 prepare/complete는 각 exit 0이다. 최근 재확인도 실제 이메일 code를 입력했다. 모바일 지문 width는 390/390이고 select는 18px이다. pending/active 반복·양측 철회·거절·controlled 만료·같은 owner 수락·새 agent 연결·새 pair 대기·기록 포화/철회/정리 뒤 수동 재시도를 실제 화면에서 확인했다.

Orca 새 page의 screenshot 1px와 full screenshot blank는 장애 증거로 보존했다. 기존 설치 Playwright/Chromium fallback의 정상 PNG 62개와 장애 2개를 모두 직접 열었다. 영상은 필요없다. 시간·포화 fixture는 own relay를 멈추고 own DB만 수정했다. 실제 24h/운영 데이터 정리는 하지 않았다. 자원·개인키·합성 메일·fixture password를 모두 회수했다.

새 캡처 두 파일의 이름을 실제 역할과 viewport에 맞췄다. 35번은 새 발신 B의 철회이고 실제 수신 B의 철회는 38–39번이다. 47번은 최초 mobile 캡처다. 별도 desktop은 57번이다. 제품 결함이나 원본 실패 보정이 아니다.

## 증거·검사

[manifest](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/evidence.json), [상태 관찰](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/state-observations.json), [fixture 변경](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/fixture-mutations.txt), [정리](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/cleanup.json), [원본 대조](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/original-preserved.json)를 보존한다. 원본 관련 65개 파일의 SHA256은 동일하다. clean 기록 HEAD의 lint·strict·diffcheck 결과를 완료 전에 기록한다.

예정 범위는 보고서·좁은 시나리오·직접 PNG/manifest·정리/검증 증거다. SIZE-001의 누적 PLANS는 기존 기록을 보존하며 자기 결과만 append한다. SIZE-002가 발생하면 두 viewport·두 POLICY 표의 필수 관측과 완료 로그/검증 JSON이 추가분의 원인이다. 제품 변경을 분리할 과제가 아니므로 증거를 삭제해 줄 수를 줄이지 않는다. DEP 변경 0이고 DESIGN 자동 검사 미지원은 실제 직접 관측으로 보완한다.

D03은 기술 아키텍처여서 보존했다. 실제 D04 화면 원천과 designer UX 원천에 재검수 링크를 추가한다. 원천 상태는 review/draft를 유지하며 stamp로 이번 과제 키를 연결한다.

## 완료·인계

제품 수락은 coor가 고정 TESTER/OPS 결과와 합쳐 결정한다. 실메일·실제 24h·운영 공개/배포·운영 부하/복원·노우↔다닷은 별도 후속이다. worker 완료 전문을 현재 인박스에 기록하고 work.py finish로 로그에 보존한다. 빈 인박스·로그 전문 일치를 확인한 뒤 커밋·role push·최종 clean HEAD lint를 완료하고 worker_done을 한 번 보낸다.

## 공백 검사 초기 실패와 재실행

첫 기록 3096e02의 staged 공백 검사는 Docker Compose 정리 로그 10줄의 후행 공백으로 exit 2였다. 체크 실패 뒤 셸이 계속돼 첫 커밋이 생성됐다. 원본 출력·재현 exit 2·명령·정리 로그 원문을 [diffcheck-initial.json](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/diffcheck-initial.json)에 보존했다. 정리 로그의 줄 끝 공백만 정규화했다. 정리 명령의 COMMAND_EXIT:0과 제품/PNG/관측은 바꾸지 않았다. 이후 공백 검사를 기준 ref부터 다시 실행한다.

## 커밋 후 검사와 완료 기록

기록 HEAD `c1b2df18b1a287e63fd7150c50b20c02018fe2e8`의 FullOps `--from 458798c2ee15c179edacfd6f94ebb9896d26f411`은 exit 0이다. ERROR 0·WARNING 1·실행 불가 0이며 등록 product-lint/product-test도 각각 exit 0이다. 원본 [JSON](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/record-lint.json)과 [출력](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-UI-FIX/record-lint.txt)을 보존한다.

유일한 WARNING은 SIZE-001의 기존 누적 PLANS 886줄(기준 877, 상한 500)이다. designer는 자기 결과만 append했다. 기존 통합/실패 기록을 삭제하거나 coor 문서를 분할하지 않았다. 실제 lint 대상은 5파일·추가 147줄이다. SIZE-002·DEP 경고는 없고 의존성 변경도 0이다. 실제 Git 첫 기록은 80파일·추가 474줄이며 64개는 binary PNG다. 보고서·시나리오·두 viewport/상태의 직접 증거와 보존 대조가 지시된 결과여서 별도 제품 과제로 분리하지 않았다.

strict는 exit 0(13개·문제 0·경고 0)이다. 초기 Compose 로그 후행 공백은 원문과 exit 2를 보존한 뒤 정규화했다. 최종 ref-to-HEAD 공백 검사는 exit 0이다. 이 기록 보존과 work.py finish 뒤 새 깨끗한 HEAD에서 같은 lint를 마지막으로 실행하고 preamble worker_done에 전체 SHA·결정·인계 링크를 고정한다.
