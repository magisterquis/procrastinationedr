package event

/*
 * event_test.go
 * Tests for event.go
 * By J. Stuart McMurray
 * Created 20250823
 * Last Modified 20250823
 */

import (
	"context"
	"slices"
	"testing"
)

func TestParseAuditLines(t *testing.T) {
	var (
		have = []string{
			`			type=EXECVE msg=audit(1755989505.288:419): argc=2 a0="id" a1="-u"`,
			`type=CWD msg=audit(1755989505.288:419): cwd="/home/wurzel"`,
			`type=PATH msg=audit(1755989505.288:419): item=0 name="/usr/bin/id" inode=719 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PATH msg=audit(1755989505.288:419): item=1 name="/usr/bin/id" inode=719 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PATH msg=audit(1755989505.288:419): item=2 name="/lib64/ld-linux-x86-64.so.2" inode=2058 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PROCTITLE msg=audit(1755989505.288:419): proctitle=6964002D75`,
			`type=SYSCALL msg=audit(1755989505.304:420): arch=c000003e syscall=59 success=yes exit=0 a0=5596bbb85e50 a1=5596bbb85ad0 a2=5596bbb8cbc0 a3=e445143c66797f60 items=3 ppid=1390 pid=1392 auid=1000 uid=1000 gid=1000 euid=1000 suid=1000 fsuid=1000 egid=1000 sgid=1000 fsgid=1000 tty=(none) ses=1 comm="echo" exe="/usr/bin/echo" subj=unconfined key=(null)ARCH=x86_64 SYSCALL=execve AUID="wurzel" UID="wurzel" GID="wurzel" EUID="wurzel" SUID="wurzel" FSUID="wurzel" EGID="wurzel" SGID="wurzel" FSGID="wurzel"`,
			`type=EXECVE msg=audit(1755989505.304:420): argc=2 a0="/bin/echo" a1=7878227979`,
			`type=CWD msg=audit(1755989505.304:420): cwd="/home/wurzel"`,
			`type=PATH msg=audit(1755989505.304:420): item=0 name="/bin/echo" inode=644 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PATH msg=audit(1755989505.304:420): item=1 name="/bin/echo" inode=644 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PATH msg=audit(1755989505.304:420): item=2 name="/lib64/ld-linux-x86-64.so.2" inode=2058 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PROCTITLE msg=audit(1755989505.304:420): proctitle=2F62696E2F6563686F007878227979`,
			`type=USER_END msg=audit(1755990125.031:421): pid=1374 uid=1000 auid=1000 ses=1 subj=unconfined msg='op=PAM:session_close grantors=pam_limits,pam_permit,pam_unix acct="root" exe="/usr/bin/sudo" hostname=? addr=? terminal=/dev/pts/0 res=success'UID="wurzel" AUID="wurzel"`,
			`type=CRED_DISP msg=audit(1755990125.031:422): pid=1374 uid=1000 auid=1000 ses=1 subj=unconfined msg='op=PAM:setcred grantors=pam_permit acct="root" exe="/usr/bin/sudo" hostname=? addr=? terminal=/dev/pts/0 res=success'UID="wurzel" AUID="wurzel"`,
			`type=SYSCALL msg=audit(1755990136.219:423): arch=c000003e syscall=59 success=yes exit=0 a0=5622ac481010 a1=5622ac39a920 a2=5622ac38ec20 a3=3995628483f56db0 items=3 ppid=1367 pid=1403 auid=1000 uid=1000 gid=1000 euid=0 suid=0 fsuid=0 egid=1000 sgid=1000 fsgid=1000 tty=pts0 ses=1 comm="sudo" exe="/usr/bin/sudo" subj=unconfined key=(null)ARCH=x86_64 SYSCALL=execve AUID="wurzel" UID="wurzel" GID="wurzel" EUID="root" SUID="root" FSUID="root" EGID="wurzel" SGID="wurzel" FSGID="wurzel"`,
			`--`,
			`type=PATH msg=audit(1755990136.219:423): item=0 name="/usr/bin/sudo" inode=10634 dev=fe:01 mode=0104755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PATH msg=audit(1755990136.219:423): item=1 name="/usr/bin/sudo" inode=10634 dev=fe:01 mode=0104755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PATH msg=audit(1755990136.219:423): item=2 name="/lib64/ld-linux-x86-64.so.2" inode=2058 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PROCTITLE msg=audit(1755990136.219:423): proctitle=7375646F007461696C002D6E00313030002F7661722F6C6F672F61756469742F61756469742E6C6F67`,
			`type=USER_ACCT msg=audit(1755990136.223:424): pid=1403 uid=1000 auid=1000 ses=1 subj=unconfined msg='op=PAM:accounting grantors=pam_permit acct="wurzel" exe="/usr/bin/sudo" hostname=? addr=? terminal=/dev/pts/0 res=success'UID="wurzel" AUID="wurzel"`,
			`type=USER_CMD msg=audit(1755990136.223:425): pid=1403 uid=1000 auid=1000 ses=1 subj=unconfined msg='cwd="/home/wurzel" cmd=7461696C202D6E20313030202F7661722F6C6F672F61756469742F61756469742E6C6F67 exe="/usr/bin/sudo" terminal=pts/0 res=success'UID="wurzel" AUID="wurzel"`,
			`type=CRED_REFR msg=audit(1755990136.223:426): pid=1403 uid=1000 auid=1000 ses=1 subj=unconfined msg='op=PAM:setcred grantors=pam_permit acct="root" exe="/usr/bin/sudo" hostname=? addr=? terminal=/dev/pts/0 res=success'UID="wurzel" AUID="wurzel"`,
			`--`,
			`type=PATH msg=audit(1755990136.227:428): item=0 name="/usr/bin/tail" inode=758 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PATH msg=audit(1755990136.227:428): item=1 name="/usr/bin/tail" inode=758 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PATH msg=audit(1755990136.227:428): item=2 name="/lib64/ld-linux-x86-64.so.2" inode=2058 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PROCTITLE msg=audit(1755990136.227:428): proctitle=7461696C002D6E00313030002F7661722F6C6F672F61756469742F61756469742E6C6F67`,
			`type=USER_END msg=audit(1755990136.239:429): pid=1403 uid=1000 auid=1000 ses=1 subj=unconfined msg='op=PAM:session_close grantors=pam_limits,pam_permit,pam_unix acct="root" exe="/usr/bin/sudo" hostname=? addr=? terminal=/dev/pts/0 res=success'UID="wurzel" AUID="wurzel"`,
			`type=CRED_DISP msg=audit(1755990136.239:430): pid=1403 uid=1000 auid=1000 ses=1 subj=unconfined msg='op=PAM:setcred grantors=pam_permit acct="root" exe="/usr/bin/sudo" hostname=? addr=? terminal=/dev/pts/0 res=success'UID="wurzel" AUID="wurzel"`,
			`type=SYSCALL msg=audit(1755990152.207:431): arch=c000003e syscall=59 success=yes exit=0 a0=5622ac37fb90 a1=5622ac39c6b0 a2=5622ac38ec20 a3=8 items=3 ppid=1367 pid=1407 auid=1000 uid=1000 gid=1000 euid=1000 suid=1000 fsuid=1000 egid=1000 sgid=1000 fsgid=1000 tty=pts0 ses=1 comm="grep" exe="/usr/bin/grep" subj=unconfined key=(null)ARCH=x86_64 SYSCALL=execve AUID="wurzel" UID="wurzel" GID="wurzel" EUID="wurzel" SUID="wurzel" FSUID="wurzel" EGID="wurzel" SGID="wurzel" FSGID="wurzel"`,
			`type=SYSCALL msg=audit(1755990152.207:432): arch=c000003e syscall=59 success=yes exit=0 a0=5622ac4814e0 a1=5622ac480eb0 a2=5622ac38ec20 a3=8 items=3 ppid=1367 pid=1406 auid=1000 uid=1000 gid=1000 euid=0 suid=0 fsuid=0 egid=1000 sgid=1000 fsgid=1000 tty=pts0 ses=1 comm="sudo" exe="/usr/bin/sudo" subj=unconfined key=(null)ARCH=x86_64 SYSCALL=execve AUID="wurzel" UID="wurzel" GID="wurzel" EUID="root" SUID="root" FSUID="root" EGID="wurzel" SGID="wurzel" FSGID="wurzel"`,
			`type=EXECVE msg=audit(1755990152.207:431): argc=4 a0="grep" a1="-C" a2="3" a3="PROCTITLE"`,
			`type=CWD msg=audit(1755990152.207:431): cwd="/home/wurzel"`,
			`type=PATH msg=audit(1755990152.207:431): item=0 name="/usr/bin/grep" inode=1736 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PATH msg=audit(1755990152.207:431): item=1 name="/usr/bin/grep" inode=1736 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PATH msg=audit(1755990152.207:431): item=2 name="/lib64/ld-linux-x86-64.so.2" inode=2058 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PROCTITLE msg=audit(1755990152.207:431): proctitle=67726570002D4300330050524F435449544C45`,
			`type=EXECVE msg=audit(1755990152.207:432): argc=5 a0="sudo" a1="tail" a2="-n" a3="100" a4="/var/log/audit/audit.log"`,
			`type=CWD msg=audit(1755990152.207:432): cwd="/home/wurzel"`,
			`type=PATH msg=audit(1755990152.207:432): item=0 name="/usr/bin/sudo" inode=10634 dev=fe:01 mode=0104755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PATH msg=audit(1755990152.207:432): item=1 name="/usr/bin/sudo" inode=10634 dev=fe:01 mode=0104755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PATH msg=audit(1755990152.207:432): item=2 name="/lib64/ld-linux-x86-64.so.2" inode=2058 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PROCTITLE msg=audit(1755990152.207:432): proctitle=7375646F007461696C002D6E00313030002F7661722F6C6F672F61756469742F61756469742E6C6F67`,
			`type=USER_ACCT msg=audit(1755990152.211:433): pid=1406 uid=1000 auid=1000 ses=1 subj=unconfined msg='op=PAM:accounting grantors=pam_permit acct="wurzel" exe="/usr/bin/sudo" hostname=? addr=? terminal=/dev/pts/0 res=success'UID="wurzel" AUID="wurzel"`,
			`type=USER_CMD msg=audit(1755990152.211:434): pid=1406 uid=1000 auid=1000 ses=1 subj=unconfined msg='cwd="/home/wurzel" cmd=7461696C202D6E20313030202F7661722F6C6F672F61756469742F61756469742E6C6F67 exe="/usr/bin/sudo" terminal=pts/0 res=success'UID="wurzel" AUID="wurzel"`,
			`type=CRED_REFR msg=audit(1755990152.215:435): pid=1406 uid=1000 auid=1000 ses=1 subj=unconfined msg='op=PAM:setcred grantors=pam_permit acct="root" exe="/usr/bin/sudo" hostname=? addr=? terminal=/dev/pts/0 res=success'UID="wurzel" AUID="wurzel"`,
			`--`,
			`type=PATH msg=audit(1755990152.215:437): item=0 name="/usr/bin/tail" inode=758 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PATH msg=audit(1755990152.215:437): item=1 name="/usr/bin/tail" inode=758 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PATH msg=audit(1755990152.215:437): item=2 name="/lib64/ld-linux-x86-64.so.2" inode=2058 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
			`type=PROCTITLE msg=audit(1755990152.215:437): proctitle=7461696C002D6E00313030002F7661722F6C6F672F61756469742F61756469742E6C6F67`,
			`type=PROCTITLE msg=audit(1755977624.084:77): proctitle="(systemd)"`,
			`type=PROCTITLE msg=audit(1755977714.039:122): proctitle="id"`,
		}
		want = [][]string{
			[]string{"id", "-u"},
			[]string{"/bin/echo", "xx\"yy"},
			[]string{"sudo", "tail", "-n", "100", "/var/log/audit/audit.log"},
			[]string{"tail", "-n", "100", "/var/log/audit/audit.log"},
			[]string{"grep", "-C", "3", "PROCTITLE"},
			[]string{"sudo", "tail", "-n", "100", "/var/log/audit/audit.log"},
			[]string{"tail", "-n", "100", "/var/log/audit/audit.log"},
			[]string{"(systemd)"},
			[]string{"id"},
		}
		lch         = make(chan string, len(have))
		ach         = make(chan []string, len(have))
		ctx, cancel = context.WithCancel(t.Context())
	)
	defer cancel()
	for _, l := range have {
		lch <- l
	}
	close(lch)

	/* Can we parse a bunch of lines? */
	if err := ParseAuditLines(ctx, ach, lch); nil != err {
		t.Fatalf("Error: %s", err)
	}

	/* Did it work? */
	var got [][]string
	for a := range ach {
		got = append(got, a)
	}
	if len(got) != len(want) {
		t.Errorf("Got %d argvs, expected %d", len(got), len(want))
	}
	for i := range min(len(got), len(want)) {
		if !slices.Equal(got[i], want[i]) {
			t.Errorf(
				"Argv at index %d incorrect\n"+
					" got: %q\n"+
					"want: %q",
				i,
				got,
				want,
			)
		}
	}
	if len(got) > len(want) {
		for _, v := range got[len(want):] {
			t.Errorf("Extra argv: %#v", v)
		}
	} else if len(want) > len(got) {
		for _, v := range want[len(got):] {
			t.Errorf("Missing argv: %q", v)
		}
	}
}

func TestHandleLine(t *testing.T) {
	for n, c := range map[string]struct {
		have string
		want []string
	}{"wrongtype": {
		have: `type=PATH msg=audit(1755989158.921:397): item=0 name="/usr/bin/tail" inode=758 dev=fe:01 mode=0100755 ouid=0 ogid=0 rdev=00:00 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0OUID="root" OGID="root"`,
	}, "hex": {
		have: `type=PROCTITLE msg=audit(1755989158.921:397): proctitle=7461696C002D66002F7661722F6C6F672F61756469742F61756469742E6C6F67`,
		want: []string{"tail", "-f", "/var/log/audit/audit.log"},
	}, "login shell": {
		have: `type=PROCTITLE msg=audit(1755989356.388:409): proctitle="-bash"`,
		want: []string{"-bash"},
	}, "no arguments": {
		have: `type=PROCTITLE msg=audit(1755989356.408:411): proctitle="id"`,
		want: []string{"id"},
	}, "double quotes": {
		have: `type=PROCTITLE msg=audit(1755989505.304:420): proctitle=2F62696E2F6563686F007878227979`,
		want: []string{"/bin/echo", `xx"yy`},
	}, "one argument": {
		have: `type=PROCTITLE msg=audit(1755989449.412:415): proctitle=6964002D75`,
		want: []string{"id", "-u"},
	}} {
		t.Run(n, func(t *testing.T) {
			var (
				ach         = make(chan []string, 1)
				ctx, cancel = context.WithCancel(t.Context())
			)
			defer cancel()

			/* Handle the line itself. */
			if err := handleLine(ctx, ach, c.have); nil != err {
				t.Fatalf("Error: %s", err)
			}
			close(ach) /* Would be done by ParseAuditLines. */

			/* Were we supposed to get argv? */
			if nil == c.want {
				for a := range ach {
					t.Errorf(
						"Expected no argv, got %q",
						a,
					)
				}
				return
			}

			/* Did we get the right one? */
			got, ok := <-ach
			if !ok {
				t.Fatalf("Got no argv")
			}
			if !slices.Equal(got, c.want) {
				t.Errorf(
					"Incorrect argv\n"+
						"have: %s\n"+
						" got: %q\n"+
						"want: %q",
					c.have,
					got,
					c.want,
				)
			}
		})
	}
}
