package ext

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/extensionlist"
	"github.com/darthbanana13/artifact-selector/pkg/github"
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

var ExtensionContentType = map[string][]string{
	"deb":            {"application/octet-stream", "application/vnd.debian.binary-package", "application/x-debian-package"},
	"tar.gz":         {"application/gzip", "application/x-gtar", "application/x-gzip"},
	"zip":            {"application/zip"},
	"asc":            {"application/pgp-signature"},
	"txt":            {"text/plain"},
	"tar.xz":         {"application/x-xz"},
	"tar.bz2":        {"application/x-bzip2", "application/x-bzip"},
	"tbz":            {"application/x-bzip2", "application/x-bzip1-compressed-tar"},
	"tar.zst":        {"application/octet-stream"},
	"appimage":       {"application/vnd.appimage"},
	"appimage.zsync": {"application/octet-stream"},
	"sha256sum":      {"application/octet-stream"},
	"sha256":         {"application/octet-stream"},
	"json":           {"application/json"},
	"msi":            {"application/x-msi"},
	"exe":            {"application/x-msdownload"},
	"rpm":            {"application/x-rpm"},
	"apk":            {"application/vnd.android.package-archive"},
	"dmg":            {"application/x-apple-diskimage"},
	"pkg":            {"application/octet-stream"},
	"":								{"application/octet-stream"}, //TODO: Make a special value for this called LINUXBINARY
}

type Ext struct {
	TargetExts []string
}

// TODO: Handle errors for extensions with unknown content types
func NewExtFilter(targetExts []string) (IExt, error) {
	return &Ext{TargetExts: targetExts}, nil
}

func (e *Ext) SetTargetExts(targetExts []string) error {
	e.TargetExts = targetExts
	return nil
}

func (e *Ext) FilterArtifact(artifact github.Artifact) (github.Artifact, bool) {
	for _, ext := range e.TargetExts {
		if e.HasExtension(artifact, ext) {
			return artifact, true
		}
	}
	return artifact, false
}

func (e *Ext) HasExtension(artifact github.Artifact, ext string) bool {
	if ext == "" {
		return IsLinuxExecutable(artifact)
	}
	return DoesEndWithExtension(artifact.FileName, ext) && HasContentType(artifact.ContentType, ext)
}

func DoesEndWithExtension(fileName, ext string) bool {
	return strings.HasSuffix(strings.ToLower(fileName), fmt.Sprintf(".%s", ext))
}

func IsLinuxExecutable(artifact github.Artifact) bool {
	hasNoExtension := !strings.Contains(artifact.FileName, ".")
	// Maybe the author used . for separators in the filename but it actually has no extension
	hasKnownExtension := extensionlist.IsKnownExtension(filepath.Ext(artifact.FileName))
	hasLinuxExecutableContentType := HasContentType(artifact.ContentType, "")
	return (hasNoExtension || !hasKnownExtension) && hasLinuxExecutableContentType
}

func HasContentType(contentType string, ext string) bool {
	return slices.Contains(ExtensionContentType[ext], strings.ToLower(contentType))
}
