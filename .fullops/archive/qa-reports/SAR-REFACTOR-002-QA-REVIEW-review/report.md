---
title: 독립 리팩터링 QA의 고정 결과 검토
status: review
updated: 2026-10-14
owner: coor
tasks: [SAR-REFACTOR-002-TESTER, SAR-REFACTOR-002-QA-REVIEW]
summary: 고정 QA 소스·실제 경계 실행·최종 등록 gate와 인계의 독립 검토를 수락한다.
---

# SAR-REFACTOR-002-QA-REVIEW

## 범위와 독립성

base/merge-base는 실제 tester 지시서와 최종 lint의 `bcb1dd6629b508a30ec39a5d3141c8ce0a558597`이며 head는 최종 `e9b5cbabde61335fb9e6f4eb52ca2505042732eb`다. QA가 실행한 제품은 bcb, QA 착수 기록은8c08a1b이다. 제품 코드가 바뀐 새 후보로 혼동하지 않는다.

검토자는 GPT-6 Astra high, 실제 CODEX_THREAD_ID `01a139ab-4840-7d71-b94b-e033b4b85efb`다. tester 작성 세션 `01a139b0-8386-7221-bb61-44399fb01ea1` 및 원 DEV 부모/자식과 다르다. 설치 FullOps1.4.2의 snapshot/prepare로 만든 detached `C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-1ad71516546644c4a216ea3e53223a41`은 읽기 전용이다. 결과만 coor QA에 썼다.

AGENTS/FULLOPS, fullops-common-0.3.3 README·coding-style·testing·security, project, document-writing, fullops-review 및 open-code-review-delegate를 적용했다. 10개 governing 입력은 이미 직접 읽은 DEV final3e5와 Git blob이 동일하다. OCR1.13.0은 고정 파일/규칙 준비에만 사용했고 의미 판단은 검토자가 수행했다. test level lite이며 필수 등록 gate를 생략하지 않았다.

## 소스·경계·증거 판단

확정 critical/high/medium/low 결함은 0건이다. 전체44개 `(path,status)` 중 reviewed32/skipped12이며 result.json에 개별 사유를 기록했다. 새 QA 코드·원 실패 코드·실행 helper·현재 인계/phase/context·원시 실행과 benchmark를 읽었다. 연결 제품 구현은 FINAL-MAIN/CODE에서 직접 검토한 동일 bytes만 재사용했다. 모든 변경 파일의 SHA256은 evidence-manifest.json에 있다.

[FINAL-MAIN](../SAR-REFACTOR-002-FINAL-MAIN-REVIEW-review/report.md)에 보존한 24개 QA 원문은 exact e9 Git blob이며 두 신규 assertion hash는 input-identity.json과 일치한다. cmd/internal/adapters/db/scripts/Go/Make/Docker scope는 bcb→e9 diff0이다. 최종 DEV3e5와 이 범위 차이는 이미 검토한 중복 Inbox assertion 한 줄뿐이다. bcb 독립 QA를 final3e5의 신규 실행으로 표현하지 않는다.

Go 신규 assertion은 queued/delivered × active/revoked 4조합에서 TTL 1ns 전·정각, transport/lease/claim 제거, delivered/received 보존, 참여자 전용24h 이력과 replay 제거를 검사한다. persist 이전 ACK는 거부하고 이후 수락한다. 합성 claim 직접필드와 실제 DB HTTP claim 회귀의 범위를 구분했다.

callback 마지막15초 lease에서 3회 한도, 새 도착의 lease 우회 금지, 정확 만료 뒤 새 event/첫 시도, 과거 event 늦은 ACK의 새 lease 보호, 정상 성공 idle을 검사한다. 관계·메시지·구독·설치 만료·새 도착 없음의 5개 음성 경계도 있다. 기존 호출 간 current-authority 회귀를 같은 실행에서 재검증했다.

Inbox 신규 matrix는 unread 두 본문의 원 JSON 전체를 EACCES/EIO/EPERM/EACCES마다 비교한다. ENOENT만 통과하고 순서대로 두 본문을 수신하며 read marker/중복/빈 Inbox를 검사한다. mock은 fresh marker stat에만 한정하고 finally 복원한다. 원래 prune-after-claim 손실 구현은 이 assertion을 통과할 수 없다. 오래된 marker 삭제 분기는 실행하지 않는다.

최초 Go 두 delivered 실패의 원 테스트·native1은 보존됐다. ACK completion은 received이며 sweep은 빈 completion만 실패로 바꾸므로 수정 기대값은 기존 계약과 일치한다. 제품 결함 RED/GREEN으로 바꾸어 기록하지 않았다.

## 실제 실행과 측정

독립 qa-race 로그의16개 최상위 검사, Node 기존 회귀/새 matrix, 별도 합성 Postgres17 HTTP/claim4개 검사(5.512초)는 실제 native0이다. commands.jsonl·prepare/post-clean 및 run.py의 fixed SHA/clean 검사와 일치한다. 원본 DEV 전체 integration111.338초는 재사용으로 분리했고 최초 adapter build 선행 실패와 bcb FullOps ANTI-004 ERROR1을 숨기지 않는다.

