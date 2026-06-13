package extos

import (
	"math"
	"slices"

	"github.com/darthbanana13/artifact-selector/internal/filter"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/ext"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os/decorator"
	"github.com/darthbanana13/artifact-selector/internal/funcdecorator"
)

// NOTE: technically, the tar extensions would work for any *nix system however
// the most likely answer if the OS string was found in the filename is that
// it is for a Linux distro. If this changes in the future, please update the
// map.
var ExtOSMap = map[string]string{
	"deb":			"debian",
	"rpm":			"rhel",
	"appimage":	"linux",
	"tar.zst":	"linux",
	"tar.gz":		"linux",
	"tar.xz":		"linux",
	"tbz":			"linux",
	"apk":			"android",
	"exe":			"windows",
	"dmg":			"macos",
}

type ExtOsDecorator struct {
	os.IOS
}

func ExtOsConstructorDecorator() funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(ofc decorator.Constructor) decorator.Constructor {
		return func(targetOS string) (os.IOS, error) {
			of, err := ofc(targetOS)
			if err != nil {
				return of, err
			}
			of, err = NewExtOSDecorator(of)
			if err != nil {
				return of, err
			}
			return of, nil
		}
	}
}

func NewExtOSDecorator(os os.IOS) (os.IOS, error) {
	return &ExtOsDecorator{
		IOS: os,
	}, nil
}

func (eod *ExtOsDecorator) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	filteredArtifact, keep := eod.IOS.FilterArtifact(artifact)
	osAlias, ok := ExtOSMap[artifact.GetStringMetadata(ext.MetadataKey)]
	if !ok {
		return filteredArtifact, keep
	}
	foundOSIndex, ok := filteredArtifact.Metadata[os.MetadataOSIndexKey]
	if !ok {
		foundOSIndex = math.MaxInt
	}
	targetAliases := eod.IOS.TargetAliases()
	if i := slices.Index(targetAliases, osAlias); i >= 0 && len(targetAliases) - i > foundOSIndex.(int) {
		filteredArtifact.AddMetadata(os.MetadataOSNameKey, osAlias, os.MetadataOSIndexKey, len(targetAliases)-i)
		return filteredArtifact, true
	}
	return filteredArtifact, keep
}
