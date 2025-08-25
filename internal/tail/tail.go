// Package tail - Tail audit.log
package tail

/*
 * tail.go
 * Tail audit.log
 * By J. Stuart McMurray
 * Created 20250823
 * Last Modified 20250823
 */

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"golang.org/x/sync/errgroup"
)

// Tail sends lines written to the file named fn to ch.  It just wraps tail(1).
// Tail will close ch before returning.
func Tail(ctx context.Context, ch chan<- string, fn string) error {
	defer close(ch)

	eg, ectx := errgroup.WithContext(ctx)

	/* Wrap tail(1). */
	cmd := tailCmd(ectx, fn)
	tout, err := cmd.StdoutPipe()
	if nil != err {
		return fmt.Errorf("setting up tail(1)'s output: %w", err)
	}
	ebuf := new(bytes.Buffer)
	cmd.Stderr = ebuf

	/* Start it going. */
	if err := cmd.Start(); nil != err {
		return fmt.Errorf("starting tail(1): %w", err)
	}

	/* Actually run tail(1). */
	eg.Go(func() error {
		if err := cmd.Wait(); nil != err && 0 != ebuf.Len() {
			return fmt.Errorf(
				"running tail(1) (stderr: %q): %w",
				ebuf,
				err,
			)
		} else if nil != err {
			return fmt.Errorf("running tail(1): %w", err)
		}
		return nil
	})

	/* Read lines into ch. */
	eg.Go(func() error {
		defer tout.Close()
		scanner := bufio.NewScanner(tout)
		/* Read a line.  Would be nice if we could cancel this. */
		for scanner.Scan() {
			/* Try to send it to ch. */
			select {
			case ch <- scanner.Text():
			case <-ectx.Done():
				return nil
			}
		}
		/* Process probably died. */
		if err := scanner.Err(); nil != err {
			return fmt.Errorf("reading tail(1)'s output: %w", err)
		}
		return nil
	})

	/* Bit of a hack, but there's no good way to know if the process died
	because the context was done. */
	if err := eg.Wait(); nil != err && nil == ctx.Err() {
		return err
	}
	return nil
}

// tailCmd returns an OS-suitable tail invocation.
func tailCmd(ctx context.Context, fn string) *exec.Cmd {
	args := make([]string, 0)

	/* Per-OS flags. */
	switch runtime.GOOS {
	case "linux":
		args = append(args, "--follow=name", "--silent")
	case "darwin":
		args = append(args, "tail", "-F", "--silent")
	default: /* Works for OpenBSD, at least? */
		args = append(args, "-f")
	}
	args = append(args, fn)

	return exec.CommandContext(ctx, "tail", args...)
}
