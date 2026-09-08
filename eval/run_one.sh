#!/usr/bin/env bash
# Run one agent once against a fresh copy of the fixture and evaluate the result.
# usage: run_one.sh <agent: claude|codex|opencode> <cond: base|pr> <rep>
# Models: CLAUDE_MODEL (default sonnet), CODEX_MODEL (default gpt-6-astra), OPENCODE_MODEL (default opencode/deepseek-v4-flash)
set -u
S="$(cd "$(dirname "$0")" && pwd)"
agent=$1; cond=$2; rep=$3
name="${cond}-${agent}-${rep}"
R="$S/runs/$name"
rm -rf "$R"; mkdir -p "$R"
cp -R "$S/fixture/." "$R/proj"
mkdir -p "$R/proj/.claude/skills/use-modern-go"
cp -R "$S/skill-$cond/." "$R/proj/.claude/skills/use-modern-go/"
cd "$R/proj" && git init -q && git add -A && git -c user.email=eval@example.invalid -c user.name=eval commit -qm init

PROMPT='Modernize the Go code in this module to current Go idioms. Use the use-modern-go skill in .claude/skills/use-modern-go (read its SKILL.md) as the source of truth: follow its workflow (call its list subcommand via the wrapper script, then apply the returned guidelines). Edit the files in place. Do not add third-party dependencies. When finished, briefly list which guidelines you applied and which you deliberately skipped and why.'

start=$(date +%s)
case "$agent" in
  claude)
    claude -p "$PROMPT" --model "${CLAUDE_MODEL:-sonnet}" \
      --dangerously-skip-permissions --setting-sources project --no-session-persistence \
      --max-turns 80 --max-budget-usd 6 --output-format stream-json --verbose \
      < /dev/null > "$R/stream.jsonl" 2> "$R/stderr.txt" ;;
  codex)
    codex exec -m "${CODEX_MODEL:-gpt-6-astra}" -c model_reasoning_effort="${CODEX_REASONING:-medium}" \
      --dangerously-bypass-approvals-and-sandbox --skip-git-repo-check --ephemeral \
      -C "$R/proj" --json -o "$R/last.md" "$PROMPT" \
      < /dev/null > "$R/stream.jsonl" 2> "$R/stderr.txt" ;;
  opencode)
    opencode run --dir "$R/proj" -m "${OPENCODE_MODEL:-opencode/deepseek-v4-flash}" --auto --pure --format json "$PROMPT" \
      < /dev/null > "$R/stream.jsonl" 2> "$R/stderr.txt" ;;
  *) echo "unknown agent $agent" >&2; exit 2 ;;
esac
echo "exit=$? secs=$(( $(date +%s) - start ))" > "$R/meta.txt"

cd "$R/proj"
git diff > "$R/diff.patch"
{
  echo "## $name"
  if go vet ./... >"$R/vet.txt" 2>&1; then echo "compiles=yes"; else echo "compiles=NO"; fi
  echo "sort: SortFunc=$(grep -c 'slices.SortFunc' sortstable.go) SortStableFunc=$(grep -c 'slices.SortStableFunc' sortstable.go) SliceStable=$(grep -c 'sort.SliceStable' sortstable.go)"
  echo "cmpor: cmp.Or=$(grep -c 'cmp.Or' fallback.go)"
  echo "tctx: t.Context=$(grep -c 't.Context()' store_test.go) Background=$(grep -c 'context.Background' store_test.go)"
  echo "bloop: b.Loop=$(grep -c 'b.Loop()' sum_bench_test.go) bN=$(grep -c 'b.N' sum_bench_test.go)"
  echo "atomic: typed=$(grep -cE 'atomic\.(Int64|Int32|Bool)' stats.go) untyped=$(grep -cE 'atomic\.(Add|Load|Store)Int' stats.go)"
  echo "new: $(grep -oE 'new\([^)]*\)' ptr.go | sort | uniq -c | tr '\n' ';')"
  cp "$S/hidden/zz_verify_test.go" .
  echo "--- hidden tests:"; go test ./... 2>&1 | grep -E 'FAIL|UNSTABLE|EAGER|JSONBROKEN|CONFIGBROKEN|flush|^ok|cannot|undefined|panic' | head -12
  echo "--- bench:"; go test -run xxx -bench BenchmarkSum -benchtime=100x 2>&1 | grep -E 'panic|Benchmark|FAIL|ok|runtime error' | head -4
  rm -f zz_verify_test.go
} > "$R/eval.txt" 2>&1
cat "$R/eval.txt"
