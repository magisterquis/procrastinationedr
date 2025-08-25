#!/bin/ksh
#
# readme.t
# Make sure the readme output is correct
# By J. Stuart McMurray
# Created 20250825
# Last Modified 20250825

set -euo pipefail

. t/shmore.subr

tap_plan 4

TD=t/testdata/readme
LOGFILE="$TD/audit.log"
PROCTITLES="$TD/proctitles"
SPECS="$TD/procrastinationedr.specs"
README=README.md

# Make sure the specs haven't changed.
GOT=$(<"$SPECS")
WANT=$(
        awk '/^   cat >\/etc\/procrastinationedr.specs/,/   _eof/' "$README" |
                (egrep -v '^   (cat >|_eof)' ||:) |
                sed 's/^   //'
)
tap_is "$GOT" "$WANT" "Test specs are the same as in the README" "$0" $LINENO

# Make sure we get the right events
GOT=$(go run . \
        -audit-log "$LOGFILE" \
        -filters "$SPECS"     \
        -no-tail              \
        -no-timestamps        \
        <"$LOGFILE"
)
WANT=$(
        ed -s "$README" <<'_eof'
/6. Observe the output/++
+,/```/-p
_eof
)
WANT=$(print -r "$WANT" | sed 's/^   //')
tap_is "$GOT" "$WANT" "Got correct events from procrastinationedr" "$0" $LINENO

# Make sure the Perl in the README works.
PERLOL=$(egrep -o 'perl -pe .*' "$README")
GOT=$(($(print -r "$PERLOL" | wc -l)))
tap_is "$GOT" 1 "Found the Perl one-liner" "$0" $LINENO
GOT=$(eval $PERLOL $LOGFILE)
WANT=$(<"$PROCTITLES")
tap_is "$GOT" "$WANT" "Got correct events from the Perl one-liner" "$0" $LINENO

# vim: ft=sh
