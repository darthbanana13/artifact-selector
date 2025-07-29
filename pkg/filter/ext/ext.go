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
	//TODO: Make a special value for this called LINUXBINARY
	"": {"application/octet-stream"},
}

type ExtFilter struct {
	El         *extensionlist.ExtensionList
	TargetExts []string
}

// TODO: Handle errors for extensions with unknown content types
func NewOSFilter(targetExts []string) (*ExtFilter, error) {
	el, err := extensionlist.NewExtensionList()
	return &ExtFilter{El: el, TargetExts: targetExts}, err
}

// TODO: Refactor to a smaller version
func (ef *ExtFilter) Filter(releases github.ReleasesInfo) github.ReleasesInfo {
	var filteredArtifacts []github.Artifact

	for _, ext := range ef.TargetExts {
		for _, artifact := range releases.Artifacts {
			if ef.HasExtension(artifact, ext) {
				filteredArtifacts = append(filteredArtifacts, artifact)
			}
		}
	}
	releases.Artifacts = filteredArtifacts
	return releases
}

func (ef *ExtFilter) HasExtension(artifact github.Artifact, ext string) bool {
	if ext == "" {
		return ef.IsLinuxExecutable(artifact)
	}
	return ef.DoesEndWithExtension(artifact.FileName, ext) && ef.HasContentType(artifact.ContentType, ext)
}

func (ExtFilter) DoesEndWithExtension(fileName, ext string) bool {
	return strings.HasSuffix(strings.ToLower(fileName), fmt.Sprintf(".%s", ext))
}

func (ef *ExtFilter) IsLinuxExecutable(artifact github.Artifact) bool {
	hasNoExtension := !strings.Contains(artifact.FileName, ".")
	// Maybe the author used . for separators in the filename but it actually has no extension
	hasKnownExtension := ef.El.IsKnownExtension(filepath.Ext(artifact.FileName))
	hasLinuxExecutableContentType := ef.HasContentType(artifact.ContentType, "")
	return (hasNoExtension || !hasKnownExtension) && hasLinuxExecutableContentType
}

func (ExtFilter) HasContentType(contentType string, ext string) bool {
	return slices.Contains(ExtensionContentType[ext], strings.ToLower(contentType))
}
