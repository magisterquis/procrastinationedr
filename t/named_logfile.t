#!/bin/ksh
#
# named_logfile.t
# Make sure we can read a simple auditd log
# By J. Stuart McMurray
# Created 20250901
# Last Modified 20250901

set -euo pipefail

. t/shmore.subr

TD=t/testdata/simple_logfile
TMPD=$(mktemp -d)
FILTERS="$TD/procrastinationedr.specs"
AUDITLF="$TD/audit.log"
OUTFILE="$TMPD/out"
LOGFILE="$TMPD/log"
trap 'rm -rf "$TMPD"; tap_done_testing' EXIT

# Run logging both to stdout and a file.
go run . \
        -audit-log     "$AUDITLF" \
        -filters       "$FILTERS" \
        -log           "$LOGFILE" \
        -max-runtime    1m        \
        -no-tail                  \
        -no-timestamps \
        <"$AUDITLF" >&2
RET=$?
tap_is "$RET" 0 "Ran happily with logfile" "$0" $LINENO
if [[ 0 -ne "$RET" ]]; then
        tap_diag "Log:"
        tap_diag "$(<"$LOGFILE")"
fi
go run . \
        -audit-log     "$AUDITLF" \
        -filters       "$FILTERS" \
        -max-runtime    1m        \
        -no-tail                  \
        -no-timestamps            \
        >$OUTFILE
RET=$?
tap_is "$RET" 0 "Ran happily with stdout" "$0" $LINENO
if [[ 0 -ne "$RET" ]]; then
        tap_diag "Log:"
        tap_diag "$(<"$OUTFILE")"
fi

# Make sure the logfile gets the same output as we'd get on stdout.
GOT=$(($(wc -c <"$OUTFILE")))
tap_isnt "$GOT" 0 "Got output" "$0" $LINENO
GOT=$(diff -u "$OUTFILE" "$LOGFILE" ||:)
tap_is "$GOT" "" "No difference between stdout and logfile" "$0" $LINENO

wait

# vim: ft=sh
