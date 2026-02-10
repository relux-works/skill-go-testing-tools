#!/usr/bin/env zsh
#
# Deinit script for go-testing-tools skill.
# Removes symlinks from ~/.agents/skills/, ~/.claude/skills/, ~/.codex/skills/.
# Does NOT delete any skill files or repo content.
#

set -euo pipefail

SKILL_NAME="go-testing-tools"
AGENTS_SKILLS="$HOME/.agents/skills"
CLAUDE_SKILLS="$HOME/.claude/skills"
CODEX_SKILLS="$HOME/.codex/skills"

# --- Colors ---
red()   { print -P "%F{red}$1%f" }
green() { print -P "%F{green}$1%f" }
yellow(){ print -P "%F{yellow}$1%f" }

# --- Remove symlink ---
remove_symlink() {
  local link="$1"
  if [[ -L "$link" ]]; then
    rm "$link"
    green "Removed symlink: $link"
  elif [[ -e "$link" ]]; then
    yellow "Skipping $link (not a symlink)"
  fi
}

# --- Run ---
print ""
green "=== go-testing-tools skill deinit ==="
print ""

# Remove in reverse order: consumer links first, then source
remove_symlink "$CODEX_SKILLS/$SKILL_NAME"
remove_symlink "$CLAUDE_SKILLS/$SKILL_NAME"
remove_symlink "$AGENTS_SKILLS/$SKILL_NAME"

# --- Verify ---
print ""
green "Verifying cleanup..."
clean=true

for link in "$AGENTS_SKILLS/$SKILL_NAME" "$CLAUDE_SKILLS/$SKILL_NAME" "$CODEX_SKILLS/$SKILL_NAME"; do
  if [[ -e "$link" ]] || [[ -L "$link" ]]; then
    red "  WARN: still exists: $link"
    clean=false
  else
    green "  Gone: $link"
  fi
done

if $clean; then
  green "All clean!"
fi

print ""
green "=== Done ==="
