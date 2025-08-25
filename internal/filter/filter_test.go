package filter

/*
 * filter_test.go
 * Tests for filter.go
 * By J. Stuart McMurray
 * Created 20250824
 * Last Modified 20250824
 */

import (
	"slices"
	"testing"
)

func TestFilter(t *testing.T) {
	var (
		specList = `
1 ^a.*
1 b.*
2 c.*
2 d.*
3 ^e.*
3 f.*
`
		have = [][]string{
			{"abc", "123"},
			{"xxx", "yyy", "zzz"},
			{"foo"},
			{"7et", "", "", "what"},
			{"d", "e", "ffffff"},
		}
		want = []Hit{
			{Level: 0x1, Args: "abc 123"},
			{Level: 0x3, Args: "foo"},
			{Level: 0x2, Args: "d e ffffff"},
		}
		argvs = make(chan []string, len(have))
		hits  = make(chan Hit, len(have))
	)

	specs, err := ParseSpecs([]byte(specList))
	if nil != err {
		t.Fatalf("Error parsing specs: %s", err)
	}

	for _, v := range have {
		argvs <- v
	}
	close(argvs)
	Filter(t.Context(), hits, argvs, specs)

	got := make([]Hit, 0, len(hits))
	for hit := range hits {
		got = append(got, hit)
	}

	if !slices.Equal(got, want) {
		t.Errorf(
			"Incorrect filtered hits\n got: %#v\nwant: %#v",
			got,
			want,
		)
	}
}
