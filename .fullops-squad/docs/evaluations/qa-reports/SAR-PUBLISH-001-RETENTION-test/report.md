---
title: 운영 백업 14일 자동 정리 예약과 14일 초과분 삭제 결과
status: draft
updated: 2026-10-22
owner: ops
tasks: [SAR-PUBLISH-001-RETENTION]
summary: "D021에 따라 운영 knowslink-state/backups의 14일 초과 백업 9개를 지우고 매일 03:17 KST 자동 정리 cron을 예약한 대상 목록·실행·확인 결과를 기록한다."
---

# SAR-PUBLISH-001-RETENTION 결과

운영 서버에 백업 14일 보존 정리를 만들고 14일 초과 백업 9개를 지웠다. 남은 백업은 11개이고 14일 초과분은 0이다. 매일 03:17 KST에 정리가 돌도록 cron을 예약했다. 사용자 결정 D021(2026-10-22) 범위 안이다. 실패·rollback은 없다.

## 범위·기준

- Task `task_48ab82885d36`, Dispatch `ctx_bff9096a83f7`. 기준 main `d5695ff1`. 하위 에이전트는 쓰지 않았다.
- 바꾼 것: 운영 `knowslink-state/backups/`의 14일 초과 `*.dump` 9개 삭제, 서버 `~/.local/bin/knowslink-prune-backups.sh` 설치, 서버 사용자 crontab 1줄 추가.
- 바꾸지 않은 것: 백업 디렉터리 밖 파일, DB·volume, 컨테이너, Tunnel·ingress, 운영 env, Cloud Console, 다른 서비스.
- 시각은 서버 기준이다(KST = UTC+9). 실행 시각 2026-10-21T23:47:03Z(08:47:03 KST 2026-10-22).

## 정리 방식

- 기존 운영 방식에 맞춰 사용자 cron을 썼다. 서버에는 이미 사용자 crontab(`myknows-pull.sh` 매시)이 있고 `cron`이 active다. systemd는 `sudo` 암호가 필요해 쓰지 않았다.
- 스크립트: repo [`deploy/knowslink/prune-backups.sh`](../../../../deploy/knowslink/prune-backups.sh). 서버 사본은 같은 파일이다(1036bytes, SHA256 앞 16자 `a3c3f25cf83e5dc2`, 모드 755).
- 규칙: 백업 디렉터리 바로 아래의 일반 파일 `*.dump` 중 수정 시각이 14일(20160분)을 넘은 것만 지운다. 하위 폴더·심볼릭 링크·`*.dump`가 아닌 파일·디렉터리 밖은 건드리지 않는다. 가장 최근 `*.dump` 1개는 14일을 넘어도 남긴다(새 백업이 멈춰도 마지막 복구점을 지우지 않기 위해서다). 로그는 삭제·보존 개수 한 줄뿐이다.
- 예약: `17 3 * * * ~/.local/bin/knowslink-prune-backups.sh ~/deploy/knowslink-state/backups >> ~/deploy/knowslink-state/backup-prune.log 2>&1`. 로그 파일은 `knowslink-state/backup-prune.log`(0600)다. 기존 crontab은 `knowslink-state/crontab.pre-retention001`(0600, 2줄)에 보관했다.
- 기준을 `find -mtime +14`가 아니라 정확한 14일(`-mmin +20160`)로 했다. `-mtime +14`는 15일 이상만 지워 문안 ‘최대 14일 보관’과 맞지 않는다. 그래서 지시서의 ‘6개’ 대신 14.0~14.7일 지난 2026-10-07 백업 3개가 더 들어가 9개를 지웠다. 모두 14일 초과다.

## 삭제 대상 목록 (지우기 전, 이름·날짜·크기만)

| 파일 | 수정일 | bytes |
|---|---|---|
| `28bd1bb-20261004T044744Z.dump` | 2026-10-04 | 13920 |
| `66d1057-20261003T104032Z.dump` | 2026-10-03 | 6538 |
| `8a7ad36-20261003T105604Z.dump` | 2026-10-03 | 6968 |
| `8a7ad36-20261003T105623Z.dump` | 2026-10-03 | 6969 |
| `8a7ad36-20261003T105630Z.dump` | 2026-10-03 | 6968 |
| `f824015-20261003T110011Z.dump` | 2026-10-03 | 6969 |
| `0911c2c-20261007T130747Z.dump` | 2026-10-07 | 7873 |
| `0911c2c-20261007T140237Z.dump` | 2026-10-07 | 7875 |
| `d08903a-20261007T145512Z.dump` | 2026-10-07 | 8046 |

삭제 전 파일은 20개였다. 가장 최근은 `4787e85-pre-publish001-20261021T224540Z.dump`(rollback용 SAR-PUBLISH-001-OPS 백업, 2026-10-22)이며 남겼다.

## 실행 결과

| 확인 | 결과 |
|---|---|
| 첫 실행 로그 | `2026-10-21T23:47:03Z deleted=9 kept=11` (종료 0) |
| 남은 백업 수 | 11 (모두 0600) |
| 가장 오래된 날짜 | 2026-10-10 (`04a65b20-20261009T171623Z.dump`, 약12일) |
| 14일 초과 | 0 |
| rollback 백업 | `4787e85-pre-publish001-20261021T224540Z.dump` 있음 |
| crontab | 정리 줄 1개 추가, 기존 `myknows-pull.sh` 줄·주석 줄 유지 |
| cron 데몬 | `active` |
| 다음 실행 | 2026-10-23 03:17 KST(2026-10-22 18:17Z) |

다음 실행 시각은 crontab 줄과 cron 데몬 active로 확인했다. 아직 예약 실행이 한 번도 돈 것은 아니다. 2026-10-23 이후 `backup-prune.log`에 줄이 생기는지 확인하는 명령은 D12에 적었다.

## 로컬 검증 (lite)

임시 폴더 fixture에서 실행했다(종료코드 모두 0).

- 20일·15일·13일·1일 된 `*.dump`와 30일 된 `notes.txt`, 디렉터리 밖 30일 된 `*.dump`: 20일·15일 2개만 지워졌고 `notes.txt`와 밖의 파일은 그대로다. 같은 명령을 다시 실행하면 `deleted=0`이다.
- `--list`는 지우지 않고 대상만 보인다. 끝에 `/`가 붙은 디렉터리 인자도 된다.
- 모두 14일 초과(40일·20일)인 경우 가장 최근 1개(20일)가 남는다. 빈 디렉터리는 `deleted=0 kept=0`, 종료 0이다.
- 서버에서 `bash -n` 통과 후 실행했다.

## 한계·미실행

- 예약된 cron의 첫 자동 실행은 아직 일어나지 않았다.
- 서버 `/etc` 설정·systemd·sudo는 쓰지 않았다. 사용자 crontab이라 서버 재설치나 사용자 변경 시 사라질 수 있다.
- 백업 안의 개인정보 삭제 보장은 이 정리의 범위가 아니다. Postgres WAL archive는 꺼져 있다(직전 확인값 `archive_mode=off`, 이번에 다시 확인하지 않았다).
- rollback용 `4787e85-pre-publish001` 백업도 2026-11-05 이후에는 가장 최근 1개가 아니게 되면 정리 대상이 된다. 지시서는 이번 시점의 보존만 요구했다. 더 오래 둘 이유가 생기면 파일을 다른 곳에 복사해야 한다.
- 정책 문안 변경과 배포는 이 과제 범위가 아니다(dev SAR-PUBLISH-001-P03, 배포는 coor).

## 문서 HEAD의 FullOps 검사

완료 보고의 검증 항목에 적는다(아래 인박스 완료 보고와 같다).
