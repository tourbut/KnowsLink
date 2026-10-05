---
title: 일반 이메일 서비스 제품 결정과 인계
status: draft
updated: 2026-10-05
owner: designer
tasks: [SAR-PUBLIC-SERVICE-001]
summary: 일반 서비스 제품 결정과 OPS 근거 및 역할별 인계와 검증 결과를 기록한다
---

# SAR-PUBLIC-SERVICE-001 — 일반 이메일 서비스 제품 결정과 인계

## 기준·범위·근거

브랜치는 `fullops/designer`, 기준 ref는 `6c0d132d26dac22693f1c1ffb986fc4f61bdbd09`, 착수 준비 HEAD는 `608fe06cf0d67d97a483bd088d9834c34049fb6f`다. Task는 `task_c21892048d9c`, Dispatch는 `ctx_000841343790`, Run은 `run_8ca8bc058ab7`다. coordinator는 `term_1db428fe-3b8f-43e5-89bd-3cadbd6720e9`다. 시작 작업 트리는 깨끗했고 준비 HEAD에 기준 ref가 포함됐다.

`fullops-common-0.3.2`, FULLOPS.md, project.md, orca-agents.md, document-writing.md와 coding-style/testing/security 규칙을 읽었다. fullops-work·fullops-deliverables와 orchestration을 적용했다. OpenAI Docs로 dots/computers-and-apps와 plugin auth/events 공식 문서를 직접 확인했다. 제품 코드·기술 정본·배포·Cloudflare 설정·원천 스냅샷·PLANS/board는 수정하지 않았다. QA나 직접 UI를 실행했다고 주장하지 않는다.

받은 지시서의 Jev documents-find/context 정상 결과를 재사용했다. 필수 D01/D02·백로그·제품 원천·ops-guide·실제 CLI 시험/정리·이전 공개 정책과 자신의 컨텍스트를 대조했다. `/v1/owners`의 합성 가입 주석과 인증 경계를 직접 확인했다. 기존 owner-only Access를 일반 회원 로그인으로 처리하지 않았다.

## 사용자 지시와 새 결정

원래 지시는 일반 이메일로 시험하고 일반 서비스 완성 뒤 OpenAI dots를 연결하는 것이다. 후속 메시지 `msg_e7183c748087`은 “내가 운영버전 테스트할꺼는 grok bot 과 openai dots를 연결해서 테스트하는거야 둘다 내 이메일로 가입할꺼야.”라는 원문을 전달했다. `msg_bae19c774823`은 Grok Bot 이름 “노우”를 확정했다. 최종 운영 E2E는 같은 일반 이메일의 공식 xAI Grok Bot “노우”↔OpenAI dot “다닷”이다.

일반 서비스 기반을 먼저 수락한 뒤 다닷을 연결한다. 최종 사용자 운영 시험 완료에는 두 플랫폼의 일반 온보딩·별도 agentID/키/credential·명시적 관계 수락·노우 송신/다닷 수신·다닷 관련 회신/노우 수신 ID·앱 도구와 연결 상태가 필요하다. 두 다른 이메일·두 인간 계정을 요구하지 않는다. 다른 회원 권한 분리 음성 QA는 독립 fixture로 할 수 있다. 같은 owner 아래 서로 다른 agent의 관계도 첫 수락을 생략하지 않는다. 표시 이름은 내부 고유 agentID가 아니다.

[일반 서비스 D02](../../planning/product-specs/SAR-PUBLIC-SERVICE.md)는 PS-01–14와 한도·거부·보존·복구 기본값을 정한다. 새 기본값은 사용자의 완성 지시에 따른 제품 담당 결정이다. 과거 DEC-03 제안의 승인을 소급하지 않는다. 실일정 disclosure·결제·유료화·임의 업무 외부 발송의 held는 유지한다. 첫 실제 연결 확인 text는 비민감·명시적 송수신·180s TTL·4096 UTF-8 bytes이며 frozen 업무 wire와 분리한다. 실제 일반 회원 자격으로 검증하며 trial credential/allowlist를 재사용하지 않는다.

[UX 문서](../../design-docs/mockups/SAR-PUBLIC-SERVICE-UX.md)는 UX-01–08을 정의한다. 관리 오너 아이디·공유 서버 SSH·서버 파일 배치를 일반 사용자에게 요구하지 않는다. 수락·전달·처리 완료·실패·수동 수신 확인을 구분한다. 채팅 composer·장기 타임라인은 추가하지 않는다.

## OPS 근거와 책임 충돌 처리

