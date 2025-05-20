#!/bin/bash

# DO NOT EDIT
# This file is taken from go-perun, any changes should be done there first.
# https://github.com/perun-network/go-perun

cn="$(dirname $(readlink -f $0))/copyright-header"
n=$(wc -l $cn | cut -d ' ' -f 1)

function check_cr() {
  diff $cn <(head -n${n} $1 | sed -e 's/20\(19\|2[0-9]\)/20XX/') -q > /dev/null
  if [ $? -ne 0 ]; then
    echo $1
    diff $cn <(head -n${n} $1)
  fi
}

if [ $# -ne 0 ]; then
  code=0
  for f in "$@"; do
    check_cr $f
    [ $? -ne 0 ] && code=1
  done
  exit $code
fi

find contracts/ migrations/ src/ -type f -exec $0 {} +
