#!/usr/bin/env bash
echo -en "$CONFIG" > /go_template

## start
exec /go-kube-haproxy "$@" -template /go_template
