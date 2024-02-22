#!/usr/bin/env bash

v=$(cat VERSION.txt)
$(which go) build -ldflags="-X 'tmplhate/cmd.Version=${VERSION:-$v}'"
