package utils

import "strings"

// Basename extracts the basename from a module name.
func Basename(module string) string {
	parts := strings.Split(module, ".")
	return parts[len(parts)-1]
}
