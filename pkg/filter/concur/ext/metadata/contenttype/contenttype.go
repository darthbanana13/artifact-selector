package contenttype

import (
	"slices"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext"
)

const (
	MISSING   = "missing"
	UNKNOWN   = "unknown"
	MISSMATCH = "missmatch"
	MATCH     = "match"
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
	ext.LINUXBINARY:  {"application/octet-stream"},
}

func FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	if artifact.ContentType == "" {
		artifact.Metadata = filter.AddMetadata(artifact.Metadata, "contentType", MISSING)
		return artifact, true
	}
	vals, ok := ExtensionContentType[artifact.Metadata["ext"].(string)]
	if !ok {
		artifact.Metadata = filter.AddMetadata(artifact.Metadata, "contentType", UNKNOWN)
	} else if !slices.Contains(vals, artifact.ContentType) {
		artifact.Metadata = filter.AddMetadata(artifact.Metadata, "contentType", MISSMATCH)
	} else {
		artifact.Metadata = filter.AddMetadata(artifact.Metadata, "contentType", MATCH)
	}
	return artifact, true
}
