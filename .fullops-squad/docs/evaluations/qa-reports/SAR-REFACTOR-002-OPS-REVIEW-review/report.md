---
title: Astra 운영 배포 결과 독립 리뷰
status: draft
updated: 2026-10-14
owner: coor
tasks: [SAR-REFACTOR-002-OPS]
summary: 고정 운영 결과 SHA의 배포·복원·보존·실패와 재사용 gate를 독립 검토한다.
---

# SAR-REFACTOR-002-OPS-REVIEW 리뷰

고정 결과 SHA `3562662f5f836debd22a255afcbf52ea930e965b`의 운영 결과 문서와 증거를 수락한다. 미해결 critical/high는 없다. 실제 배포 제품은 `8a93e38e190f1b9e7b6bb183e578e410f5dd97c6`이며 결과 문서 SHA와 구분한다. 이 판단은 검토자의 독립 판단이다. OCR은 대상·규칙·기록 검사를 제공한다.

## 독립성·기준·범위

- 검토자: Codex / GPT-6 Astra high. 실제 CODEX_THREAD_ID `01a13a34-bed7-7b00-884f-9afbc3bc9a33`. 작성자 실제 session `01a13a26-fcf7-7182-a926-8b14d44eaec6`과 다르다.
- base/merge-base: `8a93e38e190f1b9e7b6bb183e578e410f5dd97c6`; head: `3562662f5f836debd22a255afcbf52ea930e965b`.
- 읽기 전용 detached snapshot: `C:/Users/shin/orca/workspaces/KnowsLink/.fullops-review-a8147ab97d774fa685b11fc96603488e`. 기록은 coor의 이 리뷰 폴더에만 작성했다. snapshot 수정·원격 실행·배포·삭제·이동·commit을 하지 않았다.
- 기준: AGENTS.md, FULLOPS.md, fullops-common-0.3.3의 README/coding-style/testing/security, project.md, document-writing.md, D07/D12/D13, 보존된 두 OPS 지시서. base→head에서 공통 규칙·제품 lint 설정·OCR rule 변경은 없다. 설치 FullOps 1.4.2 fullops-review와 open-code-review-delegate를 적용했다.
- OCR 1.13.0. 전체44파일, OCR 대상35·제외9. reviewed33·skipped11, 전체 대비75%. 파일별 이유는 result.json에 기록했다. 과거 탐색 자료의 전문 의미 재검토는 생략했다. 현재 packet230쌍·outcomes와 필수 운영 증거는 확인했다.
- test level lite, subagent level standard. 이번 결과는 문서·task-local 실행 증거다. 수락 제품의 전체 코드를 다시 감사하지 않았다. UI/template/CSS·의존성·runtime diff0이므로 시각·테마·새 의존성 검사는 해당 없음이다.

## 발견 사항과 처리

선행 검토에서 Rates 수치, 현재 배포 문서 목차, 공개 운영 문서 front matter, smoke 판정 범위, 이전 smoke 입력 의존성의 다섯 보완을 전달했다. 최종 SHA에서 모두 해결됐다. Rates는 baseline4, before/frozen/after5로 정확히 구분한다. 집계만으로 정확한 기여 요청을 단정하지 않는다. metadata 본문 의미는 미검증으로 기록한다. manifest 사본·원 경로·해시와 재현 배치를 보존했다.

baseline 사본 의심은 동시 작성 중 크기 목록에 근거한 예비 의심이었다. 실제 JSON의 11:27:54Z·이전2747 SHA와 원본/사본 해시를 확인해 철회했다. 현재 결함으로 보고하지 않는다.

이전 이미지 `220edb718b…`는 build 직전 inspect 성공 후 전환 뒤 inspect 종료1이었다. 원 post-check.log를 유지한다. 원인에 대한 Docker 내부 감사는 미실행이다. 같은 이전 소스2747의 고유 rollback 태그 `12f1313a…` 재빌드는 실제 image inspect·network-none/read-only 실행 파일 존재/해시 검사로 뒷받침된다. 이를 원 이미지 동일성이나 애플리케이션 rollback 성공으로 표현하지 않는다. 다음 배포 전에 고유 rollback 태그를 확보한다는 운영 보완도 기록했다. 현재 healthy 서비스와 DB를 유지한 판단을 수락한다.

## 검증 근거

예비 원본 검토 뒤 최종 snapshot의 사본 해시·내용 연결을 확인했다. [verification.json](verification.json)에 고정 SHA 검증 결과를 기록했다.

