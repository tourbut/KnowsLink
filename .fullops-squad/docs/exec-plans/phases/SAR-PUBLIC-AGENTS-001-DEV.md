---
title: 일반 회원 agent 연결과 키 및 관계의 구현·검증 기록
status: draft
updated: 2026-10-06
owner: dev
tasks: [SAR-PUBLIC-AGENTS-001-DEV]
summary: 일반 회원 agent 연결·키별 자격·관계·한도와 자동 검증 및 후속 QA를 기록한다
---

# SAR-PUBLIC-AGENTS-001-DEV — 일반 회원 연결·키·관계 실행 기록

## 기준과 짧은 기술 계획

준비 HEAD·main·origin/main은 `fe642109a0ba8bc8940d3c688dbacbcebe4a2f8e`다. lint 기준은 `d2f7ba5aeb6644fd2b27fe6ada5b61db5976933b`이며 준비 HEAD의 조상임을 확인했다. Task는 `task_23b84e8ba5e6`, Dispatch는 `ctx_69e4d4d35bc6`다. 적용 기준은 fullops-common-0.3.3·FULLOPS·project·document-writing·PS-04–07/해당 PS-11/개인정보·UX-04–05다. 정본 인박스는 `handovers/to_dev.md`다.

기존 회원 세션·Go template·singleton transaction·Ed25519·pair 세대를 재사용했다. owner는 세션에서만 결정한다. 연결은 발급·PoP 준비·owner 지문 확인·클라이언트의1회 완료로 분리했다. 최근5분 재인증은 agent 생성·grant 발급/확인·키/agent 철회에 요구한다. grant는 owner·agent·client=node-local·등록/회전에 묶인 최대10분 수단이다. 발급만으로 연결하지 않는다.

이전 agent-wide credential은 활성 키와 분리돼 있었다. 공개 회원 연결은 Key.Credential로 키별 bearer 해시를 저장한다. 공개 연결 완료는 legacy Agent.Credential을 제거한다. 선택 철회와 회전 뒤 해당 key credential의 모든 agent API 인증을 거부한다. 기존 synthetic-only owner API의 저장 자격은 회귀 호환으로 유지한다. 이를 일반 회원 인증으로 사용하지 않는다.

새 라이브러리·frontend·migration·frozen wire 변경은 없다. Go와 Node stdlib·기존 Adapter.request만 사용했다. 기존 동일 버전의 crypto·filesystem·HTTP 사용 근거를 재사용했다. Context7 신규 라이브러리 조회는 해당 없음이다. Jev omit 후보인 adapters/README는 지원 클라이언트 설치 안내를 위해 읽고 갱신했다.

## 구현과 한도

- 연결은 waiting → prepared → approved → consumed다. 실패·취소·만료는 새 키 미연결이며 기존 키를 철회하지 않는다. 준비 키의 PoP는 token·owner·agent·client·mode·kid·public에 결속한다. 준비 뒤 공개키를 바꿀 수 없다. 동일 kid를 재사용하지 않는다.
- 키 등록은 동시 활성3개다. 회전 완료는 기존 키 전체를 철회하고 새 key credential을 발급한다. credential과 개인키를 화면·tool 출력에 반환하지 않는다. credential은 클라이언트 완료 응답에 한 번만 전달한다. 개인키는 클라이언트 로컬 생성이다.
- Node22 로컬 CLI는 준비·완료를 분리한다. 새0700폴더와0600파일만 사용하고 기존 디렉터리를 덮어쓰지 않는다. 완료 때 소유권·권한·symlink를 확인한다. HTTPS root/loopback HTTP root만 허용하며 redirect·userinfo·query·path를 거부한다. 완료 응답을 잃으면 새 회전으로 복구한다.
- 전체 활성 agent200·owner5를 집행한다. 미연결 등록 agent도 slot을 소비한다. agent 전체 철회로 slot을 회수한다. inactive owner의 agent는 인가·집계에서 제외한다.
- active pair400·owner20, pending200·송신owner10, pending24h를 집행한다. pending은 active slot을 쓰지 않는다. B-owner만 수락/거절하고 양측이 철회한다. same-owner의 두 agent도 첫 수락을 요구한다. 브라우저 결정은 현재 Generation을 비교해 옛 수락 화면의 재사용을 막는다.
- HTTP 신규 principal40·전체200/60s, cleanup principal20·전체100/60s, 익명IP30/60s를 DB rate 목록에 저장한다. valid grant는 회원 principal이다. invalid grant는 source IP다. key별 agent credential은 stable agent ID로 집계하므로 key/프로세스 교체로 rate를 초기화하지 않는다. 거부 요청도 기록한다.
- grant 상태는 token 해시로만 보관하며 만료 뒤24h에 삭제한다. 완료한 grant 만료는 활성 키를 철회하지 않는다. 이메일은 agent ID·grant·pair·contacts·로그에 쓰지 않는다. owner 비활성·agent 철회는 기존 current/lease/claim 경계를 차단한다.

## UI와 검사 범위

홈의 새 agent/키·지문·상태·지원 client·재확인·관계 초대/수락/거절/철회를 구현했다. 연결 화면은 대상·권한·기한·지문·확인/취소·대기/완료/만료를 표시한다. 상태를 색만으로 구분하지 않는다. 입력에 label을 쓰고 긴 ID는 readonly 입력으로 복사할 수 있다. 기존 memberStyle을 재사용했다.
Go template CSS에는 별도 theme/Tailwind/shadcn·디자인 전용 lint가 없다. 해당 없음이며 make lint와 FullOps lint를 적용한다. DESIGN 규칙이 Go 문자열 CSS를 검사하지 않는 한계는 유지한다. 직접 UX-04–05 시각 검수는 designer 후속이다.

