package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
	assert.Equal(
		t,
		[]string{},
		expandPaths([]string{""}),
	)
	assert.Equal(
		t,
		[]string{"/"},
		expandPaths([]string{"/"}),
	)
	assert.Equal(
		t,
		[]string{"a"},
		expandPaths([]string{"a"}),
	)
	assert.Equal(
		t,
		[]string{"a", "a/b", "a/b/c", "a/b/d"},
		expandPaths([]string{"a/b/c", "a/b/d"}),
	)
	assert.Equal(
		t,
		[]string{"/", "/a", "/a/b"},
		expandPaths([]string{"/a/b"}),
	)
	assert.Equal(
		t,
		[]string{"/", "/a", "/a/b", "/a/b/"},
		expandPaths([]string{"/a/b/"}),
	)
	assert.Equal(
		t,
		[]string{"/", "/a", "/a/b", "/a/c"},
		expandPaths([]string{"/a/b", "/a/c"}),
	)
}
