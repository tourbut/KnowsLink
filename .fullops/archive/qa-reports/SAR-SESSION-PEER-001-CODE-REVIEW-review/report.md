---
title: SAR-SESSION-PEER-001-CODE-REVIEW 리뷰
status: review
updated: 2026-10-17
owner: reviewer
summary: bdc650004647 병합 전 독립 리뷰 보고서. 발송 승인 본문 절단(high)으로 수정 요청
---

# SAR-SESSION-PEER-001-CODE-REVIEW 리뷰

- 검토자 / CLI / 모델: coor가 배정한 독립 reviewer `sp001-code-review-opus` / Claude Code / `claude-opus-5-5`. 실제 host 세션 ID는 `b04c7447-9cde-431b-a0f9-159dcf137ac1`이다. result.json `independence.reviewer_session`에는 snapshot 생성 기록과 같은 값 `sp001-code-review-opus`를 쓴다. 구현자 세션은 `ctx_7d74d5471bb0`이다.
- base SHA / head SHA / merge-base: `37d491c1c6a94ed48c4c1d6b892dd5a0f9c19916` / `bdc650004647462cccb0aa69debb3397b7c43b3c` / `37d491c1c6a94ed48c4c1d6b892dd5a0f9c19916`.
- 검토 공간: 읽기 전용 detached snapshot `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-e58ed2bd0af54bfba9899548a54fd965`(HEAD `bdc6500`, clean). snapshot은 수정하지 않았다. 결과는 coor 기록 체크아웃에 썼다.
- OCR 버전 / 적용 규칙: open-code-review v1.13.0 (1672852). `rules.json`(rule_sha256 `f1061481…de8`)과 `.fullops-squad/rules/common/`(README·coding-style·testing·security), `project.md`. OCR 자동 판정이 아니라 AI 검토 결과다.
- 요구사항·완료 기준 원천: [SAR-SESSION-PEER-001](../../../design-docs/SAR-SESSION-PEER-001.md) SP-01–12·V01–07·'하나의 DEV 인계', DEV 지시서와 완료 보고(`handovers/logs/2026-10-17_to_dev.md`의 `SAR-SESSION-PEER-001-DEV` 절), [DEV QA 보고](../SAR-SESSION-PEER-001-DEV-test/report.md), [기술 계획](../../../exec-plans/phases/SAR-SESSION-PEER-001-DEV.md).
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 31 / 26 / 5 / 18 / 13. 지시서 승인 범위(adapters·scripts·mod·테스트·기술/설치 문서) 안이다. relay Go·SQL·compose·internal은 바뀌지 않았다.
- lint(`lint.json`) ERROR / WARNING / 실행 불가: 0 / 5 / 0. `product-lint`(`make lint`)·`product-test`(`make test`) 모두 exit 0, head `bdc6500`, test_level lite. prepare가 등록 worktree의 같은 base·HEAD·설정 결과를 복사했다. 재실행하지 않았다.
- SIZE-001: `core.ts` 337, `mcp.ts` 401, `text.ts` 372줄(상한 300). core·mcp는 4줄 이하 증가로 기존 초과다. `text.ts`는 watch 확장으로 53줄 늘었다. watch 세션 로직을 `session-lock.ts` 쪽으로 옮길 여지가 있으나 차단 사유는 아니다.
- SIZE-002: 지시서 예상 '큰 규모일 수 있다' / 실제 코드 추가 605줄 / 모드·lock·watch·문서가 한 기능으로 결속돼 분리하지 않았다는 완료 보고 사유를 수락한다.
- DEP-001: `adapters/package.json`은 `scripts.format:check`·`scripts.test`만 바뀌었다. dependencies·devDependencies 변경 없음을 diff로 확인했다. 모드는 `claude-code` 런타임 API만 import한다.
- UI 디자인: 웹 화면 변경 없음. Claude Code toast·명령 출력은 짧은 한국어 문구이며 발신 agent·외부 메시지 표시를 구분한다. frontend.json tools가 비어 디자인 lint는 해당 없음(지시서 명시).

## 검토 범위

reviewed 18건: 코드 9건(`register.ts`·`register.test.ts`·`hooks.json`·`session-lock.ts`(+test)·`text.ts`·`mcp.ts`·`core.ts`·`package_plugin.py`), `package.json`, `SKILL.md`, 문서 6건(D03·D10·D11·project.md·adapters README·기술 계획)과 DEV QA report.md. 관련 기존 코드로 `mcp.ts`의 `startAuto`(`KNOWSLINK_AUTO_RECEIVE=off`)와 `knowslink_text_send/receive` 핸들러, `inbox.ts`의 `AutoReceiver.stop`을 함께 읽었다.

