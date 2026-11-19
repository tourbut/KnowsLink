---
title: SAR-PUBLIC-MESSAGES-001-DEV — 일반 text·receipt·gate·자원 보호
status: draft
updated: 2026-10-06
owner: dev
tasks: [SAR-PUBLIC-MESSAGES-001-DEV]
summary: 일반 text·receipt·gate·자원 보호의 분석과 구현·검증을 기록한다
---

# SAR-PUBLIC-MESSAGES-001-DEV — 일반 text·receipt·gate·자원 보호

## 기준과 기술 계획

기준 SHA는 `7efbaa349a5857eb1eac859a955ec3a09c91f800`이다. 준비 HEAD는 `f210b565dd3d0a0a43af9e9a3b1fceb43dcdaf93`이다. fullops-common-0.3.3·FULLOPS·project·문서 규칙·PS04/06/07/08–11·UX06/07·MVP C1–C5를 직접 읽었다.
connections.go는 연결 자격만 발급한다. 시험 text는 두 agent allowlist에 묶여 있다. 기존 AGENTS 기록의 일반 text·receipt·gate·HTTP 동시 처리 미구현을 이번 범위에서 해소한다. frozen relay.v1 필드는 바꾸지 않는다.

1. `knowslink.text.v1` 별도 closed signed wire와 `/v1/text/*`를 추가한다. 기존 Message·lease·current-auth·멱등·sweep를 재사용한다. 명시 송신·1회 관련 답장만 지원한다. 원문은 ACK 또는 만료에 지우고 metadata는 24h 보존한다.
2. 기존 회원 홈에서 요청 ID로 자기 receipt를 조회한다. 선택 agent·요청/답장 ID·전달·처리 상태·TTL·수동 pull·실패 복구를 표시한다. 자유 composer나 장기 목록은 만들지 않는다.
3. queue100·pending gate100·receipt20000·claim4·raw32768을 공통 수락 경계에 적용한다. H/R 포화는 부모 TTL의 명시 실패 종료로 처리한다. HTTP 신규16·정리4는 즉시 거부하는 process 채널과 공유 DB 입장 기록으로 제한한다. 게이트 GET/approve는 신규, 검증된 deny는 정리다.
4. 기존 Node 연결 폴더를 읽는 text CLI와 opt-in MCP 도구를 추가한다. 실제 로컬 Node 프로세스 두 개와 일반 신원 HTTP·Postgres를 왕복 검사한다. 운영 계정·실메일·최종 노우↔다닷은 후속이다.

예상 변경은 제품/자동 검사 12–18파일과 기술 원천/증거다. 새 의존성·migration·frontend·wire registry 변경은 없다. 별도 theme/design lint가 없어 기존 memberStyle·Go template·make lint를 쓴다. 직접 시각 판정은 designer 후속이다.

## 탐색과 API 근거

기존 Jev code/doc find·context JSON은 기준 SHA에서 coordinator가 생성했다. 재실행은 기존 결과 보존 오류(exit1)라 원본을 재사용했다. connections.go 충돌 가능성은 연결 자격과 메시지 권한의 범위 차이였고 제품 전제를 바꾸지 않았다. connect.ts의 미전송 부분은 로컬 원문으로 읽었다.
Context7 Go resolve는 `Monthly quota exceeded`로 조회 불가였다. 공식 https://pkg.go.dev/net/http (go1.27.1)와 https://pkg.go.dev/context 를 확인했다. MaxBytesReader·CrossOriginProtection·request context deadline·TimeoutHandler의 취소 경계를 확인하고 기존 stdlib와 pgx transaction 패턴을 재사용한다. 새 SDK API는 쓰지 않는다.

## 검증과 후속

검증 명령·원본 실패·exit·고정 SHA와 코드/metadata 차이는 아래 증거와 완료 전문에 보존한다. 독립 OPS 보안 리뷰·TESTER QA·designer UX06/07 직접 검수는 coor 후속이다. Free 제한과 운영 보호를 유지한다.

## 구현 결과와 판단

