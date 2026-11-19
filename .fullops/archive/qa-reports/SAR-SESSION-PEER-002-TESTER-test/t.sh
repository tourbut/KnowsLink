#!/bin/bash
# t.sh send|key|read A|B [text]
export MSYS_NO_PATHCONV=1
case "$2" in A) H=term_7ffaa33e-81e7-4f07-b2b8-0d982ae72c67;; B) H=term_d370f451-f8f0-4bc1-9827-7bc46e3a7312;; *) H=$2;; esac
case "$1" in
 send) orca terminal send --terminal $H --text "$3" --enter >/dev/null;;
 key) orca terminal send --terminal $H --text "$3" >/dev/null;;
 read) orca terminal read --terminal $H --screen | tail -${3:-30};;
esac
