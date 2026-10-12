---
title: 사람 peer 개편 실제 계정 확인
status: draft
updated: 2026-10-12
owner: coor
tasks: [SAR-HUMAN-PEER-001, SAR-HUMAN-PEER-001-OPS]
summary: 실제 Google 로그인과 운영 적용 전후 회원 화면 확인 범위를 기록한다
---

# 실제 계정 확인

사용자가 자기 이메일을 제공하고 Orca Google 화면에서 로그인을 직접 완료했다. coor는 운영 `https://link.knowslog.com`의 인증된 회원 화면에서 확인된 마스킹 이메일, 활성 에이전트 2개, 그중 키 등록 에이전트 1개와 기존 관계 1개를 확인했다. Google 비밀번호·OTP·cookie·개인키는 조회하거나 기록하지 않았다. 이 확인은 운영 적용 전의 이전 UI다.

수락 제품은 `26094e90648aa71e8d07b381e9891acb37438730`이며 main/origin과 당시 모든 idle clean 역할에 동기화했다. 제품 exact HEAD 등록 `make lint`/`make test` native0, FullOps ERROR0/WARNING19/unavailable0이다. 최종 게이트는 task 전용 `AppData/Local/KnowsLinkDevTools/humanpeer001-coor/final-lint.json`, Git 공용 `fullops-integration/SAR-HUMAN-PEER-001-final.json`에 보존했다. 경고의 SIZE·fixture SEC·SLOP·test 스크립트 DEP는 DEV 및 독립 코드 리뷰의 수락 이유를 적용한다. 새 패키지 의존성은 없다.

운영 적용은 `SAR-HUMAN-PEER-001-OPS`, task `task_256846c828a2`, dispatch `ctx_e9796171e27a`로 인계했다. fresh Claude sonnet5.5/high의 requested/effective와 turn_started를 확인했다. 운영 backup/복원 리허설·migration·보존 비교와 적용 결과는 OPS 완료 기록에서 별도로 확인한다. 사용자 표시 이름과 기존 두 에이전트 이름을 질문했고 아직 응답 전이다. 실제 친구 요청·외부 대화 발송은 수행하지 않았다.

제품 수락의 제한은 원 tester scenario/report와 UI delta·scoped visual 기록을 따른다. 실제 키보드/200%, 독립 live 벽시계·queued/leased/gate 등 not_run을 유지한다. 검수 세션 환경자격 자기로그 노출은 PLANS/visual final 문서에 값 없이 보존했고 해당 dispatch 종료·권한 폐기를 확인했다. 개별 hook/launch token 별도 회전은 미확인이다.

## 적용 후 확인

OPS가 제품 `26094e9` 적용·보호된 신규 backup 및 격리 복원 리허설·goose2/trigger1·relay healthy·local verify0와 기존 상태 보존을 보고했다. 원 보고와 실패는 OPS 폴더에 보존한다. coor는 기존 로그인 탭을 reload해 `/home/name`으로 이동함을 확인했다. Google에서 관측한 본인 이름으로 표시 이름을 저장했고 새 대화 홈·친구·내 에이전트 화면을 확인했다. 사용자 명칭 답변 전 기본 구분 이름을 쓰겠다고 알린 뒤 기존 두 agent를 '내 에이전트 1'과 '내 에이전트 2'로 저장했다. 각각 연결 필요/등록됨 상태를 다시 GET으로 확인했다. 키를 새로 발급하거나 철회하지 않았다.

사용자가 요청한 자기 현재·향후 agent 연결 안내를 읽고 checkbox를 선택해 동의를 저장했다. 새 대화에서 서로 다른 두 이름의 발신/수신 선택이 나타났고, 다시 내 에이전트 GET에서 자기 연결 동의 prompt가 사라졌음을 확인했다. 실제 친구는0이며 외부 요청·메시지 전송은0이다. 상대의 코드와 상호 수락이 있어야 친구의 에이전트 목록이 나타난다. 미연결 agent의 실제 클라이언트 설치·키 연결은 이 확인으로 완료되지 않는다.

Orca click 후 화면 이동이 확인되지 않아 관측한 native form의 `requestSubmit()`으로 profile/agent-name/self-connect를 제출했다. 필드 입력과 응답의 실제 GET 표시를 확인했으며 직접 API 호출·cookie 조회·DB 쓰기는 하지 않았다. 이 fallback을 실제 마우스/키보드 입력 PASS로 표시하지 않는다.

운영 screenshot `AppData/Local/KnowsLinkDevTools/humanpeer001-coor/live-agents.png`를 직접 열어 두 이름·닫힌 상세·새 메뉴·이전보다 짧은 기본 문구를 확인했다. 관측 viewport1201×841/DPR1.095445, 가로 overflow=false다. 최초 `rect<44` 단언이 여러 요소를 잡았으나 재측정은 computed min-height44px/실제 rect43.988964px로 0.012px 미만의 subpixel 반올림이다. 이를 제품 제어 크기 결함으로 올리지 않는다. 이 관측은 독립1440/390/320 시각검사를 대체하지 않는다. 키보드·실200%·mobile 실제 계정·실친구 왕복은 not_run이다.

운영 DB 복원과 별도로 backup의 암호화 및 최대7일 보존 집행은 이 작업에서 확인하지 않았다. D02의 보존 요구가 검증됐다고 주장하지 않는다. 보호된 파일 권한과 복원 검증은 OPS 보고에 있는 확인 범위다.
