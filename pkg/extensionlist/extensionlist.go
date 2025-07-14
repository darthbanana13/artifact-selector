package extensionlist

import (
	"fmt"
	"os"
	"strings"
	// "unicode"
)

type ExtensionList struct {
	exts map[string]bool
}

func NewExtensionList() (*ExtensionList, error) {	
	el := &ExtensionList{}
	dat, err := os.ReadFile("pkg/extensionlist/extensionlist.txt")
	if err != nil {
		return el, err
	}
	lines := strings.Split(string(dat), "\n")
	el.exts = make(map[string]bool, len(lines))

	for _, ext := range lines {
		el.exts[ext] = true
	}
	return el, nil
}

func (el *ExtensionList) IsKnownExtension(ext string) bool {
	if ext == "" || ext == "." {
		return true
	}
	ext = strings.ToLower(ext)
	if !strings.HasPrefix(ext, ".") {
		ext = fmt.Sprintf(".%s", ext)
	}
	if _, ok := el.exts[ext]; ok {
		return true
	}
	return false
}