- bundle SHA·FETCH_HEAD·조상 관계와 clean checkout, 이전2747→8a의 schema/deploy/Go dependency diff0을 확인했다. DEV3e5→8a→결과356의 runtime/build/test/config 입력도 diff0이다.
- writer 정지·2초 안정 epoch1001476·다른 DB client0, 새0600/상위0700 백업20039bytes·TOC25, network-none 복원의 전체 frozen data digest 일치와 goose2/guard1/HumanModel1 및 downgrade 거부를 확인했다.
- stop 11:30:57Z→up 11:31:13Z는16초다. 새 이미지482d3197… healthy/restart0·migration0, 약3분 relay/Tunnel error0의 실제 로그를 확인했다.
- before→frozen은 모든 기록된 map digest가 같고 frozen→after는 Rates만 다르다. Remote 전체 map·schema·설정·private 파일 권한·기존15백업·공유31컨테이너는 같다. 보존 assertion33개를 고정 JSON과 대조했다. private row·본문·자격·callback 값을 출력하지 않았다.
- public smoke21개는 상태코드와 선택 보안/redirect 헤더 검사다. 실제 발송·receive·ACK·로그인·metadata 본문 계약·호스트 UI 검사는 포함하지 않는다. 입력21개 계약·스크립트 archive·입력 원본/사본 hash가 일치한다.
- evidence-copy-audit16건의 원본/사본 해시를 직접 확인했다. 첫·보완 attempt 완료 전문2개와 빈 OPS 인박스, 최종 packet/outcomes230쌍·동일 attempt/input hash를 확인했다. 과거 partial/unknown과 실패를 현재 성공으로 바꾸지 않았다.
- `git diff --check`는 종료0이다. 실행 명령은 기존 commands.txt와 logs에 남아 있다. reviewer는 운영 호출을 다시 실행하지 않았다.

## lint와 재사용의 한계

최종 HEAD의 실제 FullOps1.4.2 결과 corrected-lint.json은 ERROR0·WARNING1·unavailable2이며 원 종료코드는1이다. 현재 PATH에서 make 부재를 독립 확인했다. product-lint(kind lint)와 product-test(kind test)는 unavailable·exit_code null을 그대로 유지한다. [lint-original.json](lint-original.json)은 원본 바이트 사본이며 [lint-provenance.json](lint-provenance.json)은 원 해시와 변경 범위를 기록한다. [lint.json](lint.json)에는 두 reason만 구체적으로 채웠다. 새 HEAD에서 테스트를 실행해 통과한 것으로 바꾸지 않았다.

등록 gate의 원 실행 SHA `3e5af5e5ef9ca993baa6745aaf5b2b1f157ace46`, Linux make lint/test EXIT0과 원 gate SHA256을 확인했다. 지시서와 공통 testing 규칙의 입력 동일성 재사용에 따라 이 결과를 수락한다. 새 변경이 없어 전체 제품 gate를 반복하지 않는다. SIZE-001은 기존 PLANS1802→1811줄이며 coor 준비 기록9줄 증가다. OPS가 기존 계획을 자르거나 역할 밖 재구성할 이유가 없다. SIZE-002·DEP-001은 없다. 약4.47만 추가 줄의 대부분은 기존 탐색 packet/history와 값 없는 증거·완료 archive이며 예상 운영 문서 범위 안의 기록 보존이다.

## 대화 미참조 인계 점검

산출물 인덱스와 docs/README.md에서 시작했다. 문서 승인 상태 draft를 배포 실패로 해석하지 않으며 오래된 날짜별 운영 절은 당시 기록으로 구분한다.

| 확인 항목 | 정본 경로/절 | 결과 | 한계·후속 |
|---|---|---|---|
| 현재 요구·이유 | handovers/logs/2026-10-14_to_ops.md의 두 지시서 | 확인 | 배포 승인·보존·무삭제·최신 문서 보완 |
| 구조·구현/미완료 | docs/README.md, D07, OPS report의 범위 | 확인 | 동일 runtime, singleton/row lock와 로컬 crash window 한계 유지 |
| 실행·검증 | OPS report·commands.txt·manifest·원 로그 | 확인 | 원 실패와 smoke/gate 한계 보존 |
| 운영·복구 | docs/operations.md, D12/D13 최신 절 | 확인 | 현재 DB 보존, 같은 writer 후보, 실제 rollback 미실행 |
| 다음 담당 | README Grok 업데이트, D13 최신 절 | 확인 | owner 로컬 bundle/MCP/watcher 갱신; coor main 통합 |
| 로컬 링크·절 | docs 목차→현재 OPS report, D12/D13→OPS report, report→증거 | 확인 | 외부 문서 사이트 재조회·전체 과거 anchor 감사 생략 |
| snapshot 정리 후 접근 | coor canonical 리뷰 폴더와 결과356의 커밋된 증거 | 확인 | snapshot만 소유권 확인 후 coordinator 정리; 외부 private 원본 보존 |

## 결론과 남은 위험

운영 결과 문서3562662의 통합을 수락한다. review.py check는 기록의 일관성을 검증하며 실제 배포 테스트를 대신하지 않는다. 원 제품 수락과 이번 새 배포 증거, 현재 환경 실행 불가를 분리했다.

실계정 왕복·장시간 알림·최대 부하·Google 로그인·Grok 현지 업데이트·실제 rollback은 이번 검토에서 수행하지 않았다. 원 이미지의 사후 조회 실패 원인은 미확정이며, 복구 후보는 동일 소스 재빌드다. 다음 배포 전 이미지 태그 고정과 실제 가용성 확인을 유지한다.
