---
title: SAR-SESSION-PEER-001-CODE-REVIEW-3 리뷰
status: review
updated: 2026-10-17
owner: reviewer
summary: df525350b8d3 병합 전 독립 재리뷰 보고서. FIX-2의 같은 세션 subagent deny와 V04 해석 문서를 확인했고 새 결함 없이 수락, 이전 low 3건은 범위 밖 후속으로 승계
---

# SAR-SESSION-PEER-001-CODE-REVIEW-3 리뷰

- 검토자 / CLI / 모델: coor가 배정한 독립 reviewer `sp001-code-review-3-opus` / Claude Code 2.1.302 / `claude-opus-5-5`. 실제 host 세션 ID는 `4d91c0a1-d395-4484-a388-15937008030c`이다. result.json `independence.reviewer_session`에는 snapshot 생성 기록과 같은 값 `sp001-code-review-3-opus`를 쓴다. 구현자 세션은 `ctx_a900924558cd`이다. 이전 reviewer(`sp001-code-review-opus`, `sp001-code-review-2-opus`)와 다른 세션이다.
- base SHA / head SHA / merge-base: `37d491c1c6a94ed48c4c1d6b892dd5a0f9c19916` / `df525350b8d38a86c92a4f8aeb5c37449741ff9a` / `37d491c1c6a94ed48c4c1d6b892dd5a0f9c19916`.
- 검토 공간: 읽기 전용 detached snapshot `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-fc0beaa674b84871832cd6135baa798b`(HEAD `df52535`, `git status --porcelain` 0줄). snapshot은 수정하지 않았다. 결과는 coor 기록 체크아웃에 썼다.
- OCR 버전 / 적용 규칙: open-code-review v1.13.0 (1672852). `rules.json`(rule_sha256 `f1061481…de8`, CODE-REVIEW-2와 같음)과 `.fullops-squad/rules/common/`, `project.md`. `999c607..df52535`에서 규칙 파일 변경은 없다. OCR 자동 판정이 아니라 AI 검토 결과다.
- 요구사항·완료 기준 원천: 설계 SAR-SESSION-PEER-001 SP-04(사용자 결정 2026-10-17로 같은 세션 subagent 포함)와 SP-01–12, DEV-FIX-2 지시서와 완료 보고(`handovers/logs/2026-10-17_to_dev.md`의 `SAR-SESSION-PEER-001-DEV-FIX-2` 절), 직전 수락 리뷰 [CODE-REVIEW-2](../SAR-SESSION-PEER-001-CODE-REVIEW-2-review/report.md)(999c607 기준).
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 59 / 38 / 21 / 18 / 41. 이번 수정 범위(`999c607..df52535`)는 19파일이며 코드는 `register.ts`(+5줄)·`register.test.ts`(+14줄) 2개다. 나머지는 D10·D11·adapters README·DEV QA 보고와 증거·하네스 기록이다. 지시서의 승인 범위(adapters/hooks/·관련 테스트·문서) 안이다. relay·Codex·Hermes·`adapters/src/` 변경은 없다.
- lint(`lint.json`) ERROR / WARNING / 실행 불가: 0 / 5 / 0. `product-lint`(`make lint`)·`product-test`(`make test`) 모두 exit 0, head `df52535`, test_level lite.
- SIZE-001: `core.ts` 337, `mcp.ts` 401, `text.ts` 372줄. FIX-2는 세 파일을 바꾸지 않았다. 기존 판단(기존 초과·차단 아님)을 유지한다. 지시서도 범위 밖으로 둔다. `register.ts`는 244줄로 상한 안이다.
- SIZE-002: 추가 811줄(CODE-REVIEW-2 시점 784 + FIX-2 코드·문서 27). 늘어난 코드는 deny 분기 5줄과 테스트 14줄이다. 분할 필요가 없다.
- DEP-001: `adapters/package.json`은 FIX-2에서 바뀌지 않았다. main 대비 scripts만 바뀌고 dependencies·devDependencies 변경이 없다(CODE-REVIEW-2 확인 승계).
- UI 디자인: 웹 화면 변경 없음(지시서 명시). 디자인 lint는 해당 없음.

## 검토 범위

