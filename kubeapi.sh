#!/usr/bin/env bash
echo -en "$CONFIG" > /template

## start
exec /go-kube-haproxy "$@" \
  -template /template \
  -output $HAPROXY_CONFIG_PATH \
  -reloadcmd "kill -s HUP $(pidof haproxy-systemd-wrapper)"
