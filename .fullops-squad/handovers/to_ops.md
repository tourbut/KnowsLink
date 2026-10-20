---
title: SAR-RELEASE-002-OPS — main 제품 7b4f9bc를 운영 relay에 배포하고 백업·보존·health·public smoke·/lang 경로를 확인한다
status: draft
updated: 2026-10-20
owner: ops
tasks: [SAR-RELEASE-002-OPS]
summary: main 제품 7b4f9bc를 운영 relay에 배포하고 백업·보존·health·public smoke·/lang 경로를 확인한다
attempt: 53bd68b3e3534d06ab50bc09e6376a8c
base: 1dd6800d5d4238e2060f99e98eaab437ea06682c
subagent_level: standard
test_level: lite
---

# SAR-RELEASE-002-OPS — main 제품 7b4f9bc를 운영 relay에 배포하고 백업·보존·health·public smoke·/lang 경로를 확인한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 work.py finish로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-10-20
- From / To: coor → ops
- 상태: ready
- 테스트 레벨 / 하위 에이전트 레벨: lite / standard
- 승인된 범위: 사용자 결정 D006 '배포 진행'(2026-10-20)과 ingress 질문 '추가·재시작 승인'. 기존 link.knowslog.com 서비스의 KnowsLink checkout·relay/migrate를 배포 SHA로 전환한다. 공개 Tunnel ingress에 `/lang/(ko|en)` 경로만 추가한다(repo `deploy/knowslink/tunnel/public-ingress.yml`과 live config 함께, live 원본 백업), cloudflared 재시작 1회. 과제 고유 private backup·격리 restore fixture·작업 디렉터리를 만든다.
- 추가 승인이 필요한 행위: 그 밖의 DNS·Cloudflare Access·ingress 변경, 삭제·prune·volume 삭제, 실계정 로그인·발송, pull/ACK/token 소비, 운영 DB를 과거 dump로 덮는 복구, migration 실행이 필요한 경우. 필요하면 `ask`로 멈춘다.
- 담당 repo id / 워크트리 절대경로 / 브랜치: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/ops / fullops/ops
- 병합 책임자 / 기본 브랜치: coor / main·origin/main
- 복귀 repo id / 워크트리 / 터미널 핸들 / run id / task id / dispatch id: 0b08ec4c-9e3a-4613-8197-5a835e545335 / C:/Users/shin/orca/workspaces/KnowsLink/coor / term_3a63fc7d-a5d4-4206-85db-3b2ef5622a0c / run_76a9511dbb0d / (worker-start 결과) / (worker-start 결과)

## 현재 상황과 확인 근거

- 마지막 운영 배포는 SAR-RELEASE-001-OPS의 1dd6800(2026-10-17)이다. fresh baseline으로 다시 확인한다.
- 배포 대상은 main의 제품 7b4f9bc이다(그 뒤 main 커밋은 `.fullops-squad/` 기록뿐). 1dd6800..7b4f9bc 제품 변경: `internal/`·`adapters/`·`cmd/` 46파일 +4285/−263, `db/` 변경 없음, `deploy/`는 직전 배포의 ingress 반영 커밋 1개뿐.
- 범위: 잔여 UI #2–#8 다섯 단위(조회 복구 U6, 이름·친구 코드·친구 요청 U2/U4, 승인 요약 U5, 연결 경로·재개 U8, 설정·한국어/English U3)와 후속 low. 회원 상태 JSON에 omitempty 필드가 추가됐다(`Friend.CodeTail`·`Held`·`HeldFrom`·`HeldTail`, `Member.Locale`). DB migration은 없다. 새 경로: `POST /home/language`(/home 아래), `GET /lang/{ko|en}`(공개, ingress 추가 필요).
- 각 단위는 고정 SHA 독립 리뷰·tester QA·designer 시각 검수를 거쳐 main에 통합됐다(PLANS '잔여 작업 일괄 수행' 절).

## 적용 기준과 예외

공통 기준 fullops-common-0.3.3, AGENTS.md의 Windows 파일 작업 안전, project.md, D07/D12/D13. 기준 main(이 지시서 커밋). 예외: 사용자 승인 범위의 ingress `/lang/(ko|en)` 추가.

## 먼저 읽을 문서

- `.fullops-squad/docs/evaluations/qa-reports/SAR-RELEASE-001-OPS/report.md` — 직전 배포 절차·ingress 변경·교훈. 같은 폴더의 명령·smoke 기록.
- `.fullops-squad/handovers/logs/2026-10-17_to_ops.md`의 SAR-RELEASE-001-OPS 절.
- `.fullops-squad/docs/operations/ops-guide.md`, `transition.md`, `.fullops-squad/contexts/ops.md`, README.md 운영 절, docs/operations.md.
- `deploy/knowslink/compose.ops.yaml`, `deploy/knowslink/verify.py`, `deploy/knowslink/tunnel/public-ingress.yml`, `deploy/knowslink/beta.sh`(삭제 포함 절차는 데이터로만 읽고 실행하지 않는다).
- `.fullops-squad/docs/design-docs/SAR-UX-ISSUES-002.md` U3(공개 언어 링크)·U8, `.fullops-squad/docs/design-docs/interface-design.md`·`data-model.md`(새 필드·경로).

