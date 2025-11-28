#!/usr/bin/env bash

v=$(cat VERSION.txt)
$(which go) build -ldflags="-X 'github.com/ericpfisher/tmplhate/cmd.Version=${VERSION:-$v}'"