## 자동 검증 증거

모든 명령은 레포 루트에서 실행한다. 명령의 원래 종료코드를 저장했고 파이프로 가리지 않았다. 최종 로그 위치는 [QA 증거](../../evaluations/qa-reports/SAR-PUBLIC-AGENTS-001-DEV/)다.

| 명령 | 결과·대상 |
|---|---|
| make lint | exit0. Go format/vet/module·TypeScript format/ESLint/type·Compose·scripts |
| make test | exit0. Go race·연결/한도 단위·MCP/시험/로컬 CLI 검증 |
| make verify-mvp | exit0. 실제 격리 Postgres·migration·Go integration race·TS synthetic/seed/trial 왕복 |
| git diff --check | exit0 |
| strict deliverables | 완료 게이트에서 실행·결과 기록 |
| FullOps lint --from d2f7ba5… | 코드 커밋 뒤 실행·결과 기록 |

TestPublicAgentHTTP는 타 회원 agent/grant 바꿔치기·agent로 owner권한 우회·5분 재인증·cross-origin·완료 전 권한 부재·동시 grant 완료1회·회전·선택 철회·옛 credential send/pull/persist/ACK/claim/authorize/gate-consume/result 거부·동일 owner 명시 수락·동시 수락·옛 화면/새 세대·pending만료·동시 owner agent cap·재시작·포화 중 정리를 검증했다.
TestConnectionApprovalAndKeyCredentials·TestConnectionFailureExpiryAndCancellation·TestAgentAndPairCapacity는 key3·agent200/5·pair400/20·pending200/10/24h의 경계/초과·PoP·token/client/owner 결속·철회와 직렬화 복원을 검사했다. connect.test.ts는 실제 Node crypto·HTTP fixture·0700/0600·덮어쓰기 금지·URL·1회 완료를 검사했다. 기존 신원·gate·업무·시험 회귀도 같은 검증에서 실행했다.

`make verify-mvp`의 relay-stop → Go → relay-up 순서를 유지했다. [기존 invalid_lease 원인](SAR-PUBLIC-IDENTITY-001-DEV-TRIAL-DIAG.md)은 외부 relay sweep writer였다. 이번에 같은 fixture DB writer를 격리했다. 기존 실패를 삭제하거나 재실행 PASS로 대체하지 않았다.
SQL·migration·sqlc·table 변경이 없어 make generate/schema는 실행하지 않았다. D08은 기존 실제 schema 정본을 유지했다. bundle MCP 동작을 바꾸지 않아 외부 Grok CLI plugin 설치 검사는 반복하지 않았다.

## 미검증·범위·후속

실메일·공개·운영 배포·외부 계정 연결·과금은 실행하지 않았다. 지원 구현은 Node 로컬 CLI다. Grok Bot/다닷 실제 설치·앱 권한·OAuth·실메시지 왕복 성공을 주장하지 않는다. UX-06·일반 text·queue/gate/receipt 용량·HTTP 동시 처리16/4·claim4는 MESSAGES/공개 수락의 후속 범위다. 계정 비활성화 UI·30일 계정 매핑 삭제는 기존 후속 상태다. 철회 metadata 최소24h와 OPS의 실제 상태 보존량/복원 수락은 유지한다.

DEV 완료는 최종 수락이 아니다. coor는 최종 고정 SHA에서 별도 OPS 보안 delta 리뷰·TESTER 교차 계정 QA·designer 직접 UX-04–05 검수를 배정한다. 미해결 critical/high와 필수 검사 실패는 병합을 차단한다. 독립 리뷰나 사람 시각 검수를 DEV 자동 검사로 대체하지 않는다.

## QA/UI 인계

1. 고정 후보에서 make verify-mvp로 교차 계정·회전·철회·관계·한도 HTTP 증거를 재현한다. synthetic database 외 운영 DB를 사용하지 않는다.
2. README의 로컬 SMTP sink·host relay 절차로 localhost의 일반 회원 화면을 연다. 코드 확인 전 화면·로그에 인증값을 남기지 않는다.
3. 자기 홈에서 Node client를 선택한다. 자기 Node 환경에서 README의 prepare를 실행하고 브라우저 지문을 비교한다. confirm 뒤 complete를 실행한다. grant 상태와 키 활성 여부를 분리해서 확인한다.
4. 캡처는 미연결 홈·공개키 지문 확인·완료/취소/만료·활성/철회 키·수신 초대/수락/철회 관계다. 이메일·코드·grant token·credential·개인키를 캡처 전에 마스킹한다. 긴 ID의 읽기/복사·키보드 이동·label·오류 문구를 확인한다.
5. TESTER는 교차 계정1쌍과 same-owner2 agents를 검사한다. 실제 이메일/플랫폼 검증은 fixture PASS와 분리한다. 시간 변화가 정지 캡처로 판정 불가능할 때만 영상을 만든다.

## 변경 규모와 산출물

보안 경계를 연결 발급·사용·철회·공유 HTTP·회원 UI·클라이언트·회귀까지 하나의 후보에서 검사해야 하므로 한 과제로 유지했다. SIZE 경고를 허용하되 테스트·원천 문서·독립 검수는 생략하지 않았다. 최종 수치는 완료 보고에 기록한다.
D03/D05/D06/D07/D09/D10·README·adapter 안내·project·DEV context를 갱신했다. D08은 schema 변경 없음으로 유지했다. 제품 D02·UX 정본과 board는 변경하지 않았다.
