package main

import (
	"fmt"
	"regexp"
	"strings"
)

type matcher func(path string) bool
type matchers []matcher

func (m matchers) match(path string) bool {
	for _, matcher := range m {
		if matcher(path) {
			return true
		}
	}
	return false
}

func newFileMatcher(root string, pattern string) matcher {
	pattern = regexp.QuoteMeta(pattern)
	pattern = strings.ReplaceAll(pattern, "\\*\\*", ".*")
	pattern = strings.ReplaceAll(pattern, "\\*", "[^/]*")

	starting := ""
	if pattern[0] == '/' {
		starting = "^"
	}
	regex := regexp.MustCompile(
		fmt.Sprintf(`%s%s(/.+)?$`, starting, pattern))

	if root != "" && root[len(root)-1] != '/' {
		root += "/"
	}
	return func(path string) bool {
		if strings.Index(path, root) != 0 {
			return false
		}
		if root != "" && root != "/" {
			path = path[len(root)-1:]
		}
		return regex.MatchString(path)
	}
}

func newFileMatchers(path string, patterns []string) matchers {
	var list matchers
	for _, pattern := range patterns {
		list = append(list, newFileMatcher(path, pattern))
	}
	return list
}
