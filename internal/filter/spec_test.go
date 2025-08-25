package filter

/*
 * spec_test.go
 * Parse filter specifications
 * By J. Stuart McMurray
 * Created 20250824
 * Last Modified 20250824
 */

import (
	"reflect"
	"regexp"
	"testing"
)

func TestParseSpec(t *testing.T) {
	have := `
1 a.*
2 c.*

# A comment
1 b.*
2 e.*



2 d.*


`
	want := Specs{
		Sets: []RegexSet{{
			Level: 1,
			REs: []*regexp.Regexp{
				regexp.MustCompile(`a.*`),
				regexp.MustCompile(`b.*`),
			},
		}, {
			Level: 2,
			REs: []*regexp.Regexp{
				regexp.MustCompile(`c.*`),
				regexp.MustCompile(`d.*`),
				regexp.MustCompile(`e.*`),
			},
		}},
	}
	got, err := ParseSpecs([]byte(have))
	if nil != err {
		t.Fatalf("Error: %s", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"ParseSpec failed\n"+
				"have:\n%s\n"+
				" got: %#v\n"+
				"want: %#v",
			have,
			got,
			want,
		)
	}
}