skipped 13건과 확인 방법:
- `docs/evaluations/jev/` 6건: 하네스 자동 생성 탐색 기록. 수락 판단에 쓰지 않는다.
- `handovers/logs/2026-10-17_to_dev.md`: 자동 생성 기록이지만 SAR-SESSION-PEER-001-DEV 지시서와 완료 보고 절을 직접 읽었다. 패킷 블록은 열람하지 않았다.
- `fixture-agents.json`, `toast-a-terminal.png`: DEV QA report.md 서술로 존재를 확인했다. 이미지 내용 판정은 하지 않았다.
- 로그 4건: `mod-test.log`(3 pass·src 10개 로드 실패, report 서술과 일치)와 `watcher-smoke.log`(finding 5)를 직접 읽었다. `adapters-check-win.log`·`adapters-test-win.log`는 lint.json의 등록 gate 통과로 대체했다.

## 발견 사항

| # | 심각도 | 파일·줄 | 내용과 영향 | 수정 요청 |
|---|---|---|---|---|
| 1 | high | `adapters/hooks/register.ts:199-212` | 발송 승인 대화상자가 `text`를 `slice(0, 200)`으로 자르고 절단 표시·전체 길이를 보이지 않는다. text는 최대 4096 bytes다. 소유자는 앞 200자만 보고 `보내기`를 고르며, 모델(외부 본문 주입 포함)이 201자 이후에 붙인 내용은 확인 없이 발송된다. SP-05 '본문을 확인한 소유자의 명시 발송'과 SP-06 '그 대상·본문에만 적용'을 만족하지 못한다. 재현: 연결 세션에서 250자 text로 `knowslink_text_send` 호출 → 대화상자에 앞 200자만, 승인 시 250자 전부 발송. | 전체 본문을 보이거나, 표시 한도를 넘으면 `(전체 N자 중 앞 200자)`를 넣고 전체 확인 수단이 없으면 deny한다. 201자 이상 본문의 대화상자 문구를 회귀 테스트로 검사한다. |
| 2 | medium | `adapters/hooks/register.ts:128-132` | `knowslink_session_replaced` 처리가 toast만 띄우고 `link`를 해제하지 않는다. 해제는 자식 프로세스 종료 뒤 `finally`에서만 일어난다. 교체된 watcher가 열린 handle(SSE 65초 timeout, fetch keep-alive 등)로 늦게 끝나면 그동안 옛 세션의 text 도구가 통과한다. `knowslink_text_receive`는 inbox가 비면 MCP 프로세스가 relay에서 직접 pull한다. README·D03의 '옛 세션은 10초 안에 더 받지 않는다'(SP-09)와 다르다. Claude 세션 takeover 뒤 deny는 관측되지 않았다(CLI smoke만). | replaced 처리에서 `link === current`이면 즉시 `link = null`·상태줄 초기화·자식 정지를 수행한다. |
| 3 | medium | `adapters/hooks/register.test.ts:1-49` | 지시서 완료 기준 'V05(외부 본문 명령 미실행, 승인 없는 답장 0)를 자동 테스트로 확인'이 부분 충족이다. 자동 테스트는 미연결 deny만 다룬다. 연결 세션의 `취소`·ask 실패(`-p`)·인자 누락 deny 경로는 자동 검사가 없다. DEV QA 관측도 모두 `보내기`이며 취소 경로와 명령·경로·승인 요청을 담은 외부 본문은 관측하지 않았다. | 연결 상태 fixture(spawn 이벤트 모의 등)로 취소·ask 거부·인자 누락이 `next` 미도달·deny임을 검사한다. 불가능하면 그 이유와 수동 취소 관측을 보고서에 남긴다. |
| 4 | low | `adapters/hooks/register.ts:202-207` | 대화상자에 `reply_to`가 보이지 않아 새 메시지와 답장을 구별할 수 없다(SP-06 '그 대상'). | finding 1 수정 때 reply_to 유무를 함께 표시한다. |
| 5 | low | `SAR-SESSION-PEER-001-DEV-test/watcher-smoke.log:1-3` | 보고서는 `smokeA2 busy exit 2`를 관측했다고 적지만 로그에는 bound/replaced/bound 3줄만 있다. | 두 번째 watcher의 출력·종료코드를 증거에 더하거나 서술을 로그 범위로 줄인다. |
| 6 | low | `scripts/package_plugin.py:125` | 주석의 `/knowslink on`은 실제 명령 `/knowslink-session on`과 다르다. | 주석을 고친다. |

