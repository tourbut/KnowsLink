---
title: SAR-GOOGLE-CONNECT-002-DEV 독립 리뷰
status: draft
updated: 2026-10-10
owner: tester
tasks: [SAR-GOOGLE-CONNECT-002-TESTER]
summary: Access 없는 Google 연결 후보의 독립 고정 SHA 리뷰와 lite QA 결과
---

# SAR-GOOGLE-CONNECT-002-DEV 리뷰

- 검토자 / CLI / 모델: tester, Claude Code, Sonnet 5.5. 실제 세션 `c6511eeb-207c-46e2-8d42-33378e01d5bc`. 구현자 세션 `01a1218c-e248-7962-a471-d3a01468503b`와 다르다. AI 판단이며 OCR 자동 판정이 아니다.
- base SHA / head SHA / merge-base: `5af28c9fcfc289a845f9730aed34d32308e8b5a6` / `dbbe2f17f353282430f28b6e3ed98e72d4450e69` / base와 동일.
- OCR 버전 / 적용 규칙: open-code-review v1.12.13 delegate 모드(LLM 미사용). review/rule.json sha256 `f1061481…94de8`. 공통 기준 fullops-common-0.3.3, 예외 없음. 읽은 문서: common README·coding-style·testing·security, project.md, FULLOPS.md, contexts/tester.md.
- 요구사항·완료 기준 원천: `handovers/to_tester.md`, DEV exec-plan·QA report, D12 ops-guide 최신 절.
- 전체 변경 / OCR 대상 / 제외 / reviewed / skipped: 37 / 28 / 9(qa-reports 증거 로그) / 28 reviewed / 9 skipped(manifest·report.md로 존재·표 확인).
- lint: ERROR 0 / WARNING 0 / 실행 불가 0. DEV 최종 HEAD gate를 재사용했다. 출처 `D:/workspace/KnowsLink/.git/worktrees/dev/fullops-gate/google002-final-lint.json`, sha256 `209bbbec6ff6ee0a2b46978ae215f638f1f3b5637ab1432a9b0d8f0b534fffaf`, head=dbbe2f17, config_sha256은 merge-base의 lint.json과 일치. product-lint·product-test 모두 실제 exit 0. 새 SHA 실행이 아니라 같은 조건의 증거 복사다(`lint.json`). Go 테스트 일부는 gate 출력에서 `(cached)`였다.
- SIZE-002: 예상 300줄 이내 / 실제 추가 217줄(gate summary) / 차이 없음.
- DEP-001: 변경 없음. 새 제품 의존성·SQL·UI 없음(package*.json 미변경 확인).
- UI 디자인: 해당 없음. 화면·테마 변경이 없고 문구·manifest·서비스 안내만 바뀌었다.

## 검토 범위

28개 대상을 모두 reviewed로 기록했다. 제품 코드(login.ts, beta.sh, public-ingress.yml, test_public_config.sh, device_test.go, package_plugin.py, manifest, Makefile)는 diff 전체와 private-files.ts·check_public_ingress.py·relay 라우팅을 읽었다. 문서·packet JSON은 요약·정합성 위주.

## 발견 사항

1. low / public-ingress.yml / `/home/../owner`는 `/home(/.*)?`에 매칭(rule0). relay ServeMux가 정리·301 후 재요청은 rule1 404, owner 인증도 필요. 노출 없음. 기존 allowlist 상속. 해결(수락).
2. low / test_public_config.sh / `stat -c %a`가 Windows Git Bash에서 실패(exit1). 서버(Linux)에서는 통과. 제품 결함 아님.
3. low / beta.sh / validator 실패 뒤 후보 파일이 남아 재실행이 거부됨. live 파일 불변. 운영자는 후보 삭제 뒤 재실행.
4. medium(잔여, 수락) / ops-guide D12 / 익명 Device cap2000+만료 후24h tombstone으로 신규 연결 포화 가능. 기존 회원·키·철회는 영향 없음. coor가 공개 전 인지해야 한다.

critical/high 없음.

## 보안 핵심 확인

