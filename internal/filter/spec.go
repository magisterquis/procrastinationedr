package filter

/*
 * spec.go
 * Parse filter specifications
 * By J. Stuart McMurray
 * Created 20250824
 * Last Modified 20250824
 */

import (
	"bytes"
	"cmp"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// SpecRE extracts the level and regex from a filter line.
var SpecRE = regexp.MustCompile(`^\s*(\d+)\s+(.*)`)

// RegexSet is a set of regular expressions at a given level.
type RegexSet struct {
	Level uint
	REs   []*regexp.Regexp
}

// Specs holds a set of RegexSets at varying levels.
type Specs struct {
	Sets []RegexSet /* Sorted by Level. */
}

// ParseSpecsFromFile is a convenience function to slurp a file and pass it
// to ParseSpecs.
func ParseSpecsFromFile(fn string) (Specs, error) {
	b, err := os.ReadFile(fn)
	if nil != err {
		return Specs{}, fmt.Errorf("reading file: %w", err)
	}
	spec, err := ParseSpecs(b)
	if nil != err {
		return Specs{}, fmt.Errorf("parsing specs: %w", err)
	}
	return spec, nil
}

// ParseSpecs parses a filter specification into a Spec.  The specification
// should contain one regex per line, preceeded by a severity number and a
// space, e.g.
//
//	10 ^rm -(rf|fr) /
//	20 ^pkill -[A-Z0-9]+
func ParseSpecs(spec []byte) (Specs, error) {
	sets := make(map[uint]RegexSet)

	for l := range bytes.Lines(spec) {
		l = bytes.TrimRight(l, "\r\n")
		/* Skip blanks and comments. */
		if 0 == len(l) || '#' == l[0] {
			continue
		}
		/* Extract the important bits. */
		level, re, err := parseSpecLine(l)
		if nil != err {
			return Specs{}, fmt.Errorf("parsing %q: %w", l, err)
		}
		/* Looks good, save it. */
		set := sets[level]
		set.Level = level /* Eh. */
		set.REs = append(set.REs, re)
		sets[level] = set
	}

	/* Save ALL the sets. */
	var ret Specs
	for _, set := range sets {
		/* sort -u */
		slices.SortFunc(set.REs, func(a, b *regexp.Regexp) int {
			return strings.Compare(a.String(), b.String())
		})
		set.REs = slices.Compact(set.REs)
		set.REs = slices.Clip(set.REs)
		/* Save this level. */
		ret.Sets = append(ret.Sets, set)
	}

	/* Sort by level. */
	slices.SortFunc(ret.Sets, func(a, b RegexSet) int {
		return cmp.Compare(a.Level, b.Level)
	})

	return ret, nil
}

// parseSpecLine extract the level and regex from a line in the spec.
func parseSpecLine(l []byte) (uint, *regexp.Regexp, error) {
	/* Get the good bits. */
	ms := SpecRE.FindSubmatch(l)
	if 3 != len(ms) {
		return 0, nil, fmt.Errorf("invalid line")
	}
	/* Level. */
	level, err := strconv.ParseUint(string(ms[1]), 0, 0)
	if nil != err {
		return 0, nil, fmt.Errorf("invalid level %q: %w", ms[1], err)
	}
	/* Regex. */
	re, err := regexp.Compile(string(ms[2]))
	if nil != err {
		return 0, nil, fmt.Errorf("invalid regex %q: %w", re, err)
	}

	return uint(level), re, nil
}
