#!/bin/bash
set -eu

v() {
	printf '%q ' "$@"
	echo
	"$@"
}

echo '# Checking out code...' 
v git config --global advice.detachedHead false
v mkdir -p /build
v cd /build
printf '%q %q https://"${KALE_CREDS}"@github.com/%q.git %q\n' git clone "$KALE_REPO" project
git clone "https://${KALE_CREDS}@github.com/${KALE_REPO}.git" project
v cd project
if [[ "$KALE_REF" == refs/heads/* ]] && v git checkout "${KALE_REF#refs/heads/}"; then
	v git reset --hard "$KALE_REV"
else
	v git checkout "$KALE_REV"
fi

echo
echo '# Performing preflight check...'
if ! test -f Dockerfile; then
    printf 'ERROR: Git commit %s does not contain a file "Dockerfile"\n' "$KALE_REV" >&2
    exit 1
fi

echo
echo '# Building...'
v exec kaniko-executor --context="dir://${PWD}" "$@"
