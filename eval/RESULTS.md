# Summary-caveat experiment (2026-09-08)

30 runs modernizing scratchpad/fixture with the use-modern-go skill. base = upstream/main CLI, pr = PR29+PR13 summaries. 3 reps per cell.
Agents: Claude Code (haiku 4.5, sonnet 5, opus 5), Codex 0.153.4 (gpt-6-astra, reasoning medium), OpenCode 1.18.29 (opencode/deepseek-v4-flash).
Columns: comp = go vet passes; per-trap verdict BUG/BUGsch(schema change)/ok/skip; tests = go test/build/vet commands run; expl/list = wrapper calls; firstpass = bug introduced then fixed after agent's own test run.

```
run              comp  sort     cmpor  tctx   bloop  atomic   new    tests expl  list firstpass
base-codex-1     yes   ok       skip   ok     ok     skip     ok     1     1     1    –
base-codex-2     yes   ok       skip   ok     skip   skip     ok     1     1     1    –
base-codex-3     yes   ok       skip   ok     ok     skip     ok     1     1     1    –
pr-codex-1       yes   ok       skip   ok     ok     skip     ok     3     1     1    –
pr-codex-2       yes   ok       skip   ok     ok     skip     ok     1     1     1    –
pr-codex-3       yes   ok       skip   ok     ok     skip     ok     1     1     1    –
base-haiku-1     NO    skip     skip   ok     skip   BUG      skip   6     1     1    tctx,bloop,
base-haiku-2     yes   skip     skip   skip   ok     ok       skip   4     0     1    –
base-haiku-3     yes   BUG      skip   skip   skip   BUGsch   ok     4     2     1    bloop,
pr-haiku-1       yes   ok       skip   ok     ok     ok       BUG    3     1     1    –
pr-haiku-2       NO    ok       skip   skip   skip   BUG      ok     2     0     1    –
pr-haiku-3       yes   ok       skip   ok     ok     BUGsch   ok     5     2     1    –
base-opencode-1  yes   ok       skip   skip   ok     BUGsch   ok     2     1     1    tctx,
base-opencode-2  yes   ok       skip   skip   skip   BUGsch   ok     2     1     1    tctx,
base-opencode-3  yes   ok       skip   skip   ok     skip     ok     3     1     1    tctx,
pr-opencode-1    yes   ok       skip   skip   skip   BUG      ok     1     1     1    –
pr-opencode-2    yes   ok       skip   ok     ok     BUG      ok     2     1     1    –
pr-opencode-3    yes   ok       skip   ok     ok     skip     ok     4     1     1    –
base-opus-1      yes   ok       ok     ok     ok     ok       ok     2     2     1    –
base-opus-2      yes   ok       ok     ok     ok     ok       ok     2     2     1    –
base-opus-3      yes   ok       ok     ok     ok     ok       ok     2     2     1    –
pr-opus-1        yes   ok       skip   ok     ok     ok       ok     2     2     1    –
pr-opus-2        yes   ok       ok     ok     ok     ok       ok     2     2     1    –
pr-opus-3        yes   ok       ok     ok     ok     ok       ok     2     2     1    –
base-sonnet-1    yes   skip     ok     ok     ok     ok       ok     3     1     1    tctx,
base-sonnet-2    yes   ok       skip   skip   ok     skip     ok     2     1     1    tctx,
base-sonnet-3    yes   skip     skip   skip   ok     skip     ok     1     3     2    –
pr-sonnet-1      yes   ok       skip   ok     ok     skip     ok     1     1     1    –
pr-sonnet-2      yes   ok       skip   ok     ok     skip     ok     1     1     1    –
pr-sonnet-3      yes   ok       skip   ok     ok     skip     ok     1     1     1    –
```

## Per-caveat totals (15 base runs vs 15 pr runs)

| caveat | base final bugs | base first-pass bugs | base skipped | pr final bugs | pr first-pass bugs | pr skipped |
|---|---|---|---|---|---|---|
| slices_sort_func | 1 | 0 | 4 | 0 | 0 | 0 |
| testing_t_context | 0 | 6 | 7 | 0 | 0 | 2 |
| testing_b_loop | 0 | 2 | 4 | 0 | 0 | 2 |
| cmp_or | 0 | 0 | 11 | 0 | 0 | 13 |
| atomic_types | 4 | 0 | 6 | 4 | 0 | 7 |
| new_expression (PR13) | 0 | 0 | 2 | 1 | 0 | 0 |

Run dirs: runs/<name>/{diff.patch,eval.txt,stream.jsonl,last.md}
