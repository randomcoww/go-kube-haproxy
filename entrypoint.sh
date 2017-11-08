#!/usr/bin/env bash

echo -en "$TEMPLATE" > /go_template

## start
exec /go-kube-haproxy $@ -template /go_template