중점(999c607..df52535):
- `adapters/hooks/register.ts:192-199`: KnowsLink text 도구(`TOOL` 정규식이 send/receive/receipt 모두 포함)이고 `e.agentId`가 있으면 연결 확인·승인 대화상자·`next` 전에 deny한다. 비 KnowsLink 도구는 그 앞에서 `next(e)`로 통과하므로 다른 도구에 영향이 없다. 예외가 나면 기존 `.catch`가 fail-closed로 막는다.
- `agentId` 판별 근거: 이 세션의 Claude Code 2.1.302 모드 API 타입(`plugin-authoring/types/claude-code.d.ts`)을 직접 읽었다. `ToolCallInput = ToolCallEnvelope & AgentLoop & ToolRequestMeta`(13304행)이고, `AgentLoop.agentId`(194행)는 "subagent·teammate 루프에서는 그 id, 메인 루프에서는 없음"이다. 또 "Pinned: a different value is refused, one left out is kept"이므로 위쪽 hook이 `next({...e})` 재작성으로 값을 지우거나 바꿀 수 없다. workflow agent와 엔진 자체 fork도 id를 가지므로 deny 쪽(fail-closed)으로 처리된다. KnowsLink text 도구 인자(`peer`·`text`·`reply_to` 등, `mcp.ts`)에 `agentId`가 없어 envelope 필드와 충돌하지 않는다.
- 실제 세션 증거: `fix2/subagent-deny-session.txt`와 `fix2/transcript-tool-results.txt`에서 subagent(`isSidechain=True`, `agentId=a71d8a5d…`)의 `knowslink_text_receive`가 `is_error=True`와 deny 문구로 끝났고, 메인(`agentId=None`)의 같은 호출은 `received`였다. 별도 프로세스 자식은 자기 모드 인스턴스의 `link`가 null이라 기존 미연결 deny를 받는다(코드 경로 확인).
- `register.test.ts` 새 테스트: 연결 fixture에서 `agentId` 있는 send·receive deny, 대화상자 0, 같은 상태에서 메인 send(`sent` 1)·receive 허용. 기존 7건은 그대로이며 기대값 완화·삭제가 없다. receipt는 같은 정규식·분기라 별도 단언이 없다(위험 낮음, finding 아님).
- 문서: D10(module-design) register.ts 행·SP-04 해석·V04 해석, D11(user-guide) 7·8단계, adapters README 7·8단계, DEV QA 보고 FIX-2 절을 코드와 사용자 결정(차단 / 현재 동작 유지)에 대조했다. README의 거부 문구는 코드 원문과 같고 D11은 핵심 부분을 인용한다. V04 해석은 코드 변경 없이 tester 관측(fullops/tester a32d0ab V04)을 근거로 한다. SP-04 원문("소유자가 자식을 별도 agent로 연결한 경우에만")과도 모순이 없다. 같은 프로세스 subagent는 별도로 `on`할 수 없기 때문이다.

전체 diff(main 37d491c 대비) 회귀 확인: FIX-2가 바꾸지 않은 `hooks.json`·`package.json`·`SKILL.md`·`core.ts`·`mcp.ts`·`session-lock.ts`(+test)·`text.ts`·`package_plugin.py`·D03·project.md·기술 계획은 999c607 이후 변경이 없다. CODE-REVIEW-2의 검토를 승계하고 새 분기와의 상호작용만 다시 확인했다. 새 분기는 deny만 추가하므로 미연결 deny, 호출별 승인, 4096자 초과 deny, replaced 즉시 해제, 수신 뒤 대기 수 갱신, toast 본문 비노출 경로는 그대로다. 미연결 세션의 subagent가 미연결 문구 대신 subagent 문구를 받는 차이만 있고, 둘 다 deny라 영향이 없다.

skipped 41건과 확인 방법:
- `docs/evaluations/jev/`(history 포함) 18건: 하네스 자동 생성 탐색 기록. 수락 판단에 쓰지 않는다.
- `handovers/logs/2026-10-17_to_dev.md`: DEV-FIX-2 지시서의 해야 할 일·적용 기준·완료 기준·범위 밖과 완료 보고 절만 직접 읽었다.
- FIX-2 증거 `fix2/` 6건은 모두 직접 읽었다(result.json reason에 파일별 내용). `adapters-check-win.log`의 typecheck는 `tsconfig.json` include가 `src/**/*.ts`라 hooks를 검사하지 않는다. hooks의 `e.agentId` 접근은 위의 d.ts 계약과 `claude plugin test` 로드·통과로 확인했다.
- 이전 DEV·FIX 증거 16건(`fix/` 11건, 루트 로그·png·fixture 5건): 999c607 이후 변경 없음. CODE-REVIEW-2의 확인을 승계한다.

