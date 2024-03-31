#!/usr/bin/env bash

app_version=$(cat VERSION.txt)
cmd="buildx build"
OLD_IFS=$IFS
IFS=$'\n'
params=("--platform ${PLATFORM:-linux/amd64}"
"--tag tmplhate:${VERSION:-$app_version}"
"--file Dockerfile"
".")
IFS=$OLD_IFS

build () {
  ${RUNNER:-docker} $cmd ${params[@]}
  ${RUNNER:-docker} tag tmplhate:${VERSION:-$app_version} tmplhate:latest
}

build
