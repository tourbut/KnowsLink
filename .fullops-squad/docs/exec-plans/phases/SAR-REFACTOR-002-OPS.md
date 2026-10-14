---
title: Astra 리팩터링 운영 적용 실행 기록
status: draft
updated: 2026-10-14
owner: ops
tasks: [SAR-REFACTOR-002-OPS]
summary: 고정 제품의 보존 배포와 이전 이미지 조회 실패의 안전한 복구 후보 준비 및 후속을 추적한다.
---

# SAR-REFACTOR-002-OPS

사용자가 기존 운영 배포를 명시 승인했다. 준비HEAD63ba2167352a4b6b1ec2df8782841ba768963176에서 인박스 전문·기준·current packet을 확인했다. 단일 writer 전환을 직접 직렬 수행했고 중간 승인을 요구하지 않았다. 제품 코드·공유 설정은 수정하지 않았다.

1. fresh baseline의 이전제품2747fe98653d5c6946f41c219fc256395333d7fe·clean·healthy·권한/Remote/공유 digest를 기록했다.
2. 수락8a93e38e190f1b9e7b6bb183e578e410f5dd97c6의 bundle/hash/조상·infra diff0과 final DEV 실행 입력 동일성을 확인했다. 기존 relay 가동 중 새 이미지를 build했다.
3. 11:30:57Z relay stop 뒤 epoch1001476·다른 DB client0을 확인했다. 새0600 backup/list와 network-none 복원의 전체 data digest·schema/guard·downgrade 거부를 검증했다.
4. 11:31:13Z 새 relay healthy를 확인했다. 약16초 중단이다. 보존33개·공개21개 검사를 통과했고 약3분 relay/Tunnel error0이다.
5. 후속 검사에서 이전 이미지 주소 조회가 종료1이었다. 원본을 보존하고 coor에 escalation했다. 현재 서비스를 유지하며 같은 이전 SHA의 별도 archive·rollback 태그 이미지를 만들었다. 원 이미지 동일성·실제 rollback 성공은 미확인이다.
6. D12/D13·공개 운영/업데이트 안내·QA·packet outcomes·완료 보고 전문을 작성했다. work.py finish가 종료0으로 날짜별 archive를 보존했고 빈 인박스를 확인했다. coor는 결과 고정 SHA의 독립 리뷰와 main/origin 통합을 맡는다.

복구는 현재 DB 새 백업과 같은 HumanModel writer의 검토된 소스/이미지 후보로만 한다. 과거 DB 덮기·Down·삭제/이동 정리·volume 삭제·key 변경은 하지 않는다. Grok 로컬 bundle·watcher는 별도 owner 후속이다. Dots 기존 Remote 연결은 유지한다.

원 SHA·UTC/KST·이미지/백업 hash·현재 구독·실패 원문·재사용gate·미실행·보존 경로는 [QA 보고](../../evaluations/qa-reports/SAR-REFACTOR-002-OPS/report.md)에 있다. 지시서의 과거 partial/unknown·시점 충돌은 수동 보완했으며 의미 검색의 완전성이나 과거 전문 재감사를 주장하지 않는다.
