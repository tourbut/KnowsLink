---
title: 전반 리팩터링·코드 정비·불필요 파일 정리·로직 점검
status: draft
updated: 2026-11-05
owner: coor
tasks: [SAR-REFACTOR-003]
summary: 사용자 지시로 coor가 직접 수행한 3차 전반 정비의 기준·방법·반영·보류·검증을 기록한다.
---

# SAR-REFACTOR-003

## 착수

- 사용자 요청(2026-11-04): "전반적인 리팩토링 및 코드 정비 불필요한 파일정리 로직점검을 진행해 다른 워커에게 지시하지말고 너가 직접 수행해".
- 역할 worker 배정 없이 coor 세션이 직접 수행했다. 사용자 지시로 coor 배정 규칙의 예외다. `jev_route` 분류와 역할 인박스는 만들지 않았다.
- 점검·구현·리뷰는 coor 세션 안의 하위 에이전트 workflow로 나눴다. 하위 에이전트 effort는 사용자 지시로 high 이하로 둔다.
- 기준 main은 `952f6eaab00be1888dbe7dfbe7886ab1853094e1`다. 작업 브랜치는 `fullops/coor`다.
- 완료 조건:
  - 확인된 결함·죽은 코드·불필요 파일을 기존 동작 계약 안에서 정리한다.
  - Linux 고정 SHA에서 `make lint`·`make test`·Postgres race integration이 모두 종료 0이다.
  - FullOps lint ERROR 0이다.
  - 구현자와 다른 리뷰어의 고정 SHA 독립 리뷰에서 critical/high/medium이 0이다.

## 방법

1. **기준선.** 시작할 때 이 워크트리는 2,346개 파일이 `eol=lf` 규칙과 달리 CRLF로 체크아웃돼 있었다.
   - 내용 해시가 index blob과 같은지 확인한 뒤 LF로 맞췄다. Git 내용 변경은 없다.
   - Linux 기준선은 운영 서버가 아니라 로컬 Docker 컨테이너(node 22 + Go 1.27.1 + 격리 Postgres)에서 돌렸다.
   - `952f6eaa`는 `make lint`·`make test`·build·plugin verify·integration이 통과했다.
   - integration의 `abnormal_clock` 2건은 Docker VM 시계 보정으로 생긴 간헐 실패이며, 3회 재실행에서 통과했다.
2. **점검.** 정적 분석(staticcheck·deadcode·tsc 엄격 모드)을 먼저 돌린 뒤 영역 8곳을 점검했다.
   - 점검 영역은 Go core, 회원·인증, 화면·i18n, 원격 MCP, TS adapter, 스크립트·배포, 문서·파일, 교차 계약이다.
   - 코드 영역마다 로직 관점과 정리 관점의 점검자를 두고, 보강 점검자가 놓친 항목을 찾았다.
   - 모든 지적은 반박 검증을 거쳤다. 버그성 지적은 검증자 2명이 보았다.
   - 결과는 162건이다. 확인 159건, 이견 1건, 반박 2건이다.
3. **구현.** 영역 7곳은 각자 격리 git worktree에서 반영·시험·커밋했고, coor가 병합했다.
   - 병합 충돌 4곳(`text_events.go`, `store.go`, `member.go`·`public_text_test.go`, `sqlc.yaml`)은 양쪽 변경을 모두 살려 해결했다.
   - 병합 뒤 쓰이지 않게 된 EN 키 6개를 지웠다.
   - 영역을 가로지르는 죽은 handler rate 코드(F026/F059)는 병합 뒤 별도 단계로 지웠다(`51b3eb52`, 141줄 삭제).
4. **리뷰.** 고정 SHA `3140acde`를 영역 5곳과 보안 관점의 리뷰어가 보고, 지적마다 반박 검증을 거쳤다.
   - 확인된 지적은 low 1건이다. 수정은 `97195a76`이며 red/green hook 시험을 붙였다.
   - 그 뒤 변경분(`3140acde..97195a76`)은 별도 리뷰어가 다시 보았다.

## 반영 요약

전체 결과는 120개 파일, 1,107줄 추가, 2,023줄 삭제다. 제품 코드만 보면 97개 파일, 1,030줄 추가, 1,958줄 삭제다. 항목별 내역은 [QA 보고](../../evaluations/qa-reports/SAR-REFACTOR-003/report.md)에 있다.

