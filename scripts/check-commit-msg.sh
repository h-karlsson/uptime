#!/bin/sh

subject=$(head -n 1 "$1")

case "$subject" in
  "Merge "* | "Revert "* | "fixup! "* | "squash! "*) exit 0 ;;
esac

types='feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert'
if printf '%s\n' "$subject" | grep -Eq "^($types)(\([^)]+\))?!?: .+"; then
  exit 0
fi

echo "Invalid commit message: \"$subject\"" >&2
echo "Expected \"<type>(<scope>)?: <description>\", e.g. \"feat: add status page\"" >&2
echo "Types: $(echo "$types" | tr "|" " ")" >&2
exit 1
