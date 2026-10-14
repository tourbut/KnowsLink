---
title: 공개 첫페이지 제품 결정과 DEV·tester 인계
status: review
updated: 2026-10-14
owner: designer
tasks: [SAR-LANDING-GUIDE-001]
summary: 서비스 소개·설치·사용예시·로그인의 카피와 상태를 정하고 원천·검증·후속 책임을 연결한다
---

# SAR-LANDING-GUIDE-001 실행 기록

## 범위와 고정 기준

제품 코드 없이 공개 안내 제품/UX 설계를 작성했다. 결과 정본은 [D04](../../design-docs/mockups/service-member-ui.md), 요구는 [D02 LG-01–08](../../planning/product-specs/SAR-HUMAN-PEER.md), 구현 예정 사용 흐름은 [D11](../../operations/user-guide.md)이다. 정확한 결과 커밋은 이 문서를 포함하는 `fullops/designer`의 worker_done `SHA:`다. 설계 완료와 제품 구현·동작·시각·운영 수락은 다르다.

- 기준 ref: `3230ee566ccf455e58a6fc3b43f5cd4b3a72aba0`.
- 실제 착수 HEAD: `93271bf0d2aa8b6677a1f3d3135dc5dd708cddf8`; 준비 문서와 탐색 원천이 포함된 깨끗한 체크아웃이었다.
- attempt: `38f444fbc6944d27be42b0f8a4f27b7b`.
- Task `task_eb743a73b6d8`, Dispatch `ctx_074184d23af0`, worker `term_f05a63c9-ed45-4529-9183-1b0e357222ab`.
- 복귀: coor `term_38ae16dc-0ca4-4beb-ab00-d1e1b641f061`, Run `run_5534969da2d0`.
- 실행 codex/gpt-6-astra/high, 주 스킬 frontend-design, FullOps frontend/work/orca/deliverables와 Orca orchestration 적용. 최신 사용자 지정이 과거 Sol/Astra 제외 기록보다 우선한다.
- fullops-common-0.3.3·project.md·문서 규칙을 적용했다. 단일 화면 설계라 하위 worker는 만들지 않았다. PLANS/board·설정·제품 코드·운영은 변경하지 않았다.

## 결정과 기존 지시 변경 근거

사용자 최신 요청은 첫페이지에서 서비스 소개·설치방법·사용방법·실제 기능에 맞는 사용예시를 읽고 로그인해 사용하도록 했다. 따라서 기존 D04의 로고·짧은 문구·로그인 버튼만 남기던 방향을 공개 안내에서 대체한다. 과거 최소문구·가짜 회원 메시지 금지·실패/수락 이력은 삭제하지 않았다. 공개 문서의 명확히 표시된 가상 예시와 회원 대화 데이터는 다르다.

1. 공개 `/`는 로그인 상태와 무관하게 읽는다. 유효 세션은 '대화로 이동'으로 `/home`에 들어간다. 기존 자동 홈 redirect 변경을 LG-02로 명시했다.
2. 첫 화면은 목적·로그인 다음 행동·사람 아래 에이전트 관계를 보여준다. 설치 명령과 로그인 폼을 반복하지 않는다.
3. 시작 순서는 이름·에이전트·친구 수락·대상 선택/승인이다. 웹 설치 불필요와 agent 연결 필요를 함께 안내한다.
4. Grok 로컬 Command MCP와 Dots Remote HTTPS MCP를 구분한다. 접힌 요약·펼친 절차·준비물·실제 Command/Arguments·installer 실행 조건·비밀 아닌 신규 설정을 제공한다. 긴 소스 준비·버전별 상세는 README로 연결한다. 신규/업데이트와 키 보관 위치를 분리한다.
5. 예시는 웹 직접 작성과 agent 발송 요청 두 가지다. 합성 이름과 일반 text만 사용한다. 자동 답장·일정 조회·예약·실제 성공 사례를 꾸미지 않는다.
6. Google-only/both/email-only/none/error/expired/authed 상태를 정했다. 공개 안내에 개인정보를 렌더하지 않고 실제 인증·Remote 복귀 경계는 보존한다.
7. 기존 여섯 색과 글꼴·44px 동작·focus를 유지한다. 공개 hero/h2와 구획 폭·간격만 확대한다. 1440/390/320 와이어프레임과 자기 비평 두 차례를 정본에 남겼다.

## 사실 확인과 범위 판단

README, 현재 구조/상세/운영 문서, OPS 기록과 member start/auth/template·human route를 대조했다. `startPage`의 유효세션 redirect, `Google`/`EmailLogin` 분기, POST action과 기존 오류·코드 확인 화면을 확인했다. human route의 이름 미설정 복귀와 내 계정 경로를 확인했다. 기술 API나 새 라이브러리를 설계하지 않았으므로 이번 작업의 Context7 조회는 해당 없음이다. 호스트 메뉴·요금제·Events 최신 지원을 새로 검증했다고 주장하지 않고 레포의 설치 정본과 계정별 확인 조건을 인계했다.