- **동작 수정(문서화된 계약에 맞춤)**
  - 시험(trial) claim이 서비스 전체 claim 4칸을 만료까지 잡던 문제
  - 컨테이너에 tzdata가 없어 IANA 시간대 `schedule.commit`이 거부되던 문제
  - 웹 대화 입력 4096바이트 한국어가 8 KiB form 상한에 막히던 문제
  - 잘못된 form 본문을 그대로 처리하던 문제
  - SIGTERM 때 SSE stream이 Shutdown을 막던 문제
  - remote `/mcp`·token의 DB 장애를 인증 실패로 알리던 문제
  - webhook 4 KiB 초과 응답을 실패로 보고 재전송하던 문제
  - 대화 보관 안내의 7일을 D021의 14일로 고침
  - 폐지된 '활성 N/5' 카운터(SAR-CONNECT-SIMPLE-001 L-01)
  - 장치 확인 재표시의 소유자 확인 누락
  - adapter 송신 TTL 여유(180→170 s)
  - 엣지 HTML 오류 본문을 키·친구 문제로 분류하던 문제
  - persist/ACK 실패 뒤 도착 알림이 사라지던 문제
  - 손상된 pending.json으로 MCP 서버가 종료되던 문제
  - shell guard의 wrapper 옵션 값 우회
  - `verify.py` regression·public 판정
  - `test_public_config.sh`의 git 의존
  - `schema.py`의 CRLF 해시(D08 갱신)
- **죽은 코드·불필요 파일 제거**
  - 기동 캘린더 정리 `calendar_retire.go`와 legacy 필드, `KNOWSLINK_CALENDAR_KEY`. 파일 스스로 정한 제거 조건인 운영 연결 0을 충족했다.
  - 실행되지 않던 handler 쪽 요청 budget 코드
  - 쓰이지 않는 EN 번역 키
  - 폐지된 owner-only Access 도구: `access_apply.py`, `access_trial_plan.py`, `tunnel/config.yml.tmpl`, `beta.sh render-config·expose·seed·owner-login`
  - `api_rate.go`, 루트 `.npmrc`, `db/queries/.gitkeep`, 일회성 publish 브라우저 fixture
- **문서**
  - README·docs의 운영 현재 상태를 `c6dfc39a`로 고침
  - Grok 설치 고정 SHA를 `c6dfc39a`로 고침
  - 삭제·변경에 맞춰 D03·D04·D05·D06·D08·D09·D10·D11·D12·D13의 해당 줄을 갱신하고, 이력 절에는 제거 메모를 붙였다.

## 보류(사용자·기획 판단 필요, 미반영)

- **F017:** 친구 요청과 취소를 반복하면 서비스 전체 친구 행 상한(20000)이 찰 수 있다. 회원당 평생 요청 상한 N은 새 제품 규칙이다.
- **F001:** 원격 설치가 끝나면 agent 전체를 철회하므로 회원이 따로 추가한 로컬 키도 함께 철회된다. 철회 범위는 인증 정책 판단이다.
- **F003:** 익명 DCR로 원격 클라이언트 상한 1000을 30일간 채울 수 있다. 검증에서 의견이 갈렸다.
- **F055:** AGENTS.md 삭제 금지의 해제 조건이 모호하다. 안전 규칙이라 사용자 결정 사항이다.
- **F057:** 실행되지 않은 리뷰 지시서가 남아 있다. FullOps 기록이라 보존한다.
- **F160:** 보관된 QA에서만 쓰는 브라우저 fixture 3개다. 향후 시각 검수에 재사용할 수 있어 보존한다.
- **F131·F138:** 설치기 수정 2건이다. 게시된 0.1.1은 바꿀 수 없으므로 다음 설치기 버전에서 처리한다.
- **F139·F043:** Postgres digest와 릴리스 버전이 여러 곳에 복사돼 있다. 정리 효과가 작아 보류한다.
- **F147:** 'node-local' 표시 이름은 agent에 보이는 API 값도 바꾸므로 보류한다.
- **F056:** 이전 영상 렌더 v01·v02는 보존 결정이 기록돼 있고 영상 수락 전이라 보류한다.

## 검증

- Windows(이 호스트) 검사:
  - toolchain gofmt, `go vet`(기본·모든 fixture 태그), Go 단위 시험이 통과했다.
  - adapters check·test, `claude plugin test adapters` hook 13건이 통과했다.
  - `package_plugin.py --verify`, shell guard·Get-Registered 시험이 통과했다.
- Linux 컨테이너 검사: `29ce3d47`(병합)과 `3140acde`(리뷰 SHA)에서 install·`make lint`·`make test`·build·plugin verify·migrate·race integration이 모두 종료 0이다. 최종 후보 결과는 QA 보고에 적는다.
- FullOps lint `--from 952f6eaa`: ERROR 0이다.
  - 실행 불가 2는 등록 명령 `make lint`·`make test`다. Windows PATH에 make가 없고 public-config 시험은 Linux 파일 모드가 필요하다. 같은 SHA의 Linux 결과로 대체했다.
  - WARNING은 테스트 증가로 인한 파일 크기, 캘린더 정리와 함께 지운 테스트 2개, 400줄 초과 변경 규모다.
- 운영 배포, 실계정 시험, 게시 파일(`/install/0.1.1/`) 교체는 하지 않았다. adapter 변경은 다음 플러그인 릴리스 때 버전을 올려 배포한다.

## 남은 일

- 사용자 승인: main 병합·push, 그 뒤 운영 relay 배포(D12 절차).
- 보류 항목의 사용자·기획 판단.
- 다음 플러그인·설치기 릴리스 때 버전 상향(0.1.2)과 게시.
