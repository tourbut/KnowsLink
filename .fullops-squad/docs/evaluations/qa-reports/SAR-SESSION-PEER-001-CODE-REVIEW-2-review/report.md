---
title: SAR-SESSION-PEER-001-CODE-REVIEW-2 리뷰
status: review
updated: 2026-10-17
owner: reviewer
summary: 999c6079a778 병합 전 독립 재리뷰 보고서. 이전 high 1·medium 2·low 3 해소 확인, 새 low 3건과 함께 수락
---

# SAR-SESSION-PEER-001-CODE-REVIEW-2 리뷰

- 검토자 / CLI / 모델: coor가 배정한 독립 reviewer `sp001-code-review-2-opus` / Claude Code / `claude-opus-5-5`. 실제 host 세션 ID는 `726edd1b-874f-4a51-8267-7b490706b854`이다. result.json `independence.reviewer_session`에는 snapshot 생성 기록과 같은 값 `sp001-code-review-2-opus`를 쓴다. 구현자 세션은 `ctx_4bc77f1f98fc`이다. 이전 리뷰 reviewer(`sp001-code-review-opus`)와 다른 세션이다.
- base SHA / head SHA / merge-base: `37d491c1c6a94ed48c4c1d6b892dd5a0f9c19916` / `999c6079a7780f1284da89c187cf1a9d593c9f1d` / `37d491c1c6a94ed48c4c1d6b892dd5a0f9c19916`.
- 검토 공간: 읽기 전용 detached snapshot `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-a5f301afdc044ac8a32f787012ca4e0a`(HEAD `999c607`, `git status --porcelain` 0줄). snapshot은 수정하지 않았다. 결과는 coor 기록 체크아웃에 썼다.
- OCR 버전 / 적용 규칙: open-code-review v1.13.0 (1672852). `rules.json`(rule_sha256 `f1061481…de8`, 이전 리뷰와 같음)과 `.fullops-squad/rules/common/`(README·coding-style·testing·security), `project.md`. `bdc6500..999c607`에서 규칙 파일 변경은 없다. OCR 자동 판정이 아니라 AI 검토 결과다.
- 요구사항·완료 기준 원천: 설계 SAR-SESSION-PEER-001 SP-01–12, 이전 리뷰 [report.md](../SAR-SESSION-PEER-001-CODE-REVIEW-review/report.md)·result.json, DEV-FIX 지시서와 완료 보고(`handovers/logs/2026-10-17_to_dev.md`의 `SAR-SESSION-PEER-001-DEV-FIX` 절), DEV QA 보고의 'FIX 재검증' 절.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 47 / 32 / 15 / 18 / 29. 수정 커밋 범위(`bdc6500..999c607`)는 25파일이며 코드는 `register.ts`·`register.test.ts`·`package_plugin.py` 3개다. 지시서 범위(고칠 항목 6건, D10·D11·DEV QA 보완) 안이다. relay·Codex·Hermes 변경은 없다.
- lint(`lint.json`) ERROR / WARNING / 실행 불가: 0 / 5 / 0. `product-lint`(`make lint`)·`product-test`(`make test`) 모두 exit 0, head `999c607`, test_level lite.
- SIZE-001: `core.ts` 337, `mcp.ts` 401, `text.ts` 372줄. 수정 커밋은 세 파일을 바꾸지 않았다. 이전 리뷰 판단(기존 초과·차단 아님)을 유지한다. `register.ts`는 239줄로 상한 안이다.
- SIZE-002: 추가 784줄(이전 605 + FIX 테스트 위주 179). 지시서가 회귀 테스트 추가를 요구했고 늘어난 줄은 대부분 `register.test.ts` fixture다. 분할하지 않은 이유를 수락한다.
- DEP-001: `adapters/package.json`은 FIX에서 바뀌지 않았다. scripts만 바뀌고 dependencies·devDependencies 변경이 없음을 전체 diff로 다시 확인했다.
- UI 디자인: 웹 화면 변경 없음. 승인 대화상자 문구는 대상·새 메시지/답장·전체 글자 수·본문을 구분해 보인다. frontend.json tools가 비어 디자인 lint는 해당 없음(지시서 명시).

## 검토 범위

reviewed 18건: 코드 `register.ts`(전체 파일)·`register.test.ts`·`package_plugin.py`, FIX가 바꾸지 않은 `hooks.json`·`package.json`·`SKILL.md`·`core.ts`·`mcp.ts`·`session-lock.ts`(+test)·`text.ts`, 문서 D03·D10·D11·project.md·adapters README·기술 계획·DEV QA report.md. FIX가 바꾸지 않은 파일은 main `37d491c` 대비 전체 diff로 FIX와의 상호작용(replaced 이벤트 생성 시점, TEXT_MAX와 zod `max(4096)`·peer 형식 검증)을 확인했다.

