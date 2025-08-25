procrastinationed~r~
====================
Quite possibly the world's worst Linux ED~R~.  Really just reformats and
filters `auditd(8)`'s `PROCTITLE` lines.  Probably should have been a
[Perl](#Perl) one-liner.

There is no R.  In practice, it's an ED but somewhat less capable than `ed(1)`.

The general idea is to have `auditd(8)` log process executions and
procrastinationedr filter and print them.  It watches `audit.log` with
`tail -f`(ish).

"Features"
----------
- Confugurable event priorities
- Uses `auditd(8)`
- Easy configuration
- Easy deployment
- Integrates with any logging solution which can read its stdout
- No messing about with `argv`'s vectorness

Quickstart
----------
1. Have Go [installed](./https://go.dev/doc/install)
1. Have `auditd(8)` installed and running.
   ```sh
   apt install auditd # Or you distro of choice's equivalent
   ```
2. Tell auditd to log process executions
   ```sh
   auditctl -a always,exit -F arch=b64 -S execve -F success=1 -F auid!=0
   ```
3. Set some logging specifications
   ```sh
   cat >/etc/procrastinationedr.specs <<_eof
   10 ^rm -[rf]{2,}
   10 /etc/shadow
   5  ^sudo
   5  emacs
   _eof
   ```
4. Build the code and start watching for scary processes
   ```sh
   go install github.com/magisterquis/procrastinationedr@latest
   procrastinationedr -h # For just in case
   procrastinationedr
   ```
5. Do malicious things
   ```sh
   cat /etc/shadow             # Steal hashes
   rm -fr /*                   # Remove the French langauge pack
   emacs /root/.ssh/id_ed25519 # Oh dear.
   ```
6. Observe the output
   ```
   10 cat /etc/shadow
   10 rm -fr /bin /boot /dev /etc /home /lib /lib64 /lost+found /media /mnt /opt /proc /root /run /sbin /srv /sys /tmp /usr /var
   5 emacs /root/.ssh/id_ed25519
   5 /usr/bin/emacs -no-comp-spawn --batch -l /tmp/emacs-async-comp-debian-startup-ySq3rw.el
   5 /usr/bin/emacs -no-comp-spawn --batch -l /tmp/emacs-async-comp-time-date.el-he29zT.el
   ```

Usage
-----
```
Usage: procrastinationedr [options]

Quite possibly the world's worst Linux ED!R.

The filter specifications file should consist of lines with a numerical
priority and a regular expression matching a process's space-joined argv, e.g.

10 ^rm -[rf]{2,}
10 /etc/shadow
5  ^sudo
5  emacs

Options:
  -audit-log logfile
    	Audit logfile (default "/var/log/audit/audit.log")
  -filters file
    	Filter specifications file (default "/etc/procrastinationedr.specs")
  -max-runtime duration
    	Optional maximum run duration
  -no-tail
    	Stop upon reading the end of the logfile
  -no-timestamps
    	Don't print timestamps
```

Config
------
The only real configuration is via the filter specifications file, which
consists of prioritized regular expressions used to decide whether to print
program executions logged by `auditd(8)`.

The priority is an integer.  Higher numbers take priority.

The regular expression is, well, a regular expression.  Each program's argument
vector will be joined with spaces and checkcked against each regular expression
in decreasing numerical order.  Should one match, it will be logged.  See the
quickstart for an example.

Perl
----
This whole project could have been the below with a bit of `egrep(1)`.
```sh
sudo tail -f /var/log/audit/audit.log | perl -pe 's/.*proctitle=\"(.*)\"/\1/||s/.*proctitle=(.*)/pack"H*",$1/e||s/.*//s;;s/\0/ /g;'
```
