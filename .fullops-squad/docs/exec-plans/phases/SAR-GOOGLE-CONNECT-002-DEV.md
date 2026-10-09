---
title: 비용 없는 Google 클라이언트 연결 보완
status: draft
updated: 2026-10-10
owner: dev
tasks: [SAR-GOOGLE-CONNECT-002-DEV]
summary: Access 없는 Tunnel 후보와 기존 Google 플러그인 연결의 최소 보완 및 검증
---

# 비용 없는 Google 클라이언트 연결 보완

## 기술 계획과 범위

기준 ref는 5af28c9, 준비 커밋은 fc09f00이다. fullops-common-0.3.3과 project.md를 적용하며 예외는 없다. DEV는 구현·배포 후보·기술문서·합성 검증을 맡는다. 운영 적용·독립 리뷰·실계정 수락·시각 검수는 coor 후속이다. 로컬 Docker는 사용하지 않는다.

1. 기존 public-ingress allowlist를 재사용해 Access 없는 config.public.yml을 별도 생성한다. 허용 경로 이외는 404이며 운영 config와 DNS를 바꾸지 않는다. 기존 Access 기반 도구는 그대로 유지한다.
2. 플러그인 manifest·marketplace·README·skill의 synthetic-only 안내를 고친다. 기본 held와 비명시 발송 금지는 유지한다. MCP/login 구현과 OAuth 자격은 재사용한다.
3. 익명 Device cap2000과 만료 뒤24h tombstone 보존을 유지한다. 포화·보존 경계·회복과 기존 다른 회원 callback 거부를 DEV 회귀로 고정한다. 조기 삭제나 새 IP 정책은 추가하지 않는다.
4. D03 D10 D11 D12 D13과 packet-outcomes를 갱신한다. 운영 서버 고유 임시 checkout·격리 DB에서 make lint/test·Google 경계를 검사하고 최종 깨끗한 HEAD gate를 남긴다.

호출/영향 범위는 beta.sh dispatch, public-ingress, check_public_ingress, Makefile, package_plugin, 두 manifest, 패키지 skill/README, device_test와 기존 login/MCP/Google tests다. 새 의존성·SQL·UI 변경은 없다. 예상 제품/문서 변경은 300줄 이내이며 증거 파일은 별도다.

## 근거와 판단

Context7 Cloudflare resolve-library-id는 monthly quota 초과였다. 공식 Tunnel config와 Access path 문서를 확인했다. Tunnel은 위에서 아래 첫 일치 규칙과 필수 catch-all을 사용한다. origin required:false는 edge Access 앱을 해제하지 않는다. 문서: https://developers.cloudflare.com/cloudflare-one/networks/connectors/cloudflare-tunnel/do-more-with-tunnels/local-management/configuration-file/ 및 https://developers.cloudflare.com/cloudflare-one/access-controls/policies/app-paths/ .

coor msg_3319725732d2는 dashboard API GET으로 현재 앱 두 개와 trial 앱404를 확인했다. 두 KnowsLink 앱만 백업·제거하고 새 Access 정책·가입·청구 동의를 만들지 않는 경로를 D12로 인계한다. 현재 trial은 운영에서 비활성이며 코드는 보존한다. 다른 호스트가 있으면 generated candidate를 통째로 덮지 않고 기존 config에 해당 호스트 블록만 병합한다.

기존 medium 익명 신규연결 DoS 한계는 남는다. cap2000이 차면 최초 요청 만료 후24h까지 새 Device를 막을 수 있다. 기존 회원·기존 키 사용·철회는 이 cap에 묶이지 않는다. 한도를 키우거나 tombstone을 지우면 제품 보존/남용 정책이 바뀌므로 이 최소 과제에서는 변경하지 않는다.
