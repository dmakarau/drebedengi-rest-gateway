#!/bin/sh
# PreToolUse hook: block commits, PRs and issues that carry AI attribution.
CMD=$(python3 -c 'import sys,json; print(json.load(sys.stdin).get("tool_input",{}).get("command",""))' 2>/dev/null)

echo "$CMD" | grep -qE 'git commit|gh (pr|issue) (create|edit|comment)' || exit 0

if echo "$CMD" | grep -qiE 'Co-Authored-By: *Claude|Generated with.*Claude|🤖'; then
  echo "Remove Co-Authored-By and AI attribution from the message" >&2
  exit 2
fi