일반 text wire·현재 양쪽 회원/agent/키/관계 세대 확인·ACK 후 원문 삭제·1회 관련 회신·자기 receipt ID 조회를 구현했다. 일반 text에는 claim·gate·업무 효과를 허용하지 않는다. receipt의 queued/leased/delivered는 처리 완료와 분리한다. 관련 회신도 별도 요청 ID와 transport 상태를 가진다.
회원 gate는 저장된 서명·typed body·digest·현재 권한·기한을 다시 확인한다. hint는 이스케이프한 보조 정보다. 누락/만료/철회된 검증 자료는 approve를 열지 않는다. 명시 deny 경로는 정리 슬롯을 사용한다. schedule.query의 기본 deny와 실행 불가 commit은 유지한다.
기존 rate 호출을 공통 HTTP 입장 경계로 모았다. 한 입장 요청은 schema·CSRF·권한 거부를 포함해 공유 DB rate를 한 번 소비한다. process 채널 포화는 DB에 대기열을 만들지 않고 즉시429를 반환한다. 따라서 process 입장 이전의 용량 거부는 DB rate 호출을 하지 않는다. 공유 동시 슬롯 포화는 rate를 소비한다. 신규16/정리4·10s 요청/본문 제한·30s crash 슬롯 만료를 적용한다.
공통 message queue100·pending gate100·receipt20000·claim4·raw32768을 frozen/trial/public 경로에 적용했다. high priority는 상한을 우회하지 않는다. ACK·철회·deny는 별도 budget으로 유지한다. receipt 포화 상태에서 H/R을 만들 수 없으면 부모 TTL에서 실패하고 claim을 회수한다.
Node CLI는 기존0700 연결 폴더와0600 key/credential 파일을 읽는다. MCP는 public-node 명시 opt-in·confirmed=true·수동 receive/receipt만 지원한다. 반환 text는 untrusted data다. 재시도 오류에는 안전한 HTTP 코드·error·retry_at만 노출한다. credential·claim·서명본문은 오류에 넣지 않는다.

## 영향 산출물

D03 architecture/tech-stack, D05 interface, D06 data-model, D07 database, D09 CRUD, D10 module 원천을 갱신했다. D08은 table·SQL·migration·sqlc 계약이 바뀌지 않아 생성하지 않았다. D01/D02/D04 제품 정본·UX mockups는 designer 소유이므로 수정하지 않았다. 기존 memberStyle을 사용했다. 별도 theme/design lint와 Tailwind/shadcn 스택은 없다. make lint는 기존 포맷·정적 검사이며 직접 시각 검수를 대신하지 않는다.
PLANS·project·DEV context·README·adapter README·KnowsLink agent skill을 현재 구현에 맞췄다. 새 package·의존성 버전·lockfile·schema는 없다. 공유 row lock과 기존 State를 재사용했고 새 추상 계층이나 composer는 만들지 않았다.

## 검사 범위와 원본 실패

아래 증거 디렉터리의 로그는 각 명령 stdout/stderr이며 후행 공백만 정규화했다. 종료코드는 각 .exit에 저장했다. 파이프의 끝 명령을 통과 근거로 쓰지 않았다. 첫 lint/test는 실제 실패로 유지한다.

| 검사 | 결과와 수정 근거 |
|---|---|
| 초기 go test ./internal/relay | exit1. 기존 홈 전체의 자동이라는 단어 금지 assertion이 새 수동 안내의 자동 wake 없음 문구에 걸렸다. 관계 영역 assertion으로 범위를 맞췄다. 초기 stdout은 터미널 출력이며 원본 파일이 없어 실패 사실을 이 표에 기록한다. |
| make lint 첫 실행 | exit2. text.ts Prettier 위반. 기존 formatter로 수정했다. |
| make test 첫 실행 | exit2. MCP bundle이 imported connect/text CLI main을 실행했다. relayBase를 core로 이동하고 text.js basename·module URL을 함께 확인했다. |
| make verify-mvp 첫 실행 | exit2. local Node/MCP 왕복은 통과했다. 테스트 서명 helper가 KST를 UTC Z로 표기해 replay 경합이422였다. UTC 직렬화로 고쳤다. |
| make verify-mvp 두 번째 | exit2. revoked command 검사에 text envelope를 사용해 schema422가 먼저 발생했다. 실제 pull/persist/ack command와 send wire로 고쳤다. 제품 권한을 완화하지 않았다. |
| make lint/test/verify-mvp 최종 | 각 exit0. 마지막 구현을 다시 검사했다. 고정 커밋의 FullOps는 아래 완료 증거에서 별도로 연결한다. |
| make verify-grok-plugin | 격리 HOME의 실제 로컬 CLI/plugin discovery 검사다. 실제 외부 계정 추론·노우↔다닷 관측과 구분한다. 최종 exit는 증거 .exit를 따른다. |
| receipt template 자동 검사 | exit0. 소유권·queued/ACK/실패 구분·TTL/수동 안내·XSS·composer 없음 검사다. |
| headless Chrome 390×844 | exit0. home·queued receipt·expired receipt의 scrollWidth=innerWidth=390. 정확한 Go template의 공개 fixture만 사용했다. 이미지와 JSON은 DEV 레이아웃 측정이며 designer 직접 시각 PASS가 아니다. |

