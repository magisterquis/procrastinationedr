package tail

/*
 * tail_test.go
 * Tests for tail.go
 * By J. Stuart McMurray
 * Created 20250823
 * Last Modified 20250823
 */

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
)

func TestTail(t *testing.T) {
	/* Do we get an error if the file is missing? */
	t.Run("no_file", func(t *testing.T) {
		/* Filename, but it doesn't exist. */
		fn := filepath.Join(t.TempDir(), "kittens")
		/* Tail should barf. */
		ch := make(chan string)
		if err := Tail(t.Context(), ch, fn); nil == err {
			t.Errorf("Did not get error with missing file")
		}
	})

	/* Can we cancel the context after a few lines are read? */
	t.Run("empty_file/cancel_after_read", func(t *testing.T) {
		var (
			ch          = make(chan string, 1024)
			ctx, cancel = context.WithCancel(t.Context())
			have        = []string{"foo", "bar", "tridge"}
			terr        error
			wg          sync.WaitGroup
		)
		defer cancel()
		/* File to tail. */
		fn := filepath.Join(t.TempDir(), "kittens")
		f, err := os.Create(fn)
		if nil != err {
			t.Fatalf("Error creating file %s: %s", f.Name(), err)
		}
		defer f.Close()
		/* Start the tail and send some lines. */
		wg.Go(func() { terr = Tail(ctx, ch, f.Name()) })
		wg.Go(func() {
			for _, l := range have {
				fmt.Fprintf(f, "%s\n", l)
			}
		})
		/* Should get the lines ok. */
		got := make([]string, 0, len(have))
		for range len(have) {
			got = append(got, <-ch)
		}
		if !slices.Equal(got, have) {
			t.Errorf(
				"Incorrect lines:\n got: %q\nwant: %q",
				got,
				have,
			)
		}
		/* Should have got a predictable error. */
		cancel()
		wg.Wait()
		if nil != terr {
			t.Errorf("Tail error: %s", err)
		}
	})

	/* Can we cancel the context after a few lines are read? */
	t.Run("not_empty_file", func(t *testing.T) {
		var (
			ch          = make(chan string, 1024)
			ctx, cancel = context.WithCancel(t.Context())
			existing    = []string{"kittens", "moose"}
			have        = []string{"foo", "bar", "tridge"}
			terr        error
			wg          sync.WaitGroup
		)
		defer cancel()
		/* File to tail. */
		fn := filepath.Join(t.TempDir(), "kittens")
		if err := os.WriteFile(
			fn,
			[]byte(strings.Join(existing, "\n")+"\n"),
			0600,
		); nil != err {
			t.Fatalf("Error creating file %s: %s", fn, err)
		}
		f, err := os.OpenFile(fn, os.O_WRONLY|os.O_APPEND, 0600)
		if nil != err {
			t.Fatalf("Error opening file %s: %s", fn, err)
		}
		defer f.Close()
		/* Start the tail and send some more lines. */
		wg.Go(func() { terr = Tail(ctx, ch, fn) })
		wg.Go(func() {
			for _, l := range have {
				fmt.Fprintf(f, "%s\n", l)
			}
		})
		/* Should get the lines ok. */
		got := make([]string, 0, len(have))
		want := append(
			slices.Clone(existing),
			have...,
		)
		for range len(want) {
			got = append(got, <-ch)
		}
		if !slices.Equal(got, want) {
			t.Errorf(
				"Incorrect lines:\n"+
					"existing: %q\n"+
					"   added: %q\n"+
					"     got: %q\n"+
					"    want: %q",
				existing,
				have,
				got,
				want,
			)
		}
		/* Shouldn't get an error cancelling the context. */
		cancel()
		wg.Wait()
		if nil != terr {
			t.Errorf("Tail error: %s", err)
		}
	})
}
