// Package event - Turn audit log lines into events
package event

/*
 * event.go
 * Turn audit log lines into events
 * By J. Stuart McMurray
 * Created 20250823
 * Last Modified 20250823
 */

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
)

// ParseAuditLines parses PROCTITLE lines from audit.log into argument vectors.
// argvs will be closed before ParseAuditLines returns.
func ParseAuditLines(
	ctx context.Context,
	argvs chan<- []string,
	lines <-chan string,
) error {
	defer close(argvs)

	for {
		select {
		case <-ctx.Done():
			return nil
		case line, ok := <-lines:
			if !ok {
				return nil
			}
			if err := handleLine(ctx, argvs, line); nil != err {
				return err
			}
		}
	}
}

// handleLine handles a single audit log line.
func handleLine(ctx context.Context, argvs chan<- []string, line string) error {
	/* Work Make sure we have an argument vector. */
	if !strings.HasPrefix(line, "type=PROCTITLE") {
		return nil
	}

	/* Get the important bit. */
	_, line, ok := strings.Cut(line, "proctitle=")
	if !ok {
		return fmt.Errorf("missing start of argv: %q", line)
	}

	/* A single argument is easy. */
	if strings.HasPrefix(line, `"`) && strings.HasSuffix(line, `"`) {
		line = strings.TrimPrefix(line, `"`)
		line = strings.TrimSuffix(line, `"`)
		select {
		case argvs <- []string{line}:
		case <-ctx.Done():
		}
		return nil
	}

	/* Should just be hex, then. */
	bs, err := hex.DecodeString(line)
	if nil != err {
		return fmt.Errorf("invalid hex argv: %q", line)
	}

	select {
	case argvs <- strings.Split(string(bs), "\x00"):
	case <-ctx.Done():
	}

	return nil
}
