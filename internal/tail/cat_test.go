package tail

/*
 * cat_test.go
 * Tests for cat.go
 * By J. Stuart McMurray
 * Created 20250825
 * Last Modified 20250825
 */

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestCat(t *testing.T) {
	/* Test file. */
	have := []string{
		"foo",
		"bar",
		"tridge",
	}
	fn := filepath.Join(t.TempDir(), "kittens")
	if err := os.WriteFile(
		fn,
		[]byte(strings.Join(have, "\n")+"\n"),
		0600,
	); nil != err {
		t.Errorf("Error creating %s: %s", fn, err)
	}

	/* cat the file. */
	ch := make(chan string)
	eg, ctx := errgroup.WithContext(t.Context())
	eg.Go(func() error { return Cat(ctx, ch, fn) })

	/* Collect the lines. */
	var got []string
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case l, ok := <-ch:
				if !ok {
					return nil
				}
				got = append(got, l)
			}
		}
	})

	/* Did it work? */
	if err := eg.Wait(); nil != err {
		t.Fatalf("Error: %s", err)
	}
	if !slices.Equal(got, have) {
		t.Errorf("Incorrect lines read\n got: %q\nwant: %q", got, have)
	}
}
