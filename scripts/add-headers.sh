#!/usr/bin/env bash

# © 2025 Thor Thor
# Contact: codethor@gmail.com
# LinkedIn: https://www.linkedin.com/in/thor-thor0
# SPDX-License-Identifier: MIT

set -euo pipefail

NAME="Thor Thor"
EMAIL="codethor@gmail.com"
LINKEDIN="https://www.linkedin.com/in/thor-thor0"
SPDX_ID="MIT"
YEAR="$(date +%Y)"

exts=(
  go rs py sh bash zsh rb yml yaml toml ini cfg nix ps1 js ts jsx tsx c cc cpp h hh hpp
  java kt kts scala cs php css scss html htm xml sql hs lua dart
)

special_names=(
  "Dockerfile"
  "Dockerfile.*"
  "Makefile"
  "makefile"
)

skip_dirs=(
  ".git"
  "node_modules"
  "vendor"
  "target"
  "build"
  "dist"
  ".next"
  ".venv"
  "venv"
  "__pycache__"
  ".tox"
  ".mypy_cache"
  "coverage"
  ".idea"
  ".vscode"
)

is_text() {
  [[ -f "$1" ]] && LC_ALL=C grep -Iq . "$1"
}

skip_file() {
  case "$1" in
    *.min.*|*.map|*.json|*.ipynb) return 0 ;;
    *) return 1 ;;
  esac
}

has_header() {
  grep -qE 'SPDX-License-Identifier:|codethor@gmail.com' "$1"
}

comment_style() {
  local ext="$1"
  case "$ext" in
    py|sh|bash|zsh|rb|yml|yaml|toml|ini|cfg|nix|ps1) echo "hash" ;;
    css|scss) echo "block" ;;
    html|htm|xml) echo "html" ;;
    sql|hs|lua) echo "dash" ;;
    php) echo "php" ;;
    *) echo "slash" ;;
  esac
}

make_header() {
  local style="$1"
  case "$style" in
    hash)
      cat <<EOF
# © ${YEAR} ${NAME}
# Contact: ${EMAIL}
# LinkedIn: ${LINKEDIN}
# SPDX-License-Identifier: ${SPDX_ID}
EOF
      ;;
    slash)
      cat <<EOF
// © ${YEAR} ${NAME}
// Contact: ${EMAIL}
// LinkedIn: ${LINKEDIN}
// SPDX-License-Identifier: ${SPDX_ID}
EOF
      ;;
    block)
      cat <<EOF
/*
 * © ${YEAR} ${NAME}
 * Contact: ${EMAIL}
 * LinkedIn: ${LINKEDIN}
 * SPDX-License-Identifier: ${SPDX_ID}
 */
EOF
      ;;
    dash)
      cat <<EOF
-- © ${YEAR} ${NAME}
-- Contact: ${EMAIL}
-- LinkedIn: ${LINKEDIN}
-- SPDX-License-Identifier: ${SPDX_ID}
EOF
      ;;
    html)
      cat <<EOF
<!--
  © ${YEAR} ${NAME}
  Contact: ${EMAIL}
  LinkedIn: ${LINKEDIN}
  SPDX-License-Identifier: ${SPDX_ID}
-->
EOF
      ;;
    php)
      cat <<EOF
/*
 * © ${YEAR} ${NAME}
 * Contact: ${EMAIL}
 * LinkedIn: ${LINKEDIN}
 * SPDX-License-Identifier: ${SPDX_ID}
 */
EOF
      ;;
    *)
      return 1
      ;;
  esac
}

insert_header() {
  local file="$1"
  local ext="$2"
  local style
  style="$(comment_style "$ext")"
  local header
  header="$(make_header "$style")"

  local tmp
  tmp="$(mktemp)"
  trap 'rm -f "$tmp"' RETURN

  case "$style" in
    hash)
      if head -n1 "$file" | grep -q "^#!"; then
        {
          head -n1 "$file"
          echo ""
          echo "$header"
          tail -n +2 "$file"
        } >"$tmp"
      else
        {
          echo "$header"
          echo ""
          cat "$file"
        } >"$tmp"
      fi
      ;;
    html)
      if head -n1 "$file" | grep -Eqi '^<\?xml|^<!DOCTYPE'; then
        {
          head -n1 "$file"
          echo ""
          echo "$header"
          tail -n +2 "$file"
        } >"$tmp"
      else
        {
          echo "$header"
          echo ""
          cat "$file"
        } >"$tmp"
      fi
      ;;
    php)
      if head -n1 "$file" | grep -q '^<\?php'; then
        {
          head -n1 "$file"
          echo "$header"
          tail -n +2 "$file"
        } >"$tmp"
      else
        {
          echo "$header"
          echo ""
          cat "$file"
        } >"$tmp"
      fi
      ;;
    *)
      {
        echo "$header"
        echo ""
        cat "$file"
      } >"$tmp"
      ;;
  esac

  mv "$tmp" "$file"
}

should_handle_file() {
  local file="$1"
  local base
  base="$(basename "$file")"
  for pattern in "${special_names[@]}"; do
    if [[ "$base" == $pattern ]]; then
      echo "hash"
      return 0
    fi
  done
  local ext="${file##*.}"
  for candidate in "${exts[@]}"; do
    if [[ "$candidate" == "$ext" ]]; then
      echo "$ext"
      return 0
    fi
  done
  return 1
}

main() {
  local find_cmd=(find .)
  for dir in "${skip_dirs[@]}"; do
    find_cmd+=(-path "./${dir}" -prune -o)
  done

  find_cmd+=(-type f \( )
  for ext in "${exts[@]}"; do
    find_cmd+=(-name "*.${ext}" -o)
  done
  for name in "${special_names[@]}"; do
    find_cmd+=(-name "${name}" -o)
  done
  unset 'find_cmd[${#find_cmd[@]}-1]'
  find_cmd+=(\) -print0)

  "${find_cmd[@]}" | while IFS= read -r -d '' file; do
    skip_file "$file" && continue
    is_text "$file" || continue
    has_header "$file" && continue

    local ext
    if ! ext="$(should_handle_file "$file")"; then
      continue
    fi

    if [[ "$ext" == "hash" ]]; then
      ext="sh"
    fi

    insert_header "$file" "$ext"
    echo "Added header: $file"
  done
}

main "$@"

echo "Done."

