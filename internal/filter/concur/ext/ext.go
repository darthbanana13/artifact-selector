package ext

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/darthbanana13/artifact-selector/internal/extensionlist"
	"github.com/darthbanana13/artifact-selector/internal/filter"
)

// TODO: Theoretically we should be able to tell what extensions we'd like to
//  search for by default by the name of the OS/distro
// var osExtensionsMap = map[string][]string{
//   "windows":  {"msi", "exe"},
//   "macos":    {"dmg", "pkg"},
//   "android":  {"apk"},
//   "generic":  {"zip"},
//   "linux":    {"appimage", "tar.gz", "tar.xz", ""},
// }
//
// var LinuxDistroExtensions = map[string][]string{
//   "deb": {"debian", "ubuntu"},
//   "rpm": {"fedora", "redhat", "rhel"},
//   "apk": {"android"},
// }

const (
	LinuxBinary = "LINUXBINARY"
)

type Ext struct {
	TargetExts []string
}

func NewExt(targetExts []string) (IExt, error) {
	e := &Ext{}
	err := e.SetTargetExts(targetExts)
	return e, err
}

func (e *Ext) SetTargetExts(targetExts []string) error {
	e.TargetExts = targetExts
	return nil
}

func (e *Ext) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	for _, ext := range e.TargetExts {
		if HasExtension(artifact.FileName, artifact.ContentType, ext) {
			artifact.Metadata, _ = filter.AddMetadata(artifact.Metadata, "ext", ext)
			return artifact, true
		}
	}
	return artifact, false
}

func HasExtension(fileName, contentType, ext string) bool {
	if ext == LinuxBinary {
		return IsLinuxExecutable(fileName, contentType)
	}
	return DoesEndWithExtension(fileName, ext)
}

func DoesEndWithExtension(fileName, ext string) bool {
	return strings.HasSuffix(strings.ToLower(fileName), fmt.Sprintf(".%s", ext))
}

func IsLinuxExecutable(fileName, contentType string) bool {
	hasNoExtension := !strings.Contains(fileName, ".")
	// Maybe the author used . for separators in the filename but it actually has no extension
	hasKnownExtension := extensionlist.IsKnownExtension(filepath.Ext(fileName))
	return hasNoExtension || !hasKnownExtension
}