Grok 수신 process와 webhook routine, Dots Events 구독과 화면 지시, 발송마다 사용자 승인, 최대180초 전달 기한·최대24h 웹 사본을 원천과 맞췄다. 최신 OPS의 미실행 실계정 왕복·장시간·Grok 현지 업데이트·실제 rollback 한계는 유지한다. 운영 성공 통계·내부 SHA·검수 상태를 공개 카피로 사용하지 않는다.

D02/D04/D11만 갱신했다. D10은 프로그램 책임·기술 계획이므로 이번에는 no_change이며 DEV가 구현 후 영향을 갱신한다. D01 사업·D03 구조·D05 인터페이스·D06/07/08/09 데이터·D12 운영·D13 이행 계약은 바꾸지 않았다. 과거 D02/D04 원천은 이미 현행 링크를 제공하므로 재작성하지 않았다. PLANS 추천은 coor 소유라 no_change다.

## 탐색 원본과 착수 확인

원 packet HEAD는 기준3230ee5였고 실제 착수는93271bf였다. 공식 `jev_packet.py --force --write-inbox`로 갱신했다. 첫 실행은 완료 보고 절이 없는 인박스에서 블록을 삽입하지 못했다. 완료 보고 절을 추가하고 명시 D02/D04/D11 옵션으로 다시 실행했다. 두 이전 packet은 history에 보존했다. 이후 공식 `jev_packet.check(..., required=True)`가 HEAD·attempt·지시 fingerprint·원천 hash를 통과했다. 이 검사는 문서 편집 전에 수행했다.

- 지시 fingerprint: `f6b9d61c2e2fb0faefa1851b8e55afc12f02a8f93a357dd1e2f9156606f0486c`.
- 최종 착수 packet fingerprint: `5ea3d41dc336e3c31e8b15f90ccb2ecbb717726008de536eebd7eeda5ada2664`.
- 원본 history: [기준 원본](../../evaluations/jev/history/SAR-LANDING-GUIDE-001-packet-32691c08d32341a6a2ef3965ce58d264.json), [중간 재생성](../../evaluations/jev/history/SAR-LANDING-GUIDE-001-packet-db5a05183a314f0d820e731e12b50fa1.json).
- [현재 packet](../../evaluations/jev/SAR-LANDING-GUIDE-001-packet.json)은51경로·partial true·unknown19다. 과거 producer의 HEAD/지시 연결 미확인은 새 의미 검색 PASS로 바꾸지 않았다.
- [원천 감사](../../evaluations/qa-reports/SAR-LANDING-GUIDE-001-test/source-audit.json)는 원본/현재 후보·잔여·큰 파일86경로의 hash·제목·관련 발췌를 보존한다. 전체 역사 문서 재감사는 아니다. 기존 최소 로그인과 Sol 충돌은 위 최신 요청으로 해소했다. 큰 member.go는 start/auth 구간, PLANS·OPS·과거 logs는 해당 사실/제목을 확인했다.
- [outcomes](../../evaluations/jev/SAR-LANDING-GUIDE-001-packet-outcomes.json)는 모든 비선택 path/category의 처리와 이유를 연결한다. 추천 document_update는 수정 의무로 해석하지 않았다.

## 검증과 인계

검사 명령·원 종료코드·제품 입력 동일성과 미실행은 [설계 QA](../../evaluations/qa-reports/SAR-LANDING-GUIDE-001-test/report.md)에 기록한다. 설계의 링크·front matter·문구·추적성 검사는 실제 UI 동작·접근성·시각 수락이 아니다.

공식 work.py new로 비어 있던 [DEV 인박스](../../../handovers/to_dev.md)와 [tester 인박스](../../../handovers/to_tester.md)를 만들었다. 각각 SAR-LANDING-GUIDE-001-DEV/TESTER이고 queued다. 설계 독립 리뷰/main 통합 뒤 DEV를 배정하며 tester는 구현 SHA와 fixture가 준비된 뒤 배정한다. 현재 인박스 base는 준비93271bf이며 coor가 실제 배정 시 최신 통합 후보·각 역할 packet을 연결한다. 오래된 designer task/dispatch는 사용하지 않는다.

DEV는 기술 계획·구현·필수 make lint/test·고정 SHA·격리 서버/fixture를 맡는다. tester는 실제 합성 동작과 인증 경계, designer는 실제 캡처 시각 검수, 별도 세션은 fixed SHA 코드/문서 리뷰를 맡는다. 모두 Astra high를 유지한다. 미해결 critical/high는 수락을 차단한다. coor가 PLANS/board와 실제 main/origin 통합을 관리한다.

제품 구현·새 브라우저·실계정·설치 실행·외부 송신·운영 배포는 수행하지 않았다. 임시 검증·탐색 JSON은 QA/jev 경로에 보존했고 삭제·이동 정리는 없다.

coor 후속 msg_6593a1e11cb6은 설치 안내가 README 링크만으로 끝나지 않는지 점검하도록 요청했다. 제품 관점에서 수용해 Grok 접기에 최소 실제 등록 값과 준비물을 보강했다. 복사 버튼·설치 자동화·새 의존성은 추가하지 않았다. 설정 예시는 신규 전용이며 기존 agent/키의 재등록을 유도하지 않는다.
