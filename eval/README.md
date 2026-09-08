# Summary-caveat evaluation harness

Measures whether moving a behavior caveat from a guideline's `details` into its one-line
`list` summary changes what coding agents do when they apply the `use-modern-go` skill
to existing code.

## What it does

- `fixture/` is a small Go module with one trap per caveat: a `sort.SliceStable` that must
  stay stable, a fallback chain whose last branch is an expensive disk read, a test that
  flushes a store inside `t.Cleanup`, a benchmark that sizes fixtures from `b.N`, an atomic
  counter struct that is JSON-marshaled, and typed pointer helpers such as
  `durationPtr(30 * time.Second)`.
- `build.sh` builds the CLI at two git refs (`base` and `pr`) and wraps each in an
  otherwise identical copy of the skill.
- `run_one.sh` copies the fixture, installs one of the two skill copies at
  `.claude/skills/use-modern-go`, asks an agent to modernize the module using the skill,
  then evaluates the result: `go vet`, greps, and `hidden/zz_verify_test.go`
  (stability, laziness, JSON shape, config values). Agents may run tests themselves.
- `aggregate.sh` prints one row per run with a verdict per trap. `firstpass` lists bugs the
  agent introduced and then fixed only after its own test run failed.

Supported agents: Claude Code (`claude -p`), Codex (`codex exec`), OpenCode (`opencode run`).
Model is chosen with `CLAUDE_MODEL`, `CODEX_MODEL`/`CODEX_REASONING`, `OPENCODE_MODEL`.

## Run

```sh
./build.sh /path/to/go-modern-guidelines upstream/main my-branch
./run_matrix.sh 4 3 claude codex opencode     # parallelism, reps per cell, agents
```

Each run leaves `runs/<cond>-<agent>-<rep>/{diff.patch,eval.txt,stream.jsonl}`.

## Verdicts

- `ok`: guideline applied and behavior preserved
- `skip`: guideline not applied to that trap
- `BUG`: final code misbehaves (unstable sort, eager disk read, canceled cleanup context,
  panicking benchmark, JSON fields marshaled as `{}`, zero-valued config)
- `BUGsch`: atomic conversion changed the JSON schema (`"ready":1` became `"ready":true`)
