---
title: SAR-SYSTEM-DIAGRAM-001-DEV 문서 QA
status: draft
updated: 2026-10-14
owner: dev
tasks: [SAR-SYSTEM-DIAGRAM-001-DEV]
summary: 링크·소스·Mermaid 구조·제품 게이트 재사용과 검증 한계를 기록한다
---

# SAR-SYSTEM-DIAGRAM-001-DEV — 문서 검증 보고

## 대상과 결과

문서 기준은 `8f255f9991787e14b38052029d334b9536160fdb`다. 착수 HEAD는 지시서·packet 준비 커밋 `8b3c7906029be2c22138e3416b0962a9aa70c2f5`다. 운영 제품은 2026-10-13 관측 `2747fe98653d5c6946f41c219fc256395333d7fe`로 구분했다.

작성 산출물은 [현재 시스템 상세 설계](../../../../../docs/system-design.md)와 README·docs 목차·입문 architecture의 연결이다. 제품 코드·계정·배포는 변경하지 않았다. [기술 계획](../../../exec-plans/phases/SAR-SYSTEM-DIAGRAM-001-DEV.md)에 설계 판단을 기록했다.

실제 작성 세션은 terminal `term_1c1c452e-99a8-4b66-a7e3-1a7a64b269de`, task `task_25715fa3ad5d`, dispatch `ctx_c2a550211b14`다. provider session ID는 별도 확인하지 않았다. 작성자 자기 검토는 독립 리뷰를 대신하지 않는다.

## 소스 의미 대조

| 검토 주제 | 직접 근거와 판정 |
|---|---|
| 구성·기동 | cmd/relay는 설정·pool·Google/선택 SMTP·Cleanup·RemoteEvents를 시작한다. cmd/migrate는 별도 SQL-only Up이다 |
| 물리 경계 | Compose의 internal DB network·loopback relay·선택 Tunnel, 운영 override와 ingress 정확한 경로를 확인했다 |
| 관계 | human.go authority의 self/friend/legacy 출처와 generation, store.current의 owner·키·권한 재검사를 확인했다 |
| 인증·wire | google.go의 state/nonce/PKCE·ID token, connections.go의 키 PoP, public_text.go와 text.ts의 Ed25519, human_messages.go의 owner wire를 확인했다 |
| Grok 저장·ACK | text.receive는 callback store 뒤 persist/ACK한다. Inbox.store의 sync/rename과 ID marker, AutoReceiver와 MCP take/notify를 확인했다 |
| Dots 저장·ACK | remoteTool receive의 persist→설치 Inbox→transport ACK와 별도 knowslink_ack 삭제, sweepRemote의 만료/현재권한 제거를 확인했다 |
| Events·SSE | SSE 빈 신호·LISTEN·heartbeat, Remote callback 검증·lease·HMAC·재시도·구독 TTL을 대조했다 |
| 동시성·보관 | relay.sql FOR UPDATE·epoch CAS, store clock/sweep/lease, Conversation 24h, migration2 guard와 의도적 Down 거부를 확인했다 |
| 운영 한계 | 기존 OPS 현행 절과 보고의2747fe9 관측·직전26094e9 실계정 범위·격리복원·rollback 미실행을 구별했다 |

## 실행 근거

- Mermaid extractor 첫 실행은 QA 폴더가 없어 output 저장에서 실패했다. 문법 실패가 아니다. 폴더 생성 뒤 같은 설치 도구를 재실행했다. 사용자 파일을 삭제·이동하지 않았다.
- [Mermaid IR](mermaid-ir.json)은 설치 diagram-design extractor의 이번 문서 구조 추출 결과다. flowchart 9 drawable node/2 container/10edge, Grok 5participant/12message, Dots 4participant/11message다. 추출은 렌더 PASS가 아니다.
- [정적 검사](static-checks.json)는 링크·fragment·fence·front matter·그림 밀도·방향·선택 조건과 정확한 소스 token/순서를 확인한다. [실행 원문](check-command.txt)은 stdin 검사의 재현 증거다.
- 정적 검사 결과는 로컬링크113·소스검사20·그림3·오류0·EXIT0이다. `deliverables.py --strict`는13개·문제0·경고0·EXIT0이며 `git diff --check`도EXIT0이다.
- [첫 정적 결과](static-checks-first.json)의 EXIT1을 보존했다. 생성 전 증거 링크와 기준 이후의 하네스 board JSON을 제품변경으로 분류한 검사 범위 문제였다. 증거를 생성하고 제품 실행 입력의 판정 범위를 바로잡아 위 재검사를 통과했다. 제품을 수정하지 않았다.
- [gate 재사용](gate-reuse.json)은 exact cd4dae8와 현재 제품 입력의 Git diff 및 기존 Linux 명령 EXIT0 근거다. 제품 경로에서 adapter README만 문서 변경이며 실행 입력 변경은 없다.
- [packet audit](packet-audit.json)과 [이번 outcomes](../../jev/SAR-SYSTEM-DIAGRAM-001-DEV-packet-outcomes.json)는34경로/36쌍과 unknown20건을 이번 확인 범위로 처리한다. 과거192쌍 결과를 복사하지 않았다.
- `git diff --check`, `deliverables.py --strict`, 커밋 후 FullOps lint의 실제 결과는 완료 기록에 고정한다. lint의 Windows product-lint/test unavailable은 신규 PASS로 바꾸지 않는다.

## 검증하지 않은 범위·후속

실제 Mermaid 렌더·한국어 잘림·브라우저 시각 검수는 coor가 고정 SHA에서 수행한다. coor의 2026-10-14 follow-up이 역할을 명시했다. worker는 새 패키지나 서비스를 기동하지 않았다.

새 로그인·실계정 송신·receive/ACK·호스트 설치/업데이트·Events 갱신·재시작·운영 상태·부하·백업/복원·전체 제품 테스트는 실행하지 않았다. 문서 변경에는 그 실행이 필요하지 않다. 기존 제품 게이트는 입력 동일성이 확인된 과거 실행으로만 재사용한다.

소스 token 검사는 의미 검토를 보조하며 전체 동작 증명이 아니다. 상대 링크 검사는 레포 파일/fragment 대상으로 한정한다. 외부 URL의 장기 가용성은 보장하지 않는다. 독립 내용 리뷰·직접 렌더 검수·main 통합/push는 coor 후속이다.
