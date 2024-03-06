package util

import "strings"

// SplitNameWithOwner split nameWithOwner string into owner and name string
// e.g. cloudwego/hertz => cloudwego hertz
func SplitNameWithOwner(s string) (string, string) {
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