skipped 29건과 확인 방법:
- `docs/evaluations/jev/` 12건: 하네스 자동 생성 탐색 기록. 수락 판단에 쓰지 않는다.
- `handovers/logs/2026-10-17_to_dev.md`: DEV-FIX 지시서의 고칠 항목·완료 기준과 완료 보고 절만 직접 읽었다.
- FIX 증거 `fix/` 11건 중 `dialog-full-body.txt`·`replaced-deny.txt`·`debug-ui-lines.txt`·`watch-busy.out`·`watch-take*.out`·`mod-test.log`·`adapters-test-win.log`를 직접 읽었다. `adapters-check-win.log`는 lint.json의 등록 gate로 대체했고 `mod-test.full.log`는 요약 로그로 대체했다.
- 이전 DEV 증거 5건(`fixture-agents.json`, `toast-a-terminal.png`, 로그 3건): bdc6500 이후 변경 없음. 이전 리뷰의 확인을 승계한다.

## 이전 발견 사항 해소 확인

| # | 심각도 | 이전 내용 | 확인 결과 |
|---|---|---|---|
| 1 | high | 승인 대화상자가 본문을 200자에서 표시 없이 자름 | 해소. `register.ts:208-221`이 본문 전체와 `본문 전체 N자`를 보이고 4096자 초과는 묻지 않고 deny한다. 정상 본문은 4096 bytes 이하이므로 UTF-16 길이도 4096 이하여서 유효한 본문이 거부되지 않는다. 테스트가 201자 이상 본문 전체·글자 수·4097자 deny를 단언한다. 255자 실제 대화상자 관측이 있다 |
| 2 | medium | replaced 뒤 link 미해제 | 해소. `register.ts:129-139`가 `link === current`일 때 즉시 해제·상태줄 초기화·자식 정지를 한다. 이후 `finally`는 `link !== current`라 중복 toast가 없다. 정지된 자식이 값 없이 끝날 때 `ended.code` 접근으로 생기던 미처리 rejection을 `?? { code: null }`이 막는다(`off` 경로에도 이득). 테스트와 실제 세션 관측(`replaced-deny.txt`, `debug-ui-lines.txt` 04:10:19)이 있다 |
| 3 | medium | 연결 세션의 취소·ask 실패·인자 누락 자동 검사 부재 | 해소. 연결 fixture로 취소·dialog 예외·Other 자유 입력·인자 누락이 MCP send에 도달하지 않음(`sent` 0)을 검사한다. 악성 외부 본문 수신 뒤 Bash·send·대화상자·추가 spawn 0, toast에 본문 없음, 발신자 제어문자 제거를 검사한다. 모델이 지시를 따르지 않는지는 자동화 불가 사유를 DEV QA 보고와 D10에 남겼다 |
| 4 | low | reply_to 미표시 | 해소. `새 메시지` 또는 `답장(reply_to <ID>)` |
| 5 | low | smoke 서술과 로그 불일치 | 해소. 서술을 로그 3줄 범위로 고치고 `fix/watch-busy.out`(busy 한 줄, `exit 2`)을 추가 |
| 6 | low | `package_plugin.py` 주석 명령 이름 | 해소 |

## 새 발견 사항

| # | 심각도 | 파일·줄 | 내용과 영향 | 수정 요청 |
|---|---|---|---|---|
| 7 | low | `adapters/hooks/register.ts:217-228` | `link` 확인은 대화상자를 열기 전에만 한다. 대화상자가 열린 동안 `replaced`나 `/knowslink-session off`로 연결이 풀려도 소유자가 `보내기`를 고르면 발송된다. 소유자가 본 본문 그대로이므로 SP-05/06은 지켜진다. 그러나 README·D03·D11의 "이후 그 세션의 text 도구는 모두 거부된다"와 어긋난다. 재현: 연결 세션에서 send 대화상자를 연 채 다른 세션이 `--take` → 10초 안 해제 toast → 대화상자에서 `보내기` → 발송 | ask 뒤 `link`가 같은 Link가 아니면 deny한다 |
| 8 | low | `adapters/hooks/register.ts:17` | `CONTROL`은 Cc·Cf만 공백으로 바꾼다. 변이 선택자(U+FE00–FE0F, U+E0100–E01EF, 범주 Mn)는 대화상자에서 보이지 않는다. 외부 본문에 조종된 모델이 이 문자로 데이터를 숨겨 붙이면 소유자는 `전체 N자` 숫자 외에는 알기 어렵다. 태그 문자(Cf)는 이미 공백으로 바뀐다 | 보이지 않는 결합 문자를 가시 표식으로 바꾸거나 포함 시 경고·deny한다 |
| 9 | low | DEV QA `report.md:70` | 전체 본문 표시는 255자로만 관측했다. 최대 길이(4096자, 한글 약 1365자) 본문이 AskUserQuestion에서 잘리지 않는지는 관측하지 않았다. UI가 긴 질문을 자르면 #1이 UI 계층에서 재현된다 | tester V05/V06에서 최대 길이 본문 대화상자를 한 번 관측해 증거로 남긴다 |