## 발견 사항

새 발견 사항은 없다.

CODE-REVIEW-2의 미해결 low 3건은 FIX-2 지시서가 "리뷰 low 3"을 범위 밖으로 명시했으므로 그대로 승계한다(result.json findings, `resolved: false`).

| # | 심각도 | 파일·줄 | 내용 | 상태 |
|---|---|---|---|---|
| 7 | low | `adapters/hooks/register.ts:222-233` | 대화상자가 열린 동안 replaced·off로 연결이 풀려도 `보내기`면 발송된다 | 미해결, 후속 |
| 8 | low | `adapters/hooks/register.ts:17` | 보이지 않는 변이 선택자(Mn)가 대화상자에 그대로 들어간다 | 미해결, 후속 |
| 9 | low | DEV QA `report.md:70` | 최대 길이 본문 대화상자 표시 미관측 | 미해결, tester V05/V06 |

줄 번호는 df52535 기준으로 다시 맞췄다(#7은 deny 분기 5줄 추가로 5줄 밀림).

## 검증 및 남은 제약

- 재사용 증거: `lint.json`(attempt `3393904f…`, base `37d491c`, head `df52535`, config_sha256 `9c49bb2d…`, policy_sha256 `6d797d34…`, product-lint·product-test exit 0). prepare가 준비한 결과이며 DEV 지시서에 따라 Linux 격리 clone에서 실행됐다. 환경 동일성은 자동 판정하지 않는다. coor 지시에 따라 lint.py·전체 테스트를 재실행하지 않았다.
- FIX-2 자체 증거: `fix2/mod-test.log` hooks 8 pass(src 10행은 기존 kit 로드 실패로 판정 제외), `fix2/mod-test-package.log` 8 pass 0 fail, `fix2/adapters-test-win.log` PASS.
- 미검증: foreground subagent·teammate 실제 세션은 관측하지 않았다. d.ts 계약은 loop 단위라 background 여부와 무관하다. receipt 도구 deny 단언 없음(같은 분기).
- 남은 일: low 3건 후속(coor가 PLANS.md에 유지), tester 재확인, main 통합은 coor 담당이다.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과 | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `design-docs/SAR-SESSION-PEER-001.md` SP-04, D10 'SP-04 해석'·'V04 해석' 문단, DEV-FIX-2 지시서 '현재 상황' | 확인 | 사용자 결정 날짜·내용이 D10·README·D11·DEV QA에 같게 기록됨. SP-04 표 자체에는 확장 문구가 없고 D10이 해석을 담는다 |
| 구조와 구현/미완료 상태 | D10 register.ts 행, DEV QA FIX-2 절 | 확인 | 코드와 일치 |
| 실행·검증 방법과 증거 | DEV QA FIX-2 표와 `fix2/` 증거 | 확인 | foreground subagent 미관측을 한계로 적음 |
| 운영·복구(해당 시) | adapters README 7·8단계·한계 문단, D11 7·8단계 | 확인 | 운영 배포는 범위 밖 |
| 다음 작업·담당·재개 조건 | DEV-FIX-2 완료 보고 '남은 일' | 확인 | 재리뷰(이 문서)·tester 재확인·main 통합은 coor |
| 로컬 링크·절 접근/지원 한계 | DEV QA report.md → `fix2/*` 상대 링크 | 확인 | 링크 대상 6건이 모두 커밋에 있음 |
| snapshot 정리 후 정본 접근 | 이 폴더(coor 기록 체크아웃)의 result.json·report.md·lint.json | 확인 | snapshot 경로 없이 열람 가능 |

## 검토 결론

수락한다. 대상은 SHA `df525350b8d38a86c92a4f8aeb5c37449741ff9a`이다. 같은 세션 subagent·teammate의 KnowsLink text 도구 deny는 사용자 결정(SP-04 확장)대로 동작하고, `e.agentId` 판별은 Claude Code 2.1.302 모드 API 계약과 실제 세션 관측으로 뒷받침된다. V04 해석 문서는 사용자 결정(현재 동작 유지)과 tester 관측에 맞다. main 대비 전체 diff에서 회귀는 없다. CODE-REVIEW-2의 low 3건은 후속으로 남는다. `review.py check`는 기록 검사이며 AI 검토 내용이나 테스트 성공을 자동 보증하지 않는다.
