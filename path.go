package main

import (
	"path"
	"path/filepath"
)

func endsWithSlash(s string) string {
	if s == "" || s[len(s)-1] == '/' {
		return s
	}
	return s + "/"
}

func expandPaths(paths []string) []string {
	uniquePaths := map[string]struct{}{}
	for _, p := range paths {
		p = path.Clean(p)
		for ; p != "" && p != "." && p != "/"; p = filepath.Dir(p) {
			uniquePaths[p] = struct{}{}
		}
		if p == "/" {
			uniquePaths["/"] = struct{}{}
		}
	}

	if len(uniquePaths) == 0 {
		return nil
	}

	result := make([]string, 0, len(uniquePaths))
	for p := range uniquePaths {
		result = append(result, p)
	}
	return result
}
