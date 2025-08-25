package tail

/*
 * cat.go
 * Don't tail -f, just cat
 * By J. Stuart McMurray
 * Created 20250825
 * Last Modified 20250825
 */

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
)

// Cat reads the lines from fn into ch./ Cat will close ch before returning.
func Cat(ctx context.Context, ch chan<- string, fn string) error {
	defer close(ch)
	f, err := os.Open(fn)
	if nil != err {
		return fmt.Errorf("opening %s: %w", fn, err)
	}
	defer f.Close()
	defer context.AfterFunc(ctx, func() { f.Close() })

	/* Stop reading when the context is done. */
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ch <- scanner.Text()
	}
	if err := scanner.Err(); nil != err &&
		!(errors.Is(err, syscall.EBADF) && nil != ctx.Err()) {
		return fmt.Errorf("reading %s: %w", fn, err)
	}
	return nil
}
