// Program procrastinationedr - Quite possibly the world's worst Linux ED!R
package main

/*
 * procrastinationedr.go
 * Quite possibly the world's worst Linux ED!R
 * By J. Stuart McMurray
 * Created 20250823
 * Last Modified 20250901
 */

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/magisterquis/procrastinationedr/internal/event"
	"github.com/magisterquis/procrastinationedr/internal/filter"
	"github.com/magisterquis/procrastinationedr/internal/tail"
	"golang.org/x/sync/errgroup"
)

func main() {
	/* Command-line flags. */
	var (
		auditLog = flag.String(
			"audit-log",
			"/var/log/audit/audit.log",
			"Audit `logfile`",
		)
		filterFile = flag.String(
			"filters",
			"/etc/procrastinationedr.specs",
			"Filter specifications `file`",
		)
		logFile = flag.String(
			"log",
			"",
			"Optional `file` to which to write logs",
		)
		maxRuntime = flag.Duration(
			"max-runtime",
			0,
			"Optional maximum run `duration`",
		)
		noTail = flag.Bool(
			"no-tail",
			false,
			"Stop upon reading the end of the logfile",
		)
		noTimestamps = flag.Bool(
			"no-timestamps",
			false,
			"Don't print timestamps",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %s [options]

Quite possibly the world's worst Linux ED!R.

The filter specifications file should consist of lines with a numerical
priority and a regular expression matching a process's space-joined argv, e.g.

10 ^rm -[rf]{2,}
10 /etc/shadow
5  ^sudo
5  emacs

Options:
`,
			filepath.Base(os.Args[0]),
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Stop at some point, if we're meant to. */
	ctx := context.Background()
	if 0 != *maxRuntime {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, *maxRuntime)
		defer cancel()
	}

	/* Set up logging. */
	lf := os.Stdout
	if "" != *logFile {
		var err error
		if lf, err = os.OpenFile(
			*logFile,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND,
			0644,
		); nil != err {
			log.Fatalf(
				"Error opening logfile %s: %s",
				*logFile,
				err,
			)
		}
		defer lf.Close()
	}
	log.SetOutput(lf)
	if *noTimestamps {
		log.SetFlags(0)
	}

	/* Work out our filters. */
	var specs filter.Specs
	if "" != *filterFile {
		var err error
		specs, err = filter.ParseSpecsFromFile(*filterFile)
		if nil != err {
			log.Fatalf(
				"Failed to read filter specifications "+
					"from %s: %s",
				*filterFile,
				err,
			)
		}
	}

	eg, ctx := errgroup.WithContext(ctx)

	/* Tail audit.log */
	lines := make(chan string)
	if !*noTail {
		eg.Go(func() error { return tail.Tail(ctx, lines, *auditLog) })
	} else {
		eg.Go(func() error { return tail.Cat(ctx, lines, *auditLog) })
	}

	/* Parse into argvs. */
	argvs := make(chan []string)
	eg.Go(func() error { return event.ParseAuditLines(ctx, argvs, lines) })

	/* Filter to just the interesting ones. */
	hits := make(chan filter.Hit)
	eg.Go(func() error {
		filter.Filter(ctx, hits, argvs, specs)
		return nil
	})

	/* And print the interesting ones. */
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case hit, ok := <-hits:
				if !ok {
					return nil
				}
				log.Printf("%d %s", hit.Level, hit.Args)
			}
		}
	})

	/* Wait for things to happen. */
	if err := eg.Wait(); nil != err {
		log.Fatalf("Fatal error: %s", err)
	}
}
