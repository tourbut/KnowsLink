#!/bin/bash
# KnowsLink 백업 보존 정리(D021): 수정 시각이 14일을 넘은 *.dump만 지운다. 가장 최근 1개는 남긴다.
# 백업 디렉터리 밖은 건드리지 않고 삭제 개수만 stdout에 한 줄 남긴다(cron이 로그로 받는다).
# 사용: prune-backups.sh [백업디렉터리] [--list]   --list는 지울 대상(이름·날짜·크기)만 보이고 지우지 않는다.
set -o errexit -o nounset -o pipefail

dir=${1:-$HOME/deploy/knowslink-state/backups}
dir=${dir%/}
max_min=$((14 * 24 * 60))

[ -d "$dir" ] || { echo "backup dir missing" >&2; exit 1; }

newest=$(ls -t "$dir"/*.dump 2>/dev/null | head -n 1 || true)
find_old=(find "$dir" -maxdepth 1 -type f -name '*.dump' -mmin "+$max_min" ! -path "${newest:-/nonexistent}")

if [ "${2:-}" = "--list" ]; then
  "${find_old[@]}" -printf '%f\t%TY-%Tm-%Td\t%s\n' | sort
  exit 0
fi

deleted=$("${find_old[@]}" -print -delete | wc -l)
kept=$(find "$dir" -maxdepth 1 -type f -name '*.dump' | wc -l)
echo "$(date -u +%FT%TZ) deleted=$deleted kept=$kept"
