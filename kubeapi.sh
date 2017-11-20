#!/usr/bin/env bash
echo -en "$CONFIG" > /template

## start
exec /go-kube-haproxy "$@" \
  -template /template \
  -output $HAPROXY_CONFIG_PATH
  -reloadcmd "haproxy -f $HAPROXY_CONFIG_PATH -p $HAPROXY_PID_PATH -sf $(cat $HAPROXY_PID_PATH)"
