---
title: SAR-CONNECT-SIMPLE-001-TESTER — 에이전트 한도 제거와 연결 단순화 독립 QA 시나리오
status: draft
updated: 2026-10-27
owner: tester
tasks: [SAR-CONNECT-SIMPLE-001-QA]
summary: "격리 relay fixture에서 실제 배포 클라이언트(plugin.js)와 Handler로 CS-AC-01–16을 확인하는 시나리오, 허용 요청, fixture 제어, 통과 조건"
---

# SAR-CONNECT-SIMPLE-001-TESTER — 독립 QA 시나리오

대상은 main `3f1da403086e4449a7b3df214c441cd1b5972db6`의 제품 코드다(검증 브랜치 `fullops/tester` `5c750943`, 제품 경로 diff 없음). 규칙은 [D02](../../planning/product-specs/SAR-CONNECT-SIMPLE-001.md), 화면은 [D04 절](../../design-docs/mockups/service-member-ui.md#연결-단순화와-내-에이전트-정리--sar-connect-simple-001)이다. 결과는 [QA 보고서](../qa-reports/SAR-CONNECT-SIMPLE-001-QA-test/report.md)에 있다. 테스트 레벨은 lite다. Jev 웹 조작 도구는 쓰지 않는다. 이유: 비가역 생성의 순서 통제, 실제 클라이언트 프로세스, Windows 폴더 권한 측정, 상태 시계 이동이 필요하다.

## 공통 조건

- origin은 `http://127.0.0.1:18951`(주 fixture), 이전 시나리오 회귀는 `fixture2-manifest.json`의 임의 loopback 포트 하나다. 운영 relay·운영 agent·실계정·실메일·실제 메시지 발송·사용자 실제 `~/.claude`·`~/.codex`는 쓰지 않는다.
- DB는 격리 `postgres:17` 컨테이너 `kl-csqa-pg`(loopback 15531, 라벨 `knowslink.tester.qa=csimple001`, 512MiB)의 `knowslink`(주)와 `knowslink2`(회귀)다. 검사 뒤 컨테이너를 정지했다.
- fixture는 실제 `Service.Handler()`를 그대로 쓴다. 추가한 것은 fixture 제어 경로와 Google 대역뿐이다(`connect_simple_qa_fixture_test.go`, 태그 `integration connectsimpleqa`).
  - 허용 제어: `POST /fixture/as/{a..e}`(Google 발급자 회원 세션), `GET /fixture/state?m=`(읽기), `POST /fixture/exp?connection=&delta=`(시도 하나의 만료 시각 이동), `POST /fixture/fill-members`(회원 상한 채움), `/fixture/drop-complete?on=`(키 저장 응답 한 건을 서버에 보내지 않음), `POST /fixture/nameless?m=`(이전 relay가 남긴 이름 없는 agent 한 개 씨앗), `/fixture/hits`(API 요청 수).
  - Google 대역: `/fixture/google/auth`(계정 선택 화면, PKCE 확인)와 `/fixture/google/token`(서명한 ID 토큰). 실제 Google 화면·로그인은 아니다.
- 합성 회원 A–E(`csqa-x@example.test`)만 쓴다. 계정 C·E는 친구다(제품의 요청·수락 폼으로 만들었다).
- 클라이언트는 이 SHA의 `adapters` 빌드 `dist/plugin.js setup-connect`(설치기가 쓰는 실제 진입점)와 `dist/login.js`(시작·폴링 한 번씩), MCP stdio 클라이언트다. 출력은 JSON 줄 그대로 저장한다.
- 폴더 실측은 D:(`D:/workspace/KnowsLink/.devtools/connect-simple-001-qa/work/…`)와 C:(`.fullops-squad/local/SAR-CONNECT-SIMPLE-001-QA/work-c…`)의 새 폴더에서만 한다. 삭제·이동 정리는 하지 않았다.
- relay는 익명 연결 API에 IP당 분당 30회, `/home/*`에 회원당 분당 40회 한도가 있다. 드라이버는 429를 기다린다.

## covers와 독립 checks

| covers | 독립 checks (`harness/`) | 통과 조건 |
|---|---|---|
| CS-AC-01 | `s_limits.py` | 회원 A가 웹에서 활성 agent 8개를 만든다. 모두 303이다. KO·EN 화면 7쪽에 '회원당'·'5개'·'먼저 철회'·'남은'·'한도' 문구가 없다. |
| CS-AC-02 | `s_limits.py` | 회원 B가 agent 10개를 만들고 철회한다. 철회 기록 10개가 남고(`revoked`) 새 agent가 만들어진다. |
| CS-AC-03 | DEV 근거 재사용 | `TestAgentAndPairCapacity`. 임계 주입이 필요해 fixture로 만들지 않는다(not_run, 아래 사유). |
| CS-AC-04 | `s_reasons.py`(보조), DEV 근거 | B가 24시간 안 20개를 만든 뒤 21번째 웹 생성과 연결 동의가 `agent_rate`와 시각으로 거절된다. 기존 agent 수가 같다. 시계 주입은 DEV 근거. |
| CS-AC-05 | `s_attempts.py` | C가 취소·relay 만료·앱 종료(대기)·로그인만 한 시도·동의 뒤 미저장 시도 5종을 만든다. '내 에이전트' 카드와 서버 agent 수가 0이고 친구 E 화면의 C agent가 0이다. 다음 정상 연결이 성공하고 양쪽 화면에 1이 보인다. |
| CS-AC-06 | `s_cleanup.py` | '이 시도 정리' 뒤 서버에서 즉시 사라진다. 동의 뒤 저장이 유실된 시도는 기한 후에도 정상 완료된다. 기한 후 23시간은 남고 25시간은 정리된다. 정리된 시도는 완료할 수 없다. 키 있는 agent는 그대로다. |
| CS-AC-07 | `s_attempts.py`, `s_names.py`, `s_nohost.py`, 브라우저 | `--host claude` 동의 화면 기본 이름 `Claude Code`, 두 번째 `Claude Code 2`, `codex`는 `Codex`, 호스트 없음은 `내 에이전트 N`, 웹 추가 칸은 `내 에이전트 N`이고 만든 뒤 번호가 오른다. Dots 원격은 fixture로 불가(not_run). |
| CS-AC-08 | `s_names.py`, `s_names2.py`, `ego-round-01.mjs` | 이름 없는 agent 카드가 목록 앞 절(`#unnamed`)에 이름 입력과 저장 버튼을 접지 않고 보인다. 저장하면 이름이 정해지고 절이 사라진다. |
| CS-AC-09 | `s_grok.py`(Grok 한 호스트), DEV 근거 | Grok `knowslink_status`가 내 agent 표시 이름을 보인다. Claude Code·Codex와 보낸·받은 출력은 DEV 스냅숏. |
| CS-AC-10 | `s_folders.py` | C:·D: 각각 취소 시도 파일 폴더·대기 시도 폴더·빈 폴더·새 폴더·연결 완료 폴더에서 사용자 정리 없이 연결 성공한다. 남은 파일은 `.unused-<시각>`으로 옮겨지고 지워지지 않는다. 권한은 소유자 현재 사용자, 상속 차단, 허용은 현재 사용자뿐이다. 상위가 파일인 경로는 `private_folder`다. |
| CS-AC-11 | `s_expiry_real.py`, `s_attempts.py`, `s_folders.py`, `s_reasons.py` | 사유 코드·문구·다음 행동: `expired`(실제 10분 대기)·`cancelled`·`unconfirmed` 경로·`private_folder`·`unreachable`·`rate_limited`·`agent_rate` KO·EN. 사유 없는 `failed`가 없다. `service_capacity`·`signed_out`은 DEV 근거·웹 전용. |
| CS-AC-12 | `s_browser.py main`, `ego-round-04.mjs` | 시작한 회원 A는 완료 화면과 이름, 동의 뒤 저장 전 화면, 취소·만료 화면을 본다. 다른 회원 B와 로그아웃은 현행 일반 문구이고 이름·상태가 없다. |
| CS-AC-13 | `s_browser.py main·signup` | 다른 계정, 새 가입 닫힘, 로그인 기한 지남이 서로 다른 문구다. 응답에 다른 계정의 이메일·이름·ID가 없다. |
| CS-AC-14 | `s_grok.py` | 이 SHA의 Grok 패키지 zip과 설치기 `defaults.json` 형태로, 환경 변수 없이 동의 전 요청 0건·확인 없는 connect 거절·동의 뒤 연결·자동 폴더. 환경 변수 설치는 잘못된 defaults보다 env가 이긴다. defaults 없는 패키지는 held다. |
| CS-AC-15 | `package_plugin.py --host grok --verify` | 종료코드 0. |
| CS-AC-16 | `s_security.py` | 디스크의 비밀 문자열(private.pem 본문·credential·토큰) 12–39개가 모든 캡처·로그에 0건이다. 패턴 검색 0건이다. 비친구 회원의 이메일·agent ID가 화면에 없다. |
| 이전 시나리오 회귀 | `regress-ux002/checks.py`, `mutate.py` | [SAR-UX-ISSUES-002-CONNECT-TESTER](SAR-UX-ISSUES-002-CONNECT-TESTER.md)의 재개·만료·취소·device 승인 뒤 미저장 행이 그대로다. 규칙이 바뀐 기대(한도·선택지 목록)는 따로 적는다. |

## 시나리오 밖 (인계)

- CS-AC-03(전체 201개·임계 주입)·CS-AC-04의 시계 주입·`service_capacity`: DEV 테스트가 근거다. fixture로 2000개 agent를 만들지 않았다.
- Dots 원격 경로의 기본 이름과 이름 없는 agent 생성 여부: fixture에 원격 OAuth 클라이언트가 없다.
- Claude Code·Codex 도구의 보낸·받은·상태 출력 이름(CS-AC-09 일부)과 실제 Grok Bot·Linux `install_bot_mcp.sh` 실행.
- 실제 Google 화면의 계정 선택·실패 문구, 1440·390·320px·키보드(CS-AC-17): designer 시각 검수(ART).
- 실연결 재검증(CS-AC-18): 운영 배포 뒤 같은 과제의 2단계로 따로 배정한다. 1단계 PASS와 섞지 않는다.
