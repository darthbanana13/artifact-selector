package ext

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/extensionlist"
	"github.com/darthbanana13/artifact-selector/pkg/filter"
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
	LINUXBINARY = ""
)

var ExtensionContentType = map[string][]string{
	"deb":            {"application/octet-stream", "application/vnd.debian.binary-package", "application/x-debian-package"},
	"tar.gz":         {"application/gzip", "application/x-gtar", "application/x-gzip"},
	"zip":            {"application/zip"},
	"asc":            {"application/pgp-signature"},
	"txt":            {"text/plain"},
	"sh":             {"application/x-sh"},
	"sig":            {"application/pgp-signature"},
	"minisig":        {"application/octet-stream"},
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
	LINUXBINARY:      {"application/octet-stream"},
}

type Ext struct {
	TargetExts []string
}

func NewExtFilter(targetExts []string) (IExt, error) {
	return &Ext{TargetExts: targetExts}, nil
}

func (e *Ext) SetTargetExts(targetExts []string) error {
	e.TargetExts = targetExts
	return nil
}

func (e *Ext) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	for _, ext := range e.TargetExts {
		if HasExtension(artifact.Source.FileName, artifact.Source.ContentType, ext) {
			artifact.Metadata["ext"] = ext
			return artifact, true
		}
	}
	return artifact, false
}

func HasExtension(fileName, contentType, ext string) bool {
	if ext == LINUXBINARY {
		return IsLinuxExecutable(fileName, contentType)
	}
	return DoesEndWithExtension(fileName, ext) && HasContentType(contentType, ext)
}

func DoesEndWithExtension(fileName, ext string) bool {
	return strings.HasSuffix(strings.ToLower(fileName), fmt.Sprintf(".%s", ext))
}

func IsLinuxExecutable(fileName, contentType string) bool {
	hasNoExtension := !strings.Contains(fileName, ".")
	// Maybe the author used . for separators in the filename but it actually has no extension
	hasKnownExtension := extensionlist.IsKnownExtension(filepath.Ext(fileName))
	hasLinuxExecutableContentType := HasContentType(contentType, LINUXBINARY)
	return (hasNoExtension || !hasKnownExtension) && hasLinuxExecutableContentType
}

func HasContentType(contentType string, ext string) bool {
	return slices.Contains(ExtensionContentType[ext], strings.ToLower(contentType))
}
