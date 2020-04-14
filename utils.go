package main

func endsWithSlash(s string) string {
	if s == "" || s[len(s)-1] == '/' {
		return s
	}
	return s + "/"
}
