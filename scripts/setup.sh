#!/usr/bin/env zsh
#
# Setup script for go-testing-tools skill.
# Runs prerequisites check, then creates symlink chain:
#   ~/.agents/skills/go-testing-tools  -> <this-repo>
#   ~/.claude/skills/go-testing-tools  -> ~/.agents/skills/go-testing-tools
#   ~/.codex/skills/go-testing-tools   -> ~/.agents/skills/go-testing-tools
#

set -euo pipefail

REPO_DIR="$(cd "$(dirname "$0")/.." && pwd)"
SKILL_NAME="go-testing-tools"
AGENTS_SKILLS="$HOME/.agents/skills"
CLAUDE_SKILLS="$HOME/.claude/skills"
CODEX_SKILLS="$HOME/.codex/skills"

# --- Colors ---
red()   { print -P "%F{red}$1%f" }
green() { print -P "%F{green}$1%f" }
yellow(){ print -P "%F{yellow}$1%f" }

# --- 1. Run prerequisites check ---
run_prerequisites() {
  green "Running prerequisites check..."
  if bash "$REPO_DIR/scripts/check-tools.sh"; then
    green "Prerequisites OK"
  else
    red "Prerequisites check failed. Fix issues above before continuing."
    exit 1
  fi
}

# --- 2. Create symlink (replace if exists) ---
create_symlink() {
  local target="$1"
  local link="$2"

  if [[ -L "$link" ]]; then
    local existing
    existing="$(readlink "$link")"
    if [[ "$existing" == "$target" ]]; then
      green "Symlink already correct: $link -> $target"
      return
    fi
    rm "$link"
    yellow "Replaced existing symlink: $link"
  elif [[ -d "$link" ]]; then
    yellow "Skipping $link (real directory exists — remove manually if needed)"
    return
  fi

  mkdir -p "$(dirname "$link")"
  ln -s "$target" "$link"
  green "Created symlink: $link -> $target"
}

# --- 3. Verify ---
verify() {
  local link="$1"
  if [[ -L "$link" ]] && [[ -e "$link" ]]; then
    green "  OK: $link"
  elif [[ -L "$link" ]]; then
    red "  BROKEN: $link (dangling symlink)"
  else
    red "  MISSING: $link"
  fi
}

# --- Run ---
print ""
green "=== go-testing-tools skill setup ==="
print ""

run_prerequisites
print ""

# Symlink chain: repo -> ~/.agents -> ~/.claude + ~/.codex
create_symlink "$REPO_DIR" "$AGENTS_SKILLS/$SKILL_NAME"
create_symlink "$AGENTS_SKILLS/$SKILL_NAME" "$CLAUDE_SKILLS/$SKILL_NAME"
create_symlink "$AGENTS_SKILLS/$SKILL_NAME" "$CODEX_SKILLS/$SKILL_NAME"

# --- Verify ---
print ""
green "Verifying..."
verify "$AGENTS_SKILLS/$SKILL_NAME"
verify "$CLAUDE_SKILLS/$SKILL_NAME"
verify "$CODEX_SKILLS/$SKILL_NAME"

# Check SKILL.md reachability (if it exists yet)
if [[ -f "$CLAUDE_SKILLS/$SKILL_NAME/SKILL.md" ]]; then
  green "  SKILL.md reachable via Claude Code"
elif [[ -f "$REPO_DIR/SKILL.md" ]]; then
  green "  SKILL.md exists in repo (will be reachable after symlinks)"
else
  yellow "  SKILL.md not yet created (expected — skill is being built)"
fi

print ""
green "=== Done ==="
