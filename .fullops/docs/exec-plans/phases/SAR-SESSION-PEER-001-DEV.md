---
title: Claude Code 작업 세션 연결 기술 구현
status: review
updated: 2026-10-17
owner: dev
tasks: [SAR-SESSION-PEER-001-DEV]
summary: Claude Code 플러그인 모드와 세션 결속 watcher로 SP-01–12의 Claude↔Claude 왕복을 구현하고 격리 relay에서 검증한다.
---

# 기술 계획

기준 HEAD는 `283ea63`(main `37d491c`)이다. 공통 기준은 fullops-common-0.3.3, 제품 정본은 [SP-01–12](../../design-docs/SAR-SESSION-PEER-001.md)와 HP-01–16이다. 테스트 lite, 하위 위임 standard이지만 변경이 adapters 한 곳에 모여 하위 에이전트는 쓰지 않았다.

| 판단 | 결정과 근거 |
|---|---|
| 호스트 버전 | Claude Code `2.1.302`(Windows 11, Node 22.22.2). plugin-authoring 스킬이 제공한 같은 빌드의 `claude-code.d.ts`로 API를 확인했다. Context7은 쓰지 않았다. |
| 수신 방식 | Channels가 아니라 기존 relay pull이다. Channels는 research preview·조직 정책·세션 opt-in·permission relay를 포함해 SP-07을 만족하기 어렵고 설치 조건이 더 많다. 모드가 `plugin.js watch <폴더> --session <구분값>` 자식 프로세스를 띄우고 그 stdout 이벤트만 읽는다. 본문은 이벤트에 없다. |
| `knowslink on` 형태 | 슬래시 명령 `/knowslink-session on|off|status [--take]`. `/knowslink`는 패키지 skill이 이미 `/knowslink:knowslink`로 차지해 등록이 거부됐다. |
| toast 경로 | 모드의 `$.ui.toast`(12초)와 `$.ui.status`. 도착 알림이며 '세션 수신 확인'이 아니다. |
| 세션 결속 | `<폴더>/session.lock` `{session,pid,since}` 배타 claim(`session-lock.ts`). 살아 있는 다른 세션은 busy, `--take`만 교체, 교체된 watcher는 10초 lock 감시로 자진 종료, 죽은 pid는 회수. 구분값은 세션 ID 앞 8자다. |
| 승인 경계 | `tool.call` hook이 미연결 세션의 `knowslink_text_*`를 deny한다. 발송은 호출마다 소유자 `AskUserQuestion`(`보내기`/`취소`)을 거친다. `-p`·취소·인자 누락은 deny. |
| 켜지 않은 세션 | Claude 설치의 `--connect`가 `KNOWSLINK_AUTO_RECEIVE=off`를 쓰므로 MCP 프로세스는 스스로 pull하지 않는다. Codex 설정은 바꾸지 않았다. |
| 지원 수준 | '기존 세션 수신'(도착 toast·세션 결속 로컬 사본; 작업 중·대기 양쪽 관측). 본문 모델 표시는 '다음 조회'(사용자 요청 → `knowslink_text_receive`). |

새 의존성은 없다. 변경 파일은 `adapters/src/session-lock.ts`(+test), `adapters/src/text.ts`, `adapters/src/mcp.ts`, `adapters/src/core.ts`, `adapters/hooks/{hooks.json,register.ts,register.test.ts}`, `adapters/package.json`, `scripts/package_plugin.py`, 문서(adapters README·SKILL·D03·D10·D11·project)다. relay Go 코드·SQL·compose는 바꾸지 않았다.

## 구현 결과와 추적

| 규칙 | 구현·검증 경로 |
|---|---|
| SP-01 선택 | 설치·Google 연결만으로 수신하지 않음(`KNOWSLINK_AUTO_RECEIVE=off`). `on` 전 text 도구 deny(`register.test.ts` 1, 세션 A off 뒤 호출 거부 관측). |
| SP-02 정확한 대상 | lock claim. CLI smoke: 두 번째 `--session`은 `knowslink_session_busy` exit 2. `--take`로 교체 뒤 옛 watcher `knowslink_session_replaced` 종료. |
| SP-05/06 명시 발송·답장 승인 | 모든 send에 대화상자. V01에서 A→B 1회, B→A 3회 모두 `보내기` 선택 뒤에만 queued. 인자 누락 호출은 대화상자 없이 deny. |
| SP-07 신뢰 경계 | 모델은 `knowslink_text_receive` 결과(`untrusted:true`)로만 본문을 본다. 두 세션 모두 본문 지시 실행·자동 답장 0. |
| SP-08 작업 보존 | B가 100초 Bash 작업 중(1분 26초 시점) 도착: 상태줄 `대기 1`, 작업은 1분 48초에 정상 완료. |
| SP-09 연결 수명 | `off` 뒤 수신 도구 deny 관측. 모듈 reload로 자식이 끝나면 미연결로 보고하고 재 `on`이 필요하다. |
| SP-10 한도 | 원 메시지 TTL 뒤 `reply_to` 답장은 relay `422 expired`로 거부됐고 새 메시지로 보냈다. 4096/4097은 기존 `mcp.test.ts`. |
| SP-11 상태 | 상태줄 `연결|대기 N`, toast에 발신 agent·대기 수·다음 행동만. 작업 경로·프롬프트는 공유하지 않는다. |
| SP-12 근거 | 격리 relay·합성 회원 A/B·Claude Code 두 세션. 실제 운영 relay·실계정·Codex·Hermes는 미실행. |

## 실행 증거와 한계

[DEV QA 보고](../../evaluations/qa-reports/SAR-SESSION-PEER-001-DEV-test/report.md)에 메시지 ID·시각·세션 구분값·toast 캡처·로그를 보존했다. 상세 한계는 보고서에 있다.

- 한 번의 모드 소스 교체가 열린 세션을 hot reload해 자식 watcher가 끝났다. 제품 설치(플러그인 폴더 고정)에서는 일어나지 않지만 문서에 재 `on` 절차를 적었다.
- 첫 발송 시도에서 모델이 `peer_agent_id`를 써 빈 peer 대화상자가 떴다. 인자 누락 deny를 추가해 재현되지 않았다.
- 상태줄 대기 수는 도착·수신 때 갱신한다. 다른 프로세스가 읽은 메시지는 반영하지 않는다.
- Windows는 `off`·종료 때 lock 파일을 지우지 않는다. 죽은 pid 회수로 처리한다.

## 보존 자원

검증 자원은 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/session-peer-001-dev/`(설치 두 개·private 폴더·fixture 스크립트·로그·원본 캡처)와 Docker `knowslink-sp001-dev` 프로젝트(Postgres·relay, 정지 상태)·`knowslink-sp001-dev-mailsink2` 컨테이너다. 삭제·이동은 하지 않았다. 실제 회원·운영 서버·사용자 전역 Claude 설정은 바꾸지 않았다.
