<!-- fullops:start -->
## FullOps

이 레포는 `.fullops/fullops.toml`이 있을 때만 FullOps 하네스를 쓴다. 규약의 지도는 `.fullops/guide/FULLOPS.md`다.

- 세션을 시작하면 `fullops status`를 먼저 실행한다. 역할, 현재 인박스, 병렬 방침, lead의 할 일이 거기 나온다. Claude Code가 아닌 호스트에서는 고정 위치 실행본을 전체 경로로 부른다: `~/.fullops/bin/fullops`(Windows는 `%USERPROFILE%\.fullops\bin\fullops.exe`).
- 역할은 체크아웃한 브랜치로 정해진다. 기본 브랜치(dev 모드는 주 담당 브랜치)의 세션이 lead이고, 역할 브랜치 `fullops/<역할>`의 세션이 그 역할의 worker다. Orca 배정(preamble)을 받은 세션은 그 배정의 worker다.
- 현재 과제의 정본은 역할 인박스 `.fullops/inbox/to_<역할>.md` 하나다. 인박스는 `fullops task new`로만 만들고, 다음 과제는 `.fullops/PLANS.md`의 대기 행으로 둔다. 과제 이름으로 된 지시서 파일은 만들지 않는다.
- 코드를 바꾸면 커밋하고 `fullops lint --from <기준>`을 통과시킨다. 통과 기록 없이 턴을 끝내면 종료가 막힌다.
- 맡은 과제가 끝나면 `fullops done`을 실행하고 출력된 본문으로 Orca preamble의 `worker_done`을 한 번 보낸다.
- 배정은 `fullops dispatch`, 완료 보고의 병합은 `fullops merge`로 한다. 날것 `worker-start`와 터미널에 지시서를 붙여 넣는 배정은 쓰지 않는다.
- lead가 사용자에게 물을 때는 질문 카드로만 묻고(Claude Code는 `AskUserQuestion`) 추천안을 1번 선택지에 두어 label에 `(추천)`을 붙인다. 카드는 열린 결정으로 기록된다. 답을 받으면 `fullops decision list --open`의 ID로 `fullops decision answer <ID> --answer "<답>"`을 실행한다. 카드 도구가 없는 호스트는 묻기 전에 `fullops decision add`로 적는다.
- 막힌 이유는 `fullops why`로 본다.
<!-- fullops:end -->
