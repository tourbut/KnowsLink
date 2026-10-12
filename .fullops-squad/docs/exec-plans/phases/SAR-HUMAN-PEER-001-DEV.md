---
title: 사람 친구·이름·연속 대화 기술 구현
status: review
updated: 2026-10-12
owner: dev
tasks: [SAR-HUMAN-PEER-001-DEV]
summary: 기존 singleton JSON과 전송 인가를 유지하며 사람 친구와 웹 대화를 구현하고 검증한다.
---

# 기술 계획

준비 HEAD는 `59d419d0b7c4fca46cabafc04b5a47c6233c944a`다. 공통 기준은 fullops-common-0.3.3이며 제품 정본은 HP-01–16과 현행 D04다. startup packet을 같은 HEAD에서 원 옵션으로 재생성했다. startup check는 예외 없이 통과했다. partial 호출자는 일반 검색으로 보완한다.

1. 기존 Member/Agent에 이름을 추가한다. 신원·키·credential은 그대로 둔다. 기존 회원의 자기 agent 범위 확대는 별도 동의를 저장한다.
2. 사람 친구 상태와 세대를 별도로 저장한다. legacy pair는 정확한 두 endpoint만 유지한다. 공유 인가 함수가 친구·자기 동의·legacy의 권한 출처와 세대를 비교한다. 해제·차단은 legacy보다 우선한다.
3. 웹 발송은 세션 owner 작성의 별도 wire다. 개인키를 수집하거나 agent 서명을 주장하지 않는다. agent signed text는 기존 검증을 유지한다. 각 새 발송은 독립 TTL을 가진다.
4. 새 text의 당시 참여자·표시 이름·본문을 기존 Message 안에 저장한다. 참여 owner의 유효 세션만 24h 이내 과거 복사본을 읽는다. 전송·claim 권한과 분리한다. 기존 원문은 소급 복구하지 않는다.
5. JSON 모델 버전과 DB writer 차단 trigger로 옛 writer의 필드 유실을 막는다. migration은 기존 pending을 종료하고 ID·키·receipt·gate·원문 보존 경계를 유지한다. OPS는 새 writer를 확인한 뒤 이행한다.
6. D04 공용 토큰과 서버 폼으로 전체 UI를 구현한다. 자식 하나는 templates/CSS만, 다른 하나는 격리 QA 환경만 소유한다. 부모가 backend·adapter·최종 검사를 통합한다.

영향 범위는 store/current/ingest/lease/ACK/claim/gate/result, HTTP contacts/keys, SSE/remote, identity, member handlers, adapter receive, SQL migration, fixtures/tests와 D03/D05–D11/README다. pair 참조와 current 호출자를 수동 검색한다. legacy·친구·전송이 같은 인가를 공유하므로 기능과 권한 이행을 분리 배포하지 않는다. SIZE 경고는 전체 P0 기능의 결속으로 설명하며 승인 경계를 낮추지 않는다.

UI 실행은 codex / gpt-6.1-sol / high, frontend-design이다. frontend check는 configuration_ready/exit 0이며 브라우저 QA 성공을 뜻하지 않는다. D04의 두 자기 비평과 여섯 토큰을 그대로 적용한다. DESIGN 자동 검사는 Go 문자열 CSS를 판정하지 못하므로 실제 320/390/1440 화면을 확인한다.

