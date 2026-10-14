---
title: 공개 첫페이지 설계 검증과 미실행 경계
status: review
updated: 2026-10-14
owner: designer
tasks: [SAR-LANDING-GUIDE-001]
summary: 제품 사실·설계 추적성·문서 검사와 과거 제품 gate 재사용 및 후속 QA를 구분한다
---

# 설계 검증 결과

현재 대상은 공개 첫페이지 설계 문서다. 제품 코드·테스트·설정·의존성은 변경하지 않았다. 실제 UI 동작·시각 수락은 not_run이며 DEV 구현 뒤 독립 QA와 designer가 검증한다.

## 기준과 사실 대조

착수 HEAD는 `93271bf0d2aa8b6677a1f3d3135dc5dd708cddf8`, lint 기준은 `3230ee566ccf455e58a6fc3b43f5cd4b3a72aba0`다. [design-check.json](design-check.json)에 원 packet/현재 packet의 실제 hash와 producer 상태, 제품 입력 비교와 기존 gate hash를 보존한다. [source-audit.json](source-audit.json)은86경로의 필요한 구간을 기록한다.

| 설계 주장 | 읽은 근거 | 판정 |
|---|---|---|
| 사람 수락 아래 현재·향후 활성 agent 연결 | D02 HP-03/06, README 웹에서 시작하기 | 범위 일치 |
| 웹 직접 작성과 연결 agent 발송 요청은 다름 | D02 HP-09/10, docs/architecture.md 송수신·서명 | 가상 예시의 작성 주체와 발송 승인을 명시 |
| 웹 사용 설치 불필요, 실제 대화에는 agent 필요 | README 웹에서 시작하기, D04 빈 계정 흐름 | 소개/설치 카피에 두 조건을 함께 명시 |
| Grok 로컬, Dots Remote와 서버 키 보관 | README 신규 설치/연결, docs/system-design.md | 상세 실행은 README로 연결, 키 복사 금지 |
| 수신·표시·자동답장은 별개 | README 자동 수신/Events, docs/operations.md 수신 알림 | process/routine/구독/표시 조건과 매 발송 승인 명시 |
| 대기·전달됨·읽음, 180초·24h 사본 경계 | D02 HP-11/12, docs/architecture.md 저장과 보관 | 의미·외부 사본 한계를 유지 |
| 유효세션 / 동작 변경 | member.go startPage 64행 이하 | 기존303 /home을 의도적으로 공개 안내 유지로 변경하는 LG-02 인계 |
| 기존 최소 로그인 방향 변경 | D04 과거 로그인 행, 현재 designer 원 지시서 | 2026-10-14 최신 요청 근거와 원문 보존 |

## 문서·구성 검사

frontend.py --repo . --host codex --model gpt-6-astra --effort high --check는 exit0, configuration_ready, problems0이다. tools=[]이며 인증·browser_qa·mcp_connection은 not_run이다. 주 스킬 frontend-design과 기존 D04 여섯 토큰을 적용했다. 새 설치·테마·디자인 검사기는 없다.

문서 검사 명령은 레포 루트에서 실행한다. 실제 종료코드와 결과는 [document-validation.json](document-validation.json)에 기록한다. 검사 대상은 변경 Markdown의 필수 front matter와 로컬 링크, D02/D04/D11 추적성, 전체 deliverables strict, JSON 파싱, git diff --check다. fragment와 원격 링크의 실제 HTTP 성공은 자동 로컬 링크 검사의 범위 밖이며 README의 해당 제목은 별도로 대조한다.

실행 결과는 변경 Markdown10개 front matter/로컬 링크 오류0, README 설치/업데이트 제목4개 존재, 산출물13개 strict 문제0·경고0·exit0, JSON 파싱과 git diff --check exit0이다. 공식 packet 완료 계약58쌍도 통과했다.

LG-01–08은 D02·D04·DEV/tester 인박스에 연결됐다. 공개 카피에는 이메일·키·실제 메시지·운영 SHA가 없다. D04의 구체 카피/상태/폭·토큰·와이어프레임과 두 차례 자기 비평을 검토했다. 이 검토는 독립 리뷰나 실제 렌더 PASS가 아니다.

## 제품 gate와 실행 불가

기존 수락 DEV `3e5af5e5ef9ca993baa6745aaf5b2b1f157ace46`와 현재 작업 트리의 cmd/internal/adapters/db/deploy/scripts·의존성·빌드·설정 입력 diff0을 직접 확인했다. 과거 [verified-final-lint.json](../SAR-REFACTOR-002-FINAL-MAIN-REVIEW-review/verified-final-lint.json)의 SHA256 `f64c96869ea21f83260458cc9cabce6c4c79bad282e43133dca0a290e59e102c`도 일치했다. 그 기록의 Linux make lint/test는 각각 exit0이다. 새 문서 HEAD의 실행 PASS나 현재 Windows 환경 동일성으로 바꾸지 않는다.

설계 작업에서는 제품 실행 입력을 바꾸지 않았으므로 기존 제품 gate를 위 정확한 범위로 재사용한다. FullOps lint는 커밋 뒤 기준3230ee5로 별도 실행하고 결과·미실행을 완료 보고에 연결한다. 등록 product-lint/product-test가 Windows make 부재로 실행 불가하면 새 PASS로 보고하지 않는다. 필수 검사를 비활성화하거나 실행 환경을 바꾸지 않는다.

## 미실행·후속

- 제품 구현·Go 회귀의 새 실행·실제 브라우저/스크린 리더·Google 왕복·운영 메일·MCP 연결·설치 실행은 not_run이다. 문서 전용 설계 범위이기 때문이다.
- 별도 디자인 lint·다중테마는 프로젝트에 없으며 해당 없음이다. Go 문자열 CSS를 일반 DESIGN 규칙이 검증한다고 주장하지 않는다.
- DEV가 구현 SHA의 핵심 auth/template와 필수 make lint/test를 수행한다. tester는 격리 서버에서 상태표·실제 POST·redirect·정보 비노출·키보드·200%·overflow를 독립 확인한다.
- designer는 같은 SHA의 1440·390·320 hero·설치·가상예시·오류 실제 캡처를 판정한다. 별도 세션 fixed SHA 코드/문서 리뷰와 coor main/origin 통합은 후속이다.
- 최신 OPS의 Grok 현지 업데이트·실계정 재왕복·장시간 알림·다음 Events 갱신·실제 rollback 미검증은 그대로 남는다. 이번 카피가 그 성공을 보장하지 않는다.

생략한 것은 새 상품·가격·검색·다크모드·새 이미지/모션·설치 자동화다. 공개 로그인 진입 변경의 실제 회귀 위험은 DEV/tester의 구현 고정 검사로 확인해야 한다.
