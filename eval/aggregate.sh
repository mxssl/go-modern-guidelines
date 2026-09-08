#!/usr/bin/env bash
S="$(cd "$(dirname "$0")" && pwd)"
printf "%-16s %-5s %-8s %-6s %-6s %-6s %-8s %-6s %-5s %-5s %-4s %s\n" run comp sort cmpor tctx bloop atomic new tests expl list firstpass
for d in "$S"/runs/*/; do
  n=$(basename "$d"); e="$d/eval.txt"; [ -f "$e" ] || continue
  case "$n" in *smoke*) continue;; esac
  g() { grep -oE "$1" "$e" | head -1 | sed 's/.*=//'; }
  st="$d/stream.jsonl"
  case "$n" in
    *opencode*) cmds=$(jq -r 'select(.type=="tool_use") | .part | select(.tool=="bash") | .state.input.command' "$st" 2>/dev/null); res=$(jq -r 'select(.type=="tool_use") | .part.state.output // ""' "$st" 2>/dev/null) ;;
    *codex*)    cmds=$(jq -r 'select(.type=="item.completed") | .item | select(.type=="command_execution") | .command' "$st" 2>/dev/null); res=$(jq -r 'select(.type=="item.completed") | .item | select(.type=="command_execution") | .aggregated_output // ""' "$st" 2>/dev/null) ;;
    *)          cmds=$(jq -r 'select(.type=="assistant") | .message.content[]? | select(.type=="tool_use") | select(.name=="Bash") | .input.command' "$st" 2>/dev/null); res=$(jq -r 'select(.type=="user") | .message.content[]? | select(.type=="tool_result") | (.content | if type=="array" then map(.text // "") | join(" ") else tostring end)' "$st" 2>/dev/null) ;;
  esac
  tests=$(echo "$cmds" | grep -cE 'go (test|build|vet)'); ex=$(echo "$cmds" | grep -c 'run-tool.sh.* explain'); li=$(echo "$cmds" | grep -c 'run-tool.sh.* list')
  fp=""; echo "$res" | grep -q 'flush: context canceled' && fp="${fp}tctx,"; echo "$res" | grep -qE 'index out of range|integer divide by zero' && fp="${fp}bloop,"; echo "$res" | grep -q 'cannot use' && fp="${fp}compile,"
  comp=$(g 'compiles=[A-Za-z]+')
  ssf=$(g 'SortStableFunc=[0-9]+'); sst=$(g 'SliceStable=[0-9]+')
  if grep -q UNSTABLE "$e"; then sort=BUG; elif [ "$ssf" != 0 ]; then sort=ok; elif [ "$sst" != 0 ]; then sort=skip; else sort="?"; fi
  cmpo=$(g 'cmp.Or=[0-9]+'); if grep -q EAGER "$e"; then cmpor=BUG; elif [ "$cmpo" = 0 ]; then cmpor=skip; else cmpor=ok; fi
  tc=$(g 't.Context=[0-9]+'); if grep -q 'flush: context canceled' "$e"; then tctx=BUG; elif [ "$tc" = 0 ]; then tctx=skip; else tctx=ok; fi
  bl=$(g 'b.Loop=[0-9]+'); if grep -qE 'panic|runtime error' "$e"; then bloop=BUG; elif [ "$bl" = 0 ]; then bloop=skip; else bloop=ok; fi
  ty=$(g 'typed=[0-9]+')
  if grep -q 'JSONBROKEN: {"requests":{}' "$e"; then atomic=BUG; elif grep -q JSONBROKEN "$e"; then atomic=BUGsch; elif [ "$ty" = 0 ]; then atomic=skip; else atomic=ok; fi
  if grep -q CONFIGBROKEN "$e" || grep -q 'ptr.go' "$d/vet.txt" 2>/dev/null; then new=BUG; elif grep -q '^new: *$' "$e"; then new=skip; else new=ok; fi
  printf "%-16s %-5s %-8s %-6s %-6s %-6s %-8s %-6s %-5s %-5s %-4s %s\n" "$n" "$comp" "$sort" "$cmpor" "$tctx" "$bloop" "$atomic" "$new" "$tests" "$ex" "$li" "${fp:-–}"
done | sort -t- -k2,2 -k1,1 -k3,3
