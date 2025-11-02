#!/usr/bin/env bash
set -euo pipefail
root="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
cd "$root"

normalize_go_mod() {
  [[ -f go.mod ]] || return 0
  # Strip UTF-8 BOM on first line (if present)
  awk 'NR==1{sub(/^\xef\xbb\xbf/,"")} {print}' go.mod > go.mod.tmp && mv go.mod.tmp go.mod
  # Convert CRLF -> LF if detected
  if file go.mod 2>/dev/null | grep -qi 'CRLF'; then
    awk '{ sub(/\r$/,""); print }' go.mod > go.mod.tmp && mv go.mod.tmp go.mod
  fi
}

guess_module_from_remote() {
  local remote mod_guess=""
  remote="$(git remote get-url origin 2>/dev/null || true)"
  if [[ "$remote" =~ ^git@github\.com:([^[:space:]]+?)(\.git)?$ ]]; then
    mod_guess="github.com/${BASH_REMATCH[1]}"
  elif [[ "$remote" =~ ^https?://github\.com/([^[:space:]]+?)(\.git)?$ ]]; then
    mod_guess="github.com/${BASH_REMATCH[1]}"
  fi
  if [[ -z "$mod_guess" && -f cmd/mldsa/main.go ]]; then
    local imp
    imp="$(grep -Eo 'github\.com/[^"]+/code/clean' cmd/mldsa/main.go | head -n1 || true)"
    [[ -n "$imp" ]] && mod_guess="${imp%/code/clean}"
  fi
  echo "${mod_guess:-github.com/codethor0/ml-dsa-debug-whitepaper}"
}

read_module_line() {
  awk '/^module[[:space:]]+/ {print $2; exit}' go.mod 2>/dev/null || true
}

write_or_fix_module() {
  local cur="$1" want="${2%.git}"
  if [[ -z "$cur" ]]; then
    if [[ ! -f go.mod || ! -s go.mod ]]; then
      printf 'module %s\n\ngo 1.22\n' "$want" > go.mod
      echo "created go.mod with module $want"
    else
      { echo "module $want"; cat go.mod; } > go.mod.tmp && mv go.mod.tmp go.mod
      echo "inserted module line: $want"
    fi
  elif [[ "$cur" != "$want" ]]; then
    awk -v m="$want" '{if(!done && $1=="module"){print "module " m; done=1; next} print}' go.mod > go.mod.tmp && mv go.mod.tmp go.mod
    echo "rewrote module: $cur -> $want"
  else
    echo "module OK: $want"
  fi
}

fix_imports() {
  local mod="$1" want="$1/code/clean"
  if [[ -f cmd/mldsa/main.go ]] && grep -qE 'github\.com/[^"]+/code/clean' cmd/mldsa/main.go; then
    sed -E -i '' "s|github\.com/[^\" ]+/code/clean|$want|g" cmd/mldsa/main.go
    echo "imports synced to: $want"
  fi
}

normalize_go_mod
cur_mod="$(read_module_line || true)"
mod_guess="$(guess_module_from_remote)"
write_or_fix_module "$cur_mod" "$mod_guess"
fix_imports "${mod_guess%.git}"