- allowlist+404: 실제 cloudflared 2026.9.1에서 `/`, callback, `/home`, `/connect/<43>`, `/v1/connect/start`, `/v1/text/send`, `/v1/keys/a/k`는 rule0. `/owner`, `/v1/owners`, `/v1/keys`, `/v1/authorize`, `/v1/test/pull`, `/healthz`, owner-revoke/key-revoke/agents, 미등록·대소문자·`//owner`·다른 hostname은 rule1(404).
- render-only: 0600, 새 파일 exclusive(`set -C`), 재실행 거부, 잘못된 UUID·validator 실패(exit7) 전파, live config 불변을 서버에서 직접 확인(test_public_config.sh exit0, 수동 렌더 exit0/재렌더 exit1).
- 운영 순서·인계(D12 4~6): ingress 차단 확인 뒤 두 KnowsLink 앱(`fc81b205…`, `bd210310…`)만 백업·제거, 광범위 Everyone·새 Bypass·청구 동의 금지, 실패 시 Tunnel 중지·member allowlist+404 복구, owner/admin 무인증 fallback 금지. 완전하다. trial은 운영 앱 부재라 공개 복구하지 않는다고 명시. 과거 Access 계획(001)은 이력, 002 절이 정본임이 문서에서 구분된다.
- Google 세션/소유권·키 분리: device_test가 기존 다른 회원 callback 시 회원·세션 수·owner·state 불변, 두 클라이언트 별도 agent·credential을 고정한다. 로컬 `go test ./internal/relay -run 'TestDevice|TestGoogle' -count=1` exit0.
- 상위 폴더: `mkdir(dirname, recursive)` 뒤 `privateDirectory`의 exclusive `mkdir`+ACL 검증이 그대로다. Node22.22.2 Windows 직접 확인: 새 상위 `a/b/c` 생성(private.pem, login.json), 기존 폴더·빈 기존 폴더 모두 EEXIST로 거부하고 기존 폴더는 비어 있음.
- manifest/skill/README는 9도구·명시 동의·held 기본·Access 불필요·결제 금지와 일치. mcp.test가 도구 9개를 검증한다.

## 검증 및 남은 제약 (lite QA, 직접 실행)

| 검사 | 실제 exit |
|---|---|
| Windows Node22.22.2(SHA256 검증한 공식 zip, 임시 폴더) `npm ci`/build | 0 / 0 |
| login.test, connect.test, trial-boundaries.test, mcp.test | 0 / 0 / 0 / 0 |
| 새 상위폴더·기존 폴더 거부 probe | 0 |
| 서버 고유 temp `test_public_config.sh`, `check_public_ingress.py` | 0 / 0 |
| 서버 실제 cloudflared validate + rule 26경로 | 모두 기대 rule |
| 서버 전후 docker ps 비교 | cmp 0 |
| 로컬 Go TestDevice/TestGoogle | 0 |

처음 mcp.test 실패는 내 임시 복사본에 `internal/relay/registry.json`을 빠뜨린 탓이며 보완 후 통과했다. `.env.server` 값은 메모리에서만 사용했고 argv·로그·Git에 넣지 않았다. 서버 임시 디렉터리·컨테이너는 정리했다. 운영 서비스는 변경하지 않았고 로컬 Docker는 쓰지 않았다. DEV 전체 검사는 반복하지 않았다. 실제 Google 로그인·Bot 설치·관계 수락·대화 왕복·운영 적용은 미검증이며 coor 후속이다. 합성 통과로 대신하지 않는다.

## 대화 미참조 인계 점검

| 확인 항목 | 정본 경로/절 | 결과 | 누락·오래된 정보·후속 |
|---|---|---|---|
| 현재 요구와 결정 이유 | exec-plans/phases/SAR-GOOGLE-CONNECT-002-DEV.md, D12 | 확인 | 없음 |
| 구조와 구현/미완료 상태 | design-docs architecture·module-design·tech-stack, adapters/README.md | 확인 | 운영 적용·실계정 수락 미완료로 명시됨 |
| 실행·검증 방법과 증거 | qa-reports/SAR-GOOGLE-CONNECT-002-DEV/report.md, server-checks.sh | 확인 | 증거 로그는 skipped, 요약표로 대조 |
| 운영·복구 | ops-guide.md `비용 없는 Tunnel 적용` 절 | 확인 | 위험·복구 순서 완전 |
| 다음 작업·담당·재개 조건 | D12 6항, DEV report 끝 | 확인 | coor가 적용·실계정·관계·왕복 수행 |
| 로컬 링크·절 접근/지원 한계 | README의 ops-guide 앵커 | 확인 | 앵커는 한글 제목 기반, 렌더러에 따라 한계 |
| snapshot 정리 후 정본 접근 | 이 리뷰 폴더 | 확인 | 정리 전 coor 병합 필요 |

## 검토 결론

수락. critical/high 미해결 없음. 잔여: medium 신규 연결 포화 한계, low 3건. 실제 운영·Google·Bot·관계·대화 수락은 coor 후속.
