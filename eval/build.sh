#!/usr/bin/env bash
# Build the CLI at two git refs and wrap each in a copy of the skill.
# usage: build.sh <repo-path> <base-ref> <pr-ref>
set -eu
S="$(cd "$(dirname "$0")" && pwd)"
repo=$1; base=$2; pr=$3
mkdir -p "$S/bin"
for pair in "base:$base" "pr:$pr"; do
  cond=${pair%%:*}; ref=${pair#*:}
  tmp=$(mktemp -d)
  git -C "$repo" archive "$ref" | tar -x -C "$tmp"
  (cd "$tmp" && go build -o "$S/bin/gmg-$cond" .)
  d="$S/skill-$cond"; rm -rf "$d"; mkdir -p "$d/scripts"
  cp "$tmp/plugin/skills/use-modern-go/SKILL.md" "$d/"
  printf '#!/usr/bin/env sh\nexec "%s/bin/gmg-%s" "$@"\n' "$S" "$cond" > "$d/scripts/run-tool.sh"
  chmod +x "$d/scripts/run-tool.sh"
  rm -rf "$tmp"
  echo "built $cond from $ref"
done
echo; echo "list output differences (go 1.26):"
diff <("$S/bin/gmg-base" list --go-version 1.26) <("$S/bin/gmg-pr" list --go-version 1.26) || true
