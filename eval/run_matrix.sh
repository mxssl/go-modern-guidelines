#!/usr/bin/env bash
# usage: run_matrix.sh <parallelism> <reps> <agent>... ; e.g. run_matrix.sh 4 3 claude codex opencode
set -u
S="$(cd "$(dirname "$0")" && pwd)"
par=$1; reps=$2; shift 2
for a in "$@"; do for c in base pr; do for r in $(seq 1 "$reps"); do echo "$a $c $r"; done; done; done \
  | xargs -P "$par" -L 1 bash -c '"'"$S"'/run_one.sh" "$0" "$1" "$2" > /dev/null 2>&1; echo "done $0 $1 $2"'
"$S/aggregate.sh"
