---
title: SAR-AUTO-RECEIVE-001-REVIEW 독립 리뷰
status: draft
updated: 2026-10-10
owner: tester
tasks: [SAR-AUTO-RECEIVE-001-TESTER]
summary: 자동수신 고정 SHA af7627d의 독립 snapshot 리뷰와 lite QA 결론
---

# SAR-AUTO-RECEIVE-001-REVIEW 리뷰

- 검토자 / CLI / 모델: tester, Grok 4.7. 검토자 세션 `01a1261a-d89d-7f03-8185-0553132b8265`. 구현자 세션 `c1bdb696-62e1-4806-84e2-5dc48623fad1`과 다르다. OCR delegate는 LLM을 호출하지 않았고 판정은 이 세션이 했다.
- base SHA / head SHA / merge-base: `dda6130d7023903175bbfa4036a7e052011e2c83` / `af7627d8cca54df55e856f48d226f78dbf2b9d0f` / base와 같다.
- snapshot: `C:\Users\shin\orca\workspaces\KnowsLink\.fullops-review-4170cd9ba9f6484496ac46ed76664cb3`. detached HEAD, porcelain 없음, read_only. 제품 코드와 snapshot은 수정하지 않았다.
- OCR 버전 / 적용 규칙: open-code-review v1.12.13. `review/rule.json` sha256 `f1061481dcdb66affe2f0f69a6c2354d276babcf23434dedfe6b622171994de8`. 그룹 1은 markdown, 그룹 2는 json, 그룹 3은 TypeScript다. fullops-common-0.3.3.
- 요구사항·완료 기준 원천: 현재 tester 인박스, DEV 실행 기록, architecture.md 170행, 설치본 MCP SDK 1.32.0.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 32 / 32 / 0 / 32 / 0.
- lint(`lint.json`) ERROR / WARNING / 실행 불가와 사유: 0 / 7 / 0. 출처는 이 세션의 Linux 고유 임시 clone이다. Node v22.22.2, npm 10.9.7, `npm ci` 뒤 plugin `fullops-squad/1.3.0/scripts/lint.py --from dda6130d7023903175bbfa4036a7e052011e2c83`. product-lint와 product-test exit 0. lint.json sha256 `0648aca688a400ae1f5c9c1fa7485adcde77ede21410edb09db2dca1101f1581`. 의존성 없이 실행한 앞선 시도는 prettier와 tsc가 없어 exit 1이었고 통과 증거로 쓰지 않았다. Windows `cloudflared` 부재는 이 Linux 결과와 분리한다. 81c2bec과 318274f의 Linux lint는 head가 달라 재사용하지 않았다.
- SIZE-002: 지시서는 줄 예산을 주지 않았다. lint 추가 1180줄은 자동수신과 선택 wake 한 기능과 FullOps 기록이다. inbox, wake, 수동 receive를 나누면 한 계약이 끊긴다. 경고는 수용하고 수정 요청으로 올리지 않는다.
- DEP-001: `adapters/package.json`의 test 스크립트에 `node dist/inbox.test.js`만 추가됐다. `package-lock.json` blob `e69dc9c0d2f42a584717527a2e090aeaee789052`는 기준과 같다. 새 의존성이 없고 SDK는 기존 1.32.0이다.
- UI 디자인: 새 화면과 테마가 없다. DESIGN 경고는 0이다. 회원 receipt 문장은 기존 HTML 템플릿의 문구 불일치이며 medium 발견으로 남긴다.

## 검토 범위

32개 `(path, status)`를 모두 reviewed로 기록했다. 제외 파일은 없다.

제품 코드 다섯 파일은 snapshot에서 읽었다. `text.ts` `receive`는 수신자, 만료, 서명 검증 뒤에 `store`를 호출하고 그 다음 persist와 ACK를 한다. `inbox.ts`는 같은 ID의 `.json` 또는 `.read`가 있으면 다시 알리지 않고, `take`는 rename으로 claim한다. `mcp.ts`의 `knowslink_text_receive`는 inbox를 먼저 읽고, 비어 있으면 `auto.serial`로 pull해 자동 loop와 수동 receive가 나란히 lease하지 않는다. 알림 data는 `event`, `id`, `from`, `pending`, `next`이고 메시지 원문이 없다. `grok-wake.ts`는 기본 off이고 주소는 `127.0.0.1`이며 prompt에는 ID와 대기 수만 들어간다.

문서와 JSON은 자동수신 절, 완료 보고, 패킷 정체성을 대조했다. 과거 QA 문자열은 이 계약의 수락 근거로 쓰지 않았다.

운영 정정: `review.py check --task-key`는 `SAR-AUTO-RECEIVE-001-DEV`다. tester 인박스에 적힌 TESTER 키로 과거 구현 SHA의 패킷을 요구하지 않는다. tester 패킷은 현재 QA의 finish와 완료 SHA가 검사한다. 코디네이터 메시지 msg_aef16591e294와 msg_9ed204c7307a를 반영했다. 고정 SHA와 제품 코드는 바꾸지 않았다.

