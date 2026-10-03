---
title: KnowsLink 서비스 기획 반입과 프로젝트 준비 요청
status: draft
updated: 2026-10-03
owner: coor
tasks: [SAR-PREP-002]
summary: 최신 서비스 기획을 기준으로 초기 골격 다음의 MVP 개발 준비 범위를 전달한다
---

# SAR-PREP-002 — 사용자 요청

이제 서비스 기획 레포에서 우리가 구상중인 서비스 기획안 불러와서 프로젝트 준비해

원천은 https://github.com/tourbut/service-design/tree/main/silent-agent-relay 다. main을 새 임시 clone으로 직접 확인했다. 최신 커밋은 7bc9ea190ea549fae8b047e850247a19322fc9c3이다. 기존 404ff834 스냅샷과 달라진 파일은 product.md와 decisions.md뿐이다. 변경은 A2A v0.3.0 개념 검토 완료, wire 비호환과 A2A wire 필드 미수입, owner 승인 대체 금지, webhook OFF 유지다. 나머지 원천 파일은 동일하다. 원천은 upstream 바이트 그대로 반입하며 이전 버전은 Git 이력에 보존한다.

이미 초기 Go relay·TypeScript adapter·Postgres·Compose·lint 골격의 독립 QA와 main 수락을 완료했다. 기존 골격을 다시 만들지 않는다. 이번 요청은 전체 서비스 기획의 개발 준비이며 제품 MVP 구현과 운영 배포를 자동으로 추가하지 않는다. 제품 목표·잠긴 규칙·MVP 포함/제외·사용자 흐름·수락 기준·미정 결정과 기능 단위 후속 인계를 준비한다. 기술 계획·구현·검증은 후속 담당 dev의 같은 과제다.
