package extensionlist

import (
	"fmt"
	"strings"
)

func IsKnownExtension(ext string) bool {
	if ext == "" || ext == "." {
		return true
	}
	ext = strings.ToLower(ext)
	if !strings.HasPrefix(ext, ".") {
		ext = fmt.Sprintf(".%s", ext)
	}
	if _, ok := exts[ext]; ok {
		return true
	}
	return false
}
