---
title: SAR-CONNECT-SIMPLE-001-QA — 에이전트 한도 제거와 연결 단순화 독립 검증(대기 초안)
status: draft
updated: 2026-10-26
owner: tester
tasks: [SAR-CONNECT-SIMPLE-001-QA]
summary: DEV 완료 SHA에서 한도 제거·연결 시도 정리·실패 사유·재시도 폴더·브라우저 화면을 격리 fixture로 독립 검증하는 tester 지시서 초안
---

# SAR-CONNECT-SIMPLE-001-QA — 에이전트 한도 제거와 연결 단순화 독립 검증

이 파일은 대기 초안이다. `handovers/to_tester.md`는 SAR-AUTOSETUP-001-QA(D040 대기)가 쓰고 있다. coor는 tester 인박스가 비고 DEV 완료 SHA가 나온 뒤 `work.py new --role tester --key SAR-CONNECT-SIMPLE-001-QA`로 인박스를 만들고 이 본문을 옮긴다. 이 파일 자체를 dispatch 지시서로 쓰지 않는다.

- 작성일: 2026-10-26
- From / To: designer(제품 기획) → coor → tester
- 상태: blocked — DEV 완료 SHA 대기
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위: fullops/tester의 시나리오·qa-reports 작성, DEV 완료 SHA의 격리 로컬 relay와 합성 회원 fixture(A·B·C), 이 PC의 Windows 폴더 실측(레포 `.fullops-squad/local/` 또는 지정 테스트 폴더 안, D: 드라이브 포함), release 패키지 검사, 커밋·push.
- 2단계 승인 범위(사용자 결정 2026-10-26, 아래 '실연결 재검증' 절): 배포된 relay에서 Grok Bot 쪽 `knowslog01@gmail.com` 회원과 Claude Code 쪽 `scshin88@gmail.com` 회원의 실제 연결, 두 회원 사이 친구 요청·수락, ego 브라우저에서 에이전트가 Google 로그인·지문 비교·이름·연결 동의를 직접 진행.
- 추가 승인이 필요한 행위: 1단계에서 운영 relay·운영 agent 사용, 실제 메시지 발송(2단계 포함), 비밀번호·2단계 인증 코드 입력, 사용자 실제 `~/.claude`·`~/.codex` 설정 변경, 위 두 회원 외 계정 사용, 기존 운영 agent 철회·정리, 제품 코드 수정(결함은 보고만), 제품 규칙 변경.
- 상위 의도: 없음. [D02 SAR-CONNECT-SIMPLE-001](../../docs/planning/product-specs/SAR-CONNECT-SIMPLE-001.md) CS-AC-01–CS-AC-16의 독립 동작 확인.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/tester / fullops/tester
- 병합 책임자 / 기본 브랜치: coor / main(origin)
- 복귀: coor 워크트리 C:/Users/shin/orca/workspaces/KnowsLink/coor (run id·task id·dispatch id는 dispatch 후 기록)

## 현재 상황과 확인 근거

사용자가 에이전트 연결 한도 제거와 연동 단순화를 요청했다(2026-10-26). designer가 D02 SAR-CONNECT-SIMPLE-001과 D04 절을 정했다. DEV가 [to_dev.md](../to_dev.md)로 구현한다. 새 흐름(진행 중인 연결, 연결 링크 상태 화면, Google 실패 구분)과 연결 상태 변경이 있어 lite에서도 tester 독립 QA를 배정한다.

## 적용 기준과 예외

- 공통 기준: `fullops-common-0.3.6` — `.fullops-squad/rules/common/README.md`, `testing.md`, `security.md`.
- 프로젝트 정본: `.fullops-squad/project.md`, 웹 실행은 `fullops-test` web-engine과 `docs/evaluations/scenarios/README.md`.
- 기준: DEV 완료 SHA(배정 때 기록). ID 정의는 D02 SAR-CONNECT-SIMPLE-001.
- 예외: 없음.

## 먼저 읽을 문서