## 발견 사항

- medium, 미해결. `internal/relay/member.go` 27행. 회원 receipt가 자동 wake가 없으니 수동 receive를 실행하라고 안내한다. `internal/relay/public_messages_integration_test.go` 91행이 그 문장을 기대한다. relay는 push하지 않는다. public-node 클라이언트는 inbox로 자동 pull한 뒤 화면 표시는 `knowslink_text_receive`가 한다. 데이터 유실은 아니다. 안내가 자동 저장과 어긋난다. tester는 제품 코드를 고치지 않는다.
- low, 미해결. `adapters/README.md` 269행. `pkill -f 'plugin.js watch'`는 그 호스트의 모든 watcher를 끝낸다. 연결 폴더마다 watcher가 있으면 한 폴더만 중지되지 않는다.
- low, 미해결. `adapters/README.md` 218행과 `adapters/skills/knowslink/SKILL.md` 38행. 문서는 알림 payload에서 `next`를 빠뜨리거나 `id`, `from`, `pending`만 적는다. 코드는 `event`, `id`, `from`, `pending`, `next`를 보낸다. `next`는 안내 문장이고 원문이 아니다.
- 루트 `README.md` 148행과 adapters README 185행의 자동 wake 문장은 공개 메시지와 MVP-003 시험 절이다. 새 자동수신 절은 `hostNotice: sent_unverified`를 노우 턴으로 쓰지 않는다.

critical과 high는 없다.

## 검증 및 남은 제약

Linux lint와 test는 위 lint.json이다. `make test`의 Go race와 `npm test`가 포함되고 둘 다 exit 0이다.

이 세션의 Windows detached worktree `C:\Users\shin\AppData\Local\Temp\knowslink-autorecv-qa`에서 Node 22.22.2로 `node dist/mcp.test.js` exit 0, `python scripts/package_plugin.py --verify` exit 0, `node dist/inbox.test.js` exit 0을 확인했다. Windows zip sha256은 `666704d3d8da0dd20360ebf63d20581216819ede1b2dc3aeb203f90c0837aa7f`다. node_modules는 같은 lockfile blob의 dev 체크아웃 junction이다. Linux 패키지 로그 `D:/workspace/KnowsLink/.git/fullops-gate/coor-af7627d-linux-package.log`는 새 `npm ci` 188패키지, Node v22.22.2, zip sha256 `f7d09684870514b0cfd368458aeaa99d968493abb92c52850ba6d46edb0c4eb6`, 추출 MCP 검사 exit 0, remote exit 0이다. 줄바꿈 때문에 zip 해시는 다르다. coor가 본 exit 1은 재현되지 않았고, 테스트가 assertion 본문을 삼키므로 coor 로그가 없으면 그 예외 종류는 특정하지 못한다. 제품 회귀로 확정하지 않는다. 운영 설치는 이 QA가 수행하지 않는다.

설치본 SDK `dist/esm/server/index.js`의 `isMessageIgnored`는 session에 logging level이 없으면 false다. 따라서 `sendLoggingMessage`는 `notifications/message`를 보낸다. 호스트가 그 알림을 모델 턴으로 올리는 계약은 SDK에 없다. Context7 문서 조회는 월간 한도로 실패했다.

`D:/workspace/KnowsLink/.git/fullops-gate/coor-autorecv-81c-preflight.json`은 수동 receive 없이 `lastSuccessAt`가 갱신됐고 pending 0이다. 호스트 알림이나 노우 턴의 증거가 아니다.

실행하지 않은 것: 실서버 배포, 실계정과 관계 변경, 그록봇 UI, 실제 gateway token을 쓰는 wake, 추가 과금. D12는 이 때문에 추가 수정하지 않았다. snapshot은 reviewer dispatch가 release되기 전에는 정리하지 않는다.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과(확인/미확인/해당 없음) | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | deliverables README, architecture.md 170행, 인박스 | 확인 | 자동 저장과 노우 턴을 구분한다. |
| 구조와 구현/미완료 상태 | inbox.ts, text.ts, mcp.ts, grok-wake.ts | 확인 | 자동 pull과 선택 wake는 구현됐다. 실제 호스트 표시는 미완이다. |
| 실행·검증 방법과 증거 | lint.json, inbox.test.js, mcp.test.js | 확인 | Linux 필수 명령 exit 0. Windows cloudflared 실패는 환경 한계다. |
| 운영·복구(해당 시) | ops-guide.md 344행 | 확인 | 절차는 있으나 이번 QA가 운영 적용을 하지 않았다. |
| 다음 작업·담당·재개 조건 | member.go 27행, coor 설치 보류 | 확인 | 문구는 designer/coor. 설치와 실제 Bot 수신은 coor 후속이다. |
| 로컬 링크·절 접근/지원 한계 | architecture.md 앵커, README 216행 | 확인 | payload 문서만 코드보다 짧다. |
| snapshot 정리 후 정본 접근 | 이 report와 result.json | 확인 | snapshot이 남아 있다. 정리 전 정본은 tester 체크아웃의 이 디렉터리다. |
