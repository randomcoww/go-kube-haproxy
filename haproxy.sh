#!/usr/bin/env bash
# echo -en "$CONFIG" > /etc/haproxy/haproxy.cfg

## start
rm -f $HAPROXY_PID_PATH
exec haproxy "$@" -V \
  -f $HAPROXY_CONFIG_PATH \
	-p $HAPROXY_PID_PATH