Context7 resolve-library-id는 quota 초과로 실패했다. 공식 [html/template](https://pkg.go.dev/html/template)와 [net/http CrossOriginProtection](https://pkg.go.dev/net/http#CrossOriginProtection) 문서로 출력 이스케이프와 기존 Origin 방어를 확인했다. 새 의존성은 없다.

검증은 합성 A/B/C의 실제 로그인·DB·HTTP happy path, 미수락/제삼자/위조/CSRF/세대/철회/4096bytes/180s/24h/legacy 경계, 기존 signed transport·gate 회귀를 포함한다. 최종 커밋 SHA에서 등록 make lint/test와 FullOps lint를 실행한다. 실계정·실메일·외부 발송·운영 배포는 실행하지 않는다. 임시 파일은 보존한다.

## 구현 결과와 추적

| 요구사항 | 구현·검증 경로 |
|---|---|
| HP-01 이름·신원 | Member/Agent.Name, validName, profile/agent-name와 create/device/remote consent. trim 후1–40 Unicode, Control/Bidi_Control 거부, 중복·40개 astral·ZWJ 이름 허용. Google issuer/subject·owner·키 연속성 유지 |
| HP-02–07 사람 친구 | 무작위 코드/회전, pending 방향·24h·중복/반대 제출, 수락/거절/취소/해제/양측 차단/해제, 출처+세대, 현재/향후 agent 동적 목록, 기존 자기 범위 명시 동의 |
| HP-08–12 대화 | 방향 없는 agent쌍, 고정 발신/수신, 실제 owner/agent 구분, 웹 별도 wire, 새 root별 독립180s,24h 당시 사본, 멱등 충돌·응답 유실 조회·본문 회복, 실제 receipt 상태 |
| HP-13–15 철회·이행 | 공유 authority/current를 keys·send·lease·persist/ACK·claim·gate/result·SSE·remote에서 사용. 해제/차단은 legacy보다 우선, 새 수락은 옛 권한 복원 없음. 구버전 writer trigger, unknown JSON 버전 거부 |
| HP-16 후보 완료 | 실제 합성 OTP/API/DB와 브라우저 캡처, 등록 lint/test를 인계. 제품 최종 수락은 coordinator의 독립 tester/designer/reviewer 후속 |

기존 singleton·키/credential·서명 wire·gate CSRF·Origin/CSP·rate·안전 정리 예산을 재사용했다. 새 의존성은 없다. package.json 변경은 owner wire 회귀를 기존 test 명령에 추가한 것이다. 업무 M/H/R이나 이미 삭제된 text를 대화로 복구하지 않는다. 원격 receive는 기존 messages 배열을 보존하고 별도 authorship map/untrusted=true를 추가했다.

새 product 화면은 human_http/human_pages에 모으고 기존 legacy·직접 연결·receipt·gate를 유지했다. D04 두 자기 비평을 반영해 이름/대화가 첫 화면이며, 기술·키·계정·안전 설명은 관련 동작/상세에 놓았다. 외부 자원·JS·장식 카드·그라디언트·가짜 읽음/온라인은 추가하지 않았다. desktop 세 열·mobile 문서 흐름, 긴 이름·중복 구분·빈 선택기·오류 회복을 검증했다.

D03/D05–D11·README·adapter README/설치 skill·project를 갱신했다. 실제 SQL00002와 모든 migration을 반영하는 D08 생성기를 추가했다. D12/D13은 기술 이행·재개 조건만 연결하고 현재 운영 기록을 보존했다. designer 제품 정본과 coor PLANS/board는 수정하지 않았다. 출처/인가/UI/보존이 하나의 기능이므로 SIZE 경고를 이유로 안전 경계를 분리하지 않았다.

## 실행 증거와 실패 기록

첫 구현 commit은 `4eaefa6e0a5cf623f0d8d8080d918086f86030ec`, 다음 통합은 `f567abff18b8f6ce77fc79e6601ba7163ae7f9ff`, 회귀 보완은 `7f83fd6f3ae9945e292a86212c0ab566ea282b7a`다. 최종40SHA와 같은 HEAD의 등록 검사 결과는 worker_done과 외부 final-lint.json에 고정한다. 자기 참조를 위해 문서 속 SHA를 과거 실행 SHA로 바꾸지 않는다.

실제 런타임·실행기 버전·고유 clone/DB·SMTP·loopback·공개 API fixture·PNG hash/크기·명령 exit와 보존 경로는 [runtime 보고서](../../evaluations/qa-reports/SAR-HUMAN-PEER-001-DEV-test/runtime/report.md)에 있다. Go1.27.1/Node22.22.2/Postgres17.11과 실제 SMTP sink(.invalid)만 사용했다. 키는 client private folder에 생성했고 로그인 세션/코드/키/원문을 로그로 전달하지 않았다.

- 전체 cmd/internal unit, tagged integration/browser 컴파일, human 이름/친구/대화/철회/gate/remote 경계, adapter owner wire self-check, D08 source/hash/trigger self-check는 통과했다.
- `4eaefa6`의 TestHumanPostgres는 실제 A/B/C 로그인·명명·connect proof·친구·6회 왕복·signed 발송·CSRF/C 거부·해제·이력·DB downgrade 차단을 통과했다. 최신 fixture는 A2/B2/C1과 실제 자기 발송/ACK, connection-private temp SQL legacy JSON 변환 검사를 추가했다. 이 temp fixture는 실제 public schema/guard를 재사용하며 다른 시험 singleton을 바꾸지 않는다.
- `f567abf` 전체 tagged SQL은 exit1이었다. 세 실패는 legacy 이미 active인 정확한 쌍의 중복 accept303와 localized gate 상태 기대값이었다. 제품 구현과 상태 불변을 대조하고 기존 승인·거절 비활성·원문 검증을 유지한 채 assertions를 고쳤다. 이 실패 로그를 보존하며 PASS로 재분류하지 않는다.
- 실제5화면×3폭·21개 viewport PNG는 `4eaefa6`의 인증/API/DB fixture다. 320×640/390×844 DPR1, CSS1440×900 DPR0.75이다. DOM overflow0/44px 미만0. 부모도 desktop 홈·320 agent·390 대화를 직접 열어 확인했다. 이후 delta는 이름 UTF-16 ceiling80, 상태 unavailable 번역, remote authorship metadata, 이름 검증/시험/문서다. 캡처를 최종 SHA 실행으로 표시하지 않는다.
- 최초 full PNG는 Orca native surface의 반복 타일로 증거 부적합이었다. 폐기 삭제 없이 보존하고 지원 viewport-section 캡처로 교정했다. 실제1440px/DPR1 이미지라고 주장하지 않는다.
- 키보드는 document.hasFocus=false/Tab BODY/Enter 미작동으로 blocked_engine_focus이며 PASS가 아니다. 시각·키보드·200% 확대·실 mobile 가상 키보드·스크린리더 최종 수락은 독립 designer/tester에 인계한다.

## 탐색·협업·보존과 후속

packet의63개 비선택 path/category와 unknown18/remaining25 수동 검토를 [outcomes](../../evaluations/jev/SAR-HUMAN-PEER-001-DEV-packet-outcomes.json)에 기록했다. authority/current/Pair 및 모든 전송·gate·SSE·remote 호출자를 rg로 확인했다. LSP/동적 전체 그래프나 운영 DB 전체 분석을 완료했다고 주장하지 않는다.

하위2개만 사용했다. UI worker는 template/CSS·기존 integration fixture를, runtime worker는 자기 외부 helper/합성 환경/증거만 소유했다. 부모가 backend/adapter/SQL와 통합·실패 수정·최종 게이트를 담당했다. 하위 자체 검토를 독립 리뷰로 표시하지 않는다. nested Run 이후 parent Run check의 consumer_fenced는 coordinator가 실제 Dispatch 무교체를 확인했고 child Run mailbox check로 복구했다. task/dispatch 권한을 우회하거나 임의 재배정하지 않았다.

임시 검증 파일은 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/humanpeer001-runtime-ctx9ded`와 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/human-peer-001-ui`에 보존했다. 원격 고유 `/tmp/knowslink-humanpeer001-*`, 내부 Docker network/DB/SMTP/relay, QA client folder도 runtime 보고서대로 보존했다. 삭제·이동 정리·운영 배포·실메일/Google·외부 상대·비용 발송은 실행하지 않았다.

coor는 고정 최종 후보의 독립 QA·직접 시각·별도 세션 리뷰와 main 통합을 조정한다. OPS는 수락 뒤 D07의 쓰기 중단·private backup/복원 리허설·migration·새 writer·철회/legacy 검증을 적용한다. 과거 backup 전체 복원은 trigger만으로 안전을 보장하지 않으므로 권한 불명확 시 전달을 차단한다. 상품 quota·dark theme·첨부·검색/내보내기·장기 보관·자동 답장·업무 효과는 추가하지 않았다.

최신 전체 SQL 회귀는 정확한 7f83fd6f3ae9945e292a86212c0ab566ea282b7a에서 go test -race -tags=integration ./internal/relay -count=1, 118.664초, native exit0으로 통과했다. 별도 DB 생성·migration도 exit0이며 필터 없이 실제 HumanPostgres/Legacy를 포함한다. 원본 실패와 재실행 로그는 runtime 보고서의 고정 경로에 보존했다.