run.py는 known_hosts/RejectPolicy, 인증값 출력 가림, task 고유 clone, 외부 overlay, loopback 합성 DB CPU1/512MiB, native 종료코드, finally DB stop을 사용한다. 운영/계정 발송을 수행하지 않는다. reviewer는 실행 helper를 다시 실행하거나 fixture를 삭제하지 않았다.

before/after 모든 benchmark 중앙값을 각 원시 표본에서 재계산해 일치했다. 10000개 sweep98.0%, transactionCPU24.8%, 4구독claim68.6% 감소와 1구독8.5% 증가가 맞다. 공유host의 짧은 순차3×3 측정이며 실제DB7.8% 관측과 독립 CPU 측정을 구분한다. throughput·신뢰구간·운영 부하 보장은 없다. benchmark 줄끝 공백2줄의 정규화와 외부 원본 보존도 명시했다.

최종 actual FullOps1.4.2 lint는 외부 qa-refactor002/final-lint.json이다. head e9/base bcb, ERROR0/WARNING1/unavailable0, completion.json의 native0 및 실제 gate/commands.jsonl과 일치한다. config/policy hash는 DEV final의 동일 등록 입력이다. task-local make.cmd는 실제 run.py로 전달하고 원 종료코드를 반환한다. exact e9/clean Linux clone의 make lint와 make test가 모두0이며 원문을 execution-manifest와 함께 복사했다. Go vet/module·TS format/lint/typecheck·설정 검사와 Go race·adapter 전체 체인·shell guard가 실제 실행됐다. lint의 예상 config 거부 stderr는 실패 숨김이 아니다. 이전 bcb lint의 ERROR1을 이 QA의 PASS로 덮지 않는다.

SIZE-001 WARNING1은 PLANS1794줄(이전1790)의 순증4줄이다. coor 준비/실패 보고 이력이며 과거 증거 삭제나 문서 분할을 요구하지 않는다. 새 dependency·공개 schema·제품/UI 변경이 없어 DEP/시각/테마 검수는 해당 없다.

## 생성 증거와 인계

packet214개 경로의217 path/category가 outcomes와 정확히 같고 attempt/input SHA가 일치한다. completed20/no_change197 및 모든 고유 reason·uncertainty를 읽었다. route의 과거 producer SHA 불일치와 unknown/partial을 현재 의미 검색 성공으로 승격하지 않았다. 현재 task/role/attempt/actual session과 원 지시서를 직접 확인했다.

생성 find/context/history/packet의 추천 의미 전체 재감사는 생략했다. 변경 JSON은 모두 parse했고 hash·head/attempt·현재 outcomes 결속을 확인했다. packet-audit220개 중215개는 고정 착수 Git LF/CRLF hash와 일치한다.4개 과거 handover는 보존 tester의 mixed-newline bytes hash와 일치하고 정규화 내용은 고정 Git과 같다. 완료 전 동적 to_tester hash1개는 고정 착수/완료 상태와 달라 당시 bytes 재현을 주장하지 않는다. 원 지시서·attempt·완료 체크·전문 보고가 final archive에 있고 현재 inbox0byte임을 직접 확인해 의미 인계를 보완했다. 이 한계는 제품/QA assertion hash와 별개다.

과거 tester context는 원문 그대로 archive에 보존되고 신규 교훈4줄만 추가됐다. 현재 context는 원 archive 링크의 contexts 기준을 명시한다. 오래된 과거 링크 전체 렌더는 생략했다. PLANS/board/cleanup의 coor 기록은 당시 진행/실패·자동 cleanup 한계를 유지한다.

대화 미참조 점검은 deliverables/README에서 시작했다. 요구와 다음 담당은 archive, 구조/계약은 D03/D05/D07/D09/D10 및 system-design, 구현/실행/한계는 QA report·phase·input/commands, 운영 복구는 기존 operations를 확인했다. 이번 QA는 제품·SQL·운영 변경이 없어 DNN을 다시 쓰지 않은 판단이 타당하다. 각 문서의 로컬 보고서/근거 경로가 존재하며 main 통합 책임은 coor다.

## 결론

고정 e9 독립 QA 소스·문서·검증 근거를 수락한다. 자동 review.py check는 기록 정합성 검사이며 실제 테스트나 의미 검토를 대신하지 않는다. main/origin 통합은 coor 책임이며 이 검토에서 수행하지 않았다.

미실행: 운영·실계정·외부 callback/host wake·UI·장시간/최대부하·백업복원·중복 전체 test. 기존 local claim 이후 crash/read/write window 및 JSONB singleton/row lock 한계는 남는다. 현재 snapshot과 task-local fixture는 보존했다.

최종 `review.py check --task-key SAR-REFACTOR-002-TESTER`는 native0, reviewed32/skipped12/total44로 통과했다. 원 출력은 check.txt다.
