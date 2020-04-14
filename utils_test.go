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
