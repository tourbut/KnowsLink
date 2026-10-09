---
title: Google 클라이언트 연결 DEV 검사
status: draft
updated: 2026-10-10
owner: dev
tasks: [SAR-GOOGLE-CONNECT-001-DEV]
summary: Windows Google 연결 핵심 검사와 운영 수락 한계를 기록한다
---

# Google 클라이언트 연결 DEV 검사

검사 담당은 DEV이며 테스트 레벨은 lite다. 기준 ref는 `f9f7675be6c9502bd6ab2f810bd42ff62a26a174`이고 착수 HEAD는 `f6daf938e80ed28c9e9500baad6fa3c61ecd502c`다. commit 전 변경 후보의 검사이며 최종 고정 SHA와 commit 후 등록 게이트 결과는 worker_done 및 Git 공용 디렉터리의 `fullops-gate/google-connect-lint.json`을 따른다. 독립 리뷰·QA와 제품 운영 수락은 부모 담당이다.

| 명령 | 실제 종료코드 | 증거 |
|---|---|---|
| `make lint` | 0 | [product-lint.log](product-lint.log) |
| `make test` | 0 | [product-test.log](product-test.log) |
| `go test -race -tags integration ./internal/relay -run 'TestGoogleDeviceHTTP\|TestGoogleHTTP\|TestEmailIdentity' -count=1` | 0 | [integration.log](integration.log) |

product-lint는 gofmt/vet/module 검증, Prettier/ESLint/TypeScript, Compose, shell/Python 및 공개 ingress 경계 self-check를 수행했다. product-test는 Go race와 MCP initialize/discovery/held, 시험 경계, 수동 connect, Google login 로컬 서명·자동 저장·별도 키·권한·만료 검사를 수행했다. 통합 검사는 격리 Postgres 17과 합성 RSA OAuth/PKCE 서버로 같은 Google 회원의 두 연결·별도 agent/credential·명시적 동의·최근 인증·Origin·수동 confirm 우회 거부·replay 및 기존 Google/이메일 회귀를 검사했다.

Windows의 기존 Node24/make 부재를 로컬 Node22.22.2/npm10 및 w64devkit2.10.0으로 복구했다. 등록 명령과 실패 기준은 변경하지 않았다. 초기 registry digest/format 실패는 LF checkout으로 해결했다. 초기 bundle CLI 부작용은 entry 파일명 검사로 해결했다. WindowsPS Set-Acl module 충돌은 .NET ACL로 해결했다. 새 최근 인증 테스트의 403 기대값은 기존 422 계약에 맞게 고쳤다. 이 중간 실패들을 최종 PASS로 소급하지 않는다.

직접 UI 검수는 기존 Service handler를 쓰는 합성 localhost fixture에서 Chrome으로 수행했다. desktop과 mobile 390×844에서 키 지문과 동의 버튼을 확인했고 document scrollWidth/clientWidth는 모두 375였다. 명시적 동의 후 approved 화면과 잘못된 요청 거부 화면을 확인했다. 화면 캡처 파일은 저장하지 않았다. 임시 프로세스와 격리 DB는 검사 후 정리한다.

실제 Google 계정·Cloudflare Access/Tunnel 적용·외부 Bot 설치·실제 왕복은 실행하지 않았다. 기존 개인키·비공개 .env 링크와 전달 WIP를 보존했다. 실제 운영 검증·복구는 D12, 사용자 절차는 D11, 추적성·기술 선택은 실행 기록과 packet-outcomes를 따른다. SQL schema·dependency 변경은 없으며 기존 held/fail은 유지한다.
