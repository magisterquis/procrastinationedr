// Package filter - Filter processes by regex
package filter

/*
 * filter.go
 * Filter processes by regex
 * By J. Stuart McMurray
 * Created 20250824
 * Last Modified 20250830
 */

import (
	"context"
	"strings"
)

// Hit contains a space-joined argv regex match.
type Hit struct {
	Level uint
	Args  string
}

// Filter sends hits from argv which match a spec to hits, which will be closed
// before Filter returns.
func Filter(
	ctx context.Context,
	hits chan<- Hit,
	argvs <-chan []string,
	specs Specs,
) {
	defer close(hits)

	for {
		select {
		case argv, ok := <-argvs:
			if !ok {
				return
			}
			filterArgv(ctx, hits, argv, specs)
		case <-ctx.Done():
			return
		}
	}
}

// filterArgv sends argv as a hit to hits if it matches a spec.
func filterArgv(
	ctx context.Context,
	hits chan<- Hit,
	argv []string,
	specs Specs,
) {
	/* Because we're EDR, we assume argv is a single string and that
	newlines are really spaces. */
	args := strings.Join(argv, " ")
	args = strings.ReplaceAll(args, "\n", " ")

	/* Check ALL the specs. */
	for _, set := range specs.Sets {
		for _, re := range set.REs {
			if !re.MatchString(args) {
				continue
			}
			select {
			case hits <- Hit{
				Level: set.Level,
				Args:  args,
			}:
			case <-ctx.Done():
			}
			return
		}
	}
}
