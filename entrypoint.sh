#!/bin/sh
echo -en "$CONFIG" > /template

## start
exec /goapp "$@" \
  -template /template
