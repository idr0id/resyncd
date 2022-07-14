package main

import "testing"

type TestCase struct {
	root            string
	pattern         string
	positiveMatches []string
	negativeMatches []string
}

func TestFileMatcher(t *testing.T) {
	// arrange
	testCases := []TestCase{
		{
			root:    "/",
			pattern: "foo",
			positiveMatches: []string{
				"/foo",
				"/foo/bar",
				"/foo/bar/baz",
				"/foobar/foo",
			},
			negativeMatches: []string{
				"/foobar",
				"/foobar/baz",
			},
		},
		{
			root:    "/",
			pattern: "*bar*",
			positiveMatches: []string{
				"/bar",
				"/foobar",
				"/barbaz",
				"/foobarbaz",
				"/foo.bar.baz",
				"/foo_bar_baz",
				"/foo/bar",
				"/foo/bar/bar",
				"/foo/bar/baz",
				"/foo/baz/bar",
				"/foo/foobarbaz",
			},
			negativeMatches: []string{
				"/foo",
				"/foo/baz/bat",
			},
		},
		{
			root:    "/",
			pattern: "/*bar*",
			positiveMatches: []string{
				"/bar",
				"/bar/baz",
			},
			negativeMatches: []string{
				"/foo",
				"/foo/bar",
				"/foo/bar/baz",
			},
		},
		{
			root:    "/",
			pattern: "bar.*",
			positiveMatches: []string{
				"/bar.",
				"/bar.baz",
				"/foo/bar.baz",
				"/foo/bar/bar.baz",
			},
			negativeMatches: []string{
				"/bar",
				"/foo/barbaz",
			},
		},
		{
			root:    "/",
			pattern: "**/baz",
			positiveMatches: []string{
				"/baz",
				"/foo/baz",
				"/foo/bar/baz",
			},
			negativeMatches: []string{
				"/",
			},
		},
		{
			root:    "/",
			pattern: "/foo",
			positiveMatches: []string{
				"/foo",
				"/foo/bar",
				"/foo/bar/baz",
			},
			negativeMatches: []string{
				"/",
				"/bar",
				"/bar/foo",
				"/bar/foo/baz",
				"/bar/baz/foo",
			},
		},
		{
			root:    "/foo",
			pattern: "**/baz",
			positiveMatches: []string{
				"/foo/baz",
				"/foo/bar/bat/baz",
				"/foo/bar/baz",
				"/foo/bar/baz/bat",
			},
			negativeMatches: []string{
				"/foo",
				"/baz",
			},
		},
	}

	for _, tc := range testCases {
		for _, path := range tc.positiveMatches {
			matcher := newFileMatcher(tc.root, tc.pattern)
			if matcher(path) == false {
				t.Errorf(
					"path `%s` should be matched (root: `%s`, pattern: `%s`)",
					path,
					tc.root,
					tc.pattern)
			}
		}

		for _, path := range tc.negativeMatches {
			matcher := newFileMatcher(tc.root, tc.pattern)
			if matcher(path) == true {
				t.Errorf(
					"path `%s` should not be matched (root: `%s`, pattern: `%s`)",
					path,
					tc.root,
					tc.pattern)
			}
		}
	}
}
