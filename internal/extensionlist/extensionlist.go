package extensionlist

import (
	"strings"
)

func IsKnownExtension(ext string) bool {
	if ext == "" || ext == "." {
		return true
	}
	ext = strings.ToLower(ext)
	ext = strings.TrimPrefix(ext, ".")
	if _, ok := exts[ext]; ok {
		return true
	}
	return false
}
