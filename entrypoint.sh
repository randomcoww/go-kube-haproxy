#!/usr/bin/env bash
echo -en "$CONFIG" > /template

## start
exec /go-kube-haproxy "$@" -template /template
