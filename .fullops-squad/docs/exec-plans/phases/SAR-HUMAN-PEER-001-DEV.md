---
title: 사람 친구·이름·연속 대화 기술 구현
status: draft
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
