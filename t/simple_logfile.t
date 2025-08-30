#!/bin/ksh
#
# simple_logfile.t
# Make sure our code is up-to-date and doesn't have debug things.
# By J. Stuart McMurray
# Created 20250823
# Last Modified 20250830

set -euo pipefail

. t/shmore.subr

tap_plan 15

TMPD=$(mktemp -d)
trap 'rm -rf "$TMPD"; tap_done_testing' EXIT
TD=t/testdata/simple_logfile
LOGFILE="$TMPD/audit.log"
NEWLOGS="$TMPD/audit.log.new"
FILTERS="$TMPD/procrastinationedr.specs"

# Get temporary files ready.
cp "$TD"/* "$TMPD"

go run . \
        -audit-log     "$LOGFILE"      \
        -filters       "$FILTERS"      \
        -max-runtime    1m             \
        -no-timestamps </dev/null 2>&1 |&
PID=$!

# Should get the last line of the original logfile
WANT="2 /lib/systemd/systemd-networkd-wait-online"
read -pr
tap_ok $?               "Read line from original logfile"   "$0" $LINENO
tap_is "$REPLY" "$WANT" "Correct hit from original logfile" "$0" $LINENO

# Should get lines from new entries.
cat <"$NEWLOGS" >>"$LOGFILE"
tap_ok $? "Added new logs to logfile" "$0" $LINENO

I=0
set -A WANTS \
        "6 ps awwwfux" \
        "3 uname -a" \
        '1 touch foo"bar' \
        '1 sh foo"bar' \
        "10 /usr/bin/clear_console -q" \
        "10 bash -c  chmod 0755 bin/vulnerableserver sudo setcap cap_net_bind_service+ep bin/vulnerableserver && unset SSH_CLIENT SSH_CONNEC"

for WANT in "${WANTS[@]}"; do
        : $((I++))
        read -pr
        N="$I/${#WANTS[@]}"
        tap_ok  $?              "Read hit $N"    "$0" $LINENO
        tap_is "$REPLY" "$WANT" "Hit $N correct" "$0" $LINENO
done

# Kill the EDR
kill $PID
wait

# vim: ft=sh