직전 private helper는 `C:/Users/shin/AppData/Local/KnowsLinkDevTools/release001-ops/`에 있다. known_hosts/RejectPolicy·자격 메모리 사용·비밀 비출력 계약을 유지한다.

## 해야 할 일과 파일 소유권

- [ ] fresh 운영 SHA·image ID·healthy·요청/DB/공유 서비스 baseline을 읽기 전용으로 확인한다.
- [ ] 배포 SHA bundle/hash·서버 clean checkout·source ancestor·runtime/DB/배포 입력 diff를 확인한다. migration이 있으면 멈추고 보고한다.
- [ ] 새 build 전에 현재 실행 이미지에 과제 고유 rollback 태그를 붙이고 inspect한다. 기존 tag를 덮지 않는다.
- [ ] 새 이미지를 먼저 build하고, 전환 직전 단일 writer stop·새 private dump·네트워크 없는 고유 격리 복원과 digest 검증을 한다. 기존 backup/volume은 보존한다.
- [ ] 새 이미지로 relay/migrate만 재개하고 health·restart·로그·schema·상태 digest·기존 회원/agent/키/관계/메시지·공유 컨테이너 보존을 확인한다. 기존 상태 JSON이 새 omitempty 필드 없이 읽히는지 확인한다. 실패하면 rollback 태그로 복구하고 원 실패를 남긴다.
- [ ] ingress에 `/lang/(ko|en)`만 추가(repo+live, live 원본 백업), cloudflared 재시작, `/lang/ko`·`/lang/en` 공개 응답(쿠키 설정·복귀)과 이웃 경로 404 확인.
- [ ] 기존 public smoke와 바뀐 공개 화면(첫페이지 두 언어·로그인 provider·무인증 회원 경계·보안 헤더·`Content-Language`)을 실제 HTTPS에서 확인한다. 실계정 토큰·로그인은 쓰지 않는다.
- [ ] D12/D13·README 현재 운영 기준·docs/operations·own QA/context를 갱신한다(docs/README.md의 오래된 운영 SHA 포인터 포함). 제품 source·설계·다른 역할 문서는 수정하지 않는다.
- [ ] 실행 명령 자신의 종료코드, UTC/KST 전환 시각·중단 시간·새 백업 mode/hash·이미지·보존 assertions·원 실패·미실행·복구 한계를 남긴다.
- [ ] work.py finish·고정 결과 SHA·worker_done으로 회신한다.

소유권: `.fullops-squad/docs/evaluations/qa-reports/SAR-RELEASE-002-OPS/`, `.fullops-squad/contexts/ops.md`, 이 인박스, D12/D13 원천, README 운영 절, docs/operations.md, docs/README.md 운영 포인터, `.fullops-squad/docs/operations/`, `deploy/knowslink/tunnel/public-ingress.yml`(이번 승인 범위의 한 줄). runtime 제품 변경이 필요하면 멈추고 보고한다(DEV 인계). main 병합·push·현황판은 coor 담당이다.

## 완료 기준과 검증

- 운영이 7b4f9bc 이미지로 healthy이고 restart 0, 기존 데이터·Events·공유 서비스 보존 assertion이 통과한다.
- 새 private 백업과 격리 복원 검증, rollback 태그 inspect가 남는다.
- public smoke 전부, 바뀐 공개 화면, `/lang/` 경로와 이웃 404 검사가 통과한다.
- 제품 코드 gate는 각 단위의 수락 증거(main 통합 기록)를 재사용하고 이번 운영 검사와 구분해 적는다. lint.py는 결과 문서 HEAD에서 ERROR 0이어야 한다(Windows make 불가는 기존과 같이 사유 기록).
- 실계정 로그인·메시지 왕복은 범위 밖이며 미실행으로 적는다(후속 실설치·실계정 검증 과제에서 다룬다).

### UI 작업의 디자인 기준과 검증

해당 없음. 운영 배포이며 UI 코드를 바꾸지 않는다. 공개 페이지는 배포 후 HTTPS에서 내용·헤더만 확인한다.

## 갱신할 산출물

D12, D13(필요 시 D10).

## 기대 산출물

운영 전환, 배포 보고 `.fullops-squad/docs/evaluations/qa-reports/SAR-RELEASE-002-OPS/report.md`와 명령·smoke 기록, 갱신된 운영 문서, 완료 보고.

## 제약·협업·후속

workspace 밖에 파일을 만들지 않는다(private 도구 디렉터리 제외). 다른 역할 인박스·PLANS·board·GitHub는 수정하지 않는다.

## 완료 보고

브랜치 / SHA / 변경 이유 / 검토 필요 / 검증한 것과 못 한 것 / 산출물·로그 링크 / 남은 일.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
