package main

func endsWithSlash(s string) string {
	if s == "" || s[len(s)-1] == '/' {
		return s
	}
	return s + "/"
}

func expandPaths(paths []string) []string {
	var list []string
	for _, path := range paths {
		size := len(path)
		if size == 0 {
			continue
		}

		current := ""
		for i := 0; i < size; i++ {
			if path[i] == '/' && i > 0 {
				list = append(list, current)
			}
			current += string(path[i])
			if (path[i] == '/' && i == 0) || i == size-1 {
				list = append(list, current)
			}
		}
	}
	return unique(list)
}

func unique(s []string) []string {
	if len(s) == 0 {
		return []string{}
	}

	keys := make(map[string]bool, len(s))
	var list []string
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}