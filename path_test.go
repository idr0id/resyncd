package main

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_endsWithSlash(t *testing.T) {
	assert.Equal(t, "", endsWithSlash(""))
	assert.Equal(t, "/", endsWithSlash("/"))
	assert.Equal(t, "x/", endsWithSlash("x"))
	assert.Equal(t, "x/", endsWithSlash("x/"))
	assert.Equal(t, "x/y/", endsWithSlash("x/y"))
	assert.Equal(t, "x/y/", endsWithSlash("x/y/"))
}

func Test_expandPaths(t *testing.T) {
	tests := []struct {
		paths []string
		want  []string
	}{
		{
			paths: []string{""},
			want:  nil,
		},
		{
			paths: []string{"/"},
			want:  []string{"/"},
		},
		{
			paths: []string{"a"},
			want:  []string{"a"},
		},
		{
			paths: []string{"a/b/c", "a/b/d"},
			want:  []string{"a", "a/b", "a/b/c", "a/b/d"},
		},
		{
			paths: []string{"/a/b"},
			want:  []string{"/", "/a", "/a/b"},
		},
		{
			paths: []string{"/a/b/"},
			want:  []string{"/", "/a", "/a/b"},
		},
		{
			paths: []string{"/a/b", "/a/c"},
			want:  []string{"/", "/a", "/a/b", "/a/c"},
		},
	}

	for _, tc := range tests {
		expanded := expandPaths(tc.paths)
		slices.Sort(expanded)

		assert.Equal(t, tc.want, expanded)
	}
}