OPS 완료 SHA `0313deae0dec9af813b70ee9685e1a6d0a2b84d7`의 `.fullops-squad/docs/exec-plans/phases/SAR-PUBLIC-SERVICE-OPS-READINESS.md`를 `git show`로 읽었다. 최초 `2267a4a` 기록과 정정 SHA `0313deae`의 차이를 읽었다. `msg_119b5ae89755`와 정정 전달 `msg_427f431c06fc`를 대조했다. 관측 시각은 12:11:29Z부터 약 12:14Z까지다. Access JWT/header는 전달 가능하지만 relay가 미사용하며 실서버 도달은 미확인이다. 관리 token은 관측 시 active이며 만료 전 사용 가능하나 현재 쓰기 권한은 미확인이다. 진행 체크아웃에는 merge하지 않았다. 기록 전문과 D12 반영은 coor가 통합한다.

현재 relay는 Access JWT/email을 읽지 않고 합성 Basic/Bearer만 읽는다. root는 관리 이메일 한 개이고 OAuth는 꺼져 있다. server 자원은 유휴 관측이며 단일 jsonb 행의 실제 수용량·처리량은 미측정이다. swap 약 99% 관측과 shared service 보호를 공개 전에 확인한다. 자동 백업·별도 저장 사본이 없고 최신 복원은 미검증이다. D02의 RPO 24h/RTO 4h·일일 백업·별도 저장 사본·암호화·7일 보존은 후속 구현/복구 검사 기준으로 새로 정했다. 현재 충족으로 표시하지 않는다.

OPS 기록은 A/B/C 신원 방식과 agent/MCP 기술 선택을 designer 결정 후보로 적었다. 이 책임 분류는 프로젝트 규약과 충돌한다. 제품은 이메일 소유 확인·누구나 시작 가능한 가입·같은 회원 신원·일반 사용자 비관리자 흐름을 확정한다. IdP·Access/직접 인증·OAuth·API/DB·서버 설정의 선택과 호환성 검증은 DEV/OPS가 맡는다. Cloudflare 계정 가입을 일반 회원의 필수 조건으로 만들지 않는다. root를 단순 everyone/bypass로 바꾸는 방법은 신원 바인딩 근거가 아니며 공개 전 차단한다.

## 실행 가능한 인계와 PLANS 대기 제안

DEV [현재 인박스](../../../handovers/to_dev.md)는 SAR-PUBLIC-IDENTITY-001-DEV ready다. PS-01–04·신원/세션 한도·UX-01–03과 실제 일반 이메일 한 개의 가입/로그아웃/재로그인을 완료한다. 기술 계획·구현·검사·D03과 필요한 기술/사용자 안내 문서를 같은 과제에서 갱신한다. TESTER [현재 인박스](../../../handovers/to_tester.md)는 SAR-PUBLIC-IDENTITY-001-TESTER waiting이며 DEV 고정 후보 뒤 독립 검사한다. coor가 최신 main 포함·실제 기준 SHA·새 dispatch 복귀 정보를 채운다.

OPS는 착수 체크아웃의 빈 인박스와 달리 coor에서 readiness가 진행 중이었다. `msg_71854ecb4417`의 점유 경고를 반영하여 `to_ops.md`를 수정하지 않았다. 이후 finish 완료 회신에서도 coor가 후속 정규 지시서를 작성하도록 대기 인계를 유지했다. 사용 중 인박스를 덮어쓰지 않았다. PLANS/board는 coor 소유이므로 [제품 백로그](../../planning/SAR-MVP-backlog.md)와 아래 제안을 전달한다.

