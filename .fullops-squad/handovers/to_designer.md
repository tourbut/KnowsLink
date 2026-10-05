---
title: SAR-PUBLIC-IDENTITY-001-UI — 고정 이메일 신원 후보의 UX01–03을 실제 브라우저에서 직접 검수한다
status: draft
updated: 2026-10-05
owner: designer
tasks: [SAR-PUBLIC-IDENTITY-001-UI]
summary: 고정 이메일 신원 후보의 UX01–03을 실제 브라우저에서 직접 검수한다
---

# SAR-PUBLIC-IDENTITY-001-UI — 일반 이메일 가입·세션 직접 시각 검수

- 상태 ready. 고정 제품 후보59b66ada8b36802484cc6d7e22523257b50572cc, 코드a446d89ff288c4243ad6d7f8780a778517154584, 기준94533b207b456c0560800fe30a7c90b2b5887c6e.
- 소유: 기록 checkout /home/shin/orca/workspaces/KnowsLink/fullops-coor의 이 designer 인박스·해당 logs 전문, docs/evaluations/qa-reports/SAR-PUBLIC-IDENTITY-001-UI.md와 필요한 UI 증거, 자기 실행 기록뿐. 제품 코드·다른 기록·PLANS/board·다른 인박스·GitHub·CF·배포는 수정하지 않는다.
- 복귀 term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9, Run run_8ca8bc058ab7. 새 세션 dispatch preamble을 사용한다.

## 적용 기준과 예외

fullops-common-0.3.2, FULLOPS.md/project.md/document-writing.md/orca-agents.md, coding-style/testing/security. 제품 SAR-PUBLIC-SERVICE PS01–04·해당PS11·UX01–03과 frozen SAR-MVP C1–C5. 사용자 일반 서비스 구현 지시에 따른 직접 UI 검수다. 새로운 제품수치/화면/코드 구현이나 기술 승인을 맡지 않는다. 최종 시험은 같은 일반 이메일의 노우↔다닷이다. 이번 로컬 fixture 직접 UI를 실제 사람 이메일/운영 시험 PASS로 대체하지 않는다.

## 먼저 읽을 문서

필수 규칙과 제품/UX 정본, DEV 실행 기록 SAR-PUBLIC-IDENTITY-001-DEV.md, README의 일반 이메일 로그인 확인 절차, 사용자 가이드, internal/relay/member.go·identity.go·http.go·scripts/mail_sink.py. 탐색 근거는 coor가 아래 추가한다.

## 지시 전제와 충돌

기존 owner-only Basic 화면은 일반 회원 가입 화면이 아니다. 이번 UX03 범위는 신원·세션·자기 owner·로그아웃이다. agent/key/pair/회원 비활성화는 다음 기능이다. 준비 중 표시는 실패로 바꾸지 말고 기능이 아직 미구현이라는 표시와 실제 성공버튼 없음 여부를 검사한다. 실제 이메일·운영 SMTP 미설정은 실제 확인만 미실행이다. fixture의 로컬 SMTP 받은 코드로 직접 브라우저 시각/흐름 검사 자체를 수행한다. 자동 QA/보안 리뷰는 별도 담당이므로 전체를 중복하지 않는다.

## 해야 할 일·완료 기준

1. 자신의 제품 실행 checkout /tmp/knowslink-public-identity-ui-59b66ad의 detached59를 확인하고 별도 임시 PostgreSQL·SMTP sink·loopback relay를 시작한다. 기존 운영 서비스/QA의 임시 환경을 변경하지 않는다. 코드를 수정하지 않는다. 필요한 실행/fixture 준비를 독립 수행한다.
2. Orca 실제 내장 브라우저에 새 전용 페이지를 열어 그 page ID를 명시적으로 제어한다. 사용자 기존 Cloudflare/다른 페이지를 탐색하거나 변경하지 않는다. orca-cli browser reference를 읽는다.
3. UX01 시작 화면의 일반 이메일 입력·로그인 절차, UX02 확인대기·마스킹·오답·만료·제한·발송 실패의 안내와 다음 동작, UX03 실제 회원 홈·빈 agent/관계 준비중·현재 로그아웃·전체 로그아웃/재확인·오래된 세션 화면을 직접 확인한다. 기존 회원 gate 영향 화면은 변경 영향에 필요한 범위만 확인한다.
4. 성공/실패의 문구·입력 레이블·가독성·레이아웃·키보드·모바일 폭을 판정한다. 테스트는 기능 판정에 필요한 최소 조작으로 한다. 실제 직접 화면 관측과 자동 근거를 구분한다. 스타일 변경을 새로 구현하지 않는다. 새 장식·별도 frontend·채팅 UI를 요구하지 않는다.
5. 필요한 화면만 캡처하고 이메일·코드·cookie·token은 캡처/보고 전에 숨긴다. dummy fixture도 값 전문을 보고하지 않는다. 보고에 환경·실제 후보SHA·브라우저page ID·경로/근거·PASS/FAIL/미실행·심각도/재현을 쓴다. 단계별 필수 실패/미해결critical/high가 있으면 수락 불가를 보고하고 임의 수정하지 않는다.
6. 자신의 임시 실행 자원을 정리하고 다른 사용자 자료를 보존한다. 이 인박스에 전문 완료 보고를 쓰고 work.py finish --role designer --key SAR-PUBLIC-IDENTITY-001-UI로 보존한다. 자신의 소유 파일만 commit하고 다른 root/reviewer 기록을 함께 add/commit하지 않는다. fixed59·결과SHA·독립관측/실제확인미실행·finding을 이 Run worker_done으로 회신한다.

## 완료 보고

실행 시각 검수자가 작성한다.

## coor 탐색 근거

SAR-PUBLIC-IDENTITY-001-UI-{find,documents-find,context}.json은 코드/문서를 분리한20개 후보와 필수 규칙의 근거다. keep은 README·member.go·mail_sink.py·identity.go·http.go·실행준비 compose/main/verify_mvp, DEV 실행 기록·이번 UX/제품 정본·필수규칙이다. 이전UI 보고·로그와 담당 context는 과거 기준/실패의 참고다. 현재후보59를 직접 확인한다. context의 추천을 참고하되 필수정본·passage 미송신자료를 제외하지 않는다. 기존trial 조작 run_trial.py는 필요시 확인만 하며 종료된trial을실행하지 않는다.
