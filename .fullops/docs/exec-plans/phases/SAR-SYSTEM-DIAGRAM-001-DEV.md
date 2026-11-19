---
title: SAR-SYSTEM-DIAGRAM-001-DEV 기술 계획
status: draft
updated: 2026-10-14
owner: dev
tasks: [SAR-SYSTEM-DIAGRAM-001-DEV]
summary: 현행 Mermaid 상세 설계의 범위·소스 판단·검증과 후속을 기록한다
---

# SAR-SYSTEM-DIAGRAM-001-DEV — 기술 계획과 작성 결과

## 기준·담당·범위

- Purpose: documentation. 기준 `8f255f9991787e14b38052029d334b9536160fdb`, fullops-common-0.3.3, test lite, subagent standard.
- 실제 작성자는 이 Dispatch의 Codex 세션이다. terminal `term_1c1c452e-99a8-4b66-a7e3-1a7a64b269de`, task `task_25715fa3ad5d`, dispatch `ctx_c2a550211b14`, run `run_5534969da2d0`다. 확인하지 않은 provider session ID는 만들지 않았다.
- 기존 DEV 사용자 터미널과 과거 자료를 보존했다. 선택 하위 위임은 사용하지 않았다. 독립 고정 SHA 내용 리뷰·실제 렌더 검수는 coor 후속이다.
- 제품 코드·설정·의존성·DB·계정·서비스·배포·외부 메시지·삭제·이동은 변경하지 않았다.

## 짧은 기술 계획

1. 필수 규약과 packet 원문을 읽고 실제 소스의 물리·신뢰 경계를 확인한다.
2. 상세 문서에 전체 구성 flowchart와 Grok·Dots sequence를 작성한다. 관계 모델은 표와 본문으로 설명한다.
3. 입문 architecture와 두 목차를 상세 문서에 연결한다. D03/D12 원천은 수정하지 않는다.
4. 링크·fragment·fence·밀도·소스 의미·제품 입력 동일성을 검사한다. packet outcomes와 완료 전문을 보존한다.

착수 계획은 인박스에 먼저 기록했다. packet instruction identity를 보존하려고 같은 계획을 이번 phase와 [착수 원문](../../evaluations/qa-reports/SAR-SYSTEM-DIAGRAM-001-DEV/technical-plan.txt)에 옮겼다. 요구·기준·attempt는 변경하지 않았다.

## 설계 결정과 실제 결과

[현재 시스템 상세 설계](../../../../docs/system-design.md)는 소스 기준 SHA와 마지막 운영 제품 SHA를 분리한다. 2026-10-13/2747fe9 운영 관측과 직전26094e9 실계정·알림 관측을 구별했다.

전체 구성은 TB의 9개 drawable node·10개 연결이다. Cloudflare와 Tunnel, relay와 DB, migration을 분리했다. Google 신원은 별도 외부 노드다. SMTP와 Grok routine은 표·본문에 둔다. 구성도의 응답 생략·Dots 외부 역할 묶음·물리 경계 의미를 설명했다.

Grok sequence는 5참여자·12message다. 로컬 agent API 송신→pull→검증→파일 sync/rename→persist/transport ACK→notice/wake→도구 읽기를 구분한다. 선택 SSE는 본문 없는 신호다. 별도 MCP/watcher process의 파일 marker 경계도 설명했다.

Dots sequence는 4참여자·11message다. DB NOTIFY·Events callback·OAuth receive·같은 transaction의 서버 Inbox와 transport ACK·호스트 읽음 ACK를 구분한다. 설치 Inbox가 무기한 보관되지 않는 사실을 source sweep으로 확인했다.

본문은 owner wire·Ed25519 agent wire·PKCE 설치 OAuth·사람 친구 출처/generation·멱등·180초 기한·24시간 사본·singleton JSONB·row lock/CAS·LISTEN/NOTIFY·writer guard·운영 격리와 백업 제약을 설명한다. 구성요소 표와 각 절의 실제 상대 소스 링크를 제공했다. 설치 본문과 D03/D12 이력을 복제하지 않았다.

## 스킬·공식 근거

FullOps work, Ponytail full, Caveman, Orca orchestration을 적용했다. diagram-design은 사용자 지정대로 구조·밀도 원칙만 참고했다. native Markdown이므로 HTML/SVG 출력·전역 스타일·온보딩은 적용하지 않았다.

Context7 Mermaid resolve는 `Monthly quota exceeded`로 실패했다. 선택 library ID나 query 성공을 만들지 않았다. 공식 [flowchart](https://mermaid.js.org/syntax/flowchart.html)와 [sequenceDiagram](https://mermaid.js.org/syntax/sequenceDiagram.html)을 2026-10-14에 확인했다. 기본 subgraph·Unicode label·참여자·message·opt·note 구문만 사용했다. 새 OpenAI/벤더 API 계약을 추가하지 않고 기존 소스와 운영 기록을 사용했다.

`deliverables.py --stamp`는 `.fullops-squad` 밖 경로를 허용하지 않는다. 루트 사용자문서의 front matter는 같은 설치 모듈의 `render`로 생성했다. 하네스 phase·QA는 정규 `--stamp`로 생성한다. D01–D13 source 상태와 내용은 그대로다.

## 검증과 인계

이번 명령·결과·한계는 [QA 보고](../../evaluations/qa-reports/SAR-SYSTEM-DIAGRAM-001-DEV/report.md)에 있다. Mermaid 구조 추출은 실제 렌더나 화면 검수 PASS가 아니다. coor가 별도 QA 탭에서 고정 SHA 렌더·한글·잘림과 독립 내용 리뷰를 수행한다.

제품 전체 테스트·새 의존성·새 상설 검사 도구는 추가하지 않았다. stdin 검사 재현 원문은 QA 증거로만 보존한다. exact cd4dae8 Linux gate는 제품 입력 동일성을 확인한 기존 실행 근거이며 이번 문서 HEAD의 신규 제품 실행이 아니다. Windows make 부재는 unavailable로 기록한다.

규모 경고가 발생하면 상세 문서와 이번 검증 기록의 필요 범위를 설명한다. 약300줄 상세 문서는 요청된200–350줄 범위 안에 있다. 과거 PLANS와 증거를 삭제하거나 이동해서 경고를 숨기지 않는다.

완료 보고 전문과 원 지시서는 `work.py finish`로 날짜별 DEV 로그에 보존한다. coor가 고정 SHA 리뷰·시각 검수·main 일반 통합과 push를 담당한다. 이번 worker는 제품 수락이나 운영 재확인을 주장하지 않는다.
