#!/bin/sh

if git diff --cached --name-only | grep -q '^.env$'; then
    echo "❌ .env file is forbidden from the repository."
    exit 1
fi