Postgres 통합은 일반 코드 로그인과 prepare/회원 확인/complete로 두 agent를 만들었다. 실제 분리 Node 연결 프로세스·두 MCP 프로세스·text CLI를 사용했다. 송신→상대 수동 수신/ACK→관련 회신→첫 송신자 수신/ACK→요청·회신 ID 연결과 receipt를 확인했다. 응답 로그에는 metadata ID만 남긴다. 합성 계정 fixture이며 실메일은 보내지 않았다.
HTTP 통합은 일반 신원·타회원 거부·pending 관계·멱등8개 경합·내용 충돌·generic/text 경계·키 철회·unpair/reaccept 세대·기한·재시작을 검사했다. shared queue99+8개 high send는1개만 수락했다. claim5개 경합은4개만 수락했다. 두 relay 인스턴스의 HTTP 신규16/정리4 경합과 종료 회수를 검사했다. 게이트 검증본문·힌트XSS·누락 서명·포화 approve 거부·deny 유지도 검사했다. 기존 frozen/trial 및 격리 자원 회수는 make verify-mvp에서 함께 확인했다.

## 미검증과 인계

coor가 고정 후보 SHA를 받은 뒤 별도 OPS 세션에서 보안·자원 보호 독립 리뷰를 배정한다. TESTER는 PS08–11 독립 QA를 수행한다. designer는 UX06/07 직접 UI를 검수한다. 이 후속은 DEV 자체 자동 검증으로 대체하지 않는다. 미해결 critical/high와 필수 실패는 main 수락을 차단한다.
운영 배포·실메일·실제 외부 계정·최종 노우↔다닷·실24h 관측은 미검증이다. OPS는 공개 수락 전에 최대 JSONB 크기·CPU·처리량·백업/복원·WAL/backup 잔존을 측정한다. 합성 가입 unset·운영DB 합성owner0·기존 low L-A/L-B·Free 제한 조건을 유지한다. DEV 로컬 통과를 공개 수락으로 표시하지 않는다.

## 증거 위치

원본 실패와 최종 실행은 [QA 증거](../../evaluations/qa-reports/SAR-PUBLIC-MESSAGES-001-DEV/)에 보존했다. .exit는 Python subprocess 또는 원래 exec 반환 종료코드다. 후행 공백만 제거했다. 직접 실행의 최종 lint/test/MVP/plugin은 모두 exit0다. 격리 Compose DB·container·network는 종료 후 회수했다.
고정 코드 커밋은 증거 head.txt에 기록한다. 완료 metadata/archive 최종 HEAD의 FullOps 결과는 레포 밖 고정 파일과 같은 Run worker_done에 연결한다. SHA를 자기 커밋 본문에 순환 기록하지 않는다.

## 고정 후보 게이트와 규모

코드 SHA `ef5c571cb56364b59944197a14830124b033433f`에서 lint.py --from `7efbaa349a5857eb1eac859a955ec3a09c91f800`는 exit0이다. product-lint/product-test 각 exit0이며 ERROR0·WARNING11·실행불가0이다. fullops-code.json/log/exit와 head.txt를 QA 증거에 보존했다. 마지막 metadata/archive 커밋에서도 같은 기준 게이트를 다시 실행한다.
SIZE-001 8개는 PLANS933/core316/connections_test390/http632/member341/member_agents316/integration403/store490이다. 불필요 wrapper를 삭제하고 새 wire/capacity/receipt 파일로 경계를 분리했다. 기존 frozen·권한 handler를 더 분할하는 무관한 변경은 줄였다. 새 통합 파일은6개 실제 process/DB 검사와 최소 helper다. SIZE-002는2081추가줄/400예산이다. 예상 제품12–18파일 대비 실제20파일은 중앙 HTTP budget·cleanup과 CLI bundle 근본 원인 수정 때문이다. 전체75파일은 제품20/문서16/증거39이며 선행 coordinator board/Jev는 DEV 생성분이 아니다. 같은 권한·용량 경로를 나누어 미완료 상태로 릴리스하지 않는다.
SEC-001 1개는 시험용 ClaimToken=active이며 실제 비밀이 아니다. SLOP-004 1개는 plugin 검증의 사용자 PASS 출력이다. DEP0·DESIGN0·새 의존성/lock/version0이다. Go template 스타일은 이 정규식 검사로 직접 시각 수락하지 않는다. 독립 reviewer는 경고와 분할 판단을 검토한다.
