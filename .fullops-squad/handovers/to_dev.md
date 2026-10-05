---
title: SAR-PUBLIC-IDENTITY-001-REVIEW — 이메일 신원·세션 구현을 별도 고성능 세션에서 독립 검토한다
status: draft
updated: 2026-10-05
owner: dev
tasks: [SAR-PUBLIC-IDENTITY-001-REVIEW]
summary: 이메일 신원·세션 구현을 별도 고성능 세션에서 독립 검토한다
---

# SAR-PUBLIC-IDENTITY-001-REVIEW — 일반 이메일 신원·세션 독립 보안 리뷰

- 상태: ready. 구현자와 다른 새 Claude Opus5.5high 세션.
- 복귀: term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9, Run run_8ca8bc058ab7. 새 dispatch preamble 사용.
- 기록 체크아웃: /home/shin/orca/workspaces/KnowsLink/fullops-coor. 소유 파일은 이 인박스·해당 logs 전문과 SAR-PUBLIC-IDENTITY-001-REVIEW-review 디렉터리 및 자기 리뷰 실행 기록뿐이다. 제품·PLANS·board·다른 인박스·GitHub·배포/CF를 수정하지 않는다.

## 적용 기준과 예외

fullops-common-0.3.2와 FULLOPS.md/project.md/document-writing.md/orca-agents.md, coding-style/testing/security를 적용한다. 제품 PS01–04·신원 PS11·UX01–03, frozen C1–C5를 유지한다. 독립 위험 변경 리뷰는 고성능 모델에 배정하는 fullops-review 규약을 적용한다. Jev의 최종 추천도 Claude Opus5.5high다.

리뷰 기준 base=94533b207b456c0560800fe30a7c90b2b5887c6e, head=59b66ada8b36802484cc6d7e22523257b50572cc. 깨끗한 detached 읽기 전용 snapshot은 /tmp/knowslink-public-identity-review-59b66ad다. snapshot의 tracked 파일을 바꾸거나 그 안에서 build/test/install하지 않는다. 필요한 재현은 별도 scratch 사본에서 실행한다. 결과는 기록 checkout에만 작성한다. 구현자의 실제 provider 세션 ID는 coor가 아래 추가한다. 자신의 실제 별도 세션 ID를 확인해 independence에 기록한다.

## 먼저 읽을 문서

.fullops-squad/FULLOPS.md, project.md, rules/common/README.md와 세 규칙, document-writing.md, 제품 SAR-PUBLIC-SERVICE.md·SAR-MVP.md, UX SAR-PUBLIC-SERVICE-UX.md, DEV 실행 기록 SAR-PUBLIC-IDENTITY-001-DEV.md, review/rule.json. open-code-review-delegate 및 fullops-review 스킬을 적용한다. 별도 탐색 근거는 coor 추가 절을 따른다. preview.json의 모든 변경 파일은 리뷰 범위다.

## 해야 할 일·완료 기준

1. prepare된 preview/rules/result/report를 읽고 945..59의 변경과 관련 경계를 검토한다. 이메일 코드 발급·검증·SMTP 실패·추측/타이밍·Challenge 식별·신원 연속성·동시 첫 가입·rate/restart·세션 fixation/CSRF/현재·전체 logout·cross-owner/gate·공개 synthetic 우회·민감정보 노출·기존 frozen 업무 회귀를 확인한다.
2. 구현자 검사 참고와 자신의 관측을 구분한다. 실제 운영 SMTP/사용자 이메일이 없는 사람 확인을 PASS로 만들지 않는다. 이번 단계와 다음 AGENTS/MESSAGES의 명시된 범위를 구분한다. 미해결 제품 기준 위반은 영향·심각도와 재현을 남긴다.
3. 모든 (path,status) 항목을 reviewed/skipped+구체적 reason으로 판정한다. 제외 증거는 요약·manifest·무결성 근거로 확인한다. finding은 정확한 파일/줄·심각도·resolved와 재현을 쓴다.
4. 원본 DEV checkout /home/shin/orca/workspaces/KnowsLink/fullops-dev가 clean/head59b66ad임을 확인한 뒤 lint --from94533b2를 실행해 리뷰 디렉터리 lint.json에 보존한다. snapshot은 계속 읽기 전용이다. report에 규칙/refs/독립 세션/검사별 대상 SHA·종료코드·미실행·warning 영향을 쓴다.
5. review.py check --repo 기록checkout --key SAR-PUBLIC-IDENTITY-001-REVIEW --from945fullSHA --to59fullSHA --task-key SAR-PUBLIC-IDENTITY-001-DEV를 통과시킨다. 이 검사는 기록 일치이며 자동 의미 수락이 아님을 구분한다. critical/high 또는 필수 실패가 있으면 수락 불가로 보고하고 코드를 직접 고치지 않는다.
6. 이 인박스에 완료 보고 전문을 쓰고 work.py finish --role dev --key SAR-PUBLIC-IDENTITY-001-REVIEW로 보존한다. 자신의 소유 파일만 git add/commit 한다. root/coor의 다른 자료를 함께 커밋하지 않는다. 제품 후보59와 리뷰 결과SHA·finding·미실행·check를 현재 Run worker_done으로 회신한다.

## 완료 보고

실행 검토자가 작성한다.

## coor 탐색·독립 세션 근거

구현자의 실제 Claude provider 세션은 e0666abc-1cf2-49ac-a288-45a8043404bc다. 해당 session JSONL의 task_568f0a5f227c와 실제 cwd·2026-10-05T12:48:57.301Z를 대조했다. 검토자는 자기의 다른 provider 세션 ID를 기록한다.

탐색 근거는 SAR-PUBLIC-IDENTITY-001-REVIEW-{find,documents-find,context}.json이다. 코드/문서를 분리했고 fallback 없음, conflict_ids/caution_ids 없음이다. keep은 identity.go·README·identity_test.go·member.go·mail.go·http.go·mail_sink.py와 필수 규칙, DEV 실행 기록이다. 추가 keep: identity_integration_test.go·store.go·cmd/relay/main.go(실제 변경 파일). 다른 keep인 deploy verify.py·test_messages_integration_test.go·verify_setup.py·registry.json·기존 제품 문서 리뷰·리뷰 템플릿·QA 인박스는 관련 검사/이력 근거로 읽는다. .gitignore는 omit? 추천이며 필요한 경우 확인한다. preview의 전체 변경 검토를 줄이지 않는다.

지시 전제와 충돌을 구분한다. 기존 제품 문서의 수락은 코드/실제 운영 수락이 아니다. deploy verify.py의 public302는 이전 owner-only 경계 검사이다. 이번 실제 이메일 positive와 바꾸어 쓰지 않는다. 기존 gate/timeout 검사는 해당 기능 회귀 근거이며 이번 신원 QA를 대신하지 않는다. 필요한 검사 재현은 별도 scratch에 수행하고 snapshot을 수정하지 않는다.