1. [D02 SAR-CONNECT-SIMPLE-001](../../docs/planning/product-specs/SAR-CONNECT-SIMPLE-001.md)
2. [D04 연결 단순화 절](../../docs/design-docs/mockups/service-member-ui.md#연결-단순화와-내-에이전트-정리--sar-connect-simple-001)
3. DEV 완료 보고와 `docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-DEV-test/report.md`
4. 재사용할 기존 시나리오: `docs/evaluations/scenarios/SAR-UX-ISSUES-002-CONNECT-TESTER.md`(연결 재개·만료·취소 회귀)

## 해야 할 일과 파일 소유권

- [ ] 새 시나리오 `docs/evaluations/scenarios/SAR-CONNECT-SIMPLE-001-TESTER.md`를 쓴다. covers와 독립 checks, 격리 fixture origin·포트·DB 식별자, 허용 GET·명시 POST를 선언한다.
- [ ] 한도: 회원 A가 활성 agent 7개를 만든다(CS-AC-01). 철회 10개 보존 회원이 새 agent를 만든다(CS-AC-02). 화면에 '회원당', '5개', '먼저 철회' 문구가 없다.
- [ ] 연결 시도: 취소·만료·동의 뒤 미저장 시도 5회 뒤 '내 에이전트' 수와 친구 B 화면의 A agent 수가 같고, 다음 정상 연결이 성공한다(CS-AC-05). '이 시도 정리'가 즉시 지운다. 정리 전 저장 완료가 정상 완료된다(CS-AC-06, 기한 정리는 DEV 시계 주입 근거 재사용).
- [ ] 이름: 내 컴퓨터·웹 추가 경로의 기본 이름과 두 번째 번호(CS-AC-07). Dots 원격은 fixture로 가능한 범위만, 불가하면 DEV 테스트 근거와 not_run 사유. 이름 없는 기존 agent의 바로 보이는 이름 저장(CS-AC-08).
- [ ] 재시도 폴더: 남은 시도 파일 폴더, 미리 만든 빈 폴더, D: 드라이브 폴더에서 다시 연결(CS-AC-10). 옮긴 폴더가 지워지지 않았는지 확인한다.
- [ ] 실패 사유: E 표의 `expired`·`cancelled`·`unconfirmed`·`private_folder`·`unreachable`을 주입해 클라이언트 출력의 사유 코드·문구·다음 행동을 KO·EN으로 확인한다(CS-AC-11). `rate_limited`·`agent_rate`·`service_capacity`는 DEV 테스트 근거 재사용 가능.
- [ ] 브라우저: 완료 링크 재방문을 시작한 회원 A·다른 회원 B·로그아웃으로 확인한다(CS-AC-12). Google 실패 세 원인은 fixture로 가능한 범위만, 불가하면 DEV 테스트 근거와 not_run 사유(CS-AC-13).
- [ ] Grok release 기본값: 패키지 검사로 환경 변수 없는 설정에서 동의 전 서버 요청 0건, 동의 뒤 연결(CS-AC-14). 기존 환경 변수 설정이 그대로 동작한다.
- [ ] 보안 회귀: 출력·로그·화면에서 키·토큰·credential 검색 0건, 다른 회원 상태·이름·이메일 비노출(CS-AC-16).
- [ ] 기존 시나리오 `SAR-UX-ISSUES-002-CONNECT-TESTER.md`의 재개·만료·취소 핵심을 다시 돌린다.

파일 소유권: tester는 시나리오와 qa-reports만 쓴다. 제품 코드·D02·D04는 고치지 않는다.

## 완료 기준과 검증

- 위 항목마다 PASS/FAIL/not_run과 근거가 `docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-QA-test/report.md`에 있다. 결함은 심각도(critical/high/medium/low)와 재현 절차를 적는다.
- DEV의 `make lint`·`make test`·`package_plugin.py --verify` 결과는 같은 SHA·환경이면 재사용한다. 새 실패가 있으면 다시 실행한다.
- FullOps lint: `python3 C:/Users/shin/.fullops-squad/runtime/fullops-squad/packages/02fb5d63e340de209eeb446307db1c2b7f2264fa94de1a6835be87f00c3d5f29/scripts/lint.py --repo C:/Users/shin/orca/workspaces/KnowsLink/tester --from <배정 기준 SHA>` ERROR 0.
- 증거 형식: `.md`는 front matter 필수, 원시 로그는 `.txt`/`.log`. 커밋하는 캡처는 report.md가 링크하는 것만이다. 로그와 임시 파일은 `.fullops-squad/local/SAR-CONNECT-SIMPLE-001-QA/`.
- 시간 예산: fixture 준비 포함 2시간. 넘으면 진행 상황을 escalation으로 알린다.

### 실연결 재검증 — 2단계(사용자 결정 2026-10-26)

1단계(격리 fixture) PASS와 운영 배포(ops, 사용자 승인) 뒤 실행한다. 1단계와 같은 지시서의 별도 절로 기록하고 PASS를 섞지 않는다.

1. 회원 구성: Grok Bot 쪽은 `knowslog01@gmail.com` 회원, Claude Code 쪽은 `scshin88@gmail.com` 회원이다. 서로 다른 두 회원이므로 친구 요청·수락 흐름을 포함한다.
2. 브라우저 단계: Google 로그인, 계정 선택, 지문 비교, 에이전트 이름, 연결 동의, 친구 코드 요청·수락을 사용자에게 넘기지 않는다. 에이전트가 ego 브라우저에서 모두 진행한다. 지문은 클라이언트 출력과 브라우저 화면을 에이전트가 직접 대조하고 결과를 기록한다.
3. 멈추는 경우: Google이 비밀번호·2단계 인증 코드·새 기기 확인을 요구하면 입력하지 않고 coor에 `ask`로 알린다. 다른 Google 계정이 선택돼 있으면 W-02 화면을 확인한 뒤 지정 계정으로 바꾼다.
4. 측정: 호스트별 T-01–T-04 사람 행동 수 대신 에이전트가 대신한 행동과 사람이 꼭 한 행동(있다면)을 순서대로 기록한다. 연결 실패가 나면 E 표 사유 코드·문구·다음 행동이 보였는지 기록한다. 실패 시도가 '내 에이전트' 목록에 늘지 않는지 확인한다(CS-AC-05).
5. 확인 범위: 두 회원 각각 연결 완료(키 등록됨), 기본 이름, 링크 재방문 완료 화면(W-01), 친구 요청·수락 뒤 서로의 agent가 이름과 함께 보이는지, 도구 결과의 표시 이름(A-06).
6. 메시지 발송: 두 회원 사이 시험 메시지 1건은 coor를 통해 사용자 승인을 받은 뒤에만 보낸다.
7. 운영의 기존 agent(`agent_c18976751d5a96d8adfbb0` 등)와 이름 없는 agent는 건드리지 않는다. 2단계에서 새로 만든 agent도 사용자 지시 없이 철회하지 않는다.

### designer 시각 검수(ART) 인계 조건

UI 변경(M·W)이 있으므로 designer 직접 시각 검수를 배정한다. 조건과 범위:

- 시점: DEV 완료 SHA가 나오면 tester와 병행할 수 있다. 같은 fixture를 쓰면 상태를 바꾸는 POST 순서를 coor가 조정한다.
- 대상: '내 에이전트'의 진행 중인 연결·이름을 정할 에이전트 절, 연결 링크 화면(대기·동의 직후·완료·만료·취소), 연결 중 Google 실패 화면. KO·EN × 1440·390·320px.
- 판정 기준: D04 절의 문구 표와 '자기 비평과 수락'(CS-AC-17). 보고는 `docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-ART-test/report.md`.
- 캡처는 판정 항목에 연결된 것만 커밋한다.

## 갱신할 산출물

없음(QA 보고서와 시나리오만). D02·D04·D11 변경이 필요해 보이면 완료 보고 `남은 일`에 적는다.

## 기대 산출물

- `docs/evaluations/scenarios/SAR-CONNECT-SIMPLE-001-TESTER.md`
- `docs/evaluations/qa-reports/SAR-CONNECT-SIMPLE-001-QA-test/report.md`

## 제약·협업·후속

- 1단계는 운영 relay·운영 agent·실제 상대 발송 금지. 2단계는 위 '실연결 재검증' 절의 범위만 쓴다. 운영의 이름 없는 agent와 `agent_c18976751d5a96d8adfbb0`은 건드리지 않는다.
- D: 드라이브 실측 폴더는 검증 뒤 사용자 데이터와 섞이지 않게 테스트 폴더 안에만 만든다. 삭제·이동 정리는 `AGENTS.md`의 'Windows 파일 작업 안전'을 따른다.
- SAR-AUTOSETUP-001-QA(D040)와 별개다. 그 과제의 결과를 이 과제 PASS로 쓰지 않는다.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
항목당 1~2줄로 쓰고 상세는 `qa-reports/SAR-CONNECT-SIMPLE-001-QA-test/report.md` 링크로 대신한다. worker_done 본문 첫 줄은 `[완료] SAR-CONNECT-SIMPLE-001-QA | 브랜치 fullops/tester | SHA <브랜치 HEAD>` 한 줄로 쓴다.
