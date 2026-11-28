---
id: handover-dev
title: SAR-MCP-SPEC-CHECK-001-DEV — 배포 전 원격 /mcp가 MCP 2026-07-28 공식 문서의 필수 요구 사항을 맞추는지 점검하고 빠진 것을 보고한다
status: approved
updated: 2026-11-28
owner: dev
summary: 배포 전 원격 /mcp가 MCP 2026-07-28 공식 문서의 필수 요구 사항을 맞추는지 점검하고 빠진 것을 보고한다
key: SAR-MCP-SPEC-CHECK-001-DEV
role: dev
goal: 배포 전 원격 /mcp가 MCP 2026-07-28 공식 문서의 필수 요구 사항을 맞추는지 점검하고 빠진 것을 보고한다
base: ce4109b440790694789dde3d00cde4f3b0513828
route: direct
from: lead
round: 0
budget: 2
test_level: lite
parallel: on
---
# SAR-MCP-SPEC-CHECK-001-DEV — 배포 전 원격 /mcp가 MCP 2026-07-28 공식 문서의 필수 요구 사항을 맞추는지 점검하고 빠진 것을 보고한다

현재 과제는 이 역할 인박스에서 관리한다. 같은 과제의 후속은 본문에 반영하고 다른 과제는 PLANS.md에 대기시킨다. 완료하면 완료 보고 전문을 채우고 `fullops done`으로 지시서·결과를 로그에 보존한다.

- 작성일: 2026-11-28
- From / To: coor(lead) / dev
- 상태: ready
- 테스트 레벨 / 병렬 방침: lite / on
- 승인된 범위: 배포 전 점검(주인 지시, 쿠어 msg_884b83260650·msg_77ae3023376a). 코드는 읽기와 분석만 한다. 작은 필수 누락(한두 파일, 기존 클라이언트를 깨지 않는 것)은 고쳐도 된다. 큰 누락은 고치지 말고 보고한다.
- 추가 승인이 필요한 행위: 원격 push, 운영 배포, 인증·OAuth 흐름 변경(CIMD 지원 추가 포함), 기존 클라이언트를 깨는 변경. 필요하면 멈추고 `ask`로 묻는다.
- 상위 의도: 없음

## 현재 상황과 확인 근거

- 운영 배포 대기 묶음(origin/main `ce4109b4`):
  - SAR-MCP-UPGRADE-001: 원격 `/mcp` 2026-07-28 dual-era. 헤더 검사 -32020, 미지원 버전 -32022, resultType·ttlMs/cacheScope, discover instructions
  - SAR-MSG-TTL-1H-001-DEV
- 주인이 공식 문서 https://modelcontextprotocol.io/docs/2026-07-28 대비 필수 요구 사항을 배포 전에 다시 확인하라고 했다. 주인 위키(concepts/mcp-2026-07-28-spec)에서 나온 후보는 다음과 같다:
  1. DCR deprecated, CIMD(Client ID Metadata Documents) 권장. 인가 응답의 issuer 검증 필수
  2. 서버 주도 요청(elicitation 등) 금지, MRTR로 대체
  3. tools/list의 ttlMs·cacheScope 힌트
  4. Roots·Sampling·Logging deprecated
- 참고: Extensions 개요 https://modelcontextprotocol.io/extensions/overview
- 지금 Grok Bot과 Dots는 OAuth DCR(loopback redirect 허용, PKCE S256)로 붙는다.

## 적용 기준과 예외

- 1차 출처(위 공식 문서, 스펙 본문, 관련 SEP)로만 판단한다. 기억이나 위키 문구를 근거로 삼지 않는다. 확인하지 못하면 '미확인'으로 적는다.
- MUST(필수)와 SHOULD(권장)를 구분한다. 서버 쪽 의무와 클라이언트 쪽 의무도 구분한다.

## 먼저 읽을 문서

- `.fullops/tasks/SAR-MCP-UPGRADE-001/analysis.md`, `qa/report.md`
- `internal/relay`의 remote_era.go, remote_mcp.go, remote_tools.go, OAuth(DCR·authorize·token·metadata) 처리

## 해야 할 일과 파일 소유권

- [ ] 공식 문서의 서버 쪽 MUST를 표로 정리한다. 열은 다음과 같다:
  - 요구 사항
  - 근거 링크
  - 지금 상태: 충족, 누락, 해당 없음
  - 파일:줄
  대상은 헤더, 버전 협상과 오류 코드, 인증, 서버 주도 요청, deprecated 기능이다.
- [ ] 인증 관련:
  - CIMD를 지원하지 않을 때 2026-07-28 클라이언트가 연결에 실패하는지 판정한다. DCR fallback이 MUST인지 SHOULD인지도 본다.
  - 인가 응답에 `iss`를 넣는 것이 서버 MUST인지 확인한다(RFC 9207 연계). 넣는다면 지금 응답에 있는지 본다.
  - Claude Code v2, Dots, Grok Bot이 CIMD와 DCR 중 무엇을 쓰는지 공식 문서로 확인한다.
- [ ] MRTR: 서버가 먼저 거는 요청을 지금 보내는 곳이 있는지 확인한다. 5순위 Elicitation 판단을 고친다.
- [ ] analysis.md의 순위가 바뀌면 고치고, 바뀐 점을 적는다.
- [ ] 결론을 낸다: 이번 묶음을 그대로 배포해도 되는지, 아니면 배포 전에 고칠 필수 누락이 있는지. 작은 누락을 고쳤다면 표적 테스트와 lint를 붙인다.

## 완료 기준과 검증

- `tasks/SAR-MCP-SPEC-CHECK-001-DEV/qa/report.md`에 MUST 표, 인증 판정, MRTR 판정, 결론을 적는다.
- 코드를 고쳤으면 표적 테스트와 `fullops lint --from <base>`를 통과한다.

### UI 작업의 디자인 기준과 검증

해당 없음.

## 갱신할 산출물

- `.fullops/tasks/SAR-MCP-UPGRADE-001/analysis.md`(순위가 바뀌는 경우)

## 기대 산출물

- qa 보고. 코드 수정이 있으면 그 코드와 테스트.

## 제약·협업·후속

- 큰 누락(예: CIMD 지원)은 새 과제 후보로 적는다. lead가 PLANS 행과 주인 결정으로 넘긴다.
- 하위 에이전트: 공식 문서 조사는 읽기 전용 하위 에이전트로 나눠도 된다.

## 완료 보고

첫 줄은 `[완료] SAR-MCP-SPEC-CHECK-001-DEV | 브랜치 fullops/dev | SHA <HEAD> | 배포 판정 <그대로 배포|배포 전 수정 필요>`다. 원격 push는 승인 전이다. 그래서 `fullops done`을 실행하지 않는다. 로컬 커밋 뒤 worker_done만 보낸다.
처리는 `fullops-work`, 회신은 `fullops-orca` 스킬을 따른다.
