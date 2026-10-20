---
title: SAR-HOST-VERIFY-001-LIVE — 운영에서 Claude Code 세션·Grok Bot 신규 설치와 실계정 왕복·알림을 검증한다
status: draft
updated: 2026-10-20
owner: tester
tasks: [SAR-HOST-VERIFY-001-LIVE]
summary: 운영에서 Claude Code 세션·Grok Bot 신규 설치와 실계정 왕복·알림을 검증한다
attempt: d42304b32b524a52b58d877c6bae2eae
base: f591dc0079e8fd5e488e2260da506e1efb89627b
subagent_level: standard
test_level: lite
---

# SAR-HOST-VERIFY-001-LIVE — 운영에서 Claude Code 세션·Grok Bot 신규 설치와 실계정 왕복·알림을 검증한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → tester
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위 / 추가 승인이 필요한 행위: 사용자 결정 D008·D013(2026-10-20)에 따라 운영(link.knowslog.com, 제품 7b4f9bc)에서 Claude Code 작업 세션과 Grok Bot의 깨끗한 신규 설치·연결·친구 수락·실제 메시지 왕복·수신 알림을 검증한다. 이 PC에 새 설치 폴더를 만들고 공개 설치 안내대로 설치하는 것은 승인됨. 각 단계의 실계정 로그인, 친구 요청·수락, 메시지 발송, Grok 앱 설정 변경, 기존 설치·키의 철회·삭제는 실행 직전에 `ask`로 coor를 통해 사용자 승인을 받는다(로그인은 사용자가 직접 한다). 운영 서버·DB·ingress 변경, 기존 계정의 에이전트·키 삭제, force push는 하지 않는다.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/tester / fullops/tester
- 병합 책임자 / 기본 브랜치: coor / main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

- 설계 U8 호스트 표(`.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` '근거에 따른 호스트 표시')는 Claude Code 작업 세션을 '시험 환경 확인(Claude Code 2.1.302, Windows)', Grok Bot을 '기존 계정 사용 확인'으로 둔다. 깨끗한 신규 설치와 운영 실계정 왕복은 미검증이다. 이 과제가 그 검증이다. GitHub #8 댓글에도 후속으로 적었다.
- 이전 실계정 검증: `.fullops-squad/docs/evaluations/qa-reports/SAR-HUMAN-PEER-002-LIVE/report.md`(Grok knowslog01@gmail.com·Dots scshin88@gmail.com, Orca 브라우저 프로필 분리, Grok Bot native computer-use, evidence.json), 알림 관측 `SAR-HUMAN-PEER-003-NOTIFY/report.md`, 세션 QA `SAR-SESSION-PEER-001-TESTER-test/`·`SAR-SESSION-PEER-002-TESTER-test/`.
- 설치 안내: README.md(Grok 신규 설치), adapters/README.md('Claude Code 작업 세션 연결'), 운영 화면의 '어디에서 사용할까요?' 안내.

## 적용 기준과 예외

fullops-common-0.3.3(보안: 개인 키·쿠키·친구 코드·토큰을 증거에 넣지 않음), `.fullops-squad/project.md`, 설계 U8 '후속 검증표' 항목(호스트별 제품/버전·OS·계정 권한·소스 SHA·패키지 검사·깨끗한 앱 설치·도구 검색·인증·이름 등록·친구 수락·왕복·재개·검증 날짜·근거를 각각 기록, 소스 PASS로 앱 PASS를 채우지 않음). 서버 API를 직접 호출해 에이전트를 흉내 내지 않는다.

## 먼저 읽을 문서

위 설계 U8 절, 이전 LIVE·NOTIFY 보고와 evidence 형식, README.md·adapters/README.md 설치 절, `.fullops-squad/contexts/tester.md`, 운영 배포 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-RELEASE-002-OPS/report.md`.

## 해야 할 일과 파일 소유권

- [ ] 시나리오 `.fullops-squad/docs/evaluations/scenarios/SAR-HOST-VERIFY-001-LIVE.md`: 계정 배정(사용자에게 확인: 예 Claude Code 쪽 Google 계정과 Grok 쪽 계정), 단계, 사용자 승인 지점, 증거 형식, 중단 조건.
- [ ] Claude Code 작업 세션: 새 폴더에 공개 안내대로 깨끗한 설치 → 운영 연결(웹 동의는 사용자) → 에이전트 이름·키 등록 확인 → `/knowslink-session on` → 친구 상대(Grok 계정)와 메시지 왕복(발송은 승인 대화상자에서 사용자 승인) → 도착 toast·본문 조회 확인. 기존 사용자 Claude 설정·다른 프로젝트를 바꾸지 않는다.
- [ ] Grok Bot: 공개 안내의 신규 설치 경로로 새 설치(또는 기존 설치와 분리된 새 연결) → 운영 연결 → 친구 수락 → 왕복 → 루틴/webhook 수신 알림 확인. Grok 앱 조작이 필요하면 computer-use와 사용자 승인으로 진행한다.
- [ ] 단계마다 호스트·버전·OS·소스 SHA·명령 자신의 종료코드·화면 관측을 U8 후속 검증표 형식으로 기록한다. 실패·막힘은 원 오류와 함께 보존하고 우회하지 않는다.
- [ ] 결과 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-HOST-VERIFY-001-LIVE/report.md`와 evidence.json(비밀 제외), 호스트별 '신규 설치 검증 완료' 판정(조건을 모두 통과한 경우만).
- [ ] 완료 보고 → work.py finish → 커밋 → worker_done. 설계 호스트 표 갱신이 필요하면 coor에 알린다(designer 담당).

소유 파일: 위 시나리오·QA 폴더, own route/inbox/log. 설치 폴더는 workspace 밖일 수 있으나 경로를 보고에 남기고 다른 사용자 자료를 바꾸지 않는다.

## 완료 기준과 검증

두 호스트 각각 신규 설치→연결→친구→왕복→알림의 단계별 PASS/FAIL/not_run과 근거가 있다. 실패는 재현·원 오류·담당(DEV/ops/사용자) 제안과 함께 남긴다. 운영 데이터는 사용자 승인 범위만 바뀐다.

### UI 작업의 디자인 기준과 검증

해당 없음(실계정 동작 검증).

## 갱신할 산출물

없음(설계 표 갱신은 designer).

## 기대 산출물

시나리오, 결과 보고·evidence, 완료 보고.

## 제약·협업·후속

사용자 결정이 필요한 지점마다 `ask`로 멈추고, 답을 받기 전 실행하지 않는다. 사용자 로그인 화면에 자격을 대신 입력하지 않는다. 제품 코드·다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