설계 대조 결과(문제없음 확인):
- SP-01: Claude `--connect` 설치는 `KNOWSLINK_AUTO_RECEIVE=off`라 MCP가 스스로 pull하지 않는다. 미연결 세션의 `knowslink_text_*`는 hook이 deny하며 `.catch`가 `next` 전 실패를 deny로 닫는다(fail-closed).
- SP-02: `session.lock`의 `wx` 배타 생성, 살아 있는 다른 세션 busy, `--take`만 교체를 확인했다. take·dead pid 경합 때의 비원자 쓰기와 Windows pid 재사용은 low 한계로 판단했다. 같은 agent의 세션들은 같은 private inbox를 공유하므로 다른 agent로의 오배송은 없다.
- SP-05/06: 모든 send가 호출마다 소유자 대화상자를 거치며 `취소`·ask 실패는 deny다. 자동 답장 경로는 없다. 단 본문 확인은 finding 1.
- SP-07: watcher stdout 이벤트와 toast에는 본문이 없다. toast의 `from`은 제어문자를 지우고 64자로 자른다. Channels·permission relay를 쓰지 않는다. 본문은 기존 `untrusted:true` MCP 결과로만 모델에 들어간다. SKILL.md가 재시도·자동 재연결 금지를 지시한다.
- SP-08: 모드는 interrupt·폴링 재촉을 하지 않고 toast·상태줄만 갱신한다. DEV QA에서 100초 작업 중 도착 후 작업이 정상 완료됐다.
- SP-09: `off`는 `link`를 즉시 해제하고 자식을 정지한다. 이후 text 도구 deny가 관측됐다. takeover 쪽은 finding 2.

## 검증 및 남은 제약

- 재사용 증거: `lint.json`(attempt `a3bfc848…`, base `37d491c`, head `bdc6500`, config_sha256 `9c49bb2d…`, policy_sha256 `6d797d34…`). prepare가 등록 worktree의 `fullops-gate/lint-result.json`에서 같은 base·HEAD·설정으로 통과한 결과를 바이트 그대로 복사했다. DEV 지시서에 따라 Linux 격리 clone에서 실행한 결과다. 환경 동일성은 자동 판정하지 않으며 이 리뷰는 lint.py·전체 테스트를 재실행하지 않았다(coor 지시).
- 테스트 삭제·감소·기대값 완화는 없다. `session-lock.test.js`가 test 목록에 추가됐다.
- 미검증: 연결 세션의 승인 취소·외부 본문 주입(V05), takeover 뒤 옛 세션 deny, V03/V04/V07, Codex·Hermes, 운영 relay. V01–07 독립 QA는 tester 담당이다.
- 남은 일: finding 1(high)·2·3(medium)을 수정한 새 SHA로 재리뷰한다. low 3건은 같은 수정에 포함하기를 권한다.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과 | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `design-docs/SAR-SESSION-PEER-001.md` 제품 규칙·하나의 DEV 인계, 기술 계획 표 | 확인 | Channels 대신 relay pull, `/knowslink-session` 명명 이유가 기록됨 |
| 구조와 구현/미완료 상태 | D03 'Claude Code 작업 세션 연결' 절, D10 SAR-SESSION-PEER-001-DEV 표, project.md 현재 상태 | 확인 | D03의 'replaced 뒤 10초 안 종료' 서술은 finding 2 수정 후 재확인 필요 |
| 실행·검증 방법과 증거 | project.md 모드 검사 명령, DEV QA report.md, adapters README | 확인 | smoke 로그 증거 부족(finding 5) |
| 운영·복구(해당 시) | adapters README 한계 문단(reload 뒤 재 on, Windows lock 회수), D11 6단계 | 확인 | 운영 배포는 범위 밖 |
| 다음 작업·담당·재개 조건 | 기술 계획 '실행 증거와 한계', 지시서 완료 보고 '남은 일' | 확인 | 이 리뷰의 수정 요청이 추가됨. coor가 PLANS.md에 기록 |
| 로컬 링크·절 접근/지원 한계 | D11 → `../../../adapters/README.md#claude-code-작업-세션-연결--sar-session-peer-001-dev`, README → DEV QA report | 확인 | 상대 경로와 GitHub anchor 형식 일치 |
| snapshot 정리 후 정본 접근 | 이 폴더(coor 기록 체크아웃)의 result.json·report.md·lint.json | 확인 | snapshot 경로 없이 열람 가능 |

## 검토 결론

수정 요청이며 현재 SHA `bdc6500`은 수락하지 않는다. 핵심 신뢰 경계(미연결 deny, 본문 비노출 toast, 호출별 승인, 자동 답장 0, permission relay 미사용)는 설계대로 구현됐다. 그러나 승인 대화상자가 본문을 표시 없이 200자로 잘라 SP-05/06의 본문 확인을 보장하지 못한다(high). replaced 뒤 link 미해제와 승인·취소 경로 자동 검증 부재(medium)도 함께 고친다. `review.py check`는 미해결 high로 실패하는 것이 정상이다. check 통과 여부는 AI 검토 내용이나 테스트 성공을 자동 보증하지 않는다.