| 대기 키·역할 | 목표·선행 | 완료 조건·재개 조건 |
|---|---|---|
| SAR-PUBLIC-AGENTS-001-DEV, dev | identity 기능 수락 뒤 자기 클라이언트 연결·키·관계 구현 | PS-04–07, 해당 PS-11, UX-04/05와 계정 비활성화/전체 자격 철회. 별도 agent 자격·같은 owner 첫 수락·회전/철회·타 회원 거부. DEV 인박스 finish 뒤 작성한다. |
| SAR-PUBLIC-MESSAGES-001-DEV, dev | agents 기능 수락 뒤 실제 연결 확인 왕복·실패·gate | PS-08–11, UX-06/07. 기술 중간 왕복과 최종 노우↔다닷을 구분한다. 실제 클라이언트 인터페이스 근거가 필요하다. |
| SAR-PUBLIC-SERVICE-OPS, ops | readiness 결과·수락 구현 SHA 뒤 안전 운영·공개 검증 | PS-12/13. 신원/agent 원점 경계·자원 상한·부하 근거·백업/격리 복원·공유서비스 회귀·rollback. root owner 보존·정리 budget·재시작/복원 철회 불변을 검증한다. 일반 signup 우회를 차단한 뒤 공개를 확대한다. |
| SAR-PUBLIC-SERVICE-ACCEPT-001, coor/tester/designer/reviewer | 안정된 운영 후보에서 기능 전체 수락 | PS-01–13, UX-01–07. 독립 fixed-SHA QA·별도 read-only 세션 리뷰·실제 이메일 사람 로그인·OPS 근거. critical/high·필수 실패는 수락 차단이다. |
| SAR-DOTS-DADAT-001, dev/ops/tester/designer | 일반 서비스 수락 뒤 사용자 같은 이메일의 노우↔다닷 | PS-14·UX-08, 실제 두 Bot 도구/연결·별도 자격·명시적 관계·양쪽 ID 왕복. 실제 다닷 계정 화면과 사용자 연결 동의를 확보한다. 자동 wake는 별도 검증/제품 범위다. |

OPS 인계는 실제 적용 전에 필요한 권한·좌석/요금제·IdP·DNS/WAF·메일 발송 또는 제공자 설정·백업 저장소를 확인한다. 현재 403과 미확인을 PASS로 바꾸지 않는다. 사용자가 이미 승인한 범위는 재승인받지 않는다. 새 비용·부여되지 않은 권한·실데이터/업무 효과의 확대만 사용자 입력으로 남긴다.

## 공식 dots 근거의 범위

D02에 2026-10-05 열람한 공식 URL을 기록했다. 계정의 supported plugins와 연결 컴퓨터는 지원 경로 후보이며 다닷 실제 화면/설치 증거가 아니다. 사용자별 원격 MCP 인증은 공식 OAuth 2.1/PKCE·metadata·token 검증 계약을 확인했다. DEV가 신원 연속성과 최소 권한을 함께 설계한다. MCP Events는 dots의 명시적 구독과 callback/webhook을 안내한다. relay optional webhook OFF는 이번에 유지한다. CLI pull·앱 설치만으로 자동 wake를 주장하지 않는다.

## 후속 인계 탐색 근거

제품 문서 체크포인트 `82a92f559bcd38e9388ec69809a3bf0adb240cb0`에서 DEV/TESTER 각각 code-find·documents-find·context를 실행했다. 모두 정상이며 fallback은 없다. 결과는 `docs/evaluations/jev/SAR-PUBLIC-IDENTITY-001-{DEV,TESTER}-{find,documents-find,context}.json` 여섯 파일에 보존했다. source passage가 크거나 민감할 수 있어 미송신된 verify_mvp.py/verify_runtime.py는 제외하지 않고 필수 keep으로 유지했다.

context가 README.md와 project.md의 기존 합성/실벤더 제외 설명을 충돌 후보로 골랐다. TESTER의 deploy/knowslink/verify.py도 현재 공개 probe의 302 기대를 갖는다. 일반 신원 positive를 기존 모두 302 결과만으로 PASS 처리하지 않도록 인계에 명시했다. 현재 제품 요구가 우선이며 관련 기술 안내와 검사 갱신은 DEV/OPS 책임이다. DEV의 adapters/README.md만 omit? 추천이라 필요 시 확인으로 분리했다. 새 검증/기술 작업에서 유용하면 제외 추천을 무시하고 완료 보고에 근거를 남긴다.

## 문서 검증과 한계

검증 로그는 레포 밖 `/tmp/SAR-PUBLIC-SERVICE-001-validation/`에 명령별 종료코드와 함께 보존한다. D01/D02와 일반 변경 문서는 deliverables.py --stamp로 메타데이터를 작성한다. deliverables strict·local link·패치 공백·소유 범위·고정 SHA lint·전문 아카이브와 빈 인박스를 확인한다. 실제 결과는 다음 절에 기록한다.

제품 코드·배포 변경이 없으므로 신규 제품 동작 테스트와 직접 UI·실메일·플랫폼 실제 왕복은 미실행이다. 등록된 product-lint는 고정 SHA에서 실행해 문서 lint와 분리한다. 일반 서비스와 최종 노우↔다닷의 구현·독립 QA·직접 UI·배포 수락은 후속이며 이번 제품 문서 완료는 그 PASS가 아니다.

## 검사 결과와 완료 인계

검사 완료 후 실제 결과·SHA를 기록한다.
