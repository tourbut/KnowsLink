---
title: 현재 시스템 구성도·상세 문서 독립 리뷰
status: review
updated: 2026-10-14
owner: coor
tasks: [SAR-SYSTEM-DIAGRAM-001-DEV, SAR-SYSTEM-DIAGRAM-001-REVIEW]
summary: 고정 문서 SHA의 구성·두 수신 흐름·소스/링크와 실제 Mermaid 렌더 검토
---

# 현재 시스템 구성도·상세 문서 독립 리뷰

## 결론·독립성·범위

고정 head 9f84c9d7b510fbf6d84d0d16450116346327d66a의 문서 결과를 수락한다. base/merge-base는8b3c7906029be2c22138e3416b0962a9aa70c2f5다. 미해결 critical/high와 구체적 수정 요청은0이다. 전체20파일 중16 reviewed/4 skipped다. 생성 탐색/IR와 일회성 worker 검사 전문의 생략 범위·대체 확인을 result.json에 기록했다.

검토자는 실제 Codex 세션01a12eb2-a9c9-7973-a811-fd4b5729f7b7다. 작성자의 실제 provider 세션은01a13829-64b2-7450-b1c7-7158189845a3다. worker가 보고 당시 미확인으로 남긴 ID를 coor가 session_meta의 ID/cwd/source와 초기 user preamble의 task_25715fa3ad5d/ctx_c2a550211b14 일치로 독립 확인했다. preamble 본문이나 capability는 출력하지 않았다. 작성자 보고의 당시 미확인을 허위로 수정하지 않았다.

별도 detached snapshot은 `C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-7f3475f6e1244b23ab9d32cd7b43314a`다. 결과는 coor에만 쓰고 snapshot은 읽기 전용·clean으로 유지했다. OCR v1.13.0 custom Markdown/JSON 규칙, fullops-common-0.3.3 코딩·테스트·보안, project.md, 문서 작성 규칙, 원 인박스와 완료 기준을 적용했다. 같은 기준 파일의 source8f255f9와 후보를 사용했다. Codex 호스트 AI의 내용 판단이며 OCR 자동 의미 판정이 아니다.

## 내용과 인계 확인

README·docs 목차·입문 architecture의 전체 diff와 상세 문서311줄을 읽었다. 전체 구성9노드/10연결, Grok5참여자/12메시지, Dots4참여자/11메시지를 실제 기동·Compose/ingress·Google/Remote OAuth·text/Inbox·Remote MCP/Events·store/SQL/writer 근거와 대조했다. source 기준8f255f9, 마지막 운영관측2026-10-13/2747fe9, 직전26094e9 실계정 시험을 구분한다.

Cloudflare와 cloudflared, 물리 경계와 논리 요청, 별도 migration을 정확히 구분한다. Grok은 store/sync/rename→persist/transport ACK→notice/wake→도구 take이고 Dots는 같은 transaction의 persist/설치 Inbox/transport ACK와 호스트 읽음 ACK를 분리한다. Google 신원과 owner/agent/친구 출처·generation, 서명과 암호화, TTL·웹/로컬/Remote 사본, JSONB singleton·행 lock/CAS·NOTIFY, writer와 과거 backup 복원의 한계를 확인했다. 구현에 없는 서비스나 자동 답장·수신 보장을 추가하지 않았다.

정본 인덱스에서 docs/README.md→system-design.md의 구조/수신/인증/저장/운영 설명과 실제 소스 링크를 찾았다. 설치는 기존 README, 복구는 operations/D07, 상세 protocol과 모듈 책임은 D03/D05/D06 정본으로 연결된다. 운영 재확인·장시간·다음 Events 갱신·실계정 후속은 명시돼 있다. 원 인박스 전문이 날짜 로그에 그대로 있고 현재 인박스가 비었음을 직접 대조했다. snapshot 이후에도 리뷰·렌더 근거와 본문은 정본 체크아웃에 남는다.

## 검증·제약

읽기 전용 snapshot에서 기존 일회성 링크 helper를 재사용해 변경 사용자문서4개의 로컬링크101개·fragment·닫힌 fence와 내부 런타임/실계정 식별자 경계를 독립 확인했다. 오류0/EXIT0이다. [독립 링크 결과](independent-links.json)는 worker의 QA/phase를 포함한115개보다 좁은 별도 범위다. worker의 source20·strict13·구조추출3과 최초 실패/수정·archive 검증도 대조했다. 구조 추출은 렌더 성공으로 취급하지 않는다.

별도 Orca QA 페이지에서 공식 사용법에 따른 Mermaid12.2.0 ESM으로 SVG3개를 실제 생성했다. 한글 라벨이 browser snapshot에 있으며 각 SVG의 content bounds가 viewBox 안이다. 고정 후보와 세 그래프 원문 SHA256이 모두 같은 것을 확인한 뒤 재사용했다. [최종 입력 연결](../SAR-SYSTEM-DIAGRAM-001-RENDER/final-binding.json), [실행 결과](../SAR-SYSTEM-DIAGRAM-001-RENDER/render-result.json), [일회성 재현](../SAR-SYSTEM-DIAGRAM-001-RENDER/render-preview.txt)을 보존한다. 이는 고정 HEAD의 새 실행이나 GitHub 렌더러 검사가 아니다.

스크린샷은 Page.captureScreenshot의 `browser page did not draw a frame`으로 실패했다. 탭 선택·재snapshot 뒤에도 같아 추가 엔진으로 우회하지 않았다. [캡처 한계](../SAR-SYSTEM-DIAGRAM-001-RENDER/capture-unavailable.json)를 남기며 픽셀 시각 검수·모바일/GitHub 표시 PASS를 주장하지 않는다. native Markdown 문서 변경이므로 제품 UI/테마/디자인 lint는 해당 없음이다. 문서의 텍스트 설명과 실제 SVG 구조/기하를 내용 수락 근거로 사용한다. 선택형 HTML/SVG 스타일 질문은 이 Mermaid 문서의 완료 조건으로 삼지 않았다.

고정 후보 lint는 ERROR0/WARNING1/unavailable2/EXIT1이다. Windows make 부재를 명령별 reason에 기록했다. 제품 입력 diff는 exact cd4dae8 대비 adapters/README.md뿐이며 코드·설정·테스트 입력은 같다. 기존 Linux gate의 실제 head·두 명령 passed/EXIT0·원본 hash를 확인해 재사용했다. 현재 문서 HEAD 제품 검사를 PASS로 바꾸지 않았다.

SIZE-001은 기존 PLANS1774→1776 경고다. 과거 기록은 보존했다. worker 기준8f255f9의 SIZE-002 추가404줄에는 준비·연결·metadata·phase가 포함된다. 사용자 상세 문서311줄은 요청 범위200–350 안이며 구성/수신/인증/저장을 함께 읽는 설명이므로 더 분리하지 않았다. DEP-001은 없고 제품 의존성 변경0이다.

완료 메시지 msg_e5dc918a82d5는 subject에 정확한 SHA가 있고 body에는 SHA label이 없었다. 원 메시지를 보존하고 실제 fullops/dev HEAD와 subject를 대조해 integration의 SHA metadata만 보완했다. main 일반 push·조상 검증 뒤 ACK한다. 운영·계정·발송·배포·전체 제품·장시간 테스트는 새로 실행하지 않았다. snapshot은 보존하고 담당 coor가 종료·소유권·증거 보존 조건을 확인한 뒤 정리한다.