수정으로 생긴 차단 회귀는 없다. 확인한 회귀 후보:
- `current.stop()` 뒤 같은 버퍼의 남은 이벤트는 처리되지만 `link`가 null이라 text 도구 통과로 이어지지 않는다. 옛 세션의 `on` 재실행은 새 watcher가 busy로 막힌다(SP-02 유지).
- 미연결 deny, `.catch` fail-closed, toast 본문 비노출, 호출별 승인, 자동 답장 0, permission relay 미사용은 그대로다.
- 테스트 삭제·감소·기대값 완화는 없다. 기존 3건이 유지되고 4건이 추가됐다.

## 검증 및 남은 제약

- 재사용 증거: `lint.json`(attempt `d3a594d0…`, base `37d491c`, head `999c607`, config_sha256 `9c49bb2d…`, policy_sha256 `6d797d34…`, product-lint·product-test exit 0). prepare가 준비한 결과이며 DEV 지시서에 따라 Linux 격리 clone에서 실행됐다. 환경 동일성은 자동 판정하지 않는다. coor 지시에 따라 lint.py·전체 테스트를 재실행하지 않았다.
- FIX 자체 증거: `fix/mod-test.log` hooks 7 pass, `fix/adapters-test-win.log` PASS 10줄.
- 관측 이상: 첫 takeover 시도에서 Claude 세션 `c0d543f9`가 takeover watcher의 `timeout` 종료 시각에 함께 끝났고 원인은 미확인이다. 두 번째 시도(background watcher)에서는 재현되지 않았다. 모드 코드에서 세션 종료를 일으킬 경로는 찾지 못했다. tester가 takeover를 반복할 때 같이 관찰하기를 권한다.
- 미검증: 모델이 외부 본문 지시를 따르지 않는지(V05 수동), V03/V04/V07, Codex·Hermes, 운영 relay. V01–07 독립 QA는 tester 담당이다.
- 남은 일: low 3건은 후속 과제 또는 tester 재검증 범위로 PLANS.md에 남긴다(coor). 수락을 막지 않는다.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과 | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | `design-docs/SAR-SESSION-PEER-001.md` SP-01–12, DEV-FIX 지시서 '고칠 항목' | 확인 | 수정 이유가 이전 리뷰 경로와 함께 기록됨 |
| 구조와 구현/미완료 상태 | D03 'Claude Code 작업 세션 연결' 절, D10 register.ts 행, project.md | 확인 | D03의 replaced 즉시 해제 서술이 코드와 일치. 대화상자 중 교체 예외는 finding 7 |
| 실행·검증 방법과 증거 | D10 V05 문단, DEV QA 'FIX 재검증' 절과 `fix/` 증거 | 확인 | 최대 길이 대화상자 증거 없음(finding 9) |
| 운영·복구(해당 시) | adapters README 한계 문단, D11 3·6단계 | 확인 | 운영 배포는 범위 밖 |
| 다음 작업·담당·재개 조건 | DEV-FIX 완료 보고 '남은 일' | 확인 | 재리뷰(이 문서)·tester 재검증·main 통합은 coor |
| 로컬 링크·절 접근/지원 한계 | DEV QA report.md → `fix/*` 상대 링크 | 확인 | 링크 대상 파일이 모두 커밋에 있음 |
| snapshot 정리 후 정본 접근 | 이 폴더(coor 기록 체크아웃)의 result.json·report.md·lint.json | 확인 | snapshot 경로 없이 열람 가능 |

## 검토 결론

수락한다. 대상은 SHA `999c6079a7780f1284da89c187cf1a9d593c9f1d`이다. 이전 리뷰의 high 1·medium 2·low 3이 모두 해소됐다. 핵심 신뢰 경계는 약해지지 않았고 수정이 만든 차단 회귀는 없다. 새 low 3건은 후속으로 다룬다. `review.py check`는 기록 검사이며 AI 검토 내용이나 테스트 성공을 자동 보증하지 않는다.
