---
title: 회원 UI 운영 계정 확인
status: review
updated: 2026-10-11
owner: coor
tasks: [SAR-SERVICE-UI-001, SAR-SERVICE-UI-001-OPS]
summary: 운영 Google 재인증과 기존 에이전트 목록 및 반응형·키보드 스모크 결과
---

# 운영 계정 확인

Orca embedded browser에서 https://link.knowslog.com/home 을 확인했다. 기존 Google 계정으로 `로그인 다시 확인` → Google 계정 선택 → 계속 → `Google 로그인 확인 완료` → `홈으로 이동`을 실행했다. callback의 authorization_failed는 없었다. 비밀번호·세션·토큰·개인키를 수집하지 않았다.

같은 계정의 기존 에이전트 4개와 기존 관계 2개가 표시됐다. 키 등록 3개와 미연결 1개를 구분했다. 관계 상대의 식별자는 기존 연결 정보다. 다른 소유자의 에이전트가 내 목록에 표시되는 현상은 관찰되지 않았다. 권한 경계의 공격 테스트를 새로 수행한 결과로 해석하지 않는다. 새 에이전트 생성·키 발급·철회·관계 변경·메시지 발송은 하지 않았다.

`home-widths.json`은 실제 인증 화면의 320·390·1440px 측정이다. 모두 main 가로 넘침과 44px 미만 조작 영역이 0건이었다. `home-1440.png`는 실제 계정 화면의 대표 캡처이며 이메일은 제품에서 마스킹됐다. 식별자는 공개 식별자다.

키보드에서 본문 건너뛰기를 Enter로 실행해 #main으로 이동했다. 다음 Tab은 `새 에이전트 만들기`로 이동했고 solid 3px outline이 표시됐다. `keyboard-skip.json`과 `keyboard-tab.json`에 관찰값을 보존했다. 모든 폼의 키보드 조작을 반복하지 않았다.

운영 배포 SHA·DB/state/Google/Tunnel·공유 서비스 보존은 OPS의 별도 전후 증거가 정본이다. 이 보고서는 브라우저 관찰만 기록한다. 제품 고정 후보의 기능 QA와 시각 검수는 FIX-TESTER/FIX-VISUAL 보고서를 재사용한다. 스크린리더·200% 확대·다른 브라우저·장시간 수신은 미검증이다.
