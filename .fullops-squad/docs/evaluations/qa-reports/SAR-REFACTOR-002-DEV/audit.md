---
title: KnowsLink 구조·성능 감사
status: draft
updated: 2026-10-14
owner: dev
tasks: [SAR-REFACTOR-002-DEV]
summary: 제품 상한 내 보관 이력의 반복 인가 비용과 모듈별 점검 결과를 기록한다.
---

# 감사 범위

KnowsLink는 사람이 수락한 관계 안에서 웹과 agent의 text를 전달한다. 하나의 Go 프로세스가 웹·HTTP·OAuth/MCP를 처리한다. Postgres singleton JSONB가 상태·권한·멱등·lease를 직렬화한다.

예상 부하는 현재 제품 상한인 활성 회원 100명, 활성 agent 200개, queued/leased 100개, receipt 20000개 이하다. 운영 부하 관측값이 아니라 코드의 보호 상한이다. 측정은 100명/200 agent와 이미 만료된 delivered 사본 100/1000/10000개를 사용한다. 24시간 보관 중인 일반적인 비활성 이력이며 최대 본문이나 모든 상한의 동시 포화를 대표하지 않는다.

## 확인한 지도

| 경로 | 확인 내용 |
|---|---|
| README·docs/system-design.md | 사람 관계, 실제 운영 관측 날짜와 소스 기준 구분 |
| cmd/relay/main.go·cmd/migrate | 설정·readiness·별도 migration·cleanup/events background 시작 |
| Makefile·go.mod·adapters/package.json·Dockerfile·compose.yaml | 고정 Go/Node/DB·lint/test/build·private DB·Tunnel 기본 OFF |
| store.go·db/queries/relay.sql·writer migration | row lock·DB 시계·epoch CAS·실패 시 상태 복원·rate 보존·NOTIFY commit |
| cleanup_admission.go·capacity.go | snapshot은 입장 분류에만 사용, DB 안 재검증·소유자별 정리 예산·실패 admission 회수 |
| human.go·public_text.go·identity.go | 현재 권한 출처/세대, owner·회원 구분, 서명·멱등·TTL·보관 |
| Go remote 및 TS adapter | 별도 하위 보고서로 검사 경로·수정·회귀를 취합 |

## 개선할 항목

1. **만료된 이력에 반복하는 현재 권한 순회** (`internal/relay/store.go`, `State.sweep`).
   입력은 100명/200 agent와 10000개 만료 delivered 이력이다. 매 transaction의 sweep이 각 메시지마다 publicOwner·authority·humanMember를 호출해 회원 map을 반복 탐색한다. Windows Go1.27.1 baseline 3회에서 sweep은 35.8–50.5ms이며 160000B/10000 allocations다. 만료 처리 자체가 transport·claim을 제거하므로 만료 전 메시지에만 current를 실행하는 최소 수정으로 반복 탐색을 없앤다. delivered receipt·24시간 사본·멱등 보관과 미만료 철회는 유지한다.

2. **전체 JSONB와 전역 row lock의 구조적 처리량 한계** (`Service.transaction`).
   10000개 fixture의 JSON은 6925918bytes다. baseline transaction CPU는 116–186ms이며 SQL·네트워크 시간은 포함하지 않는다. HTTP는 입장·업무·퇴장 transaction을 사용하므로 이 값을 요청 처리량으로 환산하면 안 된다. 기존 코드의 명시적 단기 설계 한계이며 새 결함으로 과장하지 않는다. 이번 변경은 중복 순회를 줄이고 DB 저장 형식을 보존한다. 정규화·다중 writer·캐시 도입에는 별도 호환·복구 계획과 제품 부하 목표가 필요하다.

3. **읽음 claim 뒤 housekeeping 오류의 본문 유실** (`adapters/src/inbox.ts`). EACCES 재현에서 receive는 실패하고 unread 본문도 사라졌다. housekeeping을 claim 전에 완료하고 동시 stat ENOENT만 허용했다. 수정 전 exit1과 수정 후 exit0, 호출자·패키징 감사 범위는 [adapter 보고](adapter-report.md)에 있다. claim 이후 process crash·본문 read/write 오류의 기존 한계는 남는다.

4. **callback 최종 lease 뒤 새 도착 지연과 반복 검색** (`internal/relay/remote_events.go`). 세 번 finish 없이 claim한 뒤 새 메시지가 도착하면 원래 코드는 t+47초에서 t+3600초 구독 만료까지 대기했다. 새 pending은 즉시 다음 claim을 예약하고, 한 호출의 agent별 검색을 재사용한다. 10000 records/4구독 claim 중앙값은 590438→162572ns/op다. 한 구독의 비용 증가와 실제 범위는 [remote 보고](remote-report.md)에 기록했다.

## 미검사와 한계

운영 계정·실메시지·실제 host 화면·장시간 구독 갱신·운영 최대 부하·백업 복원은 검사하지 않는다. 과거 QA는 계약과 환경 참고이며 이번 통과 근거가 아니다. 전 파일의 완전 호출 그래프나 보안 무결점을 주장하지 않는다. UI template·CSS 외형은 변경 대상이 아니다.
