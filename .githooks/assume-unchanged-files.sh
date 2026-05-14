#!/bin/sh

# Get repo root path
REPO_ROOT=$(git rev-parse --show-toplevel)

# Loop over all .env* files in the repo root
for file in "$REPO_ROOT"/.env*; do
    [ -f "$file" ] && git update-index --assume-unchanged "$file"
done
