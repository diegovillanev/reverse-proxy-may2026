#!/bin/sh

# Check for any changes to committed .env* files (excluding `.env`)
for file in $(git diff --cached --name-only | grep '^\.env' | grep -v '^\.env$'); do
    # If the file is already in the repo and changed
    if git ls-files --error-unmatch "$file" >/dev/null 2>&1; then
        echo "❌ Changes to $file are forbidden."
        exit 1
    fi
done
